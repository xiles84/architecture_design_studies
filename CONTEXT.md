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

### Book initiative and the AI work pipeline — **in progress, next task ready**

- Goal: [`goal-20260922T025912Z-data-architecture-book`](docs/ai-work/goals/goal-20260922T025912Z-data-architecture-book/GOAL.md)
  — an evidence-backed data-architecture book built with Typst in Podman, with every claim traceable
  to a run, tag, digest and signed analysis.
- Planning checkpoint integrated and tagged `repo/data-architecture-book-handoff-v1`.
- The queue CLI (v1) is **complete and integrated**, tagged `repo/ai-task-queue-v1`, with capability
  aliases tagged `repo/ai-capability-aliases`.
- This document's rewrite is `task-20260922T025912Z-history-status-reconcile`; it imports EH-01,
  EH-02 and Studies 03–05 by reference (see [`docs/ai-work/tasks/`](docs/ai-work/tasks/) and the
  archived history) and fixes public status drift.
- **Active book evidence source:** `book/evidence/v5/claims.json` (re-scopes five claims whose
  statements or limits asserted more than their own evidence supported; tag target
  `repo/book-evidence-registry-v5`; see
  [`book/evidence/README.md`](book/evidence/README.md)). v4 retired the two stale placement/endpoint
  gaps, and in doing so carries every still-active v3 claim
  forward content-identically and adds claim lifecycle metadata (`status`, `gap_kind`,
  `superseded_by`, `closed_dimensions`, `remaining_dimensions`). v3
  (`book/evidence/v3/`, tag `repo/book-evidence-registry-v3`), v2 and v1 are **frozen historical
  reference**; v1's claim-to-cell resolution failed (36 of 43 names resolve nowhere).
  `tools/evidence validate` defaults to v5; `validate-v5/-v4/-v3/-v2` pin a package explicitly.
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
→ [Study](studies/02-ticket-booking/) · [reports](studies/02-ticket-booking/reports/) · [REPORTS.md](studies/02-ticket-booking/REPORTS.md)

### Study 03 — reserved seating — **complete and HIGH-validated**

`venue → event → seat`: 40-minute holds, seat-map reads, all-or-nothing multi-seat blocks, 14 designs.
The v1 matrix and both repair runs are measured, reported, analysed and signed; the v2
operational-reports runs (3b) are measured, analysed and HIGH-validated (2026-09-22): r06 is
unanswerable in every design, L2 answers the 100k-seat map 4.5-11x the row designs on YugabyteDB,
and L3's point lookups collapse to 26 ops/s.
→ [Study](studies/03-reserved-seating/) · [reports](studies/03-reserved-seating/reports/)

### Study 04 — configuration portal — **v1 measured, independently reviewed; v2 protocol next**

`product → installed product → configuration entry`: 18 designs planned, 10 implemented, measured on
PostgreSQL (`small`) with two negative controls. The independent HIGH review landed 2026-09-22
(`task-20260922T025912Z-study04-independent-review`, signed analysis
`20260921T1215Z-small--deepseek-flash--2026-09-22`): it accepted the 93.7% control loss, the
15.2-15.8x rollup read gain at a 6.85x write cost, and the deterministic 10.5x document storage
win; it narrowed the 1.76x lock advantage to an upper bound and found that the ~1.7x identical-SQL
attribution floor is real but its proposed cell-order cause is not corroborated by the recorded
order. Next: the v2 completion protocol (`task-20260922T025912Z-study04-v2-protocol`, now unblocked),
which must recover the missing designs and add cardinality, cadence/churn, repeated randomized
trials, equal-total resources, placement and YugabyteDB, plus the three controls the review
requires (repeated-design instrument control, `-retries 1`, writer sweep).
→ [Study](studies/04-configuration-portal/) · [reports](studies/04-configuration-portal/reports/)

### Study 05 — external cache — **v1 reviewed; v2 runs landed (churn/scale, resources/engines, open-loop)**

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
