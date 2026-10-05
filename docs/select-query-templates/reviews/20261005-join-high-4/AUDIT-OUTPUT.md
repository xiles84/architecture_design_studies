# Audit output — fourth HIGH join/CTE review (v4 topology runs)

## 1. Every v4 run: validity, samples, controls, build, digest

| run | topology/budget/phase/placement | valid/cells | 5-sample cells | controls fired | digest in join.md |
|---|---|---|---|---|---|
| colocated-b1 | cluster3/2 CPU 3g/blended/colocated | 30/30 | 30 | 2/2 | true |
| colocated-b2 | cluster3/2 CPU 3g/blended/colocated | 30/30 | 30 | 2/2 | true |
| colocated-b3 | cluster3/2 CPU 3g/blended/colocated | 30/30 | 30 | 2/2 | true |
| colocated-i1 | cluster3/2 CPU 3g/isolated/colocated | 30/30 | 30 | 2/2 | true |
| colocated-i2 | cluster3/2 CPU 3g/isolated/colocated | 30/30 | 30 | 2/2 | true |
| colocated-i3 | cluster3/2 CPU 3g/isolated/colocated | 30/30 | 30 | 2/2 | true |
| distributed-b1 | cluster3/2 CPU 3g/blended/distributed | 30/30 | 30 | 2/2 | true |
| distributed-b2 | cluster3/2 CPU 3g/blended/distributed | 30/30 | 30 | 2/2 | true |
| distributed-b3 | cluster3/2 CPU 3g/blended/distributed | 30/30 | 30 | 2/2 | true |
| distributed-i1 | cluster3/2 CPU 3g/isolated/distributed | 30/30 | 30 | 2/2 | true |
| distributed-i2 | cluster3/2 CPU 3g/isolated/distributed | 30/30 | 30 | 2/2 | true |
| distributed-i3 | cluster3/2 CPU 3g/isolated/distributed | 30/30 | 30 | 2/2 | true |
| equal-b1 | single/6 CPU 9g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| equal-b2 | single/6 CPU 9g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| equal-b3 | single/6 CPU 9g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| equal-i1 | single/6 CPU 9g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| equal-i2 | single/6 CPU 9g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| equal-i3 | single/6 CPU 9g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-b1 | single/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-b2 | single/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-b3 | single/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-i1 | single/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-i2 | single/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| single-i3 | single/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-b1 | cluster3/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-b2 | cluster3/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-b3 | cluster3/2 CPU 3g/blended/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-i1 | cluster3/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-i2 | cluster3/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |
| tablet-i3 | cluster3/2 CPU 3g/isolated/noncolocated-single-tablet | 30/30 | 30 | 2/2 | true |

Runs: 30. Distinct client image ids: 1. Distinct repo commits: 1.

## 2. Placement evidence (cluster runs)

| run | tablets (yb_table_properties) | leaders after run (placement.txt) | donation leaders before | gate files | gate note |
|---|---|---|---|---|---|
| colocated-b1 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| colocated-b2 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| colocated-b3 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| colocated-i1 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 |  |  | false |
| colocated-i2 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| colocated-i3 | charity=1c person=1c donation=1c donation_rank_idx=1c donation_person_idx=1c | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 |  |  | false |
| distributed-b1 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| distributed-b2 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| distributed-b3 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| distributed-i1 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| distributed-i2 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| distributed-i3 | charity=3 person=3 donation=3 donation_rank_idx=3 donation_person_idx=3 | charity:yb-n1+yb-n2+yb-n3 donation:yb-n1+yb-n2+yb-n3 donation_person_idx:yb-n1+yb-n2+yb-n3 donation_rank_idx:yb-n1+yb-n2+yb-n3 person:yb-n1+yb-n2+yb-n3 | yb-n1+yb-n2+yb-n3 | prepared ok | true |
| tablet-b1 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 |  |  | false |
| tablet-b2 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| tablet-b3 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 |  |  | false |
| tablet-i1 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| tablet-i2 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 |  |  | false |
| tablet-i3 | charity=1 person=1 donation=1 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 |  |  | false |

## 3i. join_where per-run medians (ms), phase isolated

