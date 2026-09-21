# Study 05 — external cache throughput and consistency — run `devcheck-all`

Generated 2026-09-21T14:28:11Z from 20 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `devcheck-all` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `df7728df338467368098c8eaadfda154f6a7dcb8` |
| `git describe` | `study-05/v1-harness-7-gdf7728d` |
| Working tree dirty at start | false |
| Benchmark image | `localhost/cachebench:1` (`b65a77f01c4b8276ad1`) |
| Resource framing | `db-only` |
| **Inputs digest** | `0f7f36b53f72e7fa` |
| Cells | 20 |
| Distinct scenarios | 20 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 20; cells that failed: 5.
  - failed: `pg-single/legacy-na-memory-aside-strict-coord` (STRICT CONTRACT VIOLATED: 1 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/legacy-na-memory-through-relaxed-coord` (IMPOSSIBLE CACHE VALUE: 87 read(s) returned a state the database never held; the cell is invalid)
  - failed: `pg-single/legacy-na-redis-through-strict-coord` (STRICT CONTRACT VIOLATED: 1 stale-after-ack read(s); the cell cannot support a performance conclusion)
  - failed: `pg-single/owned-opt-redis-through-relaxed-coord` (phase audit: ledger replay failed: INV-8: donor 7 recent slice is not the deterministic newest-20 slice; INV-5: donor 16 holds 9 donations, the ledger promised …)
  - failed: `pg-single/owned-pess-redis-aside-strict-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the committed payload was not accepted as committed, so the dete…)
- Negative controls that fired: 1. Controls that did NOT fire: 0.
- Strict cells that violated their contract: 2.
  - `pg-single/legacy-na-memory-aside-strict-coord`: 1 stale-after-ack read(s)
  - `pg-single/legacy-na-redis-through-strict-coord`: 1 stale-after-ack read(s)
- Cells with impossible cache values: 2 (**these cells are invalid under every freshness policy**).
  - `pg-single/legacy-na-memory-through-relaxed-coord`: 87 impossible cache value(s)
  - `pg-single/owned-pess-redis-aside-strict-coord`: 1 impossible cache value(s)
