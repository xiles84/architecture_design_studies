package selectbench

import (
	"strings"
	"testing"
)

// M1 regression: the YB matched arm must request a range-sharded rank index,
// because a bare `(rank)` index is hash-sharded and cannot serve `rank <= t`.
func TestJoinDomainIndexArms(t *testing.T) {
	dflt := JoinDomainSchemaSQLOptions("", "")
	if !strings.Contains(dflt, "donation (rank);") {
		t.Fatal("default arm must keep the engine's default index form")
	}
	rangeArm := JoinDomainSchemaSQLOptions("", IndexArmYBRange)
	if !strings.Contains(rangeArm, "donation (rank ASC);") {
		t.Fatal("range arm must create `donation (rank ASC)`")
	}
	if strings.Contains(rangeArm, "donation (rank);") {
		t.Fatal("range arm must not also create the default hash index")
	}
}

func TestJoinDomainDefaultIndexArm(t *testing.T) {
	if got := JoinDomainDefaultIndexArm("yugabyte", ""); got != IndexArmYBDefaultHash {
		t.Fatalf("YB default arm = %q", got)
	}
	if got := JoinDomainDefaultIndexArm("postgres", ""); got != IndexArmPGBtree {
		t.Fatalf("PG default arm = %q", got)
	}
	if got := JoinDomainDefaultIndexArm("postgres", IndexArmYBRange); got != IndexArmYBRange {
		t.Fatalf("explicit arm must be honoured, got %q", got)
	}
	if JoinDomainIndexArmApplies("postgres", IndexArmYBRange) {
		t.Fatal("the YB range arm is not applicable to PostgreSQL")
	}
	if !JoinDomainIndexArmApplies("yugabyte", IndexArmYBRange) {
		t.Fatal("the YB range arm is applicable to YugabyteDB")
	}
}

// The access-path metrics must distinguish a range index scan from a hash index
// full scan, because that distinction is the whole point of the arm.
func TestParseAccessMetrics(t *testing.T) {
	rangePlan := ` Finalize Aggregate
   ->  Index Only Scan using donation_rank_idx on donation  (cost=24.10..25.23 rows=21 width=0)
         Index Cond: (rank <= 20)
         Heap Fetches: 0
 Storage Index Rows Scanned: 20
 Storage Rows Scanned: 20
 Buffers: shared hit=4 read=2
 Execution Time: 0.383 ms`
	m := ParseAccessMetrics(rangePlan)
	if !m.PlanIndexScan || m.PlanSeqScan {
		t.Fatalf("range plan should be an index scan, got index=%v seq=%v", m.PlanIndexScan, m.PlanSeqScan)
	}
	if m.StorageRowsScanned != 20 || m.IndexRowsScanned != 20 {
		t.Fatalf("rows scanned = %d/%d", m.StorageRowsScanned, m.IndexRowsScanned)
	}
	if m.BufSharedHit != 4 || m.BufSharedRead != 2 {
		t.Fatalf("buffers = %d/%d", m.BufSharedHit, m.BufSharedRead)
	}

	hashPlan := ` Aggregate
   ->  Index Only Scan using donation_rank_idx on donation  (cost=24.10..10024.23 rows=20 width=0)
         Storage Filter: (rank <= 20)
 Rows Removed by Filter: 199980
 Execution Time: 57.0 ms`
	h := ParseAccessMetrics(hashPlan)
	if h.StorageRowsScanned != 0 {
		t.Fatalf("hash plan has no storage scan count, got %d", h.StorageRowsScanned)
	}
	if h.RowsRemovedByFilter != 199980 {
		t.Fatalf("rows removed = %d", h.RowsRemovedByFilter)
	}
}

// Review 3, P1: the "non-colocated" arm was one tablet per relation on one
// leader. The distributed arm must pre-split the relations and the range index.
func TestJoinDomainDistributedPlacement(t *testing.T) {
	ddl := JoinDomainSchemaSQLOptions("distributed", IndexArmYBRange)
	if strings.Count(ddl, "SPLIT INTO 3 TABLETS") != 4 {
		t.Fatalf("distributed arm must split three tables + the person index, got %d", strings.Count(ddl, "SPLIT INTO 3 TABLETS"))
	}
	if !strings.Contains(ddl, "donation (rank ASC) SPLIT AT VALUES ((66667), (133334))") {
		t.Fatal("distributed arm must pre-split the range rank index at the thirds")
	}
	if strings.Contains(ddl, "colocation") {
		t.Fatal("the distributed arm is non-colocated; it must not carry a colocation clause")
	}
	single := JoinDomainSchemaSQLOptions("noncolocated-single-tablet", IndexArmYBRange)
	if strings.Contains(single, "SPLIT") {
		t.Fatal("the single-tablet arm must not split")
	}
}
