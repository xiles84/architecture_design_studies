# Project context

**What this repository is.** Architecture Design Studies measures what architecture decisions
actually cost. Each study changes one modelling decision, builds designs that differ only in that
decision, loads identical deterministic data, and reports reads, writes, storage and correctness
gates with full provenance. Everything runs in Podman; nothing needs host installation.

**How to read this file.** It is the *current* state only: one status per study and per task, the
rules that are live now, and the open gaps. The full pre-rewrite history — including every earlier
status narrative and receipt — is archived untouched under [`docs/history/`](docs/history/); do not
duplicate that history here. The queue is authoritative for task state; `CONTEXT.md` is a readable
summary, not a second ledger.

## Rules that are live

- **Queue runtime provenance (2026-10-03):** the wrapper rejects stale/unmarked
  queue images against the current source digest and requires an explicit tested
  rebuild. Recovery already persists the new worktree claim; the SELECT worker's
  reported blocker exposed an obsolete runtime image. Repair checks and decision
  log: [`docs/queue-runtime-20261003/`](docs/queue-runtime-20261003/VALIDATION.md).
  Completion milestone: `repo/queue-runtime-repair-20261003` on local main.

- **Never push or pull.** All work is local. The owner controls remotes. Milestone tags are
  annotated, never moved, deleted or reused.
- **One current status per task/study.** A completed task leaves one receipt at its canonical path;
  later conclusions supersede by reference, never by silently editing an earlier signed artifact.
- **Measurements are serial.** The whole machine reports as a single measurement slot
  ([`host-zenbook-ux5406sa`](docs/environments/host-zenbook-ux5406sa.md)); only one run may hold
  `ads-run-lock` at a time. The queue enforces this, and `recover` refuses while the lock is held.
  Non-measuring tasks (docs, reviews, protocols, toolchains) may run in parallel.
- **Correctness gates timing.** A design that answers wrongly reports no timings.
- **AI work is queued.** Goals, tasks, immutable events, reviews, leases and the integration lock
  live in [`docs/ai-work/`](docs/ai-work/WORKFLOW.md); the CLI is [`tools/queue/`](tools/queue/README.md).
  Brainstorming has its own workflow in [`docs/ai-work/brainstorms/`](docs/ai-work/brainstorms/WORKFLOW.md).
- **Capability vocabulary.** Canonical `HIGH`/`LOW`; aliases `leader`/`master`→HIGH and
  `worker`/`follower`/`slave`→LOW are accepted as input only. `primary`/`replica` are datastore
  terms, not capability declarations. See [`AGENTS.md`](AGENTS.md).
- **Work runs on the owner's current machine.** A task that needs a resource this machine does not
  have (a second host, an engine image that is not local) is **not** expected work. Those tasks stay
  `proposed` with a `DEFERRED.md` in their directory and are listed in
  [`docs/optional-todos.md`](docs/optional-todos.md), to revisit only if the resource appears. No
  agent releases a deferred task by itself.

## Current state

### HIGH task sweep — **cache and queue repair integrated; SELECT harness validated (2026-10-03)**

The owner declared HIGH and asked this session to iterate through HIGH work. The
[dated sweep](docs/high-task-sweeps/2026-10-03/REPORT.md) records the live baseline:
29 HIGH tasks completed, 18 unfinished and **zero eligible**, with no awaiting-review
tasks. Twelve proposals have explicit resource deferrals; four feasible Study 09
proposals remain behind the empirical program's T1/T2-before-T5 gate; the cache protocol
comparison was conditional at the snapshot and is now released (Study 05 below);
SELECT synthesis waits for its four measuring workers and
their separately claimed harness. The Edition 3 pilot also awaits the owner's preview
approval before integration; its PDF was requested to open and approval asked in chat.
The audit had zero errors and three existing coordination warnings, retained in its
inputs. The sweep receipt is intended for local-main integration and annotated milestone
`repo/high-task-sweep-20261003`; its queue status records the observed outcome. This
receipt completes an audit, not the scientific backlog.

Execution continued after the owner's renewed request: the conditional cache
comparison is now measured, reviewed and integrated into local main. Seventeen
other HIGH tasks remain ineligible (12 resource deferrals, four T1/T2-gated
hierarchy proposals, SELECT synthesis awaiting measurements). No new eligible
task appeared during the serial matrix. The original sweep remains an immutable
baseline; Study 05 below records the new evidence and observed integration.

The owner's further fix request exposed an obsolete queue image; the tested
image/source provenance repair is completed on local main at
`repo/queue-runtime-repair-20261003`. HIGH then recovered the handed-off SELECT
harness, corrected material defects and validated it on both single-node engines.
Its completion milestone is `study-01/select-v1-harness`; queue receipt records
the local-main integration outcome. There are 32 completed HIGH tasks and 17
gated HIGH tasks after the additional queue repair; SELECT harness itself is a
LOW-minimum prerequisite executed/reviewed in this HIGH session. Its bounded
validation does not complete SELECT evidence synthesis or its four measuring tasks.

### Book initiative and the AI work pipeline — **in progress, next task ready**

- **SELECT templates requested (2026-10-02, HIGH GPT-6 / Codex desktop):** the owner asked
  for alternative SELECT forms, examples and measured regimes where they fit or fail.
  [Eight-task program](docs/select-query-templates/PLAN.md): HIGH protocol/semantic atlas,
  LOW Go oracle/harness, four LOW measurement families (join/CTE, existence/absence,
  aggregates/top-N, access paths/pagination), HIGH evidence synthesis, LOW book integration.
  Scope extends Study 01's existing relational domain and preserves Edition 3's pilot gate;
  priority remains below the existing T1 workers. Task ids start `task-20261002T131729Z-select-`.
  Publication was planning only. Protocol is complete and the corrected harness
  is under bounded engine validation. The publication's local-main integration is recorded
  by [its queue receipt](docs/ai-work/tasks/2026/10/task-20261002T131729Z-select-query-task-publication/STATUS.md)
  and annotated milestone `repo/select-query-task-program-v1`; only this publication is
  complete; measuring families, evidence synthesis and book integration remain queued.
  Queue state remains authoritative.

- **SELECT protocol and harness (2026-10-03):** `select-protocol` is completed. The
  handed-off harness was resumed by HIGH GPT-6 through epoch-4 recovery after
  the queue runtime repair. Ordering, keyset ties/NULLs, unmatched-parent top-N,
  live SQL controls, timing/provenance and plan hierarchy were corrected.
  Final PostgreSQL and YugabyteDB checks each passed 17 templates, 53 live fixture
  cells, four wrong SQL controls and 51 verified client samples. Exact input hashes
  match their generated reports. Both databases stopped and the lock was released.
  HIGH review passed; completion/integration is recorded by the harness queue
  receipt and annotated tag `study-01/select-v1-harness`. The four measuring tasks
  can proceed after queue completion; their matrices were not run by these checks.
  Correction evidence: `studies/01-charity-tree/reports/20261003T-select-high-correction-check.md`.

- **SELECT protocol archive (2026-10-02, HIGH deepseek-flash):** its
  deliverable [`select-query-templates.md`](studies/01-charity-tree/protocols/select-query-templates.md)
  is written — premeasurement predictions P1–P12 with falsifiers, six result contracts, the REQUEST
  seven-form atlas with exact fixture outputs (NULLs, duplicate child matches, unmatched parents, ties,
  empty sets), the sizing/placement and rollup/rolldown/embedding coverage, non-applicability of
  optimistic/pessimistic concurrency, and bounded LOW handoffs for the four measuring families. No
  comparative matrix from this protocol has run. Reviewed (self-review) and merged into local `main` at `b405f88a` under tag
  `study-01/select-v0-protocol`; harness validation is tracked above.

- Goal: [`goal-20260922T025912Z-data-architecture-book`](docs/ai-work/goals/goal-20260922T025912Z-data-architecture-book/GOAL.md)
  — an evidence-backed data-architecture book built with Typst in Podman, with every claim traceable
  to a run, tag, digest and signed analysis.
- Planning checkpoint integrated and tagged `repo/data-architecture-book-handoff-v1`.
- **Edition 3 — concept first (planned 2026-10-01, HIGH claude-opus-5-5):** the owner asked for the book
  to teach concepts first, use the case studies to show when each fits, and keep text short. Rules:
  [`book/STYLE.md`](book/STYLE.md); goal amendment `20261001T155900Z-concept-first-edition-3`. Phase 1 is
  the pilot task `task-20261001T155900Z-book-ed3-pilot` (ready, LOW): template macros, compact claim
  badges, "How to use this book", chapter 4 Rollup and the Study 04 profile; the owner approves the style
  before phase 2 converts the rest. Tag `repo/data-architecture-book-ed3-plan`. **T2–T5 workers:** the
  book's prose is being moved into concept chapters; re-scope a changed number wherever it is quoted,
  including the new case boxes.
- The queue CLI (v1) is **complete and integrated**, tagged `repo/ai-task-queue-v1`, with capability
  aliases tagged `repo/ai-capability-aliases`.
- This document's rewrite is `task-20260922T025912Z-history-status-reconcile`; it imports EH-01,
  EH-02 and Studies 03–05 by reference (see [`docs/ai-work/tasks/`](docs/ai-work/tasks/) and the
  archived history) and fixes public status drift.
- **Active book evidence source:** `book/evidence/v6/claims.json` (tag target
  `repo/book-evidence-registry-v6`; see
  [`book/evidence/README.md`](book/evidence/README.md)). v6 carries all 29 v5 claims forward with their
  statements, limits, confounds, anchors, support cells, trials and lifecycle **byte-identical** and
  adds the empirical-architecture method metadata: `evidence_form` (isolate / package /
  not_applicable — currently 10 / 9 / 10, with **0 interaction**), `attribution` (20 of 29 withhold a
  component-level attribution), `control_coverage` (7 present, 15 absent, 7 not applicable) and a
  `replication` record (28 single-run, 1 run-level, **0 hardware-independent**). Four method lints are
  machine-checked with failing fixtures; `optional` `prediction` and `crossover` fields are defined and
  unused (0 recorded predictions, 0 bounded crossovers). Assignment of every value is auditable in
  `book/evidence/v6/method_migration.json`; the pass is analysed in `book/evidence/v6/ANALYSIS.md`.
  v5 (`book/evidence/v5/`, tag `repo/book-evidence-registry-v5`) re-scoped five claims whose statements
  or limits asserted more than their own evidence supported and remains a **frozen predecessor**. v4
  retired the two stale placement/endpoint gaps and added claim lifecycle metadata (`status`,
  `gap_kind`, `superseded_by`, `closed_dimensions`, `remaining_dimensions`). v3
  (`book/evidence/v3/`, tag `repo/book-evidence-registry-v3`), v2 and v1 are **frozen historical
  reference**; v1's claim-to-cell resolution failed (36 of 43 names resolve nowhere).
  `tools/evidence validate` defaults to v6; `validate-v6/-v5/-v4/-v3/-v2` pin a package explicitly.
- **Claim-versus-limit audit (2026-09-24, task `book-evidence-v5-consistency`):** all 29 active
  claims were audited for the class "the statement asserts more than its limits, trials or
  anchors support". Five were re-scoped and 24 passed unchanged: `v2-10` dropped an
  unsupported "at no measurable cost", `v2-16` dropped a write-path cost its own limit
  disclaims, `v2-07`'s limit no longer cites three runs the claim cannot show, `v2-14`'s limit
  is scoped to its range, and `v3-06` keeps its statement with corrected limits because the
  study's `PRIMARY KEY ((person_id) HASH, ...)` makes the mapping a configuration fact. v5
  adds a repetition-versus-trials lint; `book/evidence/v5/SUPERSESSION_LEDGER.md` and
  `ANALYSIS.md` record every change and the per-claim verdict table.
- **Warning for whoever takes the delta tasks:** re-scoping a claim does not fix prose that
  paraphrases the old wording. `book/concepts/expiry-and-clock-authority.typ` and
  `book/chapters/minor-variants.typ` still repeat phrasing v5 retired, and
  `book/concepts/cache-consistency.typ` carries the `v2-16` sentence. Owned by
  `book-cache-chapter-v3` and `book-review2-delta`.
