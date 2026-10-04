---
analysis_id: 20261004T-select-join-cte--deepseek-flash--2026-10-04
run_id: select-join-pg1/pg2/pg3/yb1/yb2/yb3-20261004
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek Flash via opencode CLI (LOW executor session)"
analyzed_at: 2026-10-04
inputs_digest: "pg1=070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac; pg2=112577024534a9e91932f78d3ba8b9e530e0caea0aeefdf851737e578be7bab6; pg3=05b2d7a99e38c77c84baf413a769f531179682921a2419bf73da89153b2d00ac; yb1=f650acde1ebb5541080360bca3e467863611e00322d3f580cefdd36ba90e0281; yb2=25e5f7bf01de5149bed501f089eebe900d7b462c2927eedaee302436b664d198; yb3=2d8394b9cc2a9a0f649367486f9d994fcfc1538e3e7d151bc75a1e8febe82e46"
repo_commit: 7e148174b5211225f6b4f243fc02e0de417a0586
supersedes:
status: current
headline: "Six inner-join spellings are plan-identical within an engine; forcing MATERIALIZED on a join-containing CTE costs up to 88x on PostgreSQL but only ~1.6x on YugabyteDB, and the reuse crossing exists only on PostgreSQL."
---

# Analysis — join/filter spellings and CTE materialization — deepseek-flash

> LOW-executor interpretation. It is my own reading of the generated reports and does
> **not** authorize book promotion; a HIGH review (`task-20261002T131729Z-select-synthesis`,
> then `task-20261004T142623Z-select-book-high`) owns that. Generated per-run reports are
> `results/select-join-{pg,yb}{1,2,3}-20261004/join.md`; the replicated report is
> `reports/20261004T-select-join-replicated.md`.

## TL;DR

- **Pick a spelling for readability, not for speed — unless you force materialization.**
  On both engines the six inner-join spellings and the filtered derived table produced the
  same plan structure as `JOIN … WHERE`, and every within-run ratio between them was
  0.86–1.52x (PostgreSQL) and 0.91–1.04x (YugabyteDB). No universal winner (P1/P2 not
  falsified, and partly the useful "no resolved difference" answer).
- **Forcing `MATERIALIZED` is only expensive when the CTE contains the join.** Materializing
  a filtered-input CTE is near-neutral (0.83–1.11x PG, 0.97–1.08x YB). Materializing a
  *joined-output* CTE with the predicate applied outside cost **6.8x–88x on PostgreSQL**
  at 0.01–50% selectivity, but only **1.02x–1.61x on YugabyteDB**. The syntax is portable;
  its cost is not (P3).
- **The reuse crossing is PostgreSQL-only.** On PostgreSQL a multi-reference CTE is
  materialized by default; `NOT MATERIALIZED` inlines and is *faster* at 1% consumer
  selection (0.29–0.74x) but *slower* at 100% (2.06–4.29x at 16 references) — a bounded
  crossing over {1,2,4,8,16}. On YugabyteDB `NOT MATERIALIZED` was uniformly 8.9–11.2x
  slower at 16 references at both consumer selections, with no crossing (P4 engine-specific).
- **Where the two engines differ is plan shape, not spelling.** PostgreSQL uses
  `facts_rank_idx` for the 1% predicate; YugabyteDB chose a Seq Scan with
  `Storage Filter: (rank <= 2000)`. That single fact explains the two headline contrasts.

## What I was looking at

Six whole single-node runs, all correctness-gated: 3× PostgreSQL 17.11 and 3× YugabyteDB
2025.2.6, `seed 42`, 200,000 fact rows / 2,000 customers, 5 trials per cell, a seeded
per-run cell permutation with interleaved samples and a first-cell-repeat position control.
**91/91 cells were valid in every run, 0 invalid/timeout/error**, and both negative controls
fired where they could (they are recorded per run). I deliberately ignored the 3-node YB,
placement, blended-write and row-width arms: they are recorded as gaps below and I did not
measure them.

## Conclusions

**The six spellings are one bag and, in practice, one plan.** `j02_join_on`,
`j03_filtered_derived`, `j04_cte_filtered`, `j05_cte_joined_outer` and `j06_comma_join`
all had `same plan structure = true` against `j01_join_where` at every threshold on both
engines, and their ratios sat inside the within-run drift (see the position control, which
ranged 0.61–1.35 across the six runs). A team choosing between them should therefore choose
on readability, portability and optimizer control, not on a measured speed ranking. The one
PostgreSQL cell that looks like an exception — `filtered_derived` at 1.52x on a 0.42 ms cell
— shares the plan and is inside that run's 0.80 position-control drift, so I do not treat it
as resolved.

