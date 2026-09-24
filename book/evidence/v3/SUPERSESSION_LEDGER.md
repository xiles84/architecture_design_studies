# Supersession ledger — evidence registry v2 → v3

**Attribution.** Written by the executor `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI;
effort not exposed) on 2026-09-23, task `task-20260923T235907Z-book-evidence-study05-v3`,
implementing the evidence-update clause of
`brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`
(contribution `contribution-20260923T233727Z-synthesis-01-3bd104`).

**What this ledger is.** The single attributed record of what changed from
`book/evidence/v2/claims.json` (v2, frozen) to the active `book/evidence/v3/claims.json` (v3), and
which clause of the stale `v2-gap-04-study05-unmeasured-regimes` is superseded versus retained.

## Frozen predecessors are unchanged

- v1 milestone tag **`repo/book-evidence-registry-v1`** → commit `b7116cf12ef3eb832c53c44911fbefd7b87c3ec6`;
  `sha256(book/evidence/claims.json) = 0768644e0bb5c1d831b5e9a3972f604bd586186f20c45694b7a3ac066d747f01`.
- v2 milestone tag **`repo/book-evidence-registry-v2`** → commit `075dd94672b7fdecdcf1f83e07bb3e7e9b87f764`;
  `sha256(book/evidence/v2/claims.json) = 976e3635f5d2f769c3ab71e35e13083dc13bc728a82be32e10a7bfbe5f12e756`,
  `sha256(book/evidence/v2/schema.json) = 3cf39309351f0f8fab375c66827736f6f3150bb148511ff306b493014e715238`,
  `sha256(book/evidence/v2/confounds.json) = 7d4163917daa789ba46758eb577aae1abcbf212a6f84c8b591c113e305a48816`,
  `sha256(book/evidence/v2/coverage.json) = c8e6fda901f4bdf8ec2f348edfa9c9185777f205f4be833bce3fb25806de4931`.
- `git diff repo/book-evidence-registry-v2 -- book/evidence/v2/` is **empty**; neither v1 nor v2 was
  edited, and no v1/v2 tag was moved, deleted or reused.

## v3 is v2 carried forward plus a Study 05 update

- **Carried forward unchanged:** 22 of the 23 v2 claims (`v2-01` … `v2-16`, `v2-gap-01`,
  `v2-gap-02`, `v2-gap-03`, `v2-gap-05`, `v2-gap-06`, `v2-gap-07`). Their statements, anchors,
  support keys, limits, confounds, reviews and supersessions are byte-identical in content to v2.
- **Not carried forward:** `v2-gap-04-study05-unmeasured-regimes`, retired by reference (below).
- **Added:** `v3-01` … `v3-06` (numeric) and `v3-gap-01` (gap).
- **Retirement chain:** v3's `retired_predecessor_claims` names both
  `claim-gap-04-study02-03-validation-open` (retired in v1→v2) and
  `v2-gap-04-study05-unmeasured-regimes` (retired here), so the whole chain stays readable.

## `v2-gap-04` clause by clause

| v2-gap-04 clause | Verdict | v3 disposition |
|---|---|---|
| real-TTL churn unmeasured | **superseded** | `v3-01-churn-crosses-real-ttl`, `v3-02-churn-600-through-vs-aside`: 601 s and 1801 s cells crossed 2.0 and 6.0 real 300 s TTL boundaries with 0 wrong reads |
| repeated trials unmeasured | **superseded** (in-run) | the 2026-09-23 cells report 3 trials per measurement, and the two 601 s churn cells are a write-through/write-aside sibling pair (`v3-02`); a repeated *churn* cell is still absent and stated |
| cache resource accounting unmeasured | **superseded** | `v3-03-equal-total-framing-labelled`: a separately labelled `equal-total` arm now exists beside `db-only` |
| YugabyteDB/cluster unmeasured | **superseded** | `v3-04-yb-single-strict-relaxed`, `v3-05-yb-cluster3-strict-relaxed` on RF=1 and RF=3 with spread endpoints |
| verified placement unmeasured | **partly superseded, still a gap** | `v3-06-colocated-vs-noncolocated-locality` measures the locality pair, but no tablet/leader distribution was produced; the verification gap is retained in `v3-gap-01` |
| medium scale unmeasured | **retained** | `v3-gap-01`: the harness maximum is 3,000 people; the 8k/80k/800k tiers were not run |
| open-loop demand unmeasured | **retained** | `v3-gap-01`; owned by `task-20260922T200831Z-study05-v2-open-loop` |
| hard TTL expiry | **retained** | `v3-gap-01`: `hard expiries = 0` in all three churn cells; only probabilistic early expiry fired |

## New claims and the cell each resolves to

Every new numeric claim resolves at least one structured support key to a token that appears
verbatim in its primary report; the report's inputs digest appears in both that report and
`ANALYSIS.md`.

| Claim | Primary run | Inputs digest | Run tag | Report token(s) resolved |
|---|---|---|---|---|
| `v3-01-churn-crosses-real-ttl` | `20260923T1100Z-churn1800-through` | `8ad59fb36a65c17f` | `run/05-cache-consistency/20260923T1100Z-churn1800-through` | `1801.3`, `37063`, `2879077` |
| `v3-02-churn-600-through-vs-aside` | `20260923T1020Z-churn600-through` | `75af94d4c846e528` | `run/05-cache-consistency/20260923T1020Z-churn600-through` | `601.0`, `1467721`, `197196` |
| `v3-03-equal-total-framing-labelled` | `20260923T1150Z-res2-equaltotal-pg` | `e0dfc55f3c363944` | `run/05-cache-consistency/20260923T1150Z-res2-equaltotal-pg` | `equal-total`, `11k` |
| `v3-04-yb-single-strict-relaxed` | `20260923T1152Z-res2-yb1` | `6f04a8faf4a428bf` | `run/05-cache-consistency/20260923T1152Z-res2-yb1` | `596`, `514` |
| `v3-05-yb-cluster3-strict-relaxed` | `20260923T1157Z-res2-yb3` | `f2133069805cf82d` | `run/05-cache-consistency/20260923T1157Z-res2-yb3` | `942`, `846` |
| `v3-06-colocated-vs-noncolocated-locality` | `20260923T1215Z-res2-placement-yb3` | `a18817d592f767c5` | `run/05-cache-consistency/20260923T1215Z-res2-placement-yb3` | `882`, `673`, `852`, `520` |
| `v3-gap-01-study05-remaining-regimes` | `20260923T1100Z-churn1800-through` | `8ad59fb36a65c17f` | `run/05-cache-consistency/20260923T1100Z-churn1800-through` | none (gap; `gap_basis` cites the hard-expiries zero) |

`v3-02` also carries `20260923T1040Z-churn600-aside` (`cf16e598ad9b32f8`) as a corroborating anchor;
`v3-03` carries `20260923T1148Z-res2-dbonly-pg` (`b28d8dd331c4c968`) as a corroborating anchor.

## Book prose corrected by reference

- `book/chapters/scenarios.typ` — the "Where a scenario stops" gap paragraph no longer says real
  TTL churn, repeated trials and cache resource accounting are unmeasured; it names the measured
  churn/equal-total/engine evidence and keeps medium scale, placement distribution, hard expiry and
  open-loop demand as gaps. The retired `v2-gap-04` card is replaced by the v3 cards.
- `book/README.md` — the active registry is named as `book/evidence/v3/claims.json`.
- `book/lib/evidence.typ` — loads v3 and renders the union index.
- **Not edited, owned elsewhere:** `book/chapters/evidence-and-reproduction.typ` line 7 still names
  `book/evidence/v2/claims.json` in prose. The statement stays true in substance (v3 carries the v2
  claim ids forward), but the path is stale; that file belongs to
  `task-20260923T235922Z-book-figure-proofset` (`book/chapters/`) and is left for it. `book/dist/build-manifest.json`
  is the generated manifest of the last pre-v3 build and still records the v2 registry; it is
  regenerated only by a build, which this task does not run.

## Build wiring moved to v3

- `book/build.sh` — the evidence digest, the `claims_total`/`claims_indexed` check and the manifest's
  `evidence_registry` now read `book/evidence/v3/claims.json`. This was **required**, not cosmetic:
  retiring `v2-gap-04` means the v2 claim-id count no longer matches what the page renders, so a
  build still digesting v2 would fail its own verification (`claims=22/23`).
- `book/main.typ` — the title-page provenance line, the "must resolve to a claim" sentence, the
  registry table row and the evidence-index heading name v3.
- `book/dist/` is a generated artefact and is not regenerated here: this task performs no book build.
  `book/build.sh` and `book/dist` also sit in `task-20260923T235913Z-book-figure-provenance`'s owned
  paths and `book/main.typ` in `task-20260923T235922Z-book-figure-proofset`'s; the v3 references must
  be preserved, not reverted, when those tasks run.
