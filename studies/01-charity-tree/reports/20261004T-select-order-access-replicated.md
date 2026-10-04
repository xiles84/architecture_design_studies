# SELECT access-path matrix — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 6 whole runs summarised. Values are each run's own median; the across-run spread is shown, never merged into one number.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-order-access-postgres-1-20261004 | postgres | single | deep | 42 |
| select-order-access-postgres-2-20261004 | postgres | single | deep | 42 |
| select-order-access-postgres-3-20261004 | postgres | single | deep | 42 |
| select-order-access-yugabyte-1-20261004 | yugabyte | single | deep | 42 |
| select-order-access-yugabyte-2-20261004 | yugabyte | single | deep | 42 |
| select-order-access-yugabyte-3-20261004 | yugabyte | single | deep | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | depth=10000;spelling=keyset | 0.491327, 0.451205, 0.469857 | 0.469857 | 8.5 | 3 |
| postgres | single | depth=10000;spelling=offset | 2.055720, 2.133488, 2.074541 | 2.074541 | 3.7 | 3 |
| postgres | single | depth=1000;spelling=keyset | 0.891678, 1.025613, 1.041333 | 1.025613 | 14.6 | 3 |
| postgres | single | depth=1000;spelling=keyset;position=last | 0.878565 | 0.878565 | 0.0 | 1 |
| postgres | single | depth=1000;spelling=offset | 1.218841, 1.252553, 1.286011 | 1.252553 | 5.4 | 3 |
| postgres | single | depth=100;spelling=keyset | 0.900655, 0.998355, 0.952554 | 0.952554 | 10.3 | 3 |
| postgres | single | depth=100;spelling=offset | 0.846225, 0.855089, 0.884186 | 0.855089 | 4.4 | 3 |
| postgres | single | depth=10;spelling=keyset | 1.049490, 0.948652, 0.983635 | 0.983635 | 10.3 | 3 |
| postgres | single | depth=10;spelling=offset | 0.857571, 0.819740, 0.899717 | 0.857571 | 9.3 | 3 |
| postgres | single | depth=1;spelling=keyset | 1.005356, 0.935028, 0.903688 | 0.935028 | 10.9 | 3 |
| postgres | single | depth=1;spelling=offset | 0.778301, 0.852926, 0.845620 | 0.845620 | 8.8 | 3 |
| postgres | single | shape=broad/cast;spelling=or | 1.407117, 1.245411, 1.324284 | 1.324284 | 12.2 | 3 |
| postgres | single | shape=broad/cast;spelling=union | 1.954402, 1.727877, 1.976684 | 1.954402 | 12.7 | 3 |
| postgres | single | shape=broad/cast;spelling=union_all | 1.723704, 1.927458, 1.793518 | 1.793518 | 11.4 | 3 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=or | 1.211515, 1.202983, 1.341688 | 1.211515 | 11.4 | 3 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=union | 1.685351, 1.713714, 1.665216 | 1.685351 | 2.9 | 3 |
| postgres | single | shape=broad/expression;index=expression_index;spelling=union_all | 1.666657, 1.553154, 1.572190 | 1.572190 | 7.2 | 3 |
| postgres | single | shape=broad/expression;spelling=or | 1.224890, 1.472910, 1.227995 | 1.227995 | 20.2 | 3 |
| postgres | single | shape=broad/expression;spelling=union | 1.799820, 1.833869, 1.712651 | 1.799820 | 6.7 | 3 |
| postgres | single | shape=broad/expression;spelling=union_all | 1.686649, 1.856047, 1.805940 | 1.805940 | 9.4 | 3 |
| postgres | single | shape=broad/indexable;spelling=or | 1.208164, 1.279928, 1.297230 | 1.279928 | 7.0 | 3 |
| postgres | single | shape=broad/indexable;spelling=union | 1.744248, 1.921675, 1.788678 | 1.788678 | 9.9 | 3 |
| postgres | single | shape=broad/indexable;spelling=union_all | 1.840260, 1.799185, 1.781235 | 1.799185 | 3.3 | 3 |
| postgres | single | shape=selective/cast;spelling=or | 0.922202, 0.863186, 0.878092 | 0.878092 | 6.7 | 3 |
| postgres | single | shape=selective/cast;spelling=union | 0.881560, 0.954056, 0.813842 | 0.881560 | 15.9 | 3 |
| postgres | single | shape=selective/cast;spelling=union_all | 0.871287, 0.856929, 0.861363 | 0.861363 | 1.7 | 3 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=or | 0.290745, 0.271766, 0.282496 | 0.282496 | 6.7 | 3 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=union | 0.364596, 0.433833, 0.332002 | 0.364596 | 27.9 | 3 |
| postgres | single | shape=selective/expression;index=expression_index;spelling=union_all | 0.322052, 0.384552, 0.422961 | 0.384552 | 26.2 | 3 |
| postgres | single | shape=selective/expression;spelling=or | 0.879715, 0.970147, 0.931551 | 0.931551 | 9.7 | 3 |
| postgres | single | shape=selective/expression;spelling=union | 0.858044, 0.923608, 0.869854 | 0.869854 | 7.5 | 3 |
| postgres | single | shape=selective/expression;spelling=union_all | 0.884706, 0.835096, 0.873422 | 0.873422 | 5.7 | 3 |
| postgres | single | shape=selective/indexable;spelling=or | 0.354475, 0.424355, 0.374895 | 0.374895 | 18.6 | 3 |
| postgres | single | shape=selective/indexable;spelling=union | 0.353186, 0.382915, 0.428903 | 0.382915 | 19.8 | 3 |
| postgres | single | shape=selective/indexable;spelling=union_all | 0.333533, 0.395387, 0.375063 | 0.375063 | 16.5 | 3 |
| yugabyte | single | depth=10000;spelling=keyset | 4.348607, 4.091508, 4.368005 | 4.348607 | 6.4 | 3 |
| yugabyte | single | depth=10000;spelling=offset | 6.098029, 5.657803, 6.370120 | 6.098029 | 11.7 | 3 |
| yugabyte | single | depth=1000;spelling=keyset | 5.657415, 5.515618, 5.944456 | 5.657415 | 7.6 | 3 |
| yugabyte | single | depth=1000;spelling=offset | 4.544106, 4.399545, 4.946153 | 4.544106 | 12.0 | 3 |
| yugabyte | single | depth=100;spelling=keyset | 5.943401, 5.479473, 6.481624 | 5.943401 | 16.9 | 3 |
| yugabyte | single | depth=100;spelling=offset | 4.552814, 4.316243, 4.783925 | 4.552814 | 10.3 | 3 |
| yugabyte | single | depth=10;spelling=keyset | 5.961666, 5.750417, 5.698592 | 5.750417 | 4.6 | 3 |
| yugabyte | single | depth=10;spelling=offset | 4.695505, 4.345947, 4.310410 | 4.345947 | 8.9 | 3 |
| yugabyte | single | depth=1;spelling=keyset | 5.941060, 5.535822, 6.102109 | 5.941060 | 9.5 | 3 |
| yugabyte | single | depth=1;spelling=offset | 4.433638, 4.351427, 4.481934 | 4.433638 | 2.9 | 3 |
| yugabyte | single | depth=1;spelling=offset;position=last | 4.329512 | 4.329512 | 0.0 | 1 |
| yugabyte | single | shape=broad/cast;spelling=or | 6.191972, 6.324200, 6.541074 | 6.324200 | 5.5 | 3 |
| yugabyte | single | shape=broad/cast;spelling=union | 10.073319, 8.814304, 9.874776 | 9.874776 | 12.7 | 3 |
| yugabyte | single | shape=broad/cast;spelling=union_all | 9.171276, 8.582023, 8.839165 | 8.839165 | 6.7 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=or | 4.811472, 5.338468, 4.918317 | 4.918317 | 10.7 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=union | 7.166819, 7.478430, 7.451968 | 7.451968 | 4.2 | 3 |
| yugabyte | single | shape=broad/expression;index=expression_index;spelling=union_all | 6.639805, 7.154042, 7.321621 | 7.154042 | 9.5 | 3 |
| yugabyte | single | shape=broad/expression;spelling=or | 6.476607, 6.550784, 6.576510 | 6.550784 | 1.5 | 3 |
| yugabyte | single | shape=broad/expression;spelling=union | 9.538274, 9.132888, 9.553295 | 9.538274 | 4.4 | 3 |
| yugabyte | single | shape=broad/expression;spelling=union_all | 9.331759, 8.647626, 9.077006 | 9.077006 | 7.5 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=or | 6.172420, 5.758767, 6.002137 | 6.002137 | 6.9 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=union | 9.360888, 8.673496, 9.221627 | 9.221627 | 7.5 | 3 |
| yugabyte | single | shape=broad/indexable;spelling=union_all | 8.929360, 8.585931, 8.840938 | 8.840938 | 3.9 | 3 |
| yugabyte | single | shape=selective/cast;spelling=or | 3.976737, 3.641978, 4.024952 | 3.976737 | 9.6 | 3 |
| yugabyte | single | shape=selective/cast;spelling=union | 6.602031, 6.269404, 6.523992 | 6.523992 | 5.1 | 3 |
| yugabyte | single | shape=selective/cast;spelling=union_all | 6.503088, 6.065570, 6.437079 | 6.437079 | 6.8 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=or | 3.094914, 3.300704, 3.194597 | 3.194597 | 6.4 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=union | 4.846780, 5.508483, 5.084925 | 5.084925 | 13.0 | 3 |
| yugabyte | single | shape=selective/expression;index=expression_index;spelling=union_all | 4.969379, 5.334283, 5.232050 | 5.232050 | 7.0 | 3 |
| yugabyte | single | shape=selective/expression;spelling=or | 4.125288, 3.717254, 3.989447 | 3.989447 | 10.2 | 3 |
| yugabyte | single | shape=selective/expression;spelling=union | 6.553698, 6.354504, 6.689376 | 6.553698 | 5.1 | 3 |
| yugabyte | single | shape=selective/expression;spelling=union_all | 6.374217, 6.242359, 6.783315 | 6.374217 | 8.5 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=or | 3.618328, 3.213270, 3.612128 | 3.612128 | 11.2 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=union | 6.067404, 5.996545, 6.175522 | 6.067404 | 2.9 | 3 |
| yugabyte | single | shape=selective/indexable;spelling=union_all | 6.173958, 5.658665, 6.247534 | 6.173958 | 9.5 | 3 |

