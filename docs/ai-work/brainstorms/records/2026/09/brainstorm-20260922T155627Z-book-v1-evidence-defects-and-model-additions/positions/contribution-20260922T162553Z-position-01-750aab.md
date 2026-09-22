# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `contribution-20260922T162553Z-position-01-750aab` |
| Slot | `position-01` |
| Actor | model=gpt-5 tool=openai-codex effort=unknown session=session-20260922T160724Z-a4fe79 capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-22T16:25:53Z` |
| Confidence | high |

## Summary

Require an additive signed claim-qualification audit before numeric book prose uses the current registry. Lead the model additions with answerability-before-performance, then a multidimensional design vector and evidence-confidence cards.

## Disagreements

The current schema-valid registry is not publication-ready as an exact cell resolver, and Study 03 L3 is not verified physical-colocation evidence.

## Contribution

# Independent position — repair the evidence boundary, then teach answerability first

## Position

The strongest evidence correction is an additive, signed **claim-qualification audit** that
sits between the current `book/evidence/claims.json` and book prose. The current registry is
schema-valid but is not yet a reliable claim-to-cell resolver: most cell identifiers do not
exist in the cited reports, several run-failure counts contradict the reports, available run
tags are omitted, and multi-run conclusions are represented as if one run supplied them. The
book must not silently repair the registry in prose or treat validator success as scientific
validation. A new attributed audit should name each defect and either provide exact evidence
anchors or mark the claim ineligible for v1. Existing reports, analyses, manifests, tags, and
the registry remain unchanged.

The strongest model addition is **answerability before performance**: decide which facts must
survive before choosing indexes, layouts, or engines. Study 02 directly measures the cost and
benefit of preserving sale/refund history; Study 03 directly shows that every measured
current-state design erased the hold history needed to size the hold window. This gives the
reader a decision rule more fundamental than a throughput ranking: if the model discards a
business event, no later query optimization can recover the question.

Everything recommended below is book/audit work over committed artifacts. It requires no
build, run, database start, new design, or change to a protected artifact.

## Evidence inspected

- Repository rules and methodology: `AGENTS.md`, `CONTEXT.md`, `LESSONS_LEARNED.md`,
  `docs/methodology.md`, and `docs/environments/host-zenbook-ux5406sa.md`.
- Book contract and queue scope: the book goal and `BACKLOG.md`; the evidence-registry task's
  `BRIEF.md`, submitted `RESULT.md`, HIGH `REVIEW.md`, and current
  `book/evidence/{claims.json,claims.schema.json,README.md}`.
- The cited reports, analyses, manifests, READMEs, handoffs, report protocols, progress logs,
  SQL catalog documentation, and discussions under Studies 01–05. I did not read another
  brainstorm position or draft.
- Git's committed run tags. No database or benchmark was started.

## A. Evidence corrections, strongest first

### A1. Add a signed semantic claim audit; do not publish directly from the current cell map

**Defect.** In 11 of the 12 non-gap claims, at least one `provenance.cells[].name` does not
occur in the cited report. For claims 01–10 every listed cell and every winner is absent.
Claim 11 resolves three scenario ids but invents `postgres-control` and `redis-control`;
only claim 12 resolves all of its names. These are not harmless display aliases:

- Claim 01 calls entries `d1-indexed`, `d2-indexed-denormalised-key`, and
  `d3-consolidated-aggregate`, while the source analysis defines D1 as no indexes, D2 as
  indexed, D3 as indexed plus copied charity key, and D4/D5 as the rollups.
- Claim 04 lists `w1-overbooked-control` and `w2-overbooked-control`; the cited report and
  analysis name the controls C1 and H0.
- Claims 05–10 mix operations or result dimensions (`r05-refund-report`, `ledger-storage`,
  `organiser-edit-latency`) with designs, although the field is called `cells` and the
  README says it lists measured cells.

The evidence-registry result explicitly says these are curator summaries and were not
cross-checked against result JSON; its HIGH review records that as a residual risk. The
task brief nevertheless requires the controlled pair, valid cells, and producing tag for
every numeric claim. Validator success therefore proves schema consistency, not that a
claim resolves to the measurement it names.

**Correction.** Add a new signed audit keyed by `claim_id`. For each atomic claim it records
canonical topology/design/scenario/operation ids exactly as the report emits them, source
status, winner eligibility, every run/digest/tag/commit used, and one of `eligible`,
`eligible_with_qualification`, or `not_eligible_for_v1`. Keep the old registry unchanged and
cite it as the object audited.

**Reader effect.** A link lands on the actual measured cell instead of a plausible-sounding
alias, and “valid winner” becomes checkable rather than self-declared.

**Cost / prohibition.** Artifact inspection and a new Markdown/JSON audit only. No rerun,
new design, database, or edit to an existing report, analysis, manifest, tag, or registry.

### A2. Separate run completeness, claim-support status, and control scope

**Defect.** `failed_cells: 0` hides material source status:

- Claim 04 cites `20260913T021206Z`, whose report says 42 cells and 2 failed YugabyteDB C2
  cells; it also reports C5 under-booking on all three topologies.
- Claim 07 cites `20260921T183722Z`, whose report says 1 of 6 cells failed/was partial and
  whose X1 cluster medians use two rather than three race trials.
- Claim 08 cites `20260915T002411Z`, whose report says 4 of 41 cells failed and lists early
  rejections in designs intended to be correct. Those refusals are the subject of the
  claim, so they cannot be flattened to generic `valid` cells.
- Claim 02's report says 22 cells have errors, missing verification, or invariant failures;
  its analysis carefully distinguishes negative-control firings, deliberately skipped
  incomparable operations, and contention outcomes from initial correctness-gate success.
  The registry loses that distinction.

Control evidence is also run-specific. Study 02's read-only `20260920T234953Z` reports
0 of 6 controls fired because contention was absent; the P3/X1 race runs report 0 of 0;
Study 03's read-only `20260921T231328Z` reports 0 of 0. Claim 06 nevertheless contains the
correctness phrase “not one under-booked seat anywhere.” Its source analysis says the audit's
detection power is inherited, not demonstrated in that run.

**Correction.** The audit must distinguish `run_failed_cells`, exact `claim_support_cells`,
`partial_cells`, `initial_gate`, `post_operation_audit`, and `negative_control_scope` with
values such as `in_run_fired`, `present_not_exercised`, `absent_inherited`, or `none_needed`.
Phrase a clean result as “none observed under this workload” unless the relevant control
fired in the same regime.

**Reader effect.** “Correct” no longer conflates initial answer checks, post-write
invariants, a deliberately failing control, and a complete performance cell.

**Cost / prohibition.** New audit prose/data from the reports and analyses only. It must not
reclassify old cells inside their source artifacts.

### A3. Restore the actual provenance graph, including multi-run support

**Defect.** Git contains run tags for claims 02–10, but all nine numeric claims store
`run_tag: null`. Including gap records, 12 registry entries omit an extant tag. Examples are
`run/01-charity-tree/20260916T090036Z-v3`,
`run/02-ticket-booking/20260913T021206Z`,
`run/03-reserved-seating/20260915T002411Z`, and
`run/04-configuration-portal/20260921T1215Z-small`; their manifests also name those tags.

Some conclusions cannot be represented by one provenance object:

- Claim 01's signed analysis front matter names sixteen run/digest inputs, but the registry
  attributes its compound headline only to `20260912-small`, which has no `repo_commit` or
  run tag in the registry.
- Claim 08's analysis explicitly relies on the matrix, repeated race, and repair runs
  `20260915T173255Z`, `20260915T232736Z`, and `20260916T000706Z`, while the registry names
  only the matrix.
- `claim-gap-04-study02-03-validation-open` is a two-study assertion with one Study 02 run
  and one `study` field; Study 03 has no evidence anchor in the record.

**Correction.** Let the additive audit attach a list of evidence anchors to an atomic
claim. Mark Claim 01 `legacy_provenance_incomplete` and cite all committed digests its
analysis names; do not manufacture a code revision for the legacy runs. Record the full
manifest commit and existing immutable tag wherever present.

**Reader effect.** Tags identify producing code, digests identify data, and analyses identify
interpretation; a multi-run conclusion visibly remains multi-run.

**Cost / prohibition.** Git/report inspection only. Never create a replacement tag, fill a
legacy null by guesswork, or imply that a code tag contains the later result files.

### A4. Atomize compound claims and make `trials` mean one thing

**Defect.** Several headline records bundle independently supported propositions: Claim 01
combines rollups, indexes, contention, and embedding; Claim 04 combines invariant evidence,
race throughput, backoff, and unrelated event edits; Claim 10 combines the c1/c2 arbitration
pair, x1's lost updates, and n3's trigger cost. One status and one winner cannot qualify all
parts honestly.

The `trials` field is inconsistent with its own strength contract. Claim 03 records five
per-cell trials and is `repeated_controlled`; Claim 02 records `trials: 1` although its
analysis states three independent trials per cell. Claims 06 and 07 also record one while
their race evidence uses `-race-trials 3`; their analyses separately warn that a whole cell
was not repeated. Claim 01 compresses a multi-run analysis to one trial. “Trial” currently
means run, cell, or inner race repetition depending on the claim.

**Correction.** Split compounds into one decision/outcome per audited item. Record separate
counts for independent loads/cells, inner operation trials, and whole-run replications; keep
the strength conservative when independence is unclear.

**Reader effect.** A recommendation can be accepted or rejected without inheriting an
unrelated clause, and “repeated” has a stable meaning.

**Cost / prohibition.** Editorial decomposition in the new audit/book evidence layer only;
no stronger label than the sources warrant.

### A5. Publish the confounds beside the numbers, not in a distant caveat

**Defects and required qualifications.**

- Study 01 D20→D22 is a clean read-side location/index comparison but a confounded write
  comparison because D22 maintains D4's whole five-column/two-parent rollup package. The
  recency analysis and `LESSONS_LEARNED.md` say so. Claim 02 must not price “flag versus
  rollup” writes.
- Study 01 D2→D3 measures a package of copied key, SQL, and indexes. The repeated analysis
  explicitly says copying the key alone did not reproduce the gain. Book taxonomy must not
  label the 1.75x score as a pure rolldown effect.
- Study 03's row/document seat-map contrast is a family/package comparison, not a
  one-decision controlled pair; present the scan/point-lookup trade as observed for those
  complete layouts, not as a universal document effect.
- Study 05's 2.2–3.5x cache claim is an **add-cache** outcome. The survey manifest says
  `resource_framing: db-only`, with no cache CPU/memory budget, and the signed analysis names
  the missing equal-total Redis comparison. It is not a resource-efficiency result.
- Study 02's 32-buyer and 128-buyer PostgreSQL P3/X1 arms disagree in sign; the former has
  92–283% P3 spreads and one partial cell, while the latter has tighter spreads. Preserve
  “PostgreSQL unsettled/order-sensitive” rather than averaging or selecting a ratio.
- Claim 12's limit says strict/relaxed differences span −8.3% to +12.3%, but its source table
  also contains −14.2%. The book must use the actual −14.2% to +12.3% observed range and still
  call it single-trial/noise, never “strict is free.”

**Reader effect.** The reader learns which decision a ratio prices and which dimensions moved
with it.

**Cost / prohibition.** Qualification text only; no cross-run pooled ratio and no causal
upgrade.

### A6. Add a v1 methodology-coverage matrix and keep “planned” separate from “measured”

**Defect.** Current coverage is easy to overread as compliance with today's minimums:

- Study 01's equal-total YugabyteDB work uses a single query endpoint; its signed discussion
  says endpoint CPU throttled while peers had spare capacity. It is a deployment result, not
  balanced cluster capacity or isolated replication cost.
- Study 02 `REPORTS.md` records physical colocation as an untouched gap. Its v2 read-only run
  did not exercise controls, and its P3/X1 race runs included no negative control.
- Study 03 `REPORTS.md` calls L3 the colocation axis, but the committed report describes only
  the intended hash-key layout (“one tablet per event, or spread by section”), and the signed
  analysis says tablet balance is unexplained and asks for tablet-level metrics. Without
  explicit placement evidence, L3 is a key-layout/sharding-intent result, not completed
  colocated/non-colocated coverage. Study 03 also has no separately labelled equal-total
  topology condition.
- Studies 04 and 05 explicitly leave YugabyteDB, cluster, placement, equal-total resources,
  scale/regimes, and repeats open; Study 04 implements only 10 of 18 planned designs, and
  Study 05 does not charge Redis to an equal-total budget.

The matrix should show, per study and claim: measured, measured-but-confounded, planned,
unsupported, not applicable with reason, or gap. It should also show review state: Study 04
and 05 independent HIGH reviews and Study 02/03 v2 validation are still open in `CONTEXT.md`.

**Reader effect.** A well-designed protocol is not mistaken for completed evidence, and old
results are not retroactively claimed to meet later rules.

**Cost / prohibition.** A book/audit table from committed plans and results. Do not execute
the already queued review or protocol tasks and do not create a substitute plan.

### A7. Restore gap records promised by the evidence task

**Defect.** The evidence-registry brief requires initial gap records for native datastore
families, real-network topology, and incomplete Studies 04/05 coverage. The registry contains
the Study 04/05 gaps, hold-window answerability, and validation-open status, but no native
datastore-family gap and no real-network/balanced-endpoint gap. The book backlog already owns
future protocols for both, so this is a v1 disclosure omission, not a request to re-plan them.

**Correction.** Add these as entries in the new qualification/gap audit and book coverage
matrix, referencing the existing backlog tasks. Do not create or redesign those tasks here.

**Reader effect.** The book cannot be mistaken for evidence about native cross-major stores,
real networks, host failure, or balanced cluster capacity.

**Cost / prohibition.** Documentation only; no implementation or measurement.

## B. Model additions, strongest first

### B1. Answerability-before-performance: current state versus event history

**Model.** Start with the questions that must remain answerable after correction, refund,
expiry, deletion, and reassignment. Classify each fact as current state, derived current
state, or durable event history. If a required question needs an erased event, choose a
history/ledger model before optimizing reads.

**Direct evidence.** Study 02 run `20260920T234953Z`, digest `160bd80c49bbd892`, report and
signed analysis: r05 is unanswerable in 14 of 15 designs; X1's append-only `sale_event` is the
only measured design that answers the refund questions, at +35% storage and 8x load time.
Runs `20260921T183722Z` and `20260921T205212Z` then price its write/race/refund cost, with the
PostgreSQL disagreement preserved. Study 03 run `20260921T231328Z`, digest
`2f5b619430c2e7a7`: r06 is unanswerable in every design because release clears the evidence;
its `REPORTS.md` names a hold-event table only as an unmeasured remedy.

**Evidence label.** Direct for Study 02's ledger and both studies' answerability failures;
analogy/gap for transferring a ledger remedy to Study 03.

**Reader change.** “Can I answer the business question?” becomes the first branch of the
decision tree, ahead of throughput.

**Cost / prohibition.** One book framework and examples from committed artifacts. Do not
present Study 03's unbuilt hold ledger as measured.

### B2. A multidimensional design vector, not one major-model label

**Model.** Describe a design as orthogonal choices:

1. authoritative state: current rows/documents versus event history;
2. information placement: derive, index, rolldown, rollup, full embedding, bounded embedding;
3. maintenance owner: database trigger, atomic application transaction, asynchronous/non-atomic;
4. arbitration: optimistic condition/CAS, pessimistic lock, serializable, intentionally unsafe;
5. cache: none/local/shared × aside/through × relaxed/strict;
6. physical deployment: node count, replication, endpoint routing, resource framing, verified placement;
7. workload/regime: query boundary, cardinality, contention, mutation mix, freshness contract.

**Direct evidence.** Study 04's n1/n2/n3/n4/c1/c2/x1/x2 controlled variants; Study 05's
scenario id already encodes model/version/backend/strategy/freshness/writers; Study 01's
D3/D4/D5/D6/D9/D10; Study 03's S/E/K/L families. These show that “document,” “normalized,”
or a vendor name alone does not identify the relevant decision.

**Reader change.** The reader can compare one changed coordinate at a time and recognize a
package comparison when several coordinates move.

**Cost / prohibition.** Taxonomy and mapping table only; no claim that every cross-product was
measured.

### B3. Derive → index existing fact → materialize → copy/cache decision ladder

**Model.** Prefer the least stateful mechanism that meets the read requirement. First derive;
then index a fact already maintained; only then materialize/roll down/embed/cache. For every
new copy, name maintenance owner, atomicity boundary, mutation coverage, drift detector, and
write/storage cost.

**Direct evidence.** Study 01's D20 versus D22 recency result (indexing an already maintained
`last_donation_at` beats adding a new flag for the measured reads, with the write confound kept
visible); Study 04's n1/n2/n3/n4/x2 (rolldown and rollup, trigger versus application, drift
control); Study 05's strict cache fences and controls. Study 01 D4/D5 and D9/D10 add the
mutation-coverage and concurrency failure modes.

**Evidence label.** Direct inside each study; analogy when expressed as one general ladder.

**Reader change.** Denormalization stops being a yes/no recommendation and becomes a sequence
of increasingly expensive commitments.

**Cost / prohibition.** Synthesis of existing designs. Do not quote a cross-study numeric
ranking.

### B4. Contention surface and blast-radius model

**Model.** For each write, mark the serialization unit and every unrelated operation that
shares it. Compare optimistic retries with lock waiting at low and high contention; distinguish
logical correctness from throughput and collateral latency.

**Direct evidence.** Study 04 run `20260921T1215Z-small`, digest `ab0f6ef5e5500775`: c2 lock
beats c1 version check under one hot key with 16 writers, while x1 loses 93.7% of acknowledged
updates; single-run/hot-key only. Study 01's parent rollups show contention depends on how many
parents receive writes. Study 02's counter on the event row blocks unrelated edits and its
ledger delays the organiser edit. Study 03's L2 section document moves contention to a document
CAS loop.

**Evidence label.** Direct examples per study; qualitative analogy across domains.

**Reader change.** “One extra row/column” is evaluated by the number of writers and operations
it couples, not by local statement count.

**Cost / prohibition.** Diagram/table from existing evidence; no universal lock-over-CAS rule.

### B5. Query-boundary/locality portfolio

**Model.** A layout has a vector of query effects, not a scalar speed. Map each query by its
natural boundary: point row, parent slice, cross-parent range, aggregate, or history. Optimize
the dominant boundary and show the anti-query that the choice harms.

**Direct evidence.** Study 01 repeated run `20260913T125342Z-v3`, digest
`0bc5900a6d7a5d3a`: D3 beats D2 on selected charity queries because copied key, SQL, and indexes
work together, while donor-history is similar and global sum is slower. Study 03 run
`20260921T231328Z`: section documents accelerate the YugabyteDB seat-map scan while L3 point
lookups fall to 26.4 ops/s. Study 01 full embedding is poor for cross-donor recency.

**Reader change.** “Documents are faster” and “normalization is slower” are replaced with a
query-shape decision and an explicit counterexample.

**Cost / prohibition.** Existing numbers kept within their runs; no cross-study ratio.

### B6. Evidence-confidence and benchmark-reading model

**Model.** Give every book claim a compact evidence card: atomic claim, direct/analogy/gap,
trial structure, noise/spread, gate and same-run control status, complete/partial cells,
environment/resource framing, topology/endpoint/placement status, run tag, digest, commit,
analysis, review state, and falsifier. A result can be fast but invalid, correct but
unlicensed by a control, or reproducible but not comparable.

**Direct evidence.** `docs/methodology.md`; the claim-audit defects above; Study 02's P3 order
effect and missing controls; Study 03's partial/repair runs; Study 05's successive harness
defects and final all-green run; `LESSONS_LEARNED.md` on one-trial noise, coordinated omission,
single query endpoints, and provenance failures.

**Reader change.** The book teaches how to reject a seductive number before it teaches how to
apply one.

**Cost / prohibition.** Book template and existing metadata only; it must not imply statistical
confidence that was never measured.

### B7. Failure-mode catalogue organized by lost information or broken invariant

**Model.** Each entry has: symptom, discarded/duplicated fact, invariant, detector/control,
safe mitigation, and evidence strength. Seed it with lost history (Study 02/03), stale derived
state (Study 01/04), lost update (Study 04), unsafe confirmation/expiry (Study 03), stale cache
publication and multi-instance incoherence (Study 05), and cluster endpoint imbalance (Study 01).

**Evidence label.** Direct within the cited study; analogy when transferred.

**Reader change.** Readers can recognize why a model fails, not merely memorize a winning
design code.

**Cost / prohibition.** Curated book table only. Harness bugs remain methodology examples and
must not be mislabelled as database-design failures.

## Assumptions, alternatives, and dissent

- I assume v1 may add a new attributed audit/appendix without mutating the completed evidence
  task's artifacts. If the book pipeline requires `claims.json` as its only input, then the
  safe alternative is to mark affected claims ineligible and block their numeric prose until
  a separately authorized correction artifact is linked.
- I reject the alternative “the aliases are close enough.” Several aliases change design
  identity or name operations as cells, so they defeat exact provenance.
- I reject treating Study 03 L3 as completed physical-colocation evidence. The committed corpus
  supports a hash-key/layout intent and measured outcomes; it does not contain the explicit
  placement verification required by current methodology.
- I do not recommend a new benchmark, design, or study for v1. Unmeasured items stay gaps and
  the already queued protocols remain untouched.

## Confidence and falsifiers

**Confidence: high** on the registry/provenance/status defects because they are literal
contradictions between committed JSON, reports, manifests, tags, and task records. **High** on
answerability-before-performance because two studies directly record unanswerable operational
questions and one measures an event ledger. **Medium-high** on the remaining analytic models:
they are faithful syntheses but cross-study generalization must stay labelled analogy.

This position would be falsified or narrowed by a committed, pre-existing semantic mapping that
proves the registry aliases resolve one-to-one to actual result cells, by placement artifacts
that verify Study 03's claimed colocated/non-colocated scenarios, or by evidence that the omitted
tags do not identify the cited producing commits. I found none in the declared evidence packet.

## Recommended synthesis order

1. Require the additive claim-qualification audit before numeric book prose.
2. Put answerability-before-performance at the front of the design decision model.
3. Use the multidimensional design vector and derive/index/materialize ladder for taxonomy.
4. Add contention, query-boundary, and evidence-confidence models.
5. Publish all current-methodology and future-study items only as named gaps; create no tasks
   unless the owner later gives the separate `CREATE TASKS FROM BRAINSTORM` command.