- **External review of the Edition 1 draft (2026-09-24, GPT-5.6 Sol, high effort):** verified
  against the committed PDF, manifest and registry by the `book-evidence-v4` session. Accepted:
  **R04** (the active registry held two gaps contradicted by newer evidence in the same package) and
  **R03** (one chapter named a v2 registry path while the cover said v3) — both closed by
  `task-20260924T102917Z-book-evidence-v4`. Accepted in part: **R01** (the §15.2 table is physically
  broken) and the figure-numbering/tracking defects, queued as separate tasks. **Rejected: R02's
  severity** — the committed artefact renders real provenance for commit, describe, dirty state,
  clock, image digest, evidence digest and tree hash; the only literal `unknown` is inside Typst's own
  `0.15.1 (unknown commit)` version string. The residual is real and now hardened: an input-less
  compile renders an *Unverified development build* banner, and `build.sh` fails if that banner
  reaches a release PDF.
- **Cache chapter rebuilt (R01, 2026-09-24):** `book/concepts/cache-consistency.typ` no longer holds the
  three-column mechanism table whose collapsed middle column produced five pages of hyphenated
  fragments. The mechanisms are now cards (`mechanism-card` in `book/lib/config.typ`), grouped into
  source correctness and copy correctness, each with Race closed / Mechanism / Guarantee / Does not
  solve / Evidence, followed by a purpose matrix across correctness, freshness, load and recovery.
  The chapter also states the strict-after-acknowledgement contract formally, separates invalidation
  from publication fencing, and names the `wrong_reads` counter as a strict-comparator count. Measured
  effect on the rendered text: orphan-hyphen fragment lines fall from 99 (worst page 33, 17 on the
  page) to 30 (worst page 34, 3), and the chapter shrinks by three pages.
- **External-review remediation queue (2026-09-24):** two tasks are complete and integrated —
  `task-20260924T102917Z-book-evidence-v4` (v4 package, retired gaps, single registry path, build
  gates; tag `repo/book-evidence-registry-v4`) and `task-20260924T105000Z-book-cache-chapter-v2`
  (mechanism cards replacing the collapsed table, plus the cache terminology fixes; tag
  `repo/data-architecture-book-cache-chapter-v2`). Four are **published and ready** for a LOW
  executor: `task-20260924T110000Z-book-text-pass-v2` (R05-R12, R19-R22, R25, R29),
  `task-20260924T110001Z-book-figures-v2` (R26-R28, and it clears the two waivers in
  `book/evidence/check-waivers.txt`), `task-20260924T110002Z-book-pedagogy-v2` (R23-R25 and the
  visual tranches), and `task-20260924T110003Z-book-v2-release` (rebuild `book/dist`, tag
  `repo/data-architecture-book-v2`), which depends on the other three. The tracked PDF in
  `book/dist/` is still the Edition 1 build until that release task runs.
- **Template-token guard (2026-09-24, task `book-template-guard`):** Edition 2's draft rendered raw
  `#registry-label` / `#registry-version` tokens on the title page, the provenance page, the Chapter 18
  heading, the index prose, Chapter 8 and §15.4, because the tokens were written inside backticks (Typst
  raw text is not evaluated) and inside a plain string argument. All six sites now use `#raw(...)` or
  string concatenation. `book/build.sh` fails when the extracted page text contains `#registry-`,
  `#build-`, `#source-`, `${`, `{{` or `UNKNOWN_PLACEHOLDER`, and `book/evidence/check.sh` fails on a
  backticked `#` token or a build-input token inside a string. The review exchange behind this is
  recorded at `book/evidence/reviews/20260924-edition-2-review-response--deepseek-flash.md`.
