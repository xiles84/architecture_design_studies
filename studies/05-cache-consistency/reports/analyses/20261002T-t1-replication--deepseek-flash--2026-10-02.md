---
analysis_id: 20261002T-t1-replication--deepseek-flash--2026-10-02
run_id: 20261002T185604Z, 20261002T185911Z, 20261002T190350Z, 20261002T190748Z, 20261002T191221Z, 20261002T191405Z, 20261002T191609Z, 20261002T194046Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH session, Deep Code CLI; T1 replication of Study 05's cache, churn and engine claims across eight whole runs."
analyzed_at: 2026-10-02
inputs_digest: per-run a8684db17dce1b34 (A1 cache, relaxed first), 6d9110fa6e235d18 (A2 cache, strict first), d92f414e80be3b12 (C1 yb-single, relaxed first), a67c106081445252 (C2 yb-single, strict first), d94ceb9845fcf33c (B1 db-only), 62a4db1d826c6a9b (B2 equal-total), a9f216283e7ff7e2 (D1 churn 600, through first), 95c934ce25faf269 (D2 churn 600, aside first)
inputs: 20261002T185604Z@a8684db17dce1b34, 20261002T185911Z@6d9110fa6e235d18, 20261002T190350Z@d92f414e80be3b12, 20261002T190748Z@a67c106081445252, 20261002T191221Z@d94ceb9845fcf33c, 20261002T191405Z@62a4db1d826c6a9b, 20261002T191609Z@a9f216283e7ff7e2, 20261002T194046Z@95c934ce25faf269
repo_commit: 37ce688
supersedes:
status: current
headline: The cache-throughput half of v2-14 reproduces (2.24x/2.54x) but its p99 half does not (about 3x, not an order of magnitude); the strict/relaxed nulls hold on PostgreSQL and yb-single with the primary's yb direction not reproducing; the 600 s churn reproduces its 2.0 TTL crossings and 0 wrong reads, but one cell recorded an impossible cache value the gate did not catch.
---

# Analysis — T1 replication of Study 05's cache, churn and engine claims — deepseek-flash

## TL;DR

- **`v2-14`'s throughput half reproduces; its tail half does not.** The cached warm read was **2.54x** and
  **2.24x** the no-cache baseline (primary band 2.2-3.5x), but p99 fell only from **27.6/26.9 ms to
  10.1/8.5 ms** — about **3x**, not the primary's **order of magnitude** (15.8-27.1 ms → 0.8-2.0 ms).
- **`v2-16`'s null holds.** The strict/relaxed warm read differs by **7.3%** in one order and **5.7%** in
  the other, with the **sign flipping**, and both differences sit inside the cells' **8.7%-22.8%** trial
  spread. "Strict costs nothing measurable in read throughput" is supported again.
- **`v3-03`'s labelling holds; its point figures do not.** Both framings are recorded and kept apart —
  db-only **8,199** ops/s (cache budget `none`) and equal-total **7,931** ops/s (cache 1 CPU / 768 MiB),
  **3.3%** apart inside the **17.7%/18.4%** spread. The primary's 11k/8.0k did not reproduce.
- **`v3-04`'s null holds, its direction does not.** On yb-single the strict/relaxed read is **1.043** and
  **1.105**, inside spreads of **29.9%-74.7%** and **21.7%-51.1%**; both fresh runs lean *strict*-faster,
  the opposite of the primary's relaxed-596-vs-strict-514. Nothing here is beyond spread.
- **`v3-02`'s crossings and wrong reads reproduce; one cell found a correctness hole.** Both strategies
  crossed **2.003-2.004** real 300 s TTL boundaries with **0 wrong reads** in all four cells (and hard
  expiry fired 0 times, again). But `20261002T191609Z`'s **aside-strict** cell recorded **1 impossible
  cache value** in 1,173,729 reads while its gate passed — the study's floor says an impossible value must
  always be zero, so this is a **limit on `v3-02` and a gap in the gate**, not a clean replication.
- **`v3-01`, `v3-05` and `v3-06` are not re-measured** and are recorded as deferred/gap targets (see
  feasibility): the 1800 s horizon was not run, and the two `yb-cluster3`/placement claims cannot run on
  this machine and have no placement evidence in the harness.

## What I was looking at

Eight whole runs on the published budgets, one benchmark at a time, all tagged from clean trees:

