# Project context

The living state of this repository. Updated whenever a study starts, finishes, or
changes shape — so that anyone (or any future session) picking this up knows where things
stand without reading the git log.

**Last updated:** 2026-09-13

---

## Rule for every AI session: commit and tag, never push

**From 2026-09-13, by the owner's instruction:** AI sessions working in this repository
**commit and tag their own work, and never push** (the owner pushes and pulls). Reports of
the same study will reach different conclusions as studies are enhanced with new designs and
questions; tags are how a reader finds out which code, SQL and data produced each conclusion.
Untagged work means an unexplained change in conclusions.

- Commit at every meaningful step; never end a session with finished work uncommitted.
- `run/<study>/<run-id>` tags the commit behind each run (`run-study.sh --tag`).
- `study-NN/vX-<label>` tags study milestones — and, when enhancing a study, the state before
  and after the enhancement, so `git diff <before> <after>` explains a change in conclusions.
- Annotated tags only; never move, delete or reuse a tag; never rewrite history; never push.

Full rule: [CLAUDE.md → "Hard rule: commit and tag, never push"](CLAUDE.md).

### Tags so far

| Tag | Marks |
|---|---|
| `study-01/v1` | study 01 code and reports behind its first analysis (`795420b`) |
| `study-01/v2-second-analysis` | study 01 with the independent GPT-6 analysis and regenerated reports |
| `study-02/v1-harness` | study 02 designs, harness and shared platform as first run (`774d258`) |
| `run/02-ticket-booking/20260913T021206Z` | commit that produced study 02's `small` matrix (`7570648`) |
| `study-02/v1-analysis` | study 02 report and first signed analysis (`2bc7e1c`) |

---

## What this project is

An experimental record of **how architecture and data-modelling decisions actually affect
performance**, measured rather than asserted.

The distinguishing constraints:

- Claims are backed by **real runs**, not reasoning about what ought to be faster.
- Every test runs **inside Podman containers**, so a reader can reproduce it without
  installing anything but Podman.
- Every result names the **environment** it came from, documented down to the CPU
  topology and the specific ways that machine can mislead a benchmark.
- Every study produces a **report**; superseded reports move to `reports/outdated/`
  rather than being deleted.
- **Correctness gates timing.** A design that answers wrongly is never reported as fast.

Implementation language for tooling is **Go**.

## Repository layout

```
docs/
  methodology.md          the rules every study follows
  replication.md          how a reader reproduces a run
  environments/           one page per machine results were produced on
infra/
  versions.env            pinned images + per-node resource budget
  lib.sh                  shared podman helpers
  pg-single.sh            PostgreSQL, 1 node
  yb-single.sh            YugabyteDB, 1 node (RF=1)
  yb-cluster3.sh          YugabyteDB, 3 nodes (RF=3)
platform/                 shared Go module (adsplatform): core / ports / adapters
studies/
  01-charity-tree/        study 01 (see below) — own harness, predates platform/
  02-ticket-booking/      study 02 (see below) — built on platform/
```

Everything belonging to one study (SQL, harness, runner, image name, results, reports,
analyses) lives in that study's directory; only study-independent code is shared, in
`platform/`, laid out as ports and adapters (`platform/README.md`).

**Versioning.** Tags: `study-01/v1` marks the code behind study 01's published results;
`run/<study>/<run-id>` marks the commit a run was produced from (`run-study.sh --tag`, study
02 onwards). Results, manifests and reports record commit, `git describe` and a dirty flag.

Each study directory holds: `README.md` (the question and the designs), `sql/` (one
directory per design), `diagrams/` (PlantUML sources + rendered SVG), `harness/` (Go),
`run-study.sh`, `results/<run-id>/`, `reports/`.

## Current state

### Study 01 — tree structures (charity → person → donation)

**Status:** all experiments (survey, C, D, E) complete; D9 cache bug found, diagnosed and fixed as D10; two signed analyses written, including an independent GPT-6 review of the available evidence.

Eight designs, each a single deliberate change from its neighbour:

