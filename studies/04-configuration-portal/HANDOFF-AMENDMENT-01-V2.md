# Handoff amendment 01 — Study 04 v2 completion protocol

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-study04-v2-protocol` (HIGH planning), responding to
`task-20260922T025912Z-study04-independent-review`
(`book/evidence/reviews/20260922-study04-independent-review.json`).
**Required tag on integration:** `study-04/v2-handoff`.
**Reads with:** [`HANDOFF.md`](HANDOFF.md) (v1), which this amendment extends and does not replace.

## 0. What this amendment changes

The v1 handoff defined eighteen designs and a full workload catalogue; the v1 run implemented ten
designs at one cardinality, one cadence and PostgreSQL single-node only. The evidence registry
records that gap (`v2-gap-03`) plus three review-mandated controls. This amendment:

1. adds the three required controls (A–C below) to the v1 matrix;
2. defines the cardinality, cadence/churn, repeated-trial, equal-total-resource, placement and
   YugabyteDB phases;
3. fixes the acceptance criteria the v2 report must meet **before** any number is published;
4. splits the work into separately claimable execution tasks (section 7).

It does **not** run anything. No design, dataset or runner change is authorised by this file alone;
each execution task below carries its own acceptance gate.

## 1. The three review-mandated controls (blocking for republishing any ratio)

**Control A — repeated-design instrument control (the ~1.7x read floor).**
Run design `n1_rows_indexed` as both the first and the last cell of a phase, ≥3 repeats, with a fresh
database per cell or a recorded randomized cell order. Purpose: measure the identical-SQL floor
instead of assuming ~1.7x. Acceptance: report the floor as a measured range with the cell order.
The v1 hypothesis that a late cell inherits a dirtier cache is **not** evidence until this runs.

**Control B — `-retries 1` contention control.**
Run the optimistic design (`c1_optimistic_version`) with retries disabled once, under the same 16
writers and one key. Purpose: separate the lock's advantage from extra per-retry harness charging.
Acceptance: the 1.76x figure may be republished only with this control, and then only as "one point
on a writer curve at 16 writers, an upper bound".

**Control C — writer sweep.**
Run `c1_optimistic_version` and `c2_pessimistic_lock` at 1, 4, 16 and 64 writers on one key.
Purpose: turn the single point into a curve and find the crossover. Acceptance: report throughput and
abort/retry rate per writer count; do not extrapolate beyond 64.

## 2. Cardinality phase

Dataset: the v1 deterministic generator (seed 42), extended to tiers **1, 10, 30, 60, 120, 500**
configuration entries per product (v1 ran 3600 rows only, where the whole table is cached).
Designs: `n0_rows_unindexed`, `n1_rows_indexed`, `n3_rollup_trigger`, `n4_rollup_app`.
Deliverable: a crossover table (read latency and ops/s per tier). Acceptance: the unindexed design
must show a scan cost that grows with tier; if it does not, the tier is still cached and the phase
must stop and report that rather than publish a flat line.

## 3. Cadence and churn phase

Cadence definitions remain those of HANDOFF.md §6. Run three regimes at the 120-entry tier:
quiet (1 update/min), steady (1 update/s), burst (30 updates/s for 60 s, then quiet). Workloads:
the v1 `w*` write set plus the overview read. Deliverable: staleness distribution and rollup lag for
trigger versus application maintenance. Acceptance: state the TTL/refresh interval actually crossed
in each regime; a regime in which no expiry or refresh fired is a gap, not a null result.

## 4. Repeated and randomized trials

Every phase in sections 2, 3 and 6 runs **≥3 independent trials per cell**, with a fresh database
load per trial (strongest) or a recorded randomized cell order (weaker, labelled). Report median,
min–max and the trial count; never a single trial where a spread exists. This supersedes the v1
single-trial convention for v2 cells.

## 5. Equal-total resources, placement and YugabyteDB

**Equal-total control.** Alongside the v1 per-node baseline (2 CPU / 3 GiB per database node),
run a separately labelled cell that holds the *total* database budget constant across topologies
(1 node × 2 CPU / 3 GiB vs 3 nodes × 0.67 CPU / 1 GiB). Report it as the equal-total arm, never
pooled with the per-node arm.

**Placement.** The v1 placement designs (`y1_*`, `y2_*`) were never executed. Execute them on
YugabyteDB 3-node RF=3 **and verify actual placement** (tablet/leader distribution recorded in the
run log). Data colocated vs non-colocated must hold logical data, operations, replication and
resources constant. A run without a placement record is an explicit gap.

**YugabyteDB coverage.** Repeat the design set of sections 2 and 6 on YugabyteDB 1-node (RF=1) and
3-node (RF=3). Report per engine; do not mix engines in one ranking.

## 6. Remaining designs

Build and run the eight unbuilt v1 designs (the exact list is HANDOFF.md §3); each must pass the
correctness gate before timing. Any design that cannot be built is recorded as a coverage gap with
the reason, not dropped silently.

## 7. Separately claimable execution tasks

1. **Study 04 v2 controls** — controls A, B and C only. Deliverable: a short controls report; gate
   for republishing the 1.76x and the read floor. (published with this amendment)
2. **Study 04 v2 cardinality and cadence** — sections 2, 3 and the repeats of section 4.
3. **Study 04 v2 remaining designs** — section 6 plus the equal-total control and the YugabyteDB
   placement/engine coverage of section 5.

Each task takes the benchmark lock, runs one matrix at a time, and produces a generated report plus
a signed analysis under `studies/04-configuration-portal/`. None may publish a design ratio that
fails its correcteness gate or lacks its negative control.

## 8. Acceptance criteria for the v2 report

- Three controls present, fired and reported before any ratio.
- ≥3 trials per cell, or a recorded randomized order with the limitation stated.
- Cardinality crossover reported per tier; no flat line where the tier was still cached.
- Cadence regimes state the refresh/expiry actually crossed.
- Equal-total arm labelled and separate; placement verified with a recorded distribution.
- YugabyteDB and PostgreSQL reported separately; failed cells named, never winners.
- Every number resolved to a v2 evidence claim before it enters the book.
