# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `postgres`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-agg-pg-2-20261003`
- run tag: `run/study-01/select-agg-pg-2-20261003`
- environment: `xiles-podman-select-aggregate-20261003`
- repo commit: `33b769017da9648bef92cb1a804697e42c957e90`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-agg-pg-2-20261003` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 0.291771 | 2.765842 | 0.240582 | 9.4795 | 0.8246 | 0.0870 |
| 10 | 0.320708 | 2.876002 | 0.253065 | 8.9677 | 0.7891 | 0.0880 |
| 100 | 0.460999 | 2.930769 | 0.502608 | 6.3574 | 1.0903 | 0.1715 |
| 1000 | 1.903942 | 3.216808 | 2.315719 | 1.6896 | 1.2163 | 0.7199 |
| 10000 | 7.532297 | 6.861774 | 22.264445 | 0.9110 | 2.9559 | 3.2447 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 0.223831 | 0.234856 | 0.216995 | 1.0493 | 0.9695 | 0.9239 |
| 10 | 0.205327 | 0.215903 | 0.226328 | 1.0515 | 1.1023 | 1.0483 |
| 100 | 0.344501 | 0.437996 | 0.277735 | 1.2714 | 0.8062 | 0.6341 |
| 1000 | 1.017237 | 1.816070 | 0.721542 | 1.7853 | 0.7093 | 0.3973 |
| 10000 | 7.999212 | 16.629260 | 5.040590 | 2.0789 | 0.6301 | 0.3031 |

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
