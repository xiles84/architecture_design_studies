# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 18 whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.

## Compatibility classes

Runs are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.

| class | runs | inputs |
|---|---|---|
| `eb75d021938d` | select-join-topo-yb-cluster-bo1-20261004, select-join-topo-yb-cluster-bo2-20261004, select-join-topo-yb-cluster-bo3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-topo-yb-cluster-bo1-20261004`@`79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| `3033324e094c` | select-join-topo-yb-cluster-i1-20261004, select-join-topo-yb-cluster-i2-20261004, select-join-topo-yb-cluster-i3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-topo-yb-cluster-i1-20261004`@`d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| `76aa820fc736` | select-join-topo-yb-equal-bo1-20261004, select-join-topo-yb-equal-bo2-20261004, select-join-topo-yb-equal-bo3-20261004, select-join-topo-yb-single-bo1-20261004, select-join-topo-yb-single-bo2-20261004, select-join-topo-yb-single-bo3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-topo-yb-equal-bo1-20261004`@`79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| `b5f93048b28f` | select-join-topo-yb-equal-i1-20261004, select-join-topo-yb-equal-i2-20261004, select-join-topo-yb-equal-i3-20261004, select-join-topo-yb-single-i1-20261004, select-join-topo-yb-single-i2-20261004, select-join-topo-yb-single-i3-20261004 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-topo-yb-equal-i1-20261004`@`d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |

## Runs

