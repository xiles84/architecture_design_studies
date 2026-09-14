# Study 01 v3 — final artifact validation

## TL;DR

- All 261 result cells match their manifests' clean producing commit, immutable image
  and environment; all four annotated run tags resolve to those producing commits.
- All four regenerated reports retain their original input digests and index the new
  final analysis as current. The final analysis is approximately 800 words, with four
  signed companions carrying the detailed comparisons.
- The matched memory run has twelve identical initial dataset summaries, 168 passed
  initial checks, 1,512 passed post-mutation checks and zero read errors.

Author: GPT-6 via Codex desktop, 2026-09-14. Environment:
`host-zenbook-ux5406sa`. Data preservation commit: `e869b62`.
Completion milestone: `study-01/v3-enhancements`.

## Scope and results

| Run | Cells | Inputs digest | Producing revision |
|---|---:|---|---|
| `20260913T124917Z-v3` | 14 | `d2d1cc8380b381f2` | `41a1490` |
| `20260913T125342Z-v3` | 100 | `0bc5900a6d7a5d3a` | `c084c73` |
| `20260913T172624Z-v3` | 135 | `bd95d304cb39e2de` | `9b7aa48` |
| `20260914T104721Z-v3` | 12 | `9ead901cb74c4183` | `d1b18bf` |

The 135-cell report retains six D9 negative-control cache failures. No such cell is
included in its valid-rate summaries. Every other cell passes the reporter's recorded
gate, audit and operation-error criteria. Finite windows, rejected arrivals and the
cache audit's ordered-ID contract remain explicit in the signed discussions.

The final-phase memory medians and spreads were recomputed from the twelve JSON files
and checked against the generated per-trial tables. Other cited figures were reviewed
against the mechanism and stress/deployment reports during analysis. Local links in
the final analysis, new discussion companions, protocol and report index resolve.
Both 2026-09-12 analysts' files have an empty Git diff against
`study-01/v2-before-enhancements`. `git diff --check` passes.

## Report regeneration provenance

Each report was first regenerated with its manifest's recorded benchmark image.
For the three September 13 runs, the validated reporter image from `d1b18bf` was then
applied to retain the added growth-read, arrival-count/retry and audit columns and to
represent unavailable YugabyteDB storage as unavailable. This changes presentation,
not the producing revision or input bytes. Both intervening and preceding report
editions are preserved under `outdated/<run-id>.md.previous-N.md`.

- September 13 recorded image:
  `37a233a73a18324a4749a831a53963a97be6f332b9caf579c7cb6158760e101e`.
- September 14 recorded image and corrected reporter:
  `5f0b982c0a01eb95afcc05d7b859883b8889e765ef8892c8c44c1bf21b609f3b`.

Report rendering ran in Podman with network disabled, 0.25 CPU and 512 MiB; it reads
saved files and starts no databases. Study 03 held the measurement lock during this
documentation continuation. Its containers and lock were left with that run.

No benchmark was rerun for this artifact review. The previously completed containerized
Go race tests, vet, shell syntax checks and reporter regression test are recorded in
[the tooling validation report](20260914-v3-report-validation.md). The new designs'
earlier validation is in [the harness report](20260913-v3-harness-validation.md).
