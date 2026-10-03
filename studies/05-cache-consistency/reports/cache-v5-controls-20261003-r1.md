# Cache publication and deletion measurements

## TL;DR

Fixed selection rule: count every result and control; publish timing only for a cell whose complete correctness gate passed. 10/10 cells passed; 57 controls passed. Inputs digest: `02900bf8b0429be4` (32 JSON files). This report contains measured facts; interpretation belongs in a signed analysis.

Environment: `xiles-podman-4cpu-2g-config-protocol-20261003`. Producing commit: `8a6e7f0943d2eff975eba05c048971c8ae55da6a`; run tag: `run/study-05/cache-v5-controls-20261003-r1`; client image: `9f2d2eccd1aa903db28d5fbe8bd95133f78ef363ae777e37fa6dcd9daf3d8d59`.

Single-node PostgreSQL and Redis share this laptop's Podman VM. Service budget: PostgreSQL 1 CPU/512 MiB plus Redis 1 CPU/256 MiB; client 1 CPU/512 MiB; planned reserve 1 CPU/768 MiB against machine configuration 4 CPU/2 GiB. The guest's actual resources and service cgroup counters are retained beside each cell. All source requests use one endpoint. Three-node/equal-total cluster and multi-host placement comparisons are resource/coverage gaps. Colocation, rollup, rolldown and embedding have no controlled contrast in this keyed projection. These are short, closed-loop measurements on 128 keys, with per-key client writer ordering, not saturation, donor-portal throughput or uncoordinated-writer evidence.

## Every admitted trial

