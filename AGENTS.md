# Working in this repository

Instructions for **every AI agent** working here, whichever vendor or tool runs it (Claude
Code reads this file through `CLAUDE.md`; OpenAI Codex and other agents read `AGENTS.md`
directly). This is the single source of these rules: change them here, nowhere else.

Before starting, read [`CONTEXT.md`](CONTEXT.md) for current state (including work other
sessions have in progress) and [`docs/methodology.md`](docs/methodology.md) for the rules
that govern every study.

**Every task starts by creating or reusing its own worktree and ends by merging its
completed changes into local `main`.** This is the owner's standing authorization;
no separate merge permission is needed. Follow the task lifecycle below.

## What this project is

An experimental record of how architecture and data-modelling decisions affect
performance, **measured rather than asserted**. The output is reproducible runs and
written reports, not opinions. Several analysts — people and AI models from different
vendors — are expected to read the same data and may disagree; that is by design.

## Non-negotiables

These come from the project owner and override convenience:

1. **Everything runs in Podman containers.** Never instruct a reader to install a
   database, a language toolchain or a SQL client on their host. The benchmark client is
   containerised too. `podman compose` is acceptable in principle but no compose provider
   ships with Podman, so orchestration uses the plain CLI.
2. **Go** is the implementation language for tooling.
3. **Every test has a report.** When a report is superseded, move it to
   `reports/outdated/` — never delete it, never silently edit it.
4. **Hardware and environment are documented** and named by every result file.
   Results are comparable only within one environment.
5. **`CONTEXT.md` and `LESSONS_LEARNED.md` are kept current.** Update them as part of the
   work, not afterwards.
6. **Database studies compare single-node and 3-node clusters**, and report read and
   write throughput plus `EXPLAIN ANALYZE` as the floor, expanding with other tools where
   they help.
7. **The AI proposes the design combinations.** Denormalisation, replicating data between
   tables, embedding children into parents, index choices, ACID/isolation levels,
   optimistic vs pessimistic concurrency — these are expected to be generated and
   compared, not asked about.
8. **AI agents work in isolated worktrees, commit, merge completed tasks into `main`,
   and tag their own work. They never push or pull.** Remote synchronization is the owner's.
9. **One benchmark on the machine at a time.** Agents may write code and analyses in
   parallel; they may not measure in parallel.
10. **Never edit another analyst's file.** Not to shorten it, restructure it or correct it.
11. **Always assume another AI agent is working in this repository at the same time**, and
    work so that neither interferes with the other. Findings reach `main` by merging; the
    agent that merges later reconciles both sides so `main` reads as one project, not a
    patchwork of sessions ("Hard rule: assume a concurrent agent").

## Required comparisons for future studies

**Owner clarification, 2026-09-14:** these requirements apply to every new study and
new measurement plan for an existing study. Record their coverage in the study's
protocol before running; existing published results are not retroactively relabelled.

1. **Calculate and validate cluster and client sizing.** Document available host/VM
   resources, client and system reserves, node count, per-node CPU/memory and aggregate
   budget. Multi-node comparisons include a separately labelled equal-total-budget
   control alongside the per-node baseline. Calculate client readers, writers and
   connection-pool allowance explicitly, check that the generator can supply the
   intended demand, and measure per-node throttling and endpoint distribution. A
   single query endpoint must not silently stand in for balanced cluster access.
2. **Measure rollup, rolldown and embedding wherever applicable.** Compare derived
   answers with stored parent aggregates (rollup), copied parent information on children
   (rolldown), and embedded child information. Include read gains, mutation/maintenance
   costs, storage and correctness under the relevant workload; isolate decisions with
   controlled pairs instead of changing several mechanisms together.
3. **Every multi-node database study includes colocated and non-colocated data
   scenarios.** Hold logical data, operations, replication and resources constant;
   verify actual placement and measure its effect. Running containers on the same
   laptop is not evidence of data colocation. An unsupported placement mechanism is
   an explicit coverage gap, not a completed comparison.
4. **Concurrency studies include optimistic and pessimistic strategies at minimum.**
   Compare a conditional/version-checked approach with a lock-before-update approach
   under the same business invariant and comparable workload. Include contention,
   retries/aborts, latency and correctness; add other strategies when applicable.

