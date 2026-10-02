# Result — T1: replicate Study 04's concurrency claims (`v2-12`, `v2-17`)

**Task:** `task-20261001T133759Z-t4-study04-replication` · **claim:** `claim-e93843f335962c8b` (epoch 1)
**Attempt:** `attempt-20261002T161943Z-1db2e5` · **worker:** deepseek-flash, Deep Code CLI, LOW executor

**Outcome: delivered, committed and submitted for HIGH review.** Both claims are promoted to
`run_replicated` with re-scoped statements. The task's `integration_required` and its review policy put the
merge and the required tag `study-04-t1-replication-v1` with HIGH; nothing here is merged.

## What was delivered

| artifact | where |
|---|---|
| pre-registered protocol (P1-P6 with falsifiers, required-comparison coverage) | `studies/04-configuration-portal/T1-CONCURRENCY-REPLICATION.md` (commit `9d47a10`) |
| 14 whole runs, tagged by the runner | `run/04-configuration-portal/20261002T1624*-w1,w2,w16,w64,w16-r0` and `…20261002T1653*/1654*/1656*/1657*` |
| signed analysis with the scored predictions | `studies/04-configuration-portal/reports/analyses/20261002T-t1-concurrency-replication--deepseek-flash--2026-10-02.md` |
| registry: both claims re-scoped and promoted, anchors added | `book/evidence/v6/claims.json` |
| supersession ledger entry | `book/evidence/v6/SUPERSESSION_LEDGER.md` |
| concurrency prose re-scoped | `book/chapters/minor-variants.typ`, `book/concepts/concurrency-control.typ` |
| context, lessons and the v5-to-v6 carry-forward set | `CONTEXT.md`, `LESSONS_LEARNED.md`, `tools/evidence/internal/evidence/v6_test.go` |
| result commit | `3adc5d3` (worktree `.worktrees/t1-study04`, branch `study-04/t1-replication`) |

## Results

- **`v2-12`'s hedge is retired by its own control.** At 16 writers `c2`/`c1` is **1.637x/1.573x** at the
  default retry budget against **1.690x/1.556x** at `-retries 0`, one pair per design order. In a matched
  order the budget moves the ratio only **1.1%-3.2%**, so "could not separate it from per-retry harness
  work" is no longer true. The claim is a band of five whole runs, **1.56x-1.76x**, with the negative
  control reproducing **93.7%/93.9%** lost at 16 writers (claim: 93.7%) and zero lost in either correct
  design at every level.
- **`v2-17`'s crossover reproduces and is not order-dependent.** The version-check-first runs read
  **0.948x and 0.905x at 1 writer** and 1.035x at 2; the lock-first runs **0.939x** at 1 and 1.106x/1.166x
  at 2; 1.56x-1.69x at 16 and 1.73x-1.96x at 64 in both orders. One 1-writer run measured 3.85x and is
  reported as an **unexplained excursion** — `c1` read 306 ops/s there against 1,286 in a run with the same
  design order, and no client/database throttling was recorded — not as the order effect the first draft
  of the analysis claimed.

## Decisions logged

1. **Added six runs to create a genuine reversed design order.** The tranche's first eight runs were
   produced with eight different `--order-seed` values and treated as "both design orders", but every
   manifest's executed `design_order_pg-single` put `c2_pessimistic_lock` before `c1_optimistic_version`.
   A seed permutes the whole design list; it is not a control on one pair. I read the executed order from
   the manifests, selected seeds whose executed order puts `c1` first, and re-ran levels 1, 2, 16 and 64
   plus a 16-writer zero-retry control. Confidence: high. Affects both claims, the analysis and the book
   prose; it is the reason the conclusions below differ from the first draft.
2. **Corrected the central finding.** The first draft read the 3.85x 1-writer run as an order effect. The
   evidence does not support that: the other run with the same pair order read 0.939x, `c1`'s own throughput
   moved 306 → 1,286 ops/s, and the version-check-first runs put that level at 0.905x/0.948x. The excursion
   is kept in the evidence and in the claim's limits, and called unexplained. Confidence: high. Affects the
   analysis, `v2-17`'s statement/limits and the book paragraph on the crossing.
3. **Scored P1 `falsified` for `v2-17`.** The protocol's named falsifier ("a ratio above 1 at 1 writer")
   fired in one run, so the prediction's verdict is `falsified` even though the crossing reproduces in the
   order-controlled runs; the prediction text says both. **Flagged for HIGH:** the registry records the
   prediction's verdict, so a `run_replicated` claim can carry `verdict: falsified`, which is the
   registry/verdict mismatch T2's review already noted. If HIGH prefers a different convention, that is a
   registry-wide decision, not this task's to change.
4. **Ran the 1-writer level with `c1`/`c2` only** (`20261002T165533Z`, fully green), because
   `x1_lost_update_control` cannot lose an update at one writer and the gate correctly fails that cell.
   Three runs that included `x1` at 1 writer record a failed cell; they are reported in the limits as
   **not applicable**, not hidden or discarded. Confidence: high. Affects `v2-17`'s control_coverage note.
5. **Matched the retry-budget comparison by design order.** The first draft compared 1.64x (c2 first,
   default) against 1.69x (c2 first, zero) — that pair was order-matched, but it also described the
   c1-first pair as the same; both pairs are now recorded separately (3.2% and 1.1%). Confidence: high.
6. **`v2-12`'s attribution is no longer withheld**, because the retry control removes the per-retry-work
   explanation the withholding existed for; `v2-17`'s control_coverage is now `present` because `x1` fires
   at 2 writers and above. Confidence: medium-high (both are direct consequences of decisions 1 and 4).
7. **Counting conventions:** `v2-12` anchors = primary + the four 16-writer runs (5), `whole_run_replications`
   and `replication.runs` = 5; `v2-17` anchors = the nine v6 anchors + all fourteen tranche runs (23),
   runs = 23. `across_run_spread_pct` = 12.4 (`v2-12`) and 8.3 (`v2-17`), both over the 16-writer ratios.
   Confidence: medium; the field is "how many whole runs support a claim", and each single-level run is a
   whole run, but HIGH may prefer a narrower count.

## Validation

- `tools/evidence validate-v6`: **31 claims, 0 errors** (`book/evidence/v6/claims.json`)
- `go test ./tools/evidence/...`: **ok** (includes the v5-to-v6 carry-forward set with `v2-12` added)
- `tools/queue/queue audit`: **0 errors, 3 warnings** — two are the known `t2` live-ref/expired-lease
  warnings, one is the pre-existing `t4` `live_ref_drift` from the claim setting `last_event` empty
- 14 run reports regenerated; each now indexes
  `20261002T-t1-concurrency-replication--deepseek-flash--2026-10-02`
- book scratch build (`book/build.sh --out dist/t4-check.pdf`): `source_dirty: false`, evidence digest
  `8d78efca…` matches `book/evidence/v6/claims.json`, 88 pages, 6/6 fonts, 31/31 claims, no unverified
  banner, no unresolved tokens; scratch artefacts removed, the tracked PDF is for the integration build

## Where the measurement is weak

One fresh-load trial per cell; the 4-, 8- and 32-writer levels were not re-measured; `x1` cannot fire at
1 writer; scale `small`, one hot key, one host, PostgreSQL single node; absolute throughput is day-specific
and the declared budgets still exceed the 2 GiB podman machine.

## Remaining for HIGH

Review the Decision Log and the promoted claims; if accepted, integrate `study-04/t1-replication` into
local `main`, create `study-04-t1-replication-v1`, complete the task and remove `.worktrees/t1-study04`.
