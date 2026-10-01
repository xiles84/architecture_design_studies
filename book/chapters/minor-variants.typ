#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Minor variants")

#marker("chapter", "minor variants")
A minor variant is one controlled change inside a family. The corpus has measured pairs for some;
the rest are named with their evidence status so a reader can tell the difference.

#heading(level: 2, "Measured controlled pairs")
#heading(level: 3, "Index selection")
*Read score* — used from here on — is this corpus's blended within-run composite of a family's read
queries. It compares variants inside one run and nothing else: it is not the throughput or the
latency of any single query, and it must not be quoted across runs.

Adding the query index was the first large read win in the charity tree (D2 indexed, about 10-11x the
minimal design's read score in two whole runs, 13.3x in the first); flattening the foreign key reached
about 14-16x (20.5x in the first).
#registry-card("v2-01-normalized-index-first")

#heading(level: 3, "Optimistic versus pessimistic concurrency")
Under one hot key and 16 writers, lock-before-update delivered 1.76x the version-checked design — an
upper bound, because the run could not separate the lock's advantage from extra per-retry work. An
unchecked read-modify-write lost 93.7% of acknowledged updates.
#registry-card("v2-12-lock-vs-version-upper-bound")

#heading(level: 3, "Trigger versus application maintenance")
Trigger-maintained rollups cost 4.0–5.7x on whole-configuration replacement across three
randomized-order runs (median 4.94x), and their latency tail is 22–46x worse (p99 61–109 ms against
2.2–2.7 ms); the application rollup kept writes near the reference. Repeating the measurement inside that
one configuration gave 35–115% disagreement, so quote the range and never a point.
#registry-card("v2-13-trigger-rollup-cost")

#heading(level: 3, "Strict versus relaxed freshness")
*Strict after acknowledgement* means a read that starts after a write to the same key is acknowledged
never returns an older committed value. *Relaxed* is a family of weaker contracts, not one contract:
compared with the strict promise just stated, a relaxed contract may permit a committed but older value
according to its declared semantics. Only a *bounded-staleness* member declares a maximum Δ; eventual and
best-effort freshness need not, and the relaxed contracts measured here declared no maximum, so they are
unbounded-relaxed. No member permits a dirty, torn or impossible value. Every controlled pair differed
by −14.2% to +12.3% in read throughput — noise. The visible cost of strictness is on the write path,
where a publication fence is added against the stale-fill race (a reader republishing S0 after a writer
commits S1 and invalidates).
#registry-card("v2-16-strict-freshness-read-cost")

#heading(level: 2, "Named but not measured as controlled pairs")
#list(
  [*Full versus bounded embedding.* Embedding was measured as a family, not against a bounded
   variant.],
  [*Locks / CAS / isolation levels.* The corpus has CAS and lock designs, but not an isolation-level
   sweep under one invariant.],
  [*Pre-created versus on-demand inventory.* Measured in ticketing as a family contrast with
   arbitration, not isolated from concurrency strategy.],
  [*Cache-aside versus write-through.* Both appear; the three-instance failure is about fencing, not
   about the two patterns head-to-head.],
)
#gap[
  No variant above is a cross-study ranking, and none is an equal-total-resource comparison. Treat
  the unmeasured variants as open, not as null results.
]