| ID | Design | Isolates |
|---|---|---|
| D1 | normalized-minimal | the floor: 3NF, no secondary indexes |
| D2 | normalized-indexed | indexes alone (SQL identical to D1) |
| D3 | flattened-fk | denormalising the grandparent key onto the grandchild |
| D8 | flattened-nofk | the cost of enforcing referential integrity |
| D4 | rollup-trigger | consolidating aggregates upward, via triggers |
| D5 | rollup-app | same aggregates, maintained by the application |
| D6 | embedded-jsonb | folding the **whole** child table into the parent row |
| D9 | embedded-hybrid | folding a **bounded** slice (newest 20) into the parent, keeping the table |
| D10 | embedded-hybrid-locked | D9 with a concurrency-correct cache trigger — prices correctness itself |
| D7 | yb-child-colocated | physical data placement (YugabyteDB only) |

D6 and D9 are the two ends of the embedding question. D6 takes it to its conclusion and
runs into unbounded write amplification — appending one donation rewrites a document that
grows with the donor's entire history. D9 caps the array at 20, which makes the rewrite a
constant, and keeps the real donation table underneath as the source of truth.

Three topologies: PostgreSQL 1-node, YugabyteDB 1-node, YugabyteDB 3-node RF=3.
Twelve queries, three write operations, correctness-gated, plans captured per cell.

**Verified:** all designs return identical answers on both engines (12/12 checks).

**Correctness beyond the answer check.** Two audits run after the write benchmark:

- **rollup audit** (D4, D5) — recomputes every aggregate from the donation table and
  counts parent rows that disagree, plus the signed drift in the headline total. This is
  what makes "optimistic vs pessimistic" a real comparison rather than a throughput race:
  a strategy that is faster and silently loses increments has not won.
- **embedded-cache audit** (D9) — checks the bounded slice still holds exactly the newest
  donations, in the right order. Compares donation-id sequences rather than raw JSONB,
  because the loader and PostgreSQL render timestamps differently in JSON and a byte
  comparison would call a working cache broken.

**Tail latency.** p50/p90/p95/p99/max are recorded for every cell. p99.9 and p99.99 are
recorded only where the sample count supports them (10 000 and 100 000 respectively, so at
least ten samples sit beyond the quantile); below that the field is omitted rather than
estimated, because a p99.9 drawn from 363 samples is the maximum relabelled. Deep tails
are also biased low by coordinated omission — the harness is closed-loop — so they are
compared between designs and never quoted as SLO figures.

Deep tails were added after the survey matrix started, and raw samples are not persisted,
so the survey's own cells carry p99 and max but not p99.9. The follow-up stages
(D9 pass, experiment C, trials run) rebuild the image and do carry them — which is where
they matter most, since contention shows in the tail long before the median.

**Which tables get written.** Reads span all three tables, and *which* table answers a
question is itself the design variable. Writes originally all originated at `donation`,
with `person` and `charity` written only as side effects (trigger, app rollup, embedded
array). Two operations that originate at the **parent** have since been added, because
both are strongly design-discriminating and their absence hid costs that fall specifically
on the embedded and rollup designs:

| Op | Why it discriminates |
| --- | --- |
| `update_person` — a donor edits their email | Touches no indexed column and no child data; about as cheap an update as exists. But D6 has made the person row *contain the donation history*, so the cost depends on the row that trivial edit lands on. Genuinely uncertain: PostgreSQL does not rewrite an unchanged TOASTed column, so a heavy donor's out-of-line array survives, while a median donor's ~2 kB array sits inline and is re-versioned in full. |
| `delete_person` — erasure, donor and all their donations | Identical logical outcome everywhere, very different work. Normalised designs delete N children then the parent, and in D4 every one of those child deletes fires a trigger that recomputes the parent's MIN/MAX from the survivors. D6 deletes one row and the whole history goes with it. |

`delete_person` is excluded from the *mixed* workload on purpose: erasure removes donors
that concurrent inserts are still picking at random, and the resulting foreign-key
violations would be an artefact of the harness rather than a property of any design. It is
measured in the isolated write benchmark, where it has the id space to itself.

**Reads and writes are always reported together.** Every design here buys read speed
with something, so the generated report opens with a *trade-off at a glance* table putting
the purchase and the price in one row: a geometric-mean read score, insert throughput,
the number of indexes maintained on every `donation` insert, and storage. In the
controlled-pair tables, read differences below the error bar are suppressed but **write
rows are always shown** — a hidden price tag reads as no price at all.

**The report checks variation using D4/D5 reads.** These designs have identical read SQL
and serve as a control for cell-to-cell variation. Regenerating the survey from digest
`558e89403fd016cc` gives 1.69x maximum disagreement across 36 comparisons (1.04x median),
correcting the earlier 1.56x note. The report uses the maximum to suppress small pairwise
claims. This is a descriptive check, not a confidence interval or a universal significance
threshold; important comparisons still need repeated trials.

