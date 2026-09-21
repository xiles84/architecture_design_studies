package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"adsplatform/ports"
)

// Correctness gates timing. Nothing in this study reports a throughput before
// these checks pass, because a design that answers the question wrongly, quickly,
// is worth nothing (methodology 5).
//
// Every expectation here is computed in Go from the generated dataset and from the
// operations the client saw acknowledged. Nothing is checked against another query
// -- that would only prove the queries agree with each other.

// GateInfo is the result of the pre-timing verification of one cell.
type GateInfo struct {
	Passed   bool     `json:"passed"`
	Checks   int      `json:"checks"`
	Failures []string `json:"failures,omitempty"`
	// ReadStatements run, so a design that failed because a statement would not
	// execute is distinguishable from one that returned wrong rows.
	Statements []string `json:"statements"`
}

// verifyReads answers every read statement and compares what came back against
// the dataset. It runs before any timing.
func verifyReads(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, sample int) GateInfo {
	g := GateInfo{Passed: true}
	fail := func(format string, a ...any) {
		g.Passed = false
		if len(g.Failures) < 20 {
			g.Failures = append(g.Failures, fmt.Sprintf(format, a...))
		}
	}

	ips := ds.Installations
	if sample > 0 && sample < len(ips) {
		// Every installation if the fleet is small, otherwise a spread: the first,
		// the last, and an even stride. Checking only the first N would miss a
		// boundary bug in the loader.
		stride := len(ips) / sample
		if stride < 1 {
			stride = 1
		}
		var picked []Installation
		for i := 0; i < len(ips) && len(picked) < sample; i += stride {
			picked = append(picked, ips[i])
		}
		ips = picked
	}

	for _, ip := range ips {
		g.Checks++
		g.Statements = addOnce(g.Statements, sR01)
		rows, err := cat.args(sR01, map[string]any{"installed_product_id": ip.ID})
		if err != nil {
			fail("bind r01: %v", err)
			continue
		}
		qr, err := db.Query(ctx, cat.stmt(sR01).SQL, rows...)
		if err != nil {
			fail("r01 on %d: %v", ip.ID, err)
			continue
		}
		got, rev, err := scanKeyValues(qr)
		if err != nil {
			fail("r01 scan on %d: %v", ip.ID, err)
			continue
		}
		// The result arrives sorted by key; the dataset is in key-generation
		// order, which is deliberately not alphabetical (the section prefix
		// cycles). Comparing the two by index would pair the wrong rows and
		// report a correct design as broken -- which is exactly what the first
		// dev check did before this line existed.
		want := sortedEntries(ds.Entries[ip.ID])
		if len(got) != len(want) {
			fail("INV-1: installation %d returned %d keys, dataset has %d", ip.ID, len(got), len(want))
			continue
		}
		for i := range want {
			if got[i].Key != want[i].Key || got[i].Value != want[i].Value {
				fail("INV-1: installation %d key %q: got %q want %q", ip.ID, want[i].Key, got[i].Value, want[i].Value)
				break
			}
		}
		if rev != 0 {
			fail("INV-5: a freshly loaded installation reports revision %d, expected 0", rev)
		}
	}

	// The four remaining reads must execute and return a sane shape. Their
	// correctness is checked by the audit after the writing phases; here the point
	// is that a design which cannot answer one of them fails the gate rather than
	// quietly reporting four reads where the catalogue promised five.
	probe := ips[0]
	for _, name := range []string{sR02, sR03, sR04, sR05} {
		vals := map[string]any{
			"installed_product_id": probe.ID,
			"key":                  ds.Entries[probe.ID][0].Key,
			"known_revision":       int64(0),
			"scope_kind":           "product_definition",
			"scope_id":             int64(probe.ProductDefID),
		}
		a, err := cat.args(name, vals)
		if err != nil {
			fail("bind %s: %v", name, err)
			continue
		}
		g.Checks++
		g.Statements = addOnce(g.Statements, name)
		qr, err := db.Query(ctx, cat.stmt(name).SQL, a...)
		if err != nil {
			fail("%s: %v", name, err)
			continue
		}
		n, err := drain(qr)
		if err != nil {
			fail("%s scan: %v", name, err)
			continue
		}
		if n == 0 {
			fail("%s returned no rows for a loaded portal", name)
		}
	}
	return g
}

// AuditInfo is the post-write reconciliation. It is run after every writing,
// contention, churn, failure-injection and mixed-workload phase, never only once
// at the end: a design that drifts and then repairs itself between phases would
// otherwise pass.
type AuditInfo struct {
	Passed bool `json:"passed"`
	Checks int  `json:"checks"`
	// KeyValueMismatches counts installations whose database state disagrees with
	// what the client was told it had published (INV-3, INV-10).
	KeyValueMismatches int `json:"key_value_mismatches"`
	RevisionMismatch   int `json:"revision_mismatch"`
	// DesignChecks are the design's own a_*_mismatches statements, each of which
	// must return zero.
	DesignChecks map[string]int `json:"design_checks,omitempty"`
	Failures     []string       `json:"failures,omitempty"`
}

