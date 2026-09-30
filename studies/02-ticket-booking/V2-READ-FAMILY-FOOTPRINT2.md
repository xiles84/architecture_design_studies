# Re-measuring the read family at footprint 2 (study 02)

**Question.** The 2026-09-30 calibration fixed the read phase's client footprint at **2 workers** (client
unthrottled, every cell latency-bound at `ops/s × p50 = 1.9`). The four designs of the 2026-09-29 re-run were
measured at the old footprint — 8 workers, a client throttled for two thirds of its periods — and the three
read-hazard limits say their read figures are **not quotable**. This tranche re-measures those four at the new
footprint, so the family has one footing, and asks whether the *ordering* the family is read for survives the
change.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,explain,read`; the published
budgets (`DB_CPUS=2`, `CLIENT_CPUS=2`) on the fixed 4-CPU / 2 GiB machine; duration 5 s per measurement; the
dataset seed and loader; designs `p3_precreated_cas` (the O(tier) baseline), `x1_cas_ledger` (whose "2x slower
across the board" the re-run retired), `c1_count_naive` (the negative control `q05` turns on) and
`r1_inventory_row` (the O(1) availability design). **The factor varied between this tranche and the previous
one is the read footprint** — 2 workers now, 8 then — and the two orders within this tranche.

**Design.** Two whole runs, the four designs in reversed order:

| # | order | read workers |
|---|---|---|
| 1 | `p3`, `x1`, `c1`, `r1` | 2 (the runner's default) |
| 2 | `r1`, `c1`, `x1`, `p3` | 2 |

Compared with `20260929T165533Z` / `20260929T170353Z` (same designs, same phases, **8** workers) and with
`20260930T101521Z` (the default-path verification, `p3` and `r1` at 2 workers).

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| G1 | **The client stays unthrottled at the new default** in both runs, for all four designs: `nr_throttled = 0` in every read phase, as the calibration and the verification run found for two designs. | Any throttled read phase, which would mean the four-design set costs more client CPU than the two-design set did, and the default needs re-checking per design set. |
| G2 | **Every cell is latency-bound**: `ops_per_sec × p50 ≈ 2` (against 2 workers) in both designs of both orders, within ±20%. | Products far from 2 (say above 4), which would mean the client is again the constraint at this footprint. |
| G3 | **The ordering the family is read for survives**: `r1` ≫ `p3` on `q01` at 1,000 seats, and `c1` ≫ `p3`/`x1` on `q05` — the two design facts the footprint-8 runs established. | Either ordering inverting or collapsing to parity. |
| G4 | **The footprint-8 levels do not reproduce, and that is the point**: `r1`'s `q01@1000` lands near 20k (the calibration's 2-worker figure) rather than the 53k the 8-worker run read, and `p3`'s tier-100,000 cell stays in the 600s (store-bound at any footprint). | `r1` at ~53k, which would mean the run did not use the new default. |
| G5 | **Correctness gates stay green**: audit consistent, no failed cells, both trees clean and tagged. | Any failed cell. |

**What a pass decides.** It gives four of the fifteen read designs a quotable footing at a stated footprint,
and it decides whether the three read-hazard limits can be relaxed from "the read figures are not quotable" to
"not quotable *at footprint 8*; re-measured at footprint 2 for four designs, the remaining eleven owed". It
does not re-measure those eleven, and it does not make levels portable between footprints.

**Scope and gaps, stated in advance.** Four designs of fifteen, one topology, one engine, one scale, one run
per order, one host; the machine is a stated project constraint; the YB half of the read family is untouched
and the 2 GiB machine governs whether it can be measured at all.
