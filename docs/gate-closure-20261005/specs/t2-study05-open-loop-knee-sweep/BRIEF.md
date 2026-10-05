# T2: bracket Study 05's open-loop knee with at least five offered rates

**Task:** `task-20261005T151500Z-t2-study05-open-loop-knee-sweep` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-05/t2-open-loop-knee` · **Worktree:** `.worktrees/t2-study05-open-loop-knee-sweep` · **Study:** `studies/05-cache-consistency`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

The 2026-09-23 open-loop phase measured too few rates to tell the client's knee from the server's. The charter asks for rate points that separate them.

## What to do

1. Calibrate the client's own ceiling first (a no-op or cache-only target), then choose at least five openloop-rates that bracket the server knee and stay below the client ceiling.
2. At each rate record offered, completed and queued work, p50/p99, and client and database CFS throttling; a rate where the client throttles is labelled client-limited.
3. Report per-run completion ratio and latency against offered rate with sign counts; state the knee interval only where observations bracket it.
4. Register with a crossover record (levels at least 5).

## Acceptance criteria

- The client ceiling is measured and every reported rate is below it or labelled client-limited
- At least five offered rates, each in at least three whole runs
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
