# A reference cell drawn from the operation being compared (Study 04)

**Question.** The last two tranches established that this host drifts through a session — enough to
double the apparent scatter — and that a reference cell measured before and after caught it (+44%). But
that reference was the *contention* phase on a different design pair, so it bounded the session's
movement without measuring the drift of the quantity the interaction claim actually compares. This
tranche repeats the instrument on the **same operation**, and takes the reading three times through the
tranche rather than only at its ends, so the drift becomes something the preceding and following cells can
be read against.

**Design.** The reference is the interaction cell itself at the middle size — designs
`n3_rollup_trigger` and `n4_rollup_app`, `pg-single`, `small`, `verify,write`, tier 60, operation
`w05_replace_all` — measured three times, interleaved with two comparison cells:

| Position | Cell |
|---|---|
| 1 | reference, tier 60 |
| 2 | comparison, tier 10 |
| 3 | reference, tier 60 |
| 4 | comparison, tier 240 |
| 5 | reference, tier 60 |

Five runs in one session. The reference's three readings are the instrument; the two comparison cells are
there to see whether a single reading of a cell can be placed against the reference that surrounds it.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| S1 | **The same-operation instrument works.** The three tier-60 reference readings give a `n4/n3` ratio, and their spread is measured; the *ratio's* spread is smaller than the spread of either design's absolute throughput, because a shared session term cancels in the ratio and only the differential part remains. | The ratio's spread is as large as or larger than the absolutes' spreads, which would mean the two designs do not share the session term and a ratio is no more portable than a throughput. |
| S2 | **A single cell can be placed.** The two comparison cells fall inside the corresponding bands the interleaved tranche measured (tier 10 ≈ 1.3–3.2x, tier 240 ≈ 10–24x). | Either cell falls outside its band, which would mean one reading plus a reference is not enough to locate a cell. |
| S3 | **The reference readings move monotonically with session time**, as the drift did in the previous two tranches. | A non-monotone sequence, which would mean the drift is not a simple time trend and the instrument needs more than endpoints. |
| S4 | **Correctness holds in every run.** 0 audit failures, 0 rollup mismatches. | Any failure, which voids that cell. |

Recorded before measurement. Evidence form: `interaction` (the claim is unchanged); this tranche adds an
instrument, not a new claim.

## Scope

| Field | Value |
|---|---|
| Runs | 5: three reference readings, two comparison cells |
| Reference | tier 60, the same designs, operation and phases as the comparison |
| Prior evidence | `study-04/v2-interaction-repeat`, `study-04/v2-interleaved-sweep` |

**Boundaries.** One session, one host, three reference readings and one reading of each comparison cell.
This cannot establish that the instrument is reliable — only whether it is consistent within one session —
and a single reading of a cell remains a single reading no matter how well the surrounding session is
characterised.
