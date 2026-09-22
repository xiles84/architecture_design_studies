# Execution Handoff — Study 04, configuration portal

| Field | Value |
|---|---|
| Handoff ID / revision | `EH-04` revision 1 |
| Planner | **DeepSeek HIGH**; model id exposed by the session: `deepseek-flash`; effort setting **not exposed**; tool identity **not exposed**; planned 2026-09-21 |
| Starting source / main revision | local `main` = `df13f2a6baeded28bf5d1216b6ec6d622e259850e` (`study-01/v3-reviewed`); worktree `.worktrees/study04-configuration-portal`; branch `study-04/configuration-portal` |
| Checkpoint tag | `study-04/v0-handoff` (this document). Never moved, never reused. |
| Status | ready — amendment AM-01 (ER-01, the WSL→Podman bridge) is the first execution action |
| Next setting | **HIGH — implement this committed handoff using the same DeepSeek HIGH model in WSL** |
| Progress / escalations | `PROGRESS.md` (task-owned), `ESCALATIONS.md` (task-owned) |
| Environment | `host-zenbook-ux5406sa`, live-reconfirmed 2026-09-21 (see §8) |

## 0. How to use this document

This is the specification of record for Study 04. It is a *plan plus an execution map*: every
decision that shapes a number is taken here; everything the executor may choose freely is listed
in §19. It does not replace the methodology ([`docs/methodology.md`](../../docs/methodology.md)),
the repository rules ([`AGENTS.md`](../../AGENTS.md)) or the study's own `README.md`.

The owner's instruction for this task is unusual and is recorded plainly: **one DeepSeek HIGH
agent performs every phase** — planning, implementation, measurement, validation, analysis and
integration. There is therefore no LOW executor and no model switch. The HIGH/LOW workflow is
kept as a *provenance checkpoint* (this document, its tags and its amendments) rather than as a
division of labour. The absence of an independent model's review is a limitation of the study's
conclusions and is stated in the final analysis.

`docs/environments/host-zenbook-ux5406sa.md` bounds what every number here can mean: 8
heterogeneous cores shared by client and server, no real network, laptop thermals, and named
resource conditions that must never be pooled.

## 1. The request, and how it was interpreted

The owner asked for a new study of a **configuration portal**: a product-configuration service
used by other products, where every creation, modification, publication and retrieval of
configuration goes through the portal and its database.

The study must measure how design choices move performance, storage, contention, correctness and
maintenance under: data-design alternatives (row-per-key, embedding, rollup, rolldown, hybrid,
immutable snapshots), configuration cardinality up to a mandatory maximum of 60 entries, value
and serialized-document size, update cadence from daily to secondly, update contention,
optimistic versus pessimistic concurrency, single-node versus three-node operation, equal-total
resource deployment, and colocated versus non-colocated distributed data.

Interpretation decisions taken here, because the request leaves them open:

- **One shared logical dataset for every design and engine**, fixed seed, bounded skew (§5).
- **Cadence is a rate, not a wait.** "Once per day" means a calculated fleet demand of
  *installed products / day*, not a day of wall-clock time (§6).
- **The value 60 is a first-class tier** and is never averaged away (§5).
- **Every design is measured with the same operations**, and the operations are the six reads and
  six writes the owner listed (§7).

## 2. Study identity, glossary and product semantics

**Study id:** `04-configuration-portal` — directory `studies/04-configuration-portal/`.
**Benchmark image:** `localhost/configbench:1` (`study.env`), so rebuilding this study's client can
never change another study's matrix.

### 2.1 Glossary (also in `README.md`)

| Term | Meaning |
|---|---|
| Product definition | a product type or catalogue entry. Owns no configuration. |
| Installed product | one deployed instance of a product definition, inside an environment, optionally attributed to a business unit. **Configuration belongs to the installed product**, not the definition. |
| Environment | the environment containing an installed product (e.g. production, staging, development). |
| Business unit | an optional ownership or scope dimension. Several installed products of the same definition may exist in one environment — one per business unit. |
| Configuration entry | one named configuration key and its value. |
| Effective configuration | the complete current configuration returned to an installed product. |
| Configuration revision | the monotonic version of an installed product's effective configuration. |
| Publication | making a new configuration revision current. |

### 2.2 Product semantics — identical for every design

These are contract, not implementation. A design that cannot meet them is a design that fails,
not a design that defines a different product.

1. Configuration is scoped to an **installed product**. Two installed products of the same
   definition never share an entry.
2. A complete-configuration read returns **one coherent committed revision**. It must never
   combine keys from two revisions.
3. A publication **atomically** makes a batch of changes current, bumping the revision by exactly
   one.
4. Revisions increase monotonically per installed product and never repeat.
5. An **acknowledged** update is durable: once the portal has told a caller the change is
   published, no later read may fail to see it.
6. Installed-product **metadata** (display name, flags) can be updated independently of
   configuration without publishing a revision.
7. "Check whether a newer revision exists" must be answerable **without** reading the complete
   configuration.

## 3. Design catalogue — 18 designs

Every design is defined by `sql/<id>/` and registered in `harness/designs.go`. All eighteen share
one logical schema shape (product definitions, installed products, environments, business units,
configuration entries) and differ only in the mechanism named in the "isolates" column. Where a
design changes more than one mechanism it is split, or the narrowed claim is stated in the table's
Notes column and in the study `README.md`.

### 3.1 Normalized family

