# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-i2-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `4d80c4a5b4cee2cb40c3189670cb317f74e28f78`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-i2-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `db129e012d85d78aca8b727eb02c65792acca57222e947eec69bd26ec93db187`
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
| 20 (0.01%) | 1.9808 | 1.9932 | 1.5264 | 1.7328 | 1.8442 | 1.9107 | 1.8041 | 89.1566 | 1.8512 | 1.6286 |
| 200 (0.1%) | 3.7880 | 3.1924 | 3.2900 | 3.3905 | 2.7809 | 3.0823 | 3.7587 | 80.7199 | 2.9830 | 3.1820 |
| 2000 (1%) | 9.7930 | 8.7578 | 8.8885 | 9.0793 | 9.9455 | 9.8563 | 9.0221 | 85.3225 | 9.4778 | 8.4301 |
| 20000 (10%) | 61.8723 | 63.1383 | 63.4888 | 62.9863 | 66.5280 | 61.8222 | 66.4610 | 86.3475 | 62.9674 | 65.7911 |
| 100000 (50%) | 86.9464 | 86.5040 | 86.0126 | 85.1563 | 92.3968 | 89.4241 | 88.3656 | 92.0068 | 87.7690 | 86.7234 |
| 180000 (90%) | 108.2321 | 106.8943 | 110.0270 | 109.3550 | 115.9215 | 111.6208 | 111.1171 | 102.4661 | 110.4011 | 109.4383 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.9808 | 1.9932 | 1.0063 | false |
| 20 | filtered_derived vs join_where | 1.9808 | 1.5264 | 0.7706 | false |
| 20 | cte_filtered vs join_where | 1.9808 | 1.7328 | 0.8748 | false |
| 20 | cte_filtered_mat vs join_where | 1.9808 | 1.8442 | 0.9311 | false |
| 20 | cte_filtered_notmat vs join_where | 1.9808 | 1.9107 | 0.9646 | false |
| 20 | cte_joined_outer vs join_where | 1.9808 | 1.8041 | 0.9108 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.9808 | 89.1566 | 45.0112 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.9808 | 1.8512 | 0.9346 | false |
| 20 | comma_join vs join_where | 1.9808 | 1.6286 | 0.8222 | false |
| 200 | join_on vs join_where | 3.7880 | 3.1924 | 0.8428 | false |
| 200 | filtered_derived vs join_where | 3.7880 | 3.2900 | 0.8685 | false |
| 200 | cte_filtered vs join_where | 3.7880 | 3.3905 | 0.8950 | false |
| 200 | cte_filtered_mat vs join_where | 3.7880 | 2.7809 | 0.7341 | false |
| 200 | cte_filtered_notmat vs join_where | 3.7880 | 3.0823 | 0.8137 | false |
| 200 | cte_joined_outer vs join_where | 3.7880 | 3.7587 | 0.9923 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.7880 | 80.7199 | 21.3093 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.7880 | 2.9830 | 0.7875 | false |
| 200 | comma_join vs join_where | 3.7880 | 3.1820 | 0.8400 | false |
| 2000 | join_on vs join_where | 9.7930 | 8.7578 | 0.8943 | false |
| 2000 | filtered_derived vs join_where | 9.7930 | 8.8885 | 0.9076 | false |
| 2000 | cte_filtered vs join_where | 9.7930 | 9.0793 | 0.9271 | false |
| 2000 | cte_filtered_mat vs join_where | 9.7930 | 9.9455 | 1.0156 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.7930 | 9.8563 | 1.0065 | false |
| 2000 | cte_joined_outer vs join_where | 9.7930 | 9.0221 | 0.9213 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.7930 | 85.3225 | 8.7126 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.7930 | 9.4778 | 0.9678 | false |
| 2000 | comma_join vs join_where | 9.7930 | 8.4301 | 0.8608 | false |
| 20000 | join_on vs join_where | 61.8723 | 63.1383 | 1.0205 | false |
| 20000 | filtered_derived vs join_where | 61.8723 | 63.4888 | 1.0261 | false |
| 20000 | cte_filtered vs join_where | 61.8723 | 62.9863 | 1.0180 | false |
| 20000 | cte_filtered_mat vs join_where | 61.8723 | 66.5280 | 1.0752 | false |
| 20000 | cte_filtered_notmat vs join_where | 61.8723 | 61.8222 | 0.9992 | false |
| 20000 | cte_joined_outer vs join_where | 61.8723 | 66.4610 | 1.0742 | false |
| 20000 | cte_joined_outer_mat vs join_where | 61.8723 | 86.3475 | 1.3956 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 61.8723 | 62.9674 | 1.0177 | false |
| 20000 | comma_join vs join_where | 61.8723 | 65.7911 | 1.0633 | false |
| 100000 | join_on vs join_where | 86.9464 | 86.5040 | 0.9949 | false |
| 100000 | filtered_derived vs join_where | 86.9464 | 86.0126 | 0.9893 | false |
| 100000 | cte_filtered vs join_where | 86.9464 | 85.1563 | 0.9794 | false |
| 100000 | cte_filtered_mat vs join_where | 86.9464 | 92.3968 | 1.0627 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.9464 | 89.4241 | 1.0285 | false |
| 100000 | cte_joined_outer vs join_where | 86.9464 | 88.3656 | 1.0163 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.9464 | 92.0068 | 1.0582 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.9464 | 87.7690 | 1.0095 | false |
| 100000 | comma_join vs join_where | 86.9464 | 86.7234 | 0.9974 | false |
| 180000 | join_on vs join_where | 108.2321 | 106.8943 | 0.9876 | false |
| 180000 | filtered_derived vs join_where | 108.2321 | 110.0270 | 1.0166 | false |
| 180000 | cte_filtered vs join_where | 108.2321 | 109.3550 | 1.0104 | false |
| 180000 | cte_filtered_mat vs join_where | 108.2321 | 115.9215 | 1.0710 | false |
| 180000 | cte_filtered_notmat vs join_where | 108.2321 | 111.6208 | 1.0313 | false |
| 180000 | cte_joined_outer vs join_where | 108.2321 | 111.1171 | 1.0267 | false |
| 180000 | cte_joined_outer_mat vs join_where | 108.2321 | 102.4661 | 0.9467 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 108.2321 | 110.4011 | 1.0200 | false |
| 180000 | comma_join vs join_where | 108.2321 | 109.4383 | 1.0111 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.7328 | 1.8442 | 1.0643 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.7328 | 1.9107 | 1.1026 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8041 | 89.1566 | 49.4192 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8041 | 1.8512 | 1.0261 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.3905 | 2.7809 | 0.8202 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.3905 | 3.0823 | 0.9091 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.7587 | 80.7199 | 21.4755 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.7587 | 2.9830 | 0.7936 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.0793 | 9.9455 | 1.0954 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.0793 | 9.8563 | 1.0856 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.0221 | 85.3225 | 9.4570 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.0221 | 9.4778 | 1.0505 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 62.9863 | 66.5280 | 1.0562 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 62.9863 | 61.8222 | 0.9815 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 66.4610 | 86.3475 | 1.2992 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 66.4610 | 62.9674 | 0.9474 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 85.1563 | 92.3968 | 1.0850 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 85.1563 | 89.4241 | 1.0501 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 88.3656 | 92.0068 | 1.0412 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 88.3656 | 87.7690 | 0.9932 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 109.3550 | 115.9215 | 1.0600 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 109.3550 | 111.6208 | 1.0207 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 111.1171 | 102.4661 | 0.9221 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 111.1171 | 110.4011 | 0.9936 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.5416 | 1 | `08171b32ab32e67e` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.9331 | 1 | `06fc577b962275bc` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 9.3131 | 2 | `cd8ab13ee2e31b8a` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.7073 | 2 | `c64f9bbd828d5c08` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.1684 | 4 | `15d111d2c537806e` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.8771 | 4 | `1cec421a490be70d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 9.5930 | 8 | `4d7d9c6570e762bd` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.7196 | 8 | `b79b4ac28e1dee8b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 10.0268 | 16 | `8ed4cf98543615d2` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.4214 | 16 | `036da7e3e98be779` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.1139 | 1 | `9081b0db7576b82a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.6135 | 1 | `d9f17516fd69c3ce` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 9.3026 | 2 | `b3a3ebc0b071eee8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.9169 | 2 | `57a28938e63a984c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.4049 | 4 | `98e235f909649301` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.8718 | 4 | `1384b10b627c6210` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.7730 | 8 | `88a95b60ff9c0825` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.9762 | 8 | `4aece0680dd4db81` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 10.2498 | 16 | `329ce5f6f6cbe364` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.3550 | 16 | `a52550c160030633` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.3083 | 1 | `e37bd4b030ba1c8d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.2527 | 1 | `ac688dab41a52008` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 11.2013 | 2 | `ce45afe5bc55e3b3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.8364 | 2 | `8f13a32b120dc8bb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 22.8469 | 4 | `34b843591abbe534` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.1775 | 4 | `506ab26b1348335c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 57.3604 | 8 | `4d5de15f35b45d1b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 62.3365 | 8 | `945fdca57bd345aa` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 100.2314 | 16 | `c9d956836c9d95db` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 87.1351 | 16 | `337c3aaedab6f343` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3932 | false |
| notmat | 1 | 1% | 1.1172 | false |
| mat | 1 | 100% | 1.0858 | false |
| notmat | 1 | 100% | 1.0403 | false |
| mat | 2 | 1% | 0.9989 | false |
| notmat | 2 | 1% | 1.2028 | false |
| mat | 2 | 100% | 1.0241 | false |
| notmat | 2 | 100% | 1.8188 | false |
| mat | 4 | 1% | 1.0289 | false |
| notmat | 4 | 1% | 2.7970 | false |
| mat | 4 | 100% | 0.9994 | false |
| notmat | 4 | 100% | 3.3995 | false |
| mat | 8 | 1% | 0.9145 | false |
| notmat | 8 | 1% | 5.9794 | false |
| mat | 8 | 100% | 1.0264 | false |
| notmat | 8 | 100% | 6.4135 | false |
| mat | 16 | 1% | 1.0222 | false |
| notmat | 16 | 1% | 9.9963 | false |
| mat | 16 | 100% | 0.9936 | false |
| notmat | 16 | 100% | 8.3612 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 92.396814 | 93.963310 | 1.0170 | false |

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

