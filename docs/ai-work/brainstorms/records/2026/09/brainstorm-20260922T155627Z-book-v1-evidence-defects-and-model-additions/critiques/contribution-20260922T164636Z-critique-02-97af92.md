# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `contribution-20260922T164636Z-critique-02-97af92` |
| Slot | `critique-02` |
| Actor | model=gpt-5 tool=openai-codex effort=unknown session=session-20260922T160724Z-a4fe79 capability=HIGH role=reviewer |
| Capability input | HIGH |
| Submitted | `2026-09-22T16:46:36Z` |
| Confidence | high |

## Summary

Both positions identify a publication blocker. Preserve registry v1 and create one attributed versioned correction as the exclusive book source; validate structured support cells, run status, controls, provenance and trial dimensions. Keep the fixed six-family taxonomy and lead the analytic additions with answerability-before-performance.

## Disagreements

A narrative audit is safe only as the correction ledger for the sole active evidence version; exact full-registry cell-resolution counts still need mechanical regeneration; answerability-first leads by editorial judgement, not measured ranking.

## Contribution

# Cross-review — version the correction, keep one active evidence source

## Verdict

Both positions identify the right publication blocker: the current registry cannot safely
drive numeric book prose. The strongest defect is not the omitted tags by themselves but the
fact that a schema-valid cell map can identify the wrong design and conceal the source run's
failed or partial cells. The strongest model addition is answerability-before-performance,
provided it is labelled as a decision rule derived from two measured answerability failures,
not as a measured ranking.

The positions should be combined, but neither proposed correction mechanism is sufficient on
its own. Position 01's free-standing audit would become a second source of truth if the book
continued to consume `book/evidence/claims.json`. Position 02 and critique 01's in-place repair
would erase the distinction between the accepted v1 registry and its correction. The evidence
packet resolves this apparent ambiguity: the registry is an immutable v1 milestone at
`repo/book-evidence-registry-v1`; the brainstorm boundary says a correction is a new,
attributed artifact citing the old one; and the book build records an evidence-registry digest.
The admissible combination is therefore a **versioned evidence-correction package that is the
book's sole active claim source**, while v1 remains unchanged and cited as the corrected
artifact.

No item below requires a database, build, new design, new study, or measurement.

## Findings confirmed and priority

### 1. Block direct book consumption of registry v1

The positions and critique 01 converge on literal contradictions:

- Claim 01's seven curator names do not occur in its report; worse, D1-D3 are shifted and D4/D5
  are reversed relative to the source design identities, so the declared winner can name the
  wrong mechanism.
- Claim 04's five names do not occur in its report. The report says 42 cells and 2 failed while
  the claim says `failed_cells: 0`; the validator only reconciles that value with the claim's
  own selected array.
- Claim 06 says no under-booking occurred, but its cited race run reports `0 of 0` negative
  controls. “None observed” is support for that workload, not demonstrated detector power in
  that run.
- Claim 12 narrows its stated range to -8.3%..+12.3% although the cited analysis table includes
  -14.2%.
- The corrected tag count is 9 of 12 numeric claims with `run_tag: null` despite an extant tag,
  not position 02's 11 of 13. Claim 01 is the separate legacy-provenance case: no producing tag
  or commit can be recovered from its chain without invention.

These are direct evidence defects. A reader can be led to the wrong design, an incomplete view
of run status, or a false precision range. They outrank the conservative `trials` issue, which
can understate evidence rather than misdirect it.

### 2. Use a new versioned registry plus a correction ledger, not a parallel narrative audit

The synthesis should recommend a new attributed evidence version (for example,
`book/evidence/v2/claims.json` with a signed correction ledger that cites
`repo/book-evidence-registry-v1`). The exact path is an implementation choice for a later
owner-authorized task; the invariant is what matters:

1. v1 is never edited and remains resolvable by its tag;
2. each changed or split claim records the v1 claim id and correction reason;
3. the book and build manifest consume and digest exactly one active version;
4. a claim cannot silently fall back to v1;
5. validation is upgraded against committed source artifacts before the corrected version is
   accepted.

This answers critique 01's frozen-registry question without requiring an owner choice between
two unsafe extremes. It also avoids re-planning the completed registry task: this is a narrowly
scoped correction to defects found after acceptance, not a second design of the original task.
If versioned registry support proves impossible, position 01's signed audit is acceptable only
if it becomes the **sole** book input and v1 is excluded from the prose pipeline. Publishing
both remains unacceptable.

**Cost:** inspect committed artifacts, create one versioned claim set and correction ledger,
extend tests/validator, and point the book build at one digest. **Must not require:** edits to
v1, reports, analyses, discussions, results, manifests or tags; a rerun; guessed provenance; or
stronger evidence labels.

### 3. Strengthen semantic validation, not string matching alone

Critique 01's proposed rule that every `cells[].name` occur in the report is a useful smoke
test but not a sufficient resolver. A coincidental string match does not establish the measured
cell, and forcing a selected claim's `failed_cells` to equal the whole report's failed count
would conflate two scopes. The corrected schema/validator should instead require:

- a structured support key such as topology + design + scenario/operation + trial or aggregate,
  proven to exist in the cited report/result artifacts;
- separate `run_cells` totals/status from the atomic claim's `support_cells`, including partial,
  skipped, diagnostic, negative-control and failed outcomes;
- winner eligibility only from the exact valid support cells;
- negative-control status as same-run fired, present-but-not-exercised, inherited, absent, or
  not applicable;
- every extant producing tag and commit, or an explicit `legacy_provenance_incomplete` reason;
- multiple anchors for conclusions that genuinely depend on several runs;
- distinct counts for whole-run replication, fresh-load/cell trials and inner race iterations.

