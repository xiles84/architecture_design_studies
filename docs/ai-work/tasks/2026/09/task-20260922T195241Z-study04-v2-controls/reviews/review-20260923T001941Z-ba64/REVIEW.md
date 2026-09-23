# Review — Study 04 v2 controls

**Task:** `task-20260922T195241Z-study04-v2-controls`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration** after one review-required correction (committed `938d420`).

## Decision Log reviewed

1. No harness change needed; controls use the committed Study 04 code — accept; the harness Go and
   platform/infra trees are byte-identical to `f1d9b3c` (`git diff` empty).
2. Control A measured as cross-phase reproducibility, within-run position recorded as the open half
   — accept as a documented partial (see conditions).
3. First contention attempt discarded, not massaged — accept; the five failed runs are kept,
   untagged, and excluded from the analysis.
4. One contention measurement per sweep point — accept; named as a weakness.
5. Run tags created manually on the producing commit — accept, and the cause was fixed (below).

## Review-required correction

The runner's dirty check counted untracked generated `reports/` files, so every run after the first
reported `repo_dirty: true` and silently skipped `--tag`. That is a provenance defect: a run that
cannot be tagged must fail loudly. The reviewer required a fix; `run-study.sh` now filters generated
results/reports explicitly and the tree is clean. The correction is committed at `938d420` and does
not touch the harness, so the control numbers are unaffected.

## Independent checks

- `c2/c1` recomputed from the result JSON: 0.89 (1 writer), 1.26 (4), 1.58 (16), 1.60 (64) — matches
  the analysis.
- Control B: c1 retries=1 731.4 ops/s vs default 552.2; c2 869.7 → 1.19x / 1.57x — matches.
- All eight valid cells: `gate.passed = true`, `failed_cells: []`, `lost_updates = 0`.
- Five reduced-phase runs are recorded as failed with `INV-3/INV-10: missing key`, kept in
  `results/`, quoted nowhere.
- Eight run tags exist on `f1d9b3c0…`.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Control A: identical-SQL floor as a measured range | partly met — 0.9%–6.7% measured across fresh-load phases; within-run position effect explicitly open |
| Control B: `-retries 1` control; 1.76x restated only with it | met — restated as 1.19x–1.57x |
| Control C: 1/4/16/64 writer sweep | met |
| Correctness gate passes; one matrix at a time | met |
| Generated report + signed analysis | met |

## Conditions on the outcome

The book's Study 04 claims may now quote the writer sweep and the narrowed lock advantage. They must
not claim a measured within-run position effect; that control belongs to a follow-up task (add a
runner mode that places one design at both ends of a phase without overwriting its result).

## Integration

Approve and integrate into local `main`; tag `study-04/v2-controls`.
