# Working in this repository

Instructions for **every AI agent** working here, whichever vendor or tool runs it (Claude
Code reads this file through `CLAUDE.md`; OpenAI Codex and other agents read `AGENTS.md`
directly). This is the single source of these rules: change them here, nowhere else.

Before starting, read [`CONTEXT.md`](CONTEXT.md) for current state (including work other
sessions have in progress) and [`docs/methodology.md`](docs/methodology.md) for the rules
that govern every study.

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
8. **AI agents commit and tag their own work, and never push.** The owner pushes and
   pulls.
9. **One benchmark on the machine at a time.** Agents may write code and analyses in
   parallel; they may not measure in parallel.
10. **Never edit another analyst's file.** Not to shorten it, restructure it or correct it.

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

## Hard rule: working in parallel

Several agents can work at the same time. They **must not share one working folder**: a
folder has one checked-out branch, and two agents editing the same files — or one
switching branches — change each other's work mid-edit. Use **git worktrees**: one folder
per agent, one branch per folder, one shared repository (commits and tags are visible to
all worktrees at once, with nothing to push between them).

```bash
git worktree add ../ads-<study>-<topic> -b <study>/<topic>
```

- **Branch names:** `<study>/<topic>` (e.g. `study-02/retry-backoff`) or `repo/<topic>`.
- **Merging into `main`** is done by the owner, or by an agent the owner asks to. Expect
  conflicts in `CONTEXT.md` and `LESSONS_LEARNED.md`; resolve them by keeping every
  session's entries, never by choosing one side. Tag after the merge if it changes a study.
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