| ID | Schema | Isolates | Tradeoff expected | Correctness risks |
|---|---|---|---|---|
| `n0_rows_unindexed` | `config_entry(ip_id, key, value, updated_rev)` PK `(ip_id,key)`; **no** secondary index on `key` | the index control: what the R5 search index costs and buys | slower R5, cheaper W4/W2 | none specific |
| `n1_rows_indexed` | as `n0` **plus** `config_entry_by_key(key, ip_id)`; `installed_product` carries `current_revision` | **the reference** — row per configuration, indexed | — (this is the floor) | none specific |
| `n2_rows_rolldown` | `n1` **plus** `product_definition_id`, `environment_id`, `business_unit_id` copied onto every `config_entry` row, indexed `(product_definition_id, key)` | rolldown: parent information copied onto children | R5 without a join; W2/W4 write a wider row, and moving an installed product is a multi-row rewrite | copied columns can drift from their parents (INV-8) |
| `n3_rollup_trigger` | `n1` **plus** `ip_rollup(config_count, current_revision, last_modified, content_hash)` maintained by `AFTER` triggers on the writing tables | parent rollups maintained by the database | R4 free of aggregation; every write pays trigger work and a hot row | rollup can drift (INV-9) |
| `n4_rollup_app` | `n1` **plus** the same `ip_rollup` columns, maintained **by the harness inside the publishing transaction** | rollup maintenance by the application (pairs with `n3`) | same read gain, different write shape and failure surface | drift if a write path forgets it (INV-9) |

### 3.2 Document family

| ID | Schema | Isolates | Tradeoff expected | Correctness risks |
|---|---|---|---|---|
| `d1_doc_row` | `config_document(ip_id, revision, doc jsonb, content_hash, bytes)`, one row per ip; publication **replaces the whole document** | a separate one-to-one document row | one row per read; every publication rewrites the whole document, so write cost grows with cardinality | lost update if two publications race (INV-3) |
| `d2_doc_on_parent` | the document lives **as a column on `installed_product` itself** | embedding the document in the parent row | one row touched; a metadata-only update (W6) now rewrites the configuration too | same, plus W6 amplification |
| `d3_doc_sections` | one document row **per configuration section** (`config_section_document(ip_id, section, doc, ...)`), plus a per-ip revision row | namespace/section-sharded documents | publication rewrites only touched sections: bounded write amplification as cardinality grows | a complete read must assemble all sections at one revision (INV-2) |
| `d4_doc_jsonb_path` | as `d1` but a single-key change is an **`jsonb_set` on the existing row** | whole-document replacement versus JSON path modification | smaller WAL and MVCC churn per update, but an update now rewrites a TOASTed value anyway at large sizes | same as `d1` |

### 3.3 Hybrid and immutable families

| ID | Schema | Isolates | Tradeoff expected | Correctness risks |
|---|---|---|---|---|
| `h1_rows_plus_readview` | `n1`'s rows remain the source of truth; a materialized read representation (`config_readview(ip_id, revision, doc)`) is maintained in the **same transaction** as the write | normalized source **plus** a materialized read representation | fast complete reads without giving up row-per-key maintenance | two representations must agree after every write path (INV-7, INV-10) |
| `s1_snapshot_pointer` | `config_snapshot(ip_id, revision, doc, content_hash, bytes)` is append-only; `installed_product.current_revision` is the **atomic pointer** | immutable per-revision snapshots with an atomic current-revision pointer | publication is one row insert plus one pointer update; readers read the pointer then its snapshot; history is free | a reader that reads the pointer and then the snapshot must get *that* snapshot (INV-2) |
| `s2_append_history` | append-only `config_history(ip_id, key, value, revision, created_at)` plus a materialized current-state table | row-per-key immutable revisions with materialized current state | full history and auditability; two writes per change, and the current state can drift | drift between history and current state (INV-7) |

### 3.4 Placement family (YugabyteDB only)

| ID | Schema | Isolates |
|---|---|---|
| `y1_colocated` | `n1`'s exact schema declared `WITH (colocation = true)` — an installed product's rows share one tablet | physical data colocation |
| `y2_noncolocated` | `n1`'s exact schema with default (hash-sharded, distributed) placement | the non-colocated half of the pair |

`y1` and `y2` differ **only** in placement. `y2` must reproduce `n1` within the noise floor; the
report carries that as an internal consistency check (same logical data, same operations, same
replication factor, same resource budget — `docs/methodology.md` §6a).

### 3.5 Concurrency family

`c1`, `c2` and `x1` all use `n1`'s schema except where the table says otherwise. The *strategy* is
a decision in the harness (`harness/designs.go`), and the SQL decides what each statement does —
exactly as study 03 does it.

| ID | Strategy | Isolates | Notes |
|---|---|---|---|
| `c1_optimistic_version` | `config_entry.version` is read, then the publication is a guarded `UPDATE ... WHERE version = $expected`; on a zero-row result the harness re-reads and retries, bounded by the deadline | optimistic version-checked update with bounded retries | the row needs a `version` column, so `x1` is a separate table definition |
| `c2_pessimistic_lock` | `SELECT ... FOR UPDATE` on the installed-product row before reading or writing any entry | lock-before-update | same invariant, different arbitration |
| `x1_lost_update_control` | **negative control**: read the current value, then write the new value back — no version check, no lock, no `WHERE` on the old value | the guard itself | *must* be seen to lose one of two concurrent updates (INV-3) |
| `x2_rollup_drift_control` | **negative control**: `n4`'s rollup maintained **outside** the publishing transaction, in a second transaction after commit | atomic versus non-atomic application-maintained derived data | *must* be seen to drift under concurrency (INV-9, INV-10) |

