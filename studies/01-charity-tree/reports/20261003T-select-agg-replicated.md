# Study 01 SELECT aggregate & top-N — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 6 whole runs summarised, each run's own median shown; never merged across engines.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-agg-pg-1b-20261003 | postgres | single | sweep | 42 |
| select-agg-pg-2-20261003 | postgres | single | sweep | 42 |
| select-agg-pg-3-20261003 | postgres | single | sweep | 42 |
| select-agg-yb-1-20261003 | yugabyte | single | sweep | 42 |
| select-agg-yb-2-20261003 | yugabyte | single | sweep | 42 |
| select-agg-yb-3-20261003 | yugabyte | single | sweep | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | family=aggregate;stmt=correlated;parent_max=1 | 0.163634, 0.240582, 0.267138 | 0.240582 | 43.0 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=10 | 0.191366, 0.253065, 0.262873 | 0.253065 | 28.3 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=100 | 0.422182, 0.502608, 0.513632 | 0.502608 | 18.2 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=1000 | 2.203104, 2.315719, 2.241320 | 2.241320 | 5.0 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=10000 | 20.178764, 22.264445, 20.454044 | 20.454044 | 10.2 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=1 | 0.234977, 0.291771, 0.273383 | 0.273383 | 20.8 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=10 | 0.224404, 0.320708, 0.273469 | 0.273469 | 35.2 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=100 | 0.415305, 0.460999, 0.418194 | 0.418194 | 10.9 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=1000 | 1.726389, 1.903942, 1.659597 | 1.726389 | 14.2 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=10000 | 7.020799, 7.532297, 7.710870 | 7.532297 | 9.2 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=1;position=last | 0.349337 | 0.349337 | 0.0 | 1 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=1 | 2.798993, 2.765842, 2.805736 | 2.798993 | 1.4 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=10 | 2.760325, 2.876002, 2.669221 | 2.760325 | 7.5 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=100 | 2.942386, 2.930769, 2.896816 | 2.930769 | 1.6 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=1000 | 3.136784, 3.216808, 3.235251 | 3.216808 | 3.1 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=10000 | 5.763481, 6.861774, 6.210246 | 6.210246 | 17.7 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=1 | 0.124118, 0.216995, 0.240613 | 0.216995 | 53.7 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=10 | 0.169402, 0.226328, 0.223350 | 0.223350 | 25.5 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=100 | 0.259799, 0.277735, 0.248373 | 0.259799 | 11.3 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=1000 | 0.649699, 0.721542, 0.713099 | 0.713099 | 10.1 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=10000 | 5.428776, 5.040590, 5.167448 | 5.167448 | 7.5 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=1 | 0.195276, 0.234856, 0.171264 | 0.195276 | 32.6 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=10 | 0.197841, 0.215903, 0.224811 | 0.215903 | 12.5 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=100 | 0.421836, 0.437996, 0.427217 | 0.427217 | 3.8 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=1000 | 1.693071, 1.816070, 1.637665 | 1.693071 | 10.5 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=10000 | 16.774430, 16.629260, 15.987555 | 16.629260 | 4.7 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=1 | 0.191538, 0.223831, 0.194375 | 0.194375 | 16.6 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=10 | 0.208794, 0.205327, 0.227256 | 0.208794 | 10.5 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=100 | 0.320492, 0.344501, 0.319639 | 0.320492 | 7.8 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=1000 | 0.994342, 1.017237, 0.975977 | 0.994342 | 4.1 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=10000 | 7.745943, 7.999212, 8.049526 | 7.999212 | 3.8 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=1 | 3.473802, 3.993887, 4.020777 | 3.993887 | 13.7 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=10 | 8.241740, 8.229113, 9.763410 | 8.241740 | 18.6 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=100 | 52.509607, 56.826322, 59.848449 | 56.826322 | 12.9 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=1000 | 521.916135, 519.449263, 492.632448 | 519.449263 | 5.6 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=10000 | 8460.299924, 8716.346747, 8426.879726 | 8460.299924 | 3.4 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=1 | 3.741284, 3.866297, 3.998084 | 3.866297 | 6.6 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=10 | 3.827498, 3.952024, 4.408331 | 3.952024 | 14.7 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=100 | 4.860858, 5.474165, 5.568177 | 5.474165 | 12.9 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=1000 | 9.964763, 9.926294, 10.700134 | 9.964763 | 7.8 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=10000 | 17.272811, 17.183780, 18.906516 | 17.272811 | 10.0 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=1;position=last | 3.677096 | 3.677096 | 0.0 | 1 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=1 | 11.118851, 12.130249, 12.107479 | 12.107479 | 8.4 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=10 | 11.598654, 12.482785, 14.378997 | 12.482785 | 22.3 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=100 | 10.810412, 12.297712, 12.489613 | 12.297712 | 13.7 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=1000 | 11.866113, 12.648519, 12.935225 | 12.648519 | 8.5 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=10000 | 21.890487, 22.390241, 21.406016 | 21.890487 | 4.5 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=1 | 5.739392, 5.690649, 5.915677 | 5.739392 | 3.9 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=10 | 5.472999, 6.002842, 6.106900 | 6.002842 | 10.6 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=100 | 5.282858, 5.408136, 5.638299 | 5.408136 | 6.6 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=1000 | 6.029429, 6.296445, 6.679306 | 6.296445 | 10.3 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=10000 | 16.527319, 17.214960, 21.405415 | 17.214960 | 28.3 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=1 | 3.550300, 3.538668, 3.664861 | 3.550300 | 3.6 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=10 | 6.509639, 6.795568, 7.337406 | 6.795568 | 12.2 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=100 | 39.966097, 38.662160, 39.160220 | 39.160220 | 3.3 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=1000 | 302.125610, 335.158272, 361.462456 | 335.158272 | 17.7 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=10000 | 5497.831655, 5679.717299, 5862.099402 | 5679.717299 | 6.4 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=1 | 6.019929, 6.229836, 5.842895 | 6.019929 | 6.4 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=10 | 5.458151, 5.740875, 6.331968 | 5.740875 | 15.2 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=100 | 5.657493, 6.108811, 5.694723 | 5.694723 | 7.9 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=1000 | 6.741464, 7.351668, 7.239546 | 7.239546 | 8.4 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=10000 | 20.192528, 20.694888, 21.080825 | 20.694888 | 4.3 | 3 |

