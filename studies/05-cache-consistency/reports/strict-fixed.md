# Study 05 — external cache throughput and consistency — run `strict-fixed`

Generated 2026-09-21T16:42:16Z from 5 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `strict-fixed` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `f9d9f57894de7e50094d75beb3ac6887841b8577` |
| `git describe` | `study-05/v2.2-residual-classified-dirty` |
| Working tree dirty at start | true |
| Benchmark image | `localhost/cachebench:1` (`ba2d7567777164ffa2d`) |
| Resource framing | `db-only` |
| **Inputs digest** | `419d7fd52f5418a9` |
| Cells | 5 |
| Distinct scenarios | 5 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 5; cells that failed: 1.
  - failed: `pg-single/legacy-na-redis-aside-strict-coord` (phase ackcheck: ACK CHECK FAILED: 1 of 310 acknowledgements were not visible to a second session immediately afterwards (concurrent: acknowledged donor 2 and a …)
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No strict cell recorded a stale-after-ack read.
- No cell reported an impossible / uncommitted cache value.
- Highest and lowest measured throughput across all cells: 162k (`pg-single/l:na:mem:wth:str/hotspot_reads`) and 5.9k (`pg-single/ref-roll/cacheable_90_reads`).

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `419d7fd52f5418a9`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `legacy-na-memory-aside-strict-coord` | legacy | na | memory | aside | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-strict-coord` | legacy | na | redis | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `ref-rollup-trigger` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | warm_read_only | cacheable | 17k | 0.01 | 1.01 | 226 | 37364 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | cacheable | 17k | 0.01 | 1.16 | 226 | 35356 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | cacheable | 18k | 0.01 | 1.04 | 221 | 39650 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | cacheable | 157k | 0.02 | 0.10 | 214 | 330435 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_99_reads | total app | 13k | 0.02 | 1.17 | 222 | 28037 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_90_reads | total app | 5.4k | 0.03 | 2.06 | 223 | 11411 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | warm_read_only | cacheable | 12k | 0.02 | 1.53 | 225 | 26759 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | cacheable | 13k | 0.02 | 1.84 | 219 | 29346 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | cacheable | 11k | 0.02 | 2.71 | 226 | 23493 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | cacheable | 162k | 0.02 | 0.07 | 220 | 353010 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | app_99_reads | total app | 11k | 0.02 | 2.65 | 222 | 23550 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | app_90_reads | total app | 5.3k | 0.04 | 4.97 | 218 | 11334 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | warm_read_only | cacheable | 15k | 0.22 | 1.96 | 193 | 30859 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | cacheable | 17k | 0.20 | 0.98 | 172 | 33497 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | cacheable | 15k | 0.22 | 1.03 | 157 | 30200 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | hotspot_reads | cacheable | 17k | 0.18 | 0.85 | 194 | 34621 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_99_reads | total app | 10.0k | 0.16 | 1.65 | 226 | 21823 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_90_reads | total app | 4.2k | 0.15 | 2.99 | 219 | 9187 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 10k | 0.30 | 4.16 | 158 | 21241 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 15k | 0.22 | 1.15 | 173 | 30728 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 12k | 0.29 | 1.89 | 164 | 23688 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | cacheable | 17k | 0.17 | 0.92 | 194 | 36006 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 8.9k | 0.23 | 4.24 | 221 | 19107 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 3.9k | 0.25 | 9.27 | 223 | 8368 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | warm_read_only | cacheable | 8.8k | 0.56 | 3.13 | 38 | 17811 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_99_reads | cacheable | 7.8k | 0.62 | 11 | 37 | 15657 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | cacheable_90_reads | cacheable | 5.9k | 0.81 | 29 | 36 | 11899 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | hotspot_reads | cacheable | 6.6k | 0.75 | 26 | 35 | 13236 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_99_reads | total app | 8.8k | 0.59 | 6.05 | 36 | 17673 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ref-roll` | app_90_reads | total app | 6.3k | 0.76 | 21 | 44 | 12606 | 0 | 0.0 | 0 | 0.00% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 1.5k | 1.13 | 3.57 | 3031 | 0 | 0.0 |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 1.1k | 1.62 | 4.56 | 2158 | 0 | 0.0 |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 684 | 2.40 | 8.86 | 1370 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 560 | 3.12 | 9.58 | 1121 | 0 | 0.0 |
| pg-single | `ref-roll` | hotspot_writes | 1.4k | 1.20 | 3.88 | 2788 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | cold_fill | 400 | 400 | 0 | 0.00% | 107 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | warm_read_only | 49341 | 49341 | 0 | 0.00% | 40445 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_99_reads | 29121 | 29047 | 0 | 0.00% | 23866 | 0.00% | 0 | 0 | — | — | 0 | 598 | 91 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | 73794 | 73720 | 0 | 0.00% | 60822 | 0.00% | 0 | 0 | — | — | 0 | 1191 | 91 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_90_reads | 84398 | 84199 | 0 | 0.00% | 68646 | 0.00% | 0 | 0 | — | — | 0 | 993 | 256 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | 132627 | 132428 | 0 | 0.00% | 108454 | 0.00% | 0 | 0 | — | — | 0 | 2364 | 256 | 0 | — |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | 414919 | 414919 | 0 | 0.00% | 414919 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 414919 | 414919 | 0 | 0.00% | 414919 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_1 | 57994 | 57879 | 0 | 0.00% | 46633 | 0.00% | 0 | 0 | — | — | 0 | 2167 | 163 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_3 | 15891 | 15809 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2377 | 278 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 96 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | warm_read_only | 32885 | 32885 | 0 | 0.00% | 26651 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_99_reads | 24683 | 24575 | 0 | 0.00% | 20197 | 0.00% | 0 | 0 | — | — | 0 | 987 | 133 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | 62299 | 62191 | 0 | 0.00% | 51378 | 0.00% | 0 | 0 | — | — | 0 | 1059 | 133 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_90_reads | 72992 | 72570 | 0 | 0.00% | 59779 | 0.00% | 0 | 0 | — | — | 0 | 987 | 554 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | 103442 | 103020 | 0 | 0.00% | 84747 | 0.00% | 0 | 0 | — | — | 0 | 2094 | 554 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | 415242 | 415242 | 0 | 0.00% | 415242 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 415242 | 415242 | 0 | 0.00% | 415242 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_1 | 63863 | 63507 | 0 | 0.00% | 51718 | 0.00% | 0 | 0 | — | — | 0 | 5091 | 490 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_3 | 19138 | 18695 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2802 | 700 | 0 | — |
| pg-single | `l:na:red:asd:str` | cold_fill | 400 | 400 | 0 | 0.00% | 146 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | warm_read_only | 37638 | 37638 | 0 | 0.00% | 36862 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_99_reads | 22676 | 22567 | 0 | 0.00% | 22353 | 0.00% | 0 | 0 | — | — | 0 | 565 | 130 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | 65984 | 65875 | 0 | 0.00% | 65651 | 0.00% | 0 | 0 | — | — | 0 | 1356 | 130 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_90_reads | 74966 | 74716 | 0 | 0.00% | 73514 | 0.00% | 0 | 0 | — | — | 0 | 1075 | 330 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | 113484 | 113234 | 0 | 0.00% | 111945 | 0.00% | 0 | 0 | — | — | 0 | 2472 | 330 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_reads | 44551 | 44551 | 0 | 0.00% | 44551 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 44551 | 44551 | 0 | 0.00% | 44551 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_1 | 31271 | 31151 | 0 | 0.00% | 29845 | 0.00% | 0 | 0 | — | — | 0 | 2646 | 155 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_3 | 35832 | 35702 | 0 | 0.00% | 35104 | 0.00% | 0 | 0 | — | — | 0 | 2526 | 166 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cold_fill | 400 | 400 | 0 | 0.00% | 148 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 28329 | 28329 | 0 | 0.00% | 27572 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 20122 | 19967 | 0 | 0.00% | 19942 | 0.00% | 0 | 0 | — | — | 0 | 820 | 183 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 56637 | 56482 | 0 | 0.00% | 56449 | 0.00% | 0 | 0 | — | — | 0 | 2085 | 183 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 64696 | 64265 | 0 | 0.00% | 64162 | 0.00% | 0 | 0 | — | — | 0 | 1405 | 576 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 95305 | 94874 | 0 | 0.00% | 94771 | 0.00% | 0 | 0 | — | — | 0 | 2534 | 576 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | 43158 | 43158 | 0 | 0.00% | 43158 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 43158 | 43158 | 0 | 0.00% | 43158 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_1 | 38212 | 37965 | 0 | 0.00% | 37204 | 0.00% | 0 | 0 | — | — | 0 | 2341 | 327 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_3 | 42586 | 42297 | 0 | 0.00% | 42235 | 0.00% | 0 | 0 | — | — | 0 | 5679 | 380 | 0 | — |
| pg-single | `ref-roll` | warm_read_only | 22785 | 22785 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | app_99_reads | 18936 | 18905 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 603 | 66 | 0 | — |
| pg-single | `ref-roll` | cacheable_99_reads | 38880 | 38849 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1898 | 66 | 0 | — |
| pg-single | `ref-roll` | app_90_reads | 51229 | 51029 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1118 | 454 | 0 | — |
| pg-single | `ref-roll` | cacheable_90_reads | 66198 | 65998 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2133 | 454 | 0 | — |
| pg-single | `ref-roll` | hotspot_reads | 17131 | 17131 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ref-roll` | hotspot_writes | 17131 | 17131 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 0-75s | 610451 | 6890 | 1.1% | 0 to 0.25 |
| pg-single | `l:na:mem:wth:str` | 0-75s | 578358 | 7585 | 1.3% | 0 to 0.25 |
| pg-single | `o:opt:red:wth:str` | 0-75s | 244940 | 2214 | 0.9% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 39676 | 565 | 0 | — | 565 | 140 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 517µs) |
| pg-single | `l:na:mem:wth:str` | 32015 | 518 | 0 | — | 518 | 118 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 700µs) |
| pg-single | `l:na:red:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:wth:str` | 2987 | 199 | 0 | — | 199 | 23 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.024ms) |
| pg-single | `ref-roll` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |

| Topology | Scenario | readers | keys | elapsed ms | db loads | loads/key | fallbacks | duplicate fills | wrong | impossible |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 16 | 8 | 209.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:str` | 16 | 8 | 210.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str` | 16 | 8 | 231.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | 16 | 8 | 226.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 1 | 8 | memory | 5.3k | 0.0 | 17k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:str` | 3 | 8 | memory | 3.6k | 0.0 | 3.2k | 15891 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:wth:str` | 1 | 8 | memory | 6.6k | 0.0 | 19k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:str` | 3 | 8 | memory | 3.9k | 0.0 | 4.3k | 19138 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:red:asd:str` | 1 | 8 | redis | 2.3k | 0.0 | 11k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str` | 3 | 8 | redis | 2.5k | 0.0 | 12k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 1 | 8 | redis | 3.4k | 0.0 | 12k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 3 | 8 | redis | 3.8k | 0.0 | 14k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |

## Fault injection and controls

A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly
has not been fast.

| Topology | Scenario | Fault | reproduced | wrong observed | impossible observed | correctness |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |

## Ledger replay (audit)

| Topology | Scenario | Phase # | passed | checks | content mismatches | donation-set mismatches | aggregate mismatches | recent-slice mismatches | design checks | first failure |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:wth:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `o:opt:red:wth:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `ref-roll` | 1 | true | 14 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 a_rollup_drift=0 |  |

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | memory | 256.0 KiB | — | 0 | 31251 | 32962 | 32858 | 104 | 0 | 0 | 12000 | 0 | 20 | 0 | 0 | 1 |
| pg-single | `l:na:mem:wth:str` | memory | 256.0 KiB | 235.2 KiB | 111 | 24479 | 29920 | 29751 | 169 | 2 | 0 | 10159 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | redis | 256.0 KiB | — | 0 | 0 | 4463 | 4321 | 142 | 2 | 0 | 6331 | 0 | 20 | 0 | 0 | 0 |
| pg-single | `ref-roll` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Redis accounting (from `INFO`)

Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in
behaviour between the two is partly a policy difference and is stated as such.

| Topology | Scenario | redis_version | maxmemory | maxmemory_policy | maxmemory_samples | used_memory | used_memory_peak | keyspace_hits | keyspace_misses | expired_keys | evicted_keys | total_commands_processed | total_net_input_bytes | total_net_output_bytes | used_cpu_sys | used_cpu_user |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `o:opt:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1834616 | 3980328 | 275480 | 15897 | 1 | 0 | 352875 | 32874638 | 536017852 | 5.509348 | 3.568606 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 334 | 1 | 34 | db-only | 8 |
| pg-single | `l:na:mem:wth:str` | 340 | 11 | 182 | db-only | 8 |
| pg-single | `l:na:red:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:str` | 344 | 168 | 30436 | db-only | 8 |
| pg-single | `ref-roll` | 198 | 139 | 10724 | db-only | 8 |

## Storage

PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a
different query, so none is reported there rather than a number read as bytes on disk.

| Topology | Scenario | Relation | Bytes |
|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:asd:str` | donation | 4.0 MiB |
| pg-single | `l:na:mem:asd:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:mem:asd:str` | person | 168.0 KiB |
| pg-single | `l:na:mem:wth:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:wth:str` | donation | 4.0 MiB |
| pg-single | `l:na:mem:wth:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:mem:wth:str` | person | 168.0 KiB |
| pg-single | `l:na:red:asd:str` | charity | 32.0 KiB |
| pg-single | `l:na:red:asd:str` | donation | 4.0 MiB |
| pg-single | `l:na:red:asd:str` | load_donation | 2.3 MiB |
| pg-single | `l:na:red:asd:str` | person | 168.0 KiB |
| pg-single | `o:opt:red:wth:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:wth:str` | charity | 32.0 KiB |
| pg-single | `o:opt:red:wth:str` | donation | 4.0 MiB |
| pg-single | `o:opt:red:wth:str` | load_donation | 2.3 MiB |
| pg-single | `o:opt:red:wth:str` | person | 168.0 KiB |
| pg-single | `ref-roll` | charity | 32.0 KiB |
| pg-single | `ref-roll` | donation | 4.0 MiB |
| pg-single | `ref-roll` | load_donation | 2.3 MiB |
| pg-single | `ref-roll` | person | 280.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `ref-roll` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |

## Ledger assertion: the requirement must never be ahead of the database

After every writing phase the harness compares, for every donor, the content the database
actually holds against the ledger's freshness requirement. A requirement BEHIND a committed
but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails
the cell, because a bypass read that returned such a state would be an accounting defect
rather than a design finding.

| Topology | Scenario | Phase | donors | mismatches | passed | first example |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:mem:asd:str` | faults | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:mem:wth:str` | faults | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | warm | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | mixed | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | hotspot | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | stampede | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | instances | 800 | 0 | true |  |
| pg-single | `l:na:red:asd:str` | faults | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | warm | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | mixed | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | hotspot | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | stampede | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | instances | 800 | 0 | true |  |
| pg-single | `o:opt:red:wth:str` | faults | 800 | 0 | true |  |
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

