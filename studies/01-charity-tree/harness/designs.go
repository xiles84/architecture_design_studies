package main

import "fmt"

// Design describes one schema variant under test. The flags tell the loader and
// the write driver how to shape data for this design; everything a reader does
// is driven by the design's own queries.sql instead.
type Design struct {
	ID      string
	Title   string
	Summary string
	Engines []string // engines this design can run on

	Triggers          bool // has triggers.sql, applied after the bulk load
	CharityOnDonation bool // donation table carries a denormalised charity_id
	Rollups           bool // person/charity carry consolidated aggregate columns
	Embedded          bool // donations live inside the person row as JSONB
	AppRollup         bool // the harness, not a trigger, maintains the aggregates
	RecentCache       bool // person carries a bounded newest-first cache of children
	SumOnly           bool // only total_donated_cents is materialised; other reads use D3 SQL
}

// recentCacheSize is the bound on D9's embedded slice. It is a design constant,
// not a tuning knob: the catalogue's q09 asks for exactly this many rows, and a
// query asking for one more would fall off the cache entirely.
const recentCacheSize = 20

var designs = []Design{
	{ID: "d11_copied_key", Title: "copied-key-only", Summary: "D2 reads/indexes with D3's copied key and its FK; no charity access path yet.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d12_recency_index", Title: "recency-index-only", Summary: "D11 plus charity/date index, with D2 SQL unchanged.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d13_recency_sql", Title: "recency-rewritten", Summary: "D12 with only q02/q05/q12 recency reads rewritten to use the copied key.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d17_sum_sql", Title: "sum-rewritten", Summary: "D13 with only q08 rewritten; no additional sum index.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d14_sum_plain", Title: "sum-plain-index", Summary: "D17 plus a plain charity index; separates index width from covering payload.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d15_sum_covering", Title: "sum-covering-index", Summary: "D14 index includes amount_cents; D3 differs only in q03/q04 ranking rewrites.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true},
	{ID: "d16_sum_rollup", Title: "sum-only-trigger", Summary: "D3 plus sums on both parents; keeps parent locks but omits count/extrema maintenance.", Engines: []string{"postgres", "yugabyte"}, CharityOnDonation: true, Triggers: true, SumOnly: true},
	{
		ID: "d1_normalized_minimal", Title: "normalized-minimal",
		Summary: "Textbook 3NF. Primary and foreign keys only, no secondary indexes.",
		Engines: []string{"postgres", "yugabyte"},
	},
	{
		ID: "d2_normalized_indexed", Title: "normalized-indexed",
		Summary: "D1 plus secondary indexes tuned to the query catalogue. Byte-identical SQL to D1.",
		Engines: []string{"postgres", "yugabyte"},
	},
	{
		ID: "d3_flattened_fk", Title: "flattened-fk",
		Summary:           "D2 plus the grandparent key (charity_id) denormalised onto donation.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true,
	},
	{
		ID: "d4_rollup_trigger", Title: "rollup-trigger",
		Summary:           "D3 plus consolidated aggregates on person and charity, maintained by triggers.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true, Rollups: true, Triggers: true,
	},
	{
		ID: "d5_rollup_app", Title: "rollup-app",
		Summary:           "Same columns as D4, but the application maintains the aggregates.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true, Rollups: true, AppRollup: true,
	},
	{
		ID: "d6_embedded_jsonb", Title: "embedded-jsonb",
		Summary:  "The donation child table folded into the person row as a JSONB array.",
		Engines:  []string{"postgres", "yugabyte"},
		Embedded: true,
	},
	{
		ID: "d9_embedded_hybrid", Title: "embedded-hybrid",
		Summary:           "D3 plus a bounded 20-element newest-first cache of donations on the person row.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true, Triggers: true, RecentCache: true,
	},
	{
		ID: "d10_embedded_hybrid_locked", Title: "embedded-hybrid-locked",
		Summary:           "D9 with a concurrency-correct cache trigger (merge-sort on insert, lock before rebuild); reads byte-identical.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true, Triggers: true, RecentCache: true,
	},
	{
		ID: "d8_flattened_nofk", Title: "flattened-nofk",
		Summary:           "D3 with the FOREIGN KEY constraints removed; everything else byte-identical.",
		Engines:           []string{"postgres", "yugabyte"},
		CharityOnDonation: true,
	},
	{
		ID: "d7_yb_child_colocated", Title: "yb-child-colocated",
		Summary:           "D3's columns, with donation sharded by person_id so one person's donations share a tablet.",
		Engines:           []string{"yugabyte"},
		CharityOnDonation: true,
	},
}

func designByID(id string) (Design, error) {
	for _, d := range designs {
		if d.ID == id {
			return d, nil
		}
	}
	return Design{}, fmt.Errorf("unknown design %q", id)
}

func (d Design) supports(engine string) bool {
	for _, e := range d.Engines {
		if e == engine {
			return true
		}
	}
	return false
}
