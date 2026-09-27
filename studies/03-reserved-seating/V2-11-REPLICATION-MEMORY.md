# `v2-11` replication, third pass: the same reversal under a raised memory envelope

**Why a third pass.** The second pass (`20260927T111110Z`) lost two of its three cells before they measured
anything: connections terminated by the server (SQLSTATE 57P01), with the database log recording server
processes **killed by signal 9**, 122 lines in one cell. The topology allots the database **3 GB** and pins
YugabyteDB's own `memory_limit_hard_bytes` to **1.5 GiB**, and the kills are that limit being enforced
inside the container. So the second pass did not fail because of design order, and repeating it unchanged
would fail the same way.

**What changes: the resource envelope, and only that.**

| Setting | Passes 1 and 2, and the 2026-09-21 baseline | This pass |
|---|---|---|
| database container memory | 3 GB | **5 GB** (`ADS_DB_MEMORY=5g`) |
| YugabyteDB `memory_limit_hard_bytes` | 1.5 GiB | **3 GiB** (`ADS_YB_TSERVER_MEMORY=3221225472`) |
| designs, topology, scale, phases, trials, seed | — | unchanged |

**What that costs, stated before measuring.** Absolutes from this pass are **not** comparable with any
earlier run: a bigger memory budget changes caches and plan choices, and `README`/`make` conventions in this
repository treat results as comparable only within one resource envelope. What remains comparable is the
**ratio between the two designs inside this run**, which is the quantity `v2-11` claims, and the AB/BA
question against pass 1. Every absolute quoted from this pass must be labelled as the 5 GB envelope.

**Order.** `l3_section_sharded, l2_section_document, l1_claim_rows` — the same as pass 2, so if the loads
now complete, this run is the AB/BA partner of pass 1.

**Host state is recorded.** `record_host_state()` (added to `infra/lib.sh` with this pass) writes the
machine's memory, CPU count, container budgets and running containers to `host-state.txt` beside the
results, so a cell lost here can be attributed to a number rather than to a log reading afterwards.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| W1 | **The raised envelope removes the load failures.** All three cells complete their load and reach verification. | Any cell dying at load again, which would mean the memory limit is not the cause and the killing mechanism is still unknown. |
| W2 | **The reversal reproduces the claim's direction.** `l2_section_document` reads the 100k-seat map at least **4x** faster than `l1_claim_rows`. | A multiple below 4x, which would make the claim's direction depend on the order or on the envelope. |
| W3 | **The two passes agree.** This pass's multiple and pass 1's (27.8x) fall within **±40%** of each other. | More than 40% apart, which would put the claim's multiple inside the position/session band this corpus keeps measuring. |
| W4 | **The correctness gate stays green** for every cell that reaches it. | Any cell reaching verification and failing it. |

Recorded before measurement. If W2 and W3 hold with all cells green, `v2-11` is registered `run_replicated`
with `order: abba` and both passes as anchors; the 2026-09-21 baseline remains its primary anchor.

## Scope and boundaries

One engine, one topology, one scale, one host; the PostgreSQL half of the claim and the 3-node topology
remain unmeasured. A pass here says the sharded half and the engine-dependence are still *open*, not that
they are wrong.
