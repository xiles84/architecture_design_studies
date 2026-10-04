# Study 01 SELECT aggregate & top-N — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 6 whole runs summarised, each run's own median shown; never merged across engines.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-order-aggregate-postgres-1-20261004 | postgres | single | sweep | 42 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | sweep | 42 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | sweep | 42 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | sweep | 42 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | sweep | 42 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | sweep | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | family=aggregate;stmt=correlated;parent_max=1 | 0.254698, 0.202625, 0.261305 | 0.254698 | 23.0 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=10 | 0.254588, 0.219296, 0.493301 | 0.254588 | 107.6 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=100 | 0.700587, 0.460843, 0.625201 | 0.625201 | 38.3 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=1000 | 1.727673, 1.732768, 1.681847 | 1.727673 | 2.9 | 3 |
| postgres | single | family=aggregate;stmt=correlated;parent_max=10000 | 15.907986, 16.241235, 16.652408 | 16.241235 | 4.6 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=1 | 0.283191, 0.434762, 0.289159 | 0.289159 | 52.4 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=10 | 0.499808, 0.512776, 0.272098 | 0.499808 | 48.2 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=100 | 0.406413, 0.365654, 0.370106 | 0.370106 | 11.0 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=1000 | 1.365769, 1.607923, 1.465880 | 1.465880 | 16.5 | 3 |
| postgres | single | family=aggregate;stmt=post_join;parent_max=10000 | 5.995085, 6.515128, 6.161754 | 6.161754 | 8.4 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=1 | 2.230215, 2.483171, 2.210538 | 2.230215 | 12.2 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=10 | 2.295233, 2.220210, 2.421641 | 2.295233 | 8.8 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=100 | 2.325254, 2.388580, 2.309264 | 2.325254 | 3.4 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=1000 | 2.587679, 2.510928, 2.576308 | 2.576308 | 3.0 | 3 |
| postgres | single | family=aggregate;stmt=pre_agg;parent_max=10000 | 4.751486, 4.995993, 4.910587 | 4.910587 | 5.0 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=1 | 0.203181, 0.195156, 0.181084 | 0.195156 | 11.3 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=10 | 0.381249, 0.247913, 0.150273 | 0.247913 | 93.2 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=100 | 0.423064, 0.238245, 0.247900 | 0.247900 | 74.6 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=1000 | 0.737363, 0.567822, 0.518560 | 0.567822 | 38.5 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=10000 | 3.806609, 3.606980, 3.660584 | 3.660584 | 5.5 | 3 |
| postgres | single | family=topn;stmt=distinct_on;parent_max=1000;position=last | 0.523844 | 0.523844 | 0.0 | 1 |
| postgres | single | family=topn;stmt=lateral;parent_max=1 | 0.228814, 0.250637, 0.231808 | 0.231808 | 9.4 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=10 | 0.260120, 0.404105, 0.258035 | 0.260120 | 56.2 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=100 | 0.608117, 0.616913, 0.351821 | 0.608117 | 43.6 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=1000 | 1.365204, 1.385665, 1.409829 | 1.385665 | 3.2 | 3 |
| postgres | single | family=topn;stmt=lateral;parent_max=10000 | 12.350192, 11.945957, 12.706788 | 12.350192 | 6.2 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=1 | 0.240248, 0.256439, 0.413695 | 0.256439 | 67.6 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=10 | 0.179630, 0.266148, 0.499921 | 0.266148 | 120.3 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=100 | 0.350775, 0.317123, 0.295158 | 0.317123 | 17.5 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=1000 | 0.803642, 0.868835, 0.846486 | 0.846486 | 7.7 | 3 |
| postgres | single | family=topn;stmt=rownumber;parent_max=10000 | 6.238146, 5.747159, 5.765574 | 5.765574 | 8.5 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=1 | 4.186747, 4.085893, 8.266088 | 4.186747 | 99.8 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=10 | 7.060701, 7.243868, 17.758949 | 7.243868 | 147.7 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=100 | 46.598187, 43.470839, 124.504364 | 46.598187 | 173.9 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=1000 | 425.472067, 416.727342, 1170.377956 | 425.472067 | 177.1 | 3 |
| yugabyte | single | family=aggregate;stmt=correlated;parent_max=10000 | 8695.735167, 6726.910867, 19981.692225 | 8695.735167 | 152.4 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=1 | 3.687866, 3.316064, 7.739420 | 3.687866 | 119.9 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=10 | 3.668820, 3.684170, 8.885454 | 3.684170 | 141.6 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=100 | 5.392189, 5.063425, 10.484977 | 5.392189 | 100.5 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=1000 | 10.254603, 9.475192, 25.463011 | 10.254603 | 155.9 | 3 |
| yugabyte | single | family=aggregate;stmt=post_join;parent_max=10000 | 16.081415, 16.910202, 37.412629 | 16.910202 | 126.1 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=1 | 11.858074, 9.607655, 24.434710 | 11.858074 | 125.0 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=10 | 10.035164, 11.998627, 24.411916 | 11.998627 | 119.8 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=100 | 13.228366, 9.343454, 28.382864 | 13.228366 | 143.9 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=1000 | 12.362386, 13.310056, 25.265151 | 13.310056 | 96.9 | 3 |
| yugabyte | single | family=aggregate;stmt=pre_agg;parent_max=10000 | 17.622484, 18.925754, 46.054655 | 18.925754 | 150.2 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=1 | 5.444005, 5.008627, 12.146865 | 5.444005 | 131.1 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=10 | 5.006473, 5.483670, 11.742605 | 5.483670 | 122.8 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=100 | 4.714137, 5.675520, 11.744019 | 5.675520 | 123.9 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=1000 | 5.818641, 5.835957, 14.592044 | 5.835957 | 150.3 | 3 |
| yugabyte | single | family=topn;stmt=distinct_on;parent_max=10000 | 15.099795, 18.050407, 37.909284 | 18.050407 | 126.4 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=1 | 3.685713, 4.072752, 7.011562 | 4.072752 | 81.7 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=10 | 5.691520, 5.742642, 14.570725 | 5.742642 | 154.6 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=100 | 32.039796, 32.109549, 80.777657 | 32.109549 | 151.8 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=1000 | 286.601844, 276.741397, 799.642915 | 286.601844 | 182.4 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=10000 | 5963.188427, 4660.097965, 13886.221819 | 5963.188427 | 154.7 | 3 |
| yugabyte | single | family=topn;stmt=lateral;parent_max=10000;position=last | 10395.212656 | 10395.212656 | 0.0 | 1 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=1 | 5.295108, 5.823031, 13.013492 | 5.823031 | 132.5 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=10 | 5.068200, 4.954549, 11.701965 | 5.068200 | 133.1 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=100 | 4.859502, 4.560806, 13.231090 | 4.859502 | 178.4 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=1000 | 6.320737, 6.282085, 16.313609 | 6.320737 | 158.7 | 3 |
| yugabyte | single | family=topn;stmt=rownumber;parent_max=10000 | 21.648975, 18.108292, 48.069590 | 21.648975 | 138.4 | 3 |

