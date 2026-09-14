# Lessons learned

Things that cost time, produced a wrong number, or would have produced a wrong number if
they had not been caught. Written down so they are not rediscovered.

Findings **about databases** belong in study reports. This file is about **running the
experiments**.

---

## Measurement

### Repeating a timing window is not repeating database preparation

Study 01's original `-trials` repeats phases on one load. It cannot supply independent
loads or alternate design order, and its report loader keys by topology/design, so simply
adding trial files would silently retain only one. The v3 follow-up uses one result per
freshly loaded trial and a separate grouping reporter. Preparation, order and exact
charity keys travel with every result; largest/smallest plans bracket the timed reads.

### Detailed explanations need a home outside the final decision document

The owner valued the D2/D3 mechanism analysis but found the final document too long.
Methodology 11b preserves that depth in signed discussion companions, including exchanges
between analysts, while the final analysis links to them. Do not solve length by deleting
evidence, rewriting another analyst, or removing disagreements.

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

### Explain a design change through the work it actually changes

The D2/D3 review exposed a reporting gap: "charity copied down" named a logical schema
decision but omitted the two new indexes and query rewrites that made it useful. A reader
could reasonably mistake the 1.55x overall score for a general benefit of storing another
column. Describe the schema, index and SQL changes together, then show which queries
changed and which did not. Distinguish a geometric-mean score over isolated queries from
throughput measured under an actual workload mix.

### Read execution counters before repeating the intended explanation

An `Index Only Scan` node is not evidence of zero heap access: Study 01's saved D3 sum
plan reports 31,293 heap fetches. SQL comments and even an analysis template described the
intended zero-fetch path; those are hypotheses until the execution counters confirm them.
Record the sampled key and when the plan was captured. A plan for the largest charity
does not establish the exact speedup of a benchmark that samples all charities.

### Check what a reported baseline and error bar actually mean

The regenerated survey's maximum D4/D5 disagreement is 1.69x, correcting an earlier 1.56x
context note. Such a maximum is a descriptive control, not a statistical confidence
interval. The cache-fix mixed run also starts at 4:4 readers:writers, so its displayed
100% baseline is not an all-reader baseline. A mixed run that reduces reader count as it
adds writers cannot isolate write interference from its retained-throughput percentage.
Independent analyses should inspect these definitions before carrying headlines forward.

Evidence and database-specific implications are in the
[GPT-6 Study 01 analysis](studies/01-charity-tree/reports/analyses/20260912-study01--gpt-6--2026-09-12.md).

---

## Building study 02 (overbooking)

### Ask the engine what it is doing; do not recall it

While writing study 02's YugabyteDB setup, the note "READ COMMITTED silently runs as
Snapshot Isolation unless `yb_enable_read_committed_isolation` is set" went into four files,
from memory of older releases. A two-minute check on the pinned image —
`SHOW yb_effective_transaction_isolation_level` inside an RC transaction, on a node started
without the flag — said `read committed`. The claim was false for the version under test,
and it would also have cast doubt on study 01's YugabyteDB concurrency cells, which in fact
ran real RC. The harness now records the effective isolation in every result, and the flag
is kept only to pin the behaviour. Engine semantics that a design depends on are read from
the engine at run time, never written down from recollection.

### A timeout that stops new work but not retries in progress is not a timeout

The race gave each event 60 s: after that, no buyer starts a new booking. On YugabyteDB at
SERIALIZABLE, buyers already *inside* a booking kept retrying serialization failures — up
to 1 000 attempts each, some statements taking tens of seconds (`Timed out waiting
kResponseSent`) — and a 10-seat race ran for many minutes, with 14 597 errors and an
organiser edit p99 of 166 s. The deadline now also stops a booking from starting another
attempt (work already in flight finishes, so no commit of unknown outcome is created), and
each size tier has a time budget after which remaining events are reported as not raced.

### The same rule broken twice means the rule needs a mechanism

"Never edit a bash script that is currently running" was already written down. Study 02's
runner was edited mid-run anyway, by a tool that rewrites files in place. Only a dev check
was at risk, but a written rule that failed once will fail again. Study 02's runner now
re-executes itself from a temporary copy at start-up, which makes the original safe to edit
at any time.

### A race's throughput needs the right clock

The first race reported sales divided by the whole race's wall time. For a 10-seat event
with 32 buyers, most of that time is spent telling the other 22 buyers "sold out", so the
number measured the rejection path, not selling. Sales are now divided by the time to the
last sale, and the rejection latency is reported on its own.

### A negative control is the audit's proof

C1 (count, then insert, at READ COMMITTED) was added knowing it is wrong. On its first
PostgreSQL dev run it overbooked every race event on every tier — 39 events, 866 extra
seats — and the audit caught all of them while the correct designs came back clean. Without
that control a clean audit could have meant "correct" or "not enough contention". The report
now states per topology whether each control fired.

### Several sessions can share one repository and one container runtime

While study 02 was built, another session was writing a second analysis of study 01 and
regenerating its reports in the same working tree. Two consequences. First, the infra
scripts reuse container names (`pg-single`, `yb-n1`…), so one session's `up` silently
destroys another's database mid-run; study 02's runner refuses to start if they exist.
Second, a repository-wide "dirty tree" check would mark every run dirty because of someone
else's documentation edits; the check is scoped to the code that can change this study's
numbers, and commits include only the paths the committing session changed.

### A shared tail figure across unrelated queries is a hazard signal

On the dev runs, reads that share nothing — a primary-key lookup, a band aggregate, an
availability count — all showed p99 near 65 ms, on both engines. Unrelated operations
converging on one tail value point to the environment, not the designs: here, most likely
CFS quota throttling of 2-CPU containers. The client now records its `cpu.stat` per phase and
the runner captures the database containers' before and after each cell, so the report can
show the throttling instead of asking a reader to trust the suspicion.

