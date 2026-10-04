# Colocated / non-colocated placement (review R1 / handoff C4)

LOW `deepseek-flash` / opencode CLI, 2026-10-04. Feature check with error
evidence; no placement success is invented.

## What was tried

The domain tables accept an explicit YugabyteDB placement:
`JoinDomainSchemaSQLPlacement` appends `WITH (colocation = true|false)` to
`charity`, `person` and `donation`. `--placement colocated|noncolocated`
selects it; the effective placement is recorded in provenance settings and
therefore in the compatibility class. An empty value keeps the engine default,
so every earlier run is unchanged.

## Feature check on the pinned engine (2025.2.6.0-b111)

| run | tag | outcome |
|---|---|---|
| `select-join-topo-yb-colocated-20261004` | `run/study-01/select-join-topo-yb-colocated-20261004` | **rejected**: `ERROR: cannot set colocation true on a non-colocated database (SQLSTATE 0A000)` from `CREATE TABLE … WITH (colocation = true)`; recorded as the run error, no cells |
| `select-join-topo-yb-noncolocated-20261004` | `run/study-01/select-join-topo-yb-noncolocated-20261004` | **supported**: cluster3, RF3, domain arm, bounded subset, 30/30 valid; placement `noncolocated`; endpoint ops 123/122/122 |

The cluster's default `yugabyte` database is non-colocated, so
`colocation = true` on a table is rejected. `colocation = false` is accepted and
is the supported baseline.

## Conclusion and gap

- The **non-colocated** arm is measured and is the baseline physical placement
  for the cluster comparisons.
- The **colocated** arm is an explicit coverage gap: it cannot be expressed on
  the default database. Expressing it requires a database created with
  `COLOCATION = true` (or a colocation tablegroup) and a client DSN pointing at
  it; `infra/yb-cluster3.sh` does not create such a database.
- The gap is recorded with the exact engine error; the topology script change to
  add a colocated database is the next step, not an assumed result.

## Report fix found by this check

An early failure (DDL/feature check) has no cells and no feasibility evidence.
`RenderJoinRun` used to refuse to render such a run, so the feature-check error
was written to `join.json` but could not be read back. It now renders a report
with a **Run error** banner, and the failed colocated run has a readable report.
