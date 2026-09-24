# Diagram fidelity audit — Study 03 (reserved seating)

**Audit revision:** `5b1f8df0c0ef86c7a894b84fa44080127610b9db`.
**Renderer:** `infra/diagram-render.sh`, `PLANTUML_IMAGE` digest
`sha256:9b9ee6af54a86ea6ac53805e9da4da7fd913711cfe6c19b6ac1a41c398223aeb` (PlantUML 1.2026.8).
**Scope:** the study's 5 `.puml` sources (2 structural, 3 sequence). Each factual annotation is traced to
its SQL or signed analysis and classified current / revision-scoped / stale / unsupported. No source in
this study needed correction; no measured result, report or signed analysis was altered.

| Source | Factual annotation(s) | Verdict | Evidence at the audit revision |
|---|---|---|---|
| `00_overview.puml` | venue → section → seat and the four design families | current | `studies/03-reserved-seating/README.md`; `sql/` layout families |
| `s_arbitration.puml` | two buyers on the same marked seats; Study 02's `SKIP LOCKED` does not apply because a marked seat may not be skipped | current (design reasoning, not a measurement) | `sql/s0_check_then_hold_rc` … `sql/s4_check_then_hold_serializable`; study README |
| `l_layout.puml` | L1 pre-created rows (publish 100,000 seats = 100,000 rows); L2 section document (one row per section, holds rewrite it); L3 section-sharded placement | current; the 100,000-seat figure is the study's published largest tier | `sql/l1_claim_rows`, `sql/l2_section_document`, `sql/l3_section_sharded`; `studies/03-reserved-seating/README.md` (1,000 / 10,000 / 100,000 tiers) |
| `e_expiry.puml` | expiry designs and the "sweeper stopped minutes 45–85" outage window | current (scenario parameter) | `sql/e0_app_clock_expiry`, `sql/e1_sweeper_expiry`, `sql/e2_cart_expiry`; study README/harness phase |
| `k_checkout.puml` | K1 payment window; S1 checked confirmation refuses after paying; K0 naive confirmation is a negative control that sells over an expired hold; S1r retries a transient refusal seen on YugabyteDB | current | `sql/k0_naive_confirm`, `sql/k1_payment_window`, `sql/s1_conditional_update`, `sql/s1r_confirm_retry`; the S1r refusal is ER-01 and v2-10 |

The three sequence sources describe application order and controls, not observed performance; the K0
control is labelled as a control and fires as documented.
