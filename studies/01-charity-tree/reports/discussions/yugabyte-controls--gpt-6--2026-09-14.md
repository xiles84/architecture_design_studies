---
discussion_id: yugabyte-controls--gpt-6--2026-09-14
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

# Discussion — YugabyteDB exceptions and deployment controls

[Final analysis](../analyses/20260914-study01-v3--gpt-6--2026-09-14.md) ·
[Measurements](../20260913T172624Z-v3.md) ·
[Run manifest](../../results/20260913T172624Z-v3/manifest.yaml)

## What was controlled

All comparisons use YugabyteDB 2025.2.6.0-b111, freshly generated seed-42 data,
independent initial Go correctness gates, three trials and recorded effective READ
COMMITTED isolation. The single-node FK and cache conditions retain the original
2 CPU/3 GiB per-node budget. The equal-total deployment condition compares one
6 CPU/9 GiB node at replication factor 1 with three 2 CPU/3 GiB nodes at replication
factor 3. The client always has 2 CPU/2 GiB on this same eight-core laptop.

Rates below are medians of trial rates; spread is `(max − min) / median`. Ratios of
group medians and within-trial ratios are distinguished. The generated report retains
every trial, including unsuccessful correctness controls. No DocDB physical storage
size was collected; container memory limits are not database storage measurements.

## Foreign keys cost something, but the earlier multipliers did not repeat

| Comparison | D3, with FKs | D8, without FKs | D8/D3 median-rate ratio |
|---|---:|---:|---:|
| Single-node inserts, 30 seconds | 356.74/s | 448.54/s | 1.26× |
| Three-node whole-donor erasure | 122.29/s | 136.02/s | 1.11× |

Single-node insert spreads were **4.9% and 7.4%**. D8/D3 ratios within the three trials
were approximately **1.22, 1.24 and 1.36**. This supports a repeatable operation-specific
cost at this budget, with a smaller magnitude than the earlier isolated 1.8× result.
It does not establish that removing integrity checks is a good application tradeoff.

Erasure used **45,000 initial donors and 1,020,455 donations**. Each measured phase
erased 11,257 donors, taking **79.42–95.19 seconds** rather than the requested 120-second
ceiling. The finite operation budget ends these phases first; concurrency permits the
small overshoot beyond the nominal quarter-population count. These are much longer
windows than the earlier subsecond erasure samples, but still change the dataset as
they run. The paired D8/D3 ratios were **1.08, 1.04 and 1.11**: the earlier roughly
1.51× erasure advantage did not reproduce at this larger scale.

The changed scale and longer windows prevent calling this a correction of the old
measurements. The result narrows their applicability. A transaction trace or
statement-by-statement FK-work comparison would be needed to attribute the remaining
gap to specific distributed checks rather than the whole constraint package.

## D10 passes the exercised cache contract on both topologies

The cache race directs 90% of scheduled inserts to **one donor**, offers 250 inserts/s,
and runs 32 writers with four fixed donor-portal readers. D9 is the intentionally
defective negative control; D10 changes sorting and locking in cache maintenance.

**D9 failed the cache audit in all six trials**, at warmup and after the measured
window, with the hot donor's cached ID sequence disagreeing with the base table.
**D10 passed all six**, as well as all separately loaded correction and deletion
cells. A correctness failure excludes D9 from valid-capacity conclusions even though
the database acknowledged its writes without an operation error.

| Topology | D10 race completed inserts/s | D9 diagnostic inserts/s | D10 correction/s | D9 correction/s | D10 deletion/s | D9 deletion/s |
|---|---:|---:|---:|---:|---:|---:|
| One node | 14.68 | 14.84 | 232.81 | 255.17 | 185.16 | 198.30 |
| Three nodes | 16.62 | 16.94 | 192.44 | 227.07 | 147.18 | 164.27 |

The race throughput difference is small, but the workload is already severely
overloaded: D10 rejects **6,743–6,814 of 7,500 offered requests**, depending on trial
and topology. Successful scheduled-response p99 has group medians of **23.48 seconds**
on one node and **21.82 seconds** on three. All accepted inserts and warmup counts
reconcile. Those facts show what this single hot row can complete under the specified
queue, not an acceptable operating point.

The separate mutation tests suggest that D10's correction/delete maintenance costs
more than D9's on YugabyteDB: median-rate ratios D10/D9 are about **0.91/0.93** on one
node and **0.85/0.90** on three. Those costs must accompany the correctness benefit;
the earlier PostgreSQL statement of nearly free repair cannot simply be transferred.

