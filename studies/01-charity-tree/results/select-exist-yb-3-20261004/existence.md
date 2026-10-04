# Study 01 SELECT membership & absence — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 21 cells, 21 valid, engine `yugabyte`, topology `single`, scale `fanout`, seed 42.

## Provenance

- run id: `select-exist-yb-3-20261004`
- run tag: `run/study-01/select-exist-yb-3-20261004`
- environment: `xiles-podman-select-existence-20261004`
- repo commit: `805b5b85828387ee6a74912e815dbbfe56105efa`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-exist-yb-3-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
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
| 1 | 8.685197 | 8.945586 | 8.685973 | 9.179738 | 1.0300 | 1.0569 | 1.0262 |
| 10 | 4.167820 | 4.332868 | 4.294224 | 3.920944 | 1.0396 | 0.9408 | 0.9049 |
| 100 | 3.515803 | 3.570068 | 3.919756 | 3.588145 | 1.0154 | 1.0206 | 1.0051 |

## P7 — absence (parent identity set)

All three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.

| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |
|---|---|---|---|
| 1 | 8.002279 | 7.842389 | 8.273174 |
| 10 | 3.723942 | 3.941372 | 4.064619 |
| 100 | 3.632083 | 3.655629 | 3.867141 |

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
| family=absence;stmt=left_is_null;fanout=1 | valid | `68323d5ddc34b966` |
| family=absence;stmt=left_is_null;fanout=10 | valid | `c46cbe8221360b15` |
| family=absence;stmt=left_is_null;fanout=100 | valid | `c46cbe8221360b15` |
| family=absence;stmt=not_exists;fanout=1 | valid | `4befbd47cba4cb7d` |
| family=absence;stmt=not_exists;fanout=10 | valid | `4befbd47cba4cb7d` |
| family=absence;stmt=not_exists;fanout=100 | valid | `4befbd47cba4cb7d` |
| family=absence;stmt=not_in_notnull;fanout=1 | valid | `aee973d7d8e0c259` |
| family=absence;stmt=not_in_notnull;fanout=10 | valid | `aee973d7d8e0c259` |
| family=absence;stmt=not_in_notnull;fanout=100 | valid | `aee973d7d8e0c259` |
| family=membership;stmt=exists;fanout=1 | valid | `7829bec10abd7f80` |
| family=membership;stmt=exists;fanout=10 | valid | `ae8547df9c874293` |
| family=membership;stmt=exists;fanout=100 | valid | `78118f4ada896e85` |
| family=membership;stmt=in;fanout=1 | valid | `7829bec10abd7f80` |
| family=membership;stmt=in;fanout=10 | valid | `ae8547df9c874293` |
| family=membership;stmt=in;fanout=100 | valid | `78118f4ada896e85` |
| family=membership;stmt=in_notnull;fanout=1 | valid | `83a17ba911b18be2` |
| family=membership;stmt=in_notnull;fanout=10 | valid | `aff0d67ec775c95a` |
| family=membership;stmt=in_notnull;fanout=100 | valid | `969ccb966a5b07af` |
| family=membership;stmt=join_distinct;fanout=1 | valid | `a5173b86da1d211f` |
| family=membership;stmt=join_distinct;fanout=10 | valid | `09764a71fb76da6d` |
| family=membership;stmt=join_distinct;fanout=100 | valid | `09764a71fb76da6d` |

## Recorded limits

- one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- qualifying proportion, parent hit fraction and hot-parent skew are not swept in this bounded run
- YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