| arm | t=20 | t=2000 | t=180000 | joined MAT/default t=2000 | MAT t=2000 | MAT t=180000 |
|---|---|---|---|---|---|---|
| single | 2.23, 1.71, 1.70 | 9.49, 9.74, 8.38 | 108.8, 126.8, 104.4 | 9.55, 8.31, 10.38 | 89.0, 91.4, 85.7 | 109.3, 113.0, 105.1 |
| equal | 1.69, 1.62, 1.83 | 7.83, 8.50, 7.85 | 99.4, 101.1, 100.9 | 8.12, 9.34, 9.51 | 77.3, 78.9, 78.5 | 96.1, 94.3, 96.9 |
| tablet | 1.87, 2.75, 2.61 | 9.29, 10.94, 11.51 | 142.8, 146.7, 107.7 | 11.75, 13.72, 11.98 | 125.5, 122.4, 124.2 | 134.1, 105.2, 137.7 |
| distributed | 3.31, 2.87, 3.18 | 7.84, 8.02, 8.10 | 75.5, 79.1, 77.1 | 8.50, 8.53, 8.31 | 65.1, 68.0, 65.3 | 101.6, 103.8, 100.9 |
| colocated | 2.26, 2.39, 1.61 | 6.28, 6.57, 5.96 | 89.4, 126.1, 92.9 | 16.80, 18.01, 11.54 | 101.2, 105.8, 64.9 | 118.4, 121.6, 122.1 |

Across-run medians (median of the three run medians) and ratios to the single node:

| arm | t=20 | t=2000 | t=180000 | /single t=20 | /single t=2000 | /single t=180000 | MAT t=180000 /single | MAT t=2000 /single |
|---|---|---|---|---|---|---|---|---|
| single | 1.71 | 9.49 | 108.84 | 1.00 | 1.00 | 1.00 | 1.00 | 1.00 |
| equal | 1.69 | 7.85 | 100.92 | 0.99 | 0.83 | 0.93 | 0.88 | 0.88 |
| tablet | 2.61 | 10.94 | 142.83 | 1.53 | 1.15 | 1.31 | 1.23 | 1.40 |
| distributed | 3.18 | 8.02 | 77.05 | 1.86 | 0.85 | 0.71 | 0.93 | 0.73 |
| colocated | 2.26 | 6.28 | 92.93 | 1.33 | 0.66 | 0.85 | 1.11 | 1.14 |

Run-range separation of the analysis's pairs (a vs b):

| pair (a vs b) | t=20 | t=2000 | t=180000 | MAT/default t=2000 |
|---|---|---|---|---|
| distributed vs tablet | a>b all runs | a<b all runs | a<b all runs | a<b all runs |
| colocated vs distributed | a<b all runs | a<b all runs | a>b all runs | a>b all runs |
| distributed vs single | a>b all runs | a<b all runs | a<b all runs | overlap |
| colocated vs single | overlap | a<b all runs | overlap | a>b all runs |
| tablet vs single | overlap | overlap | overlap | a>b all runs |
| equal vs single | overlap | overlap | a<b all runs | overlap |

## 3b. join_where per-run medians (ms), phase blended

| arm | t=20 | t=2000 | t=180000 | joined MAT/default t=2000 | MAT t=2000 | MAT t=180000 |
|---|---|---|---|---|---|---|
| single | 1.84, 1.90, 2.20 | 8.93, 9.82, 8.69 | 106.1, 110.5, 111.7 | 9.79, 9.25, 9.55 | 82.0, 87.5, 85.2 | 99.7, 125.9, 106.2 |
| equal | 2.35, 2.19, 1.82 | 10.59, 10.39, 8.74 | 119.3, 114.3, 109.3 | 9.57, 9.43, 8.66 | 94.7, 99.5, 80.0 | 119.3, 111.5, 103.5 |
| tablet | 2.09, 3.25, 2.97 | 12.30, 11.31, 11.78 | 117.5, 121.9, 116.5 | 8.08, 10.41, 12.76 | 91.5, 136.9, 122.3 | 114.6, 150.2, 105.9 |
| distributed | 3.30, 3.14, 2.76 | 8.33, 8.35, 8.24 | 87.7, 83.1, 80.2 | 9.15, 8.35, 7.78 | 78.3, 69.8, 69.4 | 117.6, 111.9, 114.4 |
| colocated | 2.40, 2.21, 2.34 | 7.49, 6.86, 6.74 | 97.4, 128.1, 140.3 | 9.25, 11.75, 16.10 | 66.9, 70.2, 111.1 | 99.1, 129.2, 130.7 |

Across-run medians (median of the three run medians) and ratios to the single node:

