---
analysis_id: 20260927T-v2-11-memory-envelope--deepseek-flash--2026-09-27
run_id: 20260927T142039Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-03/v2-11-memory-envelope."
analyzed_at: 2026-09-27
inputs_digest: febd951548f154ca (run 20260927T142039Z, database container 5 GB and tserver limit 3 GiB); pass 1 is 0e717c5c3a629cd6 (3 GB / 1.5 GiB)
repo_commit: dc76b7804873a10b8eeb62398cf0ecf0d56220aa
status: current
headline: With the pair run in the opposite order under a raised memory envelope, v2-11's multiple reproduces at 24.4x against pass 1's 27.8x - so the claim registers as run_replicated - while l3_section_sharded fails its load for the third time, now with a different symptom, so the sharded half stays unmeasured.
---

# Analysis — study 03 `v2-11`, memory-bounded pass: the reversal succeeded, `l3` still did not — deepseek-flash

## TL;DR

- **The AB/BA pair now exists.** This pass ran `l2_section_document` **before** `l1_claim_rows`; pass 1 ran
  them the other way round. The 100k-seat map multiple is **24.4x** here against **27.8x** there — a **12%**
  difference, inside the 40% the protocol allowed. The claim's direction survives having its two designs
  swap places, so `v2-11` can now be registered `run_replicated` with `order: abba`.
- **Raising the envelope did not rescue `l3_section_sharded`.** W1 is falsified: it failed again, and its
  failure *changed shape* — from connections terminated mid-`COPY` (SQLSTATE 57P01) to
  `dial tcp 10.89.0.41:5433: connect: connection refused`, i.e. nothing listening when the load began. Three
  passes, three failures, in three different positions, under two memory envelopes.
- **It did rescue `l1_claim_rows`**, which failed its load in pass 2 and completed here. So the memory limit
  was real and worth fixing; it was not the whole story.
- **The claim's second half is still unmeasured**, and now for a reason that has outlived a plausible fix.
  The sharded design's load cannot complete on this host in this harness, so its collapse is neither
  confirmed nor refuted.
- **Correctness is clean where it ran**: `l1` and `l2` both passed 65/65 against Go-computed truth.

## The three passes

| Pass | Order (pair) | Envelope | `l1` 100k map | `l2` 100k map | Pair multiple | `l3` |
|---|---|---|---|---|---|---|
| 2026-09-21 (baseline) | l1 then l2 | 3 GB / 1.5 GiB | 149 | 3.1k | 20.8x | 15.9 ops/s on another question |
| pass 1 `…234227Z` | **l1 then l2** | 3 GB / 1.5 GiB | 41 | 1,138 | **27.8x** | load failed (57P01) |
| pass 2 `…111110Z` | l3 first | 3 GB / 1.5 GiB | load failed | 351 | — | load failed (57P01) |
| pass 3 `…142039Z` | **l2 then l1** | **5 GB / 3 GiB** | 34 | 830 | **24.4x** | load failed (connection refused) |

Two designs, two orders, two envelopes: the multiple is 27.8x and 24.4x. The absolutes moved much further
(`l1` 149 → 41 → 34; `l2` 3.1k → 1,138 → 830), which is the same pattern studies 02 and 04 showed.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| W1 | the raised envelope removes the load failures | **falsified** | `l3` failed again, with a new symptom |
| W2 | the reversal reproduces the direction (≥4x) | **confirmed** | 24.4x |
| W3 | the two passes agree within ±40% | **confirmed** | 27.8x against 24.4x, −12% |
| W4 | the correctness gate stays green where it runs | **confirmed** | `l1` and `l2` 65/65 |

## What this changes

1. **`v2-11` becomes `run_replicated`** — three whole runs, the pair measured in both orders, three trials
   per cell — with `order: abba` and both replications as anchors. Its strength was the only thing the
   earlier passes could not supply, and the lint was right to refuse it then.
2. **`l3`'s failure is now a named, reproducible blocker rather than an incident.** Two symptoms, three
   positions, two envelopes: the honest statement is that the sharded collapse *cannot be tested here* with
   this harness, which is different from having tested it and found nothing.
3. **A memory limit is a confound to be declared, not just a limit.** This pass's absolutes are not
   comparable with the 3 GB runs; only the ratio is, and the claim says so.
4. **`record_host_state` earned its place immediately**: it captured `db_budget: cpus=2 memory=5g`,
   `yb_tserver_memory_bytes: 3221225472` and the host's 15.7 GB beside the results, which is the attribution
   that was missing when two cells died in pass 2.

## Where this analysis is weak

1. **`l3` is unexplained.** A startup race, a schema-size problem on that design's load, and a host
   limitation are all consistent with "connection refused" on the first cell; nothing here separates them,
   and the next diagnostic step is to run `l3` alone against a topology left up.
2. **The two compared passes differ in envelope as well as order**, so this is not a clean AB/BA: the 12%
   agreement could be hiding a memory effect and an order effect that happen to cancel. A fourth pass at
   5 GB in pass 1's order would separate them.
3. **Within-cell spreads are 207–340%** on the questions quoted, so the multiples rest on medians of three
   noisy trials; the agreement between passes is more convincing than either pass alone.
4. **One engine, one topology, one scale, one host.** The engine-dependence that makes the claim interesting
   is still untested, and the 3-node topology remains unmeasured.

## Next experiments

1. **Run `l3` alone** against a topology that is already up, with the host state recorded during the load,
   to turn "connection refused" into a cause.
2. **A fourth pass at 5 GB in pass 1's order**, which would make the comparison a clean AB/BA rather than an
   order-and-envelope change together.
3. **The PostgreSQL half**, which is what would make `v2-11` an engine-dependent finding rather than a
   YugabyteDB one.
