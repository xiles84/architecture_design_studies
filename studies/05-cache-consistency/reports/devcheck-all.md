# Study 05 — external cache throughput and consistency — run `devcheck-all`

Generated 2026-09-21T14:02:43Z from 27 result file(s). This report contains measurements only; conclusions live in
signed analyses under `reports/analyses/`.

| Field | Value |
|---|---|
| Study | 05-cache-consistency |
| Run id | `devcheck-all` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `fd599da9776019d6b117f01d21dcc96077300a8a` |
| `git describe` | `study-05/v1-harness-2-gfd599da` |
| Working tree dirty at start | false |
| Benchmark image | `localhost/cachebench:1` (`8d64eab6e10a1b57080`) |
| Resource framing | `db-only` |
| **Inputs digest** | `6d5b21608914fd00` |
| Cells | 27 |
| Distinct scenarios | 27 |
| Hard TTL | 300 s (probabilistic early expiry p = clamp(1 − remaining/300 s, 0, 1)) |
| Trialling | 1 trial(s) per measurement; throughput is the median |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 27; cells that failed: 27.
  - failed: `pg-single/ctl-stale-invalidation` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/ctl-unchecked-rmw` (no cache backend for "none")
  - failed: `pg-single/legacy-na-memory-aside-relaxed-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/legacy-na-memory-aside-strict-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/legacy-na-memory-through-relaxed-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/legacy-na-memory-through-strict-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/legacy-na-none-none-none-coord` (no cache backend for "none")
  - failed: `pg-single/legacy-na-redis-aside-relaxed-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/legacy-na-redis-aside-relaxed-ext20` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/legacy-na-redis-aside-strict-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/legacy-na-redis-aside-strict-ext20` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/legacy-na-redis-through-relaxed-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/legacy-na-redis-through-strict-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/owned-opt-memory-aside-relaxed-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/owned-opt-memory-aside-strict-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/owned-opt-memory-through-relaxed-coord` (phase faults: fault cache-unavailable could not be reproduced: the outage produced an impossible cache value)
  - failed: `pg-single/owned-opt-memory-through-strict-coord` (phase faults: fault dirty-cache-write-audit-rejects-impossible-version could not be reproduced: the primed read was not fresh, so the audit cannot be calibrated)
  - failed: `pg-single/owned-opt-none-none-none-coord` (no cache backend for "none")
  - failed: `pg-single/owned-opt-redis-aside-relaxed-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/owned-opt-redis-aside-strict-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/owned-opt-redis-through-relaxed-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/owned-opt-redis-through-strict-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/owned-pess-none-none-none-coord` (no cache backend for "none")
  - failed: `pg-single/owned-pess-redis-aside-strict-coord` (phase faults: fault lease-holder-death could not be reproduced: could not take the lease (ok=false err=dial tcp: lookup ads-redis on 10.89.0.1:53: no such host))
  - failed: `pg-single/ref-embedded-locked` (no cache backend for "")
  - failed: `pg-single/ref-normalized-indexed` (no cache backend for "")
  - failed: `pg-single/ref-rollup-trigger` (no cache backend for "")
- Negative controls that fired: 1. Controls that did NOT fire: 1.
  - did not fire: `pg-single/ctl-unchecked-rmw` — the associated correctness claim is **not demonstrated**
- No strict cell recorded a stale-after-ack read.
- Cells with impossible cache values: 3 (**these cells are invalid under every freshness policy**).
  - `pg-single/legacy-na-memory-through-relaxed-coord`: 9 impossible cache value(s)
  - `pg-single/legacy-na-redis-aside-strict-ext20`: 1 impossible cache value(s)
  - `pg-single/owned-opt-memory-through-relaxed-coord`: 4832 impossible cache value(s)
- Highest wrong-read rate recorded: 70.91% of all reads (`pg-single/o:opt:mem:asd:str/cacheable_99_reads`), with 40692 wrong read(s) in 57385.

This list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `6d5b21608914fd00`.

## Scenarios

Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.
The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.

| Topology | Scenario | Model | Ver | Backend | Strategy | Freshness | Writers | Strict policy | Gate | Strict viol. | Impossible | Control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale-invalidation` | legacy | na | memory | aside | relaxed | coord | invalidation-only | aborted | 0 | 0 | fired |
| pg-single | `ctl-unchecked-rmw` | owned | unsafe | none | none | none | coord | invalidation-only | aborted | 0 | 0 | NOT fired |
| pg-single | `legacy-na-memory-aside-relaxed-coord` | legacy | na | memory | aside | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-memory-aside-strict-coord` | legacy | na | memory | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-memory-through-relaxed-coord` | legacy | na | memory | through | relaxed | coord | invalidation-only | aborted | 0 | 9 |  |
| pg-single | `legacy-na-memory-through-strict-coord` | legacy | na | memory | through | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-none-none-none-coord` | legacy | na | none | none | none | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-relaxed-coord` | legacy | na | redis | aside | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-relaxed-ext20` | legacy | na | redis | aside | relaxed | ext20 | authoritative-read | aborted | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-strict-coord` | legacy | na | redis | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-redis-aside-strict-ext20` | legacy | na | redis | aside | strict | ext20 | authoritative-read | aborted | 0 | 1 |  |
| pg-single | `legacy-na-redis-through-relaxed-coord` | legacy | na | redis | through | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `legacy-na-redis-through-strict-coord` | legacy | na | redis | through | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-memory-aside-relaxed-coord` | owned | opt | memory | aside | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-memory-aside-strict-coord` | owned | opt | memory | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-memory-through-relaxed-coord` | owned | opt | memory | through | relaxed | coord | invalidation-only | aborted | 0 | 4832 |  |
| pg-single | `owned-opt-memory-through-strict-coord` | owned | opt | memory | through | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-none-none-none-coord` | owned | opt | none | none | none | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-aside-relaxed-coord` | owned | opt | redis | aside | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-aside-strict-coord` | owned | opt | redis | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-through-relaxed-coord` | owned | opt | redis | through | relaxed | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-opt-redis-through-strict-coord` | owned | opt | redis | through | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-pess-none-none-none-coord` | owned | pess | none | none | none | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `owned-pess-redis-aside-strict-coord` | owned | pess | redis | aside | strict | coord | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `ref-embedded-locked` | ref |  |  |  |  |  | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `ref-normalized-indexed` | ref |  |  |  |  |  | invalidation-only | aborted | 0 | 0 |  |
| pg-single | `ref-rollup-trigger` | ref |  |  |  |  |  | invalidation-only | aborted | 0 | 0 |  |

## Throughput

A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is
not a fast result, and this table is built so it cannot be read as one.

