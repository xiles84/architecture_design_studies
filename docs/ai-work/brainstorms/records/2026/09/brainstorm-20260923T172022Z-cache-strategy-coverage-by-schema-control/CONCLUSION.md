# Brainstorm conclusion

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T233727Z-synthesis-01-3bd104` |
| Slot | `synthesis-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=session-20260923T174300Z-deepseek-flash capability=HIGH role=integrator |
| Capability input | HIGH |
| Submitted | `2026-09-23T23:37:27Z` |
| Confidence | high |

## Summary

Adopt a guarantee-first cache decision flow with a protocol table carrying an explicit evidence-status column: name the contract, establish writer-observation completeness, identify source capabilities, then choose topology and publication protocol, with performance controls separate. Teach the studied mechanisms, keep every number on an active claim, and route the newer Study 05 results through a versioned evidence update.

## Disagreements

Preserved: how much unmeasured deletion/outbox protocol design belongs in the book; whether cache-side fences are uniquely necessary (resolved as fence-or-equivalent, recorded); which unmeasured tactics warrant controlled pairs.

## Contribution

## Decision

Adopt a **guarantee-first cache decision flow** and a protocol table with an explicit
evidence-status column. Teach the mechanisms Study 05 demonstrates, keep every number tied to an
active registered claim, and route the newer 2026-09-23 Study 05 results through a versioned evidence
update before the book quotes them.

The flow, in order:

1. **Name the contract.** Strict after acknowledgement, bounded staleness, session/read-your-writes,
   or best effort. Define impossible/uncommitted cache values as forbidden separately from permitted
   stale committed values.
2. **Establish change-observation completeness.** Can every mutation — including bulk jobs,
   migrations, deletes, recreates and external writers — participate in the protocol? If not, a strict
   *cache-hit* proof is unavailable; use an authoritative read or a complete change source, or
   downgrade the contract explicitly.
3. **Identify source capabilities.** A frozen schema still supports service-side generations and
   authoritative reads; an owned schema can add a monotonic version, a transactional outbox and a
   deletion generation. These improve proof and recovery options; they do not confer freshness.
4. **Choose topology and publication protocol.** Process-local versus shared decides the coordination
   domain; cache-aside versus write-through decides who fills or republishes. State the exact ordering
   proof: two-fence invalidation, version validation on hit, conditional put-if-newer plus a
   synchronous acknowledgement barrier, authoritative read, or another demonstrated equivalent.
5. **Add performance controls separately.** A bounded lease suppresses duplicate fills; TTL and early
   expiry bound residency and load. Neither supplies currentness.

## Why

Against the four decision criteria:

- **Map book coverage to Study 05.** The book (`concepts/cache-consistency.typ`, family 6 of
  `major-families.typ`, `scenarios.typ`, `minor-variants.typ`) states the gain, the failure mode and a
  freshness warning but has no selection map, so a reader cannot choose when the schema is frozen or
  when writes bypass the adapter.
- **Separate measured from untested.** Active claims support: an add-cache gain of 2.2–3.5x and p99
  15.8–27.1 ms → 0.8–2.0 ms under `db-only` framing (`v2-14`); a multi-instance mechanism in which
  three logical instances on one shared Redis cache, write-through + relaxed, recorded 87.81% wrong
  reads on the owned model (`38749/44126`) and 87.24% on legacy (`35287/40449`) while every
  single-instance relaxed cell recorded zero (`v2-15`); and a strict/relaxed read-throughput
  difference of −14.2%..+12.3%, inside the run's ~20% noise, with the cost on the write/coordination
  path (`v2-16`). The mechanism record is the stale-fill race: a reader snapshots committed S0, a
  writer commits S1 and invalidates, and the reader republishes S0 — so "publish after commit" is
  necessary and not sufficient; a strict writer fences before *and* after commit, and the contract
  boundary is the **acknowledgement**, not the commit
  (`reports/discussions/20260921-cache-fences-and-ledger.md` §§1–3).
- **Compare correctness, cost, coordination, failures, complexity.** Achieved by naming the
  mechanisms apart — database mutation lock, fill lease, atomic cache publication fence/CAS, source
  version token, outbox — because they solve different races and were not isolated against each other.
- **Test whether versioning and soft deletes remove locks.** They do not: the owned relaxed
  three-instance cell still failed, and neither outbox consumer semantics nor a database soft-delete
  design was measured.

**Two corrections carried into this conclusion, both verified.**

1. The three-instance **private-memory** zeros are a *bypass* result, not safety. In
   `reports/20260921T-survey3.md` the memory `instances_3` arms show cache hits 0 with 15k–21k bypass
   reads and the note "separate in-process LRUs and separate leases; a write through one instance is
   invisible to the other two". They measure a cache that is effectively absent, and must never be
   cited as evidence that a process-local cache is consistent.
