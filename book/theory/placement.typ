#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#pagebreak(weak: true)
#heading(level: 2, numbering: none, "Cluster placement")
A logical layout is an intention. A placement claim requires observed tablet, leader and endpoint behavior, alongside the resource budget.

Strategies: #link(<theory-placement>)[Colocation and distributed access].

#strategy("Colocation and distributed access", "theory-placement", [Place related work deliberately, then verify where the engine actually stores and serves it.],
  [Joining related rows can involve several nodes. Containers sharing a laptop do not demonstrate database-level data colocation.],
  ([Hold logical data and replication fixed.], [Verify actual data and leader placement.], [Separate balanced endpoints from equal resource budgets.],),
  ([Related operations align with a verified placement key.], [Node and client budgets are explicit.], [Per-node distribution is observed.],),
  ([Layout labels stand in for placement evidence.], [Cluster nodes silently receive a larger total budget.], [One endpoint is assumed to represent balanced cluster capacity.],),
  ([May add or avoid cross-node work.], [Replication and coordination costs remain.], [Replicated data and engine metadata.], [Placement must preserve the logical answer and durability contract.], [Inspect per-node resources and distribution.],),
  figure-evidence("../assets/fig-topology-physical-boundary.svg", "structure", "conceptual illustration", "Colocation and distributed access: mechanism and trade-off.", [Logical keys and placement labels are separate from observed tablets, leaders and serving endpoints. Containers sharing a host do not establish physical data colocation.]),
  ([The same logical answer and resource framing are compared.], [Tables, tablets, replicas and leaders.], [Hot partitions and serving endpoints.], [No automatic read/write relocation.], [Replicas and any colocated derived state.], [Replicas or copies under their consistency contract.], [Unverified placement, skew or unbalanced serving.], [Match resources and workload; verify placement before interpreting effects.],), "case-cache", measured: true)
