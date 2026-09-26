# Replication: the seat map on YugabyteDB (study 03, `v2-11`)

**Question.** `v2-11-section-seat-map` claims that on YugabyteDB a stored section document answers the
100k-seat map far faster than one row per seat, while the variant that shards by section collapses point
lookups; and that on PostgreSQL no difference can be attributed to the design. It rests on one run, one
measurement per configuration, from 2026-09-21 — the same footing as `v2-06` had before its replication,
and study 03 has never been re-run at all (its runner could not even build an image from a WSL path until
2026-09-26). This tranche re-measures the comparison on one topology with three repeats per cell.

**What is held constant.** Designs `l1_claim_rows`, `l2_section_document` and `l3_section_sharded`; topology
`yb-single`; scale `small`; phases `verify,read`; dataset seed 42. Only the session differs.

**Baseline being tested against.** From `run/03-reserved-seating/20260921T231328Z`, YugabyteDB 1 node:
seat map of a section at 100k seats — L1 **149** ops/s, L2 **3.1k**; seats available per section at 100k —
L1 **62.2**, L2 **157**, L3 **15.9**; seats available per section at 10 seats — L3 **28.3**.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| U1 | **The document's advantage on the 100k seat map reproduces.** `l2_section_document` reads the 100k-seat map at least **4x** faster than `l1_claim_rows` on `yb-single`. | A ratio below 4x, which would mean the 20.8x of the original run was a session or load artefact rather than a property of the two designs. |
| U2 | **The sharded variant still collapses.** `l3_section_sharded` records under **100 ops/s** on at least one "seats available per section" row while L1 and L2 stay above it. | L3 staying above 100 ops/s on every such row, which would mean the collapse was specific to that run. |
| U3 | **The magnitudes move by no more than a session term.** Every cell falls within **±40%** of the 2026-09-21 value for the same design and question. | Any cell outside ±40%, which would mean `v2-11`'s figures are not reproducible in size and should be re-scoped to directions. |
| U4 | **The correctness gate stays green.** All designs pass verification and no cell is dropped for a failed check. | Any failure, which voids that cell. |

Recorded before measurement. Evidence form stays `package` with attribution withheld, as `v2-11` is a
comparison of two whole designs rather than a one-decision pair; this tranche raises its replication tier.

**Coverage gap, stated in advance.** The PostgreSQL half of the claim is **not** re-measured here, and
neither is the 3-node topology, so a pass on U1–U2 says nothing about the engine-dependence that makes the
claim interesting. It also does not re-run the deliberately wrong controls (`s0_check_then_hold_rc`,
`k0_naive_confirm`), so the gate here is a pass-gate rather than a fired control.

## Scope

| Field | Value |
|---|---|
| Runs | one, `-trials 3` per cell |
| Topology | `yb-single` only |
| Existing evidence | `run/03-reserved-seating/20260921T231328Z` (primary anchor of `v2-11`) |

**Boundaries.** One session, one engine, one topology, one scale. The claim's most interesting content is
that the difference is *engine-dependent*; a single-engine replication cannot test that, only whether this
engine's half resolves.
