#!/usr/bin/env bash
# Second pass over run 20260912-small, plus the two focused experiments.
#
# The survey matrix was launched before D9 existed and before repeated trials were
# implemented. Rather than discard three hours of valid measurements, this fills
# the gaps:
#
#   1. D9 (embedded-hybrid) across all three topologies, into the SAME run id.
#      The manifest records it as a separate pass with its own binary id, so a
#      reader can see it was not measured in the first sweep.
#   2. Experiment C — concurrency control on the hot rollup row.
#   3. A repeated-trials run on the controlled pairs whose survey numbers are too
#      close, or too implausible, to argue from. The foreign-key delete result is
#      the specific offender: the survey has the design WITHOUT constraints
#      deleting slower, which cannot be a real effect.
#
# Long-running and unattended. Each stage is independent; a failure in one does
# not stop the others.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"

# lib.sh turns on `set -e`. Turn it back off: this script is deliberately
# resilient -- a failed cell must be recorded and the matrix must continue.
# With -e active, the first non-zero podman exit would kill the whole run
# before the failure-handling code below ever executed.
set +e

SURVEY_RUN="${1:-20260912-small}"
TRIALS_RUN="${2:-20260912-trials}"

log "stage 1/4 — D9 across all topologies, appended to $SURVEY_RUN"
bash "$STUDY_DIR/run-study.sh" \
  --scale small --duration 10s --warmup 3s --conns 8 \
  --designs d9_embedded_hybrid \
  --run-id "$SURVEY_RUN" || warn "stage 1 failed"

log "stage 2/4 — experiment C, concurrency control (pg-single)"
bash "$STUDY_DIR/run-concurrency.sh" \
  --topology pg-single --scale small --duration 8s \
  --conns-list "1,4,8,16,32" || warn "stage 2 failed"

# Five trials on the pairs that carry conclusions. Reads are skipped: the read
# differences in the survey are order-of-magnitude and not in doubt, whereas the
# write differences are tens of percent and sit inside the noise band.
log "stage 3/4 — repeated-trials run on the controlled pairs (writes only)"
bash "$STUDY_DIR/run-study.sh" \
  --scale small --duration 8s --warmup 3s --conns 8 --trials 5 \
  --topologies pg-single \
  --designs d3_flattened_fk,d8_flattened_nofk,d4_rollup_trigger,d9_embedded_hybrid,d6_embedded_jsonb \
  --write-ops insert,update,delete \
  --run-id "$TRIALS_RUN" || warn "stage 3 failed"

log "stage 4/4 — generating reports"
gen() {
  local run="$1" out="$2" mode="${3:-report}"
  podman run --rm \
    -v "$(hostpath "$STUDY_DIR/results"):/results" \
    -v "$(hostpath "$STUDY_DIR/reports"):/reports" \
    "$BENCH_IMAGE" -cmd "$mode" -results "/results/$run" -report-out "/reports/$out" \
    || warn "report generation failed for $run"
}
gen "$SURVEY_RUN" "$SURVEY_RUN.md"
gen "$TRIALS_RUN"  "$TRIALS_RUN.md"

log "done. reports in $STUDY_DIR/reports"
