---
analysis_id: <run_id>--<analyst>--<YYYY-MM-DD>
run_id: <run_id, e.g. 20260912-small>
environment: <environment id, e.g. host-zenbook-ux5406sa>
analyst: <short stable id, e.g. claude-opus-5, gpt-x, jane-doe>
analyst_kind: <ai | human>
analyst_version: <full identification, e.g. "Claude Opus 5 (claude-opus-5) via Claude Code">
analyzed_at: <YYYY-MM-DD>
inputs_digest: <paste the digest from the generated report — this is what identifies the data>
supersedes: <analysis_id you are replacing, or leave blank>
status: <current | superseded>
headline: <one sentence, shown in the report's index table>
---

# Analysis — <run id> — <analyst>

> Delete this quote block when you fill the template in.
>
> **Before writing:** open the generated report, find the `inputs_digest`, and check its
> analyses table. If an analysis already exists for that digest, read it. Write a new one
> to **disagree, extend, or bring a different perspective** — not to restate. Never edit
> someone else's analysis; write your own and reference theirs by `analysis_id`.

## What I was looking at

State the run, the scale, the topologies, and anything you deliberately ignored. If cells
failed or are missing, say so here rather than quietly working around them.

## Conclusions

The interpretation. Lead with what a reader should *do differently*, not with what the
numbers were — the numbers are already in the generated report and do not need repeating.

Cite specific measurements when you make a claim, so a later analyst can check you:

> D2 → D3 removed a join from six queries but moved `q08` from 432 to 7 200 ops/s only
> because the covering index turned it into an index-only scan — see the plan in
> `plans/d3_flattened_fk.txt`, which shows no heap fetches.

## Where I think the measurement is weak

Every analysis should include this. Name the places where you do not trust the numbers,
or where the environment's limits (shared cores, no real network, single trial) could
plausibly explain what you observed.

## What I would measure next

Concrete, runnable next experiments — ideally phrased so someone can turn them into a
`run-study.sh` invocation or a new design directory.

## Disagreements with other analyses

If you have read another analysis of this digest and reached a different conclusion, say
so explicitly and say why. Disagreement between analysts is a finding, not a problem to
be tidied away.
