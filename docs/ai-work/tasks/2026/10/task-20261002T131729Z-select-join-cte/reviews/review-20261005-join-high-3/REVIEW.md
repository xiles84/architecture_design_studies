# Third independent HIGH review — changes requested (narrow)

Claude Opus 5.5 / Claude Code desktop / HIGH; selected effort unknown; 2026-10-05.
Submission `99233466`, result `a2f3fa9f`, attempt `attempt-20261004T215518Z-975148`.

Read the [full review](../../../../../../../select-query-templates/reviews/20261005-join-high-3/REVIEW.md),
[checks](../../../../../../../select-query-templates/reviews/20261005-join-high-3/CHECKS.md)
and [corrective handoff 3](../../../../../../../select-query-templates/reviews/20261005-join-high-3/CORRECTIVE-HANDOFF.md).

Accepted: the Decision Log, and 37 one-build v3b runs (all valid, controls
firing, digests stable). The YB range index is used (110 vs 200,000 rows
scanned). Classes are arm-complete and limits truthful. The topology report
re-renders byte-identically. Both phases cover the domain matrix. The
colocated database was measured. vet/tests pass and the queue audit is clean.
The domain matrix is final and is not re-run.

Changes requested:

1. The "non-colocated" cluster arm is one tablet per relation with every leader
   on the same node (6/6 runs). It is not distributed, so add a pre-split,
   leader-spread `distributed` arm with a placement gate, and re-run the 30-run
   topology subset under one build.
2. v3's P4 verdict is inverted. YB at the 1% consumer has a replicated sign
   change between n=1 and n=2. PG has no crossing at 1%, and its 100% n=1 point
   is unresolved.
3. Minor reporting fixes; then a v4 successor.

NEXT MODEL: LOW
