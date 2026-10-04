# Containerized source check — SELECT order correction

Generated from recorded facts only; no separate database was started for this
check (the matrices themselves already ran under `ads-run-lock`).

Every order-corrected run builds its benchmark image with the source check as a
build step:

```
STEP 9/9: RUN go vet ./... && go test ./... && CGO_ENABLED=0 GOFLAGS=-trimpath go build -ldflags="-s -w" -o /out/bench ./harness
?   	charitytree	[no test files]
ok  	charitytree/harness	0.172s
ok  	charitytree/selectbench	0.384s
```

- Toolchain: `docker.io/library/golang:1.26-bookworm` (the pinned builder in
  `infra/versions.env`).
- The 18 runs all completed this step (18 build logs in the run chain), so the
  order machinery, the report renderers, the manifest command and every existing
  test pass inside the container, not only on the host.
- The order-specific tests are `selectbench/order_test.go`:
  `TestPlanOrderIsAPermutationAndDeterministic`,
  `TestOrderSeedDependsOnRunID`,
  `TestPlanSequenceReportsTheExecutedAxesInOrder`,
  `TestInterleavedMatchesSequentialStatusAndTrialCounts`,
  `TestInterleavedSkipsCellsWithNoTrials`.
- The final image is ~a few MB from `scratch`; the client is the only process in
  it and it holds no database driver beyond the pinned `pgx`.

Limits: this check runs the same source the runs used; it is not an independent
reimplementation, and it exercises the fake querier for the cell loop rather than
a database (the live oracle gate runs in each matrix).