**Analysis provenance.** Generated reports carry an **inputs digest** — a content hash
over every result file — and index the signed analyses written about it. Several
analysts, including different AI models, are expected to interpret the same data; the
digest is what lets a later one tell whether the data has already been analysed. See
methodology rule 11 and `reports/analyses/TEMPLATE.md`.

### Runs so far (all `small` scale, environment `host-zenbook-ux5406sa`)

Run ids starting `20260913-` were produced on **2026-09-12 UTC** — the date in the name was
hand-set by mistake (LESSONS_LEARNED → "Never hand-set a date"). Manifests hold true times.

| Run id | What | State |
|---|---|---|
| `20260912-small` | Survey: all designs × PG 1-node, YB 1-node, YB 3-node (D9 added in a second pass) | ✅ 26 cells. Its `update`/`delete` numbers share one table across write phases — use `20260913-writes-isolated` for writes |
| `20260912T180830Z-concurrency` | Experiment C: D3/D4/D5 × 1–32 writers × strategies × isolation | ✅ 31 cells |
| `20260912-mixed` | Experiment D: readers:writers splits, dashboard mix | ✅ 5 designs — surfaced the D9 cache bug |
| `20260912-trials` | Repeated-trials writes | ⚠️ `delete` **invalid** (finite-pool bug); report in `reports/outdated/` |
| `20260912-parentwrites` | `update_person` / `delete_person` | ⚠️ `delete_person` **invalid** |
| `20260913-writes-isolated` | Every write op on a fresh load, 3 trials, PG + YB 3-node | ✅ replaces both invalid passes |
| `20260913-history-*`, `-top-1000charities`, `-contention-*`, `-portal-*`, `-size-*` | Experiment E: limitation regimes | ✅ D7 3-node unbounded needed a second pass (YB killed a session mid-reload) |
| `20260913-d9-cache-race` | Reproduction of D9 cache corruption, with diagnosed examples | ✅ 3–5 wrong caches per 30 s window |
| `20260913-cache-fix`, `-cache-fix-writes` | D9 vs D10 (fixed trigger): correctness and cost | ✅ D10 always consistent, 0.96–0.98x D9 writes |

**Analyses:** both remain current and are indexed by the 16 input runs via their digests.

- [Claude Opus 5 — first analysis](studies/01-charity-tree/reports/analyses/20260912-study01--claude-opus-5--2026-09-12.md): overall conclusions and the regimes where a usually slower design becomes useful.
- [GPT-6 — independent final analysis](studies/01-charity-tree/reports/analyses/20260912-study01--gpt-6--2026-09-12.md): reviewed the existing results without running new benchmarks; explains D2/D3 through indexes, SQL and actual plan work, and records narrower disagreements about embedding, foreign keys and application rollups.

**D2/D3 interpretation:** the 1.548x survey ratio is a geometric-mean score, not measured
mixed-workload throughput. D3 copies `charity_id`, adds two charity indexes and rewrites
the charity queries. Those six queries score 2.71x higher together; the six with unchanged
SQL score 0.884x. The saved charity-feed plans show fewer examined donations and donor
lookups. The sum plan still has 31,293 heap fetches despite its `Index Only Scan` label.
The mechanism is selective access and changed query work; the exact overall gain needs
repeated trials. Do not describe it as a universal benefit from a copied column.

**Writing a final analysis:** follow [methodology rule 11](docs/methodology.md) and the
[analysis template](studies/01-charity-tree/reports/analyses/TEMPLATE.md). Regenerate the
input reports, record every run and digest, read existing analyses, and write a separate
signed file with recommendations, linked measurements, mechanisms, limitations, proposed
follow-ups and explicit disagreements. Define aggregate scores, distinguish observations
from explanations, and identify untested behavior. Never edit another analyst's file.
Regenerate the reports after adding the analysis so their indexes include it, and update
this context and the process lessons. Study 01's final analysis does not cover Study 02.

### Experiment E — limitation regimes

Owner's instruction: where a design's limitation would be acceptable in some domain, test
that regime too, so the conclusion can say when a losing design is actually the right one.
The owner deferred to the AI which studies stay unbounded and whether the argument extends
to other levels of the tree. Decisions:

