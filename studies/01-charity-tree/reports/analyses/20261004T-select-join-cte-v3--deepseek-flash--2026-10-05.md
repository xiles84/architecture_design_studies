---
analysis_id: 20261004T-select-join-cte-v3--deepseek-flash--2026-10-05
run_id: select-join-v3b-domain-{pg,yb}-{i,b}{1,2,3}-20261005; select-join-v3b-topo-{single,equal,cluster}-{i,b}{1,2,3}-20261005; select-join-v3b-coloc-{i,b}{1,2,3}-20261005
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek Flash via opencode CLI (LOW executor session)"
analyzed_at: 2026-10-05
inputs_digest: "see the per-run reports under reports/ (each carries its own join.json digest); this analysis aggregates the v3b arm"
repo_commit: a05d80e4
supersedes: 20261004T-select-join-cte-v2--deepseek-flash--2026-10-04
status: current
headline: "With a range-sharded rank index the engine contrast is PostgreSQL ~97x vs YugabyteDB ~42x for forced MATERIALIZED at 20 rows, not the 1.4-1.7x the hash-indexed v2 reported; colocation removes the cluster read penalty at a point read but makes the materializing plan dearer."
---

# Analysis (v3 successor) — join/filter spellings and CTE materialization, matched range-index arm — deepseek-flash

This supersedes my `…-v2--deepseek-flash--2026-10-04` (itself superseding v1; both
are in `reports/outdated/`). Review 2 found that the v2 YugabyteDB arm had **no
usable rank index**: `CREATE INDEX … (rank)` is hash-sharded, 0 of 273 YB plans
used `donation_rank_idx`, and every "selective" YB cell scanned all 200,000 rows.
The v2 engine contrast was therefore an index artefact. This analysis uses the
range-index arm (`rank ASC`, `index_arm=yb-range-rank`), whose plans show
`Index Cond` and 20 index rows scanned at threshold 20; the physical-arm gate
confirms index use at 20/200/2000 before every run. All figures below are
across-run medians from the generated reports
`reports/select-join-v3b-domain-replicated-20261005.md` and
`reports/select-join-v3b-topology-20261005.md`.

## TL;DR

- **The v2 engine contrast was wrong in magnitude and had the wrong mechanism.**
  With the matched range index, forcing `MATERIALIZED` on the joined-output CTE
  costs **~97x on PostgreSQL and ~42x on YugabyteDB** at 20 rows (23-24x and
  ~8.7-9.0x at 2,000). Both engines lose a selective path; YB loses about half
  as much, not one-fortieth.
- **P4 reuse crossing:** still PostgreSQL-only, and only at the full consumer.
  At 16 references the `NOT MATERIALIZED / MATERIALIZED` ratio is 3.96 (PG,
  isolated) and 11.73 (YB, isolated) at the full consumer, and 0.71 (PG) versus
  11.29 (YB) at the 1% consumer. YB at one reference is **mixed-sign/unresolved**
  (0.846 isolated, 0.929 blended at 1%; 1.011 / 0.922 at 100%).
- **Convergence (P12):** at 180,000 rows (90%) both engines converge to ~1x
  (PG forced materialization 1.07-1.15x, YB 0.96-0.98x).
- **Topology:** the extra budget of the equal-total single node buys nothing; a
  non-colocated cluster is 1.4-1.6x slower than single node; the **colocated**
  cluster matches single node at a point read (1.85 ms vs 2.02 ms at 20 rows)
  and cuts the 90% scan from 157.7 to 135.1 ms, but its materializing plan is
  dearer (16.2x vs 9.9x over its default at 2,000 rows).
- **Blended phase:** the writer held 10.0 updates/s with 0 failures, rotated
  evenly across endpoints ([59,59,58] on cluster), and reader/write overlap and
  cgroup throttling are recorded per run.

## What I was looking at

- **Domain matrix, 91 cells, 5 trials, position control, both phases, both
  engines**: PostgreSQL 17.11 (`pg-btree-rank`) and YugabyteDB 2025.2.6.0-b111
  (`yb-range-rank`), 3 whole runs each. 91/91 valid in every run; both semantic
  controls fire.
- **Topology subset, 30 cells, 5 trials, both phases**: single 2 CPU/3 GiB,
  equal-total single 6 CPU/9 GiB, non-colocated cluster3 RF3, and **colocated
  cluster3** (database `select_coloc`, `COLOCATION = true`), 3 whole runs each;
  30/30 valid everywhere, and the colocated runs record `is_colocated = true`
  with one tablet per table.
- The v3 (pre-gate-fix) runs and the v2 hash-arm runs are retained as labelled
  cross-build/physical-arm observations; they are not pooled here.

## Conclusions

