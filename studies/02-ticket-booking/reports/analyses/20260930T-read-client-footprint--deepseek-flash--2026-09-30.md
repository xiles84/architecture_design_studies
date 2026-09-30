---
analysis_id: 20260930T-read-client-footprint--deepseek-flash--2026-09-30
run_id: 20260930T095448Z (8 workers), 20260930T095859Z (4), 20260930T100257Z (2), plus the pre-commit wiring check 20260930T101054Z and the verification run 20260930T101521Z at the new default
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; the read-client footprint calibration that implements the owner's 2026-09-30 decision."
analyzed_at: 2026-09-30
inputs_digest: 9657526bfd1d78f7 (8 workers), 3b9612501e134ef3 (4), b868b00d489af428 (2), 2341a0d87ebe5414 (default verification)
inputs: 20260930T095448Z@9657526bfd1d78f7, 20260930T095859Z@3b9612501e134ef3, 20260930T100257Z@b868b00d489af428, 20260930T101521Z@2341a0d87ebe5414
repo_commit: 386a610888c8f1f98129c484113d8a96fcf46135
status: current
headline: The read phase's 8 workers were the client's bottleneck - at 2 workers the client is unthrottled (0 of 631 periods against 422 of 633), the database's throttling falls from 900 to 109 periods, and every read cell becomes latency-bound at 1.9 operations in flight for 2 workers; the study's read phase now defaults to 2.
---

# Analysis — sizing the read client down (study 02) — deepseek-flash

## TL;DR

- **The lever is the worker count, and the curve is clean.** Client throttling: **422/633 periods (66.7%) at 8
  workers → 33/632 (5.2%) at 4 → 0/631 at 2**; the database's own throttling falls with it, **900 → 560 →
  109** periods, because the client stops competing for the same four CPUs.
- **At 2 workers every cell is latency-bound, which is what a read benchmark should be.** The product
  `ops_per_sec × p50` — operations in flight — is **1.9 for both designs at every tier**, against 2 workers,
  which is Little's law holding to within 5%. At 8 workers the same product is 3.2–5.3 *for 8 workers*: the
  client was so starved it could not keep its own connections busy.
- **The ordering is the measurement that travels**: `r1` stays ahead of `p3` on `q01` at every footprint
  (53.2k vs 15.7k; 27.1k vs 13.9k; 20.2k vs 11.8k at 1,000 seats), and `p3`'s tier-100,000 cell is
  **621–631 ops/s at all three footprints** — a genuinely store-bound cell that the client's footprint does
  not move.
- **What this changes**: `run-study.sh` now defaults the **read phase** to 2 workers
  (`READ_CONNS`/`--read-conns`), with a comment naming this evidence; the harness grew a `-read-conns` flag so
  the read phase can be sized without touching the write, race and hold phases, which keep the published
  `-conns 8` and stay comparable with every published result. The manifest records `client_read_conns`.
- **Verified**: a first attempt ran before the change was committed, so the runner refused to tag it
  (`repo_dirty: true`, kept as the wiring check and not cited); the committed verification run
  `20260930T101521Z` records `client_read_conns: 2`, **0 throttled periods** in both designs' read phases, and
  the expected signatures - `r1` flat at 20.2/17.9/18.0/19.2/20.2k across the five tiers, `p3` falling
  15,061 to 606 ops/s - so the default path is wired end to end, not just measured.

## What I was looking at

Three whole runs on `pg-single`, scale `small`, phases `verify,explain,read`, the published budgets
(`DB_CPUS=2`, `CLIENT_CPUS=2`) and duration, designs `r1_inventory_row` and `p3_precreated_cas`, differing
only in `-conns`: 8 (the published footprint), 4 and 2. Both trees clean, all three tagged at `386a610`, no
failed cells. Client cgroup counters are per phase; the database's are the runner's before/after capture.

| `-conns` | client read phase | database, same cell | `r1` q01@1000 | `p3` q01@1000 | ops in flight |
|---|---|---|---|---|---|
| 8 | 111.6 s CPU, **422/633 throttled** | 900 throttled periods | 53,244 | 15,732 | 5.3 / 3.2 |
| 4 | 103.1 s, 33/632 | 560 | 27,109 | 13,875 | 3.8 / 3.1 |
| 2 | 63.4 s, **0/631** | **109** | 20,172 | 11,792 | **1.9 / 1.9** |

