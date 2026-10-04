# Execution brief — Stop template tokens rendering in the PDF, and fail the build if they ever do

Implements findings from the third review round of the *Edition 1/2* draft. The review is verified
against source: every listed finding below was reproduced before being scheduled, and no arithmetic
error is claimed.

## Goal

Stop template tokens rendering in the PDF, and fail the build if they ever do. `task.json` holds the acceptance criteria, and they are the contract.

## Decisions already made (do not re-litigate)

- **Never assert what was not measured.** Where a claim's text exceeds its limits, the unsupported
  assertion is dropped; no new measurement is commissioned by this task.
- **A tagged registry package is never edited in place.** Wording changes to a claim are a new
  versioned package with a supersession ledger. v1-v4 stay frozen and no tag moves.
- **The book renders claim text verbatim.** If a statement needs a different wording, the registry
  changes; the book does not paraphrase a claim into something the registry does not say.
- **Numbers and interpretation stay separate.** A mechanistic sentence is allowed and useful; it must
  be labelled as mechanism, never presented as this corpus's measurement.
- **Do not commit `book/dist/`.** Build to `dist/preview.pdf` + `dist/preview-manifest.json` and
  delete both before committing; the release task rebuilds the tracked artefact once.

## Constraints

- Everything builds through `book/build.sh` (Podman only, its own `ads-book-build-lock`, never the
  benchmark lock). This task measures nothing.
- Never edit another analyst's report, analysis or digest, and never edit a frozen evidence package.
- Stage explicit paths; never `git add -A`.
- Keep `CONTEXT.md` and `LESSONS_LEARNED.md` current as part of the work.

## Environment note

In this WSL + Windows-podman setup the queue CLI's own git commits can leave `.git/index.lock`
behind, and the next queue command then fails with "another git process" when none is running.
Check `ps -ef | grep git` and `podman ps` before removing it, then remove it and retry.

## Escalate only if

- a finding cannot be fixed without changing a number, a report or a signed analysis;
- a limit can only be satisfied by a new measurement, or
- two accepted findings conflict.
