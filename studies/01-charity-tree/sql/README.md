# The SQL catalogue

One directory per design. These files are **both the executable definition of a design and
its documentation** — they are compiled into the benchmark binary with `go:embed`, so a
result can never have been produced by SQL that differs from the SQL sitting here.

```
<design>/
├── schema.sql     tables (and, in D7, the physical key layout)
├── indexes.sql    everything built after the data lands
├── queries.sql    the 12 read questions
├── writes.sql     insert / update / delete, plus D5's locking variants
└── triggers.sql   only where a design maintains something automatically (D4, D9)
```

## Annotated-statement format

`queries.sql` and `writes.sql` hold **named** statements:

```sql
-- name: q03_top_donor_charity
-- params: charity_id
-- Free-form documentation. Everything before the first line of SQL is metadata;
-- comments after the body starts are kept as part of the statement.
SELECT ...;
```

- `name:` is the id used in results, reports and plan dumps.
- `params:` lists parameter names in `$1..$n` order, or `none`. The harness binds them
  **by name**, which is what lets one benchmark driver execute statements whose signatures
  differ between designs — D1's insert takes six parameters and D3's takes seven, and
  neither needs a special case in the driver.

Parameter names the harness knows: `charity_id`, `person_id`, `donation_id`, and for
writes `amount_cents`, `currency`, `donated_at`, `note`, `version`.

## Rules for writing a design

**Every design must answer all twelve questions correctly.** That is enforced, not
assumed: `-cmd verify` checks each answer against a value computed independently in Go
from the generated dataset, and a failure aborts the cell before any timing is recorded.

**Write each query the best way the design allows.** D6's queries exploit the fact that
its JSONB array is stored in chronological order, so "first and last donation" is two
array subscripts rather than an unnest. Where a query still has to scan everything, that
is a genuine property of the design — but it must not be a strawman.

**Prefer designs that differ from an existing one by exactly one decision.** D1 and D2
share byte-identical `queries.sql`; D3 and D8 differ only in three `REFERENCES` clauses;
D3 and D7 differ only in a `PRIMARY KEY` definition. A design that changes three things at
once produces a number nobody can attribute.

**Comment the *why*, not the syntax.** These files are read as prose by people deciding
how to model their own tree.
