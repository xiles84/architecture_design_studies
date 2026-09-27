# Diagnosis: why does `l3_section_sharded` never finish its load? (study 03)

**What is known, without further measurement.** `l3_section_sharded` has failed its load in all three
passes of the `v2-11` replication:

| Pass | Its position | Envelope | Symptom |
|---|---|---|---|
| `20260926T234227Z` | third (last) | 3 GB / 1.5 GiB | connections terminated mid-`COPY`, SQLSTATE 57P01 |
| `20260927T111110Z` | first | 3 GB / 1.5 GiB | connections terminated mid-`COPY`, SQLSTATE 57P01 |
| `20260927T142039Z` | first | 5 GB / 3 GiB | `dial tcp 10.89.0.41:5433: connect: connection refused` |

Position does not explain it (it has failed first, last and in between). The memory envelope does not
explain all of it (a raised envelope changed the symptom rather than removing it). Two symptoms point at two
mechanisms: a database that kills its own backends under memory pressure, and a database that is **not
listening** when the load starts.

**This diagnostic separates one pair of explanations.** Running `l3` **alone** — one design, one topology,
no other cell before or after it — distinguishes:

- **H1, intrinsic**: `l3`'s own load cannot complete in this harness at this scale, whether or not anything
  runs before it. Then the sharded half of `v2-11` is blocked by a harness or schema defect, which is a
  finding about the tooling rather than about the design.
- **H2, cumulative or sequencing**: `l3`'s load completes when it is the only cell. Then the failure is
  about the state the earlier cells leave behind (memory, connections, a restart), and the claim's sharded
  half may still be measurable with a dedicated run.

**Envelope.** Database container 5 GB, tserver limit 3 GiB — the same as the pass in which `l1` and `l2`
both succeeded, so that resource pressure is not the variable being tested and `l3` gets the most
favourable conditions this host can offer it.

**Host state is recorded** by `record_host_state()` beside the results, and the database log is kept, so a
third symptom would still be attributable.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| D1 | **H2: `l3` loads successfully when it is the only cell**, because two earlier cells in a session leave the database in a state its load cannot survive. | Its load failing with nothing before it, which makes H1 true and the sharded half unmeasurable here rather than difficult. |
| D2 | **If it loads, its numbers are usable.** All read questions pass verification against Go-computed truth, and the "seats available per section" questions come out far slower than `l1`/`l2`'s, as the 2026-09-21 run recorded (15.9 and 28.3 ops/s there). | It loads but fails verification, which would make the design's answers wrong under this schema rather than merely slow. |
| D3 | **The host state shows no memory exhaustion**, since only one cell runs. | Memory or swap near exhaustion at the snapshot, which would send the memory question back to the topology rather than to the sequence. |

Recorded before measurement. This is a **diagnostic**: it produces evidence about the harness and about
whether an existing claim can be tested, not a new claim. Whatever it shows, it will not be used to change
`v2-11`'s registered numbers — only, if it succeeds, to note that the sharded half is measurable by a
dedicated run.

## Scope

| Field | Value |
|---|---|
| Runs | one, `--designs l3_section_sharded` only, `--trials 3`, `yb-single`, `small`, phases `verify,read` |
| Compared with | the three failed cells above |
| Boundary | One host, one engine, one scale. A success here does not recover the failed cells; it says they were lost to sequencing rather than to the design. |
