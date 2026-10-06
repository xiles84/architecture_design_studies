package evidence

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	planSpecRel     = "book/appendix/query-plans.spec.json"
	planDataRel     = "book/appendix/query-plans.json"
	planSQLRel      = "book/appendix/query-plans.sql.json"
	planRegistryRel = "book/evidence/v6/claims.json"
)

func loadPlanFixture(t *testing.T) (string, *PlanSpec, *PlanSQL, *V2Document) {
	t.Helper()
	root := repoRoot(t)
	spec, err := LoadPlanSpec(filepath.Join(root, filepath.FromSlash(planSpecRel)))
	if err != nil {
		t.Fatal(err)
	}
	q, err := LoadPlanSQL(filepath.Join(root, filepath.FromSlash(planSQLRel)))
	if err != nil {
		t.Fatal(err)
	}
	reg, _, err := LoadV2(filepath.Join(root, filepath.FromSlash(planRegistryRel)))
	if err != nil {
		t.Fatal(err)
	}
	return root, spec, q, reg
}

// TestTrimPlanKeepsTreeAndDecisiveAttributes pins the trim: node lines lose only
// their cost estimate, listed attributes survive with their indentation, and the
// bulk (Output, Buffers, RocksDB metrics, estimates) is dropped.
func TestTrimPlanKeepsTreeAndDecisiveAttributes(t *testing.T) {
	plan := strings.Join([]string{
		"Hash Join  (cost=2.82..54.26 rows=21 width=16) (actual time=0.017..0.043 rows=20 loops=1)",
		"  Output: p.person_id, ch.name",
		"  Hash Cond: (p.charity_id = ch.charity_id)",
		"  Buffers: shared hit=84",
		"  ->  Seq Scan on donation d  (cost=21.10..18842.09 rows=200000 width=16) (actual time=0.590..88.429 rows=200000 loops=1)",
		"        Storage Table Read Requests: 196",
		"          Metric rocksdb_number_db_seek: 196.000",
		"        Estimated Seeks: 196",
		"        Rows Removed by Filter: 7",
		"Planning:",
		"  Buffers: shared hit=20",
		"Planning Time: 0.133 ms",
		"Execution Time: 0.054 ms",
		"",
	}, "\n")
	got, total, kept := TrimPlan(plan)
	want := strings.Join([]string{
		"Hash Join (actual time=0.017..0.043 rows=20 loops=1)",
		"  Hash Cond: (p.charity_id = ch.charity_id)",
		"  ->  Seq Scan on donation d (actual time=0.590..88.429 rows=200000 loops=1)",
		"        Storage Table Read Requests: 196",
		"        Rows Removed by Filter: 7",
		"Planning Time: 0.133 ms",
		"Execution Time: 0.054 ms",
	}, "\n")
	if got != want {
		t.Fatalf("trim mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
	if total != 13 || kept != 7 {
		t.Fatalf("line counts = %d/%d, want 13/7", kept, total)
	}
}

// TestPlanShapeIgnoresNumbersNotNodes: two executions of one plan share a shape;
// changing an access path changes it. This is the meaning of "same node tree".
func TestPlanShapeIgnoresNumbersNotNodes(t *testing.T) {
	a := "Hash Join  (cost=1..2 rows=3 width=4) (actual time=0.1..0.2 rows=20 loops=1)\n  ->  Seq Scan on person p  (cost=0..1 rows=1 width=8) (actual time=0.3..0.4 rows=2000 loops=1)\n"
	b := "Hash Join  (cost=9..9 rows=9 width=4) (actual time=5.1..7.2 rows=99 loops=1)\n  ->  Seq Scan on person p  (cost=0..9 rows=9 width=8) (actual time=1.3..2.4 rows=7 loops=1)\n"
	c := "Hash Join  (cost=1..2 rows=3 width=4) (actual time=0.1..0.2 rows=20 loops=1)\n  ->  Index Scan using person_pkey on person p  (cost=0..1 rows=1 width=8) (actual time=0.3..0.4 rows=2000 loops=1)\n"
	if PlanShape(a) != PlanShape(b) {
		t.Fatalf("numbers changed the shape:\n%s\n%s", PlanShape(a), PlanShape(b))
	}
	if PlanShape(a) == PlanShape(c) {
		t.Fatal("a different access path kept the same shape")
	}
}

// TestPlanAppendixDataIsCurrent is the staleness gate: the committed data must be
// exactly what the generator builds from the committed spec, registry and results.
func TestPlanAppendixDataIsCurrent(t *testing.T) {
	root, spec, q, reg := loadPlanFixture(t)
	data, errs := BuildPlanData(root, planSpecRel, planSQLRel, planRegistryRel, spec, q, reg)
	for _, e := range errs {
		t.Errorf("unexpected error: %v", e)
	}
	if len(errs) > 0 {
		return
	}
	want, err := EncodePlanData(data)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(planDataRel)))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n")), want) {
		t.Fatalf("%s is stale; run `go run . plans --repo ../.. --write` in tools/evidence", planDataRel)
	}
}

