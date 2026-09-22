# Project context

The living state of this repository. Updated whenever a study starts, finishes, or
changes shape — so that anyone (or any future session) picking this up knows where things
stand without reading the git log.

**Last updated:** 2026-09-22

---

## Rule for every AI session: commit and tag, never push

**From 2026-09-13, by the owner's instruction:** AI sessions working in this repository
**commit and tag their own work, and never push** (the owner pushes and pulls). Reports of
the same study will reach different conclusions as studies are enhanced with new designs and
questions; tags are how a reader finds out which code, SQL and data produced each conclusion.
Untagged work means an unexplained change in conclusions.

- Commit at every meaningful step; never end a session with finished work uncommitted.
- `run/<study>/<run-id>` tags the commit behind each run (`run-study.sh --tag`).
- `study-NN/vX-<label>` tags study milestones — and, when enhancing a study, the state before
  and after the enhancement, so `git diff <before> <after>` explains a change in conclusions.
- Annotated tags only; never move, delete or reuse a tag; never rewrite history; never push.

Full rule: [AGENTS.md → "Hard rule: commit and tag, never push"](AGENTS.md).

## Rule for every AI session: agents from any vendor, working in parallel

**From 2026-09-13, by the owner's instruction:** agents from different vendors (Claude,
OpenAI and others) work in this repository, sometimes at the same time.

- **Instructions live in [`AGENTS.md`](AGENTS.md)**, the tool-neutral file. `CLAUDE.md` only
  imports it. Never write a rule only in `CLAUDE.md`.
- **Every task starts in its own worktree:** create one from local `main` or reuse this
  task's existing folder/branch after checking ownership and status. Never share a
  working folder or develop directly on `main`.
- **Every completed task merges into local `main`:** the owner authorized this as the
  default on 2026-09-14. Resolve conflicts and validate in the task worktree, preserve
  all sessions' work, then integrate into `main` and tag the result. No extra routine
  merge approval is required; pushing and pulling remain the owner's responsibility.
