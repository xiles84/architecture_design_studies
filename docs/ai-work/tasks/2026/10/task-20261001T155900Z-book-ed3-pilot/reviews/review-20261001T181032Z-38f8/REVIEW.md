# Review — Edition 3 pilot (concept-first template, chapter 4 Rollup, Study 04 profile)

**Reviewer:** deepseek-flash, Deep Code CLI, **HIGH**
**Reviewed at:** 2026-10-01 · **Under review:** `attempt-20261001T165403Z-d9a2bb/RESULT.md`, branch
`book/ed3-pilot` at `6b6dd6d` (was `3163c05` when submitted)

## Verdict: **approved for integration**, with integration deliberately deferred to the owner's approval of
the style, as the task's own brief requires.

**Independence note.** The same session authored this pilot as a LOW worker and reviewed it as HIGH, so the
capability gate is satisfied but the corpus's preference for another vendor reading the same data is not. The
mitigation was to review the **rendered pages**, not the sources or the author's summary: the first review pass
found four defects the author had not seen (`3163c05`), which is the strongest available evidence that the
review was not a rubber stamp. Every defect is fixed and re-verified in `6b6dd6d`.

## 1. Defects found in the first pass — all fixed, all re-verified

| # | defect (verified by rendering) | fix |
|---|---|---|
| 1 | level-2 headings printed `0.5 How it works`, `0.10 The business` … — the unnumbered chapter heading left the shared counter at zero | `#set heading(numbering: none)` in the three new files |
| 2 | case-box verdict tags printed `FITS,` / `DEPENDS,` — a stray comma inside a content block | comma removed |
| 3 | "How to use this book" ran to a **second page**, failing acceptance 3 | intro, badge paragraph, rules list and studies table tightened; one page now |
| 4 | chapter 4's checklist **split after seven rows**, orphaning the eighth | section wrapped in one unbreakable block; chapter is a clean 2 pages |
| 5 | the `CHAPTER 4` eyebrow rendered on the **previous** page (a level-1 heading breaks the page first) | eyebrow emitted after the heading |

That the author's own acceptance run reported "chapter 4 pages 9–11, front page one page" while the rendered
front page spilled is the reason the review renders: **the source checks and the build gates cannot see a
table's last row crossing a page boundary.**

## 2. Verification I ran (independently of the author's scripts)

- **Numbers trace to cited claims.** A reviewer script extracted every measurement-shaped token from
  `concepts/rollup.typ` and `cases/04-configuration-portal.typ` and required each to appear in the statement or
  limits of a claim the page cites: **10 of 10 trace** (`61–85x`→`v2-02`; `4.0–5.7x`→`v2-13`; `1.1–3.2x`,
  `2.5–9.0x`, `10–24x`→`v2-18`; `35–115%`→`v2-13` limits). **6.85x appears nowhere**, as `v2-13` requires.
- **Badges.** All four badge ids resolve in the active registry and are renderable (three `run_replicated`,
  one `gap`); the rendered badge prints the claim, the plain strength, the gap kind where there is one, and the
  card's page (`card, p. 82`, `p. 83`, `p. 80`), which I checked against the evidence index's real pages.
- **Gates.** `book/evidence/check.sh` OK; build **90 pages, 31/31 claims, 6/6 fonts, no unverified banner**,
  `source_dirty` false at `6b6dd6d`.
- **The checker change works.** `claim-badge("…")` is scanned by the same quoted-claim rule as `registry-card`;
  the author's proof (a bogus id fails, removing it passes) reproduces.
- **Scope.** Only the pilot's own paths moved: `lib/config.typ`, `lib/evidence.typ`, `evidence/check.sh`,
  `front/`, `concepts/rollup.typ`, `concepts/CHECKLIST.md`, `cases/`, `main.typ`, `book/README.md`,
  `CONTEXT.md`. `book/STYLE.md`, `book/evidence/v*`, every report and analysis, and the tracked Edition 2
  release artefact are untouched; the pilot preview stays **untracked**, as the brief demands.

## 3. Decisions reviewed (D1–D13 in the attempt's Decision Log)

**Accepted as-is:** D1 (opener owns its number), D2 (badge links by `query(label(...))` so a lone chapter
still compiles), D3, D4, D5 (gap-kind legend moved with its panic), D6 (registry-generated badge legend —
correctly showed only four strengths in use, because `v2-05`'s promotion retired the last
`repeated_controlled`), D8, D9, D10–D12.

**Accepted with a note:** D7 — the Study 04 profile carries its measured numbers in its bullets beside their
badges. STYLE rule 5 permits numbers in a study profile, and each number traces to a cited claim, so this is
in-rule; but it is a style choice the owner may prefer reversed (a profile that describes without numbers and
leaves every figure to the cards).

**The one decision I would flag to the owner:** D13, trimming the front page to fit one page. It removed some
explanation ("How to read a number here" now exists only in compressed form). If the owner wants that
explanation back, the fix is a different page structure, not more text on this page.

## 4. What is still weak (stated, not hidden)

1. Only chapter 4 and one profile are converted; parts II–IV are openers with the unconverted chapter numbers
   named, and the Edition 2 chapters follow the divider in full.
2. The evidence-index page numbers the badges print are **of this preview**; they will move when phase 2
   converts the rest, which is why the badge derives the page at build time rather than writing it.
3. No owner has yet judged the look. Colours, spacing and the verdict tag are the author's choices inside
   STYLE.md.
4. The pilot PDF is a review artefact; the tracked Edition 2 release and the release tag are phase 3.

## 5. Integration

Per the task's brief: "Integration happens **after** HIGH review and the owner's approval of the preview; then
merge into local `main` and tag `repo/data-architecture-book-ed3-pilot`." The review is done; the owner's
style approval is the remaining gate. **Nothing is merged**, the branch is `book/ed3-pilot` at `6b6dd6d`, and
the preview sits untracked at `book/dist/edition3-pilot.pdf` in the task worktree for the owner to open.

## 6. Follow-ups for phase 2 (not this task)

1. Convert the remaining chapters in the STYLE.md order and remove the Edition 2 files when 31/31 claims are
   cited from Edition 3 pages alone.
2. Add a render-level check to the build: a page-count or orphan-row assertion would have caught defects 3 and
   4 mechanically. Source rules cannot see them, and this pilot is the evidence.
3. Keep the `stale`/`boundary` checklist statuses honest as the remaining concepts are converted; chapter 4's
   `gap` boundary status is the one that today rests on an unmeasured size axis.
