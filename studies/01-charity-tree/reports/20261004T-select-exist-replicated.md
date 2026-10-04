# Study 01 SELECT membership & absence — replicated whole runs (generated measurements)

_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._

**TL;DR:** 6 whole runs summarised, each run's own median shown, never merged across engines or fanout levels.

## Runs

| run | engine | topology | scale | seed |
|---|---|---|---|---|
| select-exist-pg-1b-20261004 | postgres | single | fanout | 42 |
| select-exist-pg-2-20261004 | postgres | single | fanout | 42 |
| select-exist-pg-3-20261004 | postgres | single | fanout | 42 |
| select-exist-yb-1-20261004 | yugabyte | single | fanout | 42 |
| select-exist-yb-2-20261004 | yugabyte | single | fanout | 42 |
| select-exist-yb-3-20261004 | yugabyte | single | fanout | 42 |

## Cells across runs (per engine)

| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |
|---|---|---|---|---|---|---|
| postgres | single | family=absence;stmt=left_is_null;fanout=1 | 1.622063, 1.535428, 1.502038 | 1.535428 | 7.8 | 3 |
| postgres | single | family=absence;stmt=left_is_null;fanout=10 | 0.738915, 0.705977, 0.764227 | 0.738915 | 7.9 | 3 |
| postgres | single | family=absence;stmt=left_is_null;fanout=100 | 0.716452, 0.629697, 0.627223 | 0.629697 | 14.2 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=1 | 1.704936, 1.338610, 1.428544 | 1.428544 | 25.6 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=10 | 0.683355, 0.691725, 0.693355 | 0.691725 | 1.4 | 3 |
| postgres | single | family=absence;stmt=not_exists;fanout=100 | 0.242070, 0.220823, 0.340455 | 0.242070 | 49.4 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=1 | 2.515485, 1.795916, 1.802199 | 1.802199 | 39.9 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=10 | 0.792541, 0.803190, 0.771177 | 0.792541 | 4.0 | 3 |
| postgres | single | family=absence;stmt=not_in_notnull;fanout=100 | 0.724409, 0.677712, 0.674786 | 0.677712 | 7.3 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=1 | 2.195002, 2.131830, 1.830639 | 2.131830 | 17.1 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=10 | 1.038681, 1.017468, 1.089431 | 1.038681 | 6.9 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=100 | 0.237410, 0.264915, 0.298123 | 0.264915 | 22.9 | 3 |
| postgres | single | family=membership;stmt=exists;fanout=1;position=last | 1.856280 | 1.856280 | 0.0 | 1 |
| postgres | single | family=membership;stmt=in;fanout=1 | 1.995845, 1.828571, 1.923109 | 1.923109 | 8.7 | 3 |
| postgres | single | family=membership;stmt=in;fanout=10 | 1.006198, 1.014969, 1.045929 | 1.014969 | 3.9 | 3 |
| postgres | single | family=membership;stmt=in;fanout=100 | 0.268753, 0.297646, 0.268862 | 0.268862 | 10.7 | 3 |
| postgres | single | family=membership;stmt=in_notnull;fanout=1 | 2.221116, 2.187308, 1.823572 | 2.187308 | 18.2 | 3 |
| postgres | single | family=membership;stmt=in_notnull;fanout=10 | 1.047735, 0.987456, 1.026126 | 1.026126 | 5.9 | 3 |
| postgres | single | family=membership;stmt=in_notnull;fanout=100 | 0.278093, 0.296680, 0.267351 | 0.278093 | 10.5 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=1 | 3.202929, 3.078255, 2.693038 | 3.078255 | 16.6 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=10 | 1.402433, 1.319520, 1.706179 | 1.402433 | 27.6 | 3 |
| postgres | single | family=membership;stmt=join_distinct;fanout=100 | 1.043884, 1.056820, 1.076828 | 1.056820 | 3.1 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=1 | 8.584150, 8.013683, 7.842389 | 8.013683 | 9.3 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=10 | 3.589667, 3.988009, 3.941372 | 3.941372 | 10.1 | 3 |
| yugabyte | single | family=absence;stmt=left_is_null;fanout=100 | 3.520505, 3.493895, 3.655629 | 3.520505 | 4.6 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=1 | 8.304120, 7.930311, 8.002279 | 8.002279 | 4.7 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=10 | 3.424493, 3.827463, 3.723942 | 3.723942 | 10.8 | 3 |
| yugabyte | single | family=absence;stmt=not_exists;fanout=100 | 3.683749, 3.547002, 3.632083 | 3.632083 | 3.8 | 3 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=1 | 8.905806, 8.276875, 8.273174 | 8.276875 | 7.6 | 3 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=10 | 3.776185, 4.056180, 4.064619 | 4.056180 | 7.1 | 3 |
| yugabyte | single | family=absence;stmt=not_in_notnull;fanout=100 | 3.919785, 3.579715, 3.867141 | 3.867141 | 8.8 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=1 | 9.020588, 9.042504, 8.685197 | 9.020588 | 4.0 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=10 | 3.963112, 4.376581, 4.167820 | 4.167820 | 9.9 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=100 | 3.682824, 3.176320, 3.515803 | 3.515803 | 14.4 | 3 |
| yugabyte | single | family=membership;stmt=exists;fanout=1;position=last | 8.997854 | 8.997854 | 0.0 | 1 |
| yugabyte | single | family=membership;stmt=in;fanout=1 | 8.916762, 8.986667, 8.945586 | 8.945586 | 0.8 | 3 |
| yugabyte | single | family=membership;stmt=in;fanout=10 | 4.038378, 4.354579, 4.332868 | 4.332868 | 7.3 | 3 |
| yugabyte | single | family=membership;stmt=in;fanout=100 | 3.509264, 3.141402, 3.570068 | 3.509264 | 12.2 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=1 | 9.478730, 8.858875, 8.685973 | 8.858875 | 8.9 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=10 | 4.416506, 4.425877, 4.294224 | 4.416506 | 3.0 | 3 |
| yugabyte | single | family=membership;stmt=in_notnull;fanout=100 | 3.773528, 3.773027, 3.919756 | 3.773528 | 3.9 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=1 | 9.362084, 9.434328, 9.179738 | 9.362084 | 2.7 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=10 | 3.783003, 4.025412, 3.920944 | 3.920944 | 6.2 | 3 |
| yugabyte | single | family=membership;stmt=join_distinct;fanout=100 | 3.607218, 3.374191, 3.588145 | 3.588145 | 6.5 | 3 |

