# Study 03 v2 — the reports a reserved-seating operation is actually asked for

**Protocol.** Written by HIGH before any code. Nothing here is a result.

**Owner's request (2026-09-15):** *"In the other cases, let's include reports and queries
that similar real cases like ours need and let's integrate them to our scenarios and
tests."*

Study 03 measures six read questions: the section seat map, an event's sections, its
availability, the seats of one hold, a customer's tickets, a ticket by id. Those are the
buyer's screens. A venue also runs a box office and an operations desk during the drop,
and their questions land on exactly the rows the hold mechanism contends on.

---

## 1. The reports

Added to **every** design's `queries.sql`. **No schema and no index changes** to the
fourteen designs: the study is measured and signed (`study-03/v1-analysis`), and its
race, lifecycle and storage numbers must stay valid. What these reports cost on the
indexes each design already has is the number worth having.

| # | Name | Params | The real question |
|---|---|---|---|
| r01 | `r01_holds_expiring_soon` | `event_id`, `within` | which holds expire in the next N minutes — the operations desk's view, and what actually sizes a sweeper |
| r02 | `r02_seat_status_lookup` | `event_id`, `seat_id` | who holds or owns this seat right now — the call the box office makes while a customer is on the phone |
| r03 | `r03_section_sales_window` | `event_id`, `since`, `until` | confirmed sales per section inside a window — selling pace, section by section |
| r04 | `r04_event_recent_confirmations` | `event_id` | the last 50 confirmed sales — the operations feed |
| r05 | `r05_customers_last_purchase_window` | `since`, `until` | customers whose **last** purchase falls in the window — the recency question of [`../01-charity-tree/RECENCY.md`](../01-charity-tree/RECENCY.md) §2, in this domain, with the same two window regimes |
| r06 | `r06_hold_funnel_window` | `event_id`, `since`, `until` | holds created, confirmed, expired and abandoned inside a window — the funnel that says whether 40 minutes is the right number |

## 2. The finding this will produce before it is run

**r06 is unanswerable in every design, and r01 means different things in different
ones.** `w_release_expired` sets the seat row back to
`status='available', hold_id=NULL, customer_id=NULL, hold_expires_at=NULL`: the hold that
expired leaves no trace at all. So an expired hold is indistinguishable from a seat that
was never held, and the abandonment rate — the number that decides the length of the
checkout window — cannot be computed from the database. The same erasure is why a
released hold cannot be told apart from an expired one.

r01 is answerable wherever a hold is visible on the seat row, but in **E0** (application
clock) the answer depends on whose clock is asked, which is the very thing that design
exists to expose; the report states this instead of printing one number.

Implementation records, per design and mechanically from the SQL, whether each report is
`answerable`, `partial` (with the caveat text) or `unanswerable`. An unanswerable report
is not timed and never appears as a zero.

This matters more here than in study 02: study 03's signed analysis reports transient
refusals and sweeper behaviour, and a venue cannot audit either of them after the fact
from the state these designs keep.

## 3. No new design in this iteration — and why

The obvious fix is a hold-history table (`hold_event`: created, confirmed, released,
expired, by whom, on whose clock) written in the same transaction as each transition,
which would make r06 answerable and give the transient-refusal question an audit trail.
**HIGH's decision is to specify it and not build it yet**, because it touches all five
E/K designs' write paths, and study 03's main matrix costs seventeen hours — a change of
that size needs its own protocol, its own controlled pairs (hold history on / off, one
decision) and its own measured run, not a rider on a reporting change.

What is done here instead: the reports are added and measured, the answerability table is
published, and the gap is written down as the study's next design question with the
mechanism already specified. The analysis will name it as the first follow-up.

## 4. Correctness

- The correctness gate is extended with the new reports: each answerable report is
  checked against values computed independently in Go from the generated dataset, before
  any timing.
- The invariants INV-1..INV-8 are unchanged, and the negative controls (S0 check-then-act,
  E0 application clock, K0 naive confirmation) must still fire in the phases where they
  fired before. A reporting change that made a control stop firing would be a harness bug,
  and the dev checks check exactly that.
- r01 and r02 are read against the *same* rows the race contends on. They are measured in
  the read phase only; running them inside the race is a separate experiment (§6).

## 5. Sizing

Unchanged from study 03's established budget: 2 CPUs / 3 GiB per database node, 2 CPUs /
2 GiB for the client, connections spread over all three YugabyteDB nodes, race buyers as
published. The reports are added to the read phase, so the race's sizing is untouched.
Client saturation for r03/r05/r06 is measured in the dev checks before the matrix is
sized, and recorded.

Requirement mapping (AGENTS.md "Required comparisons"), as an extension of an existing
measurement plan: **1 — sizing** restated above, unchanged. **2 — rollup/rolldown/
embedding:** covered by the study's existing layout designs — L1 claim rows, L2 the
section document (embedding), L3 section-sharded — which these reports are measured
across; the reports add no new materialisation. **3 — colocation:** unchanged; L3's
placement lever is already the study's colocation axis and the reports are measured on
it. **4 — concurrency strategies:** unchanged; S1 (conditional update, optimistic) versus
S2/S3 (lock, NOWAIT — pessimistic) remains the study's spine and no report touches the
write path.

## 6. What this will not be able to say

- Nothing about a production SLO, for the reasons the signed analysis already states.
- Nothing about reports running **during** a drop: they are measured in the read phase.
  Whether the operations desk's `r01` competes with the buyers it is watching is a real
  question and an experiment this does not run.
- Nothing about abandonment, expiry rates or refusal history: the state needed to answer
  them is not kept (§2). This is a measured gap, not an oversight.
