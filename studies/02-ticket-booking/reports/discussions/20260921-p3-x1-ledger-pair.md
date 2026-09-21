---
discussion_id: p3-x1-ledger-pair--deepseek-flash--2026-09-21
analysis_id: 20260921T205212Z--deepseek-flash--2026-09-21
run_id: 20260921T205212Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), Deep Code CLI, no effort setting exposed by the session"
analyzed_at: 2026-09-21
inputs_digest: 6c88f6463c80c6ac
inputs: 20260921T205212Z@6c88f6463c80c6ac, 20260921T183722Z@4a02ec8fb0991166
repo_commit: a04fd2060d5f6a7050a66c252d66149169054e8c
responds_to: 20260913T021206Z--claude-opus-5--2026-09-13, 20260921T183722Z--deepseek-flash--2026-09-21, 20260920T234953Z--deepseek-flash--2026-09-21
---

# Discussion — what does an append-only sale ledger cost the design it is added to?

[Final analysis](20260921T205212Z--deepseek-flash--2026-09-21.md) ·
[Sibling analysis](20260921T183722Z--deepseek-flash--2026-09-21.md) ·
[Reports analysis](20260920T234953Z--deepseek-flash--2026-09-21.md) ·
[Measurements 3b-3](../20260921T205212Z.md) · [Measurements 3b-2](../20260921T183722Z.md) ·
[REPORTS.md protocol](../../REPORTS.md)

## Question and controlled comparison

**One decision:** add an append-only `sale_event` table to P3 (`p3_precreated_cas`), written
in the *same* statement as the sale and the refund. X1 is P3's schema plus that table and
its four indexes, P3's queries plus three report statements answered from the ledger, and
P3's `w_sell_seat` unchanged; only `w_cancel_ticket` differs, and only because a refund
must name the buyer it refunds.

**Held fixed:** arbitration (compare-and-set), the seat rows, the indexes on them, the read
questions, the dataset (seed 42), the per-node budget (2 CPUs / 3 GiB), the client (2
CPUs / 2 GiB, 8 readers), the harness and the gate.

**Inseparable change, stated rather than hidden:** the refund path. P3's cancel is
`UPDATE ticket … WHERE ticket_id = $1 AND status = 'sold'`; X1's is a `FOR UPDATE` read of
that row inside the same statement, then the same update, then a ledger insert. AM-04.1
chose the locking read over a timestamp or a sequence because neither reconstructs commit
order under concurrency (RR-ER-01), and AM-04.2 explicitly includes it in the pair's price
rather than claiming a clean single-decision pair. Per AM-04.2, quoted here so the claim is
checkable:

> **P3 → X1, write axis:** the sell path is P3's, plus one ledger insert in the same
> statement. The refund path is P3's, plus one ledger insert and the locking read that
> insert needs to name the refunded buyer. The sell-out race — the headline experiment —
> exercises only the first.

## Observed measurements

