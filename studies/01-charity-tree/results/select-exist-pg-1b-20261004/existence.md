# Study 01 SELECT membership & absence — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 21 cells, 21 valid, engine `postgres`, topology `single`, scale `fanout`, seed 42.

## Provenance

- run id: `select-exist-pg-1b-20261004`
- run tag: `run/study-01/select-exist-pg-1b-20261004`
- environment: `xiles-podman-select-existence-20261004`
- repo commit: `805b5b85828387ee6a74912e815dbbfe56105efa`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-exist-pg-1b-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- client CPU quota: `200000 100000`, memory: `2147483648`

## Fanout levels (constant 10,000 child rows)

| fanout | parents | children | dataset fingerprint |
|---|---|---|---|
| 1 | 10000 | 10000 | `8548e4c065b17e93` |
| 10 | 1000 | 10000 | `b8b1346604686909` |
| 100 | 100 | 10000 | `5f1fcd4d0d1062f8` |

## P6 — membership (parent identity set)

Each row is one fanout level; ms are this run's medians. `in/exists`, `jd/exists` and `jd/in` are the B spelling over A at the same level and predicate.

| fanout | exists ms | in ms | in_notnull ms | join_distinct ms | in/exists | jd/exists | jd/in |
|---|---|---|---|---|---|---|---|
| 1 | 2.195002 | 1.995845 | 2.221116 | 3.202929 | 0.9093 | 1.4592 | 1.6048 |
| 10 | 1.038681 | 1.006198 | 1.047735 | 1.402433 | 0.9687 | 1.3502 | 1.3938 |
| 100 | 0.237410 | 0.268753 | 0.278093 | 1.043884 | 1.1320 | 4.3970 | 3.8842 |

## P7 — absence (parent identity set)

All three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.

| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |
|---|---|---|---|
| 1 | 1.704936 | 1.622063 | 2.515485 |
| 10 | 0.683355 | 0.738915 | 0.792541 |
| 100 | 0.242070 | 0.716452 | 0.724409 |

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
