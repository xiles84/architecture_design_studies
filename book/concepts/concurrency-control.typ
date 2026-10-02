#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Concurrency control")

#marker("concept", "concurrency control")
Two strategies at minimum: check a version and retry, or take a lock before updating. Both protect
the same invariant; they differ in where contention is paid and what the failure looks like.

#heading(level: 2, "Optimistic (version-checked)")
The writer reads, checks that the version is unchanged, writes conditionally, and retries on a
mismatch. Under low contention it is close to free; under high contention retries become the
workload. Under one hot key and 16 writers, lock-before-update measured about 1.6–1.8x the optimistic
design across four fresh runs in both design orders, and the retry budget does not explain the gap:
removing it moves the ratio by 1.1% to 3.2%.

#heading(level: 2, "Pessimistic (lock-before-update)")
The writer waits for the row, then updates. With pessimistic locking, contention is primarily paid as
waiting/blocking rather than optimistic conflict retries; depending on the workload that may increase
*or* decrease observed latency. Application-level retries can still be required — for deadlocks, lock
timeouts, serialization failures or other transient failures — so the retry loop does not universally
disappear. On a single engine primitive, `SKIP LOCKED` and compare-and-set both avoid the counter-row
bottleneck measured in ticketing.

#figure-evidence(
  "../assets/fig-optimistic-vs-pessimistic.svg",
  "sequence",
  "conceptual illustration",
  "The invariant recheck, not the lock, rejects an operation that has become invalid.",
  [Two parallel sequences over one row at version 5. *Optimistic*: both writers read version 5; R1's
   `UPDATE … WHERE version = 5` updates one row and moves the version to 6, so R2's conditional update
   touches zero rows and R2 re-reads and retries. *Pessimistic*: R1 takes `SELECT … FOR UPDATE`, R2 waits,
   R1 rechecks the invariant, updates and commits, then R2 acquires the lock and rechecks the invariant
   before deciding. In both panels the lock or the version check serialises the writers, and the
   *invariant recheck* is what rejects the operation that has become invalid. The figure draws `v2-12`'s
   mechanism and asserts no rate.],
  image-width: 72%,
)

#heading(level: 2, "The negative control matters here")
An unchecked read-modify-write lost 93.7% of acknowledged updates in the same run that measured the
lock advantage. A concurrency result without its lost-update control is not evidence.

#registry-card("v2-12-lock-vs-version-upper-bound")

#heading(level: 2, "Where the two designs cross")

The lock is not simply faster. Swept across concurrent writers on one hot key, it gives up about 5–10%
of the version check's throughput at one writer — 0.905x, 0.939x and 0.948x in three whole runs — and
leads from two writers up: 1.035x, 1.106x and 1.166x at 2, 1.56x–1.69x at 16, and 1.73x–1.96x at 64,
reaching 1.68x at 32 in the first sweep. The movement is on the optimistic side: its throughput falls
about 2.3x from one writer to 64 while its conflict count rises from 0 to about 25,000. That
correlation is not an explanation, and the obvious one was tested and rejected — with the retry budget
set to zero the ratio at 16 writers does not narrow (1.60 against 1.61 in the first control, 1.1%–3.2%
in the replications), and the design that never retries moves by the same amount as the one that does.
One run is an exception worth naming: a single 1-writer run measured 3.85x, not because its order
changed but because the version check itself read 306 ops/s there against 1,286 in another run at the
same level and order, with no throttling recorded. So "which is faster" has no answer without the
concurrency level, and the interval where the ordering reverses — between one and two writers — is what
a recommendation should name.
#registry-card("v2-17-optimistic-pessimistic-crossover")

#heading(level: 2, "Boundaries")
One key, one engine, small scale. The 16-writer ratio is a band of five whole runs on two days; the
32-writer point rests on the primary alone, and one 1-writer run is an unexplained exception.

#heading(level: 2, "The checklist")

The eight questions, answered for this concept; the status in brackets is the evidence behind the
answer.

- *Invariant.* The business rule holds under concurrent writers: no lost increment, no oversold seat,
  no double sale. [direct]
- *Where the information lives.* In the current row the writers contend for, plus whatever version or
  lock column guards it. [direct]
- *Where the contention is.* One hot key or row; the corpus measures one event or one installed
  product. [direct]
- *What work moved from reads to writes.* Conflict resolution itself: a version check or a lock turns
  a read-then-decide into one guarded write. [mechanism]
- *What state is derived.* Nothing in this concept; a counter or rollup a concurrent writer maintains
  is derived state, which is its own concept. [mechanism]
- *What can become stale.* The value a writer read before losing the race — which is why the guard has
  to be part of the write. [mechanism]
- *Where it stops working.* The ordering is workload-dependent: at one writer the version check leads,
  and from four writers the lock leads by 1.26–1.60x. [mechanism]
- *The experiment that would settle it.* A writer-count sweep of the same pair with retries disabled,
  so the lock's advantage is separated from per-retry work. [direct]
