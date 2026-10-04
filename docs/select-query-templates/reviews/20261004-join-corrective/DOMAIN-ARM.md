# Study 01 domain arm — fixture, projection and rank proxy (review R2 / C2)

LOW `deepseek-flash` / opencode CLI, 2026-10-04, task
`task-20261002T131729Z-select-join-cte`. This is the mandatory Study 01 domain
binding; the synthetic customers/facts fixture is retained unchanged as the
microfixture arm.

## What is measured

The family's reusable catalogue is written against the projection
`customers(id, name)` and `facts(id, customer_id, rank, amount, payload)`. The
domain arm loads Study 01's real entities into `selectjoin_domain` and exposes
that projection through two views, so **the measured SQL is byte-identical
between arms** and only the data's origin differs:

```sql
person  (person_id, charity_id -> charity, full_name)
donation(donation_id, person_id -> person, rank, amount_cents, note)

CREATE VIEW selectjoin_domain.customers AS
  SELECT p.person_id AS id, ch.name AS name
  FROM person p JOIN charity ch ON ch.charity_id = p.charity_id;

CREATE VIEW selectjoin_domain.facts AS
  SELECT d.donation_id AS id, d.person_id AS customer_id,
         d.rank AS rank, d.amount_cents AS amount, d.note AS payload
  FROM donation d;
```

- **Projection:** `(person_id, charity_name, donation_id, amount_cents)`. The
  parent is a person, the child a donation, and the person's charity is the
  relationship the Charity donations case turns on. `AsJoinFixture` in
  `selectbench/joindomain.go` is the single place this binding is defined.
- **Generated keys:** `person_id` 1..P assigned round-robin over `charity_id`
  1..C; `donation_id` 1..N assigned round-robin over `person_id`;
  `donation.rank` is a seeded Fisher–Yates permutation of 1..N. The permutation
  is the same generator the synthetic arm uses, so every threshold selects an
  exact count with no sampling error.
- **Rank proxy:** Study 01's real selective predicate is a date or amount range,
  which cannot produce the reviewed exact thresholds. `rank <= threshold` is the
  proxy; the label travels with the result and the report states it.
- **Foreign keys:** `NOT NULL`, matching Study 01's schema. The optional-match
  contract is exercised by persons with no qualifying donation; the synthetic
  microfixture keeps the orphan-NULL-FK case as its own evidence.
- **Indexes:** `donation_rank_idx(rank)`, `donation_person_idx(person_id)`, both
  primary keys. No per-spelling index, so a spelling pair shares one physical arm.

## Scale and matrix

`cte-domain` = 50 charities, 2,000 persons, 200,000 donations — the same shape as
the approved `cte` microfixture. `cte-domain-small` (20/200/20,000) is used for
gates and smoke checks only. The matrix is unchanged: the 60-cell spelling grid
(10 spellings x 6 exact thresholds) plus the 30-cell reuse ladder (3 policies x
5 reference counts x 2 consumer selections), 5 trials per cell, seeded
interleaving and the actual-first-cell repeat.

## Live validation (no matrix yet)

| run | tag | result |
|---|---|---|
| `select-join-domain-feas-pg-20261004` | `run/study-01/select-join-domain-feas-pg-20261004` | retained: exposed two probes still bound to the synthetic schema (threshold counts `-1`, rank index not seen) |
| `select-join-domain2-feas-pg-20261004` | `run/study-01/select-join-domain2-feas-pg-20261004` | PG gate passes; threshold counts exact (2/20/200/2,000/10,000/18,000); rank index used; `x01` and `x02` controls fire; MATERIALIZED shows `CTE Scan` |

Unit tests (`selectbench/joindomain_test.go`) cover the projection binding,
determinism, the rank permutation and the load/view order; `go vet` and
`go test` pass in Podman.

## Open

The corrected single-node matrix (3 whole runs per engine, 5 trials) has not
been run; it is the next step, under the new run ids and the calibrated 1-CPU
client. The 3-node/cluster and blended phases remain C4 and C5.