| Topology | Scenario | Phase | Endpoint | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % | wrong reads | wrong % of reads |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | warm_read_only | cacheable | 145k | 0.05 | 0.29 | 3.40 | 58360 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `ctl-stale` | cacheable_99_reads | cacheable | 110k | 0.06 | 0.40 | 107 | 44742 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `ctl-stale` | cacheable_90_reads | cacheable | 55k | 0.01 | 0.19 | 201 | 32950 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `ctl-stale` | hotspot_reads | cacheable | 112k | 0.05 | 0.34 | 191 | 45001 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `ctl-stale` | app_99_reads | total app | 9.4k | 0.01 | 2.85 | 213 | 5090 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `ctl-stale` | app_90_reads | total app | 2.1k | 0.02 | 113 | 207 | 1073 | 0 | 0.0 | 41 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | warm_read_only | cacheable | 132k | 0.05 | 0.31 | 8.21 | 52894 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | cacheable_99_reads | cacheable | 101k | 0.06 | 0.41 | 4.15 | 40524 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | cacheable_90_reads | cacheable | 56k | 0.03 | 0.34 | 214 | 33823 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | hotspot_reads | cacheable | 96k | 0.06 | 0.44 | 13 | 38641 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | app_99_reads | total app | 8.3k | 0.01 | 3.35 | 209 | 4612 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:rel` | app_90_reads | total app | 2.1k | 0.02 | 161 | 218 | 1215 | 0 | 0.0 | 46 | 0.01% |
| pg-single | `l:na:mem:asd:str` | warm_read_only | cacheable | 136k | 0.05 | 0.28 | 106 | 54510 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | cacheable | 97k | 0.06 | 0.40 | 6.87 | 38990 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | cacheable | 63k | 0.01 | 0.27 | 221 | 31654 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | cacheable | 121k | 0.05 | 0.40 | 6.31 | 48707 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:asd:str` | app_99_reads | total app | 14k | 0.01 | 1.55 | 222 | 7765 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:asd:str` | app_90_reads | total app | 2.2k | 0.02 | 158 | 211 | 1228 | 0 | 0.0 | 107116 | 25.44% |
| pg-single | `l:na:mem:wth:rel` | warm_read_only | cacheable | 153k | 0.04 | 0.30 | 5.12 | 61218 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:rel` | cacheable_99_reads | cacheable | 111k | 0.05 | 0.41 | 3.12 | 44782 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:rel` | cacheable_90_reads | cacheable | 34k | 0.01 | 1.07 | 228 | 18430 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:rel` | hotspot_reads | cacheable | 90k | 0.07 | 0.44 | 6.19 | 36199 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:rel` | app_99_reads | total app | 27k | 0.02 | 5.23 | 39 | 10813 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:rel` | app_90_reads | total app | 8.8k | 0.02 | 7.12 | 214 | 4868 | 0 | 0.0 | 8179 | 1.86% |
| pg-single | `l:na:mem:wth:str` | warm_read_only | cacheable | 132k | 0.05 | 0.31 | 2.91 | 53190 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | cacheable | 151k | 0.04 | 0.31 | 4.60 | 60846 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | cacheable | 59k | 0.01 | 0.36 | 226 | 32217 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | cacheable | 105k | 0.04 | 0.38 | 105 | 46638 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:mem:wth:str` | app_99_reads | total app | 14k | 0.01 | 2.41 | 213 | 8518 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:mem:wth:str` | app_90_reads | total app | 3.4k | 0.01 | 17 | 218 | 1847 | 0 | 0.0 | 950 | 0.19% |
| pg-single | `l:na:red:asd:rel` | warm_read_only | cacheable | 594 | 13 | 37 | 41 | 240 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | cacheable_99_reads | cacheable | 349 | 15 | 107 | 108 | 144 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | cacheable_90_reads | cacheable | 599 | 13 | 21 | 23 | 248 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | hotspot_reads | cacheable | 618 | 13 | 21 | 22 | 248 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | app_99_reads | total app | 492 | 11 | 46 | 94 | 230 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel` | app_90_reads | total app | 577 | 13 | 74 | 103 | 234 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | warm_read_only | cacheable | 7.0k | 0.78 | 17 | 26 | 2802 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_99_reads | cacheable | 5.5k | 0.93 | 18 | 34 | 2193 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_90_reads | cacheable | 3.7k | 1.35 | 32 | 41 | 1491 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_reads | cacheable | 4.8k | 1.11 | 28 | 32 | 1994 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | app_99_reads | total app | 5.6k | 0.87 | 13 | 39 | 2334 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:rel:ext` | app_90_reads | total app | 4.3k | 0.88 | 23 | 65 | 1871 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | warm_read_only | cacheable | 543 | 14 | 20 | 21 | 224 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | cacheable | 504 | 14 | 39 | 40 | 216 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | cacheable | 454 | 14 | 47 | 47 | 184 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | hotspot_reads | cacheable | 428 | 16 | 52 | 53 | 176 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_99_reads | total app | 566 | 12 | 91 | 95 | 234 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str` | app_90_reads | total app | 802 | 10 | 18 | 45 | 329 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | warm_read_only | cacheable | 5.2k | 0.97 | 22 | 31 | 2079 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | cacheable_99_reads | cacheable | 5.5k | 0.89 | 22 | 35 | 2222 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | cacheable_90_reads | cacheable | 2.7k | 1.77 | 39 | 41 | 1066 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | hotspot_reads | cacheable | 4.6k | 1.04 | 17 | 39 | 1868 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | app_99_reads | total app | 3.9k | 1.21 | 16 | 73 | 1631 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:asd:str:ext` | app_90_reads | total app | 4.5k | 0.90 | 16 | 40 | 1832 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | warm_read_only | cacheable | 514 | 11 | 52 | 53 | 208 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | cacheable_99_reads | cacheable | 497 | 14 | 46 | 46 | 200 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | cacheable_90_reads | cacheable | 391 | 14 | 54 | 63 | 160 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | hotspot_reads | cacheable | 496 | 15 | 29 | 29 | 208 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | app_99_reads | total app | 686 | 12 | 42 | 42 | 282 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:rel` | app_90_reads | total app | 686 | 11 | 49 | 52 | 284 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | warm_read_only | cacheable | 501 | 12 | 42 | 44 | 208 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | cacheable | 327 | 16 | 100 | 104 | 136 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | cacheable | 691 | 11 | 19 | 22 | 280 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | hotspot_reads | cacheable | 655 | 12 | 17 | 19 | 270 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | app_99_reads | total app | 522 | 13 | 41 | 42 | 214 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `l:na:red:wth:str` | app_90_reads | total app | 527 | 14 | 35 | 72 | 224 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:mem:asd:rel` | warm_read_only | cacheable | 204k | 0.03 | 0.20 | 7.88 | 81783 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:rel` | cacheable_99_reads | cacheable | 111k | 0.06 | 0.41 | 2.27 | 44598 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:rel` | cacheable_90_reads | cacheable | 55k | 0.01 | 0.26 | 225 | 31002 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:rel` | hotspot_reads | cacheable | 95k | 0.06 | 0.43 | 115 | 38126 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:rel` | app_99_reads | total app | 8.4k | 0.01 | 3.33 | 221 | 5048 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:rel` | app_90_reads | total app | 2.2k | 0.02 | 114 | 218 | 1141 | 0 | 0.0 | 118 | 0.03% |
| pg-single | `o:opt:mem:asd:str` | warm_read_only | cacheable | 126k | 0.04 | 0.30 | 173 | 63333 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:asd:str` | cacheable_99_reads | cacheable | 105k | 0.06 | 0.41 | 2.23 | 41961 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:asd:str` | cacheable_90_reads | cacheable | 91k | 0.06 | 0.47 | 34 | 36266 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:asd:str` | hotspot_reads | cacheable | 132k | 0.04 | 0.27 | 183 | 63577 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:asd:str` | app_99_reads | total app | 13k | 0.01 | 2.37 | 225 | 7468 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:asd:str` | app_90_reads | total app | 2.0k | 0.02 | 157 | 218 | 1087 | 0 | 0.0 | 130148 | 26.63% |
| pg-single | `o:opt:mem:wth:rel` | warm_read_only | cacheable | 128k | 0.05 | 0.35 | 8.70 | 51259 | 0 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:rel` | cacheable_99_reads | cacheable | 136k | 0.05 | 0.30 | 4.37 | 54614 | 0 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:rel` | cacheable_90_reads | cacheable | 54k | 0.01 | 0.17 | 208 | 30473 | 0 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:rel` | hotspot_reads | cacheable | 91k | 0.06 | 0.47 | 43 | 36583 | 0 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:rel` | app_99_reads | total app | 31k | 0.01 | 4.34 | 39 | 12364 | 0 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:rel` | app_90_reads | total app | 6.0k | 0.01 | 18 | 165 | 2677 | 6 | 0.0 | 6622 | 1.44% |
| pg-single | `o:opt:mem:wth:str` | warm_read_only | cacheable | 134k | 0.04 | 0.35 | 18 | 53875 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:mem:wth:str` | cacheable_99_reads | cacheable | 109k | 0.06 | 0.41 | 3.10 | 43734 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:mem:wth:str` | cacheable_90_reads | cacheable | 59k | 0.02 | 0.22 | 213 | 31704 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:mem:wth:str` | hotspot_reads | cacheable | 96k | 0.04 | 0.34 | 211 | 41648 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:mem:wth:str` | app_99_reads | total app | 11k | 0.01 | 3.93 | 218 | 5881 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:mem:wth:str` | app_90_reads | total app | 3.1k | 0.01 | 52 | 185 | 1471 | 0 | 0.0 | 1115 | 0.26% |
| pg-single | `o:opt:red:asd:rel` | warm_read_only | cacheable | 405 | 14 | 54 | 54 | 176 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | cacheable_99_reads | cacheable | 686 | 11 | 15 | 16 | 280 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | cacheable_90_reads | cacheable | 607 | 11 | 29 | 29 | 248 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | hotspot_reads | cacheable | 587 | 13 | 19 | 20 | 240 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | app_99_reads | total app | 489 | 14 | 42 | 50 | 204 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:rel` | app_90_reads | total app | 811 | 10 | 22 | 34 | 330 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | warm_read_only | cacheable | 773 | 9.80 | 24 | 26 | 312 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | cacheable_99_reads | cacheable | 432 | 14 | 59 | 64 | 176 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | cacheable_90_reads | cacheable | 443 | 13 | 57 | 62 | 184 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | hotspot_reads | cacheable | 526 | 12 | 56 | 71 | 216 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | app_99_reads | total app | 845 | 9.92 | 16 | 18 | 350 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:asd:str` | app_90_reads | total app | 790 | 9.43 | 50 | 53 | 323 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | warm_read_only | cacheable | 435 | 17 | 46 | 47 | 176 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | cacheable_99_reads | cacheable | 392 | 16 | 54 | 58 | 160 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | cacheable_90_reads | cacheable | 448 | 16 | 48 | 49 | 184 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | hotspot_reads | cacheable | 515 | 15 | 25 | 27 | 208 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | app_99_reads | total app | 610 | 15 | 20 | 21 | 246 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:rel` | app_90_reads | total app | 481 | 13 | 54 | 55 | 204 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | warm_read_only | cacheable | 649 | 11 | 31 | 32 | 264 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | cacheable | 585 | 12 | 40 | 41 | 241 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | cacheable | 416 | 15 | 51 | 53 | 178 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | cacheable | 424 | 15 | 63 | 71 | 175 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_99_reads | total app | 879 | 9.16 | 39 | 49 | 358 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:opt:red:wth:str` | app_90_reads | total app | 569 | 12 | 60 | 67 | 230 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | warm_read_only | cacheable | 664 | 11 | 21 | 23 | 272 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | cacheable_99_reads | cacheable | 560 | 14 | 20 | 23 | 231 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | cacheable_90_reads | cacheable | 529 | 13 | 40 | 44 | 216 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | hotspot_reads | cacheable | 435 | 14 | 58 | 60 | 182 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | app_99_reads | total app | 702 | 12 | 24 | 31 | 289 | 0 | 0.0 | 0 | 0.00% |
| pg-single | `o:pess:red:asd:str` | app_90_reads | total app | 478 | 15 | 40 | 42 | 199 | 0 | 0.0 | 0 | 0.00% |

