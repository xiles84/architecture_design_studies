# Supersession ledger — evidence registry v4 → v5

**Attribution.** Written by `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI; effort setting
not exposed) on 2026-09-24, task `book/evidence-v5-consistency`, implementing the claim-versus-limit
findings of the third review round (GPT-5.6 Sol, high effort, 2026-09-24 08:22 −03:00) and the sweep
those findings implied.

**What this ledger is.** The single attributed record of what changed from
`book/evidence/v4/claims.json` (v4) to the active `book/evidence/v5/claims.json` (v5). It exists because
five claims asserted more in their statement or limits than their own evidence allowed — the class of
error this book is written to prevent — and because one of the five is resolved differently from the way
the reviewer proposed.

## Frozen predecessors are unchanged

- v1 `repo/book-evidence-registry-v1`, v2 `repo/book-evidence-registry-v2` (commit
  `075dd94672b7fdecdcf1f83e07bb3e7e9b87f764`), v3 `repo/book-evidence-registry-v3` (commit
  `d8ed10fcc0a3a13f02ec832b421e35e9829c4a23`) and v4 `repo/book-evidence-registry-v4` (commit
  `dd53ac29c4c817ad555036df152d64b5d3807447`) are untouched, and no tag was moved, deleted or reused.
- `sha256(book/evidence/v4/claims.json) = 814f40949d162f564ce2a831099727c2f1237a32629f99a2a24dd006db4c47c4`
  is the v4 predecessor this package was derived from.

## v5 is v4 with five claims re-scoped

24 of the 29 claims are carried forward **byte-identical**. Only these five changed, and in each case
the change removes an assertion the claim's own evidence did not support. No number was added, removed or
altered; no anchor was repointed; no support key changed.

### `v2-07-refund-answerability-and-storage` — limit corrected

**Before:** "the +33-35% storage delta is deterministic and stable across three runs"
**After:** "the +33-35% storage delta is a relation-size difference measured in the primary run and
deterministic for the same dataset and design; it is not a repeated timed result"

**Why:** the claim declares `trials.whole_run_replications = 1` and its single anchor is
`20260920T234953Z`. Nothing in that anchor, its report or its analysis evidences three runs for the
storage delta, so the limit asserted repetition the claim could not show. The verifiable content — a
deterministic relation-size difference — is kept. **Open, and recorded rather than papered over:** two
sibling Study 02 result directories exist (`20260913T021010Z`, `20260913T021206Z`); if a later pass
verifies that each measured the same delta, they can be added as corroborating anchors and the
repetition claim restored on evidence.

### `v2-10-guarded-confirm-refusals` — statement corrected

**Before:** "…retrying the refused statement once removed every refusal at no measurable cost; on
PostgreSQL the guarantee held everywhere."
**After:** "…retrying the refused statement once removed every refusal; on PostgreSQL the guarantee held
everywhere. This claim is about refusals, not throughput, so no cost conclusion follows from it."

**Why:** "at no measurable cost" is a throughput conclusion asserted by a claim whose own limit says it
is about refusals, not throughput. The reviewer's other option — registering a throughput claim — is
rejected here because this task commissions no measurement; if the owner wants that cost quantified it
belongs in a new Study 02/03 measurement task.

### `v2-14-cache-throughput-gain` — limit scoped

**Before:** "no YugabyteDB, cluster, verified placement or open-loop demand"
**After:** "this 2.2-3.5x range is supported only by the PostgreSQL single-node + Redis cells; do not
transfer the range to YugabyteDB, to a cluster, or to open-loop demand"

**Why:** as a limit on the 2.2–3.5× result the original was correct in intent, but after the v3 ingest it
reads as a statement about the whole corpus, which now does contain YugabyteDB and cluster cache cells.
The scoped version cannot be misread that way.

### `v2-16-strict-freshness-read-cost` — statement corrected

**Before:** "…inside the run's ~20% single-trial noise floor. The visible cost is on the write path (more
fences), so this must never be restated as 'strict freshness is free'."
**After:** "…inside the run's ~20% single-trial noise floor. This claim measures read throughput only: it
does not quantify the write-path cost of a strict protocol, and it must never be restated as 'strict
freshness is free'."

**Why:** the claim's scope is "read throughput only"; a write-path cost conclusion was not supported by
it. The guard against quoting strict freshness as free is retained, and the mechanism (extra write-path
fencing) remains in the book's prose, labelled as mechanism rather than as this corpus's measurement.

### `v3-06-colocated-vs-noncolocated-locality` — limits corrected, statement kept

**Statement:** unchanged. It says the colocated layout "keeps one donor's donations in one tablet".
**Limits, before:** "the runner could not produce a tablet/leader distribution … so the placement labels
are intent, not verified placement".
**Limits, after:** two separate limits — the logical partition mapping follows deterministically from
`PRIMARY KEY ((person_id) HASH, ...)`, so one donor's rows hash to one tablet *by configuration*; and
physical tablet/leader placement was not observed, so this is engine locality on one host and not
verified colocation.

**Why this differs from the reviewer's proposal.** They offered to soften the statement to "is designed
to keep", on the grounds that the card's own limit denies observation. The study's schema settles it:
`studies/05-cache-consistency/sql/reference/y1_colocated/schema.sql` declares
`PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)` while `y2_noncolocated` declares
`PRIMARY KEY (donation_id)`. In YugabyteDB a hash-keyed primary key places one person's rows in one
tablet by construction, so the statement describes a configuration fact and the *limits* were the part
that under-explained it. Softening the statement would have weakened a claim the configuration supports.

## The sweep, and what it found

All 29 active claims were reviewed for the class "the statement asserts more than its limits, trials or
anchors support", not only the five reported. The audit's per-claim verdicts are in
`ANALYSIS.md`; 24 passed unchanged. Two observations were recorded without a change:

- **`inner_iterations` on `v2-08` and `v2-09`** reads as the buyer count (128 and 32) rather than as
  repeat iterations inside a cell. If the Study 02 protocol defines it that way the field is right and
  this is only a naming hazard; if not, the field is being misused. Not changed here because verifying it
  needs the study's own definition, and inventing a correction would be exactly the failure this task
  exists to fix.
- **`v2-12`'s "one installed product"** is an opaque harness term, already owned by the book-side
  terminology task rather than by the registry.

## Validation

```text
tools/evidence validate-v5         0 errors
tools/evidence validate            0 errors, 29 claims, active source book/evidence/v5/claims.json
go test ./...                      green, including the new repetition-versus-trials lint fixtures
git diff repo/book-evidence-registry-v4 -- book/evidence/v4/   empty
v4 vs v5, 24 claims normalised     byte-identical
```
