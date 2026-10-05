# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-b2-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `eee9aedad20795de40569abd06dda87c0e0280d29c5f1e953bc46100bd5a2306`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL; shared buffers and rows removed parsed per cell where printed
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: pg-single
- setting `fact_count`: 200000
- setting `index_arm`: pg-btree-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 2 CPU 3g
- setting `phase`: blended
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
| pg-single | 921 | 2 |

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
| 20 (0.01%) | 0.4696 | 0.6495 | 0.4575 | 0.4949 | 0.4873 | 0.6742 | 0.5445 | 45.5078 | 0.4843 | 0.5483 |
| 200 (0.1%) | 0.8139 | 0.7423 | 0.7752 | 0.7677 | 0.8287 | 0.8808 | 0.7947 | 41.8275 | 0.7546 | 0.7244 |
| 2000 (1%) | 1.9468 | 1.9598 | 2.0176 | 1.7193 | 2.3317 | 1.6827 | 2.1140 | 48.2378 | 1.9725 | 1.7484 |
| 20000 (10%) | 7.7079 | 6.9970 | 6.9336 | 6.9677 | 9.5616 | 7.1754 | 7.0109 | 49.9614 | 9.8826 | 7.2381 |
| 100000 (50%) | 49.5673 | 49.4591 | 49.1861 | 46.6222 | 63.4422 | 50.2679 | 50.5956 | 65.3221 | 49.9366 | 48.0317 |
| 180000 (90%) | 47.9526 | 83.2492 | 84.1793 | 82.4899 | 86.5141 | 47.4730 | 81.2971 | 89.8946 | 83.2593 | 82.6754 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4696 | 0.6495 | 1.3829 | true |
| 20 | filtered_derived vs join_where | 0.4696 | 0.4575 | 0.9741 | true |
| 20 | cte_filtered vs join_where | 0.4696 | 0.4949 | 1.0538 | true |
| 20 | cte_filtered_mat vs join_where | 0.4696 | 0.4873 | 1.0377 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4696 | 0.6742 | 1.4357 | true |
| 20 | cte_joined_outer vs join_where | 0.4696 | 0.5445 | 1.1594 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4696 | 45.5078 | 96.9024 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4696 | 0.4843 | 1.0313 | true |
| 20 | comma_join vs join_where | 0.4696 | 0.5483 | 1.1675 | true |
| 200 | join_on vs join_where | 0.8139 | 0.7423 | 0.9121 | true |
| 200 | filtered_derived vs join_where | 0.8139 | 0.7752 | 0.9525 | true |
| 200 | cte_filtered vs join_where | 0.8139 | 0.7677 | 0.9433 | true |
| 200 | cte_filtered_mat vs join_where | 0.8139 | 0.8287 | 1.0182 | false |
| 200 | cte_filtered_notmat vs join_where | 0.8139 | 0.8808 | 1.0822 | true |
| 200 | cte_joined_outer vs join_where | 0.8139 | 0.7947 | 0.9764 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.8139 | 41.8275 | 51.3941 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.8139 | 0.7546 | 0.9272 | true |
| 200 | comma_join vs join_where | 0.8139 | 0.7244 | 0.8901 | true |
| 2000 | join_on vs join_where | 1.9468 | 1.9598 | 1.0067 | true |
| 2000 | filtered_derived vs join_where | 1.9468 | 2.0176 | 1.0364 | true |
| 2000 | cte_filtered vs join_where | 1.9468 | 1.7193 | 0.8832 | true |
| 2000 | cte_filtered_mat vs join_where | 1.9468 | 2.3317 | 1.1977 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.9468 | 1.6827 | 0.8644 | true |
| 2000 | cte_joined_outer vs join_where | 1.9468 | 2.1140 | 1.0859 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.9468 | 48.2378 | 24.7783 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.9468 | 1.9725 | 1.0132 | true |
| 2000 | comma_join vs join_where | 1.9468 | 1.7484 | 0.8981 | true |
| 20000 | join_on vs join_where | 7.7079 | 6.9970 | 0.9078 | true |
| 20000 | filtered_derived vs join_where | 7.7079 | 6.9336 | 0.8996 | true |
| 20000 | cte_filtered vs join_where | 7.7079 | 6.9677 | 0.9040 | true |
| 20000 | cte_filtered_mat vs join_where | 7.7079 | 9.5616 | 1.2405 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.7079 | 7.1754 | 0.9309 | true |
| 20000 | cte_joined_outer vs join_where | 7.7079 | 7.0109 | 0.9096 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.7079 | 49.9614 | 6.4819 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.7079 | 9.8826 | 1.2821 | true |
| 20000 | comma_join vs join_where | 7.7079 | 7.2381 | 0.9391 | true |
| 100000 | join_on vs join_where | 49.5673 | 49.4591 | 0.9978 | true |
| 100000 | filtered_derived vs join_where | 49.5673 | 49.1861 | 0.9923 | true |
| 100000 | cte_filtered vs join_where | 49.5673 | 46.6222 | 0.9406 | true |
| 100000 | cte_filtered_mat vs join_where | 49.5673 | 63.4422 | 1.2799 | false |
| 100000 | cte_filtered_notmat vs join_where | 49.5673 | 50.2679 | 1.0141 | true |
| 100000 | cte_joined_outer vs join_where | 49.5673 | 50.5956 | 1.0207 | true |
| 100000 | cte_joined_outer_mat vs join_where | 49.5673 | 65.3221 | 1.3178 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 49.5673 | 49.9366 | 1.0075 | true |
| 100000 | comma_join vs join_where | 49.5673 | 48.0317 | 0.9690 | true |
| 180000 | join_on vs join_where | 47.9526 | 83.2492 | 1.7361 | true |
| 180000 | filtered_derived vs join_where | 47.9526 | 84.1793 | 1.7555 | true |
| 180000 | cte_filtered vs join_where | 47.9526 | 82.4899 | 1.7202 | true |
| 180000 | cte_filtered_mat vs join_where | 47.9526 | 86.5141 | 1.8042 | false |
| 180000 | cte_filtered_notmat vs join_where | 47.9526 | 47.4730 | 0.9900 | true |
| 180000 | cte_joined_outer vs join_where | 47.9526 | 81.2971 | 1.6954 | true |
| 180000 | cte_joined_outer_mat vs join_where | 47.9526 | 89.8946 | 1.8747 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 47.9526 | 83.2593 | 1.7363 | true |
| 180000 | comma_join vs join_where | 47.9526 | 82.6754 | 1.7241 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4949 | 0.4873 | 0.9847 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4949 | 0.6742 | 1.3624 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.5445 | 45.5078 | 83.5777 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.5445 | 0.4843 | 0.8895 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7677 | 0.8287 | 1.0794 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7677 | 0.8808 | 1.1473 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7947 | 41.8275 | 52.6347 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7947 | 0.7546 | 0.9496 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7193 | 2.3317 | 1.3562 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7193 | 1.6827 | 0.9787 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 2.1140 | 48.2378 | 22.8178 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 2.1140 | 1.9725 | 0.9330 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.9677 | 9.5616 | 1.3723 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.9677 | 7.1754 | 1.0298 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.0109 | 49.9614 | 7.1262 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.0109 | 9.8826 | 1.4096 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 46.6222 | 63.4422 | 1.3608 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 46.6222 | 50.2679 | 1.0782 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 50.5956 | 65.3221 | 1.2911 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 50.5956 | 49.9366 | 0.9870 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 82.4899 | 86.5141 | 1.0488 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 82.4899 | 47.4730 | 0.5755 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 81.2971 | 89.8946 | 1.1058 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 81.2971 | 83.2593 | 1.0241 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4208 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.4534 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.4884 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3819 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.5611 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.6043 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.7971 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.9761 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.4058 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.8540 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.3541 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4626 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.7366 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.3829 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.8130 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.6825 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.1038 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.0082 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.3993 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.6452 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.5110 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.1385 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.8286 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.5362 | 2 | `ea98347d41c49747` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8768 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.4318 | 4 | `55b877a44f5f2e35` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.5132 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.6938 | 8 | `d8980dd706cdce1b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.5556 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.9760 | 16 | `813a4381ac770391` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.5941 | false |
| notmat | 1 | 1% | 1.2142 | true |
| mat | 1 | 100% | 1.0038 | false |
| notmat | 1 | 100% | 0.8717 | true |
| mat | 2 | 1% | 1.0997 | true |
| notmat | 2 | 1% | 0.3330 | false |
| mat | 2 | 100% | 1.0004 | true |
| notmat | 2 | 100% | 1.9044 | false |
| mat | 4 | 1% | 1.0983 | true |
| notmat | 4 | 1% | 0.3424 | false |
| mat | 4 | 100% | 1.0300 | true |
| notmat | 4 | 100% | 2.8537 | false |
| mat | 8 | 1% | 1.1097 | true |
| notmat | 8 | 1% | 0.5410 | false |
| mat | 8 | 100% | 1.0108 | true |
| notmat | 8 | 100% | 3.2573 | false |
| mat | 16 | 1% | 0.9981 | true |
| notmat | 16 | 1% | 0.7504 | false |
| mat | 16 | 100% | 0.9458 | true |
| notmat | 16 | 100% | 3.8858 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.794676 | 0.681240 | 0.8573 | true |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING btree (charity_id)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING btree (person_id)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING btree (donation_id)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING btree (rank)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING btree (person_id)`

