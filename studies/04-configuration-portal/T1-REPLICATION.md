# Study 04 — T1 replication protocol: the rollup-maintenance pair, three runs

**Written before any measurement.** This file is the `## Prediction` record for a run-level
replication of the strongest isolation pair in the corpus. It is committed before the first run and is
not edited afterwards; verdicts are added to the signed analysis, not here.

Why this pair: `n3_rollup_trigger` and `n4_rollup_app` store identical rollup columns and read them
with byte-identical SQL; they differ only in **who maintains them**. `x2_rollup_drift_control` is
byte-identical to `n4` except that the harness recomputes the rollup after commit, so it is the
negative control for the publication rule. The corpus's claim `v2-13-trigger-rollup-cost` rests on one
run; this protocol is the minimal slice that can raise it above `single_run_directional`.

## Scope

| Field | Value |
|---|---|
| Study | `studies/04-configuration-portal/` |
| Environment | `host-zenbook-ux5406sa` (8 CPUs, 16,496,418,816 B visible in the VM, measured 2026-09-25) |
| Topology | `pg-single` only (the claim is PostgreSQL-only) |
| Scale | `small` (the scale the cited run used) |
| Designs | `n3_rollup_trigger`, `n4_rollup_app`, `x2_rollup_drift_control` |
| Phases | `verify,write` |
| Trials per cell | 3 |
| Duration / warmup | 5 s / 2 s (the runner defaults) |
| Runs | 3 whole runs, each with its own order seed and run tag |

## Prediction

| Id | Comparison | Expected direction | Expected mechanism | Expected crossover / regime | Falsifier |
|---|---|---|---|---|---|
| P1 | `n3` vs `n4` on `w05_replace_all` | n4 materially faster (the cited run records 83.7 vs 593 ops/s, ≈6.9x) | a per-row trigger fires once per entry in the replaced configuration; the application path recomputes once | the gap should shrink as entries per installation fall; not measured here | the three-run median ratio falls inside the report's ±30% suppression floor, or reverses in any run |
| P2 | `n3` vs `n4` on the read phases | no attributable difference | both store the same columns and run byte-identical read SQL | — | any read differs by more than 30% in the three-run median |
| P3 | `x2` drift audit | non-zero drift, control fires in every run | the rollup is recomputed after the commit, so a reader can observe a stale value | — | zero drift recorded in any run (the control did not fire; the invariant is then not demonstrated) |
| P4 | the n3/n4 ratio across runs | the ordering n4 > n3 is stable in all 3 runs | a structural, per-row cost is not a scheduling artefact | — | the sign changes between runs |

## Evidence forms and attribution

- `n3` → `n4` is an **isolate** comparison (who maintains the rollup; schema, SQL and indexes matched).
- `n4` → `x2` is an **isolate** comparison at the extreme: the SQL files are byte-identical and only the
  harness call site moves.
- There is no package in this slice, so any measured difference may be attributed to the maintenance
  decision — provided P2 holds (no read difference) and the controls fire.

## Controls

| Control | Deliberately wrong how | Invariant it must violate | Must fire in |
|---|---|---|---|
| `x2_rollup_drift_control` | rollup recomputed after the mutation commits | the published rollup must agree with its entries | every run, write phase |
| the correctness gate | — | every design must reproduce the Go-computed answers before timing | every run, every cell |

## Boundaries

This slice establishes nothing about: other scales; YugabyteDB; cadence, cardinality or concurrency
regimes; the read gain of the rollup family (`v2-13`'s 15x read gain is a different comparison and is
not re-measured here); or any design outside the three named. A `run_replicated` strength on `v2-13`
covers only the write-cost ratio on `pg-single` at `small` scale.

## Acceptance

- Every cell passes its gate and `x2` records drift in every run; a run in which the control does not
  fire is reported as not demonstrating the invariant and is not pooled.
- The three runs are reported with their own order seeds and per-run medians, plus the across-run
  spread `(max − min) / median` on the ratio.
- `v2-13` moves to `run_replicated` only if P1 and P4 hold and the recorded `prediction` field carries
  this protocol's verdict.
