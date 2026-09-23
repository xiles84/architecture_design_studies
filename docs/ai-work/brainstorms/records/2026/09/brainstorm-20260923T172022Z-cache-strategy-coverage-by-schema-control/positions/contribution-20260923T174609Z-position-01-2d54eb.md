# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T174609Z-position-01-2d54eb` |
| Slot | `position-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=session-20260923T174300Z-deepseek-flash capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-23T17:46:09Z` |
| Confidence | high |

## Summary

Map caches on two axes: who arbitrates the copy (cache-side fence, DB version token, authoritative read) and whether the schema may change (legacy vs owned), with aside/through, memory/redis and relaxed/strict as parameters. Ownership did not buy consistency (owned 87.81% vs legacy 87.24% wrong reads under relaxed write-through), and the strict arms came from a fence that fires before and after commit. Record soft deletes, outbox-driven invalidation, partitioned process-local caches, hard-TTL expiry and placement as gaps.

## Disagreements

Anticipated: relaxed staleness is permitted so it should not be called wrong; one machine and one trial cannot be generalized; version columns, outboxes and soft deletes are standard and may be recommended.

## Contribution

## Position

Restructure the cache material as **two axes and a set of parameters**, not a flat list of patterns:

- **Axis A — who arbitrates the copy:** cache-side fence | database version token | authoritative read.
- **Axis B — may the schema change:** legacy (frozen) | owned (version column, outbox, trigger).

`aside`/`through` (how a value is filled/published), `memory`/`redis` (where the copy lives) and
`relaxed`/`strict` (what freshness is promised) are **parameters on those axes, not the axes
themselves**. Under that map the book's existing sentences survive, but three of them need
qualification: ownership does not buy consistency, a lease is not a fence, and the lock-avoidance
result from Study 04 does not transfer to caches. Then record the untested tactics — process-local
caching with real partitioning, database soft deletes/tombstones, outbox-driven invalidation — as
gaps, not as advice.

### What Study 05 actually supports (measured)

From `studies/05-cache-consistency/reports/20260921T-survey3.md` (27 cells, 0 failed, both controls
fired) and the two 2026-09-23 analyses):

- **The gain:** 2.2–3.5x throughput, p99 15.8–27.1 ms → 0.8–2.0 ms, `resource_framing: db-only`
  (`v2-14`). Redis is measured as a shared cache, never as an authoritative store.
- **The failure:** three logical instances on a shared Redis cache, write-through + **relaxed**,
  recorded **87.81%** wrong reads on the owned model (`38749/44126`, `o:opt:red:wth:rel`) and
  **87.24%** on legacy (`35287/40449`); every single-instance relaxed cell recorded zero
  (`v2-15`). The two models are within 0.6 points of each other.
- **The repair:** the strict arms of the same shared-cache topology recorded **zero** wrong reads
  (`o:opt:red:wth:str`: 58 807 hits, 0 wrong; `o:opt:red:asd:str`: 58 129 hits, 0 wrong).
- **Freshness cost:** every controlled strict/relaxed pair differed by −14.2%..+12.3% in *read*
  throughput, inside the run's ~20% single-trial noise; the visible cost is on the write path
  (`v2-16`).
- **Churn:** 600 s and 1800 s cells crossed 2.00 and 6.00 real 300 s TTL boundaries with 0 wrong
  reads over 1.37M–2.88M reads; but `expired_hard = 0` — only probabilistic early expiry fired
  (11 321–37 063). "The TTL was crossed in time" is measured; "the hard TTL expired an entry" is
  not.
- **Faults:** cache-unavailable → 20 reads bypassed, 0 wrong; lease-holder death → the abandoned
  lease was stealable after 100 ms with no manual intervention.

### Axis B changes which mechanisms are *available*, not the required guarantee

The invariant is the same for both models: *a cache value may be published only from a committed
state, and invalidation must refuse a stale publication, not merely delete it*. What differs:

