# Progress — Study 04 v2 remaining designs

**Task:** `task-20260922T195241Z-study04-v2-remaining-designs`
**Attempt:** `attempt-20260923T093624Z-dc62ef` (claim `claim-d0fd650bcd112190`, epoch 1)
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-remaining-designs`

Measured: one new design (d2) and its comparison to d1. No parallel measurement.

## Decision Log

1. **Take the "build what is SQL-only, record the rest as a gap" branch the acceptance allows.**
   Seven of the eight designs need harness changes (new write modes or Kinds) or an unverified
   engine mechanism; implementing them hastily would violate the correctness standard. Confidence:
   High.
2. **`d2_doc_on_parent` is genuinely SQL-only.** The harness drives document designs by statement
   name (`wd1_read_doc`, `wd1_write_doc`); d2 ships the same names against `installed_product`
   columns. Verified by a tiny dev check that passed the gate before the reported run. Confidence:
   High.
3. **`d4_doc_jsonb_path` is not SQL-only.** The document write path hands the harness's new whole
   document to `wd1_write_doc`; a per-key `jsonb_set` needs a new harness write mode. Recorded as a
   gap. Confidence: High.
4. **`d3_doc_sections`, `h1`, `s1`, `s2` need new Kinds/read paths** (section assembly, readview
   maintenance, snapshot pointer, history-plus-current). Recorded as gaps. Confidence: High.
5. **`y1_colocated` needs the colocation DDL verified on the pinned YugabyteDB 2025.2.6 image and a
   placement-recording step**; running `y2` alone proves nothing. Recorded as a gap. Confidence: High.
6. **The equal-total arm needs a runner flag** to split the total budget across nodes; not built.
   Confidence: High.
7. **The d2 result is reported as measured, including the unconfirmed hypothesis.** The predicted
   W6 amplification did not appear at small cardinality; that is the finding, not a reason to
   re-run until it does. Confidence: High.

## Residual risks

- d2 was measured at one cardinality and one run; the amplification may appear with larger documents.
- The unbuilt seven leave the v2 coverage matrix's "remaining designs" cell partly open.
