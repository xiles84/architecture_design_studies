# Self-review — Edition 2 revision against the external review (2026-09-24)

**Attribution.** Written by `deepseek-flash` (DeepSeek, HIGH-capability session, Deep Code CLI; effort
setting not exposed) on 2026-09-24. This is the executor's final self-review of the revision requested by
the owner; it is **not** an independent review, and the reader should weigh that.

**Artefact reviewed.** `book/dist/data-architecture-reference.pdf`, commit of the release task,
61 pages, 29/29 claims indexed, 6/6 fonts embedded, evidence digest present, no unverified banner, no
unresolved template tokens. Tags: `repo/data-architecture-book-visual-correctness-a`,
`repo/data-architecture-book-visual-tranche-b`, `repo/data-architecture-book-v2-1`.

## Item-by-item status

| Ref | Item | Status | Evidence |
|---|---|---|---|
| A | Figure 6 (lease vs fencing token) redraw and semantic distinction | **Completed** | `fig_lease_vs_fence.puml`: 5 lifelines; coordinator no longer publishes to the cache; two panels + a distinction note separating lease / sink-side fencing token / source-generation validation. Text equivalent and caption updated. |
| B | Figure 7 (stale-fill race) version provenance and precise naming | **Completed** | `fig_cache_stale_fill.puml`: explicit `S0/v0 → S1/v1`; named source-generation validation; R2 authoritative read; caption and text equivalent state the stale-fill-race scope only. |
| C | Strict vs relaxed freshness "bounded window" | **Completed** | One stale site existed (`chapters/minor-variants.typ`); the concept chapter and glossary were already correct. Now all three agree. |
| D | Rollup description terminology | **Completed** | `chapters/major-families.typ` (same transaction vs publication) and `chapters/how-to-choose.typ` (write claim scoped to the measured workload). |
| E | Information-placement text-equivalent contradiction | **Completed** | `concepts/information-placement.typ`. |
| F | Pessimistic-locking categorical language | **Completed** | `concepts/concurrency-control.typ`; 1.76x stays scoped. |
| G | Append-only history / ledger concepts | **Completed** | `concepts/derived-state-and-history.typ` (derived state vs retained history; fundamental claim) and `chapters/major-families.typ` ("append and fence" removed). |
| 2.1 | Normalized vs rolldown vs rollup figures | **Completed** | `fig_families_er` (three ER panels with cardinalities) and `fig_rollup_vs_rolldown` (direction comparison). |
| 2.2 | Embedded documents | **Completed** | `fig_embedded_vs_normalized`. |
| 2.3 | Append-only history state-vs-history | **Completed** | `fig_history_vs_current`. |
| 2.4 | Cache as external read copy | **Completed** | `fig_cache_read_copy`. |
| 2.5 | Indexes and hot rows | **Completed** | `fig_index_vs_hotrow`. |
| 2.6 | Optimistic vs pessimistic concurrency | **Completed** | `fig_optimistic_vs_pessimistic`, with the invariant recheck named. |
| 2.7 | Derived state vs history | **Completed** | `fig_derived_vs_history`. |
| 2.8 | Topology / distributed placement | **Completed** | `fig_topology_physical_boundary`, including the one-physical-host boundary and the UNVERIFIED physical-placement band. |
| 2.9 | Benchmark-interpretation decision flow | **Rejected, with reason** | The reading chapter keeps its checklist as a table (searchable, cannot be letter-spaced) and the prose now states that a gap is not a zero effect. Adding a competing flowchart would duplicate it. |
| 3 | Cache/terminology audit across the book | **Completed** | Eight terms distinct in `glossary.typ`; the publication-fence card and matrix row name source-generation validation; the chapter intro separates it from a sink-side token and cache CAS. |
| 4 | Editorial / layout pass | **Partially completed** | New figures were placed with concept → diagram → takeaway → evidence card ordering; two tall sequences render below full width rather than increase density; the 61-page growth is recorded. A whole-book whitespace rebalance was not performed. |
| 5 | Validation pass | **Completed with a stated limit** | Every new/changed figure re-read against its prose; arrows, lifelines, version progression and cardinalities checked; no general rule inferred from one result; every retained number still maps to its claim; all pages rasterised and the figure pages inspected; the release was rebuilt through `book/build.sh`. The limit: the visual inspection is one model's, not an independent read. |

## Where this is weak

- **No independent review or reader test.** Both queue tasks were self-reviewed by the same session.
- **Two figures render below full width** (74% and 72%) to avoid footer overflow; legible in the raster
  but not confirmed with print readers.
- **Item 2.9 was rejected**, and item 4 is partial by design; both are recorded rather than smoothed over.
- **The book is now 61 pages**, up from 54. No figure was dropped, and no accepted item was traded away
  for page count.
- **`book/evidence/v5/` was not touched**, so every measured value and coverage gap is unchanged.

## Reproof

The revision is reproducible from the merged state: `book/assets/export.sh --check`, `book/evidence/check.sh`
and `book/build.sh` all pass; the figure manifest records the pinned renderer and per-asset hashes.
