---
analysis_id: 20261005T-select-synthesis--claude-opus-5-5--2026-10-05
run_id: "select-order-{access,aggregate,existence}-{postgres,yugabyte}-{1,2,3}-20261004; select-join-v3b-domain-{pg,yb}-{i,b}{1,2,3}-20261005; select-join-v4-topo-{single,equal,tablet,distributed,colocated}-{i,b}{1,2,3}-20261005"
environment: host-zenbook-ux5406sa
analyst: claude-opus-5-5
analyst_kind: ai
analyst_version: "Claude Opus 5.5 (claude-opus-5-5) via Claude Code desktop; HIGH session, analyst role; selected effort xhigh"
analyzed_at: 2026-10-05
inputs_digest: baaaab3622c5170faf371f64a51bf4cf2c9a02aa644fcb8c82c143016211f0c6
repo_commit: "ca2a38fb (access/aggregate/existence), 7243fee0 (join/CTE domain), ba821e7b (placement); per-run commits and digests in the appendix"
supersedes:
status: current
headline: "No SELECT spelling wins everywhere. The six inner-join spellings run the same plan; the answer contract, the engine, the index and the parent or page count decide the rest, and every winner below holds only in its measured regime."
---

# Analysis — Study 01 SELECT synthesis: which template, and where it stops — claude-opus-5-5

Synthesis for `task-20261002T131729Z-select-synthesis`. It scores the twelve
predictions of the [SELECT protocol](../../protocols/select-query-templates.md)
against the four measuring families. It also gives the selection matrix the book
uses. Every figure is a per-run value from the generated
[evidence tables](../20261005T-select-synthesis-evidence.md), made from saved,
tagged runs by `docs/select-query-templates/synthesis-20261005/evidence.go`.
Registered claims are `v6-sel-*` in `book/evidence/v6/claims.json`.

## TL;DR

- **Write the join the readable way; the spelling does not matter.**
  `JOIN … WHERE`, `JOIN … ON`, a filtered derived table, either default CTE and
  a comma join ran the **same plan as `JOIN … WHERE` in all 120 pairs**
  (2 engines × 2 phases × 6 selectivities × 5 spellings, three runs each). What
  costs is **forcing `MATERIALIZED`** on a CTE whose filter sits outside it:
  84–101x on PostgreSQL and 34–49x on YugabyteDB at 0.01% selectivity, shrinking
  to about 1x at 90%.
- **Choose by answer first.** `WHERE` on an optional child changes the answer;
  nullable `NOT IN` returned **zero rows** in all 18 checks (6 runs × 3 fanout levels). Then pick
  `EXISTS`/`NOT EXISTS` for membership/absence. On PostgreSQL,
  `JOIN + DISTINCT` costs 1.4–1.7x at fanout 1–10 and **3.6–4.2x at fanout
  100**. On YugabyteDB the differences stay within 16%.
- **Per-parent work flips with the parent count, and it is engine-specific.**
  On PostgreSQL a correlated aggregate or `LATERAL` probe is cheapest for one
  parent and loses from about 100. On YugabyteDB the correlated aggregate loses
  at every count, reaching **6.7–20 s at 10,000 parents (398–541x)**, and
  `LATERAL` wins only at one parent.
- **Keyset pagination pays only past a depth.** With this spelling and no
  matching composite index, keyset is 6–29% slower than `OFFSET` up to depth 100
  and **4.2–4.7x faster at depth 10,000** on PostgreSQL. `UNION ALL` never beat
  `OR` by a resolved margin.
- **Biggest doubt:** the YugabyteDB access and existence runs used the default
  hash-sharded indexes. **None of their 167 plans used an index.** Their
  YugabyteDB numbers describe a scan-only physical design, not YugabyteDB with
  an ordered index.

## What I was looking at