The first measurement settled it, and not the way the suspicion was framed: the **client was
never throttled** (0 periods in every phase), while the **YugabyteDB container was throttled
in 40–65% of its CFS periods in every cell** — up to 720 s of throttled time in one cell. On
this laptop, YugabyteDB's throughput and tails are bounded by its 2-CPU quota. That applies to
study 01's YugabyteDB cells too, which ran under the same budget without this measurement.

### A server-side failure needs the server's logs, captured before teardown

A YugabyteDB cell failed with "the database system is shutting down" from the YSQL layer. The
runner had captured `podman logs`, which for `yugabyted` holds only start-up chatter — the
PostgreSQL-layer and tserver logs live under the data directory — and then removed the
container. The cause could not be established. The runner now copies those log tails out of
the container on failure, before teardown.

It reproduced in the `small` matrix, and the captured tail still was not enough: the last 400
lines of the tserver log were all deadlock-detector timeouts (78 438 of them). The cause was
found by grepping the full logs inside the container while it was still up: under C2's
SERIALIZABLE lock storm on a CPU-throttled node, the tserver's RPCs to its master stalled for
~29 s, the **YSQL lease expired**, and the tserver killed every SQL session. A fixed-size log
tail is a guess at where the cause is; when a flood precedes the failure, it guesses wrong.
Search the full log for the moment of failure, not the end of the file.

### Provenance that can fail silently will fail silently

The runner captured the repository version with `git -C "$REPO" … 2>/dev/null || echo unknown`.
Under Git Bash with path conversion disabled (lib.sh does this so podman works), git.exe could
not resolve `/c/extra/...`. The first tagged matrix therefore recorded `repo_commit: unknown`,
`repo_dirty: false` and a `run_tag` that was never created — a run that looked traceable and
was not. It was caught by reading the console within minutes, and the run is kept with an
`INVALID.md`. Two fixes: git gets a host-style path through `hostpath()`, exactly like
podman; and the runner now **refuses to start** when it cannot read the commit or the working
tree status, because a fallback value in a provenance field is a fabricated fact.

### Compressed timings must be calibrated per engine

The holds experiment compresses a real checkout (minutes) into milliseconds: a 250 ms TTL.
Calibrated on PostgreSQL, it was degenerate on YugabyteDB, where taking a hold itself took
longer than the TTL — about 96% of holds expired before the buyer could pay, and the
experiment measured expiry rather than late payments. One `-hold-time-scale` factor now
scales every hold timing together, which keeps their ratios, and therefore the share of late
payments, unchanged across engines.

### Rules an agent cannot see do not exist for it

The repository's instructions lived in `CLAUDE.md`, which only Claude Code reads. An OpenAI
Codex session worked in the same tree and, among useful work, rewrote another analyst's
signed file in place — against a rule it had no reason to load. The instructions now live in
`AGENTS.md` (read by Codex and other agents) and `CLAUDE.md` imports it, so every agent reads
one text. Rules also stopped assuming a vendor: identify the model in commit messages however
the tool does it, and name vendors only where provenance needs them.

### Parallel agents need separate folders and a shared lock

Two sessions in one working folder edited `CONTEXT.md` and `LESSONS_LEARNED.md` at the same
time, and each had to avoid committing the other's half-finished changes. Branches alone do
not fix this — a folder has one checked-out branch — but git worktrees do: one folder and
branch per agent, one repository. Measurements are the part that cannot be parallelised at
all: container names are shared by every study and the machine's cores by every run. A lock
that lives where the contended resource lives (a podman volume, created atomically, visible
to every worktree, shell and agent) turns that rule from a convention into a refusal. It was
tested for a second acquirer, a guarded `down`, nesting, normal exit and SIGTERM, which
releases only after the current foreground command returns.

---

## Building study 03 (reserved seating)

### Name the domain regime in the question

Study 02 used "seat" for a unit of capacity: `seat_no` was an admission number assigned by
the system, and its holds claimed capacity, not a place. The owner read it as a study of
physical seats and asked for "the same study, with marked seats". The difference is not
cosmetic: with marked seats the invariant becomes a single-row fact and the hard problems
move to holds, expiry and conflicts over one chosen seat. State the regime (general admission
or reserved seating) in the first sentence of a study, and define the domain words before the
designs use them.

### Timestamps compared by a gate need one precision

The first PostgreSQL gate failed S1 by 1 ms on the loaded holds' expiry. The load clock came
from `SELECT now()` (microseconds) and the offsets added to it were milliseconds, so the
database and the Go truth disagreed below the millisecond. The load clock is now truncated to
milliseconds before anything is derived from it. The gate caught it; a timing run would
not have.

### A warmup on a finite pool can consume the measurement

Isolated hold, release and cancel benchmarks draw from a fixed pool of seats, holds and
tickets. A warmup by duration drained the `tiny` release pool on PostgreSQL before the first
trial, which then measured 0 releases without an error. Warmup on a finite pool is now one
trial's worth of operations, so warmup and trials stay inside the pool.

### A probe for an outage must run inside the outage

The first E1 sweeper-outage probe found 0 of 0 unavailable seats: the sweeper resumed on its
timer before the probe had read, so the negative condition could not show. The outage now
ends only after the probe has run (4 of 4 and 9 of 9 on the next check). A detector that
cannot fire is not evidence of a correct design.

### A compare-and-set loop without backoff measures its own spinning

L2 (a section document written with a version compare-and-set) retried a lost CAS
immediately. On 10-seat races the buyers spun against each other: 21 seats/s. With the same
jittered backoff used for retryable errors, 206 seats/s. A lost CAS is a retry like any
other, and a design's number without backoff describes the harness.
