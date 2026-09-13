# Reports

Two kinds of file live here, and the split is deliberate.

## Generated reports — `<run-id>.md`

Produced by the harness from a results directory. **Numbers only, no interpretation.**
Regenerate at any time:

```bash
podman run --rm -v "$PWD/../results:/results" -v "$PWD:/reports" \
  localhost/ticketbench:1 -cmd report -results /results/<run-id> -report-out /reports/<run-id>.md
```

Each one ends with an **inputs digest** — a content hash over every result file in the run.
That digest, not the run id, identifies the data: a run can be extended with extra cells
after the fact, and the digest changes when it is.

It also names the **repository version** every cell was produced from (commit, `git
describe`, and whether the working tree was dirty), and the run tag
`run/02-ticket-booking/<run-id>` when the run was started with `--tag`. That identifies the
code: to revisit a run with a new question, `git checkout` the tag to reproduce it exactly,
or `git diff <tag>` to see what has changed in the study since.

## Analyses — `analyses/<run-id>--<analyst>--<date>.md`

What the numbers *mean*, written by a person or a model, signed with frontmatter that
records who, which version, when, and **which digest they read**.

This project expects several analysts — including different AI models — to interpret the
same data. Different analysts notice different things, and **where two disagree, that
disagreement is a finding** and stays visible rather than being edited away.

To add one: copy [`analyses/TEMPLATE.md`](analyses/TEMPLATE.md), fill the frontmatter, write
your conclusions, then regenerate the report so its index picks you up.

Rules:

- **Check the index first.** If an analysis already exists for this digest, read it. Write
  a new one to disagree, extend, or bring a different perspective — never to restate.
- **Never edit another analyst's file.** Write your own and cite theirs by `analysis_id`.
- **Name where you think the measurement is weak.** An analysis with no stated doubts has
  not been done carefully.

## `outdated/`

When a study is re-run and the conclusions change, the superseded report and its analyses
move here. They are **never deleted and never silently edited**, so any conclusion that was
ever published can be traced back to the run that produced it.
