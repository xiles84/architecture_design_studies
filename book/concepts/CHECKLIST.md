# The concept checklist

Every concept chapter of the book answers the same eight questions, so a reader can carry the frame
from one concept to the next and out of the book entirely. The checklist is a reading frame, **not an
exercise set**: the book contains no reader exercises.

For each answer the chapter states an **evidence status**:

| Status | Means |
|---|---|
| `direct` | a registered measured claim answers it |
| `mechanism` | the mechanism is explained and implemented, but the number is not registered |
| `analogy` | transferred from a different engine, family or workload, and labelled as such |
| `gap` | no evidence; shown as a gap rather than answered by analogy |

## The eight questions

1. **What is the invariant?** The rule that must hold, and the failure it forbids.
2. **Where does the information live?** Which table, row, document, column or copy holds the answer.
3. **Where is the contention?** The single row, key, tablet or lock that concurrent work collides on.
4. **What work moved from reads to writes?** What was paid once at write time instead of on every read.
5. **What state is derived?** Which value is computed from other values rather than stored as a fact.
6. **What can become stale?** Which copy or derived value can disagree with its source, and under what
   contract it is allowed to.
7. **Under which workload does this decision stop working?** The boundary — the regime that removes the
   design's premise — and the design that overtakes it there.
8. **Which experiment would prove or disprove the expectation?** The isolation pair, the control that
   must fire, and the sweep that would locate the crossing.

## Current answer status per concept

This table is checked against `book/concepts/*.typ` and the active evidence registry by
`book/evidence/check.sh`; a concept chapter with no row here, or a row naming a chapter that does not
exist, fails the build's source checks.

| Concept chapter | Invariant | Information | Contention | Reads→writes | Derived | Stale | Boundary | Experiment |
|---|---|---|---|---|---|---|---|---|
| `information-placement` | direct | direct | direct | direct | direct | mechanism | gap | direct |
| `indexes-and-hot-rows` | direct | direct | direct | mechanism | mechanism | mechanism | mechanism | direct |
| `concurrency-control` | direct | direct | direct | mechanism | mechanism | mechanism | mechanism | direct |
| `derived-state-and-history` | direct | direct | direct | direct | direct | mechanism | gap | direct |
| `expiry-and-clock-authority` | direct | direct | mechanism | mechanism | mechanism | direct | gap | direct |
| `cache-consistency` | direct | direct | direct | mechanism | mechanism | direct | gap | direct |
| `distributed-placement` | gap | gap | gap | gap | gap | gap | gap | gap |
| `reading-benchmarks` | mechanism | mechanism | mechanism | mechanism | mechanism | mechanism | mechanism | direct |

Two rows are deliberately almost empty: `distributed-placement` has no hardware-independent evidence
(one host, unverified physical placement, single-endpoint history) and `reading-benchmarks` is a method
chapter that teaches how to read a number rather than reporting one.

## Boundary work this checklist exposes

The `gap` and `mechanism` cells are the book's honest to-do list. The largest are:

- **Native datastore families** — every embedding and cache answer is relational plus Redis-as-cache.
- **Physical placement and network** — no real inter-node network, no observed tablet/leader placement.
- **Append-only history under concurrency** — the ledger's own concurrency is not isolated from its CAS.
- **Hard TTL expiry** — crossed in time, never observed to remove an entry.
- **Crossovers** — most sweeps are two points; only the writer sweep and the cardinality sweep bound a
  crossing, and neither carries an across-run interval yet.