| run | class | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|---|
| select-join-topo-yb-cluster-bo1-20261004 | `eb75d021938d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-cluster-bo2-20261004 | `eb75d021938d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-cluster-bo3-20261004 | `eb75d021938d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-cluster-i1-20261004 | `3033324e094c` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-cluster-i2-20261004 | `3033324e094c` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-cluster-i3-20261004 | `3033324e094c` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-bo1-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-bo2-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-bo3-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-i1-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-i2-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-equal-i3-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-bo1-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-bo2-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-bo3-20261004 | `76aa820fc736` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-i1-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-i2-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-topo-yb-single-i3-20261004 | `b5f93048b28f` | yugabyte | single | cte-domain | 200000 | 42 |

## Cells across runs (per compatibility class, engine and topology)

| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|---|
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 153.6320, 157.5145, 121.5300 | 153.6320 | 23.4 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 61.0815, 63.7881, 59.0641 | 61.0815 | 7.7 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 56.9925, 58.2179, 62.8824 | 58.2179 | 10.1 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 170.6905, 116.2442, 132.0531 | 132.0531 | 41.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 136.1561, 92.8607, 137.8927 | 136.1561 | 33.1 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 92.5430, 91.6260, 95.8888 | 92.5430 | 4.6 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 123.1130, 163.0366, 131.8830 | 131.8830 | 30.3 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 58.0963, 58.3723, 62.3660 | 58.3723 | 7.3 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 60.1670, 62.4607, 61.8067 | 61.8067 | 3.7 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 156.3926, 159.0656, 161.7484 | 159.0656 | 3.4 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 58.1714, 61.2836, 59.7423 | 59.7423 | 5.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 61.6067, 58.8678, 71.1768 | 61.6067 | 20.0 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 63.1134, 61.0842, 66.4819 | 63.1134 | 8.6 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 59.3168, 66.6020, 68.7383 | 66.6020 | 14.1 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 57.6754, 64.1190, 57.7788 | 57.7788 | 11.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 60.2633, 60.6322, 63.0079 | 60.6322 | 4.5 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 59.7394, 61.0438, 58.2497 | 59.7394 | 4.7 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 56.4671, 59.9284, 63.2386 | 59.9284 | 11.3 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 58.5864, 60.6726, 64.8867 | 60.6726 | 10.4 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 58.1184, 64.7898, 58.7467 | 58.7467 | 11.4 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 56.7674, 61.6748, 63.6636 | 61.6748 | 11.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 60.1879, 60.1288, 60.0949 | 60.1288 | 0.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 57.4937, 59.8558, 61.2062 | 59.8558 | 6.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 59.9555, 61.4999, 62.3945 | 61.4999 | 4.0 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 533.9507, 607.9929, 598.5338 | 598.5338 | 12.4 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 554.3707, 615.0791, 593.5546 | 593.5546 | 10.2 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 59.1565, 59.4324, 70.3594 | 59.4324 | 18.8 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 59.3696, 62.8535, 64.7306 | 62.8535 | 8.5 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 137.5559, 140.4409, 188.3161 | 140.4409 | 36.1 | 3 |
| `3033324e094c` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 153.7558, 177.7283, 145.7134 | 153.7558 | 20.8 | 3 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 125.5740, 117.4256, 144.3187, 125.0343, 125.5266, 125.4940 | 125.5103 | 21.4 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 60.7990, 53.9288, 61.7175, 57.5677, 57.8894, 59.7128 | 58.8011 | 13.2 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 64.6677, 51.7414, 63.4538, 61.7730, 59.0295, 60.3632 | 61.0681 | 21.2 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 130.1365, 113.7144, 121.1784, 125.7757, 124.6410, 122.8594 | 123.7502 | 13.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 107.1603, 92.4960, 97.6875, 108.0833, 94.9424, 96.5852 | 97.1363 | 16.0 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 96.3998, 77.4101, 95.3780, 104.8964, 92.3406, 100.2182 | 95.8889 | 28.7 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 132.8978, 132.1863, 125.5814, 136.3596, 158.4339, 129.8621 | 132.5420 | 24.8 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 68.1762, 49.1918, 61.2431, 60.9876, 56.4210, 55.3569 | 58.7043 | 32.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 60.9297, 50.2477, 65.8494, 60.0741, 59.0918, 63.4369 | 60.5019 | 25.8 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 130.7126, 129.3549, 130.7395, 119.6240, 140.8515, 128.5681 | 130.0337 | 16.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 55.7718, 50.3190, 60.4591, 62.6396, 57.2546, 59.5696 | 58.4121 | 21.1 | 6 |
| `76aa820fc736` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 64.3848, 50.2280, 63.9741, 59.6312, 64.0663, 59.0996 | 61.8026 | 22.9 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 58.4504, 51.0631, 64.7180, 62.2506, 56.5063, 58.1751 | 58.3127 | 23.4 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 62.3052, 55.3749, 64.6120, 65.7374, 58.9286, 59.5917 | 60.9485 | 17.0 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 64.3880, 60.7967, 65.0221, 54.4408, 58.3898, 64.2076 | 62.5022 | 16.9 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 69.2072, 54.3752, 72.9628, 61.9992, 60.0259, 61.3180 | 61.6586 | 30.1 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 62.7746, 56.0266, 63.2297, 56.8826, 60.5417, 56.1946 | 58.7121 | 12.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 57.7688, 49.6663, 67.6621, 57.4907, 57.5371, 59.1869 | 57.6529 | 31.2 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 60.3551, 50.7252, 61.2280, 60.6356, 60.0821, 63.1612 | 60.4953 | 20.6 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 66.1339, 50.5845, 64.2402, 58.7067, 57.9146, 63.7872 | 61.2470 | 25.4 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 58.1759, 50.8713, 71.6489, 56.9979, 58.5865, 57.2049 | 57.6904 | 36.0 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 61.7149, 50.9506, 61.1923, 54.4799, 55.6929, 54.8806 | 55.2867 | 19.5 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 65.0029, 54.6479, 72.4158, 59.9407, 55.2996, 53.8015 | 57.6202 | 32.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 59.7647, 51.1746, 61.8597, 57.1308, 59.5685, 58.5001 | 59.0343 | 18.1 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 441.8267, 441.0225, 514.7835, 786.3588, 843.0020, 680.3916 | 597.5876 | 67.3 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 444.6848, 382.3463, 462.0861, 672.4271, 653.0132, 712.6552 | 557.5496 | 59.2 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 56.1354, 52.6655, 65.2811, 59.6673, 56.3947, 62.8064 | 58.0310 | 21.7 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 64.1528, 52.8947, 67.5904, 66.9390, 64.4059, 62.6936 | 64.2794 | 22.9 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 145.9511, 105.9753, 144.6326, 329.3815, 189.9891, 200.8992 | 167.9701 | 133.0 | 6 |
| `76aa820fc736` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 142.9396, 121.4790, 139.1589, 175.1936, 169.2588, 187.6542 | 156.0992 | 42.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 115.0769, 124.9935, 131.2919, 115.3082, 114.2169, 148.6556 | 120.1509 | 28.7 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 57.7261, 57.8967, 58.6695, 54.5759, 58.2810, 57.4356 | 57.8114 | 7.1 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 57.7575, 56.1353, 59.6361, 57.1971, 53.2324, 65.6797 | 57.4773 | 21.7 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 120.7959, 119.4418, 117.5039, 115.5422, 111.5070, 126.1569 | 118.4728 | 12.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 85.2691, 91.4008, 93.3009, 84.7686, 102.3405, 107.0941 | 92.3509 | 24.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 88.1047, 92.6476, 101.9812, 85.3482, 88.2535, 110.8418 | 90.4506 | 28.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 116.7473, 119.8374, 128.3253, 111.8797, 112.3960, 137.0404 | 118.2923 | 21.3 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 57.9434, 55.7310, 60.1860, 54.7661, 57.5154, 60.1146 | 57.7294 | 9.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 56.3890, 56.1235, 61.0918, 58.5522, 54.8947, 60.5303 | 57.4706 | 10.8 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 119.2619, 116.6264, 121.7237, 119.9078, 117.9654, 123.4313 | 119.5849 | 5.7 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 55.4153, 54.6696, 57.6055, 56.7329, 58.9060, 61.8880 | 57.1692 | 12.6 | 6 |
| `b5f93048b28f` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 56.4833, 58.2504, 58.8241, 57.8108, 56.8327, 64.8960 | 58.0306 | 14.5 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 57.1678, 53.8209, 62.5562, 57.1951, 56.8409, 64.7629 | 57.1814 | 19.1 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 56.9487, 57.9235, 59.7381, 57.3638, 58.3907, 59.1468 | 58.1571 | 4.8 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 58.9152, 54.5555, 60.5576, 60.3669, 58.4089, 57.1655 | 58.6620 | 10.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 57.0732, 54.9889, 62.9782, 56.9828, 59.2367, 68.7158 | 58.1549 | 23.6 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 56.8895, 55.4522, 62.5908, 54.9054, 56.8430, 65.3764 | 56.8662 | 18.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 53.7868, 54.4170, 57.7512, 55.5775, 55.3589, 63.0049 | 55.4682 | 16.6 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 56.0107, 56.8127, 60.7020, 56.9096, 57.3420, 58.0999 | 57.1258 | 8.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 56.2572, 56.8694, 59.7208, 54.1710, 57.9552, 62.8693 | 57.4123 | 15.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 55.7032, 59.7615, 57.0674, 53.3344, 54.7696, 57.9389 | 56.3853 | 11.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 57.7261, 59.5332, 57.7332, 53.3539, 58.2328, 57.2066 | 57.7296 | 10.7 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 57.3291, 55.0674, 57.6409, 56.3732, 55.9978, 69.4263 | 56.8512 | 25.3 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 53.7189, 57.1266, 57.3764, 53.5907, 57.7456, 65.6005 | 57.2515 | 21.0 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 376.9908, 382.2118, 393.5135, 559.5926, 588.9804, 626.3086 | 476.5530 | 52.3 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 382.4606, 388.9698, 411.2624, 577.3888, 586.5704, 697.1195 | 494.3256 | 63.7 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 58.3360, 57.2648, 61.5403, 55.0814, 56.2393, 57.4687 | 57.3667 | 11.3 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 57.9457, 58.8432, 59.5980, 61.7294, 57.5775, 63.1438 | 59.2206 | 9.4 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 123.4679, 125.9317, 136.3275, 160.9036, 150.9267, 199.8259 | 143.6271 | 53.2 | 6 |
| `b5f93048b28f` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 124.7187, 127.9154, 140.1911, 179.9540, 141.8750, 181.3559 | 141.0331 | 40.2 | 6 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 166.0411, 165.4338, 164.4606 | 165.4338 | 1.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 61.9044, 69.3139, 60.6358 | 61.9044 | 14.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 66.4139, 70.1515, 60.8721 | 66.4139 | 14.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 170.7248, 147.4668, 153.3465 | 153.3465 | 15.2 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 147.6074, 103.1969, 164.2476 | 147.6074 | 41.4 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 147.2876, 140.0756, 133.6519 | 140.0756 | 9.7 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 131.4017, 181.3808, 130.1848 | 131.4017 | 39.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 66.5217, 64.2143, 64.7563 | 64.7563 | 3.6 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 64.8700, 83.9865, 69.1991 | 69.1991 | 27.6 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 171.7928, 146.4789, 181.1888 | 171.7928 | 20.2 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 64.2191, 62.0984, 59.7879 | 62.0984 | 7.1 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 63.2172, 66.9318, 60.6704 | 63.2172 | 9.9 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 66.1929, 67.4369, 57.9646 | 66.1929 | 14.3 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 65.1117, 64.7954, 60.7695 | 64.7954 | 6.7 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 66.1884, 66.4623, 59.2939 | 66.1884 | 10.8 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 65.6087, 68.4372, 60.3320 | 65.6087 | 12.4 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 64.0627, 66.1701, 63.5566 | 64.0627 | 4.1 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 68.8162, 69.8505, 62.7718 | 68.8162 | 10.3 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 62.6415, 66.0063, 61.7703 | 62.6415 | 6.8 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 64.2062, 64.9486, 60.4983 | 64.2062 | 6.9 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 61.4513, 67.8635, 63.2637 | 63.2637 | 10.1 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 62.6346, 68.6306, 59.0598 | 62.6346 | 15.3 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 64.8848, 68.3145, 57.2060 | 64.8848 | 17.1 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 60.8916, 64.8082, 65.1011 | 64.8082 | 6.5 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 892.2641, 756.3907, 916.4125 | 892.2641 | 17.9 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 637.0188, 634.1322, 605.5264 | 634.1322 | 5.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 63.0788, 67.4081, 59.8556 | 63.0788 | 12.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 67.4091, 70.9761, 58.3541 | 67.4091 | 18.7 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 284.0113, 209.0527, 322.6704 | 284.0113 | 40.0 | 3 |
| `eb75d021938d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 174.6363, 184.0309, 148.9859 | 174.6363 | 20.1 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9640 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 2.2985 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0359 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0506 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 2.3299 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0261 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9665 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9938 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7649 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 0.9284 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9530 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9547 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0274 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0128 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 4.4333 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.8848 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.5377 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9463 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 13.4798 | >1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9861 | <1 |
| select-join-topo-yb-cluster-bo1-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.7835 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.1162 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 1.6618 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0341 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0481 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 2.0928 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.2548 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.1294 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.0067 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.2383 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.0211 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0142 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0028 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0371 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0324 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 3.1593 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9278 | <1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.6346 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9788 | <1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 11.2163 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0024 | >1 |
| select-join-topo-yb-cluster-bo2-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.7867 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0142 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 2.7472 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0831 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0033 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 2.2029 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1406 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9077 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.8463 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7185 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.0669 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0095 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9789 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9672 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9001 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 5.0769 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0371 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.3735 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0657 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 15.8099 | >1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9955 | <1 |
| select-join-topo-yb-cluster-bo3-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.9643 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0500 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 2.3406 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9987 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9251 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 1.5022 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9766 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9823 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.0914 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7872 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 0.9843 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0257 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9987 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9852 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9624 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.3026 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0618 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.7229 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9283 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.4602 | >1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9798 | <1 |
| select-join-topo-yb-cluster-i1-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.3459 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0409 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 1.5153 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9525 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9890 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 1.5565 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0610 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9902 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.7308 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0250 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 0.9619 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9269 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9917 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0366 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9805 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.3007 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0262 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.9657 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9933 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 9.9534 | >1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9728 | <1 |
| select-join-topo-yb-cluster-i2-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.2351 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9886 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 2.3081 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0439 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8835 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 1.3472 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8684 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.7514 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.8164 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.8154 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.1019 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.2177 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9538 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0273 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0508 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 3.2329 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9867 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.3042 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9760 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 9.0030 | >1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.8546 | <1 |
| select-join-topo-yb-cluster-i3-20261004 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.6350 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0901 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.9214 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.2224 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0044 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.4972 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9463 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9607 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9956 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0167 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9035 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.8718 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8917 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9270 | <1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0355 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.3250 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0345 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.4743 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0326 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.5590 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0615 | >1 |
| select-join-topo-yb-equal-bo1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 7.1372 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0717 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.8382 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9776 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0301 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5412 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0004 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9078 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.8791 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0219 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.8367 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.8663 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9370 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9728 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9754 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 1.8915 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0304 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.4459 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9934 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.6368 | >1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9135 | <1 |
| select-join-topo-yb-equal-bo2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 6.9047 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0208 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6158 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0130 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9919 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.4909 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0293 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.1039 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9269 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9605 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1019 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0040 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8387 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9264 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1453 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.2874 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9142 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.0567 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9461 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.9543 | >1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9942 | <1 |
| select-join-topo-yb-equal-bo3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 7.1517 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0417 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.5387 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0456 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0226 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5598 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9983 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9649 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0129 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9789 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9455 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9902 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0114 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0153 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0077 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.1703 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9987 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.3188 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9798 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 6.5945 | >1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9879 | <1 |
| select-join-topo-yb-equal-i1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 6.7159 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0590 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6719 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0194 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9637 | <1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5905 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9635 | <1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0717 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0241 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0275 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0954 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0497 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0826 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0701 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9931 | <1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.2710 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0498 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.3507 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0556 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.1015 | >1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9818 | <1 |
| select-join-topo-yb-equal-i2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 6.7152 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0185 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6197 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0448 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0138 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.7337 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0386 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0786 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9653 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0542 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9424 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0162 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9167 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9463 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9209 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.1781 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9935 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.4275 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9704 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 6.2906 | >1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9997 | <1 |
| select-join-topo-yb-equal-i3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 6.8844 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9190 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.7255 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9736 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0359 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.7591 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0074 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0452 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0514 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1399 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0470 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0960 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8787 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0797 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0538 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 5.7906 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9937 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.0473 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9741 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 12.6321 | >1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.8930 | <1 |
| select-join-topo-yb-single-bo1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.2290 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0111 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6582 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9854 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9214 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.4413 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9224 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.8912 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.8849 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1248 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0034 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9658 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9278 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0730 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9134 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.1382 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0353 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.9417 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0633 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 14.9187 | >1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9828 | <1 |
| select-join-topo-yb-single-bo2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.0814 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0024 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6214 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9293 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0214 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.6958 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0734 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9761 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9556 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0101 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.8909 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9782 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8950 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0224 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9574 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.5751 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9884 | <1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.1705 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0857 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 11.6956 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0704 | >1 |
| select-join-topo-yb-single-bo3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.9590 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9620 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.4942 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9653 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9894 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.4763 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0128 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9616 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9636 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9330 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.8835 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9124 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9363 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0833 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0267 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9306 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9643 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2379 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9950 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.7839 | >1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9443 | <1 |
| select-join-topo-yb-single-i1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.0654 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9894 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.7374 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9764 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9366 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5529 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9659 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9682 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9453 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9528 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9377 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9629 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9831 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9720 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9851 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6551 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0431 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.5628 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0088 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 10.3619 | >1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9925 | <1 |
| select-join-topo-yb-single-i2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.0456 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9281 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.7304 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9713 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0121 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.7080 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9327 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.2044 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0221 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1103 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0135 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0053 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8325 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9189 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0619 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.0565 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0412 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.8784 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.8971 | <1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.6708 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0629 | >1 |
| select-join-topo-yb-single-i3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.7863 | >1 |

## Recorded limits

- [select-join-topo-yb-cluster-bo1-20261004] explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- [select-join-topo-yb-cluster-bo1-20261004] one host, single client: no independent hardware, real-network or equal-total-budget comparison
- [select-join-topo-yb-cluster-bo1-20261004] no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- [select-join-topo-yb-cluster-bo1-20261004] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
