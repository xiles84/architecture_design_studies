---
analysis_id: 20260927T-v2-08-tier-axis--deepseek-flash--2026-09-27
run_id: 20260927T224424Z (p3 then x1), 20260927T225107Z (x1 then p3), 20260927T225817Z (x1 then p3), 20260927T230535Z (p3 then x1)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW/worker-capability session, Deep Code CLI; executing the tier-axis tranche of study 02's v2-08 replication."
analyzed_at: 2026-09-27
inputs_digest: 0ada8aefa00f60f9 (A1), 3af58e45f03b6428 (B1), 362e803f13bb5acb (B2), f4a028309578488c (A2)
inputs: 20260927T224424Z@0ada8aefa00f60f9, 20260927T225107Z@3af58e45f03b6428, 20260927T225817Z@362e803f13bb5acb, 20260927T230535Z@f4a028309578488c
repo_commit: 3eed9380ad28e6ccce2c5208b94e71440ad3fe9f
status: current
headline: The ledger's race cost is monotone in tier - 1.08x at 10 seats to 2.51x at 100,000 - so the claim's 1.1-2.3x band describes its middle three tiers and both endpoint tiers are boundaries, with the 10-seat direction unreliable (60% of pairs).
---

# Analysis — study 02 `v2-08`, the ledger's cost across the tier axis — deepseek-flash

## TL;DR

- **The cost is monotone in tier.** Pooled run-median `p3/x1` across four runs, both design orders, five
  trials per cell: **1.075, 1.237, 1.409, 1.845, 2.505** at 10 / 100 / 1,000 / 10,000 / 100,000 seats. No
  adjacent inversion in any single run either.
- **Both endpoints sit outside the claim's `1.1–2.3x` band**, reproducing the primary run's 1.06 and 2.35.
  The band describes tiers 100–10,000; at 10 seats the ledger is nearly free, at 100,000 it costs ~2.5x.
- **At 10 seats the cost is indistinguishable from zero.** The run-median ratio is 1.075 (below the floor),
  but the per-trial median is **1.131** (at the floor) and the ledger is slower in only **12 of 20 pairs
  (60%)** — the only tier where the direction is unreliable. This is a boundary, not a value.
- **At 100,000 seats the ledger is slower in 20 of 20 pairs** (pooled per-trial median 2.456, range
  1.883–3.155), so the top endpoint is a solid boundary — with the caveat that the tier is **one event per
  trial**.
- **Correctness green in all four runs**: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors.
- **Still not a re-scoping.** The band and "tight spreads" are signed; this tranche supplies the boundary
  evidence the HIGH decision needs, alongside the tier-100 settle's
  (`20260927T-v2-08-tier100-settle--deepseek-flash--2026-09-27`).

## What I was looking at

Four whole runs on `pg-single`, scale `small`, phases `verify,race`, designs `p3_precreated_cas` and
`x1_cas_ledger`, **128 buyers**, **5 race trials**, all **five tiers** in every run, two runs in each design
order (`20260927T224424Z` A, `20260927T225107Z` B, `20260927T225817Z` B, `20260927T230535Z` A). The primary
anchor `20260921T205212Z` measured the same five tiers once, in one order, three trials; the 2026-09-27
replication and settle covered tiers 100/1,000/10,000 and 100/1,000 respectively, so **the 10-seat and
100,000-seat tiers had never been repeated**. All four trees were clean, all four are tagged
`run/02-ticket-booking/<run-id>` at `3eed938`, and no cell failed.

## Conclusions

**The tier is the axis the band was always hiding.** The claim reads "the ledger costs 1.1–2.3x of the
sell-out race at 128 buyers", which presents a range as though it were the cost. It is not: the cost is a
monotone function of how big the event is. Pooled run medians rise 1.075 → 1.237 → 1.409 → 1.845 → 2.505,
and every individual run shows the same ordering (e.g. A1: 1.023, 1.398, 1.634, 1.902, 2.558). The band's
floor and ceiling are simply the 100-seat and 10,000-seat tiers in the primary run, so **a reader who quotes
"1.1–2.3x" is quoting two interior points and calling them a range**. This is the same defect shape the
corpus already recorded for `v2-13` ("A ratio quoted without its axis is not a finding"), now measured
directly: the axis is tier, and the reading at either end is outside the band.

