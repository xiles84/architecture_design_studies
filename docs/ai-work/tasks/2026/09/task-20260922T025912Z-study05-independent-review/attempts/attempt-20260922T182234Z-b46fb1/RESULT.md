# Result — Study 05 independent HIGH review

**Task:** `task-20260922T025912Z-study05-independent-review`
**Branch:** `repo/study05-independent-review`, base after the Study 04 review
**Required tag:** `repo/study05-independent-review`
**Capability:** HIGH / reviewer (model `deepseek-flash`)
**Benchmark:** none — artifact-only review

## What was delivered

1. A signed independent review analysis
   (`studies/05-cache-consistency/reports/analyses/20260921-cache-consistency-allgreen--deepseek-flash--2026-09-22.md`),
   indexed under the all-green report.
2. A machine-readable claim-verdict file
   (`book/evidence/reviews/20260922-study05-independent-review.json`).
3. Report regenerated so its index lists the review (generated-at line + one row only).
4. `CONTEXT.md` Study 05 status and open-gap 4 updated; one `LESSONS_LEARNED.md` entry.

## Verdicts

| Claim | Verdict |
|---|---|
| `claim-11-cache-throughput-gain` | **accepted** (warm cacheable reads, `small`, one trial, `db-only` framing) |
| `claim-12-strict-freshness-cost` | **accepted** (read throughput only; not "strict is free") |
| `claim-gap-02-study05-unmeasured-regimes` | **accepted**; four requirements added to `study05-v2-protocol` |

## Material defect found

The all-green analysis's three-instance wrong-read table **swaps the legacy and owned rows**. From
the run JSON:

| Cell | wrong reads | reads | rate | cache hits |
|---|---:|---:|---:|---:|
| `l:na:red:wth:rel` (legacy) | 35 287 | 40 449 | 87.24% | 40 442 |
| `o:opt:red:wth:rel` (owned) | 38 749 | 44 126 | 87.81% | 44 122 |

The generated report already attributes 87.81% to `o:opt:red:wth:rel`; only the derived prose
swapped the labels. The signed analysis is not edited; the correction is recorded for the evidence
correction task and the corrected mapping is stated in the review analysis.

## Requirements handed to `study05-v2-protocol`

1. Repeated trials (and ideally a second machine) for the three-instance result.
2. Real-TTL churn across the 300 s boundary.
3. Equal-total framing that charges the Redis container's budget.
4. Execution of the placement pair, or an explicit "colocation untested" statement.

## Scope requirements preserved

Single trial; PostgreSQL only; `small` (800 donors); closed loop; no equal-total accounting; no
real TTL crossing; no YugabyteDB, cluster or verified colocation; no open-loop demand; Redis
measured as a shared cache and never as an authoritative key-value store. This review is the first
independent reading of Study 05.

## Independently re-derived figures

- legacy no-cache 7 794 → memory-aside-strict 20 233 (2.60x), redis-aside-strict 18 667 (2.39x).
- owned no-cache 5 805 → redis-aside-strict 20 380 (3.51x), memory-aside-strict 16 613 (2.86x).
- owned-pess no-cache 7 043 → redis-aside-strict 15 782 (2.24x).
- Controlled strict-vs-relaxed read differences inside -8.3 % to +12.3 % (~20 % noise floor).
