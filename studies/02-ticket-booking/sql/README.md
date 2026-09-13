# SQL catalogue — study 02

One directory per design. Every design has the same five files:

| File | Contents |
|---|---|
| `schema.sql` | tables and constraints, applied before the bulk load |
| `indexes.sql` | secondary indexes, built after the load |
| `queries.sql` | the read catalogue — `q01`…`q05`, identical names and result shapes in every design |
| `writes.sql` | every statement the booking, cancellation, publishing and hold flows issue |
| `audit.sql` | how to recompute the truth from this design's tables, used by the overbooking audit |

Named statements use the repository's catalogue format:

```sql
-- name: w_take_seat
-- params: event_id
-- Free-form documentation lines.
UPDATE event SET seats_sold = seats_sold + 1 WHERE event_id = $1 AND seats_sold < capacity
RETURNING seats_sold;
```

The Go harness decides **which** statements run, in **which** transaction, at **which**
isolation level (`harness/designs.go` → `Strategy`, `Isolation`, `LockEvent`). The SQL
files decide **what** each statement does. Where two designs differ by one decision, the
difference is usually visible as a diff of these files:

```bash
diff -r sql/c1_count_naive sql/c2_count_serializable   # identical: the decision is the isolation level
diff -r sql/h0_hold_naive_confirm sql/h1_hold_checked_confirm   # one WHERE clause
```

## Portability between PostgreSQL and YugabyteDB

The same files run on both engines. Two rules make that work:

- **Index column order is chosen for YugabyteDB's defaults.** YugabyteDB hash-shards the
  first column of a primary key or index and range-orders the rest, so
  `(event_id, seat_no)` means "all of one event's seats together, in seat order" on both
  engines. A column that must be range-scanned on its own (the hold expiry) is written
  `(expires_at ASC)`, which PostgreSQL accepts and YugabyteDB reads as range sharding.
- **Parameters in `INSERT ... SELECT` and `generate_series` are cast explicitly.** Without
  a cast, a parameter in a select list has no type to infer from.
