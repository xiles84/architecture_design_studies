#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("How to choose")

#marker("chapter", "how to choose")
This chapter is the decision path. It front-loads the two questions that dominate every later
trade-off: what must the schema be able to answer, and who arbitrates the contended row.

#heading(level: 2, "Step 0 — answerability before performance")
#direct[
  Write down the questions the data must answer *after* an undo, a cancellation or a correction.
  Two measured answerability failures make this a rule rather than a preference: refunds were
  unanswerable in 14 of 15 ticketing designs, and no reserved-seating design could size a hold
  window because an expired hold is erased.
]
#registry-card("v2-07-refund-answerability-and-storage")
#registry-card("v2-gap-02-hold-funnel-unanswerable")
If a required question is unanswerable, no index or cache fixes it. Choose a family that can
answer it — usually append-only history — and accept its maintenance cost.

#heading(level: 2, "Step 1 — place each fact once")
Walk the derive → index → materialize → copy/cache ladder (see the *information placement* concept)
for each question. Prefer the lowest rung that meets the requirement, because each rung up adds
maintenance, storage or a consistency contract.

#heading(level: 2, "Step 2 — name the contended surface")
Find the rows many writers touch at once. If there is one, choose an arbitration strategy
explicitly: pre-created units with compare-and-set, `SKIP LOCKED`, or a version-checked update. The
counter-row design was 4–8x slower than pre-created seats and blocked unrelated edits.

#heading(level: 2, "Step 3 — decide who maintains derived state")
For every rollup or copy, name the maintainer and the failure it must never have. A trigger makes
the guarantee structural; application maintenance keeps writes cheap but must be atomic with its
publication. The measured trigger cost was 6.85x on whole-configuration replacement.

#heading(level: 2, "Step 4 — state the evidence level and the gap")
Every recommendation carries a claim id, its strength, its confounds and its coverage gap. A
recommendation without a gap statement is over-claiming.

#heading(level: 2, "A worked decision")
#list(
  [A ticketing system must never oversell and must report refunds → append-only ledger plus a
   pre-created-seat arbiter, +33–35% storage.],
  [A configuration portal reads its overview far more than it writes → application-maintained
   rollup; a trigger only if a forgotten write is unacceptable.],
  [A public donor portal serves warm profiles → cache with a named contract and a publication fence
   against the stale-fill race (a reader republishing S0 after a writer commits S1 and invalidates),
   accepting that the cache's budget is outside the `db-only` comparison.],
)
