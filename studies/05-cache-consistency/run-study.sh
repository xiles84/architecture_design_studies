#!/usr/bin/env bash
# Run study 05 across scenarios and topologies.
#
#   ./run-study.sh --scale tiny --topologies pg-single --phases verify,explain,calibrate,warm
#   ./run-study.sh --tag --scale small --topologies pg-single
#   ./run-study.sh --tag --scale medium --topologies pg-single \
#       --designs legacy-na-redis-aside-relaxed-coord --phases verify,calibrate,warm,mixed,churn \
#       --churn-duration 11m
#
# Each (topology, scenario) pair is one independent container run writing its own
# JSON. A matrix that fails on one cell leaves every other cell's result behind;
# failures are recorded in the manifest, with the database and cache logs captured
# before teardown.
#
# Provenance: the repository commit, `git describe`, and whether the study's code had
# uncommitted changes are captured ONCE, before anything runs, and passed into every
# cell and the manifest. With --tag, a clean tree is tagged so the run can always be
# checked out again; a dirty tree is never tagged.
set -uo pipefail

# Run from a private copy of this script. Bash reads a script incrementally while it
# runs, so editing the file mid-matrix can garble what it executes next
# (LESSONS_LEARNED: "Never edit a bash script that is currently running").
if [[ -z "${ADS_RUNNER_COPY:-}" ]]; then
  _study_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  _copy="$(mktemp -t run-study.XXXXXX)"
  cp "${BASH_SOURCE[0]}" "$_copy"
  ADS_RUNNER_COPY="$_copy" ADS_STUDY_DIR="$_study_dir" exec bash "$_copy" "$@"
fi
STUDY_DIR="$ADS_STUDY_DIR"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
source "$STUDY_DIR/study.env"
export YB_EXTRA_TSERVER_FLAGS
# lib.sh turns on `set -e`; a failed cell must be recorded, not abort the matrix.
set +e

SCALE="small"
DURATION=""
WARMUP=""
CONNS="8"
WRITE_CONNS="2"
TRIALS="1"
INSTANCES="1"
CHURN_DURATION="0"
STAMPEDE="16"
STAMPEDE_KEYS="8"
HOT_KEYS="4"
CAPACITY_KB="0"
CACHE_FITS="no"
REDIS_MAXMEMORY_MB="32"
FRAME="db-only"
SEED="42"
FAULT_SEED="4242"
SAMPLE="12"
TOPOLOGIES="pg-single"
DESIGNS=""
PHASES="verify,explain,calibrate,warm,mixed,hotspot,stampede,instances,churn,faults,audit"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
TAG="no"
EXTRA=""
ORDER_SEED=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scale) SCALE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --warmup) WARMUP="$2"; shift 2 ;;
    --conns) CONNS="$2"; shift 2 ;;
    --write-conns) WRITE_CONNS="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --instances) INSTANCES="$2"; shift 2 ;;
    --churn-duration) CHURN_DURATION="$2"; shift 2 ;;
    --stampede-readers) STAMPEDE="$2"; shift 2 ;;
    --stampede-keys) STAMPEDE_KEYS="$2"; shift 2 ;;
    --hot-keys) HOT_KEYS="$2"; shift 2 ;;
    --cache-capacity-kib) CAPACITY_KB="$2"; shift 2 ;;
    --cache-fits) CACHE_FITS="yes"; shift ;;
    --redis-maxmemory-mb) REDIS_MAXMEMORY_MB="$2"; shift 2 ;;
    --resource-frame) FRAME="$2"; shift 2 ;;
    --seed) SEED="$2"; shift 2 ;;
    --fault-seed) FAULT_SEED="$2"; shift 2 ;;
    --sample) SAMPLE="$2"; shift 2 ;;
    --topologies) TOPOLOGIES="$2"; shift 2 ;;
    --designs) DESIGNS="$2"; shift 2 ;;
    --phases) PHASES="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --tag) TAG="yes"; shift ;;
    --extra) EXTRA="$2"; shift 2 ;;
    --order-seed) ORDER_SEED="$2"; shift 2 ;;
    -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
    *) die "unknown flag $1" ;;
  esac
done

