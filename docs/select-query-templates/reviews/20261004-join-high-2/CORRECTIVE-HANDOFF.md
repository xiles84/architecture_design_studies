# Execution Handoff Update 2 — join/CTE correction

HIGH Claude Opus 5.5 / Claude Code desktop; effort unknown; 2026-10-04.
Task `task-20261002T131729Z-select-join-cte`, branch `study-01/select-join-cte`,
worktree `.worktrees/select-join-cte`. Reviewed submission `ad6deb42` /
result `cd84859f`. Read [REVIEW.md](REVIEW.md) and [CHECKS.md](CHECKS.md).
This amends the [original handoff](../../JOIN-HANDOFF-20261004.md) and
[update 1](../20261004-join-high/CORRECTIVE-HANDOFF.md). Everything they
require and this file does not change still applies. No acceptance is narrowed.

## Resume

Claim the same task as LOW from the retained worktree (commands in update 1),
create a new attempt, `queue guard` before mutations and heartbeat at least
every 15 minutes. **Never pass `--now` to a live queue command.** If the lease
lapses, recover with the real clock and wait out the grace period. Preserve
every existing result, report, tag, attempt file and signed analysis. New
measurements use new run IDs (suggested prefix `select-join-v3-…-20261005`)
and annotated producing tags. Integrate current local `main` first if it
advanced.

## Decisions made by HIGH

1. **The matched YugabyteDB index arm is range-sharded.** On YB, create the
   rank index as `donation (rank ASC)`. Change nothing else: PostgreSQL btree
   DDL, the person index and primary keys stay as they are, because those
   serve equality only. Label the arm as a setting, for example
   `index_arm=yb-range-rank`. The existing hash-index YB runs stay valid as the
   separately labelled physical arm `yb-default-hash-rank`. The hash/range pair
   is a one-decision physical comparison, not a syntax effect. The synthetic
   microfixture is not re-run; label its YB hash index in the analysis.
2. **Feasibility gate for the YB matched arm.** Record `pg_indexes.indexdef`
   for every index. Before any matrix, the default `j01_join_where` at
   thresholds 20, 200 and 2,000 must use `donation_rank_idx`.
   - If the range index still is not chosen, record plans, estimates and
     statistics, then escalate narrowly. Do not add hints or change the
     threshold axis.
   - Add YB `Storage Rows Scanned` (and PostgreSQL rows/buffers where already
     parsed) per cell, so the index regime is visible in the report.
3. **One frozen client build for every new run.** Commit all code and report
   fixes first, then run from that commit.
4. **Domain matrix, both phases, both engines** (91 cells, 5 trials, seeded
   interleaving, actual-first repeat, 3 whole runs each):
   - PostgreSQL single 2 CPU/3 GiB: isolated ×3 and blended ×3.
   - YugabyteDB single 2 CPU/3 GiB, range arm: isolated ×3 and blended ×3.
   - The earlier `domain2` runs remain a cross-build replication observation.
5. **Topology subset on the range arm:** single, equal-total single and
   cluster3 non-colocated, each isolated ×3 and blended ×3 (18 runs).
   - Budgets and the subset matrix are unchanged.
   - The blended writer rotates across endpoints in cluster arms and records
     per-endpoint writer operations. This removes the single-endpoint gap; the
     offered rate stays 10 updates/s.
6. **Colocated arm, actually attempted.** On the RF3 cluster, create a
   database with `COLOCATION = true`. Load the identical domain data and range
   index arm, with the same resources and endpoints, and point the client DSN
   at it. Run the cluster3 subset isolated ×3 and blended ×3.
   - Placement evidence for both cluster arms:
     - per-table colocation (`yb_table_properties`);
     - tablet count per table and index;
     - leader node per tablet, using `yb-admin`/`yb_local_tablets`, whichever
       the pinned image supports (LOW decides; log it).
   - Only an error from a genuinely colocated database is an allowed gap.
     Same-host containers are never colocation evidence.
7. **Reporting fixes** (failing-before/passing-after tests):
   - The compatibility class includes every arm factor: node and client
     budget, topology, endpoint mode, placement, phase and index arm.
     A test asserts the existing single and equal-total topology runs land in
     different classes.
   - Recorded limits are derived from the run's own settings; a blended run's
     report must not say "no concurrent writers". Fix the stale
     `distributed_rpc` text.
   - Render per-node and client `cpu.stat` deltas (`nr_periods`,
     `nr_throttled`, `throttled_usec`) in per-run and replicated reports. The
     replicated report also lists per-run endpoint operations, writer
     offered/completed/failed, retries and overlap.
   - State that sample latency includes retry and backoff time, and show
     per-cell retries.
   - Re-render the existing topology report into a corrected file; move the
     pooled one to `reports/outdated/` with a replacement reference. Prove
     digest stability by rendering twice.
8. **Decision Log.** The new attempt's own `PROGRESS.md` logs, as
   retroactive entries, the unlogged decisions of `c30935`: the retry policy,
   view binding, domain2 rerun, writer endpoint and the placement conclusion.
   It also logs every new decision. Do not edit earlier attempt files.
9. **Signed successor (v3)** superseding v2. Move v1 and v2 to
   `reports/outdated/`; they are your own files.
   - Cite only generated-report figures.
   - Give P1–P4 and P12 separate verdicts, with P5 as a control only.
   - Report YB per index arm, and the hash/range pair as a physical factor.
   - Correct the rollup/rolldown/embedding justification: this
     query-spelling experiment holds Study 01's physical arm fixed, so it is
     non-applicable here, not absent from the data model.
   - Keep YB reuse n=1 mixed-sign/unresolved.
   - Compare colocated and non-colocated, and the topology arms per budget.
   - Name the remaining weaknesses.

## Acceptance for resubmission

- Range-arm indexdef and selective plan use are verified, or a narrow
  escalation is filed with evidence.
- All runs above complete with correctness gates and controls firing; any
  failed cell is reported with its reason.
- Both cluster arms have placement evidence. The colocated arm is attempted on
  a colocated database.
- Reports separate every arm factor, derive limits truthfully, render
  throttling, endpoint, writer and retry evidence, and re-render with
  identical digests.
- The new attempt has its own RESULT.md and Decision Log. v3 is published,
  with v1/v2 moved.
- `go vet ./...` and `go test ./...` pass in Podman; queue audit shows 0
  errors; the tree is clean; then submit.

## Escalate only if

The range index cannot be made usable on the pinned YB; the colocated database
cannot hold the fixture within budget; or the client cannot sustain the writer
demand on a cluster arm. Everything else: Decide → Log → Continue.

NEXT MODEL: LOW
