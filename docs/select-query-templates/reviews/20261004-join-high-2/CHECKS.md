# Checks — second HIGH join/CTE review

Claude Opus 5.5 / Claude Code desktop / HIGH; 2026-10-04. Read-only: no
database or benchmark was started; `ads-run-lock` was absent and no containers
were running when the review began.

## Commands

```bash
podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
  docker.io/library/golang:1.26-bookworm \
  go run ./docs/select-query-templates/reviews/20261004-join-high-2/audit.go
podman run --rm -v "<worktree>:/repo:ro" -w /repo/studies/01-charity-tree \
  docker.io/library/golang:1.26-bookworm sh -c 'go vet ./... && go test ./...'
tools/queue/queue audit
```

## Results outside the audit program

- `go vet` / `go test`: PASS (`charitytree` no test files; `charitytree/harness`
  0.151 s; `charitytree/selectbench` 0.268 s).
- Queue audit: 91 tasks, 0 errors, 2 pre-existing warnings on the unrelated
  completed `task-20261001T133759Z-t2-study02-replication`.
- `domain2` report digests equal `sha256(join.json)` for all six runs
  (pg1 `347d12b6…`, pg2 `d4a7ba24…`, pg3 `ccd61591…`, yb1 `f03b1ad7…`,
  yb2 `0b8ec8dc…`, yb3 `e497ebb2…`).
- Equal-total cgroups: `cpu.max 600000 100000`, `memory.max 9663676416`;
  per-node single/cluster: `200000 100000`, `3221225472`.
- Integrated order-corrected families, first YB/PG run each (index scans /
  seq scans in saved plans): access YB 0/51, PG 29/26; aggregate YB 20/37,
  PG 46/10; existence YB 0/44, PG 4/40. The shared teaching DDL declares
  `PRIMARY KEY` and `CREATE INDEX … (amount)` without `ASC`.
- Queue time: event `recovered` and event `submitted` both record
  `occurred_at 2026-10-04T20:00:00Z`; their archive commits `6dc48392` and
  `ad6deb42` were written at 19:52:12Z and 21:23:41Z. Claim
  `claim-03114d688bd5fe53` expired 19:36:30Z; commit `ce28a962` is 19:47:03Z.
- Producing run tags inventoried before review: [TAGS-BEFORE.txt](TAGS-BEFORE.txt).

## Audit program output

## 1. Rank index use in saved EXPLAIN ANALYZE plans (domain2 matrix)

| run | valid cells | plans naming donation_rank_idx | plans with Seq Scan on donation + Storage Filter on rank |
|---|---|---|---|
| select-join-domain2-pg1-20261004 | 91 | 70 | 0 |
| select-join-domain2-pg2-20261004 | 91 | 70 | 0 |
| select-join-domain2-pg3-20261004 | 91 | 70 | 0 |
| select-join-domain2-yb1-20261004 | 91 | 0 | 73 |
| select-join-domain2-yb2-20261004 | 91 | 0 | 72 |
| select-join-domain2-yb3-20261004 | 91 | 0 | 73 |

## 2. Topology subset arms (manifest) and join_where medians (ms)

| arm | topology | node budget | phase | t=20 run medians | t=2000 run medians | t=180000 run medians |
|---|---|---|---|---|---|---|
| single-i | single | 2 CPU 3g | isolated | 56.7, 58.9, 61.9 | 57.8, 56.8, 64.9 | 119.9, 118.0, 123.4 |
| equal-i | single | 6 CPU 9g | isolated | 55.4, 54.7, 57.6 | 56.5, 58.3, 58.8 | 119.3, 116.6, 121.7 |
| cluster-i | cluster3 | 2 CPU 3g | isolated | 58.2, 61.3, 59.7 | 61.6, 58.9, 71.2 | 156.4, 159.1, 161.7 |
| single-bo | single | 2 CPU 3g | blended | 62.6, 57.3, 59.6 | 59.6, 64.1, 59.1 | 119.6, 140.9, 128.6 |
| equal-bo | single | 6 CPU 9g | blended | 55.8, 50.3, 60.5 | 64.4, 50.2, 64.0 | 130.7, 129.4, 130.7 |
| cluster-bo | cluster3 | 2 CPU 3g | blended | 64.2, 62.1, 59.8 | 63.2, 66.9, 60.7 | 171.8, 146.5, 181.2 |

Generated report classes containing runs of more than one node budget:

- class `76aa820fc736` pools node budgets 2 CPU 3g + 6 CPU 9g
- class `b5f93048b28f` pools node budgets 2 CPU 3g + 6 CPU 9g

## 3. Generated per-run limits contradicting the run

- blended runs whose report says `no concurrent writers ... blended phase is a recorded gap`: 9 of 9

## 4. Contract pairs per run (domain2)

| run | NOT MAT / MAT c=100% n=1 | n=2 | n=16 | c=1% n=1 | c=1% n=16 | joined-output MAT / default t=20 | t=2000 | t=180000 |
|---|---|---|---|---|---|---|---|---|
| select-join-domain2-pg1-20261004 | 0.904 | 1.729 | 4.427 | 0.204 | 0.828 | 66.91 | 26.13 | 1.18 |
| select-join-domain2-pg2-20261004 | 1.054 | 1.692 | 4.555 | 0.256 | 0.937 | 87.49 | 31.44 | 1.05 |
| select-join-domain2-pg3-20261004 | 0.887 | 1.606 | 4.389 | 0.185 | 0.836 | 81.44 | 29.14 | 1.17 |
| select-join-domain2-yb1-20261004 | 1.125 | 1.880 | 10.310 | 1.037 | 9.477 | 1.56 | 1.39 | 0.99 |
| select-join-domain2-yb2-20261004 | 0.988 | 2.051 | 9.298 | 1.018 | 9.285 | 1.43 | 1.44 | 0.97 |
| select-join-domain2-yb3-20261004 | 1.007 | 1.795 | 9.262 | 1.079 | 9.764 | 1.58 | 1.54 | 0.93 |

## 5. Placement feature-check runs

- select-join-topo-yb-colocated-20261004: placement=colocated cells=0 error="gate fixture load: unsupported by engine: ERROR: cannot set colocation true on a non-colocated database (SQLSTATE 0A000)"
- select-join-topo-yb-noncolocated-20261004: placement=noncolocated cells=30 error=""
