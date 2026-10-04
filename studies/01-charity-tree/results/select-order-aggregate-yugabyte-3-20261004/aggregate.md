# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-order-aggregate-yugabyte-3-20261004`
- run tag: `run/study-01/select-order-aggregate-yugabyte-3-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-order-aggregate-yugabyte-3-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## Execution order

- order seed `42194605`; interleaved rounds: 5
- executed sequence (one sample of every cell per round):
  1. family=topn;stmt=lateral;parent_max=10000
  2. family=topn;stmt=rownumber;parent_max=10000
  3. family=topn;stmt=rownumber;parent_max=1
  4. family=topn;stmt=lateral;parent_max=1000
  5. family=topn;stmt=lateral;parent_max=10
  6. family=aggregate;stmt=post_join;parent_max=10000
  7. family=topn;stmt=rownumber;parent_max=1000
  8. family=topn;stmt=distinct_on;parent_max=1
  9. family=topn;stmt=distinct_on;parent_max=100
  10. family=aggregate;stmt=post_join;parent_max=10
  11. family=topn;stmt=lateral;parent_max=1
  12. family=topn;stmt=distinct_on;parent_max=10000
  13. family=topn;stmt=distinct_on;parent_max=1000
  14. family=aggregate;stmt=correlated;parent_max=1
  15. family=topn;stmt=distinct_on;parent_max=10
  16. family=aggregate;stmt=correlated;parent_max=1000
  17. family=aggregate;stmt=pre_agg;parent_max=1
  18. family=aggregate;stmt=correlated;parent_max=10000
  19. family=topn;stmt=rownumber;parent_max=100
  20. family=topn;stmt=rownumber;parent_max=10
  21. family=aggregate;stmt=pre_agg;parent_max=100
  22. family=aggregate;stmt=post_join;parent_max=1
  23. family=aggregate;stmt=correlated;parent_max=10
  24. family=aggregate;stmt=post_join;parent_max=1000
  25. family=aggregate;stmt=post_join;parent_max=100
  26. family=topn;stmt=lateral;parent_max=100
  27. family=aggregate;stmt=pre_agg;parent_max=10000
  28. family=aggregate;stmt=correlated;parent_max=100
  29. family=aggregate;stmt=pre_agg;parent_max=10
  30. family=aggregate;stmt=pre_agg;parent_max=1000

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 7.739420 | 24.434710 | 8.266088 | 3.1572 | 1.0681 | 0.3383 |
| 10 | 8.885454 | 24.411916 | 17.758949 | 2.7474 | 1.9987 | 0.7275 |
| 100 | 10.484977 | 28.382864 | 124.504364 | 2.7070 | 11.8745 | 4.3866 |
| 1000 | 25.463011 | 25.265151 | 1170.377956 | 0.9922 | 45.9638 | 46.3238 |
| 10000 | 37.412629 | 46.054655 | 19981.692225 | 1.2310 | 534.0895 | 433.8691 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 13.013492 | 7.011562 | 12.146865 | 0.5388 | 0.9334 | 1.7324 |
| 10 | 11.701965 | 14.570725 | 11.742605 | 1.2452 | 1.0035 | 0.8059 |
| 100 | 13.231090 | 80.777657 | 11.744019 | 6.1051 | 0.8876 | 0.1454 |
| 1000 | 16.313609 | 799.642915 | 14.592044 | 49.0169 | 0.8945 | 0.0182 |
| 10000 | 48.069590 | 13886.221819 | 37.909284 | 288.8775 | 0.7886 | 0.0027 |

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
