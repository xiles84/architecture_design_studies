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
The writer waits for the row, then updates. Latency per operation is higher and the waiting is
visible, but the retry loop disappears. On a single engine primitive, `SKIP LOCKED` and
compare-and-set both avoid the counter-row bottleneck measured in ticketing.

#heading(level: 2, "The negative control matters here")
An unchecked read-modify-write lost 93.7% of acknowledged updates in the same run that measured the
lock advantage. A concurrency result without its lost-update control is not evidence.

#registry-card("v2-12-lock-vs-version-upper-bound")

#heading(level: 2, "Boundaries")
One key, one engine, one trial, small scale. The ratio is an upper bound and a point, not a curve.
