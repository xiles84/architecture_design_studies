# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 12 whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.

## Compatibility classes

Runs are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.

| class | runs | inputs |
|---|---|---|
| `d227e8adbf1f` | select-join-v3b-domain-pg-b1-20261005, select-join-v3b-domain-pg-b2-20261005, select-join-v3b-domain-pg-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-domain-pg-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/library/postgres:17.11`@`7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`engine default`; phase=`blended`; index_arm=`pg-btree-rank`; endpoints=`single endpoint` |
| `41b04c169c1e` | select-join-v3b-domain-pg-i1-20261005, select-join-v3b-domain-pg-i2-20261005, select-join-v3b-domain-pg-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-domain-pg-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/library/postgres:17.11`@`7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`engine default`; phase=`isolated`; index_arm=`pg-btree-rank`; endpoints=`single endpoint` |
| `90c87a23d12c` | select-join-v3b-domain-yb-b1-20261005, select-join-v3b-domain-yb-b2-20261005, select-join-v3b-domain-yb-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-domain-yb-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`engine default`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `30b5211d80e9` | select-join-v3b-domain-yb-i1-20261005, select-join-v3b-domain-yb-i2-20261005, select-join-v3b-domain-yb-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-domain-yb-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`engine default`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |

## Runs

| run | class | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|---|
| select-join-v3b-domain-pg-b1-20261005 | `d227e8adbf1f` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-pg-b2-20261005 | `d227e8adbf1f` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-pg-b3-20261005 | `d227e8adbf1f` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-pg-i1-20261005 | `41b04c169c1e` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-pg-i2-20261005 | `41b04c169c1e` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-pg-i3-20261005 | `41b04c169c1e` | postgres | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-b1-20261005 | `90c87a23d12c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-b2-20261005 | `90c87a23d12c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-b3-20261005 | `90c87a23d12c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-i1-20261005 | `30b5211d80e9` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-i2-20261005 | `30b5211d80e9` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-domain-yb-i3-20261005 | `30b5211d80e9` | yugabyte | single | cte-domain | 200000 | 42 |

## Cells across runs (per compatibility class, engine and topology)

| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|---|
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=100000 | 91.7894, 89.3394, 86.1487 | 89.3394 | 6.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=180000 | 129.6052, 105.3497, 110.8970 | 110.8970 | 21.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=180000;position=last | 108.9677 | 108.9677 | 0.0 | 1 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=20 | 1.9191, 2.1732, 1.8744 | 1.9191 | 15.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=200 | 3.2192, 3.8043, 3.3052 | 3.3052 | 17.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=2000 | 10.4146, 9.4514, 9.6032 | 9.6032 | 10.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=comma_join;threshold=20000 | 67.2274, 70.8262, 62.5773 | 67.2274 | 12.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=100000 | 95.5874, 90.6630, 88.0537 | 90.6630 | 8.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=180000 | 108.9847, 109.4802, 105.9201 | 108.9847 | 3.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20 | 1.7762, 1.9058, 1.7147 | 1.7762 | 10.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=200 | 3.6180, 3.6271, 3.9767 | 3.6271 | 9.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=2000 | 9.0673, 8.9177, 8.4272 | 8.9177 | 7.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20000 | 68.0607, 66.6370, 62.5907 | 66.6370 | 8.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 96.8301, 92.1453, 109.9779 | 96.8301 | 18.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 121.2151, 117.1653, 112.4860 | 117.1653 | 7.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 1.8862, 2.0027, 2.1491 | 2.0027 | 13.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 3.6131, 2.9714, 2.8068 | 2.9714 | 27.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 10.3554, 9.3338, 9.1607 | 9.3338 | 12.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 67.1924, 68.6005, 67.4564 | 67.4564 | 2.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 91.5077, 87.1437, 87.4766 | 87.4766 | 5.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 111.1850, 106.7813, 104.0539 | 106.7813 | 6.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 2.4434, 1.9788, 1.7721 | 1.9788 | 33.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 3.2571, 3.5907, 3.5011 | 3.5011 | 9.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 9.5377, 9.6720, 9.4982 | 9.5377 | 1.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 68.1455, 67.3947, 64.2181 | 67.3947 | 5.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 91.6823, 88.7379, 86.4482 | 88.7379 | 5.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 111.0178, 107.2786, 107.5754 | 107.5754 | 3.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 2.0400, 1.7806, 1.8999 | 1.8999 | 13.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=200 | 4.2330, 3.4536, 2.9245 | 3.4536 | 37.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.6825, 9.8780, 10.0949 | 9.8780 | 4.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 65.8666, 66.8383, 62.7383 | 65.8666 | 6.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 100.7127, 98.7858, 91.1245 | 98.7858 | 9.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 113.7440, 102.8562, 100.4942 | 102.8562 | 12.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 88.4830, 86.6027, 80.1526 | 86.6027 | 9.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 86.6851, 86.1154, 84.0988 | 86.1154 | 3.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 85.8304, 85.4882, 81.4994 | 85.4882 | 5.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 90.7192, 108.7087, 85.5257 | 90.7192 | 25.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | 93.6706 | 93.6706 | 0.0 | 1 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 91.1699, 85.3816, 86.0598 | 86.0598 | 6.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 115.3694, 111.0756, 107.4375 | 111.0756 | 7.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.9852, 1.9109, 1.5736 | 1.9109 | 21.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 3.7862, 3.4811, 3.1082 | 3.4811 | 19.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 10.0079, 9.4694, 9.8936 | 9.8936 | 5.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 64.8759, 67.8046, 62.8089 | 64.8759 | 7.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=100000 | 89.8128, 91.1028, 98.7745 | 91.1028 | 9.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=180000 | 118.3836, 109.8656, 108.1524 | 109.8656 | 9.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20 | 1.9032, 2.1162, 2.2550 | 2.1162 | 16.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=200 | 3.6986, 2.9990, 2.7614 | 2.9990 | 31.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=2000 | 9.3428, 9.4266, 8.8739 | 9.3428 | 5.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20000 | 67.4144, 64.2108, 63.4797 | 64.2108 | 6.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=100000 | 89.8909, 93.8932, 85.3678 | 89.8909 | 9.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=180000 | 107.3997, 119.6913, 107.7808 | 107.7808 | 11.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=20 | 1.9798, 2.0419, 2.2524 | 2.0419 | 13.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=200 | 3.1592, 3.0523, 3.1484 | 3.1484 | 3.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=2000 | 8.4765, 10.4809, 8.7099 | 8.7099 | 23.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_on;threshold=20000 | 62.5451, 66.8211, 64.3026 | 64.3026 | 6.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=100000 | 86.5406, 87.5344, 86.8166 | 86.8166 | 1.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 108.4982, 110.6720, 108.7572 | 108.7572 | 2.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.1515, 1.5563, 1.7500 | 1.7500 | 34.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=200 | 3.1810, 3.8137, 3.4474 | 3.4474 | 18.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 9.6082, 9.6494, 9.4030 | 9.6082 | 2.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=bag;stmt=join_where;threshold=20000 | 66.7985, 65.5114, 64.6913 | 65.5114 | 3.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 9.3403, 9.0276, 9.4162 | 9.3403 | 4.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 10.1317, 10.4865, 9.8073 | 10.1317 | 6.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.2801, 6.9661, 7.3953 | 7.2801 | 5.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.7633, 8.8803, 7.6531 | 8.7633 | 14.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=1% | 8.1555, 9.1840, 9.5131 | 9.1840 | 14.8 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=100% | 9.2403, 9.2138, 8.3448 | 9.2138 | 9.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.0565, 9.2588, 8.0887 | 8.0887 | 14.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 9.3743, 9.2110, 8.6979 | 9.2110 | 7.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1% | 8.7641, 8.8404, 9.1989 | 8.8404 | 4.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1%;position=last | 8.3816 | 8.3816 | 0.0 | 1 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=100% | 9.0791, 9.6147, 8.3036 | 9.0791 | 14.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.4584, 9.6595, 9.3403 | 9.4584 | 3.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.1585, 9.5307, 9.2089 | 9.5307 | 10.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.5311, 8.1452, 8.1884 | 8.1884 | 4.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.8117, 8.4285, 7.9805 | 8.4285 | 9.9 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=1% | 8.7456, 8.7867, 8.7947 | 8.7867 | 0.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=100% | 8.5401, 8.9510, 8.1335 | 8.5401 | 9.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.3248, 9.2661, 8.2328 | 8.3248 | 12.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.4238, 9.1258, 8.0887 | 8.4238 | 12.3 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=1% | 9.2116, 8.7731, 8.5094 | 8.7731 | 8.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=100% | 8.8693, 9.2261, 8.5534 | 8.8693 | 7.6 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 106.8276, 109.3732, 69.5902 | 106.8276 | 37.2 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 91.3246, 112.5226, 111.7905 | 111.7905 | 19.0 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 6.9244, 7.1228, 6.7443 | 6.9244 | 5.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.5240, 8.9821, 7.8115 | 8.5240 | 13.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=1% | 13.1376, 12.8306, 12.5784 | 12.8306 | 4.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=100% | 15.8540, 16.5561, 15.6587 | 15.8540 | 5.7 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 24.3854, 25.7260, 23.0083 | 24.3854 | 11.1 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 31.0500, 32.6066, 28.4025 | 31.0500 | 13.5 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=1% | 56.9874, 56.0457, 53.9595 | 56.0457 | 5.4 | 3 |
| `30b5211d80e9` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=100% | 59.7160, 63.1631, 58.2327 | 59.7160 | 8.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=100000 | 55.6649, 48.4774, 50.1389 | 50.1389 | 14.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=180000 | 53.9519, 72.7510, 75.3407 | 72.7510 | 29.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=20 | 0.4313, 0.5237, 0.7213 | 0.5237 | 55.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=200 | 0.8444, 0.7080, 0.7560 | 0.7560 | 18.0 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=2000 | 2.1599, 2.1973, 1.9456 | 2.1599 | 11.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=comma_join;threshold=20000 | 7.2673, 7.5892, 8.3874 | 7.5892 | 14.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=100000 | 46.3914, 44.6471, 57.8186 | 46.3914 | 28.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=180000 | 80.1814, 78.3874, 50.2553 | 78.3874 | 38.2 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=20 | 0.5432, 0.4773, 0.6643 | 0.5432 | 34.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=200 | 0.8460, 0.7325, 0.7709 | 0.7709 | 14.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=2000 | 2.1408, 1.7271, 2.0126 | 2.0126 | 20.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered;threshold=20000 | 7.4998, 7.3109, 7.6167 | 7.4998 | 4.1 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 60.7805, 60.4186, 63.8025 | 60.7805 | 5.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 88.7252, 80.3723, 80.8848 | 80.8848 | 10.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 0.5306, 0.4179, 0.7395 | 0.5306 | 60.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 0.8011, 0.7822, 0.7268 | 0.7822 | 9.5 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 2.3362, 2.1855, 2.0906 | 2.1855 | 11.2 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 9.5370, 9.6622, 9.5962 | 9.5962 | 1.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 66.2708, 48.1178, 50.1601 | 50.1601 | 36.2 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 77.5869, 75.4187, 49.2250 | 75.4187 | 37.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=180000;position=last | 77.9609 | 77.9609 | 0.0 | 1 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 0.4771, 0.4587, 0.4287 | 0.4587 | 10.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 1.1406, 1.0393, 1.0479 | 1.0479 | 9.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 1.8521, 2.1203, 2.0417 | 2.0417 | 13.1 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 7.4806, 7.5541, 7.0524 | 7.4806 | 6.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20;position=last | 0.6081 | 0.6081 | 0.0 | 1 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 47.8812, 46.0983, 46.5513 | 46.5513 | 3.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 80.1983, 64.4969, 73.4084 | 73.4084 | 21.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20 | 0.4948, 0.4778, 0.4971 | 0.4948 | 3.9 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=200 | 0.7971, 0.8495, 0.9989 | 0.8495 | 23.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 2.0924, 2.0023, 1.7816 | 2.0023 | 15.5 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 7.6227, 7.4931, 7.1662 | 7.4931 | 6.1 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20;position=last | 0.4149 | 0.4149 | 0.0 | 1 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 63.2750, 67.0791, 66.9646 | 66.9646 | 5.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 82.1193, 78.6928, 76.6270 | 78.6928 | 7.0 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 49.0254, 48.1081, 48.4540 | 48.4540 | 1.9 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 50.5228, 48.7891, 51.7686 | 50.5228 | 5.9 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 50.3981, 47.9916, 47.5223 | 47.9916 | 6.0 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 54.0679, 53.8350, 55.8272 | 54.0679 | 3.7 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 67.6797, 46.4157, 46.8125 | 46.8125 | 45.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 77.1722, 89.8642, 49.7390 | 77.1722 | 52.0 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 0.5795, 0.4173, 0.7524 | 0.5795 | 57.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 1.0369, 1.0311, 1.0361 | 1.0361 | 0.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 2.2592, 2.2780, 1.8431 | 2.2592 | 19.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 7.4865, 7.7251, 7.1995 | 7.4865 | 7.0 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=100000 | 47.1748, 45.2736, 54.3888 | 47.1748 | 19.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=180000 | 52.7558, 79.8069, 88.3348 | 79.8069 | 44.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=20 | 0.6767, 0.7151, 0.5387 | 0.6767 | 26.1 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=200 | 1.0697, 0.9991, 0.8377 | 0.9991 | 23.2 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=2000 | 2.0170, 1.7610, 1.8511 | 1.8511 | 13.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=filtered_derived;threshold=20000 | 7.2144, 7.6786, 7.1140 | 7.2144 | 7.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=100000 | 64.6798, 48.9945, 48.0515 | 48.9945 | 33.9 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=180000 | 74.8171, 77.7107, 66.8251 | 74.8171 | 14.5 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=20 | 0.4506, 0.7237, 0.4750 | 0.4750 | 57.5 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=200 | 0.7477, 0.9253, 0.9292 | 0.9253 | 19.6 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=2000 | 2.2845, 1.9979, 1.7797 | 1.9979 | 25.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_on;threshold=20000 | 7.7137, 7.5124, 7.3328 | 7.5124 | 5.1 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=100000 | 48.5971, 49.6616, 52.2632 | 49.6616 | 7.4 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=180000 | 73.1975, 79.8609, 70.0986 | 73.1975 | 13.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=20 | 0.4848, 0.4308, 0.7061 | 0.4848 | 56.8 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=200 | 0.7903, 1.0708, 0.8913 | 0.8913 | 31.5 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=2000 | 1.8787, 2.1664, 1.8410 | 1.8787 | 17.3 | 3 |
| `41b04c169c1e` | postgres | single | family=bag;stmt=join_where;threshold=20000 | 7.6594, 7.6453, 7.6768 | 7.6594 | 0.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=16;consumer=1% | 3.8264, 3.4726, 3.8915 | 3.8264 | 10.9 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=16;consumer=100% | 3.8993, 3.8954, 4.0729 | 3.8993 | 4.6 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=1;consumer=1% | 0.6139, 0.5191, 0.7328 | 0.6139 | 34.8 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=1;consumer=100% | 2.6169, 2.5270, 2.7212 | 2.6169 | 7.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=2;consumer=1% | 2.9484, 2.4462, 2.4471 | 2.4471 | 20.5 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=2;consumer=100% | 2.7170, 2.5693, 3.0303 | 2.7170 | 17.0 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=4;consumer=1% | 3.1623, 3.1445, 2.7572 | 3.1445 | 12.9 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=4;consumer=100% | 3.1052, 3.2507, 3.0613 | 3.1052 | 6.1 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=8;consumer=1% | 3.1344, 2.9094, 2.8666 | 2.9094 | 9.2 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=default;refs=8;consumer=100% | 3.2971, 3.5924, 3.2170 | 3.2971 | 11.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=16;consumer=1% | 3.7634, 3.6563, 3.8925 | 3.7634 | 6.3 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=16;consumer=100% | 4.2304, 4.3386, 4.0885 | 4.2304 | 5.9 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=1;consumer=1% | 2.8372, 2.6265, 2.5651 | 2.6265 | 10.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=1;consumer=100% | 2.6763, 2.6417, 2.7046 | 2.6763 | 2.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=2;consumer=1% | 3.0586, 2.5726, 2.7639 | 2.7639 | 17.6 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=2;consumer=100% | 2.7738, 2.5574, 2.9788 | 2.7738 | 15.2 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=4;consumer=1% | 2.7508, 2.6363, 2.7532 | 2.7508 | 4.3 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=4;consumer=100% | 2.9854, 3.2909, 3.0078 | 3.0078 | 10.2 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=8;consumer=1% | 3.4181, 3.0287, 2.9525 | 3.0287 | 15.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=mat;refs=8;consumer=100% | 3.2834, 3.1866, 3.2280 | 3.2280 | 3.0 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=1% | 2.6808, 2.9461, 2.6584 | 2.6808 | 10.7 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=100% | 16.7457, 16.8836, 16.1016 | 16.7457 | 4.7 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=1% | 0.4855, 0.4685, 0.4815 | 0.4815 | 3.5 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=100% | 2.8192, 2.5348, 2.6732 | 2.6732 | 10.6 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=1% | 0.6380, 0.6589, 0.8109 | 0.6589 | 26.2 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=100% | 5.3424, 5.2653, 5.0722 | 5.2653 | 5.1 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=1% | 0.9916, 1.1180, 0.9861 | 0.9916 | 13.3 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=100% | 8.1854, 8.2070, 7.5602 | 8.1854 | 7.9 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=1% | 1.7151, 1.8230, 1.4900 | 1.7151 | 19.4 | 3 |
| `41b04c169c1e` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=100% | 10.7261, 10.3942, 10.5705 | 10.5705 | 3.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=100000 | 86.7790, 85.8352, 87.6625 | 86.7790 | 2.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=180000 | 106.0419, 106.6831, 109.0250 | 106.6831 | 2.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=20 | 1.8820, 1.9759, 1.8361 | 1.8820 | 7.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=200 | 3.0425, 2.8256, 3.3524 | 3.0425 | 17.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=2000 | 8.5453, 8.9554, 8.7519 | 8.7519 | 4.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=20000 | 62.7880, 61.9678, 63.3152 | 62.7880 | 2.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=comma_join;threshold=20;position=last | 1.7584 | 1.7584 | 0.0 | 1 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=100000 | 89.5374, 84.5817, 87.2704 | 87.2704 | 5.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=180000 | 106.8858, 107.8939, 109.4673 | 107.8939 | 2.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20 | 1.5390, 1.9160, 2.2096 | 1.9160 | 35.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=200 | 3.7898, 3.3555, 3.1969 | 3.3555 | 17.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=2000 | 9.1381, 9.1511, 9.7351 | 9.1511 | 6.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered;threshold=20000 | 63.1623, 64.1243, 64.3350 | 64.1243 | 1.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 91.9753, 94.6522, 90.8677 | 91.9753 | 4.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 116.0640, 114.7388, 114.6585 | 114.7388 | 1.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 1.9284, 2.2389, 2.2421 | 2.2389 | 14.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 3.0267, 2.8621, 3.0646 | 3.0267 | 6.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 10.5616, 8.8490, 9.8200 | 9.8200 | 17.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 64.9904, 66.3241, 65.6807 | 65.6807 | 2.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 86.6649, 88.9944, 86.0191 | 86.6649 | 3.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 111.4675, 107.3462, 109.6103 | 109.6103 | 3.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 2.5929, 2.2016, 1.6458 | 2.2016 | 43.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 3.6592, 3.6881, 3.2804 | 3.6592 | 11.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 9.1046, 9.8422, 9.8942 | 9.8422 | 8.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 60.9284, 62.0557, 64.9550 | 62.0557 | 6.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 88.3008, 85.5083, 87.2393 | 87.2393 | 3.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 108.3577, 108.8847, 112.5127 | 108.8847 | 3.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 2.1226, 1.9081, 2.3552 | 2.1226 | 21.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=200 | 3.5976, 2.8575, 3.1365 | 3.1365 | 23.6 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.2196, 8.6397, 9.5881 | 9.2196 | 10.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 62.7514, 63.6403, 64.1595 | 63.6403 | 2.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 94.9390, 97.9048, 98.4019 | 97.9048 | 3.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 107.2638, 106.4139, 104.3876 | 106.4139 | 2.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 82.0061, 87.1743, 81.2771 | 82.0061 | 7.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 82.9406, 81.9603, 86.7706 | 82.9406 | 5.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 83.4664, 82.1125, 82.8729 | 82.8729 | 1.6 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 84.6818, 85.6688, 88.2650 | 85.6688 | 4.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 88.0015, 86.2338, 84.9575 | 86.2338 | 3.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 107.1141, 108.6003, 109.0088 | 108.6003 | 1.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.7864, 1.9944, 1.9394 | 1.9394 | 10.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 2.9775, 3.6449, 2.8491 | 2.9775 | 26.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 8.9502, 8.2550, 9.3775 | 8.9502 | 12.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 65.8629, 61.1197, 63.4140 | 63.4140 | 7.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=100000 | 94.8574, 85.8603, 87.7463 | 87.7463 | 10.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=180000 | 110.4637, 112.5280, 112.0664 | 112.0664 | 1.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20 | 1.6527, 1.7353, 2.1083 | 1.7353 | 26.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=200 | 2.9905, 3.2006, 3.7681 | 3.2006 | 24.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=2000 | 9.6755, 9.3179, 8.3933 | 9.3179 | 13.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=filtered_derived;threshold=20000 | 63.8833, 63.1170, 62.9096 | 63.1170 | 1.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=100000 | 87.6960, 103.7985, 88.8777 | 88.8777 | 18.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=180000 | 108.8605, 105.7222, 107.8579 | 107.8579 | 2.9 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=20 | 1.7592, 1.9269, 2.1824 | 1.9269 | 22.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=200 | 3.5448, 2.9406, 3.6614 | 3.5448 | 20.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=2000 | 9.1065, 9.8402, 9.6440 | 9.6440 | 7.6 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_on;threshold=20000 | 63.8258, 65.1925, 63.9069 | 63.9069 | 2.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=100000 | 86.3689, 86.6515, 87.0115 | 86.6515 | 0.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 111.4336, 105.1388, 112.8999 | 111.4336 | 7.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.1445, 1.8141, 2.1835 | 2.1445 | 17.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=200 | 2.7167, 2.6196, 2.9529 | 2.7167 | 12.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 9.0819, 8.9702, 9.5336 | 9.0819 | 6.2 | 3 |
| `90c87a23d12c` | yugabyte | single | family=bag;stmt=join_where;threshold=20000 | 63.0249, 63.5257, 63.6895 | 63.5257 | 1.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 9.3432, 9.1396, 8.8372 | 9.1396 | 5.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 9.3055, 10.1888, 10.4310 | 10.1888 | 11.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.0894, 6.3401, 6.8390 | 6.8390 | 11.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.7224, 7.8966, 8.4109 | 8.4109 | 9.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=1% | 8.4098, 8.3935, 7.9880 | 8.3935 | 5.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=2;consumer=100% | 8.4711, 8.5670, 8.8538 | 8.5670 | 4.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.3626, 8.5645, 8.5764 | 8.5645 | 2.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 8.6106, 8.0042, 8.2380 | 8.2380 | 7.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1% | 8.8291, 8.1899, 8.2337 | 8.2337 | 7.8 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=1%;position=last | 8.2537 | 8.2537 | 0.0 | 1 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=default;refs=8;consumer=100% | 8.9518, 9.4045, 9.1456 | 9.1456 | 4.9 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 8.9113, 9.3859, 9.2632 | 9.2632 | 5.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 9.9545, 9.7356, 10.4876 | 9.9545 | 7.6 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.1718, 7.9332, 8.0568 | 8.0568 | 3.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1%;position=last | 7.4565 | 7.4565 | 0.0 | 1 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.8633, 8.8313, 7.9181 | 8.8313 | 10.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=1% | 8.8764, 7.9289, 8.1190 | 8.1190 | 11.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=2;consumer=100% | 8.3026, 8.2272, 8.7531 | 8.3026 | 6.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.9820, 8.5997, 9.0837 | 8.9820 | 5.4 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.9492, 8.9128, 8.5006 | 8.9128 | 5.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=1% | 9.1722, 8.0892, 8.4037 | 8.4037 | 12.9 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=mat;refs=8;consumer=100% | 9.8114, 9.1632, 9.3121 | 9.3121 | 7.0 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 89.5516, 86.1408, 83.9057 | 86.1408 | 6.6 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 109.1164, 108.3816, 99.4849 | 108.3816 | 8.9 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 7.4837, 6.9825, 7.6069 | 7.4837 | 8.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.3393, 7.9434, 8.1437 | 8.1437 | 4.9 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=1% | 12.7368, 12.0652, 11.7632 | 12.0652 | 8.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=100% | 14.9117, 16.2400, 14.7076 | 14.9117 | 10.3 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 21.9993, 23.1679, 23.0552 | 23.0552 | 5.1 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 31.2281, 28.4809, 28.8889 | 28.8889 | 9.5 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=1% | 89.8282, 57.0332, 57.8272 | 57.8272 | 56.7 | 3 |
| `90c87a23d12c` | yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=100% | 56.1452, 58.0818, 56.2829 | 56.2829 | 3.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=100000 | 54.8663, 48.0317, 47.7892 | 48.0317 | 14.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=180000 | 78.4841, 82.6754, 75.9969 | 78.4841 | 8.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=20 | 0.5194, 0.5483, 0.5582 | 0.5483 | 7.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=200 | 0.7180, 0.7244, 0.8816 | 0.7244 | 22.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=2000 | 1.6872, 1.7484, 1.9171 | 1.7484 | 13.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=comma_join;threshold=20000 | 7.5551, 7.2381, 6.9819 | 7.2381 | 7.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=100000 | 50.4568, 46.6222, 47.7381 | 47.7381 | 8.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=180000 | 50.1551, 82.4899, 78.6819 | 78.6819 | 41.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=20 | 0.5001, 0.4949, 0.4595 | 0.4949 | 8.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=200 | 0.9974, 0.7677, 0.9393 | 0.9393 | 24.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=2000 | 2.0456, 1.7193, 1.9683 | 1.9683 | 16.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered;threshold=20000 | 7.2616, 6.9677, 6.8580 | 6.9677 | 5.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 61.8963, 63.4422, 58.4158 | 61.8963 | 8.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | 58.7612 | 58.7612 | 0.0 | 1 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 69.6185, 86.5141, 84.1072 | 84.1072 | 20.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 0.6091, 0.4873, 0.6517 | 0.6091 | 27.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 0.7521, 0.8287, 0.7618 | 0.7618 | 10.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 2.0084, 2.3317, 2.3361 | 2.3317 | 14.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 9.4334, 9.5616, 9.2961 | 9.4334 | 2.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | 9.3948 | 9.3948 | 0.0 | 1 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 49.3899, 50.2679, 44.3648 | 49.3899 | 12.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 82.0364, 47.4730, 77.7130 | 77.7130 | 44.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 0.5153, 0.6742, 0.4534 | 0.5153 | 42.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 1.0120, 0.8808, 0.9332 | 0.9332 | 14.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 1.8422, 1.6827, 1.5831 | 1.6827 | 15.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 7.0651, 7.1754, 6.8165 | 7.0651 | 5.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 45.6567, 50.5956, 49.2412 | 49.2412 | 10.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 64.5399, 81.2971, 79.0973 | 79.0973 | 21.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20 | 0.4750, 0.5445, 0.4514 | 0.4750 | 19.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=200 | 0.9248, 0.7947, 0.7179 | 0.7947 | 26.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 1.8930, 2.1140, 1.7460 | 1.8930 | 19.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 7.2137, 7.0109, 7.4247 | 7.2137 | 5.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer;threshold=200;position=last | 0.6812 | 0.6812 | 0.0 | 1 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 63.2256, 65.3221, 65.4836 | 65.3221 | 3.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 93.0220, 89.8946, 90.7410 | 90.7410 | 3.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 47.0695, 45.5078, 45.6655 | 45.6655 | 3.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 47.4744, 41.8275, 45.6113 | 45.6113 | 12.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 43.7699, 48.2378, 43.6970 | 43.7699 | 10.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 52.9223, 49.9614, 45.3416 | 49.9614 | 15.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 46.6147, 49.9366, 48.8755 | 48.8755 | 6.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 79.4112, 83.2593, 44.9244 | 79.4112 | 48.3 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 0.5622, 0.4843, 0.6821 | 0.5622 | 35.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 0.7677, 0.7546, 0.7267 | 0.7546 | 5.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 1.9513, 1.9725, 2.0480 | 1.9725 | 4.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 7.5818, 9.8826, 7.1753 | 7.5818 | 35.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=100000 | 45.9740, 49.1861, 48.0782 | 48.0782 | 6.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=180000 | 48.5609, 84.1793, 50.7758 | 50.7758 | 70.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=20 | 0.5129, 0.4575, 0.4485 | 0.4575 | 14.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=200 | 0.7617, 0.7752, 0.8574 | 0.7752 | 12.3 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=2000 | 1.7229, 2.0176, 1.7262 | 1.7262 | 17.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=filtered_derived;threshold=20000 | 7.2271, 6.9336, 6.7937 | 6.9336 | 6.3 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=100000 | 46.3995, 49.4591, 51.3712 | 49.4591 | 10.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=180000 | 64.7109, 83.2492, 81.3918 | 81.3918 | 22.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=20 | 0.4160, 0.6495, 0.6739 | 0.6495 | 39.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=200 | 0.8461, 0.7423, 0.7083 | 0.7423 | 18.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=2000 | 1.9423, 1.9598, 1.6713 | 1.9423 | 14.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_on;threshold=20000 | 7.1561, 6.9970, 6.9207 | 6.9970 | 3.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=100000 | 61.1514, 49.5673, 47.5226 | 49.5673 | 27.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=180000 | 85.1346, 47.9526, 62.4210 | 62.4210 | 59.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=20 | 0.4642, 0.4696, 0.5992 | 0.4696 | 28.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=200 | 0.7512, 0.8139, 0.7896 | 0.7896 | 7.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=2000 | 1.8015, 1.9468, 2.0829 | 1.9468 | 14.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=bag;stmt=join_where;threshold=20000 | 7.6169, 7.7079, 6.9500 | 7.6169 | 10.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=16;consumer=1% | 3.0162, 3.4058, 3.1171 | 3.1171 | 12.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=16;consumer=100% | 3.7032, 3.8540, 3.5601 | 3.7032 | 7.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=1;consumer=1% | 0.6132, 0.4208, 0.4388 | 0.4388 | 43.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=1;consumer=100% | 2.3676, 2.4534, 2.5920 | 2.4534 | 9.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=2;consumer=1% | 2.7435, 2.4884, 2.6486 | 2.6486 | 9.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=2;consumer=100% | 2.6144, 2.3819, 2.4816 | 2.4816 | 9.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=4;consumer=1% | 2.7322, 2.5611, 2.5308 | 2.5611 | 7.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=4;consumer=100% | 2.7126, 2.6043, 2.5929 | 2.6043 | 4.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=8;consumer=1% | 2.4437, 2.7971, 2.6456 | 2.6456 | 13.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=default;refs=8;consumer=100% | 3.2566, 2.9761, 3.0538 | 3.0538 | 9.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=16;consumer=1% | 3.3480, 3.3993, 3.1948 | 3.3480 | 6.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=16;consumer=100% | 3.4463, 3.6452, 3.4397 | 3.4463 | 6.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=1;consumer=1% | 2.6243, 2.3541, 2.4904 | 2.4904 | 10.8 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=1;consumer=100% | 2.5801, 2.4626, 2.7651 | 2.5801 | 11.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=2;consumer=1% | 2.6620, 2.7366, 2.5532 | 2.6620 | 6.9 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=2;consumer=100% | 2.7935, 2.3829, 2.4868 | 2.4868 | 16.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=4;consumer=1% | 2.5024, 2.8130, 2.8890 | 2.8130 | 13.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=4;consumer=100% | 2.8295, 2.6825, 3.0520 | 2.8295 | 13.1 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=8;consumer=1% | 3.0563, 3.1038, 2.8292 | 3.0563 | 9.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=mat;refs=8;consumer=100% | 3.1890, 3.0082, 3.1845 | 3.1845 | 5.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=1% | 2.2999, 2.5556, 2.3842 | 2.3842 | 10.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=16;consumer=100% | 15.8456, 14.9760, 15.4788 | 15.4788 | 5.6 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=1% | 0.4883, 0.5110, 0.6732 | 0.5110 | 36.2 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=1;consumer=100% | 2.5503, 2.1385, 2.4696 | 2.4696 | 16.7 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=1% | 0.7900, 0.8286, 0.6784 | 0.7900 | 19.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=2;consumer=100% | 4.8823, 4.5362, 4.6611 | 4.6611 | 7.4 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=1% | 1.0573, 0.8768, 0.9517 | 0.9517 | 19.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=4;consumer=100% | 7.5347, 7.4318, 7.5850 | 7.5347 | 2.0 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=1% | 1.4932, 1.5132, 1.3720 | 1.4932 | 9.5 | 3 |
| `d227e8adbf1f` | postgres | single | family=reuse;policy=notmat;refs=8;consumer=100% | 9.3082, 9.6938, 9.8098 | 9.6938 | 5.2 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | join_on vs join_where | 0.8963 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 1.1049 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 1.0773 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.3121 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.1101 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.0233 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 101.4022 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.2111 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20 | comma_join vs join_where | 1.1189 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | join_on vs join_where | 1.1263 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 1.0139 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 1.3277 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 1.0012 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.3471 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.2310 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 63.1960 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0220 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@200 | comma_join vs join_where | 0.9558 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | join_on vs join_where | 1.0781 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.9563 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.1355 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1148 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0226 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 1.0508 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 24.2961 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0831 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 0.9365 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | join_on vs join_where | 0.9395 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9488 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9534 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.2385 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9276 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9471 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.9480 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9954 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 0.9919 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | join_on vs join_where | 0.7588 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.7518 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.8251 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.0122 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.8077 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.7466 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0339 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.7623 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 0.8972 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | join_on vs join_where | 0.7601 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.5704 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 0.5891 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 0.8177 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9636 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.7581 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0926 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9328 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 0.9219 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 4.2799 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.7964 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0898 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0772 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 0.9703 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2880 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 1.0685 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.8675 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 0.9159 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3870 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 1.0431 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.7777 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.2507 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.6110 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 0.9793 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 2.8583 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 1.1100 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7625 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 0.9306 | <1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.2788 | >1 |
| select-join-v3b-domain-pg-b1-20261005 | postgres | single | position | position control first vs last | 0.9959 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | join_on vs join_where | 1.3829 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 0.9741 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 1.0538 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.0377 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.4357 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.1594 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 96.9024 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0313 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20 | comma_join vs join_where | 1.1675 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | join_on vs join_where | 0.9121 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 0.9525 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 0.9433 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 1.0182 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.0822 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.9764 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 51.3941 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9272 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@200 | comma_join vs join_where | 0.8901 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | join_on vs join_where | 1.0067 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.0364 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 0.8832 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1977 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.8644 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 1.0859 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 24.7783 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0132 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 0.8981 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | join_on vs join_where | 0.9078 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.8996 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9040 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.2405 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9309 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9096 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.4819 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.2821 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 0.9391 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | join_on vs join_where | 0.9978 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9923 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.9406 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.2799 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0141 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 1.0207 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3178 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0075 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 0.9690 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | join_on vs join_where | 1.7361 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 1.7555 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.7202 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.8042 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9900 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.6954 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.8747 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.7363 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 1.7241 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 5.5941 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 1.2142 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0038 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 0.8717 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 1.0997 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.3330 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 1.0004 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.9044 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 1.0983 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3424 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 1.0300 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.8537 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.1097 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5410 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 1.0108 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.2573 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 0.9981 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7504 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 0.9458 | <1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 3.8858 | >1 |
| select-join-v3b-domain-pg-b2-20261005 | postgres | single | position | position control first vs last | 0.8573 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | join_on vs join_where | 1.1248 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 0.7485 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 0.7668 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.0877 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 0.7567 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 0.7535 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 76.2163 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1384 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20 | comma_join vs join_where | 0.9317 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | join_on vs join_where | 0.8970 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 1.0858 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 1.1896 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.9647 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.1817 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.9091 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 57.7620 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9203 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@200 | comma_join vs join_where | 1.1165 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | join_on vs join_where | 0.8024 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.8287 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 0.9450 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1216 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.7600 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.8383 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 20.9793 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9833 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 0.9204 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | join_on vs join_where | 0.9958 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9775 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9868 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.3376 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9808 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 1.0683 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.5240 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0324 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 1.0046 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | join_on vs join_where | 1.0810 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 1.0117 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 1.0045 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.2292 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9335 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 1.0362 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3779 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0285 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 1.0056 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | join_on vs join_where | 1.3039 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.8134 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.2605 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.3474 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 1.2450 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.2672 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.4537 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7197 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 1.2175 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 5.6759 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 1.5344 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0668 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 0.9528 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 0.9640 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2561 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 1.0021 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.8783 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 1.1416 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3761 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 1.1771 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.9253 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.0694 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5186 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 1.0428 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.2123 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0249 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7649 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 0.9662 | <1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.3478 | >1 |
| select-join-v3b-domain-pg-b3-20261005 | postgres | single | position | position control first vs last | 1.0059 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | join_on vs join_where | 0.9295 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 1.3957 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 1.1204 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.0945 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 0.9841 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.0206 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 101.1198 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1953 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20 | comma_join vs join_where | 0.8895 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | join_on vs join_where | 0.9460 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 1.3534 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 1.0705 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 1.0136 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.4431 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.0086 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 63.9266 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.3120 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@200 | comma_join vs join_where | 1.0684 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | join_on vs join_where | 1.2159 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.0736 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.1395 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.2435 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9858 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 1.1137 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 26.8254 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.2025 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 1.1497 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | join_on vs join_where | 1.0071 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9419 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9792 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.2451 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9767 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9952 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 7.0590 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9774 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 0.9488 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | join_on vs join_where | 1.3309 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9707 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.9546 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.2507 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 1.3637 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.9853 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3020 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.3927 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 1.1454 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | join_on vs join_where | 1.0221 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.7207 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.0954 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.2121 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0600 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.0956 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.1219 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0543 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 0.7371 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 4.6216 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.7908 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0227 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0773 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 1.0374 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2164 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 1.0209 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.9663 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 0.8699 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3136 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 0.9614 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.6361 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.0905 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5472 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 0.9959 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.2532 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 0.9835 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7006 | <1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 1.0849 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.2945 | >1 |
| select-join-v3b-domain-pg-i1-20261005 | postgres | single | position | position control first vs last | 1.0048 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | join_on vs join_where | 1.6797 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 1.6597 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 1.1079 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 0.9699 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.0647 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.1089 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 111.6643 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9686 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20 | comma_join vs join_where | 1.2156 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | join_on vs join_where | 0.8642 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 0.9331 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 0.6840 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.7305 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 0.9706 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.7934 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 45.5643 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9630 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@200 | comma_join vs join_where | 0.6612 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | join_on vs join_where | 0.9222 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.8129 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 0.7972 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.0088 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9787 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9243 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 22.1530 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0515 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 1.0143 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | join_on vs join_where | 0.9826 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 1.0044 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9563 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.2638 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9881 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9801 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 7.0415 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0104 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 0.9927 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | join_on vs join_where | 0.9866 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9116 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.8990 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.2166 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9689 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.9282 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3507 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9346 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 0.9762 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | join_on vs join_where | 0.9731 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.9993 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 0.9815 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.0064 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9444 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.8076 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9854 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1253 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 0.9110 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 5.0600 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.9027 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 1.0454 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0031 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 1.0517 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2694 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 0.9954 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 2.0493 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 0.8384 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3556 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 1.0124 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.5246 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.0410 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.6266 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 0.8870 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 2.8934 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0529 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.8484 | <1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 1.1138 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.3342 | >1 |
| select-join-v3b-domain-pg-i2-20261005 | postgres | single | position | position control first vs last | 0.8684 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | join_on vs join_where | 0.6728 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | filtered_derived vs join_where | 0.7630 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_filtered vs join_where | 0.9409 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.0474 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 0.6071 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_joined_outer vs join_where | 0.7040 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 68.6250 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0657 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20 | comma_join vs join_where | 1.0215 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | join_on vs join_where | 1.0425 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | filtered_derived vs join_where | 0.9399 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_filtered vs join_where | 0.8649 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.8154 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.1757 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.1207 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 58.0797 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.1624 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@200 | comma_join vs join_where | 0.8482 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | join_on vs join_where | 0.9667 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.0055 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.0932 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1356 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.1090 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9678 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 25.8137 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0011 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@2000 | comma_join vs join_where | 1.0568 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | join_on vs join_where | 0.9552 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9267 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9922 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.2500 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9187 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9335 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 7.2722 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9378 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@20000 | comma_join vs join_where | 1.0926 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | join_on vs join_where | 0.9194 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | filtered_derived vs join_where | 1.0407 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_filtered vs join_where | 1.1063 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.2208 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9598 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.8907 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.2813 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.8957 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@100000 | comma_join vs join_where | 0.9594 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | join_on vs join_where | 0.9533 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | filtered_derived vs join_where | 1.2602 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_filtered vs join_where | 0.7169 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.1539 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.7022 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 1.0472 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0931 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7096 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | bag@180000 | comma_join vs join_where | 1.0748 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=1 c=1% | mat vs default | 3.5003 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.6570 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=1 c=100% | mat vs default | 0.9939 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=1 c=100% | notmat vs default | 0.9824 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=2 c=1% | mat vs default | 1.1294 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.3313 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=2 c=100% | mat vs default | 0.9830 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.6738 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=4 c=1% | mat vs default | 0.9986 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3577 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=4 c=100% | mat vs default | 0.9825 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.4696 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=8 c=1% | mat vs default | 1.0300 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5198 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=8 c=100% | mat vs default | 1.0034 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.2858 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0003 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.6831 | <1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=16 c=100% | mat vs default | 1.0038 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | reuse n=16 c=100% | notmat vs default | 3.9533 | >1 |
| select-join-v3b-domain-pg-i3-20261005 | postgres | single | position | position control first vs last | 1.4186 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 0.8203 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.7707 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.7177 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.8992 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.2091 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9898 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 38.2394 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8330 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 0.8776 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 1.3048 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.1008 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.3950 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.1141 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.3469 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.3243 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 30.5303 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0960 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 1.1199 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0027 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.0654 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 1.0062 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 1.1629 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0025 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0152 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1904 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9855 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9409 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 1.0127 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 1.0136 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0022 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0312 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9667 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9957 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3436 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0450 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9962 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 1.0154 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0983 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0367 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0649 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0034 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0224 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0992 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0189 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0047 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9769 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9913 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9592 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0416 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0003 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9724 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9626 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9612 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9516 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1527 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0556 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0162 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9561 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.0555 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.5145 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9801 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.7603 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0741 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6307 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0393 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.6267 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0389 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 10.1742 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0960 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.2720 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9538 | <1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.5847 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0697 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.7260 | >1 |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte | single | position | position control first vs last | 0.9343 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 1.0622 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9566 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 1.0562 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.2342 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.2136 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0518 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 48.0538 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0994 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 1.0892 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 1.1225 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.2218 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.2809 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.0926 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.4079 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.0908 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 31.2870 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.3914 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 1.0786 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0970 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.0388 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 1.0202 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9865 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0972 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9631 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1539 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9203 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9983 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 1.0262 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9936 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0094 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0441 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9769 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.0018 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3486 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9621 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9755 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 1.1979 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 0.9909 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 0.9761 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0923 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0270 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 0.9868 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.1299 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9952 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 0.9906 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 1.0055 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 1.0703 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 1.0262 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0913 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0210 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0356 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0121 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0329 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 1.0147 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2513 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1013 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1184 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0059 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.9446 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.4374 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9603 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.8956 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0041 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7051 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.1135 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5583 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 0.9877 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 6.9638 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.9743 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.1760 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0269 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.4250 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9555 | <1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.6373 | >1 |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte | single | position | position control first vs last | 0.9399 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 0.9995 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9656 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 1.0120 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.0268 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 0.7538 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0786 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 37.2235 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8882 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 0.8409 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 1.2399 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.2761 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.0826 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.0378 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.1109 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.0622 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 29.3845 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9648 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 1.1353 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0116 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.8804 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 1.0211 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 1.0300 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0378 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0057 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.6927 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9836 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9180 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 1.0034 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9878 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0101 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0313 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 1.0199 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.0074 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3859 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9957 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9941 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 1.0214 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0084 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0030 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0443 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9886 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0026 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.1309 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9764 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0075 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9553 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9926 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9696 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0156 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9709 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9966 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9246 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9655 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9657 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1781 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1123 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9414 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9682 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.0164 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.4726 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9886 | <1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.6612 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0591 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6882 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0319 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5068 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0207 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 7.0233 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0182 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.1541 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0482 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.4946 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0054 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.5374 | >1 |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte | single | position | position control first vs last | 1.0024 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 0.9202 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.8846 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.8256 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.8767 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.1357 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9482 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 41.1269 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9227 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 0.8920 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 0.9931 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.1627 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.1374 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.1359 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.0239 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.3307 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 27.2512 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.1903 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 1.0120 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 0.8822 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.9724 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.9437 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 1.0778 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9927 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0077 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.9331 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0416 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 1.0839 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9363 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 1.0092 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0189 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0059 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 1.0202 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9860 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3581 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9712 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.0064 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 1.0387 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0378 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.1045 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.1189 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0574 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0594 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.1638 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0535 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0607 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9899 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 1.0911 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 1.0045 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.1172 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 1.0248 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0232 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0483 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0633 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 1.1945 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1718 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9511 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0055 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9727 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.0724 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.6109 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9242 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.7158 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0333 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.0268 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.8986 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3122 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0511 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 6.5024 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.9769 | <1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.5773 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0126 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 11.4373 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0026 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.0138 | >1 |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte | single | position | position control first vs last | 0.8408 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 1.3120 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 1.3597 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 1.2246 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.2868 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.2715 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1441 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 55.6465 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.2278 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 1.3964 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 0.8003 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 0.7864 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 0.9511 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.7791 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 0.9415 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 0.9056 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 22.5803 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9128 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 0.9975 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0862 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.9769 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.9242 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9673 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0023 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0237 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.8595 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9813 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9795 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 1.0200 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9801 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 1.0172 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0472 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 1.0287 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.0203 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.6594 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0350 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.0811 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 1.0726 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0408 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0357 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0527 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9955 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0137 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.1285 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9754 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0206 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 1.0815 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9927 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9892 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0587 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9648 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9693 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9294 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0036 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9519 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1693 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0225 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9491 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0115 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.9567 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.3971 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9715 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.7969 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0008 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7785 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9908 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5400 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 0.9924 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 6.3398 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.9596 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 6.5694 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0700 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 12.1154 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9088 | <1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.7302 | >1 |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte | single | position | position control first vs last | 1.0877 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | join_on vs join_where | 1.2870 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | filtered_derived vs join_where | 1.2886 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.9798 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 1.2281 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.0126 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0856 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 45.8008 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8992 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20 | comma_join vs join_where | 1.0711 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | join_on vs join_where | 0.9133 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | filtered_derived vs join_where | 0.8010 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.1535 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.8142 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.0156 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 0.8483 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 24.3946 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9016 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@200 | comma_join vs join_where | 0.9587 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | join_on vs join_where | 0.9263 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.9437 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.8962 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9742 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0101 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0736 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.6674 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0522 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@2000 | comma_join vs join_where | 1.0213 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9940 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9813 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 0.9675 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0427 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9927 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9698 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3221 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9709 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9673 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | join_on vs join_where | 0.9833 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.1377 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0143 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.2668 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0076 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 0.9958 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0496 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9913 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@100000 | comma_join vs join_where | 0.9923 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9910 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9944 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9739 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0343 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9568 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9891 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9240 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9879 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | bag@180000 | comma_join vs join_where | 1.0197 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1073 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9120 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0428 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0207 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.9245 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.3222 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9747 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.8765 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0178 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.8445 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9300 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2655 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=8 c=1% | mat vs default | 0.9250 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 5.8659 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0301 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 7.0129 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9919 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.3905 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9390 | <1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.3987 | >1 |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte | single | position | position control first vs last | 0.9112 | <1 |

## Per-run endpoint, writer, retry and throttling evidence

| run | topology | phase | index arm | endpoints (ops) | writer offered/completed/failed | writer rate/s | reader retries | overlap | node throttled periods (usec) | client throttled periods (usec) |
|---|---|---|---|---|---|---|---|---|---|---|
| select-join-v3b-domain-pg-b1-20261005 | single | blended | pg-btree-rank | single | 335/335/0 | 10.01 | 0 | 93/637 | pg-single 0 (0) | 287 (16046843) [periods 364] |
| select-join-v3b-domain-pg-b2-20261005 | single | blended | pg-btree-rank | single | 340/340/0 | 10.02 | 0 | 131/637 | pg-single 0 (0) | 290 (16327785) [periods 369] |
| select-join-v3b-domain-pg-b3-20261005 | single | blended | pg-btree-rank | single | 349/349/0 | 10.03 | 0 | 114/637 | pg-single 0 (0) | 307 (16877989) [periods 377] |
| select-join-v3b-domain-pg-i1-20261005 | single | isolated | pg-btree-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | pg-single 0 (0) | 314 (15550957) [periods 399] |
| select-join-v3b-domain-pg-i2-20261005 | single | isolated | pg-btree-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | pg-single 0 (0) | 307 (17621363) [periods 383] |
| select-join-v3b-domain-pg-i3-20261005 | single | isolated | pg-btree-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | pg-single 0 (0) | 325 (15767589) [periods 404] |
| select-join-v3b-domain-yb-b1-20261005 | single | blended | yb-range-rank | single | 526/526/0 | 10.01 | 0 | 261/637 | yb-single 76 (2994430) | 247 (8107224) [periods 604] |
| select-join-v3b-domain-yb-b2-20261005 | single | blended | yb-range-rank | single | 523/523/0 | 10.01 | 1 | 234/638 | yb-single 72 (2544553) | 251 (8872000) [periods 605] |
| select-join-v3b-domain-yb-b3-20261005 | single | blended | yb-range-rank | single | 527/527/0 | 10.01 | 1 | 219/638 | yb-single 68 (2696404) | 242 (10359570) [periods 610] |
| select-join-v3b-domain-yb-i1-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 80 (3524146) | 254 (9634765) [periods 647] |
| select-join-v3b-domain-yb-i2-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 83 (3005333) | 239 (9067213) [periods 626] |
| select-join-v3b-domain-yb-i3-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 72 (2826828) | 242 (9119270) [periods 594] |

## Recorded limits

- [select-join-v3b-domain-pg-b1-20261005] explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- [select-join-v3b-domain-pg-b1-20261005] one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- [select-join-v3b-domain-pg-b1-20261005] one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- [select-join-v3b-domain-pg-b1-20261005] physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- [select-join-v3b-domain-pg-b1-20261005] distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- [select-join-v3b-domain-pg-b1-20261005] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
- [select-join-v3b-domain-pg-i1-20261005] isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- [select-join-v3b-domain-yb-b1-20261005] explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- [select-join-v3b-domain-yb-b1-20261005] physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- [select-join-v3b-domain-yb-b1-20261005] YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
