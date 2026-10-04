#!/bin/sh
set -eu
q() { /usr/local/bin/queue "$@" --capability HIGH --role integrator --model GPT-6 --tool Codex-desktop --effort unknown --session-id session-20261002-select-query-tasks; }
q integrate --task task-20261002T131729Z-select-query-task-publication --tag-message 'SELECT task program: eight dependent tasks for examples, correct answer contracts, controlled crossover measurements and concept-first book publication. Agent: GPT-6 / Codex desktop / effort unknown / HIGH.'
q complete --task task-20261002T131729Z-select-query-task-publication --summary 'Eight SELECT tasks and durable request/program merged into local main and annotated-tagged. No benchmarks or book claims created during publication.'
/usr/local/bin/queue audit --json
/usr/local/bin/queue brainstorm-audit --json
/usr/local/bin/queue list --capability HIGH --eligible
