# Review — task-20260924T123500Z-book-template-guard

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review. The checks were run against the committed sources and the rebuilt PDF, not
from memory of writing them.

## What was checked

| Check | Method | Result |
|---|---|---|
| All six leak sites fixed | `grep -rn '`#\|\"#registry-'` over the book sources | none remain |
| The page names the registry | extracted text | `book/evidence/v4/claims.json` appears 4×; the heading reads "Evidence registry (active v4 claims)" |
| No token anywhere in the page | count of `#registry-label` / `#registry-version` in the extracted text | 0 |
| The source rule fires | negative control: a backticked token | `[FAIL] template-interpolation` with the line named |
| The page gate fires | negative control: a literal `{{`, which the source rule does not look for | build aborts: `unresolved template token(s) rendered in the PDF: {{` |
| The build is green afterwards | `book/build.sh --out dist/preview.pdf` | 50 pages, 29/29 claims, 6/6 fonts, digest present, banner absent, no unresolved tokens |
| Nothing else moved | `git diff --stat` on the branch | three sources, two scripts, one new review record, context and lessons |
| Owner-edited files respected | inspected `book-evidence-v4` before cleanup | both clean; nothing deleted or overwritten |

## Judgement

1. **The root cause is correctly identified.** Two distinct causes were present — raw text inside
   backticks, and a `#name` inside a string argument — and fixing only the first would have left the
   Chapter 18 heading broken. The fix uses `#raw(...)` where the page wants monospace and concatenation
   where it needs a string.
2. **The right gate was added, and it is the one that generalises.** The source rule is useful and cheap,
   but the gate that would have caught all six occurrences is the assertion on the rendered page text.
   Every previous gate in this book checked structure and all passed while the page was wrong; this one
   reads what the reader reads.
3. **The negative controls are the evidence.** Two were run, each targeting a different gate, and the
   `{{` case specifically demonstrates that the page gate catches a token the source rule deliberately
   ignores. Reverting was verified (`git diff --stat` shows only the intended patches).
4. **The review exchange is recorded rather than settled in chat**, including the disagreement that was
   not conceded: `v3-06`'s statement stands on the configuration evidence, and the limit is what changes.

## Where this task is weak

- **The token list is a denylist.** `#registry-`, `#build-`, `#source-`, `${`, `{{` and
  `UNKNOWN_PLACEHOLDER` cover this book's inputs and the reviewer's list; a future input named something
  else, leaked in the same way, would be caught by the source rule (which is name-agnostic for the
  backticked form) but not by the page gate. The pair is deliberately complementary for that reason.
- **`{{` and `${` are plausible in legitimate prose.** A Typst source could legitimately render those
  characters as text; if that ever happens, the fix is a documented waiver, not a weakened gate.
- **The visual check remains textual.** I verified the page text and the manifest, not every rendered
  page as an image.
- **The reviewer's release verdict is still outstanding.** They will judge a release artefact, which does
  not exist until the release task runs; this task deliberately left `book/dist/` alone.

## Required tag

`repo/data-architecture-book-template-guard`, created by `queue integrate` on the integrated state.
