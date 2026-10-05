# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 30 whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.

## Compatibility classes

Runs are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.

| class | runs | inputs |
|---|---|---|
| `09b4d2db03b2` | select-join-v4-topo-colocated-b1-20261005, select-join-v4-topo-colocated-b2-20261005, select-join-v4-topo-colocated-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-colocated-b1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`colocated`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `7c269cf982fd` | select-join-v4-topo-colocated-i1-20261005, select-join-v4-topo-colocated-i2-20261005, select-join-v4-topo-colocated-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-colocated-i1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`colocated`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `6b8c9e0a77ae` | select-join-v4-topo-distributed-b1-20261005, select-join-v4-topo-distributed-b2-20261005, select-join-v4-topo-distributed-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-distributed-b1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`distributed`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `ed469cd9aa6d` | select-join-v4-topo-distributed-i1-20261005, select-join-v4-topo-distributed-i2-20261005, select-join-v4-topo-distributed-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-distributed-i1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`distributed`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `ffdfe4180cce` | select-join-v4-topo-equal-b1-20261005, select-join-v4-topo-equal-b2-20261005, select-join-v4-topo-equal-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-equal-b1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=6 CPU 9g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `e917d1a74006` | select-join-v4-topo-equal-i1-20261005, select-join-v4-topo-equal-i2-20261005, select-join-v4-topo-equal-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-equal-i1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=6 CPU 9g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `fae5cef41884` | select-join-v4-topo-single-b1-20261005, select-join-v4-topo-single-b2-20261005, select-join-v4-topo-single-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-single-b1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `0380110419e3` | select-join-v4-topo-single-i1-20261005, select-join-v4-topo-single-i2-20261005, select-join-v4-topo-single-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-single-i1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`single`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`single endpoint` |
| `64b081a2c954` | select-join-v4-topo-tablet-b1-20261005, select-join-v4-topo-tablet-b2-20261005, select-join-v4-topo-tablet-b3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-tablet-b1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`blended`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |
| `5e9b02ba9a42` | select-join-v4-topo-tablet-i1-20261005, select-join-v4-topo-tablet-i2-20261005, select-join-v4-topo-tablet-i3-20261005 | env=`host-zenbook-ux5406sa`; dataset=`fcf26e08f377ec51`; client=`localhost/charitybench-select:select-join-v4-topo-tablet-i1-20261005`@`b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`; engine=`docker.io/yugabytedb/yugabyte:2025.2.6.0-b111`@`561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`; topology=`cluster3`; node=2 CPU 3g; client_budget=1 CPU 2g; placement=`noncolocated-single-tablet`; phase=`isolated`; index_arm=`yb-range-rank`; endpoints=`round-robin across 3 endpoints` |

## Runs

| run | class | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|---|
| select-join-v4-topo-colocated-b1-20261005 | `09b4d2db03b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-colocated-b2-20261005 | `09b4d2db03b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-colocated-b3-20261005 | `09b4d2db03b2` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-colocated-i1-20261005 | `7c269cf982fd` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-colocated-i2-20261005 | `7c269cf982fd` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-colocated-i3-20261005 | `7c269cf982fd` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-b1-20261005 | `6b8c9e0a77ae` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-b2-20261005 | `6b8c9e0a77ae` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-b3-20261005 | `6b8c9e0a77ae` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-i1-20261005 | `ed469cd9aa6d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-i2-20261005 | `ed469cd9aa6d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-distributed-i3-20261005 | `ed469cd9aa6d` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-b1-20261005 | `ffdfe4180cce` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-b2-20261005 | `ffdfe4180cce` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-b3-20261005 | `ffdfe4180cce` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-i1-20261005 | `e917d1a74006` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-i2-20261005 | `e917d1a74006` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-equal-i3-20261005 | `e917d1a74006` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-b1-20261005 | `fae5cef41884` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-b2-20261005 | `fae5cef41884` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-b3-20261005 | `fae5cef41884` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-i1-20261005 | `0380110419e3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-i2-20261005 | `0380110419e3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-single-i3-20261005 | `0380110419e3` | yugabyte | single | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-b1-20261005 | `64b081a2c954` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-b2-20261005 | `64b081a2c954` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-b3-20261005 | `64b081a2c954` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-i1-20261005 | `5e9b02ba9a42` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-i2-20261005 | `5e9b02ba9a42` | yugabyte | cluster3 | cte-domain | 200000 | 42 |
| select-join-v4-topo-tablet-i3-20261005 | `5e9b02ba9a42` | yugabyte | cluster3 | cte-domain | 200000 | 42 |

## Cells across runs (per compatibility class, engine and topology)

| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|---|
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 118.7366, 131.4230, 102.9684 | 118.7366 | 24.0 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 2.0433, 1.9701, 1.7732 | 1.9701 | 13.7 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.3226, 10.9920, 8.2525 | 9.3226 | 29.4 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 109.2512, 113.0238, 105.1295 | 109.2512 | 7.2 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 88.7571, 86.7975, 77.6468 | 86.7975 | 12.8 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 89.0481, 91.3876, 85.6812 | 89.0481 | 6.4 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 121.2623, 108.0538, 105.8212 | 108.0538 | 14.3 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.0240, 1.8743, 1.8929 | 1.8929 | 7.9 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.7642, 10.2569, 9.5892 | 9.7642 | 6.8 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 108.8445, 126.7814, 104.4094 | 108.8445 | 20.6 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.2348, 1.7054, 1.6959 | 1.7054 | 31.6 | 3 |
| `0380110419e3` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 9.4894, 9.7444, 8.3845 | 9.4894 | 14.3 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 10.0155, 10.5321, 9.8091 | 10.0155 | 7.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 9.5782, 9.9350, 9.8147 | 9.8147 | 3.6 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 7.0662, 7.4711, 6.3985 | 7.0662 | 15.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 8.7268, 9.0043, 7.5035 | 8.7268 | 17.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.6680, 10.2502, 8.7135 | 8.7135 | 18.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 8.8600, 9.3801, 7.8543 | 8.8600 | 17.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 9.1091, 10.0384, 8.4722 | 9.1091 | 17.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 9.7011, 10.2148, 9.3765 | 9.7011 | 8.6 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.4473, 8.7035, 7.8159 | 8.4473 | 10.5 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.8105, 8.6382, 8.3748 | 8.6382 | 5.0 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.6703, 9.9818, 7.6961 | 8.6703 | 26.4 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.9971, 9.3092, 7.5740 | 8.9971 | 19.3 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 98.9009, 119.7908, 70.0928 | 98.9009 | 50.3 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 87.3387, 94.6692, 109.8145 | 94.6692 | 23.7 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 7.0045, 7.2312, 6.0996 | 7.0045 | 16.2 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 9.2390, 9.5693, 7.6852 | 9.2390 | 20.4 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 25.8475, 26.7395, 19.5558 | 25.8475 | 27.8 | 3 |
| `0380110419e3` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 31.5749, 32.6548, 28.4197 | 31.5749 | 13.4 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 131.4460, 93.6174, 102.6479 | 102.6479 | 36.9 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.2153, 2.7091, 2.6428 | 2.6428 | 18.7 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 7.2387, 5.9729, 6.9049 | 6.9049 | 18.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 99.0774, 129.2045, 130.6704 | 129.2045 | 24.5 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 107.8110, 70.6817, 73.9177 | 73.9177 | 50.2 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 66.9229, 70.1970, 111.1469 | 70.1970 | 63.0 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 137.6312, 127.7908, 130.7481 | 130.7481 | 7.5 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.8738, 2.3715, 1.8447 | 1.8738 | 28.1 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 7.0996, 6.9862, 7.4350 | 7.0996 | 6.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 97.3977, 128.0755, 140.2607 | 128.0755 | 33.5 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.4015, 2.2136, 2.3416 | 2.3416 | 8.0 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 7.4860, 6.8613, 6.7446 | 6.8613 | 10.8 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 7.1613, 7.3384, 6.9366 | 7.1613 | 5.6 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 8.0430, 7.8544, 8.3399 | 8.0430 | 6.0 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 1.4670, 1.6814, 1.5879 | 1.5879 | 13.5 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 6.5721, 7.1381, 6.4450 | 6.5721 | 10.5 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 6.6974, 6.5491, 6.5818 | 6.5818 | 2.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 7.3625, 6.2227, 6.2500 | 6.2500 | 18.2 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 7.5863, 7.6111, 7.7631 | 7.6111 | 2.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 7.7655, 8.0260, 8.1085 | 8.0260 | 4.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 6.6272, 6.6599, 6.1899 | 6.6272 | 7.1 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 6.2434, 5.8177, 6.5031 | 6.2434 | 11.0 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 6.8936, 5.9674, 6.4149 | 6.4149 | 14.4 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 6.2739, 6.8122, 6.1680 | 6.2739 | 10.3 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 17.1661, 12.6650, 16.4177 | 16.4177 | 27.4 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 83.8295, 73.9947, 84.0777 | 83.8295 | 12.0 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 1.6655, 1.4840, 1.6901 | 1.6655 | 12.4 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 5.7105, 5.9803, 6.7069 | 5.9803 | 16.7 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 3.6780, 4.2791, 4.5916 | 4.2791 | 21.4 | 3 |
| `09b4d2db03b2` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 23.8469, 22.7672, 23.1546 | 23.1546 | 4.7 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 143.9169, 140.2097, 143.3530 | 143.3530 | 2.6 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.7232, 2.6650, 2.7767 | 2.7232 | 4.1 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 10.6796, 8.9196, 10.3747 | 10.3747 | 17.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 134.0807, 105.2482, 137.7133 | 134.0807 | 24.2 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 120.5125, 84.2837, 116.8543 | 116.8543 | 31.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 125.4706, 122.3975, 124.2378 | 124.2378 | 2.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 141.0170, 107.6607, 141.4773 | 141.0170 | 24.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.7147, 15.2788, 1.9744 | 2.7147 | 490.1 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 10.5949, 11.6093, 11.5250 | 11.5250 | 8.8 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 142.8339, 146.7331, 107.7438 | 142.8339 | 27.3 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 1.8712, 2.7500, 2.6144 | 2.6144 | 33.6 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 9.2925, 10.9394, 11.5129 | 10.9394 | 20.3 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 9.9788, 10.0465, 9.4064 | 9.9788 | 6.4 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 9.9252, 10.8146, 10.4564 | 10.4564 | 8.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 8.4291, 6.9484, 7.7494 | 7.7494 | 19.1 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 8.0878, 8.6243, 8.6454 | 8.6243 | 6.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 9.2660, 10.0216, 9.4074 | 9.4074 | 8.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 9.6243, 9.7735, 9.7193 | 9.7193 | 1.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 9.9818, 10.2975, 9.3783 | 9.9818 | 9.2 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 11.3022, 10.6065, 9.7549 | 10.6065 | 14.6 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 9.1057, 9.5220, 9.0683 | 9.1057 | 5.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 8.3138, 8.4367, 8.1893 | 8.3138 | 3.0 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 9.1150, 9.4803, 8.3419 | 9.1150 | 12.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 9.2572, 8.2991, 9.7356 | 9.2572 | 15.5 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 75.3324, 71.3659, 77.9223 | 75.3324 | 8.7 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 89.6101, 95.7612, 113.6569 | 95.7612 | 25.1 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 7.6431, 8.3977, 7.4002 | 7.6431 | 13.1 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 9.3544, 9.5104, 9.6020 | 9.5104 | 2.6 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 27.2331, 22.9815, 22.0028 | 22.9815 | 22.8 | 3 |
| `5e9b02ba9a42` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 28.1329, 33.3348, 33.0502 | 33.0502 | 15.7 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 124.9643, 154.4661, 113.5581 | 124.9643 | 32.7 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.0321, 2.8920, 2.8027 | 2.8027 | 30.7 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 11.3222, 13.1497, 9.5828 | 11.3222 | 31.5 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 114.5862, 150.2236, 105.9493 | 114.5862 | 38.6 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 131.1921, 133.5939, 128.7220 | 131.1921 | 3.7 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 91.4519, 136.8571, 122.2600 | 122.2600 | 37.1 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 153.9764, 147.0753, 117.0366 | 147.0753 | 25.1 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 3.1898, 3.2901, 2.8511 | 3.1898 | 13.8 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 11.2379, 9.7356, 11.6188 | 11.2379 | 16.8 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 117.4869, 121.9118, 116.4572 | 117.4869 | 4.6 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.0892, 3.2479, 2.9713 | 2.9713 | 39.0 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 12.3004, 11.3076, 11.7788 | 11.7788 | 8.4 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 11.2682, 11.6859, 10.0611 | 11.2682 | 14.4 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 11.3531, 11.7102, 10.1862 | 11.3531 | 13.4 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 8.0753, 7.5779, 6.6946 | 7.5779 | 18.2 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 10.4798, 10.1817, 9.9197 | 10.1817 | 5.5 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 10.1721, 10.5773, 9.8833 | 10.1721 | 6.8 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 9.6506, 9.8072, 10.0376 | 9.8072 | 3.9 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 10.5595, 10.4245, 10.6292 | 10.5595 | 1.9 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 11.7126, 11.0817, 12.0248 | 11.7126 | 8.1 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 10.2271, 8.4837, 9.5204 | 9.5204 | 18.3 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 11.0907, 10.4883, 9.3698 | 10.4883 | 16.4 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 10.3487, 10.0566, 11.3920 | 10.3487 | 12.9 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 10.6480, 9.7623, 10.6557 | 10.6480 | 8.4 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 88.6976, 94.5046, 85.9427 | 88.6976 | 9.7 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 99.4867, 117.8989, 88.2381 | 99.4867 | 29.8 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 7.0769, 7.3729, 8.3972 | 7.3729 | 17.9 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 10.6736, 8.6449, 9.5145 | 9.5145 | 21.3 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 28.8182, 28.4491, 27.5815 | 28.4491 | 4.3 | 3 |
| `64b081a2c954` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 31.3932, 35.7623, 31.0676 | 31.3932 | 15.0 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 87.2105, 86.2752, 86.2732 | 86.2752 | 1.1 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 3.6563, 2.9347, 3.5478 | 3.5478 | 20.3 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 8.5606, 8.3636, 8.9192 | 8.5606 | 6.5 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 117.6337, 111.8808, 114.3705 | 114.3705 | 5.0 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 73.0584, 74.8211, 69.2856 | 73.0584 | 7.6 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 78.3079, 69.8051, 69.3664 | 69.8051 | 12.8 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 83.8236, 93.9067, 81.3291 | 83.8236 | 15.0 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 3.5023, 3.2771, 2.7102 | 3.2771 | 24.2 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.5904, 8.3742, 8.7445 | 8.7445 | 13.9 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 87.7122, 83.0955, 80.2491 | 83.0955 | 9.0 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 3.2999, 3.1432, 2.7574 | 3.1432 | 17.3 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 8.3302, 8.3521, 8.2406 | 8.3302 | 1.3 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 9.0413, 8.5467, 8.7935 | 8.7935 | 5.6 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 10.1480, 9.4018, 8.4914 | 9.4018 | 17.6 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 5.6389, 5.1566, 5.0454 | 5.1566 | 11.5 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 7.1104, 7.8833, 7.1589 | 7.1589 | 10.8 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 8.1952, 8.0836, 7.3875 | 8.0836 | 10.0 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 9.0057, 8.0861, 7.8845 | 8.0861 | 13.9 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 8.8969, 10.2841, 8.5272 | 8.8969 | 19.7 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 10.9469, 9.7893, 8.8962 | 9.7893 | 20.9 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 8.0929, 7.6900, 6.9311 | 7.6900 | 15.1 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 8.4941, 7.6363, 7.7094 | 7.7094 | 11.1 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 8.7506, 7.8430, 7.4220 | 7.8430 | 16.9 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 8.2004, 7.9377, 7.7442 | 7.9377 | 5.7 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 93.5290, 88.6595, 87.9638 | 88.6595 | 6.3 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 107.6959, 108.7505, 88.9342 | 107.6959 | 18.4 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 5.8764, 5.6639, 5.3236 | 5.6639 | 9.8 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 7.5990, 7.0357, 7.3270 | 7.3270 | 7.7 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 20.0919, 17.9784, 19.7329 | 19.7329 | 10.7 | 3 |
| `6b8c9e0a77ae` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 27.2655, 27.2062, 26.7053 | 27.2062 | 2.1 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 92.7433, 97.4336, 126.7825 | 97.4336 | 34.9 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.2971, 1.8564, 2.2459 | 2.2459 | 19.6 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 6.0258, 5.8743, 5.6248 | 5.8743 | 6.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 118.3854, 121.5728, 122.0927 | 121.5728 | 3.0 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 66.3033, 68.8441, 98.4898 | 68.8441 | 46.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 101.2475, 105.7687, 64.9106 | 101.2475 | 40.4 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 121.3268, 92.6233, 91.2534 | 92.6233 | 32.5 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.2220, 1.6661, 2.5284 | 2.2220 | 38.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 6.6017, 6.9282, 6.9368 | 6.9282 | 4.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 89.3831, 126.0745, 92.9268 | 92.9268 | 39.5 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 2.2624, 2.3887, 1.6140 | 2.2624 | 34.2 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 6.2830, 6.5681, 5.9581 | 6.2830 | 9.7 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 7.0203, 7.5374, 7.3896 | 7.3896 | 7.0 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 7.3026, 8.3265, 7.6066 | 7.6066 | 13.5 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 1.5258, 1.5832, 1.5382 | 1.5382 | 3.7 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 5.8739, 6.5227, 6.2159 | 6.2159 | 10.4 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 6.7233, 6.1805, 6.1913 | 6.1913 | 8.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 5.6241, 6.6765, 6.2755 | 6.2755 | 16.8 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 6.7819, 7.1684, 7.0122 | 7.0122 | 5.5 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 6.9184, 7.8844, 7.7867 | 7.7867 | 12.4 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 5.6936, 6.3761, 5.5235 | 5.6936 | 15.0 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 6.0484, 5.9511, 6.0351 | 6.0351 | 1.6 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 5.5072, 7.0613, 6.1713 | 6.1713 | 25.2 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 5.7661, 6.8154, 6.8657 | 6.8154 | 16.1 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 14.2082, 15.5538, 14.5390 | 14.5390 | 9.3 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 73.6152, 73.4265, 74.0736 | 73.6152 | 0.9 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 1.5917, 1.6063, 1.5234 | 1.5917 | 5.2 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 6.6999, 6.3530, 6.1468 | 6.3530 | 8.7 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 3.6095, 3.7971, 3.5399 | 3.6095 | 7.1 | 3 |
| `7c269cf982fd` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 20.4266, 22.8080, 21.3547 | 21.3547 | 11.2 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 96.1079, 98.6336, 102.2606 | 98.6336 | 6.2 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 1.9096, 1.7659, 1.8214 | 1.8214 | 7.9 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.5134, 8.4416, 8.2564 | 8.4416 | 14.9 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 96.0912, 94.3181, 96.8530 | 96.0912 | 2.6 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 74.1659, 76.2120, 77.6508 | 76.2120 | 4.6 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 77.2926, 78.8539, 78.4964 | 78.4964 | 2.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 97.3263, 103.1642, 101.7018 | 101.7018 | 5.7 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.9447, 1.5243, 1.7431 | 1.7431 | 24.1 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.4829, 8.0655, 7.7667 | 8.0655 | 21.3 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 99.4009, 101.0593, 100.9159 | 100.9159 | 1.6 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 1.6950, 1.6249, 1.8281 | 1.6950 | 12.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 7.8345, 8.5039, 7.8514 | 7.8514 | 8.5 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 8.7344, 8.3040, 8.0490 | 8.3040 | 8.3 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 9.3374, 9.8119, 9.4933 | 9.4933 | 5.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 5.9480, 5.7841, 6.4792 | 5.9480 | 11.7 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 6.9311, 7.7434, 7.8252 | 7.7434 | 11.5 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 7.8210, 7.1551, 7.3689 | 7.3689 | 9.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 8.0383, 7.4379, 7.6248 | 7.6248 | 7.9 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 8.5172, 8.0252, 8.8759 | 8.5172 | 10.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 9.1577, 8.7639, 8.7376 | 8.7639 | 4.8 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 7.1422, 7.4250, 8.5634 | 7.4250 | 19.1 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 7.4771, 7.8509, 7.0869 | 7.4771 | 10.2 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 8.4432, 8.1427, 8.4822 | 8.4432 | 4.0 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 8.4080, 7.9761, 7.9357 | 7.9761 | 5.9 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 68.4239, 73.1574, 72.1884 | 72.1884 | 6.6 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 77.3955, 79.8006, 75.7994 | 77.3955 | 5.2 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 6.3419, 5.9633, 7.1454 | 6.3419 | 18.6 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.5447, 7.5505, 8.3915 | 8.3915 | 11.8 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 20.5445, 22.2108, 21.7533 | 21.7533 | 7.7 | 3 |
| `e917d1a74006` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 25.8026, 28.0347, 25.7612 | 25.8026 | 8.8 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=180000 | 78.6368, 80.5871, 77.9346 | 78.6368 | 3.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=20 | 2.8150, 2.7460, 2.6984 | 2.7460 | 4.2 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer;threshold=2000 | 7.6579, 7.9676, 7.8671 | 7.8671 | 3.9 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 101.6240, 103.7744, 100.8829 | 101.6240 | 2.8 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 63.0618, 67.3520, 64.7019 | 64.7019 | 6.6 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 65.0933, 67.9742, 65.3485 | 65.3485 | 4.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 81.1422, 81.4507, 74.6723 | 81.1422 | 8.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.6038, 2.7759, 3.1138 | 2.7759 | 18.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 8.0165, 7.7897, 7.2819 | 7.7897 | 9.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=180000 | 75.4595, 79.1312, 77.0517 | 77.0517 | 4.8 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=20 | 3.3142, 2.8698, 3.1798 | 3.1798 | 14.0 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=bag;stmt=join_where;threshold=2000 | 7.8363, 8.0216, 8.1008 | 8.0216 | 3.3 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=1% | 7.7865, 7.9709, 8.1518 | 7.9709 | 4.6 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=16;consumer=100% | 7.9741, 8.8635, 8.8872 | 8.8635 | 10.3 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=1% | 5.4924, 5.2445, 5.4597 | 5.4597 | 4.5 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=1;consumer=100% | 7.0748, 6.9410, 7.2126 | 7.0748 | 3.8 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=1% | 6.7126, 7.1358, 7.4319 | 7.1358 | 10.1 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=default;refs=4;consumer=100% | 7.0435, 7.3744, 7.1322 | 7.1322 | 4.6 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=1% | 7.5942, 8.1524, 8.4047 | 8.1524 | 9.9 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=16;consumer=100% | 8.3630, 8.9812, 8.2484 | 8.3630 | 8.8 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=1% | 6.9168, 7.2612, 6.8174 | 6.9168 | 6.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=1;consumer=100% | 6.6983, 7.2869, 7.4503 | 7.2869 | 10.3 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=1% | 7.0160, 7.6119, 6.9621 | 7.0160 | 9.3 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=mat;refs=4;consumer=100% | 6.8674, 7.3757, 6.7784 | 6.8674 | 8.7 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=1% | 71.1265, 91.2518, 76.1795 | 76.1795 | 26.4 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=16;consumer=100% | 89.1630, 93.4224, 102.1203 | 93.4224 | 13.9 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=1% | 4.8997, 5.2130, 5.0769 | 5.0769 | 6.2 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=1;consumer=100% | 6.4239, 6.8842, 6.4237 | 6.4239 | 7.2 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=1% | 17.3076, 18.7986, 15.3370 | 17.3076 | 20.0 | 3 |
| `ed469cd9aa6d` | yugabyte | cluster3 | family=reuse;policy=notmat;refs=4;consumer=100% | 23.4287, 26.4841, 25.3041 | 25.3041 | 12.1 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 104.5941, 112.0943, 106.1600 | 106.1600 | 7.1 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 2.0173, 1.9335, 2.0982 | 2.0173 | 8.2 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 8.3822, 9.4653, 8.9308 | 8.9308 | 12.1 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 99.6518, 125.9413, 106.2035 | 106.2035 | 24.8 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 82.3613, 93.7461, 83.8518 | 83.8518 | 13.6 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 82.0311, 87.5076, 85.2482 | 85.2482 | 6.4 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 108.2737, 109.8973, 108.6185 | 108.6185 | 1.5 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 2.1840, 2.4198, 2.1172 | 2.1840 | 13.9 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 8.7744, 10.8565, 9.3874 | 9.3874 | 22.2 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 106.0982, 110.4733, 111.6797 | 110.4733 | 5.1 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 1.8427, 1.8952, 2.2035 | 1.8952 | 19.0 | 3 |
| `fae5cef41884` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 8.9347, 9.8186, 8.6908 | 8.9347 | 12.6 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 8.8534, 9.5944, 9.3468 | 9.3468 | 7.9 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 9.8040, 11.1998, 11.0585 | 11.0585 | 12.6 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 5.8544, 7.2793, 6.4892 | 6.4892 | 22.0 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 9.0690, 8.4978, 8.9336 | 8.9336 | 6.4 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 8.2388, 8.5868, 9.2765 | 8.5868 | 12.1 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 8.8107, 8.6191, 9.9674 | 8.8107 | 15.3 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 10.0662, 10.6204, 9.0380 | 10.0662 | 15.7 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.5017, 10.0090, 10.3431 | 10.3431 | 4.8 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 8.1215, 8.3394, 8.8567 | 8.3394 | 8.8 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 8.2521, 9.2849, 8.7948 | 8.7948 | 11.7 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 7.9730, 9.8650, 8.4168 | 8.4168 | 22.5 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 9.1206, 9.0289, 9.9249 | 9.1206 | 9.8 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 87.5949, 109.9561, 84.6282 | 87.5949 | 28.9 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 120.6135, 114.4688, 118.6404 | 118.6404 | 5.2 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 7.4205, 7.5097, 7.5682 | 7.5097 | 2.0 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.5940, 8.5587, 8.2308 | 8.5587 | 4.2 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 22.8712, 24.4475, 22.4200 | 22.8712 | 8.9 | 3 |
| `fae5cef41884` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 29.2379, 31.0923, 32.6549 | 31.0923 | 11.0 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 122.2914, 121.1096, 106.6266 | 121.1096 | 12.9 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 1.8522, 2.1198, 1.8870 | 1.8870 | 14.2 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 9.8925, 10.5466, 9.2359 | 9.8925 | 13.2 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 119.3151, 111.5058, 103.5151 | 111.5058 | 14.2 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 90.5090, 95.6558, 80.5408 | 90.5090 | 16.7 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 94.6921, 99.4964, 80.0139 | 94.6921 | 20.6 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 114.5344, 114.8506, 105.5875 | 114.5344 | 8.1 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 1.9874, 2.3058, 1.8585 | 1.9874 | 22.5 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 9.2867, 9.3068, 9.5321 | 9.3068 | 2.6 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 119.2955, 114.3228, 109.3388 | 114.3228 | 8.7 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=join_where;threshold=20 | 2.3514, 2.1925, 1.8150 | 2.1925 | 24.5 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 10.5890, 10.3871, 8.7390 | 10.3871 | 17.8 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 11.5325, 11.1651, 9.0766 | 11.1651 | 22.0 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 10.5784, 10.7672, 9.6206 | 10.5784 | 10.8 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 6.8576, 7.6845, 7.3748 | 7.3748 | 11.2 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 9.1768, 9.4879, 8.5258 | 9.1768 | 10.5 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 9.0018, 8.4036, 8.6189 | 8.6189 | 6.9 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 9.1655, 9.6585, 8.0977 | 9.1655 | 17.0 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 10.4806, 9.9975, 8.6664 | 9.9975 | 18.1 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 10.4087, 10.2790, 9.3057 | 10.2790 | 10.7 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 9.5822, 8.6920, 7.3748 | 8.6920 | 25.4 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 9.0049, 10.2125, 8.1132 | 9.0049 | 23.3 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 9.2065, 9.4010, 8.2582 | 9.2065 | 12.4 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 10.4165, 9.6774, 8.4199 | 9.6774 | 20.6 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 79.7530, 87.6861, 69.8838 | 79.7530 | 22.3 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 89.1317, 101.1018, 83.0334 | 89.1317 | 20.3 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 7.5263, 7.8425, 7.4604 | 7.5263 | 5.1 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 8.9860, 8.3364, 7.7323 | 8.3364 | 15.0 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 25.3471, 24.8587, 21.1959 | 24.8587 | 16.7 | 3 |
| `ffdfe4180cce` | yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 33.2972, 36.8032, 27.6865 | 33.2972 | 27.4 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9224 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 44.8924 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7803 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9670 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.9398 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9484 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.3496 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.0172 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.4131 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 4.5175 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.1353 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9500 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.8689 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0293 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.5492 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.8521 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.2390 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0593 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.3971 | >1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9655 | <1 |
| select-join-v4-topo-colocated-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.4227 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.2239 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 31.9309 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0713 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8705 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 10.2308 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0182 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.7310 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.0088 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9978 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.9610 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.8826 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.8150 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.8378 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9112 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6534 | <1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0947 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.6587 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0372 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 1.7259 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0218 | >1 |
| select-join-v4-topo-colocated-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.4208 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.1287 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 31.5676 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7878 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0238 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 16.4793 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1024 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.7318 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9316 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9322 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.8982 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0644 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0090 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0406 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9746 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6976 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9869 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.7047 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.1191 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.3668 | >1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9722 | <1 |
| select-join-v4-topo-colocated-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.0814 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0153 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 29.3071 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9821 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9591 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 16.1144 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0507 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0376 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3245 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.3574 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.7316 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0432 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0297 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.1406 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.8191 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.5369 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0252 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.6319 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9660 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.0239 | >1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9474 | <1 |
| select-join-v4-topo-colocated-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.0806 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.7771 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 28.8208 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.6975 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8944 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 16.1034 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0548 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.7728 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9643 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7347 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 4.0273 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0146 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9124 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9740 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.1425 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.6144 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0208 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4162 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9510 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 2.0636 | >1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9469 | <1 |
| select-join-v4-topo-colocated-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.8184 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.3915 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 61.0218 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.5665 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9441 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 10.8945 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1643 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.3643 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3139 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9820 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 3.5909 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9904 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9709 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9889 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9968 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 0.5717 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0941 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4029 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9489 | <1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 1.9675 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0237 | >1 |
| select-join-v4-topo-colocated-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.7380 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.1080 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 22.1399 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0613 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0277 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 9.4005 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1513 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9943 | <1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3411 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9557 | <1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.4352 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0421 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.1946 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0687 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0678 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.4517 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9106 | <1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.0276 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9840 | <1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 10.3447 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0787 | >1 |
| select-join-v4-topo-distributed-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.6125 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9337 | <1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 23.8041 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0426 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0014 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.3578 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0027 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0383 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3464 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1301 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.4913 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0984 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9687 | <1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.8925 | <1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9702 | <1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.2241 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9816 | <1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.3646 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.2033 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 10.3736 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0412 | >1 |
| select-join-v4-topo-distributed-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 11.5670 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.2867 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 25.1273 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9829 | <1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.0823 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.4176 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0611 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0751 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.4252 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0135 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.3737 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.0551 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0769 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0235 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0047 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.6711 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9822 | <1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.3871 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9697 | <1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 10.0033 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0477 | >1 |
| select-join-v4-topo-distributed-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.4734 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8494 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 19.0278 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7856 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9772 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.3067 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0230 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0421 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3467 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0753 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2594 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.8921 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9468 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9080 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0452 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.5784 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9750 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.3263 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9753 | <1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 9.1346 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0488 | >1 |
| select-join-v4-topo-distributed-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 11.1815 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9568 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 23.4691 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9673 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9933 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.4739 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9711 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0184 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3114 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0293 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.3845 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9940 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0498 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9918 | <1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0667 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.6344 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0002 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.5913 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0228 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 11.4482 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0133 | >1 |
| select-join-v4-topo-distributed-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.5401 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8486 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 20.3477 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9792 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9712 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 8.0669 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8989 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0115 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.3093 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9691 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2487 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9299 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0330 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.8906 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9368 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.0637 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9504 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.5479 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0310 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 9.3451 | >1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9281 | <1 |
| select-join-v4-topo-distributed-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 11.4907 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.7877 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 38.4919 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.8452 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9342 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.9425 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8770 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0251 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0002 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9601 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.3973 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0975 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9813 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9792 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0227 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.8158 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.1365 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.6329 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9088 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 6.9155 | >1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9840 | <1 |
| select-join-v4-topo-equal-b1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.4258 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9668 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 43.6288 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0517 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0154 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.5788 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8960 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0594 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9754 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0046 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1311 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0206 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0764 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.8786 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1187 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9581 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0020 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.8105 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.8954 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.8536 | >1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9547 | <1 |
| select-join-v4-topo-equal-b2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.3898 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0397 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 44.3749 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0239 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0569 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1559 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0908 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9752 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9467 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9657 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.0000 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0116 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9516 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9069 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9581 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.4592 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0398 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4190 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9548 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.6994 | >1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9673 | <1 |
| select-join-v4-topo-equal-b3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.6308 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1266 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 43.7565 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1473 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.2143 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.8657 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.2104 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9669 | <1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9667 | <1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9791 | <1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2008 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0662 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0788 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.2328 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0796 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6269 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0460 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2099 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9751 | <1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.8338 | >1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9808 | <1 |
| select-join-v4-topo-equal-i1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.2887 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0868 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 46.9017 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9381 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9927 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.2727 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9484 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9760 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9333 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0208 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2837 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0310 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0139 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9751 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1380 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 3.1042 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0724 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.7692 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9664 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.8099 | >1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.8932 | <1 |
| select-join-v4-topo-equal-i2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 8.1330 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9963 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 42.4757 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9535 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0516 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.9978 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9892 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0133 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9597 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0078 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.3217 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1028 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9056 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0724 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1511 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9520 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0408 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3786 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.1027 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.9686 | >1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9204 | <1 |
| select-join-v4-topo-equal-i3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 7.9845 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0947 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 44.6949 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1852 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9382 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.1812 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9821 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9858 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9392 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0205 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.3872 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.2675 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9099 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9476 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9677 | <1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.7760 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0352 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.3184 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.1370 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.8939 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0712 | >1 |
| select-join-v4-topo-single-b1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 12.3025 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0202 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 49.4653 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.2768 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9640 | <1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 8.9124 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1057 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0147 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.1400 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9948 | <1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1456 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0316 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0926 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0072 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.1489 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.8471 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0475 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.6074 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.1069 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 11.4605 | >1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.8937 | <1 |
| select-join-v4-topo-single-b2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.2206 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9522 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 38.0535 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9608 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0276 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.8090 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0802 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9506 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9510 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9726 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.3648 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.1663 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9845 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9213 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9073 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.4168 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9957 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.2762 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9670 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.0543 | >1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9353 | <1 |
| select-join-v4-topo-single-b3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 10.7284 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 0.9143 | <1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 39.7161 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9057 | <1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9824 | <1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.3840 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0290 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0909 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0037 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.1141 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1954 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9913 | <1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0096 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0587 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 1.0003 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9819 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0155 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.5638 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9095 | <1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.8748 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0128 | >1 |
| select-join-v4-topo-single-i1-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.1185 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.1552 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 50.8956 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0990 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.1280 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 9.3785 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0526 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0366 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.8915 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.8523 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.1650 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9679 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9593 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0627 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9738 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.6087 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9924 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4813 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9531 | <1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 11.3738 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0282 | >1 |
| select-join-v4-topo-single-i2-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.5289 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0456 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 45.7852 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.1162 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9843 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 10.2190 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1437 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9862 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.0069 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0135 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=1% | mat vs default | 1.2215 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9533 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.1161 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0242 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.8832 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.2443 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=100% | mat vs default | 0.9643 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.6184 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.8637 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 7.1457 | >1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9554 | <1 |
| select-join-v4-topo-single-i3-20261005 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.1888 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9727 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 62.7943 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.5268 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9205 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 7.4349 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9136 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0636 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9753 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.3106 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.2665 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.8764 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0583 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.0185 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.0174 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.8331 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.1034 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.2530 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9371 | <1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 7.8715 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.0317 | >1 |
| select-join-v4-topo-tablet-b1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.7630 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.8904 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 41.1328 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.0130 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.1629 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 12.1031 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.8610 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.2670 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.2322 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.2064 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.1195 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9729 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0301 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.8491 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9508 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.6896 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9954 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.6465 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.8921 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.0871 | >1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9463 | <1 |
| select-join-v4-topo-tablet-b2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.0680 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9432 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 43.3218 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.9595 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8136 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 10.3796 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9864 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9751 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9098 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.0050 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.4221 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.2543 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9446 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 0.9592 | <1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 1.1527 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.7907 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0616 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.0951 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0565 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.5421 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.1805 | >1 |
| select-join-v4-topo-tablet-b3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.6625 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.4553 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 64.4035 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 1.4508 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 1.1493 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 13.5023 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1402 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.0076 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.9387 | <1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9873 | <1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.0803 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9067 | <1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 1.0279 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.1566 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9837 | <1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.9390 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.9619 | <1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 2.9231 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0003 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 7.5493 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 1.1387 | >1 |
| select-join-v4-topo-tablet-i1-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 9.0285 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 0.9691 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 30.6491 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 5.5560 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.8154 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 11.1887 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0612 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 0.9555 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 0.7173 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 0.7337 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.3704 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 1.2086 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9782 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.1027 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.9460 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.2932 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 0.8491 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4107 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 1.0250 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 7.1036 | >1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9808 | <1 |
| select-join-v4-topo-tablet-i2-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 8.8548 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer vs join_where | 1.0620 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_mat vs join_where | 44.6956 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@20 | cte_joined_outer_notmat vs join_where | 0.7552 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer vs join_where | 0.9011 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_mat vs join_where | 10.7912 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0011 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer vs join_where | 1.3305 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_mat vs join_where | 1.2782 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | bag@180000 | cte_joined_outer_notmat vs join_where | 1.3131 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | mat vs default | 1.1702 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=1% | notmat vs default | 0.9549 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | mat vs default | 0.9472 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=1 c=100% | notmat vs default | 1.1106 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | mat vs default | 0.8867 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=1% | notmat vs default | 2.3389 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | mat vs default | 1.0017 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=4 c=100% | notmat vs default | 3.4005 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | mat vs default | 0.9970 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=1% | notmat vs default | 8.2840 | >1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | mat vs default | 0.9329 | <1 |
| select-join-v4-topo-tablet-i3-20261005 | yugabyte | cluster3 | reuse n=16 c=100% | notmat vs default | 10.8696 | >1 |