## Conclusions

**The published footprint was measuring the client, and the owner's direction fixes it.** At 8 workers the
client spends 111.6 s of CPU inside a 63.3 s window and loses two thirds of its periods; its throughput is
then set by how fast a throttled two-core container can issue queries, not by the store. Dropping to two
workers costs two thirds of the headline throughput (53.2k → 20.2k on `r1`) and buys a cell whose operations
in flight equal the workers — the definition of a latency-bound measurement — with the client at half its
quota and the database's own throttling down by a factor of eight.

**4 workers is the aggressive option, 2 the safe one.** At 4 the client is throttled 5.2% of periods and uses
103 s of CPU against a 126 s allowance: only 18% headroom, which a slower session could consume — the same
host moved a cell's level by 40% between two sessions this week. At 2 the headroom is 50%. The default is
therefore **2**, and `--read-conns 4` remains available for a run that wants throughput and can show the
client unthrottled in its own result.

**What this does not fix, and what it does.** It does not make the read family's levels portable — they are
now honest levels at a stated footprint, and the corpus's rule that levels move between sessions still
applies. It does make them *comparable*, which they were not while the client was the constraint, and it
leaves the ratios untouched: the design orderings (`r1` ≫ `p3` on `q01`; `p3`'s store-bound 100,000-seat
cell; the `q05` and `r03` results) are the same at every footprint. That is the property the read tables
should be read for.

**Prediction verdicts** (recorded in `READ-CLIENT-FOOTPRINT.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| F1 | throttling falls with the footprint; 2 workers is 0 or near 0 | **confirmed** — 66.7% → 5.2% → 0% |
| F2 | at the unthrottled footprints `ops/s × p50` ≈ the worker count | **confirmed** — 1.9 against 2 workers for both designs at every tier; 3.8/3.1 at four |
| F3 | the `r1` ≫ `p3` ordering survives every footprint | **confirmed** — and `p3`'s 100,000-seat cell is 621–631 at all three |
| F4 | the database's throttling falls too | **confirmed** — 900 → 560 → 109 periods |
| F5 | correctness gates stay green | **confirmed** — no failed cells at any footprint |

## Where I think the measurement is weak

1. **Two designs of fifteen.** `r1` and `p3` are the extremes of read cost; the other thirteen, including the
   counter family that was throttled hardest in the 2026-09-20 anchor, are unmeasured at the new footprint,
   so the family's re-measurement is still owed.
2. **One run per footprint.** The curve is monotone and the mechanism (workers → demand → throttle) is
   visible in three independent quantities, but a single run per point cannot separate footprint from
   session; 4 sits close enough to the edge that a repeat could tip it.
3. **`ops/s × p50 ≈ workers` assumes p50 stands for the mean.** It is a diagnostic, not a proof; a proper
   latency-bound check would use the mean and the observed concurrency (Little's law with measured L).
4. **One topology, one engine, one scale, one host** — and the machine's 4 CPUs / 2 GiB are now a stated
   project constraint rather than an open question.
5. **The read phase's absolute levels are now footprint-specific.** Any future reader comparing a new read
   run with an old one must check `client_read_conns` in the manifest; the reports print the option, and the
   default change is recorded here and in `CONTEXT.md`.

## What I would measure next

1. **Re-measure the read family at the new footprint** — the four designs of the 2026-09-29 re-run first,
   then the remaining eleven — so the family has one footing and its tables can be read again.
2. **The counter family (`p4`, `c4`, `r2`)** specifically, since those were the designs whose read phases
   were most throttled in the primary anchor.
3. **A repeat at 4 workers** to decide whether the aggressive option is safe enough to offer as a documented
   alternative or should be dropped.

## Disagreements with other analyses

None. This implements the direction set by `20260930T-read-client-sizing` (the client's quota is a symptom of
the total allocation) and by `20260930T-read-q01-plan-check` (the flat cell is O(1) by design), and it leaves
both intact. It supersedes nothing.
