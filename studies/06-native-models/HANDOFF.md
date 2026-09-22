# Execution Handoff — Study 06, native storage models

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-native-major-model-protocol` (HIGH planning).
**Required tag on integration:** `study-06/v0-handoff`.

## 1. The question

Do the *native* major storage models — document, wide-column and authoritative key-value — behave as
the book's families predict, when the configuration-domain semantics, operations, resources,
correctness gates and topology controls are identical to the relational baseline?

**This is not** an extension of JSONB-as-a-document result or of Redis-as-a-cache. The registry
records that boundary as `v2-gap-06-native-datastore-families`: JSONB is not evidence for a native
document database, and Redis was measured only as a cache.

## 2. Technology selection and pinning

| Role | Technology | Pin | Why |
|---|---|---|---|
| Relational baseline | PostgreSQL | `docker.io/library/postgres:17.11` (already in `versions.env`) | The corpus baseline; Study 04's configuration-portal schema and operations are reused unchanged. |
| Native document | MongoDB | `docker.io/library/mongo:8.0` | A document store whose native unit is the embedded child; the closest native test of the embedded-documents family. |
| Native wide-column | ScyllaDB | `docker.io/scylladb/scylla:6.2` | Partition + clustering keys are the native form of "one parent, many children", and the native test of the rolldown/derived-state boundary. |
| Native key-value (authoritative) | Valkey | `docker.io/valkey/valkey:8.1` | A key-value store used as the *system of record*, not a cache; the native test of the external-read-copies family's boundary. |
| Append-only history (native) | PostgreSQL ledger | `postgres:17.11` log schema | History has no separate native engine in scope; the ledger is included so the family is represented without claiming a new technology. |

Pinning rule: each technology is pinned to the exact tags above. If a tag has moved or does not
exist, the harness task resolves and records the **exact tag plus image digest** in `versions.env`
and every run manifest before any measurement; a moved tag is a run blocker, not a silent upgrade.
MongoDB/ScyllaDB/Valkey versions are chosen as the current stable lines and are confirmed at harness
build; the protocol fixes the *line*, the harness fixes the patch and digest.

## 3. Identical semantics

The Study 04 configuration-portal semantics are reused verbatim: one product has many configuration
entries; the operations are read-whole-configuration, point-read one entry, replace-whole-
configuration, move one entry between products, and read a derived overview. Each technology
implements the same logical data and the same five operations; the SQL/query catalogue keeps the
same statement names and result shapes where the technology allows, and a documented native analogue
where it does not.

## 4. Correctness gate

Study 04's invariants INV-1…INV-13 apply unchanged. Each technology adds one native check:
MongoDB a document schema-validation rule, ScyllaDB a lightweight-transaction or
read-before-write check, Valkey atomicity per key. Every design has at least one negative control
that must be seen to violate the invariant before timing is reported.

## 5. Resources and topology controls

- Per-node budget identical to the corpus: 2 CPU / 3 GiB per database node, 2 CPU / 2 GiB client;
  the aggregate budget and a separately labelled equal-total arm are recorded.
- v0 runs each technology single-node only. Cluster and RF variants are **deferred to Study 07**
  (topology) so this study changes one thing — the storage model — at a time.
- Client readers, writers and connection pools are calculated and recorded; load-generator capacity
  is checked against the intended demand.

## 6. Deliverables and separately claimable execution tasks

1. **Native-model harness and environment** — Containerfiles, `versions.env` pins with digests,
   dataset loader and native correctness gate for MongoDB, ScyllaDB and Valkey; probe scripts.
2. **Document vs relational** — the five operations on MongoDB and the PostgreSQL baseline.
3. **Wide-column vs relational** — the five operations on ScyllaDB and the baseline.
4. **Authoritative key-value vs relational** — the five operations on Valkey-as-record and the
   baseline; explicitly not a cache comparison.

Each task takes the benchmark lock, runs one matrix at a time, and produces a generated report plus a
signed analysis under `studies/06-native-models/`. No task publishes a cross-technology ranking; each
compares its technology to the baseline within one run.

## 7. Acceptance criteria for the study report

- Every technology resolves to a pinned tag and recorded digest, in `versions.env` and the manifest.
- Identical logical data and operations across technologies, with native analogues documented.
- Correctness gate passes and every negative control fires; failed cells named, never winners.
- Per-node and equal-total arms labelled and separate; single-node only in v0.
- Each comparison is within-run; no cross-technology pooled ranking.
- Every number resolved to a v2 evidence claim before it enters the book.
