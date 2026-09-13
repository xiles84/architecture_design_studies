---
discussion_id: counter-cost-vs-index-cost--claude-opus-5--2026-09-13
analysis_id: 20260913T021206Z--claude-opus-5--2026-09-13
run_id: 20260913T021206Z
environment: host-zenbook-ux5406sa
analyst: claude-opus-5
analyst_kind: ai
analyst_version: "Claude Opus 5 (model id claude-opus-5), in the Claude Code desktop app; the session that designed and ran study 02"
analyzed_at: 2026-09-13
inputs_digest: 71bcee71d725d43d
inputs: 20260913T021206Z@71bcee71d725d43d
repo_commit: 7570648c5220
responds_to: the project owner's question of 2026-09-13 (below); no other analysis
---

# Discussion — why doesn't the counter's cost reappear as index maintenance?

[Final analysis](../analyses/20260913T021206Z--claude-opus-5--2026-09-13.md) ·
[Measurements](../20260913T021206Z.md)

## Question and controlled comparison

The owner asked: P3 (compare-and-set on pre-created seats) outperformed P4 (SKIP LOCKED plus
a counter), C3 (lock the event, then count) and C4 (guarded counter). The counter plausibly
costs something — but P3 must also pay, on every sale, for maintaining the index used to find
an available seat. Why did that cost not cancel the counter's?

**Short answer: the counter's cost is not work, it is waiting.** Updating a counter is cheap.
What is expensive is that every buyer of an event must update the *same row*, a row lock is
held until the transaction commits, and so buyers form a queue whose speed is set by the
whole transaction's duration. Index maintenance is real work too, but each sale touches
*different* index entries, and the structure that protects an index page — a latch — is held
only for the instant the page is modified, not until commit. Work done in parallel does not
queue; work on one shared row does.

The cleanest controlled pair is **P2 → P4**: identical tables, identical indexes, identical
seat pick (`FOR UPDATE SKIP LOCKED`) and identical sale; P4 adds exactly one statement,
`UPDATE event SET seats_sold = seats_sold + 1`, in the same transaction. Any index cost is
the same on both sides of that pair. C3, C4 and R1 show the same effect without pre-created
seats; P3 differs from P2 in more than one way (no locking read and fewer round trips — see
the final analysis, conclusion 3), so it is the fastest example, not the controlled one.

## Observed measurements

From the [report](../20260913T021206Z.md), PostgreSQL 1-node, sell-out race:

| | P3 CAS | P2 skip-locked | P4 = P2 + counter | C4 counter | R1 counter, own row | C3 lock + count |
|---|---:|---:|---:|---:|---:|---:|
| sold/s, 1k seats | 6.2k | 3.6k | 816 | 854 | 866 | 953 |
| sold/s, 10k seats | 5.9k | 3.6k | 842 | 831 | 834 | 721 |
| sold/s, 100k seats | 6.9k | 3.4k | 853 ⏱ | 835 ⏱ | 868 ⏱ | 485 ⏱ |
| attempts per sale, 10k | 1.09 | 1.00 | 1.06 | 1.00 | 1.00 | 1.00 |
| sale latency p99 (ms), 10k | 55 | 65 | 209 | 171 | 171 | 64 |
| DB container: throttled / runnable CFS periods, whole cell | 861 / 1 093 (79%) | 1 133 / 1 275 (89%) | 1 500 / 1 944 (77%) | 517 / 1 811 (29%) | 468 / 1 807 (26%) | 638 / 1 847 (35%) |

⏱ = the 100 000-seat race hit its 60 s timeout.

Four observations follow directly:

1. **Adding the counter alone cut P2 by 4.3x** (3.6k → 842/s at 10k seats) while leaving every
   index and every other statement unchanged.
2. **Every design that serialises buyers on one row converges at ~820–900/s** — P4 with
   pre-created tickets, C4 and R1 with created ones, at every event size. A ceiling shared by
   designs that do different amounts of work is a queue's ceiling, not a cost of the work.
3. **They do not retry; they wait.** C4, R1 and C3 have exactly 1.00 attempts per sale, and
   their p99 sale latency is 3x P3's.
4. **The database was idle-ish during the counter designs.** Their container was throttled
   in only 26–35% of CPU periods, against 79–89% for P2 and P3. The counter designs left CPU
   unused because buyers were blocked on a lock, not busy. (P4's 77% is the exception that
   fits: P4 also contains P2's spin on "all seats locked" for small events, which burns CPU
   on the ten-seat tier of the same cell.)

And from the captured plans
([`c4_counter_guard.txt`](../../results/20260913T021206Z/pg-single/plans/c4_counter_guard.txt),
[`p4_precreated_counter.txt`](../../results/20260913T021206Z/pg-single/plans/p4_precreated_counter.txt)):
the counter update executes in **0.046 ms** (15 buffer hits), the ticket insert with its
foreign-key check in **0.070 ms**. If C4's 854 sales/s are strictly serial on the event row,
each transaction holds that row for **~1.2 ms** — roughly **ten times the database work it
does on the row**.

## Mechanism and alternatives

