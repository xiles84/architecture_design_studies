# Study protocol — <study id, short name>

| Field | Value |
|---|---|
| Study | `studies/<NN-name>/` |
| Protocol revision | <n> |
| Planner | <actual model, known effort, tool, date; HIGH> |
| Starting main revision | <verified commit> |
| Status | <draft / ready / amended / executed-awaiting-review> |
| Required tag on integration | `study-NN/v0-handoff` |
| Numbering | <Study NN is X; the next free number is Y> |

A protocol is written **before** any code that measures. It fixes the question, the
prediction, the comparison forms, the controls, the workload levels and the boundaries, so
that a result can be *wrong* rather than merely reported. The repository's method is
[`docs/methodology.md`](../methodology.md); this template is its per-study form.

## 1. Question and invariant

- The real-world scenario, in one paragraph a practitioner would recognise.
- The architectural question, stated as a decision, not a topic.
- The invariant(s) that must hold, each with the failure it forbids. If the study has no
  invariant, say so and state what correctness instead means.

## 2. Prediction (recorded before any measurement)

**This section is written before the first run and is not edited afterwards.** A prediction
added after a result is a retrofit; if one must be added, mark it `unrecorded`.

| Id | Comparison | Expected direction | Expected mechanism | Expected crossover / regime | Falsifier | Recorded at commit |
|---|---|---|---|---|---|---|
| P1 | <design A vs B> | <A > B / no difference> | <why> | <the level at which the ordering should reverse, if any> | <the observation that would refute it> | <commit> |

Each prediction gets a verdict in the final analysis: `confirmed`, `falsified` or
`unresolved`. A falsified prediction is a result and is reported as one — Study 04's
`d2` metadata-amplification prediction that did not appear is the model.

## 3. Designs, pairs and the three forms of evidence

Every comparison is one of three forms. The form decides what may be claimed from it.

| Comparison | Form (`isolate` / `interaction` / `package`) | What changes | Attribution rule |
|---|---|---|---|
| <D1 → D2> | isolate | exactly one decision, SQL/indexes otherwise matched | the delta may be attributed to that decision |
| <A × B> | interaction | two decisions crossed | neither main effect may be attributed alone; report the interaction |
| <package vs baseline> | package | several decisions at once | **nothing may be attributed to a component unless an isolation pair supports it**; otherwise record `attribution_withheld` |

A package is allowed and often realistic. It is never a licence to name one component as
the cause. If the package's headline has no supporting isolation pair, say so in the
statement and in the limits.

## 4. Correctness gate and negative controls

| Control | Deliberately wrong how | Invariant it must violate | Must fire in | Read-path control? |
|---|---|---|---|---|
| <X0> | <the wrong mechanism> | <INV-n> | <phase / topology> | <yes/no> |

Rules carried from methodology 5 and 5a:

- The gate compares against values computed independently from the generated dataset, never
  one query against another. A failing cell reports no timing.
- Every invariant has at least one deliberately wrong design that is expected to break it.
- **Read-only phases need a read-path control** — a deliberately stale or wrong copy that is
  seen to return the wrong answer. A control that cannot fire under a read-only workload is
  recorded as `control_absent`, and the absence of violations in the other designs is then
  *not evidence* of their correctness.
- Invariants are two-sided: a "never over/under" check is paired with the opposite direction
  and with a reconciliation of what the client was told against what the store holds.

## 5. Workload dimensions and levels

| Dimension | Levels | Why these levels |
|---|---|---|
| <cardinality / writers / contention / hit rate / …> | <e.g. 1, 4, 16, 64> | <the crossover or boundary the sweep is meant to find> |

Rules:

- Prefer levels chosen to locate a **crossover** or a **boundary** over levels chosen to
  produce a larger throughput number.
- A sweep claims a crossing only with enough levels to bound it; a two-point comparison is a
  direction, not a curve, and must be labelled as such.
- Change the shape, hold the volume: a regime that removes a limitation keeps the total work
  comparable (methodology 8a).

## 6. Regimes that remove each limitation

| Design | Its limitation | Regime that removes it | Designs re-run |
|---|---|---|---|

The regime map is complete, including limitations whose acceptable regime is already covered
by another experiment or which no benchmark can establish.

## 7. Replication, order and variability

| Field | Plan |
|---|---|
| Whole runs per cell | <n ≥ 3 for a conclusion that argues a small difference> |
| Trials per cell within a run | <n> |
| Design order | <randomized per run / interleaved / AB-BA / fixed, and why> |
| Position control | <same design first and last in a phase? yes/no> |
| Across-run variability | <how it is reported; median plus across-run spread> |
| Tier | `run-level` (same machine) or `hardware-independent` (≥2 hosts) |

A same-machine repeat is `run-level` replication and is never described as
hardware-independent; the environment hazards of the host still apply.

## 8. Coverage of the owner's required comparisons

Map each requirement to a named pair, or record a justified non-applicability or an explicit
unsupported gap. Absence of a scenario is never evidence that a mechanism has no effect.

- cluster and client sizing (per-node and equal-total arms, pools, endpoint distribution)
- rollup, rolldown and embedding
- colocated and non-colocated data (multi-node studies only)
- optimistic and pessimistic concurrency (concurrency studies)
- any other decision this study varies

## 9. Boundaries and non-applicability

State, before running, what this study will **not** establish: the workloads, engines,
topologies, scales and failure modes outside its scope. These become the study's boundary
section, not a footnote.

## 10. Acceptance

- Every prediction carries a verdict; every control's firing is recorded per topology.
- Every claim is tagged with its evidence form and, if a package, its attribution status.
- Crossover and boundary results state the dimension, the levels and the crossing interval.
- Failed cells are named and never reported as winners.
- Every result names its environment, pinned images, run tag and inputs digest.
- No number enters the book before it resolves to a registered claim.
