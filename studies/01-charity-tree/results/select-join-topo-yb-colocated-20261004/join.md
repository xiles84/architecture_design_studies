# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 0 cells, 0 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

**Run error:** gate fixture load: unsupported by engine: ERROR: cannot set colocation true on a non-colocated database (SQLSTATE 0A000)

## Provenance

- run id: `select-join-topo-yb-colocated-20261004`
- run tag: ``
- environment: ``
- repo commit: ``
- engine version: ``
- database engine image: `unknown` @ `unknown`
- benchmark client image: `unknown` @ `unknown`
- inputs digest: `3767fc07360b21b6b2a5ecc092f2ad077b020375e039ea9dea711f56abf46768`
- dataset fingerprint: ``
- client CPU quota: ``, memory: ``

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 2000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 2000 (1%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 180000 (90%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | missing | | | | |
| default | 1 | 100% | missing | | | | |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | missing | | | | |
| default | 4 | 100% | missing | | | | |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | missing | | | | |
| default | 16 | 100% | missing | | | | |
| mat | 1 | 1% | missing | | | | |
| mat | 1 | 100% | missing | | | | |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | missing | | | | |
| mat | 4 | 100% | missing | | | | |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | missing | | | | |
| mat | 16 | 100% | missing | | | | |
| notmat | 1 | 1% | missing | | | | |
| notmat | 1 | 100% | missing | | | | |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | missing | | | | |
| notmat | 4 | 100% | missing | | | | |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | missing | | | | |
| notmat | 16 | 100% | missing | | | | |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|

## Matched blended workload

No writer record: this run predates the phase field or is feasibility-only.

## Execution order

- order seed `0`, 0 interleaved rounds

## Live negative controls

None recorded.

## Excluded cells

None: every cell is valid.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
