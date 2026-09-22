# Result — history import and status reconciliation

**Task:** `task-20260922T025912Z-history-status-reconcile`
**Branch:** `repo/ai-work-history-import`, base `795c335`
**Required tag:** `repo/ai-work-history-v1`

## Decision Log

1. **Imports are reference records, not re-execution.** Each historical effort (EH-01, EH-02,
   Studies 03–05) gets a queue task record whose two-event chain (`in_progress → completed`) uses the
   same bootstrap transition the planning checkpoint uses, cites the original handoffs/reports and
   the primary tag, and states that no new work was done. No source artifact is edited.
2. **The pre-rewrite `CONTEXT.md` is archived, never merged away.** It is copied whole to
   `docs/history/context/20260922T132054Z-CONTEXT.md`; the live file is rewritten to current state
   only and links to the archive so historical facts survive without duplicate live receipts.
3. **Public statuses are reconciled to evidence, not to optimism.** Study 02/03 v2 runs are stated as
   measured and analysed, with the phase-3b HIGH validation named as the only open item; Study 05 is
   stated as complete for its single-small evidence set with its limits explicit.
4. **The methodology change is an order fix only.** The 11a block moves between 11 and 11b; every word
   is unchanged.
5. **The separate validation task is released, not created.** `study02-03-v2-validation` already
   existed as a proposed record, so it was released to `ready` rather than duplicated; its dependency
   on `evidence-registry` keeps it unclaimable until that registry lands.

## Evidence

- Archived status narrative: `docs/history/context/20260922T132054Z-CONTEXT.md` (1 183 lines).
- Concise live status: `CONTEXT.md`; archive index: `docs/history/README.md`.
- Imported records: `docs/ai-work/tasks/2026/09/task-20260914T180000Z-study01-integration-eh01/`,
  `…-20260915T210000Z-recency-and-reports-eh02/`, `…-20260915T120000Z-reserved-seating-study03/`,
  `…-20260920T120000Z-configuration-portal-study04/`, `…-20260921T090000Z-cache-consistency-study05/`.
- Navigation and status: `README.md`, `studies/02-ticket-booking/README.md`,
  `studies/03-reserved-seating/README.md`, `studies/05-cache-consistency/README.md`.
- Methodology heading order: `docs/methodology.md` sections 11, 11a, 11b, 12, 12a, 12b.
- Released task: `task-20260922T025912Z-study02-03-v2-validation` → `ready`
  (`event-20260922T132432Z-released`).

## Validation performed

- `queue audit` on the branch: **23 tasks, 0 errors**; the single `live_ref_drift` warning is the
  in-flight claim's live blob and clears on the first checkpoint (see residual risks).
- Relative-link check across every edited Markdown file: **0 broken links**.
- Methodology heading order confirmed 11 → 11a → 11b → 12 → 12a → 12b with unchanged body text.
- Protected artifacts (`studies/*/results`, `studies/*/reports/analyses`, `…/discussions`,
  `docs/handoffs/`) were not modified: this branch's diff touches only `README.md`, `CONTEXT.md`,
  `docs/methodology.md`, `docs/history/`, `docs/ai-work/tasks/` and the three study READMEs.

## Residual risks / limitations

- The queue archive's state for the imported work is a **reference summary**: the authoritative
  record remains the original handoffs, reports, runs and tags under `docs/handoffs/` and `studies/`.
- The `live_ref_drift` audit warning during an active claim is a queue-v1 false positive. Fixing it
  would touch `tools/queue`, outside this task's `owned_paths`; it is recorded here and in
  `PROGRESS.md` for a future queue task.
- The six scientific protocol tasks remain `proposed`; releasing them is a later HIGH gate.
