---
analysis_id: book-evidence-v5-consistency--deepseek-flash--2026-09-24
run_id: 20260920T234953Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task book/evidence-v5-consistency."
analyzed_at: 2026-09-24
repo_commit: 317df0e970d06f7a42fb9b3a981093e34dc07ddf
supersedes:
status: current
headline: All 29 active claims audited for "the statement asserts more than its limits, trials or anchors support"; five were re-scoped in a v5 package, 24 passed unchanged, and the one the reviewer proposed to soften is instead kept with corrected limits because the study schema proves the mapping.
---

# Analysis — book evidence v5 (claim-versus-limit audit) — deepseek-flash

> Deliverable of task `book/evidence-v5-consistency`, implementing the third review round's finding
> class. This is an **audit and re-scoping** analysis: it measures nothing, adds no number and alters no
> measured value. It changes what five claims *say about themselves*.

## TL;DR

- The active registry moves from `book/evidence/v3`-derived v4 to `book/evidence/v5`. 24 of 29 claims
  are byte-identical to v4; five are re-scoped.
- The class audited is narrow and specific: **a claim whose statement, or one of whose limits, asserts
  more than the claim's own trials, anchors and limits allow.** It is the same class as the v4 defect
  (two gaps kept alive after their evidence had moved), one layer down.
- Five claims were affected, not four and not twenty: the reviewer found all five, and the sweep found
  no others. That is worth recording as evidence that their reading of the draft was careful.
- One disposition differs from the reviewer's proposal: `v3-06` keeps its statement and gets corrected
  limits, because the study's primary key makes the statement a configuration fact.
- Two observations are recorded without a change, so they are visible rather than silently fixed.

## 1. Method

For each of the 29 claims, three questions were asked against the committed JSON and the cited reports:

1. **Statement vs limits.** Does the statement assert a kind of conclusion — a cost, a direction, a
   physical outcome, a repetition — that one of the limits explicitly disclaims?
2. **Statement vs trials.** Does the statement or a limit cite a repetition count larger than the
   claim's declared `trials`, without a corroborating anchor?
3. **Limit vs corpus.** Does a limit read as a statement about the whole corpus when it is a scope limit
   on one claim's range, especially after a later package broadened the corpus?

Question 3 is the one that only becomes visible over time: v2-14's "no YugabyteDB, cluster…" was a
correct limit when written and became a false impression once the v3 ingest added exactly those cells.

## 2. Verdicts for all 29 claims

| Claim | Verdict |
|---|---|
| `v2-01-normalized-index-first` | passes — the blended-measure limit covers the ratios |
| `v2-02-rollup-family` | passes |
| `v2-03-embedding-not-ahead` | passes |
| `v2-04-recency-maintained-index` | passes — `three trials per cell` matches `fresh_load_trials_per_cell = 3` |
| `v2-05-d3-read-advantage-repeated` | passes — five fresh loads declared and cited |
| `v2-06-hot-drop-designs` | passes |
| `v2-07-refund-answerability-and-storage` | **changed** — limit cited repetition the claim cannot show |
| `v2-08-ledger-race-cost-128-buyers` | passes; observation on `inner_iterations` (below) |
| `v2-09-postgres-ledger-ratio-unsettled` | passes; same observation |
| `v2-10-guarded-confirm-refusals` | **changed** — statement asserted a throughput conclusion |
| `v2-11-section-seat-map` | passes — the card already labels itself a package comparison |
| `v2-12-lock-vs-version-upper-bound` | passes — the upper-bound wording is honest; "one installed product" is a book-side term |
| `v2-13-trigger-rollup-cost` | passes |
| `v2-14-cache-throughput-gain` | **changed** — limit read as a corpus-wide gap |
| `v2-15-three-instance-staleness` | passes — single trial and the swapped analysis rows are both stated |
| `v2-16-strict-freshness-read-cost` | **changed** — statement asserted a write-path cost |
| `v2-gap-01-ledger-load-multiplier-unstable` | passes — the two sibling runs it cites are declared in `trials` |
| `v2-gap-02-hold-funnel-unanswerable` | passes — schema limitation, correctly badged |
| `v2-gap-03-study04-missing-designs` | passes |
| `v2-gap-06-native-datastore-families` | passes |
| `v3-01-churn-crosses-real-ttl` | passes — the zero hard expiries are stated as an open clause |
| `v3-02-churn-600-through-vs-aside` | passes — two cells in separate runs, with the sibling anchor present |
| `v3-03-equal-total-framing-labelled` | passes — the difference is explicitly not attributed |
| `v3-04-yb-single-strict-relaxed` | passes — no winner attributed |
| `v3-05-yb-cluster3-strict-relaxed` | passes — opposite direction stated, no engine ranked |
| `v3-06-colocated-vs-noncolocated-locality` | **changed** — limits under-explained a statement the schema supports |
| `v3-gap-01-study05-remaining-regimes` | passes — partially superseded, remaining dimensions named |
| `v4-gap-01-physical-placement-unverified` | passes |
| `v4-gap-02-endpoint-distribution-partial` | passes |

