---
analysis_id: 20260930T-read-one-worker--deepseek-flash--2026-09-30
run_id: 20260930T232926Z (order A), 20260930T233744Z (order B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; one read worker, and the correction it forced."
analyzed_at: 2026-09-30
inputs_digest: 4706076ef2c6abc4 (order A), 1878ba0fe59000b2 (order B)
inputs: 20260930T232926Z@4706076ef2c6abc4, 20260930T233744Z@1878ba0fe59000b2
repo_commit: 77431c983d42030679345136184f791f241c6c08
status: current
headline: With the mean latency, in-flight is exactly 1.000 in all 72 cells at one worker and exactly 2.00 at two - every read cell in the family is perfectly latency-bound, so the "client-decode-bound" cells of the previous tranche were an artefact of using the median in Little's law on skewed latency distributions (mean/p50 up to 2.7x), and the q02/q05 caveat is withdrawn.
---

# Analysis — one read worker, and the artefact it exposed (study 02) — deepseek-flash

## TL;DR

- **The separation I went looking for does not exist, and the reason is my own statistic.** With the **mean**
  latency, `ops_per_sec × mean` is **1.000 at one worker in all 72 cells** of the four designs and both orders
  (min 0.996, median 0.997) and **exactly 2.00** in every cell at two workers. Every read cell in this family
  is perfectly latency-bound: the store bounds all of them, at both footprints.
- **The previous tranche's "client-decode-bound" cells were a median artefact.** The same cells read 0.34–1.00
  with `p50` at one worker and 0.75–1.92 at two, and the gap is exactly the latency skew: p3's
  `q05_band_tickets_sold` has `p50 = 0.275 ms` and `mean = 0.744 ms` (2.7x), so a median-based in-flight
  count understates the concurrency by the same factor.
- **The explanation I gave was also wrong on its own terms.** I attributed the low cells to the client
  decoding **28 rows**; all three `q05` variants return **one row** (`COUNT(*)`/`SUM(...)`, no `GROUP BY`), and
  `q02` is the only 28-row query in the set. The skew, not the row count, is what moved.
- **Withdrawn**: the registry caveat that `q02`/`q05` "levels are not store numbers". They are store numbers at
  a stated footprint, like every other cell, and the family's levels are therefore store measurements — with
  the footprint attached, and with the concurrency term (1 worker vs 2 vs 8) setting the magnitude.
- **What survives from the tranche**: `I3` and `I4` confirmed (no throttling at half the demand — 0 of ~631
  periods in all eight read phases, client CPU 25.6–41.9 s of a 126 s allowance; no failed cells; both trees
  clean and tagged at `77431c9`), and the *diagnostic*, repaired: use the **mean**.

## What I was looking at

Two whole runs on `pg-single`, scale `small`, phases `verify,explain,read`, published budgets and duration,
four designs (`p3_precreated_cas`, `x1_cas_ledger`, `c1_count_naive`, `r1_inventory_row`) in reversed order,
at `--read-conns 1`. Both tagged at `77431c9`, no failed cells. Compared with the same four designs at two
workers (`20260930T164028Z`, `20260930T164834Z`).

| statistic | one worker (72 cells) | two workers (selected cells) |
|---|---|---|
| `ops/s × p50` | 0.34 – 1.00 | 0.75 – 1.92 |
| **`ops/s × mean`** | **0.996 – 1.000** | **2.00 exactly** |
| latency skew (mean ÷ p50) | 1.0 – 2.7x | 1.0 – 2.7x |

## Conclusions

**Little's law wants the mean, and the corpus's diagnostic has been using the median.** `ops/s × p50` is only
an in-flight count when the latency distribution is tight. Where it is not — the heavier aggregate queries,
where `mean/p50` reaches 2.7x — the product silently measures the *median* concurrency, which looks exactly
like a worker that is idle between fetches. That is what I read last tranche as "the client's row decoding
bounds the loop", and it was available to be caught then: the same JSON carries `mean_ms` next to `p50_ms`.

**One worker is the right instrument for the question, and it gave an exact answer.** At one worker the
product is bounded above by 1 by construction, so agreement with 1.000 cannot be an artefact of the worker
count, and it came out at 1.000 across four designs, nine query cells and two orders. The read phase is a
store measurement, end to end, at both footprints — which also means the levels from the two-worker tranches
are store numbers at a stated concurrency, and only their *magnitude* is concurrency-dependent
(`throughput ≈ workers ÷ mean latency`).

**Nothing in the family's substantive results changes.** The orderings (`r1` ≫ `p3` on `q01`; `c1` ≫ `p3` on
`q05`; the ledger's `r05` at ~32k; `r03`'s gap) were never based on this diagnostic, and the previous
tranche's level comparisons stand; what falls is one explanatory sentence and one registry caveat. The lesson
is the expensive kind: a plausible mechanism ("the client decodes 28 rows") that fit the numbers and was
wrong in both directions — the row count was one, and the residue was skew.

