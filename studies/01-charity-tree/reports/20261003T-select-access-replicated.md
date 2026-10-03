# SELECT access-path matrix — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 7 whole runs summarised. Values are each run's own median; the across-run spread is shown, never merged into one number.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-access-pg-1-20261003 | postgres | single | deep | 42 |
| select-access-pg-2-20261003 | postgres | single | deep | 42 |
| select-access-pg-3-20261003 | postgres | single | deep | 42 |
| select-access-pg-pos-20261003 | postgres | single | deep | 42 |
| select-access-yb-1-20261003 | yugabyte | single | deep | 42 |
| select-access-yb-2-20261003 | yugabyte | single | deep | 42 |
| select-access-yb-3-20261003 | yugabyte | single | deep | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | depth=10000;spelling=keyset | 0.487803, 0.427452, 0.486106, 0.457182 | 0.471644 | 12.8 | 4 |
| postgres | single | depth=10000;spelling=offset | 2.131418, 2.171211, 2.080571, 1.982273 | 2.105994 | 9.0 | 4 |
| postgres | single | depth=1000;spelling=keyset | 1.004620, 0.871633, 0.930056, 0.893647 | 0.911852 | 14.6 | 4 |
| postgres | single | depth=1000;spelling=offset | 1.207265, 1.234014, 1.270472, 1.192947 | 1.220639 | 6.4 | 4 |
| postgres | single | depth=100;spelling=keyset | 0.970568, 0.920448, 0.903991, 0.891731 | 0.912219 | 8.6 | 4 |
| postgres | single | depth=100;spelling=offset | 0.880969, 0.867012, 0.931667, 0.902112 | 0.891541 | 7.3 | 4 |
| postgres | single | depth=10;spelling=keyset | 0.960248, 0.955337, 1.066094, 0.939465 | 0.957793 | 13.2 | 4 |
| postgres | single | depth=10;spelling=offset | 0.879777, 0.820651, 0.909782, 0.827168 | 0.853473 | 10.4 | 4 |
| postgres | single | depth=1;spelling=keyset | 1.013290, 0.969571, 0.986581, 0.931700 | 0.978076 | 8.3 | 4 |
| postgres | single | depth=1;spelling=offset | 0.884729, 0.814058, 0.819929, 0.811851 | 0.816993 | 8.9 | 4 |
| postgres | single | shape=broad/cast;spelling=or | 1.271384, 1.329183, 1.332954, 1.290834 | 1.310008 | 4.7 | 4 |
| postgres | single | shape=broad/cast;spelling=union | 1.971766, 1.684447, 1.720254, 1.696035 | 1.708144 | 16.8 | 4 |
| postgres | single | shape=broad/cast;spelling=union_all | 1.902759, 1.927206, 1.848675, 1.875858 | 1.889309 | 4.2 | 4 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=or | 1.626503, 1.258192, 1.344989, 1.286213 | 1.315601 | 28.0 | 4 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=union | 1.833793, 1.552454, 1.673449, 1.651826 | 1.662637 | 16.9 | 4 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=union_all | 2.038856, 1.617218, 1.880469, 1.629259 | 1.754864 | 24.0 | 4 |
| postgres | single | shape=broad/expression;spelling=or | 2.196484, 1.222194, 1.277451, 1.355900 | 1.316676 | 74.0 | 4 |
| postgres | single | shape=broad/expression;spelling=union | 1.758007, 1.839686, 1.747708, 1.633648 | 1.752858 | 11.8 | 4 |
| postgres | single | shape=broad/expression;spelling=union_all | 1.766537, 1.940389, 1.759378, 1.925405 | 1.845971 | 9.8 | 4 |
| postgres | single | shape=broad/indexable;spelling=or | 1.321904, 1.177284, 1.152580, 1.834578 | 1.249594 | 54.6 | 4 |
| postgres | single | shape=broad/indexable;spelling=union | 1.798644, 1.914575, 1.812832, 1.798113 | 1.805738 | 6.4 | 4 |
| postgres | single | shape=broad/indexable;spelling=union_all | 1.658562, 1.865953, 1.848096, 1.679676 | 1.763886 | 11.8 | 4 |
| postgres | single | shape=selective/cast;spelling=or | 0.872541, 0.841638, 0.849832, 0.991944 | 0.861187 | 17.5 | 4 |
| postgres | single | shape=selective/cast;spelling=union | 0.876240, 0.737233, 0.809649, 0.851168 | 0.830408 | 16.7 | 4 |
| postgres | single | shape=selective/cast;spelling=union_all | 0.767099, 0.773480, 0.854721, 0.898037 | 0.814100 | 16.1 | 4 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=or | 0.288316, 0.265954, 0.250803, 0.264256 | 0.265105 | 14.2 | 4 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=union | 0.309812, 0.300824, 0.283988, 0.318259 | 0.305318 | 11.2 | 4 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=union_all | 0.318378, 0.279172, 0.314781, 0.336840 | 0.316580 | 18.2 | 4 |
| postgres | single | shape=selective/expression;spelling=or | 0.824664, 0.827197, 0.895906, 0.863765 | 0.845481 | 8.4 | 4 |
| postgres | single | shape=selective/expression;spelling=union | 0.843858, 0.806994, 0.828612, 0.867208 | 0.836235 | 7.2 | 4 |
| postgres | single | shape=selective/expression;spelling=union_all | 0.846114, 0.824431, 0.858703, 0.841242 | 0.843678 | 4.1 | 4 |
| postgres | single | shape=selective/indexable;spelling=or | 0.310498, 0.253115, 0.338102, 0.318574 | 0.314536 | 27.0 | 4 |
| postgres | single | shape=selective/indexable;spelling=or;position=last | 0.281925 | 0.281925 | 0.0 | 1 |
| postgres | single | shape=selective/indexable;spelling=union | 0.339654, 0.298157, 0.369664, 0.307091 | 0.323373 | 22.1 | 4 |
| postgres | single | shape=selective/indexable;spelling=union_all | 0.280965, 0.296781, 0.464899, 0.325088 | 0.310935 | 59.2 | 4 |
| yugabyte | single | depth=10000;spelling=keyset | 4.685009, 4.998019, 11.682140 | 4.998019 | 140.0 | 3 |
| yugabyte | single | depth=10000;spelling=offset | 8.699910, 7.418655, 16.887356 | 8.699910 | 108.8 | 3 |
| yugabyte | single | depth=1000;spelling=keyset | 7.721890, 6.822253, 16.402403 | 7.721890 | 124.1 | 3 |
| yugabyte | single | depth=1000;spelling=offset | 5.392865, 5.836768, 12.319564 | 5.836768 | 118.7 | 3 |
| yugabyte | single | depth=100;spelling=keyset | 6.519669, 7.006257, 17.572849 | 7.006257 | 157.8 | 3 |
| yugabyte | single | depth=100;spelling=offset | 5.810455, 6.360430, 14.058823 | 6.360430 | 129.7 | 3 |
| yugabyte | single | depth=10;spelling=keyset | 8.163653, 8.254318, 16.486891 | 8.254318 | 100.8 | 3 |
| yugabyte | single | depth=10;spelling=offset | 6.419176, 6.560888, 13.942363 | 6.560888 | 114.7 | 3 |
| yugabyte | single | depth=1;spelling=keyset | 7.588768, 7.936387, 16.705487 | 7.936387 | 114.9 | 3 |
| yugabyte | single | depth=1;spelling=offset | 6.020191, 6.739808, 14.194939 | 6.739808 | 121.3 | 3 |
| yugabyte | single | shape=broad/cast;spelling=or | 5.366463, 6.986492, 16.740603 | 6.986492 | 162.8 | 3 |
| yugabyte | single | shape=broad/cast;spelling=union | 11.366105, 11.844361, 25.919721 | 11.844361 | 122.9 | 3 |
| yugabyte | single | shape=broad/cast;spelling=union_all | 11.271674, 10.534562, 23.939678 | 11.271674 | 118.9 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=or | 7.570656, 7.119667, 19.239256 | 7.570656 | 160.1 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=union | 8.046960, 11.705941, 26.107827 | 11.705941 | 154.3 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=union_all | 7.561875, 10.532490, 25.699853 | 10.532490 | 172.2 | 3 |
| yugabyte | single | shape=broad/expression;spelling=or | 5.237082, 7.530308, 16.901565 | 7.530308 | 154.9 | 3 |
| yugabyte | single | shape=broad/expression;spelling=union | 7.989638, 11.142094, 27.355109 | 11.142094 | 173.8 | 3 |
| yugabyte | single | shape=broad/expression;spelling=union_all | 7.314392, 10.005592, 30.641029 | 10.005592 | 233.1 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=or | 5.001992, 8.036357, 19.663053 | 8.036357 | 182.4 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=union | 7.570294, 9.568491, 25.104787 | 9.568491 | 183.3 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=union_all | 7.029154, 10.539168, 22.170745 | 10.539168 | 143.7 | 3 |
| yugabyte | single | shape=selective/cast;spelling=or | 2.983885, 4.038333, 11.706365 | 4.038333 | 216.0 | 3 |
| yugabyte | single | shape=selective/cast;spelling=union | 5.144890, 6.878930, 17.062551 | 6.878930 | 173.2 | 3 |
| yugabyte | single | shape=selective/cast;spelling=union_all | 5.063655, 7.142692, 16.400881 | 7.142692 | 158.7 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=or | 3.323720, 4.153415, 9.738460 | 4.153415 | 154.4 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=union | 6.326547, 7.431587, 16.287180 | 7.431587 | 134.0 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=union_all | 6.220068, 7.027629, 15.602001 | 7.027629 | 133.5 | 3 |
| yugabyte | single | shape=selective/expression;spelling=or | 3.143710, 4.051543, 12.202750 | 4.051543 | 223.6 | 3 |
| yugabyte | single | shape=selective/expression;spelling=union | 4.946131, 7.521499, 19.574730 | 7.521499 | 194.5 | 3 |
| yugabyte | single | shape=selective/expression;spelling=union_all | 5.245558, 7.011259, 21.653088 | 7.011259 | 234.0 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=or | 2.822762, 3.547935, 8.586247 | 3.547935 | 162.4 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=or;position=last | 3.944609 | 3.944609 | 0.0 | 1 |
| yugabyte | single | shape=selective/indexable;spelling=union | 4.899867, 6.531859, 14.989680 | 6.531859 | 154.5 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=union_all | 4.969804, 6.273535, 14.813902 | 6.273535 | 156.9 | 3 |

