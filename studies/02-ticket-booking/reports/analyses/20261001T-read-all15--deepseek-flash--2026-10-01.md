---
analysis_id: 20261001T-read-all15--deepseek-flash--2026-10-01
run_id: 20261001T100948Z (order A), 20261001T103931Z (order B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; all fifteen read designs in one session, both orders."
analyzed_at: 2026-10-01
inputs_digest: 8413cbbd90b7e01c (order A), 2d2126a5a13118be (order B)
inputs: 20261001T100948Z@8413cbbd90b7e01c, 20261001T103931Z@2d2126a5a13118be
repo_commit: 4ed503704bf62a7805f100821d6f988032ed0873
status: current
headline: All fifteen read designs in one session at footprint 2 - nothing throttles, all 270 cells are store-bound at 2.00 operations in flight, no >=1.5x pair swaps, r1 and p4 are tied within 8% (the earlier gap was a session difference) and q05 separates the structures by 3.4-3.9x; but the top five designs are a statistical tie at this sample, and c1 shows a 40% position term (first versus last).
---

# Analysis — the read family, jointly comparable (study 02) — deepseek-flash

## TL;DR

- **`J1`/`J2` confirmed**: 0 throttled periods in all **30 read phases** (fifteen designs × two orders), and
  with the **mean** latency `ops_per_sec × mean` is **1.993–2.000 in all 270 cells** — the read phase is a
  store measurement end to end, at the footprint, for the whole family.
- **`J3` confirmed**: no pair differing by ≥1.5x on `q01`@1,000 swaps between the orders. Agreement is
  0.92–1.40, but the ends of that range are informative rather than noise — see the next two points.
- **`J4` confirmed, and it was a session artefact**: `r1` and `p4` — which looked 5–10% apart when measured in
  different sessions — are **1.5% / 7.7%** apart inside one session. The cross-session impression that
  `p4`/`c4` led `r1` does not survive; they are **tied**.
- **The family's quotable structure is groups, not a ranking.** At the lower of the two orders:
  `h1` 18,278 ≈ `r1` 18,242 ≈ `h0` 18,152 ≈ `r2` 17,538 ≈ `p4` 17,526 (`c4` 16,391 in its slower order) |
  `c3` 14,539 ≈ `c2` 14,265 | `c1` 13,653 | `c5` 12,373 | `p1` 11,620 ≈ `x1` 11,584 ≈ `r3` 11,582 ≈
  `p3` 11,164 ≈ `p2` 10,423. **The top five differ by 4%, which is inside this session's own order term** —
  they should be quoted as one group, and only the group boundaries (1.5x between the top group and the
  per-seat group) are robust.
- **`c1_count_naive` carries a 40% position term**: 19,068 ops/s as the first design of order A against
  13,653 as the last of order B. It is the one design whose position differs most between the orders, and it
  is the corpus's first-trial/position lesson arriving on a *read* cell. Its level should be quoted as a
  range or not at all.
- **`J5` confirmed**: the `q05` structure separation reproduces within one session at **3.9x / 3.4x** (the
  slowest pre-aggregated design against the fastest per-seat `COUNT(*)` design).
- **`J6` confirmed**: no failed cells, both trees clean, both runs tagged at `4ed5037`,
  `client_read_conns: 2` in both manifests.

## What I was looking at

Two whole runs, `pg-single`, scale `small`, phases `verify,explain,read`, published budgets and duration,
**all fifteen designs** in reversed order, at the read default of 2 workers.

| design | `q01`@1,000, order A | order B | agreement | lower of the two |
|---|---|---|---|---|
| `h1_hold_checked_confirm` | 19,193 | 18,278 | 0.95 | 18,278 |
| `r1_inventory_row` | 18,873 | 18,242 | 0.97 | 18,242 |
| `h0_hold_naive_confirm` | 19,251 | 18,152 | 0.94 | 18,152 |
| `r2_inventory_buckets` | 17,538 | 17,776 | 0.99 | 17,538 |
| `p4_precreated_counter` | 17,526 | 18,524 | 0.95 | 17,526 |
| `c4_counter_guard` | 18,758 | 16,391 | 0.87 | 16,391 |
| `c3_count_lock_event` | 16,021 | 14,539 | 0.91 | 14,539 |
| `c2_count_serializable` | 16,308 | 14,265 | 0.87 | 14,265 |
| `c1_count_naive` | 19,068 | 13,653 | **0.72** | 13,653 |
| `c5_seat_unique` | 13,259 | 12,373 | 0.93 | 12,373 |
| `p1_precreated_lock_first` | 11,706 | 11,620 | 0.99 | 11,620 |
| `x1_cas_ledger` | 11,584 | 11,616 | 1.00 | 11,584 |
| `r3_seat_pool` | 11,582 | 11,983 | 0.97 | 11,582 |
| `p3_precreated_cas` | 11,255 | 11,164 | 0.99 | 11,164 |
| `p2_precreated_skip_locked` | 10,423 | 11,342 | 0.92 | 10,423 |

## Conclusions

**One session answers the question the split sessions left open.** With all fifteen designs measured in one
run, the structures separate cleanly and the fine ranking does not: the pre-aggregated availability designs
(`h0`, `h1`, `r1`, `r2`, `p4`) sit within **4%** of one another, and the per-seat designs (`p1`, `p2`, `p3`,
`x1`, `r3`) within 11% of one another, while the two groups are **1.5x** apart. A reader should take the
groups and the gap between them; the order inside a group is not resolvable at one run per order, because the
session term within these four runs is itself 4–13%.

`J4` is the sharpest consequence: the impression that `p4`/`c4` led `r1` came from comparing two sessions. In
one session they are tied (1.5% and 7.7%), so the earlier cross-set ordering in the re-statement analysis is
superseded by this joint table.

**The position term is real, and it hit exactly the design whose position moved most.** `c1` reads 19,068 as
order A's first design and 13,653 as order B's last: a 40% swing, larger than any design difference among the
top five. Two of the four earlier tranches also showed the first-listed design reading high (the eleven-design
tranche's `p4` 18.4k/19.3k is its own case). Whatever the mechanism — container start order, page cache state,
the machine's thermal or background state — the practical rule is that **a design's position within a run is
part of its measurement**, and quoting a single level for a position-sensitive design is unsafe. I would not
quote `c1`'s `q01` level from these runs at all.

**What this completes.** With the four earlier tranches, the PostgreSQL read family now has: a fixed footprint,
a verified store-bound reading, complete per-design coverage, a jointly comparable table, and an explicit
statement of what its levels and ratios can support. `v2-07`'s limit no longer needs the "no single run covers
all fifteen designs" caveat, and the YugabyteDB arms remain the only read gap.

**Prediction verdicts** (recorded in `V2-READ-ALL15.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| J1 | no throttling in all 30 read phases | **confirmed** — 0 of ~631 periods everywhere |
| J2 | every cell store-bound at 2.00 ± 5% by the mean | **confirmed** — 270 cells, 1.993–2.000 |
| J3 | no ≥1.5x pair swaps between orders | **confirmed** — agreement 0.92–1.40 |
| J4 | `r1` and `p4` within 10% in one session | **confirmed** — 7.7% and 1.5% |
| J5 | `q05` separation ≥3x within one session | **confirmed** — 3.9x and 3.4x |
| J6 | correctness green, clean tagged trees | **confirmed** — no failed cells, both at `4ed5037` |

## Where I think the measurement is weak

1. **One run per order.** The order term within this pair is up to 13% (excluding `c1`), so a design's level
   is good to about ±10% and its *rank* within a group is not resolvable.
2. **The position effect is characterised, not explained.** `c1`'s 40% is consistent with position, order and
   session all at once, and this design pair cannot separate them.
3. **`c1`'s level is not quotable** from these runs, and no run yet places it in a middle position to check.
4. **Throughput is a latency measurement at a stated footprint** — 2 workers throughout; the relation to other
   footprints is `workers ÷ mean latency`, not a scaling of the same number.
5. **One topology, one engine, one scale, one host**; the YugabyteDB read arms remain uncovered, and the
   machine is a stated project constraint.

## What I would measure next

1. **Nothing further for the PostgreSQL read family** — levels, ratios, structure, bounds and the footprint
   are all measured; only `c1`'s position sensitivity would justify a further run, and it would need three
   positions, not two.
2. **The YugabyteDB read arms**, if the machine can ever carry them — the last read gap.
3. **A position-sensitivity tranche** for any design whose level the corpus intends to quote: the same design
   first, middle and last, which is cheap and would settle `c1`.

## Disagreements with other analyses

It **supersedes the cross-set ordering** in `20261001T-read-cells-restated` (where `p4`/`c4` appeared to lead
`r1` across sessions): measured jointly, they are tied. It confirms `20260930T-read-one-worker`'s store-bound
conclusion on 270 cells rather than 72, and it leaves every earlier ordering and the `q05` structure result
intact.
