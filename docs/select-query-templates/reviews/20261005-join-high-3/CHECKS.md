# Checks — third HIGH join/CTE review

Claude Opus 5.5 / Claude Code desktop / HIGH; 2026-10-05. Read-only: no
database or benchmark was started; `ads-run-lock` was absent and no containers
were running.

## Commands

```bash
podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
  docker.io/library/golang:1.26-bookworm \
  go run ./docs/select-query-templates/reviews/20261005-join-high-3/audit.go
podman run --rm -v "<worktree>:/repo:ro" -w /repo/studies/01-charity-tree \
  docker.io/library/golang:1.26-bookworm sh -c 'go vet ./... && go test ./...'
# topology report re-render, twice, producing image, inputs copied read-only:
podman run --rm --network none -v "<copy-of-24-topology-runs>:/results:ro" -v "<out>:/out" \
  localhost/charitybench-select:select-join-v3b-topo-single-i1-20261005 \
  -cmd join-report -results /results -report-out /out/r1.md
tools/queue/queue audit
```

## Results outside the audit program

- `go vet` / `go test`: PASS (`charitytree/harness` 0.149 s,
  `charitytree/selectbench` 0.251 s).
- Topology report re-render: both outputs and the committed
  `reports/select-join-v3b-topology-20261005.md` hash to
  `f973dc849b0fc13358225626981e8d50101842ee101b79e7fee4b3da3fd11d6d`.
- Headline ratios from the generated domain report: PostgreSQL isolated
  `cte_joined_outer_mat` / `cte_joined_outer` at threshold 20 = 48.4540 / 0.4948
  = 97.9; YugabyteDB isolated = 86.6027 / 1.8999 = 45.6.
- Cluster tservers report `ysql_num_shards_per_tserver=2`
  (`resources/node.json`), while `yb_table_properties` shows 1 tablet per relation.
- Queue: claim `claim-f4eafb4a52e60123` at 21:55:18Z, submitted 00:27:55Z
  (commit 00:29:41Z) — real clock. Audit: 91 tasks, 0 errors, 2 pre-existing
  warnings on an unrelated completed task.
- Producing run tags inventoried before review: [TAGS-BEFORE.txt](TAGS-BEFORE.txt).

## Reuse ladder: NOT MATERIALIZED / MATERIALIZED per run (n = 1, 2, 4, 8, 16)

Per-run medians of five samples; computed from saved `join.json`.

| run | 1% consumer | 100% consumer |
|---|---|---|
| select-join-v3b-domain-pg-i1-20261005 | 0.17, 0.21, 0.36, 0.50, 0.71 | 1.05, 1.93, 2.74, 3.27, 3.96 |
| select-join-v3b-domain-pg-i2-20261005 | 0.18, 0.26, 0.42, 0.60, 0.81 | 0.96, 2.06, 2.49, 3.26, 3.89 |
| select-join-v3b-domain-pg-i3-20261005 | 0.19, 0.29, 0.36, 0.50, 0.68 | 0.99, 1.70, 2.51, 3.27, 3.94 |
| select-join-v3b-domain-pg-b1-20261005 | 0.19, 0.30, 0.42, 0.49, 0.69 | 0.99, 1.75, 2.66, 2.92, 4.60 |
| select-join-v3b-domain-pg-b2-20261005 | 0.22, 0.30, 0.31, 0.49, 0.75 | 0.87, 1.90, 2.77, 3.22, 4.11 |
| select-join-v3b-domain-pg-b3-20261005 | 0.27, 0.27, 0.33, 0.48, 0.75 | 0.89, 1.87, 2.49, 3.08, 4.50 |
| select-join-v3b-domain-yb-i1-20261005 | 0.81, 1.50, 2.93, 6.19, 11.29 | 0.97, 1.86, 3.69, 6.73, 8.99 |
| select-join-v3b-domain-yb-i2-20261005 | 0.87, 1.46, 2.78, 6.39, 11.32 | 1.07, 1.85, 3.57, 6.85, 11.81 |
| select-join-v3b-domain-yb-i3-20261005 | 0.82, 1.43, 2.79, 6.34, 7.45 | 0.98, 1.93, 3.51, 6.81, 12.14 |
| select-join-v3b-domain-yb-b1-20261005 | 0.92, 1.43, 2.45, 9.79, 10.05 | 0.94, 1.80, 3.49, 5.72, 10.96 |
| select-join-v3b-domain-yb-b2-20261005 | 0.88, 1.52, 2.69, 7.05, 9.18 | 0.90, 1.97, 3.20, 6.34, 11.13 |
| select-join-v3b-domain-yb-b3-20261005 | 0.94, 1.45, 2.54, 6.88, 9.06 | 1.03, 1.68, 3.40, 6.04, 9.49 |

