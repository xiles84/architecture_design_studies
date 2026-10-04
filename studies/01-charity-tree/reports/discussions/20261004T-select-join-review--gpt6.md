---
discussion_id: 20261004T-select-join-review--gpt6
analysis_id: 20261004T-select-join-cte--deepseek-flash--2026-10-04
run_id: select-join-pg1-20261004
environment: host-zenbook-ux5406sa
analyst: GPT-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop, HIGH; selected effort unknown"
analyzed_at: 2026-10-04
inputs_digest: 070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac
inputs: "pg2@112577024534a9e91932f78d3ba8b9e530e0caea0aeefdf851737e578be7bab6; pg3@05b2d7a99e38c77c84baf413a769f531179682921a2419bf73da89153b2d00ac; yb1@f650acde1ebb5541080360bca3e467863611e00322d3f580cefdd36ba90e0281; yb2@25e5f7bf01de5149bed501f089eebe900d7b462c2927eedaee302436b664d198; yb3@2d8394b9cc2a9a0f649367486f9d994fcfc1538e3e7d151bc75a1e8febe82e46"
repo_commit: 7e148174b5211225f6b4f243fc02e0de417a0586
responds_to: 20261004T-select-join-cte--deepseek-flash--2026-10-04
---

# HIGH qualification — join/CTE evidence requires correction

[Original LOW analysis](../analyses/20261004T-select-join-cte--deepseek-flash--2026-10-04.md)
· [Replicated measurements](../20261004T-select-join-replicated.md)
· [Independent checks](../../../../docs/select-query-templates/reviews/20261004-join-high/CHECKS.md)
· [Review and corrective route](../../../../docs/select-query-templates/reviews/20261004-join-high/REVIEW.md)

## TL;DR

The raw single-node evidence has correct input hashes, medians, randomized
orders and firing controls. Completion and book promotion require correction:
the resource rejection contradicts saved guest resources, mandatory domain,
blended and distributed arms are absent, and report regeneration changes the
advertised input hash. Preserve the six original runs as their actual synthetic
parent/fact observations.

## What the matched policy pairs show

For PostgreSQL, joined-output MATERIALIZED/default at rank threshold 20 is
73.1056, 92.8513 and 85.3413 across the three runs. At threshold 20,000 it is
6.7868, 7.3624 and 7.1675. These large local observations are consistent with
the captured full-join materialization boundary. They do not establish a
portable numerical rule or repair the missing workload/domain coverage.

For YugabyteDB at threshold 180,000 the same paired ratio is **1.0172,
0.8844, 0.9442**: mixed sign. The first run cannot supply the range for a
replicated claim. The small broad-predicate effect remains unresolved.

P4 compares MATERIALIZED with NOT MATERIALIZED, rather than default with
NOT MATERIALIZED: default at one reference already folds. At PostgreSQL's
100% consumer endpoint, NOT MATERIALIZED/MATERIALIZED is **0.9080–0.9687**
at one reference, **1.6957–2.1235** at two, and **4.0048–4.3493** at sixteen.
The first endpoint is near equality; do not call it a resolved reference-count
crossover without considering variability. At the 1% consumer endpoint,
sixteen-reference ratios are **0.6950–0.7868**. These two consumer endpoints
show a changed direction, but do not constitute a five-level consumer-selection
sweep or establish its crossing location.

The first-to-last position control is 0.6056–1.3470 across the six runs.
PostgreSQL run 3 has 88/91 cells with over 20% sample spread. The original
analysis's phrase that a 1.52x effect is "inside 0.80 drift" is not a valid
arithmetic or significance test. Compare paired observations with sample and
run spread, retaining unresolved outcomes. P12 means the original broad-predicate
convergence prediction; it still needs its own verdict. A matching shape
signature is evidence about nodes, not proof of complete plan or timing equality.

## Resource and provenance evidence

All runs save MemTotal 16,109,780 kB and online CPUs 0-7. The machine inspect
configuration says 2 GiB/4 CPUs, but effective WSL runtime is ~15.36 GiB/8 CPUs.
The latter controls sizing. This contradicts the claimed single-node memory
failure and leaves the skipped topology/placement subset unsupported.

The unchanged pg1 input hash is `070df793…4f2cfac`. Re-rendering with its
original producing image advertises `7c879564…824b6c`, because the renderer
appends an extra newline before hashing. The original measurements remain
intact; fix derived reporting and test exact digest stability. Plans use
ANALYZE/BUFFERS on both engines, omitting the required PG VERBOSE and YB DIST
evidence. An RPC "not parsed" caveat cannot replace counters never requested.

## Limits and next action

No new database runs were made by this review. Current timings describe one
synthetic 200,000-fact fixture, one host, one sequential reader and zero writers.
The original signed analysis remains unchanged. The task receives a targeted
LOW handoff for the approved domain, calibrated topology/placement and matched
blended phases, plan/provenance repairs and a new signed interpretation. No
universal materialization or optimizer conclusion is accepted from these inputs.