**M1 corrected — forcing materialization of a join-containing CTE.** Against the
joined-output CTE's own default plan, PostgreSQL's forced `MATERIALIZED` is
96.1x (blended) and 97.9x (isolated) at 20 rows, 23.1x / 24.0x at 2,000, and
1.15x / 1.07x at 180,000. With the range index, **YugabyteDB is 38.6x (blended)
and 45.6x (isolated) at 20 rows**, 9.0x / 8.7x at 2,000, and 0.98x / 0.96x at
180,000. The direction of the v2 claim survives — PostgreSQL pays more — but the
magnitude does not: the earlier "1.4-1.7x on YB" described a hash-index arm that
had no selective path to lose, so the comparison was between an index scan and a
full scan. **Action:** any YB CTE-materialization statement must name the index
arm; a hash-sharded secondary index makes the engine look insensitive to a hint
that is in fact expensive. Every run now records `indexdef`, per-cell
`Storage Rows Scanned` and the gate result.

**P4 reuse pair.** The contract pair `NOT MATERIALIZED / MATERIALIZED` at 16
references is 3.958 (PG isolated) and 11.730 (YB isolated) at the full consumer,
and 0.712 (PG) versus 11.294 (YB) at the 1% consumer; blended values are 4.491
and 10.888, and 0.712 and 9.299. YB is therefore *not* uniformly slower: at one
reference it is 0.846 (isolated) and 0.929 (blended) at the 1% consumer, and
1.011 / 0.922 at the full consumer — **mixed sign, unresolved**. The crossing
(the point where not-materializing turns from a win into a loss) exists only on
PostgreSQL and only at the unselective consumer. **Action:** do not publish a YB
reuse crossing; the honest statement is "on PostgreSQL, NOT MATERIALIZED wins at
a selective consumer and loses once the consumer is unselective and the CTE is
reused more than once; on YugabyteDB the penalty grows monotonically with the
reuse count and the one-reference endpoint is inside run noise."

**P1 default spellings and P12 convergence.** With the range arm the default
joined-output CTE runs at 0.97-1.00 of join+WHERE on PostgreSQL and 0.85-0.93 on
YugabyteDB; all six default spellings are oracle-identical, and at 90% every
forced policy converges to ~1x on both engines. The positions are therefore
equivalent and their timings converge at high selectivity; the report still
shows sizable per-run spread on PostgreSQL at the 20-row threshold, so those
cells are reported as medians of three runs, not as single measurements.

**P5 controls.** Both deliberately-wrong controls (`x01`, `x02`) were detected in
every run. Control only.

**Topology, per budget and placement.** `join_where` across-run medians:
single 2 CPU/3 GiB 2.02 ms at 20 rows / 105.4 ms at 180,000; equal-total
6 CPU/9 GiB 1.85 / 107.9 — **the extra budget changes nothing** on this
workload. Non-colocated cluster3 is 3.22 ms (1.6x single) and 157.7 ms (1.5x);
colocated cluster3 is 1.85 ms (equal to single) and 135.1 ms (0.86x the
non-colocated cluster). Colocation keeps the point read local and reduces the
broad scan, but the forced-materializing plan costs 16.2x its default at 2,000
rows on the colocated cluster versus 9.9x on the non-colocated cluster: one
tablet per node concentrates the materialising work. Blended runs change little
(the same comparisons move by under ~10%), and the writer is balanced across the
cluster's three endpoints.

**Blended phase and throttling.** Every blended arm sustained the offered
10.0 updates/s with zero failed updates and reconciled its ledger; e.g. the
cluster3 blended run delivered 176 completed updates split [59,59,58] across
endpoints with 68 of 210 reader operations overlapping a commit. Per-node
cgroup throttling is recorded (e.g. yb-n3: 15 throttled periods, 650 ms over the
run) and rendered in the report. Reader retries of transient conflicts were 0 in
the v3b runs — with the range index the reads are short enough not to be
restarted by the writer.

## Where the measurement is weak

- Three whole runs per arm; the across-run median hides between-run variance,
  and the per-run spread is still large on PostgreSQL at the selective end.
- One host: cluster nodes are containers on one machine. Colocation is evidenced
  by `yb_table_properties`/`yb_local_tablets` and `yb-admin` leaders, **not** by
  independent hardware; the leader evidence shows one tablet per colocated table.
- The writer is a fixed maintenance load on one counter; it is not an
  optimistic/pessimistic concurrency comparison.
- The client is one CPU / 2 GiB; a cell that flattens may be a client ceiling.
- Rows-scanned metrics are parsed from the engine's own plan text; an unparsed
  cell reports zero and is not distinguished from a genuinely small scan.

## Coverage and gaps

- Measured pairs: six default spellings; forced `MATERIALIZED` /
  `NOT MATERIALIZED` against the default for the filtered and joined-output CTEs;
  the reuse pair over five reuse counts and two consumers; single vs equal-total
  single vs non-colocated cluster vs colocated cluster, isolated and blended.
- Non-applicable **here**: rollup/rolldown/embedding. This is a query-spelling
  experiment that holds Study 01's physical arm fixed; the study's data model
  does have parent-child structure (charity -> person -> donation) and those
  decisions are measured in Study 01's other arms, so they are out of scope for
  *this* comparison rather than absent from the model.
- Outstanding: more than three runs per arm; a colocated single-node arm;
  cluster write paths under contention beyond the fixed writer.
