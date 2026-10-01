# Lessons learned

Things that cost time, produced a wrong number, or would have produced a wrong number if
they had not been caught. Written down so they are not rediscovered.

Findings **about databases** belong in study reports. This file is about **running the
experiments**.

---

## Measurement

### A model handoff must preserve decisions, evidence and the next role

The owner separates expensive reasoning from routine execution: HIGH writes an
Execution Handoff, LOW carries it out, and HIGH reviews validation and results.
A bare request to switch models loses the decisions the executor needs. Commit the
handoff with exact scope, steps, acceptance criteria and permitted choices; log unmapped
decisions as Escalation Required. Expected waits remain execution work. Each iteration
names the next role and any agreed model/effort, and context distinguishes a model-switch
checkpoint from a finished, merged task. These rules now live in AGENTS.md for all studies.

The final review must verify the actual tag targets and main ancestry: a committed
receipt can only name its own final tag as intended before that tag is created. EH-01's
receipt retained this forward-looking wording even after integration succeeded. The
next HIGH review records the observed commit/tag outcome without rewriting that
historical receipt. When only a documentation merge remains, HIGH can give LOW explicit
preservation and completion conditions so a routine lock wait does not create another
unnecessary review cycle. Keep the already merged implementation distinct from any
new review documents still awaiting integration.

### Isolation needs an explicit integration step

A completed Study 01 branch contained the results and updated rules while `main`
still exposed the old context. Committing and tagging made the work traceable but
did not make it the default for the next task. The owner's 2026-09-14 rule now pairs
an isolated worktree at every task's start with a local merge at completion. Integrate
current `main` in the task worktree, preserve concurrent contributions when resolving
conflicts, validate, then update `main` and verify commit reachability. Wait before
changing shared scripts used by a running benchmark. A tagged branch alone is not a
completed integration; remote pushes and pulls remain separate owner actions.

Check the live lock's branch/worktree labels as well as `git worktree list`. During
the 2026-09-15 HIGH review, the peer Study 03 matrix still identified the main checkout
despite the standing isolation rule. Keep working in the owned worktree, leave the
active run in place, and wait for a safe shared-main update; do not assume another
agent has already isolated its work merely because the rule exists.

### Worktree isolation prevents collisions; reconciliation prevents a patchwork

Separate worktrees let agents edit independently, but they do not make two completed
branches read like one project. The agent merging later owns that integration: bring
current `main` into the task branch, reduce living documents to one current status,
combine duplicate lessons without losing either session's evidence, and keep rules,
terminology, READMEs and indexes consistent. Preserve attributed analyses, discussions,
handoffs, progress logs, escalations, reports, results, manifests, tags and history.
A real disagreement gets a separate signed response; it is not smoothed away during
reconciliation. This Study 01 closeout encountered Study 03 first as a running matrix
and later as a completed signed study; retaining both as current status would have made
the repository read like stitched session notes. The canonical hard rule is in AGENTS.md.

### A study-specific finding is not automatically a future-study requirement

