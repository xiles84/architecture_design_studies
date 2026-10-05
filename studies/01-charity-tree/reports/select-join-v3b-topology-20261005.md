# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 24 whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.

## Compatibility classes

Runs are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.

| class | runs | inputs |
|---|---|---|
| `37981fca29b2` | select-join-v3b-coloc-b1-20261005, select-join-v3b-coloc-b2-20261005, select-join-v3b-coloc-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-coloc-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`colocated`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `76a52f217951` | select-join-v3b-coloc-i1-20261005, select-join-v3b-coloc-i2-20261005, select-join-v3b-coloc-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-coloc-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`colocated`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `cda7c8de7aac` | select-join-v3b-topo-cluster-b1-20261005, select-join-v3b-topo-cluster-b2-20261005, select-join-v3b-topo-cluster-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-cluster-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `02ca49dbb93f` | select-join-v3b-topo-cluster-i1-20261005, select-join-v3b-topo-cluster-i2-20261005, select-join-v3b-topo-cluster-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-cluster-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `f18b27a60278` | select-join-v3b-topo-equal-b1-20261005, select-join-v3b-topo-equal-b2-20261005, select-join-v3b-topo-equal-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-equal-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=6 CPU 9g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `a69c8639a6a3` | select-join-v3b-topo-equal-i1-20261005, select-join-v3b-topo-equal-i2-20261005, select-join-v3b-topo-equal-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-equal-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=6 CPU 9g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `835f4b827e5c` | select-join-v3b-topo-single-b1-20261005, select-join-v3b-topo-single-b2-20261005, select-join-v3b-topo-single-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-single-b1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `eb80af426f20` | select-join-v3b-topo-single-i1-20261005, select-join-v3b-topo-single-i2-20261005, select-join-v3b-topo-single-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v3b-topo-single-i1-20261005`@`f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |

## Runs

| run | class | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | `37981fca29b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-coloc-b2-20261005 | `37981fca29b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-coloc-b3-20261005 | `37981fca29b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-coloc-i1-20261005 | `76a52f217951` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-coloc-i2-20261005 | `76a52f217951` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-coloc-i3-20261005 | `76a52f217951` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-b1-20261005 | `cda7c8de7aac` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-b2-20261005 | `cda7c8de7aac` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-b3-20261005 | `cda7c8de7aac` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-i1-20261005 | `02ca49dbb93f` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-i2-20261005 | `02ca49dbb93f` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-cluster-i3-20261005 | `02ca49dbb93f` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-b1-20261005 | `f18b27a60278` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-b2-20261005 | `f18b27a60278` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-b3-20261005 | `f18b27a60278` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-i1-20261005 | `a69c8639a6a3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-i2-20261005 | `a69c8639a6a3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-equal-i3-20261005 | `a69c8639a6a3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-b1-20261005 | `835f4b827e5c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-b2-20261005 | `835f4b827e5c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-b3-20261005 | `835f4b827e5c` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-i1-20261005 | `eb80af426f20` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-i2-20261005 | `eb80af426f20` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v3b-topo-single-i3-20261005 | `eb80af426f20` | yugabyte | single | cte-domain | 200000 | 42 |

## Cells across runs (per compatibility class, engine and topology)

| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|---|
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 106.5822, 168.7502, 155.3366 | 155.3366 | 40.0 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.8159, 2.1342, 2.1286 | 2.1342 | 32.2 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 12.4619, 10.0756, 12.3677 | 12.3677 | 19.3 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 156.2962, 152.8720, 155.5775 | 155.5775 | 2.2 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 119.9445, 146.1342, 133.5754 | 133.5754 | 19.6 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 121.8307, 141.3146, 89.1868 | 121.8307 | 42.8 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 150.1061, 120.7980, 153.0432 | 150.1061 | 21.5 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.7984, 3.0783, 2.3107 | 2.7984 | 27.4 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 12.2331, 12.0708, 11.6497 | 12.0708 | 4.8 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 124.2879, 168.3196, 157.6872 | 157.6872 | 27.9 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.5730, 3.2214, 3.2375 | 3.2214 | 20.6 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 12.2691, 11.2850, 12.5324 | 12.2691 | 10.2 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 11.4626, 11.4274, 10.8586 | 11.4274 | 5.3 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 10.6587, 12.8053, 12.3066 | 12.3066 | 17.4 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 7.2216, 7.6610, 7.2377 | 7.2377 | 6.1 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 8.2658, 10.8564, 9.6954 | 9.6954 | 26.7 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 9.7563, 11.3636, 10.6321 | 10.6321 | 15.1 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 8.9563, 9.7799, 11.0389 | 9.7799 | 21.3 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 9.4616, 12.1679, 9.5029 | 9.5029 | 28.5 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 11.2661, 12.2825, 11.0011 | 11.2661 | 11.4 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 10.6110, 8.7797, 10.5054 | 10.5054 | 17.4 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 9.9394, 10.5959, 8.6484 | 9.9394 | 19.6 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 9.1470, 11.0561, 9.9617 | 9.9617 | 19.2 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 9.8541, 10.8504, 10.4243 | 10.4243 | 9.6 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 80.5107, 92.6179, 83.4273 | 83.4273 | 14.5 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 127.2619, 103.2266, 99.7441 | 103.2266 | 26.7 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 8.6892, 8.5162, 7.9456 | 8.5162 | 8.7 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 11.6640, 10.0453, 9.2810 | 10.0453 | 23.7 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 23.5047, 26.6072, 30.2478 | 26.6072 | 25.3 | 3 |
| `02ca49dbb93f` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 28.4156, 38.5641, 35.9026 | 35.9026 | 28.3 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 140.6864, 144.9383, 140.0713 | 140.6864 | 3.5 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.1549, 2.7007, 2.0925 | 2.1549 | 28.2 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 7.9813, 7.2531, 7.5938 | 7.5938 | 9.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 104.2239, 147.7394, 109.2234 | 109.2234 | 39.8 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 118.7996, 78.8698, 120.2742 | 118.7996 | 34.9 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 121.9556, 126.6278, 80.3757 | 121.9556 | 37.9 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 105.9562, 154.6765, 145.1275 | 145.1275 | 33.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.7018, 2.5863, 2.1770 | 2.5863 | 20.3 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 7.6128, 7.2220, 7.8713 | 7.6128 | 8.5 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 115.9706, 149.1025, 141.1788 | 141.1788 | 23.5 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.4074, 2.4009, 2.4475 | 2.4074 | 1.9 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 7.9846, 8.0471, 7.2744 | 7.9846 | 9.7 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 8.2966, 8.7312, 8.2549 | 8.2966 | 5.7 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 8.4698, 8.2241, 8.5219 | 8.4698 | 3.5 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 1.4230, 2.0189, 1.6301 | 1.6301 | 36.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 7.6359, 6.7427, 7.0792 | 7.0792 | 12.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 7.2355, 7.8086, 7.3308 | 7.3308 | 7.8 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 6.3531, 6.8956, 6.7632 | 6.7632 | 8.0 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 7.8036, 8.1164, 8.2010 | 8.1164 | 4.9 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 9.2489, 9.2145, 8.5150 | 9.2145 | 8.0 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 6.8552, 6.4261, 6.2293 | 6.4261 | 9.7 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 6.9557, 6.4124, 7.5868 | 6.9557 | 16.9 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 6.4670, 8.1476, 7.0739 | 7.0739 | 23.8 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 7.3732, 8.2863, 7.6239 | 7.6239 | 12.0 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 16.9441, 16.9673, 17.3321 | 16.9673 | 2.3 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 86.9239, 96.9797, 89.5045 | 89.5045 | 11.2 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 1.6953, 1.8658, 1.7691 | 1.7691 | 9.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 7.0788, 7.8148, 6.4264 | 7.0788 | 19.6 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 3.8805, 4.9077, 4.8691 | 4.8691 | 21.1 | 3 |
| `37981fca29b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 24.6238, 27.0590, 23.8457 | 24.6238 | 13.0 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 131.2419, 138.5225, 137.1160 | 137.1160 | 5.3 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 1.8661, 2.4730, 2.5959 | 2.4730 | 29.5 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 6.5196, 7.0161, 6.9564 | 6.9564 | 7.1 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 99.3995, 134.5914, 102.3280 | 102.3280 | 34.4 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 111.9835, 109.5407, 114.8409 | 111.9835 | 4.7 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 112.1343, 112.8755, 113.2101 | 112.8755 | 1.0 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 134.7846, 101.1228, 137.9459 | 134.7846 | 27.3 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.4683, 2.3242, 2.5683 | 2.4683 | 9.9 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 7.0956, 6.9640, 6.6207 | 6.9640 | 6.8 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 135.1399, 134.4102, 138.9327 | 135.1399 | 3.3 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.4578, 1.7407, 1.8504 | 1.8504 | 38.8 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 7.0962, 7.1472, 7.2996 | 7.1472 | 2.8 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 7.6562, 7.7729, 8.0656 | 7.7729 | 5.3 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 7.6698, 8.9044, 7.8559 | 7.8559 | 15.7 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 1.9257, 1.6259, 1.6902 | 1.6902 | 17.7 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 6.7059, 6.6151, 6.2446 | 6.6151 | 7.0 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 6.3128, 6.9971, 6.9732 | 6.9732 | 9.8 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 6.2539, 6.2222, 6.5665 | 6.2539 | 5.5 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 7.9898, 7.5222, 7.5354 | 7.5354 | 6.2 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 8.0604, 8.2524, 8.3128 | 8.2524 | 3.1 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 6.3418, 6.2895, 6.3418 | 6.3418 | 0.8 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 6.1495, 5.9706, 6.2955 | 6.1495 | 5.3 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 6.7352, 7.6848, 6.7033 | 6.7352 | 14.6 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 6.4506, 6.9561, 6.3369 | 6.4506 | 9.6 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 12.9642, 17.2355, 13.0658 | 13.0658 | 32.7 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 82.3852, 80.5005, 78.9774 | 80.5005 | 4.2 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 1.6589, 1.3403, 1.6383 | 1.6383 | 19.4 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 6.4259, 6.5568, 6.7292 | 6.5568 | 4.6 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 5.0546, 4.6905, 4.8904 | 4.8904 | 7.4 | 3 |
| `76a52f217951` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 21.2756, 21.7213, 21.4972 | 21.4972 | 2.1 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 118.7013, 119.2145, 126.5837 | 119.2145 | 6.6 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 2.0942, 2.0619, 1.8656 | 2.0619 | 11.1 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 10.1714, 10.3959, 9.2314 | 10.1714 | 11.4 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 128.7761, 110.1336, 104.8720 | 110.1336 | 21.7 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 87.6993, 90.6717, 103.4108 | 90.6717 | 17.3 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 91.2837, 93.2029, 89.0833 | 91.2837 | 4.5 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 118.9832, 117.6864, 116.2756 | 117.6864 | 2.3 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.9719, 2.2333, 2.0867 | 2.0867 | 12.5 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.9870, 9.3842, 10.1596 | 9.9870 | 7.8 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 117.0644, 114.8487, 114.1303 | 114.8487 | 2.6 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 1.9905, 1.9582, 2.0190 | 1.9905 | 3.1 | 3 |
| `835f4b827e5c` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 10.4775, 9.6818, 9.8052 | 9.8052 | 8.1 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 10.4380, 10.3767, 9.5205 | 10.3767 | 8.8 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 10.5890, 10.4461, 10.0204 | 10.4461 | 5.4 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.3647, 7.3297, 7.6182 | 7.3647 | 3.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 9.3756, 10.0718, 8.3830 | 9.3756 | 18.0 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 10.0690, 10.1148, 9.6179 | 10.0690 | 4.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 9.7171, 9.4066, 9.2042 | 9.4066 | 5.5 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.6320, 11.1652, 9.6711 | 9.6711 | 15.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.5088, 10.8862, 10.0578 | 10.5088 | 7.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 9.5621, 9.0427, 9.7282 | 9.5621 | 7.2 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.3882, 10.6477, 9.3078 | 9.3078 | 24.3 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 9.6526, 10.1524, 9.3609 | 9.6526 | 8.2 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.9148, 9.7514, 9.2607 | 9.2607 | 9.0 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 86.5949, 102.1546, 82.4492 | 86.5949 | 22.8 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 107.0287, 129.2604, 122.0000 | 122.0000 | 18.2 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 8.1603, 7.6100, 7.9971 | 7.9971 | 6.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.7926, 8.7710, 8.9194 | 8.7926 | 1.7 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 25.2081, 39.2456, 26.0548 | 26.0548 | 53.9 | 3 |
| `835f4b827e5c` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 31.4259, 32.7979, 30.7302 | 31.4259 | 6.6 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 103.4669, 109.0072, 100.2009 | 103.4669 | 8.5 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 1.7763, 2.1751, 1.8490 | 1.8490 | 21.6 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.6070, 8.9026, 9.5234 | 9.5234 | 7.4 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 103.9693, 102.5118, 104.1211 | 103.9693 | 1.5 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 86.1375, 85.4364, 81.7325 | 85.4364 | 5.2 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 85.8868, 85.3347, 78.1647 | 85.3347 | 9.0 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 110.4309, 106.1922, 104.1177 | 106.1922 | 5.9 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.9820, 1.9179, 1.5638 | 1.9179 | 21.8 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.8174, 9.3560, 9.0075 | 9.3560 | 8.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 108.9713, 107.9493, 101.6297 | 107.9493 | 6.8 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 1.8066, 2.2810, 1.8455 | 1.8455 | 25.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 9.3060, 9.9750, 8.7323 | 9.3060 | 13.4 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 9.0918, 9.5951, 9.9309 | 9.5951 | 8.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 10.0379, 10.1698, 9.5039 | 10.0379 | 6.6 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.2165, 7.0131, 6.8676 | 7.0131 | 5.0 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.2153, 8.3700, 8.9415 | 8.3700 | 8.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.6983, 8.3882, 8.7095 | 8.6983 | 3.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 8.8821, 8.8771, 8.5290 | 8.8771 | 4.0 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.1317, 9.3735, 9.0820 | 9.1317 | 3.2 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.2754, 10.4579, 9.5326 | 10.2754 | 9.0 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.7694, 8.6235, 7.9369 | 8.6235 | 9.7 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.9848, 8.7043, 7.8878 | 8.7043 | 12.6 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.7947, 9.3029, 8.3203 | 8.7947 | 11.2 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.6993, 9.2324, 8.1596 | 8.6993 | 12.3 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 72.2798, 81.2836, 70.9807 | 72.2798 | 14.3 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 83.9565, 83.8727, 78.5057 | 83.8727 | 6.5 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 7.4272, 6.6781, 7.5141 | 7.4272 | 11.3 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 9.4589, 8.0905, 7.8245 | 8.0905 | 20.2 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 24.6731, 24.7592, 22.4908 | 24.6731 | 9.2 | 3 |
| `a69c8639a6a3` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 30.2062, 31.5258, 28.6019 | 30.2062 | 9.7 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 122.2483, 163.2515, 163.1814 | 163.1814 | 25.1 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.8744, 2.0955, 2.3839 | 2.3839 | 32.7 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 10.0018, 10.7292, 11.7377 | 10.7292 | 16.2 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 150.5998, 159.5452, 152.1131 | 152.1131 | 5.9 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 134.0262, 132.4396, 138.5601 | 134.0262 | 4.6 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 136.6209, 139.7980, 99.0413 | 136.6209 | 29.8 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 159.4722, 128.5588, 119.4367 | 128.5588 | 31.1 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.2663, 2.0223, 2.1410 | 2.1410 | 11.4 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 10.6172, 10.5862, 10.0796 | 10.5862 | 5.1 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 162.2235, 126.9917, 160.0020 | 160.0020 | 22.0 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.1715, 3.1531, 2.9576 | 2.9576 | 33.2 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 11.2935, 11.4720, 10.1490 | 11.2935 | 11.7 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 11.3908, 10.5318, 11.7796 | 11.3908 | 11.0 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 11.4519, 11.8938, 12.4767 | 11.8938 | 8.6 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 7.6020, 8.4139, 8.5783 | 8.4139 | 11.6 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 9.8510, 9.9445, 10.6686 | 9.9445 | 8.2 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 10.9476, 9.2930, 10.0524 | 10.0524 | 16.5 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 10.8663, 11.2436, 10.8159 | 10.8663 | 3.9 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 11.2329, 11.2256, 11.2960 | 11.2329 | 0.6 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 11.3404, 11.0195, 10.9831 | 11.0195 | 3.2 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 9.7062, 10.2495, 10.6174 | 10.2495 | 8.9 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 9.1000, 10.2450, 10.2798 | 10.2450 | 11.5 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 10.0172, 9.9866, 10.9263 | 10.0172 | 9.4 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 11.0697, 9.4559, 10.8191 | 10.8191 | 14.9 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 92.7883, 94.9560, 96.1861 | 94.9560 | 3.6 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 99.7729, 105.4313, 116.9069 | 105.4313 | 16.3 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 8.9580, 7.9312, 8.6992 | 8.6992 | 11.8 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 10.1255, 9.7422, 10.0981 | 10.0981 | 3.8 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 32.3579, 31.2613, 27.5075 | 31.2613 | 15.5 | 3 |
| `cda7c8de7aac` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 36.9767, 38.2764, 38.4589 | 38.2764 | 3.9 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 105.3376, 102.6497, 142.8837 | 105.3376 | 38.2 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 1.8114, 2.4953, 2.3719 | 2.3719 | 28.8 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 10.1694, 8.2152, 12.2369 | 10.1694 | 39.5 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 99.5245, 97.6008, 133.1028 | 99.5245 | 35.7 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 77.5173, 77.8521, 117.2604 | 77.8521 | 51.0 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 88.2206, 77.5481, 114.9470 | 88.2206 | 42.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 100.0272, 101.2703, 137.6820 | 101.2703 | 37.2 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.4150, 1.9253, 2.3966 | 2.3966 | 20.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 8.7812, 9.0017, 12.0151 | 9.0017 | 35.9 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 105.4197, 101.0970, 138.4676 | 105.4197 | 35.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.0159, 1.9586, 2.0646 | 2.0159 | 5.3 | 3 |
| `eb80af426f20` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 8.9062, 8.5163, 12.4329 | 8.9062 | 44.0 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 8.8007, 8.5079, 11.6433 | 8.8007 | 35.6 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 9.8736, 9.1553, 12.0289 | 9.8736 | 29.1 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 6.9298, 6.6818, 8.0703 | 6.9298 | 20.0 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.8855, 7.5329, 10.6639 | 8.8855 | 35.2 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 9.0704, 7.9716, 11.1384 | 9.0704 | 34.9 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 9.6797, 8.4075, 10.5764 | 9.6797 | 22.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.6487, 8.7898, 12.0625 | 9.6487 | 33.9 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 9.0270, 9.3698, 12.2842 | 9.3698 | 34.8 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.5555, 8.3475, 10.2226 | 8.5555 | 21.9 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.1036, 8.5228, 10.6396 | 8.5228 | 29.8 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.4936, 8.8278, 10.7332 | 8.8278 | 25.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.1360, 7.9760, 11.5837 | 8.1360 | 44.3 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 76.5929, 99.2313, 115.8743 | 99.2313 | 39.6 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 86.7718, 111.5181, 157.0835 | 111.5181 | 63.0 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 6.4678, 6.4928, 8.6349 | 6.4928 | 33.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.2653, 7.4394, 10.8463 | 8.2653 | 41.2 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 22.2474, 22.1093, 29.9783 | 22.2474 | 35.4 | 3 |
| `eb80af426f20` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 28.0235, 30.1732, 38.7751 | 30.1732 | 35.6 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 118.5620, 117.7177, 120.0761 | 118.5620 | 2.0 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 1.8665, 2.1018, 1.9828 | 1.9828 | 11.9 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 10.0710, 9.5416, 9.7466 | 9.7466 | 5.4 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 111.3526, 110.7488, 108.7952 | 110.7488 | 2.3 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 93.8055, 87.0739, 87.9115 | 87.9115 | 7.7 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 88.3222, 90.7287, 87.5445 | 88.3222 | 3.6 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 115.5723, 114.5476, 115.0683 | 115.0683 | 0.9 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.2151, 1.7972, 2.0601 | 2.0601 | 20.3 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 10.7140, 9.8266, 11.2835 | 10.7140 | 13.6 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 119.1081, 117.6850, 112.1092 | 117.6850 | 5.9 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.0736, 1.8796, 1.8267 | 1.8796 | 13.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 9.6234, 9.6159, 9.6112 | 9.6159 | 0.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 10.2419, 10.6348, 9.8102 | 10.2419 | 8.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 10.6998, 10.1686, 10.2137 | 10.2137 | 5.2 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.8273, 6.9667, 7.7286 | 7.7286 | 11.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.6493, 8.8254, 9.3265 | 8.8254 | 7.7 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.9034, 8.7022, 9.4821 | 8.9034 | 8.8 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 9.3293, 9.7623, 9.2396 | 9.3293 | 5.6 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.6948, 10.7166, 9.6120 | 9.6948 | 11.4 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.3947, 10.9798, 10.4970 | 10.4970 | 5.6 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 9.8492, 8.9680, 9.0351 | 9.0351 | 9.8 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.3108, 8.3939, 9.6020 | 8.3939 | 15.4 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 9.3678, 9.1844, 8.7467 | 9.1844 | 6.8 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 9.1888, 9.4250, 10.0081 | 9.4250 | 8.7 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 81.4039, 75.7946, 78.5458 | 78.5458 | 7.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 91.4780, 95.7041, 87.7698 | 91.4780 | 8.7 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 6.9936, 7.7151, 6.9386 | 6.9936 | 11.1 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.5113, 9.1769, 8.8536 | 8.8536 | 7.5 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 24.8511, 24.0662, 22.8222 | 24.0662 | 8.4 | 3 |
| `f18b27a60278` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 33.2752, 33.3407, 30.9524 | 33.2752 | 7.2 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8951 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 49.3471 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.1223 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9996 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 15.2739 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9534 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.2131 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.8987 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9136 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 4.8176 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.1914 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9109 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9270 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.8938 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.5363 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1606 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.8758 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9406 | <1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.0423 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0920 | >1 |
| select-join-v3b-coloc-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.2628 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.1249 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 32.8499 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0772 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9013 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 15.7359 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8975 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9721 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9909 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0374 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.1830 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9242 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9510 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.1590 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0434 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6285 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.2017 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.9241 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9296 | <1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 1.9433 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.1204 | >1 |
| select-join-v3b-coloc-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 11.7921 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8550 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 49.1414 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.8895 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0439 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 11.0490 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0820 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9922 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.7737 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0280 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.8214 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0852 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0717 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9078 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9650 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6642 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1272 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.5258 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9935 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.0996 | >1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9992 | <1 |
| select-join-v3b-coloc-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.5029 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.7593 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 45.5629 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0043 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9187 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 15.8020 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9999 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9712 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.7355 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9974 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.2932 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.8614 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9170 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9582 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0669 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.8007 | <1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0315 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4020 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0436 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 1.6933 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0509 | >1 |
| select-join-v3b-coloc-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.7415 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.4206 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 62.9277 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.3352 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9817 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 15.7929 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9744 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0306 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.0013 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7523 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.8682 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.8243 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9026 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9912 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0983 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6703 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1179 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4909 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9677 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.2174 | >1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9268 | <1 |
| select-join-v3b-coloc-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.0405 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.4029 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 62.0628 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.3880 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9530 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 15.5091 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9070 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9869 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.7365 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9929 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.7522 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9693 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0082 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0776 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9613 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.7013 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9650 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.2738 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9343 | <1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 1.6199 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0582 | >1 |
| select-join-v3b-coloc-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.0532 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.3237 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 61.7213 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0437 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8856 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 12.0973 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9401 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.7536 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9283 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9830 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2768 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.1784 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9238 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0279 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9150 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.9557 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0187 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4029 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9861 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.1459 | >1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9903 | <1 |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.7123 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.6646 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 42.0027 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.6414 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9353 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 12.1860 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9228 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.2855 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.2563 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0123 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2182 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9426 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0302 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9797 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0746 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 3.3640 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.8410 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4043 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0659 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 9.0161 | >1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9265 | <1 |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.8644 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8060 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 46.8484 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7239 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.1565 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 9.7587 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9932 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0199 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9507 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7465 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2377 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0141 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9636 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9465 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0869 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.7364 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0003 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.5558 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9590 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.1655 | >1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.8803 | <1 |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.3700 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0944 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 46.6166 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0876 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0157 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 9.9299 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9971 | <1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.8575 | <1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.2575 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.2077 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.4693 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.2032 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.2025 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.4111 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9376 | <1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.4092 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1002 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.1727 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.8254 | <1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 7.0238 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0570 | >1 |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 11.9398 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.6625 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 45.3631 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9556 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8928 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 12.5223 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0696 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0026 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9082 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7177 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.1460 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.1116 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9760 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9253 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9729 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.3414 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1095 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.9432 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0648 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.1049 | >1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9592 | <1 |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.0613 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.6575 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 41.2586 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7137 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9869 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 7.1165 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9296 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9851 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9866 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9705 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.4515 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0978 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.8920 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9573 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9370 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.8450 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9443 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.2524 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.8751 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 7.6830 | >1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.8939 | <1 |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.1049 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9001 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 45.2382 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0682 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0465 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1778 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1133 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9954 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9349 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9703 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2583 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.8935 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9609 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9841 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0522 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7912 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9849 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5667 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9466 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.9482 | >1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9715 | <1 |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.5495 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1182 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 46.3262 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9562 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9923 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.4353 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0219 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0003 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9411 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9733 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2873 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1074 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9511 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0398 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0554 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7655 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9654 | <1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4152 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0077 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.1270 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0798 | >1 |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.4117 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0854 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 48.1255 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1278 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0141 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1086 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1740 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0711 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9704 | <1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0264 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1691 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.8978 | <1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0295 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9493 | <1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9224 | <1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.4069 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0832 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3500 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9798 | <1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.0065 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0277 | >1 |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.5933 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9832 | <1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 47.6794 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0971 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0323 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.2291 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0549 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9495 | <1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9541 | <1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0134 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2152 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0292 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0937 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.1514 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0111 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.8365 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9794 | <1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4008 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0044 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.9500 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0237 | >1 |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.3640 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9536 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 37.4552 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8408 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.8925 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.5548 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9379 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0098 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9496 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9837 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2296 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9522 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0399 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9666 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1090 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9517 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0400 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5514 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9769 | <1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.4713 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0283 | >1 |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.2472 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0019 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 44.2871 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8473 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0906 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.9512 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0315 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9859 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0245 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0245 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1557 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0941 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8822 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.8751 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9553 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.5823 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9567 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3535 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9145 | <1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.1474 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0030 | >1 |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.2604 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0521 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 44.0586 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9906 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9708 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.7124 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9532 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0140 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.1000 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0164 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2984 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1080 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.8947 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9378 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9586 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.5035 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9174 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2341 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9228 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.2961 | >1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9924 | <1 |
| select-join-v3b-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.1075 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0530 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 46.3035 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1405 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0737 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.6266 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9693 | <1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0380 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9589 | <1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0247 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2337 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0382 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0572 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.8709 | <1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0037 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.8800 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0367 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4867 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0760 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.8446 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0421 | >1 |
| select-join-v3b-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 12.3740 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9240 | <1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 51.2181 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0335 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9415 | <1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.0853 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0361 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.1091 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9189 | <1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0188 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2770 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0497 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1103 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0640 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9733 | <1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7090 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0061 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3387 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0158 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.6602 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0037 | >1 |
| select-join-v3b-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 12.1752 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.8985 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 38.4533 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1980 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.1418 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.9055 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9860 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9992 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9441 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9488 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2346 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9333 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9120 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9302 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9364 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.4527 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.8405 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.8951 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0964 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.7031 | >1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9143 | <1 |
| select-join-v3b-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.7882 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.2740 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 39.7486 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9830 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9646 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1059 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0570 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0154 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9654 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0017 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2493 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9717 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1314 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9876 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1074 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7735 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9487 | <1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5888 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0331 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 11.6634 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0234 | >1 |
| select-join-v3b-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 12.1807 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1488 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 56.7947 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1608 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9842 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.2454 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9664 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0319 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9613 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9943 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2667 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0700 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9977 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0171 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9636 | <1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6914 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0952 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.6662 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0360 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.9520 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0212 | >1 |
| select-join-v3b-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 13.0588 | >1 |

## Per-run endpoint, writer, retry and throttling evidence

| run | topology | phase | index arm | endpoints (ops) | writer offered/completed/failed | writer rate/s | reader retries | overlap | node throttled periods (usec) | client throttled periods (usec) |
|---|---|---|---|---|---|---|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 128/128/0 | 10.01 | 0 | 61/210 | yb-n1 5 (781983); yb-n2 9 (1719228); yb-n3 0 (0) | 99 (5328518) [periods 283] |
| select-join-v3b-coloc-b2-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 140/140/0 | 10.03 | 0 | 66/210 | yb-n1 5 (572002); yb-n2 12 (1377765); yb-n3 0 (0) | 90 (4295606) [periods 294] |
| select-join-v3b-coloc-b3-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 133/133/0 | 10.05 | 0 | 61/210 | yb-n1 5 (827497); yb-n2 6 (771131); yb-n3 5 (95830) | 94 (4855464) [periods 293] |
| select-join-v3b-coloc-i1-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 9 (827369); yb-n2 5 (684797); yb-n3 9 (388943) | 95 (4547990) [periods 278] |
| select-join-v3b-coloc-i2-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 6 (987249); yb-n2 6 (915820); yb-n3 5 (297311) | 90 (4635477) [periods 280] |
| select-join-v3b-coloc-i3-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 6 (872217); yb-n2 13 (1573888); yb-n3 0 (0) | 89 (4579043) [periods 274] |
| select-join-v3b-topo-cluster-b1-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 176/176/0 | 10.02 | 0 | 68/210 | yb-n1 7 (702031); yb-n2 5 (656221); yb-n3 11 (314027) | 131 (5059249) [periods 348] |
| select-join-v3b-topo-cluster-b2-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 175/175/0 | 10.03 | 0 | 75/210 | yb-n1 4 (689694); yb-n2 6 (551110); yb-n3 6 (65130) | 133 (5278383) [periods 348] |
| select-join-v3b-topo-cluster-b3-20261005 | cluster3 | blended | yb-range-rank | 3 (168/168/168) | 178/178/0 | 10.02 | 0 | 78/210 | yb-n1 5 (553282); yb-n2 6 (969611); yb-n3 13 (354926) | 128 (5343505) [periods 352] |
| select-join-v3b-topo-cluster-i1-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 6 (564235); yb-n2 4 (533493); yb-n3 15 (650060) | 109 (4231433) [periods 315] |
| select-join-v3b-topo-cluster-i2-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 9 (756973); yb-n2 7 (791742); yb-n3 12 (228304) | 134 (6738599) [periods 357] |
| select-join-v3b-topo-cluster-i3-20261005 | cluster3 | isolated | yb-range-rank | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 6 (737504); yb-n2 23 (1814258); yb-n3 0 (0) | 128 (4594089) [periods 354] |
| select-join-v3b-topo-equal-b1-20261005 | single | blended | yb-range-rank | single | 156/156/0 | 10.03 | 1 | 68/211 | yb-single 0 (0) | 136 (5587070) [periods 315] |
| select-join-v3b-topo-equal-b2-20261005 | single | blended | yb-range-rank | single | 155/155/0 | 10.04 | 0 | 72/210 | yb-single 0 (0) | 130 (5468048) [periods 312] |
| select-join-v3b-topo-equal-b3-20261005 | single | blended | yb-range-rank | single | 153/153/0 | 10.01 | 0 | 69/210 | yb-single 0 (0) | 132 (5401798) [periods 313] |
| select-join-v3b-topo-equal-i1-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 115 (3564778) [periods 278] |
| select-join-v3b-topo-equal-i2-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 117 (3892542) [periods 284] |
| select-join-v3b-topo-equal-i3-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 123 (5104909) [periods 278] |
| select-join-v3b-topo-single-b1-20261005 | single | blended | yb-range-rank | single | 156/156/0 | 10.02 | 0 | 73/210 | yb-single 93 (5614109) | 135 (5133819) [periods 313] |
| select-join-v3b-topo-single-b2-20261005 | single | blended | yb-range-rank | single | 160/160/0 | 10.02 | 0 | 74/210 | yb-single 88 (7444711) | 123 (5450145) [periods 313] |
| select-join-v3b-topo-single-b3-20261005 | single | blended | yb-range-rank | single | 160/160/0 | 10.02 | 0 | 65/210 | yb-single 88 (6158308) | 141 (6213880) [periods 335] |
| select-join-v3b-topo-single-i1-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 79 (5163078) | 117 (4723409) [periods 301] |
| select-join-v3b-topo-single-i2-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 80 (7407562) | 107 (4808165) [periods 280] |
| select-join-v3b-topo-single-i3-20261005 | single | isolated | yb-range-rank | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 114 (9189479) | 185 (7627490) [periods 404] |

## Recorded limits

- [select-join-v3b-coloc-b1-20261005] explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- [select-join-v3b-coloc-b1-20261005] one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- [select-join-v3b-coloc-b1-20261005] one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- [select-join-v3b-coloc-b1-20261005] cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- [select-join-v3b-coloc-b1-20261005] readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- [select-join-v3b-coloc-b1-20261005] physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- [select-join-v3b-coloc-b1-20261005] YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- [select-join-v3b-coloc-b1-20261005] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
- [select-join-v3b-coloc-i1-20261005] isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- [select-join-v3b-topo-equal-b1-20261005] one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
