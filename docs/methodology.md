# Methodology

The rules every study in this repository follows. They exist so that two results in this
repo can be compared to each other, and so that a reader can tell what a number does and
does not mean.

## 1. Everything runs in containers

No study may require the reader to install a database, a language toolchain, or a CLI
client on their own machine. Podman is the only prerequisite. Benchmark clients run in
containers too, not on the host.

This is not tidiness for its own sake: a benchmark run against whatever PostgreSQL
happened to be on the author's laptop is not reproducible, and a reader who has to
install five things to check a claim will not check it.

**Orchestration uses the plain `podman` CLI rather than `podman compose`.** Neither
`podman-compose` nor `docker-compose` ships with Podman, so requiring one would violate
the rule above for the sake of slightly shorter scripts.

## 2. Results are only comparable within an environment

Hardware changes benchmark numbers by more than most design decisions do. Every result
file names its environment id, and every environment has a specification page under
[`docs/environments/`](environments/) that records the CPU, memory, storage, container
runtime, and — critically — the **known measurement hazards** of that machine.

A result from one environment must never be compared against a result from another.

## 3. Image versions are pinned

Floating tags (`:latest`, `:17`) make a result unreproducible and make two results in
this repository silently incomparable. Every image is pinned to an exact tag in
[`infra/versions.env`](../infra/versions.env), and changing one is a documented event.

## 4. Calculate resources and distinguish per-node from equal-total comparisons

A single-node engine and one node of a three-node cluster get exactly the same CPU and
memory budget. This makes "single vs cluster" answer the question an operator actually
asks: *what do I get when I add two more machines?*

The equal-total framing asks a different question: "is one larger machine or several
smaller ones better at the same resource budget?" Under the owner's 2026-09-14
requirements in [AGENTS.md](../AGENTS.md#required-comparisons-for-future-studies), future
multi-node comparisons include this control as a separately labelled condition with
its own recorded budgets and results. Never pool it with the per-node baseline.

Before measurement, write the resource calculation in the study protocol:

- Use the live resources available to the container VM or each physical host. Reserve
  explicit CPU and memory for the client, OS and supporting processes. Check each
  host's allocation sum against its available resources; any deliberate oversubscription
  is a separate stress condition.
- With a database budget of `C` CPUs and `M` memory across `N` equal workers/nodes,
  allocate `C/N` and `M/N` per node. The matched single node receives `C` and `M`.
  Account for coordinator/master processes inside those totals and record effective
  engine settings and minimum supported allocations. For example, three nodes at
  2 CPU/3 GiB match one node at 6 CPU/9 GiB; this is an example, not a universal size.
- Distinguish database workers/nodes from benchmark client workers. Specify reader
  and writer counts, the offered arrival schedule, queue bound and pool calculation.
  For one connection per active worker, the pool allowance must cover readers plus
  writers plus explicitly counted control connections. Document any other connection
  model. Hold these choices fixed across a controlled pair; concurrency sweeps are
  labelled experimental variables, not implicit multipliers of node count.
- Calibrate the client in reported preliminary checks, then freeze its settings before
  the comparison. Record delivered demand, queueing/rejections and client throttling
  so a load-generator limit is visible. Use balanced connections across supported query
  endpoints for cluster capacity; retain a single-endpoint condition when it answers a
  separate question. Record per-node CPU throttling, memory and connection distribution.

These calculations establish a reproducible budget, not an optimum inferred from CPU
count alone. Validate them with measured saturation and independent trials. A node-count
increase does not itself justify claiming a replication or placement cost.

Limits are CFS quotas (`--cpus`), never core pins (`--cpuset-cpus`). See the environment
notes for why: on a CPU that mixes performance and efficiency cores, pinning silently
biases a container onto slow cores with no signal in the output.

## 5. Correctness gates timing

**A design that answers the question wrongly, quickly, is worth nothing.**

Before any timing is recorded, every design must reproduce the same answers, verified
against values computed independently — in application code, from the generated dataset —
rather than against another query. A failing check aborts the cell and reports no
timings.

Denormalised keys, trigger-maintained rollups and embedded documents are all
opportunities to drift out of sync. A benchmark that never checks would happily report a
broken design as the fastest one, and that failure mode is not hypothetical: it is the
single most likely way for a study like this to produce a confident wrong answer.

Ties are handled explicitly. With a million timestamps truncated to the millisecond, two
rows can share a maximum, and SQL is free to return either; identity checks accept any
member of the tied set rather than demanding one particular row.

## 5a. A correctness check is proven by a design that fails it

An audit that has never caught anything has not been shown to work — it may be checking
the wrong thing, at the wrong moment, under too little contention. Study 01's gate proved
itself by accident, catching harness bugs. From study 02 on, it is proven on purpose.

When a study's question is an invariant (no overbooking, no lost increment, no stale
cache), it includes at least one **negative control**: a deliberately wrong design that is
expected to violate the invariant under the study's workload. The generated report states
for every topology whether each control fired. If a control did not fire, the absence of
violations in the other designs of that experiment is reported as **not evidence** of
their correctness, only of insufficient contention — and the workload is made harsher,
not the conclusion softer.

