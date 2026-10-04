# Book style — Edition 3: concept first, studies as proof

Owner decision, 2026-10-01 (goal amendment `20261001T155900Z-concept-first-edition-3`).
This file is the single source of the book's didactic rules. Change them here.

## The promise to the reader

- Each chapter teaches **one concept**.
- The **case studies show when it fits** — and when it does not.
- Text is **short**: bullets and tables first, prose only to connect them.

## Book structure

| Part | Chapters |
|---|---|
| Front | Title · Edition and provenance · **How to use this book** (one page) |
| I — Placing data | 1 Where a fact lives (the ladder) · 2 Indexes · 3 Rolldown · 4 Rollup · 5 Embedding |
| II — Many writers | 6 Hot rows · 7 Optimistic vs pessimistic · 8 Expiry and the clock |
| III — Time and copies | 9 History vs current state · 10 Caches and freshness · 11 Where a cluster puts data |
| IV — Putting it together | 12 The case studies (one profile per study) · 13 How to choose · 14 Reading a benchmark |
| Appendix | Glossary · Evidence and reproduction · Evidence registry (full cards) |

Order follows learning dependency: each chapter uses only ideas from earlier ones.

## Where each claim goes

Every renderable claim in the active registry is cited at least once (31 today).

| Chapter | Claims | Studies |
|---|---|---|
| 1 Where a fact lives | `v2-gap-06` | — (overview) |
| 2 Indexes | `v2-01`, `v2-04` | 01 |
| 3 Rolldown | `v2-05` | 01 |
| 4 Rollup | `v2-02`, `v2-13`, `v2-18`, `v2-gap-03` | 01, 04 |
| 5 Embedding | `v2-03`, `v2-11` | 01, 03 |
| 6 Hot rows | `v2-06` | 02 |
| 7 Optimistic vs pessimistic | `v2-12`, `v2-17` | 04 |
| 8 Expiry and the clock | `v2-10`, `v2-gap-02` | 03 |
| 9 History vs current state | `v2-07`, `v2-08`, `v2-09`, `v2-gap-01`, `v2-gap-02` | 02, 03 |
| 10 Caches and freshness | `v2-14`, `v2-15`, `v2-16`, `v3-01`, `v3-02`, `v3-03`, `v3-gap-01` | 05 |
| 11 Where a cluster puts data | `v3-04`, `v3-05`, `v3-06`, `v4-gap-01`, `v4-gap-02` | 05 |
| 12 Case studies | any, as badges | 01–05 |

Short ids above abbreviate the full claim ids in the registry.

## Old file → new home

| Edition 2 source | Edition 3 home |
|---|---|
| `chapters/how-to-choose.typ` | 13 How to choose (shortened; numbers move to boxes) |
| `chapters/major-families.typ` | chapters 2–5, 9, 10 (one family per concept) |
| `chapters/minor-variants.typ` | the concept that owns each variant (e.g. trigger vs app → 4) |
| `chapters/cross-family-comparisons.typ` | "Compared with" rows in fit tables; 13 How to choose |
| `chapters/topologies.typ` | 11 Where a cluster puts data |
| `chapters/scenarios.typ` | 12 Case studies; analogies stay labelled |
| `chapters/evidence-and-reproduction.typ` | Appendix |
| `concepts/information-placement.typ` | 1 Where a fact lives |
| `concepts/indexes-and-hot-rows.typ` | split: 2 Indexes, 6 Hot rows |
| `concepts/derived-state-and-history.typ` | split: 4 Rollup, 9 History vs current state |
| other `concepts/*.typ` | same concept, rewritten to the template |

## Concept chapter template

Fixed order. Target **2–4 pages**.

