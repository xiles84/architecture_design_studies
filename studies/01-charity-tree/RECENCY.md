# Study 01 v4 — "who donated last in a period"

**Protocol.** Written by HIGH before any code, in the form methodology 6a asks for: the
question, its exact semantics, the designs that answer it, the controlled pairs, how the
answers are verified, what is audited, what is measured, and what the measurement cannot
say. Nothing here is a result.

**Owner's request (2026-09-15):** *"could we also include the situation where we have a
query to know the people that made the LAST donation last week (or another period of
time). One of the options I thought is adding a flag in the donations saying
last_donation and index with it, but if you have other ideas we can add to the study."*

The flag is one of the designs below (D20). It is measured against the alternatives
rather than assumed, which is the whole method of this repository.

---

## 1. Why this question is not one of the twelve

The existing catalogue asks *"when was the last donation"* of one charity (q02) or one
donor (q06): **one key in, one row out**. The new question inverts it —

> give me every **donor** whose **most recent** gift falls inside a window

— and that is a different shape. The answer is not a row selected by a key; it is a
**set defined by a per-donor aggregate**, and no single index over `donation` orders
donors by their own maximum. Every design below is a different answer to *where the
"this is the donor's latest gift" fact lives*: nowhere (derive it), on the child row
(a flag), on the parent row (a rollup), or inside the parent document (embedding).

It is also the first question in this study where **the window matters**, and that
changes what "correct" means — see §2.

**Scope note.** Study 01's README excludes "anything requiring materialised views or
time-bucketed tables" at the owner's request. The owner's 2026-09-15 request lifts that
exclusion **for this question only**. Materialised views and time-bucket tables remain
out of scope and no design below uses one; the work stays inside the study's subject,
which is the shape of the tree and where derived facts are placed.

---

## 2. Exact semantics, and the trap in them

Let `last(p)` be the greatest `donated_at` among donor `p`'s donations, over the whole
history. For a window `[since, until)`:

> **Q:** the donors for which `since ≤ last(p) < until`.

Two parameter regimes are measured, because they are **not the same question** even
though they share one SQL text:

| Regime | Window | Business meaning | Note |
|---|---|---|---|
| **trailing** | `until` = end of the dataset's time span, `since` = `until − 7 days` | "who gave last week" — the thank-you / re-engagement list | here `last(p) ∈ window` ⟺ `p donated in the window` |
| **historical** | `until` = end − 90 days, `since` = `until − 7 days` | "who gave that week **and has not given since**" — the lapsed-donor list | here the two are **different sets** |

The trap is that in the trailing regime the cheap wrong answer —
`SELECT DISTINCT person_id FROM donation WHERE donated_at >= $1` — is accidentally
right, because nothing can be newer than the end of history. In the historical regime it
is wrong: it returns everyone *active* in that week, including donors who are still
giving today. **A design is only credited with answering this question if it answers
both regimes**, and the correctness gate checks both.

This is also why the window is bound as a parameter pair and never written as
`now() - interval '7 days'`: the generated history ends at a fixed instant
(`2026-09-01`, `harness/gen.go`), so a wall-clock window would drift with the calendar
and make two runs incomparable. The harness derives `until`/`since` from the dataset,
and records them in the result file.

### The four statements

Added to **every** design's `queries.sql`, each expressed the way that design's own
schema allows:

| # | Name | Params | Answer |
|---|---|---|---|
| q13 | `q13_donors_last_gift_window` | `since`, `until` | the 100 most recent such donors: `person_id`, `full_name`, `last_at`, newest first |
| q14 | `q14_donors_last_gift_window_count` | `since`, `until` | how many such donors — a count cannot stop early, so it prices the same access path without the `LIMIT` |
| q15 | `q15_charity_donors_last_gift_window` | `charity_id`, `since`, `until` | q13 restricted to one charity — the segment a charity actually mails |
| q16 | `q16_lapsed_donors_count` | `since` | donors whose `last(p) < since` — the complement, with a result set two orders of magnitude larger |

q13 and q14 are deliberately a pair: `LIMIT 100` lets an ordered access path stop after
100 rows, and the count denies it that. A design that wins q13 by early termination and
loses q14 has been described exactly, instead of being called "faster".

