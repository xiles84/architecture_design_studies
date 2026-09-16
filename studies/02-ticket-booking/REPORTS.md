# Study 02 v2 — the reports a ticketing system is actually asked for

**Protocol.** Written by HIGH before any code. Nothing here is a result.

**Owner's request (2026-09-15):** *"In the other cases, let's include reports and queries
that similar real cases like ours need and let's integrate them to our scenarios and
tests."*

Study 02 measures five read questions today: availability, a band's events, a customer's
tickets, a ticket by id, tickets sold across a band. They are the questions a *buyer's*
screen asks. A real ticketing company also runs a back office, and its questions land on
the same tables that the overbooking designs reshaped — which makes them design-
discriminating in a way the buyer's screens are not.

---

## 1. The reports

Added to **every** design's `queries.sql`. **No schema and no index changes** to the
fourteen existing designs: their published write, race and storage numbers must stay
valid, and what these reports cost *on the indexes a transactional design already has*
is exactly the interesting number.

| # | Name | Params | The real question |
|---|---|---|---|
| r01 | `r01_event_sales_window` | `event_id`, `since`, `until` | how the drop went: tickets sold and revenue inside a window |
| r02 | `r02_event_recent_buyers` | `event_id` | the last 50 sales with the buyer — the operations feed, and the first thing looked at when something goes wrong |
| r03 | `r03_customers_last_purchase_window` | `since`, `until` | customers whose **last** purchase falls in the window — the same recency question study 01 asks of donors ([`../01-charity-tree/RECENCY.md`](../01-charity-tree/RECENCY.md) §2), with the same two window regimes |
| r04 | `r04_band_sellthrough` | `band_id` | per event: capacity, sold, percentage — the management report that decides the next tour |
| r05 | `r05_refunds_window` | `since`, `until` | cancellations and the money returned, inside a window |
| r06 | `r06_outstanding_holds` | `event_id` | how much capacity is sitting in baskets right now (H designs only) |

r03 uses the trailing and historical regimes defined in study 01's protocol, for the
same reason: in a trailing window "last purchase in the window" and "purchased in the
window" accidentally coincide, and the difference only appears in a historical one.

## 2. The finding this will produce before it is run

**No current design can answer r05, and several cannot answer r01 honestly.**
`w_cancel_ticket` sets the row back to `status='available', customer_id=NULL,
sold_at=NULL` (P family), or removes the ticket (C and R families). The sale is not
recorded anywhere after it is undone. So:

- **r05 is unanswerable in all fourteen designs.** Not slow — impossible.
- **r01 is answerable only for sales that are still live**: a ticket sold and refunded
  inside the window vanishes from the window's revenue, which is not what a finance
  report means by "sold in this period".
- r02 and r04 are answerable everywhere, at costs that differ by design.
- r06 is answerable only in H0/H1, which are the only designs with a hold.

That is a genuine architectural result and belongs in the study rather than being
designed away: **an inventory model chosen purely to prevent overbooking can make the
business's own accounting unanswerable**, and the fix costs something on the hot path.
The measurement is what that something is.

Implementation records, per design and mechanically (from the SQL, not from opinion),
whether each report is `answerable`, `partial` (a stated caveat, like r01 above) or
`unanswerable`, and the generated report prints that table beside the throughput one. An
`unanswerable` report is not timed and is never shown as a blank or a zero.

## 3. The one new design

| ID | Directory | From | The single change |
|---|---|---|---|
| **X1** | `x1_cas_ledger` | P3 (`p3_precreated_cas`) | an append-only `sale_event` table — `(event_id, seat_no, customer_id, kind ∈ {sold, cancelled}, at, price_cents)` — written in the **same transaction** as the sale and the cancellation. r01, r03 and r05 are answered from it |

P3 is the base because the published analysis makes it one of the two fastest correct
designs for a hot drop: the ledger's cost is measured where it hurts most, not where it
would disappear into a slower design's noise.

What X1 measures:

- the cost of a second insert in the selling transaction, in the **sell-out race** (the
  headline experiment) and in isolated writes — one extra row write, one more index to
  maintain, one more thing to contend on;
- what the ledger buys: r01, r03 and r05 become answerable, and r01 becomes *correct*
  across refunds;
- whether the ledger's own append point becomes hot at the 100 000-seat tier.

**Controlled pair: P3 → X1** — identical arbitration, identical seat rows, identical
queries except the three answered from the ledger. One decision.

**Deliberately not done in this iteration** (HIGH decision, to keep the measured matrix
affordable): a second ledger variant on C4, a separate `order` aggregate, and event
sourcing as the primary model. If X1's race cost is small, the C4 variant is the obvious
next rung and the analysis will say so.

## 4. Correctness

- The existing correctness gate is extended with the new reports: each answerable report
  is checked against values computed independently in Go from the generated dataset,
  before any timing.
- X1 gets a **ledger reconciliation audit** after every writing phase: the sales and
  cancellations in `sale_event` must reconstruct exactly the live ticket state, and must
  match the sales the harness saw commit. Study 02 already reconciles harness-observed
  sales against tickets; this extends the same check to the ledger, and it is the check
  that would catch a ledger written outside the selling transaction.
- The negative controls (C1, H0) are unchanged and must still fire.

## 5. Sizing

Unchanged from study 02's established budget: 2 CPUs / 3 GiB per database node, 2 CPUs /
2 GiB for the client, 32 buyers per race event, readers at the study default. The report
queries are added to the **read** phase, which is not the contended one, so they do not
change the race's sizing. Client saturation for r01/r03/r05 is measured in the dev
checks before the matrix is sized, and recorded.

Requirement mapping (AGENTS.md "Required comparisons"): this is an extension of an
existing study's measurement plan, so it records coverage rather than re-running the
study. **1 — sizing:** unchanged and restated above. **2 — rollup/rolldown/embedding:**
the ledger is neither; it is an append-only record, and the study's existing counter
(C4), inventory row (R1) and bucket (R2) designs already cover materialised aggregates.
r04's sell-through is the rollup-versus-derive question for this study, and it is
measured against C4/R1/R2's stored counters as they stand. **3 — colocation:** not
extended here; study 02's existing YugabyteDB topologies are unchanged, and this is
recorded as an untouched gap. **4 — concurrency strategies:** unchanged; P3 (optimistic
CAS) versus P1/C3 (pessimistic) is already the study's spine, and X1 inherits P3's.

## 6. What this will not be able to say

- Nothing about a production SLO (closed loop, shared cores, quota-bound YugabyteDB).
- Nothing about a reporting workload running *concurrently* with a drop: the reports are
  measured in the read phase. Whether a finance query throttles a sell-out is a separate
  experiment, and the analysis will list it.
- Nothing about long-horizon ledger growth: the ledger is as old as the run.
