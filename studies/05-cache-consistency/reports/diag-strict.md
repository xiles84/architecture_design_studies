# Study 05 — external cache throughput and consistency — run `diag-strict`

Generated 2026-09-21T15:36:06Z from 3 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `diag-strict` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `4a747f5d76d619c528060a9b3b38427d33d1c92d` |
| `git describe` | `study-05/v0.1-handoff-amendment-01-1-g4a747f5-dirty` |
| Working tree dirty at start | true |
| Benchmark image | `localhost/cachebench:1` (`398de52c518c488d2ef`) |
| Resource framing | `db-only` |
| **Inputs digest** | `6c435d614c358b5d` |
| Cells | 3 |
| Distinct scenarios | 3 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 3; cells that failed: 3.
  - failed: `pg-single/legacy-na-memory-through-strict-coord` (STRICT CONTRACT VIOLATED: 338 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/legacy-na-redis-through-strict-coord` (STRICT CONTRACT VIOLATED: 67 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/owned-opt-redis-through-strict-coord` (STRICT CONTRACT VIOLATED: 42 stale-after-ack read(s); the cell cannot support a performance conclusion)
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- Strict cells that violated their contract: 3.
  - `pg-single/legacy-na-memory-through-strict-coord`: 338 stale-after-ack read(s)
  - `pg-single/legacy-na-redis-through-strict-coord`: 67 stale-after-ack read(s)
  - `pg-single/owned-opt-redis-through-strict-coord`: 42 stale-after-ack read(s)
- No cell reported an impossible / uncommitted cache value.
- Highest wrong-read rate recorded: 0.30% of all reads (`pg-single/l:na:mem:wth:str/cacheable_90_reads`), with 338 wrong read(s) in 112519.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `6c435d614c358b5d`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | aborted | 338 | 0 |  |
| pg-single | `legacy-na-redis-through-strict-coord` | legacy | na | redis | through | strict | coord | invalidation-only | aborted | 67 | 0 |  |
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | aborted | 42 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | warm_read_only | cacheable | 15k | 0.01 | 1.62 | 226 | 31505 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | cacheable | 14k | 0.01 | 1.92 | 224 | 30913 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | cacheable | 16k | 0.01 | 1.54 | 224 | 32995 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | cacheable | 145k | 0.03 | 0.13 | 223 | 301371 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:mem:wth:str` | app_99_reads | total app | 11k | 0.02 | 2.59 | 225 | 23900 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:mem:wth:str` | app_90_reads | total app | 5.6k | 0.03 | 10 | 219 | 11806 | 0 | 0.0 | 400 | 0.04% |
| pg-single | `l:na:red:wth:str` | warm_read_only | cacheable | 6.7k | 0.50 | 34 | 54 | 13416 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | cacheable | 7.1k | 0.51 | 29 | 122 | 14659 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | cacheable | 12k | 0.34 | 3.08 | 118 | 23573 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `l:na:red:wth:str` | hotspot_reads | cacheable | 9.4k | 0.36 | 3.48 | 242 | 19912 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `l:na:red:wth:str` | app_99_reads | total app | 8.6k | 0.24 | 4.34 | 223 | 17971 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `l:na:red:wth:str` | app_90_reads | total app | 4.3k | 0.23 | 7.94 | 224 | 9228 | 0 | 0.0 | 300 | 0.10% |
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 7.4k | 0.45 | 24 | 113 | 14810 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 9.8k | 0.39 | 2.96 | 49 | 19557 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 10k | 0.39 | 2.54 | 52 | 20043 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | cacheable | 7.1k | 0.46 | 25 | 225 | 14237 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 10k | 0.18 | 3.60 | 224 | 21014 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 4.0k | 0.22 | 8.26 | 226 | 8817 | 0 | 0.0 | 46 | 0.01% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 765 | 2.03 | 17 | 1530 | 0 | 0.0 |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 506 | 3.49 | 10 | 1018 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 513 | 3.56 | 7.61 | 1026 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 109 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | warm_read_only | 38901 | 38901 | 0 | 0.00% | 31859 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_99_reads | 24905 | 24777 | 0 | 0.00% | 20263 | 0.00% | 0 | 0 | — | — | 0 | 716 | 155 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | 63186 | 63058 | 0 | 0.00% | 51945 | 0.00% | 0 | 0 | — | — | 0 | 1029 | 155 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_90_reads | 73919 | 73479 | 25 | 0.03% | 60483 | 0.04% | 3 | 1 | 6.00 | 21 | 8 | 921 | 539 | 0 | mixed: cache hit served a superseded entry=24; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | 112519 | 111766 | 338 | 0.30% | 92382 | 0.37% | 3 | 1 | 1670 | 2784 | 319 | 2032 | 539 | 0 | mixed: cache hit served a superseded entry=327; mixed: the refilling reader published before the write committed=11 |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | 350883 | 350883 | 0 | 0.00% | 350883 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 350883 | 350883 | 0 | 0.00% | 350883 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_1 | 45710 | 45399 | 9 | 0.02% | 37085 | 0.02% | 1 | 1 | 8295 | 8355 | 9 | 2890 | 412 | 0 | unattributed: cache hit served a superseded entry=7; unattributed: the refilling reader published before the write committed=2 |
| pg-single | `l:na:mem:wth:str` | instances_3 | 11014 | 10770 | 28 | 0.25% | 0 | 0.00% | 3 | 2 | 13 | 35 | 10 | 3192 | 327 | 0 | unattributed: authoritative read returned a superseded snapshot=28 |
| pg-single | `l:na:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 145 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | warm_read_only | 16815 | 16815 | 0 | 0.00% | 16134 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | app_99_reads | 18430 | 18230 | 46 | 0.25% | 18221 | 0.25% | 1 | 1 | 67 | 152 | 41 | 473 | 190 | 0 | mixed: cache hit served a superseded entry=45; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | 36928 | 36728 | 46 | 0.12% | 36701 | 0.13% | 1 | 1 | 67 | 152 | 41 | 702 | 190 | 0 | mixed: cache hit served a superseded entry=45; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | app_90_reads | 45623 | 45046 | 67 | 0.15% | 45065 | 0.15% | 3 | 1 | 43 | 152 | 41 | 803 | 626 | 0 | mixed: cache hit served a superseded entry=65; mixed: the refilling reader published before the write committed=2 |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | 74973 | 74396 | 67 | 0.09% | 74414 | 0.09% | 3 | 1 | 43 | 152 | 41 | 2226 | 626 | 0 | mixed: cache hit served a superseded entry=65; mixed: the refilling reader published before the write committed=2 |
| pg-single | `l:na:red:wth:str` | hotspot_reads | 24042 | 24042 | 0 | 0.00% | 24042 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 24042 | 24042 | 0 | 0.00% | 24042 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_1 | 32879 | 32510 | 30 | 0.09% | 31867 | 0.09% | 3 | 1 | 5.00 | 30 | 10 | 2520 | 430 | 0 | unattributed: cache hit served a superseded entry=30 |
| pg-single | `l:na:red:wth:str` | instances_3 | 31135 | 30798 | 44 | 0.14% | 30817 | 0.14% | 3 | 1 | 5.00 | 12 | 21 | 2842 | 370 | 0 | unattributed: cache hit served a superseded entry=41; unattributed: the refilling reader published before the write committed=3 |
| pg-single | `o:opt:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 150 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 17689 | 17689 | 0 | 0.00% | 16995 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 21862 | 21676 | 0 | 0.00% | 21578 | 0.00% | 0 | 0 | — | — | 0 | 740 | 226 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 46959 | 46773 | 0 | 0.00% | 46666 | 0.00% | 0 | 0 | — | — | 0 | 1479 | 226 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 55162 | 54655 | 4 | 0.01% | 54509 | 0.01% | 1 | 1 | 12 | 21 | 4 | 1101 | 644 | 0 | mixed: cache hit served a superseded entry=4 |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 80342 | 79797 | 42 | 0.05% | 79689 | 0.05% | 2 | 1 | 1575 | 2674 | 38 | 2290 | 644 | 0 | mixed: cache hit served a superseded entry=42 |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | 18460 | 18460 | 0 | 0.00% | 18460 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 18460 | 18460 | 0 | 0.00% | 18460 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_1 | 31410 | 31157 | 0 | 0.00% | 30434 | 0.00% | 0 | 0 | — | — | 0 | 2675 | 322 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_3 | 24283 | 23972 | 0 | 0.00% | 23904 | 0.00% | 0 | 0 | — | — | 0 | 2454 | 414 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 0-75s | 512209 | 5405 | 1.1% | 0 to 0.25 |
| pg-single | `l:na:red:wth:str` | 0-75s | 177274 | 1342 | 0.8% | 0 to 0.25 |
| pg-single | `o:opt:red:wth:str` | 0-75s | 169482 | 1167 | 0.7% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 32312 | 489 | 0 | — | 489 | 143 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.946ms) |
| pg-single | `l:na:red:wth:str` | 2263 | 180 | 0 | — | 180 | 23 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 2.033ms) |
| pg-single | `o:opt:red:wth:str` | 2243 | 180 | 0 | — | 180 | 21 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.072ms) |