Claim 02's `trials: 1` versus three fresh-load trials proves that the last concept is ambiguous,
but it does not by itself prove the registry is stronger than the evidence: `1` may mean one
matrix execution. Clarify the dimensions and retain the conservative strength until the source
supports promotion.

**Reader effect:** a green validator would mean “this exact atomic claim resolves to this exact
eligible evidence,” not merely “this JSON is internally consistent.”

### 4. Carry comparability and completeness qualifications beside each number

Position 01's confound list should survive synthesis:

- D2 to D3 prices copied key + SQL + indexes, not rolldown alone.
- D20 to D22 is informative for read placement/indexing but confounded for write cost by D22's
  larger rollup-maintenance package.
- Study 03 row/document comparisons are package comparisons, not a universal document effect.
- Study 05's cache gain is add-cache, db-only resource framing, not equal-total resource
  efficiency.
- Study 02's 32- and 128-buyer P3/X1 arms disagree in direction; the PostgreSQL conclusion is
  unsettled/order-sensitive, not a pooled ratio.
- Claim 12 must show the full -14.2%..+12.3% observed single-trial/noise range.

The book also needs the proposed coverage matrix, with `measured`,
`measured-but-confounded`, `planned`, `gap`, and justified `not_applicable`, plus review state.
That table is disclosure, not execution of the already queued reviews or scientific protocols.

### 5. Apply one placement standard to all three studies

Critique 01's consistency challenge is correct and can be resolved from the current corpus:

- Study 01 D3/D7 and D20/D24 directly support a primary-key/layout change and observed access
  behavior; a saved plan showing one storage request versus two is mechanism evidence, not
  tablet-level verification of physical colocation. The D20/D24 analysis calls the performance
  difference a null result but supplies no actual-placement evidence.
- Study 02 already records physical colocation as an untouched gap.
- Study 03 L3 supports section hash-key/sharding intent and observed performance; its analysis
  says tablet balance is unexplained.

Therefore none should be advertised as satisfying the current colocated/non-colocated minimum.
Label Study 01 and Study 03 as layout/placement intent with measured consequences, and retain
physical placement as a common gap. This does not erase their useful measurements.

## Model additions: keep the fixed taxonomy, add decision axes

Position 02's claim that four new measured **major families** sit outside the book's six
conflicts with the goal request: `REQUEST.md` explicitly fixes the six major families and
already lists optimistic/pessimistic concurrency, locks/CAS/isolation and inventory creation as
minor variants; it also already reserves an `expiry-and-clock-authority` concept chapter. The
evidence is valuable, but the classification must change:

1. **Answerability-before-performance** leads as a decision rule. Direct evidence: Study 02
   `r05` is unanswerable in 14 of 15 measured designs and X1 prices the one measured ledger
   remedy; Study 03 `r06` is unanswerable in every measured design. A hold-event ledger for
   Study 03 remains analogy/gap. Phrase the rule as derived from those failures, not measured
   superiority of “asking first.”
2. **Multidimensional design vector** maps the fixed major family plus information placement,
   maintenance owner, arbitration, expiry/clock authority, unit of change, cache contract,
   deployment and workload. This absorbs position 02's proposed arbitration, expiry,
   aggregate-boundary and derived-current-value “families” without changing the fixed taxonomy.
3. **Derive → index an existing fact → materialize → copy/cache ladder** is the strongest
   when-to framework, with each step naming atomicity, maintenance owner, mutation coverage,
   detector, storage and write cost. Evidence is direct within Studies 01, 04 and 05 and analogy
   across them.
4. **Contention surface and arbitration** combines position 01's blast-radius model with
   position 02's measured arbiters. It is a minor-variant/decision chapter, not a seventh major
   family, and it must not become a universal lock-over-CAS rule.
5. **Query-boundary portfolio, evidence-confidence cards and failure-mode catalogue** remain
   high-value analytic models. They teach readers which query a layout helps, when evidence is
   ineligible or merely directional, and which lost fact or broken invariant caused failure.

**Cost:** taxonomy mapping, decision tables and prose over committed designs. **Must not
require:** a cross-study numeric ranking, a claim that every design-vector combination was
measured, an upgrade in evidence strength, or presenting an unbuilt Study 03 ledger as direct
evidence.

## Recommended synthesis order

1. Declare registry v1 ineligible as the direct numeric source and require one versioned,
   attributed correction package as the exclusive book input.
2. Atomize claims and add source-resolving, run-status, control-scope, multi-run provenance and
   trial-dimension validation.
3. Publish the confounds, placement/common-methodology gaps, and missing native-datastore plus
   real-network/balanced-endpoint gap disclosures.
4. Lead the decision model with answerability-before-performance, explicitly labelled a derived
   rule.
5. Preserve the fixed six-family taxonomy and add the design vector, derive/index/materialize
   ladder, contention/arbitration, query-boundary and benchmark-reading models beneath it.

## Remaining dissent and uncertainty

- A signed narrative audit may still be useful to humans, but it should be the correction
  ledger attached to the one active evidence version, not an independently maintained claim
  database.
- Exact counts beyond the independently checked claim-01 and claim-04 mappings should be
  regenerated mechanically before the correction artifact is published; critique 01 did not
  independently reproduce position 01's full 11-of-12 cell-resolution count.
- The current corpus establishes that actual placement is not documented to today's standard;
  it does not prove the intended layouts failed to colocate data.
- “Answerability first” is strongly supported as a reasoning rule, but its ordering over all
  other decision rules is editorial judgement, not a measured comparison.

## Confidence

High on the publication blocker, versioning resolution, fixed-taxonomy constraint, placement
qualification and ordering of misleading versus conservative defects. Medium-high on the
combined analytic-model ordering because it is editorial synthesis over direct per-study
evidence rather than a measured cross-study ranking.
