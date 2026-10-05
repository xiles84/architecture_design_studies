# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-i3-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-i3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `68a755a4e33271cef1b47a2a6e22925bcf393e554a509635321141c29deda022`
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
| 20 (0.01%) | 1.7500 | 2.2524 | 2.2550 | 1.7147 | 2.1491 | 1.7721 | 1.8999 | 80.1526 | 1.5736 | 1.8744 |
| 200 (0.1%) | 3.4474 | 3.1484 | 2.7614 | 3.9767 | 2.8068 | 3.5011 | 2.9245 | 84.0988 | 3.1082 | 3.3052 |
| 2000 (1%) | 9.4030 | 8.7099 | 8.8739 | 8.4272 | 9.1607 | 9.4982 | 10.0949 | 81.4994 | 9.8936 | 9.6032 |
| 20000 (10%) | 64.6913 | 64.3026 | 63.4797 | 62.5907 | 67.4564 | 64.2181 | 62.7383 | 85.5257 | 62.8089 | 62.5773 |
| 100000 (50%) | 86.8166 | 85.3678 | 98.7745 | 88.0537 | 109.9779 | 87.4766 | 86.4482 | 91.1245 | 86.0598 | 86.1487 |
| 180000 (90%) | 108.7572 | 107.7808 | 108.1524 | 105.9201 | 112.4860 | 104.0539 | 107.5754 | 100.4942 | 107.4375 | 110.8970 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.7500 | 2.2524 | 1.2870 | false |
| 20 | filtered_derived vs join_where | 1.7500 | 2.2550 | 1.2886 | false |
| 20 | cte_filtered vs join_where | 1.7500 | 1.7147 | 0.9798 | false |
| 20 | cte_filtered_mat vs join_where | 1.7500 | 2.1491 | 1.2281 | false |
| 20 | cte_filtered_notmat vs join_where | 1.7500 | 1.7721 | 1.0126 | false |
| 20 | cte_joined_outer vs join_where | 1.7500 | 1.8999 | 1.0856 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.7500 | 80.1526 | 45.8008 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.7500 | 1.5736 | 0.8992 | false |
| 20 | comma_join vs join_where | 1.7500 | 1.8744 | 1.0711 | false |
| 200 | join_on vs join_where | 3.4474 | 3.1484 | 0.9133 | false |
| 200 | filtered_derived vs join_where | 3.4474 | 2.7614 | 0.8010 | false |
| 200 | cte_filtered vs join_where | 3.4474 | 3.9767 | 1.1535 | false |
| 200 | cte_filtered_mat vs join_where | 3.4474 | 2.8068 | 0.8142 | false |
| 200 | cte_filtered_notmat vs join_where | 3.4474 | 3.5011 | 1.0156 | false |
| 200 | cte_joined_outer vs join_where | 3.4474 | 2.9245 | 0.8483 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.4474 | 84.0988 | 24.3946 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.4474 | 3.1082 | 0.9016 | false |
| 200 | comma_join vs join_where | 3.4474 | 3.3052 | 0.9587 | false |
| 2000 | join_on vs join_where | 9.4030 | 8.7099 | 0.9263 | false |
| 2000 | filtered_derived vs join_where | 9.4030 | 8.8739 | 0.9437 | false |
| 2000 | cte_filtered vs join_where | 9.4030 | 8.4272 | 0.8962 | false |
| 2000 | cte_filtered_mat vs join_where | 9.4030 | 9.1607 | 0.9742 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.4030 | 9.4982 | 1.0101 | false |
| 2000 | cte_joined_outer vs join_where | 9.4030 | 10.0949 | 1.0736 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.4030 | 81.4994 | 8.6674 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.4030 | 9.8936 | 1.0522 | false |
| 2000 | comma_join vs join_where | 9.4030 | 9.6032 | 1.0213 | false |
| 20000 | join_on vs join_where | 64.6913 | 64.3026 | 0.9940 | false |
| 20000 | filtered_derived vs join_where | 64.6913 | 63.4797 | 0.9813 | false |
| 20000 | cte_filtered vs join_where | 64.6913 | 62.5907 | 0.9675 | false |
| 20000 | cte_filtered_mat vs join_where | 64.6913 | 67.4564 | 1.0427 | false |
| 20000 | cte_filtered_notmat vs join_where | 64.6913 | 64.2181 | 0.9927 | false |
| 20000 | cte_joined_outer vs join_where | 64.6913 | 62.7383 | 0.9698 | false |
| 20000 | cte_joined_outer_mat vs join_where | 64.6913 | 85.5257 | 1.3221 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 64.6913 | 62.8089 | 0.9709 | false |
| 20000 | comma_join vs join_where | 64.6913 | 62.5773 | 0.9673 | false |
| 100000 | join_on vs join_where | 86.8166 | 85.3678 | 0.9833 | false |
| 100000 | filtered_derived vs join_where | 86.8166 | 98.7745 | 1.1377 | false |
| 100000 | cte_filtered vs join_where | 86.8166 | 88.0537 | 1.0143 | false |
| 100000 | cte_filtered_mat vs join_where | 86.8166 | 109.9779 | 1.2668 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.8166 | 87.4766 | 1.0076 | false |
| 100000 | cte_joined_outer vs join_where | 86.8166 | 86.4482 | 0.9958 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.8166 | 91.1245 | 1.0496 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.8166 | 86.0598 | 0.9913 | false |
| 100000 | comma_join vs join_where | 86.8166 | 86.1487 | 0.9923 | false |
| 180000 | join_on vs join_where | 108.7572 | 107.7808 | 0.9910 | false |
| 180000 | filtered_derived vs join_where | 108.7572 | 108.1524 | 0.9944 | false |
| 180000 | cte_filtered vs join_where | 108.7572 | 105.9201 | 0.9739 | false |
| 180000 | cte_filtered_mat vs join_where | 108.7572 | 112.4860 | 1.0343 | false |
| 180000 | cte_filtered_notmat vs join_where | 108.7572 | 104.0539 | 0.9568 | false |
| 180000 | cte_joined_outer vs join_where | 108.7572 | 107.5754 | 0.9891 | false |
| 180000 | cte_joined_outer_mat vs join_where | 108.7572 | 100.4942 | 0.9240 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 108.7572 | 107.4375 | 0.9879 | false |
| 180000 | comma_join vs join_where | 108.7572 | 110.8970 | 1.0197 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.7147 | 2.1491 | 1.2534 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.7147 | 1.7721 | 1.0335 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8999 | 80.1526 | 42.1874 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8999 | 1.5736 | 0.8282 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.9767 | 2.8068 | 0.7058 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.9767 | 3.5011 | 0.8804 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 2.9245 | 84.0988 | 28.7568 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 2.9245 | 3.1082 | 1.0628 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.4272 | 9.1607 | 1.0870 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.4272 | 9.4982 | 1.1271 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.0949 | 81.4994 | 8.0733 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.0949 | 9.8936 | 0.9801 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 62.5907 | 67.4564 | 1.0777 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 62.5907 | 64.2181 | 1.0260 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 62.7383 | 85.5257 | 1.3632 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 62.7383 | 62.8089 | 1.0011 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 88.0537 | 109.9779 | 1.2490 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 88.0537 | 87.4766 | 0.9934 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 86.4482 | 91.1245 | 1.0541 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 86.4482 | 86.0598 | 0.9955 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 105.9201 | 112.4860 | 1.0620 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 105.9201 | 104.0539 | 0.9824 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 107.5754 | 100.4942 | 0.9342 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 107.5754 | 107.4375 | 0.9987 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.3953 | 1 | `846668f655537063` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.6531 | 1 | `482ba77b89be360d` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 9.5131 | 2 | `4229abf29ccf38e7` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.3448 | 2 | `f335dbc51026aafc` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.0887 | 4 | `51dc66cfcb62d2de` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.6979 | 4 | `375df480dde4d735` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 9.1989 | 8 | `3ec43a5bb1b5d330` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 8.3036 | 8 | `c09d73a7237f6c85` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.4162 | 16 | `4c97f6b88b38e168` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8073 | 16 | `9411a4a0819e15a5` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.1884 | 1 | `d5e33384d5ccd4ae` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.9805 | 1 | `bf05b4b8a624efed` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.7947 | 2 | `cf63e64f9e6bb1ea` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.1335 | 2 | `c5e044c5f9149434` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.2328 | 4 | `ace426f656bdbf01` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.0887 | 4 | `11be97c1a74568fa` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.5094 | 8 | `20c9f3cf0a104a7a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 8.5534 | 8 | `66f23d60fb516b36` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.3403 | 16 | `d6706693e512ee36` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.2089 | 16 | `a503e2f1e048fc43` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.7443 | 1 | `685e824af039222e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.8115 | 1 | `703c9a596122acd6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.5784 | 2 | `a699839476f8588f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.6587 | 2 | `f7b22ae7928dc362` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 23.0083 | 4 | `18ea37ad9d2e68ee` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.4025 | 4 | `a2174c2c2baa0777` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 53.9595 | 8 | `ed9151428fc23b17` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 58.2327 | 8 | `e7e77c47ca765130` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 69.5902 | 16 | `2774f5e799e816ea` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 111.7905 | 16 | `a67683dd7271dd48` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1073 | false |
| notmat | 1 | 1% | 0.9120 | false |
| mat | 1 | 100% | 1.0428 | false |
| notmat | 1 | 100% | 1.0207 | false |
| mat | 2 | 1% | 0.9245 | false |
| notmat | 2 | 1% | 1.3222 | false |
| mat | 2 | 100% | 0.9747 | false |
| notmat | 2 | 100% | 1.8765 | false |
| mat | 4 | 1% | 1.0178 | false |
| notmat | 4 | 1% | 2.8445 | false |
| mat | 4 | 100% | 0.9300 | false |
| notmat | 4 | 100% | 3.2655 | false |
| mat | 8 | 1% | 0.9250 | false |
| notmat | 8 | 1% | 5.8659 | false |
| mat | 8 | 100% | 1.0301 | false |
| notmat | 8 | 100% | 7.0129 | false |
| mat | 16 | 1% | 0.9919 | false |
| notmat | 16 | 1% | 7.3905 | false |
| mat | 16 | 100% | 0.9390 | false |
| notmat | 16 | 100% | 11.3987 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 9.198869 | 8.381560 | 0.9112 | false |

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

