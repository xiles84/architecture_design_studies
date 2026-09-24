#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card, registry-label

#heading("Cache consistency")

#marker("concept", "cache consistency")
A cache is a read copy with a freshness contract. The design question is not "which cache pattern?"
but "what does a hit *prove*, and who can observe every change?".

*Invalidation and publication fencing are different operations.* Invalidation removes or marks an
existing cache entry stale. It cannot stop a reader that loaded an older value from publishing it
afterwards — which is exactly the race the stale-fill figure below draws. A *publication fence* is the
additional check that refuses that fill when the source moved after the reader took its snapshot.
Writing "invalidation is a fence" makes the fence sound like a deletion, and the deletion is the part
that does not close the race. This is the concept the external-read-copies family exists to hold.

#heading(level: 2, "The flow: contract first, pattern last")

#list(
  [*1. Name the contract.* _Strict after acknowledgement_: a read that starts after a write to the
   same key is acknowledged must not return an older committed value. Or a *bounded-staleness*
   budget, or *session / read-your-writes*, or best effort. State which one, on the page.],
  [*2. Establish change-observation completeness.* Can every mutation participate — bulk jobs,
   migrations, deletes, recreates, external writers? If not, a strict *cache-hit* proof is
   unavailable: use an authoritative read, a complete change source, or downgrade the contract
   explicitly.],
  [*3. Identify the source capabilities.* A frozen schema still supports service-side generations and
   authoritative reads; an owned schema can add a monotonic version, a transactional outbox and a
   deletion generation. These improve proof and recovery; they do not confer freshness.],
  [*4. Choose topology and the publication protocol.* Process-local versus shared decides the
   coordination domain; cache-aside versus write-through decides who fills or republishes. State the
   ordering proof — one of the named mechanisms below.],
  [*5. Add performance controls separately.* A bounded fill lease suppresses duplicate fills while it
   is held; a hard TTL bounds how long an entry is *eligible to serve* and probabilistic early expiry
   spreads the refill work. None of these supplies currentness.],
)

*Forbidden is separate from stale.* A value that is impossible or was never committed is a defect
under every contract, not a relaxation of one. *Relaxed* is a family of weakenings rather than one
contract, so which member is in force has to be named:

#list(
  [*Strict after acknowledgement* — the contract the comparator below enforces.],
  [*Session / read-your-writes* — a client can see its own committed writes.],
  [*Bounded staleness* (≤ Δ) — a hit may lag the latest commit by at most a declared Δ.],
  [*Eventual freshness* — convergence is required; no distance or deadline is declared.],
  [*Best effort* — nothing beyond a valid committed value is promised.],
)

None of these permits a dirty, torn or uncommitted value, and every measured cell here recorded zero
impossible values — that is the negative control on this sentence. The study's measured relaxed cells
are **unbounded-relaxed**: they permit committed-but-older values relative to the strict comparator
without declaring a maximum staleness. They are therefore not bounded-staleness contracts, and their
violation count says nothing about a Δ window.

#heading(level: 2, "The strict contract, written down")

*Strict after acknowledgement* is the contract the measured cells are judged against, so it is worth
stating formally rather than in prose. Let `W1` be a write to a key, `ack(W1)` the moment its
acknowledgement is observed, `R1` a read of the same key, and `begin(R1)` the moment it starts.

#align(center)[
  #box(inset: 6pt, stroke: 0.5pt + palette.rule, radius: 3pt)[
    `begin(R1) > ack(W1)` $⇒$ `R1` must not return a version older than `W1`
    \
    #text(size: 8pt, fill: palette.muted)[unless a later write to the key supersedes `W1`]
  ]
]

Every other contract is a weakening of that line, and the weakening must be named: a
*bounded-staleness* budget, *session / read-your-writes*, or *best effort*. The comparator is what the
harness count below is measured against — a "wrong read" is a read that violates this line, not a read
that returns something impossible.

#heading(level: 2, "The named mechanisms, one primary failure mode each")

A mechanism is not a synonym for freshness. Each one addresses one primary failure mode or control
objective, and the "does not solve" line matters as much as the guarantee: several of these can be
present and the contract still be broken.

