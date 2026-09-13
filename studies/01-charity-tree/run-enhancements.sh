#!/usr/bin/env bash
# Study 01 v3 follow-ups. Entrypoint: run-study.sh --suite enhancements --tag.
# Every cell has its own fresh load, correctness gate, JSON, plans and resource logs.
set -euo pipefail
STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
[[ "${1:-}" == enhancements ]] && shift
EXPERIMENTS="reads,mechanisms,growth,contention,exceptions,deployment"
TRIALS=5
DURATION=3s
WARMUP=1s
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-v3"
TAG=no
while [[ $# -gt 0 ]]; do
  case "$1" in
    --experiments) EXPERIMENTS="$2"; shift 2 ;;
    --trials) TRIALS="$2"; shift 2 ;;
    --duration) DURATION="$2"; shift 2 ;;
    --warmup) WARMUP="$2"; shift 2 ;;
    --run-id) RUN_ID="$2"; shift 2 ;;
    --tag) TAG=yes; shift ;;
    *) die "unknown flag $1" ;;
  esac
done
[[ "$RUN_ID" =~ ^[a-zA-Z0-9_-]+$ ]] || die "unsafe run id"
[[ "$TRIALS" =~ ^[1-9][0-9]*$ ]] || die "trials must be positive"
for suite in ${EXPERIMENTS//,/ }; do
  case "$suite" in verify|reads|mechanisms|growth|contention|exceptions|deployment) ;; *) die "unknown experiment $suite" ;; esac
done
need_podman
run_lock_acquire "study-01 v3 $RUN_ID ($EXPERIMENTS)"
for c in pg-single yb-single yb-n1 yb-n2 yb-n3; do
  container_exists "$c" && die "container $c already exists; refusing to alter another run"
done
GIT=(git -C "$(hostpath "$REPO")")
COMMIT="$("${GIT[@]}" rev-parse HEAD)"
STATUS="$("${GIT[@]}" status --porcelain -- infra studies/01-charity-tree ':!studies/01-charity-tree/results' ':!studies/01-charity-tree/reports')"
[[ -z "$STATUS" ]] || die "study/infra code is dirty; commit before a reported run"
[[ "$TAG" == yes ]] || die "reported follow-ups require --tag"
OUT="$STUDY_DIR/results/$RUN_ID"
[[ ! -e "$OUT" ]] || die "run id already exists; create a new immutable run"
RUN_TAG="run/01-charity-tree/$RUN_ID"
"${GIT[@]}" tag -a "$RUN_TAG" "$COMMIT" -m "Study 01 v3 $RUN_ID; experiments $EXPERIMENTS" || die "run tag failed"
DESCRIBE="$("${GIT[@]}" describe --tags --always)"
mkdir -p "$OUT"
log "building the committed benchmark once"
podman build -q -t localhost/charitybench:3 -f "$(hostpath "$STUDY_DIR/Containerfile")" "$(hostpath "$STUDY_DIR")" > "$OUT/image-id.txt"
IMAGE="$(podman image inspect localhost/charitybench:3 --format '{{.Id}}')"
ENVIRONMENT="${BENCH_ENVIRONMENT:-host-zenbook-ux5406sa}"
export YB_EXTRA_TSERVER_FLAGS=yb_enable_read_committed_isolation=true
{
  echo "run_id: $RUN_ID"
  echo "environment: $ENVIRONMENT"
  echo "repo_commit: $COMMIT"
  echo "repo_describe: $DESCRIBE"
  echo "repo_dirty: false"
  echo "run_tag: $RUN_TAG"
  echo "image_id: $IMAGE"
  echo "started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "experiments: $EXPERIMENTS"
  echo "independent_trials: $TRIALS"
  echo "duration: $DURATION"
  echo "warmup: $WARMUP"
  echo "client_budget: 2 CPUs, 2 GiB"
  echo "cells:"
} > "$OUT/manifest.yaml"

