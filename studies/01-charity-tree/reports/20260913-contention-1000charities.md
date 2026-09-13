# Experiment C — concurrency control on the hot rollup row

| | |
|---|---|
| Run id | `20260913-contention-1000charities` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Topology | PostgreSQL, 1 node |
| Engine | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` |
| Scale | `small` — 5000 people, 108568 donations |
| Operation | `insert` — append one donation and update both parent rollups |

## The setup

Ten charities; every donation in the system updates one of their rows. The largest
charity holds 30% of all donors, so under load its row is a point every writer
must pass through. The contention is structural, not contrived — it is what
consolidating an aggregate onto a parent *means*.

| Row | What maintains the aggregates |
|---|---|
| `D3 (no rollup)` | nothing — the ceiling. What the insert costs if no aggregate exists. |
| `D4 trigger` | an AFTER trigger, inside the same statement |
| `D5 blind` | plain UPDATE; the engine's row lock decides the order |
| `D5 forupdate` | SELECT ... FOR UPDATE first — pessimistic, lock then act |
| `D5 optimistic` | read version, UPDATE ... WHERE version = n, retry on loss |

## Throughput as writers are added (inserts/s)

A strategy that wins at 4 writers and collapses at 32 has not solved anything.

| Writers | D3 no rollup | D4 trigger | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|
| 8 | 6.1k | 5.2k | 2.4k | 2.4k | 
| 32 | 8.8k | 6.2k | 3.1k | 3.0k | 

### median latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|
| 8 | 1.25 | 1.42 | 2.51 | 2.47 | 
| 32 | 1.84 | 2.18 | 4.34 | 4.52 | 

### p99 latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|
| 8 | 2.10 | 2.78 | 24 | 24 | 
| 32 | 50 | 61 | 64 | 66 | 

### p99.9 latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|
| 8 | 5.70 | 9.35 | 33 | 39 | 
| 32 | 54 | 64 | 67 | 72 | 

A dash means fewer than 10 000 samples in that cell — too few for p99.9 to be
anything but the worst request relabelled. Use the `worst observed` row instead.

### worst observed (ms)

| Writers | D3 no rollup | D4 trigger | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|
| 8 | 12 | 33 | 41 | 49 | 
| 32 | 60 | 70 | 74 | 85 | 

> Tails here are closed-loop and therefore optimistic: a stalled server stalls
> the load generator too, so the stall is counted once rather than charged to
> every request it would have delayed. Compare across strategies, which all pay
> the same bias; do not read these as SLO numbers.

## Isolation level at fixed concurrency

Held at 8 writers so the isolation level is the only variable. Stricter isolation
is expected to cost *retries* rather than latency, which is why retries are counted
here rather than hidden inside the throughput number.

| Strategy | Isolation | inserts/s | p99 ms | Retries | Errors |
|---|---|---:|---:|---:|---:|
| forupdate | read committed | 2.4k | 24 | 0 | 0 |
| optimistic | read committed | 2.4k | 24 | 268 | 0 |

## Did the aggregates survive? (the result that matters)

After each contended run, every rollup is recomputed from the donation table and
compared. **A strategy that is faster and silently loses increments has not won
anything** — it has published a total that is wrong in a way nothing in the system
would ever notice. `drift` is what the rollups claim was donated minus what the
donation rows actually add up to.

| Writers | Variant | Isolation | Person rows wrong | Charity rows wrong | Drift (cents) |
|---:|---|---|---:|---:|---:|
| 8 | D4 trigger | read committed | 0 | 0 | 0 |
| 8 | D5 forupdate | read committed | 0 | 0 | 0 |
| 8 | D5 optimistic | read committed | 0 | 0 | 0 |
| 32 | D4 trigger | read committed | 0 | 0 | 0 |
| 32 | D5 forupdate | read committed | 0 | 0 | 0 |
| 32 | D5 optimistic | read committed | 0 | 0 | 0 |

No strategy lost an update in this run. That is a result about *these* strategies
on *this* engine, not a general guarantee: every one of them performs its
increment inside a single UPDATE statement, which the engine executes atomically
under its own row lock. A read-modify-write split across two statements — the
shape most application code reaches for first — has no such protection.

---

Raw results: `20260913-contention-1000charities/`.

## Conclusions and analysis provenance

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so that several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, that disagreement is itself a finding and is
left visible rather than resolved by editing one of them.

| | |
|---|---|
| Run id | `20260913-contention-1000charities` |
| Result files | 8 |
| **Inputs digest** | `a9424bb1e126154f` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `a9424bb1e126154f` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `a9424bb1e126154f` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`a9424bb1e126154f`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
