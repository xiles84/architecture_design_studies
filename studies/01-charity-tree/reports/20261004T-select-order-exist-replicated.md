# Study 01 SELECT membership & absence — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 6 whole runs summarised, each run's own median shown, never merged across engines or fanout levels.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-order-existence-postgres-1-20261004 | postgres | single | fanout | 42 |
| select-order-existence-postgres-2-20261004 | postgres | single | fanout | 42 |
| select-order-existence-postgres-3-20261004 | postgres | single | fanout | 42 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | fanout | 42 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | fanout | 42 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | fanout | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | family=absence;stmt=left_is_null;fanout=1 | 1.627203, 1.746013, 1.534601 | 1.627203 | 13.0 | 3 |
| postgres | single | family=absence;stmt=left_is_null;fanout=10 | 0.730025, 0.767388, 0.757808 | 0.757808 | 4.9 | 3 |
| postgres | single | family=absence;stmt=left_is_null;fanout=100 | 0.690735, 0.653910, 0.643376 | 0.653910 | 7.2 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=1 | 1.508953, 1.642936, 1.742235 | 1.642936 | 14.2 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=10 | 0.711377, 0.715728, 0.720443 | 0.715728 | 1.3 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=100 | 0.211010, 0.241261, 0.233455 | 0.233455 | 13.0 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=1 | 1.790824, 1.841238, 1.900687 | 1.841238 | 6.0 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=10 | 0.823264, 0.806535, 0.804945 | 0.806535 | 2.3 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=100 | 0.667245, 0.715350, 0.665071 | 0.667245 | 7.5 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=1 | 1.848299, 2.219894, 2.144528 | 2.144528 | 17.3 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=10 | 1.046127, 1.031887, 0.988961 | 1.031887 | 5.5 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=100 | 0.258465, 0.283102, 0.268153 | 0.268153 | 9.2 | 3 |
| postgres | single | family=membership;stmt=in;fanout=1 | 2.398619, 1.901901, 2.068100 | 2.068100 | 24.0 | 3 |
| postgres | single | family=membership;stmt=in;fanout=10 | 0.984078, 1.019102, 1.086892 | 1.019102 | 10.1 | 3 |
| postgres | single | family=membership;stmt=in;fanout=100 | 0.282208, 0.266355, 0.272380 | 0.272380 | 5.8 | 3 |
| postgres | single | family=membership;stmt=in;fanout=1;position=last | 1.876086 | 1.876086 | 0.0 | 1 |
| postgres | single | family=membership;stmt=in_notnull;fanout=1 | 2.338639, 2.182980, 2.608755 | 2.338639 | 18.2 | 3 |
| postgres | single | family=membership;stmt=in_notnull;fanout=10 | 1.006667, 1.073310, 1.109824 | 1.073310 | 9.6 | 3 |
| postgres | single | family=membership;stmt=in_notnull;fanout=100 | 0.294357, 0.303358, 0.261529 | 0.294357 | 14.2 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=1 | 3.122387, 3.103390, 3.202961 | 3.122387 | 3.2 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=10 | 1.442715, 1.468340, 1.427715 | 1.442715 | 2.8 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=100 | 1.050418, 1.030008, 1.133713 | 1.050418 | 9.9 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=1 | 10.765602, 11.292946, 10.883540 | 10.883540 | 4.8 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=10 | 5.035923, 6.125948, 5.325422 | 5.325422 | 20.5 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=100 | 4.965958, 5.036845, 4.730838 | 4.965958 | 6.2 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=1 | 10.356080, 11.226823, 10.519197 | 10.519197 | 8.3 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=10 | 4.893101, 5.790081, 4.902011 | 4.902011 | 18.3 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=100 | 4.693781, 4.979181, 4.554114 | 4.693781 | 9.1 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=1;position=last | 11.714273 | 11.714273 | 0.0 | 1 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=1 | 11.298211, 11.192847, 10.769454 | 11.192847 | 4.7 | 3 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=10 | 5.087358, 6.061813, 5.337235 | 5.337235 | 18.3 | 3 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=100 | 5.147238, 5.407629, 4.732829 | 5.147238 | 13.1 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=1 | 11.417515, 11.933313, 11.183983 | 11.417515 | 6.6 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=10 | 5.339044, 6.301696, 5.650517 | 5.650517 | 17.0 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=100 | 4.966368, 5.003641, 4.733273 | 4.966368 | 5.4 | 3 |
| yugabyte | single | family=membership;stmt=in;fanout=1 | 11.571033, 12.625299, 11.632632 | 11.632632 | 9.1 | 3 |
| yugabyte | single | family=membership;stmt=in;fanout=10 | 5.686797, 6.512661, 5.301876 | 5.686797 | 21.3 | 3 |
| yugabyte | single | family=membership;stmt=in;fanout=100 | 4.984291, 4.863765, 4.704860 | 4.863765 | 5.7 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=1 | 11.970625, 12.851232, 12.216653 | 12.216653 | 7.2 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=10 | 5.773285, 6.766019, 5.733477 | 5.773285 | 17.9 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=100 | 5.269807, 5.493082, 5.289940 | 5.289940 | 4.2 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=1 | 12.677520, 13.133680, 13.014586 | 13.014586 | 3.5 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=10 | 5.414712, 5.993135, 5.654911 | 5.654911 | 10.2 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=100 | 5.024495, 5.071254, 4.793187 | 5.024495 | 5.5 | 3 |

