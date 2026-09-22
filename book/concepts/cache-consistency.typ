#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Cache consistency")

#marker("concept", "cache consistency")
A cache is a read copy with a freshness contract. Invalidation is a fence — a rule that refuses a
publication whose input moved — not a best-effort deletion. This is the concept the external-
read-copies family exists to hold.

#heading(level: 2, "The gain")
Adding a cache to warm cacheable reads bought 2.2–3.5x in throughput and an order of magnitude in
p99 (15.8–27.1 ms to 0.8–2.0 ms). This is an *add-cache* result: the Redis container's own budget
is outside the comparison, so it is not an equal-total-resource efficiency claim.

#heading(level: 2, "The failure mode")
Three logical application instances sharing one cache under write-through + relaxed freshness
recorded about 87% wrong reads; every single-instance relaxed cell recorded zero. A private
in-memory lease cannot coordinate instances, and a database version token does not remove the
window by itself. The corrected attribution: 87.81% belongs to the owned model, 87.24% to legacy.

#heading(level: 2, "Freshness is a contract, not a speed knob")
Every controlled strict/relaxed pair differed by −14.2% to +12.3% in read throughput — inside the
run's noise. The visible cost sits on the write path, where strictness adds fences. Do not read
"strict freshness costs nothing" from the read table.

#registry-card("v2-14-cache-throughput-gain")
#registry-card("v2-15-three-instance-staleness")
#registry-card("v2-16-strict-freshness-read-cost")

#heading(level: 2, "Boundaries")
Redis was measured as a cache, never as an authoritative key-value store. JSONB is not a native
document store. Transfer either result to a native family only as a labelled analogy.