**The pair was run twice**, with the buyer count and the design order changed together
(AM-03.12's 3b-2/3b-3):

| Arm | Buyers | Order | Digest | Race (P3 → X1) | Churn (P3 → X1) |
|---|---:|---|---|---|---|
| 3b-2 `20260921T183722Z` | 32 | P3 first | `4a02ec8fb0991166` | PG 4.1k→5.7k (X1 faster, P3 spread 92 %); yb-single 459→224; yb-cluster3 269→161 | PG 2.9k→4.3k; yb-single 401→170 |
| 3b-3 `20260921T205212Z` | 128 | X1 first | `6c88f6463c80c6ac` | PG 7.0k→3.9k (P3 faster, spreads 4–6 %); yb-single 374→206; yb-cluster3 413→187 | PG 4.6k→3.4k; yb-single 255→145; yb-cluster3 269→131 |

10k-seat tier shown; the pattern holds at every tier. Full tables:
`results/20260921T205212Z/<topology>/x1_cas_ledger.json` and `…/p3_precreated_cas.json`,
and the same names under `results/20260921T183722Z/`. Plans:
`results/20260921T205212Z/<topology>/plans/{p3_precreated_cas,x1_cas_ledger}.txt`.

Derived figures used here: *race throughput* is successful sales over the time from the
first to the last sale, median of three trials, with each trial on a fresh load; *spread* is
(max−min)/median over those three trials; *attempts per sale* is (sales + retries) ÷ sales;
*isolated book/cancel* are ops/s on a fresh load with the crowd absent.

Other measurements in the two arms, all supporting the mechanism below:

- **Cancel, isolated, per engine** — `yb-single` 575→271, `yb-cluster3` 404→164 (2.1x,
  2.5x); PostgreSQL 4.6k→5.2k (the reverse, i.e. unresolved). Book: 379→250,
  258→189, 4.1k→4.3k.
- **Sale latency p99** at 128 buyers, `yb-single`: 471→1382 (10 seats), 1113→2291 (10k).
  Attempts per sale are *unchanged*: 13.07→13.16 (10 seats PG), 1.19→1.20 (10k PG).
- **Organiser edit p99** (an unrelated event-page update every 20 ms during the race):
  PostgreSQL 65/68/66→69/86/99 ms at 1k/10k/100k; `yb-single` down to 196 ms (P3) against
  699 (X1) at the 1k tier.
- **Correctness:** every X1 cell reports 0 ledger mismatches and 0 attribution mismatches at
  load and after every race and churn tier, both engines, both arms — including the
  `yb-cluster3` 128-buyer cell that the pre-AM-04 audit could not handle. Under-booked seats
  after the churn sweep: 0 for both designs on every topology.
- **Costs the ledger pays elsewhere:** +33 % storage (46.6→61.8 MiB on `pg-single`);
  answerability of r01/r03/r05, measured in the sibling reports run
  (`20260920T234953Z@160bd80c49bbd892`).

## Mechanism and alternatives

**Observed work.** The sell statement becomes: the same `UPDATE … WHERE ticket_id = $1 AND
status = 'available' AND version = $2` returning a seat, plus one `INSERT INTO sale_event`.
A seat row is 1 row per seat; the ledger adds one row and one index entry per *sale*
(measured, not inferred: the storage delta and the load delta are in the reports). The
refund statement adds a `SELECT … FOR UPDATE` of the same row it is about to update — an
extra index lookup and lock acquisition on PostgreSQL, an extra round trip on YugabyteDB.

**What the numbers establish.** (a) The ledger's cost is concentrated where the row lock is
contended: the race and churn on YugabyteDB (~2x), the isolated refund on YugabyteDB
(~2x), while isolated book and everything on PostgreSQL at 128 buyers is 1.1–1.5x. (b) The
mechanism is *slower attempts*, not *more attempts*: attempts-per-sale is identical to two
decimal places between the designs, so the ledger is not causing extra CAS retries; it is
lengthening the critical section each retry occupies. (c) The extra write is visible to
unrelated writers on the same instance (organiser edit p99), so it is not a local cost.

**Hypothesis I cannot establish from these two runs.** The PostgreSQL sign flip between
32 and 128 buyers. Candidates, in the order I would test them: (i) design order within a
topology — the first design measured in each topology has been compared in both orders
(P3-first in 3b-2, X1-first in 3b-3) and the *first* design is slower in 3b-2's PostgreSQL
cells; (ii) the interleaved trial structure (race 1, churn 1, race 2, churn 2, race 3,
churn 3), where each trial's fresh load follows the previous trial's writes, so the trial
order interacts with WAL/checkpoint state; (iii) host state (this is one laptop with other
work on it) during whichever cell runs first. None of these is separable from the data
collected, and the report's spread column is the honest witness: it is 92–283 % in exactly
the cells where the ratio flips.

**Alternatives the ledger's design could have used**, and why they are not measured:
a `RETURNING OLD`-based attribution (not available on either engine version here);
ordering the ledger by `at` or by `sale_event_id` (both fail to reconstruct commit order —
`now()` is fixed at transaction start, YSQL hands out sequence blocks per connection, both
demonstrated in `results/devchecks/am03-dc04-rrer01/`); a separate `order` aggregate; event
sourcing as the primary model. The first three are why the refund reads the row; the last
two are named in REPORTS.md §3 as deliberately deferred.

## Exchange with other analyses

- `20260913T021206Z--claude-opus-5--2026-09-13` (study 02's v1 matrix) established that P3
  is one of the two fastest correct designs for a hot drop and that YugabyteDB's numbers
  there are quota-bound rather than engine limits. That is why X1 was built on P3: the
  ledger's cost is measured where it hurts most. I have not asked its author to endorse
  anything here, and the v1 analysis does not cover the ledger.
- `20260920T234953Z--deepseek-flash--2026-09-21` (the reports run, same session as this
  one) is the other half of the same decision: it prices what the ledger *buys*. It also
  found a ~2x cell-level read difference between P3 and X1 in a read-only run, which is
  the same order-effect suspicion this discussion raises for the race — two independent
  sightings of the same measurement hazard, neither of them proof.
- `20260921T183722Z--deepseek-flash--2026-09-21` and
  `20260921T205212Z--deepseek-flash--2026-09-21` are the two arms' analyses and they
  emphasise different arms. That disagreement is deliberate and stays visible: one arm has
  the tighter trials, the other the tighter claim about the disagreement itself.

## Where the measurement is weak

- **No negative control in either arm.** The design list is P3 and X1 only, so no wrong
  design was present to be caught; the audits' sensitivity under this contention is
  inherited, not demonstrated. A gate that passes with no control is weaker than one that
  passes beside a control that fired.
- **Order and buyer count changed together**, which is the confound named above.
- **One complete `yb-cluster3` X1 cell is missing from 3b-2** (YugabyteDB RPC timeout in
  race trial 3); its medians are over 2 trials. The same cell completed in 3b-3.
- **3 trials per tier, one cell per design per topology**, so there is no cell-level
  repetition; the within-cell spread cannot detect a per-cell systematic effect.
- **Quota-bound engines**: the database container was throttled for most of every cell
  (`yb-cluster3` P3 7822/29687 periods, X1 13168/83280), and the denominators differ
  because X1's cells take longer. Ratios between designs measured under the same quota are
  argued; absolute throughput is not.
- **Closed-loop buyers, shared laptop cores, no real network** between the three
  "nodes" — a placement result, not a distributed-systems one.
- **The 100k tier times out** on both YugabyteDB topologies for both designs, so the
  largest drop is measured partially (P3 31 % sold, X1 15–16 %).

## Discriminating next experiment

Run the pair **interleaved inside one cell**: same topology, same dataset, same buyer
count, but alternate the designs tier-by-tier (`p3, x1, p3, x1, …`) with the same number of
trials each, at 32 and at 128 buyers. If the PostgreSQL ratio stays with the design (X1
slower at both), the ledger's cost is real and the 3b-2 result was host state; if the ratio
follows the *position* in the sequence, the metric is measuring the machine rather than the
design and every P3 → X1 number in the study needs re-running with the orders balanced. The
outcome would change conclusion 1 of the 3b-3 analysis from a range (1x–2.3x) to a single
figure, and it would retire or confirm the second sighting of the cell-order hazard in the
reports run.
