# Lessons learned

Things that cost time, produced a wrong number, or would have produced a wrong number if
they had not been caught. Written down so they are not rediscovered.

Findings **about databases** belong in study reports. This file is about **running the
experiments**.

---

## Measurement

### A correctness gate catches harness bugs, not just design bugs

The verification pass was added to catch designs that drift out of sync. The first thing
it actually caught was two bugs in the measurement apparatus:

1. `-cmd verify` did not load data, so it verified against whatever the *previous* run's
   write benchmark had left in the database. Every check failed, with timestamps from the
   benchmark's own inserts as the giveaway.
2. `SUM(bigint)` returns **`numeric`** in PostgreSQL, not `bigint` — because a sum of
   bigints can overflow one. pgx decodes it as `pgtype.Numeric`, and a converter that
   only knew `int64` silently produced **zero**. Every total read as 0.

Both would have been invisible in a timing-only benchmark: the queries ran, returned
rows, and were fast. **A benchmark with no correctness gate does not fail — it lies.**

### Comparing designs requires identical data, from a fixed seed

Obvious in principle, easy to lose in practice. Every cell regenerates the dataset from
the same seed rather than sharing a database, so no cell can be contaminated by a
previous one's writes.

### Write benchmarks destroy the dataset they measured

Insert, update and delete all mutate the data the read benchmark just used. Ordering is
therefore fixed — load → verify → explain → read → write — and each cell reloads from
scratch. Any subsequent command run against a used database is measuring rubble.

### Draining result rows is not optional

A query whose rows are never fetched has not necessarily been executed. The read loop
iterates every row, otherwise a `LIMIT`-less aggregate could be "measured" without the
server ever materialising it.

### Fresh keys per execution, or you are benchmarking the buffer cache

Repeating one key measures cache residency, not the access path.

### One trial produces confident nonsense

The first survey run reported D3 (with foreign keys) deleting at 6.2k/s and D8 (without
them) at 3.0k/s. Removing a constraint cannot make deletes twice as slow; the number was
noise — a checkpoint or an autovacuum pass landing inside one ten-second window.

It is worth dwelling on how convincing that number looked. It was 2x, comfortably past any
"ignore small differences" threshold, and it appeared in a controlled pair built
specifically to isolate one variable. Nothing about its presentation said "noise".

Hence `-trials N`: report the **median** across trials, not the mean, and carry the
**spread** `(max−min)/median` as an error bar. The median discards one bad trial; a mean
absorbs it and launders it into a finding. Survey runs may use one trial to cover the
matrix; any number used to argue a conclusion needs several.

### A deep percentile needs samples behind it, or it is just the maximum in disguise

Adding p99.9 looked free — every sample is already kept, so it is one more index into a
sorted slice. Checking the actual sample counts first was what made it honest: within one
run, cells ranged from **363 to 392 154 samples**.

At 392 000 samples, p99.9 is backed by ~392 observations and is a real measurement. At 363
it is the top 0.36 of a single sample — the maximum, wearing a percentile's clothes. And
the low-count cells are precisely the slow queries a reader would most want a tail figure
for, so the failure mode is not rare, it is targeted at the interesting cases.

Hence the thresholds: p99.9 needs 10 000 samples, p99.99 needs 100 000 — ten observations
beyond the quantile, minimum. Below that the field is omitted rather than estimated, and
`max` is reported instead, which is at least what it claims to be.

### Closed-loop benchmarks understate their own tails

This harness has each worker wait for its response before issuing the next request. That
means a server stall also stalls the load generator, so the stall gets recorded once
instead of being charged to every request that would have arrived during it —
*coordinated omission*.

The bias is negligible at p50, mild at p99, and large at p99.9 and beyond. So the deeper
the tail, the less the absolute number means. Deep tails here are a **floor** on what an
open-loop client would see, and are only compared between designs, which all pay the same
bias. Publishing them as SLO figures would be wrong, and the report says so wherever they
appear.

### Finite pools must be measured by count, not by duration

Deletes consume rows, and each row can only be deleted once. The harness measured them
like everything else — for a fixed duration — and in the repeated-trials pass that went
badly: 3 s of warmup plus 5 × 8 s of trials needed ~258 000 deletes against a 108 081-row
table. The pool ran out partway through, and every subsequent trial measured workers
spinning on instant "no such row" failures: **111 million errors, 0 ops/s**, reported as
the result. Erasure (`delete_person`, a 5 000-row pool) drained during *warmup* alone.

The correctness gate did not catch it, because nothing answered a question wrongly — the
measurement itself was empty. Two fixes:

1. Destructive operations get a **fixed operation count per trial**, sized so warmup plus
   every trial fits in *half* the pool (the last rows of a table are not representative).
   Throughput is that count over the time it actually took.
2. Exhausting the pool is a **stop** signal (`errExhausted`), never an error to count and
   loop past. A worker that hits it exits, and the cell is flagged as suspect.

