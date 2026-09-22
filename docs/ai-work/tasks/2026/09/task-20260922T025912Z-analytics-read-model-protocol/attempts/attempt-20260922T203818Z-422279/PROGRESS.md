# Progress — Study 08 analytics/read-model protocol

**Task:** `task-20260922T025912Z-analytics-read-model-protocol`
**Attempt:** `attempt-20260922T203818Z-422279` (claim `claim-64e8882459beb8d5`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/analytics-read-model-handoff-v1`

Planning only: no image pulled, no database started, no measurement.

## Decision Log

1. **Four strategies under one workload and one correctness contract**: base OLTP, PostgreSQL
   materialized view, in-engine search copy, ClickHouse analytical copy. Confidence: High.
2. **Isolate the copy, not the engine, in the in-engine search arm.** A native search engine is a
   separately labelled extension so the study does not change two factors at once. Confidence: High.
3. **ClickHouse is pinned by line in the protocol and by tag+digest in the harness**; a moved tag
   blocks the run. Confidence: High.
4. **Measure refresh, staleness, storage and correctness for every strategy**, not just query gain;
   each copy needs a firing stale/duplicate control. Confidence: High.
5. **v0 single-node**, cluster topology belongs to Study 07. Confidence: High.
6. **Four execution tasks, harness first.** Confidence: High.

## Residual risks

- Three workloads may not separate the strategies; the harness task may add a fourth rather than
  report a tie.
- ClickHouse's loader cadence is part of the design, not a constant; the columnar task must measure
  the cadence it configures.