## Audit program output

## 1. Every v3b run: validity, samples, controls, build and digest

| run | arm settings | valid/cells | 5 samples | controls fired | report digest = sha256(join.json) |
|---|---|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | yugabyte/cluster3/2 CPU 3g/blended/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-coloc-b2-20261005 | yugabyte/cluster3/2 CPU 3g/blended/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-coloc-b3-20261005 | yugabyte/cluster3/2 CPU 3g/blended/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-coloc-i1-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-coloc-i2-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-coloc-i3-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/colocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-domain-pg-b1-20261005 | postgres/single/2 CPU 3g/blended/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-pg-b2-20261005 | postgres/single/2 CPU 3g/blended/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-pg-b3-20261005 | postgres/single/2 CPU 3g/blended/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-pg-i1-20261005 | postgres/single/2 CPU 3g/isolated/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-pg-i2-20261005 | postgres/single/2 CPU 3g/isolated/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-pg-i3-20261005 | postgres/single/2 CPU 3g/isolated/engine-default/pg-btree-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-b1-20261005 | yugabyte/single/2 CPU 3g/blended/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-b2-20261005 | yugabyte/single/2 CPU 3g/blended/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-b3-20261005 | yugabyte/single/2 CPU 3g/blended/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-i1-20261005 | yugabyte/single/2 CPU 3g/isolated/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-i2-20261005 | yugabyte/single/2 CPU 3g/isolated/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-domain-yb-i3-20261005 | yugabyte/single/2 CPU 3g/isolated/engine-default/yb-range-rank | 91/91 | 91 | 2/2 | true |
| select-join-v3b-feas-yb-range-20261005 | yugabyte/single/2 CPU 3g/isolated/engine-default/yb-range-rank | 0/0 | 0 | 2/2 | true |
| select-join-v3b-topo-cluster-b1-20261005 | yugabyte/cluster3/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-cluster-b2-20261005 | yugabyte/cluster3/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-cluster-b3-20261005 | yugabyte/cluster3/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-cluster-i1-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-cluster-i2-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-cluster-i3-20261005 | yugabyte/cluster3/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-b1-20261005 | yugabyte/single/6 CPU 9g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-b2-20261005 | yugabyte/single/6 CPU 9g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-b3-20261005 | yugabyte/single/6 CPU 9g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-i1-20261005 | yugabyte/single/6 CPU 9g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-i2-20261005 | yugabyte/single/6 CPU 9g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-equal-i3-20261005 | yugabyte/single/6 CPU 9g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-b1-20261005 | yugabyte/single/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-b2-20261005 | yugabyte/single/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-b3-20261005 | yugabyte/single/2 CPU 3g/blended/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-i1-20261005 | yugabyte/single/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-i2-20261005 | yugabyte/single/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |
| select-join-v3b-topo-single-i3-20261005 | yugabyte/single/2 CPU 3g/isolated/noncolocated/yb-range-rank | 30/30 | 30 | 2/2 | true |

Client image ids across v3b runs: 1 distinct

## 2. Index use in domain plans (cells with plan_index_scan / plan_seq_scan)

