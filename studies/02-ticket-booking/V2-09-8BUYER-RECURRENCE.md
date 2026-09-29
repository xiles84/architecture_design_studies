# Does the 8-buyer ledger collapse recur? (study 02, `v2-09`)

**Question.** In the 2026-09-28 buyer-count sweep, the 8-buyer `p3`-first cell's **ledger arm** (`x1`) fell
from ~4,840 to ~970 ops/s between its second and third trials and **stayed** there for the rest of the cell,
while the same level's other order sat flat at ~5,000. The sweep reported the cell rather than discarding
it, and called the cause unresolved. This tranche asks whether that collapse **recurs**, and if so whether
it is a property of the ledger arm, of a repeated fresh load, or of the session.

**What the existing artifact already says** (from `20260928T014723Z`, so the prediction below is informed
rather than blind):

- The change is **persistent, not a transient stall**: trials 1–2 read 1.05 s wall and p50 1.25 ms; trials
  3–5 read 5.2 s wall and p50 6.7 ms. Five times the wall time per operation, held for three trials.
- It is **not host-CPU exhaustion of the database**: the cell's `cpu.stat` shows `nr_throttled` rising 32 →
  68 but `throttled_usec` rising only ~4.2 ms across the whole cell.
- It appears in the **ledger arm only**: the same run's CAS arm read 5,454 / 5,143 / 5,496 / 5,154 / 4,492 —
  flat within 18%.
- It did **not** appear in the other 8-buyer order (`x1` first), whose ledger arm was flat ~5,000, so at
  least one of {order, session} matters.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **8 buyers per event**; **tier 1,000 only**; `-race-trials 5`; dataset seed
42. The factor varied between runs is the **design order**.

**Design.** Six whole runs, alternating order — `A`,`B`,`B`,`A`,`A`,`B`, where `A` is `p3`,`x1` and `B` is
`x1`,`p3` — giving twelve design-cells (six per arm) at one buyer count. Six runs rather than two because
the event is defined by recurrence: one occurrence in three runs says little either way.

## Prediction

Recorded before measurement. A **collapse** is a drop of **≥2x** between consecutive trials of one cell,
persisting in the remaining trials (not a single-trial spike).

| Id | Prediction | Falsifier |
|---|---|---|
| C1 | **The collapse recurs.** At least **one** of the twelve cells shows a collapse. | **No** cell collapses, which makes the 2026-09-28 event a one-off and the earlier cell safe to keep as a datum. |
| C2 | **It is a ledger-arm effect.** Every collapsing cell is an `x1` (ledger) cell; no pure CAS cell collapses. | A CAS cell collapsing, which would make it a harness/session effect rather than the ledger's. |
| C3 | **It is persistent within the cell**, not a spike: the trials after the drop stay within 30% of the dropped level. | A cell that drops and recovers on the next trial. |
| C4 | **It is more likely in the CAS-first order** (`A`), where the ledger arm runs second and meets the state three fresh loads have already left. | Collapses appearing only in `B`, or equally in both orders. |
| C5 | **Correctness stays green** in all six runs: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors. | Any failure, which voids that run's cells. |

**What a pass can and cannot decide.** It decides whether the event recurs and which arm it belongs to. It
does **not** find the mechanism: the harness records the database's cgroup counters but not the client's,
and no trial-level clock or plan is captured, so a client-side stall and a database-side state change are
not separable from these runs. It says nothing about other buyer counts or tiers.

**Coverage gap, stated in advance.** `pg-single` only, one host, 8 buyers, tier 1,000 only; six runs is
enough for recurrence, not for a rate. The 2 GiB podman-machine ceiling applies.

## Scope

| Field | Value |
|---|---|
| Runs | six, `pg-single`, two designs, `--race-buyers 8 -race-trials 5 -race-tiers 1000` |
| Compared with | `20260928T014723Z` (the collapsed cell), `20260928T014640Z` (the flat one) |
| Unit of evidence | the per-trial series of each design-cell, not its median |
| Boundary | One host, one engine, one buyer count, one tier |
