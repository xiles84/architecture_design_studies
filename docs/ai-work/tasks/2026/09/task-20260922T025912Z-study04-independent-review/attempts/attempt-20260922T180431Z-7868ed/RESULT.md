# Result — Study 04 independent HIGH review

**Task:** `task-20260922T025912Z-study04-independent-review`
**Branch:** `repo/study04-independent-review`, base after `repo/study02-03-v2-high-validation`
**Required tag:** `repo/study04-independent-review`
**Capability:** HIGH / reviewer (model `deepseek-flash`)
**Benchmark:** none — artifact-only review

## What was delivered

1. A signed independent review analysis
   (`studies/04-configuration-portal/reports/analyses/20260921T1215Z-small--deepseek-flash--2026-09-22.md`),
   indexed under the run's report.
2. A machine-readable claim-verdict file
   (`book/evidence/reviews/20260922-study04-independent-review.json`).
3. The report regenerated so its index lists the review (only the `Generated` line and one analyses
   row changed).
4. `CONTEXT.md` Study 04 status and open-gap 3 updated.

## Verdicts

| Claim | Verdict |
|---|---|
| `claim-10-lock-vs-version-check` | **narrowed** — keep 93.7% control loss and 6.85x trigger cost; the 1.76x lock advantage is an upper bound (single trial, one key/16 writers, PostgreSQL only, `-retries 1` control not run) |
| `claim-gap-01-study04-missing-designs` | **accepted** — the gaps are already owned by `task-20260922T025912Z-study04-v2-protocol`; the review requires three added controls |

## Independently re-derived figures

- `x1`: 3 416 acknowledged, counter 215, 3 201 lost = 93.7%.
- `c2` 887 vs `c1` 503 ops/s = 1.76x (upper bound); `c1` 545 retries / 10 560 conflicts.
- Rollup: `r04` 867 → 13 197 (`n3`, 15.2x) / 13 693 (`n4`, 15.8x); `w05_replace_all` 573 → 83.7
  (`n3`, 6.85x, p99 2.5 → 42 ms), 593 (`n4`).
- Storage: `n1` `config_entry` 1 032 192 B vs `d1` `config_document` 98 304 B = 10.5x.
- Identical-SQL read spread: `r01` 9 148 → 15 141 (1.655x), `r02` 16 308 → 26 801 (1.643x),
  `r03` 25 114 → 40 010 (1.593x), `n1` slowest in all three.

## Material findings beyond the original analysis

- The ~1.7x attribution floor is real and must be preserved, but the proposed cell-order cause is
  **not corroborated by the recorded order** (slowest identical-SQL cell is 8th of 10; the last two
  cells beat cells 6–8). The repeated-design control is required to settle it.
- The analysis's "65% spread" label comes from the `r01` comparison (1.655x) while its parenthetical
  example is `r03` (1.593x); presentation only.
- Missing designs/topologies are coverage gaps, never null results.

## Requirements handed to `study04-v2-protocol`

1. Repeated-design instrument control (`n1` first and last, ≥3 trials, fresh database per cell or
   recorded randomized order).
2. `-retries 1` contention control to separate the lock's advantage from harness charging.
3. An explicit statement that 1.76x is one point on a writer sweep, not a curve.
