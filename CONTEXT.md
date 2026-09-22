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
- **Next ready work:** `task-20260922T025912Z-history-status-reconcile` (this task). Then
  `evidence-registry` → the three HIGH reviews + `typst-toolchain` → `book-synthesis-v1` →
  `book-release-review-v1`, then the six scientific protocol tasks.
- **Book release tag target:** `repo/data-architecture-book-v1` (not yet created).
- **Brainstorm concluded** (tag `repo/book-evidence-and-model-brainstorm-v1`):
  `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` — two positions
  (gpt-5, deepseek-flash), two cross-reviews, one synthesis. It found registry v1 ineligible as the
  book's numeric source (11 of 12 numeric claims have at least one declared cell name absent from the
  report they cite; 36 of 43 names resolve nowhere; `failed_cells: 0` where the report says 2 failed)
  and decided one versioned, attributed evidence-correction package as the sole active claim source,
  with the fixed six-family taxonomy kept. No measurement was authorized and **no tasks were
  created**; the owner's next command is
  `CREATE TASKS FROM BRAINSTORM brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions`.

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

### Study 04 — configuration portal — **v1 measured; review and v2 protocol open**

`product → installed product → configuration entry`: 18 designs planned, 10 implemented, measured on
PostgreSQL (`small`) with two negative controls. Open: the independent HIGH review of the proposed
book claims, then the v2 completion protocol (which must recover the missing designs and add
cardinality, cadence/churn, repeated randomized trials, equal-total resources, placement and
YugabyteDB).
→ [Study](studies/04-configuration-portal/) · [reports](studies/04-configuration-portal/reports/)

### Study 05 — external cache — **v1 measured; review and v2 protocol open**

Cache-aside/write-through, strict vs relaxed freshness, wrongness of a relaxed cache and cache-side
invalidation fences, on PostgreSQL + Redis. The single-`small` all-green run is measured and analysed,
and its limits are explicit (single trial, unmeasured TTL/churn, no YugabyteDB/cluster/placement,
open-loop demand). Open: the independent HIGH review, then the v2 protocol.
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
2. **Evidence registry** — v1 (`repo/book-evidence-registry-v1`) is schema-valid but not eligible
   as the book's direct numeric source; the correction package is
   `task-20260922T170845Z-book-evidence-correction-v2`, and the phase-3b claim verdicts it must
   honour are in `book/evidence/reviews/20260922-phase-3b-high-validation.json`.
3. **Study 04 v2** — missing designs, repeated randomized trials, cardinality, cadence/churn,
   equal-total resources, placement, YugabyteDB.
4. **Study 05 v2** — real-TTL churn, medium scale, repeats, cache resource accounting,
   YugabyteDB/cluster/placement, open-loop demand.
5. **Book v1** — synthesis and independent release review not yet done; no final PDF yet.

## Where the long history lives

- Archived `CONTEXT.md` editions: [`docs/history/context/`](docs/history/context/)
- Imported historical receipts (EH-01, EH-02, Studies 03–05): [`docs/ai-work/tasks/`](docs/ai-work/tasks/)
- Original handoffs and review receipts: [`docs/handoffs/`](docs/handoffs/) (read-only history)
- Per-study detail: each study's `reports/`, `HANDOFF.md`, `PROGRESS.md`, `ESCALATIONS.md`
