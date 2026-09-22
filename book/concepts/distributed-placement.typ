#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Distributed placement")

#marker("concept", "distributed placement")
Node count, replication factor and actual data placement are three separate facts. A 3-node cluster
on one laptop has the first two and none of the third's network consequences.

#heading(level: 2, "What the corpus measures")
Studies 01–03 run PostgreSQL single-node and YugabyteDB single-node and 3-node RF=3 cells. The
multi-node cells show placement-shaped effects — a section-sharded layout collapsed point lookups
to 26 ops/s on YugabyteDB — but the balance of tablets across nodes is unexplained.

#heading(level: 2, "What it does not measure")
Physical colocation is unverified, and balanced client access across endpoints is unproven. A
single query endpoint stands in for cluster access in every multi-node run. Both are recorded gaps,
not null results.

#registry-card("v2-gap-05-colocation-unverified")
#registry-card("v2-gap-07-real-network-balanced-endpoints")
#registry-card("v2-gap-06-native-datastore-families")

#heading(level: 2, "Boundaries")
Do not read a single-node-versus-cluster ratio as production capacity guidance. The nodes share the
host's cores and have no real network between them.
