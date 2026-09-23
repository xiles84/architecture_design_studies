---
analysis_id: 20260923T00116Z-controls--deepseek-flash--2026-09-23
run_id: 20260923T00116Z-ctlC-w16
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Author of these control runs."
analyzed_at: 2026-09-23
inputs_digest: 7f3b085773289fb3
repo_commit: f1d9b3c0157949c5ff3671951f97da15363eccaa
supersedes:
status: current
headline: The Study 04 v2 controls settle the writer curve (optimistic wins at 1 writer; the lock costs 1.26x–1.60x from 4 to 64) and narrow the 1.76x claim to a ~1.19x upper bound once retry work is removed, but the within-run identical-SQL position effect remains unmeasured.
---

# Analysis — Study 04 v2 controls — deepseek-flash

> This analysis is the deliverable of `task-20260922T195241Z-study04-v2-controls`, the first of the
> three execution tasks created by `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md`
> (tag `study-04/v2-handoff`). It reads the eight generated reports listed below and their result
> JSON. It edits no earlier analysis. Five earlier contention attempts are recorded as **failed and
> excluded** (section 6).

## TL;DR

- **Control C (writer sweep) is measured and clean.** Under one hot key, the optimistic design is
  ahead at 1 writer (1329 vs 1177 ops/s) and the lock is ahead from 4 writers on: 1.26x at 4,
  **1.58x at 16**, 1.60x at 64. Every cell passed the correctness gate with **0 lost updates**.
- **Control B narrows the v1 1.76x claim.** With retries disabled the optimistic design reached
  **731 ops/s** at 16 writers, *faster* than its default 552 — so the lock's advantage is not retry
  charging alone. Against the retries-disabled optimistic design the lock is **1.19x**, against the
  default it is 1.57x. The v1 "1.76x upper bound" is therefore still an upper bound, and the
  defensible floor for the lock's advantage is ~1.19x.
- **Control A (identical-SQL reproducibility) is partly measured.** The same design and SQL ran in
  three separate phases against fresh databases: read throughput reproduces within **2.1% / 0.9% /
  6.7% / 0.9% / 4.1%**. The v1 run's larger *within-run* identical-SQL spread is a different effect,
  and this control does **not** measure it — see the weak-measurement section.

## 1. What was run

All cells: PostgreSQL 17.11 single node, scale `small`, environment `host-zenbook-ux5406sa`, commit
`f1d9b3c0157949c5ff3671951f97da15363eccaa`, one benchmark lock at a time.

| Run | Design(s) | Writers | Trials | Cells failed | Inputs digest | Run tag |
|---|---|---|---|---|---|---|
| `20260922T2348Z-controlA1` | n1_rows_indexed | – (read) | 3 | 0 | `9ca380472f2033ca` | `run/04-configuration-portal/20260922T2348Z-controlA1` |
| `20260923T0000Z-controlA2` | n1_rows_indexed | – (read) | 3 | 0 | `4643c2a37fa8aaa3` | `run/…/20260923T0000Z-controlA2` |
| `20260923T0005Z-controlA3` | n1_rows_indexed | – (read) | 3 | 0 | `0ae57aca9b145cbe` | `run/…/20260923T0005Z-controlA3` |
| `20260923T0010Z-ctlB-r1` | c1 (retries=1) | 16 | 1 | 0 | `eb28c3ea26ea0522` | `run/…/20260923T0010Z-ctlB-r1` |
| `20260923T0011Z-ctlC-w1` | c1, c2 | 1 | 1 | 0 | `d897600aa583d03a` | `run/…/20260923T0011Z-ctlC-w1` |
| `20260923T0014Z-ctlC-w4` | c1, c2 | 4 | 1 | 0 | `1e5cdb992e027ec8` | `run/…/20260923T0014Z-ctlC-w4` |
| `20260923T00116Z-ctlC-w16` | c1, c2 | 16 | 1 | 0 | `7f3b085773289fb3` | `run/…/20260923T00116Z-ctlC-w16` |
| `20260923T00164Z-ctlC-w64` | c1, c2 | 64 | 1 | 0 | `624799b41ae48039` | `run/…/20260923T00164Z-ctlC-w64` |

(Phases: `verify,read` for A; `verify,read,write,contention` for B and C.)

## 2. Control A — identical-SQL reproducibility

`n1_rows_indexed` run alone in three separate phases, fresh schema and load each time, three trials
per read:

