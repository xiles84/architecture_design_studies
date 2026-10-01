---
analysis_id: 20261001T-read-cells-restated--deepseek-flash--2026-10-01
run_id: six footprint-2 PostgreSQL runs — 20260930T164028Z, 20260930T164834Z, 20260930T201732Z, 20260930T203913Z, 20260930T232926Z, 20260930T233744Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; re-states the primary anchors' PostgreSQL read cells on the footprint-2 footing."
analyzed_at: 2026-10-01
inputs_digest: ecd8fd8a6935831a, dde0616ddec5a6ad, 17c1586f3e78956d, 8ca41db9b81e969d, 4706076ef2c6abc4, 1878ba0fe59000b2
inputs: 20260930T164028Z@ecd8fd8a6935831a, 20260930T164834Z@dde0616ddec5a6ad, 20260930T201732Z@17c1586f3e78956d, 20260930T203913Z@8ca41db9b81e969d, 20260930T232926Z@4706076ef2c6abc4, 20260930T233744Z@1878ba0fe59000b2
repo_commit: 2bc1bf1b2d41
status: current
headline: The primary anchors' PostgreSQL read cells now have a replacement footing - all fifteen designs measured at footprint 2 in both orders across four runs (plus two one-worker runs), every cell store-bound and order-stable, so the footprint-8 read levels and ratios may be replaced rather than caveated; the YugabyteDB read arms of both anchors remain uncovered, and no single run covers all fifteen designs, so cross-tranche orderings stay session-confounded.
---

# Analysis — the primary anchors' read cells, re-stated (study 02) — deepseek-flash

## TL;DR

- **Coverage: the PostgreSQL arm is complete, the YugabyteDB arms are not.** Both primary anchors
  (`20260913T021206Z`, `20260920T234953Z`) ran **three topologies** — `pg-single`, `yb-single`,
  `yb-cluster3` — with a client throttled for two thirds of its read periods. All **fifteen** study-02 designs
  now have a **`pg-single` footprint-2** read measurement (four runs, both orders each, plus two one-worker
  runs); the two YugabyteDB arms remain unmeasured at any footprint and stay as they were.
- **The replacement rule, with its size measured**: a footprint-8 read level is replaced by its footprint-2
  counterpart from the same design, and the difference is **1.09x–2.22x**, in rank order with how many
  periods the old client lost per design. A *ratio* between two designs also changes: `r1` over `p3` on
  `q01`@1,000 was **2.07x** at footprint 8 and is **1.16x / 1.48x** at footprint 2.
- **The quotable ordering at footprint 2** (`q01`@1,000 seats, order-stable to within 5% between orders):
  `p4` 18.4–19.3k ≈ `c4` 18.4–18.6k > `r1` 17.1–17.8k ≈ `h1` 17.5–19.2k ≈ `r2` 17.5–17.7k > `h0` 16.7–17.4k >
  `p3` 11.5–15.3k ≈ `c2`/`c3` 13.8–14.2k ≈ `c1` 13.4–14.2k > `c5` 12.2–13.3k ≈ `r3` 12.0k ≈
  `p2`/`p1` 11.4–11.7k ≈ `x1` 10.4–11.6k.
- **`q05` (tickets sold across a band) separates the structures by 4–8x**: the pre-aggregated and bounded
  designs read **9.1–15.7k ops/s** while the designs that must count pre-created per-seat rows read
  **1.9–2.2k** (`p3`, `p1`, `p2`, `x1`). This is the family's clearest design signal and it is order-stable.
- **No book prose quotes these levels** — checked, not assumed: the read figures in the draft are Study 01's
  index score, Study 05's cache reads and `v2-11`'s YugabyteDB section map, none of which is this family. So
  nothing downstream needed rewriting.
- **`v2-07` gains the two `x1` runs as corroborating anchors**: its refund-answer (`r05`) result reproduces at
  footprint 2 at **32,164 / 32,362 ops/s** against the primary's 30.3k/33.0k, i.e. within 6% across a
  2.7x client-footprint change.

