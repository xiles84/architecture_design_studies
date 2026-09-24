# Diagram fidelity audit — Study 02 (ticket booking)

**Audit revision:** `5b1f8df0c0ef86c7a894b84fa44080127610b9db`.
**Renderer:** `infra/diagram-render.sh`, `PLANTUML_IMAGE` digest
`sha256:9b9ee6af54a86ea6ac53805e9da4da7fd913711cfe6c19b6ac1a41c398223aeb` (PlantUML 1.2026.8).
**Scope:** the study's 5 `.puml` sources (4 structural, 1 sequence). Each factual annotation is traced to
its SQL or its signed analysis and classified current / revision-scoped / stale / unsupported. No source
in this study needed correction; no measured result, report or signed analysis was altered.

| Source | Factual annotation(s) | Verdict | Evidence at the audit revision |
|---|---|---|---|
| `00_overview.puml` | the cross-row `sold ≤ capacity` invariant; the list of mechanisms that can serialise the last seat; "fourteen ways" | current | `studies/02-ticket-booking/README.md` (14 strategies incl. two controls); `sql/*/`; signed analyses (v2-06) |
| `c_created.puml` | C1 count (RC) is a negative control that overbooks; C2 SERIALIZABLE aborts with 40001; C3 event FOR UPDATE; C4 guarded counter; C5 `MAX(seat_no)+1` with `ON CONFLICT DO NOTHING` | current | `sql/c1_count_naive` … `sql/c5_seat_unique`; the C1 overbook and the two control firings are recorded in `reports/analyses/20260913T021206Z--claude-opus-5--2026-09-13.md` and v2-06 |
| `h_holds.puml` | H0 naive confirm re-sells a seat (negative control); H1's `expires_at > now()` in the confirming statement is the only difference; the database clock decides | current | `sql/h0_hold_naive_confirm`, `sql/h1_hold_checked_confirm`; conflict C of the Study 03/02 work and v2-06/v2-10 |
| `p_precreated.puml` | P1 lock-first, P2 skip-locked, P3 compare-and-set, P4 skip+counter, X1 CAS+ledger; the ledger makes "sold"/"cancelled" inseparable from the ticket's state change; pre-creation cost at 100,000 seats | current (the 100,000-seat figure is the study's published large-event scale) | `sql/p1_precreated_lock_first` … `sql/x1_cas_ledger`; `studies/02-ticket-booking/README.md`; v2-07/v2-08 |
| `r_extra_tables.puml` | R1 inventory row isolates the organiser's edit from the sale row; R2 shards the counter into at most `min(8, capacity)` buckets; R3 pre-created seat tokens | current | `sql/r1_inventory_row`, `sql/r2_inventory_buckets`, `sql/r3_seat_pool`; the bucket bound is `min(8, capacity)` in the schema |

All five sources are design-contract and sequence descriptions; none carries a measured performance
number, so none is a "observed result" figure. The two deliberate negative controls (`C1`, `H0`) are
labelled as controls on the page and fire as documented.
