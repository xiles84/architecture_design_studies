# Execution Handoff — join/filter and CTE case

Owner update: [request](../book-join-theory-20261004/REQUEST.md), 2026-10-04.
HIGH planner: GPT-6 / Codex desktop / effort unknown. This dated amendment
controls join/CTE execution where the original protocol's bounded table conflicts.
It preserves the original protocol and the previous LOW worker's authored plan.

## Goal and current state

Extend **Study 01 charity/person/donation**; customers/orders are only the theory
fixture. Do not create a second study or migrate the historical harness.
Theory already teaches the six spellings and CTE/LEFT JOIN boundaries.
The shared SELECT harness and three other measurement families are integrated.
The join/CTE task is ready with a yielded planning checkpoint, not measured:
`task-20261002T131729Z-select-join-cte`, branch `study-01/select-join-cte`,
worktree `.worktrees/select-join-cte`. Its previous PROGRESS.md is preserved.

Claim through the queue as LOW; do not copy another worker's attempt.
Integrate current local main into the task branch before implementation so this
amendment and the randomized-order correction reach the resumed worktree.
Read the protocol, prior Decision Log, current harness and this handoff.

## Decisions resolved by HIGH before execution

1. Six **default** spellings are one bag contract: JOIN+WHERE, JOIN+ON filter,
   filtered derived table, filtered-input CTE, joined-output CTE+outer filter,
   comma join. Hold projection, parameters, data, indexes, statistics and
   prepared-plan policy fixed. The sixth spelling must be added if absent.
2. MATERIALIZED/NOT MATERIALIZED apply only to the **two CTE spellings**;
   do not fabricate policy cells for direct joins or derived tables. This gives
   ten forms per selectivity level: six defaults plus four forced-policy forms.
   Selectivity levels remain **0.01%, 0.1%, 1%, 10%, 50%, 90%**.
   Filtered-input materialization is a separate placement pair from the CTE
   of all joined rows filtered outside; report which boundary blocks pushdown.
3. Reuse is a separate contract and matrix: default/MATERIALIZED/NOT MATERIALIZED
   at **1, 2, 4, 8, 16 references**, with **1% and 100% consumer selection**.
   HIGH approves the previous LOW proposal to expand beyond the original {1,8}.
   Five levels permit a bounded crossing; a missing crossing is a valid result.
4. Start at **200,000 fact rows** with deterministic shuffled ranks covering
   1..200,000 and fixed parent keys; threshold counts are 20, 200, 2,000,
   20,000, 100,000 and 180,000. A dedicated rank predicate can make exact
   selectivity reproducible; label its mapping to the teaching amount predicate.
   Keep the fact count and rank/index set fixed across valid syntax pairs.
   Calibration may reject this size, but must not silently alter these axes.
5. For reuse, use a side-effect-free deterministic token such as
   `md5(repeat(id::text, 32))`, verify support/volatility on each pinned engine,
   and consume it in each reference (e.g. sum of ASCII codes of its first
   character). Return `(consumer_id, checksum)` for every reference; compute
   the independent oracle in Go. An unused expression is not expensive-work
   evidence. Separate per-consumer filters stay outside the common CTE.
   Reject optimizer-elided work through actual plan/loop evidence.
6. One-use spellings and reuse contracts get separate results and predictions.
   Judge P1-P4 and P12; P5 is a semantic control, never a speed competitor.
   Do not require a visible normalized-plan difference to retain a reproducible
   timing difference: investigate planning/client/counters and leave causal
   attribution unresolved when the available evidence cannot explain it.
7. PostgreSQL 17 single-node and pinned YugabyteDB single-node are required.
   Include a preselected 3-node YB RF3 subset: selectivity {0.01%,1%,90%},
   forms JOIN+WHERE and joined-output CTE under three policies; reuse {1,4,16}
   under the three policies at both consumer selections. Per-node baseline and
   equal-total control are separate. Cluster-only subsets cannot establish a
   five-level crossing; single-node curves provide that bounded claim.
8. Verify colocated and non-colocated placement on that subset with matched
   data/replication/resources. If the pinned image or current host cannot express
   or fit an arm, preserve an explicit gap with the failing feature/resource
   evidence; never equate same-host containers with data colocation.

## Ordered LOW work

1. Record live VM/host resources, system/client reserves, database totals,
   per-node budgets, reader/writer/control connection allowance and endpoints.
   Calibrate generator demand, client and per-node throttling; freeze settings.
   All database starts and checks hold `ads-run-lock`; wait for another holder.
2. Add missing embedded SQL, large fixture and independent Go oracles in the
   additive SELECT namespace. Keep the existing fixed writer and snapshot/ledger
   contracts. Compare stored rollup/rolldown/embedding only as separately labelled
   physical arms; this query-spelling experiment changes none of them.
   Optimistic/pessimistic strategy variation is non-applicable here.
3. Run a bounded feasibility check per engine: syntax support, actual CTE
   folding/materialization, selective counts, effective index use, consumed
   expensive work, full result draining and client capacity. Record each check
   in a generated report. Unsupported CTE policies are explicit cells, not defaults.
4. Prove the bag oracle on duplicate matches, NULL amounts, empty input and
   unmatched parents. For the optional-match contract, ON filtering and filtered
   child input must agree; a null-rejecting child WHERE substitution must fail.
   Keep identity-based membership separate from full joined-row timing.
5. Commit code before runs; use producing run tags and a distinct image/output
   namespace. Use **three whole runs per engine/arm, five trials per cell**,
   distinct seeded cell permutations, interleaved samples, and the **actual first
   cell repeated at the end**. Save seeds, actual sequence and drift. Read the
   order-correction signed HIGH review before interpreting variability.
6. Drain all rows; record read throughput, end-to-end latency distributions,
   sample counts, result rows/bytes, EXPLAIN ANALYZE, normalized signatures,
   estimates/actuals/loops, buffers, memory/temp spills and supported RPC counters.
   Keep isolated reads and a matched blended read/write phase separate. A cell
   timeout is **120 seconds**; save it as timeout, never substitute a faster query.
7. Generate measurements-only reports; submit RESULT.md, Decision Log, coverage
   matrix and exact inputs digests for HIGH review. A LOW signed interpretation
   is optional and separately attributed; it does not authorize book promotion.
   Local-main measurement integration still requires HIGH review.

## Acceptance and escalation

Every supported six-format/policy cell passes the independent oracle before
timing. Plans prove the claimed policy. Reports retain equal/no-resolved-difference,
adverse, invalid, timeout and unsupported outcomes; no universal syntax winner.
Use the agreed fixed schema/index arm for syntax attribution. Any different
index, placement or budget is a separately labelled factor.

If feasibility rejects the scientific fixture, expensive-work expression,
contract or required matrix, record the facts and escalate narrowly before
running affected cells. Routine code choices are Decide -> Log -> Continue.
Feature/resource failures can be logged as explicit gaps; they cannot support
a conclusion for the missing arm. A lock wait alone needs no escalation.

## Route after LOW submission

Existing HIGH `task-20261002T131729Z-select-synthesis` reviews the Decision Log,
corrected evidence and join results, signs analysis and registers eligible claims.
New HIGH `task-20261004T142623Z-select-book-high` then adds measured variations to
the existing Charity donations case and updates the selection guide/theory links.
It must retain uncertainty and unsupported regimes, rebuild and visually inspect
the book, and merge/tag locally. The original LOW `select-book` assignment is
superseded by this owner update; do not execute it.

NEXT MODEL: LOW