- **Legacy (frozen):** there is no version column, no outbox, no trigger and no reliable
  `updated_at`, so the fence has to live in the cache/service (Study 05's Redis Lua fence + lease).
- **Owned:** you may add `person.cache_version` and `cache_outbox`, which improves optimistic
  concurrency and gives a durable post-commit invalidation hook — but it did **not** remove the
  multi-instance staleness: owned 87.81% vs legacy 87.24%.

So the book's "owned vs legacy" distinction is real and useful, but it is an *availability* axis,
not a correctness axis. The correctness axis is the fence, and it is required on both sides.

### Which "avoiding locks" claims need qualification

1. **`v2-12` (lock vs version, 1.76x) is a database-row claim and an upper bound.** It is
   PostgreSQL, small scale, one hot key, 16 writers, and the run could not separate the lock's
   advantage from extra per-retry work. It must not be read as "caches make locks unnecessary".
   `how-to-choose.typ` already uses it correctly for contended rows; the cache chapter must not
   inherit it.
2. **A version token does not remove the window by itself — this is measured, not inferred.** The
   owned model (version + outbox) still recorded ~87% stale under relaxed write-through. Any
   sentence of the form "add a version column and you can drop the fence/lease" is contradicted by
   that cell.
3. **A lease is coordination, not a fence.** It suppresses duplicate fills and survives holder
   death (stealable after 100 ms), but it does not order "the database snapshot I read" against
   "the publish I perform". The gate-passing strict designs needed the fence.
4. **The fence has a timing requirement.** A fill can capture the *new* fence token and read the
   *pre-commit* snapshot, then publish it, because the post-commit fence has not happened yet — so a
   strict writer must fence **before and after** commit (study 05 discussion §1). "Publish only
   after commit" is necessary and not sufficient.
5. **The contract is on the acknowledgement, not the commit.** A read that begins after a write was
   *acknowledged* must not return an older value; treating commit as the boundary made correct
   cells look wrong (`LESSONS_LEARNED.md`, study 01 harness-manufactures-failure lesson).
6. **The locks that do exist for caches in the corpus are database locks**, in the owned model's
   pessimistic path (`w_lock_person`, `w_lock_persons2`, ordered for two-parent reassignment). They
   order the mutation; they do not provide cache freshness.

The honest summary: **the measured way to keep a shared cache strict is a cache-side fence; no
measured design removed both lock and fence.**

### `aside` vs `through`, and `memory` vs `redis`, are not consistency choices

- At 600 s churn with the strict owned model, through and aside were both correct: 197 196 vs
  182 282 writes, 12 517 vs 11 321 probabilistic expiries, 0 wrong reads each. Study 05 has **no**
  head-to-head throughput pair isolating aside from through, so the book must not imply a ratio; the
  current `minor-variants.typ` sentence ("not head-to-head") is right.
- `memory` vs `redis` is *where the copy lives*, and it changes what the failure looks like. In the
  one-vs-three table the 3-instance **memory** arm has separate in-process LRUs and separate leases:
  cache hits collapse to ~0 (the report records hits 0 and 15k–21k bypass reads) and wrong reads are
  0 **because the reads bypassed the cache**. That is a no-hit/bypass result, not evidence that a
  process-local cache is consistent — the book must say so explicitly, because the far more
  dangerous reading ("process-local is safe, Redis is not") is exactly what those zeros invite.
- The 87% failure is therefore best described as: *a shared cache, a relaxed contract and
  multi-instance writes*, with a private-LRU arm that avoided it only by not serving from cache.

### Book-visible coverage gaps to record (each a gap, never a null result)

- **Database soft deletes / tombstones.** There is no soft-delete or tombstone *schema* anywhere in
  Study 05. The only tombstone in the corpus is the cache-side pre-write tombstone used by a fault
  injection (`db-rollback-after-pre-write-tombstone`). "Soft deletes let you avoid invalidation" is
  entirely untested.
- **Outbox-driven invalidation.** `cache_outbox` and `a_outbox_orphans` exist in the owned schema,
  but no claim isolates outbox consumption, ordering or durability as a tactic; it is present, not
  evidenced.
- **Process-local caching with per-key partitioning.** Never built; the memory 3-instance arm
  bypassed, so the correct multi-instance form of a private cache is unmeasured.
- **Hard TTL expiry, refreshes, write-rate burst.** `expired_hard = 0`; refresh counters absent; the
  5×-for-120 s burst has no write-rate knob.
- **Scale and eviction.** The harness tops out at 3 000 people, not the requested 8k/80k/800k, and
  Redis evictions were 0 (cache fit the working set), so no eviction effect is claimed.
- **Placement and topology.** The colocated/non-colocated pair is ~1.3–1.5x faster colocated
  (y1 882/852/558 vs y2 673/557/387) but the tablet/leader distribution was never produced, and one
  host is not colocation evidence.
- **Equal-total** is now labelled (`resource_framing: equal-total`, warm 10 837 vs 8 027 db-only)
  but is a single run and must never be pooled with `db-only`.
- **Open-loop demand** is still a proposed task.

## Reasoning and trade-offs

- The map answers the question's "which guarantees are supported" precisely: only fence-based
  strictness is supported across instances; relaxed is supported only with its wrong-read rate
  printed beside throughput, which the report already does.
- It also tells a reader what ownership buys: fewer moving parts and a durable hook (version +
  outbox), but not consistency — a genuinely useful and slightly counter-intuitive result.
- Cost: reworking `concepts/cache-consistency.typ` and family 6 of `major-families.typ` to state the
  two axes, and adding a labelled gap list. No new measurement is proposed here.
- Deliberately excluded: any statement that relaxed staleness is a bug. It is a *promise*; the book's
  job is to define "wrong" relative to the stated contract and never let a permitted staleness read
  as a defect or as freshness.

## Assumptions and evidence gaps

Observed: the 27-cell all-green report and its wrong-read table; the one-vs-three table's per-arm
notes; `v2-14/15/16` and `v2-gap-04`; the fence discussion §1–§3; the churn and resources/engines
analyses; `sql/README.md` mutation lists and audit list. Inferred, not measured: that the
process-local arm's zeros are a bypass effect (the report's own notes say the instances have
separate stores and that writes are invisible to the others — but there is no dedicated experiment
isolating "per-key partitioned private cache"); that the fence is the *only* mechanism that could
have produced the strict arms' zeros (true within this study, not a proof of necessity in general).

