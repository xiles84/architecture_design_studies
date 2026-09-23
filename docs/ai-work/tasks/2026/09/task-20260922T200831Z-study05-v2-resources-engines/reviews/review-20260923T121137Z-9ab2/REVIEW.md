# Review — Study 05 v2 equal-total, engines and placement

**Task:** `task-20260922T200831Z-study05-v2-resources-engines`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration**, with the placement-distribution clause recorded as an open
gap.

## Decision Log reviewed

1. Fix the multi-endpoint defect rather than drop the cluster arm — accept, and the fix is the right
   one: first-endpoint-only would have been a false cluster measurement.
2. Re-run the whole matrix on one coherent commit — accept; pre-fix failed runs are kept and excluded.
3. Equal-total labelled, never pooled — accept.
4. Placement distribution reported as a gap — accept; the honest reading.
5. One host is not colocation evidence — accept; the study report says so.

## Independent checks

- Equal-total vs db-only recomputed from JSON: warm 10,837 vs 8,027; cacheable-99 7,876 vs 8,130;
  app-99 7,945 vs 7,873. `resource_framing` present in both.
- yb-cluster3 cells report `connection_nodes: 3`; yb-single 1. All eight cells `gate.passed = true`,
  `failed_cells: []`, strict violations 0, impossible cache values 0, ledger mismatches 0.
- Colocated/non-colocated recomputed: y1 882/852/558 vs y2 673/557/387.
- `platform/adapters/pgxdb/spread.go` and the Study 05 `splitList` build clean; single-endpoint path
  returns the pool unchanged.
- Run tags exist for the reported cells.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Equal-total arm separately labelled, never pooled | met |
| Strict/relaxed on YB 1-node RF=1 and 3-node RF=3, per engine | met (endpoints spread over 3) |
| Colocated vs non-colocated with tablet/leader distribution recorded | **speed pair met; distribution NOT recorded** (server list only) — gap |
| Gate passes; controls fired; one matrix at a time | met |

## Conditions

The book and registry may quote the equal-total arm, the per-engine strict/relaxed pair and the
colocated-vs-non-colocated speed difference **on one host**. They must not claim evidenced physical
colocation or a recorded tablet/leader distribution. The placement gap stays open for a future run
that can list tablets on this image (or for two hosts, which is deferred).

## Integration

Approve and integrate; tag `study-05/v2-resources-engines`.
