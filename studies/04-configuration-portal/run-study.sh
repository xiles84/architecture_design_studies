#!/usr/bin/env bash
# Run study 04 across designs and topologies.
#
#   ./run-study.sh --scale tiny --topologies pg-single --designs n1_rows_indexed
#   ./run-study.sh --tag --scale small --topologies pg-single
#
# Each (topology, design) pair is one independent container run writing its own
# JSON. A matrix that fails on one cell leaves every other cell's result behind;
# failures are recorded in the manifest, with the database logs captured before
# teardown.
#
# Provenance: the repository commit, `git describe`, and whether the study's code
# had uncommitted changes are captured ONCE, before anything runs, and passed into
# every cell and the manifest. With --tag, a clean tree is tagged so the run can
# always be checked out again; a dirty tree is never tagged.
set -uo pipefail

# Run from a private copy of this script. Bash reads a script incrementally while
# it runs, so editing the file mid-matrix can garble what it executes next
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
DURATION="5s"
WARMUP="2s"
CONNS="8"
LOAD_CONNS="4"
TRIALS="1"
ENTRIES_PER_IP=""
FLEET=""
REGIME="small"
CARD_MODE="exact"
CADENCE="none"
BURST="1"
CONCURRENCY="16"
TOPOLOGIES="pg-single,yb-single,yb-cluster3"
DESIGNS=""
PHASES="verify,explain,read,write,contention"
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
    --load-conns) LOAD_CONNS="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --entries-per-ip) ENTRIES_PER_IP="$2"; shift 2 ;;
    --fleet) FLEET="$2"; shift 2 ;;
    --value-regime) REGIME="$2"; shift 2 ;;
    --cardinality-mode) CARD_MODE="$2"; shift 2 ;;
    --cadence) CADENCE="$2"; shift 2 ;;
    --burst-size) BURST="$2"; shift 2 ;;
    --concurrency) CONCURRENCY="$2"; shift 2 ;;
    --topologies) TOPOLOGIES="$2"; shift 2 ;;
    --designs) DESIGNS="$2"; shift 2 ;;
    --phases) PHASES="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --tag) TAG="yes"; shift ;;
    --extra) EXTRA="$2"; shift 2 ;;
    --order-seed) ORDER_SEED="$2"; shift 2 ;;
    -h|--help) sed -n '2,12p' "$0"; exit 0 ;;
    *) die "unknown flag $1" ;;
  esac
done

ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
[[ -z "$ORDER_SEED" ]] && ORDER_SEED="$(printf '%s' "$RUN_ID" | cksum | cut -d' ' -f1)"
OUT="$STUDY_DIR/results/$RUN_ID"
mkdir -p "$OUT"

ALL_DESIGNS="n0_rows_unindexed n1_rows_indexed n2_rows_rolldown n3_rollup_trigger n4_rollup_app d1_doc_row c1_optimistic_version c2_pessimistic_lock x1_lost_update_control x2_rollup_drift_control"

# ---------------------------------------------------------------------------
# Repository version. "Dirty" means uncommitted changes in the code that can
# change THIS study's results: the shared platform, the infra scripts and the
# study itself. Results, reports and analyses are excluded -- they are what a run
# produces -- and so are other studies and the docs, which several analysts may be
# editing at the same time without affecting a single number here.
# ---------------------------------------------------------------------------
GIT=(git -C "$(hostpath "$REPO")")
REPO_COMMIT="$("${GIT[@]}" rev-parse HEAD 2>/dev/null)" \
  || die "cannot read the repository commit -- refusing to produce results that cannot be traced to code"
# Dirty means uncommitted changes to code that can change results. Generated
# results and reports are what a run PRODUCES, not code, so they are filtered
# out explicitly: git's `:!` pathspec exclusion does not suppress untracked files
# under those directories, which made every run after the first report a dirty
# tree and silently skip --tag (LESSONS_LEARNED, Study 04 v2 controls 2026-09-23).
CODE_PATHS=(platform infra "studies/$STUDY_ID" .containerignore)
STATUS="$("${GIT[@]}" status --porcelain -- "${CODE_PATHS[@]}" | grep -vE 'studies/[^/]+/(results|reports)/' || true)" \
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
# (infra/lib.sh). Taken before anything builds an image.
run_lock_acquire "study ${STUDY_ID} run-study.sh run ${RUN_ID}"

BUSY="$(podman ps --format '{{.Names}}' | grep -E '^(pg-single|yb-single|yb-n[123])$' || true)"
[[ -n "$BUSY" ]] && die "database containers already running ($BUSY) -- another run may be in progress"