## Pair ratios per run (per engine)

| run | engine | topology | family | level | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|---|
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1 | pre_agg vs post_join | 7.8753 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1 | correlated vs post_join | 0.8994 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1 | correlated vs pre_agg | 0.1142 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10 | pre_agg vs post_join | 4.5922 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10 | correlated vs post_join | 0.5094 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10 | correlated vs pre_agg | 0.1109 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 100 | pre_agg vs post_join | 5.7214 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 100 | correlated vs post_join | 1.7238 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 100 | correlated vs pre_agg | 0.3013 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.8947 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1000 | correlated vs post_join | 1.2650 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.6677 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.7926 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10000 | correlated vs post_join | 2.6535 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.3480 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1 | lateral vs rownumber | 0.9524 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1 | distinct_on vs rownumber | 0.8457 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1 | distinct_on vs lateral | 0.8880 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10 | lateral vs rownumber | 1.4481 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10 | distinct_on vs rownumber | 2.1224 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10 | distinct_on vs lateral | 1.4657 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 100 | lateral vs rownumber | 1.7336 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 100 | distinct_on vs rownumber | 1.2061 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 100 | distinct_on vs lateral | 0.6957 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1000 | lateral vs rownumber | 1.6988 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1000 | distinct_on vs rownumber | 0.9175 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 1000 | distinct_on vs lateral | 0.5401 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10000 | lateral vs rownumber | 1.9798 | >1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10000 | distinct_on vs rownumber | 0.6102 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | topn | 10000 | distinct_on vs lateral | 0.3082 | <1 |
| select-order-aggregate-postgres-1-20261004 | postgres | single | position | 0 | position control first vs last | 0.7104 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1 | pre_agg vs post_join | 5.7116 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1 | correlated vs post_join | 0.4661 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1 | correlated vs pre_agg | 0.0816 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10 | pre_agg vs post_join | 4.3298 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10 | correlated vs post_join | 0.4277 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10 | correlated vs pre_agg | 0.0988 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 100 | pre_agg vs post_join | 6.5324 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 100 | correlated vs post_join | 1.2603 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 100 | correlated vs pre_agg | 0.1929 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.5616 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1000 | correlated vs post_join | 1.0776 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.6901 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.7668 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10000 | correlated vs post_join | 2.4928 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.2509 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1 | lateral vs rownumber | 0.9774 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1 | distinct_on vs rownumber | 0.7610 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1 | distinct_on vs lateral | 0.7786 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10 | lateral vs rownumber | 1.5183 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10 | distinct_on vs rownumber | 0.9315 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10 | distinct_on vs lateral | 0.6135 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 100 | lateral vs rownumber | 1.9453 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 100 | distinct_on vs rownumber | 0.7513 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 100 | distinct_on vs lateral | 0.3862 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1000 | lateral vs rownumber | 1.5949 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1000 | distinct_on vs rownumber | 0.6535 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 1000 | distinct_on vs lateral | 0.4098 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10000 | lateral vs rownumber | 2.0786 | >1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10000 | distinct_on vs rownumber | 0.6276 | <1 |
| select-order-aggregate-postgres-2-20261004 | postgres | single | topn | 10000 | distinct_on vs lateral | 0.3019 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1 | pre_agg vs post_join | 7.6447 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1 | correlated vs post_join | 0.9037 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1 | correlated vs pre_agg | 0.1182 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10 | pre_agg vs post_join | 8.8999 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10 | correlated vs post_join | 1.8130 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10 | correlated vs pre_agg | 0.2037 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 100 | pre_agg vs post_join | 6.2395 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 100 | correlated vs post_join | 1.6892 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 100 | correlated vs pre_agg | 0.2707 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1000 | pre_agg vs post_join | 1.7575 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1000 | correlated vs post_join | 1.1473 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 1000 | correlated vs pre_agg | 0.6528 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10000 | pre_agg vs post_join | 0.7969 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10000 | correlated vs post_join | 2.7025 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | aggregate | 10000 | correlated vs pre_agg | 3.3911 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1 | lateral vs rownumber | 0.5603 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1 | distinct_on vs rownumber | 0.4377 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1 | distinct_on vs lateral | 0.7812 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10 | lateral vs rownumber | 0.5162 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10 | distinct_on vs rownumber | 0.3006 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10 | distinct_on vs lateral | 0.5824 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 100 | lateral vs rownumber | 1.1920 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 100 | distinct_on vs rownumber | 0.8399 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 100 | distinct_on vs lateral | 0.7046 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1000 | lateral vs rownumber | 1.6655 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1000 | distinct_on vs rownumber | 0.6126 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 1000 | distinct_on vs lateral | 0.3678 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10000 | lateral vs rownumber | 2.2039 | >1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10000 | distinct_on vs rownumber | 0.6349 | <1 |
| select-order-aggregate-postgres-3-20261004 | postgres | single | topn | 10000 | distinct_on vs lateral | 0.2881 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1 | pre_agg vs post_join | 3.2154 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1 | correlated vs post_join | 1.1353 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.3531 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10 | pre_agg vs post_join | 2.7353 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10 | correlated vs post_join | 1.9245 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.7036 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 100 | pre_agg vs post_join | 2.4532 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 100 | correlated vs post_join | 8.6418 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 100 | correlated vs pre_agg | 3.5226 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 1.2055 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1000 | correlated vs post_join | 41.4908 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 34.4167 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.0958 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10000 | correlated vs post_join | 540.7320 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 493.4455 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1 | lateral vs rownumber | 0.6961 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1 | distinct_on vs rownumber | 1.0281 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1 | distinct_on vs lateral | 1.4771 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10 | lateral vs rownumber | 1.1230 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10 | distinct_on vs rownumber | 0.9878 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10 | distinct_on vs lateral | 0.8796 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 100 | lateral vs rownumber | 6.5932 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 100 | distinct_on vs rownumber | 0.9701 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1471 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1000 | lateral vs rownumber | 45.3431 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.9206 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0203 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10000 | lateral vs rownumber | 275.4490 | >1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10000 | distinct_on vs rownumber | 0.6975 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0025 | <1 |
| select-order-aggregate-yugabyte-1-20261004 | yugabyte | single | position | 0 | position control first vs last | 1.7432 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1 | pre_agg vs post_join | 2.8973 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1 | correlated vs post_join | 1.2322 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.4253 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10 | pre_agg vs post_join | 3.2568 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10 | correlated vs post_join | 1.9662 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.6037 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 100 | pre_agg vs post_join | 1.8453 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 100 | correlated vs post_join | 8.5853 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 100 | correlated vs pre_agg | 4.6525 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 1.4047 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1000 | correlated vs post_join | 43.9809 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 31.3092 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.1192 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10000 | correlated vs post_join | 397.8019 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 355.4369 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1 | lateral vs rownumber | 0.6994 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1 | distinct_on vs rownumber | 0.8601 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1 | distinct_on vs lateral | 1.2298 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10 | lateral vs rownumber | 1.1591 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10 | distinct_on vs rownumber | 1.1068 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10 | distinct_on vs lateral | 0.9549 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 100 | lateral vs rownumber | 7.0403 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 100 | distinct_on vs rownumber | 1.2444 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1768 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1000 | lateral vs rownumber | 44.0525 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.9290 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0211 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10000 | lateral vs rownumber | 257.3461 | >1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10000 | distinct_on vs rownumber | 0.9968 | <1 |
| select-order-aggregate-yugabyte-2-20261004 | yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0039 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1 | pre_agg vs post_join | 3.1572 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1 | correlated vs post_join | 1.0681 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1 | correlated vs pre_agg | 0.3383 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10 | pre_agg vs post_join | 2.7474 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10 | correlated vs post_join | 1.9987 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10 | correlated vs pre_agg | 0.7275 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 100 | pre_agg vs post_join | 2.7070 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 100 | correlated vs post_join | 11.8745 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 100 | correlated vs pre_agg | 4.3866 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1000 | pre_agg vs post_join | 0.9922 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1000 | correlated vs post_join | 45.9638 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 1000 | correlated vs pre_agg | 46.3238 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10000 | pre_agg vs post_join | 1.2310 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10000 | correlated vs post_join | 534.0895 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | aggregate | 10000 | correlated vs pre_agg | 433.8691 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1 | lateral vs rownumber | 0.5388 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1 | distinct_on vs rownumber | 0.9334 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1 | distinct_on vs lateral | 1.7324 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10 | lateral vs rownumber | 1.2452 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10 | distinct_on vs rownumber | 1.0035 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10 | distinct_on vs lateral | 0.8059 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 100 | lateral vs rownumber | 6.1051 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 100 | distinct_on vs rownumber | 0.8876 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 100 | distinct_on vs lateral | 0.1454 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1000 | lateral vs rownumber | 49.0169 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1000 | distinct_on vs rownumber | 0.8945 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 1000 | distinct_on vs lateral | 0.0182 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10000 | lateral vs rownumber | 288.8775 | >1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10000 | distinct_on vs rownumber | 0.7886 | <1 |
| select-order-aggregate-yugabyte-3-20261004 | yugabyte | single | topn | 10000 | distinct_on vs lateral | 0.0027 | <1 |

## Recorded limits

- [select-order-aggregate-postgres-1-20261004] one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- [select-order-aggregate-postgres-1-20261004] single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- [select-order-aggregate-postgres-1-20261004] YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