Study 01's v3 context and discussions recorded resource, placement and concurrency
lessons, but the owner had to ask whether every later study would apply them. Keeping
context current does not by itself establish minimum experimental coverage. The
2026-09-14 clarification is now explicit in
[AGENTS.md](AGENTS.md#required-comparisons-for-future-studies), with a protocol coverage
map and methodology guidance. Include applicable rollup/rolldown/embedding controls,
multi-node colocated/non-colocated placement, and optimistic/pessimistic concurrency.

### Equal aggregate resources can hide a bottleneck at the query endpoint

Study 01 v3 gave three database nodes the same total CPU/memory as one larger node,
but sent cluster SQL through only one smaller endpoint. The aggregate budget alone
could not distinguish replication work from the distribution of query capacity.
Calculate per-host, per-node and client budgets, account for supporting processes,
and record endpoint distribution and per-node throttling. Treat client readers/writers
and connection-pool capacity as separate sizing decisions. Calibrate and freeze them
across controlled pairs; verify offered demand and rejections instead of assuming a
fixed worker count can drive every topology. A formula supplies a starting budget,
while measurements establish its limits.

### A memory comparison must hold the generated dataset fixed

The first v3 growth matrix used medium/history-multiplier=1 at 3 GiB and multiplier=2
at 256 MiB. A D6/D3 donor-read reversal there could involve both resource pressure and
the different dataset. The follow-up repeats multiplier=2 under both configurations
in one run, with the same seed and mutation counts. Keep the original observations,
but require this control before attributing the difference to the memory configuration.

The completed matched run `20260914T104721Z-v3` reproduced the reversal in all
three trials per configuration with identical initial dataset summaries. Budget
settings also changed D6's allocated footprint despite fixed logical data. Report
the configuration effect separately from its unproven mechanism, and disclose that
the resource blocks ran in a fixed order even though design order alternated.

### Zero operation errors does not mean offered demand was served

The v3 arrival tests can have no SQL errors and exact acknowledgement reconciliation
while rejecting most offered requests at the bounded queue. Local node-stop trials
also preserved acknowledged counts while successful response p99 reached about fifteen
seconds. Report offered, rejected, accepted, completed, errors, drain time and latency
together; use a no-fault condition at the same arrival rate to quantify the fault's
incremental effect. Count preservation alone is not an availability or latency result.

### Result hashes need an explicit checkout byte policy

This Windows repository has `core.autocrlf=true`. New container-written JSON had LF
both in the working tree and index, but no attribute protecting it against conversion
on a later checkout. Study 01 now pins v3 result JSON to `text eol=lf`; otherwise an
unchanged Git revision can produce a different raw-byte input digest on another
checkout. The attribute is scoped to new v3 runs so historical bytes are not silently
rewritten. Verify `git ls-files --eol` and `git check-attr` when adding new run formats.

Study 03 then met the same thing from the other side, which shows what it costs when the
attribute is missing: its `small` matrix hashed to `a56ce92ce38b8204` in the checkout that
produced it and to `b8913f899a9584ae` in a fresh worktree, with every measurement identical
and no result file changed in git. A report regenerated there would have published a digest
no reader could reproduce. Studies 02 and 03 are now pinned in the repository-root
`.gitattributes`, and regenerating the four study 03 reports restored each run's original
digest. When a digest changes, check the bytes before the data: the same measurements must
not hash two ways.

### Database preparation can change during a short read experiment

The v3 ANALYZE-only diagnostic started a D15 sum plan with 31,293 heap fetches and
ended with 218. Background maintenance was allowed to run; an initial plan does not
describe every timed execution. Explicit vacuum preparation produced zero heap fetches
in both bracketing plans for all five D15 trials and both targeted charity sizes.
Record preparation and before/after execution counters, and describe an evolving
ANALYZE-only state as a diagnostic rather than a fixed dirty-page control. See run
`20260913T125342Z-v3` and its signed mechanism discussion.

### Repeating a timing window is not repeating database preparation

Study 01's original `-trials` repeats phases on one load. It cannot supply independent
loads or alternate design order, and its report loader keys by topology/design, so simply
adding trial files would silently retain only one. The v3 follow-up uses one result per
freshly loaded trial and a separate grouping reporter. Preparation, order and exact
charity keys travel with every result; largest/smallest plans bracket the timed reads.

The new catalogue-control test immediately caught a source-generation error: an
unanchored match took the documentation's example `-- name:` for the first statement.
Start catalogues at a line-anchored real query declaration, and test the exact unchanged
SQL between controlled variants before starting databases.

An arrival generator must conserve demand: offered = accepted + rejected, and accepted =
completed + failed after drain. Its overload test uses a service slower than the offered
rate and requires visible rejections and queue delay. Successful-request percentiles
must be shown alongside failure/rejection counts, not reported as all-request latency.

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

### A fixed arrival rate measures compliance, not capacity

Study 01's recency hot-donor sweep varied the writer count (1, 4, 8, 16) while holding the
offered arrival rate at 500/s. Every healthy cell completed exactly 10,000 requests at
500.0/s — the offered rate — at every writer count, because the rate, not any design, was
the binding constraint. The sweep therefore answered "who can keep up with 500/s?" and not
"how much can each design do?", and the only design it distinguished was the one that
*failed* to keep up (D23's compare-and-set, which lost 22–56% of its offered load at 16
writers). An open-loop arrival experiment that fixes the rate across a concurrency sweep
measures a threshold, not a curve. Either scale the rate with the worker count, or set it
well above the expected ceiling and read the drop rate — and say which of the two you did,
because "500/s, no errors" looks like a result and is not one.

### A controlled pair can be clean on one axis and confounded on another

RECENCY.md paired D20 (a flag on the child) against D22 (a rollup on the parent) to ask
where a derived fact should live. On the **read** side the pair is exactly one decision:
both answer the same four statements, one from a partial index on `donation`, the other
from an index on `person`. On the **write** side it is not: D22 inherits D4's whole rollup
package — five maintained columns across two parent tables — while D20 maintains one
boolean. So D20's 6,219 inserts/s against D22's 3,436 prices four extra columns and a
second hot row, not the flag-versus-rollup decision the pair was built to isolate. The
protocol's own controlled-pair table did not distinguish the two axes, and the analysis had
to state the confound instead of a result. **When registering a pair, check it separately
for each axis it will be measured on**; a pair can need a third design (here, a
`last_donation_at`-only rollup) to become one decision on the axis you care about.

### A negative control needs the RIGHT kind of contention, not just SOME contention

D21's unguarded flag-race control (RECENCY.md, v4) fired reliably under 8 ordinary
connections spreading inserts across 500 donors at maximum closed-loop throughput — but
did not fire in the paced, open-loop "arrival" experiment at 8 hot-donor writers, even
though every insert targeted the SAME single donor. Raising that experiment's writers to
16 made it fire. The scheduled arrival rate (500/s across 8 workers ≈ one request every
16ms per worker) gave each transaction more room to commit before the next one started,
narrowing the window the race needs; a tight retry loop with no pacing does not. A
negative control that fails to fire under one load shape is not evidence the design is
safe — it may only mean that load shape does not create the overlap the control depends
on. Where the handoff maps a fallback (here: more writers), try it before concluding
anything; where none is mapped, the right question is what property of the load actually
matters, not just how much of it there was.

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

### "Publish only after the commit" does not make cache-aside safe

A reader that begins its cache fill *before* a writer's invalidation holds a state that is committed and
**superseded**. Publishing it afterwards puts back exactly the value the writer removed. Study 05 measured
~14 000 stale-after-ack reads in a strict cell whose only rule was "publish after the commit".

What closes it is a per-key **invalidation fence in the cache**: the writer advances the fence and removes
the value atomically; a fill captures the fence *before* its database snapshot and stamps it on the entry;
a publish is refused unless the key's fence is unchanged. Two details are load-bearing, and both were found
by failing cells rather than by reasoning:

* a strict writer must fence **before and after** the commit — one fence leaves a window in which a reader
  captures the *new* fence and then reads the *pre-mutation* state, so the check passes;
* the ledger's freshness requirement must move at the **acknowledgement**, not at the commit, because the
  contract is written against the ack and a strict writer fences before it acks. Treating the commit as the
  boundary made correct cells look wrong.

Study 05's `publishes_refused_by_fence` (61–85 refusals per cell) is the evidence that the race is real, not
theoretical. The DB version token is **not** what fixes this; the fence is, and it works for a legacy model
that cannot be modified. A version token is still what lets a multi-instance *local* cache validate a hit.

### A harness that mis-attributes one write fabricates correctness findings

Study 05's hotspot phase built a mutation and then overwrote its `PersonID` to force it onto a hot donor.
For a correct/delete/reassign the donation belonged to the donor chosen inside the builder, so the database
changed a row the ledger attributed to someone else and the two diverged. The gate reported it as an
**impossible cache value** and failed a fault cell whose own logic was correct.

The general rule: a mutation must be **built for** its target, never adjusted afterwards. The same class of
bug appeared twice more in the same session — a correction written as an absolute amount made two concurrent
corrections order-sensitive (fixed by making both the SQL and the ledger relative), and a *single* pending
state slot lost one of two simultaneously-committed states (fixed with a reference-counted set).

### A cache that was never started looks exactly like a cache that is merely cold

Study 05's runner called `need_redis` for a helper named `needs_redis`, so no cache container was ever
started — and because a cache error degrades to an authoritative read by design, every cache scenario
"passed" its phases while measuring nothing but the database. Only the fault that demanded a working lease
failed, and it looked like a lease bug for four dev iterations.

Two lessons, and the cheap one first: check the helper name you are calling. The durable one: **a scenario
that needs a cache must prove the cache answers before it measures anything** (a `PING` gate in the cell,
not a promise in the runner). Absence must never be indistinguishable from a cold start.

### A zero value that means "unset" is a cache that does not exist

`hasCache()` was `backend != "none"`. The zero value of the backend field is `""`, so three reference cells
that simply did not set it were treated as cache scenarios and panicked on an empty instance list instead of
running as the no-cache baselines they were. Ask for the backends you support explicitly
(`backend == memory || backend == redis`); never test a field against one sentinel value when the field has
a zero value that is neither.

### A phase that measures a different deployment must not leave state behind for the next one

Study 05's three-instance phase measures a different deployment through its own instances. Its writes cannot
invalidate the base adapter's stores, so entries written before it survived it and a later strict read
looked like a violation of a deployment that never produced it. The phase now starts and ends cold. Any
phase that swaps in a different set of components owns the state it leaves.

### Write-phase counters and end-state gauges are not the same reading

Study 05's cache snapshot happens once, at the end of a cell, after the fault phases have flushed — so
`items` and `resident_bytes` describe an empty cache while `evictions`, `fills` and `publishes` describe the
whole cell. Both are worth reporting, but a reader who takes the gauge for steady state is misled. Read
gauges at the moment they mean something, or say in the report which moment that is.

### A bypass read cannot be stale — if one is, the accounting is wrong, not the design

Study 05's strict failures were first read as a cache-design problem. Instrumenting the first few
stale reads in full showed that some of them came from an **authoritative database read** returning
a state one and two versions behind the requirement. A bypass read consults no cache and runs in one
repeatable-read transaction, so no cache mechanism can produce that: the *ledger's* freshness
requirement was ahead of the database's own state, i.e. the harness was manufacturing violations.

The lesson generalises past caches: **each wrong-read bucket must name a mechanism that could
physically produce it, and a bucket whose mechanism does not exist in that path is a harness bug.**
Practically: keep the first few failing observations in full — source, both fences, both sequences,
versions behind, whether a write overlapped — in the result file. A count tells you a rule failed;
only the observation tells you which rule, and the source field is what separates "the design is
wrong" from "the accounting is wrong". Add the assertion that makes the impossible case fail loudly
(`Required()` must never be ahead of the database) before drawing any conclusion from the counter.

### A mutation must enforce the owner it observed, or the ledger and the database diverge

Study 05 built each write from the ledger's state and then applied it to the child row **by id
only**. Under deliberately concurrent writes a second writer could move or delete that row in
between, so the statement acted on a row that belonged to somebody else (or matched nothing and
succeeded silently) while the ledger applied the change to the person it had observed. The two
diverged, and the harness then reported impossible cache values and audit mismatches that belonged
to no design.

Two rules follow, and both are cheap:

* **Every statement that acts on a child row carries the parent it was observed under**
  (`WHERE donation_id = $1 AND person_id = $2`). A refused write is a *correct* outcome, not an
  error.
* **A statement that matched no row means nothing changed** — so the ledger must not change either.
  Report it as its own outcome (`errNoEffect`), acknowledge it as a no-op, and count it, so a
  workload that has quietly stopped doing work cannot hide inside a clean run.

One more, from the same repair: when a mutation becomes *relative* in the ledger (a correction is a
delta), it must become relative in **every** SQL catalogue that pairs with it. Three reference
schemas still used an absolute `SET amount_cents = $2` and silently wrote the delta as the amount.

### A ledger that disagrees with the database manufactures correctness findings

Study 05's most expensive lesson, because it cost four rounds of investigation and briefly made the study's
headline a harness artifact. The ledger's history must be appended in the **database's commit order**. Two
concurrent writers on one key can commit in one order and reach the handler in the other, after which a
legitimate read looks "one state behind" and a correct design is failed.

Two diagnostics are worth their weight, and both were what finally settled it:

* **A design with NO cache showing the cache's symptom proves the harness is at fault.** The violation
  reproduced in a plain database-only reference cell, which no cache mechanism can explain.
* **Separate the sequential case from the concurrent one.** The acknowledgement checker writes and then
  verifies from a SECOND database session; with one writer it reported zero violations in 150 attempts
  before and after the fix, which killed "the ack is early" and left "the bookkeeping is out of order".
  A writer verifying through its own pool can see its own commit and turn a real defect into a clean result.

The fix is small: serialise writes that touch one key (sharded mutexes, never across keys), which leaves
inter-key parallelism — the thing contention phases measure — untouched.

### The checker can have the bug it is looking for

Study 05's acknowledgement checker captured the freshness requirement *after* its verification read. A
concurrent writer could therefore move the requirement between the read and the capture, and the checker
manufactured a violation — repeatedly, and plausibly enough to be believed. The rule it violated is the
study's own central rule: **capture the requirement at the start of the read, before any I/O**. When you
write a checker, re-derive its invariants from first principles rather than trusting that it is immune to
the mistake it detects.

### A negative control is judged by its fault, never by the invariant it breaks

Two defects in one: the unsafe-write control had to lose a *counter* as well as break the version invariant,
so a run where the invariant visibly broke had the harness refusing to license the guarded designs; and the
controls were being subjected to the ledger and acknowledgement assertions, which they exist to violate.
A control's verdict is "did its own fault fire", with the evidence recorded; every other gate must skip it,
and the skip belongs in the report rather than in a silent branch.

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

### A new query's parameter has more than one binding site

Study 01's `readVals` binds the values every read query draws from, but it is not the
only fixed-values map in the harness: `ExplainAll` (plan capture) builds its own, and it
had no `since`/`until` when q13-q16 (RECENCY.md, v4) were added. Every design's `-cmd
full` run failed at "capturing plans" with "no value bound for parameter \"since\""
until this second binding site was found and given the same window. Adding a query
parameter means finding every place a fixed-vals map exists for that catalogue, not just
the one the read benchmark uses — grep for the sibling parameter names already in scope
(`charity_id`, `person_id`, `donation_id`) rather than trusting one call site.

### An embedded document's field names are a second source of truth

D6/D9/D10 store donations as JSONB with single-character keys (`t` for `donated_at`,
`i` for the id — see `load.go`'s `jsonDonation`), because a million elements' worth of
verbose keys is real bytes. New SQL against that document has to use the actual on-disk
key, not the column name it stands for. q13-q16 (RECENCY.md, v4) were first written with
`e->>'donated_at'`, which returns SQL `NULL` for every row rather than erroring, so
every check silently reported zero matching donors until the correctness gate caught it
(exactly what the gate is for). The existing queries in the same file already used
`e->>'t'` correctly; checking them first would have caught this before running anything.

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

### A validator that only checks a document against itself cannot detect a wrong identifier

The evidence registry passed `tools/evidence validate` with `0 errors` while its `cells` field
named designs that do not exist in the reports it cites and its `failed_cells` contradicted them.
Neither is detectable the way the checks were written: a cell name was compared only against the
registry's own list, and `failed_cells` against the registry's own selected cells, never against
the cited run. A computed audit of the twelve numeric claims found 11 with at least one declared
cell name absent from their report and 36 of 43 names resolving nowhere; the worst mapping is not a
stale alias but a shifted and partly reversed one (registry `d4-rollup-application` where the
corpus's D4 is the trigger design), so the declared winner of a headline claim can point a reader
at a different mechanism. The lesson is not "add more fields": it is that a closed-world check
proves internal consistency and nothing else. Every identifier that claims to resolve to an
external artefact needs at least one assertion made **against that artefact**, and a green
validator must not be read as evidence about the source.

Evidence: `book/evidence/claims.json`; `studies/01-charity-tree/reports/20260912-small.md`;
`studies/02-ticket-booking/reports/20260913T021206Z.md`; brainstorm
`brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions`, tag
`repo/book-evidence-and-model-brainstorm-v1`.

### Two matched cells can swap identities in prose while every number stays real

The all-green Study 05 analysis prints its three-instance wrong-read table with the legacy and
owned rows exchanged: 38 749 / 87.81 % belongs to `o:opt:red:wth:rel` (the report's own line) and
35 287 / 87.24 % to `l:na:red:wth:rel`, but the signed prose attributes them the other way round.
Both rates, both denominators and both cache-hit counts are measured facts; only the two model
labels moved, so no numeric check catches it and a reader would conclude the wrong model collapses.
**When an analysis compares two designs that differ by one dimension, check each row's identity
against the report's own cell id before quoting a rate** — a derived claim takes the model name
from the prose, but the report and the result JSON are the only arbiters of which cell produced it.
Evidence: `studies/05-cache-consistency/reports/analyses/20260921-cache-consistency-allgreen.md`;
`studies/05-cache-consistency/reports/20260921T-survey3.md`;
`studies/05-cache-consistency/results/20260921T-survey3/pg-single/{l:na:red:wth:rel,o:opt:red:wth:rel}.json`.

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

### Human capability labels need one canonical machine value

People use several familiar words for the same AI-pipeline relationship. Letting each
word become a separate queue value would fragment eligibility checks and historical
queries, while borrowing `primary`/`replica` would collide with this repository's data
topology vocabulary. Accept user-facing aliases only in an explicit session declaration,
normalize them immediately to `HIGH` or `LOW`, and preserve the raw input separately.
Legacy `master`/`slave` can be parsed for compatibility without being generated or
recommended; `leader`/`worker` are the preferred human-facing pair.

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

### `RETURNING` shows what the update left, not what it erased

X1 (study 02 v2) writes its refund ledger row from the `UPDATE` that cancels the ticket:
`RETURNING …, customer_id`. That update has just set `customer_id = NULL`, and `RETURNING`
yields the post-update row on PostgreSQL 17 and on YSQL (PostgreSQL 15); `RETURNING OLD`
exists only from PostgreSQL 18. Every refund was logged without its buyer. Nothing errored,
and the design passed its gate on both engines. When a side record describes a change, take
the "before" values from something that still holds them. Here that is the ledger row of the
sale being reversed.

### An audit only ever run on state built to agree with it proves nothing

X1's ledger reconciliation audit was reported "consistent" on both engines. It had run only
at load, where the loader writes the ledger and the tickets from the same generated sales, so
it could not have disagreed. It never ran after a write, never after the race it exists for,
and it counted events per seat without checking whom they named, so the NULL-buyer bug above
was invisible to it. Treat a new audit like a negative control. Show it firing, on a real
defect or an injected one, and run it after the phases that can break the invariant, before
any "consistent" counts as evidence (EH-02 AM-03).

### A lesson has to be named in the next handoff to reach it

"A new query's parameter has more than one binding site" was learned in study 01's phase 1 of
this same task. Phase 3a's handoff (AM-02.4) did not cite it, and both studies' `ExplainAll`
again lacked the new window parameters. The mistake would have surfaced only as
`NOT CAPTURED` plans in a measured run. A planner writing a handoff should search
LESSONS_LEARNED for the mechanisms it touches (parameters, plans, audits, loaders) and cite
the entries by name.

### A NOT NULL column downstream of a NULL-producing bug turns silent corruption into an outright failure

X1's original `w_cancel_ticket` logged the refund's buyer from the cancelling `UPDATE`'s own
`RETURNING`, which yields the row *after* the update -- already NULLed. On its own, that
would have been a silent data-quality bug: a refund row nobody could attribute. But
`sale_event.customer_id` is declared `NOT NULL`, so every one of those inserts violated the
constraint and the whole cancelling transaction rolled back -- meaning X1 could not
successfully cancel a single ticket, a functional break far more visible than the bug that
caused it. Neither AM-02's dev checks nor its unit tests caught this, because none of them
exercised a cancellation on X1 (the isolated `write:cancel` op's default warmup exhausted its
finite pool before the measured window even started, so `0.0 ops/s` with `errors=0` looked
clean). A schema constraint can convert a data-correctness bug into a load-bearing failure
faster than any test written to check the data directly -- which is a reason to keep such
constraints, and a reason a "0 ops, 0 errors" write-benchmark line deserves a second look
before being read as "nothing happened here."

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

### A guarded statement that matches nothing cannot conflict

On YugabyteDB, under the race's load, a confirmation `UPDATE ... WHERE hold_id = $h AND status =
'held' AND hold_expires_at > now()` matched 0 of a valid hold's seats: 8 times in 40
10 000-seat races of one diagnosis, more often on three nodes. A `SELECT` in the same
transaction showed the hold valid, and the identical `UPDATE` re-issued in that transaction
matched every seat. The server logged nothing for that backend. Turning off expression
pushdown, wait queues or query-layer retries did not remove it (ER-01,
`results/devchecks/er01-*`). A statement that writes nothing gives the engine no write
conflict to detect: if its read is stale, the refusal is silent rather than a retryable error.
Two things made it findable: the ledger classifying every refusal by its margin, and a
diagnostic that repeats the refused statement before rolling back. A standalone probe of the
same statements never reproduced it; only the full workload did.

### `podman machine ssh` from Git Bash can leave a file named `NUL`

On Windows, `podman machine ssh` wrote the VM's host key to a file called `NUL` in the working
directory, a reserved name that breaks `git add`. Delete it after use (`rm ./NUL` in Git Bash), or
read the VM's state through `podman` commands instead.

### A lag measured against loaded state measures the load

Study 03's lifecycle report showed a release lag of about 45 000 human minutes on the 10-seat
tier for every design with lazy expiry (dc3, dc12). The first tier's sweeper was the first to
release the dataset's holds that were loaded already expired, and their "lag" was the age of the
load. The numbers looked absurd only because the time was compressed; at real speed the same
bug would have reported a plausible half hour. Lag now counts only holds granted during the
phase. Any metric measured from a timestamp should say which rows may contribute to it.

### Instrumentation must not be able to delete a design's results

Study 03's `small` matrix lost the lifecycle phase of four YugabyteDB cells. In three of them no design
statement had failed. The harness's own sold-seat monitor query timed out while the node was saturated,
and the harness treated the first monitor error as fatal. In the fourth, a harness reload's `ANALYZE` did
the same. The designs' own statements, which timed out at the same moments, were retried and counted as
errors. An observer query should have the same tolerance as the workload it
observes (retry, count, move on). Otherwise a busy engine produces missing numbers, and a report shows
"failed" where it should show "slow".

### A classification that needs a diagnostic is silent where the diagnostic does not run

AM-01 split early rejections into transient and persistent by re-issuing the refused statement. L2
(a section document checked in the application) has no statement to re-issue, so its refusals were
counted as "0 transient". They had the same shape as the transient refusals of every other design: a
read 20–100 ms after the hold's commit that did not show it. The report tables printed `n/a`, but the
TL;DR listed L2 among designs with invariant violations. When a class depends on a diagnostic, report
"unclassified" wherever the diagnostic cannot run, in every place the number appears.

## Operational reports and the recency question (EH-02)

### Nothing reconstructs commit order; write the check so it does not need to

Learned by study 02's X1 ledger and decided in AM-04; recorded here in this session's words,
as that amendment asked. The ledger needed to know *which of two events for a seat came
first* — a refund must name the buyer of the sale it reverses, and the reconciliation audit
had to read each seat's history as sold → cancelled → sold. Two candidate keys were tried
and both are wrong, on different engines:

- **`now()`** is the transaction *start* time in PostgreSQL, not the commit time. Under a
  CAS retry loop a later-starting transaction can commit first, so ordering by `at` inverted
  a cancel and the resale that necessarily followed it (seat 87/10: `at` said the second
  sale preceded the cancellation that freed the seat).
- **A `BIGSERIAL`** hands out cached blocks per connection on YugabyteDB, so `sale_event_id`
  is not commit order there either — 53–246 attribution "problems" per run, identical across
  repeats, and `ALTER SEQUENCE … CACHE 1` did not change what YSQL reported.

The fix was not a better key. The refund reads the buyer **from the row it is about to
clear**, under `FOR UPDATE`, in the same statement, and the audit was rewritten order-free
(no `LAG`, no `ROW_NUMBER`): a buyer may not be refunded more times than they bought a
seat, and a live sale must appear in the ledger verbatim. Both checks hold after every race
and churn tier on both engines, including at the 128-buyer contention that produced the
false positives. **An audit that needs to know "which record came first" should be
rewritten so it does not need to know.**

### A finished run is not a recorded run

The previous session's 3b-1 (45 cells, 2 h 49 min) completed at 02:38:52Z and its results
were left untracked when the session ended. The `run/…` tag existed and pointed at the
producing commit, the report existed, and none of it was in the history: `git status` in
that worktree showed 185 untracked files. A reader of `git log` would have seen the run's
tag and no data. **Commit a run's output before starting anything else after it**, and treat
"the tag exists" as no evidence at all that the results were committed — the runner tags
the *code* before the run, not the output after it.

### A worktree link that works for one shell can be invisible to the other

This repository is worked from WSL and from Git Bash on Windows, and a worktree's link form
decides which of them can use it. The 2026-09-20 session's worktree held Windows-form paths,
so WSL git called it "prunable"; `git worktree repair` then wrote `/mnt/c/...`, which Git
Bash's `git.exe` cannot resolve — repairing for one shell broke the other. The form that
works for both is a **relative** `gitdir` in the worktree's own `.git` file
(`gitdir: ../../.git/worktrees/<name>`). The *admin* side (`.git/worktrees/<name>/gitdir`)
should stay absolute for whichever shell owns the checkout: a relative path there makes
`git worktree list` print a relative path and mark the worktree prunable, and a later
`git worktree prune` in any session would then delete the entry.

### An absent podman client looks exactly like a held lock

`run_lock_acquire` creates the lock volume and reads any failure as "someone else is
running". On WSL, where this branch's `infra/lib.sh` had no client, `podman volume create`
failed with "command not found", the diagnostic then failed too, and the runner printed
`another benchmark holds this machine's lock: (lock volume vanished)` before refusing. The
message names the wrong cause. `(lock volume vanished)` means "the create failed, and so did
the inspect" — check `command -v podman` before believing another session holds the lock.
Main's resolver (`repo/wsl-podman-bridge`) fixes the client; the message is still ambiguous.

### A control that is not in the run cannot fire

Study 02's reports matrix (3b-1) has phases `verify,explain,read`, so the two negative
controls had no contention to fail under and the report says "did not fire — this experiment
did not generate enough contention". The race pair (3b-2/3b-3) measured only P3 and X1, so
its control line reads "**0 of 0 fired**". In both cases the gate passed and the audits were
consistent — and in neither case has the audit been *shown* to catch overbooking in that
regime. Methodology 5a's point survives contact: a correct design passing beside no wrong
design is weak evidence. **Include at least one control in every run whose correctness is
claimed**, even when the run's question is throughput.

### A design's measured speed can follow its position in the cell order

The same P3 → X1 pair was measured twice, changing the buyer count and the design order
together: at 32 buyers with P3 first, X1 looked 2–3x *faster* in the race; at 128 buyers
with X1 first, P3 was 1.1–2.3x faster with tight trials. The tell was in the report's own
spread column — P3's three trials disagreed by 92–283 % in exactly the cells where it
looked slow. A third sighting came from a read-only run: X1's buyer reads were ~2x slower
than P3's including a primary-key lookup whose plan and buffer counts were identical
(`Index Scan using ticket_pkey`, `shared hit=4`), i.e. the whole cell was slower per
operation. **A pair measured once per cell cannot separate the design from the machine.**
Alternate the order between arms, or better, interleave the two designs tier-by-tier inside
one cell; and when a pair table prints a ratio, read the trials' spread before the ratio.

### A time guard is a pre-run decision; re-check it against the actuals

AM-03.11 sized phase 3b at ~8 h from a calibration and set a 10 h escalation guard, which
was applied mechanically before the runs. The runs then took 2 h 49 min, 2 h 09 min,
2 h 17 min and several hours against 2.5 h, 1.7 h, 1.7 h and 2 h — because 128-buyer and
14-design tiers on a CFS-throttled engine cost far more than the small calibration cells
suggested. The guard's letter was satisfied and its spirit was not. **Size a run from
per-tier time budgets measured on the slowest topology, and re-estimate after the first arm
rather than after the last.**

### A derived claim can be narrower than its own run when a headline is copied across arms

The queue derives each book claim verbatim from a signed analysis headline. Study 02's
32-buyer arm (`20260921T183722Z`) has a body that states the ledger's YugabyteDB race cost as
"1.7x to 5.7x", but its headline prints "1.7-2.1x" — which is the *sibling* 128-buyer arm's
range — and its conclusion-3 cancel figures (575→271, 404→164) are the sibling's too, while
its own JSON gives 579→270 and 485→157. The derived `claim-07` therefore understated the
cost by widening nothing, and a reader trusting the headline would have been wrong by 3x at
the 10-seat tier. **When two arms of the same pair are written in one session, re-derive
every headline number from the arm's own result JSON before taking it as a claim**; a
headline is not evidence, the run is.

### A one-shot load measurement is a machine reading, not a design cost

The ledger's `pg-single` load cost was measured at 6410 ms against P3's 826 ms (7.8x) in the
Study 02 reports run, and at 1062/1106 ms against 881/868 ms (1.2x) in the two sibling runs
of the same pair, with an empty harness diff between the commits. The storage delta was
stable (+33-35 %) across all three; the load multiplier moved 6x. **Quote storage and
schema-shape deltas from a load; do not publish a load-time multiplier from a single cell's
one-shot load**, and if the load cost matters, measure it in its own repeated load-only
cell.

## Building the AI work queue (queue v1)

### Relative worktree paths are portable, but old Git cannot *list* them

The repository stores a linked worktree's `.git` link and the
`.git/worktrees/<name>/gitdir` back-pointer as relative paths, so the same registration
resolves under WSL Git, native Windows Git and a Podman bind mount. Git reads relative links
even when it does not write them — `git -C <worktree>` works and Windows Git 2.53 handles
them — but **WSL Git 2.43 resolves the back-pointer against the current directory in
`git worktree list`, reports the worktree `prunable`, and prints the path relative**, so a
tool that checks "is this directory a registered worktree?" through `worktree list` decides
the worktree does not exist. Detect it from the filesystem and by opening it
(`git rev-parse --git-common-dir` inside the directory), not from `worktree list`.

### A container that mounts one worktree hides the repository

Mounting only the invoking worktree put `<common>/.git` outside the container, and every git
command inside failed with `not a git repository: <common>/.git/worktrees/<name>`. Mount the
**common repository root** at the same absolute path the shell uses. The same check exposes a
second trap: task records store `canonical_worktree` relative to the common root, so a CLI
running inside a linked worktree must resolve it against `--git-common-dir`'s parent, not
the current working tree.

### The container has no git author

A commit made by a CLI inside the pinned image failed with git's "Author identity unknown":
the repository's `user.name`/`user.email` lived in the host's global config, which the
container does not see. Forward `GIT_AUTHOR_*`/`GIT_COMMITTER_*` from the host repository's
own identity in the wrapper; do not invent an author inside the tool.

### An event that cites a file must commit the file

The queue's `publish`, `review` and `escalate` wrote `task.json`/`REVIEW.md`/`ESCALATIONS.md`
and then committed only the event and `STATUS.md`. The files stayed untracked while the event
claimed they existed, so a fresh clone (and, later, the merge) would have lost them. The test
suite caught it only because a worktree checks out a *committed* tree: `TestFullLifecycle`
and `TestFreshCloneReconstruction` failed with a missing `task.json`. **When a record names
evidence, commit the evidence in the same step as the record**, and let a test read the
artifact from a clean checkout rather than from the author's working directory.

## Building durable AI deliberation

### A guard at command entry does not prove ownership at commit

The first brainstorm submission implementation checked its slot claim, then waited for
the local-main integration lock and wrote the contribution. During that wait the lease
could expire, another leader could recover the same slot with a higher epoch, and the old
writer could still commit before its final ref deletion failed. **Revalidate and
CAS-refresh ownership immediately before the durable mutation, under the mutation's
serialization lock.** The final implementation extends the slot lease under the archive
lock before writing, then deletes exactly that refreshed ref after commit. The same
review found the administrative counterpart: cancellation must reject a live claim but
CAS-clear an expired one, or an abandoned lease becomes a permanent administrative
blocker.

## Building the book toolchain (Typst)

### A schema-valid registry can still resolve to nothing

`tools/evidence validate` on registry v1 reported 0 errors while 36 of its 43 declared cell
names appeared nowhere in the reports they cited. Schema validation checked two internal lists
against each other, not the claim against the corpus. **A validator earns its keep only when
its subject resolves to something outside itself.** The v2 resolver requires each support key's
token to appear verbatim in the cited report, and the old failure is pinned as a regression
fixture (`v1-diagnostic` must print 12/11/43/36) so it cannot return silently.

### `grep -q` in a pipe under `pipefail` can invert a true match

The book build tested whether the evidence digest appeared in the PDF text with
`printf '%s' "$TEXT" | grep -q "$DIGEST"` under `set -o pipefail`. The text was larger than the
pipe buffer, `grep -q` exited on the first match, `printf` took SIGPIPE, and the pipeline's
non-zero status made the test false — the digest *was* on the page while the manifest said
`false`. **Use a here-string (`grep -q "$needle" <<<"$haystack"`) for presence tests in
`pipefail` scripts**, and verify a build artefact by reading it, not by trusting the exit code
of the thing that produced it.

### A verification claim must be computed, not asserted

The first build reported "3 fonts" and "0 embedded" because `pdffonts`' `emb` column moves with
the font type (CID Type 0C inserts two extra fields). The fix was to match the stable embedded
triple (`yes yes yes`). The same run counted 26 "claim tokens" for a 23-claim registry because
the regex matched tags and prose (`v2-protocol`). **When a manifest states a number, derive it
from the artefact and make the build fail when it disagrees** — the count is then a check, not a
decoration.

### Pin the renderer, and let it carry its fonts

Typst ships its default fonts inside the binary, so the PDF embeds Libertinus Serif and DejaVu
Sans Mono with no host font installation. Pinning the image by immutable manifest digest
(`ghcr.io/typst/typst@sha256:032e…`) fixes version, font set and rendering behaviour together;
the manifest records the digest so a later reader knows which toolchain produced the page.

### A committed binary artefact must be exempt from the text line-ending rule

The root `.gitattributes` pins `book/** text eol=lf` so sources hash identically across Windows
and WSL. That rule also caught the generated PDF: git normalised CRLF to LF on commit, so the
committed blob's sha256 (`a03b54…`) did not equal the PDF's real hash (`da433f…`) that the build
manifest records. **Exempt binaries explicitly** (`book/dist/**/*.pdf -text`) and verify with
`git cat-file -p :path | sha256sum` that the stored blob matches the recorded hash — a manifest
whose hash does not describe the committed bytes is worse than no hash.

## Running the Study 04 v2 controls (2026-09-23)

### A shared-ledger harness needs its full phase sequence

The contention phase audits the design against one ledger built for the whole cell. Running
`--phases verify,explain,contention` (skipping `read,write`) makes that ledger expect a state the
skipped phases would have produced, so the audit fails with
`INV-3/INV-10: installation N is missing key …`, the runner records the cell as **failed**, and the
contention numbers are unusable. The failure is a harness-ordering artefact, not a design defect:
re-running with `verify,read,write,contention` passed every cell. **When a harness keeps one ledger
per cell, a reduced phase set must be a supported mode, not an accident** — either prepare the ledger
for the reduced sequence or refuse to run the subset. The five failed runs are kept in `results/`
and excluded from the analysis.

### The runner's dirty check counts generated reports, so later runs are not tagged

`run-study.sh` computes "dirty" with `git status --porcelain -- platform infra studies/<id>
:!studies/*/results :!studies/*/reports`. The pathspec exclusions do not suppress untracked
`reports/*.md` files, so after the first run every later run reports `repo_dirty: true` and the
runner **silently skips `--tag`**. No tracked code changed — `git status --porcelain
--untracked-files=no` over the code paths is empty — but the runs lose their automatic tags. The
tags were recreated by hand on the same commit. **A provenance check must not treat its own
generated artefacts as a dirty tree**, and a runner that declines to tag must fail loudly rather
than continue quietly. Studies 01–03 use the same pattern and should be audited.

### Identical SQL reproduces within ~7% across phases; a larger within-run spread has another cause

Running the same design (`n1_rows_indexed`) alone in three fresh-load phases reproduced every read
within 0.9%–6.7% (most ≤2%). That is the cross-phase instrument floor on this host. It is far below
the ~1.6x identical-SQL spread seen within the v1 run, so that larger spread is not run-to-run drift
— it is a within-run position/order effect, and it remains unmeasured because the runner cannot
place one design first and last in a single phase without overwriting its own result file.
**Measure the effect in the position domain before naming its mechanism.**

### A cadence arm that records delivery is not a staleness measurement

Study 04's cadence phase proves the scheduler fired (offered == started == completed, queue depth
within bound) and gives scheduling lag, but it does not stamp derived values, so it cannot say how
old a rollup was when a reader saw it. The protocol asked for staleness and rollup lag; the harness
had no field for it. **Delivering the measured half and naming the missing field is the honest
result** — a delivery rate is not a staleness distribution, and "quiet/steady/burst" describes the
update cadence, not the age of the answer. Closing it needs a harness change, not a bigger matrix.

