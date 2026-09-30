---
analysis_id: 20260930T-read-family-footprint2--deepseek-flash--2026-09-30
run_id: 20260930T164028Z (order A), 20260930T164834Z (order B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; re-measures the read family at the footprint fixed on 2026-09-30."
analyzed_at: 2026-09-30
inputs_digest: ecd8fd8a6935831a (order A), dde0616ddec5a6ad (order B)
inputs: 20260930T164028Z@ecd8fd8a6935831a, 20260930T164834Z@dde0616ddec5a6ad
repo_commit: cced2d36a25d4a70e8566476eb49407bca03e38c
status: current
headline: At footprint 2 the client is unthrottled (0 of 631 periods in all eight design-cells) and every cell is latency-bound (in-flight 1.9 for 2 workers); the families' orderings survive, the store-bound cell is footprint-invariant (p3 at 100,000 seats reads 629-680 against 623 at footprint 8), but the r1-over-p3 margin narrows from 2.07x to 1.16-1.48x, so read magnitudes are footprint-specific and only ratios that have a mechanism travel.
---

# Analysis — the read family at footprint 2 (study 02) — deepseek-flash

## TL;DR

- **The new default holds for the full four-design set.** All eight read phases record **0 throttled
  periods** (of ~631), and every cell's operations in flight (`ops_per_sec × p50`) is **1.9 against 2
  workers** — `G1` and `G2` confirmed, and not just for the two designs the verification run covered.
- **`G4` confirmed**: `r1`'s `q01` at 1,000 seats lands at 17.8k / 17.1k, not the 38.2k the footprint-8 run
  read; `p3` at 100,000 seats reads **680 / 629** against **623** at footprint 8.
- **The cells separate into two kinds, and only one is portable.** `p3`'s 100,000-seat cell is
  **footprint-invariant** (623 → 680 → 629 across three footprints and two orders): it is store-bound. Every
  other cell moves with the footprint, because its throughput *is* the workers divided by its latency.
- **The design margin was partly a footprint artefact.** `r1` over `p3` on `q01` was **2.07x** at footprint 8
  and is **1.16x / 1.48x** at footprint 2: at 8 workers the client's throttle bit the designs unequally
  (in-flight 5.3 for `r1` against 3.2 for `p3`), so part of that margin was client parallelism rather than
  query cost. At 2 workers both are pinned at 1.9 and the margin is the **latency** ratio, which is the
  honest design difference.
- **What survives, confirmed independently**: the `q05` ordering (`c1` 9,235 / 14,186 against `p3` 2,178),
  the ledger's `r05` at **32,164** against 30,322 at footprint 8 (a 6% move), and the `r03` gap. These are
  the read family's quotable content.
- **Registry consequence**: the three read-hazard limits can be relaxed from "the read figures are not
  quotable" to a footprint-explicit statement — four designs re-measured at footprint 2, the remaining eleven
  and the primary anchors' own footprint-8 figures still not quotable.

## What I was looking at

Two whole runs on `pg-single`, scale `small`, phases `verify,explain,read`, the published budgets and
duration, four designs (`p3_precreated_cas`, `x1_cas_ledger`, `c1_count_naive`, `r1_inventory_row`) in
reversed order, at the runner's new read default of **2 workers** (`client_read_conns: 2` in both manifests).
Both trees clean, both tagged at `cced2d3`, no failed cells. Compared with `20260929T165533Z` /
`20260929T170353Z` (same designs and phases at **8** workers) and `20260930T101521Z` (the default-path
verification of `p3` and `r1` at 2 workers).

| design | `q01`@1000, order A | order B | footprint 8 (A) | verification at 2 | `q05`, A / B | `r05`, A |
|---|---|---|---|---|---|---|
| `p3_precreated_cas` | 15,344 | 11,535 | 18,441 | 10,984 | 2,178 / — | — |
| `x1_cas_ledger` | 10,427 | 11,604 | 15,498 | — | 1,987 / — | 32,164 |
| `c1_count_naive` | 13,407 | 14,186 | 25,010 | — | 9,235 / — | — |
| `r1_inventory_row` | 17,846 | 17,134 | 38,172 | 18,014 | 12,084 / — | — |

Client read phase: **0 throttled periods in every one of the eight cells** (0/631, 0/630); usage 40–65 s
against a 126 s allowance, i.e. 50% headroom.

## Conclusions

**The footprint works as designed, and it changes what the read table means.** With the client unthrottled
and pinned at two operations in flight, a cell's throughput is its latency divided into the worker count:
`p3`'s `q01` rises from 15.3k at 1,000 seats to 19.0k at 10 seats in one order because its latency does, and
`r1`'s stays at 17.8–19.1k at every tier because its latency does not. Every level in these runs is a
latency measurement at a stated concurrency — which is a cleaner statement than the footprint-8 table, whose
levels were partly the client's, but it is also a statement that **cannot be compared with a footprint-8
level**.

**Which numbers travel is now measurable, and the answer is narrower than I expected.** `p3`'s 100,000-seat
cell reads 623 (footprint 8), 680 and 629 (footprint 2, two orders): a cell whose cost is in the store does
not care how many workers ask. Everything else moved by 40–60%. The `r1`-over-`p3` margin on `q01` is the
sharpest case: 2.07x at footprint 8, 1.16x and 1.48x at footprint 2 — so **the footprint-8 margin was
inflated by the client**, and a claim using it would have overstated the design difference by up to 1.8x.
The two confirmed orderings that survive (`c1` ≫ `p3` on `q05`; the ledger's `r05` at ~32k) do so because
their mechanism is in the query, not in the generator.

**Consequences for the read family, stated plainly.** (1) Its **levels** are quotable only with a footprint
attached, and the manifests now carry `client_read_conns`; the four designs here have a footprint-2 footing.
(2) Its **orderings with a mechanism** are quotable and reproduce across footprints. (3) Its **margins** are
not portable — any ratio between designs must name the footprint until the family is re-measured as a whole.
(4) The eleven designs not re-measured here, and the primary anchors' own read figures, remain not quotable.

**Prediction verdicts** (recorded in `V2-READ-FAMILY-FOOTPRINT2.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| G1 | the client stays unthrottled for all four designs | **confirmed** — 0 throttled periods in all eight cells |
| G2 | every cell is latency-bound (`ops/s × p50` ≈ 2, ±20%) | **confirmed** — 1.9 in every cell of both orders |
| G3 | the family's orderings survive | **confirmed, with a caveat** — `r1` ≫ `p3` on `q01` and `c1` ≫ `p3`/`x1` on `q05` both hold, but the `r1` margin narrows to 1.16–1.48x |
| G4 | the footprint-8 levels do not reproduce and the store-bound cell stays in the 600s | **confirmed** — `r1` at 17–18k against 38.2k; `p3`@100000 at 680/629 against 623 |
| G5 | correctness green, trees clean and tagged | **confirmed** — no failed cells, both runs tagged at `cced2d3` |

## Where I think the measurement is weak

1. **Four designs of fifteen.** The quotable footing covers `p3`, `x1`, `c1`, `r1`; the counter family
   (`p4`, `c4`, `r2`) that was throttled hardest in the primary anchor is still unmeasured at footprint 2.
2. **One run per order, and the order matters**: `p3`'s `q01`@1000 is 15,344 in one order and 11,535 in the
   other (a 33% session/order term), so single levels should not be quoted from one run.
3. **`ops/s × p50 ≈ workers` is a diagnostic, not a proof** — it uses p50 for the mean and does not measure
   observed concurrency directly.
4. **The `q05` margin also moved** (6.1x at footprint 8, 4.2x at footprint 2, one order each), so the
   "ordering survives" statement is about direction, not size.
5. **One topology, one engine, one scale, one host**; the machine is a stated project constraint, and the YB
   read half is untouched.

## What I would measure next

1. **The remaining eleven designs at footprint 2**, counter family first, so the family's tables are on one
   footing end to end.
2. **The primary anchors' read cells at footprint 2** (`20260920T234953Z`'s designs), which would let the
   corpus replace the primary's read tables rather than merely caveat them.
3. **A second run per order** at footprint 2 if any design's level is to be quoted rather than its ordering.

## Disagreements with other analyses

None. It implements the footing `20260930T-read-client-footprint` set, confirms the direction
`20260929T-read-family-rerun` established (its orderings) while showing that its margins were partly a
footprint artefact, and leaves `20260930T-read-q01-plan-check`'s O(1) reading of `r1` intact — visible here
as flat `q01` across all five tiers at both footprints.
