#!/usr/bin/env bash
# Phase A probe: the Redis capability and policy check (dev check, no measurement).
#
#   ./probe-cache.sh
#
# Every primitive the study's cache policy depends on is exercised against the real
# server BEFORE any measurement: SET NX PX with a unique token, a
# compare-and-delete release that refuses a wrong token, version-fenced publication,
# the exact hard TTL, and the recorded maxmemory policy. A probe that cannot satisfy
# one of them exits non-zero, so it cannot be skimmed.
set -uo pipefail
STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
source "$STUDY_DIR/study.env"
set +e

OUT="${STUDY_DIR}/results/devchecks/phase-a-redis"
mkdir -p "$OUT"

run_lock_acquire "study ${STUDY_ID} probe-cache.sh (cache capability probe)"
podman build -q -t "$BENCH_IMAGE" -f "$(winpath "$STUDY_DIR/Containerfile")" "$(winpath "$REPO")" >/dev/null \
  || die "image build failed"

bash "$REPO/infra/redis.sh" up || die "redis failed to start"
trap 'bash "$REPO/infra/redis.sh" down >/dev/null 2>&1' EXIT

log "probing Redis primitives and policy"
podman run --rm --network "$NETWORK" --cpus "$CLIENT_CPUS" --memory "$CLIENT_MEMORY" \
  "$BENCH_IMAGE" -cmd probe -redis-addr "${REDIS_NAME:-ads-redis}:6379" -out /dev/stdout \
  | tee "$OUT/probe.json"
status=${PIPESTATUS[0]}

podman exec "${REDIS_NAME:-ads-redis}" redis-cli info > "$OUT/redis-info.txt" 2>&1 || true
podman exec "${REDIS_NAME:-ads-redis}" redis-cli config get maxmemory >> "$OUT/redis-info.txt" 2>&1 || true
podman image inspect "$REDIS_IMAGE" --format '{{.Id}}' > "$OUT/redis-image-id.txt" 2>&1 || true

if [[ $status -ne 0 ]]; then
  warn "the Redis capability probe FAILED; see $OUT/probe.json"
  exit $status
fi
log "probe passed; evidence in $OUT"
