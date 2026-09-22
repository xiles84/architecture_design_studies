# Owner request

Captured verbatim from the 2026-09-22 implementation request. The surrounding chat
already established that the current session capability is HIGH and that “typist” means
Typst (`.typ`).

---

PLEASE IMPLEMENT THIS PLAN:
# Evidence-backed data architecture book and AI work pipeline

## Summary

Current baseline is clean local and remote-tracking `main` at `f27cb14`. The repository has enough evidence for a strong beginner-to-master reference on data architecture, concurrency, caching, reporting, and database topology—but not yet a general software-design encyclopedia.

The immediate problems are:

- Reader-facing documentation is stale: the root index omits Studies 03–05; several study READMEs still say “planning” or “in progress.”
- There is no book-level taxonomy, scenario guide, or claim-to-evidence index.
- `CONTEXT.md` mixes current state with history and contains contradictory/duplicate status records.
- Major native datastore families are not yet compared. JSONB is not evidence for a native document database, and Redis was tested only as a cache.
- Topology evidence is laptop-bound; Studies 04–05 are single-trial and PostgreSQL-only.
- The HIGH/LOW workflow exists, but there is no atomic task claim, lease, recovery protocol, structured receipt, or parent/child correlation.

Chosen product decisions:

- Scope v1 as a measured data architecture and state-design reference.
- Define a major model as a design family and a minor model as a controlled variant.
- Treat technology, topology, workload, and design family as separate axes.
- Author the assembled book in Typst and build its PDF in Podman.
- Keep operational material in Markdown and machine-readable queue/evidence state in JSON.
- Publish current evidence with explicit confidence and gaps; do not wait for every future study.

## AI queue and durable history

Create `docs/ai-work/` as the authoritative archive:

```text
docs/ai-work/
  WORKFLOW.md
  SCHEMA.md
  goals/<goal-id>/
    REQUEST.md
    GOAL.md
    amendments/
  tasks/<yyyy>/<mm>/<task-id>/
    task.json
    BRIEF.md
    STATUS.md
    events/<timestamp>-<event-id>.json
    attempts/<attempt-id>/{PROGRESS.md,RESULT.md,ESCALATIONS/}
    reviews/<review-id>/REVIEW.md
```

Rules:

- `REQUEST.md` preserves this owner request verbatim; `GOAL.md` records HIGH’s interpretation and success criteria.
- Every new repository chat asks whether it is a HIGH or LOW session unless the user already stated it.
- Record `session_capability` separately from `work_role` (`planner`, `executor`, `reviewer`, `analyst`, or `integrator`).
- HIGH may claim any task. LOW may claim only tasks whose `minimum_capability` is LOW.
- Each task records stable task, goal, root, parent, and originating-handoff IDs, allowing every report, escalation, correction, and review to be traced back to the HIGH request that created it.
- Queue ordering is: dependencies ready, capability eligible, `not_before` reached, priority descending, creation time ascending, task ID ascending.
- Events are immutable individual files linked by `previous_event_id`; corrections create later events rather than rewriting history.
- `CONTEXT.md` becomes a compact current-state index, never the queue’s state authority.

Implement a Go CLI, executed from a pinned Podman image, with:

```text
queue session-start
queue list
queue next
queue publish
queue claim
queue guard
queue heartbeat
queue checkpoint
queue yield
queue escalate
queue submit
queue review
queue approve
queue request-changes
queue integrate
queue complete
queue recover
queue audit
queue review-candidates --since 7d
```

Use Git compare-and-swap refs for live coordination:

- `refs/ads-queue/live/<task-id>` holds the latest claim/lease JSON blob.
- `refs/ads-queue/locks/main-integration` serializes updates to local `main`.
- `git update-ref <ref> <new> <observed-old>` makes claims atomic across all worktrees sharing this repository.
- The live ref is host-local coordination; committed task events are the portable archive and reconstruct state after a fresh clone.
- This guarantees exclusivity within this shared local Git repository. Multi-host clones would require a later central coordinator.

Lifecycle:

```text
proposed → ready → claimed → in_progress
                         ↘ yielded → ready
                         ↘ blocked_high → ready
                         ↘ awaiting_review
                              ↘ changes_requested → ready
                              ↘ approved_for_integration → completed
```

Use a two-hour default lease, heartbeat at least every 15 minutes, and an explicit extension before long benchmarks. Recovery requires an expired lease, a ten-minute grace period, no task-owned benchmark lock, a worktree/branch inspection, and a compare-and-swap claim epoch change. Recovery always reuses and preserves the existing worktree; it never resets or discards unfinished files.

Migrate existing handoffs by reference without rewriting them. Import completed EH-01, EH-02, Studies 03–05, and create a separate ready HIGH task for the still-open Study 02/03 v2 validation. Coverage gaps become `proposed`, not automatically authorized benchmark work.

## Typst book and evidence contract

Create this reader layer:

```text
book/
  main.typ
  chapters/
    how-to-choose.typ
    major-families.typ
    minor-variants.typ
    cross-family-comparisons.typ
    topologies.typ
    scenarios.typ
    evidence-and-reproduction.typ
  concepts/
    information-placement.typ
    indexes-and-hot-rows.typ
    concurrency-control.typ
    derived-state-and-history.typ
    expiry-and-clock-authority.typ
    cache-consistency.typ
    distributed-placement.typ
    reading-benchmarks.typ
  evidence/claims.json
  assets/
  Containerfile
  README.md
  dist/data-architecture-reference.pdf
  dist/build-manifest.json
```

Pin Typst and all fonts in the container. The build manifest records source commit, dirty state, Typst image/version/digest, evidence-registry digest, release date, PDF hash, and build command.

The taxonomy is fixed as:

