# HIGH correction decision log

GPT-6 via Codex desktop; HIGH declared by owner; selected effort unknown.

Resumed the canonical task after its completed-work handoff, expired lease and
ten-minute CAS recovery grace. The new queue image persisted epoch 4 correctly.
Merged current local main in this task's own worktree; earlier worker progress,
results and report remain unchanged.

Reviewed the prior Decision Log first. Accepted additive packaging, semantic
groups, NULL representation, fixed total order, tested fixtures and containerized
checks. Corrected material gaps:

1. Preserve SQL row arrival order for ordered pagination and top-N. Sorting in
   verification hid wrong ORDER BY. Fix keyset mixed DESC/ASC comparisons, tie
   continuation and NULL cursors; include first-page keyset rather than skipping
   it. Confidence High. Alternative: skip keyset, which leaves the contract hole.
2. Use inner lateral top-N: the oracle requires no row from a parent without
   children. Add explicit top-N output order. The oracle contract was retained.
   Confidence High. Affects aggregate/top-N measurement workers.
3. Run four wrong SQL controls on the live teaching fixture, requiring successful
   queries and identical column shapes. Errors/type differences cannot count as
   detection. Change the name-dedup/early-limit controls to compatible projections.
   Confidence High. Affects all SELECT workers.
4. Verify once before EXPLAIN/timing and verify every drained sample. A later
   mismatch discards all samples. Measure three client round-trip/drain samples,
   capture server planning/execution separately, and report decoded JSON size as
   such. Unknown timings are N/A; one whole run establishes no resolved winner.
   Confidence High. Affects reports and synthesis.
5. Retain join hierarchy/outer-inner order in signatures instead of sorting node
   inventories. Confidence High. Alternative: a node-set fingerprint cannot prove
   identical structure. Affects plan explanation.
6. Replace unchecked runtime provenance with clean-source run tags, unique client
   images, named environment/settings, resource/cgroup snapshots and fail-closed
   locked entry. Distinguish unexpected engine failures from timeouts and supported
   feature gaps. Confidence High. All historical artifacts remain immutable.

No large matrix or performance conclusion is part of this harness validation.

NEXT MODEL: HIGH
