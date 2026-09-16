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
