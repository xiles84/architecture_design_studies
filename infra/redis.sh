#!/usr/bin/env bash
# Shared Redis cache under test (study 05). Pinned image from infra/versions.env.
#
# Usage: redis.sh up | down | cli | info | logs
#
# The policy is configured EXPLICITLY and recorded by the harness's INFO reading,
# because the study compares Redis's approximate LRU against an exact in-process
# LRU: leaving the policy at whatever the image defaults to would make the
# comparison a comparison of two unknowns.
#
# `--save ''` and `--appendonly no`: the cache's contents are not a result and must
# not survive a restart. Persistence would also make a "restart" fault measure the
# reload of a dataset rather than an empty cache.
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

NAME="${REDIS_NAME:-ads-redis}"
MAXMEM="${ADS_REDIS_MAXMEMORY:-32mb}"
POLICY="${ADS_REDIS_POLICY:-allkeys-lru}"
SAMPLES="${ADS_REDIS_SAMPLES:-5}"

# Default is the ADD-CACHE framing's cache budget. The equal-total framing deducts
# these resources from the database instead of adding them, and the runner passes
# both budgets explicitly so the two conditions are never pooled.
CACHE_CPUS="${ADS_CACHE_CPUS:-1}"
CACHE_MEMORY="${ADS_CACHE_MEMORY:-805306368}"

cmd_up() {
  run_lock_guard
  need_podman
  net_ensure
  rm_container "$NAME"
  log "starting $NAME ($REDIS_IMAGE, ${CACHE_CPUS} cpus, ${CACHE_MEMORY} bytes, maxmemory=${MAXMEM}, ${POLICY}, samples=${SAMPLES})"
  podman run -d --name "$NAME" --hostname "$NAME" \
    --network "$NETWORK" \
    --cpus "$CACHE_CPUS" --memory "$CACHE_MEMORY" \
    "$REDIS_IMAGE" \
    redis-server --save '' --appendonly no \
      --maxmemory "$MAXMEM" --maxmemory-policy "$POLICY" --maxmemory-samples "$SAMPLES" \
      --port 6379 >/dev/null
  local start now
  start=$(date +%s)
  while true; do
    if podman exec "$NAME" redis-cli ping >/dev/null 2>&1; then
      break
    fi
    now=$(date +%s)
    (( now - start > 60 )) && die "$NAME not ready after 60s"
    sleep 1
  done
  log "$NAME ready on container port 6379"
  podman exec "$NAME" redis-cli info memory | grep -E '^(maxmemory|maxmemory_policy|maxmemory_samples|used_memory):' | sed 's/^/    /'
}

cmd_down() {
  run_lock_guard
  rm_container "$NAME"
  podman volume rm -f ads-redis-data >/dev/null 2>&1 || true
  log "$NAME removed"
}

case "${1:-up}" in
  up) cmd_up ;;
  down) cmd_down ;;
  cli) shift; podman exec -it "$NAME" redis-cli "$@" ;;
  info) podman exec "$NAME" redis-cli info ;;
  logs) podman logs --tail "${2:-80}" "$NAME" ;;
  dsn) echo "${NAME}:6379" ;;
  *) die "usage: $0 {up|down|cli|info|logs|dsn}" ;;
esac