### A predicted maintenance amplification can be absent below the size the mechanism needs

Study 04's `d2_doc_on_parent` embeds the JSON document in the parent row. The design catalogue
predicted that an ordinary metadata update (`w06_update_metadata`) would then rewrite the TOASTed
document and cost more than d1's short-row update. Measured at `scale small`: publication 865.9 vs
858.8 ops/s and metadata 1651.0 vs **1672.7** ops/s — no amplification at all. Below the TOAST
threshold the document is inline in the tuple, so rewriting it is nearly free. **A design hypothesis
about size-dependent cost must be tested at the size the mechanism needs**, and a negative result at
the wrong scale is a statement about the scale, not about the design. The d1 ↔ d2 pair now exists;
the amplification test belongs at the 500-entry tier and the large value regime.

### Crossing a TTL boundary in time is not the same as observing the TTL expire an entry

Study 05 v2 churn ran 600 s and 1800 s under load, crossing 2.00 and 6.00 real 300 s TTL boundaries
with zero wrong reads in up to 2.88M reads. But `expired_hard` was **0** in every cell: only
probabilistic early expiry fired (11.3k–37.1k), because the designs refresh or evict entries before
their hard deadline. The duration-based counter (`boundaries crossed = elapsed/300`) and the
event-based counter (`hard expiries fired`) measure different things, and a report that quotes only
the first would overclaim. **When a phase proves a time span, check whether the event it is named for
actually fired**; a zero is a finding, and the cell that would exercise it (no refresh on a subset,
or early expiry disabled) is the follow-up.