- client: nr_periods 594, nr_throttled 242, throttled_usec 9119270 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 668, nr_throttled 72, throttled_usec 2826828 | quota 200000 100000, memory 3221225472

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
| family=reuse;policy=default;refs=8;consumer=1%;position=last | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
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
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 84.2588 | 86.1487 | 90.2768 | 5 | 11.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 102.5198 | 110.8970 | 139.4749 | 5 | 9.0 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.6548 | 1.8744 | 1.9320 | 5 | 533.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.7714 | 3.3052 | 3.8680 | 5 | 302.6 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.9409 | 9.6032 | 10.6154 | 5 | 104.1 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 61.7235 | 62.5773 | 65.4950 | 5 | 16.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 84.5202 | 88.0537 | 92.4204 | 5 | 11.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 103.2516 | 105.9201 | 108.7509 | 5 | 9.4 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.3656 | 1.7147 | 1.8366 | 5 | 583.2 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.8934 | 3.9767 | 4.2112 | 5 | 251.5 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.3698 | 8.4272 | 9.5384 | 5 | 118.7 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 60.6285 | 62.5907 | 63.7697 | 5 | 16.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 93.0932 | 109.9779 | 117.4243 | 5 | 9.1 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 110.8330 | 112.4860 | 114.7608 | 5 | 8.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 2.0521 | 2.1491 | 2.2108 | 5 | 465.3 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.6082 | 2.8068 | 3.0290 | 5 | 356.3 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 8.3912 | 9.1607 | 10.4135 | 5 | 109.2 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 65.0866 | 67.4564 | 72.0181 | 5 | 14.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 85.1229 | 87.4766 | 90.8358 | 5 | 11.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 102.7243 | 104.0539 | 108.8096 | 5 | 9.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.6087 | 1.7721 | 1.9719 | 5 | 564.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.2525 | 3.5011 | 3.8103 | 5 | 285.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 9.2561 | 9.4982 | 9.9904 | 5 | 105.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 63.2512 | 64.2181 | 64.9283 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 84.0955 | 86.4482 | 88.5747 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 104.0237 | 107.5754 | 113.6604 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6887 | 1.8999 | 2.1168 | 5 | 526.3 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.8021 | 2.9245 | 3.4195 | 5 | 341.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.1750 | 10.0949 | 11.4144 | 5 | 99.1 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 59.4349 | 62.7383 | 66.4742 | 5 | 15.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 87.2506 | 91.1245 | 92.3272 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 100.1209 | 100.4942 | 104.7578 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 74.4998 | 80.1526 | 88.1990 | 5 | 12.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 79.6560 | 84.0988 | 100.5185 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 80.0278 | 81.4994 | 88.9413 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 84.0644 | 85.5257 | 89.6862 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 83.1510 | 86.0598 | 87.0075 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 101.1954 | 107.4375 | 123.7982 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.5266 | 1.5736 | 1.7688 | 5 | 635.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.6037 | 3.1082 | 3.2383 | 5 | 321.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.2907 | 9.8936 | 10.3106 | 5 | 101.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 60.8418 | 62.8089 | 66.8173 | 5 | 15.9 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 91.1314 | 98.7745 | 113.2674 | 5 | 10.1 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 105.0255 | 108.1524 | 114.8349 | 5 | 9.2 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 2.1586 | 2.2550 | 2.3490 | 5 | 443.5 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.5977 | 2.7614 | 3.0525 | 5 | 362.1 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.4190 | 8.8739 | 9.3219 | 5 | 112.7 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 62.2395 | 63.4797 | 67.2693 | 5 | 15.8 |
| family=bag;stmt=join_on;threshold=100000 | valid | 85.1915 | 85.3678 | 88.3986 | 5 | 11.7 |
| family=bag;stmt=join_on;threshold=180000 | valid | 102.4709 | 107.7808 | 118.3472 | 5 | 9.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.9352 | 2.2524 | 2.4818 | 5 | 444.0 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.9866 | 3.1484 | 3.5920 | 5 | 317.6 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.4330 | 8.7099 | 9.6543 | 5 | 114.8 |
| family=bag;stmt=join_on;threshold=20000 | valid | 59.6409 | 64.3026 | 72.2880 | 5 | 15.6 |
| family=bag;stmt=join_where;threshold=100000 | valid | 85.0681 | 86.8166 | 89.7022 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.9289 | 108.7572 | 117.8050 | 5 | 9.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6722 | 1.7500 | 1.9777 | 5 | 571.4 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.3316 | 3.4474 | 3.6050 | 5 | 290.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.0505 | 9.4030 | 9.6379 | 5 | 106.3 |
| family=bag;stmt=join_where;threshold=20000 | valid | 61.6513 | 64.6913 | 68.3252 | 5 | 15.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.3583 | 9.4162 | 10.4917 | 5 | 106.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.2811 | 9.8073 | 10.2308 | 5 | 102.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.2223 | 7.3953 | 8.0059 | 5 | 135.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.3522 | 7.6531 | 8.1507 | 5 | 130.7 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.7739 | 9.5131 | 38.0799 | 5 | 105.1 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.0309 | 8.3448 | 8.9282 | 5 | 119.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.6256 | 8.0887 | 8.5254 | 5 | 123.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.4741 | 8.6979 | 9.0195 | 5 | 115.0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.5497 | 9.1989 | 9.4814 | 5 | 108.7 |
| family=reuse;policy=default;refs=8;consumer=1%;position=last | valid | 8.0566 | 8.3816 | 9.4042 | 5 | 119.3 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 7.9478 | 8.3036 | 8.6759 | 5 | 120.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.2445 | 9.3403 | 11.0336 | 5 | 107.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.5103 | 9.2089 | 9.3521 | 5 | 108.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.8191 | 8.1884 | 8.2033 | 5 | 122.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.5561 | 7.9805 | 8.6415 | 5 | 125.3 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.0939 | 8.7947 | 9.1479 | 5 | 113.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 7.7586 | 8.1335 | 8.8728 | 5 | 122.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.1386 | 8.2328 | 8.6349 | 5 | 121.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.9012 | 8.0887 | 8.3135 | 5 | 123.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.0501 | 8.5094 | 8.9989 | 5 | 117.5 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 8.2137 | 8.5534 | 10.5816 | 5 | 116.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 68.4822 | 69.5902 | 80.2444 | 5 | 14.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 108.7829 | 111.7905 | 115.6154 | 5 | 8.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.3906 | 6.7443 | 7.5681 | 5 | 148.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.5622 | 7.8115 | 9.2977 | 5 | 128.0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 12.1121 | 12.5784 | 13.4868 | 5 | 79.5 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 15.0373 | 15.6587 | 17.5170 | 5 | 63.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.1092 | 23.0083 | 25.7368 | 5 | 43.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.4840 | 28.4025 | 35.8006 | 5 | 35.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 52.3630 | 53.9595 | 56.9314 | 5 | 18.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 56.3602 | 58.2327 | 68.6931 | 5 | 17.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42559523`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=8;consumer=1%
  2. family=bag;stmt=cte_filtered;threshold=180000
  3. family=reuse;policy=default;refs=2;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  6. family=bag;stmt=cte_filtered_notmat;threshold=200
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  8. family=bag;stmt=comma_join;threshold=100000
  9. family=bag;stmt=comma_join;threshold=2000
  10. family=bag;stmt=cte_filtered_mat;threshold=20000
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  12. family=bag;stmt=cte_filtered;threshold=20000
  13. family=bag;stmt=join_where;threshold=200
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  15. family=bag;stmt=cte_filtered;threshold=2000
  16. family=reuse;policy=default;refs=16;consumer=100%
  17. family=reuse;policy=mat;refs=1;consumer=1%
  18. family=reuse;policy=default;refs=4;consumer=1%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  20. family=bag;stmt=join_on;threshold=20
  21. family=bag;stmt=cte_filtered_mat;threshold=180000
  22. family=bag;stmt=comma_join;threshold=20000
  23. family=reuse;policy=default;refs=4;consumer=100%
  24. family=reuse;policy=default;refs=8;consumer=100%
  25. family=bag;stmt=join_on;threshold=2000
  26. family=reuse;policy=notmat;refs=8;consumer=100%
  27. family=bag;stmt=filtered_derived;threshold=100000
  28. family=bag;stmt=cte_filtered_mat;threshold=20
  29. family=bag;stmt=join_where;threshold=100000
  30. family=bag;stmt=filtered_derived;threshold=180000
  31. family=reuse;policy=default;refs=1;consumer=1%
  32. family=reuse;policy=mat;refs=4;consumer=100%
  33. family=reuse;policy=notmat;refs=1;consumer=100%
  34. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  35. family=reuse;policy=notmat;refs=2;consumer=100%
  36. family=reuse;policy=mat;refs=1;consumer=100%
  37. family=bag;stmt=cte_joined_outer;threshold=200
  38. family=reuse;policy=mat;refs=8;consumer=1%
  39. family=reuse;policy=default;refs=1;consumer=100%
  40. family=reuse;policy=mat;refs=16;consumer=100%
  41. family=bag;stmt=filtered_derived;threshold=2000
  42. family=bag;stmt=cte_filtered_mat;threshold=2000
  43. family=bag;stmt=join_where;threshold=20
  44. family=bag;stmt=filtered_derived;threshold=200
  45. family=bag;stmt=cte_filtered;threshold=20
  46. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  47. family=reuse;policy=notmat;refs=16;consumer=100%
  48. family=bag;stmt=join_where;threshold=180000
  49. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  50. family=reuse;policy=mat;refs=2;consumer=100%
  51. family=bag;stmt=cte_joined_outer_mat;threshold=200
  52. family=reuse;policy=mat;refs=2;consumer=1%
  53. family=bag;stmt=comma_join;threshold=180000
  54. family=bag;stmt=cte_filtered;threshold=200
  55. family=reuse;policy=mat;refs=8;consumer=100%
  56. family=bag;stmt=cte_filtered_notmat;threshold=180000
  57. family=reuse;policy=notmat;refs=2;consumer=1%
  58. family=bag;stmt=comma_join;threshold=20
  59. family=reuse;policy=notmat;refs=1;consumer=1%
  60. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  61. family=bag;stmt=join_on;threshold=180000
  62. family=bag;stmt=join_on;threshold=100000
  63. family=bag;stmt=join_on;threshold=20000
  64. family=bag;stmt=filtered_derived;threshold=20
  65. family=bag;stmt=comma_join;threshold=200
  66. family=reuse;policy=default;refs=2;consumer=100%
  67. family=bag;stmt=cte_joined_outer;threshold=20000
  68. family=bag;stmt=join_where;threshold=20000
  69. family=bag;stmt=cte_filtered_notmat;threshold=2000
  70. family=bag;stmt=cte_joined_outer;threshold=20
  71. family=bag;stmt=cte_filtered_mat;threshold=200
  72. family=reuse;policy=notmat;refs=4;consumer=1%
  73. family=reuse;policy=mat;refs=16;consumer=1%
  74. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  75. family=bag;stmt=cte_joined_outer;threshold=100000
  76. family=bag;stmt=cte_joined_outer;threshold=2000
  77. family=bag;stmt=join_on;threshold=200
  78. family=bag;stmt=cte_joined_outer_mat;threshold=20
  79. family=bag;stmt=join_where;threshold=2000
  80. family=reuse;policy=mat;refs=4;consumer=1%
  81. family=reuse;policy=notmat;refs=4;consumer=100%
  82. family=bag;stmt=cte_filtered_notmat;threshold=20
  83. family=reuse;policy=notmat;refs=16;consumer=1%
  84. family=bag;stmt=cte_filtered_mat;threshold=100000
  85. family=bag;stmt=cte_filtered;threshold=100000
  86. family=reuse;policy=notmat;refs=8;consumer=1%
  87. family=bag;stmt=cte_filtered_notmat;threshold=20000
  88. family=bag;stmt=cte_filtered_notmat;threshold=100000
  89. family=reuse;policy=default;refs=16;consumer=1%
  90. family=bag;stmt=filtered_derived;threshold=20000
  91. family=reuse;policy=default;refs=8;consumer=1%;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `0e5ac17981b74ed8` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `9fd2bb90c12ecd26` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `33729cac30e2819a` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `39a51092121ea48a` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `50f6c37d58e1023f` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `5cb851f29600e31e` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `f0bebc6a0b4eaede` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `f5ededdde6485102` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `8fbeca8c331bbd50` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `0ba4454d70f7ec95` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `bbb04b31d37cd6d6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `304f52adf17de5fa` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `c59434aa3af1a5ab` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `9a6f63589b4456bb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `ff92fba7b0438a41` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `b4355d89b78bd615` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `302576cd63672bcf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `3d617be8368c4fb5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `0417ed40199bcd39` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `e4be80f1cc372e78` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `105c10aa650eb6d4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `4ffce1b6e35c46e4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `921e19e3daeee2a4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `61d985a551155dfc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `7de275075fa265be` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `16ae5938ccde6dbd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `a020681b8c22217c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `588a30d0c697c9d9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `5eb83f49731edb32` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `355dbcd5ab65935a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `b11c6f9051d632d0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2387acae402a9fac` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d68911a3ff3af2a1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `4e11b494c7b91dea` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `10d6108af084742d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `8a41fa9635e1ed7e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `08bf938789236d6f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `124eb378531bbdd5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `7ecd7c70ac981605` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `2b0bac54b80faeef` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `16991cd1c16e2de7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `5eb7e53b0e1c04b9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `3a89e16e48af7109` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `beec2dafbabbef22` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `deae89a0f623362f` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `3690af8f6d33e919` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `531695f4d0f4c943` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `db73db4a04ff7c1d` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `4e0903eeb9f17f88` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `fdb659cd900032a8` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `363a421ef625a3d6` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `c8f56a695071ec75` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `4a21816ff1c72a25` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `5c33c874e90df46a` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `7365944a632c3248` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `bd6882786916ceb4` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `d315e996f3984fb2` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `fbe57236b4262948` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `496ed2b882bc212e` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `3622bb1511e08719` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4c97f6b88b38e168` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `9411a4a0819e15a5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `846668f655537063` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `482ba77b89be360d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `4229abf29ccf38e7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `f335dbc51026aafc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `51dc66cfcb62d2de` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `375df480dde4d735` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `3ec43a5bb1b5d330` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1%;position=last | r_default_n8_c1 | valid | `930481925ae9cd3b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `c09d73a7237f6c85` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d6706693e512ee36` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a503e2f1e048fc43` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `d5e33384d5ccd4ae` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `bf05b4b8a624efed` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `cf63e64f9e6bb1ea` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `c5e044c5f9149434` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ace426f656bdbf01` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `11be97c1a74568fa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `20c9f3cf0a104a7a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `66f23d60fb516b36` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2774f5e799e816ea` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a67683dd7271dd48` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `685e824af039222e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `703c9a596122acd6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `a699839476f8588f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `f7b22ae7928dc362` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `18ea37ad9d2e68ee` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `a2174c2c2baa0777` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `ed9151428fc23b17` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `e7e77c47ca765130` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