- Major families: normalized facts, rolldown/duplicated data, rollups/materialized aggregates, embedded documents, append-only history/snapshots, and external read copies/caches.
- Minor variants: index selection, full versus bounded embedding, trigger versus application maintenance, optimistic versus pessimistic concurrency, locks/CAS/isolation, pre-created versus on-demand inventory, cache-aside versus write-through, and strict versus relaxed freshness.
- Topology: node count, replication, placement, endpoint routing, resource framing, network conditions, and failure behavior.
- Technology: PostgreSQL, YugabyteDB, Redis, and future engines.
- Workload: volume, cardinality, skew, contention, read/write mix, history length, failures, and consistency requirement.

Every family, variant, topology, and scenario follows one progressive structure:

1. Plain-language quick choice.
2. Where it thrives and perishes.
3. Common scenarios.
4. Reusable mechanism explanation.
5. Read, write, storage, correctness, and operational costs.
6. Direct evidence with run, digest, tag, environment, and analysis.
7. Unmeasured boundaries and transfer risks.
8. Reproduction links.

`claims.json` records each recommendation with one of:

- `repeated_controlled`
- `single_run_directional`
- `mechanism_supported`
- `analogy`
- `gap`
- `invalid`

Every numeric statement must identify its controlled comparison, run, digest, producing tag, environment, topology, report, and signed analysis. Cross-study numeric rankings are forbidden unless workload and semantics are explicitly demonstrated comparable. Scenario examples such as supermarkets, e-commerce, auctions, and top-10 dashboards must say whether they are directly measured or an analogy.

Rewrite the living navigation documents while preserving immutable evidence:

- Update the root README to index all five studies and the book.
- Correct stale Study 02, 03, and 05 status sections.
- Archive the existing long `CONTEXT.md` snapshot, then replace it with one current state per study and links to task history.
- Remove duplicate live receipts only after preserving them in the archive.
- Reorder methodology headings without changing their rules.
- Never rewrite raw results, generated report history, signed analyses, discussions, prior handoffs, or producing tags.

## Ordered task queue

| Order | Task | Capability | Result |
|---|---|---:|---|
| 0 | Publish this request and handoff from an isolated worktree | HIGH | Goal archive, bootstrap task records, commit, annotated handoff tag, merge to local `main` |
| 1 | Implement `queue` v1 and its Podman wrapper | LOW | Atomic claims, leases, recovery, role filtering, event validation, integration lock |
| 2 | Import historical tasks and reconcile public statuses | LOW | Queryable history, concise current context, accurate root/study navigation |
| 3 | Build the evidence registry and validation tools | LOW | Claims extracted with provenance and confidence; no new interpretation |
| 4a | Validate Study 02/03 v2 evidence | HIGH | Signed review or targeted correction tasks |
| 4b | Independently review Study 04 | HIGH | Provisional claims accepted, narrowed, or converted to follow-up tasks |
| 4c | Independently review Study 05 | HIGH | Same, without rewriting its existing analysis |
| 5 | Build the Typst toolchain and book shell | LOW | Reproducible container build, navigation, assets, validation, draft PDF |
| 6 | Write the synthesis chapters | HIGH | Major/minor/topology/scenario guidance grounded in reviewed claims |
| 7 | Independent release review | HIGH | Claim audit, beginner/master readability review, final PDF and `repo/data-architecture-book-v1` tag |

Tasks 4a–4c may run in parallel. Database measurements remain serialized by `ads-run-lock`.

Future scientific tasks enter as `proposed` and require a HIGH protocol before execution:

1. Complete Study 04’s eight missing designs, randomized repeated trials, cardinality crossover, cadence/churn, YugabyteDB, equal-total resources, and placement pair.
2. Complete Study 05’s sustained real-TTL churn, medium scale, repeated trials, equal-total cache accounting, YugabyteDB, cluster placement, and open-loop workload.
3. Run a real topology study with balanced endpoints, independent hosts, injected latency, node/endpoint loss, per-node and equal-total budgets.
4. Add a native cross-major comparison using the configuration domain.
5. Add analytical/search/read-model comparisons for dashboards and top-N workloads.
6. Add a true hierarchy-representation study: adjacency list, materialized path, closure table, nested sets, and bounded embedding.

These studies improve later editions but do not block an honest v1 whose gaps are visible.

## Verification and acceptance

Queue validation must cover:

- Twenty parallel claimers with exactly one winner.
- HIGH/LOW filtering and HIGH consumption of LOW-eligible tasks.
- Dependency, priority, FIFO, `not_before`, and seven-day queries.
- Crashes before worktree creation, with dirty files, after commits, and before submission.
- Expired-lease recovery and rejection of the old claimant’s next guard.
- No recovery while the task owns the benchmark lock.
- Exactly one concurrent local-main integrator.
- Invalid event chains, malformed identities/timestamps, spec-digest changes, and orphan refs.
- Fresh-clone reconstruction from committed archives.
- Native Windows Git and WSL Git operating on the same common Git directory.

Book validation must ensure:

- The Typst book builds entirely in Podman.
- Root navigation lists every study and no public status contradicts current evidence.
- Every numeric recommendation resolves to current evidence and never cites a failed/invalid cell as a winner.
- Every scenario is labelled direct evidence or analogy.
- Every topology claim states nodes, replication, resource framing, endpoint distribution, placement evidence, and network limitation.
- All local links, run tags, report digests, analyses, and supersession links resolve.
- The PDF and build manifest identify the exact source and evidence digest.
- Current work remains isolated, committed, tagged, merged locally, and never pushed or pulled by an AI.

Because Plan Mode prevented writing the archival handoff, keep the current HIGH role for the short publication checkpoint in task 0. Its committed handoff will then route the queue implementation to LOW.

---
