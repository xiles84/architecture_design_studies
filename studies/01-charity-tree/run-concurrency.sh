#!/usr/bin/env bash
# Experiment C — concurrency control on the hot rollup row.
#
#   ./run-concurrency.sh --topology pg-single
#
# D4 and D5 store identical aggregates and read them identically. They differ
# only in HOW the parent rows get updated, which makes them the right place to
# ask what optimistic and pessimistic concurrency control actually cost.
#
# The contention is structural, not contrived: there are ten charities and every
# donation in the system updates one of their rows. The largest charity holds 30%
# of all donors, so under load its row is the single point every writer must pass.
#
# Two things are measured per cell, and the second matters more:
#
#   1. throughput and retry count
#   2. whether the aggregates STILL AGREE with the donation rows afterwards
#
# A strategy that is faster and silently loses increments has not won anything —
# it has published a wrong total that nothing in the system would ever notice.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"

# lib.sh turns on `set -e`. Turn it back off: this script is deliberately
# resilient -- a failed cell must be recorded and the matrix must continue.
# With -e active, the first non-zero podman exit would kill the whole run
# before the failure-handling code below ever executed.
set +e
run_lock_acquire "study 01 run-concurrency.sh $*"

TOPOLOGY="pg-single"
SCALE="small"
DURATION="10s"
WARMUP="3s"
CONNS_LIST="1,4,8,16,32"
STRATEGIES="blind,forupdate,optimistic"
ISOLATIONS="rc,rr,ser"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-concurrency"

# Extra harness flags passed straight to the benchmark container, e.g.
# --extra "-max-per-person 20 -charities 1000". Used by run-regimes.sh.
EXTRA=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --topology) TOPOLOGY="$2"; shift 2 ;;
    --scale) SCALE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --conns-list) CONNS_LIST="$2"; shift 2 ;;
    --strategies) STRATEGIES="$2"; shift 2 ;;
    --isolations) ISOLATIONS="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --extra) EXTRA="$2"; shift 2 ;;
    -h|--help) sed -n '2,25p' "$0"; exit 0 ;;
    *) die "unknown flag $1" ;;
  esac
done

case "$TOPOLOGY" in
  pg-single)   ENGINE=postgres; DSN="postgres://bench:bench@pg-single:5432/bench?sslmode=disable"; UP="$REPO/infra/pg-single.sh" ;;
  yb-single)   ENGINE=yugabyte; DSN="postgres://yugabyte@yb-single:5433/yugabyte?sslmode=disable"; UP="$REPO/infra/yb-single.sh" ;;
  yb-cluster3) ENGINE=yugabyte; DSN="postgres://yugabyte@yb-n1:5433/yugabyte?sslmode=disable";     UP="$REPO/infra/yb-cluster3.sh" ;;
  *) die "unknown topology $TOPOLOGY" ;;
esac

ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
OUT="$STUDY_DIR/results/$RUN_ID/$TOPOLOGY"
mkdir -p "$OUT"

need_podman
log "building benchmark image"
podman build -q -t "$BENCH_IMAGE" \
  -f "$(hostpath "$STUDY_DIR/Containerfile")" "$(hostpath "$STUDY_DIR")" >/dev/null || die "image build failed"

bash "$UP" up || die "$TOPOLOGY failed to start"
record_topology "$OUT/topology.yaml" pg-single yb-single yb-n1 yb-n2 yb-n3

{
  echo "run_id: $RUN_ID"
  echo "experiment: C — concurrency control on the hot rollup row"
  echo "started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "environment: $ENVIRONMENT"
  echo "topology: $TOPOLOGY"
  echo "scale: $SCALE"
  echo "duration_per_cell: $DURATION"
  echo "conns_list: $CONNS_LIST"
  echo "strategies: $STRATEGIES"
  echo "isolations: $ISOLATIONS"
} > "$STUDY_DIR/results/$RUN_ID/manifest.yaml"

# cell <design> <strategy> <isolation> <conns> <label>
#
# Each cell reloads from scratch. Without that, the second cell would be writing
# into a table the first cell already grew, and the comparison would drift.
cell() {
  local design="$1" strategy="$2" isolation="$3" conns="$4" label="$5"
  log "[$TOPOLOGY] $label"

  podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    "$BENCH_IMAGE" -cmd load -engine "$ENGINE" -design "$design" -scale "$SCALE" \
    -load-conns 4 $EXTRA -dsn "$DSN" >/dev/null 2>&1 \
    || { warn "load failed for $label"; return 1; }

  podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$OUT"):/results" \
    "$BENCH_IMAGE" \
      -cmd write -engine "$ENGINE" -topology "$TOPOLOGY" -design "$design" -scale "$SCALE" \
      -conns "$conns" -duration "$DURATION" -warmup "$WARMUP" \
      -strategy "$strategy" -isolation "$isolation" -write-ops insert \
      -environment "$ENVIRONMENT" -run-id "$RUN_ID" $EXTRA \
      -dsn "$DSN" -out "/results/$label.json" \
    || warn "write benchmark failed for $label"
}

# --- Part 1: the contention curve -------------------------------------------
# How each strategy degrades as writers are added. A strategy that wins at 4
# connections and collapses at 32 has not solved the problem.
for c in ${CONNS_LIST//,/ }; do
  cell d4_rollup_trigger blind rc "$c" "c${c}_d4_trigger_rc"
  for s in ${STRATEGIES//,/ }; do
    cell d5_rollup_app "$s" rc "$c" "c${c}_d5_${s}_rc"
  done
  # The no-rollup baseline: what the insert would cost if nobody maintained an
  # aggregate at all. Everything above is measured against this ceiling.
  cell d3_flattened_fk blind rc "$c" "c${c}_d3_norollup_rc"
done

# --- Part 2: isolation levels at fixed concurrency ---------------------------
# Held at 8 writers so the only variable is the isolation level. Serializable is
# expected to cost retries rather than latency, which is why retries are counted.
for s in ${STRATEGIES//,/ }; do
  for i in ${ISOLATIONS//,/ }; do
    [[ "$i" == "rc" ]] && continue   # already covered by part 1
    cell d5_rollup_app "$s" "$i" 8 "c8_d5_${s}_${i}"
  done
done

echo "finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$STUDY_DIR/results/$RUN_ID/manifest.yaml"
log "results in $OUT"