The invalid runs are kept with an `INVALID.md` beside them rather than deleted.

### Write phase order distorted a headline number by 30x

Recorded below as a suspected confound; now measured. The survey ran insert → update →
delete against one table and reported **D6 deleting 28 donations/s**. The same delete on a
freshly loaded table: **~830/s**. The insert and update phases had grown D6's arrays and
bloated its GIN index first, and the delete phase paid for all of it.

That is not a subtle bias — it inverts how bad D6 looks by more than an order of magnitude,
and it would have gone straight into the conclusion. Every write operation now runs against
its own fresh load, and is audited on the state *it* produced. Survey `update`/`delete`
figures from `20260912-small` carry this confound and should be read alongside the
isolated re-run.

### A trigger-maintained cache can pass every isolated audit and still be wrong

D9's cache trigger passed the verification pass and the survey's post-write audit. It was
only under *concurrent* reads and writes (experiment D) that the audit found one wrong donor
cache — small enough to dismiss as a fluke.

**Correction, found later:** concurrent reads were never required. Once audits ran after
*each* write operation, D9 was caught corrupting 10 of 5 000 caches on a plain 8-writer
insert burst. The survey's audit passed only because it ran once, at the end of the cell —
after the update and delete phases had rebuilt the caches of thousands of donors and
**healed most of the damage**. The audit's position masked the bug it existed to catch. An
audit belongs directly after the operation under suspicion, not after whatever happens to
run last. A dedicated 4 × 30 s reproduction found
**3–5 wrong caches in every window**, rising with writer count. Recording five diagnosed
examples per audit, rather than a count, is what made it fixable: they showed two
unrelated bugs.

- **Ordering race.** Every one was an adjacent swap of consecutive donation ids. The
  trigger *prepended*, so the cache was in commit order; the table is in `donated_at`
  order. Concurrent inserts for one donor commit out of timestamp order.
- **Lost update.** A committed donation missing while later ones were present. The rebuild
  on the update/delete path is one `UPDATE ... SET col = (SELECT ... FROM donation)`. It
  waits on the person row lock; PostgreSQL re-checks the *target row* after the wait but
  does not re-run the subquery against a fresh snapshot, so it writes a cache built before
  the competing insert committed.

Lessons: correctness checks for denormalised data must run under concurrency, not just
after it; audits should record examples, not tallies; and "the database handles
concurrency" is true of a single statement's target row, not of everything that statement
reads. The fix is a separate design (D10) rather than an edit to D9, so both the bug and the
price of fixing it stay reproducible.

### Write operations in one cell are not independent

`insert`, `update` and `delete` run in sequence against the same table. By the time
`delete` runs, the table has been grown by the insert phase and churned by the update
phase — and by *different amounts* in each design, because each design's earlier phases
ran at different speeds. A design with faster inserts hands its delete phase a larger,
more bloated table.

The ordering is fixed and documented, but it means cross-design `delete` comparisons carry
a confound that `insert` comparisons do not. Worth isolating in a dedicated run before
drawing conclusions from it.

### The same value can be stored two ways in one JSONB column

The bulk loader writes timestamps as Go's RFC 3339 (`...T23:55:20.197Z`); PostgreSQL's own
`jsonb_build_object()` in the D9 trigger writes `...T23:55:20.197+00:00`. Both parse to the
same instant and every query casts before use, so nothing broke — but a cache audit that
compared raw JSONB would have called every trigger-touched row a mismatch and reported a
perfectly working design as broken.

The audit compares **donation-id sequences** instead, which tests what actually matters
(right children, right order, none missing) and is immune to representation. The general
lesson: when application code and the database both write into the same column, they will
not agree on formatting, so never build a correctness check on byte equality.

---

## Hardware and containers

### Heterogeneous CPUs make core pinning a trap

The test host is an Intel Core Ultra 7 258V: 4 performance cores plus 4 low-power
efficiency cores. The guest kernel sees 8 undifferentiated CPUs and cannot tell Podman
which is which, so `--cpuset-cpus` can land a container entirely on E-cores and run it
materially slower **with no signal in the output**.

The repository therefore uses `--cpus` CFS quotas so every container gets an equal share
of whatever it lands on, and reports spread across trials so that scheduling luck appears
as variance rather than hiding inside a single number.

### The container VM's memory is not the host's memory

WSL2 defaults to half of host RAM. A 32 GiB laptop presented 15.4 GiB to Podman. The
3-node topology asks for 3 × 3 GiB plus a 2 GiB client and fits only just. Always check
`podman info` rather than trusting the host spec.

### `podman machine` can exist without ever having been started

The default machine on this host had been created two months earlier and never booted, at
4 CPUs / 2 GiB — far too small for the 3-node cluster. Worth verifying before planning a
run around resource budgets.

---

## Tooling friction

### Git Bash rewrites POSIX paths before Podman sees them