## Pair ratios per run (per engine)

| engine | topology | pair | run ids | ratios (B/A) | signs |
|---|---|---|---|---|---|
| postgres | single | depth=1 offset vs keyset | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.2917, 1.0963, 1.0687 | >1, >1, >1 |
| postgres | single | depth=10 offset vs keyset | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.2238, 1.1573, 1.0933 | >1, >1, >1 |
| postgres | single | depth=100 offset vs keyset | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.0643, 1.1675, 1.0773 | >1, >1, >1 |
| postgres | single | depth=1000 offset vs keyset | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.7316, 0.8188, 0.8097 | <1, <1, <1 |
| postgres | single | depth=10000 offset vs keyset | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.2390, 0.2115, 0.2265 | <1, <1, <1 |
| postgres | single | position control first vs last | select-order-access-postgres-1-20261004 | 0.9853 | <1 |
| postgres | single | shape=broad/cast or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.3889, 1.3874, 1.4926 | >1, >1, >1 |
| postgres | single | shape=broad/cast or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.2250, 1.5476, 1.3543 | >1, >1, >1 |
| postgres | single | shape=broad/expression or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.4694, 1.2451, 1.3947 | >1, >1, >1 |
| postgres | single | shape=broad/expression or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.3770, 1.2601, 1.4706 | >1, >1, >1 |
| postgres | single | shape=broad/indexable or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.4437, 1.5014, 1.3788 | >1, >1, >1 |
| postgres | single | shape=broad/indexable or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.5232, 1.4057, 1.3731 | >1, >1, >1 |
| postgres | single | shape=selective/cast or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.9559, 1.1053, 0.9268 | <1, >1, <1 |
| postgres | single | shape=selective/cast or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.9448, 0.9928, 0.9809 | <1, <1, <1 |
| postgres | single | shape=selective/expression or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.9754, 0.9520, 0.9338 | <1, <1, <1 |
| postgres | single | shape=selective/expression or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 1.0057, 0.8608, 0.9376 | >1, <1, <1 |
| postgres | single | shape=selective/indexable or vs union | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.9964, 0.9023, 1.1441 | <1, <1, >1 |
| postgres | single | shape=selective/indexable or vs union_all | select-order-access-postgres-1-20261004, select-order-access-postgres-2-20261004, select-order-access-postgres-3-20261004 | 0.9409, 0.9317, 1.0004 | <1, <1, >1 |
| yugabyte | single | depth=1 offset vs keyset | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.3400, 1.2722, 1.3615 | >1, >1, >1 |
| yugabyte | single | depth=10 offset vs keyset | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.2697, 1.3232, 1.3221 | >1, >1, >1 |
| yugabyte | single | depth=100 offset vs keyset | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.3054, 1.2695, 1.3549 | >1, >1, >1 |
| yugabyte | single | depth=1000 offset vs keyset | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.2450, 1.2537, 1.2018 | >1, >1, >1 |
| yugabyte | single | depth=10000 offset vs keyset | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 0.7131, 0.7232, 0.6857 | <1, <1, <1 |
| yugabyte | single | position control first vs last | select-order-access-yugabyte-1-20261004 | 0.9765 | <1 |
| yugabyte | single | shape=broad/cast or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.6268, 1.3937, 1.5097 | >1, >1, >1 |
| yugabyte | single | shape=broad/cast or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.4812, 1.3570, 1.3513 | >1, >1, >1 |
| yugabyte | single | shape=broad/expression or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.4727, 1.3942, 1.4526 | >1, >1, >1 |
| yugabyte | single | shape=broad/expression or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.4408, 1.3201, 1.3802 | >1, >1, >1 |
| yugabyte | single | shape=broad/indexable or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.5166, 1.5061, 1.5364 | >1, >1, >1 |
| yugabyte | single | shape=broad/indexable or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.4467, 1.4909, 1.4730 | >1, >1, >1 |
| yugabyte | single | shape=selective/cast or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.6602, 1.7214, 1.6209 | >1, >1, >1 |
| yugabyte | single | shape=selective/cast or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.6353, 1.6655, 1.5993 | >1, >1, >1 |
| yugabyte | single | shape=selective/expression or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.5887, 1.7095, 1.6768 | >1, >1, >1 |
| yugabyte | single | shape=selective/expression or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.5452, 1.6793, 1.7003 | >1, >1, >1 |
| yugabyte | single | shape=selective/indexable or vs union | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.6769, 1.8662, 1.7097 | >1, >1, >1 |
| yugabyte | single | shape=selective/indexable or vs union_all | select-order-access-yugabyte-1-20261004, select-order-access-yugabyte-2-20261004, select-order-access-yugabyte-3-20261004 | 1.7063, 1.7610, 1.7296 | >1, >1, >1 |

## Recorded limits

- [select-order-access-postgres-1-20261004] no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- [select-order-access-postgres-1-20261004] one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- [select-order-access-postgres-1-20261004] single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host
