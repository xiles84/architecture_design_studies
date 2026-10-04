# Review — task-20260924T110000Z-book-text-pass-v2

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, against the committed chapters and the rebuilt PDF.

## What was checked

| Check | Method | Result |
|---|---|---|
| Page 3 describes the real evidence system | rendered text | scope sentence plus direct/mechanism/analogy, the four `gap_kind` values, partially superseded and retired |
| The structure claim is scoped | source and page | "Every *major-family* section…", with variant/scenario/concept named as different |
| Topology no longer claims a pure topology effect | source | layout *package*, consensus scoped to RF=3 with the asynchronous case named |
| Normalization is not tied to overbooking | source | says it decides where a fact lives and needs a constraint plus a strategy |
| Rollup wording is transactional first | source | one transaction, or a named delivery/idempotency/reconciliation mechanism; "publication" reserved |
| Append-only claims answerability, not correctness | source | plus the four things it does not settle |
| Universal costs removed | source | "writes are cheap *in the measured cells of this corpus*" |
| `read score` is defined before use | page order | definition precedes the first multiplier |
| Trial vocabulary is defined | source and page | in-run repeat distinguished from a run replication |
| Ranges agree between prose and cards | rendered text | card ranges are en dashes; zero ASCII-hyphen ranges in card text |
| The build and figure gates still pass | `build.sh`, `export.sh --check` | 51 pages, 29/29 claims, 6/6 fonts, digest present, no banner, no unresolved tokens |
| Nothing else moved | `git diff --stat` | chapters, page 3, two lib files, export.sh, the figure manifest, context and lessons |

## Judgement

1. **The page-3 rewrite is the most valuable change here.** A reader's first encounter with the
   evidence system was a definition that had been overtaken by three registry versions. It now names
   the vocabulary the cards actually use, including the two lifecycle states no earlier revision
   mentioned.
2. **Scoping the consensus claim is the kind of correction that keeps a reference book honest.**
   "Replication turns a write into a consensus write" is true of the measured topology and false in
   general, and the fix says which one the corpus measured.
3. **Separating transaction from publication fixes the sentence the reviewer flagged twice.** The
   original mixed a database rollup with cache semantics; the replacement states the obligation in
   terms of where the derived state lives.
4. **`typeset-ranges` is the right layer.** Editing the registry to contain en dashes would put
   presentation in the evidence record; hand-editing every sentence would drift. The transform sits at
   render and touches no digit.
5. **The `source_revision` fix was found by building, not by reading.** It is the strongest evidence
   that the figure gate does real work: it caught a stale hash produced by a manifest written before
   its own commit.

## Where this task is weak

- **Two criteria are judgement calls, not checkable ones**: whether every family summary separates the
  mechanism from the observation, and whether every forward reference carries a definition. I fixed
  the instances I found and a different reader would find more.
- **The formula behind `read score` is still absent.** The book says what it is and is not; it does not
  publish the weighting, because no published weighting exists.
- **A display transform means the extracted PDF text is not byte-identical to the registry text.** That
  is deliberate, documented, and limited to the hyphen between digits, but a tool comparing the two
  will see it.
- **No reviewer has seen this revision**, and the release artefact will be their first look at it.

## Required tag

`repo/data-architecture-book-text-pass-v2`, created by `queue integrate` on the integrated state.
