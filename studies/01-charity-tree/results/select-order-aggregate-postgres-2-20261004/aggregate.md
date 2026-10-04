# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `postgres`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-order-aggregate-postgres-2-20261004`
- run tag: `run/study-01/select-order-aggregate-postgres-2-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-order-aggregate-postgres-2-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## Execution order

- order seed `42168493`; interleaved rounds: 5
- executed sequence (one sample of every cell per round):
  1. family=aggregate;stmt=pre_agg;parent_max=10000
  2. family=aggregate;stmt=post_join;parent_max=10000
  3. family=topn;stmt=lateral;parent_max=100
  4. family=topn;stmt=lateral;parent_max=1
  5. family=topn;stmt=rownumber;parent_max=10000
  6. family=topn;stmt=lateral;parent_max=10
  7. family=aggregate;stmt=pre_agg;parent_max=100
  8. family=aggregate;stmt=correlated;parent_max=100
  9. family=topn;stmt=distinct_on;parent_max=10000
  10. family=aggregate;stmt=post_join;parent_max=10
  11. family=topn;stmt=lateral;parent_max=1000
  12. family=aggregate;stmt=post_join;parent_max=1000
  13. family=topn;stmt=rownumber;parent_max=10
  14. family=topn;stmt=distinct_on;parent_max=10
  15. family=topn;stmt=distinct_on;parent_max=100
  16. family=aggregate;stmt=correlated;parent_max=1000
  17. family=aggregate;stmt=correlated;parent_max=10000
  18. family=aggregate;stmt=post_join;parent_max=1
  19. family=aggregate;stmt=pre_agg;parent_max=1000
  20. family=aggregate;stmt=pre_agg;parent_max=10
  21. family=topn;stmt=distinct_on;parent_max=1000
  22. family=topn;stmt=rownumber;parent_max=1000
  23. family=topn;stmt=rownumber;parent_max=100
  24. family=topn;stmt=distinct_on;parent_max=1
  25. family=aggregate;stmt=correlated;parent_max=10
  26. family=topn;stmt=lateral;parent_max=10000
  27. family=aggregate;stmt=pre_agg;parent_max=1
  28. family=topn;stmt=rownumber;parent_max=1
  29. family=aggregate;stmt=post_join;parent_max=100
  30. family=aggregate;stmt=correlated;parent_max=1

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 0.434762 | 2.483171 | 0.202625 | 5.7116 | 0.4661 | 0.0816 |
| 10 | 0.512776 | 2.220210 | 0.219296 | 4.3298 | 0.4277 | 0.0988 |
| 100 | 0.365654 | 2.388580 | 0.460843 | 6.5324 | 1.2603 | 0.1929 |
| 1000 | 1.607923 | 2.510928 | 1.732768 | 1.5616 | 1.0776 | 0.6901 |
| 10000 | 6.515128 | 4.995993 | 16.241235 | 0.7668 | 2.4928 | 3.2509 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 0.256439 | 0.250637 | 0.195156 | 0.9774 | 0.7610 | 0.7786 |
| 10 | 0.266148 | 0.404105 | 0.247913 | 1.5183 | 0.9315 | 0.6135 |
| 100 | 0.317123 | 0.616913 | 0.238245 | 1.9453 | 0.7513 | 0.3862 |
| 1000 | 0.868835 | 1.385665 | 0.567822 | 1.5949 | 0.6535 | 0.4098 |
| 10000 | 5.747159 | 11.945957 | 3.606980 | 2.0786 | 0.6276 | 0.3019 |

## Controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x05_aggregate_double_count | sel_aggregate_sel | true | 10000 |  |
| x06_topn_rank_ties | sel_topn_sel | true | 10021 |  |

## Plans

| axis | status | plan signature |
|---|---|---|
| family=aggregate;stmt=correlated;parent_max=1 | valid | `39f14eb5f1f73324` |
| family=aggregate;stmt=correlated;parent_max=10 | valid | `39f14eb5f1f73324` |
| family=aggregate;stmt=correlated;parent_max=100 | valid | `39f14eb5f1f73324` |
| family=aggregate;stmt=correlated;parent_max=1000 | valid | `39f14eb5f1f73324` |
| family=aggregate;stmt=correlated;parent_max=10000 | valid | `242506261e61b710` |
| family=aggregate;stmt=post_join;parent_max=1 | valid | `622fe39b15b719fe` |
| family=aggregate;stmt=post_join;parent_max=10 | valid | `622fe39b15b719fe` |
| family=aggregate;stmt=post_join;parent_max=100 | valid | `622fe39b15b719fe` |
| family=aggregate;stmt=post_join;parent_max=1000 | valid | `289dbf6a8a6f2019` |
| family=aggregate;stmt=post_join;parent_max=10000 | valid | `4e6500440fe6ec01` |
| family=aggregate;stmt=pre_agg;parent_max=1 | valid | `2377b93d8eb95272` |
| family=aggregate;stmt=pre_agg;parent_max=10 | valid | `2377b93d8eb95272` |
| family=aggregate;stmt=pre_agg;parent_max=100 | valid | `2377b93d8eb95272` |
| family=aggregate;stmt=pre_agg;parent_max=1000 | valid | `2377b93d8eb95272` |
| family=aggregate;stmt=pre_agg;parent_max=10000 | valid | `4a4d598e44ab1e93` |
| family=topn;stmt=distinct_on;parent_max=1 | valid | `f5d24ec026e7a349` |
| family=topn;stmt=distinct_on;parent_max=10 | valid | `f5d24ec026e7a349` |
| family=topn;stmt=distinct_on;parent_max=100 | valid | `6c7ad82ee706b43e` |
| family=topn;stmt=distinct_on;parent_max=1000 | valid | `b8d4f6586a433b67` |
| family=topn;stmt=distinct_on;parent_max=10000 | valid | `5c980bdf1ec853c3` |
| family=topn;stmt=lateral;parent_max=1 | valid | `621b74ad37190f4a` |
| family=topn;stmt=lateral;parent_max=10 | valid | `621b74ad37190f4a` |
| family=topn;stmt=lateral;parent_max=100 | valid | `2d1ebe491fe72168` |
| family=topn;stmt=lateral;parent_max=1000 | valid | `c3b0330f87551451` |
| family=topn;stmt=lateral;parent_max=10000 | valid | `7f402299b06ee0f8` |
| family=topn;stmt=rownumber;parent_max=1 | valid | `54617db4d8f1c625` |
| family=topn;stmt=rownumber;parent_max=10 | valid | `54617db4d8f1c625` |
| family=topn;stmt=rownumber;parent_max=100 | valid | `6f5d11c9978725f5` |
| family=topn;stmt=rownumber;parent_max=1000 | valid | `55087658b0d8a4ea` |
| family=topn;stmt=rownumber;parent_max=10000 | valid | `2a503fea9af15ad6` |

## Recorded limits

- one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
