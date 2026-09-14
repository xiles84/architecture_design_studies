# SQL catalogue — study 03

One directory per design. Every design has the same five files:

| File | Contents |
|---|---|
| `schema.sql` | tables and constraints, applied before the bulk load; its header explains the design |
| `indexes.sql` | secondary indexes, built after the load |
| `queries.sql` | the read catalogue — `q01`…`q06`, identical names and result shapes in every design |
| `writes.sql` | every statement the hold, checkout, sweeper, refund and publishing flows issue |
| `audit.sql` | how to recompute the truth from this design's tables, for the audit and the ledger |

Named statements use the repository's catalogue format (`-- name:`, `-- params:`). The Go
harness decides **which** statements run, in **which** transaction, at **which** isolation
level (`harness/designs.go`); the SQL files decide **what** each statement does.

## Designs that differ by one decision — check with `diff -r`

```bash
diff -r sql/s0_check_then_hold_rc sql/s4_check_then_hold_serializable   # first lines only: the decision is the isolation level
diff -r sql/s1_conditional_update sql/k0_naive_confirm                  # w_confirm_seats (and headers)
diff -r sql/s1_conditional_update sql/l3_section_sharded                # the primary key, one index (and headers)
diff -r sql/s1_conditional_update sql/s2_lock_then_update               # w_lock_seats added (and headers)
diff -r sql/s2_lock_then_update sql/s3_lock_nowait                      # FOR UPDATE NOWAIT
diff -r sql/s1_conditional_update sql/e0_app_clock_expiry               # now() -> $app_now in the logic
diff -r sql/s1_conditional_update sql/e1_sweeper_expiry                 # the expiry predicates
diff -r sql/s1_conditional_update sql/k1_payment_window                 # checkout_started_at, w_begin_checkout
```

E2, L1 and L2 change the tables and are read as whole designs.

## Conventions

- **Every statement that grants, extends, releases or sells returns `now() AS db_now`.** The
  harness's ledger judges theft, honored holds and late sales by the database clock, even in
  E0, whose logic uses the application's clock.
- **Blocks are `seat_ids INT[]`**, always sorted ascending by the harness, and always inside
  one section: every seat predicate carries `event_id` and `section_no`.
- **All-or-nothing is the caller's job.** A multi-row statement returns the seats it took;
  fewer than requested means roll back.
- **Parameters in `INSERT … SELECT` and `unnest` are cast explicitly** (study 02's rule).
- **Index column order is chosen for YugabyteDB's defaults**, which hash-shard the first
  column; range-scanned expiry columns are written `ASC`.
- The loader runs any `w_load_*` statement once, after the bulk copy (L1 writes its sold
  claims' `'infinity'` marker that way). The explain phase skips them.
