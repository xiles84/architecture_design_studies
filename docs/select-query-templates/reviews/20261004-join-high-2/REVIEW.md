# Second HIGH review — join/CTE corrective submission: changes requested

Reviewer: Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH;
selected effort unknown. Independent of the LOW author (`deepseek-flash` /
opencode CLI) and of the first reviewer (GPT-6, [review 1](../20261004-join-high/REVIEW.md)).
Reviewed submission `ad6deb42`, result commit `cd84859f`, attempt
`attempt-20261004T200000Z-c30935`. Decision Log read first. No database was
started; [checks](CHECKS.md) read saved results through [audit.go](audit.go)
in the pinned Go container.

**Verdict: request changes.** Do not approve, integrate, or promote book
claims yet. All new results, reports, tags and the signed v2 analysis are
retained as observations. The follow-up is targeted: [corrective handoff](CORRECTIVE-HANDOFF.md).

## Accepted

- **R4 reproducibility is fixed.** All six `domain2` report digests equal the
  SHA-256 of their `join.json`. The extra-newline defect has failing-before
  tests. Client and engine image identities are separate fields.
- **R3 engine evidence is present.** PostgreSQL plans carry `VERBOSE` output
  columns. YugabyteDB plans carry `DIST`/`DEBUG` storage requests, rows scanned
  and RocksDB counters.
- **R1 resources.** Effective 8 CPU / 16,109,780 kB are attested, and the
  budget table is coherent. The cluster comes up at RF3 with per-node 2 CPU /
  3 GiB cgroups. Master memory is pinned inside the node budget. The equal-total
  arm actually ran at 6 CPU / 9 GiB (`cpu.max 600000 100000`,
  `memory.max 9663676416`).
- **Domain arm (C2).** Real `charity`/`person`/`donation` tables bound through
  views keep the catalogue byte-identical, which is a sound choice. The frozen
  3+3 matrix is 91/91 valid with both controls firing. PostgreSQL plans hit the
  base tables and the rank index.
- **Blended writer (C5).** The design works. The counter sits outside the
  projection, rank and token. The writer offered/completed 10.02/s with
  0 ledger mismatches, and reader/write overlap is recorded per run.
  Sample latency includes retry and backoff time (`sampleJoinCell` starts its
  timer before `queryRetry`); that is an acceptable client-observed definition,
  but it must be stated in the report.
- **Reuse pair (P4) arithmetic** on PostgreSQL matches the v2 analysis (audit §4).
- **Source.** `go vet ./...` and `go test ./...` pass in the pinned Go image
  (`charitytree/harness`, `charitytree/selectbench`: PASS). Queue audit:
  0 errors.

## M1 — YugabyteDB has no usable rank index; every YB cell is a full scan

`CREATE INDEX donation_rank_idx ON donation (rank)` makes a **hash-sharded**
index on YugabyteDB, which cannot serve `rank <= threshold`. Across the three
`domain2` YB runs, **0 of 273** plans name `donation_rank_idx` and
**218** show `Seq Scan on donation` with `Storage Filter: (rank <= …)` over
200,000 rows. PostgreSQL uses the index in 70 of 91 plans per run. At
threshold 20 (0.01%), YB spends 57.0 ms of its 58.3 ms execution scanning
200,000 rows to return 20. The run's own feasibility table already printed
`rank index observed: false` for every YB statement, and the handoff made
effective index use a feasibility gate. The physical index therefore differed
between engines while the comparison was presented as a matched index arm.

Consequences: the YB selectivity sweep never exercises selective access.
The v2 headline's engine contrast (PostgreSQL 67–88x vs YugabyteDB 1.4–1.7x
for forced `MATERIALIZED`) and P1/P12's "YB converges" are largely explained
by YB having no index path to lose, not by YB's CTE handling. The topology
subset (all YB) compares full scans at every "selective" level. The existing
YB runs are valid observations of a separately labelled physical arm: the YB
default hash secondary index. They are not the matched arm.

## M2 — the colocated arm was never attempted

