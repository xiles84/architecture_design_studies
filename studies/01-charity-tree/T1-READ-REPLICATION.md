# T1: replicating Study 01's read-score claims (run-level)

**Goal.** Study 01's `v2-01`, `v2-02` and `v2-03` are the corpus's oldest surviving claims and all three are
`single_run_directional`, anchored to one 2026-09-12 small-scale run with **a single measurement per
configuration**. This tranche re-runs that experiment as **whole runs, twice, in reversed design order** — the
T1 pattern the programme uses — and asks whether the *ratios and orderings those claims quote* reproduce.

**The claims under test** (statements as registered):

| claim | statement (abridged) |
|---|---|
| `v2-01-normalized-index-first` | "the first large read win is the query index: D2 indexed scores **13.3x** D1 minimal's read score; flattening the foreign key instead (D3 flat+FK) reaches **20.5x**" |
| `v2-02-rollup-family` | "parent rollups win the aggregate read score (D4 rollup/trg **73.8x** and D5 rollup/app **69.9x** against D1 minimal) and pay for it on inserts (**2.6k/s** and **2.1k/s** against **6.2k/s**)" |
| `v2-03-embedding-not-ahead` | "embedding the donation child in the donor row (D6 embedded) never came out ahead of the normalized designs on reads: **2.4x** D1 minimal's read score, below both indexed designs, with a **0.57x** insert rate" |

**What is held constant.** The primary anchor's own configuration: **scale `small`**, `duration_per_op 10s`,
the same design set (`d1`, `d2`, `d3`, `d4`, `d5`, `d6`, `d8`, `d9`), the same write ops, the same dataset
parameters, study 01's own harness (`-cmd full`) and its correctness gate. The environment is the same host
and the same published budgets. **The factors varied are the run and the design order**: one run with the
designs in `d1…d9` order, one reversed.

**Design.** Two whole runs on `pg-single` only. The primary anchor's YugabyteDB arms (`yb-single`,
`yb-cluster3`) are **out of scope** and stay unreplicated: the podman machine has 2 GiB and the owner has
fixed the machine, so a YB replication is not measurable here. That is a stated coverage gap, not a result.

## Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage in this plan |
|---|---|
| 1. cluster and client sizing, equal-total-budget control | Budgets are the study's published ones (`DB_CPUS=2`, `CLIENT_CPUS=2` on a 4-CPU host) and are unchanged from the primary anchor. **Not applicable:** this plan compares one topology against its own past, not cluster sizes, so there is no equal-total-budget arm to add; the sizing arithmetic for the host is the documented one. |
| 2. rollup, rolldown and embedding | **Covered, and this is the plan's core.** `d4_rollup_trigger` and `d5_rollup_app` (rollup), `d6_embedded_jsonb` and `d9_embedded_hybrid` (embedding), each against `d1`/`d2`/`d3`; read scores and insert rates both reported. Rolldown (copied parent information on children) has no Study 01 design and is **outstanding**. |
| 3. colocated and non-colocated data | **Not applicable here** (PostgreSQL single node, no replication, no placement mechanism) — and **outstanding** for the YugabyteDB arms, which this plan excludes. |
| 4. optimistic and pessimistic concurrency | **Not applicable** to a read-score/insert-rate comparison with no concurrency phase; Study 01's concurrency work is in its own protocols and is untouched by this plan. |

## Prediction

Recorded before measurement. Following the corpus's own lesson that **ranges travel and point figures do
not** (the 2026-09-26 replication moved every absolute by 26–72% while its range reproduced), the predictions
name bands and orderings, not the quoted numbers.

| Id | Prediction | Falsifier |
|---|---|---|
| K1 | **The read-score orderings reproduce in both runs**: `d3` > `d2` > `d6` > `d1`, and `d4`/`d5` above `d2`. | Any inversion of those four relations in either run. |
| K2 | **The quoted bands survive**: `d2`/`d1` ≥ **6x** and `d3`/`d1` ≥ **10x** (against the quoted 13.3x and 20.5x), and `d4`/`d1`, `d5`/`d1` ≥ **25x** (against 73.8x and 69.9x). | Any of those four ratios below its band, in either run. |
| K3 | **Embedding stays not-ahead**: `d6`/`d1` between **1.5x and 6x** and below both `d2` and `d3` (against the quoted 2.4x). | `d6` at or above `d2`, or outside the band. |
| K4 | **The write penalty is real and ordered**: `d1`'s insert rate is the highest, and `d4`, `d5`, `d6` are each **below half** of it (against 0.42x, 0.34x and 0.57x). | Any of `d4`, `d5`, `d6` at or above half of `d1`'s insert rate. |
| K5 | **The two runs agree with each other**: for each of the four ratios above, the two runs' values lie within a **factor of 2** of each other, and no ≥1.5x pair of designs swaps between the runs. | A ratio differing by more than 2x between the runs, or a ≥1.5x pair inverting. |
| K6 | **Correctness gates stay green**: every cell passes study 01's own verification, no failed cells, both runs from a clean committed tree, and each run tagged after the fact (study 01's runner predates `--tag`). | Any failed cell, a dirty tree, or an untagged run. |

**What a pass decides.** If K1–K6 hold, the three claims have **two whole runs in reversed order** with a
scored prediction, which is the T1 evidence the programme asks for, and the registry can move them from
`single_run_directional` to `run_replicated` — with the YugabyteDB gap and the small scale still named. If
they do not, the claims keep their strength and gain a limit saying which relation failed to reproduce.

**Known limits, stated in advance.** Small scale; `pg-single` only; one engine; one host; two runs; the
read score is Study 01's own blended measure over several queries, so a claim about it is a claim about that
blend rather than about any single query.
