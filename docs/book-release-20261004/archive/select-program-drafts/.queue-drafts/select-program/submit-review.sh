#!/bin/sh
set -eu
q() { /usr/local/bin/queue "$@" --capability HIGH --role reviewer --model GPT-6 --tool Codex-desktop --effort unknown --session-id session-20261002-select-query-tasks; }
q guard --task task-20261002T131729Z-select-query-task-publication
q submit --task task-20261002T131729Z-select-query-task-publication --summary 'Eight SELECT task briefs published; durable program and check report committed; no benchmarks run.'
q review --task task-20261002T131729Z-select-query-task-publication --review-id review-20261002-select-publication --review-file .queue-drafts/select-program/PUBLICATION_REVIEW.md --summary 'HIGH self-review accepts semantic partitions, bounded execution, dependencies and evidence gates.'
q approve --task task-20261002T131729Z-select-query-task-publication --summary 'Publication accepted after zero-error queue audits and scoped validation; independent empirical review remains a downstream requirement.'
/usr/local/bin/queue audit --json
/usr/local/bin/queue brainstorm-audit --json
