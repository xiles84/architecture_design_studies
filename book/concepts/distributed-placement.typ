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

Study 05 adds two things those cells did not have: a key-locality pair, where the same donor's
donations are held in one tablet in one layout and hashed across tablets in the other, and
three-node cells whose client spreads operations over all three endpoints rather than one.

#heading(level: 2, "What it does not measure")
Two different limitations, and they should not be collapsed into one sentence.

*Physical placement is unverified.* The locality pair ran, but the runner produced no tablet or
leader distribution, so the layout labels record intent, not observed placement. That is an
*instrument* gap: it needs a measurement the harness cannot take, not just another run.

*Endpoint distribution is partial.* The 2026-09-23 Study 05 three-node cells spread operations
round-robin over three endpoints, but the earlier multi-node cells still reach the cluster through a
single endpoint. Their ratios remain valid for a client that does one endpoint's worth of cluster
access, and none of it is a network result.

Both are recorded gaps, not null results.

#registry-card("v4-gap-01-physical-placement-unverified")
#registry-card("v4-gap-02-endpoint-distribution-partial")
#registry-card("v2-gap-06-native-datastore-families")

#heading(level: 2, "Boundaries")
Do not read a single-node-versus-cluster ratio as production capacity guidance. The nodes share the
host's cores and have no real network between them.
