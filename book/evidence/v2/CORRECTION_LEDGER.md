# Correction ledger — evidence registry v1 → v2

**Attribution.** Written by the LOW execution session `deepseek-flash` on 2026-09-22, task
`task-20260922T170845Z-book-evidence-correction-v2`, implementing the conclusion of
`brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions`
(contribution `contribution-20260922T170212Z-synthesis-01-7465d9`).

**What this ledger is.** The single attributed record of every change from
`book/evidence/claims.json` (v1) to the active correction package `book/evidence/v2/claims.json`
(v2). The brainstorm rejected a free-standing second registry; this ledger is what makes one
active version plus a cited predecessor admissible.

## v1 is frozen and unchanged

- v1 milestone tag: **`repo/book-evidence-registry-v1`** → commit `b7116cf12ef3eb832c53c44911fbefd7b87c3ec6`.
- `git diff repo/book-evidence-registry-v1 -- book/evidence/claims.json book/evidence/claims.schema.json`
  is **empty**; the two v1 files are byte-for-byte unchanged.
- For the record: `sha256(book/evidence/claims.json) = 0768644e0bb5c1d831b5e9a3972f604bd586186f20c45694b7a3ac066d747f01`;
  `sha256(book/evidence/claims.schema.json) = 0be2735aaffc4d781a44142ec454d891b26b0d6853c2d4139613ae58ce39b244`.
- No v1 tag was moved, deleted or reused.

## The mechanical defect v1 had

`tools/evidence v1-diagnostic` recomputes the brainstorm's count from the committed artefacts:

```text
v1-diagnostic: 12 numeric claims, 11 with a missing name, 36/43 names unresolved
```

v1 was *schema-valid* and reported `0 errors` while 84% of its declared cell names occurred
nowhere in the reports they cited. Claim-01's mapping did not merely rename: it shifted the
low-numbered designs and reversed D4/D5 (`d4-rollup-application` where the corpus's D4 is the
rollup **trigger**), so its declared winner could point a reader at the wrong mechanism. Claim-04
declared `failed_cells: 0` for a run whose report says "Cells: 42 … 2 failed". The failure was
**claim-to-cell resolution**, not falsity: most v1 statements are copied verbatim from signed
analyses and are probably right.

## Verdicts and changes

Strength is never promoted. An entry marked *accepted* carries the same or a narrower scope in v2.