## Writes (isolated per operation, then the hotspot race)

| Topology | Scenario | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | hotspot_writes | 1.3k | 1.19 | 3.01 | 537 | 0 | 0.0 |
| pg-single | `l:na:mem:asd:rel` | hotspot_writes | 837 | 2.07 | 7.64 | 336 | 0 | 0.0 |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 1.6k | 1.14 | 2.51 | 636 | 0 | 0.0 |
| pg-single | `l:na:mem:wth:rel` | hotspot_writes | 522 | 3.60 | 7.54 | 210 | 0 | 0.0 |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 865 | 2.24 | 3.74 | 347 | 0 | 0.0 |
| pg-single | `l:na:red:asd:rel` | hotspot_writes | 113 | 15 | 33 | 46 | 0 | 0.0 |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_writes | 95.6 | 16 | 100 | 40 | 0 | 0.0 |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 112 | 16 | 29 | 47 | 0 | 0.0 |
| pg-single | `l:na:red:asd:str:ext` | hotspot_writes | 105 | 15 | 56 | 43 | 0 | 0.0 |
| pg-single | `l:na:red:wth:rel` | hotspot_writes | 123 | 14 | 33 | 50 | 0 | 0.0 |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 72.0 | 23 | 63 | 30 | 0 | 0.0 |
| pg-single | `o:opt:mem:asd:rel` | hotspot_writes | 530 | 3.20 | 9.26 | 214 | 0 | 0.0 |
| pg-single | `o:opt:mem:asd:str` | hotspot_writes | 453 | 3.85 | 10 | 182 | 0 | 0.0 |
| pg-single | `o:opt:mem:wth:rel` | hotspot_writes | 611 | 3.07 | 5.94 | 246 | 0 | 0.0 |
| pg-single | `o:opt:mem:wth:str` | hotspot_writes | 398 | 4.35 | 13 | 160 | 0 | 0.0 |
| pg-single | `o:opt:red:asd:rel` | hotspot_writes | 136 | 13 | 39 | 57 | 0 | 0.0 |
| pg-single | `o:opt:red:asd:str` | hotspot_writes | 127 | 14 | 32 | 52 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:rel` | hotspot_writes | 88.9 | 20 | 46 | 37 | 0 | 0.0 |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 69.7 | 26 | 45 | 29 | 0 | 0.0 |
| pg-single | `o:pess:red:asd:str` | hotspot_writes | 55.6 | 27 | 111 | 23 | 0 | 0.0 |

## Wrong-read accounting (per phase)

Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database
on every hit. `ahead` counts reads that returned a later committed state than they required, which is
legitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and
is kept separate from both fresh and wrong.

| Topology | Scenario | Phase | reads | fresh | wrong | wrong % all | hits | wrong % hits | keys | max behind | stale p50 ms | stale p99 ms | streak | TTF p99 ms | concurrent | impossible | causes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | warm_read_only | 70615 | 70615 | 0 | 0.00% | 70575 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | app_99_reads | 4859 | 4838 | 0 | 0.00% | 4798 | 0.00% | 0 | 0 | — | — | 0 | 45 | 21 | 0 | — |
| pg-single | `ctl-stale` | cacheable_99_reads | 60634 | 60613 | 0 | 0.00% | 60572 | 0.00% | 0 | 0 | — | — | 0 | 51 | 21 | 0 | — |
| pg-single | `ctl-stale` | app_90_reads | 61794 | 61738 | 8 | 0.01% | 61592 | 0.01% | 2 | 1 | 4.00 | 5.00 | 4 | 155 | 51 | 0 | other evidenced cause=8 |
| pg-single | `ctl-stale` | cacheable_90_reads | 100435 | 100379 | 8 | 0.01% | 100081 | 0.01% | 2 | 1 | 4.00 | 5.00 | 4 | 322 | 51 | 0 | other evidenced cause=8 |
| pg-single | `ctl-stale` | hotspot_reads | 52597 | 52597 | 0 | 0.00% | 52597 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | hotspot_writes | 52597 | 52597 | 0 | 0.00% | 52597 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `ctl-stale` | instances_1 | 27253 | 27202 | 25 | 0.09% | 26871 | 0.09% | 1 | 1 | 7.00 | 19 | 11 | 434 | 32 | 0 | other evidenced cause=25 |
| pg-single | `ctl-stale` | instances_3 | 2465 | 2430 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 263 | 90 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 18 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | warm_read_only | 62517 | 62517 | 0 | 0.00% | 62477 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | app_99_reads | 5029 | 5010 | 0 | 0.00% | 4971 | 0.00% | 0 | 0 | — | — | 0 | 61 | 19 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | cacheable_99_reads | 55620 | 55601 | 0 | 0.00% | 55559 | 0.00% | 0 | 0 | — | — | 0 | 164 | 19 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | app_90_reads | 56886 | 56824 | 19 | 0.03% | 56685 | 0.03% | 1 | 2 | 11 | 21 | 19 | 294 | 45 | 0 | other evidenced cause=19 |
| pg-single | `l:na:mem:asd:rel` | cacheable_90_reads | 96840 | 96778 | 19 | 0.02% | 96575 | 0.02% | 1 | 2 | 11 | 21 | 19 | 294 | 45 | 0 | other evidenced cause=19 |
| pg-single | `l:na:mem:asd:rel` | hotspot_reads | 47393 | 47393 | 0 | 0.00% | 47393 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | hotspot_writes | 47393 | 47393 | 0 | 0.00% | 47393 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:rel` | instances_1 | 29195 | 29162 | 8 | 0.03% | 28824 | 0.03% | 2 | 1 | 4.00 | 9.00 | 7 | 224 | 28 | 0 | other evidenced cause=8 |
| pg-single | `l:na:mem:asd:rel` | instances_3 | 3126 | 3082 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 306 | 138 | 0 | — |
| pg-single | `l:na:mem:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 12 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | warm_read_only | 64013 | 64013 | 0 | 0.00% | 63973 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | app_99_reads | 7144 | 3299 | 3824 | 53.53% | 7057 | 54.19% | 9 | 1 | 71 | 1454 | 373 | 135 | 34 | 0 | other evidenced cause=3824 |
| pg-single | `l:na:mem:asd:str` | cacheable_99_reads | 54677 | 24925 | 29731 | 54.38% | 54590 | 54.46% | 9 | 1 | 450 | 1418 | 16701 | 135 | 34 | 0 | other evidenced cause=29731 |
| pg-single | `l:na:mem:asd:str` | app_90_reads | 55846 | 25750 | 30047 | 53.80% | 55622 | 54.02% | 9 | 1 | 448 | 1419 | 16712 | 316 | 70 | 0 | other evidenced cause=30047 |
| pg-single | `l:na:mem:asd:str` | cacheable_90_reads | 93877 | 63745 | 30083 | 32.05% | 93563 | 32.15% | 9 | 1 | 449 | 1435 | 16712 | 530 | 70 | 0 | other evidenced cause=30083 |
| pg-single | `l:na:mem:asd:str` | hotspot_reads | 57432 | 57432 | 0 | 0.00% | 57432 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | hotspot_writes | 57432 | 57432 | 0 | 0.00% | 57432 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:asd:str` | instances_1 | 27604 | 14148 | 13431 | 48.66% | 27221 | 49.34% | 8 | 2 | 603 | 904 | 9205 | 395 | 39 | 0 | other evidenced cause=13431 |
| pg-single | `l:na:mem:asd:str` | instances_3 | 2820 | 2772 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 256 | 134 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | warm_read_only | 81406 | 81406 | 0 | 0.00% | 81366 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | app_99_reads | 11832 | 10572 | 1260 | 10.65% | 11832 | 10.65% | 11 | 3 | 7.00 | 962 | 71 | 31 | 17 | 0 | other evidenced cause=1260 |
| pg-single | `l:na:mem:wth:rel` | cacheable_99_reads | 66829 | 65569 | 1260 | 1.89% | 66829 | 1.89% | 11 | 3 | 7.00 | 962 | 71 | 32 | 17 | 0 | other evidenced cause=1260 |
| pg-single | `l:na:mem:wth:rel` | app_90_reads | 71794 | 69243 | 2549 | 3.55% | 71751 | 3.55% | 12 | 4 | 5.00 | 525 | 71 | 51 | 42 | 0 | other evidenced cause=2549 |
| pg-single | `l:na:mem:wth:rel` | cacheable_90_reads | 94344 | 91793 | 2549 | 2.70% | 94027 | 2.71% | 12 | 4 | 5.00 | 525 | 71 | 152 | 42 | 0 | other evidenced cause=2549 |
| pg-single | `l:na:mem:wth:rel` | hotspot_reads | 44369 | 44369 | 0 | 0.00% | 44369 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | hotspot_writes | 44369 | 44369 | 0 | 0.00% | 44369 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:rel` | instances_1 | 22068 | 21493 | 561 | 2.54% | 21672 | 2.59% | 9 | 4 | 6.00 | 168 | 27 | 160 | 14 | 9 | other evidenced cause=561 |
| pg-single | `l:na:mem:wth:rel` | instances_3 | 2208 | 2185 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 208 | 86 | 0 | — |
| pg-single | `l:na:mem:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | warm_read_only | 65651 | 65651 | 0 | 0.00% | 65611 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | app_99_reads | 8860 | 8671 | 164 | 1.85% | 8792 | 1.87% | 8 | 1 | 13 | 1032 | 16 | 153 | 36 | 0 | other evidenced cause=164 |
| pg-single | `l:na:mem:wth:str` | cacheable_99_reads | 80120 | 79931 | 164 | 0.20% | 80052 | 0.20% | 8 | 1 | 13 | 1032 | 16 | 190 | 36 | 0 | other evidenced cause=164 |
| pg-single | `l:na:mem:wth:str` | app_90_reads | 82015 | 81685 | 258 | 0.31% | 81830 | 0.32% | 8 | 2 | 7.00 | 1032 | 16 | 190 | 99 | 0 | other evidenced cause=258 |
| pg-single | `l:na:mem:wth:str` | cacheable_90_reads | 121949 | 121619 | 258 | 0.21% | 121427 | 0.21% | 8 | 2 | 7.00 | 1032 | 16 | 210 | 99 | 0 | other evidenced cause=258 |
| pg-single | `l:na:mem:wth:str` | hotspot_reads | 53150 | 53150 | 0 | 0.00% | 53150 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | hotspot_writes | 53150 | 53150 | 0 | 0.00% | 53150 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:mem:wth:str` | instances_1 | 29442 | 29296 | 106 | 0.36% | 28906 | 0.37% | 5 | 2 | 4.00 | 39 | 8 | 218 | 73 | 0 | other evidenced cause=106 |
| pg-single | `l:na:mem:wth:str` | instances_3 | 2969 | 2913 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 201 | 140 | 0 | — |
| pg-single | `l:na:red:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | warm_read_only | 336 | 336 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | app_99_reads | 236 | 231 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 67 | 5 | 0 | — |
| pg-single | `l:na:red:asd:rel` | cacheable_99_reads | 436 | 431 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 70 | 5 | 0 | — |
| pg-single | `l:na:red:asd:rel` | app_90_reads | 642 | 612 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 94 | 32 | 0 | — |
| pg-single | `l:na:red:asd:rel` | cacheable_90_reads | 962 | 932 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 94 | 32 | 0 | — |
| pg-single | `l:na:red:asd:rel` | hotspot_reads | 304 | 304 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | hotspot_writes | 304 | 304 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel` | instances_1 | 609 | 569 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 596 | 45 | 0 | — |
| pg-single | `l:na:red:asd:rel` | instances_3 | 511 | 478 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 51 | 37 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | warm_read_only | 3517 | 3517 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | app_99_reads | 2486 | 2473 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 14 | 28 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_99_reads | 5133 | 5120 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 19 | 28 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | app_90_reads | 6894 | 6829 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 73 | 157 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | cacheable_90_reads | 8770 | 8705 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 145 | 157 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_reads | 2430 | 2430 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | hotspot_writes | 2430 | 2430 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | instances_1 | 2611 | 2577 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 237 | 82 | 0 | — |
| pg-single | `l:na:red:asd:rel:ext` | instances_3 | 2013 | 1998 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 58 | 47 | 0 | — |
| pg-single | `l:na:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | warm_read_only | 294 | 294 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_99_reads | 263 | 262 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 41 | 1 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_99_reads | 535 | 534 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 41 | 1 | 0 | — |
| pg-single | `l:na:red:asd:str` | app_90_reads | 864 | 830 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 130 | 37 | 0 | — |
| pg-single | `l:na:red:asd:str` | cacheable_90_reads | 1091 | 1057 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 130 | 37 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_reads | 240 | 240 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | hotspot_writes | 240 | 240 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_1 | 671 | 641 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 538 | 37 | 0 | — |
| pg-single | `l:na:red:asd:str` | instances_3 | 489 | 458 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 154 | 44 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | warm_read_only | 2455 | 2455 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | app_99_reads | 1746 | 1737 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 9.00 | 14 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cacheable_99_reads | 4627 | 4618 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 13 | 14 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | app_90_reads | 6242 | 6193 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 115 | 90 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | cacheable_90_reads | 7784 | 7735 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 248 | 90 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | hotspot_reads | 2149 | 2149 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | hotspot_writes | 2149 | 2149 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | instances_1 | 4188 | 4163 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 136 | 70 | 0 | — |
| pg-single | `l:na:red:asd:str:ext` | instances_3 | 3662 | 3623 | 0 | 0.00% | 0 | 0.00% | 1 | 0 | — | — | 0 | 93 | 74 | 1 | — |
| pg-single | `l:na:red:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | warm_read_only | 320 | 320 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | app_99_reads | 299 | 298 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2.00 | 1 | 0 | — |
| pg-single | `l:na:red:wth:rel` | cacheable_99_reads | 579 | 578 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 2.00 | 1 | 0 | — |
| pg-single | `l:na:red:wth:rel` | app_90_reads | 839 | 801 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 213 | 39 | 0 | — |
| pg-single | `l:na:red:wth:rel` | cacheable_90_reads | 1055 | 1017 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 213 | 39 | 0 | — |
| pg-single | `l:na:red:wth:rel` | hotspot_reads | 232 | 232 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | hotspot_writes | 232 | 232 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:rel` | instances_1 | 585 | 551 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 367 | 37 | 0 | — |
| pg-single | `l:na:red:wth:rel` | instances_3 | 477 | 457 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 415 | 24 | 0 | — |
| pg-single | `l:na:red:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | warm_read_only | 296 | 296 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | app_99_reads | 234 | 233 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 25 | 2 | 0 | — |
| pg-single | `l:na:red:wth:str` | cacheable_99_reads | 418 | 417 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 25 | 2 | 0 | — |
| pg-single | `l:na:red:wth:str` | app_90_reads | 627 | 597 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 71 | 32 | 0 | — |
| pg-single | `l:na:red:wth:str` | cacheable_90_reads | 979 | 949 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 75 | 32 | 0 | — |
| pg-single | `l:na:red:wth:str` | hotspot_reads | 350 | 350 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | hotspot_writes | 350 | 350 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_1 | 568 | 543 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 91 | 29 | 0 | — |
| pg-single | `l:na:red:wth:str` | instances_3 | 393 | 378 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 76 | 22 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 15 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | warm_read_only | 94364 | 94364 | 0 | 0.00% | 94324 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | app_99_reads | 5948 | 5914 | 16 | 0.27% | 5888 | 0.27% | 1 | 1 | 5.00 | 6.00 | 16 | 17 | 20 | 0 | other evidenced cause=16 |
| pg-single | `o:opt:mem:asd:rel` | cacheable_99_reads | 60196 | 60162 | 16 | 0.03% | 60135 | 0.03% | 1 | 1 | 5.00 | 6.00 | 16 | 48 | 20 | 0 | other evidenced cause=16 |
| pg-single | `o:opt:mem:asd:rel` | app_90_reads | 61289 | 61214 | 32 | 0.05% | 61081 | 0.05% | 2 | 1 | 6.00 | 23 | 16 | 133 | 47 | 0 | other evidenced cause=32 |
| pg-single | `o:opt:mem:asd:rel` | cacheable_90_reads | 98220 | 98145 | 32 | 0.03% | 97883 | 0.03% | 2 | 1 | 6.00 | 23 | 16 | 270 | 47 | 0 | other evidenced cause=32 |
| pg-single | `o:opt:mem:asd:rel` | hotspot_reads | 45831 | 45831 | 0 | 0.00% | 45831 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | hotspot_writes | 45831 | 45831 | 0 | 0.00% | 45831 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:rel` | instances_1 | 24267 | 24219 | 22 | 0.09% | 23912 | 0.09% | 1 | 1 | 6.00 | 15 | 15 | 249 | 29 | 0 | other evidenced cause=22 |
| pg-single | `o:opt:mem:asd:rel` | instances_3 | 7243 | 7214 | 0 | 0.00% | 7021 | 0.00% | 0 | 0 | — | — | 0 | 284 | 34 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 13 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | warm_read_only | 76360 | 76360 | 0 | 0.00% | 76320 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | app_99_reads | 7094 | 2462 | 4616 | 65.07% | 7007 | 65.88% | 9 | 2 | 58 | 1616 | 502 | 205 | 29 | 0 | other evidenced cause=4616 |
| pg-single | `o:opt:mem:asd:str` | cacheable_99_reads | 57385 | 16677 | 40692 | 70.91% | 57298 | 71.02% | 9 | 2 | 498 | 1386 | 17673 | 205 | 29 | 0 | other evidenced cause=40692 |
| pg-single | `o:opt:mem:asd:str` | app_90_reads | 58566 | 17500 | 41021 | 70.04% | 58344 | 70.31% | 9 | 2 | 497 | 1385 | 17693 | 286 | 66 | 0 | other evidenced cause=41021 |
| pg-single | `o:opt:mem:asd:str` | cacheable_90_reads | 104057 | 62991 | 41021 | 39.42% | 103833 | 39.51% | 9 | 2 | 497 | 1385 | 17693 | 290 | 66 | 0 | other evidenced cause=41021 |
| pg-single | `o:opt:mem:asd:str` | hotspot_reads | 76580 | 76580 | 0 | 0.00% | 76580 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | hotspot_writes | 76580 | 76580 | 0 | 0.00% | 76580 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:asd:str` | instances_1 | 23440 | 20618 | 2798 | 11.94% | 23196 | 12.06% | 7 | 2 | 526 | 804 | 1420 | 245 | 36 | 0 | other evidenced cause=2798 |
| pg-single | `o:opt:mem:asd:str` | instances_3 | 8467 | 8439 | 0 | 0.00% | 8223 | 0.00% | 0 | 0 | — | — | 0 | 270 | 40 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 9 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | warm_read_only | 63620 | 63620 | 0 | 0.00% | 63580 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | app_99_reads | 11639 | 10297 | 1342 | 11.53% | 11639 | 11.53% | 10 | 2 | 12 | 2977 | 44 | 63 | 12 | 0 | other evidenced cause=1342 |
| pg-single | `o:opt:mem:wth:rel` | cacheable_99_reads | 74016 | 72674 | 1342 | 1.81% | 74016 | 1.81% | 10 | 2 | 12 | 2977 | 44 | 63 | 12 | 0 | other evidenced cause=1342 |
| pg-single | `o:opt:mem:wth:rel` | app_90_reads | 76714 | 74869 | 1753 | 2.29% | 76705 | 2.29% | 11 | 2 | 9.00 | 2976 | 44 | 89 | 20 | 91 | other evidenced cause=1753 |
| pg-single | `o:opt:mem:wth:rel` | cacheable_90_reads | 115301 | 110956 | 1753 | 1.52% | 115158 | 1.52% | 11 | 2 | 9.00 | 2976 | 44 | 89 | 20 | 2591 | other evidenced cause=1753 |
| pg-single | `o:opt:mem:wth:rel` | hotspot_reads | 42387 | 42387 | 0 | 0.00% | 42387 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | hotspot_writes | 42387 | 42387 | 0 | 0.00% | 42387 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:rel` | stampede | 128 | 112 | 0 | 0.00% | 120 | 0.00% | 1 | 0 | — | — | 0 | — | 0 | 16 | — |
| pg-single | `o:opt:mem:wth:rel` | instances_1 | 26195 | 24109 | 432 | 1.65% | 25862 | 1.67% | 8 | 2 | 5.00 | 206 | 18 | 95 | 7 | 1650 | other evidenced cause=432 |
| pg-single | `o:opt:mem:wth:rel` | instances_3 | 7992 | 7476 | 0 | 0.00% | 7826 | 0.00% | 1 | 0 | — | — | 0 | 427 | 43 | 484 | — |
| pg-single | `o:opt:mem:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 14 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | warm_read_only | 65960 | 65960 | 0 | 0.00% | 65920 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | app_99_reads | 5563 | 5344 | 193 | 3.47% | 5504 | 3.51% | 6 | 1 | 19 | 1001 | 21 | 87 | 29 | 0 | other evidenced cause=193 |
| pg-single | `o:opt:mem:wth:str` | cacheable_99_reads | 58582 | 58363 | 193 | 0.33% | 58523 | 0.33% | 6 | 1 | 19 | 1001 | 21 | 106 | 29 | 0 | other evidenced cause=193 |
| pg-single | `o:opt:mem:wth:str` | app_90_reads | 60022 | 59654 | 318 | 0.53% | 59885 | 0.53% | 7 | 1 | 11 | 1001 | 21 | 199 | 68 | 0 | other evidenced cause=318 |
| pg-single | `o:opt:mem:wth:str` | cacheable_90_reads | 100105 | 99737 | 318 | 0.32% | 99899 | 0.32% | 7 | 1 | 11 | 1001 | 21 | 199 | 68 | 0 | other evidenced cause=318 |
| pg-single | `o:opt:mem:wth:str` | hotspot_reads | 49243 | 49243 | 0 | 0.00% | 49243 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | hotspot_writes | 49243 | 49243 | 0 | 0.00% | 49243 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 120 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:mem:wth:str` | instances_1 | 33336 | 33223 | 93 | 0.28% | 33012 | 0.28% | 6 | 1 | 8.00 | 263 | 15 | 265 | 28 | 0 | other evidenced cause=93 |
| pg-single | `o:opt:mem:wth:str` | instances_3 | 7376 | 7344 | 0 | 0.00% | 7203 | 0.00% | 0 | 0 | — | — | 0 | 298 | 40 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | warm_read_only | 280 | 280 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | app_99_reads | 228 | 226 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 17 | 2 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | cacheable_99_reads | 566 | 564 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 17 | 2 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | app_90_reads | 877 | 844 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 130 | 42 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | cacheable_90_reads | 1197 | 1164 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 366 | 42 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | hotspot_reads | 312 | 312 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | hotspot_writes | 312 | 312 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | instances_1 | 505 | 478 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 652 | 28 | 0 | — |
| pg-single | `o:opt:red:asd:rel` | instances_3 | 444 | 413 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 91 | 34 | 0 | — |
| pg-single | `o:opt:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | warm_read_only | 416 | 416 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | app_99_reads | 386 | 381 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 69 | 6 | 0 | — |
| pg-single | `o:opt:red:asd:str` | cacheable_99_reads | 626 | 621 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 69 | 6 | 0 | — |
| pg-single | `o:opt:red:asd:str` | app_90_reads | 915 | 872 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 72 | 47 | 0 | — |
| pg-single | `o:opt:red:asd:str` | cacheable_90_reads | 1147 | 1104 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 228 | 47 | 0 | — |
| pg-single | `o:opt:red:asd:str` | hotspot_reads | 288 | 288 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | hotspot_writes | 288 | 288 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:asd:str` | instances_1 | 749 | 724 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 52 | 29 | 0 | — |
| pg-single | `o:opt:red:asd:str` | instances_3 | 660 | 615 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 271 | 48 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | warm_read_only | 263 | 263 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | app_99_reads | 247 | 246 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 75 | 1 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | cacheable_99_reads | 463 | 462 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 75 | 1 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | app_90_reads | 650 | 632 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 75 | 19 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | cacheable_90_reads | 872 | 854 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 76 | 19 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | hotspot_reads | 256 | 256 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | hotspot_writes | 256 | 256 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | instances_1 | 511 | 483 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 340 | 30 | 0 | — |
| pg-single | `o:opt:red:wth:rel` | instances_3 | 542 | 518 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 89 | 26 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | warm_read_only | 367 | 367 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_99_reads | 364 | 352 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 10 | 12 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_99_reads | 685 | 673 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 10 | 12 | 0 | — |
| pg-single | `o:opt:red:wth:str` | app_90_reads | 895 | 860 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 52 | 37 | 0 | — |
| pg-single | `o:opt:red:wth:str` | cacheable_90_reads | 1105 | 1070 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 52 | 37 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_reads | 231 | 231 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | hotspot_writes | 231 | 231 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_1 | 651 | 629 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 233 | 23 | 0 | — |
| pg-single | `o:opt:red:wth:str` | instances_3 | 572 | 539 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 242 | 42 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cold_fill | 40 | 40 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | warm_read_only | 368 | 368 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | app_99_reads | 316 | 314 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 29 | 2 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cacheable_99_reads | 611 | 609 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 29 | 2 | 0 | — |
| pg-single | `o:pess:red:asd:str` | app_90_reads | 796 | 756 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 86 | 44 | 0 | — |
| pg-single | `o:pess:red:asd:str` | cacheable_90_reads | 1076 | 1036 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 116 | 44 | 0 | — |
| pg-single | `o:pess:red:asd:str` | hotspot_reads | 246 | 246 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | hotspot_writes | 246 | 246 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | stampede | 128 | 128 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | — | 0 | 0 | — |
| pg-single | `o:pess:red:asd:str` | instances_1 | 532 | 514 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 165 | 23 | 0 | — |
| pg-single | `o:pess:red:asd:str` | instances_3 | 580 | 555 | 0 | 0.00% | 0 | 0.00% | 0 | 0 | — | — | 0 | 367 | 32 | 0 | — |

