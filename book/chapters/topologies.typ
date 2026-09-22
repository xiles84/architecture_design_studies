#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Topologies")

#marker("chapter", "topologies")
Topology is a separate axis from design family: node count, replication, placement, endpoint
routing, resource framing, network conditions and failure behaviour. The measured corpus covers
PostgreSQL single-node and YugabyteDB single-node/3-node RF=3; placement and real-network coverage
are gaps.
#registry-card("v2-gap-05-colocation-unverified")
#registry-card("v2-gap-07-real-network-balanced-endpoints")
#registry-card("v2-11-section-seat-map")
