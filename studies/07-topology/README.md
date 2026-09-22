# Study 07 — topology

**Question:** how do node count, replication, placement, routing, resource budget, network and
failure affect a data-architecture design *separately*?

**Status:** protocol only; **no measurement yet**. The decision-complete handoff is
[`HANDOFF.md`](HANDOFF.md) (tag `repo/topology-study-handoff-v1`).

**Why it exists:** every multi-node result in Studies 01–05 is three containers on one laptop with
shared cores and no real network. The evidence registry records this as
`v2-gap-07-real-network-balanced-endpoints` and `v2-gap-05-colocation-unverified`. This study is the
one that can close them.

**Non-negotiables:** independent hosts, a real network path, balanced client endpoints with a
recorded per-endpoint distribution, one axis varied at a time, correctness gates timing, and a
negative control that fires for every failure experiment.

**Runner status:** no runner exists yet. Execution task 1 builds the second-host environment page,
endpoint-balancing instrumentation, the netem helper and the node-kill helper; tasks 2–4 run the
axes.