# ---------------------------------------------------------------------------
# Resource framings. The two conditions are NEVER pooled: db-only and equal-total
# keep the same total SERVICE budget, add-cache spends more. The cache's memory
# limit is in bytes because the memory backend's resident footprint is charged
# against it and a rounded "0.75g" would be a number nobody could audit.
# ---------------------------------------------------------------------------
case "$FRAME" in
  db-only)     export ADS_DB_CPUS="${ADS_DB_CPUS:-2}" ADS_DB_MEMORY="${ADS_DB_MEMORY:-3g}" ;;
  add-cache)   export ADS_DB_CPUS="${ADS_DB_CPUS:-2}" ADS_DB_MEMORY="${ADS_DB_MEMORY:-3g}" \
                      ADS_CACHE_CPUS="${ADS_CACHE_CPUS:-1}" ADS_CACHE_MEMORY="${ADS_CACHE_MEMORY:-805306368}" ;;
  equal-total) export ADS_DB_CPUS="${ADS_DB_CPUS:-1}" ADS_DB_MEMORY="${ADS_DB_MEMORY:-2415919104}" \
                      ADS_CACHE_CPUS="${ADS_CACHE_CPUS:-1}" ADS_CACHE_MEMORY="${ADS_CACHE_MEMORY:-805306368}" ;;
  *) die "unknown resource frame $FRAME" ;;
esac
export ADS_REDIS_MAXMEMORY="${REDIS_MAXMEMORY_MB}mb"

ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
[[ -z "$ORDER_SEED" ]] && ORDER_SEED="$(printf '%s' "$RUN_ID" | cksum | cut -d' ' -f1)"
OUT="$STUDY_DIR/results/$RUN_ID"
mkdir -p "$OUT"

# ---------------------------------------------------------------------------
# Repository version. "Dirty" means uncommitted changes in the code that can
# change THIS study's results: the shared platform, the infra scripts and the study
# itself. Results, reports and analyses are excluded -- they are what a run
# produces -- and so are other studies and the docs, which several analysts may be
# editing at the same time without affecting a single number here.
# ---------------------------------------------------------------------------
GIT=(git -C "$(hostpath "$REPO")")
REPO_COMMIT="$("${GIT[@]}" rev-parse HEAD 2>/dev/null)" \
  || die "cannot read the repository commit -- refusing to produce results that cannot be traced to code"
CODE_PATHS=(platform infra "studies/$STUDY_ID" .containerignore ':!studies/*/results' ':!studies/*/reports')
STATUS="$("${GIT[@]}" status --porcelain -- "${CODE_PATHS[@]}")" \
  || die "git status failed -- cannot tell whether the study code is committed"
REPO_DIRTY="false"
[[ -n "$STATUS" ]] && REPO_DIRTY="true"
REPO_DESCRIBE="$("${GIT[@]}" describe --tags --always 2>/dev/null || echo "$REPO_COMMIT")"
[[ "$REPO_DIRTY" == "true" ]] && REPO_DESCRIBE="${REPO_DESCRIBE}-dirty"
RUN_TAG=""
if [[ "$TAG" == "yes" ]]; then
  if [[ "$REPO_DIRTY" == "true" ]]; then
    warn "study code has uncommitted changes; NOT tagging (commit first)"
    "${GIT[@]}" status --porcelain -- "${CODE_PATHS[@]}" | head -20 >&2
  else
    RUN_TAG="run/${STUDY_ID}/${RUN_ID}"
    if "${GIT[@]}" rev-parse -q --verify "refs/tags/$RUN_TAG" >/dev/null; then
      log "tag $RUN_TAG already exists (extending a run); leaving it in place"
    else
      "${GIT[@]}" tag -a "$RUN_TAG" -m "Study ${STUDY_ID}: code that produced run ${RUN_ID}" \
        && log "tagged $REPO_COMMIT as $RUN_TAG"
    fi
  fi
fi
[[ "$REPO_DIRTY" == "true" ]] && warn "running from a dirty tree: results will say so"

# One measurement on this machine at a time, across every worktree and session
# (infra/lib.sh). Taken before anything builds an image or starts a container.
run_lock_acquire "study ${STUDY_ID} run-study.sh run ${RUN_ID}"

BUSY="$(podman ps --format '{{.Names}}' | grep -E '^(pg-single|yb-single|yb-n[123]|ads-redis)$' || true)"
[[ -n "$BUSY" ]] && die "database or cache containers already running ($BUSY) -- another run may be in progress"

