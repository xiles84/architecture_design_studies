# Result — Edition 3 pilot: concept-first template, chapter 4 Rollup, the Study 04 profile

**Task:** `task-20261001T155900Z-book-ed3-pilot` · **claim:** `claim-765cf86eb1b252c2` (epoch 1) ·
**attempt:** `attempt-20261001T165403Z-d9a2bb` · **worker:** deepseek-flash, Deep Code CLI, LOW executor
**Outcome: delivered, awaiting HIGH review and the owner's approval of the preview.**

## 1. The preview

| | |
|---|---|
| Preview PDF | `book/dist/edition3-pilot.pdf` (**untracked**, in this worktree) |
| Preview manifest | `book/dist/edition3-pilot-manifest.json` (**untracked**) |
| Total pages | **90** (after the review fixes in section 6) |
| New pages | **p. 3** How to use this book (1 page) · **pp. 8–9** chapter 4 Rollup (2 pages, part I opener on p. 7) · **p. 13** Study 04 profile (1 page) · Edition 2 divider starts **p. 14** |
| Build | `book/build.sh --out dist/edition3-pilot.pdf --manifest dist/edition3-pilot-manifest.json` |
| Open it | `xdg-open book/dist/edition3-pilot.pdf`, or copy the file out of the worktree; the text layer can be read with `pdftotext -layout` inside the pinned image |

Build gates: **claims indexed 31/31**, fonts embedded 6/6, no unverified banner, `book/evidence/check.sh`
OK, tree reported dirty (pilot sources uncommitted, as expected for a review build).

## 2. Acceptance evidence

**2.1 Macros exist and fail loudly (steps 2–3, acceptance 1).** All seven components are in
`lib/config.typ`; `claim-badge` and `strength-label` are in `lib/evidence.typ`. Each guard was exercised with
bad input and the build stopped with a message naming the rule, the value and the file to read:

| Input | Message it produced |
|---|---|
| a 26-word key idea | `concept-opener Rollup: the key idea has 26 words; STYLE.md allows 25. Shorten it…` |
| `verdict: "maybe"` | `case-box S: verdict "maybe" is not one of fits, does not fit, depends (STYLE.md, case box)` |
| 7 checklist rows | `checklist-table: needs exactly 8 rows, got 7 (methodology 21; the eight questions are fixed)` |
| status `"sorted"` | `checklist-table: unknown status "sorted"; use one of direct, mechanism, analogy, gap` |
| 4 remember bullets | `remember: needs 1-3 bullets, got 4 (STYLE.md, template 9)` |
| a bogus badge id | `no claim with id v2-99-nope in book/evidence/v6/claims.json; do not write a number without a claim` |

**2.2 The checker catches a bad badge (step 4, acceptance 1).** `book/evidence/check.sh`'s quoted-claim rule
now scans `claim-badge("…")` alongside `registry-card`/`claim-by-id`. Proof, run in this worktree:

```
$ printf '#import "../lib/evidence.typ": claim-badge\n#claim-badge("v2-99-bogus-id")\n' > scratch/bogus.typ
$ book/evidence/check.sh
  [FAIL] unknown-claim-id: v2-99-bogus-id is not a claim in v6
check: 1 failure(s)
$ rm scratch/bogus.typ && book/evidence/check.sh
check: OK
```

**2.3 Badges link to the card and print its page.** Rendered on chapter 4's case-box page: 
`v2-02-rollup-family · replicated · card, p. 82`, `v2-13-trigger-rollup-cost · replicated · card, p. 83`,
`v2-18-maintenance-size-interaction · replicated · card, …`, `v2-gap-03-study04-missing-designs · gap ·
coverage gap · card, p. 80`. The link targets are the labelled cards in `evidence-index()`.

**2.4 Chapter 4 follows the template and cites its four claims (acceptance 2).** Ten sections in order:
opener (25-word key idea), the problem, the figure (`assets/fig-rollup-vs-rolldown.svg`, reused, with a text
equivalent), five "how it works" bullets, a 3×4 fit table, the five-row cost table, two case boxes, a gap
callout, three remember bullets, the checklist with its eight rows. Pages **8–9**.

**2.5 Numbers in the new pages, and the claim each appears in (acceptance 4).** Copied as written; no value
a claim says did not reproduce is quoted (in particular **6.85x is absent**).

| Page | Box / field | Text as printed | Source |
|---|---|---|---|
| 10 | case box, Study 01, `result` | `about 61–85x the minimal normalized design on aggregate reads` | `v2-02` statement ("D4 rollup/trg about 61-85x … D1 minimal in two whole runs") |
| 10 | case box, Study 01, `result` | `about a third to two-fifths of its rate` | `v2-02` statement ("about a third to two-fifths of D1's insert rate in the same two runs") |
| 10 | case box, Study 04, `result` | `4.0–5.7x slower with the trigger` | `v2-13` statement ("costs 4.0-5.7x the application-maintained rollup") |
| 10 | case box, Study 04, `result` | `about 1.1–3.2x at 10 entries, 2.5–9.0x at 60 and 10–24x at 240` | `v2-18` statement ("roughly 1.1-3.2x … at 10 entries, 2.5-9.0x at 60 and 10-24x at 240") |
| 10 | gap callout | `Eight of Study 04's eighteen planned designs are not built` | `v2-gap-03` statement ("Eight of the eighteen planned Study 04 configuration-portal designs are not implemented") |
| 15 | profile, "What it showed" 1 | `4.0–5.7x the application-maintained rollup` | `v2-13` statement (as above) |
| 15 | profile, "What it showed" 2 | `roughly 1.1–3.2x at 10 entries, 2.5–9.0x at 60 and 10–24x at 240` | `v2-18` statement (as above) |
| 15 | profile, "What it showed" 3 | `disagreed by 35–115%` | `v2-13` limit ("disagreed by 35-115% across the three runs") |

