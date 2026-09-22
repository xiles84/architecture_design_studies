---
discussion_id: 20260921-configuration-portal-mechanisms--deepseek-flash--2026-09-21
run_id: 20260921T1215Z-small
environment: host-zenbook-ux5406sa
author: deepseek-flash
author_kind: ai
author_version: "DeepSeek HIGH (session model id deepseek-flash; effort setting and tool identity not exposed)"
written_at: 2026-09-21
inputs_digest: ab0f6ef5e5500775
repo_commit: 5e15abb7b15e3236af9a6c91dcfcf0f6289bde81
companion_to: 20260921T1215Z-small--deepseek-flash--2026-09-21
---

# Discussion — configuration portal mechanisms, and one artefact I could not remove

Companion to the signed analysis of `20260921T1215Z-small`. This file holds the SQL-level and
plan-level reasoning behind that analysis's conclusions, plus the one measurement problem I found
and did not solve. It is written by the same agent that designed, implemented and measured the
study — see the provenance note at the end.

## 1. Why the lock beat the version check, in statements

Both designs perform the same logical operation: read a configuration key that carries a counter,
add one, write it back, bump the installation's revision — all in one transaction. They differ in
one statement.

`c1_optimistic_version` reads `(value, version)` and then:

```sql
UPDATE config_entry
SET value = $3, version = version + 1, updated_at = now()
WHERE installed_product_id = $1 AND key = $2 AND version = $4;
```

A zero-row result is not an error; it is the arbitration signal, and the harness reads the row
again and retries, bounded by `-retries` (default 8). `c2_pessimistic_lock` issues, before reading
anything:

```sql
SELECT current_revision FROM installed_product WHERE id = $1 FOR UPDATE;
```

and then writes unconditionally. The competing transaction blocks on that row lock and proceeds
when the first commits; at READ COMMITTED its subsequent read sees the committed value.

**Observation.** 16 writers on one key of one installation: `c2` 3 539 acknowledged at 886 ops/s,
`c1` 2 065 at 503 with 545 retries and 10 560 conflicts.

**Mechanism.** Every `c1` retry is a wasted round of work: read, compute, write, discover the
version moved, roll back. At 16 concurrent writers on one row the probability that the version is
still current when the guarded `UPDATE` runs is low, so most attempts are discarded. `c2` replaces
that discarded work with waiting, and waiting on a row lock in-process costs less than a full
transaction restart — on PostgreSQL, where the lock manager queues waiters fairly. 545 retries for
2 065 successes means roughly one attempt in five was thrown away; the measured ratio of 1.76x is
consistent with that.

**Alternative explanation I cannot exclude from this run.** `c1` records 10 560 "conflicts" but
only 545 *retries*, because the counter `conflicts` counts every zero-row guarded write and the
retry counter counts only ops that exhausted their loop. If the harness is charging more work to
`c1` than the measurement itself (a per-retry `rand` allocation, ledger locking), part of the 1.76x
is the harness. This is testable: the same race with `-retries 1` would show whether the gap is
proportional to retry count. **I did not run it, and the 1.76x should be read as an upper bound on
the lock's advantage until someone does.**

## 2. Why `x1` loses 93.7 % and why that is the whole point

`x1` uses the identical statements as `c2` with the lock removed; its write is

```sql
UPDATE config_entry SET value = $3, version = version + 1, updated_at = now()
WHERE installed_product_id = $1 AND key = $2;
```

Two transactions read 5, both write 6 and 6, both commit. Under READ COMMITTED nothing forbids
this: the second `UPDATE` re-reads the row but writes an explicit value, so it neither fails nor
notices. 3 416 acknowledged increments, counter 215, 3 201 lost updates.

The design of the *detector* matters as much as the design of the control, and it took two
attempts to get right. The first version compared the stored counter against the ledger's value for
that key — the value the last successful writer recorded. That compares a number with itself and
reported the control as a pass. The detector that works compares the counter against **the counter's
starting value plus the number of acknowledged increments**; each increment adds exactly one, so
the two must agree, and any gap is a lost update by construction. This is the general shape of a
reconciliation audit worth having: derive the expectation from the *operations the client saw
acknowledged*, never from a value the client last happened to write.

`x2`'s detector is the same idea applied to derived data. `x2`'s SQL is byte-identical to `n4`'s;
the harness calls `w_recompute_rollup` in a second transaction after the publication commits,
computing the count from a value read *inside* the publishing transaction. Concurrent publications
therefore each write a count that is already stale, and `a_rollup_mismatches` — which recomputes
the count and hash from `config_entry` and compares them with the stored columns — recorded 60
mismatches. Both controls fired; that is what licenses reading the other designs' clean audits as
evidence rather than as silence.

## 3. The rollup's price is per row, and that is the trigger's whole character

`n3_rollup_trigger` attaches one function per row:

