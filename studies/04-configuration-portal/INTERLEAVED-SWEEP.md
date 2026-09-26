# Interleaved sweep with a reference cell (Study 04)

**Question.** The repeat of the maintenance × replacement-size interaction found a trend that reproduces
but magnitudes that spread up to 55%, and the spread turned out to be **drift by time**: throughput fell
for both designs through the session, so the repetition index was confounded with the session's state.
This tranche tests the fix the lesson prescribes — vary the thing whose effect you want to exclude — and
adds the instrument that turns the drift into a number.

**Design.**

1. **A reference cell, measured before and after everything else.** The reference is the contention phase
   at 16 writers for `c2_pessimistic_lock`, a design that never retries and so is insensitive to retry
   policy. If the session is stable, the two readings agree; the difference between them *is* the
   session's drift.
2. **A Latin square over the three sizes.** Nine runs in three sweeps, each size occupying each position
   in the sweep exactly once:

   | Sweep | 1st | 2nd | 3rd |
   |---|---|---|---|
   | A | 10 | 60 | 240 |
   | B | 60 | 240 | 10 |
   | C | 240 | 10 | 60 |

   So "position in the session" and "size" are no longer the same thing, which the sequential repeat
   could not manage.

**What is held constant.** Designs `n3_rollup_trigger` and `n4_rollup_app`, topology `pg-single`, scale
`small`, phases `verify,write`, three measurement repeats per configuration, operation `w05_replace_all`,
one session. Nothing else changes.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| L1 | **The trend survives interleaving.** Within every sweep the ratio rises with the size: tier 10 lowest, tier 240 highest. | Any sweep in which the ordering breaks, which would mean the trend depended on the sequential design rather than on the size. |
| L2 | **The session drift is measurable and modest.** The reference cell's `c2` throughput differs by less than **20%** between the reading taken before the first sweep and the reading taken after the last. | A difference of 20% or more, which would mean this session drifted as much as the previous one and the reference instrument has caught it. |
| L3 | **Interleaving reduces the apparent spread.** At each size, the three interleaved runs agree *more closely* than the sequential tranche did (15.1% at tier 10, 54.7% at 60, 51.9% at 240). | A spread at any size as wide as or wider than the sequential one, which would mean the spread is not positional and the sequential result was not an artefact of ordering. |
| L4 | **Correctness holds in every run.** All gates pass, 0 audit failures, 0 rollup mismatches. | Any failure, which voids that cell. |

Recorded before measurement. Evidence form: `interaction` (unchanged); the tranche adds an
order-controlled replication of the same claim and an instrument for the session term.

## Scope

| Field | Value |
|---|---|
| Runs | 9 interaction runs (3 sweeps × 3 sizes) plus 2 reference runs |
| Reference cell | `c2_pessimistic_lock`, `contention` at 16 writers, phases `verify,contention` |
| Existing evidence | the sequential tranche, `study-04/v2-interaction-repeat` |

**Boundaries.** Still one host, one day, one scale, three sizes. A Latin square removes the *confound*
between size and position but cannot remove drift itself — it distributes it. The reference cell measures
the drift and lets the two runs at the ends be compared; it does not make the session stable.
