---
analysis_id: 20260930T-read-client-sizing--deepseek-flash--2026-09-30
supersedes: 20260929T-read-family-rerun--deepseek-flash--2026-09-29
run_id: 20260930T000544Z (client 1 core), 20260930T000759Z (client 2), 20260930T001013Z (client 3); and the re-run's 20260929T165533Z, 20260929T170353Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; the client-budget calibration that corrects its own previous attribution."
analyzed_at: 2026-09-30
inputs_digest: 8edd7768cc45c254 (client 1), 058075bb2f18736d (client 2), 959431d049995a6e (client 3), 853e38158dac6460, dd6f3d3bad2f27f4 (the re-run's two orders)
inputs: 20260930T000544Z@8edd7768cc45c254, 20260930T000759Z@058075bb2f18736d, 20260930T001013Z@959431d049995a6e, 20260929T165533Z@853e38158dac6460, 20260929T170353Z@dd6f3d3bad2f27f4
repo_commit: cb3099f5c37386e9dfc84f7ea3a8a745c5e49650
status: current
headline: Giving the read client a third core removes its throttling entirely and makes the reads slower - so the client's quota was a symptom, not the ceiling, and r1's flat q01 is not a client floor; the read phase on this machine is bound by the database's share, and the fix is a smaller client footprint or a bigger machine, not client sizing.
---

# Analysis — the read client's quota was a symptom, not the ceiling (study 02) — deepseek-flash

*Supersedes `20260929T-read-family-rerun--deepseek-flash--2026-09-29`, whose measurements stand and whose
attribution does not. That analysis read `r1_inventory_row`'s flat `q01` as "a client floor"; this one varies
the client's budget and shows the flatness survives when the client is not throttled at all. The superseded
file moves to `reports/outdated/`.*

## TL;DR

- **The client's quota is not the bind.** At `CLIENT_CPUS=3` the client's throttling is **0 of 633 periods**
  (it was 422 of 633 at the study's two-core budget), and the reads get **slower**, not faster: `r1`'s
  `q01_event_availability` falls from 52,884 to 43,726 ops/s at 10 seats and from 45,042 to 41,628 at 100,000.
- **Falsified prediction, again in the informative direction.** I predicted a sized client would stop being
  flat. It is still flat (43.7k → 41.6k across a 10,000x tier range, ratio 1.05), so **a flat cell is not
  evidence of a client ceiling** — my previous analysis over-read that signature. What the calibration shows
  is that the level is set by the *machine's* 4 CPUs shared between a 2-core database and the client: the
  database's throttled periods rise **238 → 847 → 1052** as the client's share grows.
- **The study's budget is inside the bind, not outside it.** `DB_CPUS=2 + CLIENT_CPUS=2 = 4` = the podman
  machine's whole CPU count, and the declared memories (`3g + 2g`) are 2.5x the machine's 2 GiB. The read
  phase therefore cannot be repaired by re-allocating cores on this machine; it needs a **smaller client
  read footprint** (fewer, longer-lived connections) or a **bigger machine**.
- **The measurement to quote from the read family is its ordering, not its levels** — `q05`, the `q03`/`r03`
  gap and the `r05` answerability result all reproduce across sessions and budgets, which is the part with a
  mechanism.
- **A capture caveat found on the way**: the database container's `throttled_usec` exceeds the wall time its
  own `nr_periods` implies (124.5 s of throttling in 106.3 s of periods), so only its *period counts* and
  direction are read here. The report's hazard table implies those seconds are directly interpretable.

## What I was looking at

Three runs, identical but for the client's quota — 1, 2 and 3 cores against an unchanged two-core database on
a four-CPU machine — one design (`r1_inventory_row`, the audit's most client-bound), `pg-single`, scale
`small`, phases `verify,read`. The client's budget became settable for this (`ADS_CLIENT_CPUS`, added in the
same change as the protocol; the database's budget already had an override and the client's did not, which is
why this question had not been asked in thirteen days of read runs). Both containers' cgroup counters are
recorded in every run, and the client's are per phase.

| `CLIENT_CPUS` | client read phase: CPU / periods / throttled | database, same cell: throttled periods |
|---|---|---|
| 1 | 63.3 s / 633 / **633 (100%)** | 238 |
| 2 (published) | 111.7 s / 633 / **422 (66.7%)** | 847 |
| 3 | 117.3 s / 633 / **0 (0%)** | 1052 |

## Conclusions

**Removing the client's throttle makes the reads slower.** This is the whole result in one line. At three
cores the client is never throttled and uses 117 s of CPU over 63 s of wall — more total work than at two
cores — yet every `q01` tier lands below its two-core value (52,884 → 43,726 at 10 seats; 50,031 → 40,264 at
100; 45,047 → 39,599 at 1,000; 45,020 → 39,650 at 10,000; 45,042 → 41,628 at 100,000). The third core came
out of the database's share of the same four CPUs, and the database's throttled periods rose by a quarter.
**The client's throttling was the visible end of a machine that is 100% allocated before any experiment
starts.**

**What the flat cell actually is.** At two cores `r1`'s `q01` runs 45.0–52.9k across the five tiers; at three
cores, 39.6–43.7k. It is flat in both, so the flatness is not caused by the client's quota, and my previous
analysis's inference — "a number that does not move when the query gets 10,000x harder is the client's
ceiling" — is wrong as stated. Two readings survive the evidence and this tranche cannot separate them:
`r1_inventory_row`'s availability read may be genuinely **O(1)** (one inventory row per event, which is the
design's stated point), or it may be sitting on a **machine-level ceiling** (~40–53k) that no tier can push
through. Distinguishing them needs a plan check on that query at two tiers, not another budget. What is
settled is that it is **not a client floor**.

**Consequences for the read family.** (1) Its levels are not quotable on this machine, and no re-allocation
of cores fixes that; the actionable change is the client's read **footprint** — fewer, longer-lived
connections — or a larger machine, both owner decisions because they alter a study constant. (2) Its
**orderings** remain quotable and are the part to cite: `q05` fast where there are no per-seat rows, the
`q03`/`r03` gap of two orders of magnitude, `r05` answered only by the ledger at ~32k. (3) My previous
analysis's headline must not be cited; this file replaces it.

**Prediction verdicts** (recorded in `READ-CLIENT-SIZING.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| S1 | at two cores the re-run's ~420/634 throttle share reproduces | **confirmed** — 422/633 (66.7% against 66.2%) |
| S2 | at three cores the client's throttling falls below 10% of periods | **confirmed** — 0/633 |
| S3 | a sized client stops being flat (100k at least 2x below 10 seats) | **falsified** — 43,726 → 41,628, ratio 1.05: still flat, and *lower* overall |
| S4 | the bind moves: the database's throttling rises when the client's share grows | **confirmed** — 238 → 847 → 1052 throttled periods |
| S5 | correctness gates stay green | **confirmed** — no failed cells in any budget |

## Where I think the measurement is weak

1. **One design and one run per budget.** `r1` is the extreme case by construction; the other fourteen
   designs' share sensitivity is unmeasured, and a single run per budget cannot separate the budget effect
   from session drift — the two-core `q01` level moved from ~38k in the re-run to ~45–53k here.
2. **The machine is not a controlled instrument.** With the two budgets summing to the machine's whole CPU
   count, both containers' throttling responds to each other; a three-way allocation (2/2, 2/3, 3/2) would be
   needed to attribute the level to the database's share rather than to the client's, and the machine has no
   spare core to do it with.
3. **The database's throttled *time* is internally inconsistent on this host** (124.5 s in 106.3 s of
   periods), so the database column above is period counts only.
4. **The memory question is untouched and worse than the CPU one**: declared `3g + 2g` against a 2 GiB
   machine, still the owner's open decision.
5. **One tier set, one scale, one topology, one engine, one host** — as for the rest of study 02.

## What I would measure next

1. **A plan check at two tiers** to settle whether `r1`'s `q01` is O(1) by design, which decides whether its
   flatness is a finding or a ceiling.
2. **The client's read footprint** — halve its read connections and re-measure `r1` at the published
   two-core budget: if the throughput rises, the client's per-connection work was the cost and the fix is in
   the harness rather than the machine.
3. **A three-way budget matrix** (database 2 / client 2, 2/3, 3/2) on a machine with a spare core, to
   attribute the level properly — which is the owner's machine decision.

## Disagreements with other analyses

It **supersedes my own** `20260929T-read-family-rerun` on attribution, keeping its measurements (the
hazard's reproducibility, the `q05` and `r03` findings, and the non-reproduction of X1's 2x read penalty) and
correcting its "floor is the client" headline. It also qualifies the hazard-table reading implicit in the
`study02-client-quota-audit`: the client's throttling there is real and worth recording, but on this machine
it is downstream of the total allocation rather than the cause of the read levels.
