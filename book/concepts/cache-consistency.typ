#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Cache consistency")

#marker("concept", "cache consistency")
A cache is a read copy with a freshness contract. Invalidation is a fence, not a deletion;
publication must be refused when the fence moved. A private in-memory lease does not coordinate
instances sharing a store.
#registry-card("v2-14-cache-throughput-gain")
#registry-card("v2-15-three-instance-staleness")
#registry-card("v2-16-strict-freshness-read-cost")