| family | runs | engines and topology | order | cells |
|---|---|---|---|---|
| access paths (P10, P11) | `select-order-access-*` ×6 | PG / YB single node | seeded per run, interleaved | 206/206 valid |
| aggregate and top-N (P8, P9) | `select-order-aggregate-*` ×6 | PG / YB single node | seeded per run, interleaved | 182/182 valid |
| membership and absence (P6, P7) | `select-order-existence-*` ×6 | PG / YB single node | seeded per fanout level | 128/128 valid |
| join/CTE domain (P1–P5, P12) | `select-join-v3b-domain-*` ×12 | PG btree; YB range index; isolated and blended | seeded per run | 1,092/1,092 valid |
| placement (required comparison) | `select-join-v4-topo-*` ×30 | YB single, equal-total, 3-node single-tablet / distributed / colocated | seeded per run | 900/900 valid |

Every applicable negative control fired in every run:

- overlap 36/36;
- aggregate double count and rank ties 6/6 each;
- nullable `NOT IN` and membership multiplier 18/18 each;
- optional-match `WHERE` and limit-before-join 12/12 each.

Page equality under a fixed snapshot held at all 30 depth checks.

The fixed-order primaries (`select-access`, `select-agg`, `select-exist`) are
not re-scored. GPT-6's [evidence review](../discussions/20261004T-select-evidence--gpt6.md)
showed that their order is confounded. The corrected runs replace them for every
verdict below. The customers/orders teaching fixture carries no timing claim.

## Prediction verdicts

A sign counts only when all runs agree. Position drift is the change in one
repeated cell between the first and last place in a run:

- access 1.5–2.4%;
- existence 13–22%;
- aggregate 29–74%;
- join domain up to 42% (one cell: 0.87 and 1.42).

A stable sign smaller than its family's drift is reported as unresolved in
magnitude.

| P | Verdict | Evidence (per-run ranges) |
|---|---|---|
| P1 six inner-join spellings | **confirmed** | 120/120 spelling pairs ran `join_where`'s plan in every run. 27 had a stable sign, all with identical plans, sizes 0.2–40%, and directions that flip between phases. None is a spelling effect. |
| P2 derived table vs plain join | **confirmed** | Same plan in all 24 threshold × set pairs. |
| P3 default vs `MATERIALIZED` CTE | **confirmed, refined** | Filter outside the CTE: PG 84–101x at 0.01%, 51–64x at 0.1%, 23–27x at 1%. YB 34–49x, 20–29x, 8.1–9.5x. Filter inside the CTE: PG ≤1.6x, YB 0.71–1.25x. Where the predicate sits decides the cost. |
| P4 reuse crossing | **confirmed in 3 of 4 regimes; falsified at PG 1%** | `NOT MATERIALIZED`/`MATERIALIZED` crosses between 1 and 2 references at PG 100% consumer, YB 1% and YB 100%. At n=2 it is 1.43–2.06; at n=16, 3.9–12.1. At the PG 1% consumer `NOT MATERIALIZED` wins at every n: 0.17–0.81. |
| P5 LEFT JOIN `ON` vs `WHERE` | **confirmed (answer contract)** | `x01_optional_where` fired 12/12. |
| P6 membership spellings | **confirmed on PG; not observed on YB** | `JOIN + DISTINCT`/`EXISTS`: PG 1.40–1.69 (fanout 1), 1.38–1.44 (10), 3.64–4.23 (100). YB 1.10–1.16, mixed 0.95–1.01, 1.01. `IN` vs `EXISTS`: no resolved difference (PG 0.86–1.30 mixed). |
| P7 absence and `NOT IN` | **confirmed** | Nullable `NOT IN` returned 0 rows in 18/18. PG `NOT EXISTS` is fastest at fanout 100 (`LEFT JOIN … IS NULL` 2.71–3.27x, NOT-NULL `NOT IN` 2.85–3.16x). YB within 1–10%. |
| P8 aggregate crossover | **confirmed on PG; falsified on YB** | PG correlated/post-join: 0.47–0.90 (1 parent), mixed (10), 1.26–1.72 (100), 2.49–2.70 (10,000). YB correlated loses at every level (1.07–1.23 at 1, 398–541 at 10,000), and pre-aggregation never beats post-join with a stable sign (one run 0.99 at 1,000). |
| P9 per-parent top-N | **confirmed on YB; small side unresolved on PG** | `LATERAL`/`ROW_NUMBER`: YB 0.54–0.70 (1), 1.12–1.25 (10), 257–289 (10,000). PG 0.56–0.98 (1, within drift), mixed (10), 1.19–1.95 (100), 1.98–2.20 (10,000). |
| P10 `OR` vs `UNION ALL` | **falsified in the tested shapes** | PG broad shapes: `UNION ALL` 1.23–1.55x, `UNION` 1.25–1.50x slower. Selective shapes within 0.86–1.14. YB 1.32–1.87x slower, but no index was used. |
| P11 keyset vs `OFFSET` | **confirmed in direction; "flat keyset cost" not tested** | keyset/`OFFSET`: PG 1.06–1.29 (depth 1–100), 0.73–0.82 (1,000), 0.21–0.24 (10,000). YB 1.20–1.36 (1–1,000), 0.69–0.72 (10,000). PG keyset fell from 0.90–1.01 to 0.45–0.49 ms as depth grew. |
| P12 convergence with selectivity | **confirmed** | The six defaults never diverge, so nothing needs to converge. The forced-`MATERIALIZED` joined CTE converges from 84–101x (PG) and 34–49x (YB) to 1.02–1.44x and 0.93–1.03x at 90%. |

