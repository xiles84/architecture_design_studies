package selectbench

import "testing"

const seqScan = `Seq Scan on orders o  (cost=0.00..1.24 rows=2 width=8) (actual time=0.010..0.012 rows=2 loops=1)
  Filter: (amount > 100)
  Buffers: shared hit=1
Planning Time: 0.100 ms
Execution Time: 0.050 ms`

// The same plan shape and filter, with different statistics and timings.
const seqScanResampled = `Seq Scan on orders o  (cost=0.00..1.40 rows=3 width=8) (actual time=0.020..0.030 rows=3 loops=1)
  Filter: (amount > 100)
  Buffers: shared hit=2
Planning Time: 0.200 ms
Execution Time: 0.090 ms`

const seqScanOtherPredicate = `Seq Scan on orders o  (cost=0.00..1.24 rows=1 width=8) (actual time=0.008..0.009 rows=1 loops=1)
  Filter: (amount > 200)
  Buffers: shared hit=1
Planning Time: 0.110 ms
Execution Time: 0.040 ms`

const indexScan = `Index Scan using orders_amount_idx on orders o  (cost=0.14..8.16 rows=2 width=8) (actual time=0.015..0.020 rows=2 loops=1)
  Index Cond: (amount > 100)
  Buffers: shared hit=3
Planning Time: 0.120 ms
Execution Time: 0.060 ms`

func TestNormalizeStripsVolatileMetadata(t *testing.T) {
	lines := Normalize(seqScan)
	want := []string{"Seq Scan on orders o", "Filter: (amount > 100)"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines %q, want %d %q", len(lines), lines, len(want), want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

// A re-run of the same plan must not be reported as a plan change.
func TestSignatureIgnoresCostsRowsAndTimings(t *testing.T) {
	if Signature(seqScan) != Signature(seqScanResampled) {
		t.Errorf("same plan with different statistics produced different signatures: %s vs %s",
			Signature(seqScan), Signature(seqScanResampled))
	}
	if !SameStructure(seqScan, seqScanResampled) {
		t.Error("same plan with different statistics was not the same structure")
	}
}

// A different access path is a real plan difference and must be visible.
func TestSignatureChangesWithAccessPath(t *testing.T) {
	if Signature(seqScan) == Signature(indexScan) {
		t.Error("Seq Scan and Index Scan produced the same signature")
	}
	if SameStructure(seqScan, indexScan) {
		t.Error("Seq Scan and Index Scan were reported as the same structure")
	}
}

// The same shape with a different predicate is the same structure but a
// different plan: the filter is part of the meaning, not the shape.
func TestStructureIgnoresPredicateButSignatureKeepsIt(t *testing.T) {
	if !SameStructure(seqScan, seqScanOtherPredicate) {
		t.Error("same shape with a different filter should share a structure")
	}
	if Signature(seqScan) == Signature(seqScanOtherPredicate) {
		t.Error("a different filter must change the full signature")
	}
}

// Print order must not change the shape signature.
func TestStructureIsOrderIndependent(t *testing.T) {
	a := "Hash Join  (cost=1.0..2.0 rows=2 width=8)\n  -> Seq Scan on orders o\n  -> Seq Scan on customers c"
	b := "-> Seq Scan on customers c\n-> Seq Scan on orders o\nHash Join  (cost=9.9..9.9 rows=99 width=8)"
	if !SameStructure(a, b) {
		t.Errorf("the same node set in a different print order was not the same structure:\n%q\n%q",
			Structure(a), Structure(b))
	}
}
