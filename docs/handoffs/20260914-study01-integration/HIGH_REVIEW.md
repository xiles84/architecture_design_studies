# HIGH review — Study 01 v3 integration

## TL;DR

EH-01's implementation and local integration are accepted. The final analysis is concise,
its detailed comparisons live in signed discussions, and all requested future-study
requirements are canonical repository rules. No scientific correction or new benchmark
is required for this integration. Only this review's documentation merge remains,
waiting on the active Study 03 run under [EH-02](FINALIZATION_HANDOFF.md).

Reviewer: GPT-6 via Codex desktop, HIGH role, 2026-09-15; selected effort not exposed.
Reviewed closeout: `db77fe60b5bfe78b29ec04db9150f33cc1a0bb49`.
Current-main snapshot: `c7d4067e1c870602e9f794d9067c8163018523dd`.
Environment of saved measurements: `host-zenbook-ux5406sa`.
This is a validation record of existing artifacts, not a new signed interpretation of
the measurements or an independent experimental replication.

## Integration and preservation

The following checks returned exit code 0:

```text
git diff --quiet 0efd760 db77fe6 -- studies/01-charity-tree infra
git diff --quiet 025b72f db77fe6 -- studies/02-ticket-booking studies/03-reserved-seating platform
git diff --quiet db77fe6 c7d4067 -- studies/01-charity-tree infra platform AGENTS.md docs
git merge-base --is-ancestor db77fe6 c7d4067
git diff --check 025b72f db77fe6 -- AGENTS.md CONTEXT.md LESSONS_LEARNED.md docs infra
```

Both `repo/worktree-to-main-workflow` and `study-01/v3-integrated` are annotated tags
targeting `db77fe6`. `study-01/v3-enhancements` still targets
`6bf64851e0943da4b210cdada005bbde43b668f0`. Both September 12 analysts' files exist and
have empty diffs against `study-01/v2-before-enhancements`. No published analysis,
discussion, result or generated report was changed in this review.

The later main commits belong to Study 03: AM-02, calibration, matrix status and failure
diagnoses. They change no Study 01 source, shared infrastructure or platform. They were
incorporated into this task's isolated worktree, retaining the other task's contributions.
This review does not decide Study 03's open escalations or validate its ongoing matrix.

## Artifact and numerical verification

Read-only PowerShell file inspection recomputed the reporter's SHA-256 digest algorithm:
hash each JSON's bytes; sort relative slash-normalized paths ordinally; hash the lines
`path hash` with LF terminators; retain the first 16 hex characters. Each digest equals
the report, signed analysis and EH-01 register. All cells name the same clean producing
commit, image ID and environment as their own manifest; all run tags are annotated and
resolve to those manifest commits. All four reports index the final analysis.

| Run | Cells | Recomputed digest | Recorded invalid cells |
|---|---:|---|---:|
| `20260913T124917Z-v3` | 14 | `d2d1cc8380b381f2` | 0 |
| `20260913T125342Z-v3` | 100 | `0bc5900a6d7a5d3a` | 0 |
| `20260913T172624Z-v3` | 135 | `bd95d304cb39e2de` | 6 D9 negative controls |
| `20260914T104721Z-v3` | 12 | `9ead901cb74c4183` | 0 |

The six D9 concurrent-insert cells each record a real cache mismatch, including warmup
and post-write audits, across three trials on each YugabyteDB topology. All six matched
D10 trials actually ran the cache audit and recorded zero mismatches. The JSON field is
`rollup_audit`, which includes cache checks; an absent generic `audit` field is not a
passing audit. The generated report excludes the six invalid cells from rate summaries.

Headline medians recomputed from saved JSON agree with the final analysis's rounding:

| Check | Recomputed value |
|---|---|
| Five-load D2 / D3 read scores | 5,071.75 / 8,897.79; ratio 1.75 |
| D14 / D15 charity sums | 759.39 / 2,495.97 operations/s |
| Longer-history D3 / D6 final donor reads | 27,497.76 / 6,737.78 operations/s |
| Longer-history D3 / D6 global sums | 27.43 / 5.06 operations/s |
| Matched 256 MiB D3 / D6 final donor reads | 4,578.18 / 16,034.40 operations/s |
| Matched 3 GiB D3 / D6 final donor reads | 25,222.57 / 15,415.07 operations/s |
| Hotspot D5 / D4 / D16 completed inserts | 433.31 / 246.64 / 241.34 operations/s |
| YB one-node D3 / D8 inserts | 356.74 / 448.54 operations/s |
| YB three-node D3 / D8 donor erasure | 122.29 / 136.02 operations/s |
| Equal-budget one-node / cluster read scores | 1,721.99 / 726.79; ratio 2.37 |
| Equal-budget one-node / cluster inserts | 1,140.82 / 270.93 operations/s; ratio 4.21 |