need_podman
log "building benchmark image $BENCH_IMAGE (context: repository root)"
# winpath, not hostpath, for the build. A bind-mount source tolerates /mnt/c/...
# because the engine resolves it inside the machine, but the *build context* is
# resolved by the podman client on Windows, which reads a leading "/" as the root of
# the current drive. winpath() is a no-op under Git Bash (where hostpath already
# returned the C:/ form), so this is correct on both.
podman build -q -t "$BENCH_IMAGE" -f "$(winpath "$STUDY_DIR/Containerfile")" "$(winpath "$REPO")" >/dev/null \
  || die "image build failed"

BENCH_IMAGE_ID="$(podman image inspect "$BENCH_IMAGE" --format '{{.Id}}' 2>/dev/null | cut -c1-19)"
REDIS_IMAGE_ID="$(podman image inspect "$REDIS_IMAGE" --format '{{.Id}}' 2>/dev/null | cut -c1-19)"

# The scenario registry is read from the binary, so the runner and the registry
# cannot disagree about which scenarios exist.
ALL_DESIGNS="$(podman run --rm "$BENCH_IMAGE" -cmd list | awk '{print $1}' | tr '\n' ' ')"

{
  echo "- pass_started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "  study: $STUDY_ID"
  echo "  run_id: $RUN_ID"
  echo "  repo_commit: $REPO_COMMIT"
  echo "  repo_describe: $REPO_DESCRIBE"
  echo "  repo_dirty: $REPO_DIRTY"
  echo "  run_tag: ${RUN_TAG:-none}"
  echo "  environment: $ENVIRONMENT"
  echo "  scale: $SCALE"
  echo "  seed: $SEED"
  echo "  fault_seed: $FAULT_SEED"
  echo "  phases: $PHASES"
  echo "  duration_per_measurement: ${DURATION:-<scale default>}"
  echo "  warmup: ${WARMUP:-<scale default>}"
  echo "  client_conns: $CONNS"
  echo "  writer_conns: $WRITE_CONNS"
  echo "  trials: $TRIALS"
  echo "  logical_instances: $INSTANCES"
  echo "  churn_duration: $CHURN_DURATION"
  echo "  stampede_readers: $STAMPEDE"
  echo "  stampede_keys: $STAMPEDE_KEYS"
  echo "  hot_keys: $HOT_KEYS"
  echo "  cache_capacity_kib: $CAPACITY_KB"
  echo "  cache_fits_working_set: $CACHE_FITS"
  echo "  sample: $SAMPLE"
  echo "  resource_framing: $FRAME"
  echo "  budget_database: cpus=$DB_CPUS memory=$DB_MEMORY"
  echo "  budget_cache: cpus=${ADS_CACHE_CPUS:-none} memory=${ADS_CACHE_MEMORY:-none}"
  echo "  budget_client: cpus=$CLIENT_CPUS memory=$CLIENT_MEMORY"
  echo "  redis_maxmemory: ${ADS_REDIS_MAXMEMORY}"
  echo "  order_seed: $ORDER_SEED"
  echo "  topologies: $TOPOLOGIES"
  echo "  designs: ${DESIGNS:-all}"
  echo "  images:"
  echo "    postgres: $PG_IMAGE"
  echo "    yugabyte: $YB_IMAGE"
  echo "    redis: $REDIS_IMAGE"
  echo "  bench_image: $BENCH_IMAGE"
  echo "  bench_image_id: $BENCH_IMAGE_ID"
  echo "  redis_image_id: $REDIS_IMAGE_ID"
  echo "  podman: $(podman --version)"
  echo "  podman_client: $ADS_PODMAN"
  echo "  extra_flags: \"$EXTRA\""
} >> "$OUT/manifest.yaml"

FAILED=()