| run | topology | designs | phases | trials | what it carries |
|---|---|---|---|---|---|
| `20261002T185604Z` (A1) | pg-single | none / redis-through-relaxed / redis-through-strict | verify,warm,mixed | 3 | `v2-14`, `v2-16` |
| `20261002T185911Z` (A2) | pg-single | same | verify,warm,mixed | 3 | `v2-14`, `v2-16` |
| `20261002T190350Z` (C1) | yb-single | redis-through-relaxed / redis-through-strict | verify,warm,mixed | 3 | `v3-04` |
| `20261002T190748Z` (C2) | yb-single | same | verify,warm,mixed | 3 | `v3-04` |
| `20261002T191221Z` (B1) | pg-single | redis-through-strict, framing `db-only` | verify,warm,mixed | 3 | `v3-03` |
| `20261002T191405Z` (B2) | pg-single | same, framing `equal-total` | verify,warm,mixed | 3 | `v3-03` |
| `20261002T191609Z` (D1) | pg-single | through-strict / aside-strict, churn 10m | verify,calibrate,warm,mixed,churn | 3 | `v3-02` |
| `20261002T194046Z` (D2) | pg-single | same | verify,calibrate,warm,mixed,churn | 3 | `v3-02` |

The executed pair order was read from each manifest's `design_order_pg-single` / `design_order_yb-single`
and is genuinely reversed within each pair: A1 relaxed→strict, A2 strict→relaxed; C1 relaxed→strict, C2
strict→relaxed; D1 through→aside, D2 aside→through (this was the lesson from the `v2-12`/`v2-17`
replication, where eight different seeds all produced the same order).

**Warm read (per cell), PostgreSQL arms:**

| run | none | relaxed-thru | strict-thru | cached/none | strict/relaxed |
|---|---|---|---|---|---|
| A1 (relaxed first) | 2,969 | 7,527 | 8,076 | **2.54x** | 1.073 |
| A2 (strict first) | 3,336 | 7,484 | 7,057 | **2.24x** | 0.943 |

p99 moved 27.6/26.9 ms → 10.1/8.5 ms (cached) — a 2.7-3.2x cut, not 10x. Cache hit rate was >98% of
reads in every cached cell; every gate passed, `wrong_reads` 0 and `impossible_cache_values` 0 in the
warm/mixed phases.

**YugabyteDB 1-node strict/relaxed** (C1/C2): strict 6,511/6,945 ops/s, relaxed 6,244/6,287; ratio
1.043/1.105; spread 74.7%/29.9% and 51.1%/21.7%; gates green.

**Resource framing** (B1/B2): `db-only` warm 8,199 ops/s (cache budget none) and `equal-total` 7,931
ops/s (cache `cpus=1 memory=805306368`); both record `resource_framing`; 3.3% apart.

**Churn 601 s** (D1/D2): both strategies crossed 2.003-2.004 boundaries with 0 wrong reads; hard expiry
0 in all four cells (probabilistic early expiry 10,938-12,364); D1 aside-strict recorded **1 impossible
cache value**.

| prediction | verdict |
|---|---|
| P1 `v2-14` cached read ≥1.5x baseline in both orders | **confirmed on throughput** (2.54x, 2.24x). The claim's p99 clause is **not** reproduced and is reported as such. |
| P2 `v2-16` strict/relaxed inside spread, no consistent direction | **confirmed** (7.3%/5.7%, sign flips, inside 8.7%-22.8%) |
| P3 `v3-03` both framings labelled and separate, ratio not claimed | **confirmed** (both record framing and budget; 3.3% apart inside spread) |
| P4 `v3-04` strict/relaxed inside spread, no consistent direction | **confirmed as a null**; the primary's direction did not reproduce (both fresh runs lean the other way, inside spread) |
| P5 `v3-02` both strategies ≥2 TTL crossings, 0 wrong reads, no winner | **confirmed on crossings and wrong reads** (2.00, 0 in all four cells); a separate correctness breach (1 impossible value) is reported and is not one of P5's named falsifiers |
| P6 gates and provenance | **confirmed**: every cell's gate and audit passed, every throughput has its wrong-read count, all eight runs tagged from clean trees — but the gate did not catch the impossible value, which is a coverage gap, not a pass |
| P7 `v3-01` 1800 s re-run or deferred | **deferred as planned**: the 1800 s horizon was not run; the 600 s pair is the re-test and the primary's 1800 s figure stays the primary's |
| P8 `v3-05`, `v3-06` recorded as coverage gaps, no promotion | **confirmed**: both are recorded as gaps, nothing promoted |

