# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 6 whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.

## Compatibility classes

Runs are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.

| class | runs | inputs |
|---|---|---|
| `7a2d4d325c21` | select-join-domain-pg1-20261004, select-join-domain-pg2-20261004, select-join-domain-pg3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-domain-pg1-20261004`@`fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`; engine=`docker.io/library/postgres:17.11`@`7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1` |
| `2d7a72310ff1` | select-join-domain-yb1-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-domain-yb1-20261004`@`fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| `88af7c567a3a` | select-join-domain-yb2-20261004, select-join-domain-yb3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-domain-yb2-20261004`@`bc7ce0737685b8d0125c7291b019ef64d8dcf34c681bcc75569f23b2d083414b`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |

## Runs

| run | class | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|---|
| select-join-domain-pg1-20261004 | `7a2d4d325c21` | postgres | single | cte-domain | 200000 | 42 |
| select-join-domain-pg2-20261004 | `7a2d4d325c21` | postgres | single | cte-domain | 200000 | 42 |
| select-join-domain-pg3-20261004 | `7a2d4d325c21` | postgres | single | cte-domain | 200000 | 42 |
| select-join-domain-yb1-20261004 | `2d7a72310ff1` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-domain-yb2-20261004 | `88af7c567a3a` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-domain-yb3-20261004 | `88af7c567a3a` | yugabyte | single | cte-domain | 200000 | 42 |

## Cells across runs (per compatibility class, engine and topology)

| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|---|
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=100000 | 123.0144 | 123.0144 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=180000 | 220.4352 | 220.4352 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=20 | 80.1080 | 80.1080 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=200 | 74.1292 | 74.1292 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=2000 | 68.7936 | 68.7936 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=comma_join;threshold=20000 | 126.1999 | 126.1999 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=100000 | 135.5453 | 135.5453 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=180000 | 170.7941 | 170.7941 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20 | 78.4165 | 78.4165 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=200 | 92.6889 | 92.6889 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=2000 | 72.9088 | 72.9088 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20000 | 98.1257 | 98.1257 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 155.3808 | 155.3808 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 170.7817 | 170.7817 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 71.1056 | 71.1056 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 93.0677 | 93.0677 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 84.5401 | 84.5401 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 126.7666 | 126.7666 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 129.8132 | 129.8132 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 166.3519 | 166.3519 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 79.3469 | 79.3469 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 84.4803 | 84.4803 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 84.4173 | 84.4173 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 88.4380 | 88.4380 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 193.0855 | 193.0855 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 191.6006 | 191.6006 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 79.1075 | 79.1075 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=200 | 81.3453 | 81.3453 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 73.2526 | 73.2526 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 97.0167 | 97.0167 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 233.9574 | 233.9574 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 220.1819 | 220.1819 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 188.3907 | 188.3907 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 200.7583 | 200.7583 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 201.0103 | 201.0103 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 148.2526 | 148.2526 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 118.9194 | 118.9194 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 185.1344 | 185.1344 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 70.4831 | 70.4831 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 82.3018 | 82.3018 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 100.7849 | 100.7849 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 96.6108 | 96.6108 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=100000 | 143.6310 | 143.6310 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=180000 | 174.2613 | 174.2613 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20 | 98.3247 | 98.3247 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=200 | 90.3503 | 90.3503 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=2000 | 91.5311 | 91.5311 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20000 | 88.0785 | 88.0785 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=100000 | 118.4268 | 118.4268 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=180000 | 175.2461 | 175.2461 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=180000;position=last | 188.5695 | 188.5695 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=20 | 76.8356 | 76.8356 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=200 | 70.6495 | 70.6495 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=2000 | 77.8327 | 77.8327 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_on;threshold=20000 | 87.3950 | 87.3950 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=100000 | 120.5154 | 120.5154 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 161.5619 | 161.5619 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 105.8785 | 105.8785 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=200 | 73.1880 | 73.1880 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 76.7624 | 76.7624 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=bag;stmt=join_where;threshold=20000 | 93.8910 | 93.8910 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 108.3874 | 108.3874 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 88.5839 | 88.5839 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 74.3600 | 74.3600 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 82.7491 | 82.7491 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=1% | 90.2165 | 90.2165 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=100% | 78.0111 | 78.0111 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 91.3477 | 91.3477 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 80.0318 | 80.0318 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1% | 86.2296 | 86.2296 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=100% | 74.0249 | 74.0249 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 75.6634 | 75.6634 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 78.4442 | 78.4442 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 78.6449 | 78.6449 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 76.3532 | 76.3532 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=1% | 77.1046 | 77.1046 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=100% | 95.3046 | 95.3046 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 84.3075 | 84.3075 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 75.4162 | 75.4162 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=1% | 70.0125 | 70.0125 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=100% | 76.7111 | 76.7111 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 979.1275 | 979.1275 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 1157.2493 | 1157.2493 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 81.0691 | 81.0691 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 79.1151 | 79.1151 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=1% | 243.0256 | 243.0256 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=100% | 201.6763 | 201.6763 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 285.9821 | 285.9821 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 257.7474 | 257.7474 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=1% | 566.6426 | 566.6426 | 0.0 | 1 |
| `2d7a72310ff1` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=100% | 538.3197 | 538.3197 | 0.0 | 1 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=100000 | 127.8640, 85.7783, 86.0306 | 86.0306 | 48.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=180000 | 119.9697, 130.0475, 174.9010 | 130.0475 | 42.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=20 | 1.0066, 0.7429, 1.0890 | 1.0066 | 34.4 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=200 | 1.5537, 1.4749, 1.4856 | 1.4856 | 5.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=2000 | 3.3154, 3.5948, 6.2662 | 3.5948 | 82.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=comma_join;threshold=20000 | 19.7013, 21.3029, 18.0951 | 19.7013 | 16.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=100000 | 111.2374, 81.8470, 112.8479 | 111.2374 | 27.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=180000 | 125.9838, 173.9902, 137.4895 | 137.4895 | 34.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=20 | 0.9902, 0.8937, 1.0551 | 0.9902 | 16.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=200 | 1.5094, 1.1353, 1.6146 | 1.5094 | 31.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=2000 | 4.1071, 3.5722, 3.6275 | 3.6275 | 14.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered;threshold=20000 | 16.9481, 16.9754, 22.0181 | 16.9754 | 29.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 125.4378, 80.1030, 119.3190 | 119.3190 | 38.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 142.6691, 139.4390, 124.3221 | 139.4390 | 13.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 0.9213, 1.0511, 1.1285 | 1.0511 | 19.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 1.3681, 1.4749, 1.4828 | 1.4749 | 7.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 3.6700, 4.6169, 4.9725 | 4.6169 | 28.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 18.9753, 14.4377, 33.6297 | 18.9753 | 101.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 132.2922, 79.9130, 130.0051 | 130.0051 | 40.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 162.8678, 125.3316, 156.8097 | 156.8097 | 23.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 0.9278, 1.1268, 1.0677 | 1.0677 | 18.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 1.7450, 1.2630, 1.4606 | 1.4606 | 33.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 3.4413, 3.0845, 5.0778 | 3.4413 | 57.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 16.6135, 15.7976, 20.0279 | 16.6135 | 25.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 127.4178, 73.7816, 135.8569 | 127.4178 | 48.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 157.1349, 130.8156, 149.6988 | 149.6988 | 17.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20 | 0.9676, 0.6325, 1.1123 | 0.9676 | 49.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=200 | 1.5839, 1.1862, 1.4381 | 1.4381 | 27.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 3.4964, 2.6834, 3.3905 | 3.3905 | 24.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 22.6977, 14.6243, 34.5062 | 22.6977 | 87.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 147.3646, 151.0271, 166.4652 | 151.0271 | 12.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 212.6815, 149.2009, 196.5463 | 196.5463 | 32.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 91.2885, 96.0967, 91.0875 | 91.2885 | 5.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 90.6451, 69.8653, 71.3702 | 71.3702 | 29.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 95.2411, 71.4844, 91.6398 | 91.6398 | 25.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 117.0927, 91.0315, 120.0458 | 117.0927 | 24.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 106.3374, 84.1068, 82.8897 | 84.1068 | 27.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 150.0152, 92.3036, 166.3514 | 150.0152 | 49.4 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 0.8834, 0.9156, 1.0277 | 0.9156 | 15.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 1.5029, 1.3588, 1.7451 | 1.5029 | 25.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 3.4607, 3.4040, 5.4642 | 3.4607 | 59.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 22.3343, 16.7984, 14.2443 | 16.7984 | 48.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=100000 | 120.0549, 82.7441, 112.1246 | 112.1246 | 33.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=180000 | 168.1826, 101.6153, 108.4891 | 108.4891 | 61.4 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=20 | 1.0375, 0.6565, 0.9249 | 0.9249 | 41.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=200 | 1.6059, 1.3723, 1.5025 | 1.5025 | 15.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=2000 | 3.5197, 3.0841, 3.2125 | 3.2125 | 13.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=filtered_derived;threshold=20000 | 15.8680, 15.7120, 10.1589 | 15.7120 | 36.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=100000 | 104.3220, 90.2344, 81.9913 | 90.2344 | 24.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=100000;position=last | 72.5330 | 72.5330 | 0.0 | 1 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=180000 | 175.9199, 98.4095, 139.1107 | 139.1107 | 55.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=20 | 0.8795, 0.9895, 1.0379 | 0.9895 | 16.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=200 | 1.5769, 1.0699, 1.7515 | 1.5769 | 43.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=2000 | 3.4136, 2.2548, 3.3856 | 3.3856 | 34.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_on;threshold=20000 | 15.4189, 15.6357, 13.9540 | 15.4189 | 10.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=100000 | 122.3986, 85.7644, 111.3829 | 111.3829 | 32.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=100000;position=last | 83.7135 | 83.7135 | 0.0 | 1 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=180000 | 158.4710, 98.1927, 103.4024 | 103.4024 | 58.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=20 | 1.0654, 0.8339, 1.1177 | 1.0654 | 26.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=200 | 1.4746, 1.1568, 1.6067 | 1.4746 | 30.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=2000 | 3.2721, 3.2849, 3.5932 | 3.2849 | 9.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=bag;stmt=join_where;threshold=20000 | 19.0600, 17.1928, 24.8175 | 19.0600 | 40.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=16;consumer=1% | 6.5675, 4.4871, 6.0684 | 6.0684 | 34.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=16;consumer=100% | 7.4641, 5.3225, 6.5124 | 6.5124 | 32.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=1;consumer=1% | 1.0877, 0.9778, 1.1366 | 1.0877 | 14.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=1;consumer=100% | 4.7786, 4.6296, 4.4159 | 4.6296 | 7.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=2;consumer=1% | 4.8223, 3.4488, 4.5282 | 4.5282 | 30.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=2;consumer=100% | 6.1419, 4.2866, 4.6779 | 4.6779 | 39.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=4;consumer=1% | 4.9065, 3.8070, 5.6812 | 4.9065 | 38.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=4;consumer=100% | 5.2060, 3.9164, 4.2442 | 4.2442 | 30.4 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=8;consumer=1% | 5.9973, 5.8029, 6.1257 | 5.9973 | 5.4 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=default;refs=8;consumer=100% | 6.1096, 5.4107, 6.3234 | 6.1096 | 14.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=16;consumer=1% | 6.4718, 4.6714, 6.8189 | 6.4718 | 33.2 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=16;consumer=100% | 7.9919, 7.0619, 7.4622 | 7.4622 | 12.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=1;consumer=1% | 4.8967, 4.6471, 5.0182 | 4.8967 | 7.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=1;consumer=100% | 5.1171, 3.5360, 5.0806 | 5.0806 | 31.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=2;consumer=1% | 4.6856, 3.7076, 5.5009 | 4.6856 | 38.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=2;consumer=100% | 4.8363, 3.7604, 5.1618 | 4.8363 | 29.0 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=2;consumer=100%;position=last | 4.9819 | 4.9819 | 0.0 | 1 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=4;consumer=1% | 5.0619, 5.3971, 5.6445 | 5.3971 | 10.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=4;consumer=100% | 5.5407, 4.4396, 3.6059 | 4.4396 | 43.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=8;consumer=1% | 5.3902, 3.9185, 5.7651 | 5.3902 | 34.3 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=mat;refs=8;consumer=100% | 5.6833, 4.2207, 4.3625 | 4.3625 | 33.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=1% | 5.6004, 3.7988, 5.1315 | 5.1315 | 35.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=100% | 33.2215, 27.5992, 53.1498 | 33.2215 | 76.9 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=1% | 1.1817, 0.6936, 1.0737 | 1.0737 | 45.5 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=100% | 4.4720, 3.0746, 4.5309 | 4.4720 | 32.6 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=1% | 1.6276, 1.0894, 1.4273 | 1.4273 | 37.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=100% | 11.5251, 11.5578, 15.5487 | 11.5578 | 34.8 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=1% | 1.7496, 1.3025, 1.9974 | 1.7496 | 39.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=100% | 23.1480, 13.4820, 19.4619 | 19.4619 | 49.7 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=1% | 2.9625, 3.0246, 2.9976 | 2.9976 | 2.1 | 3 |
| `7a2d4d325c21` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=100% | 20.4004, 16.9787, 23.5775 | 20.4004 | 32.3 | 3 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=100000 | 178.5353, 188.7832 | 183.6593 | 5.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=180000 | 256.3938, 174.5075 | 215.4507 | 38.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=20 | 71.3233, 75.0968 | 73.2100 | 5.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=200 | 95.7748, 72.1468 | 83.9608 | 28.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=2000 | 74.1813, 99.5767 | 86.8790 | 29.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=comma_join;threshold=20000 | 125.7364, 115.4632 | 120.5998 | 8.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=100000 | 152.3758, 139.6369 | 146.0064 | 8.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=180000 | 166.7558, 232.7039 | 199.7299 | 33.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20 | 91.8535, 93.0575 | 92.4555 | 1.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=200 | 83.3770, 110.6672 | 97.0221 | 28.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=2000 | 89.6260, 74.6743 | 82.1501 | 18.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20000 | 127.3603, 124.0584 | 125.7094 | 2.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 136.0317, 150.6452 | 143.3385 | 10.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 215.6168, 205.2236 | 210.4202 | 4.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000;position=last | 198.2346 | 198.2346 | 0.0 | 1 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 83.4065, 103.9383 | 93.6724 | 21.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 74.1115, 78.2634 | 76.1875 | 5.4 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 96.6450, 82.7469 | 89.6959 | 15.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 95.2844, 89.2912 | 92.2878 | 6.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 125.1805, 156.0797 | 140.6301 | 22.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 235.0769, 209.7476 | 222.4123 | 11.4 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 100.7976, 73.5282 | 87.1629 | 31.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 75.9594, 109.2058 | 92.5826 | 35.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 80.1999, 74.9934 | 77.5967 | 6.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 120.9388, 101.3577 | 111.1482 | 17.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 167.3406, 171.9248 | 169.6327 | 2.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 243.0776, 168.7611 | 205.9194 | 36.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 90.9569, 85.2483 | 88.1026 | 6.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=200 | 94.2939, 95.4511 | 94.8725 | 1.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 85.2398, 84.1131 | 84.6764 | 1.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 116.0596, 94.0174 | 105.0385 | 21.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 179.2701, 172.9103 | 176.0902 | 3.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 167.8199, 229.4704 | 198.6451 | 31.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 135.0339, 214.1780 | 174.6059 | 45.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 194.2958, 136.4335 | 165.3646 | 35.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 144.3756, 144.2880 | 144.3318 | 0.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 217.0542, 226.3498 | 221.7020 | 4.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 166.4903, 128.2294 | 147.3598 | 26.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 200.0032, 184.6368 | 192.3200 | 8.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000;position=last | 220.6716 | 220.6716 | 0.0 | 1 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 81.2951, 73.8782 | 77.5867 | 9.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 71.0185, 79.0972 | 75.0579 | 10.8 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 99.7000, 74.7956 | 87.2478 | 28.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 105.9995, 120.7276 | 113.3636 | 13.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=100000 | 192.4074, 177.5349 | 184.9712 | 8.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=180000 | 231.3123, 243.6141 | 237.4632 | 5.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20 | 78.9719, 93.5598 | 86.2658 | 16.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=200 | 97.2641, 113.0063 | 105.1352 | 15.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=2000 | 82.5772, 121.5485 | 102.0629 | 38.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20000 | 116.9710, 105.5479 | 111.2594 | 10.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=100000 | 138.1147, 145.1905 | 141.6526 | 5.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=180000 | 247.3008, 183.7657 | 215.5333 | 29.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=20 | 109.0463, 110.1639 | 109.6051 | 1.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=200 | 93.6990, 83.4269 | 88.5630 | 11.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=2000 | 97.2891, 80.4195 | 88.8543 | 19.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_on;threshold=20000 | 97.7491, 98.6984 | 98.2237 | 1.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=100000 | 175.1617, 174.8875 | 175.0246 | 0.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 267.9862, 225.4778 | 246.7320 | 17.2 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 82.6690, 87.1451 | 84.9070 | 5.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=200 | 104.4598, 103.9895 | 104.2247 | 0.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 102.6413, 106.0996 | 104.3705 | 3.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=bag;stmt=join_where;threshold=20000 | 90.0696, 101.1035 | 95.5865 | 11.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 95.6351, 73.8573 | 84.7462 | 25.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 89.4784, 93.9867 | 91.7325 | 4.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 97.8882, 92.4769 | 95.1825 | 5.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 75.9922, 85.9608 | 80.9765 | 12.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=1% | 80.2171, 91.0910 | 85.6540 | 12.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=100% | 104.6142, 76.7118 | 90.6630 | 30.8 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 95.9475, 88.9721 | 92.4598 | 7.5 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 84.1687, 81.1366 | 82.6526 | 3.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1% | 82.8818, 72.9649 | 77.9233 | 12.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=100% | 105.0926, 70.5471 | 87.8199 | 39.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 73.6282, 82.0594 | 77.8438 | 10.8 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 83.1863, 78.3027 | 80.7445 | 6.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 87.8541, 87.8327 | 87.8434 | 0.0 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 87.4619, 102.5196 | 94.9908 | 15.9 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=1% | 98.5504, 94.1085 | 96.3294 | 4.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=100% | 84.7355, 100.8830 | 92.8093 | 17.4 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 86.6428, 77.3296 | 81.9862 | 11.4 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 94.4433, 89.5087 | 91.9760 | 5.4 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=1% | 99.9630, 90.8158 | 95.3894 | 9.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=100% | 79.2663, 76.8639 | 78.0651 | 3.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 1033.6538, 1052.2198 | 1042.9368 | 1.8 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 1035.6494, 1160.0005 | 1097.8249 | 11.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 98.4189, 77.4532 | 87.9360 | 23.8 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 110.1305, 97.9613 | 104.0459 | 11.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=1% | 216.0452, 204.0947 | 210.0700 | 5.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=100% | 250.5153, 219.6598 | 235.0875 | 13.1 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 344.3300, 359.5680 | 351.9490 | 4.3 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 337.3900, 285.6017 | 311.4959 | 16.6 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=1% | 592.9320, 621.2279 | 607.0800 | 4.7 | 2 |
| `88af7c567a3a` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=100% | 656.6843, 635.6290 | 646.1567 | 3.3 | 2 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | join_on vs join_where | 0.8255 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 0.9738 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 0.9294 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 0.8647 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 0.8709 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 0.9083 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 85.6869 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8292 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20 | comma_join vs join_where | 0.9448 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | join_on vs join_where | 1.0694 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 1.0891 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 1.0236 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.9278 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.1834 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.0741 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 61.4701 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0192 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@200 | comma_join vs join_where | 1.0536 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | join_on vs join_where | 1.0432 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.0757 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.2552 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1216 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0517 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 1.0686 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 29.1067 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0576 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 1.0132 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | join_on vs join_where | 0.8090 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.8325 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.8892 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 0.9956 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.8716 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 1.1909 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.1434 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.1718 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 1.0336 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | join_on vs join_where | 0.8523 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9809 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.9088 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.0248 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0808 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 1.0410 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.2040 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.8688 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 1.0447 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | join_on vs join_where | 1.1101 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 1.0613 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 0.7950 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 0.9003 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0277 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.9916 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.3421 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9466 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 0.7570 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 4.5019 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 1.0865 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0708 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 0.9359 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 0.9717 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.3375 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 0.7874 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.8765 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 1.0317 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3566 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 1.0643 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 4.4464 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.8988 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.4940 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 0.9302 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.3391 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 0.9854 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.8527 | <1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 1.0707 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.4508 | >1 |
| select-join-domain-pg1-20261004 | postgres | single | position | position control first vs last | 1.0301 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | join_on vs join_where | 1.1866 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 0.7873 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 1.0717 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.2605 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.3512 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 0.7584 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 115.2330 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0979 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20 | comma_join vs join_where | 0.8908 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | join_on vs join_where | 0.9249 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 1.1863 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 0.9814 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 1.2750 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.0919 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.0254 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 60.3971 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.1747 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@200 | comma_join vs join_where | 1.2750 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | join_on vs join_where | 0.6864 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.9389 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.0875 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.4055 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9390 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.8169 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 21.7614 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0362 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 1.0943 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | join_on vs join_where | 0.9094 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9139 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9874 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 0.8398 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9189 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.8506 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 5.2947 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9771 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 1.2391 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | join_on vs join_where | 1.0521 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9648 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.9543 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 0.9340 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9318 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.8603 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.7610 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9807 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 1.0002 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | join_on vs join_where | 1.0022 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 1.0349 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.7719 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.4201 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 1.2764 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.3322 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.5195 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9400 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 1.3244 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 4.7527 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.7094 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 0.7638 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 0.6641 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 1.0750 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.3159 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 0.8772 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 2.6963 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 1.4177 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3421 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 1.1336 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 3.4425 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.6753 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5212 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 0.7801 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.1380 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0411 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.8466 | <1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 1.3268 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 5.1853 | >1 |
| select-join-domain-pg2-20261004 | postgres | single | position | position control first vs last | 0.9761 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | join_on vs join_where | 0.9286 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 0.8275 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 0.9440 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.0097 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 0.9553 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 0.9952 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 81.4974 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9195 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20 | comma_join vs join_where | 0.9743 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | join_on vs join_where | 1.0901 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 0.9351 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 1.0049 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.9229 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 0.9091 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.8951 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 44.4190 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0861 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@200 | comma_join vs join_where | 0.9246 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | join_on vs join_where | 0.9422 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.8940 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.0095 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.3839 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.4132 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9436 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 25.5034 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.5207 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 1.7439 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | join_on vs join_where | 0.5623 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.4093 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.8872 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.3551 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.8070 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 1.3904 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 4.8371 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.5740 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 0.7291 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | join_on vs join_where | 0.7361 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 1.0067 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 1.0132 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.0713 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 1.1672 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 1.2197 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.4945 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.7442 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 0.7724 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | join_on vs join_where | 1.3453 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 1.0492 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.3297 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.2023 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 1.5165 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.4477 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.9008 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.6088 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 1.6915 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 4.4152 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.9447 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 1.1505 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0260 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 1.2148 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.3152 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 1.1035 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 3.3239 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 0.9935 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3516 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 0.8496 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 4.5856 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.9411 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.4893 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 0.6899 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.7286 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 1.1237 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.8456 | <1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 1.1458 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 8.1613 | >1 |
| select-join-domain-pg3-20261004 | postgres | single | position | position control first vs last | 0.8846 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 0.7257 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9287 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.7406 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.6716 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 0.7494 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.7472 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.7793 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.6657 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 0.7566 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 0.9653 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.2345 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.2665 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.2716 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.1543 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.1115 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 2.7431 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.1245 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 1.0129 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0139 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.1924 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.9498 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 1.1013 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0997 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9543 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 2.6186 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.3129 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.8962 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9308 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9381 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0451 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.3501 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9419 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.0333 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.5790 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0290 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.3441 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 0.9827 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.1918 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.1247 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.2893 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0772 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.6022 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.9413 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9868 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0207 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 1.0847 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 1.0786 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 1.0571 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0571 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0296 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.1859 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.3628 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1459 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 1.3644 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0576 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0902 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9227 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9561 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.8547 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 2.6938 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 1.2217 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 2.5852 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9229 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.1307 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9423 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2206 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 0.8119 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 6.5713 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0363 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 7.2721 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.6981 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.0336 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.8855 | <1 |
| select-join-domain-yb1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 13.0639 | >1 |
| select-join-domain-yb1-20261004 | yugabyte | single | position | position control first vs last | 1.0760 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 1.3191 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9553 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 1.1111 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.0089 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.2193 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1003 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6334 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9834 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 0.8628 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 0.8970 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 0.9311 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 0.7982 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.7095 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 0.7272 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 0.9027 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 1.8600 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.6799 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 0.9169 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 0.9479 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.8045 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.8732 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9416 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 0.7814 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.8305 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.4066 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9713 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.7227 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 1.0853 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 1.2987 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.4140 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0579 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 1.3427 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.2886 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 2.4098 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.1769 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.3960 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 0.7885 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0985 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 0.8699 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 0.7766 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.7147 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 0.9553 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0235 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9505 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0193 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9228 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.8631 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.6223 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 0.8046 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.8772 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9071 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.6262 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7463 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9567 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.8975 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0054 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1509 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.4492 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.2285 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 2.6933 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.8100 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 2.3947 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9030 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.5887 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.1221 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 4.0085 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.2061 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 7.1539 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.7543 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.2486 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.7699 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 10.8083 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9297 | <1 |
| select-join-domain-yb2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.5743 | >1 |
| select-join-domain-yb2-20261004 | yugabyte | single | position | position control first vs last | 0.9194 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 1.2641 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 1.0736 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 1.0678 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.1927 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 0.8437 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9782 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 2.4577 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8478 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 0.8617 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 0.8023 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.0867 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.0642 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.7526 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.0502 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 0.9179 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 1.3120 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.7606 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 0.6938 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 0.7580 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.1456 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.7038 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.7799 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 0.7068 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.7928 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.3599 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.7050 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9385 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9762 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 1.0440 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.2270 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 0.8832 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 1.0025 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9299 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 2.2388 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.1941 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.1420 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 0.8302 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0151 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 0.7984 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 0.8614 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.8925 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 0.9831 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 0.9887 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.7332 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0795 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 0.8150 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 1.0804 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 1.0320 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 0.9102 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9302 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.7485 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0177 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.8189 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.7739 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9498 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.8375 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1926 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.1396 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.0331 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 2.2406 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 1.3151 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 2.8634 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.8691 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 4.0414 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.1032 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5200 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.2447 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 8.5141 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0895 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 9.0100 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.1111 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 14.2467 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.8331 | <1 |
| select-join-domain-yb3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 12.3422 | >1 |
| select-join-domain-yb3-20261004 | yugabyte | single | position | position control first vs last | 1.1952 | >1 |

## Recorded limits

- [select-join-domain-pg1-20261004] explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- [select-join-domain-pg1-20261004] one host, single client: no independent hardware, real-network or equal-total-budget comparison
- [select-join-domain-pg1-20261004] no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- [select-join-domain-pg1-20261004] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
- [select-join-domain-yb1-20261004] explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
