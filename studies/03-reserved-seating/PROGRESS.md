# Study 03 — execution progress

Executor log for [HANDOFF.md](HANDOFF.md). Facts only: what was done, which commit, which
mapped decisions were taken and why. Interpretation belongs to the analysis phase.

**Executor:** Claude Opus 5 (`claude-opus-5`), Claude Code desktop. Sessions from 2026-09-14.

## Steps

| Step | What | State | Commit / tag |
|---|---|---|---|
| 1 | Scaffold, README, PROGRESS | done | `967bc82` |
| 2 | Study 02 terminology (Part A) | done | `study-02/v1.2-terminology` |
| 3 | Engine probe | done — every answer as assumed, no fallback used | `53e277e` |
| 4 | Platform: ErrLockNotAvailable | done — study 02 image still builds | `241e06f`, `repo/platform-lock-not-available` |
| 5 | SQL catalogue, 13 designs (14 with S1r, AM-01) | done | `13941ff` |
| 6 | Harness (17 unit tests pass in the build) | done | `eb7e4ff` |
| 7 | Diagrams | done — five diagrams rendered | `8b6d8f9` |
| 8 | Dev checks and calibration | dc1–dc11 done; ER-01 raised, decided (AM-01) | `0a9b2da`; ER-01 decision `e70c842`, `study-03/v0.1-handoff-amendment-01` |
| 8a | AM-01: transient-refusal class, refusal diagnostics, S1r, report; dev checks dc12–dc14 | done — §10.3 as amended by AM-01.5 holds (19 unit tests pass in the build) | `e10dc6c`, `7ab4b67`, `f43fc0f`, `03d7e8a`; tag `study-03/v1-harness` |
| 8b | AM-02 calibration (dc15, S1 `small`, 3 topologies) | done — passes; rule 1 applies, projection ≈ 15.7 h (rule 2) | `3d7aa89` |
| 9 | Main matrix (`small`) | done — 41 cells, 4 failed (diagnosed; ER-03); ER-04 raised | `3e36cfb`; run tag `run/03-reserved-seating/20260915T002411Z` |
| 10 | Repeated race trials | done — 27 cells, 0 failed, 3 h 33 min (from worktree `.worktrees/study03-measurement`) | `3d1d1b0`; run tag `run/03-reserved-seating/20260915T173255Z` |
| 11 | Context, lessons, README; ready for analysis | done — merged into `main`; ER-03 and ER-04 open (non-blocking) | tag `study-03/v1-measured` |

## Mapped decisions (handoff §11)

| # | Decision | Reason |
|---|---|---|
| M1 | Every statement that grants, extends, releases, sells or refunds returns `clock_timestamp() AS db_clock` beside `now() AS db_now`; the ledger replays a seat's events in `db_clock` order and judges expiry-based steals by `db_now`, release-based ones by order. | `now()` is the transaction start (probe P1). A hold that waited for a buyer's release lock carries a `now()` older than the release, so judging by `now()` alone would report a false theft. Unit test `TestGrantThatWaitedForARelease`. Implements INV-2 exactly; the invariant is unchanged. |
| M2 | K0's naive `w_confirm_seats` also sets `customer_id = $customer_id`. | Without it, confirming a seat the sweeper had already made available violates the sold-row CHECK and errors, so the control would fail by error instead of overselling. Code that "already checked" names the buyer too. |
| M3 | L1's loaded sales get their `'infinity'` claims from a `w_load_sold_claims` statement the loader runs after the bulk copy (the explain phase skips `w_load_*`). | pgx cannot encode or scan `'infinity'` into Go time (probe P4); writing it in SQL keeps the design's representation exactly as specified. |
| M4 | E2's `w_expired_carts` returns `event_id` as well; the E2 sweeper attributes released seats to carts through the ledger. | Needed to record which hold each swept seat belonged to. |
| M5 | L2 hold decision uses the read's `now()`; the ledger's grant clock is the compare-and-set write's `clock_timestamp()`. | As §6.12 specifies; recorded for clarity. |
| M6 | S0 lists both `race` and `lifecycle` as experiments where it is expected to fire. | §3.5 names both; §10.3's pass criterion checks the race. |
| M7 | The guard margin G is 2 minutes of real time in the race and isolated writes (whose holds last a real 40 minutes) and 2 human minutes in the lifecycle. | §3.4 defines G as 2 human minutes = 5% of the TTL; this keeps it 5% of the TTL on both clocks. |
| M8 | Buyer seeds are deterministic per (event, buyer), not time-based. | Same random stream in every design, as §3.10 requires. |
| M9 | A lost compare-and-set (L2) backs off like engine contention (jittered, 2 ms base, 50 ms cap). | §3.10 names engine errors only; a lost write is the same contention to a buyer. Without it L2 spun: 21 seats/s on 10-seat races in dc3, 206 in dc4. |
| M10 | Finite-pool write ops (hold, release, cancel) warm up for one trial's worth of operations instead of a duration. | A duration warmup drained the tiny release pool before the first trial (dc2: 0 releases). Keeps warmup plus trials inside half the pool, as `FiniteBudget` intends. |
| M11 | Loaded holds' base clock is the database now() truncated to milliseconds. | The gate failed S1 by 1 ms (dc1): microsecond now() plus millisecond offsets vs millisecond truth. |
| M12 | The sweeper outage ends when the outage probe has run, not at minute 85 by wall clock. | dc3: the sweeper resumed on its own tick first and released every expired seat, so E1's probe saw nothing (0/0). dc4: E1 4/4 and 9/9 unavailable. |
| M13 | Release lag counts only holds granted during the phase, not holds loaded already expired. | Harness bug: the 10-seat lifecycle tier reported ~45 000 human minutes of lag in every lazy design (dc3, dc12), the age of the load. |

