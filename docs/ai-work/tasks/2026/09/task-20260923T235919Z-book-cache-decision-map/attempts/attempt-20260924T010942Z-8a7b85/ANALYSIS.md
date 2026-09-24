---
analysis_id: book-cache-decision-map--deepseek-flash--2026-09-23
run_id: none (book authorship against the active evidence registry)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH capability, Deep Code CLI; effort not exposed. Executor of task-20260923T235919Z-book-cache-decision-map."
analyzed_at: 2026-09-23
inputs_digest: 814f40949d162f564ce2a831099727c2f1237a32629f99a2a24dd006db4c47c4
repo_commit: c98e94fcb37773b809280d7338875b96b2d17e86
supersedes:
status: current
headline: The book cache material now follows a guarantee-first decision flow with a four-status protocol table, names the race each mechanism closes, and states the process-local zero-wrong-read arms as bypass results; every number resolves to the active v3 registry.
---

# Analysis — cache decision map — deepseek-flash

> Deliverable of `task-20260923T235919Z-book-cache-decision-map`, implementing the book clause of
> `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`. Authorship only: no
> measurement, no database, no benchmark lock.

## TL;DR

- The cache concept now runs **contract → change-observation completeness → source capabilities →
  topology/publication protocol → performance controls**, with topology and fill path demoted to
  choices made after the guarantee, not guarantees themselves.
- A protocol table gives every named mechanism the one race it closes and an **evidence-status**
  column with four values: active registered claim, newer signed but unregistered evidence, mechanism
  demonstrated, proposed / unmeasured.
- The process-local arms' zero wrong reads are stated as a **bypass** result (0 hits over
  15,578–21,756 reads), and *relaxed* is defined only against the strict comparator, never as a dirty
  or impossible value.

## What changed and why

| File | Change |
|---|---|
| `book/concepts/cache-consistency.typ` | rewritten as the guarantee-first decision map: the five-step flow, the four-status protocol table, the comparator definition of relaxed, and the 2026-09-23 churn/engine/placement evidence |
| `book/chapters/major-families.typ` | the cache family's mechanism/costs/boundaries now defer to the flow and name the stale-fill race; equal-total card added |
| `book/chapters/minor-variants.typ` | strict/relaxed defined against the comparator; the write-path fence names its race |
| `book/chapters/scenarios.typ` | donor-portal scenario corrected to 87.81/87.24 and the bypass caveat |
| `book/chapters/how-to-choose.typ` | the worked decision names the contract and the stale-fill race |

## Numeric provenance

Every number in the reworked material resolves to an active `book/evidence/v3/claims.json` claim:
2.2–3.5x, 15.8–27.1 ms and 0.8–2.0 ms (`v2-14`); 87.81% and 87.24% (`v2-15`); −14.2%..+12.3%
(`v2-16`); 11k and 8.0k (`v3-03`); 1,801 s / 6.0 / 0 / 2,879,077 / 601 s / 2.0 (`v3-01`, `v3-02`);
514 / 596 (`v3-04`), 942 / 846 (`v3-05`), 882 / 673 / 852 / 520 (`v3-06`).

**One honest exception, recorded rather than hidden.** The bypass counts 0 hits and 15,578–21,756
reads are the anchor report's figures (`studies/05-cache-consistency/reports/20260921T-survey3.md`,
the primary anchor of registered claim `v2-15`), not a number in the claim's statement. Promoting them
into the claim would mean editing `book/evidence/`, which this task's `forbidden_paths` excludes. They
are cited as report figures beside the registered claim; a later evidence task may fold them into a
claim.

## Where this authorship is weak

1. **No measurement.** This is prose over existing claims; it adds no evidence and cannot strengthen a
   claim.
2. **The unmeasured branches remain unmeasured.** Deletion/tombstone protocols and the outbox consumer
   are marked *proposed / unmeasured*, deliberately not filled with words.
3. **The bypass exception above.** The recommendation that a private cache needs per-hit version
   validation rests on the owned process-local arms, which are mechanism-level, not a registered rate.
4. **The book is not rebuilt into `book/dist/`.** `book/dist` is a forbidden path here; the sources
   compile (verified in the pinned image), and the next full build writes
   `evidence_digest_sha256 = 814f40949d162f564ce2a831099727c2f1237a32629f99a2a24dd006db4c47c4`.
