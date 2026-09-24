#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Concurrency control")

#marker("concept", "concurrency control")
Two strategies at minimum: check a version and retry, or take a lock before updating. Both protect
the same invariant; they differ in where contention is paid and what the failure looks like.

#heading(level: 2, "Optimistic (version-checked)")
The writer reads, checks that the version is unchanged, writes conditionally, and retries on a
mismatch. Under low contention it is close to free; under high contention retries become the
workload. The measured upper bound put lock-before-update at 1.76x the optimistic design under one
hot key and 16 writers, with the caveat that the run could not separate the lock's advantage from
extra per-retry harness work.

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

#heading(level: 2, "Boundaries")
One key, one engine, one trial, small scale. The ratio is an upper bound and a point, not a curve.
