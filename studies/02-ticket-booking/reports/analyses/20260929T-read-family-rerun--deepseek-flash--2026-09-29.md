---
analysis_id: 20260929T-read-family-rerun--deepseek-flash--2026-09-29
run_id: 20260929T165533Z (p3,x1,c1,r1), 20260929T170353Z (r1,c1,x1,p3)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; the read-family re-run that follows the client-quota audit."
analyzed_at: 2026-09-29
inputs_digest: 853e38158dac6460 (order A), dd6f3d3bad2f27f4 (order B)
inputs: 20260929T165533Z@853e38158dac6460, 20260929T170353Z@dd6f3d3bad2f27f4
repo_commit: a58bf1015cbd905e1455870d8d60931be6e401dd
status: current
headline: The read phase's client throttling is not a period artefact - it reproduces today, design-dependent (r1 loses 51% of its periods, flat at ~38k ops/s across five tiers, which is a client floor and not a query cost), so the read family's numbers need client sizing rather than a re-run; X1's 2x read penalty does not reproduce and the q05 and r03 findings do.
---

# Analysis — the read family after the audit: the floor is the client (study 02) — deepseek-flash

## TL;DR

- **The hazard is not history.** I predicted today's client would not be throttled in the read phase. It is:
  `p3` 2.3 s and 5.9 s, `x1` 2.6 s and 4.6 s, `c1` 5.4 s, and **`r1_inventory_row` 32.3 s of ~63 s of
  periods** — in **both** orders, within 0.1 s of each other. **P1 is falsified**, and falsified in the
  useful direction: this is a sizing defect in the read phase, not a bad afternoon in September.
- **The tell is a flat cell.** `r1`'s `q01_event_availability` reads **37.7–38.7k ops/s at every tier** —
  10, 100, 1,000, 10,000 and 100,000 seats, a 10,000x input range. `p3` over the same range falls 25,967 →
  595. A number that does not move when the query gets 10,000x harder is not measuring the query; it is the
  client's ceiling. `r1`'s reads are a **client floor**, and its `q05` (17.4k) is on the same floor.
- **X1's read penalty was a cell effect, now shown twice.** The old run had X1 at 0.29–0.73 of P3 across the
  buyer reads; today it is **0.77–1.05** in one order and **0.97–1.05** in the other. **P4 confirmed**: no
  ~2x ledger read cost, and the small residual in one order is order-dependent.
- **The two findings the old analysis drew still hold.** `q05` sell-through per band is **C1 12.2–12.7k
  against P3 2.0–2.1k and X1 1.9–2.0k** (P2 confirmed; old: 5.5k against 1.7k), and the gap inside X1
  between `q03_customer_tickets` and `r03` is **143–256x** (P3 confirmed; old 78–191x). `r05` is answered
  only by the ledger at **30.3k and 33.0k** against the old run's 32.2k.
- **Nothing is retracted, and the read family still is not quotable** — for a different reason than the audit
  gave. The audit said "that run's client was throttled"; the truth is "the read phase throttles the client
  for the client-heavy designs, today, reproducibly". The fix is client sizing (or fewer, longer-lived read
  connections), and until then the four designs measured here have a client-bound floor.

## What I was looking at

Two whole runs on `pg-single`, scale `small`, phases `verify,explain,read` — the original's phase set, with
no writes or races — four designs each in reversed order: `p3_precreated_cas`, `x1_cas_ledger`,
`c1_count_naive`, `r1_inventory_row`. Both trees clean, both tagged at `a58bf10`, no failed cells. The
comparison is against `20260920T234953Z` (`pg-single`, all fifteen designs, same phases and scale), the run
whose read phase the audit found throttled by 2.3–58.1 s depending on design.

## Conclusions

**The read phase is client-bound for exactly the designs that do the most client-side work, and it is
reproducible.** Per design, client CPU over the read phase and the throttling it cost:

| Design | order A | order B |
|---|---|---|
| `p3_precreated_cas` | 55.4 s CPU, 70/633 periods (2.3 s) | 56.2 s CPU, 71/632 (5.9 s) |
| `x1_cas_ledger` | 56.0 s, 65/633 (2.6 s) | 56.3 s, 68/633 (4.6 s) |
| `c1_count_naive` | 76.1 s, 76/633 (5.4 s) | 75.7 s, 70/633 (5.4 s) |
| **`r1_inventory_row`** | **111.9 s, 420/634 (32.3 s)** | **111.9 s, 422/633 (32.3 s)** |

Two cores of quota, ~63 s of periods: `r1` uses 1.77 cores on average and loses half its periods. That is
not a database measurement, and the flat `q01` proves it independently of the throttling counters: 37,748 /
38,340 / 38,172 / 38,707 / 37,893 ops/s at the five tiers. The old run showed the same design as its most
throttled (58.1 s), so this is a durable property of `r1_inventory_row`'s read phase under a two-core client,
not a coincidence of one session.

**The audit's conclusion, corrected.** The audit wrote that the 2026-09-20 read figures were "a floor set by
the client's quota" and implied the current harness would not do that. The first half is right and now
better supported; the second half is wrong. `v2-06`, `v2-07` and `v2-gap-01` carry limits that say their
reports' read figures are not quotable — that stays true, and the reason is now the phase's sizing rather
than one session's luck. This is the owner's 2026-09-14 sizing requirement arriving in the read phase: the
generator must be able to supply the intended demand before its numbers can be read as the store's.

**What reproduces is what the designs say, and it is the ordering rather than the level.** The `q05` result
(C1 ≈ 12.4k against P3/X1 ≈ 2.0k) and the `r03` gap (≈ 150–250x inside X1, with `r03` pinned at 159–226
ops/s across every design and both orders) both survive a session change, and both are *ratios* whose
mechanism is in the query: no per-seat rows means no per-seat aggregate, and `r03`'s per-customer sort costs
what it costs. The absolutes moved — `q03` from 31.8k to 32–43k, `q05` from 5.5k to 12.4k on C1 — which is
the corpus's standing ratio-over-absolute lesson in a third setting.

**The X1 read penalty was a cell effect, and the old analysis was right to suspect one.** With the client
throttled ≤ 4.6 s for `x1` and `p3` alike, the buyer reads come out within 5% in one order and within 25% in
the other — not the 0.29–0.73 of the old run. The remaining difference is order-dependent, so it belongs to
session state and not to the ledger's read path; the analysis that flagged it, without being able to name a
cause, was pointing at the right thing.

**Prediction verdicts** (recorded in `V2-READ-FAMILY-RERUN.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| P1 | today's client is not throttled in these read cells | **falsified** — `r1` loses 420/634 and 422/633 periods; the others 2.3–5.9 s |
| P2 | `q05` is much faster on C1 than P3/X1 | **confirmed** — 12.2–12.7k against 2.0–2.1k and 1.9–2.0k |
| P3 | `q03` is ~100x `r03` inside X1 | **confirmed** — 189x and 169x; `r05` answered only by the ledger at 30.3k/33.0k |
| P4 | X1's ~2x read penalty does not reproduce | **confirmed** — ratios 0.77–1.05 and 0.97–1.05 against the old 0.29–0.73 |
| P5 | correctness gates stay green | **confirmed** — no failed cells, both trees clean |

## Where I think the measurement is weak

1. **Four designs of fifteen.** The other eleven, including the two other heavily throttled ones
   (`p4_precreated_counter`, `c4_counter_guard`), were not re-measured; `r1` being the extreme case is
   evidence, not a census.
2. **One client size.** The two-core client quota is the study's constant, so "client-bound" here means
   "bound at this budget"; nothing says how much headroom `r1` would need. That is the sizing calculation the
   owner's rule asks for and this tranche does not complete.
3. **The flat-`q01` argument assumes the query's cost should fall with tier.** It does for `p3` and `c1`
   over the same range, which is what makes `r1`'s flatness diagnostic; but a design whose read is genuinely
   O(1) would look the same.
4. **The YB half is untouched** — the original quoted 5.3 → 29.9 and 7.2 → 42.2 ops/s on `r03` for
   `yb-single` and `yb-cluster3`, and the 2 GiB machine ceiling still governs whether that can be re-measured
   at all.
5. **One host, one engine per cell, one session per order** — and the read phase's own 10-second windows are
   short enough that a client floor and a store ceiling are hard to separate without more trials per cell.

## What I would measure next

1. **Size the read client** before quoting any read level: raise the client's CPU budget (or cut its
   connections) until the throttling counters are zero for the heaviest design, and record the point at which
   `r1`'s `q01` starts to fall with tier. That is the measurement the read family is missing.
2. **Re-measure the other eleven designs** at the sized budget, so the read tables have one footing.
3. **Make a flat cell an automatic flag**: the report already prints client throttling beside each cell, and
   a cell whose value does not move across the tier axis while its neighbours do is a candidate client floor.

## Disagreements with other analyses

It corrects an implication of my own `study02-client-quota-audit` (that the hazard was that session's) while
confirming its finding, and it agrees with — and explains — the 2026-09-21 analysis's suspicion that X1's
read penalty was "a cell-level effect". It leaves `v2-07`'s answerability result untouched: `r05` is still
answered only by the ledger, at the same ~32k.
