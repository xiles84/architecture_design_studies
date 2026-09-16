# Execution Handoff — recency question (study 01) and operational reports (studies 02, 03)

| Field | Value |
|---|---|
| Handoff ID / revision | **EH-02 / v1** |
| Planner | Claude Opus 5, setting `ultracode` (HIGH), Claude Code desktop, 2026-09-15 |
| Starting source / main revision | local `main` at `7123da6`; task worktree `.worktrees/recency-reports`, branch `repo/recency-and-reports`, created from `7123da6` |
| Checkpoint tag | `repo/recency-reports-handoff-v1` (annotated; never moved) |
| Status | ready |
| Next setting | **LOW — Claude Opus 5, setting `high`** |
| Progress / escalations | [`PROGRESS.md`](PROGRESS.md) · [`ESCALATIONS.md`](ESCALATIONS.md) (this task owns both) |

## Objective and boundaries

The owner asked, on 2026-09-15, for two things:

1. **Study 01** — a query for *"the people that made the LAST donation last week (or
   another period of time)"*, with their own suggestion (a `last_donation` flag on
   `donation`, indexed) as one candidate among others.
2. **Studies 02 and 03** — *"reports and queries that similar real cases like ours need"*,
   integrated into the existing scenarios and tests.

HIGH has written the two protocols that answer them. **Read them first; they are the
specification, this handoff is the execution order:**

- [`studies/01-charity-tree/RECENCY.md`](../../../studies/01-charity-tree/RECENCY.md) —
  the question's exact semantics, the seven new designs, controlled pairs, correctness,
  experiments, sizing, limitations.
- [`studies/02-ticket-booking/REPORTS.md`](../../../studies/02-ticket-booking/REPORTS.md)
- [`studies/03-reserved-seating/REPORTS.md`](../../../studies/03-reserved-seating/REPORTS.md)

**This handoff covers phase 1 only: study 01's implementation and its dev checks, up to
and including a `small` calibration. It stops before any reported measurement and before
studies 02 and 03 are touched.** Phases 2 (the measured study-01 matrix) and 3 (studies 02
and 03) are specified in the protocols but are deliberately not authorised yet: their
sizing depends on the calibration produced in step 19, which is exactly the kind of
decision that belongs to HIGH.

**Owned paths** (this task may create or edit only these):

```
studies/01-charity-tree/RECENCY.md            (already written by HIGH)
studies/01-charity-tree/sql/d18_*..d24_*/     (new)
studies/01-charity-tree/sql/d1_*..d17_*/queries.sql   (append four statements ONLY)
studies/01-charity-tree/harness/*.go
studies/01-charity-tree/run-enhancements.sh
studies/01-charity-tree/run-study.sh          (design lists only)
studies/01-charity-tree/diagrams/**
studies/01-charity-tree/README.md
studies/01-charity-tree/results/devchecks/**  (new subdirectories only)
studies/02-ticket-booking/REPORTS.md          (already written by HIGH)
studies/03-reserved-seating/REPORTS.md        (already written by HIGH)
docs/handoffs/20260915-recency-and-reports/**
CONTEXT.md, LESSONS_LEARNED.md                (this task's own entries only)
```

**Must not be touched:** any file under `reports/` or `results/` of any study except the
new `results/devchecks/` directories created here; any existing analysis or discussion;
`.worktrees/study01-v3`, `.worktrees/study01-review-audit` and anything in them; any
existing `schema.sql`, `indexes.sql`, `writes.sql` or `triggers.sql` of D1–D17; any
existing query text in D1–D17's `queries.sql`; any published tag.

**Why D1–D17 keep byte-identical schema, indexes and writes:** every published storage,
insert, update and delete number for this study must stay valid and comparable. Only
their read catalogue grows, by exactly four statements appended at the end.

## Decisions made by HIGH

### D-1. What the question is, precisely

Donor `p` qualifies for window `[since, until)` when `MAX(donated_at)` over **all** of
`p`'s donations falls inside it. Two regimes are measured, `trailing` and `historical`
(RECENCY.md §2). The window is derived from the dataset, never from `now()`.

