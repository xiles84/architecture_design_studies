# Progress — Study 04 v2 staleness field

**Task:** `task-20260923T012236Z-study04-v2-staleness-field`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-staleness-field`

Measured: one cadence cell per maintained-aggregate design (n3, n4), 12 probes each.

## Decision Log

1. **Probe the design's own aggregate through `r04_list_installations`** rather than adding a
   design-specific statement, so the probe works for every design and needs no SQL change. Confidence: High.
2. **Publish a real child change and update the ledger**, so the audit after the cadence phase still
   reconciles and the gate is unaffected. Confidence: High.
3. **Do not probe x2 (drift control).** Its maintenance is out of transaction and its own audit is
   the measurement; probing it with the atomic publish path would time out and say nothing. Confidence: High.
4. **Report the poll interval as the resolution floor.** The observed p50 ~4 ms is dominated by the
   2 ms poll plus one round trip, so the result is "no window above ~5 ms", not "4 ms". Confidence: High.
5. **Keep the probe additive**: a design with no aggregate simply gets no `rollup_lag_ms` field
   instead of a zero that would read as "measured and zero". Confidence: High.

## Residual risks

- Sub-millisecond windows are below the probe's resolution.
- One regime, one scale, PostgreSQL only; 12 samples.