q15 is the segment query, and it is the one that re-touches this study's spine: in
D1/D2/D11/D12 the charity of a donation is only known through `person`, so the query
joins; from D3 onward `donation.charity_id` is there to be filtered directly. In this
dataset a person belongs to exactly one charity and their donations carry it, so
"their last gift to this charity" and "their last gift" coincide — the filter changes
the access path, not the answer.

---

## 3. The designs

Every new design is **exactly one decision** away from a design that already exists, so
a measured difference has one explanation. Nothing in D1–D17 changes: no schema, no
index, no write path, no existing query text. They gain the four statements above and
nothing else, so every published result of theirs stays valid and comparable.

| ID | Directory | From | The single change | Engines |
|---|---|---|---|---|
| **D18** | `d18_recency_probe` | D3 | q13–q16 rewritten as a **parent-driven probe**: walk `person`, and take each donor's newest gift with one descent of the existing `donation_person_time_idx`. Schema, indexes and write path byte-identical to D3 | both |
| **D19** | `d19_recency_window_sql` | D3 | q13–q16 rewritten **window-first** over the existing `donation_time_idx`: read the window, then confirm each candidate is that donor's maximum. Schema, indexes and writes byte-identical to D3 | both |
| **D20** | `d20_recency_flag` | D3 | `donation.is_last_donation BOOLEAN NOT NULL DEFAULT false`, a **partial index** on it, and trigger maintenance that takes the donor's `person` row lock first — the owner's proposal, made concurrency-safe | both |
| **D21** | `d21_recency_flag_unguarded` | D20 | the same maintenance **without** the person-row lock — **negative control**, expected to leave two flagged rows for one donor under concurrent inserts | both |
| **D22** | `d22_recency_rollup_idx` | D4 | one index: `person (last_donation_at DESC)`. D4 already *stores* `last_donation_at`; it has no global access path to it | both |
| **D23** | `d23_recency_rollup_app_idx` | D5 | the same one index on the **application-maintained** rollup | both |
| **D24** | `d24_recency_flag_colocated` | D20 | `donation`'s primary key rebuilt so one donor's donations share a tablet (D7's placement lever), flag and index unchanged | yugabyte |

### What is deliberately *not* a new design

**Embedding does not get one, and that is the finding.** D6 holds the entire history
inside the donor row and D9/D10 hold the newest 20; both answer q13–q16 by unnesting or
by falling back to the `donation` table, and neither offers an access path that orders
*donors* by their own latest gift. An index on `(recent_donations->0->>'donated_at')`
does not rescue it: the cast from the document's text timestamp to `timestamptz` is
`STABLE`, not `IMMUTABLE`, so PostgreSQL will not build the index, and indexing the raw
text makes correctness depend on the rendering of a JSON value. The conclusion to
*measure* is therefore: for a cross-parent recency question, embedding pays its write
cost and buys nothing, while materialising one scalar on the parent (D22/D23) buys the
whole query. D6, D9 and D10 are measured on q13–q16 to show the size of that gap, not
excused from it.

### The maintenance that the flag actually needs

This is the part that makes D20 a design rather than a column. Specified here so that
implementation is transcription, not invention.

- **Insert.** Lock the donor: `SELECT 1 FROM person WHERE person_id = NEW.person_id FOR
  UPDATE`. Then clear any currently-flagged donation of that donor that is not newer
  than the arriving one, and set `NEW.is_last_donation` only if no strictly newer
  donation exists. The lock is what makes "exactly one flagged row per donor" hold under
  concurrency; D21 omits it and, at READ COMMITTED, two simultaneous inserts for one
  donor each fail to see the other's uncommitted row and both set the flag.
- **Backdated insert.** Falls out of the same statement: the new row is not newer, so it
  is not flagged and the existing flag is not cleared. It is measured as its own write
  op (§6) because an append-only generator never produces one by accident, and it is
  precisely where a naive "always move the flag to the new row" implementation breaks.
  It needs **no new SQL in any design**: `w_insert_donation` already binds `donated_at`,
  so the op is the existing statement with a timestamp one day before the dataset's
  epoch start — always older than every generated donation, so "the flag must not move"
  is an exact expectation rather than a probabilistic one.
- **Delete.** If the deleted row carried the flag, promote the donor's next newest
  donation. This is the same asymmetry the study already found in D4: a `MAX` cannot be
  decremented, it has to be recomputed.
- **Amount correction.** `donated_at` does not move, so the flag does not move. The
  trigger must not fire work here, and the write benchmark shows whether it does.
- **Bulk load.** The loader sets the flag directly (it already knows each donor's newest
  donation) and the triggers are attached afterwards, exactly as D4's rollups are
  loaded. Loading through the trigger would measure the loader.

D22/D23 need no new maintenance: `person.last_donation_at` is already maintained by
D4's trigger and D5's application path respectively, and those paths are already
measured. The single index is the whole change — which is the point of the pair.

---

## 4. Controlled pairs

| Pair | Isolates |
|---|---|
| D3 → D18 | **SQL formulation alone**: one index probe per donor against a plain `GROUP BY`, on identical bytes of schema |
| D3 → D19 | **SQL formulation alone**, second formulation: window-first against a plain `GROUP BY` |
| D18 → D19 | which formulation suits which regime — the probe costs one index descent per donor whatever the window holds, the window-first plan costs what the window holds and then one probe per candidate |
| D3 → D20 | **materialising the fact on the child**: one boolean, one partial index, one trigger |
| **D20 → D21** | **the guard**: identical schema, identical index, identical queries; the person-row lock is the only difference |
| D4 → D22 | **one index**: the same stored rollup with and without a global access path |
| D5 → D23 | the same, on the application-maintained rollup |
| D22 → D23 | **trigger vs application** maintenance of the same column, i.e. pessimistic vs optimistic |
| D20 → D22 | **where the derived fact lives**: on the child (one flag per donation, index churn on every insert) or on the parent (one scalar per donor, a hot row per donor) |
| D3 → D6 / D9 / D10 | what embedding does with a cross-parent recency question |
| D20 → D24 | **data placement only**, identical SQL (YugabyteDB) |
| yb-single → yb-cluster3 | adding two more nodes |

D22 and D23 read identically by construction and extend the study's existing D4/D5
noise-floor control to the new statements.

---

## 5. Correctness

**Verification (the gate, before any timing).** The harness computes, in Go, each
donor's true `last(p)` from the generated dataset, and from it the expected answer to
q13–q16 in **both** window regimes. Checked per design, per cell:

- q14 and q16: exact counts.
- q13 and q15: the returned rows must be the expected donors, in non-increasing
  `last_at` order, with the boundary tie handled the way the existing gate handles ties
  (any member of a tied set at the `LIMIT` boundary is accepted).
- A design failing any of these is not timed, and the cell is recorded as failed.

**Recency audit (after every phase that writes).** Recomputes `last(p)` from the
`donation` table and reports, per design:

| Symptom | Designs it can affect |
|---|---|
| a donor with ≠ 1 flagged donation | D20, D21, D24 |
| a flagged donation that is not the donor's newest | D20, D21, D24 |
| `person.last_donation_at` disagreeing with the table | D4, D5, D22, D23 |
| a cache head that is not the donor's newest | D9, D10 |

The audit reports counts **and up to ten examples**, following the D9 cache audit: a
count alone has never been enough to diagnose one of these.

**Negative control, methodology 5a.** D21 must be **seen** to fail the audit in the
hot-donor experiment. If it does not, the experiment did not create enough contention
and the other designs' clean audits mean nothing that run — the report must say so
rather than claim correctness.

---

## 6. What is measured

Four groups, all inside the existing v3 enhancement runner and lock discipline, every
trial on its own fresh load.

1. **`recency-reads`** — q13–q16 in both window regimes, every design, every topology.
   Reports per-query throughput and latency, and a per-question winner table. The
   existing twelve-question geometric mean is **not** extended with these four: the old
   score must keep meaning what it meant, and the new questions get their own score.
2. **`recency-maintenance`** — isolated writes on fresh loads: `insert`,
   `insert_backdated` (new), `delete`, `update` (amount correction), with the recency
   audit after each. This is where the flag's price is paid: a second row update plus
   partial-index churn on the hot path, and a recompute on delete.
3. **`recency-hot-donor`** — 1, 4, 8, 16 writers inserting donations **for the same
   donor**, which is where per-donor maintenance serialises. D20 (pessimistic, person-row
   lock), D22 (trigger rollup), D23 (application CAS, optimistic) and D21 (no guard,
   control). Records throughput, retries/aborts, latency tails, and the audit. This is
   the study's optimistic-vs-pessimistic requirement for this question.
4. **`recency-placement`** — D24 against D20 on `yb-cluster3`, identical SQL, with
   physical placement evidence captured the way D7's is, plus per-node CPU throttling
   and connection distribution. Running both in one container on one laptop is **not**
   evidence of colocation; the placement listing is.

### Required comparisons (AGENTS.md), mapped

| Requirement | How this protocol covers it | Gap |
|---|---|---|
| 1 — calculated sizing | §7 | client saturation must be re-measured for q14/q16, which are much slower per op than the existing reads |
| 2 — rollup / rolldown / embedding | rollup D22, D23; derived fact on the child D20, D21, D24; embedding D6, D9, D10 measured on the same statements (§3, "not a new design") | none for this question |
| 3 — colocated vs non-colocated | D20 → D24 on `yb-cluster3`, group 4 | single physical host; no real network separation (unchanged from v3) |
| 4 — optimistic vs pessimistic | D20 (lock) vs D23 (CAS) vs D22 (trigger), group 3, under one invariant: exactly one true "latest gift" per donor | a third strategy (deferred/asynchronous maintenance) stays unimplemented, as in v3 |

---

## 7. Sizing

Unchanged from the study's established budget, restated here because methodology 4
requires the calculation to be in the protocol rather than inferred:

| Resource | Value | Reasoning |
|---|---|---|
| VM available | 8 CPUs, 16 496 418 816 B (≈15.4 GiB) | measured live, `docs/environments/host-zenbook-ux5406sa.md` |
| Per database node | 2 CPUs, 3 GiB (`infra/versions.env`) | held constant per node, so "1 node vs 3" answers "what do two more machines buy me" |
| yb-cluster3 total | 6 CPUs, 9 GiB | 3 × the per-node budget |
| Benchmark client | 2 CPUs, 2 GiB | one container, separate from the database containers |
| Host reserve | ≥ 2 CPUs on the 3-node topology (8 − 6 database − 2 client is already oversubscribed by 2) | recorded as a known oversubscription of this laptop, not a clean reserve; it is why every YugabyteDB number here is CPU-quota bound |
| Equal-total control | 1 node at 6 CPUs / 9 GiB, separately labelled | v3 already runs this; the recency reads are added to it |
| Readers | 8 connections (study default) | held fixed across every pair |
| Writers | 1 / 4 / 8 / 16 in group 3 only | the labelled experimental variable |
| Pool allowance | readers + writers + 2 control connections | the harness's existing rule |

**Client calibration is a stop condition, not an assumption.** q14 and q16 are whole-table
aggregates in most designs; at 8 connections they may saturate the database long before
the client, or (in D20/D22) return so fast that the client becomes the limit. The dev
checks measure delivered demand and client CPU throttling for the four new statements
before the matrix is sized, and the result is recorded.

---

## 8. What this will not be able to say

Stated before the numbers exist, so it cannot be tuned to them afterwards.

- **Nothing about a production SLO.** Closed-loop harness, shared laptop cores,
  CPU-quota-bound YugabyteDB containers, no real network. Tails are compared between
  designs, never quoted.
- **Nothing about a dataset whose donors keep donating at a different rate.** The
  window's selectivity is a property of the generated history; the trailing-7-day
  window over a 5.7-year span selects a small fraction of donors, and a system whose
  donors are mostly active would move every design's cost.
- **Nothing about index maintenance at scale.** The partial index in D20 is small
  *because* it holds one row per donor; at a scale where most donors have one donation
  it approaches a full index, and this dataset does not reach that regime.
- **Nothing about asynchronous maintenance.** A queue or logical-decoding path is the
  obvious fourth answer and is not implemented here.
- **Placement conclusions are bounded** by running three "nodes" on one laptop.
