# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-smoke-20261004`
- run tag: `run/study-01/select-join-topo-yb-smoke-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `e1311d7bb951f79c2c41334384d3ee96ee24dc40`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-smoke-20261004` @ `e382b6ece46966039aba3697efec3a676d503b1def6837481f4a946e7b7d312a`
- inputs digest: `a9a71ecf4450643b80ea8e5e5878c8068efde8d4ef5d18a0a56f08eda5985b19`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `endpoint_mode`: round-robin across 3 endpoints
- setting `endpoints`: yb-n1,yb-n2,yb-n3
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: cluster3
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 123 | 1 |
| yb-n2 | 122 | 1 |
| yb-n3 | 122 | 1 |

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 2000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 2000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 180000 | 180000 |
| 20 | 20 |
| 2000 | 2000 |

| statement | support | rank index observed | forced policy evidence |
|---|---|---|---|
| j01_join_where | valid | false |  |
| j02_join_on | valid | false |  |
| j03_filtered_derived | valid | false |  |
| j04_cte_filtered | valid | false | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | false | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | false |  |
| j11_optional_on | valid | false |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 76.1711 | missing | missing | missing | missing | missing | 70.9197 | 188.9457 | 75.2028 | missing |
| 2000 (1%) | 74.0807 | missing | missing | missing | missing | missing | 87.6699 | 111.0882 | 115.0912 | missing |
| 180000 (90%) | 194.9578 | missing | missing | missing | missing | missing | 220.6920 | 337.7387 | 184.0202 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 76.1711 | 70.9197 | 0.9311 | false |
| 20 | cte_joined_outer_mat vs join_where | 76.1711 | 188.9457 | 2.4805 | false |
| 20 | cte_joined_outer_notmat vs join_where | 76.1711 | 75.2028 | 0.9873 | false |
| 2000 | cte_joined_outer vs join_where | 74.0807 | 87.6699 | 1.1834 | false |
| 2000 | cte_joined_outer_mat vs join_where | 74.0807 | 111.0882 | 1.4996 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 74.0807 | 115.0912 | 1.5536 | false |
| 180000 | cte_joined_outer vs join_where | 194.9578 | 220.6920 | 1.1320 | false |
| 180000 | cte_joined_outer_mat vs join_where | 194.9578 | 337.7387 | 1.7324 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 194.9578 | 184.0202 | 0.9439 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 70.9197 | 188.9457 | 2.6642 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 70.9197 | 75.2028 | 1.0604 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 87.6699 | 111.0882 | 1.2671 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 87.6699 | 115.0912 | 1.3128 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 220.6920 | 337.7387 | 1.5304 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 220.6920 | 184.0202 | 0.8338 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 75.2596 | 1 | `c0ec83b0e2d373d9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 79.1764 | 1 | `d3c3eede1fb5b303` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 80.1985 | 4 | `672954c2915b76fd` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 71.5059 | 4 | `65f3dda2e574384b` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 76.7113 | 16 | `bbf16e5177e9466b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 121.4401 | 16 | `00057f2c7c64e64a` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 86.4523 | 1 | `62f5c387950adf82` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 105.8087 | 1 | `b9e0f15da58fdc81` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 85.9589 | 4 | `74ea8343fc9dde30` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 112.9628 | 4 | `90416df2d8539abd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 78.3127 | 16 | `d7771248b81bd97e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 147.5060 | 16 | `95df444b0faa460a` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 71.8230 | 1 | `3fb9d20563d5cec3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 80.4940 | 1 | `17175e4615225cf9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 250.8866 | 4 | `bc396007394b1c5c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 237.1761 | 4 | `c64e560cf66f0664` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 850.8240 | 16 | `9fc2b44cb2093d3d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 767.9829 | 16 | `39eaa245fe7de32f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1487 | false |
| notmat | 1 | 1% | 0.9543 | false |
| mat | 1 | 100% | 1.3364 | false |
| notmat | 1 | 100% | 1.0166 | false |
| mat | 4 | 1% | 1.0718 | false |
| notmat | 4 | 1% | 3.1283 | false |
| mat | 4 | 100% | 1.5798 | false |
| notmat | 4 | 100% | 3.3169 | false |
| mat | 16 | 1% | 1.0209 | false |
| notmat | 16 | 1% | 11.0912 | false |
| mat | 16 | 100% | 1.2146 | false |
| notmat | 16 | 100% | 6.3240 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 220.6920 | 220.6920 | 220.6920 | 1 | 4.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 70.9197 | 70.9197 | 70.9197 | 1 | 14.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 87.6699 | 87.6699 | 87.6699 | 1 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 337.7387 | 337.7387 | 337.7387 | 1 | 3.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 188.9457 | 188.9457 | 188.9457 | 1 | 5.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 111.0882 | 111.0882 | 111.0882 | 1 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 184.0202 | 184.0202 | 184.0202 | 1 | 5.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 75.2028 | 75.2028 | 75.2028 | 1 | 13.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 115.0912 | 115.0912 | 115.0912 | 1 | 8.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 194.9578 | 194.9578 | 194.9578 | 1 | 5.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 76.1711 | 76.1711 | 76.1711 | 1 | 13.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 74.0807 | 74.0807 | 74.0807 | 1 | 13.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 76.7113 | 76.7113 | 76.7113 | 1 | 13.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 121.4401 | 121.4401 | 121.4401 | 1 | 8.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 75.2596 | 75.2596 | 75.2596 | 1 | 13.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 79.1764 | 79.1764 | 79.1764 | 1 | 12.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 80.1985 | 80.1985 | 80.1985 | 1 | 12.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 71.5059 | 71.5059 | 71.5059 | 1 | 14.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 78.3127 | 78.3127 | 78.3127 | 1 | 12.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 147.5060 | 147.5060 | 147.5060 | 1 | 6.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 86.4523 | 86.4523 | 86.4523 | 1 | 11.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 105.8087 | 105.8087 | 105.8087 | 1 | 9.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 85.9589 | 85.9589 | 85.9589 | 1 | 11.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 112.9628 | 112.9628 | 112.9628 | 1 | 8.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 850.8240 | 850.8240 | 850.8240 | 1 | 1.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 767.9829 | 767.9829 | 767.9829 | 1 | 1.3 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 71.8230 | 71.8230 | 71.8230 | 1 | 13.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 80.4940 | 80.4940 | 80.4940 | 1 | 12.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 250.8866 | 250.8866 | 250.8866 | 1 | 4.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 237.1761 | 237.1761 | 237.1761 | 1 | 4.2 |

## Execution order

- order seed `42254165`, 1 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=bag;stmt=join_where;threshold=2000
  4. family=reuse;policy=default;refs=1;consumer=100%
  5. family=bag;stmt=cte_joined_outer;threshold=20
  6. family=bag;stmt=join_where;threshold=20
  7. family=reuse;policy=default;refs=4;consumer=100%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  9. family=reuse;policy=mat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=mat;refs=1;consumer=100%
  12. family=reuse;policy=mat;refs=16;consumer=100%
  13. family=reuse;policy=notmat;refs=4;consumer=1%
  14. family=reuse;policy=default;refs=16;consumer=100%
  15. family=reuse;policy=mat;refs=4;consumer=1%
  16. family=bag;stmt=join_where;threshold=180000
  17. family=reuse;policy=default;refs=1;consumer=1%
  18. family=reuse;policy=mat;refs=1;consumer=1%
  19. family=bag;stmt=cte_joined_outer;threshold=2000
  20. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  21. family=reuse;policy=notmat;refs=16;consumer=1%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  23. family=bag;stmt=cte_joined_outer_mat;threshold=20
  24. family=reuse;policy=default;refs=16;consumer=1%
  25. family=reuse;policy=notmat;refs=4;consumer=100%
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=reuse;policy=mat;refs=16;consumer=1%
  28. family=reuse;policy=default;refs=4;consumer=1%
  29. family=reuse;policy=notmat;refs=16;consumer=100%
  30. family=bag;stmt=cte_joined_outer;threshold=180000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `cc46b61998bf9b79` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `958e3c073924a40a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `ca5cb605db0d49b2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `42df360f6d7c1a56` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `8a10641dc310386f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `0511fd67a568fd24` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `31f59bf092d56a49` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `14502e98cfcc1253` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `44b8faf826012124` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `04c0de41fa9b98cf` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `6a8e1f566481483e` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `43ac8de45dfc3241` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `bbf16e5177e9466b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `00057f2c7c64e64a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c0ec83b0e2d373d9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `d3c3eede1fb5b303` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `672954c2915b76fd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `65f3dda2e574384b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d7771248b81bd97e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `95df444b0faa460a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `62f5c387950adf82` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `b9e0f15da58fdc81` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `74ea8343fc9dde30` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `90416df2d8539abd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `9fc2b44cb2093d3d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `39eaa245fe7de32f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `3fb9d20563d5cec3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `17175e4615225cf9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `bc396007394b1c5c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c64e560cf66f0664` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
