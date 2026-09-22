# Progress — Study 07 topology protocol

**Task:** `task-20260922T025912Z-topology-study-protocol`
**Attempt:** `attempt-20260922T201438Z-54b7b8` (claim `claim-2c1288df37770024`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/topology-study-handoff-v1`

Planning only: no database started, no measurement.

## Decision Log

1. **Independent hosts are a hard prerequisite, not a variable.** A second environment page and a
   real network path gate every cell; a container-on-one-laptop run is invalid here. Confidence: High.
2. **Balanced endpoints must be measured.** Every client run records the per-endpoint connection and
   operation distribution; a single endpoint is a labelled control or a gap. Confidence: High.
3. **One axis at a time**, against a fixed baseline, so attribution is possible. Confidence: High.
4. **Correctness and negative controls apply to failure experiments too** (a no-op netem injection or
   a kill with no recorded failover is a harness gap). Confidence: High.
5. **Four execution tasks**: env/harness first (blocking), then node/replication/equal-total,
   placement/routing, network/failure. Each takes the benchmark lock. Confidence: High.
6. **Numbering:** this is Study 07; the native-models study is 06. Confidence: High.

## Residual risks

- A second host may not be available; the env/harness task must then escalate rather than run on one
  host under a "cluster" label.
- netem inside containers may behave differently from a real link; the no-op injection control is
  what detects that.
