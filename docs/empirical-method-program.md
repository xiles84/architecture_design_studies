# Empirical-architecture method program (2026-09-25)

**Owner charter, 2026-09-25.** The book's identity is empirical architecture: explain a mechanism,
predict what it will do, isolate it, sweep the regimes where it matters, combine it into a realistic
package, measure, explain the result, and state where the evidence stops. Prefer experiments that find
**crossovers and boundaries** over experiments that produce a larger throughput number. Controlled
single-variable experiments, interaction experiments and realistic multi-decision packages are three
different forms of evidence, and a package result may never be attributed to one component without an
isolation experiment. Performance stays secondary to correctness.

**Method layer already implemented** (task `repo/empirical-method`, tag
`repo/empirical-method-v1`): methodology sections 16–21, the study-protocol template
(`docs/templates/STUDY_PROTOCOL.md`), the evidence registry v6 with the four method lints and their
failing fixtures (`book/evidence/v6/`), and the concept checklist (`book/concepts/CHECKLIST.md`).

**Nothing here is executed.** This file orders the work and carries publish-ready task specs. Publish
them with `tools/queue publish --spec <task.json> --brief <BRIEF.md>` from a machine with Podman; the
local `.queue-drafts/` directory is git-ignored and is **not** a durable handoff.

## Ordering rule

**T1 and T2 come before T5.** No new domain (Study 06, 08, 09) or topology work (Study 07) is published
while a headline single-run claim from T1 is outstanding, unless the domain task is the only way to
close a hole the charter names. T3 and T4 may run in parallel with T1/T2 where they touch different
studies, because measurement is serial anyway.

## T1 — run-level replication and order (first)

Re-measure the single-run headline claims as ≥3 whole runs each, with an independently randomized
design order per run and ≥3 trials per cell, reporting the median and the **across-run** spread.

- Targets: 01 `v2-01`…`v2-05` (also fixing their missing manifest/tag), 02 `v2-06`/`v2-08`/`v2-09`,
  03 `v2-11`, 04 `v2-12`/`v2-13`, 05 `v2-14`/`v2-15`/`v2-16`, `v3-04`/`v3-05`/`v3-06`.
- Extend the seeded design shuffle already present in the 03/04/05 runners to 01 and 02.
- Include the within-run **position control**: the same design first and last in one phase.
- Acceptance: a claim moves to strength `run_replicated` only if its controls fire in **every** run and
  its `replication` record validates under the v6 lints.

## T2 — crossover and boundary sweeps (the preferred result)

Each sweep is ≥5 levels with repeats, and its deliverable is a monotonicity/reversal verdict plus a
crossing interval with spread. Two results already measured should be **registered** first, because
they exist and only lack a claim:

- Study 04 writer sweep: `c2/c1 = 0.89 / 1.26 / 1.58 / 1.60` at 1/4/16/64 writers — a sign reversal
  between 1 and 4 writers. Register as the corpus's first `crossover` claim.
- Study 04 `d2` metadata-amplification prediction that did not appear — register as the first
  `prediction` with verdict `falsified`.

Then extend: writer count for 04 `c1`/`c2` with the retries-disabled control; cardinality for 04
row/rollup designs below 1 and above 500 entries; hot-parent/skew for 01 `D3`/`D4`/`D6`;
donations-per-donor for 01 `D6`/`D9`; cache capacity as a fraction of working set for 05; more
open-loop rate points to separate the client knee from the server knee.

## T3 — interaction experiments (the four currently unmeasured)

Register no interaction claim until the crossed factors are actually crossed:

- append-only history × concurrency: a new controlled pair, ledger-under-CAS vs ledger-under-lock,
  which also explains why 02 `P3→X1` is only a package result today.
- rollup × concurrency: extend Study 01's Experiment C to the trigger design at low contention.
- cache × database representation × instance count × hit rate: cross the factors in one run.
- sharding × access pattern: point vs range vs scan for 01 `D7`, 03 `L3`, 05 `y1`/`y2` in one run.

## T4 — package attribution isolation

For every package result the book uses, either build the intermediate design that isolates the
mechanism or record `attribution_withheld`. Concrete: decompose 02 `X1` into ledger-storage and
ledger-write arms; add a matched row design to 03 `L2`/`L3`; build the tablet/leader-distribution
instrument **before** `y1`/`y2` is called a placement result.

## T5 — gated domain and topology expansion

- **Study 09 hierarchy** — feasible on this machine with current images; adopt the prediction and
  evidence-form template, and treat its move-cost-vs-subtree-size axis as a crossover source.
- **Study 06 native models** — needs MongoDB/ScyllaDB/Valkey images; closes `v2-gap-06`.
- **Study 08 read copies** — needs ClickHouse/search images.
- **Study 07 topology** — needs ≥2 hosts; it is the only path to hardware-independent replication and
  to the replication/network/failure cells the charter's boundary rule demands.

## Read-path controls (applies across T1–T4)

Fifteen active claims record `control_coverage: absent`, seven of them read-only phases where the
study's controls could not fire. Every read-only phase added or re-run from now on carries a read-path
control — a deliberately stale or wrong copy that is seen to return the wrong answer — or its claim
records `control_coverage: absent` and draws no correctness inference (methodology 20).

## Blocked in the current environment

No Podman and no Typst on this host, so measurement, the queue CLI and the PDF build cannot run here.
The deferred items in this program are blocking on **resources**, not on decisions:
`.queue-drafts/` cannot be committed, `book/build.sh` cannot be run, and no cell can be measured.