**Contract boundary:** the inherited audit compares the ordered donation IDs. It does
not validate every cached amount, currency, timestamp or note, and these experiments
do not overlap insert, correction and deletion streams in the same phase. Initial
Go verification and isolated mutation audits do not establish all interleavings or
reassignment behavior. D10 passed the recorded contract; a claim of unrestricted
concurrent cache correctness would exceed these measurements.

## Equal total resources still favored one larger node in this setup

| D3 measure | One node, 6 CPU/9 GiB | Three nodes, total 6 CPU/9 GiB |
|---|---:|---:|
| Median twelve-query geometric score | 1,721.99 | 726.79 |
| Isolated inserts/s, 30-second windows | 1,140.82 | 270.93 |
| Completed inserts/s at 500 offered/s | 499.95 | 220.18 |
| Median read rate during those arrivals | 210.70/s | 38.90/s |
| Median successful scheduled-response p99 | 42.06 ms | 1,393.62 ms |
| Rejected arrivals, out of 15,000 per trial | 0, 0, 0 | 8,107, 8,167, 8,157 |

The single-node score is **2.37×** the cluster's; isolated insert throughput is
**4.21×**. The score summarizes twelve separate read tests rather than a blended
application workload. In the actual arrival workload the cluster rejected over half
the offered requests, although every accepted write succeeded and reconciled.

This comparison changes replication and the distribution of resources together, as a
deployment choice would. It does **not** isolate the intrinsic cost of replication.
All cluster SQL enters through `yb-n1`, so only two CPUs are available at the query
endpoint. The [trial-1 read-cell counters](../../results/20260913T172624Z-v3/yb-cluster3-standard/logs/0075-equal-total-d3_flattened_fk-t1-yb-n1-after.txt)
minus their [before snapshot](../../results/20260913T172624Z-v3/yb-cluster3-standard/logs/0075-equal-total-d3_flattened_fk-t1-yb-n1-before.txt)
show quota throttling in **69.1% of CPU scheduling periods** at n1, versus **0% and
10.9%** at n2/n3. These counters span the entire cell, including load and verification.
The imbalance supports testing balanced client connections before generalizing that
equal-resource sharding cannot improve performance.

RF=3 uses majority replication; it does not wait for all three replicas on every
write. This experiment measures the complete deployment, not an isolated quorum
round trip. [YugabyteDB replication architecture](https://docs.yugabyte.com/stable/architecture/docdb-replication/replication/)

## Local node stop: count preservation with a visible service penalty

In each trial, the runner stopped **yb-n3**, ten seconds into a 30-second window of
100 offered inserts/s. The client remained connected through **yb-n1**. The fault
logs record a stop and restart in all three trials; the regenerated reports retain
accepted, rejected, completed and reconciled counts.

| Trial | Completed of 3,000 offered | Rejected | Elapsed including drain | Successful response p99 | Count after restart |
|---|---:|---:|---:|---:|---:|
| 1 | 1,753 | 1,247 | 40.52 s | 15.016 s | 109,934 |
| 2 | 1,749 | 1,251 | 40.06 s | 15.055 s | 109,930 |
| 3 | 1,749 | 1,251 | 30.03 s | 15.020 s | 109,930 |

Recovered counts match `108,081 initial + 100 warmup acknowledgments + completed`
in every trial ([trial-1 fault](../../results/20260913T172624Z-v3/yb-cluster3-standard/logs/0078-local-node-stop-d3_flattened_fk-t1-fault.txt)
and [recovered count](../../results/20260913T172624Z-v3/yb-cluster3-standard/logs/0078-local-node-stop-d3_flattened_fk-t1-recovered-count.txt)).
There were no terminal read/write errors; the bounded client instead rejected demand
while accepted requests waited. Preserving acknowledged counts does not mean that
availability met a latency target.

This is a local non-query-node stop with a grace period. It is not loss of the client
endpoint, a power failure, a network partition, or proof of durability across physical
hosts. There is also no matching no-fault 100/s condition, so the entire throughput gap
cannot be assigned quantitatively to the stop. An independent-host run and endpoint
failover remain necessary deployment tests.

## Exchange with the earlier analyses and the remaining doubts

This responds to the published analysis IDs
`20260912-study01--gpt-6--2026-09-12` and
`20260912-study01--claude-opus-5--2026-09-12`; neither author participated in this new
review. The first asked for these exception checks. The new data supports its warning
against broad "foreign keys are free" claims, while reducing the old operation-specific
multipliers. It extends D10's recorded correctness evidence to both YB topologies and
adds a measured local-failure outcome to the earlier untested resilience framing.

Three trials, shared laptop cores, CPU quotas, a single query endpoint and a local
container network remain important limits. The author also implemented and operated
the harness. Next, spread client connections across nodes, repeat FK/erasure windows
at several data sizes, check complete cache payloads under mixed mutation streams,
and use independently documented hosts for network and failure testing.