- **Cache chapter corrected (2026-09-24, task `book-cache-chapter-v3`):** the mechanism cards no longer
  overclaim. The database lock card is split into *pessimistic lock + invariant recheck* (the lock
  serialises; rejecting the loser is the recheck's contribution) and *conditional update /
  compare-and-set* (the expectation goes into the write). The fill lease now says duplicate fills are
  suppressed only while a lease is valid and honoured, the TTL bounds *eligibility to serve* rather than
  refill load, and the outbox is a conditional prerequisite for asynchronously published coherence
  rather than a universal one. The purpose matrix gained a *Change observation* column, so durable
  observation is no longer filed under mutation correctness. "Relaxed" is now presented as a family of
  weakenings with bounded staleness as one member, and the study's measured cells are labelled
  unbounded-relaxed. §15.7's causal "only because" became an observation plus the mechanism those cells
  are consistent with, and the `v2-16` write-path sentence is a labelled mechanism note the claim does
  not quantify. Verified: 51 pages, 29/29 claims, 6/6 fonts, no banner, no unresolved tokens, orphan
  hyphenated fragment lines 29 (worst page 3).
- **Figures re-exported and de-duplicated (2026-09-24, task `book-figures-v2`):** Typst is now the only
  owner of figure numbering — the two book-owned PlantUML sources no longer draw "Figure 1/3" in their
  titles, so the page shows one number per figure (verified: Figures 1-4, captions only). The
  letter-spaced labels are fixed at the export step rather than per figure: PlantUML emits `textLength`
  and `lengthAdjust="spacing"` on every `<text>`, Typst honours them, and **all five** assets carried
  them, so `book/assets/export.sh` now strips both and records `postprocess: strip-textlength-v1` per
  figure in `figure-manifest.json` (schema_version 2). `check.sh` gained a fatal `figure-stretch` rule,
  the two `figure-number` waivers are deleted, and the marker tracking dropped from 0.8pt to 0.3pt.
  Both figure rules were verified by negative control. `check.sh`'s own registry default also moved to
  v5 — it had been left at v4, so a standalone run checked the frozen package.
- **Text and terminology pass (2026-09-24, task `book-text-pass-v2`):** page 3 now defines the evidence
  system the book actually uses — direct, mechanism, analogy, the four `gap_kind` values, partially
  superseded and retired — and scopes "the same progressive structure" to major-family sections.
  Normalization no longer implies it prevents double-selling; application-maintained rollups are
  described as transactional first and "published" only when a copy crosses the authoritative
  boundary; append-only history claims better *answerability*, not better correctness; the topology
  chapter scopes its consensus claim to the measured RF=3 topology and calls the seat-map result a
  layout *package* rather than "one topology change"; `read score` is defined at first use; and the
  trial vocabulary distinguishes an in-run repeat from an independent replication. Numeric ranges in
  the generated cards are typeset with en dashes by `typeset-ranges` in `lib/config.typ`, so prose and
  cards agree while the registry text stays verbatim ASCII.
  Also fixed: `book/assets/export.sh` recorded `source_revision` as the blob at HEAD, which a `--write`
  in a dirty tree stores as a value `--check` can never reproduce; it now records the committed blob
  only when that blob describes the file, else `uncommitted`.
- **Reader-facing framework added (2026-09-24, task `book-pedagogy-v2`):** the Concepts divider is now
  a real part opener ("Part II — Concepts") with a two-sentence lead and a pointer to the glossary, which
  is new (`book/glossary.typ`) and is the single definition site for 26 canonical terms: claim, strength
  values, direct/mechanism/analogy, the four `gap_kind` values, partially superseded and retired, run /
  replication / trial / in-run repeat, `read_score`, `wrong_reads`, the strict contract and the relaxed
  family, resource framing, colocation, rollup/rolldown/embedding, source vs copy correctness,
  invalidation vs fence, confound, and family/variant/scenario. Tranche 1 of the review's visuals is in:
  a book-owned decision-path figure in Chapter 2 (V1) and a freshness-timeline figure in Chapter 15
  (V7), both with a form and evidence-status label and an adjacent text equivalent; the benchmark-reading
  checklist (V11) is a table. The book is 56 pages with seven registered figures, and the chapter-17
  "floor, not an SLO" wording is scoped in the same file.
- **Review-2 delta applied (2026-09-24, task `book-review2-delta`):** repeat citations now render a
  one-line capsule carrying the id, the strength, **every** limit and a pointer to the full card, while
  the first citation and the registry index keep the full card — 34 citations collapsed, and the book is
  54 pages from 56. The counter is Typst state in document order, so it cannot drift from the page the
  way a build-time index could. Partially superseded cards now separate *Moved to successor claims* from
  *Still open in this claim* and label the statement as unchanged from its source registry. The page-3
  legend is generated from the registry's `gap_kind` values in a fixed order, and a new kind with no
  written definition fails the build instead of going unlisted. The placement ladder no longer calls the
  cache "the fastest read" or "the only step that adds a consistency contract"; copy-out is the step that
  crosses the transaction boundary. A lease-versus-fence figure (V8's companion) joins the cache chapter.
- **Edition 2 released (2026-09-24, task `book-v2-release`):** `book/dist/data-architecture-reference.pdf`
  and its manifest are rebuilt from fully merged sources and committed. 54 pages, 29/29 claims indexed,
  6/6 fonts embedded, evidence digest present, no unverified banner, no unresolved template tokens, and
  the manifest's edition now appears on the page (asserted by the build). The tracked PDF is the artefact
  the external reviewer asked for; the previous 49-page Edition 1 build is replaced, and its state is
  recoverable from the tag `repo/data-architecture-book-v1`-era history rather than from this file.
  Two release blockers were found and fixed while producing it: the manifest hard-coded `Edition 1` beside
  a title page reading Edition 2, and the figure manifest's `source_revision` field made staleness depend
  on repository state, so the staleness comparison now uses content hashes only.
- **Edition 2 revised (2026-09-24, tasks `book-visual-correctness-a` + `book-visual-tranche-b`, released
  by `book-v2-1-release`):** the owner-supplied external review of the Edition 2 PDF was applied. Technical
  items A–G: relaxed freshness is a family (only bounded staleness declares Δ); application-maintained
  rollups are transactional-first and "published" only across the authoritative boundary; the placement
  text equivalent no longer calls all five rungs a copy; pessimistic locking is waiting/blocking with
  retries still possible; derived state is separated from retained history and append-only no longer
  implies a fence. The cache chapter's publication-fence card is renamed **source-generation validation**
  and the glossary defines eight distinct mechanisms. Figures: the lease-vs-fence and stale-fill figures
  were redrawn, and nine explanatory figures were added (families ER, rollup vs rolldown, embedding,
  index vs hot row, optimistic vs pessimistic, derived vs history, history vs current, cache read copy,
  topology physical boundary); the book now has 16 figures and is **61 pages**. Item 2.9 (a
  quote-readiness flowchart) was rejected with reason in favour of the existing table; the visual review
  is one model's. No measured number, claim id or gap changed; `book/evidence/v5/` is untouched. The
  change is recorded in `book/CHANGELOG.md` and
  `book/evidence/reviews/20260924-edition-2-revision--deepseek-flash.md`; tags
  `repo/data-architecture-book-visual-correctness-a`, `...-visual-tranche-b`, `repo/data-architecture-book-v2-1`.
- **Book tasks read the active registry only:** `book/lib/evidence.typ` derives its path from the
  single `registry_version` build input; no source may name a version, and
  `book/evidence/check.sh` (run by `book/build.sh`) fails the build if one does. The tracked PDF in
  `book/dist/` is regenerated once, at release, from fully merged sources.
- **Book toolchain:** `book/` holds the modular Typst reader layer, pinned by
  `book/Containerfile` to Typst `0.15.1` (image digest
  `sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f`). `book/build.sh`
  compiles and verifies the draft in Podman and writes `book/dist/build-manifest.json`; the toolchain
  build was 30 pages with 3/3 fonts embedded, all 23 v2 claim ids indexed and the registry digest
  present in the text.
- **Book figure layer:** `book/FIGURES.md` settles the import mechanism (a generated, read-only
  `book/assets/` export chosen by a two-variant pilot) and `book/assets/export.sh` renders each
  canonical study `.puml` with the pinned PlantUML renderer into `book/assets/`, writing
  `book/assets/figure-manifest.json` (source revision + SHA-256, renderer digest, embedded-bytes
  SHA-256). `book/build.sh` runs `export.sh --check` before compiling and fails on a missing asset, a
  changed source or a stale render; `source_tree_hash_sha256` covers the figure layer and the
  referenced sources. The proof set is now authored (`book-figure-proofset`): four figures, each with
  two independent labels (form: structure/sequence/state; evidence status: conceptual illustration /
  implemented design contract / negative control / observed result) and an adjacent text equivalent —
  placement and ownership (structure, book source), two buyers on the same marked seats (sequence,
  study-03 reuse, labelled control), the stale-fill race (sequence, book source, observed result
  illustrating `v2-15` with no embedded rate) and expiry/clock authority (sequence, study-03 reuse,
  labelled control). Book-owned sources live in `book/assets/sources/`; the two-reader check is
  recorded in `book/FIGURES.md` (both passes were one model) and the rule is methodology 15.
- **Diagram fidelity audit (2026-09-23):** all 30 committed study diagram sources were classified
  against their SQL/harness/plans and the audit is recorded per study in
  `studies/0{1,2,3}-*/diagrams/FIDELITY.md`. Three annotations were corrected: D3's "six of twelve
  queries lose a join" (now the current three of sixteen) and its unscoped "never touching the heap"
  (now a named prepared-read result with the 1,756 / 31,293 zero-fetch and 31,293-fetch trial-1
  caveat), and D10's unsupported "per 30 s" corruption window (now per concurrent split, citing
  `20260913-d9-cache-race`). Corrected sources were re-rendered with the pinned helper.
- **Book draft (synthesis):** `book-synthesis-v1` wrote the v1 prose — six family chapters with the
  progressive structure, eight concept chapters, variants, comparisons, topologies and scenarios —
  grounded only in the 23 active v2 claims, with confounds and gaps labelled on the page. It reads
  `book/evidence/v2/` only. No measurement was performed.
- **Book released (v1):** `book-release-review-v1` reviewed the synthesis, re-validated the registry,
  audited every number against the 23 active v2 claims (54 numeric tokens, 0 unresolved) and every
  claim against the prose (23/23 cited), and fixed three reader-visible defects (uncited study
  coverage gaps, a title-page provenance rendering fault, a blank front-matter page). The released
  artefact is `book/dist/data-architecture-reference.pdf`, 36 pages, PDF sha256
  `58d61df0600a46e9b019aa83067475add4e3a2b986345c6041deeea3c8516cb1`, tag
  `repo/data-architecture-book-v1`. It remains a v1: single-run, single-host evidence, and the two
  cited coverage gaps are open work.
- **Work queue state (2026-09-23, end of the book-pipeline execution block):** 38 tasks completed,
  **0 ready, 0 claimed, 0 in progress**; 17 proposed. The six scientific protocol tasks
  (`study04-v2-protocol`, `study05-v2-protocol`, `analytics-read-model-protocol`,
  `hierarchy-model-protocol`, `native-major-model-protocol`, `topology-study-protocol`) are complete,
  and the block they defined has been executed or measured (Studies 01–05 v2 runs, the book evidence
  v3, the figure layer and the cache decision map). The proposed work is the next execution tier —
  Study 06 (`native-*`), Study 07 (`topology-*`), Study 08 (`analytics-*`), the `hierarchy-*` arms and
  `cache-untested-protocols` — none of it ready; publish it with the queue before starting.
- **Brainstorms concluded (3)** — [`docs/ai-work/brainstorms/CONCLUDED.md`](docs/ai-work/brainstorms/CONCLUDED.md)
  is the state authority; each record below links its own `CONCLUSION.md`.
  - `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` (tag
    `repo/book-evidence-and-model-brainstorm-v1`) — two positions (gpt-5, deepseek-flash), two
    cross-reviews, one synthesis. It found registry v1 ineligible as the book's numeric source and
    decided one versioned, attributed evidence-correction package as the sole active claim source,
    with the fixed six-family taxonomy kept. Published work:
    `task-20260922T170845Z-book-evidence-correction-v2` (completed).
    [CONCLUSION](docs/ai-work/brainstorms/records/2026/09/brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions/CONCLUSION.md)
  - `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` — three positions, two
    critiques, one synthesis: diagram sources need one pinned renderer and a provenance gate, the
    committed sources need a fidelity/scope audit, and the book needs a small labelled figure proof
    set rather than diagram dumps. Published work: `task-20260923T235910Z-diagram-renderer-pin`
    (completed), `task-20260923T235913Z-book-figure-provenance` (completed),
    `task-20260923T235916Z-diagram-fidelity-audit` (completed),
    `task-20260923T235922Z-book-figure-proofset` (completed),
    `task-20260923T235931Z-brainstorm-conclusions-index` (this task, claimed).
    [CONCLUSION](docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies/CONCLUSION.md)
  - `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` — three positions, two
    critiques, one synthesis: the cache material must be guarantee-first, and the Study 05 v2 results
    must be registered before the book cites them; unmeasured deletion/outbox protocols stay labelled
    as such. Published work: `task-20260923T235907Z-book-evidence-study05-v3` (completed),
    `task-20260923T235919Z-book-cache-decision-map` (completed),
    `task-20260923T235925Z-cache-untested-protocols` (proposed).
    [CONCLUSION](docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control/CONCLUSION.md)

### Empirical-architecture method layer — **implemented; experiments not started**

- Charter (owner, 2026-09-25): the book's identity is empirical architecture — explain, **predict**,
  isolate, sweep regimes, combine, measure, explain, state the boundary. Crossover/boundary results are
  preferred over larger throughput numbers; controlled single-variable, interaction and package
  comparisons are three distinct evidence forms, and a package may not be attributed to a component
  without an isolation experiment. Replication comes **before** new domains; no reader exercises.
- Implemented in `task repo/empirical-method` (tag `repo/empirical-method-v1`):
  - `docs/methodology.md` sections 16–21 (prediction precondition, three evidence forms, replication
    tiers and ordering, crossovers/boundaries, read-path controls, the eight-question concept
    checklist).
  - `docs/templates/STUDY_PROTOCOL.md` — the per-study protocol template with a committed-before-code
    `## Prediction` section, evidence-form table, read-path controls, sweep levels and replication plan.
  - `book/evidence/v6/` — the method-metadata registry (see the evidence bullet above), with four
    machine-checked lints each proven by a failing fixture.
  - `book/concepts/CHECKLIST.md` — the eight-question reading frame with a per-concept evidence-status
    table, checked against `book/concepts/*.typ` by a new `concept-checklist` rule in
    `book/evidence/check.sh`.
  - `docs/empirical-method-program.md` — the ordered T1–T5 program (replication → crossover sweeps →
    interactions → package isolation → gated domain/topology expansion) with the ordering rule that
    blocks T5 while T1 is outstanding.
- **Blocker corrected 2026-09-25:** the host has no Podman or Typst *in the WSL distro*, but Windows
  Podman 5.8.1 (`/mnt/c/Program Files/RedHat/Podman/podman.exe`, machine `podman-machine-default`,
  already running) is reachable through `infra/lib.sh`, and the pinned `localhost/ads-book:1` Typst
  image is cached. `podman info` reports 8 CPUs / 16,496,418,816 B, matching `host-zenbook-ux5406sa`.
  Measurement and the book build both run here; the earlier "no build can run" note was wrong.
  `.queue-drafts/` is still git-ignored (`.git/info/exclude`), so the program document remains the
  durable record.
- **Verified 2026-09-25:** the book builds against the v6 registry — `book/build.sh` to a scratch path
  gives 61 pages, 6/6 fonts embedded, 29/29 claims indexed, evidence digest present, no unverified
  banner, no unresolved tokens, edition on page.
- **Pagination repair (2026-09-29):** visual review found full evidence cards and labelled callouts
  split across page boundaries, occasionally leaving only a title or label at a page foot. Both now
  render as non-breakable blocks, so Typst moves the complete unit to the next page. A pinned Typst
  scratch compile produced a 76-page review PDF with the affected cards and callouts intact; no claim,
  figure or evidence-package content changed. Integrated into local `main` at `e3e3606` and tagged
  `repo/book-pdf-page-breaks-v1`.
- **Heading rhythm pass (2026-09-29):** visual review found that consecutive chapter subsections
  read as a continuous block because their surrounding gaps were too narrow. The shared Typst heading
  styles now give levels 1, 2 and 3 more room above and below, while preserving the chapter page break.
  A pinned Typst scratch compile produced a 76-page review PDF with the revised hierarchy visible;
  no claim, figure or evidence-package content changed. Integrated into local `main` at `59dead0` and
  tagged `repo/book-heading-spacing-v1`.
- **Found and fixed on this task:** `book/CHANGELOG.md` named a fixed registry path
  (`book/evidence/v5/`), which fails the `registry-literal` rule in `book/evidence/check.sh` on `main`
  itself. Reworded to name the then-active package without a path literal.
- **Done 2026-09-25 (task `book/t1-correction`, tags `repo/book-evidence-registry-v6` and
  `repo/book-t1-correction-v1`):** the eight-question checklist is answered inside all eight concept
  chapters and `check.sh` fails if a chapter lacks it; `v2-13` is re-scoped to the replicated range
  (strength `run_replicated`, a `falsified` prediction, four anchors) and the superseded 6.85x is
  retired; the four prose sites that repeated it are corrected; the tracked PDF is rebuilt on v6
  (65 pages, 29/29 claims, digest present). The new `concept-checklist` rule was negative-controlled:
  deleting the section from one chapter in a scratch copy fails the check with that chapter named,
  while the real tree passes.
- **The book now reads without experimental-design vocabulary (2026-09-26, tags
  `repo/book-plain-language-v1` and `repo/book-cards-plain-language-v1`):** every limitation, in the prose
  and in the evidence cards, states what was measured, then why it limits the conclusion, and only then
  names the term or the id. Eleven book sources were rewritten first, then twenty claims'
  `statement`/`limits` — a *rewording*, declared separately from a *re-scoping* in the carry-forward test,
  which fails if a reworded claim's confounds, trials or anchors move. No evidence id or results field name
  was removed.
- **The three coverage gaps named here on 2026-09-25 are each closed for one claim:** a bounded crossover
  exists (`v2-17`), an `interaction` claim exists (`v2-18`), and run-level replication exists (`v2-13`).
  What remains is breadth, not kind: 29 of 31 claims are still single-run, and the other named interactions
  are unmeasured.
- **The interaction was then repeated three times per size (2026-09-26, tag `study-04/v2-interaction-repeat`):**
  the ordering by size holds in all three repetitions, so the interaction is reproducible in direction, but
  the spread within a size is 15% at 10 entries and 55% at 60, so `v2-18`'s figures are now **bands**
  (~2.1–2.4x, 2.5–4.4x, 13–24x) and it is `run_replicated`. The scatter is a **session-level drift**, not
  noise: absolute throughput fell as the session ran on for both designs (at 60 entries n4 read 631 → 205 →
  206 ops/s). The repeat also dropped the tier-60 median to 3.42x, below the 3.99–5.69x `v2-13` replicated,
  so the corroboration between the two claims is now recorded as approximate rather than a match.
- **The sweeps were then interleaved, and the drift measured (2026-09-26, tag `study-04/v2-interleaved-sweep`):**
  nine runs in a Latin square so size and position in the session stop being the same thing, wrapped in a
  reference cell (`c2` contention at 16 writers) measured immediately before and after. The ordering by
  size held in **all three sweeps**, so the interaction does not depend on the sequential layout — but the
  spread within a size roughly **doubled** (101%, 98%, 104% against 15%, 55%, 52%), because adjacent
  repeats share the session's state and so *understate* the uncertainty. The reference cell put a number on
  the drift: **445 → 640 ops/s, +44%**. `v2-18` then carried the ranges six runs per size actually span
  (~1.3–3.2x, 2.5–9.0x, 10–24x) with `replication runs=6, order=interleaved, spread=104.1`; the two lower
  bands touch when both tranches are pooled, though every individual sweep ordered them correctly.
- **A reference cell on the same operation, read three times through a tranche (2026-09-26, tag
  `study-04/v2-same-operation-reference`):** the ratio proved **more portable than the absolutes** (40%
  spread against 49% and 114%), because a session term shared by both designs cancels in the ratio. But two
  predictions failed and both matter: the session **wanders rather than drifts** (the three reference
  readings went 8.64 → 5.21 → 8.53, so reading only the ends would have reported no movement while the
  middle sat 39% below both), and **a band from six runs did not contain a seventh** — tier 10 measured
  1.06 against a 1.34–3.15 range. `v2-18` now carries 7/9/7 runs per size, bands of 1.06–3.15 / 2.51–8.95 /
  10.20–24.20, and `replication runs=7, spread=134.7`; 23 anchors. The bands describe what has been seen,
  not a range a future reading must fall inside.
- **The first non-study-04 replication (2026-09-26, tag `study-02/v2-06-replication`):** `v2-06-hot-drop-designs`
  — the corpus's headline overbooking claim, one run since 2026-09-13 — was re-measured twice, in opposite
  design order (AB/BA). The **4–8x range reproduced in both orders** (P3 4.30–7.06x and 4.64–6.99x),
  so the claim holds, but **every absolute throughput fell 26–72%** below the original, and
  `p2_precreated_skip_locked` moved **58–76% with its position in the session** while the counters moved
  within 15%. `v2-06` is now `run_replicated` (3 runs, AB/BA, spread 76). Two infrastructure findings came
  with it: study 02's runner passed `hostpath` where the image build needs `winpath` (study 04 already had
  the fix), so **no study 02 claim could be re-run from a WSL path at all**; and the runner **silently
  skipped the `--tag`** on the second run, so that tag was created by hand on the recorded commit.
