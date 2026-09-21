# Study 05 — external cache throughput and consistency — run `verify-ledger`

Generated 2026-09-21T15:57:08Z from 6 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `verify-ledger` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `43a299008b47b23d332dd9ab0b39405a2544e326` |
| `git describe` | `study-05/v1.1-diagnostics-dirty` |
| Working tree dirty at start | true |
| Benchmark image | `localhost/cachebench:1` (`ab21cfd658926ebefc5`) |
| Resource framing | `db-only` |
| **Inputs digest** | `b905f5495b14d41f` |
| Cells | 6 |
| Distinct scenarios | 6 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 6; cells that failed: 4.
  - failed: `pg-single/legacy-na-memory-through-strict-coord` (STRICT CONTRACT VIOLATED: 7 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/legacy-na-redis-through-strict-coord` (STRICT CONTRACT VIOLATED: 122 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/owned-opt-redis-through-strict-coord` (STRICT CONTRACT VIOLATED: 32 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/ref-rollup-trigger` (IMPOSSIBLE CACHE VALUE: 114 read(s) returned a state the database never held; the cell is invalid)
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- Strict cells that violated their contract: 3.
  - `pg-single/legacy-na-memory-through-strict-coord`: 7 stale-after-ack read(s)
  - `pg-single/legacy-na-redis-through-strict-coord`: 122 stale-after-ack read(s)
  - `pg-single/owned-opt-redis-through-strict-coord`: 32 stale-after-ack read(s)
- Cells with impossible cache values: 1 (**these cells are invalid under every freshness policy**).
  - `pg-single/ref-rollup-trigger`: 114 impossible cache value(s)
- Highest and lowest measured throughput across all cells: 4.9k (`pg-single/ref-norm/cacheable_99_reads`) and 3.1k (`pg-single/ref-norm/hotspot_reads`).
- Highest wrong-read rate recorded: 14.08% of all reads (`pg-single/l:na:mem:wth:str/instances_1`), with 6554 wrong read(s) in 46544.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `b905f5495b14d41f`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | aborted | 7 | 0 |  |
| pg-single | `legacy-na-redis-through-strict-coord` | legacy | na | redis | through | strict | coord | invalidation-only | aborted | 122 | 0 |  |
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | aborted | 32 | 0 |  |
| pg-single | `ref-embedded-locked` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |
| pg-single | `ref-normalized-indexed` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |
| pg-single | `ref-rollup-trigger` | ref |  | none |  |  |  | invalidation-only | aborted | 0 | 114 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | warm_read_only | cacheable | 15k | 0.01 | 1.92 | 222 | 31605 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | cacheable | 13k | 0.02 | 1.79 | 220 | 28445 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | cacheable | 15k | 0.02 | 2.33 | 216 | 31452 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | cacheable | 126k | 0.02 | 0.10 | 218 | 270449 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:mem:wth:str` | app_99_reads | total app | 11k | 0.02 | 2.41 | 225 | 24826 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:mem:wth:str` | app_90_reads | total app | 5.3k | 0.03 | 6.48 | 222 | 11667 | 0 | 0.0 | 6642 | 0.62% |
| pg-single | `l:na:red:wth:str` | warm_read_only | cacheable | 8.6k | 0.46 | 11 | 145 | 17481 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | cacheable | 9.0k | 0.45 | 5.73 | 53 | 17960 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | cacheable | 10k | 0.37 | 3.31 | 53 | 20962 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `l:na:red:wth:str` | hotspot_reads | cacheable | 13k | 0.28 | 2.46 | 210 | 28072 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `l:na:red:wth:str` | app_99_reads | total app | 8.5k | 0.22 | 4.91 | 225 | 18226 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `l:na:red:wth:str` | app_90_reads | total app | 3.8k | 0.24 | 10 | 225 | 8380 | 0 | 0.0 | 539 | 0.16% |
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 7.8k | 0.49 | 13 | 191 | 15707 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 8.3k | 0.51 | 9.52 | 47 | 16574 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 8.5k | 0.52 | 8.21 | 50 | 17329 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | cacheable | 9.2k | 0.43 | 5.76 | 205 | 19178 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 9.4k | 0.21 | 4.37 | 224 | 19760 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 4.1k | 0.20 | 13 | 219 | 8761 | 0 | 0.0 | 105 | 0.03% |
| pg-single | `ref-emb` | warm_read_only | cacheable | 3.4k | 1.61 | 32 | 37 | 6760 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-emb` | cacheable_99_reads | cacheable | 3.6k | 1.31 | 30 | 48 | 7146 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-emb` | cacheable_90_reads | cacheable | 3.9k | 1.42 | 20 | 41 | 7858 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-emb` | hotspot_reads | cacheable | 4.5k | 1.15 | 19 | 46 | 9002 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-emb` | app_99_reads | total app | 5.3k | 0.96 | 28 | 38 | 10555 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-emb` | app_90_reads | total app | 2.8k | 1.73 | 35 | 58 | 5593 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | warm_read_only | cacheable | 4.7k | 1.13 | 23 | 33 | 9469 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | cacheable_99_reads | cacheable | 4.9k | 1.23 | 14 | 30 | 9793 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | cacheable_90_reads | cacheable | 3.6k | 1.58 | 29 | 44 | 7262 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | hotspot_reads | cacheable | 3.1k | 1.75 | 31 | 39 | 6282 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | app_99_reads | total app | 3.2k | 1.52 | 26 | 53 | 6437 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-norm` | app_90_reads | total app | 2.6k | 1.93 | 33 | 41 | 5260 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | warm_read_only | cacheable | 3.5k | 1.49 | 29 | 37 | 7077 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_99_reads | cacheable | 4.0k | 1.25 | 31 | 40 | 8006 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_90_reads | cacheable | 4.0k | 1.26 | 33 | 47 | 8055 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | hotspot_reads | cacheable | 4.2k | 1.22 | 31 | 53 | 8378 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_99_reads | total app | 4.7k | 0.96 | 26 | 48 | 9474 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_90_reads | total app | 4.6k | 1.06 | 24 | 34 | 9139 | 0 | 0.0 | 0 | 0.00% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 935 | 1.85 | 6.17 | 1872 | 0 | 0.0 |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 720 | 2.58 | 5.66 | 1440 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 659 | 2.63 | 7.77 | 1320 | 0 | 0.0 |
| pg-single | `ref-emb` | hotspot_writes | 593 | 3.26 | 6.20 | 1186 | 0 | 0.0 |
| pg-single | `ref-norm` | hotspot_writes | 1.4k | 1.11 | 4.58 | 2799 | 0 | 0.0 |
| pg-single | `ref-roll` | hotspot_writes | 649 | 2.92 | 21 | 1299 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 101 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | warm_read_only | 36921 | 36921 | 0 | 0.00% | 30328 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_99_reads | 24713 | 24597 | 0 | 0.00% | 20095 | 0.00% | 0 | 0 | — | — | 0 | 314 | 143 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | 59396 | 59280 | 0 | 0.00% | 48693 | 0.00% | 0 | 0 | — | — | 0 | 1248 | 143 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_90_reads | 69816 | 69398 | 7 | 0.01% | 56812 | 0.01% | 2 | 1 | 3.00 | 11 | 3 | 666 | 516 | 0 | mixed: cache hit served a superseded entry=7 |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | 105951 | 105533 | 7 | 0.01% | 86722 | 0.01% | 2 | 1 | 3.00 | 11 | 3 | 2500 | 516 | 0 | mixed: cache hit served a superseded entry=7 |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | 353764 | 353764 | 0 | 0.00% | 353764 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 353764 | 353764 | 0 | 0.00% | 353764 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_1 | 46544 | 39675 | 6554 | 14.08% | 37661 | 17.40% | 2 | 1 | 1717 | 2797 | 6551 | 2846 | 425 | 0 | unattributed: cache hit served a superseded entry=6552; unattributed: the refilling reader published before the write committed=2 |
| pg-single | `l:na:mem:wth:str` | instances_3 | 12612 | 12289 | 74 | 0.59% | 0 | 0.00% | 4 | 1 | 42 | 3046 | 26 | 2496 | 365 | 0 | unattributed: authoritative read returned a superseded snapshot=74 |
| pg-single | `l:na:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 148 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | warm_read_only | 22756 | 22756 | 0 | 0.00% | 22020 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | app_99_reads | 19252 | 19017 | 96 | 0.50% | 19069 | 0.50% | 1 | 1 | 19 | 38 | 79 | 495 | 162 | 0 | mixed: cache hit served a superseded entry=95; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | 41695 | 41460 | 96 | 0.23% | 41501 | 0.23% | 1 | 1 | 19 | 38 | 79 | 1560 | 162 | 0 | mixed: cache hit served a superseded entry=95; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | app_90_reads | 49240 | 48688 | 122 | 0.25% | 48739 | 0.25% | 2 | 1 | 15 | 38 | 79 | 1075 | 536 | 0 | mixed: cache hit served a superseded entry=121; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | 77284 | 76732 | 122 | 0.16% | 76782 | 0.16% | 2 | 1 | 15 | 38 | 79 | 2530 | 536 | 0 | mixed: cache hit served a superseded entry=121; mixed: the refilling reader published before the write committed=1 |
| pg-single | `l:na:red:wth:str` | hotspot_reads | 36947 | 36947 | 0 | 0.00% | 36947 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 36947 | 36947 | 0 | 0.00% | 36947 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_1 | 33988 | 33594 | 47 | 0.14% | 32903 | 0.14% | 2 | 1 | 23 | 58 | 35 | 1725 | 446 | 0 | unattributed: cache hit served a superseded entry=45; unattributed: the refilling reader published before the write committed=2 |
| pg-single | `l:na:red:wth:str` | instances_3 | 28132 | 27749 | 56 | 0.20% | 27762 | 0.20% | 2 | 1 | 13 | 209 | 25 | 2621 | 418 | 0 | unattributed: cache hit served a superseded entry=55; unattributed: the refilling reader published before the write committed=1 |
| pg-single | `o:opt:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 152 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 19312 | 19312 | 0 | 0.00% | 18600 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 20191 | 20011 | 0 | 0.00% | 19950 | 0.00% | 0 | 0 | — | — | 0 | 280 | 221 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 41261 | 41081 | 0 | 0.00% | 41006 | 0.00% | 0 | 0 | — | — | 0 | 1055 | 221 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 49463 | 48950 | 32 | 0.06% | 48876 | 0.07% | 2 | 1 | 14 | 41 | 18 | 1690 | 607 | 0 | mixed: cache hit served a superseded entry=32 |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 70985 | 70472 | 32 | 0.05% | 70398 | 0.05% | 2 | 1 | 14 | 41 | 18 | 3063 | 607 | 0 | mixed: cache hit served a superseded entry=32 |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | 25894 | 25894 | 0 | 0.00% | 25894 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 25894 | 25894 | 0 | 0.00% | 25894 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_1 | 27366 | 27050 | 39 | 0.14% | 26356 | 0.15% | 1 | 1 | 50 | 107 | 39 | 2487 | 365 | 0 | unattributed: cache hit served a superseded entry=38; unattributed: the refilling reader published before the write committed=1 |
| pg-single | `o:opt:red:wth:str` | instances_3 | 27848 | 27576 | 2 | 0.01% | 27496 | 0.01% | 1 | 1 | 1.00 | 1.00 | 2 | 2577 | 367 | 0 | unattributed: cache hit served a superseded entry=2 |
| pg-single | `ref-emb` | warm_read_only | 8493 | 8493 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-emb` | app_99_reads | 10832 | 10813 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 446 | 44 | 0 | — |
| pg-single | `ref-emb` | cacheable_99_reads | 20356 | 20337 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 3078 | 44 | 0 | — |
| pg-single | `ref-emb` | app_90_reads | 26220 | 26125 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1629 | 240 | 0 | — |
| pg-single | `ref-emb` | cacheable_90_reads | 35688 | 35593 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 3078 | 240 | 0 | — |
| pg-single | `ref-emb` | hotspot_reads | 11098 | 11098 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-emb` | hotspot_writes | 11098 | 11098 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-norm` | warm_read_only | 11884 | 11884 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-norm` | app_99_reads | 6441 | 6424 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 993 | 31 | 0 | — |
| pg-single | `ref-norm` | cacheable_99_reads | 18321 | 18304 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2195 | 31 | 0 | — |
| pg-single | `ref-norm` | app_90_reads | 23394 | 23317 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1088 | 228 | 0 | — |
| pg-single | `ref-norm` | cacheable_90_reads | 32701 | 32624 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 3026 | 228 | 0 | — |
| pg-single | `ref-norm` | hotspot_reads | 8481 | 8481 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-norm` | hotspot_writes | 8481 | 8481 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | warm_read_only | 8883 | 8883 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | app_99_reads | 9741 | 9348 | 0 | 0.00% | 0 | 0.00% | 1 | 0 | — | — | 0 | 483 | 63 | 10 | — |
| pg-single | `ref-roll` | cacheable_99_reads | 19447 | 19054 | 0 | 0.00% | 0 | 0.00% | 1 | 0 | — | — | 0 | 1958 | 63 | 10 | — |
| pg-single | `ref-roll` | app_90_reads | 28055 | 27239 | 0 | 0.00% | 0 | 0.00% | 3 | 0 | — | — | 0 | 3326 | 381 | 47 | — |
| pg-single | `ref-roll` | cacheable_90_reads | 38208 | 37375 | 0 | 0.00% | 0 | 0.00% | 3 | 0 | — | — | 0 | 3732 | 381 | 47 | — |
| pg-single | `ref-roll` | hotspot_reads | 10627 | 10627 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | hotspot_writes | 10627 | 10627 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 0-75s | 508475 | 9462 | 1.9% | 0 to 0.25 |
| pg-single | `l:na:red:wth:str` | 0-75s | 196414 | 1500 | 0.8% | 0 to 0.25 |
| pg-single | `o:opt:red:wth:str` | 0-75s | 168744 | 1405 | 0.8% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 35001 | 530 | 0 | — | 530 | 126 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.479ms) |
| pg-single | `l:na:red:wth:str` | 2465 | 185 | 0 | — | 185 | 18 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.669ms) |
| pg-single | `o:opt:red:wth:str` | 2433 | 200 | 0 | — | 200 | 19 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 3.44ms) |
| pg-single | `ref-emb` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |
| pg-single | `ref-norm` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |
| pg-single | `ref-roll` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |

