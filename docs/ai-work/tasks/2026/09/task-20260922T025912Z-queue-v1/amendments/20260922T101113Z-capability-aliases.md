# Task amendment — capability alias parsing

**Amendment id:** `20260922T101113Z-capability-aliases`
**Goal amendment:** [`20260922T101113Z-capability-aliases`](../../../../../goals/goal-20260922T025912Z-data-architecture-book/amendments/20260922T101113Z-capability-aliases.md)

Queue v1 must:

1. accept case-insensitive `high`, `leader`, and `master` as canonical `HIGH`;
2. accept case-insensitive `low`, `worker`, `follower`, and `slave` as canonical `LOW`;
3. reject `primary` and `replica` as capability declarations;
4. use only canonical `HIGH`/`LOW` for eligibility and persisted state;
5. preserve the raw declaration as `session_capability_input` when it is available;
6. never emit or recommend legacy `master`/`slave` terms in help or generated prompts;
7. test every alias, mixed case, whitespace handling, invalid terms, and separation from
   `work_role`.

This is a clarification of session input, not a change to the task's state machine or
published eligibility rules.
