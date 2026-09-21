---
discussion_id: 20260921-cache-fences-and-ledger
study: 05-cache-consistency
run_id: 20260921T-survey
inputs_digest: 8ff86dc9e5228e86
repo_commit: 567778ff7091f7748a3a6096cec3b98b408adec7
author: deepseek-flash
author_kind: ai
author_version: deepseek-flash, DeepSeek HIGH role (model id, exact version, tool and effort not exposed by the client), WSL2/Podman
written_at: 2026-09-21T15:35:00Z
---

# Study 05 — the fence, the ledger and the failures: mechanisms behind the numbers

This is the companion to
[`analyses/20260921-cache-consistency-survey.md`](../analyses/20260921-cache-consistency-survey.md).
It records the mechanisms, the race timelines and the harness decisions that the
analysis states as results, so a later analyst can attack them directly.

## 1. Why "publish only after the commit" is not enough

The protocol's conservative floor is that a cache value may be published only if
it corresponds to a successfully committed database state. That floor is
necessary and it is not sufficient. Timeline (all events real; the counts are
from the dev-check iterations recorded in `PROGRESS.md`):

```
reader R                          writer W
----------                        ----------
                                  (nothing yet)
Get(key) -> miss
lease acquired
readPortalSnapshot -> state S0
                                  BEGIN; UPDATE; COMMIT   (S1 now committed)
                                  DELETE key              (invalidation)
Put(key, S0)   <-- ACCEPTED
                                  acknowledge to caller
reader R2 starts, requires S1
Get(key) -> HIT, returns S0  -->  STALE AFTER ACK
```

R's snapshot predates the commit, so R's value is committed *and* superseded.
Nothing in "publish after commit" prevents R from putting it back after W removed
it. The first implementation of this study did exactly this and recorded ~14 000
stale-after-ack reads in a strict cell under coordinated writers.

### The fix in this study

A per-key **invalidation fence** in the cache itself:

* `Fence(key)` advances the key's fence **and** removes the value, atomically
  (one mutex in the in-process store, one Lua script in Redis);
* `FenceOf(key)` reads it;
* a fill captures the fence **before** it takes its database snapshot and stamps
  it on the entry;
* `Put` is refused unless the stored fence still equals the entry's.

R's `Put` is now refused, R still serves its own (fresh, for its own requirement)
value, and R2's later fill reads a post-commit state. The run shows 61–85 refused
publications per cell, so the race is not theoretical.

### Why a strict writer must fence twice

Fencing only before the mutation leaves this window:

```
W: Fence (F -> F+1), value removed
R: FenceOf -> F+1
R: readPortalSnapshot -> S0        (W has not committed yet)
W: COMMIT (S1)
                                  <- R may still publish S0 at F+1: accepted
W: Fence (F+1 -> F+2)              <- too late for R
```

R captured *the new fence* and read *the old state*, so the fence check passes.
The post-commit fence is what makes the token sound: after it, `FenceOf` can only
return a value that postdates the commit. A strict writer therefore fences before
**and** after the commit. Removing the second fence re-introduced thousands of
stale reads in the dev checks.

### Why "acknowledged" and not "committed" is the right boundary

The contract is "a read that begins after a write to that key was
**acknowledged** must not return an older value". A strict writer leaves the
cache invalidated *before* it acknowledges, so the ack is the earliest moment at
which a served value can be judged against the write. The ledger therefore:

* applies the **data** change at commit time (`markPendingLocked`);
* moves the **requirement** at the acknowledgement (`Ack`), which every write
  path calls where it returns success to the caller — after the fences for a
  strict writer, after the best-effort update for a relaxed one, immediately for
  an external writer whose commit is its own acknowledgement;
* classifies a value matching a committed-but-unacknowledged state as `ahead`,
  never `impossible`.

Treating the commit as the boundary made correct cells look wrong: a read served
the superseded value while the writer was still fencing, and the accounting
called it a violation. That is the same class of error the repository's
`LESSONS_LEARNED.md` already describes from study 01 — the harness manufacturing
the failure it was looking for.

## 2. The ledger's incomplete model of concurrent writes

The oracle models each key's committed history. Two things break it under
deliberately conflicting concurrent writes, and both are *harness* limitations
rather than design findings:

