# Execution brief — cache chapter: cards instead of a collapsed table, and the terminology fixes

Implement findings **R01** (the broken §15.2 table) and **R13–R18** of the external review of the
*Edition 1* draft. R01 is the review's most defensible production defect: five pages of the book's
most conceptual section were hyphenated into vertical fragments because a table column collapsed.

## Goal

`book/concepts/cache-consistency.typ` must state the freshness contract formally, name each mechanism
with what it closes *and* what it does not solve, and stop blending source arbitration with copy
ordering.

## Decisions already made (do not re-litigate)

- **Cards, not a resized table.** One card per mechanism: Race closed / Mechanism / Guarantee / Does
  not solve / Evidence. The mandatory "does not solve" line is the point — several of these
  mechanisms can be present and the contract still broken.
- **Grouping is R16.** Source correctness (database mutation lock, conditional update) arbitrates the
  mutation; copy correctness (invalidation, publication fence/CAS, source version validated at the
  hit boundary, fill lease, TTL and early expiry, transactional outbox) orders what leaves the
  transaction. Do not merge them into one list again.
- **The evidence-status vocabulary is unchanged.** The table carried four classes — active registered
  claim / newer signed but unregistered evidence / mechanism demonstrated / proposed-unmeasured — and
  every card keeps the class it had. The fill lease is "newer signed but unregistered"; the outbox is
  "proposed / unmeasured" and must stay visibly unmeasured.
- **R13.** Invalidation is *not* a fence. Invalidation removes or marks an entry stale; it cannot stop
  an in-flight reader from publishing an older value. The publication fence is the extra check. The
  chapter's own stale-fill race is the counter-example and should say so.
- **R15.** Write the strict-after-acknowledgement contract formally with `W1`, `ack(W1)`, `R1`,
  `begin(R1)`, including the "unless superseded" caveat.
- **R14.** Name the harness counter (`wrong_reads`) and say it counts strict-comparator violations over
  committed-but-older values, never corruption; flag that the bounded-window question needs a bounded
  contract to be meaningful.
- **R17/R18.** The outbox is a prerequisite for strictness, not a guarantee of freshness, and a
  purpose matrix should say per mechanism whether it addresses correctness, freshness, load or
  recovery.
- **No new numbers.** Every figure in the chapter stays attached to its existing registry card.

## Constraints

- Do not edit `book/evidence/**`: the registry is owned by `task-20260924T102917Z-book-evidence-v4`
  and already carries v4.
- `book/lib/config.typ` is shared: add the mechanism components there rather than inventing a second
  visual language in the chapter, and keep the existing callout and marker definitions untouched.
- Do not commit `book/dist/`; build to a scratch path and delete it.
- Never edit another analyst's report or analysis.

## Acceptance criteria

See `task.json`. The gate that matters most is the measurable one: count orphan-hyphen lines in the
extracted page text before and after (Edition 1 had 99, with 17 on the worst page) and record both
numbers. A green compile is not evidence that the page is readable.

## Escalate only if

- a mechanism's evidence status cannot be preserved without inventing a registry claim;
- the card layout cannot be made readable without changing the visual system in a way that would
  restyle the whole book.
