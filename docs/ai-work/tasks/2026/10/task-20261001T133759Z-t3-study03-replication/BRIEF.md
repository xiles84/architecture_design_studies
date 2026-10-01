# T1: replicate Study 03's claim v2-10 on a machine that can carry it, or record the gap

**Task:** `task-20261001T133759Z-t3-study03-replication` · **capability:** LOW may execute · **role:** executor · **benchmark:** required
**Branch:** `study-03/t1-replication` · **Worktree:** `.worktrees/t1-study03`

## Why this exists

The corpus's claims are only as strong as their replication. A survey on 2026-10-01 found the registry at
9 `run_replicated`, 1 `repeated_controlled`, 1 `mechanism_supported`, 13 `single_run_directional` and 7
`gap`. This task is one of five that take the remaining single-run claims, study by study, and either
replicate them run-level or state precisely what could not be replicated and why.

**Do not re-read the old reports and call it replication.** The unit of work is a fresh whole run.

## What to do

1. Write and commit the protocol first: the claims under test, the run configuration (matching the primary's
   own scale and durations), the designs, the run order, and **predictions in bands and orderings with named
   falsifiers** - not point figures, because a point figure from one run did not survive the last three
   replications.
2. Record the owner's required-comparison coverage table in the protocol before running (sizing; rollup,
   rolldown and embedding; colocated versus non-colocated data; optimistic versus pessimistic concurrency),
   marking each covered, not applicable with a reason, or outstanding.
3. Run at least two whole runs in reversed or AB/BA order, one at a time under `ads-run-lock`, from a clean
   committed tree; tag every run `run/<study>/<run-id>`. v2-10 rests on guarded-confirmation refusals, which needed the YugabyteDB arm; that arm is the part the machine may not carry.
4. Analyse the fresh runs and score every prediction honestly, including the ones you get wrong.
5. Update the registry: promote what replicated, with the statement re-scoped to what actually travelled;
   leave the rest with a limit naming the relation that failed.
6. Re-scope any book prose quoting a superseded point figure; update `CONTEXT.md` and `LESSONS_LEARNED.md`;
   regenerate the reports; merge into local `main` and tag.

## Acceptance criteria

- A feasibility check on the fixed 2 GiB machine is done before any run and its result is recorded
- If it fits: two whole runs in reversed order, predictions scored, claim promoted or limited
- If it does not fit: the claim keeps its strength and gains an explicit coverage gap naming the resource limit - no partial run is presented as a replication
- A protocol with predictions is committed BEFORE any run, naming its bands or orderings and its falsifiers
- At least two whole runs per target, in reversed or AB/BA design order, tagged run/<study>/<run-id>
- Every cell passes its correctness gate; failures abort the cell and are reported
- A signed analysis reads the fresh runs, scores every prediction honestly, and names where the measurement is weak
- Each target claim is either promoted to run_replicated with its statement re-scoped to what replicated, or keeps its strength with a limit saying which relation did not reproduce
- Any book prose quoting a superseded point figure is re-scoped, and CONTEXT.md and LESSONS_LEARNED.md are updated
- The task is merged into local main and tagged; the worktree is removed

## Validation

- tools/evidence validate-v6: 0 errors
- go test ./tools/evidence/... passes, including the v5-to-v6 carry-forward sets
- tools/queue/queue audit: 0 errors
- each regenerated report indexed by its analysis digest
- book build clean (source_dirty false) with the evidence digest matching the registry

## Escalate only if

A target cannot be measured at all on the fixed machine (record it as a coverage gap instead of escalating),
or a promotion would require changing a claim's scientific question (never do that - escalate).
