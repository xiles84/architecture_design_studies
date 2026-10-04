# Execution Handoff Update — join/CTE correction

HIGH GPT-6 / Codex desktop; effort unknown; 2026-10-04.
Task `task-20261002T131729Z-select-join-cte`; branch `study-01/select-join-cte`;
retained worktree `.worktrees/select-join-cte`. Reviewed submission `937388a0`,
result `89479eb6`, producing source `7e148174`. This is a targeted amendment to
the original [HIGH handoff](../../JOIN-HANDOFF-20261004.md), not a restart of the
SELECT program. Acceptance is not narrowed to get the submission approved.

## Resume and protect existing work

```bash
tools/queue/queue session-start --capability LOW --role executor --model <actual-model> --tool <actual-tool> --effort <known-or-unknown>
tools/queue/queue claim --task task-20261002T131729Z-select-join-cte --capability LOW --role executor --model <actual-model> --tool <actual-tool> --effort <known-or-unknown>
```

Run from the retained task worktree. Read the generated claim output for the new
attempt, then `queue guard` before mutations and heartbeat at least every 15
minutes. Preserve both old attempts. Do not change any original result, signed
analysis or producing tag. Commit fixes before every new reported run; use
new run IDs and annotated producing tags. Do not use the existing six original
run names. Integrate current local main into this branch before correction if
it advanced, preserving other sessions' documents.

## Ordered corrections and decisions

1. **Fix reporting and validate saved evidence first.** Correct the extra-newline
   input digest, separate client and engine provenance, and isolate replicated
   reports by environment, dataset, index/settings and phase. Add failing-before,
   passing-after tests for these defects, regenerate into new review outputs,
   and record a validation report. Old original reports/results remain retained;
   if publishing a replacement report, move the superseded generated report to
   `outdated/` with an explicit replacement reference. Add the missing live
   semantic checks, engine-specific EXPLAIN options, volatility/policy evidence
   and explicit query-throughput/latency summaries. Derived queries/second from
   sequential latency must be labelled as derived; it is not a capacity test.

2. **Implement the Study 01 domain arm.** Keep the reusable join catalogue,
   shuffled-rank fixture, independent bag/optional/checksum oracles, separate
   contracts and seeded interleaving. Bind the measured parent/child fixture
   explicitly to person/donation with its charity relationship; document the
   projection, generated keys and rank proxy. Retain the original customers/facts
   evidence as a synthetic microfixture. Pure alias/table-name differences are
   not a new scientific factor, but the new domain results must identify their
   actual source and input bytes. Do not migrate Study 01's historical harness
   or touch its published design results. Run the existing approved single-node
   60-cell spelling + 30-cell reuse matrix (3 whole runs/engine, 5 trials/cell)
   for the corrected arm; unsupported policies stay visible. No extra axis or
   larger fixture is requested.

3. **Use effective WSL resources, not machine configuration metadata.** Recheck
   `/proc/meminfo`, online CPUs and effective cgroup limits under the run lock.
   The reviewed runtime was 8 CPUs / 16,109,780 kB. The standard 3×2 CPU/3 GiB
   database budget plus the original 2-CPU client leaves no system CPU reserve.
   The preselected feasibility budget is therefore database **6 CPU/9 GiB total**,
   client **1 CPU/2 GiB**, system **1 CPU/2 GiB**: **8 CPU/13 GiB total**, within
   the saved runtime. These are ceilings/reserves, not proof of attainable load.
   Calibrate the one-CPU client, including writer demand and throttling, before
   freezing it. Do not change the VM or host allocation. If live resources differ
   or the client demonstrably cannot supply demand, record the calibration and
   escalate the affected resource decision narrowly; other fixes continue.

4. **Complete the approved bounded topology subset.** Use a matching 1-CPU client
   for new single/cluster comparisons. Keep standard single-node 2 CPU/3 GiB and
   3-node 2 CPU/3 GiB per node as the per-node comparison. The separately labelled
   equal-total single receives **6 CPU/9 GiB**, compared to the same three-node
   aggregate. Include master/coordinator memory inside these totals. Keep these
   new calibrated subsets separate from original 2-CPU-client observations.
   The original handoff subset is unchanged: bag selectivity {0.01%,1%,90%},
   JOIN+WHERE and joined-output CTE under three policies; reuse {1,4,16}, three
   policies, both consumer selections. Three whole runs/arm, five trials,
   seeded interleaving and actual-first-cell repeat. Implement round-robin
   endpoint selection for the sequential reader, record per-endpoint operations
   and connections, and size the pool for readers, writer, endpoint setup and
   controls explicitly. Verify RF3, node limits, per-node throttling and tablet
   placement. Compare colocated/non-colocated schemas on this subset with matched
   logical data, resources and replication. Actual unsupported placement is an
   allowed gap after a feature check with error/evidence; no invented data-placement
   success and no inference from same-host containers. Do not silently combine
   budgets, endpoint modes or placement arms.

5. **Add the matched blended workload without changing projected answers.**
   Use one fixed writer that increments an independent maintenance counter on
   selected existing donation rows. The counter is excluded from bag projection,
   rank, keys and the MD5 id token. Writer updates never alter measured answer
   values. This preserves the independent read oracles while exercising row/index
   maintenance. Use one writer at a recorded offered rate of **10 committed updates/s**,
   rotating deterministically through a fixed set; record offered/completed rate,
   retries/failures, latency, per-key counts and final Go-ledger validation.
   Isolated controls load the same counter column and data, with writer disabled.
   Hold writer, transactions and demand fixed across spellings. Run separate
   labelled isolated and blended phases for the domain matrix and bounded
   topology subsets, with fresh fixture/reset between phases and correctness
   controls per phase. Start writer before phase warmup; record per-sample
   overlap and mutation counts so a short cell with no write overlap cannot be
   claimed as a concurrent measurement. The writer is not an optimistic versus
   pessimistic comparison. A demonstrated inability to sustain this demand is
   an explicit client-limited outcome or narrow escalation, never a silent change
   to rate or answer contract.

6. **Publish corrected interpretation and coverage.** Use new signed files;
   the old author may publish a successor with `supersedes`, but cannot silently
   edit the original. Read [CHECKS.md](CHECKS.md) for independent matched ratios.
   Give P1–P4/P12 and P5 control separate verdicts with original definitions.
   Do not claim a crossover merely because the consumer endpoints differ or
   because the near-equal n=1 sign is repeated; unresolved is accepted. No new
   consumer sweep is required. Keep small effects unresolved where sample or
   run variability dominates. Record absent hardware/network coverage plainly.

## Validation, submission and route

Record every check/calibration in generated or task-owned reports. All database
starts, feature checks and matrices acquire `ads-run-lock`; wait for another
holder and never remove an active lock. Mark long runs in CONTEXT.md. Preserve
all original input digests and tag objects in TAGS-BEFORE.txt. Every new result
names source, client/engine image identities, phase, budgets, endpoints, physical
placement, dataset, controls and five samples. Save actual EXPLAIN ANALYZE
VERBOSE/BUFFERS (PG) and DIST (YB); unsupported flags require explicit feature
evidence. New reports must re-render without digest changes and must not pool
incompatible resource, placement, environment or phase settings.

Publish the new attempt's PROGRESS.md, RESULT.md, exact digests and comparison
coverage. Submit the same task for HIGH review with a clean tree. HIGH owns
approval and integration; no completion tag or book publication before it.
Leave synthesis and the HIGH book task dependency-gated; never execute the
superseded LOW book task. No push/pull. Routine implementation details follow
Decide -> Log -> Continue; escalate only material scientific/resource failures.

NEXT MODEL: LOW