- index arm `pg-btree-rank`; placement `engine default`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 200: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 369, nr_throttled 290, throttled_usec 16327785 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 369, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=comma_join;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200;position=last | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 0 | 0 | 1286/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 0 | 0 | 1286/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 0 | 0 | 1286/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 0 | 0 | 1286/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 0 | 0 | 1286/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=join_on;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 34/0 | 1983 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 0 | 0 | 416/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 26/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 60/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 104/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 272/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 46.6337 | 48.0317 | 55.8311 | 5 | 20.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 76.1720 | 82.6754 | 94.9055 | 5 | 12.1 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.4451 | 0.5483 | 0.5864 | 5 | 1823.8 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6669 | 0.7244 | 0.7783 | 5 | 1380.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.6501 | 1.7484 | 2.1640 | 5 | 571.9 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.8472 | 7.2381 | 7.4826 | 5 | 138.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 43.0493 | 46.6222 | 59.3180 | 5 | 21.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 78.1199 | 82.4899 | 91.9688 | 5 | 12.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4210 | 0.4949 | 0.6404 | 5 | 2020.6 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6454 | 0.7677 | 0.7885 | 5 | 1302.6 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.6168 | 1.7193 | 1.9765 | 5 | 581.6 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.4756 | 6.9677 | 10.9851 | 5 | 143.5 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 61.9536 | 63.4422 | 64.0360 | 5 | 15.8 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 61.7028 | 86.5141 | 100.7446 | 5 | 11.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.3907 | 0.4873 | 0.5147 | 5 | 2051.9 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7730 | 0.8287 | 0.8982 | 5 | 1206.8 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.9178 | 2.3317 | 2.4793 | 5 | 428.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.9763 | 9.5616 | 10.1927 | 5 | 104.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 48.9008 | 50.2679 | 54.2516 | 5 | 19.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 42.2346 | 47.4730 | 51.8019 | 5 | 21.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.6261 | 0.6742 | 0.8463 | 5 | 1483.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.8264 | 0.8808 | 0.9102 | 5 | 1135.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.5330 | 1.6827 | 2.5213 | 5 | 594.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.8694 | 7.1754 | 7.3630 | 5 | 139.4 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 46.1116 | 50.5956 | 69.3893 | 5 | 19.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 66.4027 | 81.2971 | 90.8478 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4831 | 0.5445 | 0.5679 | 5 | 1836.6 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.6744 | 0.7947 | 0.8571 | 5 | 1258.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 2.0204 | 2.1140 | 2.5576 | 5 | 473.0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.7437 | 7.0109 | 34.9234 | 5 | 142.6 |
| family=bag;stmt=cte_joined_outer;threshold=200;position=last | valid | 0.5907 | 0.6812 | 0.7582 | 5 | 1467.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 62.6607 | 65.3221 | 66.7433 | 5 | 15.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 81.1710 | 89.8946 | 100.5071 | 5 | 11.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 42.4305 | 45.5078 | 47.5482 | 5 | 22.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 40.0839 | 41.8275 | 51.1462 | 5 | 23.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 44.3040 | 48.2378 | 54.4883 | 5 | 20.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 47.8449 | 49.9614 | 55.7043 | 5 | 20.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 47.4219 | 49.9366 | 72.1332 | 5 | 20.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 71.9031 | 83.2593 | 89.8443 | 5 | 12.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4509 | 0.4843 | 0.5576 | 5 | 2064.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6535 | 0.7546 | 0.7971 | 5 | 1325.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.8928 | 1.9725 | 2.1060 | 5 | 507.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.9375 | 9.8826 | 28.7584 | 5 | 101.2 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 46.6500 | 49.1861 | 51.2507 | 5 | 20.3 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 44.7886 | 84.1793 | 93.2287 | 5 | 11.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4100 | 0.4575 | 0.5073 | 5 | 2185.9 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.7430 | 0.7752 | 0.8339 | 5 | 1290.0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.9608 | 2.0176 | 2.0763 | 5 | 495.6 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.4476 | 6.9336 | 8.3110 | 5 | 144.2 |
| family=bag;stmt=join_on;threshold=100000 | valid | 45.7530 | 49.4591 | 52.2552 | 5 | 20.2 |
| family=bag;stmt=join_on;threshold=180000 | valid | 69.3951 | 83.2492 | 97.7960 | 5 | 12.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.6282 | 0.6495 | 0.7007 | 5 | 1539.8 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.6473 | 0.7423 | 0.9062 | 5 | 1347.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.9067 | 1.9598 | 2.1789 | 5 | 510.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.5336 | 6.9970 | 32.4885 | 5 | 142.9 |
| family=bag;stmt=join_where;threshold=100000 | valid | 48.5549 | 49.5673 | 50.9427 | 5 | 20.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 44.8905 | 47.9526 | 52.0981 | 5 | 20.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.4071 | 0.4696 | 0.5476 | 5 | 2129.4 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6473 | 0.8139 | 1.0524 | 5 | 1228.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.8935 | 1.9468 | 2.3628 | 5 | 513.7 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.5619 | 7.7079 | 8.5381 | 5 | 129.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.2980 | 3.4058 | 3.4517 | 5 | 293.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.5489 | 3.8540 | 3.8777 | 5 | 259.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.3669 | 0.4208 | 0.4606 | 5 | 2376.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.3345 | 2.4534 | 2.6554 | 5 | 407.6 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.3740 | 2.4884 | 2.6266 | 5 | 401.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.2969 | 2.3819 | 2.5590 | 5 | 419.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.5484 | 2.5611 | 2.6560 | 5 | 390.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.5572 | 2.6043 | 2.6346 | 5 | 384.0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.7020 | 2.7971 | 3.1940 | 5 | 357.5 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.9082 | 2.9761 | 3.1040 | 5 | 336.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.2577 | 3.3993 | 3.6017 | 5 | 294.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.3942 | 3.6452 | 3.9014 | 5 | 274.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.2484 | 2.3541 | 2.4437 | 5 | 424.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.3774 | 2.4626 | 2.5490 | 5 | 406.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.6479 | 2.7366 | 2.7877 | 5 | 365.4 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.2922 | 2.3829 | 2.6350 | 5 | 419.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.7639 | 2.8130 | 2.9394 | 5 | 355.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.6298 | 2.6825 | 2.8481 | 5 | 372.8 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 3.0721 | 3.1038 | 3.3018 | 5 | 322.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.8344 | 3.0082 | 3.1482 | 5 | 332.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.4443 | 2.5556 | 3.3798 | 5 | 391.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.8425 | 14.9760 | 15.7777 | 5 | 66.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4568 | 0.5110 | 6.2877 | 5 | 1957.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.1052 | 2.1385 | 3.5635 | 5 | 467.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.5892 | 0.8286 | 1.2001 | 5 | 1206.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.4248 | 4.5362 | 5.3995 | 5 | 220.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8397 | 0.8768 | 0.9289 | 5 | 1140.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.9090 | 7.4318 | 7.6564 | 5 | 134.6 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3005 | 1.5132 | 1.5709 | 5 | 660.9 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 9.4913 | 9.6938 | 10.9841 | 5 | 103.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 340, committed 340, failed 0 over 33.935 s
- offered rate 10.019/s, completed rate 10.019/s
- writer operations per endpoint: [340]
- reader operations overlapping a committed write: 131 of 637
- reader retries of transient conflicts: 0
- per key: 340 keys, min 1, max 1 per key
- ledger: 340 keys reconciled, 0 mismatch
- update latency: avg 2.683 ms, min 0.812, max 51.216 over 340

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42096721`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=200
  2. family=bag;stmt=join_where;threshold=200
  3. family=bag;stmt=filtered_derived;threshold=20
  4. family=bag;stmt=cte_filtered_notmat;threshold=20000
  5. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  6. family=bag;stmt=join_on;threshold=20000
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=comma_join;threshold=2000
  9. family=bag;stmt=cte_filtered;threshold=200
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=bag;stmt=cte_joined_outer;threshold=20
  12. family=bag;stmt=join_where;threshold=180000
  13. family=bag;stmt=cte_filtered_notmat;threshold=20
  14. family=bag;stmt=comma_join;threshold=180000
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  16. family=bag;stmt=filtered_derived;threshold=20000
  17. family=reuse;policy=notmat;refs=8;consumer=100%
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=bag;stmt=cte_filtered_notmat;threshold=100000
  20. family=bag;stmt=cte_filtered_mat;threshold=20000
  21. family=reuse;policy=mat;refs=8;consumer=1%
  22. family=bag;stmt=filtered_derived;threshold=200
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  25. family=bag;stmt=join_on;threshold=2000
  26. family=reuse;policy=mat;refs=1;consumer=1%
  27. family=reuse;policy=notmat;refs=2;consumer=100%
  28. family=reuse;policy=notmat;refs=1;consumer=1%
  29. family=bag;stmt=cte_filtered;threshold=20
  30. family=reuse;policy=notmat;refs=4;consumer=1%
  31. family=reuse;policy=mat;refs=4;consumer=100%
  32. family=reuse;policy=default;refs=4;consumer=100%
  33. family=bag;stmt=join_where;threshold=100000
  34. family=bag;stmt=cte_joined_outer;threshold=180000
  35. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  36. family=bag;stmt=join_on;threshold=100000
  37. family=reuse;policy=notmat;refs=16;consumer=1%
  38. family=reuse;policy=default;refs=4;consumer=1%
  39. family=reuse;policy=mat;refs=16;consumer=100%
  40. family=reuse;policy=notmat;refs=1;consumer=100%
  41. family=bag;stmt=cte_filtered_mat;threshold=20
  42. family=reuse;policy=default;refs=8;consumer=100%
  43. family=reuse;policy=notmat;refs=16;consumer=100%
  44. family=bag;stmt=filtered_derived;threshold=100000
  45. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  46. family=bag;stmt=cte_joined_outer;threshold=100000
  47. family=reuse;policy=mat;refs=2;consumer=1%
  48. family=bag;stmt=cte_filtered;threshold=2000
  49. family=bag;stmt=join_on;threshold=200
  50. family=reuse;policy=default;refs=2;consumer=1%
  51. family=reuse;policy=default;refs=2;consumer=100%
  52. family=bag;stmt=join_where;threshold=20
  53. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  54. family=reuse;policy=default;refs=8;consumer=1%
  55. family=reuse;policy=notmat;refs=4;consumer=100%
  56. family=bag;stmt=cte_filtered;threshold=100000
  57. family=reuse;policy=default;refs=1;consumer=100%
  58. family=bag;stmt=cte_joined_outer_mat;threshold=20
  59. family=bag;stmt=cte_filtered_mat;threshold=200
  60. family=bag;stmt=cte_filtered_mat;threshold=2000
  61. family=bag;stmt=cte_filtered_notmat;threshold=180000
  62. family=bag;stmt=cte_filtered;threshold=180000
  63. family=reuse;policy=notmat;refs=2;consumer=1%
  64. family=bag;stmt=comma_join;threshold=20
  65. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  66. family=bag;stmt=filtered_derived;threshold=180000
  67. family=bag;stmt=comma_join;threshold=20000
  68. family=bag;stmt=join_where;threshold=20000
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  70. family=reuse;policy=default;refs=16;consumer=1%
  71. family=bag;stmt=comma_join;threshold=200
  72. family=bag;stmt=cte_filtered;threshold=20000
  73. family=bag;stmt=join_where;threshold=2000
  74. family=bag;stmt=cte_filtered_mat;threshold=100000
  75. family=bag;stmt=cte_filtered_mat;threshold=180000
  76. family=bag;stmt=cte_joined_outer;threshold=2000
  77. family=bag;stmt=cte_filtered_notmat;threshold=2000
  78. family=reuse;policy=mat;refs=2;consumer=100%
  79. family=bag;stmt=cte_joined_outer;threshold=20000
  80. family=reuse;policy=default;refs=16;consumer=100%
  81. family=reuse;policy=mat;refs=1;consumer=100%
  82. family=bag;stmt=comma_join;threshold=100000
  83. family=bag;stmt=join_on;threshold=20
  84. family=bag;stmt=join_on;threshold=180000
  85. family=bag;stmt=filtered_derived;threshold=2000
  86. family=bag;stmt=cte_joined_outer_mat;threshold=200
  87. family=bag;stmt=cte_filtered_notmat;threshold=200
  88. family=reuse;policy=notmat;refs=8;consumer=1%
  89. family=reuse;policy=default;refs=1;consumer=1%
  90. family=reuse;policy=mat;refs=8;consumer=100%
  91. family=bag;stmt=cte_joined_outer;threshold=200;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `4bb1561e58681728` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `41f23de3437828d5` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `5573434208edb047` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `09a028d269057a3e` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `4b4ce4c83ae7fcc8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `6a10a6fd7e521978` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `d4ca4a9d1c0fd3e3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `a398c2b3383a151d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `d4df6a0ef2f05082` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `5db776896b6f593f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200;position=last | j05_cte_joined_outer | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `4bb1561e58681728` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `41f23de3437828d5` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `5573434208edb047` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `5573434208edb047` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5573434208edb047` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `09a028d269057a3e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `549a07108ab01ba9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `813a4381ac770391` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `ea98347d41c49747` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `55b877a44f5f2e35` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `d8980dd706cdce1b` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
