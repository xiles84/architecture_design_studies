# Escalations — recency question and operational reports (EH-02)

Unmapped decisions raised during execution of [EH-02](HANDOFF.md). One entry per
decision, using `docs/templates/ESCALATION_REQUIRED.md`. Identifiers are `RR-ER-NN`.

LOW appends entries and never edits the handoff; HIGH decides each entry here and
publishes the corresponding amendment in the handoff's **Amendments** section. The
original observation and options stay visible after a decision.

The handoff lists seven concrete triggers (ER trigger 1–7). They are examples, not the
complete set: anything the handoff does not map is an escalation. A benchmark lock held
by another task is **not** one — wait for it, and say so in `PROGRESS.md`.

## Escalation Required — RR-ER-01

- Raised by: Claude Sonnet 5 (LOW role), Claude Code desktop, 2026-09-16.
- Raised at: 2026-09-16T (dev check `am03-dc01-ledger-controls`, PostgreSQL, `pg-single`, `tiny`).
- Handoff ID / revision: EH-02, [AM-03.1](HANDOFF.md#amendments) (the audit statement) and
  [AM-03.3](HANDOFF.md#amendments) (the proof procedure), producing commit `b2830ad`.
- Status: open
- Blocked steps: AM-03.3 step 3's "consistent" proof; AM-03.8's dev checks wherever they
  exercise X1's ledger audit under concurrent writes (race, churn); AM-03.12's phase-3b
  race pair (3b-2, 3b-3), which measures X1 specifically under contention.
- Next setting: **HIGH — decide the ordering key for `a_ledger_attribution` (or accept the
  gap)**

### Observation and evidence

AM-03.3 step 1 (natural control, unfixed SQL) confirmed defect 1 exactly as described,
and worse: `sale_event.customer_id NOT NULL` turns the NULL-customer bug into an outright
failure of every cancellation, not just a silent data-quality problem. With
`-cmd full -phases verify,write,churn -write-ops cancel -race-trials 1` on the unfixed
`x1_cas_ledger`, `-churn-pct 100` forced every successful buy to attempt a cancel: 100/100
attempts errored (`step1d.log`), and manual confirmation shows the update's own
`RETURNING customer_id` is NULL after the update
(`step1-natural-control-unfixed.log`, `step1d.log`; direct `psql` check in this session's
transcript). AM-03.4's fix (deriving the buyer from the ledger row of the sale being
reversed) was applied and rebuilt; cancellations now succeed (`errors=0` across all three
churn tiers, `step3-natural-control-fixed.log`).

**But the audit still reports inconsistent after the fix, which AM-03.3 step 3 did not
expect:**

```
churn      10 seats x10 ...  ledger audit: INCONSISTENT — 0 net mismatches, 3 attribution problems
churn     100 seats x2  ...  ledger audit: INCONSISTENT — 0 net mismatches, 6 attribution problems
churn    1000 seats x1  ...  ledger audit: INCONSISTENT — 0 net mismatches, 6 attribution problems
```

Every reported violation is a false positive in the **audit's own ordering**, not a real
defect in X1. Direct inspection of one flagged seat (event 86, seat 3):

```
 sale_event_id | event_id | seat_no | customer_id |   kind    |              at
---------------+----------+---------+-------------+-----------+-------------------------------
          6968 |       86 |       3 |         920 | sold      | 2026-09-17 00:03:08.870883+00
          6980 |       86 |       3 |        1955 | sold      | 2026-09-17 00:03:08.871494+00
          6978 |       86 |       3 |         920 | cancelled | 2026-09-17 00:03:08.872916+00
```

`sale_event_id` (6968 < 6978 < 6980) gives the only order consistent with the CAS
invariant: sold(920) → cancelled(920) → sold(1955) (a resale can only succeed after the
seat is released). Sorted by `at`, the second sale appears to precede the cancellation
that must have preceded it. **`at` is `now()`, which PostgreSQL fixes at transaction
*start*, not commit** — under the CAS retry loop's concurrent goroutines, a later-starting
transaction can commit before an earlier-starting one reaches its own commit, so `at`
does not track commit order the way the audit's comment assumes ("the newest event ...
Ordered by `at`, then id").

I re-ran the same check with the window ordered by `sale_event_id` alone (no `at` in the
sort key) against the same database state: zero violations. `sale_event_id` (a plain
`BIGSERIAL`, `nextval()` called at each INSERT's own execution, not at transaction start)
tracks real execution order here in a way `at` does not.

### Unmapped decision

AM-03.1 specified this exact SQL, ordered by `at, then id`, **because of a stated concern
about the other engine**: "YSQL caches sequence values per connection, so BIGSERIAL order
is not commit order there." I have not yet tested this statement against YugabyteDB, so I
cannot confirm whether reversing the priority (order by `sale_event_id` first) would
introduce the failure mode AM-03.1 was written to avoid, there. Reordering the audit's
own SQL changes something HIGH specified as an explicit directive ("Append this
statement, with its comments kept in the file"), and doing so wrong would make the
correctness gate the review is trying to strengthen just as unreliable as the one it
replaces, in the other direction. I am not confident this is a decision within LOW's
mapped authority to make unilaterally, given HIGH named a cross-engine concern I cannot
evaluate without a YugabyteDB run I have not yet been asked to do at this step.

Options, as I see them:

1. **Order by `sale_event_id` only.** Fixes the confirmed PostgreSQL false positive.
   Risks reintroducing exactly the YSQL problem AM-03.1's comment warned about, untested.
2. **Keep `at` primary, but widen the tolerance**: treat two rows within some small
   `Δt` (e.g. a few hundred ms, larger than one CAS retry's real latency but smaller than
   any human-scale gap) as unordered, and fall back to `sale_event_id` only inside that
   window. Works on PostgreSQL by the same evidence above (the misordered pair's `at`
   values differ by 0.6 ms); needs a YugabyteDB run to check it does not paper over a
   real YSQL sequence-caching artifact, which could plausibly appear at a similar or
   larger Δt.
3. **Use the engine's own commit-order signal instead of either column** — e.g.
   `pg_current_xact_id()`/`txid_current()` captured by the writing statement itself
   (still fixed at transaction start, so no better than `at`), or restructure so the
   ledger insert records a value guaranteed monotonic with commit (harder: no such
   built-in on either engine without an explicit "commit sequence" table of its own,
   which is a schema change AM-03 does not authorise for X1).
4. **Accept the gap as a stated weakness of the audit under concurrent load** and check
   attribution only in the isolated `write:cancel` op (single ticket, no concurrent
   writers, `at` cannot misorder there) plus a **manual read of a small number of
   post-race seat histories** rather than an automated full-table pass. Cheapest, but
   weakens exactly the automated proof AM-03.1 was written to get.

I have not implemented any of these. I ran option 1 only as a diagnostic (shown above,
not committed) to confirm the ordering is the cause, not the fix.

### Independent work

Not blocked, and continuing: AM-03.5 (post-write `r01`/`r05` checks against harness
counters — reads ticket-table aggregates, not the ledger's chronological order), AM-03.6
(value comparisons; the study 02 `r03` `partial` declaration and REPORTS.md note, already
recorded), AM-03.7 (`ExplainAll` window parameters; `pairs`; diagram; README), AM-03.8's
dev checks for study 03 and for study 02 designs other than X1's race/churn cells, and
AM-03.9's records.

### HIGH decision

<pending>


### LOW addendum, 2026-09-20 — tried to resolve it under the revised autonomy rules; evidence says it cannot be resolved locally

Under the 2026-09-17 rule ("decide, log, continue") I tried to close RR-ER-01 myself instead
of waiting. Evidence, all in `studies/02-ticket-booking/results/devchecks/am03-dc04-rrer01/`:

1. **`w_cancel_ticket`'s own buyer lookup has the same `at` bug as the audit** — not only the
   audit. On PostgreSQL a seat's history was `sold(604) → cancelled(604) → sold(538, earlier
   `at`) → cancelled(**604**)`: the second refund named the *stale* buyer because the subquery
   ordered by `at DESC`. So X1's real write path can attribute a refund to the wrong customer
   under concurrent resale. Counts stay right (r01/r05 checks pass); attribution does not.
2. **Ordering by `sale_event_id` alone (in both the subquery and the audit) fixed PostgreSQL**:
   8/8 high-contention runs (64 buyers, 30–50% churn) plus a standard-parameter run
   (46/46 gate, 15 ledger audits, all report checks) were fully consistent.
3. **It is worse on YugabyteDB**: 53–75 (tier 10) and 234–246 (tier 100) attribution problems,
   *identical across runs* (deterministic, not timing noise). Pinning `ALTER SEQUENCE … CACHE 1`
   had no effect: YSQL still reports `Cache 100`. The original `at, id` ordering had looked
   consistent on YB only because those runs were low-contention.
4. So no single ordering key is correct on both engines, and the PostgreSQL fix regresses
   YugabyteDB. I reverted the three SQL files to the committed (at-ordered) state so the tree
   is not worse on either engine than before; the experiments are kept as evidence only.

**Escalation Required (revised bar: this is real correctness risk in a measured design, on
one engine, with no local fix).**

- Decision needed: how X1 attributes a refund correctly on both engines.
- Why it cannot be deferred: phase 3b would measure X1 with wrong refund attribution.
- Options: (a) **Order-free attribution [LOW's recommendation]**: read the buyer from the ticket
  row itself before the cancelling UPDATE, `WITH old AS (SELECT customer_id FROM ticket WHERE
  ticket_id = $1 AND status='sold' FOR UPDATE), cancelled AS (UPDATE …)`. Sales/cancels of a seat
  already serialize on that row, so no sequence or timestamp is involved. Cost: `FOR UPDATE`
  moves the row lock a statement earlier (the UPDATE takes it anyway), which HIGH's AM-03.4
  rejected as "changes the cancel's locking"; I think that rejection should be revisited given
  this evidence. The audit then needs an order-free formulation too (e.g. each refund's buyer
  must equal a buyer of some earlier sale of that seat, plus the existing net/live checks).
  (b) Restrict X1 to PostgreSQL and report YugabyteDB as an explicit coverage gap.
  (c) Keep PostgreSQL's `sale_event_id` ordering, document YugabyteDB attribution as unproven.
- Work affected: AM-03.10–.13 (calibration and the four measured runs) for X1.
- Work that can safely continue: study 03 entirely (no ledger); study 02's non-X1 designs; the
  study 02 reports matrix (3b-1) if X1 is excluded or clearly flagged.
- Decision Log: D1 (order by `sale_event_id`, PostgreSQL-only) — Medium confidence, **reverted**;
  D2 (`CACHE 1`) — no effect, reverted.

NEXT MODEL: HIGH