- **Always assume another AI is active:** work only in the task's own worktree, never
  disturb another task's branch, files, containers or run, and inspect the shared lock
  before measurement. The agent merging later reconciles current status, lessons, rules,
  indexes and terminology so `main` reads as one project while attributed artifacts and
  disagreements remain intact. The canonical rule is in
  [AGENTS.md](AGENTS.md#hard-rule-assume-a-concurrent-agent-the-later-merger-reconciles).
- **One measurement at a time on this machine.** Code and analyses in parallel; dev checks
  and matrices never. Runners take the **benchmark lock** (`run_lock_acquire` in
  `infra/lib.sh`, a podman volume named `ads-run-lock`); a second runner from any worktree
  or tool refuses to start and names the holder. If the lock is held, wait.
- **Before a long run, add it to the runs table below as "running".**

### Tags so far

| Tag | Marks |
|---|---|
| `study-04/v0-handoff` | study 04 (configuration portal) Execution Handoff EH-04 rev 1, escalation log, progress log, scoped LF policy |
| `study-04/v0.1-handoff-amendment-01` | study 04: AM-01 executed — WSL reaches the documented Podman engine; phase B bind-mount probe |
| `repo/wsl-podman-bridge` | infra/lib.sh: engine resolver + `podman()` + `winpath()`; studies 01-03 untouched |
| `repo/open-loop-arrivals` | platform: the open-loop arrival driver cadence experiments need |
| `repo/data-architecture-book-handoff-v1` | HIGH planning checkpoint for the evidence-backed Typst book and durable AI task queue; queue v1 is the next LOW task |
| `repo/ai-capability-aliases` | session declarations accept leader/master and worker/follower/slave aliases while queue state remains canonical HIGH/LOW |
| `study-04/v1-harness` | study 04: 10-design SQL catalogue, harness, gate, ledger, audits, report generator; dev-checked |
| `run/04-configuration-portal/20260921T1215Z-small` | commit that produced study 04's small `pg-single` matrix |
| `study-04/v1-measured` | study 04: small matrix measured and reported; both controls fired |
| `study-04/v1-analysis` | study 04: signed analysis and mechanism discussion companion |
| `study-01/v1` | study 01 code and reports behind its first analysis (`795420b`) |
| `study-01/v2-second-analysis` | study 01 with the independent GPT-6 analysis and regenerated reports (`918a90f`) |
| `study-01/v2-before-enhancements` | study 01 before the v3 follow-ups (GPT-6 session) |
| `study-01/v3-harness` | seven controlled variants, independent-load runner and initial validation |
| `study-01/v3-verified` | all new variants pass PostgreSQL and YugabyteDB verification |
| `study-01/v3-mechanisms` | five-load D2/D3 results, concise analysis and mechanism companion |
| `study-01/v3-followup-results` | 135-cell growth, contention and deployment result set |
| `study-01/v3-enhancements` | completed local follow-ups, matched memory control, concise final analysis and signed discussions |
| `study-02/v1-harness` | study 02 designs, harness and shared platform as first run (`774d258`) |
| `run/02-ticket-booking/20260913T021206Z` | commit that produced study 02's `small` matrix (`7570648`) |
| `study-02/v1-analysis` | study 02 report and first signed analysis (`2bc7e1c`) |
| `repo/agents-md-and-run-lock` | AGENTS.md as the tool-neutral instructions, parallel-work rules, benchmark lock |
| `repo/study-comparison-minimums` | future-study requirements for calculated sizing, information placement, data colocation and concurrency strategies |
| `repo/worktree-to-main-workflow` | integrated task workflow: every task uses an isolated worktree and merges completed changes into local main |
| `study-01/v3-integrated` | integrated Study 01 v3 local enhancements, final analysis and discussion companions |
| `study-01/v3-high-review` | HIGH acceptance checkpoint for the integrated v3 evidence and EH-02 revision 1 |
| `study-01/v3-high-review-r2` | amended HIGH checkpoint that preserved the peer's main-checkout reservation through Study 03 step 11 |
| `study-01/v3-reviewed` | final integrated Study 01 v3 state: accepted HIGH review, reconciled project context and LOW completion receipt |
| `repo/study01-integration-handoff-v1` | HIGH planning checkpoint for the LOW integration steps and subsequent HIGH review |
| `study-02/v1.1-counter-discussion` | study 02: discussion companion "counter cost vs index cost" and a dated note in the analysis answering the owner's question |
| `study-03/v0-handoff` | study 03 (reserved seating) Execution Handoff and escalation log, before any code |
| `study-02/v1.2-terminology` | study 02 README terminology section (general admission); no SQL, harness, report or analysis change |
| `repo/platform-lock-not-available` | platform: `ErrLockNotAvailable` (SQLSTATE 55P03) appended for study 03's NOWAIT design |
| `study-03/v0.1-handoff-amendment-01` | study 03: ER-01 decided (transient refusals on YugabyteDB); AM-01 adds design S1r, refusal diagnostics and amended dev-check criteria |
| `study-03/v1-harness` | study 03: 14 designs (13 + S1r), harness and SQL as dev-checked on tiny (dc1–dc14) |
| `study-03/v0.2-handoff-amendment-02` | study 03: ER-02 decided; AM-02 adds the `small` calibration run and duration rules for steps 9–10 |
| `repo/concurrent-agents-reconciliation` | AGENTS.md hard rule: always assume a concurrent agent; the agent that merges later reconciles shared documents into one coherent project |
| `run/03-reserved-seating/20260915T002411Z` | commit that produced study 03's `small` main matrix (`7dbdd11`) |
| `run/03-reserved-seating/20260915T173255Z` | commit that produced study 03's repeated race (`52a9975`) |
| `study-03/v1-measured` | study 03 measured: main matrix, repeated race, diagnoses, generated reports, documents; ready for analysis (ER-03, ER-04 open) |
| `study-03/v0.3-handoff-amendment-03` | study 03: ER-03 and ER-04 decided; AM-03 adds tolerant instrumentation, an L2 refusal diagnostic and a repair run |
| `run/03-reserved-seating/20260915T232736Z` | commit that produced study 03's ER-03 repair run (`4044dd5`) |
| `run/03-reserved-seating/20260916T000706Z` | commit that produced study 03's ER-04 repair run (`4044dd5`) |
| `study-03/v1.1-repairs` | study 03: AM-03 harness, both repair runs, regenerated reports, and results/reports pinned to LF |
| `study-03/v1-analysis` | study 03: signed analysis of the small matrix, report index, README and context complete |
| `repo/recency-reports-handoff-v1` | HIGH planning checkpoint: study 01's recency protocol, studies 02/03 reporting protocols and EH-02, before any code |
| `study-01/v4-harness` | study 01: D18–D24, catalogue/harness/runner/diagram changes for the recency question, dev-checked on `tiny` (PostgreSQL and YugabyteDB, single and 3-node) and calibrated on `small`; no reported run yet |
| `run/01-charity-tree/20260916T090036Z-v3` | commit that produced study 01's measured recency matrix (192 cells, 0 failed) |
| `study-01/v4.1-handoff-amendment-01` | EH-02 AM-01: phase 1 validated and accepted; five repairs specified before phase 2 |
| `study-01/v4-measured` | study 01: recency matrix measured (AM-01 repairs applied and held under the real run); ready for analysis |
| `study-01/v4-analysis` | study 01: signed analysis of the recency matrix (digest `1b05f142fd9e06ef`), report index regenerated |
| `study-01/v4.2-handoff-amendment-02` | EH-02 AM-02: phase 3a authorised (studies 02/03 operational reports, implementation and dev checks); two protocol claims verified across all 28 designs; measured runs still gated |
| `study-02/v2-reports-devchecked` | study 02: r01-r06 and the new X1 design (append-only sale ledger) implemented, dev-checked correct on PostgreSQL and YugabyteDB `tiny`; phase 3b (measured runs) not yet authorised |
| `study-03/v2-reports-devchecked` | study 03: r01-r06 implemented (r06 unanswerable everywhere, verified), dev-checked correct on PostgreSQL and YugabyteDB `tiny`, including L3; phase 3b not yet authorised |
| `repo/recency-reports-handoff-am03` | EH-02 AM-03: phase 3a reviewed. Eleven defects found (X1 refunds lose their buyer; ledger audit unproven); repairs, `small` calibration and four phase-3b runs specified. The two `v2-reports-devchecked` tags above mark the unrepaired state |
| `study-02/v2.1-reports-repaired` | study 02: the AM-03 repairs plus their dev-check evidence and records, on one commit. Annotation records that RR-ER-01 was still open at that state and was decided in AM-04 (`study-02/v2.2-ledger-attribution`) |
| `study-03/v2.1-reports-repaired` | study 03: the AM-03 repairs (window parameters bound in `ExplainAll`, value comparisons, r06's unanswerability pinned by test) plus their dev-check evidence |
| `run/02-ticket-booking/20260920T234953Z` | commit that produced study 02's operational-reports matrix (45 cells, 0 failed) — the phase-3b run the 2026-09-20 session finished but did not commit |
| `run/02-ticket-booking/20260921T182546Z` | **marks no run**: created by a launch that then refused on the benchmark lock. No cells, no manifest, no report |
| `run/02-ticket-booking/20260921T183722Z` | commit that produced the P3 → X1 race pair at 32 buyers (6 cells, one partial on `yb-cluster3`) |
| `run/02-ticket-booking/20260921T205212Z` | commit that produced the P3 → X1 race pair at 128 buyers (6 cells, 0 failed) |
| `run/03-reserved-seating/20260921T231328Z` | commit that produced study 03's operational-reports matrix (41 cells, 0 failed) |
| `study-02/v2-measured` | study 02 v2 measured: the reports matrix and both race arms, each with its report |
| `study-03/v2-measured` | study 03 v2 measured: the reports matrix (digest `2f5b619430c2e7a7`) with its report |
| `study-02/v2-analysis` | study 02 v2 analysed: four signed analyses (reports, 32 buyers, 128 buyers) and the P3 → X1 discussion companion, report indexes regenerated |
| `study-03/v2-analysis` | study 03 v2 analysed: the signed analysis of the reports matrix, report index regenerated |
| `study-05/v0-handoff` | study 05 (external cache) Execution Handoff EH-05 rev 1, escalation log, progress log, scoped LF policy, Redis image pin |
| `study-05/v1-harness` | study 05: scenario registry, both cache backends, lease, oracle, wrong-read accounting, faults, report, containerised runner |
| `run/05-cache-consistency/20260921T-survey` | commit that produced study 05's `small` `pg-single` survey (`567778f`) |
| `study-05/v1-measured` | study 05: `small` survey measured (27 cells: 21 core + 3 reference + controls), reported |
| `study-05/v1-analysis` | study 05: signed analysis and mechanism companion, including the nine failed cells and the coverage gaps |

Check `git tag -n1` for the authoritative list; this table can lag behind a session that
has not updated it yet.

---

## Requirements for studies from now on — owner clarification, 2026-09-14

The canonical requirements are in
[AGENTS.md](AGENTS.md#required-comparisons-for-future-studies), with implementation
guidance in methodology sections 4 and 6a. Every new study and future measurement plan
for an existing study records:

- Calculated database-node and client-worker sizing, per-host resource accounting,
  an equal-total-budget cluster control, and measured endpoint/client bottlenecks.
- Rollup, rolldown and embedding comparisons wherever applicable, including the cost
  of maintaining the duplicated information and its correctness.
- Colocated versus non-colocated data in multi-node database scenarios, verified by
  physical-placement evidence and controlled independently of data and resource budget.
- Optimistic and pessimistic concurrency as the minimum when concurrency is involved,
  with additional strategies where applicable.

Protocols map each requirement to planned/measured pairs or an explicit reason/gap.
These are prospective requirements; they do not imply that previous runs measured
every dimension. Context and lessons alone had not made all four requirements mandatory;
the new AGENTS.md section does so explicitly.

**Study 01 v3 HIGH review accepted and integrated, 2026-09-20 — complete:**
EH-01's final local integration is
verified at `db77fe60b5bfe78b29ec04db9150f33cc1a0bb49`. Both annotated tags,
`repo/worktree-to-main-workflow` and `study-01/v3-integrated`, resolve to that commit,
which is an ancestor of current local `main`. The 261 saved cells retain their producing
commits, images, environment and four report digests. The six expected D9 failures stay
excluded from performance conclusions. HIGH accepts the existing local enhancements,
concise final analysis, signed discussions and prospective study requirements.

The task branch incorporated committed main through `7123da6`, including the completed
and signed Study 03 v1 analysis and the repository-wide concurrent-agent reconciliation
rule. Its former main-checkout reservation ended with Study 03 step 11. LOW reconciled
the living documents without changing attributed artifacts, committed merge checkpoint
`0e7c05fa0b112ab33023b07af032500dab4aae8c`, and fast-forwarded local `main` to it under
the owned integration lock after the concurrent dev check released the runtime. The
separate `.worktrees/recency-reports` branch, files, commits and containers were not
changed or adopted. This completion receipt is the final mapped EH-02 update and is
tagged `study-01/v3-reviewed`; no further HIGH pass or model switch is required. The
[Execution Handoff](docs/handoffs/20260914-study01-integration/HANDOFF.md),
[progress](docs/handoffs/20260914-study01-integration/PROGRESS.md),
[execution receipt](docs/handoffs/20260914-study01-integration/EXECUTION_RESULT.md) and
[escalation log](docs/handoffs/20260914-study01-integration/ESCALATIONS.md) preserve the
iteration; the [HIGH review](docs/handoffs/20260914-study01-integration/HIGH_REVIEW.md)
records the acceptance evidence. Study 03's completed state and decisions remain in its
own section and logs.

## Model roles for every study

[AGENTS.md](AGENTS.md#model-roles-execution-handoff-and-escalation-required) defines
the owner's general workflow: HIGH plans and writes an Execution Handoff, LOW executes
its mapped steps, unmapped decisions become Escalation Required, and HIGH validates and
analyses. Each iteration names the next model level and any agreed effort. Routine
waiting remains LOW work. No automatic model change is implied.

Current task: HIGH planning/review = GPT-6 via Codex (exact selected effort not exposed);
EH-01's executor recorded GPT-5 via Codex desktop under the user-selected LOW role
(exact effort not exposed). HIGH accepted EH-01, and GPT-5 completed EH-02's mapped
reconciliation, safe local-main integration and final tag in the LOW role. This task has
no remaining model iteration.
Study 03's existing model/effort mapping below remains specific to that study.

## What this project is

An experimental record of **how architecture and data-modelling decisions actually affect
performance**, measured rather than asserted.

The distinguishing constraints:

- Claims are backed by **real runs**, not reasoning about what ought to be faster.
- Every test runs **inside Podman containers**, so a reader can reproduce it without
  installing anything but Podman.
- Every result names the **environment** it came from, documented down to the CPU
  topology and the specific ways that machine can mislead a benchmark.
- Every study produces a **report**; superseded reports move to `reports/outdated/`
  rather than being deleted.
- **Correctness gates timing.** A design that answers wrongly is never reported as fast.

Implementation language for tooling is **Go**.

## Repository layout

```
docs/
  methodology.md          the rules every study follows
  replication.md          how a reader reproduces a run
  environments/           one page per machine results were produced on
infra/
  versions.env            pinned images + per-node resource budget
  lib.sh                  shared podman helpers
  pg-single.sh            PostgreSQL, 1 node
  yb-single.sh            YugabyteDB, 1 node (RF=1)
  yb-cluster3.sh          YugabyteDB, 3 nodes (RF=3)
platform/                 shared Go module (adsplatform): core / ports / adapters
studies/
  01-charity-tree/        study 01 (see below) — own harness, predates platform/
  02-ticket-booking/      study 02 (see below) — built on platform/
  03-reserved-seating/    study 03 (see below) — built on platform/; measured and analysed (v1-analysis)
```

Everything belonging to one study (SQL, harness, runner, image name, results, reports,
analyses) lives in that study's directory; only study-independent code is shared, in
`platform/`, laid out as ports and adapters (`platform/README.md`).

**Versioning.** Tags: `study-01/v1` marks the code behind study 01's published results;
`run/<study>/<run-id>` marks the commit a run was produced from (`run-study.sh --tag`, study
02 onwards). Results, manifests and reports record commit, `git describe` and a dirty flag.

Each study directory holds: `README.md` (the question and the designs), `sql/` (one
directory per design), `diagrams/` (PlantUML sources + rendered SVG), `harness/` (Go),
`run-study.sh`, `results/<run-id>/`, `reports/`.

## Current state

### Task — data architecture book and durable AI work queue (2026-09-22): **planning published, merged into local `main`**

The owner's request and HIGH interpretation are preserved under
[`goal-20260922T025912Z-data-architecture-book`](docs/ai-work/goals/goal-20260922T025912Z-data-architecture-book/).
The authoritative operating contract is [`docs/ai-work/WORKFLOW.md`](docs/ai-work/WORKFLOW.md),
with JSON conventions in [`SCHEMA.md`](docs/ai-work/SCHEMA.md). The publication checkpoint
is tagged `repo/data-architecture-book-handoff-v1`.

The first executable item is
[`task-20260922T025912Z-queue-v1`](docs/ai-work/tasks/2026/09/task-20260922T025912Z-queue-v1/),
now `ready` for a LOW executor. It implements atomic claims, leases, recovery, guards,
event validation, and the local-main integration lock. The book, evidence registry,
history import, status reconciliation, reviews, and synthesis remain queued behind it.
This planning checkpoint started no benchmark and does not itself provide queue CLI or
book implementation evidence.

Owner amendment `20260922T101113Z-capability-aliases` lets a new chat declare `leader`
or legacy `master` for HIGH, and `worker`, `follower`, or legacy `slave` for LOW. The
pipeline normalizes these inputs to canonical `HIGH`/`LOW`; it does not use
`primary`/`replica`, which remain datastore-topology terms. The queue-v1 task has a
linked implementation amendment and must include alias parser tests.

### Task — recency question and operational reports, EH-02 (2026-09-15 → 2026-09-21): **complete, merged into `main`**

**Worktree `.worktrees/recency-reports`, branch `repo/recency-and-reports`, from `7123da6`.**
Owner's request of 2026-09-15: study 01 gains the question *"which people made their last
donation in a period"* (their own suggestion — a `last_donation` flag on `donation`,
indexed — measured as one design among several), and studies 02 and 03 gain the reports a
real operation asks for.

- Planning is published: [`studies/01-charity-tree/RECENCY.md`](studies/01-charity-tree/RECENCY.md),
  [`studies/02-ticket-booking/REPORTS.md`](studies/02-ticket-booking/REPORTS.md),
  [`studies/03-reserved-seating/REPORTS.md`](studies/03-reserved-seating/REPORTS.md) and
  [EH-02](docs/handoffs/20260915-recency-and-reports/HANDOFF.md) with its
  [progress](docs/handoffs/20260915-recency-and-reports/PROGRESS.md) and
  [escalation log](docs/handoffs/20260915-recency-and-reports/ESCALATIONS.md).
- **Phase 1 complete (LOW, Claude Sonnet 5, 2026-09-15/16):** four new statements
  appended to every existing catalogue (D1–D17, byte-identical schema/indexes/writes),
  seven new designs (D18–D24), loader/verifier/audit/report/runner changes, dev checks
  on `tiny` across all three topologies, and a `small` calibration cell. Three real bugs
  were found and fixed by the dev-check gate itself: the embedded designs (D6/D9/D10)
  read the wrong JSON key for the new questions; `ExplainAll`'s plan-capture path had a
  second, separate fixed-parameter map that never got `since`/`until`; the "arrival"
  experiment's audit point never called the new `AuditRecency`. All three are fixed,
  committed and re-verified. **D21 (the negative control) was seen to fail twice**: under
  ordinary 8-connection spread demand (3/500 donors) and, after the paced hot-donor
  experiment needed 16 writers instead of 8 to reproduce it, there too (1/500, with D20
  staying consistent under the identical contention) — see `LESSONS_LEARNED.md` for why
  the paced experiment needed more writers. All 22 PostgreSQL and 24 YugabyteDB-capable
  designs pass the correctness gate (20/20 checks: the original twelve plus q13–q16 in
  both window regimes) on all three topologies. The partial index in D20 was accepted by
  YugabyteDB without the mapped fallback. Calibration: D3's four recency questions run at
  160–450 ops/s at `small` scale (no supporting index, as expected), client CPU ~4% of
  its budget (not the bottleneck). Full [progress log](docs/handoffs/20260915-recency-and-reports/PROGRESS.md).
  Checkpoint tag: `study-01/v4-harness`. Not merged to `main`; no reported run yet.
- **Phase 1 validated by HIGH (2026-09-16) and accepted**, with five repairs required
  before any measured run ([AM-01](docs/handoffs/20260915-recency-and-reports/HANDOFF.md#amendments)).
  The repairs are all in the path a *reported* run takes, which no dev check exercised:
  the maintenance group passes four write ops where the experiment path accepts one;
  `experiment.go` never calls the new `AuditRecency`; `experimentProblems` cannot see a
  recency-audit failure, so the negative control could fire invisibly (a methodology 5a
  violation in the reporting path); the new reporting section lives in `report.go` while
  the runner generates its report with `report_enhancements.go`; and `delete_person`
  was never exercised on the flag designs, where the delete trigger fires once per child
  row. Three of the five follow from EH-02's own under-specification of the reporting
  path, not from execution.
- **Phase 2 complete (2026-09-16, LOW, Claude Sonnet 5).** All five AM-01 repairs
  committed (`ff79475`, `2c6f515`) and `delete_person` dev-checked clean on D20/D21
  first (AM-01.5: 6,000+ ops/s, 0 errors, audits consistent, no escalation). The
  measured recency matrix ran via `run-study.sh --suite enhancements --experiments
  recency-reads,recency-maintenance,recency-hot-donor,recency-placement --trials 3
  --tag`: **192 cells, 0 failed processes**, well under the ~4-hour estimate. Run tag
  `run/01-charity-tree/20260916T090036Z-v3` (auto-created by the runner). The AM-01
  repairs held under the real run — D21's negative control fired in every
  `recency-maintenance` trial and in the `hot-donor-w16` sweep, visible in the generated
  report exactly where AM-01.3/1.4 put it. One new finding for the analysis: D23
  (optimistic CAS) recorded write errors (up to 18/10000) under the hottest single-donor
  contention — a genuine abort-under-contention measurement, not a defect. Report:
  [`reports/20260916T090036Z-v3.md`](studies/01-charity-tree/reports/20260916T090036Z-v3.md).
  Results: `results/20260916T090036Z-v3/`. Nothing is running; lock released.
- **Analysis published (2026-09-16, HIGH, Claude Opus 5).**
  [`20260916T090036Z-v3--claude-opus-5--2026-09-16`](studies/01-charity-tree/reports/analyses/20260916T090036Z-v3--claude-opus-5--2026-09-16.md),
  digest `1b05f142fd9e06ef`. **Answer to the owner's question:** the `last_donation` flag
  works and is cheap (+0.95% storage, 6,219 inserts/s), but where the parent already
  carries `last_donation_at` — as D4/D5 have since the original survey — **one index on it
  is 2.9x faster on PostgreSQL and 7.9x faster on YugabyteDB**, with no new column, trigger
  or write path. The flag's guard (a person-row lock) costs nothing measurable while the
  unguarded control corrupted 3–7 donors per trial: take the lock. Embedding buys nothing
  for a cross-parent recency question (D6 is 4.9x *slower* than the plain aggregate); the
  per-donor index probe D18 is slower than the naive aggregate on PostgreSQL and 22x slower
  on YugabyteDB; colocation (D20 vs D24) is a null result. Two limitations the analysis
  names as the planner's own errors: no no-flag write baseline in the maintenance group,
  and a hot-donor sweep that was rate-limited below capacity. The report was regenerated
  with the run's own pinned image so its analyses index resolves; no measurement changed.
- **Phase 3a authorised (2026-09-16, HIGH):** studies 02 and 03's operational reports,
  **implementation and dev checks only**, per
  [AM-02](docs/handoffs/20260915-recency-and-reports/HANDOFF.md#amendments). Two protocol
  claims are now verified across all 28 designs rather than sampled, and stand as findings
  in their own right: **study 02 destroys the sale on cancellation in every design** (C/R/H
  delete the ticket, P resets it), so no design can report refunds and `r01`'s revenue is
  wrong across a refund; and **no study-03 design retains any hold history**, so the
  abandonment funnel cannot be computed at all. AM-02 adds the reports with zero change to
  any existing statement, schema, index or write path; declares answerability in Go and
  pins it to the SQL with a unit test; adds exactly one new design (study 02's X1 = P3 plus
  an append-only sale ledger, with a reconciliation audit); and adds none to study 03.
  Phase 3b — the measured runs — stays gated until the dev-check numbers exist.
- **Phase 3a implemented and dev-checked (2026-09-16, LOW, Claude Sonnet 5).** AM-02.1–.5
  executed. All 15 study-02 designs (14 + X1) and all 14 study-03 designs pass the
  correctness gate on PostgreSQL and YugabyteDB 1-node at `tiny`, and C1 and H0 still fire.
  The dev checks caught two bugs: a duplicated `CREATE TABLE` in X1's schema, and study 03's
  r01 truth missing lazy-expired holds.
- **Phase 3a reviewed: not accepted as is; phase 3b sized (2026-09-16, HIGH, Claude
  Opus 5, [AM-03](docs/handoffs/20260915-recency-and-reports/HANDOFF.md#amendments)).**
  Eleven defects. The central one: **X1 records every refund with a NULL buyer**, because
  `RETURNING` yields the post-update row. X1's ledger audit only counts, so it cannot see
  that; it never runs after the race; and in the dev checks it was only ever run on the
  ledger the loader seeded, where it is consistent by construction. Other defects:
  count-only report checks; study 02's r03 wrongly declared answerable across refunds; the
  explain phase lacking the window parameters (a repeat of this task's phase-1 lesson); and
  X1 missing from `pairs`, the diagram and the README. AM-03 orders the audit, then two
  proofs that it fires (on the unfixed SQL and on injected faults), then the fix, then a
  `small` calibration. Four measured runs are sized by rules with a 10 h guard: study 02's
  reports matrix, a P3 → X1 race pair at 32 and at 128/64 buyers, and study 03's reports
  matrix. Tags `study-0{2,3}/v2-reports-devchecked` mark the unrepaired state. **Next: LOW
  (Claude Sonnet 5) executes AM-03.** Nothing is running.
- **AM-03.1–.9 executed; RR-ER-01 open, blocking .10–.13 (2026-09-16/17, LOW, Claude
  Sonnet 5).** X1's refund attribution fixed — the buyer now comes from the ledger row of
  the sale being reversed, not the post-update `RETURNING` (which was NULL, and since
  `sale_event.customer_id` is `NOT NULL`, every cancellation on X1 was failing outright
  before this fix, confirmed by a dev check). The new `a_ledger_attribution` audit that
  proves this still has a false-positive gap: it orders each seat's history by `at`
  (PostgreSQL's `now()`, fixed at transaction *start*), which does not track commit order
  under the CAS retry loop's concurrent writers — confirmed by direct inspection
  (`sale_event_id` gives the only order consistent with the CAS invariant; `at` does not).
  **[RR-ER-01](docs/handoffs/20260915-recency-and-reports/ESCALATIONS.md) is open**:
  reordering the check was specified by HIGH for a stated YugabyteDB reason not yet tested,
  so LOW did not reorder it unilaterally. Everything else in AM-03 passed dev-checked:
  study 02's r03 declared `partial` on non-ledger designs; r02/r03 checked by value; both
  studies' `ExplainAll` fixed (repeating this task's own phase-1 lesson); X1 registered in
  `pairs`, the diagram and the README; all 29 designs' correctness gates pass on both
  engines with the repaired SQL; both negative controls still fire; every post-write report
  check (using the harness's own booked/cancelled counters, unaffected by the ordering bug)
  passes. **Calibration and the four measured runs are not started, correctly gated on
  RR-ER-01** — running them while the audit's own correctness is in question would produce
  numbers the correctness gate cannot yet vouch for. **Next: HIGH decides RR-ER-01.**
  Nothing is running; lock released.
- **RR-ER-01 decided; phase 3b unblocked (2026-09-20, HIGH, Claude Opus 5,
  [AM-04](docs/handoffs/20260915-recency-and-reports/HANDOFF.md#amendments)).** LOW's
  experiments proved the question was wrong, not just the answer: `at` is fixed at
  transaction START and `sale_event_id` is handed out in per-connection cached blocks by
  YSQL, so **neither reconstructs commit order**, each failing on a different engine
  (PostgreSQL clean but YugabyteDB deterministically wrong, 246 attribution problems;
  `ALTER SEQUENCE … CACHE 1` had no effect). LOW also found the same `at` bug in
  `w_cancel_ticket` itself, not only in the audit — X1's real write path could name a stale
  buyer on a refund. **Decision: stop reconstructing order.** The refund reads its buyer
  from the ticket row it is clearing, under `FOR UPDATE`, in the same statement; the
  attribution audit is rewritten order-free (per-buyer refunds bounded by sales, plus the
  live sale recorded verbatim). This revises AM-03.4's rejection of `FOR UPDATE`, which was
  mine: the cancelling `UPDATE` takes that same row lock moments later anyway, so the extra
  lookup is part of the ledger's cost, and AM-04.2 restates the P3 → X1 pair precisely
  rather than dropping the claim. LOW re-checks X1 only (both engines, including the
  high-contention setting that exposed the bug), then runs calibration and the four
  measured runs straight through. **Next: LOW (Claude Sonnet 5) at AM-04.1.**
- **PHASE 3B COMPLETE — nothing is running; the benchmark lock is free (runs 2026-09-20/21).** AM-04 passed (`5e3a3f8`, tag `study-02/v2.2-ledger-attribution`); calibration done (`am03-cal-small/SUMMARY.md`): study 02 reports at 10 s, study 03 at 5 s, race B2 = 128. All four runs are measured, committed, reported and analysed, one at a time under the lock:

| Run | Study | What | Cells | Digest | Producing tag |
|---|---|---|---|---|---|
| 3b-1 `20260920T234953Z` | 02 | every design's answerable reports (r01–r06) beside the buyer reads, 10 s | 45, 0 failed | `160bd80c49bbd892` | `run/02-ticket-booking/20260920T234953Z` |
| 3b-2 `20260921T183722Z` | 02 | P3 → X1 at 32 buyers, 3 race + 3 churn trials | 6, **1 partial** | `4a02ec8fb0991166` | `run/02-ticket-booking/20260921T183722Z` |
| 3b-3 `20260921T205212Z` | 02 | the same pair, X1 first, 128 buyers | 6, 0 failed | `6c88f6463c80c6ac` | `run/02-ticket-booking/20260921T205212Z` |
| 3b-4 `20260921T231328Z` | 03 | each layout's reports (r01–r06) beside six buyer reads, 5 s | 41, 0 failed | `2f5b619430c2e7a7` | `run/03-reserved-seating/20260921T231328Z` |

  - **What the four runs say**, in one paragraph each. **Study 02 reports (3b-1):** a ticketing model chosen to prevent overbooking cannot report refunds at all — r05 is `unanswerable` in 14 of 15 designs and r01/r03 are the wrong answer in those same 14, because cancellation erases the sale; X1's append-only `sale_event` answers all three (r05 32k ops/s, r01 32.4k) at +35 % storage, and the recency report r03 is 100–600x the cost of the buyer's reads. **Study 02 race pair (3b-2/3b-3):** on YugabyteDB the ledger costs about half the sell-out race (~2x) and 2.1–2.5x the isolated refund, consistently in both arms and at the 128-buyer contention that the pre-AM-04 audit could not survive; on PostgreSQL the two arms disagree in sign and the disagreement tracks buyer count and design order, which changed together, so that axis is reported as a range (1x–2.3x), not a point. Nothing under-booked anywhere. **Study 03 reports (3b-4):** r06 — the hold funnel that would size the 40-minute window — is unanswerable in every design on every topology; on YugabyteDB the section document answers the 100k-seat map 4.5–11x faster than per-seat rows while the section-sharded variant collapses point lookups to 26 ops/s; and r05 is again 50–100x every other report, the third study to show that shape.
  - **Analyses and milestones.** Signed by `deepseek-flash` as a new analyst (four analyses + one discussion companion), tagged `study-02/v2-analysis` and `study-03/v2-analysis`; the measured states are `study-02/v2-measured` and `study-03/v2-measured`. `study-02/v2-measured` sits on 3b-3's commit, the commit holding study 02's last report — AM-03.13 said to tag "after 3b-4", and since 3b-4 measures study 03 that reading gives each study a tag on its own measured state. **The handoff's HIGH validation of phase 3b (Claude Opus 5) remains open**, as does the invitation for a second analyst; the analyses state this themselves.
  - **Coverage gaps this phase did not close**, all named in the analyses: no control design was in the race runs (3b-2/3b-3 measure only P3 and X1) and the reports runs have no writing phase, so no negative control fired in any of the four; `yb-cluster3/x1_cas_ledger` in 3b-2 is partial (2-trial medians) after a YugabyteDB RPC timeout; single cell per design per topology; closed-loop buyers; quota-throttled engines; no open-loop or SLO work. The next experiments are listed in each analysis and the pair's discussion companion.
  - **Resumption note, 2026-09-21.** The 2026-09-20 session ended after 3b-1 finished (23:50:01Z → 02:38:52Z, 45 cells, 0 failed) but before committing its output; the results and report were left untracked in `.worktrees/recency-reports`, whose `.git` link had been written with a Windows path and needed `git worktree repair`. Both are dealt with: the link is repaired (and made *relative*, so WSL git and Git Bash's git.exe both resolve it — the previous form worked for only one of them), the CRLF phantom edits the checkout had accumulated were normalized from its own LF index, and 3b-1 is committed as `b0d9b69`. The task's own model-role map (HIGH = Claude Opus 5, LOW = Claude Sonnet 5) is not the session that finished it: `deepseek-flash` executed AM-03.12/AM-03.13 and signed the analyses as a new analyst.
  - **One run tag marks no run.** The refused 3b-2 launch of 2026-09-21 created `run/02-ticket-booking/20260921T182546Z` before dying on the lock; it has no cells, no manifest and no report. Tags are never moved or deleted, so it stays and is recorded here. Its empty results directory was removed and 3b-2 used a new id.
  - **Environment, and why runs go through Git Bash.** From WSL the branch's `infra/lib.sh` had no podman client, so `infra/lib.sh` now carries main's resolver (`2e31fa4`, 58 additive lines, no study code). The measured runs themselves still execute under Git Bash, exactly as 3b-1 and the AM-04 dev checks did, where `cygpath` makes `hostpath()` the Windows form `podman.exe` needs. 3b-2 and 3b-3 were produced with `a24c10933bd0`, study 03's run with `a145af66f330`; the four reports were regenerated with those pinned images and every diff is index-only.
  - **Integration receipt.** The task branch was integrated into local `main` by merging main into it and then fast-forwarding main: merge commit `df371ae` (parents `3cb3840` and `84ac2b9`), and **this receipt commit**, which carries the final integrated state and is tagged `repo/recency-reports-integrated`. Every task commit named on this page is reachable from `main` — verified with `git merge-base --is-ancestor` for `56edc26`, `b0d9b69`, `2e31fa4`, `dad86d0`, `a04fd20`, `ab1b76e`, `e190749`, `2dbaed9`, `b5f58e6`, `feb4d5f`, `3cb3840` and `cb62711`. **Reconciliation was confined to `CONTEXT.md` and `LESSONS_LEARNED.md`**, the only files both sides had edited: the tags table now lists both sides' tags once (no duplicate row), the task's heading says complete, and studies 04 and 05 keep their own sections, status and receipts untouched. Nothing under `studies/04-*`, `studies/05-*`, `infra/redis.sh` or the study-01 integration handoff differs from main — the merge brought them in unchanged. Checks in the merged tree, in the pinned `golang:1.26-bookworm` container: `go test ./...` green for the platform and studies 01–03, `go vet ./...` clean for all four, `bash -n` clean on the runners. `gofmt -l` flags four files — `platform/core/measure/arrival.go` (main's; identical to the branch base, so main changed it after this branch's base) and `studies/03-reserved-seating/harness/{ledger,lifecycle,load}.go` (identical to main and to the base): all pre-existing on `main`, none in this task's footprint, and another study's files were not reformatted. Nothing was pushed or pulled, and no tag was moved, deleted or reused.
  - **Worktrees retired, 2026-09-21.** The task is finished and its branch is merged, so the worktrees it no longer needs were removed: this task's own `.worktrees/recency-reports`, study 04's `.worktrees/study04-configuration-portal` and study 05's `.worktrees/study05-cache-consistency` (all clean, all on merged commits), the empty leftover `.worktrees/study03-measurement` folder, and the dead registration for `C:/Users/xiles/.codex/worktrees/study05-cache-prompt/…` whose directory no longer existed. All ten local branches and all their commits remain in the repository; only working folders were removed. `git worktree list` now shows the main checkout alone. **One residue, recorded rather than hidden:** the empty directory `.worktrees/recency-reports` could not be deleted — Windows reports "being used by another process" and WSL "Permission denied / Device or resource busy" — because a live process still holds it as its working directory, almost certainly the terminal that ran phase 3b. It is empty, it is in `.git/info/exclude`, git has no registration for it, and it disappears once that terminal is closed (or from Explorer, or on the next reboot); nothing depends on it.

### Study 01 — tree structures (charity → person → donation)

**v3 local enhancements completed (2026-09-14):** the owner approved all six follow-ups in
GPT-6's "What I would measure next". The protocol is in
[`ENHANCEMENTS.md`](studies/01-charity-tree/ENHANCEMENTS.md); the before-state is tagged
`study-01/v2-before-enhancements`. Independently loaded trials, mechanism variants,
growth/churn, fixed-reader contention, YB exceptions and deployment controls were
measured without migrating the original harness. Real network separation needs other hosts.

The implementation was developed in `.worktrees/study01-v3`, branch
`study-01/measurement-enhancements` (GPT-6 through Codex). D11–D17 and the separate
`-cmd experiment` / `-cmd report-enhancements` path are implemented. The v3 runner
creates a run tag, pins the image ID, takes the shared benchmark lock and keeps every
fresh-load trial. Containerized catalogue/control, scheduler overload and grouping tests
pass. All seven variants also passed the live PostgreSQL/YugabyteDB gates (14 checks per cell).
The implementation is already integrated into local `main` at `study-01/v3-integrated`;
the subsequent HIGH review and LOW completion receipt are integrated at
`study-01/v3-reviewed`.

**Completed:** `20260913T124917Z-v3`, 14 successful verification cells
(seven new variants on PostgreSQL and YugabyteDB single-node), run tag
`run/01-charity-tree/20260913T124917Z-v3`, harness `41a1490`, milestone
`study-01/v3-harness`. All 196 checks passed; both topologies were removed and the lock released.

**Completed:** `20260913T125342Z-v3`, 100 cells covering five independent trials
of D2/D3 and the mechanism variants, including inserts and a D2/D3 blended workload.
Run tag: `run/01-charity-tree/20260913T125342Z-v3`; digest `0bc5900a6d7a5d3a`.
All cells passed their correctness gates and reported no operation errors. D2/D3 median
12-query scores were 5,071.75 / 8,897.79, with 15.2% / 9.1% spread. The
[concise analysis](studies/01-charity-tree/reports/analyses/20260913T125342Z-v3--gpt-6--2026-09-13.md)
links a separate [mechanism discussion](studies/01-charity-tree/reports/discussions/d2-d3-mechanisms--gpt-6--2026-09-13.md).
The preceding gate is tagged `study-01/v3-verified`.

**Completed:** `20260913T172624Z-v3`, three independent trials
for growth/history/memory, fixed-reader contention, YB FK/cache exceptions, equal-total
budgets and local node-stop recovery. Run tag `run/01-charity-tree/20260913T172624Z-v3`.
Podman was restarted after the interactive pause; live VM resources were rechecked:
8 CPUs, 16,496,418,816 bytes and the same kernel. The preceding 100-cell run is tagged
`study-01/v3-mechanisms`. All 135 processes completed at 2026-09-13 19:55:53 UTC;
digest `bd95d304cb39e2de`. Six D9 negative-control cells failed their cache audits
(three trials on each YB topology); the other 129 cells reported no operation or
invariant errors. On resuming 2026-09-14, no database containers or benchmark lock
remained. The result set is tagged `study-01/v3-followup-results`.

All six constrained-memory cells completed their nine growth phases with no gate
failures. The matched follow-up below holds the data fixed while repeating both
memory configurations. Reporting now exposes growth read rates, arrival retries/counts and audits; containerized race
tests, vet and shell syntax checks passed on 2026-09-14. See the tooling report
`reports/20260914-v3-report-validation.md`. The completed matrix used its original
pinned image throughout.

**Completed:** matched memory controls
`20260914T104721Z-v3`, run tag `run/01-charity-tree/20260914T104721Z-v3`, source
`d1b18bf`. Twelve cells: D3/D6 × 256 MiB/3 GiB × three fresh-load trials, with
identical medium/history-multiplier=2 data and nine mutation phases. All twelve
initial gates and 108 post-mutation gates passed (1,680 individual checks), with
no read errors. Finished 2026-09-14 11:18:42 UTC; digest `9ead901cb74c4183`.
Final-phase donor-read medians at 256 MiB were D3 4,578.18/s and D6 16,034.40/s;
at 3 GiB they were D3 25,222.57/s and D6 15,415.07/s. The same dataset now
supports a memory-configuration-dependent reversal; it does not isolate the container
ceiling from PostgreSQL memory settings. This run removed its database and released
its lock. During the 2026-09-14 continuation, Study 03's historical
`devchecks/dc11-yb3-subset` held the machine lock; that wait ended and Study 03 v1 later
completed. Current runtime ownership always comes from the live lock and current context.

**Reporting change:** methodology 11b keeps the final signed analysis concise and moves
detailed design comparisons and analyst exchanges to signed `reports/discussions/`
companions, using `docs/templates/DISCUSSION.md`. Original published analyses stay intact.

**Final v3 analysis:** [enhancement conclusions](studies/01-charity-tree/reports/analyses/20260914-study01-v3--gpt-6--2026-09-14.md)
links four signed [discussion companions](studies/01-charity-tree/reports/discussions/README.md).
All 261 cells completed across four runs; six expected D9 cache-control failures remain
invalid for performance conclusions. The final milestone is `study-01/v3-enhancements`;
source/run tags and digests remain independent of the report/analysis commit. The branch
was integrated locally at `db77fe6`; the owner controls remote pushes. HIGH accepted
the existing evidence on 2026-09-15; EH-02 covers only the later review-document merge.
The [final artifact validation](studies/01-charity-tree/reports/20260914-v3-artifact-validation.md)
records provenance/link checks and the preserved editions from report regeneration.

**Earlier v1/v2 baseline:** survey and experiments C/D/E completed; D9's cache bug was
diagnosed and addressed by D10. Two signed analyses remain preserved for those inputs.
The v3 enhancements and current execution state are recorded above.

Baseline designs D1–D10 (the seven v3 additions are documented above and in the README):

| ID | Design | Isolates |
|---|---|---|
| D1 | normalized-minimal | the floor: 3NF, no secondary indexes |
| D2 | normalized-indexed | indexes alone (SQL identical to D1) |
| D3 | flattened-fk | denormalising the grandparent key onto the grandchild |
| D8 | flattened-nofk | the cost of enforcing referential integrity |
| D4 | rollup-trigger | consolidating aggregates upward, via triggers |
| D5 | rollup-app | same aggregates, maintained by the application |
| D6 | embedded-jsonb | folding the **whole** child table into the parent row |
| D9 | embedded-hybrid | folding a **bounded** slice (newest 20) into the parent, keeping the table |
| D10 | embedded-hybrid-locked | D9 with a concurrency-correct cache trigger — prices correctness itself |
| D7 | yb-child-colocated | physical data placement (YugabyteDB only) |

D6 and D9 are the two ends of the embedding question. D6 takes it to its conclusion and
runs into unbounded write amplification — appending one donation rewrites a document that
grows with the donor's entire history. D9 caps the array at 20, which makes the rewrite a
constant, and keeps the real donation table underneath as the source of truth.

Three topologies: PostgreSQL 1-node, YugabyteDB 1-node, YugabyteDB 3-node RF=3.
Twelve queries, three write operations, correctness-gated, plans captured per cell.

**Verified:** all designs return identical answers on both engines (12/12 checks).

**Correctness beyond the answer check.** Two audits run after the write benchmark:

- **rollup audit** (D4, D5) — recomputes every aggregate from the donation table and
  counts parent rows that disagree, plus the signed drift in the headline total. This is
  what makes "optimistic vs pessimistic" a real comparison rather than a throughput race:
  a strategy that is faster and silently loses increments has not won.
- **embedded-cache audit** (D9) — checks the bounded slice still holds exactly the newest
  donations, in the right order. Compares donation-id sequences rather than raw JSONB,
  because the loader and PostgreSQL render timestamps differently in JSON and a byte
  comparison would call a working cache broken.

**Tail latency.** p50/p90/p95/p99/max are recorded for every cell. p99.9 and p99.99 are
recorded only where the sample count supports them (10 000 and 100 000 respectively, so at
least ten samples sit beyond the quantile); below that the field is omitted rather than
estimated, because a p99.9 drawn from 363 samples is the maximum relabelled. Deep tails
are also biased low by coordinated omission — the harness is closed-loop — so they are
compared between designs and never quoted as SLO figures.

Deep tails were added after the survey matrix started, and raw samples are not persisted,
so the survey's own cells carry p99 and max but not p99.9. The follow-up stages
(D9 pass, experiment C, trials run) rebuild the image and do carry them — which is where
they matter most, since contention shows in the tail long before the median.

**Which tables get written.** Reads span all three tables, and *which* table answers a
question is itself the design variable. Writes originally all originated at `donation`,
with `person` and `charity` written only as side effects (trigger, app rollup, embedded
array). Two operations that originate at the **parent** have since been added, because
both are strongly design-discriminating and their absence hid costs that fall specifically
on the embedded and rollup designs:

| Op | Why it discriminates |
| --- | --- |
| `update_person` — a donor edits their email | Touches no indexed column and no child data; about as cheap an update as exists. But D6 has made the person row *contain the donation history*, so the cost depends on the row that trivial edit lands on. Genuinely uncertain: PostgreSQL does not rewrite an unchanged TOASTed column, so a heavy donor's out-of-line array survives, while a median donor's ~2 kB array sits inline and is re-versioned in full. |
| `delete_person` — erasure, donor and all their donations | Identical logical outcome everywhere, very different work. Normalised designs delete N children then the parent, and in D4 every one of those child deletes fires a trigger that recomputes the parent's MIN/MAX from the survivors. D6 deletes one row and the whole history goes with it. |

`delete_person` is excluded from the *mixed* workload on purpose: erasure removes donors
that concurrent inserts are still picking at random, and the resulting foreign-key
violations would be an artefact of the harness rather than a property of any design. It is
measured in the isolated write benchmark, where it has the id space to itself.

**Reads and writes are always reported together.** Every design here buys read speed
with something, so the generated report opens with a *trade-off at a glance* table putting
the purchase and the price in one row: a geometric-mean read score, insert throughput,
the number of indexes maintained on every `donation` insert, and storage. In the
controlled-pair tables, read differences below the error bar are suppressed but **write
rows are always shown** — a hidden price tag reads as no price at all.

**The report checks variation using D4/D5 reads.** These designs have identical read SQL
and serve as a control for cell-to-cell variation. Regenerating the survey from digest
`558e89403fd016cc` gives 1.69x maximum disagreement across 36 comparisons (1.04x median),
correcting the earlier 1.56x note. The report uses the maximum to suppress small pairwise
claims. This is a descriptive check, not a confidence interval or a universal significance
threshold; important comparisons still need repeated trials.

**Analysis provenance.** Generated reports carry an **inputs digest** — a content hash
over every result file — and index the signed analyses written about it. Several
analysts, including different AI models, are expected to interpret the same data; the
digest is what lets a later one tell whether the data has already been analysed. See
methodology rule 11 and `reports/analyses/TEMPLATE.md`.

### Runs so far (all `small` scale, environment `host-zenbook-ux5406sa`)

Run ids starting `20260913-` were produced on **2026-09-12 UTC** — the date in the name was
hand-set by mistake (LESSONS_LEARNED → "Never hand-set a date"). Manifests hold true times.

| Run id | What | State |
|---|---|---|
| `20260912-small` | Survey: all designs × PG 1-node, YB 1-node, YB 3-node (D9 added in a second pass) | ✅ 26 cells. Its `update`/`delete` numbers share one table across write phases — use `20260913-writes-isolated` for writes |
| `20260912T180830Z-concurrency` | Experiment C: D3/D4/D5 × 1–32 writers × strategies × isolation | ✅ 31 cells |
| `20260912-mixed` | Experiment D: readers:writers splits, dashboard mix | ✅ 5 designs — surfaced the D9 cache bug |
| `20260912-trials` | Repeated-trials writes | ⚠️ `delete` **invalid** (finite-pool bug); report in `reports/outdated/` |
| `20260912-parentwrites` | `update_person` / `delete_person` | ⚠️ `delete_person` **invalid** |
| `20260913-writes-isolated` | Every write op on a fresh load, 3 trials, PG + YB 3-node | ✅ replaces both invalid passes |
| `20260913-history-*`, `-top-1000charities`, `-contention-*`, `-portal-*`, `-size-*` | Experiment E: limitation regimes | ✅ D7 3-node unbounded needed a second pass (YB killed a session mid-reload) |
| `20260913-d9-cache-race` | Reproduction of D9 cache corruption, with diagnosed examples | ✅ 3–5 wrong caches per 30 s window |
| `20260913-cache-fix`, `-cache-fix-writes` | D9 vs D10 (fixed trigger): correctness and cost | ✅ D10 always consistent, 0.96–0.98x D9 writes |

**Analyses:** both remain current and are indexed by the 16 input runs via their digests.

- [Claude Opus 5 — first analysis](studies/01-charity-tree/reports/analyses/20260912-study01--claude-opus-5--2026-09-12.md): overall conclusions and the regimes where a usually slower design becomes useful.
- [GPT-6 — independent final analysis](studies/01-charity-tree/reports/analyses/20260912-study01--gpt-6--2026-09-12.md): reviewed the existing results without running new benchmarks; explains D2/D3 through indexes, SQL and actual plan work, and records narrower disagreements about embedding, foreign keys and application rollups.

**D2/D3 interpretation:** the 1.548x survey ratio is a geometric-mean score, not measured
mixed-workload throughput. D3 copies `charity_id`, adds two charity indexes and rewrites
the charity queries. Those six queries score 2.71x higher together; the six with unchanged
SQL score 0.884x. The saved charity-feed plans show fewer examined donations and donor
lookups. The sum plan still has 31,293 heap fetches despite its `Index Only Scan` label.
The mechanism is selective access and changed query work; the exact overall gain needs
repeated trials. Do not describe it as a universal benefit from a copied column.

**Writing a final analysis:** follow [methodology rule 11](docs/methodology.md) and the
[analysis template](studies/01-charity-tree/reports/analyses/TEMPLATE.md). Regenerate the
input reports, record every run and digest, read existing analyses, and write a separate
signed file with recommendations, linked measurements, mechanisms, limitations, proposed
follow-ups and explicit disagreements. Define aggregate scores, distinguish observations
from explanations, and identify untested behavior. Never edit another analyst's file.
Regenerate the reports after adding the analysis so their indexes include it, and update
this context and the process lessons. Study 01's final analysis does not cover Study 02.

### Experiment E — limitation regimes

Owner's instruction: where a design's limitation would be acceptable in some domain, test
that regime too, so the conclusion can say when a losing design is actually the right one.
The owner deferred to the AI which studies stay unbounded and whether the argument extends
to other levels of the tree. Decisions:

- **Unbounded history and 10 large charities remain the primary profile** — a charity's
  donation history genuinely grows. Regimes are reported beside it, never instead of it.
- **Bounded history (cap 20 per donor, donation volume held constant)** is run only for the
  designs whose limitation depends on per-donor history: D3 (baseline), D4 (delete trigger
  recompute), D6 (whole-array rewrite), D9 (cache = entire history at cap ≤ 20), and D7 on
  YugabyteDB (per-donor tablet). Not D1/D2/D5/D8 — their limitations don't depend on it.
- **The argument does extend to the top of the tree.** 1 000 small charities instead of 10
  large ones (people and donations held constant) targets D4/D5's single hot rollup row per
  charity and D6's charity-scoped unnesting. Asked both as a full cell and as a contention
  sweep at 8 and 32 writers.
- **Donor-portal workload** (reads that never cross donors) for D3/D6/D9 — the regime that
  removes D6's worst weakness regardless of history length.
- **Small dataset** (`tiny` vs `small`) for D1/D2 — does indexing matter at all when tiny?
- Regimes already covered elsewhere are recorded in the map in `run-regimes.sh` (D4 at low
  concurrency → experiment C; D5 insert-only ledger → experiment C; D8 → not a performance
  regime).

Harness support: `-max-per-person`, `-charities`, `-read-mix portal`, `-cmd report-compare`
(same designs under two regimes, flags where the per-question **winner** changes).

### Follow-up coverage and remaining questions

V3 completed the five-load D2/D3 mechanism comparisons, medium/long-history growth,
constrained-memory trials, fixed-reader contention, YB FK/cache checks, equal-total
budgets and local node-stop diagnostics. Scheduled arrivals now expose queueing,
rejections and acknowledged-write reconciliation. The matched memory run above is
the remaining active control; none of these runs certifies a production latency SLO.

- Independent hosts, real network separation, query-endpoint failure and a matching
  no-fault arrival condition remain external deployment work.
- Spread YB client connections across all nodes before interpreting a single-query-node
  bottleneck as a general limit of equal-resource sharding.
- Extend cache audits from ordered IDs to complete payloads under overlapping mutation
  streams; test reassignment invariants separately.
- Longer isolated mutations and covering-index reads under sustained churn would narrow
  the remaining variance and preparation uncertainty.
- Asynchronous rollups (queue / logical decoding) remain a separate unimplemented design.

### Study 02 — avoiding overbooking (band → event → ticket)

**Status:** designs, harness, runner, diagrams and README done (commit `774d258`, tag
`study-02/v1-harness`; runner fix `7570648`). Dev checks on `tiny` in `results/devchecks/`.

| Run id | What | State |
|---|---|---|
| `20260913T021010Z` | first `small` matrix launch | ❌ **INVALID** — aborted at the first cell; provenance not captured (see its INVALID.md) |
| `20260913T021206Z` | `small` matrix: 14 designs × PG 1-node, YB 1-node, YB 3-node; tag `run/02-ticket-booking/20260913T021206Z`; digest `71bcee71d725d43d` | ✅ 40 of 42 cells; C2 failed on both YB topologies (YSQL lease expiry, diagnosed) |

**Analysis:** `studies/02-ticket-booking/reports/analyses/20260913T021206Z--claude-opus-5--2026-09-13.md`
— first signed analysis (by the model that built the study; conflict of interest stated).
Headline: no correct design oversold and both controls did on every engine; for a hot drop,
pre-created rows taken by CAS or SKIP LOCKED sell 4–8x faster than any single counter row; a
counter on the event row also blocks unrelated edits; C5 under-sells after refunds; C2 is
unusable for a hot drop. Main doubts: DB CPU quota throttling everywhere, no retry backoff
(likely understates R2), one trial, 60 s timeouts truncate the 100k tier, single YB query node.

The owner's question: events of 10 / 100 / 1 000 / 10 000 / 100 000 seats, no overbooking;
compare pre-created tickets against tickets created on booking, other strategies, and extra
tables such as reservations. Fourteen designs in four families:

| Family | Designs |
|---|---|
| pre-created tickets | P1 lock-first (FOR UPDATE) · P2 skip-locked · P3 CAS · P4 skip-locked + counter |
| created on booking | **C1 count at RC — negative control** · C2 same SQL at SERIALIZABLE · C3 lock event + count · C4 guarded counter · C5 unique seat + CHECK |
| extra table | R1 inventory row · R2 inventory buckets (sharded counter) · R3 pre-created seat-slot pool |
| reservation | **H0 hold, naive confirm — negative control** · H1 hold, confirm checked by DB clock |

Experiments per cell: correctness gate (5 reads + audit on load) → plans for reads **and**
writes → reads (availability per size tier) → isolated writes on fresh loads (spread
booking, cancel, publish per tier) → **sell-out race** (32 buyers per unsold event, per tier,
organiser editing the event) → **churn race** (10% cancel; measures under-booking) →
**holds** (H only: basket, abandonment, expiry, sweeper, late payments). An overbooking
audit (over capacity, duplicate seats, derived-state drift, ledger reconciliation) runs
after every writing phase.

**Dev-check observations (tiny, not results, not for conclusions):**

- The negative control C1 overbooked every race event on PostgreSQL; the audit caught all.
- C5 under-booked in the churn race on every tier, as its design predicts (max+1 never refills refund holes).
- C2 on the 10 000-seat event: 15 attempts per sale on PostgreSQL; on YugabyteDB the race
  hung until the deadline was made to bound retries (LESSONS_LEARNED).
- YugabyteDB 1-node sells one to two orders of magnitude fewer seats/s than PostgreSQL on this
  laptop; hot-row designs (C4, P4, R1, H) around 15–20/s, R2 buckets ~87/s. Reads show a shared
  ~65 ms p99. **Measured:** the client is never CPU-throttled; the YugabyteDB container is
  throttled in 40–65% of CFS periods in every cell (up to 720 s per cell). YugabyteDB numbers
  here are quota-bound — this also applies to study 01's YugabyteDB cells.
- C2 on YugabyteDB 1-node fails reproducibly (dev check and `small` matrix) with "the database
  system is shutting down". **Cause established:** SERIALIZABLE lock storm on a CPU-throttled
  node → tserver RPCs to master stall ~29 s → YSQL lease expires → tserver kills all SQL
  sessions. See `results/20260913T021206Z/yb-single/logs/c2_count_serializable.DIAGNOSIS.md`.
- H0's control fired on YugabyteDB too. Holds at a 250 ms TTL were degenerate there (96% of
  holds expired before payment); YugabyteDB topologies now run holds at `-hold-time-scale 10`.
- YugabyteDB 2025.2.6 runs READ COMMITTED as real RC even without
  `yb_enable_read_committed_isolation` (checked); the flag is set to pin it.

**Harness:** `studies/02-ticket-booking/harness/` on `platform/`. Flags worth knowing:
`-phases`, `-race-trials N` (fresh load per trial; report shows median and spread),
`-race-tier-budget`, `-race-timeout`, `-hold-*`, `-tiers`.

**Open for study 02** (from the analysis's "what I would measure next"): retry-backoff variants
(expect R2 to overtake R1); repeated race trials on PostgreSQL and YugabyteDB 3-node; a buyers
sweep (4–128); a single-statement SKIP LOCKED variant to separate round trips from locking in
P3's lead; YugabyteDB with connections spread over all nodes; a larger per-node CPU budget; an
asynchronous allocator (virtual waiting room); a second analysis by a different model.

**Terminology review (planned, 2026-09-13):** study 02 is **general admission** — `seat_no` is
an admission number, and H0/H1 hold a unit of capacity, not a place. The review found no
test that is wrong under that reading. Its README gets a terminology section (no SQL,
harness, report or analysis change, no re-run), tag `study-02/v1.2-terminology`: step 2 of
study 03's handoff.

### Study 03 — reserved seating: choose seats, keep them 40 minutes (venue → event → seat)

**Status (2026-09-16, 01:30 UTC): complete — measured, analysed and signed** (tag
`study-03/v1-analysis`). Every escalation is decided and executed: ER-01 (AM-01), ER-02 (AM-02),
ER-03 and ER-04 (AM-03). The signed analysis is
[`20260915T002411Z--claude-opus-5--2026-09-16`](studies/03-reserved-seating/reports/analyses/20260915T002411Z--claude-opus-5--2026-09-16.md).
Nothing is running; the benchmark lock is free.

**What study 03 concluded** (details in the analysis):

- On YugabyteDB every correct design refused holds that were valid, and every such refusal was
  transient — the same statement moments later matched all the seats. PostgreSQL showed none.
- S1r, which retries a short confirmation once, recorded zero refusals; all 140 retries sold, at
  a cost inside the trial spread. It is the recommended defence on YugabyteDB.
- Publishing, not selling, separates the layouts: a 100 000-seat event costs 541 ms (PostgreSQL)
  and 3 216 ms (YugabyteDB) with per-seat rows against 0.89 ms and 1.56 ms for claim-on-hold.
- Avoid a JSONB section document for hot selling (28 976–33 712 compare-and-set retries per race
  tier on the cluster), and SERIALIZABLE on YugabyteDB (0.0–0.1 seats/s, every event timed out).

**Open for a future session:** a second analyst's view is invited, especially on L3, whose 4–8×
slowdown under contention this analysis reports without explaining. The analysis lists five next
experiments.

| Run | What | Result | Inputs digest |
|---|---|---|---|
| `20260915T002411Z` ([report](studies/03-reserved-seating/reports/20260915T002411Z.md)) | `small` main matrix, 3 topologies, 14 designs, 17 h 6 min | 41 cells, 4 failed (S4 and E2 on both YugabyteDB topologies, each diagnosed beside its logs); 15/15 controls fired | `a56ce92ce38b8204` |
| `20260915T173255Z` ([report](studies/03-reserved-seating/reports/20260915T173255Z.md)) | repeated race, 3 trials, 1 000- and 10 000-seat tiers, pg-single + yb-cluster3, 3 h 33 min | 27 cells, 0 failed; both controls fired | `29296b1fe8dea2e3` |
| `20260915T232736Z` ([report](studies/03-reserved-seating/reports/20260915T232736Z.md)) | repair (ER-03): E2 and S4 lifecycle on both YugabyteDB topologies, 40 min | 4 cells, 0 failed, no violation | `2ccece48793fe693` |
| `20260916T000706Z` ([report](studies/03-reserved-seating/reports/20260916T000706Z.md)) | repair (ER-04): L2 race on both YugabyteDB topologies, 11 min | 2 cells, 0 failed | `903de88de503ded2` |

The two repair runs come from a later commit than the matrix (AM-03's harness); the analysis states
which numbers come from which run.

Facts for the analysis:
- **PostgreSQL:** no invariant violation and no early rejection in any correct design.
- **YugabyteDB, transient refusals:** every correct design with a guarded statement recorded
  transient refusals, most on three nodes.
- **S1r:** 0 early rejections anywhere. Every confirmation it retried after a short match sold:
  35 in the matrix, 105 in the repeated race.
- **ER-03 (decided and executed):** the four failed cells had lost their **whole** lifecycle phase to
  an instrumentation query or a reload `ANALYZE` timing out, not to a design statement. AM-03 made
  both tolerant and re-ran those cells: all four completed, E2 with no violation on either topology,
  S4 collapsing at 0.0–0.1 confirmed seats/s — measured instead of missing.
- **ER-04 (decided and executed):** L2's YugabyteDB early rejections could not be classified, and the
  report read them as design failures. AM-03 gave L2 its own diagnostic — re-read the section
  document once on a refusal — and re-ran its race. On yb-cluster3 **all 17 early rejections were
  transient**: the document re-read showed the hold valid. So the ER-01 behaviour is not limited to
  statements filtering on columns the hold just wrote; a plain read of a recently committed row can
  miss it too. Reports now print "transient: not recorded" where no diagnostic ran (K0 always, L2
  before this change) instead of an unmeasured zero.
- **Provenance fix:** a run's inputs digest is a hash over its result bytes, and Windows line-ending
  conversion on checkout changed them (the matrix hashed two ways with identical measurements).
  Studies 02 and 03 are now pinned to `text eol=lf`, matching study 01's policy; every report
  reproduces its run's original digest.

How the plan got here: ER-02 was decided by AM-02 (tag `study-03/v0.2-handoff-amendment-02`). A
`small` calibration (dc15) replaced the duration extrapolation, and rule 1 dropped the 100 000-seat
race tier on YugabyteDB, because it sold 2–3% before timing out. The repeated race ran from study 03's
own worktree (`.worktrees/study03-measurement`, branch `study-03/measurement`), merged into `main` at
step 11. *Deviation, stated plainly:* between the task-worktree rule (`b290958`) and the end of the
main matrix, study 03 worked directly in the checkout of `main`. The
specification is the [Execution Handoff](studies/03-reserved-seating/HANDOFF.md) (tag `study-03/v0-handoff`); step-by-step
state, commits and mapped decisions are in [`PROGRESS.md`](studies/03-reserved-seating/PROGRESS.md);
unmapped decisions in [`ESCALATIONS.md`](studies/03-reserved-seating/ESCALATIONS.md).

- Done: study 02 terminology (tag `study-02/v1.2-terminology`); engine probe (every assumption
  held); platform error class for NOWAIT (tag `repo/platform-lock-not-available`); 13-design SQL
  catalogue; Go harness with 17 unit tests; tiny dev checks on PostgreSQL (all gates pass, no
  violation in any correct design, all four controls fire after one harness fix) and YugabyteDB
  1-node and 3-node (gates pass, controls fire).
- **ER-01 (decided):** on YugabyteDB, designs whose confirmation is a guarded statement
  occasionally refuse a valid hold; the identical statement re-issued in the same transaction
  matches every seat (a *transient refusal*). Not caused by expression pushdown, wait queues or
  query-layer retries (`diagnose-er01.sh`, `results/devchecks/er01-*`); never on PostgreSQL.
  Decision: measured and reported as its own class (still an INV-3 violation), plus design S1r
  (retry a short confirmation once) to measure the defence. AM-01 in HANDOFF.md §16.
- **AM-01 dev checks (dc12–dc14):**
  - PostgreSQL showed no early rejection.
  - On YugabyteDB, every early rejection in a correct design was transient: 21 on one node, 137 on
    three.
  - S1r showed none; every one of its 43 retries after a short match sold.
- **ER-02 (decided):** the projected duration exceeded §10.4. AM-02 adds a `small` calibration run
  (dc15) and fixed rules:
  - a 2-minute YugabyteDB race tier budget if needed;
  - the YugabyteDB 100 000-seat race tier dropped if it cannot finish;
  - the repeated race limited to the 1 000- and 10 000-seat tiers;
  - a new escalation only above 30 h (main) or 12 h (repeated race).

**The owner's workflow** (standard notation from 2026-09-14):

- A **higher** model/effort does the planning that needs thinking power, written as the
  "Execution Handoff". It also decides every "Escalation Required" (`ER-NN`) and writes the
  result analysis.
- A **lower** model/effort executes the handoff. It records any decision the handoff does not
  map as "Escalation Required", then continues with unblocked work.
- Every iteration ends by naming the next step's model: higher or lower.

Study 03 mapping: higher = Claude Opus 5, setting `ultracode` (study 03's earlier documents say
"ultra"); lower = Claude Opus 5, setting `high` (they say "high").

**The question:** buyers choose specific seats and keep them for 40 minutes while they pay,
and must never find a held seat gone. With marked seats, a unique key already prevents a
seat being sold twice. The hard parts become:

- a hold that excludes everyone else until it expires, and then frees the seat;
- conflicts over the *same* seat (SKIP LOCKED does not apply);
- all-or-nothing multi-seat blocks;
- the seat map as the hot read;
- expiry racing payment;
- whose clock decides.

**Plan in brief:**

- **13 designs:**
  - S0–S4, arbitration: conditional update, lock, NOWAIT, SERIALIZABLE; S0 is a check-then-act control;
  - E0–E2, expiry: lazy, sweeper, cart; E0 is an application-clock control;
  - K0–K1, checkout: K0 is a naive-confirmation control, K1 a payment window;
  - L1–L3, layout: claim rows with a unique arbiter, section document, section-sharded (L3 on YugabyteDB only).
- **Invariants INV-1..INV-8:** no double sale, no theft, honored hold, no leaked seats,
  all-or-nothing, no late sale, ledger reconciliation, no under-selling.
- **Experiments:** a hot-drop race with seat choice and deferred confirmers on real
  40-minute holds, and a compressed-time lifecycle with a sweeper outage.
- **Connections** are spread over all three YugabyteDB nodes.

### Study 04 — configuration portal (product → installed product → configuration entry)

**Status (2026-09-21):** measured and analysed on PostgreSQL single-node. Tagged
`study-04/v0-handoff`, `study-04/v0.1-handoff-amendment-01`, `study-04/v1-harness`,
`study-04/v1-measured`, `study-04/v1-analysis`. Nothing is running; the benchmark lock is free. One
agent, DeepSeek HIGH (`deepseek-flash`; effort and tool identity not exposed by the session),
planned, implemented, measured, validated and analysed this study — there is no LOW executor and no
model switch, which the analysis states as a limitation.

**Measured run `20260921T1215Z-small`** (tag `run/04-configuration-portal/20260921T1215Z-small`,
inputs digest `ab0f6ef5e5500775`): ten designs at scale `small` on `pg-single`, ten cells, none
failed, both negative controls fired. Headlines, all from
[the signed analysis](studies/04-configuration-portal/reports/analyses/20260921T1215Z-small--deepseek-flash--2026-09-21.md):
a lock-before-update delivered 886 ops/s against a version check's 503 at identical correctness,
with 545 retries and 10 560 conflicts on the optimistic side; the unchecked read-modify-write
acknowledged 3 416 increments and left the counter at 215 (**3 201 lost updates, 93.7 %**); a parent
rollup makes the portal overview ~15x faster while the trigger variant costs 6.9x on
whole-configuration replacement; the same logical configuration is 10.5x smaller as one document.
**The run's principal weakness**: reads whose SQL is identical across the eight row designs spread
by 65 %, so no read difference below ~1.7x in this digest is attributable. Clause: ten of the
eighteen designs are implemented; cadence, churn, deployment controls and repeated trials are
mapped and not run, and the analysis records them as coverage gaps.

The portal is used by other products: all configuration creation, modification, publication and
retrieval goes through it and its database. Configuration belongs to an **installed product** (one
deployment of a **product definition**, in an **environment**, optionally per **business unit**),
never to the product definition itself; several installations of one definition may share an
environment.

- Specification: [`studies/04-configuration-portal/HANDOFF.md`](studies/04-configuration-portal/HANDOFF.md)
  (`EH-04` revision 1). Escalations: `ESCALATIONS.md`; step log: `PROGRESS.md`.
- **18 designs**: `n0`–`n4` normalized (index control, reference, rolldown, trigger rollup, app
  rollup); `d1`–`d4` document (one-to-one row, embedded on the parent, section-sharded, JSON path
  update); `h1` normalized rows plus a materialized read representation; `s1` immutable snapshots
  with an atomic revision pointer, `s2` append-only history plus materialized current state;
  `y1`/`y2` colocated vs non-colocated (YugabyteDB only); `c1`/`c2` optimistic vs pessimistic
  concurrency; `x1`/`x2` the two negative controls (lost update, rollup drift).
- **INV-1…INV-13** verified in Go from the generated dataset; correctness gates timing.
- **Cardinality tiers 1, 10, 30, 60, 120, 500** — 60 mandatory — plus a bounded skewed
  distribution. The controlled comparisons hold total entries, then total serialized bytes,
  approximately constant while entries per installed product changes; a **fixed fleet** is a
  separately labelled third scenario.
- **Cadence is a rate, not a wait:** λ = installed products / period, i.e. 0.0058–500 updates/s for
  a 500-product fleet. Jittered arrivals and synchronized bursts; every calculated capacity is
  labelled a projection and never presented as a measured temporal result.
- This session's measurement scope, agreed with the owner, is **prove the pipeline**: phases A–C
  plus a reduced phase D (a small `pg-single` matrix). Full-breadth D and phases E–J are mapped in
  the handoff for a later session.

**04-ER-01 (decided):** the preflight assumed a native WSL `podman`; there is none. The documented
engine is reachable only through the Windows `podman.exe`, which from WSL already reports the same
server, images, volumes and benchmark-lock state. AM-01 decides that WSL drives that engine through
`podman.exe` via an additive resolver in `infra/lib.sh`, and that a WSL-local podman is never
initialised — a second engine would carry no shared benchmark lock, which is the failure the lock
exists to prevent. Recorded in `ESCALATIONS.md`.

**Integration (2026-09-21).** This task's changes merge into local `main` as a **fast-forward**:
`main` was at `df13f2a` when the task branch was cut from it and did not move, so no merge commit
and no reconciliation of a concurrent agent's edits were required. The main checkout's working
tree carried 1970 CRLF-vs-LF phantom modifications on arrival; `git diff --ignore-cr-at-eol --quiet`
exited 0 over all of them (every differing file differed only by CR at end of line, nothing staged,
nothing untracked), so the working tree was normalized to its own LF index and `git status` is now
clean there. No index entry, no commit and no tag was changed by that normalization. The task
commits `bd2b825`, `8f02c53`, `2e95990`, `5e15abb`, `73eaf14` and `40f3294` are reachable from
`main`; the integrated state is tagged `study-04/v1-integrated`. Nothing was pushed or pulled.

**Environment:** `host-zenbook-ux5406sa` re-verified live on 2026-09-21 — 8 CPUs,
16 496 422 912 bytes, kernel `6.6.87.2-microsoft-standard-WSL2`, host ASUS Zenbook S 14 UX5406SA.
`podman machine inspect` still shows its stale `init` value of 4 CPUs / 2048 MiB; the live guest is
the environment page's 8 CPU / ≈15.36 GiB, and no machine was created, resized or started.

### Study 05 — external cache throughput and consistency (donor portal)

**Question.** For a database-backed donor portal, what does an external cache buy in throughput and latency
over the same-run no-cache baseline, what does strict freshness cost, how often and how badly is a relaxed
cache wrong, and what can a *legacy* database model that may not be modified do compared with an owned one
that may carry a version token and an outbox?

**Where.** `studies/05-cache-consistency/` — EH-05 rev 1 (`HANDOFF.md`), `PROGRESS.md`, `ESCALATIONS.md`
(two decided items, `05-ER-01`/`05-ER-02`), 21 core scenarios plus 3 database-layout reference cells and 2
controls, SQL in `sql/{legacy,owned,reference}/`, the Go harness in `harness/`, runners `run-study.sh` and
`probe-cache.sh`, the cache topology script `infra/redis.sh`.

**Scenario ids are a product of dimensions**, not a family of copies:
`<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`, e.g.
`owned-opt-redis-aside-strict-coord` or `legacy-na-memory-through-relaxed-ext20`.

**The run.** `20260921T-survey`: scale `small` (800 donors), PostgreSQL 17.11 `pg-single` plus pinned Redis
7.4.11, resource framing `db-only` (database 2 CPUs/3 GiB), seed 42, fault seed 4242, code `567778f`,
inputs digest **`8ff86dc9e5228e86`**, report `reports/20260921T-survey.md`.

**Result, in one line.** Caching added 2.0×–4.4× over the same-run no-cache baseline (3 914 → 15 018 warm
reads/s for the best cell); strict freshness cost nothing measurable in read throughput in a clean pair;
a single-instance relaxed cache was wrong for 0.04 %–0.3 % of reads while writes were in flight and for
none in the warm read-only phase; a **three-instance shared-Redis cache was wrong for ~86 % of reads**, the
run's headline; and **nine of 27 cells failed their correctness rules** and support no conclusion.

**Correctness regime (the deliverable that matters most).** Content-hash freshness judged by an independent
Go oracle and an operation ledger, not by re-reading the database on every hit; every payload-producing read
in one repeatable-read transaction; every fill and refresh under a per-key lease; an exact byte-bounded LRU
in process and `allkeys-lru` with recorded `maxmemory`/`maxmemory-samples` in Redis; a 300 s hard TTL that is
never shortened plus probabilistic early expiry verified by fake-clock unit tests at 0/75/150/225/300 s; six
deterministic fault phases; a stale-read negative control that fires; and a dirty-cache-write audit that
injects a record from no committed state and requires rejection.

**The finding that changes how caching should be built here:** "publish only after the commit" does *not*
make cache-aside safe. A reader that begins its fill before a writer's invalidation holds a committed but
**superseded** state and can republish it. The study closes that with a cache-side per-key **invalidation
fence** (atomic in both backends, captured before the snapshot, checked on publish) that a strict writer
advances **before and after** its commit, and by moving the ledger's freshness requirement to the
**acknowledgement**. `publishes_refused_by_fence` (61–85 per cell) is the evidence that the race is real.
See `reports/discussions/20260921-cache-fences-and-ledger.md`.

**Do not reuse these numbers without reading the analysis' weakness list.** Single trials; no churn phase
(so the TTL was never crossed in a measured run); no equal-total framing; no `medium` scale; no YugabyteDB
or three-node cell, so the colocation requirement is an explicit gap; one laptop, shared cores, no real
network; and an unresolved ledger limitation for concurrently conflicting writes of the same key.

**One provenance caveat, stated up front:** the three reference cells were re-run after a harness fix, so those
three were produced by a later code state than the other 24 cells while the manifest records one commit for the
run. The fix changes `hasCache()` from a comparison against `"none"` to an explicit backend test, which is a
no-op for every cell whose backend is set — i.e. all 24 others. The analysis repeats this in its weakness list.

**Post-survey repair (same session, after the analysis).** The strict `through` failures were traced
to a *harness* defect, not a design one: a mutation selected a child row from the ledger and then
acted on it by id only, so a concurrent writer could move or delete it in between and the database
and the ledger diverged. Every child-row statement now enforces the owner the harness observed (all
five schema directories), a statement matching no row is an acknowledged no-op that leaves the
ledger alone and is counted (`write_noops_refused_by_owner`), and corrections are relative in every
schema including the reference ones. AM-01's required assertion is implemented — after every writing
phase the harness compares all 800 donors' database content against the ledger's requirement and
fails the cell if the requirement is ahead: **0 mismatches in warm, mixed, hotspot and stampede; 1
donor after `instances`**. A second measurement defect was closed in the same pass: violations
recorded by the instances-phase arms were invisible to the cell's acceptance check, so one strict
cell had passed with 11 hidden stale reads. With both fixed, two reference cells
(`ref-normalized-indexed`, `ref-embedded-locked`) are green again, and the strict `through` residual
is now a genuine cache-design finding (7–32 stale reads, every assertion passing, 260–625 refused
publications per cell) rather than an ambiguity. Evidence: `results/verify-ledger/`.

**Corrections after the repair (runs `verify-ledger`, `verify-rollup`).** Two results changed the study's
reading, and both are in the analysis (§6c/§6d) and AM-03. First, the legacy in-process *aside* strict
cell — the one the first analysis leaned on as proof that the fence closed the refill race — is **not
reliably clean**: re-run it recorded 10 stale-after-ack reads per ~53 000 hits in the mixed phase, with
every ledger assertion passing and zero impossible values. So strict freshness as implemented is **not
achieved reliably in any strict arm** (aside or through, memory or Redis, legacy or owned); that is now
the study's single open scientific question. Second, `ref-rollup-trigger`'s 114 "impossible" values were
a harness artifact of the weaker impossible-value rule; with the source-based rule restored the cell
fails for the right reason — its own `a_rollup_drift` audit reports 3 mismatches and the stored total is
2 173 cents above the ledger — so the trigger-maintained rollup is **not covered** by that reference as
it stands. The two green reference cells (`ref-normalized-indexed`, `ref-embedded-locked`) and both
controls still behave.

**Residual classified (AM-04).** The few stale reads that survive the fence are, in the legacy
in-process aside cell, `fill` and `bypass` samples — **database reads**, not cache hits — each exactly
one state behind the requirement and resolving before the phase-end assertion runs. The next step is
pinned down: record `(key, ack seq, commit time)` per mutation, compare a stale database read's snapshot
time with the acknowledgement time, and add "a fresh connection must already see the state the
requirement moved to, immediately after each ack" as an assertion.

**Resolution (run `20260921T-survey3`, digest `03fac739d0835768`, clean commit `48fe1e3`, tag
`run/05-cache-consistency/20260921T-survey3`).** The remaining "correctness failures" were four defects in
the harness, three of which were *manufacturing* findings rather than detecting them: a mutation applied by
id only (fixed with owner-guarded statements and a distinct no-op outcome), the ledger's history order not
following the database's commit order (proved by a **no-cache** design showing it, fixed by serialising
writes that touch one key), the acknowledgement checker capturing its requirement *after* its read, and a
control judged by the wrong evidence. With all four fixed the full matrix is green: **27 of 27 cells pass
every gate**, both controls fire, 6 ledger assertions per cell with zero mismatches, 310 acknowledgement
verifications per cell with zero violations, zero impossible values, zero stale-after-ack reads.

**Final measured results (single run, `small`).** Caching buys 2.2×–3.5× in warm cacheable throughput and
an order of magnitude in p99 (15.8–27.1 ms → 0.8–2.0 ms). **Strict freshness cost nothing measurable in read
throughput** in any of the eight controlled pairs (all differences inside the ~20 % noise floor). Every
single-instance relaxed cell recorded **zero** wrong reads across ~1.2 million measured reads; the only
wrong reads in the whole matrix are the **three-instance shared-Redis through** cells at **87.2 % and
87.8 %**, confirmed in two runs and in both database models. The lease gives exactly one database load per
key under a 16-reader stampede (8 loads, 0 duplicate fills). Under 20 % external writers the legacy cache's
value collapses to the no-cache baseline (7 664 vs 7 794), and its strict sibling only reaches 8 269 by
reading authoritatively.

**Coverage gaps, unchanged and explicit:** no sustained-churn phase (the 300 s TTL was never crossed in a
measured run), no equal-total framing, no `medium` scale, no repeated trials, no YugabyteDB cell, no
three-node cluster, no colocation evidence, no open-loop SLO work.

**Status: measured, analysed and integrated into `main`.** Findings, coverage gaps and the five leads for a
second analyst are in `reports/analyses/20260921-cache-consistency-allgreen.md`, which supersedes the
survey analysis (now in `reports/outdated/`).

### Study 05 — final integration receipt

Branch `study-05/cache-consistency` (worktree `.worktrees/study05-cache-consistency`) was merged into local
`main` by fast-forward at each milestone; the final state is `main` = the commit carrying this paragraph, with
every task commit reachable — `fa1a57a`, `92f8d0b`, `91f2ce0`, `567778f`, `cac02f5`, `a1f09a4`, `72bb543`,
`4a747f5`, `43a2990`, `ead9b18`, `0682f5e`, `f9d9f57`, `48fe1e3` (the code the reported run was produced
from) and `372f437` (final analysis, lessons, context). Annotated tags: `study-05/v0-handoff`,
`study-05/v0.1-handoff-amendment-01`, `study-05/v1-harness`, `study-05/v1-measured`, `study-05/v1-analysis`,
`study-05/v1-integrated`, `study-05/v1.1-diagnostics`, `study-05/v2-owner-guard`,
`study-05/v2.1-impossible-attribution`, `study-05/v2.2-residual-classified`, `study-05/v3-allgreen`,
`study-05/v3-analysis`, and the reported runs `run/05-cache-consistency/20260921T-survey` and
`run/05-cache-consistency/20260921T-survey3`. No tag was moved, deleted or reused; nothing was pushed or
pulled and no remote was touched. The benchmark lock was released and the study's containers (`pg-single`,
`ads-redis`) were stopped and removed by the runner after every pass. Post-integration checks on `main`:
`gofmt -l .` clean; `go vet ./...` clean; `go test ./...` green for the study, in the pinned Go container
(the platform's timing-sensitive `core/measure` test remains the one caveat recorded above). The nine failed cells, the
rollup reference drift, the unrun churn/equal-total/topology arms and the three-instance confirmation are
open and are listed in the analysis (sections 3, 6 and 7). The next session should take those, not re-run
this matrix unchanged.

### Study 05 — integration receipt

The task worktree `.worktrees/study05-cache-consistency` on branch `study-05/cache-consistency` was merged
into local `main` by fast-forward at task completion. The task commits are `fa1a57a` (handoff),
`92f8d0b` (SQL catalogue), `91f2ce0` (harness), `fd599da`, `e5d1e1a`-class dev-check fixes, `567778f`
(the code the reported run was produced from), `cac02f5` (measured survey, signed analysis, mechanism
companion, context and lessons) and the receipt commit that carries this paragraph. All are reachable from
`main`; the integrated state is tagged `study-05/v1-integrated`. The producing run is tagged
`run/05-cache-consistency/20260921T-survey`. The benchmark lock was released and the study's containers
(`pg-single`, `ads-redis`) were stopped and removed by the runner; nothing was pushed or pulled, and no tag
was moved, deleted or reused.

**Post-integration checks re-run on the combined state**, inside the pinned `golang:1.26-bookworm`
container, against the `main` checkout mounted read-only: `gofmt -l .` clean; `go vet ./...` clean for the
study; `go test ./...` green for the study (14 tests, including the fake-clock probabilistic-expiry tests at
0/75/150/225/300 s, the exact-LRU byte bound, the lease rule, the fence rule and the oracle's
freshness/acknowledgement semantics). The platform is unchanged by this task, so its tests are unchanged
**with one caveat worth recording**: `go test ./...` across the platform failed once on
`adsplatform/core/measure` and the same package passed when run alone immediately afterwards. The timing
tests in that package are sensitive to parallel package execution and to load on this machine — the same
hazard `LESSONS_LEARNED.md` records for measurement generally. It is study 04's code, it is not part of this
task, and it is reported here rather than fixed.

## Decisions taken, and why

| Decision | Reason |
|---|---|
| Plain `podman` CLI, not compose | No compose provider ships with Podman; requiring one would break the "install nothing" rule |
| Per-node resource budget held constant | Makes "1 node vs 3 nodes" mean "what do two more machines buy me" |
| `--cpus` quotas, never `--cpuset-cpus` | This host mixes P-cores and E-cores that the guest kernel cannot distinguish; pinning would silently bias results |
| SQL in files, embedded via `go:embed` | Readable as documentation, and the binary provably runs the SQL shipped beside it |
| Bulk rollup at load, triggers attached after | Loading with per-row triggers would measure the loader, not the design; trigger cost is measured by the *write* benchmark |
| Queries benchmarked in isolation | A blended number hides which query moved, and attribution is the point |
| Bounded long tail in generated data | An unbounded power law would make the embedded design a strawman |
| Donations generated in time order | Mirrors an append-only stream, so surrogate key and timestamp correlate as in production |

## Open questions worth a future study

- Extend the v3 equal-total-budget comparison to balanced query endpoints, larger data
  and independent physical hosts.
- Partitioning in PostgreSQL (`PARTITION BY HASH` + local indexes) as the non-distributed
  equivalent of D7's placement lever.
- Where the embedded design's crossover lies as array length grows — the write cost is
  O(history), so there is a document size at which it stops paying.
- Deferred / asynchronous rollups (queue, logical decoding) as a fourth point between
  D3's read cost and D4's write cost.