**Prediction verdicts** (recorded in `READ-ONE-WORKER.md` before the runs), read honestly:

| Id | Prediction | Verdict |
|---|---|---|
| I1 | store-bound cells ≥ 0.9 in-flight at one worker | **falsified by the statistic I used** — two `q03` cells read 0.89 with `p50`; with the mean they are 1.00. The prediction was right about the cells and wrong about the formula |
| I2 | `q02`/`q05` ≤ 0.8 (client-bound) | **falsified** — 0.34–0.38 for `p3`/`x1` but 0.95–0.97 for `r1` with `p50`, and 1.00 for all three with the mean: they are store-bound, and my hypothesis was wrong |
| I3 | one worker does not throttle | **confirmed** — 0 of ~631 periods in all eight read phases; client CPU 25.6–41.9 s (about 0.4x the two-worker figures) |
| I4 | correctness green, clean tagged trees | **confirmed** — no failed cells, both runs tagged at `77431c9` |

## Where I think the measurement is weak

1. **Four designs, one topology, one engine, one scale, one host** — though the 72-cell agreement of
   `ops/s × mean` with 1.000 is uniform enough that the statistic's validity does not depend on the sample.
2. **One run per order.** The order agreement is excellent (0.95–1.00 on `q01`), but no level here should be
   quoted from a single run.
3. **The instrument has an upper bound of 1**, so it cannot detect a *store* that is itself the constraint from
   a store that is merely slow; it only says where the worker's time goes.
4. **Throughput at one worker is a latency measurement, not a capacity measurement**, and is not comparable
   with the two- or eight-worker levels beyond the `workers ÷ mean latency` relation.
5. **The mean is itself sensitive to tail outliers** (p99 up to 67 ms in these runs); the median is robust but
   the wrong quantity here. A trimmed mean would be the better instrument if this were carried further.

## What I would measure next

1. **Nothing further on this question** — it is answered, and the answer is that the diagnostic was misapplied.
2. **Re-state the two primary anchors' read cells at footprint 2**, now that the levels may be quoted as store
   measurements with the footprint attached.
3. **If a level is ever quoted from one run, quote the mean latency beside it**, so a reader can check the
   Little's-law relation rather than take the throughput on trust.

## Disagreements with other analyses

It **corrects my own `20260930T-read-family-11-footprint2`**: the sentence "cells whose query returns many
rows are bounded by the client's row decoding (in-flight 0.7–1.5) rather than the store (1.8–2.0)" is wrong in
both parts — the statistic was the median, and neither of the two queries returns 28 rows (`q05` returns one).
It also **withdraws the registry caveat** that clause produced. It agrees with
`20260930T-read-client-footprint`'s mechanism (worker count sets the demand) and with
`20260930T-read-q01-plan-check` on the design, and it leaves every ordering in the family untouched.
