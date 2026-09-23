package main

// This file is the study's design registry: for each design, what mechanism it
// uses, which decision it isolates, and which of its SQL statements the harness
// must be able to find. The SQL itself is in ../sql/<id>/.

// Kind is the storage mechanism. It decides which statements the harness drives;
// the SQL files decide what those statements do.
type Kind string

const (
	Rows Kind = "rows" // one row per configuration key
	Doc  Kind = "doc"  // one JSON document per installed product
)

// RollupMode is who maintains the parent aggregates that r04 reads.
type RollupMode string

const (
	RollupNone    RollupMode = "none"    // the read aggregates every time
	RollupTrigger RollupMode = "trigger" // n3: the database maintains them
	RollupApp     RollupMode = "app"     // n4: the application maintains them, in the transaction
	RollupDrift   RollupMode = "drift"   // x2: the application maintains them, after commit -- the control
)

// ConcMode is how concurrent publications of the same key arbitrate.
type ConcMode string

const (
	// ConcNone writes an explicit value. Last writer wins, which is correct:
	// there is nothing to lose.
	ConcNone ConcMode = "none"
	// ConcOptimistic reads a version and writes back only if it has not moved.
	ConcOptimistic ConcMode = "optimistic"
	// ConcPessimistic locks the parent row before reading anything.
	ConcPessimistic ConcMode = "pessimistic"
	// ConcUnsafe is the negative control: read-modify-write with no guard.
	ConcUnsafe ConcMode = "unsafe"
)

// Design describes one way of keeping a portal's configuration.
type Design struct {
	ID      string
	Short   string
	Title   string
	Family  string
	Summary string

	Kind     Kind
	Rollup   RollupMode
	Rolldown bool
	Conc     ConcMode

	// ReadModifyWrite means this design's modify op derives the new value from
	// the value it just read. A design that writes an explicit value cannot lose
	// an update -- only be overwritten -- so this flag is what makes INV-3
	// testable at all.
	ReadModifyWrite bool

	// NegativeControl designs are measured like every other design but their
	// speed is never reported as a plain number: a design that loses updates
	// quickly has not been fast (methodology 5a).
	NegativeControl bool

	// YBOnly marks SQL that uses YugabyteDB-only syntax.
	YBOnly bool

	// Pair names the design this one differs from by exactly one decision.
	Pair string

	// Needs lists statements that must exist in the design's catalogue. A
	// missing statement is a harness bug and must fail at load, not silently
	// change what the design does.
	Needs []string

	// Risk names the correctness risk this design carries, for the report.
	Risk string
}

// requirable statement names, so a typo in a design's Needs is a compile error
// rather than a runtime surprise.
const (
	sR01 = "r01_effective_config"
	sR02 = "r02_read_key"
	sR03 = "r03_revision_check"
	sR04 = "r04_list_installations"
	sR05 = "r05_search_key"

	sRevBump  = "w_revision_bump"
	sW01      = "w01_modify_key"
	sW02      = "w02_add_key"
	sW03      = "w03_delete_key"
	sW04Up    = "w04_batch_upsert"
	sW04Del   = "w04_batch_delete"
	sW05Clear = "w05_replace_clear"
	sW05Fill  = "w05_replace_fill"
	sW06      = "w06_update_metadata"

	sAuditKeys  = "a_key_values"
	sAuditRev   = "a_revision"
	sAuditCount = "a_counts"
)

var baseReads = []string{sR01, sR02, sR03, sR04, sR05}

// rowDesign is the statement set shared by every row-per-key design. Keeping it
// in one place is what makes "these designs differ by one decision" checkable
// with diff -r rather than by memory.
func rowDesign(needs ...string) []string {
	out := append([]string{}, baseReads...)
	out = append(out, sRevBump, sW01, sW02, sW03, sW04Up, sW04Del, sW05Clear, sW05Fill, sW06,
		sAuditKeys, sAuditRev, sAuditCount, "w_load_entries")
	return append(out, needs...)
}

