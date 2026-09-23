# Study 05 — external cache throughput and consistency — run `20260923T1137Z-res-equaltotal-pg`

Generated 2026-09-23T11:25:43Z from 1 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `20260923T1137Z-res-equaltotal-pg` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `a9f59f010d027cc5dd8de2852f139d2376c72dba` |
| `git describe` | `run/05-cache-consistency/20260923T1135Z-res-dbonly-pg` |
| Working tree dirty at start | false |
| Run tag | `run/05-cache-consistency/20260923T1137Z-res-equaltotal-pg` |
| Benchmark image | `localhost/cachebench:1` (`86683f6286f160b834c`) |
| Resource framing | `equal-total` |
| **Inputs digest** | `b56bc32cd4145fb0` |
| Cells | 1 |
| Distinct scenarios | 1 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 3 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 1; cells that failed: 0.
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No strict cell recorded a stale-after-ack read.
- No cell reported an impossible / uncommitted cache value.
- Highest and lowest measured throughput across all cells: 7.3k (`pg-single/o:opt:red:wth:str/warm_read_only`) and 6.7k (`pg-single/o:opt:red:wth:str/cacheable_90_reads`).

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `b56bc32cd4145fb0`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | pass | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 7.3k | 0.46 | 9.13 | 686 | 48630 | 0 | 11.7 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 7.1k | 0.47 | 7.25 | 641 | 48759 | 0 | 2.3 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 6.7k | 0.37 | 4.26 | 698 | 51946 | 0 | 21.9 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 8.0k | 0.17 | 2.89 | 688 | 58177 | 0 | 15.1 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 3.0k | 0.19 | 6.33 | 667 | 22660 | 0 | 21.3 | 0 | 0.00% |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 154 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 52154 | 52154 | 0 | 0.00% | 51360 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 52737 | 52505 | 0 | 0.00% | 52427 | 0.00% | 0 | 0 | — | — | 0 | 1277 | 323 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 106082 | 105850 | 0 | 0.00% | 105772 | 0.00% | 0 | 0 | — | — | 0 | 1916 | 323 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 124005 | 123349 | 0 | 0.00% | 123095 | 0.00% | 0 | 0 | — | — | 0 | 2441 | 926 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 179709 | 179053 | 0 | 0.00% | 178799 | 0.00% | 0 | 0 | — | — | 0 | 8075 | 926 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | 0-75s | 230159 | 4939 | 2.1% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | 6645 | 209 | 190 | 684 | 209 | 14 | 0 | 500ms | not yet calibrated |

*The stampede phase did not run for any cell in this run.*

## Fault injection and controls

*No fault phase ran in this run.*

## Ledger replay (audit)

*No audit phase ran in this run.*

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 4871 | 4577 | 294 | 0 | 0 | 6362 | 0 | 0 | 0 | 0 | 0 |

## Redis accounting (from `INFO`)

Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in
behaviour between the two is partly a policy difference and is stated as such.

| Topology | Scenario | redis_version | maxmemory | maxmemory_policy | maxmemory_samples | used_memory | used_memory_peak | keyspace_hits | keyspace_misses | expired_keys | evicted_keys | total_commands_processed | total_net_input_bytes | total_net_output_bytes | used_cpu_sys | used_cpu_user |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 3497816 | 3568128 | 256555 | 9488 | 5 | 0 | 314721 | 24230069 | 504839849 | 7.366374 | 3.086728 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | 417 | 189 | 20222 | equal-total | 8 |

## Storage

PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a
different query, so none is reported there rather than a number read as bytes on disk.

| Topology | Scenario | Relation | Bytes |
|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:wth:str` | charity | 32.0 KiB |
| pg-single | `o:opt:red:wth:str` | donation | 4.0 MiB |
| pg-single | `o:opt:red:wth:str` | load_donation | 2.3 MiB |
| pg-single | `o:opt:red:wth:str` | person | 168.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |

## Ledger assertion: the requirement must never be ahead of the database

After every writing phase the harness compares, for every donor, the content the database
actually holds against the ledger's freshness requirement. A requirement BEHIND a committed
but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails
the cell, because a bypass read that returned such a state would be an accounting defect
rather than a design finding.

| Topology | Scenario | Phase | donors | mismatches | passed | first example |
|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | mixed | 800 | 0 | true |  |

## Limitations stated with the numbers

- The client, the cache and the database share the same eight heterogeneous cores of a laptop. Absolute
  throughput is lower than a two-machine setup would give, and only relative comparisons within this
  environment are meaningful.
- There is no real network. A shared cache and a distributed database both look better here than they
  would across availability zones; EXPLAIN (ANALYZE, DIST) RPC counts are the portable signal.
- All containers run on one WSL2 machine. This is NOT evidence of data colocation or of network
  behaviour; the placement pair reports engine evidence and physical colocation remains untested.
- The workloads are closed-loop per operation within a phase (measure.Run). Deep tails are therefore
  optimistic floors and are compared between scenarios, never quoted as SLO figures. The open-loop
  scheduler is used only where the report says so.
- A single trial has no error bar. Where trials were repeated the spread is shown per measurement, and
  no conclusion may rest on a difference below the measured noise floor.
- One DeepSeek HIGH agent planned, built, measured and analysed this study. There is no independent
  model review of the implementation, the oracle or the conclusions.
- A cache hit executes no SQL and therefore has no plan. The saved database operation is what a hit
  avoids; the plans file records the statement it avoids.

