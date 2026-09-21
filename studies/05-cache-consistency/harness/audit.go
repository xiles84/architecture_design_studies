package main

import (
	"context"
	"fmt"
	"sort"

	"adsplatform/ports"
)

// The audit: the ledger replay the protocol requires. After every writing phase the
// database is compared against the oracle, whose state was built from the generated
// dataset and the operations the client saw acknowledged -- never from the database.
//
// Three comparisons per sampled donor, because they fail differently:
//
//   - the DONATION ID SET catches a write that was acknowledged but vanished, or one
//     that was rejected and landed anyway (the source-write invariant);
//   - the AGGREGATES catch a stored rollup that drifted;
//   - the PORTAL CONTENT HASH catches everything else, including an embedded slice
//     that lost or reordered an element.
func auditCell(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, orc *oracle, sample int) AuditInfo {
	a := AuditInfo{Passed: true, DesignChecks: map[string]int{}}
	fail := func(format string, args ...any) {
		a.Passed = false
		if len(a.Failures) < 20 {
			a.Failures = append(a.Failures, fmt.Sprintf(format, args...))
		}
	}
	usesVersion := cat.design.usesVersion()

	for _, p := range samplePeople(ds, sample) {
		a.Checks++

		// 1. The exact donation id set.
		gotIDs, err := queryInt64Column(ctx, db, cat, sAuditDonationIDs, map[string]any{"person_id": p.ID})
		if err != nil {
			fail("%s on donor %d: %v", sAuditDonationIDs, p.ID, err)
			continue
		}
		wantIDs := orc.ExpectedDonationIDs(p.ID)
		if !equalInt64(gotIDs, wantIDs) {
			a.DonationSetMismatch++
			fail("INV-5: donor %d holds %d donations, the ledger promised %d (%s)",
				p.ID, len(gotIDs), len(wantIDs), setDiff(wantIDs, gotIDs))
		}

		// 2. The aggregates, recomputed from the child table.
		var count, total int64
		args, err := cat.args(sAuditAggregate, map[string]any{"person_id": p.ID})
		if err != nil {
			fail("bind %s: %v", sAuditAggregate, err)
			continue
		}
		if err := db.QueryRow(ctx, cat.stmt(sAuditAggregate).SQL, args...).Scan(&count, &total); err != nil {
			fail("%s on donor %d: %v", sAuditAggregate, p.ID, err)
			continue
		}
		want, ok := orc.ExpectedContent(p.ID)
		if !ok {
			fail("oracle has no donor %d", p.ID)
			continue
		}
		if count != want.DonationCount || total != want.DonationTotalCents {
			a.AggregateMismatch++
			fail("INV-6: donor %d aggregates are (%d, %d), the ledger promised (%d, %d)",
				p.ID, count, total, want.DonationCount, want.DonationTotalCents)
		}

		// 3. The portal content, through the study's one read implementation.
		got, _, err := readPortalSQL(ctx, db, cat, usesVersion, p.ID)
		if err != nil {
			fail("portal read on donor %d: %v", p.ID, err)
			continue
		}
		if got.ContentHash() != want.ContentHash() {
			a.KeyValueMismatches++
			fail("INV-7: donor %d portal content differs from the ledger-derived expectation (%s)",
				p.ID, firstDifferencePortal(want, got))
		}
		if ids := recentIDs(got); !equalInt64(ids, orc.ExpectedRecentIDs(p.ID)) {
			a.RecentSliceMismatch++
			fail("INV-8: donor %d recent slice is not the deterministic newest-20 slice", p.ID)
		}

		// 4. The owned model's version/outbox invariant.
		if usesVersion {
			a.Checks++
			n, err := scalarInt(ctx, db, cat, sAuditVersionMonotonic, map[string]any{"person_id": p.ID})
			if err != nil {
				fail("%s on donor %d: %v", sAuditVersionMonotonic, p.ID, err)
			} else if n != 0 {
				fail("INV-9: donor %d cache_version disagrees with its outbox event count", p.ID)
			}
		}
	}

	// 5. The design's own drift audits, where the design has them. Each must return
	//    zero, and a missing statement is a harness bug rather than a pass.
	for _, name := range []string{sAuditRolldownDrift, sAuditRollupDrift, sAuditEmbeddedDrift, sAuditOutboxOrphans} {
		if !cat.has(name) {
			continue
		}
		a.Checks++
		n, err := scalarInt(ctx, db, cat, name, nil)
		if err != nil {
			fail("%s: %v", name, err)
			continue
		}
		a.DesignChecks[name] = int(n)
		if n != 0 {
			fail("INV-10: %s reports %d mismatches", name, n)
		}
	}
	return a
}

func queryInt64Column(ctx context.Context, db ports.DB, cat *catalogue, name string, vals map[string]any) ([]int64, error) {
	args, err := cat.args(name, vals)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(ctx, cat.stmt(name).SQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, nil
}

func setDiff(want, got []int64) string {
	inGot := map[int64]bool{}
	for _, v := range got {
		inGot[v] = true
	}
	inWant := map[int64]bool{}
	for _, v := range want {
		inWant[v] = true
	}
	var missing, extra []int64
	for _, v := range want {
		if !inGot[v] {
			missing = append(missing, v)
		}
	}
	for _, v := range got {
		if !inWant[v] {
			extra = append(extra, v)
		}
	}
	return fmt.Sprintf("missing %v, unexpected %v", firstN(missing, 5), firstN(extra, 5))
}

func firstN(xs []int64, n int) []int64 {
	if len(xs) > n {
		return xs[:n]
	}
	return xs
}

// checkLedgerMatchesDB asserts the one property every strict conclusion rests on: the
// ledger's freshness requirement for a key is never AHEAD of the content the database
// actually holds.
//
// It exists because the opposite was observed. An authoritative BYPASS read -- one
// that consults no cache and runs in a single repeatable-read transaction -- returned
// a state one and two versions behind the requirement, which no cache mechanism can
// cause. That means the requirement was wrong, and a cell cannot be judged until this
// passes (AM-01, finding 3a).
//
// A requirement BEHIND a committed-but-unacknowledged state is legitimate and is not
// counted: that state classifies as ahead.
func (c *cell) checkLedgerMatchesDB(ctx context.Context, phase string) error {
	usesVersion := c.d.usesVersion()
	lc := LedgerCheck{Phase: phase, People: len(c.ds.People), Passed: true}
	for _, p := range c.ds.People {
		content, _, err := readPortalSQL(ctx, c.db, c.cat, usesVersion, p.ID)
		if err != nil {
			return fmt.Errorf("ledger assertion after %s: %w", phase, err)
		}
		dbHash := content.ContentHash()
		reqHash, reqSeq := c.orc.Required(p.ID)
		if dbHash == reqHash {
			continue
		}
		if k, _, _, _ := c.orc.Classify(p.ID, dbHash, reqHash, reqSeq); k == KindStale {
			lc.Mismatches++
			if len(lc.Examples) < 5 {
				lc.Examples = append(lc.Examples,
					fmt.Sprintf("donor %d: the database holds a state the requirement is ahead of (requirement seq %d)", p.ID, reqSeq))
			}
		}
	}
	lc.Passed = lc.Mismatches == 0
	c.res.Ledger = append(c.res.Ledger, lc)
	if !lc.Passed {
		return fmt.Errorf("LEDGER ASSERTION FAILED after %s: the freshness requirement is ahead of the database for %d donor(s); every strict number from this cell is provisional (%s)",
			phase, lc.Mismatches, lc.Examples[0])
	}
	return nil
}
