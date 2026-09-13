# Architecture Design Studies

**Measuring what architecture decisions actually cost, instead of arguing about it.**

Every claim in this repository is backed by a run you can reproduce on your own machine
with nothing installed but Podman. Every result names the hardware it came from. Every
design is checked for *correctness* before it is allowed to report a speed.

---

## Why this exists

"Denormalise for read performance." "Foreign keys are expensive." "Embedding children is
faster than joining." "Sharding scales your reads."

All of these are sometimes true. The interesting question is never *whether* a technique
helps — it is **how much, at what cost, and on which query**. That question has a number
attached, and the number is measurable.

So each study takes one modelling decision, builds two schemas that differ **only** in
that decision, loads identical data into both, and measures them side by side — reads,
writes, storage, load time, and the query plans underneath.

## What is here

| | |
|---|---|
| [`studies/01-charity-tree/`](studies/01-charity-tree/) | **Tree structures.** A `charity → person → donation` hierarchy across 8 schema designs and 3 database topologies. |
| [`studies/02-ticket-booking/`](studies/02-ticket-booking/) | **Avoiding overbooking.** `band → event → ticket` with events of 10 to 100 000 seats: 14 strategies (pre-created vs created tickets, locks, counters, constraints, extra inventory and reservation tables, two deliberately wrong controls) under sell-out races. |
| [`platform/`](platform/) | The shared Go benchmark core — measurement, SQL catalogue, provenance — laid out as ports and adapters. |
| [`docs/methodology.md`](docs/methodology.md) | The rules every study follows, and why. |
| [`docs/replication.md`](docs/replication.md) | How to reproduce any run. |
| [`docs/environments/`](docs/environments/) | One page per machine results were produced on, including how that machine can mislead a benchmark. |
| [`infra/`](infra/) | Podman orchestration: PostgreSQL 1-node, YugabyteDB 1-node, YugabyteDB 3-node RF=3. |
| [`CONTEXT.md`](CONTEXT.md) | Current project state — what is done, what is running, what is open. |
| [`LESSONS_LEARNED.md`](LESSONS_LEARNED.md) | Everything that cost time or nearly produced a wrong number. |

## Study 01 at a glance

A three-level tree, which is the most common shape in business software:

```
charity ──< person ──< donation
```

Every question about it is either **navigate down** (a donor's recent gifts — easy) or
**aggregate up** (who gives the most, what is the total, when was the last one — where
the cost is, because the answer lives in the leaves and the question is asked at the
root). The eight designs are eight different answers to *when you pay for that*:

![designs](studies/01-charity-tree/diagrams/rendered/00_overview.svg)

Each design differs from its neighbour by exactly one decision, so any measured
difference has exactly one explanation:

| Pair | Isolates |
|---|---|
| D1 → D2 | secondary indexes, on **byte-identical SQL** |
| D2 → D3 | denormalising the grandparent key down a level |
| D3 → D8 | enforcing referential integrity (**foreign keys on vs off**) |
| D3 → D4 | consolidating aggregates upward at write time |
| D4 → D5 | trigger-maintained vs application-maintained rollups |
| D3 → D6 | embedding the child table inside the parent row |
| D3 → D7 | physical data placement only (YugabyteDB) |
| 1 node → 3 nodes | adding two more machines |

→ [Full study](studies/01-charity-tree/)

## Study 02 at a glance

A show company must never sell more tickets than an event has seats:

```
band ──< event (10 … 100 000 seats) ──< ticket
```

The invariant spans rows, so each design has to decide where the capacity check lives and
what serialises two buyers reaching for the last seat — a locked seat, a locked parent, a
counter, a unique index, the isolation level, a pool of seat tokens, or a hold with an
expiry. The headline experiment is a **sell-out race**: a crowd of buyers arrives at an
unsold event at once. Every phase is audited for overbooking **and** under-booking, and two
deliberately wrong designs must be seen to fail, proving the audit works.

![designs](studies/02-ticket-booking/diagrams/rendered/00_overview.svg)

→ [Full study](studies/02-ticket-booking/)

## Running it yourself

```bash
cd studies/01-charity-tree && ./run-study.sh --scale small --duration 10s
```

Podman is the only prerequisite. The benchmark client, both database engines and even the
diagram renderer run in containers — nothing is installed on your machine. See
[docs/replication.md](docs/replication.md).

## The rules

These are what separate a study from a benchmark screenshot. In full in
[docs/methodology.md](docs/methodology.md); the four that matter most:

1. **Correctness gates timing.** Before any measurement, every design must reproduce the
   same answers, verified against values computed independently from the generated
   dataset. A design that answers wrongly reports no timings at all. *This is not
   theoretical — the gate's first catches were two bugs in the measurement code that
   would otherwise have published confident wrong numbers.*

2. **Results are comparable only within an environment.** Hardware moves benchmark numbers
   more than most design decisions do. Every result names its environment, and every
   environment page documents the specific ways that machine can mislead you.

3. **Plans are captured, not just times.** Wall-clock alone confuses "did less work" with
   "had a warmer cache". PostgreSQL runs carry `EXPLAIN (ANALYZE, BUFFERS)`; YugabyteDB
   runs carry `EXPLAIN (ANALYZE, DIST)`, because on a distributed store the **RPC count**
   is the signal that travels to real hardware and the milliseconds are not.

4. **Reports expire.** A report describes one run, of one code state, on one machine. When
   a re-run changes the conclusions, the superseded report moves to `reports/outdated/` —
   never deleted, never silently edited.

## Reading results honestly

These runs happen on a **laptop**, with the benchmark client sharing eight cores with the
database and no real network between "nodes". That means:

- **Absolute throughput is understated.** Every design pays the same tax, so relative
  comparisons hold.
- **Distributed topologies look better than they would in production**, because
  cross-node consensus costs microseconds over loopback instead of milliseconds across
  availability zones.

Both limits are stated on every report rather than buried, and the environment page lists
three more.

## License

See [LICENSE](LICENSE).
