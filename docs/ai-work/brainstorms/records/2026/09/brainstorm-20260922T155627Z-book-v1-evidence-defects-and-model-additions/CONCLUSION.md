# Brainstorm conclusion

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `contribution-20260922T170212Z-synthesis-01-7465d9` |
| Slot | `synthesis-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=unknown capability=HIGH role=reviewer |
| Capability input | leader |
| Submitted | `2026-09-22T17:02:12Z` |
| Confidence | high |

## Summary

Registry v1 is not eligible as the book's numeric source: 11 of 12 numeric claims have at least one declared cell name absent from the cited report and 36 of 43 names (84 per cent) resolve nowhere, so one versioned, attributed correction package becomes the sole active claim source while v1 stays immutable at repo/book-evidence-registry-v1. The correction must atomize claims and carry source-resolving support keys, run-level versus support-cell status, same-run control scope, multi-run anchors, trial dimensions and a legacy_provenance_incomplete mark for claim-01; the six confounds, the coverage matrix, placement as a common gap and the two missing gap records are published beside the numbers. The fixed six-family taxonomy is kept, with answerability-before-performance leading the decision models as a derived rule. No measurement, no new study, no tasks.

## Disagreements

Preserved: a human-facing narrative audit is admissible only as the correction ledger of the one active evidence version; answerability-first is editorial ordering, not a measured ranking; the block is a claim-to-cell resolution failure, not a claim that v1's statements are false; physical placement is undocumented to today's standard, not proven uncolocated; and this synthesis was written by the session that authored position-02 and critique-01, so the owner's later review should weigh that reduced independence. The mechanical count that critique-02 asked for is now computed: 11 of 12 claims affected, 36 of 43 names unresolved, and the corrected tag count is 9 of 12.

## Contribution

# Brainstorm conclusion

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `CONCLUSION.md` · slot `synthesis-01` |
| Synthesizer | deepseek-flash, WSL2 shell + `tools/queue`, effort not exposed by the client; capability `HIGH` (raw input `leader`), role `analyst` |
| Submitted | UTC timestamp recorded by the CLI at submission |
| Confidence | **high** on the publication blocker and the versioning resolution; **medium-high** on the ordering of the analytic models, which is editorial synthesis rather than a measured ranking |

## Decision

1. **Registry v1 is not eligible as the book's numeric source.** `book/evidence/claims.json`
   (milestone `repo/book-evidence-registry-v1`) stays immutable and cited; **one versioned,
   attributed evidence-correction package becomes the sole active claim source**, and the book
   build digests exactly one version. No numeric book prose is written from v1.
2. **That package must atomize claims and carry source-resolving support** — a structured support
   key (topology + design + scenario/operation + trial/aggregate) proven to exist in the cited
   report, run-level cell totals and status kept separate from the claim's support cells,
   same-run negative-control scope, every extant producing tag and commit or an explicit
   `legacy_provenance_incomplete` reason, multiple anchors for multi-run conclusions, and distinct
   counts for whole-run replication, fresh-load trials and inner iterations.
3. **Confounds and coverage are published beside each number**, not in a distant caveat: the six
   confounds below, a coverage matrix (`measured` / `measured-but-confounded` / `planned` /
   `gap` / justified `not_applicable`) with review state, and the missing native-datastore and
   real-network/balanced-endpoint disclosures. Physical placement is stated as an open gap for
   all three studies that gesture at it.
4. **The fixed six-family taxonomy is kept.** Model additions are concept chapters and decision
   axes beneath it, in this order: answerability-before-performance (labelled a rule *derived*
   from two measured answerability failures), the multidimensional design vector, the
   derive → index-an-existing-fact → materialize → copy/cache ladder, contention
   surface/arbitration, query-boundary portfolio, evidence-confidence cards, failure-mode
   catalogue.
5. **Nothing is measured, and no task is created.** Every item is registry, validator or book
   work over committed artefacts.

### The six confounds to publish beside each number

1. Study 01 D2→D3 prices a copied key **plus** its SQL and indexes, not rolldown alone — the
   repeated analysis says copying the key by itself did not reproduce the gain.
2. Study 01 D20→D22 is a clean read-side location/index comparison but a confounded **write**
   comparison, because D22 also maintains D4's whole multi-column, two-parent rollup package.
3. Study 03's row-versus-document seat-map contrast is a family/package comparison, not a
   one-decision pair; present the scan/point-lookup trade for those complete layouts only.
4. Study 05's 2.2–3.5× is an **add-cache** result under `resource_framing: db-only`; it is not
   equal-total resource efficiency, and the missing Redis budget is an explicit gap.
5. Study 02's 32-buyer and 128-buyer P3/X1 arms disagree in sign; keep "PostgreSQL
   unsettled/order-sensitive" rather than pooling or selecting a ratio.
6. Claim-12 must show the full observed single-trial range **−14.2 %…+12.3 %** and keep calling it
   noise — never "strict freshness is free".

## Why — applying the decision criteria

**Criteria 2 and 5 (evidence first; honest strength).** The blocker is not a matter of taste: it
is computed. Across the twelve numeric claims, **11 have at least one declared cell name that
does not occur anywhere in the report they cite**, and **36 of the 43 declared names (84 %) do not
occur at all**: claim-01 0/7, claim-02 0/2, claim-03 0/2, claim-04 0/5, claim-05 0/3, claim-06
0/3, claim-07 0/2, claim-08 0/3, claim-09 0/3, claim-10 0/4, claim-11 3/5 (the failures are
exactly `postgres-control` and `redis-control`), claim-12 4/4. This mechanically confirms
position-01's 11-of-12 count, which critique-02 asked to have regenerated rather than asserted,
and it retires critique-01's partial 2-of-12 check. Claim-01's mapping is not merely renamed: it
shifts D1–D3 by one relative to the report's `D1 minimal / D2 indexed / D3 flat+FK` labels and
**reverses D4/D5** (`d4-rollup-application` where the corpus's D4 is rollup *trigger*), so its
declared winner can lead a reader to the wrong mechanism. The same field reports
`failed_cells: 0` for a run whose own report says "Cells: 42 … 2 failed", and it narrows
claim-12's observed range to −8.3 %…+12.3 % while the cited table also contains −14.2 %. The
corrected tag count is **9 of 12 numeric claims** omitting an extant tag (not position-02's
"11 of 13", which conflated gap records and miscounted the numeric set; position-02 is immutable,
so critique-01 and this conclusion are the correction of record).

**Criterion 1 (everything derivable from committed artefacts).** Verified in passing and worth
recording, because the book's reproduction chapter depends on it: run tags mark the *producing
code state*, not the results. For `20260921T-survey3`, `20260921T1215Z-small` and
`20260921T205212Z` the tag is an ancestor of the commit that added the results, and for
claim-11 the tag commit **equals** the registry's `repo_commit` (`48fe1e33…`, with the results
committed later in `372f437`). A check of the form "does the tag's tree contain the results"
therefore fails on 16 of 17 tags **by construction and is not a defect**. Filling the nine
missing `run_tag` values is sound, and the book must say that a checked-out run tag does not
contain that run's result files.

**Criteria 3 and 4 (smallest effective set; immutability).** The versioning decision resolves the
question critique-01 could not settle from the evidence packet, using the boundary rule that a
correction is a new attributed artefact citing the old one. It is a narrow post-acceptance
correction, not a second design of the completed registry task, and it needs no edit to v1,
reports, analyses, discussions, results, manifests or tags.

**Criterion 6 (no new measurement).** Nothing selected requires a run, a build, a design or a
database start. Study 01's missing commit for claim-01 is published as
`legacy_provenance_incomplete`, never manufactured.

## Rejected alternatives

| Rejected | Why |
|---|---|
| Position-01's free-standing signed audit alongside v1 | Two active sources of the same facts, the failure `LESSONS_LEARNED.md` records for duplicated field names. Survives only as the **correction ledger** attached to the one active version. |
| In-place repair of `claims.json` (position-02 / critique-01 framing) | Erases the distinction between the accepted v1 milestone and its correction. Withdrawn. |
| Position-02's four new **major families** (arbitration, expiry/clock authority, unit of change, derived current value) | Conflicts with the owner's fixed six-family taxonomy in `REQUEST.md`, which already treats optimistic/pessimistic control, expiry and clock authority, derived-state maintenance and answerability as concept chapters or minor variants. Reclassified as decision axes beneath the fixed taxonomy; the *evidence* is unchanged and valuable. |
| Declaring all affected claims ineligible and blocking their prose | Over-broad: claim-12 fully resolves, and 10 claims are repairable. Only claim-01 needs the legacy-provenance label. |
| Critique-01's rule "`failed_cells` must equal the report's failed count" | Scope-conflating, as critique-02 showed: the report's failures are run-level while a claim's cells are a selected subset. Replaced by separate run-level and support-level fields. |
| "The aliases are close enough", or publishing v1 and the correction together | Defeats exact provenance, or restores two sources. |
| Re-running Study 01 so claim-01 gains a tag | A new measurement; out of scope, and unnecessary when the gap can be labelled. |

## Remaining dissent

- **Preserved from critique-02:** a human-facing narrative audit may still be worth having; it is
  admissible only as the correction ledger of the one active version, never as an independently
  maintained claim database.
- **Preserved from critique-02:** "answerability first" is supported as a reasoning rule, and its
  placement above every other rule is editorial judgement, not a measured comparison. The
  conclusion adopts it with that label rather than quietly upgrading it.
- **Preserved from critique-01:** the block on v1 is a *scope* decision, not a claim that v1's
  statements are false. Most v1 statements are copied verbatim from signed analyses and are
  probably right; what fails is the resolution from claim to measured cell.
- **Preserved:** the corpus shows actual placement is *undocumented* to today's standard
  (Study 01 D7/D24, Study 02's own recorded gap, Study 03 L3 whose analysis calls tablet balance
  unexplained). It does not show that the intended layouts failed to colocate data; both readings
  must stay visible.
- **Not resolved, and not blocking:** whether the versioned-package mechanism is implementable in
  the current build. The fallback is stated under review triggers.

## Assumptions, gaps and transfer risks

- **Independence caveat, stated plainly:** this synthesis was written by deepseek-flash, the same
  session that authored `position-02` and `critique-01`. The workflow prefers a synthesizer that
  is not one of the position authors, and the plan assigned this slot to GPT-5 Sol, which was
  active but not available in this session; the CLI had already handed the slot here and the
  archive records every actor, so the authorship is visible rather than hidden. The conclusion
  therefore leans on the two computed checks and on critique-02's independent resolutions rather
  than on its author's preferences, and the owner's later review should weigh it accordingly.
- **Coverage of verification:** the 43-name resolution table, the tag/ancestry checks, the
  claim-04 report line, the claim-12 table rows and the claim-02 trial sentence were verified
  directly. Position-01's remaining claim-by-claim findings (the claim-02/06/07 trial fields
  beyond claim-02, the claim-05…10 operation-versus-design mixing, the claim-08 multi-run
  reliance) were **not** independently reproduced here and are carried on their author's evidence.
- **Scope of reading:** neither position nor either critique read all 56 reports; a defect visible
  only inside an unopened report body remains invisible to this conclusion.
- **Transfer risk:** all evidence remains one laptop, one environment
  (`host-zenbook-ux5406sa`), heterogeneous cores and no real network. Nothing decided here changes
  that, and the book must keep saying so.

## Falsification and review triggers

1. A committed one-to-one mapping from registry names to report labels would reduce the cell
   finding to a documentation gap; none was found, and `claims.json`'s own `generated_from` note
   describes the cells as curated.
2. If versioned registry support proves impossible, the fallback is position-01's audit **as the
   sole input** with v1 excluded from prose. Publishing both reopens this conclusion.
3. If `trials` is defined in the registry as whole-run replications only, claim-02's `trials: 1`
   is correct and only the field's documentation is defective — a much smaller finding.
4. A committed placement-verification artefact (or an explicit statement that none exists) for
   Study 01 D7/D24 or Study 03 L3 closes or confirms the placement gap.
5. Reopen if `repo/book-evidence-registry-v1` is moved or re-tagged, or if a v2 appears without
   the correction ledger and its cross-reference to v1.

## Task-ready next steps

Described only; the owner must separately issue `CREATE TASKS FROM BRAINSTORM <id>` before any of
these may be published or claimed.

1. **Declare v1 ineligible as the book's numeric source** and point the book build's evidence
   digest at exactly one active version (with the versioned package's location decided in that
   task, not here).
2. **Build the versioned evidence-correction package plus its correction ledger**, citing
   `repo/book-evidence-registry-v1`, splitting compound claims, and carrying the structured
   support key, run-level status, control scope, multi-run anchors, trial dimensions and the
   `legacy_provenance_incomplete` mark for claim-01.
3. **Upgrade the schema and `tools/evidence` validator** so a green run means "this atomic claim
   resolves to this eligible evidence", including the 43-name resolution test as a regression
   check and the corrected tag resolution (9 numeric claims plus 3 gap records).
4. **Add the disclosure chapter**: the six confounds, the coverage matrix with review state,
   physical placement as a common gap, and the two missing gap records (native datastore families,
   real-network/balanced endpoints).
5. **Write the decision-model chapters** in the order decided above, each labelled direct
   evidence, analogy or gap, with no cross-study numeric ranking and no upward strength change.