- **Every build context migrated to one helper (2026-09-26, tag `repo/buildpath-v1`):** the
  `hostpath`-instead-of-`winpath` defect was found in **four more runners** — study 01's four scripts and
  study 03's runner and diagnosis script — so **studies 01, 02 and 03 could not build their benchmark images
  from a WSL path at all**, which is the structural reason none of their claims had ever been replicated.
  `buildpath()` now lives in `infra/lib.sh` with the reason written down, all six remaining build sites use
  it, bind-mount sources keep `hostpath` (correct for them), and study 03's image was rebuilt with it to
  confirm (exit 0). The lesson is recorded too: a correction written inside one study's copy of a shared
  script protects only that copy.
- **The same defect class, second instance: the dirty-tree check (2026-09-26, tag `repo/dirty-check-v1`).**
  Study 04 fixed this on 2026-09-23 and wrote the reason inside its own runner — git's `:!` pathspec
  exclusion does not suppress **untracked** files, so every run after the first report looks dirty and
  `--tag` is skipped *without failing*. **Studies 02 and 03 still carried it**, which is exactly why the
  second study 02 replication produced results and a report but no run tag: the validator caught the missing
  tag, not the runner. Both now filter results and reports explicitly, like studies 04 and 05. Verified with
  a control pair: a probe file under `results/` leaves the new check empty and made the old one non-empty,
  while an uncommitted file under `platform/` is still detected. The two fixes together mean the reason no
  non-study-04 claim had ever been replicated was **two shared-infrastructure defects, each fixed in one
  study's copy**.
- **BLOCKER — sizing defect, found 2026-09-27 (tag `study-03/l3-diagnosis`): the podman machine has 2 GiB of
  memory and every container runs inside it**, while `infra/versions.env` allots the database 3 GB and the
  client 2 GB. Five gigabytes of declared budgets on a two-gigabyte machine: **no container limit can be
  enforced and the VM's own ceiling is the real one**, so the per-node budgets recorded in every manifest
  have never described the environment. This is the mechanism behind the three `l3_section_sharded` load
  failures (server processes **killed with signal 9**, 122 lines in one cell), and it is why raising the
  container limits to 5 GB changed the symptom instead of removing it. The machine was also found **stopped**
  (state `stopped`, hence "lock volume vanished" and `connect: connection refused`); it has been started, and
  the configuration is unchanged. **Decision needed from the owner — resize the machine
  (`podman machine set --memory 8192`) or shrink every study's declared budgets to fit 2 GiB.** Either is an
  environment change: runs before and after are not comparable, and the change must be recorded as such.
  Consequences: study 07's 3-node topology cannot be attempted here at all, and every published run's
  declared envelope carries this gap.
  **Scope refined 2026-10-02 (t3):** the blocker is a **multi-node** limit. A single-node YugabyteDB arm
  starts, runs and completes cleanly here — `20261002T150337Z-feas` on `yb-single` did, which is what made the
  `v2-10` replication possible — so `yb-cluster3` and study 07's three-node topology remain unmeasurable while
  `yb-single` does not. Every published YB run before this date still carries the declared-envelope gap.
- **Study 03's replication history, corrected (passes `study-03/v2-11-replication`, `…-order-reversed`;
  both superseded by `study-03/v2-11-memory-envelope` and then `study-03/close-out`).** An earlier version of
  this bullet said neither pass registered, which was true when written and is not true now: `v2-11` is
  `run_replicated` with **five whole runs** and both engines. The first pass confirmed the direction
  (1,138 ops/s against 41, a **27.8x** multiple); the order-reversed pass lost two cells to memory pressure
  (`killed with signal 9`, 122 lines in one cell); the memory-envelope pass supplied the AB/BA pair
  (**27.8x and 24.4x**, a 12% difference) and registered the claim; the PostgreSQL half gave **0.83x**; and
  the close-out pass reversed the PostgreSQL pair for **0.88x**. The engine contrast — about 25x on
  YugabyteDB against essentially none on PostgreSQL — is therefore the claim's most replicated element.
  One correction travelled with this work: an analysis of the first pass had recorded the wrong
  `inputs_digest` (the first 16-hex string in the report rather than the labelled one).
- **`v2-11` registered after a third pass, and the vocabulary gap closed itself (2026-09-27, tag
  `study-03/v2-11-memory-envelope`).** The retry raised the envelope (database 3 GB → 5 GB, tserver limit
  1.5 GiB → 3 GiB) and ran the pair in the **opposite** order, which is what the lint had been demanding:
  the 100k-seat multiple came out **24.4x against pass 1's 27.8x, a 12% difference**, so `v2-11` is now
  `run_replicated` with `order: abba` and three anchors. The raised envelope fixed `l1`'s load (which failed
  in pass 2) but **not `l3`**, which failed a third time — in three different positions, under two
  envelopes, and now with a different symptom (`connection refused` instead of SQLSTATE 57P01). The sharded
  half of the claim therefore cannot be tested on this host with this harness, which is not the same as
  having tested it and found nothing. `record_host_state()` captured the budgets and host memory beside the
  results on its first use.
- **Remaining for the method layer:** publishing T1–T5 as queue tasks (`tools/queue` works now), closing
  the breadth gaps above, and the standing follow-ups — **16 of 31 claims are still `single_run_directional`**
  (7 `gap`, 1 `mechanism_supported`, 1 `repeated_controlled`, **6 `run_replicated`**, after `v2-08` was
  registered on 2026-09-27 and `v2-09` on 2026-09-28), the other
  named interactions are unmeasured, `l3_section_sharded`'s load needs a diagnosis, and the
  reference cell should now be taken **from the same operation being compared**, between every run rather
  than only at the ends of a tranche, so the drift becomes a per-run covariate.

### Study 01 — charity tree — **complete; recency question analysed**

`charity → person → donation`, 8 schema designs × 3 PostgreSQL/YugabyteDB topologies, correctness-gated.
v1, v2 and the v3 independent review are integrated and tagged (`study-01/…`, `repo/recency-reports-integrated`);
the recency question (EH-02) is measured and analysed. Repository minimums for future studies are
tagged `repo/study-comparison-minimums`.
→ [Study](studies/01-charity-tree/) · [reports](studies/01-charity-tree/reports/)

### Study 02 — avoiding overbooking — **complete and HIGH-validated**

`band → event → ticket`, 14 strategies including two deliberately wrong controls. The original
sell-out race and the v2 operational-reports phase (3b) are measured, analysed and — as of
2026-09-22 — **HIGH-validated** (`task-20260922T025912Z-study02-03-v2-validation`; signed
analyses `20260921T205212Z--deepseek-flash--2026-09-22`, `20260921T231328Z--deepseek-flash--2026-09-22`).
The validation accepted the ledger's ~2x YugabyteDB race cost at 128 buyers and the refund
answerability result, left the PostgreSQL race ratio *unsettled by design order*, and required
three artifact-only corrections (widened `claim-07` range, dropped `claim-05` "8x load", retired
`claim-gap-04`); verdicts in `book/evidence/reviews/20260922-phase-3b-high-validation.json`.

**The 100-seat reversal of the `v2-08` replication was settled (2026-09-27, six runs `20260927T22*`, tag
`study-02/v2-08-tier100-settle`).** The earlier tranche's 0.99x was a run median taken from a cell whose own
trials disagree by 50–92%. Six more runs (5 race trials per cell, three in each design order, tier-1,000 as a
positive control) read `p3/x1` at **1.279, 1.333, 1.257, 1.712, 1.328 and 1.476** at 100 seats — every run
inside the claim's `1.1–2.3x` band — and the ledger was slower in **29 of 30 individual trial pairs** (pooled
per-trial median 1.391x). The tier-1,000 control read 1.304–1.712 in all six runs and correctness was green
(46/46, 48/48, 0 bad counters). So the band survives and **"tight spreads" does not**: the 100-seat cell spread
47–77% and the ratio moved 34% across the six runs (54% across all eight). `v2-08` then recorded 8 whole runs
with the two earlier replication passes as anchors; **re-scoping its signed wording is still a HIGH decision.** A
first-trial warm-up appeared in 6 of 6 runs. Analysis
`studies/02-ticket-booking/reports/analyses/20260927T-v2-08-tier100-settle--deepseek-flash--2026-09-27.md`; the
earlier replication analysis's two reports still say "None yet" because it lacks the `inputs:` frontmatter line
the report index needs (recorded in `LESSONS_LEARNED.md`) — not edited here, as it is another session's file.
The tranche is merged into local `main` and tagged `study-02/v2-08-tier100-settle`; `CONTEXT.md`,
`LESSONS_LEARNED.md` and the v6 registry were updated in the same change. Two registry points a later
reader should know: the record's `across_run_spread_pct` was corrected from a within-cell figure (57.0)
to the run-to-run spread of the tier-100 ratio (54.1), and the two earlier replication passes were added
as anchors, which the earlier replication analysis had promised but not recorded. The tracked book draft was
rebuilt on the new registry in the same change (`book/dist/`, 72 pages, 31/31 claims, evidence digest
`06f3cc01…`), as every recent claim-registry change has done.

**The ledger's cost is monotone in tier (2026-09-27, four runs `20260927T224424Z`/`…225107Z`/`…225817Z`/
`…230535Z`, tag `study-02/v2-08-tier-axis`).** Four whole runs measured **all five tiers in every run**, in
both design orders, 5 trials per cell. The pooled run-median `p3/x1` rose **1.075, 1.237, 1.409, 1.845, 2.505**
at 10 / 100 / 1k / 10k / 100k seats — monotone in the pooled medians and in every individual run — so the
claim's `1.1–2.3x` band describes its **middle three tiers**, not the axis. Both endpoint tiers had never been
repeated: at 10 seats the run medians are 1.023–1.180 (pooled 1.075, below the floor) but the per-trial median
is 1.131 and the ledger is slower in only **12 of 20** pairs, while at 100,000 seats it is slower in **20 of
20** (pooled per-trial median 2.456, above the 2.3x top — but that tier is **one event** per trial, so an
event-specific effect is indistinguishable from a tier effect). Correctness was green (0 bad counters in all
four runs). `v2-08` now records **12 whole runs, 13 anchors and 15 limits**, and its tier-100 ratio spans
0.992–1.712 over the twelve. Analysis
`studies/02-ticket-booking/reports/analyses/20260927T-v2-08-tier-axis--deepseek-flash--2026-09-27.md`. **The
axis is now the main input to the HIGH re-scoping: the band is an interior summary and "tight spreads" fails
at the small tiers.** Merged into local `main` and tagged `study-02/v2-08-tier-axis`; the tracked book draft
was rebuilt on the new registry.

**`v2-08` re-scoped by HIGH (2026-09-27, task `repo/v2-08-rescope`).** HIGH validated both worker tranches
(their numbers were recomputed from the raw results and reproduce exactly) and then took the decision they
had flagged. The signed statement no longer carries the flat `1.1–2.3x` band or the words "tight spreads":
it names the tier axis (about 1.1x at 10 seats to about 2.5x at 100,000 seats on PostgreSQL), labels the
YugabyteDB, refund and organiser-edit figures as single-run and not re-measured, and keeps
`run_replicated` (12 whole runs). `book/evidence/v6/SUPERSESSION_LEDGER.md` gained a **post-tag section**
recording this decision and the earlier `v2-06` and `v2-11` re-scopings, which the test's `changed` set had
listed while the ledger still said "one exception"; `ANALYSIS.md` gained the matching update. The three book
prose sites that quoted the removed band (`chapters/scenarios.typ`, `chapters/major-families.typ`,
`concepts/derived-state-and-history.typ`) were corrected in the same change and the tracked draft rebuilt.
No other claim's statement or evidence changed.

