# Client sizing for the read phase (study 02)

**Question.** The 2026-09-29 re-run showed the read phase's client throttling is not a period artefact:
`r1_inventory_row` loses 420 of 634 CFS periods to its two-core quota and its `q01_event_availability` reads
sit flat at ~38k ops/s across all five tiers. Before any read level can be quoted again, the owner's sizing
rule (2026-09-14, methodology 4) asks a question this tranche answers: **is the bind the client's own quota,
or the machine?** The podman machine has **4 CPUs**; study 02 budgets the database 2 and the client 2, so the
two already sum to the whole machine and the client cannot be sized up without taking from the database. That
trade is the measurement.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,read`; the dataset seed and
loader; **the database's budget** (`DB_CPUS=2`, `DB_MEMORY=3g`); the client's memory (`CLIENT_MEMORY=2g`);
one design, `r1_inventory_row`, the design the audit found most client-bound. **The factor varied is
`CLIENT_CPUS`**, through an `ADS_CLIENT_CPUS` override added to `infra/lib.sh` in this change (the database's
budget already had one; the client's did not, which is why this question had not been asked).

**Design.** Three runs, identical but for the client's quota:

| # | `CLIENT_CPUS` | database | machine |
|---|---|---|---|
| 1 | 1 | 2 | 4 CPUs |
| 2 | **2** (the study's published budget) | 2 | 4 CPUs |
| 3 | 3 | 2 | 4 CPUs |

Budget 2 is the reproduction of the re-run's finding; 1 and 3 bracket it. When client + database exceed the
machine's 4 CPUs (run 3), the containers share the machine under their own quotas — which is exactly the
regime the owner's "equal-total-budget" control is for, and the reason both containers' throttling is
recorded in the same cell.

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| S1 | **At `CLIENT_CPUS=2` the finding reproduces**: `r1`'s read phase is throttled for roughly half its periods (~420/634) and used ~1.8 cores. | A throttle share far from the re-run's, which would make the re-run itself unstable. |
| S2 | **The quota is the bind for the client**: at `CLIENT_CPUS=3` the client's throttling falls below 10% of its periods, i.e. the three-core quota is not what limits it. | Throttling staying above ~25% at three cores, which would point at the machine rather than the quota. |
| S3 | **A sized client is no longer flat**: at `CLIENT_CPUS=3`, `r1`'s `q01_event_availability` falls with the tier instead of sitting at one value — the 100,000-seat cell at least **2x** below the 10-seat cell (it was 37,893 against 37,748 at two cores). | The cell staying flat within ±10% across tiers even when the client is unthrottled, which would make "flat means client-bound" wrong and the flatness a property of the query. |
| S4 | **The bind moves rather than disappears**: at `CLIENT_CPUS=3`, with the same 4-CPU machine and an unchanged two-core database quota, the **database's** throttling in the same cell rises. | Database throttling unchanged or lower, which would mean the machine had spare capacity all along. |
| S5 | **Correctness gates stay green** in all three runs. | Any failed cell. |

**Scope and gaps, stated in advance.** One design, one topology, one engine, one tier set, three budgets, one
run each: this is a calibration, not a comparison of designs, and none of its runs is a study-02 result to be
compared with the published matrix. It does not measure what a *larger machine* would do, and it does not
settle the memory question — the same machine offers 2 GiB while the study's declared budgets sum to 5 GiB,
which is the owner's separate open decision.