| Arm / protocol / strategy / deletion | Phase / trial | Read/s | Write/s | Read p50/p99 ms | Write p50/p99 ms | DB p50 ms | Ack p50/p99 ms | Client lock p99 ms | Source queries/read | Relay queries / delivered |
|---|---|---:|---:|---|---|---:|---|---:|---:|---|
| A / version / cas / hard | read / 1 | 7427.2 | 0.0 | 0.350 / 14.425 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 1 | 0.0 | 1145.2 | 0.000 / 0.000 | 1.654 / 4.684 | 1.645 | 0.002 / 0.010 | 0.000 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 1 | 5462.9 | 236.4 | 0.477 / 52.407 | 2.812 / 69.000 | 2.808 | 0.001 / 0.031 | 0.000 | 1.042928 | 0 / 0 |
| A / version / cas / hard | read / 2 | 7150.0 | 0.0 | 0.389 / 62.720 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 2 | 0.0 | 981.9 | 0.000 / 0.000 | 1.985 / 3.457 | 1.974 | 0.002 / 0.008 | 0.000 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 2 | 5684.0 | 217.2 | 0.436 / 62.329 | 2.594 / 71.994 | 2.589 | 0.001 / 0.028 | 0.000 | 1.036115 | 0 / 0 |
| A / version / cas / hard | read / 3 | 7430.1 | 0.0 | 0.354 / 14.141 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 3 | 0.0 | 1143.4 | 0.000 / 0.000 | 1.675 / 2.829 | 1.669 | 0.002 / 0.005 | 0.000 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 3 | 5602.3 | 296.2 | 0.484 / 62.576 | 2.532 / 66.698 | 2.529 | 0.001 / 0.047 | 0.000 | 1.052868 | 0 / 0 |
| A / sync / cas / hard | read / 1 | 15246.9 | 0.0 | 0.178 / 0.842 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | write / 1 | 0.0 | 937.7 | 0.000 / 0.000 | 2.109 / 2.913 | 1.735 | 0.186 / 0.365 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | mixed / 1 | 12299.3 | 199.7 | 0.190 / 1.688 | 2.908 / 68.760 | 2.373 | 0.272 / 0.625 | 0.000 | 0.028827 | 0 / 0 |
| A / sync / cas / hard | read / 2 | 14945.5 | 0.0 | 0.177 / 0.813 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | write / 2 | 0.0 | 895.7 | 0.000 / 0.000 | 2.192 / 2.961 | 1.795 | 0.199 / 0.459 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | mixed / 2 | 9712.6 | 231.9 | 0.253 / 1.948 | 3.181 / 68.586 | 2.535 | 0.318 / 0.851 | 0.000 | 0.037656 | 0 / 0 |
| A / sync / cas / hard | read / 3 | 16299.5 | 0.0 | 0.162 / 0.530 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | write / 3 | 0.0 | 898.8 | 0.000 / 0.000 | 2.168 / 3.559 | 1.792 | 0.185 / 0.372 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / cas / hard | mixed / 3 | 12636.8 | 288.7 | 0.208 / 1.440 | 2.609 / 70.288 | 2.030 | 0.231 / 0.685 | 0.000 | 0.035455 | 0 / 0 |
| A / outbox / lock / hard | read / 1 | 15840.2 | 0.0 | 0.162 / 0.583 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 69 / 0 |
| A / outbox / lock / hard | write / 1 | 0.0 | 336.6 | 0.000 / 0.000 | 4.797 / 13.463 | 2.279 | 2.736 / 10.341 | 0.000 | 0.000000 | 108 / 171 |
| A / outbox / lock / hard | mixed / 1 | 15277.1 | 105.5 | 0.177 / 0.950 | 7.064 / 76.668 | 2.596 | 3.881 / 71.665 | 0.000 | 0.007684 | 48 / 53 |
| A / outbox / lock / hard | read / 2 | 16785.0 | 0.0 | 0.163 / 0.490 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 69 / 0 |
| A / outbox / lock / hard | write / 2 | 0.0 | 325.8 | 0.000 / 0.000 | 5.085 / 9.953 | 2.308 | 2.775 / 7.068 | 0.012 | 0.000000 | 111 / 164 |
| A / outbox / lock / hard | mixed / 2 | 11281.2 | 95.3 | 0.228 / 1.541 | 8.767 / 72.498 | 2.928 | 5.128 / 69.303 | 0.000 | 0.009243 | 44 / 53 |
| A / outbox / lock / hard | read / 3 | 13480.8 | 0.0 | 0.204 / 0.942 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 73 / 0 |
| A / outbox / lock / hard | write / 3 | 0.0 | 343.5 | 0.000 / 0.000 | 4.910 / 11.058 | 2.157 | 2.845 / 7.427 | 0.000 | 0.000000 | 112 / 172 |
| A / outbox / lock / hard | mixed / 3 | 12644.5 | 106.4 | 0.200 / 1.247 | 7.772 / 74.384 | 2.811 | 4.395 / 69.789 | 0.000 | 0.009416 | 49 / 59 |
| A / sync / lock / hard | read / 1 | 12082.5 | 0.0 | 0.209 / 5.567 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | write / 1 | 0.0 | 296.0 | 0.000 / 0.000 | 4.950 / 27.804 | 3.882 | 0.324 / 4.940 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | mixed / 1 | 9947.2 | 150.5 | 0.224 / 6.061 | 6.118 / 72.592 | 4.997 | 0.406 / 2.143 | 0.018 | 0.027712 | 0 / 0 |
| A / sync / lock / hard | read / 2 | 14899.2 | 0.0 | 0.186 / 0.847 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | write / 2 | 0.0 | 624.3 | 0.000 / 0.000 | 2.755 / 10.182 | 2.322 | 0.225 / 1.276 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | mixed / 2 | 11414.0 | 204.6 | 0.236 / 3.768 | 4.784 / 59.200 | 3.780 | 0.352 / 2.055 | 0.000 | 0.025936 | 0 / 0 |
| A / sync / lock / hard | read / 3 | 15921.5 | 0.0 | 0.171 / 2.448 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | write / 3 | 0.0 | 468.3 | 0.000 / 0.000 | 3.618 / 11.849 | 2.965 | 0.287 / 1.767 | 0.000 | 0.000000 | 0 / 0 |
| A / sync / lock / hard | mixed / 3 | 10940.9 | 202.5 | 0.249 / 3.275 | 4.456 / 59.389 | 3.558 | 0.327 / 1.897 | 0.000 | 0.033212 | 0 / 0 |
| A / version / lock / hard | read / 1 | 6615.7 | 0.0 | 0.425 / 47.248 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / lock / hard | write / 1 | 0.0 | 780.8 | 0.000 / 0.000 | 2.347 / 5.449 | 2.339 | 0.002 / 0.021 | 0.000 | 0.000000 | 0 / 0 |
| A / version / lock / hard | mixed / 1 | 4428.8 | 125.3 | 0.526 / 61.225 | 4.361 / 82.157 | 4.352 | 0.001 / 0.040 | 0.001 | 1.029103 | 0 / 0 |
| A / version / lock / hard | read / 2 | 6535.9 | 0.0 | 0.408 / 30.262 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / lock / hard | write / 2 | 0.0 | 870.5 | 0.000 / 0.000 | 2.187 / 4.119 | 2.174 | 0.002 / 0.005 | 0.000 | 0.000000 | 0 / 0 |
| A / version / lock / hard | mixed / 2 | 5553.9 | 226.8 | 0.492 / 57.196 | 3.248 / 64.544 | 3.206 | 0.001 / 0.008 | 0.000 | 1.039129 | 0 / 0 |
| A / version / lock / hard | read / 3 | 6591.7 | 0.0 | 0.394 / 60.804 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / lock / hard | write / 3 | 0.0 | 883.0 | 0.000 / 0.000 | 2.037 / 5.736 | 2.030 | 0.002 / 0.004 | 0.000 | 0.000000 | 0 / 0 |
| A / version / lock / hard | mixed / 3 | 5161.2 | 141.3 | 0.420 / 53.856 | 5.383 / 73.957 | 5.339 | 0.001 / 0.049 | 0.000 | 1.026224 | 0 / 0 |
| A / outbox / cas / hard | read / 1 | 13342.0 | 0.0 | 0.200 / 1.506 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 68 / 0 |
| A / outbox / cas / hard | write / 1 | 0.0 | 339.9 | 0.000 / 0.000 | 5.093 / 10.535 | 2.145 | 3.064 / 7.311 | 0.035 | 0.000000 | 109 / 172 |
| A / outbox / cas / hard | mixed / 1 | 13372.5 | 103.7 | 0.180 / 1.312 | 8.825 / 72.057 | 2.827 | 5.206 / 67.826 | 0.000 | 0.008693 | 48 / 58 |
| A / outbox / cas / hard | read / 2 | 13764.3 | 0.0 | 0.183 / 1.247 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 67 / 0 |
| A / outbox / cas / hard | write / 2 | 0.0 | 273.3 | 0.000 / 0.000 | 6.415 / 15.773 | 2.585 | 3.777 / 8.319 | 0.072 | 0.000000 | 98 / 138 |
| A / outbox / cas / hard | mixed / 2 | 13241.3 | 84.6 | 0.176 / 2.502 | 10.699 / 72.717 | 3.537 | 5.714 / 63.291 | 0.000 | 0.007723 | 40 / 43 |
| A / outbox / cas / hard | read / 3 | 15174.8 | 0.0 | 0.168 / 1.276 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 0.000000 | 67 / 0 |
| A / outbox / cas / hard | write / 3 | 0.0 | 304.0 | 0.000 / 0.000 | 5.403 / 12.368 | 2.182 | 3.527 / 10.310 | 0.000 | 0.000000 | 102 / 154 |
| A / outbox / cas / hard | mixed / 3 | 11873.5 | 96.2 | 0.201 / 1.710 | 8.156 / 73.837 | 2.788 | 4.694 / 69.792 | 0.005 | 0.009748 | 46 / 54 |
| A / version / cas / hard | read / 1 | 6808.4 | 0.0 | 0.363 / 53.264 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 1 | 0.0 | 1115.8 | 0.000 / 0.000 | 1.695 / 2.670 | 1.689 | 0.002 / 0.005 | 0.000 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 1 | 5374.0 | 284.7 | 0.502 / 58.592 | 2.716 / 67.344 | 2.707 | 0.001 / 0.047 | 1.119 | 1.052255 | 0 / 0 |
| A / version / cas / hard | read / 2 | 7290.4 | 0.0 | 0.362 / 24.531 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 2 | 0.0 | 1178.9 | 0.000 / 0.000 | 1.657 / 2.576 | 1.650 | 0.002 / 0.005 | 0.025 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 2 | 5397.9 | 296.5 | 0.490 / 57.577 | 2.639 / 66.823 | 2.633 | 0.001 / 0.044 | 0.000 | 1.053534 | 0 / 0 |
| A / version / cas / hard | read / 3 | 6112.1 | 0.0 | 0.442 / 64.498 | 0.000 / 0.000 | 0.000 | 0.000 / 0.000 | 0.000 | 1.000000 | 0 / 0 |
| A / version / cas / hard | write / 3 | 0.0 | 1166.1 | 0.000 / 0.000 | 1.661 / 3.150 | 1.653 | 0.002 / 0.006 | 0.000 | 0.000000 | 0 / 0 |
| A / version / cas / hard | mixed / 3 | 5695.3 | 306.8 | 0.461 / 60.740 | 2.519 / 68.513 | 2.513 | 0.001 / 0.020 | 0.000 | 1.055870 | 0 / 0 |
| B / sync / cas / soft | delete-recreate / 1 | 0.0 | 840.5 | 0.000 / 0.000 | 2.250 / 3.666 | 1.880 | 0.194 / 0.361 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / soft | delete-recreate / 2 | 0.0 | 872.0 | 0.000 / 0.000 | 2.249 / 3.242 | 1.858 | 0.199 / 0.424 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / soft | delete-recreate / 3 | 0.0 | 825.4 | 0.000 / 0.000 | 2.295 / 4.391 | 1.905 | 0.204 / 0.440 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / hard | delete-recreate / 1 | 0.0 | 898.3 | 0.000 / 0.000 | 2.171 / 3.100 | 1.818 | 0.189 / 0.348 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / hard | delete-recreate / 2 | 0.0 | 813.2 | 0.000 / 0.000 | 2.418 / 3.267 | 1.983 | 0.216 / 0.407 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / hard | delete-recreate / 3 | 0.0 | 808.5 | 0.000 / 0.000 | 2.416 / 3.210 | 1.996 | 0.219 / 0.403 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / tombstone | delete-recreate / 1 | 0.0 | 804.4 | 0.000 / 0.000 | 2.422 / 3.317 | 2.004 | 0.220 / 0.414 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / tombstone | delete-recreate / 2 | 0.0 | 787.6 | 0.000 / 0.000 | 2.484 / 3.213 | 2.066 | 0.228 / 0.501 | 0.000 | 0.000000 | 0 / 0 |
| B / sync / cas / tombstone | delete-recreate / 3 | 0.0 | 795.6 | 0.000 / 0.000 | 2.463 / 3.256 | 2.031 | 0.226 / 0.418 | 0.000 | 0.000000 | 0 / 0 |

