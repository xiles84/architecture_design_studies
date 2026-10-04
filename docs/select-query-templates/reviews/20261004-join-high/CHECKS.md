# Independent saved-input audit — HIGH GPT-6 / Codex desktop

No database starts or new timings. Hashes normalize checkout CRLF to committed LF.

| run | input SHA-256 | cells | controls | samples >20% spread | last/first | actual guest memory kB |
|---|---|---|---|---|---|---|
| select-join-pg1-20261004 | `070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac` | 91 | 2/2 | 25/91 | 0.8006 | 16109780 |
| select-join-pg2-20261004 | `112577024534a9e91932f78d3ba8b9e530e0caea0aeefdf851737e578be7bab6` | 91 | 2/2 | 38/91 | 0.8963 | 16109780 |
| select-join-pg3-20261004 | `05b2d7a99e38c77c84baf413a769f531179682921a2419bf73da89153b2d00ac` | 91 | 2/2 | 88/91 | 0.6056 | 16109780 |
| select-join-yb1-20261004 | `f650acde1ebb5541080360bca3e467863611e00322d3f580cefdd36ba90e0281` | 91 | 2/2 | 46/91 | 0.8609 | 16109780 |
| select-join-yb2-20261004 | `25e5f7bf01de5149bed501f089eebe900d7b462c2927eedaee302436b664d198` | 91 | 2/2 | 46/91 | 1.3470 | 16109780 |
| select-join-yb3-20261004 | `2d8394b9cc2a9a0f649367486f9d994fcfc1538e3e7d151bc75a1e8febe82e46` | 91 | 2/2 | 8/91 | 1.0102 | 16109780 |

## Materialized joined-output / default joined-output (paired per-run ratios)

| engine | threshold | ratios in run order |
|---|---|---|
| postgres | 20 | 73.1056, 92.8513, 85.3413 |
| postgres | 200 | 73.9123, 47.6644, 56.8073 |
| postgres | 2000 | 26.6177, 30.6338, 27.4406 |
| postgres | 20000 | 6.7868, 7.3624, 7.1675 |
| postgres | 100000 | 1.4197, 1.4576, 1.2333 |
| postgres | 180000 | 1.4145, 1.3707, 1.2678 |
| yugabyte | 20 | 1.5592, 1.5073, 1.5269 |
| yugabyte | 200 | 1.5596, 1.5325, 1.5315 |
| yugabyte | 2000 | 1.4476, 1.6317, 1.4778 |
| yugabyte | 20000 | 1.3118, 1.3803, 1.3084 |
| yugabyte | 100000 | 1.0577, 1.0497, 1.0845 |
| yugabyte | 180000 | 1.0172, 0.8844, 0.9442 |

## NOT MATERIALIZED / MATERIALIZED (actual P4 pair, not /default)

| engine | references | consumer id-space % | ratios in run order |
|---|---|---|---|
| postgres | 1 | 1 | 0.2066, 0.2139, 0.1948 |
| postgres | 1 | 100 | 0.9338, 0.9687, 0.9080 |
| postgres | 2 | 1 | 0.2716, 0.2173, 0.2013 |
| postgres | 2 | 100 | 1.8339, 2.1235, 1.6957 |
| postgres | 4 | 1 | 0.4197, 0.3777, 0.4279 |
| postgres | 4 | 100 | 2.8919, 2.5650, 1.9924 |
| postgres | 8 | 1 | 0.5387, 0.5406, 0.5644 |
| postgres | 8 | 100 | 2.7486, 3.1994, 4.0087 |
| postgres | 16 | 1 | 0.6950, 0.7868, 0.7798 |
| postgres | 16 | 100 | 4.2011, 4.0048, 4.3493 |
| yugabyte | 1 | 1 | 1.1097, 1.0432, 1.0389 |
| yugabyte | 1 | 100 | 1.0184, 0.9955, 1.0421 |
| yugabyte | 2 | 1 | 1.7076, 1.7987, 1.7214 |
| yugabyte | 2 | 100 | 1.7674, 1.6969, 1.7703 |
| yugabyte | 4 | 1 | 3.0470, 2.7194, 2.4946 |
| yugabyte | 4 | 100 | 3.0861, 2.6638, 2.3941 |
| yugabyte | 8 | 1 | 5.0116, 4.9979, 4.6362 |
| yugabyte | 8 | 100 | 5.1479, 5.0772, 5.4213 |
| yugabyte | 16 | 1 | 9.0169, 9.3047, 9.9375 |
| yugabyte | 16 | 100 | 10.3730, 10.1194, 9.5438 |

## Broad default-spelling pair ratios to join+WHERE

| engine | threshold | overall min | overall max |
|---|---|---|---|
| postgres | 100000 | 0.9035 | 1.1113 |
| postgres | 180000 | 0.8595 | 1.0810 |
| yugabyte | 100000 | 0.9632 | 1.1138 |
| yugabyte | 180000 | 0.8476 | 1.0471 |

All saved-input arithmetic, original input digests, seeded orders, forced materialization and control checks passed. This does not accept absent workload/topology coverage or the author's interpretation.
