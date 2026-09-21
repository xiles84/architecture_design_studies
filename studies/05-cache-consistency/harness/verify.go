package main

import (
	"context"
	"fmt"
	"strings"

	"adsplatform/ports"
)

// Correctness gates timing (methodology 5). Nothing in this study reports a
// throughput before these checks pass, and every expectation is computed in Go from
// the generated dataset and the operation ledger -- never against another query,
// which would only prove the queries agree with each other.

// readPortalSQL is the study's one implementation of "read the portal view",
// shared by the correctness gate, the audit and the cache adapter's fill path. A
// separate gate implementation could drift from the measured one and then certify a
// payload the harness never actually reads.
//
// The two statements run inside ONE repeatable-read transaction, so the returned
// content is a state the database really held. Reading them independently could
// stitch two states into one that never existed, and the content hash would then
// describe an impossible value.
func readPortalSQL(ctx context.Context, db ports.DB, cat *catalogue, usesVersion bool, personID int64) (PortalContent, int64, error) {
	tx, err := db.Begin(ctx, ports.RepeatableRead)
	if err != nil {
		return PortalContent{}, 0, err
	}
	defer tx.Rollback(ctx)

	args, err := cat.args(sPortalPerson, map[string]any{"person_id": personID})
	if err != nil {
		return PortalContent{}, 0, err
	}
	rows, err := tx.Query(ctx, cat.stmt(sPortalPerson).SQL, args...)
	if err != nil {
		return PortalContent{}, 0, err
	}
	content, err := scanPortalPerson(rows, personID)
	if err != nil {
		return PortalContent{}, 0, err
	}

	args, err = cat.args(sPortalRecent, map[string]any{"person_id": personID})
	if err != nil {
		return PortalContent{}, 0, err
	}
	rows, err = tx.Query(ctx, cat.stmt(sPortalRecent).SQL, args...)
	if err != nil {
		return PortalContent{}, 0, err
	}
	recent, err := scanPortalRecent(rows)
	if err != nil {
		return PortalContent{}, 0, err
	}
	content.Recent = recent
	content.Normalize()

	var version int64
	if usesVersion {
		a, err := cat.args(sPersonVersion, map[string]any{"person_id": personID})
		if err != nil {
			return PortalContent{}, 0, err
		}
		if err := tx.QueryRow(ctx, cat.stmt(sPersonVersion).SQL, a...).Scan(&version); err != nil {
			return PortalContent{}, 0, err
		}
	}
	_ = tx.Rollback(ctx)
	return content, version, nil
}

// samplePeople spreads a sample across the donor list: the first, the last and an
// even stride. Checking only the first N would miss a boundary bug in the loader.
func samplePeople(ds *Dataset, sample int) []Person {
	ps := ds.People
	if sample <= 0 || sample >= len(ps) {
		return ps
	}
	stride := len(ps) / sample
	if stride < 1 {
		stride = 1
	}
	var out []Person
	for i := 0; i < len(ps) && len(out) < sample; i += stride {
		out = append(out, ps[i])
	}
	return out
}

