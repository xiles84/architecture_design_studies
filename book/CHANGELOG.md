# CHANGELOG — *Data Architecture Reference*

## Edition 3 preview 6 — 2026-10-05

Adds the appendix *Query plans behind the SELECT claims*, which the owner asked for. It quotes the
stored `EXPLAIN ANALYZE` plans for every measured SELECT claim (`v6-sel-01`..`18`); `v6-sel-07` is a
negative control that records rows, not a plan, and the appendix says so. There are 45 plans in 25 exhibits,
with a reading note per claim. The plans are extracted by `tools/evidence plans` from committed result
files. Each exhibit's run must be an anchor of its claim, and the file's sha256 must equal the
anchor's digest. Every same or different node-tree statement is recomputed, and `go test` fails if
the generated data is stale. The trim rule, kept attributes and the full plan's location are printed.
The Charity template guide and the reproduction section link to the appendix. Milestone:
`repo/book-plan-appendix-v1`.

## Edition 3 preview 5 — 2026-10-05

Answers the owner's question of when a derived table outperforms the other SELECT forms. The
Charity case now says plainly that the filtered derived table ran the same plan as `JOIN … WHERE` in
every measured pair and never outperformed (`v6-sel-01`). It adds the one measured derived table that
changes the plan, pre-aggregation (`v6-sel-18`, registered from existing evidence): it lost at
1–1,000 donors on PostgreSQL, dipped only inside the drift at 10,000, and never won on YugabyteDB.
Theory explains when a derived table stops being a spelling; the template guide names both results.
Milestone: `repo/book-derived-table-clarity-v1`.

## Edition 3 preview 4 — 2026-10-05

Publishes the reviewed SELECT results in the Charity donations case: join forms and the CTE
fence, CTE reuse, membership and absence, totals and each donor's top donation, pagination and
disjunctions, verified three-node placement, and a compact template guide. The 18 claims of the
SELECT synthesis (`v6-sel-*`) are cited. YugabyteDB's scan-only access plans and the missing cluster
arms are shown as a coverage gap.

Each SELECT theory strategy now links to its measured variation; theory keeps expectations without
numbers. The `repeated_controlled` badge reads "repeated — no prior prediction". The glossary adds
tablet and leader, semi-join and anti-join, keyset pagination and position drift. The cache case's
placement gaps name the verified Charity placement. Frozen evidence, signed analyses and producing
tags are preserved. Milestone: `repo/book-select-query-templates-high-v1`.

## Edition 3 preview 3 — 2026-10-04

Adds theory lessons for six equivalent join/filter forms, LEFT JOIN semantic
boundaries and CTE folding/materialization. Study 01 gains a clearly planned
case variation, without new timings. LOW resumes the existing measurement task;
HIGH reviews its evidence and publishes the later case results. The new HIGH
publication task supersedes the original LOW publication assignment.

Frozen evidence, signed reports and producing tags are preserved. Milestone:
`repo/book-join-theory-20261004`.

## Edition 3 preview 2 — 2026-10-04

Replaces the mixed pilot/older-edition reading path with the owner's requested hierarchy:
Theory → six topics → twenty strategies, then five business cases with matching variations,
thrive/die situations and compact measured summaries. Adds concise SELECT examples, shared
lesson layouts and linked navigation. Detailed evidence stays in the appendix.

All 31 active claims remain cited. Original figures, frozen evidence and signed analyses are
preserved; older chapter sources remain in the repository. Queue guidance is conceptual and
new SELECT timing results are not promoted. Milestone: `repo/data-architecture-book-v3-preview-2`.

## Edition 3 preview 1 — 2026-10-04

Owner-authorized local release of the reviewed concept-first pilot combined with all committed
repository work. Adds the reading guide, validated chapter components, compact linked evidence
badges, the Rollup chapter and the Study 04 profile. The retained Edition 2 chapters and concepts
are explicitly marked as awaiting conversion. This preview does not claim the full rewrite is done.

The active v6 evidence package remains unchanged; SELECT correction results are preserved in the
repository, without adding their timing claims to the book. Built with the canonical pinned Podman
toolchain, source/figure/evidence gates and rendered review. Details and cleanup provenance:
`docs/book-release-20261004/VALIDATION.md`. Release tag: `repo/data-architecture-book-v3-preview-1`.

## Edition 2 revision — 2026-09-24

This revision applies the owner-supplied external review of the Edition 2 PDF. It changes prose and
figures only: **no measured number, claim id, trial count, evidence classification, confound or gap was
changed**, and the then-active v5 evidence package is untouched. The revision was executed as two queue tasks
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
