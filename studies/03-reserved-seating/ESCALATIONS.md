# Study 03 — Escalations

Decisions the [Execution Handoff](HANDOFF.md) did not map, raised by an executor session
(Opus 5, `high`) for the planning session (Opus 5, `ultracode`) to decide.

## Protocol

1. **Executor:** append a new `ER-NN` entry using the template below. Never edit an existing
   entry. Commit it with a message starting `Study 03: Escalation Required ER-NN`. Continue
   with unblocked work. If nothing is unblocked, end the session with:
   `Escalation Required: ER-NN — switch to Opus 5, ultracode`.
2. **Ultra session:** append a `Decision` block under the entry. Never rewrite the
   executor's text. If the decision changes the specification, also append an amendment to
   `HANDOFF.md` §16 and tag it `study-03/v0.N-handoff-amendment-NN`. Commit.
3. **Executor, next session:** read every entry with a new decision before resuming.

A decision records the model that took it, and when. Entries stay in this file for good:
they explain why the study looks the way it does.

## Template

```markdown
## ER-NN — <short title>

- raised_by: <model id, setting, tool>
- raised_at: <YYYY-MM-DD HH:MM UTC>
- handoff_sections: <§ numbers>
- trigger: <§12 item number, or "unmapped">
- status: open

**Blocked:** what cannot proceed.

**Context:** what was observed, with file paths, commands and outputs (quote them; link
logs under results/devchecks/).

**Options:**
1. <option> — consequence for what is measured and how results are judged
2. <option> — consequence

**Executor's recommendation:** <option and why, or "none">

**Work continuing meanwhile:** <steps>

### Decision

- decided_by: <model id, setting, tool>
- decided_at: <YYYY-MM-DD HH:MM UTC>
- decision: <option or new instruction>
- rationale:
- handoff_amendment: <AM-NN, or none>
```

---

## ER-01 — YugabyteDB: a confirmation's UPDATE misses a valid hold that the same transaction then reads

- raised_by: claude-opus-5, setting `high`, Claude Code desktop (executor session)
- raised_at: 2026-09-14 15:10 UTC
- handoff_sections: §3.4 (INV-3), §10.3 (dev-check pass criteria), §12 trigger 2
- trigger: §12 item 2 — a correct design violates an invariant in dev checks
- status: open

**Blocked:** tagging `study-03/v1-harness` and starting the main matrix (steps 8–9). §10.3
requires zero INV-1..INV-7 violations in designs meant to be correct; on YugabyteDB 1-node
the seat-row designs record occasional **early rejections** (INV-3), so the criterion fails.

**Context.** Every occurrence has the same shape:

1. A buyer's hold commits (conditional multi-row `UPDATE`, 1–4 seats, real TTL 40 minutes).
2. About 10–130 ms later the same buyer confirms, in a new READ COMMITTED transaction:
   `SELECT now()`, then S1's `w_confirm_seats` (`UPDATE event_seat SET status='sold' ...
   WHERE event_id AND section_no AND seat_id = ANY($seats) AND hold_id = $h AND status =
   'held' AND hold_expires_at > now() RETURNING ...`).
3. The UPDATE takes about 90–500 ms and **matches 0 rows**. The harness rolls back and the
   ledger records an early rejection, "refused 39m59.9s before the hold's expiry".
4. **Inside the same transaction, before the rollback**, a plain SELECT of the exact seat
   ids the UPDATE was given shows every seat `held`, by this hold, `hold_expires_at >
   now()` = true. A read after the rollback shows the same.

Verbatim example (`results/devchecks/dc10-yb1-early-intx/yb-single/s1_conditional_update.json`):

> refused 39m59.868s before the hold's expiry; the confirming statement matched 0 of 2
> seats (seat ids [5522 5523]); right after, the database showed 2 of the hold's seats
> validly held by it; same seats re-read inside the refusing transaction: [5522 held this
> hold valid=true; 5523 held this hold valid=true] at now() 2026-09-14T15:02:50.821789Z;
> re-read after: [...same...] at now() 2026-09-14T15:02:50.986288Z; grant now()
> 2026-09-14T15:02:50.689923Z, confirm now() 2026-09-14T15:02:50.821789Z

The captured YugabyteDB plan for `w_confirm_seats` is a primary-key Index Scan with the
whole predicate pushed down as a `Storage Filter`, `hold_expires_at > now()` included
(`results/devchecks/dc5-yb1-all/yb-single/plans/s1_conditional_update.txt`).

**Where it has been seen** (all `tiny`, commit `b1b5812` or earlier; counts are refused
confirmations of holds with at least G left):

| Run | Engine | What ran | Occurrences |
|---|---|---|---|
| dc3, dc4 | PostgreSQL 1-node | 12 designs, all phases | **0** |
| dc5 | YugabyteDB 1-node | 13 designs, all phases | S1 1, S2 1, S3 1 (1 000-seat race), K1 1; 0 in E1, E2, L1, L2, L3, S4; E0 3 (control) |
| dc6 | YugabyteDB 1-node | S1, S2, S3, K1; race 1k + 10k | 0 |
| dc7 | YugabyteDB 1-node | S1, all phases | 0 |
| dc8 | YugabyteDB 1-node | S1, S2; 4 × 10k race each | S1 1, S2 1 |
| dc9 | YugabyteDB 1-node | S1; 6 × 10k race | 3 (at 3.4 s, 6.5 s and 19 s after the gate — not tied to the deadline) |
| dc10 | YugabyteDB 1-node | S1; 6 × 10k race | 2 |

Roughly one per two 10 000-seat races, out of ≈1 500 immediate confirmations per race;
always in the race (hot contention for the best blocks), never yet in the lifecycle of a
correct design. The ledger shows no theft, no ticket mismatch and no partial hold at any
occurrence, and every deferred 40-minute confirmation in these races succeeded.

**Checks that rule out the harness** (§12 item 2 asks for two):

1. Ledger unit tests pass (17), including rejection classes and lock-wait ordering.
2. The in-transaction re-read (item 4 above) proves the seat ids, hold id and validity the
   UPDATE was given match rows the transaction itself can see. `section_no` is not in that
   re-read, but the hold's own UPDATE carried the same `section_no` and matched.

**A standalone reproduction did not reproduce it** (`probe/` P13,
`results/devchecks/engine-probe-*.md`): hold-then-confirm loops with 32 workers for 3 minutes
on YugabyteDB, with and without hot-block contention, S1's partial expiry index and CHECK
constraints, and pools of 8 and 80 connections — 20 011, 7 232 and 5 673 holds, 0 anomalies
(PostgreSQL: 136 060 holds, 0). What the race has and P13 lacks: 32 buyers also running
10 000-row section aggregates (`q02`) and seat-map reads, a fully loaded catalogue, FKs,
CPU-quota saturation of the YugabyteDB container.

**Options:**

1. **Record it as a measured engine behaviour and run the matrix unchanged.** Keep every
   design's SQL. The report counts these refusals and marks the cells, and the analysis
   discusses them as a YugabyteDB READ COMMITTED behaviour under load that breaks the
   owner's guarantee for seat-row designs. Needs an amendment to §10.3 so the matrix may
   start with this class present on YugabyteDB. *Consequence:* YugabyteDB cells of
   S1/S2/S3/K1 (and probably L3, E1, E0) carry ❌ for INV-3; no conclusion is hidden.
2. **Report it as its own class**, "refused although valid inside the same transaction",
   separate from early rejections caused by theft or expiry. The harness already records
   the evidence for each occurrence. Combinable with option 1. *Consequence:* a report
   change plus a ledger/report field; distinguishes an engine anomaly from a design failure.
3. **Time-boxed investigation before the matrix** (for example 2 hours): keep the
   YugabyteDB container after a reproducing race and search the tserver and postgres logs
   at the occurrence timestamps (study 02 lesson: search the whole log); try the same race
   on yb-cluster3; try the confirmation without the pushed-down `now()` predicate as a
   diagnostic (not a design change); try the standalone probe with q02-style aggregate load.
   *Consequence:* delays the matrix; may identify a YugabyteDB defect worth reporting upstream.
4. **Add a new design** (not a fix of S1): "S1r — confirmation re-reads and retries once
   when it matches fewer rows than a hold it can still see valid". *Consequence:* measures
   the price of defending against the anomaly; must not replace S1 (the handoff forbids
   fixing a design so it passes).

**Executor's recommendation:** 1 + 2, with the 3 cluster check folded into the yb-cluster3
dev check. The behaviour is reproducible within the study's own workload, the evidence per
occurrence is complete, and it is exactly the kind of "the seat was gone although my hold
was valid" failure the owner asked this study to find. 4 is worth considering as a
follow-up, not a blocker.

**Work continuing meanwhile:** yb-cluster3 dev checks (§10.3 subset), rendering diagrams,
README and PROGRESS updates. No `small` run is started.

## ER-01 addendum — yb-cluster3 dev check (new evidence, same question)

- raised_by: claude-opus-5, setting `high`, Claude Code desktop (executor session)
- raised_at: 2026-09-14 16:35 UTC
- relates_to: ER-01 (this is not a separate decision; it adds evidence for it)

`results/devchecks/dc11-yb3-subset` (tiny, YugabyteDB 3 nodes RF=3, connections spread over all
three nodes, commit `d7bb6c3`), race phase, refused confirmations of holds with at least G
left:

| Design | Race early rejections | Note |
|---|---:|---|
| S1 conditional-update | 3 | 1 000- and 10 000-seat races |
| K1 payment-window | 6 | 100-, 1 000- and 10 000-seat races |
| L3 section-sharded | 7 | 100-, 1 000- and 10 000-seat races |
| L1 claim-rows | 7 | the upsert layout: `UPDATE seat_claim ... WHERE hold_id AND expires_at > now()` |
| E1 sweeper-expiry | 12 | 10-, 1 000- and 10 000-seat races |
| L2 section-document | **0** | its check runs in the application on a read, then a version compare-and-set |
| K0 naive-confirm (control) | **0** | its confirmation filters on seat ids only, not on hold, status or expiry |
| E0 (control) | 12 | |
| S0 (control) | 152 | theft, as designed |

The per-occurrence evidence has the same shape as on one node. For example, S-family: "matched 0 of
6 seats (seat ids [480 ... 485]) ... re-read inside the refusing transaction: [480 held this
hold valid=true; ... 485 held this hold valid=true]". The lifecycle of every correct design had
no early rejection on the cluster.

Two facts that may help the decision:

1. **The class is wider than seat-row designs.** L1's claim table shows it too. The designs
   whose refusing statement does not filter on the hold's own just-written columns (K0 by seat
   id only, L2 by application-side check plus version) show none, in 8 races each.
2. **It is more frequent on three nodes** (every tier from 100 seats) than on one (10 000-seat
   races, one 1 000-seat case).

Everything else on yb-cluster3 met §10.3: every gate 51/51. The controls fired: S0 in the
race and the lifecycle, E0 in the lifecycle, K0 in the lifecycle (6 late sales, 5 sales
without the hold), and the E1 outage probe (3 of 3 unavailable on both tiers).

### Decision (ER-01)

- decided_by: claude-opus-5, setting `ultracode`, Claude Code desktop (planning session)
- decided_at: 2026-09-14 17:40 UTC
- status: **decided**
- decision: **options 1 + 2 + 4**, after a time-boxed investigation (option 3, done by this
  session). It is recorded as a measured YugabyteDB behaviour, not a harness bug. It is
  reported as its own class, the **transient refusal**, and a new design **S1r** measures
  the defence. The matrix starts once the AM-01 dev checks pass.
- handoff_amendment: **AM-01** (HANDOFF.md §16), tag `study-03/v0.1-handoff-amendment-01`

**Investigation** (`diagnose-er01.sh`; S1; YugabyteDB 1 node; `tiny`; 10 000-seat race;
8 races per configuration; `results/devchecks/er01-20260914T164638Z` and
`results/devchecks/er01-20260914T170750Z`):

| Configuration | What it rules in or out | Early rejections |
|---|---|---:|
| A default | reproduces within the study's own workload | 2 |
| B `yb_enable_expression_pushdown=off` (the plan's `Storage Filter` became a query-layer `Filter`, checked in the captured plan) | not pushdown | 2 |
| C tserver `enable_wait_queues=false` | not wait-on-conflict | 2 |
| D default + `yb_debug_log_internal_restarts=on`, harness re-issues the refused statement in the same transaction | the miss is **transient**; the refusing backend logged no restart or error at that moment | 1 |
| E `yb_max_query_layer_retries=0` + restart logging, same re-issue | not an artefact of internal statement retries | 1 |

The evidence for D and E is the same, per occurrence:

1. The confirmation's `UPDATE` matched **0 of N** seats (N = 2 in D, 6 in E).
2. A `SELECT` of the same seats inside the same transaction showed every seat held by this
   hold, valid.
3. **The identical `UPDATE`, issued again in the same transaction, matched N of N.**
4. The refusing backend (identified by `pg_backend_pid()`) wrote no log line in that window,
   and no log line named those seats.
5. In the same second, other backends' hold statements on neighbouring seats failed with
   `Value write after transaction start: doc ht (…) > read time (…)`. Their read times were
   about 300 ms older than recent commits.

`SUMMARY.md` of round 2 prints "2 issued again …" per configuration because each refusal's
detail appears twice in the result JSON. There was one refusal per configuration.

**What this establishes, and what it does not.** Under this workload's CPU saturation, on
YugabyteDB 2025.2.6 READ COMMITTED, a guarded `UPDATE` can behave as if it evaluated its
predicate against a state without a commit acknowledged to the same client before the
statement began, and silently match nothing. The miss is transient: the same statement
a moment later, in the same transaction, sees the hold. A statement that matches no rows
writes nothing, so no write-conflict check can turn the stale read into a retryable error;
that is consistent with the silence in the log. The engine's internal mechanism is **not**
established: the tserver's own logs at default verbosity say nothing, and finding it would
take engine-level tracing beyond this study. Two facts from dc11 fit this reading: the
refusals appear only in designs whose refusing statement filters on columns the hold just
wrote, and more often on three nodes.

**Rationale.**

- *Why not a harness bug:* the ledger tests pass; the per-occurrence evidence is complete
  and internally consistent; the re-issue proves the harness bound the right values.
- *Why run the matrix rather than keep investigating:* this failure is exactly what the
  owner asked the study to find ("I don't want to find out at payment that the seat is not
  available"). Measuring how often each design, engine and topology shows it is the study's
  job. Explaining the engine's internals is not.
- *Why a separate class:* an early rejection caused by theft, skew or a wrong expiry (S0, E0)
  is a design failure. A transient refusal is the engine refusing a hold it can see a moment
  later. Mixing them would make S1 on YugabyteDB look like S0. The class is still an INV-3
  violation, because the buyer was refused, and the cells still carry ❌.
- *Why S1r:* the owner needs to know what to do about it. "Retry a short confirmation once"
  is the smallest application-side defence and is safe by construction: the retried
  statement is still guarded, so it cannot sell an expired, released or stolen hold.
  S1 stays unchanged; §12 forbids fixing a design so it passes, and S1r is a new id that
  differs by one decision.
- *Why not an engine setting as a design:* none of the settings tried (B, C, E) removed it, and
  study designs are data-model and transaction decisions, not engine tuning.

**For the owner (not the executor):** the evidence bundle
(`results/devchecks/er01-*`, ESCALATIONS ER-01) is enough for an upstream report to
YugabyteDB, should you want one.