**Where the counter's time goes (hypothesis, consistent with the numbers).** A C4 sale is
BEGIN → guarded `UPDATE event` (row lock taken) → `INSERT ticket` → COMMIT (row lock released,
after the commit record is flushed). The lock is held across two client–server round trips,
the insert, and the commit flush. The row work is ~0.05 ms; the other ~1.1 ms is round trips
through a CPU-throttled client and server and the commit. Every second buyer spends that
time waiting, then the next, so the event sells at 1 / (lock hold time) no matter how many
buyers or CPUs there are. Moving the counter to its own row (R1) does not change this at all
(866 vs 854/s), which also rules out "the event row is wide" as the explanation.

**Why P3's index maintenance does not queue.** Selling a pre-created seat is an `UPDATE` of
`status`, `customer_id` and `sold_at`. In PostgreSQL that is **not** a HOT update, because
`status` appears in the partial indexes' predicates and `customer_id` is indexed. So each
sale does pay index work, as the owner expected: new entries for the new row version in
`ticket_pkey`, `ticket_event_seat_uq` and `ticket_customer_idx` (the partial "available" index
gets no new entry; the old one becomes dead until vacuum). But:

- those entries belong to **different seats**, so concurrent sales modify different
  positions of the B-trees, mostly different leaf pages;
- a B-tree page is protected by a **latch**, held for the microseconds it takes to change the
  page and released immediately — **not held until COMMIT** like a row lock;
- in P3 even the row lock on the sold seat lives only for one autocommit statement.

So P3's index cost is paid **in parallel** by 32 buyers on the CPU they have — which is why P3
is the design that saturates the CPU (79% of periods throttled) instead of leaving it idle.

**Where P3's toll does reappear — as retries, not waiting.** When two buyers pick the same
seat, one compare-and-set matches zero rows and retries. That is P3's form of contention, and
it scales with how crowded the free seats are: 3.65 attempts per sale on ten-seat events,
1.83 at 100 seats, 1.02 at 100 000. It is visible in throughput exactly where expected — P3's
slowest tier is the ten-seat one (2.5k/s). A lost CAS costs one short statement and no one
waits for it.

**Dead index entries** are a second, smaller place the toll reappears: sold seats leave dead
entries in `ticket_available_idx` that later candidate searches may step over until vacuum
(or PostgreSQL marks them dead on first visit). That grows the cost of a search, but a search
holds no lock anyone else needs.

**C3 is the counter's problem plus a growing one.** It locks the event row, then counts the
tickets inside the lock. The count's cost grows with tickets sold, so the lock hold time grows
with it: 953/s at 1k seats, 485/s at 100 000.

**YugabyteDB, same shape.** DocDB detects conflicts per key. Index writes for different seats
are different keys and do not conflict; the counter is one key every buyer writes. The
counter designs sold 15.5–22 seats/s on both YugabyteDB topologies against P3's 122–402/s, and
the organiser's unrelated edit to the event row waited 2.0–2.4 s behind C4's buyers versus
23–62 ms when the counter had its own row (R1).

**Alternatives considered.**
- *The counter update is expensive* — contradicted by the plan (0.046 ms) and by P4 ≈ C4 ≈ R1.
- *Index maintenance is cheap enough to ignore* — not shown: the captured P3 sale plan was
  taken with an already-sold ticket and matched zero rows, so it does not measure index
  maintenance, and P3's CPU saturation says the work is real. The argument does not need it
  to be cheap, only to be parallel.
- *P3 wins only because it uses fewer round trips* — a real confound for P3 vs P2 (final
  analysis, conclusion 3), but not for the counter question: P2 → P4 has the same round-trip
  structure except one added statement inside the locked transaction.

## Exchange with other analyses

No other analysis of digest `71bcee71d725d43d` exists yet. This companion answers a question
from the project owner and extends conclusions 2 and 3 of the linked final analysis.

## Where the measurement is weak

- **Lock hold time is inferred, not measured.** ~1.2 ms comes from throughput under the
  assumption of strict serialisation; the harness does not record lock waits
  (`pg_locks`/wait events) or per-statement round-trip time.
- **Commit flush versus round trips is not separated.** Both lie inside the lock window; this
  run cannot say which dominates.
- **CPU-throttling counters cover the whole cell**, not the race phase alone, so they mix the
  race with reads, writes and loads (P4's figure shows how a cell's other tiers can dominate).
- **One trial**, 32 buyers, single YugabyteDB query node — as in the final analysis.
- The **index work of a P3 sale is not isolated** by any measurement here (see alternatives).

## Discriminating next experiment

1. **Shorten the lock window without changing the counter.** Add a C4 variant that claims the
   seat and inserts the ticket in **one autocommit statement**
   (`WITH s AS (UPDATE event … WHERE seats_sold < capacity RETURNING event_id) INSERT INTO
   ticket … SELECT … FROM s`). If waiting on round trips dominates, it should sell several
   times faster than C4 while staying serialised; if it barely moves, the commit flush
   dominates.
2. **Diagnostic, not a design:** repeat C4 and P3 with `synchronous_commit = off`. A large C4
   gain and a small P3 gain would confirm that time spent holding the lock through commit is
   the cost.
3. **Buyers sweep** (`-race-buyers 4,16,64`): a queue-bound design stays flat as buyers grow;
   a CPU-bound one (P3) rises until the CPU saturates.
4. **Record lock waits**: sample `pg_stat_activity.wait_event` during the race; counter designs
   should show buyers in `Lock: transactionid`/`tuple`, P3 should not.
