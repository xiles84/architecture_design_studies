---
discussion_id: embedding-growth--gpt-6--2026-09-14
analysis_id: 20260914-study01-v3--gpt-6--2026-09-14
run_id: 20260913T172624Z-v3
environment: host-zenbook-ux5406sa
analyst: gpt-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop; exact runtime build not exposed; implementer and experiment operator"
analyzed_at: 2026-09-14
inputs_digest: bd95d304cb39e2de
inputs: 20260913T172624Z-v3@bd95d304cb39e2de, 20260914T104721Z-v3@9ead901cb74c4183
repo_commit: 9b7aa48b329ca44e341fe2aedb6d39de8d4cb70d
responds_to: 20260912-study01--gpt-6--2026-09-12, 20260912-study01--claude-opus-5--2026-09-12
status: current
---

# Discussion — growing donor histories and memory pressure

[Final analysis](../analyses/20260914-study01-v3--gpt-6--2026-09-14.md) ·
[Growth measurements](../20260913T172624Z-v3.md) ·
[Matched memory measurements](../20260914T104721Z-v3.md)

## A growing database, not just another initial row count

D3 stores donations in an indexed child table; D6 embeds each donor's donation history
as JSONB. Each fresh load receives three deterministic cycles of **1,000 inserts,
1,000 amount corrections and 500 deletions**, for 1,500 additional surviving donations.
The Go reference dataset changes only after a successful SQL mutation. Fourteen checks
run after each phase; reads run afterward, without overlapping mutations. All 18 cells
and 162 post-mutation gates in the original growth matrix passed.

Each phase measures global sums, charity sums and recent donor history with eight read
workers. The nominal one-second read window can take longer while in-flight requests
drain. These are short probes of each database state, not sustained capacity or tail
measurements. The mutation streams are sequential and their rates are not comparable
to the separate eight-writer isolated-write benchmark.

The standard medium dataset contains **1,020,455 donations across 45,000 donors**.
The longer-history condition contains **918,808 donations across 5,000 donors**, with
an average history of 183.76 and a maximum of 3,952. Both are finite generated
histories without an imposed per-donor cap; neither represents indefinite growth.

## Longer histories preserve D3's advantage on these standard-memory reads

The table uses **the last deletion phase of cycle three**, then takes medians across
the three independently loaded trials. It does not treat nine phases on one evolving
database as nine independent repetitions.

| Condition at 3 GiB | D3 global sums/s | D6 global sums/s | D3 recent donor reads/s | D6 recent donor reads/s |
|---|---:|---:|---:|---:|
| Medium, ordinary histories | 22.50 | 5.18 | 33,098.02 | 22,399.48 |
| Longer histories, fewer donors | 27.43 | 5.06 | 27,497.76 | 6,737.78 |

The longer-history condition illustrates both costs of embedding: cross-donor
aggregation still processes JSON values, and even a recent-20 query has to open the
donor's history. The SQL expands the array and orders its elements before limiting
the answer. In contrast, D3's donor/date index can stop after the requested entries.
These are different generated regimes, so their difference alone does not quantify
a history-length effect independently of donor population and total row count.

Mutation rates reinforce the tradeoff. In the first cycle of the longer-history
condition, D3's median correction/deletion rates were **1,889.44/1,950.59 per second**;
D6's were **91.06/71.96**. D6 varies greatly across subsequent phases and trials, so
those particular ratios should not be generalized. Each cycle also targets different
donors; a faster later phase does not prove the database improved as it aged.

Median initial-to-final relation sizes in that condition changed **208.06→208.52 MiB
for D3** and **138.95→162.28 MiB for D6**. D6 remained smaller, but its allocated
footprint expanded much more than its roughly 0.16% logical row-count growth. Rewritten
arrays, allocation and maintenance timing are plausible contributors; relation sizes
alone do not separate live payload from reusable/dead space or prove persistent bloat.

## Why the memory control matters

The original constrained condition uses **2,009,226 donations, 45,000 donors and a
256 MiB database limit**, with 64 MiB shared buffers and explicitly smaller memory
settings. Its standard-medium counterpart has different data, so comparing those two
conditions would confound resource configuration and dataset shape.

