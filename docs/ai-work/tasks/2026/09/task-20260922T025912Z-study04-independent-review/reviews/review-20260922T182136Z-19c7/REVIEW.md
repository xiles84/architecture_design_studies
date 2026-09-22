# Review — Study 04 independent HIGH review

**Task:** `task-20260922T025912Z-study04-independent-review`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Reviewed:** `attempt-20260922T180431Z-7868ed`, submitted at `6f5dd31`, branch
`repo/study04-independent-review`.
**Decision:** **approve for integration.** No corrective handoff.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | Re-derive from result JSON, not the report | Accept. It is how every figure below was checked. |
| 2 | Keep the ~1.7x floor, separate artefact from explanation | Accept, and it is the review's strongest contribution: the spread is real (1.59-1.66x), the proposed order cause is not corroborated by the manifest order. |
| 3 | Narrow `claim-10`, do not reject it | Accept. Two clauses re-derive exactly; only the 1.76x clause needs the author's own upper-bound caveat. |
| 4 | No new task; extend `study04-v2-protocol` scope | Accept. The protocol already depends on this review and owns the gaps; three concrete controls are added to it. |
| 5 | Do not edit `claims.json` v1 | Accept. Verdicts are additive, matching the phase-3b review pattern. |
| 6 | Report regeneration changes only the generated-at line and one index row | Accept and verified. |

All decisions are local, reversible and inside the brief.

## Result reviewed

- The signed analysis is present, correctly scoped, and its weaknesses section is honest.
- Verdicts are backed by figures the reviewer re-derived: `x1` 3 201/3 416 lost; `c2` 887 vs `c1`
  503; `r04` 867 → 13 197/13 693; `w05` 573 → 83.7; storage 1 032 192 → 98 304 B; identical-SQL
  reads 1.655x/1.643x/1.593x.
- The order check (slowest identical-SQL cell is 8th of 10; the last two cells beat cells 6–8) is
  reproducible from `manifest.yaml` `design_order_pg-single` and the result JSON; it is a genuine
  correction of the mechanism claim, not a preference.
- The three requirements for `task-20260922T025912Z-study04-v2-protocol` are specific and cheap and
  do not authorize measurement.

## Conditions on the outcome

`claim-10`'s first clause must not be quoted without the upper-bound scope; the v2 protocol must
carry the three added controls. No measurement is authorized by this review.

## Integration

Approve and integrate into local `main`; tag `repo/study04-independent-review`.
