# Validation receipt — independent HIGH join/CTE review

GPT-6 / Codex desktop / HIGH; selected effort unknown; 2026-10-04.

- Saved-input audit: PASS. Pinned `golang:1.26-bookworm` container, no network,
  repository read-only; `go run ./docs/select-query-templates/reviews/20261004-join-high/audit.go`.
  [Generated output](CHECKS.md) reproduces six original digests, 546 medians,
  orders, samples and controls, plus paired ratios and effective guest memory.
- Source verification: PASS. `podman build -t localhost/charitybench-high-review:20261004
  -f studies/01-charity-tree/Containerfile studies/01-charity-tree` executed
  `go vet ./... && go test ./...` and compiled the benchmark client. Root module:
  no test files; `charitytree/harness`: PASS (0.153 s);
  `charitytree/selectbench`: PASS (0.270 s).
- Producing-image report regeneration: FAIL (expected defect reproduction).
  Original image `localhost/charitybench-select:select-join-pg1-20261004`, unchanged
  input mounted read-only; `-cmd join-report -results /inputs -report-out /review/rerender-pg1.md`.
  [Output](rerender-pg1.md) advertises the wrong digest; the original report/input
  remain unchanged. See R4 in the review. No new measured run was performed.
- Queue baseline: 91 tasks, 0 errors, 2 existing completed-T2 live-ref warnings
  (`expired_lease` and `live_ref_drift` for `task-20261001T133759Z-t2-study02-replication`).
  This review does not recover or mutate that unrelated task.
- Lock inspection: `ads-run-lock` absent; no database starts by HIGH review.
- All ten original producing tags were inventoried before review in
  [TAGS-BEFORE.txt](TAGS-BEFORE.txt). Completion tag `study-01/select-v2-join-cte`
  is absent, as required while approval is withheld.

Outcome: a committed, tagged HIGH correction checkpoint, with queue transition
to `changes_requested`. This is not the task completion or a local-main
measurement merge. LOW reclaims the retained task branch and creates a new attempt.
