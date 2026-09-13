# Working in this repository

Context for AI sessions. Read [`CONTEXT.md`](CONTEXT.md) for current state and
[`docs/methodology.md`](docs/methodology.md) for the rules that govern every study.

## What this project is

An experimental record of how architecture and data-modelling decisions affect
performance, **measured rather than asserted**. The output is reproducible runs and
written reports, not opinions.

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

## Hard rule: correctness gates timing

A design that answers the question wrongly, quickly, is worth nothing.

Every design must pass `-cmd verify` — answers checked against values computed
independently in Go from the generated dataset — before any timing is reported. A failure
aborts the cell.

Do not weaken this to make a run pass. The gate's first two catches were bugs in the
*harness*, not in the designs (see `LESSONS_LEARNED.md`); a benchmark without it does not
fail loudly, it publishes wrong numbers quietly.

## Hard rule: numbers and conclusions are separate files

A generated report (`reports/<run-id>.md`) holds **measurements and no interpretation**.
Conclusions go in a **signed analysis** under `reports/analyses/`, created from
`reports/analyses/TEMPLATE.md`.

This project expects **several analysts, including different AI models**, to read the same
data. So, before writing an analysis:

1. Generate the report and read its **inputs digest** — a content hash over every result
   file in the run. That, not the run id, identifies the data: a run can be extended with
   extra cells afterwards.
2. Check the report's analyses index. **If an analysis already exists for that digest,
   read it first.** Write a new one only to disagree, extend, or bring a different
   perspective — never to restate what is already there.
3. **Never edit another analyst's file.** Write your own and cite theirs by `analysis_id`.
   Where two analyses disagree, that disagreement is a finding and stays visible.
4. Fill in the frontmatter honestly — `analyst`, `analyst_kind`, `analyst_version`,
   `analyzed_at`, `run_id`, `inputs_digest`. A future reader needs to know which model, at
   which version, on what date, looking at what.
5. Every analysis must include a section naming **where the measurement is weak**. An
   analysis with no stated doubts has not been done carefully.
6. **Every final report opens with a TL;DR** (methodology 11a). In a generated report it is
   measured facts selected by fixed rules; in a signed analysis it is the actionable
   conclusion, backed by measurements cited in the body.

Regenerate the report after adding an analysis so the index picks it up.

## Layout

```
docs/            methodology, replication, environments/
infra/           versions.env + one script per topology
platform/        shared Go module `adsplatform` — hexagonal (see platform/README.md)
  core/          measure, catalog, provenance, inspect — pure, driver-free
  ports/         DB / Tx / Row / ErrorClass interfaces
  adapters/      pgxdb (the only driver import), filestore, markdown, cgroup
studies/NN-name/ everything belonging to ONE study, nothing shared
  README.md      the question, the designs, how to run it
  study.env      the study's image name and engine flags
  sql/<design>/  schema.sql indexes.sql queries.sql writes.sql [audit.sql] [triggers.sql]
  diagrams/      PlantUML sources + rendered/*.svg
  harness/       Go benchmark client (main.go is the entry adapter)
  run-study.sh   the matrix runner
  results/<run-id>/   JSON, plans/, logs/, manifest.yaml
  reports/ + reports/analyses/ + reports/outdated/
```

**Shared code goes in `platform/`, study code stays in its study.** Core logic depends only
on `ports/`; only a study's `main.go` imports adapters. Study 01 predates the platform and
keeps its own harness (tag `study-01/v1`); do not migrate it without a re-run that shows
identical numbers.

**Provenance.** Start any run that will be analysed from a committed tree with
`run-study.sh --tag` (tag `run/<study>/<run-id>`). Results record commit, describe and
dirty flag; analyses record `repo_commit`. Several sessions and analysts may be editing the
repo at once: commit only the paths you changed, and never commit someone else's edits.

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
whose signatures differ between designs. Adding a design means adding a directory under
`sql/` and an entry in `harness/designs.go` — no driver changes.

## Adding a design

1. `sql/<id>/` with the four (or five) files. Keep the comments: these files are the
   study's documentation as much as its executable content.
2. Register it in `harness/designs.go`, setting the flags that tell the loader how to
   shape data (`CharityOnDonation`, `Rollups`, `Embedded`, `Triggers`, `AppRollup`).
3. Add it to `ALL_PG_DESIGNS` / `ALL_YB_DESIGNS` in `run-study.sh`.
4. Add a `.puml` diagram and re-render with `./diagrams/render.sh svg`.
5. Add it to `designOrder` and `designShort` in `harness/report.go`, and add a controlled
   pair in `pairs` if it isolates one decision against an existing design.
6. Verify on both engines before running a matrix.

*Study 02 differs:* designs carry `audit.sql` and are registered in `designs` (with a
`Strategy`, `Isolation` and loader flags) in `harness/designs.go`; `designShort` and the
report's order come from that slice, so steps 3 and 5 reduce to `ALL_DESIGNS` in
`run-study.sh` and `pairs` in `harness/report.go`. If the study's question is an invariant,
keep at least one **negative control** that must be seen to violate it (methodology 5a).

**Prefer designs that differ from an existing one by exactly one decision.** A design
that changes three things at once produces a number nobody can attribute.

## Environment gotchas

- **Windows/Git Bash:** always route host paths through `hostpath()` from `infra/lib.sh`
  before handing them to podman. An unconverted path mounts an empty directory *without
  erroring* and the run produces no files.
- **YugabyteDB** binds YSQL to its `--advertise_address`, not loopback, and ships `ysqlsh`
  rather than `psql`.
- **Do not rebuild the benchmark image while a matrix is running** — later cells would run
  a different binary than earlier ones.
- Writing Go/SQL through shell heredocs has repeatedly failed on quoting. Write those
  files directly.

- **Check YugabyteDB's effective isolation, don't recall it.** Older releases ran READ
  COMMITTED as Snapshot Isolation unless `yb_enable_read_committed_isolation=true` was set;
  on the pinned 2025.2.6 image RC is effective by default (verified 2026-09-12). Study 02
  sets the flag anyway to pin the behaviour, and every result records
  `engine_info.effective_isolation`.
- **Infra container names are shared by every study** (`pg-single`, `yb-n1`…). Study 02's
  runner refuses to start if they already exist; check `podman ps` before any run.
- Study 02's runner re-executes from a temp copy of itself, so it is safe to edit mid-run.
  Study 01's runners are not.

## Style

- Comments in SQL and Go explain *why a design is interesting*, not what the syntax does.
  These files are read as documentation.
- Reports state their limitations in the report, not in a footnote. The laptop-shared-cores
  and no-real-network caveats belong wherever conclusions are drawn.
- Report what was measured. If a cell failed, say which and why.
