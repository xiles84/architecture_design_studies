# Study 05 — external cache throughput and consistency — run `ackcheck`

Generated 2026-09-21T16:31:55Z from 3 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `ackcheck` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `f9d9f57894de7e50094d75beb3ac6887841b8577` |
| `git describe` | `study-05/v2.2-residual-classified-dirty` |
| Working tree dirty at start | true |
| Benchmark image | `localhost/cachebench:1` (`a895003effe0cd8d0dd`) |
| Resource framing | `db-only` |
| **Inputs digest** | `94c01ccf3da86706` |
| Cells | 3 |
| Distinct scenarios | 3 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 3; cells that failed: 2.
  - failed: `pg-single/legacy-na-memory-aside-strict-coord` (phase ackcheck: ACK CHECK FAILED: 8 of 310 acknowledgements were not visible to a second session immediately afterwards (acknowledged donor 2 and a SECOND sessi…)
  - failed: `pg-single/ref-normalized-indexed` (phase ackcheck: ACK CHECK FAILED: 3 of 310 acknowledgements were not visible to a second session immediately afterwards (acknowledged donor 3 and a SECOND sessi…)
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No strict cell recorded a stale-after-ack read.
- No cell reported an impossible / uncommitted cache value.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `94c01ccf3da86706`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `legacy-na-memory-aside-strict-coord` | legacy | na | memory | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `ref-normalized-indexed` | ref |  | none |  |  |  | invalidation-only | aborted | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

*No cache hits were observed in this run, so there is no expiry evidence to report.*

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:mem:wth:str` | 20 | 0 | 0 | — | 0 | 0 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 621µs) |
| pg-single | `ref-norm` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |

*The stampede phase did not run for any cell in this run.*

## Fault injection and controls

*No fault phase ran in this run.*

## Ledger replay (audit)

*No audit phase ran in this run.*

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed | write no-ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:str` | memory | 256.0 KiB | 239.7 KiB | 116 | 0 | 345 | 336 | 9 | 0 | 0 | 650 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-norm` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:wth:str` | 8 | 0 | 0 | db-only | 8 |
| pg-single | `ref-norm` | 0 | 0 | 0 | db-only | 8 |

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
| pg-single | `ref-norm` | charity | 32.0 KiB |
| pg-single | `ref-norm` | donation | 3.0 MiB |
| pg-single | `ref-norm` | load_donation | 2.3 MiB |
| pg-single | `ref-norm` | person | 168.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `l:na:mem:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `ref-norm` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |

## Ledger assertion: the requirement must never be ahead of the database

After every writing phase the harness compares, for every donor, the content the database
actually holds against the ledger's freshness requirement. A requirement BEHIND a committed
but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails
the cell, because a bypass read that returned such a state would be an accounting defect
rather than a design finding.

*No ledger assertion ran in this run.*

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