- client: nr_periods 624, nr_throttled 266, throttled_usec 10545057 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 691, nr_throttled 71, throttled_usec 2647320 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
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
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 85.8823 | 86.7234 | 93.2625 | 5 | 11.5 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 107.2738 | 109.4383 | 116.7438 | 5 | 9.1 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.4900 | 1.6286 | 2.0480 | 5 | 614.0 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.7615 | 3.1820 | 4.0626 | 5 | 314.3 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.2315 | 8.4301 | 9.5253 | 5 | 118.6 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 61.1569 | 65.7911 | 67.8132 | 5 | 15.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 82.2105 | 85.1563 | 88.4194 | 5 | 11.7 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 102.1065 | 109.3550 | 114.2971 | 5 | 9.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.7078 | 1.7328 | 1.8678 | 5 | 577.1 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.3524 | 3.3905 | 4.1091 | 5 | 294.9 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.7806 | 9.0793 | 9.8690 | 5 | 110.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 62.6171 | 62.9863 | 70.8401 | 5 | 15.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 88.8805 | 92.3968 | 100.8969 | 5 | 10.8 |
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | valid | 93.5058 | 93.9633 | 94.9275 | 5 | 10.6 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 109.3426 | 115.9215 | 116.8148 | 5 | 8.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.6258 | 1.8442 | 2.4880 | 5 | 542.2 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.6915 | 2.7809 | 2.9687 | 5 | 359.6 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.9026 | 9.9455 | 11.5550 | 5 | 100.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 62.3981 | 66.5280 | 69.5015 | 5 | 15.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 86.6158 | 89.4241 | 113.7266 | 5 | 11.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 100.8394 | 111.6208 | 115.4820 | 5 | 9.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.6732 | 1.9107 | 2.4387 | 5 | 523.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.8479 | 3.0823 | 3.4097 | 5 | 324.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 9.4475 | 9.8563 | 10.3965 | 5 | 101.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 60.2886 | 61.8222 | 64.2429 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 83.0292 | 88.3656 | 90.7745 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 102.3386 | 111.1171 | 116.4192 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7522 | 1.8041 | 2.1026 | 5 | 554.3 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 3.5730 | 3.7587 | 3.8166 | 5 | 266.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.2357 | 9.0221 | 9.4755 | 5 | 110.8 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 61.6540 | 66.4610 | 82.1280 | 5 | 15.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 89.3528 | 92.0068 | 95.6352 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.3061 | 102.4661 | 106.9909 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 84.6896 | 89.1566 | 97.6013 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 75.5910 | 80.7199 | 84.7383 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 77.3506 | 85.3225 | 90.3783 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 85.3529 | 86.3475 | 95.7527 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 85.4035 | 87.7690 | 92.0248 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 106.1952 | 110.4011 | 115.0039 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7314 | 1.8512 | 2.2408 | 5 | 540.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.9319 | 2.9830 | 3.3428 | 5 | 335.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.2516 | 9.4778 | 10.5532 | 5 | 105.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 60.7776 | 62.9674 | 67.3787 | 5 | 15.9 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 83.2483 | 86.0126 | 89.5862 | 5 | 11.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 104.2112 | 110.0270 | 114.3430 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.4479 | 1.5264 | 1.6446 | 5 | 655.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.9190 | 3.2900 | 17.2614 | 5 | 304.0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.5046 | 8.8885 | 9.4797 | 5 | 112.5 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 60.8290 | 63.4888 | 68.5802 | 5 | 15.8 |
| family=bag;stmt=join_on;threshold=100000 | valid | 80.9855 | 86.5040 | 87.7072 | 5 | 11.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 104.7623 | 106.8943 | 110.8582 | 5 | 9.4 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.7624 | 1.9932 | 2.4361 | 5 | 501.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.8538 | 3.1924 | 3.8788 | 5 | 313.2 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.5841 | 8.7578 | 9.2829 | 5 | 114.2 |
| family=bag;stmt=join_on;threshold=20000 | valid | 59.9162 | 63.1383 | 67.8764 | 5 | 15.8 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.3674 | 86.9464 | 90.2626 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.8204 | 108.2321 | 115.9978 | 5 | 9.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7376 | 1.9808 | 2.1884 | 5 | 504.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.4416 | 3.7880 | 4.0981 | 5 | 264.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3606 | 9.7930 | 10.5195 | 5 | 102.1 |
| family=bag;stmt=join_where;threshold=20000 | valid | 59.7359 | 61.8723 | 62.9977 | 5 | 16.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.2411 | 10.0268 | 10.2597 | 5 | 99.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.9060 | 10.4214 | 11.3380 | 5 | 96.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.3860 | 6.5416 | 6.8756 | 5 | 152.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.4708 | 7.9331 | 8.3643 | 5 | 126.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.8858 | 9.3131 | 9.6848 | 5 | 107.4 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.6650 | 8.7073 | 9.1534 | 5 | 114.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.5489 | 8.1684 | 8.7283 | 5 | 122.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.4405 | 8.8771 | 9.3347 | 5 | 112.6 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 9.4001 | 9.5930 | 9.8906 | 5 | 104.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.7200 | 9.7196 | 10.1945 | 5 | 102.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.5469 | 10.2498 | 10.3671 | 5 | 97.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.0127 | 10.3550 | 10.8640 | 5 | 96.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.7394 | 9.1139 | 9.3710 | 5 | 109.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.2403 | 8.6135 | 9.8062 | 5 | 116.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.8747 | 9.3026 | 9.8661 | 5 | 107.5 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.2686 | 8.9169 | 9.1267 | 5 | 112.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.0067 | 8.4049 | 9.3352 | 5 | 119.0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.6058 | 8.8718 | 9.1068 | 5 | 112.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.3153 | 8.7730 | 8.9922 | 5 | 114.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.3338 | 9.9762 | 10.7666 | 5 | 100.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 71.1473 | 100.2314 | 107.7861 | 5 | 10.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 79.5576 | 87.1351 | 127.0121 | 5 | 11.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.0689 | 7.3083 | 8.0679 | 5 | 136.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.8978 | 8.2527 | 8.9588 | 5 | 121.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 10.8597 | 11.2013 | 11.5066 | 5 | 89.3 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.7371 | 15.8364 | 16.9409 | 5 | 63.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.5944 | 22.8469 | 23.6438 | 5 | 43.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.4085 | 30.1775 | 33.7435 | 5 | 33.1 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 53.0544 | 57.3604 | 59.2672 | 5 | 17.4 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 56.8520 | 62.3365 | 84.6357 | 5 | 16.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42785812`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=100000
  2. family=reuse;policy=mat;refs=8;consumer=100%
  3. family=bag;stmt=filtered_derived;threshold=200
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=reuse;policy=mat;refs=8;consumer=1%
  6. family=reuse;policy=default;refs=4;consumer=1%
  7. family=bag;stmt=join_where;threshold=180000
  8. family=reuse;policy=mat;refs=16;consumer=100%
  9. family=bag;stmt=join_on;threshold=20000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  11. family=bag;stmt=cte_joined_outer;threshold=180000
  12. family=reuse;policy=mat;refs=2;consumer=1%
  13. family=bag;stmt=cte_filtered_notmat;threshold=20000
  14. family=reuse;policy=mat;refs=16;consumer=1%
  15. family=bag;stmt=cte_joined_outer;threshold=100000
  16. family=bag;stmt=join_where;threshold=2000
  17. family=bag;stmt=cte_filtered;threshold=100000
  18. family=reuse;policy=default;refs=16;consumer=1%
  19. family=reuse;policy=default;refs=1;consumer=100%
  20. family=bag;stmt=comma_join;threshold=2000
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  22. family=bag;stmt=join_where;threshold=20000
  23. family=reuse;policy=default;refs=2;consumer=100%
  24. family=bag;stmt=join_on;threshold=100000
  25. family=reuse;policy=notmat;refs=2;consumer=100%
  26. family=reuse;policy=notmat;refs=8;consumer=100%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  28. family=reuse;policy=notmat;refs=1;consumer=1%
  29. family=bag;stmt=cte_filtered_notmat;threshold=180000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  31. family=bag;stmt=filtered_derived;threshold=20000
  32. family=bag;stmt=comma_join;threshold=100000
  33. family=bag;stmt=cte_joined_outer;threshold=200
  34. family=reuse;policy=notmat;refs=1;consumer=100%
  35. family=bag;stmt=filtered_derived;threshold=2000
  36. family=bag;stmt=join_on;threshold=200
  37. family=reuse;policy=notmat;refs=8;consumer=1%
  38. family=bag;stmt=join_on;threshold=20
  39. family=bag;stmt=cte_filtered_mat;threshold=20
  40. family=reuse;policy=mat;refs=1;consumer=100%
  41. family=bag;stmt=cte_joined_outer;threshold=20000
  42. family=bag;stmt=join_on;threshold=180000
  43. family=bag;stmt=cte_filtered_mat;threshold=2000
  44. family=bag;stmt=join_where;threshold=20
  45. family=reuse;policy=default;refs=8;consumer=100%
  46. family=bag;stmt=cte_filtered_mat;threshold=20000
  47. family=reuse;policy=notmat;refs=16;consumer=100%
  48. family=bag;stmt=cte_filtered;threshold=200
  49. family=reuse;policy=mat;refs=4;consumer=100%
  50. family=bag;stmt=cte_joined_outer_mat;threshold=20
  51. family=bag;stmt=cte_filtered;threshold=20
  52. family=bag;stmt=filtered_derived;threshold=20
  53. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  54. family=reuse;policy=notmat;refs=4;consumer=100%
  55. family=bag;stmt=cte_filtered;threshold=20000
  56. family=reuse;policy=default;refs=16;consumer=100%
  57. family=bag;stmt=comma_join;threshold=20000
  58. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  59. family=reuse;policy=default;refs=8;consumer=1%
  60. family=reuse;policy=notmat;refs=4;consumer=1%
  61. family=bag;stmt=cte_joined_outer_mat;threshold=200
  62. family=reuse;policy=mat;refs=2;consumer=100%
  63. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  64. family=bag;stmt=comma_join;threshold=20
  65. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  66. family=bag;stmt=join_where;threshold=200
  67. family=bag;stmt=cte_filtered;threshold=2000
  68. family=bag;stmt=comma_join;threshold=200
  69. family=bag;stmt=join_where;threshold=100000
  70. family=reuse;policy=default;refs=2;consumer=1%
  71. family=bag;stmt=cte_filtered_notmat;threshold=200
  72. family=reuse;policy=notmat;refs=16;consumer=1%
  73. family=bag;stmt=cte_filtered_notmat;threshold=100000
  74. family=reuse;policy=mat;refs=1;consumer=1%
  75. family=bag;stmt=cte_filtered_mat;threshold=200
  76. family=bag;stmt=join_on;threshold=2000
  77. family=reuse;policy=default;refs=4;consumer=100%
  78. family=bag;stmt=cte_joined_outer;threshold=2000
  79. family=bag;stmt=filtered_derived;threshold=180000
  80. family=bag;stmt=cte_filtered_notmat;threshold=2000
  81. family=bag;stmt=cte_filtered_mat;threshold=180000
  82. family=bag;stmt=filtered_derived;threshold=100000
  83. family=bag;stmt=cte_filtered;threshold=180000
  84. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  85. family=reuse;policy=default;refs=1;consumer=1%
  86. family=reuse;policy=notmat;refs=2;consumer=1%
  87. family=bag;stmt=cte_filtered_notmat;threshold=20
  88. family=reuse;policy=mat;refs=4;consumer=1%
  89. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  90. family=bag;stmt=comma_join;threshold=180000
  91. family=bag;stmt=cte_filtered_mat;threshold=100000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `1420904bba94f0e6` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `74d69fbc863bc32f` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `e22f6cd93a59d19c` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `19baca70fba42e9a` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `7254d1afd1939ffc` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `2d79cda199927b4b` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `75401ae00d110fd5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `bf9cfabd1b4a1187` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `df8162a470f8086d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `f158d3219b39663d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `8e05eea5331a742a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `28ecf098f3d36bca` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `588172d0936daf28` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | j04m_cte_filtered_mat | valid | `f1d26aa32c74902d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `7eb06b86f1b16688` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `a6d33fa411f7f279` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `9b7c17fd322645fb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `2188321edf413942` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `c3db636a19865af4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `bab4e2565ae60cf8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `7aa2d43de01a2d41` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `9c2bfd00ee713541` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `be7dfd72fcf357e1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `5bd081ac8b56668a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `4c65843995fe231d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `4b14d63561e9b457` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `3d488a2193674b34` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ebd2b2747181c1f7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `39cb5cb9a6db4c94` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `64204fb9ef5df9a3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `55d0de96f98db4f8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `e5c00d41157aa87a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9cf454291f29297c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ecc5802156d85657` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `eeb0988ab21e1d41` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `bbd54cbf4391414a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `418f6c9c1c3409d0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `03a4179788362cac` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `14c1d2edf11192e7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `073778157edadb1e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `6ac5e602e4be0da1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `d8a773ef960ec71d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `bf7de12a511c2bff` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `58513fc3957662d6` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `01b503194db828cf` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `7941545cec47b3df` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `be0b9fc195d1611b` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `d88cf3fdb3ef3f46` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `0d2cf4ca9a9d0c14` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `7fbe900261d04483` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `a66fb83fa7bb529d` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `49a6d2ea3e546881` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `389682e9bb4d70bc` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `6db1667540284700` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `63d64d761bf12e65` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `58687c22d4dd2ddf` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `bfb02da22b5705ed` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `3eedc6f507a014f1` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `722df808f535356b` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `24a2045fa55fa95a` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `22cfcb2035e85d85` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `8ed4cf98543615d2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `036da7e3e98be779` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `08171b32ab32e67e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `06fc577b962275bc` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `cd8ab13ee2e31b8a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `c64f9bbd828d5c08` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `15d111d2c537806e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `1cec421a490be70d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `4d7d9c6570e762bd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `b79b4ac28e1dee8b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `329ce5f6f6cbe364` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a52550c160030633` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9081b0db7576b82a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `d9f17516fd69c3ce` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `b3a3ebc0b071eee8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `57a28938e63a984c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `98e235f909649301` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `1384b10b627c6210` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `88a95b60ff9c0825` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `4aece0680dd4db81` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `c9d956836c9d95db` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `337c3aaedab6f343` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e37bd4b030ba1c8d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ac688dab41a52008` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `ce45afe5bc55e3b3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `8f13a32b120dc8bb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `34b843591abbe534` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `506ab26b1348335c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `4d5de15f35b45d1b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `945fdca57bd345aa` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
