# Execution Handoff Update 3 — distributed placement arm and v4

HIGH Claude Opus 5.5 / Claude Code desktop; effort unknown; 2026-10-05.
Task `task-20261002T131729Z-select-join-cte`, branch `study-01/select-join-cte`,
worktree `.worktrees/select-join-cte`. Reviewed submission `99233466` /
result `a2f3fa9f`. Read [REVIEW.md](REVIEW.md) and [CHECKS.md](CHECKS.md).
This narrows the remaining work; everything accepted in round 2 stands.
**Do not re-run the domain matrix.**

## Resume

Claim the same task as LOW from the retained worktree. Use the real clock,
`guard` before mutations and heartbeat at least every 15 minutes. Preserve all
existing runs, reports, tags, attempts and signed files. Integrate current
local `main` first if it advanced. New run IDs, for example
`select-join-v4-topo-…-20261005`, each with an annotated producing tag.

## Decisions made by HIGH

1. **New placement arm `distributed`** (non-colocated, cluster only), created
   in the default database:
   - `charity`, `person` and `donation`: `SPLIT INTO 3 TABLETS`.
   - `donation_person_idx`: `SPLIT INTO 3 TABLETS`.
   - Range rank index:
     `CREATE INDEX donation_rank_idx ON donation (rank ASC) SPLIT AT VALUES ((66667), (133334))`.
   - Logical data, the range index arm, RF3, budgets and endpoints stay
     unchanged.
   - Relabel the existing engine-default arm `noncolocated-single-tablet`
     (setting value); it is not "distributed".
2. **Placement gate before measuring** a distributed run:
   - `yb_table_properties.num_tablets` = 3 for `donation` and
     `donation_rank_idx`;
   - the leaders of `donation`'s tablets are on **three distinct nodes**;
   - wait for the load balancer up to 5 minutes (LOW chooses the polling
     command; log it);
   - record the listing before and after the run.

   If leaders will not spread, record the listings and escalate narrowly.
   Do not move leaders by hand. Colocated and single-tablet arms record their
   placement but are not gated on spread.
3. **One new frozen build, full topology subset (30 runs):**
   - single 2 CPU/3 GiB, equal-total single 6 CPU/9 GiB, and the three
     cluster3 placements (`distributed`, `noncolocated-single-tablet`,
     `colocated`);
   - each arm isolated ×3 and blended ×3;
   - unchanged 30-cell subset, 5 trials, seeded interleaving and actual-first
     repeat.

   The `v3b` topology runs remain a labelled cross-build observation.
4. **Placement in the class key and the report.** The replicated topology
   report shows, per run:
   - tablets per relation;
   - leader nodes;
   - endpoint and writer operations;
   - throttling.
5. **Minor reporting fixes.**
   - Render PostgreSQL `storage_rows_scanned` as n/a.
   - Rename or document the whole-plan `plan_index_scan`/`plan_seq_scan`
     flags so they cannot be read as the rank access path.
   - Record the tservers' effective `ysql_num_shards_per_tserver` and the
     observed tablet count. State the cause if evidence shows it; otherwise
     mark it unexplained.
6. **v4 signed successor** superseding v3; move v3 to `reports/outdated/` and
   extend `REPLACEMENTS-select-join.md`.
   - P4 must follow the per-run signs in [CHECKS.md](CHECKS.md):
     - YB 1% consumer: a replicated sign change between n=1 (0.81–0.94) and
       n=2 (1.43–1.52); state the small size of the n=1 effect.
     - YB 100% at n=1: unresolved.
     - PG 1%: no crossing in 1–16.
     - PG 100%: crossing at or below n=2, with n=1 unresolved.
   - Topology and placement claims cite the new report and name the measured
     placement: leaders and tablets, not container locality. Compare
     colocated with distributed as the required pair, and single-tablet with
     distributed as the splitting effect.
   - Keep the accepted v3 domain conclusions (M1 corrected figures, P1, P5,
     P12) unless new evidence changes them.

## Acceptance for resubmission

- The distributed arm passes its placement gate, with saved before and after
  listings.
- 30 topology runs from one build, all valid with controls firing; failures
  reported.
- The topology report separates the five arms × two phases, renders
  placement, and re-renders byte-identically.
- The attempt has its own PROGRESS/Decision Log and RESULT. v4 is published
  and v3 moved.
- `go vet`/`go test` pass in Podman; queue audit shows 0 errors; the tree is
  clean.

## Escalate only if

The pinned engine rejects the split DDL, or leaders cannot be spread within
the gate's wait. Everything else: Decide → Log → Continue.

NEXT MODEL: LOW
