# SELECT comparison (generated measurements)

_This file contains measurements only. Interpretation belongs in a signed analysis._

## Provenance

| Field | Value |
|---|---|
| repo commit | `3d1dfb34340b0b1dee18fdfbaf3a609ace89ddba` |
| benchmark image | `postgres` |
| image digest | `unknown` |
| dataset seed | 42 |
| dataset fingerprint | `f05b8420d983e1c8` |
| setting: scale | tiny |
| setting: topology | single |

## Cells

| run | topology | statement | group | status | median ms | spread % | trials | plan structure |
|---|---|---|---|---|---|---|---|---|
| 20261003T181731Z | single | s09_not_exists | sel_absence | valid | 0.000 | 0.0 | 1 | `2326095297e70a77` |
| 20261003T181731Z | single | s10_left_anti | sel_absence | valid | 0.000 | 0.0 | 1 | `7db40c0f03ba467e` |
| 20261003T181731Z | single | s12_aggregate_join | sel_aggregate | valid | 0.000 | 0.0 | 1 | `d5c188dc90e96d4c` |
| 20261003T181731Z | single | s13_aggregate_preagg | sel_aggregate | valid | 0.000 | 0.0 | 1 | `54409710085549d0` |
| 20261003T181731Z | single | s01_join_where | sel_full_rows | valid | 0.000 | 0.0 | 1 | `77adb06e1c9f3cc4` |
| 20261003T181731Z | single | s02_filtered_derived | sel_full_rows | valid | 0.000 | 0.0 | 1 | `032d3aac03c11060` |
| 20261003T181731Z | single | s03_cte_filtered | sel_full_rows | valid | 0.000 | 0.0 | 1 | `032d3aac03c11060` |
| 20261003T181731Z | single | s04_join_on | sel_full_rows | valid | 0.000 | 0.0 | 1 | `77adb06e1c9f3cc4` |
| 20261003T181731Z | single | s05_comma_join | sel_full_rows | valid | 0.000 | 0.0 | 1 | `77adb06e1c9f3cc4` |
| 20261003T181731Z | single | s06_exists_parent | sel_membership | valid | 0.000 | 0.0 | 1 | `443d469bfd0ec22b` |
| 20261003T181731Z | single | s07_in_parent | sel_membership | valid | 0.000 | 0.0 | 1 | `443d469bfd0ec22b` |
| 20261003T181731Z | single | s08_join_distinct_parent | sel_membership | valid | 0.000 | 0.0 | 1 | `02160b70ce0b88eb` |
| 20261003T181731Z | single | s11_optional_on | sel_optional_match | valid | 0.000 | 0.0 | 1 | `7db40c0f03ba467e` |
| 20261003T181731Z | single | s16_page_offset | sel_pagination | valid | 0.000 | 0.0 | 1 | `8277f870b7a26176` |
| 20261003T181731Z | single | s14_topn_rownumber | sel_topn | valid | 0.000 | 0.0 | 1 | `0299eee0672bdcbd` |
| 20261003T181731Z | single | s15_topn_lateral | sel_topn | valid | 0.000 | 0.0 | 1 | `1fefdee7c854cb89` |

## Excluded cells

None: every cell is valid.

## Comparisons

### Same plan (no plan difference to explain)

- `s01_join_where` vs `s04_join_on` on single: same plan structure
- `s01_join_where` vs `s05_comma_join` on single: same plan structure
- `s02_filtered_derived` vs `s03_cte_filtered` on single: same plan structure
- `s04_join_on` vs `s05_comma_join` on single: same plan structure
- `s06_exists_parent` vs `s07_in_parent` on single: same plan structure

### No resolved difference

- `s09_not_exists` vs `s10_left_anti` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s12_aggregate_join` vs `s13_aggregate_preagg` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s01_join_where` vs `s02_filtered_derived` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s01_join_where` vs `s03_cte_filtered` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s02_filtered_derived` vs `s04_join_on` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s02_filtered_derived` vs `s05_comma_join` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s03_cte_filtered` vs `s04_join_on` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s03_cte_filtered` vs `s05_comma_join` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s06_exists_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s07_in_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference
- `s14_topn_rownumber` vs `s15_topn_lateral` on single: medians differ by 0.000 ms, inside the 0.0%/0.0% spread — no resolved difference

### Resolved differences

None.

### Not comparable (never ranked)

- - `s09_not_exists` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions

## Predictions this run informs

- P1-P5 join/CTE spellings
- P6-P7 membership/absence
- P8-P9 aggregates/top-N
- P10-P11 predicates/pagination

Verdicts are assigned in the signed analysis, not here.
