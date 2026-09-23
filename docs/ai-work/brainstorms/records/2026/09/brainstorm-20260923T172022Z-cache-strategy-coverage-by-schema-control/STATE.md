# Current brainstorm state

```text
BRAINSTORM STATUS
ID: brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control
State: ongoing
Stage: collecting_positions
Summary:
  Map caches on two axes: who arbitrates the copy (cache-side fence, DB version token, authoritative read) and whether the schema may change (legacy vs owned), with aside/through, memory/redis and relaxed/strict as parameters. Ownership did not buy consistency (owned 87.81% vs legacy 87.24% wrong reads under relaxed write-through), and the strict arms came from a fence that fires before and after commit. Record soft deletes, outbox-driven invalidation, partitioned process-local caches, hard-TTL expiry and placement as gaps.
  Progress: positions 1/3; critiques 0/2; synthesis pending
Disagreements: Anticipated: relaxed staleness is permitted so it should not be called wrong; one machine and one trial cannot be generalized; version columns, outboxes and soft deletes are standard and may be recommended.
Next action: CONTINUE BRAINSTORM brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control
```
