---
analysis_id: 20261004T-select-join-cte-v2--deepseek-flash--2026-10-04
run_id: select-join-domain2-pg1/pg2/pg3/yb1/yb2/yb3-20261004; select-join-topo-yb-{single,equal,cluster}-{i1,i2,i3,bo1,bo2,bo3}-20261004
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek Flash via opencode CLI (LOW executor session)"
analyzed_at: 2026-10-04
inputs_digest: "pg1=347d12b60a3bffc23ae0193584db62e42e0b967778dd90d8c6ccd293ddf7ac1a; pg2=d4a7ba24e4ff2986f6b0640e3a4e72ab1b07969591ea625b3228c29d0aa6edbe; pg3=ccd615919f584e75b6fd4640f737bcf8fc0b7c3a1bf370b6d2d483399a86c704; yb1=f03b1ad718dce4985f6849dd8f4cf33687d09daea675106ed2992eb3a1aecb2b; yb2=0b8ec8dc842ac10e68c5fcb0659308b343a97c2bf0cd279c1f18cc70e744ad04; yb3=e497ebb27c95425f238f8a12e2529f92a0346468baaa0a32449f38cef0c479cc"
repo_commit: e2bc34524b5a46ca103fd7a299d1d15991e4ce0b
supersedes: 20261004T-select-join-cte--deepseek-flash--2026-10-04
status: current
headline: "On the real person/donation/charity domain, forcing MATERIALIZED on a join-containing CTE costs 67–88x on PostgreSQL at a selective threshold and at most ~1.7x on YugabyteDB; the reuse crossing exists only on PostgreSQL at the full-consumer endpoint and is unresolved at one reuse; a 3-node cluster is ~1.33x slower only at the 90% scan."
---

# Analysis (successor) — join/filter spellings and CTE materialization on the Study 01 domain — deepseek-flash

This supersedes my own `20261004T-select-join-cte--deepseek-flash--2026-10-04`. That
analysis read a *synthetic* microfixture (customers/facts) and a mis-derived reuse
pair; a HIGH review (`docs/select-query-templates/reviews/20261004-join-high/`)
rejected its P4 interpretation. This version uses the real
person/donation/charity fixture loaded through binding views, so the measured
catalogue SQL is byte-identical between arms, and it reports the reuse pair the
contract actually defines.

## TL;DR

- **P1 (default spellings):** all six default spellings are answer-equivalent;
  only PostgreSQL shows timing spread wide enough at selective thresholds to
  read two spellings as different, and that spread is not reproducible across
  runs, so treat those cells as unresolved, not as a ranking.
- **P3/P4a (forced `MATERIALIZED` on the joined-output CTE):** real and
  engine-specific. PostgreSQL pays **67–88x** at 20 rows, 26–31x at 2000, and
  ~1.0–1.4x at 100k/180k. YugabyteDB pays **1.43–1.66x** at its worst.
- **P4b (the reuse pair):** the pair is `NOT MATERIALIZED / MATERIALIZED`. On
  PostgreSQL with a full consumer it is 0.89–1.05 at one reuse (**unresolved**),
  1.61–1.73 at two, 2.88–2.98 at eight, 4.39–4.56 at sixteen; with a 1% consumer
  it is 0.18–0.94 (notmat faster). On YugabyteDB it is 1.01–10.3 (notmat slower)
  at every reuse. The crossing exists only on PostgreSQL, and **not** at one
  reuse.
- **P12 (50–90% convergence):** converged within ~10% on YugabyteDB; only
  partial on PostgreSQL (a filtered CTE reached 1.58x at 180k).
- **Topology (new):** a balanced 3-node RF3 cluster is equal to single-node at
  1% selectivity and ~1.33x slower at the 90% scan; the equal-total single node
  (6 CPU/9 GiB) is not faster than the per-node single (2 CPU/3 GiB) on this
  workload. The blended phase sustained 10.02 committed updates/s with 0 ledger
  mismatches.

## What I was looking at

- **Frozen domain matrix, 3 runs/engine, one client build:** 91 cells (10 bag
  spellings x 6 thresholds + 3 policies x 5 reuse counts x 2 consumer
  selections + a randomized position control), 5 trials, seeded interleaving,
  `EXPLAIN ANALYZE` on every cell, both semantic negative controls firing in all
  six runs. 91/91 cells valid in each. PostgreSQL 17.11 and YugabyteDB
  2025.2.6.0-b111, one laptop, containers CPU-capped (client 1 CPU/2 GiB).
- **Topology subset, 18 runs:** domain arm, 30 cells, 5 trials, three arms —
  single (2 CPU/3 GiB), equal-total single (6 CPU/9 GiB), 3-node RF3 (2 CPU/3 GiB
  per node) — each isolated and blended (one fixed writer at 10 updates/s), 3
  whole runs per arm. All 30/30 valid.
- I ignored the byte-level provenance of the run's own report beyond its digest;
  the numbers below are medians of per-run medians, which understates
  between-run variance. The synthetic arm is retained only as a fidelity check
  and is not interpreted here.

## Conclusions

**P1 — default spellings.** Answers were identical to the oracle in every cell,
so no default spelling is semantically wrong. Timings: on YugabyteDB every
default spelling sits within 0.90–1.08 of the join+WHERE baseline at all six
thresholds. On PostgreSQL two cells spread widely across runs at the selective
end (`comma_join` 0.97–1.68 at 20 rows; `join_on` 0.64–1.53 at 20 rows), and the
spread direction flips between runs. With a per-run position effect of
0.69–1.11 (the position-control cell re-timed at the end of the run), a single
run cannot separate those spellings. **Action:** do not claim a default-spelling
ranking on PostgreSQL at selective thresholds; use at least three interleaved
runs and the across-run median, as done here.

