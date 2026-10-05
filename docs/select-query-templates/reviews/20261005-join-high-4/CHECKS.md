# Checks — fourth HIGH join/CTE review

Claude Opus 5.5 / Claude Code desktop / HIGH; selected effort xhigh; 2026-10-05.
Read-only: no database or benchmark was started. `ads-run-lock` was absent and
no containers were running when the review began.

## Commands

```bash
podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
  docker.io/library/golang:1.26-bookworm \
  go run ./docs/select-query-templates/reviews/20261005-join-high-4/audit.go
podman run --rm -v "<worktree>:/repo:ro" -w /repo/studies/01-charity-tree \
  docker.io/library/golang:1.26-bookworm sh -c 'go vet ./... && go test ./...'
# topology report re-render, twice, producing image, a copy of the 30 runs mounted read-only:
podman run --rm --network none -v "<copy-of-30-v4-runs>:/results:ro" -v "<out>:/out" \
  localhost/charitybench-select:select-join-v4-topo-single-i1-20261005 \
  -cmd join-report -results /results -report-out /out/r1.md
tools/queue/queue audit
```

The full audit output is [AUDIT-OUTPUT.md](AUDIT-OUTPUT.md).

## Results

- **Provenance.** All 30 `select-join-v4-topo-*` runs have `repo_commit`
  `ba821e7b`, one client image id (`b6ea3c007cf4…`) and one engine image id
  (`561d5e6d…`). All 30 producing tags are annotated and point at `ba821e7b`.
  Budgets: node 2 CPU/3 GiB (6 CPU/9 GiB equal-total), client 1 CPU/2 GiB.
  Code after `ba821e7b` changes only the runner's derived `leader-summary.txt`
  capture, not the client.
- **Validity.** 30/30 cells valid in every run, 5 samples in every cell, both
  negative controls fire in every run, and every per-run report contains
  `sha256(join.json)`.
- **Aggregate inputs digest.** sha256 over [V4-INPUTS.txt](V4-INPUTS.txt)
  (one `run sha256(join.json)` line per run, sorted; LF bytes as committed) =
  `acc5d12126ae90644bbca88ef4323ff6dfc063f7400c1e70752aacc448ef3c18`. The
  files are `eol=lf`, so the committed blobs hash identically.
- **Placement gate.** All six `distributed` runs have three tablets for every
  relation, `placement-prepared` and `placement-gate.ok`, the gate note, and
  donation leaders on yb-n1+yb-n2+yb-n3 **before** the cells. The post-run
  `placement.txt` shows every relation still led on all three nodes.
  `noncolocated-single-tablet` and `colocated` runs have one tablet per
  relation (colocated: one shared tablet) led on **one node per run**. The
  node varies: tablet arm yb-n2 in 3 runs and yb-n3 in 3; colocated arm yb-n3
  in 4 and yb-n2 in 2.
- **Report.** Both re-renders and the committed
  `reports/select-join-v4-topology-20261005.md` hash to
  `c2e683090d40e3cee51733ab7ac0cd8a12243d50eb40467d1a301f4dfb7aae9a`. The
  report renders, per run, tablets, donation leaders, endpoint and writer
  operations, retries, overlap and node/client throttling.
- **Blended writer.** 113–175 offered = completed, 0 failed, 0 ledger
  mismatches; cluster writer and reader operations rotate evenly over three
  endpoints (e.g. reader 168/168/168, writer 41/41/40).
- **Throttling.** Single-node 2 CPU node: 71–92 throttled periods (5.2–6.4 s)
  per run; equal-total node 0; cluster nodes 0–26 per node. The 1 CPU client
  is throttled in 32–46% of its CFS periods in every run (e.g. 161/347 in
  distributed-b1).
- **Shard flag.** All 30 runs record `ysql_num_shards_per_tserver=2`. It is
  read from the container's arguments (`podman inspect`), so it is the
  configured value, not one read back from a running tserver. Unsplit
  hash-keyed tables (`integer PRIMARY KEY`) report one tablet on three
  tservers, where 2 × 3 = 6 would follow if the flag applied. Cause:
  **unexplained**.
- `go vet` / `go test`: PASS (`charitytree/harness` 0.150 s,
  `charitytree/selectbench` 0.257 s).
- Queue: claimed 11:19:13Z, first run commit 11:27Z, submitted 12:37:12Z;
  real clock. Audit: 91 tasks, 0 errors, 2 pre-existing warnings on the
  completed `t2-study02-replication` task.
- Producing and checkpoint tags inventoried before review:
  [TAGS-BEFORE.txt](TAGS-BEFORE.txt).

## Run-range separation (`join_where`; a pair is "separated" when every run of one arm lies beyond every run of the other)

Across-run medians in ms; ratio = a / b.

