# `v2-11` replication, second pass: the pair in reverse order (study 03)

**Question.** The first replication of `v2-11` (run `20260926T234227Z`) ran the designs in the order
`l1, l2, l3` and confirmed the claim's direction — the stored document read the 100k-seat map 27.8x faster
than one row per seat — but the registry lint refuses `run_replicated` for a fixed order, because on this
host a design's position in the session has already been shown to move its own throughput by 58–76%
(study 02, `v2-06`, 2026-09-26). This pass runs the **same two designs in the opposite order**, so the
comparison is AB/BA and the claim can be registered if it survives.

**Order, and why `l3` goes first.** `l3_section_sharded, l2_section_document, l1_claim_rows`. Two things
are being tested at once:

- `l2` before `l1` reverses the pair that the claim rests on.
- `l3` failed its **load** — three consecutive connection terminations, SQLSTATE 57P01 — when it ran *last*.
  Putting it first with a fresh host separates two explanations that the previous run could not: a design
  whose load cannot complete at all, and one that could not complete after two other designs had loaded.

**What is held constant.** Designs, topology (`yb-single`), scale (`small`), phases (`verify,read`), three
trials per cell, dataset seed 42. Only the order and the session differ from `run 1`.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| V1 | **The reversal reproduces the direction.** `l2_section_document` still reads the 100k-seat map at least **4x** faster than `l1_claim_rows`. | A multiple below 4x, which would mean the 27.8x of run 1 (and the 20.8x of 2026-09-21) depends on the order the designs were run in. |
| V2 | **The two orders agree.** Run 1's multiple and this one's fall within **±40%** of each other. | More than 40% apart, which would put the claim's multiple in the same position-effect band study 02 measured. |
| V3 | **`l3`'s failure was positional.** Run first, `l3_section_sharded` completes its load and reaches verification. | It fails again when run first, which would make the failure intrinsic to the design or to this host's capacity rather than to its position — and would mean the sharded half of the claim cannot be tested here at all. |
| V4 | **The correctness gate stays green** for whichever cells reach it. | Any cell that reaches verification and fails it. |

Recorded before measurement. If V1–V2 hold, `v2-11` is registered as `run_replicated` with `order: abba` and
the pair's two anchors; if they do not, the claim stays `single_run_directional` and this run becomes the
evidence for a position-effect limitation instead.

## Scope

| Field | Value |
|---|---|
| Runs | one, `-trials 3` per cell, order `l3, l2, l1` |
| Compared with | `20260926T234227Z` (order `l1, l2, l3`) and `20260921T231328Z` |
| Boundary | One engine, one topology, one scale, one host. The PostgreSQL half of the claim remains unmeasured, so engine-dependence is still untested. |