1. **Ordering.** `noteCommitted` applies a mutation as soon as its transaction
   returns. Two writers committing the same key can have their transactions
   commit in one order and their `noteCommitted` calls run in another. For
   set-valued mutations (`person_update`, `insert`, `delete`) the resulting
   ledger state can then be one the database never held — which surfaces either
   as an `impossible` value or as an audit mismatch. The observed audit failure
   (a donation present in the ledger and absent from the database) has this
   shape.
2. **Aggregates.** An absolute amount correction is order-sensitive for the same
   reason. This one *was* fixed properly rather than argued away: corrections are
   now relative in both the SQL and the ledger
   (`amount_cents = amount_cents + $2`), so two concurrent corrections produce
   the same total whichever order they commit in. A 795-cent aggregate mismatch
   disappeared with that change.

The unfixed part is the ordering of set-valued mutations. Two honest routes
exist: serialise conflicting writes per key inside the harness (which changes the
contention the hotspot phase exists to measure), or restrict the write mix to
order-insensitive mutations. **Neither was done, and the analysis declares the
consequence rather than hiding it.**

## 3. The pending-state race

The oracle learns of a committed state just after the commit returns. Before
`Ack`, that state lives in a pending slot. With a *single* slot, two writers
committing before either acknowledges lose one of the two states from the
history entirely; a reader that then serves it from the cache is told it
corresponds to no committed state. That produced 3–131 `impossible` values per
cell. Pending states are now a reference-counted set and acknowledgements carry
the exact state they acknowledge.

The complementary rule: a value that arrived from a **database read** cannot be
impossible — the read runs in one repeatable-read transaction, so what it
returned is a state the database held, and the only reason the ledger may not
know it is its own recording lag. Such a value is classified `ahead`. A value
that arrived from a **cache hit** is exactly the case the detector exists to
catch, and an unknown hash on a hit is still impossible. This asymmetry is what
lets the dirty-cache-write audit keep its teeth.

## 4. The controlled comparisons, and why some of them are not reportable

The study's controlled pairs change exactly one dimension:

| Pair | Single decision isolated | Reportable? |
|---|---|---|
| `{legacy,owned}-na-none-none-none-coord` vs `{legacy,owned}-opt-none-none-none-coord` | the version column and the outbox | yes |
| `{...}-memory-aside-{relaxed,strict}-coord` | the freshness contract | yes |
| `{...}-redis-aside-relaxed-coord` vs `{...}-memory-aside-relaxed-coord` | where the cache lives | yes, with the caveat that Redis's LRU is approximate and its capacity is not the in-process capacity |
| `{...}-redis-aside-relaxed-{coord,ext20}` | coordinated vs external writers | partly: the external cell's strict sibling failed |
| `legacy-na-redis-aside-strict-ext20` vs `legacy-na-redis-aside-strict-coord` | external writers under strict freshness | **no** — the external cell failed its contract |
| `{legacy,owned}-opt-...-coord` vs the `{legacy,memory}` three-instance arms | one instance versus three | yes, and it is the run's headline |
| `{owned-opt,owned-pess}-...` | optimistic versus pessimistic concurrency | yes, both passed their gates |

## 5. What the harness does NOT claim

* It does not claim that a local in-memory cache is strict in a multi-instance
  deployment. Its strict policy for that case is a database version check (owned)
  or an authoritative read (legacy), and the result records which.
* It does not claim that containers on one laptop demonstrate colocation or
  network behaviour. The placement pair (`ref-y1-colocated`,
  `ref-y2-noncolocated`) exists as SQL and was **never executed**.
* It does not claim that the three-instance result generalises. It is one trial
  at one scale on one machine.
* It never reports a stale or impossible result as fast. The generated report
  prints a relaxed cell's wrong-read count and rate beside every throughput
  figure, and the analysis does the same.

## 6. Open leads, in the order I would take them

1. The strict **through** strategy fails in both models and both backends while
   the strict **aside** strategy does not. The asymmetry points at the
   post-commit republish sharing a store with concurrent fillers. Every failing
   cell's causes are in its JSON; the shortest path is
   `owned-opt-redis-through-strict-coord` (35 stale reads).
2. `ref-rollup-trigger`'s 76 impossible values need a drift diagnosis before the
   rollup requirement can be called covered at all.
3. The three-instance 85 % wrong-read rate needs a second machine and repeated
   trials before it is treated as a property of shared caches rather than of this
   arrangement.

Analyst: `deepseek-flash`, acting as DeepSeek HIGH. No independent review.