### A comma-separated DSN list is not a DSN

Study 05's runner hands a multi-node cell `postgres://…yb-n1…,postgres://…yb-n2…,postgres://…yb-n3…`.
A single `pgxpool.ParseConfig` rejects it (`failed to configure TLS (sslmode is invalid)`), so **every**
YugabyteDB 3-node cell failed at connect and the cluster arm looked like a broken study. The tempting
one-line "fix" — use the first endpoint — would connect every client to one node's SQL layer and then
report the result as a cluster number, which is exactly what a multi-node study must not do. Study 03
already solved this with one pool per endpoint and a round-robin wrapper; the reusable version now
lives in `platform/adapters/pgxdb` (`OpenSpread`) and every result records `connection_nodes`. **A
multi-endpoint study must prove it spread its connections, or say it did not.**

### A probe's resolution is the smallest finding it can report

Study 04's new staleness probe polls the design's aggregate every 2 ms until it reflects an
acknowledged child change. It measured p50 ~4 ms for both trigger- and application-maintained
rollups — but that figure is mostly the poll interval plus one round trip, not the design's window.
The defensible statement is "no window above roughly 5 ms", and the analysis says exactly that. **A
measurement can only be as precise as the instrument, and an instrument built to answer "is there a
window?" cannot answer "how long is it?" below its own sampling period.** When the expected effect is
near the resolution, say the bound, not the number.

### A root-relative path compiles only under one root

The book's evidence module read `/evidence/v2/claims.json`. A leading slash makes a Typst path
**root-relative**, so the book compiled with `book/build.sh` (which passes `--root book`) and failed
everywhere else: an editor's language server roots the project at the repository, where the same path
becomes `<repo>/evidence/v2/claims.json` and the import dies in milliseconds
(`file not found (searched at …/evidence/v2/claims.json)`). A path **without** a leading slash is
resolved against the file that names it, so it works under either root. **A build script that
happens to set the right root hides a source that only works under that root**; make paths
file-relative unless there is a reason to be root-relative, and test the book from both roots.

## Working the AI work queue and measuring open-loop demand (2026-09-23)

### A stale empty `.git/index.lock` makes the queue's git steps fail, and the message can look like a conflict

`queue integrate` / `queue complete` run git in a shared worktree. Twice in one session a
zero-byte `.git/index.lock` left by a killed git process made the command fail; one of those
failures printed `merging … into main conflicted; resolve it in the main worktree and retry`, while
`git merge` of the same branch completed cleanly. Check `ls .git/index.lock`, confirm no git or
podman process is running (`ps`, `podman ps`), remove the lock, and retry — do not start resolving a
conflict the message invented. **When a tool reports a conflict, verify it with git before acting on
it.**

### A second integrate of an already-merged branch can be a no-op that reports failure

After a task was integrated, a follow-up documentation commit on the same branch was merged into main
by hand because `queue integrate` kept reporting a conflict. The queue state still moved to
`completed` afterwards. Re-integrating an already-integrated branch is not a supported way to land a
correction; make the correction part of the task's own commit range, or treat it as a small new
change with its own worktree and tag.

### A single-trial closed-loop floor is not a capacity reference for an open-loop phase

The same Study 05 cell (`owned-opt-redis-aside-relaxed-coord`, `pg-single`, `db-only`) measured a
warm closed-loop floor of **7,787**, **19,921** and **21,175** ops/s on three consecutive runs of one
session — a 2.7x spread on shared, throttled laptop cores. Any open-loop conclusion that compares a
delivered rate against "the" closed-loop floor from another run is measuring the machine's mood.
Measure the floor in the same cell, and state the spread.

### The open-loop saturation flag is fixed-rule and can fire on a fully delivered rate

`measure.RunOpenLoop` marks `client_saturated` when any arrival is dropped **or** when the p99
scheduling lag exceeds one arrival slot. At 5,000 ops/s offered the cell delivered 4,969 ops/s with
**0 dropped**, yet the flag was set because the p99 lag (1.12 ms) exceeded the 200 µs slot. Read the
counts (`dropped`, `offered` vs `completed`), not the boolean: `client_saturated = true` with
`dropped = 0` and `delivered ≈ offered` is a scheduler-jitter artefact, not a saturated client.

## Retiring evidence gaps and single-sourcing the registry (2026-09-24)

### A superseding claim must retire the gap it closes, in the same version

v3 of the book evidence registry added `v3-06-colocated-vs-noncolocated-locality` (the Study 05
key-locality pair, executed) and `v3-05-yb-cluster3-strict-relaxed` (three-node cells whose client
spreads over all three endpoints) while carrying `v2-gap-05` and `v2-gap-07` forward
**byte-identical** from v2. The two active gaps then said "Study 05's placement pair was never
executed" and "balanced client access across cluster endpoints is not proven", and the book printed
them beside the claims that contradicted them. v3's own `coverage.json` **had** been updated for the
endpoint spread, so two files in one package disagreed and only one of them reached the page.
**When a new run closes part of an existing gap, retire or narrow that gap in the same change**, and
record which clause closed (`closed_dimensions`) and which successor owns the rest. v4 enforces the
mechanical part: a claim in `claims[]` is active or partially superseded, a retired claim may not
also be active, and a gap must declare its kind. The judgement call — *does this new evidence close
that old sentence?* — stays human, and that is precisely where v3 failed.

### A gap must say what sort of gap it is

"A gap callout is missing evidence" collapsed four different findings into one badge:
a **schema limitation** (no representation can answer the question at any speed — the hold funnel),
a **coverage gap** (nobody ran that regime), an **instrument gap** (the runner cannot take the
measurement — tablet/leader distribution) and an **unstable measurement** (measured, not
publishable). They imply different next actions: the first cannot be fixed by running anything, the
second and third can, and the fourth needs more trials. The registry now carries `gap_kind` and the
book badges from it.

### An artefact review must read the manifest before asserting missing provenance

An external review of the 49-page draft raised, as a release blocker, that every provenance field
said `unknown`. The committed PDF renders real values for commit, describe, dirty state, build
clock, image digest, evidence digest and source-tree hash, and `book/dist/build-manifest.json` — in
the same directory — records all of them. The only literal `unknown` was inside Typst's own
`0.15.1 (unknown commit)` version string. **Check the primary artefact and its manifest before
reporting a provenance failure**; a review that asserts from a lossy text extraction repeats the
error it is reviewing for. The residual was real and is now hardened: a compile without the build
script's inputs renders a visible *Unverified development build* banner, and `build.sh` fails if
that banner appears in a release PDF.

### A count gate over the rendered page catches layout defects a compile cannot

`build.sh` counts how many registry claim ids actually appear in the extracted PDF text. Adding the
lifecycle note to each evidence card made one card's meta text long enough to squeeze the id column,
so `v3-gap-01-study05-remaining-regimes` broke across a line and the gate reported 28/29. The
compile was perfectly happy; the page was not. **A count derived from the rendered artefact is a
layout check as well as a content check** — keep ids on their own line, and keep the gate.

### Static source checks turn a "we must remember" rule into a build failure

The stale registry path in prose was a one-line mistake that survived a whole release. It is now a
rule the build runs before compiling (`book/evidence/check.sh`): one active registry path, no
retired claim rendered as evidence, every quoted id resolved, every gap declared. Where a rule
cannot hold yet, the exception goes in `check-waivers.txt` with an owner and a reason, and the
checker prints it — a waiver is visible debt, not a silent bypass. The Typst loader is the second
line of defence and names the successor of a retired claim, because that is the string an author
needs.

### A shared image tag needs its own lock

Every worktree builds the same `localhost/ads-book:1` image, so two concurrent builds race on the
tag even though neither touches a database. `book/build.sh` now takes an `ads-book-build-lock`
podman volume (the same atomic-create pattern as the benchmark lock) and never touches
`ads-run-lock`, which stays reserved for measurement.

### A table column sized `auto` beside prose will eat the page

`book/concepts/cache-consistency.typ` declared five mechanisms in a table with
`columns: (auto, 1fr, auto)`. The long evidence-status cells claimed almost the whole text width, the
`1fr` middle column collapsed to a sliver, and Typst hyphenated the prose into vertical fragments —
`se - ri - alises`, `asyn - chro - nous` — for five pages, in the book's most conceptual section.
Counting orphan-hyphen lines in the extracted page text makes the defect measurable rather than
aesthetic: 99 before, with 17 on the worst page; 30 after, worst page 3. **Never put `auto` on a
column whose cells hold sentences**, and prefer cards to a wide table when each row needs four or five
labelled fields — the labels then sit beside their content instead of in a column header the reader
has to look back for.