- Highest and lowest measured throughput across all cells: 170k (`pg-single/o:opt:mem:wth:rel/warm_read_only`) and 2.9k (`pg-single/l:na:red:asd:rel:ext/hotspot_reads`).
- Highest wrong-read rate recorded: 48.72% of all reads (`pg-single/o:opt:red:wth:rel/instances_3`), with 2258 wrong read(s) in 4635.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `0f7f36b53f72e7fa`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale-invalidation` | legacy | na | memory | aside | relaxed | coord | invalidation-only | pass | 0 | 0 | fired |
| pg-single | `legacy-na-memory-aside-relaxed-coord` | legacy | na | memory | aside | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-memory-aside-strict-coord` | legacy | na | memory | aside | strict | coord | invalidation-only | aborted | 1 | 0 |  |
| pg-single | `legacy-na-memory-through-relaxed-coord` | legacy | na | memory | through | relaxed | coord | invalidation-only | aborted | 0 | 87 |  |
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-relaxed-coord` | legacy | na | redis | aside | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-relaxed-ext20` | legacy | na | redis | aside | relaxed | ext20 | authoritative-read | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-strict-coord` | legacy | na | redis | aside | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-strict-ext20` | legacy | na | redis | aside | strict | ext20 | authoritative-read | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-through-relaxed-coord` | legacy | na | redis | through | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `legacy-na-redis-through-strict-coord` | legacy | na | redis | through | strict | coord | invalidation-only | aborted | 1 | 0 |  |
| pg-single | `owned-opt-memory-aside-relaxed-coord` | owned | opt | memory | aside | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-memory-aside-strict-coord` | owned | opt | memory | aside | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-memory-through-relaxed-coord` | owned | opt | memory | through | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-memory-through-strict-coord` | owned | opt | memory | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-redis-aside-relaxed-coord` | owned | opt | redis | aside | relaxed | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-redis-aside-strict-coord` | owned | opt | redis | aside | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-opt-redis-through-relaxed-coord` | owned | opt | redis | through | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | pass | 0 | 0 |  |
| pg-single | `owned-pess-redis-aside-strict-coord` | owned | pess | redis | aside | strict | coord | invalidation-only | aborted | 0 | 1 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | warm_read_only | cacheable | 144k | 0.04 | 0.32 | 7.03 | 57947 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ctl-stale` | cacheable_99_reads | cacheable | 110k | 0.06 | 0.39 | 6.78 | 44201 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ctl-stale` | cacheable_90_reads | cacheable | 100k | 0.06 | 0.42 | 6.92 | 39963 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ctl-stale` | hotspot_reads | cacheable | 139k | 0.04 | 0.33 | 5.07 | 56014 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ctl-stale` | app_99_reads | total app | 7.6k | 0.01 | 3.17 | 220 | 4407 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `ctl-stale` | app_90_reads | total app | 2.0k | 0.02 | 117 | 216 | 1151 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | warm_read_only | cacheable | 132k | 0.05 | 0.32 | 3.85 | 52947 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | cacheable_99_reads | cacheable | 124k | 0.05 | 0.36 | 6.22 | 49658 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | cacheable_90_reads | cacheable | 49k | 0.01 | 0.25 | 223 | 28268 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | hotspot_reads | cacheable | 91k | 0.04 | 0.35 | 222 | 51116 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | app_99_reads | total app | 6.6k | 0.01 | 3.52 | 223 | 3796 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:rel` | app_90_reads | total app | 2.0k | 0.02 | 120 | 219 | 969 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:asd:str` | warm_read_only | cacheable | 127k | 0.05 | 0.32 | 9.14 | 51136 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | cacheable | 115k | 0.06 | 0.36 | 4.91 | 46285 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | cacheable | 57k | 0.03 | 0.38 | 217 | 32106 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | cacheable | 100k | 0.05 | 0.39 | 200 | 45611 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_99_reads | total app | 7.9k | 0.01 | 2.52 | 215 | 4591 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:asd:str` | app_90_reads | total app | 2.0k | 0.02 | 158 | 223 | 1047 | 0 | 0.0 | 4 | 0.00% |
| pg-single | `l:na:mem:wth:rel` | warm_read_only | cacheable | 127k | 0.05 | 0.33 | 8.25 | 51168 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:rel` | cacheable_99_reads | cacheable | 136k | 0.04 | 0.32 | 2.33 | 54801 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:rel` | cacheable_90_reads | cacheable | 14k | 0.01 | 6.21 | 193 | 6624 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:rel` | hotspot_reads | cacheable | 16k | 0.01 | 4.70 | 34 | 6278 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:rel` | app_99_reads | total app | 29k | 0.02 | 5.07 | 35 | 11768 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:rel` | app_90_reads | total app | 9.8k | 0.02 | 7.92 | 215 | 4013 | 0 | 0.0 | 6582 | 1.82% |
| pg-single | `l:na:mem:wth:str` | warm_read_only | cacheable | 128k | 0.05 | 0.32 | 14 | 51261 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | cacheable | 83k | 0.05 | 0.34 | 195 | 48649 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | cacheable | 69k | 0.01 | 0.31 | 196 | 39925 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | cacheable | 71k | 0.06 | 0.46 | 186 | 38393 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | app_99_reads | total app | 8.5k | 0.01 | 5.79 | 210 | 4534 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:mem:wth:str` | app_90_reads | total app | 3.1k | 0.01 | 9.75 | 212 | 1844 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | warm_read_only | cacheable | 12k | 0.39 | 3.36 | 143 | 4874 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel` | cacheable_99_reads | cacheable | 15k | 0.31 | 2.64 | 36 | 5886 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel` | cacheable_90_reads | cacheable | 10.0k | 0.48 | 7.11 | 38 | 4005 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel` | hotspot_reads | cacheable | 7.9k | 0.46 | 14 | 167 | 3170 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel` | app_99_reads | total app | 3.2k | 0.27 | 7.04 | 224 | 1827 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel` | app_90_reads | total app | 1.6k | 0.30 | 164 | 212 | 855 | 0 | 0.0 | 3 | 0.01% |
| pg-single | `l:na:red:asd:rel:ext` | warm_read_only | cacheable | 3.6k | 1.46 | 18 | 24 | 1511 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_99_reads | cacheable | 3.3k | 1.60 | 31 | 35 | 1347 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_90_reads | cacheable | 3.9k | 1.47 | 19 | 26 | 1579 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_reads | cacheable | 2.9k | 1.89 | 31 | 34 | 1155 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | app_99_reads | total app | 3.3k | 1.52 | 33 | 40 | 1323 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | app_90_reads | total app | 3.1k | 1.75 | 30 | 51 | 1245 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | warm_read_only | cacheable | 24k | 0.22 | 1.00 | 35 | 9726 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | cacheable | 25k | 0.24 | 1.03 | 27 | 9997 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | cacheable | 16k | 0.28 | 3.10 | 46 | 6381 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | hotspot_reads | cacheable | 18k | 0.30 | 1.59 | 34 | 7109 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_99_reads | total app | 4.9k | 0.15 | 2.88 | 212 | 2955 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_90_reads | total app | 1.8k | 0.20 | 165 | 220 | 1096 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | warm_read_only | cacheable | 4.4k | 1.35 | 12 | 22 | 1761 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | cacheable_99_reads | cacheable | 3.2k | 1.42 | 34 | 46 | 1295 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | cacheable_90_reads | cacheable | 3.4k | 1.73 | 11 | 36 | 1382 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | hotspot_reads | cacheable | 3.1k | 1.55 | 36 | 39 | 1260 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | app_99_reads | total app | 4.7k | 1.04 | 18 | 30 | 1904 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | app_90_reads | total app | 3.7k | 1.41 | 25 | 36 | 1482 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | warm_read_only | cacheable | 13k | 0.42 | 5.95 | 33 | 5056 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:rel` | cacheable_99_reads | cacheable | 11k | 0.42 | 3.43 | 40 | 4270 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:rel` | cacheable_90_reads | cacheable | 9.1k | 0.52 | 6.92 | 34 | 3647 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:rel` | hotspot_reads | cacheable | 9.8k | 0.42 | 3.76 | 44 | 3940 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:rel` | app_99_reads | total app | 10k | 0.42 | 9.25 | 37 | 4291 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:rel` | app_90_reads | total app | 5.0k | 0.52 | 16 | 44 | 2028 | 0 | 0.0 | 4047 | 5.86% |
| pg-single | `l:na:red:wth:str` | warm_read_only | cacheable | 12k | 0.37 | 3.12 | 37 | 5001 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | cacheable | 11k | 0.40 | 3.33 | 43 | 4266 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | cacheable | 10k | 0.43 | 3.64 | 42 | 4033 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `l:na:red:wth:str` | hotspot_reads | cacheable | 8.4k | 0.55 | 8.72 | 40 | 3381 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `l:na:red:wth:str` | app_99_reads | total app | 5.4k | 0.26 | 4.89 | 215 | 2828 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `l:na:red:wth:str` | app_90_reads | total app | 2.2k | 0.28 | 114 | 218 | 1140 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:mem:asd:rel` | warm_read_only | cacheable | 89k | 0.05 | 0.32 | 217 | 52500 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | cacheable_99_reads | cacheable | 125k | 0.05 | 0.35 | 2.27 | 50119 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | cacheable_90_reads | cacheable | 50k | 0.01 | 0.26 | 224 | 28713 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | hotspot_reads | cacheable | 101k | 0.06 | 0.39 | 196 | 42873 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | app_99_reads | total app | 6.6k | 0.01 | 5.43 | 207 | 3679 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | app_90_reads | total app | 2.2k | 0.01 | 111 | 223 | 1301 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | warm_read_only | cacheable | 111k | 0.05 | 0.30 | 163 | 52028 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | cacheable_99_reads | cacheable | 111k | 0.06 | 0.36 | 116 | 44573 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | cacheable_90_reads | cacheable | 107k | 0.05 | 0.41 | 5.44 | 43026 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | hotspot_reads | cacheable | 110k | 0.04 | 0.30 | 196 | 57780 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | app_99_reads | total app | 6.9k | 0.01 | 3.42 | 227 | 3948 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:str` | app_90_reads | total app | 1.8k | 0.02 | 168 | 216 | 888 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:rel` | warm_read_only | cacheable | 170k | 0.03 | 0.30 | 3.75 | 68377 | 0 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:rel` | cacheable_99_reads | cacheable | 131k | 0.05 | 0.34 | 1.84 | 52777 | 0 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:rel` | cacheable_90_reads | cacheable | 39k | 0.01 | 0.79 | 219 | 19264 | 0 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:rel` | hotspot_reads | cacheable | 115k | 0.05 | 0.38 | 106 | 46045 | 0 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:rel` | app_99_reads | total app | 34k | 0.02 | 4.19 | 32 | 14000 | 0 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:rel` | app_90_reads | total app | 8.2k | 0.01 | 14 | 154 | 3610 | 3 | 0.0 | 7771 | 1.51% |
| pg-single | `o:opt:mem:wth:str` | warm_read_only | cacheable | 120k | 0.05 | 0.43 | 10 | 48402 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:str` | cacheable_99_reads | cacheable | 107k | 0.06 | 0.39 | 5.74 | 42798 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:str` | cacheable_90_reads | cacheable | 74k | 0.01 | 0.20 | 210 | 35290 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:str` | hotspot_reads | cacheable | 96k | 0.06 | 0.42 | 5.57 | 38388 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:str` | app_99_reads | total app | 7.5k | 0.01 | 6.04 | 217 | 4203 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:wth:str` | app_90_reads | total app | 1.9k | 0.02 | 109 | 210 | 988 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | warm_read_only | cacheable | 13k | 0.38 | 1.78 | 145 | 5061 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:rel` | cacheable_99_reads | cacheable | 11k | 0.39 | 2.27 | 44 | 4406 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:rel` | cacheable_90_reads | cacheable | 10k | 0.45 | 2.81 | 41 | 4348 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:rel` | hotspot_reads | cacheable | 9.4k | 0.45 | 5.32 | 47 | 3766 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:rel` | app_99_reads | total app | 2.5k | 0.28 | 108 | 208 | 1486 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:rel` | app_90_reads | total app | 1.1k | 0.34 | 183 | 216 | 559 | 0 | 0.0 | 4 | 0.01% |
| pg-single | `o:opt:red:asd:str` | warm_read_only | cacheable | 12k | 0.36 | 3.56 | 202 | 5088 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | cacheable_99_reads | cacheable | 12k | 0.39 | 2.58 | 43 | 4636 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | cacheable_90_reads | cacheable | 11k | 0.43 | 3.52 | 41 | 4265 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | hotspot_reads | cacheable | 9.7k | 0.46 | 6.33 | 34 | 3873 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | app_99_reads | total app | 2.8k | 0.29 | 112 | 220 | 1661 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | app_90_reads | total app | 968 | 0.34 | 178 | 212 | 574 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | warm_read_only | cacheable | 12k | 0.39 | 5.72 | 38 | 4802 | 0 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:rel` | cacheable_99_reads | cacheable | 11k | 0.39 | 2.97 | 44 | 4312 | 0 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:rel` | cacheable_90_reads | cacheable | 7.7k | 0.49 | 30 | 52 | 3079 | 0 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:rel` | hotspot_reads | cacheable | 7.3k | 0.53 | 39 | 42 | 2934 | 0 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:rel` | app_99_reads | total app | 9.6k | 0.39 | 11 | 67 | 3881 | 0 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:rel` | app_90_reads | total app | 6.3k | 0.27 | 15 | 213 | 2563 | 1 | 0.0 | 4340 | 6.33% |
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 14k | 0.35 | 1.74 | 186 | 5455 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 12k | 0.39 | 3.00 | 41 | 4854 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 11k | 0.46 | 5.51 | 39 | 4378 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | cacheable | 11k | 0.40 | 4.54 | 43 | 4337 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 5.8k | 0.23 | 5.50 | 220 | 3108 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 1.8k | 0.29 | 111 | 224 | 962 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | warm_read_only | cacheable | 11k | 0.43 | 4.83 | 36 | 4644 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | cacheable_99_reads | cacheable | 11k | 0.44 | 6.17 | 35 | 4405 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | cacheable_90_reads | cacheable | 11k | 0.47 | 4.59 | 34 | 4257 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | hotspot_reads | cacheable | 7.5k | 0.54 | 14 | 36 | 3008 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | app_99_reads | total app | 3.3k | 0.24 | 19 | 223 | 1811 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | app_90_reads | total app | 1.0k | 0.32 | 192 | 224 | 567 | 0 | 0.0 | 0 | 0.00% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | hotspot_writes | 939 | 2.07 | 3.66 | 377 | 0 | 0.0 |
| pg-single | `l:na:mem:asd:rel` | hotspot_writes | 535 | 3.53 | 6.25 | 216 | 0 | 0.0 |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 552 | 3.39 | 8.45 | 222 | 0 | 0.0 |
| pg-single | `l:na:mem:wth:rel` | hotspot_writes | 892 | 1.84 | 6.26 | 358 | 0 | 0.0 |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 460 | 4.04 | 13 | 185 | 0 | 0.0 |
| pg-single | `l:na:red:asd:rel` | hotspot_writes | 727 | 2.42 | 7.58 | 292 | 0 | 0.0 |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_writes | 1.2k | 1.62 | 3.56 | 472 | 0 | 0.0 |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 1.2k | 1.58 | 3.60 | 480 | 0 | 0.0 |
| pg-single | `l:na:red:asd:str:ext` | hotspot_writes | 887 | 2.11 | 4.60 | 356 | 0 | 0.0 |
| pg-single | `l:na:red:wth:rel` | hotspot_writes | 526 | 3.52 | 6.41 | 212 | 0 | 0.0 |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 369 | 5.04 | 9.42 | 149 | 0 | 0.0 |
| pg-single | `o:opt:mem:asd:rel` | hotspot_writes | 778 | 2.23 | 7.36 | 314 | 0 | 0.0 |
| pg-single | `o:opt:mem:asd:str` | hotspot_writes | 546 | 2.92 | 12 | 219 | 0 | 0.0 |
| pg-single | `o:opt:mem:wth:rel` | hotspot_writes | 507 | 3.60 | 7.08 | 204 | 0 | 0.0 |
| pg-single | `o:opt:mem:wth:str` | hotspot_writes | 395 | 4.51 | 11 | 159 | 0 | 0.0 |
| pg-single | `o:opt:red:asd:rel` | hotspot_writes | 516 | 3.27 | 10 | 209 | 0 | 0.0 |
| pg-single | `o:opt:red:asd:str` | hotspot_writes | 441 | 4.13 | 8.65 | 178 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:rel` | hotspot_writes | 326 | 5.47 | 13 | 131 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 416 | 4.31 | 9.04 | 167 | 0 | 0.0 |
| pg-single | `o:pess:red:asd:str` | hotspot_writes | 918 | 1.99 | 5.00 | 368 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | cold_fill | 40 | 40 | 0 | 0.00% | 12 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | warm_read_only | 67441 | 67441 | 0 | 0.00% | 67401 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | app_99_reads | 4609 | 4590 | 0 | 0.00% | 4554 | 0.00% | 0 | 0 | — | — | 0 | 82 | 20 | 0 | — |
| pg-single | `ctl-stale` | cacheable_99_reads | 59127 | 59108 | 0 | 0.00% | 59072 | 0.00% | 0 | 0 | — | — | 0 | 82 | 20 | 0 | — |
| pg-single | `ctl-stale` | app_90_reads | 60347 | 60303 | 0 | 0.00% | 60147 | 0.00% | 0 | 0 | — | — | 0 | 144 | 48 | 0 | — |
| pg-single | `ctl-stale` | cacheable_90_reads | 109550 | 109506 | 0 | 0.00% | 109343 | 0.00% | 0 | 0 | — | — | 0 | 405 | 48 | 0 | — |
| pg-single | `ctl-stale` | hotspot_reads | 70546 | 70546 | 0 | 0.00% | 70546 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | hotspot_writes | 70546 | 70546 | 0 | 0.00% | 70546 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | instances_1 | 33782 | 33757 | 0 | 0.00% | 33425 | 0.00% | 0 | 0 | — | — | 0 | 240 | 28 | 0 | — |
| pg-single | `ctl-stale` | instances_3 | 2624 | 2574 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 204 | 131 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | warm_read_only | 63073 | 63073 | 0 | 0.00% | 63033 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | app_99_reads | 3922 | 3899 | 0 | 0.00% | 3858 | 0.00% | 0 | 0 | — | — | 0 | 159 | 24 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | cacheable_99_reads | 62681 | 62658 | 0 | 0.00% | 62617 | 0.00% | 0 | 0 | — | — | 0 | 159 | 24 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | app_90_reads | 63948 | 63900 | 0 | 0.00% | 63740 | 0.00% | 0 | 0 | — | — | 0 | 173 | 51 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | cacheable_90_reads | 99526 | 99478 | 0 | 0.00% | 99167 | 0.00% | 0 | 0 | — | — | 0 | 173 | 51 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | hotspot_reads | 60184 | 60184 | 0 | 0.00% | 60184 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | hotspot_writes | 60184 | 60184 | 0 | 0.00% | 60184 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | instances_1 | 31191 | 31167 | 0 | 0.00% | 30801 | 0.00% | 0 | 0 | — | — | 0 | 204 | 28 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | instances_3 | 2269 | 2222 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 351 | 125 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 17 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | warm_read_only | 62431 | 62431 | 0 | 0.00% | 62391 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_99_reads | 4801 | 4772 | 1 | 0.02% | 4725 | 0.02% | 1 | 1 | 80 | 80 | 1 | 41 | 30 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | 59263 | 59234 | 1 | 0.00% | 59186 | 0.00% | 1 | 1 | 80 | 80 | 1 | 193 | 30 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:mem:asd:str` | app_90_reads | 60467 | 60406 | 1 | 0.00% | 60216 | 0.00% | 1 | 1 | 80 | 80 | 1 | 149 | 67 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | 99047 | 98986 | 1 | 0.00% | 98727 | 0.00% | 1 | 1 | 80 | 80 | 1 | 305 | 67 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | 52587 | 52587 | 0 | 0.00% | 52587 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 52587 | 52587 | 0 | 0.00% | 52587 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_1 | 33313 | 33287 | 0 | 0.00% | 32943 | 0.00% | 0 | 0 | — | — | 0 | 194 | 31 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_3 | 2940 | 2898 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 229 | 131 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 16 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | warm_read_only | 62898 | 62898 | 0 | 0.00% | 62858 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | app_99_reads | 11761 | 10697 | 1064 | 9.05% | 11761 | 9.05% | 12 | 2 | 13 | 1121 | 56 | 25 | 18 | 0 | mixed: cache hit served a superseded entry=1064 |
| pg-single | `l:na:mem:wth:rel` | cacheable_99_reads | 74986 | 73922 | 1064 | 1.42% | 74986 | 1.42% | 12 | 2 | 13 | 1121 | 56 | 25 | 18 | 0 | mixed: cache hit served a superseded entry=1064 |
| pg-single | `l:na:mem:wth:rel` | app_90_reads | 79204 | 77071 | 2024 | 2.56% | 79179 | 2.56% | 12 | 4 | 6.00 | 1118 | 56 | 60 | 57 | 17 | mixed: cache hit served a superseded entry=2024 |
| pg-single | `l:na:mem:wth:rel` | cacheable_90_reads | 87444 | 83970 | 2024 | 2.31% | 87319 | 2.32% | 12 | 4 | 6.00 | 1118 | 56 | 72 | 57 | 17 | mixed: cache hit served a superseded entry=2024 |
| pg-single | `l:na:mem:wth:rel` | hotspot_reads | 7506 | 5682 | 0 | 0.00% | 7506 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | hotspot_writes | 7506 | 5682 | 0 | 0.00% | 7506 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | stampede | 128 | 112 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | instances_1 | 27946 | 27334 | 406 | 1.45% | 27576 | 1.47% | 10 | 5 | 6.00 | 77 | 25 | 64 | 52 | 53 | unattributed: cache hit served a superseded entry=406 |
| pg-single | `l:na:mem:wth:rel` | instances_3 | 3139 | 3089 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 231 | 161 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | warm_read_only | 63933 | 63933 | 0 | 0.00% | 63893 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_99_reads | 5090 | 5061 | 0 | 0.00% | 5011 | 0.00% | 0 | 0 | — | — | 0 | 162 | 33 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | 63000 | 62971 | 0 | 0.00% | 62921 | 0.00% | 0 | 0 | — | — | 0 | 229 | 33 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_90_reads | 64527 | 64462 | 0 | 0.00% | 64319 | 0.00% | 0 | 0 | — | — | 0 | 195 | 90 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | 111919 | 111854 | 0 | 0.00% | 111621 | 0.00% | 0 | 0 | — | — | 0 | 229 | 90 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | 44496 | 44496 | 0 | 0.00% | 44496 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 44496 | 44496 | 0 | 0.00% | 44496 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_1 | 35869 | 35838 | 0 | 0.00% | 35490 | 0.00% | 0 | 0 | — | — | 0 | 212 | 52 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_3 | 2750 | 2714 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 169 | 106 | 0 | — |
| pg-single | `l:na:red:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | warm_read_only | 5608 | 5608 | 0 | 0.00% | 5568 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | app_99_reads | 1752 | 1736 | 0 | 0.00% | 1729 | 0.00% | 0 | 0 | — | — | 0 | 71 | 19 | 0 | — |
| pg-single | `l:na:red:asd:rel` | cacheable_99_reads | 8971 | 8955 | 0 | 0.00% | 8948 | 0.00% | 0 | 0 | — | — | 0 | 71 | 19 | 0 | — |
| pg-single | `l:na:red:asd:rel` | app_90_reads | 9798 | 9750 | 1 | 0.01% | 9671 | 0.01% | 1 | 1 | 75 | 75 | 1 | 198 | 59 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:asd:rel` | cacheable_90_reads | 14909 | 14861 | 1 | 0.01% | 14777 | 0.01% | 1 | 1 | 75 | 75 | 1 | 263 | 59 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:asd:rel` | hotspot_reads | 4117 | 4117 | 0 | 0.00% | 4117 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | hotspot_writes | 4117 | 4117 | 0 | 0.00% | 4117 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | instances_1 | 5303 | 5281 | 0 | 0.00% | 5194 | 0.00% | 0 | 0 | — | — | 0 | 426 | 29 | 0 | — |
| pg-single | `l:na:red:asd:rel` | instances_3 | 5094 | 5064 | 1 | 0.02% | 5000 | 0.02% | 1 | 1 | 47 | 47 | 1 | 201 | 38 | 0 | unattributed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:asd:rel:ext` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | warm_read_only | 2128 | 2128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | app_99_reads | 1397 | 1393 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 9.00 | 13 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_99_reads | 3072 | 3068 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 9.00 | 13 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | app_90_reads | 4196 | 4147 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 99 | 138 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_90_reads | 6066 | 6017 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 147 | 138 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_reads | 1564 | 1564 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_writes | 1564 | 1564 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | instances_1 | 2170 | 2149 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 305 | 78 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | instances_3 | 2301 | 2272 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 283 | 93 | 0 | — |
| pg-single | `l:na:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | warm_read_only | 10231 | 10231 | 0 | 0.00% | 10191 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_99_reads | 2911 | 2889 | 0 | 0.00% | 2861 | 0.00% | 0 | 0 | — | — | 0 | 71 | 29 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | 15200 | 15178 | 0 | 0.00% | 15150 | 0.00% | 0 | 0 | — | — | 0 | 71 | 29 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_90_reads | 16323 | 16261 | 0 | 0.00% | 16138 | 0.00% | 0 | 0 | — | — | 0 | 287 | 80 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | 23741 | 23679 | 0 | 0.00% | 23549 | 0.00% | 0 | 0 | — | — | 0 | 358 | 80 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_reads | 8748 | 8748 | 0 | 0.00% | 8748 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 8748 | 8748 | 0 | 0.00% | 8748 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_1 | 7575 | 7545 | 0 | 0.00% | 7399 | 0.00% | 0 | 0 | — | — | 0 | 500 | 42 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_3 | 6072 | 6043 | 0 | 0.00% | 5957 | 0.00% | 0 | 0 | — | — | 0 | 329 | 38 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | warm_read_only | 2067 | 2067 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | app_99_reads | 2003 | 1989 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 70 | 27 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cacheable_99_reads | 4012 | 3998 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 70 | 27 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | app_90_reads | 5402 | 5328 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 192 | 207 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cacheable_90_reads | 7178 | 7104 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 378 | 207 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | hotspot_reads | 1683 | 1683 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | hotspot_writes | 1683 | 1683 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | instances_1 | 2272 | 2224 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 180 | 146 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | instances_3 | 2253 | 2214 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 142 | 120 | 0 | — |
| pg-single | `l:na:red:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | warm_read_only | 5942 | 5942 | 0 | 0.00% | 5902 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | app_99_reads | 4263 | 4092 | 171 | 4.01% | 4263 | 4.01% | 9 | 3 | 17 | 1160 | 24 | 62 | 48 | 0 | mixed: cache hit served a superseded entry=171 |
| pg-single | `l:na:red:wth:rel` | cacheable_99_reads | 9357 | 9186 | 171 | 1.83% | 9357 | 1.83% | 9 | 3 | 17 | 1160 | 24 | 62 | 48 | 0 | mixed: cache hit served a superseded entry=171 |
| pg-single | `l:na:red:wth:rel` | app_90_reads | 11234 | 10606 | 628 | 5.59% | 11234 | 5.59% | 10 | 4 | 11 | 1090 | 25 | 101 | 142 | 0 | mixed: cache hit served a superseded entry=628 |
| pg-single | `l:na:red:wth:rel` | cacheable_90_reads | 15702 | 15074 | 628 | 4.00% | 15702 | 4.00% | 10 | 4 | 11 | 1090 | 25 | 174 | 142 | 0 | mixed: cache hit served a superseded entry=628 |
| pg-single | `l:na:red:wth:rel` | hotspot_reads | 5070 | 5070 | 0 | 0.00% | 5070 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | hotspot_writes | 5070 | 5070 | 0 | 0.00% | 5070 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | instances_1 | 5669 | 5430 | 234 | 4.13% | 5631 | 4.16% | 9 | 4 | 12 | 619 | 13 | 223 | 65 | 0 | unattributed: cache hit served a superseded entry=234 |
| pg-single | `l:na:red:wth:rel` | instances_3 | 6645 | 4430 | 2215 | 33.33% | 6645 | 33.33% | 21 | 102 | 813 | 5949 | 757 | 96 | 121 | 0 | unattributed: cache hit served a superseded entry=2215 |
| pg-single | `l:na:red:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | warm_read_only | 5986 | 5986 | 0 | 0.00% | 5946 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | app_99_reads | 2958 | 2929 | 1 | 0.03% | 2914 | 0.03% | 1 | 1 | 19 | 19 | 1 | 24 | 39 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | 8841 | 8812 | 1 | 0.01% | 8797 | 0.01% | 1 | 1 | 19 | 19 | 1 | 24 | 39 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:wth:str` | app_90_reads | 9875 | 9810 | 1 | 0.01% | 9754 | 0.01% | 1 | 1 | 19 | 19 | 1 | 144 | 100 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | 15472 | 15407 | 1 | 0.01% | 15351 | 0.01% | 1 | 1 | 19 | 19 | 1 | 221 | 100 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `l:na:red:wth:str` | hotspot_reads | 4591 | 4591 | 0 | 0.00% | 4591 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 4591 | 4591 | 0 | 0.00% | 4591 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_1 | 5956 | 5926 | 0 | 0.00% | 5850 | 0.00% | 0 | 0 | — | — | 0 | 406 | 44 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_3 | 5839 | 5798 | 0 | 0.00% | 5759 | 0.00% | 0 | 0 | — | — | 0 | 611 | 63 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 15 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | warm_read_only | 61922 | 61922 | 0 | 0.00% | 61882 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | app_99_reads | 3713 | 3692 | 0 | 0.00% | 3657 | 0.00% | 0 | 0 | — | — | 0 | 47 | 22 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | cacheable_99_reads | 64366 | 64345 | 0 | 0.00% | 64310 | 0.00% | 0 | 0 | — | — | 0 | 47 | 22 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | app_90_reads | 65685 | 65639 | 0 | 0.00% | 65478 | 0.00% | 0 | 0 | — | — | 0 | 136 | 51 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | cacheable_90_reads | 99573 | 99527 | 0 | 0.00% | 99235 | 0.00% | 0 | 0 | — | — | 0 | 265 | 51 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | hotspot_reads | 51689 | 51689 | 0 | 0.00% | 51689 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | hotspot_writes | 51689 | 51689 | 0 | 0.00% | 51689 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | instances_1 | 26143 | 26118 | 0 | 0.00% | 25861 | 0.00% | 0 | 0 | — | — | 0 | 303 | 25 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | instances_3 | 5814 | 5787 | 0 | 0.00% | 5618 | 0.00% | 0 | 0 | — | — | 0 | 444 | 30 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 12 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | warm_read_only | 61728 | 61728 | 0 | 0.00% | 61688 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | app_99_reads | 4016 | 3999 | 0 | 0.00% | 3944 | 0.00% | 0 | 0 | — | — | 0 | 67 | 20 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | cacheable_99_reads | 57369 | 57352 | 0 | 0.00% | 57296 | 0.00% | 0 | 0 | — | — | 0 | 78 | 20 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | app_90_reads | 58205 | 58161 | 0 | 0.00% | 58002 | 0.00% | 0 | 0 | — | — | 0 | 156 | 59 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | cacheable_90_reads | 110203 | 110159 | 0 | 0.00% | 109996 | 0.00% | 0 | 0 | — | — | 0 | 156 | 59 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | hotspot_reads | 72806 | 72806 | 0 | 0.00% | 72806 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | hotspot_writes | 72806 | 72806 | 0 | 0.00% | 72806 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | instances_1 | 29756 | 29729 | 0 | 0.00% | 29434 | 0.00% | 0 | 0 | — | — | 0 | 343 | 31 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | instances_3 | 9814 | 9782 | 0 | 0.00% | 9532 | 0.00% | 0 | 0 | — | — | 0 | 292 | 42 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | warm_read_only | 83694 | 83694 | 0 | 0.00% | 83654 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | app_99_reads | 13696 | 12192 | 1504 | 10.98% | 13696 | 10.98% | 11 | 1 | 8.00 | 964 | 47 | 22 | 29 | 0 | mixed: cache hit served a superseded entry=1504 |
| pg-single | `o:opt:mem:wth:rel` | cacheable_99_reads | 78876 | 77372 | 1504 | 1.91% | 78876 | 1.91% | 11 | 1 | 8.00 | 964 | 47 | 22 | 29 | 0 | mixed: cache hit served a superseded entry=1504 |
| pg-single | `o:opt:mem:wth:rel` | app_90_reads | 82163 | 79998 | 2163 | 2.63% | 82140 | 2.63% | 12 | 2 | 6.00 | 963 | 47 | 80 | 38 | 0 | mixed: cache hit served a superseded entry=2163 |
| pg-single | `o:opt:mem:wth:rel` | cacheable_90_reads | 103436 | 101271 | 2163 | 2.09% | 103182 | 2.10% | 12 | 2 | 6.00 | 963 | 47 | 80 | 38 | 0 | mixed: cache hit served a superseded entry=2163 |
| pg-single | `o:opt:mem:wth:rel` | hotspot_reads | 54975 | 54975 | 0 | 0.00% | 54975 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | hotspot_writes | 54975 | 54975 | 0 | 0.00% | 54975 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | instances_1 | 30417 | 29973 | 437 | 1.44% | 29955 | 1.46% | 9 | 2 | 5.00 | 102 | 25 | 96 | 10 | 0 | unattributed: cache hit served a superseded entry=437 |
| pg-single | `o:opt:mem:wth:rel` | instances_3 | 11073 | 11041 | 0 | 0.00% | 10838 | 0.00% | 0 | 0 | — | — | 0 | 246 | 37 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 15 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | warm_read_only | 60041 | 60041 | 0 | 0.00% | 60001 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | app_99_reads | 3870 | 3846 | 0 | 0.00% | 3805 | 0.00% | 0 | 0 | — | — | 0 | 40 | 28 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | cacheable_99_reads | 56190 | 56166 | 0 | 0.00% | 56125 | 0.00% | 0 | 0 | — | — | 0 | 161 | 28 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | app_90_reads | 57270 | 57221 | 0 | 0.00% | 57115 | 0.00% | 0 | 0 | — | — | 0 | 230 | 66 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | cacheable_90_reads | 99257 | 99208 | 0 | 0.00% | 99016 | 0.00% | 0 | 0 | — | — | 0 | 230 | 66 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | hotspot_reads | 45005 | 45005 | 0 | 0.00% | 45005 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | hotspot_writes | 45005 | 45005 | 0 | 0.00% | 45005 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | instances_1 | 36667 | 36645 | 0 | 0.00% | 36370 | 0.00% | 0 | 0 | — | — | 0 | 463 | 32 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | instances_3 | 5813 | 5788 | 0 | 0.00% | 5655 | 0.00% | 0 | 0 | — | — | 0 | 329 | 32 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | warm_read_only | 5761 | 5761 | 0 | 0.00% | 5721 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | app_99_reads | 1773 | 1757 | 0 | 0.00% | 1747 | 0.00% | 0 | 0 | — | — | 0 | 10 | 16 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | cacheable_99_reads | 7462 | 7446 | 0 | 0.00% | 7436 | 0.00% | 0 | 0 | — | — | 0 | 10 | 16 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | app_90_reads | 7999 | 7962 | 1 | 0.01% | 7900 | 0.01% | 1 | 1 | 6.00 | 6.00 | 1 | 429 | 45 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `o:opt:red:asd:rel` | cacheable_90_reads | 13189 | 13152 | 1 | 0.01% | 13085 | 0.01% | 1 | 1 | 6.00 | 6.00 | 1 | 473 | 45 | 0 | mixed: cache hit served a superseded entry=1 |
| pg-single | `o:opt:red:asd:rel` | hotspot_reads | 4645 | 4645 | 0 | 0.00% | 4645 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | hotspot_writes | 4645 | 4645 | 0 | 0.00% | 4645 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | instances_1 | 4905 | 4882 | 0 | 0.00% | 4788 | 0.00% | 0 | 0 | — | — | 0 | 507 | 28 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | instances_3 | 5714 | 5688 | 2 | 0.04% | 5627 | 0.04% | 1 | 1 | 804 | 804 | 2 | 246 | 32 | 0 | unattributed: cache hit served a superseded entry=2 |
| pg-single | `o:opt:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 15 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | warm_read_only | 5741 | 5741 | 0 | 0.00% | 5701 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | app_99_reads | 1612 | 1590 | 0 | 0.00% | 1579 | 0.00% | 0 | 0 | — | — | 0 | 25 | 25 | 0 | — |
| pg-single | `o:opt:red:asd:str` | cacheable_99_reads | 8153 | 8131 | 0 | 0.00% | 8120 | 0.00% | 0 | 0 | — | — | 0 | 25 | 25 | 0 | — |
| pg-single | `o:opt:red:asd:str` | app_90_reads | 8719 | 8667 | 0 | 0.00% | 8600 | 0.00% | 0 | 0 | — | — | 0 | 257 | 61 | 0 | — |
| pg-single | `o:opt:red:asd:str` | cacheable_90_reads | 13293 | 13241 | 0 | 0.00% | 13171 | 0.00% | 0 | 0 | — | — | 0 | 497 | 61 | 0 | — |
| pg-single | `o:opt:red:asd:str` | hotspot_reads | 4851 | 4851 | 0 | 0.00% | 4851 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | hotspot_writes | 4851 | 4851 | 0 | 0.00% | 4851 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | instances_1 | 5105 | 5072 | 0 | 0.00% | 4962 | 0.00% | 0 | 0 | — | — | 0 | 450 | 39 | 0 | — |
| pg-single | `o:opt:red:asd:str` | instances_3 | 4379 | 4354 | 0 | 0.00% | 4301 | 0.00% | 0 | 0 | — | — | 0 | 314 | 33 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 16 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | warm_read_only | 6348 | 6348 | 0 | 0.00% | 6308 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | app_99_reads | 4011 | 3789 | 222 | 5.53% | 4011 | 5.53% | 7 | 1 | 34 | 1182 | 27 | 94 | 37 | 0 | mixed: cache hit served a superseded entry=222 |
| pg-single | `o:opt:red:wth:rel` | cacheable_99_reads | 9733 | 9511 | 222 | 2.28% | 9733 | 2.28% | 7 | 1 | 34 | 1182 | 27 | 94 | 37 | 0 | mixed: cache hit served a superseded entry=222 |
| pg-single | `o:opt:red:wth:rel` | app_90_reads | 11792 | 11134 | 656 | 5.56% | 11791 | 5.56% | 9 | 2 | 8.00 | 1056 | 27 | 109 | 108 | 0 | mixed: cache hit served a superseded entry=656 |
| pg-single | `o:opt:red:wth:rel` | cacheable_90_reads | 16393 | 15735 | 656 | 4.00% | 16392 | 4.00% | 9 | 2 | 8.00 | 1056 | 27 | 109 | 108 | 0 | mixed: cache hit served a superseded entry=656 |
| pg-single | `o:opt:red:wth:rel` | hotspot_reads | 3857 | 3857 | 0 | 0.00% | 3857 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | hotspot_writes | 3857 | 3857 | 0 | 0.00% | 3857 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | instances_1 | 7725 | 7373 | 326 | 4.22% | 7691 | 4.24% | 8 | 2 | 5.00 | 134 | 13 | 145 | 74 | 0 | unattributed: cache hit served a superseded entry=326 |
| pg-single | `o:opt:red:wth:rel` | instances_3 | 4635 | 2357 | 2258 | 48.72% | 4635 | 48.72% | 16 | 31 | 1961 | 4618 | 487 | 1.00 | 25 | 0 | unattributed: cache hit served a superseded entry=2258 |
| pg-single | `o:opt:red:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 17 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 6023 | 6023 | 0 | 0.00% | 5983 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 3129 | 3107 | 0 | 0.00% | 3088 | 0.00% | 0 | 0 | — | — | 0 | 169 | 32 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 9439 | 9417 | 0 | 0.00% | 9398 | 0.00% | 0 | 0 | — | — | 0 | 175 | 32 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 10263 | 10210 | 0 | 0.00% | 10136 | 0.00% | 0 | 0 | — | — | 0 | 175 | 94 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 16261 | 16208 | 0 | 0.00% | 16134 | 0.00% | 0 | 0 | — | — | 0 | 229 | 94 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | 5573 | 5573 | 0 | 0.00% | 5573 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 5573 | 5573 | 0 | 0.00% | 5573 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_1 | 5208 | 5176 | 0 | 0.00% | 5100 | 0.00% | 0 | 0 | — | — | 0 | 373 | 46 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_3 | 5550 | 5522 | 0 | 0.00% | 5485 | 0.00% | 0 | 0 | — | — | 0 | 308 | 46 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | warm_read_only | 5281 | 5281 | 0 | 0.00% | 5241 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | app_99_reads | 2128 | 2108 | 0 | 0.00% | 2089 | 0.00% | 0 | 0 | — | — | 0 | 14 | 22 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cacheable_99_reads | 7502 | 7482 | 0 | 0.00% | 7461 | 0.00% | 0 | 0 | — | — | 0 | 157 | 22 | 0 | — |
| pg-single | `o:pess:red:asd:str` | app_90_reads | 8052 | 8006 | 0 | 0.00% | 7933 | 0.00% | 0 | 0 | — | — | 0 | 159 | 51 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cacheable_90_reads | 13188 | 13142 | 0 | 0.00% | 13061 | 0.00% | 0 | 0 | — | — | 0 | 390 | 51 | 0 | — |
| pg-single | `o:pess:red:asd:str` | hotspot_reads | 3917 | 3917 | 0 | 0.00% | 3917 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | hotspot_writes | 3917 | 3917 | 0 | 0.00% | 3917 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | stampede | 128 | 112 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | instances_1 | 3117 | 2924 | 0 | 0.00% | 3002 | 0.00% | 0 | 0 | — | — | 0 | 342 | 32 | 0 | — |
| pg-single | `o:pess:red:asd:str` | instances_3 | 4703 | 4440 | 0 | 0.00% | 4622 | 0.00% | 1 | 0 | — | — | 0 | 421 | 37 | 1 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

| Topology | Scenario | Age bucket | hits | early expiries | observed % | specified p (band) |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 0-75s | 280715 | 589 | 0.2% | 0 to 0.25 |
| pg-single | `l:na:mem:asd:rel` | 0-75s | 253185 | 598 | 0.2% | 0 to 0.25 |
| pg-single | `l:na:mem:asd:str` | 0-75s | 246648 | 550 | 0.2% | 0 to 0.25 |
| pg-single | `l:na:mem:wth:rel` | 0-75s | 185259 | 262 | 0.1% | 0 to 0.25 |
| pg-single | `l:na:mem:wth:str` | 0-75s | 255500 | 575 | 0.2% | 0 to 0.25 |
| pg-single | `l:na:red:asd:rel` | 0-75s | 34656 | 160 | 0.5% | 0 to 0.25 |
| pg-single | `l:na:red:asd:str` | 0-75s | 55844 | 110 | 0.2% | 0 to 0.25 |
| pg-single | `l:na:red:wth:rel` | 0-75s | 38948 | 57 | 0.1% | 0 to 0.25 |
| pg-single | `l:na:red:wth:rel` | <0s | 2 | 0 | 0.0% | 0 |
| pg-single | `l:na:red:wth:str` | 0-75s | 37497 | 73 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:mem:asd:rel` | 0-75s | 244285 | 576 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:mem:asd:str` | 0-75s | 283456 | 633 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:mem:wth:rel` | 0-75s | 282604 | 670 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:mem:wth:str` | 0-75s | 246047 | 529 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:red:asd:rel` | 0-75s | 33866 | 61 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:red:asd:str` | 0-75s | 32986 | 67 | 0.2% | 0 to 0.25 |
| pg-single | `o:opt:red:wth:str` | 0-75s | 38275 | 68 | 0.2% | 0 to 0.25 |

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 775 | 83 | 0 | — | 83 | 10 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.302ms) |
| pg-single | `l:na:mem:asd:rel` | 946 | 106 | 0 | — | 106 | 7 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.81ms) |
| pg-single | `l:na:mem:asd:str` | 852 | 105 | 0 | — | 105 | 19 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.191ms) |
| pg-single | `l:na:mem:wth:rel` | 425 | 30 | 0 | — | 30 | 1 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.582ms) |
| pg-single | `l:na:mem:wth:str` | 864 | 96 | 0 | — | 96 | 9 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.284ms) |
| pg-single | `l:na:red:asd:rel` | 318 | 80 | 0 | — | 80 | 11 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.231ms) |
| pg-single | `l:na:red:asd:rel:ext` | 0 | 0 | 0 | — | 0 | 0 | 0 | 500ms | no fill samples; lease left at the 500 ms default (NOT calibrated) |
| pg-single | `l:na:red:asd:str` | 336 | 86 | 0 | — | 86 | 4 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 806µs) |
| pg-single | `l:na:red:asd:str:ext` | 0 | 0 | 0 | — | 0 | 0 | 0 | 500ms | no fill samples; lease left at the 500 ms default (NOT calibrated) |
| pg-single | `l:na:red:wth:rel` | 98 | 22 | 0 | — | 22 | 0 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 903µs) |
| pg-single | `l:na:red:wth:str` | 233 | 67 | 0 | — | 67 | 6 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.576ms) |
| pg-single | `o:opt:mem:asd:rel` | 921 | 102 | 0 | — | 102 | 6 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.026999ms) |
| pg-single | `o:opt:mem:asd:str` | 828 | 83 | 0 | — | 83 | 12 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.376ms) |
| pg-single | `o:opt:mem:wth:rel` | 740 | 54 | 0 | — | 54 | 5 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 2.608ms) |
| pg-single | `o:opt:mem:wth:str` | 747 | 96 | 0 | — | 96 | 6 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.829ms) |
| pg-single | `o:opt:red:asd:rel` | 203 | 74 | 0 | — | 74 | 5 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.567ms) |
| pg-single | `o:opt:red:asd:str` | 225 | 85 | 0 | — | 85 | 9 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.68ms) |
| pg-single | `o:opt:red:wth:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:wth:str` | 249 | 73 | 0 | — | 73 | 6 | 0 | 100ms | 4 x measured p99 fill latency, raised to the 100 ms floor (p99 fill = 1.626ms) |
| pg-single | `o:pess:red:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |

| Topology | Scenario | readers | keys | elapsed ms | db loads | loads/key | fallbacks | duplicate fills | wrong | impossible |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 16 | 8 | 210.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:rel` | 16 | 8 | 210.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:rel` | 16 | 8 | 213.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:str` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel` | 16 | 8 | 227.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel:ext` | 16 | 8 | 46.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str` | 16 | 8 | 227.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str:ext` | 16 | 8 | 55.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:rel` | 16 | 8 | 230.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | 16 | 8 | 228.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:rel` | 16 | 8 | 210.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:str` | 16 | 8 | 211.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:wth:rel` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:wth:str` | 16 | 8 | 211.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:rel` | 16 | 8 | 227.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:str` | 16 | 8 | 228.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:rel` | 16 | 8 | 213.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | 16 | 8 | 226.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:pess:red:asd:str` | 16 | 8 | 227.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 1 | 8 | memory | 1.9k | 0.0 | 48k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `ctl-stale` | 3 | 8 | memory | 2.9k | 0.0 | 3.0k | 2624 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:asd:rel` | 1 | 8 | memory | 1.9k | 0.0 | 41k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:rel` | 3 | 8 | memory | 2.7k | 0.0 | 2.4k | 2269 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:asd:str` | 1 | 8 | memory | 1.8k | 0.0 | 49k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:str` | 3 | 8 | memory | 3.7k | 0.0 | 2.7k | 2940 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:wth:rel` | 1 | 8 | memory | 7.1k | 0.0 | 30k | 0 | 406 | 53 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:rel` | 3 | 8 | memory | 4.1k | 0.0 | 2.7k | 3139 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:wth:str` | 1 | 8 | memory | 2.1k | 0.0 | 51k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:str` | 3 | 8 | memory | 2.8k | 0.0 | 3.2k | 2750 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:red:asd:rel` | 1 | 8 | redis | 1.2k | 0.0 | 9.3k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel` | 3 | 8 | redis | 1.4k | 0.0 | 7.6k | 0 | 1 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel:ext` | 1 | 8 | redis | 2.5k | 0.0 | 2.6k | 2170 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel:ext` | 3 | 8 | redis | 2.1k | 0.0 | 2.8k | 2301 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str` | 1 | 8 | redis | 1.6k | 0.0 | 15k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str` | 3 | 8 | redis | 1.3k | 0.0 | 10k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str:ext` | 1 | 8 | redis | 2.9k | 0.0 | 2.4k | 2272 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str:ext` | 3 | 8 | redis | 2.9k | 0.0 | 2.2k | 2253 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:rel` | 1 | 8 | redis | 4.3k | 0.0 | 7.6k | 0 | 234 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:rel` | 3 | 8 | redis | 6.2k | 0.0 | 7.8k | 0 | 2215 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 1 | 8 | redis | 1.3k | 0.0 | 11k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 3 | 8 | redis | 2.0k | 0.0 | 8.9k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:rel` | 1 | 8 | memory | 1.9k | 0.0 | 40k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:rel` | 3 | 8 | memory | 1.2k | 0.0 | 11k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:asd:str` | 1 | 8 | memory | 1.5k | 0.0 | 41k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:str` | 3 | 8 | memory | 1.5k | 0.0 | 18k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:wth:rel` | 1 | 8 | memory | 6.6k | 0.0 | 33k | 0 | 437 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:wth:rel` | 3 | 8 | memory | 1.8k | 0.0 | 15k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:wth:str` | 1 | 8 | memory | 1.1k | 0.0 | 52k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:wth:str` | 3 | 8 | memory | 1.2k | 0.0 | 10.0k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:red:asd:rel` | 1 | 8 | redis | 1.1k | 0.0 | 9.1k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:rel` | 3 | 8 | redis | 1.4k | 0.0 | 10k | 0 | 2 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:str` | 1 | 8 | redis | 1.2k | 0.0 | 8.4k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:str` | 3 | 8 | redis | 916 | 0.0 | 9.0k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:rel` | 1 | 8 | redis | 6.3k | 0.0 | 11k | 0 | 326 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:rel` | 3 | 8 | redis | 608 | 0.0 | 6.4k | 0 | 2258 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 1 | 8 | redis | 1.3k | 0.0 | 8.6k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 3 | 8 | redis | 1.3k | 0.0 | 9.5k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:pess:red:asd:str` | 1 | 8 | redis | 1.0k | 0.0 | 5.9k | 0 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:pess:red:asd:str` | 3 | 8 | redis | 797 | 0.0 | 8.7k | 0 | 0 | 1 | shared store: invalidation is visible to every instance |

