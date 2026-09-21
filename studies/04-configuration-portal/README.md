# Study 04 — configuration portal

How should a configuration portal store and serve configuration, and what does each choice cost?

A configuration portal is a service other products use: they create, modify, publish and retrieve
their configuration through it and its database. This study builds the same portal in a catalogue
of eighteen designs and measures what each one costs to read, to write, to keep correct under
concurrency, and to hold. **Ten are implemented and dev-checked; eight are mapped and not yet
built**, and the coverage gap is stated in the handoff and the analysis rather than implied away.

The study's specification of record is [`HANDOFF.md`](HANDOFF.md) (`EH-04` revision 1); the
step-by-step record is [`PROGRESS.md`](PROGRESS.md); unmapped decisions and the controls that
depend on them are in [`ESCALATIONS.md`](ESCALATIONS.md).

## The domain

| Term | Meaning |
|---|---|
| Product definition | a product type or catalogue entry. Owns no configuration. |
| Installed product | one deployed instance of a product definition, in an environment, optionally attributed to a business unit. |
| Environment | the environment containing an installed product. |
| Business unit | an optional ownership or scope dimension. |
| Configuration entry | one named configuration key and its value. |
| Effective configuration | the complete current configuration returned to an installed product. |
| Configuration revision | the monotonic version of an installed product's effective configuration. |
| Publication | making a new configuration revision current. |

**Configuration belongs to an installed product, not to the product definition.** Two installations
of the same definition in the same environment — one per business unit — have entirely separate
configuration, and INV-11 exists to prove they never cross.

## The question, in five reads and six writes

| Read | | Write | |
|---|---|---|---|
| `r01_effective_config` | the complete effective configuration | `w01_modify_key` | modify one key |
| `r02_read_key` | one key | `w02_add_key` | add one key |
| `r03_revision_check` | is there a newer revision? | `w03_delete_key` | delete one key |
| `r04_list_installations` | the portal overview | `w04_publish_batch` | atomically publish a batch |
| `r05_search_key` | one key across installations | `w05_replace_all` | atomically replace everything |
| | | `w06_update_metadata` | update metadata, independently of configuration |

## The designs

| ID | Design | Isolates | Pairs with |
|---|---|---|---|
| `n1_rows_indexed` | row per key, indexed | **the reference** | — |
| `n0_rows_unindexed` | the same, without the search index | the index | `n1` |
| `n2_rows_rolldown` | parent keys copied onto every row | rolldown | `n1` |
| `n3_rollup_trigger` | count, revision, timestamp, hash maintained by a trigger | rollup maintenance | `n4` |
| `n4_rollup_app` | the same rollups, maintained by the application | who maintains them | `n3` |
| `d1_doc_row` | one JSON document row per installation, replaced whole | document storage | `n1` |
| `c1_optimistic_version` | version-checked write with bounded retries | optimistic arbitration | `c2` |
| `c2_pessimistic_lock` | lock the parent, then read and write | pessimistic arbitration | `c1` |
| `x1_lost_update_control` | **negative control**: unchecked read-modify-write | the guard itself | `c1` |
| `x2_rollup_drift_control` | **negative control**: rollup maintained after commit | atomic maintenance | `n4` |

`x2`'s SQL files are byte-identical to `n4`'s. The only difference is where the harness calls
`w_recompute_rollup`. That is the strongest form of a controlled pair this catalogue can express.

**Mapped but not yet implemented** (see `HANDOFF.md` §3 and the amendment in `PROGRESS.md`):
`d2_doc_on_parent`, `d3_doc_sections`, `d4_doc_jsonb_path`, `h1_rows_plus_readview`,
`s1_snapshot_pointer`, `s2_append_history`, and the colocated/non-colocated pair `y1`/`y2`.

## Correctness first

Thirteen invariants (INV-1…INV-13, `HANDOFF.md` §4) are verified **in Go from the generated
dataset and from the operations the client saw acknowledged** — never by comparing one query to
another, which would only prove the queries agree with each other. A cell that fails a gate
reports no timing.

The gate earns its place. Building this study it caught three harness bugs before a single number
was reported:

1. the expected configuration was compared against the read result **by position**, while the read
   returns keys in alphabetical order and the generator produces them in section order — which
   reported the reference design as broken;
2. two statements used `$1` without declaring `-- params:`, so 2 263 write operations ran against a
   server answering "there is no parameter $1" and every one of them was counted as an error;
3. the lost-update detector compared the stored counter against the **last value written** — a
   number with itself — and cheerfully reported 583 acknowledged increments and a counter of 39 as
   a pass.

**Negative controls must fire.** `x1` loses updates that were acknowledged, and `x2` lets the
stored rollup drift from its entries. If either does not fire, the associated correctness claim is
recorded as *not demonstrated*, the workload is hardened, and the run repeats. The invariant is
never weakened to make a run pass.

## Running it

Everything runs in Podman; nothing is installed on the host.

```bash
./probe-bind-mount.sh                      # phase B: does the worktree reach the container?
./probe-engines.sh pg-single,yb-single     # phase A: what do the engines actually do?
./run-study.sh --scale tiny --topologies pg-single        # dev checks; not results
./run-study.sh --tag --scale small --topologies pg-single  # a reported matrix
```

Every runner takes the shared benchmark lock (`infra/lib.sh`), so two studies never measure on this
machine at once. A run writes `results/<run-id>/`, a generated report to `reports/<run-id>.md`, and
an inputs digest over the result bytes.

## What the numbers mean, and what they do not

`docs/environments/host-zenbook-ux5406sa.md` bounds every number here: eight heterogeneous cores
shared by the client and the database, no real network, a laptop that throttles. Relative
comparisons between designs are the point; absolute throughput is not a capacity plan. Deep tails
come from a closed-loop generator and are optimistic flooring, compared between designs and never
quoted as SLOs. The cadence phase reports an offered rate that is a **calculation** — installed
products divided by the update period — and any compressed-rate run is labelled as the
compressed-time validation it is.

This study was planned, implemented, measured, validated and analysed by **one** model. That is a
limitation of its conclusions, stated in the analysis, and a second analyst's view is invited.
