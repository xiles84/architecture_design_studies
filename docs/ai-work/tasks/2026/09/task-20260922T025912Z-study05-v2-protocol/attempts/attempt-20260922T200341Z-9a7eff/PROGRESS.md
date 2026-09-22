# Progress — Study 05 v2 completion protocol

**Task:** `task-20260922T025912Z-study05-v2-protocol`
**Attempt:** `attempt-20260922T200341Z-9a7eff` (claim `claim-00b6ff53e971b287`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `study-05/v4-handoff`

Planning only: no database started, no measurement.

## Deliverables

- `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md` — decision-complete protocol.
- A pointer from `HANDOFF.md` to the amendment.
- Three execution tasks published: `task-20260922T200831Z-study05-v2-churn-scale`,
  `…-resources-engines`, `…-open-loop`.
- `CONTEXT.md` gaps 4 and 5 updated (gap 5 was stale: the book v1 release is done).

## Decision Log

1. **Amend, do not rewrite.** The v1 survey and its numbers stay as the baseline; the amendment adds
   the unmeasured regimes and gates future ratios. Confidence: High.
2. **Churn must cross the TTL.** A churn cell that crossed no refresh/expiry is reported as a gap,
   not a null result. Confidence: High.
3. **Scale must report working-set-to-RAM and eviction.** Throughput without an eviction record is
   not a cache-hit number. Confidence: High.
4. **Equal-total accounting is a separate arm**, never pooled with `db-only`; this is the review's
   `conf-04` fix. Confidence: High.
5. **Placement needs a recorded tablet/leader distribution**, or it is a gap. Confidence: High.
6. **Open-loop demand is its own task** using the platform arrival driver; closed-loop tails are not
   reused. Confidence: High.
7. **Split into three tasks** so each matrix is one benchmark-lock session with its own report and
   signed analysis. Confidence: High.

## Residual risks

- 800k donors may exceed the container budget; the scale task must report the ratio it could
  configure rather than silently caching.
- The arrival driver's saturating behaviour on this host is untested; the open-loop task states the
  offered rate and any shortfall.