### A labelled evidence unit must not split across pages

Visual review of the book found a full claim card broken after its identifier and a coverage callout
broken after its label. The next page contained the statement, limits and source, but a reader scanning
the page foot could not tell what the title qualified. Full evidence cards and labelled callouts now use
non-breakable Typst blocks, so the renderer inserts the page break before the whole unit. Rendered
artefacts need this check: text extraction can prove a card's words exist, but it cannot show that their
label was stranded on the preceding page.

### Heading hierarchy needs a visible pause at every level

The book's level-2 and level-3 headings initially had enough space to avoid collisions but not enough
to separate consecutive ideas after a dense paragraph. Define the hierarchy gaps once in the shared
Typst heading styles, then inspect a dense chapter page as well as a chapter opener. The chapter's
page-break policy and the subsection rhythm are complementary: the first controls entry to a chapter,
and the second helps readers scan its internal structure.

### The queue CLI can leave a stale git lock behind on WSL + Windows podman

Four queue commands in a row failed with "Another git process seems to be running", naming a
`.git/index.lock` that no process held. The queue CLI runs in a podman container and drives git
against the shared working tree; in this WSL + Windows-podman setup its commit step can leave the lock
file behind, and the next command refuses to start. The publish event chain of the first task had to
be committed by hand for that reason. **Check `ps -ef | grep git` and `podman ps` before removing a
lock** — an empty lock with no git process and no running container is stale — then remove it and
retry. Do not conclude that a concurrent agent is holding it without that check, and do not widen the
habit into removing locks that a live process owns.

### A template token inside backticks renders literally, and only a page-text gate catches it

Replacing a hard-coded registry path with an interpolated one looked like the right fix, and it shipped
a draft whose page 2 read **Evidence registry `#registry-label`**, whose Chapter 18 heading read
**Evidence registry (active `#registry-version` claims)**, and which leaked the same tokens in two more
chapters. Six sites, two causes: Typst does not evaluate inside backticks (raw text is not evaluated),
and it does not evaluate a `#name` inside a plain string argument. Both read as correct source to
anyone reviewing the diff, which is why the diff-level review missed them and a reviewer reading the
rendered page did not.

**A rendered-artefact assertion is the only thing that catches this class.** `check.sh` now fails a
backticked `#` token or a build-input token inside a string, but the gate that would have caught every
occurrence is the one that reads the page: `build.sh` greps the extracted text for `#registry-`,
`#build-`, `#source-`, `${`, `{{` and `UNKNOWN_PLACEHOLDER` and aborts the build. Both were verified by
negative control: a backticked token fails the source rule, and a literal `{{` — which the source rule
deliberately does not look for — fails the build with the offending token named.

The lesson is wider than Typst: **when a generator interpolates a value into prose, assert on the
generated output, never on the generator's source.** Every earlier gate in this book checked structure
(claims indexed, fonts embedded, digest present) and all of them passed while the page said
`#registry-label`.

### A registry re-scope does not fix prose that paraphrases the retired wording

v5 removed "at no measurable cost" from `v2-10` and the write-path cost from `v2-16`, and the cards
re-rendered from the registry immediately. The book's **prose** did not: two chapters still repeated the
retired sentences, and a reader would have found the claim corrected beside prose asserting it. Fixing a
claim and fixing what the book says about that claim are two jobs, and the second one has no generator.

**Search the prose for the wording you removed, by phrase, as part of the same task.** The check that
would catch it mechanically is a lint over each supersession ledger's "before" strings against the book
sources; until that exists, the phrase search is the control. This is the same shape as the earlier
lesson that a superseding claim must retire the gap it closes — the artefact was updated and the thing
that quotes it was not.

