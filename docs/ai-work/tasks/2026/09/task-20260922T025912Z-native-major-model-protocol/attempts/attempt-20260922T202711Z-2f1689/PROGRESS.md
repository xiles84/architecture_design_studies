# Progress — Study 06 native storage models protocol

**Task:** `task-20260922T025912Z-native-major-model-protocol`
**Attempt:** `attempt-20260922T202711Z-2f1689` (claim `claim-62febfbcce0b28c5`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `study-06/v0-handoff`

Planning only: no image pulled, no database started, no measurement.

## Decision Log

1. **Four technologies, each native to a family**: PostgreSQL (baseline), MongoDB (document),
   ScyllaDB (wide-column), Valkey authoritative (key-value). A PostgreSQL ledger represents history
   without claiming a new engine. Confidence: High.
2. **Pin the line in the protocol, the patch and digest in the harness.** The protocol fixes
   MongoDB 8.0 / ScyllaDB 6.2 / Valkey 8.1; the harness task resolves the exact tag and records the
   digest in `versions.env` and every manifest. A moved tag blocks the run. Confidence: High.
3. **Reuse Study 04 semantics verbatim** so the storage model is the only changed factor.
   Confidence: High.
4. **v0 is single-node.** Cluster/RF variants go to Study 07 so one thing changes at a time.
   Confidence: High.
5. **Valkey is the system of record, not a cache** — the boundary the registry demands be stated.
   Confidence: High.
6. **Four execution tasks, harness first**, each with its own report and signed analysis.
   Confidence: High.

## Residual risks

- Upstream tags may move between planning and harness build; the harness records the digest and
  blocks on a moved tag.
- A native key-value store cannot express the whole-configuration read the same way; the report must
  document the native analogue rather than fake the operation.