Controls are measured like every other design and their speed is **never** reported as a plain
number (`docs/methodology.md` §5a). The report states, for every topology, whether each control
fired.

### 3.6 Considered and deliberately excluded

| Candidate | Decision | Reason |
|---|---|---|
| Foreign-key enforcement control | excluded | Study 01's `D8` already isolated the cost of enforcing referential integrity on the same engine pair. Here every design keeps the same FK shape, so a no-FK variant would move one mechanism already measured in this repository and add a cell that answers nothing new. |
| Configuration-section *partitioning* at the row level | excluded | `d3_doc_sections` already isolates the section axis at the document level. A partitioned-row variant would change the section unit *and* the physical partitioning at once, producing a number nobody can attribute. |
| Whole-document replacement vs JSON path modification | **included** as `d4_doc_jsonb_path` | it isolates exactly one decision (how the update is expressed) against `d1`. |
| Row-per-key immutable revisions | **included** (narrowed) as `s2_append_history` | the "row per key per revision" shape is carried by the history table; a *separate* design that also dropped the materialized current state would change two mechanisms. |
| Document-per-revision snapshots | **included** as `s1_snapshot_pointer` | identical to the required "immutable snapshots with an atomic current-revision pointer". |
| Append-only history with materialized current state | **included** as `s2_append_history` | see above. |

### 3.7 Controlled pairs (the report's `pairs` table)

| Pair | The one decision it isolates |
|---|---|
| `n1` ↔ `n0` | the secondary search index |
| `n1` ↔ `n2` | rolldown of parent keys onto configuration rows |
| `n3` ↔ `n4` | trigger-maintained versus application-maintained rollup |
| `d1` ↔ `d2` | separate one-to-one document row versus document embedded on the parent |
| `d1` ↔ `d4` | whole-document replacement versus JSON path modification |
| `d2` ↔ `d3` | one document per installed product versus per-section sharding |
| `n1` ↔ `d1` | row-per-configuration versus a single document |
| `n1` ↔ `h1` | source rows alone versus source rows plus a materialized read representation |
| `d1` ↔ `s1` | mutable document versus immutable snapshot with an atomic pointer |
| `n1` ↔ `c1` | no guard versus optimistic guard (only meaningful under contention) |
| `c1` ↔ `c2` | optimistic version check versus lock-before-update |
| `c1` ↔ `x1` | the guard versus its absence (the control) |
| `n4` ↔ `x2` | atomic versus non-atomic application-maintained derived data (the control) |
| `y1` ↔ `y2` | colocation (YugabyteDB only) |

## 4. Correctness contract

### 4.1 Invariants INV-1 … INV-13

Every invariant is verified **in Go from the generated dataset** — never by comparing one query to
another (`docs/methodology.md` §5). A failing gate aborts its cell and the cell reports no timing.

| ID | Invariant | Verified by |
|---|---|---|
| INV-1 | one current value per (installed product, key) pair | Go rebuilds the expected current map from the generated change log and compares it to a full read, key by key |
| INV-2 | a complete read contains keys from exactly one revision | the read returns its revision; Go recomputes the expected key/value map for **that** revision from the log and compares |
| INV-3 | no acknowledged update is lost | Go holds the set of acknowledged publications and requires every one to be present in the final effective configuration, and the revision to have advanced by exactly the number of acknowledged publications on that product |
| INV-4 | concurrent changes to **different** keys are preserved | injection: N writers each own a distinct key; every one must survive |
| INV-5 | revisions increase monotonically | the harness records the revision sequence per product and rejects any non-increase or repeat |
| INV-6 | added keys appear after publication | acknowledged `w02_add_key` must be visible in the next complete read |
| INV-7 | deleted keys do not remain in embedded or materialized representations | the read document, the section shards, the readview and the materialized current state are all checked to contain no deleted key |
| INV-8 | rolldown fields match their authoritative parents | `n2`'s copied columns are joined back to their parents and compared |
| INV-9 | rollup count, revision, timestamp and hash match independently derived values | Go recomputes count, current revision and content hash from the generated change log and from a full read; the timestamp is checked to move forward only when a publication happened |
| INV-10 | database state reconciles with client-observed committed outcomes | the harness's ledger of acknowledged operations is reconciled against the database after every writing phase |
| INV-11 | configuration never crosses installed-product, environment or business-unit boundaries | every read cross-checks that no returned key belongs to another product, and that a product's rows carry its own environment/business unit |
| INV-12 | an atomic batch is completely visible or completely absent | a batch is published and then read; a partial application is a violation |
| INV-13 | retried ambiguous outcomes do not create unexplained duplicate publications | every retry is logged with its attempt number and the observed revision deltas are reconciled to the acknowledged set |

Ties are handled explicitly: where two rows could legitimately answer the same question, the check
accepts any member of the tied set rather than demanding one particular row
(`docs/methodology.md` §5).

### 4.2 Negative controls

| Control | Invariant it must violate | What "fired" means |
|---|---|---|
| `x1_lost_update_control` | INV-1 and INV-3 | at least one of two concurrent updates to the same key is missing from the final state, while both were acknowledged |
| `x2_rollup_drift_control` | INV-9 and INV-10 | the stored rollup disagrees with the value derived from the entries after a concurrent publication |

If a control does **not** fire, the associated correctness claim is recorded as **not
demonstrated**, the workload is hardened by the pre-authorised escalation ladder in §20, and the
run is repeated. The invariant is never weakened. Every topology reports, for each control,
whether it fired.

