# Result — Study 01 SELECT corrective round 3

**Task:** `task-20261002T131729Z-select-join-cte`
**Attempt:** `attempt-20261005T111913Z-13271e`; claim `claim-901d6661ceb83a37`
**Capability/role:** LOW executor (model `deepseek-flash`, opencode CLI)
**Handoff executed:** `docs/select-query-templates/reviews/20261005-join-high-3/CORRECTIVE-HANDOFF.md`

## Delivered

**Distributed placement arm, gated.** `distributed` pre-splits `charity`,
`person`, `donation` and `donation_person_idx` into three tablets and creates
`donation_rank_idx` as `(rank ASC) SPLIT AT VALUES ((66667),(133334))`. The gate
checks `num_tablets = 3` for `donation` and the rank index from the engine's own
metadata and three distinct tablet leaders through a runner handshake
(`yb-admin list_tablets`, polled up to five minutes, before/after listings
saved). All six distributed runs passed; every relation's leaders are
`yb-n1:9100, yb-n2:9100, yb-n3:9100`. The engine default arm is relabelled
`noncolocated-single-tablet`, and the colocated and single-tablet arms record
their single-leader placement without a spread gate.

**30 topology runs from one build:** single 2 CPU/3 GiB, equal-total 6 CPU/9 GiB,
`distributed`, `noncolocated-single-tablet`, `colocated`, each isolated and
blended, three whole runs, 30 cells, five trials, position control. 30/30 valid
in every run; both controls fire; no invalid, timeout or error cell.

**Report.** `reports/select-join-v4-topology-20261005.md` separates the five arms
× two phases into distinct classes (budget, topology, placement, phase, index
arm) and renders, per run, tablets per relation, tablet leaders, endpoint and
writer operations, and per-node/client throttling. It re-renders
byte-identically. PostgreSQL `storage_rows_scanned` is `n/a`; the scan flags are
named whole-plan; the engine's effective `ysql_num_shards_per_tserver` is
recorded.

**v4 signed analysis** supersedes v3 (moved to `reports/outdated/`;
`REPLACEMENTS-select-join.md` extended). P4 now follows the per-run signs:
a replicated YB sign change between one reference (0.81–0.94, 6/6 below 1) and
two (1.43–1.52, 6/6 above 1) at the 1% consumer — a small (6–19%) but consistent
effect; YB at the full consumer n=1 mixed/unresolved; PG at the 1% consumer no
crossing in 1–16; PG at the full consumer crossing at or below n=2 with n=1
unresolved. Topology claims cite the new report and name the measured placement;
the required pair is colocated vs distributed, and single-tablet vs distributed
is the splitting effect.

**Headline results (across-run medians, `join_where`).** Against the single node
(1.71 / 9.49 / 108.84 ms at 20 / 2,000 / 180,000 rows): the single-tablet
cluster is 1.53x / 1.15x / 1.31x; `distributed` is 1.86x / 0.85x / 0.71x;
`colocated` is 1.32x / 0.66x / 0.85x. Splitting (single-tablet -> distributed)
gives 1.22x slower at a point read and 0.73x / 0.54x at 2,000 / 180,000 rows.
Colocated's forced-materializing plan costs 17.24x its default at 2,000 rows
against 8.31x distributed. The equal-total budget still buys nothing.

## Gate

- `go vet ./...` and `go test ./...` pass in `docker.io/library/golang:1.26-bookworm`
  under Podman.
- `tools/queue/queue audit`: 0 errors.
- Tree clean; benchmark lock free; no containers running; no push; completion tag
  not created.

## Notes / gaps

- Cluster nodes are containers on one host; placement is evidenced by tablets,
  leaders and listings, not by independent hardware.
- Three runs per arm; the distributed point-read penalty is inside the spread.
- The v3b and v3 runs remain labelled cross-build observations.
- No book claim is promoted; HIGH review owns promotion.
