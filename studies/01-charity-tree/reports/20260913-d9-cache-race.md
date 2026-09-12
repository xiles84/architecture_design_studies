# Experiment D — reads and writes running concurrently

| | |
|---|---|
| Run id | `20260913-d9-cache-race` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `small` — 5000 people, 108081 donations |

## Why this experiment exists

Every other measurement in this study runs reads alone, then writes alone. That
gives clean attribution, but it cannot see the costs that only exist when the two
overlap: MVCC bloat accumulating *while* readers scan, contention on rows writers
keep updating, autovacuum and compaction stealing CPU, and buffer-cache
competition.

All of those penalise designs that buy read speed with redundancy — which is most
of the designs here. **Isolated measurement therefore flatters exactly the designs
under scrutiny**, and this is where that bias gets checked instead of assumed away.

### The workload

A donor portal plus a dashboard. Weights sum to 100.

| Read query | Weight | | Write op | Weight |
|---|---:|---|---|---:|
| A person's 20 most recent donations | 30 | | insert | 85 |
| Donation by id (point lookup) | 20 | | update | 7 |
| Charity activity feed (last 50, with names) | 15 | | update_person | 5 |
| How many donations a person made | 10 | | delete | 3 |
| Information of the last donation (one charity) | 8 | | |
| First and last donation of a person | 5 | | |
| Total donated (one charity) | 5 | | |
| Who donates the most | 4 | | |
| Total donated (global) | 3 | | |

Splits are **readers:writers**. The first is all-readers and is the baseline —
same query mix, same data, same process, so the only variable across splits is
the presence of writers. The dataset is reloaded before every split, so no split
inherits the bloat the previous one created.

## PostgreSQL, 1 node

| Design | Split | Reads/s | **Read throughput kept** | Writes/s | Read p99 (ms) | Aggregates still correct |
|---|---|---:|---:|---:|---:|---|
| D9 hybrid | 4:4 | 2.8k | **100%** | 907 | 2.20 | **NO — 3 rows wrong** |
| D9 hybrid | 2:6 | 1.2k | **43%** | 1.7k | 1.12 | **NO — 5 rows wrong** |
| D9 hybrid | 4:4 | 2.7k | **96%** | 907 | 2.25 | **NO — 3 rows wrong** |
| D9 hybrid | 2:6 | 1.5k | **54%** | 2.0k | 1.08 | **NO — 5 rows wrong** |

## Which questions suffer most under write load

Percentage of the all-readers throughput retained at the heaviest write split.
A design can hold its aggregate throughput while one specific question collapses,
and the blended number would never show it.

### PostgreSQL, 1 node

| Question | D9 hybrid | 
|---|---:|
| Information of the last donation (one charity) | 55% | 
| Who donates the most | 55% | 
| First and last donation of a person | 53% | 
| Total donated (global) | 53% | 
| Total donated (one charity) | 52% | 
| A person's 20 most recent donations | 54% | 
| Donation by id (point lookup) | 54% | 
| How many donations a person made | 53% | 
| Charity activity feed (last 50, with names) | 55% | 

> Readers and writers share the same eight host cores here, so some of the loss
> below is simply CPU being taken by the writers rather than interference in the
> database. The comparison that survives that caveat is **between designs at the
> same split** — they all lose the same CPU, so a design that keeps less of its
> throughput than another is losing it to something other than scheduling.

## Conclusions and analysis provenance

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so that several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, that disagreement is itself a finding and is
left visible rather than resolved by editing one of them.

| | |
|---|---|
| Run id | `20260913-d9-cache-race` |
| Result files | 1 |
| **Inputs digest** | `70f87f5d5667bff8` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `70f87f5d5667bff8` ✅ | current | Rollups win aggregate questions by 60-130x but cost 2.4x on inserts and collapse under many writers unless parents are many and small; plain indexes answer everything else; embedding never wins outright on PostgreSQL and a naive cache trigger silently corrupts data. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`70f87f5d5667bff8`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