- **Unbounded history and 10 large charities remain the primary profile** — a charity's
  donation history genuinely grows. Regimes are reported beside it, never instead of it.
- **Bounded history (cap 20 per donor, donation volume held constant)** is run only for the
  designs whose limitation depends on per-donor history: D3 (baseline), D4 (delete trigger
  recompute), D6 (whole-array rewrite), D9 (cache = entire history at cap ≤ 20), and D7 on
  YugabyteDB (per-donor tablet). Not D1/D2/D5/D8 — their limitations don't depend on it.
- **The argument does extend to the top of the tree.** 1 000 small charities instead of 10
  large ones (people and donations held constant) targets D4/D5's single hot rollup row per
  charity and D6's charity-scoped unnesting. Asked both as a full cell and as a contention
  sweep at 8 and 32 writers.
- **Donor-portal workload** (reads that never cross donors) for D3/D6/D9 — the regime that
  removes D6's worst weakness regardless of history length.
- **Small dataset** (`tiny` vs `small`) for D1/D2 — does indexing matter at all when tiny?
- Regimes already covered elsewhere are recorded in the map in `run-regimes.sh` (D4 at low
  concurrency → experiment C; D5 insert-only ledger → experiment C; D8 → not a performance
  regime).

Harness support: `-max-per-person`, `-charities`, `-read-mix portal`, `-cmd report-compare`
(same designs under two regimes, flags where the per-question **winner** changes).

### Not yet done (from the analysis's "what I would measure next")

- Repeated, alternating D2/D3 trials and variants that separate the copied key, new indexes and SQL rewrites; capture representative plans during equivalent database preparation.
- A `medium` (1 M donations) run: every result so far fits in memory.
- YugabyteDB 1-node FK re-measurement with trials (one unreplicated 1.8x figure).
- Longer repeated YugabyteDB 3-node donor-erasure measurements with/without FKs, and D10 cache correctness/cost on YugabyteDB.
- Asynchronous rollups (queue / logical decoding) — the design Note R1 points at but nobody built.
- Total-budget sharding framing (1 node x 6 CPU vs 3 nodes x 2 CPU).
- An open-loop load generator before any tail figure is treated as an SLO.

### Study 02 — avoiding overbooking (band → event → ticket)

**Status:** designs, harness, runner, diagrams and README done (commit `774d258`, tag
`study-02/v1-harness`; runner fix `7570648`). Dev checks on `tiny` in `results/devchecks/`.

| Run id | What | State |
|---|---|---|
| `20260913T021010Z` | first `small` matrix launch | ❌ **INVALID** — aborted at the first cell; provenance not captured (see its INVALID.md) |
| `20260913T021206Z` | `small` matrix: 14 designs × PG 1-node, YB 1-node, YB 3-node; tag `run/02-ticket-booking/20260913T021206Z`; digest `71bcee71d725d43d` | ✅ 40 of 42 cells; C2 failed on both YB topologies (YSQL lease expiry, diagnosed) |

**Analysis:** `studies/02-ticket-booking/reports/analyses/20260913T021206Z--claude-opus-5--2026-09-13.md`
— first signed analysis (by the model that built the study; conflict of interest stated).
Headline: no correct design oversold and both controls did on every engine; for a hot drop,
pre-created rows taken by CAS or SKIP LOCKED sell 4–8x faster than any single counter row; a
counter on the event row also blocks unrelated edits; C5 under-sells after refunds; C2 is
unusable for a hot drop. Main doubts: DB CPU quota throttling everywhere, no retry backoff
(likely understates R2), one trial, 60 s timeouts truncate the 100k tier, single YB query node.

The owner's question: events of 10 / 100 / 1 000 / 10 000 / 100 000 seats, no overbooking;
compare pre-created tickets against tickets created on booking, other strategies, and extra
tables such as reservations. Fourteen designs in four families:

| Family | Designs |
|---|---|
| pre-created tickets | P1 lock-first (FOR UPDATE) · P2 skip-locked · P3 CAS · P4 skip-locked + counter |
| created on booking | **C1 count at RC — negative control** · C2 same SQL at SERIALIZABLE · C3 lock event + count · C4 guarded counter · C5 unique seat + CHECK |
| extra table | R1 inventory row · R2 inventory buckets (sharded counter) · R3 pre-created seat-slot pool |
| reservation | **H0 hold, naive confirm — negative control** · H1 hold, confirm checked by DB clock |