List measured pairs, justified non-applicability and outstanding gaps in each protocol
and final analysis. Methodology sections 4 and 6a explain the sizing and control details.

## Model roles: cost-efficient Execution Handoff, Decision Log and Escalation

**Owner workflow, 2026-09-14, revised 2026-09-17 for cost efficiency:** HIGH and LOW are
roles chosen through the user's model/effort settings. Record the actual model/tool and
selected effort when known; never infer an unavailable setting or switch models silently.

### Session capability and queued work

**Owner workflow, 2026-09-22:** at the start of every new repository chat, ask whether
the session is **HIGH** or **LOW** unless the user already stated it. Record that answer
as the session capability; do not infer it from a model name. Capability and work role
are separate: a HIGH session may plan, review, analyse, integrate, or execute a worker
task, while a LOW session may claim only work whose minimum capability is LOW.

Durable requests, task briefs, attempts, events, escalations, results, and reviews live
under [`docs/ai-work/`](docs/ai-work/). Once the queue CLI is available, use its atomic
claim and guard commands before acting. During bootstrap, a task's canonical branch name
and registered worktree are its exclusive claim: create them from current local `main`,
and if either already exists, inspect and resume that task rather than starting a second
copy. `CONTEXT.md` is an overview, never the queue's state authority.

**The objective is to reserve HIGH for work that genuinely benefits from deeper
reasoning, and let LOW execute as much as possible autonomously.** The preferred pattern
is **HIGH planning → one long LOW execution → HIGH validation**, not frequent switching.
Minimize unnecessary HIGH usage, unnecessary model switches, repeated context, premature
escalation, and re-litigating decisions HIGH already made.

1. **HIGH plans.** Understand the objective, identify constraints, make the
   architectural/strategic/high-impact decisions, resolve significant ambiguities,
   define acceptance criteria, and divide the work into tasks LOW can safely execute.
   Do not perform routine implementation here when it can be delegated. Publish a
   compact **Execution Handoff** (self-contained, not a repeat of the full reasoning
   history): goal, current state where relevant, decisions already made, tasks to
   execute, constraints, acceptance criteria, known risks/assumptions, and the specific
   conditions that would genuinely need immediate escalation. `docs/templates/
   EXECUTION_HANDOFF.md` is the template when its structure fits; a shorter handoff that
   still contains those elements is fine. End with `NEXT MODEL: LOW`.
2. **LOW executes with substantial autonomy.** Default behavior is **Decide → Log →
   Continue**, not "encounter ambiguity → escalate." LOW does not reconsider a decision
   HIGH already made unless execution proves an assumption wrong or impossible to
   follow — that itself is reported (Decision Log or escalation, by its cost), not
   silently overridden. LOW may decide, without escalating, whenever a choice is local
   to the implementation, has a conventional or sensible answer, is reasonably
   reversible, does not change the intended architecture, does not materially change
   the requirements, and is unlikely to invalidate substantial downstream work. LOW does
   not change the scientific question, design, published acceptance criteria or signed
   conclusions to make a run pass — that is never a "local, reversible" decision.
3. **The Decision Log records meaningful decisions, not every choice.** Log a decision
   only when it could plausibly affect correctness, architecture, performance,
   maintainability, security, behavior, compatibility, acceptance criteria, or the final
   result — in the task's `PROGRESS.md` or a dedicated log, whichever the task already
   uses. Each entry: the decision, a short reason, confidence (High/Medium/Low), the
   alternative considered (when relevant), and which tasks it potentially affects. Keep
   entries short; this log is what HIGH reviews at validation instead of re-deriving
   everything LOW touched.
4. **Escalate only when continuing alone costs more than invoking HIGH now.** Use
   `ESCALATION REQUIRED` (`docs/templates/ESCALATION_REQUIRED.md`, logged in the task's
   `ESCALATIONS.md`) when a decision could invalidate substantial downstream work,
   would significantly alter the architecture, would materially change the
   requirements, carries significant security/correctness/data-loss risk, is difficult
   or expensive to reverse, contradicts a critical HIGH assumption, or genuinely cannot
   be resolved from the existing instructions. **Do not escalate merely because**
   several reasonable options exist, confidence is imperfect, a minor detail was not
   explicitly specified, LOW would simply prefer HIGH decide, or the choice can easily
   be reviewed later — decide, log, and continue instead. An escalation states: the
   decision needed, why it cannot safely be deferred, the relevant facts found, the
   options (and a recommendation if LOW has one), what work is affected, and what work
   can safely continue in the meantime — then `NEXT MODEL: HIGH`. Continue unrelated
   mapped work rather than stopping the whole task.
