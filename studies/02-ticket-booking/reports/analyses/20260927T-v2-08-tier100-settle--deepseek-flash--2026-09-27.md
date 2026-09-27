---
analysis_id: 20260927T-v2-08-tier100-settle--deepseek-flash--2026-09-27
run_id: 20260927T220823Z (p3 then x1), 20260927T220927Z (x1 then p3), 20260927T221035Z (x1 then p3), 20260927T221151Z (p3 then x1), 20260927T221309Z (p3 then x1), 20260927T221420Z (x1 then p3)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW/worker-capability session, Deep Code CLI; executing the tier-100 settle tranche of study 02's v2-08 replication."
analyzed_at: 2026-09-27
inputs_digest: aa6a391a407e00fd (pass A1, pg-single), 0400a8c02823e97b (B1), cbdf513f3bec02e0 (B2), 170a5377488f46de (A2), 1fe5231a4ef25173 (A3), 60936cc50483a8f3 (B3)
inputs: 20260927T220823Z@aa6a391a407e00fd, 20260927T220927Z@0400a8c02823e97b, 20260927T221035Z@cbdf513f3bec02e0, 20260927T221151Z@170a5377488f46de, 20260927T221309Z@1fe5231a4ef25173, 20260927T221420Z@60936cc50483a8f3
repo_commit: e256e77c958ff468af10bae0b624573919c96e2a
status: current
headline: The 100-seat reversal was a run-median artefact of a cell whose trials disagree by half their own median; with 30 trial pairs the ledger is slower in 29, every run median is inside the claim's 1.1-2.3x band, and the thing that still fails at this tier is "tight spreads".
---

# Analysis — study 02 `v2-08`, tier-100 settle — deepseek-flash

## TL;DR

- **The 0.99x reversal does not survive more evidence.** Over six new runs (5 race trials each, three in
  each design order), the tier-100 `p3/x1` run median was **1.279, 1.333, 1.257, 1.712, 1.328, 1.476** —
  every one above the claim's 1.1x floor, so the band is not widened by this evidence.
- **The ledger was slower in 29 of 30 individual trial pairs.** The one exception (0.940x, run
  `20260927T221035Z`, trial 3) is the same kind of reading as the earlier 0.97x: a single pair, not a
  property of the tier. Pooled per-trial median ratio **1.391x**.
- **"Tight spreads" is still false at this tier.** Every run's 100-seat cell spread **47–77%** for `p3`
  and 26–54% for `x1`; the run-to-run spread of the ratio is **34%** over these six runs and **54%** over
  all eight. This, not the band, is what HIGH's re-scoping should fix.
- **The harness still sees the cost at the neighbouring tier.** All six tier-1,000 ratios are **> 1.30**
  (1.304–1.712), so the negative reading at tier 100 was a cell property, not a broken run.
- **Correctness green everywhere**: 46/46 and 48/48 verifications, 0 over-booked, 0 under-booked,
  0 timed-out events in all six runs.
- **I did not re-scope the claim.** The `1.1–2.3x` band and "tight spreads" are signed wording; this
  tranche supplies the evidence and leaves the decision with HIGH, as the 2026-09-27 replication
  analysis (`20260927T-v2-08-replication--deepseek-flash--2026-09-27`) already recorded.

## What I was looking at

Six whole runs of study 02 on `pg-single`, scale `small`, phases `verify,race`, designs
`p3_precreated_cas` and `x1_cas_ledger`, 128 buyers per event, **5 race trials per cell**, tiers **100 and
1,000**. The run orders were three `p3,x1` (`20260927T220823Z`, `20260927T221151Z`, `20260927T221309Z`)
and three `x1,p3` (`20260927T220927Z`, `20260927T221035Z`, `20260927T221420Z`). All six trees were clean,
all six are tagged `run/02-ticket-booking/<run-id>` at `e256e77`, and no cell failed. The tier-1,000 cell
is a **positive control**: it is the neighbouring tier where both earlier passes saw a cost, so it tells
us whether a null result at tier 100 is a real absence or a broken session.

I ignored the other three tiers, the YugabyteDB topologies, the refund ratio and the organiser-edit
effect, all named as a coverage gap in the pre-registration
(`studies/02-ticket-booking/V2-08-TIER100-SETTLE.md`).

## Conclusions

**The reversal was an artefact of taking a median from a cell that cannot support one.** The previous
pass read 0.99x because `p3`'s three trials ranged 2,343–4,065 ops/s and the reversed pass's middle trial
happened to favour the ledger; its own analysis said the tier was intrinsically noisy and to settle it
with more pairs. Thirty pairs say the opposite of a reversal: the ledger is slower in 29, and the pooled
per-trial median ratio is 1.391x. The single 0.940x pair is not a different regime — it is the tail of a
distribution whose spread is 50–77% of its median. A claim phrased as a band over a central tendency
survives this; a reader who took any one trial as the ratio would not.

