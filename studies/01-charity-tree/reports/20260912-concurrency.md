# Experiment C — concurrency control on the hot rollup row

| | |
|---|---|
| Run id | `20260912T180830Z-concurrency` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Topology | PostgreSQL, 1 node |
| Engine | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` |
| Scale | `small` — 5000 people, 108081 donations |
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

| Writers | D3 no rollup | D4 trigger | D5 blind | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|---:|
| 1 | 1.9k | 1.8k | 1.2k | 960 | 962 | 
| 4 | 4.2k | 3.4k | 2.2k | 1.8k | 1.5k | 
| 8 | 7.0k | 3.0k | 2.3k | 2.0k | 1.6k | 
| 16 | 9.4k | 1.7k | 1.4k | 1.8k | 1.2k | 
| 32 | 8.9k | 1.1k | 1.0k | 1.7k | 534 | 

### median latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 blind | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|---:|
| 1 | 0.51 | 0.54 | 0.82 | 1.00 | 0.99 | 
| 4 | 0.91 | 0.96 | 1.73 | 1.95 | 2.10 | 
| 8 | 1.11 | 1.71 | 2.48 | 2.63 | 2.79 | 
| 16 | 1.39 | 2.60 | 3.20 | 2.83 | 5.46 | 
| 32 | 1.83 | 2.53 | 2.79 | 3.15 | 18 | 

### p99 latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 blind | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|---:|
| 1 | 0.87 | 0.97 | 1.55 | 1.79 | 1.82 | 
| 4 | 1.38 | 3.03 | 4.05 | 5.04 | 8.13 | 
| 8 | 1.63 | 14 | 15 | 15 | 33 | 
| 16 | 17 | 79 | 88 | 88 | 75 | 
| 32 | 50 | 304 | 317 | 201 | 374 | 

### p99.9 latency (ms)

| Writers | D3 no rollup | D4 trigger | D5 blind | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|---:|
| 1 | 1.39 | 1.65 | — | — | — | 
| 4 | 2.02 | 5.28 | 5.88 | 6.30 | 14 | 
| 8 | 5.10 | 24 | 26 | 22 | 44 | 
| 16 | 23 | 115 | 147 | 138 | — | 
| 32 | 56 | — | — | 329 | — | 

A dash means fewer than 10 000 samples in that cell — too few for p99.9 to be
anything but the worst request relabelled. Use the `worst observed` row instead.

### worst observed (ms)

| Writers | D3 no rollup | D4 trigger | D5 blind | D5 forupdate | D5 optimistic | 
|---:|---:|---:|---:|---:|---:|
| 1 | 5.27 | 11 | 21 | 5.20 | 21 | 
| 4 | 20 | 15 | 17 | 28 | 27 | 
| 8 | 17 | 56 | 41 | 27 | 59 | 
| 16 | 36 | 205 | 216 | 240 | 105 | 
| 32 | 60 | 973 | 1048 | 492 | 897 | 

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
| blind | read committed | 2.3k | 15 | 0 | 0 |
| blind | repeatable read | 1.6k | 31 | 15839 | 0 |
| blind | serializable | 324 | 101 | 33388 | 28 |
| forupdate | read committed | 2.0k | 15 | 0 | 0 |
| forupdate | repeatable read | 1.4k | 30 | 14693 | 0 |
| forupdate | serializable | 308 | 118 | 26175 | 18 |
| optimistic | read committed | 1.6k | 33 | 15245 | 0 |
| optimistic | repeatable read | 1.3k | 37 | 14148 | 0 |
| optimistic | serializable | 297 | 110 | 29942 | 31 |

## Did the aggregates survive? (the result that matters)

After each contended run, every rollup is recomputed from the donation table and
compared. **A strategy that is faster and silently loses increments has not won
anything** — it has published a total that is wrong in a way nothing in the system
would ever notice. `drift` is what the rollups claim was donated minus what the
donation rows actually add up to.

| Writers | Variant | Isolation | Person rows wrong | Charity rows wrong | Drift (cents) |
|---:|---|---|---:|---:|---:|
| 1 | D4 trigger | read committed | 0 | 0 | 0 |
| 1 | D5 blind | read committed | 0 | 0 | 0 |
| 1 | D5 forupdate | read committed | 0 | 0 | 0 |
| 1 | D5 optimistic | read committed | 0 | 0 | 0 |
| 4 | D4 trigger | read committed | 0 | 0 | 0 |
| 4 | D5 blind | read committed | 0 | 0 | 0 |
| 4 | D5 forupdate | read committed | 0 | 0 | 0 |
| 4 | D5 optimistic | read committed | 0 | 0 | 0 |
| 8 | D4 trigger | read committed | 0 | 0 | 0 |
| 8 | D5 blind | read committed | 0 | 0 | 0 |
| 8 | D5 blind | repeatable read | 0 | 0 | 0 |
| 8 | D5 blind | serializable | 0 | 0 | 0 |
| 8 | D5 forupdate | read committed | 0 | 0 | 0 |
| 8 | D5 forupdate | repeatable read | 0 | 0 | 0 |
| 8 | D5 forupdate | serializable | 0 | 0 | 0 |
| 8 | D5 optimistic | read committed | 0 | 0 | 0 |
| 8 | D5 optimistic | repeatable read | 0 | 0 | 0 |
| 8 | D5 optimistic | serializable | 0 | 0 | 0 |
| 16 | D4 trigger | read committed | 0 | 0 | 0 |
| 16 | D5 blind | read committed | 0 | 0 | 0 |
| 16 | D5 forupdate | read committed | 0 | 0 | 0 |
| 16 | D5 optimistic | read committed | 0 | 0 | 0 |
| 32 | D4 trigger | read committed | 0 | 0 | 0 |
| 32 | D5 blind | read committed | 0 | 0 | 0 |
| 32 | D5 forupdate | read committed | 0 | 0 | 0 |
| 32 | D5 optimistic | read committed | 0 | 0 | 0 |

No strategy lost an update in this run. That is a result about *these* strategies
on *this* engine, not a general guarantee: every one of them performs its
increment inside a single UPDATE statement, which the engine executes atomically
under its own row lock. A read-modify-write split across two statements — the
shape most application code reaches for first — has no such protection.

---

Raw results: `20260912T180830Z-concurrency/`.
