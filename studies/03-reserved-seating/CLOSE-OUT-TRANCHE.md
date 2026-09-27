# Close-out tranche: the PostgreSQL pair reversed, and `l3` alone (study 03)

Two open items from the `v2-11` work, measured in one session because measurement time is the scarce
resource here, not the number of runs. They are independent, and each is registered separately below.

## Part 1 — the PostgreSQL pair in the opposite order

The PostgreSQL half was measured once, with `l1_claim_rows` before `l2_section_document`, and gave a
**0.83x** multiple. The YugabyteDB pair has an AB/BA control (27.8x and 24.4x); the PostgreSQL pair does
not, and on that engine the multiple is small enough that an order effect could matter more than the
difference being measured. This pass runs **`l2` before `l1`**.

| Id | Prediction | Falsifier |
|---|---|---|
| P5 | **The PostgreSQL half survives reversal.** The 100k-seat map multiple stays between **0.6x and 1.7x** with the pair reversed. | A multiple outside that band, which would mean the "no attributable difference" half depends on which design runs first. |
| P6 | **The engine contrast survives.** The PostgreSQL multiple remains at least **10x smaller** than the YugabyteDB one. | A contrast below 10x, which would weaken the engine-dependence this claim now rests on. |
| P7 | **Correctness stays green.** Both cells pass verification. | Any cell failing the gate. |

## Part 2 — `l3_section_sharded` alone on `yb-single`

`l3`'s load has failed in every pass it has been part of, and the earlier diagnostic that would have run it
alone never started because the podman machine was stopped. It is run **alone** here: one design, nothing
before or after it.

| Id | Prediction | Falsifier |
|---|---|---|
| D1 | **The load completes when nothing else runs first**, which would make the failure a sequencing or capacity effect rather than a property of the design. | It failing again with an empty session around it, which makes the design's load the problem and leaves the sharded half of `v2-11` measurable only by fixing that. |
| D2 | **If it loads, its answers are right**: verification passes against Go-computed truth. | Loading but failing verification, which would make the design's answers wrong rather than merely slow. |

**Explicitly not claimed:** Part 2 is a diagnostic. A success would say the earlier failures were about the
state other cells leave behind; it would **not** retroactively validate the three failed cells, and it will
not be folded into `v2-11`'s registered numbers. The claim's sharded element stays as the 2026-09-21 run
left it.

## Environment and boundaries

Both parts run under the same 2 GiB podman machine as every other run on this host; the declared container
budgets cannot be enforced and that remains an open owner decision recorded in `CONTEXT.md`. Part 2 uses the
default envelope (database 3 GB declared) so it is comparable with the passes whose failures it is
explaining, rather than with the 5 GB pass. One scale, one host, one engine per part.

---

## Results

*(To be written after measuring. This file is the pre-registration; the results belong in a signed analysis
under `reports/analyses/`, and nothing may be added here that was not predicted before the runs.)*
