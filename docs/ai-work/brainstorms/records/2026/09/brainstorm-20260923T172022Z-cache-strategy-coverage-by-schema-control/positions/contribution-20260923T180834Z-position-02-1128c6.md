# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T180834Z-position-02-1128c6` |
| Slot | `position-02` |
| Actor | model=unknown tool=Codex desktop effort=unknown session=unknown capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:08:34Z` |
| Confidence | medium |

## Summary

Map cache choices across schema control, writer control, placement/update pattern, and freshness contract. Explain why versions and tombstones can help but do not by themselves remove cache fences or every lock.

## Disagreements

Version-only coordination, TTL-only legacy behavior, and the value of a measured soft-delete pair remain open; exact book scope depends on reviewed evidence.

## Contribution

## Evidence examined

- `book/concepts/cache-consistency.typ`, `book/chapters/{major-families,minor-variants,scenarios}.typ`, and `book/evidence/v2/claims.json`: the book states a freshness contract and the broad cache gain, but offers no decision map for schema control, writer control, local versus shared cache, or deletion. Active claims `v2-14`, `v2-15`, and `v2-16` are scoped to the original small single-trial survey; `v2-gap-04` still describes regimes that later Study 05 work has begun to address.
- `studies/05-cache-consistency/README.md` and `sql/README.md`: the study crosses owned/legacy database models, memory/Redis backends, aside/through strategies, relaxed/strict freshness, and coordinated/external writers. Owned adds `person.cache_version`, transactional `cache_outbox`, version changes and ordered locking; legacy forbids cache-specific schema changes. Neither README documents a soft-delete/tombstone design.
- `studies/05-cache-consistency/reports/analyses/20260921-cache-consistency-allgreen--deepseek-flash--2026-09-22.md`: the independent review accepts the scoped cache gain and strict-read result, corrects the owned/legacy attribution of the roughly 87% three-instance wrong-read finding, and warns that Redis was only a cache.
- `studies/05-cache-consistency/reports/discussions/20260921-cache-fences-and-ledger.md` sections 1–3: an old committed fill can repopulate the cache after a newer commit and invalidation; pre- and post-commit fences plus an acknowledgement boundary address this race. Its later residual-failure sections are explicitly superseded and are not current findings.
- `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`, the later signed analyses in `reports/analyses/`, `CONTEXT.md`, and `docs/methodology.md`: v2 adds some TTL, resource and engine evidence, with stated gaps. A book update must resolve any new number through a reviewed active evidence claim rather than silently promoting the current v1 book claim.

## Recommendation

Give the book a cache decision map organized by **four independent choices**, not a single legacy-versus-owned ladder:

| Choice | Questions the reader must answer |
|---|---|
| Source and schema control | Can every relevant mutation atomically advance a version or record an outbox event? Can deletes and re-creates be distinguished? |
| Writer control | Do all writes pass through the caching service, or can external writers bypass its invalidation? |
| Cache placement and update pattern | Process-local or shared cache? Cache-aside or synchronous write-through? These are separate axes; in-memory is a backend, not a third update pattern. |
| Contract | Is an older committed value acceptable after a write acknowledgement? What are the limits on staleness, database validation, outage behavior and refill load? |

Then show the mechanism on one key's timeline: database commit, cache invalidation/fence, fill, publication, and acknowledgement. Define *dirty or impossible publication* separately from *stale committed read*. Mark the point at which a later read must see the new value. This is the clearest way to explain why a valid database version token can coexist with a stale shared-cache publication.

For a frozen legacy schema, describe what coordinated writers can do with post-commit cache fences and per-key fill leases, and state the guarantee that disappears when external writers bypass them. Strictness then needs an authoritative read or a reliable external change signal; that can increase database traffic or coordination cost. For an owned schema, show the transactional version/outbox as stronger change evidence and a way to reject obsolete values, while retaining cross-instance publication fencing and write-order rules. Separate a database lock for contended business writes, a cache fill lease for stampede control, and a version/CAS check; they solve different races. Do not say versioning removes all locks: the owned model itself uses ordered locking, and the measured three-instance owned relaxed cell still returned 87.81% wrong reads in the original scoped run (`v2-15`).

Treat soft delete or tombstone as a **candidate deletion protocol**, not a measured winner. It may preserve deletion identity and a monotonic generation across delete/re-create, but it creates retention, privacy, cleanup and resurrection rules and does not alone prevent a stale fill or an external writer from bypassing invalidation. Ask whether the book needs a conceptual deletion section and whether a future controlled pair is justified; label both as unmeasured by Study 05.

Use a compact trade-off table for each viable regime: database read amplification, write amplification/outbox cost, cache and lease coordination, wrong-read exposure under external writes, outage fallback, eviction/TTL behavior, and operational burden. Attach the active claim IDs only to cells those claims support. Current v2 Study 05 work should enter the book only through a new reviewed evidence version with its remaining hard-expiry, scale, repetition and placement limits explicit.

## Reasoning and trade-offs

The book presently gives a correct warning and a few scoped numbers, but a reader cannot infer which design to choose when the schema is frozen or when external writers exist. The Study 05 matrix already supplies the axes and the fence mechanism. A decision map makes those constraints visible without claiming that cache-aside beats write-through head-to-head; the book itself says that controlled pair is not established. The additional page space and maintenance of a matrix are justified only if every branch names its guarantee and failure mode.

## Assumptions and evidence gaps

The Study 05 domain is one donor portal, not a universal cache policy. The original headline is one small, closed-loop PostgreSQL run on one host; the later v2 analyses extend some regimes but are not yet incorporated into the active book registry. There is no measured soft-delete pair in the inspected packet, and the schema/version discussion does not prove freedom from locks, leases or database reads. Guarantees with external legacy writers depend on a reliable notification or validation mechanism that this study's frozen legacy model explicitly lacks.

## Alternatives

Keep only the current high-level warning: shortest book, but it leaves the user's schema-control question unanswered. Present a large catalog of every caching tactic: more complete-looking, but mixes orthogonal choices and unmeasured claims. A narrower matrix with links to detailed Study 05 mechanisms is preferable.

## Anticipated disagreements

One view may treat version columns as enough to avoid coordination; the stale-fill timeline and owned three-instance cell are concrete counterexamples for the studied contract. Another may prefer TTL-only bounded staleness as a simpler legacy option; that is valid if the contract permits stale reads, but its wrong-read and load profile must be labelled rather than called strict.

## Falsification tests

Try to place every implemented Study 05 cell into the four-axis map without ambiguity; if a cell needs a fifth axis, revise the model. Trace a stale refill, an external write, a delete/re-create, a lease expiry, and a cache outage through each proposed guarantee. Check the final prose against `v2-14`/`v2-15`/`v2-16` and current signed analyses. A controlled soft-delete result or a version-only strict three-instance run with an independent oracle could change the deletion and lock recommendations.

**Confidence:** high that the book lacks this decision map and that versioning alone is insufficient in the measured case; medium on untested soft-delete guidance.
