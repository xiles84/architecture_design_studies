# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T184844Z-critique-01-2f385a` |
| Slot | `critique-01` |
| Actor | model=gpt-5.6-sol tool=Codex effort=high session=01a0cf34-bce0-7162-843a-e5ec65412925 capability=HIGH role=reviewer |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:48:44Z` |
| Confidence | high |

## Summary

Converge on a guarantee-first decision flow while narrowing claims about schema ownership and universal fence necessity. Preserve dissent on how much unmeasured protocol guidance belongs in the book and whether schema control or writer observation should lead.

## Disagreements

How much unmeasured outbox/deletion protocol design belongs in the book; whether schema control or writer-observation completeness should lead the decision map.

## Contribution

# Cross-review — guarantee boundaries before mechanism names

**Brainstorm:** `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`  
**Slot:** `critique-01`  
**Author:** `gpt-5.6-sol`, Codex, high effort, HIGH capability, reviewer  
**Confidence:** high on the evidence and guarantee boundaries; medium on the final book layout

## Evidence inspected

I read the brainstorm question and state, all three independent positions, every path in the evidence packet, the current Study 05 signed analyses and fence discussion, and the relevant harness paths that select strict policies and order fences, commits, publication, and acknowledgement. I treated the active book evidence registry as the authority for book claims and the 2026-09-23 analyses as newer source evidence that is not yet registered. No task or measurement was created.

## Overall assessment

All three positions identify the central stale-fill race correctly and agree on the most important editorial correction: cache-aside versus write-through and process-local versus Redis are implementation dimensions, not freshness guarantees. They also correctly separate database mutation locks, fill leases, and publication fences.

The convergence is real, but the proposed conclusions need two corrections before synthesis:

1. **Study 05 demonstrates sufficiency of one coordinated two-fence protocol; it does not prove a cache-side fence is necessary for every strict design.** A synchronous acknowledgement barrier, authoritative validation on every hit, a transactionally coupled cache, or another protocol could provide an equivalent ordering proof. Positions 01 and 02 occasionally turn “the only strict mechanism implemented here” into “the required mechanism in general.” Position 03 is closer with “fence or proven equivalent barrier.” The book should use that formulation.
2. **The owned and legacy cells do not isolate the value of schema ownership.** The roughly 87% relaxed shared-cache cells show that the owned package did not rescue that particular uncoordinated publication path. They do not establish that ownership “did not buy consistency,” as position 01 states. The owned design bundles a version column, outbox, different SQL shape, and optimistic/pessimistic writers; the relaxed cell deliberately does not claim a strict contract. Study 05 separately shows an owned version token can validate process-local hits across instances. The defensible conclusion is narrower: *a source version that is not consulted at the relevant cache-hit/publication boundary does not make a relaxed shared cache strict*.

## Response to position 01

Position 01 contributes the strongest numerical audit and correctly catches several traps: the private-memory three-instance zeros are validation/bypass outcomes rather than proof that independent LRUs coordinate; hard TTL expiry is still unexercised; outbox consumption and database soft deletes are unmeasured; and active book claims must remain tied to their registered digest.

Its two-axis map is too compressed. “Who arbitrates the copy” combines three unlike things: a cache mutation primitive, a database token, and an authoritative read. Those mechanisms can coexist, and they answer different questions. A database version says which source state a payload represents; a fence/CAS orders publication; an authoritative read establishes currentness at read time. Treating them as alternatives risks teaching readers to choose one when strict designs often need a composition.

The phrase “the correctness axis is the fence, and it is required on both sides” also exceeds the evidence. The harness itself uses three strict policies: shared invalidation fencing, database version validation for multi-instance local caches, and authoritative reads when no cache proof exists. The measured result supports the implemented fence protocol for coordinated shared caches, not a universal necessity theorem.

Position 01 is right that a lease is not a freshness fence. It should go further: a generation fence is not necessarily a lock. Study 05’s Redis Lua operation is atomic coordination, but it does not serialize a critical section in the same way as a database row lock or a fill lease. The book should avoid turning “some coordination remains” into “a lock remains.”

## Response to position 02

Position 02 supplies the most useful four-choice inventory and correctly requires a one-key event timeline. Its distinction between dirty/impossible publication and stale committed reads should be retained. Its soft-delete treatment is appropriately cautious.

The four choices are not fully independent as written. Writer control determines whether a version/outbox is complete; cache placement determines whether invalidation is visible; and the freshness contract determines whether those omissions matter. Presenting them as independent axes invites impossible cells, such as “strict shared cache with bypassing writers and no authoritative validation.” The final map should be a decision sequence with feasibility gates, followed by a trade-off table for viable protocols.

Position 02 also says an owned schema should retain “cross-instance publication fencing.” That is sound for the Study 05 design, but should be written as “atomic conditional publication or an equivalent barrier.” A payload carrying source version 8 can be installed with `put-if-newer` and reject version 7 without a separate fence counter; however, it can still serve version 7 after version 8 commits but before the cache learns of version 8. That example proves both halves: version ordering can replace one stale-overwrite mechanism, but it cannot establish currentness after acknowledgement without validation, synchronous propagation, or a barrier coupled to acknowledgement.

## Response to position 03

Position 03 has the best guarantee boundary: frozen schema does not make strict reads impossible; it removes a cheap cache-hit proof when writes cannot all be observed. Its distinction between database soft delete and cache tombstone, including delete/recreate ABA and tombstone-retention hazards, should survive synthesis. It also correctly insists that asynchronous outbox delivery is not strict after acknowledgement unless acknowledgement waits for the relevant cache effect or reads validate against authority.

Its proposed source/writer-by-cache table still merges evidence levels inside cells. For example, the shared frozen-schema row combines a mechanism demonstrated on one host with a general statement about all instances participating. The latter is a protocol precondition, not a measured comparative result. Each row needs separate fields for `guarantee`, `required assumptions`, `Study 05 evidence`, and `unmeasured failure boundary`.

Position 03’s statement that one process has one coordination domain also needs a qualification: it is true only while there is one cache-owning process and no parallel writer outside that process. Restarts, background jobs, direct database writes, and rolling deployment can create a second domain without changing the nominal architecture. The decision map should ask about *all mutation paths over the guarantee interval*, not only the steady-state instance count.

## Proposed convergence

Organize the book section as a decision sequence rather than a flat taxonomy or orthogonal cube:

1. **Name the contract.** Choose strict after acknowledgement, bounded staleness, session/read-your-writes, or best effort. Define impossible/uncommitted values as forbidden separately from permitted stale committed values.
2. **Establish change-observation completeness.** Can every mutation, including bulk jobs, migrations, deletes, recreates, and external writers, participate in the protocol? If no, a strict cache hit requires authoritative validation or a complete external change source; otherwise downgrade the contract explicitly.
3. **Identify source capabilities.** A frozen schema can still support service-side generations and authoritative reads. An owned schema can add a monotonic source version, transactional outbox, and deletion generation. These capabilities improve proof and recovery options; they do not confer freshness by themselves.
4. **Choose topology and publication protocol.** Process-local versus shared decides the coordination domain. Cache-aside versus write-through decides who fills or republishes. For strictness, state the exact ordering proof: two-fence invalidation, version validation on hit, conditional `put-if-newer` plus synchronous acknowledgement barrier, authoritative read, or another demonstrated equivalent.
5. **Add performance controls separately.** A bounded lease suppresses duplicate fills; TTL and early expiry bound residency/load behavior; neither supplies currentness. Report database reads per logical read, write and acknowledgement latency, cache operations, wrong-read exposure, recovery lag, and operational dependencies.

A compact table should then use these columns:

| Protocol | Contract | Complete writer observation required? | DB reads on hit | Write/ack cost | Shared coordination | Failure behavior | Evidence status |
|---|---|---:|---:|---|---|---|---|

The evidence status must distinguish:

- **Active registered measurement:** `v2-14`, `v2-15`, and `v2-16`, with their original small, single-trial, closed-loop, `db-only` scope.
- **Newer signed source evidence, not yet an active book claim:** TTL-duration churn, equal-total framing, YugabyteDB/endpoints, and intended placement pairs, retaining the hard-expiry, scale, repetition, and placement-proof gaps.
- **Mechanism demonstrated:** the stale-fill timeline, before/after-commit fence protocol, acknowledgement boundary, version validation, authoritative bypass, and lease-holder recovery within the study’s tests.
- **Proposed/unmeasured:** outbox consumer semantics, soft-delete protocol, delete/recreate retention, CDC, high-hit partitioned local caches, cache eviction pressure, and isolated comparisons among version-check-on-hit, synchronous fences, and asynchronous outbox.

## Lock-avoidance wording

The synthesis should reject the sentence “versions or tombstones remove locks” and also avoid the mirror overstatement “strict caches require locks.” Use named mechanisms:

- A source version plus atomic conditional write can reject an older publication without a fill lock.
- It cannot prove that the cache knows the latest acknowledged source version.
- A fill lease reduces duplicate database work and needs token-safe expiry/release; it is not a freshness proof.
- A database row lock or optimistic CAS protects a business mutation; Study 05 does not isolate a general winner for that choice.
- A cache generation/fence orders fill publication against invalidation. It is shared coordination, but need not be described as a lock.
- An asynchronous outbox provides durable change identity and replay. Without an acknowledgement barrier or read validation, delivery lag still permits stale-after-ack hits.
- A database soft delete preserves identity only if its generation remains monotonic across delete/recreate. A cache tombstone prevents resurrection only for a stated retention/in-flight interval or with an authoritative fallback.

## Material dissent to preserve

There are two reasonable book-scope choices that the evidence does not settle:

1. **How much unmeasured protocol design belongs in the book.** One approach includes a conceptual deletion/outbox section clearly labelled as inference; the conservative approach records these only as gaps until controlled evidence exists. I recommend a short conceptual section because deletion and event lag are correctness hazards readers need to recognize, with no comparative recommendation or numeric claim.
2. **Whether to lead with schema control or writer observation.** Schema control is the user-facing constraint, but writer observation is the actual guarantee gate. I recommend leading the decision flow with the contract and writer completeness, then branching on schema control. A legacy system with one fully controlled adapter can support a stronger cache protocol than an “owned” database with sanctioned bypass paths.

## Assumptions and falsification

This critique assumes Study 05’s strict contract: a read starting after acknowledgement of a write to the same key must not return an older committed value; overlapping reads are classified separately; impossible or uncommitted values are forbidden in all arms.

The proposed convergence should change if a reviewed controlled pair shows any of the following:

- source-version `put-if-newer` with no synchronous invalidation or validation remains strict across a forced commit-to-notification delay;
- an async outbox preserves strict-after-ack semantics while acknowledgement does not wait and reads do not validate;
- a deletion protocol survives tombstone expiry, eviction, key reuse, and delayed old fillers without a monotonic retained generation or authoritative check;
- a high-hit multi-instance process-local design stays strict without per-hit database validation and without complete write routing;
- repeated, externally offered-load runs materially change the shared relaxed-cache failure mechanism or the conclusion that strict read-path overhead was below the original run’s noise floor.

## Recommendation to the synthesizer

Keep the positions’ shared race analysis and evidence cautions, but synthesize a **guarantee-first decision flow**. State that Study 05 demonstrates a sufficient coordinated fence protocol and several fallback policies; it does not prove one universal cache mechanism. Narrow the owned/legacy conclusion to what the controlled evidence supports, distinguish atomic coordination from locks, and preserve the two scope disagreements above.
