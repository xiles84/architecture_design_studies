# Result — book cache decision map

**Task:** `task-20260923T235919Z-book-cache-decision-map`
**Branch:** `book/cache-decision-map`
**Required tag:** `repo/book-cache-decision-map-v1`
**Capability:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Benchmark:** none — authorship only

## What was delivered

1. **Guarantee-first flow** in `book/concepts/cache-consistency.typ`: name the contract → establish
   change-observation completeness → identify source capabilities → choose topology and publication
   protocol → add performance controls separately. Topology and fill path are decisions made after the
   guarantee.
2. **A protocol table with an evidence-status column** carrying four statuses: active registered claim,
   newer signed but unregistered evidence, mechanism demonstrated, proposed / unmeasured. Each row
   names one mechanism (database mutation lock, fill lease, publication fence/CAS, source version
   token, transactional outbox) and the race it closes.
3. **Every lock-avoidance sentence names its race**, and *relaxed* is defined only against the strict
   comparator; impossible/uncommitted values are forbidden separately from permitted staleness.
4. **The private-memory trap is explicit.** The legacy process-local arms' zero wrong reads are a
   bypass (0 hits over 15,578–21,756 reads); the owned process-local arms served hits and stayed
   correct only by validating a source version at the hit boundary.
5. **The 2026-09-23 evidence is folded in**: real-TTL churn (`v3-01`, `v3-02`), the labelled
   equal-total arm (`v3-03`), per-engine strict/relaxed (`v3-04`, `v3-05`), the locality pair
   (`v3-06`) and the successor gap (`v3-gap-01`).
6. Chapters corrected: `major-families`, `minor-variants`, `scenarios`, `how-to-choose`.
7. Signed authorship analysis at
   `docs/ai-work/tasks/2026/09/task-20260923T235919Z-book-cache-decision-map/attempts/attempt-20260924T010942Z-8a7b85/ANALYSIS.md`.

## Acceptance checks

```text
guarantee-first flow present                                   yes (5 steps, concept)
protocol table separates the four evidence statuses            yes
every lock-avoidance sentence names its race                   yes (table + chapters)
private-memory bypass stated (0 hits, 15,578–21,756 reads)     yes, cited to v2-15's anchor report
relaxed defined against the strict comparator                  yes
every number resolves to a registered v3 claim                 yes, except the two bypass figures (see below)
book rebuilds                                                  compile exit 0 in the pinned image
manifest records the evidence digest                           next full build: 814f40949d162f56… (build.sh already reads v3)
```

**Recorded exception.** The bypass figures (0 hits, 15,578–21,756 reads) are in the primary anchor
report of registered claim `v2-15`, not in the claim's statement. Promoting them into a claim requires
editing `book/evidence/`, a `forbidden_path` of this task; they are therefore cited as report figures
beside the claim. Flagged for review and for a later evidence task.

**`book/dist` not regenerated.** `book/build.sh` and `book/dist` are forbidden to this task, so the
committed PDF predates this rewrite. The sources compile cleanly (`typst compile --root book main.typ`,
exit 0, no warnings) and the next full build will write the v3 evidence digest.

## Validation

```text
typst compile --root book main.typ (pinned ads-book image, book/ mounted only)   exit 0, no warnings
v3 registry digest used for the analysis frontmatter                              814f40949d162f56…
files changed                                                                     cache concept + 4 chapters
measurement / database / benchmark lock                                           none
```