**`MATERIALIZED` costs what it forbids the planner to do.** `WITH e AS MATERIALIZED
(SELECT * FROM facts WHERE rank <= :t)` still filters before joining and stays within
0.83–1.11x of the default on PostgreSQL. But `WITH e AS MATERIALIZED (SELECT … FROM
customers c JOIN facts f …) … WHERE rank <= :t` computes the *whole* join (≈200,000 rows)
before applying the predicate: its absolute median is ~37–48 ms at every threshold from
0.01% to 90%, so the ratio against the inlined spelling explodes as the predicate gets more
selective — 88x at 0.01%, 27x at 1%, 6.8x at 10%, 1.4x at 50–90%. The same change on
YugabyteDB costs only 1.61x at 0.01% and 1.02–1.31x at the wider levels, because
YugabyteDB already scans the fact table and applies a storage filter for *every* spelling;
materializing the join adds work but removes nothing that was being used. This is the
clearest example in the family of a boundary being a property of the engine's plan, not of
the SQL keyword.

**The reuse decision is even less portable.** On PostgreSQL, `r_default` at 4 references had
the same plan structure as `r_mat` — the default already materializes a multi-referenced CTE
— so `MATERIALIZED` was redundant. `NOT MATERIALIZED` inlined and the reference count became
a multiplier on whichever work survived inlining: at 1% consumer selection the inlined plan
computes the token for ~1% of the CTE output, so it was faster (0.29–0.74x at references
2–16); at 100% it recomputes the token once per reference, so it grew to 2.06x/2.86x/3.05x/
4.29x at 2/4/8/16 references. That is a bounded crossing in the second consumer-selection
dimension, exactly as P4 predicted — **on PostgreSQL**. On YugabyteDB `NOT MATERIALIZED` at
16 references was 490 ms/585 ms against 55 ms/52 ms (8.9x/11.2x) at 1%/100% consumer
selection, because each inlined reference re-scans the fact table and the consumer filter
cannot reduce that scan. On that engine there is no consumer-selection level where inlining
wins, so P4's reversal is **not reproduced** in this fixture and must not be stated as a
portable rule.

**Practical reading.** Prefer the plain join; use a CTE for clarity; do not add
`MATERIALIZED` "to help" unless you have measured the plan, because it can turn a selective
filter into a full join. For repeated expensive work, the default is usually already
materialized, and forcing `NOT MATERIALIZED` is a decision that helps only where the engine
can push a selective filter into the inlined copy — an index-visible regime on PostgreSQL,
and not at all in the scans YugabyteDB chose here.

## Where I think the measurement is weak

- **One host, single node, shared cores.** No real network, no independent hardware, no
  3-node or placement comparison; the declared 3 GB database and 2 GB client budgets cannot
  both be enforced on this 2 GiB podman machine, so every level is a level on this machine.
- **Within-run drift is not negligible.** The repeated-design position control read
  0.80/0.90/0.61 on PostgreSQL and 0.86/1.35/1.01 on YugabyteDB. Any PostgreSQL ratio below
  roughly 1.5–2x, and especially sub-millisecond cells, is not resolved by three runs at
  five trials; the 1.52x derived-table cell is in that class.
- **YugabyteDB plan signatures are approximate.** The normalizer was written for PostgreSQL
  text plans; YugabyteDB emits a similar but not identical format, so its
  `same plan structure` column is coarser evidence than on PostgreSQL. The `rank_index_used`
  observation, by contrast, was checked by hand and is solid.
- **The reuse token is synthetic.** `md5(repeat(id::text,32))` is intentionally expensive and
  side-effect-free; absolute reuse costs do not transfer to a real consumer, and the
  1%/100% consumer levels only bound the crossing, not its exact location.
- **One selectivity axis, one fact size, narrow rows, no concurrent writers, no blended
  phase.** The widest threshold returns 179,911 rows rather than 180,000 because 89 facts
  have a NULL parent key; that is labelled and the oracle agrees, but it means the top level
  is 89.96%, not exactly 90%.
- **Five trials under-sample tails.** The reports carry planning/execution times, buffers and
  plan text, but no p99; the reuse and 90% cells are the most likely to have heavy tails.

## What I would measure next

1. The handoff's 3-node YB RF=3 subset **with balanced endpoints** and actual tablet/leader
   placement evidence, plus the colocated/non-colocated pair — the two dimensions this block
   does not touch. That needs a spread-aware client endpoint and either a larger machine or
   the documented resource exception.
2. A consumer-selection sweep (1, 5, 10, 25, 50, 100%) at fixed 16 references, to locate the
   PostgreSQL crossing rather than bracket it, and the same on YugabyteDB to confirm no
   crossing exists.
3. A real expensive consumer and a row-width/fact-count sweep, to test whether the
   PostgreSQL P3 penalty at 0.01% is a fixed join cost or scales with fact width.
4. A matched blended read/write phase and a plan-level `enable_` control set, to separate the
   spelling from the join algorithm the planner picks.

## Disagreements with other analyses

None read: no prior analysis cites these digests. The earlier fixed-order SELECT analyses
are not about this family and are not pooled here.
