# Progress and Decision Log — Edition 3 pilot

**Task:** `task-20261001T155900Z-book-ed3-pilot` · **claim:** `claim-765cf86eb1b252c2` (epoch 1) · **attempt:**
`attempt-20261001T165403Z-d9a2bb` · **worker:** deepseek-flash, Deep Code CLI, LOW executor

The Decision Log records the decisions that affect correctness, the build or the review. Each entry: the
decision, why, confidence, the alternative considered, and what it affects.

## Decisions

**D1 — `concept-opener` owns its own number instead of relying on Typst's heading counter.**
*Why:* the opener takes the chapter number the owner's structure fixes (4 = Rollup), but the Edition 3 parts
now precede the Edition 2 chapters, so Typst's automatic level-1 counter would number the pilot chapter "1".
*Alternative:* renumber with a counter reset — fragile, and it would have renumbered every Edition 2 chapter.
*Confidence:* high. *Affects:* `lib/config.typ`, main.typ's structure, the contents.

**D2 — the badge resolves through `claim-by-id` and prints the card's page with `query(label(...))`.**
*Why:* the brief requires a badge to link to the full card and print its page, and an unknown or retired id
must still panic. `query` returns an empty array when the index is absent, so a chapter compiled alone in an
editor prints "card in the evidence registry index" instead of failing.
*Alternative:* `locate`, which panics when the label is missing — that would break a single-chapter preview.
*Confidence:* high. *Verification:* the pilot renders `v2-02-rollup-family · replicated · card, p. 84` and
three more badges on page 10. *Affects:* `lib/evidence.typ`, `concepts/rollup.typ`, `cases/*`.

**D3 — the `card-label` definition sits above `evidence-index()`.** *Why:* Typst resolves names in order, and
the first build failed with "unknown variable: card-label" because the helper was appended after its use.
*Confidence:* high. *Affects:* `lib/evidence.typ`.

**D4 — `remember` takes bullets as separate arguments, not an array.** *Why:* the natural call is
`remember([a], [b], [c])`; the first draft took one array argument and rejected that call. *Confidence:*
high. *Affects:* `lib/config.typ`, chapter 4.

**D5 — the gap-kind legend moved from the provenance page to the front page, with its panic.**
*Why:* the brief moves "how to read a number" to the front page; the legend is generated from the registry
and its panic for an undefined kind is the gate that keeps the count honest, so both moved together rather
than the legend being copied. *Confidence:* high. *Affects:* `front/how-to-use.typ`, `main.typ`.

**D6 — the badge legend is generated from the registry, including a label for `repeated_controlled`.**
*Why:* STYLE.md names five strengths but the registry carries six; STYLE says a strength added later must
get a label or the build fails, so `strength-label` panics on an unknown one and the page lists what the
registry actually uses. *Alternative:* hard-code the five and let `repeated_controlled` render raw. That is
what STYLE forbids. *Confidence:* high. *Affects:* `lib/evidence.typ`, `front/how-to-use.typ`.

**D7 — the Study 04 profile carries its measured numbers in the "What it showed" bullets.**
*Why:* the first draft listed the findings with badges but no numbers, which is allowed but weak; STYLE rule
5 permits numbers in a study profile beside their badge, and the numbers are copied from the claims'
statements. *Confidence:* medium (a style judgement the owner may reverse). *Affects:*
`cases/04-configuration-portal.typ`.

**D8 — the unconverted chapters are named by number in each part opener and listed in full after the
divider.** *Why:* the brief keeps them in the build under a divider; naming their chapter numbers in the part
openers means a reader can see what the converted edition will contain without a separate placeholder page.
*Confidence:* high. *Affects:* `main.typ`, the contents.

**D9 — no waiver was added to `book/evidence/check-waivers.txt`.** *Why:* the checker's rules all pass once
the new chapter has its `CHECKLIST.md` row; a waiver would hide a real rule. *Confidence:* high.

## Things worth flagging to the reviewer

- The pilot's own build writes `dist/edition3-pilot.pdf` and `dist/edition3-pilot-manifest.json`, which the
  brief requires to stay **untracked**. `book/dist/pilot.txt` (my page-number extraction) was deleted after
  use; nothing else was left in `dist/`.
- `book/STYLE.md`, `book/evidence/v*` and every report/analysis are untouched, as the brief requires.
- The tracked Edition 2 release artefact (`dist/data-architecture-reference.pdf`) was **not** rebuilt: the
  pilot is a separate output, and rebuilding the release artefact is phase 3 of the rollout.
