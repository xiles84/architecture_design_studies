# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)

_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._

**TL;DR:** 6 whole runs summarised; each run's own median is shown and never merged across engines, topologies or thresholds.

## Runs

| run | engine | topology | scale | fact rows | seed |
|---|---|---|---|---|---|
| select-join-pg1-20261004 | postgres | single | cte | 200000 | 42 |
| select-join-pg2-20261004 | postgres | single | cte | 200000 | 42 |
| select-join-pg3-20261004 | postgres | single | cte | 200000 | 42 |
| select-join-yb1-20261004 | yugabyte | single | cte | 200000 | 42 |
| select-join-yb2-20261004 | yugabyte | single | cte | 200000 | 42 |
| select-join-yb3-20261004 | yugabyte | single | cte | 200000 | 42 |

## Cells across runs (per engine and topology)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | family=bag;stmt=comma_join;threshold=100000 | 37.6825, 37.0255, 60.0191 | 37.6825 | 61.0 | 3 |
| postgres | single | family=bag;stmt=comma_join;threshold=180000 | 45.5845, 46.1841, 74.2699 | 46.1841 | 62.1 | 3 |
| postgres | single | family=bag;stmt=comma_join;threshold=20 | 0.4979, 0.5505, 0.5099 | 0.5099 | 10.3 | 3 |
| postgres | single | family=bag;stmt=comma_join;threshold=200 | 0.5782, 0.6202, 0.9139 | 0.6202 | 54.1 | 3 |
| postgres | single | family=bag;stmt=comma_join;threshold=2000 | 1.3679, 1.6504, 2.2905 | 1.6504 | 55.9 | 3 |
| postgres | single | family=bag;stmt=comma_join;threshold=20000 | 6.2381, 5.5945, 7.3755 | 6.2381 | 28.5 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=100000 | 37.8795, 39.8509, 60.2255 | 39.8509 | 56.1 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=180000 | 48.0588, 49.7117, 60.0565 | 49.7117 | 24.1 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=20 | 0.4700, 0.4103, 0.5586 | 0.4700 | 31.6 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=200 | 0.6489, 0.5609, 0.8483 | 0.6489 | 44.3 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=2000 | 1.7570, 1.4668, 2.1902 | 1.7570 | 41.2 | 3 |
| postgres | single | family=bag;stmt=cte_filtered;threshold=20000 | 5.9134, 5.5818, 7.4979 | 5.9134 | 32.4 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 41.9282, 44.9175, 62.9446 | 44.9175 | 46.8 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 48.8781, 48.4859, 71.2601 | 48.8781 | 46.6 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 0.3920, 0.6016, 0.6535 | 0.6016 | 43.5 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 0.5953, 0.5746, 0.8686 | 0.5953 | 49.4 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 1.5628, 1.5003, 2.2407 | 1.5628 | 47.4 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 6.5758, 6.1904, 7.8886 | 6.5758 | 25.8 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_mat;threshold=20;position=last | 0.3958 | 0.3958 | 0.0 | 1 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 38.4509, 37.4683, 58.0338 | 38.4509 | 53.5 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 45.6820, 44.6485, 51.6417 | 45.6820 | 15.3 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 0.4427, 0.5527, 0.5953 | 0.5527 | 27.6 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 0.8553, 0.8413, 0.9135 | 0.8553 | 8.4 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 1.6232, 1.3624, 2.1657 | 1.6232 | 49.5 | 3 |
| postgres | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 5.8440, 5.8829, 7.0544 | 5.8829 | 20.6 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 37.6850, 36.4659, 58.9167 | 37.6850 | 59.6 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 45.0621, 47.2707, 66.3503 | 47.2707 | 45.0 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=20 | 0.5048, 0.3993, 0.5549 | 0.5048 | 30.8 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=200 | 0.5446, 0.8494, 0.7998 | 0.7998 | 38.1 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 1.4107, 1.2700, 1.8898 | 1.4107 | 43.9 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 6.1904, 5.8073, 7.3033 | 6.1904 | 24.2 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 53.5020, 53.1524, 72.6632 | 53.5020 | 36.5 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 63.7415, 64.7961, 84.1167 | 64.7961 | 31.4 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 36.9053, 37.0778, 47.3537 | 37.0778 | 28.2 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 40.2530, 40.4844, 45.4371 | 40.4844 | 12.8 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 37.5496, 38.9063, 51.8563 | 38.9063 | 36.8 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 42.0132, 42.7557, 52.3460 | 42.7557 | 24.2 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | 36.2864 | 36.2864 | 0.0 | 1 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 36.8440, 36.7364, 61.9921 | 36.8440 | 68.5 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 42.5202, 46.9925, 58.6311 | 46.9925 | 34.3 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 0.3844, 0.3957, 0.6638 | 0.3957 | 70.6 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 0.7009, 0.7032, 1.1907 | 0.7032 | 69.7 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 1.5546, 1.6029, 1.9320 | 1.6029 | 23.5 | 3 |
| postgres | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 5.6211, 6.1446, 8.6984 | 6.1446 | 50.1 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=100000 | 37.1299, 39.8349, 49.7021 | 39.8349 | 31.6 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=180000 | 44.1086, 42.5508, 67.8667 | 44.1086 | 57.4 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=20 | 0.6337, 0.5796, 0.5433 | 0.5796 | 15.6 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=200 | 0.7846, 0.6290, 1.1272 | 0.7846 | 63.5 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=2000 | 1.7095, 1.5771, 1.9151 | 1.7095 | 19.8 | 3 |
| postgres | single | family=bag;stmt=filtered_derived;threshold=20000 | 5.7220, 5.9304, 7.6269 | 5.9304 | 32.1 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=100000 | 34.9793, 37.8893, 53.4579 | 37.8893 | 48.8 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=180000 | 45.3185, 40.9529, 62.4330 | 45.3185 | 47.4 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=20 | 0.3605, 0.3628, 0.4731 | 0.3628 | 31.0 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=200 | 0.7784, 0.8235, 0.7535 | 0.7784 | 9.0 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=2000 | 1.3183, 1.6320, 1.9229 | 1.6320 | 37.0 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=20000 | 5.5920, 5.8614, 9.1660 | 5.8614 | 61.0 | 3 |
| postgres | single | family=bag;stmt=join_on;threshold=200;position=last | 0.6231 | 0.6231 | 0.0 | 1 |
| postgres | single | family=bag;stmt=join_where;threshold=100000 | 38.7161, 37.6218, 54.1921 | 38.7161 | 42.8 | 3 |
| postgres | single | family=bag;stmt=join_where;threshold=180000 | 46.7475, 47.6480, 68.7031 | 47.6480 | 46.1 | 3 |
| postgres | single | family=bag;stmt=join_where;threshold=20 | 0.4177, 0.3890, 0.5173 | 0.4177 | 30.7 | 3 |
| postgres | single | family=bag;stmt=join_where;threshold=200 | 0.6714, 0.8207, 0.9129 | 0.8207 | 29.4 | 3 |
| postgres | single | family=bag;stmt=join_where;threshold=2000 | 1.4119, 1.3901, 1.9756 | 1.4119 | 41.5 | 3 |
| postgres | single | family=bag;stmt=join_where;threshold=20000 | 6.1398, 6.0667, 7.6702 | 6.1398 | 26.1 | 3 |
| postgres | single | family=reuse;policy=default;refs=16;consumer=1% | 3.0468, 3.1462, 4.0233 | 3.1462 | 31.0 | 3 |
| postgres | single | family=reuse;policy=default;refs=16;consumer=100% | 3.3187, 3.1885, 4.6956 | 3.3187 | 45.4 | 3 |
| postgres | single | family=reuse;policy=default;refs=1;consumer=1% | 0.4576, 0.6543, 0.7223 | 0.6543 | 40.5 | 3 |
| postgres | single | family=reuse;policy=default;refs=1;consumer=100% | 2.2153, 1.9196, 2.6865 | 2.2153 | 34.6 | 3 |
| postgres | single | family=reuse;policy=default;refs=2;consumer=1% | 2.3980, 2.1791, 3.0633 | 2.3980 | 36.9 | 3 |
| postgres | single | family=reuse;policy=default;refs=2;consumer=100% | 2.3958, 2.3629, 2.9653 | 2.3958 | 25.1 | 3 |
| postgres | single | family=reuse;policy=default;refs=4;consumer=1% | 2.5367, 2.3796, 3.0217 | 2.5367 | 25.3 | 3 |
| postgres | single | family=reuse;policy=default;refs=4;consumer=100% | 2.4807, 2.4233, 3.3216 | 2.4807 | 36.2 | 3 |
| postgres | single | family=reuse;policy=default;refs=8;consumer=1% | 2.8126, 2.7135, 3.4752 | 2.8126 | 27.1 | 3 |
| postgres | single | family=reuse;policy=default;refs=8;consumer=100% | 2.9185, 3.1735, 3.9767 | 3.1735 | 33.3 | 3 |
| postgres | single | family=reuse;policy=mat;refs=16;consumer=1% | 3.2646, 3.2338, 3.8049 | 3.2646 | 17.5 | 3 |
| postgres | single | family=reuse;policy=mat;refs=16;consumer=100% | 3.3909, 3.6705, 4.4203 | 3.6705 | 28.0 | 3 |
| postgres | single | family=reuse;policy=mat;refs=1;consumer=1% | 2.1583, 2.2570, 3.2872 | 2.2570 | 50.0 | 3 |
| postgres | single | family=reuse;policy=mat;refs=1;consumer=100% | 2.4804, 2.4236, 3.0881 | 2.4804 | 26.8 | 3 |
| postgres | single | family=reuse;policy=mat;refs=2;consumer=1% | 2.5836, 2.6064, 3.8939 | 2.6064 | 50.3 | 3 |
| postgres | single | family=reuse;policy=mat;refs=2;consumer=100% | 2.6858, 2.2794, 3.1538 | 2.6858 | 32.6 | 3 |
| postgres | single | family=reuse;policy=mat;refs=4;consumer=1% | 2.6090, 2.2673, 3.1626 | 2.6090 | 34.3 | 3 |
| postgres | single | family=reuse;policy=mat;refs=4;consumer=100% | 2.4531, 2.7523, 4.4975 | 2.7523 | 74.3 | 3 |
| postgres | single | family=reuse;policy=mat;refs=8;consumer=1% | 2.4965, 2.5571, 3.3360 | 2.5571 | 32.8 | 3 |
| postgres | single | family=reuse;policy=mat;refs=8;consumer=100% | 3.2429, 2.9563, 3.5456 | 3.2429 | 18.2 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=16;consumer=1% | 2.2689, 2.5444, 2.9672 | 2.5444 | 27.4 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=16;consumer=100% | 14.2454, 14.6995, 19.2252 | 14.6995 | 33.9 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=1;consumer=1% | 0.4459, 0.4827, 0.6402 | 0.4827 | 40.3 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=1;consumer=100% | 2.3161, 2.3477, 2.8040 | 2.3477 | 20.8 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=2;consumer=1% | 0.7016, 0.5664, 0.7837 | 0.7016 | 31.0 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=2;consumer=100% | 4.9255, 4.8402, 5.3479 | 4.9255 | 10.3 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=4;consumer=1% | 1.0949, 0.8563, 1.3533 | 1.0949 | 45.4 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=4;consumer=100% | 7.0942, 7.0597, 8.9609 | 7.0942 | 26.8 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=8;consumer=1% | 1.3449, 1.3824, 1.8828 | 1.3824 | 38.9 | 3 |
| postgres | single | family=reuse;policy=notmat;refs=8;consumer=100% | 8.9133, 9.4583, 14.2132 | 9.4583 | 56.0 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=100000 | 89.8597, 94.6713, 82.9680 | 89.8597 | 13.0 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=180000 | 111.7094, 114.8336, 104.8518 | 111.7094 | 8.9 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=20 | 56.1530, 56.2012, 51.8888 | 56.1530 | 7.7 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=200 | 52.0503, 59.1745, 51.5427 | 52.0503 | 14.7 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=2000 | 51.7476, 59.6796, 50.2710 | 51.7476 | 18.2 | 3 |
| yugabyte | single | family=bag;stmt=comma_join;threshold=20000 | 68.5227, 69.9753, 62.1412 | 68.5227 | 11.4 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=100000 | 90.2188, 93.8807, 82.5497 | 90.2188 | 12.6 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=180000 | 107.3187, 113.9327, 102.7894 | 107.3187 | 10.4 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=20 | 53.4507, 55.3352, 48.9501 | 53.4507 | 11.9 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=200 | 52.4404, 58.5707, 48.8216 | 52.4404 | 18.6 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=2000 | 52.9426, 57.1560, 50.4312 | 52.9426 | 12.7 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered;threshold=20000 | 66.8677, 70.3850, 60.1950 | 66.8677 | 15.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=100000 | 97.3519, 101.3366, 87.8927 | 97.3519 | 13.8 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=180000 | 114.6110, 130.7327, 108.7353 | 114.6110 | 19.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20 | 52.6065, 55.1612, 47.0736 | 52.6065 | 15.4 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=200 | 53.3720, 55.2458, 48.6444 | 53.3720 | 12.4 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=2000 | 54.3065, 57.4772, 48.3419 | 54.3065 | 16.8 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000 | 70.4617, 71.3000, 61.9860 | 70.4617 | 13.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | 96.0397 | 96.0397 | 0.0 | 1 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=100000 | 86.9815, 94.8852, 83.2714 | 86.9815 | 13.4 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=180000 | 109.1626, 118.1213, 104.2315 | 109.1626 | 12.7 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20 | 55.8474, 56.1700, 50.4773 | 55.8474 | 10.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=200 | 53.3432, 56.9247, 49.4430 | 53.3432 | 14.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=2000 | 53.7409, 58.5656, 49.0466 | 53.7409 | 17.7 | 3 |
| yugabyte | single | family=bag;stmt=cte_filtered_notmat;threshold=20000 | 66.4306, 70.5793, 61.4589 | 66.4306 | 13.7 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=100000 | 90.1083, 96.6752, 81.3736 | 90.1083 | 17.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=180000 | 102.6674, 122.3637, 103.9512 | 103.9512 | 18.9 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20 | 55.8498, 57.1734, 50.3063 | 55.8498 | 12.3 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=200 | 56.8219, 58.7040, 50.1562 | 56.8219 | 15.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=2000 | 60.4125, 56.2968, 51.4149 | 56.2968 | 16.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer;threshold=20000 | 66.2433, 68.1864, 62.0711 | 66.2433 | 9.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=100000 | 95.3078, 101.4762, 88.2493 | 95.3078 | 13.9 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000 | 104.4287, 108.2137, 98.1551 | 104.4287 | 9.6 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=180000;position=last | 99.1604 | 99.1604 | 0.0 | 1 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20 | 87.0800, 86.1778, 76.8114 | 86.1778 | 11.9 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=200 | 88.6169, 89.9664, 76.8135 | 88.6169 | 14.8 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=2000 | 87.4509, 91.8593, 75.9808 | 87.4509 | 18.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_mat;threshold=20000 | 86.8973, 94.1146, 81.2126 | 86.8973 | 14.8 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | 90.2313, 93.4475, 83.5110 | 90.2313 | 11.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | 108.8500, 116.5785, 105.3809 | 108.8500 | 10.3 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20 | 56.1114, 54.7682, 48.8650 | 54.7682 | 13.2 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=200 | 52.0684, 57.9382, 51.3568 | 52.0684 | 12.6 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | 57.2107, 57.1141, 51.4716 | 57.1141 | 10.0 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | 68.3329, 70.8516, 61.7382 | 68.3329 | 13.3 | 3 |
| yugabyte | single | family=bag;stmt=cte_joined_outer_notmat;threshold=20000;position=last | 58.8255 | 58.8255 | 0.0 | 1 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=100000 | 93.6500, 94.9401, 82.6168 | 93.6500 | 13.2 | 3 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=180000 | 110.6885, 117.8714, 102.5922 | 110.6885 | 13.8 | 3 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=20 | 51.9190, 56.2280, 54.3175 | 54.3175 | 7.9 | 3 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=200 | 58.7998, 55.3489, 49.4475 | 55.3489 | 16.9 | 3 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=2000 | 54.4900, 57.0162, 51.0314 | 54.4900 | 11.0 | 3 |
| yugabyte | single | family=bag;stmt=filtered_derived;threshold=20000 | 65.5495, 71.1113, 61.3544 | 65.5495 | 14.9 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=100000 | 86.5841, 100.1651, 82.6664 | 86.5841 | 20.2 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=180000 | 111.9913, 125.0222, 102.6127 | 111.9913 | 20.0 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=20 | 51.3924, 58.2015, 55.5802 | 55.5802 | 12.3 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=200 | 52.2678, 55.9788, 51.3271 | 52.2678 | 8.9 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=2000 | 58.4189, 56.5501, 50.8117 | 56.5501 | 13.5 | 3 |
| yugabyte | single | family=bag;stmt=join_on;threshold=20000 | 63.5956, 70.6003, 60.1219 | 63.5956 | 16.5 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=100000 | 89.8912, 89.9328, 83.3312 | 89.8912 | 7.3 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=180000 | 121.1290, 119.3936, 105.8799 | 119.3936 | 12.8 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=20 | 54.1855, 56.9529, 50.2806 | 54.1855 | 12.3 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=200 | 53.6080, 56.2774, 48.5288 | 53.6080 | 14.5 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=2000 | 54.0418, 57.6083, 49.5256 | 54.0418 | 15.0 | 3 |
| yugabyte | single | family=bag;stmt=join_where;threshold=20000 | 69.6583, 73.9714, 62.0317 | 69.6583 | 17.1 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=16;consumer=1% | 54.8308, 58.1748, 53.9796 | 54.8308 | 7.7 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=16;consumer=100% | 52.4722, 59.5489, 52.8587 | 52.8587 | 13.4 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=1;consumer=1% | 56.3565, 60.8007, 54.5585 | 56.3565 | 11.1 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=1;consumer=100% | 55.8016, 59.6535, 52.2319 | 55.8016 | 13.3 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=2;consumer=1% | 51.8952, 56.9339, 52.5271 | 52.5271 | 9.6 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=2;consumer=100% | 56.1597, 57.8339, 48.3360 | 56.1597 | 16.9 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=4;consumer=1% | 54.1220, 60.2149, 51.7899 | 54.1220 | 15.6 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=4;consumer=100% | 52.9323, 56.8766, 51.9626 | 52.9323 | 9.3 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=8;consumer=1% | 55.4008, 56.3602, 49.7483 | 55.4008 | 11.9 | 3 |
| yugabyte | single | family=reuse;policy=default;refs=8;consumer=100% | 57.1514, 55.3011, 50.9801 | 55.3011 | 11.2 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=16;consumer=1% | 54.3485, 59.2150, 50.0097 | 54.3485 | 16.9 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=16;consumer=100% | 56.4226, 58.3002, 53.1911 | 56.4226 | 9.1 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=1;consumer=1% | 51.3329, 55.9581, 50.0198 | 51.3329 | 11.6 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=1;consumer=100% | 55.8041, 57.8953, 50.8898 | 55.8041 | 12.6 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=2;consumer=1% | 55.3564, 56.3527, 51.3200 | 55.3564 | 9.1 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=2;consumer=100% | 55.0306, 56.6137, 49.3620 | 55.0306 | 13.2 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=4;consumer=1% | 51.5940, 55.3893, 51.0538 | 51.5940 | 8.4 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=4;consumer=100% | 59.0504, 57.5499, 51.9624 | 57.5499 | 12.3 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=8;consumer=1% | 56.6268, 58.1511, 50.6843 | 56.6268 | 13.2 | 3 |
| yugabyte | single | family=reuse;policy=mat;refs=8;consumer=100% | 51.4913, 58.1307, 49.7646 | 51.4913 | 16.2 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=1% | 490.0560, 550.9748, 496.9699 | 496.9699 | 12.3 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=16;consumer=100% | 585.2708, 589.9626, 507.6440 | 585.2708 | 14.1 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=1% | 56.9631, 58.3730, 51.9662 | 56.9631 | 11.2 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=1;consumer=100% | 56.8319, 57.6320, 53.0314 | 56.8319 | 8.1 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=1% | 94.5241, 101.3628, 88.3428 | 94.5241 | 13.8 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=2;consumer=100% | 97.2606, 96.0695, 87.3859 | 96.0695 | 10.3 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=1% | 157.2088, 150.6242, 127.3600 | 150.6242 | 19.8 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=4;consumer=100% | 182.2373, 153.2996, 124.4023 | 153.2996 | 37.7 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=1% | 283.7898, 290.6342, 234.9841 | 283.7898 | 19.6 | 3 |
| yugabyte | single | family=reuse;policy=notmat;refs=8;consumer=100% | 265.0696, 295.1429, 269.7885 | 269.7885 | 11.1 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | contract | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-join-pg1-20261004 | postgres | single | bag@20 | join_on vs join_where | 0.8631 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 1.5172 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 1.1253 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 0.9384 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.0600 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.2087 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 88.3604 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9203 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20 | comma_join vs join_where | 1.1921 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | join_on vs join_where | 1.1593 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 1.1685 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 0.9664 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.8866 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.2739 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.8111 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 59.9526 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0440 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@200 | comma_join vs join_where | 0.8612 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | join_on vs join_where | 0.9337 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.2107 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.2444 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1068 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.1497 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9991 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 26.5944 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1010 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 0.9688 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | join_on vs join_where | 0.9108 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9320 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9631 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.0710 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9518 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 1.0082 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.8428 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9155 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 1.0160 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | join_on vs join_where | 0.9035 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9590 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 0.9784 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.0830 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9932 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.9734 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3819 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9516 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 0.9733 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | join_on vs join_where | 0.9694 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.9435 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.0281 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.0456 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9772 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.9639 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.3635 | >1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9096 | <1 |
| select-join-pg1-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 0.9751 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 4.7170 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.9745 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 1.1197 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0455 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 1.0774 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2926 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 1.1210 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 2.0559 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 1.0285 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.4316 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 0.9889 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.8597 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.8876 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.4782 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 1.1112 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.0541 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0715 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7447 | <1 |
| select-join-pg1-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 1.0218 | >1 |
| select-join-pg1-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.2925 | >1 |
| select-join-pg1-20261004 | postgres | single | position | position control first vs last | 0.8006 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | join_on vs join_where | 0.9326 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 1.4900 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 1.0547 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.5465 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.4209 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.0265 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 95.3136 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0173 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20 | comma_join vs join_where | 1.4151 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | join_on vs join_where | 1.0034 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 0.7664 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 0.6834 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.7001 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.0251 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 1.0349 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 49.3263 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.8567 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@200 | comma_join vs join_where | 0.7556 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | join_on vs join_where | 1.1740 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 1.1345 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.0552 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.0793 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9801 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9137 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 27.9888 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.1531 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 1.1873 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | join_on vs join_where | 0.9662 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9775 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9201 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.0204 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9697 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9573 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 7.0476 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.0129 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 0.9222 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | join_on vs join_where | 1.0071 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 1.0588 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 1.0592 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.1939 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9959 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 0.9693 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.4128 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 0.9765 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 0.9841 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | join_on vs join_where | 0.8595 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.8930 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 1.0433 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.0176 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9370 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.9921 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.3599 | >1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9862 | <1 |
| select-join-pg2-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 0.9693 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 3.4495 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.7377 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 1.2625 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.2230 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 1.1961 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2599 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 0.9647 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 2.0484 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 0.9528 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.3599 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 1.1357 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.9132 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.9424 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5094 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 0.9316 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 2.9804 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 1.0279 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.8087 | <1 |
| select-join-pg2-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 1.1512 | >1 |
| select-join-pg2-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.6102 | >1 |
| select-join-pg2-20261004 | postgres | single | position | position control first vs last | 0.8963 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | join_on vs join_where | 0.9146 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | filtered_derived vs join_where | 1.0503 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_filtered vs join_where | 1.0799 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_filtered_mat vs join_where | 1.2634 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_filtered_notmat vs join_where | 1.1509 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer vs join_where | 1.0727 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer_mat vs join_where | 91.5453 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.2832 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20 | comma_join vs join_where | 0.9857 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | join_on vs join_where | 0.8254 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | filtered_derived vs join_where | 1.2348 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_filtered vs join_where | 0.9293 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_filtered_mat vs join_where | 0.9515 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_filtered_notmat vs join_where | 1.0007 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer vs join_where | 0.8762 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer_mat vs join_where | 49.7734 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.3044 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@200 | comma_join vs join_where | 1.0011 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | join_on vs join_where | 0.9733 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | filtered_derived vs join_where | 0.9694 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_filtered vs join_where | 1.1086 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_filtered_mat vs join_where | 1.1342 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0962 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer vs join_where | 0.9565 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer_mat vs join_where | 26.2481 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9779 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@2000 | comma_join vs join_where | 1.1594 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | join_on vs join_where | 1.1950 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | filtered_derived vs join_where | 0.9944 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_filtered vs join_where | 0.9775 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_filtered_mat vs join_where | 1.0285 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9197 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer vs join_where | 0.9522 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer_mat vs join_where | 6.8246 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | cte_joined_outer_notmat vs join_where | 1.1340 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@20000 | comma_join vs join_where | 0.9616 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | join_on vs join_where | 0.9865 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | filtered_derived vs join_where | 0.9171 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_filtered vs join_where | 1.1113 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_filtered_mat vs join_where | 1.1615 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0709 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer vs join_where | 1.0872 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.3408 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.1439 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@100000 | comma_join vs join_where | 1.1075 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | join_on vs join_where | 0.9087 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | filtered_derived vs join_where | 0.9878 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_filtered vs join_where | 0.8741 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_filtered_mat vs join_where | 1.0372 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_filtered_notmat vs join_where | 0.7517 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer vs join_where | 0.9658 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer_mat vs join_where | 1.2244 | >1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.8534 | <1 |
| select-join-pg3-20261004 | postgres | single | bag@180000 | comma_join vs join_where | 1.0810 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=1 c=1% | mat vs default | 4.5509 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=1 c=1% | notmat vs default | 0.8863 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=1 c=100% | mat vs default | 1.1495 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=1 c=100% | notmat vs default | 1.0438 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=2 c=1% | mat vs default | 1.2711 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=2 c=1% | notmat vs default | 0.2558 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=2 c=100% | mat vs default | 1.0635 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=2 c=100% | notmat vs default | 1.8035 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=4 c=1% | mat vs default | 1.0466 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=4 c=1% | notmat vs default | 0.4479 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=4 c=100% | mat vs default | 1.3540 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=4 c=100% | notmat vs default | 2.6977 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=8 c=1% | mat vs default | 0.9600 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=8 c=1% | notmat vs default | 0.5418 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=8 c=100% | mat vs default | 0.8916 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=8 c=100% | notmat vs default | 3.5741 | >1 |
| select-join-pg3-20261004 | postgres | single | reuse n=16 c=1% | mat vs default | 0.9457 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=16 c=1% | notmat vs default | 0.7375 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=16 c=100% | mat vs default | 0.9414 | <1 |
| select-join-pg3-20261004 | postgres | single | reuse n=16 c=100% | notmat vs default | 4.0943 | >1 |
| select-join-pg3-20261004 | postgres | single | position | position control first vs last | 0.6056 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 0.9485 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9582 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.9864 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.9709 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.0307 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0307 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.6071 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 1.0355 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 1.0363 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 0.9750 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.0968 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 0.9782 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.9956 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 0.9951 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.0600 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 1.6531 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 0.9713 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 0.9709 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0810 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.0083 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.9797 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 1.0049 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9944 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.1179 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.6182 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0586 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 0.9575 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9130 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9410 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 0.9599 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 1.0115 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9537 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9510 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.2475 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9810 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9837 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 0.9632 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0418 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0036 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0830 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9676 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0024 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0603 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0038 | >1 |
| select-join-yb1-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 0.9996 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9246 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9138 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.8860 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 0.9462 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9012 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.8476 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.8621 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.8986 | <1 |
| select-join-yb1-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9222 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9109 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 1.0108 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 1.0000 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0185 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 1.0667 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.8214 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9799 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.7319 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9533 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.9047 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.1156 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 3.4428 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0221 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 5.1225 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.9010 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 4.6380 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9912 | <1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 8.9376 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0753 | >1 |
| select-join-yb1-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 11.1539 | >1 |
| select-join-yb1-20261004 | yugabyte | single | position | position control first vs last | 0.8609 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 1.0219 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 0.9873 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.9716 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.9685 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 0.9863 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0039 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.5131 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9616 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 0.9868 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 0.9947 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 0.9835 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.0408 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 0.9817 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.0115 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.0431 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 1.5986 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0295 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 1.0515 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 0.9816 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 0.9897 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 0.9921 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9977 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 1.0166 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 0.9772 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5946 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 0.9914 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 1.0360 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9544 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9613 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 0.9515 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 0.9639 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9541 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 0.9218 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.2723 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9578 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 0.9460 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 1.1138 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 1.0557 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 1.0439 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.1268 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 1.0551 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 1.0750 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.1284 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0391 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 1.0527 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 1.0471 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9873 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9543 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0950 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9893 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 1.0249 | >1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9064 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9764 | <1 |
| select-join-yb2-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9618 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9204 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9601 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9705 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 0.9661 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.9898 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.7804 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 0.9789 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.6611 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9199 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.5014 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0118 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.6953 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0318 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 5.1567 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 1.0512 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 5.3370 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 1.0179 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.4710 | >1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 0.9790 | <1 |
| select-join-yb2-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.9072 | >1 |
| select-join-yb2-20261004 | yugabyte | single | position | position control first vs last | 1.3470 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | join_on vs join_where | 1.1054 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | filtered_derived vs join_where | 1.0803 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered vs join_where | 0.9735 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered_mat vs join_where | 0.9362 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_filtered_notmat vs join_where | 1.0039 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer vs join_where | 1.0005 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_mat vs join_where | 1.5277 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | cte_joined_outer_notmat vs join_where | 0.9718 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20 | comma_join vs join_where | 1.0320 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | join_on vs join_where | 1.0577 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | filtered_derived vs join_where | 1.0189 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered vs join_where | 1.0060 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered_mat vs join_where | 1.0024 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_filtered_notmat vs join_where | 1.0188 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer vs join_where | 1.0335 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer_mat vs join_where | 1.5828 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | cte_joined_outer_notmat vs join_where | 1.0583 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@200 | comma_join vs join_where | 1.0621 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | join_on vs join_where | 1.0260 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | filtered_derived vs join_where | 1.0304 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered vs join_where | 1.0183 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered_mat vs join_where | 0.9761 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_filtered_notmat vs join_where | 0.9903 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer vs join_where | 1.0381 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_mat vs join_where | 1.5342 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | cte_joined_outer_notmat vs join_where | 1.0393 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@2000 | comma_join vs join_where | 1.0150 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | join_on vs join_where | 0.9692 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | filtered_derived vs join_where | 0.9891 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered vs join_where | 0.9704 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered_mat vs join_where | 0.9993 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_filtered_notmat vs join_where | 0.9908 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer vs join_where | 1.0006 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_mat vs join_where | 1.3092 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | cte_joined_outer_notmat vs join_where | 0.9953 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@20000 | comma_join vs join_where | 1.0018 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | join_on vs join_where | 0.9920 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | filtered_derived vs join_where | 0.9914 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered vs join_where | 0.9906 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered_mat vs join_where | 1.0547 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_filtered_notmat vs join_where | 0.9993 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer vs join_where | 0.9765 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_mat vs join_where | 1.0590 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | cte_joined_outer_notmat vs join_where | 1.0022 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@100000 | comma_join vs join_where | 0.9956 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | join_on vs join_where | 0.9691 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | filtered_derived vs join_where | 0.9689 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered vs join_where | 0.9708 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered_mat vs join_where | 1.0270 | >1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_filtered_notmat vs join_where | 0.9844 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer vs join_where | 0.9818 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_mat vs join_where | 0.9270 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | cte_joined_outer_notmat vs join_where | 0.9953 | <1 |
| select-join-yb3-20261004 | yugabyte | single | bag@180000 | comma_join vs join_where | 0.9903 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=1 c=1% | mat vs default | 0.9168 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=1 c=1% | notmat vs default | 0.9525 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=1 c=100% | mat vs default | 0.9743 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=1 c=100% | notmat vs default | 1.0153 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=2 c=1% | mat vs default | 0.9770 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=2 c=1% | notmat vs default | 1.6819 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=2 c=100% | mat vs default | 1.0212 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=2 c=100% | notmat vs default | 1.8079 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=4 c=1% | mat vs default | 0.9858 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=4 c=1% | notmat vs default | 2.4592 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=4 c=100% | mat vs default | 1.0000 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=4 c=100% | notmat vs default | 2.3941 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=8 c=1% | mat vs default | 1.0188 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=8 c=1% | notmat vs default | 4.7235 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=8 c=100% | mat vs default | 0.9762 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=8 c=100% | notmat vs default | 5.2920 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=16 c=1% | mat vs default | 0.9265 | <1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=16 c=1% | notmat vs default | 9.2066 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=16 c=100% | mat vs default | 1.0063 | >1 |
| select-join-yb3-20261004 | yugabyte | single | reuse n=16 c=100% | notmat vs default | 9.6038 | >1 |
| select-join-yb3-20261004 | yugabyte | single | position | position control first vs last | 1.0102 | >1 |

## Recorded limits

- [select-join-pg1-20261004] one host, single client: no independent hardware, real-network or equal-total-budget comparison
- [select-join-pg1-20261004] no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- [select-join-pg1-20261004] a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