A deletion/recreation pair counts as two acknowledged writes; write latencies refer to each individual write. Trial duration includes drained in-flight operations. Source queries/read counts snapshot and validation statements, excluding separately reported relay queries and post-trial audits. Fault-control timings below describe orchestrated pauses, not spontaneous outage distributions.

## Correctness, resource counters and fault controls

### Cell 1: A / version / cas / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 155648. Client cgroup periods/throttled: 64/34; throttled microseconds: 6365003.

- paused-old-fill: pass=true; old_publish=true served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.018..0.019 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.009..0.010 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.006..0.006 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.083 ms
Execution Time: 0.030 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.006..0.007 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.027 ms
Execution Time: 0.015 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.005..0.005 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.004..0.004 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.004..0.004 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=33 read=2
  I/O Timings: shared read=0.012
Planning Time: 0.111 ms
Execution Time: 0.019 ms
```

### Cell 2: A / sync / cas / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 131072. Client cgroup periods/throttled: 64/36; throttled microseconds: 8407590.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.025..0.027 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.012..0.013 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.010..0.010 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.099 ms
Execution Time: 0.040 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.006..0.007 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.030 ms
Execution Time: 0.015 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.004..0.004 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.003..0.003 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.003..0.003 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=33 read=2
  I/O Timings: shared read=0.016
Planning Time: 0.152 ms
Execution Time: 0.025 ms
```