**`v2-09` re-scoped by its own sweep (2026-09-28, eleven runs `20260928T01*`, tag
`study-02/v2-09-buyer-sweep`).** Leader and follower in one session: HIGH pre-registered a buyer-count sweep
with order controls, LOW executed it. `v2-09`'s swap rested on a primary comparison that moved buyer count
(32 → 128) **and** design order together — its own limit said neither could be isolated. Crossing them at
1,000 seats found the **ledger slower in 55 of 55 trial pairs**, both orders agreeing on the sign at every
level (8/16/32/64/128 buyers), the magnitude **non-monotone** (pooled medians 3.08, 1.32, 1.23, 1.59, 1.57)
with **no crossing of ratio 1**; the 32-buyer cell the primary read as 0.585 now reads 1.233/1.219, so the
swap was the confound (and a five-day-old harness state). Two pre-registered predictions — a monotone rise
and a bounded crossing — were **falsified**; the reference cell moved 8.6%; one 8-buyer cell changed state
mid-run and is reported, not discarded. `v2-09` is now `run_replicated` (11 runs, alternating order) with
B2 recorded as a falsified prediction and the ratio's non-portability as a limit. Analysis
`studies/02-ticket-booking/reports/analyses/20260928T-v2-09-buyer-sweep--deepseek-flash--2026-09-28.md`; the
v6 ledger's post-tag section now records four re-scopings, and the two book prose sites that called the
ledger cost "order-sensitive" were corrected. Merged into local `main` and tagged; the tracked draft was
rebuilt on the new registry.

**The sweep's tier weakness was then closed (2026-09-29, two runs `20260929T00*`, tag
`study-02/v2-09-lowbuyer-tiers`).** Two more whole runs at 32 buyers across all five tiers show the same
stable direction at every tier: all ten tier × order cells above 1 (1.021–1.606) and the 32-buyer ratio below
the 128-buyer one everywhere, so the primary's **tier-wide reversal does not reproduce anywhere** on the
axis. T1–T5 all confirmed. The tranche added two anchors and two limits to `v2-09` and changed no statement;
it also found the small-tier floor — at 10 and 100 seats only 3/10 and 6/10 trial pairs favour CAS, with
spreads of 86–174% — while the 1,000-seat cell reproduced the sweep's ratio a day later within 1%. Analysis
`studies/02-ticket-booking/reports/analyses/20260929T-v2-09-lowbuyer-tiers--deepseek-flash--2026-09-29.md`.

**The 8-buyer collapse was explained, not a ledger effect (2026-09-29, six runs `20260929T01*`, tag
`study-02/v2-09-8buyer`).** The sweep's 8-buyer cell had fallen 5x mid-cell and was kept unexplained. Six
alternating runs (twelve cells, sixty trials) found **four** low-throughput events (619–655 ops/s against a
session normal of ~2,500) falling in **both designs** — three of the four in the CAS arm — and both orders,
falsifying the pre-registered arm-specificity and order predictions. Three events were a single trial that
recovered; the original three-trial run looked persistent only because its slow trials were the cell's last
three and so set its median. So the low mode belongs to the 8-buyer level's very short trials, **not to the
ledger**, and no 8-buyer magnitude is quoted. `v2-09` gained six anchors and three limits with no statement
change; the `v2-09` limit that had called it the ledger arm's is corrected. Cause still open — the artifact
carries the database's cgroup counters but not the client's, so the fix is instrumentation, not more cells.
Analysis `studies/02-ticket-booking/reports/analyses/20260929T-v2-09-8buyer-recurrence--deepseek-flash--2026-09-29.md`
— **superseded** the same day by the client-quota audit below.

**The client's own quota, read at last (2026-09-29, no new measurement, tag
`repo/client-quota-audit`).** The harness has written `client_cpu` since the study's first commit — one entry
per **trial** for the race phase — and no analysis had read it. Reading it **refuted the previous
paragraph's own inference**: in all twelve 8-buyer cells the client's throttling is **zero** and the slow
trials use *more* client CPU (2.6–2.8 s against ~1.9 s), so the client waits rather than starves, and with
the database's ~4 ms the 8-buyer low mode is a latency excursion in the response path, not a quota effect.
The audit also found **35 client-throttled cells out of 775** across study 02, all in the **read phase** of
two anchors (`20260920T234953Z`, `20260913T021206Z`) and throttled by a **design-dependent** amount within one
run (2.3 s to 58.1 s), so those reports' read figures are not quotable; `v2-06`, `v2-07` and `v2-gap-01` each
gained a limit saying so, and no statement, strength or anchor changed for any claim. `v2-09`'s limit about
the 8-buyer cause was corrected. The superseded inference is archived under `reports/outdated/`. Also fixed
here: root `CONTEXT.md`/`LESSONS_LEARNED.md` are now pinned to LF in `.gitattributes`, because a CRLF working
copy made both files read as modified with zero content difference and would have blocked this merge.

**Read phase re-run: the floor is the client (2026-09-29, two runs `20260929T16/17*`, tag
`study-02/read-family-rerun`).** The audit implied the read-phase throttling was that session's. Four designs
(`p3`, `x1`, `c1`, `r1`), both orders, same phases and scale, say otherwise: **`r1_inventory_row` loses 420 of
634 and 422 of 633 client CFS periods today** (`p3` 2.3–5.9 s, `x1` 2.6–4.6 s, `c1` 5.4 s), and its
`q01_event_availability` reads are **flat at ~38k ops/s across all five tiers** (10 → 100,000 seats) where
`p3` falls 25,967 → 595 — a flat cell over a 10,000x input range is the client's ceiling, not the query's
cost. So the read family's numbers are **client-bound** for client-heavy designs at the study's two-core
client budget: the fix is **client sizing**, not another re-run. P1 falsified; P2–P5 confirmed — the designs'
own ordering reproduces (`q05` C1 12.2–12.7k against P3/X1 ≈ 2.0k; the `q03`/`r03` gap 143–256x; `r05`
answered only by the ledger at 30.3k/33.0k against 32.2k), and X1's "~2x slower across the board" read penalty
does **not** reproduce (0.77–1.05, 0.97–1.05 against 0.29–0.73), confirming the primary analysis's own
suspicion. No claim statement, strength, confound or anchor changed; the three read-hazard limits were
extended. Analysis
`studies/02-ticket-booking/reports/analyses/20260929T-read-family-rerun--deepseek-flash--2026-09-29.md`
— **superseded** the next day by the budget calibration below.

**Client budget calibration: the quota was a symptom, not the ceiling (2026-09-30, three runs
`20260930T00*`, tag `study-02/read-client-sizing`).** `infra/lib.sh` could vary the database's budget from the
environment but not the client's, so the question had never been asked. Three budgets (client 1 / 2 / 3 cores
against a fixed two-core database on the **4-CPU** machine): at **three cores the client is throttled 0 times
in 633 periods** and the reads are **slower** (52,884 → 43,726 ops/s at 10 seats), because the core comes out
of the database's share (its throttled periods rise **238 → 847 → 1052**). S1/S2/S4/S5 confirmed; **S3
falsified** — `r1`'s `q01` is still flat when the client is unthrottled (43.7k → 41.6k), so my earlier "the
floor is the client" reading was wrong: a flat cell is a ceiling *or* a genuinely O(1) design. **Resolved the
same day by a plan check** (`20260930T-read-q01-plan-check`, no new measurement): `r1`'s `q01` is
`Index Scan using event_inventory_pkey … rows=1`, cost **2.49**, **0.012 ms**, 3 buffers, against a 690-row
table that no tier changes — **O(1) by construction**, while `p3`'s `COUNT(*)` over tickets touches 58,927
rows for cost 4,517 and 6.774 ms. So the flatness is the design's point and the *level* is the machine's
point-lookup rate. Consequence: on this machine the read family's **levels** cannot be repaired
by re-allocating cores (the budgets already sum to 4 CPUs, and the declared memories 3g+2g are 2.5× the 2 GiB
present); the actionable options are a **smaller client read footprint** (fewer, longer-lived connections) or
a **bigger machine** — both owner decisions. The family's **orderings** remain quotable. `CONTEXT` and
`LESSONS` corrected; the three read-hazard limits were rewritten to the corrected attribution.
Analysis
`studies/02-ticket-booking/reports/analyses/20260930T-read-client-sizing--deepseek-flash--2026-09-30.md`.

**Owner decision, 2026-09-30: the machine is fixed, and the read phase gets a smaller client footprint.**
The project assumes the computer it has and will not spend on more resources, so neither open question about
bigger hardware is open any more: `DB_CPUS=2 + CLIENT_CPUS=2` on the 4-CPU / 2 GiB podman machine is the
environment, and the read family's levels are to be repaired by sizing the generator down. Calibrated the
same day at three read-worker counts (tag `study-02/read-client-footprint`, three runs plus a verification
run): **client throttling 422/633 → 33/632 → 0/631** at 8/4/2 workers, the database's own throttling
**900 → 560 → 109** periods, and at 2 workers every read cell is latency-bound (`ops/s × p50 = 1.9` for 2
workers, for both designs at every tier) — the diagnostic that says the client is no longer the constraint.
`r1`'s lead over `p3` on `q01` survives every footprint, and `p3`'s 100,000-seat cell is 621–631 ops/s at all
three, i.e. genuinely store-bound. **Implemented**: the harness gained `-read-conns` and `run-study.sh` now
defaults `READ_CONNS=2` for the read phase only, leaving the write, race and hold phases at the published
`-conns 8` so their published numbers stay comparable; every manifest records `client_read_conns`, and read
levels are comparable only within one footprint. F1–F5 all confirmed. **Owed next: the read family's
re-measurement at the new footprint** (the four designs of the 2026-09-29 re-run first, then the remaining
eleven). Analysis
`studies/02-ticket-booking/reports/analyses/20260930T-read-client-footprint--deepseek-flash--2026-09-30.md`.

**The read family's first footprint-2 re-measurement (2026-09-30, two runs `20260930T16*`, tag
`study-02/read-family-footprint2`).** The four designs of the 2026-09-29 re-run, both orders, at the new
default: **0 throttled periods in all eight read phases** and every cell latency-bound (`ops/s × p50 = 1.9`
for 2 workers), so G1/G2/G5 confirmed and G4 confirmed (`r1` at 17–18k against 38.2k at footprint 8). Two
findings to carry: the **store-bound cell is footprint-invariant** (`p3` at 100,000 seats: 623 at footprint
8, 680 / 629 at footprint 2), and **the design margin is not portable** — `r1` over `p3` on `q01` was 2.07x
at footprint 8 but 1.16x / 1.48x here, because the old margin was partly the client's unequal throttling
(5.3 operations in flight against 3.2) rather than query cost. The three read-hazard limits are now
**footprint-explicit** (four designs at footprint 2; the eleven not re-measured, including both primary
anchors, still not quotable); no claim statement, strength, confound or anchor changed. Analysis
`studies/02-ticket-booking/reports/analyses/20260930T-read-family-footprint2--deepseek-flash--2026-09-30.md`.

**The read family is complete at footprint 2 (2026-09-30, two runs `20260930T20*`, tag
`study-02/read-family-11`).** The remaining eleven designs, both orders: **0 throttled periods in all 22 read
phases**, so all fifteen designs now share one non-client-bound footing (client read CPU 40.8–68.0 s against
a 126 s allowance), and the orders agree on `q01` to within 0.91–1.03. Two refinements: the in-flight
diagnostic distinguishes store-bound cells (1.8–2.0, the point-lookup and single-row queries) from
**client-decode-bound** ones (0.7–1.5, `q02`/`q05`, the two 28-row queries) — correcting the previous
comment that every cell was latency-bound — and the primary anchor's per-design read inflation is now
quantified at **1.09x to 2.22x, in rank order with the periods its client lost per design**. The three
read-hazard limits are complete (fifteen designs at footprint 2; the primary anchors' own figures still not
quotable). Analysis
`studies/02-ticket-booking/reports/analyses/20260930T-read-family-11-footprint2--deepseek-flash--2026-09-30.md`.

**One read worker, and a correction (2026-09-30, two runs `20260930T23*`, tag `study-02/read-one-worker`).**
The tranche meant to separate store-bound cells from client-bound ones falsified the separation instead: with
the **mean** latency, in-flight is **1.000 in all 72 cells** at one worker and **2.00** at two, so every read
cell in the family is a store measurement at a stated footprint. The previous day's claim that `q02`/`q05`
were client-bound was wrong twice — `q05` returns one row, not 28, and the low counts came from using the
**median** in Little's law on a skewed distribution (`mean/p50` up to 2.7x). **The registry caveat is
withdrawn** (`q02`/`q05` levels are store numbers, quotable with the footprint attached) and `LESSONS` carries
the corrected entry plus "Little's law wants the mean". Read-phase client CPU 25.6–41.9 s, 0 throttled periods.
Analysis
`studies/02-ticket-booking/reports/analyses/20260930T-read-one-worker--deepseek-flash--2026-09-30.md`.

