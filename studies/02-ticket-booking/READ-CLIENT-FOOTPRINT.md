# A smaller read client footprint (study 02)

**Owner decision, 2026-09-30.** The machine is fixed: the project assumes the computer it has and will not
spend on more resources. The read phase's levels are therefore to be repaired by **shrinking the client's
footprint**, not by re-allocating cores and not by a larger machine. This tranche calibrates the one lever
that does that.

**The lever, measured first.** The read phase runs `-conns` concurrent workers, each issuing one query at a
time on its own connection; the published default is **8** (`run-study.sh`, `CONNS`). At two client cores
that footprint saturates the quota: on 2026-09-29/30 the client spent **111.7 s of CPU in a 63.3 s window**
and lost **422 of 633 periods** to throttling, while the same runs showed the reads getting *slower* when the
client was given a third core (taken from the database). So the demand, not the quota, is what has to fall.

**What is held constant.** Topology `pg-single`; scale `small`; phases `verify,explain,read`; the published
budgets (`DB_CPUS=2`, `CLIENT_CPUS=2`); duration 5 s per measurement; designs `r1_inventory_row` (the cell
whose flatness is now explained as O(1) by its plan) and `p3_precreated_cas` (the O(tier) baseline). **The
factor varied is `-conns`**: the read workers.

**Design.** Three whole runs, identical but for the read workers:

| # | `-conns` | note |
|---|---|---|
| 1 | **8** | the published footprint; must reproduce the saturation |
| 2 | 4 | half the demand |
| 3 | 2 | a quarter of the demand |

## Prediction

Recorded before measurement.

| Id | Prediction | Falsifier |
|---|---|---|
| F1 | **Throttling falls with the footprint**: run 1 reproduces the ~420/633 throttle share, run 2 is well below it, run 3 is **0 or near 0**. | Run 3 still throttled above ~10% of periods, which would mean the client's per-operation cost is so high that even two workers exceed two cores, and the fix must be in the harness's per-op work rather than in concurrency. |
| F2 | **Throughput falls with the footprint and the cell becomes latency-bound**: at the unthrottled footprints, `ops_per_sec × p50` is close to the worker count (Little's law), against a client-bound cell at 8 workers where the product exceeds it. | The product staying far from the worker count at 2 workers, which would mean something other than the workers is limiting. |
| F3 | **The ordering is what travels**: `r1` stays far ahead of `p3` on `q01` at every footprint (the design difference is a data-model fact), while the absolute levels move by 2x or more. | The ordering inverting or collapsing at the smaller footprints. |
| F4 | **The database's own throttling falls too**, since the client stops competing for the same four CPUs. | Database throttling unchanged or higher. |
| F5 | **Correctness gates stay green** at every footprint. | Any failed cell. |

**What a pass decides.** Which read-worker count leaves the client unthrottled and the cell latency-bound —
the configuration the read family should be measured at from now on, as a labelled change to the read phase
only (the write, race and hold phases keep the published `-conns 8`, so their published numbers stay
comparable). It does not re-measure the family; that is the tranche after this one.

**Scope and gaps, stated in advance.** Two designs of fifteen, one topology, one engine, one scale, one run
per footprint, one host — and the machine's 4 CPUs / 2 GiB are now a stated constraint of the project rather
than an open question. Nothing here changes the published matrix.
