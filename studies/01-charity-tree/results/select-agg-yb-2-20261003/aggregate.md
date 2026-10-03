# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 31 cells, 31 valid, engine `yugabyte`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-agg-yb-2-20261003`
- run tag: `run/study-01/select-agg-yb-2-20261003`
- environment: `xiles-podman-select-aggregate-20261003`
- repo commit: `33b769017da9648bef92cb1a804697e42c957e90`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-agg-yb-2-20261003` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 3.866297 | 12.130249 | 3.993887 | 3.1374 | 1.0330 | 0.3293 |
| 10 | 3.952024 | 12.482785 | 8.229113 | 3.1586 | 2.0823 | 0.6592 |
| 100 | 5.474165 | 12.297712 | 56.826322 | 2.2465 | 10.3808 | 4.6209 |
| 1000 | 9.926294 | 12.648519 | 519.449263 | 1.2742 | 52.3306 | 41.0680 |
| 10000 | 17.183780 | 22.390241 | 8716.346747 | 1.3030 | 507.2427 | 389.2922 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 6.229836 | 3.538668 | 5.690649 | 0.5680 | 0.9135 | 1.6081 |
| 10 | 5.740875 | 6.795568 | 6.002842 | 1.1837 | 1.0456 | 0.8833 |
| 100 | 6.108811 | 38.662160 | 5.408136 | 6.3289 | 0.8853 | 0.1399 |
| 1000 | 7.351668 | 335.158272 | 6.296445 | 45.5894 | 0.8565 | 0.0188 |
| 10000 | 20.694888 | 5679.717299 | 17.214960 | 274.4503 | 0.8318 | 0.0030 |

## Controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x05_aggregate_double_count | sel_aggregate_sel | true | 10000 |  |
| x06_topn_rank_ties | sel_topn_sel | true | 10021 |  |

## Repeated-design position control

The first matrix cell repeated at the end of the phase; a large last/first ratio is a session position effect, not a spelling effect.

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 3.866297 | 3.677096 | 0.9511 | true |

## Plans

| axis | status | plan signature |
|---|---|---|
| family=aggregate;stmt=correlated;parent_max=1 | valid | `d729f8dfabb20429` |
| family=aggregate;stmt=correlated;parent_max=10 | valid | `e2148140c3d78dc2` |
| family=aggregate;stmt=correlated;parent_max=100 | valid | `64e0cf364832f8cc` |
| family=aggregate;stmt=correlated;parent_max=1000 | valid | `889741b382830e22` |
| family=aggregate;stmt=correlated;parent_max=10000 | valid | `f00fd7d86260bf76` |
| family=aggregate;stmt=post_join;parent_max=1 | valid | `d87f09989d187150` |
| family=aggregate;stmt=post_join;parent_max=10 | valid | `92dec2d5c59f8b30` |
| family=aggregate;stmt=post_join;parent_max=100 | valid | `1744c0bb5c600761` |
| family=aggregate;stmt=post_join;parent_max=1000 | valid | `e1ae5e54743ba5ba` |
| family=aggregate;stmt=post_join;parent_max=10000 | valid | `7a8b4ce432ee2d17` |
| family=aggregate;stmt=post_join;parent_max=1;position=last | valid | `d87f09989d187150` |
| family=aggregate;stmt=pre_agg;parent_max=1 | valid | `fa60d810f18daf06` |
| family=aggregate;stmt=pre_agg;parent_max=10 | valid | `abe3eb36aa0dd96e` |
| family=aggregate;stmt=pre_agg;parent_max=100 | valid | `f36bc7cf0492ab4a` |
| family=aggregate;stmt=pre_agg;parent_max=1000 | valid | `90b3df74c2a035ec` |
| family=aggregate;stmt=pre_agg;parent_max=10000 | valid | `30475639f838149f` |
| family=topn;stmt=distinct_on;parent_max=1 | valid | `4d1cdb55f23e53bd` |
| family=topn;stmt=distinct_on;parent_max=10 | valid | `8aa69d04f15c2fd1` |
| family=topn;stmt=distinct_on;parent_max=100 | valid | `2df9cee787e82e86` |
| family=topn;stmt=distinct_on;parent_max=1000 | valid | `7df9d49082f18bcc` |
| family=topn;stmt=distinct_on;parent_max=10000 | valid | `0737e39fc0546857` |
| family=topn;stmt=lateral;parent_max=1 | valid | `5ddbfa1b3131d1e6` |
| family=topn;stmt=lateral;parent_max=10 | valid | `4a044949a7be2b0c` |
| family=topn;stmt=lateral;parent_max=100 | valid | `9c313abcdb61169e` |
| family=topn;stmt=lateral;parent_max=1000 | valid | `0d6c1fbb8a85adfc` |
| family=topn;stmt=lateral;parent_max=10000 | valid | `11b29ad565240cdd` |
| family=topn;stmt=rownumber;parent_max=1 | valid | `8bd5802160b47a69` |
| family=topn;stmt=rownumber;parent_max=10 | valid | `eb7b94bad446123f` |
| family=topn;stmt=rownumber;parent_max=100 | valid | `d8dc6b4a6c4f5c52` |
| family=topn;stmt=rownumber;parent_max=1000 | valid | `396542271f535da3` |
| family=topn;stmt=rownumber;parent_max=10000 | valid | `dbf23f9de1481e68` |

## Recorded limits

- one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
