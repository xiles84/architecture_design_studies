#!/usr/bin/env bash
# Run study 01 across every design and every topology.
#
#   ./run-study.sh --scale medium --duration 10s
#
# Each (topology, design) pair is one independent invocation of the benchmark
# container writing its own JSON file. That matters: a three-hour matrix that
# aborts on the fifth of twenty-one cells should still leave four usable results
# behind, not an empty directory. Failures are recorded and the matrix continues.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"

# lib.sh turns on `set -e`. Turn it back off: this script is deliberately
# resilient -- a failed cell must be recorded and the matrix must continue.
# With -e active, the first non-zero podman exit would kill the whole run
# before the failure-handling code below ever executed.
set +e

SCALE="medium"
DURATION="10s"
WARMUP="3s"
CONNS="8"
LOAD_CONNS="4"
TOPOLOGIES="pg-single,yb-single,yb-cluster3"
DESIGNS=""
# Child-originating writes first, then the two that originate at the PARENT.
# Leaving the parent ones out hides costs that fall specifically on the embedded
# and rollup designs -- see sql/*/writes.sql.
WRITE_OPS="insert,update,delete,update_person,delete_person"
TRIALS="1"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
KEEP_UP="no"

# Extra harness flags passed straight to the benchmark container, e.g.
# --extra "-max-per-person 20 -charities 1000". Used by run-regimes.sh.
EXTRA=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scale) SCALE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --warmup) WARMUP="$2"; shift 2 ;;
    --conns) CONNS="$2"; shift 2 ;;
    --load-conns) LOAD_CONNS="$2"; shift 2 ;;
    --topologies) TOPOLOGIES="$2"; shift 2 ;;
    --designs) DESIGNS="$2"; shift 2 ;;
    --write-ops) WRITE_OPS="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --keep-up) KEEP_UP="yes"; shift ;;
    --extra) EXTRA="$2"; shift 2 ;;
    -h|--help) sed -n '2,12p' "$0"; exit 0 ;;
    *) die "unknown flag $1" ;;
  esac
done

ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
OUT="$STUDY_DIR/results/$RUN_ID"
mkdir -p "$OUT"

ALL_PG_DESIGNS="d1_normalized_minimal d2_normalized_indexed d3_flattened_fk d8_flattened_nofk d4_rollup_trigger d5_rollup_app d6_embedded_jsonb d9_embedded_hybrid"
ALL_YB_DESIGNS="$ALL_PG_DESIGNS d7_yb_child_colocated"

need_podman
log "building benchmark image"
podman build -q -t "$BENCH_IMAGE" \
  -f "$(hostpath "$STUDY_DIR/Containerfile")" "$(hostpath "$STUDY_DIR")" >/dev/null \
  || die "image build failed"

# ---------------------------------------------------------------------------
# Manifest: what this run is, before any of it happens.
#
# Appended, never truncated. A run id can legitimately be filled in more than one
# pass -- a design added after the first sweep, a failed cell re-run -- and each
# pass records its own settings and its own image id. Overwriting would erase the
# record of how the earlier cells were produced, which is precisely the context a
# later reader needs to know whether the cells are comparable.
# ---------------------------------------------------------------------------
if [[ -f "$OUT/manifest.yaml" ]]; then
  warn "run id $RUN_ID already has results; appending a new pass to its manifest"
  echo "" >> "$OUT/manifest.yaml"
fi
{
  echo "- pass_started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "  run_id: $RUN_ID"
  echo "  environment: $ENVIRONMENT"
  echo "  scale: $SCALE"
  echo "  duration_per_op: $DURATION"
  echo "  warmup_per_op: $WARMUP"
  echo "  client_conns: $CONNS"
  echo "  load_conns: $LOAD_CONNS"
  echo "  write_ops: $WRITE_OPS"
  echo "  trials_per_cell: $TRIALS"
  echo "  topologies: $TOPOLOGIES"
  echo "  designs: ${DESIGNS:-all}"
  echo "  images:"
  echo "    postgres: $PG_IMAGE"
  echo "    yugabyte: $YB_IMAGE"
  echo "  bench_image_id: $(podman image inspect "$BENCH_IMAGE" --format '{{.Id}}' 2>/dev/null | cut -c1-19)"
  echo "  budget_per_db_node: cpus=$DB_CPUS memory=$DB_MEMORY"
  echo "  budget_client: cpus=$CLIENT_CPUS memory=$CLIENT_MEMORY"
  echo "  podman: $(podman --version)"
} >> "$OUT/manifest.yaml"