## Probabilistic early expiration: observed vs specified

The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical
check on the implementation, alongside the unit tests that pin the same ages with a fake clock.

*No cache hits were observed in this run, so there is no expiry evidence to report.*

## Leases and stampede prevention

| Topology | Scenario | acquired | contended | waits | wait p99 ms | timeouts | fallbacks | duplicate fills | lease duration | calibration |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `ctl-rmw` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:mem:asd:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:mem:wth:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:mem:wth:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:asd:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:asd:rel:ext` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:asd:str:ext` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:wth:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `l:na:red:wth:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:mem:asd:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:mem:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:mem:wth:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:mem:wth:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:asd:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:wth:rel` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:opt:red:wth:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `o:pess:red:asd:str` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `ref-emb` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `ref-norm` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |
| pg-single | `ref-roll` | 0 | 0 | 0 | — | 0 | 0 | 0 | — | — |

| Topology | Scenario | readers | keys | elapsed ms | db loads | loads/key | fallbacks | duplicate fills | wrong | impossible |
|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:rel` | 16 | 8 | 211.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | 16 | 8 | 210.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:rel` | 16 | 8 | 213.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:str` | 16 | 8 | 211.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel` | 16 | 8 | 137.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel:ext` | 16 | 8 | 134.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str` | 16 | 8 | 131.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str:ext` | 16 | 8 | 50.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:rel` | 16 | 8 | 160.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | 16 | 8 | 181.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:rel` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:str` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:wth:rel` | 16 | 8 | 212.0 | 8 | 1.00 | 0 | 0 | 0 | 16 |
| pg-single | `o:opt:mem:wth:str` | 16 | 8 | 223.0 | 8 | 1.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:rel` | 16 | 8 | 119.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:str` | 16 | 8 | 133.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:rel` | 16 | 8 | 171.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | 16 | 8 | 136.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |
| pg-single | `o:pess:red:asd:str` | 16 | 8 | 168.0 | 0 | 0.00 | 0 | 0 | 0 | 0 |

