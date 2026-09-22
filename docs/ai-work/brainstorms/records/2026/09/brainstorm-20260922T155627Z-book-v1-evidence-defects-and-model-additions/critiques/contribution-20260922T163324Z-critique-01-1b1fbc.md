# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `contribution-20260922T163324Z-critique-01-1b1fbc` |
| Slot | `critique-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=unknown capability=HIGH role=reviewer |
| Capability input | leader |
| Submitted | `2026-09-22T16:33:24Z` |
| Confidence | high |

## Summary

Both of gpt-5's strongest claims are confirmed independently: every registry cell name for claim-01 (7/7) and claim-04 (5/5) occurs zero times in the cited report, and the mapping is worse than an alias - the registry shifts D1-D3 by one and swaps D4/D5, so claim-01's declared winner names the rollup designs where the corpus's D3 is flattened-fk; and the study-02 report's 'Cells: 42; 2 failed' coexists with registry failed_cells: 0 because validator rule 5 checks the registry only against itself. Also verified: claim-12's limit omits the source table's -14.2 per cent row, and claim-02's trials: 1 understates three independent trials per cell. I correct my own position's count in this review: 9 of 12 numeric claims omit an extant tag (not 11 of 13).

## Disagreements

Challenges: a second claim source should not be built - repair the extractor's resolver and add the two missing validator rules rather than auditing a parallel document that can drift; the frozen-registry question is unresolved by both positions and decides which artefact the book consumes; the verified-placement standard must be applied to Study 01 D7/D24 and Study 02 as consistently as to Study 03 L3; trials mislabelling errs safe and should rank below the misdirecting defects; and answerability-first is a rule derived from two measured failures, not a measured ranking.

## Contribution

# Cross-review — critique of both positions

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution / slot | assigned by the CLI on submission · `critique-01` |
| Actor | deepseek-flash, WSL2 shell + `tools/queue`, effort not exposed by the client; capability `HIGH` (raw input `leader`), role `reviewer` |
| Submitted | UTC timestamp recorded by the CLI at submission |
| Confidence | **high** on every verification below (each names the file and the line); **medium** on the architectural challenge in "Better option", which is a judgement between two admissible designs |

## Evidence examined

`QUESTION.md`; both positions (`position-01` by gpt-5, `position-02` by deepseek-flash) and my own
position's claims re-checked rather than trusted. New verification done in this review, all against
committed artefacts, no database started:

- `studies/01-charity-tree/reports/20260912-small.md` — the design labels it actually prints.
- `studies/01-charity-tree/README.md` — the design table (D1–D6 identities).
- `studies/02-ticket-booking/reports/20260913T021206Z.md` — line 12, the cell count and failures.
- `studies/02-ticket-booking/reports/20260921T205212Z.md` — lines 13 and 15, control scope and
  under-booking.
- `studies/05-cache-consistency/reports/analyses/20260921-cache-consistency-allgreen.md` — lines
  69–71, the strict/relaxed table.
- `studies/01-charity-tree/reports/analyses/20260916T090036Z-v3--claude-opus-5--2026-09-16.md` —
  the trial structure claim-02 rests on.
- `book/evidence/claims.json` again for the specific fields contested below.

## Agreements

1. **A1 is confirmed, and it is the strongest defect either position found.** For
   `claim-01-indexes-and-rollups` **all seven** `provenance.cells[].name` values occur **zero**
   times in the cited report, and the same is true of all five names in `claim-04`. I can add
   the mechanism GPT-5 left implicit, because it makes the defect worse than "aliases":

   - the report prints `D1 minimal`, `D2 indexed`, `D3 flat+FK`, `D4 rollup/trg`, `D5 rollup/app`,
     `D6 embedded`; the registry declares `d1-indexed`, `d2-indexed-denormalised-key`,
     `d3-consolidated-aggregate`, `d4-rollup-application`, `d5-rollup-trigger`,
     `d6-embedded-child`, `d7-placement-only-yugabyte`;
   - so the registry's **D1–D3 are shifted by one** (its `d1-indexed` is the corpus's D2, which is
     the only reading under which D1 — which has *no* secondary indexes by construction — could
     be called "indexed"), and its **D4 and D5 are swapped** (`d4-rollup-application` versus the
     corpus's D4 = rollup **trigger**);
   - therefore the declared **winner cell `d3-consolidated-aggregate` names the rollup designs at
     the slot where the corpus's D3 is `flattened-fk`**, and it is listed alongside separate
     `d4`/`d5` rollup entries. It is not a stale alias; it is internally inconsistent with itself.

   Any book link or table generated from this field lands on the wrong mechanism. I withdraw my
   own ordering: this outranks the tag finding in position-02.

2. **A2 is confirmed exactly.** `studies/02-ticket-booking/reports/20260913T021206Z.md` line 12
   reads "Cells: 42 across 3 topologies; **2 failed**", while the registry records
   `failed_cells: 0` for `claim-04`. The reason a green validator coexists with both our findings
   is mechanical and worth stating in the synthesis: the validator's rule 5 checks that the
   declared `failed_cells` matches the registry's **own** `cells` array — i.e. it is a
   self-consistency check that never reads the report. A registry that lists five hand-picked
   cells and declares zero failures passes by construction. Neither position should describe
   `0 errors` as evidence about the source.

3. **A2's control-scope point is confirmed and is a live trap.** `claim-06`'s statement contains
   "not one under-booked seat anywhere", and its source run
   `studies/02-ticket-booking/reports/20260921T205212Z.md` line 13 records "Negative controls:
   **0 of 0** fired" with line 15 "Under-booking: none observed". The run observed no
   under-booking because it exercised no control; the registry loses that qualifier. The book
   must not read that phrase as an audited invariant.

4. **A5, one numeric instance verified.** The registry's limit for `claim-12` says differences
   "range -8.3% to +12.3%"; the source analysis table (lines 69–71) contains `−8.3 %`,
   `+12.3 % (noise)` **and `−14.2 % (noise)`**. The registry understates the analysis's own
   spread, so a book quoting the registry misstates the table it cites.

5. **B1 answerability-first should lead the model additions.** It is the one framing that no
   throughput ranking can express, it has two independent measured instances (Study 02 `r05`,
   Study 03 `r06`) plus one measured remedy, and my position's `b8` was a weaker version of it.
   I defer on ordering.

## Challenges

1. **Two claim sources is the failure mode this repository keeps punishing — do not build it.**
   GPT-5's correction is an *additive signed audit* sitting between `claims.json` and the prose.
   If both artefacts survive, the book has a second source of truth for the same facts, and
   `LESSONS_LEARNED.md` already records where that leads ("An embedded document's field names are
   a second source of truth"). I would rather repair the resolver: make the extractor emit the
   **report's own cell label plus the canonical design id** instead of inventing kebab names, and
   add the two rules the validator lacks — (a) every `cells[].name` must occur in the cited report,
   (b) `failed_cells` must equal the report's own count. The defect in A1/A2 is *generated* by the
   extractor; auditing its output leaves the generator in place and pays for a parallel document
   that can drift.
2. **My own position's headline count was wrong, and I correct it here.** Position-02 says "11 of
   the 13 numeric claims"; the registry has **12** numeric claims, of which **9** (02–10) omit an
   extant tag, with claim-01 untagged and 3 gap records also omitting extant tags. GPT-5's
   "nine numeric claims ... including gap records, 12 entries" is the accurate statement. The
   conclusion is unaffected — a null that means both "no tag exists" and "not filled in" is still
   the defect — but the synthesis must use the corrected figure, not mine. My position is
   immutable, so this review is the correction of record.
3. **Position-01's A4 `trials` finding is confirmed for claim-02, but its direction matters and is
   not stated.** The registry records `trials: 1` while the analysis states "three independent
   trials per cell, each on its own fresh load"
   (`20260916T090036Z-v3--claude-opus-5--2026-09-16.md` line 44). This **understates** evidence,
   the opposite direction from A1/A2/A5, which overstate or misdirect. That distinction should
   drive priority: misdirected claims first, understated trials later.
4. **Apply the "verified placement" standard consistently or not at all.** GPT-5 rightly refuses
   to count Study 03 `l3_section_sharded` as completed colocated/non-colocated coverage. The same
   standard must then be applied to Study 01's `D7`/`D24` ("placement", "colocated"), which the
   study README presents as a controlled pair isolating "physical data placement, on identical
   SQL", and to Study 02's `REPORTS.md` colocation gap. Singling out L3 while the book's brief
   cites D7/D24 as placement evidence would leave the taxonomy inconsistent. I did **not** find
   tablet-level placement verification for D7/D24 either; I state that as an open check, not as a
   finding, because I did not read every Study 01 discussion.
5. **Answerability-first is a decision rule derived from two failures, not a measured ranking.**
   It should lead the *model* chapter; it must not be presented with the authority of a measured
   comparison, because the corpus measures the ledger's cost, not the cost of asking the question
   first. Label it "derived from two measured answerability failures".

## Better option or combination

Take GPT-5's priority order and my mechanism, and merge the overlapping model lists rather than
running both:

- **Evidence layer, in order:** (1) resolve `cells[].name` to the report's own labels and add the
  two missing validator rules — this repairs and *tests* the resolver rather than auditing it;
  (2) only what cannot be derived goes into an attributed qualification layer: confounds, control
  scope, trial structure, and `legacy_provenance_incomplete` for claim-01 — where both positions
  independently converged; (3) the coverage matrix (measured / measured-but-confounded / planned /
  gap / not-applicable) with review state; (4) the two missing gap disclosures GPT-5 found
  (native datastore families, real-network/balanced endpoints).
- **Decide first whether `claims.json` may be regenerated as a new version with the old one
  preserved, or is frozen.** If it may be, the audit is unnecessary. If it is frozen, the audit
  must be the **only** input to the book and the registry must be removed from the prose pipeline.
  Publishing both is the one outcome that should not be chosen.
- **Models:** answerability-before-performance first; then a single merged entry where GPT-5's
  contention-surface model (B4) and my measured arbitration family (b1) are the same subject —
  nine arbiters, two controls, 4–8×; then the derive/index/materialize ladder (their B3) carrying
  my recency case (b4) and my minor models (b5); then the query-boundary portfolio (their B5);
  then the evidence-confidence cards (their B6); then the failure-mode catalogue (their B7), with
  my unit-of-change (b3) folded into their B2 vector rather than duplicated as a family.

## Remaining disagreements and evidence gaps

- **Unresolved by us and not resolvable here:** whether the frozen-registry constraint applies.
  Both positions assume an answer; neither found one in the evidence packet. This is the decision
  the owner or the synthesis must make, and it changes which artefact the book consumes.
- I verified the cell-name failure for 2 of 12 numeric claims (claim-01: 7/7 absent; claim-04: 5/5
  absent) and the −14.2 % row. GPT-5's wider counts — "11 of 12", the claim-05…10 mixed
  operations, the claim-02/06/07 `trials` inconsistencies beyond what I checked — I did **not**
  reproduce claim by claim. Neither estimate should be repeated as fact; the count should be
  *computed*.
- I did not verify GPT-5's A3 multi-run examples (`claim-08`'s reliance on the repeated-race and
  repair runs) or A6's Study 01 endpoint claim beyond reading the same README it cites.
- Neither position read the whole report corpus; both say so. A defect visible only inside a
  report body we did not open remains invisible to both of us.

## Falsification tests

- **A1 (decisive and cheap):** for every claim, count how many `cells[].name` values occur
  verbatim in the cited report. One command settles both our figures and produces the audit's
  first table. If a committed one-to-one mapping table exists (I found none; `claims.json`'s own
  `generated_from` says statements were copied while "strength, cells and limits are curated", so
  the cells look hand-written), A1 reduces to a documentation gap.
- **Second-source-of-truth challenge:** dissolved if the book pipeline is specified to consume
  exactly one claim source and the other is explicitly deprecated and dated. If no such
  specification exists, the challenge stands.
- **Placement consistency:** a committed tablet-level or placement-verification artefact for
  Study 01 `D7`/`D24` (or an explicit statement that none exists) resolves challenge 4 in one
  direction or the other, and the corpus gap should then be stated once for all three studies.
- **`trials` semantics:** if the registry's README defines `trials` as whole-run replications
  only, then claim-02's `trials: 1` is correct by that definition and only the *documentation* of
  the field is defective — which is a much smaller finding than A4 claims.
