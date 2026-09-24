# Execution brief — Book Edition 2 revision A

Revision A of the owner-directed revision of *Data Architecture Reference*. It implements the
technical-correctness items A–G and the tranche-A visuals (2.1, 2.3+2.7, 2.8) from the owner-supplied
external review, plus the redraw of Figures 6 and 7 and the cache-terminology audit. The owner approved
this plan on 2026-09-24; the plan is the authority, this brief is self-contained.

## Goal

Make the book's explanations of leases, fencing tokens, source-generation validation, relaxed freshness,
rollups, placement, concurrency and history **correct and unambiguous**, and add the four highest-value
explanatory visuals, without touching a single measured number or registry claim.

## Decisions already made (do not re-litigate)

- **Registry `v5` is read-only.** No claim id, statement, limit, trial count, strength, gap_kind,
  confound or number may change. If a correction needs a claim change, escalate; do not edit evidence.
- **Figure 6 shows three separated panels**: (1) lease-only/unsafe; (2) classic sink-side fencing token
  (the sink remembers the highest accepted token, so arrival order at the sink matters); (3) a labelled
  distinction from source-generation publication validation. Never conflate them.
- **Figure 7 is source-generation publication validation** with explicit `S0`/`v0` → `S1`/`v1`
  provenance and is scoped to the stale-fill race only. The strict-after-acknowledgement contract stays
  in the existing freshness-timeline figure (Figure 5).
- **2.9 stays a table.** Do not add a "can I quote this number?" flowchart.
- **2.2, 2.4, 2.5 and 2.6 belong to the later Tranche B task.** Do not author them here.
- **One recurring Donor / Donation domain** for new figures. Do not redraw study-owned figures
  (expiry, arbitration); those keep their study sources.
- **No figure embeds a rate.** An `observed result` figure names its active claim; every other new
  figure is a `conceptual illustration`. Form is `structure` or `sequence`.
- **Do not touch `book/dist/`.** Build to `book/build.sh --out dist/preview.pdf --manifest
  dist/preview-manifest.json` and delete both before committing.

## Work items

### 1. Prose corrections (exact sites)

1. `book/chapters/minor-variants.typ`, "Strict versus relaxed freshness": replace the "older … for a
   bounded window" sentence. Say relaxed is a family of weaker contracts and only bounded staleness
   declares a maximum Δ; the measured relaxed cells declared none. Match `concepts/cache-consistency.typ`
   and `glossary.typ`.
2. `book/chapters/major-families.typ`, Rollups → Mechanism: replace "application-maintained is cheaper
   but must be atomic with the publication". Say: where base fact and rollup share a database they are
   updated in one transaction; where they cannot, the copy crosses the authoritative boundary and
   publication, delivery, idempotency and reconciliation become relevant. Do not generalize "writes are
   cheaper"; scope it to the measured configuration-portal workload.
3. `book/chapters/how-to-choose.typ`, Step 3: "Application maintenance keeps writes cheap" becomes
   "kept writes near the reference in the measured configuration-portal workload" (keep the 6.85x trigger
   cost scoped as-is).
4. `book/concepts/information-placement.typ`, the figure text equivalent: replace "Each of the five
   placements is a copy … derive keeps no copy" with "Derive stores no additional copy. The remaining
   rungs introduce some stored derived structure or copy — an index, a materialized aggregate, a
   copied-down attribute or an external read copy."
5. `book/concepts/concurrency-control.typ`, Pessimistic section: replace "Latency per operation is
   higher … the retry loop disappears" with: contention is paid mainly as waiting/blocking rather than
   optimistic conflict retries; observed latency may rise or fall by workload; application-level retries
   can still be required for deadlocks, lock timeouts, serialization failures or other transient
   failures. Keep 1.76x scoped to one hot key, 16 writers, one run.
6. `book/concepts/derived-state-and-history.typ`: split **derived state** (rollup, rolldown, external
   read copy) from **retained history** (event/ledger, audit, temporal, snapshots). Replace "the only
   structures that can answer a question after the event was undone" with the fundamental claim: a
   current-state-only representation cannot reconstruct information it has destroyed. Remove
   "a fact stored twice" as the framing for all three.
7. `book/chapters/major-families.typ`, Append-only history → Costs: remove "writes append and fence".
   State that append-only storage does not by itself settle concurrency, ordering, duplicates or
   atomicity; arbitration is specified separately.

### 2. Cache terminology audit

Add or align these eight terms in `book/glossary.typ` (the single definition site), and use them
precisely in `book/concepts/cache-consistency.typ` and in figure text:

| Term | Meaning |
|---|---|
| invalidation | removes/marks a resident cache value stale; orders nothing about an in-flight fill |
| fill lease | bounds duplicate concurrent fills while ownership is valid; not a freshness proof |
| fencing token | monotonically orders operations **at a protected sink**; the sink rejects a token older than one it already accepted |
| source-generation validation | before publishing, check the authoritative current generation still equals the observed one; refuse if it moved |
| cache CAS | conditionally modifies **cache state** against an explicitly stated expected value/generation |
| hit validation | verifies a cached value against authoritative/version state before serving it |
| transactional outbox | makes a committed change durably observable for asynchronous propagation; it does not by itself make every hit current |
| TTL / early expiry | controls eligibility/load/recovery; not a freshness proof unless the declared contract is bounded by it |