### D-2. The four statements, and the six formulations

Names and parameters are fixed (RECENCY.md §2). Every design's catalogue gets all four.
The canonical SQL is below; transcribe it, adapting only identifier names to that
design's schema. Keep the comment style of the surrounding files — these files are the
study's documentation.

**Formulation A — no copied key (D1, D2, D11, D12).** `q15` reaches the charity through
`person`:

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT person_id, MAX(donated_at) AS last_at
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) >= $1 AND MAX(donated_at) < $2
         ORDER BY 2 DESC
         LIMIT 100) t
  JOIN person p ON p.person_id = t.person_id
 ORDER BY t.last_at DESC;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM (SELECT person_id
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) >= $1 AND MAX(donated_at) < $2) t;

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT d.person_id, MAX(d.donated_at) AS last_at
          FROM donation d
          JOIN person owner ON owner.person_id = d.person_id
         WHERE owner.charity_id = $1
         GROUP BY d.person_id
        HAVING MAX(d.donated_at) >= $2 AND MAX(d.donated_at) < $3
         ORDER BY 2 DESC
         LIMIT 100) t
  JOIN person p ON p.person_id = t.person_id
 ORDER BY t.last_at DESC;

-- name: q16_lapsed_donors_count
-- params: since
SELECT COUNT(*) AS donor_count
  FROM (SELECT person_id
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) < $1) t;
```

**Formulation B — copied key (D3, D7, D8, D13, D14, D15, D16, D17).** Identical to A
except `q15`, which filters `donation.charity_id = $1` directly and drops the join to
`owner`. This is the same one-decision difference the study already measures in q02/q05/q12.

**Formulation C — parent rollup (D4, D5, and byte-identical in D22, D23).** Donors with
no donations have `last_donation_at IS NULL` and are correctly excluded by the
comparisons:

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT person_id, full_name, last_donation_at AS last_at
  FROM person
 WHERE last_donation_at >= $1 AND last_donation_at < $2
 ORDER BY last_donation_at DESC
 LIMIT 100;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM person
 WHERE last_donation_at >= $1 AND last_donation_at < $2;

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT person_id, full_name, last_donation_at AS last_at
  FROM person
 WHERE charity_id = $1 AND last_donation_at >= $2 AND last_donation_at < $3
 ORDER BY last_donation_at DESC
 LIMIT 100;

-- name: q16_lapsed_donors_count
-- params: since
SELECT COUNT(*) AS donor_count
  FROM person
 WHERE last_donation_at < $1;
```

**Formulation D — whole-history document (D6).** `person.donations` is a chronological
JSONB array; the aggregate is taken over the unnested document, which is the design's
honest cost:

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT p.person_id, p.full_name, l.last_at
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX((e->>'donated_at')::timestamptz) AS last_at
                       FROM jsonb_array_elements(p.donations) e) l
 WHERE l.last_at >= $1 AND l.last_at < $2
 ORDER BY l.last_at DESC
 LIMIT 100;
```

`q14`/`q16` are the same shape with `COUNT(*)` and no `ORDER BY`/`LIMIT`; `q15` adds
`p.charity_id = $1`.

**Formulation E — bounded cache (D9, D10).** The cache is newest-first, so element 0 is
the donor's latest gift *when the cache is correct* — which is the design's own claim,
and this query is where a corrupted cache produces a wrong answer:

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT p.person_id, p.full_name,
       (p.recent_donations->0->>'donated_at')::timestamptz AS last_at
  FROM person p
 WHERE (p.recent_donations->0->>'donated_at')::timestamptz >= $1
   AND (p.recent_donations->0->>'donated_at')::timestamptz <  $2
 ORDER BY last_at DESC
 LIMIT 100;
```

Same pattern for q14/q15/q16. **Do not add an index for this**: the cast is `STABLE`, not
`IMMUTABLE`, so PostgreSQL refuses an expression index on it, and that refusal is part of
the finding (RECENCY.md §3). If a plain `CREATE INDEX` is attempted and rejected, that is
expected — do not work around it, and do not switch the design to the table.

