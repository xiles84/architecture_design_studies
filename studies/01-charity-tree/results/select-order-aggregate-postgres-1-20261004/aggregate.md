# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 31 cells, 31 valid, engine `postgres`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-order-aggregate-postgres-1-20261004`
- run tag: `run/study-01/select-order-aggregate-postgres-1-20261004`
- environment: `xiles-podman-select-order-correction-20261004`
- repo commit: `ca2a38fbf9e53da00344578258138db9e01d3545`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-order-aggregate-postgres-1-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## Execution order

- order seed `42877742`; interleaved rounds: 5
- executed sequence (one sample of every cell per round):
  1. family=topn;stmt=distinct_on;parent_max=1000
  2. family=topn;stmt=distinct_on;parent_max=10000
  3. family=aggregate;stmt=post_join;parent_max=10
  4. family=topn;stmt=rownumber;parent_max=1
  5. family=topn;stmt=rownumber;parent_max=10
  6. family=aggregate;stmt=pre_agg;parent_max=1000
  7. family=topn;stmt=lateral;parent_max=10
  8. family=aggregate;stmt=pre_agg;parent_max=10000
  9. family=topn;stmt=distinct_on;parent_max=10
  10. family=topn;stmt=lateral;parent_max=10000
  11. family=aggregate;stmt=correlated;parent_max=100
  12. family=aggregate;stmt=pre_agg;parent_max=100
  13. family=aggregate;stmt=post_join;parent_max=100
  14. family=aggregate;stmt=correlated;parent_max=10
  15. family=aggregate;stmt=pre_agg;parent_max=10
  16. family=aggregate;stmt=pre_agg;parent_max=1
  17. family=topn;stmt=rownumber;parent_max=10000
  18. family=topn;stmt=distinct_on;parent_max=100
  19. family=topn;stmt=rownumber;parent_max=100
  20. family=aggregate;stmt=post_join;parent_max=1
  21. family=aggregate;stmt=post_join;parent_max=1000
  22. family=topn;stmt=lateral;parent_max=1
  23. family=aggregate;stmt=correlated;parent_max=1000
  24. family=aggregate;stmt=correlated;parent_max=1
  25. family=topn;stmt=rownumber;parent_max=1000
  26. family=aggregate;stmt=post_join;parent_max=10000
  27. family=topn;stmt=lateral;parent_max=100
  28. family=topn;stmt=lateral;parent_max=1000
  29. family=topn;stmt=distinct_on;parent_max=1
  30. family=aggregate;stmt=correlated;parent_max=10000
  31. family=topn;stmt=distinct_on;parent_max=1000;position=last

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 0.283191 | 2.230215 | 0.254698 | 7.8753 | 0.8994 | 0.1142 |
| 10 | 0.499808 | 2.295233 | 0.254588 | 4.5922 | 0.5094 | 0.1109 |
| 100 | 0.406413 | 2.325254 | 0.700587 | 5.7214 | 1.7238 | 0.3013 |
| 1000 | 1.365769 | 2.587679 | 1.727673 | 1.8947 | 1.2650 | 0.6677 |
| 10000 | 5.995085 | 4.751486 | 15.907986 | 0.7926 | 2.6535 | 3.3480 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 0.240248 | 0.228814 | 0.203181 | 0.9524 | 0.8457 | 0.8880 |
| 10 | 0.179630 | 0.260120 | 0.381249 | 1.4481 | 2.1224 | 1.4657 |
| 100 | 0.350775 | 0.608117 | 0.423064 | 1.7336 | 1.2061 | 0.6957 |
| 1000 | 0.803642 | 1.365204 | 0.737363 | 1.6988 | 0.9175 | 0.5401 |
| 10000 | 6.238146 | 12.350192 | 3.806609 | 1.9798 | 0.6102 | 0.3082 |

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
| family=topn;stmt=distinct_on;parent_max=1000;position=last | valid | `b8d4f6586a433b67` |
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