5. **HIGH resolves an escalation narrowly.** Review only as much context as the specific
   decision needs, not the whole project. Reply with a small **Execution Handoff
   Update**: the decision, changed instructions, affected tasks, any new constraint or
   acceptance criterion — not a restatement of what did not change. End with
   `NEXT MODEL: LOW`. LOW resumes from the original handoff plus this update plus its
   existing Decision Log, redoing only what the decision actually invalidated.
6. **HIGH validates by reviewing the Decision Log first, then the result.** For each
   logged LOW decision, judge only whether it is acceptable, negligible, needs
   adjustment, or materially degraded the result — do not order rework merely to
   substitute HIGH's preferred but equivalent style for an acceptable LOW choice.
   Trigger rework only for a decision or defect that materially affects correctness,
   requirements, architecture, important performance characteristics, security,
   maintainability, acceptance criteria or overall quality. A real problem gets a
   **targeted** corrective handoff (what's wrong, which tasks are affected, what must
   change, what stays untouched) — never a restart of the whole task. End with
   `NEXT MODEL: LOW` for the correction, or, when nothing meaningful remains, close with
   a brief conclusion (what was completed, important decisions, LOW decisions reviewed
   and accepted, outcome, known limitations, relevant future work) and `TASK COMPLETE`.
7. **Every response ends with exactly one routing line:** `NEXT MODEL: HIGH`,
   `NEXT MODEL: LOW`, or `TASK COMPLETE` — plus the concrete model/effort when the task
   has agreed one. If the user states which model/effort is currently active, say
   explicitly whether they need to switch before continuing.

A model-switch handoff is a checkpoint in the same unfinished task. Commit and tag that
checkpoint, update context/lessons, and state what remains; do not claim the task is
complete or the changes merged merely because the planning iteration ended. The normal
task-completion merge rule still applies once the task's acceptance conditions are met.
An expected wait for a benchmark lock stays with LOW; waiting alone is not a decision
that needs logging or escalation.

## Hard rule: commit and tag, never push

Studies are revisited: new designs, new questions, new analysts. Two reports of the same
study can reach different conclusions, and the **only** way to know why is to know which
code, SQL and data produced each. Without tags, that history is lost. So every agent:

1. **Commits at every meaningful step** — a design added, a harness change, a run's
   results, a report, an analysis, updates to `CONTEXT.md` / `LESSONS_LEARNED.md`. Never
   leave finished work uncommitted at the end of a session. Unless the owner says otherwise
   for a specific change, this rule applies.
2. **Tags, with annotated tags (`git tag -a`):**
   - `run/<study>/<run-id>` — the commit that produced a run. Start every run that will be
     reported from a committed tree with the runner's `--tag`; a dirty tree is not tagged.
   - `study-NN/vX-<label>` — a milestone of a study: designs or harness as first run
     (`v1-harness`), a report plus analysis set (`v1-analysis`), a new analyst's analysis, and
     — **when enhancing a study** — the state *before* the change (if not already tagged) and
     the state *after* it (`v2-<what changed>`). A later reader diffs the two tags to see
     exactly what changed between conclusions.
   - `repo/<label>` — repository-wide changes that belong to no single study (rules,
     infrastructure, the shared platform).
   Check `git tag` first: tag names are shared by every agent and every worktree.
3. **Never moves, deletes or reuses a tag, and never rewrites history** (no amend of
   published work, no rebase, no force). A corrected state gets a new tag.
4. **Never pushes, never pulls, never changes remotes.** Those are the owner's.
5. **Records the tag** in what it produces: results and reports carry the run tag
   automatically; analyses carry `repo_commit`; `CONTEXT.md` lists tags and runs.
6. **Identifies itself in the commit message** — model and tool, e.g. a
   `Co-Authored-By:` trailer or a closing line — so `git log` shows which agent did what.