#mechanism-group("Source correctness — arbitrating the mutation itself",
  [These protect the authoritative row. They decide *who wins* a contended update; they say nothing
   about a copy that has already left the transaction.])

#mechanism-card(
  "Pessimistic lock + invariant recheck",
  [two database writers updating the same row],
  [`SELECT … FOR UPDATE` serialises the writers, and the writer re-reads the row under the lock and
   rechecks the transition's precondition before applying it],
  [serialisation: the two writers do not interleave. Rejecting the obsolete operation is the
   *recheck's* contribution, not the lock's — a lock alone guarantees neither],
  [a stale read in another connection, a lost write outside the lock's scope, or any copy of the row],
  [active registered claim — `v2-12`],
)

#mechanism-card(
  "Conditional update / compare-and-set",
  [an update applied to a row whose state has already moved on],
  [the expected state or version goes into the write itself (`WHERE … AND status = …`); a zero-row
   update reports that another writer got there first],
  [the transition applies only while the row still satisfies the expectation, so the loser is rejected
   rather than merely serialised],
  [the retry cost the loser now pays, and any copy of the row],
  [active registered claim — `v2-06`],
)

#mechanism-group("Copy correctness — ordering what leaves the transaction",
  [These protect the cache. One of them closes the stale-fill race; two of them — the fill lease and
   the TTL — are load controls that prove nothing about freshness.])

#mechanism-card(
  "Invalidation",
  [an entry that is still resident after its source changed],
  [the writer removes or marks the entry stale after committing],
  [once the invalidation has been applied, subsequent lookups of that cache entry miss or see it
   marked stale — call it what it is: a deletion with a message attached. Reads and fills already in
   flight are outside this guarantee],
  [the stale-fill race. A reader that snapshotted S0 before the invalidation can still publish S0
   afterwards; deletion orders nothing about a fill already in flight],
  [mechanism used by every measured write-through cell; on its own it never made a shared relaxed
   cache strict],
)

#mechanism-card(
  "Publication fence / CAS",
  [a reader that snapshotted committed state S0 republishes it *after* a writer commits S1 and
   invalidates],
  [before publishing, verify that the source generation still equals the generation the reader
   observed; refuse the fill if it moved],
  [the stale-fill race specifically — "publish after commit" is necessary and *not* sufficient. In
   the measured strict writer protocol the writer fences before and after the commit and treats the
   acknowledgement, not the commit, as the contract boundary; a strict contract can also be
   established without writer-side fencing, by synchronous source validation or an authoritative
   read],
  [an unobserved writer, a dirty source read, or an indefinitely stale hit that is never revalidated],
  [active registered claim — `v2-15`],
)

#mechanism-card(
  "Source version token validated at the hit boundary",
  [a hit that never consults the source version],
  [every hit compares the entry's version against the source's current version and misses on a
   mismatch],
  [a process-local cache can be answerable: the hit carries its own proof — *provided* every
   relevant mutation advances the token atomically, including deletes, recreates and writers outside
   the application, and the validating read is authoritative enough for the declared contract],
  [the read cost it adds, and a shared relaxed cache that never consults the token — an unconsulted
   token is decoration, not a guarantee],
  [mechanism demonstrated — the owned process-local arms of `20260921T-survey3`; no registered rate
   for the validation cost],
)

#mechanism-card(
  "Fill lease",
  [duplicate concurrent fills of the same key — a stampede],
  [a bounded lease lets one filler work while the others wait, fall back or time out],
  [Effect: duplicate fills are suppressed *while a valid lease is held*. Duplicate work can reappear
   if the lease expires or ownership is lost before the first filler finishes — which is exactly why a
   lease is not a fence],
  [*any* freshness question. It is a load control; it guarantees nothing about the value filled],
  [newer signed but unregistered evidence — the 2026-09-23 Study 05 churn cell demonstrates the lease
   path (`20260923T1100Z-churn1800-through`); no registry claim states a lease rate],
)

