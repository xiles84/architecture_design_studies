---
analysis_id: 20260929T-v2-09-8buyer-recurrence--deepseek-flash--2026-09-29
run_id: 20260929T010633Z (A), 20260929T010807Z (B), 20260929T010947Z (B), 20260929T011128Z (A), 20260929T011301Z (A), 20260929T011444Z (B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; the 8-buyer collapse recurrence test."
analyzed_at: 2026-09-29
inputs_digest: 76f431bf0eb136fc (A1), 2e5d35a847eda9ee (B1), a86da7ea90c90c9a (B2), 2eeac99e2287cfa9 (A2), 4a03a51500ee3a3a (A3), 0dbf1c3abb037c71 (B3)
inputs: 20260929T010633Z@76f431bf0eb136fc, 20260929T010807Z@2e5d35a847eda9ee, 20260929T010947Z@a86da7ea90c90c9a, 20260929T011128Z@2eeac99e2287cfa9, 20260929T011301Z@4a03a51500ee3a3a, 20260929T011444Z@0dbf1c3abb037c71
repo_commit: a253658bc64080ef7c9cef8fc6af71eb03a6aa7a
status: current
headline: The 8-buyer low-throughput mode recurs (4 events in 60 trials, ~4x below the session's ~2,500 ops/s) but it is not ledger-specific and not order-specific - it hit both designs; the persistence of the original cell was an artefact of where the event fell in the cell.
---

# Analysis — study 02 `v2-09`, does the 8-buyer collapse recur? — deepseek-flash

## TL;DR

- **It recurs.** Four low-throughput events in 60 trials (12 cells, six runs): readings of **619–655 ops/s**
  against a session normal of ~2,500 — a ~4x drop. So the 2026-09-28 event was not a one-off.
- **It is not the ledger's.** Three of the four events are in the **CAS** arm (`p3`), and only one in the
  ledger arm; both arms and both orders are affected. My C2 and C4 predictions are **falsified**.
- **Persistence is the exception.** Three events were a single trial that recovered on the next one; only
  one spanned consecutive trials, and that was the original cell's three-trial run.
- **The original "collapse" was mostly an artefact of position.** Its ledger arm's last three trials were
  all low, so the cell median was the low mode and the earlier two trials looked like the anomaly. Had the
  event fallen one trial earlier and recovered, it would have read as a spike.
- **The session's normal is half the sweep's** (~2,500 vs ~5,000 ops/s at the same cell), another
  day/session term, which is why only within-cell series are read here.
- **Correctness green**: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors in all twelve cells.
- **Consequence**: no 8-buyer magnitude is quotable, and the cause is still open — the artifact carries the
  database's cgroup counters but not the client's, so a client stall cannot be separated from a
  database-state change.

## What I was looking at

Six whole runs on `pg-single`, scale `small`, phases `verify,race`, designs `p3_precreated_cas` and
`x1_cas_ledger`, **8 buyers per event**, **tier 1,000 only**, 5 trials, orders alternating `A B B A A B`
(six cells per arm). All trees clean, all tagged at `a253658`, no failed cells. The unit is the **per-trial
series** of each design-cell, never its median.

## Conclusions

**The low mode is real, occasional and design-blind.** Across 60 trials there are four readings of 619–655
ops/s where the session's normal is ~2,500: `p3` trial 4 of `010807Z` (655, recovery to 2,546), `x1` trial 5
of `010947Z` (619, cell end), `p3` trial 4 of `011301Z` (623, recovery to 2,642), and the original cell's
three-trial run. That is 4 events in 6 runs, roughly 7% of trials, and **three of the four are in the CAS
arm**. The pre-registered prediction that this would be a ledger-arm effect is falsified, as is the
prediction that it would favour the CAS-first order; the one event that fell in the ledger arm was in the
other order. What looked like "the ledger arm collapses" in one cell was a 1-in-6 roll of a slow trial
landing at the end of a cell, where it sets the median.

**Persistence needs a definition, and the honest one is positional.** A drop that recovers is a spike; a
drop that continues is a regime. This tranche saw three spikes and one run of consecutive low trials — and
the only reason the original cell read as a regime rather than a spike is that its slow trials were its last
three, so the median moved with them. **A cell whose low trials happen to be terminal looks like a state
change; a cell whose low trial is in the middle looks like an outlier.** With five trials and a 7% event
rate, the position of the event decides its name, which means neither name is the mechanism.

**The measurement is therefore fragile at this buyer count in a specific way.** At 8 buyers each trial sells
only a handful of seats, so a trial's duration is dominated by a few operations; one slow operation is a 4x
drop, and a run of them moves a median of five. That is why the sweep's 8-buyer level gave the 3.08 pooled
median that this tranche has now explained: it was not the ledger, and it was not a regime. No 8-buyer
magnitude is quotable from either tranche, and the `v2-09` limit has been corrected to say so.

**What is still not known.** The artifact records the database container's `cpu.stat` (the original cell's
throttling was ~4 ms across the cell, i.e. negligible) but nothing from the client container, and no plan or
clock per trial. So "client stall" and "database state change" remain indistinguishable here, and the fix is
instrumentation — the client's own cgroup counters beside the database's — rather than more runs.

**Prediction verdicts** (recorded in `V2-09-8BUYER-RECURRENCE.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| C1 | at least one of twelve cells shows a ≥2x persistent drop | **confirmed** — but see C3 |
| C2 | every collapsing cell is a ledger cell | **falsified** — three of the four low events are CAS cells |
| C3 | drops persist rather than recover | **falsified as a general rule** — three of four events recovered; only the terminal one appears persistent |
| C4 | collapses favour the CAS-first order | **falsified** — the one terminal event was in the ledger-first order |
| C5 | correctness green | **confirmed** — 0 bad counters in all twelve cells |

## Where I think the measurement is weak

1. **Six runs, 60 trials.** Four events is enough to show the mode is real and design-blind; it is not
   enough to estimate its rate, and a 7% estimate has a wide interval.
2. **One tier and one buyer count.** Nothing here says whether the mode exists at 32 or 128 buyers — the
   sweep's tier-1000 cells there spread only 12–42%, so it may be specific to very low concurrency, where
   each trial is a handful of operations.
3. **No client instrumentation.** The mechanism cannot be assigned from this artifact.
4. **The session's absolute level is half the previous day's**, so the 619–655 readings are only comparable
   with this session's ~2,500; the ratio (4x) travels, the absolutes do not.
5. **The 2 GiB podman-machine ceiling** applies, as to the rest of study 02.

## What I would measure next

1. **Instrument the client container** (its cgroup `cpu.stat` beside the database's) for a tranche of these
   runs; without it the mechanism is unassignable and more cells only sharpen the rate.
2. **The same test at 32 buyers**, where each trial is longer, to see whether the mode is a property of very
   short trials or of the pair.
3. **A discard rule, pre-registered**: since a low trial can be identified by its own duration and the
   harness already records `wall_s`, a future protocol could state a minimum wall time below which a trial
   is void — decided before the run, not after seeing which cells it spoils.

## Disagreements with other analyses

None. It corrects the emphasis of the 2026-09-28 sweep analysis, which reported the collapsed cell as "the
ledger arm's" mid-cell state change; on the recurrence evidence the mode is not the ledger's, and the
sweep's own caution (that it kept the cell rather than discarding it) is what made this test possible.
