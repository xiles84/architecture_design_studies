#!/usr/bin/env bash
# Additive SELECT validation runner. All database access is serialized by the
# shared lock; immutable source/image identities accompany the new results.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$HERE/../.." && pwd)"
source "$REPO/infra/lib.sh"
SCALE="${SCALE:-tiny}"
SEED="${SEED:-42}"
ENGINE="postgres"
TRIALS=3
FAMILY="select"
PAGE_LIMIT=5
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
TAG=""
ENVIRONMENT="${BENCH_ENVIRONMENT:-}"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --scale) SCALE="$2"; shift 2 ;;
    --seed) SEED="$2"; shift 2 ;;
    --engine) ENGINE="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --family) FAMILY="$2"; shift 2 ;;
    --page-limit) PAGE_LIMIT="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --tag) TAG="$2"; shift 2 ;;
    --environment) ENVIRONMENT="$2"; shift 2 ;;
    *) die "unknown argument $1" ;;
  esac
done
[[ "$SCALE" =~ ^(tiny|small|medium|large|deep)$ ]] || die "unknown scale $SCALE"
[[ "$FAMILY" =~ ^(select|access)$ ]] || die "family must be select or access"
[[ "$TRIALS" =~ ^[1-9][0-9]*$ ]] || die "trials must be positive"
[[ "$PAGE_LIMIT" =~ ^[1-9][0-9]*$ ]] || die "page limit must be positive"
[[ "$RUN_ID" =~ ^[a-zA-Z0-9_-]+$ ]] || die "invalid run id"
[[ "$TAG" == run/study-01/* ]] || die "provide an immutable --tag run/study-01/<id>"
[[ -n "$ENVIRONMENT" && -f "$REPO/docs/environments/$ENVIRONMENT.md" ]] || die "provide a documented --environment"
if [[ "$FAMILY" == access ]]; then
  CLIENT_CMD="select-access"; OUT_JSON="access.json"; OUT_MD="access.md"
else
  CLIENT_CMD="select"; OUT_JSON="select.json"; OUT_MD="select.md"
fi
case "$ENGINE" in
  postgres) TOPOLOGY_SCRIPT="$REPO/infra/pg-single.sh"; NODE=pg-single; ENGINE_IMAGE="$PG_IMAGE"; VOLUME=ads-pg-single-data ;;
  yugabyte) TOPOLOGY_SCRIPT="$REPO/infra/yb-single.sh"; NODE=yb-single; ENGINE_IMAGE="$YB_IMAGE"; VOLUME=ads-yb-single-data ;;
  *) die "engine must be postgres or yugabyte" ;;
esac
OUTDIR="${OUTDIR:-$HERE/results/$RUN_ID}"
[[ ! -e "$OUTDIR" ]] || die "output already exists: $OUTDIR"
GIT_ROOT="$(hostpath "$REPO")"
# Exclude this study's results/ and reports/ explicitly: git's `:!` pathspec
# does not suppress UNTRACKED files, so a previous run's output would otherwise
# make the tree look dirty and block every later run (see LESSONS_LEARNED).
STATUS="$(git -C "$GIT_ROOT" status --porcelain --untracked-files=all -- infra platform studies/01-charity-tree | grep -vE 'studies/01-charity-tree/(results|reports)/' || true)"
[[ -z "$STATUS" ]] || die "commit source and reports before a SELECT run: $STATUS"
SOURCE_COMMIT="$(git -C "$GIT_ROOT" rev-parse HEAD)"
IMAGE="localhost/charitybench-select:$RUN_ID"
run_lock_acquire "study-01 SELECT $FAMILY $RUN_ID $ENGINE"
STARTED=no
cleanup() {
  local rc=$?
  if [[ "$STARTED" == yes ]]; then
    podman logs "$NODE" >"$OUTDIR/resources/database.log" 2>&1 || true
    bash "$TOPOLOGY_SCRIPT" down >/dev/null 2>&1 || true
  fi
  run_lock_release >/dev/null 2>&1 || true
  return "$rc"
}
trap cleanup EXIT
! podman container exists "$NODE" || die "$NODE already exists; preserve its owner"
! podman volume exists "$VOLUME" || die "$VOLUME already exists; inspect its owner"
! git -C "$GIT_ROOT" rev-parse --verify "refs/tags/$TAG" >/dev/null 2>&1 || die "tag exists; choose a new run id"
podman build -t "$IMAGE" -f "$(buildpath "$(hostpath "$HERE/Containerfile")")" "$(buildpath "$(hostpath "$HERE")")"
IMAGE_ID="$(podman image inspect "$IMAGE" --format '{{.Id}}')"
ENGINE_IMAGE_ID="$(podman image inspect "$ENGINE_IMAGE" --format '{{.Id}}')"
git -C "$GIT_ROOT" tag -a "$TAG" "$SOURCE_COMMIT" -m "Study 01 SELECT validation $RUN_ID; source $SOURCE_COMMIT; actor in task/run commit archive"
mkdir -p "$OUTDIR/resources"
printf 'environment=%s\nrepo_commit=%s\nrun_tag=%s\nengine=%s\nengine_image=%s\nengine_image_id=%s\nclient_image=%s\nclient_image_id=%s\nnode_budget=%s CPU %s\nclient_budget=%s CPU %s\ntrials=%s\nwhole_runs=1\n' \
  "$ENVIRONMENT" "$SOURCE_COMMIT" "$TAG" "$ENGINE" "$ENGINE_IMAGE" "$ENGINE_IMAGE_ID" "$IMAGE" "$IMAGE_ID" \
  "$DB_CPUS" "$DB_MEMORY" "$CLIENT_CPUS" "$CLIENT_MEMORY" "$TRIALS" >"$OUTDIR/manifest.env"
podman machine inspect >"$OUTDIR/resources/machine.json"
podman image inspect "$IMAGE" "$ENGINE_IMAGE" >"$OUTDIR/resources/images.json"
STARTED=yes
bash "$TOPOLOGY_SCRIPT" up
DSN="$(bash "$TOPOLOGY_SCRIPT" dsn)"
podman inspect "$NODE" >"$OUTDIR/resources/node.json"
podman exec "$NODE" cat /proc/meminfo >"$OUTDIR/resources/guest-memory.txt"
podman exec "$NODE" cat /sys/devices/system/cpu/online >"$OUTDIR/resources/guest-cpus.txt"
podman exec "$NODE" cat /sys/fs/cgroup/cpu.stat >"$OUTDIR/resources/node-cpu-before.txt"
set +e
podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
  -v "$(winpath "$(hostpath "$OUTDIR")"):/results" \
  -e "ADS_RUN_LOCK_HELD=$ADS_RUN_LOCK_HELD" -e "ADS_REPO_COMMIT=$SOURCE_COMMIT" \
  -e "ADS_RUN_TAG=$TAG" -e "ADS_BENCH_IMAGE=$IMAGE" -e "ADS_BENCH_IMAGE_ID=$IMAGE_ID" -e "ADS_ENGINE_IMAGE_ID=$ENGINE_IMAGE_ID" \
  "$IMAGE" -cmd "$CLIENT_CMD" -engine "$ENGINE" -topology single -scale "$SCALE" -seed "$SEED" \
  -trials "$TRIALS" -page-limit "$PAGE_LIMIT" -environment "$ENVIRONMENT" -run-id "$RUN_ID" -dsn "$DSN" \
  -out "/results/$OUT_JSON" -report-out "/results/$OUT_MD" >"$OUTDIR/resources/client.log" 2>&1
RESULT_RC=$?
set -e
cat "$OUTDIR/resources/client.log"
podman exec "$NODE" cat /sys/fs/cgroup/cpu.stat >"$OUTDIR/resources/node-cpu-after.txt"
printf 'exit_code=%s\n' "$RESULT_RC" >>"$OUTDIR/manifest.env"
exit "$RESULT_RC"
