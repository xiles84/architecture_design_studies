# Result — Study 04 v2 remaining designs

**Task:** `task-20260922T195241Z-study04-v2-remaining-designs`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-remaining-designs`
**Measured:** `d1_doc_row` vs `d2_doc_on_parent`, PostgreSQL 17.11 single node, scale small, 3 trials,
commit `c1288c71…`, tag `run/04-configuration-portal/20260923T0955Z-d2-vs-d1`.

## Delivered

- **`d2_doc_on_parent` implemented** (`sql/d2_doc_on_parent/`, `harness/designs.go`, `ALL_DESIGNS`):
  the document, guard revision, content hash and byte count as columns on `installed_product`.
  A tiny dev check and the reported run both passed the correctness gate.
- **Comparison measured:** publication 865.9 vs 858.8 ops/s; metadata update **1651.0 vs 1672.7
  ops/s** (the predicted W6 amplification did **not** appear at small cardinality); r05 search
  22898 vs 14433 ops/s (single-run, unattributed). `failed_cells: []` in both cells.
- **Signed analysis:** `reports/analyses/20260923T0955Z-d2-vs-d1--deepseek-flash--2026-09-23.md`.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Build and run the eight unbuilt designs | **one of eight built and run** (d2); seven recorded as gaps with reasons below |
| Record each unbuilt design as a coverage gap with the reason | met (below) |
| Equal-total-resource arm | **not built** — needs a runner flag splitting the total budget |
| YugabyteDB coverage | **not run** — no d2/YB run; needs the placement mechanism verified first |
| Verified placement with a recorded distribution | **not built** — needs a placement-recording step |
| Correctness gate passes; one matrix at a time | met for the one run performed |

## Coverage gaps (each with the reason it was not built)

| Design | Reason |
|---|---|
| `d3_doc_sections` | needs harness support to assemble a complete answer from several section rows at one revision (a new read path/Kind), not only SQL |
| `d4_doc_jsonb_path` | needs a harness write mode that applies a per-key `jsonb_set`; the document write path currently hands over the whole new document |
| `h1_rows_plus_readview` | needs a new Kind: rows plus a materialized read view maintained in the same transaction |
| `s1_snapshot_pointer` | needs a new Kind: append-only snapshots plus an atomic pointer read |
| `s2_append_history` | needs a new Kind: history plus materialized current state, with drift audit |
| `y1_colocated` | needs the colocation DDL verified on the pinned YugabyteDB 2025.2.6 image, plus placement recording |
| `y2_noncolocated` | same placement recording; meaningful only together with `y1` |
| equal-total arm | needs a runner flag that fixes the total budget and splits it across nodes |
| placement verification | needs a tablet/leader distribution recorded into the run log |

## Notes

The one measured design produced a negative result against its own hypothesis, reported as such.
The seven gaps are harness work, not measurement time, and are the honest remainder of this task.
