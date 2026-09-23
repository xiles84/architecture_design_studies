# Review — Study 04 v2 staleness field

**Task:** `task-20260923T012236Z-study04-v2-staleness-field`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration.**

## Decision Log reviewed

1. Probe through `r04_list_installations`, no design-specific SQL — accept; it works for every
   design and leaves the catalogue untouched.
2. A real child change with the ledger updated — accept; the audit still reconciled, which is the
   risk this decision carried.
3. x2 not probed — accept; its maintenance is out of transaction and its own audit measures it.
4. The 2 ms poll reported as a resolution floor — accept; this is the difference between a
   measurement and an overclaim.
5. A design with no aggregate gets no field — accept; an absent field cannot be misread as zero.

## Independent checks

- Recomputed from result JSON: n3 p50 4.30 / max 6.88 ms; n4 p50 4.01 / max 4.34 ms; 12 probes each;
  `rollup_lag_timeouts` empty in both.
- Both cells: `gate.passed = true`, `failed_cells: []`, cadence audits passed.
- Dev check `20260923T1300Z-devcheck-stale` passed before the reported run; run tag
  `run/04-configuration-portal/20260923T1305Z-staleness-120` exists on commit `d9f63b0c…`.
- The probe only adds keys it records in the ledger; `git diff` shows no change to the design SQL.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Cadence phase stamps the derived value and records the age | met |
| Staleness and rollup lag reported for trigger vs application at 120 entries | met (both, with the resolution bound stated) |
| Gate passes; change re-verified before publishing a number | met |
| One matrix at a time under the benchmark lock | met |

## Conditions

The book and registry may say Study 04's trigger- and application-maintained rollups show **no
staleness window above roughly 5 ms** at 120 entries on one PostgreSQL host. They must not quote a
sub-millisecond figure, and must not claim the probe characterises a tail.

## Integration

Approve and integrate; tag `study-04/v2-staleness-field`.
