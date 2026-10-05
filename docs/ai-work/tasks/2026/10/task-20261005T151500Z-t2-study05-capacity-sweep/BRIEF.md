# T2: sweep Study 05 cache capacity as a fraction of the working set

**Task:** `task-20261005T151500Z-t2-study05-capacity-sweep` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-05/t2-capacity-sweep` · **Worktree:** `.worktrees/t2-study05-capacity-sweep` · **Study:** `studies/05-cache-consistency`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

v2-14's cache-through gain was measured with a cache that holds the working set. The charter asks where the gain shrinks or reverses as the cache holds less of it.

## What to do

1. Use cache-capacity-kib / redis-maxmemory-mb with cache-fits to set at least five capacities as fractions of the measured working set (for example 5%, 10%, 25%, 50%, 100%, 150%), recording the eviction policy and observed hit rate.
2. Compare the no-cache baseline with cache-through (and cache-aside if the primary did) at each fraction, keeping the freshness gate and its negative control.
3. Report throughput and p99 ratios to the baseline per run with sign counts; state where the gain falls inside the baseline's spread, with an interval only where bracketed.
4. Register with a crossover record (levels at least 5).

## Acceptance criteria

- At least five capacity fractions, each in at least three whole runs, with hit rate recorded per cell
- No cell serves an impossible or wrong value beyond the study's floor; a breach aborts the cell
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