# A persisted plan makes a killed runner's missing cells visible too.
PLAN="$OUT/cells.tsv"
: > "$PLAN"
cell() { printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "$@" >> "$PLAN"; }
enabled() { [[ ",$EXPERIMENTS," == *",$1,"* ]]; }
ordered() {
  local trial="$1"; shift
  if (( trial % 2 )); then printf '%s\n' "$@"; else local i; for ((i=$#;i>0;i--)); do printf '%s\n' "${!i}"; done; fi
}
for ((t=1;t<=TRIALS;t++)); do
  if enabled verify; then
    for topo in pg-single yb-single; do
      order=0
      for d in d11_copied_key d12_recency_index d13_recency_sql d17_sum_sql d14_sum_plain d15_sum_covering d16_sum_rollup; do
        order=$((order+1)); cell "$topo" standard "catalogue-gate" "$d" verify small "$t" "$order" "-preparation analyze"
      done
    done
  fi
  if enabled reads; then
    order=0
    for d in $(ordered "$t" d2_normalized_indexed d3_flattened_fk); do
      order=$((order+1)); cell pg-single standard vacuum-small "$d" reads small "$t" "$order" "-preparation vacuum -targeted"
      cell pg-single standard insert-small "$d" writes small "$t" "$order" "-write-ops insert"
      cell pg-single standard blended-small "$d" arrival small "$t" "$order" "-arrival-rate 500 -duration 15s"
    done
  fi
  if enabled mechanisms; then
    order=0
    for d in $(ordered "$t" d11_copied_key d12_recency_index d13_recency_sql d17_sum_sql d14_sum_plain d15_sum_covering); do
      order=$((order+1)); cell pg-single standard vacuum-small "$d" reads small "$t" "$order" "-preparation vacuum -targeted"
      cell pg-single standard insert-small "$d" writes small "$t" "$order" "-write-ops insert"
    done
    for d in $(ordered "$t" d14_sum_plain d15_sum_covering); do
      cell pg-single standard analyze-sum "$d" reads small "$t" 1 "-preparation analyze -queries q08_total_donated_charity -targeted"
    done
  fi
  if enabled growth; then
    for d in $(ordered "$t" d3_flattened_fk d6_embedded_jsonb); do
      cell pg-single standard medium-growth "$d" growth medium "$t" 1 "-cycles 3 -batch 1000 -duration 1s"
      cell pg-single standard long-history "$d" growth small "$t" 1 "-history-multiplier 8 -cycles 3 -batch 1000 -duration 1s"
      cell pg-single constrained memory-growth "$d" growth medium "$t" 1 "-history-multiplier 2 -cycles 3 -batch 1000 -duration 1s"
    done
  fi
  if enabled contention; then
    for hot in 0 0.9; do
      for rate in 250 2500; do
        order=0
        for d in $(ordered "$t" d4_rollup_trigger d5_rollup_app d16_sum_rollup); do
          order=$((order+1)); cell pg-single standard "hot${hot}-rate${rate}" "$d" arrival small "$t" "$order" "-strategy forupdate -read-workers 4 -write-workers 32 -arrival-rate $rate -hot-probability $hot -hot-donors 10 -duration 20s"
        done
      done
    done
    for d in $(ordered "$t" d4_rollup_trigger d16_sum_rollup); do
      for op in update delete; do cell pg-single standard "extrema-$op" "$d" writes small "$t" 1 "-write-ops $op -duration 15s"; done
    done
  fi
  if enabled exceptions; then
    order=0
    for d in $(ordered "$t" d3_flattened_fk d8_flattened_nofk); do
      order=$((order+1)); cell yb-single standard fk-insert "$d" writes small "$t" "$order" "-write-ops insert -duration 30s"
      cell yb-cluster3 standard fk-erasure "$d" writes medium "$t" "$order" "-write-ops delete_person -duration 120s"
    done
    for topo in yb-single yb-cluster3; do
      order=0
      for d in $(ordered "$t" d9_embedded_hybrid d10_embedded_hybrid_locked); do
        order=$((order+1)); cell "$topo" standard cache-race "$d" arrival small "$t" "$order" "-read-mix portal -arrival-rate 250 -hot-probability 0.9 -hot-donors 1 -write-workers 32 -duration 30s"
        for op in update delete; do cell "$topo" standard "cache-$op" "$d" writes small "$t" "$order" "-write-ops $op -duration 15s"; done
      done
    done
  fi
  if enabled deployment; then
    for topo in yb-single yb-cluster3; do
      budget=standard; [[ "$topo" == yb-single ]] && budget=large
      cell "$topo" "$budget" equal-total d3_flattened_fk reads small "$t" 1 "-targeted"
      cell "$topo" "$budget" equal-total d3_flattened_fk writes small "$t" 1 "-write-ops insert -duration 30s"
      cell "$topo" "$budget" equal-total-arrival d3_flattened_fk arrival small "$t" 1 "-arrival-rate 500 -duration 30s"
      if [[ "$topo" == yb-cluster3 ]]; then
        cell "$topo" "$budget" local-node-stop d3_flattened_fk arrival small "$t" 1 "-arrival-rate 100 -duration 30s"
      fi
    done
  fi
done

# Group startup by topology/budget while preserving alternating trial order.
LC_ALL=C sort -s -t $'\t' -k1,1 -k2,2 "$PLAN" > "$OUT/ordered-cells.tsv"
active=""; active_budget=""; failed=0; number=0
cleanup() {
  if [[ -n "$active" ]]; then bash "$REPO/infra/$active.sh" down; fi
  run_lock_release
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
while IFS=$'\t' read -r topo budget condition design mode scale trial order extra; do
  if [[ "$active/$active_budget" != "$topo/$budget" ]]; then
    if [[ -n "$active" ]]; then bash "$REPO/infra/$active.sh" down; fi
    active=""; active_budget="$budget"
    unset ADS_DB_CPUS ADS_DB_MEMORY ADS_PG_SHARED_BUFFERS ADS_PG_EFFECTIVE_CACHE ADS_PG_WORK_MEM ADS_PG_MAINTENANCE_MEM ADS_YB_TSERVER_MEMORY ADS_YB_MASTER_MEMORY
    if [[ "$budget" == constrained ]]; then
      export ADS_DB_CPUS=2 ADS_DB_MEMORY=256m ADS_PG_SHARED_BUFFERS=64MB ADS_PG_EFFECTIVE_CACHE=192MB ADS_PG_WORK_MEM=4MB ADS_PG_MAINTENANCE_MEM=32MB
    elif [[ "$budget" == large ]]; then
      export ADS_DB_CPUS=6 ADS_DB_MEMORY=9g ADS_YB_TSERVER_MEMORY=4831838208 ADS_YB_MASTER_MEMORY=1610612736
    fi
    if ! bash "$REPO/infra/$topo.sh" up > "$OUT/start-$topo-$budget.log" 2>&1; then
      echo "  - startup_failed: $topo/$budget" >> "$OUT/manifest.yaml"
      failed=$((failed+1)); bash "$REPO/infra/$topo.sh" down; continue
    fi
    active="$topo"
  fi
  number=$((number+1))
  name="$(printf '%04d' "$number")-$condition-$design-t$trial"
  dir="$OUT/$topo-$budget"
  mkdir -p "$dir/plans" "$dir/logs"
  case "$topo" in
    pg-single) engine=postgres; dbs=(pg-single) ;;
    yb-single) engine=yugabyte; dbs=(yb-single) ;;
    yb-cluster3) engine=yugabyte; dbs=(yb-n1 yb-n2 yb-n3) ;;
  esac
  dsn="$(bash "$REPO/infra/$topo.sh" dsn)"
  record_topology "$dir/$name-topology.yaml" "${dbs[@]}"
  memory="$(podman inspect "${dbs[0]}" --format '{{.HostConfig.Memory}}')"
  for c in "${dbs[@]}"; do
    podman exec "$c" cat /sys/fs/cgroup/cpu.stat /sys/fs/cgroup/memory.current /sys/fs/cgroup/memory.stat /sys/fs/cgroup/io.stat > "$dir/logs/$name-$c-before.txt"
  done
  read -r -a args <<< "$extra"
  log "[$number] $topo/$budget $condition $design trial $trial"
  echo "  - started: $topo-$budget/$name" >> "$OUT/manifest.yaml"
  fault_pid=""
  if [[ "$condition" == local-node-stop ]]; then
    (
      while ! grep -q arrival_measurement_started "$dir/logs/$name-client.txt" 2>/dev/null; do
        [[ -f "$dir/logs/$name-finished.flag" ]] && exit 0
        sleep 1
      done
      sleep 10
      date -u +%Y-%m-%dT%H:%M:%SZ > "$dir/logs/$name-fault.txt"
      podman stop -t 5 yb-n3 >> "$dir/logs/$name-fault.txt" 2>&1
    ) &
    fault_pid=$!
  fi
  if podman run --rm --name ads-study01-client --network "$NETWORK" --cpus 2 --memory 2g \
      -e "BENCH_REPO_COMMIT=$COMMIT" -e "BENCH_REPO_DESCRIBE=$DESCRIBE" -e BENCH_REPO_DIRTY=false \
      -e "BENCH_RUN_TAG=$RUN_TAG" -e "BENCH_IMAGE_ID=$IMAGE" \
      -v "$(hostpath "$dir"):/results" "$IMAGE" -cmd experiment -experiment "$mode" \
      -dsn "$dsn" -engine "$engine" -topology "$topo" -design "$design" -scale "$scale" \
      -conns 8 -load-conns 4 -duration "$DURATION" -warmup "$WARMUP" -trials 1 \
      -run-id "$RUN_ID" -environment "$ENVIRONMENT" -condition "$condition" -trial-id "$trial" -order "$order" \
      -db-memory-bytes "$memory" "${args[@]}" -out "/results/$name.json" -explain-out "/results/plans/$name.txt" \
      > "$dir/logs/$name-client.txt" 2>&1; then
    echo "    process_status: completed" >> "$OUT/manifest.yaml"
  else
    echo "    process_status: failed" >> "$OUT/manifest.yaml"
    failed=$((failed+1)); warn "cell failed; saved $dir/logs/$name-client.txt"
  fi
  if [[ -n "$fault_pid" ]]; then
    : > "$dir/logs/$name-finished.flag"
    wait "$fault_pid" || true
    podman start yb-n3 >> "$dir/logs/$name-fault.txt" 2>&1
    wait_sql yb-n3 ysqlsh yb-n3 5433 yugabyte yugabyte 180
    sql_in yb-n1 ysqlsh yb-n1 5433 yugabyte yugabyte "SELECT COUNT(*) FROM donation" > "$dir/logs/$name-recovered-count.txt"
  fi
  for c in "${dbs[@]}"; do
    podman logs --tail 150 "$c" > "$dir/logs/$name-$c.log" 2>&1 || true
    podman exec "$c" cat /sys/fs/cgroup/cpu.stat /sys/fs/cgroup/memory.current /sys/fs/cgroup/memory.stat /sys/fs/cgroup/io.stat > "$dir/logs/$name-$c-after.txt" || true
  done
done < "$OUT/ordered-cells.tsv"
echo "finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$OUT/manifest.yaml"
echo "failed_processes: $failed" >> "$OUT/manifest.yaml"
podman run --rm -v "$(hostpath "$STUDY_DIR"):/study" "$IMAGE" -cmd report-enhancements \
  -results "/study/results/$RUN_ID" -report-out "/study/reports/$RUN_ID.md"
log "completed $RUN_ID; failed processes $failed; report generated"
(( failed == 0 ))