## Fault injection and controls

A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly
has not been fast.

| Topology | Scenario | Fault | reproduced | wrong observed | impossible observed | correctness |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `ctl-stale` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `ctl-stale` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `ctl-stale` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `ctl-stale` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `ctl-stale` | suppressed-invalidation-produces-a-stale-read | true | 1 | 0 | one suppressed post-commit invalidation produced 1 stale-after-ack read(s) that the oracle detected |
| pg-single | `ctl-stale` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:mem:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:mem:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:red:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:asd:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:asd:rel:ext` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:rel:ext` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `l:na:red:asd:rel:ext` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 500ms; no manual intervention required |
| pg-single | `l:na:red:asd:rel:ext` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:asd:rel:ext` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:asd:rel:ext` | legacy-external-writer-bypasses-adapter | true | 0 | 0 | relaxed: the external write was invisible to the cache; the read returned fresh (source bypass) and any staleness is counted |
| pg-single | `l:na:red:asd:rel:ext` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state and the detector classified it impossible (this scenario reads the database authoritatively, so there is no cache path to exercise) |
| pg-single | `l:na:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:asd:str:ext` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:str:ext` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:asd:str:ext` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 500ms; no manual intervention required |
| pg-single | `l:na:red:asd:str:ext` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:asd:str:ext` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:asd:str:ext` | legacy-external-writer-bypasses-adapter | true | 0 | 0 | strict: the read could not trust the cache, so it bypass and returned fresh |
| pg-single | `l:na:red:asd:str:ext` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state and the detector classified it impossible (this scenario reads the database authoritatively, so there is no cache path to exercise) |
| pg-single | `l:na:red:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:red:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:wth:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:wth:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:wth:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `l:na:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `l:na:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:red:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:red:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:red:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:mem:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:mem:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:asd:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:mem:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:mem:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:wth:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:wth:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:wth:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:red:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:asd:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:asd:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:red:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:wth:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:wth:rel` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:wth:rel` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:opt:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:opt:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:red:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:red:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:red:wth:str` | dirty-cache-write-audit-rejects-impossible-version | true | 0 | 0 | the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned impossible (source hit) |
| pg-single | `o:pess:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (fill -> fill); the next read returned fresh |
| pg-single | `o:pess:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:pess:red:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:pess:red:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:pess:red:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:pess:red:asd:str` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |

## Ledger replay (audit)

| Topology | Scenario | Phase # | passed | checks | content mismatches | donation-set mismatches | aggregate mismatches | recent-slice mismatches | design checks | first failure |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:asd:rel` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:asd:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:wth:rel` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:mem:wth:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:asd:rel` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:asd:rel:ext` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:asd:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:asd:str:ext` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:wth:rel` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `l:na:red:wth:str` | 1 | true | 13 | 0 | 0 | 0 | 0 | a_rolldown_drift=0 |  |
| pg-single | `o:opt:mem:asd:rel` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:mem:asd:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:mem:wth:rel` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:mem:wth:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:red:asd:rel` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:red:asd:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |
| pg-single | `o:opt:red:wth:rel` | 1 | false | 26 | 1 | 1 | 1 | 2 | a_outbox_orphans=0 a_rolldown_drift=0 | INV-8: donor 7 recent slice is not the deterministic newest-20 slice |
| pg-single | `o:opt:red:wth:str` | 1 | true | 26 | 0 | 0 | 0 | 0 | a_outbox_orphans=0 a_rolldown_drift=0 |  |

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstone fences | version validations | bypass reads | external writes | unrecorded states confirmed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | memory | 48.0 KiB | 4.3 KiB | 2 | 0 | 258 | 256 | 2 | 0 | 719 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:mem:asd:rel` | memory | 48.0 KiB | 2.1 KiB | 1 | 150 | 412 | 409 | 3 | 0 | 499 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | memory | 48.0 KiB | 2.1 KiB | 1 | 65 | 362 | 346 | 16 | 0 | 0 | 1023 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:mem:wth:rel` | memory | 48.0 KiB | 6.4 KiB | 3 | 129 | 1295 | 1292 | 3 | 2 | 0 | 2 | 0 | 20 | 0 | 3273 |
| pg-single | `l:na:mem:wth:str` | memory | 48.0 KiB | 6.4 KiB | 3 | 92 | 851 | 745 | 106 | 2 | 0 | 1009 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:red:asd:rel` | redis | 48.0 KiB | — | 0 | 0 | 181 | 173 | 8 | 0 | 518 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:red:asd:rel:ext` | redis | 48.0 KiB | — | 0 | 0 | 0 | 0 | 0 | 0 | 617 | 2 | 0 | 9943 | 146 | 0 |
| pg-single | `l:na:red:asd:str` | redis | 48.0 KiB | — | 0 | 0 | 249 | 228 | 21 | 0 | 0 | 1531 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:red:asd:str:ext` | redis | 48.0 KiB | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 1260 | 0 | 11113 | 132 | 0 |
| pg-single | `l:na:red:wth:rel` | redis | 48.0 KiB | — | 0 | 0 | 689 | 686 | 3 | 2 | 0 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | redis | 48.0 KiB | — | 0 | 0 | 544 | 475 | 69 | 2 | 0 | 743 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:mem:asd:rel` | memory | 48.0 KiB | 2.2 KiB | 1 | 129 | 392 | 388 | 4 | 0 | 620 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:mem:asd:str` | memory | 48.0 KiB | 2.1 KiB | 1 | 0 | 256 | 239 | 17 | 0 | 0 | 911 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:mem:wth:rel` | memory | 48.0 KiB | 6.5 KiB | 3 | 250 | 1163 | 1159 | 4 | 2 | 0 | 8 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:mem:wth:str` | memory | 48.0 KiB | 6.4 KiB | 3 | 88 | 695 | 651 | 44 | 2 | 0 | 805 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:red:asd:rel` | redis | 48.0 KiB | — | 0 | 0 | 159 | 155 | 4 | 0 | 385 | 2 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:red:asd:str` | redis | 48.0 KiB | — | 0 | 0 | 174 | 168 | 6 | 0 | 0 | 669 | 0 | 20 | 0 | 0 |
| pg-single | `o:opt:red:wth:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | redis | 48.0 KiB | — | 0 | 0 | 556 | 510 | 46 | 2 | 0 | 755 | 0 | 20 | 0 | 0 |
| pg-single | `o:pess:red:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Redis accounting (from `INFO`)

Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in
behaviour between the two is partly a policy difference and is stated as such.