#mechanism-card(
  "TTL and early expiry",
  [an entry that would otherwise live forever, and the thundering herd that follows a mass expiry],
  [a hard TTL bounds how long an entry is *eligible to serve*; probabilistic early expiry spreads the
   refill work before many entries reach the same deadline],
  [Effect: logical validity is bounded, and refill work is spread. Neither is a hard upper bound on
   refill traffic],
  [currentness. An entry inside its TTL is *permitted* to be stale, which is the opposite of a
   freshness proof, and early expiry does not by itself cap refill load],
  [active registered claim — `v3-01`, `v3-02`; hard expiry has still never been observed to fire],
)

#mechanism-card(
  "Transactional outbox",
  [a change lost between the commit and its asynchronous notification],
  [the change is written in the same transaction as the mutation and consumed by the publisher, so
   the notification is durable and ordered after the commit],
  [a committed change cannot disappear between the commit and its notification, and the delivery
   path is recoverable],
  [freshness, and it is not universally required. *When* coherence depends on asynchronous
   publication, an outbox is what makes the committed change durably observable; strictness can also
   be established by synchronous source validation or an authoritative read, with no asynchronous step
   to protect. Even with an outbox, strict after-acknowledgement hits still need an acknowledgement
   policy and a read protocol that closes the window],
  [proposed / unmeasured — no committed measurement exists],
)

#heading(level: 2, "What each mechanism is actually for")

The dimensions are the point, and the header is short so no word breaks in a narrow column:
*mutation* is whether the mechanism makes the source transition correct, *observation* whether a
committed change is durably seen, *freshness* whether a read is proved current, *load* whether it
changes how much work the system does, and *recovery* whether it helps reconcile after a failure.
Reading a table like this is faster than holding eight cards in your head.

#table(
  columns: (1.4fr, 1fr, 1.05fr, 1.1fr, 0.95fr, 0.8fr),
  stroke: 0.4pt + palette.rule,
  inset: 4pt,
  align: (left, center, center, center, center, center),
  table.header(
    [*Mechanism*], [*Mutation*], [*Observation*], [*Freshness*], [*Load*], [*Recovery*],
  ),
  [Pessimistic lock + recheck], [yes], [—], [—], [may cost], [—],
  [Conditional update / CAS], [yes], [—], [—], [may cost], [—],
  [Invalidation], [—], [partial: clears, does not fence], [partial], [—], [—],
  [Publication fence / CAS], [—], [—], [yes, race-specific], [—], [—],
  [Source version at hit boundary], [—], [—], [yes], [read cost], [—],
  [Fill lease], [—], [—], [—], [yes, while held], [—],
  [TTL / early expiry], [—], [—], [bounds eligibility, not proof], [yes], [fallback only],
  [Transactional outbox], [—], [yes], [prerequisite in async designs], [—], [yes],
)

#text(size: 8.5pt, fill: palette.muted)[*Evidence status.* _Active registered claim_: the active
#raw(registry-label) states it. _Newer signed but unregistered evidence_: a signed
2026-09-23 study analysis demonstrates the mechanism, but no registry claim states a rate for it.
_Mechanism demonstrated_: the harness or study demonstrates the mechanism without a registered
number. _Proposed / unmeasured_: no committed measurement exists.]

#heading(level: 2, "Freshness is a contract, not a speed knob")
Every controlled strict/relaxed pair differed by −14.2% to +12.3% in read throughput — inside the
run's single-trial noise. Mechanically a strict protocol adds write-path fencing, but the claim beside
this paragraph measures read throughput only and does not quantify that cost, so quote no write-path
figure from it. What the read result does forbid is the other direction: never restate this as "strict
freshness is free".
#registry-card("v2-16-strict-freshness-read-cost")

#heading(level: 2, "The gain, and the framing it does not have")
Adding a cache to warm cacheable reads bought 2.2–3.5x throughput and an order of magnitude in p99
(15.8–27.1 ms to 0.8–2.0 ms). That is an *add-cache* result under `db-only` framing: the cache's own
CPU and memory are outside the comparison (`conf-04`). A separately labelled equal-total arm now
exists — the same design and topology read 11k ops/s equal-total against 8.0k db-only — and the two
framings are never pooled.
#registry-card("v2-14-cache-throughput-gain")
#registry-card("v3-03-equal-total-framing-labelled")

