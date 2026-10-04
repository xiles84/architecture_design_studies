# Applied business chapters — source review

Parent task: `task-20261004T132000Z-book-theory-cases`.
Actor: HIGH, GPT-6 through Codex desktop; selected effort unknown.
Owned branch/worktree: `repo/book-applied-cases-20261004`,
`.worktrees/book-applied-cases-20261004`.

## Decisions

- Use one chapter per business and variation sections per decision. Each variation
  links to its theory strategy. Reason: follow the owner's theory/business hierarchy.
  Confidence: High. Affects all five case chapters.
- Mark observations, inferred fit boundaries and unmeasured regimes explicitly.
  Reason: a business selection rule must not sound like a measured crossover.
  Confidence: High. Alternative: an unqualified thrive/die column.
- Use compact summaries and badges; full cards stay in the parent's evidence appendix.
  Reason: avoid repeated registry prose in teaching chapters. Confidence: High.
- Correct the old configuration profile's document identifier to `d1_doc_row`, as
  registered in the study README. Do not claim a numeric document comparison there.
  Confidence: High. Affects only the configuration case.
- Retain the active registry's older Study 04 gap wording transparently alongside
  later targeted replications. Do not silently rewrite immutable registry content.
  Confidence: High. Affects interpretation of the configuration gap.

## Evidence check

Read each active v6 claim statement and relevant limits, plus the source study
READMEs. Summary bands copy claims; no derived ratios or absolute levels whose
replication failed are promoted. The SELECT-family numerical programme is not
added to the book's claim set.

The cache case exposes the impossible-value churn gate failure, missing physical
placement verification, partial endpoint balance and resource-framing boundary.
The seating map gain is an engine-specific package comparison; guarded-confirm
refusals are a correctness result with no throughput conclusion.

## Checks

Run `check-cases.go` in a Podman Go container: five labelled business chapters,
all 31 active claims cited, no unknown claims or full cards, audited prose lines
at most 60 words. Exact stdout is in `SOURCE-CHECK.txt`.

Integrated compilation, pagination and rendered proof belong to the parent task.
No independent PDF build, benchmark, database start or evidence mutation was done.
