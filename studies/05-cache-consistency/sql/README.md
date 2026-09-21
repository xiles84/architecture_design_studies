# The Study 05 SQL catalogue

Designs are readable SQL files, embedded into the benchmark binary with `go:embed`, so the
SQL that produced a result is provably the SQL sitting beside it. Named statements use the
repository format:

```sql
-- name: r_portal_person
-- params: person_id
-- Free-form documentation lines.
SELECT ...;
```

`params` names are bound by the harness, which is what lets one driver execute statements
whose signatures differ between designs.

## Layout

```
legacy/d3_legacy/          the LEGACY database model — Study 01's D3, frozen.
                           It is also the flattened-FK reference cell: the reference set
                           measures it with the cache backend set to `none`.
owned/d3_owned/            the OWNED database model: D3 + person.cache_version +
                           cache_outbox + atomic version changes + ordered reassignment.
reference/normalized_indexed/   rolldown reference: no copied charity_id on donation.
reference/rollup_trigger/       rollup reference: stored person aggregates, trigger-maintained.
reference/embedded_locked/      embedding reference: bounded newest-20 slice on person,
                                maintained by Study 01's D10 (concurrency-correct) trigger.
reference/y1_colocated/         placement reference (YugabyteDB only): child keyed by person.
reference/y2_noncolocated/      placement reference (YugabyteDB only): child keyed by id.
```

Each directory holds `schema.sql`, `indexes.sql`, `[triggers.sql]`, `queries.sql`,
`writes.sql` and `audit.sql`. A design's registry entry may point `queries.sql`,
`writes.sql` and `audit.sql` at another directory via `SQLDir`; that is how the two
placement cells reuse the flattened baseline's operations, so their pair differs by
exactly one decision. A design that has a `triggers.sql` gets it applied **after** the
bulk load, with its backfill, so the loader does not pay a per-row rebuild.

## The uniform read shape

Every design returns the same column list for `r_portal_person` and `r_portal_recent`:

| Statement | Columns |
|---|---|
| `r_portal_person` | `person_id, full_name, email, joined_at, charity_id, charity_name, charity_country, donation_count, donation_total_cents` |
| `r_portal_recent` | `donation_id, amount_cents, currency, donated_at, note` |
| `r_charity_recent` | `donation_id, amount_cents, donated_at, full_name` (the blended workload's uncacheable read) |

That is what lets one canonical Go encoder produce the content hash for every design, so
the hash compares payloads rather than JSON renderers. A design that returned a different
shape would be measuring its own formatting.

## The mutations

Five logical mutations, present in every design: `w_donation_insert`,
`w_donation_correct`, `w_donation_delete`, `w_person_update`, `w_donation_reassign`.
A sixth statement, `w_load_donations`, moves the staging table into the design's tables.
The owned model adds `w_version_bump`, `w_outbox_insert`, `w_version_cas` (optimistic) and
`w_lock_person` / `w_lock_persons2` (pessimistic, and ordered locking for the two-parent
reassignment).

## The audits

`a_person_aggregate`, `a_donation_ids` and `a_recent_ids` exist in every design and are
cross-checks against the Go oracle. `a_rolldown_drift`, `a_rollup_drift`,
`a_embedded_drift`, `a_version_monotonic` and `a_outbox_orphans` exist only where the
invariant they check exists, and each must return zero.