**P3/P4a — forcing the materialization of a join-containing CTE.** This is the
strongest and most reproducible result. Against the joined-output CTE's own
default plan, PostgreSQL's forced `MATERIALIZED` is 66.9/87.5/81.4x at 20 rows,
41.7/62.7/49.9x at 200, 26.1/31.4/29.1x at 2000, and 1.0–1.4x at 100k/180k
(three runs, in that order). The non-materialized hint is 0.76–1.23 everywhere.
On YugabyteDB the same forced materialization is 1.43–1.66x at its worst and
below 1.0 at the largest threshold. The mechanism is visible in the plans: on
PostgreSQL the forced materialization turns a nested-loop/index path into a
materialize-plus-scan that also blocks the selective pushdown; on YugabyteDB the
same hint is absorbed by a hash join over a distributed scan, so the penalty is
mild. **Action:** the engine, not the SQL spelling, decides whether the hint is
dangerous; the same query text must be measured per engine.

**P4b — the reuse pair.** The contract's pair is NOT MATERIALIZED against
MATERIALIZED, holding the reuse count and consumer constant. On PostgreSQL with
the full consumer the medians are 0.904 (n=1, runs 0.89–1.05), 1.69 (n=2,
1.61–1.73), 2.36 (n=4, 2.31–2.44), 2.95 (n=8, 2.88–2.98), 4.43 (n=16,
4.39–4.56). Only the n>=2 endpoints separate; **n=1 is unresolved** and must not
be used to place a crossing. With the 1% consumer the pair inverts: 0.20–0.84 at
n=1..16 (notmat faster). On YugabyteDB the pair is 1.04/1.01 (n=1), 1.88–2.22
(n=2), 2.60–2.88 (n=4), 5.07–5.55 (n=8), 9.26–10.31 (n=16) — notmat is slower
throughout, roughly linear in the reuse count, with no crossing. **Action:** the
crossing is a PostgreSQL-only, consumer-dependent effect; the honest statement is
"on PostgreSQL, NOT MATERIALIZED wins at a highly selective consumer and loses
once the consumer is unselective and the CTE is reused more than once."

**P5 — semantic negative controls.** Both deliberately-wrong controls
(`x01`, `x02`) were detected in all six runs. This says nothing about
performance; it is the evidence that the gate can fail.

**P12 — convergence at 50–90%.** At 50% and 90% selectivity the default
spellings converge on YugabyteDB (0.90–1.04). On PostgreSQL they converge only
partially at 90% (a filtered CTE reached 1.58x and comma-join 0.75–1.51), so the
convergence claim holds on YugabyteDB and is unresolved on PostgreSQL.

**Topology (new, subset).** Across three runs each: at 1% selectivity
join+WHERE medians are 58.9 ms (single), 55.4 ms (equal-total single), 59.7 ms
(cluster) — indistinguishable, i.e. cluster round-trip and per-node throttling
do not show at a point read. At 90% they are 119.9 / 119.3 / 159.1 ms — the
cluster is ~1.33x slower, consistent with shipping the wider result set between
nodes. The equal-total single node is *not* faster than the per-node single
node, so the extra 3x budget does not remove a client- or plan-bound limit here.
The forced-materialization ratio at the 2000-row threshold is 1.54–1.66
(isolated) and 1.50–2.11 (blended, highest on the cluster): the concurrency
penalty lands on the materializing plan too, but the run-to-run spread is
comparable to the effect, so I do not rank the arms on it.

**Blended phase.** The writer offered and completed 10.02 updates/s in every
blended arm (0 failed), reconciled its Go ledger against the database counters
with 0 mismatches, and 154–173 reader operations per arm ran while a write
committed. Reader retries of transient conflicts were few (median 0–3 per run):
YugabyteDB restarted the longest unmaterialized multi-reference reads (SQLSTATE
40001) under the live writer. Those are recorded on the cell, not hidden, and
never timed without a verification pass. This is the first evidence in the study
that the *concurrency* cost of the reuse choice is engine-specific, but it is
three runs and must not be read as a rate.

## Where the measurement is weak

- Three runs of five trials per engine/arm; the across-run median hides
  between-run variance, and the position control moved by up to x1.11/x0.69
  within a run, so effects below ~20% are not resolved.
- One laptop; database and client share eight cores; the containers are
  CPU-capped even when alone. Cluster nodes are containers on the same host, so
  latency is not a network measurement and colocation is not evidenced.
- The writer is bound to one SQL endpoint; cluster readers are round-robin, but
  writes are not, so the cluster write path is under-measured.
- The colocated arm is an explicit gap: the pinned engine rejects
  `colocation = true` on the default database (SQLSTATE 0A000). Only the
  non-colocated baseline was measured.
- PostgreSQL cells exceeded 20% sample spread in 48–63 of 91 cells, so its
  timings rest on the median, not on tight samples.

## Coverage and gaps

- Measured pairs: default-vs-join+WHERE (six spellings x six thresholds x two
  engines), forced MATERIALIZED / NOT MATERIALIZED vs default (two CTE variants),
  the reuse pair (five reuse counts x two consumers x two engines), single vs
  equal-total single vs 3-node cluster (isolated and blended).
- Non-applicable: optimistic/pessimistic concurrency comparison (this task has
  one fixed maintenance writer, no contended business invariant);
  rollup/rolldown/embedding (no parent aggregate in this study's data model).
- Outstanding: colocated placement, cluster-balanced writes, and more than three
  runs per arm.
