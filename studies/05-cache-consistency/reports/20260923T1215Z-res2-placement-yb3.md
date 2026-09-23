# Study 05 — external cache throughput and consistency — run `20260923T1215Z-res2-placement-yb3`

Generated 2026-09-23T12:02:06Z from 2 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `20260923T1215Z-res2-placement-yb3` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `7ca602457d5336dde8f481093b1637a21f11a6b4` |
| `git describe` | `run/05-cache-consistency/20260923T1157Z-res2-yb3-1-g7ca6024` |
| Working tree dirty at start | false |
| Run tag | `run/05-cache-consistency/20260923T1215Z-res2-placement-yb3` |
| Benchmark image | `localhost/cachebench:1` (`86aaa3803d9358d2aba`) |
| Resource framing | `db-only` |
| **Inputs digest** | `a18817d592f767c5` |
| Cells | 2 |
| Distinct scenarios | 2 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 3 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 2; cells that failed: 0.
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No strict cell recorded a stale-after-ack read.
- No cell reported an impossible / uncommitted cache value.
- Highest and lowest measured throughput across all cells: 882 (`yb-cluster3/y1-coloc/warm_read_only`) and 520 (`yb-cluster3/y2-noncol/cacheable_99_reads`).

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `a18817d592f767c5`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| yb-cluster3 | `ref-y1-colocated` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |
| yb-cluster3 | `ref-y2-noncolocated` | ref |  | none |  |  |  | invalidation-only | pass | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | warm_read_only | cacheable | 882 | 7.52 | 36 | 786 | 4820 | 0 | 40.0 | 0 | 0.00% |
| yb-cluster3 | `y1-coloc` | cacheable_99_reads | cacheable | 852 | 7.36 | 37 | 48 | 5130 | 0 | 6.7 | 0 | 0.00% |
| yb-cluster3 | `y1-coloc` | cacheable_90_reads | cacheable | 801 | 7.61 | 37 | 55 | 4800 | 0 | 9.5 | 0 | 0.00% |
| yb-cluster3 | `y1-coloc` | app_99_reads | total app | 444 | 10 | 68 | 1006 | 2898 | 1 | 41.4 | 0 | 0.00% |
| yb-cluster3 | `y1-coloc` | app_90_reads | total app | 558 | 11 | 53 | 136 | 3349 | 12 | 18.6 | 0 | 0.00% |
| yb-cluster3 | `y2-noncol` | warm_read_only | cacheable | 673 | 8.64 | 43 | 490 | 3908 | 0 | 22.2 | 0 | 0.00% |
| yb-cluster3 | `y2-noncol` | cacheable_99_reads | cacheable | 520 | 12 | 45 | 57 | 3149 | 0 | 9.2 | 0 | 0.00% |
| yb-cluster3 | `y2-noncol` | cacheable_90_reads | cacheable | 557 | 11 | 46 | 61 | 3429 | 0 | 22.3 | 0 | 0.00% |
| yb-cluster3 | `y2-noncol` | app_99_reads | total app | 387 | 14 | 63 | 669 | 2479 | 1 | 27.2 | 0 | 0.00% |
| yb-cluster3 | `y2-noncol` | app_90_reads | total app | 309 | 20 | 66 | 110 | 2155 | 11 | 50.0 | 0 | 0.00% |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | warm_read_only | 4894 | 4894 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| yb-cluster3 | `y1-coloc` | app_99_reads | 2622 | 2609 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2321 | 16 | 0 | — |
| yb-cluster3 | `y1-coloc` | cacheable_99_reads | 8161 | 8148 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 5734 | 16 | 0 | — |
| yb-cluster3 | `y1-coloc` | app_90_reads | 10832 | 10747 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 3563 | 106 | 0 | — |
| yb-cluster3 | `y1-coloc` | cacheable_90_reads | 16083 | 15998 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 5603 | 106 | 0 | — |
| yb-cluster3 | `y2-noncol` | warm_read_only | 3981 | 3981 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| yb-cluster3 | `y2-noncol` | app_99_reads | 2269 | 2268 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 1069 | 3 | 0 | — |
| yb-cluster3 | `y2-noncol` | cacheable_99_reads | 5652 | 5651 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2926 | 3 | 0 | — |
| yb-cluster3 | `y2-noncol` | app_90_reads | 7448 | 7396 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2926 | 75 | 0 | — |
| yb-cluster3 | `y2-noncol` | cacheable_90_reads | 11225 | 11173 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 6478 | 75 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

*No cache hits were observed in this run, so there is no expiry evidence to report.*

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |
| yb-cluster3 | `y2-noncol` | 0 | 0 | 0 | — | 0 | 0 | 0 | n/a | no cache in this scenario |

*The stampede phase did not run for any cell in this run.*

## Fault injection and controls

*No fault phase ran in this run.*

## Ledger replay (audit)

*No audit phase ran in this run.*

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| yb-cluster3 | `y2-noncol` | none | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | 415 | 0 | 0 | db-only | 8 |
| yb-cluster3 | `y2-noncol` | 434 | 0 | 0 | db-only | 8 |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| yb-cluster3 | `y2-noncol` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |

## Ledger assertion: the requirement must never be ahead of the database

After every writing phase the harness compares, for every donor, the content the database
actually holds against the ledger's freshness requirement. A requirement BEHIND a committed
but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails
the cell, because a bypass read that returned such a state would be an accounting defect
rather than a design finding.

| Topology | Scenario | Phase | donors | mismatches | passed | first example |
|---|---|---|---|---|---|---|
| yb-cluster3 | `y1-coloc` | warm | 800 | 0 | true |  |
| yb-cluster3 | `y1-coloc` | mixed | 800 | 0 | true |  |
| yb-cluster3 | `y2-noncol` | warm | 800 | 0 | true |  |
| yb-cluster3 | `y2-noncol` | mixed | 800 | 0 | true |  |

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

