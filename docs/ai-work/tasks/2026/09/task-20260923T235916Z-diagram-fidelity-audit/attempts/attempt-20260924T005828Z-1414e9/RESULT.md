# Result — diagram fidelity audit

**Task:** `task-20260923T235916Z-diagram-fidelity-audit`
**Branch:** `repo/diagram-fidelity-audit`
**Required tag:** `repo/diagram-fidelity-audit-v1`
**Capability:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Benchmark:** none — audit and re-render only

## What was delivered

1. **A per-study audit of all 30 diagram sources** (`FIDELITY.md` under each study's `diagrams/`):
   every factual annotation is traced to its SQL/harness/plan/report at revision
   `5b1f8df0c0ef86c7a894b84fa44080127610b9db` and classified current / revision-scoped / stale /
   unsupported.
2. **Three corrections in two sources:**
   - `d3_flattened_fk.puml` — "six of twelve queries lose a join" → "three of the current sixteen
     (q02, q08, q15)"; and "index-only scan, never touching the heap" scoped to a prepared-read
     result with the zero-fetch and 31,293-fetch trial-1 numbers and revision named.
   - `d10_embedded_hybrid_locked.puml` — the unsupported "per 30 s" window replaced by "per
     concurrent read/write split", citing `reports/20260913-d9-cache-race.md` (digest
     `70f87f5d5667bff8`); the supported count and the two-bug mechanism are kept.
3. **Re-rendered corrected sources** with the pinned helper; only the D3 and D10 SVGs changed bytes,
   the other 18 study-01 SVGs are byte-identical.
4. `CONTEXT.md` records the audit.

## Classification summary

```text
study              sources  current  revision-scoped  stale  unsupported
01-charity-tree       20       18*          1            1        1
02-ticket-booking      5        5           0            0        0
03-reserved-seating    5        5           0            0        0
```
`*` D3 contributes one stale + one revision-scoped annotation in a single source; D10 one unsupported
unit. Study 01's `d4` "seven of twelve" is the revision-scoped count and is left as written; `d7`'s
"40 donations" is current-but-illustrative.

## Validation

```text
audit revision                                  5b1f8df0c0ef86c7a894b84fa44080127610b9db
d3/d10 join and heap evidence                   sql/d2|d3/queries.sql, run 20260913T125342Z-v3 discussion
d10 corruption evidence                         reports/20260913-d9-cache-race.md, digest 70f87f5d5667bff8
re-render                                       pinned helper, digest 9b9ee6af…, only d3+d10 SVGs changed
other 18 study-01 SVGs                          byte-identical after re-render
measured results / reports / analyses altered   none
```

## Conclusion

The committed diagram corpus now says only what its cited artefacts support, with the preparation and
revision conditions stated where a result is conditional. No measurement was performed and no database
or benchmark lock was touched.
