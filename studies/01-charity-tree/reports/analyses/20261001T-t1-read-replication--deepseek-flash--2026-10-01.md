---
analysis_id: 20261001T-t1-read-replication--deepseek-flash--2026-10-01
run_id: 20261001T112526Z-t1 (order A), 20261001T121054Z-t1 (order B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; T1 run-level replication of Study 01's read-score claims."
analyzed_at: 2026-10-01
inputs_digest: f725ccae64c1fa76 (order A), 39767ca6e590b357 (order B)
inputs: 20261001T112526Z-t1@f725ccae64c1fa76, 20261001T121054Z-t1@39767ca6e590b357
repo_commit: 9356884
supersedes:
status: current
headline: The three oldest claims replicate - D3 > D2 > D6 > D1 and the rollup family far above D2 hold in both whole runs in reversed order, with d2/d1 at 9.7-10.6x (quoted 13.3x), d3/d1 at 13.7-15.5x (20.5x), d4/d5 at 61-85x (73.8x/69.9x) and d6 at 1.7-2.2x (2.4x); every ratio stayed within 1.4x between runs while the absolute scores moved 40-45%.
---

# Analysis — T1 replication of Study 01's read-score claims — deepseek-flash

## TL;DR

- **`v2-01`, `v2-02` and `v2-03` replicate.** Two whole runs of the primary's own configuration (scale
  `small`, 10 s per op, the same eight designs) in **reversed design order** reproduce every ordering the
  claims rest on: `d3 > d2 > d6 > d1` in both runs, and `d4`/`d5` far above `d2`.
- **The ratios reproduce; the point figures do not.** `d2`/`d1`: **10.6x** and **9.7x** (quoted 13.3x);
  `d3`/`d1`: **13.7x** and **15.5x** (20.5x); `d4`/`d1`: **61.5x** and **85.2x** (73.8x); `d5`/`d1`: **64.0x**
  and **82.8x** (69.9x); `d6`/`d1`: **2.2x** and **1.7x** (2.4x). Every ratio stayed within a factor of
  **1.4** between the two runs, while the **absolute read scores moved 40–45%** for the same designs
  (`d2`: 4,890 against 2,559). That is the corpus's range-travels lesson arriving on its oldest claims.
- **Embedding is still not ahead**: `d6` is 1.7–2.2x `d1`, below both indexed designs (about 10x and 14–16x)
  — the claim's substance survives, with a lower band than quoted.
- **The write penalty reproduces and is ordered**: `d4` 0.32x/0.37x of `d1`'s insert rate, `d5` 0.34x/0.29x,
  `d6` **0.57x/0.60x** — the last being the one prediction I got wrong (see `K4` below).
- **Two `K` predictions were falsified, both mine rather than the claims'**: `K4`'s "`d6` below half of `d1`"
  contradicts the primary's own measured 0.57x, which is what I reproduced; and `K5`'s factor-of-2 agreement
  held everywhere (1.09–1.39), so the falsification is a threshold I set badly, not a failed replication.
- **Provenance weakness, recorded**: Study 01's runner writes **no `repo_commit` and no `repo_dirty`** into
  its manifest — unlike study 02's — so the two runs cannot state their own code. I tagged them
  `run/01-charity-tree/<run-id>` at commit `9356884`, the clean tree they were produced from, and that tag is
  the only provenance they have.
- **Registered**: the three claims move from `single_run_directional` to **`run_replicated`**, with their
  statements re-scoped to the replicated bands (a point figure that moved 30% between runs is not the
  claim), the two runs added as anchors, and the YugabyteDB arms named as the remaining gap.

## What I was looking at