FAILED=()

run_cell() {
  local topo="$1" engine="$2" dsn="$3" design="$4"
  local jdir="$OUT/$topo" pdir="$OUT/$topo/plans"
  mkdir -p "$jdir" "$pdir"

  log "[$topo] $design"
  podman run --rm \
    --network "$NETWORK" \
    --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$jdir"):/results" \
    -v "$(hostpath "$pdir"):/plans" \
    "$BENCH_IMAGE" \
      -cmd full \
      -engine "$engine" \
      -topology "$topo" \
      -design "$design" \
      -scale "$SCALE" \
      -conns "$CONNS" \
      -load-conns "$LOAD_CONNS" \
      -duration "$DURATION" \
      -warmup "$WARMUP" \
      -trials "$TRIALS" \
      -write-ops "$WRITE_OPS" \
      -environment "$ENVIRONMENT" \
      -run-id "$RUN_ID" \
      $EXTRA \
      -dsn "$dsn" \
      -out "/results/$design.json" \
      -explain-out "/plans/$design.txt"

  if [[ $? -ne 0 ]]; then
    warn "[$topo] $design FAILED -- continuing with the rest of the matrix"
    # Capture the database side BEFORE the topology is torn down. The first
    # failure on the 3-node cluster lost its logs with the containers, leaving
    # only the client-side error to diagnose from.
    local c names
    case "$topo" in
      pg-single) names="pg-single" ;;
      yb-single) names="yb-single" ;;
      yb-cluster3) names="yb-n1 yb-n2 yb-n3" ;;
    esac
    for c in $names; do
      podman logs --tail 300 "$c" > "$jdir/$design.$c.log" 2>&1 || true
    done
    FAILED+=("$topo/$design")
  fi
}

designs_for() {
  if [[ -n "$DESIGNS" ]]; then
    echo "${DESIGNS//,/ }"
  elif [[ "$1" == "postgres" ]]; then
    echo "$ALL_PG_DESIGNS"
  else
    echo "$ALL_YB_DESIGNS"
  fi
}

for topo in ${TOPOLOGIES//,/ }; do
  case "$topo" in
    pg-single)
      bash "$REPO/infra/pg-single.sh" up || { warn "pg-single failed to start"; continue; }
      record_topology "$OUT/pg-single/topology.yaml" pg-single
      for d in $(designs_for postgres); do
        run_cell pg-single postgres "postgres://bench:bench@pg-single:5432/bench?sslmode=disable" "$d"
      done
      [[ "$KEEP_UP" == "yes" ]] || bash "$REPO/infra/pg-single.sh" down
      ;;
    yb-single)
      bash "$REPO/infra/yb-single.sh" up || { warn "yb-single failed to start"; continue; }
      record_topology "$OUT/yb-single/topology.yaml" yb-single
      for d in $(designs_for yugabyte); do
        run_cell yb-single yugabyte "postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable" "$d"
      done
      [[ "$KEEP_UP" == "yes" ]] || bash "$REPO/infra/yb-single.sh" down
      ;;
    yb-cluster3)
      bash "$REPO/infra/yb-cluster3.sh" up || { warn "yb-cluster3 failed to start"; continue; }
      record_topology "$OUT/yb-cluster3/topology.yaml" yb-n1 yb-n2 yb-n3
      for d in $(designs_for yugabyte); do
        run_cell yb-cluster3 yugabyte "postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable" "$d"
      done
      [[ "$KEEP_UP" == "yes" ]] || bash "$REPO/infra/yb-cluster3.sh" down
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

log "results in $OUT"
((${#FAILED[@]})) && warn "${#FAILED[@]} cell(s) failed: ${FAILED[*]}"
exit 0
