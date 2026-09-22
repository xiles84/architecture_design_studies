# Goal amendment — deliberate multi-leader brainstorms

**Amendment id:** `20260922T114507Z-multi-leader-brainstorms`

**Recorded at:** `2026-09-22T11:45:07Z`

**Owner request:** provide an explicit way for several leader/HIGH sessions to challenge
an idea, preserve the discussion, distinguish ongoing from concluded work, and report
when a conclusion is ready for separately authorized task creation.

## Decision

Add a durable brainstorm archive and Git-CAS contribution claims alongside, but separate
from, the implementation queue. The owner starts it only with `START BRAINSTORM:` and
continues it with `CONTINUE BRAINSTORM <id>`. Defaults are three blind independent
positions, two cross-reviews, and one synthesis. A synthesis may preserve dissent.

Every state exposes a compact summary. Conclusion emits `BRAINSTORM ENDED` and explicitly
states that no tasks were created. Only a later
`CREATE TASKS FROM BRAINSTORM <id>` authorizes a chosen HIGH session to publish tasks.
Those task records declare `originating_brainstorm_id` and are then linked back.

Ongoing and concluded items have separate generated indexes, while immutable records
stay at stable paths. Contributions and events form the historical record. Reusable
lessons are added only after later outcomes support them; the debate itself remains in
the brainstorm archive.

## Effect

This amendment extends AI coordination and task provenance. It does not change study
evidence, authorize a benchmark, weaken HIGH/LOW task eligibility, or make consensus a
correctness criterion.
