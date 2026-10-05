# Book: appendix of the EXPLAIN ANALYZE plans behind the SELECT claims

**Task:** `task-20261005T190500Z-book-plan-appendix` · **capability:** HIGH · **benchmark:** none

## Owner request (2026-10-05)

"I believe the tests we have performed analyze explains on the database queries. could we have an appendix with them?"

## Decisions (HIGH; the owner asked not to be asked)

- Scope: the 18 measured SELECT claims (`v6-sel-01`..`18`). Every SELECT cell already stores one `EXPLAIN ANALYZE` execution in its result JSON (`cells[].Plan`), captured after the correctness check and before the timed samples. Studies 02/03 write plans too but are outside this first appendix (follow-up).
- One representative run per exhibit: a registry anchor of the claim, normally the first run of the set; never a new measurement.
- Mechanical extraction in Go (`tools/evidence plans`): verify claim, anchor and sha256 digest; keep node lines and decisive attributes; drop Output, Buffers, YugabyteDB per-request metrics and estimates, and cost estimates on node lines. The full plan stays in the committed result file, whose path and axis are printed.
- The generator checks every plan-shape statement the reading notes make (same / different node tree).
- Notes name the decisive node and do not add conclusions beyond the cited claim.

NEXT MODEL: HIGH
