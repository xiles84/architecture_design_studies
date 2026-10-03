#!/usr/bin/env bash
# A bounded experiment under the shared measurement lock; never touches another
# study's database names, volume, results, images or signed analyses.
export ADS_DB_CPUS=1 ADS_DB_MEMORY=512m ADS_CLIENT_CPUS=1 ADS_CLIENT_MEMORY=512m
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$HERE/../../../infra/lib.sh"
PREFIX="${1:?usage: run.sh <new-prefix> [duration] [whole-runs]}"
DURATION="${2:-3s}"
WHOLE_RUNS="${3:-3}"
[[ "$PREFIX" =~ ^[a-zA-Z0-9._-]+$ ]] || die "invalid run prefix"
[[ "$WHOLE_RUNS" =~ ^[1-9][0-9]*$ ]] || die "invalid whole-run count"
IMAGE="localhost/cache-protocol-v5:20261003"
PG=cp-pg
CACHE=cp-redis
VOL=ads-cp-pg-data
STUDY="$HERE/.."
RESULTS="$STUDY/results"
REPORTS="$STUDY/reports"
need_podman
gitroot="$(hostpath "$REPO_ROOT")"
COMMIT="$(git -C "$gitroot" rev-parse HEAD)"
[[ -z "$(git -C "$gitroot" status --porcelain)" ]] || die "commit the full producing tree first"
for ((whole=1;whole<=WHOLE_RUNS;whole++)); do
 [[ ! -e "$RESULTS/$PREFIX-r$whole" && ! -e "$REPORTS/$PREFIX-r$whole.md" ]] || die "result/report already exists"
 ! git -C "$gitroot" rev-parse "refs/tags/run/study-05/$PREFIX-r$whole" >/dev/null 2>&1 || die "run tag already exists"
done
run_lock_acquire "study05 protocol v5 $PREFIX"
started=0
cleanup_services(){
 if [[ "$started" == 1 ]]; then
  podman rm -f "$PG" "$CACHE" >/dev/null 2>&1 || true
  podman volume rm "$VOL" >/dev/null 2>&1 || true
  started=0
 fi
}
cleanup(){ cleanup_services; run_lock_release; }
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
# Existing names are evidence of an interrupted/foreign owner, not permission
# to delete them. A fresh run uses only names proven absent while holding lock.
for c in "$PG" "$CACHE"; do
 ! podman container exists "$c" || die "owned service name already exists: $c; inspect before retry"