**The band is not the problem; the word "tight" is.** All six run medians (1.257–1.712) and the pooled
median (1.391) sit inside the claim's `1.1–2.3x`. So the third pass **confirms the band at 128 buyers and
tier 100** rather than forcing it open — which narrows the HIGH decision considerably: what the evidence
does not support is "tight spreads". At tier 100 the `p3` cell's spread was 47%, 48%, 51%, 62%, 64% and
77% across the six runs, and the ratio itself moved 34% across six runs (54% across all eight, 0.992–1.712).
A ratio whose run-to-run spread is a third of its own value is a direction, not a point; even the direction
has a sub-1 tail in 1 of 30 pairs.

**The tier is not the whole story: there is a first-trial effect.** In **all six** runs `p3`'s trial 1 was
below its own run median (2,236 vs 3,843; 2,377 vs 3,613; 2,516 vs 2,849; 2,451 vs 3,993; 2,260 vs 3,601;
1,839 vs 2,984). That is the same within-cell position effect study 02's `v2-06` replication and study 04's
T1 work found, and it is large enough here to move a three-trial median by a tier's worth of ratio: in the
first pass, dropping trial 1 would have raised the tier-100 ratio. Any future tier-100 number should either
run enough trials that one warm-up trial cannot set the median, or state how it is handled — this tranche
does the former (5 trials), and the effect is a limit, not something the claim currently records.

**The control did its job.** Tier 1,000 read 1.304–1.712 in all six runs, so the harness and the session
could detect the cost throughout; the tier-100 result is therefore interpretable as a cell property rather
than a failed run. Correctness was green in every run (46/46 `p3`, 48/48 `x1`, zero bad counters), so no
cell is void.

**Prediction verdicts** (recorded in `V2-08-TIER100-SETTLE.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| P1 | `p3` faster in ≥20 of 30 pairs | **confirmed** — 29 of 30 |
| P2 | pooled median ratio in [1.10, 1.60] | **confirmed** — 1.391 |
| P3 | tier-1,000 ratio > 1.10 in every run | **confirmed** — min 1.304 |
| P4 | correctness green in every run | **confirmed** — 0 bad cells |
| P5 | tier-100 spread > 40% in ≥4 of 6 runs | **confirmed** — 6 of 6 |

## Where I think the measurement is weak

1. **Thirty pairs, one host, one engine, the 2 GiB machine ceiling.** The podman machine's 2 GiB against
   declared 3 GB + 2 GB budgets is still unenforced (CONTEXT, 2026-09-27), so the absolute ops/s carry the
   same envelope gap as every other study 02 number; only the within-run ratio is being read here.
2. **Pairing trials by index is a convention, not a control.** The five trials of `p3` and `x1` are
   independent fresh-load measurements; pairing trial *i* with trial *i* assumes a shared session term at
   that position. It is the same assumption the two earlier passes' medians made implicitly, and it is why
   the first-trial effect above matters — but a differently paired sample could move the 29/30.
3. **The tier-100 cell is measured near a floor.** The whole tier races 20 events of 100 seats against 128
   buyers; at that size the sell-out is short and the per-second rate is dominated by the client's fixed
   work, which is a plausible mechanism for both the wide spread and the small ratio. This tranche does not
   separate client-limited from server-limited behaviour.
4. **No session drift was visible across these six runs** (`p3` medians 3,843 / 3,613 / 2,849 / 3,993 /
   3,601 / 2,984), but six runs in ten minutes cannot rule one out the way study 04's reference cell did.

## What I would measure next

1. **The re-scoping decision (HIGH)** — on this evidence the band stands and "tight spreads" does not; the
   narrowest honest fix is to drop or qualify the spread wording and state the tier-100 uncertainty.
2. **A within-run position control at tier 100** — the same design first and last in one cell — to put a
   number on the first-trial effect rather than inferring it from medians.
3. **Tiers 10 and 10,000**, the two the `v2-08` replication left out, so the "cost is a function of tier"
   picture is anchored at both ends.
4. **The YugabyteDB half**, untouched by two PostgreSQL tranches and asserted by the same claim.

## Disagreements with other analyses

None with the two analyses this extends. `20260927T-v2-08-replication--deepseek-flash--2026-09-27`
reported the reversal and flagged it as possibly noise; this tranche resolves that specific doubt in the
direction the earlier analysis's own weakness section predicted, and leaves its re-scoping flag standing.
No other analysis covers these digests.
