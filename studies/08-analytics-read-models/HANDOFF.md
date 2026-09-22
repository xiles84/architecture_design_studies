# Execution Handoff — Study 08, analytical, search and read-model copies

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-analytics-read-model-protocol` (HIGH planning).
**Required tag on integration:** `repo/analytics-read-model-handoff-v1`.
**Numbering:** Studies 06 (native models), 07 (topology) exist; this is Study 08.

## 1. The question

For top-N and dashboard workloads, what does each read-model strategy cost and return: base OLTP
queries, rollups/materialized views, a search copy, and an analytical copy? The corpus knows the
read gain and the maintenance cost of a rollup in isolation; it does not know how the four compare
under one workload, one correctness contract and one resource budget.

## 2. Read-model strategies

| Strategy | Implementation | Why |
|---|---|---|
| Base OLTP | PostgreSQL 17.11, the Study 04 configuration-portal schema | The floor: no extra copy, pay per query. |
| Rollup / materialized view | PostgreSQL materialized view over the same base, refreshed on a stated schedule | The family already measured in isolation; here it competes with other copies. Shared refresh shapes are counted too. |
| Search copy | PostgreSQL full-text/trigram index over the same base, kept in sync by trigger or application | A search copy inside one engine isolates the *copy*, not the engine. A native search engine (OpenSearch line) is a separately labelled extension, deferred to the native study to avoid changing two factors. |
| Analytical copy | ClickHouse (native columnar) as a separate read copy loaded from the base | A genuine analytical engine, so the "analytical copy" is not a renamed index. |

Pinning: PostgreSQL 17.11 as today; ClickHouse and (if used) OpenSearch are pinned to an exact line
in this protocol and to an exact tag plus image digest in `versions.env` by the harness task. A moved
tag blocks the run.

## 3. Workloads

Three **top-N/dashboard** workloads, each with a deterministic dataset (seed 42) and a recorded
offered rate:

1. top-10 products by configuration-change count in a window (dashboard),
2. top-N donors by lifetime total from the charity dataset (donor portal),
3. recent-activity feed for one product (search copy's home turf).

Each workload runs against every strategy, with the same logical answer. The base OLTP query is the
reference; a strategy is only a result if it returns the same answer under the same correctness gate.

## 4. What is measured for every strategy

- **Query gain**: throughput and p50/p99 latency for the top-N read, open-loop and closed-loop.
- **Refresh / maintenance**: refresh interval actually crossed, refresh duration, and write
  throughput under live refresh.
- **Staleness**: the age of the copy at read time, and the wrong-answer window when the base changes.
- **Storage**: base, copy and index bytes separately.
- **Correctness**: the copy must return the base's answer at the stated freshness contract; a
  deliberately stale control must be seen to differ.

## 5. Correctness and negative controls

Study 04's invariants apply to the base. Each copy adds: a refresh-boundary test (read during a
refresh must be defined: blocked, old snapshot, or new snapshot, and recorded), a staleness control
that must be seen to return a stale answer, and a duplicate/lost-copy control for the loader. A copy
that cannot demonstrate its negative control is not evidence.

## 6. Resources and topology controls

Per-node 2 CPU / 3 GiB for the base engine and 2 CPU / 2 GiB client as today; each copy gets a
separately recorded budget and a separately labelled equal-total arm (total resource held constant).
v0 is single-node PostgreSQL; a ClickHouse cluster is out of scope (Study 07 owns topology). Client
readers/writers and connection pools are calculated and recorded.

## 7. Deliverables and separately claimable execution tasks

1. **Analytics harness and pins** — ClickHouse Containerfile, `versions.env` pins with digests,
   deterministic loader, correctness gate and refresh/staleness instrumentation.
2. **Rollup and materialized-view arms** — refresh schedules and the top-N workloads.
3. **Search-copy arm** — PostgreSQL full-text/trigram search copy plus the trigger/application sync
   variants.
4. **Analytical-copy arm** — ClickHouse copy, loader cadence, and the top-N workloads.

Each task takes the benchmark lock, runs one matrix at a time, and produces a generated report plus a
signed analysis under `studies/08-analytics-read-models/`. No cross-strategy pooled ranking across
different runs.

## 8. Acceptance criteria for the study report

- Every strategy returns the base's logical answer under the stated freshness contract.
- Refresh interval, refresh duration, staleness and storage reported per strategy.
- Negative controls fire (stale copy, duplicate/lost loader row, refresh-boundary behaviour).
- Per-node and equal-total arms labelled and separate; single-node only.
- Search copy inside one engine and any native search engine are labelled separately.
- Every number resolved to a v2 evidence claim before it enters the book.