### Cell 3: A / outbox / lock / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 90112. Client cgroup periods/throttled: 65/36; throttled microseconds: 8888956.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- relay-acknowledgement-barrier-and-recovery: pass=true; served generation=2; blocked/recovery/stale-duration ms=50.216/5.297/0.000.
- outbox-without-ack-barrier: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/33.559.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.018..0.019 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.011..0.012 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.005..0.005 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.089 ms
Execution Time: 0.036 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.007..0.008 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.031 ms
Execution Time: 0.017 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.003..0.004 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.002..0.002 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.002..0.002 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=28 read=2
  I/O Timings: shared read=0.025
Planning Time: 0.250 ms
Execution Time: 0.019 ms
```

### Cell 4: A / sync / lock / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 90112. Client cgroup periods/throttled: 73/39; throttled microseconds: 5907579.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.032..0.034 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.024..0.025 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.005..0.005 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.130 ms
Execution Time: 0.048 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.010..0.011 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.048 ms
Execution Time: 0.023 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.006..0.006 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.005..0.005 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.005..0.005 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=33 read=2
  I/O Timings: shared read=0.016
Planning Time: 0.161 ms
Execution Time: 0.017 ms
```

### Cell 5: A / version / lock / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 90112. Client cgroup periods/throttled: 66/44; throttled microseconds: 6038046.

