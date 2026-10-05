# Result — query-plan appendix (Edition 3 preview 6)

**Task:** `task-20261005T190500Z-book-plan-appendix`. **Attempt:** `attempt-20261005T190559Z-87d0aa`.
**Agent:** Claude Opus 5.5 / Claude Code desktop / HIGH / executor; effort xhigh.

## Delivered

- **Appendix:** *Query plans behind the SELECT claims*, pp. 75–87, with 45 trimmed `EXPLAIN ANALYZE`
  plans in 25 exhibits and a reading note per claim.
- **Coverage:** every measured SELECT claim. `v6-sel-07` has no exhibit because its control records
  rows, not a plan, and the appendix says so.
- **Generator:** `tools/evidence plans`. It quotes only runs that are anchors of their claim and whose
  digest matches, and it recomputes every node-tree statement. `go test` fails when the data is stale.

Validation and labelled self-review: `docs/book-plan-appendix-20261005/VALIDATION.md`. Evidence
validate: 50 claims, 0 errors. The canonical build verified 112 pages and 50/50 claims, PDF sha256
`c88b7247…`.

## Decision Log

1. **Scope: the Study 01 SELECT claims only.** These are the claims the owner has been reading.
   Studies 02 and 03 also store plans and are a follow-up.
   - Confidence: High.
2. **One anchored run per exhibit, usually run 1 of the set.** A plan is one instrumented
   execution, and the appendix says the claims' ratios come from timed samples.
   - Confidence: High.
3. **Trim, don't print the full text.** YugabyteDB `DIST` plans run to 96 lines, mostly per-request
   RocksDB metrics.
   - The rule is fixed in code and printed in the appendix.
   - The untrimmed text is located by result path and axis.
   - Confidence: High.
4. **The generator lives in `tools/evidence`, not a one-off script.** It reuses the registry loader,
   and its tests give the staleness gate.
   - Confidence: High.
5. **The colocated pair is expected to differ.** The generator found the colocated run used
   primary-key index scans on the small tables. The note reports this and claims nothing more.
   - Confidence: High.
6. **`v6_test.go` was gofmt-aligned.** The change is whitespace only; it was left over from the
   synthesis task.
   - Confidence: High.
