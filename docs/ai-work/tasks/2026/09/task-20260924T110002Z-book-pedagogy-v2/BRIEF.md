# Execution brief — Add the reader-facing framework: part opener, glossary and the first visual tranche

Implement the remaining findings of the external review of the *Edition 1* draft that belong to this
task. The review is verified: R02's severity was rejected (the committed artefact carries real
provenance), R01, R03 and R04 are already fixed and integrated.

## Goal

Add the reader-facing framework: part opener, glossary and the first visual tranche. Read `task.json` for the acceptance criteria; they are the contract.

## Decisions already made (do not re-litigate)

- **No new numbers, no new measurements.** Every figure stays attached to an existing registry claim
  in `book/evidence/v4/`, and nothing here may add or restate one.
- **The registry is not yours to edit.** v4 is active and owned by
  `task-20260924T102917Z-book-evidence-v4`. If you believe a claim needs changing, escalate.
- **Do not touch `book/dist/`.** The tracked PDF is regenerated once, by the release task, from fully
  merged sources. Build to `book/build.sh --out dist/preview.pdf --manifest dist/preview-manifest.json`
  and delete both before committing.
- **Anything that cannot be sourced to committed evidence is a gap, not a claim.** Label it as such on
  the page rather than smoothing it into the prose.
- **Do not edit another analyst's report, analysis or digest**, and do not edit the frozen evidence
  packages.

## Constraints

- Every build runs through `book/build.sh` (Podman only); it takes its own `ads-book-build-lock` and
  never touches the benchmark lock. This task measures nothing.
- If the source checker reports a violation, fix the cause or add an owned waiver line to
  `book/evidence/check-waivers.txt`; do not widen the rule.
- Stage explicit paths; never `git add -A`.
- Keep `CONTEXT.md` and `LESSONS_LEARNED.md` current as part of the work.

## Environment note

In this WSL + Windows-podman setup the queue CLI's own git commits can leave `.git/index.lock` behind,
which makes the next queue command fail with "another git process" when none is running. Check
`ps -ef | grep git` and `podman ps` before removing it, then remove it and retry.

## Escalate only if

- an accepted finding cannot be implemented without changing a published claim, a report or an
  analysis;
- two accepted findings conflict with each other or with a documented rule.
