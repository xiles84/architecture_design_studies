# Reports

Measurements, concise final analyses, and detailed discussion companions live here.

For the completed v3 enhancements, start with the
[final analysis](analyses/20260914-study01-v3--gpt-6--2026-09-14.md).
The [discussion index](discussions/README.md) leads to the full mechanisms and responses
to earlier analysts. Original analyses remain available for their original input sets.

## Generated reports — `<run-id>.md`

Produced by the harness from a results directory. **Numbers only, no interpretation.**
Regenerate at any time:

```bash
podman run --rm -v "$PWD/../results:/results" -v "$PWD:/reports" \
  localhost/charitybench:1 -cmd report -results /results/<run-id> -report-out /reports/<run-id>.md
```

Each one ends with an **inputs digest** — a content hash over every result file in the run.
That digest, not the run id, identifies the data: a run can be extended with extra cells
after the fact, and the digest changes when it is.

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

## Discussions — `discussions/<topic>--<analyst>--<date>.md`

Detailed SQL and plan explanations, design comparisons, and exchanges between models or
human analysts belong here. Use [the discussion template](../../../docs/templates/DISCUSSION.md).
Keep final analyses around 1,000 words and link companions beside the relevant conclusions.
Both documents carry author, run, digest and commit provenance and link to one another.

## `outdated/`

When a study is re-run and the conclusions change, the superseded report and its analyses
move here. They are **never deleted and never silently edited**, so any conclusion that was
ever published can be traced back to the run that produced it.
