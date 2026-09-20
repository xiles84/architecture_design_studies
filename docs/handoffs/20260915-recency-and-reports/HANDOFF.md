# Execution Handoff — recency question (study 01) and operational reports (studies 02, 03)

| Field | Value |
|---|---|
| Handoff ID / revision | **EH-02 / v1, amended by AM-01 and AM-02 (2026-09-16)** |
| Planner | Claude Opus 5, setting `ultracode` (HIGH), Claude Code desktop, 2026-09-15 |
| Starting source / main revision | local `main` at `7123da6`; task worktree `.worktrees/recency-reports`, branch `repo/recency-and-reports`, created from `7123da6` |
| Checkpoint tag | `repo/recency-reports-handoff-v1` (annotated; never moved) |
| Status | phases 1 and 2 complete and analysed; AM-02 authorises phase 3a (studies 02/03 reporting, implementation and dev checks) |
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

### AM-01 — phase 1 accepted with five repairs; phase 2 authorised after them

- Decided by: Claude Opus 5, setting `ultracode` (HIGH), Claude Code desktop, 2026-09-16.
- Reviewing: LOW's phase 1 at commit `129d431`, tag `study-01/v4-harness`.
- Next setting: **LOW — Claude Opus 5, setting `high`** (or a smaller model; see AM-01.7).

**Phase 1 is accepted.** All eight acceptance conditions are met and the evidence is
real: 22 PostgreSQL and 24 YugabyteDB designs pass a gate that was *demonstrated* to
work by catching three genuine bugs, and D21's negative control was seen to fail twice
under two different load shapes, with examples. The lesson LOW drew about *why* the
paced experiment needed 16 writers where the closed-loop one needed 8 is the kind of
finding this study exists to produce, and it is recorded properly.

**But phase 2 cannot start yet.** Reviewing the code against the runner it will actually
execute under found four defects and one untested risk. Three of them exist because
**this handoff's D-5 under-specified the reporting path**, which is my error, not the
executor's: D-5 named `harness/report.go` without noticing that
`run-enhancements.sh` — the runner every recency group uses — ends by calling
`-cmd report-enhancements`, which is `harness/report_enhancements.go`, a different
reporter that D-5 never mentioned.

