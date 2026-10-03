# Study 01 SELECT aggregate & top-N — generated measurements

_Measurements only. Interpretation is in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `postgres`, topology `single`, scale `sweep`, seed 42.

## Provenance

- run id: `select-agg-pg-1b-20261003`
- run tag: `run/study-01/select-agg-pg-1b-20261003`
- environment: `xiles-podman-select-aggregate-20261003`
- repo commit: `08a3c55396e01e07799a1014ee4c99b4b59f964c`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-agg-pg-1b-20261003` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- dataset fingerprint: `ae31e63e6f27a89f`
- client CPU quota: `200000 100000`, memory: `2147483648`
- client throttling: before "nr_throttled 0", after "nr_throttled 0"

## P8 — parent aggregate (per selected parent)

Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.

| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |
|---|---|---|---|---|---|---|
| 1 | 0.234977 | 2.798993 | 0.163634 | 11.9118 | 0.6964 | 0.0585 |
| 10 | 0.224404 | 2.760325 | 0.191366 | 12.3007 | 0.8528 | 0.0693 |
| 100 | 0.415305 | 2.942386 | 0.422182 | 7.0849 | 1.0166 | 0.1435 |
| 1000 | 1.726389 | 3.136784 | 2.203104 | 1.8170 | 1.2761 | 0.7023 |
| 10000 | 7.020799 | 5.763481 | 20.178764 | 0.8209 | 2.8741 | 3.5011 |

## P9 — per-parent top-N at N=1

All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.

| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |
|---|---|---|---|---|---|---|
| 1 | 0.191538 | 0.195276 | 0.124118 | 1.0195 | 0.6480 | 0.6356 |
| 10 | 0.208794 | 0.197841 | 0.169402 | 0.9475 | 0.8113 | 0.8563 |
| 100 | 0.320492 | 0.421836 | 0.259799 | 1.3162 | 0.8106 | 0.6159 |
| 1000 | 0.994342 | 1.693071 | 0.649699 | 1.7027 | 0.6534 | 0.3837 |
| 10000 | 7.745943 | 16.774430 | 5.428776 | 2.1656 | 0.7009 | 0.3236 |

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
