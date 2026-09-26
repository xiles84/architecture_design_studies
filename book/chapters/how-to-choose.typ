#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("How to choose")

#marker("chapter", "how to choose")
This chapter is the decision path. It front-loads the two questions that dominate every later
trade-off: what must the schema be able to answer, and who arbitrates the contended row.

#figure-evidence(
  "../assets/fig-decision-tree.svg",
  "decision tree",
  "conceptual illustration",
  "The decision path, in the order the decisions bind.",
  [What must the schema answer — after an undo, a cancellation or a correction? If the current
   representation cannot answer it, change the representation before the speed. If a row is
   contended, arbitrate it: a conditional update, a lock plus an invariant recheck, or pre-created
   rows. If it is not, leave the row canonical and index the question. Name the maintainer of any
   derived read. State the evidence level and the gap before quoting a number. The figure orders the
   chapters' recommendations and asserts no measured value of its own.],
)

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
the guarantee structural. Application maintenance kept writes near the reference in the measured
configuration-portal workload, and its obligation is stated
against where the derived state lives: when the base fact and the rollup are in the same database,
update both in one transaction; when they cannot share one atomic transaction, name the delivery,
idempotency and reconciliation mechanism that keeps the derived state recoverable. Only the second case
is *publication*, because only there does a copy escape the transaction that wrote the original — leaving
a moment when a reader can see the new value while the copy is still old.
The measured trigger cost was 4.0–5.7x on whole-configuration replacement across three runs (median
4.94x) — a range rather than a single figure, because repeating the measurement inside that configuration
gave results that disagreed with each other by more than we would like.

#heading(level: 2, "Step 4 — say how strong the evidence is")
Every recommendation states how well it is supported, what it could not separate, and what has not been
measured. A recommendation that does not say what is missing claims more than the evidence allows.

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
