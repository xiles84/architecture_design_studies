# T1: replicating `v2-07` — refund answerability and the ledger's storage cost

**Task:** `task-20261001T133759Z-t2-study02-replication` · **branch** `study-02/t1-remaining`

**Goal.** `v2-07` is Study 02's last `single_run_directional` claim, and unlike the study's other claims it
asserts **two schema facts, not timings**: that in 14 of 15 designs the refunds report (`r05`) is
unanswerable and the sales/revenue (`r01`) and last-purchase (`r03`) reports are wrong across a refund; and
that the append-only ledger (`x1_cas_ledger`) that fixes answerability costs about 33–35% more storage. Its own
limit says as much: *"answerability is a schema fact, not a timing"* and *"the +33-35% storage delta is a
relation-size difference … not a repeated timed result"*. A claim whose evidence is a deterministic fact still
rests on one run, so this tranche re-derives the fact in **two whole runs in both design orders** and reports
what the current code and schema produce.

**What is held constant.** The study's own runner, harness, correctness gates, dataset at scale `small`, and
the published budgets (client 2 CPUs / 2 GiB against the fixed host). **Varied: the run and the design order**
(15 designs forward, then 15 in reverse). **Phases: `verify` only** — the load, the correctness gates, the
relation sizes and the report coverage all come from the load path, and no part of this claim depends on a
timing, so running the read/write/race phases would add cost and nothing to the question.

**Measurement definitions, fixed before running** (this is where the primary is ambiguous — its analysis says
"35% storage", which is the tables-only figure, while all relations give 36.2%):

| quantity | rule |
|---|---|
| storage, all relations | `sum(stats.relations[].total_bytes)` for the design |
| storage, tables only | the same sum over `kind == "table"` |
| storage, indexes only | the same sum over `kind == "index"` |
| the delta | `x1_cas_ledger / p3_precreated_cas − 1` |
| answerability | `report_coverage[<key>].status` per design: `answerable` / `partial` / `unanswerable` |

The primary's own data re-derived by these rules: all relations **+36.2%**, tables only **+35.2%**, indexes
only **+37.9%**; `r05` unanswerable in **14 of 15** designs (only `x1` answerable), with `r01` and `r03`
`partial` in the same 14 and `r04`/`r02` answerable everywhere. So the claim's *count* is exact in its own
data and its *percentage* depends on which relations are counted.

## Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage |
|---|---|
| 1. cluster and client sizing, equal-total-budget control | Published budgets, unchanged; the client is not exercised (no timing in this claim). **Not applicable:** one topology, no cluster-size comparison, no equal-total arm. |
| 2. rollup, rolldown and embedding | **Not applicable** to this claim: it is about write-time history (an append-only ledger), not derived state. The ledger is neither a rollup nor a rolldown; embedding is what Study 02's `c*` and `h*` families already do, and this tranche measures only their report coverage. |
| 3. colocated and non-colocated data | **Not applicable** (single-node PostgreSQL; no replication). Outstanding for the YugabyteDB arms, which this tranche excludes. |
| 4. optimistic and pessimistic concurrency | **Not applicable**: no concurrency arm runs here. `v2-07`'s own read side was re-stated on the footprint-2 footing in September and is not re-opened. |

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| M1 | **The answerability count is unchanged**: `r05` is `unanswerable` in exactly 14 of 15 designs in both runs, and `x1_cas_ledger` the only `answerable` one. | Any other count, or a different design answering `r05`. |
| M2 | **`r01` and `r03` stay `partial` in the same 14** and `r02`/`r04` stay `answerable` in all 15, in both runs. | Any status change. |
| M3 | **The storage delta reproduces within ±3 percentage points of the primary's tables-only figure**: between +32% and +38% on every rule (all / tables / indexes), in both runs. | Any rule outside +32…+38% in either run. |
| M4 | **The delta is a property of the design, not of the run**: the all-relations figure differs between the two runs by less than 1% relative. | A run-to-run difference of 1% or more. |
| M5 | **The ledger stays the largest design** and the four pre-created designs stay within 1% of each other (they differ by the counter, not the storage). | The ledger not the largest, or the p-family spread above 1%. |
| M6 | **Correctness gates stay green** in every cell: `verify` passes with no failures, `load_audit` shows 0 over-booked, 0 duplicate seats, 0 ledger mismatches, and both runs are tagged from a clean tree. | Any failed cell or dirty tree. |

**What a pass decides.** `v2-07` moves to `run_replicated` for its **schema half** (two whole runs, both
orders), with its statement re-scoped to name the storage rule that the percentage is measured on and with a
limit keeping its remaining single-run halves visible. A failure in M1/M2 keeps the claim's strength and adds
a limit naming the relation that moved; a failure in M3 does the same for the percentage.

**Known limits, stated in advance.** Specification PostgreSQL only; scale `small`; the answerability statuses
are the harness's own audit of its own schema, so a second run tests reproducibility and not independence of
the audit's rule; and no timing is measured here, so nothing in this tranche speaks to the ledger's load or
read cost.

## The survey: what of Study 02 can this machine still measure?

The task asks for a ranked survey rather than an attempt at everything, because the owner has fixed the
machine at 2 GiB and the study's YugabyteDB arms do not fit.

| # | target | what is un-replicated | cost on the fixed machine | rank |
|---|---|---|---|---|
| 1 | `v2-08`'s **PostgreSQL refund and organiser-edit arms** (2.1–2.5x, up to 5x) | three statement figures rest on the primary run alone, while the race arm now has twelve whole runs | `pg-single`, phases `race,churn`, two runs ≈ 1 h | **1 — highest value per minute** |
| 2 | `v2-09`'s **buyer × tier surface** (10/100/10k/100k seats at 8–64 buyers) | the sweep covers one tier; the axis is monotone at 32 buyers only | `pg-single`, `race`, two runs ≈ 40 min | 2 |
| 3 | `v2-07`'s storage and answerability | this tranche — deterministic facts, no timing | `pg-single`, `verify` only, two runs ≈ 10 min | 3 — done here |
| 4 | `v2-06`'s read-hazard limits | one PostgreSQL read family, re-stated on the footprint-2 footing; its YB arm is not | `pg-single` cheap; `yb-*` **does not fit** | 4 |
| 5 | `v2-11`'s **YugabyteDB read ratio** | the whole YB half | **does not fit**: a cluster node plus the client exceeds the 2 GiB machine | 5 — record as a coverage gap, not a plan |
| 6 | any `yb-cluster3` arm of any claim | — | **does not fit** (the owner's fixed machine) | — |

Items 1 and 2 are the queue's `t2` follow-ups; item 3 is this tranche; items 4–6 are recorded here so the
gap is visible rather than implied.