On Windows, `podman run -w /src` becomes `-w "C:/Program Files/Git/src"`. Setting
`MSYS_NO_PATHCONV=1` stops that, but then `$(pwd)` produces `/c/extra/...`, which the
native Windows Podman binary does not understand as a bind-mount source — and it does not
error. It **silently mounts an empty directory**, so the run appears to succeed and
produces no files.

Fix, applied in `infra/lib.sh`: disable the conversion, then convert host paths
explicitly with `cygpath -m` in a `hostpath()` helper used for every bind mount and build
context.

### Never edit a bash script that is currently running

Bash does not read a script into memory up front. It parses one compound command at a
time and keeps a byte offset into the file, so editing the file in place while it runs can
leave that offset pointing into the middle of a different line.

A large `for` loop is parsed in full before its first iteration, so a run already inside
one is safe — but everything *after* the loop has not been read yet. Adding lines earlier
in the file shifts those bytes and can garble the tail.

This happened here: `run-study.sh` was edited to append rather than truncate its manifest
while a matrix was mid-run, which put the trailing manifest block at risk. The measured
results were never in danger — each cell is written by its own container — but the run
record was.

Two mitigations, in order of preference:

1. Edit a copy and switch scripts between runs.
2. If an editor must touch the file, prefer tools that write a new file and rename over the
   old one (`sed -i` does this): the running shell keeps its open file descriptor on the
   original inode and finishes reading the version it started with.

### Heredocs are a poor way to write source files

Writing Go and SQL through shell heredocs failed on quoting often enough to be a net
time loss. Source files are written directly; heredocs are reserved for short,
comment-free content.

---

## YugabyteDB specifics

### YSQL does not listen on loopback inside its own container

`yugabyted start --advertise_address=<name>` binds YSQL to that address only. A readiness
probe against `127.0.0.1:5433` inside the container is refused even though the server is
healthy — which looks exactly like a failed startup. Probe the container's own hostname.

### The image ships `ysqlsh`, not `psql`

Readiness helpers must take the client binary as a parameter rather than assuming `psql`.

### Node count does not imply replication factor

A three-node cluster can serve RF=1 and look identical from the outside. On this image
(2025.2.6.0) `yugabyted` raises the universe to RF=3 by itself as the third node joins;
older versions did not. The cluster script asserts `numReplicas == 3` from
`yb-admin get_universe_config` and **refuses to report an RF=1 universe as a 3-node
cluster**.

### Nodes must join one at a time

Starting all three simultaneously races the master election and intermittently produces a
two-node universe with the third orphaned.

### Repeated DDL churn on a YugabyteDB cluster gets sessions killed

Isolating write operations means reloading before each one: drop, create, COPY, index
backfill, repeated five times per cell. On the 3-node cluster, D7's cell survived its reads
and inserts and then died in a reload with `FATAL: terminating connection due to
administrator command` (57P01) while building an index — the server terminated the
session, and the whole cell was discarded.

Two lessons. First, a load is idempotent (it begins with a full drop), so a
*server-initiated* termination should be retried, narrowly — schema errors, verification
failures and timeouts are real results and still fail the cell. Second, the runner tore the
cluster down before anyone could read its logs, so the cause could not be established.
`run-study.sh` now captures each database container's log tail into the results directory
*before* teardown whenever a cell fails.

### `pg_total_relation_size()` is meaningless on YugabyteDB

Data lives in DocDB, not PostgreSQL heap files, so the size functions return zero. The
harness records that fact explicitly rather than recording a zero that reads like a
measurement.

### Tablet counts must be pinned, not auto-detected

YugabyteDB derives shards-per-table from the CPU count it observes, which varies with how
the container happened to be scheduled. A benchmark whose shard count moves between runs
is not a benchmark. `ysql_num_shards_per_tserver` is set explicitly.

### `EXPLAIN (ANALYZE, DIST)` is the measurement that travels

On a single machine there is no real network, so wall-clock deltas between 1-node and
3-node understate distribution costs badly. RPC counts do not depend on where the nodes
are, so they are the portable signal and the reason the plans are captured at all.

---

## Process

### Never hand-set a date that the clock can supply

Experiment E was launched with `REGIME_DATE=20260913` typed by hand, on 2026-09-12 UTC. Every
run id it produced — `20260913-history-cap20` and fourteen others — therefore claims a date
one day later than the one on which it ran. The run directories were not renamed, because
reports and analyses already cite them; each manifest carries the true UTC timestamps, and
this note is the correction. Run ids are identifiers, but readers will take a date in one at
face value, so the scripts' `date -u` default should simply have been left alone.

### Make the matrix resilient, not atomic

A multi-hour matrix that aborts on cell five should leave four usable results behind. Each
(topology, design) pair is an independent container run writing its own file; failures are
recorded in the manifest and the matrix continues.

### Record what was actually running, not what was requested

`record_topology` inspects the live containers for image, CPU and memory limits after
startup. A script argument says what was *asked for*; only inspection says what *ran*.
