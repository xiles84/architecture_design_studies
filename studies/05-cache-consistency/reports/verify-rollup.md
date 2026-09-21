# Study 05 — external cache throughput and consistency — run `verify-rollup`

Generated 2026-09-21T16:20:21Z from 4 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `verify-rollup` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `ead9b18e5d3631c5ba379b7a37dd1b0062e16608` |
| `git describe` | `study-05/v2-owner-guard-dirty` |
| Working tree dirty at start | true |
| Benchmark image | `localhost/cachebench:1` (`ec45913dd98e5b1965c`) |
| Resource framing | `db-only` |
| **Inputs digest** | `27649ef92847a133` |
| Cells | 4 |
| Distinct scenarios | 4 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 4; cells that failed: 2.
  - failed: `pg-single/legacy-na-memory-aside-strict-coord` (STRICT CONTRACT VIOLATED: 10 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/ref-rollup-trigger` (phase audit: ledger replay failed: INV-7: donor 1 portal content differs from the ledger-derived expectation (total 39637861 vs 39635688); INV-10: a_rollup_drif…)
- Negative controls that fired: 1. Controls that did NOT fire: 0.
- Strict cells that violated their contract: 1.
  - `pg-single/legacy-na-memory-aside-strict-coord`: 10 stale-after-ack read(s)
- No cell reported an impossible / uncommitted cache value.
- Highest and lowest measured throughput across all cells: 3.7k (`pg-single/ref-norm/warm_read_only`) and 2.9k (`pg-single/ref-norm/cacheable_90_reads`).
- Highest wrong-read rate recorded: 0.17% of all reads (`pg-single/l:na:mem:asd:str/instances_3`), with 22 wrong read(s) in 12581.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `27649ef92847a133`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale-invalidation` | legacy | na | memory | aside | relaxed | coord | invalidation-only | pass | 0 | 0 | fired |
| pg-single | `legacy-na-memory-aside-strict-coord` | legacy | na | memory | aside | strict | coord | invalidation-only | aborted | 10 | 0 |  |
| pg-single | `ref-normalized-indexed` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |
| pg-single | `ref-rollup-trigger` | ref |  | none |  |  |  | invalidation-only | aborted | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | warm_read_only | cacheable | 14k | 0.01 | 1.88 | 225 | 31724 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `ctl-stale` | cacheable_99_reads | cacheable | 14k | 0.02 | 2.08 | 224 | 30627 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `ctl-stale` | cacheable_90_reads | cacheable | 15k | 0.02 | 1.95 | 225 | 30327 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `ctl-stale` | hotspot_reads | cacheable | 127k | 0.02 | 0.14 | 223 | 271793 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `ctl-stale` | app_99_reads | total app | 10k | 0.02 | 2.01 | 223 | 22391 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `ctl-stale` | app_90_reads | total app | 4.3k | 0.04 | 7.53 | 223 | 9340 | 0 | 0.0 | 6 | 0.00% |
| pg-single | `l:na:mem:asd:str` | warm_read_only | cacheable | 13k | 0.02 | 2.42 | 224 | 28489 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | cacheable | 15k | 0.01 | 1.52 | 224 | 31937 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | cacheable | 13k | 0.02 | 2.03 | 227 | 28906 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | cacheable | 123k | 0.03 | 0.15 | 223 | 271262 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_99_reads | total app | 8.7k | 0.02 | 3.27 | 227 | 18195 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_90_reads | total app | 4.0k | 0.05 | 27 | 220 | 8659 | 0 | 0.0 | 42 | 0.00% |
| pg-single | `ref-norm` | warm_read_only | cacheable | 3.7k | 1.34 | 26 | 42 | 7444 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | cacheable_99_reads | cacheable | 3.6k | 1.47 | 27 | 51 | 7287 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | cacheable_90_reads | cacheable | 2.9k | 1.80 | 34 | 40 | 5785 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | hotspot_reads | cacheable | 3.1k | 1.69 | 34 | 51 | 6167 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | app_99_reads | total app | 2.9k | 1.59 | 37 | 55 | 5875 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | app_90_reads | total app | 2.7k | 1.81 | 34 | 64 | 5383 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | warm_read_only | cacheable | 4.1k | 1.29 | 29 | 39 | 8252 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_99_reads | cacheable | 5.1k | 1.07 | 16 | 41 | 10298 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_90_reads | cacheable | 4.1k | 1.27 | 21 | 42 | 8252 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | hotspot_reads | cacheable | 4.0k | 1.28 | 30 | 39 | 7907 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_99_reads | total app | 4.2k | 1.20 | 27 | 53 | 8431 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_90_reads | total app | 4.0k | 1.25 | 26 | 81 | 8098 | 0 | 0.0 | 0 | 0.00% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | hotspot_writes | 1.0k | 1.73 | 3.70 | 2100 | 0 | 0.0 |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 1.0k | 1.37 | 4.55 | 2086 | 0 | 0.0 |
| pg-single | `ref-norm` | hotspot_writes | 1.1k | 1.13 | 6.84 | 2236 | 0 | 0.0 |
| pg-single | `ref-roll` | hotspot_writes | 633 | 2.91 | 8.29 | 1266 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | cold_fill | 400 | 400 | 0 | 0.00% | 105 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | warm_read_only | 37322 | 37322 | 0 | 0.00% | 30540 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | app_99_reads | 23522 | 23472 | 0 | 0.00% | 19148 | 0.00% | 0 | 0 | — | — | 0 | 542 | 51 | 0 | — |
| pg-single | `ctl-stale` | cacheable_99_reads | 60117 | 60067 | 0 | 0.00% | 49579 | 0.00% | 0 | 0 | — | — | 0 | 1245 | 51 | 0 | — |
| pg-single | `ctl-stale` | app_90_reads | 68519 | 68380 | 0 | 0.00% | 55773 | 0.00% | 0 | 0 | — | — | 0 | 1132 | 161 | 0 | — |
| pg-single | `ctl-stale` | cacheable_90_reads | 103437 | 103298 | 0 | 0.00% | 84439 | 0.00% | 0 | 0 | — | — | 0 | 2227 | 161 | 0 | — |
| pg-single | `ctl-stale` | hotspot_reads | 344504 | 344504 | 0 | 0.00% | 344504 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | hotspot_writes | 344504 | 344504 | 0 | 0.00% | 344504 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | instances_1 | 41273 | 41170 | 6 | 0.01% | 33212 | 0.02% | 1 | 1 | 2.00 | 7.00 | 6 | 2474 | 111 | 0 | unattributed: cache hit served a superseded entry=5; unattributed: the refilling reader published before the write committed=1 |
| pg-single | `ctl-stale` | instances_3 | 18829 | 18713 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2366 | 313 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cold_fill | 400 | 400 | 0 | 0.00% | 94 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | warm_read_only | 34282 | 34282 | 0 | 0.00% | 27951 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_99_reads | 20016 | 19958 | 0 | 0.00% | 16322 | 0.00% | 0 | 0 | — | — | 0 | 1260 | 70 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | 57489 | 57431 | 0 | 0.00% | 47137 | 0.00% | 0 | 0 | — | — | 0 | 1665 | 70 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_90_reads | 65120 | 64933 | 10 | 0.02% | 52659 | 0.02% | 1 | 1 | 12 | 26 | 10 | 1156 | 226 | 0 | mixed: cache hit served a superseded entry=9; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | 99523 | 99336 | 10 | 0.01% | 80909 | 0.01% | 1 | 1 | 12 | 26 | 10 | 2539 | 226 | 0 | mixed: cache hit served a superseded entry=9; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | 349903 | 349903 | 0 | 0.00% | 349903 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 349903 | 349903 | 0 | 0.00% | 349903 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_1 | 44610 | 44494 | 0 | 0.00% | 35937 | 0.00% | 0 | 0 | — | — | 0 | 4439 | 173 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_3 | 12581 | 12499 | 22 | 0.17% | 0 | 0.00% | 2 | 1 | 21 | 33 | 15 | 2207 | 177 | 0 | unattributed: authoritative read returned a superseded snapshot=22 |
| pg-single | `ref-norm` | warm_read_only | 9885 | 9885 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-norm` | app_99_reads | 7021 | 7007 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 688 | 30 | 0 | — |
| pg-single | `ref-norm` | cacheable_99_reads | 16165 | 16151 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2258 | 30 | 0 | — |
| pg-single | `ref-norm` | app_90_reads | 21392 | 21317 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 883 | 187 | 0 | — |
| pg-single | `ref-norm` | cacheable_90_reads | 29042 | 28967 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 3106 | 187 | 0 | — |
| pg-single | `ref-norm` | hotspot_reads | 8663 | 8663 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-norm` | hotspot_writes | 8663 | 8663 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | warm_read_only | 10346 | 10346 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | app_99_reads | 9562 | 9159 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 442 | 36 | 0 | — |
| pg-single | `ref-roll` | cacheable_99_reads | 21826 | 21423 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1879 | 36 | 0 | — |
| pg-single | `ref-roll` | app_90_reads | 30122 | 29038 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 844 | 330 | 0 | — |
| pg-single | `ref-roll` | cacheable_90_reads | 40418 | 38223 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2668 | 330 | 0 | — |
| pg-single | `ref-roll` | hotspot_reads | 10783 | 5397 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | hotspot_writes | 10783 | 5397 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 0-75s | 492695 | 6331 | 1.3% | 0 to 0.25 |
| pg-single | `l:na:mem:asd:str` | 0-75s | 494700 | 6419 | 1.3% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 31880 | 525 | 0 | — | 525 | 117 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 2.643ms) |
| pg-single | `l:na:mem:asd:str` | 31108 | 537 | 0 | — | 537 | 138 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.46ms) |
| pg-single | `ref-norm` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |
| pg-single | `ref-roll` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |

| Topology | Scenario | readers | keys | elapsed ms | db loads | loads/key | fallbacks | duplicate fills | wrong | impossible |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 1 | 8 | memory | 4.2k | 0.0 | 12k | 0 | 6 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `ctl-stale` | 3 | 8 | memory | 3.9k | 0.0 | 4.6k | 18829 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:asd:str` | 1 | 8 | memory | 3.4k | 0.0 | 14k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:str` | 3 | 8 | memory | 2.7k | 0.0 | 2.8k | 12581 | 22 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |

## Fault injection and controls

A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly
has not been fast.

| Topology | Scenario | Fault | reproduced | wrong observed | impossible observed | correctness |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `ctl-stale` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `ctl-stale` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `ctl-stale` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `ctl-stale` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `ctl-stale` | suppressed-invalidation-produces-a-stale-read | true | 1 | 0 | one suppressed post-commit invalidation produced 1 stale-after-ack read(s) that the oracle detected |
| pg-single | `ctl-stale` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |

## Ledger replay (audit)

| Topology | Scenario | Phase # | passed | checks | content mismatches | donation-set mismatches | aggregate mismatches | recent-slice mismatches | design checks | first failure |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:asd:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `ref-norm` | 1 | true | 12 | 0 | 0 | 0 | 0 |  |  |
| pg-single | `ref-roll` | 1 | false | 14 | 1 | 0 | 0 | 0 | a_rolldown_drift=0 a_rollup_drift=3 | INV-7: donor 1 portal content differs from the ledger-derived expectation (total 39637861 vs 39635688) |

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | memory | 256.0 KiB | 4.4 KiB | 2 | 24439 | 25696 | 25673 | 23 | 0 | 4136 | 2 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | memory | 256.0 KiB | 2.2 KiB | 1 | 23472 | 24840 | 24750 | 90 | 0 | 0 | 7731 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `ref-norm` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-roll` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 347 | 10 | 177 | db-only | 8 |
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ref-norm` | 206 | 28 | 1350 | db-only | 8 |
| pg-single | `ref-roll` | 0 | 0 | 0 | db-only | 8 |

