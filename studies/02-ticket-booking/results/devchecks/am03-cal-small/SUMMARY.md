# AM-03.10 calibration (dev check, not a measured run) — LOW, Claude Sonnet 5, 2026-09-20

`small`, `-conns 8 -duration 5s -warmup 2s`, tag `study-02/v2.2-ledger-attribution` (`5e3a3f8`).
Engine cells: PostgreSQL 17.11 `pg-single`; YugabyteDB 2025.2.6 `yb-cluster3` (3 x 2 CPU / 3 GiB).

## Cell wall times and load (`-phases verify,explain,read`, yb-cluster3)

| Study | Design | Wall | Load (copy / index) | Gate |
|---|---|---:|---|---|
| 02 | P3 | 157 s | 23.7 s (4.4 / 16.8) | 46/46 |
| 02 | X1 | 245 s | 34.2 s (5.5 / 25.2) | 48/48 |
| 03 | S1 | 223 s | 22.3 s (5.6 / 11.6) | 65/65 |
| 03 | L2 | 212 s | 16.3 s (2.6 / 7.9) | 65/65 |

## Report throughput at 5 s (executions = ops/s x 5)

| Report | P3 | X1 |
|---|---:|---:|
| r01 trailing / historical | 844 / 1161 | 2428 / 2081 |
| r02 | 902 | 937 |
| r03 trailing / historical | **7.0 / 6.9 (35 exec.)** | 34 / 38 |
| r04 | 43 | 64 |
| r05 trailing / historical | unanswerable | 3365 / 3605 |

Study 03 (S1 / L2): r01 3553 / 2474; r02 3382 / 992; r03 1484-1658 / 1647-1738;
r04 1512 / 1187; **r05 50-55 / 52-53 (~260 exec.)**. No errors anywhere.

## Race generator capacity (pg-single, 128 buyers, `verify,race`, 1 trial)

P3 32 s wall, X1 58 s wall, both `errors=0`; client CPU throttling **0 of 0 periods**
(`client_cpu.race#1.nr_periods=0`). X1 sells at roughly half of P3's rate from the 10-seat tier
up (e.g. 100 000 seats: 6376 vs 3009 sold/s) — the ledger's cost is already visible.

## Sizing rules applied (AM-03.11)

- **R1:** study 02's reports run needs **10 s** — P3's `r03` gives 35 executions at 5 s, under
  100 (70 at 10 s, above the 50 floor, so it is a result and a named weakness, not a blocker).
  Study 03 stays at **5 s** (slowest report ~260 executions).
- **R2:** client throttled 0% and no errors at 128 buyers, so **B2 = 128**.
- **R3 estimate** (extrapolated from these cells; yb-cluster3 read cells ~4-5 min, 10 s adds
  ~80 s): 3b-1 ~2.6 h, 3b-4 ~2.0 h, 3b-2 ~1.7 h, 3b-3 ~1.7 h, **total ~8 h < 10 h guard**.