| pair (a vs b) | phase | t=20 | t=2,000 | t=180,000 | MAT/default at 2,000 |
|---|---|---|---|---|---|
| distributed vs single-tablet | isolated | 3.18 / 2.61 = 1.22, **separated** | 0.73, **separated** | 0.54, **separated** | 8.50 vs 11.98, **separated** |
| distributed vs single-tablet | blended | 1.06, overlap | 0.71, **separated** | 0.71, **separated** | overlap |
| colocated vs distributed | isolated | 0.71, **separated** | 0.78, **separated** | 1.21, **separated** | 16.80 vs 8.50, **separated** |
| colocated vs distributed | blended | 0.75, **separated** | 0.82, **separated** | 1.54, **separated** | **separated** |
| distributed vs single node | isolated | 1.86, **separated** | 0.85, **separated** | 0.71, **separated** | overlap |
| distributed vs single node | blended | 1.66, **separated** | 0.93, **separated** | 0.75, **separated** | **separated** |
| distributed vs equal-total single | isolated | 1.88, **separated** | 1.02, overlap | 0.76, **separated** | — |
| distributed vs equal-total single | blended | 1.43, **separated** | 0.80, **separated** | 0.73, **separated** | — |
| single-tablet vs single node | isolated | 1.53, overlap | 1.15, overlap | 1.31, overlap | **separated** |
| single-tablet vs single node | blended | 1.57, overlap | 1.32, **separated** | 1.06, **separated** | overlap |
| colocated vs single node | isolated | 1.33, overlap | 0.66, **separated** | 0.85, overlap | **separated** |
| colocated vs single node | blended | 1.24, **separated** | 0.77, **separated** | 1.16, overlap | overlap |
| equal-total vs single node | isolated | 0.99, overlap | 0.83, overlap | 0.93, **separated** | overlap |
| equal-total vs single node | blended | 1.16, overlap | 1.16, overlap | 1.03, overlap | overlap |

MAT/default values in this table are medians of per-run ratios (colocated
16.80x, distributed 8.50x isolated). The v4 analysis's 17.24x and 8.31x are
ratios of across-run medians (101.2475 / 5.8743 and 65.3485 / 7.8671); both
readings separate the arms in every run.

## v4 analysis statements checked against the data

| statement in v4 | check |
|---|---|
| Topology table (single 1.71/9.49/108.84; equal 1.70/7.85/100.92; single-tablet 2.61/10.94/142.83; distributed 3.18/8.02/77.05; colocated 2.26/6.28/92.93) | matches the report |
| Splitting effect 1.22x / 0.73x / 0.54x (body) | matches |
| Colocated vs distributed 2.26 vs 3.18, 6.28 vs 8.02, 92.93 vs 77.05; 17.24x vs 8.31x | matches; separated in every run, both phases |
| P4 ladder: YB 1% n=1 0.81–0.94, n=2 1.43–1.52; YB 100% n=1 0.90–1.07; PG 1% 0.17–0.81; PG 100% n=1 0.87–1.05, n=2 1.70–2.06 | matches the per-run v3b domain ratios |
| TL;DR: "Against the single-tablet cluster … 0.85x at 2,000 rows and 0.71x at 180,000 rows but 1.19x at the 20-row point read" | **wrong baseline**: 0.85x/0.71x are against the single node; against single-tablet they are 0.73x/0.54x. 1.19x matches no computed pair (1.22x against single-tablet, 1.86x against the single node) |
| Headline: "the materializing CTE from 1.25x into 0.87x"; body: colocated "materializing plan is 1.81x the single node's" | **ratio of penalties, described as plan cost**: 11.98/9.55, 8.31/9.55 and 17.24/9.55. The forced-MAT plan's own time at 2,000 rows is 1.40x (single-tablet), 0.73x (distributed) and 1.14x (colocated) the single node's |
| "The blended phase moves these comparisons by under ~10%" | **false**: colocated vs single node at 180,000 rows goes 0.85x → 1.16x; equal-total vs single at 2,000 rows 0.83x → 1.16x; single-tablet vs single at 180,000 rows 1.31x → 1.06x; colocated `join_where` at 180,000 rows 92.93 → 128.08 ms |
| "single-tablet and colocated runs record one leader node (yb-n3 and yb-n2 respectively)" | **inaccurate**: one leader node per run, varying between yb-n2 and yb-n3 in both arms |
| Single-tablet cluster "1.53x at a point read and 1.31x at 90% (its one leader carries every relation)" | ratios match, but **not resolved**: overlap at every threshold in the isolated phase |
| Shard flag: "a one-tablet outcome for a relation that was not explicitly split is expected behaviour" | **contradicts the recorded flag** (2 per tserver × 3 = 6 expected); the value is the configured one; cause unexplained |
| `inputs_digest` frontmatter | a description, not a digest; the aggregate digest above identifies the 30 inputs |
