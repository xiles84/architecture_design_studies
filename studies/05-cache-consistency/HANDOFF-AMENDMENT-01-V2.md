# Handoff amendment 01 — Study 05 v2 completion protocol

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-study05-v2-protocol` (HIGH planning), responding to
`task-20260922T025912Z-study05-independent-review`
(`book/evidence/reviews/20260922-study05-independent-review.json`).
**Required tag on integration:** `study-05/v4-handoff`.
**Reads with:** [`HANDOFF.md`](HANDOFF.md) (v1), which this amendment extends and does not replace.

## 0. What this amendment changes

The v1 survey established the cache result (2.2–3.5x, p99 an order of magnitude) and the
three-instance staleness mechanism (~87% wrong reads) on PostgreSQL + Redis at 800 donors, one
trial per measurement, closed loop, `resource_framing: db-only`. The registry records the unmeasured
regimes as `v2-gap-04`. This amendment defines:

1. real-TTL churn with the 300 s TTL actually crossed;
2. medium scale and a tested ceiling;
3. ≥3 repeated and randomized trials;
4. cache resource accounting (equal-total control);
5. YugabyteDB / cluster / verified placement;
6. open-loop demand;
7. the acceptance criteria the v2 report must meet before any number is published;
8. separately claimable execution tasks (section 7).

It runs nothing.

## 1. Real-TTL churn

The v1 run never crossed the 300 s TTL in a measured cell (only in unit tests). Run a phase whose
duration is a multiple of the TTL: 600 s, 1800 s and a burst at 5× the steady write rate for 120 s.
Workloads: the warm cacheable read set plus the write that invalidates it (write-through and
write-behind variants). Acceptance: each regime records how many refreshes and expiries fired; a
regime that crossed none is a gap, not a null result. Report the wrong-read rate separately from
throughput, as the v1 report did.

## 2. Medium scale and a ceiling

Scales: 8k, 80k and 800k donors (dataset generator extended, seed 42). At each scale run the same
controlled pair (`legacy` and `owned` cache models, memory-aside and redis-aside, strict and
relaxed freshness) so cardinality is isolated from model. Acceptance: report the working-set-to-RAM
ratio actually configured; if Redis evicts (evicted_keys > 0), the run must report it and the read
throughput is not a cache-hit result.

## 3. Repeated and randomized trials

Every cell ≥3 independent trials; report median, min–max and trial count. Where the database cannot
be reloaded cheaply, use a recorded randomized cell order and label the weaker design. A figure
whose spread exceeds the effect is reported as inconclusive, not as a ratio.

## 4. Cache resource accounting (equal-total control)

Alongside the v1 `db-only` framing, run a separately labelled **equal-total** arm: hold the total CPU
and memory constant and charge a share to the Redis container (e.g. DB 2 CPU / 2 GiB + cache
1 CPU / 1 GiB vs DB-only 2 CPU / 2 GiB). Report the two framings separately; never pool them. This
is the gap the v1 review named (`conf-04`).

## 5. YugabyteDB, cluster and verified placement

Repeat the core strict/relaxed pair on YugabyteDB 1-node (RF=1) and 3-node (RF=3), and on the
colocated vs non-colocated pair that v1 never executed. **Verify actual placement**: record the
tablet/leader distribution for the cached and uncached tables in the run log. A run without a
placement record is a gap. Report per engine; no cross-engine ranking.

## 6. Open-loop demand

Add an open-loop arrival phase (fixed arrival rate, no closed-loop backpressure) at two rates, with
the platform's arrival driver. Report achieved throughput, overflow/generated-but-not-attempted
operations and latency percentiles. Closed-loop results remain floors; open-loop results must state
the offered rate and the server-side saturation point if reached.

## 7. Separately claimable execution tasks

1. **Study 05 v2 churn and scale** — sections 1, 2 and the repeats of section 3.
2. **Study 05 v2 resource accounting and engines** — sections 4 and 5.
3. **Study 05 v2 open-loop demand** — section 6.

Each takes the benchmark lock, runs one matrix at a time, and produces a generated report plus a
signed analysis under `studies/05-cache-consistency/`. None may publish a wrong-read rate without
the same run's throughput and the Redis eviction/placement record.

## 8. Acceptance criteria for the v2 report

- TTL/refresh actually crossed in every churn cell, with counts.
- Working-set-to-RAM ratio and eviction reported at every scale.
- ≥3 trials per cell, or a recorded randomized order with the limitation stated.
- `db-only` and equal-total arms labelled and separate.
- Placement recorded; YugabyteDB and PostgreSQL reported separately.
- Wrong-read rate and throughput from the same cell; failed cells named, never winners.
- Every number resolved to a v2 evidence claim before it enters the book.
