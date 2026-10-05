# Progress — SELECT synthesis (HIGH analyst)

**Task:** `task-20261002T131729Z-select-synthesis`. **Attempt:**
`attempt-20261005T133549Z-3446d3`, claim `claim-520f55fdc01e5bf7` (epoch 1, real clock).
**Agent:** Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH / analyst; effort xhigh.
**Worktree:** `.worktrees/select-synthesis`; branch `study-01/select-synthesis`.

The four dependencies are completed:

- access paths, aggregate/top-N and existence, with their order correction;
- join/CTE, approved by review 4 on 2026-10-05.

No database or benchmark was started.

## Decision Log

1. **Score only the randomized-order runs for P6–P11.** GPT-6's evidence review showed the
   fixed-order primaries were order-confounded. The 18 `select-order-*` runs replace them.
   - Confidence: High.
   - Alternative: pool both corpora (rejected; it would mix schedules).
   - Affects: every P6–P11 verdict.
2. **"Same plan" means an identical normalized plan tree.** The tree keeps node types, relations
   and index names; costs, row estimates and timings are removed. This is the P1/P2 falsifier's
   "identified plan difference".
   - Confidence: Medium-High.
   - Alternative: raw EXPLAIN text, which differs on every run.
   - Affects: `v6-sel-01`.
3. **Stable-sign rule.** A direction counts only when all runs agree, and its size is compared with
   the family's position drift. A stable sign inside the drift is "unresolved in magnitude"
   (examples: PG `LATERAL` at one parent, PG pre-aggregation at 10,000 parents).
   - Confidence: High.
   - Affects: P8, P9 and claim scopes.
4. **Generate a numbers-only evidence report for the synthesis.**
   `docs/select-query-templates/synthesis-20261005/evidence.go` writes
   `reports/20261005T-select-synthesis-evidence.md`, which prints each run's `sha256(json)` and the
   cell tokens. The registry requires an anchor's digest verbatim in its report and analysis, and
   the SELECT cross-run reports print no digest. The generator is deterministic: two generations
   were byte-identical.
   - Confidence: Medium-High.
   - Alternatives: change the harness renderers and rebuild images (out of scope); anchor to the
     input manifests (they hold digests but no cells).
   - Affects: all `v6-sel-*` anchors.
5. **Anchor environment.** The order-correction runs name the environment page
   `xiles-podman-select-order-correction-20261004`, which documents the same host and slot. Their
   anchors use `host-zenbook-ux5406sa` (the schema requires `host-*`), and each such claim says so
   in a limit.
   - Confidence: Medium.
   - Affects: `v6-sel-05` … `v6-sel-15`.
6. **The run-tag pattern accepts `run/study-0N/…`.** The SELECT runners tag runs as
   `run/study-01/…`, and the validator accepted only `run/0N-…`. The pattern now accepts both
   namespaces, and a test pins the accepted and rejected forms.
   - Confidence: High.
   - Alternative: null run tags with `legacy_provenance_incomplete` (false; the tags exist).
   - Affects: `tools/evidence` and all `v6-sel-*`.
7. **Strength honesty.**
   - `run_replicated` only where a direction was recorded before the runs: protocol P1–P12; review
     3's placement prediction for `v6-sel-16`; the fixed-order `DISTINCT ON` hypothesis re-tested
     under randomized order for `v6-sel-12`.
   - `repeated_controlled` otherwise: `v6-sel-06`, `-08`, `-14`, `-17`.
   - No-difference claims (`v6-sel-01`, `-06`) state that they are bounded non-observations, not
     equivalence results. The schema has no equivalence strength.
   - Confidence: Medium.
8. **Narrow `v4-gap-01` and `v4-gap-02`** (registry rule: re-read the gaps the new evidence
   touches). Both become `partially_superseded` by `v6-sel-16`/`v6-sel-17`. Their statements now
   say where placement and endpoint balance are verified; anchors, trials and confounds are
   unchanged. They are listed as re-scoped in the carry-forward test and the ledger.
   - Confidence: High.
9. **No SELECT report regeneration.** No SELECT renderer emits an analyses index (checked every
   SELECT cross-run report), so the brief's regeneration step is a no-op. Re-rendering would only
   reproduce the same bytes.
   - Confidence: High.
10. **YugabyteDB physical-arm limit is a registered gap** (`v6-sel-gap-01`), not a footnote. YB
    access plans used an index in 0/103 cells and existence in 0/64. Their YB conclusions describe
    scan plans.
    - Confidence: High.
    - Affects: P10, P11 and the book handoff.
11. **Not registered:** PG pre-aggregation at 10,000 parents (0.77–0.80x, inside the 29% drift)
    and the unexplained stable YB-blended 200-row spelling signs (identical plans). Both are
    reported in the analysis only.
    - Confidence: High.

## Status

Complete against the brief's acceptance list.

- Signed analysis with verdicts for P1–P12 and a selection matrix.
- 18 registered claims, validator green.
- Book handoff for `select-book-high`.
- CONTEXT and LESSONS updated.

Submitted for labelled HIGH self-review.
