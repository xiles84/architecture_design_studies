# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-i2-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `6da46b4f1e2f9b5782d3d345e55a8b8f87df08804be165c472414a92cdefedb0`
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
| 20 (0.01%) | 1.5563 | 2.0419 | 2.1162 | 1.9058 | 2.0027 | 1.9788 | 1.7806 | 86.6027 | 1.9109 | 2.1732 |
| 200 (0.1%) | 3.8137 | 3.0523 | 2.9990 | 3.6271 | 2.9714 | 3.5907 | 3.4536 | 86.1154 | 3.4811 | 3.8043 |
| 2000 (1%) | 9.6494 | 10.4809 | 9.4266 | 8.9177 | 9.3338 | 9.6720 | 9.8780 | 85.4882 | 9.4694 | 9.4514 |
| 20000 (10%) | 65.5114 | 66.8211 | 64.2108 | 66.6370 | 68.6005 | 67.3947 | 66.8383 | 108.7087 | 67.8046 | 70.8262 |
| 100000 (50%) | 87.5344 | 93.8932 | 91.1028 | 90.6630 | 92.1453 | 87.1437 | 88.7379 | 98.7858 | 85.3816 | 89.3394 |
| 180000 (90%) | 110.6720 | 119.6913 | 109.8656 | 109.4802 | 117.1653 | 106.7813 | 107.2786 | 102.8562 | 111.0756 | 105.3497 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.5563 | 2.0419 | 1.3120 | false |
| 20 | filtered_derived vs join_where | 1.5563 | 2.1162 | 1.3597 | false |
| 20 | cte_filtered vs join_where | 1.5563 | 1.9058 | 1.2246 | false |
| 20 | cte_filtered_mat vs join_where | 1.5563 | 2.0027 | 1.2868 | false |
| 20 | cte_filtered_notmat vs join_where | 1.5563 | 1.9788 | 1.2715 | false |
| 20 | cte_joined_outer vs join_where | 1.5563 | 1.7806 | 1.1441 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.5563 | 86.6027 | 55.6465 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.5563 | 1.9109 | 1.2278 | false |
| 20 | comma_join vs join_where | 1.5563 | 2.1732 | 1.3964 | false |
| 200 | join_on vs join_where | 3.8137 | 3.0523 | 0.8003 | false |
| 200 | filtered_derived vs join_where | 3.8137 | 2.9990 | 0.7864 | false |
| 200 | cte_filtered vs join_where | 3.8137 | 3.6271 | 0.9511 | false |
| 200 | cte_filtered_mat vs join_where | 3.8137 | 2.9714 | 0.7791 | false |
| 200 | cte_filtered_notmat vs join_where | 3.8137 | 3.5907 | 0.9415 | false |
| 200 | cte_joined_outer vs join_where | 3.8137 | 3.4536 | 0.9056 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.8137 | 86.1154 | 22.5803 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.8137 | 3.4811 | 0.9128 | false |
| 200 | comma_join vs join_where | 3.8137 | 3.8043 | 0.9975 | false |
| 2000 | join_on vs join_where | 9.6494 | 10.4809 | 1.0862 | false |
| 2000 | filtered_derived vs join_where | 9.6494 | 9.4266 | 0.9769 | false |
| 2000 | cte_filtered vs join_where | 9.6494 | 8.9177 | 0.9242 | false |
| 2000 | cte_filtered_mat vs join_where | 9.6494 | 9.3338 | 0.9673 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.6494 | 9.6720 | 1.0023 | false |
| 2000 | cte_joined_outer vs join_where | 9.6494 | 9.8780 | 1.0237 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6494 | 85.4882 | 8.8595 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6494 | 9.4694 | 0.9813 | false |
| 2000 | comma_join vs join_where | 9.6494 | 9.4514 | 0.9795 | false |
| 20000 | join_on vs join_where | 65.5114 | 66.8211 | 1.0200 | false |
| 20000 | filtered_derived vs join_where | 65.5114 | 64.2108 | 0.9801 | false |
| 20000 | cte_filtered vs join_where | 65.5114 | 66.6370 | 1.0172 | false |
| 20000 | cte_filtered_mat vs join_where | 65.5114 | 68.6005 | 1.0472 | false |
| 20000 | cte_filtered_notmat vs join_where | 65.5114 | 67.3947 | 1.0287 | false |
| 20000 | cte_joined_outer vs join_where | 65.5114 | 66.8383 | 1.0203 | false |
| 20000 | cte_joined_outer_mat vs join_where | 65.5114 | 108.7087 | 1.6594 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 65.5114 | 67.8046 | 1.0350 | false |
| 20000 | comma_join vs join_where | 65.5114 | 70.8262 | 1.0811 | false |
| 100000 | join_on vs join_where | 87.5344 | 93.8932 | 1.0726 | false |
| 100000 | filtered_derived vs join_where | 87.5344 | 91.1028 | 1.0408 | false |
| 100000 | cte_filtered vs join_where | 87.5344 | 90.6630 | 1.0357 | false |
| 100000 | cte_filtered_mat vs join_where | 87.5344 | 92.1453 | 1.0527 | false |
| 100000 | cte_filtered_notmat vs join_where | 87.5344 | 87.1437 | 0.9955 | false |
| 100000 | cte_joined_outer vs join_where | 87.5344 | 88.7379 | 1.0137 | false |
| 100000 | cte_joined_outer_mat vs join_where | 87.5344 | 98.7858 | 1.1285 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 87.5344 | 85.3816 | 0.9754 | false |
| 100000 | comma_join vs join_where | 87.5344 | 89.3394 | 1.0206 | false |
| 180000 | join_on vs join_where | 110.6720 | 119.6913 | 1.0815 | false |
| 180000 | filtered_derived vs join_where | 110.6720 | 109.8656 | 0.9927 | false |
| 180000 | cte_filtered vs join_where | 110.6720 | 109.4802 | 0.9892 | false |
| 180000 | cte_filtered_mat vs join_where | 110.6720 | 117.1653 | 1.0587 | false |
| 180000 | cte_filtered_notmat vs join_where | 110.6720 | 106.7813 | 0.9648 | false |
| 180000 | cte_joined_outer vs join_where | 110.6720 | 107.2786 | 0.9693 | false |
| 180000 | cte_joined_outer_mat vs join_where | 110.6720 | 102.8562 | 0.9294 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 110.6720 | 111.0756 | 1.0036 | false |
| 180000 | comma_join vs join_where | 110.6720 | 105.3497 | 0.9519 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.9058 | 2.0027 | 1.0509 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.9058 | 1.9788 | 1.0383 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.7806 | 86.6027 | 48.6372 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.7806 | 1.9109 | 1.0732 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.6271 | 2.9714 | 0.8192 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.6271 | 3.5907 | 0.9899 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.4536 | 86.1154 | 24.9350 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.4536 | 3.4811 | 1.0080 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.9177 | 9.3338 | 1.0467 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.9177 | 9.6720 | 1.0846 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.8780 | 85.4882 | 8.6544 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.8780 | 9.4694 | 0.9586 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 66.6370 | 68.6005 | 1.0295 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 66.6370 | 67.3947 | 1.0114 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 66.8383 | 108.7087 | 1.6264 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 66.8383 | 67.8046 | 1.0145 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 90.6630 | 92.1453 | 1.0163 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 90.6630 | 87.1437 | 0.9612 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 88.7379 | 98.7858 | 1.1132 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 88.7379 | 85.3816 | 0.9622 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 109.4802 | 117.1653 | 1.0702 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 109.4802 | 106.7813 | 0.9753 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 107.2786 | 102.8562 | 0.9588 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 107.2786 | 111.0756 | 1.0354 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.9661 | 1 | `dfa1d96b873aca6c` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.8803 | 1 | `4139211825588cd1` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 9.1840 | 2 | `b6729d1b0604e1b4` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 9.2138 | 2 | `2b62b395ae587b82` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 9.2588 | 4 | `b40c1f2e26d8a06f` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.2110 | 4 | `f712359a3084c3b9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.8404 | 8 | `c621222deda8b7c8` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.6147 | 8 | `9b421fcd3ba03dcc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.0276 | 16 | `19994fa5f00c2beb` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.4865 | 16 | `0a013b03c7a8c65b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.1452 | 1 | `8691134162a40cb4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.4285 | 1 | `8734b8f410b9555a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.7867 | 2 | `87354a102df5cdf0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.9510 | 2 | `0e363e37f2a4029a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 9.2661 | 4 | `2faf8d5f9682057e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.1258 | 4 | `72aaf78bc42bfe7d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.7731 | 8 | `330bbbc118076618` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.2261 | 8 | `8887b5846a50c727` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.6595 | 16 | `88f29bf556c17e95` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.5307 | 16 | `6f5e9eb71dc47b10` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.1228 | 1 | `2467ef064158e039` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.9821 | 1 | `9f6bdf043e701982` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.8306 | 2 | `65b7482a24cf23e4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 16.5561 | 2 | `55873739a4dfcf6c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 25.7260 | 4 | `34343bde060174d5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 32.6066 | 4 | `e5d8c94a9fff3555` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 56.0457 | 8 | `c375095cd92842da` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 63.1631 | 8 | `afb3bcf8e7dba648` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 109.3732 | 16 | `313ed902889cdf2d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 112.5226 | 16 | `b3709fdca6cd1051` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1693 | false |
| notmat | 1 | 1% | 1.0225 | false |
| mat | 1 | 100% | 0.9491 | false |
| notmat | 1 | 100% | 1.0115 | false |
| mat | 2 | 1% | 0.9567 | false |
| notmat | 2 | 1% | 1.3971 | false |
| mat | 2 | 100% | 0.9715 | false |
| notmat | 2 | 100% | 1.7969 | false |
| mat | 4 | 1% | 1.0008 | false |
| notmat | 4 | 1% | 2.7785 | false |
| mat | 4 | 100% | 0.9908 | false |
| notmat | 4 | 100% | 3.5400 | false |
| mat | 8 | 1% | 0.9924 | false |
| notmat | 8 | 1% | 6.3398 | false |
| mat | 8 | 100% | 0.9596 | false |
| notmat | 8 | 100% | 6.5694 | false |
| mat | 16 | 1% | 1.0700 | false |
| notmat | 16 | 1% | 12.1154 | false |
| mat | 16 | 100% | 0.9088 | false |
| notmat | 16 | 100% | 10.7302 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 86.115421 | 93.670573 | 1.0877 | false |

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

- client: nr_periods 626, nr_throttled 239, throttled_usec 9067213 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 699, nr_throttled 83, throttled_usec 3005333 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | valid | false | true | 202050 | 0 | 0/0 | 199800 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 86.3503 | 89.3394 | 91.7094 | 5 | 11.2 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 103.6140 | 105.3497 | 110.3754 | 5 | 9.5 |
| family=bag;stmt=comma_join;threshold=20 | valid | 2.0561 | 2.1732 | 2.4297 | 5 | 460.2 |
| family=bag;stmt=comma_join;threshold=200 | valid | 3.2846 | 3.8043 | 3.8434 | 5 | 262.9 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 9.0749 | 9.4514 | 38.6641 | 5 | 105.8 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 65.6872 | 70.8262 | 73.9805 | 5 | 14.1 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 87.1560 | 90.6630 | 95.2916 | 5 | 11.0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 106.0324 | 109.4802 | 115.1599 | 5 | 9.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.4919 | 1.9058 | 2.1740 | 5 | 524.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.3265 | 3.6271 | 4.0058 | 5 | 275.7 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.6113 | 8.9177 | 9.5997 | 5 | 112.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 63.1083 | 66.6370 | 74.6185 | 5 | 15.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 87.6343 | 92.1453 | 95.4931 | 5 | 10.9 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 108.3034 | 117.1653 | 120.8436 | 5 | 8.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.8591 | 2.0027 | 2.1653 | 5 | 499.3 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.7424 | 2.9714 | 3.1334 | 5 | 336.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.0742 | 9.3338 | 9.7486 | 5 | 107.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 67.8490 | 68.6005 | 70.3120 | 5 | 14.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 85.4484 | 87.1437 | 95.5688 | 5 | 11.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 105.9326 | 106.7813 | 109.1565 | 5 | 9.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.7951 | 1.9788 | 2.3921 | 5 | 505.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.5074 | 3.5907 | 4.1715 | 5 | 278.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.8659 | 9.6720 | 10.4579 | 5 | 103.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 63.1591 | 67.3947 | 71.0376 | 5 | 14.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 88.5529 | 88.7379 | 90.7225 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 101.8681 | 107.2786 | 114.0239 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6525 | 1.7806 | 1.9323 | 5 | 561.6 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.9657 | 3.4536 | 17.5835 | 5 | 289.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.3503 | 9.8780 | 10.6581 | 5 | 101.2 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 66.0063 | 66.8383 | 70.6298 | 5 | 15.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 96.0181 | 98.7858 | 117.9981 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.6215 | 102.8562 | 114.6735 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 81.9411 | 86.6027 | 88.3481 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 85.3588 | 86.1154 | 92.9246 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.2480 | 85.4882 | 97.0592 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 91.1163 | 108.7087 | 129.6968 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | valid | 91.7152 | 93.6706 | 95.6495 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 84.2688 | 85.3816 | 86.9914 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 105.0022 | 111.0756 | 116.9917 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6893 | 1.9109 | 24.4910 | 5 | 523.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.9426 | 3.4811 | 3.6484 | 5 | 287.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.2470 | 9.4694 | 10.0770 | 5 | 105.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 65.4555 | 67.8046 | 73.5510 | 5 | 14.7 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 86.8031 | 91.1028 | 92.9902 | 5 | 11.0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 106.2220 | 109.8656 | 123.4240 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 2.0178 | 2.1162 | 2.4248 | 5 | 472.6 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.7624 | 2.9990 | 3.7384 | 5 | 333.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.7660 | 9.4266 | 10.4127 | 5 | 106.1 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 61.6362 | 64.2108 | 74.1805 | 5 | 15.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 88.2508 | 93.8932 | 98.8883 | 5 | 10.7 |
| family=bag;stmt=join_on;threshold=180000 | valid | 104.4187 | 119.6913 | 122.3478 | 5 | 8.4 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.5815 | 2.0419 | 2.1452 | 5 | 489.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.9959 | 3.0523 | 3.1141 | 5 | 327.6 |
| family=bag;stmt=join_on;threshold=2000 | valid | 9.8796 | 10.4809 | 10.6316 | 5 | 95.4 |
| family=bag;stmt=join_on;threshold=20000 | valid | 63.7782 | 66.8211 | 72.6767 | 5 | 15.0 |
| family=bag;stmt=join_where;threshold=100000 | valid | 86.8729 | 87.5344 | 89.2993 | 5 | 11.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 99.5708 | 110.6720 | 122.4848 | 5 | 9.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.3723 | 1.5563 | 1.8846 | 5 | 642.5 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.0889 | 3.8137 | 3.9736 | 5 | 262.2 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3480 | 9.6494 | 10.1417 | 5 | 103.6 |
| family=bag;stmt=join_where;threshold=20000 | valid | 62.3795 | 65.5114 | 68.5928 | 5 | 15.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.5853 | 9.0276 | 9.6243 | 5 | 110.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.9052 | 10.4865 | 11.1313 | 5 | 95.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.5759 | 6.9661 | 7.4069 | 5 | 143.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.5452 | 8.8803 | 9.5268 | 5 | 112.6 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.7687 | 9.1840 | 9.4781 | 5 | 108.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.5532 | 9.2138 | 9.9404 | 5 | 108.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.5028 | 9.2588 | 9.5797 | 5 | 108.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.0337 | 9.2110 | 10.9555 | 5 | 108.6 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.6103 | 8.8404 | 9.1911 | 5 | 113.1 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 9.1324 | 9.6147 | 10.2220 | 5 | 104.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.4021 | 9.6595 | 11.5296 | 5 | 103.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.2718 | 9.5307 | 10.5738 | 5 | 104.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.5011 | 8.1452 | 8.7356 | 5 | 122.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.0204 | 8.4285 | 9.5066 | 5 | 118.6 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.7494 | 8.7867 | 9.5259 | 5 | 113.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.4776 | 8.9510 | 9.4645 | 5 | 111.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.7744 | 9.2661 | 10.2609 | 5 | 107.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.8408 | 9.1258 | 9.7530 | 5 | 109.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.5007 | 8.7731 | 10.2872 | 5 | 114.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.0160 | 9.2261 | 9.8473 | 5 | 108.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.5397 | 109.3732 | 137.6102 | 5 | 9.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.2744 | 112.5226 | 133.1600 | 5 | 8.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.9854 | 7.1228 | 7.1466 | 5 | 140.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.2672 | 8.9821 | 9.4817 | 5 | 111.3 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 12.4442 | 12.8306 | 13.6036 | 5 | 77.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 15.5712 | 16.5561 | 17.1257 | 5 | 60.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 24.1276 | 25.7260 | 28.1393 | 5 | 38.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 31.4925 | 32.6066 | 33.9688 | 5 | 30.7 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 55.3178 | 56.0457 | 58.3403 | 5 | 17.8 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 60.6004 | 63.1631 | 64.1780 | 5 | 15.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42728134`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=200
  2. family=bag;stmt=filtered_derived;threshold=180000
  3. family=bag;stmt=join_where;threshold=100000
  4. family=reuse;policy=notmat;refs=16;consumer=1%
  5. family=bag;stmt=join_where;threshold=200
  6. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  7. family=bag;stmt=comma_join;threshold=20000
  8. family=bag;stmt=cte_joined_outer_mat;threshold=20
  9. family=bag;stmt=join_on;threshold=100000
  10. family=reuse;policy=default;refs=4;consumer=100%
  11. family=reuse;policy=mat;refs=16;consumer=100%
  12. family=bag;stmt=cte_filtered;threshold=20000
  13. family=bag;stmt=cte_filtered_notmat;threshold=20000
  14. family=reuse;policy=default;refs=16;consumer=100%
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  16. family=reuse;policy=mat;refs=8;consumer=1%
  17. family=bag;stmt=filtered_derived;threshold=20000
  18. family=bag;stmt=join_on;threshold=2000
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=bag;stmt=cte_filtered_notmat;threshold=20
  21. family=bag;stmt=filtered_derived;threshold=200
  22. family=bag;stmt=cte_filtered;threshold=20
  23. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  24. family=reuse;policy=default;refs=2;consumer=1%
  25. family=bag;stmt=cte_filtered_notmat;threshold=180000
  26. family=bag;stmt=comma_join;threshold=200
  27. family=reuse;policy=default;refs=1;consumer=1%
  28. family=bag;stmt=cte_joined_outer;threshold=180000
  29. family=reuse;policy=notmat;refs=8;consumer=100%
  30. family=bag;stmt=cte_filtered;threshold=200
  31. family=bag;stmt=comma_join;threshold=2000
  32. family=reuse;policy=mat;refs=16;consumer=1%
  33. family=bag;stmt=cte_joined_outer;threshold=200
  34. family=reuse;policy=notmat;refs=4;consumer=100%
  35. family=bag;stmt=cte_filtered_mat;threshold=20000
  36. family=reuse;policy=notmat;refs=16;consumer=100%
  37. family=bag;stmt=cte_filtered;threshold=180000
  38. family=reuse;policy=default;refs=2;consumer=100%
  39. family=bag;stmt=cte_filtered_mat;threshold=200
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  41. family=bag;stmt=filtered_derived;threshold=20
  42. family=bag;stmt=cte_filtered_notmat;threshold=2000
  43. family=reuse;policy=mat;refs=4;consumer=100%
  44. family=bag;stmt=filtered_derived;threshold=2000
  45. family=reuse;policy=mat;refs=1;consumer=100%
  46. family=bag;stmt=join_on;threshold=200
  47. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  48. family=bag;stmt=join_where;threshold=180000
  49. family=bag;stmt=cte_filtered_mat;threshold=100000
  50. family=bag;stmt=cte_joined_outer;threshold=100000
  51. family=bag;stmt=cte_filtered_notmat;threshold=200
  52. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  53. family=bag;stmt=comma_join;threshold=100000
  54. family=reuse;policy=notmat;refs=2;consumer=1%
  55. family=bag;stmt=comma_join;threshold=180000
  56. family=reuse;policy=notmat;refs=1;consumer=1%
  57. family=bag;stmt=join_on;threshold=20
  58. family=bag;stmt=cte_filtered_mat;threshold=2000
  59. family=reuse;policy=default;refs=16;consumer=1%
  60. family=reuse;policy=notmat;refs=8;consumer=1%
  61. family=reuse;policy=default;refs=8;consumer=100%
  62. family=bag;stmt=cte_filtered_mat;threshold=20
  63. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  64. family=bag;stmt=join_on;threshold=20000
  65. family=reuse;policy=mat;refs=2;consumer=100%
  66. family=bag;stmt=cte_filtered;threshold=100000
  67. family=reuse;policy=mat;refs=4;consumer=1%
  68. family=bag;stmt=filtered_derived;threshold=100000
  69. family=bag;stmt=cte_joined_outer;threshold=20000
  70. family=bag;stmt=comma_join;threshold=20
  71. family=bag;stmt=join_where;threshold=20
  72. family=bag;stmt=join_where;threshold=2000
  73. family=reuse;policy=notmat;refs=1;consumer=100%
  74. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  75. family=bag;stmt=cte_joined_outer;threshold=2000
  76. family=reuse;policy=mat;refs=8;consumer=100%
  77. family=reuse;policy=default;refs=4;consumer=1%
  78. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  79. family=reuse;policy=mat;refs=2;consumer=1%
  80. family=reuse;policy=default;refs=8;consumer=1%
  81. family=bag;stmt=cte_joined_outer;threshold=20
  82. family=bag;stmt=join_where;threshold=20000
  83. family=reuse;policy=notmat;refs=2;consumer=100%
  84. family=bag;stmt=cte_filtered;threshold=2000
  85. family=bag;stmt=cte_filtered_mat;threshold=180000
  86. family=reuse;policy=default;refs=1;consumer=100%
  87. family=reuse;policy=mat;refs=1;consumer=1%
  88. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  89. family=bag;stmt=cte_filtered_notmat;threshold=100000
  90. family=bag;stmt=join_on;threshold=180000
  91. family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `fcf4b9c69b1a52bb` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `8b1e2e5842c2cb54` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `9197fdcf4115b6f0` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `7a35777e40bfc4d1` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `f0ccb016766b37e4` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `878c756ba3c65383` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `cd6d71c8f2cced0f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `dd32a1a004b4c456` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `93edaba6bcd201f9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `532ab385feaef55f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `73c51b3adcbcf17a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `296a7657fe46a55e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `84dd60f095507605` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `e20abfbbffcad95a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `062d00a9f91c3d71` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `d9b632581534a492` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `7631e23d64ed1aa8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `617b2e0217b2685a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `5df707c50c7b1a15` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `21f2d860be31c61e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `759bd76960445cbe` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `f2aaaa24081f04ec` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `c670eadda67b7029` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `9dce10a4f1064398` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `c0bc8e767fd5a834` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `1ba260cb22230dc0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d65143136b44041f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `74c745958566b4a3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `578108b3463a8482` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `54c6cb792623279f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `7f6e7f021243e287` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `b1e52aba26901e07` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b1011a808987bc91` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `ba3833e7ed4a0186` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `0176f586882cc4d2` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `1ccd033a9f30b001` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | j05m_cte_joined_outer_mat | valid | `69af17efcdf14067` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `5745a927c226d54a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `0b8483a8a5608562` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `93a9cdf253364d91` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `b7902a09aaf423c9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `95293b946e7c67c4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `5a6a88a0afc34cd9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `ead1cddd55f95349` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `f643cc135f70cc12` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `bb6b8aa3f2d3a4cd` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `749ef7c81fc76675` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `75440d8c1e0f1411` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `63df59198454df14` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `0f2d4617656a345a` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `7ac58f3053247ebc` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `a6a0ad94f3b6d9fa` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `42cbbe8dd30adb6f` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `e4029722008ef5d2` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `d1dab5895e0eb659` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `1fe1aabb017564c5` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `5c2b4fb1948c02a5` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `e1bac74ee4c73d7f` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `5417065a2fca8e4d` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `50f23da51d84191d` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `2053da40891a34c4` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `19994fa5f00c2beb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `0a013b03c7a8c65b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `dfa1d96b873aca6c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `4139211825588cd1` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `b6729d1b0604e1b4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `2b62b395ae587b82` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `b40c1f2e26d8a06f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f712359a3084c3b9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `c621222deda8b7c8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `9b421fcd3ba03dcc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `88f29bf556c17e95` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `6f5e9eb71dc47b10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `8691134162a40cb4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `8734b8f410b9555a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `87354a102df5cdf0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `0e363e37f2a4029a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2faf8d5f9682057e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `72aaf78bc42bfe7d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `330bbbc118076618` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `8887b5846a50c727` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `313ed902889cdf2d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `b3709fdca6cd1051` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2467ef064158e039` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9f6bdf043e701982` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `65b7482a24cf23e4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `55873739a4dfcf6c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `34343bde060174d5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `e5d8c94a9fff3555` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `c375095cd92842da` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `afb3bcf8e7dba648` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
