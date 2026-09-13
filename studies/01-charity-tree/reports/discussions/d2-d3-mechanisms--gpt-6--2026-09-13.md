---
discussion_id: d2-d3-mechanisms--gpt-6--2026-09-13
analysis_id: 20260913T125342Z-v3--gpt-6--2026-09-13
run_id: 20260913T125342Z-v3
environment: host-zenbook-ux5406sa
analyst: gpt-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop; exact runtime build not exposed; this session implemented and ran the enhancements"
analyzed_at: 2026-09-13
inputs_digest: 0bc5900a6d7a5d3a
inputs: 20260913T125342Z-v3@0bc5900a6d7a5d3a
repo_commit: c084c73b0cade873406c57fa65ce61f5a26c505a
responds_to: 20260912-study01--gpt-6--2026-09-12, 20260912-study01--claude-opus-5--2026-09-12
---

# Discussion — what makes D3 faster than D2?

[Concise analysis](../analyses/20260913T125342Z-v3--gpt-6--2026-09-13.md) ·
[Measurements and every trial](../20260913T125342Z-v3.md) ·
[Execution order and provenance](../../results/20260913T125342Z-v3/manifest.yaml)

## Question and controlled comparison

D3 copies each donation's charity key, adds indexes, and rewrites several queries.
Those are separable decisions. These trials ask which changes remove actual database
work, and whether the earlier one-trial improvement survives repeated fresh loads.

Every cell regenerated the same 108,081 donations, 5,000 donors and 10 charities.
PostgreSQL 17.11 had 2 CPUs and 3 GiB; the containerized client had 2 CPUs and 2 GiB.
Read trials used explicit `VACUUM (ANALYZE)`, eight connections, one second of warmup
and three measured seconds per query. D2/D3 alternated order across five trials. The
six intermediate variants reversed their order on even trials. This was not a fully
randomized experiment: D2/D3 occupied the earlier block and queries kept fixed order.

The ladder below follows implementation order rather than numeric design ID. Each
step after D11 changes only the named SQL or index decision. **D11 includes the extra
foreign key as well as the copied column**; its storage/write delta cannot identify
the column alone. The sequence is a set of conditional comparisons, not a factorial
estimate of independent effects that can be added together.

| Design | Change from preceding row | Median 12-query score | Spread |
|---|---|---:|---:|
| D2 | Indexed normalized baseline | 5,071.75 | 15.2% |
| D11 | Copy charity key plus FK; retain D2 reads/indexes | 5,171.13 | 7.6% |
| D12 | Add charity/date index | 4,925.65 | 12.1% |
| D13 | Rewrite charity recency queries q02/q05/q12 | 6,796.49 | 10.3% |
| D17 | Rewrite charity sum q08 | 7,287.79 | 8.1% |
| D14 | Add plain charity index | 7,380.71 | 8.3% |
| D15 | Add amount as INCLUDE payload to that index | 7,860.27 | 13.2% |
| D3 | Rewrite ranking queries q03/q04 | 8,897.79 | 9.1% |

Each trial's score is `exp(mean(log(query ops/s)))` over exactly twelve isolated
queries; the displayed number is the median of five such scores. It is not aggregate
application requests/s. Spread is `(max − min) / median`, not a confidence interval.
D3/D2 paired score ratios were **1.529, 1.746, 1.830, 1.754 and 1.774**. All favored
D3, with a median of 1.754. That strengthens the direction of the earlier finding,
without establishing a universal multiplier.

## The charity feed: the SQL must use the index

Merely adding the key and date index produced no aggregate improvement: D11 and D12
remain close to D2. Changing the recency SQL at identical D12 indexes raised the median
q12 feed rate from **1,080.57 to 13,715.95/s**. The query can now select the charity's
latest donations directly before looking up the donor names.

The smallest charity is particularly revealing. It contains 1,756 donations versus
31,293 in the largest. In the [D2 trial-1 after plan](../../results/20260913T125342Z-v3/pg-single-standard/plans/0001-vacuum-small-d2_normalized_indexed-t1-after.txt),
the global date index walks **2,393 donation entries** and performs **1,033 donor
lookups** to return 50 matches for the smallest charity; the plan reports 3,131 shared
buffer hits. For the largest charity, a match is easier to find: 151 donation entries,
135 donor lookups and 410 hits.

With the rewritten SQL, the [D15 trial-1 after plan](../../results/20260913T125342Z-v3/pg-single-standard/plans/0017-vacuum-small-d15_sum_covering-t1-after.txt)
reads **50 donation entries** for either charity. It uses 20 donor lookups and 86 hits
for the smallest, 44 lookups and 139 hits for the largest. The donor lookup remains;
the expensive search for donations belonging to that charity largely disappears.
These are execution counters, not a prediction based on the index name.

D2 versus D3 median q12 rates were **535.46 versus 16,066.15/s** for the smallest
charity (30.0×), but **5,246.45 versus 13,496.10/s** for the largest (2.57×).
Uniformly sampled charity IDs therefore conceal substantial sensitivity to charity
size. A production workload weighted toward the largest charity would have a smaller
feed benefit than the uniformly sampled average suggests.

## Charity sums: separate a useful rewrite from avoiding heap fetches

