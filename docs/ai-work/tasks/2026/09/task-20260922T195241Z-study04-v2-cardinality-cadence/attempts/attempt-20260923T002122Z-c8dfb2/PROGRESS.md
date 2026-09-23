# Progress — Study 04 v2 cardinality and cadence

**Task:** `task-20260922T195241Z-study04-v2-cardinality-cadence`
**Attempt:** `attempt-20260923T002122Z-c8dfb2` (claim `claim-6019513121be64fc`, epoch 1)
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-cardinality-cadence`

Measured: nine PostgreSQL cells (6 cardinality tiers, 3 cadence regimes); one benchmark lock at a time.

## Decision Log

1. **One run per cardinality tier** because the runner takes a single `--entries-per-ip`; tiers
   1/10/30/60/120/500, each with n0/n1/n3/n4 and three read trials. Confidence: High.
2. **Cadence regimes map to the harness's cadence options**: quiet = `minutely`, steady =
   `secondly`, burst = `secondly` with `--burst-size 30`, at 120 entries per product. Confidence: High.
3. **All nine cells passed the gate** (`failed_cells: []`); no cell was excluded. Confidence: High.
4. **The cardinality hypothesis was tested, not assumed**: the unindexed design does show a growing
   scan cost (r04 11.2k → 188 ops/s), so the phase legitimately ran to 500 rather than stopping.
   Confidence: High.
5. **Staleness and derived-state rollup lag are not recorded by the harness.** The cadence phase
   reports scheduling lag and delivery counts, not the age of a derived value at read time. This
   half of the cadence acceptance is delivered as a named gap, not as a fabricated number.
   Confidence: High.
6. **Runs were auto-tagged** this time: the dirty-check fix from `study-04/v2-controls` is in this
   worktree's runner, so `--tag` worked after the first run. Confidence: High.

## Residual risks

- The staleness/rollup-lag field needs a harness change and is the honest limit of the cadence arm.
- Cardinality reads and cadence writes are separate runs; tier 1 is a degenerate point.