done
! podman volume exists "$VOL" || die "owned data volume already exists; inspect before retry"
net_ensure
log "building separate protocol client"
podman build -t "$IMAGE" -f "$(buildpath "$HERE/Containerfile")" "$(buildpath "$REPO_ROOT")"
IMAGE_ID="$(podman image inspect "$IMAGE" --format '{{.Id}}')"
start_services(){
 podman volume create "$VOL" >/dev/null
 started=1
 podman run -d --name "$PG" --network "$NETWORK" --cpus 1 --memory 512m \
  -e POSTGRES_USER=bench -e POSTGRES_PASSWORD=bench -e POSTGRES_DB=bench \
  -v "$VOL":/var/lib/postgresql/data "$PG_IMAGE" \
  -c shared_buffers=64MB -c effective_cache_size=256MB -c work_mem=4MB \
  -c maintenance_work_mem=64MB -c max_connections=40 -c track_io_timing=on >/dev/null
 wait_sql "$PG" psql 127.0.0.1 5432 bench bench 120
 podman run -d --name "$CACHE" --network "$NETWORK" --cpus 1 --memory 256m \
  "$REDIS_IMAGE" redis-server --save '' --appendonly no --maxmemory 64mb \
  --maxmemory-policy volatile-lru --maxmemory-samples 5 >/dev/null
 for ((i=0;i<60;i++)); do
  if podman exec "$CACHE" redis-cli ping >/dev/null 2>&1; then return; fi
  sleep 1
 done
 die "cache did not become ready"
}
capture(){
 local dir="$1" label="$2"
 mkdir -p "$dir"
 podman inspect "$PG" "$CACHE" > "$dir/$label-containers.json"
 for c in "$PG" "$CACHE"; do
  podman exec "$c" sh -c 'cat /sys/fs/cgroup/cpu.stat; cat /sys/fs/cgroup/memory.events; cat /sys/fs/cgroup/memory.current' > "$dir/$label-$c-cgroup.txt"
 done
}
failed=0
for ((whole=1;whole<=WHOLE_RUNS;whole++)); do
 RUN="$PREFIX-r$whole"
 TAG="run/study-05/$RUN"
 OUT="$RESULTS/$RUN"
 mkdir -p "$OUT/cells" "$OUT/resources" "$REPORTS"
 git -C "$gitroot" tag -a "$TAG" "$COMMIT" -m "Study 05 bounded cache protocol run $RUN; GPT-6 via Codex desktop"
 start_services
 podman machine inspect > "$OUT/resources/machine-config.json"
 podman machine ssh cat /proc/meminfo > "$OUT/resources/podman-guest-memory.txt"
 podman machine ssh nproc > "$OUT/resources/podman-guest-cpus.txt"
 record_host_state "$OUT/resources/runner-host.txt"
 podman image inspect "$PG_IMAGE" "$REDIS_IMAGE" "$IMAGE" > "$OUT/resources/images.json"
 # Bind-mount proof is read back in the container before admitting measurements.
 printf '%s\n' "$RUN" > "$OUT/resources/bind-sentinel.txt"
 mountroot="$(winpath "$(hostpath "$OUT")")"
 probe="$(podman run --rm -v "$mountroot:/probe:ro" --entrypoint /bin/cat localhost/ads-queue:1 /probe/resources/bind-sentinel.txt)"
 [[ "$probe" == "$RUN" ]] || die "bind-mount probe failed"
 seed=$((42+whole*104729))
 podman run --rm "$IMAGE" --cmd order --seed "$seed" > "$OUT/order.txt"
 printf '%s\n' "run_id: $RUN" "repo_commit: $COMMIT" "run_tag: $TAG" "environment: xiles-podman-4cpu-2g-config-protocol-20261003" \
  "client_image: $IMAGE_ID" "ordering_seed: $seed" "whole_run: $whole/$WHOLE_RUNS" \
  "trials_per_phase: 3" "duration: $DURATION" "keys: 128" \
  "resources: pg=1cpu/512MiB redis=1cpu/256MiB client=1cpu/512MiB reserve=1cpu/768MiB" \
  "cluster: NOT MEASURED (resource coverage gap)" "placement: single-node fixed; multi-host contrast NOT MEASURED" > "$OUT/manifest.yaml"
 cell(){
  local id="$1" arm="$2" protocol="$3" strategy="$4" kind="$5"
  capture "$OUT/resources/$id" before
  log "$RUN $id"
  if podman run --rm --network "$NETWORK" --cpus 1 --memory 512m -v "$mountroot:/results" "$IMAGE" \
    --run-id "$RUN" --commit "$COMMIT" --tag "$TAG" --image "$IMAGE_ID" \
    --arm "$arm" --protocol "$protocol" --strategy "$strategy" --kind "$kind" \
    --duration "$DURATION" --trials 3 --seed "$seed" --out "/results/cells/$id.json" \
    > "$OUT/resources/$id/client.log" 2>&1; then
   printf 'cell_%s: passed\n' "$id" >> "$OUT/manifest.yaml"
  else
   printf 'cell_%s: failed\n' "$id" >> "$OUT/manifest.yaml"
   failed=1
   cat "$OUT/resources/$id/client.log" >&2
  fi
  capture "$OUT/resources/$id" after
 }
 while IFS= read -r id; do
  protocol="${id%%-*}"
  rest="${id#*-}"
  strategy="${rest%%-*}"
  cell "A-$id" A "$protocol" "$strategy" hard
 done < "$OUT/order.txt"
 # Rotate deletion package position across whole runs, retaining every package.
 case $((whole%3)) in
  1) packages=(soft hard tombstone);;
  2) packages=(hard tombstone soft);;
  0) packages=(tombstone soft hard);;
 esac
 for kind in "${packages[@]}"; do cell "B-$kind" B sync cas "$kind"; done
 analysesroot="$(winpath "$(hostpath "$REPORTS/analyses")")"
 mkdir -p "$REPORTS/analyses"
 podman run --rm -v "$mountroot:/results:ro" -v "$(winpath "$(hostpath "$REPORTS")"):/reports" \
  -v "$analysesroot:/analyses:ro" "$IMAGE" --cmd report --inputs /results --out "/reports/$RUN.md"
 podman run --rm -v "$mountroot:/results:ro" "$IMAGE" --cmd verify --inputs /results \
  > "$OUT/resources/verification.txt" 2>&1 || failed=1
 podman logs "$PG" > "$OUT/resources/postgresql.log" 2>&1
 podman logs "$CACHE" > "$OUT/resources/redis.log" 2>&1
 cleanup_services
done
exit "$failed"
