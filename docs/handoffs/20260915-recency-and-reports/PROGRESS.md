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
