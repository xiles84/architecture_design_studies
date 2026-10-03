# Check report — Study 01 SELECT harness (selectbench)

| Field | Value |
|---|---|
| Task | `task-20261002T131729Z-select-harness` (LOW executor) |
| Attempt | `attempt-20261003T014215Z-c6df8d` |
| Date | 2026-10-03 |
| Host image | `docker.io/library/golang:1.26-bookworm` (go1.26.8) |
| Module | `charitytree` (`studies/01-charity-tree/go.mod`) |
| New package | `studies/01-charity-tree/selectbench/` (additive; the historical harness is untouched) |
| Benchmark lock | not required: no database is started and no measurement is taken (pure Go unit/fixture checks) |

## Commands and results

```
podman run --rm -v <study>:/src -w /src golang:1.26-bookworm sh -c '
  export PATH=$PATH:/usr/local/go/bin
  gofmt -l selectbench/     # empty: formatted
  go vet ./selectbench/...  # clean
  go build ./...            # BUILD_OK: the existing harness still compiles
  go test -v ./selectbench/...'
```

```
--- gofmt (expect no files) ---           (none)
--- vet ---                               VET_OK
--- test ---
TestCatalogLoadsEveryStatementWithAGroup   PASS
TestEveryStatementHasAnEvaluator           PASS
TestGroupGuardRefusesCrossGroupAndControls PASS
TestGroupsAgreeWithOracles                 PASS
TestContractsAreDistinctOnTheFixture       PASS
TestPaginationKeysetMatchesOffsetPage      PASS
TestControlsAreDetected                    PASS
TestFixtureExercisesEveryRequiredEdgeCase  PASS
TestNormalizeStripsVolatileMetadata        PASS
TestSignatureIgnoresCostsRowsAndTimings    PASS
TestSignatureChangesWithAccessPath         PASS
TestStructureIgnoresPredicateButSignatureKeepsIt PASS
TestStructureIsOrderIndependent            PASS
PASS  charitytree/selectbench  0.004s
```

The last five cover plan-signature normalization (`plan.go`): costs, row estimates, actual
times, buffer counts and print order do not change a signature, while a different access path
(Seq Scan vs Index Scan) does, and a different predicate changes the full signature but not the
shape. A plan difference is therefore reported only when the plan really changed.

## What the checks establish

- **Semantic groups cannot be pooled accidentally.** Every catalogue statement carries a group id;
  `Catalog.Comparable` refuses a pair whose groups differ and refuses any pair involving a
  deliberately-wrong control, with an explaining error. `TestGroupGuardRefusesCrossGroupAndControls`
  exercises all four refusals and the unknown-name case.
- **Each deliberately-wrong read-path fixture is detected.** `TestControlsAreDetected` shows all four
  controls differ from their contract oracle: the `WHERE` variant drops preserved parents, nullable
  `NOT IN` collapses to the empty set, name-level `DISTINCT` merges two identities, and limiting before
  the join truncates the bag.
- **Independent Go oracles.** `TestGroupsAgreeWithOracles` checks every valid template against a
  from-scratch oracle for its contract (bag, membership, absence, optional match, aggregate, top-N,
  ordered page), and `TestPaginationKeysetMatchesOffsetPage` shows keyset and offset agree under the
  declared total order.
- **The fixture exercises the required edge cases.** `TestFixtureExercisesEveryRequiredEdgeCase` fails
  if the data lacks a NULL amount, a NULL child key, duplicate qualifying children, an unmatched parent
  or an all-miss set — a passing test on inadequate data would prove nothing.
- **Old Study 01 execution remains reproducible.** `go build ./...` compiles the historical harness
  unchanged; `selectbench` is a separate package with no edits to existing files.
- **The comparison report cannot carry an interpretation, pool a control or rank a failed cell**
  (`report.go`, 6 tests): `RenderComparison` refuses a non-empty interpretation, names
  invalid/unsupported/timed-out cells and excludes them from ranking, refuses cross-group and
  control pairs, recognises two runs of the same plan as `same plan structure`, resolves a
  difference only when it exceeds the cells' spread, and otherwise reports `no resolved difference`.

## Not covered by this check (recorded as remaining work)

The DB-backed catalogue execution, plan capture/signatures, provenance JSON, the comparison report and
the `-cmd select` entrypoint (with `ads-run-lock`) are **not implemented yet**. The oracles and controls
here are the correctness floor those must build on; no number in this report is a performance result.