Negative controls are measured like every other design, and their speed is never shown
as a plain number: a design that breaks the invariant quickly has not been fast.

Invariants have two sides. "Never oversell" is easy to satisfy by never selling; a check
for **over**-booking is paired with one for **under**-booking (the design said "sold out"
while it had seats), and with a reconciliation of what the client was told against what
the database holds.

## 6. The dataset is deterministic and shared

Every design and every engine in a run loads the *same logical dataset*, generated from a
fixed seed. A difference in a result is then a difference in the design, never a
difference in the data.

Generated data must be shaped like the real thing in the ways that matter:

- **skewed, not uniform** — uniform parents hide both skew effects and hot-row contention
- **long-tailed but bounded** — an unbounded power law makes document-oriented designs
  into strawmen rather than fair comparisons
- **ordered the way it would really arrive** — an append-only stream means the surrogate
  key and the timestamp are correlated; shuffling them makes every time-ordered index
  look worse than it is

## 6a. Design the required controlled comparisons

The minimum coverage is defined in
[AGENTS.md](../AGENTS.md#required-comparisons-for-future-studies). Each study protocol
maps those requirements to named design pairs, planned regimes and a correctness
contract. Include a reason for non-applicability or an explicit pending/unsupported gap;
never use absence of a scenario as evidence that the mechanism has no effect.

**Information placement:** start with answers derived from base data. Rollup stores
aggregates at a parent; rolldown copies a parent's key or other information to children;
embedding stores child information in its parent, possibly as a bounded recent slice.
Compare each applicable choice with the baseline. Include insert, correction, deletion
and parent-information changes where supported by the domain; check all maintained
information covered by the claimed contract. Measure reads, maintenance cost, bytes and
relevant size/skew/history regimes. Keep SQL and indexes matched when isolating a copied
field or aggregate; use intermediate designs when an optimized package changes both.

**Physical data colocation:** in every multi-node database study, compare colocated
and non-colocated placement with identical logical data, business operations, replication
factor and resource budget. Document the engine's actual placement mechanism, keys,
partitions/tablets and placement evidence. Plans and RPC/network counters help connect
the measured difference to placement. Hosting all nodes on one machine only describes
the environment; it does not satisfy the placement comparison. If the engine cannot
express a requested layout, record the unsupported case and keep the coverage gap visible.

**Concurrency:** include optimistic conditional/version checks with bounded retries
and pessimistic locking before the competing read/decision/write. Keep the business
invariant, offered work, dataset and worker counts comparable; record effective isolation
and any semantic difference. Exercise low contention and a controlled hotspot, reconcile
acknowledged outcomes, and report retries, aborts/rejections, throughput and latency.
Use an invalid negative control for the invariant. Additional strategies can extend this
minimum when the domain supports them; a database label alone does not identify the
application's concurrency strategy.

Only expand combinations that answer an attributable question; an exhaustive product of
all dimensions is not required. Generated reports retain measurements and coverage gaps;
the concise final analysis links detailed mechanism comparisons in signed discussions.

## 7. Measurement discipline

- **Warmup is discarded.** Connection setup, plan caching and cold buffers are measured
  separately or not at all, never folded into the reported number.
- **Keys are drawn fresh per execution.** Repeating one key measures the buffer cache,
  not the access path.
- **Queries are benchmarked in isolation *and* in a blended workload — both are required.**
  Isolation is the primary instrument, because one blended number hides *which* query
  moved and attribution is the entire point. But isolation alone is not sufficient, and
  the reason is specific rather than theoretical: several costs exist **only** when reads
  and writes overlap in time — MVCC bloat accumulating while readers scan, contention on
  rows writers keep updating, autovacuum and RocksDB compaction stealing CPU, buffer-cache
  competition. Every one of those penalises designs that buy read speed with redundancy,
  which is to say the designs a study like this is usually about. **Measuring only in
  isolation therefore biases the results in favour of exactly the designs under scrutiny.**
  The blended run (experiment D) is where that bias is checked, and its headline is not
  throughput but the *fraction of read throughput retained* once writers appear.
- **Percentiles are computed from every sample.** Means hide the cases a design handles
  badly, and tail latency is usually where a design decision shows up first.
- **A percentile is only reported when the samples support it.** p50, p90, p95, p99 and
  max are always recorded; p99.9 requires 10 000 samples and p99.99 requires 100 000, so
  that at least ten observations sit beyond the quantile. Cell sample counts in a single
  run span roughly 360 to 400 000 — at the low end a "p99.9" would be the single worst
  request relabelled, and it would look far more authoritative than it is. Below the
  threshold the field is omitted and `max` is the honest tail figure.
- **Deep tails from this harness are optimistic and must be read as a floor.** The
  benchmark is closed-loop: each worker waits for its own response before issuing the next
  request. When the server stalls, the load generator stalls with it, so the stall is
  recorded once instead of being charged to every request that would have arrived during
  it. This is *coordinated omission*. It barely touches p50, biases p99 mildly, and biases
  p99.9 and beyond substantially. Tails are therefore compared **between designs** — all
  of which pay the same bias — and are never presented as SLO figures. An open-loop
  generator issuing requests on a fixed schedule would be needed for that, and would be a
  separate harness.
- **Result rows are drained.** A query whose rows are never fetched has not necessarily
  been executed.
- **Destructive operations are measured by count, not duration.** A delete can only
  happen once per row. Each trial gets a fixed operation budget that fits warmup plus all
  trials inside half the pool; running out is a stop signal, never an error to loop past.
- **Each write operation gets a freshly loaded database.** Running insert, update and
  delete in sequence lets later phases inherit earlier phases' bloat — by a different
  amount per design. Measured here as a 30x distortion of one design's delete throughput.
- **A timeout bounds retries, not just new work.** An operation that retries internally
  (lost races, serialization failures) must stop retrying at the deadline too, or one
  pathological transaction keeps a "timed-out" measurement running indefinitely. Work
  already in flight is allowed to finish, so a deadline never manufactures a commit of
  unknown outcome.
- **Quota throttling is recorded, not assumed.** Containers run under CFS quotas, and a
  throttled client puts tens of milliseconds into tails that belong to no design. The
  client's `cpu.stat` is captured per phase and the databases' per cell.

## 8. A single trial is not a measurement

These runs happen on a laptop. One trial can be wrecked by an autovacuum pass, a
checkpoint, or the OS scheduler moving a container onto an efficiency core — and the
result looks exactly like a design difference.

Where a comparison matters, each cell is measured **several times** and the reported
throughput is the **median**, not the mean: the median ignores one bad trial where a mean
absorbs it. Every result also carries its per-trial numbers and a **spread** —
`(max − min) / median` — which is the honest error bar.

Two rules follow:

- A cell whose own trials disagree by more than about 20% cannot support an argument
  about a small difference between designs. The harness flags those at the point of
  measurement instead of leaving a reader to find it in the JSON.
- A difference smaller than the noise floor is not reported as a finding. The generated
  reports suppress differences under ±30% in the controlled-pair tables for exactly this
  reason.

Broad survey runs may use a single trial to cover the matrix in reasonable time; any
number used to *argue a conclusion* should come from a repeated-trial run.

## 8a. Every limitation is tested in the regime that removes it

Almost every design loses somewhere. A conclusion that stops at "design X is slow at Y" is
incomplete in a way that does damage: someone whose domain never does Y reads it and
rejects the design that would have served them best.

So when a design's weakness depends on a property of the data or the workload, the study
also measures the **regime in which that property is absent**, and asks one question of
it: *does the regime change which design wins?* If it does, the conclusion must say so as
an explicit, evidenced exception.

Rules for regimes:

- **Change the shape, hold the volume.** A bounded-history dataset keeps the same number
  of donations by adding donors; a wide-top dataset keeps the same people and donations by
  adding charities. Otherwise every design gets faster because the table got smaller, and
  the regime's own effect is unreadable.
- **The general case stays primary.** Regimes are alternatives reported beside the main
  results, never substitutes for them.
- **Only where it makes sense.** A regime is added for a design only when its limitation
  actually depends on that property. Capping history says nothing about foreign keys, so
  D8 is not re-run under it.
- **Keep the regime map complete**, including limitations whose acceptable regime is
  already covered by another experiment, or which no benchmark can establish.

## 9. Plans are captured, not just times

Wall-clock time alone confuses "did less work" with "had a warmer cache". Every run
captures the plan for every read query:

- **PostgreSQL** — `EXPLAIN (ANALYZE, BUFFERS, VERBOSE)`. Buffer counts attribute cost to
  actual block reads.
- **YugabyteDB** — `EXPLAIN (ANALYZE, DIST)`. On a distributed store the **RPC count is
  the portable cost signal**; milliseconds are an artefact of where the nodes happen to
  be. A benchmark on one machine has no real network, so RPC counts travel and
  wall-clock deltas do not.

Plans are written as readable text, not buried in escaped JSON strings. They are the most
informative artefact a run produces.

## 10. Storage is a result

A rollup column or a denormalised key buys read speed with bytes and with load time. A
report that omits the bytes is telling half the story, so every run records relation and
index sizes and the duration of each load phase.

## 11. Measurements and conclusions are separate artefacts

A generated report contains **numbers and no interpretation**. What the numbers *mean* is
written separately, in a signed analysis under `reports/analyses/`.

This is not bureaucracy. It exists because this project expects **several analysts,
including different AI models**, to read the same data. Different models bring different
priors and notice different things, and where two analyses disagree, **that disagreement
is itself a finding** — it is left visible rather than resolved by editing one of them.

Every analysis carries frontmatter identifying:

| Field | Why |
|---|---|
| `analyst`, `analyst_kind`, `analyst_version` | which model or person, at what version, said this |
| `analyzed_at` | when — models change, and a conclusion has a date |
| `run_id` | which run |
| **`inputs_digest`** | **which exact data** |
| `supersedes`, `status` | whether this replaces an earlier analysis |
| `headline` | one sentence, indexed into the generated report |

The **inputs digest** is the part that does real work. It is a content hash over every
result file in the run, so it is stable across machines and independent of filesystem
order. A run id alone cannot identify what was analysed, because a run can be extended
with extra cells after the fact — which is exactly what happened to study 01's first
matrix. The digest answers "has this precise data already been interpreted, and by whom?"

Rules for analysts:

- **Check the index before writing.** If an analysis already exists for this digest, read
  it. Write a new one to disagree, extend, or bring a different perspective — never to
  restate.
- **Never edit another analyst's file.** Write your own and reference theirs by
  `analysis_id`.
- **Every analysis names where it thinks the measurement is weak.** An analysis with no
  stated doubts has not been done carefully.
- If the digest shown in the report no longer matches the one an analysis quotes, the
  report marks that analysis stale. Its conclusions may still hold, but they were not
  drawn from what is there now.

## 11b. Keep the final analysis short; publish the reasoning beside it

The final signed analysis is the decision document: TL;DR, scope, recommendations,
the measurements needed to support them, weaknesses, and the next experiments. Aim for
roughly 1,000 words; link to detail instead of expanding every mechanism in place.

Detailed design comparisons, SQL/plan walkthroughs, and exchanges between AI models or
human analysts belong in **signed companion documents under `reports/discussions/`**.
Use `docs/templates/DISCUSSION.md`. The final analysis links to its companions near each
affected conclusion; each companion links back and records the same run/digest/commit
provenance. Additional runs must be listed explicitly. Companions may be as detailed as
the evidence requires. This split encourages explanation and disagreement, not less of it.

A discussion distinguishes observations, possible mechanisms, alternative explanations,
and the experiment that could settle a disagreement. Cite other analysts by `analysis_id`;
do not imply they participated in a new exchange. Each author writes a separate signed
document. Never rewrite another author's analysis, or silently shorten a published one.
Old analyses remain evidence for their original inputs; an explicitly superseded edition
is preserved under `outdated/` with its former location documented.

## 12. Reports expire

A report describes one run of one code state on one machine. When a study is re-run and
the conclusions change, the superseded report moves to `reports/outdated/` — it is never
deleted and never silently edited. A conclusion must always be traceable to the run that
produced it.

## 11a. Every final report opens with a TL;DR

Reports are long, and the reader who most needs the conclusion is the one least likely to
read to the end. Both kinds of report therefore start with a TL;DR:

- A **generated report's** TL;DR lists **measured facts selected by fixed rules** — which
  cells failed, whether negative controls fired, which designs broke an invariant, the
  highest and lowest valid result per condition. It ranks by the number only and says
  nothing about which design is better, so the no-interpretation rule still holds.
- A **signed analysis's** TL;DR is the interpretation: a handful of bullets a reader can act
  on — the recommendation, when it changes, and the biggest doubt — each backed by a
  measurement cited in the body.

## 12a. Every result names the code that produced it

Studies get revisited — a new design, a new question, a new analyst — sometimes long after
the run. A number is only reproducible from the exact harness and SQL that produced it, so
from study 02 on:

- Every result file, manifest and generated report records the **repository commit**,
  `git describe`, and whether the study's code (`platform/`, `infra/`, the study directory,
  excluding results and reports) had **uncommitted changes**. A dirty run is flagged as not
  reproducible from its commit alone.
- A run that will be analysed is started from a clean tree with `--tag`, which creates
  `run/<study>/<run-id>` on that commit. `git checkout` the tag reproduces the run;
  `git diff <tag>` shows what has changed in the study since.
- Milestones of a study are tagged `study-NN/<label>`. Study 01's published results predate
  this rule and are marked by `study-01/v1`.
- Signed analyses record `repo_commit` beside `inputs_digest`: the digest identifies the
  data, the commit identifies the code.
- **Enhancing a study is bracketed by tags**: the state before the enhancement and the state
  after it are both tagged, so a change in conclusions between two reports is explained by
  `git diff` rather than by memory. Tags are annotated and never moved or deleted.
- AI sessions commit and tag their work as they go and **never push**; pushing is the owner's.

## 13. Failures are reported

If a cell fails, the matrix continues and the failure is recorded in the run manifest. A
three-hour matrix that aborts on the fifth of twenty-one cells should leave four usable
results behind, not an empty directory. Reports state which cells are missing and why.
