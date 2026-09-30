# Supersession ledger — evidence registry

## v5 → v6 (2026-09-25) — method metadata only

**Attribution.** Written by `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI; effort setting
not exposed) on 2026-09-25, task `repo/empirical-method`, implementing the owner's 2026-09-25
empirical-architecture charter (prediction before measurement; three evidence forms; no attribution
from a package without an isolation pair; replication before domain expansion; crossovers and
boundaries preferred).

**What changed.** Nothing about the claims themselves. All 29 v5 claims are carried forward with their
statement, limits, confounds, anchors, support cells, trial counts, review state and lifecycle
**byte-identical**, pinned by `TestV6CarriesV5ForwardExceptTheReplicatedClaim` in `tools/evidence` —
with **one exception and two additions, all before the package tag**: `v2-13-trigger-rollup-cost` was
re-scoped by the corpus's first run-level replication, which falsified its previously cited 6.85x;
`v2-17-optimistic-pessimistic-crossover` was **added** with the corpus's first bounded crossover (seven
writer levels, crossing bracketed between 1 and 2 writers), and whose retry-based mechanism was later
withdrawn by the retry control; and `v2-18-maintenance-size-interaction` was **added** with the corpus's
first `interaction` evidence form (maintenance strategy × replacement size). v6 therefore carries 31
claims: 29 from v5, one of them re-scoped, plus two new. See `ANALYSIS.md`.
v6 adds:

| Added field | Values | Count across the 31 claims |
|---|---|---|
| `evidence_form` | isolate / package / interaction / not_applicable | 11 / 9 / 1 / 10 |
| `attribution` | withheld, isolated_by, note | 20 withhold attribution |
| `control_coverage` | present / absent / not_applicable | 7 / 17 / 7 |
| `replication` | runs, order, spread, tier | 29 single-run, 2 run-level (one a gap) |

The optional `prediction` and `crossover` fields are defined and machine-checked, and are no longer unused:
`v2-13` carries a recorded prediction with verdict `falsified`, and `v2-17` carries the corpus's first
bounded crossover (crossing interval `[1, 2] writers`).

**Plain-language rewording of the cards (2026-09-26, task `book/cards-plain-language`).** Twenty claims
had their `statement` and/or `limits` reworded so that each one says what was measured, then why it limits
the conclusion, and only then uses the project's term — the wording the book itself was given on
2026-09-26. No number, evidence id, anchor, confound or trial count moved. That is a distinct class from a
re-scoping, and `TestV6CarriesV5ForwardExceptTheReplicatedClaim` now keeps two declared sets to hold the
distinction: a **re-scoped** claim may change anything, while a **reworded** claim may change its text and
nothing else — if a reworded claim's confounds, trial counts or anchors move, the test fails. The claims
reworded are `v2-01`, `v2-02`, `v2-04`–`v2-07`, `v2-09`–`v2-11`, `v2-13`, `v2-14`, `v2-15`, `v2-gap-01`,
`v3-01`–`v3-06` and `v3-gap-01`. The registry files (`evidence_form`, `attribution`, `replication`,
`coverage.json`) are unchanged by this pass.

**How the values were assigned.** Every claim's assignment is in
[`method_migration.json`](method_migration.json), so a reader can audit the rule applied to each claim
rather than infer it. The per-claim verdicts and the full finding table are in
[`ANALYSIS.md`](ANALYSIS.md); the new coverage dimensions are in [`coverage.json`](coverage.json).

**Rules now enforced (each proven by a failing fixture).** A package claim may not name a component
without an isolation pair; a confounded claim must withhold attribution or name its pair; an unfired
control is not coverage; `run_replicated` needs two whole runs, a non-fixed order and a prediction;
`hardware_replicated` needs the hardware-independent tier; a reported reversal needs five levels and a
crossing interval.

**Retirements and supersessions.** None. v5 remains a frozen predecessor; a v6 supersession may resolve
against v1, v2, v3, v4 or v5.

---

## After the v6 tag — four re-scopings, and the `v2-08` and `v2-09` decisions (2026-09-26 … 2026-09-28)

**Attribution.** Written by `deepseek-flash` (DeepSeek, HIGH-capability, Deep Code CLI; effort not exposed)
on 2026-09-27, task `repo/v2-08-rescope`, and extended on 2026-09-28 by the same model in the
leader-and-executor session that ran `study-02/v2-09-buyer-sweep`. It records the re-scopings the T1
run-level replication programme made to the v6 package **after** its tag, and takes the `v2-08` decision two
worker tranches had explicitly left to HIGH and the `v2-09` decision that sweep forced.

**Why this section exists.** `TestV6CarriesV5ForwardExceptTheReplicatedClaim` keeps three declared sets —
added, **re-scoped** and reworded — so no claim can change silently. Four claims entered the re-scoped set
after the v6 tag, but only the tag-time exception (`v2-13`) was recorded above: the ledger said "one
exception" while the test listed five. This section closes that gap. The section above is left as written,
because it was true of the v6 tag.

**The four post-tag re-scopings.**

| Claim | Tranche tag(s) | What moved |
|---|---|---|
| `v2-06-hot-drop-designs` | `study-02/v2-06-replication` (2026-09-26) | strength → `run_replicated` (3 runs, AB/BA); statement re-scoped to carry the reproduced 4–8x range and the fact that absolute throughput fell 26–72% between sessions; limits extended |
| `v2-08-ledger-race-cost-128-buyers` | `study-02/v2-08-replication`, `study-02/v2-08-tier100-settle`, `study-02/v2-08-tier-axis` (2026-09-27) | strength → `run_replicated` (12 whole runs); statement re-scoped (below); anchors and limits extended |
| `v2-09-postgres-ledger-ratio-unsettled` | `study-02/v2-09-buyer-sweep`, `study-02/v2-09-lowbuyer-tiers` (2026-09-28/29) | strength `single_run_directional` → `run_replicated` (11 runs, alternating order); statement re-scoped (below) after its headline swap was falsified; a recorded prediction with verdict `falsified` added; the 32-buyer tier coverage then extended by a second tranche (two anchors, two limits, no further statement change) |
| `v2-11-section-seat-map` | `study-03/v2-11-memory-envelope`, `study-03/close-out` (2026-09-27) | strength → `run_replicated` (5 runs, both engines, AB/BA); the engine contrast is the replicated element; limits extended |

Each claim's own signed analyses are the primary record; this table is the registry-side index of them.

**The `v2-08` decision.** The primary run quoted the ledger at "1.1–2.3x … with tight spreads". Three
tranches of PostgreSQL race runs at 128 buyers — two whole runs over tiers 100/1,000/10,000, six tier-100
runs, and four runs covering all five tiers — established that:

- **The cost is monotone in tier.** Pooled run-median `p3/x1` is 1.075, 1.237, 1.409, 1.845 and 2.505 at 10,
  100, 1,000, 10,000 and 100,000 seats, monotone in the pooled medians and in every individual run. The old
  band is therefore an **interior summary**: its floor and ceiling are the 100- and 10,000-seat tiers.
- **"Tight spreads" is false at the small tiers.** The 100-seat cell spread 47–77%, the 10-seat cell
  16–71%, and the tier-100 ratio spans 0.992–1.712 across twelve runs (spread 54%), with the ledger slower in
  only 12 of 20 trial pairs at 10 seats.
- **The YugabyteDB, refund and organiser-edit figures were never re-measured**, and remain one run each.

Decision: the statement now **names the tier axis**, drops the flat band and the words "tight spreads", and
labels the not-re-measured arms as such. Strength stays `run_replicated`, which the PostgreSQL race
supports. The limits that referred to the old wording or the old band now speak of the *former* one; one
limit records this re-scope and one records that the YugabyteDB, refund and organiser-edit arms were not
re-measured. Nothing was retired, and the superseded wording is preserved here and in the two worker
analyses. **No other claim's statement or evidence changed in that pass.**

**The `v2-09` decision.** `v2-09` asserted that "the two designs change places depending on the workload",
resting on a primary comparison that moved buyer count (32 → 128) and design order (`p3`-first → `x1`-first)
together — its own limit said neither could be isolated. The buyer-count sweep crossed them at 1,000 seats:
eleven runs at 8/16/32/64/128 buyers, both orders, 5 trials, with a repeated reference cell. It found that
**the ledger is slower in 55 of 55 trial pairs**, that both orders agree on the sign at every level, and that
the magnitude is **non-monotone** (pooled medians 3.08, 1.32, 1.23, 1.59, 1.57) with **no crossing of ratio
1**; the 32-buyer cell that the primary read as 0.585 now reads 1.233 and 1.219. Two of the protocol's
predictions — a monotone rise with buyers and a bounded crossing — were **falsified**, and the reference cell
moved 8.6% between the two ends. Decision: the statement now says the direction is stable at 1,000 seats
across 8–128 buyers in both orders and that the primary's apparent swap came from its confound; strength
moves to `run_replicated` (11 runs, alternating order), with B2 recorded as a falsified `prediction` and the
ratio's non-portability as a limit. Nothing was retired, and the superseded wording is preserved here and in
the sweep's signed analysis.

**The `v2-09` tier extension (2026-09-29).** The sweep's first stated weakness was that it isolated the
confound at one tier. Two further whole runs at **32 buyers across all five tiers** (tag
`study-02/v2-09-lowbuyer-tiers`) show the same stable direction at every one: all ten tier × order cells
above 1 (1.021–1.606), the 32-buyer ratio below the 128-buyer one at every tier. The primary's tier-wide
reversal therefore does not reproduce **anywhere** on the axis. The tranche added two anchors and two limits
and changed no statement; it also fixed a caveat the sweep could not see — at 10 and 100 seats only 3 of 10
and 6 of 10 trial pairs favour CAS (cell spreads 86–174%), so the direction there sits at the measurement's
floor, while the 1,000-seat cell reproduced the sweep's ratio a day later within 1%.

**The 8-buyer collapse, explained (2026-09-29).** The sweep had kept an 8-buyer cell whose ledger arm fell
5x mid-cell, and neither discarded nor explained it. A recurrence test — six runs, twelve cells, sixty
trials, tag `study-02/v2-09-8buyer` — found **four low-throughput events** (619–655 ops/s against a session
normal of ~2,500) and they fall in **both** designs (three of the four in the CAS arm) and both orders. The
pre-registered arm-specificity and order predictions were falsified. Three of the four events were a single
trial that recovered; the one that spanned consecutive trials was the original cell, where the slow trials
happened to be the last three and so set the cell median. Conclusion: the low mode is a property of the
8-buyer level's very short trials, **not of the ledger**, and neither "collapse" nor "spike" is a mechanism —
only a description of where the event fell in a five-trial cell. `v2-09` gained six anchors and three limits
(one of which corrects the earlier limit that called it the ledger arm's), and no statement changed.

**The client-quota audit (2026-09-29).** The same session then read the `client_cpu` accounting the harness
has written since the study's first commit — one entry per **trial** for the race phase, per phase elsewhere
— which no analysis had used, and which had refuted the previous paragraph's own inference: the 8-buyer
slow trials show **zero** client throttling and *more* client CPU (2.6–2.8 s against ~1.9 s), so the client
was waiting, not starved. Scanning all 775 recorded client-CPU cells found **35 throttled**, concentrated in
the **read phase** of two anchors, `20260920T234953Z` and `20260913T021206Z`, and throttled by a
**design-dependent** amount within one run (2.3 s for `p2_precreated_skip_locked` to 58.1 s for
`r1_inventory_row`, out of ~1085 periods). Consequences recorded in the registry: `v2-09`'s limit about the
8-buyer cause was **corrected** (it had said the client was not recorded), and `v2-06`, `v2-07` and
`v2-gap-01` — the claims anchored to those two runs — each gained a limit stating that their reports' **read
figures are not quotable** because the client's quota set them. No statement, strength, confound or anchor
changed for any of the three: `v2-06` rests on race-phase sell rate, `v2-07` on answerability and storage,
`v2-gap-01` on the load phase, and none of those phases was throttled. The superseded inference is preserved
in `reports/outdated/20260929T-v2-09-8buyer-recurrence--deepseek-flash--2026-09-29.md`.

**The read-phase re-run (2026-09-29).** The audit implied the read-phase throttling was that session's. A
re-run of four designs (`p3`, `x1`, `c1`, `r1`), both orders, same phases and scale, shows it is **not**:
`r1_inventory_row` loses **420 of 634 and 422 of 633 CFS periods** to client throttling today (`p3` 2.3–5.9 s,
`x1` 2.6–4.6 s, `c1` 5.4 s), and its `q01_event_availability` reads sit **flat at ~38k ops/s across all five
tiers** (10 → 100,000 seats) where `p3` falls from 25,967 to 595 — a cell that does not move over a 10,000x
input range is the client's ceiling, not the query's cost. So the read family's numbers are client-bound for
the client-heavy designs at the study's two-core client budget, and the fix is **client sizing**, not another
re-run. The three limits above were updated to say so. What *does* reproduce is the designs' own ordering:
`q05` sell-through per band (C1 12.2–12.7k against P3/X1 ≈ 2.0k), the `q03`/`r03` gap (143–256x), and `r05`
answered only by the ledger (30.3k, 33.0k against the primary's 32.2k). X1's "~2x slower across the board"
read penalty **does not** reproduce (0.77–1.05 and 0.97–1.05 against the primary's 0.29–0.73), which confirms
the primary analysis's own suspicion that it was a cell-level effect. No claim's statement, strength,
confound or anchor changed; only the three read-hazard limits were extended.

**The client-budget calibration (2026-09-30), which corrects the paragraph above.** The audit and the re-run
both read the client's throttling as the read phase's ceiling. Varying the client's budget — possible only
after `ADS_CLIENT_CPUS` was added to `infra/lib.sh` in the same change — shows otherwise: at **three client
cores** the client is throttled **0 of 633 periods** and the reads are **slower** across every tier
(`r1`'s `q01` 52,884 → 43,726 at 10 seats), because the third core is taken from the database's share of the
machine's four CPUs (its throttled periods 238 → 847 → 1052). `r1`'s flat `q01` therefore survives with an
unthrottled client (ratio 1.05), so it is no longer evidence about the client; a plan check the same day
(`20260930T-read-q01-plan-check`, no new measurement) settled which: `r1`'s `q01` is a primary-key lookup of a
pre-aggregated counter in a 690-row table (cost 2.49, rows=1, 0.012 ms) against the baseline's `COUNT(*)` over
58,927 ticket rows (cost 4,517, 6.774 ms) — **O(1) by construction**, with the cell's level being the
machine's point-lookup rate. The three read-hazard limits were rewritten accordingly (their claim-specific tails unchanged: none
of `v2-06`, `v2-07`, `v2-gap-01` rests on a read level, and `v2-07`'s answerability and ~32k `r05` result
reproduce), and the superseded re-run analysis moved to `reports/outdated/`. The operational conclusion is
that the read family's levels cannot be fixed by re-allocating cores on a machine whose budgets already sum
to its CPU count: the options are a smaller client read footprint or a larger machine.

**The owner's decision and the footprint change (2026-09-30).** The owner fixed the environment — the project
assumes the computer it has and will not spend on more resources — and directed the read phase at a **smaller
client footprint**. Calibrated at three read-worker counts (tag `study-02/read-client-footprint`): client
throttling **422/633 → 33/632 → 0/631** at 8/4/2 workers, the database's own throttling **900 → 560 → 109**
periods, and at 2 workers every read cell is latency-bound (`ops_per_sec × p50 = 1.9` against 2 workers, both
designs, every tier). `r1`'s lead over `p3` on `q01` survives all three footprints and `p3`'s 100,000-seat
cell reads 621–631 ops/s in all three, so the *ordering* is what the read family should be read for.
Implemented: `-read-conns` in the harness and `READ_CONNS=2` as the runner's read-phase default (`--read-conns`
overrides; `0` follows `--conns`), leaving write, race and hold at the published 8; `client_read_conns` is
recorded in every manifest, so a reader can tell a footprint-2 read level from a footprint-8 one. **This is the
first change to a study phase's client configuration in this program**, and it is deliberately scoped to the
read phase only. No claim changed; the read-hazard limits that say read levels are not quotable stand until the
family is re-measured at the new footprint.

**Evidence.** Protocols `studies/02-ticket-booking/V2-08-REPLICATION.md`, `…-TIER100-SETTLE.md`,
`…-TIER-AXIS.md`, `…-BUYER-SWEEP.md`; signed analyses `20260927T-v2-08-replication--deepseek-flash--2026-09-27`,
`20260927T-v2-08-tier100-settle--deepseek-flash--2026-09-27`,
`20260927T-v2-08-tier-axis--deepseek-flash--2026-09-27`,
`20260928T-v2-09-buyer-sweep--deepseek-flash--2026-09-28`.

---

## Frozen predecessor text: v4 → v5

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