## 5. Dataset, cardinality and byte regimes

### 5.1 Deterministic logical dataset

Generated from a fixed seed in Go, identical for every design and every engine.

| Property | Value | Notes |
|---|---|---|
| Seed | `42` (`-seed`, repository default) | changing it invalidates every comparison |
| Product definitions | 10 | each with a name and a version string |
| Environments | 3 | `production`, `staging`, `development` |
| Business units | 4 | `payments`, `identity`, `reporting`, `retail` |
| Installed products (fleet) | `-fleet`, default **500** for cadence work, **60** for cardinality work | one installation per (definition, environment, business unit) combination, sampled with bounded skew |
| Entries per installed product | one of the exact tiers in §5.2 | the study's cardinality axis |
| Value regimes | `small` and `large` (§5.3) | both are always recorded, never pooled |
| Key popularity | bounded Zipf, exponent 1.1, capped at 20x the median | hot keys exist but are not a strawman |
| Installation popularity | bounded Zipf, exponent 1.0 | hot and cold installed products |
| Read/write skew | 80/20 (80 % of reads target 20 % of installed products) | configurable `-read-skew` |
| History retention | keep the newest `H = 24` revisions per installed product; older revisions are aged out by the same policy for every design | recorded in the result's `options` |

The dataset must be **shaped like the real thing** (`docs/methodology.md` §6): skewed, long-tailed
but bounded, and generated in arrival order so that revision and time correlate.

### 5.2 Cardinality tiers

Exact tiers, all retained:

**1, 10, 30, 60, 120, 500**

- **60 is mandatory** and is never averaged into another tier.
- Two or more tiers below 60: 1, 10, 30.
- Two or more tiers above 60: 120, 500.

Two measurements are taken over this axis:

1. **Exact configuration counts**, for locating design crossovers.
2. **A realistic skewed distribution** bounded by the configured maximum: per-installed-product
   counts drawn from a bounded Zipf over the same tiers, so the tier value is a *ceiling* rather
   than a uniform count.

### 5.3 Byte regimes and byte accounting

Two value regimes, both measured:

| Regime | Key bytes | Value bytes | Shape |
|---|---|---|---|
| `small` | ~24 | ~96 | short scalar settings (`timeout=30s`), the common case |
| `large` | ~24 | ~4096 | serialized structured settings (a nested JSON blob), the costly case |

Every result records, **per design and phase**: key bytes, value bytes, **serialized configuration
bytes** (the size of one complete effective configuration for that installed product), **returned
bytes** (what a complete read actually transferred), and **row/document overhead** (stored bytes
minus logical bytes). Key count is never equated with document size.

### 5.4 Controlled cardinality comparison

Changing entries per installed product alone also changes total volume, which would make every
design look different for the wrong reason (`docs/methodology.md` §8a: *change the shape, hold the
volume*). So the cardinality comparison is run twice:

1. **Constant total entries** — the fleet shrinks as entries per product grow, so the total number
   of configuration entries is approximately constant across tiers.
2. **Constant total serialized bytes** — the fleet shrinks so that total serialized configuration
   bytes are approximately constant across tiers.

