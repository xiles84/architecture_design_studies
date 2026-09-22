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