## Observations for the analysis phase

Facts recorded during dev checks and runs (handoff §13).

- **Calibration, PostgreSQL (dc3, tiny, S1 lifecycle):** hold p99 4.8–8.5 ms, confirmation p99 5.4–6.3 ms against a 50 ms limit (1 human minute). `-human-minute 50ms` kept.
- **Calibration, YugabyteDB 1-node (dc7, tiny, S1 lifecycle):** hold p99 76–103 ms, confirmation p99 50–71 ms against a 500 ms limit. `-human-minute 500ms` kept.
- **dc5 (YugabyteDB 1-node, tiny, 13 designs, commit `1024aa7`):** every gate 51/51. Controls fired: S0 (race, lifecycle), K0 (lifecycle: late sales), E0 (lifecycle: 8 and 48 thefts). **ER-01:** single early rejections in S1, S2, K1 (10 000-seat race) and S3 (1 000-seat race) — confirmation UPDATE matched 0 rows of a hold the same transaction then read as valid; reproduced in dc8–dc10, not in P13 or on PostgreSQL. See ESCALATIONS.md.
- **S4 (SERIALIZABLE) on YugabyteDB 1-node:** race 0.1 seats/s with hold latency p50 1–38 s and every tier timed out; lifecycle confirmed almost nothing (errors: attempt deadline). No invariant violation. Engine probe P7 showed YugabyteDB aborting serializable conflicts with 40P01.
- **L2 (section document) on PostgreSQL** in dc3 before M9: 21 seats/s on 10-seat races with p50 hold 58 ms (lost compare-and-set spinning); with backoff (dc4) 206 seats/s.
- **dc3 (PostgreSQL, tiny, 12 designs, commit `1b2801d`):** every gate 51/51; no violation in any design meant to be correct, in any phase. Controls: S0 fired in the race (3 857 thefts, 2 009 early rejections, 89 deferred confirmations refused) and the lifecycle (418 thefts); E0 fired in the lifecycle (66 thefts, 13 early rejections); K0 fired in the lifecycle (14 late sales, 10 sales without the hold). E1's outage probe did not fire because of a harness race (M12), fixed and confirmed in dc4.
- **S0's unconditional overwrite of a seat that is already sold is refused by the shared `event_seat` CHECK constraint** (held rows must have `sold_at` NULL): those attempts are errors (183 in the 10-seat race), while overwrites of held seats succeed and show as theft.
- **K1 refused no confirmation in dc3's lifecycle (0 late), S-designs 1–13 late per tier; no design recorded a boundary rejection at tiny scale on PostgreSQL.**
- Engine probe `results/devchecks/engine-probe-20260914T112716Z.md`: on YugabyteDB a
  SERIALIZABLE read-then-update conflict aborted one session with 40P01 ("deadlock
  detected ... Consider using READ COMMITTED"), on PostgreSQL with 40001. Both are retried.

## Environment notes

- 2026-09-14: the Podman machine was stopped (the host had restarted); started with
  `podman machine start`. 8 CPUs and 15.4 GiB are visible, as the environment page records.
  Images from other projects (`dnakit-*`) exist on this Podman machine. No containers of
  theirs were running when this study's work started.
- 2026-09-14 10:47–~11:26 UTC: the study 01 v3 session held the benchmark lock
  (`20260914T104721Z-v3 (memory-control)`). Study 03 waited, wrote code, and started no
  image builds or databases until it was released.

## AM-01 dev checks (step 8a)

Commit `e10dc6c` (harness) plus `7ab4b67` (release-lag fix, M13). `tiny` scale; 5 s measurements.

| Run | Topology | What ran | Result |
|---|---|---|---|
| dc12-pg-am01 | PostgreSQL | S1, S1r, K1, L1, E2, all phases | every gate 51/51; no violation; 0 early rejections; S1r retried 7 confirmations of expired holds, 0 sold |
| dc13a-yb1-s1r-all | YugabyteDB 1 node | S1r, all phases | gate 51/51; no violation; lifecycle retried 1 expired-hold confirmation, 0 sold |
| dc13b-yb1-race4 | YugabyteDB 1 node | S1, S1r, K1, L1, E1; race x4 | early rejections, all transient: S1 2, E1 5, K1 10, L1 4; S1r 0, with 6 retries, 6 sold; deferred confirmations 100% |
| dc14a-yb3-s1r-all | YugabyteDB 3 nodes | S1r, all phases | gate 51/51; no violation; race retried 17 short confirmations, 17 sold; lifecycle 2 expired-hold retries, 0 sold |
| dc14b-yb3-race2 | YugabyteDB 3 nodes | S0, S1, S1r, E0, E1, K1, L1, L3; race x2 | every gate 51/51; S0 fired (646 thefts, 355 early rejections, 19 transient); correct designs' early rejections all transient: S1 23, E1 15, K1 28, L1 38, L3 33; S1r 0, with 20 retries, 20 sold; E0 (no race control) 19 transient; deferred confirmations 100% in every correct design |

Observations (facts): no persistent early rejection appeared in any correct design; transient refusals
concentrate in the 1 000- and 10 000-seat races and are more frequent on three nodes; every
confirmation S1r retried after a short match sold on the second attempt.

## Duration projection (§10.4, AM-01.7)

Base: the AM-01 dev-check cells at the runner's default 5 s measurements (dc12, dc13a, dc14a),
per-cell reload times from `reload_ms` (dc3, dc5, dc11), scaled to `small`. Assumptions: `small` has
about 12× the seats of `tiny` and loads scale linearly; each race tier is capped by the 3-minute
tier budget plus in-flight events (60 s race timeout); YugabyteDB lifecycle keeps 2 tiers × 3 events;
PostgreSQL lifecycle runs 3 tiers × 6 events; everything else (verify, explain, reads, writes,
audits) grows 1.5×. Uncertainty about ±40%.

| Topology | `tiny` cell (measured) | `small` cell: loads + race + lifecycle + rest | Designs | Main matrix |
|---|---:|---|---:|---:|
| pg-single | 4.0 min | 1 + 5 + 2 + 4 ≈ 12 min | 13 | ≈ 2.6 h |
| yb-single | 13.0 min | 14 + 13 + 3 + 9 ≈ 39 min | 14 | ≈ 9.2 h |
| yb-cluster3 | 14.7 min | 26 + 15 + 3 + 9 ≈ 53 min | 14 | ≈ 12.4 h |
| **Total** | | | | **≈ 24 h** (threshold 14 h) |

Repeated race (step 10: pg-single and yb-cluster3, `verify,race`, 3 trials): pg-single ≈ 16 min × 13 ≈
3.5 h; yb-cluster3 ≈ (5 min load + 2 min verify + 3 × (4 min reload + 15 min race)) ≈ 65 min × 14 ≈
15 h; **total ≈ 18.5 h** (threshold 12 h). Both exceed §10.4: Escalation Required ER-02.

## AM-02 calibration (step 8b)

`devchecks/dc15-small-calibration`: S1, `small`, all phases, commit `8b5b0ef`; 2026-09-14 23:19 →
2026-09-15 00:22 UTC. Every gate 58/58; no failed cell; no violation on PostgreSQL; early
rejections on YugabyteDB all transient (yb-single 4, yb-cluster3 43); deferred confirmations 100%.

| Measure | pg-single | yb-single | yb-cluster3 |
|---|---:|---:|---:|
| Cell wall (topology written → result written) | 8.5 min | 26.4 min | 27.3 min |
| Initial load | 1.5 s | 30.6 s | 31.5 s |
| Reloads (5) | 7.5 s | 88 s | 158 s |
| Race wall, tiers 10 / 100 / 1k / 10k / 100k (s) | 4.6 / 3.0 / 4.1 / 31.3 / 60.4 | 30.8 / 21.6 / 40.5 / 122.2 / 61.0 | 31.6 / 24.2 / 39.0 / 122.1 / 60.9 |
| 100k race event | timed out, 18% sold | timed out, 2% sold | timed out, 3% sold |
| Tiers that used the full 3-min budget | none | none | none |
| Lifecycle wall, all tiers | 136 s | 322 s | 275 s |
| Remainder (verify, explain, reads, writes, audits) | ≈ 4.4 min | ≈ 14.5 min | ≈ 14.9 min |

**AM-02 rules applied:**

- Rule 1: yes on both YugabyteDB topologies (2% and 3% sold < 25%); `-race-tiers 10,100,1000,10000`
  added to `YB_HARNESS_FLAGS`.
- `M` = 8.5 × 13 + (26.4 − 1.0) × 14 × 1.15 + (27.3 − 1.0) × 14 × 1.15 = 1.8 h + 6.8 h + 7.1 h ≈
  **15.7 h** ≤ 24 h → rule 2: no other change. Rule 3 would have saved nothing (no tier used its budget).
- Repeated race (step 10, tiers 1 000 and 10 000, 3 trials), with verify taken as 2.5 min per cell:
  pg-single ≈ (0.1 + 2.5 + 3 × 0.62) × 13 ≈ 1.0 h; yb-cluster3 ≈ (0.5 + 2.5 + 3 × (0.53 + 2.69)) × 14 ×
  1.15 ≈ 3.4 h; total ≈ **4.4 h** ≤ 12 h → 3 trials.

## Main matrix (step 9) — facts

Run `20260915T002411Z`: commit `7dbdd11`, tag `run/03-reserved-seating/20260915T002411Z`, inputs digest
`a56ce92ce38b8204`, report `reports/20260915T002411Z.md`. 2026-09-15 00:24 → 17:30 UTC (17 h 6 min;
projection 15.7 h).

- **Topology times** (from `topology.yaml` and result-file times): pg-single ≈ 1 h 55 min (13 cells);
  yb-single ≈ 7 h 40 min (14 cells, ≈ 33 min/cell); yb-cluster3 ≈ 7 h 30 min (14 cells, including
  cluster start).
- **Gates:** 58/58 in every cell.
- **Negative controls:** 15 of 15 fired, per the report's controls table.
- **Failed cells (4), each with a DIAGNOSIS.md beside its logs:**
  - S4 on yb-single: lifecycle monitor timeout, SERIALIZABLE deadlock storm.
  - S4 on yb-cluster3: the reload before the lifecycle timed out, same storm.
  - E2 on yb-single and on yb-cluster3: lifecycle monitor timeout, node saturation.
  - All four completed verify, reads, writes and the race without a violation; ER-03 asks whether
    to re-run them.
- **PostgreSQL:** no violation in any correct design, and no early rejection.
- **YugabyteDB early rejections in correct designs (report TL;DR):**
  - yb-single: S3 2, E1 1, K1 2, L1 1, all transient; L2 4, unclassified (ER-04).
  - yb-cluster3: S1 41, S2 27, S3 46, E1 37, E2 7, K1 28, L1 35, L3 1, all transient; L2 13,
    unclassified (ER-04).
- **S1r, 0 early rejections on every topology.** Race retries after a short match: yb-cluster3 35,
  all 35 sold; yb-single 0; pg-single 0. Lifecycle retries (expired holds) never sold: pg 54,
  yb-single 4, yb-cluster3 2.
- **E1 sweeper-outage probe (the control condition):** pg-single 59/59, yb-single 8/8, yb-cluster3
  11/11 seats unavailable after expiry; 0 in every other design.
- **Throttling:** the yb-single database container was throttled in 63–96% of CPU periods per cell
  (S4 30%), from the per-cell `cpu.stat` snapshots; yb-cluster3 nodes less (S4 7–15%, E2 32–56%).

## Environment and coordination notes (continued)

- 2026-09-15: the owner stated that another AI agent (OpenAI Codex) works on the repository
  concurrently. AGENTS.md hard rule "assume a concurrent agent" was added (tag
  `repo/concurrent-agents-reconciliation`). Study 03 worked directly in the checkout of `main` until
  the matrix ended. Step 10 onwards runs in `.worktrees/study03-measurement` (branch
  `study-03/measurement`), which merges into `main` at step 11.

## Repeated race (step 10) — facts

Run `20260915T173255Z`: commit `52a9975`, tag `run/03-reserved-seating/20260915T173255Z`, inputs digest
`29296b1fe8dea2e3`, report `reports/20260915T173255Z.md`. 2026-09-15 17:33 → 21:06 UTC (3 h 33 min;
projection 4.4 h). `verify,race`, 3 fresh-load trials, 1 000- and 10 000-seat tiers, pg-single and
yb-cluster3.

- **Cells:** 27, 0 failed. Gates 58/58. Both S0 controls fired.
- **pg-single:** no violation and no early rejection in any correct design. Race spread across trials
  3–21% (report race table).
- **yb-cluster3 early rejections in correct designs:**
  - transient: S1 85, S2 89, S3 96, E1 95, E2 9, K1 108, L1 127, L3 10;
  - unclassified (ER-04): L2 41.
- **S1r on yb-cluster3:** 0 early rejections; 105 short confirmations retried, all 105 sold.
- The session was interrupted by a client logout at ≈ 20:30 UTC. The runner process and containers
  kept running and the run completed normally.
