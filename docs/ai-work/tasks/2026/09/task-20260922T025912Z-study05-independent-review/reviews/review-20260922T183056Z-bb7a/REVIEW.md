# Review — Study 05 independent HIGH review

**Task:** `task-20260922T025912Z-study05-independent-review`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Reviewed:** `attempt-20260922T182234Z-b46fb1`, submitted at `fde11dd`.
**Decision:** **approve for integration.** No corrective handoff.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | Re-derive from result JSON | Accept — it produced the identity-swap finding. |
| 2 | Record the swap, do not edit the signed analysis | Accept; matches the repository's immutability rule and the review-file pattern. |
| 3 | Accept `claim-11`/`claim-12` with scope, do not broaden | Accept; the required scopes are explicit and correct. |
| 4 | Redis cache-only scope requirement | Accept and important: no cell measures Redis as a store of record. |
| 5 | No new task; add four requirements to `study05-v2-protocol` | Accept; the protocol already depends on this review. |
| 6 | Do not edit `claims.json` v1 | Accept. |
| 7 | Cite the superseded discussion only for mechanism | Accept. |

All decisions are local, reversible and within the brief.

## Result reviewed

- The signed analysis is present, scoped, and names its weaknesses; the four v2 requirements are
  concrete and do not authorize measurement.
- The identity-swap correction is reproducible: `l:na:red:wth:rel` instances_3 = 35 287/40 449 =
  87.24 %, `o:opt:red:wth:rel` instances_3 = 38 749/44 126 = 87.81 %, matching the report's own
  `o:opt:red:wth:rel` line. This is a genuine defect in a signed analysis, handled correctly by
  recording rather than editing it.
- The re-derived throughput ratios match the all-green analysis: 2.60x/2.39x (legacy),
  3.51x/2.86x (owned), 2.24x (owned-pess).
- `CONTEXT.md`, the report index row, and the new `LESSONS_LEARNED.md` entry are accurate and small;
  the report regeneration changed only the generated-at line and the analyses row.

## Conditions on the outcome

`claim-11`/`claim-12` keep their scopes; the corrected model attribution and the Redis boundary must
be carried into any book use; the four v2 requirements land with `study05-v2-protocol`. No
measurement is authorized by this review.

## Integration

Approve and integrate into local `main`; tag `repo/study05-independent-review`.