// TestPlanAppendixCoversEveryMeasuredSelectClaim: every active, measured v6-sel
// claim has an exhibit, except the ones the appendix explains it cannot show.
func TestPlanAppendixCoversEveryMeasuredSelectClaim(t *testing.T) {
	_, spec, _, reg := loadPlanFixture(t)
	// v6-sel-07's evidence is a negative control that records rows, not a plan.
	noPlan := map[string]bool{"v6-sel-07-nullable-not-in-empty": true}
	shown := map[string]bool{}
	for _, e := range spec.Exhibits {
		shown[e.Claim] = true
	}
	for _, c := range reg.Claims {
		if !strings.HasPrefix(c.ClaimID, "v6-sel-") || c.Strength == "gap" || noPlan[c.ClaimID] {
			continue
		}
		if !shown[c.ClaimID] {
			t.Errorf("measured SELECT claim %s has no plan exhibit", c.ClaimID)
		}
	}
}

// TestPlanExhibitsRejectUnanchoredRunsAndWrongShapes is the negative control: the
// generator must refuse a run the claim does not cite, a shape statement the plans
// contradict, an axis that does not exist and a gap claim.
func TestPlanExhibitsRejectUnanchoredRunsAndWrongShapes(t *testing.T) {
	root, _, q, reg := loadPlanFixture(t)
	cases := []struct {
		name string
		ex   ExhibitSpec
		want string
	}{
		{"unanchored run", ExhibitSpec{ID: "x", Claim: "v6-sel-09-aggregate-correlated-postgresql-crossover",
			RunID: "select-order-aggregate-yugabyte-1-20261004", File: "aggregate.json",
			Plans: []PlanRef{{Axis: "family=aggregate;stmt=correlated;parent_max=1"}}}, "is not an anchor"},
		{"contradicted shape", ExhibitSpec{ID: "x", Claim: "v6-sel-01-inner-join-spellings-same-plan",
			RunID: "select-join-v3b-domain-pg-i1-20261005", File: "join.json",
			Plans: []PlanRef{{Axis: "family=bag;stmt=join_where;threshold=2000"}},
			Shapes: []ShapeCheck{{Axis: "family=bag;stmt=filtered_derived;threshold=2000",
				As: "family=bag;stmt=join_where;threshold=2000", Expect: "different"}}}, "expects a different one"},
		{"missing axis", ExhibitSpec{ID: "x", Claim: "v6-sel-01-inner-join-spellings-same-plan",
			RunID: "select-join-v3b-domain-pg-i1-20261005", File: "join.json",
			Plans: []PlanRef{{Axis: "family=bag;stmt=join_where;threshold=7"}}}, "no cell with axis"},
		{"gap claim", ExhibitSpec{ID: "x", Claim: "v6-sel-gap-01-yugabytedb-scan-only-and-single-node-families",
			RunID: "select-order-access-yugabyte-1-20261004", File: "access.json",
			Plans: []PlanRef{{Axis: "shape=selective/indexable;spelling=or"}}}, "is a gap"},
	}
	for _, tc := range cases {
		_, errs := BuildPlanData(root, planSpecRel, planSQLRel, planRegistryRel, &PlanSpec{Schema: PlanSpecSchema, Exhibits: []ExhibitSpec{tc.ex}}, q, reg)
		found := false
		for _, e := range errs {
			if strings.Contains(e.Error(), tc.want) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: want an error containing %q, got %v", tc.name, tc.want, errs)
		}
	}
}

// TestPlanSQLMustMatchTheCellAndItsCode: the bound SQL is accepted only for the
// exact statement and parameters the cell ran, and only while every catalogue
// source is unchanged since the producing commit. A commit that predates the join
// catalogue must be refused, and a full plan set with no matching SQL must fail.
func TestPlanSQLMustMatchTheCellAndItsCode(t *testing.T) {
	root, spec, q, reg := loadPlanFixture(t)

	stale := *q
	stale.Statements = append([]PlanSQLEntry(nil), q.Statements...)
	for i := range stale.Statements {
		if stale.Statements[i].Params != "" {
			stale.Statements[i].Params += ",extra=1"
			break
		}
	}
	_, errs := BuildPlanData(root, planSpecRel, planSQLRel, planRegistryRel, spec, &stale, reg)
	if !errorsContain(errs, "but the cell ran") {
		t.Errorf("a parameter mismatch was accepted: %v", errs)
	}

	empty := &PlanSQL{Schema: PlanSQLSchema, Catalogues: q.Catalogues}
	_, errs = BuildPlanData(root, planSpecRel, planSQLRel, planRegistryRel, spec, empty, reg)
	if !errorsContain(errs, "no SQL for") {
		t.Errorf("missing SQL was accepted: %v", errs)
	}

	c := &sourceChecker{root: root, cache: map[string]error{}}
	// ca2a38fb produced the access/aggregate/existence runs, before the join catalogue existed.
	if err := c.check("ca2a38fbf9e53da00344578258138db9e01d3545", "studies/01-charity-tree/selectbench/joincatalog.go"); err == nil {
		t.Error("a source file absent at the producing commit was accepted")
	}
	if err := c.check("ca2a38fbf9e53da00344578258138db9e01d3545", "studies/01-charity-tree/selectbench/bind.go"); err != nil {
		t.Errorf("an unchanged source was refused: %v", err)
	}
}

func errorsContain(errs []error, want string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error(), want) {
			return true
		}
	}
	return false
}