| run | index-scan cells | seq-scan cells | j01 t=20 storage rows |
|---|---|---|---|
| select-join-v3b-domain-pg-b1-20261005 | 24 | 61 | 0 |
| select-join-v3b-domain-pg-b2-20261005 | 24 | 61 | 0 |
| select-join-v3b-domain-pg-b3-20261005 | 25 | 61 | 0 |
| select-join-v3b-domain-pg-i1-20261005 | 24 | 61 | 0 |
| select-join-v3b-domain-pg-i2-20261005 | 25 | 61 | 0 |
| select-join-v3b-domain-pg-i3-20261005 | 25 | 61 | 0 |
| select-join-v3b-domain-yb-b1-20261005 | 67 | 61 | 110 |
| select-join-v3b-domain-yb-b2-20261005 | 67 | 60 | 110 |
| select-join-v3b-domain-yb-b3-20261005 | 67 | 60 | 110 |
| select-join-v3b-domain-yb-i1-20261005 | 66 | 61 | 110 |
| select-join-v3b-domain-yb-i2-20261005 | 66 | 61 | 110 |
| select-join-v3b-domain-yb-i3-20261005 | 67 | 60 | 110 |

## 3. Contract pairs per domain run

| run | MAT/default joined t=20 | t=2000 | t=180000 | NOTMAT/MAT n=1 c=1% | n=1 c=100% | n=16 c=1% | n=16 c=100% |
|---|---|---|---|---|---|---|---|
| select-join-v3b-domain-pg-b1-20261005 | 99.09 | 23.12 | 1.441 | 0.186 | 0.988 | 0.687 | 4.598 |
| select-join-v3b-domain-pg-b2-20261005 | 83.58 | 22.82 | 1.106 | 0.217 | 0.868 | 0.752 | 4.108 |
| select-join-v3b-domain-pg-b3-20261005 | 101.15 | 25.03 | 1.147 | 0.270 | 0.893 | 0.746 | 4.500 |
| select-join-v3b-domain-pg-i1-20261005 | 99.08 | 24.09 | 1.024 | 0.171 | 1.053 | 0.712 | 3.958 |
| select-join-v3b-domain-pg-i2-20261005 | 100.69 | 23.97 | 1.220 | 0.178 | 0.960 | 0.806 | 3.891 |
| select-join-v3b-domain-pg-i3-20261005 | 97.47 | 26.67 | 1.044 | 0.188 | 0.988 | 0.683 | 3.938 |
| select-join-v3b-domain-yb-b1-20261005 | 38.63 | 9.05 | 0.990 | 0.916 | 0.941 | 10.049 | 10.962 |
| select-join-v3b-domain-yb-b2-20261005 | 45.69 | 9.50 | 0.977 | 0.880 | 0.899 | 9.178 | 11.132 |
| select-join-v3b-domain-yb-b3-20261005 | 34.51 | 8.64 | 0.928 | 0.944 | 1.028 | 9.058 | 9.486 |
| select-join-v3b-domain-yb-i1-20261005 | 43.37 | 8.86 | 1.025 | 0.812 | 0.967 | 11.294 | 8.990 |
| select-join-v3b-domain-yb-i2-20261005 | 48.64 | 8.65 | 0.959 | 0.874 | 1.066 | 11.323 | 11.806 |
| select-join-v3b-domain-yb-i3-20261005 | 42.19 | 8.07 | 0.934 | 0.824 | 0.979 | 7.451 | 12.139 |

## 4. Topology arms: join_where run medians (ms) and MAT/default at t=2000

| arm | t=20 | t=180000 | joined MAT/default t=2000 |
|---|---|---|---|
| topo-single-i | 2.02, 1.96, 2.06 | 105.4, 101.1, 138.5 | 8.7, 9.4, 9.4 |
| topo-equal-i | 1.81, 2.28, 1.85 | 109.0, 107.9, 101.6 | 8.9, 9.6, 8.2 |
| topo-cluster-i | 2.57, 3.22, 3.24 | 124.3, 168.3, 157.7 | 9.8, 14.0, 7.2 |
| coloc-i | 2.46, 1.74, 1.85 | 135.1, 134.4, 138.9 | 17.2, 16.1, 16.3 |
| topo-single-b | 1.99, 1.96, 2.02 | 117.1, 114.8, 114.1 | 9.0, 9.0, 9.7 |
| topo-equal-b | 2.07, 1.88, 1.83 | 119.1, 117.7, 112.1 | 8.8, 9.5, 9.0 |
| topo-cluster-b | 2.17, 3.15, 2.96 | 162.2, 127.0, 160.0 | 13.7, 13.0, 8.4 |
| coloc-b | 2.41, 2.40, 2.45 | 116.0, 149.1, 141.2 | 15.3, 17.5, 10.6 |