| Topology | Scenario | readers | keys | elapsed ms | db loads | loads/key | fallbacks | duplicate fills | wrong | impossible |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 16 | 8 | 211.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | 16 | 8 | 226.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | 16 | 8 | 227.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 1 | 8 | memory | 5.8k | 0.0 | 13k | 0 | 6554 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:str` | 3 | 8 | memory | 3.1k | 0.0 | 2.7k | 12612 | 74 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:red:wth:str` | 1 | 8 | redis | 6.2k | 0.0 | 8.5k | 0 | 47 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 3 | 8 | redis | 4.5k | 0.0 | 8.0k | 0 | 56 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 1 | 8 | redis | 3.7k | 0.0 | 7.6k | 0 | 39 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 3 | 8 | redis | 3.7k | 0.0 | 8.3k | 0 | 2 | 0 | shared store: invalidation is visible to every instance |

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

| Topology | Scenario | Phase # | passed | checks | content mismatches | donation-set mismatches | aggregate mismatches | recent-slice mismatches | design checks | first failure |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:wth:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `o:opt:red:wth:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `ref-emb` | 1 | true | 14 | 0 | 0 | 0 | 0 | a_embedded_drift=0 a_rolldown_drift=0 |  |
| pg-single | `ref-norm` | 1 | true | 12 | 0 | 0 | 0 | 0 |  |  |
| pg-single | `ref-roll` | 1 | true | 14 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 a_rollup_drift=0 |  |

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | memory | 256.0 KiB | 6.5 KiB | 3 | 25321 | 29626 | 29001 | 625 | 2 | 0 | 7801 | 0 | 20 | 0 | 3 | 0 |
| pg-single | `l:na:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 4415 | 3928 | 487 | 2 | 0 | 6331 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 4141 | 3881 | 260 | 2 | 0 | 5663 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `ref-emb` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-norm` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-roll` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Redis accounting (from `INFO`)

Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in
behaviour between the two is partly a policy difference and is stated as such.

| Topology | Scenario | redis_version | maxmemory | maxmemory_policy | maxmemory_samples | used_memory | used_memory_peak | keyspace_hits | keyspace_misses | expired_keys | evicted_keys | total_commands_processed | total_net_input_bytes | total_net_output_bytes | used_cpu_sys | used_cpu_user |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1418064 | 3927856 | 228644 | 16130 | 1 | 0 | 309195 | 33251735 | 427706719 | 5.922485 | 2.924627 |
| pg-single | `o:opt:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1561552 | 3984272 | 196580 | 15452 | 1 | 0 | 268556 | 29129115 | 369026795 | 5.415058 | 3.007016 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ref-emb` | 210 | 15 | 299 | db-only | 8 |
| pg-single | `ref-norm` | 210 | 42 | 1473 | db-only | 8 |
| pg-single | `ref-roll` | 0 | 0 | 0 | db-only | 8 |

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
| pg-single | `ref-emb` | charity | 32.0 KiB |
| pg-single | `ref-emb` | donation | 4.0 MiB |
| pg-single | `ref-emb` | load_donation | 2.3 MiB |
| pg-single | `ref-emb` | person | 720.0 KiB |
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
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `ref-emb` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
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
| pg-single | `l:na:mem:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | faults | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:red:wth:str` | faults | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | mixed | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | hotspot | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | stampede | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | instances | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | faults | 800 | 0 | true |  |
| pg-single | `ref-emb` | warm | 800 | 0 | true |  |
| pg-single | `ref-emb` | mixed | 800 | 0 | true |  |
| pg-single | `ref-emb` | hotspot | 800 | 0 | true |  |
| pg-single | `ref-emb` | stampede | 800 | 0 | true |  |
| pg-single | `ref-emb` | instances | 800 | 0 | true |  |
| pg-single | `ref-emb` | faults | 800 | 0 | true |  |
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

