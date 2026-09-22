# Synthesis handoff — claims mapped to the fixed six-family taxonomy

**For:** the HIGH `book-synthesis-v1` task. **From:** `task-20260922T170845Z-book-evidence-correction-v2`.
This is a handoff, not chapters: it maps every active claim to a family and a decision model, and
labels each entry **direct evidence**, **analogy**, or **gap**. Evidence levels are conservative and
never promoted above the source analysis.

Active source: `book/evidence/v2/claims.json`. Do not write numeric prose from v1
(`book/evidence/claims.json`); it is frozen at `repo/book-evidence-registry-v1` and its
claim-to-cell resolution failed (see `CORRECTION_LEDGER.md`).

## Family → claims

| Family | Claims | Level |
|---|---|---|
| normalized facts | `v2-01-normalized-index-first`, `v2-06-hot-drop-designs`, `v2-10-guarded-confirm-refusals`, `v2-12-lock-vs-version-upper-bound`, `v2-gap-07-real-network-balanced-endpoints` | direct (v2-01/06/10/12); gap marker on v2-gap-07 |
| rolldown / duplicated data | `v2-05-d3-read-advantage-repeated`, `v2-gap-05-colocation-unverified` | direct; the family's headline result is the one repeated_controlled claim |
| rollups / materialized aggregates | `v2-02-rollup-family`, `v2-04-recency-maintained-index`, `v2-13-trigger-rollup-cost` | direct, all confounded as noted |
| embedded documents | `v2-03-embedding-not-ahead`, `v2-11-section-seat-map`, `v2-gap-06-native-datastore-families` | direct within PostgreSQL/state; **analogy gap** to native document databases |
| append-only history / snapshots | `v2-07-refund-answerability-and-storage`, `v2-08-ledger-race-cost-128-buyers`, `v2-09-postgres-ledger-ratio-unsettled`, `v2-gap-01-ledger-load-multiplier-unstable`, `v2-gap-02-hold-funnel-unanswerable` | direct; two gap markers |
| external read copies / caches | `v2-14-cache-throughput-gain`, `v2-15-three-instance-staleness`, `v2-16-strict-freshness-read-cost`, `v2-gap-04-study05-unmeasured-regimes` | direct within the cache framing; Redis is a cache, never an authoritative store |

## Analytic models, in the order the conclusion fixed

1. **Answerability before performance — a derived rule, not a measured ranking.** Label it
   *derived* from two measured answerability failures: Study 02's refunds window unanswerable in
   14 of 15 designs (`v2-07`) and Study 03's hold funnel unanswerable everywhere (`v2-gap-02`).
2. **The multidimensional design vector** (family, variant, technology, topology, workload) —
   every claim's `support.key` is one point on that vector.
3. **The derive → index-an-existing-fact → materialize → copy/cache ladder.** Direct evidence:
   `v2-01` (index), `v2-02`/`v2-13` (materialize), `v2-14` (copy/cache). Analogy only when
   crossing to a native family (`v2-gap-06`).
4. **Contention surface / arbitration.** Direct: `v2-06` (CAS vs counter), `v2-12` (lock vs
   version), with the upper-bound caveat.
5. **Query-boundary portfolio.** Direct: the report-cost results — `v2-09`/`v2-08` (ledger
   reports), `v2-11` (seat-map vs point-lookup), `v2-16` (strict fences on the write path).
6. **Evidence-confidence cards.** Use each claim's `strength`, `trials`, `confounds`, `limits`
   and `review` verbatim; do not round a confidence up.
7. **Failure-mode catalogue.** Direct: `v2-15` (three-instance shared cache), `v2-12`'s control
   (lost updates), `v2-07`/`v2-gap-02` (erased history), the failed Study 02 cells.

## What must be published beside the numbers

- The **six confounds** (`book/evidence/v2/confounds.json`), referenced by the claims they affect.
- The **coverage matrix** (`book/evidence/v2/coverage.json`) with review state, including the
  placement, native-datastore and real-network gaps.
- The environment limitation wherever a number appears: one laptop's shared cores, no real
  network, closed-loop demand.
- The **run-tag provenance model**: a run tag marks producing code, not a tree containing results.

## Do not do

- Do not quote a cross-study numeric ranking; no comparability record exists.
- Do not publish `claim-07`'s old "1.7-2.1x", `claim-05`'s old "8x load", or the all-green
  analysis's swapped three-instance rows; the corrections are in `v2-09`, `v2-gap-01` and `v2-15`.
- Do not read missing designs or topologies as null results (`v2-gap-03`, `v2-gap-04`).