## 5. Cluster placement: tablets and leader node per relation

| run | placement | relation tablets (yb_table_properties) | leader nodes (yb-admin) |
|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 (distinct leader nodes: 1) |
| select-join-v3b-coloc-b2-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 (distinct leader nodes: 1) |
| select-join-v3b-coloc-b3-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-coloc-i1-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-coloc-i2-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-coloc-i3-20261005 | colocated | charity=1c person=1c donation=1c donation_pkey=0 donation_rank_idx=1c donation_person_idx=1c | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-b1-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-b2-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-b3-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-i1-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-i2-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n3 donation:yb-n3 donation_person_idx:yb-n3 donation_rank_idx:yb-n3 person:yb-n3 (distinct leader nodes: 1) |
| select-join-v3b-topo-cluster-i3-20261005 | noncolocated | charity=1 person=1 donation=1 donation_pkey=0 donation_rank_idx=1 donation_person_idx=1 | charity:yb-n2 donation:yb-n2 donation_person_idx:yb-n2 donation_rank_idx:yb-n2 person:yb-n2 (distinct leader nodes: 1) |

## 6. Blended writer and endpoints

| run | reader endpoint ops | writer offered/completed/failed | writer endpoint ops | overlap/samples |
|---|---|---|---|---|
| select-join-v3b-coloc-b1-20261005 | [168 168 168] | 128/128/0 | [43 43 42] | 61/210 |
| select-join-v3b-coloc-b2-20261005 | [168 168 168] | 140/140/0 | [47 47 46] | 66/210 |
| select-join-v3b-coloc-b3-20261005 | [168 168 168] | 133/133/0 | [45 44 44] | 61/210 |
| select-join-v3b-domain-pg-b1-20261005 | [921] | 335/335/0 | [335] | 93/637 |
| select-join-v3b-domain-pg-b2-20261005 | [921] | 340/340/0 | [340] | 131/637 |
| select-join-v3b-domain-pg-b3-20261005 | [921] | 349/349/0 | [349] | 114/637 |
| select-join-v3b-domain-yb-b1-20261005 | [933] | 526/526/0 | [526] | 261/637 |
| select-join-v3b-domain-yb-b2-20261005 | [934] | 523/523/0 | [523] | 234/638 |
| select-join-v3b-domain-yb-b3-20261005 | [934] | 527/527/0 | [527] | 219/638 |
| select-join-v3b-topo-cluster-b1-20261005 | [168 168 168] | 176/176/0 | [59 59 58] | 68/210 |
| select-join-v3b-topo-cluster-b2-20261005 | [168 168 168] | 175/175/0 | [59 58 58] | 75/210 |
| select-join-v3b-topo-cluster-b3-20261005 | [168 168 168] | 178/178/0 | [60 59 59] | 78/210 |
| select-join-v3b-topo-equal-b1-20261005 | [505] | 156/156/0 | [156] | 68/211 |
| select-join-v3b-topo-equal-b2-20261005 | [504] | 155/155/0 | [155] | 72/210 |
| select-join-v3b-topo-equal-b3-20261005 | [504] | 153/153/0 | [153] | 69/210 |
| select-join-v3b-topo-single-b1-20261005 | [504] | 156/156/0 | [156] | 73/210 |
| select-join-v3b-topo-single-b2-20261005 | [504] | 160/160/0 | [160] | 74/210 |
| select-join-v3b-topo-single-b3-20261005 | [504] | 160/160/0 | [160] | 65/210 |

## 7. Generated per-run reports: stale limit text

- stale limit phrases found: 0
- topology report table rows starting with a class id: 248
