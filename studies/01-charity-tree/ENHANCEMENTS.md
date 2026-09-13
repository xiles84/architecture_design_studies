# Study 01 v3 — follow-up measurement protocol

Requested from `20260912-study01--gpt-6--2026-09-12`, “What I would measure next”.
Before tag: `study-01/v2-before-enhancements`. Implementation and results receive new
annotated v3 milestones and immutable run tags. Existing result sets stay unchanged.

| Follow-up | Implemented experiment / acceptance condition |
|---|---|
| D2/D3 repeat | At least five fresh loads per design, AB/BA order by trial; all 12 reads, geometric score and spread; q08/q12 largest/smallest charity; plans before and after reads under identical preparation |
| Mechanisms | D11 copies the key with D2 indexes/reads and D3's extra FK; D12 adds the date index; D13 adds recency rewrites; D17 adds only the sum rewrite; D14 adds a plain charity index; D15 replaces that plain index with INCLUDE(amount_cents); D3 adds the ranking rewrites. Each subsequent step isolates one mechanism |
| Growth | Medium plus longer histories, repeated deterministic insert/correct/delete cycles with Go truth checks after each phase; relation bytes, cgroup memory and disk-read counters. A separately labelled constrained-memory regime must demonstrate physical relation bytes exceeding the database memory limit; do not call one million rows out-of-memory by assumption |
| Busiest parent | Fixed reader count, controlled offered write rates and donor/charity concentration; D4, D5 lock-first and D16 sum-only trigger. Identical insert-only mix for the application comparison; correction/delete cost separately |
| Exceptions | YB one-node D3/D8 inserts; three-node erasure with a larger finite donor pool and actual elapsed windows; D9 negative control against D10 on both YB topologies, initial verification and post-write audits |
| Deployment | Explicit equal-total-budget YB one large node / three smaller nodes; scheduled open-loop arrivals with queueing latency, service time, rejected arrivals and errors. Local node-stop recovery is a separate diagnostic; physical-host network/failure validation needs independent machines |

All timing cells verify independently generated Go answers first. The new experiment
runner uses the existing study 01 loader, catalogue and driver directly; it does not
migrate the legacy study to `platform/`. New command paths preserve legacy result formats.
The experiment report groups all repeated cells without overwriting by design/topology.
Each run pins the benchmark image ID, commits, resource condition, preparation and order.

Generated reports contain measurements only. New final analyses are concise and link
to signed `reports/discussions/` documents for design mechanisms and analyst exchanges.
The previous analyses remain intact, attributable to the inputs they actually reviewed.

## External deployment validation

Real network separation cannot be represented by three containers on this laptop.
Use three documented hosts (same CPU generation and disk class), a fourth client host,
measured RTT/loss, pinned engine images and explicit budgets. Repeat the exact catalogue
and arrival schedule at RF=1/RF=3; record acknowledged writes, unavailable intervals,
and reconciliation after stopping one database host. Run from a new committed run tag.
Do not infer host-failure durability from stopping a local process or container.