**A third instance, 2026-09-25, and a harder one.** The T1 replication replaced `v2-13`'s single-run
"6.85x" with a replicated range, and the phrase search found the old figure in **four** places —
`chapters/minor-variants`, `chapters/how-to-choose`, `chapters/scenarios` and
`concepts/derived-state-and-history` — each citing the same claim card that had just been corrected. It
is harder than the earlier two because the change is a *number*, not a phrase: "6.85x" is a string a
search can find, but the sentence around it was also stale ("the application rollup kept writes near the
reference for the same read gain", where the read gain was never re-measured). Fixing a claim means
re-reading every sentence that quotes it, not just replacing the figure. The proposed lint would help;
until then, grep the retired number *and* read its paragraph.

### A checker whose scope depends on how it is invoked is not a checker

`book/evidence/check.sh` excluded hits under its own directory by comparing scanned paths against an
absolute `$HERE`. Run with a relative `--book`, the scanned paths were relative too, the comparison
failed, and the rule reported 26 failures that were all its own documentation. The rule's *scope* had
silently changed with the invocation, which is worse than a wrong rule: two agents running the same
check would have disagreed about the same tree.

Two fixes came out of it, both worth copying: **canonicalise the scan root before comparing paths with
it**, and **verify a checker's exclusion with a negative control** — the first attempt at the fix looked
correct and suppressed the very control it was supposed to catch, because the exclusion pattern was
matching the *content* of the matched line rather than its path.

### A one-sentence summary of a mechanism will overclaim before the mechanism does

The cache chapter's cards began with "Each one closes exactly one race" and "only one of them is a
performance control that proves nothing about freshness". Both were false as written — two mechanisms
are load controls, and a TTL closes no race at all — and neither sentence was in the original table.
They were introduced by the *summarising* pass, not by the content, and a card format makes a
summary sentence look authoritative because it sits in the same block as the guarantees.

The same pass merged two mechanisms into one card and gave them one guarantee: a pessimistic lock and a
conditional update are not interchangeable, because a lock serialises writers while rejection belongs to
the recheck (or to the predicate in the write). A card that carries two mechanisms ends up asserting
what only their combination provides.

**Check a summary sentence against every item it summarises, and give each mechanism its own card.**
When a sentence has a "only", "always", "exactly" or "at most", it is a claim about the set, not about
the row in front of it, and it needs the same evidence test as the numbers.

### A renderer's output attributes are part of your build contract

Every SVG in the book carried `textLength` and `lengthAdjust="spacing"` on each `<text>` element. PlantUML
emits them to place text at a guessed width, Typst honours them faithfully, and every label was therefore
stretched to the renderer's estimate instead of laid out in the real font. The symptom looked like a
per-figure problem — one long-labelled diagram was obviously letter-spaced while the short ones looked
acceptable — and it was not: all five assets had it, and a fix aimed at the visible figure would have left
the other four subtly wrong.

Two things follow. **Fix it where the artefact is generated**, and record the fix as a build input the way
the renderer digest is recorded: `export.sh` now strips the attributes and writes
`postprocess: strip-textlength-v1` into the figure manifest, so changing the step invalidates every asset
exactly as a new renderer does. And **gate the property, not the figure**: `check.sh` fails any asset
carrying either attribute, which is what stops the next re-export from quietly reintroducing it.

### A default copied into a second script is a defect waiting for an excuse

Moving the book to registry v5 meant changing the default in `build.sh` and in `lib/evidence.typ`. A third
copy lived in `book/evidence/check.sh`, and nothing failed: `build.sh` passes the version explicitly, so
every build was correct while a standalone `check.sh` run checked the **frozen v4 package** and reported
OK. Two agents could have run "the same check" against the same tree and validated different registries.

The rule this project keeps relearning: when a value becomes a build input, every script that keeps its own
copy of it must move in the same commit — and a checker run standalone deserves the same acceptance test as
the checker run by the build.

### Typeset the data, do not edit it

The registry is a JSON file, so its numeric ranges use ASCII hyphens; the page typesets a range with an
en dash. Prose and cards therefore disagreed on the same number, and the tempting fixes were both wrong:
regenerating the evidence text with typographic characters would put presentation into the evidence
record, and rewriting each sentence by hand would drift the moment a claim changed.

The fix sits at the boundary: `typeset-ranges` converts a hyphen between digits to an en dash **as the
card renders**, so the registry stays verbatim and the page reads consistently. The same rule applies to
any generated text in this repository — transform at render, keep the source a faithful record, and let a
gate (here `validate`, which resolves the registry text against its report) keep the source honest.

### A provenance field that names the last commit is not stable until the commit exists

`book/assets/export.sh` recorded `source_revision` as `git rev-parse HEAD:<source>`. Run with `--write`
while the source change was still uncommitted, it stored the *previous* blob; the next `--check`, after
the commit existed, computed a different value and failed the build on a figure nobody had touched. The
manifest was not wrong about the file — it was answering a question ("what does HEAD say?") instead of the
one that matters ("does this hash describe this content?").

It now records the committed blob only when that blob matches the file's content, and `uncommitted`
otherwise. **Prefer a field that describes the artefact to a field that describes the repository state**,
and test a manifest by writing it dirty and checking it clean — which is exactly how this surfaced.

### A glossary is the one place where duplication is not a defect

The same twenty-odd terms were defined in passing across nine concept chapters, each definition slightly
different, and two of them were three registry versions out of date. The glossary does not remove those
sentences; it gives every term **one** site a reader can be sent to, which is what makes the in-place
definitions safe to keep short.

Two details are worth copying. The glossary names the registry field where a term is also a field
(`gap_kind`, `trials`, `resource_framing`), so the vocabulary and the schema cannot drift apart silently.
And it was written against the *registry's* vocabulary rather than the prose's — that is how
`retired` and `partially superseded` came to have definitions at all.

### A figure that asserts a number is a claim without a card

The two tranche-1 figures assert none: the decision path orders the book's own recommendations, and the
freshness timeline draws the strict contract at the point where it binds. Both carry a form and
evidence-status label and a text equivalent beside them, which is the rule from the figures task applied
to figures that carry no measurement — the label is what tells a reader whether to look for a claim card.
A conceptual figure may place the numbers a claim provides; it may not introduce one.

### Let the renderer count, not the build

The review asked for a full evidence card where a claim is first established and a short capsule wherever
it is cited again. The obvious implementation is a build step that scans the sources for `registry-card`
calls, decides which is first, and passes the answer in — and it is the wrong one: the scan's notion of
document order has to be kept in step with the include order, and nothing fails when the two drift.

Typst already knows the order, because it renders in it. A `state` counter incremented as each card
renders gives the same answer with no second source of truth, and it cannot disagree with the page. The
scan-free version also needed no build input at all, which is how the vestigial one got deleted rather
than documented.

### A legend that restates its data will eventually contradict it

Page 3 listed the kinds of gap by hand. The registry had four; the page named three of them in its own
words for three versions — the same failure mode as the stale gap claims, one layer up in the prose.

The legend now reads the kinds **in use** from the active claims and looks each one up in a description
table, and a kind with no description panics with the missing name. A future package that introduces a
fifth kind therefore cannot ship a page that lists four: the build stops, and the message says which
definition is missing. Where a page explains data, derive the list from the data and write only the
explanation.

### The provenance record is an artefact too, and it drifts first

The release manifest reported `"edition": "Edition 1 (draft)"` while the title page of the PDF it described
said Edition 2. The field was a hand-written copy of a value that lives in `lib/config.typ`, and nothing
compared the two — the build checked pages, fonts, claim ids, digests and tokens, and every one of those
passed on a manifest whose edition was wrong.

Two habits come out of it. **Derive every field of a provenance record from the same source the artefact
itself uses**, and then **assert the derived value is actually in the output** — `build.sh` now reads the
edition from `config.typ` and fails if that string does not appear in the extracted page text. A manifest
is the one file a reader trusts to describe what they are holding; it cannot be the file nobody checks.

### Staleness must depend on content, not on repository state

`figure-manifest.json` carried `source_revision`, the commit that last changed each figure source. Running
`--write` before the commit recorded `uncommitted`; running `--check` after it recorded a hash, and the
build failed on a figure nobody had touched. Earlier in the same project the same field had failed the
other way round, when it recorded the *previous* commit's blob.

The field is useful and stays, as documentation of when the content arrived. What changed is the
comparison: staleness is now decided by the content hash of the source, the content hash of the asset, the
renderer identity and the postprocess step — every one a property of the artefact. **If a field describes
the repository rather than the artefact, keep it out of the equality test.**

## Revising the book's figures and prose (2026-09-24)

### A sequence figure must give every participant its prose names a lifeline
The lease-versus-fence figure described a second filler but drew only R1, the lease and the cache, so the
second filler had no lifeline and the coordinator's messages to the cache made the coordinator look like
the publisher. Its closing note then attributed the rejection to "the source generation" although no
source check was drawn. A figure that names a participant in prose owes it a lifeline, and a note may
only name the mechanism the arrows actually show.

### A fencing token and source-generation validation are different mechanisms
A sink-side fencing token is accepted or rejected *at the sink* by the highest token the sink has already
seen, so its ordering depends on a newer token having arrived. Source-generation validation consults the
authoritative generation *before publishing* and can refuse an old publish even when no newer operation
ever reached the sink. The two close related races; writing one as the other made the figure wrong while
every word in the caption stayed true.

### A tall sequence diagram overflows a non-breakable figure block into the footer
`figure-evidence` keeps the image, its caption and its text equivalent together, so a 1097-pixel-tall
sequence rendered at 100% text width was taller than the page and printed over the running footer. The
fix is a `image-width` parameter that centres a tall figure at a reduced width; shrinking the source with
renderer scaling does not help, because it shrinks text and height together. Measure the SVG aspect ratio
before placing a sequence.

### Disconnected panels can be reordered by the layout engine
Three panels with no edges between them rendered 3-1-2 on the page. Forcing the order with hidden edges
produced a diagonal that wasted half the width, and splitting the figure in two fixed both order and
legibility. Prefer one connected figure, or separate figures, over disconnected panels whose numbers must
carry the reading order.

### A renderer deprecation warning can become visible page content
`skinparam ParticipantPadding` made PlantUML print "Please use CSS style instead of skinparam
ParticipantPadding" *inside the diagram*. A successful render exit code says nothing about what the
picture contains; inspect the raster, not only whether the renderer succeeded.

### The queue's worktree on WSL plus Windows podman records a relative gitdir
`queue claim` created the expected branch and directory, but `git worktree list` marked it `prunable`
because the recorded gitdir file pointed at a relative path that no longer resolved. Running
`git worktree repair <path>` from the main checkout fixed the link, after which the claim, commits,
submit and integrate all worked. Check `git worktree list` after a claim, before trusting the folder.

## Recording the method (2026-09-25)

### A method field that cannot fail its own rule is documentation, not a gate
Adding `evidence_form` and `attribution` to the registry changed nothing on its own: the v1 registry
was schema-valid and resolved to nothing, which is the failure mode a schema cannot catch. The fields
only became a gate when each rule got a lint *and* a fixture that must fail — a package claim with no
isolation pair, a confounded claim with no withheld attribution, `control_coverage: present` beside
`controls_fired: 0`, `run_replicated` with one run and fixed order, a reversal with three levels. All
29 claims passed the lints unchanged, which is the useful result: the corpus was already honest, and
now it can no longer become dishonest silently.

### A metadata-only registry version needs a byte-identical carry-forward test
The risk of a "metadata pass" is that it quietly rewords a statement while claiming to add fields. The
guard is a test that fails on any difference in statement, limits, confounds, trials or anchor count
between the new package and its predecessor (`TestV6CarriesV5ForwardUnchanged`). It also pins claim
order, so a reordering cannot hide a change.

### A source checker needs its own negative control, or it fires on the wrong rows
The first version of the `concept-checklist` rule extracted every `` | `name` | `` row from
`CHECKLIST.md`, which also matched the four evidence-status legend rows (`direct`, `mechanism`,
`analogy`, `gap`) and reported them as chapters that do not exist. The rule now requires a row with
eight columns, and a temporary-copy test proves it still fires when a real concept row is deleted.
Every new checker rule needs the same treatment as a correctness gate: show it catching the thing it
is for.

### A book source check can already be failing on `main`
`book/evidence/check.sh` failed on `main` before this task: `book/CHANGELOG.md` named a fixed registry
path, which the `registry-literal` rule forbids, and no session had run the check since that line was
written. Run the check first, on an unmodified checkout, before attributing a failure to your own
change — and fix the pre-existing one when it is in the tree you are already editing.

### An ignored staging directory is not a handoff
`.queue-drafts/` is excluded by `.git/info/exclude`, so a new git worktree does not contain it and
nothing written there can reach `main`. Drafts staged there exist only in the folder that wrote them.
A program that must survive a session belongs in a tracked document (here
`docs/empirical-method-program.md`) with the publish commands it needs.

### A missing toolchain is a resource gap, not a decision gap — and it was not missing
This host has no Podman or Typst in the WSL distro, so the first conclusion drawn was that measurement,
the queue CLI and the PDF build were all unavailable. The honest response to that conclusion was the
verifiable part — registry, lints, templates, method, checker rule — plus a named blocker rather than a
half-run matrix or an unverified Typst edit. **That conclusion was wrong.** Windows-side Podman 5.8.1
was installed the whole time at `/mnt/c/Program Files/RedHat/Podman/podman.exe`, reachable through
`infra/lib.sh`, with the pinned Typst image cached; `podman info` reported the documented 8 CPUs /
16.5 GB. Two lessons: check the Windows side of a WSL host before declaring a container resource
absent (the repo's own `hostpath` helper and the `wsl-podman-bridge` work exist for exactly this), and a
"resource gap" claim is itself a claim that needs checking — not merely one that needs stating.

## Replicating a claim, and what it cost (2026-09-25)

### A replicated ratio can keep its direction and lose its magnitude
The registry's `v2-13` cited one run at 6.85x (83.7 vs 573 ops/s) for the trigger-vs-application rollup
pair. Three whole runs with randomized order gave 3.99x, 4.94x and 5.69x — same direction every time,
magnitude never as large. One run established a direction; it could not establish a number. Reporting
the replicated range and retiring the point estimate is the whole point of the exercise.

### A prediction's falsifier must be as strong as its expectation
P1 expected the 6.85x to reproduce, but its stated falsifier was only "the median falls inside the
suppression floor or reverses". The median was 4.94x and nothing reversed, so by its own wording the
prediction was "not falsified" while its actual expectation plainly failed. A falsifier has to name the
observation that would refute the *expectation* (here: a three-run median far from 6.85x), or it
degrades into a test that passes whatever happens.

### A cell whose own trials span the ratio cannot report a point
The trigger cell's three trials inside one run were `[40,159,162]`, `[22,124,165]` and `[61,91,93]` —
spreads of 77%, 115% and 35%, against the corpus's ~20% rule, while the application cell was flat at
1.4–6.8%. The reported ratio is therefore a range over an unstable cell, and the registry claim must
say so instead of quoting a ratio as though both sides were equally precise.

### The first trial of a cell can carry a 2–5x position penalty
In all three runs the trigger's first trial was its slowest (`40` vs a `159` median, `22` vs `124`,
`61` vs `91`). Warmup is discarded by duration, not by effect, so a design whose cost grows with cache
state or bloat starts each cell cold. Two consequences: the median is doing real work here, and the
within-run position effect the corpus already suspected is now measured — inside a single cell, before
any cross-cell ordering enters.

### A worktree belongs inside the project, whatever an example says
`AGENTS.md` gave `git worktree add ../ads-<study>-<topic>` as the example, so two tasks were started in
sibling directories outside the project — and one of them was deleted at completion, taking the
owner's view of those files with it. Every historical task `STATUS.md` records its worktree as
`.worktrees/<name>`, `tools/queue claim` creates `.worktrees/<name>`, and `.worktrees/` is already in
`.git/info/exclude`, so the in-project location was both the established convention and the obvious
one. Two lessons: when an instruction and the artefacts disagree, the artefacts are evidence and the
instruction may be stale — check how the other 30 tasks actually did it; and a path is a decision the
owner has already made, so writing outside the project is a defect rather than a preference. The
example in `AGENTS.md` is corrected.

### A seeded shuffle is not a varied order
The T1 read tranche ran three times with `--order-seed` 101, 202 and 303, and the two-design list came
out `n4_rollup_app, n3_rollup_trigger` **all three times** — so the design that read `r01` 1.52x faster
held the same first position in every run, and design identity is confounded with position. The
three-design write tranche was little better: three seeds produced only two distinct orders (202 and 303
coincided). The shuffle is not broken — testing it directly gives different orders for different seeds —
but with a short list a coincidental repeat is likely, and the honest fix is to read the realized order
out of each manifest rather than trust the method's name. "Randomized order" describes the mechanism;
only `design_order_<topology>` in the manifest describes what happened. Encoded as a limit on `v2-13`
rather than left as an assumption.

### Identical SQL and identical storage do not make a read comparison clean
n3 and n4 have byte-identical `queries.sql`, schemas that differ only in comments, and exactly equal
`storage_bytes` and load timings — and one of five reads still differed by 1.52x in all three runs. The
reusable point is that ruling out the obvious causes (a different query, a bigger table) is not the same
as having a mechanism, and a *consistent* confounded difference is still confounded: three runs agreeing
on a comparison where order never varied is three runs of the same mistake. That tranche also ran without
`explain`, so it could not test the one explanation that mattered. Run the position control first: a
replicated gap with no mechanism is an open question, not a result.

### A single-order comparison can carry a 1.5x positional term
The position control settled the r01 anomaly in two runs and 8 minutes: with `n3` first, `n4/n3` = 0.62;
with `n4` first, 1.50. Each design reads `r01` about **1.55x** faster when it holds first position, and
the plan, SQL, storage and row counts are identical in all four combinations. So a 1.5x "design
difference" was entirely positional — and the same harness measured it the other way round in the write
phase, where the *first* trial was the slowest, so the sign of a position effect is not even predictable
across phases. Two consequences: any read comparison reported from one realized order should say so, and
the cheapest possible experiment (run it the other way round) should be tried before theorising about the
mechanism. The earlier tranche's three-run consistency was worth nothing here, because all three runs
repeated the same order — replication of a comparison is not replication of an *order*.

### Extend the measurement before relaxing the rule
The corpus had a real reversal (optimistic vs pessimistic locking, crossing between 1 and 4 writers) and
could not register it: the new `crossover` lint requires five levels and a stated crossing interval, and
the sweep had four. Relaxing the lint to fit the data was the available shortcut and the wrong one — a
rule that is amended in the same task where it blocks its author is not a rule. Running three more levels
(2, 8, 32) cost under nine minutes, satisfied the lint unamended, and produced a **better** result than
the exemption would have: the crossing is now bracketed by two adjacent levels (0.89 → 1.10) instead of
by a gap from 1 to 4. When a gate blocks a result, the question to ask first is what measurement would
satisfy it; the answer is often cheaper than the argument about the gate.

### When both arms of a control move together, the control measured the run
The retry control was built to prove that the optimistic design's disadvantage at 16 writers comes from
retrying. Setting the retry budget to zero did not narrow the gap at all (`c2/c1` 1.60 against 1.61), and
the giveaway was sitting in the same table: **`c2_pessimistic_lock` never retries**, at any setting of the
flag, and it still moved 10% between the two arms. A design that cannot respond to the option under test
moving by the same amount as the design that can means the difference belongs to the *run*, not to the
option. The general form: before reading an arm-to-arm difference as an effect, check that a unit which is
insensitive to the change did not move too. That check costs nothing when the control already contains
such a unit — and this control did, which is why the false positive lived for one run instead of a week.
The consequence was writing the claim's mechanism sentence and then withdrawing it: the claim keeps its
numbers and loses its explanation, which is the correct trade.

### A falsifier's baseline must be measured the same day
P1 asked whether the lock's throughput stayed within 20% of *870 ops/s*, a figure from two days earlier,
and failed at +34% — while the same-day baseline showed +10%, i.e. inside the band. The prediction was
therefore reported as both falsified and confirmed depending on which baseline the falsifier named, which
is a defect in the prediction and not in the design. A baseline quoted in a falsifier is part of the
experiment; if it was not measured in the same session it carries whatever the day carries. Here that was
a **~30% day term** on a design that does not retry at all (869 ops/s on 2026-09-23, 1,056 on 2026-09-26,
identical settings). Two consequences: state the baseline run id in the falsifier, and re-measure a fixed
reference cell each measurement day so the day term is recorded rather than discovered — the corpus's
absolute ops/s figures are day-specific, and only within-run ratios travel.

### A null from your own query is not an absence in the artifact
Within an hour of writing it, this session's own retry-control analysis declared that the
contention-phase counter invariant "is not recorded" and downgraded its correctness verdict accordingly.
It was wrong. Every result file carries `final_counter_value`, `expected_counter_value` and
`lost_updates`, and the generated report renders them in a contention table — the check that caught the
lost-update control's 3,201 lost updates was sitting in the artifact the whole time. The mistake was
mechanical: I queried `final_value` and `expected_value`, got `null` for keys that do not exist, and read
my own miss as a property of the data. **Look up the field names in a real file before asserting that a
result is missing** — `jq 'keys'` costs one command, and it would have turned a wrong finding into a right
one. This is the second instance of the same family in one week: the earlier one declared that this host
had no Podman and no Typst, when both existed on the Windows side. The reusable rule is that **"X is
absent" is a claim about the world, not about your query**, and it needs the same evidence as any other
claim. The cost here was not just a wrong paragraph: the first version of the analysis threw away a
positive result (0 lost updates under rejected publications) and recorded a limitation that did not exist.

### A ratio quoted without its axis is not a finding
`v2-13` reported the trigger-maintained rollup at 4.0–5.7x the application-maintained one, replicated
across three runs — a solid result, and one that crossed the maintenance decision with nothing. Crossing
it with the size of a replacement gives **2.07x at 10 entries, 4.38x at 60 and 24.20x at 240**: the
replicated figure is one reading on a steep curve, and the same design is either nearly free or
catastrophic depending on a setting the claim did not mention. The lesson is not that the earlier claim
was wrong — its numbers replicated twice — but that **"A is N times B" is an incomplete sentence until
the axis is named**, because the reader will supply a size of their own. The cheapest guard is the one
this experiment used: take the decision the claim is about and cross it with the variable the mechanism
says should matter, then see whether the number survives.

### A phase that does not apply is silently absent
The first version of that experiment crossed maintenance strategy with concurrent writers. It produced no
data, no error and no warning: `harness/main.go` guards the contention phase with
`if !d.hasContention() { return nil }`, and the rollup designs have no arbitration, so the phase returned
immediately and the result files simply had `concurrency: null`. Three runs were spent before the code
was read. The general defect: **a phase that is skipped and a phase that measured nothing look identical
in the artifact**, so a harness should record "not applicable" rather than omitting the key — and an agent
meeting an empty result should read the runner's condition before designing around it. This is the third
instance this week of the same root cause: a result was absent from where I looked and I concluded it was
absent from the world. The first was a worktree path, the second a JSON field name, this one a phase
guard; the check is the same each time — read how the producing code decides, before reporting what the
output means. **The harness now records it** (`main @ 3c16e39`): each requested phase lands in `phases[]`
as `ran` or `not_applicable` with a reason, and the report lists the ones that did not apply — verified
with a negative control in a single run, where the design that does have contention records `ran` and is
deliberately absent from the table.

### Repeating a measurement in sequence measures the session too
The interaction curve was repeated three times per level, and the repeats went in wall-clock order — two
per level, level by level. The trend survived perfectly (the ordering by size held in all three
repetitions) but the magnitudes did not: 15% spread at the smallest level and 55% at the middle one,
against a pre-registered falsifier of 40%. The cause was visible in the absolutes: throughput fell
through the session for *both* designs (at the middle level the application design read 631, then 205,
then 206 ops/s), and the ratio moved because the two designs drifted by different amounts. So the
"spread" I measured is not run-to-run noise at all; it is **drift by time, and the repetition index is
confounded with it** — the same mistake the design-order work caught in the read tranche, one level up.
The fix is the same in shape: vary the thing whose effect you want to exclude. Interleave the repeats,
or run the sizes in a different order each time, or (cheapest) measure a fixed reference cell before and
after the tranche so the drift is a number rather than an inference. This is the third tranche in one
session where a session-level term dominated the effect being measured, which makes it a property of this
setup rather than an unlucky run: **on this host, a single measurement session is not a stable
instrument**, and any repeated claim should say how repetition and time were kept apart — or admit that
they were not.

### Repetition close together in time measures a narrower spread than repetition apart
The lesson above assumed the sequential repeat's spread was an *upper* bound, contaminated by the drift it
had detected. Interleaving the same nine measurements in a Latin square showed the opposite: the spread
roughly doubled (101%, 98%, 104% against 15%, 55%, 52%). The reason is that three runs of a size placed
next to each other in the session sample almost the same machine state, so they agree with each other for
reasons that have nothing to do with the design; spread them across the session and each one samples the
whole drift range. **Adjacent repeats understate uncertainty, and the arrangement of repeats in time is
part of the experiment, not an implementation detail.** A repeat design should be chosen for what it
lets you *see*, and a tranche whose repeats are contiguous should say so rather than present its narrow
spread as the measurement's precision — or measure a reference cell alongside, which is what turned "the
session probably drifted" into "+44% across this tranche".

### A pipe hides the exit status, and that is how a red state got merged
`go test ./... | tail -1` printed `FAIL`, and the command chain continued anyway: the exit status came
from `tail`, which succeeded. The commit and the merge went through with a failing test in the tree, and
it was caught only because the next command happened to re-run the suite. **A pipeline reports the last
command's status, so a check inside a chain must not be piped unless its status is captured** — use
`set -o pipefail`, or run the check unpiped, or read `PIPESTATUS`. The failure it hid was real: a new
profile check requires `replication.runs` to agree with `trials.whole_run_replications`, and a claim that
had just been relabelled to six runs still declared one trial. Both halves of this are worth keeping: the
gate worked, and the shell hid it.

### The machine wanders, it does not drift, and a band is not a bound
Two beliefs from the previous two tranches died in the next one, and both were mine.

**"Drift by time" was the wrong model.** Three readings of one reference cell taken through a single
tranche gave 8.64 → 5.21 → 8.53: down, then back up. Reading only the start and the end would have
reported "no movement" while the middle sat 39% below both. A monotone trend can be characterised by its
endpoints; a wandering one cannot, and the fix is the one this tranche happened to use — take the reading
*inside* the tranche, beside the cells it is meant to help interpret, not only at the edges.

**A band is not a bound.** Tier 10 had seven candidate ranges behind it and the seventh reading landed
below all of them (1.06 against a 1.34–3.15 range built from six runs). Presenting a min–max range as the
interval a future measurement will fall inside is a claim the data does not support; the honest phrasing
is "what has been seen, over N runs", with N beside it, because the more times you look the wider the
range gets. Both lessons point the same way: **this corpus's numbers are samples from a wandering process,
so any interval it reports is an observation record, not a tolerance** — and the ordering between designs
is the part that has so far survived every arrangement.

### The ratio is more portable than the throughput
Over the same three reference readings, the `n4/n3` ratio spread 40% while the application design's
absolute throughput spread 49% and the trigger's 114%. A session term shared by two designs cancels when
they are divided, and only the differential part survives. So the design of a repeated comparison should
prefer ratios and paired readings over absolute numbers wherever the claim allows it: it is the difference
between carrying a number that moves with the machine and one that mostly does not.

### A defect fixed in one copy of a shared script leaves the other copies broken
Study 04's runner learned to pass `winpath`, not `hostpath`, as the image-build context: handed a POSIX
path, the Windows podman reads `/mnt/c/...` as `C:\mnt\c\...` and the build dies with "context must be a
directory". Study 02's runner never got that fix, so **study 02 could not build its benchmark image from a
WSL path at all** — meaning no study 02 claim had ever been replicated, not because replication was
expensive (the run takes five minutes) but because a two-word bug in a copy of a script made it impossible.
The lesson is about where a fix is recorded: a comment inside one study's runner is not a fix for the
repository. **When a runner is corrected, check every other runner for the same defect in the same
change**, and prefer putting the correction in `infra/lib.sh` so there is one copy of it. The corpus's
convention of one runner per study, each with its own copy of the orchestration, makes this failure
repeatable by construction.

### A run that is not tagged is not reproducible, and the runner may skip it silently
Study 02's runner emitted `run/02-ticket-booking/20260926T214452Z` for the first replication and then
produced the second run's results and report without emitting any tag at all. The data was fine; the
provenance was not, and nothing warned. The registry validator caught it — an anchor whose run tag does
not exist is an error — so the gap was found by the gate rather than by a reader. The tag was created by
hand on the commit the run recorded (`repo_commit` in its manifest). **A runner's `--tag` must either
produce a tag or fail loudly; a silent skip is worse than a failure**, because the run looks complete and
is not. Check `git rev-parse` for the expected tag after any `--tag` run, before writing the analysis.

### A grep for a fingerprint can return a different fingerprint
An analysis recorded its run's `inputs_digest` as `a2a6b655e14e7fa1`, taken with
`grep -oE '[0-9a-f]{16}' report.md | head -1`. The report's actual digest line said `0e717c5c3a629cd6`; the
value that was published was simply the **first sixteen hex characters anywhere in the file**, which
happened to be a different field (a config or dataset hash). It survived a validation pass, because the
validator checks that an anchor's digest *appears* in the report and the analysis — and it did appear, in
both, from the same wrong source. **A fingerprint must be read from the line that labels it, never
extracted by shape**, and the cheap guard is to grep the label (`grep "Inputs digest"`) rather than the
pattern. Note also that the wrongness was invisible from inside the file: it was found only by comparing
against a second run, where two identical-looking digests for different data exposed it.

### A run that dies in its load is not a failed design, and the memory is the suspect
Two of three cells in study 03's order-reversed pass died during their `load`, with connections terminated
by the server (SQLSTATE 57P01). The database log showed why: `server process (PID 721) was terminated by
signal 9: Killed`, 122 such lines in one cell. That is the out-of-memory killer, and no property of the
design can produce it. Two consequences worth keeping. First, **"never reached the gate" and "failed the
gate" are different absences** and must be reported differently; conflating them would have turned a host
resource limit into evidence against two designs. Second, **position was the wrong hypothesis and the test
settled it**: `l3` failed as the last design in one run and as the first in the next, while `l1` succeeded
first and failed last — so a hypothesis about ordering was falsified by varying the order, which is the same
move this corpus keeps needing. When a cell dies, rule out the machine before theorising about the design;
record memory beside the run, and bound the database container explicitly.

### A failed check is not an absent one — the third time
Three separate defects in this corpus share one shape: **a command that failed was read as a command that
found nothing.**

- A query for `final_counter_value` used the wrong key names, returned `null`, and was written up as "the
  harness does not record this" — it did (study 04's retry control).
- `grep -oE '[0-9a-f]{16}' report.md | head -1` returned the *first* fingerprint in the file rather than the
  labelled one, and it matched nothing that was checked, so the wrong digest was published (`v2-11`).
- `podman volume inspect ads-run-lock >/dev/null 2>&1 && echo HELD || echo FREE` printed **FREE** while the
  podman machine was stopped: a failed inspection and an absent volume take the same `||` branch, so a dead
  engine looked idle. The runner's own `run_lock_acquire` refused correctly, which is what exposed it.

The fix is the same in all three: **distinguish three states — present, absent, unreadable — and refuse to
act on the third.** For the lock that means capturing the engine's exit status and reporting `UNREADABLE`
rather than `FREE`; for a query it means printing the key names that exist before concluding a field is
missing; for a fingerprint it means `grep "Inputs digest"` and never a pattern. The general rule: a negative
result is only evidence if you have shown the instrument was working when you took it.

## Settling a reversal (2026-09-27)

### A median taken over a noisy cell can manufacture a reversal, and a reversal is a count
Study 02's `v2-08` replication read 0.99x at 100 seats — the ledger *faster* — where the first pass read
1.53x, and the honest reading was that the tier's cell disagreed with itself by 50–92% of its own median.
Six more runs with five trials each (30 trial pairs instead of 6) put the ledger slower in **29 of 30**,
with a pooled per-trial median of 1.391x and every run median inside the original band (1.257–1.712). The
single sub-1 pair (0.940x) is the tail of a distribution, not a second regime. Two reusable points:
**when a reversal is the finding, judge it on the count of pairs that reverse, not on a median of medians**
— the earlier 0.99x was one trial pair in six, and a median of three trials per design let one trial set
it; and **more trials, not more nested medians, is the cheap way out** (the whole settle cost six runs of
~80 s). The registry's own rule for this is the `replication` record plus the spread limit: a cell whose
spread exceeds ~20% supports a direction, not a point ratio, and the claim should say which it is offering.

### An analysis can be indexed nowhere while looking perfectly formed
The report's "Analyses of this run" table matches an analysis by its `run_id` **or** by an `inputs:
run-id@digest, …` frontmatter line (`platform/core/provenance`). A multi-run analysis whose `run_id` is a
descriptive string and whose digests live only in `inputs_digest:` therefore matches **no run**, and every
report it covers says "None yet" — even though the `inputs_digest` is correct and the registry validator,
which only checks that the digest *appears* in report and analysis, passes. The 2026-09-27 `v2-08`
replication analysis hit this, and its two run reports were never regenerated; the settle tranche fixed
its own by adding the `inputs:` line and regenerating with the recorded benchmark image. **A multi-run
analysis needs `inputs: <run-id>@<digest>, …`, and "the report lists it" is a separate check from "the
validator is green"** — run the report generator and read the index before believing the analysis is
attached to its data.

## The tier axis (2026-09-27)

### A range whose endpoints lie outside it is an interior summary, not the cost
`v2-08` quotes the ledger at "1.1–2.3x" of the sell-out race. Sweeping the tier with the same pair in four
runs, both orders, gives a *monotone* cost: 1.075x at 10 seats, 1.237 at 100, 1.409 at 1,000, 1.845 at
10,000 and **2.505 at 100,000**. So the quoted range is just the two interior tiers the sweep's author
happened to bracket, and a reader who says "between 1.1x and 2.3x" is quoting a middle while the ends sit at
"no measurable cost" and "well over 2x". The general rule is to read the ends of a sweep before quoting its
middle: **a band drawn from a sweep is a summary of the levels measured, not a bound on the mechanism**, and
if the endpoints fall outside it the axis has to be named in the same sentence. Here it also produced a
boundary rather than a bigger number, which is the form the corpus prefers.

### A falsifier that does not name its statistic is two falsifiers
Study 02's tier-axis protocol predicted that the 10-seat ratio would fall below the claim's 1.1x floor, and
named no statistic. The run medians came out 1.023–1.180 (pooled 1.075 — confirmed) while the per-trial
median was 1.131 (falsified); the tier is also the only one where the ledger is slower in fewer than 80% of
trial pairs (12 of 20, 60%). The verdict therefore depends on which number the reader picks, which is a
defect in the prediction rather than in the data. **Write the statistic into the falsifier** — "the pooled
run-median ratio", "the per-trial median", "the count of pairs that reverse" — because with a noisy cell
those are three different experiments, and the corpus's own rule (methodology 4) is that a prediction must
name the observation that would refute it precisely enough that one observation can do it.

### A registry's test and its ledger can disagree about the same fact
`TestV6CarriesV5ForwardExceptTheReplicatedClaim` listed **four** re-scoped claims — `v2-13` at the v6 tag
plus `v2-06`, `v2-08` and `v2-11` added after it — while `book/evidence/v6/SUPERSESSION_LEDGER.md` still
said the package carried "one exception". Both files were true when written; the drift came from updating
the machine-checked set and the prose record in separate tasks, three times. The gate stayed green
throughout, so nothing failed: a reader trusting the ledger would have believed the v6 claims were frozen
while the test said otherwise. **A declared-inventory test protects against a silent change only if the
inventory's prose twin is updated in the same change**, which is why the fix here is a post-tag section in
the ledger recording all three re-scopings and the HIGH decision, rather than an edit to the sentence that
was true at the tag. When a test carries a hard-coded list of exceptions, the list and its narrative belong
in one commit. (The post-tag section has since recorded a fourth re-scoping, `v2-09`; the count in this
entry is the one true when it was written.)

## Isolating a confound (2026-09-28)

### A confounded comparison is repaired by varying the confounded factor, not by repeating the comparison
Study 02's `v2-09` claimed the two designs "change places depending on the workload", from a primary pair
that moved buyer count (32 → 128) **and** design order (`p3`-first → `x1`-first) together — two factors, one
comparison, and the claim's own limit admitted neither could be isolated. Eleven runs crossing the two
factors at one tier showed the direction is in fact **stable** (the ledger slower in 55 of 55 trial pairs,
both orders agreeing at every buyer count) and that the apparent swap was the confound. Two of the
pre-registered predictions were falsified in the process. The reusable rules: **when a claim rests on a
comparison that moved two factors, the fix is a crossed design and not more runs of the same pair** — a
confound is invariant under repetition; and **the falsified prediction is a deliverable**, because it is
what replaced a false headline with a true one. A second, quieter confound travelled with it: the primary
ran a different phase set on a five-day-old image, so the re-scope is stated as a direction rather than a
ratio, and the period difference is named as the residual.

### A cell can change state mid-run, and a median does not protect against it
In the 8-buyer, `p3`-first cell the ledger arm fell from ~4,800 to ~1,000 ops/s across its last three
trials — a 5x collapse inside one cell — while the same level's other order sat flat at ~5,000. The cell
median (1,000 for the collapsed arm) is then a description of the collapsed state, not of the cell, and the
ratio built on it (5.10) is an average of two regimes. The corpus's within-cell position and warm-up
findings are the opposite shape (a slow first trial); this is a *late* break, which no median-of-five
protects against and which a spread check only reports after the fact. Two guards apply: read the per-trial
series, not only the median, before quoting a cell, and prefer a cell whose trials agree — the tier-axis
tranche's cells did, and their ratios survived a 2.5x shift in the absolutes.

### A slow trial is not a slow design, and its position decides its name
The 8-buyer cell that fell 5x mid-run looked like a ledger-arm state change, and the corpus kept it flagged
rather than explained. A recurrence test settled it the other way: six runs, twelve cells, sixty trials
produced **four** low readings (619–655 ops/s against a session normal of ~2,500), and three of the four
were in the **CAS** arm; the pre-registered "it is the ledger's" and "it favours one order" predictions were
both falsified. Three of the four events were a single trial that recovered on the next one; the original
looked persistent only because its slow trials were the cell's **last three**, so they set the cell median
and the two normal trials read as the anomaly. Two rules follow. **A name for an event in a short cell is a
description of where it fell, not a mechanism** — "collapse" and "spike" are the same measurement at
different positions in a five-trial median, and a ~7% event rate makes position decisive. And **before
attributing a slow cell to a design, look for the same event in the other arm**: if it is there too, the
design is not the explanation, and the honest next step is to **read the instrumentation that already
exists** — see the next entry, which corrects this lesson's own first draft.

### Read the instrumentation you already have before building more
This entry's predecessor ended by prescribing client-side instrumentation. The client's cgroup counters were
already in every result file — `client_cpu`, written by `platform/adapters/cgroup` since the study's first
commit, one entry per **trial** for the race phase — and no analysis had ever opened them. Reading them
settled the question the previous tranche had declared unassignable (zero client throttling; the slow trials
consume *more* client CPU, so the client waits rather than starves) and then produced a larger finding
nobody was looking for: **35 of study 02's 775 recorded client-CPU cells were throttled**, all in read
phases, and throttled by a **design-dependent** amount inside a single run — 2.3 s for one design against
58.1 s for another out of ~1085 CFS periods, so those runs' read *rankings* can be an ordering of who lost
the client's quota. Two rules. **Before proposing new instrumentation, grep the artifact for the measurement
you want**; the platform is six days older than the analysis habit of reading it. And **when a hazard table
is generated into every report — the CPU-quota section has been there since the same commit — reading it is
part of analysing the run**, not an optional extra; the numbers were one `jq` away for thirteen days while a
mechanism was recorded as "unresolved".

### A flat cell is a ceiling *or* an O(1) design — and a throttled client may be neither the cause nor the fix
This entry was first written as "a flat cell over a wide input range is a ceiling, not a cost", on the
strength of `r1_inventory_row`'s `q01_event_availability` sitting at 37.7–38.7k ops/s across five tiers while
the baseline fell from 25,967 to 595, *and* the same design losing half its client CFS periods to throttling.
Varying the client's budget the next day falsified the causal half: at three client cores the client is
throttled **0 times in 633 periods** and the cell is **still flat** (43.7k → 41.6k, ratio 1.05) — and every
tier is *slower* than at two cores, because the third core came out of the database's share of the machine's
four CPUs (its throttled periods rose 238 → 847 → 1052). Three rules. **Flatness is a symptom with at least
two causes** — a ceiling, or a query whose cost genuinely does not grow (a precreated inventory row is the
design's stated point) — so it needs a plan check, not an inference. That check cost one `grep` in the run
directory and settled it: `r1`'s `q01` is `Index Scan using event_inventory_pkey … rows=1`, cost 2.49,
0.012 ms, three buffers, against a table that holds one row per event at every tier — O(1) by construction,
with the *level* being the machine's point-lookup rate. **Correlation with a quota is not
causation by that quota**: the client was throttled *and* flat, and removing the throttle left it flat;
before blaming a resource, vary that resource. And **when the budgets already sum to the whole machine, the
fix is a smaller footprint or a bigger machine, never a re-allocation** — giving one container a core is
taking a core from another, which the same measurement shows in the other container's counters. What survived
all three tranches is the part with a mechanism: the *ordering* (`q05` fast where there are no per-seat rows,
`r03` slow where a per-customer sort is unavoidable), which is the corpus's ratio-over-absolute lesson
arriving through a quota floor rather than through session drift.

## Sizing the generator (2026-09-30)

### The worker count is a measurement parameter, and Little's law says whether you sized it right
Study 02's read phase ran **8 workers** against a **two-core** client on a four-CPU machine: the client used
111.6 s of CPU inside a 63.3 s window and lost **422 of 633 periods** to throttling, so the read tables were
describing the client. Dropping to 2 workers changed every number in the expected direction — client
throttling 422/633 → 33/632 → **0/631** at 8/4/2, the database's own throttling **900 → 560 → 109** periods —
and it also revealed the diagnostic worth keeping: **operations in flight**, `ops_per_sec × p50`, came out at
**1.9 for 2 workers** in every cell of both designs, i.e. Little's law holding, which is what "this cell
measures the store's latency" looks like. At 8 workers the same product was 3.2–5.3 *against 8 workers*: a
generator that cannot keep its own connections busy is measuring itself. Three rules. **Check the generator's
own saturation before quoting its levels** (the counters and the hazard table are in every result). **Use
`ops/s × latency ≈ workers` as the sizing check**, because it is a property of the *cell*, not of the host,
so it travels. And **when you resize the generator, keep the other phases untouched and label the change in
every manifest** — the read footprint is now 2 with `client_read_conns` recorded, while write, race and hold
stay at the published 8, so a reader can tell which footing a number stands on. The corollary for this corpus:
a level is only quotable within one footprint, and it is the **ratios** (design orderings, and store-bound
cells like `p3`'s 621–631 ops/s at 100,000 seats across all three footprints) that travel between them.

### A resized generator changes the effect size, not just the level
Re-measuring the four read designs at footprint 2 — client unthrottled, every cell pinned at 1.9 operations in
flight — reproduced every ordering and then moved one of them by nearly half: `r1_inventory_row` over `p3`
on `q01` was **2.07x at footprint 8 and 1.16x / 1.48x at footprint 2**. The reason is worth keeping: at 8
workers the client's throttle bit the two designs *unequally* (`r1` kept 5.3 requests in flight, `p3` only
3.2), so part of the old margin was the generator's parallelism rather than the query's cost — a comparison
that was contaminated could still be directionally right and 80% too large. Two rules. **When a generator is
resized, re-derive the margins, not just the levels**: a ratio is portable across footprints only if its
mechanism is in the thing being compared, and here the two that were (`c1` ≫ `p3` on `q05`, the ledger's
`r05` at ~32k) held to within ~30%, while the one that was partly the client's did not. And **the signature of
a genuinely store-bound cell is footprint-invariance**: `p3` at 100,000 seats read 623 at footprint 8 and
680 / 629 at footprint 2 across two orders, which is what a cell whose cost is in the store looks like — worth
checking before quoting any ratio as a design property.

### In-flight counts tell you *what* bounds a cell, not just whether the generator does
**Corrected 2026-09-30** — this entry originally claimed that the two queries returning many rows had
client-bound levels. Both halves were wrong and the correction is the lesson below. The part that survives is
the *shape* of the reasoning: a cell can in principle be bounded by the store or by the harness, and the way to
tell them apart is to measure the loop's store share — at one worker, `ops/mean latency` is that share
directly, and it came out at 1.000 for all 72 cells of four designs (see "Little's law wants the mean"). What
does **not** survive is diagnosing it from `ops_per_sec × p50`, or from a plausible story about how many rows
the query returns (it was one row, and I never checked).

### Little's law wants the mean — a median-based in-flight count fakes client overhead
`ops_per_sec × p50` is only an in-flight count when the latency distribution is tight. Where it is not, the
product measures the *median* concurrency, and where the two-worker read runs showed `q05_band_tickets_sold` at
0.76 against 2 workers, the same JSON recorded `p50 = 0.275 ms` beside `mean = 0.744 ms` — a 2.7x skew, and
0.76 × 2.7 = 2.0. Recomputing with the mean gave **1.000 in all 72 cells at one worker and exactly 2.00 at
two**: every cell was perfectly latency-bound all along, and the "client overhead" I inferred was the skew.
Three rules. **Use the mean for Little's law, never the median** — the median is the robust statistic and the
wrong one here, and the mean is already recorded beside it in every result. **Pick an instrument whose scale
answers the question**: one worker bounds the concurrency above by 1, so agreement with 1.000 cannot be an
artefact of the worker count. And **a mechanism that fits the numbers is not evidence for it** — "the client
decodes 28 rows" explained the gap neatly, the query returned one row, and the check would have taken one
`awk` over the SQL. The corpus's rule about naming the statistic bites hardest when the wrong statistic
produces a *plausible finding* rather than an absurd one.

### A design's position inside a run is part of its measurement
The all-fifteen read run put the same design at the two ends of the order: `c1_count_naive` read **19,068
ops/s as the first design of one order and 13,653 as the last of the other** — a 40% swing, larger than any
design difference among the top five, and enough to invert a ranking. The same pair showed the top five
designs within **4%** of one another while the session's own order term was 4–13%, which means **a rank inside
such a group is not a measurement at all**. Three rules. **Report groups, not ranks**, when the spread inside
a group is smaller than the run-to-run term — the quotable fact here is the 1.5x gap between groups, not
`h1` beating `r1` by 0.2%. **Reverse the order and compare ends**: the pair of orders is what exposed the
position term; a single order would have shown a clean, false ranking. And **a design that moves position
between the runs cannot have a level quoted from them** — with two orders there is no way to separate position
from design for `c1`, and settling it would need three positions, not two.