| arm | t=20 | t=2000 | t=180000 | /single t=20 | /single t=2000 | /single t=180000 | MAT t=180000 /single | MAT t=2000 /single |
|---|---|---|---|---|---|---|---|---|
| single | 1.90 | 8.93 | 110.47 | 1.00 | 1.00 | 1.00 | 1.00 | 1.00 |
| equal | 2.19 | 10.39 | 114.32 | 1.16 | 1.16 | 1.03 | 1.05 | 1.11 |
| tablet | 2.97 | 11.78 | 117.49 | 1.57 | 1.32 | 1.06 | 1.08 | 1.43 |
| distributed | 3.14 | 8.33 | 83.10 | 1.66 | 0.93 | 0.75 | 1.08 | 0.82 |
| colocated | 2.34 | 6.86 | 128.08 | 1.24 | 0.77 | 1.16 | 1.22 | 0.82 |

Run-range separation of the analysis's pairs (a vs b):

| pair (a vs b) | t=20 | t=2000 | t=180000 | MAT/default t=2000 |
|---|---|---|---|---|
| distributed vs tablet | overlap | a<b all runs | a<b all runs | overlap |
| colocated vs distributed | a<b all runs | a<b all runs | a>b all runs | a>b all runs |
| distributed vs single | a>b all runs | a<b all runs | a<b all runs | a<b all runs |
| colocated vs single | a>b all runs | a<b all runs | overlap | overlap |
| tablet vs single | overlap | a>b all runs | a>b all runs | overlap |
| equal vs single | overlap | overlap | overlap | overlap |

## 4. Blended phase: endpoints, writer and ledger

| run | reader endpoint ops | writer offered/completed/failed | writer endpoint ops | ledger mismatch |
|---|---|---|---|---|
| colocated-b1 | [168 168 168] | 122/122/0 | [41 41 40] | 0 |
| colocated-b2 | [168 168 168] | 113/113/0 | [38 38 37] | 0 |
| colocated-b3 | [168 168 168] | 125/125/0 | [42 42 41] | 0 |
| distributed-b1 | [169 168 168] | 169/169/0 | [57 56 56] | 0 |
| distributed-b2 | [169 169 168] | 162/162/0 | [54 54 54] | 0 |
| distributed-b3 | [169 169 168] | 152/152/0 | [51 51 50] | 0 |
| equal-b1 | [504] | 155/155/0 | [155] | 0 |
| equal-b2 | [505] | 157/157/0 | [157] | 0 |
| equal-b3 | [504] | 137/137/0 | [137] | 0 |
| single-b1 | [504] | 146/146/0 | [146] | 0 |
| single-b2 | [504] | 155/155/0 | [155] | 0 |
| single-b3 | [504] | 152/152/0 | [152] | 0 |
| tablet-b1 | [168 168 168] | 167/167/0 | [56 56 55] | 0 |
| tablet-b2 | [168 168 168] | 175/175/0 | [59 58 58] | 0 |
| tablet-b3 | [168 168 168] | 162/162/0 | [54 54 54] | 0 |

## 5. Node CPU throttling over each run (throttled periods / periods; throttled seconds)