Rename the cache chapter's "Publication fence / CAS" mechanism card and its matrix row so they name
**source-generation validation** explicitly; the sink-side fencing token and cache CAS are distinct and
must not be folded into that card or into one arrow/box label.

### 3. Figure redesigns

**`book/assets/sources/fig_lease_vs_fence.puml`** — 5 participants: `filler R1`, `filler R2`,
`lease / token coordinator`, `source (database)`, `cache / sink`. Three panels with numbered messages:
- *Lease only (unsafe):* R1 acquires lease #7 → slow source read → lease expires → R2 acquires lease #8
  → R2 reads newer S1 → R2 publishes S1 → R1 later publishes S0 → nothing rejects it. Add the note that
  lease ownership alone does not reject a stale publish unless publication is revalidated.
- *Sink-side fencing token:* R1 gets token 7 → slow read → R2 gets token 8 → R2 publishes (S1, token 8),
  sink accepts and records 8 → R1 attempts (S0, token 7), sink rejects 7 because 8 was already accepted.
  Note that this orders operations only once a token has reached the sink.
- *Distinction panel/note:* source-generation validation checks the authoritative current generation
  before publishing, so it can refuse token 7 even if no token-8 operation reached the sink. Name both.

**`book/assets/sources/fig_cache_stale_fill.puml`** — explicit versions: R1 misses, reads source
`S0 / v0`; writer W commits `S1` and advances source version to `v1`; W invalidates the key; R1 validates
the current source version before publishing, source returns `v1`, `v1 != v0`, R1 publishes nothing;
R2 reads → miss → authoritative read → `S1 / v1` (and may safely fill). Name the mechanism
**source-generation publication validation** — do not write "version / CAS check". Update the caption
and text equivalent in `concepts/cache-consistency.typ` to say this figure demonstrates the stale-fill
race, not the full strict-after-acknowledgement contract.

### 4. New tranche-A figures

All sources live in `book/assets/sources/`, get a `figures.json` entry, are embedded once, and carry a
short caption plus a text equivalent with no embedded figure number. Use the shared `_style.puml`; add
stereotypes/skins there for authoritative (solid), derived and external-copy boxes, and use solid arrows
for synchronous and dashed for derived/asynchronous propagation. Show cardinalities and explicit
versions/generations where order matters.

| id / source | form · status | content | placed in |
|---|---|---|---|
| `fig-families-er` / `fig_families_er.puml` | structure · conceptual illustration | three stacked ER panels using Donor/Donation: normalized (display_name only on Donor, 1..* Donation); rolldown (`donation.donor_display_name <<copied>>`, fan-out when Donor.display_name changes); rollup (`donor.lifetime_total <<derived aggregate>>`, maintained on Donation change). Cardinalities on every relation | `book/chapters/major-families.typ` after the intro |
| `fig-rollup-vs-rolldown` / `fig_rollup_vs_rolldown.puml` | structure · conceptual illustration | the direction comparison: rolldown parent→children fan-out vs rollup children→parent aggregate | `book/chapters/major-families.typ` Rollups section |
| `fig-history-vs-derived` / `fig_history_vs_derived.puml` | structure · conceptual illustration | panel A current-state-only vs retained `DonationEvent` history (PLEDGED/PAID/REFUNDED); panel B canonical current fact branches to a derived current-state copy and to retained historical facts | `book/concepts/derived-state-and-history.typ` |
| `fig-topology-physical-boundary` / `fig_topology_physical_boundary.puml` | structure · conceptual illustration | PostgreSQL single node, YB RF=1, YB RF=3 drawn inside ONE physical-host boundary labelled "logical multi-node topology; no real inter-host network measured"; logical key→expected tablet mapping separated from actual physical placement marked UNVERIFIED | `book/concepts/distributed-placement.typ` |

### 5. Registry, manifest, FIGURES.md

- Add each new figure to `book/assets/figures.json` (id/source/asset/status).
- Run `book/assets/export.sh --write` to render assets and refresh `figure-manifest.json`.
- Update `book/FIGURES.md`: the figure table (form × evidence status × placement) and a note that this
  revision's visuals were inspected in the rendered PDF by one model.

## Constraints

- Podman only; `book/build.sh` takes its own `ads-book-build-lock`, never the benchmark lock. This task
  measures nothing.
- Stage explicit paths; never `git add -A`.
- Keep `CONTEXT.md` and `LESSONS_LEARNED.md` current as part of the work (this task may touch them).
- Do not edit another analyst's report, analysis, review or digest; do not edit frozen evidence packages.
- If `book/evidence/check.sh` fails, fix the cause; do not widen a rule.

## Escalate only if

- an accepted correction cannot be made without changing a published claim, report or analysis;
- two accepted items conflict with each other or with a documented rule;
- a figure cannot be made legible at page size without dropping an accepted item.
