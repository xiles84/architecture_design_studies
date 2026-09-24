# Review — book figure proof set

**Task:** `task-20260923T235922Z-book-figure-proofset`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Verdict:** **approved for integration** — no rework required
**Required tag:** `repo/book-figure-proofset-v1`

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | book-owned sources under `book/assets/sources/` | accepted — a book figure needs a canonical source inside owned paths; study sources stay canonical |
| 2 | study-03 reuse = selection + caption, not redraw | accepted — the book must not hold an editable copy, and the study sources were not touched |
| 3 | two-label mechanism + text equivalent | accepted — matches the acceptance exactly |
| 4 | figure 3 redesigned after reader 2 | accepted — the redesign makes the shared-instance claim visible |
| 5 | figure 4 relabelled, not redesigned | accepted — honest; the source is a sequence and the book cannot edit it |
| 6 | both reader passes one model, recorded | accepted — the limitation is stated, not hidden |
| 7 | PDF rebuilt, labelled draft | accepted — this task owns `main.typ`/`dist`; released v1 untouched |

## Independent checks re-run by the reviewer

```text
book/assets/export.sh --check                    5 figures verified
observed-result figure embeds a rate?             no (the only x/% strings are SVG element ids)
figure-evidence() calls in prose                  4 (one per figure)
dist/build-manifest.json figure_registry + total  present, 5
dist/build-manifest.json evidence digest          814f40949d162f56… (v3)
build result                                      49 pages, 5/5 fonts, 29/29 claims, exit 0
methodology rule 15 present                       yes
```

Acceptance criteria met: three-to-four labelled figures covering the structure, a Study 02/03 sequence
reuse, a new Study 05 stale-fill sequence and an expiry figure; two independent labels plus a text
equivalent; the observed-result figure cites `v2-15` and embeds no number; the two-reader check is
recorded; methodology 15 exists; the PDF rebuild records every figure source and embedded asset.

## Findings (non-blocking)

1. The two-reader check is not independent (one model). Recorded by the executor.
2. Figures 2 and 4 depend on study sources; a later change requires a re-export (the build fails until
   it happens, by design).
3. The rebuilt PDF is a draft and has not been release-reviewed; the released v1 remains at its tag.

## Verdict

Approved for integration. The figures are labelled, text-equivalent, provenance-tracked and gated, and
none of them smuggles a performance number into the prose.
