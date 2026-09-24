#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Topologies")

#marker("chapter", "topologies")
Topology is a separate axis from design family: node count, replication, placement, endpoint
routing, resource framing, network conditions and failure behaviour. A design that is best on one
node is not thereby best on three.

#heading(level: 2, "What was measured")
Studies 01–03 run PostgreSQL single-node and YugabyteDB single-node (RF=1) and 3-node (RF=3). Every
cell is throttled to a documented per-node budget, and every manifest records it. The cost model is
therefore "per-node budget", not "equal total budget".

#heading(level: 2, "What topology does to a design")
Replication turns a write into a consensus write, so write-heavy designs move differently from
read-heavy ones. The section-sharded seat map collapsed point lookups to 26 ops/s on YugabyteDB,
while the same family's whole-map read was 4.5–11x a row design — one topology change, opposite
directions for two access patterns.

#registry-card("v2-11-section-seat-map")

#heading(level: 2, "Open topology gaps")
#gap[
  Physical placement is unverified (`v4-gap-01`) and endpoint distribution is measured only for the
  2026-09-23 Study 05 three-node cells (`v4-gap-02`); the earlier multi-node cells reach the cluster
  through a single endpoint, and there is no real network between "cluster" nodes. Multi-node results
  here are placement experiments, not distributed-systems results.
]

#heading(level: 2, "Boundaries")
Do not scale these numbers to production capacity. The nodes share one host's eight cores, are
CPU-throttled, and communicate over a loopback-class path.