**At 10 seats the ledger's cost is not measurable above zero — and the direction is not reliable.** This is
the one place where the three statistics disagree and the disagreement is the finding. The run medians are
1.023, 1.024, 1.180, 1.126 (pooled 1.075) — below the claim's 1.1x floor, reproducing the primary's 1.06.
But the per-trial median is 1.131, and only 12 of 20 pairs put the ledger behind (the other four tiers:
18, 20, 20, 20 of 20). The cell is also the noisiest in the tranche (spreads 16–71%, pair range
0.765–2.065). So P2 is **confirmed on the run-median statistic and falsified on the per-trial one**, and
the honest statement is that 10 seats sits *at* the floor with a direction that flips in 40% of trials —
not that the ledger is measurably cheaper there. This is exactly the kind of ambiguity the analysis must
report rather than resolve by choosing the convenient statistic; my pre-registration named no statistic for
P2 and that was a defect in the protocol.

**At 100,000 seats the top boundary is real but rests on one event.** Twenty of twenty pairs put the ledger
slower, pooled per-trial median 2.456, run medians 2.249–2.637 — all above the band's 2.3x top. The tier is
a single event of 100,000 seats, so "20 pairs" is 20 readings of one event, not 20 events, and the corpus's
`n=1` caveat applies more here than anywhere else in the tranche. The direction is unambiguous; the
magnitude is one event's.

**The middle of the axis still holds.** Tiers 100, 1,000 and 10,000 read 1.237 / 1.409 / 1.845 pooled
(run medians 1.067–1.398 / 1.232–1.634 / 1.685–1.902), consistent with the settle's 1.257–1.712 at tier 100
and the earlier replication's 1.30–2.00. Pooling the settle and this tranche, tier 100 now has **ten runs**
and its ratio spans roughly 1.07–1.71 — a direction, not a point, exactly as the settle recorded.

**Prediction verdicts** (recorded in `V2-08-TIER-AXIS.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| P1 | monotone increase across the five pooled run medians | **confirmed** — 1.075 → 2.505, no inversion, and none within any run |
| P2 | 10-seat ratio below the 1.1x floor | **confirmed on run medians (1.075), falsified on the per-trial median (1.131)** — the tier sits at the floor; protocol named no statistic |
| P3 | 100,000-seat ratio above the 2.3x top | **confirmed** — run-median 2.505, per-trial 2.456 |
| P4 | ledger slower in ≥80% of pairs at every tier | **falsified at 10 seats (60%)**; confirmed at the other four (90%, 100%, 100%, 100%) |
| P5 | correctness green in every run | **confirmed** — 0 bad counters in all four |

## Where I think the measurement is weak

1. **Tier 100,000 is one event per trial.** Twenty pairs are twenty readings of a single event; a
   per-event effect (locks, hot pages in that one event's row) would look exactly like a tier effect.
2. **My P2 falsifier named no statistic**, so it is both confirmed and falsified depending on whether the
   reader takes the run median or the per-trial pair. The pre-registration should have fixed the statistic;
   the result is reported here with both.
3. **Four runs, one host, one engine, the 2 GiB podman-machine ceiling** — every number carries the same
   envelope gap as the rest of study 02, and only within-run ratios are read across sessions.
4. **The 10-seat cell is near-degenerate**: 100 events of 10 seats against 128 buyers, attempts/sale ~13, so
   the per-second rate is dominated by retry work and client scheduling rather than by the ledger.
5. **Absolute throughput still moves between runs** (tier-1,000 `p3`: 6,135 / 5,031 / 4,556 / 5,490), which
   is why only the ratio is interpreted.

## What I would measure next

1. **The re-scoping decision (HIGH)** — the band is a mid-axis range, "tight spreads" fails at the small
   tiers, and the cost is monotone in tier; the boundary evidence now exists for a statement that names the
   axis.
2. **More events at 100,000 seats** (or a tier between 10k and 100k) to separate tier from the single-event
   artefact.
3. **The YugabyteDB half** of the same claim, untouched by three PostgreSQL tranches, plus the refund ratio
   and the organiser-edit effect.
4. **A repeated-design control at 10 seats**, the one tier whose direction is unreliable, to see whether the
   60% is the tier or the harness.

## Disagreements with other analyses

None. The two 2026-09-27 analyses this extends (`20260927T-v2-08-replication--…`, `…-tier100-settle--…`)
reported the middle of the axis and left the endpoints as a stated gap; this tranche fills that gap and does
not contradict either. No other analysis covers these digests.