## Membership pair ratios per run (per engine)

| run | engine | topology | fanout | pair | ratio (B/A) | sign |
|---|---|---|---|---|---|---|
| select-order-existence-postgres-1-20261004 | postgres | single | 1 | in vs exists | 1.2977 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 1 | in_notnull vs exists | 1.2653 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 1 | join_distinct vs in | 1.3017 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 1 | join_distinct vs exists | 1.6893 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 10 | in vs exists | 0.9407 | <1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 10 | in_notnull vs exists | 0.9623 | <1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 10 | join_distinct vs in | 1.4661 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 10 | join_distinct vs exists | 1.3791 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 100 | in vs exists | 1.0919 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 100 | in_notnull vs exists | 1.1389 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 100 | join_distinct vs in | 3.7221 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 100 | join_distinct vs exists | 4.0641 | >1 |
| select-order-existence-postgres-1-20261004 | postgres | single | 1 | position control first vs last | 0.7822 | <1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 1 | in vs exists | 0.8568 | <1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 1 | in_notnull vs exists | 0.9834 | <1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 1 | join_distinct vs in | 1.6317 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 1 | join_distinct vs exists | 1.3980 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 10 | in vs exists | 0.9876 | <1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 10 | in_notnull vs exists | 1.0401 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 10 | join_distinct vs in | 1.4408 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 10 | join_distinct vs exists | 1.4230 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 100 | in vs exists | 0.9408 | <1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 100 | in_notnull vs exists | 1.0716 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 100 | join_distinct vs in | 3.8670 | >1 |
| select-order-existence-postgres-2-20261004 | postgres | single | 100 | join_distinct vs exists | 3.6383 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 1 | in vs exists | 0.9644 | <1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 1 | in_notnull vs exists | 1.2165 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 1 | join_distinct vs in | 1.5487 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 1 | join_distinct vs exists | 1.4936 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 10 | in vs exists | 1.0990 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 10 | in_notnull vs exists | 1.1222 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 10 | join_distinct vs in | 1.3136 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 10 | join_distinct vs exists | 1.4437 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 100 | in vs exists | 1.0158 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 100 | in_notnull vs exists | 0.9753 | <1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 100 | join_distinct vs in | 4.1622 | >1 |
| select-order-existence-postgres-3-20261004 | postgres | single | 100 | join_distinct vs exists | 4.2279 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 1 | in vs exists | 1.0134 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 1 | in_notnull vs exists | 1.0484 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 1 | join_distinct vs in | 1.0956 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 1 | join_distinct vs exists | 1.1104 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 10 | in vs exists | 1.0651 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 10 | in_notnull vs exists | 1.0813 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 10 | join_distinct vs in | 0.9522 | <1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 10 | join_distinct vs exists | 1.0142 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 100 | in vs exists | 1.0036 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 100 | in_notnull vs exists | 1.0611 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 100 | join_distinct vs in | 1.0081 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 100 | join_distinct vs exists | 1.0117 | >1 |
| select-order-existence-yugabyte-1-20261004 | yugabyte | single | 1 | position control first vs last | 1.1311 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 1 | in vs exists | 1.0580 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 1 | in_notnull vs exists | 1.0769 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 1 | join_distinct vs in | 1.0403 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 1 | join_distinct vs exists | 1.1006 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 10 | in vs exists | 1.0335 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 10 | in_notnull vs exists | 1.0737 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 10 | join_distinct vs in | 0.9202 | <1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 10 | join_distinct vs exists | 0.9510 | <1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 100 | in vs exists | 0.9720 | <1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 100 | in_notnull vs exists | 1.0978 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 100 | join_distinct vs in | 1.0427 | >1 |
| select-order-existence-yugabyte-2-20261004 | yugabyte | single | 100 | join_distinct vs exists | 1.0135 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 1 | in vs exists | 1.0401 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 1 | in_notnull vs exists | 1.0923 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 1 | join_distinct vs in | 1.1188 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 1 | join_distinct vs exists | 1.1637 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 10 | in vs exists | 0.9383 | <1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 10 | in_notnull vs exists | 1.0147 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 10 | join_distinct vs in | 1.0666 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 10 | join_distinct vs exists | 1.0008 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 100 | in vs exists | 0.9940 | <1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 100 | in_notnull vs exists | 1.1176 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 100 | join_distinct vs in | 1.0188 | >1 |
| select-order-existence-yugabyte-3-20261004 | yugabyte | single | 100 | join_distinct vs exists | 1.0127 | >1 |

## Recorded limits

- [select-order-existence-postgres-1-20261004] one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- [select-order-existence-postgres-1-20261004] single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- [select-order-existence-postgres-1-20261004] qualifying proportion, parent hit fraction and hot-parent skew are not swept in this bounded run
- [select-order-existence-postgres-1-20261004] YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
