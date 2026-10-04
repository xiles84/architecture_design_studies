# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-order-aggregate-yugabyte-2-20261004`
- run tag: `run/study-01/select-order-aggregate-yugabyte-2-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-order-aggregate-yugabyte-2-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## Execution order

- order seed `42725796`; interleaved rounds: 5
- executed sequence (one sample of every cell per round):
  1. family=aggregate;stmt=post_join;parent_max=1
  2. family=aggregate;stmt=correlated;parent_max=100
  3. family=aggregate;stmt=pre_agg;parent_max=100
  4. family=aggregate;stmt=correlated;parent_max=10000
  5. family=aggregate;stmt=post_join;parent_max=100
  6. family=aggregate;stmt=correlated;parent_max=10
  7. family=topn;stmt=lateral;parent_max=10000
  8. family=topn;stmt=lateral;parent_max=1
  9. family=topn;stmt=distinct_on;parent_max=1000
  10. family=topn;stmt=distinct_on;parent_max=1
  11. family=topn;stmt=rownumber;parent_max=10
  12. family=topn;stmt=rownumber;parent_max=1000
  13. family=aggregate;stmt=post_join;parent_max=1000
  14. family=aggregate;stmt=post_join;parent_max=10
  15. family=aggregate;stmt=pre_agg;parent_max=10
  16. family=aggregate;stmt=pre_agg;parent_max=10000
  17. family=topn;stmt=rownumber;parent_max=1
  18. family=topn;stmt=lateral;parent_max=1000
  19. family=topn;stmt=lateral;parent_max=10
  20. family=topn;stmt=distinct_on;parent_max=100
  21. family=topn;stmt=rownumber;parent_max=100
  22. family=aggregate;stmt=correlated;parent_max=1
  23. family=aggregate;stmt=post_join;parent_max=10000
  24. family=topn;stmt=distinct_on;parent_max=10000
  25. family=aggregate;stmt=pre_agg;parent_max=1000
  26. family=topn;stmt=rownumber;parent_max=10000
  27. family=topn;stmt=lateral;parent_max=100
  28. family=aggregate;stmt=pre_agg;parent_max=1
  29. family=topn;stmt=distinct_on;parent_max=10
  30. family=aggregate;stmt=correlated;parent_max=1000

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 3.316064 | 9.607655 | 4.085893 | 2.8973 | 1.2322 | 0.4253 |
| 10 | 3.684170 | 11.998627 | 7.243868 | 3.2568 | 1.9662 | 0.6037 |
| 100 | 5.063425 | 9.343454 | 43.470839 | 1.8453 | 8.5853 | 4.6525 |
| 1000 | 9.475192 | 13.310056 | 416.727342 | 1.4047 | 43.9809 | 31.3092 |
| 10000 | 16.910202 | 18.925754 | 6726.910867 | 1.1192 | 397.8019 | 355.4369 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 5.823031 | 4.072752 | 5.008627 | 0.6994 | 0.8601 | 1.2298 |
| 10 | 4.954549 | 5.742642 | 5.483670 | 1.1591 | 1.1068 | 0.9549 |
| 100 | 4.560806 | 32.109549 | 5.675520 | 7.0403 | 1.2444 | 0.1768 |
| 1000 | 6.282085 | 276.741397 | 5.835957 | 44.0525 | 0.9290 | 0.0211 |
| 10000 | 18.108292 | 4660.097965 | 18.050407 | 257.3461 | 0.9968 | 0.0039 |

## Controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x05_aggregate_double_count | sel_aggregate_sel | true | 10000 |  |
| x06_topn_rank_ties | sel_topn_sel | true | 10021 |  |

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
