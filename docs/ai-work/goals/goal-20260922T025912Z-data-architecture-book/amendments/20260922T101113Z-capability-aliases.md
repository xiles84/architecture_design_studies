# Goal amendment — session capability aliases

**Amendment id:** `20260922T101113Z-capability-aliases`
**Recorded at:** `2026-09-22T10:11:13Z`
**Owner request:** allow industry-style capability declarations instead of requiring
only `HIGH` and `LOW` in every new chat.

## Decision

Normalize explicit, case-insensitive session declarations as follows:

- `HIGH`, `leader`, `master` -> `HIGH`;
- `LOW`, `worker`, `follower`, `slave` -> `LOW`.

`master` and `slave` are legacy-compatible inputs. Agents and tooling accept them but do
not generate or recommend them; examples prefer `leader` and `worker`. The owner rejected
`primary` and `replica` for this purpose because they do not fit the AI pipeline and are
already meaningful topology terms in this repository.

Canonical task metadata, eligibility decisions, events, live refs, and routing remain
`HIGH`/`LOW`. Preserve the user's raw declaration separately when the schema supports it.
Capability aliases do not imply a work role and are recognized only as explicit session
declarations or answers to the capability question.

## Effect

This amendment changes session input vocabulary, queue-v1 parser requirements, and user
guidance. It does not change task eligibility, task order, evidence requirements, study
terminology, or any existing immutable task/event value.