**The primary anchors' read cells re-stated (2026-10-01, tag `study-02/read-restate`).** No new measurement:
an integration analysis over the six footprint-2 runs. All fifteen designs now have a **`pg-single`
footprint-2** read footing, so the anchors' footprint-8 read figures (inflated 1.09x–2.22x per design, and
their ratios along with them) are **replaceable rather than caveated**; `v2-07` gains the two `x1` runs as
corroborating anchors and its limit states the replacement rule. Two gaps stand: no single run covers all
fifteen designs (cross-set orderings carry a session term), and both anchors' **`yb-single`/`yb-cluster3`
read arms remain uncovered** — a resource gap on the fixed 2 GiB machine. The book needed no change (checked:
its read figures come from Study 01, Study 05 and `v2-11`, not this family). Analysis
`studies/02-ticket-booking/reports/analyses/20261001T-read-cells-restated--deepseek-flash--2026-10-01.md`.

**All fifteen read designs in one session (2026-10-01, runs `20261001T1009*`/`20261001T1039*`, tag
`study-02/read-all15`).** The joint run closes the "complete per design, not jointly" gap: **0 throttled
periods in 30 read phases**, all **270 cells store-bound at 2.00** operations in flight (by the mean), no
≥1.5x pair swaps between the two orders (agreement 0.92–1.40). Two substantive readings: **`r1` and `p4` are
tied within one session (1.5% / 7.7%)** — the 5–10% gap seen across sessions was an artefact, superseding the
cross-set ordering in the previous entry — and **the top five designs (`h0`, `h1`, `r1`, `r2`, `p4`) differ by
4%, inside the session's own order term**, so the family should be quoted as **groups** separated by 1.5x, not
as a ranking. **`c1_count_naive` carries a 40% position term** (first of order A against last of order B), so
its `q01` level is not quotable from these runs. `q05`'s structure separation reproduces within one session at
3.4–3.9x. The three read-hazard limits now carry the joint run and the quoting rules. Analysis
`studies/02-ticket-booking/reports/analyses/20261001T-read-all15--deepseek-flash--2026-10-01.md`.
All six predictions were confirmed; the first/last position warning above remains applicable.

**Study 01 T1 replication (2026-10-01, runs `20261001T1125*`/`20261001T1210*`, tag
`study-01/t1-read-replication`).** The three oldest claims — `v2-01`, `v2-02`, `v2-03`, all
`single_run_directional` from one 2026-09-12 run — were re-run as **two whole runs in reversed design order**
at the primary's own settings (`pg-single`, scale small, 10 s per op, the same eight designs). Every ordering
reproduced (`d3 > d2 > d6 > d1`; the rollup family far above `d2`) and every ratio held within 1.4x between the
runs, while absolute read scores moved 40-45% — so the claims were **re-scoped to bands** (`d2`/`d1` 10-11x
against the quoted 13.3x; `d3`/`d1` 14-16x against 20.5x; `d4`/`d5` 61-85x against 73.8x/69.9x; `d6` 1.7-2.2x
against 2.4x) and promoted to **`run_replicated`** with the two runs as anchors; the book sentence quoting them
was re-scoped too. Two predictions were falsified, both mine rather than the claims' (`K4`'s `d6` threshold
contradicted the primary's own 0.57x). **Study 01's runner records no `repo_commit`/`repo_dirty`**, so the runs
carry a manually created tag at commit `9356884` — a named harness weakness. The YB arms stay unreplicated
(2 GiB machine). Analysis
`studies/01-charity-tree/reports/analyses/20261001T-t1-read-replication--deepseek-flash--2026-10-01.md`.

**T1–T5 published (2026-10-01, HIGH).** A LOW session cannot publish — `queue publish` refuses any session that
is not HIGH — so the five T1 replication tasks were written by a HIGH session and published **ready**, one per
study whose claims are still `single_run_directional`:
`task-20261001T133759Z-t1-study01-replication` (v2-04, v2-05) · `…-t2-study02-replication` (v2-07 and a survey
of what the fixed machine can carry) · `…-t3-study03-replication` (v2-10, with a mandated feasibility check and
a coverage gap instead of a partial run if it does not fit) · `…-t4-study04-replication` (v2-12, v2-17, which
also carries the owner's optimistic/pessimistic requirement) · `…-t5-study05-replication` (the cache, churn,
strict/relaxed and colocation claims, which also carries the colocated/non-colocated requirement). Each is
`minimum_capability: LOW`, `work_role: executor`, and carries the same acceptance shape: a protocol with
predictions committed before any run, ≥2 whole runs in reversed order, honest scoring, promotion or a limit
naming what did not reproduce, book prose re-scoped, and a merged, tagged result.

**T1 completed, and the Edition 3 pilot reviewed (2026-10-01, HIGH).**
`task-20261001T133759Z-t1-study01-replication` is **completed**: its two claims are `run_replicated` on main at
`88a9c04` (tag `study-01/t1-remaining`), with the task's required tag `study-01-t1-remaining-v1` on the merge
commit. The HIGH review re-derived both claims from the committed results with its own script and they match to
the digit (1.59/1.72/1.64/1.84; group 1.44/1.43; q13-trailing 1.92/2.33); the reviewer's *first* script
disagreed because it matched query names by prefix and swept in the `q08@largest`/`q12@smallest` regime
variants — the worker was right and the reviewer was wrong, so the run's own `cells.tsv` is the authority for
what a score contains. `queue integrate` could not run because the canonical branch had been renamed mid-task
to stop a branch/tag name clash; the branch was restored at the merged commit.
`…-t2-study02-replication` was initially claimed and then **yielded unstarted** for the higher-priority
Edition 3 pilot; it later completed with HIGH review as recorded below. `t3`–`t5` remain `ready`.
`task-20261001T155900Z-book-ed3-pilot` is **approved_for_integration**, with integration deliberately deferred
to the owner's approval of the preview style as its brief requires. The review rendered the pages and found five
defects no source check could see (headings numbered `0.x`, a stray comma in the verdict tag, the front page
overflowing to a second page, the checklist split orphaning a row, the `CHAPTER 4` eyebrow on the previous
page); all five are fixed on `book/ed3-pilot` at `6b6dd6d`. The preview is 90 pages (front p. 3, chapter 4
pp. 8–9, Study 04 profile p. 13), 31/31 claims, 6/6 fonts, no banner, and stays untracked in the task worktree.

**T2 — Study 02's `v2-07` replicated, its storage figure re-scoped (2026-10-02, LOW deepseek-flash).**
`task-20261001T133759Z-t2-study02-replication` re-derived `v2-07` in **four whole runs of all fifteen designs**
on `pg-single` — two with `verify,read` (both design orders, the evidence) and two `verify`-only (a pre-check).
The answerability half reproduced exactly in both orders: `r05` unanswerable in 14 of 15 designs with
`x1_cas_ledger` the only design answering it, `r01`/`r03` partial in the same 14, `r02`/`r04` answerable in all
15; 15/15 clean verify cells per run, 0 over-booked, 0 duplicate seats, 0 ledger mismatches. The storage half
replicated as a **band whose value depends on the rule**: tables only +33.4%/+34.8% (primary 35.2%), all
relations +35.0%/+35.9% (primary 36.2%), indexes alone +37.7%. The claim is now `run_replicated` with its
statement naming the rule, `book/chapters/how-to-choose.typ`'s "+33–35% storage" re-scoped to "+33–38% storage
— about 35% by table size, 38% counting indexes too", and its `M4` (reproducible within 1%) recorded as
**falsified**: the delta moved 0.9 percentage points between two runs 29 minutes apart, because relation size
carries the load's free space. Analysis
`studies/02-ticket-booking/reports/analyses/20261002T-t1-answerability-replication--deepseek-flash--2026-10-02.md`.
The task's survey ranks what remains measurable on the fixed 2 GiB machine: `v2-08`'s PostgreSQL refund and
organiser-edit arms first, then `v2-09`'s buyer×tier surface; the `yb-cluster3` and `v2-11` YugabyteDB arms
**do not fit** and stay recorded as coverage gaps.

**T2 completed and reviewed (2026-10-02, HIGH).** The HIGH review re-derived every number independently: `r05`
unanswerable in 14 of 15 designs with `x1_cas_ledger` the only design answering it, `r01`/`r03` partial in the
same 14, tables-only storage +33.4%/+34.8% (primary 35.2%), all relations +35.0%/+35.9% (primary 36.2%), and
the index-only figure +37.4% to +37.9% across the five runs. That last one was quoted as a point (37.7%) in the
submitted statement, which undercut the claim's own band thesis; the review tightened it to the band. `v2-07`
is `run_replicated` on `main` at `fed5a06` (merge state `35877d8`), tags `study-02-t1-remaining-v1` and
`study-02/v2-07-answerability-replication` unchanged on it. **Two registry observations for the next
revision:** two `run_replicated` claims (`v2-04`, `v2-07`) now carry `verdict: falsified` because the corpus
records the *prediction's* verdict and not the *claim's*, so a reader scanning strengths against verdicts sees
a contradiction the registry has no field to explain; and the queue has no way to reap the live slot ref of a
completed task, so `queue audit` keeps reporting `expired_lease` and `live_ref_drift` for
`task-20261001T133759Z-t2-study02-replication` after its completion — a tooling gap, not project state.

**Study 01's last two claims replicated, and a runner defect (2026-10-01, runs `20261001T1444*`/`20261001T1459*`,
tag `study-01/t1-remaining`).** `v2-04` and `v2-05` are now `run_replicated`. `v2-05`'s `d3`-over-`d2` read
advantage held in **all four fresh loads across two whole runs** (1.59x, 1.72x, 1.64x, 1.84x) against the
primary's 2.47x — the direction robust, the magnitude down about 30%. `v2-04` replicated **as a query-specific
result**: the rollup index beats the flag on q13 (1.9–2.3x) and q15 trailing (2.6–3.1x), ties it on q14 and is
marginally *slower* on q16, so its statement was re-scoped from the 2.9x headline to the query dependence and
the ~1.4x group ratio (the primary's own data gives 1.52x for the group), with the 7.9x YugabyteDB figure left
unreplicated. **The tranche also found and fixed a runner defect**: `run-enhancements.sh` reads its cell plan
from a file and, on this host, every podman invocation consumes the shell's stdin — so it ran exactly one cell
per run and still reported "failed processes 0". Proven with `podman --version` alone; the plan is now read on
fd 3. The published matrices are unaffected (their manifests count 100 and 192 cells), and two one-cell
diagnostic runs from before the fix are kept in `results/` as the evidence.
→ [Study](studies/02-ticket-booking/) · [reports](studies/02-ticket-booking/reports/) · [REPORTS.md](studies/02-ticket-booking/REPORTS.md)

### Study 03 — reserved seating — **complete and HIGH-validated**

`venue → event → seat`: 40-minute holds, seat-map reads, all-or-nothing multi-seat blocks, 14 designs.
The v1 matrix and both repair runs are measured, reported, analysed and signed; the v2
operational-reports runs (3b) are measured, analysed and HIGH-validated (2026-09-22): r06 is
unanswerable in every design, L2 answers the 100k-seat map 4.5-11x the row designs on YugabyteDB,
and L3's point lookups collapse to 26 ops/s.
→ [Study](studies/03-reserved-seating/) · [reports](studies/03-reserved-seating/reports/)

### Study 04 — configuration portal — **v1 measured and reviewed; v2 tranches measured; T1 replication landed**

`product → installed product → configuration entry`: 18 designs planned, 11 implemented, measured on
PostgreSQL (`small`) with two negative controls. The independent HIGH review landed 2026-09-22
(`task-20260922T025912Z-study04-independent-review`, signed analysis
`20260921T1215Z-small--deepseek-flash--2026-09-22`): it accepted the 93.7% control loss, the
15.2-15.8x rollup read gain, and the deterministic 10.5x document storage win; it narrowed the 1.76x
lock advantage to an upper bound and found that the ~1.7x identical-SQL attribution floor is real but
its proposed cell-order cause is not corroborated by the recorded order. The v2 tranches (controls,
cardinality/cadence, staleness, `d2`) are measured — see open gap 3 below.

**T1 replication (2026-09-25, tag `study-04/v2-t1-replication`):** the first run-level replication in
the corpus. Three whole runs (`20260925T190831Z`, `20260925T191610Z`, `20260925T194447Z`; order seeds
101/202/303) of `n3_rollup_trigger` vs `n4_rollup_app` vs the `x2` drift control on `w05_replace_all`,
`pg-single`, `small`, three trials each. The prediction was committed **before** the runs
(`T1-REPLICATION.md`). Result: the direction replicated in 3/3 runs (application-maintained faster,
3.99x / 4.94x / 5.69x, median 4.94x) but the previously cited **6.85x did not reproduce**; the trigger
cell's own trials disagree by 35–115%, above the 20% precision rule, so the pair now supports a
direction and a range, not a point ratio. New: the trigger's p99 is 61–109 ms against 2.2–2.7 ms
(22–46x), and the trigger's **first trial was the slowest in every run** (1.5–5.6x below the cell
median) — a within-cell position effect. The `x2` control fired in all three runs (60 rollup
mismatches, INV-8/INV-9). The read-phase prediction was then tested in three read-only runs (tag
`study-04/v2-t1-read`): `n4` read `r01_effective_config` **1.52x** faster in all three, but the
two-design shuffle produced the **same order for all three seeds**, so `n4` ran first every time and the
effect cannot be attributed to the design — P2 stays **unresolved (confounded)**, and the tranche ran
without `explain`, so no plan explains it. **Follow-up resolved it (tag
`study-04/v2-position-control-crossover`):** run in both orders, the gap follows the *order* — with `n3`
first `n4/n3` = 0.62, with `n4` first = 1.50; each design reads `r01` about **1.55x** faster when it goes
first, and the plan, SQL, storage and row counts are identical in all four combinations. So P2 is **not
falsified**: there is no attributable read difference. Signed analyses:
`studies/04-configuration-portal/reports/analyses/20260925T-read-replication--deepseek-flash--2026-09-25.md`
(read tranche), `.../20260925T-position-control--deepseek-flash--2026-09-25.md` (position control),
`.../20260925T1-replication--deepseek-flash--2026-09-25.md` (write tranche).
`v2-13` in the active registry now states the replicated range, and the four book sites that repeated
the 6.85x were corrected in the same task (tag `repo/book-t1-correction-v1`).