## Pair ratios per run (per engine)

| engine | topology | family | level | pair | ratios (B/A) | signs |
|---|---|---|---|---|---|---|
| postgres | single | aggregate | 1 | pre_agg vs post_join | 11.9118 | >1 |
| postgres | single | aggregate | 1 | correlated vs post_join | 0.6964 | <1 |
| postgres | single | aggregate | 1 | correlated vs pre_agg | 0.0585 | <1 |
| postgres | single | aggregate | 10 | pre_agg vs post_join | 12.3007 | >1 |
| postgres | single | aggregate | 10 | correlated vs post_join | 0.8528 | <1 |
| postgres | single | aggregate | 10 | correlated vs pre_agg | 0.0693 | <1 |
| postgres | single | aggregate | 100 | pre_agg vs post_join | 7.0849 | >1 |
| postgres | single | aggregate | 100 | correlated vs post_join | 1.0166 | >1 |
| postgres | single | aggregate | 100 | correlated vs pre_agg | 0.1435 | <1 |
| postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.8170 | >1 |
| postgres | single | aggregate | 1000 | correlated vs post_join | 1.2761 | >1 |
| postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.7023 | <1 |
| postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.8209 | <1 |
| postgres | single | aggregate | 10000 | correlated vs post_join | 2.8741 | >1 |
| postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.5011 | >1 |
| postgres | single | topn | 1 | lateral vs rownumber | 1.0195 | >1 |
| postgres | single | topn | 1 | distinct_on vs rownumber | 0.6480 | <1 |
| postgres | single | topn | 1 | distinct_on vs lateral | 0.6356 | <1 |
| postgres | single | topn | 10 | lateral vs rownumber | 0.9475 | <1 |
| postgres | single | topn | 10 | distinct_on vs rownumber | 0.8113 | <1 |
| postgres | single | topn | 10 | distinct_on vs lateral | 0.8563 | <1 |
| postgres | single | topn | 100 | lateral vs rownumber | 1.3162 | >1 |
| postgres | single | topn | 100 | distinct_on vs rownumber | 0.8106 | <1 |
| postgres | single | topn | 100 | distinct_on vs lateral | 0.6159 | <1 |
| postgres | single | topn | 1000 | lateral vs rownumber | 1.7027 | >1 |
| postgres | single | topn | 1000 | distinct_on vs rownumber | 0.6534 | <1 |
| postgres | single | topn | 1000 | distinct_on vs lateral | 0.3837 | <1 |
| postgres | single | topn | 10000 | lateral vs rownumber | 2.1656 | >1 |
| postgres | single | topn | 10000 | distinct_on vs rownumber | 0.7009 | <1 |
| postgres | single | topn | 10000 | distinct_on vs lateral | 0.3236 | <1 |
| postgres | single | aggregate | 1 | pre_agg vs post_join | 9.4795 | >1 |
| postgres | single | aggregate | 1 | correlated vs post_join | 0.8246 | <1 |
| postgres | single | aggregate | 1 | correlated vs pre_agg | 0.0870 | <1 |
| postgres | single | aggregate | 10 | pre_agg vs post_join | 8.9677 | >1 |
| postgres | single | aggregate | 10 | correlated vs post_join | 0.7891 | <1 |
| postgres | single | aggregate | 10 | correlated vs pre_agg | 0.0880 | <1 |
| postgres | single | aggregate | 100 | pre_agg vs post_join | 6.3574 | >1 |
| postgres | single | aggregate | 100 | correlated vs post_join | 1.0903 | >1 |
| postgres | single | aggregate | 100 | correlated vs pre_agg | 0.1715 | <1 |
| postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.6896 | >1 |
| postgres | single | aggregate | 1000 | correlated vs post_join | 1.2163 | >1 |
| postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.7199 | <1 |
| postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.9110 | <1 |
| postgres | single | aggregate | 10000 | correlated vs post_join | 2.9559 | >1 |
| postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.2447 | >1 |
| postgres | single | topn | 1 | lateral vs rownumber | 1.0493 | >1 |
| postgres | single | topn | 1 | distinct_on vs rownumber | 0.9695 | <1 |
| postgres | single | topn | 1 | distinct_on vs lateral | 0.9239 | <1 |
| postgres | single | topn | 10 | lateral vs rownumber | 1.0515 | >1 |
| postgres | single | topn | 10 | distinct_on vs rownumber | 1.1023 | >1 |
| postgres | single | topn | 10 | distinct_on vs lateral | 1.0483 | >1 |
| postgres | single | topn | 100 | lateral vs rownumber | 1.2714 | >1 |
| postgres | single | topn | 100 | distinct_on vs rownumber | 0.8062 | <1 |
| postgres | single | topn | 100 | distinct_on vs lateral | 0.6341 | <1 |
| postgres | single | topn | 1000 | lateral vs rownumber | 1.7853 | >1 |
| postgres | single | topn | 1000 | distinct_on vs rownumber | 0.7093 | <1 |
| postgres | single | topn | 1000 | distinct_on vs lateral | 0.3973 | <1 |
| postgres | single | topn | 10000 | lateral vs rownumber | 2.0789 | >1 |
| postgres | single | topn | 10000 | distinct_on vs rownumber | 0.6301 | <1 |
| postgres | single | topn | 10000 | distinct_on vs lateral | 0.3031 | <1 |
| postgres | single | aggregate | 1 | pre_agg vs post_join | 10.2630 | >1 |
| postgres | single | aggregate | 1 | correlated vs post_join | 0.9772 | <1 |
| postgres | single | aggregate | 1 | correlated vs pre_agg | 0.0952 | <1 |
| postgres | single | aggregate | 10 | pre_agg vs post_join | 9.7606 | >1 |
| postgres | single | aggregate | 10 | correlated vs post_join | 0.9613 | <1 |
| postgres | single | aggregate | 10 | correlated vs pre_agg | 0.0985 | <1 |
| postgres | single | aggregate | 100 | pre_agg vs post_join | 6.9270 | >1 |
| postgres | single | aggregate | 100 | correlated vs post_join | 1.2282 | >1 |
| postgres | single | aggregate | 100 | correlated vs pre_agg | 0.1773 | <1 |
| postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.9494 | >1 |
| postgres | single | aggregate | 1000 | correlated vs post_join | 1.3505 | >1 |
| postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.6928 | <1 |
| postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.8054 | <1 |
| postgres | single | aggregate | 10000 | correlated vs post_join | 2.6526 | >1 |
| postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.2936 | >1 |
| postgres | single | topn | 1 | lateral vs rownumber | 0.8811 | <1 |
| postgres | single | topn | 1 | distinct_on vs rownumber | 1.2379 | >1 |
| postgres | single | topn | 1 | distinct_on vs lateral | 1.4049 | >1 |
| postgres | single | topn | 10 | lateral vs rownumber | 0.9892 | <1 |
| postgres | single | topn | 10 | distinct_on vs rownumber | 0.9828 | <1 |
| postgres | single | topn | 10 | distinct_on vs lateral | 0.9935 | <1 |
| postgres | single | topn | 100 | lateral vs rownumber | 1.3366 | >1 |
| postgres | single | topn | 100 | distinct_on vs rownumber | 0.7770 | <1 |
| postgres | single | topn | 100 | distinct_on vs lateral | 0.5814 | <1 |
| postgres | single | topn | 1000 | lateral vs rownumber | 1.6780 | >1 |
| postgres | single | topn | 1000 | distinct_on vs rownumber | 0.7307 | <1 |
| postgres | single | topn | 1000 | distinct_on vs lateral | 0.4354 | <1 |
| postgres | single | topn | 10000 | lateral vs rownumber | 1.9861 | >1 |
| postgres | single | topn | 10000 | distinct_on vs rownumber | 0.6420 | <1 |
| postgres | single | topn | 10000 | distinct_on vs lateral | 0.3232 | <1 |
| postgres | single | aggregate | 1 | position control post_join first vs last | 1.2778 | >1 |
| yugabyte | single | aggregate | 1 | pre_agg vs post_join | 2.9719 | >1 |
| yugabyte | single | aggregate | 1 | correlated vs post_join | 0.9285 | <1 |
| yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.3124 | <1 |
| yugabyte | single | aggregate | 10 | pre_agg vs post_join | 3.0303 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs post_join | 2.1533 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.7106 | <1 |
| yugabyte | single | aggregate | 100 | pre_agg vs post_join | 2.2240 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs post_join | 10.8025 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs pre_agg | 4.8573 | >1 |
| yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 1.1908 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs post_join | 52.3762 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 43.9837 | >1 |
| yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.2673 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs post_join | 489.8045 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 386.4829 | >1 |
| yugabyte | single | topn | 1 | lateral vs rownumber | 0.5898 | <1 |
| yugabyte | single | topn | 1 | distinct_on vs rownumber | 0.9534 | <1 |
| yugabyte | single | topn | 1 | distinct_on vs lateral | 1.6166 | >1 |
| yugabyte | single | topn | 10 | lateral vs rownumber | 1.1926 | >1 |
| yugabyte | single | topn | 10 | distinct_on vs rownumber | 1.0027 | >1 |
| yugabyte | single | topn | 10 | distinct_on vs lateral | 0.8408 | <1 |
| yugabyte | single | topn | 100 | lateral vs rownumber | 7.0643 | >1 |
| yugabyte | single | topn | 100 | distinct_on vs rownumber | 0.9338 | <1 |
| yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1322 | <1 |
| yugabyte | single | topn | 1000 | lateral vs rownumber | 44.8160 | >1 |
| yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.8944 | <1 |
| yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0200 | <1 |
| yugabyte | single | topn | 10000 | lateral vs rownumber | 272.2706 | >1 |
| yugabyte | single | topn | 10000 | distinct_on vs rownumber | 0.8185 | <1 |
| yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0030 | <1 |
| yugabyte | single | aggregate | 1 | pre_agg vs post_join | 3.1374 | >1 |
| yugabyte | single | aggregate | 1 | correlated vs post_join | 1.0330 | >1 |
| yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.3293 | <1 |
| yugabyte | single | aggregate | 10 | pre_agg vs post_join | 3.1586 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs post_join | 2.0823 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.6592 | <1 |
| yugabyte | single | aggregate | 100 | pre_agg vs post_join | 2.2465 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs post_join | 10.3808 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs pre_agg | 4.6209 | >1 |
| yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 1.2742 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs post_join | 52.3306 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 41.0680 | >1 |
| yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.3030 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs post_join | 507.2427 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 389.2922 | >1 |
| yugabyte | single | topn | 1 | lateral vs rownumber | 0.5680 | <1 |
| yugabyte | single | topn | 1 | distinct_on vs rownumber | 0.9135 | <1 |
| yugabyte | single | topn | 1 | distinct_on vs lateral | 1.6081 | >1 |
| yugabyte | single | topn | 10 | lateral vs rownumber | 1.1837 | >1 |
| yugabyte | single | topn | 10 | distinct_on vs rownumber | 1.0456 | >1 |
| yugabyte | single | topn | 10 | distinct_on vs lateral | 0.8833 | <1 |
| yugabyte | single | topn | 100 | lateral vs rownumber | 6.3289 | >1 |
| yugabyte | single | topn | 100 | distinct_on vs rownumber | 0.8853 | <1 |
| yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1399 | <1 |
| yugabyte | single | topn | 1000 | lateral vs rownumber | 45.5894 | >1 |
| yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.8565 | <1 |
| yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0188 | <1 |
| yugabyte | single | topn | 10000 | lateral vs rownumber | 274.4503 | >1 |
| yugabyte | single | topn | 10000 | distinct_on vs rownumber | 0.8318 | <1 |
| yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0030 | <1 |
| yugabyte | single | aggregate | 1 | position control post_join first vs last | 0.9511 | <1 |
| yugabyte | single | aggregate | 1 | pre_agg vs post_join | 3.0283 | >1 |
| yugabyte | single | aggregate | 1 | correlated vs post_join | 1.0057 | >1 |
| yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.3321 | <1 |
| yugabyte | single | aggregate | 10 | pre_agg vs post_join | 3.2618 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs post_join | 2.2148 | >1 |
| yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.6790 | <1 |
| yugabyte | single | aggregate | 100 | pre_agg vs post_join | 2.2430 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs post_join | 10.7483 | >1 |
| yugabyte | single | aggregate | 100 | correlated vs pre_agg | 4.7919 | >1 |
| yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 1.2089 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs post_join | 46.0398 | >1 |
| yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 38.0846 | >1 |
| yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.1322 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs post_join | 445.7130 | >1 |
| yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 393.6688 | >1 |
| yugabyte | single | topn | 1 | lateral vs rownumber | 0.6272 | <1 |
| yugabyte | single | topn | 1 | distinct_on vs rownumber | 1.0125 | >1 |
| yugabyte | single | topn | 1 | distinct_on vs lateral | 1.6142 | >1 |
| yugabyte | single | topn | 10 | lateral vs rownumber | 1.1588 | >1 |
| yugabyte | single | topn | 10 | distinct_on vs rownumber | 0.9645 | <1 |
| yugabyte | single | topn | 10 | distinct_on vs lateral | 0.8323 | <1 |
| yugabyte | single | topn | 100 | lateral vs rownumber | 6.8766 | >1 |
| yugabyte | single | topn | 100 | distinct_on vs rownumber | 0.9901 | <1 |
| yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1440 | <1 |
| yugabyte | single | topn | 1000 | lateral vs rownumber | 49.9289 | >1 |
| yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.9226 | <1 |
| yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0185 | <1 |
| yugabyte | single | topn | 10000 | lateral vs rownumber | 278.0773 | >1 |
| yugabyte | single | topn | 10000 | distinct_on vs rownumber | 1.0154 | >1 |
| yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0037 | <1 |

## Recorded limits

- [select-agg-pg-1b-20261003] one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- [select-agg-pg-1b-20261003] single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- [select-agg-pg-1b-20261003] YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
