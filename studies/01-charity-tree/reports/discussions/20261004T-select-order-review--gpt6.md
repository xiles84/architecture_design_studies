---
discussion_id: 20261004T-select-order-review--gpt6
analyst: GPT-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop, selected effort unknown, HIGH session"
analyzed_at: 2026-10-04
repo_commit: aa900e3fd93035f700fb5938be29765e11ca3bca
inputs_digest: "access 781fcb68ac4907615da0301344d5225445ebec91a8573687ae0ed6eb03c4a14e; aggregate 2df8f4499412a81f23a7f85d7952f25262e70abad40b69ccddac2c30300bde59; existence e6b469b428b00501e0dcd64dd277aa72d7121e5a2e248c9dea1473f9ca246e0f"
references_analysis: 20261004T-select-order-correction--deepseek-flash--2026-10-04
---

# SELECT order correction — independent interpretation qualification

## TL;DR

The corrected schedule and exact recorded inputs can be integrated. The cited
analysis's YB level-10 sign statement is incorrect, and its level-10,000 headline
conflicts with its own uncertainty caveat. Keep those effects unresolved; the
book receives no new SELECT timing claims. Original signed analyses remain intact.

## Findings

Recomputed `DISTINCT ON / ROW_NUMBER` from the corrected YB run JSONs:
level 10 is **0.987820725, 1.106794988, 1.003472921**, which is mixed-sign,
not `>1` in all three runs. At 10,000 it is **0.697483137, 0.996803398,
0.788633396**: all medians are numerically below one, but the second is only
0.32% below equality. That supports a literal sign count, not a resolved
replicated practical speed advantage. PostgreSQL levels 10 and 100 are mixed.

Actual-first position repeats have last/first ratios 0.710428920 (PG aggregate)
and 1.743230620 (YB aggregate). Their roughly 29% and 74% drift is a material
limitation. The first measurement is an interleaved-round median and the last is
a consecutive repeat, so this also compares sampling positions/schedules rather
than estimating a pure time trend. A single repeated cell does not bound every
pair's uncertainty. Access's stable sign bracket is retained as an observation
of these fixtures; it is not an estimate of every workload's crossover.

## Where the measurement is weak

Only three whole runs, five samples, one sequential reader, one shared host,
single-node engines and no independent hardware. No three-node/equal-total
comparison, placement arm or workload throughput was measured. The recorded
order and correctness controls remove the identified scheduling defect without
establishing those missing comparisons. The unfinished join/CTE family remains
a gate for the separately queued HIGH synthesis.

## Evidence

`docs/book-release-20261004/CORRECTION-INPUTS.md` is the independent exact-input
manifest; `CORRECTION-DETAILS.txt` records ratios/control counts. They verify 18
runs, 516 valid cells and 2,580 samples, three distinct orders per family/engine,
the recorded source `ca2a38fb` and unchanged annotated producing tags. Earlier
manifests' empty tag column is an exporter defect, corrected in the release
source and regression-tested; raw inputs already contain the tags.
