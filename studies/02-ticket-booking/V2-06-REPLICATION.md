# Replication: pre-created seats versus counter rows (study 02, `v2-06`)

**Question.** `v2-06-hot-drop-designs` is one of the corpus's headline claims — *"pre-created seats taken by
compare-and-set or SKIP LOCKED sell 4–8x faster than any single counter row"* — and it rests on **one run,
one measurement per configuration**, from 2026-09-13. Every alternative explanation for that number is
still open: the sales race may have favoured the pre-created designs on that particular occasion, the
machine may have been in an unusual state, and no repeat has ever been attempted on a study other than 04.
This tranche re-measures the same comparison, on the same topology, with three fresh-load repeats per tier
inside one run.

**What is held constant.** Designs `p2_precreated_skip_locked`, `p3_precreated_cas`, `c3_count_lock_event`,
`c4_counter_guard` and `r1_inventory_row`; topology `pg-single`; scale `small`; the `race` phase (the
sell-out race, 32 concurrent buyers per event); dataset seed 42. Only the session differs from the
2026-09-13 run.

**Why these five designs.** `p2` and `p3` are the claim's fast side, `c3` and `c4` are the counter rows it
beats, and `r1` is the inventory-row alternative that also beat the counters — included because a
replication that only re-measured the winners would not be evidence.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| T1 | **The headline reproduces.** At 100 and 1,000 seats, `p3_precreated_cas` sells **4–8x** faster than `c4_counter_guard` and `c3_count_lock_event`. | A ratio below 4x or above 8x at either tier against either counter, which would mean the headline's specific range was an artefact of the original session rather than a property of the designs. |
| T2 | **The second fast design reproduces too.** `p2_precreated_skip_locked` also beats both counters, by at least **3x**, at both tiers. | A ratio below 3x at either tier, which would mean `p3`'s result is specific to CAS and the claim's pairing of the two is too generous. |
| T3 | **The magnitudes move by no more than a session term.** Every cell's seats-per-second falls within **±40%** of the 2026-09-13 value for the same design and tier. | Any cell outside ±40%, which would mean the earlier run's numbers are not reproducible even in size, and `v2-06`'s figures should be re-scoped to directions. |
| T4 | **The correctness gate stays green.** Every design passes verification at both tiers and no cell is dropped for a failed check. | Any failure, which voids that cell. |

Recorded before measurement. Evidence form for the claim stays `isolate` (single-variable pairs); this
tranche raises its replication tier from one whole run toward run-level.

**Coverage gap, stated in advance.** This tranche does **not** re-run the deliberately wrong controls
(`c1_count_naive`'s overbooking, `h0_hold_naive_confirm`) that fired in the original run. So the gate here
is a pass-gate, not a fired negative control, and T4 must not be read as re-establishing the controls. That
limitation travels with the result.

## Scope

| Field | Value |
|---|---|
| Runs | one, with `-race-trials 3` (three fresh-load repeats per tier) |
| Tiers | 100 and 1,000 seats per event |
| Existing evidence | `run/02-ticket-booking/20260913T021206Z` (primary anchor of `v2-06`) |

**Boundaries.** One session, one host, one engine, two tiers out of the five the original swept (10, 1k,
10k, 100k are not re-measured), and no repeat of the wrong-design controls. Three repeats inside one run
detect a within-session spread; they cannot detect a session-level bias, which is the error the study 04
tranches kept finding and which this design therefore cannot exclude.
