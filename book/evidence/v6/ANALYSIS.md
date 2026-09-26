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

- v6 began as a **method-metadata-only** package: all 29 v5 claims are carried forward with their
  statement, limits, confounds, anchors, support cells, trial counts, review state and lifecycle
  byte-identical. One claim was then re-scoped before the package was tagged, by the corpus's first
  run-level replication (the update below);
  `TestV6CarriesV5ForwardExceptTheReplicatedClaim` now pins exactly that one exception.
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
| 10 `isolate`, 9 `package`, 10 `not_applicable`, **1 `interaction`** | The book's evidence is strong on single-decision pairs and honest about packages. The first crossed-factor result arrived with v2-18 (maintenance strategy × replacement size, see the third update below); the other interactions the corpus names as unmeasured are still unmeasured. |
| 20 of 30 withhold attribution | This is the v5 confounds (conf-01…conf-06) working as designed. The book may report these packages, but any recommendation that names one component needs the isolation experiment. |
| 15 claims with `control_coverage: absent` | The correctness gate is only as strong as the control that fires. Seven are read-only phases where the study's controls could not fire; their correctness inference is now withdrawn rather than implied. |
| 27 single-run, 3 run-level (one a gap), 0 hardware-independent | The charter's replication-before-domains ordering is not a preference but the corpus's largest confidence gap. One claim was replicated before the v6 tag (see the update below). |
| 1 recorded prediction (falsified), 1 bounded crossover | Both fields are now exercised: the T1 replication falsified its prediction, and the writer sweep registers a crossing bracketed between 1 and 2 writers. Neither is replicated - one measurement per level. |

*(These counts are the ones this analysis was written with, corrected for the two updates below.)*

## Update (2026-09-25, second): the first bounded crossover

`v2-17-optimistic-pessimistic-crossover` is a **new** claim (the 30th), not a re-scope. It was created by
extending Study 04's writer sweep from four levels to seven — 1, 2, 4, 8, 16, 32 and 64 writers — which
is what the `crossover` lint requires before a claimed reversal may be presented: five levels and a stated
crossing interval. The earlier four levels gave a reversal between 1 and 4 writers; two adjacent levels
(1 and 2) now bracket it, since `c2/c1` goes 0.89 → 1.10 there.

Three consequences for the counts above: the claim total is **30**, one claim carries a `crossover`
record, and `crossover-bounded` moves from `gap` to `measured` in `coverage.json`. The addition is
declared in `TestV6CarriesV5ForwardExceptTheReplicatedClaim`, which now permits exactly one added claim
and one re-scoped claim and nothing else.

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

## Update (2026-09-25): the first replicated claim

The T1 replication landed after this analysis was written and before v6 was tagged, so v6 carries one
claim that is no longer identical to v5.

`v2-13-trigger-rollup-cost` cited a single run: "costs 6.85x on whole-configuration replacement (83.7
vs 573 ops/s)". Three whole runs with randomized order (order seeds 101/202/303, three trials per cell)
gave 3.99x, 4.94x and 5.69x — the direction replicated in 3/3 runs, the magnitude did not. The claim is
now:

- strength `run_replicated`, with `replication {runs: 3, order: randomized, across_run_spread_pct:
  34.4, tier: run-level}` and four anchors (three new runs primary/corroborating, the original run
  corroborating);
- `prediction` recorded with verdict **`falsified`** — the first claim to exercise that field, and the
  first to exercise the `run_replicated` branch of the v6 lints;
- statement re-scoped to the observed range with the trigger cell's own 35–115% trial spread stated,
  the 22–46x p99 gap added, and the ~15x read-gain clause marked as not re-measured.

Two findings this changes in the tables above: **one** claim is now `run-level` (27 single-run, not 28),
and **one** recorded prediction now exists (falsified, not zero). The coverage dimension
`run-level-replication` moves from `gap` to `measured` for that one claim and stays open for the other
28.

The book prose repeated the superseded 6.85x in four places (`minor-variants`, `how-to-choose`,
`scenarios`, `concepts/derived-state-and-history`); all four now state the replicated range.

## Next experiments

T1 (run-level replication with randomized order) and T2 (crossover/boundary sweeps) are drafted under
`.queue-drafts/empirical-method/`. T2's first deliverable should register Study 04's writer-sweep
crossing as the corpus's first `crossover` claim and Study 04's `d2` falsified prediction as its first
recorded `prediction`, since both already exist as measurements and only lack registration.
