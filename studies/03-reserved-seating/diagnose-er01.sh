#!/usr/bin/env bash
# ER-01 diagnosis: on YugabyteDB, a confirmation's UPDATE sometimes matches none of
# the seats of a hold that the same transaction then reads as valid. This script
# runs the race that reproduces it (S1, 10 000-seat event, repeated trials) under
# chosen configurations and counts the refusals in each:
#
#   A  default                      the configuration every study 03 cell uses
#   B  yb_enable_expression_pushdown=off  (session setting, through the DSN)
#      The refusing statement's whole WHERE clause is pushed to DocDB as a
#      Storage Filter; B evaluates it in the query layer instead.
#   C  enable_wait_queues=false     (tserver flag, node restarted)
#      Replaces wait-on-conflict with fail-on-conflict.
#   D  default, yb_debug_log_internal_restarts=on
#      Logs every statement YugabyteDB restarts internally under READ COMMITTED.
#   E  yb_max_query_layer_retries=0, yb_debug_log_internal_restarts=on
#      No internal statement restart: a conflict reaches the harness as a
#      retryable error, and the harness retries the whole confirmation.
#
# Since round 2 the harness also re-issues the refused statement inside the same
# transaction (then rolls back as before) and records the backend pid, so every
# refusal says whether the miss was transient and which server log lines are its own.
#
# A diagnostic, not a measurement of any design: it changes engine settings the
# study's cells never use. Results go to results/devchecks/er01-<stamp>/.
#
#   ./diagnose-er01.sh [trials] [variants]   # default: 8 trials, variants "A B C"
#
# Takes the benchmark lock.
set -uo pipefail
STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
source "$STUDY_DIR/study.env"
set +e

TRIALS="${1:-8}"
VARIANTS="${2:-A B C}"
need_podman
run_lock_acquire "study ${STUDY_ID} diagnose-er01.sh"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
OUT="$STUDY_DIR/results/devchecks/er01-${STAMP}"
mkdir -p "$OUT"
GIT=(git -C "$(hostpath "$REPO")")
COMMIT="$("${GIT[@]}" rev-parse HEAD)"
DIRTY="false"
[[ -n "$("${GIT[@]}" status --porcelain -- platform infra "studies/$STUDY_ID" ':!studies/*/results' ':!studies/*/reports')" ]] && DIRTY="true"

log "building benchmark image $BENCH_IMAGE"
podman build -q -t "$BENCH_IMAGE" -f "$(hostpath "$STUDY_DIR/Containerfile")" "$(hostpath "$REPO")" >/dev/null || die "image build failed"

BASE_TSERVER_FLAGS="$YB_EXTRA_TSERVER_FLAGS"
DSN_BASE="postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable"

run_variant() {
  local name="$1" tflags="$2" dsn="$3"
  local dir="$OUT/$name"
  mkdir -p "$dir/plans" "$dir/server-logs"
  export YB_EXTRA_TSERVER_FLAGS="$tflags"
  log "[$name] yb-single with tserver flags: $tflags; dsn: $dsn"
  bash "$REPO/infra/yb-single.sh" up || { warn "[$name] yb-single failed to start"; return; }
  podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$dir"):/results" -v "$(hostpath "$dir/plans"):/plans" "$BENCH_IMAGE" \
      -cmd full -engine yugabyte -topology yb-single -design s1_conditional_update -scale tiny \
      -phases verify,explain,race -race-tiers 10000 -race-timeout 20s -race-trials "$TRIALS" \
      -environment "${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}" -run-id "devchecks/er01-${STAMP}-${name}" \
      -repo-commit "$COMMIT" -repo-dirty="$DIRTY" \
      -dsn "$dsn" -out /results/s1_conditional_update.json -explain-out /plans/s1_conditional_update.txt \
    2>&1 | tee "$dir/console.log"
  # The server's own account, kept whole before teardown (study 02 lesson).
  local f
  for f in $(podman exec yb-single bash -c 'find /home/yugabyte/yb_data -path "*logs*" -type f \( -name "*.log" -o -name "*.out" -o -name "*.err" \) -size +0'); do
    podman exec yb-single cat "$f" | gzip > "$dir/server-logs/$(basename "$f").gz"
  done
  bash "$REPO/infra/yb-single.sh" down
}

for v in $VARIANTS; do
  case "$v" in
    A) run_variant A-default "$BASE_TSERVER_FLAGS" "$DSN_BASE" ;;
    B) run_variant B-no-pushdown "$BASE_TSERVER_FLAGS" "${DSN_BASE}&yb_enable_expression_pushdown=off" ;;
    C) run_variant C-no-wait-queues "${BASE_TSERVER_FLAGS},enable_wait_queues=false" "$DSN_BASE" ;;
    D) run_variant D-log-restarts "$BASE_TSERVER_FLAGS" "${DSN_BASE}&yb_debug_log_internal_restarts=on" ;;
    E) run_variant E-no-internal-retries "$BASE_TSERVER_FLAGS" "${DSN_BASE}&yb_debug_log_internal_restarts=on&yb_max_query_layer_retries=0" ;;
    *) warn "unknown variant $v" ;;
  esac
done
export YB_EXTRA_TSERVER_FLAGS="$BASE_TSERVER_FLAGS"

{
  echo "# ER-01 diagnosis — ${STAMP}"
  echo
  echo "- repo_commit: $COMMIT (dirty=$DIRTY); design s1_conditional_update; tiny; race 10 000 seats; ${TRIALS} trials per configuration"
  echo
  echo "| Configuration | Early rejections (sum over trials) | Races with at least one | Trials run | Errors | Seats held/s per trial |"
  echo "|---|---:|---:|---:|---:|---|"
  for dir in "$OUT"/*/; do
    v="$(basename "$dir")"
    f="$dir/console.log"
    [ -f "$f" ] || continue
    early=$(grep -o "early rejections [0-9]*" "$f" | awk '{s+=$3} END{print s+0}')
    races=$(grep -c "early rejections [1-9]" "$f")
    trials=$(grep -c "^    race  10000 seats" "$f")
    errs=$(grep -o "errors=[0-9]*" "$f" | awk -F= '{s+=$2} END{print s+0}')
    rates=$(grep -o "^    race  10000 seats x1 *[0-9.]* seats/s" "$f" | awk '{printf "%s ", $5}')
    echo "| $v | $early | $races | $trials | $errs | $rates |"
  done
  echo
  echo "Re-issued statement inside the refusing transaction (round 2 onwards):"
  echo
  for dir in "$OUT"/*/; do
    grep -o "issued again in this transaction [a-z]* [0-9]* of [0-9]*" "$dir/s1_conditional_update.json" 2>/dev/null | sort | uniq -c | sed "s|^|- $(basename "$dir"): |"
  done
} | tee "$OUT/SUMMARY.md"
log "results: $OUT"
