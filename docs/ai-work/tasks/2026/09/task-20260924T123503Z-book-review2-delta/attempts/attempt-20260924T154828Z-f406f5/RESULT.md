# Result — review-2/3 delta: capsules, lifecycle rendering, generated legend, lease-versus-fence

**Task:** `task-20260924T123503Z-book-review2-delta`
**Branch:** `book/review2-delta`
**Required tag:** `repo/data-architecture-book-review2-delta`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — the only container started is the pinned PlantUML renderer

## What was delivered

**Capsules on repeat citations.** A claim renders its full card where it is first established and a
one-line capsule on every later citation: id, strength, **every** limit, and a pointer to the full card
in the registry index. 34 citations collapsed and the book went from 56 to 54 pages. The evidence index
calls the full card directly, because that is where a reader is sent.

**The counter is Typst state in document order**, not a build-generated citation index. My first design
scanned the sources for `registry-card` calls and decided which was first; that has a second source of
truth (the scan's idea of document order) which nothing would fail on when it drifted. The state counter
is the page's own order. The vestigial build input was deleted rather than documented.

**No limit line is dropped.** My first capsule carried only `limits.at(0)`. Several claims have three
limits, so that silently removed provenance from every repeat citation — the exact failure the book is
written against. It now carries all of them.

**Partially superseded cards.** The card separates **Moved to successor claims** (with the successor ids)
from **Still open in this claim**, and labels the statement as unchanged from its source registry. The
statement itself stays verbatim: the book does not edit claim text, and rewriting it would need a new
registry version. The previous rendering reprinted the old statement above a shorter open list, which made
a moved dimension look still open.

**The page-3 legend is generated.** The kinds of gap listed on the page are read from the active claims'
`gap_kind` values in a fixed order and looked up in a description table; a kind with no description
panics and fails the build. The page that explains the evidence system can no longer list three kinds
while the cards use four.

**Placement ladder.** Copy-out is no longer "the fastest read" or "the only step that adds a consistency
contract". It can be the lowest-latency hot read for the measured cacheable workload, and what is special
is that it crosses the transaction boundary — indexes, rollups and rolldowns carry their own maintenance
correctness.

**The lease-versus-fence figure.** Two panels: a lease that expires under a slow filler and lets the stale
owner publish, and a fencing token that rejects it. It is the visual companion to the fill-lease
correction, and it carries a form and evidence-status label and a text equivalent.

## Verified from committed files

```text
book/build.sh --out dist/preview.pdf   54 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                       evidence digest present, banner absent, no unresolved tokens
capsules rendered                      34 citations of "full card in the evidence registry index"
lifecycle rendering                    "Moved to successor claims", "Still open in this claim" and
                                       "unchanged from its source registry" each render
legend as rendered                     "schema limitation — … · coverage … · instrument … ·
                                       unstable measurement …", in that fixed order
figures                                Figure 6 is the lease-versus-fence panel; eight registered,
                                       all embedded, none carrying an embedded number
book/assets/export.sh --check          passes for all eight
check.sh                               OK, no waivers
```

## Where this task is weak

- **The capsule still repeats a claim's limits verbatim**, so a three-limit claim takes two or three
  rendered lines. Condensing them would be an editorial act on evidence text, and I chose not to; the
  book saves 2 pages rather than the 5 a condensed capsule would give.
- **The capsule loses the confound register**, which the full card carries. A reader who first meets a
  confounded number in a repeat citation now has to follow the pointer to see it.
- **`superseded_by` is rendered as a bare list of ids** with no strength or title, so "Moved to successor
  claims: v4-gap-01, v4-gap-02" asks the reader to look them up.
- **The legend's fixed order is a preference, not a rule**: I chose schema limitation first because it is
  the kind that no run can close. A different editor could reasonably order by frequency.
- **Tranche 2 of the review's visuals is still not implemented** (V2, V3, V4–V6, V9, V10, V12); only V8's
  companion figure was added here because the cache chapter needed it.

## Required tag

`repo/data-architecture-book-review2-delta`, created by `queue integrate` on the integrated state.
