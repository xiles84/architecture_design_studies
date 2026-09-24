# Supersession ledger — evidence registry v3 → v4

**Attribution.** Written by `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI; effort
setting not exposed) on 2026-09-24, task `book/evidence-v4`, implementing finding **R04** of the
external review of the *Edition 1* draft (GPT-5.6 Sol, high effort, 2026-09-24 00:25 −03:00).

**What this ledger is.** The single attributed record of what changed from
`book/evidence/v3/claims.json` (v3) to the active `book/evidence/v4/claims.json` (v4). It exists
because v3 carried two gaps forward from v2 *byte-identical* while adding the evidence that
contradicts them, and the book then printed both. A superseding claim that lands without retiring
the gap it closes is the defect class this ledger and the new lifecycle fields are meant to make
visible instead of silent.

## Frozen predecessors are unchanged

- v1 milestone tag **`repo/book-evidence-registry-v1`**; `book/evidence/claims.json` unchanged.
- v2 milestone tag **`repo/book-evidence-registry-v2`** → commit
  `075dd94672b7fdecdcf1f83e07bb3e7e9b87f764`.
- v3 milestone tag **`repo/book-evidence-registry-v3`** → commit
  `d8ed10fcc0a3a13f02ec832b421e35e9829c4a23`;
  `sha256(book/evidence/v3/claims.json) = 814f40949d162f564ce2a831099727c2f1237a32629f99a2a24dd006db4c47c4`,
  `sha256(v3/schema.json) = 67d0b37f67db6923d2ef1fc26e177aaebc8a2ee05e7a1e877dbcdba43b7d9b8f`,
  `sha256(v3/confounds.json) = ed598863db2010bbf82bc7aabbfc9fc43fbf987f54871411b16e7fdd40830150`,
  `sha256(v3/coverage.json) = 6b1a124dfaefeeec62cd6e2ace387d7332e934b90580b69fc0a41ac57a545916`.
- `git diff repo/book-evidence-registry-v3 -- book/evidence/v3/` is **empty**; v1, v2 and v3 were
  not edited, and no tag was moved, deleted or reused.

## v4 is v3 carried forward

- **Carried forward unchanged:** all 27 still-active v3 claims. Their statements, anchors, limits,
  confounds, support keys and reviews are content-identical to v3. Verified mechanically before
  commit: normalising away the four new lifecycle keys and diffing against v3 reports no change.
- **Retired:** `v2-gap-05-colocation-unverified` and `v2-gap-07-real-network-balanced-endpoints`,
  both appended to `retired_predecessor_claims` with the clause they close.
- **Added:** `v4-gap-01-physical-placement-unverified` and `v4-gap-02-endpoint-distribution-partial`.
- **Marked partially superseded:** `v3-gap-01-study05-remaining-regimes`, whose placement-distribution
  clause is now owned by the two successors. Its still-open regimes are named in
  `remaining_dimensions`, so the claim no longer overlaps its successors.
- **Retirement chain:** the chain from v1 (`claim-gap-04-study02-03-validation-open`) and v2
  (`v2-gap-04-study05-unmeasured-regimes`) is preserved and extended, not rewritten.

## The defect this closes

The v3 ingest retired `v2-gap-04` correctly, clause by clause. It did not notice that the Study 05
runs it was ingesting also closed clauses of two *other* active gaps:

| Retired claim | Clause as printed in v3 | Why it is now false | Successor |
|---|---|---|---|
| `v2-gap-05-colocation-unverified` | "Study 04 y1/y2 and Study 05's placement pair were never executed." | `v3-06-colocated-vs-noncolocated-locality` records the Study 05 key-locality pair as executed (2026-09-23, `20260923T1215Z-res2-placement-yb3`). | `v4-gap-01` |
| `v2-gap-07-real-network-balanced-endpoints` | "balanced client access across cluster endpoints is not proven." | `v3-05-yb-cluster3-strict-relaxed` runs where "the client spreads operations over all three endpoints". v3's own `coverage.json` already said so, while `claims.json` still said the opposite. | `v4-gap-02` |
| `v2-gap-05-colocation-unverified` | "no study verifies actual node placement." | **Still true.** Retained, narrowed, in `v4-gap-01`. | `v4-gap-01` |
| `v2-gap-07-real-network-balanced-endpoints` | "no real network" | **Still true.** Retained in `v4-gap-02`. | `v4-gap-02` |

The successors therefore split the retired claims rather than close them: what the new runs
measured moved to the numeric claims, and what remains unmeasured is stated precisely enough to
distinguish a *regime* gap from an *instrument* gap. Nothing new is measured here, and no number is
restated in a gap claim.

## New lifecycle metadata

v4 adds five fields; the schema in `v4/schema.json` is the contract, and
`tools/evidence validate` enforces it.

| Field | On | Meaning |
|---|---|---|
| `status` | every claim | `active` / `partially_superseded` / `superseded` / `retired`. Only `active` and `partially_superseded` may appear in `claims[]`; the other two states are expressed by moving the claim to `retired_predecessor_claims`. |
| `gap_kind` | gap claims | one or more of `coverage`, `schema_limitation`, `instrument`, `unstable_measurement`. This is what lets the book badge `v2-gap-02` (a schema limitation: no design can answer the hold funnel) differently from a regime that simply was not run. |
| `superseded_by` | partially superseded claims, retired entries | the successor claim ids |
| `closed_dimensions` | partially superseded claims, retired entries | the specific clauses now measured |
| `remaining_dimensions` | partially superseded claims | the clauses still open |

At the time of writing, no claim other than `v3-gap-01` carries `superseded_by`, and no claim other
than the two retired gaps carries `closed_dimensions`.

## Book prose corrected by reference

Edited to the v4 wording and claim ids (this task owns these files):

- `book/concepts/distributed-placement.typ` — the two cards become `v4-gap-01` / `v4-gap-02`, and
  "What it does not measure" now separates the instance count from the endpoint distribution.
- `book/chapters/topologies.typ` — the closing gap names the successor gaps and the single-endpoint
  DSN of the study 01/02 multi-node cells.
- `book/chapters/major-families.typ` — the placement sentence names `v4-gap-01`.
- `book/chapters/evidence-and-reproduction.typ` — the registry path is no longer a literal.

## Wiring moved to v4

- `book/lib/evidence.typ`, `book/main.typ` — the registry version is a single build input
  (`registry_version`), from which both the source-relative path and the display path are derived.
  No prose names a version.
- `book/build.sh` — evidence digest, claim totals and the manifest's `evidence_registry` read v4.
- `tools/evidence` — `validate-v4`, default `validate` dispatching on `schema_version == 4`, and v4
  lifecycle rules (retired claims may not also be active; a gap must declare `gap_kind`; a
  partially superseded claim must name existing successors).

## Related findings

- **R03** (the prose literal `book/evidence/v2/claims.json`) is closed here, since it is the same
  single-source-of-truth defect.
- **R02** (provenance "unknown") is *not* a defect of the committed artefact: pages 1–2 carry real
  commit, describe, dirty-state, clock, digest and hash values, and the only literal `unknown` is
  Typst's own `0.15.1 (unknown commit)` version string. The build is hardened against a *preview*
  build silently claiming release provenance (see `book/build.sh`), and that string is no longer
  printed as if it were a provenance field.
