# Progress — recency question and operational reports (EH-02)

Step-by-step state of this task. HIGH writes its planning iterations here; LOW appends
its execution log. Nobody edits another iteration's entry.

## HIGH iteration 1 — 2026-09-15

- Planner: Claude Opus 5, setting `ultracode` (HIGH role), Claude Code desktop app.
- Created the task worktree `.worktrees/recency-reports`, branch
  `repo/recency-and-reports`, from local `main` at `7123da6` (clean). One branch spans
  the three studies because the owner asked for one change: a new class of question and
  the reports that go with it.
- Checked before touching anything shared: `git worktree list` showed only
  `.worktrees/study01-v3` and `.worktrees/study01-review-audit` besides `main`; the
  podman volume `ads-run-lock` did not exist, so no benchmark was running; `CONTEXT.md`
  reports study 03 complete and nothing in flight.
- Read the owner's request, study 01's catalogue, designs, harness, loader, verifier and
  runner, study 02's and study 03's query catalogues and write paths, the methodology,
  and the required-comparison rules.
- Published three protocols:
  - [`studies/01-charity-tree/RECENCY.md`](../../../studies/01-charity-tree/RECENCY.md) —
    the owner's "who donated last in a period" question, its two window regimes, seven
    new designs (including the owner's flag proposal as D20 and its unguarded twin D21
    as the negative control), controlled pairs, correctness, four experiment groups,
    sizing and limitations.
  - [`studies/02-ticket-booking/REPORTS.md`](../../../studies/02-ticket-booking/REPORTS.md)
    — six operational reports, the answerability finding (no design can report refunds,
    because cancelling erases the sale), and one new design X1 that adds an append-only
    sale ledger to P3.
  - [`studies/03-reserved-seating/REPORTS.md`](../../../studies/03-reserved-seating/REPORTS.md)
    — six operational reports, the answerability finding (expiry erases the hold, so the
    abandonment funnel cannot be computed), and the reasoned decision **not** to add a
    hold-history design in this iteration.
- Published [EH-02 v1](HANDOFF.md), covering **phase 1 only**: study 01's implementation
  and dev checks up to a `small` calibration. Phases 2 and 3 are deliberately unauthorised
  until HIGH has the calibration, because their sizing depends on it.
- Decisions worth recording here because they shape everything downstream:
  - **D1–D17 keep byte-identical schema, indexes and write paths.** Only four read
    statements are appended to each catalogue, so every published storage and write
    number for study 01 stays valid.
  - **The legacy twelve-question geometric mean stays defined over those twelve.** A unit
    test protects it. Extending the score would have silently rewritten the meaning of
    every published read score in the study.
  - **Embedding gets no new design**, with the reason written down: an expression index on
    the document's text timestamp is impossible (`STABLE` cast) or unsound, so D6/D9/D10
    are measured on the new statements as they stand, and "embedding buys nothing for a
    cross-parent recency question" becomes a measured claim rather than an assumption.
  - **The window is derived from the dataset, never from `now()`**, so two runs stay
    comparable as the calendar moves.
  - **The study 01 README's exclusion of period queries is lifted for this question only**
    (the owner asked for it); materialised views and time-bucketed tables stay out.
- No database, container, image build or benchmark was started. No lock was taken. Local
  `main` is untouched.

**Next: LOW — Claude Opus 5, setting `high` — execute EH-02 phase 1, steps 1–21.**

## LOW iteration 1 — 2026-09-15/16

