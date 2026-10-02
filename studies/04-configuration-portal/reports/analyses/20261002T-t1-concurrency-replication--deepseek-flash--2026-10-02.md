---
analysis_id: 20261002T-t1-concurrency-replication--deepseek-flash--2026-10-02
run_id: 20261002T162400Z-w1, 20261002T162443Z-w2, 20261002T162526Z-w16, 20261002T162608Z-w64, 20261002T162654Z-w16-r0, 20261002T162737Z-w1, 20261002T162821Z-w2, 20261002T162904Z-w64, 20261002T165329Z, 20261002T165416Z, 20261002T165533Z, 20261002T165607Z, 20261002T165652Z, 20261002T165735Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; T1 replication of v2-12 and v2-17 across a writer sweep with retry controls in both design orders."
analyzed_at: 2026-10-02
inputs_digest: per-run da1a30a97de074bf (w1 c2-first), afb049c763e4f0e9 (w2 c2-first), 22884211c56091ea (w16 c2-first), 61056c05852d204f (w64 c2-first), c38e261dbe5ecf42 (w16 retries 0 c2-first), b9b7c3c225a511e4 (w1 c2-first, outlier), 4d94712efae8dfb4 (w2 c2-first), c26038b1e316b4de (w64 c2-first), 2e97f90971fc1e6c (w1 c1-first), eaf2b2d139e53514 (w2 c1-first), a7adefb1c521d044 (w1 c1-first, pair only), 804bf1be7a4d0de1 (w16 c1-first), 4800ccd6280e6e9a (w64 c1-first), 08278634fcad3d29 (w16 retries 0 c1-first)
inputs: 20261002T162400Z-w1@da1a30a97de074bf, 20261002T162443Z-w2@afb049c763e4f0e9, 20261002T162526Z-w16@22884211c56091ea, 20261002T162608Z-w64@61056c05852d204f, 20261002T162654Z-w16-r0@c38e261dbe5ecf42, 20261002T162737Z-w1@b9b7c3c225a511e4, 20261002T162821Z-w2@4d94712efae8dfb4, 20261002T162904Z-w64@c26038b1e316b4de, 20261002T165329Z@2e97f90971fc1e6c, 20261002T165416Z@eaf2b2d139e53514, 20261002T165533Z@a7adefb1c521d044, 20261002T165607Z@804bf1be7a4d0de1, 20261002T165652Z@4800ccd6280e6e9a, 20261002T165735Z@08278634fcad3d29
repo_commit: 9d47a109
supersedes:
status: current
headline: The lock's advantage at 16 writers replicates in both design orders (1.56x-1.69x) and the retry budget is refuted as its cause (1.1%-3.2%), so v2-12's upper-bound hedge is retired; the negative control reproduces the claim's 93.7% loss to a tenth; and the crossover between one and two writers reproduces too - the 3.85x reading that first looked like an order effect was one unexplained excursion in a single run.
---

# Analysis — T1 replication of `v2-12` and `v2-17` — deepseek-flash

## TL;DR

- **The lock's advantage replicates, and the hedge `v2-12` carried is refuted by its own control.** At 16
  writers `c2`/`c1` came out **1.64x** and **1.69x** with the lock running first and **1.57x** and **1.56x**
  with the version check running first, crossing the retry budget (default against `-retries 0`) with the
  design order. Removing the budget moves the ratio by **1.1% to 3.2%**, so `v2-12`'s "UPPER BOUND ...
  because the run could not separate it from extra per-retry harness work" is no longer true and is retired.
