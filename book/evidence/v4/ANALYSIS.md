---
analysis_id: book-evidence-v4-ingest--deepseek-flash--2026-09-24
run_id: 20260923T1215Z-res2-placement-yb3
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task book/evidence-v4."
analyzed_at: 2026-09-24
repo_commit: 01eb8d332337ca0275c68afe6a5f1e8d3b97da7a
supersedes:
status: current
headline: Retires the two placement/endpoint gaps that v3 carried forward unchanged after the Study 05 runs had already closed their clauses, and adds claim lifecycle metadata so a later claim cannot silently coexist with a stale gap.
---

# Analysis — book evidence v4 (placement and endpoint gap retirement) — deepseek-flash

> Deliverable of task `book/evidence-v4`, implementing finding **R04** of the external review of the
> *Edition 1* draft. This is an ingestion and retirement analysis: it **measures nothing**, adds no
> number, and edits no earlier signed artefact. It rewrites two *gap* claims whose text had become
> false, and leaves every numeric claim exactly as v3 left it.

## TL;DR

- The active registry moves from `book/evidence/v3/claims.json` to `book/evidence/v4/claims.json`.
  v4 carries all 27 still-active v3 claims forward content-identical and changes only the gap layer.
- `v2-gap-05-colocation-unverified` and `v2-gap-07-real-network-balanced-endpoints` are **retired**.
  Neither was wrong when written; both became partly false when the 2026-09-23 Study 05 runs landed
  in v3, and v3 did not revisit them. Their still-true clauses are carried, narrowed, by
  `v4-gap-01-physical-placement-unverified` and `v4-gap-02-endpoint-distribution-partial`.
- `v3-gap-01-study05-remaining-regimes` is marked `partially_superseded`: its placement-distribution
  clause is now owned by the two successors, and its `remaining_dimensions` name only what is still
  open. This removes a duplication that would otherwise let two active gap claims disagree.
- Claim lifecycle metadata (`status`, `gap_kind`, `superseded_by`, `closed_dimensions`,
  `remaining_dimensions`) is added, and `tools/evidence` enforces it. The intent is that the next
  superseding claim cannot land beside the gap it closes without the validator noticing.

## 1. What was actually wrong

Two independent things, and it matters that they are separate.

**A stale contradiction.** v3 carried `v2-gap-05` and `v2-gap-07` forward "byte-identical in content
to v2" (v3 `SUPERSESSION_LEDGER.md`, verbatim). Both had clauses that the same v3 package then
falsified:

- `v2-gap-05`: "Study 04 y1/y2 and Study 05's placement pair were never executed."
  `v3-06-colocated-vs-noncolocated-locality` records the pair as executed.
- `v2-gap-07`: "balanced client access across cluster endpoints is not proven."
  `v3-05-yb-cluster3-strict-relaxed` states the client "spreads operations over all three endpoints".

**An inconsistency inside the package.** v3's own `coverage.json` was updated for exactly this, in
its `real-network-and-balanced-endpoints` dimension: "Study 05's 3-node cells spread operations over
three endpoints (connection_nodes=3) but record only the count, not per-endpoint operation
distribution." So v3's coverage matrix and v3's claims registry disagreed with each other, and the
book printed the claims registry. That is the sharper form of the finding: it is not that evidence
was missing, it is that two files in one package said opposite things.

## 2. What each retired clause becomes

| Clause | Disposition |
|---|---|
| `v2-gap-05` "Study 05's placement pair was never executed" | **closed** → moved to the numeric claim `v3-06`; recorded in the retired entry's `closed_dimensions` |
| `v2-gap-05` "no study verifies actual node placement" | **retained** in `v4-gap-01`, with the reason it is still true: the runner produced no tablet/leader distribution, so the layout labels are intent, not observed placement |
| `v2-gap-07` "balanced client access across cluster endpoints is not proven" | **closed for the 2026-09-23 Study 05 cells** → recorded in the retired entry's `closed_dimensions` |
| `v2-gap-07` "endpoint distribution is not measured" | **retained, narrowed** in `v4-gap-02`: unmeasured for the studies 01–03 multi-node cells, which set one DSN against `yb-n1` |
| `v2-gap-07` "no real network between nodes" | **retained** in `v4-gap-02` |

The `yb-n1` fact is from the harness, not from a report: `studies/01-charity-tree/run-study.sh` and
`studies/02-ticket-booking/run-study.sh` both give the `yb-cluster3` cell the DSN
`postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable`. It is named in the claim's `gap_basis`
so a reader can check it without rerunning anything.

## 3. Anchors and digests

No new measurement is cited, so every anchor in v4 reuses a run that v3 already anchored. The
successor gaps take the newest relevant run as primary and keep the Study 03 placement run as
corroboration.

| Claim | Anchor | Inputs digest | Role | Report |
|---|---|---|---|---|
| `v4-gap-01-physical-placement-unverified` | `20260923T1215Z-res2-placement-yb3` | `a18817d592f767c5` | primary | `studies/05-cache-consistency/reports/20260923T1215Z-res2-placement-yb3.md` |
| `v4-gap-01-physical-placement-unverified` | `20260921T231328Z` | `2f5b619430c2e7a7` | corroborating | `studies/03-reserved-seating/reports/20260921T231328Z.md` |
| `v4-gap-02-endpoint-distribution-partial` | `20260923T1157Z-res2-yb3` | `f2133069805cf82d` | primary | `studies/05-cache-consistency/reports/20260923T1157Z-res2-yb3.md` |
| `v4-gap-02-endpoint-distribution-partial` | `20260921T231328Z` | `2f5b619430c2e7a7` | corroborating | `studies/03-reserved-seating/reports/20260921T231328Z.md` |

Both successor gaps have `support: []` and `strength: "gap"`, and `trials` all zero: a gap claims no
measurement of its own. The two primary runs are cited here, not restated: the numbers they
measured stay in `v3-05` and `v3-06`, which is why no figure appears in either gap statement.

## 4. What this does not do

- It does not close the physical-placement question, and it does not weaken it. Verified placement
  still needs an instrument the current runner does not have.
- It does not touch the numeric claims, the confound register, or any published report. No report,
  result file or signed analysis is edited; this file is a new signed artefact.
- It does not claim that the earlier multi-node studies are invalid. Their ratios remain valid *for
  a client reaching the cluster through one endpoint*, which is what `v4-gap-02` now says.
- It does not restate the endpoint count as a balance. `connection_nodes=3` records that three
  endpoints were configured, and the spread is round-robin by construction; neither is a measured
  per-endpoint distribution.

## 5. Where this analysis is weak

- **The `yb-n1` claim is a source read, not a measurement.** It is checkable in the committed
  runners, but no run recorded per-endpoint operation counts, so "single endpoint" is inferred from
  the DSN configuration plus the harness's connection setup. If the client library ever fans out to
  the other nodes on its own, this wording would be wrong; the honest hedge is that the cell is
  *configured* against one endpoint.
- **v3-gap-01's `remaining_dimensions` are my summary of its existing sentence**, not a new finding.
  I split the sentence into three dimensions; a reader who disagrees with the split should read the
  claim statement, which is unchanged.
- **Lifecycle metadata is only as good as the entry that maintains it.** Nothing here detects a
  *future* contradiction automatically; the validator checks that retired claims are not also
  active and that a gap declares its kind. Detecting "a newer claim closes an older gap's clause"
  remains a human review step, and that is exactly where v3 failed.
