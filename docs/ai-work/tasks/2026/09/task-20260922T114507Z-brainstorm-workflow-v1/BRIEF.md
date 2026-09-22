# Execution brief — deliberate multi-leader brainstorm workflow v1

## Goal

Add an explicit, durable brainstorm workflow in which several HIGH-capability sessions
can propose independently, cross-review one another, synthesize a conclusion, and preserve
the full history. Starting a brainstorm and turning its conclusion into implementation
tasks are separate, deliberate owner actions.

## User intents

Repository agents recognize these case-insensitive intents only when the user's first
nonblank line begins with the shown phrase:

- `START BRAINSTORM: <question>` — create a brainstorm. An ordinary question never does.
- `CONTINUE BRAINSTORM <brainstorm-id>` — claim and perform the next eligible contribution.
- `LIST BRAINSTORMS` — show ongoing and concluded summaries separately.
- `BRAINSTORM STATUS <brainstorm-id>` — show one current summary.
- `CREATE TASKS FROM BRAINSTORM <brainstorm-id>` — after conclusion, let this chosen
  session publish implementation tasks and link them to the brainstorm.
- `CANCEL BRAINSTORM <brainstorm-id>` and `SUPERSEDE BRAINSTORM <old-id> WITH <new-id>`
  are explicit administrative endings.

Do not treat casual uses of "brainstorm", "leader", or datastore leader/follower terms
as commands. Natural-language responses may be friendly, but the archive transition
occurs only through the explicit intent and CLI command.

## Durable layout

```text
docs/ai-work/brainstorms/
  WORKFLOW.md
  SCHEMA.md
  ONGOING.md                 generated index
  CONCLUDED.md               generated index
  templates/
  records/<yyyy>/<mm>/<brainstorm-id>/
    brainstorm.json          immutable question, criteria, targets and correlations
    QUESTION.md              owner question and evidence packet
    STATE.md                 generated current snapshot
    positions/<contribution-id>.md
    critiques/<contribution-id>.md
    CONCLUSION.md            absent until synthesis concludes
    events/<timestamp>-<event-id>.json
```

Indexes and `STATE.md` are disposable generated views. `brainstorm.json`, immutable
contributions, `CONCLUSION.md`, and chained event files are the portable authority.
Never move a record between ongoing and concluded directories: stable paths preserve links.

## Lifecycle and default round

```text
collecting_positions -> cross_review -> synthesis -> concluded -> tasked
                                  \-> cancelled
concluded|tasked -> superseded
```

Default targets are three independent positions and two cross-reviews, configurable at
start. The initiating leader may fill the first position. Position contributors receive
only `QUESTION.md`, the evidence paths and applicable repository rules; the CLI and brief
must not expose earlier positions before the target is complete. Critiques read all
positions. The synthesizer reads positions and critiques, applies the stated decision
criteria, and may conclude without unanimity. It preserves strongest dissent, rejected
alternatives, assumptions, confidence, evidence gaps, and falsification tests.

The workflow seeks a defensible decision, not consensus. A conclusion is allowed when all
required slots are complete and remaining objections are answered or explicitly retained.
An evidence gap that prevents a safe decision remains visible and may produce a proposed
research task only after the explicit task-creation intent.

## Coordination

Implement brainstorm commands in the existing Go queue binary and Podman wrapper. Use
Git compare-and-swap refs under `refs/ads-brainstorms/live/<brainstorm-id>/<slot-id>` for
contribution claims. Use the existing local-main integration lock when committing
brainstorm archive mutations so simultaneous submissions serialize without sharing a
working folder. A contribution claim has the queue's two-hour lease, heartbeat/guard
semantics and epoch rejection; reuse code rather than create a weaker parallel mechanism.

Only canonical HIGH sessions may start, claim, submit, conclude, cancel, supersede, or
create/link tasks. Capability aliases normalize first. A contribution claim determines
the next eligible unclaimed slot and role; parallel claimers may take different ready
slots, but never the same slot. Events and contribution ids remain unique after retries.

Archive writes are a narrow coordination operation, like queue event publication: the
CLI owns them, takes the integration lock, refuses tracked-dirty local `main`, stages only
brainstorm paths, and commits with actor identity. Agents never hand-edit another model's
contribution. Live refs are host-local; a fresh clone reconstructs the same state from
committed records.

## Commands

Choose names consistent with the existing flat CLI, preferably:

```text
queue brainstorm-start
queue brainstorm-list
queue brainstorm-status
queue brainstorm-claim
queue brainstorm-guard
queue brainstorm-heartbeat
queue brainstorm-submit
queue brainstorm-cancel
queue brainstorm-supersede
queue brainstorm-link-tasks
queue brainstorm-audit
```

`brainstorm-start` accepts a question file or explicit question, title, optional decision
criteria/evidence paths, and position/critique targets. `brainstorm-submit` accepts the
claimed slot's contribution file; a synthesis submission writes `CONCLUSION.md` and moves
to `concluded`. `brainstorm-link-tasks` requires a concluded brainstorm and already
published task ids, verifies each task exists and correlates them, then moves to `tasked`.
It never invents tasks itself: the chosen AI authors and publishes them after the owner's
`CREATE TASKS FROM BRAINSTORM` intent.

## Mandatory human-facing output

Every state-changing or status command emits this compact block (and structured JSON under
`--json`):

```text
BRAINSTORM STATUS
ID: <id>
State: ongoing|concluded|tasked|cancelled|superseded
Stage: <internal stage>
Summary: <two to four short lines>
Disagreements: <short text or none>
Next action: <exact user intent>
```

When synthesis concludes, additionally emit:

```text
BRAINSTORM ENDED
ID: <id>
Conclusion: <few lines>
Remaining dissent: <few lines or none>
Confidence: <high|medium|low>
Tasks: NOT CREATED

To proceed:
CREATE TASKS FROM BRAINSTORM <id>
```

`brainstorm-list` prints `ONGOING` and `CONCLUDED` sections and a two-to-four-line summary
for every item. A tasked brainstorm remains in the concluded section and lists linked task
ids. Empty sections are explicit.

## Content contracts

Provide templates for a question/evidence packet, independent position, critique, and
conclusion. Every contribution records brainstorm/contribution/parent ids, slot, model,
tool, effort, session, canonical capability, timestamp, files/evidence examined,
recommendation, confidence, assumptions, disagreements, alternatives, and tests that
could falsify it. Preserve raw capability input when available.

`LESSONS_LEARNED.md` receives only reusable lessons supported by later outcomes, never a
copy of the debate. `CONTEXT.md` links current brainstorm indexes and major concluded
decisions without becoming the state authority.

## Tests and acceptance

Run all Go tests in the pinned Podman image. Add tests for explicit-trigger parsing,
ordinary-question non-triggering, stage dependencies, blind independent positions,
twenty claimers per slot, simultaneous claims of different slots, leases/old epochs,
invalid chains and identities, exact summary/end output, stable generated indexes,
conclusion without unanimity, task creation refusal before conclusion, explicit task
linking, audit, and fresh-clone reconstruction. Existing queue tests must remain green.

Do not run a database or benchmark. Use Decide -> Log -> Continue for implementation
details. Escalate only if the existing Git CAS/main-lock model cannot safely serialize
brainstorm archive writes, or if enforcing blind positions requires changing the decided
portable archive boundary.

On completion submit for HIGH review; do not integrate before approval.

NEXT MODEL: LOW