```sql
CREATE TRIGGER config_entry_rollup
AFTER INSERT OR UPDATE OR DELETE ON config_entry
FOR EACH ROW EXECUTE FUNCTION refresh_ip_rollup();
```

and the function recomputes `count(*)` and the content hash over *all* of that installation's rows.
`n4_rollup_app` issues the same recomputation once, from the application, inside the publishing
transaction.

**Observation.** `r04_list_installations`: 867 (n1) → 13 197 (n3) / 13 693 (n4) ops/s. Writes:
`w05_replace_all` 572 (n1) → **83** (n3) / 593 (n4); `w04_publish_batch` 1 023 → 403 / 841. p99 on
`w05` rises from 2.5 ms to 42 ms.

**Mechanism.** A whole-configuration replacement deletes 60 rows and inserts 60. Under `n3` that is
120 trigger firings, each recomputing an aggregate over the surviving 60 rows — O(keys²) work for
one publication, which is precisely the "unbounded write amplification" study 01 measured on its
embedded design. `n4` does the same aggregate work exactly twice (once per publication, after the
statement), so its cost is O(keys) and it stays within 4 % of the reference on `w05`.

**What this does not say.** The trigger's advantage is not speed and should not be dismissed by a
speed comparison: it cannot be forgotten by a write path, and `x2` is the demonstration of what
happens when the application's maintenance is not in the same transaction. The correct framing is
that `n3` buys a structural guarantee at a measured 6.9x on the most expensive write, and `n4`
buys the speed at the cost of every write path having to remember.

## 4. The ordering artefact: 65 % spread among identical SQL, and I did not control it

This is the most important thing in the run and it is a problem with the run.

`r01_effective_config`, `r02_read_key` and `r03_revision_check` have **byte-identical SQL** in the
eight row-per-key designs. Their measured throughputs should therefore be indistinguishable.
They are not:

| Read | slowest (design) | fastest (design) | ratio |
|---|---|---|---|
| `r01_effective_config` | 9 148 (`n1`) | 15 141 (`x1`) | 1.66x |
| `r02_read_key` | 16 308 (`n1`) | 26 801 (`x1`) | 1.64x |
| `r03_revision_check` | 25 114 (`n1`) | 40 010 (`x2`) | 1.59x |

`n1` is slowest in all three and `x1` fastest in two. That pattern is not a design effect: `n1` and
`x1` differ only in a write statement that the read phase never issues.

**Mechanism.** Every cell in the matrix runs against the **same** `pg-single` container, in one
sequence, and each cell begins by dropping and recreating its schema (`DropSchema` then
`ApplySchema`). Dropping a table that a previous cell's write phase filled with dead tuples leaves
the catalog and the shared buffers in a state that differs from cell to cell, and `n1` was the
eighth of ten cells in the recorded order. The database is not the same instrument at cell 1 and
cell 10.

**Why it was not caught earlier.** The gate checks answers, not instrument stability, and it should
not: comparing two queries to each other is exactly what methodology 5 forbids. The correct control
is a *repeated* design in the matrix — the same design run first and last — and this matrix has
none.

**Experiment that would settle it.** Run `n1` twice in one matrix, once first and once last, with
three trials each, and report the difference. If the first/last gap is of the same magnitude as the
`n1`–`x1` gap, then the whole read table must be re-measured with per-cell databases. That is a
cheap experiment and it is the first item in the analysis's "what I would measure next". Until it
runs, the working rule is that **read differences below about 1.7x in digest `ab0f6ef5e5500775` are
not attributable**, which is why the analysis claims only the ~15x rollup effect and treats the
`n0`-faster-than-`n1` result as a hypothesis.

## 5. Storage, and why the comparison is clean where the timings are not

`n1` 1 032 192 B of `config_entry` against `d1` 98 304 B of `config_document`, for the same 3 600
logical entries. This measurement is a `pg_total_relation_size` call, not a timing, so the ordering
artefact does not touch it; it is the most solid number in the run. `load_entry` (688 128 B) is
excluded deliberately — it is staging that exists for every design and would otherwise dominate the
comparison. The 10.5x is a row-overhead and TOAST story, and the interesting follow-up (not run) is
how it moves at 500 entries, where the document design's whole-row rewrite should start costing
more than it saves.

## Provenance and independence

Study 04 was planned, implemented, measured, validated and analysed by **one** model, DeepSeek HIGH
(session model id `deepseek-flash`; effort setting and tool identity not exposed). There has been no
independent review by a different model or by a person, and the author of this document is also the
author of the designs being compared — including the choices of what to measure, which the
ordering artefact in §4 shows can hide a problem. A second analyst is invited. The digest
`ab0f6ef5e5500775` identifies exactly the data any reply should address, and the register of
disagreements is deliberately left open: this document is signed, final, and never edited.
