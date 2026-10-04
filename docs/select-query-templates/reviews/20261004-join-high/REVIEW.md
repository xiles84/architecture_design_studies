# HIGH review — join/CTE submission: changes requested

Reviewer: GPT-6 / Codex desktop / HIGH; effort unavailable, recorded unknown.
Independent of the LOW author (`deepseek-flash` / opencode CLI).
Reviewed result commit `89479eb6`, submission `937388a0`, producing source
`7e148174b5211225f6b4f243fc02e0de417a0586`. Reviewed D1–D8 first.

**Verdict: request changes. Do not approve, integrate measurements, or promote
book claims yet.** The existing evidence is retained, not discarded. The
task's retained, clean worktree was taken over for this review after submission;
no active executor claim was displaced. No database was started by this review.

## Accepted work and decisions

The independent [audit](CHECKS.md), implemented in [Go](audit.go) and run in
the pinned Go Podman image, verifies all six original result digests against
their reports, 546 valid cells, five samples/cell and their medians, six distinct
order seeds, 90 unique matrix cells plus the actual first-cell repeat, both
controls firing in each run, and materialized CTE scans. All ten producing tags
(six matrix and four feasibility) are annotated and pinned in [the inventory](TAGS-BEFORE.txt).

D1's separate catalogue/schema, D2's shuffled rank and fixed indexes, D3's
deterministic generated SQL, D5's explicit unsupported outcome, D6's different
fixtures for the two controls and D7's consumed checksum are sensible local
implementation choices. Keep these. D4's actual code correctly tests `CTE Scan`;
its old Decision Log wording about `Materialize` is stale, so the next attempt
must explain the correction without editing the prior author's log.
D8's claimed resource rejection is **not accepted**.

A fresh, separately named benchmark image passed `go vet ./...`, `go test ./...`
and compilation in Podman (root: no tests; harness and selectbench: PASS).
Image: `localhost/charitybench-high-review:20261004`, ID
`6c79d9035942759fbcad6ed642f4770d35fe1e15453158acfb836153ffad4935`.
This source check does not validate missing database arms.

## R1 — replace the unsupported resource rejection with actual coverage

Every run saves `/proc/meminfo` with **16,109,780 kB (~15.36 GiB)** and CPU
online `0-7` (8 CPUs). The Podman machine configuration separately says 2048 MiB
and 4 CPUs; that configuration is not the effective WSL allocation. The report
claims its 3 GiB database plus 2 GiB client cannot fit, but the saved runtime
evidence contradicts that conclusion. The required 3-node subset, equal-total
comparison and placement comparison were skipped without testing feasibility.
An unimplemented balanced-endpoint client is implementation work, not proof
that the host cannot run the arm.

Execute the resource and topology correction in [the handoff](CORRECTIVE-HANDOFF.md).
Include client capacity, reserves and actual limits, balanced endpoint counts,
per-node throttling, RF3 and actual colocation evidence. A genuine feature or
resource failure remains a gap, supported by a reported check. Do not reclassify
an unattempted arm as unsupported. The saved single-node before/after CPU stats
also exist; "no per-node throttling because one node" is incorrect.

## R2 — restore the approved domain and workload contract

The dated HIGH handoff explicitly extends charity/person/donation and limits
customers/orders to theory. The implemented `selectjoin.customers/facts`
fixture is a synthetic parent/fact experiment, with no explicit implemented
Study 01 domain mapping. The original protocol §9 excludes numeric claims for
the teaching domain. Keep this evidence as the original microfixture; produce
the required domain measurements under new run IDs, rather than silently
relabel its raw results as donor/donation evidence.

The dated handoff also requires isolated **and matched blended read/write**
phases. `client_workers = 1 sequential reader; 0 writers` documents that the
second phase is absent. A note calling it a gap does not satisfy a mandatory
arm. Keep the same fixed correct writer across syntax pairs; do not add a
concurrency-strategy comparison. Add the live empty-input and filtered-child
LEFT JOIN equivalence checks required by the handoff (the current live gate
tests ON only, while an empty input is tested only in Go).

## R3 — retain the engine evidence required to interpret plans

`pgxQuerier.Explain` currently sends `EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)`
to both engines. It omits PostgreSQL VERBOSE output and YugabyteDB DIST/RPC
evidence required by methodology §9. "RPC not parsed" is misleading: the
DIST counters were not requested. Capture the required engine-specific plans,
including consumed-expression output/loops and the supported storage/RPC
counters. Verify token volatility against the pinned engines, not only that a
constant expression can be evaluated. Gate both forced policies against their
actual behavior; a retained CTE under a requested NOT MATERIALIZED policy is
not automatically a valid inlining comparison.

## R4 — repair reproducibility and qualify interpretation

Using the **original pg1 producing image** to re-render the unchanged input
produced [this report](rerender-pg1.md) advertising digest
`7c8795647185b68a3053b1a265969a80e7a966ed65e82befb50f89b73d824b6c`.
The actual saved input and original report digest is
`070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac`.
`runSelectJoinReport` hashes `append(data, '\n')`, adding a second newline to
the already newline-terminated input. Fix it and prove rendering twice retains
the exact input digest. Separate client image/name/digest from database image
name/digest: the current provenance labels a client image as `EngineImage`
alongside a database image digest. The replicated renderer must reject mixed
environment/fixture/settings input or partition it explicitly; its current key
contains only engine, topology and axis.

Publish a **new signed successor**, preserving the original analysis. In
particular:

- Compare P4's actual pair, NOT MATERIALIZED / MATERIALIZED. Default folds at
  one reference and is a different baseline there. At PG/100% consumer the
  paired n=1 ratios are **0.9080–0.9687** (near equality), n=2
  **1.6957–2.1235**, and n=16 **4.0048–4.3493**. At PG/1% consumer the n=16
  ratios are **0.6950–0.7868**. Two consumer levels establish endpoint
  observations, not the claimed five-level crossover in consumer selection;
  the near-equal reference-count endpoint needs an unresolved verdict when
  variability prevents bounding a crossing. No new sweep is required merely
  to make a crossing appear.
- P12 is the original **50–90% convergence prediction**, not a generic promise
  to retain outcomes. Give it its own verdict. P5 is a semantic control only.
- Position last/first ranges from **0.6056 to 1.3470**. A pg1 1.52x pair is
  not arithmetically "inside 0.80 drift". Evaluate paired effects against
  samples and across-run variability rather than treating drift as a universal
  significance threshold. In PG run 3, **88/91 cells exceed 20% sample spread**.
- The matched joined-output policy ratio at YB threshold 180,000 is
  **1.0172, 0.8844, 0.9442**: mixed sign. Do not claim forcing materialization
  is slower at every level, or use pg1 alone as a replicated range. The pg
  threshold-20 matched ratio spans **73.1056–92.8513** across runs.
- The index/storage-filter distinction is observed. It does not by itself
  prove that one fact explains every timing difference. Distinguish shape
  signatures, complete plans, counters, and causal inference.

## Route and completion

The next executor claims **this same task and branch** after the queue records
`changes_requested`, creates a new attempt, reads this review and the corrective
handoff, and preserves all existing evidence, signed files and tags. No remote
sync. Resubmit for independent HIGH review. Synthesis and book publication stay
dependency-gated. The completion tag `study-01/select-v2-join-cte` remains absent
until genuine approval/integration; this review receives a separate checkpoint tag.

NEXT MODEL: LOW
