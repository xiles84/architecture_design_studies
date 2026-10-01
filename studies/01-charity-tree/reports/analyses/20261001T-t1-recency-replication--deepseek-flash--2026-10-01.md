---
analysis_id: 20261001T-t1-recency-replication--deepseek-flash--2026-10-01
run_id: 20261001T144451Z-t1 (pass 1), 20261001T145914Z-t1 (pass 2)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; T1 replication of v2-04 and v2-05, which also found and fixed a runner defect."
analyzed_at: 2026-10-01
inputs_digest: 57459a8619328b01 (pass 1), ab92855a27c3f94b (pass 2)
inputs: 20261001T144451Z-t1@57459a8619328b01, 20261001T145914Z-t1@ab92855a27c3f94b
repo_commit: 49dde3bd
supersedes:
status: current
headline: Both remaining Study 01 claims replicate - D3 beats D2 in all four fresh loads across two whole runs at 1.59-1.84x (primary 2.47x), and the recency flag loses to the rollup index on exactly the queries the primary found (q13 1.9-2.3x, q15-trailing 2.6-3.1x, tied on q14, marginally behind on q16); the claim's quoted 2.9x is its best query, and the group-wide ratio is about 1.4x in both the primary and the replication. A runner defect that silently ran one cell per run was found and fixed on the way.
---

# Analysis — T1 replication of `v2-04` and `v2-05` — deepseek-flash

## TL;DR

- **`v2-05` replicates.** `d3_flattened_fk` beats `d2_normalized_indexed` on the core read score in **all four
  fresh loads** across two whole runs — **1.59x, 1.72x, 1.64x, 1.84x** — against the primary's **2.47x** over
  its five loads in one session. The direction is robust; the magnitude came down by about 30%.
- **`v2-04` replicates as a *query-specific* result, and its quoted figure is its best case.** On the recency
  queries, `d22_recency_rollup_idx` beats `d20_recency_flag` in every trial (group geomean **1.37–1.53x**), but
  the per-query pattern is the same one the primary saw: **q13** 1.9–2.3x (primary **2.88x**, and 2.71x in the
  other regime), **q15 trailing** 2.6–3.1x (primary 2.73x), **q14 tied** (0.98–1.24x), **q16 marginally
  behind** (0.85–0.94x). So "an index on the rollup column answers the same question faster" holds for the
  window questions and not for all of them, and the group-wide ratio is **1.43x / 1.44x** here against
  **1.52x** recomputed from the primary's own data.
- **Both claims are registered `run_replicated`** with statements re-scoped to what travelled: the `v2-04`
  statement now names the query dependence and the ~1.4x group ratio instead of the 2.9x headline, and the
  `v2-05` statement carries the reproduced band.
- **A runner defect was found and fixed: `run-enhancements.sh` silently ran one cell per run.** Its cell loop
  reads the plan from a file, and **every podman invocation on this host consumes the shell's stdin**, so the
  loop ended after the first cell while the manifest still said "failed processes 0". Proven with
  `podman --version` alone (1 of 3 loop iterations without a private descriptor, 3 of 3 with it). The plan is
  now read on fd 3. The published matrices are not affected — their manifests count 100 and 192 cells — but any
  future run with this runner would have been truncated to one cell and reported as complete.
- **`L4` was falsified, and `L5` was not measurable.** `L4` required the `v2-04` group ratio to be ≥1.5x; it is
  1.37–1.53x, and my floor contradicted what the primary's own data gives (1.52x). `L5` predicted `d1` would be
  the slowest of five designs, but the scoped plan never contained `d1` — the runner's `reads` block does not
  list it — so the prediction was untestable as written. Both are my errors, not the claims'.

## What I was looking at

Two whole runs of the study's own v3/v4 protocol at its own settings (scale `small`, `duration 3s`,
`warmup 1s`, the published budgets), five designs, `pg-single`, conditions `reads` (cell name
`vacuum-small`) and `recency-reads`, two trials per run so the runner's AB/BA alternation gives both orders
inside each run. **18 cells per run**, both trees clean at `49dde3b`, both runs tagged by the runner, no failed
cell, no read or write error in any cell.

| comparison | primary | pass 1 | pass 2 |
|---|---|---|---|
| `d3`/`d2`, core read score | 2.47x (5 loads, one session) | **1.59x / 1.72x** | **1.64x / 1.84x** |
| `d22`/`d20`, recency group geomean | 1.52x (recomputed from its data) | 1.44x | 1.43x |
| `d22`/`d20`, q13 trailing (the quoted 2.9x) | 2.88x | 1.92x | 2.33x |
| `d22`/`d20`, q15 trailing | 2.73x | 3.06x | 2.55x |
| `d22`/`d20`, q16 (both regimes) | 0.89x | 0.85–0.94x | 0.89–0.94x |

