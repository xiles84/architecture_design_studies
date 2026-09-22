# Brainstorm archive schema

Schema version 1 separates immutable authority from generated navigation.

```text
docs/ai-work/brainstorms/
  WORKFLOW.md
  SCHEMA.md
  ONGOING.md
  CONCLUDED.md
  templates/
  records/yyyy/mm/<brainstorm-id>/
    brainstorm.json
    QUESTION.md
    STATE.md
    positions/<contribution-id>.md
    critiques/<contribution-id>.md
    CONCLUSION.md
    events/<timestamp>-<event-id>.json
```

`brainstorm.json` is immutable. It records the stable id, optional goal/root/parent
task correlations, title, exact question, decision criteria, evidence paths, position
and critique targets, creator identity, and creation time.

Each immutable event records schema/id/predecessor/sequence, brainstorm id, UTC time,
full actor identity, event type, resulting stage, and a short summary. Submission events
also record slot, contribution id, disagreements, and confidence. Administrative events
record linked task ids or the superseding brainstorm.

The permitted events are `started`, `position_submitted`, `critique_submitted`,
`synthesis_submitted`, `tasks_linked`, `cancelled`, and `superseded`.
The validator reconstructs the lifecycle, rejects broken predecessors, duplicate slots,
invalid actors/times/transitions, and missing contribution files.

`STATE.md`, `ONGOING.md`, and `CONCLUDED.md` are generated views and may be rebuilt.
The spec, chained events, contributions, and conclusion are the portable authority.
Records never move when their state changes, which keeps links stable.

Live claim blobs are host-local coordination, not archive state. A claim records
brainstorm/slot ids and kind, claim id and epoch, HIGH actor identity, claim/heartbeat/
expiry times, and the stage at claim time. Independent clones do not share these refs;
multi-host live coordination remains out of scope.

Queue task metadata may carry `originating_brainstorm_id`. Related task events may carry
`brainstorm_id` and `contribution_id`. These fields correlate implementation with the
decision record without making the brainstorm a task queue.