7. **Commits only what it can account for.** Stage explicit paths, never `git add -A` in a
   shared tree. If a commit must contain another session's work, say whose in the message.

## Hard rule: task worktrees and completion merges

**Owner clarification, 2026-09-14:** isolation is required at the beginning of every
task, including when only one agent is working. Create a worktree and task branch from
current local `main`, or reuse this task's existing worktree after checking its branch,
status and ownership. Never reuse another active task's folder or develop directly in
the checkout of `main`.

Several agents can work at the same time. They **must not share one working folder**: a
folder has one checked-out branch, and two agents editing the same files — or one
switching branches — change each other's work mid-edit. Use **git worktrees**: one folder
per agent, one branch per folder, one shared repository (commits and tags are visible to
all worktrees at once, with nothing to push between them).

```bash
git worktree add ../ads-<study>-<topic> -b <study>/<topic>
```

- **Branch names:** `<study>/<topic>` (e.g. `study-02/retry-backoff`) or `repo/<topic>`.
- **At task completion, merge all of this task's changes into local `main`.** Commit
  explicit owned paths, integrate the latest local `main` into the task branch, resolve
  conflicts in the task worktree, and run checks appropriate to the combined changes.
  Preserve every session's context/lessons entries and other analysts' files. Then
  fast-forward `main` to the integrated task branch; if `main` advanced, integrate the
  new commits and check again. Do not rewrite history or leave a successfully completed
  task only on its branch. A task-specific user instruction can override this default.
- **Serialize the final update of `main`.** Check its checkout for another task's
  edits or active work. Never stash, discard or commit someone else's unfinished files.
  Do not change files that an active benchmark can still read; wait until it finishes.
  An actual unresolved conflict or external blocker must be reported explicitly, not
  described as a completed merge. This is not a request for routine merge approval.
- **Verify completion and tag it.** Confirm the task commit is reachable from `main`,
  update context with the merge outcome, and create a new annotated milestone tag on
  the integrated state. Keep producing-run tags unchanged. The final response names
  the commit/tag and merge status. Do not delete a worktree containing unfinished work.
