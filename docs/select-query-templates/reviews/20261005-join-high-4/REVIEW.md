# Fourth HIGH review — join/CTE round 3: approved for integration, with a signed qualification

Reviewer: Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH;
selected effort xhigh. Independent of the LOW author (`deepseek-flash` /
opencode CLI). Reviewed submission `ddf7c435`, result `b33044e7`, attempt
`attempt-20261005T111913Z-13271e`, against
[corrective handoff 3](../20261005-join-high-3/CORRECTIVE-HANDOFF.md). I read
the Decision Log first. No database was started; [checks](CHECKS.md) read the
saved results through [audit.go](audit.go).

**Verdict: approve for integration.** Every acceptance item of handoff 3 is
met by the evidence. The v4 analysis has interpretation errors in its
TL;DR, headline and two body statements. Its numbers and pairs are traceable.
These are corrected in a signed HIGH qualification,
[`20261005T-select-join-review-4--claude-opus-5-5`](../../../../studies/01-charity-tree/reports/discussions/20261005T-select-join-review-4--claude-opus-5-5.md),
rather than by a fifth LOW round. Downstream SELECT synthesis and book work
must cite that qualification with v4. No SELECT timing claim is promoted here.

## Decision Log (round 3)

| # | decision | judgment |
|---|---|---|
| 1 | Start the stopped Podman machine | acceptable |
| 2 | Split DDL exactly as handed off, feature-checked first | acceptable |
| 3 | Leader gate through a prepared/ok file handshake with the runner | acceptable — the rejected alternative (prepare twice) would have re-created the tablets it checked |
| 4 | Runner backgrounds the client to serve the handshake | acceptable |
| 5 | Backfill `leader-summary.txt` for the 12 non-distributed cluster runs from their saved raw listings | acceptable — my audit parsed the raw `placement.txt` independently and gets the same leaders |
| 6 | Load leader data at render time, keeping `join.json` bytes unchanged | acceptable |
| 7 | Revert a post-run render fallback so the producing image reproduces the report | acceptable — the correct reproducibility choice |

The three deviations (a heredoc-consumed driver that ran only the `single`
arm, the leader column parse, the empty `SHOW` of a tserver flag) were found,
fixed and recorded. The six `single` runs from the first driver have the same
commit and client image as the other 24.

## Handoff 3 acceptance

- **Distributed arm and gate — met.** All six `distributed` runs report three
  tablets for every relation (`donation_rank_idx` split at 66,667/133,334), the
  gate files and note, and donation leaders on all three nodes before the
  cells. The post-run listing still shows every relation led on all three
  nodes.
- **30 runs, one build — met.** One commit (`ba821e7b`, all 30 annotated
  producing tags), one client image (`b6ea3c00…`), 30/30 valid in every run,
  both controls firing, every per-run digest equal to `sha256(join.json)`.
- **Report — met.** It separates five arms × two phases as classes, renders
  tablets, leaders, endpoints, writer and throttling per run, and re-renders
  byte-identically twice with the producing image (`c2e68309…`).
- **Records — met.** PROGRESS/Decision Log and RESULT are complete. v4 is
  published; v3 is in `reports/outdated/`, and `REPLACEMENTS-select-join.md` is
  extended.
- **Checks — met.** `go vet`/`go test` pass in Podman. Queue audit: 0 errors.
  The tree was clean at submission.
- **Minor items — met, except one.** PostgreSQL `storage_rows_scanned` renders
  n/a and the scan flags are named whole-plan. The shard flag is recorded,
  but as the *configured* container argument, and the analysis explains the
  one-tablet outcome wrongly. The qualification marks the cause unexplained.
  This leaves a recorded gap; it affects no measurement.

## What the data establish (robust in every run and both phases)

- **Splitting effect** (single-tablet → distributed, same cluster, budget,
  data and index): 0.73x/0.71x at 2,000 rows and 0.54x/0.71x at 180,000 rows
  (isolated/blended). The point read is 1.22x slower in the isolated phase,
  where all three runs of each arm are separated, and unresolved in the
  blended phase. This answers review 3's discriminating question: spreading
  data costs the point read and pays at scan sizes.
- **Required pair** (colocated vs distributed): colocated is faster at 20 and
  2,000 rows and slower at 180,000 rows (1.21x/1.54x). Its forced-MAT penalty
  at 2,000 rows is larger in every run (per-run medians: isolated 16.8x vs
  8.5x, blended 11.8x vs 8.4x).
- **Distributed vs the single node**, and vs the equal-total single node:
  slower at the point read and faster at 180,000 rows (0.71–0.76x). The
  equal-total control shows that budget alone does not buy the scan gain.

## Interpretation errors corrected by the qualification

1. **TL;DR baseline.** Its "against the single-tablet cluster … 0.85x … 0.71x
   … 1.19x" uses single-node ratios (0.85x, 0.71x) and a 1.19x that matches no
   pair. Against single-tablet, the figures are 0.73x, 0.54x and 1.22x.
2. **Penalty ratios described as plan costs.** The headline's "1.25x into
   0.87x" and the body's "colocated materializing plan is 1.81x the single
   node's" divide MAT/default penalties. The forced-MAT plan's own time at
   2,000 rows is 1.40x, 0.73x and 1.14x the single node's.
3. **Blended phase.** "Moves these comparisons by under ~10%" is false.
   Colocated vs single node at 180,000 rows goes 0.85x → 1.16x, and
   equal-total vs single at 2,000 rows goes 0.83x → 1.16x.
4. **Leader nodes.** The leader is one node per run, varying between yb-n2 and
   yb-n3 in both single-leader arms. It is not a fixed node per arm.
5. **Single-tablet vs single node.** "1.53x at a point read and 1.31x at 90%"
   is not resolved in the isolated phase (overlap at every threshold).
6. **Shard flag.** "Expected behaviour" contradicts the recorded value
   (2 per tserver × 3 tservers = 6). The cause is unexplained.

The frontmatter `inputs_digest` is a description. The 30 inputs are
identified by the aggregate digest in [CHECKS](CHECKS.md)
(`acc5d121…`).

## Why approve rather than request changes

The measured evidence, gates, report and P4 correction all satisfy the
handoff. The remaining defects are confined to the signed interpretation and
are fully specified above. Another LOW round would cost a full execution and
review for wording; the project's mechanism for disagreeing with a signed
analysis is a signed companion, and the GPT-6 qualification of the SELECT
order correction is the precedent. The two consumers of this result,
`select-synthesis` and `select-book-high`, are HIGH tasks that read reviews.

## Carried forward to `select-synthesis`

- Cite only the pairs above that are separated in every run, and name the
  measured placement (tablets and leaders), never container locality.
- The 1 CPU client is throttled in 32–46% of its CFS periods in every run.
  Point-read comparisons are the most exposed.
- Unsplit hash tables get one tablet despite the configured flag
  (unexplained). The shared SELECT teaching DDL's YugabyteDB hash default
  remains open, as CONTEXT records.
- Next discriminating experiment (not required): a single 6 CPU node with the
  same `SPLIT INTO 3 TABLETS` schema, to separate tablet parallelism from node
  distribution in the scan gain.

## Route

Approve → integrate into local `main` (fast-forward; `main` has not moved
since the branch point) → required tag `study-01/select-v2-join-cte` →
complete. Then `select-synthesis` becomes eligible.