The two earlier runs of this tranche (`20261001T142606Z-t1`, `20261001T142837Z-t1`) are **diagnostic runs, one
cell each**, kept in `results/` as the evidence that exposed the runner defect. They are not cited as
replication evidence.

## Conclusions

**`v2-05` is now run-replicated and its statement can keep its numbers.** Its assertion is an ordering plus a
survival claim ("survives five fresh loads"), and four further fresh loads across two whole runs in reversed
order reproduce both: `d3` above `d2` every time, by 1.59–1.84x against the primary's 2.47x. The insert-cost
half of that claim stays as imprecise as its own limit says, and this tranche does not touch it.

**`v2-04` replicates, but the claim should stop leading with 2.9x.** The primary's 2.9x is the q13-trailing
comparison — the best of the eight recency measurements — while its own group-wide ratio is 1.52x and mine is
1.43–1.44x. The honest statement is a query-dependent one: the rollup index wins the window questions
(1.9–3.1x), ties the count question, and is marginally *slower* on the lapsed-donor question (0.85–0.94x,
reproduced in both regimes and both runs). A reader choosing between the flag and the index needs that shape,
not one headline number. The re-scoped statement says so, and the 2.9x/7.9x pair moves into the limit that
names its query and its engine.

**The runner defect is the most consequential thing this tranche produced.** A truncated run that reports
"failed processes 0" is exactly the failure mode the correctness-first rule exists to prevent, and it was
found only because a scoped plan of 60 cells produced one. Two lessons are recorded: any loop whose body calls
podman must not use the shell's stdin, and **a run whose cell count is far below its plan should be treated as
a failure**, not as a fast run.

**Prediction verdicts** (recorded in `T1-RECENCY-REPLICATION.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| L1 | `d3` beats `d2` in every trial and both runs | **confirmed** — 1.59x, 1.72x, 1.64x, 1.84x |
| L2 | `d3`/`d2` between 1.2x and 3x | **confirmed** — all four inside |
| L3 | `d22` beats `d20` on the recency group in both runs | **confirmed** — 1.44x and 1.43x, and every trial above 1 |
| L4 | the `v2-04` speedup between 1.5x and 10x | **falsified as stated** — 1.37–1.53x; my floor contradicted the primary's own 1.52x |
| L5 | `d1` slowest of the five on the core score | **not measurable** — the scoped plan never contained `d1`; the prediction should have named a design the runner's `reads` block lists |
| L6 | gates green, clean tagged trees, no failed cells | **confirmed** for both runs (18 cells each, no failures, both tagged at `49dde3bd`) |

## Where I think the measurement is weak

1. **`pg-single` only.** `v2-04`'s YugabyteDB half — the quoted 7.9x — is unreplicated, and the fixed 2 GiB
   machine is why.
2. **Five designs of twenty-five, and `d1` absent**, so the tranche cannot say what the baseline did.
3. **The comparison with the primaries is day-confounded** for absolute levels; only the ratios and the
   per-query pattern are used.
4. **One trial per order per run** (`--trials 2` gives AB/BA but only one measurement per order).
5. **The scoped filters are new code** in `run-enhancements.sh`; the unfiltered path is unchanged, but a
   filtered run is a labelled subset and the manifest says so.
6. **The runner defect's blast radius was checked only by cell counts** in the published manifests (100 and
   192 cells, so those runs completed); no rerun was done to confirm their per-cell contents.

## What I would measure next

1. **The YugabyteDB arms of both claims**, when a machine can carry them — the last half of `v2-04` and the
   only half of most Study 01 claims.
2. **A scoped re-run including `d1`**, which would repair `L5` and give the tranche a baseline.
3. **An audit of other study runners for the same stdin pattern** — only `run-enhancements.sh` reads a loop
   from a file today, but the fix belongs in the shared habits, not in one script.

## Disagreements with other analyses

None. It re-states and does not contest `20260916T090036Z-v3--claude-opus-5--2026-09-16`: that analysis's
2.9x is correct **for q13 trailing**, and this tranche's contribution is to reproduce the query dependence it
also observed (its own text notes the flag is "2.9x slower ... q13 trailing") and to put the group-wide ratio
beside it. It leaves `20261001T-t1-read-replication` (the `v2-01`–`v2-03` replication) untouched.
