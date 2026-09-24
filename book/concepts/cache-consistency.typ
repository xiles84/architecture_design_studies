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
  [*5. Add performance controls separately.* A bounded fill lease suppresses duplicate fills; TTL and
   early expiry bound residency and load. Neither supplies currentness.],
)

*Forbidden is separate from stale.* A value that is impossible or was never committed is a defect
under every contract, not a relaxation of one. *Relaxed* means only "may serve a committed value that
is older than the strict comparator would allow, for a bounded window" — it never means dirty, torn
or uncommitted. Every measured cell here recorded zero impossible values; that is the negative control
on this sentence.

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

#heading(level: 2, "The named mechanisms, one per race")

A mechanism is not a synonym for freshness. Each one closes exactly one race, and the "does not solve"
line matters as much as the guarantee: five of these can be present and the contract still be broken.

#mechanism-group("Source correctness — arbitrating the mutation itself",
  [These protect the authoritative row. They decide *who wins* a contended update; they say nothing
   about a copy that has already left the transaction.])

#mechanism-card(
  "Database mutation lock",
  [two database writers updating the same row],
  [a pessimistic `SELECT … FOR UPDATE` serialises them before the update, and a conditional update
   (`WHERE … AND status = …`) refuses the loser],
  [at most one of the two writers commits the contended transition],
  [a stale read in another connection, a lost write outside the lock's scope, or any copy of the row],
  [active registered claim — `v2-12`, `v2-06`],
)

#mechanism-group("Copy correctness — ordering what leaves the transaction",
  [These protect the cache. Only one of them closes the stale-fill race, and only one of them is a
   performance control that proves nothing about freshness.])

#mechanism-card(
  "Invalidation",
  [an entry that is still resident after its source changed],
  [the writer removes or marks the entry stale after committing],
  [the entry stops being served *from that moment on* — call it what it is: a deletion with a message
   attached],
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
  [the stale-fill race specifically — "publish after commit" is necessary and *not* sufficient, so a
   strict writer fences before and after the commit and treats the acknowledgement, not the commit,
   as the contract boundary],
  [an unobserved writer, a dirty source read, or an indefinitely stale hit that is never revalidated],
  [active registered claim — `v2-15`],
)

#mechanism-card(
  "Source version token validated at the hit boundary",
  [a hit that never consults the source version],
  [every hit compares the entry's version against the source's current version and misses on a
   mismatch],
  [a process-local cache can be answerable: the hit carries its own proof],
  [the read cost it adds, and a shared relaxed cache that never consults the token — an unconsulted
   token is decoration, not a guarantee],
  [mechanism demonstrated — the owned process-local arms of `20260921T-survey3`; no registered rate
   for the validation cost],
)

#mechanism-card(
  "Fill lease",
  [duplicate concurrent fills of the same key — a stampede],
  [a bounded lease lets one filler work while the others wait, fall back or time out],
  [one filler at a time, so a cold key does not multiply its own load],
  [*any* freshness question. It is a load control; it guarantees nothing about the value filled],
  [newer signed but unregistered evidence — the 2026-09-23 Study 05 churn cell demonstrates the lease
   path (`20260923T1100Z-churn1800-through`); no registry claim states a lease rate],
)

#mechanism-card(
  "TTL and early expiry",
  [an entry that would otherwise live forever, and the thundering herd that follows a mass expiry],
  [a hard TTL bounds residency; probabilistic early expiry spreads the refill before the deadline],
  [bounded residency, and a bound on refill load],
  [currentness. An entry inside its TTL is *permitted* to be stale, which is the opposite of a
   freshness proof],
  [active registered claim — `v3-01`, `v3-02`; hard expiry has still never been observed to fire],
)

#mechanism-card(
  "Transactional outbox",
  [a change lost between the commit and its asynchronous notification],
  [the change is written in the same transaction as the mutation and consumed by the publisher, so
   the notification is durable and ordered after the commit],
  [a committed change cannot disappear before it is observed, and the delivery path is recoverable],
  [freshness. An outbox is a *prerequisite* for a strict cache, not a guarantee of one: strict
   after-acknowledgement hits still need an acknowledgement policy, a validation rule at the hit
   boundary, or an authoritative read that closes the window],
  [proposed / unmeasured — no committed measurement exists],
)

#heading(level: 2, "What each mechanism is actually for")

The dimensions are the point: a mechanism can be correct at the source, neutral for freshness, and
decisive for load — and reading a table like this is faster than holding five paragraphs in your head.

#table(
  columns: (1.5fr, 1.1fr, 1.1fr, 1fr, 0.9fr),
  stroke: 0.4pt + palette.rule,
  inset: 4pt,
  align: (left, center, center, center, center),
  table.header(
    [*Mechanism*], [*Source correctness*], [*Freshness proof*], [*Load control*], [*Recovery*],
  ),
  [Database mutation lock], [yes], [—], [sometimes a cost], [—],
  [Invalidation], [—], [partial: clears, does not fence], [—], [—],
  [Publication fence / CAS], [—], [yes, for the stale fill], [—], [—],
  [Source version at hit boundary], [—], [yes], [read overhead], [—],
  [Fill lease], [—], [—], [yes], [—],
  [TTL / early expiry], [—], [bounds residency, not proof], [yes], [fallback only],
  [Transactional outbox], [durable observation], [prerequisite], [—], [yes],
)

#text(size: 8.5pt, fill: palette.muted)[*Evidence status.* _Active registered claim_: the active
#raw(registry-label) states it. _Newer signed but unregistered evidence_: a signed
2026-09-23 study analysis demonstrates the mechanism, but no registry claim states a rate for it.
_Mechanism demonstrated_: the harness or study demonstrates the mechanism without a registered
number. _Proposed / unmeasured_: no committed measurement exists.]

#heading(level: 2, "Freshness is a contract, not a speed knob")
Every controlled strict/relaxed pair differed by −14.2% to +12.3% in read throughput — inside the
run's single-trial noise. The visible cost of strictness is on the write path, where a strict protocol
adds fences. Never restate this as "strict freshness is free".
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
process-local arms did serve hits and stayed correct only because they validated a source version at
the hit boundary. A private in-memory lease cannot coordinate instances; a source version token that
is never consulted at the hit or publication boundary does not make a shared relaxed cache strict.

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
