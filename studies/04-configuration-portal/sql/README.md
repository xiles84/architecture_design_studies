# SQL catalogue — study 04

One directory per design. Every design has the same five files:

| File | Contents |
|---|---|
| `schema.sql` | the common tables and this design's tables, applied **before** the bulk load; its header explains the design |
| `indexes.sql` | secondary indexes — and, for `n3_rollup_trigger`, the trigger — built **after** the load |
| `queries.sql` | the read catalogue: `r01`…`r05`, identical names and result shapes in every design |
| `writes.sql` | every statement the portal's write operations issue |
| `audit.sql` | how to recompute the truth from this design's own tables, for the reconciliation audit |

Named statements use the repository's catalogue format (`-- name:`, `-- params:`). The Go
harness decides **which** statements run, in **which** transaction, at **which** isolation level
(`harness/designs.go`); the SQL files decide **what** each statement does.

## The read contract

| Name | Question |
|---|---|
| `r01_effective_config` | the complete effective configuration of one installed product |
| `r02_read_key` | one key |
| `r03_revision_check` | is there a newer revision than the caller holds? |
| `r04_list_installations` | the portal overview: definition, environment, business unit, count, revision, last modification |
| `r05_search_key` | one key across installations of a definition, an environment or a business unit |

All five return the same shapes in every design, including the document design, which unrolls its
JSON object into key/value rows. A difference between two designs is then a difference in
mechanism and not in what was asked.

## Designs that differ by one decision — check with `diff -r`

```bash
diff -r sql/n1_rows_indexed sql/n0_rows_unindexed      # indexes.sql only: the search index is gone
diff -r sql/n1_rows_indexed sql/n2_rows_rolldown      # the copied parent keys, and r05 without a join
diff -r sql/n3_rollup_trigger sql/n4_rollup_app       # the trigger, and one extra write statement
diff -r sql/n4_rollup_app sql/x2_rollup_drift_control # nothing: the control is the harness's transaction boundary
diff -r sql/n1_rows_indexed sql/c1_optimistic_version # the guarded read-modify-write
diff -r sql/c2_pessimistic_lock sql/x1_lost_update_control  # wc2_lock_installed_product, and nothing else
```

`n0`, `n2`, `n3`, `n4`, `c1` and `c2` are whole designs read as designs; `x2` is the strongest
statement of "one decision" in the catalogue, because its files are byte-identical to `n4`'s.

## Conventions

- **Every publication ends with `w_revision_bump` in the same transaction as its change.** That is
  what makes "the client was told it published" and "the database holds it" the same event
  (INV-3, INV-12).
- **`w06_update_metadata` never moves the revision.** Ordinary metadata is not a publication, and on
  a design that embeds configuration in the parent row it must not rewrite the configuration
  either — that is the measurement.
- **Parameters in `INSERT … SELECT` and `unnest` are cast explicitly** (study 02's rule).
- **A statement that uses `$n` must declare `-- params:`.** The harness refuses to load a design
  where it does not, because the failure is otherwise silent until 2 263 statements have run
  against a server that answers "there is no parameter $1" every time.
- **`w_load_*` statements run once, after the bulk copy**, from the `load_entry` staging table, so
  no design pays a per-row load cost the others avoid.
- **Index column order is chosen for YugabyteDB's defaults**, which hash-shard the first column.
- **`x1` is a negative control.** It is measured like every other design and its speed is never
  reported as a plain number: a design that loses updates quickly has not been fast.
