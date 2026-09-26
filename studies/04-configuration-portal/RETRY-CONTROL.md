# Retry control — separating retry cost from guard cost (Study 04)

**Question.** `v2-17` registers a crossing between a bounded-retry version check
(`c1_optimistic_version`) and a lock-before-update (`c2_pessimistic_lock`), and states its mechanism as
"the version check pays for retries the lock does not". Both `v2-12` and `v2-17` carry the caveat that the
sweep could not separate the retry cost from the guard cost, because the harness retries on conflict and
counts the retries inside the same measurement. This control sets the retry budget to zero and measures
the same two designs at the same writer count, so that the retry work is removed from one side only.

**What the control changes.** One option: the harness's `-retries` flag, default 8, passed through the
runner as `--extra "-retries 0"`. With `Retries = 0` the optimistic path executes its loop once
(`for attempt := 0; attempt <= Retries; attempt++`) and returns "optimistic retries exhausted" on the
first conflict, while the pessimistic path is untouched because it never retries. Everything else —
designs, scale, topology, writer count, trials, phases — is held at the crossing run's values.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | The pessimistic lock is unaffected by the retry budget: at 16 writers its acknowledged throughput with `-retries 0` is within 20% of the 870 ops/s measured with `-retries 8`. | The lock's throughput moves by more than 20%, which would mean the retry budget changes something other than the optimistic path. |
| P2 | The version check's acknowledged throughput at 16 writers collapses by more than 50% against the 552 ops/s measured at `-retries 8`, or the cell fails outright. | Its throughput is unchanged, which would falsify "the version check pays for retries" and require `v2-17`'s mechanism sentence to be withdrawn. |
| P3 | Conflicts remain a property of the workload rather than of the budget: the version check still records more than 10,000 conflicts at 16 writers with `-retries 0`. | Conflicts fall to near zero, which would mean the conflict counter measures retry attempts rather than failed guarded updates. |
| P4 | The correctness gate still passes: the counter's final value equals the acknowledged increment count even though publications were rejected. | The gate fails, which is itself a finding about the gate's sensitivity to rejected publications, not a design result. |

Recorded before measurement, from the committed tree at `a0eac30`. Evidence form: `isolate` — the only
difference between the two arms is the retry budget, and the two designs are compared within one run as
before.

## Scope

| Field | Value |
|---|---|
| Topology | `pg-single` |
| Scale | `small` |
| Designs | `c1_optimistic_version`, `c2_pessimistic_lock` |
| Phases | `verify,contention` |
| Writers | 16 (the crossing run's writer count) |
| Trials | 3 |
| Extra flags | `-retries 0` |

**Boundaries.** One writer count, one host, one key, single node. This control tests the *mechanism
sentence*, not the curve: it cannot show where the crossing moves when the retry budget changes. A
rejected publication is not the same as a slow one, so a collapse in P2 is evidence that throughput at
this contention depends on re-attempting, not a measurement of how much the retries cost.
