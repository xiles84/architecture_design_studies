#!/usr/bin/env bash
# SELECT experiment runner (study 01, additive).
#
# Runs the `-cmd select` mode of the study 01 client against PostgreSQL
# single-node: it loads the deterministic teaching fixture into its own
# `selectbench` schema, executes the embedded catalogue, verifies every cell
# against the independent Go oracle, captures EXPLAIN (ANALYZE, BUFFERS) and
# writes the measurements-only report.
#
# It is additive: the historical Study 01 modes, designs and results are
# untouched, and this runner uses its own image tag and output directory.
#
# It starts a database, so it takes the benchmark lock (measurements are serial
# on this machine) and always stops the database and releases the lock.
#
# Usage:
#   studies/01-charity-tree/run-select.sh [--scale tiny|small|medium|large] [--seed N]
#                                         [--keep-up]
#
# Environment: SCALE, SEED, KEEP_UP, OUTDIR.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$HERE/../.." && pwd)"
# shellcheck source=../../infra/lib.sh
source "$REPO/infra/lib.sh"

SCALE="${SCALE:-tiny}"
SEED="${SEED:-42}"
KEEP_UP="${KEEP_UP:-no}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scale) SCALE="$2"; shift 2 ;;
    --seed) SEED="$2"; shift 2 ;;
    --keep-up) KEEP_UP="yes"; shift ;;
    *) echo "run-select.sh: unknown argument $1" >&2; exit 2 ;;
  esac
done

RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
OUTDIR="${OUTDIR:-$HERE/results/select-$RUN_ID}"
IMAGE="localhost/charitybench-select:1"

run_lock_acquire "study-01 run-select.sh scale=$SCALE seed=$SEED"

cleanup() {
  local rc=$?
  if [[ "$KEEP_UP" != "yes" ]]; then
    bash "$REPO/infra/pg-single.sh" down >/dev/null 2>&1 || true
  fi
  run_lock_release >/dev/null 2>&1 || true
  return $rc
}
trap cleanup EXIT

echo "==> starting pg-single"
bash "$REPO/infra/pg-single.sh" up
DSN="$(bash "$REPO/infra/pg-single.sh" dsn)"
echo "==> dsn: ${DSN/@*/@<hidden>}"

echo "==> building $IMAGE"
podman build -q -t "$IMAGE" -f "$(buildpath "$HERE/Containerfile")" "$(buildpath "$HERE")" >/dev/null

mkdir -p "$OUTDIR"
echo "==> select scale=$SCALE seed=$SEED -> $OUTDIR"
podman run --rm \
  --network "$NETWORK" \
  -v "$(hostpath "$OUTDIR"):/results" \
  -e ADS_REPO_COMMIT="$(git -C "$REPO" rev-parse HEAD 2>/dev/null || echo unknown)" \
  "$IMAGE" \
    -cmd select \
    -engine postgres \
    -topology single \
    -scale "$SCALE" \
    -seed "$SEED" \
    -dsn "$DSN" \
    -out "/results/select.json" \
    -report-out "/results/select.md"

echo "==> results in $OUTDIR"
