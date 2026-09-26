# Repeat: rollup maintenance × replacement size, three runs per tier (Study 04)

**Question.** `v2-18` registers the interaction between who maintains the parent rollup and the size of a
replacement, from **one run at each of three tiers** — 2.07x at 10 entries per installation, 4.38x at 60 and
24.20x at 240. A curve with no repetition cannot separate the trend from the scatter, and this project's
own experience says the scatter is not small: the retry control measured a **~30% day-to-day drift** on a
design that does not retry, and the `v2-13` replication found a configuration whose own repeats disagreed
by 35–115%. This tranche adds two more runs at each tier, so every level has three.

**What changes.** Nothing but the repetition. Same designs (`n3_rollup_trigger`, `n4_rollup_app`), same
topology (`pg-single`), same scale (`small`), same phases (`verify,write`), same tiers (10, 60, 240),
same three measurement repeats per configuration. Two additional runs per tier, in one session.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| R1 | **The trend reproduces.** The `n4/n3` ratio on `w05_replace_all` rises with the tier in every one of the three runs: at each tier the lowest ratio belongs to tier 10, the middle to tier 60 and the highest to tier 240. | Any run in which the ordering by tier breaks — that would mean the trend is not stable across repetitions. |
| R2 | **The width of the spread is stated, not assumed.** At each tier the three runs' ratios agree within **40%** of their median. | A spread wider than 40% at any tier, which would mean the level cannot carry a point or a tight range and the curve should be reported as bands. |
| R3 | **The tier-60 ratio stays inside `v2-13`'s replicated range plus the day term.** Three runs at tier 60 give a median within 30% of 4.38x. | A median more than 30% away, which would make the corroboration of `v2-13` in the earlier tranche look like a coincidence of one run. |
| R4 | **Correctness holds in every run.** All gates pass, with 0 audit failures and 0 rollup mismatches at every tier. | Any failure, which voids that cell regardless of its timing. |

Recorded before measurement. Evidence form: **`interaction`** (unchanged — two factors crossed); this
tranche raises the replication tier of that claim from one run per level toward run-level.

## Scope

| Field | Value |
|---|---|
| Runs added | 2 per tier, 6 in total |
| Tiers | 10, 60 and 240 entries per installation |
| Existing runs | `20260926T152426Z` (10), `20260926T152621Z` (60), `20260926T152812Z` (240) |

**Boundaries.** Three runs per level is still run-level replication on one machine, not
hardware-independent; it measures the scatter within one day on one host and says nothing about whether
the 2.07x / 4.38x / 24.20x shape survives a different machine. The tiers remain three points, so this
still cannot find a crossing — only whether the trend and its size are reproducible.
