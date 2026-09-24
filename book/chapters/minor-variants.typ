#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Minor variants")

#marker("chapter", "minor variants")
A minor variant is one controlled change inside a family. The corpus has measured pairs for some;
the rest are named with their evidence status so a reader can tell the difference.

#heading(level: 2, "Measured controlled pairs")
#heading(level: 3, "Index selection")
Adding the query index was the first large read win in the charity tree (D2 indexed, 13.3x the
minimal design's read score); flattening the foreign key reached 20.5x.
#registry-card("v2-01-normalized-index-first")

#heading(level: 3, "Optimistic versus pessimistic concurrency")
Under one hot key and 16 writers, lock-before-update delivered 1.76x the version-checked design — an
upper bound, because the run could not separate the lock's advantage from extra per-retry work. An
unchecked read-modify-write lost 93.7% of acknowledged updates.
#registry-card("v2-12-lock-vs-version-upper-bound")

#heading(level: 3, "Trigger versus application maintenance")
Trigger-maintained rollups cost 6.85x on whole-configuration replacement (83.7 vs 573 ops/s, p99 2.5
to 42 ms); the application rollup kept writes near the reference for the same read gain.
#registry-card("v2-13-trigger-rollup-cost")

#heading(level: 3, "Strict versus relaxed freshness")
*Strict after acknowledgement* means a read that starts after a write to the same key is acknowledged
never returns an older committed value. *Relaxed* is defined only against that comparator: it may
serve a committed value that is older, for a bounded window — never a dirty, torn or impossible
value. Every controlled pair differed by −14.2% to +12.3% in read throughput — noise. The visible cost
of strictness is on the write path, where a publication fence is added against the stale-fill race (a
reader republishing S0 after a writer commits S1 and invalidates).
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
