# Outdated reports

Superseded reports and analyses, kept rather than deleted.

A report describes one run, of one code state, on one machine. When a study is re-run and
the conclusions change, the old report moves here instead of being overwritten — so that a
claim someone read six months ago can still be traced to the data that produced it.

Nothing in this directory should be treated as current. Each file keeps its original
`run-id`, and its inputs digest still identifies exactly which measurements it was drawn
from.

The v3 reporter preserves regenerated editions as `<run-id>.md.previous-N.md`.
Their former location was `reports/<run-id>.md`; the numbered archive keeps the
previous bytes verbatim. Resolve its relative links from that original `reports/`
directory. The current edition at the former location links to the same measurement
inputs and records the current analysis index. An archived edition may therefore
have an unchanged input digest and an older presentation or index.