| run | per node |
|---|---|
| colocated-b1 | yb-n1 5/514 0.74s; yb-n2 5/507 0.58s; yb-n3 5/500 0.11s |
| colocated-b2 | yb-n1 5/498 0.58s; yb-n2 5/487 0.35s; yb-n3 2/478 0.02s |
| colocated-b3 | yb-n1 6/517 0.73s; yb-n2 5/509 0.61s; yb-n3 1/507 0.04s |
| colocated-i1 | yb-n1 4/495 0.54s; yb-n2 8/484 1.14s; yb-n3 0/483 0.00s |
| colocated-i2 | yb-n1 7/514 1.17s; yb-n2 6/503 0.83s; yb-n3 5/502 0.14s |
| colocated-i3 | yb-n1 5/494 0.64s; yb-n2 10/486 1.57s; yb-n3 0/480 0.00s |
| distributed-b1 | yb-n1 26/697 1.44s; yb-n2 6/691 1.14s; yb-n3 5/684 0.13s |
| distributed-b2 | yb-n1 9/675 1.04s; yb-n2 5/666 0.88s; yb-n3 3/661 0.14s |
| distributed-b3 | yb-n1 6/652 0.79s; yb-n2 6/644 0.90s; yb-n3 0/640 0.00s |
| distributed-i1 | yb-n1 10/594 0.90s; yb-n2 4/587 0.66s; yb-n3 1/585 0.00s |
| distributed-i2 | yb-n1 15/638 1.30s; yb-n2 6/633 0.86s; yb-n3 0/623 0.00s |
| distributed-i3 | yb-n1 4/618 0.65s; yb-n2 7/610 0.93s; yb-n3 5/608 0.26s |
| equal-b1 | yb-single 0/384 0.00s |
| equal-b2 | yb-single 0/379 0.00s |
| equal-b3 | yb-single 0/345 0.00s |
| equal-i1 | yb-single 0/335 0.00s |
| equal-i2 | yb-single 0/337 0.00s |
| equal-i3 | yb-single 0/337 0.00s |
| single-b1 | yb-single 71/368 5.31s |
| single-b2 | yb-single 92/394 5.19s |
| single-b3 | yb-single 83/383 5.82s |
| single-i1 | yb-single 82/380 5.21s |
| single-i2 | yb-single 79/392 6.38s |
| single-i3 | yb-single 79/380 5.67s |
| tablet-b1 | yb-n1 4/602 0.53s; yb-n2 13/593 2.06s; yb-n3 0/588 0.00s |
| tablet-b2 | yb-n1 6/621 0.56s; yb-n2 5/616 0.71s; yb-n3 10/608 0.33s |
| tablet-b3 | yb-n1 5/599 0.50s; yb-n2 17/592 1.56s; yb-n3 0/586 0.00s |
| tablet-i1 | yb-n1 4/602 0.58s; yb-n2 6/593 0.76s; yb-n3 5/590 0.05s |
| tablet-i2 | yb-n1 4/586 0.43s; yb-n2 5/577 0.63s; yb-n3 10/572 0.23s |
| tablet-i3 | yb-n1 4/599 0.52s; yb-n2 17/595 1.67s; yb-n3 0/588 0.00s |

## 6. P4 reuse ratios (NOT MATERIALIZED / MATERIALIZED) in the v3b domain runs, per run

| run | consumer | n=1 | n=2 | n=4 | n=8 | n=16 |
|---|---|---|---|---|---|---|
| pg-b1 | 1% | 0.186 | 0.297 | 0.423 | 0.489 | 0.687 |
| pg-b1 | 100% | 0.988 | 1.748 | 2.663 | 2.919 | 4.598 |
| pg-b2 | 1% | 0.217 | 0.303 | 0.312 | 0.488 | 0.752 |
| pg-b2 | 100% | 0.868 | 1.904 | 2.770 | 3.222 | 4.108 |
| pg-b3 | 1% | 0.270 | 0.266 | 0.329 | 0.485 | 0.746 |
| pg-b3 | 100% | 0.893 | 1.874 | 2.485 | 3.080 | 4.500 |
| pg-i1 | 1% | 0.171 | 0.209 | 0.360 | 0.502 | 0.712 |
| pg-i1 | 100% | 1.053 | 1.926 | 2.742 | 3.267 | 3.958 |
| pg-i2 | 1% | 0.178 | 0.256 | 0.424 | 0.602 | 0.806 |
| pg-i2 | 100% | 0.960 | 2.059 | 2.494 | 3.262 | 3.891 |
| pg-i3 | 1% | 0.188 | 0.293 | 0.358 | 0.505 | 0.683 |
| pg-i3 | 100% | 0.988 | 1.703 | 2.514 | 3.275 | 3.938 |
| yb-b1 | 1% | 0.916 | 1.435 | 2.449 | 9.794 | 10.049 |
| yb-b1 | 100% | 0.941 | 1.796 | 3.489 | 5.722 | 10.962 |
| yb-b2 | 1% | 0.880 | 1.522 | 2.694 | 7.051 | 9.178 |
| yb-b2 | 100% | 0.899 | 1.974 | 3.196 | 6.339 | 11.132 |
| yb-b3 | 1% | 0.944 | 1.449 | 2.538 | 6.881 | 9.058 |
| yb-b3 | 100% | 1.028 | 1.680 | 3.398 | 6.044 | 9.486 |
| yb-i1 | 1% | 0.812 | 1.502 | 2.929 | 6.186 | 11.294 |
| yb-i1 | 100% | 0.967 | 1.856 | 3.686 | 6.733 | 8.990 |
| yb-i2 | 1% | 0.874 | 1.460 | 2.776 | 6.388 | 11.323 |
| yb-i2 | 100% | 1.066 | 1.850 | 3.573 | 6.846 | 11.806 |
| yb-i3 | 1% | 0.824 | 1.430 | 2.795 | 6.341 | 7.451 |
| yb-i3 | 100% | 0.979 | 1.925 | 3.511 | 6.808 | 12.139 |
