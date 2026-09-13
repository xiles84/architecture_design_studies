# Regime comparison — unbounded vs cap20

| | unbounded | cap20 |
|---|---|---|
| Run id | `20260913-history-unbounded` | `20260913-history-cap20` |
| Environment | `host-zenbook-ux5406sa` | `host-zenbook-ux5406sa` |
| Scale | `small` | `small` |
| Profile | `unbounded` | `cap20` |
| Charities / people / donations | 10 / 5000 / 108081 | 10 / 12015 / 108089 |
| Max donations per person | 496 | 20 |

> Donation counts differ slightly between the two datasets (the capped generator
> stops at the first donor that reaches the uncapped total, and extra charities
> consume extra random draws). The difference is a few rows in ~100 000, far
> below any effect reported here.

## PostgreSQL, 1 node

### How much the regime changed each design

Ratios are **cap20 ÷ unbounded**: above 1.0x the alternative regime helped this design.

| Design | Read score | insert donation | correct a donation amount | remove one donation | donor edits their profile | erase a donor and all their donations | 
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 0.87x | 0.95x | 0.92x | 0.94x | 0.89x | 0.99x | 
| D4 rollup/trg | 0.96x | 0.91x | 0.94x | 0.97x | 0.92x | **1.6x** | 
| D6 embedded | 1.08x | 1.18x | **3.2x** | 1.33x | 0.97x | 0.96x | 
| D9 hybrid | 0.94x | 0.96x | 1.01x | 1.08x | 0.89x | 1.33x | 

### Which design is fastest, per question, in each regime

A regime that changes the **winner** is the result that belongs in a
conclusion. Rows where the winner changed are marked **→**.

| Question | Best in unbounded | Best in cap20 |
|---|---|---|
| Information of the last donation (global) | D9 hybrid (37k) | D9 hybrid (36k) |
| Information of the last donation (one charity) | D9 hybrid (36k) | D9 hybrid (35k) |
| Who donates the most | D4 rollup/trg (37k) | D4 rollup/trg (35k) |
| Top-10 donor leaderboard | D4 rollup/trg (34k) | D4 rollup/trg (31k) |
| Who donated last | D9 hybrid (36k) | D9 hybrid (33k) |
| First and last donation of a person | D4 rollup/trg (39k) | D4 rollup/trg (38k) |
| Total donated (global) | D4 rollup/trg (37k) | D4 rollup/trg (37k) |
| Total donated (one charity) | D4 rollup/trg (41k) | D4 rollup/trg (40k) |
| A person's 20 most recent donations | D3 flat+FK (32k) | D4 rollup/trg (29k) **→** |
| Donation by id (point lookup) | D3 flat+FK (36k) | D9 hybrid (35k) **→** |
| How many donations a person made | D4 rollup/trg (38k) | D4 rollup/trg (37k) |
| Charity activity feed (last 50, with names) | D3 flat+FK (14k) | D3 flat+FK (12k) |
| write: insert donation | D3 flat+FK (6.2k) | D3 flat+FK (5.9k) |
| write: correct a donation amount | D3 flat+FK (6.5k) | D3 flat+FK (6.0k) |
| write: remove one donation | D3 flat+FK (7.4k) | D3 flat+FK (7.0k) |
| write: donor edits their profile | D9 hybrid (7.2k) | D4 rollup/trg (6.5k) **→** |
| write: erase a donor and all their donations | D6 embedded (7.5k) | D6 embedded (7.2k) |

### Per-question change (cap20 ÷ unbounded)

| Question | D3 flat+FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|
| Information of the last donation (global) | 1.08x | 1.00x | **1.6x** | 0.97x | 
| Information of the last donation (one charity) | 0.88x | 0.94x | 1.39x | 0.98x | 
| Who donates the most | 0.92x | 0.93x | 1.04x | 0.89x | 
| Top-10 donor leaderboard | 0.95x | 0.93x | 0.98x | 0.84x | 
| Who donated last | 0.78x | 0.93x | 1.07x | 0.91x | 
| First and last donation of a person | **0.60x** | 0.96x | 1.02x | 0.96x | 
| Total donated (global) | 0.98x | 1.00x | 1.02x | 1.02x | 
| Total donated (one charity) | 1.07x | 0.96x | 0.87x | 1.05x | 
| A person's 20 most recent donations | 0.78x | 0.94x | 1.08x | 0.89x | 
| Donation by id (point lookup) | 0.83x | 1.09x | 1.34x | 0.98x | 
| How many donations a person made | 0.88x | 0.98x | 1.03x | 0.92x | 
| Charity activity feed (last 50, with names) | 0.84x | 0.84x | 0.73x | 0.85x | 

