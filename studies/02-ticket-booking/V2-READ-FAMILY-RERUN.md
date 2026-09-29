# Re-measuring the read family after the client-quota audit (study 02)

**Question.** The 2026-09-29 client-quota audit found that the **read phase** of `20260920T234953Z` — the run
behind study 02's read tables and two signed analyses — was **client-CPU throttled by a design-dependent
amount**: 2.3 s for `p2_precreated_skip_locked` against **58.1 s** for `r1_inventory_row`, out of ~1085 CFS
periods, with `p4_precreated_counter` (56.3 s), `c4_counter_guard` (53.9 s) and `r2_inventory_buckets`
(36.5 s) also heavily throttled. No claim's *support cell* uses a read number from that run, but a signed
analysis quotes its read throughput (`r03` 176 ops/s, `q03` 32k ops/s, and `q05` 5.5k on C1 against 1.7k on
P3/X1), and the same analysis flagged "X1's buyer reads came out ~2x slower across the board … it is a
cell-level effect" without being able to name the cause. This tranche re-measures those cells on the current
harness, where the client's throttling is recorded per trial, to see whether they are quotable and whether
the X1 read penalty was the design or that cell's conditions.

**What is held constant.** Topology `pg-single`; scale `small`; the same dataset seed and loader; the same
phase set as the original (`verify,explain,read` — the original run had no writes or races); the same read
queries and their 10-second measurement windows. The factors varied are **which designs are measured in the
same session** and their **order**.

**Design.** Two whole runs, four designs each, order reversed between them:

| # | order |
|---|---|
| 1 | `p3_precreated_cas`, `x1_cas_ledger`, `c1_count_naive`, `r1_inventory_row` |
| 2 | `r1_inventory_row`, `c1_count_naive`, `x1_cas_ledger`, `p3_precreated_cas` |

`p3` is the baseline, `x1` the ledger whose read penalty is in question, `c1` the negative control the `q05`
comparison turns on, and `r1` the design the old run throttled hardest (58.1 s).

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | **Today's client is not throttled in these cells.** Every read cell of every design in both runs records `nr_throttled = 0`, which is the point of re-running: the old numbers were taken with the client losing up to half its periods. | Any read cell with client throttling, which would mean the hazard is not a period effect and the read phase needs a design change (fewer, longer-lived connections), not a re-run. |
| P2 | **The `q05` ordering reproduces**: sell-through per band is much faster on C1 than on P3 or X1 (the old run: 5.5k against 1.7k), because a design with no per-seat rows has no per-seat aggregate to compute. | C1 not clearly fastest on `q05` in both runs. |
| P3 | **The two-orders-of-magnitude gap reproduces**: within X1, `q03_customer_tickets` is ~100x `r03` (old: 32.4k against 176 ops/s). | The gap collapsing below 20x or inverting. |
| P4 | **X1's "~2x slower across the board" read penalty does not reproduce**, on the same cells, when the client is not throttled: X1's `q01_*` and `q04_ticket_by_id` cells come within ±25% of P3's median across the two orders. This is the prediction that would attribute the old penalty to that cell's conditions rather than to the ledger. | A consistent ~2x penalty in both orders at the same cells, which would make it a property of the ledger's read path after all. |
| P5 | **Correctness gates stay green** in both runs: audit consistent, no failed cells. | Any failed cell, which voids that run. |

**Scope and gaps, stated in advance.** PostgreSQL single-node only: the original also measured `yb-single`
and `yb-cluster3`, where the analysis quoted 5.3 → 29.9 and 7.2 → 42.2 ops/s on `r03`, and those are **not**
re-measured here (the 2 GiB machine ceiling is the owner's open decision). This tranche says nothing about
whether the read numbers of the other eleven designs in the original run were distorted — only that the four
measured here are or are not. One tier set, one scale, one host, one session per order.

## Scope

| Field | Value |
|---|---|
| Runs | two, `pg-single`, four designs, `--phases verify,explain,read --scale small` |
| Compared with | `20260920T234953Z` (`pg-single`, all 15 designs, same phases) |
| Unit of evidence | per-cell `ops_per_sec` and its recorded client CPU, per design and order |
| Boundary | One engine, one topology, one scale; the YB half is a stated gap |