## Alternatives

- **Keep the flat pattern list** (current `minor-variants.typ` + family 6). Cheapest, but hides
  exactly the two things that matter most: ownership did not fix staleness, and the zeros in the
  private-cache arm are bypasses.
- **Map by strategy first** (aside/through). Wrong axis: those are plumbing, and the study's own
  churn numbers show they agree.
- **Write the extra measurements before writing the map.** Rejected: this brainstorm creates no
  tasks, and the existing evidence already supports a *qualified* map with labelled gaps.

## Disagreements

Anticipated objections, recorded without reading other independent positions:

- "The 87% is a deliberate, permitted relaxation, so it should not be called wrong." Agreed — the
  book must define it against the strict contract and present it as the *cost of the promise*, not
  as a defect. The number still belongs on the page beside throughput.
- "One trial, one machine: do not generalize." Agreed; the map must label it single-run and the
  mechanism claim must stay a mechanism claim, per its own `strength`.
- "Soft deletes/tombstones and outboxes are standard, so the book may recommend them." Only as
  labelled inference, not as measured advice, until a design and a run exist.

## Falsification tests

This recommendation changes if:

1. repeated, multi-machine, medium-scale trials show the shared-cache relaxed staleness is small or
   disappears (the 87% becomes an artifact of one arrangement);
2. a per-key-partitioned process-local cache is measured with high hit rate and zero wrong reads
   (then the "private cache = bypass" reading is wrong and the memory cell becomes a real result);
3. an outbox-drained invalidation design achieves strict freshness **without** a cache-side fence
   (then the fence-as-required claim weakens to fence-or-equivalent);
4. a database soft-delete/tombstone design is measured to avoid the fence or the lease (then the
   tombstone gap becomes a supported tactic instead of a gap).
