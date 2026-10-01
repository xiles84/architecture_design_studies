# All fifteen read designs in one session (study 02)

**Question.** The read family's footprint-2 coverage is complete **per design** but not **jointly**: the four
designs of `20260930T1640*` and the eleven of `20260930T20*` were measured in different sessions, so a
comparison spanning them carries a session term (10% on `q01` within a run pair, and up to 33% in an earlier
session). This tranche measures **all fifteen designs in one run, twice, in reversed order**, which turns
per-design coverage into a jointly comparable table and re-tests the order stability on the full set at once —
including `p3`/`x1`/`c1`/`r1` against the other eleven.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,explain,read`; published
budgets on the fixed 4-CPU / 2 GiB machine; 5 s per measurement; dataset seed and loader; the read default
(`client_read_conns: 2`, asserted in the manifests). **The factor varied is the design order.**

**Designs (15).** `c1_count_naive`, `c2_count_serializable`, `c3_count_lock_event`, `c4_counter_guard`,
`c5_seat_unique`, `h0_hold_naive_confirm`, `h1_hold_checked_confirm`, `p1_precreated_lock_first`,
`p2_precreated_skip_locked`, `p3_precreated_cas`, `p4_precreated_counter`, `r1_inventory_row`,
`r2_inventory_buckets`, `r3_seat_pool`, `x1_cas_ledger`.

**Design.** Two whole runs, the fifteen in reversed order.

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| J1 | **Nothing throttles**: `nr_throttled = 0` in all 30 read phases, as at this footprint for the two subsets. | Any throttled phase. |
| J2 | **Every cell is store-bound**: with the **mean** latency, `ops_per_sec × mean` is **2.00 ± 5%** in every cell of both runs (270 cells). | Any cell outside that band, which would reopen the question the one-worker tranche settled. |
| J3 | **The joint ordering is order-stable**: between the two orders, no pair of the fifteen differing by ≥1.5x on `q01`@1,000 swaps. | A ≥1.5x pair inverting. |
| J4 | **The previously session-confounded comparison now holds**: `r1` (17.1–17.8k in its own session) and `p4` (18.4–19.3k in the other) agree to within **10%** in this single session, i.e. the gap that looked like a design difference is a session difference. | `r1` and `p4` still separated by more than 10%, which would make it a design difference after all. |
| J5 | **The `q05` separation reproduces within one session**: the pre-aggregated/bounded designs read **≥3x** the per-seat `COUNT(*)` designs (`p1`, `p2`, `p3`, `x1`). | A ratio below 3x in either order. |
| J6 | **Correctness gates stay green**: no failed cells, both trees clean and tagged at one commit. | Any failed cell or a dirty tree. |

**What a pass decides.** It replaces the "complete per design, not jointly" caveat in `v2-07`'s limit with a
single-session table, and it is the last read measurement this study needs on PostgreSQL. It does not touch
the YugabyteDB arms (no replacement exists on the fixed machine) and it does not re-measure any other phase.

**Scope and gaps, stated in advance.** One topology, one engine, one scale, one host, one run per order; the
machine is a stated project constraint. Throughput remains a latency measurement at a stated footprint, and
levels remain comparable only within one footprint.
