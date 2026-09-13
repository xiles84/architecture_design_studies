# Pipeline checks — not study results

Runs made on the `tiny` dataset, with 1 s measurements, while the study 02 harness was being
built. They exist to prove the pipeline works on both engines, and they found real harness
problems (see LESSONS_LEARNED.md, "Building study 02"). **Do not draw conclusions about the
designs from them**, and do not analyse them.

They are kept, not deleted, because LESSONS_LEARNED and CONTEXT cite evidence from them:

| Check | What it was | Produced by | Notable |
|---|---|---|---|
| `devcheck-1` | PostgreSQL, P2 / C1 / H1 | uncommitted code before the first study 02 commit | C1 (negative control) overbooked every race event |
| `devcheck-2` | PostgreSQL all 14 designs; YugabyteDB 1-node P1–C3 | same, plus sell-out-time metric | stopped by hand: C2 on YugabyteDB hung (retries not bounded by the race deadline); C5 under-booked in churn |
| `devcheck-3` | YugabyteDB 1-node C5, R1–R3, H0, H1, C2 (small tiers) | same, plus bounded retries, tier budget, CPU probes | YugabyteDB container throttled in 40–65% of CFS periods; C2 cell lost YSQL ("shutting down"); holds degenerate at a 250 ms TTL |

Each directory holds the result JSON, plans, per-cell console logs, the runner's console log
(`console.log`) and the report generated at the time (`report.md`). The code that produced
them was never committed as such; the first commit of study 02 contains it plus the fixes
these checks prompted.
