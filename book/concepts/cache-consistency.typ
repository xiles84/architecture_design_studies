#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Cache consistency")

#marker("concept", "cache consistency")
A cache is a read copy with a freshness contract. The design question is not "which cache pattern?"
but "what does a hit *prove*, and who can observe every change?". Invalidation is a *fence* — a rule
that refuses or repairs a publication whose input moved — not a best-effort deletion. This is the
concept the external-read-copies family exists to hold.

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

#heading(level: 2, "The named mechanisms, each closing one race")

#table(
  columns: (auto, 1fr, auto),
  stroke: 0.4pt + palette.rule,
  inset: 4pt,
  [*Mechanism*], [*The race it closes*], [*Evidence status*],
  [database mutation lock],
  [two database writers updating the same row; a pessimistic `SELECT … FOR UPDATE` serialises them
   before the update, and a conditional update (`WHERE … AND status = …`) refuses the loser],
  [active registered claim — `v2-12`, `v2-06`],
  [fill lease],
  [duplicate concurrent fills of the same key (a stampede): a bounded lease lets one filler work and
   makes the others wait, fall back or time out],
  [newer signed but unregistered evidence — the 2026-09-23 Study 05 churn cell demonstrates the lease
   path (`20260923T1100Z-churn1800-through`); no registry claim states a lease rate],
  [publication fence / CAS],
  [a reader that snapshotted committed state S0 republishes it *after* a writer commits S1 and
   invalidates — "publish after commit" is necessary and *not* sufficient; a strict writer fences
   before and after the commit, and the contract boundary is the acknowledgement, not the commit],
  [active registered claim — `v2-15`],
  [source version token],
  [a hit that never consults the source version. Validating a version at the hit boundary is what
   makes a process-local cache answerable; a token that is not consulted there does not make a shared
   relaxed cache strict],
  [mechanism demonstrated — the owned process-local arms of `20260921T-survey3`; no registered rate
   for the validation cost],
  [transactional outbox],
  [a change lost between commit and publication when delivery is asynchronous; consuming the outbox
   before acknowledging the writer closes it],
  [proposed / unmeasured],
)

#text(size: 8.5pt, fill: palette.muted)[*Evidence status.* _Active registered claim_: the active
`book/evidence/v3/claims.json` states it. _Newer signed but unregistered evidence_: a signed
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
single-instance relaxed cell recorded zero. "Relaxed" there is exactly the comparator sense above: the
hits served committed-but-older values, never impossible ones.

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
