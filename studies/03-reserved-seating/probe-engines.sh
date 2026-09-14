#!/usr/bin/env bash
# Ask the pinned engines what study 03's designs rely on (HANDOFF §4), before any
# design SQL is trusted. Starts pg-single and yb-single in turn, runs probe/ in a
# Go container on the benchmark network, and writes one transcript under
# results/devchecks/.
#
#   ./probe-engines.sh
#
# Takes the benchmark lock: it starts databases.
set -uo pipefail
STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
source "$STUDY_DIR/study.env"
export YB_EXTRA_TSERVER_FLAGS
set +e

need_podman
run_lock_acquire "study ${STUDY_ID} probe-engines.sh"

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
OUT="$STUDY_DIR/results/devchecks/engine-probe-${STAMP}.md"
mkdir -p "$(dirname "$OUT")"
GIT=(git -C "$(hostpath "$REPO")")
{
  echo "# Study 03 engine probe — ${STAMP}"
  echo
  echo "- repo_commit: $("${GIT[@]}" rev-parse HEAD)"
  echo "- repo_dirty (study code): $([[ -n "$("${GIT[@]}" status --porcelain -- platform infra "studies/$STUDY_ID" ':!studies/*/results')" ]] && echo true || echo false)"
  echo "- images: $PG_IMAGE, $YB_IMAGE"
  echo "- yb_extra_tserver_flags: $YB_EXTRA_TSERVER_FLAGS"
  echo "- budget per db node: cpus=$DB_CPUS memory=$DB_MEMORY"
  echo
} > "$OUT"

probe() {
  local engine="$1" dsn="$2"
  podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$REPO"):/src" -v ads-gomodcache:/go/pkg/mod \
    -w "/src/studies/$STUDY_ID" -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false \
    "$GO_IMAGE" go run ./probe -engine "$engine" -dsn "$dsn" 2>&1 | tee -a "$OUT"
}

bash "$REPO/infra/pg-single.sh" up && probe postgres "postgres://bench:bench@pg-single:5432/bench?sslmode=disable"
bash "$REPO/infra/pg-single.sh" down
bash "$REPO/infra/yb-single.sh" up && probe yugabyte "postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable"
bash "$REPO/infra/yb-single.sh" down

log "transcript: $OUT"
