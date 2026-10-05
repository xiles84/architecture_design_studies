# T2: sweep Study 04 replacement size across at least five levels (trigger vs application rollup)

**Task:** `task-20261005T151500Z-t2-study04-cardinality-sweep` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-04/t2-cardinality-sweep` · **Worktree:** `.worktrees/t2-study04-cardinality-sweep` · **Study:** `studies/04-configuration-portal`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

v2-18 shows the trigger-maintained rollup's penalty growing with replacement size (10, 60 and 240 entries), but three sizes cannot bound a crossover. The charter asks for sizes below and above the measured range.

## What to do

1. Use the existing entries-per-ip / cardinality-mode knobs; choose at least five sizes spanning 1 to at least 1,000 entries (for example 1, 10, 60, 240, 1,000), including the three already measured.
2. Hold the designs to the isolate pair n3_rollup_trigger against n4_rollup_app; keep the rollup oracle and its negative control.
3. Report n4/n3 per run and size with sign counts; state whether the ratio reverses or approaches 1 at the smallest size, and the crossing interval only if observations bracket one.
4. Extend or supersede v2-18 with a crossover record (levels at least 5).

## Acceptance criteria

- At least five replacement sizes, each in at least three whole runs, with no rollup mismatch or an aborted cell reported
- The registry states whether the size effect reverses, with its interval or reversal_observed false
- A protocol with predictions (bands, orderings, falsifiers) and the owner's required-comparison coverage table is committed BEFORE any run
- Sizing is calculated from the effective guest (nproc and /proc/meminfo in a container), reusing the primary's per-node and client budgets; oversubscription and per-container throttling are recorded
- At least three whole runs per target in seeded randomized or AB/BA order, the executed order read from the manifest, with a position control; every run tagged run/<study>/<run-id> from a clean tree
- Every cell passes its oracle gate and a negative control fires in each phase; failures abort the cell and are reported
- A signed analysis scores every prediction, reports per-run ratios with sign counts, and names where the measurement is weak
- Registry: new or re-scoped claims in book/evidence/v6 with ledger entry and carry-forward test; run_replicated only with a prediction committed before the runs; gaps the task closes are narrowed in the same change
- CONTEXT.md and LESSONS_LEARNED.md updated; merged into local main through the queue and tagged

## Validation

- tools/evidence validate: 0 errors; go test -count=1 ./... in tools/evidence passes
- tools/queue/queue audit: 0 errors
- the study's own tests pass in Podman

## Escalate only if

A primary configuration cannot start at all after effective-guest sizing, a promotion would change a claim's scientific question, or the placement instrument cannot list leaders on the pinned engine. Everything else: Decide → Log → Continue. Book prose is out of scope.

NEXT MODEL: LOW
