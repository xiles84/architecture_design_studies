package selectbench

import (
	"strings"
	"testing"
)

// The domain arm must bind Study 01's person/donation/charity rows to the same
// projection the reusable catalogue and oracles use, deterministically, and must
// still exercise the optional-match control the family's gate depends on.

func TestJoinDomainBindsProjection(t *testing.T) {
	jf := JoinDomainEdge().AsJoinFixture()
	if len(jf.Customers) != 4 {
		t.Fatalf("expected 4 parent rows, got %d", len(jf.Customers))
	}
	// Person 1 belongs to charity "Alpha"; the projected parent name is the
	// person's charity name, which is the relationship the case turns on.
	if jf.Customers[0].ID != 1 || jf.Customers[0].Name != "Alpha" {
		t.Fatalf("person/charity binding = %+v", jf.Customers[0])
	}
	if jf.Customers[3].ID != 4 || jf.Customers[3].Name != "Beta" {
		t.Fatalf("person/charity binding = %+v", jf.Customers[3])
	}
	if len(jf.Facts) != 4 || jf.Facts[0].ID != 101 || *jf.Facts[0].CustomerID != 1 {
		t.Fatalf("donation binding = %+v", jf.Facts[0])
	}

	// At threshold 1 person 3 has no qualifying donation, so the optional
	// contract must keep it with NULLs, and the WHERE-placement control must
	// therefore answer differently.
	opt := JoinOptional(jf, 1)
	nulls := 0
	for _, r := range opt {
		if r.FactID == nil {
			nulls++
		}
	}
	if nulls == 0 {
		t.Fatal("domain edge fixture has no unmatched person; the x01 control cannot fire")
	}
	if answersEqual(opt, JoinBag(jf, 1)) {
		t.Fatal("x01_optional_where would not be detected on the domain edge fixture")
	}
}

func TestJoinDomainDeterministicAndDistinct(t *testing.T) {
	spec := JoinDomainGenSpec{Charities: 3, Persons: 20, Donations: 200}
	a := GenerateJoinDomain(42, spec)
	b := GenerateJoinDomain(42, spec)
	c := GenerateJoinDomain(43, spec)
	if JoinDomainFingerprint(a) != JoinDomainFingerprint(b) {
		t.Fatal("the same seed must fingerprint identically")
	}
	if JoinDomainFingerprint(a) == JoinDomainFingerprint(c) {
		t.Fatal("a different seed must fingerprint differently")
	}
	seen := map[int]bool{}
	for _, d := range a.Donations {
		if d.Rank < 1 || d.Rank > len(a.Donations) || seen[d.Rank] {
			t.Fatalf("rank is not a permutation of 1..N: %d", d.Rank)
		}
		seen[d.Rank] = true
	}
	// The domain projection must be a different dataset from the synthetic
	// fixture of the same size, or the arms would be the same evidence.
	synth := GenerateJoin(42, JoinGenSpec{Customers: 20, FactCount: 200})
	if JoinFingerprint(a.AsJoinFixture()) == JoinFingerprint(synth) {
		t.Fatal("domain and synthetic projections fingerprint identically")
	}
}

func TestJoinDomainLoadIsDeterministicAndOrdered(t *testing.T) {
	df := JoinDomainEdge()
	s1 := JoinDomainLoadSQL(df, 1000)
	s2 := JoinDomainLoadSQL(df, 1000)
	if len(s1) != len(s2) {
		t.Fatalf("load statement count differs: %d vs %d", len(s1), len(s2))
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			t.Fatalf("load script is not deterministic at statement %d", i)
		}
	}
	// Referential order: charity before person before donation.
	if !strings.HasPrefix(s1[0], "INSERT INTO "+JoinDomainSchema+".charity") {
		t.Fatalf("first load statement is not the charity insert: %s", s1[0])
	}
	if !strings.HasPrefix(s1[1], "INSERT INTO "+JoinDomainSchema+".person") {
		t.Fatalf("second load statement is not the person insert: %s", s1[1])
	}
	if !strings.HasPrefix(s1[2], "INSERT INTO "+JoinDomainSchema+".donation") {
		t.Fatalf("third load statement is not the donation insert: %s", s1[2])
	}
	if !strings.Contains(JoinDomainSchemaSQL(), "CREATE VIEW "+JoinDomainSchema+".customers") ||
		!strings.Contains(JoinDomainSchemaSQL(), "CREATE VIEW "+JoinDomainSchema+".facts") {
		t.Fatal("the domain schema must create the two catalogue-binding views")
	}
}
