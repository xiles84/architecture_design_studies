# Study 01 SELECT membership & absence — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 21 cells, 21 valid, engine `yugabyte`, topology `single`, scale `fanout`, seed 42.

## Provenance

- run id: `select-order-existence-yugabyte-2-20261004`
- run tag: `run/study-01/select-order-existence-yugabyte-2-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-order-existence-yugabyte-2-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- client CPU quota: `200000 100000`, memory: `2147483648`

## Fanout levels (constant 10,000 child rows)

| fanout | parents | children | dataset fingerprint |
|---|---|---|---|
| 1 | 10000 | 10000 | `8548e4c065b17e93` |
| 10 | 1000 | 10000 | `b8b1346604686909` |
| 100 | 100 | 10000 | `5f1fcd4d0d1062f8` |

## Execution order (per fanout phase)

- fanout 1: order seed `42512324`, 5 interleaved rounds
  1. family=membership;stmt=exists;fanout=1
  2. family=absence;stmt=not_in_notnull;fanout=1
  3. family=membership;stmt=in_notnull;fanout=1
  4. family=absence;stmt=left_is_null;fanout=1
  5. family=absence;stmt=not_exists;fanout=1
  6. family=membership;stmt=in;fanout=1
  7. family=membership;stmt=join_distinct;fanout=1
- fanout 10: order seed `42754368`, 5 interleaved rounds
  1. family=absence;stmt=not_exists;fanout=10
  2. family=membership;stmt=join_distinct;fanout=10
  3. family=membership;stmt=exists;fanout=10
  4. family=membership;stmt=in;fanout=10
  5. family=absence;stmt=left_is_null;fanout=10
  6. family=absence;stmt=not_in_notnull;fanout=10
  7. family=membership;stmt=in_notnull;fanout=10
- fanout 100: order seed `42115572`, 5 interleaved rounds
  1. family=membership;stmt=in_notnull;fanout=100
  2. family=absence;stmt=left_is_null;fanout=100
  3. family=membership;stmt=in;fanout=100
  4. family=absence;stmt=not_exists;fanout=100
  5. family=absence;stmt=not_in_notnull;fanout=100
  6. family=membership;stmt=join_distinct;fanout=100
  7. family=membership;stmt=exists;fanout=100

## P6 — membership (parent identity set)

Each row is one fanout level; ms are this run's medians. `in/exists`, `jd/exists` and `jd/in` are the B spelling over A at the same level and predicate.

| fanout | exists ms | in ms | in_notnull ms | join_distinct ms | in/exists | jd/exists | jd/in |
|---|---|---|---|---|---|---|---|
| 1 | 11.933313 | 12.625299 | 12.851232 | 13.133680 | 1.0580 | 1.1006 | 1.0403 |
| 10 | 6.301696 | 6.512661 | 6.766019 | 5.993135 | 1.0335 | 0.9510 | 0.9202 |
| 100 | 5.003641 | 4.863765 | 5.493082 | 5.071254 | 0.9720 | 1.0135 | 1.0427 |

## P7 — absence (parent identity set)

All three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.

| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |
|---|---|---|---|
| 1 | 11.226823 | 11.292946 | 11.192847 |
| 10 | 5.790081 | 6.125948 | 6.061813 |
| 100 | 4.979181 | 5.036845 | 5.407629 |

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