## Storage

PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a
different query, so none is reported there rather than a number read as bytes on disk.

| Topology | Scenario | Relation | Bytes |
|---|---|---|---|
| pg-single | `ctl-stale` | charity | 32.0 KiB |
| pg-single | `ctl-stale` | donation | 4.0 MiB |
| pg-single | `ctl-stale` | load_donation | 2.3 MiB |
| pg-single | `ctl-stale` | person | 168.0 KiB |
| pg-single | `l:na:mem:asd:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:asd:str` | donation | 4.0 MiB |
| pg-single | `l:na:mem:asd:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:mem:asd:str` | person | 168.0 KiB |
| pg-single | `ref-norm` | charity | 32.0 KiB |
| pg-single | `ref-norm` | donation | 3.0 MiB |
| pg-single | `ref-norm` | load_donation | 2.3 MiB |
| pg-single | `ref-norm` | person | 168.0 KiB |
| pg-single | `ref-roll` | charity | 32.0 KiB |
| pg-single | `ref-roll` | donation | 4.0 MiB |
| pg-single | `ref-roll` | load_donation | 2.3 MiB |
| pg-single | `ref-roll` | person | 280.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `ctl-stale` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `ref-norm` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `ref-roll` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |

## Ledger assertion: the requirement must never be ahead of the database

After every writing phase the harness compares, for every donor, the content the database
actually holds against the ledger's freshness requirement. A requirement BEHIND a committed
but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails
the cell, because a bypass read that returned such a state would be an accounting defect
rather than a design finding.

| Topology | Scenario | Phase | donors | mismatches | passed | first example |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | warm | 800 | 0 | true |  |
| pg-single | `ctl-stale` | mixed | 800 | 0 | true |  |
| pg-single | `ctl-stale` | hotspot | 800 | 0 | true |  |
| pg-single | `ctl-stale` | stampede | 800 | 0 | true |  |
| pg-single | `ctl-stale` | instances | 800 | 0 | true |  |
| pg-single | `ctl-stale` | faults | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | faults | 800 | 0 | true |  |
| pg-single | `ref-norm` | warm | 800 | 0 | true |  |
| pg-single | `ref-norm` | mixed | 800 | 0 | true |  |
| pg-single | `ref-norm` | hotspot | 800 | 0 | true |  |
| pg-single | `ref-norm` | stampede | 800 | 0 | true |  |
| pg-single | `ref-norm` | instances | 800 | 0 | true |  |
| pg-single | `ref-norm` | faults | 800 | 0 | true |  |
| pg-single | `ref-roll` | warm | 800 | 0 | true |  |
| pg-single | `ref-roll` | mixed | 800 | 0 | true |  |
| pg-single | `ref-roll` | hotspot | 800 | 0 | true |  |
| pg-single | `ref-roll` | stampede | 800 | 0 | true |  |
| pg-single | `ref-roll` | instances | 800 | 0 | true |  |
| pg-single | `ref-roll` | faults | 800 | 0 | true |  |

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

