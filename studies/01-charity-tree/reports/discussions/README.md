# Detailed study discussions

The [final analyses](../analyses/) give concise conclusions. These signed companions
keep the SQL explanations, controlled design comparisons and responses between analysts.

| Topic | Discussion |
|---|---|
| D2/D3: copied key, indexes, SQL, heap fetches and query weights | [Mechanism comparison](d2-d3-mechanisms--gpt-6--2026-09-13.md) |
| D4/D5/D16: full totals, lock-first application writes and sums only | [Parent contention](parent-contention--gpt-6--2026-09-14.md) |
| D3/D6: longer histories, mutation cycles and memory configurations | [Embedding and growth](embedding-growth--gpt-6--2026-09-14.md) |
| YugabyteDB: FK exceptions, D9/D10, equal budgets and a local node stop | [Distributed controls](yugabyte-controls--gpt-6--2026-09-14.md) |

Each discussion names its author, input digests and source revision, and links to its
final analysis. To extend or disagree with one, write a separate signed response using
the [discussion template](../../../../docs/templates/DISCUSSION.md), cite its
`discussion_id` or `analysis_id`, and explain which observation or assumption differs.
Do not edit another analyst's file or imply that a model participated in an exchange
it has not written. Earlier published analyses remain preserved for their own inputs.