*The stampede phase did not run for any cell in this run.*

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 1 | 8 | memory | 5.6k | 0.0 | 13k | 0 | 9 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:str` | 3 | 8 | memory | 2.2k | 0.0 | 2.7k | 11014 | 28 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:red:wth:str` | 1 | 8 | redis | 4.3k | 0.0 | 9.7k | 0 | 30 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 3 | 8 | redis | 3.9k | 0.0 | 9.6k | 0 | 44 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 1 | 8 | redis | 3.0k | 0.0 | 10.0k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 3 | 8 | redis | 4.1k | 0.0 | 6.5k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |

## Fault injection and controls

A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly
has not been fast.

| Topology | Scenario | Fault | reproduced | wrong observed | impossible observed | correctness |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |

## Ledger replay (audit)

*No audit phase ran in this run.*

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | memory | 256.0 KiB | 6.7 KiB | 3 | 26638 | 30646 | 30082 | 564 | 2 | 0 | 7177 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 3935 | 3511 | 424 | 2 | 0 | 5393 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 4082 | 3838 | 244 | 2 | 0 | 5469 | 0 | 20 | 0 | 0 |

## Redis accounting (from `INFO`)

Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in
behaviour between the two is partly a policy difference and is stated as such.

| Topology | Scenario | redis_version | maxmemory | maxmemory_policy | maxmemory_samples | used_memory | used_memory_peak | keyspace_hits | keyspace_misses | expired_keys | evicted_keys | total_commands_processed | total_net_input_bytes | total_net_output_bytes | used_cpu_sys | used_cpu_user |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1368456 | 3961592 | 401886 | 30493 | 2 | 0 | 544696 | 58567237 | 757116988 | 11.639787 | 6.116129 |
| pg-single | `o:opt:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1374888 | 3893936 | 196590 | 15479 | 1 | 0 | 267579 | 29023152 | 370592004 | 5.653999 | 3.049761 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:str` | 0 | 0 | 0 | db-only | 8 |

## Storage

PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a
different query, so none is reported there rather than a number read as bytes on disk.

| Topology | Scenario | Relation | Bytes |
|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:wth:str` | donation | 4.0 MiB |
| pg-single | `l:na:mem:wth:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:mem:wth:str` | person | 168.0 KiB |
| pg-single | `l:na:red:wth:str` | charity | 32.0 KiB |
| pg-single | `l:na:red:wth:str` | donation | 4.0 MiB |
| pg-single | `l:na:red:wth:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:red:wth:str` | person | 168.0 KiB |
| pg-single | `o:opt:red:wth:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:wth:str` | charity | 32.0 KiB |
| pg-single | `o:opt:red:wth:str` | donation | 4.0 MiB |
| pg-single | `o:opt:red:wth:str` | load_donation | 2.3 MiB |
| pg-single | `o:opt:red:wth:str` | person | 168.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |

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

