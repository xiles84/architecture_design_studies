#!/usr/bin/env bash
# Experiment D — reads and writes at the same time.
#
#   ./run-mixed.sh --topology pg-single
#
# Every other experiment in this study measures reads alone, then writes alone.
# That gives clean attribution but it cannot see the costs that only exist when
# the two overlap:
#
#   * MVCC bloat, which accumulates WHILE readers are scanning
#   * contention on rows that writers are constantly updating
#   * autovacuum (PostgreSQL) and RocksDB compaction (YugabyteDB), both triggered
#     by writes, both stealing CPU from reads
#   * buffer-cache and index-page competition
#
# Every one of those penalises designs that buy read speed with redundancy --
# which is most of the designs here. Isolated measurement therefore flatters
# exactly the designs under scrutiny, and this experiment is how that bias gets
# checked rather than assumed away.
#
# Splits are READERS:WRITERS. The first is all-readers and is the baseline: it
# runs the same weighted query mix over the same data in the same process, so the
# only thing that changes across splits is the presence of writers. The result to
# read is the RATIO between splits, not the absolute throughput.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"

# lib.sh turns on `set -e`. Turn it back off: a failed cell must be recorded and
# the run must continue.
set +e

TOPOLOGY="pg-single"
SCALE="small"
DURATION="15s"
WARMUP="5s"
SPLITS="8:0,7:1,6:2,4:4"
DESIGNS="d2_normalized_indexed,d3_flattened_fk,d4_rollup_trigger,d6_embedded_jsonb,d9_embedded_hybrid"
RUN_ID="$(date -u +%Y%m%d)-mixed"

# Extra harness flags passed straight to the benchmark container, e.g.
# --extra "-max-per-person 20 -charities 1000". Used by run-regimes.sh.
EXTRA=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --topology) TOPOLOGY="$2"; shift 2 ;;
    --scale) SCALE="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --splits) SPLITS="$2"; shift 2 ;;
    --designs) DESIGNS="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --extra) EXTRA="$2"; shift 2 ;;
    -h|--help) sed -n '2,22p' "$0"; exit 0 ;;
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
  echo "experiment: D — reads and writes running concurrently"
  echo "started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "environment: $ENVIRONMENT"
  echo "topology: $TOPOLOGY"
  echo "scale: $SCALE"
  echo "duration_per_split: $DURATION"
  echo "splits_readers_writers: $SPLITS"
  echo "designs: $DESIGNS"
  echo "note: the harness reloads the dataset before every split, so no split inherits"
  echo "      the bloat created by the one before it"
} > "$STUDY_DIR/results/$RUN_ID/manifest.yaml"

for design in ${DESIGNS//,/ }; do
  log "[$TOPOLOGY] $design"
  podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
    -v "$(hostpath "$OUT"):/results" \
    "$BENCH_IMAGE" \
      -cmd mixed -engine "$ENGINE" -topology "$TOPOLOGY" -design "$design" -scale "$SCALE" \
      -duration "$DURATION" -warmup "$WARMUP" -mix "$SPLITS" -load-conns 4 \
      -environment "$ENVIRONMENT" -run-id "$RUN_ID" $EXTRA \
      -dsn "$DSN" -out "/results/$design.json" \
    || warn "[$TOPOLOGY] $design failed — continuing"
done

echo "finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$STUDY_DIR/results/$RUN_ID/manifest.yaml"
log "results in $OUT"
