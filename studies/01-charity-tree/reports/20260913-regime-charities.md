# Regime comparison — 10 charities vs 1000 charities

| | 10 charities | 1000 charities |
|---|---|---|
| Run id | `20260913-history-unbounded` | `20260913-top-1000charities` |
| Environment | `host-zenbook-ux5406sa` | `host-zenbook-ux5406sa` |
| Scale | `small` | `small` |
| Profile | `unbounded` | `unbounded-charities1000` |
| Charities / people / donations | 10 / 5000 / 108081 | 1000 / 5000 / 108568 |
| Max donations per person | 496 | 496 |

> Donation counts differ slightly between the two datasets (the capped generator
> stops at the first donor that reaches the uncapped total, and extra charities
> consume extra random draws). The difference is a few rows in ~100 000, far
> below any effect reported here.

## PostgreSQL, 1 node

### How much the regime changed each design

Ratios are **1000 charities ÷ 10 charities**: above 1.0x the alternative regime helped this design.

| Design | Read score | insert donation | correct a donation amount | remove one donation | donor edits their profile | erase a donor and all their donations | 
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | **2.0x** | 0.96x | 0.89x | 0.90x | 0.92x | 0.87x | 
| D4 rollup/trg | 0.86x | **1.9x** | **1.7x** | **1.5x** | 0.91x | **1.5x** | 
| D6 embedded | **6.9x** | 0.92x | **1.9x** | **0.67x** | 0.99x | 0.95x | 
| D9 hybrid | **2.0x** | 0.96x | 0.95x | 0.98x | 0.91x | 0.96x | 

### Which design is fastest, per question, in each regime

A regime that changes the **winner** is the result that belongs in a
conclusion. Rows where the winner changed are marked **→**.

| Question | Best in 10 charities | Best in 1000 charities |
|---|---|---|
| Information of the last donation (global) | D9 hybrid (37k) | D3 flat+FK (37k) **→** |
| Information of the last donation (one charity) | D9 hybrid (36k) | D9 hybrid (34k) |
| Who donates the most | D4 rollup/trg (37k) | D4 rollup/trg (36k) |
| Top-10 donor leaderboard | D4 rollup/trg (34k) | D4 rollup/trg (34k) |
| Who donated last | D9 hybrid (36k) | D4 rollup/trg (33k) **→** |
| First and last donation of a person | D4 rollup/trg (39k) | D4 rollup/trg (37k) |
| Total donated (global) | D4 rollup/trg (37k) | D4 rollup/trg (19k) |
| Total donated (one charity) | D4 rollup/trg (41k) | D4 rollup/trg (37k) |
| A person's 20 most recent donations | D3 flat+FK (32k) | D3 flat+FK (30k) |
| Donation by id (point lookup) | D3 flat+FK (36k) | D4 rollup/trg (36k) **→** |
| How many donations a person made | D4 rollup/trg (38k) | D4 rollup/trg (37k) |
| Charity activity feed (last 50, with names) | D3 flat+FK (14k) | D9 hybrid (16k) **→** |
| write: insert donation | D3 flat+FK (6.2k) | D3 flat+FK (6.0k) |
| write: correct a donation amount | D3 flat+FK (6.5k) | D3 flat+FK (5.9k) |
| write: remove one donation | D3 flat+FK (7.4k) | D3 flat+FK (6.7k) |
| write: donor edits their profile | D9 hybrid (7.2k) | D9 hybrid (6.6k) |
| write: erase a donor and all their donations | D6 embedded (7.5k) | D6 embedded (7.1k) |

### Per-question change (1000 charities ÷ 10 charities)

| Question | D3 flat+FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|
| Information of the last donation (global) | 1.12x | **0.38x** | 1.02x | 0.94x | 
| Information of the last donation (one charity) | 0.91x | 0.97x | **56.2x** | 0.96x | 
| Who donates the most | **19.8x** | 0.97x | **42.6x** | **22.3x** | 
| Top-10 donor leaderboard | **20.8x** | 1.00x | **46.4x** | **20.3x** | 
| Who donated last | 0.84x | 0.97x | **37.5x** | 0.86x | 
| First and last donation of a person | 0.97x | 0.94x | 0.85x | 0.96x | 
| Total donated (global) | 1.07x | **0.51x** | 1.12x | 1.06x | 
| Total donated (one charity) | **12.7x** | 0.90x | **50.6x** | **11.9x** | 
| A person's 20 most recent donations | 0.92x | 0.87x | 0.88x | 0.89x | 
| Donation by id (point lookup) | 0.99x | 1.12x | 0.97x | 0.98x | 
| How many donations a person made | 0.90x | 0.98x | 0.99x | 0.93x | 
| Charity activity feed (last 50, with names) | 1.15x | 1.21x | **68.8x** | 1.23x | 

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
| Run id | `20260913-top-1000charities` |
| Result files | 4 |
| **Inputs digest** | `e1f7e75cbbea44dc` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `e1f7e75cbbea44dc` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `e1f7e75cbbea44dc` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`e1f7e75cbbea44dc`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