var designs = []Design{
	{
		ID: "n0_rows_unindexed", Short: "n0", Title: "row per key, no search index",
		Family: "normalized", Kind: Rows, Rollup: RollupNone, Pair: "n1_rows_indexed",
		Summary: "n1 with the secondary index on (key) removed. The primary key still serves a complete read, so only the cross-installation search is affected.",
		Risk:    "none specific", Needs: rowDesign(),
	},
	{
		ID: "n1_rows_indexed", Short: "n1", Title: "row per key, indexed (reference)",
		Family: "normalized", Kind: Rows, Rollup: RollupNone,
		Summary: "The floor: one row per key, a primary key for the complete read and one secondary index for the search. Every other design is read against this one.",
		Risk:    "none specific", Needs: rowDesign(),
	},
	{
		ID: "n2_rows_rolldown", Short: "n2", Title: "row per key with copied parent keys",
		Family: "normalized", Kind: Rows, Rollup: RollupNone, Rolldown: true, Pair: "n1_rows_indexed",
		Summary: "Rolldown: the definition, environment and business unit are copied onto every configuration row so the search never joins the parent. Costs three indexes and a wider insert.",
		Risk:    "the copied keys can drift from their parents (INV-8)", Needs: rowDesign("a_rolldown_mismatches"),
	},
	{
		ID: "n3_rollup_trigger", Short: "n3", Title: "parent rollups maintained by trigger",
		Family: "rollup", Kind: Rows, Rollup: RollupTrigger, Pair: "n4_rollup_app",
		Summary: "count, revision, last modification and content hash live on the installed product and a row trigger maintains them. The overview read stops aggregating; every write starts paying.",
		Risk:    "the rollup can drift from the entries it summarises (INV-9)", Needs: rowDesign("a_rollup_mismatches"),
	},
	{
		ID: "n4_rollup_app", Short: "n4", Title: "parent rollups maintained by the application",
		Family: "rollup", Kind: Rows, Rollup: RollupApp, Pair: "n3_rollup_trigger",
		Summary: "The same four columns, maintained by the application inside the publishing transaction. Identical schema to n3; only who does the work changes.",
		Risk:    "the rollup can drift if a write path forgets it (INV-9)", Needs: rowDesign("w_recompute_rollup", "a_rollup_mismatches"),
	},
	{
		ID: "d1_doc_row", Short: "d1", Title: "one document row per installation",
		Family: "document", Kind: Doc, Rollup: RollupNone, Pair: "n1_rows_indexed",
		Summary: "A separate one-to-one JSON document row, replaced whole on every publication. One row answers a complete read; every publication rewrites every key.",
		Risk:    "a racing publication could overwrite another (guarded by revision, INV-3)",
		Needs: []string{sR01, sR02, sR03, sR04, sR05, sRevBump, sW06, "wd1_read_doc", "wd1_write_doc",
			sAuditKeys, sAuditRev, sAuditCount, "a_document_revision_mismatches", "a_content_hash_mismatches", "w_load_documents"},
	},
	{
		ID: "d2_doc_on_parent", Short: "d2", Title: "document embedded in the parent row",
		Family: "document", Kind: Doc, Rollup: RollupNone, Pair: "d1_doc_row",
		Summary: "The same statements as d1, but the document lives as a column on installed_product instead of a separate row. The publication is the same size; the metadata update (W6) now rewrites the document too.",
		Risk:    "same guard as d1, plus W6 amplification: an ordinary metadata write rewrites TOASTed configuration",
		Needs: []string{sR01, sR02, sR03, sR04, sR05, sRevBump, sW06, "wd1_read_doc", "wd1_write_doc",
			sAuditKeys, sAuditRev, sAuditCount, "a_document_revision_mismatches", "a_content_hash_mismatches", "w_load_documents"},
	},
	{
		ID: "c1_optimistic_version", Short: "c1", Title: "optimistic version check with bounded retries",
		Family: "concurrency", Kind: Rows, Rollup: RollupNone, Conc: ConcOptimistic,
		ReadModifyWrite: true, Pair: "c2_pessimistic_lock",
		Summary: "Read the row's version, write back only if it has not moved, and retry within a deadline when it has. Losing the race is a signal, not a failure.",
		Risk:    "unbounded retries would turn contention into a timeout (bounded by the deadline, INV-13)",
		Needs:   rowDesign("wc1_read_version", "wc1_update_guarded"),
	},
	{
		ID: "c2_pessimistic_lock", Short: "c2", Title: "lock the parent before updating",
		Family: "concurrency", Kind: Rows, Rollup: RollupNone, Conc: ConcPessimistic,
		ReadModifyWrite: true, Pair: "c1_optimistic_version",
		Summary: "Take the installed product's row lock, then read and write. The competitor waits instead of being rejected. Same schema as c1, guard column included.",
		Risk:    "waiters hold connections; a long transaction serialises the whole installation",
		Needs:   rowDesign("wc2_lock_installed_product", "wc2_read_value", "wc2_write_value"),
	},
	{
		ID: "x1_lost_update_control", Short: "x1", Title: "unchecked read-modify-write (negative control)",
		Family: "control", Kind: Rows, Rollup: RollupNone, Conc: ConcUnsafe,
		ReadModifyWrite: true, NegativeControl: true, Pair: "c1_optimistic_version",
		Summary: "Read the value, write back a value derived from it, with no version guard and no lock. Two concurrent updates can both be acknowledged while one is lost.",
		Risk:    "MUST violate INV-3 under the hot-key scenario; if it does not fire, no design's correctness is demonstrated",
		Needs:   rowDesign("wx1_read_value", "wx1_write_value", "a_version_sum"),
	},
	{
		ID: "x2_rollup_drift_control", Short: "x2", Title: "rollups maintained after commit (negative control)",
		Family: "control", Kind: Rows, Rollup: RollupDrift, Conc: ConcNone,
		NegativeControl: true, Pair: "n4_rollup_app",
		Summary: "Byte-identical SQL to n4. The only difference is that the rollup recomputation runs in a second transaction after the publication committed, leaving a window in which the stored rollup is wrong.",
		Risk:    "MUST violate INV-9/INV-10 under concurrency; this is the audit that proves the rollup audits work",
		Needs:   rowDesign("w_recompute_rollup", "a_rollup_mismatches"),
	},
}

// hasContention reports whether this design has a race worth running in the
// contention phase. An empty Conc is the zero value, not a strategy: a design
// that writes an explicit value needs no arbitration, and racing it would try to
// drive statements it does not have -- which is how the document design crashed
// the phase before this test existed.
func (d Design) hasContention() bool {
	return d.Conc == ConcOptimistic || d.Conc == ConcPessimistic ||
		d.Conc == ConcUnsafe || d.Rollup == RollupDrift
}

var designIndex = func() map[string]Design {
	m := make(map[string]Design, len(designs))
	for _, d := range designs {
		m[d.ID] = d
	}
	return m
}()

func designByID(id string) (Design, bool) {
	d, ok := designIndex[id]
	return d, ok
}
