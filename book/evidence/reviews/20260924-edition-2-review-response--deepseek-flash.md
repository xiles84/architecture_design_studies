# Review exchange — Edition 2 drafts, rounds 2 and 3

**Attribution.** Response written by `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI; effort
setting not exposed) on 2026-09-24, task `book-template-guard`. The review being answered is
**GPT-5.6 Sol, high effort, 2026-09-24 08:22 −03:00** — the second review of the revised 50-page
draft, plus its follow-up reply.

**Why this file exists.** Two analysts disagreed about the same artefact, and the disagreement is a
finding. It is recorded here rather than left in a chat log, and neither position is smoothed over.

## What the reviewer examined

A **hand-built development preview**, not a `book/build.sh` release artefact. The reviewer confirmed
this in their reply: the uploaded PDF said "Unverified development build" and stated it was compiled
without the build script's provenance inputs. That single fact decides two of the round-2 findings.

## Corrections we made to the review, and their outcomes

| Our position | Outcome |
|---|---|
| The draft's provenance was **not** broken: pages 1–2 carry real commit, describe, dirty state, build clock, image digest, evidence digest and source-tree hash, and `book/dist/build-manifest.json` records them. The banner is a *preview-path marker*, so review 1's R02 was a misdiagnosis rather than a repair. | **Accepted by the reviewer.** "My second review conflated 'this preview correctly identifies missing provenance' with 'the release provenance mechanism was repaired'." Release readiness will be re-judged on a release artefact. |
| Their §18 rationale was wrong: cards are **generated from the registry**, so a re-printed card cannot drift from its claim. The length and reading-flow cost is real. | **Accepted.** Capsules on repeat citations only; the first citation keeps the full card and no limit line is dropped. |
| Their §12 explanation of the diagram spacing was speculative: the cause is PlantUML emitting `textLength` and `lengthAdjust="spacing"`, which Typst honours, not inherited tracking or a font interaction. | **Accepted.** "My source-level explanation for the latter was speculative and yours is based on the actual generator." |

## What the review found that we confirmed and fixed

- **P0, our bug.** Six sites rendered template tokens literally (`#registry-label`,
  `#registry-version`): the title-page provenance line, the provenance table, the Chapter 18 heading, the
  index prose, Chapter 8 and §15.4. Cause: interpolation written inside backticks (raw text is not
  evaluated) and inside a plain string argument. Fixed in `book-template-guard`, with a fatal gate on
  the extracted page text and a source rule in `check.sh`. Note the reviewer's count of four was
  correct *for a preview build*: the title-page line sits inside `#if release … #else [banner]`, so a
  preview shows the banner and hides that leak.
- **Five cache-card overclaims**, four of them authored in this project's own revision: the lock/CAS
  guarantee, the fill-lease "one filler at a time", the TTL refill bound, the unconditional
  "prerequisite" reading of the outbox, and the matrix placing durable observation under mutation
  correctness.
- **A contradiction we introduced ourselves:** "relaxed … for a bounded window" beside a caveat that the
  run cannot establish a bound. Bounded staleness becomes a named subtype, not the definition.
- **Thirteen further findings** in the reply, all verified against source before scheduling.

## Where we were more sceptical, and what the evidence showed

The reviewer flagged five **claim-versus-limit mismatches** and offered a wording or a
registry-version change for each. We checked every one against `book/evidence/v4/claims.json` and
against the study sources rather than taking the offered remedy:

- `v2-10` ("at no measurable cost" beside "the claim is about refusals, not throughput") — confirmed;
  the unsupported assertion is dropped rather than backed by new measurement.
- `v2-16` (a write-path cost asserted by a claim limited to "read throughput only") — confirmed; the
  claim keeps the read result and the write-path sentence becomes a labelled mechanism note.
- `v2-07` (`Trials: 1` beside a limit citing "three runs") — confirmed; the repetition count must
  either be supported by anchors or removed from the limit.
- `v2-14` (a limit that reads as a corpus-wide gap) — confirmed; scoped to the range it limits.
- `v3-06` ("keeps one donor's donations in one tablet" beside "placement labels are intent, not verified
  placement") — **confirmed as a tension, but the reviewer's preferred remedy is not the right one.**
  They offered to soften the statement to "is designed to". The study's own schema settles it:
  `sql/reference/y1_colocated/schema.sql` declares `PRIMARY KEY ((person_id) HASH, …)` while
  `y2_noncolocated` uses `PRIMARY KEY (donation_id)`. A hash-keyed primary key places one person's rows
  in one tablet **by configuration**, not by observation, so the statement stands and the *limit* is
  what must change: logical partition mapping follows from the hash key; physical tablet/leader
  placement on the three nodes was not observed. Softening the statement would have weakened a claim
  that the configuration does support.

Their own summary is worth keeping: they found no arithmetic error, and several revised cards
(`v3-03`, `v3-04`, `v3-05`, `v3-01`) are cited as good examples of the discipline the book teaches.
The remaining concern is scoping around numbers, not the numbers.

## What this exchange changed in the queue

`book-template-guard` (this file's task), `book-evidence-v5-consistency` (the sweep of all 29 claims
plus a new lint for the repetition-versus-trials class), `book-cache-chapter-v3` (the card and
taxonomy corrections) and `book-review2-delta` (the prose, lifecycle-rendering, capsule and visual
deltas). The release task runs last and produces the artefact the reviewer will judge.
