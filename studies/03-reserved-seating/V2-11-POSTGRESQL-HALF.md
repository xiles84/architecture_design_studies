# Replication: the PostgreSQL half of the seat map (study 03, `v2-11`)

**Question.** `v2-11`'s statement has two halves. The first — that a stored section document answers the
100k-seat map far faster than one row per seat on YugabyteDB — has now been replicated twice, in both design
orders, and registered. The second — *"PostgreSQL shows no attributable difference"* — has **never been
re-measured**. It is also what makes the claim interesting: without it, the finding is a fact about one
database rather than evidence that the same design decision behaves differently by engine.

**What is held constant.** Designs `l1_claim_rows`, `l2_section_document`, `l3_section_sharded`; topology
`pg-single`; scale `small`; phases `verify,read`; three trials per cell; dataset seed 42. The comparison is
against the same designs on `yb-single` (two passes, 2026-09-26/27) and against the 2026-09-21 run of the
whole matrix.

**Environment.** The podman machine has 2 GiB and the declared budgets (database 3 GB, client 2 GB) cannot be
enforced on it; that is an open owner decision and is recorded in `CONTEXT.md`. **Everything in this tranche
runs under the same ceiling as every other run on this host**, so it is comparable with them — but its
absolute numbers inherit that caveat and will be labelled with it.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | **The PostgreSQL half reproduces.** On `pg-single` the 100k-seat map multiple `l2/l1` falls between **0.6x and 1.7x** — no attributable difference — where the same pair gives 24–28x on `yb-single`. | A multiple outside that band, which would mean the design decision helps PostgreSQL too and the claim's engine-dependence is wrong. |
| P2 | **`l3`'s load failure is specific to the YugabyteDB topology.** `l3_section_sharded` completes its load on `pg-single`. | It fails on PostgreSQL as well, which would move the blame from the engine's memory behaviour to the design's own load — a harness or schema defect. |
| P3 | **The correctness gate stays green** for every cell that reaches it. | Any cell verifying and failing. |
| P4 | **The engine contrast is the finding.** `l2/l1` on `pg-single` is at least 10x smaller than on `yb-single`. | A contrast of less than 10x, which would leave the claim's engine-dependence supported only weakly by replication. |

Recorded before measurement. If P1 fails, this tranche contradicts a published conclusion and **will not be
used to rewrite it**: it will be reported as an escalation with the run as evidence, because a worker
changing a signed claim's direction to fit a new measurement is exactly the move the workflow forbids.

## Scope and boundaries

| Field | Value |
|---|---|
| Runs | one, `pg-single`, three designs, `-trials 3` |
| Compared with | `yb-single` passes 1 and 3, and the 2026-09-21 matrix run |
| Boundary | One host, two engines, one scale. The 3-node topology is still unmeasured, and the memory ceiling applies to every number here. |