## One versus three logical application instances (equal total workers)

| Topology | Scenario | instances | workers | backend | app ops/s | spread % | cacheable ops/s | bypass reads | wrong | impossible | note |
|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 1 | 8 | memory | 1.8k | 0.0 | 37k | 0 | 25 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `ctl-stale` | 3 | 8 | memory | 2.6k | 0.0 | 2.8k | 2465 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:asd:rel` | 1 | 8 | memory | 2.1k | 0.0 | 42k | 0 | 8 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:rel` | 3 | 8 | memory | 3.3k | 0.0 | 3.7k | 3126 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:asd:str` | 1 | 8 | memory | 2.2k | 0.0 | 44k | 0 | 13431 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:asd:str` | 3 | 8 | memory | 3.8k | 0.0 | 2.9k | 2820 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:wth:rel` | 1 | 8 | memory | 6.7k | 0.0 | 28k | 0 | 561 | 9 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:rel` | 3 | 8 | memory | 2.6k | 0.0 | 2.6k | 2208 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:mem:wth:str` | 1 | 8 | memory | 5.3k | 0.0 | 40k | 0 | 106 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:mem:wth:str` | 3 | 8 | memory | 3.7k | 0.0 | 3.1k | 2969 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `l:na:red:asd:rel` | 1 | 8 | redis | 783 | 0.0 | 640 | 609 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel` | 3 | 8 | redis | 746 | 0.0 | 408 | 511 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel:ext` | 1 | 8 | redis | 3.1k | 0.0 | 2.6k | 2611 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:rel:ext` | 3 | 8 | redis | 1.5k | 0.0 | 2.9k | 2013 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str` | 1 | 8 | redis | 795 | 0.0 | 734 | 671 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str` | 3 | 8 | redis | 749 | 0.0 | 352 | 489 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str:ext` | 1 | 8 | redis | 4.6k | 0.0 | 4.6k | 4188 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:asd:str:ext` | 3 | 8 | redis | 4.7k | 0.0 | 3.0k | 3662 | 0 | 1 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:rel` | 1 | 8 | redis | 722 | 0.0 | 587 | 585 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:rel` | 3 | 8 | redis | 513 | 0.0 | 509 | 477 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 1 | 8 | redis | 620 | 0.0 | 604 | 568 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `l:na:red:wth:str` | 3 | 8 | redis | 478 | 0.0 | 310 | 393 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:rel` | 1 | 8 | memory | 2.1k | 0.0 | 38k | 0 | 22 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:rel` | 3 | 8 | memory | 1.3k | 0.0 | 14k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:asd:str` | 1 | 8 | memory | 2.1k | 0.0 | 33k | 0 | 2798 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:asd:str` | 3 | 8 | memory | 1.5k | 0.0 | 15k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:wth:rel` | 1 | 8 | memory | 6.9k | 0.0 | 30k | 0 | 432 | 1650 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:wth:rel` | 3 | 8 | memory | 1.6k | 0.0 | 12k | 0 | 0 | 484 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:mem:wth:str` | 1 | 8 | memory | 2.9k | 0.0 | 57k | 0 | 93 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:mem:wth:str` | 3 | 8 | memory | 1.6k | 0.0 | 14k | 0 | 0 | 0 | separate in-process LRUs and separate leases; a write through one instance is invisible to the other two |
| pg-single | `o:opt:red:asd:rel` | 1 | 8 | redis | 642 | 0.0 | 470 | 505 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:rel` | 3 | 8 | redis | 483 | 0.0 | 545 | 444 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:str` | 1 | 8 | redis | 826 | 0.0 | 846 | 749 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:asd:str` | 3 | 8 | redis | 821 | 0.0 | 636 | 660 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:rel` | 1 | 8 | redis | 685 | 0.0 | 453 | 511 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:rel` | 3 | 8 | redis | 519 | 0.0 | 621 | 542 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 1 | 8 | redis | 842 | 0.0 | 671 | 651 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:opt:red:wth:str` | 3 | 8 | redis | 696 | 0.0 | 587 | 572 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:pess:red:asd:str` | 1 | 8 | redis | 564 | 0.0 | 601 | 532 | 0 | 0 | shared store: invalidation is visible to every instance |
| pg-single | `o:pess:red:asd:str` | 3 | 8 | redis | 613 | 0.0 | 651 | 580 | 0 | 0 | shared store: invalidation is visible to every instance |

