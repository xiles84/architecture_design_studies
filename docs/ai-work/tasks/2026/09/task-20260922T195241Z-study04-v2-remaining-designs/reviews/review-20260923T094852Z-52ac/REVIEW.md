# Review — Study 04 v2 remaining designs

**Task:** `task-20260922T195241Z-study04-v2-remaining-designs`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration.** The task used the "record each unbuilt design as a
coverage gap with the reason" branch the acceptance explicitly allows, and delivered one real,
gate-clean design and comparison.

## Decision Log reviewed

1. Build the SQL-only designs, record the rest as gaps — accept; the acceptance authorises exactly
   this, and rushing seven harness changes would risk the correctness standard.
2. `d2_doc_on_parent` is SQL-only — accept; verified by a dev check and the reported run passing the
   gate.
3. `d4` is not SQL-only — accept; the document write path carries the whole new document.
4. `d3`/`h1`/`s1`/`s2` need new Kinds — accept.
5. `y1` needs verified colocation DDL and placement recording — accept.
6. Equal-total needs a runner flag — accept.
7. Report the d2 negative result rather than re-running until it agrees — accept; this is the
   discipline the repo asks for.

## Independent checks

- `failed_cells: []` in the run manifest; both cells `gate.passed = true` (10 checks).
- Recomputed from result JSON: publication 865.9 vs 858.8; metadata 1651.0 vs 1672.7; r05 22898 vs
  14433. Matches the analysis.
- New design is registered (`designs.go`, `ALL_DESIGNS`) and its SQL ships in `sql/d2_doc_on_parent/`.
- Run tag `run/04-configuration-portal/20260923T0955Z-d2-vs-d1` exists on the recorded commit.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Build and run the eight unbuilt designs | **partly** — d2 built and run; seven recorded as gaps with reasons |
| Record each unbuilt design as a coverage gap with the reason | met |
| Equal-total-resource arm | gap, recorded |
| YugabyteDB coverage | gap, recorded |
| Placement with a recorded distribution | gap, recorded (the acceptance itself says "no placement record means gap") |
| Gate passes; one matrix at a time | met |

## Conditions

The seven gaps and the equal-total/placement gaps must not be read as completed comparisons. They
are recorded in the task `RESULT.md` and in `CONTEXT.md`; the evidence registry's coverage matrix
already carries the Study 04 cardinality cell as `planned` and must stay that way until the harness
work lands. No new task is created here: the gaps are harness work owned by future Study 04 sessions.

## Integration

Approve and integrate; tag `study-04/v2-remaining-designs`.