2. The registry is **stale relative to committed runs**. `book/evidence/v2/claims.json` was last
   committed 2026-09-22, and `v2-gap-04` still states that real-TTL churn, repeated trials and cache
   resource accounting are unmeasured; `book/chapters/scenarios.typ` repeats that in prose. The
   2026-09-23 churn and equal-total runs measured those (hard expiry still fired zero times; medium
   scale, placement proof and open-loop demand remain open). Newer results may be used as *mechanism*
   evidence now, but must not be quoted as numbers under the existing v2 claim IDs.

## Rejected alternatives

- **Keep the current flat pattern list.** Shortest, but leaves the schema-control question unanswered
  and hides that the owned package did not change the relaxed shared-cache failure rate.
- **Lead with cache-aside versus write-through.** They are fill/publication paths, not guarantees, and
  Study 05 has no head-to-head ratio for them; the 600 s churn cells agree closely (197,196 vs 182,282
  writes, 0 wrong reads each).
- **Present version, outbox or soft delete as the lock-removal answer.** Unsupported; the owned
  relaxed cell still failed, and both the outbox consumer protocol and a database soft-delete design
  are unmeasured here.
- **Promote the 2026-09-23 numbers into the book immediately under v2 IDs.** Violates the registry
  rule; the correct route is a versioned, semantically checked evidence update.

## Remaining dissent

Preserved material objections:

1. **How much unmeasured protocol design belongs in the book.** One view includes a short, clearly
   labelled conceptual branch for deletion and outbox (readers need to recognize the hazards); the
   conservative view records them only as gaps until controlled evidence exists. Not settled; the
   owner decides at task time.
2. **Whether cache-side fences are uniquely necessary.** Resolved in debate as "a fence or a
   demonstrated equivalent barrier" — the harness itself uses at least three strict policies (shared
   invalidation fencing; a database version check for multi-instance local caches, owned; authoritative
   reads, legacy). Study 05 demonstrates a *sufficient* protocol for coordinated shared caches, not a
   universal necessity. The dissent is recorded so a later reader does not read necessity into it.
3. **Which unmeasured tactics warrant controlled pairs.** Named candidates, not chosen here:
   version-validation-on-hit versus synchronous shared fence versus outbox-delivered invalidation
   under equal-total resources; soft delete/tombstone delete-and-recreate with a paused old filler;
   a high-hit partitioned process-local cache; and a crash between commit and cache effect. Selecting
   these is task work, not debate work.

## Assumptions, gaps, and transfer risks

- The assumed contract is Study 05's: a read starting after acknowledgement of a write to the same key
  must not return an older committed value; overlapping reads are classified separately; impossible or
  uncommitted values are forbidden in every arm. A bounded-staleness product may legitimately choose a
  simpler path — but must publish an explicit stale-read budget and must not be called strict.
- All evidence is one small, single-host, closed-loop PostgreSQL + Redis arrangement; later arms add
  YugabyteDB and an equal-total label but are single runs. No cross-service generalization is
  supported.
- Redis is measured as a shared cache, never as an authoritative store; one host is not placement
  evidence.
- Every protocol-table cell must carry one of: active registered measurement (with its original
  scope), newer signed but unregistered source evidence, mechanism demonstrated, or proposed/unmeasured.

## Falsification and review triggers

Reopen or supersede this conclusion if a reviewed controlled pair shows any of:

1. put-if-newer with no synchronous invalidation or read validation stays strict across a forced
   commit-to-notification delay;
2. an asynchronous outbox preserves strict-after-ack semantics with no acknowledgement barrier and no
   read validation;
3. a deletion protocol survives tombstone expiry, eviction, key reuse and a delayed old filler without
   a monotonic retained generation or authoritative check;
4. a high-hit, multi-instance process-local design stays strict without per-hit database validation
   and without complete write routing; or
5. repeated open-loop runs materially change the shared relaxed-cache failure mechanism or the
   strict-read overhead result.

## Task-ready next steps

Not authorized by this conclusion; the owner must separately issue
`CREATE TASKS FROM BRAINSTORM brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`.

1. Produce a versioned evidence update for Study 05 v2: ingest the committed churn, equal-total and
   engine/endpoint runs, supersede the stale parts of `v2-gap-04`, and correct the book gap sentence by
   supersession — never by silently editing the earlier card.
2. Rework the book cache material into the guarantee-first decision flow plus the labelled protocol
   table, with one stale-fill sequence and one delete/recreate extension, and every lock-avoidance
   sentence naming the specific race it addresses.
3. Only if comparative guidance is needed beyond current evidence, run the candidate controlled pairs
   named under remaining dissent 3.
4. Cross-brainstorm handoff: if the figure pipeline from
   `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` is built, the stale-fill
   sequence is the natural first "mechanism demonstrated" figure — it must cite registered claims only.

**Synthesizer disclosure.** This synthesis was written by the session that authored `position-01` in
this record (`deepseek-flash`, HIGH, `session-20260923T174300Z-deepseek-flash`). Both critiques
narrowed two claims in that position — the necessity of a cache-side fence, and the reading that
ownership "did not buy consistency" — and both corrections are carried above. In addition, all three
positions and both critiques were produced by two model families (`deepseek-flash` and Codex/GPT-6),
with all cross-review from one vendor tool; the owner should weigh that reduced independence.
