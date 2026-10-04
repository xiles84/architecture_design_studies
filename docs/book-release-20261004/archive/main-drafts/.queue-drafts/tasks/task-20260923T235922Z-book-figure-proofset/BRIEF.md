# Execution brief — Author the book figure proof set with provenance labels and the methodology figure rule

Read the visual-explanations brainstorm conclusion first
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies/CONCLUSION.md`),
then `book/lib/config.typ` (visual vocabulary), `book/FIGURES.md` (chosen import mechanism),
`book/concepts/` and `book/chapters/`, and `docs/methodology.md`.

## Goal

Give the book a small, question-led visual layer that improves explanation of the mechanisms prose handles
worst, without duplicating study schema diagrams and without ever letting an illustration read as a measurement.

## Decisions already made (do not re-litigate)

- **No figure quota.** The proof set is 3-4 figures; expand only on demonstrated reader benefit.
- Reuse and adapt existing Study 02/03 sequence sources rather than redrawing them; the study-owned source stays canonical.
- Form and evidence status are **two independent labels**, never one enum.
- An observed-result figure cites an active registered claim; measured numbers stay in prose and claim cards, never inside SVG text.
- The figure-import mechanism is already chosen by the provenance dependency; do not reopen it.
- Book cache prose has already been reworked by the decision-map dependency; add figures to it, do not rewrite it.

## Acceptance criteria

See the task spec. Minimum: 3-4 figures with the two labels and text equivalents; the reader check recorded
with keep/redesign/drop outcomes; the methodology figure rule added; PDF rebuilt with assets in the manifest.

## Constraints

- No measurement; no benchmark lock.
- Never edit another analyst's analysis, report or digest.
- Commit explicit paths, tag `repo/book-figures-v1`, merge through the queue lifecycle.

## Escalate only if

- the reader check shows a figure actively creates a misconception that cannot be fixed by relabelling (then
  the figure set, not just the figure, is in question).

NEXT MODEL: LOW