**The corpus's first bounded crossover (2026-09-25, tag `study-04/v2-position-control-crossover`):**
`v2-17-optimistic-pessimistic-crossover` — the writer sweep was extended from four levels to seven
(1, 2, 4, 8, 16, 32, 64) rather than relaxing the `crossover` lint, which requires five levels and a
stated interval before a reversal may be presented. `c2/c1` is 0.89 at one writer and 1.10 at two, so the
crossing is bracketed by adjacent levels; it peaks at 1.68x and does not reverse again. The registry now
has **30 claims**, one with a `crossover` record; `crossover-bounded` moved gap → measured.

**The mechanism was then tested and rejected (2026-09-26, tag `study-04/v2-retry-control`):** pre-registered
in `studies/04-configuration-portal/RETRY-CONTROL.md`, then measured by setting the harness retry budget to
zero at 16 writers. `c1` did not collapse (657 → 725 ops/s) and the ratio did not narrow (`c2/c1` 1.60
without retries against 1.61 with them, same day) — and since `c2` never retries at any setting and moved
by the same 10%, the arm-to-arm difference is run-level, not retry work. **`v2-17`'s mechanism sentence is
withdrawn**: the crossover stands, the explanation does not. The control also measured a ~30% same-settings
**day term** (`c2` 869 ops/s on 2026-09-23 against 1,056 on 2026-09-26), now a limit on `v2-17`: absolute
throughput is day-specific and only within-run ratios carry across days. The control also **confirmed the
counter invariant under rejected publications** (`c1` ended at 4,948 against 4,948 expected, `c2` at 7,917
against 7,917, both 0 lost updates) — recorded because a first version of the analysis wrongly reported
that invariant as missing: it lives in `final_counter_value` / `expected_counter_value` / `lost_updates`
in every result file, and the report's contention table renders it, so querying the non-existent keys
`final_value` / `expected_value` returned `null` and was misread as an absence in the artifact.

**The corpus's first interaction claim (2026-09-26, tag `study-04/v2-interaction`):**
`v2-18-maintenance-size-interaction` crosses rollup maintenance strategy with entries per installation
and finds the penalty is a function of size rather than a constant — the application design is **2.07x**
faster at 10 entries, **4.38x** at 60 and **24.20x** at 240, with the trigger's p99 rising 3 → 156 ms
against the application's 1 → 8 ms and the trigger's throughput falling 46x across the axis. Correctness
held everywhere (0 audit failures, 0 rollup mismatches over 600 / 3,600 / 14,400 entries). `v2-13` gained
a limit saying its 4.0–5.7x holds at the 60-entry default; the tier-60 cell of this crossing reads 4.38x,
inside that replicated range. Registry is now **31 claims**, `evidence_form` counts 11 isolate / 9 package
/ 1 interaction / 10 not_applicable, and `interaction-experiments` moved gap → measured.

**A phase that does not apply is silently absent (harness defect, 2026-09-26):** the same experiment was
first designed to cross maintenance strategy with concurrent writers, which cannot run —
`harness/main.go` guards the contention phase with `if !d.hasContention() { return nil }`, so for the
rollup designs it returns with no result and no warning (three runs produced `concurrency: null`). A
skipped phase and a phase that measured nothing were indistinguishable in the artifacts. **Fixed
2026-09-26 (tag `study-04/harness-phase-status`, `main @ 3c16e39`):** every result now carries
`phases[]` with `status: ran | not_applicable` and a reason, and the report renders a *Phases that did not
apply* table. Verified in one run with the control included: `c1_optimistic_version` records
`contention: ran` and is not listed, `n3_rollup_trigger` records
`contention: not_applicable (design has no arbitration to contend over)` and does appear.

**T1 concurrency replication: `v2-12` and `v2-17` are `run_replicated` (2026-10-02, LOW deepseek-flash, tag
`study-04/t1-replication`).** `task-20261001T133759Z-t4-study04-replication` re-ran the
optimistic-versus-pessimistic pair on `pg-single` at 1, 2, 16 and 64 writers, both design orders, with
zero-retry controls at 16 writers. The lock's 16-writer advantage replicated in all four fresh runs
(1.64x/1.57x at the default budget against 1.69x/1.56x at zero) and the retry budget moved the ratio by
only 1.1%–3.2% in a matched order, so **`v2-12`'s "upper bound" hedge is retired** by the control its own
text said it could not run; its ratio is now a band of five whole runs (1.56x–1.76x). The negative control
reproduced the claim's figure to a tenth (93.7%/93.9% lost at 16 writers; none in either correct design at
any level). The **crossover between one and two writers reproduces** (0.905x/0.939x/0.948x at 1 writer;
1.035x/1.106x/1.166x at 2; 1.56x–1.69x at 16; 1.73x–1.96x at 64) and is not order-dependent. Two
corrections worth carrying: the tranche's first eight runs had all executed the pair in the **same** design
order, so the pre-registered "both orders" control was re-run and the order read from each manifest before
this analysis was trusted; and **one 1-writer run (3.85x) is an unexplained excursion** (`c1` read 306
ops/s there against 1,286 at the same level and order) rather than the order effect the first draft
reported. The runner's gate **refuses a negative control that cannot fire**, so `x1` at 1 writer is a
failed, not-applicable cell; the analysis reports it. 14 runs tagged; analysis
`studies/04-configuration-portal/reports/analyses/20261002T-t1-concurrency-replication--deepseek-flash--2026-10-02.md`;
registry and supersession ledger updated; book concurrency prose re-scoped. The 32-, 8- and 4-writer levels
remain primary-only.

Next: the v2 completion protocol's remaining items — recover the seven unbuilt designs (`d3`, `d4`,
`h1`, `s1`, `s2`, `y1`, `y2`), equal-total resources, verified placement, and the three controls the
review requires (repeated-design instrument control, `-retries 1`, writer sweep).
→ [Study](studies/04-configuration-portal/) · [reports](studies/04-configuration-portal/reports/)

### Study 05 — external cache — **v1 reviewed; v2 runs landed (churn/scale, resources/engines, open-loop)**

**V5 protocol/deletion comparison measured (2026-10-03, HIGH GPT-6 / Codex
desktop):** `task-20260923T235925Z-cache-untested-protocols` follows the renewed
HIGH execution request. Three restarted whole runs passed **30 cells, 216 trials
and 173 control gates**, with zero candidate wrong/impossible reads. Source
statements/warmed read were 1 for version validation and 0 for synchronous/outbox.
A first/last ratio of 2.2446 and client throttling prevent a general throughput
ranking. Durable identity and retained fencing resisted deletion/recreation;
expired fencing allowed resurrection in the controls. The benchmark lock is released.

Source `fed61266`; run tags `run/study-05/cache-v5-20261003-r1`, `-r2`, `-r3`;
digests `f2054e7244c3d2b9`, `52c1eb9b8c355d19`, `76c170e2f1a3e27e`.
Premeasurement tag: `study-05/v5-protocol-20261003`; control tag:
`study-05/v5-controls-20261003`. Required final integration milestone:
`study-05/v5-untested-protocols`, verified as an annotated tag at integrated
main commit `e71af320ad213676f54db5587ad776ca1000e561`; the queue receipt records
subsequent completion metadata.
See the [signed analysis](studies/05-cache-consistency/reports/analyses/cache-v5-20261003--gpt-6--2026-10-03.md),
[whole-run reduction](studies/05-cache-consistency/reports/cache-v5-20261003-whole-runs.md),
[validation](studies/05-cache-consistency/protocolbench/VALIDATION-04.md), and
[environment](docs/environments/xiles-podman-4cpu-2g-config-protocol-20261003.md).
The keyed projection leaves cluster, multi-host placement, uncoordinated writers
and saturation as gaps. Existing donor-portal evidence/book claims are unchanged.
HIGH self-review accepted the result and the queue integrated the owned branch
into local main. The worktree is `.worktrees/study05-v5-untested-protocols`;
no remote synchronization was performed.

Cache-aside/write-through, strict vs relaxed freshness, wrongness of a relaxed cache and cache-side
invalidation fences, on PostgreSQL + Redis (Redis measured as a shared cache, never as an
authoritative store). The single-`small` all-green run is measured and analysed, and its limits are
explicit (single trial, unmeasured TTL/churn, no YugabyteDB/cluster/placement, no equal-total
accounting, open-loop demand absent) — the v2 runs have since addressed churn, equal-total and
engines/placement, and open-loop demand (see item 4). The independent HIGH review landed 2026-09-22
(`task-20260922T025912Z-study05-independent-review`, signed analysis
`20260921-cache-consistency-allgreen--deepseek-flash--2026-09-22`): it accepted `claim-11`,
`claim-12` and `claim-gap-02`, and found one material defect — the all-green analysis swaps the
legacy and owned three-instance wrong-read rows (87.81% belongs to the owned model, 87.24% to the
legacy model). Verdicts in `book/evidence/reviews/20260922-study05-independent-review.json`. Next:
the v2 protocol (`task-20260922T025912Z-study05-v2-protocol`, now unblocked) with the review's four
added requirements (repeated trials/second machine, real-TTL churn, equal-total framing, executed
placement pair).
→ [Study](studies/05-cache-consistency/) · [reports](studies/05-cache-consistency/reports/)

**T1 replication of the cache, churn and engine claims (2026-10-02, HIGH deepseek-flash, tag
`study-05-t1-replication-v1`).** `task-20261001T133759Z-t5-study05-replication` ran eight whole runs in
reversed pair order and re-tested five of the eight single-run targets. **Promoted to `run_replicated`:**
`v2-16` (strict/relaxed read difference 7.3%/5.7%, sign flipping, inside the 8.7%-22.8% spread), `v3-02`
(both churn strategies crossed 2.0 real TTL boundaries with 0 wrong reads in both orders), `v3-03` (both
resource framings recorded and kept apart, 8.2k against 7.9k), `v3-04` (yb-single strict/relaxed 1.043 and
1.105, inside spread). **Re-scoped:** `v2-14`'s throughput band reproduced (2.24x/2.54x) but its **p99
clause did not** (27 ms → 9-10 ms, about 3x, not an order of magnitude), so the claim now carries the band
and records the tail as not reproduced under this client. **Two findings to carry:** `v3-02`'s aside cell
recorded **1 impossible cache value** in 1.17M churn reads while its gate passed — the study's floor says
impossible values must be zero, so the claim keeps a limit and the **gate's coverage of impossible values
during churn is a new gap**; and the book's "p99 falls an order of magnitude" prose was corrected in the
same change. **Deferred/gaps:** `v3-01`'s 1800 s horizon was not re-run; `v3-05` and `v3-06` need
`yb-cluster3`, which the fixed 2 GiB machine cannot carry, and no tablet/leader placement record exists, so
the required colocated/non-colocated comparison remains an explicit coverage gap. Eight runs tagged; analysis
`studies/05-cache-consistency/reports/analyses/20261002T-t1-replication--deepseek-flash--2026-10-02.md`.
Merged into local `main` at `7c53e24`; the task is complete and the worktree removed.