## Conclusions

**`v2-14` is a throughput result whose tail clause needs re-scoping.** Two orders put the cached read at
2.24x and 2.54x the no-cache baseline — inside the primary's 2.2-3.5x — so the claim's headline survives.
The p99 clause does not: the cached tail here is 8.5-10.1 ms, not 0.8-2.0 ms, and the improvement is about
3x rather than an order of magnitude. The likely reason is the client: the same client budget that
Study 02 found client-bound at read-heavy designs. Either way, the claim must not keep the "order of
magnitude" for p99 on this evidence; it should carry the throughput band and record the tail as not
reproduced under this session's client.

**`v2-16` replicates as a null and should be promoted.** The strict/relaxed difference is inside the
cells' own spread and reverses sign between orders, which is exactly the relation the claim asserts. Read
throughput only — the write-path cost stays unmeasured.

**`v3-03` replicates as a labelling claim.** The point of the claim is that resources can be counted two
ways and that the two must be labelled and never pooled. Both runs record their framing and budget and the
difference is inside spread; the primary's 11k/8.0k point figures did not reproduce and should not be
quoted.

**`v3-04` replicates as a null with an unresolved direction.** Both fresh runs put strict slightly ahead,
within spread, so no design is named; but the primary's specific direction (relaxed ahead) did not
reproduce and should be recorded as not replicated.

**`v3-02`'s crossing relation replicates, but the correctness floor did not hold cleanly.** Four cells
crossed 2.00 boundaries with 0 wrong reads — the claim's stated relation — yet one cell published an
impossible cache value once. The study's own definition makes that value a breach of the floor ("must
always be zero"), and the cell's gate passed because it checks read statements and outbox orphans, not the
cache's own invariants during churn. `v3-02` therefore gains a limit naming the impossible value and the
gate gap; the gate gap is a harness finding for a follow-up, not a reason to soften the claim.

**Correctness and provenance.** Every `verify` phase passed; `wrong_reads` are 0 everywhere they are
counted; no run was produced from a dirty tree; eight run tags exist. The one breach is the impossible
value above, which the gate did not classify.

## Where I think the measurement is weak

1. **The client is the suspected p99 limit.** The primary's cached p99 (0.8-2.0 ms) is an order of
   magnitude better than anything measured here at the same client budget; the throughput matches, so it
   is the tail that moved. This was not isolated with an `explain`/plan or a client-footprint cell.
2. **The impossible value is a single observation with no cause.** One in 1,173,729 reads; not repeated in
   the other three cells; no per-phase impossible count outside churn. It is reported, not explained.
3. **`v3-04`'s spreads are large** (up to 74.7%), so only a very large effect could have been seen through
   them; the null is weak evidence, as the claim's own limit says.
4. **Churn volumes are a closed-loop floor**, hard expiry fired 0 times in every cell (only probabilistic
   early expiry fired), and the 5x burst was not run.
5. **One host, one machine, 2 GiB, declared budgets above it**; absolute throughput is session-specific;
   Redis is a shared cache and its container budget is outside the `db-only` comparison.
6. **`v3-01`'s 1800 s horizon was not re-run**, and `v3-05`/`v3-06` could not run at all on this machine.

## What I would measure next

1. **A client footprint for the cache read** (the Study 02 `READ_CONNS` lever) to test whether the p99
   clause is client-bound.
2. **The impossible value**: instrument the cache payload's provenance during churn so an impossible
   value fails the cell and names its phase; then re-run `v3-02` to see whether it recurs.
3. **`v3-01`'s 1800 s horizon** in a dedicated long run.
4. **The placement mechanism** for `v3-06`: a tablet/leader distribution record, without which the
   colocation pair remains a gap.
5. `t5` is the last of the five T1 replication tasks; the open `select-*` program is the next queue tier.

## Disagreements with other analyses

None on the data. This tranche re-scopes `v2-14` (its p99 clause) and `v3-03`/`v3-04` (point figures and
direction), promotes `v2-16` and `v3-02` on their replicated relations, and records the impossible value
as a finding against the gate rather than against the designs.
