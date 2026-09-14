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
| 8a | AM-01: transient-refusal class, refusal diagnostics, S1r, report; dev checks dc12–dc14 | implementation done (19 unit tests pass in the build); dc12 running | `e10dc6c` |
| 9 | Main matrix (`small`) | pending: after 8a | |
| 10 | Repeated race trials | pending | |
| 11 | Context, lessons, README; ready for analysis | pending | |

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
