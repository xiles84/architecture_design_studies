---
analysis_id: 20260921-cache-consistency-allgreen--deepseek-flash--2026-09-22
run_id: 20260921T-survey3
study: 05-cache-consistency
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability independent review, Deep Code CLI; effort setting not exposed. Not the author of the study."
analyzed_at: 2026-09-22
inputs_digest: 03fac739d0835768
repo_commit: 48fe1e33db5467dfa4421f354df4a6ee77becaff
supersedes:
status: current
headline: Independent HIGH review of Study 05 — the throughput and strict-cost claims re-derive cleanly, but the all-green analysis swaps the legacy and owned three-instance wrong-read rows (87.81% belongs to the owned model), and the Redis cache/authoritative-store boundary must be stated in the book.
---

# Analysis — `20260921T-survey3` — deepseek-flash, independent HIGH review

> This analysis is the review required by `task-20260922T025912Z-study05-independent-review`. It
> reads the run, the report `reports/20260921T-survey3.md`, the signed analysis
> `20260921-cache-consistency-allgreen`, the partially-superseded discussion
> `20260921-cache-fences-and-ledger`, the superseded analysis chain under `outdated/`, and the
> registry entries for Study 05. It edits none of them. It is the independent review the study
> itself says it lacks; no measurement was run and no database was started.

## TL;DR

- **The cache-throughput and strict-freshness claims re-derive exactly.** Every ratio in the
  all-green analysis's first two tables reproduces from `results/20260921T-survey3/pg-single/*.json`:
  legacy no-cache 7 794 → memory-aside-strict 20 233 (2.60x) and redis-aside-strict 18 667 (2.39x);
  owned no-cache 5 805 → redis-aside-strict 20 380 (3.51x); owned-pess no-cache 7 043 →
  15 782 (2.24x). `claim-11` and `claim-12` are **accepted** for their stated scope, and
  `claim-gap-02` is **accepted**.
- **The study's strongest result is mis-attributed in the analysis.** The all-green analysis's
  wrong-read table gives the legacy three-instance row as 38 749 / 87.81 % / 44 122 and the owned
  row as 35 287 / 87.24 % / 40 442. The run JSON is the other way round: `l:na:red:wth:rel`
  instances_3 is **35 287 / 40 449 = 87.24 %** and `o:opt:red:wth:rel` instances_3 is
  **38 749 / 44 126 = 87.81 %**. The report itself (`o:opt:red:wth:rel`) has it right; the derived
  prose swapped the two models. The book must not attribute 87.81 % to the legacy model.
- **Redis is measured as a shared cache, never as a source of truth.** Every `redis-*` cell is a
  cache backend for the same portal data in PostgreSQL, and the study never measures Redis
  durability, persistence, failover or eviction as a system of record. The book must keep that
  boundary explicit, or a reader will hear "Redis vs PostgreSQL" where the run says "shared
  external cache vs in-process cache".
- **Every limit the brief names is real and is preserved:** single trial, PostgreSQL only, `small`
  (800 donors), closed loop, `db-only` framing (no equal-total), no real 300 s TTL crossing, no
  YugabyteDB/cluster, the placement pair never executed, no open-loop demand. This review is the
  first independent reading; before it there was none.
- **The gaps go to `task-20260922T025912Z-study05-v2-protocol`,** which already exists and depends
  on this review. This review adds four concrete requirements to it (below) rather than creating a
  new task.

## What I validated

Run `20260921T-survey3`, digest `03fac739d0835768`, commit `48fe1e3`, 27 cells, 0 failed, both
controls fired, PostgreSQL 17.11 + Redis 7.4.11 on `host-zenbook-ux5406sa`, one trial per
measurement. Method: every figure below and every verdict figure was re-extracted from
`results/20260921T-survey3/pg-single/<scenario>.json` (`warm_reads`, `wrong_reads`, `instances`,
`ack_check`, `audits`), then compared with the report, the analysis and the discussion.

## Conclusions

**1. `claim-11` (cache throughput) is accepted, with its scope.** The cacheable-read gain is
2.2-3.5x and the p99 improvement is an order of magnitude, all re-derived. It applies to *warm,
cacheable* reads on one `small` PostgreSQL node in a `db-only` resource framing; it is not a
statement about uncacheable reads, writes, larger datasets, or a cache that must pay for its own
container. The analysis says all of this; the claim must carry it verbatim.

**2. `claim-12` (strict freshness cost) is accepted as a negative read result, not a licence.**
Every controlled pair's read difference is inside the run's ~20 % noise floor (re-derived spread
-8.3 % to +12.3 %), and the analysis is careful to say the cost is on the write path (more fences).
The correct book statement is "strict freshness cost nothing measurable in *read* throughput at
`small`, one trial"; it must not become "strict freshness is free".

