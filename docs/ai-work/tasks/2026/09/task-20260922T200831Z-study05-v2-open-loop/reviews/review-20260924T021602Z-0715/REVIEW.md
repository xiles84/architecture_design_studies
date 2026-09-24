# Review — Study 05 v2 open-loop demand phase

**Task:** `task-20260922T200831Z-study05-v2-open-loop`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Verdict:** **approved for integration** — no rework required
**Required tag:** `study-05/v2-open-loop`

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | use the platform `measure.RunOpenLoop`, not a local scheduler | accepted — one validated driver |
| 2 | choose rates that bracket the same cell's closed-loop floor; supersede the other two runs | accepted — 5k below / 30k above the 21.2k floor; superseded reports moved, results retained |
| 3 | reset the wrong-read log per rate | accepted — the second rate no longer inherits the first's reads |
| 4 | state that the server's knee was not isolated rather than force a number | accepted — the driver flags the client saturated, and the analysis says so |
| 5 | label the overload read-only | accepted — no writers, so 0 wrong reads is not a freshness result |
| 6 | keep run 3 canonical and record the 2.7x floor spread | accepted |

## Independent checks re-run by the reviewer

```text
report open-loop table vs the JSON                     values match (30000 -> 15599, 141710 dropped)
inputs digest in report and analysis                   50d5a426c89e96cf (both)
correctness gate                                       passed 14/14, 0 wrong reads in both regimes
closed-loop floor in the same cell                     21,175 ops/s (warm_read_only)
closed-loop not reused as open-loop                    report + analysis say so explicitly
benchmark lock / containers after the runs             lock free, no db/cache containers
reports for superseded runs                            moved to reports/outdated/, results retained
Go build / vet / harness tests                         pass
```

Acceptance criteria: two fixed offered rates with no closed-loop backpressure — met; achieved
throughput, non-attempted operations and latency percentiles per rate — met; offered rate stated and
server saturation stated if reached — met (not reached/isolated, and that is stated); closed-loop not
reused — met; wrong-read rate with the throughput from the same cell and phase — met; generated
report plus signed analysis — met.

## Findings (non-blocking)

1. One trial per rate; the closed-loop floor moved 7,787 → 19,921 → 21,175 ops/s across the three runs,
   so the analysis refuses small differences. Correct.
2. The client saturated in both regimes; the server knee needs a deeper buffer / more workers, which
   would trade against the phase's purpose. Recorded as a limitation, not a defect.
3. `v3-gap-01` in the book registry still lists open-loop demand as unmeasured; this task measures it,
   but `book/evidence` is not owned here. A later evidence task should retire that clause. Noted for
   the next book/evidence task.

## Verdict

Approved for integration. The phase measures what it claims, the numbers carry their correctness and
their generator-side caveats, and the server question is answered honestly instead of being forced.
