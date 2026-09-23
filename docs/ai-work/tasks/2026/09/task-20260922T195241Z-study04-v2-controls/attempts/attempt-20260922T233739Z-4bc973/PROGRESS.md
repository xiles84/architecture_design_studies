# Progress — Study 04 v2 controls

**Task:** `task-20260922T195241Z-study04-v2-controls`
**Attempt:** `attempt-20260922T233739Z-4bc973` (claim `claim-b3e6503a89fd4162`, epoch 1)
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-controls`

Measured: eight valid PostgreSQL cells. One benchmark lock at a time; no parallel measurement.

## Decision Log

1. **No harness change was needed.** `-retries`, `-trials`, `--order-seed` and `--concurrency`
   already exist, so the controls run on the committed Study 04 code (`f1d9b3c`) rather than a
   modified harness. Confidence: High.
2. **Control A is measured as cross-phase reproducibility, not within-run position.**
   `n1_rows_indexed` ran alone (fresh schema and load) in three separate phases. The runner cannot
   place the same design at both ends of one phase without overwriting its own result file, so the
   within-run first-versus-last effect is recorded as the open half of the control. Confidence: High.
3. **The first contention attempt was discarded, not massaged.** Running `verify,explain,contention`
   failed the audit (`INV-3/INV-10: missing key`) because the cell's shared ledger expects the state
   the skipped phases produce. The five failed runs are kept and excluded; the controls were re-run
   with `verify,read,write,contention`. This is the correctness gate doing its job. Confidence: High.
4. **Control B and C use one contention measurement per point**; the ≥3-trial rule applies to
   Control A's reads. A repeated sweep is named as the next step. Confidence: High.
5. **Run tags were created manually for the seven later runs**, because the runner's dirty check
   counts untracked generated `reports/` files as dirty although no tracked code changed. All runs
   record the same commit `f1d9b3c0…`; the tags are on it. The runner defect is recorded in
   `LESSONS_LEARNED.md`. Confidence: High.

## Results (summary; full analysis signed separately)

- Control A: identical SQL reproduces within 0.9%–6.7% across three fresh-load phases.
- Control B: retries-disabled optimistic 731 ops/s vs default 552 at 16 writers; lock 870 → 1.19x
  against retries-disabled, 1.57x against default.
- Control C: c2/c1 = 0.89 (1 writer), 1.26 (4), 1.58 (16), 1.60 (64); all cells `lost_updates = 0`.

## Residual risks

- The within-run identical-SQL position effect (the v1 ~1.6x spread) is still unmeasured.
- Single contention measurement per sweep point; PostgreSQL only; one key.
