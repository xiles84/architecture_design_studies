# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-i3-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `4d80c4a5b4cee2cb40c3189670cb317f74e28f78`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-i3-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `d2ffaa1ab94d88e6fd02f5f42352d9fb8a3eecbb7fc1f2e94e460c52a8e3f980`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_arm`: yb-range-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 2 CPU 3g
- setting `phase`: isolated
- setting `placement`: engine default
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-single | 933 | 1 |

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 200 2000 20000 100000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 200 2000 20000 100000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 100000 | 100000 |
| 180000 | 180000 |
| 20 | 20 |
| 200 | 200 |
| 2000 | 2000 |
| 20000 | 20000 |

| statement | support | rank index observed | forced policy evidence |
|---|---|---|---|
| j01_join_where | valid | true |  |
| j02_join_on | valid | true |  |
| j03_filtered_derived | valid | true |  |
| j04_cte_filtered | valid | true | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | true | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | true | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | true |  |
| j11_optional_on | valid | true |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 2.1886 | 1.8375 | 2.2142 | 1.9016 | 2.1760 | 2.2279 | 2.2924 | 82.2052 | 1.7995 | 1.8254 |
| 200 (0.1%) | 3.4524 | 3.2793 | 3.6249 | 2.7578 | 3.5286 | 3.7494 | 3.7142 | 81.4112 | 3.8035 | 3.1600 |
| 2000 (1%) | 9.2309 | 8.9381 | 8.6844 | 9.0044 | 10.5318 | 9.4721 | 9.1821 | 81.7059 | 8.1871 | 9.0737 |
| 20000 (10%) | 61.5973 | 72.8096 | 62.6936 | 64.0440 | 66.3889 | 63.3562 | 63.4425 | 85.8457 | 63.8924 | 62.3997 |
| 100000 (50%) | 86.4668 | 88.1216 | 86.1025 | 97.2561 | 90.1374 | 85.0512 | 84.6604 | 110.3123 | 84.5678 | 86.9163 |
| 180000 (90%) | 104.4806 | 111.3741 | 107.2414 | 106.3247 | 110.0455 | 106.2137 | 105.4371 | 104.9121 | 101.9260 | 103.1727 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1886 | 1.8375 | 0.8396 | false |
| 20 | filtered_derived vs join_where | 2.1886 | 2.2142 | 1.0117 | false |
| 20 | cte_filtered vs join_where | 2.1886 | 1.9016 | 0.8689 | false |
| 20 | cte_filtered_mat vs join_where | 2.1886 | 2.1760 | 0.9943 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1886 | 2.2279 | 1.0180 | false |
| 20 | cte_joined_outer vs join_where | 2.1886 | 2.2924 | 1.0474 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1886 | 82.2052 | 37.5613 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1886 | 1.7995 | 0.8223 | false |
| 20 | comma_join vs join_where | 2.1886 | 1.8254 | 0.8341 | false |
| 200 | join_on vs join_where | 3.4524 | 3.2793 | 0.9499 | false |
| 200 | filtered_derived vs join_where | 3.4524 | 3.6249 | 1.0500 | false |
| 200 | cte_filtered vs join_where | 3.4524 | 2.7578 | 0.7988 | false |
| 200 | cte_filtered_mat vs join_where | 3.4524 | 3.5286 | 1.0221 | false |
| 200 | cte_filtered_notmat vs join_where | 3.4524 | 3.7494 | 1.0860 | false |
| 200 | cte_joined_outer vs join_where | 3.4524 | 3.7142 | 1.0759 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.4524 | 81.4112 | 23.5813 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.4524 | 3.8035 | 1.1017 | false |
| 200 | comma_join vs join_where | 3.4524 | 3.1600 | 0.9153 | false |
| 2000 | join_on vs join_where | 9.2309 | 8.9381 | 0.9683 | false |
| 2000 | filtered_derived vs join_where | 9.2309 | 8.6844 | 0.9408 | false |
| 2000 | cte_filtered vs join_where | 9.2309 | 9.0044 | 0.9755 | false |
| 2000 | cte_filtered_mat vs join_where | 9.2309 | 10.5318 | 1.1409 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.2309 | 9.4721 | 1.0261 | false |
| 2000 | cte_joined_outer vs join_where | 9.2309 | 9.1821 | 0.9947 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.2309 | 81.7059 | 8.8514 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.2309 | 8.1871 | 0.8869 | false |
| 2000 | comma_join vs join_where | 9.2309 | 9.0737 | 0.9830 | false |
| 20000 | join_on vs join_where | 61.5973 | 72.8096 | 1.1820 | false |
| 20000 | filtered_derived vs join_where | 61.5973 | 62.6936 | 1.0178 | false |
| 20000 | cte_filtered vs join_where | 61.5973 | 64.0440 | 1.0397 | false |
| 20000 | cte_filtered_mat vs join_where | 61.5973 | 66.3889 | 1.0778 | false |
| 20000 | cte_filtered_notmat vs join_where | 61.5973 | 63.3562 | 1.0286 | false |
| 20000 | cte_joined_outer vs join_where | 61.5973 | 63.4425 | 1.0300 | false |
| 20000 | cte_joined_outer_mat vs join_where | 61.5973 | 85.8457 | 1.3937 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 61.5973 | 63.8924 | 1.0373 | false |
| 20000 | comma_join vs join_where | 61.5973 | 62.3997 | 1.0130 | false |
| 100000 | join_on vs join_where | 86.4668 | 88.1216 | 1.0191 | false |
| 100000 | filtered_derived vs join_where | 86.4668 | 86.1025 | 0.9958 | false |
| 100000 | cte_filtered vs join_where | 86.4668 | 97.2561 | 1.1248 | false |
| 100000 | cte_filtered_mat vs join_where | 86.4668 | 90.1374 | 1.0425 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.4668 | 85.0512 | 0.9836 | false |
| 100000 | cte_joined_outer vs join_where | 86.4668 | 84.6604 | 0.9791 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.4668 | 110.3123 | 1.2758 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.4668 | 84.5678 | 0.9780 | false |
| 100000 | comma_join vs join_where | 86.4668 | 86.9163 | 1.0052 | false |
| 180000 | join_on vs join_where | 104.4806 | 111.3741 | 1.0660 | false |
| 180000 | filtered_derived vs join_where | 104.4806 | 107.2414 | 1.0264 | false |
| 180000 | cte_filtered vs join_where | 104.4806 | 106.3247 | 1.0177 | false |
| 180000 | cte_filtered_mat vs join_where | 104.4806 | 110.0455 | 1.0533 | false |
| 180000 | cte_filtered_notmat vs join_where | 104.4806 | 106.2137 | 1.0166 | false |
| 180000 | cte_joined_outer vs join_where | 104.4806 | 105.4371 | 1.0092 | false |
| 180000 | cte_joined_outer_mat vs join_where | 104.4806 | 104.9121 | 1.0041 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 104.4806 | 101.9260 | 0.9755 | false |
| 180000 | comma_join vs join_where | 104.4806 | 103.1727 | 0.9875 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.9016 | 2.1760 | 1.1443 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.9016 | 2.2279 | 1.1716 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.2924 | 82.2052 | 35.8604 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.2924 | 1.7995 | 0.7850 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 2.7578 | 3.5286 | 1.2795 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 2.7578 | 3.7494 | 1.3596 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.7142 | 81.4112 | 21.9187 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.7142 | 3.8035 | 1.0240 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.0044 | 10.5318 | 1.1696 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.0044 | 9.4721 | 1.0519 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.1821 | 81.7059 | 8.8984 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.1821 | 8.1871 | 0.8916 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 64.0440 | 66.3889 | 1.0366 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 64.0440 | 63.3562 | 0.9893 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 63.4425 | 85.8457 | 1.3531 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 63.4425 | 63.8924 | 1.0071 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 97.2561 | 90.1374 | 0.9268 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 97.2561 | 85.0512 | 0.8745 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 84.6604 | 110.3123 | 1.3030 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 84.6604 | 84.5678 | 0.9989 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 106.3247 | 110.0455 | 1.0350 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 106.3247 | 106.2137 | 0.9990 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 105.4371 | 104.9121 | 0.9950 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 105.4371 | 101.9260 | 0.9667 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.8819 | 1 | `0616765b638183a3` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.3429 | 1 | `b1ff127ad056aacb` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 8.9097 | 2 | `d5f12afc2fe89d25` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 9.1042 | 2 | `fdb2f1296885524a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.7220 | 4 | `c306aadcbb5dbf90` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.4203 | 4 | `7abd1c83f743843a` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 9.4759 | 8 | `d4338921d36bb206` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 8.8217 | 8 | `7b6b70576d36f24c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.4192 | 16 | `3f15c0eb4038c66b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1120 | 16 | `34829ec3bc0cc069` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.5892 | 1 | `f12d9f3b6dde8ba0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.4848 | 1 | `eb08abc385501aaa` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.5170 | 2 | `fb1dfb3bbbe8198b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.7321 | 2 | `dc693daf88f53301` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.7515 | 4 | `9554effbd411f2e7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.0185 | 4 | `d9b46a7483087df6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.9534 | 8 | `7d4da4b511700bd3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.6843 | 8 | `2ab197f6a7dd417c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 10.0850 | 16 | `f91aeeed775ea037` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.0303 | 16 | `619a1a340981feb7` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.3023 | 1 | `dc73a878034f9f02` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.2538 | 1 | `a913e2b2e4582f12` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.1498 | 2 | `4ed1b9b68e4a7c5e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 16.2452 | 2 | `811ee9f157b0701d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 22.9285 | 4 | `16d69f33909f3bbc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 29.7848 | 4 | `8b038fa99d03f73e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 69.0481 | 8 | `70fe513109c4599d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 67.6760 | 8 | `5aa4fa1ef39b26c4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 82.9430 | 16 | `5ae860453e5c57d5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 121.7654 | 16 | `78bc3e6b034c3ab0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3934 | false |
| notmat | 1 | 1% | 0.9158 | false |
| mat | 1 | 100% | 1.0152 | false |
| notmat | 1 | 100% | 0.8834 | false |
| mat | 2 | 1% | 0.9559 | false |
| notmat | 2 | 1% | 1.3637 | false |
| mat | 2 | 100% | 0.9591 | false |
| notmat | 2 | 100% | 1.7844 | false |
| mat | 4 | 1% | 1.0034 | false |
| notmat | 4 | 1% | 2.6288 | false |
| mat | 4 | 100% | 1.0635 | false |
| notmat | 4 | 100% | 3.1618 | false |
| mat | 8 | 1% | 0.9449 | false |
| notmat | 8 | 1% | 7.2867 | false |
| mat | 8 | 100% | 1.0978 | false |
| notmat | 8 | 100% | 7.6715 | false |
| mat | 16 | 1% | 1.0707 | false |
| notmat | 16 | 1% | 8.8057 | false |
| mat | 16 | 100% | 1.0908 | false |
| notmat | 16 | 100% | 12.0417 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 29.784809 | 28.996027 | 0.9735 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `engine default`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 1 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 200: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 627, nr_throttled 238, throttled_usec 7251595 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 693, nr_throttled 93, throttled_usec 2803025 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 202050 | 0 | 0/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 202050 | 0 | 0/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100%;position=last | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 82.5541 | 86.9163 | 91.2830 | 5 | 11.5 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 100.5606 | 103.1727 | 112.2665 | 5 | 9.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.7507 | 1.8254 | 2.1766 | 5 | 547.8 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.8943 | 3.1600 | 3.3948 | 5 | 316.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.9001 | 9.0737 | 9.3490 | 5 | 110.2 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 59.5102 | 62.3997 | 66.9163 | 5 | 16.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 90.7765 | 97.2561 | 117.5571 | 5 | 10.3 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 103.5794 | 106.3247 | 112.0962 | 5 | 9.4 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.8648 | 1.9016 | 1.9084 | 5 | 525.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 2.4790 | 2.7578 | 3.1477 | 5 | 362.6 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.6119 | 9.0044 | 9.4119 | 5 | 111.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 62.1072 | 64.0440 | 66.4950 | 5 | 15.6 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 89.0536 | 90.1374 | 94.5829 | 5 | 11.1 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 107.2696 | 110.0455 | 113.5361 | 5 | 9.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 2.0166 | 2.1760 | 2.3905 | 5 | 459.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 3.4663 | 3.5286 | 3.5337 | 5 | 283.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.7627 | 10.5318 | 10.5946 | 5 | 95.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 63.6151 | 66.3889 | 70.3070 | 5 | 15.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 81.9145 | 85.0512 | 88.1561 | 5 | 11.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 102.5759 | 106.2137 | 110.3262 | 5 | 9.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 2.0203 | 2.2279 | 2.7449 | 5 | 448.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.6775 | 3.7494 | 4.0586 | 5 | 266.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.4631 | 9.4721 | 9.8724 | 5 | 105.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 59.3723 | 63.3562 | 64.4156 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 83.3447 | 84.6604 | 85.8897 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 102.5752 | 105.4371 | 109.7877 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0610 | 2.2924 | 2.7102 | 5 | 436.2 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 3.2519 | 3.7142 | 3.8700 | 5 | 269.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.4578 | 9.1821 | 9.7831 | 5 | 108.9 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 61.6388 | 63.4425 | 66.4566 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 107.8017 | 110.3123 | 122.2520 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.7723 | 104.9121 | 107.3553 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 76.9030 | 82.2052 | 84.2325 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 81.1389 | 81.4112 | 83.8641 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 77.2905 | 81.7059 | 84.9872 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 85.1210 | 85.8457 | 89.5562 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 83.7566 | 84.5678 | 85.4865 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 100.6143 | 101.9260 | 113.8308 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6575 | 1.7995 | 2.0035 | 5 | 555.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 3.4307 | 3.8035 | 4.1029 | 5 | 262.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.9939 | 8.1871 | 10.6768 | 5 | 122.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 62.3821 | 63.8924 | 68.3289 | 5 | 15.7 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 83.7174 | 86.1025 | 88.8078 | 5 | 11.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 101.2186 | 107.2414 | 108.3460 | 5 | 9.3 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 2.1011 | 2.2142 | 2.4705 | 5 | 451.6 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 3.4476 | 3.6249 | 3.8318 | 5 | 275.9 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.3211 | 8.6844 | 8.7936 | 5 | 115.1 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 61.3979 | 62.6936 | 65.4906 | 5 | 16.0 |
| family=bag;stmt=join_on;threshold=100000 | valid | 86.6619 | 88.1216 | 90.6441 | 5 | 11.3 |
| family=bag;stmt=join_on;threshold=180000 | valid | 105.5157 | 111.3741 | 115.8415 | 5 | 9.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.6846 | 1.8375 | 2.7159 | 5 | 544.2 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.8363 | 3.2793 | 3.4352 | 5 | 304.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.9228 | 8.9381 | 10.1996 | 5 | 111.9 |
| family=bag;stmt=join_on;threshold=20000 | valid | 63.9448 | 72.8096 | 84.2337 | 5 | 13.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.6622 | 86.4668 | 88.7316 | 5 | 11.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 98.5550 | 104.4806 | 136.8094 | 5 | 9.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1784 | 2.1886 | 2.2671 | 5 | 456.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.3565 | 3.4524 | 4.1647 | 5 | 289.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.4704 | 9.2309 | 9.3138 | 5 | 108.3 |
| family=bag;stmt=join_where;threshold=20000 | valid | 59.6822 | 61.5973 | 67.3685 | 5 | 16.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.1959 | 9.4192 | 10.0791 | 5 | 106.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5544 | 10.1120 | 11.0722 | 5 | 98.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.6183 | 6.8819 | 7.3832 | 5 | 145.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.0204 | 9.3429 | 9.9278 | 5 | 107.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.5949 | 8.9097 | 9.8232 | 5 | 112.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.9336 | 9.1042 | 9.6079 | 5 | 109.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2325 | 8.7220 | 9.3351 | 5 | 114.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.2571 | 9.4203 | 10.9107 | 5 | 106.2 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 9.2322 | 9.4759 | 10.1594 | 5 | 105.5 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.2822 | 8.8217 | 9.5810 | 5 | 113.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.9988 | 10.0850 | 10.4433 | 5 | 99.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.6461 | 11.0303 | 11.4169 | 5 | 90.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.1879 | 9.5892 | 9.8862 | 5 | 104.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.0789 | 9.4848 | 10.2240 | 5 | 105.4 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.0983 | 8.5170 | 8.8806 | 5 | 117.4 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.6611 | 8.7321 | 9.0213 | 5 | 114.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.4480 | 8.7515 | 8.9464 | 5 | 114.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.5394 | 10.0185 | 11.4698 | 5 | 99.8 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.1617 | 8.9534 | 9.9919 | 5 | 111.7 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.4028 | 9.6843 | 10.3106 | 5 | 103.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 79.0431 | 82.9430 | 86.6028 | 5 | 12.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 81.4987 | 121.7654 | 126.9215 | 5 | 8.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.8691 | 6.3023 | 6.5716 | 5 | 158.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.9898 | 8.2538 | 8.8124 | 5 | 121.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 10.8764 | 12.1498 | 13.2767 | 5 | 82.3 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 16.1089 | 16.2452 | 16.5824 | 5 | 61.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.4492 | 22.9285 | 25.2723 | 5 | 43.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.0131 | 29.7848 | 33.3969 | 5 | 33.6 |
| family=reuse;policy=notmat;refs=4;consumer=100%;position=last | valid | 26.1795 | 28.9960 | 30.6070 | 5 | 34.5 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 57.5924 | 69.0481 | 130.6317 | 5 | 14.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 55.7319 | 67.6760 | 92.2454 | 5 | 14.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42875425`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  4. family=bag;stmt=cte_joined_outer;threshold=100000
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=bag;stmt=join_where;threshold=180000
  7. family=bag;stmt=cte_filtered_mat;threshold=20
  8. family=reuse;policy=mat;refs=4;consumer=1%
  9. family=bag;stmt=cte_filtered_notmat;threshold=20000
  10. family=bag;stmt=comma_join;threshold=100000
  11. family=reuse;policy=default;refs=4;consumer=100%
  12. family=bag;stmt=filtered_derived;threshold=20000
  13. family=bag;stmt=join_where;threshold=200
  14. family=bag;stmt=join_where;threshold=2000
  15. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  16. family=reuse;policy=default;refs=2;consumer=1%
  17. family=reuse;policy=default;refs=8;consumer=100%
  18. family=bag;stmt=cte_joined_outer;threshold=180000
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  20. family=bag;stmt=filtered_derived;threshold=200
  21. family=reuse;policy=default;refs=1;consumer=1%
  22. family=bag;stmt=comma_join;threshold=180000
  23. family=bag;stmt=filtered_derived;threshold=20
  24. family=bag;stmt=join_where;threshold=100000
  25. family=bag;stmt=filtered_derived;threshold=100000
  26. family=bag;stmt=cte_filtered_mat;threshold=200
  27. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  28. family=reuse;policy=mat;refs=1;consumer=100%
  29. family=reuse;policy=notmat;refs=2;consumer=1%
  30. family=bag;stmt=filtered_derived;threshold=2000
  31. family=bag;stmt=cte_filtered;threshold=2000
  32. family=reuse;policy=notmat;refs=1;consumer=100%
  33. family=bag;stmt=comma_join;threshold=20
  34. family=reuse;policy=notmat;refs=8;consumer=1%
  35. family=bag;stmt=cte_joined_outer;threshold=200
  36. family=reuse;policy=notmat;refs=8;consumer=100%
  37. family=bag;stmt=cte_filtered_notmat;threshold=200
  38. family=bag;stmt=join_on;threshold=20000
  39. family=reuse;policy=default;refs=1;consumer=100%
  40. family=bag;stmt=cte_joined_outer_mat;threshold=200
  41. family=reuse;policy=default;refs=8;consumer=1%
  42. family=bag;stmt=cte_filtered;threshold=180000
  43. family=bag;stmt=cte_filtered_mat;threshold=20000
  44. family=reuse;policy=mat;refs=4;consumer=100%
  45. family=bag;stmt=cte_filtered_notmat;threshold=2000
  46. family=bag;stmt=cte_joined_outer;threshold=2000
  47. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  48. family=reuse;policy=notmat;refs=16;consumer=100%
  49. family=bag;stmt=cte_filtered;threshold=100000
  50. family=bag;stmt=cte_filtered_mat;threshold=2000
  51. family=bag;stmt=cte_filtered_mat;threshold=100000
  52. family=bag;stmt=join_on;threshold=100000
  53. family=bag;stmt=cte_joined_outer;threshold=20000
  54. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  55. family=bag;stmt=cte_joined_outer_mat;threshold=20
  56. family=reuse;policy=mat;refs=2;consumer=100%
  57. family=reuse;policy=notmat;refs=16;consumer=1%
  58. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  59. family=reuse;policy=mat;refs=1;consumer=1%
  60. family=bag;stmt=filtered_derived;threshold=180000
  61. family=reuse;policy=notmat;refs=2;consumer=100%
  62. family=bag;stmt=cte_filtered;threshold=20
  63. family=reuse;policy=default;refs=16;consumer=1%
  64. family=bag;stmt=cte_filtered_notmat;threshold=180000
  65. family=bag;stmt=join_where;threshold=20
  66. family=bag;stmt=join_on;threshold=180000
  67. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  68. family=reuse;policy=mat;refs=16;consumer=100%
  69. family=reuse;policy=default;refs=4;consumer=1%
  70. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  71. family=reuse;policy=mat;refs=2;consumer=1%
  72. family=bag;stmt=cte_filtered_notmat;threshold=100000
  73. family=bag;stmt=cte_filtered_notmat;threshold=20
  74. family=bag;stmt=comma_join;threshold=2000
  75. family=bag;stmt=cte_filtered_mat;threshold=180000
  76. family=bag;stmt=cte_joined_outer;threshold=20
  77. family=bag;stmt=join_on;threshold=200
  78. family=reuse;policy=default;refs=16;consumer=100%
  79. family=bag;stmt=join_on;threshold=20
  80. family=bag;stmt=join_on;threshold=2000
  81. family=reuse;policy=mat;refs=8;consumer=1%
  82. family=bag;stmt=comma_join;threshold=200
  83. family=bag;stmt=cte_filtered;threshold=200
  84. family=bag;stmt=cte_filtered;threshold=20000
  85. family=reuse;policy=notmat;refs=4;consumer=1%
  86. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  87. family=bag;stmt=join_where;threshold=20000
  88. family=reuse;policy=mat;refs=8;consumer=100%
  89. family=bag;stmt=comma_join;threshold=20000
  90. family=reuse;policy=default;refs=2;consumer=100%
  91. family=reuse;policy=notmat;refs=4;consumer=100%;position=last

