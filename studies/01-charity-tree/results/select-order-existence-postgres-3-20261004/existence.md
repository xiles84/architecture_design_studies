# Study 01 SELECT membership & absence — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 21 cells, 21 valid, engine `postgres`, topology `single`, scale `fanout`, seed 42.

## Provenance

- run id: `select-order-existence-postgres-3-20261004`
- run tag: `run/study-01/select-order-existence-postgres-3-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-order-existence-postgres-3-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- client CPU quota: `200000 100000`, memory: `2147483648`

## Fanout levels (constant 10,000 child rows)

| fanout | parents | children | dataset fingerprint |
|---|---|---|---|
| 1 | 10000 | 10000 | `8548e4c065b17e93` |
| 10 | 1000 | 10000 | `b8b1346604686909` |
| 100 | 100 | 10000 | `5f1fcd4d0d1062f8` |

## Execution order (per fanout phase)

- fanout 1: order seed `42610932`, 5 interleaved rounds
  1. family=absence;stmt=left_is_null;fanout=1
  2. family=membership;stmt=exists;fanout=1
  3. family=absence;stmt=not_exists;fanout=1
  4. family=membership;stmt=in_notnull;fanout=1
  5. family=absence;stmt=not_in_notnull;fanout=1
  6. family=membership;stmt=join_distinct;fanout=1
  7. family=membership;stmt=in;fanout=1
- fanout 10: order seed `42791344`, 5 interleaved rounds
  1. family=membership;stmt=exists;fanout=10
  2. family=absence;stmt=not_in_notnull;fanout=10
  3. family=membership;stmt=in_notnull;fanout=10
  4. family=membership;stmt=in;fanout=10
  5. family=membership;stmt=join_distinct;fanout=10
  6. family=absence;stmt=not_exists;fanout=10
  7. family=absence;stmt=left_is_null;fanout=10
- fanout 100: order seed `42861284`, 5 interleaved rounds
  1. family=membership;stmt=in_notnull;fanout=100
  2. family=membership;stmt=join_distinct;fanout=100
  3. family=absence;stmt=left_is_null;fanout=100
  4. family=absence;stmt=not_in_notnull;fanout=100
  5. family=absence;stmt=not_exists;fanout=100
  6. family=membership;stmt=exists;fanout=100
  7. family=membership;stmt=in;fanout=100

## P6 — membership (parent identity set)

Each row is one fanout level; ms are this run's medians. `in/exists`, `jd/exists` and `jd/in` are the B spelling over A at the same level and predicate.

| fanout | exists ms | in ms | in_notnull ms | join_distinct ms | in/exists | jd/exists | jd/in |
|---|---|---|---|---|---|---|---|
| 1 | 2.144528 | 2.068100 | 2.608755 | 3.202961 | 0.9644 | 1.4936 | 1.5487 |
| 10 | 0.988961 | 1.086892 | 1.109824 | 1.427715 | 1.0990 | 1.4437 | 1.3136 |
| 100 | 0.268153 | 0.272380 | 0.261529 | 1.133713 | 1.0158 | 4.2279 | 4.1622 |

## P7 — absence (parent identity set)

All three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.

| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |
|---|---|---|---|
| 1 | 1.742235 | 1.534601 | 1.900687 |
| 10 | 0.720443 | 0.757808 | 0.804945 |
| 100 | 0.233455 | 0.643376 | 0.665071 |

## Controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x02_not_in_nullable | sel_absence | true | 0 |  |
| x07_membership_multiplier | sel_membership | true | 7241 |  |
| x02_not_in_nullable | sel_absence | true | 0 |  |
| x07_membership_multiplier | sel_membership | true | 7240 |  |
| x02_not_in_nullable | sel_absence | true | 0 |  |
| x07_membership_multiplier | sel_membership | true | 7240 |  |

## Plans

| axis | status | plan signature |
|---|---|---|
| family=absence;stmt=left_is_null;fanout=1 | valid | `f16e64291eebe84c` |
| family=absence;stmt=left_is_null;fanout=10 | valid | `f16e64291eebe84c` |
| family=absence;stmt=left_is_null;fanout=100 | valid | `f16e64291eebe84c` |
| family=absence;stmt=not_exists;fanout=1 | valid | `ee53c07e94bc3d0e` |
| family=absence;stmt=not_exists;fanout=10 | valid | `ee53c07e94bc3d0e` |
| family=absence;stmt=not_exists;fanout=100 | valid | `9f4785d463ef0450` |
| family=absence;stmt=not_in_notnull;fanout=1 | valid | `576a54b067fba0b6` |
| family=absence;stmt=not_in_notnull;fanout=10 | valid | `576a54b067fba0b6` |
| family=absence;stmt=not_in_notnull;fanout=100 | valid | `576a54b067fba0b6` |
| family=membership;stmt=exists;fanout=1 | valid | `01b008bf56170088` |
| family=membership;stmt=exists;fanout=10 | valid | `b918623bbf93f107` |
| family=membership;stmt=exists;fanout=100 | valid | `64a6e80cc11ff5ba` |
| family=membership;stmt=in;fanout=1 | valid | `01b008bf56170088` |
| family=membership;stmt=in;fanout=10 | valid | `b918623bbf93f107` |
| family=membership;stmt=in;fanout=100 | valid | `64a6e80cc11ff5ba` |
| family=membership;stmt=in_notnull;fanout=1 | valid | `01b008bf56170088` |
| family=membership;stmt=in_notnull;fanout=10 | valid | `b918623bbf93f107` |
| family=membership;stmt=in_notnull;fanout=100 | valid | `64a6e80cc11ff5ba` |
| family=membership;stmt=join_distinct;fanout=1 | valid | `16eb24166018a1dd` |
| family=membership;stmt=join_distinct;fanout=10 | valid | `16eb24166018a1dd` |
| family=membership;stmt=join_distinct;fanout=100 | valid | `16eb24166018a1dd` |

## Recorded limits

- one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- qualifying proportion, parent hit fraction and hot-parent skew are not swept in this bounded run
- YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
