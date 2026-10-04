# Targeted SELECT order correction

Goal: close the HIGH review's measurement-order and provenance gaps for access paths, aggregate/top-N and existence; preserve their scientific questions, SQL contracts and immutable results.

Current state: all three families are integrated and correctness-gated. Their execution order is fixed across runs, with consecutive samples inside each cell. Original evidence remains useful as explicitly order-confounded observations. The signed HIGH companion lists material inference corrections.

Already decided:
- Keep all original results, reports, analyses and producing tags intact. Use fresh run ids and a new signed correction analysis; no re-labelling of old runs as randomized.
- Retain the same seeds, fixtures, predicates, five pagination/parent levels, three fanout levels, page size and N=1. If a question needs an expanded sweep, report it separately instead of changing these pairs.
- Add seeded randomized spelling order per independent run and interleaved sample rounds within each logically matched phase, recording seed and actual order. Reloaded existence fixtures may form separate phases; the physical-index arm remains separate and labelled.
- Keep every timed query drained and independently verified; preserve all current fixture and wrong-answer gates. Include a first/last repeat of the actual first cell in each tested engine/family. Never repeat a hard-coded cell while calling it the randomized first.
- Measure >=3 whole runs per engine/family and >=3 trials, serially under ads-run-lock. Keep saved budgets, actual guest/node/client reserves and per-node/client throttling explicit. Do not change the Podman VM profile to make a comparison fit.
- Recorded guest resources supersede stale 2 GiB guest claims. Calculate cluster per-node/equal-total/client budgets against live resources; if infeasible or a supported topology/placement measurement remains absent, record the gap. The existing single-node evidence establishes no cluster result.
- Generate reports from an explicit primary-input manifest with per-file and corpus digests; expose run ids on ratio rows and link signed analyses. Never accidentally pool pilots, failed cells, topology/resource settings or old fixed-order runs.
- Judge repeated paired ratios, within/across-run variability and position drift together. Mixed signs, missing ordering evidence or effects smaller than uncertainty stay unresolved. No universal equality, no client-capacity inference from zero CFS throttling, no unmeasured RPC attribution.
- Page depth 100 on PG is unresolved; provisional stable-sign observations bracket 10-1,000. YB DISTINCT ON at 10,000 beats ROW_NUMBER in only 2/3 old runs. NOT IN with a NULL has FALSE for matches and UNKNOWN for nonmatches, neither TRUE in WHERE.
- Integrate current main in the correction worktree and preserve concurrent join/CTE work. No writes to another task worktree. Run meaningful containerized harness/combined-source checks.

Acceptance: source records randomized/interleaved execution and retains correctness gates; new serial producing-tagged matrices have full input manifests and uncertainty-qualified signed corrections; all new checks have reports; HIGH validates before numeric promotion and local-main integration. No book rewrite or frozen evidence-package edit in this task.

Escalate only a change to question/answer/oracle/predictions or material downstream architecture. Ordinary scheduling/order representation is local and reversible.

NEXT MODEL: LOW