Experiments per cell: correctness gate (5 reads + audit on load) → plans for reads **and**
writes → reads (availability per size tier) → isolated writes on fresh loads (spread
booking, cancel, publish per tier) → **sell-out race** (32 buyers per unsold event, per tier,
organiser editing the event) → **churn race** (10% cancel; measures under-booking) →
**holds** (H only: basket, abandonment, expiry, sweeper, late payments). An overbooking
audit (over capacity, duplicate seats, derived-state drift, ledger reconciliation) runs
after every writing phase.

**Dev-check observations (tiny, not results, not for conclusions):**

- The negative control C1 overbooked every race event on PostgreSQL; the audit caught all.
- C5 under-booked in the churn race on every tier, as its design predicts (max+1 never refills refund holes).
- C2 on the 10 000-seat event: 15 attempts per sale on PostgreSQL; on YugabyteDB the race
  hung until the deadline was made to bound retries (LESSONS_LEARNED).
- YugabyteDB 1-node sells one to two orders of magnitude fewer seats/s than PostgreSQL on this
  laptop; hot-row designs (C4, P4, R1, H) around 15–20/s, R2 buckets ~87/s. Reads show a shared
  ~65 ms p99. **Measured:** the client is never CPU-throttled; the YugabyteDB container is
  throttled in 40–65% of CFS periods in every cell (up to 720 s per cell). YugabyteDB numbers
  here are quota-bound — this also applies to study 01's YugabyteDB cells.
- C2 on YugabyteDB 1-node fails reproducibly (dev check and `small` matrix) with "the database
  system is shutting down". **Cause established:** SERIALIZABLE lock storm on a CPU-throttled
  node → tserver RPCs to master stall ~29 s → YSQL lease expires → tserver kills all SQL
  sessions. See `results/20260913T021206Z/yb-single/logs/c2_count_serializable.DIAGNOSIS.md`.
- H0's control fired on YugabyteDB too. Holds at a 250 ms TTL were degenerate there (96% of
  holds expired before payment); YugabyteDB topologies now run holds at `-hold-time-scale 10`.
- YugabyteDB 2025.2.6 runs READ COMMITTED as real RC even without
  `yb_enable_read_committed_isolation` (checked); the flag is set to pin it.

**Harness:** `studies/02-ticket-booking/harness/` on `platform/`. Flags worth knowing:
`-phases`, `-race-trials N` (fresh load per trial; report shows median and spread),
`-race-tier-budget`, `-race-timeout`, `-hold-*`, `-tiers`.

**Open for study 02** (from the analysis's "what I would measure next"): retry-backoff variants
(expect R2 to overtake R1); repeated race trials on PostgreSQL and YugabyteDB 3-node; a buyers
sweep (4–128); a single-statement SKIP LOCKED variant to separate round trips from locking in
P3's lead; YugabyteDB with connections spread over all nodes; a larger per-node CPU budget; an
asynchronous allocator (virtual waiting room); a second analysis by a different model.

## Decisions taken, and why

| Decision | Reason |
|---|---|
| Plain `podman` CLI, not compose | No compose provider ships with Podman; requiring one would break the "install nothing" rule |
| Per-node resource budget held constant | Makes "1 node vs 3 nodes" mean "what do two more machines buy me" |
| `--cpus` quotas, never `--cpuset-cpus` | This host mixes P-cores and E-cores that the guest kernel cannot distinguish; pinning would silently bias results |
| SQL in files, embedded via `go:embed` | Readable as documentation, and the binary provably runs the SQL shipped beside it |
| Bulk rollup at load, triggers attached after | Loading with per-row triggers would measure the loader, not the design; trigger cost is measured by the *write* benchmark |
| Queries benchmarked in isolation | A blended number hides which query moved, and attribution is the point |
| Bounded long tail in generated data | An unbounded power law would make the embedded design a strawman |
| Donations generated in time order | Mirrors an append-only stream, so surrogate key and timestamp correlate as in production |

## Open questions worth a future study

- Total-budget framing of sharding (one big node vs three small ones), as opposed to the
  per-node framing used here.
- Partitioning in PostgreSQL (`PARTITION BY HASH` + local indexes) as the non-distributed
  equivalent of D7's placement lever.
- Where the embedded design's crossover lies as array length grows — the write cost is
  O(history), so there is a document size at which it stops paying.
- Deferred / asynchronous rollups (queue, logical decoding) as a fourth point between
  D3's read cost and D4's write cost.