**AM-01.1 — `recency-maintenance` would fail on its first cell.** `experiment.go`'s
`writes` mode refuses a cell with more than one operation ("each experiment write cell
must specify exactly one operation"). The group as written passes
`-write-ops insert,insert_backdated,delete,update` in a single cell, so **every cell in
that group errors immediately**. It was never caught because the maintenance dev check
ran `-cmd full` directly rather than through the group. Fix: emit one cell per operation,
following the `exceptions` group's existing `for op in update delete; do cell ...` shape.

**AM-01.2 — the experiment `writes` path never runs `AuditRecency`.** LOW correctly wired
it into `main.go`'s write loop and into `arrival.go`, but `experiment.go`'s `writes` case
audits only rollups. So even after AM-01.1, the maintenance group would record **no
recency audit at all** — and D21's failure would be invisible in the very group designed
to price flag maintenance. Fix: in `experiment.go`, alongside each existing
`AuditRollups` call in the `writes` and `arrival` cases, run `AuditRecency` when
`d.LastFlag || d.RecencyRollupIdx` and store it on `run.RecencyAudit` (and on
`run.Writes[0].RecencyAudit` where the rollup audit is already attached that way).

**AM-01.3 — a fired negative control would not be reported as a problem.**
`experimentProblems` in `report_enhancements.go` flags `r.Audit` and
`a.WarmupAudit` mismatches but knows nothing about the recency audits. A run where D21
corrupts flags would therefore generate a report that does not say so. **This is a
methodology 5a violation in the reporting path and is the most important of the five
repairs:** a control that fires invisibly is no control. Fix: extend
`experimentProblems` to flag `r.RecencyAudit` and `Arrival.WarmupRecencyAudit`
mismatches, in the same wording style as the existing entries.

**AM-01.4 — recency results need a recency-aware presentation in the reporter that is
actually generated.** Today the enhancements reporter would render the eight
`q13..q16@regime` entries as anonymous rows in its generic per-operation table: no
statement of which regime is which question, no per-question winner, no pairing of read
gain against maintenance cost, no audit outcome. Two things are required, and no more:

1. In `report_enhancements.go`, a section that lists the recency questions with their
   regime spelled out, and — mandatory — **the recency audit outcome per design**,
   including an explicit line naming D21 as the negative control and stating whether it
   fired. Where no recency audit ran for a design, print that it was not measured, never
   a zero (study 03's "transient: not recorded" precedent).
2. Leave the richer section LOW built in `report.go` exactly as it is. It is correct for
   the survey reporter and will be wanted when a recency matrix is ever run through
   `run-study.sh`. Do not delete it and do not duplicate its code — a shared helper is
   welcome but not required.

Note for whoever implements this: the existing `experimentReadScore` already skips any
`@`-suffixed query and requires exactly twelve, so it is **already** immune to the new
questions, and a recency-only cell correctly contributes no score rather than a zero.
Do not "fix" that function; it is behaving correctly.

**AM-01.5 — `delete_person` is untested on the flag designs.** ER trigger 5 of this
handoff named exactly this risk and the dev checks never exercised it. D20's `AFTER
DELETE` trigger fires **once per child row** and each firing takes the donor's `person`
row lock, runs an ordering query over the survivors and issues an `UPDATE` — during a
cascade that deletes the whole history in one statement. That is plausibly pathological
and it is in study 01's standard write set, so a long matrix could hit it far from a
human. Before any measured run: dev-check `-write-ops delete_person` on D20 and D21 at
`tiny` on `pg-single`, confirm the recency audit passes afterwards, and record the
observed throughput. If it errors, deadlocks, or is so slow the cell cannot finish in
five minutes, **stop and escalate** — do not redesign the trigger to make it pass.

**AM-01.6 — then run phase 2, sized as follows.** Derived from LOW's calibration
(~30 s of measured time per read cell per trial; reads are duration-boxed, so a slow
engine costs the same wall-clock as a fast one, and only load time differs):

| Group | Cells at 3 trials | Estimate |
|---|---|---|
| `recency-reads` | 12 designs × 2 topologies × 3 | ~2.2 h |
| `recency-maintenance` (after AM-01.1, incl. `delete_person`) | ~4 designs × 5 ops × 3, PostgreSQL only, minus D23's skips | ~30 min |
| `recency-hot-donor` | 4 writer levels × 4 designs × 3 | ~30 min |
| `recency-placement` | 2 designs × 2 modes × 3, `yb-cluster3` | ~45 min |

**Total ≈ 4 hours**, comfortably inside ER trigger 4's 12-hour bound. Run with
`--trials 3 --tag`, one group at a time, in the order above, and record the run in
`CONTEXT.md` as running before it starts. The hot-donor group keeps **16 writers** in its
sweep — LOW established that 8 does not reproduce the race under paced arrival, and a
sweep that cannot make its own control fire measures nothing about the others.

**AM-01.7 — model.** AM-01.1 through AM-01.3 are mechanical and well specified. AM-01.4
involves judgement about presentation and AM-01.5 may surface a real design problem. A
smaller model is acceptable for .1–.3; keep the executing model for .4 and .5, and
escalate rather than improvise if `delete_person` misbehaves.

**Phase 3 (studies 02 and 03) remains unauthorised** and is unchanged by this amendment.

**Next after execution: HIGH — review the repairs and the phase 2 results, then write the
signed analysis.**

### AM-02 — phase 3a authorised: studies 02 and 03 reporting, implementation and dev checks

- Decided by: Claude Opus 5, setting `ultracode` (HIGH), Claude Code desktop, 2026-09-16.
- Follows: phase 2 measured and analysed (`study-01/v4-analysis`, digest `1b05f142fd9e06ef`).
- Next setting: **LOW — Claude Opus 5, setting `high`.**

Study 01's half of the owner's 2026-09-15 request is finished. This authorises the other
half — the operational reports of
[study 02](../../../studies/02-ticket-booking/REPORTS.md) and
[study 03](../../../studies/03-reserved-seating/REPORTS.md) — as **implementation and dev
checks only**. No measured run; phase 3b will authorise those once this is reviewed and
the cost is known, exactly as phase 1 gated phase 2.

**AM-02.0 — two claims in those protocols are now verified, and they are the finding.**
I wrote both after reading one design each; I have since checked all 28. State them as
established:

- **Study 02: cancellation destroys the sale in every design.** C1–C5, R1–R3 and H0–H1
  `DELETE FROM ticket`; P1–P4 reset the row to `available` with `customer_id = NULL,
  sold_at = NULL`. The designs' own comments say it outright ("A refund deletes the
  ticket"; "The row stays; only its state changes"). **No design can report refunds, and
  `r01`'s revenue is wrong across a refund in all fourteen.**
- **Study 03: no design retains any hold history.** Zero designs have a hold-event table
  or any `expired_at`/`released_at` column; `w_release_expired` clears the seat row.
  **The abandonment funnel — the number that decides whether 40 minutes is right — cannot
  be computed from any design's state.**

These are results obtained by reading the schemas, and they do not need a benchmark to be
true. What needs measuring is what the *answerable* reports cost, and what making the
unanswerable ones answerable costs on the hot path.

**AM-02.1 — add the reports, changing nothing that exists.** Add `r01`–`r06` from each
study's REPORTS.md to the `queries.sql` of every design that can answer them, appended,
with **no change to any existing statement, schema, index or write path** in the fourteen
designs of either study. Study 02 is measured and analysed, study 03 is measured, analysed
and signed; their published numbers must stay valid. The same append-only diff check
phase 1 used applies: `git diff -U0` over those directories must show no removed lines.

**AM-02.2 — answerability is declared in Go and bound to the SQL by a test.** A design
answers a report iff its catalogue defines that statement. Record coverage in a
`reportCoverage` map in each study's harness — per design, per report, one of
`answerable`, `partial` (with the caveat text) or `unanswerable` (with the reason) — and
add a **unit test that fails if the declaration and the SQL disagree**: every `answerable`
or `partial` entry must have a statement of that name in that design's catalogue, and
every `unanswerable` must not. A declaration that can drift from the SQL is worth little;
one a test pins to it is worth having. The generated report prints the table, and an
unanswerable report is never timed and never shown as a zero.

**AM-02.3 — study 02 gets one new design, X1.** `x1_cas_ledger`: P3 plus an append-only
`sale_event` table written in the same transaction as the sale and the cancellation, with
`r01`, `r03` and `r05` answered from it. Copy `p3_precreated_cas` and make exactly that one
change. **Checked on both axes this time** (LESSONS_LEARNED, "a controlled pair can be
clean on one axis and confounded on another"): on the write axis P3 → X1 is one decision,
one extra insert in the same transaction; on the read axis the two deliberately differ,
because answering `r05` at all is what X1 buys. Say so in the design's own comments.
X1 also needs a **ledger reconciliation audit** — the ledger's sales and cancellations must
reconstruct the live ticket state exactly, and must match the sales the harness saw commit.
Study 03 gets **no new design**; its REPORTS.md already reasoned that out and the reasoning
still holds.

**AM-02.4 — harness work, both studies.** The read helper in each study passes exactly one
bound parameter (`db.Query(ctx, st.SQL, key(r))` in study 02's `workload.go`); the new
reports take up to three. Add a multi-parameter variant rather than reshaping the existing
one, so the existing five/six read questions keep their exact call path. Window parameters
follow study 01's rule and RECENCY.md §2: **derived from the dataset's own time span, never
from `now()`**, and `r03`/`r05` measured in both the trailing and historical regimes with
the regime in the result name. Reuse study 01's naming (`name@regime`) — both studies'
reporters already tolerate an `@` suffix.

**AM-02.5 — dev checks, then stop.** `tiny`, PostgreSQL then YugabyteDB 1-node, under the
benchmark lock, for both studies: every design passes its correctness gate including the
new reports; the answerability table matches the SQL; X1's ledger audit passes and its
negative controls (C1, H0 in study 02) still fire where they fired before. Record the
per-report throughput observed at `tiny` so phase 3b can be sized from data rather than
from my guess. **Do not start a measured run.**

**AM-02.6 — what not to do.** Do not add a hold-history design to study 03. Do not touch
study 03's signed analysis or study 02's. Do not change either study's race, churn or
lifecycle experiments — the reports belong to the read phase, and whether a finance query
throttles a live drop is a separate experiment that both protocols already name as out of
scope. Do not "fix" `experimentReadScore`-style aggregate guards in either study's
reporter without checking first whether they already exclude `@`-suffixed names, as study
01's did.

**AM-02.7 — when phase 3b is planned, size the race honestly.** Study 02's sell-out race is
where X1's cost must show. When that run is specified, the buyers-per-event count is the
experimental variable and the offered load must not be pinned below capacity across it —
see LESSONS_LEARNED, "a fixed arrival rate measures compliance, not capacity". That rule
cost this task a null result in phase 2 and must not cost it another.

**Next after execution: HIGH — review the answerability tables against the SQL, confirm
X1's ledger audit is real, and size phase 3b's measured runs from the dev-check numbers.**

### AM-03 — phase 3a review: eleven repairs, then a `small` calibration, then phase 3b

- Decided by: Claude Opus 5 (HIGH role; the effort setting is not visible to this session),
  Claude Code desktop, 2026-09-16.
- Reviews: `39f2a27..4c940f1` (LOW iteration 3, executed by Claude Sonnet 5).
- Next setting: **LOW — Claude Sonnet 5** (the owner's LOW choice for this task).

**AM-03.0 — what the review accepted, checked rather than taken from the progress log.**

- **Append-only holds.** `git diff -U0 8c0f2a4 4c940f1` over both studies' `sql/`, X1
  excluded, has zero removed lines. X1's `writes.sql` differs from P3's only in
  `w_sell_seat`, `w_cancel_ticket` and their comments; its schema adds `sale_event` and its
  indexes add four.
- **Answerability is computed from design flags and pinned to the catalogue by tests** in
  both studies, and study 03's `TestR06UnanswerableEverywhere` pins AM-02.0's finding.
- **Every dev-check log says what the progress log says.** Study 02: 15/15 designs pass on
  PostgreSQL and on YugabyteDB 1-node. Study 03: 13/13 on PostgreSQL, 14/14 on YugabyteDB.
  C1 overbooked every race event and H0 overbooked at every holds tier (`am02-dc01-pg/logs`).
- **Study 02's race cannot fall into AM-02.7's trap.** Buyers free-run (`RunRaces`); only the
  organiser's editor is paced. The buyer count is a real demand variable here.
- **Study 03's r01 binds `within` from the wall clock** (`time.Now().Add(holdsLookahead)`).
  That is a deliberate exception to RECENCY.md §2 rather than a slip: holds live in real
  time, and "what expires in the next ten minutes" has no historical regime. Accepted, and
  it must say so in `reports.go` (AM-03.9).

**The eleven defects.** The ones that matter most are 1–5: the review question was "is X1's
ledger audit real?", and on the evidence the answer was *not yet shown*.

1. **X1 records every refund without its buyer.** In `w_cancel_ticket`, `RETURNING` returns
   the row *after* the update (PostgreSQL 17, and YSQL on PostgreSQL 15; `RETURNING OLD`
   exists only from PostgreSQL 18), and the update has just set `customer_id = NULL`. Every
   `cancelled` ledger row therefore has a NULL customer. A refund ledger that cannot say who
   was refunded fails the purpose X1 was added for.
2. **The reconciliation audit counts; it does not attribute.** `a_ledger_mismatches`
   compares each seat's net count (sales minus cancellations) with its live status. It cannot
   see defect 1, a wrong customer, a wrong timestamp or price, or events in an impossible
   order.
3. **The audit never runs after the race or the churn race.** REPORTS.md §4 says "after
   every writing phase". `race.go` calls `RunAudit` per tier and never `RunLedgerAudit`, and
   the race is X1's headline experiment.
4. **The audit was never exercised after a write and never seen to fire.** X1 ran
   `-cmd verify` only. At load the loader seeds the ledger from the same `ds.Sold` as the
   tickets, so "consistent" there is guaranteed by construction and proves nothing.
5. **X1's ledger-answered reports were verified only where they cannot differ from P3.**
   At load no refund exists, so gross and net `r01` agree and `r05` is 0. A statement that
   returned a constant 0 would pass.
6. **Rendering.** The report prints only the load-time ledger audit, although write-phase
   audits are stored. The correctness table (`auditsCell`) has no ledger column.
   `LedgerAudit.String()` dereferences a nil receiver.
7. **Study 02's `r03` is declared answerable on the fourteen non-ledger designs**, but it
   suffers the same refund erasure as `r01`: when a customer's most recent purchase is
   refunded, the ticket-table formulation dates them by an earlier purchase or drops them.
   X1's `r03` comment explains the difference wrongly. Cancelling an *earlier* ticket
   changes neither formulation; only a refunded *latest* purchase does.
8. **Most report checks compare row counts, not values.** Study 02 `r02`, `r03`; study 03
   `r01` (the far cutoff never tests the predicate), `r03`, `r04`, `r05`. Study 02's `r03`
   truth is capped at 100, so both regimes probably check "100 rows". A statement that
   ignored its window would pass several of these.
9. **Plans for the new reports would not be captured.** Neither study's `ExplainAll` value
   map has `since`/`until`, and study 03's also lacks `within`. The explain phase would
   record `NOT CAPTURED: …` for `r01`/`r03`/`r05`, which breaks the EXPLAIN floor
   (non-negotiable 6). The dev checks never ran `explain`. **This is the same defect as
   LESSONS_LEARNED "A new query's parameter has more than one binding site", from this task's
   own phase 1. AM-02.4 should have cited it, and that omission is mine.**
10. **X1 skipped "Adding a design" steps 4–5.** There is no `p3_precreated_cas →
    x1_cas_ledger` entry in `pairs` (the reporter's controlled-pair section), and neither
    `diagrams/p_precreated.puml` nor the README's design table shows X1.
11. **The records overstate.** The progress entry credits the status-prefix fix to study 02
    (it is study 03's `r02`). CONTEXT.md says X1's ledger audit passes, which is true only at
    load, by construction. The `tiny` calibration covers one design per study, on YugabyteDB
    1-node only, with no explain and no client-throttling reading, so it cannot size a
    `small` run. Tags `study-02/v2-reports-devchecked` and `study-03/v2-reports-devchecked`
    therefore mark a state with these defects. They stay (tags never move), and the repaired
    state gets new tags (AM-03.13).

**AM-03.1 — strengthen the ledger audit first (X1 only; `sql/x1_cas_ledger/audit.sql`,
`harness/audit.go`).** Append this statement, with its comments kept in the file:

```sql
-- name: a_ledger_attribution
-- params: none
-- Per seat, the ledger must read as a history: sold, cancelled, sold, ... Each refund
-- names the buyer of the sale it reverses. The newest event of a seat that is sold
-- now is that sale, with the ticket's own buyer, time and price. Ordered by `at`,
-- then id: YSQL caches sequence values per connection, so BIGSERIAL order is not
-- commit order there.
WITH seq AS (
    SELECT event_id, seat_no, kind, customer_id, at, price_cents,
           LAG(kind)        OVER w AS prev_kind,
           LAG(customer_id) OVER w AS prev_customer,
           ROW_NUMBER() OVER (PARTITION BY event_id, seat_no
                              ORDER BY at DESC, sale_event_id DESC) AS recency
      FROM sale_event
    WINDOW w AS (PARTITION BY event_id, seat_no ORDER BY at, sale_event_id)
)
SELECT event_id, seat_no, 'refund without a preceding sale' AS problem
  FROM seq WHERE kind = 'cancelled' AND prev_kind IS DISTINCT FROM 'sold'
UNION ALL
SELECT event_id, seat_no, 'refund not attributed to the refunded buyer'
  FROM seq WHERE kind = 'cancelled' AND prev_kind = 'sold'
   AND (customer_id IS NULL OR customer_id <> prev_customer)
UNION ALL
SELECT event_id, seat_no, 'second sale with no refund between'
  FROM seq WHERE kind = 'sold' AND prev_kind = 'sold'
UNION ALL
SELECT t.event_id, t.seat_no, 'live sale disagrees with its newest ledger row'
  FROM ticket t
  JOIN seq l ON l.event_id = t.event_id AND l.seat_no = t.seat_no AND l.recency = 1
 WHERE t.status = 'sold'
   AND (l.kind <> 'sold' OR l.customer_id IS DISTINCT FROM t.customer_id
        OR l.at IS DISTINCT FROM t.sold_at OR l.price_cents <> t.price_cents);
```

`RunLedgerAudit` runs both statements. `LedgerAudit` gains `AttributionMismatches int64` and
up to `auditExamples` examples `{event_id, seat_no, problem}`. The gate and every printed
line use `Mismatches + AttributionMismatches`. `String()` becomes nil-safe and names both
counts. A live seat with no ledger row at all is still caught by `a_ledger_mismatches`.
**If the load-time audit reports attribution mismatches, first check that the loader writes
`ticket.sold_at` and `sale_event.at` through the same encoding (both through `CopyFrom`).
Fix the loader, never the audit.**

**AM-03.2 — the audit runs after every writing phase (study 02 harness).** In `race.go`,
after each tier's `RunAudit`, call `RunLedgerAudit` with phase `"<mode>@<tier>"`, store it
in a new `RaceResult.LedgerAudit`, and print it the way `main.go` prints write-phase ones.
Holds do not apply, because X1 has no holds. A post-write ledger inconsistency does not
abort the cell. In this study post-write audits are results, which is how C1 is seen to
fire. It is reported wherever a post-write audit violation is reported (AM-03.7).

**AM-03.3 — prove the audit fires, in this order, before touching `w_cancel_ticket`.**
Dev checks, `tiny`, `pg-single`, under the lock, into
`studies/02-ticket-booking/results/devchecks/am03-dc01-ledger-controls/`:

1. **Natural control, on the unfixed SQL:** X1 with `-cmd full -phases verify,write,churn
   -write-ops cancel -race-trials 1`. **Expected:** attribution mismatches > 0 after
   `cancel` and after every churn tier that refunded anything (defect 1). If it reports
   consistent, **stop: Escalation Required**. Either the audit or this review is wrong,
   and HIGH must know which.
2. **Injected faults:** add a harness flag `-ledger-fault none|drop-sale|wrong-customer`
   (default `none`). It is accepted only with `-cmd verify` on a `Ledger` design; anything
   else exits with an error. After load and before the audits, `drop-sale` deletes the
   `sold` row of the first sold seat of the busiest catalogue event, and `wrong-customer`
   adds 1 to that row's `customer_id`. The fault SQL is a labelled Go constant in the fault
   code, never in the design catalogue, so plans and the SQL-binding test never see it.
   **Expected:** `drop-sale` gives exactly 1 net mismatch, `wrong-customer` exactly 1
   attribution mismatch, and each run exits non-zero with the gate message. Also expected:
   `grep -- -ledger-fault run-study.sh` finds nothing.
3. Then AM-03.4, then repeat step 1. **Expected:** consistent after `cancel` and after
   every churn tier, with at least one refund present in the ledger.

**AM-03.4 — fix X1's refund attribution (`sql/x1_cas_ledger/writes.sql`).** The `UPDATE`
stays byte-identical to P3's. The buyer comes from the ledger row of the sale being reversed:

```sql
WITH cancelled AS (
    UPDATE ticket
       SET status = 'available', customer_id = NULL, sold_at = NULL
     WHERE ticket_id = $1
       AND status = 'sold'
    RETURNING event_id, seat_no, price_cents
), logged AS (
    INSERT INTO sale_event (event_id, seat_no, customer_id, kind, at, price_cents)
    SELECT c.event_id, c.seat_no,
           (SELECT s.customer_id
              FROM sale_event s
             WHERE s.event_id = c.event_id AND s.seat_no = c.seat_no AND s.kind = 'sold'
             ORDER BY s.at DESC, s.sale_event_id DESC
             LIMIT 1),
           'cancelled', now(), c.price_cents
      FROM cancelled c
    RETURNING event_id, seat_no
)
SELECT event_id, seat_no FROM logged;
```

The design comment must say why this is correct. The sale's own `sold` row committed in the
same statement as the sale, and a refund names a ticket its caller has seen sold: churn
cancels the ticket that buyer's booking just committed, and `BenchmarkCancel` cancels loaded
tickets whose rows the loader seeded. So that row is in the refund statement's snapshot.
Rejected alternatives: `RETURNING OLD` does not exist on either engine version here; a
`FOR UPDATE` CTE changes the cancel's locking; an extra column on `ticket` would make P3 → X1
two decisions. Also correct the `r03` comment (defect 7).

**AM-03.5 — post-write report checks (study 02 harness, all designs).** In `main.go`, after
the `cancel` write op and after each churn trial (after `RunRaces` returns, on that trial's
world), check these against the harness's own counters, with window
`[loadEpoch, time.Now()+24h)`:

- `r01` for the event with the most harness sales in that phase. Expected, where `r01` is
  `answerable`: `InitialSold + booked` (gross). Where it is `partial`:
  `InitialSold + booked − cancelled` (net).
- `r05` (only where answerable): `Σ cancelled` over every event.
- Tolerance `± Σ ambiguous`. Record them as `Check`s on the write/race result, print them,
  and report them like audits (AM-03.7).

A failure on X1 is a ledger correctness violation. **A failure on any other design is
Escalation Required before phase 3b** — do not adjust the expectation.

**AM-03.6 — declarations and verification strength (both studies).**

- **Study 02 `r03`: `partial` on every non-ledger design**, note: *"a customer whose most
  recent purchase was refunded is dated by an earlier purchase, or drops out, because
  cancellation erases the sale"*. HIGH decides the definition: `r03`, like `r01`, counts a
  purchase that was later refunded, because "when did this customer last buy" is a question
  about the purchase. Update `reports_test.go` if it asserts statuses. REPORTS.md §2 carries
  HIGH's dated note (added with this amendment).
- **Compare values wherever the truth is deterministic.** Replace these row-count checks:
  - study 02 `r02` and study 03 `r04`: compare the multiset of `sold_at` of the returned
    rows with the truth's top 50. Ties do not change a multiset.
  - study 02 `r03` and study 03 `r05`: per regime, compare the multiset of `last_at` with
    the truth's top 100 by recency.
  - study 03 `r03`: compare exact rows `(section_no, confirmed, revenue_cents)`. Check both
    the whole-span window and the trailing regime's window.
  - study 03 `r01`: keep the far-cutoff check and add `within = loadNow`. Truth: the number
    of `stExpired` seats of that event, or 0 when `d.ExpiryOnSweeper`. The loader's expiry
    margins keep this clear of clock skew.
  - Compare timestamps at microsecond precision on both sides. Use whatever helper keeps
    both sides in one format (permitted choice).
- **Make the recency trap visible.** In the historical regime, also compute the truth for
  the naive formulation (anyone who *purchased* in the window). If its top-100 multiset
  equals the recency one, print `WARNING: dataset cannot distinguish last-purchase from
  purchased-in-window at top-100` and record it on the check. This does not fail the gate,
  but the gate's blind spot must be visible.

**AM-03.7 — plans, rendering and design registration.**

- `ExplainAll`, both studies: add `since`/`until` = `reportWindowFor("historical")`, and in
  study 03 add `within` = `time.Now().Add(holdsLookahead)`. **Expected:** no `NOT CAPTURED`
  for any `r0x` statement in any design's plans file.
- Study 02 reporter: render every phase's ledger audit (write ops, race tiers, churn tiers).
  Add ledger and post-write report-check violations to the correctness table beside the
  existing audit cells, and wherever the TL;DR's fixed rules list audit violations.
- `pairs`: add `{"p3_precreated_cas", "x1_cas_ledger", "What does an append-only sale ledger
  in the selling statement cost, and what does it make answerable?"}`.
- X1 in `diagrams/p_precreated.puml` (re-render SVG through the container renderer) and in
  the README's design table and family list.

**AM-03.8 — dev checks after repair, `tiny`, under the lock.** Into
`results/devchecks/am03-dc02-{pg,yb}/` of each study:

| Study | Topology | Designs | `-phases` (and flags) |
|---|---|---|---|
| 02 | pg-single | all 15 | `verify,explain,read,write,churn` `-write-ops cancel -race-trials 1` |
| 02 | pg-single | C1 | `verify,race` (must overbook) |
| 02 | pg-single | H0 | `verify,holds` (must overbook) |
| 02 | yb-single | all 15 | `verify,explain` |
| 02 | yb-single | P3, X1 | `verify,explain,read,write,race,churn` `-write-ops book,cancel -race-trials 1` |
| 03 | pg-single | all 13 | `verify,explain` |
| 03 | yb-single | all 14 | `verify,explain` |
| 03 | both | S1 | `verify,explain,read` |

**Pass:** every gate passes; no `NOT CAPTURED` in any plans file; X1's ledger audits are
consistent after every writing phase on both engines; every post-write report check passes;
C1 and H0 fire; both studies' unit tests pass (`go test -c`, run, delete). Generate one report
per study from these directories and confirm the ledger lines and the P3 → X1 pair render.
Any other outcome: fix what is mapped here, otherwise Escalation Required.

**AM-03.9 — records.** Study 03 `reports.go`: a comment stating r01's wall-clock exception
(AM-03.0). LOW writes its own progress entry and does not edit iteration 3's (the
corrections are recorded in HIGH iteration 5). Update CONTEXT.md, and add LESSONS_LEARNED
entries for any new failure met.

**AM-03.10 — `small` calibration (dev checks, not measured runs).** Under the lock, into
`studies/0N-*/results/devchecks/am03-cal-small/`, `-scale small -duration 5s -warmup 2s
-conns 8`:

| Study | Topology | Designs | `-phases` (and flags) | What it sizes |
|---|---|---|---|---|
| 02 | yb-cluster3 | P3, X1 | `verify,explain,read` | report sampling (P3 scans 289 500 ticket rows; X1 the ledger) |
| 03 | yb-cluster3 | S1, L2 | `verify,explain,read` | report sampling (L2 unnests section documents) |
| 02 | pg-single | P3, X1 | `verify,race -race-trials 1 -race-buyers 128` | generator capacity at the high buyer count |

Record in a `SUMMARY.md`: each cell's wall time (timestamps around the `podman run`), load
time, every report's ops/s and executions (ops/s × 5), client CPU throttling for the `read`,
`reports` and `race#1` probes, and the race's error and attempts-per-sale columns.

**AM-03.11 — sizing rules. Apply them mechanically, and write the chosen values into the
progress log before any measured run.**

- **R1 duration, per study:** `5s` (the published matrices' value, which keeps the
  re-measured read questions comparable) if every report has ≥ 100 executions in the
  calibration. Otherwise `10s` for that study's reports run. Never above 10 s. A report still
  under 50 executions at 10 s is a result, not a blocker: the analysis names it under
  weaknesses.
- **R2 high buyer count B2:** `128` if the client was CPU-throttled in ≤ 5% of CFS periods
  during `race#1` and the race logged no connection or pool errors. Otherwise re-run that
  calibration at `64` with the same criterion, and B2 = `64`. If 64 fails too, **Escalation
  Required**.
- **R3 wall-time guard:** estimate every run below from the calibration's cell wall times.
  My estimate before calibration, from the 2026-09-13 and 2026-09-15 matrices (study 02
  averaged 9.7 min per full cell; YugabyteDB loads 21–33 s; a P3 race trial on YugabyteDB
  about 3.5 min): 3b-1 ≈ 2.5 h, 3b-2 + 3b-3 ≈ 2.5 h, 3b-4 ≈ 2.5 h. **If the calibrated total
  exceeds 10 h, Escalation Required** with the estimate. Do not drop designs, topologies or
  trials to fit.

**AM-03.12 — phase 3b: four measured runs, authorised once AM-03.8 passes and AM-03.11 is
applied.** Each run is started from a committed tree with the runner's `--tag`, one at a time,
announced in CONTEXT.md ("running — do not start databases"), followed by `-cmd report`, a
commit of results and report, and nothing else. No analysis.

| Run | Study | Command (runner flags) | Question |
|---|---|---|---|
| 3b-1 | 02 | `./run-study.sh --scale small --duration <R1> --phases verify,explain,read --tag` (all 15 designs, three topologies) | what each design's answerable reports cost, beside the five buyer reads re-measured in the same session |
| 3b-2 | 02 | `--scale small --designs p3_precreated_cas,x1_cas_ledger --phases verify,write,race,churn --extra "-write-ops book,cancel -race-trials 3 -race-buyers 32" --tag` | the ledger's cost at the published buyer count |
| 3b-3 | 02 | same, with `--designs x1_cas_ledger,p3_precreated_cas` and `-race-buyers <B2>` | the ledger's cost when contention is 4× (or 2×); the reversed design order lets the analysis see order effects |
| 3b-4 | 03 | `./run-study.sh --scale small --duration <R1> --phases verify,explain,read --tag` (all 14 designs, three topologies) | what each layout's reports cost |

Keep the runners' default topologies (`pg-single,yb-single,yb-cluster3`) and study 02's
YugabyteDB harness flags. **Stop conditions:** a failed gate on P3 or X1, any ledger audit
inconsistency, or any post-write report check failure in 3b-2/3b-3 means finish the run,
commit it, **stop, Escalation Required**. Do not re-run anything. In 3b-1/3b-4, a failed cell
is reported as the runners already do. C2's known YugabyteDB failure is in the race phases,
which these runs do not include; if it appears anyway, report it and continue.

**AM-03.13 — tags.** After AM-03.8: `study-02/v2.1-reports-repaired` and
`study-03/v2.1-reports-repaired`. After each run: the runner's `run/…` tag. After 3b-4:
`study-02/v2-measured` and `study-03/v2-measured` on the commit holding the last report.

**AM-03.14 — what not to do.** Do not change P3 or any other existing design, schema, index
or write statement (X1's cancel is the only write change). Do not change race, churn, holds
or lifecycle mechanics. Audits and checks *after* a phase are instrumentation and are
allowed. Do not add a hold-history design. Do not write analyses. Do not move tags. Do not
re-run the published matrices.

**Next after execution: HIGH — review AM-03.3's control evidence and the calibration, then
write signed analyses for 3b-1 to 3b-4 (study 02: reports and the P3 → X1 ledger pair;
study 03: reports).**

### AM-04 — Execution Handoff Update: RR-ER-01 decided (order-free attribution)

- Decided by: Claude Opus 5 (HIGH), Claude Code desktop, 2026-09-20. Next setting: **LOW — Claude Sonnet 5.**
- Delta only. Everything in AM-03 not named here stands unchanged.

**The decision: stop reconstructing order. Read the buyer from the row being cancelled.**
LOW's evidence settles it — `at` is transaction-start time, `sale_event_id` is YSQL-cached,
and neither reconstructs commit order on both engines. Both are the wrong tool: the refund
already holds the only authoritative record of who is being refunded, which is the ticket
row it is about to clear. LOW's recommended option (a) is adopted.

**AM-04.1 — `w_cancel_ticket` (X1 only) becomes order-free.** Replace the correlated
subquery with a locking read of the ticket row in the same statement:

```sql
WITH old AS (
    SELECT ticket_id, customer_id
      FROM ticket
     WHERE ticket_id = $1 AND status = 'sold'
     FOR UPDATE
), cancelled AS (
    UPDATE ticket t
       SET status = 'available', customer_id = NULL, sold_at = NULL
      FROM old
     WHERE t.ticket_id = old.ticket_id
    RETURNING t.event_id, t.seat_no, old.customer_id, t.price_cents
), logged AS (
    INSERT INTO sale_event (event_id, seat_no, customer_id, kind, at, price_cents)
    SELECT event_id, seat_no, customer_id, 'cancelled', now(), price_cents FROM cancelled
    RETURNING event_id, seat_no
)
SELECT event_id, seat_no FROM logged;
```

`FOR UPDATE` re-checks `status='sold'` against the latest row version under READ COMMITTED
and returns that version's buyer, so `old.customer_id` is exactly the buyer being refunded —
no timestamp, no sequence, nothing to reorder. **Fallback if YugabyteDB rejects a locking
read inside a CTE:** do the locking read as a separate statement (`w_cancel_lookup_buyer`)
inside `Booker.Cancel`'s existing explicit transaction, behind `if d.Ledger`. Record which
form you used as a Decision Log entry; both are acceptable, the CTE is preferred.

**AM-04.2 — I am revising AM-03.4's rejection of `FOR UPDATE`, and the controlled-pair
claim with it.** I rejected it as "changes the cancel's locking"; that was wrong on the
evidence. The `UPDATE` takes the same exclusive lock on the same single row microseconds
later anyway, so the lock held, its duration and its ordering are unchanged — what is added
is one extra index lookup and lock acquisition, which on YugabyteDB is a real round trip.
That cost **is part of the ledger's cost**, not a second decision: an append-only ledger
that records who was refunded has to read who was refunded. State the pair precisely
wherever it is described (X1's own comments, `REPORTS.md` §3, and any later analysis):

> **P3 → X1, write axis:** the sell path is P3's, plus one ledger insert in the same
> statement. The refund path is P3's, plus one ledger insert and the locking read that
> insert needs to name the refunded buyer. The sell-out race — the headline experiment —
> exercises only the first.

**AM-04.3 — `a_ledger_attribution` becomes order-free too.** Replace it entirely (no
`LAG`, no `ROW_NUMBER`, no ordering). Keep `a_ledger_mismatches` exactly as it is. The
replacement checks, per seat:

- **B — no buyer is refunded more than they bought:** for each `(event_id, seat_no,
  customer_id)`, `count(kind='cancelled') <= count(kind='sold')`. This is what catches a
  refund named to the wrong buyer: LOW's own bad trace had buyer 604 with one sale and two
  refunds of seat 87/10.
- **C — the live sale is recorded verbatim:** for every ticket with `status='sold'`, a
  `'sold'` row must exist for that seat with `customer_id`, `at` and `price_cents` equal to
  the ticket's `customer_id`, `sold_at` and `price_cents`; and for that buyer and seat,
  `count(cancelled) < count(sold)` (a buyer holding the seat cannot have refunded every
  sale they made of it).

Drop "second sale with no refund between": one row per seat plus the CAS predicate makes
a double sale structurally impossible, and `a_duplicate_seats` already checks it. Keep the
`problem` text column and the existing `LedgerAttributionIssue` plumbing.

**AM-04.4 — what to re-check, and nothing more.** Only X1's SQL changed. On `pg-single`
then `yb-single`, `tiny`: X1 with `-phases verify,explain,read,write,race,churn
-write-ops book,cancel -race-trials 2`, plus **three** high-contention repeats per engine
(`-phases verify,churn -churn-pct 50 -race-buyers 64`) — the setting that exposed the bug.
Into `results/devchecks/am04-dc01-{pg,yb}/`. **Pass:** gate passes, every ledger audit
consistent including at high contention on *both* engines, post-write report checks pass,
no `NOT CAPTURED`. Also re-run the `-ledger-fault drop-sale` and `wrong-customer` controls
once on `pg-single` and confirm the rewritten audit still fires on each. Do **not** re-run
the other 14 designs or study 03 — they are untouched and already dev-checked (`6b1926b`).

**AM-04.5 — then continue straight through AM-03.10 → .13 without returning to HIGH.**
Calibration, the sizing rules, the four measured runs, the tags. Under the 2026-09-17
rules: decide, log, continue. Escalate only on a stop condition below. Add
`study-02/v2.2-ledger-attribution` to AM-03.13's tag list, on the commit that passes
AM-04.4.

**AM-04.6 — stop conditions (unchanged in spirit, restated).** Stop and escalate only if:
X1's ledger audit is still inconsistent after AM-04.1+.3 on either engine; a post-write
report check fails on any design; the calibration projects over 10 h total; or a measured
run fails a correctness gate on P3 or X1. A slow report, an imperfect-confidence
implementation choice, or a YugabyteDB cell failing the way C2 already does are **not**
stop conditions — log and continue.

**The lesson worth keeping** (LOW to add to `LESSONS_LEARNED.md` in its own words): neither
a transaction-start timestamp nor a sequence value reconstructs commit order under
concurrency — `now()` is fixed at BEGIN, and YSQL hands out sequence blocks per connection.
An audit that needs to know "which record came first" should be rewritten so it does not
need to know.

**NEXT MODEL: LOW**
