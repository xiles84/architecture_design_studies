---
analysis_id: 20260923T1305Z-staleness--deepseek-flash--2026-09-23
run_id: 20260923T1305Z-staleness-120
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Author of the probe and the run."
analyzed_at: 2026-09-23
inputs_digest: 986addda3f648b7f
repo_commit: d9f63b0c1cf539df5f22b0be4939ba7203bf94db
supersedes:
status: current
headline: Trigger- and application-maintained rollups have no measurable staleness window at 120 entries (p50 ~4 ms, floor-limited by the 2 ms poll); their difference is write cost, not freshness.
---

# Analysis — Study 04 rollup staleness window — deepseek-flash

> Deliverable of `task-20260923T012236Z-study04-v2-staleness-field`, created by the review of the
> Study 04 v2 cardinality/cadence task because the harness recorded **no** derived-value age. It adds
> a staleness probe to the cadence phase and reports the first measurement.

## TL;DR

- **Both maintained-aggregate designs update the aggregate inside the publishing transaction, and the
  measurement shows no staleness window**: n3 (trigger) p50 4.30 ms / max 6.88 ms; n4 (application)
  p50 4.01 ms / max 4.34 ms, over 12 probes each, with **0 timeouts**.
- **The number is at the probe's resolution.** The probe polls every 2 ms, so a sub-millisecond
  window cannot be distinguished from zero; the honest statement is "below roughly 5 ms", not "4 ms".
- **The trigger/application difference is therefore cost, not freshness.** The earlier cardinality
  run showed the trigger's `w05_replace_all` at 10–11 ops/s against the application rollup's
  413–448; here the two are the same freshness. A staleness claim about Study 04's designs would be
  wrong; the staleness case in this repository is the drift control and Study 05's cache.

## 1. What was measured

The cadence phase now, for a design with a maintained aggregate, publishes one child change,
records the commit instant, then polls the design's **own** aggregate through `r04_list_installations`
until it reflects the change (`config_count` increases), recording the elapsed time. It runs at
120 entries per product on PostgreSQL 17.11, `cadence secondly`, three trials on the surrounding
phases.

| Design | Samples | p50 | p99 | max | Timeouts | Audit |
|---|---|---|---|---|---|---|
| n3_rollup_trigger | 12 | 4.30 ms | 5.54 ms | 6.88 ms | 0 | passed |
| n4_rollup_app | 12 | 4.01 ms | 4.25 ms | 4.34 ms | 0 | passed |

Correctness: `gate.passed = true`, `failed_cells: []`, no audit failure, and the probe's own added
keys are in the ledger (the audit after the cadence phase compares them).

## 2. Why the windows are (near) zero

Both designs maintain `ip_rollup` in the same transaction that changes the child rows — n3 through
database triggers, n4 through the publishing statement — so by the time a commit is acknowledged the
aggregate already reflects it. What the probe measures is therefore the observation path: the poll
interval (2 ms) plus one `r04` round trip. That is why the two designs are indistinguishable and why
the absolute figures are a floor rather than a property of the design.

The design that *does* have a window is `x2_rollup_drift_control`, which recomputes the aggregate in
a second transaction after commit; its audit is what measures that, and it is deliberately not
probed here.

## 3. Where this measurement is weak

1. **Resolution.** A 2 ms poll cannot resolve a sub-millisecond window. The result is "no window
   above ~5 ms", not a distribution below it. A tighter probe (event-timestamped, or a single
   round-trip measurement without polling) would be needed to say more.
2. **One regime, one scale, one engine.** `secondly` cadence, 120 entries per product,
   PostgreSQL 17.11 single node. Under a heavier write load the read-then-poll race could behave
   differently.
3. **Aggregate-visible, not answer-visible.** The probe reads the aggregate column, not a full
   application read; a stale snapshot read from a replica (which this study does not have) would not
   be caught.
4. **Twelve samples.** Enough to bound the window, not to characterise a tail.