**Formulation F — the flag (D20, D21, D24).**

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT d.person_id, p.full_name, d.donated_at AS last_at
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.is_last_donation
   AND d.donated_at >= $1 AND d.donated_at < $2
 ORDER BY d.donated_at DESC
 LIMIT 100;
```

q14/q16 count without the join; q15 adds `d.charity_id = $1`.

**Formulation G — parent-driven probe (D18 only), on D3's schema:**

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT p.person_id, p.full_name, l.last_at
  FROM person p
 CROSS JOIN LATERAL (SELECT MAX(d.donated_at) AS last_at
                       FROM donation d
                      WHERE d.person_id = p.person_id) l
 WHERE l.last_at >= $1 AND l.last_at < $2
 ORDER BY l.last_at DESC
 LIMIT 100;
```

**Formulation H — window-first (D19 only), on D3's schema.** Correct in both regimes: a
donor qualifies exactly when they gave inside the window and gave nothing after it:

```sql
-- name: q13_donors_last_gift_window
-- params: since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT d.person_id, MAX(d.donated_at) AS last_at
          FROM donation d
         WHERE d.donated_at >= $1 AND d.donated_at < $2
         GROUP BY d.person_id) t
  JOIN person p ON p.person_id = t.person_id
 WHERE NOT EXISTS (SELECT 1 FROM donation d2
                    WHERE d2.person_id = t.person_id
                      AND d2.donated_at >= $2)
 ORDER BY t.last_at DESC
 LIMIT 100;
```

For q16 (`last < since`) formulation H degenerates to a scan of everything before
`since`; write it as formulation B's statement and say so in a comment — a formulation
that does not apply to a question is recorded, not faked.

### D-3. The seven new designs

Directories, parents and the single change are in RECENCY.md §3. Construction rule: copy
the parent design's directory verbatim, then make exactly the one change.

| New | Copy from | Then change |
|---|---|---|
| `d18_recency_probe` | `d3_flattened_fk` | `queries.sql`: q13–q16 → formulation G. Nothing else. |
| `d19_recency_window_sql` | `d3_flattened_fk` | `queries.sql`: q13–q16 → formulation H. Nothing else. |
| `d20_recency_flag` | `d3_flattened_fk` | `schema.sql`: `is_last_donation BOOLEAN NOT NULL DEFAULT false` on `donation`; `indexes.sql`: `CREATE INDEX donation_last_flag_idx ON donation (donated_at DESC) WHERE is_last_donation;`; new `triggers.sql` (D-4); `queries.sql`: formulation F |
| `d21_recency_flag_unguarded` | `d20_recency_flag` | `triggers.sql`: the person-row lock removed. **Negative control.** |
| `d22_recency_rollup_idx` | `d4_rollup_trigger` | `indexes.sql`: one line added, `CREATE INDEX person_last_donation_global_idx ON person (last_donation_at DESC);`; `queries.sql`: formulation C |
| `d23_recency_rollup_app_idx` | `d5_rollup_app` | the same one index line; `queries.sql`: formulation C, byte-identical to D22's |
| `d24_recency_flag_colocated` | `d20_recency_flag` | `schema.sql`: `donation`'s primary key rebuilt for placement exactly as `d7_yb_child_colocated` does it. YugabyteDB only. |

### D-4. The flag's maintenance (D20), written out

`triggers.sql`, attached after the bulk load like D4's:

- `BEFORE INSERT ON donation FOR EACH ROW`:
  1. `PERFORM 1 FROM person WHERE person_id = NEW.person_id FOR UPDATE;` — the guard.
  2. `UPDATE donation SET is_last_donation = false WHERE person_id = NEW.person_id AND is_last_donation AND donated_at <= NEW.donated_at;`
  3. `NEW.is_last_donation := NOT EXISTS (SELECT 1 FROM donation WHERE person_id = NEW.person_id AND donated_at > NEW.donated_at);`
  4. `RETURN NEW;`
