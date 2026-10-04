# Study 01 SELECT membership & absence — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 21 cells, 21 valid, engine `yugabyte`, topology `single`, scale `fanout`, seed 42.

## Provenance

- run id: `select-exist-yb-1-20261004`
- run tag: `run/study-01/select-exist-yb-1-20261004`
- environment: `xiles-podman-select-existence-20261004`
- repo commit: `805b5b85828387ee6a74912e815dbbfe56105efa`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-exist-yb-1-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
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
| 1 | 9.020588 | 8.916762 | 9.478730 | 9.362084 | 0.9885 | 1.0379 | 1.0499 |
| 10 | 3.963112 | 4.038378 | 4.416506 | 3.783003 | 1.0190 | 0.9546 | 0.9368 |
| 100 | 3.682824 | 3.509264 | 3.773528 | 3.607218 | 0.9529 | 0.9795 | 1.0279 |

## P7 — absence (parent identity set)

All three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.

| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |
|---|---|---|---|
| 1 | 8.304120 | 8.584150 | 8.905806 |
| 10 | 3.424493 | 3.589667 | 3.776185 |
| 100 | 3.683749 | 3.520505 | 3.919785 |

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