## YugabyteDB, 3 nodes (RF=3)

### How much the regime changed each design

Ratios are **cap20 ÷ unbounded**: above 1.0x the alternative regime helped this design.

| Design | Read score | insert donation | correct a donation amount | remove one donation | donor edits their profile | erase a donor and all their donations | 
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 0.98x | 1.00x | — | — | 0.94x | 1.13x | 
| D6 embedded | 0.87x | 1.19x | — | — | 0.98x | 1.11x | 
| D7 yb-coloc | 0.98x | 0.93x | — | — | 0.95x | 1.17x | 

### Which design is fastest, per question, in each regime

A regime that changes the **winner** is the result that belongs in a
conclusion. Rows where the winner changed are marked **→**.

| Question | Best in unbounded | Best in cap20 |
|---|---|---|
| Information of the last donation (global) | D3 flat+FK (2.7k) | D3 flat+FK (2.4k) |
| Information of the last donation (one charity) | D3 flat+FK (2.4k) | D3 flat+FK (2.3k) |
| Who donates the most | D6 embedded (290) | D6 embedded (222) |
| Top-10 donor leaderboard | D6 embedded (293) | D6 embedded (228) |
| Who donated last | D3 flat+FK (1.8k) | D3 flat+FK (1.7k) |
| First and last donation of a person | D6 embedded (3.9k) | D3 flat+FK (3.7k) **→** |
| Total donated (global) | D6 embedded (37.9) | D6 embedded (39.8) |
| Total donated (one charity) | D3 flat+FK (324) | D3 flat+FK (309) |
| A person's 20 most recent donations | D7 yb-coloc (3.8k) | D7 yb-coloc (3.7k) |
| Donation by id (point lookup) | D3 flat+FK (3.7k) | D3 flat+FK (3.9k) |
| How many donations a person made | D6 embedded (4.1k) | D6 embedded (3.8k) |
| Charity activity feed (last 50, with names) | D3 flat+FK (172) | D3 flat+FK (200) |
| write: insert donation | D7 yb-coloc (292) | D6 embedded (297) **→** |
| write: donor edits their profile | D3 flat+FK (2.3k) | D3 flat+FK (2.1k) |
| write: erase a donor and all their donations | D6 embedded (273) | D6 embedded (304) |

### Per-question change (cap20 ÷ unbounded)

| Question | D3 flat+FK | D6 embedded | D7 yb-coloc | 
|---|---:|---:|---:|
| Information of the last donation (global) | 0.91x | 0.96x | 0.95x | 
| Information of the last donation (one charity) | 0.96x | 0.72x | 1.03x | 
| Who donates the most | 1.01x | 0.76x | 0.96x | 
| Top-10 donor leaderboard | 0.93x | 0.78x | 1.12x | 
| Who donated last | 0.95x | 0.74x | 1.03x | 
| First and last donation of a person | 0.95x | 0.93x | 0.98x | 
| Total donated (global) | 0.95x | 1.05x | 0.94x | 
| Total donated (one charity) | 0.95x | 0.79x | 0.90x | 
| A person's 20 most recent donations | 0.97x | 0.98x | 0.96x | 
| Donation by id (point lookup) | 1.04x | 1.01x | 0.92x | 
| How many donations a person made | 0.94x | 0.94x | 0.93x | 
| Charity activity feed (last 50, with names) | 1.17x | 0.80x | 1.03x | 

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
| Run id | `20260913-history-cap20` |
| Result files | 7 |
| **Inputs digest** | `57120d55f61035f5` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `57120d55f61035f5` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `57120d55f61035f5` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`57120d55f61035f5`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
