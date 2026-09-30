---
analysis_id: 20260930T-read-q01-plan-check--deepseek-flash--2026-09-30
run_id: 20260929T165533Z (the read-family re-run whose plans are read here)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; a plan check on already-recorded evidence, with no new measurement."
analyzed_at: 2026-09-30
inputs_digest: 853e38158dac6460
inputs: 20260929T165533Z@853e38158dac6460
repo_commit: 9943eef
status: current
headline: r1's flat q01 is O(1) by construction - the plan is a primary-key lookup of a pre-aggregated counter in a 690-row table (cost 2.49, rows=1, 0.012 ms) while the baseline counts 58,927 ticket rows (cost 4,517, 6.774 ms), so the flatness is the design's point and the level is the machine's point-lookup rate, not a client ceiling or an unresolved ceiling.
---

# Analysis — `r1`'s flat `q01` is the design, and the plan says so (study 02) — deepseek-flash

*Closes the one question the 2026-09-30 calibration left open ("a ceiling **or** a genuinely O(1) design …
needs a plan check"). The answer was already in the run directory: the `explain` phase had written the plans
beside the results. **No new measurement.***

## TL;DR

- **`r1_inventory_row`'s availability read is O(1) in the tier by construction.** Its plan is
  `Index Scan using event_inventory_pkey … (cost=0.28..2.49 rows=1) (actual time=0.006..0.006 rows=1)`,
  `Buffers: shared hit=3`, `Execution Time: 0.012 ms`. The queried table holds **690 rows — one per event —
  at every tier**, because the tier scales seats *per event*; the cost model does not contain the tier at
  all.
- **The baseline's is O(tier).** `p3_precreated_cas` answers the same question with
  `SELECT COUNT(*) … FROM ticket WHERE event_id=$1 AND status='available'`: `Aggregate` over an
  `Index Only Scan` touching **58,927 rows** and **760 buffers**, `Execution Time: 6.774 ms` — **564x** the
  ledgerless-inventory design's plan.
- **So `r1`'s flat ~40–53k `q01` is two facts, not a contradiction**: the query's cost does not grow with the
  tier (design), and the level at which a two-core database serves such point lookups is the machine's
  (calibration). Neither is a client ceiling — the three-core run removed the client's throttling entirely
  and left the cell flat *and* slower.
- **This is the third time in this investigation that the answer was already recorded** (the client's cgroup
  counters, then the hazard table, now the plans). The plans are written into every run at
  `results/<run-id>/<topology>/plans/<design>.txt` by the `explain` phase, and nothing in the analysis path
  had opened them.
- **Nothing in the registry changes**: no claim rests on a read level from these runs, so this resolves a
  caveat's wording rather than a claim.

## What I looked at

`studies/02-ticket-booking/results/20260929T165533Z/pg-single/plans/r1_inventory_row.txt` and
`…/p3_precreated_cas.txt` — the `EXPLAIN (ANALYZE, BUFFERS)` output the `explain` phase wrote for the
re-run analysed in `20260930T-read-client-sizing`. Both designs' `q01` sections, and the `q02` sections for
context, are read as recorded.

## Conclusions

**The flat cell is the design's stated point, measured.** `r1_inventory_row` keeps one row per event in
`event_inventory` and decrements it, so availability is a single primary-key read whose cost cannot grow with
the number of seats. That is precisely the trade the design makes, and the read table now shows it as an
artefact: at the five tiers the cell reads 37.7–38.7k (re-run) and 45.0–52.9k (calibration, two cores) —
variation that is machine state, not tier sensitivity, because there is no tier term to be sensitive to. The
baseline, counting available tickets, falls from 25,967 to 595 ops/s over the same range, which is the tier
term doing the work.

**What that leaves as the read family's honest statement.** (1) The ordering between these designs on `q01`
is a data-model fact — one pre-aggregated row against a per-seat count — not a machine artefact, and it
reproduces: `r1` is ~1.5x `p3` at the smallest tier and ~64x at the largest. (2) The *level* of `r1`'s cell is
a machine rate: it moved 38k → 45–53k between sessions and fell to 40–44k when the client's extra core was
taken from the database, and it is not attributable to a query cost that does not exist. (3) A flat cell
still deserves this check before it is called a ceiling: the check is one `grep` in the run directory, and
the alternative reading (a genuinely constant-cost plan) is the more common one for a pre-aggregated design.

**Corrections implied.** The 2026-09-30 calibration's "ceiling *or* O(1)" is resolved to **O(1), with the
level set by the machine**; the wording in `CONTEXT.md`, `LESSONS_LEARNED.md` and the supersession ledger is
updated to the resolved form. The 2026-09-29 re-run's headline ("the floor is the client") stays retracted.

## Where I think the measurement is weak

1. **One plan per statement, sampled at one event and one tier.** The argument here is structural (a
   primary-key lookup of a 690-row table has no tier term) and does not depend on the sample, but a reader
   wanting the plan *at tier 100,000* specifically will not find it in this artifact.
2. **The plans carry no tier label.** Nothing in `plans/<design>.txt` says which tier or event was explained,
   so it cannot be tied to a row of the read table without re-running the phase; the `explain` output should
   name its sample.
3. **No latency row to match.** The plan's `Execution Time` (0.012 ms) and the read table's throughput are
   not the same quantity — one operation against a concurrent client with N connections — so the 564x here
   is a plan-cost ratio, not a throughput ratio.
4. **One engine, one host, the study's own phase set** — as for everything else in study 02.

## What I would measure next

1. **Label the explain sample** (run id, design, tier, event) in the plan file, so a plan can be cited
   against the read row it belongs to.
2. **Explain `q01` at the smallest and largest tiers** in one run, which would turn the structural argument
   into a two-point measurement for a reader who wants one.
3. **The client's read footprint** remains the open harness change for the family's levels (owner decision),
   unchanged by this.

## Disagreements with other analyses

None outstanding: it resolves the ambiguity in my own `20260930T-read-client-sizing`, agrees with
`20260929T-read-family-rerun` on the reproductions it kept, and leaves the retraction of that file's headline
in place.