## Pair ratios per run (per engine)

| engine | topology | pair | ratios (B/A) | signs |
|---|---|---|---|---|
| postgres | single | depth=1 offset vs keyset | 1.1453, 1.1910, 1.2033, 1.1476 | >1, >1, >1, >1 |
| postgres | single | depth=10 offset vs keyset | 1.0915, 1.1641, 1.1718, 1.1358 | >1, >1, >1, >1 |
| postgres | single | depth=100 offset vs keyset | 1.1017, 1.0616, 0.9703, 0.9885 | >1, >1, <1, <1 |
| postgres | single | depth=1000 offset vs keyset | 0.8321, 0.7063, 0.7321, 0.7491 | <1, <1, <1, <1 |
| postgres | single | depth=10000 offset vs keyset | 0.2289, 0.1969, 0.2336, 0.2306 | <1, <1, <1, <1 |
| postgres | single | position control selective/indexable or first vs last | 0.8850 | <1 |
| postgres | single | shape=broad/cast or vs union | 1.5509, 1.2673, 1.2906, 1.3139 | >1, >1, >1, >1 |
| postgres | single | shape=broad/cast or vs union_all | 1.4966, 1.4499, 1.3869, 1.4532 | >1, >1, >1, >1 |
| postgres | single | shape=broad/expression or vs union | 0.8004, 1.5052, 1.3681, 1.2048 | <1, >1, >1, >1 |
| postgres | single | shape=broad/expression or vs union_all | 0.8043, 1.5876, 1.3773, 1.4200 | <1, >1, >1, >1 |
| postgres | single | shape=broad/indexable or vs union | 1.3606, 1.6263, 1.5728, 0.9801 | >1, >1, >1, <1 |
| postgres | single | shape=broad/indexable or vs union_all | 1.2547, 1.5850, 1.6034, 0.9156 | >1, >1, >1, <1 |
| postgres | single | shape=selective/cast or vs union | 1.0042, 0.8760, 0.9527, 0.8581 | >1, <1, <1, <1 |
| postgres | single | shape=selective/cast or vs union_all | 0.8792, 0.9190, 1.0058, 0.9053 | <1, <1, >1, <1 |
| postgres | single | shape=selective/expression or vs union | 1.0233, 0.9756, 0.9249, 1.0040 | >1, <1, <1, >1 |
| postgres | single | shape=selective/expression or vs union_all | 1.0260, 0.9967, 0.9585, 0.9739 | >1, <1, <1, <1 |
| postgres | single | shape=selective/indexable or vs union | 1.0939, 1.1780, 1.0934, 0.9640 | >1, >1, >1, <1 |
| postgres | single | shape=selective/indexable or vs union_all | 0.9049, 1.1725, 1.3750, 1.0204 | <1, >1, >1, >1 |
| yugabyte | single | depth=1 offset vs keyset | 1.2606, 1.1775, 1.1769 | >1, >1, >1 |
| yugabyte | single | depth=10 offset vs keyset | 1.2718, 1.2581, 1.1825 | >1, >1, >1 |
| yugabyte | single | depth=100 offset vs keyset | 1.1221, 1.1015, 1.2500 | >1, >1, >1 |
| yugabyte | single | depth=1000 offset vs keyset | 1.4319, 1.1688, 1.3314 | >1, >1, >1 |
| yugabyte | single | depth=10000 offset vs keyset | 0.5385, 0.6737, 0.6918 | <1, <1, <1 |
| yugabyte | single | position control selective/indexable or first vs last | 1.1118 | >1 |
| yugabyte | single | shape=broad/cast or vs union | 2.1180, 1.6953, 1.5483 | >1, >1, >1 |
| yugabyte | single | shape=broad/cast or vs union_all | 2.1004, 1.5078, 1.4300 | >1, >1, >1 |
| yugabyte | single | shape=broad/expression or vs union | 1.5256, 1.4796, 1.6185 | >1, >1, >1 |
| yugabyte | single | shape=broad/expression or vs union_all | 1.3967, 1.3287, 1.8129 | >1, >1, >1 |
| yugabyte | single | shape=broad/indexable or vs union | 1.5135, 1.1907, 1.2767 | >1, >1, >1 |
| yugabyte | single | shape=broad/indexable or vs union_all | 1.4053, 1.3114, 1.1275 | >1, >1, >1 |
| yugabyte | single | shape=selective/cast or vs union | 1.7242, 1.7034, 1.4575 | >1, >1, >1 |
| yugabyte | single | shape=selective/cast or vs union_all | 1.6970, 1.7687, 1.4010 | >1, >1, >1 |
| yugabyte | single | shape=selective/expression or vs union | 1.5733, 1.8565, 1.6041 | >1, >1, >1 |
| yugabyte | single | shape=selective/expression or vs union_all | 1.6686, 1.7305, 1.7744 | >1, >1, >1 |
| yugabyte | single | shape=selective/indexable or vs union | 1.7358, 1.8410, 1.7458 | >1, >1, >1 |
| yugabyte | single | shape=selective/indexable or vs union_all | 1.7606, 1.7682, 1.7253 | >1, >1, >1 |

## Recorded limits

- [select-access-pg-1-20261003] no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- [select-access-pg-1-20261003] one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- [select-access-pg-1-20261003] single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host
