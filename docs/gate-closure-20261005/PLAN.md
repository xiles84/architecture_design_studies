# Closing the T1/T2 gate — HIGH plan (2026-10-05)

Planner: Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH; effort xhigh.
Authority: the owner's empirical-method charter ([`docs/empirical-method-program.md`](../empirical-method-program.md)),
which orders "T1 and T2 before T5" and asks for its tasks to be published. No measurement was run
to write this plan.

## Why now

Study 09 (hierarchy) is the only new domain this machine can run, and it is gated on T1 and T2:

- **T1 is not closed:** three claims are still `single_run_directional` (`v3-01`, `v3-05`, `v3-06`).
- **T2 is mostly open:** the writer sweep is registered (`v2-17`), but the other charter sweeps
  are not done or have fewer than five levels (`v2-18` has three sizes).

Two of the T1 deferrals rested on a false premise. Under WSL2, "the fixed 2 GiB machine" is Podman's
configured profile; saved guest snapshots show 8 CPUs and ≈15.36 GiB. `yb-cluster3` has since
completed 30 placement runs here (Study 01 SELECT). The SELECT join runner also implements a working
tablet/leader instrument (`studies/01-charity-tree/run-select.sh`: `yb_table_properties` plus
`yb-admin list_tablets` before and after each run).

## Decisions (do not reopen in the LOW tasks)

1. **Sizing comes from the effective guest.** Record it per run: 8 CPUs, ≈15.36 GiB, from
   `/proc/meminfo` and `nproc` in a container. Do not use `podman machine inspect`.
   - Reuse each primary run's per-node and client budgets, so the replication compares like with
     like.
   - Where the declared CPU quotas exceed 8, record the oversubscription and per-container
     throttling, and keep going. Shortfall alone is not a reason to stop.
2. **Replication form.** At least three whole runs per target, in seeded randomized or AB/BA order,
   with the executed order read from the manifest. Each run has a position control (one cell first
   and last), and every run is tagged from a clean tree.
3. **Sweep form.** Every T2 sweep:
   - uses at least five levels, chosen to bracket the expected reversal or boundary;
   - runs at least three whole runs;
   - reports per-run ratios with sign counts and a crossing interval only where observations
     bracket it;
   - registers its claim with a `crossover` record (`levels` ≥ 5). If no reversal appears, the claim
     says so (`reversal_observed: false`); that is a result.
4. **Correctness gates timing.** Every cell is oracle-checked, and a negative control must fire in
   each phase. Study 05 first makes its gate fail a cell that serves an *impossible* cache value: the
   2026-10-02 replication recorded one in 1,173,729 reads while the gate passed.
5. **Placement is verified, not labelled.** Any multi-node arm saves `yb_table_properties` and
   `yb-admin list_tablets` before and after each run, adapted from the SELECT runner. A layout whose
   leaders cannot be listed is a coverage gap, not a result.
6. **Registry discipline:**
   - New claims go into `book/evidence/v6` with a ledger entry and the carry-forward test.
   - A claim is `run_replicated` only with a prediction committed before its runs; otherwise it is
     `repeated_controlled`.
   - Gaps a task closes are narrowed in the same change (`v3-gap-01`, `v4-gap-01`, `v2-gap-03` where
     touched).
7. **Book prose** is not edited by these tasks. A later HIGH book task publishes reviewed results.

## Tasks (published with this plan)

| Task suffix | Study | What it closes | Capability |
|---|---|---|---|
| `t1-study05-completion` | 05 | `v3-01` (1,800 s churn), `v3-05` (yb-cluster3 strict vs relaxed), `v3-06` (colocated vs non-colocated, with verified placement); the impossible-value gate | LOW |
| `t2-study04-cardinality-sweep` | 04 | `v2-18` extended from 3 to at least 5 replacement sizes (1 to ≥1,000 entries) | LOW |
| `t2-study01-donations-per-donor-sweep` | 01 | Embedding against normalized and indexed designs across at least 5 donations-per-donor levels | LOW |
| `t2-study01-hot-parent-skew-sweep` | 01 | Rolldown, rollup and embedding under at least 5 hot-parent skew levels; adds a deterministic skew knob | LOW |
| `t2-study05-capacity-sweep` | 05 | Cache benefit against capacity as a fraction of the working set (at least 5 fractions) | LOW |
| `t2-study05-open-loop-knee-sweep` | 05 | At least 5 offered rates bracketing the knee, separating client from server limits | LOW |
| `t5-gate-review` | — | HIGH reviews the six results, decides whether T1/T2 are closed, and releases Study 09 | HIGH |

Measurement is serial under `ads-run-lock`; any order works. Priority: T1 first (62), then the
sweeps (60), then the gate review (58).

## Escalation triggers (for every LOW task)

- A primary configuration cannot start at all on this machine, after the effective-guest sizing.
- A promotion would change a claim's scientific question.
- The placement instrument cannot list leaders on the pinned engine.

Everything else: Decide → Log → Continue.
