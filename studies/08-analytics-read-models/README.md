# Study 08 — analytical, search and read-model copies

**Question:** for top-N and dashboard workloads, what do four read-model strategies cost and return
— base OLTP, rollup/materialized view, search copy, and analytical copy?

**Status:** protocol only; **no measurement yet**. The decision-complete handoff is
[`HANDOFF.md`](HANDOFF.md) (tag `repo/analytics-read-model-handoff-v1`).

**Strategies:** PostgreSQL 17.11 base and materialized views, a PostgreSQL full-text/trigram search
copy, and a ClickHouse analytical copy (pinned line; exact tag and digest recorded by the harness).

**Measures per strategy:** query gain, refresh/maintenance, staleness, storage, correctness — with a
stale-copy negative control and a refresh-boundary test.

**Boundary:** v0 is single-node and PostgreSQL-based; a native search engine is a separately labelled
extension, and cluster topology belongs to Study 07.
