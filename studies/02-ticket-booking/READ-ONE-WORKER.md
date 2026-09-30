# One read worker: what bounds each cell (study 02)

**Question.** At the fixed 2-worker read footprint the in-flight count (`ops_per_sec × p50`) separates two
kinds of cell: **1.8–2.0** against 2 workers for the point-lookup and single-row queries (the store bounds
them) and **0.7–1.5** for `q02_band_events` and `q05_band_tickets_sold`, the two queries that return 28 rows
(their loop is bounded by the client's row decoding). This tranche drops the read footprint to **one worker**,
where a single worker can have at most one request in flight, so the same product reads directly as **the
fraction of the client's loop spent waiting on the store** — a measurement of the client-overhead share
rather than a qualitative label.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,explain,read`; published
budgets on the fixed 4-CPU / 2 GiB machine; 5 s per measurement; dataset seed and loader; the four designs of
the footprint-2 tranche (`p3_precreated_cas`, `x1_cas_ledger`, `c1_count_naive`, `r1_inventory_row`) so the
comparison is within the same design set. **The factor varied is the read footprint**: 1 worker now, 2 in
`20260930T164028Z` / `20260930T164834Z`.

**Design.** Two whole runs, the four designs in reversed order, `--read-conns 1`.

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| I1 | **The store-bound cells stay store-bound at one worker**: for `q01` (all five tiers), `q03_customer_tickets` and `q04_ticket_by_id`, in-flight is **≥ 0.9** — at least 90% of each worker's loop is spent waiting on the store. | Any of those cells below 0.9, which would mean the harness's own per-operation work is a tenth or more of the loop even for the cheapest queries. |
| I2 | **`q02` and `q05` are client-bound, and now quantified**: their in-flight is **≤ 0.8**, i.e. at most 80% of the loop is store wait. | Either query at ≥ 0.9, which would make its level a store measurement after all and let the family's `q05` result be quoted as a level. |
| I3 | **One worker does not throttle**: `nr_throttled = 0` in all eight read phases, with client CPU around half the 2-worker figure (≈20–35 s of a 126 s allowance). | Any throttled phase, which would be surprising at half the demand. |
| I4 | **Correctness gates stay green**: no failed cells, both trees clean and tagged at one commit. | Any failed cell or a dirty tree. |

**What a pass decides.** Whether the two many-row queries have a store-bound level at any footprint this
program can measure, or whether — as the previous tranche suspected — their levels are the client's and only
their **orderings** may be quoted. Either answer is recorded in the read-hazard limits, which currently say
`q02`/`q05` levels "are not store numbers".

**Scope and gaps, stated in advance.** One footprint (1 worker), four designs, one topology, one engine,
one scale, one run per order, one host. Throughput at one worker is a latency measurement, not a capacity
measurement: these levels are even less comparable with the 2-worker ones than those were with the 8-worker
ones, and the tranche is about **what bounds a cell**, not about how fast the designs are.