At identical indexes, D13→D17 changes q08 to sum donations by their copied charity key.
The median uniformly sampled rate rises from **343.08 to 684.94/s**. D17→D14's extra
plain index gives a smaller **684.94→759.39/s** change, within the substantial trial
variation; this experiment does not establish that the extra plain index is necessary.

D14→D15 isolates `INCLUDE (amount_cents)`, retaining the same sum SQL. Under explicit
vacuum preparation, q08 rises from **759.39 to 2,495.97/s** (3.29×). Targeting the largest
charity gives **389.65→810.05/s**; targeting the smallest gives **2,276.96→13,385.58/s**.
The effect is larger than the respective within-design spreads.

The [plain-index plan](../../results/20260913T125342Z-v3/pg-single-standard/plans/0015-vacuum-small-d14_sum_plain-t1-after.txt)
accesses heap data: 1,122 buffer hits for the largest charity, 865 for the smallest.
The covering plan visits 31,293 or 1,756 index entries with **zero heap fetches**, using
123 or 11 hits. Across all five D15 trials, both keys and both bracketing snapshots
record zero heap fetches; D3 does too. That establishes actual heap avoidance for these
prepared reads, although one snapshot is not a counter for every timed execution.

The ANALYZE-only diagnostic explains why the earlier report could observe an
`Index Only Scan` with many heap fetches. D15 trial 1 starts with **31,293 heap fetches**
for the largest charity and ends with **218** ([before](../../results/20260913T125342Z-v3/pg-single-standard/plans/0020-analyze-sum-d15_sum_covering-t1-before.txt),
[after](../../results/20260913T125342Z-v3/pg-single-standard/plans/0020-analyze-sum-d15_sum_covering-t1-after.txt)).
Its uniformly sampled q08 rate varies by 107.3% across trials. Background maintenance
was allowed to act: this is an evolving preparation diagnostic, not a controlled
fixed level of dirty pages. PostgreSQL's visibility map determines whether an index
entry still requires a heap visibility check; including the requested columns alone
does not guarantee zero fetches. [PostgreSQL 17 documentation](https://www.postgresql.org/docs/17/indexes-index-only-scans.html)

## Ranking, unchanged queries, storage and a real mix

D15→D3 retains the schema/indexes but rewrites ranking SQL. The median charity top-donor
rate rises **262.59→620.20/s**, and leaderboard **254.22→468.98/s**. Conversely, an
unchanged donor-history query remains similar across D2/D3: **35,417→36,500/s**.
The global sum is slower in D3 (**340.06→248.27/s**). A good overall score is therefore
not evidence that every query improved. Inspect plans and workload weights before
adopting the complete package.

Median initial relation sizes in the vacuum read cells were **18.13 MiB for D2** and
**25.70 MiB for D3**. Fresh-load insert medians were **7,751.57 versus 6,319.38/s**,
but D3's spread was **49.3%** over three-second windows. The larger footprint is clear;
the write-cost percentage needs longer trials before it is used for capacity planning.

The separate blended condition used four fixed readers and eight scheduled-insert
workers at **500 offered inserts/s for 15 seconds**. The connection pool limit was
20 (readers + writers + loader allowance + four), permitting all twelve workload
workers to hold connections. D2/D3 median read rates were **2,320.47/3,916.20 per second**, a 1.69×
ratio. Both completed approximately 500 inserts/s with no rejection or operation
error, and database counts reconciled all acknowledged inserts including warmup.
Successful scheduled-response p99 ranged **53.95–57.90 ms** for D2 and **52.43–55.26 ms**
for D3. These modest-load samples support this mix; they do not certify a latency SLO
or locate saturation.

## Exchange with the earlier analyses

This is a response to their published work, not a new conversation with either author.

- [`20260912-study01--gpt-6--2026-09-12`](../analyses/20260912-study01--gpt-6--2026-09-12.md)
  requested the repeated trials and decomposition. The new data supports its claim that
  the copied key, indexes and SQL must be discussed together. Its warning about the old
  plan's heap fetches remains correct for that input; the new zero-fetch observation
  answers a different, explicitly vacuum-prepared condition.
- [`20260912-study01--claude-opus-5--2026-09-12`](../analyses/20260912-study01--claude-opus-5--2026-09-12.md)
  described the feed gain as charity copied onto donations. That is directionally useful
  but incomplete as implementation advice: D11/D12 show that copying/indexing without
  rewriting the query does not reproduce the gain. Its index-based recommendation for
  ordinary recency reads is reinforced. Neither author has reviewed these new trials.

## Where the measurement is weak and what would change the conclusion

Five fresh loads reduce one source of uncertainty; they do not remove shared-laptop
noise, CPU quotas, fixed query order, short windows or the small memory-resident dataset.
The author implemented the variants and ran the experiment, so this is not a blinded
independent replication. Initial checks verify the loaded answers and targeted feeds;
this comparison does not establish a constraint protecting copied keys during donor
reassignment. The JSON digest also does not hash the separate container logs; targeted
plans are embedded in JSON as well as linked above.

Repeat longer paired D14/D15 read-and-update windows while recording visibility and
plans through the window; a covering gain that disappears with realistic churn would
narrow the recommendation to read-mostly data. Repeat D2/D3 using measured production
query/key frequencies and longer insert windows. If those weights remove the feed and
sum demand, the copied-key package may cease to justify its storage/write cost.
