# Join/filter theory and case route — validation

Scope: add theory now, amend the existing Study 01 LOW measurement task, publish
a HIGH-minimum book-result successor. No database/calibration/matrix was run.
The SQL is illustrative; engine execution is acceptance work for the LOW study.

## Canonical final build

- Edition 3 preview 3, **89 A4 pages**, 4/4 embedded fonts, 31/31 indexed claims.
- Producing source commit: `b023c889056624e0317159cdb3b1ad76cae0e8de`;
  clean tracked tree at build. Later book differences are only generated dist outputs.
- PDF SHA-256: `2b4f386351bbf156271cdbc1aa1d9bf7339393b535ab41af606ddad5d26838bb`.
- Active evidence SHA-256 unchanged:
  `6c5b2f80e98321c0c7cc90d4d6e42e8ebc91ef66ac6504bf4a732d0c5688e2db`.
- Canonical `book/build.sh` in WSL calls the pinned Podman Typst/PlantUML/PDF tools.
  Existing WSL jq is reused; the Git function converts WSL paths for native git.exe
  to avoid cross-platform index refresh. No host toolchain was installed.
- Source gate passes; all 17 registered assets reproduce. No compiler warnings,
  missing labels, unresolved template tokens or unverified banner.
- Both build transcripts remain: initial `build.log`, repaired `final-build.log`.

## Reader and evidence audits

Containerized Go `check-reader.go` passed on initial and final extracted text:
six theory groups, **22 strategies**, five cases, all six SQL examples and both
LEFT JOIN contracts rendered, 31/31 claims in teaching, and the new HIGH task's
synthesis dependency/supersedes link plus three predecessor amendments present.
The final text is `rendered-final.txt`; the first build's text remains preserved.

Containerized `tools/evidence validate-v6` passed: **31 claims, zero errors**.
The full repository was mounted read-only so shared producing tags remain visible.
Original join/synthesis/book task.json files have zero diff from baseline main;
signed analyses, reports and frozen evidence were not edited.

Queue source image-check passed with digest
`1c6c46eb05da53419a89eed0f04a964344b29e49bb8ebb43fa75a5198d1c5712`.
Audit before integration: **91 tasks, zero errors, two pre-existing Study 02
coordination warnings** (expired lease and live-ref drift). Initial claim drift
was cleared by this task's normal checkpoint. No peer coordination was changed.

## Rendered review

HIGH self-review inspected pages 4, 20-24, 53, 55 and 56 at 1300 pixels. New
diagrams carry conceptual status and text equivalents; SQL blocks are complete,
readable and unsplit. Contents and case/theory links compile. The initial page
56 contained only the case tail; a pagebreak groups the planned variation and
its gap/reproduction note. Final pages 4,55,56 were rerendered and inspected;
unchanged new theory pages retain their inspected proofs in `proof/`.
No clipping, overlap, split query or orphaned checklist row was found.

This is one HIGH reader, not independent visual review. Expectations are labelled;
there is no benchmark ranking or new numeric evidence promotion.

## Remaining execution and coordination limits

LOW must calibrate the amended 200k-row matrix, feature support, client capacity
and topology/placement arms. Infeasible/unsupported arms remain explicit gaps.
HIGH then validates and registers evidence before HIGH book publication.

Queue v1 lacks a task-supersede command. The old LOW publication task retains
historical ready state and LOW metadata, with a mandatory **do not execute**
amendment. The published successor has minimum HIGH and a `supersedes` link.
Workers must read amendments even if the old assignment is automatically offered.

Integration is confirmed by this task's queue receipt and annotated milestone
`repo/book-join-theory-20261004`; no push or pull.

TASK COMPLETE
