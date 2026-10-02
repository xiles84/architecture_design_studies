# Result — T1: replicate Study 05's cache, churn and engine claims

**Task:** `task-20261001T133759Z-t5-study05-replication` · **claim:** `claim-4d1e0f9c08825d9d` (epoch 1)
**Attempt:** `attempt-20261002T184227Z-2c6d80` · **worker:** deepseek-flash, Deep Code CLI, HIGH session, executor

**Outcome: delivered, committed and submitted for HIGH review.** Eight whole runs in reversed pair order;
four claims promoted to `run_replicated`, one re-scoped with a retired clause, three recorded as
deferred/gap. Not merged (the task's review policy and `integration_required` put the merge with a reviewer).

## What was delivered

| artifact | where |
|---|---|
| pre-registered protocol (P1-P8 with falsifiers, feasibility, coverage table) | `studies/05-cache-consistency/T1-REPLICATION.md` (commit `37ce688`) |
| 8 whole runs, tagged | `run/05-cache-consistency/{20261002T185604Z,185911Z,190350Z,190748Z,191221Z,191405Z,191609Z,194046Z}` |
| signed analysis | `studies/05-cache-consistency/reports/analyses/20261002T-t1-replication--deepseek-flash--2026-10-02.md` |
| registry: four promotions, one re-scope, three limits | `book/evidence/v6/claims.json` |
| ledger, book prose, CONTEXT, LESSONS, carry-forward test | `book/evidence/v6/SUPERSESSION_LEDGER.md`, `book/chapters/*.typ`, `book/concepts/cache-consistency.typ`, `CONTEXT.md`, `LESSONS_LEARNED.md`, `tools/evidence/internal/evidence/v6_test.go` |
| tracked book draft rebuilt | commit after `178c2dc` (`book/dist/`, 89 pages, digest `6c5b2f80…`) |
| result commit | `178c2dc` on `study-05/t1-replication` |

## Results

- **`v2-16`, `v3-02`, `v3-03`, `v3-04` promoted to `run_replicated`.**
- **`v2-14` re-scoped:** throughput 2.54x and 2.24x the no-cache baseline (primary band 2.2-3.5x), but p99
  fell 27 → 9-10 ms (about 3x), **not** the order of magnitude the primary recorded; the tail clause is
  retired and the book's three quoting sites corrected.
- **`v3-02` finding:** the 600 s crossings and 0 wrong reads reproduced in both orders, but one aside cell
  recorded **1 impossible cache value in 1,173,729 reads** while its gate passed — the study's floor says
  that must be zero, so the claim keeps a limit and the gate's coverage is recorded as a gap.
- **Deferred/gaps:** `v3-01`'s 1800 s horizon not re-run; `v3-05`/`v3-06` need `yb-cluster3` (2 GiB machine)
  and there is no tablet/leader placement record, so the required colocated/non-colocated comparison stays
  an explicit coverage gap.

## Decisions logged

1. **Measured five of eight targets and deferred three with reasons.** `yb-cluster3` cannot run on the fixed
   2 GiB machine, and the 1800 s churn was left to the 600 s pair as the mechanism re-test. Confidence: high;
   matches the brief's "re-tested or explicitly deferred".
2. **Read the executed pair order from each manifest, not the seed.** Reverse-order seeds were chosen and
   verified; A/C/D pairs genuinely reverse. Confidence: high (the t4 lesson).
3. **Retired `v2-14`'s p99 clause rather than its headline.** The throughput relation reproduced; the tail
   did not. Suspected client-bound but not isolated. Confidence: medium-high; flagged for review.
4. **Promoted `v2-16`/`v3-04` as nulls and recorded the primary's yb direction as not reproduced.**
   Confidence: high for v2-16, medium for v3-04 (large spreads).
5. **Promoted `v3-03` on its labelling relation, not its point figures** (11k/8.0k did not reproduce).
6. **Promoted `v3-02` on its crossing/wrong-read relation while adding the impossible-value limit.**
   Confidence: medium — flagged for HIGH whether a cell that breached the study's floor should be promoted.
7. **Recorded the gate gap, not a design defect.** The impossible value is an observation the gate did not
   classify; a follow-up must make it fail the cell. Confidence: high.
8. **Anchors/replication counts:** v2-14 3, v2-16 3, v3-02 4, v3-03 4, v3-04 3; `replication.order` "abba";
   `across_run_spread_pct` only where a fixed comparison exists (12.6 / 12.9), null otherwise.

## Validation

- `tools/evidence validate-v6`: **31 claims, 0 errors**
- `go test ./tools/evidence/...`: **ok** (carry-forward set updated: v2-14/v2-16/v3-02/v3-03/v3-04 moved to `changed`)
- `tools/queue/queue audit`: **0 errors, 3 warnings** (two known t2; one t5 live-ref warning that clears on submit)
- each regenerated report indexed by its analysis digest
- book build: `source_dirty: false`, evidence digest `6c5b2f80…` matches the registry, 89 pages, 6/6 fonts, 31/31 claims

## Where the measurement is weak

One machine, closed loop, 2 GiB; `v3-04`'s spreads reach 74.7%; the impossible value is one observation with
no cause; hard expiry fired 0 times; the p99 clause was not isolated; `v3-01`/`v3-05`/`v3-06` unmeasured.

## Remaining for HIGH

Review the Decision Log and the five promotions/re-scope; decide decision 6 (promote `v3-02` despite the
impossible value); if accepted, integrate `study-05/t1-replication` into `main`, create
`study-05-t1-replication-v1`, complete the task and remove `.worktrees/t1-study05`.
