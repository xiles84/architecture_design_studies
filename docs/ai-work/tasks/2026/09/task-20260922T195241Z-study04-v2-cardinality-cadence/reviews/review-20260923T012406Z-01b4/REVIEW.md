# Review — Study 04 v2 cardinality and cadence

**Task:** `task-20260922T195241Z-study04-v2-cardinality-cadence`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration**, with the one unmet acceptance clause moved to a published
follow-up task.

## Decision Log reviewed

1. One run per cardinality tier — accept; the runner has a single `--entries-per-ip`.
2. Cadence regimes mapped to `minutely`/`secondly`/`secondly+burst 30` — accept.
3. All nine cells gate-passed — accept; checked `failed_cells: []` and `gate.passed`.
4. The cardinality hypothesis was tested, and the unindexed scan cost grew 62x — accept; the phase
   correctly continued rather than stopping.
5. Staleness/rollup lag not recorded — accept the honest gap, and require it as a follow-up.
6. Runs auto-tagged — accept; the runner fix from `study-04/v2-controls` is present.

## Independent checks

- Recomputed from result JSON: r04 row designs 11205 → 188 (n0) and 11773 → 189 (n1); rollups
  13223–14709 flat across tiers. Matches the analysis.
- Cadence: quiet 5/5, steady 301/301, burst 300/300, dropped 0, queue depth 22 of 64. Matches.
- `w05_replace_all`: trigger 10.0/11.0 ops/s vs application 413.1/447.9 in quiet/steady. Matches.
- All nine manifests: `failed_cells: []`; all cells `gate.passed = true`.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Cardinality tiers 1…500 for the four designs | met |
| Unindexed scan cost grows with tier, else stop | met (grew 62x) |
| Cadence regimes at 120 entries | met |
| Staleness and rollup lag for trigger vs application | **not met** — no derived-value-age field; owned by `task-20260923T012236Z-study04-v2-staleness-field` (published) |
| ≥3 trials per cell | met |
| Report + signed analysis | met |

## Conditions

The staleness clause is a coverage gap, not a wrong number. It is a targeted follow-up task, so this
task may integrate; the book must not claim a staleness or rollup-lag measurement for Study 04 v2
until that task runs.

## Integration

Approve and integrate; tag `study-04/v2-cardinality-cadence`.