- `AFTER DELETE ON donation FOR EACH ROW`: if `OLD.is_last_donation`, promote the donor's
  next newest surviving donation (`ORDER BY donated_at DESC, donation_id DESC LIMIT 1`).
  Take the same `person` row lock first.
- No trigger on `UPDATE` of `amount_cents`: `donated_at` does not move, so the flag does
  not move. If a trigger fires there, the write benchmark will show it and that is a bug.

**D21 is step 1 removed, and nothing else.** Two concurrent inserts for one donor at READ
COMMITTED then both fail to see the other's uncommitted row, both clear nothing relevant,
and both set their own flag — two flagged rows for one donor. That is what the control
exists to demonstrate.

Comment the *why* in `triggers.sql`, as the other designs do: the lock is the design
decision, not a detail.

### D-5. Harness changes

| Where | Change |
|---|---|
| `harness/designs.go` | seven `Design` entries; new flag `LastFlag bool` on `Design` (loader materialises `is_last_donation`) for D20/D21/D24; D20/D21/D24 also `CharityOnDonation: true, Triggers: true`; D22 copies D4's flags, D23 copies D5's |
| `harness/load.go` | when `LastFlag`, set the column in the COPY stream for each donor's newest donation — the loader already knows it, in the same place `Rollups` and `RecentCache` are materialised. Triggers attach after the load, unchanged |
| `harness/bench.go` | `Binder` gains `EpochEnd`/`EpochStart` (from the dataset) and a `Window` regime; `readVals` binds `since` and `until`. Regimes: `trailing` = `[end−7d, end)`, `historical` = `[end−97d, end−90d)`. Both are measured; the per-query result name carries the regime, e.g. `q13_donors_last_gift_window@trailing`. Catalogue names stay plain |
| `harness/bench.go` (writes) | new write op `insert_backdated`: the existing `w_insert_donation` statement with `donated_at` bound to `epochStart − 24h`. No SQL change in any design. Record how many such inserts became the donor's newest (expected: zero) |
| `harness/verify.go` | truth for q13–q16 in **both** regimes, computed from the dataset: per-donor `MAX(donated_at)`, then the expected set, the expected counts, and the expected top-100 with boundary ties accepted as the existing gate does. Eight new checks per cell |
| `harness/audit.go` | `recencyAudit` (RECENCY.md §5): flagged-row count per donor, flagged row not the newest, `person.last_donation_at` drift, cache head drift. Counts **and up to ten examples** |
| `harness/report.go` | `coreQueries` = the existing twelve, and the legacy geometric mean is computed over those twelve only — the published score must keep meaning what it meant. New "Recency" section: per-query throughput by regime, a per-question winner table, the maintenance costs, the audit outcomes, and a trade-off row per design. `designOrder`/`designShort` entries for D18–D24; new controlled pairs from RECENCY.md §4 |
| `run-enhancements.sh` | four new experiment groups: `recency-reads`, `recency-maintenance`, `recency-hot-donor`, `recency-placement`, added to the `case` validation list and the enabled-group dispatch. Every group keeps the existing lock, fresh-load-per-trial, tag and manifest discipline |
| `run-study.sh` | D18–D23 appended to `ALL_PG_DESIGNS`, D24 to `ALL_YB_DESIGNS` only |

### D-6. Permitted routine choices

LOW may decide, without escalating: Go identifier names, file and function placement
inside the listed files, comment wording, the exact form of the audit's example
formatting, log text, unit-test names and table-driven test data, and the order in which
the steps below are executed where they do not depend on each other.

### D-7. Explicit non-goals for phase 1

No reported run. No tag other than the two named in step 20. No change to studies 02/03
beyond the protocol files HIGH has already written. No materialised view, no
time-bucketed table, no asynchronous/queue-based maintenance design. No migration of
study 01 to `platform/`.

## Ordered execution

Shell: **Git Bash** for anything touching podman or the runners (`hostpath()` matters).
Every command below runs from the task worktree unless stated.