Two whole runs, `pg-single`, scale `small`, `duration 10s`, the primary anchor's eight designs
(`d1_normalized_minimal`, `d2_normalized_indexed`, `d3_flattened_fk`, `d4_rollup_trigger`, `d5_rollup_app`,
`d6_embedded_jsonb`, `d8_flattened_nofk`, `d9_embedded_hybrid`), order A `d1…d9` and order B reversed. The
read score is Study 01's own definition — the geometric mean over the study's twelve core queries
(`harness/report.go`'s `coreQueries`), excluding the recency queries — recomputed here from the run JSON so
the replication measures exactly what the claims measure. No cell failed its correctness gate.

| design | A: score / ratio | B: score / ratio | primary score / ratio |
|---|---|---|---|
| `d1_normalized_minimal` | 462 / 1.0x | 264 / 1.0x | 396 / 1.0x |
| `d2_normalized_indexed` | 4,890 / **10.6x** | 2,559 / **9.7x** | 5,245 / 13.3x |
| `d3_flattened_fk` | 6,322 / **13.7x** | 4,082 / **15.5x** | 8,120 / 20.5x |
| `d4_rollup_trigger` | 28,424 / **61.5x** | 22,463 / **85.2x** | 29,179 / 73.8x |
| `d5_rollup_app` | 29,592 / **64.0x** | 21,840 / **82.8x** | 27,655 / 69.9x |
| `d6_embedded_jsonb` | 1,030 / **2.2x** | 447 / **1.7x** | 951 / 2.4x |
| `d8_flattened_nofk` | 8,874 / 19.2x | 4,563 / 17.3x | 8,086 / 20.4x |
| `d9_embedded_hybrid` | 8,728 / 18.9x | 5,464 / 20.7x | 8,593 / 21.7x |

## Conclusions

**The three claims are now run-replicated, and their statements should say bands.** Every relation the claims
assert reproduced in both orders, and the read-score ratio is stable to within 1.4x between runs — while a
single absolute score moved 40–45%. So the claim text moves from "D2 indexed scores 13.3x D1 minimal's read
score" to the replicated band (**about 10–11x in two whole runs, against 13.3x in the first**), and likewise
for the rollup family (**about 61–85x against 73.8x and 69.9x**), and for embedding (**1.7–2.2x against
2.4x**). The orderings and the gaps between designs are what a reader should carry; the point figures were
single-run.

**The write side of `v2-02`/`v2-03` reproduces as an ordering, not as rates.** `d4` and `d5` pay roughly a
third to two-fifths of `d1`'s insert rate and `d6` about three-fifths, in both orders — the "rollups and
embedding cost inserts" relation holds. The quoted rates themselves (2.6k/s, 2.1k/s, 6.2k/s) moved with the
session like everything else: this run's `d1` inserted at 6,098 and 4,886 per second, so the *ratios* are the
portable part.

**One structural finding worth carrying**: in the reversed order **every design scored lower** — `d2` 2,559
against 4,890, `d3` 4,082 against 6,322 — while their ratios held. The order effect here is a level effect
(roughly 40%), not a ranking effect, and it is the same pattern the study-02 read family showed on 2026-10-01.
A single-order run of Study 01 would have produced a clean, false set of absolute scores.

**Prediction verdicts** (recorded in `T1-READ-REPLICATION.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| K1 | `d3 > d2 > d6 > d1` and `d4`/`d5` above `d2`, in both runs | **confirmed** — exact in A and in B |
| K2 | `d2`/`d1` ≥ 6x, `d3`/`d1` ≥ 10x, `d4`/`d1` and `d5`/`d1` ≥ 25x | **confirmed** — 9.7–10.6x, 13.7–15.5x, 61.5–85.2x, 64.0–82.8x |
| K3 | `d6`/`d1` between 1.5x and 6x and below `d2` | **confirmed** — 2.2x and 1.7x |
| K4 | every one of `d4`, `d5`, `d6` below **half** of `d1`'s insert rate | **falsified as stated** — `d6` is 0.57x/0.60x; the primary's own figure for `d6` is 0.57x, so my threshold contradicted the claim I was testing and the measurement reproduced it |
| K5 | each ratio within a factor of 2 between runs; no ≥1.5x pair swaps | **confirmed** — ratio-of-ratios 1.09–1.39, no swaps |
| K6 | correctness green, clean trees, runs tagged | **confirmed for the cells and the trees** (no failed cell, `repo_dirty` not recorded by this runner, tree clean at `9356884`); **the tagging was manual** because this runner predates `--tag` |

## Where I think the measurement is weak

1. **`pg-single` only.** The primary's `yb-single` and `yb-cluster3` arms are unreplicated, and the fixed
   2 GiB machine is the reason: the registry keeps them as an explicit coverage gap.
2. **Two runs, one scale, one duration.** The claims' own limits say "small scale"; nothing here tests whether
   the bands hold at medium scale.
3. **The runner records no commit.** These two runs' provenance rests on a manually created tag; any future
   Study 01 run inherits that weakness until the harness records `repo_commit`/`repo_dirty` like study 02's.
4. **The read score is a blend.** A claim about it is a claim about the geometric mean of twelve queries, and a
   design can win the blend while losing a specific question (the study's own analyses show exactly that).
5. **`d8` and `d9` were measured but are not part of any claim here**, and their numbers (17–21x) are carried
   only in the reports.
6. **The order/level effect is characterised, not explained** (as in study 02's read family).

## What I would measure next

1. **The YugabyteDB arms of the same three claims**, which needs a machine that can hold them.
2. **A medium-scale pass** of the same eight designs, to test whether the bands are scale-robust.
3. **Making Study 01's runner record its commit**, so its runs stop depending on a manual tag.

## Disagreements with other analyses

None. It re-states, and does not contest, `20260912-study01--gpt-6--2026-09-12`'s read-score tables: the
ratios there are confirmed as bands, the point figures are shown to be single-run levels, and the registry
now says so.