| Topology | Scenario | redis_version | maxmemory | maxmemory_policy | maxmemory_samples | used_memory | used_memory_peak | keyspace_hits | keyspace_misses | expired_keys | evicted_keys | total_commands_processed | total_net_input_bytes | total_net_output_bytes | used_cpu_sys | used_cpu_user |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `l:na:red:asd:rel` | 7.4.11 | 33554432 | allkeys-lru | — | 1554232 | 2214968 | 38622 | 1461 | 2 | 0 | 50281 | 2634023 | 53503849 | 1.336593 | 0.739158 |
| pg-single | `l:na:red:asd:rel:ext` | 7.4.11 | 33554432 | allkeys-lru | — | 1223368 | 2214968 | 74643 | 2934 | 4 | 0 | 95447 | 5045091 | 103731349 | 2.656172 | 1.321947 |
| pg-single | `l:na:red:asd:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1555576 | 2521224 | 366318 | 15904 | 11 | 0 | 437890 | 28962128 | 540290983 | 10.884323 | 5.492467 |
| pg-single | `l:na:red:asd:str:ext` | 7.4.11 | 33554432 | allkeys-lru | — | 1244736 | 1378280 | 1861 | 0 | 1 | 0 | 7456 | 335586 | 11440 | 0.163426 | 0.279771 |
| pg-single | `l:na:red:wth:rel` | 7.4.11 | 33554432 | allkeys-lru | — | 1571280 | 2521224 | 115167 | 5875 | 5 | 0 | 141743 | 9596877 | 170664606 | 3.833052 | 1.966038 |
| pg-single | `l:na:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1447544 | 2521224 | 156307 | 7617 | 6 | 0 | 191388 | 13425138 | 230358878 | 5.056485 | 2.560569 |
| pg-single | `o:opt:red:asd:rel` | 7.4.11 | 33554432 | allkeys-lru | — | 1537856 | 2521224 | 232530 | 10832 | 8 | 0 | 277074 | 19360446 | 349329836 | 7.309891 | 3.697673 |
| pg-single | `o:opt:red:asd:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1547160 | 2214968 | 73815 | 2934 | 3 | 0 | 92123 | 4895569 | 103715884 | 2.569265 | 1.161631 |
| pg-single | `o:opt:red:wth:str` | 7.4.11 | 33554432 | allkeys-lru | — | 1560480 | 2521224 | 306724 | 14015 | 10 | 0 | 367351 | 25323836 | 454849822 | 9.668183 | 4.866442 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 77 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:asd:rel` | 83 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:wth:str` | 83 | 2 | 12 | db-only | 8 |
| pg-single | `l:na:red:asd:rel` | 81 | 26 | 3358 | db-only | 8 |
| pg-single | `l:na:red:asd:rel:ext` | 61 | 12 | 193 | db-only | 8 |
| pg-single | `l:na:red:asd:str` | 80 | 23 | 2256 | db-only | 8 |
| pg-single | `l:na:red:asd:str:ext` | 61 | 4 | 80 | db-only | 8 |
| pg-single | `l:na:red:wth:rel` | 67 | 43 | 4886 | db-only | 8 |
| pg-single | `l:na:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:asd:rel` | 84 | 2 | 55 | db-only | 8 |
| pg-single | `o:opt:mem:asd:str` | 82 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:wth:rel` | 78 | 6 | 58 | db-only | 8 |
| pg-single | `o:opt:mem:wth:str` | 85 | 3 | 313 | db-only | 8 |
| pg-single | `o:opt:red:asd:rel` | 81 | 25 | 3295 | db-only | 8 |
| pg-single | `o:opt:red:asd:str` | 82 | 23 | 3410 | db-only | 8 |
| pg-single | `o:opt:red:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:str` | 76 | 28 | 2622 | db-only | 8 |
| pg-single | `o:pess:red:asd:str` | 0 | 0 | 0 | db-only | 8 |

