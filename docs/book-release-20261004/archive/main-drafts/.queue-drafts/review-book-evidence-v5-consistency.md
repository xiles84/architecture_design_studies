# Review — task-20260924T123501Z-book-evidence-v5-consistency

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, as with the earlier registry tasks. Every check below was run against the
committed files, not from memory of writing them.

## What was checked

| Check | Method | Result |
|---|---|---|
| The package validates | `validate-v5 --repo ../..` | 29 claims, 0 errors |
| The default resolves v5 | `validate --repo ../..` | 29 claims, 0 errors, active source `book/evidence/v5/claims.json` |
| The frozen predecessor still passes | `validate-v4 --repo ../..` | 29 claims, 0 errors — the new lint is not retroactive |
| The lint and its escape hatch | `go test -count=1` | green, 34.8s, six new fixtures including both negative controls |
| Only five claims changed | normalise v4/v5, drop the five, diff | byte-identical for the other 24 |
| Every change is in the ledger | read `SUPERSESSION_LEDGER.md` against the diff | each of the five has before/after text and a reason |
| No number touched | scan the five before/after pairs | no digits added, removed or altered |
| The corrected claims render | extracted page text | the new `v2-07` limit and the `v3-06` hash-key limit both appear; "stable across three runs" is gone |
| The checker is path-independent | `check.sh` relative and absolute | OK both ways, with the `registry-literal` control still firing |

## Judgement

1. **The sweep was the right response to five findings.** Finding five by hand and stopping there would
   have been the cheaper move; auditing all 29 is what turns "the reviewer caught these" into "the
   registry has been checked for this class". The per-claim verdict table makes the negative result
   auditable rather than implied.
2. **`v3-06` is decided on evidence rather than on deference.** The reviewer's proposed softening was
   reasonable from where they sat; the study's `PRIMARY KEY ((person_id) HASH, …)` settles it the other
   way, and the test now guards the decision so a later pass cannot quietly weaken it.
3. **`v2-07` is the honest version of a fix I could have faked.** Adding the two sibling runs as anchors
   would have made the original limit pass while asserting something I had not verified. Recording them
   as a recoverable follow-up is the correct trade, and the analysis says so in the same words.
4. **The lint is small and its limits are published.** It catches the `v2-07` shape and nothing subtler;
   §5 of `ANALYSIS.md` says that rather than implying broader coverage than it has.
5. **The `check.sh` fix belongs here.** The v5 bump depends on the checker behaving the same way however
   it is invoked, and the failure mode (scope silently changing with the invocation) was worth the two
   extra iterations it took to get right, including the one the negative control caught.

## Where this task is weak

- **Two of the three audit questions are judgements**, and I am one analyst. A different reader could
  find a scoping overreach in one of the 24 claims I passed.
- **The registry and the prose disagree in three places right now** (the prose still paraphrases what v5
  retired). It cannot ship that way — the release runs after both prose tasks — but the state exists
  between now and then, and `CONTEXT.md` names it as blocking rather than cosmetic.
- **`inner_iterations` on `v2-08`/`v2-09` is flagged and unchanged.** Deciding it needs the Study 02
  protocol definition; changing it on a hunch would have been exactly the defect being repaired.
- **`v5/confounds.json` is a byte copy of v4's** with no version bump inside the file. That matches how
  v3→v4 was done, but the file's own `schema_version` now trails the package it belongs to.

## Required tag

`repo/book-evidence-registry-v5`, created by `queue integrate` on the integrated state.