## Selection matrix

Read across: required answer → readable default → when measured evidence
favours an alternative → when the default or an alternative fails.

| Required answer | Readable default | Measured alternative wins when | Fails / do not | Evidence |
|---|---|---|---|---|
| Joined rows with a filter (bag) | `JOIN … ON … WHERE` | Never by spelling: all six forms share a plan | Forcing `MATERIALIZED` on a CTE that holds the join with the filter outside it (84–101x PG, 34–49x YB at 0.01%) | `v6-sel-01`, `v6-sel-02` |
| A CTE reused n times | Engine default (inlines one reference, materializes two or more) | Forcing `NOT MATERIALIZED` with a selective consumer on PG (0.21–0.81 at n=2–16) | `NOT MATERIALIZED` with n≥2 otherwise (1.43–12.1x slower) | `v6-sel-03`, `v6-sel-04` |
| Parents with optional child data | `LEFT JOIN … ON (child predicate)` | — | Child predicate in `WHERE` drops unmatched parents (control fired 12/12) | `v6-sel-01` controls |
| Parents that have a match (set) | `EXISTS` (or `IN`; no resolved difference) | — | `JOIN + DISTINCT`: PG 1.4–1.7x, and 3.6–4.2x at fanout 100; `DISTINCT` on a name merges identities | `v6-sel-05`, `v6-sel-06` |
| Parents with no match | `NOT EXISTS` | — | Nullable `NOT IN` returns nothing (18/18). PG at fanout 100: `LEFT JOIN … IS NULL` 2.7–3.3x, NOT-NULL `NOT IN` 2.8–3.2x | `v6-sel-07`, `v6-sel-08` |
| One aggregate row per parent | Join then `GROUP BY` | PG, one parent: correlated scalar 0.47–0.90x | Correlated scalar from ~100 parents on PG (1.08–2.70x), and at every count on YB (up to 6.7–20 s at 10,000); summing across two child joins double-counts (control 6/6) | `v6-sel-09`, `v6-sel-10` |
| Top 1 per parent | `ROW_NUMBER()` with a total-order tiebreak | `DISTINCT ON` on PG at 10,000 parents (0.61–0.64x); `LATERAL` probe for one parent on YB (0.54–0.70x) | `LATERAL` probe at 10+ parents on YB (1.12x → 257–289x) and 100+ on PG; `ROW_NUMBER` without a tiebreak (control 6/6) | `v6-sel-11`, `v6-sel-12` |
| Rows matching either of two predicates | `OR` | Never by a resolved margin in these shapes | `UNION ALL` over overlapping arms duplicates rows (control 36/36); an expression predicate without a matching expression index (PG 3.0–3.6x slower) | `v6-sel-13`, `v6-sel-14` |
| The next page of an ordered list | `OFFSET` for shallow pages | Keyset past depth 1,000 on PG (0.73–0.82; 0.21–0.24 at 10,000), past 10,000 on scan-only YB | Keyset at depth ≤100 (6–29% slower on PG); keyset and `OFFSET` answer different page contracts once rows change | `v6-sel-15` |
| Placement of joined data on a 3-node cluster | Engine default (one tablet) | Pre-split three tablets: 0.54–0.73x at 2,000–180,000 rows | Splitting costs the point read (1.22x, isolated); colocation is 1.21–1.54x slower at 180,000 rows | `v6-sel-16`, `v6-sel-17` |