The front page and every paragraph outside a box carry no measured number: a scan of
`front/how-to-use.typ` and `cases/04-configuration-portal.typ` for `<number>x|%|ms|k/s|ops/s` finds hits only
in the profile bullets above. Paragraph limits were checked mechanically over the three new files: **no
paragraph over 60 words or 3 sentences** (the one reported hit is a Typst code block in the front page that
the checker mistook for prose).

**2.6 The front page is one page, the profile is one page (acceptance 3).** p. 3 and p. 13, confirmed by the
page map above. The first build failed this: the front page ran to a second page, found in review and fixed
(section 6).

**2.7 The build verifies (acceptance 5).** 31/31 claims indexed, 6/6 fonts, no banner, checker OK.

## 3. What was changed

| Path | Change |
|---|---|
| `book/lib/config.typ` | `book-edition` → *Edition 3 (draft)*; added `content-text`, `word-count`, `enforce`, `concept-opener`, `fit-table`, `cost-table`, `case-box`, `remember`, `checklist-table`, `part-opener` |
| `book/lib/evidence.typ` | added `strength-label`, `card-label`, `claim-badge`; labelled every card in `evidence-index()` |
| `book/evidence/check.sh` | quoted-claim rule also scans `claim-badge("…")` |
| `book/front/how-to-use.typ` | new: chapter anatomy, badge legend generated from the registry, the reading rules moved off the provenance page, the five studies at a glance; the gap-kind legend and its panic moved here |
| `book/concepts/rollup.typ` | new: chapter 4 to the template |
| `book/concepts/CHECKLIST.md` | added the `rollup` row |
| `book/cases/04-configuration-portal.typ` | new: Study 04 profile |
| `book/main.typ` | Edition 3 front page + parts I–IV with the converted files placed, the Edition 2 divider and the old files, then glossary and evidence index; the provenance page keeps provenance and points to the front page |
| `book/README.md` | layout and edition section |
| `CONTEXT.md` | records the pilot |

Untouched, as the brief requires: `book/STYLE.md`, `book/evidence/v*`, every report, analysis and registry
package, and the tracked Edition 2 release artefact in `book/dist/`.

## 4. Decisions

The Decision Log is in `PROGRESS.md` beside this file (D1–D9). The two a reviewer should look at first:
**D2** (the badge links by `query(label(...))` so a lone chapter still compiles) and **D7** (the profile
carries its numbers beside its badges — a style judgement the owner may reverse).

## 5. Known limits of this pilot

- Only chapter 4 and the Study 04 profile are converted; parts II–IV hold openers and names, and the
  Edition 2 chapters after the divider still carry their own older prose.
- The pilot PDF is a review artefact, not the release: `book/dist/data-architecture-reference.pdf` was not
  rebuilt and no tag was created (the brief's phase 3 does that).
- The figures are the existing assets; chapter 4 reuses `fig-rollup-vs-rolldown.svg` rather than adding a
  chapter-specific one.
- No owner review of the *look* has happened yet: colours, spacing and the verdict tag are my choices within
  STYLE.md, and the brief says those are decide-and-continue rather than escalate.

## 6. What the HIGH review found, and the fixes

The review was not a source read: the pages were rendered and looked at. Four defects came out of that, all
fixed in the same branch (commit `6b6dd6d`), and all verified by re-rendering.

| # | defect | why it mattered | fix |
|---|---|---|---|
| 1 | level-2 headings printed **0.5 How it works**, **0.10 The business** … in all three new files | the chapters' level-1 heading is unnumbered, so the shared counter stayed at zero and every section was numbered under a chapter that has no number | `#set heading(numbering: none)` in the three new files; the chapter number lives in its opener |
| 2 | the verdict tag rendered **`FITS,`** / **`DEPENDS,`** | a stray comma inside a content block, visible on every case box | removed; the tags now read `FITS`, `DEPENDS` |
| 3 | "How to use this book" **overflowed onto a second page** | acceptance 3 requires one page; the spill was only the studies table's last row and the closing line | intro, badge paragraph, rules list and studies table tightened; the table's inset reduced; now one page |
| 4 | chapter 4's checklist **split after seven rows**, leaving the eighth alone on a page | the checklist is the chapter's frame, and an orphaned row reads as an afterthought | the heading, intro and table are one unbreakable group; the chapter is now a clean 2 pages |

One more, found while fixing 1: the **`CHAPTER 4` eyebrow was on the wrong page**, because a level-1 heading
breaks the page first and the eyebrow was emitted before it. It is now emitted after the heading, under the
rule, where it reads as the chapter's label without touching the contents entry.

**Effect:** 90 pages instead of 92; front page p. 3, chapter 4 pp. 8–9, profile p. 13. The badge page numbers
in section 2.3 were updated to the corrected build (p. 82, p. 83, p. 80).