- **Code, SQL, docs and analyses** may be written in parallel.
- **Measurements may not.** Dev checks, matrices and anything that starts a database run
  one at a time on this machine, for two reasons: the infra scripts reuse container names
  in every study and worktree (a second `up` destroys the first run's database), and even
  with unique names two runs would share the same eight cores and silently change each
  other's numbers — this machine's database containers are CPU-throttled even when running
  alone (study 02, LESSONS_LEARNED).

**The benchmark lock enforces the second rule.** Every runner calls
`run_lock_acquire` from `infra/lib.sh`, which creates the podman volume `ads-run-lock`.
Creation is atomic and the podman machine is shared by every worktree, shell and agent, so
a second runner — from any folder, any tool — refuses to start and prints who holds the
lock. The topology scripts (`infra/pg-single.sh` etc.) refuse to start or stop databases
while someone else holds it. Rules:

- Any new script that measures or starts databases must call `run_lock_acquire` (or be
  called by a script that did). Never bypass it for a measurement.
- If the lock is held, **wait**; do not remove it. Check `podman volume inspect
  ads-run-lock` for the holder and `CONTEXT.md` for what is running.
- A stale lock (holder killed with SIGKILL, `podman ps` shows no benchmark or database
  containers) may be removed with `podman volume rm ads-run-lock`; say so in `CONTEXT.md`.
- `ADS_IGNORE_RUN_LOCK=1` exists for poking at a database by hand when you are certain
  nothing is measuring. Never set it inside a runner.
- Record a long run in `CONTEXT.md` ("running — do not start databases") when it starts.

## Hard rule: assume a concurrent agent; the later merger reconciles

**Owner rule, 2026-09-15:** another AI agent, from any vendor, is always expected to be working on
this project at the same time. Both must be able to work without interfering, and both agents'
findings end up merged on `main`. Nobody waits for the other to finish before starting.

**Never interfere:**

- Assume the other agent exists even when you cannot see it. Before touching anything
  shared, check:
  - `git worktree list`;
  - the benchmark lock (`podman volume inspect ads-run-lock`);
  - `CONTEXT.md`.
- Work only in your own task worktree and branch. Never check out, reset, merge into or
  fast-forward a branch in a folder you do not own. Never touch another task's worktree,
  branch, untracked files or containers.
- A long run's results live uncommitted in its worktree until the run ends. Keep a runner and
  its results out of any folder another agent may need, including the checkout of `main`.
- Keep edits to shared documents small and local to your own study or task. This keeps
  merges mechanical.

**The later merger reconciles.** When you merge into `main` and it has moved, you own the
integration of both sides. Integrate the latest `main` into your task branch, then make the shared
documents read as one coherent project, as if one author had kept them:

- **`CONTEXT.md`:** one current status per study or task. Where the merged commits prove an
  earlier statement stale (for example "ER-02 open" after its decision was merged), replace it with
  the current state and keep the fact of what happened. Tags and runs tables list both sides once.
- **`LESSONS_LEARNED.md`:** keep every distinct lesson under the section it belongs to. When two
  sessions wrote the same lesson, merge them into one entry that keeps both sides' evidence.
- **`AGENTS.md`, `docs/methodology.md`, templates:** one rule set. Merge two wordings of the
  same rule into one that satisfies both, and remove duplicates and contradictions. If two
  rules genuinely conflict, keep both visible, mark the conflict and ask the owner. Do not
  pick a side.
- **READMEs and indexes:** one consistent terminology and structure. Every study and report
  stays reachable.

Reconciliation edits living shared documents only. It never changes:

- signed analyses or discussions;
- another session's handoff, progress or escalation entries;
- reports, results or manifests;
- history (no rewrite, no force, no moved tags).

It never deletes another session's findings or evidence. Rule 10 still protects attributed
files. A disagreement with another session's conclusion becomes your own signed file, not an
edit. Name the reconciliation in the merge commit message, including what was consolidated, and
check the combined result before fast-forwarding `main`.

## Hard rule: correctness gates timing

A design that answers the question wrongly, quickly, is worth nothing.

Every design must pass verification — answers checked against values computed
independently in Go from the generated dataset — before any timing is reported. A failure
aborts the cell. Where a study's question is an invariant, keep at least one **negative
control** that must be seen to violate it (methodology 5a).

Do not weaken this to make a run pass. The gate's first two catches were bugs in the
*harness*, not in the designs (see `LESSONS_LEARNED.md`); a benchmark without it does not
fail loudly, it publishes wrong numbers quietly.

## Hard rule: numbers and conclusions are separate files

A generated report (`reports/<run-id>.md`) holds **measurements and no interpretation**.
Conclusions go in a **signed analysis** under `reports/analyses/`, created from that study's
`reports/analyses/TEMPLATE.md`. Detailed mechanisms, SQL/plan walkthroughs and exchanges
between analysts go in signed companions under `reports/discussions/`
(`docs/templates/DISCUSSION.md`, methodology 11b).

Before writing an analysis:

1. Generate the report and read its **inputs digest** — a content hash over every result
   file in the run. That, not the run id, identifies the data.
2. Check the report's analyses index. **If an analysis already exists for that digest,
   read it first.** Write a new one only to disagree, extend, or bring a different
   perspective — never to restate what is already there.
3. **Never edit another analyst's file** — including an analysis by an earlier session of
   your own model. Write your own and cite theirs by `analysis_id`. Where two analyses
   disagree, that disagreement is a finding and stays visible. To replace your own
   analysis, write a new file with `supersedes:` and move the old one to `outdated/`.
4. Fill in the frontmatter honestly — `analyst` (model id, not vendor brand alone),
   `analyst_kind`, `analyst_version` (model, version and tool), `analyzed_at`, `run_id`,
   `inputs_digest`, `repo_commit`.
5. Every analysis includes a section naming **where the measurement is weak**.
6. **Every final report opens with a TL;DR** (methodology 11a): measured facts selected by
   fixed rules in a generated report; the actionable conclusion in a signed analysis.
7. Check every figure you cite against the report before committing.

Regenerate the report after adding an analysis so the index picks it up, using the
benchmark image recorded in the run's manifest when it still exists.

## Layout

```
AGENTS.md        these instructions (CLAUDE.md imports it)
CONTEXT.md       current state, tags, running work
LESSONS_LEARNED.md
docs/            methodology, replication, environments/, templates/
infra/           versions.env, lib.sh (incl. the benchmark lock), one script per topology
platform/        shared Go module `adsplatform` — hexagonal (see platform/README.md)
  core/          measure, catalog, provenance, inspect — pure, driver-free
  ports/         DB / Tx / Row / ErrorClass interfaces
  adapters/      pgxdb (the only driver import), filestore, markdown, cgroup
studies/NN-name/ everything belonging to ONE study, nothing shared
  README.md      the question, the designs, how to run it
  study.env      the study's image name and engine flags (study 02 onwards)
  sql/<design>/  schema.sql indexes.sql queries.sql writes.sql [audit.sql] [triggers.sql]
  diagrams/      PlantUML sources + rendered/*.svg
  harness/       Go benchmark client (main.go is the entry adapter)
  run-*.sh       runners (each takes the benchmark lock)
  results/<run-id>/   JSON, plans/, logs/, manifest.yaml
  reports/ + reports/analyses/ + reports/discussions/ + reports/outdated/
```

**Shared code goes in `platform/`, study code stays in its study.** Core logic depends only
on `ports/`; only a study's `main.go` imports adapters. Study 01 predates the platform and
keeps its own harness (tag `study-01/v1`); do not migrate it without a re-run that shows
identical numbers.

## The SQL catalogue format

Designs are defined by readable SQL files, embedded into the binary with `go:embed`, so
the SQL that produced a result is provably the SQL sitting beside it. Named statements
look like:

```sql
-- name: q03_top_donor_charity
-- params: charity_id
-- Free-form documentation lines.
SELECT ...;
```

`params` names are bound by the harness, which is what lets one driver execute statements
whose signatures differ between designs.

## Adding a design

1. `sql/<id>/` with the design's files. Keep the comments: these files are the study's
   documentation as much as its executable content.
2. Register it in `harness/designs.go` with the flags that tell the loader how to shape
   data (study 01: `CharityOnDonation`, `Rollups`, `Embedded`, `Triggers`, `AppRollup`;
   study 02: `Strategy`, `Isolation` and loader flags in the `designs` slice).
3. Add it to the runner's design list (`ALL_PG_DESIGNS` / `ALL_YB_DESIGNS` in study 01,
   `ALL_DESIGNS` in study 02).