| Read | A1 | A2 | A3 | Spread |
|---|---|---|---|---|
| r01_effective_config | 16020 | 15968 | 15691 | 2.1% |
| r02_read_key | 28633 | 28881 | 28876 | 0.9% |
| r03_revision_check | 46967 | 44143 | 44033 | 6.7% |
| r04_list_installations | 1381 | 1393 | 1384 | 0.9% |
| r05_search_key | 28447 | 28658 | 27540 | 4.1% |

**Reading:** for the same SQL and design, phase-to-phase throughput agrees within ~7% (most reads
within ~2%). This is the cross-phase component of the instrument floor. It is *small* compared with
the v1 run's within-run identical-SQL spread, so those larger differences are **not** explained by
run-to-run instrument drift.

## 3. Control B — retries-disabled contention control

One hot key, 16 writers, one measured contention run per variant; all cells passed the gate with
`lost_updates = 0`.

| Variant | Ops/s | Retries recorded | Conflicts |
|---|---|---|---|
| c1_optimistic_version, default (8) | 552.2 | 1062 | 20218 |
| c1_optimistic_version, retries = 1 | 731.4 | 12000 | 26066 |
| c2_pessimistic_lock, default | 869.7 | 0 | 0 |

**Reading:** disabling retries *raised* the optimistic design's throughput (552 → 731 ops/s), so the
lock's advantage is not manufactured by retry charging. Measured against the retries-disabled
optimistic design, the lock is **1.19x**; against the default it is **1.57x**. Both are below the v1
figure of 1.76x. The retries counters are not directly comparable between variants (the client
records a retry even for the single bounded retry, and gives up sooner), which is itself worth
stating rather than hiding.

## 4. Control C — writer sweep

Same hot key and workload, 1/4/16/64 concurrent writers, one measured contention run per point:

| Writers | c1 (optimistic) ops/s | c2 (lock) ops/s | c2 / c1 |
|---|---|---|---|
| 1 | 1329.2 | 1177.4 | 0.89 |
| 4 | 898.8 | 1133.0 | 1.26 |
| 16 | 552.2 | 869.7 | 1.58 |
| 64 | 465.8 | 745.1 | 1.60 |

**Reading:** the crossover sits between 1 and 4 writers. Below it the version check is cheaper;
above it the lock's waiting beats the optimistic retry storm, reaching ~1.6x and flattening. The
1.76x v1 figure was one point on this curve; the curve is now measured, and no writer count shows
the lock losing. All points carry `lost_updates = 0` and a passed correctness gate.

## 5. Correctness

Every cell in the eight runs above passed the 16-check correctness gate (`gate.passed = true`); no
invariant violation was recorded by a non-control design, and `lost_updates = 0` in all contention
cells. `failed_cells` is empty in all eight manifests.

## 6. Excluded runs (failed cells, kept for the record)

A first attempt ran the contention controls with only `verify,explain,contention`. Those five runs
(`20260923T0001Z-controlB-r1`, `…-controlC-w1`, `…-w4`, `…-w16`, `…-w64`) **failed the gate** with
`INV-3/INV-10: installation N is missing key …`, because the cell's single ledger expects the state
the skipped phases would have produced. Their contention numbers are **not** evidence and are not
quoted anywhere above; the runs are kept, untagged, in `results/`. Re-running with the full phase
sequence (`verify,read,write,contention`) produced the clean cells in section 1.

## 7. Where this measurement is weak

1. **The within-run identical-SQL position effect is not measured.** Control A varied the phase, not
   the cell position inside one multi-cell phase. The v1 run's larger identical-SQL spread came from
   first-versus-last position within one run; the runner cannot currently place the same design at
   both ends of one phase without overwriting its own output file. That remains the open half of
   Control A and is the honest limit of section 2.
2. **Control C points are single contention measurements.** The ≥3-trial rule was applied to the
   read phase of Control A; the sweep used one contention run per writer count, so each ratio has
   the spread of one measurement. A repeated sweep is the natural next step.
3. **PostgreSQL only, one key, one host.** No YugabyteDB arm and no second key; the closed-loop
   generator makes the throughput figures flooring, not SLOs.
4. **Provenance flag.** The seven later runs recorded `repo_dirty: true` and therefore were not
   auto-tagged, because the runner's dirty check counts the untracked generated `reports/` files as
   dirty even though no tracked code changed (verified: `git status --porcelain --untracked-files=no`
   over the code paths is empty). The run tags were created manually on the same commit
   `f1d9b3c0…`; the runner defect is recorded in `LESSONS_LEARNED.md`.