run_cell() {
  local topo="$1" engine="$2" dsn="$3" design="$4"
  local ENGINE_FLAGS="$PG_HARNESS_FLAGS"
  [[ "$engine" == "yugabyte" ]] && ENGINE_FLAGS="$YB_HARNESS_FLAGS"
  local jdir="$OUT/$topo" pdir="$OUT/$topo/plans" ldir="$OUT/$topo/logs"
  mkdir -p "$jdir" "$pdir" "$ldir"
  local names c
  case "$topo" in
    pg-single) names="pg-single" ;;
    yb-single) names="yb-single" ;;
    yb-cluster3) names="yb-n1 yb-n2 yb-n3" ;;
  esac
  for c in $names; do
    { echo "[$c] before $(date -u +%H:%M:%S)"; podman exec "$c" cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
  done
  if [[ "$design" == *redis* ]]; then
    { echo "[ads-redis] before $(date -u +%H:%M:%S)"; podman exec ads-redis cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
  fi
  log "[$topo] $design"
  # shellcheck disable=SC2086
  podman run --rm --network "$NETWORK" \
    --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$jdir"):/results" -v "$(hostpath "$pdir"):/plans" \
    "$BENCH_IMAGE" \
      -cmd full -engine "$engine" -topology "$topo" -scenario "$design" \
      -scale "$SCALE" -seed "$SEED" -fault-seed "$FAULT_SEED" \
      -conns "$CONNS" -write-conns "$WRITE_CONNS" -trials "$TRIALS" \
      -instances "$INSTANCES" -churn-duration "$CHURN_DURATION" \
      -stampede-readers "$STAMPEDE" -stampede-keys "$STAMPEDE_KEYS" -hot-keys "$HOT_KEYS" \
      -cache-capacity-kib "$CAPACITY_KB" -redis-maxmemory-mb "$REDIS_MAXMEMORY_MB" \
      -redis-addr "ads-redis:6379" -resource-frame "$FRAME" -sample "$SAMPLE" \
      -phases "$PHASES" \
      ${DURATION:+-duration "$DURATION"} ${WARMUP:+-warmup "$WARMUP"} \
      ${CACHE_FITS:+$([[ "$CACHE_FITS" == "yes" ]] && echo -cache-fits)} \
      -environment "$ENVIRONMENT" -run-id "$RUN_ID" \
      -repo-commit "$REPO_COMMIT" -repo-describe "$REPO_DESCRIBE" -repo-dirty="$REPO_DIRTY" -run-tag "$RUN_TAG" \
      -bench-image "$BENCH_IMAGE" -bench-image-id "$BENCH_IMAGE_ID" \
      $ENGINE_FLAGS $EXTRA \
      -dsn "$dsn" -out "/results/$design.json" -explain-out "/plans/$design.txt" \
    2>&1 | tee "$ldir/$design.log"
  local status=${PIPESTATUS[0]}
  for c in $names; do
    { echo "[$c] after $(date -u +%H:%M:%S)"; podman exec "$c" cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
  done
  if [[ "$design" == *redis* ]]; then
    { echo "[ads-redis] after $(date -u +%H:%M:%S)"; podman exec ads-redis cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
    podman exec ads-redis redis-cli info >> "$ldir/$design.redis.info" 2>&1 || true
  fi
  if [[ $status -ne 0 ]]; then
    warn "[$topo] $design FAILED -- continuing with the rest of the matrix"
    for c in $names; do
      podman logs --tail 300 "$c" > "$ldir/$design.$c.log" 2>&1 || true
    done
    FAILED+=("$topo/$design")
  fi
}

designs_for() {
  local topo="$1" list="${DESIGNS//,/ }" d out=()
  for d in ${list:-$ALL_DESIGNS}; do
    if [[ "$topo" == "pg-single" && " $YB_ONLY_DESIGNS " == *" $d "* ]]; then
      continue
    fi
    out+=("$d")
  done
  printf '%s\n' "${out[@]}" | awk -v seed="$ORDER_SEED" -v topo="$topo" '
    BEGIN {
      s = seed % 100000
      for (i = 1; i <= length(topo); i++) s += i * index("abcdefghijklmnopqrstuvwxyz0123456789-", substr(topo, i, 1))
      srand(s)
    }
    { printf "%.12f %s\n", rand(), $0 }' | sort -k1,1 | cut -d' ' -f2 | tr '\n' ' '
}

needs_redis() {
  local list="${DESIGNS//,/ }" d
  for d in ${list:-$ALL_DESIGNS}; do
    [[ "$d" == *red* ]] && return 0
  done
  return 1
}

for topo in ${TOPOLOGIES//,/ }; do
  ORDER="$(designs_for "$topo")"
  echo "  design_order_${topo}: $ORDER" >> "$OUT/manifest.yaml"
  REDIS_UP="no"
  case "$topo" in
    pg-single)
      bash "$REPO/infra/pg-single.sh" up || { warn "pg-single failed to start"; FAILED+=("pg-single/*"); continue; }
      need_redis && { bash "$REPO/infra/redis.sh" up || { warn "redis failed to start"; FAILED+=("*/*"); bash "$REPO/infra/pg-single.sh" down; continue; }; REDIS_UP="yes"; }
      record_topology "$OUT/pg-single/topology.yaml" pg-single ads-redis
      for d in $ORDER; do
        run_cell pg-single postgres "postgres://bench:bench@pg-single:5432/bench?sslmode=disable" "$d"
      done
      [[ "$REDIS_UP" == "yes" ]] && bash "$REPO/infra/redis.sh" down
      bash "$REPO/infra/pg-single.sh" down
      ;;
    yb-single)
      bash "$REPO/infra/yb-single.sh" up || { warn "yb-single failed to start"; FAILED+=("yb-single/*"); continue; }
      need_redis && { bash "$REPO/infra/redis.sh" up || { warn "redis failed to start"; FAILED+=("*/*"); bash "$REPO/infra/yb-single.sh" down; continue; }; REDIS_UP="yes"; }
      record_topology "$OUT/yb-single/topology.yaml" yb-single ads-redis
      for d in $ORDER; do
        run_cell yb-single yugabyte "postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable" "$d"
      done
      [[ "$REDIS_UP" == "yes" ]] && bash "$REPO/infra/redis.sh" down
      bash "$REPO/infra/yb-single.sh" down
      ;;
    yb-cluster3)
      bash "$REPO/infra/yb-cluster3.sh" up || { warn "yb-cluster3 failed to start"; FAILED+=("yb-cluster3/*"); continue; }
      need_redis && { bash "$REPO/infra/redis.sh" up || { warn "redis failed to start"; FAILED+=("*/*"); bash "$REPO/infra/yb-cluster3.sh" down; continue; }; REDIS_UP="yes"; }
      record_topology "$OUT/yb-cluster3/topology.yaml" yb-n1 yb-n2 yb-n3 ads-redis
      # Connections are spread over ALL query endpoints, and the endpoint
      # distribution is recorded, so a single-endpoint bottleneck cannot be read as
      # a general limit of the cluster.
      for d in $ORDER; do
        run_cell yb-cluster3 yugabyte "postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable,postgres://yugabyte@yb-n2:5433/yugabyte?sslmode=disable,postgres://yugabyte@yb-n3:5433/yugabyte?sslmode=disable" "$d"
      done
      podman exec yb-n1 bash -lc "bin/yb-admin --master_addresses=yb-n1:7100,yb-n2:7100,yb-n3:7100 list_tablets ysql.yugabyte.donation 2>/dev/null" \
        > "$OUT/yb-cluster3/placement-evidence.txt" 2>&1 || true
      podman exec yb-n1 bash -lc "ysqlsh -h yb-n1 -p 5433 -U yugabyte -d yugabyte -Atc 'SELECT host, count(*) FROM yb_servers() GROUP BY host'" \
        >> "$OUT/yb-cluster3/placement-evidence.txt" 2>&1 || true
      [[ "$REDIS_UP" == "yes" ]] && bash "$REPO/infra/redis.sh" down
      bash "$REPO/infra/yb-cluster3.sh" down
      ;;
    *) warn "unknown topology $topo -- skipped" ;;
  esac
done

{
  echo "  pass_finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  if ((${#FAILED[@]})); then
    echo "  failed_cells:"
    printf '    - %s\n' "${FAILED[@]}"
  else
    echo "  failed_cells: []"
  fi
} >> "$OUT/manifest.yaml"

log "generating report"
podman run --rm -v "$(hostpath "$STUDY_DIR"):/study" "$BENCH_IMAGE" \
  -cmd report -results "/study/results/$RUN_ID" -report-out "/study/reports/$RUN_ID.md" || warn "report generation failed"

log "inputs digest"
podman run --rm -v "$(hostpath "$STUDY_DIR"):/study" "$BENCH_IMAGE" \
  -cmd digest -results "/study/results/$RUN_ID" || warn "digest failed"

log "results in $OUT"
((${#FAILED[@]})) && warn "${#FAILED[@]} cell(s) failed: ${FAILED[*]}"
exit 0
