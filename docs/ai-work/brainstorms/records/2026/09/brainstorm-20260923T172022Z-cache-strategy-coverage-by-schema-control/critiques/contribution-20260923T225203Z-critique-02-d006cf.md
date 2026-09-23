# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T225203Z-critique-02-d006cf` |
| Slot | `critique-02` |
| Actor | model=GPT-6 family (exact variant not exposed) tool=Codex desktop effort=unknown session=codex-root-20260923 capability=HIGH role=reviewer |
| Capability input | HIGH |
| Submitted | `2026-09-23T22:52:03Z` |
| Confidence | medium |

## Summary

Converged on a guarantee-first cache decision map, separating schema control, writer reachability, topology, and update path; challenged claims that versions, outboxes, or soft deletes remove all coordination.

## Disagreements

Whether cache-side fences are uniquely necessary, the compactness of the book map, and which unmeasured tactics warrant controlled studies remain open.

## Contribution

# Cross-review — cache strategy coverage by schema control

Actor: GPT-6 family, Codex desktop; exact model variant and effort are unavailable to this session. HIGH reviewer. I reviewed DeepSeek Flash `position-01`, Codex `position-02` (its archive actor model is `unknown`), and GPT-6-family `position-03` independently of critique-01.

## Evidence checked

I checked `studies/05-cache-consistency/reports/20260921T-survey3.md`, `reports/discussions/20260921-cache-fences-and-ledger.md` §§1–3, `sql/README.md`, and `book/concepts/cache-consistency.typ` against the positions. The report assigns the three-instance relaxed shared-Redis write-through wrong-read counts to legacy 35,287/40,449 (87.24%) and owned 38,749/44,126 (87.81%); it reports zero stale-after-ack reads in strict cells. The fence discussion requires a per-key generation captured before the DB snapshot and strict writer fences before and after commit. The owned SQL catalogue includes version, outbox, optimistic CAS, and ordered pessimistic locking. Its existence is not a controlled isolation of version or outbox benefit. The current book has scoped claim cards but no selection map by schema control and writer reachability.

## Challenge and convergence

DeepSeek's two-axis framing captures an essential insight: schema ownership changes available mechanisms, not the freshness contract. But `cache-side fence | DB version token | authoritative read` are not mutually exclusive choices. The measured owned path combines a DB version with cache-side coordination, and validating each hit against authority can coexist with either. The decision map should start with the contract and *complete observation of every writer*, then show available coordination and placement choices. Schema ownership, writer control, and private/shared topology are independent inputs; cache-aside/write-through are update paths under those inputs.

Position-02 and position-03 correctly resist saying that versioning or soft deletes automatically remove locks. I would sharpen their claim: a version lets a reader or publisher compare states it knows about, but it cannot detect an unseen committed version unless a complete change signal arrives or the hit checks DB authority. A transactional outbox makes an event durable, yet an asynchronous consumer introduces a post-commit window. A writer may have to wait for invalidation before acknowledging, or readers must validate/fail closed for the strict after-ack contract. A DB mutation lock, a fill lease, and an atomic cache publication fence address distinct conflicts; do not put them in one 'locks required' column.

DeepSeek is right that the single-run 87% cells cannot become a universal failure rate and that the three-instance private-memory zero wrong reads do not establish a high-hit, safe private cache. The report explicitly notes separate LRUs/leases, invisible cross-instance writes, and 15k–21k bypass reads in many such cells. Conversely, the 87% value is measured relative to the study's strict comparator, while a relaxed product contract can permit older committed values. The book should name the comparator and show both latency/throughput and stale-read exposure, never present a permitted relaxed read as an impossible or dirty read.

I question DeepSeek's statement that a cache-side fence is the *required* strict method, and its claim that ownership buys fewer moving parts. Study 05 demonstrates one fence-based strict implementation for coordinated writers; it does not exclude another complete barrier or isolate operational complexity. Position-03 gives the better guarantee boundary: with external writers invisible to the adapter, a frozen-schema service cannot promise strict *cache-hit* freshness without complete change observation. DB reads can still give strict answers. A change feed, trigger, durable event or hit validation is a candidate only when its completeness and acknowledgement ordering are demonstrated. The book should distinguish this from the measured external-writer bypass result.

Soft delete is worth a short conceptual branch because it addresses the user's specific schema-control idea. Treat it as a way to retain deletion identity/version across delete and recreate, subject to retention, cleanup and key-reuse rules. It is not an invalidation protocol. Neither a DB soft-delete pair nor outbox-consumer strictness was measured here; mark both proposed and specify the race each would need to solve. The cache-side tombstone in fault injection is a different mechanism. A future experiment should deliberately pause an old filler across delete/recreate and tombstone expiry, then audit resurrection with an independent ledger.

The book's active v2 claims and the later signed Study 05 analyses need a provenance gate. The newer churn runs cross real TTL intervals but report no hard TTL expiries, so neither 'TTL untested at all' nor 'hard expiry established' is accurate. Do not silently use later numbers under v2 claim IDs. First produce a semantically checked, versioned evidence update; the present book can state mechanisms and labelled gaps without promoting those measurements.

## Proposed next steps for agreement

1. Define each freshness promise relative to write acknowledgement, and separately prohibit dirty/impossible values. State which writers can bypass the cache adapter.
2. Build a compact decision table: schema control × writer coverage × private/shared cache, with cache-aside/write-through and relaxed/strict guarantees inside applicable cells. Each row states its DB-read, write/coordination, failure and stale-read costs.
3. Draw one stale-fill sequence and one delete/recreate extension. Show where a version, outbox, pre/post fence, lease and DB validation do or do not close the window. Label observed Study 05 behavior separately from proposed mechanisms.
4. Audit every numeric and causal statement against the active book claim registry and signed source evidence. Version any later evidence correction; retain study gaps for hard TTL, scale, placement proof and repetition as applicable.
5. Only if the book needs comparative guidance beyond the current evidence, design controlled pairs for version validation, synchronous shared fence, outbox-delivered invalidation and soft delete, under coordinated/external writers with crash and reorder injections. Do not invent a winner before those runs.

## Dissent, assumptions and falsification

I retain dissent from position-01 that cache-side fences are uniquely required and from any reading of position-02 that a four-axis matrix should enumerate every Cartesian product. The map must fit a reader's decision and can mark infeasible combinations. I assume the study's strict contract is the intended reference; a bounded-staleness product may select a simpler path with an explicit stale-read budget. A high-hit process-local strict run, or an outbox/DB-token protocol that passes post-ack, crash, reorder and external-writer tests without a cache-side fence, would change the recommendations. Confidence: high on the measured attribution and guarantee boundary; medium on book organization and unmeasured tactic costs.