## Storage

PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a
different query, so none is reported there rather than a number read as bytes on disk.

| Topology | Scenario | Relation | Bytes |
|---|---|---|---|
| pg-single | `ctl-stale` | charity | 32.0 KiB |
| pg-single | `ctl-stale` | donation | 176.0 KiB |
| pg-single | `ctl-stale` | load_donation | 64.0 KiB |
| pg-single | `ctl-stale` | person | 48.0 KiB |
| pg-single | `l:na:mem:asd:rel` | charity | 32.0 KiB |
| pg-single | `l:na:mem:asd:rel` | donation | 176.0 KiB |
| pg-single | `l:na:mem:asd:rel` | load_donation | 64.0 KiB |
| pg-single | `l:na:mem:asd:rel` | person | 48.0 KiB |
| pg-single | `l:na:mem:asd:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:asd:str` | donation | 176.0 KiB |
| pg-single | `l:na:mem:asd:str` | load_donation | 64.0 KiB |
| pg-single | `l:na:mem:asd:str` | person | 48.0 KiB |
| pg-single | `l:na:mem:wth:rel` | charity | 32.0 KiB |
| pg-single | `l:na:mem:wth:rel` | donation | 176.0 KiB |
| pg-single | `l:na:mem:wth:rel` | load_donation | 64.0 KiB |
| pg-single | `l:na:mem:wth:rel` | person | 48.0 KiB |
| pg-single | `l:na:mem:wth:str` | charity | 32.0 KiB |
| pg-single | `l:na:mem:wth:str` | donation | 176.0 KiB |
| pg-single | `l:na:mem:wth:str` | load_donation | 64.0 KiB |
| pg-single | `l:na:mem:wth:str` | person | 48.0 KiB |
| pg-single | `l:na:red:asd:rel` | charity | 32.0 KiB |
| pg-single | `l:na:red:asd:rel` | donation | 176.0 KiB |
| pg-single | `l:na:red:asd:rel` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:asd:rel` | person | 48.0 KiB |
| pg-single | `l:na:red:asd:rel:ext` | charity | 32.0 KiB |
| pg-single | `l:na:red:asd:rel:ext` | donation | 176.0 KiB |
| pg-single | `l:na:red:asd:rel:ext` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:asd:rel:ext` | person | 48.0 KiB |
| pg-single | `l:na:red:asd:str` | charity | 32.0 KiB |
| pg-single | `l:na:red:asd:str` | donation | 176.0 KiB |
| pg-single | `l:na:red:asd:str` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:asd:str` | person | 48.0 KiB |
| pg-single | `l:na:red:asd:str:ext` | charity | 32.0 KiB |
| pg-single | `l:na:red:asd:str:ext` | donation | 176.0 KiB |
| pg-single | `l:na:red:asd:str:ext` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:asd:str:ext` | person | 48.0 KiB |
| pg-single | `l:na:red:wth:rel` | charity | 32.0 KiB |
| pg-single | `l:na:red:wth:rel` | donation | 176.0 KiB |
| pg-single | `l:na:red:wth:rel` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:wth:rel` | person | 48.0 KiB |
| pg-single | `l:na:red:wth:str` | charity | 32.0 KiB |
| pg-single | `l:na:red:wth:str` | donation | 176.0 KiB |
| pg-single | `l:na:red:wth:str` | load_donation | 64.0 KiB |
| pg-single | `l:na:red:wth:str` | person | 48.0 KiB |
| pg-single | `o:opt:mem:asd:rel` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:mem:asd:rel` | charity | 32.0 KiB |
| pg-single | `o:opt:mem:asd:rel` | donation | 176.0 KiB |
| pg-single | `o:opt:mem:asd:rel` | load_donation | 64.0 KiB |
| pg-single | `o:opt:mem:asd:rel` | person | 48.0 KiB |
| pg-single | `o:opt:mem:asd:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:mem:asd:str` | charity | 32.0 KiB |
| pg-single | `o:opt:mem:asd:str` | donation | 176.0 KiB |
| pg-single | `o:opt:mem:asd:str` | load_donation | 64.0 KiB |
| pg-single | `o:opt:mem:asd:str` | person | 48.0 KiB |
| pg-single | `o:opt:mem:wth:rel` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:mem:wth:rel` | charity | 32.0 KiB |
| pg-single | `o:opt:mem:wth:rel` | donation | 176.0 KiB |
| pg-single | `o:opt:mem:wth:rel` | load_donation | 64.0 KiB |
| pg-single | `o:opt:mem:wth:rel` | person | 48.0 KiB |
| pg-single | `o:opt:mem:wth:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:mem:wth:str` | charity | 32.0 KiB |
| pg-single | `o:opt:mem:wth:str` | donation | 176.0 KiB |
| pg-single | `o:opt:mem:wth:str` | load_donation | 64.0 KiB |
| pg-single | `o:opt:mem:wth:str` | person | 48.0 KiB |
| pg-single | `o:opt:red:asd:rel` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:asd:rel` | charity | 32.0 KiB |
| pg-single | `o:opt:red:asd:rel` | donation | 176.0 KiB |
| pg-single | `o:opt:red:asd:rel` | load_donation | 64.0 KiB |
| pg-single | `o:opt:red:asd:rel` | person | 48.0 KiB |
| pg-single | `o:opt:red:asd:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:asd:str` | charity | 32.0 KiB |
| pg-single | `o:opt:red:asd:str` | donation | 176.0 KiB |
| pg-single | `o:opt:red:asd:str` | load_donation | 64.0 KiB |
| pg-single | `o:opt:red:asd:str` | person | 48.0 KiB |
| pg-single | `o:opt:red:wth:rel` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:wth:rel` | charity | 32.0 KiB |
| pg-single | `o:opt:red:wth:rel` | donation | 176.0 KiB |
| pg-single | `o:opt:red:wth:rel` | load_donation | 64.0 KiB |
| pg-single | `o:opt:red:wth:rel` | person | 48.0 KiB |
| pg-single | `o:opt:red:wth:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:opt:red:wth:str` | charity | 32.0 KiB |
| pg-single | `o:opt:red:wth:str` | donation | 176.0 KiB |
| pg-single | `o:opt:red:wth:str` | load_donation | 64.0 KiB |
| pg-single | `o:opt:red:wth:str` | person | 48.0 KiB |
| pg-single | `o:pess:red:asd:str` | cache_outbox | 24.0 KiB |
| pg-single | `o:pess:red:asd:str` | charity | 32.0 KiB |
| pg-single | `o:pess:red:asd:str` | donation | 176.0 KiB |
| pg-single | `o:pess:red:asd:str` | load_donation | 64.0 KiB |
| pg-single | `o:pess:red:asd:str` | person | 48.0 KiB |

## Correctness gates

| Topology | Scenario | gate passed | checks | statements | failures |
|---|---|---|---|---|---|
| pg-single | `ctl-stale` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:asd:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:asd:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:asd:rel:ext` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:asd:str:ext` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:wth:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:red:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `o:opt:mem:asd:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:mem:asd:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:mem:wth:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:mem:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:asd:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:asd:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:wth:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:pess:red:asd:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |

## Limitations stated with the numbers

- The client, the cache and the database share the same eight heterogeneous cores of a laptop. Absolute
  throughput is lower than a two-machine setup would give, and only relative comparisons within this
  environment are meaningful.
- There is no real network. A shared cache and a distributed database both look better here than they
  would across availability zones; EXPLAIN (ANALYZE, DIST) RPC counts are the portable signal.
- All containers run on one WSL2 machine. This is NOT evidence of data colocation or of network
  behaviour; the placement pair reports engine evidence and physical colocation remains untested.
- The workloads are closed-loop per operation within a phase (measure.Run). Deep tails are therefore
  optimistic floors and are compared between scenarios, never quoted as SLO figures. The open-loop
  scheduler is used only where the report says so.
- A single trial has no error bar. Where trials were repeated the spread is shown per measurement, and
  no conclusion may rest on a difference below the measured noise floor.
- One DeepSeek HIGH agent planned, built, measured and analysed this study. There is no independent
  model review of the implementation, the oracle or the conclusions.
- A cache hit executes no SQL and therefore has no plan. The saved database operation is what a hit
  avoids; the plans file records the statement it avoids.

