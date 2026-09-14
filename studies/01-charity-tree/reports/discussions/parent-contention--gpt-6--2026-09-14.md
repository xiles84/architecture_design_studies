---
discussion_id: parent-contention--gpt-6--2026-09-14
analysis_id: 20260914-study01-v3--gpt-6--2026-09-14
run_id: 20260913T172624Z-v3
environment: host-zenbook-ux5406sa
analyst: gpt-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop; exact runtime build not exposed; implementer and experiment operator"
analyzed_at: 2026-09-14
inputs_digest: bd95d304cb39e2de
inputs: 20260913T172624Z-v3@bd95d304cb39e2de
repo_commit: 9b7aa48b329ca44e341fe2aedb6d39de8d4cb70d
responds_to: 20260912-study01--gpt-6--2026-09-12, 20260912-study01--claude-opus-5--2026-09-12
---

# Discussion — simpler totals, lock-first writes and the busiest parent

[Final analysis](../analyses/20260914-study01-v3--gpt-6--2026-09-14.md) ·
[All trials and measurements](../20260913T172624Z-v3.md)

## The comparison and its limits

The earlier analysis asked whether D5's lock-first application code can outperform
D4's trigger under contention, and whether maintaining only sums removes the expensive
parts of a full rollup. This experiment adds **D16**, which stores only donation sums
on each person and charity. D4/D5 store counts, sums and temporal extrema; D5 supports
inserts only in this harness. The common application comparison therefore uses only
inserts. Corrections and deletions are separate D4/D16 trials on fresh data.

PostgreSQL 17.11 has 2 CPU/3 GiB. The client has 2 CPU/2 GiB, **four fixed readers and
32 writer workers**, with a pool limit of 44 connections. The read mix and key-selection
rules stay fixed while the offered insert rate changes from 250 to 2,500/s. The hot
condition directs 90% of inserts to ten donors belonging to the largest charity;
the other condition selects across donors without that additional concentration.
There are still only ten charities, so the latter does not eliminate parent contention.

Each condition has three fresh-load trials, one second of warmup and a nominal
20-second arrival window. Requests have predetermined due times, a queue capacity of
256, and visible rejection when the queue fills. Accepted work drains after the window;
reported completed/s includes this drain. All accepted and warmup inserts reconciled,
and all recorded count/sum audits passed. This tests the maintained counts and sums;
the inherited audit does not independently recheck every stored first/last timestamp.

## Lock-first D5 retains its high-contention advantage

At **2,500 offered inserts/s**, median completed rates across three trials were:

| Write concentration | D4 full trigger | D5 lock-first application | D16 sum-only trigger |
|---|---:|---:|---:|
| No extra donor hotspot | 352.41/s | 753.90/s | 483.50/s |
| 90% to ten donors in one charity | 246.64/s | 433.31/s | 241.34/s |

D5/D4 median-rate ratios are **2.14×** and **1.76×**. Under the hotspot, D5's spread
is 3.9%, versus 18.4% for D4 and 19.3% for D16. The no-extra-hotspot D4 result has
84.7% spread, so its exact multiplier is less stable. The ordering supports the earlier
GPT-6 qualification: application maintenance is not always slower than a trigger when
the application chooses a different locking strategy.

These rates are measured completion capacity under overload, not rates at which the
offered workload was fully served. In the hotspot, D4 rejected **44,502–45,341** of
50,000 arrivals per trial; D5 rejected **41,108–41,465**; D16 rejected
**44,547–45,474**. All recorded SQL error counts were zero. A report omitting rejection
would describe this workload misleadingly.

Successful scheduled-response p99 medians were **1.93 seconds for D4**, **0.93 seconds
for D5**, and **2.08 seconds for D16** in the hotspot. These include queue waiting but
exclude rejected requests. They establish neither an all-request latency percentile nor
a production SLO.

## Removing extrema does not remove the shared parent row

D16 still updates the person and charity rows synchronously for every insert. Its
hotspot completion rate is close to D4's, despite maintaining fewer values. That is
consistent with contention remaining at the shared parent row; reducing aggregate
maintenance alone does not remove that serialization point.

D16 also supplies less precomputed read work. At the hotspot's 2,500 offered writes/s,
median read rates were **5,492.75/s for D4**, **14,534.10/s for D5**, and **3,896.25/s
for D16**. D16 still derives ranking and other non-sum answers from donations. This
is an application tradeoff between different designs, not a pure measurement of lock
duration: their read implementations and physical indexes differ.

In particular, D16 retains D3's donation amount covering index; D4 substitutes indexes
on maintained donor totals and recency. A difference between them cannot all be assigned
to extrema arithmetic. The separate mutation tests remove concurrent reader work but
retain these intentional physical-design differences.

## Low offered rates can conceal inconsistent trials

Without the extra hotspot, all designs completed the 250/s offered workload without
rejection. Under the hotspot, D16 also completed all three trials without rejection.
However, **one D4 trial rejected 1,126 of 5,000 arrivals**, and **one D5 trial rejected
686**, even though their other two trials rejected none. Their near-250/s median rates
therefore do not establish that 250/s is reliably safe in this environment.

The median read rates in that low-rate hotspot were **16,490.45/s, 18,065.40/s and
8,466.30/s** for D4, D5 and D16. The simpler design avoided rejections in these samples
while yielding fewer reads; declaring one universal winner would hide that tradeoff.

## Isolated corrections and deletions do not yet give a precise extrema cost

| Operation | D4 median/s | D4 spread | D16 median/s | D16 spread |
|---|---:|---:|---:|---:|
| Correct amount | 3,273.29 | 88.1% | 3,521.78 | 4.6% |
| Delete donation | 3,022.27 | 84.5% | 3,421.70 | 88.3% |

Most trials favor the sum-only design, but the large variation prevents a reliable
percentage attribution to extrema maintenance. For deletion, both designs fell to
roughly 570–584/s in trial three after exceeding 3,000/s in the first two. That shared
slowdown is a reason to investigate the environment and maintenance activity, not to
explain the whole gap from trigger source alone. Finite deletion pools also end some
windows before 15 seconds; the generated report records actual elapsed durations.

## Exchange and next discriminating experiments

This extends `20260912-study01--gpt-6--2026-09-12` and qualifies the broader D5 cost
claim in `20260912-study01--claude-opus-5--2026-09-12`. Neither author participated in
this new review. D5's high-contention insert advantage is replicated, while its missing
correction/deletion maintenance still prevents treating it as a complete replacement.

For a cleaner extrema-cost experiment, keep the indexes identical and restrict readers
to sum questions answered identically by D4/D16. Repeat longer isolated mutations with
checkpoint/vacuum timing recorded. For usable capacity, sweep arrival rates around the
first rejection and repeat the entire condition; retain rejection counts beside tails.
Asynchronous totals remain a separate, unmeasured design with different freshness rules.

The evidence remains three trials on shared laptop cores with CPU quotas, short windows
and an assumed read mix. The author implemented and ran this harness. The companion
reports what the controls establish without promoting an overloaded median to a service
guarantee.
