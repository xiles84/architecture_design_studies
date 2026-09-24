# Diagram fidelity audit — Study 01 (charity tree)

**Audit revision:** `5b1f8df0c0ef86c7a894b84fa44080127610b9db` (repository state the sources were read at).
**Renderer:** `infra/diagram-render.sh`, `PLANTUML_IMAGE` digest
`sha256:9b9ee6af54a86ea6ac53805e9da4da7fd913711cfe6c19b6ac1a41c398223aeb` (PlantUML 1.2026.8).
**Scope:** every factual annotation in this study's 20 `.puml` sources, traced to the SQL/harness/plan/report
that supports it and classified **current / revision-scoped / stale / unsupported**. Structural field
lists are checked against the design's `schema.sql`; mechanism prose is checked against the design's
`queries.sql`/`indexes.sql`/`triggers.sql` and the cited signed artefact. No measured result, report or
signed analysis was altered.

## Changed sources

| Source | Annotation before | Verdict | Correction and evidence |
|---|---|---|---|
| `d3_flattened_fk.puml` | "six of twelve queries lose a join" | **stale** | At the current sixteen-query catalogue, D3 drops a join in three queries (q02, q08, q15) versus D2 (`sql/d3_flattened_fk/queries.sql` vs `sql/d2_normalized_indexed/queries.sql`). Reworded to "three of the current sixteen queries drop a join (q02, q08, q15)". |
| `d3_flattened_fk.puml` | "'total donated to a charity' becomes an index-only scan, never touching the heap" | **revision-scoped, unscoped** | Zero heap fetches is a prepared-read outcome, not a property of the index. Run `20260913T125342Z-v3`: the covering plan visited 1,756 / 31,293 entries with 0 heap fetches under explicit `VACUUM`, while the ANALYZE-only diagnostic still started at 31,293 heap fetches on trial 1 and fell to 218 as the visibility map advanced (`reports/discussions/d2-d3-mechanisms--gpt-6--2026-09-13.md`). Reworded so the number and its preparation condition travel together. |
| `d10_embedded_hybrid_locked.puml` | "D9's trigger corrupted 3-5 donor caches per 30 s under concurrent reads and writes" | **unsupported unit** | The count is supported by `reports/20260913-d9-cache-race.md` (3 and 5 wrong rows in the 4:4 and 2:6 splits, digest `70f87f5d5667bff8`), but no cited artefact names a 30 s window. Reworded to "3-5 donor caches wrong per concurrent read/write split" with the report and digest named. |

`d10`'s remaining prose (two bugs — an ordering race and a lost update — and their fixes) is supported by
`sql/d10_embedded_hybrid_locked/triggers.sql` (BUG 1 "ordering race (7 of 16 examples)", BUG 2 "lost
update (9 of 16 examples)") and the same report; it is retained unchanged.

## Unchanged sources

| Source | Factual annotation(s) | Verdict | Evidence at the audit revision |
|---|---|---|---|
| `00_overview.puml` | legend D1–D8; five controlled pairs; "seven designs derived" | current | the legend lists D1–D8 and the study's own pair table; "seven derived" counts D2–D8. `sql/*/schema.sql`, `studies/01-charity-tree/README.md` |
| `d1_normalized_minimal.puml` | PKs/FKs only; `donation.person_id` unindexed → full scan | current | `sql/d1_normalized_minimal/indexes.sql` (no secondary index); `schema.sql` comment |
| `d2_normalized_indexed.puml` | three secondary indexes; SQL identical to D1 | current | `sql/d2_normalized_indexed/indexes.sql`; D1/D2 `queries.sql` differ only in the header comment |
| `d4_rollup_trigger.puml` | "seven of twelve queries stop touching the donation table"; hot charity row; trigger indexes | **revision-scoped** (true, not restated as timeless) | for q01–q12 the count is 7; for the current 16-query catalogue it is 11 (`sql/d4_rollup_trigger/queries.sql`). Left as written; recorded here as catalogue-scoped. |
| `d5_rollup_app.puml` | three concurrency strategies; "reads identical to D4 by construction"; `version` columns | current | `sql/d5_rollup_app/{schema,queries,indexes}.sql`; the diagram already shows the two `version` columns |
| `d6_embedded_jsonb.puml` | element shape; GIN cannot order/aggregate across documents; append rewrites the whole row (O(history)) | current | `sql/d6_embedded_jsonb/{schema,indexes}.sql` |
| `d7_yb_child_colocated.puml` | YugabyteDB hash-sharding; per-donor locality; secondary index is a distributed table (two-hop) | current | `sql/d7_yb_child_colocated/schema.sql`, `docs/replication.md`. The illustration "one person's 40 donations" is a hypothetical (the dataset averages ~22 donations/donor), not a measurement. |
| `d8_flattened_nofk.puml` | D3 minus FKs; KEY SHARE parent locks; "same five indexes" | current | `sql/d8_flattened_nofk/{schema,indexes}.sql`; `sql/d3_flattened_fk/indexes.sql` has 5 `CREATE INDEX` |
| `d9_embedded_hybrid.puml` | bounded 20-element cache; other questions use the real table | current; "highest-QPS query in a system like this" is qualitative, not measured | `sql/d9_embedded_hybrid/schema.sql` |
| `d11_copied_key.puml` | D2 reads/indexes plus copied key and FK | current | `sql/d11_copied_key/{schema,indexes}.sql` |
| `d12_recency_index.puml` | D11 plus charity/date index | current | `sql/d12_recency_index/indexes.sql` |
| `d13_recency_sql.puml` | D12 plus recency SQL rewrites | current | `sql/d13_recency_sql/queries.sql` |
| `d14_sum_plain.puml` | D17 plus plain charity index | current | `sql/d14_sum_plain/indexes.sql` |
| `d15_sum_covering.puml` | D14 plus `INCLUDE (amount_cents)` | current | `sql/d15_sum_covering/indexes.sql` |
| `d16_sum_rollup.puml` | D3 plus parent sums only; count/extrema derived on reads | current | `sql/d16_sum_rollup/{schema,queries}.sql` |
| `d17_sum_sql.puml` | D13 plus sum SQL rewrite | current | `sql/d17_sum_sql/queries.sql` |
| `d20_recency_flag.puml` | partial index, one TRUE row per donor; guard locks the person row; delete asymmetry | current | `sql/d20_recency_flag/{indexes,triggers}.sql`; analysis `20260916T090036Z-v3--claude-opus-5--2026-09-16.md` |
| `d22_recency_rollup_idx.puml` | one index on the existing rollup; D22 vs D23 reads byte-identical | current | `sql/d22_recency_rollup_idx/indexes.sql`, `sql/d23_recency_rollup_app_idx/indexes.sql`; same analysis |

## How this was re-rendered

`d3` and `d10` were corrected and this study's diagram set was re-rendered with the pinned helper
(`studies/01-charity-tree/diagrams/render.sh`). Only
`rendered/d3_flattened_fk.svg` and `rendered/d10_embedded_hybrid_locked.svg` changed bytes; the other
18 SVGs are byte-identical, confirming the corrections are the only source change.