- Executor: Claude Sonnet 5, Claude Code desktop app, model set to `claude-sonnet-5`
  (the owner's "lower model" for this session).
- Preflight (steps 1–3): confirmed clean checkout on `repo/recency-and-reports` at
  `4c6b607`; no other worktree running a benchmark; no lock held; read RECENCY.md in
  full before touching code.
- **Step 4** — appended q13–q16 to all seventeen existing `queries.sql` files, in the
  formulation the handoff mapped for each (A/B by whether `charity_id` is on `donation`,
  C for the rollup designs, D for D6's whole document, E for D9/D10's bounded cache).
  Verified by diff: 864 insertions, zero deletions, across exactly those 17 files.
- **Step 5** — created D18–D24 by copying the mapped parent directory and making exactly
  the one change; every `diff -r` check the handoff specifies passed. (D4↔D22's
  `queries.sql` is byte-identical, not merely "differs only in indexes.sql and
  queries.sql" as the handoff's check line literally said — D4 is itself one of the
  formulation-C designs, so D22 shares its query text by construction. Noted here as a
  clerical imprecision in the handoff's own check wording, not an escalation: the actual
  required SQL content, from D-2's canonical statements, was unambiguous.)
- **Step 6** — wrote D20's `triggers.sql` (person-row lock, insert and delete
  maintenance) and D21's as D20's with exactly the lock statement removed; `diff` shows
  only that one hunk.
- **Steps 7–11** — harness changes: `designs.go` (seven new `Design` entries, `LastFlag`,
  `RecencyRollupIdx`), `load.go` (bulk-loads the flag from `DonationsByPerson`'s own
  order), `bench.go` (window regimes, `insert_backdated`, both regimes measured for
  every recency query), `verify.go` (eight new checks against truth computed
  independently, with the same boundary-tie rule the existing gate uses; pure
  `recencyMatches`/`recencyExpectedMatches`/`recencyExpectedLapsedCount` extracted for
  unit testing), `audit.go` (`AuditRecency`, `diagnoseFlagMismatch`), `main.go`
  (wiring, `coreQueries` fixing the legacy geomean to exactly the original twelve),
  `report.go` (a separate "Recency" section, design order/labels, nine new pairs).
- **Step 12** — added `TestLegacyReadScoreExcludesRecencyQueries`,
  `TestRecencyWindowRegimesDeriveFromTheDataset`,
  `TestRecencyExpectedMatchesAgainstHandBuiltDataset`,
  `TestDiagnoseFlagMismatchClassifiesBySymptom`; updated
  `TestMechanismControlsKeepUnchangedQueriesIdentical` for the sixteen-statement
  catalogue (q15 joins the "changed" set at the D12→D13 transition, for the same reason
  q02/q05/q12 already do). `go test ./...` required building a `-c` binary and running
  it directly — this machine's Application Control policy blocks executables `go test`
  builds into `%TEMP%`. All nine tests pass.
- **Step 13** — `go vet` clean. `gofmt -l` initially flagged every file in the module;
  confirmed this is the repo's pre-existing Windows CRLF checkout, not a real issue
  (LF-normalized copies of every harness file gofmt clean), except two genuine
  misalignments in `audit.go`/`main.go` from my own edits, fixed with `gofmt -w`.
- **Step 14** — `bash -n` clean on both modified runners.
- **Step 15** — `d20_recency_flag.puml` and `d22_recency_rollup_idx.puml` written and
  rendered (`diagrams/render.sh svg`, 20 diagrams total, only the two new SVGs changed);
  README's questions and designs tables extended.
- **Step 16** — committed in two groups (SQL catalogue; harness+runner+diagrams+README),
  each ending with the required attribution trailer.
- **Step 17** — dev checks, `pg-single`, `tiny`, under the benchmark lock:
  - `-cmd verify` on all 22 PostgreSQL-capable designs: 22/22 passed, 20/20 checks each,
    **after** finding and fixing a real bug — D6/D9/D10 initially returned zero rows for
    every recency question because the embedded document's on-disk timestamp key is `t`,
    not `donated_at` (see `LESSONS_LEARNED.md`). Fixed, rebuilt, reverified: all three
    pass.
  - `-cmd full -write-ops insert,insert_backdated,delete,update` for D20–D23: all
    completed; D20/D22/D23 stayed consistent through every op; **D21 (the negative
    control) failed already here**, under ordinary 8-connection spread demand — 3 of
    500 donors with more than one flagged donation, examples recorded. D23's
    delete/update were correctly skipped (D5's existing rule).
  - Hot-donor contention (the "arrival" experiment, 8 then 16 writers, all targeting one
    donor): found a second gap — the arrival experiment's post-warmup audit point ran
    `AuditRollups` but never `AuditRecency`, so `warmup_recency_audit` was simply absent
    from every result. Wired it in (`arrival.go`). At 8 writers D21 did not fail (see
    `LESSONS_LEARNED.md` on why a paced open-loop arrival experiment is a gentler test
    than a closed-loop one); at 16 writers, the handoff's mapped fallback, it did —
    1 mismatch, example recorded — while D20 stayed consistent under the identical
    contention. Acceptance criterion 4 ("D21 was seen to fail... with examples") is
    satisfied twice over, by different mechanisms.
  - A third gap, found while running D20–D23 with `-cmd full`: `ExplainAll` (plan
    capture) has its own fixed parameter map, separate from the read benchmark's, and it
    never had `since`/`until`. Fixed (bound to the trailing regime, since EXPLAIN is
    about the access path, not which regime is correct).
- **Step 18** — repeated on `yb-single` then `yb-cluster3` (torn down between): all 24
  designs (including D7 and D24, YugabyteDB-only), 20/20 checks each, both topologies.
  D20's **partial** index (`WHERE is_last_donation`) was accepted by YugabyteDB without
  needing the mapped fallback to a plain composite index.
- **Step 19** — calibration: one `small` cell, `pg-single`, D3, the four recency
  questions restricted via `-queries`, both regimes, 3 trials. q13/q14/q16 measured
  160–200 ops/s, q15 (charity-scoped) 400–450 ops/s — an order of magnitude below the
  study's typical indexed reads, expected since D3 has no supporting index for this
  question. Client CPU sampled via `podman stats` at ~4% of its 2-CPU budget (not the
  bottleneck). Wall-clock ~90s for the whole cell (8 statements × 3 trials × (3s+1s));
  reads are time-boxed by duration regardless of throughput, so no cell-duration
  blowup is expected from the full recency-reads matrix (~24 cells × ~90s ≈ 36 min,
  well under the 12-hour ER trigger).
  - All three fixes found during steps 17–19 (the `t` key, `ExplainAll`'s missing
    params, `arrival.go`'s missing audit wiring) were committed together with the
    dev-check evidence they were found by.
- Every database container was torn down between topologies and at the end; the
  benchmark lock was released; only a pre-existing, unrelated `pg-scratch` container
  (not created by this session) remains on the host.
- No escalations were raised. The one instance of the handoff's own check wording not
  quite matching its intent (step 5, D4↔D22) is recorded above as a clerical note, not
  an escalation, since the required SQL content was unambiguous from D-2.

**Next: HIGH — review this evidence, size the study-01 recency matrix (phase 2) and the
studies 02/03 report queries (phase 3) from the calibration above, and publish an
amendment authorising them.**

## HIGH iteration 2 — 2026-09-16 (validation of phase 1)

- Reviewer: Claude Opus 5, setting `ultracode` (HIGH role), Claude Code desktop.
- Reviewed LOW's phase 1 at `129d431` / `study-01/v4-harness`. **Accepted.** All eight
  acceptance conditions in the handoff are met, and the evidence is real rather than
  asserted: the gate proved itself by catching three genuine bugs, and D21's negative
  control was seen to fail under two different load shapes with examples recorded. LOW's
  finding about *why* the paced arrival experiment needed 16 writers where the
  closed-loop benchmark needed 8 is a genuine methodological result and is recorded in
  `LESSONS_LEARNED.md` where it belongs.
- I re-derived the claims rather than taking them on trust: checked the append-only diff
  property, the one-decision `diff -r` relationships, the formulation-to-design mapping,
  the window regimes' half-open interval and sentinel handling, and the tie rule in
  `checkRanked`.
- **Validation found four defects and one untested risk, none of them visible from the
  dev checks LOW ran, because they live in the path a *reported* run takes rather than
  the path a dev check takes.** They are decided in [AM-01](HANDOFF.md#amendments):
  1. `recency-maintenance` passes four write ops in one cell; `experiment.go` refuses
     more than one, so every cell in that group would have failed immediately.
  2. `experiment.go`'s `writes` and `arrival` cases never call `AuditRecency`, so the
     maintenance group would have produced no recency audit at all.
  3. `experimentProblems` does not consider the recency audits, so a fired negative
     control would not be flagged in the generated enhancements report — a methodology
     5a violation in the reporting path, and the most important of the five repairs.
  4. The recency reporting section LOW built lives in `report.go`, but every recency
     group runs through `run-enhancements.sh`, which generates its report with
     `report_enhancements.go`. **This is my specification error:** D-5 named a reporter
     without checking which one the runner calls.
  5. `delete_person` — ER trigger 5's own subject — was never exercised on the flag
     designs, where the delete trigger fires once per child row during a cascade.
- Two things I checked that turned out **not** to be problems, recorded so the next
  reader does not re-investigate them: `experimentReadScore` already skips `@`-suffixed
  queries and requires exactly twelve, so the new questions cannot contaminate the
  enhancements read score and a recency-only cell contributes no score rather than a
  zero; and the `@` suffix LOW chose coincides with an existing convention in that
  reporter (`@largest`/`@smallest` for targeted re-runs) without colliding with it.
- Phase 2 is sized in AM-01.6 at roughly four hours at three trials, inside ER trigger
  4's bound. Phase 3 remains unauthorised.

**Next: LOW — Claude Opus 5, setting `high` — execute AM-01.1 through AM-01.6.**

## LOW iteration 2 — 2026-09-16 (AM-01 repairs and phase 2)

- Executor: Claude Sonnet 5, Claude Code desktop app.
- **AM-01.1** — `run-enhancements.sh`'s `recency-maintenance` group now loops
  `for op in insert insert_backdated delete update delete_person`, one cell per op,
  matching the `exceptions` group's existing shape. `delete_person` was added to the
  group's operation list at this point, ahead of AM-01.5's dev check, since AM-01.1's
  fix already had to enumerate every operation individually.
- **AM-01.2** — `experiment.go`'s `writes` and `arrival` cases now call `AuditRecency`
  for `LastFlag`/`RecencyRollupIdx` designs, storing the result on `run.RecencyAudit`
  (and `run.Writes[0].RecencyAudit` where the rollup audit is attached that way too).
- **AM-01.3** — `experimentProblems` (`report_enhancements.go`) now flags
  `r.RecencyAudit` and `Arrival.WarmupRecencyAudit` mismatches, using the existing
  `recencyAuditBad` helper from phase 1.
- **AM-01.4** — added a per-trial recency-audit line to the per-group report output, and
  a new `recencyAuditSummary` helper that names D21 explicitly and states plainly
  whether it fired, never wording another design's mismatch as a control. Confirmed the
  generic per-operation metrics table already lists q13-q16 with their regime spelled
  out via the existing `@trailing`/`@historical` suffix, so no separate table was
  needed for that half of the requirement. New unit test:
  `TestRecencyAuditSummaryNamesD21AsTheControl`.
  - All ten harness tests pass; `go vet` and `gofmt` clean (checked by LF-normalizing
    the CRLF checkout before diffing, same as phase 1).
  - Committed at `ff79475`.
- **AM-01.5** — dev-checked `delete_person` on D20 and D21 at `tiny`, `pg-single`,
  directly via `-cmd full -write-ops delete_person` (first attempt used the default read
  duration and ran the full twelve-plus-recency read catalogue unnecessarily, making it
  look slow; a second attempt with `-duration 1s` on the reads isolated the actual
  question). D20: 6,561 delete_person/s, 0 errors, recency audit consistent (236 donors
  checked — the finite-budget benchmark had already consumed half the pool). D21:
  6,116/s, 0 errors, also consistent — `delete_person` walks the person id space
  monotonically with no concurrent contention on the same donor, so a clean audit here
  says nothing about the guard either way; it only confirms the trigger's per-child-row
  firing during a cascade does not deadlock or error. Neither cell approached the
  five-minute stop condition. No escalation. Committed at `2c6f515`.
- **AM-01.6** — ran phase 2 through the proper runner (`run-study.sh --suite
  enhancements --experiments recency-reads,recency-maintenance,recency-hot-donor,
  recency-placement --trials 3 --tag`), under the benchmark lock, recorded as running in
  `CONTEXT.md` first. First launch attempt failed at the runner's own dirty-tree check
  because a `tee`d log file had been written inside `studies/01-charity-tree/`; removed
  it and relaunched from a clean tree with output redirected outside the study
  directory. **192 cells, 0 failed processes**, completed well under the ~4-hour
  estimate (small scale, short per-query durations). Auto-tagged
  `run/01-charity-tree/20260916T090036Z-v3` by the runner itself.
  - **The AM-01 repairs held under the real run**: D21's negative control fired in
    every `recency-maintenance` trial (4, 7, 3 of 5000 donors) and in the
    `hot-donor-w16` sweep's third trial (1 of 5000) — consistent with the dev-check
    evidence — and every firing is visible in the generated report, worded exactly as
    AM-01.4 specified.
  - **One genuine, unflagged finding**: D23 (application-maintained rollup, optimistic
    CAS) recorded write errors under the hottest single-donor contention — up to 18 of
    10,000 at `w8`/`w16` — `first_error: "statement affected no rows"`, i.e. the CAS
    comparison lost the race. Recorded in the commit message as a measurement for
    HIGH's analysis to characterise, not suppressed or treated as a defect: it is
    exactly the abort-under-contention cost optimistic concurrency trades for not
    locking, and study 01 now has it measured for the first time under this exact
    contention shape.
  - Report: [`reports/20260916T090036Z-v3.md`](../../../studies/01-charity-tree/reports/20260916T090036Z-v3.md).
    Results: `results/20260916T090036Z-v3/`. Committed at `42a2f32`.
- Lock released and every database torn down after every step; verified via
  `podman volume inspect ads-run-lock` and `podman ps -a` before and after.
- No escalations raised in this iteration.

**Next: HIGH — review the AM-01 repairs and the phase-2 measured matrix, check the
D23 contention finding, and write the signed analysis.**

## HIGH iteration 3 — 2026-09-16 (validation of phase 2, and the signed analysis)

- Reviewer/analyst: Claude Opus 5, setting `ultracode` (HIGH role), Claude Code desktop.
- **AM-01 repairs accepted.** I verified each one against the artefacts rather than the
  report of them: the maintenance group emitted one cell per operation (five ops × four
  designs × three trials, with D23's three rollup-sensitive ops correctly skipped by the
  harness's existing rule); `AuditRecency` ran and is recorded in the experiment path's
  results; `experimentProblems` flagged all four D21 firings; and the report's per-trial
  lines name D21 as the control and state plainly whether it fired.
- **Phase 2 verified independently of LOW's summary.** Digest `1b05f142fd9e06ef` over 192
  result files. **All 192 cells passed their correctness gate — 22 checks each, zero
  failures anywhere in the run.** I classified all 22 "FAILED / diagnostic" cells: four
  are D21's control firing by design, nine are D23's deliberately untimed operations, nine
  are contention cells whose errors and queue rejections are themselves the measurement.
  No correct design failed anything.
- **Two things LOW's summary missed, both of which changed the reading:**
  - The hot-donor sweep is **rate-limited below capacity** — every healthy cell completed
    exactly 10,000 at 500.0/s, the offered rate, at every writer count. It measures
    compliance, not capacity, so it cannot rank the designs that kept up. D23's collapse
    is real precisely because it is a failure to keep up (2,225–5,569 requests rejected,
    23,468–28,321 retries). Recorded as a lesson and as a "measure next".
  - The D20 → D22 pair is **clean on the read axis and confounded on the write axis**:
    D22 carries D4's entire five-column, two-table rollup package, so its insert deficit
    against D20's one boolean does not price the flag-versus-rollup decision. Also
    recorded as a lesson; the missing design is a `last_donation_at`-only rollup.
  - Both are my own specification errors, and the analysis says so in those words.
- I also recovered the one comparison the reporting correctly suppresses: D21's insert
  cells are excluded from rate summaries *because* the control fires in them, which is
  exactly the guard-cost comparison. Pulled per-trial from the JSON: D20 6,345/5,846/6,219
  against D21 1,278/6,039/5,819 — indistinguishable, so **the guard is free**.
- **Signed analysis published:**
  [`20260916T090036Z-v3--claude-opus-5--2026-09-16`](../../../studies/01-charity-tree/reports/analyses/20260916T090036Z-v3--claude-opus-5--2026-09-16.md).
  Every figure in it was recomputed from the report or the raw results before committing;
  two were wrong on first writing (a 39x that is 33x, and a placement spread stated as
  2–15% that is 1.5–20%) and were corrected. The report was regenerated with the run's own
  pinned image (`d99fd495d7b3…`, still present) so the analyses index resolves; the diff is
  three lines of index and no measurement.
- Study README and `CONTEXT.md` updated; two new entries in `LESSONS_LEARNED.md`.

**Next: HIGH — decide phase 3 (studies 02 and 03's operational reports), which remains
unauthorised, or close the task at the owner's direction.**

## HIGH iteration 4 — 2026-09-16 (phase 3a planned)

- Planner: Claude Opus 5, setting `ultracode` (HIGH role), Claude Code desktop.
- The owner said continue, so this plans the second half of their 2026-09-15 request:
  the operational reports of studies 02 and 03.
- **First I verified the two claims I had made in those protocols on 2026-09-15 after
  reading only one design each.** Both hold, and are now established across all 28
  designs rather than sampled:
  - Study 02: cancellation destroys the sale in **every** design — C1–C5, R1–R3, H0, H1
    `DELETE FROM ticket`; P1–P4 reset the row and null out `customer_id`/`sold_at`. The
    designs' own comments say so. No design can report refunds, and `r01`'s revenue is
    wrong across a refund in all fourteen.
  - Study 03: **no** design retains hold history — zero hold-event tables, zero
    `expired_at`/`released_at` columns, and `w_release_expired` clears the seat row. The
    abandonment funnel cannot be computed from any design's state.
  - These are findings obtained by reading schemas; they need no benchmark to be true.
    What needs measuring is what the answerable reports cost and what making the
    unanswerable ones answerable costs on the hot path.
- I also checked the harness shape before specifying work: both studies build on
  `platform/`, and their read helpers bind exactly one parameter per query
  (`db.Query(ctx, st.SQL, key(r))`), while the new reports need up to three. AM-02.4 asks
  for a multi-parameter variant rather than a reshape, so the existing read questions keep
  their exact call path.
- **[AM-02](HANDOFF.md#amendments) authorises phase 3a: implementation and dev checks for
  both studies, no measured run.** Key decisions: reports appended with zero change to any
  existing statement, schema, index or write path; answerability declared in Go **and
  pinned to the SQL by a unit test**, because a declaration that can drift from the SQL is
  worth little; study 02 gains exactly one new design (X1, P3 plus an append-only sale
  ledger) with a ledger reconciliation audit; study 03 gains no new design, as its own
  protocol already reasoned.
- Both lessons from phase 2 are carried into the plan explicitly: X1 was checked on **both**
  axes before being registered as a controlled pair (clean on the write axis, deliberately
  different on the read axis, which is what it buys), and AM-02.7 pre-commits phase 3b's
  race sizing so a fixed offered load cannot produce a second null result.

**Next: LOW — Claude Opus 5, setting `high` — execute AM-02.1 through AM-02.5.**

## LOW iteration 3 — 2026-09-16 (AM-02, studies 02/03 operational reports)

- Executor: Claude Sonnet 5, Claude Code desktop app.
- **AM-02.1 (both studies)** — r01-r06 appended to every existing design's `queries.sql`,
  zero change to any existing statement, schema, index or write path. Study 02: shape
  splits on `d.Precreated` (unsold ticket rows exist and need a `status='sold'` filter,
  or not). Study 03: r03-r05 are the same SQL everywhere (the `ticket` table is
  byte-identical across all 14 designs); r01/r02 split on `d.Layout`
  (SeatRows/ClaimRows/Documents) plus one design (E2) whose hold expiry lives on a
  separate `hold` table rather than the seat row. Both append-only diffs verified clean.
- **AM-02.2 (both studies)** — `ReportCoverage(d Design)`, computed from existing design
  flags rather than a literal per-design map, so it cannot drift from the schema without
  the schema itself changing. `TestReportCoverageMatchesSQL` binds it to the actual
  catalogue in both studies; study 03 also gets `TestR06UnanswerableEverywhere`, pinning
  AM-02.0's verified finding so a future hold-history design would fail the test rather
  than pass silently.
- **AM-02.3** — `x1_cas_ledger`: P3 plus `sale_event`, written in the same statement as
  the sale/cancellation via a data-modifying CTE (no change to `Booker`'s control flow;
  `w_sell_seat`'s row count still means "0 = lost the race" because the driver reports
  the final statement's count). Ledger reconciliation audit added
  (`RunLedgerAudit`/`LedgerAudit`), composing with the existing harness-observation audit
  rather than duplicating it. Study 03 gets no new design, as specified.
- **AM-02.4** — Study 02: a `BenchmarkReports` sibling to `BenchmarkReads` with a proper
  multi-parameter bind (`catalog.Bind` against a `map[string]any`), never touching the
  five-question read path. Study 03: its existing `run` helper already bound an
  arbitrary `[]any`, so r01-r05 were appended directly into `BenchmarkReads`' own
  results; a small `runRegime` sibling carries the `@regime` suffix in the reported name
  without changing what `run()`'s existing tier suffix means. Both studies derive
  `since`/`until` from the dataset's load epoch (RECENCY.md section 2); study 03's r01
  additionally needs a live wall-clock cutoff, since "holds expiring soon" is a real-time
  operational question, not a historical one — bound as an absolute instant rather than
  an engine-specific INTERVAL.
- **Two bugs found and fixed before the dev checks would have caught them structurally**
  (both self-corrected while writing the harness, not silent):
  - Study 02: `runReport`'s row-count-only check for X1's r02-style lookups would have
    passed a design that always answered the same status; rewrote to a status-prefix
    check instead.
  - Collect()'s (study 03) fixed column budget of 4 forced r02's original 5-column design
    (status, hold_id, customer_id, hold_expires_at, sold_at) down to 4 by merging the two
    mutually-exclusive timestamp columns into one `COALESCE(hold_expires_at, sold_at) AS
    at` — applied across all four r01/r02 formulations (event_seat, L1, L2, E2) before
    any SQL was applied to a design directory.
- **AM-02.5 — dev checks, `tiny`, `pg-single` then `yb-single`, under the lock.** Two
  real bugs caught and fixed (recorded in a separate commit, `a74b3f8`, since they were
  found by *running* the gate, not by writing it):
  - An earlier edit had duplicated `CREATE TABLE band`/`CREATE TABLE event` in
    `x1_cas_ledger/schema.sql`. Worked in isolation on a fresh database, failed
    ("relation band already exists") the moment it ran after any other design in the
    same container — the correctness gate's whole reason for existing.
  - Study 03's r01 truth counted only `stLive` seats, but expiry is lazy in every design
    except E1 (whose sweeper the harness already runs to completion before `Verify`): a
    hold whose true expiry has passed still reads `status='held'` until swept, which
    r01's own SQL correctly returns. Fixed to count `stExpired` too, except when
    `d.ExpiryOnSweeper`.
  - After both fixes: **study 02, all 15 designs (14 + X1) pass on both engines**
    (44-46 checks each); **study 03, all 13 PostgreSQL-capable designs and all 14
    YugabyteDB-capable designs pass** (58/58 checks each), including L3.
  - C1's negative control fired in a race (every event overbooked); H0's fired in its
    own dedicated holds experiment (overbooking at every tier) — both confirmed, not
    assumed, and neither touched by the reporting changes.
  - One calibration cell per study on `yb-single` (`-phases verify,read`, one design
    each): every answerable/partial report produced real throughput (hundreds to low
    thousands of ops/s at `tiny`); unanswerable reports were correctly absent, never a
    zero. Reports regenerated from both dev-check directories; the answerability tables
    and throughput render exactly as AM-02.1/.4 specify.
  - Every database torn down, lock released, verified before and after.
- No escalations raised in this iteration.

**Next: HIGH — review the answerability tables against the SQL, confirm X1's ledger
audit is real, and size phase 3b's measured runs from the dev-check throughput above.**
