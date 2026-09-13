#!/usr/bin/env bash
# Run study 02 across designs and topologies.
#
#   ./run-study.sh --scale small
#   ./run-study.sh --topologies pg-single --designs c1_count_naive,c4_counter_guard --scale tiny
#   ./run-study.sh --tag      # tag the commit as run/02-ticket-booking/<run-id> first
#
# Each (topology, design) pair is one independent container run writing its own
# JSON. A matrix that fails on one cell leaves every other cell's result behind;
# failures are recorded in the manifest, with the database logs captured before
# teardown.
#
# Provenance: the repository commit, `git describe`, and whether the tree had
# uncommitted changes are captured ONCE, before anything runs, and passed into
# every cell and the manifest. With --tag, a clean tree is tagged so the run can
# always be checked out again; a dirty tree is never tagged.
set -uo pipefail

# Run from a private copy of this script. Bash reads a script incrementally
# while it runs, so editing the file mid-matrix can garble what it executes next
# (LESSONS_LEARNED: "Never edit a bash script that is currently running") -- and
# the rule was broken again while this study was being built. Re-executing a
# copy makes the original safe to edit at any time.
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
DURATION="5s"
WARMUP="2s"
CONNS="8"
LOAD_CONNS="4"
TRIALS="1"
TOPOLOGIES="pg-single,yb-single,yb-cluster3"
DESIGNS=""
PHASES="verify,explain,read,write,race,churn,holds"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
TAG="no"
EXTRA=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scale) SCALE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --warmup) WARMUP="$2"; shift 2 ;;
    --conns) CONNS="$2"; shift 2 ;;
    --load-conns) LOAD_CONNS="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --topologies) TOPOLOGIES="$2"; shift 2 ;;
    --designs) DESIGNS="$2"; shift 2 ;;
    --phases) PHASES="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --tag) TAG="yes"; shift ;;
    --extra) EXTRA="$2"; shift 2 ;;
    -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
    *) die "unknown flag $1" ;;
  esac
done

ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
OUT="$STUDY_DIR/results/$RUN_ID"
mkdir -p "$OUT"

ALL_DESIGNS="p1_precreated_lock_first p2_precreated_skip_locked p3_precreated_cas p4_precreated_counter c1_count_naive c2_count_serializable c3_count_lock_event c4_counter_guard c5_seat_unique r1_inventory_row r2_inventory_buckets r3_seat_pool h0_hold_naive_confirm h1_hold_checked_confirm"

# ---------------------------------------------------------------------------
# Repository version. "Dirty" means uncommitted changes in the code that can
# change THIS study's results: the shared platform, the infra scripts and the
# study itself. Results, reports and analyses are excluded -- they are what a run
# produces -- and so are other studies and the docs, which several analysts may
# be editing at the same time without affecting a single number here.
# ---------------------------------------------------------------------------
REPO_COMMIT="$(git -C "$REPO" rev-parse HEAD 2>/dev/null || echo unknown)"
CODE_PATHS=(platform infra "studies/$STUDY_ID" .containerignore ':!studies/*/results' ':!studies/*/reports')
REPO_DIRTY="false"
if [[ -n "$(git -C "$REPO" status --porcelain -- "${CODE_PATHS[@]}" 2>/dev/null)" ]]; then
  REPO_DIRTY="true"
fi
REPO_DESCRIBE="$(git -C "$REPO" describe --tags --always 2>/dev/null || echo "$REPO_COMMIT")"
[[ "$REPO_DIRTY" == "true" ]] && REPO_DESCRIBE="${REPO_DESCRIBE}-dirty"
RUN_TAG=""
if [[ "$TAG" == "yes" ]]; then
  if [[ "$REPO_DIRTY" == "true" ]]; then
    warn "study code has uncommitted changes; NOT tagging (commit first)"
    git -C "$REPO" status --porcelain -- "${CODE_PATHS[@]}" | head -20 >&2
  else
    RUN_TAG="run/${STUDY_ID}/${RUN_ID}"
    if git -C "$REPO" rev-parse -q --verify "refs/tags/$RUN_TAG" >/dev/null; then
      log "tag $RUN_TAG already exists (extending a run); leaving it in place"
    else
      git -C "$REPO" tag -a "$RUN_TAG" -m "Study ${STUDY_ID}: code that produced run ${RUN_ID}" \
        && log "tagged $REPO_COMMIT as $RUN_TAG"
    fi
  fi
fi
[[ "$REPO_DIRTY" == "true" ]] && warn "running from a dirty tree: results will say so"

# Refuse to trample another run: the infra scripts reuse container names.
BUSY="$(podman ps --format '{{.Names}}' | grep -E '^(pg-single|yb-single|yb-n[123])$' || true)"
[[ -n "$BUSY" ]] && die "database containers already running ($BUSY) -- another run may be in progress"

need_podman
log "building benchmark image $BENCH_IMAGE (context: repository root)"
podman build -q -t "$BENCH_IMAGE" -f "$(hostpath "$STUDY_DIR/Containerfile")" "$(hostpath "$REPO")" >/dev/null \
  || die "image build failed"

if [[ -f "$OUT/manifest.yaml" ]]; then
  warn "run id $RUN_ID already has results; appending a new pass to its manifest"
  echo "" >> "$OUT/manifest.yaml"
