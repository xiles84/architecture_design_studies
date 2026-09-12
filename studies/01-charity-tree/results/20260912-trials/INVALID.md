# Partially invalid — do not use the `delete` numbers in this run

**Invalidated:** 2026-09-12 (UTC). **Replaced by:** `results/20260913-writes-isolated/`.

The `delete` measurement in every cell of this run is a harness artefact, not a database
result: ~111 million errors and a reported 0 ops/s per cell.

**Cause.** Deletes consume a finite pool of rows, but were measured by *duration*. 3 s of
warmup plus 5 × 8 s of trials needed ~258 000 deletes against a 108 081-row table; the
pool ran out partway through, and every later trial measured nothing but instant
"no such row" failures in a tight loop.

**What is still valid.** Reads, `insert` and `update` are unaffected. They were, however,
still measured with write phases sharing one table (insert → update → delete), which the
isolated re-run removes — prefer the re-run for any write comparison.

This file is kept, not deleted, so any conclusion ever drawn from it can be traced. See
LESSONS_LEARNED.md → "Finite pools must be measured by count, not by duration".