## 3. The two dispositions that differ from the review

**`v3-06` keeps its statement.** The reviewer offered "is designed to keep one donor's donations
together" on the grounds that the card's limit says the placement labels are intent rather than
observation. That reasoning would be right if the statement described an *observed* placement. It
describes a *configured* one: `schema.sql` for `y1_colocated` declares
`PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)` and `y2_noncolocated` declares
`PRIMARY KEY (donation_id)`, so a donor's rows hash to one tablet in the first and spread in the second
by construction. The defect was the limit, which mentioned the missing observation without saying what
the mapping rested on. The statement is now backed by an explicit limit instead of being weakened.

**`v2-16` and `v2-10` lose the assertion rather than gain a measurement.** The reviewer's alternative —
register a performance claim — needs a measurement this task must not invent. Dropping an unsupported
conclusion is the project's own rule and is reversible: if the owner wants the write-path cost
quantified, it becomes a Study 02/05 measurement task and the claim can be extended then, on evidence.

## 4. Recorded observations, not changed

- **`inner_iterations` on `v2-08` (128) and `v2-09` (32)** reads as the buyer count rather than as repeat
  iterations inside a cell. If the Study 02 protocol defines `inner_iterations` as the buyer sweep, the
  field is correct and only its name is a hazard; if not, the field is being misused. I did not change it:
  establishing which is true needs the study's protocol definition, and guessing would repeat the
  failure this task repairs.
- **The `v2-07` repetition basis may be recoverable.** Two sibling Study 02 result directories exist. If a
  later pass verifies that each measured the same storage delta, they can be added as corroborating
  anchors and the three-run limit restored — on evidence this time.

## 5. Where this analysis is weak

- **The audit is one analyst's reading, not a rule.** Questions 1 and 3 in §1 are natural-language
  judgements; only question 2 is now mechanical (the new repetition-versus-trials lint). A second analyst
  reading all 29 claims might find a scoping overreach I did not.
- **The new lint is deliberately shallow.** It compares number words against the declared `trials` and
  the presence of a corroborating anchor. It catches the `v2-07` shape and nothing subtler; it cannot see
  a cost conclusion in a refusals claim, which is why `v2-10` and `v2-16` needed human reading.
- **I did not chase the storage delta into the sibling runs.** The reworded limit is honest either way,
  but recovering the repetition evidence would have been the stronger fix, and I stopped at the point
  where continuing was a study-level investigation rather than an evidence-registry task.
- **No report, result file, digest or earlier signed analysis was edited**, and no number was changed;
  the five changes are wording and attribution only. That also means any reader who quoted the old
  wording has a diff to consult rather than a silent edit — which is the point of the version bump.