fi
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
  echo "  duration_per_measurement: $DURATION"
  echo "  warmup: $WARMUP"
  echo "  client_conns: $CONNS"
  echo "  trials: $TRIALS"
  echo "  phases: $PHASES"
  echo "  extra_flags: \"$EXTRA\""
  echo "  yugabyte_harness_flags: \"$YB_HARNESS_FLAGS\""
  echo "  topologies: $TOPOLOGIES"
  echo "  designs: ${DESIGNS:-all}"
  echo "  images:"
  echo "    postgres: $PG_IMAGE"
  echo "    yugabyte: $YB_IMAGE"
  echo "  yb_extra_tserver_flags: $YB_EXTRA_TSERVER_FLAGS"
  echo "  bench_image: $BENCH_IMAGE"
  echo "  bench_image_id: $(podman image inspect "$BENCH_IMAGE" --format '{{.Id}}' 2>/dev/null | cut -c1-19)"
  echo "  budget_per_db_node: cpus=$DB_CPUS memory=$DB_MEMORY"
  echo "  budget_client: cpus=$CLIENT_CPUS memory=$CLIENT_MEMORY"
  echo "  podman: $(podman --version)"
} >> "$OUT/manifest.yaml"

FAILED=()

run_cell() {
  local topo="$1" engine="$2" dsn="$3" design="$4"
  # Engine-specific harness flags come first so --extra can still override them.
  local ENGINE_FLAGS=""
  [[ "$engine" == "yugabyte" ]] && ENGINE_FLAGS="$YB_HARNESS_FLAGS"
  local jdir="$OUT/$topo" pdir="$OUT/$topo/plans" ldir="$OUT/$topo/logs"
  mkdir -p "$jdir" "$pdir" "$ldir"
  local names c
  case "$topo" in
    pg-single) names="pg-single" ;;
    yb-single) names="yb-single" ;;
    yb-cluster3) names="yb-n1 yb-n2 yb-n3" ;;
  esac
  # The database containers' CPU throttling over the whole cell: the server-side
  # half of the quota hazard the report shows for the client.
  for c in $names; do
    { echo "[$c] before $(date -u +%H:%M:%S)"; podman exec "$c" cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
  done
  log "[$topo] $design"
  # shellcheck disable=SC2086
  podman run --rm --network "$NETWORK" \
    --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$jdir"):/results" -v "$(hostpath "$pdir"):/plans" \
    "$BENCH_IMAGE" \
      -cmd full -engine "$engine" -topology "$topo" -design "$design" \
      -scale "$SCALE" -conns "$CONNS" -load-conns "$LOAD_CONNS" \
      -duration "$DURATION" -warmup "$WARMUP" -trials "$TRIALS" -phases "$PHASES" \
      -environment "$ENVIRONMENT" -run-id "$RUN_ID" \
      -repo-commit "$REPO_COMMIT" -repo-describe "$REPO_DESCRIBE" -repo-dirty="$REPO_DIRTY" -run-tag "$RUN_TAG" \
      $ENGINE_FLAGS $EXTRA \
      -dsn "$dsn" -out "/results/$design.json" -explain-out "/plans/$design.txt" \
    2>&1 | tee "$ldir/$design.log"
  local status=${PIPESTATUS[0]}
  for c in $names; do
    { echo "[$c] after $(date -u +%H:%M:%S)"; podman exec "$c" cat /sys/fs/cgroup/cpu.stat 2>&1; } >> "$ldir/$design.dbcpu.txt"
  done
  if [[ $status -ne 0 ]]; then
    warn "[$topo] $design FAILED -- continuing with the rest of the matrix"
    for c in $names; do
      podman logs --tail 300 "$c" > "$ldir/$design.$c.log" 2>&1 || true
      # yugabyted's stdout holds none of the server logs. The YSQL (PostgreSQL
      # layer) and tserver logs live under the data directory; without them a
      # "database system is shutting down" in study 02's dev run could not be
      # diagnosed after teardown.
      if [[ "$engine" == "yugabyte" ]]; then
        podman exec "$c" bash -c 'for f in /home/yugabyte/yb_data/logs/tserver/postgresql-*.log /home/yugabyte/yb_data/logs/tserver/yb-tserver.*; do [ -f "$f" ] && { echo "=== $f"; tail -n 400 "$f"; }; done' \
          >> "$ldir/$design.$c.log" 2>&1 || true
      fi
    done
    FAILED+=("$topo/$design")
  fi
}

designs_for() {
  local list="${DESIGNS//,/ }"
  echo "${list:-$ALL_DESIGNS}"
}

for topo in ${TOPOLOGIES//,/ }; do
  case "$topo" in
    pg-single)
      bash "$REPO/infra/pg-single.sh" up || { warn "pg-single failed to start"; FAILED+=("pg-single/*"); continue; }
      record_topology "$OUT/pg-single/topology.yaml" pg-single
      for d in $(designs_for); do
        run_cell pg-single postgres "postgres://bench:bench@pg-single:5432/bench?sslmode=disable" "$d"
      done
      bash "$REPO/infra/pg-single.sh" down
      ;;
    yb-single)
      bash "$REPO/infra/yb-single.sh" up || { warn "yb-single failed to start"; FAILED+=("yb-single/*"); continue; }
      record_topology "$OUT/yb-single/topology.yaml" yb-single
      for d in $(designs_for); do
        run_cell yb-single yugabyte "postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable" "$d"
      done
      bash "$REPO/infra/yb-single.sh" down
      ;;
    yb-cluster3)
      bash "$REPO/infra/yb-cluster3.sh" up || { warn "yb-cluster3 failed to start"; FAILED+=("yb-cluster3/*"); continue; }
      record_topology "$OUT/yb-cluster3/topology.yaml" yb-n1 yb-n2 yb-n3
      for d in $(designs_for); do
        run_cell yb-cluster3 yugabyte "postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable" "$d"
      done
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

log "results in $OUT"
((${#FAILED[@]})) && warn "${#FAILED[@]} cell(s) failed: ${FAILED[*]}"
exit 0
