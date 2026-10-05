# T2: sweep hot-parent skew for Study 01 rolldown, rollup and embedding

**Task:** `task-20261005T151500Z-t2-study01-hot-parent-skew-sweep` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-01/t2-hot-parent-skew` · **Worktree:** `.worktrees/t2-study01-hot-parent-skew-sweep` · **Study:** `studies/01-charity-tree`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

Copied keys (D3), maintained totals (D4) and embedded children (D6) concentrate write work on parents; their cost under a hot parent is the boundary the charter names and no run has measured.

## What to do

1. Add a deterministic, seeded skew knob to the generator and the Go oracle: the share of new donations that land on the hottest charity (and its donors), at least five levels from uniform to heavy (for example uniform, 5%, 20%, 50%, 80%).
2. Test the knob first: the generated distribution matches its declared share, and the oracle agrees with the database at every level.
3. Compare D2, D3, D4 and D6 on blended reads and on inserts at each level; record lock waits or retries where the engine exposes them.
4. Report per-run ratios and sign counts and any reversal with its interval; register with a crossover record (levels at least 5).

## Acceptance criteria

- The skew knob is tested and recorded in every manifest; at least five levels, each in at least three whole runs
- Write-side contention is reported beside read gains for every design
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
