# Regime comparison — small vs tiny

| | small | tiny |
|---|---|---|
| Run id | `20260913-size-small` | `20260913-size-tiny` |
| Environment | `host-zenbook-ux5406sa` | `host-zenbook-ux5406sa` |
| Scale | `small` | `tiny` |
| Profile | `unbounded` | `unbounded` |
| Charities / people / donations | 10 / 5000 / 108081 | 10 / 500 / 8999 |
| Max donations per person | 496 | 462 |

> Donation counts differ slightly between the two datasets (the capped generator
> stops at the first donor that reaches the uncapped total, and extra charities
> consume extra random draws). The difference is a few rows in ~100 000, far
> below any effect reported here.

## PostgreSQL, 1 node

### How much the regime changed each design

Ratios are **tiny ÷ small**: above 1.0x the alternative regime helped this design.

| Design | Read score | insert donation | correct a donation amount | remove one donation | donor edits their profile | erase a donor and all their donations | 
|---|---:|---:|---:|---:|---:|---:|
| D1 minimal | **10.6x** | 1.05x | 1.05x | 1.04x | 1.02x | **8.4x** | 
| D2 indexed | **2.6x** | 0.95x | 1.03x | 1.06x | 0.98x | 1.01x | 

### Which design is fastest, per question, in each regime

A regime that changes the **winner** is the result that belongs in a
conclusion. Rows where the winner changed are marked **→**.

| Question | Best in small | Best in tiny |
|---|---|---|
| Information of the last donation (global) | D2 indexed (34k) | D2 indexed (35k) |
| Information of the last donation (one charity) | D2 indexed (18k) | D2 indexed (24k) |
| Who donates the most | D2 indexed (288) | D2 indexed (4.2k) |
| Top-10 donor leaderboard | D2 indexed (291) | D2 indexed (3.7k) |
| Who donated last | D2 indexed (19k) | D2 indexed (24k) |
| First and last donation of a person | D2 indexed (32k) | D2 indexed (33k) |
| Total donated (global) | D2 indexed (352) | D2 indexed (4.5k) |
| Total donated (one charity) | D2 indexed (356) | D2 indexed (5.8k) |
| A person's 20 most recent donations | D2 indexed (31k) | D2 indexed (20k) |
| Donation by id (point lookup) | D1 minimal (39k) | D1 minimal (37k) |
| How many donations a person made | D2 indexed (34k) | D2 indexed (34k) |
| Charity activity feed (last 50, with names) | D2 indexed (1.1k) | D1 minimal (2.7k) **→** |
| write: insert donation | D2 indexed (6.5k) | D1 minimal (6.5k) **→** |
| write: correct a donation amount | D1 minimal (6.5k) | D1 minimal (6.9k) |
| write: remove one donation | D2 indexed (7.0k) | D2 indexed (7.4k) |
| write: donor edits their profile | D1 minimal (6.6k) | D1 minimal (6.8k) |
| write: erase a donor and all their donations | D2 indexed (6.2k) | D2 indexed (6.3k) |

### Per-question change (tiny ÷ small)

| Question | D1 minimal | D2 indexed | 
|---|---:|---:|
| Information of the last donation (global) | **11.9x** | 1.04x | 
| Information of the last donation (one charity) | **13.7x** | 1.32x | 
| Who donates the most | **13.8x** | **14.5x** | 
| Top-10 donor leaderboard | **14.0x** | **12.6x** | 
| Who donated last | **13.4x** | 1.28x | 
| First and last donation of a person | **12.7x** | 1.03x | 
| Total donated (global) | **12.9x** | **12.9x** | 
| Total donated (one charity) | **13.6x** | **16.2x** | 
| A person's 20 most recent donations | **12.6x** | **0.63x** | 
| Donation by id (point lookup) | 0.95x | 1.00x | 
| How many donations a person made | **12.2x** | 0.99x | 
| Charity activity feed (last 50, with names) | **13.9x** | **1.9x** | 

> Both runs are single-trial surveys unless their manifests say otherwise, so the
> ratios carry the same error bar as the generated reports they came from (the
> D4/D5 control put that at up to ~1.5x). Treat a change under 1.5x as noise and a
> change in winner between two designs within 1.5x of each other as a tie.

## Conclusions and analysis provenance

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so that several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, that disagreement is itself a finding and is
left visible rather than resolved by editing one of them.

| | |
|---|---|
| Run id | `20260913-size-tiny` |
| Result files | 2 |
| **Inputs digest** | `8c20043cd30c3848` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `8c20043cd30c3848` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `8c20043cd30c3848` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`8c20043cd30c3848`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