| Step | Exact action or command | Expected evidence | Failure action |
|---|---|---|---|
| 1 | `cd /c/extra/code/architecture_design_studies/.worktrees/recency-reports && git status --porcelain && git log --oneline -1` | clean tree, HEAD = this handoff's checkpoint commit on `repo/recency-and-reports` | if the branch is not this one, stop — you are in another task's folder (ER) |
| 2 | `git worktree list`, `podman volume inspect ads-run-lock`, read `CONTEXT.md` | no other task is running a benchmark; record the holder if the lock exists | lock held → **wait** and re-check; never remove it. Waiting is mapped work, not an escalation |
| 3 | Read `studies/01-charity-tree/RECENCY.md` end to end | — | a contradiction between it and this handoff is an ER |
| 4 | Append q13–q16 to the seventeen existing `sql/d*/queries.sql`, per D-2's formulation map. Append only; change no existing byte | `git diff --stat` shows only additions in those files; `git diff -U0 -- studies/01-charity-tree/sql | grep '^-'` prints nothing but file headers | any existing line changed → revert that file and redo |
| 5 | Create the seven new design directories per D-3 (copy the parent, then the one change) | `diff -r sql/d3_flattened_fk sql/d18_recency_probe` differs only in `queries.sql`; same check for D19; `diff -r sql/d4_rollup_trigger sql/d22_recency_rollup_idx` differs only in `indexes.sql` and `queries.sql`; `diff sql/d22_.../queries.sql sql/d23_.../queries.sql` is empty | a third file differing → fix; do not "improve" the parent |
| 6 | Write `sql/d20_recency_flag/triggers.sql` per D-4, and `d21`'s as D20's minus the lock | `diff sql/d20_recency_flag/triggers.sql sql/d21_recency_flag_unguarded/triggers.sql` shows exactly the lock statement and its comment | — |
| 7 | `harness/designs.go` + `harness/load.go` per D-5 | `go build ./...` in `studies/01-charity-tree` | — |
| 8 | `harness/bench.go`: window binding, regimes, `insert_backdated` | builds; `-cmd list` shows the new designs | — |
| 9 | `harness/verify.go`: eight new checks | builds | — |
| 10 | `harness/audit.go`: `recencyAudit` | builds | — |
| 11 | `harness/report.go`: `coreQueries`, recency section, design order, pairs | builds | — |
| 12 | Unit tests for the parts that do not need a database: window derivation for both regimes, the expected-set computation against a hand-built dataset, the audit's classification, and a test that the legacy geometric mean over a fixture run is **unchanged** by the new queries | `go test ./...` passes; the geomean test is the one that protects published scores | geomean changed → the new queries leaked into `coreQueries`; fix |
| 13 | `gofmt -l harness` and `go vet ./...` | both clean | — |
| 14 | `bash -n run-enhancements.sh run-study.sh` after the runner edits | clean | — |
| 15 | Diagrams and documentation: a `.puml` for `d20_recency_flag` and one for `d22_recency_rollup_idx` (the two designs that change where the fact lives), rendered with `./diagrams/render.sh svg`; study 01's README gains the four questions in its questions table and D18–D24 in its designs table, each linking its diagram where one exists | the SVGs exist and the README tables list them | render script unavailable in the container → record it and continue; it is not a blocker |
| 16 | Commit the implementation (explicit paths, no `git add -A`), message ending with the `Co-Authored-By` trailer for the executing model | one commit per coherent group, at least: SQL, harness, runner, diagrams+docs | — |
| 17 | **Dev checks, `tiny`, PostgreSQL first.** Take the lock through the runner; never bypass it. Run the verify group for all designs on `pg-single`, then the maintenance and hot-donor groups for D20–D23 at 8 writers. Results to `results/devchecks/recency-dc01-pg/` | every correct design passes all eight new checks in both regimes; `insert_backdated` never moves a flag; **D21 fails the recency audit** with examples; D20/D22/D23 pass it | a correct design failing the gate → fix the design or the harness, never the gate. D21 **not** failing → the contention is too low: raise writers to 16 and record it; if it still does not fire, that is an ER |
| 18 | Repeat step 17 on `yb-single`, then `yb-cluster3`, one at a time | same expectations. Also confirm YugabyteDB accepts the **partial** index in D20 | partial index rejected → fall back to `CREATE INDEX ... ON donation (is_last_donation, donated_at DESC)` for YugabyteDB only, record the deviation in PROGRESS.md, and continue: this is a mapped decision |
| 19 | **Calibration.** One `small` cell, `pg-single`, design D3, the recency read group only, three trials. Record wall-clock per query per regime and the client's CPU throttling | a measured per-cell duration that HIGH can multiply out | a cell exceeding 20 minutes → stop and report; do not start a matrix |
| 20 | Commit results and notes; tag `study-01/v4-harness` (annotated) on that commit; update `PROGRESS.md`, `CONTEXT.md` (study 01 section + the runs/tags tables) and `LESSONS_LEARNED.md` if anything was learned | `git tag -n1 study-01/v4-harness` resolves; `CONTEXT.md` says what is implemented and that nothing is running | tag name already taken → **stop** (ER); never move or reuse a tag |
| 21 | **STOP.** Report to HIGH: what was implemented, every dev-check outcome, the calibration number, the D21 evidence, any deviation recorded in step 18, and the open escalations | the final response names the commit and the tag, and states that no reported run has been made | — |

