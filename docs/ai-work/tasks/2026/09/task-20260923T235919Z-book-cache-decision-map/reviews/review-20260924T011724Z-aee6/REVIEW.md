# Review — book cache decision map

**Task:** `task-20260923T235919Z-book-cache-decision-map`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Verdict:** **approved for integration** — no rework required, one recorded exception carried forward
**Required tag:** `repo/book-cache-decision-map-v1`

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | guarantee-first flow, topology/fill path last | accepted — the acceptance's flow order exactly |
| 2 | four-status protocol table; lease = newer signed but unregistered | accepted — the lease path is demonstrated without a registered rate, so the status is honest |
| 3 | relaxed defined against the strict comparator; impossible forbidden separately | accepted |
| 4 | process-local bypass stated without inventing a claim | accepted with the recorded exception (below) |
| 5 | add-cache gain keeps `db-only`; equal-total arm beside it, never pooled | accepted |
| 6 | 2026-09-23 evidence added as registered v3 claims; hard-expiry zero as a gap | accepted |
| 7 | `config.typ` unchanged | accepted — no new primitive was needed |
| 8 | no full build; `book/dist` forbidden | accepted — distillate is the next build's job |

## Independent checks re-run by the reviewer

```text
typst compile --root book main.typ (pinned image, book/ only)   exit 0, no warnings
every lock-avoidance sentence names a race                      yes (table + chapters)
protocol table has the four status values                       yes
relaxed defined against strict comparator                       yes
bypass counts (0 hits, 15,578–21,756 reads) vs report           match reports/20260921T-survey3.md
other numbers vs v3 claim statements                            match (v2-14/15/16, v3-01..v3-06)
files changed                                                   cache concept + 4 chapters + archive
```

## Recorded exception (accepted, not hidden)

The two bypass figures are the anchor report's numbers, not the statement of registered claim `v2-15`.
Promoting them into the claim requires editing `book/evidence/`, which is a `forbidden_path` of this
task. They are cited beside the claim and flagged for a later evidence task. This is the honest
resolution of the task's own tension between "state the bypass counts" and "every number resolves to a
registered claim"; it does not misstate the evidence.

## Findings (non-blocking)

1. The committed `book/dist/` PDF predates the rewrite; the prose reaches a PDF at the next build.
2. Deletion/tombstone and outbox protocols remain proposed / unmeasured by design.

## Verdict

Approved for integration. The cache material now teaches the guarantee before the pattern, names every
race, and cannot be read as saying a private cache is safe or that relaxed means dirty.