| v1 claim | v2 claim(s) | Verdict | Reason and corrected fields |
|---|---|---|---|
| `claim-01-indexes-and-rollups` | `v2-01-normalized-index-first`, `v2-02-rollup-family`, `v2-03-embedding-not-ahead` | **split + corrected** | Compound claim (indexes, rollups, embedding) split into three. Cell names replaced by the report's own labels (`D2 indexed`, `D3 flat+FK`, `D4 rollup/trg`, `D5 rollup/app`, `D6 embedded`). `legacy_provenance_incomplete: true` (no run tag, no commit). Confound `conf-01` attached to the D2→D3 clause. |
| `claim-02-recency-rollup-index` | `v2-04-recency-maintained-index` | corrected | Cells replaced by the report's `d20_recency_flag` / `d22_recency_rollup_idx`; `trials` corrected from `1` to three fresh-load trials per cell; confound `conf-02` (write axis) added. |
| `claim-03-d3-read-advantage-repeated` | `v2-05-d3-read-advantage-repeated` | corrected | Cells replaced by `d3_flattened_fk` / `d2_normalized_indexed`; run tag restored (`run/01-charity-tree/20260913T125342Z-v3`); `repeated_controlled` kept (five fresh loads). |
| `claim-04-hot-drop-sellout` | `v2-06-hot-drop-designs` | corrected | Cells replaced by `P3 CAS`, `C1 count (RC)`, `H0 hold/naive`; **`failed_cells` corrected 0 → 2 of 42**; run tag restored. Control scope recorded as same-run fired. |
| `claim-05-refunds-unanswerable-without-ledger` | `v2-07-refund-answerability-and-storage`, `v2-gap-01-ledger-load-multiplier-unstable` | **split + narrowed** | The answerability clause and the +33-35% storage delta are kept; the "8x load time" clause is **removed** and recorded as a gap because the same pair/code measured 1.2x in two sibling runs. Read rows of that run are recorded as not attributable. |
| `claim-06-ledger-cost-128-buyers` | `v2-08-ledger-race-cost-128-buyers` | accepted with scope | Figures re-derived and kept; scope fixed to 128 buyers, one small run; confound `conf-05` attached. |
| `claim-07-p3-trials-disagree` | `v2-09-postgres-ledger-ratio-unsettled` | **narrowed** | The headline range "1.7-2.1x" was the *sibling* arm's; v2 states the run's own 1.7-5.7x and keeps "PostgreSQL unsettled". Two anchors (both arms). |
| `claim-08-guarded-confirm-refusals` | `v2-10-guarded-confirm-refusals` | corrected | Cells replaced by `S1 conditional` / `S1r retry`; run tag restored. |
| `claim-09-section-seat-map` | `v2-11-section-seat-map` | corrected | Cells replaced by `L2 document` / `L3 sharded`; confound `conf-03` (family/package, not one decision) added. |
| `claim-10-lock-vs-version-check` | `v2-12-lock-vs-version-upper-bound`, `v2-13-trigger-rollup-cost` | **split + narrowed** | Compound claim split. The 1.76x clause is narrowed to an **upper bound** (single trial, one key/16 writers, `-retries 1` not run); the 93.7% control loss and the 6.85x trigger cost are kept. Cells use the report's design ids. |
| `claim-11-cache-throughput-gain` | `v2-14-cache-throughput-gain` | accepted with scope | Cells replaced by `l:na:db`, `l:na:mem:asd:str`, `o:opt:red:asd:str`; `resource_framing: db-only` and the Redis-is-a-cache boundary added; confound `conf-04` attached. |
| `claim-12-strict-freshness-cost` | `v2-16-strict-freshness-read-cost` | accepted, range widened | Cells replaced by the report's scenario ids; the observed range is widened from v1's -8.3%..+12.3% to the full **-14.2%..+12.3%**, still labelled noise; confound `conf-06` attached. |
| `claim-gap-01-study04-missing-designs` | `v2-gap-03-study04-missing-designs` | accepted | Kept as a gap with a `gap_basis`; owned by `task-20260922T025912Z-study04-v2-protocol`. |
| `claim-gap-02-study05-unmeasured-regimes` | `v2-gap-04-study05-unmeasured-regimes` | accepted | Kept as a gap; owned by `task-20260922T025912Z-study05-v2-protocol`. |
| `claim-gap-03-hold-window-sizing` | `v2-gap-02-hold-funnel-unanswerable` | accepted | Kept as a gap with a `gap_basis`. |
| `claim-gap-04-study02-03-validation-open` | — | **retired** | Its content was "validation is still open"; the validation completed 2026-09-22. Listed under `retired_v1_claims`. |

## New in v2 (no v1 counterpart)

- `v2-15-three-instance-staleness` — the Study 05 shared-cache result, with the corrected model
  attribution (87.81% is `o:opt:red:wth:rel`, 87.24% is `l:na:red:wth:rel`).
- `v2-gap-05-colocation-unverified`, `v2-gap-06-native-datastore-families`,
  `v2-gap-07-real-network-balanced-endpoints` — the placement, native-datastore and
  real-network/balanced-endpoint gaps the conclusion requires beside the numbers.

## Provenance model (run tags do not contain their results)

A run tag marks the **producing code state**, not a tree that contains the results. For
`20260921T-survey3`, `20260921T1215Z-small` and `20260921T205212Z` the tag is an ancestor of the
commit that added the results, and for `claim-11` the tag commit equals the registry's
`repo_commit`. A check of the form "does the tag's tree contain the results" therefore fails by
construction and is not a defect. `v2-01`'s run predates the tagging convention and is marked
`legacy_provenance_incomplete`; its tag and commit are never invented.