need_podman
log "building benchmark image $BENCH_IMAGE (context: repository root)"
# winpath, not hostpath, for the build. A bind-mount source tolerates /mnt/c/...
# because the engine resolves it inside the machine, but the *build context* is
# resolved by the podman client on Windows, which reads a leading "/" as the root
# of the current drive: /mnt/c/extra/... becomes C:\mnt\c\extra\... and the build
# fails with "context must be a directory". winpath() is a no-op under Git Bash,
# where hostpath() already returned the C:/ form, so this is correct on both.
podman build -q -t "$BENCH_IMAGE" -f "$(winpath "$STUDY_DIR/Containerfile")" "$(winpath "$REPO")" >/dev/null \
  || die "image build failed"

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
  echo "  entries_per_ip: ${ENTRIES_PER_IP:-<scale default>}"
  echo "  fleet: ${FLEET:-<scale default>}"
  echo "  value_regime: $REGIME"
  echo "  cardinality_mode: $CARD_MODE"
  echo "  cadence: $CADENCE"
  echo "  burst_size: $BURST"
  echo "  duration_per_measurement: $DURATION"
  echo "  warmup: $WARMUP"
  echo "  client_conns: $CONNS"
  echo "  trials: $TRIALS"
  echo "  contention_writers: $CONCURRENCY"
  echo "  phases: $PHASES"
  echo "  extra_flags: \"$EXTRA\""
  echo "  postgres_harness_flags: \"$PG_HARNESS_FLAGS\""
  echo "  yugabyte_harness_flags: \"$YB_HARNESS_FLAGS\""
  echo "  yugabyte_only_designs: $YB_ONLY_DESIGNS"
  echo "  order_seed: $ORDER_SEED"
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
  echo "  podman_client: $ADS_PODMAN"
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
  log "[$topo] $design"
  # shellcheck disable=SC2086
  podman run --rm --network "$NETWORK" \
    --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$jdir"):/results" -v "$(hostpath "$pdir"):/plans" \
    "$BENCH_IMAGE" \
      -cmd full -engine "$engine" -topology "$topo" -design "$design" \
      -scale "$SCALE" -conns "$CONNS" -load-conns "$LOAD_CONNS" \
      -duration "$DURATION" -warmup "$WARMUP" -trials "$TRIALS" -phases "$PHASES" \
      -value-regime "$REGIME" -cardinality-mode "$CARD_MODE" \
      -cadence "$CADENCE" -burst-size "$BURST" -concurrency "$CONCURRENCY" \
      ${ENTRIES_PER_IP:+-entries-per-ip "$ENTRIES_PER_IP"} \
      ${FLEET:+-fleet "$FLEET"} \
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
      if [[ "$engine" == "yugabyte" ]]; then
        podman exec "$c" bash -c 'for f in /home/yugabyte/yb_data/logs/tserver/postgresql-*.log /home/yugabyte/yb_data/logs/tserver/yb-tserver.*; do [ -f "$f" ] && { echo "=== $f"; tail -n 400 "$f"; }; done' \
          >> "$ldir/$design.$c.log" 2>&1 || true
      fi
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

for topo in ${TOPOLOGIES//,/ }; do
  ORDER="$(designs_for "$topo")"
  echo "  design_order_${topo}: $ORDER" >> "$OUT/manifest.yaml"
  case "$topo" in
    pg-single)
      bash "$REPO/infra/pg-single.sh" up || { warn "pg-single failed to start"; FAILED+=("pg-single/*"); continue; }
      record_topology "$OUT/pg-single/topology.yaml" pg-single
      for d in $ORDER; do
        run_cell pg-single postgres "postgres://bench:bench@pg-single:5432/bench?sslmode=disable" "$d"
      done
      bash "$REPO/infra/pg-single.sh" down
      ;;
    yb-single)
      bash "$REPO/infra/yb-single.sh" up || { warn "yb-single failed to start"; FAILED+=("yb-single/*"); continue; }
      record_topology "$OUT/yb-single/topology.yaml" yb-single
      for d in $ORDER; do
        run_cell yb-single yugabyte "postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable" "$d"
      done
      bash "$REPO/infra/yb-single.sh" down
      ;;
    yb-cluster3)
      bash "$REPO/infra/yb-cluster3.sh" up || { warn "yb-cluster3 failed to start"; FAILED+=("yb-cluster3/*"); continue; }
      record_topology "$OUT/yb-cluster3/topology.yaml" yb-n1 yb-n2 yb-n3
      for d in $ORDER; do
        run_cell yb-cluster3 yugabyte "postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable,postgres://yugabyte@yb-n2:5433/yugabyte?sslmode=disable,postgres://yugabyte@yb-n3:5433/yugabyte?sslmode=disable" "$d"
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

log "inputs digest"
podman run --rm -v "$(hostpath "$STUDY_DIR"):/study" "$BENCH_IMAGE" \
  -cmd digest -results "/study/results/$RUN_ID" || warn "digest failed"

log "results in $OUT"
((${#FAILED[@]})) && warn "${#FAILED[@]} cell(s) failed: ${FAILED[*]}"
exit 0