### Infrastructure and environment

- Environments: [`docs/environments/`](docs/environments/) — one page per machine, with the ways it
  can mislead a benchmark. All current results come from `host-zenbook-ux5406sa`.
- Podman orchestration: [`infra/`](infra/) — PostgreSQL 1-node, YugabyteDB 1-node and 3-node RF=3,
  Redis, and the pinned queue image `localhost/ads-queue:1`.
- Diagram renderer: `infra/diagram-render.sh` renders every study's `diagrams/*.puml` with
  `PLANTUML_IMAGE` from [`infra/versions.env`](infra/versions.env), pinned by digest
  (`sha256:9b9ee6af…`, PlantUML 1.2026.8); the three study `diagrams/render.sh` scripts delegate to
  it and each study's `diagrams/RENDERER.md` records the digest actually used.
- Replication: [`docs/replication.md`](docs/replication.md). Methodology: [`docs/methodology.md`](docs/methodology.md).

## Open gaps (current, not historical)

1. **Study 02/03 v2 validation** — **closed 2026-09-22.** The phase-3b runs passed every
   AM-03/AM-04 acceptance check; the signed HIGH validation is
   `task-20260922T025912Z-study02-03-v2-validation` (analyses `20260921T205212Z--deepseek-flash--2026-09-22`
   and `20260921T231328Z--deepseek-flash--2026-09-22`). Residual, non-blocking gaps named by it:
   no control fired in the measured runs, the PostgreSQL ledger race ratio is unsettled by design
   order, and the ledger's load multiplier needs its own repeated load-only cell.
2. **Evidence registry** — **closed 2026-09-22.** The active source is `book/evidence/v2/`
   (tag target `repo/book-evidence-registry-v2`); v1 (`repo/book-evidence-registry-v1`) is frozen
   historical reference. `tools/evidence validate` reports `23 claims, 0 errors` on v2, and
   `v1-diagnostic` keeps the frozen audit (11/12 claims, 36/43 names) as a regression fixture.
   Correction ledger and synthesis handoff are in `book/evidence/v2/`. The review verdicts the
   correction honoured are in `book/evidence/reviews/`.
3. **Study 04 v2** — protocol complete; **controls and cardinality/cadence run 2026-09-23** (tags
   `study-04/v2-controls`, `study-04/v2-cardinality-cadence`). `HANDOFF-AMENDMENT-01-V2.md` (tag
   `study-04/v2-handoff`) defines all the work. Measured so far: identical SQL reproduces within
   0.9%–6.7% across fresh-load phases; the lock's advantage at 16 writers is 1.19x against a
   retries-disabled optimistic design and 1.57x against default; the writer sweep gives
   c2/c1 = 0.89 / 1.26 / 1.58 / 1.60 at 1/4/16/64 writers with 0 lost updates; across cardinality
   tiers 1–500 the row designs collapse 62x on `r04_list_installations` (11.8k → 188 ops/s) while
   the rollups stay flat (13.9k–14.7k, ~74x apart at 500), and the trigger's `w05_replace_all` is
   10–11 ops/s against the application rollup's 413–448 across quiet/steady/burst cadences. Two
   named gaps remain: the **within-run identical-SQL position effect** (the runner cannot place one
   design first and last in a phase) and — until 2026-09-23 — **derived-state staleness/rollup lag**,
   which the staleness task then measured (tag `study-04/v2-staleness-field`): n3 (trigger) p50
   4.30 ms / max 6.88, n4 (application) p50 4.01 / max 4.34, 12 probes each, 0 timeouts, so both
   atomic designs have **no window above the probe's ~5 ms resolution**; their difference is write
   cost, not freshness. Analysis `20260923T1305Z-staleness--deepseek-flash--2026-09-23`. Signed analyses
   `20260923T00116Z-controls--deepseek-flash--2026-09-23`,
   `20260923T0030Z-cardinality-cadence--deepseek-flash--2026-09-23`. The third task
   (`…-study04-v2-remaining-designs`, tag `study-04/v2-remaining-designs`) built and measured **one**
   of the eight unbuilt designs: `d2_doc_on_parent` — its predicted W6 metadata amplification did
   **not** appear at small cardinality (1651 vs 1673 ops/s, publication 866 vs 859), an honest
   negative result. The other seven (`d3`, `d4`, `h1`, `s1`, `s2`, `y1`, `y2`) are recorded as
   coverage gaps with the harness change each needs, alongside the equal-total arm and placement
   verification. Verdicts in `book/evidence/reviews/20260922-study04-independent-review.json`.
4. **Study 05 v2** — protocol complete; **real-TTL churn run 2026-09-23** (tag
   `study-05/v2-churn-scale`). `HANDOFF-AMENDMENT-01-V2.md` (tag `study-05/v4-handoff`) defines
   real-TTL churn, medium scale and a tested ceiling, repeated randomized trials, equal-total cache
   resource accounting, YugabyteDB/cluster/verified placement and open-loop demand. The churn task
   measured 600 s write-through, 600 s write-aside and 1800 s write-through cells that crossed
   **2.00 and 6.00 real 300 s TTL boundaries** with **0 wrong reads** in 1.37M–2.88M reads each.
   Two caveats are recorded: the **hard-expiry path fired 0 times** (only probabilistic early expiry
   did, 11.3k–37.1k), and refreshes are not recorded. The 5×-for-120 s burst (no write-rate knob)
   and the 8k/80k/800k donor scales (harness maximum is 3,000 people) are gaps. The resources/engines
   task then ran (tag `study-05/v2-resources-engines`): the **equal-total arm is labelled** (warm
   10,837 vs 8,027 db-only; mixed arms within ~1.03x), the **strict/relaxed pair ran on yb-single
   RF=1 and yb-cluster3 RF=3** with all three endpoints spread (strict 514 vs relaxed 596 app-99 on
   one node; strict 942 vs relaxed 846 on three), and the **colocated/non-colocated pair** favours
   colocated ~1.3–1.5x. A blocking harness defect was fixed: a comma-separated DSN list reached one
   pgx pool, failing every cluster cell; `pgxdb.OpenSpread` now spreads operations over endpoints and
   `connection_nodes` is recorded. The **tablet/leader distribution is still not evidenced** (the
   placement file holds only the server list), and one host is not colocation evidence. Signed
   analyses `20260923T1100Z-churn--deepseek-flash--2026-09-23`,
   `20260923T1215Z-resources-engines--deepseek-flash--2026-09-23`. The **open-loop demand phase then
   ran 2026-09-23** (tag `study-05/v2-open-loop`, run `20260923T0255Z-openloop`, digest
   `50d5a426c89e96cf`): the harness gained an `openloop` phase on the platform's `measure.RunOpenLoop`,
   and the canonical cell offered **5,000 ops/s → 4,969 delivered (99.4%, 0 dropped)** and
   **30,000 ops/s → 15,599 delivered with 141,710 of 300,004 arrivals (47.2%) never attempted**. The
   client's bounded buffer, not the server, was the first limit (closed-loop floor 21,175 ops/s in
   the same cell), so 15,599/s is a client-limited lower bound; 0 wrong reads in both regimes.
   Analysis `20260923T0255Z-openloop--deepseek-flash--2026-09-23`. No published task remains; verdicts
   in `book/evidence/reviews/20260922-study05-independent-review.json`.
5. **Book v1** — **released 2026-09-22** (tag `repo/data-architecture-book-v1`): the 36-page PDF at
   `book/dist/data-architecture-reference.pdf` (sha256 `58d61df0…`), synthesised from the 23 active
   v2 claims and independently reviewed. Remaining book work is future editions and the coverage
   gaps the protocols below will close; no v1 acceptance check is open. The sources were then made
   **root-independent** (tag `repo/book-preview-paths-v1`, commit `89d8a0d`): `lib/evidence.typ`
   used a root-relative path that broke editor preview when the language server rooted at the
   repository instead of `book/`. That rebuilt artefact (sha256 `379b31f9…`, 36 pages, 4/4 fonts,
   23/23 claims) was then the current build. `book-figure-proofset` later rebuilt from the v3 registry
   and the figure proof set: `book/dist/data-architecture-reference.pdf` sha256 `d6d3c029…`, 49 pages,
   5/5 fonts embedded, 29/29 v3 claims indexed, evidence digest `814f4094…` present, built from a clean
   tree — the current **draft** build, not a reviewed release; the released v1 PDF remains at its tag.
6. **Study 07 topology** — protocol complete 2026-09-22: `studies/07-topology/HANDOFF.md` (tag
   `repo/topology-study-handoff-v1`) separates node count, replication, placement, routing, resource
   budget, network and failure; it requires independent hosts and balanced endpoints, and defines
   correctness and negative controls. Four separately claimable execution tasks are published
   (`…-topology-env-harness`, `…-node-replication`, `…-placement-routing`, `…-network-failure`);
   **none may run on this machine** — it needs a second host, so the four tasks are deferred
   ([`docs/optional-todos.md`](docs/optional-todos.md)); this study is the path to close
   `v2-gap-05` and `v2-gap-07` when such hosts exist.

7. **Study 06 native models** — protocol complete 2026-09-22: `studies/06-native-models/HANDOFF.md`
   (tag `study-06/v0-handoff`) compares PostgreSQL 17.11, MongoDB 8.0, ScyllaDB 6.2 and Valkey 8.1
   (authoritative, not a cache) under identical configuration-domain semantics, operations, resources
   and correctness gates; v0 is single-node, cluster variants belong to Study 07. Four separately
   claimable execution tasks are published (`…-native-models-harness`, `…-document-vs-relational`,
   `…-widecolumn-vs-relational`, `…-keyvalue-vs-relational`); none has run. Path to close
   `v2-gap-06-native-datastore-families`. **Deferred**: the engine images (MongoDB, ScyllaDB,
   Valkey) are not on this machine; see [`docs/optional-todos.md`](docs/optional-todos.md).

8. **Study 08 read models** — protocol complete 2026-09-22:
   `studies/08-analytics-read-models/HANDOFF.md` (tag `repo/analytics-read-model-handoff-v1`)
   compares base OLTP, rollup/materialized views, an in-engine search copy and a ClickHouse analytical
   copy for top-N/dashboard workloads, measuring query gain, refresh/maintenance, staleness, storage
   and correctness. Four separately claimable execution tasks are published
   (`…-analytics-harness`, `…-rollup-arms`, `…-search-copy`, `…-columnar-copy`); none may run on
   this machine — they need ClickHouse and a search engine image, so they are deferred
   ([`docs/optional-todos.md`](docs/optional-todos.md)).

9. **Study 09 hierarchy** — protocol complete 2026-09-22: `studies/09-hierarchy/HANDOFF.md` (tag
   `repo/hierarchy-study-handoff-v1`) compares adjacency list, materialized path, closure table,
   nested sets and bounded embedding under identical operations, covering reads, moves, inserts,
   deletes, storage and concurrency with cycle/orphan/depth-bound controls. Four separately claimable
   execution tasks are published (`…-hierarchy-harness`, `…-reads-storage`, `…-writes-moves`,
   `…-concurrency-engine`); none has run.

## Where the long history lives

- Archived `CONTEXT.md` editions: [`docs/history/context/`](docs/history/context/)
- Imported historical receipts (EH-01, EH-02, Studies 03–05): [`docs/ai-work/tasks/`](docs/ai-work/tasks/)
- Original handoffs and review receipts: [`docs/handoffs/`](docs/handoffs/) (read-only history)
- Per-study detail: each study's `reports/`, `HANDOFF.md`, `PROGRESS.md`, `ESCALATIONS.md`