Restart behaviour: every step is idempotent except the dev-check runs, which create a new
`results/devchecks/<name>/` directory each time — use a fresh name (`-dc02`, `-dc03`) and
keep the earlier one.

## Escalation Required

Open an entry in [`ESCALATIONS.md`](ESCALATIONS.md) (template:
`docs/templates/ESCALATION_REQUIRED.md`), stop the dependent step, continue the
independent ones, and do **not** edit this handoff:

- **ER trigger 1** — a correct design cannot pass the new gate without changing the
  question, the window semantics, the limit, or the tie rule.
- **ER trigger 2** — D21 does not fail the recency audit even at 16 writers (step 17).
  The control is the proof that the audit works; a silent control invalidates the other
  designs' clean audits.
- **ER trigger 3** — an engine rejects something the protocol depends on, other than the
  partial index already mapped in step 18.
- **ER trigger 4** — the calibration projects a matrix longer than **12 hours**, or a
  single `small` cell longer than 20 minutes.
- **ER trigger 5** — the flag's maintenance cannot be made correct for `delete_person`
  (cascade erasure) without changing another design's write path.
- **ER trigger 6** — any instruction here contradicts `AGENTS.md`, the methodology, or
  RECENCY.md.
- **ER trigger 7** — another agent's work appears in the shared documents and reconciling
  it needs a judgement about someone else's finding.

A benchmark lock held by another task is **not** an escalation: wait for it.

## Acceptance, integration and next role

Phase 1 is accepted by HIGH when all of these are true and evidenced in `PROGRESS.md`:

1. D1–D17 differ from `main` only by four appended statements per `queries.sql`, proven
   by the diff check in step 4.
2. The seven new designs differ from their parents by exactly the one change, proven by
   the `diff -r` checks in step 5.
3. Every correct design passes the eight new checks on all three topologies, in both
   regimes, and `insert_backdated` never moved a flag.
4. **D21 was seen to fail the recency audit**, with examples recorded.
5. The legacy twelve-question geometric mean is unchanged, proven by the unit test.
6. `go vet`, `gofmt`, `go test` and `bash -n` are clean.
7. The calibration number exists.
8. Nothing is running, no lock is held, no database container is left up.

**Integration.** Phase 1 is a **checkpoint inside an unfinished task**, not a completed
task: commit and tag it, update the context, and leave the branch in place. Do **not**
merge `repo/recency-and-reports` into `main` yet — the merge happens when the task's
acceptance conditions are met, after the measured runs and the signed analyses. Never
push, never pull, never move a tag, never rewrite history.

**Next after execution: HIGH — review the dev-check evidence, size the measured matrix
from the calibration, and publish amendment AM-01 authorising phase 2 (the study-01
recency matrix) and phase 3 (studies 02 and 03).**

## Amendments

*(none yet — HIGH appends dated, attributed `AM-NN` entries here; previous instructions
are marked superseded, never erased)*
