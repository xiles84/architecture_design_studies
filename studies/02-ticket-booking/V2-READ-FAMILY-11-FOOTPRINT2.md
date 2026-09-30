# The remaining eleven read designs at footprint 2 (study 02)

**Question.** Four of the study's fifteen designs now have a quotable read footing at the fixed 2-worker
footprint (`20260930T164028Z`, `20260930T164834Z`). This tranche gives the **remaining eleven** the same
footing, so the read family's tables can be read as one set rather than four quotable designs and eleven
footprint-8 ones whose client was the constraint.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,explain,read`; the published
budgets on the fixed 4-CPU / 2 GiB machine; 5 s per measurement; the dataset seed and loader; the read
default itself (`client_read_conns: 2`, asserted in the manifests). The factor varied between the two runs is
the **design order**.

**Designs (11).** `c2_count_serializable`, `c3_count_lock_event`, `c4_counter_guard`, `c5_seat_unique`,
`h0_hold_naive_confirm`, `h1_hold_checked_confirm`, `p1_precreated_lock_first`, `p2_precreated_skip_locked`,
`p4_precreated_counter`, `r2_inventory_buckets`, `r3_seat_pool`.

**Design.** Two whole runs, the eleven in reversed order:

| # | order |
|---|---|
| 1 | `c2`, `c3`, `c4`, `c5`, `h0`, `h1`, `p1`, `p2`, `p4`, `r2`, `r3` |
| 2 | `r3`, `r2`, `p4`, `p2`, `p1`, `h1`, `h0`, `c5`, `c4`, `c3`, `c2` |

Reference: `20260920T234953Z` (`pg-single`, same eleven plus four others, phases `verify,explain,read`, but
**8 workers**, a ten-day-older harness, and a client whose read phase was throttled 73–719 periods depending
on design). That comparison is **day- and harness-confounded** and is used only as a stated expectation, never
as a scored prediction.

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| H1 | **No design throttles at the new footprint**: all 22 read phases record `nr_throttled = 0` (of ~631 periods), including the five the primary anchor throttled hardest (`c4` 714, `h1` 719, `p4` 712, `h0` 672, `r2` 637 of 1085). | Any throttled read phase, which would mean some designs demand more client CPU than two cores at this footprint and need their own footing. |
| H2 | **Every cell is latency-bound**: `ops_per_sec × p50` is within **1.6–2.4** (≈2 ±20%) in every cell of both orders, as the four-design tranche found. | Products outside that band, which would mean the read phase is not measuring the store for those designs. |
| H3 | **The ordering is not an order artefact**: between the two runs, no pair of these designs differing by **≥1.5x** on `q01`@1,000 swaps places. | A ≥1.5x pair inverting between orders, which would make the set's ordering unquotable. |
| H4 | **The client never needs more than half its quota**: read-phase CPU ≤ **63 s** (one core) for every design in both orders. | Any read phase above 63 s, meaning a design's read mix nearly saturates the two-core client even at 2 workers. |
| H5 | **Correctness gates stay green** in both runs: audit consistent, no failed cells, both trees clean and tagged at one commit. | Any failed cell, or a dirty tree. |

**Expectation, not scored** (day- and harness-confounded): every design's read level should sit **below** its
`20260920T234953Z` value, by a factor of roughly 1.5–4x, because at 8 workers the client's throttle and its
higher in-flight count both inflated those levels — and the five heavily throttled designs should show the
largest falls.

**What a pass decides.** Together with the four-design tranche, all fifteen designs have a footprint-2 read
footing, so the family's tables can be compared within one footing and the read-hazard limits can move from
"four designs re-measured, eleven owed" to a complete statement. It does not make levels portable between
footprints, and it does not re-measure the primary anchors' own cells at footprint 2 (that is a separate
tranche, since it needs the primary's four extra designs' cells as well).

**Scope and gaps, stated in advance.** Eleven designs, one topology, one engine, one scale, one run per
order, one host; the machine is a stated project constraint; the YB read half remains untouched.