## Fault injection and controls

A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly
has not been fast.

| Topology | Scenario | Fault | reproduced | wrong observed | impossible observed | correctness |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `ctl-stale` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `ctl-stale` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `ctl-stale` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `ctl-stale` | cache-aside-fill-races-writer | true | 1 | 0 | a fill and a write on donor 4 released together; 1 wrong and 0 impossible reads |
| pg-single | `ctl-stale` | suppressed-invalidation-produces-a-stale-read | true | 1 | 0 | one suppressed post-commit invalidation produced 1 stale-after-ack read(s) that the oracle detected |
| pg-single | `ctl-stale` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `l:na:mem:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:mem:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:rel` | cache-aside-fill-races-writer | true | 1 | 0 | a fill and a write on donor 4 released together; 1 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `l:na:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `l:na:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `l:na:mem:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `l:na:mem:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:rel` | cache-aside-fill-races-writer | true | 1 | 0 | a fill and a write on donor 4 released together; 1 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:rel` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `l:na:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `l:na:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `l:na:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `l:na:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `l:na:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `l:na:red:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:rel` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `l:na:red:asd:rel` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `l:na:red:asd:rel:ext` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:rel:ext` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `l:na:red:asd:rel:ext` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `l:na:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:asd:str` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `l:na:red:asd:str:ext` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:asd:str:ext` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:asd:str:ext` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `l:na:red:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:wth:rel` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `l:na:red:wth:rel` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `l:na:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `l:na:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `l:na:red:wth:str` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `o:opt:mem:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:asd:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:mem:asd:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:asd:rel` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:asd:rel` | cache-aside-fill-races-writer | true | 1 | 0 | a fill and a write on donor 4 released together; 1 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:asd:rel` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `o:opt:mem:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:mem:asd:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:asd:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:asd:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:asd:str` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `o:opt:mem:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:wth:rel` | cache-update-failure-after-commit | true | 1 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned stale and was counted |
| pg-single | `o:opt:mem:wth:rel` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:wth:rel` | cache-unavailable | false | 0 | 2 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 2 impossible) |
| pg-single | `o:opt:mem:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (hit -> fill); the next read returned fresh |
| pg-single | `o:opt:mem:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:mem:wth:str` | lease-holder-death | true | 0 | 0 | an abandoned lease blocked a second holder and was stealable after 100ms; no manual intervention required |
| pg-single | `o:opt:mem:wth:str` | cache-unavailable | true | 0 | 0 | 20 reads during the outage all returned a committed state (20 bypassed the cache, 0 wrong, 0 impossible) |
| pg-single | `o:opt:mem:wth:str` | cache-aside-fill-races-writer | true | 0 | 0 | a fill and a write on donor 4 released together; 0 wrong and 0 impossible reads |
| pg-single | `o:opt:mem:wth:str` | dirty-cache-write-audit-rejects-impossible-version | false | 0 | 0 | pending |
| pg-single | `o:opt:red:asd:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `o:opt:red:asd:rel` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `o:opt:red:asd:rel` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `o:opt:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `o:opt:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:asd:str` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `o:opt:red:wth:rel` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `o:opt:red:wth:rel` | cache-update-failure-after-commit | true | 0 | 0 | relaxed: the post-commit cache failure left the previous committed value in place; the read returned fresh and was counted (no stale read observed in this single draw; the churn phases carry the rate) |
| pg-single | `o:opt:red:wth:rel` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `o:opt:red:wth:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `o:opt:red:wth:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:opt:red:wth:str` | lease-holder-death | false | 0 | 0 | pending |
| pg-single | `o:pess:red:asd:str` | db-rollback-after-pre-write-tombstone | true | 0 | 0 | tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (bypass -> bypass); the next read returned fresh |
| pg-single | `o:pess:red:asd:str` | cache-update-failure-after-commit | true | 0 | 0 | strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned fresh |
| pg-single | `o:pess:red:asd:str` | lease-holder-death | false | 0 | 0 | pending |

## Ledger replay (audit)

*No audit phase ran in this run.*

## Cache accounting

| Topology | Scenario | backend | capacity B | resident B | items | evictions | fills | publishes | fenced | publish failures | invalidations | tombstones | version validations | bypass reads | external writes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ctl-rmw` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:mem:wth:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:db` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:rel:ext` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:asd:str:ext` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `l:na:red:wth:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:wth:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:mem:wth:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:db` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:rel` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:opt:red:wth:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:pess:db` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `o:pess:red:asd:str` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-emb` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-norm` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| pg-single | `ref-roll` | — | — | — | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

## Client and resource conditions

CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong
to no scenario.

| Topology | Scenario | client periods | client throttled | client throttled ms | framing | workers |
|---|---|---|---|---|---|---|
| pg-single | `ctl-stale` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ctl-rmw` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:asd:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:mem:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:db` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:asd:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:asd:rel:ext` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:asd:str:ext` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `l:na:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:asd:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:mem:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:db` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:asd:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:rel` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:opt:red:wth:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:pess:db` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `o:pess:red:asd:str` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ref-emb` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ref-norm` | 0 | 0 | 0 | db-only | 8 |
| pg-single | `ref-roll` | 0 | 0 | 0 | db-only | 8 |

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
| pg-single | `ctl-rmw` | false | 0 |  |  |
| pg-single | `l:na:mem:asd:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:asd:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:rel` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:mem:wth:str` | true | 13 | r_portal_person r_portal_recent r_charity_recent |  |
| pg-single | `l:na:db` | false | 0 |  |  |
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
| pg-single | `o:opt:db` | false | 0 |  |  |
| pg-single | `o:opt:red:asd:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:asd:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:wth:rel` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:opt:red:wth:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `o:pess:db` | false | 0 |  |  |
| pg-single | `o:pess:red:asd:str` | true | 14 | r_portal_person r_portal_recent r_charity_recent a_outbox_orphans |  |
| pg-single | `ref-emb` | false | 0 |  |  |
| pg-single | `ref-norm` | false | 0 |  |  |
| pg-single | `ref-roll` | false | 0 |  |  |

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