At the last phase of that first constrained experiment, D6's median recent-donor rate
was **14,125.69/s**, versus **3,728.70/s for D3**. Global sums still favored D3
(**13.89 versus 2.51/s**). The observed reversal prompted a matched run that repeats
the larger dataset, seed, operations, client and CPU budgets under both memory
configurations in the same session.

The matched run **20260914T104721Z-v3**, digest `9ead901cb74c4183`, was produced by
commit `d1b18bf437daf972860261e53b422dc87905cc09`. All twelve cells have identical
initial dataset summaries. Their twelve initial and 108 post-mutation gates passed,
with **1,680 individual checks and no read errors**. The final-phase comparison is:

| Memory configuration | Design | Recent donor reads/s | Donor-rate spread | Global sums/s | Initial→final MiB |
|---|---|---:|---:|---:|---:|
| 256 MiB | D3 | 4,578.18 | 19.6% | 13.91 | 458.02→458.73 |
| 256 MiB | D6 | 16,034.40 | 10.6% | 2.51 | 243.13→261.62 |
| 3 GiB | D3 | 25,222.57 | 8.6% | 13.80 | 458.38→458.96 |
| 3 GiB | D6 | 15,415.07 | 8.8% | 2.47 | 230.45→245.12 |

These are medians of three fresh-load trials, with spread `(max−min)/median`. D6's
donor-read median is **3.50× D3's under pressure**, while D3's is **1.64× D6's at
standard memory**. Every trial has the same direction within its configuration.
D3 retains its global-sum advantage under both. This supports a memory-dependent
donor-read exception on these inputs; it does not make embedding the general winner.

Design order alternates AB/BA by trial, but all constrained cells precede all standard
cells. The numbered files and `ordered-cells.tsv` record execution order; the growth
JSON's `settings.order=1` is only a within-pair placeholder. Regime order is therefore
not balanced against host drift. D6's allocated relation size also differs between
memory configurations despite identical logical data, leaving maintenance/allocation
behavior as part of the observed configuration effect.

## Interpreting footprint and plans cautiously

In the first constrained run, D3's median relation footprint was **457.96→458.73 MiB**;
D6's was **242.82→261.57 MiB**. D6 crossed the 256 MiB container limit during the
mutation sequence while D3 exceeded it from the start. Relation bytes exceeding the
limit establish a larger stored footprint, not that every page belongs to the active
working set or every miss reaches the SSD.

The [first D3 cell's cgroup after snapshot](../../results/20260913T172624Z-v3/pg-single-constrained/logs/0001-memory-growth-d3_flattened_fk-t1-pg-single-after.txt)
compared with its [before snapshot](../../results/20260913T172624Z-v3/pg-single-constrained/logs/0001-memory-growth-d3_flattened_fk-t1-pg-single-before.txt)
records approximately 9 GB of guest block-device reads across the cell. This includes
loading, verification and plans as well as timed work. Guest/host caching means it
cannot be equated to physical SSD traffic.

The saved full-catalogue q09 plan uses a donor with a long history: D3 visits 20 index
entries, while D6 expands 1,000 JSON elements. Both saved plans are warm, and D3 is
faster for that key. That does not explain the uniform-donor throughput reversal under
pressure by itself. D6's smaller/co-located footprint is a plausible mechanism; memory
controls now establish the reversal with fixed logical data. Representative
per-history-tier I/O measurements are still needed to distinguish footprint, cache
residency and query-key selection as mechanisms.

## Exchange and remaining weaknesses

This extends the growth request in `20260912-study01--gpt-6--2026-09-12`. The earlier
`20260912-study01--claude-opus-5--2026-09-12` recommendation against embedding for
cross-donor aggregation remains supported. Its broad statement that embedding never
won PostgreSQL reads needs to remain scoped to its earlier memory-resident inputs.
Neither author is represented as having reviewed this new data.

The experiment ages the data through a small net addition, not weeks of churn. It has
three independent trials, fixed phase/query order, short read windows, shared laptop
cores and no concurrent mutation/read stream. The memory configurations change shared
buffers and other memory settings as well as the container ceiling. The author built
and ran the harness. Next, measure history-length tiers, page/TOAST work during the
timed donor reads, and longer churn with vacuum behavior captured before turning an
observed niche into a deployment recommendation.
