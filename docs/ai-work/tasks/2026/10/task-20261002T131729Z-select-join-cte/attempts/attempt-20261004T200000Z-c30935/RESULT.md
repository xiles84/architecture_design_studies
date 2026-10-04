# Result — Study 01 SELECT join/filter spellings and CTE materialization (corrective cycle)

**Task:** `task-20261002T131729Z-select-join-cte`
**Capability/role:** LOW executor (model `deepseek-flash`, opencode CLI)
**Branch:** `study-01/select-join-cte` (worktree `.worktrees/select-join-cte`)
**Measured:** yes — 30 select-join runs, all cells verified against the independent Go oracle.

## Delivered (the six corrective items)

1. **Reporting/evidence fixes.** Compatibility classes key on environment,
   client/engine image digests, dataset fingerprint and settings (not image
   names that carry the run id); engine-specific `EXPLAIN` options; a live
   semantic gate; latency and derived-throughput sections; a per-run latency and
   endpoint table. Provenance separates the database engine image from the
   benchmark client image. Tests: `harness/select_join_report_test.go`.
2. **Study 01 domain arm** (`selectbench/joindomain.go`): real
   `charity`/`person`/`donation` tables with the catalogue bound through views,
   so the measured SQL is byte-identical to the synthetic arm. Frozen 3+3 matrix
   `select-join-domain2-{pg,yb}{1,2,3}-20261004`: 91/91 valid each, both
   semantic controls fire, one client build per engine.
3. **Resources.** Attested 8 CPU / 16,109,780 kB effective; D8 resource
   rejection retracted; 1-CPU client calibrated (`RESOURCES.md`).
4. **Bounded topology subset.** Three arms — single (2 CPU/3 GiB), equal-total
   single (6 CPU/9 GiB), 3-node RF3 (2 CPU/3 GiB/node) — isolated and blended,
   3 whole runs each, 5 trials: 18 runs, all 30/30 valid; round-robin endpoints
   with per-node throttling recorded. Placement feature check: `noncolocated`
   supported, `colocated` rejected by the pinned engine (SQLSTATE 0A000) and
   recorded as an explicit gap.
5. **Matched blended writer.** Independent `maintenance_updates` counter outside
   the projection/rank/token; one writer at an offered 10 updates/s over a fixed
   1000-key set; isolated and blended phases; offered/completed/failed rates,
   latency, per-key spread, a Go ledger and reader-write overlap. Sustained
   10.02/s with 0 ledger mismatches; transient SQLSTATE 40001 conflicts on the
   longest unmaterialized reads are retried up to 10 times with backoff and
   counted, with the oracle re-verified after every retry (`BLENDED.md`).
6. **Signed successor analysis**
   `reports/analyses/20261004T-select-join-cte-v2--deepseek-flash--2026-10-04.md`,
   superseding my original. It reports the contract's actual reuse pair
   (`NOT MATERIALIZED / MATERIALIZED`), leaves the one-reuse endpoint unresolved,
   keeps mixed signs, and states the PostgreSQL-only crossing. Published reports:
   `reports/select-join-domain2-*.md`, `reports/select-join-topology-subset-20261004.md`.

## Acceptance criteria check

| Correction | Result |
|---|---|
| R4 digest/provenance/partitioning; R3 engine EXPLAIN; live semantic gate | met |
| Study 01 domain arm, retained synthetic microfixture, frozen 3+3 matrix | met (91/91 each) |
| Effective resources; retract resource rejection; 1-CPU client | met |
| Bounded single/cluster/equal-total subset with round-robin and per-node throttling | met (18 runs, 30/30 each) |
| Colocated vs non-colocated | partial — non-colocated measured, colocated an explicit engine gap with error evidence |
| Blended read/write phase, isolated pair, ledger | met (10.02/s, 0 ledger mismatch) |
| Successor signed analysis with corrected P1–P5/P12 | met |

## Notes and gaps

- No book claim is promoted here; HIGH review owns promotion.
- Open gaps: a colocated database arm (requires a database created with
  `COLOCATION = true`), cluster-balanced writes (the writer uses one endpoint),
  and more than three runs per arm.
- The per-run analyses index does not bind the successor analysis because its
  `inputs_digest` is a composite of six run digests; the individual published
  reports carry their own digests.
