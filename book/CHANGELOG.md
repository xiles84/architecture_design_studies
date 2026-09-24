# CHANGELOG — *Data Architecture Reference*, Edition 2 revision (2026-09-24)

This revision applies the owner-supplied external review of the Edition 2 PDF. It changes prose and
figures only: **no measured number, claim id, trial count, evidence classification, confound or gap was
changed**, and `book/evidence/v5/` is untouched. The revision was executed as two queue tasks
(`book-visual-correctness-a`, `book-visual-tranche-b`) and released by
`task-20260924T170502Z-book-v2-1-release`. Tags: `repo/data-architecture-book-visual-correctness-a`,
`repo/data-architecture-book-visual-tranche-b`, `repo/data-architecture-book-v2-1`.

## Technical corrections

- **Relaxed freshness is a family, not one contract.** "Relaxed" no longer reads as permitting an older
  committed value "for a bounded window". Only *bounded staleness* declares a maximum Δ; the measured
  relaxed cells declared none and are labelled unbounded-relaxed. (`chapters/minor-variants.typ`;
  consistent with `concepts/cache-consistency.typ` and the glossary.)
- **Application-maintained rollups: transaction first, publication only across the boundary.** The
  rollup family no longer says an application rollup "must be atomic with the publication". Where base
  fact and rollup share a database they are updated in one transaction; where they cannot, the copy
  crosses the authoritative boundary and delivery, idempotency and reconciliation become relevant.
  (`chapters/major-families.typ`, `chapters/how-to-choose.typ`.)
- **Placement text equivalent no longer contradicts its own figure.** "Each of the five placements is a
  copy" became "derive stores no additional copy; the remaining rungs introduce some stored derived
  structure or copy". (`concepts/information-placement.typ`.)
- **Pessimistic locking is not a categorical latency/retry claim.** It is framed as waiting/blocking
  rather than optimistic retries; observed latency may rise or fall by workload; application retries can
  still be required for deadlocks, lock timeouts, serialization failures or other transient failures.
  (`concepts/concurrency-control.typ`.)
- **Derived state is separated from retained history.** A rollup/rolldown/external copy optimizes an
  answer from current state; an event/ledger/audit/temporal history/snapshot preserves information that
  would otherwise disappear. "The only structures that can answer" became the current-state-only claim,
  and append-only storage no longer implies a concurrency fence. (`concepts/derived-state-and-history.typ`,
  `chapters/major-families.typ`.)
- **Cache terminology audit.** `glossary.typ` now defines invalidation, fill lease, fencing token,
  source-generation validation, cache CAS, hit validation, transactional outbox and TTL/early expiry
  separately; the cache chapter's "Publication fence / CAS" card and matrix row were renamed to name
  **source-generation validation**, and the chapter intro states how it differs from a sink-side fencing
  token and from cache CAS. (`concepts/cache-consistency.typ`, `glossary.typ`.)

## Clarified claims

- No measured value moved. The clarifications that keep existing numbers scoped are: the strict/relaxed
  read-throughput range (−14.2% to +12.3%) stays a read-only result with no write-path figure; the
  pessimistic 1.76x stays scoped to one hot key, 16 writers, one run; the application-rollup write claim
  is scoped to the measured configuration-portal workload and does not generalise "writes are cheaper";
  and the measured relaxed cells are explicitly *not* bounded-staleness experiments.

## Added visuals

- Normalized / rolldown / rollup on one Donor–Donation relation (ER, with cardinalities).
- Rolldown vs rollup: parent→children fan-out against children→parent aggregation.
- Normalized child rows vs a bounded child collection inside the parent (embedding).
- The index path vs the hot row: two different tools.
- Optimistic version-checking vs pessimistic lock-before-update, with the invariant recheck named as
  what rejects an invalidated operation.
- Derived state vs retained history: one canonical fact, two stored branches.
- Current state vs retained history: undoing the state is not erasing the evidence.
- The cache as a non-authoritative read copy beside the authoritative database.
- Logical topologies inside one physical host, separating logical partitioning from unverified physical
  placement.

## Modified diagrams

- **A lease, a fencing token, and a source generation.** Redrawn with five lifelines (R1, R2,
  lease/token coordinator, source, cache/sink). The coordinator no longer appears to publish to the
  cache, the second filler has its own lifeline, and two panels plus a distinction note separate a lease,
  a sink-side fencing token (the sink remembers the highest accepted token, so arrival order at the sink
  matters) and source-generation validation (the authoritative generation is consulted before
  publishing).
- **The stale-fill race.** Redrawn with explicit `S0/v0 → S1/v1` provenance, the mechanism named
  source-generation validation, an authoritative-read path for the second reader, and an explicit scope
  note that this figure demonstrates the stale-fill race only — the strict-after-acknowledgement contract
  and its acknowledgement boundary stay in the freshness-timeline figure.
- The combined history figure was split into two for reading order and legibility.

## Unresolved questions / gaps

- **Item 2.9 (a "can I quote this number?" flowchart) was deliberately not added.** The existing
  reading chapter keeps its checklist as a table, which is searchable and cannot be letter-spaced; the
  "a gap is not a zero effect" point is made in the prose. This is a reasoned rejection, recorded.
- **The visual inspection is one model's.** Figures were rasterised and inspected for clipping, overlap
  and label size, but no independent reader or comprehension test was run.
- **Two tall sequences render below full width** (lease-vs-fence at 74%, optimistic-vs-pessimistic at
  72%) to avoid overflowing the page footer. Legible in the raster; a print reader may disagree, and the
  fallback is to split each into two full-width figures.
- **The book grew from 54 to 61 pages** across the revision.
- **No new measurement was taken**, so every coverage gap in the registry is unchanged.

## Release

- Artefact `book/dist/data-architecture-reference.pdf`, 61 pages, 29/29 claims indexed, 6/6 fonts
  embedded, evidence digest present, no unverified banner, no unresolved template tokens.
- Tag `repo/data-architecture-book-v2-1` (hyphen, not `v2.1`, because the queue task-id slug cannot
  contain a dot).
