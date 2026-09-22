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
- **Active book evidence source:** `book/evidence/v2/claims.json` (the correction package, tag
  target `repo/book-evidence-registry-v2`; see [`book/evidence/README.md`](book/evidence/README.md)).
  v1 (`book/evidence/claims.json`, tag `repo/book-evidence-registry-v1`) is **frozen historical
  reference only** — its claim-to-cell resolution failed (36 of 43 names resolve nowhere) and no
  book prose may be written from it. `tools/evidence validate` defaults to v2.
- **Gate on the next book task:** the correction package is integrated (tag
  `repo/book-evidence-registry-v2`), so `book-synthesis-v1` may now be released; it must read
  `book/evidence/v2/` only.
- **Book toolchain:** `book/` holds the modular Typst reader layer, pinned by
  `book/Containerfile` to Typst `0.15.1` (image digest
  `sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f`). `book/build.sh`
  compiles and verifies the draft in Podman and writes `book/dist/build-manifest.json`; the toolchain
  build was 30 pages with 3/3 fonts embedded, all 23 v2 claim ids indexed and the registry digest
  present in the text.
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
- **Next ready work:** the six scientific protocol tasks (`study04-v2-protocol`,
  `study05-v2-protocol`, `analytics-read-model-protocol`, `hierarchy-model-protocol`,
  `native-major-model-protocol`, `topology-study-protocol`).
- **Brainstorm concluded** (tag `repo/book-evidence-and-model-brainstorm-v1`):
  `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` — two positions
  (gpt-5, deepseek-flash), two cross-reviews, one synthesis. It found registry v1 ineligible as the
  book's numeric source and decided one versioned, attributed evidence-correction package as the sole
  active claim source, with the fixed six-family taxonomy kept; the correction task implements it.

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

### Study 05 — external cache — **v1 measured, independently reviewed; v2 protocol next**

Cache-aside/write-through, strict vs relaxed freshness, wrongness of a relaxed cache and cache-side
invalidation fences, on PostgreSQL + Redis (Redis measured as a shared cache, never as an
authoritative store). The single-`small` all-green run is measured and analysed, and its limits are
explicit (single trial, unmeasured TTL/churn, no YugabyteDB/cluster/placement, no equal-total
accounting, open-loop demand absent). The independent HIGH review landed 2026-09-22
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
3. **Study 04 v2** — protocol complete 2026-09-22:
   `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md` (tag `study-04/v2-handoff`) defines
   the three review-mandated controls, the cardinality and cadence phases, repeated randomized
   trials, the equal-total-resource arm, verified placement and YugabyteDB coverage. Three separately
   claimable execution tasks are published (`…-study04-v2-controls`, `…-cardinality-cadence`,
   `…-remaining-designs`); none has run. Verdicts in
   `book/evidence/reviews/20260922-study04-independent-review.json`.
4. **Study 05 v2** — the v1 review is done; the open work is
   `task-20260922T025912Z-study05-v2-protocol` (real-TTL churn, medium scale, repeats, cache
   resource accounting, YugabyteDB/cluster/placement, open-loop demand) plus the review's four
   added requirements (repeated trials/second machine for the three-instance result, real-TTL
   churn, equal-total framing, executed placement pair). Verdicts in
   `book/evidence/reviews/20260922-study05-independent-review.json`.
5. **Book v1** — synthesis and independent release review not yet done; no final PDF yet.

## Where the long history lives

- Archived `CONTEXT.md` editions: [`docs/history/context/`](docs/history/context/)
- Imported historical receipts (EH-01, EH-02, Studies 03–05): [`docs/ai-work/tasks/`](docs/ai-work/tasks/)
- Original handoffs and review receipts: [`docs/handoffs/`](docs/handoffs/) (read-only history)
- Per-study detail: each study's `reports/`, `HANDOFF.md`, `PROGRESS.md`, `ESCALATIONS.md`