## Per-run endpoint, writer, retry and throttling evidence

| run | topology | phase | index arm | placement | tablets (donation, rank idx) | donation leaders | endpoints (ops) | writer offered/completed/failed | writer rate/s | reader retries | overlap | node throttled periods (usec) | client throttled periods (usec) |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| select-join-v4-topo-colocated-b1-20261005 | cluster3 | blended | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 122/122/0 | 10.06 | 0 | 59/210 | yb-n1 5 (737945); yb-n2 5 (584832); yb-n3 5 (107882) | 92 (3535021) [periods 268] |
| select-join-v4-topo-colocated-b2-20261005 | cluster3 | blended | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 113/113/0 | 10.08 | 0 | 58/210 | yb-n1 5 (577296); yb-n2 5 (350507); yb-n3 2 (20643) | 83 (2948089) [periods 258] |
| select-join-v4-topo-colocated-b3-20261005 | cluster3 | blended | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 125/125/0 | 10.07 | 0 | 56/210 | yb-n1 6 (730253); yb-n2 5 (613130); yb-n3 1 (39795) | 94 (4919178) [periods 275] |
| select-join-v4-topo-colocated-i1-20261005 | cluster3 | isolated | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n2:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 4 (544740); yb-n2 8 (1135954); yb-n3 0 (0) | 81 (3918876) [periods 251] |
| select-join-v4-topo-colocated-i2-20261005 | cluster3 | isolated | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 7 (1167843); yb-n2 6 (828293); yb-n3 5 (138350) | 91 (4489964) [periods 261] |
| select-join-v4-topo-colocated-i3-20261005 | cluster3 | isolated | yb-range-rank | colocated | donation=1, donation_rank_idx=1 | yb-n2:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 5 (643974); yb-n2 10 (1569323); yb-n3 0 (0) | 81 (3498020) [periods 247] |
| select-join-v4-topo-distributed-b1-20261005 | cluster3 | blended | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (169/168/168) | 169/169/0 | 10.01 | 1 | 66/211 | yb-n1 26 (1440554); yb-n2 6 (1142087); yb-n3 5 (126167) | 161 (6353847) [periods 347] |
| select-join-v4-topo-distributed-b2-20261005 | cluster3 | blended | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (169/169/168) | 162/162/0 | 10.02 | 2 | 61/212 | yb-n1 9 (1040129); yb-n2 5 (879690); yb-n3 3 (138115) | 144 (5823602) [periods 325] |
| select-join-v4-topo-distributed-b3-20261005 | cluster3 | blended | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (169/169/168) | 152/152/0 | 10.03 | 2 | 62/212 | yb-n1 6 (789477); yb-n2 6 (898858); yb-n3 0 (0) | 139 (5200375) [periods 313] |
| select-join-v4-topo-distributed-i1-20261005 | cluster3 | isolated | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 10 (895807); yb-n2 4 (664926); yb-n3 1 (2385) | 128 (4130668) [periods 277] |
| select-join-v4-topo-distributed-i2-20261005 | cluster3 | isolated | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 15 (1301458); yb-n2 6 (863774); yb-n3 0 (0) | 138 (5229216) [periods 308] |
| select-join-v4-topo-distributed-i3-20261005 | cluster3 | isolated | yb-range-rank | distributed | donation=3, donation_rank_idx=3 | yb-n1:9100,yb-n2:9100,yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 4 (646315); yb-n2 7 (934653); yb-n3 5 (264570) | 133 (5526846) [periods 292] |
| select-join-v4-topo-equal-b1-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 155/155/0 | 10.03 | 0 | 71/210 | yb-single 0 (0) | 134 (5710025) [periods 305] |
| select-join-v4-topo-equal-b2-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 157/157/0 | 10.00 | 1 | 71/211 | yb-single 0 (0) | 120 (4535597) [periods 301] |
| select-join-v4-topo-equal-b3-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 137/137/0 | 10.05 | 0 | 67/210 | yb-single 0 (0) | 113 (3388981) [periods 277] |
| select-join-v4-topo-equal-i1-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 102 (4339856) [periods 267] |
| select-join-v4-topo-equal-i2-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 111 (4428811) [periods 264] |
| select-join-v4-topo-equal-i3-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 0 (0) | 113 (3923859) [periods 269] |
| select-join-v4-topo-single-b1-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 146/146/0 | 10.03 | 0 | 69/210 | yb-single 71 (5308192) | 116 (4144397) [periods 285] |
| select-join-v4-topo-single-b2-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 155/155/0 | 10.04 | 0 | 72/210 | yb-single 92 (5185062) | 125 (5884815) [periods 314] |
| select-join-v4-topo-single-b3-20261005 | single | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 152/152/0 | 10.04 | 0 | 72/210 | yb-single 83 (5824482) | 125 (4978502) [periods 301] |
| select-join-v4-topo-single-i1-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 82 (5214799) | 110 (3733757) [periods 298] |
| select-join-v4-topo-single-i2-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 79 (6376858) | 109 (3557702) [periods 301] |
| select-join-v4-topo-single-i3-20261005 | single | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-single:9100 | single | 0/0/0 | 0.00 | 0 | 0/0 | yb-single 79 (5674243) | 119 (3940006) [periods 295] |
| select-join-v4-topo-tablet-b1-20261005 | cluster3 | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n2:9100 | 3 (168/168/168) | 167/167/0 | 10.04 | 0 | 68/210 | yb-n1 4 (529102); yb-n2 13 (2062689); yb-n3 0 (0) | 125 (4478417) [periods 329] |
| select-join-v4-topo-tablet-b2-20261005 | cluster3 | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 175/175/0 | 10.02 | 0 | 76/210 | yb-n1 6 (562946); yb-n2 5 (709435); yb-n3 10 (331576) | 125 (5125823) [periods 346] |
| select-join-v4-topo-tablet-b3-20261005 | cluster3 | blended | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n2:9100 | 3 (168/168/168) | 162/162/0 | 10.02 | 0 | 72/210 | yb-n1 5 (496966); yb-n2 17 (1563387); yb-n3 0 (0) | 119 (4324233) [periods 321] |
| select-join-v4-topo-tablet-i1-20261005 | cluster3 | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 4 (575890); yb-n2 6 (764364); yb-n3 5 (49224) | 108 (4874043) [periods 327] |
| select-join-v4-topo-tablet-i2-20261005 | cluster3 | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n3:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 4 (426573); yb-n2 5 (634129); yb-n3 10 (226197) | 116 (4947765) [periods 307] |
| select-join-v4-topo-tablet-i3-20261005 | cluster3 | isolated | yb-range-rank | noncolocated-single-tablet | donation=1, donation_rank_idx=1 | yb-n2:9100 | 3 (168/168/168) | 0/0/0 | 0.00 | 0 | 0/0 | yb-n1 4 (521500); yb-n2 17 (1672194); yb-n3 0 (0) | 115 (4670066) [periods 311] |

## Recorded limits

- [select-join-v4-topo-colocated-b1-20261005] explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- [select-join-v4-topo-colocated-b1-20261005] one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- [select-join-v4-topo-colocated-b1-20261005] one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- [select-join-v4-topo-colocated-b1-20261005] cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- [select-join-v4-topo-colocated-b1-20261005] readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- [select-join-v4-topo-colocated-b1-20261005] physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- [select-join-v4-topo-colocated-b1-20261005] YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- [select-join-v4-topo-colocated-b1-20261005] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
- [select-join-v4-topo-colocated-i1-20261005] isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- [select-join-v4-topo-distributed-b1-20261005] distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- [select-join-v4-topo-equal-b1-20261005] one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