// verifyReads is the pre-timing gate. It answers every read statement for a spread
// of donors and compares the result with the oracle, which is built from the
// generated dataset rather than from the database.
func verifyReads(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, orc *oracle, sample int) GateInfo {
	g := GateInfo{Passed: true}
	fail := func(format string, a ...any) {
		g.Passed = false
		if len(g.Failures) < 20 {
			g.Failures = append(g.Failures, fmt.Sprintf(format, a...))
		}
	}
	usesVersion := cat.design.usesVersion()

	for _, p := range samplePeople(ds, sample) {
		g.Checks++
		g.Statements = addOnce(g.Statements, sPortalPerson)
		g.Statements = addOnce(g.Statements, sPortalRecent)
		got, version, err := readPortalSQL(ctx, db, cat, usesVersion, p.ID)
		if err != nil {
			fail("portal view for donor %d: %v", p.ID, err)
			continue
		}
		want, ok := orc.ExpectedContent(p.ID)
		if !ok {
			fail("oracle has no donor %d", p.ID)
			continue
		}
		if got.ContentHash() != want.ContentHash() {
			fail("INV-1: donor %d portal content differs from the dataset-derived expectation (%s)", p.ID, firstDifferencePortal(want, got))
			continue
		}
		if usesVersion && version != 0 {
			fail("INV-2: a freshly loaded donor reports cache_version %d, expected 0", version)
		}
		// The recent slice's ORDER is part of the cached representation's contract,
		// so it is checked separately from the (order-insensitive) hash.
		if ids := recentIDs(got); !equalInt64(ids, orc.ExpectedRecentIDs(p.ID)) {
			fail("INV-3: donor %d recent slice order differs from the deterministic (donated_at DESC, donation_id DESC) order", p.ID)
		}
	}

	// The charity feed must execute and return rows; its correctness is the
	// rolldown comparison, checked by the audit's a_rolldown_drift.
	if len(ds.Charities) > 0 {
		ch := ds.Charities[0]
		g.Checks++
		g.Statements = addOnce(g.Statements, sCharityRecent)
		args, err := cat.args(sCharityRecent, map[string]any{"charity_id": ch.ID})
		if err != nil {
			fail("bind %s: %v", sCharityRecent, err)
		} else {
			rows, err := db.Query(ctx, cat.stmt(sCharityRecent).SQL, args...)
			if err != nil {
				fail("%s: %v", sCharityRecent, err)
			} else {
				n, err := drain(rows)
				if err != nil {
					fail("%s scan: %v", sCharityRecent, err)
				} else if n == 0 {
					fail("%s returned no rows for a loaded charity", sCharityRecent)
				}
			}
		}
	}

	// The owned model's token must start at a state the audits can check.
	if usesVersion {
		g.Checks++
		g.Statements = addOnce(g.Statements, sAuditOutboxOrphans)
		if n, err := scalarInt(ctx, db, cat, sAuditOutboxOrphans, nil); err != nil {
			fail("%s: %v", sAuditOutboxOrphans, err)
		} else if n != 0 {
			fail("INV-4: %s reports %d outbox events naming a donor that does not exist", sAuditOutboxOrphans, n)
		}
	}
	return g
}

// ---------------------------------------------------------------- helpers

func recentIDs(c PortalContent) []int64 {
	out := make([]int64, 0, len(c.Recent))
	for _, d := range c.Recent {
		out = append(out, d.ID)
	}
	return out
}

func equalInt64(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func firstDifferencePortal(want, got PortalContent) string {
	switch {
	case want.DonationCount != got.DonationCount:
		return fmt.Sprintf("count %d vs %d", want.DonationCount, got.DonationCount)
	case want.DonationTotalCents != got.DonationTotalCents:
		return fmt.Sprintf("total %d vs %d", want.DonationTotalCents, got.DonationTotalCents)
	case want.FullName != got.FullName:
		return fmt.Sprintf("full_name %q vs %q", want.FullName, got.FullName)
	case want.Email != got.Email:
		return fmt.Sprintf("email %q vs %q", want.Email, got.Email)
	case len(want.Recent) != len(got.Recent):
		return fmt.Sprintf("recent slice %d vs %d elements", len(want.Recent), len(got.Recent))
	}
	return fmt.Sprintf("first differing element: %v vs %v", want.Recent, got.Recent)
}

// scalarInt runs a statement that returns one integer, such as an
// a_*_drift audit. A statement that returns no row is a failure, not a zero.
func scalarInt(ctx context.Context, db ports.DB, cat *catalogue, name string, vals map[string]any) (int64, error) {
	if !cat.has(name) {
		return 0, fmt.Errorf("scenario %s has no statement %s", cat.design.ID, name)
	}
	a, err := cat.args(name, vals)
	if err != nil {
		return 0, err
	}
	var n int64
	if err := db.QueryRow(ctx, cat.stmt(name).SQL, a...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func addOnce(xs []string, v string) []string {
	for _, x := range xs {
		if x == v {
			return xs
		}
	}
	return append(xs, v)
}

func drain(rows ports.Rows) (int, error) {
	defer rows.Close()
	n := 0
	for rows.Next() {
		n++
	}
	return n, rows.Err()
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 160 {
		s = s[:160] + "…"
	}
	return strings.ReplaceAll(s, "|", "\\|")
}