#heading(level: 2, "The failure mode, stated against the strict comparator")
Three logical application instances sharing one Redis cache on write-through + *relaxed* freshness
recorded 87.81% wrong reads on the owned model and 87.24% on the legacy model, while every
single-instance relaxed cell recorded zero.

*What the counter counts.* The harness metric is named `wrong_reads`, and it counts reads that violated
the strict comparator above. In this cell those were committed-but-older values — the legitimate
consequence of choosing a relaxed contract — and *not* dirty, torn, uncommitted or impossible values.
Read it as *strict-freshness violations*, never as corruption: every measured cell here recorded zero
impossible values. A *bounded* staleness contract would raise a further question this run cannot
answer — whether the stale values also left the permitted window — so the counter is only evidence
about staleness when the contract is bounded and the window is known.

The process-local arms are the trap. The legacy process-local arms' zero wrong reads are a *bypass*
result, not safety: they served *0* cache hits across 15,578–21,756 reads, so the cache was
effectively absent (`reports/20260921T-survey3.md`, the anchor of the registered claim). The owned
process-local arms did serve hits and validated the source version at the hit boundary; those cells
recorded no strict-comparator violations. Validation at the hit boundary is what *can* supply that
proof — but whether it is the only thing that could is not tested here, because no unvalidated owned
arm was run. Read it as the mechanism these cells are consistent with, not as an ablation. A private
in-memory lease cannot coordinate instances; a source version token that is never consulted at the hit
or publication boundary does not make a shared relaxed cache strict.

#figure-evidence(
  "../assets/fig-cache-stale-fill.svg",
  "sequence",
  "observed result",
  "The stale-fill race, and the fence that closes it.",
  [A reader snapshots committed state S0, a writer commits S1 and invalidates the key, and then the
   reader publishes S0 as a fill — so the next reader hits a committed-but-older value. A fenced
   protocol refuses that fill by checking the source version before publishing. The figure illustrates
   the registered claim beside it and embeds no rate of its own. Source:
   `book/assets/sources/fig_cache_stale_fill.puml`],
  claim: "v2-15",
)
#registry-card("v2-15-three-instance-staleness")

#heading(level: 2, "Churn, engines and placement (2026-09-23)")
Sustained churn crossed real 300 s TTL boundaries under load and stayed correct: a 1,801 s
write-through cell crossed 6.0 boundaries with 0 wrong reads in 2,879,077 reads, and the two 601 s
write-through / write-aside cells crossed 2.0 each with 0 wrong reads. The hard-expiry path fired
*0* times — only probabilistic early expiry did — so "an entry reached its hard TTL" is still not
demonstrated and is carried as a gap below. The strict/relaxed pair now also runs per engine, with the
client spread over the endpoints: on YugabyteDB 1-node strict read 514 against relaxed 596 app-99
ops/s, on 3-node 942 against 846; the direction differs per engine and neither is attributed. The
colocated/non-colocated key pair favours the colocated layout (warm 882 against 673 ops/s,
cacheable-99 852 against 520) as engine locality on one host, not verified colocation.
#registry-card("v3-01-churn-crosses-real-ttl")
#registry-card("v3-02-churn-600-through-vs-aside")
#registry-card("v3-04-yb-single-strict-relaxed")
#registry-card("v3-05-yb-cluster3-strict-relaxed")
#registry-card("v3-06-colocated-vs-noncolocated-locality")

#heading(level: 2, "What is still open")
#gap[
  The hard TTL has still never been observed to remove an entry; medium scale beyond the harness
  maximum, the tablet/leader placement distribution and open-loop arrival demand are unmeasured.
  Treat these as open, never as null results.
]
#registry-card("v3-gap-01-study05-remaining-regimes")

#heading(level: 2, "Boundaries")
Redis was measured as a cache, never as an authoritative key-value store. JSONB is not a native
document store. Transfer either result to a native family only as a labelled analogy.
