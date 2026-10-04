#!/bin/sh
set -eu
q() { /usr/local/bin/queue "$@" --capability HIGH --role planner --model GPT-6 --tool Codex-desktop --effort unknown --session-id session-20261002-select-query-tasks; }
q guard --task task-20261002T131729Z-select-query-task-publication
q publish --spec .queue-drafts/select-program/select-protocol/task.json --brief .queue-drafts/select-program/select-protocol/BRIEF.md
q publish --spec .queue-drafts/select-program/select-harness/task.json --brief .queue-drafts/select-program/select-harness/BRIEF.md
q publish --spec .queue-drafts/select-program/select-join-cte/task.json --brief .queue-drafts/select-program/select-join-cte/BRIEF.md
q publish --spec .queue-drafts/select-program/select-existence/task.json --brief .queue-drafts/select-program/select-existence/BRIEF.md
q publish --spec .queue-drafts/select-program/select-aggregate-topn/task.json --brief .queue-drafts/select-program/select-aggregate-topn/BRIEF.md
q publish --spec .queue-drafts/select-program/select-access-paths/task.json --brief .queue-drafts/select-program/select-access-paths/BRIEF.md
q publish --spec .queue-drafts/select-program/select-synthesis/task.json --brief .queue-drafts/select-program/select-synthesis/BRIEF.md
q publish --spec .queue-drafts/select-program/select-book/task.json --brief .queue-drafts/select-program/select-book/BRIEF.md
q audit --json
q brainstorm-audit --json