- paused-old-fill: pass=true; old_publish=true served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.026..0.027 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.012..0.013 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.012..0.012 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.082 ms
Execution Time: 0.047 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.004..0.004 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.020 ms
Execution Time: 0.009 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.005..0.005 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.004..0.004 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.004..0.004 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=33 read=2
  I/O Timings: shared read=0.035
Planning Time: 0.172 ms
Execution Time: 0.019 ms
```

### Cell 6: A / outbox / cas / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 90112. Client cgroup periods/throttled: 66/36; throttled microseconds: 7089651.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- relay-acknowledgement-barrier-and-recovery: pass=true; served generation=2; blocked/recovery/stale-duration ms=50.216/5.898/0.000.
- outbox-without-ack-barrier: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/34.246.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.022..0.023 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.013..0.014 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.006..0.006 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.090 ms
Execution Time: 0.035 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.008..0.009 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.031 ms
Execution Time: 0.017 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.002..0.002 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.001..0.001 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.001..0.001 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=28 read=2
  I/O Timings: shared read=0.014
Planning Time: 0.139 ms
Execution Time: 0.012 ms
```

### Cell 7: A / version / cas / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 204800. Client cgroup periods/throttled: 64/35; throttled microseconds: 6044407.

- paused-old-fill: pass=true; old_publish=true served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.028..0.029 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.012..0.013 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.014..0.014 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.081 ms
Execution Time: 0.050 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.019..0.020 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.083 ms
Execution Time: 0.039 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.003..0.004 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.003..0.003 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.002..0.003 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=35 read=2
  I/O Timings: shared read=0.012
Planning Time: 0.129 ms
Execution Time: 0.033 ms
```

### Cell 8: B / sync / cas / soft

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 212992. Client cgroup periods/throttled: 23/0; throttled microseconds: 0.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-delete: pass=true; durable generation and fences rejected old fill; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-marker-window: pass=true; fencing metadata remains present after the marker window; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-recreate: pass=true; served generation=3; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- expiring-fence-resurrection-control: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.018..0.020 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.007..0.008 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.009..0.009 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.100 ms
Execution Time: 0.032 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.018..0.019 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.069 ms
Execution Time: 0.037 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.002..0.002 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.001..0.002 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.001..0.001 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=36 read=2
  I/O Timings: shared read=0.017
Planning Time: 0.163 ms
Execution Time: 0.013 ms
```