## What I was looking at

Six runs, all tagged and clean, all `pg-single`, scale `small`:

| runs | designs | read workers | orders |
|---|---|---|---|
| `20260930T164028Z`, `20260930T164834Z` | `p3`, `x1`, `c1`, `r1` | 2 | A / B |
| `20260930T201732Z`, `20260930T203913Z` | the other eleven | 2 | A / B |
| `20260930T232926Z`, `20260930T233744Z` | `p3`, `x1`, `c1`, `r1` | **1** | A / B |

## Conclusions

**What may now be quoted, and how.** A read figure from this family is quotable as a **store measurement at a
stated footprint**, and its magnitude scales as `workers ÷ mean latency` between footprints. The primary
anchors' own read figures are not quotable — their inflation is design-dependent (1.09x–2.22x, measured
against the old client's per-design throttling) — and neither are their read *ratios*, because the inflation
was not uniform across designs (`r1`/`p3` halved when re-measured). Substituting the footprint-2 figures is
therefore the correct repair, and it is what the limits now instruct.

**Tranches are not interchangeable sessions.** The four-design set and the eleven-design set were measured in
different runs, hours apart; their *within-tranche* orderings are order-stable, but a comparison spanning the
two (`r1` 17.8k against `p4` 18.4k, say) carries the session term those tranches showed elsewhere (up to 10% on
`q01`, and 33% in an earlier session). A single footprint-2 run covering all fifteen designs has not been
made, and until it is, cross-set comparisons should be stated as such. This is the honest limit of "complete
coverage": complete *per design*, not joint.

**The two YugabyteDB arms are the remaining gap, and the reason is resource, not method.** Both anchors
measured `yb-single` and `yb-cluster3` read phases with the same throttled client, so their YB read figures
carry the same defect and **no replacement exists** — the 2 GiB podman machine governs whether that half can
be measured at all, and the owner has fixed the machine. Recording that as a coverage gap is the correct
outcome; re-labelling the old YB read figures as quotable is not.

**Prediction:** none. This is an integration analysis over already-published runs; no new measurement was
taken and no prediction was registered, so none is scored. The `v2-07` reproduction and the coverage audit are
the substantive checks it performs.

## Where I think the measurement is weak

1. **Cross-set comparability** — as above: fifteen designs covered, but in two sessions.
2. **The `pg-single` arm only.** Both anchors' `yb-single` and `yb-cluster3` read arms remain not quotable, with
   no replacement measurable on this machine.
3. **One run per order per set.** The order agreement is 0.91–1.03, but a level quoted from a single run still
   carries its session.
4. **The primary anchors also ran non-read phases** (race, load, hold), whose figures this analysis does not
   touch: the race and load bases of `v2-06`, `v2-08`, `v2-09` and `v2-gap-01` are unaffected by the footprint
   change and were not re-stated.
5. **`v2-11`'s YugabyteDB read claim** (section document 4.5–11x) is a different environment at a different
   footprint and is untouched here; whether its client was throttled is not something this analysis checked.

## What I would measure next

1. **One footprint-2 run covering all fifteen designs**, which would turn "complete per design" into "complete
   and jointly comparable" for the PostgreSQL arm — the single most useful next read measurement.
2. **The `v2-11` YugabyteDB read claim's client**, checked against the same hazard table, since that claim
   quotes a read *ratio* from a footprint nobody has audited.
3. **Nothing further on the PostgreSQL level question** — it is answered.

## Disagreements with other analyses

None. It depends on `20260930T-read-client-footprint` (the footprint and its mechanism),
`20260930T-read-family-footprint2` and `20260930T-read-family-11-footprint2` (the measurements), and
`20260930T-read-one-worker` (which established that every cell is store-bound and withdrew the `q02`/`q05`
caveat). It re-states, and does not contest, `20260920T234953Z--deepseek-flash--2026-09-21`'s read tables.
