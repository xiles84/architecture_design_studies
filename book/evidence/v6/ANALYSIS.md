---
analysis_id: book-evidence-v6-method-metadata--deepseek-flash--2026-09-25
run_id: none (no measurement)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task repo/empirical-method."
analyzed_at: 2026-09-25
repo_commit: 0fc710f
supersedes: book-evidence-v5-consistency--deepseek-flash--2026-09-24 (the v5 text remains in book/evidence/v5/ANALYSIS.md; this file is the v6 package's own analysis)
status: current
headline: v6 adds evidence_form, attribution, control_coverage and replication-tier metadata to all 29 claims without changing one statement or limit, and the pass shows the corpus has 0 registered interaction claims, 0 recorded predictions, 0 bounded crossovers, 20 of 29 withheld attributions and 15 claims with no fired control.
---

# Analysis — book evidence v6 (method-metadata pass) — deepseek-flash

## TL;DR

- v6 is a **method-metadata-only** package: all 29 v5 claims are carried forward with their statement,
  limits, confounds, anchors, support cells, trial counts, review state and lifecycle byte-identical,
  pinned by `TestV6CarriesV5ForwardUnchanged`.
- It adds the metadata the empirical-architecture charter requires: `evidence_form`, `attribution`,
  `control_coverage`, and a `replication` record, plus optional `prediction` and `crossover`.
- Four rules are now machine-checked with failing fixtures: no attribution from a package; a confounded
  claim must withhold attribution or name its isolation pair; an unfired control is not coverage; a
  reversal must be bounded by five levels and a crossing interval.
- Honest findings about the corpus: **0** interaction claims, **0** recorded predictions, **0** bounded
  crossovers, **28** of 29 claims carried by a single run, **20** of 29 with withheld attribution, and
  **15** claims (7 of them read-only phases) with no fired control.

## Scope

This analysis covers the registry only. No run was executed, no measurement was attempted, and no
claim's content was changed. Its provenance is the v5 package and the migration mapping in
[`method_migration.json`](method_migration.json).

## Why the fields exist

The corpus's methodology already required isolation pairs (6a), regimes (8a) and negative controls
(5a). What it could not do is *state* which form a comparison took, or refuse a claim that attributed a
package to one component. v6 makes those statements machine-checkable, so a future claim cannot quietly
re-introduce the v1 registry's failure mode — a claim whose name resolved to nothing — in a subtler
form: a claim whose headline names a mechanism its own evidence cannot isolate.

## The findings, and what each implies

| Finding | Implication |
|---|---|
| 10 `isolate`, 9 `package`, 10 `not_applicable`, **0 `interaction`** | The book's evidence is strong on single-decision pairs and honest about packages, but it has never registered a crossed-factor result. Every interaction the corpus names as unmeasured is genuinely unmeasured. |
| 20 of 29 withhold attribution | This is the v5 confounds (conf-01…conf-06) working as designed. The book may report these packages, but any recommendation that names one component needs the isolation experiment. |
| 15 claims with `control_coverage: absent` | The correctness gate is only as strong as the control that fires. Seven are read-only phases where the study's controls could not fire; their correctness inference is now withdrawn rather than implied. |
| 28 single-run, 1 run-level, 0 hardware-independent | The charter's replication-before-domains ordering is not a preference but the corpus's largest confidence gap. |
| 0 recorded predictions, 0 bounded crossovers | Both exist in the corpus — Study 04's falsified `d2` prediction, the writer-sweep crossing between `c1` and `c2` — but neither reaches the book's numeric source because neither is a registered claim. |

## Where this analysis is weak

- The `evidence_form` and `attribution` values are *my* classification of what each v5 claim records.
  The mapping file makes every assignment auditable, but a second analyst could reasonably classify
  Study 02's family comparison or Study 05's add-cache result differently; where two readings differ,
  that disagreement should become its own signed document, not an edit to this one.
- `control_coverage: absent` conflates two situations: a read-only phase whose controls could not fire,
  and a claim whose run simply had no control. The distinction is visible in
  `run_level.controls_present/controls_fired` and in each claim's note, but it is not a separate field.
- A `run-level` tier is recorded for the one claim with multiple whole runs, but those runs were not
  ordered independently for that claim's question, so the tier describes run count, not replication
  quality.
- Nothing here is compile-verified for the book's Typst layer; the eight-question checklist is a
  book-owned Markdown contract (`book/concepts/CHECKLIST.md`) checked by `book/evidence/check.sh`, and
  its integration into the Typst chapters still needs a Podman/Typst session.

## Next experiments

T1 (run-level replication with randomized order) and T2 (crossover/boundary sweeps) are drafted under
`.queue-drafts/empirical-method/`. T2's first deliverable should register Study 04's writer-sweep
crossing as the corpus's first `crossover` claim and Study 04's `d2` falsified prediction as its first
recorded `prediction`, since both already exist as measurements and only lack registration.