### Cell 9: B / sync / cas / hard

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 163840. Client cgroup periods/throttled: 23/1; throttled microseconds: 2805.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-delete: pass=true; durable generation and fences rejected old fill; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-marker-window: pass=true; fencing metadata remains present after the marker window; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-recreate: pass=true; served generation=3; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- expiring-fence-resurrection-control: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.010..0.011 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.006..0.007 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.002..0.002 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.083 ms
Execution Time: 0.022 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.004..0.005 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.021 ms
Execution Time: 0.011 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.002..0.002 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.001..0.002 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.001..0.001 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=36 read=2
  I/O Timings: shared read=0.025
Planning Time: 0.143 ms
Execution Time: 0.011 ms
```

### Cell 10: B / sync / cas / tombstone

Gate: `true`. Error: ``. Wrong/impossible/errors: 0/0/0. Object rows after final trial: 128; source relation/index bytes: 163840. Client cgroup periods/throttled: 24/1; throttled microseconds: 8683.

- paused-old-fill: pass=true; old_publish=false served_generation=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- uncommitted-value: pass=true; impossible key=1 generation=1 required=1; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- stale-read-with-protocol-disabled: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- unguarded-lost-update: pass=true; two intended increments: expected=3 actual=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-delete: pass=true; durable generation and fences rejected old fill; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- absence-marker-expired: pass=true; GET observed the 100-ms marker absent; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-marker-window: pass=true; fencing metadata remains present after the marker window; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- old-fill-after-recreate: pass=true; served generation=3; blocked/recovery/stale-duration ms=0.000/0.000/0.000.
- expiring-fence-resurrection-control: pass=true; stale-after-ack key=1 generation=1 required=2; blocked/recovery/stale-duration ms=0.000/0.000/0.000.

Redis evicted keys before/after: 0/0. Cache policy: volatile-lru; maxmemory: 67108864 bytes.

snapshot plan:

```text
Nested Loop Left Join  (cost=0.30..16.35 rows=1 width=17) (actual time=0.024..0.025 rows=1 loops=1)
  Buffers: shared hit=4
  ->  Index Scan using protocol_generation_pkey on protocol_generation g  (cost=0.15..8.17 rows=1 width=16) (actual time=0.010..0.010 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
  ->  Index Scan using protocol_object_pkey on protocol_object o  (cost=0.15..8.17 rows=1 width=17) (actual time=0.012..0.012 rows=1 loops=1)
        Index Cond: (key_id = '1'::bigint)
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=4
Planning Time: 0.113 ms
Execution Time: 0.054 ms
```

version plan:

```text
Index Scan using protocol_generation_pkey on protocol_generation  (cost=0.15..8.17 rows=1 width=8) (actual time=0.008..0.009 rows=1 loops=1)
  Index Cond: (key_id = '1'::bigint)
  Buffers: shared hit=2
Planning Time: 0.032 ms
Execution Time: 0.018 ms
```

pending plan:

```text
Limit  (cost=0.12..2.59 rows=32 width=30) (actual time=0.003..0.003 rows=0 loops=1)
  Buffers: shared hit=1
  ->  LockRows  (cost=0.12..59.53 rows=770 width=30) (actual time=0.002..0.002 rows=0 loops=1)
        Buffers: shared hit=1
        ->  Index Scan using protocol_outbox_pending on protocol_outbox  (cost=0.12..51.83 rows=770 width=30) (actual time=0.001..0.002 rows=0 loops=1)
              Filter: (NOT delivered)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=36 read=2
  I/O Timings: shared read=0.035
Planning Time: 0.271 ms
Execution Time: 0.019 ms
```

## Signed analyses index

No signed analysis currently indexes this digest.