## Membership pair ratios per run (per engine)

| engine | topology | fanout | pair | ratios (B/A) | signs |
|---|---|---|---|---|---|
| postgres | single | 1 | in vs exists | 0.9093 | <1 |
| postgres | single | 1 | in_notnull vs exists | 1.0119 | >1 |
| postgres | single | 1 | join_distinct vs in | 1.6048 | >1 |
| postgres | single | 1 | join_distinct vs exists | 1.4592 | >1 |
| postgres | single | 10 | in vs exists | 0.9687 | <1 |
| postgres | single | 10 | in_notnull vs exists | 1.0087 | >1 |
| postgres | single | 10 | join_distinct vs in | 1.3938 | >1 |
| postgres | single | 10 | join_distinct vs exists | 1.3502 | >1 |
| postgres | single | 100 | in vs exists | 1.1320 | >1 |
| postgres | single | 100 | in_notnull vs exists | 1.1714 | >1 |
| postgres | single | 100 | join_distinct vs in | 3.8842 | >1 |
| postgres | single | 100 | join_distinct vs exists | 4.3970 | >1 |
| postgres | single | 1 | in vs exists | 0.8577 | <1 |
| postgres | single | 1 | in_notnull vs exists | 1.0260 | >1 |
| postgres | single | 1 | join_distinct vs in | 1.6834 | >1 |
| postgres | single | 1 | join_distinct vs exists | 1.4439 | >1 |
| postgres | single | 10 | in vs exists | 0.9975 | <1 |
| postgres | single | 10 | in_notnull vs exists | 0.9705 | <1 |
| postgres | single | 10 | join_distinct vs in | 1.3001 | >1 |
| postgres | single | 10 | join_distinct vs exists | 1.2969 | >1 |
| postgres | single | 100 | in vs exists | 1.1236 | >1 |
| postgres | single | 100 | in_notnull vs exists | 1.1199 | >1 |
| postgres | single | 100 | join_distinct vs in | 3.5506 | >1 |
| postgres | single | 100 | join_distinct vs exists | 3.9893 | >1 |
| postgres | single | 1 | in vs exists | 1.0505 | >1 |
| postgres | single | 1 | in_notnull vs exists | 0.9961 | <1 |
| postgres | single | 1 | join_distinct vs in | 1.4004 | >1 |
| postgres | single | 1 | join_distinct vs exists | 1.4711 | >1 |
| postgres | single | 10 | in vs exists | 0.9601 | <1 |
| postgres | single | 10 | in_notnull vs exists | 0.9419 | <1 |
| postgres | single | 10 | join_distinct vs in | 1.6313 | >1 |
| postgres | single | 10 | join_distinct vs exists | 1.5661 | >1 |
| postgres | single | 100 | in vs exists | 0.9018 | <1 |
| postgres | single | 100 | in_notnull vs exists | 0.8968 | <1 |
| postgres | single | 100 | join_distinct vs in | 4.0051 | >1 |
| postgres | single | 100 | join_distinct vs exists | 3.6120 | >1 |
| yugabyte | single | 1 | in vs exists | 0.9885 | <1 |
| yugabyte | single | 1 | in_notnull vs exists | 1.0508 | >1 |
| yugabyte | single | 1 | join_distinct vs in | 1.0499 | >1 |
| yugabyte | single | 1 | join_distinct vs exists | 1.0379 | >1 |
| yugabyte | single | 10 | in vs exists | 1.0190 | >1 |
| yugabyte | single | 10 | in_notnull vs exists | 1.1144 | >1 |
| yugabyte | single | 10 | join_distinct vs in | 0.9368 | <1 |
| yugabyte | single | 10 | join_distinct vs exists | 0.9546 | <1 |
| yugabyte | single | 100 | in vs exists | 0.9529 | <1 |
| yugabyte | single | 100 | in_notnull vs exists | 1.0246 | >1 |
| yugabyte | single | 100 | join_distinct vs in | 1.0279 | >1 |
| yugabyte | single | 100 | join_distinct vs exists | 0.9795 | <1 |
| yugabyte | single | 1 | in vs exists | 0.9938 | <1 |
| yugabyte | single | 1 | in_notnull vs exists | 0.9797 | <1 |
| yugabyte | single | 1 | join_distinct vs in | 1.0498 | >1 |
| yugabyte | single | 1 | join_distinct vs exists | 1.0433 | >1 |
| yugabyte | single | 10 | in vs exists | 0.9950 | <1 |
| yugabyte | single | 10 | in_notnull vs exists | 1.0113 | >1 |
| yugabyte | single | 10 | join_distinct vs in | 0.9244 | <1 |
| yugabyte | single | 10 | join_distinct vs exists | 0.9198 | <1 |
| yugabyte | single | 100 | in vs exists | 0.9890 | <1 |
| yugabyte | single | 100 | in_notnull vs exists | 1.1879 | >1 |
| yugabyte | single | 100 | join_distinct vs in | 1.0741 | >1 |
| yugabyte | single | 100 | join_distinct vs exists | 1.0623 | >1 |
| yugabyte | single | 1 | in vs exists | 1.0300 | >1 |
| yugabyte | single | 1 | in_notnull vs exists | 1.0001 | >1 |
| yugabyte | single | 1 | join_distinct vs in | 1.0262 | >1 |
| yugabyte | single | 1 | join_distinct vs exists | 1.0569 | >1 |
| yugabyte | single | 10 | in vs exists | 1.0396 | >1 |
| yugabyte | single | 10 | in_notnull vs exists | 1.0303 | >1 |
| yugabyte | single | 10 | join_distinct vs in | 0.9049 | <1 |
| yugabyte | single | 10 | join_distinct vs exists | 0.9408 | <1 |
| yugabyte | single | 100 | in vs exists | 1.0154 | >1 |
| yugabyte | single | 100 | in_notnull vs exists | 1.1149 | >1 |
| yugabyte | single | 100 | join_distinct vs in | 1.0051 | >1 |
| yugabyte | single | 100 | join_distinct vs exists | 1.0206 | >1 |

## Recorded limits

- [select-exist-pg-1b-20261004] one sequential reader, no concurrent writers; latency of a drained answer, not throughput
- [select-exist-pg-1b-20261004] single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host
- [select-exist-pg-1b-20261004] qualifying proportion, parent hit fraction and hot-parent skew are not swept in this bounded run
- [select-exist-pg-1b-20261004] YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained
