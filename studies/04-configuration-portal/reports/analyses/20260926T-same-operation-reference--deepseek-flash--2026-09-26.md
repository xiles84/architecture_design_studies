---
analysis_id: 20260926T-same-operation-reference--deepseek-flash--2026-09-26
run_id: 20260926T205817Z, 20260926T210029Z, 20260926T210243Z, 20260926T210436Z, 20260926T210636Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/same-operation-reference."
analyzed_at: 2026-09-26
inputs_digest: per-run 4b46fc9abeb565cd (reference 1), 1deeda3df6421186 (tier 10), dceeeda9bdd7951b (reference 2), 2060e361b6abbbbd (tier 240), d9632ac797c6e714 (reference 3)
repo_commit: 71517b8
status: current
headline: A reference cell on the same operation shows the ratio is more portable than the absolutes (40% spread against 49-114%), the session wanders rather than drifts in one direction, and one tier-10 reading fell outside the band six earlier runs had established.
---

# Analysis — Study 04 same-operation reference: the ratio is portable, the session wanders — deepseek-flash

## TL;DR

- **The ratio travels better than the absolutes, as predicted.** Over the three reference readings the
  `n4/n3` ratio spread **40.2%**, against **49%** for `n4`'s absolute throughput and **114%** for `n3`'s.
  A session term shared by both designs cancels in the ratio, so the ratio is the quantity to carry
  between cells and the absolute numbers are the quantity to distrust. S1 confirmed.
- **The session does not drift; it wanders.** The three reference readings went **8.64 → 5.21 → 8.53** —
  down and then back up, not monotone. S3 is falsified, and it falsifies the *model* the previous two
  tranches were using: "drift by time" implied a trend an endpoint reading could characterise, and this
  sequence is not a trend.
- **One reading fell outside a band built from six runs.** Tier 10 measured **1.06**, below the
  1.34–3.15 range that the sequential and interleaved tranches had established, so S2 is falsified for
  tier 10. A band from six runs was not wide enough to contain a seventh reading. Tier 240's 24.17 sat
  just inside its band (10.20–24.20).
- **Correctness held in all five runs**: 0 audit failures, 0 rollup mismatches. S4 confirmed.
- **`v2-18`'s tier-10 band widens to 1.1–3.2x**, and the claim's three levels now rest on 7, 9 and 7 runs
  with ratio spreads of 101%, 135% and 67%.

## The tranche

| Position | Run | Cell | n3 ops/s | n4 ops/s | n4/n3 |
|---|---|---|---|---|---|
| 1 | 20260926T205817Z | reference, tier 60 | 22 | 190 | 8.64 |
| 2 | 20260926T210029Z | comparison, tier 10 | 225 | 239 | **1.06** |
| 3 | 20260926T210243Z | reference, tier 60 | 71 | 370 | 5.21 |
| 4 | 20260926T210436Z | comparison, tier 240 | 6 | 145 | 24.17 |
| 5 | 20260926T210636Z | reference, tier 60 | 43 | 367 | 8.53 |

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| S1 | the ratio's spread is smaller than either absolute's spread | **confirmed** | 40.2% against 49% (`n4`) and 114% (`n3`) |
| S2 | both comparison cells fall inside the established bands | **falsified at tier 10** | 1.06 against a 1.34–3.15 band; tier 240's 24.17 was inside |
| S3 | the reference readings move monotonically with session time | **falsified** | 8.64 → 5.21 → 8.53 |
| S4 | correctness holds in every run | **confirmed** | 0 audit failures, 0 rollup mismatches in five runs |

## What this changes

1. **`v2-18`'s bands and counts.** Pooling every run of each cell: tier 10 spans **1.06–3.15** over 7 runs
   (spread 101%), tier 60 **2.51–8.95** over 9 (spread 135%), tier 240 **10.20–24.20** over 7 (spread 67%).
   The ordering still holds in every individual arrangement, and the two lower bands still touch when
   everything is pooled — 3.15 against 2.51.
2. **The drift model is corrected, again.** The previous lesson said a session drifts and endpoints measure
   it. The endpoints here would have reported "8.64 at the start, 8.53 at the end, no drift" while the
   middle reading sat 39% below both. **A non-monotone instrument needs readings inside the tranche, not
   only at its ends** — which is exactly why this design took three, and the three are what caught it.
3. **A six-run band is not a bound.** One in seven readings fell outside it. Bands in this corpus should be
   described as what was seen, with the run count beside them, rather than as a range a future reading will
   fall inside.

## Where this analysis is weak

1. **Three reference readings cannot establish the instrument's reliability**, only that it disagreed with
   a simpler model within one session.
2. **The comparison cells were measured once each**, so their placement against the reference is a single
   observation — and tier 10 is precisely the one that fell outside, which is one observation doing a lot
   of work.
3. **The +44% contention reference from the previous tranche and this tranche's −39% mid-tranche excursion
   are not comparable** quantities: one is a different operation, the other the same one. Together they
   show movement exists; neither quantifies "the" drift.
4. **One host, one day.** Nothing here says the wandering is a property of this machine rather than of
   this day.
5. **No mechanism.** Cache state, thermal behaviour, container reuse and background load all predict
   movement; none is tested.

## Next experiments

1. **Take the reference reading between every run**, not three times per tranche, so the wandering becomes
   a covariate beside each measured cell rather than a surrounding symptom.
2. **Seek the mechanism**: record the host's CPU frequency and load beside each run, so the movement can be
   attributed to something rather than merely observed.
3. **Repeat a single cell many times in one session** (say ten runs of tier 10) to characterise the
   distribution instead of its extremes, which is what the bands currently are.
