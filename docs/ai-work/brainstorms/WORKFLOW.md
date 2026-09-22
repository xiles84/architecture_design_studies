# Deliberate multi-leader brainstorms

This workflow lets several HIGH-capability AI sessions challenge a question before any
implementation task is authorized. It is opt-in: ordinary questions and casual uses of
the word “brainstorm” do not create or advance a record.

## Owner commands

Put one of these phrases at the start of the first nonblank line:

- `START BRAINSTORM: <question>` creates a durable discussion.
- `CONTINUE BRAINSTORM <id>` asks the current HIGH session to claim and fill the next
  available position, critique, or synthesis slot.
- `LIST BRAINSTORMS` shows ongoing and concluded items separately.
- `BRAINSTORM STATUS <id>` shows the compact current state.
- `CREATE TASKS FROM BRAINSTORM <id>` authorizes the chosen HIGH session to turn an
  existing conclusion into queue tasks. It is deliberately separate from conclusion.
- `CANCEL BRAINSTORM <id>` ends an ongoing item without a conclusion.
- `SUPERSEDE BRAINSTORM <old-id> WITH <new-id>` retains both records and marks the
  replacement.

The natural-language phrase is an agent instruction. The agent validates it with
`queue brainstorm-intent`, then uses the matching queue command. Starting requires a
short title; the agent derives one from the question unless the owner supplies it.

## Round structure

The default is three independent positions, two cross-reviews, and one synthesis:

```text
collecting_positions -> cross_review -> synthesis -> concluded -> tasked
        |                    |              |
        +--------------------+--------------+-> cancelled
concluded | tasked -> superseded
```

Position authors read the question, evidence packet, and repository rules, but do not
read prior positions until the independent round closes. This is a workflow constraint,
not a secrecy boundary: the portable archive is committed Git content. Critique authors
read every position. The synthesizer reads all positions and critiques and preserves
material dissent; agreement is not required.

Every contribution states the evidence inspected, recommendation, reasoning,
assumptions, alternatives, disagreements, confidence, and what evidence could falsify
it. A contributor may identify a blocking evidence gap. Synthesis may retain that gap or
recommend research, but it still does not create a task.

## Claiming and recovery

Only HIGH sessions may mutate brainstorm state. Run `brainstorm-claim` after the owner
asks to continue. Each ready slot has a two-hour Git compare-and-swap lease under
`refs/ads-brainstorms/live/<brainstorm-id>/<slot-id>`; heartbeat at least every
15 minutes and guard before submission. Several sessions may claim different slots, but
only one wins any one slot.

An expired slot is reclaimed by the next `brainstorm-claim`; its claim epoch increases,
so the old claimant's next guard or submit is rejected. Unsubmitted draft text belongs
to that session and is not portable. Submitted contributions and events are committed
and reconstructible from a fresh clone.

Archive mutations are the same narrow coordination exception as task-event publication:
the CLI takes `refs/ads-queue/locks/main-integration`, refuses tracked edits in local
`main`, stages only owned brainstorm paths, commits them, and releases the lock. Agents
never hand-edit generated state or another contributor's file.

## Compact outputs

Every status names the id, overall state, internal stage, a short summary, progress,
remaining disagreements, and the exact next owner action. Lists always contain separate
`ONGOING` and `CONCLUDED` sections; cancelled and superseded records are historical
endings in the concluded index.

The synthesis response additionally prints exactly `BRAINSTORM ENDED`, states that
tasks were not created, and shows:

```text
CREATE TASKS FROM BRAINSTORM <id>
```

Only after that explicit owner command may a HIGH session publish implementation tasks.
Each task must declare `originating_brainstorm_id`; `brainstorm-link-tasks` verifies
the correlation before moving the record to `tasked`.

## History and lessons

The entire discussion remains under the stable record path described in
[SCHEMA.md](SCHEMA.md). Do not copy the debate into `LESSONS_LEARNED.md). Add a lesson
only when later implementation or evidence demonstrates a reusable outcome.