- **The crossover reproduces between one and two writers.** In the version-check-first order the ratio is
  **0.948x** and **0.905x** at 1 writer and **1.035x** at 2; in the lock-first order it is **0.939x** at 1
  and **1.106x**/**1.166x** at 2. The claim's 1.68x at 32 was not re-measured, but the divergence above the
  crossing does reproduce: **1.56x-1.69x at 16** and **1.73x-1.96x at 64**, in both orders.
- **One run at 1 writer is an excursion, and my first reading of it was wrong.** `20261002T162737Z-w1`
  measured 3.85x because `c1` itself read **306 ops/s** there while the other run with the *same* design
  order read **1,286** - not because the order changed. The version-check-first runs added after that first
  draft put the same level at 0.905x and 0.948x. The excursion is reported, not discarded, and no client or
  database throttling explains it.
- **The negative control reproduces the claim's own figure exactly.** `x1_lost_update_control` lost **93.7%**
  and **93.9%** of acknowledged updates at 16 writers (the claim says 93.7%), 50.0% at 2, 97.6%-98.0% at
  64, and **both correct designs lost nothing at every level**.
- **Correctness was green for the compared designs throughout**: every `c1`/`c2` cell passed its gate and
  audit. The only failed cells are the `x1` control at **1 writer**, where a lost update is structurally
  impossible; that is a not-applicable cell and a defect in this tranche's first protocol, not in a design.

## What I was looking at

Fourteen whole runs on `pg-single` at scale `small`, phases `verify,contention`, designs
`c1_optimistic_version`, `c2_pessimistic_lock` and (except where noted) `x1_lost_update_control`. The first
eight were produced by this task's earlier attempt: levels 1, 2, 16 and 64 writers with the lock executed
before the version check, then levels 1, 2 and 64 and a 16-writer zero-retry control, also lock-first. A
check of the manifests, however, showed that **every one of those eight runs executed `c2` before `c1`** -
the shuffle had not produced a reversed pair order - so the tranche's committed protocol, which requires
both design orders at each level, was not actually executed.

This analysis therefore adds six runs that put **`c1` before `c2`** (order seeds chosen so the executed
order is the reversed one, verified in each manifest): 1, 2, 16 and 64 writers at the default retry budget
and a 16-writer zero-retry control. The 1-writer level is run twice in this direction: once with all three
designs (`20261002T165329Z`, whose `x1` cell fails the gate as expected) and once with only `c1`/`c2`
(`20261002T165533Z`), which is fully green.

The study's runner shuffles the design list per run from `--order-seed`, and the manifest records the
executed order. In the table below, **"c2 first" means the lock executed before the version check** and
"c1 first" the reverse. The pair order is the control; the run id order is not.

| run | writers | retries | pair order | c1 ops/s | c2 ops/s | c2/c1 | c1 conflicts | x1 loss | cells failed |
|---|---|---|---|---|---|---|---|---|---|
| 20261002T162400Z-w1 | 1 | default | c2 first | 1,286 | 1,208 | **0.939x** | 0 | 0% (cannot fire) | x1 |
| 20261002T162737Z-w1 | 1 | default | c2 first | **306** | 1,179 | **3.849x** | 0 | 0% (cannot fire) | x1 |
| 20261002T165329Z | 1 | default | c1 first | 1,286 | 1,219 | **0.948x** | 0 | 0% (cannot fire) | x1 |
| 20261002T165533Z | 1 | default | c1 first | 1,328 | 1,202 | **0.905x** | 0 | (not run) | none |
| 20261002T162443Z-w2 | 2 | default | c2 first | 1,128 | 1,247 | 1.106x | 7,885 | 50.0% | none |
| 20261002T162821Z-w2 | 2 | default | c2 first | 1,117 | 1,302 | 1.166x | 7,903 | 50.0% | none |
| 20261002T165416Z | 2 | default | c1 first | 1,231 | 1,274 | 1.035x | 8,562 | 50.0% | none |
| 20261002T162526Z-w16 | 16 | default | c2 first | 704 | 1,152 | 1.637x | 25,583 | 93.7% | none |
| 20261002T165607Z | 16 | default | c1 first | 752 | 1,183 | 1.573x | 27,077 | 93.9% | none |
| 20261002T162608Z-w64 | 64 | default | c2 first | 563 | 1,018 | 1.809x | 25,283 | 98.0% | none |
| 20261002T162904Z-w64 | 64 | default | c2 first | 517 | 1,015 | 1.963x | 23,344 | 97.8% | none |
| 20261002T165652Z | 64 | default | c1 first | 569 | 984 | 1.728x | 25,588 | 97.6% | none |
| 20261002T162654Z-w16-r0 | 16 | **0** | c2 first | 671 | 1,134 | 1.690x | 25,257 | 93.8% | none |
| 20261002T165735Z | 16 | **0** | c1 first | 737 | 1,146 | 1.556x | 26,948 | 93.9% | none |

`x1`'s loss percentage is lost updates over acknowledged updates; at 1 writer there is no concurrent
overwrite, so the control cannot lose one and the gate fails the cell by design (methodology 5a).

| prediction | verdict |
|---|---|
| P1 crossover at the same place in both orders | **falsified by one run** — holds in the c1-first order (0.948x, 0.905x at 1; 1.035x at 2) and in one of the two c2-first runs (0.939x at 1; 1.106x at 2), but `20261002T162737Z-w1` measured 3.85x at 1 writer with `c1` at 306 ops/s. The falsifier ("a ratio above 1 at 1 writer") fired once. |
| P2 the 16-writer ratio ≥ 1.4x in both orders | **confirmed** — 1.637x/1.573x at the default budget and 1.690x/1.556x at zero |
| P3 the retry budget explains < 10% of it | **confirmed at both orders** — 3.2% with the lock first (1.637x against 1.690x) and 1.1% with the version check first (1.573x against 1.556x) |
| P4 `c1` below half its 1-writer throughput at 64, conflicts rising | **confirmed** — 1,328 → 569 and 1,286 → 563 ops/s (0.43x and 0.44x), conflicts 0 → ~25,000 |
| P5 the negative control fires, correct designs clean | **confirmed from 2 writers up** — x1 lost 50.0% at 2 and 93.7%-98.0% at 16 and 64; at 1 writer it cannot fire and the gate fails the cell |
| P6 gates and provenance | **confirmed with the 1-writer exception** — every `c1`/`c2` cell passed its gate and audit, all fourteen runs were tagged from clean trees, and the three `x1` cells at 1 writer are reported as failed and not applicable |

## Conclusions

**`v2-12` keeps its direction and loses its hedge.** The four fresh 16-writer ratios (1.56x, 1.57x, 1.64x,
1.69x) bracket the primary's 1.76x from below, and the control the claim said it could not run has now
been run in **both design orders**: at a matched order the retry budget changes the ratio by 1.1%-3.2%,
and `c2`, which never retries, moves by the same amount as `c1`. The sentence calling the 1.76x an upper
bound because of per-retry harness work is therefore removed, and the claim records the ratio as a band of
five whole runs on two days. This also resolves the tension the registered claims carried between them:
`v2-17`'s limit had already recorded that the retry explanation failed, and `v2-12` still doubted it.

**`v2-17`'s crossing is reproduced and is not order-dependent.** In the version-check-first order the
1-writer ratio is 0.948x and 0.905x, and in the lock-first order it is 0.939x (with one exception), so the
crossing sits where the primary placed it. The single 3.849x reading comes from a run in which `c1`
measured 306 ops/s while another run with the same pair order measured 1,286 - a 4.2x movement in the
*version check's own* throughput at the level where neither design contends. Because the later,
order-controlled runs give 0.905x and 0.948x there, this is an unexplained single-run excursion, **not**
the order effect the first draft of this analysis reported. The claim's four- and eight-writer levels were
not re-measured, so its 32-writer point still rests on the primary alone.

**The negative control is the cleanest replication here.** 93.7% lost at 16 writers is the claim's figure
to a tenth. The control also exposed a protocol defect worth recording: because a 1-writer run has no
concurrency, `x1` cannot break the invariant there, and including it made three of the tranche's runs
report a failed cell. Those cells are not-applicable, and the tranche's own 1-writer evidence now comes
from a run of the pair alone plus the version-check-first run.

## Where I think the measurement is weak

1. **One run at level 1 is unexplained.** `20261002T162737Z-w1` moved `c1` from 1,286 to 306 ops/s with the
   same pair order, the same settings and no recorded throttling (client and database cgroup counters both
   zero). It is reported, not smoothed: it is why P1 is scored falsified by one run rather than confirmed.
2. **One fresh-load trial per cell.** No level has a within-run spread; only the order pair and the retry
   pair give any repeat at all, and the two retry comparisons each hold order fixed while the two orders
   are single readings.
3. **The 32-writer point was not re-measured**, and the primary's 4- and 8-writer levels were not either, so
   the curve is the primary's between 2 and 16 writers.
4. **The `x1` control cannot fire at 1 writer**, which is a coverage gap in the primary sweep too (the
   primary ran no control at all).
5. **Scale `small`, one hot key, one host, PostgreSQL single node**, and absolute throughput is day-specific:
   only within-run ratios travel. The declared budgets still exceed the 2 GiB podman machine (the standing
   environment blocker), so no container limit is enforceable.

## What I would measure next

1. **The 1-writer level with repeats in three positions**, to characterise the 306 ops/s excursion the way
   the corpus did for Study 02's `c1_count_naive`.
2. **The 32-writer level**, so the curve is covered at every point the claim names.
3. **The 4- and 8-writer levels**, which the primary measured once and nothing has repeated.
4. `t5` (Study 05's cache and churn claims) remains the next queued replication.

## Disagreements with other analyses

None on the data. This tranche **resolves a tension inside the registry**: `v2-12`'s upper-bound wording is
superseded by the retry control that `v2-17`'s limit already cites, and the re-scoped statement says so.
