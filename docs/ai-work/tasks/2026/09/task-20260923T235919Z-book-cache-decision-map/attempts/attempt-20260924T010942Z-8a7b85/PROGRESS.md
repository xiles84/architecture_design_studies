# Progress — rework the book cache material into the guarantee-first decision map

**Task:** `task-20260923T235919Z-book-cache-decision-map`
**Branch:** `book/cache-decision-map`
**Attempt:** `attempt-20260924T010942Z-8a7b85` (claim `claim-97506e777d7b433d`, epoch 1)
**Capability/role:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Required tag:** `repo/book-cache-decision-map-v1`

Authorship only. No database started, no benchmark lock taken, no measurement.

## Decision Log

1. **Guarantee-first flow, topology last.** The concept now runs contract → change-observation
   completeness → source capabilities → topology/publication protocol → performance controls.
   Cache-aside vs write-through and process-local vs shared are decisions made *after* the guarantee,
   not guarantees. Confidence: High.
2. **The protocol table carries the four evidence statuses.** Each row names one mechanism and the one
   race it closes. The status column separates active registered claim (`v2-12`/`v2-06` lock,
   `v2-15` fence), newer signed but unregistered evidence (the 2026-09-23 churn cell's lease path — no
   registry claim states a lease rate), mechanism demonstrated (the owned process-local version
   check), and proposed / unmeasured (the outbox consumer). Confidence: High.
3. **Relaxed is defined against the strict comparator, and impossible is forbidden separately.**
   "Relaxed" may serve a committed-but-older value for a bounded window; it never means dirty, torn or
   uncommitted. Confidence: High.
4. **The process-local bypass result is stated without inventing a claim.** The legacy arms' zero
   wrong reads are a bypass (0 hits over 15,578–21,756 reads) — the anchor report's figures
   (`20260921T-survey3.md`, the primary anchor of `v2-15`). Adding them to the claim itself would mean
   editing `book/evidence/`, which is forbidden here, so they are cited as report figures beside the
   registered claim. Confidence: Medium (recorded exception; see ANALYSIS.md).
5. **The add-cache gain keeps its `db-only` framing and now sits beside the labelled equal-total arm.**
   2.2–3.5x / p99 15.8–27.1 → 0.8–2.0 ms (`v2-14`); warm 11k vs 8.0k equal-total vs db-only (`v3-03`).
   Never pooled. Confidence: High.
6. **The 2026-09-23 evidence is added where the book was silent.** Real-TTL churn duration, the
   600 s through/aside pair, the per-engine strict/relaxed pair and the locality pair now appear, each
   as a registered v3 claim, with the hard-expiry zero carried as a gap. Confidence: High.
7. **`book/lib/config.typ` was left unchanged.** The protocol table uses existing Typst primitives and
   the existing palette; no new helper was justified. Confidence: High.
8. **No full `book/build.sh` run.** `book/dist` and `book/build.sh` are forbidden to this task; the
   sources were compiled in the pinned image with `--root book` (exit 0). The next full build records
   `evidence_digest_sha256 = 814f40949d162f564ce2a831099727c2f1237a32629f99a2a24dd006db4c47c4`.
   Confidence: High.

## What was produced

- `book/concepts/cache-consistency.typ` — the guarantee-first decision map
- `book/chapters/{major-families,minor-variants,scenarios,how-to-choose}.typ` — corrected cache pointers
- `docs/ai-work/.../attempts/attempt-20260924T010942Z-8a7b85/ANALYSIS.md` — the signed authorship analysis

## Residual risks / limitations

- The bypass counts are report figures behind a registered claim, not claim numbers (decision 4).
- Deletion/tombstone and outbox protocols remain proposed / unmeasured by design.
- The committed `book/dist/` PDF still predates this rewrite; the cache prose changes reach a PDF only
  at the next build (owned by `book-figure-proofset` / the release task).
