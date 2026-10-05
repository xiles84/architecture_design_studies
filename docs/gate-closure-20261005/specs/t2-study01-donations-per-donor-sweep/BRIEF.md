# T2: sweep donations per donor for Study 01 embedding against normalized and indexed designs

**Task:** `task-20261005T151500Z-t2-study01-donations-per-donor-sweep` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-01/t2-donations-per-donor` · **Worktree:** `.worktrees/t2-study01-donations-per-donor-sweep` · **Study:** `studies/01-charity-tree`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

v2-03 found embedding (D6) not ahead of the indexed designs at the study's one child-count distribution. Embedding's expected advantage and its write penalty both depend on how many children a parent holds; the charter asks for that axis.

## What to do

1. Use the existing max-per-person knob to set at least five donations-per-donor levels (for example 1, 5, 20, 100, 500), holding total donations comparable across levels (methodology 8a: change the shape, hold the volume).
2. Compare D2 (indexed), D6 (embedded JSONB) and D9 (embedded hybrid), plus D1 as the baseline, on reads and inserts at each level.
3. Report per-run ratios and sign counts; state where, if anywhere, embedding overtakes the indexed design on reads, and how its insert cost grows.
4. Register the result with a crossover record (levels at least 5); this task covers the owner's required embedding comparison for this axis.

## Acceptance criteria

- At least five donations-per-donor levels, each in at least three whole runs, every cell oracle-checked
- Read and write sides are both reported; an embedding read gain is never shown without its insert cost
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
