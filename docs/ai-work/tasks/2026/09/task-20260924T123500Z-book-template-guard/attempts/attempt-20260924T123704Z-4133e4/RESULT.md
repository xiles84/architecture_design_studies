# Result — template-token guard

**Task:** `task-20260924T123500Z-book-template-guard`
**Branch:** `book/template-guard`
**Required tag:** `repo/data-architecture-book-template-guard`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only; no database, no measurement

## What was delivered

1. **All six leak sites fixed.** `book/main.typ` (title-page provenance line, provenance table, Chapter 18
   heading, index prose), `book/chapters/evidence-and-reproduction.typ`, and
   `book/concepts/cache-consistency.typ` §15.4. The heading needed string concatenation
   (`"Evidence registry (active " + registry-version + " claims)"`), because a `#name` inside a plain
   string argument is never evaluated — the second of the two causes. The rest use `#raw(registry-label)`,
   which evaluates and keeps the monospace look.
2. **A source rule.** `book/evidence/check.sh` fails on a backticked token starting with `#` (raw text is
   not evaluated) or on a build-input token inside a string argument.
3. **A rendered-page gate**, which is what actually catches this class: `book/build.sh` greps the extracted
   PDF text for `#registry-`, `#build-`, `#source-`, `${`, `{{` and `UNKNOWN_PLACEHOLDER`, records
   `unresolved_tokens_in_pdf` in the manifest, and aborts the build.
4. **The review exchange recorded**, attributed, at
   `book/evidence/reviews/20260924-edition-2-review-response--deepseek-flash.md`: the reviewer's
   confirmation that they examined a hand-built preview, their acceptance of our two corrections, the
   thirteen new findings, and the `v3-06` disposition that differs from their suggestion.
5. **`CONTEXT.md` and `LESSONS_LEARNED.md`** updated with the defect, the gates and the general lesson.

## Verified from committed files

```text
check.sh (source rule)                OK
build.sh --out dist/preview.pdf       50 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                      evidence digest present, unverified banner absent,
                                      unresolved tokens in pdf: false  — verified
rendered registry references          book/evidence/v4/claims.json appears 4x in the page text
                                      heading reads "Evidence registry (active v4 claims)"
leftover tokens in the page text      0
owner-edited files in book-evidence-v4  both clean; the scratch control file is absent
```

### Negative controls

| Control | Expected | Observed |
|---|---|---|
| A backticked `` `#registry-label` `` added to a source | source rule fails | `[FAIL] template-interpolation: book/main.typ:157 … backticks a # token, which renders literally instead of evaluating` |
| A literal `{{` added to a source — the source rule deliberately does not look for it | build fails on the page text | `verification FAILED: unresolved template token(s) rendered in the PDF: {{` |

Both were introduced, observed and reverted; `book/main.typ` was restored from a pre-injection copy and
its diff shows only the four intended patches.

## Housekeeping

The two files the owner edited by hand in the completed `book-evidence-v4` worktree were inspected before
any cleanup: `book/_negative-control.typ` is absent and
`book/concepts/cache-consistency.typ` has no uncommitted change there. Nothing was deleted or overwritten.

## Conclusion

The release blocker is closed at the source and, more importantly, would now fail the build rather than
reach a reader. No claim, number, report or analysis was touched; the registry stays v4, and the version
bump is the next task's business.
