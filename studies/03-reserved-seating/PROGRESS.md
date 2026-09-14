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
| 5 | SQL catalogue, 13 designs | done | `13941ff` |
| 6 | Harness (17 unit tests pass in the build) | done | `eb7e4ff` |
| 7 | Diagrams | in progress | |
| 8 | Dev checks and calibration | in progress | |
| 9 | Main matrix (`small`) | pending | |
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

## Observations for the analysis phase

Facts recorded during dev checks and runs (handoff §13).

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
