# platform — the shared benchmark core

Everything a study needs that is **not** about that study's domain: how to measure, how to
talk to a database, how to tie numbers to the data, code and analysts that produced them.
Studies import it as the Go module `adsplatform` through a `replace` directive, so a study is
always built against the platform code of the same commit.

It is laid out as a hexagon (ports and adapters):

```
                  ┌───────────────────────── adapters (driving) ─────────────────────────┐
                  │  studies/NN/harness/main.go — the CLI entry point of each study       │
                  └──────────────────────────────────┬───────────────────────────────────┘
                                                     │ calls
          ┌──────────────────────────────────────────▼──────────────────────────────────────────┐
          │ core/                   pure, study-independent logic                                │
          │   measure/     latency summaries, trial medians and spreads, the closed-loop driver  │
          │   catalog/     the annotated-SQL catalogue format, DDL scripts, parameter binding    │
          │   provenance/  inputs digest, repository version, the signed-analyses index          │
          │   inspect/     EXPLAIN capture, engine version and effective isolation, load retry   │
          └──────────────────────────────────────────┬──────────────────────────────────────────┘
                                                     │ depends only on
                                   ┌─────────────────▼─────────────────┐
                                   │ ports/   DB, Tx, Row, Rows,        │
                                   │          Isolation, ErrorClass     │
                                   └─────────────────▲─────────────────┘
                                                     │ implemented by
                  ┌──────────────────────────────────┴───────────────────────────────────┐
                  │ adapters/ (driven)                                                    │
                  │   pgxdb/      PostgreSQL and YugabyteDB via pgx — the only driver import │
                  │   filestore/  JSON results and readable text plans                    │
                  │   markdown/   report formatting rules                                  │
                  │   cgroup/     the container's own CPU-throttling counters              │
                  └───────────────────────────────────────────────────────────────────────┘
```

**The rule:** `core/` and a study's use cases import `ports/`, never a driver. Only a study's
`main.go` wires an adapter in. That keeps the measurement logic independent of the engine
under test, and it means a new engine (or a non-SQL store) is a new adapter rather than a
rewrite of every study.

## What lives here, and what does not

| Here | In the study |
|---|---|
| How to run N workers for a duration or a fixed count, with warmup and trials | Which operations to run, with which keys |
| Percentiles, and when a percentile has too few samples to report | What the operations mean for the design |
| Classifying an engine error as retryable, a unique violation, an ambiguous commit | What a retry means for the design (a lost race, a serialization failure) |
| Parsing `-- name:` / `-- params:` SQL catalogues | The SQL itself |
| Digest, repository version, analyses index | The report's tables |

Each rule in `core/measure` encodes a lesson that once produced a wrong number — finite pools
measured by count, the median rather than the mean, no p99.9 without 10 000 samples. See
[`LESSONS_LEARNED.md`](../LESSONS_LEARNED.md).

## Study 01

Study 01 predates this module and keeps its own harness. Moving it onto the platform would
change the binary behind results that have already been published and analysed; it should
happen only together with a re-run that shows the numbers are unchanged. Its tag
`study-01/v1` marks the code that produced them.