The twelve matched-memory cells have one identical initial dataset summary, 168 passing
initial checks, 1,512 passing post-mutation checks and zero growth-read errors. Local
node-stop trials reconcile acknowledged counts, reject 1,247/1,251/1,251 out of 3,000
requests, and record successful scheduled-response p99 of 15.016/15.055/15.020 seconds.
These observations support the existing qualified wording, including rejected demand,
the single cluster query endpoint, and the limits of a laptop deployment.

## Owner requirements and report structure

- AGENTS.md requires calculated node/client budgets and measured validation, including
  equal-total controls. Methodology section 4 specifies `C/N`, `M/N`, reserves,
  connection-pool accounting, client calibration and endpoint distribution. This is
  reproducible sizing guidance; it does not claim a CPU formula finds an optimal load.
- Applicable rollup, rolldown and embedding controls include maintenance, storage and
  correctness costs. Every multi-node study includes colocated/non-colocated scenarios,
  and every concurrency study includes optimistic/pessimistic strategies. Protocols
  must disclose unsupported or missing coverage rather than count it as completed.
- These requirements apply to new studies and future plans for existing studies;
  historical measurements are not relabelled as having tested every new dimension.
- Methodology 11b and the discussion template place detailed mechanisms and analyst
  responses outside the final decision document. The existing final analysis has
  751 whitespace-delimited words before its input register and links four signed
  companions. Earlier analysts are cited without implying new participation.
- AGENTS.md requires an owned worktree, meaningful commits, immutable annotated tags,
  final local-main integration, and explicit HIGH/LOW handoff and escalation roles.

Relevant local file links in the reviewed policy, living documents, final analysis,
discussion companions and handoff documents resolve (97 links across 17 files).
Prior containerized tests, vet,
shell validation and artifact checks remain applicable because the tested source is
unchanged; see the [prior integration review](../../../studies/01-charity-tree/reports/20260914-v3-integration-validation.md).
No benchmark, image build or report regeneration was run for this review.

## Closeout findings and limits

The LOW receipt stopped at intended final-tag wording, and living context retained
"next HIGH" plus an unfinished sentence. Git establishes that the actual merge/tag
steps succeeded. This HIGH checkpoint corrects the living context and records those
facts here, preserving the historical receipt. These are documentation fixes, not a
reason to rerun measurements or rewrite an analyst's conclusions.

An attempted second-agent audit returned no findings because it reached a usage limit.
This acceptance rests on the primary review above, not an independent second sign-off.
Real-network/physical-host validation and broader cache-payload mutation audits remain
explicitly unmeasured follow-ups in the existing analysis and protocol.

At 2026-09-15 09:43 UTC, Study 03 held `ads-run-lock` for `20260915T002411Z`. Its live
labels name branch `main` and worktree `C:/extra/code/architecture_design_studies`,
which also had uncommitted matrix outputs. The owner subsequently confirmed that another
AI is active and that both agents must use separate worktrees. This task uses
`.worktrees/study01-v3`; its read-only helper used `.worktrees/study01-review-audit`.
The active main-based run is left in place until its owning task finishes; moving it
mid-run is outside this handoff. Its run and files remain untouched.
The existing enhancements and rules are already on main; only this new review checkpoint
awaits integration. **Next: LOW model / low effort — execute EH-02.** If its mapped
preservation checks pass, LOW may complete the task without another HIGH review.

**Later HIGH check:** main then advanced to `7a9aea2`, changing only context to reserve
its checkout explicitly until Study 03 step 11 is committed. The task worktree retains
that incoming entry. EH-02 revision 2 also requires that reservation to end; a brief
idle gap between matrix runs does not permit a main update. The initial review remains
at `study-01/v3-high-review`; the amended checkpoint is `study-01/v3-high-review-r2`.
