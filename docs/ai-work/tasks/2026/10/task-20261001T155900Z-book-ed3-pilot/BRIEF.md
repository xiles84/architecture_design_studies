# Edition 3 pilot: concept-first template, chapter 4 Rollup and the Study 04 profile

**Task:** `task-20261001T155900Z-book-ed3-pilot` · **capability:** LOW may execute · **role:** executor · **benchmark:** none
**Branch:** `book/ed3-pilot` · **Worktree:** `.worktrees/book-ed3-pilot`
**Planner:** claude-opus-5-5, Claude Code desktop, HIGH, 2026-10-01

## Goal

The owner wants the book to teach **concepts first** and use the **case studies to show when each
concept fits**, with **short text**. This pilot builds the new template and one complete example, so
the owner can approve the style before the whole book is converted.

Read [`book/STYLE.md`](../../../../../../book/STYLE.md) first. It is the contract for every page.

## Decisions already made (do not re-litigate)

- Structure, chapter template, case box, study profile, writing rules and evidence badges: as in
  `book/STYLE.md`. Change requests to STYLE.md go in the Decision Log, not into the file.
- Pilot scope: the front page "How to use this book", chapter **4 Rollup**, the **Study 04** profile.
  Nothing else is converted.
- Evidence: `claim-badge(id)` inline; full cards only in the evidence index.
- No new numbers, no measurement, no registry edits. `book/evidence/v*` is read-only.
- The Edition 2 chapters stay in the build, after the new parts, under a divider that says they await
  conversion. Every renderable claim must still be cited (31 today).
- `book-edition` becomes `"Edition 3 (draft)"`. Tracked `book/dist/` files are not touched.

## Steps

1. **Claim** with the queue; check `git worktree list`, `CONTEXT.md` and `podman volume inspect
   ads-book-build-lock` before building.
2. **Macros** in `book/lib/config.typ`: `concept-opener`, `fit-table`, `cost-table`, `case-box`,
   `remember`, `checklist-table`, and a part opener. Each validates its input and panics with a clear
   message: `verdict` ∈ {fits, does not fit, depends}; checklist has exactly 8 rows; status ∈ {direct,
   mechanism, analogy, gap}.
3. **Badges** in `book/lib/evidence.typ`: `claim-badge(id)` as STYLE.md describes. Resolve the claim
   through `claim-by-id` (so unknown or retired ids still panic). Give each full card in
   `evidence-index()` a label so the badge links to it and prints its page.
4. **Checker**: extend the quoted-claim rule in `book/evidence/check.sh` to `claim-badge("…")`.
   Prove it: a badge with a bogus id must fail the check, then remove the bogus id.
5. **Front page** `book/front/how-to-use.typ`, one page: the anatomy of a concept chapter, the badge
   legend, and the five studies at a glance (study · business · main concepts). Move "How to read a
   number here" off the provenance page into it, shortened. Keep the gap-kind legend generated from the
   registry, and keep its panic for an undefined kind.
6. **Chapter 4** `book/concepts/rollup.typ`, to the template. Cite `v2-02`, `v2-13`, `v2-18` and
   `v2-gap-03`. Reuse `assets/fig-rollup-vs-rolldown.svg`. Add its row to `book/concepts/CHECKLIST.md`.
7. **Profile** `book/cases/04-configuration-portal.typ`, to the profile template, one page. Check the
   domain words against `studies/04-configuration-portal/README.md`.
8. **Skeleton** in `book/main.typ`: the Edition 3 parts in STYLE.md order, with only the converted
   files placed; then the Edition 2 divider and the old files; then glossary and evidence index.
   Update `book/README.md` (layout).
9. **Build** `book/build.sh --out dist/edition3-pilot.pdf --manifest dist/edition3-pilot-manifest.json`.
   Keep both files **untracked** in the worktree for the owner's review; never commit them.
10. **Self-check** and write `RESULT.md` (see acceptance). Update `CONTEXT.md`; add to
    `LESSONS_LEARNED.md` only if something was learned. Commit explicit paths, submit for review.

## Style sample (pins the voice; wording may change, the numbers rule may not)

```typst
#concept-opener(4, "Rollup",
  [Store a total on the parent so reads skip the children — and pay for it on every child write.])

A page shows each charity's total. Summing every donation on each view repeats the same work.
Storing the total once makes that read a single row.

#case-box(
  study: "Study 01 · charity donations",
  verdict: "fits",
  situation: [A charity site shows totals per charity and per donor.],
  compared: [Totals stored on the parent (by trigger or by the application) against totals
    computed from the donations on every read.],
  result: [Stored totals scored about 61–85x the minimal normalized design on aggregate reads.
    Inserts ran at about a third to two-fifths of its rate.],
  lesson: [A rollup fits when totals are read much more often than donations are written.],
  watch-out: [The read gain and the insert penalty reproduced; the absolute insert rates did not.],
  evidence: [#claim-badge("v2-02-rollup-family")],
)

#case-box(
  study: "Study 04 · configuration portal",
  verdict: "depends",
  situation: [Products read their configuration from a portal; an overview rolls each one up.],
  compared: [The overview kept by a database trigger against kept by the application.],
  result: [Replacing a whole configuration was 4.0–5.7x slower with the trigger. The gap grew with
    size: about 1.1–3.2x at 10 entries, 10–24x at 240.],
  lesson: [A trigger buys a guarantee no code path can forget, paid on every row written.
    Check how large your writes are.],
  watch-out: [Small scale, one machine, and repeats disagreed widely: trust the ordering, not the ratio.],
  evidence: [#claim-badge("v2-13-trigger-rollup-cost") #claim-badge("v2-18-maintenance-size-interaction")],
)
```

Never quote 6.85x: `v2-13` says it did not reproduce.

## Acceptance

1. The macros, `claim-badge` and the checker change exist and fail loudly on bad input (step 4 proof
   recorded in `RESULT.md`).
2. Chapter 4 follows the template, is 2–4 pages, and cites its four claims by badge.
3. The Study 04 profile is about one page; "How to use this book" is one page.
4. In the new files: no paragraph over 60 words or 3 sentences; numbers only in case boxes and the
   profile; `RESULT.md` lists each box → its numbers → the claim each number appears in.
5. `book/evidence/check.sh` passes; the build verifies with every renderable claim indexed and cited.
6. `RESULT.md` gives the page numbers of the new pages in the preview, the total page count, and how
   to open the preview.
7. `CONTEXT.md` updated; Decision Log in `PROGRESS.md` of the attempt; submitted for review.

Integration happens **after** HIGH review and the owner's approval of the preview; then merge into
local `main` and tag `repo/data-architecture-book-ed3-pilot`.

## Constraints

- Stage explicit paths; never `git add -A`. Never push or pull.
- Do not edit another analyst's file, any report, analysis, registry package or `book/STYLE.md`.
- `book/build.sh` takes `ads-book-build-lock`; it never needs the benchmark lock. T2–T5 replication
  tasks may be measuring at the same time; that does not block this task.
- If `.git/index.lock` is left behind by the queue CLI, check `ps` and `podman ps` first, then remove it.

## Escalate only if

- a STYLE.md rule cannot be met without changing a claim, a gap's visibility or a build gate;
- the badge-to-card link cannot work and the only fallback would hide a gap or a strength;
- a T2–T5 task changes a claim this pilot cites in a way that makes a case box wrong.

Everything else — macro look, colours, heading levels, file names inside `book/front` and
`book/cases` — decide, log, continue.

NEXT MODEL: LOW
