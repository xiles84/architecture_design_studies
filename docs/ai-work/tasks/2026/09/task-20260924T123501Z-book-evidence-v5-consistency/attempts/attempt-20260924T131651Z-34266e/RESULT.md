# Result — book evidence v5: the claim-versus-limit audit

**Task:** `task-20260924T123501Z-book-evidence-v5-consistency`
**Branch:** `book/evidence-v5-consistency`
**Required tag:** `repo/book-evidence-registry-v5`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only; the audit reads committed results, it measures nothing

## What was delivered

1. **A sweep of all 29 active claims**, not just the five reported, for the class "the statement asserts
   more than its limits, trials or anchors support". Five failed; 24 passed unchanged. The per-claim
   verdict table is in `book/evidence/v5/ANALYSIS.md`.
2. **`book/evidence/v5/`** — a new package with `claims.json`, `schema.json`, `confounds.json`,
   `coverage.json`, an attributed `SUPERSESSION_LEDGER.md` recording every before/after, and a signed
   `ANALYSIS.md`.
3. **The five re-scopes**: `v2-10` (dropped an unsupported cost conclusion), `v2-16` (dropped a
   write-path cost the claim's scope disclaims), `v2-07` (a limit no longer cites three runs the claim
   cannot show), `v2-14` (a limit scoped to its range), `v3-06` (statement kept, limits corrected).
4. **A new lint plus its tests**: `validate-v5` fails a claim citing more repetitions than its declared
   trials support, with no corroborating anchor. Six new fixtures, including both negative controls —
   the lint fires on the `v2-07` shape and passes once a corroborating anchor exists.
5. **Fixed a path-dependence bug in `check.sh`** that the v5 work exposed (below).
6. **`CONTEXT.md` and `LESSONS_LEARNED.md`** updated, including the follow-up warning for the two delta
   tasks.

## Verified from committed files

```text
tools/evidence validate-v5 --repo ../..   29 claims, 0 errors
tools/evidence validate     --repo ../..   29 claims, 0 errors, active source book/evidence/v5/claims.json
tools/evidence validate-v4 --repo ../..   29 claims, 0 errors   (the lint is not retroactive)
go test -count=1 ./internal/evidence/     ok (34.8s), 6 new v5 fixtures
check.sh --book book                      OK     (was 26 spurious failures)
check.sh --book <absolute>                OK
book/build.sh --out dist/preview.pdf      50 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                          evidence digest present, banner absent, no unresolved tokens
v4 vs v5, 24 untouched claims normalised  byte-identical
```

### Negative controls

| Control | Expected | Observed |
|---|---|---|
| Restore `v2-07`'s "stable across three runs" beside `Trials: 1` | lint fails, naming the phrase | `cites "three runs" but declares trials 1 whole-run … with no corroborating anchor` (unit test) |
| Add a corroborating anchor to the same sentence | lint passes | passes (the escape hatch works) |
| A chapter naming `book/evidence/v4/claims.json` | `registry-literal` fires | `[FAIL] registry-literal: …/information-placement.typ:54 …`, then reverted |

## The one disposition that differs from the review

The reviewer proposed softening `v3-06` to "is designed to keep one donor's donations together", on the
grounds that the card's limit says the placement labels are intent rather than observation. I kept the
statement instead, because the study's schema makes it a configuration fact:
`sql/reference/y1_colocated/schema.sql` declares `PRIMARY KEY ((person_id) HASH, donated_at DESC,
donation_id ASC)` against `PRIMARY KEY (donation_id)` in `y2_noncolocated`. One donor's rows hash to one
tablet by construction. The defective part was the *limit*, which mentioned the missing observation
without saying what the mapping rested on; it now says both. Softening the statement would have
weakened a claim the configuration supports, and `TestV5KeepsTheV306Statement` guards against a later
pass doing it anyway.

## A defect found and fixed on the way

`book/evidence/check.sh` excluded hits under its own directory by comparing scanned paths against an
absolute `$HERE`. Run with a relative `--book` — as I did by accident — the scanned paths were relative
too, the comparison failed, and the checker reported 26 failures that were all its own documentation.
The rule's *scope* had silently changed with the invocation. Canonicalising the scan root fixes it in
both forms.

Recorded because it nearly shipped: **my first attempt at that fix suppressed the negative control it
was supposed to satisfy**, because the exclusion pattern matched the *content* of a matched line rather
than its path. The control caught it. A checker's exemption deserves its own negative control.

## Where this task is weak

- **The audit is one analyst's reading.** Two of the three questions are natural-language judgements;
  only the repetition-versus-trials question is now mechanical. A second analyst reading all 29 claims
  might find a scoping overreach I did not.
- **The new lint is deliberately shallow** and the analysis says so: it cannot see a cost conclusion
  inside a refusals claim, which is exactly why `v2-10` and `v2-16` needed a reader.
- **The `v2-07` repetition evidence is left recoverable, not recovered.** Two sibling Study 02 result
  directories exist; verifying whether each measured the same storage delta is a study-level
  investigation, and I chose the honest reword over an unverified anchor. The analysis names it.
- **`inner_iterations` on `v2-08`/`v2-09` is flagged and unchanged**, because deciding whether it is a
  misuse needs the Study 02 protocol definition. Guessing there would have repeated the failure being
  repaired.
- **The registry and the prose now disagree in three places** (the prose still paraphrases what v5
  retired). This is by design — prose is owned by `book-cache-chapter-v3` and `book-review2-delta`, and
  the release runs after both — but a build taken before those land would pair a corrected claim with
  the old sentence. `CONTEXT.md` states it as a blocking follow-up rather than a note.

## Conclusion

The registry now says about each of its 29 claims only what that claim's evidence supports, the one
shape of overreach that can be detected mechanically is detected mechanically, and the four frozen
packages are untouched. No number changed: this was a wording and attribution task, and the diff against
the v4 tag is the proof.
