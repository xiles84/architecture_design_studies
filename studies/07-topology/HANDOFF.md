# Execution Handoff — Study 07, topology

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-topology-study-protocol` (HIGH planning).
**Required tag on integration:** `repo/topology-study-handoff-v1`.
**Numbering:** Study 06 is the native-models study; this is Study 07.

## 1. The question

How do node count, replication, physical placement, endpoint routing, resource budget, network
conditions and failure behaviour *separately* affect a data-architecture design? The corpus has
single-laptop multi-container "clusters" and a shared-core caveat; this study is the first that can
call its topology a topology.

## 2. Non-negotiables

- **Independent hosts.** Nodes are separate machines or VMs with a real network path, not containers
  on one laptop. A topology cell run on one host is invalid for this study.
- **Balanced endpoints.** Every client run must connect through a load balancer or a client-side
  round-robin, and must record the per-endpoint connection and operation distribution. A single
  endpoint is a named gap, not a cluster result.
- **Separate axes.** Node count, replication factor, placement, routing, resource budget, network
  and failure are varied one at a time against a fixed baseline.
- **Correctness gates timing**, and every failure experiment has a negative control that must be
  seen to violate the invariant.

## 3. Environment prerequisite (execution task 1)

Record at least two hosts in `docs/environments/` with CPU, memory, disk, kernel, network path,
clock source and the ways each can mislead a benchmark. A second environment page is a blocking
prerequisite: without it the study cannot state "independent hosts".

## 4. Axes and cells

| Axis | Baseline | Variants |
|---|---|---|
| Node count | 1 database node | 3 nodes, same total budget and same per-node budget (two separate arms) |
| Replication | RF=1 | RF=3 synchronous, RF=3 with one follower lagging |
| Placement | default distribution | pinned/colocated vs deliberately non-colocated, with the tablet/leader distribution recorded |
| Routing | one balanced endpoint | single endpoint (control), round-robin, least-connections |
| Resource budget | per-node fixed | equal-total control (charge all nodes to one budget) |
| Network | LAN | added latency (netem) at 1 ms, 10 ms, 50 ms RTT; one packet-loss level |
| Failure | none | kill one node mid-run; pause one node; lose one network link |

Workload: the Study 02 ticketing race and the Study 03 seat-map read are the fixed workloads, so
results can be compared across topologies without changing the business operation.

## 5. Correctness and negative controls per axis

- Node count/replication: the ticketing invariant must hold at every level; the unchecked
  read-modify-write control must oversell.
- Placement: the same logical read must return the same answer colocated and non-colocated; a
  deliberately stale replica control must be seen to serve stale data.
- Routing: a deliberately single-endpoint arm must show the skew the balanced arm removes.
- Network: the added-latency arm must change latency in the direction injected; a no-op injection is
  a harness failure, not a finding.
- Failure: the killed-node arm must record the failover event; the "kill and observe no effect" case
  is reported as a coverage gap.

## 6. Resources and sizing (per methodology §4)

Document host resources, client/system reserves, node count, per-node CPU/memory, aggregate budget
and the equal-total arm; calculate client readers, writers and connection-pool allowance; measure
per-node throttling and endpoint distribution. Never pool the per-node and equal-total arms.

## 7. Deliverables and separately claimable execution tasks

1. **Topology env and harness** — second host environment page, endpoint-balancing and distribution
   instrumentation, netem helper, node-kill helper; correctness gate unchanged. (infrastructure)
2. **Topology node/replication/resource arms** — node count, replication and the equal-total arm.
3. **Topology placement/routing arms** — placement with recorded distribution, routing and the
   single-endpoint control.
4. **Topology network/failure arms** — netem latency/loss and the failure experiments with controls.

Each task takes the benchmark lock, runs one matrix at a time across the hosts, and produces a
generated report plus a signed analysis under `studies/07-topology/`.

## 8. Acceptance criteria for the study report

- Every cell records the host list, the network path and the per-endpoint distribution.
- Placement is verified with a recorded tablet/leader distribution.
- Per-node and equal-total arms labelled and separate; no cross-arm pooling.
- Correctness gate passes and every negative control is seen to fire.
- No container-on-one-laptop cell is labelled a cluster result.
- Every number resolved to a v2 evidence claim before it enters the book.