A third, **separately labelled** scenario is the **fixed fleet**: fleet held constant (60
installed products) while entries per product vary over the tiers. It answers a different question
(one portal's storage growth) and is never pooled with the two controlled comparisons.

## 6. Cadence: definitions and calculations

**Cadence is not throughput and not concurrency.** Three different quantities, never conflated:

- **Cadence** — how often one installed product publishes;
- **Total fleet arrival rate** — the offered load the portal must absorb;
- **Concurrent writers** — how many writers touch one installed product at once.

The defining relation, which the handoff computes rather than asserts:

```
global offered update rate  λ  =  installed products  /  update period
```

Worked through for the default cadence fleet of **500 installed products**:

| Regime | Period | λ (offered updates/s) | How it is measured |
|---|---|---|---|
| daily | 86 400 s | 0.0058 | fixed-count low-rate validation + pre-aged history; **projection only** where stated |
| hourly | 3 600 s | 0.139 | fixed-count low-rate validation + pre-aged history |
| minutely | 60 s | 8.33 | jittered open-loop arrivals, measured directly |
| secondly | 1 s | 500 | jittered open-loop arrivals **and** synchronized bursts, measured directly |

Rules that follow, and that the harness must obey:

- Low-rate regimes are **never** measured by waiting. They use a fixed operation count, a
  pre-aged history state, and an explicitly labelled capacity calculation.
- **Never present a calculated capacity projection as a directly measured temporal result.** Every
  result records whether its number is `measured` or `projected`, and projections are labelled in
  the report.
- Arrivals are **jittered** across the period in the primary regime, and **synchronized** in the
  burst regime (all products in a burst publish at the same instant, `-burst-size` products per
  burst tick). Both are required; they are different workloads.
- History growth is simulated by **accelerated history generation** *before* measurement, so the
  measured update lands on a realistically deep history rather than on a fresh table.

Recorded for every cadence cell: offered operations, started operations, completed operations,
queue depth, rejected operations, scheduling lag, throughput, latency percentiles, and client
saturation. The generator's own limit must be visible, not inferred
(`docs/methodology.md` §4).

## 7. Workloads

### 7.1 Reads (`queries.sql`, identical names and result shapes in every design)

| Op | Description | Reads a complete configuration? |
|---|---|---|
| `r01_effective_config` | the complete effective configuration of one installed product | yes — must be one coherent revision (INV-2) |
| `r02_read_key` | one configuration key | no |
| `r03_revision_check` | is a newer revision available than the caller's? | no — must not transfer the configuration |
| `r04_list_installations` | list installed products with product definition, environment, business unit, configuration count, current revision, last modification time | no (rollup-dependent) |
| `r05_search_key` | one key across installations of a product definition, an environment, or a business unit | no (the index/rolldown axis) |

### 7.2 Writes (`writes.sql`)

| Op | Description |
|---|---|
| `w01_modify_key` | modify one existing key |
| `w02_add_key` | add one key |
| `w03_delete_key` | delete one key |
| `w04_publish_batch` | atomically publish a batch of changes (revision +1) |
| `w05_replace_all` | atomically replace the complete configuration |
| `w06_update_metadata` | update ordinary installed-product metadata, independent of configuration |

`w01`–`w03` are single-change publications. `w04` and `w05` are the atomic multi-change paths.

### 7.3 Contention scenarios

| Scenario | Meaning |
|---|---|
| same key, same installed product | the true hot spot (`x1` must fire here) |
| different keys, same installed product | contention on the parent/rollup/pointer rather than on the entry |
| different installed products | no logical conflict at all — the baseline for "cost of one update" |

Blended workload: isolated reads and isolated writes are always reported **together with** a
blended run whose headline is the **fraction of read throughput retained** once writers appear
(`docs/methodology.md` §7).

## 8. Sizing, connections and environment

### 8.1 Live resource calculation (recomputed 2026-09-21, not copied)

| Quantity | Value | Source |
|---|---|---|
| Guest CPUs | 8 | `podman machine ssh -- nproc` |
| Guest memory | 16 496 422 912 B (≈15.36 GiB) | `free -b` in the podman machine |
| Guest kernel | `6.6.87.2-microsoft-standard-WSL2` | `uname -r` |
| Host | ASUS Zenbook S 14 UX5406SA, Intel Core Ultra 7 258V, 8 logical CPUs, 31.48 GB | Windows `Win32_ComputerSystem` / `Win32_Processor` |

Reservations, then the database budget:

| Reserve | CPU | Memory |
|---|---|---|
| Podman machine + WSL VM + OS + supporting work | 2 | ~4.4 GiB |
| Benchmark client | 2 | 2 GiB |
| **Database budget** | **6** | **9 GiB** |

Conditions (all named in every result; never pooled):

| Condition | Budget | Note |
|---|---|---|
| per-node baseline | 2 CPU / 3 GiB **per node** | `pg-single`, `yb-single`, and each of `yb-n1..n3` |
| **equal-total control** | 6 CPU / 9 GiB on a **single** node | labelled separately; answers a different question |
| equal-total cluster | 3 × (2 CPU / 3 GiB) | the same 6 CPU / 9 GiB spread over three nodes |

Limits are CFS quotas (`--cpus`), **never** core pins — this host mixes P-cores and E-cores that
the guest kernel cannot distinguish (`docs/environments/host-zenbook-ux5406sa.md`).

### 8.2 Connection budget

| Role | Count |
|---|---|
| readers | 8 |
| writers | 8 |
| load streams (COPY) | 4 |
| control connections (revision check, marker) | 1 |
| administrative/inspection | 1 |
| **allowance** | **22** |

Checked against the engine's pool, recorded per cell. `yb-cluster3` spreads connections over
**all three** nodes through the harness's comma-separated DSN list; a single query endpoint is
never silently used as a stand-in for balanced cluster access. Per-node connection distribution is
recorded.

### 8.3 Topologies

`pg-single`, `yb-single`, `yb-cluster3` — the repository's pinned scripts
(`infra/pg-single.sh`, `infra/yb-single.sh`, `infra/yb-cluster3.sh`). Images are pinned in
`infra/versions.env`; this study changes neither.

**Colocation is verified, not assumed.** Running all three nodes inside one WSL2 VM is *not*
evidence of data colocation and gives no real network-latency comparison. Placement evidence is
captured from the engine's own catalogs and `EXPLAIN (ANALYZE, DIST)` **RPC counts**, which are the
portable distributed-cost signal.

## 9. Experiments

| Phase | Name | Content | This session |
|---|---|---|---|
| A | Engine and semantic probes | isolation, JSON update behaviour, transaction timestamps, error classification, `EXPLAIN` support, placement mechanism, query endpoints | **yes** |
| B | WSL bind-mount probe | proof that a file in the task worktree is visible at the expected container path and that output written to the mounted results path survives | **yes** |
| C | Tiny dev checks | every design, invariant, negative control, read, write, report generator and digest check on `tiny` | **yes** |
| D | Broad survey | all retained designs, primary topologies, cardinality tiers, isolated reads/writes, plans, storage | **reduced** — one small matrix on `pg-single` |
| E | Cardinality crossover | repeated trials around the observed boundary, including 60 | mapped, later |
| F | Concurrency | optimistic, pessimistic, naturally non-conflicting and the unsafe control under identical workloads | mapped, later |
| G | Cadence and mixed workload | daily/hourly/minutely/secondly targets, jittered arrivals, synchronized bursts | mapped, later |
| H | Sustained churn | state growth, vacuum, WAL, compaction, document rewrite cost, retained history, materialization correctness | mapped, later |
| I | Deployment controls | per-node baseline, equal-total resources, balanced endpoints, verified colocated/non-colocated placement | mapped, later |
| J | Repeated conclusion runs | every number used in a conclusion gets repeated trials, median, individual values and spread | mapped, later |

### 9.1 Cell counts and duration

For a single trial and one topology, `-cmd full` visits one cell per design: 16 designs on
PostgreSQL (`y1`/`y2` are YugabyteDB-only), 18 on YugabyteDB. Estimated wall-clock from study 03's
observed ~25 min/cell at `small`:

| Matrix | Cells | Estimated |
|---|---|---|
| Phase C dev checks, `tiny`, 3 designs × 2 topologies | 6 | ~20 min |
| Phase D reduced (this session), `small`, 4 designs × `pg-single` | 4 | ~1.5 h |
| Phase D full, `small`, 16 × pg-single + 18 × yb-single + 18 × yb-cluster3 | 52 | ~22 h |
| Phases E–I | ~60 | ~30 h+ |
| Phase J repeats (3 trials on the conclusion set) | ~40 | ~20 h+ |

The full campaign is therefore a **multi-session** effort and is explicitly *not* this session's
scope (§22). This session ends after Phase A, B, C and the reduced Phase D.

### 9.2 Unbiased matrix-reduction rules

Used only where a full matrix cannot fit, and never to rescue a convenient conclusion:

1. Reduce **topologies** before designs: a design not run on an engine is a coverage gap, a design
   not run at all is a lost comparison.
2. Never drop a **negative control** or the reference `n1`.
3. Never drop the **60** tier.
4. Never pool two conditions with different budgets, replication factors or placement.
5. Reduce **trials** last, and only in a survey phase; every conclusion comes from Phase J repeats.
6. Any reduction is written into the run's manifest with the reason, and is listed as a coverage
   gap in the report and the analysis.

## 10. Harness implementation

Go, on the shared `adsplatform` module, hexagonal like studies 02 and 03: `main.go` is the only
file that imports an adapter.

### 10.1 Files

```
studies/04-configuration-portal/harness/
  main.go        flags, phases, wiring (the entry adapter)
  designs.go     the 18 designs: their tables, statements, strategy and flags
  dataset.go     deterministic generation: definitions, installations, entries, change log
  load.go        schema, bulk load, index build, pre-aged history
  verify.go      INV-1..INV-13 gates; computes expectations independently from the dataset
  audit.go       post-phase audits for every writing/contention/churn phase
  workload.go    the six reads and six writes, isolated and blended
  cadence.go     cadence scenarios, jittered and synchronized-burst schedules
  contention.go  the three contention scenarios
  report.go      the generated report (numbers only, TL;DR by fixed rules)
  *_test.go      unit tests
```

### 10.2 Mode of a cell

Exactly study 02/03's shape so the reporting code and the provenance rules are shared
understanding: `-cmd full` runs `verify → explain → read → write → mixed → audit` for one design on
one topology and writes one JSON result plus one readable plans file. A cell that fails a gate
aborts and writes no timing.

### 10.3 Flags (beyond study 02/03's, whose meanings are kept)

| Flag | Meaning |
|---|---|
| `-entries-per-ip` | the cardinality tier (1, 10, 30, 60, 120, 500) |
| `-fleet` | number of installed products |
| `-value-regime` | `small` \| `large` |
| `-cardinality-mode` | `exact` \| `constant-entries` \| `constant-bytes` \| `fixed-fleet` \| `skewed` |
| `-cadence` | `daily` \| `hourly` \| `minutely` \| `secondly` \| `none` |
| `-burst-size` | products publishing together in a synchronized burst |
| `-concurrency` | concurrent writers per installed product |
| `-retries` | bound on optimistic retries |
| `-phases` | as study 03 |

### 10.4 Platform change

One additive file: `platform/core/measure/arrival.go` — an open-loop arrival scheduler (jittered
and synchronized schedules) recording offered/started/completed, queue depth, rejections and
scheduling lag. The existing closed-loop driver is left untouched, so studies 01–03 binaries are
unchanged. This is the study's only shared-code change and is tagged `repo/wsl-podman-bridge`'s
sibling milestone.

## 11. SQL catalogue conventions

One directory per design with `schema.sql`, `indexes.sql`, `queries.sql`, `writes.sql`,
`audit.sql`, in the repository's annotated format (`-- name:`, `-- params:`), embedded with
`go:embed` so the SQL that produced a result is provably the SQL beside the binary.

Rules:

- Read statement names and result shapes are **identical in every design** (`r01`…`r05`).
- Parameter binding is by name through `-- params:`; the harness decides *which* statement runs,
  in *which* transaction, at *which* isolation level.
- Explicit casts on parameters in `INSERT … SELECT` and `unnest` (study 02's rule).
- `w_load_*` statements run once after the bulk copy; the explain phase skips them.
- Index column order is chosen for YugabyteDB's defaults (the first column is hash-sharded).
- Every design's `audit.sql` recomputes the truth **from that design's own tables**, which is what
  makes the audit independent of the query it checks.

## 12. Exact files created by this study

```
studies/04-configuration-portal/
  README.md  HANDOFF.md  ESCALATIONS.md  PROGRESS.md  .gitattributes  study.env  sqlfs.go
  Containerfile  run-study.sh  probe-engines.sh  probe-bind-mount.sh
  sql/{README.md,<18 design dirs>/{schema,indexes,queries,writes,audit}.sql}
  harness/*.go
  diagrams/{*.puml,render.sh,rendered/*.svg}
  results/  reports/{,analyses/,discussions/,outdated/}
```

Shared files touched (each additive, each named in its commit):

- `infra/lib.sh` — the WSL→Podman engine resolver and `winpath()`;
- `platform/core/measure/arrival.go` — the open-loop scheduler;
- `CONTEXT.md`, `LESSONS_LEARNED.md` — kept current as part of the work.

## 13. Exact commands

```bash
# --- worktree (already created) -------------------------------------------------
cd /mnt/c/extra/code/architecture_design_studies
git worktree add .worktrees/study04-configuration-portal -b study-04/configuration-portal

# --- engine reachability --------------------------------------------------------
export PATH="$PATH:/mnt/c/Program Files/RedHat/Podman"   # only needed for probes
podman.exe version && podman.exe ps --all && podman.exe volume inspect ads-run-lock

# --- phase A/B probes (each takes the benchmark lock) ---------------------------
cd studies/04-configuration-portal
./probe-bind-mount.sh          # phase B
./probe-engines.sh             # phase A

# --- phase C dev checks (tiny, not results) -------------------------------------
./run-study.sh --scale tiny --topologies pg-single,yb-single \
  --designs n1_rows_indexed,c1_optimistic_version,x1_lost_update_control \
  --phases verify,explain,read,write,audit --run-id <id>

# --- phase D reduced (this session) ---------------------------------------------
./run-study.sh --tag --scale small --topologies pg-single \
  --designs n1_rows_indexed,n2_rows_rolldown,d1_doc_row,c1_optimistic_version

# --- report regeneration after an analysis is added ------------------------------
podman run --rm -v "$(hostpath "$PWD"):/study" "${BENCH_IMAGE:-localhost/configbench:1}" \
  -cmd report -results "/study/results/<run-id>" -report-out "/study/reports/<run-id>.md"
```

Every command that starts a container or measures acquires the benchmark lock through
`run_lock_acquire` in `infra/lib.sh`. Nothing here starts a second database while another session
holds the lock.

## 14. Expected results

- Phase A: a written probe record stating the **effective isolation** on each engine, whether
  `jsonb_set` and whole-document replacement differ measurably, which `EXPLAIN` options are
  supported, the **actual colocation mechanism** on this image (and whether it can be expressed at
  all), and the live endpoint list.
- Phase B: proof that the task worktree is reachable inside the container at the expected path and
  that a file written from the container survives on the host, plus the answer to "does the engine
  accept `/mnt/c/…` bind sources or must they be `C:/…`?".
- Phase C: all designs pass verification on both engines; all reads and writes execute; plans are
  captured; report generation and digest computation succeed; **both negative controls fire** on at
  least the contended scenario.
- Reduced Phase D: four `pg-single` cells with valid gates, generated report, inputs digest, and
  repeatability across two regenerations.

## 15. Acceptance criteria

A milestone is complete only when **all** of the following hold:

1. `go vet ./...` and `go test ./...` pass **inside the image build**, so a binary whose tests fail
   cannot exist.
2. Every measured cell passed its correctness gate; failures are listed with their cause.
3. Both negative controls fired on every topology where they were run, or the affected claim is
   recorded as not demonstrated.
4. Every result carries environment, repository commit, `git describe`, dirty flag, image identity,
   and every report carries an **inputs digest**.
5. The digest is stable across two independent regenerations and across a WSL/Windows-style
   checkout byte round trip (this is what the study `.gitattributes` is for).
6. The generated report contains no interpretation and opens with a TL;DR selected by fixed rules.
7. The signed analysis states its weaknesses and discloses single-model provenance.
8. The run was produced from a **clean, committed tree** tagged `run/04-configuration-portal/<id>`.
9. `main` contains the task commit and the integrated milestone is tagged.

## 16. Stop conditions

Work stops and the checkpoint is committed when any of these is true:

- the benchmark lock is held by another session (wait is not a stop; it is LOW-style routine work);
- the EOL-only proof over the main checkout's dirty paths fails (a real edit exists);
- a Phase A/B probe contradicts a premise of this handoff and no mapped rule covers it;
- a negative control still will not fire after the §20 harshening ladder;
- a cell's results are ambiguous in a way INV-1…INV-13 cannot classify.

In every case: commit the safe checkpoint, write the escalation, and end the response with
`NEXT MODEL: HIGH — continue <specific remaining work> using DeepSeek HIGH in WSL`.

## 17. Escalation triggers

Escalations are recorded in `ESCALATIONS.md` using `docs/templates/ESCALATION_REQUIRED.md`, the
dependent work stops, HIGH decides in a separate dated section, an amendment `AM-NN` is published
here, and the whole thing is committed and tagged `study-04/v0.N-handoff-amendment-NN`.

| Trigger | Blocks |
|---|---|
| **ER-01** — no native WSL `podman`; only the Windows `podman.exe` reaches the documented engine | all container work; AM-01 decides the bridge |
| the bind-mount probe fails or mounts an empty directory | all measurement |
| colocated placement cannot be expressed or evidenced on the pinned image | the colocation requirement (recorded as a coverage gap, never as a completed comparison) |
| a control does not fire after the harshening ladder | the associated correctness claim |
| the main checkout's dirt is not EOL-only | integration |

## 18. Decision Log — what the executor may decide freely and must record

Routine, reversible choices belong in `PROGRESS.md`, not in an escalation:

- exact synthetic key names and value strings inside the declared byte budgets;
- the order in which cells run (from the run's order seed);
- transaction batching, statement timeouts and retry back-off shape, within the declared bounds;
- internal struct, file and function naming;
- which statement in a design carries a given op, as long as the op's contract in §7 is met;
- duration/warmup per phase at `tiny`, and the exact dev-check command line;
- the precise wording of a generated table, as long as it carries measurements only.

**Never** changed silently (they are the scientific content): the question, design semantics,
invariants, negative controls, workload meaning, acceptance criteria and interpretation rules.

## 19. Reporting, tagging and integration

- Commit every meaningful step with **explicit paths**, identifying the agent (model id and role)
  in the message. Never `git add -A`.
- Tags, all annotated, checked with `git tag` first, never moved, deleted or reused:
  `study-04/v0-handoff` → `study-04/v0.1-handoff-amendment-01` → `repo/wsl-podman-bridge` →
  `study-04/v1-harness` → `run/04-configuration-portal/<run-id>` → `study-04/v1-measured` →
  `study-04/v1-analysis` → `study-04/v1-integrated`.
- Generated reports hold measurements and no interpretation; conclusions live in a signed analysis
  under `reports/analyses/`, with mechanisms in a signed companion under `reports/discussions/`.
- Integration: prove the main checkout's dirt is EOL-only, normalize the working tree to its LF
  index, integrate the latest local `main` into the task branch, reconcile the shared documents,
  run the combined checks, fast-forward `main`, verify reachability, tag the integrated state.
- **Never push, never pull, never change remotes.** The broken cross-platform worktree
  `.worktrees/recency-reports` is reported to the owner and left untouched.

## 20. Pre-authorised harshening ladder (for a control that does not fire)

In order, stopping as soon as the control fires:

1. narrow the contending key set (more writers on fewer keys);
2. raise `-concurrency` up to the connection allowance in §8.2;
3. hold the writers inside their transaction longer (`-hold-ms`);
4. add jitter to arrival so the race is not accidentally serialised by the scheduler;
5. run the scenario on a freshly loaded database with no other phase before it.

Each rung is recorded with its evidence. The invariant is never weakened.

## 21. Amendments

*(HIGH appends dated, attributed `AM-NN` decisions here. Previous text is never erased; superseded
instructions are marked as superseded.)*

### AM-01 — ER-01 decided: drive the documented Windows Podman engine from WSL

**Decided by:** DeepSeek HIGH (`deepseek-flash`), 2026-09-21.

**Evidence:** `podman` is absent from the WSL distro (`command -v podman` fails; no binary in
`/usr/bin` or `/usr/local/bin`). `podman.exe` exists at
`/mnt/c/Program Files/RedHat/Podman/podman.exe`, reports client 5.8.1 (windows/amd64) against
server 5.8.5 (linux/amd64) over the `podman-machine-default` SSH connection, and from WSL already
sees the documented storage: images `localhost/seatbench:1` and `localhost/ticketbench:1`, volumes
including `ads-results`, the exited container `pg-scratch`, and **no** `ads-run-lock` (free). The
podman machine also sees the repository at `/mnt/c/extra/code/architecture_design_studies`, the
same path the WSL distro uses. Live guest resources are 8 CPUs and 16 496 422 912 bytes with kernel
`6.6.87.2-microsoft-standard-WSL2`, matching `docs/environments/host-zenbook-ux5406sa.md`.

**Decision:** WSL drives the engine through `podman.exe`. `infra/lib.sh` gains an explicit,
additive engine resolver (`PODMAN`/`ADS_PODMAN` → `podman` on PATH → the Windows client) and a
`podman()` function so every existing call site, and `command -v podman`, keep working. A WSL-local
podman is **never** initialised: a second engine would have different storage and no shared
benchmark lock, which is precisely the failure the lock exists to prevent. `winpath()` is added for
the case where `podman.exe` requires Windows-form bind sources; Phase B decides whether it is
needed, and `hostpath()`'s existing behaviour is preserved unless the probe proves otherwise.

**Supersedes:** nothing. This amendment *adds* the engine resolution the repository previously got
from having `podman` on `PATH`.

**Probe outcome (measured 2026-09-21, `results/devchecks/phase-b-bind-mount/`):** the phase B
bind-mount probe ran before any measurement and both spellings mount the study tree correctly from
WSL:

| Mount source as passed | Sentinel visible | Container read-back identical | Host write survived | Verdict |
|---|---|---|---|---|
| `/mnt/c/...` (`hostpath`) | yes | yes | yes | **WORKS** |
| `C:/...` (`winpath`) | yes | yes | yes | **WORKS** |
| `/mnt/c/.../definitely-not-here` (control) | no | no | no | **FAILS** |

`RUNNER_DECISION=hostpath`. The third row is the control that makes the first two meaningful: a
probe that reports "works" for a path that cannot exist would be worthless. Because `hostpath()`
already works, it is left **byte-for-byte unchanged** and `winpath()` remains an unused, documented
fallback for a host where the Windows form is required. The runner therefore reads exactly like the
other studies' runners; no new path handling is introduced anywhere.

**Next setting:** **HIGH — execute this handoff using the same DeepSeek HIGH model in WSL.**

---

## Amendment 01 (v2 completion)

`HANDOFF-AMENDMENT-01-V2.md` extends this handoff with the three review-mandated controls
(repeated-design instrument control, `-retries 1` contention control, writer sweep), the cardinality
and cadence phases, repeated randomized trials, the equal-total-resource arm, verified placement and
YugabyteDB coverage. It does not measure; it defines the acceptance criteria and the separately
claimable execution tasks. Read it before running any v2 cell.
