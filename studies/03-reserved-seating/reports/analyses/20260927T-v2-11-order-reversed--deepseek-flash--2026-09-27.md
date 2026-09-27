---
analysis_id: 20260927T-v2-11-order-reversed--deepseek-flash--2026-09-27
run_id: 20260927T111110Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-03/v2-11-order-reversed."
analyzed_at: 2026-09-27
inputs_digest: b356727900544927 (run 20260927T111110Z); compared with 0e717c5c3a629cd6 (run 20260926T234227Z)
repo_commit: 2d623f3495f3dba312a963e979a8cf13e18b16ca
status: current
headline: The order-reversed replication produced one usable cell out of three - two loads were killed by memory pressure inside the Yugabyte container - so v2-11 stays single_run_directional, and the blocker is host memory rather than design order.
---

# Analysis — study 03 `v2-11`, order-reversed attempt: the machine refused the experiment — deepseek-flash

## TL;DR

- **The reversal did not happen.** Two of the three cells — `l3_section_sharded` (run first) and
  `l1_claim_rows` (run last) — died during their **load**, so the pair the claim rests on could not be
  compared in the opposite order. V1 and V2 are **not evaluable**.
- **`l3`'s earlier failure was not positional.** V3 is falsified: it failed again when run first with a fresh
  host, and the mechanism is now visible — its server processes were **killed with SIGKILL**
  ("server process (PID 721) was terminated by signal 9: Killed", 122 such lines in one cell's database log).
  That is the out-of-memory killer, not a design property.
- **The one cell that ran moved by 58–69%** from its counterpart in the first replication
  (`l2_section_document`: 100k map 1,138 → 351 ops/s; seats-available 1,145 → 484) with within-cell spreads
  of **289%, 710% and 318%**. On this evidence a single cell is not a measurement.
- **`v2-11` stays `single_run_directional`**, correctly and now for two independent reasons: no order control,
  and the reversal that would have provided it was not produced.
- **The next attempt must change the resource envelope, not the design order.** Repeating this run as it
  stands would waste another thirteen minutes and lose cells the same way.

## What happened

| Cell | Position in run 2 | Outcome |
|---|---|---|
| `l3_section_sharded` | first | load failed 3/3, connection terminated (SQLSTATE 57P01) |
| `l2_section_document` | second | completed; 100k map read 351 ops/s |
| `l1_claim_rows` | third | load failed 3/3, same error |

| Cell | Position in run 1 | Outcome |
|---|---|---|
| `l1_claim_rows` | first | completed; 100k map read 41 ops/s |
| `l2_section_document` | second | completed; 100k map read 1,138 ops/s |
| `l3_section_sharded` | third | load failed 3/3, same error |

Read together: the cell that failed twice occupied two different positions, and in run 2 a cell that had
succeeded in run 1 failed. **Position does not predict it; memory pressure does.**

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| V1 | the reversal reproduces the direction | **not evaluable** | `l1` produced no measurement |
| V2 | the two orders' multiples agree within ±40% | **not evaluable** | no second multiple exists |
| V3 | `l3`'s failure was positional | **falsified** | it failed when run first; the cause is SIGKILL |
| V4 | the correctness gate stays green | **partial** | `l2` passed; the other two never reached the gate |

## What this changes

1. **`v2-11` gains two limits**: that an order-reversed attempt was made and produced no comparable pair
   because two loads were killed, and that the unmeasured half of the claim is now blocked by host memory
   rather than by cost or by any property of the designs.
2. **A resource limit is now the binding constraint on study 03 replication.** The host reports 15 GB total
   and 14 GB available between runs; during a load the database's processes are killed outright. Any
   three-design `yb-single` run on this machine carries that risk, and a lost cell must keep being reported
   as *never reaching the gate* rather than as a failed one.
3. **My own count of what a replication is worth.** Six cells have been attempted across the two passes;
   four produced data and two did not. None of the four can be combined into a claim-level statement,
   because the two passes do not agree and one pass has half a pair missing.

## Where this analysis is weak

1. **The diagnosis is a symptom reading, not a measurement.** SIGKILL plus a container memory limit is the
   most economical explanation, but no memory samples were taken during the load; the container's limit is
   inferred from the topology scripts, not observed at the moment of death.
2. **The 14 GB "available" figure was taken after the run**, when the containers had been torn down. It does
   not contradict the diagnosis, but it does not support it either.
3. **One surviving cell cannot separate the two sessions.** The 58–69% movement between run 1 and run 2 for
   `l2` is consistent with the session-level variability this corpus has measured three times, and with the
   smaller load that memory pressure would have produced; this tranche cannot distinguish them.
4. **The claim's engine-dependence is still untested**, as is the sharded collapse, and now neither can be
   tested here until the memory question is settled.

## Next experiments

1. **Bound the database's memory explicitly** — give the `yb-single` container a documented limit and a
   smaller scale, then repeat the reversal. A cell that dies should die for a reason the run can name.
2. **Record host and container memory beside every run**, which turns "we think it was OOM" into a number.
3. **Only then** repeat the PostgreSQL half of the comparison, which is what makes `v2-11` interesting.