## Live negative controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x01_optional_where | sel_join_optional | true | 1 |  |
| x02_join_limit_before | sel_join_bag | true | 7 |  |

## Excluded cells

None: every cell is valid.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `95f16159d035b389` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `fbe450f5d127dff5` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `99bab6b3aefe3cd4` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `a755a17cf05ea9f5` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `1be5d3259e12f9b1` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `f9000c6da6d653fa` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `f0480e46cf44a2d8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `d7c0c3e77a1955f5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `1330098f6c0eb91c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `c4f3c51a45ed5151` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `a579a4f800c9a9f8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `1bd914d6aa6d108b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `4b98fc4c948708c8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `f41924ba33b76a96` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `c76a0001d8e29f6c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `6a73654ffd3ab4df` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `be1a3428a431626e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `1e1b43d2f497a882` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `0695a4a5be319379` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `3d790520e86f986a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `7ea982bb64cdcfac` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `2bf5ad4833f2a9a7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `50c52f501c3bfa12` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `1daf5a96514ec325` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `feb4116246f5af50` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `b78bf6318c8e462f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `911faf6db2125d53` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `d2fb2995f91d71f0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `b2c0b3fe7590dde4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `6a7a016e4f8b956a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `512ed335e27e36c6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3f12afa747e5c1eb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `213e55ef02312e17` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `92e9c608ac253dfa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `9bc54aa2e9ea3967` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `05fde3710fb5713d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `f7bdf780961dd100` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `7f342a6bfa704a03` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `24980e57c6bab629` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `6c883318b73eff0f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `454e6d84df8a7fe8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `ec02ec0d52ee465d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `7c68062b4ef84b60` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `6338863d026d0912` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `dab0d9de6134bb65` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `930c06a5a11fc1a2` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `51f4a8931c42c25d` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `5535ddb3a6697235` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `d1f976c1e1c7006c` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `eb4229f94c908eed` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `cea97450868cc64f` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `c890c46957378104` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `73fc3de08a3d51e6` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `7e53ccb0c27e8224` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `3991598d9c993ec9` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `8919477864fc9b55` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `f482e4ff37668020` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `14498602861d6ffb` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `8a811d145cfe566c` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `14cb88b863f1f94d` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `3f15c0eb4038c66b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `34829ec3bc0cc069` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0616765b638183a3` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b1ff127ad056aacb` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `d5f12afc2fe89d25` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `fdb2f1296885524a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `c306aadcbb5dbf90` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `7abd1c83f743843a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `d4338921d36bb206` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `7b6b70576d36f24c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `f91aeeed775ea037` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `619a1a340981feb7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f12d9f3b6dde8ba0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eb08abc385501aaa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `fb1dfb3bbbe8198b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `dc693daf88f53301` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `9554effbd411f2e7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `d9b46a7483087df6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `7d4da4b511700bd3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `2ab197f6a7dd417c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `5ae860453e5c57d5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `78bc3e6b034c3ab0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `dc73a878034f9f02` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `a913e2b2e4582f12` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `4ed1b9b68e4a7c5e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `811ee9f157b0701d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `16d69f33909f3bbc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `8b038fa99d03f73e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100%;position=last | r_notmat_n4_c100 | valid | `5248132992f8717e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `70fe513109c4599d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `5aa4fa1ef39b26c4` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
