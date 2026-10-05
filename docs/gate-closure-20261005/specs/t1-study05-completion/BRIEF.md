# T1: complete Study 05 replication (v3-01, v3-05, v3-06) and gate impossible cache values

**Task:** `task-20261005T151500Z-t1-study05-completion` · **capability:** LOW · **role:** executor · **benchmark:** required
**Branch:** `study-05/t1-completion` · **Worktree:** `.worktrees/t1-study05-completion` · **Study:** `studies/05-cache-consistency`

Read the HIGH plan first: [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md). Its decisions (sizing from the effective guest, replication and sweep form, correctness gates, verified placement, registry discipline) are fixed; do not reopen them.

## Why this exists

Three Study 05 claims are the corpus's last single-run claims and hold the T1 gate open. Two were deferred on a false premise ("yb-cluster3 does not fit the fixed 2 GiB machine"); the effective guest is 8 CPUs and about 15.36 GiB, and yb-cluster3 has completed 30 runs here since. The 2026-10-02 replication also recorded one impossible cache value that the gate did not classify.

## What to do

1. Make the churn gate fail a cell that serves an impossible cache value (a value never written for that key, or older than the floor the study defines), and add a negative control that is seen to fire. This comes first: v3-01 re-runs the churn phase.
2. Re-run v3-01's 1,800 s churn against the fixed 300 s TTL in at least three whole runs, reporting TTL crossings, wrong and impossible reads, and hard versus early expiry.
3. Re-run v3-05 (yb-cluster3 strict write-through against relaxed, client spread over all three endpoints) in at least three whole runs in AB/BA order, with per-endpoint operations recorded.
4. Re-run v3-06 (y1-coloc against y2-noncol) in at least three whole runs. Save yb_table_properties and yb-admin list_tablets before and after every run, adapting studies/01-charity-tree/run-select.sh, and state the observed tablets and leader nodes per layout.
5. Promote, re-scope or limit v3-01, v3-05 and v3-06; narrow v3-gap-01 (its table and leader clause) and v4-gap-01 (the Study 05 key-locality clause) if the placement listing closes them.

## Acceptance criteria

- The impossible-value gate fails a cell and its negative control fires; the change is tested
- v3-01, v3-05 and v3-06 each have at least three fresh whole runs, or an explicit recorded reason why a configuration could not start after effective-guest sizing
- v3-06's layouts are described by observed tablets and leader nodes, never by labels alone
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
