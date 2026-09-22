# Study 06 — native storage models

**Question:** do the native document, wide-column and authoritative key-value models behave as the
book's families predict, when the configuration-domain semantics, operations, resources, correctness
gates and topology controls match the relational baseline?

**Status:** protocol only; **no measurement yet**. The decision-complete handoff is
[`HANDOFF.md`](HANDOFF.md) (tag `study-06/v0-handoff`).

**Technologies (pinned):** PostgreSQL 17.11 (baseline), MongoDB 8.0 (document), ScyllaDB 6.2
(wide-column), Valkey 8.1 (authoritative key-value), plus a PostgreSQL ledger for the history family.
Exact patch tags and image digests are recorded by the harness task before any measurement.

**Why it exists:** the evidence registry records `v2-gap-06-native-datastore-families` — JSONB is not
a native document store and Redis was measured only as a cache. This study is the first native
cross-model comparison.

**Boundary:** v0 is single-node only; cluster and RF variants belong to Study 07. Each comparison is
within-run against the baseline; no cross-technology ranking.