4. Add or update a `.puml` diagram and re-render with `./diagrams/render.sh svg`.
5. Study 01: add it to `designOrder` and `designShort` in `harness/report.go`. Both
   studies: add a controlled pair in `pairs` if it isolates one decision.
6. Verify on both engines before running a matrix.

**Prefer designs that differ from an existing one by exactly one decision.** A design
that changes three things at once produces a number nobody can attribute.

## Environment gotchas

- **Windows/Git Bash:** route host paths through `hostpath()` from `infra/lib.sh` before
  handing them to podman **or git**. An unconverted path mounts an empty directory without
  erroring, and made git record a run's commit as "unknown".
- **YugabyteDB** binds YSQL to its `--advertise_address`, not loopback, and ships `ysqlsh`
  rather than `psql`.
- **Check YugabyteDB's effective isolation, don't recall it.** On the pinned 2025.2.6 image
  READ COMMITTED is effective by default (verified 2026-09-13); older releases ran it as
  Snapshot Isolation. Study 02 pins it with a flag and records
  `engine_info.effective_isolation`.
- **Do not rebuild a benchmark image while a matrix using it is running.** Each study has
  its own image name (`study.env`), so studies do not collide; two worktrees of the same
  study do — the benchmark lock serialises them.
- **Do not edit a bash runner that is running.** Study 02's runner re-executes from a copy of
  itself and is safe; study 01's runners are not.
- Writing Go/SQL through shell heredocs has repeatedly failed on quoting. Write those
  files directly. Never create files named `NUL`, `CON`, `AUX` etc. (Windows reserved names
  break `git add`).

## Style

- Comments in SQL and Go explain *why a design is interesting*, not what the syntax does.
  These files are read as documentation.
- Reports state their limitations in the report, not in a footnote. The laptop-shared-cores
  and no-real-network caveats belong wherever conclusions are drawn.
- Report what was measured. If a cell failed, say which and why.
- Name vendors only where provenance needs it (analyst ids, tool-specific setup). Rules are
  written for any agent.