// audit reconciles the database against the ledger of acknowledged operations.
// The ledger is built by the harness as it publishes, so this is a comparison
// between what the client was promised and what the database holds -- the
// invariant that makes a fast wrong design impossible to report as a fast
// correct one.
func audit(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, led *ledger, sample int) AuditInfo {
	a := AuditInfo{Passed: true, DesignChecks: map[string]int{}}
	fail := func(format string, args ...any) {
		a.Passed = false
		if len(a.Failures) < 20 {
			a.Failures = append(a.Failures, fmt.Sprintf(format, args...))
		}
	}

	expected, revisions := led.snapshot()

	ips := ds.Installations
	if sample > 0 && sample < len(ips) {
		stride := max(len(ips)/sample, 1)
		var picked []Installation
		for i := 0; i < len(ips) && len(picked) < sample; i += stride {
			picked = append(picked, ips[i])
		}
		ips = picked
	}

	for _, ip := range ips {
		want, touched := expected[ip.ID], false
		if _, ok := revisions[ip.ID]; ok {
			touched = true
		}
		if !touched {
			continue
		}
		a.Checks++
		vals, err := cat.args(sAuditKeys, map[string]any{"installed_product_id": ip.ID})
		if err != nil {
			fail("bind %s: %v", sAuditKeys, err)
			continue
		}
		qr, err := db.Query(ctx, cat.stmt(sAuditKeys).SQL, vals...)
		if err != nil {
			fail("%s on %d: %v", sAuditKeys, ip.ID, err)
			continue
		}
		// The audit statement returns (key, value) with no revision: it is the
		// design's own account of its state, not the read catalogue's answer.
		got, err := scanPairs(qr)
		if err != nil {
			fail("%s scan on %d: %v", sAuditKeys, ip.ID, err)
			continue
		}
		if diff := firstDifference(want, got); diff != "" {
			a.KeyValueMismatches++
			fail("INV-3/INV-10: installation %d %s", ip.ID, diff)
		}

		a.Checks++
		vals, err = cat.args(sAuditRev, map[string]any{"installed_product_id": ip.ID})
		if err == nil {
			var rev int64 = -1
			if qerr := db.QueryRow(ctx, cat.stmt(sAuditRev).SQL, vals...).Scan(&rev); qerr != nil {
				fail("%s on %d: %v", sAuditRev, ip.ID, qerr)
			} else if wantRev := revisions[ip.ID]; rev != wantRev {
				a.RevisionMismatch++
				fail("INV-5: installation %d is at revision %d, the ledger says %d", ip.ID, rev, wantRev)
			}
		}
	}

	// The design's own audits. Anything named *_mismatches returns one integer
	// and must return zero: a design whose derived data has drifted is caught
	// here even when the primary state still agrees.
	for _, name := range cat.order {
		if !strings.HasSuffix(name, "_mismatches") {
			continue
		}
		a.Checks++
		qr, err := db.Query(ctx, cat.stmt(name).SQL)
		if err != nil {
			fail("%s: %v", name, err)
			continue
		}
		var mismatches int
		if qr.Next() {
			if err := qr.Scan(&mismatches); err != nil {
				fail("%s scan: %v", name, err)
			}
		}
		qr.Close()
		a.DesignChecks[name] = mismatches
		if mismatches != 0 {
			fail("INV-8/INV-9: %s reports %d mismatches", name, mismatches)
		}
	}
	return a
}

// ---------------------------------------------------------------- helpers

type keyValue struct {
	Key   string
	Value string
}

// scanKeyValues reads a three-or-two-column result of (key, value[, revision])
// and returns the pairs in key order. The revision, when present, is returned
// separately so that INV-2 can be checked without depending on column order.
func scanKeyValues(rows ports.Rows) ([]keyValue, int64, error) {
	defer rows.Close()
	var out []keyValue
	var rev int64 = -1
	for rows.Next() {
		var k, v string
		var r int64
		// The read catalogue always returns the revision last. A design that
		// returned a different shape would fail here rather than silently.
		if err := rows.Scan(&k, &v, &r); err != nil {
			return nil, -1, err
		}
		out = append(out, keyValue{k, v})
		rev = r
	}
	if err := rows.Err(); err != nil {
		return nil, -1, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, rev, nil
}

// scanPairs reads a two-column (key, value) result and returns it in key order.
// The write catalogue's audit statements have this shape rather than the read
// catalogue's three columns, and conflating the two silently mis-scans every row.
func scanPairs(rows ports.Rows) ([]keyValue, error) {
	defer rows.Close()
	var out []keyValue
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out = append(out, keyValue{k, v})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}

func firstDifference(want map[string]string, got []keyValue) string {
	wi, gi := want, map[string]string{}
	for _, kv := range got {
		gi[kv.Key] = kv.Value
	}
	for k, v := range wi {
		gv, ok := gi[k]
		if !ok {
			return fmt.Sprintf("is missing key %q", k)
		}
		if gv != v {
			return fmt.Sprintf("holds %q=%q, the client was told %q", k, gv, v)
		}
	}
	for k := range gi {
		if _, ok := wi[k]; !ok {
			return fmt.Sprintf("still holds deleted key %q", k)
		}
	}
	return ""
}

// sortedEntries returns a copy of the entries ordered by key, matching the order
// every read statement returns.
func sortedEntries(es []Entry) []Entry {
	out := append([]Entry(nil), es...)
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
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
