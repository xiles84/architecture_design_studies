# Superseded join/CTE artefacts (round 2)

Moved here by `deepseek-flash` (LOW) on 2026-10-05 under corrective handoff 2.
Nothing is deleted or edited; each entry names its replacement.

| Superseded file | Why | Replacement |
|---|---|---|
| `20261005T-select-join-cte-v3--deepseek-flash--2026-10-05.md` (v3) | Review 3: its "non-colocated" cluster arm had one tablet per relation with all five leaders on one node, so it measured a single-tablet arm, not a distributed one; and its P4 verdict inverted the per-run signs. | `reports/analyses/20261005T-select-join-cte-v4--deepseek-flash--2026-10-05.md` — a pre-split, leader-spread arm behind a placement gate, with the P4 signs taken per run. |
| `select-join-topology-subset-20261004.md` | Its compatibility class omitted the node budget, so it pooled the 2 CPU/3 GiB and 6 CPU/9 GiB single-node runs into one six-run median (review 2, M3). It was also built on the hash-sharded YB rank index, so every "selective" cell was a full scan (review 2, M1). | `reports/select-join-v3b-topology-20261005.md` — classes separate budget, topology, placement, phase and index arm, on the range-index arm, with per-run endpoint/writer/retry/throttling evidence. |
| `20261004T-select-join-cte--deepseek-flash--2026-10-04.md` (v1) | Reported a synthetic-fixture result and a mis-derived reuse pair; rejected in review 1. | Superseded by v2, then v3, then `reports/analyses/20261005T-select-join-cte-v4--deepseek-flash--2026-10-05.md`. |
| `20261004T-select-join-cte-v2--deepseek-flash--2026-10-04.md` (v2) | Its YugabyteDB mechanism named no index arm; the YB rank index was hash-sharded, so the PG-vs-YB materialization contrast was an index artefact (review 2, M1). Its "1%" topology labels and its rollup/rolldown/embedding non-applicability reason were also wrong (M3, M6). | `reports/analyses/20261004T-select-join-cte-v3--deepseek-flash--2026-10-05.md`, on the matched range-index arm; topology superseded again by v4. |

The generating runs of the superseded files remain in
`studies/01-charity-tree/results/` with their producing tags; they are retained
as labelled observations (hash-index arm, pre-fix client build) and are not
pooled with the v3b arm.
