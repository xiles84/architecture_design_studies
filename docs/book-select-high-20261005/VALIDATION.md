# Edition 3 preview 4 — validation (task `select-book-high`)

Claude Opus 5.5 / Claude Code desktop / HIGH; effort xhigh; 2026-10-05. No database or benchmark was
started.

## Canonical build

`book/build.sh` ran in WSL Ubuntu. It drives the Windows `podman.exe` through `infra/lib.sh`, using
the pinned Typst image `sha256:032e2922…`, the PlantUML renderer `sha256:9b9ee6af…` and the
`localhost/ads-book:1` PDF tools.

| Build | Source commit | Tree | Pages | Fonts | Claims indexed | PDF sha256 |
|---|---|---|---|---|---|---|
| 1 ([log](build-1.log)) | `bc6a2301` | clean | 98 | 4/4 | 49/49 | `3738a8c2…` (superseded) |
| 2, final ([log](build.log)) | `bf03e408` | clean | 98 | 4/4 | 49/49 | `be83c9422389a1b2e92265eae3d7291e7695b8a72119e6674a5cbe8a733b3799` |

Both builds pass the gates:

- source check (`check.sh`), figure manifest (17 figures reproduced);
- evidence digest present in the PDF;
- no unverified banner, no unresolved tokens;
- the edition on the title page is "Edition 3 preview 4".

Build 1's rendered review found the layout problems and the stale figure listed below. They were
fixed and rebuilt; build 1's text and proofs are kept as `rendered-1.txt` and `proof-1/`.

## Reader audit

[`check-reader.go`](check-reader.go) runs in the golang container over `pdftotext` output. All 34
checks pass on the final text ([output](check-reader.out)):

- **Renders:**
  - the six new variations and the template guide;
  - the access-path note;
  - the placement cross-link;
  - the new badge label;
  - the edition.
- **No longer renders:**
  - the unverified banner;
  - "Planned variation";
  - the old "no new SELECT ratios" sentence;
  - the old badge label;
  - the old "does not yet promote" coverage note.
- **Claims:** every one of the 18 `v6-sel-*` ids appears at least twice (badge and card).

## Rendered review (HIGH self-review)

Pages 3, 29, 51, 53, 56–59, 67 and 74 were rendered at 90 dpi and inspected. Final proofs of the
changed pages are in `proof/`.

**Found in build 1, fixed in build 2:**

- Short paragraphs ended in a badge that wrapped alone, which stretched the justified line. Each new
  badge now follows a manual line break.
- `(amount + 0)` split across lines. It is now unbreakable.
- The guide table hyphenated cells and broke `v6-sel-…` ids mid-id. The table is now unjustified with
  unbreakable ids.
- "Study 01's" rendered with a prime. The sentence was reworded.
- **Substantive:** the placement figure (p. 51) still said tablet/leader placement was "UNVERIFIED in
  this corpus" and "cannot be observed with this harness". Its PlantUML source now says placement is
  verified only in the Charity SELECT placement runs. It was re-exported with the pinned renderer, and
  only that asset changed.

**Found in neither build:** clipping, overlapping text, a split table or an orphaned heading on the
inspected pages.

## Evidence and preservation

- The registry is unchanged in this task (`v6`, 49 claims; validated 0 errors by the synthesis task).
- Book numbers are copied from claim statements. No new ratio was computed.
- No signed analysis, report, result, producing tag or frozen evidence package was edited.

## Limits

- One HIGH reader inspected the pages; this is not independent visual review.
- The SELECT numbers carry their claims' limits: one laptop, one sequential reader, scan-only
  YugabyteDB access/existence plans, and bracketed crossings.