| # | Section | How | Limit |
|---|---|---|---|
| 1 | Opener | `concept-opener(number, title, idea)` | key idea ≤ 25 words |
| 2 | The problem | prose | ≤ 3 sentences, one concrete example |
| 3 | Picture | `figure-evidence` (methodology 15) | one figure; text equivalent ≤ 3 sentences |
| 4 | How it works | bullets | 3–5 bullets, ≤ 20 words each |
| 5 | When it fits | `fit-table(use: …, avoid: …)` | 3–4 rows per side |
| 6 | What it costs | `cost-table(reads, writes, storage, correctness, operations)` | one phrase per row |
| 7 | In the studies | `case-box(…)` | 1–3 boxes |
| 8 | Not measured yet | `gap` callout + gap badges | only when a gap applies; ≤ 3 lines |
| 9 | Remember | `remember(…)` | ≤ 3 bullets |
| 10 | The checklist | literal `#heading(level: 2, "The checklist")` + `checklist-table` | 8 rows; answer ≤ 15 words + status |

The literal checklist heading stays because `book/evidence/check.sh` looks for it.

## Case box

One study, one lesson about the concept.

| Field | Content |
|---|---|
| `study` | "Study 04 · configuration portal" |
| `verdict` | `fits` / `does not fit` / `depends` — shown as a coloured tag |
| `situation` | 1 sentence: the business need |
| `compared` | 1 line: the designs, in plain words |
| `result` | ≤ 2 sentences — **the only place a number may appear** |
| `lesson` | 1 sentence: what this says about when the concept fits |
| `watch-out` | 1 plain limit: the one that most changes how to read the result |
| `evidence` | one or more `claim-badge("…")` |

## Study profile (chapter 12)

One page per study with registered claims (01–05). Studies without claims get one line in a
"Planned studies" list, labelled as not measured.

1. **The business** — 2 sentences.
2. **The question** — 1 sentence.
3. **Designs compared** — table: design id · what it does · concept chapter.
4. **What it showed** — 3–5 bullets, each with a badge.
5. **Not measured** — gap badges.
6. **Reproduce** — study path and the run tags named by the cited claims.

## Writing rules

1. A paragraph has **≤ 3 sentences and ≤ 60 words**. Aim for sentences ≤ 25 words.
2. Prefer a bullet list or a table to a paragraph.
3. One idea per paragraph. Put the point first.
4. Plain words first. Define a term at first use or link it to the glossary.
5. **Numbers live only in case boxes and study profiles**, next to their badge.
6. **Every number in a box appears in the statement or limits of a claim the box cites.**
   Copy it as written; do not compute a new ratio. Never quote a value a claim says did not reproduce.
7. One caveat per box. The full limits stay on the card in the evidence registry.
8. No cross-study numeric ranking. Analogies keep the `analogy` callout.
9. No reader exercises (methodology 21).
10. Explain a thing once; elsewhere, link to it.
11. Break a dense explanation into short labelled steps, a comparison table or
    small callouts. New or complex concepts need a visual showing the relationship,
    sequence or failure boundary, plus a brief text equivalent. A diagram must teach
    the mechanism; decorative images do not replace an explanation.

## Evidence on the page

- `claim-badge(id)` renders one line: claim id · plain strength · gap kind (if any) · link and page
  of the full card.
- Plain strengths: `run_replicated` → *replicated*, `single_run_directional` → *single run —
  direction only*, `mechanism_supported` → *mechanism*, `analogy` → *analogy*, `gap` → *gap* plus
  its kind. A strength the registry adds later must get a label, or the build fails.
- Full cards appear **only** in the evidence registry index.
- Gaps are always visible on the page: a `gap` callout with the badge.

## Rollout

| Phase | Task | Exit |
|---|---|---|
| 1 Pilot | macros, front page, new skeleton, chapter 4 Rollup, Study 04 profile | owner approves the style |
| 2 Conversion | every other chapter and profile; old files removed | 31/31 claims cited; rules met |
| 3 Release | rebuild `book/dist`, CHANGELOG, tag `repo/data-architecture-book-v3` | release manifest clean |