## Conclusions

**Answer first, spelling last.** No inner-join spelling, CTE wording or comma
join moved a plan. A reader who writes the join the way the question reads loses
nothing. The cost only appears when a hint forces the engine to materialize a
joined result before filtering it. The evidence is about where the predicate
lives, not how the SQL is laid out. With the filter inside a forced-materialized
CTE, PostgreSQL paid at most 1.6x; with the same filter outside, it paid
84–101x.

**The engine's default CTE policy was right in 3 of 4 regimes.** Both engines
inline a one-reference CTE and materialize from two references. The plans were
identical to forced `MATERIALIZED` in every n≥2 run. Overriding it helped only
on PostgreSQL with a selective consumer, where `NOT MATERIALIZED` stayed 1.2–4.8x
faster at 2–16 references. That is a narrow, teachable exception, not a rule
to inline by default.

**Membership and absence are semantic choices first.** The
nullable-`NOT IN` empty set is the strongest result in the corpus: 18 of 18
checks, both engines, every fanout. Teach it as three-valued logic: the predicate
is `FALSE` or `UNKNOWN`, never `TRUE` for any outer row once the subquery holds
a NULL (GPT-6's correction to "UNKNOWN for every row"). The performance
difference is a PostgreSQL fact. `JOIN + DISTINCT` and the non-semi-join
anti-joins pay a fanout-dependent tax there and not on YugabyteDB. Fanout
changes with the parent count (10,000 / 1,000 / 100 parents at fanout 1 / 10 /
100), so the growth cannot be attributed to fanout alone.

**Per-parent work is where the engines part most.** The correlated aggregate and
the `LATERAL` probe are textbook "small result" templates. On PostgreSQL they
win only at one parent, and a correlated aggregate already loses at 100 parents.
On this single-node YugabyteDB each probe presumably becomes a separate storage request through
an equality index: the correlated aggregate grew to 6.7–20 s at 10,000 parents,
while PostgreSQL took 16 ms. The book should state the parent count with any
probe template. The per-request mechanism is a hypothesis: no RPC counters were
parsed.

**Pagination is a contract and a depth.** Keyset and `OFFSET` returned the same
page at every depth under a fixed snapshot. Keyset is not free: this
NULL-ordered spelling has a three-way `OR` predicate and no `(amount, id)`
composite index, so it is a bitmap scan with a top-N sort. It loses at shallow
pages and wins from depth 1,000 on PostgreSQL. The textbook "flat keyset cost"
assumes a matching ordered index, which was not measured.

**`OR` is the default for disjunctions.** The predicted `UNION ALL` advantage
needs an `OR` that blocks an index. PostgreSQL served each indexable disjunction
with one BitmapOr, so the advantage never appeared. The physical factor that did
matter is an expression index for an expression predicate: 3.0–3.6x on
PostgreSQL. That is attributable to the index, not to the SQL.

**Placement is a separate decision from spelling.** On the 3-node cluster,
pre-splitting the relations into three tablets with leaders on three nodes made
2,000–180,000-row reads 0.54–0.73x the single-tablet cluster, and 180,000-row reads 0.71–0.76x the
single node and the equal-total single node. It cost the 20-row read
1.22x (isolated). Colocation won small reads and lost the broad scan. See the
[review-4 qualification](../discussions/20261005T-select-join-review-4--claude-opus-5-5.md).

## Where I think the measurement is weak

- **YugabyteDB physical design.** The teaching DDL's single-column indexes are
  hash-sharded on YugabyteDB. Index nodes appeared in 0/103 YB access plans and
  0/64 YB existence plans, against 51/103 and 12/64 on PostgreSQL. P10 and P11
  on YugabyteDB therefore compare scan plans; their crossing depths and `OR`
  margins must not be taught as YugabyteDB behaviour with an ordered index.
  Only the join domain has a range-index arm (199–201/273 YB plans used an
  index). Registered as `v6-sel-gap-01`.
- **Single node for three families.** Access, aggregate and existence have no
  3-node, equal-total or placement arm; only the join family does.
- **Latency of one sequential reader.** Five samples per cell, no throughput or
  tail claims. In the placement runs the 1 CPU client was throttled in 32–46%
  of its CFS periods.
- **Position drift.** Aggregate drift is 29% (PG) and 74% (YB). PG
  pre-aggregation beats post-join only at 10,000 parents (0.77–0.80), inside
  that drift, so it is an observation, not a registered claim. The domain cells
  drift by up to 42%, larger than every stable spelling sign.
- **Bracketed, not located.** Every crossing is between two adjacent levels: PG
  keyset in 100–1,000, PG correlated in 1–100, YB `LATERAL` in 1–10, reuse in
  1–2.
- **One laptop.** The cluster nodes share eight cores in one VM, so nothing here
  is a network or independent-hardware result.

## What I would measure next

1. A YugabyteDB range-index arm (`(amount ASC)`, `(customer_id, amount DESC)`)
   for the access and existence families, gated on index use, to replace the
   scan-only YB numbers.
2. Keyset with a matching composite index `(amount DESC, id)`, at depths 1–10,000
   plus 200/500, on both engines, to test the flat-cost component of P11.
3. Aggregate family with ≥15 interleaved rounds and ≥5 whole runs, to resolve
   PG pre-aggregation at 1,000–10,000 parents and `DISTINCT ON` at 10–100.
4. Existence with parent count held fixed while fanout varies, to separate the
   two.
5. The distributed schema on the equal-total single node with
   `SPLIT INTO 3 TABLETS`, to separate tablet parallelism from node spread.

## Disagreements with other analyses

- The original access analysis said keyset wins "an order of magnitude" at depth
  10,000 and "OR is never worse". I agree with GPT-6: the PG advantage is
  4.2–4.7x. `OR` is not resolved-worse in the tested shapes, which is not the
  same as never worse.
- The original aggregate analysis said "`DISTINCT ON` is the fastest top-N
  spelling at scale on both engines". The corrected runs support it only on
  PostgreSQL at 1,000–10,000 parents. YB at 10,000 includes a 0.997 run, which
  is practically equal, as GPT-6's qualification says.
- The original aggregate analysis placed YugabyteDB's correlated crossover
  between 1 and 10 parents. The corrected runs show the correlated aggregate
  already slower at one parent (1.07–1.23x), so YugabyteDB has no crossover.
- The v4 join analysis is read with my
  [review-4 qualification](../discussions/20261005T-select-join-review-4--claude-opus-5-5.md),
  which corrects its TL;DR baseline, penalty wording and blended-phase claim.

## Coverage of the owner's required comparisons

- **Sizing:**
  - Single-node families: DB 2 CPU/3 GiB, client 2 CPU/2 GiB, one sequential
    reader, pool of 2.
  - Placement: per-node and equal-total arms, reader and writer endpoints
    balanced over three nodes, per-node and client throttling recorded.
  - The effective guest is 8 CPUs and ≈15.36 GiB.
- **Rollup, rolldown, embedding:** not applicable. These are syntax pairs over
  one physical model; Study 01's D3–D11 remain the information-placement
  evidence.
- **Colocated vs non-colocated:** measured for the join workload with verified
  tablets and leaders (`v6-sel-16`, `v6-sel-17`). Not measured for the other
  three families (gap).
- **Optimistic vs pessimistic:** not applicable. No arbitration was varied; the
  blended writer is a fixed 10/s maintenance load.

## Appendix — inputs (run, set, sha256 of the result JSON)

Each row is one anchored run; the evidence tables print the same digests, its
producing tag and commit. Aggregate inputs digest:
`baaaab3622c5170faf371f64a51bf4cf2c9a02aa644fcb8c82c143016211f0c6`.

| run | set | sha256(json) |
|---|---|---|
| select-order-access-postgres-1-20261004 | access-pg | `0783debd960ffb0e76f9cf0771a11a3c065e433a6e1efa11519f06447204a0f4` |
| select-order-access-postgres-2-20261004 | access-pg | `357620cf7fb2687698f082b5d1a7294fe72f41fb8c30d93def9d9156a43c12f5` |
| select-order-access-postgres-3-20261004 | access-pg | `f0b0f2c1bf41c51187b6a1fe78d73d735bd6d0a5fcd98ac849bfc37fc0c56aa0` |
| select-order-access-yugabyte-1-20261004 | access-yb | `9a3b8b777b99388a1cb983808a53bb0f976aad693cb14e7dd4f4fe8da8b99d7f` |
| select-order-access-yugabyte-2-20261004 | access-yb | `818eaca45753f59f49626bb1135e68dec19cf9359c611f6048632cbb478e4b9d` |
| select-order-access-yugabyte-3-20261004 | access-yb | `ee27b23999ddd2405173c950ed2128986d7f723807e30232c02190a5532c2e52` |
| select-order-aggregate-postgres-1-20261004 | agg-pg | `76268c647e965cd579b85330b32ccdf716c959230c56fe0d71216ce23ab39dc2` |
| select-order-aggregate-postgres-2-20261004 | agg-pg | `b742e387f0cd26e3a730227db33212c88c0535251e947928e0e6c480cbce0b35` |
| select-order-aggregate-postgres-3-20261004 | agg-pg | `5285f21e2330e7118cb9fd7d3fc06576ad2fcc8d6b5c0619825f0b3726cd4823` |
| select-order-aggregate-yugabyte-1-20261004 | agg-yb | `21ad333d1e4ff8a9da998c1a748ec21beca168a49d122da6c11658dcb9ced477` |
| select-order-aggregate-yugabyte-2-20261004 | agg-yb | `e21df42d7c564cc456eba1e11a5d0ec43eaed59d7849ccdfd21b1be390eb2ab7` |
| select-order-aggregate-yugabyte-3-20261004 | agg-yb | `971d1d0df2e057b6aaae034b0834b4a706559e22842bdec96ecb822cddd03014` |
| select-order-existence-postgres-1-20261004 | exist-pg | `a1cd2e20e7fd50961ed5562a099cd46e5b993e69847f8822daf4e3696d74f042` |
| select-order-existence-postgres-2-20261004 | exist-pg | `5a12edf5d2b7a1df53338b019e7e77cec89e64d80302731659073759283f4f39` |
| select-order-existence-postgres-3-20261004 | exist-pg | `b2cce0c8c6ad02cc39a45e93ddfe29dcd26f23988af6ebc3700a63965286eafd` |
| select-order-existence-yugabyte-1-20261004 | exist-yb | `b8e15f58e4e81a72c956ddd1e23bb656a932cf3a33ce8b41f59111be6fe9a82c` |
| select-order-existence-yugabyte-2-20261004 | exist-yb | `703b92c3644821ed9ddb95fb5a5e4566ea8d342cc63a998e02f5be3a11dc6749` |
| select-order-existence-yugabyte-3-20261004 | exist-yb | `62979fd70cec76ac2cfc6f2441d73310d61cd0108485bb3abae839bbf451f262` |
| select-join-v3b-domain-pg-i1-20261005 | dom-pg-i | `724642ee2a9420bed8795540dfb5132831794dbb5f4c0588ace27cf0707614da` |
| select-join-v3b-domain-pg-i2-20261005 | dom-pg-i | `093d13091d8585b37c203fd328653867313ca3d810fe35372267895672eaae8f` |
| select-join-v3b-domain-pg-i3-20261005 | dom-pg-i | `6a2111f1912b1cd00d196e68d818df6cf08c060976d579f884c1b6e033110d3d` |
| select-join-v3b-domain-pg-b1-20261005 | dom-pg-b | `4ad32c4a1ab1d6ca5968ee0a7bb75f9580e86aa1e6ec5df69ce60d8ec0dd3f8c` |
| select-join-v3b-domain-pg-b2-20261005 | dom-pg-b | `eee9aedad20795de40569abd06dda87c0e0280d29c5f1e953bc46100bd5a2306` |
| select-join-v3b-domain-pg-b3-20261005 | dom-pg-b | `db176a1e371b3c5b22dd520f8087ebf675cf7129560391c68294592c11370cb6` |
| select-join-v3b-domain-yb-i1-20261005 | dom-yb-i | `aed6b10e8425f3871462de7a93c38b198680edb2c1ef539e7a94d85c90c8ffa6` |
| select-join-v3b-domain-yb-i2-20261005 | dom-yb-i | `6da46b4f1e2f9b5782d3d345e55a8b8f87df08804be165c472414a92cdefedb0` |
| select-join-v3b-domain-yb-i3-20261005 | dom-yb-i | `68a755a4e33271cef1b47a2a6e22925bcf393e554a509635321141c29deda022` |
| select-join-v3b-domain-yb-b1-20261005 | dom-yb-b | `28deac3dcf98eeb5da0432bc8fb437001b425cc7b793e3e2c6f78a35de377235` |
| select-join-v3b-domain-yb-b2-20261005 | dom-yb-b | `5d23235fbf9206d0c6da526df936bb9cd18736e0d6d4515beb92705a6d59389a` |
| select-join-v3b-domain-yb-b3-20261005 | dom-yb-b | `824b5579aeac034ee9c3afe73c771ad5f097531d215069a6a29649c67c1e30b7` |
| select-join-v4-topo-single-i1-20261005 | topo-single-i | `5ff0161edb2e32cb4c63da0123451f452b614e6c9ef4330c1c4724292cd9c92a` |
| select-join-v4-topo-single-i2-20261005 | topo-single-i | `a63c107a532494fd6f8f9bbe8f49ca6e095cc14c0bead22fc848ac0e6a61693e` |
| select-join-v4-topo-single-i3-20261005 | topo-single-i | `a99d6402480ba00cece3985fbe62f122e1519382e6e02ed4cd47ed5b9218d42d` |
| select-join-v4-topo-single-b1-20261005 | topo-single-b | `8321c3ddbaeb64c88aa2315bd66bd90c9cacdf42177e4b4540c1749b92205d27` |
| select-join-v4-topo-single-b2-20261005 | topo-single-b | `82d73027d7f16910d7cbf4e49cc60d32f8da47ef39a3a5c1104607ef6a11db63` |
| select-join-v4-topo-single-b3-20261005 | topo-single-b | `339cf3f56a1f3e809ce65c22ec1e1227e135ef71f45fc7267205cd18f1a94d81` |
| select-join-v4-topo-equal-i1-20261005 | topo-equal-i | `64dcb81ee841e6d5e8e980097756dfe625516647b960dc1693ac2570685697d7` |
| select-join-v4-topo-equal-i2-20261005 | topo-equal-i | `bffe080fefcb13b42d5b1f89a1906d70060534d141e743c0fe13234b1449c26e` |
| select-join-v4-topo-equal-i3-20261005 | topo-equal-i | `eb1eda2ea2f9f4df1f54c640b91c49ca45a903827d648ce2474d1535b5e99432` |
| select-join-v4-topo-equal-b1-20261005 | topo-equal-b | `5b2b8c1a474427784366133e9a46b5d6b023eb5d2fdc72210a01c7cbea6181b0` |
| select-join-v4-topo-equal-b2-20261005 | topo-equal-b | `81698f69d454c8a835180e6bdc7c6f7a0c1c3ed6785717fbb80932bc01c8eacd` |
| select-join-v4-topo-equal-b3-20261005 | topo-equal-b | `f6b9c08cf1e49b8c28b767a485a5bb4fbeee295d6c2605ab8743972334a3e1dd` |
| select-join-v4-topo-tablet-i1-20261005 | topo-tablet-i | `6122526384aa46083db9172dc68a117e37915dd37a9807571f85973be1ae4a3a` |
| select-join-v4-topo-tablet-i2-20261005 | topo-tablet-i | `020584f5fd06c11a3b471d4e25265c7ac45ecf5d310ed241e5e9ff4a7d1395b1` |
| select-join-v4-topo-tablet-i3-20261005 | topo-tablet-i | `34447d84de019b9cadd9dc7faf9ad95497d2c6212ce3eb4bfca4f79b5d799b9d` |
| select-join-v4-topo-tablet-b1-20261005 | topo-tablet-b | `cc9d1cbf384650d419ca7fb3f16bcc9905806cda53b4f640ebe25024346aaaa2` |
| select-join-v4-topo-tablet-b2-20261005 | topo-tablet-b | `9bbb7a64c3dc758819f23782c3d16987a59c490cd657008ceaf5faf1bb6fc16a` |
| select-join-v4-topo-tablet-b3-20261005 | topo-tablet-b | `00fd931d2a48a57a27f502788a1ea816a3a69d7365f1016c6e129d4a8e0d9a6b` |
| select-join-v4-topo-distributed-i1-20261005 | topo-distributed-i | `1e89ee032810798eedb2bc9d723a7bd8386f9717146a338ac038608541db9805` |
| select-join-v4-topo-distributed-i2-20261005 | topo-distributed-i | `c379b5634cc0c83d6c5dda73cbabb8240e851f6bdedc2c45499d4b6853647554` |
| select-join-v4-topo-distributed-i3-20261005 | topo-distributed-i | `261a8263c1f684e8bb2ab17c74a9216b2f4c722b6c9cbc7dad47db99e13674dc` |
| select-join-v4-topo-distributed-b1-20261005 | topo-distributed-b | `ca8e03004e7caf4a78d0441c47b3804975c22e3924cc1156f775731c60a23b24` |
| select-join-v4-topo-distributed-b2-20261005 | topo-distributed-b | `233446b75e88c47c953204b767c1ef54ff47c91967edbe07448dceab512448a7` |
| select-join-v4-topo-distributed-b3-20261005 | topo-distributed-b | `1e62a6d00eb0c539c464a0418f5fa923fbde0556ece7132df02b25d848c5a782` |
| select-join-v4-topo-colocated-i1-20261005 | topo-colocated-i | `771f30fddb4e5880544603c9cc85e1424602ea23efe11e298b93cb35f0f11a38` |
| select-join-v4-topo-colocated-i2-20261005 | topo-colocated-i | `88e0feb9497ff35fc6cadf0e9abcbdc8d5d90c7655d3c6b6b9caf84aa32aad26` |
| select-join-v4-topo-colocated-i3-20261005 | topo-colocated-i | `efa57cff77782a8c5bad2db94a32b3437bb84a5ca5d51709bff14055937684fa` |
| select-join-v4-topo-colocated-b1-20261005 | topo-colocated-b | `93eb5f6087c05d62860f236404ceb449cc9526da746005adc3a396506d030522` |
| select-join-v4-topo-colocated-b2-20261005 | topo-colocated-b | `a3dc08acf943eb0f29692a94c144eee7b13abc3364ca7ecfd00f31ec0ad332c0` |
| select-join-v4-topo-colocated-b3-20261005 | topo-colocated-b | `b4e20666adf6f00ecb4771f4d9c949f8d857571a0fd4e190953dbe1b4fafaa89` |