**3. The three-instance finding is real and important; its model attribution must be corrected.**
Independent of the swap, both three-instance `through`+relaxed cells are ~87 % wrong while every
single-instance relaxed cell is 0 % wrong — a 87-vs-0 separation far outside any noise, and the
right shape for the mechanism (a private in-memory lease cannot coordinate instances sharing a
store, and the `through` strategy republishes what a peer then reads). The correction the book
needs: **87.81 % (38 749/44 126) is the `owned-opt-redis-through-relaxed` cell; 87.24 %
(35 287/40 449) is `legacy-na-redis-through-relaxed`.** The analysis's "87.8 % and 86 % in two
runs" also needs its second figure re-checked against its own run before publication — this review
re-derived only the `20260921T-survey3` pair.

**4. The Redis boundary is a hard scope requirement.** The run's Redis cells measure an external
shared cache; `backend: redis` is a cache backend, and the "redis-aside vs redis-through"
comparison is about where invalidation is visible, not about Redis's own guarantees. Redis is
never authoritative here, and no cell measures it as one. The book must say so wherever it cites
these numbers; otherwise the 2.2-3.5x is read as a database comparison it is not.

**5. The rollup/rolldown/embedding coverage is directionally usable, with one caveat kept.** The
analysis uses `ref-rollup-trigger` (9 412 reads/s, clean `a_rollup_drift` this run),
`ref-embedded-locked` (7 007 reads/s, worst p99) and the legacy flattened model as the rolldown
baseline. The review accepts these as within-run reference points and keeps the analysis's own
condition: the trigger rollup's earlier drift failure is why the clean audit is credible, and the
design should be re-tested at `medium` before it is relied on.

**6. The superseded chain is used correctly.** The discussion `20260921-cache-fences-and-ledger`
is marked partially superseded and its §1-§3 mechanisms (fence, acknowledgement boundary,
pending-state race) are cited as mechanism, while its residual-failure sections are explicitly
retired. That is the right treatment; nothing in this review depends on a retracted number.

**7. The four v2 requirements this review hands to the protocol.** No new task:
`task-20260922T025912Z-study05-v2-protocol` already owns the gaps. It must add, explicitly:
(a) **repeated trials and ideally a second machine** for the three-instance result, since it is
single-trial on one laptop and is the study's most-quoted number; (b) **real-TTL churn** across
the 300 s boundary, which is verified only by unit tests today; (c) **equal-total framing** so the
Redis container's CPU/memory is charged to the cached cell, not left outside the comparison; and
(d) **execution of the placement pair** (`ref-y1-colocated`/`ref-y2-noncolocated`) with actual
placement verified, or an explicit statement that colocation is untested. Items (a)-(d) are
protocol scope, not measurements here.

## Where I think the measurement is weak

- **One trial per measurement, one machine, closed loop.** The ~20 % noise floor is larger than
  every strict-vs-relaxed difference; only the 87-vs-0 three-instance separation is safely outside
  it. Deep percentiles are floors.
- **The corrected attribution shows the derived prose is not yet trustworthy without a check
  against JSON.** A reader who quoted the analysis's table would have mis-reported which model
  loses 87.8 %.
- **No sustained-churn phase**, so the TTL rule's field behaviour is unmeasured; no equal-total
  resource accounting, so the cache is subsidised; no YugabyteDB, cluster or executed placement
  pair.
- **Cache gauges vs counters.** `items`/`resident_bytes` are end-of-cell snapshots; the analysis
  already says to read counters as evidence. The review agrees.
- **No independent review before this one.** One agent designed, fixed and analysed the study; this
  review is the first outside reading and it found a material attribution error, which is itself
  evidence for why the gap mattered.

## What I would measure next

1. **Repeat the three-instance `through`+relaxed cells** (legacy and owned) with ≥ 3 trials and, if
   possible, a second machine; report the wrong-read rate and the superseded-entry count each time.
2. **Churn across the real 300 s TTL** with the early-expiry formula active, measuring stale reads,
   refills and load.
3. **Equal-total**: rerun the headline cache pairs with the Redis container's budget charged, so
   the 2.2-3.5x is a fair comparison.
4. **A larger scale and an executed placement pair**, which are the two coverage items most likely
   to change the ranking.

These belong to `task-20260922T025912Z-study05-v2-protocol`.

## Disagreements with other analyses

I agree with the conclusions of `20260921-cache-consistency-allgreen` and with the mechanism record
in the partially-superseded discussion. I **disagree with one table in the all-green analysis**:
its legacy and owned three-instance wrong-read rows are swapped (conclusion 3 above); the report's
own `o:opt:red:wth:rel` line and the run JSON are the arbiter. I also require the Redis
cache/authoritative-store boundary and the second-run "86 %" figure to be stated explicitly before
either is cited in the book. No other figure I checked needed a change.