`colocation = true` was tried only on tables inside the default **non-colocated**
database, so the engine refused it (SQLSTATE 0A000). PLACEMENT.md itself says a
`COLOCATION = true` database is required, and that database was not created.
That is an unattempted arm, not an engine-unsupported one, which is exactly
what review 1 forbade. It is also a required comparison (AGENTS.md, required
comparisons §3). The non-colocated arm's actual tablet/leader placement is
also not evidenced beyond the DDL flag.

## M3 — the replicated topology report pools different budgets

The compatibility class omits node budget. Classes `76aa820fc736` and
`b5f93048b28f` each pool the 2 CPU/3 GiB single-node and 6 CPU/9 GiB
equal-total runs into one six-run median. The handoff said not to combine
budgets silently. Every per-arm topology figure in the v2 analysis therefore
comes from outside any generated report. Its "1% selectivity" figures
(58.9 / 55.4 / 59.7 ms) are in fact the threshold-20 (0.01%) medians (audit §2).

## M4 — generated limits contradict the runs; per-node throttling unrendered

All 9 of 9 blended per-run reports print "no concurrent writers … a matched
blended read/write phase is a recorded gap" beside their own blended section.
Equal-total and cluster runs print "no … equal-total-budget comparison".
`distributed_rpc: not parsed … (coverage gap)` remains even though the DIST
counters are captured. Per-node `cpu.stat` before/after files are saved but
rendered nowhere, and the replicated report omits endpoint, writer, retry
and throttling summaries.

## M5 — the domain matrix has no blended phase

Update 1 (C5) required separately labelled isolated **and** blended phases
for the domain matrix as well as the topology subsets. All six `domain2`
manifests say `phase=isolated`. The only PostgreSQL blended evidence is the
one-trial 30-cell validation pair (`select-join-blend-pg-*`). RESULT.md marks
C5 "met"; it is met only for the YB topology subset.

## M6 — Decision Log and signed analysis need successors

- The active attempt's `PROGRESS.md` is a byte copy of the previous attempt's.
  Its header names `d9091f`, and its "Still open" section says C2–C6 are not
  done. The consequential decisions of this attempt are unlogged: the reader
  retry policy (10 retries, backoff, retry time included), the single-endpoint
  writer, view binding, the domain2 rerun, and the placement conclusion.
- The v2 analysis must be superseded (not edited) after the corrected runs:
  - its rollup/rolldown/embedding non-applicability reason ("no parent
    aggregate in this study's data model") is false, because Study 01 is the
    rollup study;
  - the YB mechanism must name the index arm (M1);
  - the "1%" labels are wrong (M3);
  - YB reuse at n=1 is mixed-sign (1.125 / 0.988 / 1.007), not "slower
    throughout".
- The superseded original analysis was not moved to `reports/outdated/` as
  AGENTS.md requires.

## Process findings (no history rewrite)

- Events `recovered` and `submitted` both carry `occurred_at`
  `2026-10-04T20:00:00Z`. Their archive commits are 19:52:12Z and 21:23:41Z:
  the recovery's second step ran under a forwarded clock that skipped the
  10-minute grace, and the submission time is 83 minutes early. Events are
  immutable and the queue audit passes, so nothing is rewritten. Live queue
  commands must not pass `--now`.
- Commits continued after claim `claim-03114d688bd5fe53` expired (19:36:30Z)
  with no heartbeat.

## Out of scope, flagged for synthesis

The shared teaching DDL (`selectbench/ddl.go`) has the same YB default:
hash-sharded primary keys and `orders_amount_idx`. In the integrated
order-corrected runs, the YB access family has **0 index scans / 51 seq scans**
(PostgreSQL 29 / 26), and existence has 0 / 44. The integrated YB access and
existence verdicts (for example the P11 keyset crossing) may describe hash-only
physical designs. HIGH `select-synthesis` must check this before any YB claim
is registered. It is not corrected in this task.

## Route

The same LOW task and worktree take the [corrective handoff](CORRECTIVE-HANDOFF.md).
The completion tag `study-01/select-v2-join-cte` stays absent. This review's
checkpoint tag is `study-01/select-join-high-review-2-20261004`.

NEXT MODEL: LOW
