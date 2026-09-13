#!/usr/bin/env bash
# Three-node YugabyteDB cluster, replication factor 3.
#
# Same engine, same image, same flags and the same PER-NODE resource budget as
# yb-single. The only variable is the node count, which is what makes the pair a
# usable experiment on sharding and replication rather than two unrelated runs.
#
# Usage: yb-cluster3.sh up | down | ysqlsh | status | logs <n> | dsn
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

NODES=(yb-n1 yb-n2 yb-n3)
DB=yugabyte
USER=yugabyte

TSERVER_FLAGS="memory_limit_hard_bytes=1610612736,ysql_num_shards_per_tserver=2,yb_num_shards_per_tserver=2"
MASTER_FLAGS="memory_limit_hard_bytes=536870912"

# Optional per-study tserver flags; see yb-single.sh. Unset = study 01's nodes.
TSERVER_FLAGS="${TSERVER_FLAGS}${YB_EXTRA_TSERVER_FLAGS:+,$YB_EXTRA_TSERVER_FLAGS}"

start_node() {
  local name="$1" join="$2"
  local extra=()
  [[ -n "$join" ]] && extra=(--join="$join")

  podman run -d --name "$name" --hostname "$name" \
    --network "$NETWORK" \
    --cpus "$DB_CPUS" --memory "$DB_MEMORY" \
    -v "ads-${name}-data":/home/yugabyte/yb_data \
    "$YB_IMAGE" \
    bin/yugabyted start --background=false \
      --advertise_address="$name" \
      --base_dir=/home/yugabyte/yb_data \
      --callhome=false \
      --tserver_flags="$TSERVER_FLAGS" \
      --master_flags="$MASTER_FLAGS" \
      "${extra[@]}" >/dev/null
}

cmd_up() {
  need_podman
  net_ensure
  for n in "${NODES[@]}"; do rm_container "$n"; podman volume rm -f "ads-${n}-data" >/dev/null 2>&1 || true; done

  log "starting ${NODES[0]} (cluster seed)"
  start_node "${NODES[0]}" ""
  wait_sql "${NODES[0]}" ysqlsh "${NODES[0]}" 5433 "$USER" "$DB" 300

  # Nodes join one at a time. Starting all three at once races the master
  # election and intermittently produces a two-node universe with the third
  # node orphaned -- which then quietly reports RF=1 results as if they were RF=3.
  for n in "${NODES[@]:1}"; do
    log "joining $n"
    start_node "$n" "${NODES[0]}"
    wait_sql "$n" ysqlsh "$n" 5433 "$USER" "$DB" 300
  done

  # On this image (2025.2.6.0) yugabyted raises the universe to RF=3 by itself as
  # the third node joins, so no placement step is needed. Older versions created
  # the universe at RF=1 and left it there; if that happens the assertion in
  # cmd_status fails loudly rather than letting an RF=1 universe be reported as a
  # replicated cluster.
  if [[ "$(current_rf)" != "3" ]]; then
    log "replication factor is not 3 yet; configuring data placement"
    podman exec "${NODES[0]}" bash -lc \
      "bin/yugabyted configure data_placement --base_dir=/home/yugabyte/yb_data --fault_tolerance=zone" >/dev/null 2>&1 || true
  fi

  cmd_status
}

# current_rf reads the live replica count straight out of the universe config.
# This is the authoritative answer; node count alone does not imply replication.
current_rf() {
  podman exec "${NODES[0]}" bash -lc \
    "bin/yb-admin --master_addresses=${NODES[0]}:7100,${NODES[1]}:7100,${NODES[2]}:7100 get_universe_config 2>/dev/null" \
    | grep -oE '"numReplicas":[0-9]+' | head -1 | cut -d: -f2
}

cmd_status() {
  local n_servers rf
  n_servers=$(sql_in "${NODES[0]}" ysqlsh "${NODES[0]}" 5433 "$USER" "$DB" "SELECT count(*) FROM yb_servers()" | tr -d '[:space:]')
  rf=$(current_rf)
  log "tservers: ${n_servers}   replication factor: ${rf:-unknown}"
  if [[ "$n_servers" != "3" ]]; then
    die "expected 3 tservers, found ${n_servers} -- refusing to report this as a 3-node cluster"
  fi
  if [[ -n "$rf" && "$rf" != "3" ]]; then
    die "expected replication factor 3, found ${rf}"
  fi
}

cmd_down() {
  for n in "${NODES[@]}"; do
    rm_container "$n"
    podman volume rm -f "ads-${n}-data" >/dev/null 2>&1 || true
  done
  log "3-node cluster removed (including data volumes)"
}

case "${1:-up}" in
  up) cmd_up ;;
  down) cmd_down ;;
  status) cmd_status ;;
  ysqlsh) shift; podman exec -it "${NODES[0]}" bash -lc "ysqlsh -h ${NODES[0]} -p 5433 -U $USER -d $DB $*" ;;
  logs) podman logs --tail "${3:-80}" "${2:-${NODES[0]}}" ;;
  dsn) echo "postgres://${USER}@${NODES[0]}:5433/${DB}?sslmode=disable" ;;
  *) die "usage: $0 {up|down|status|ysqlsh|logs|dsn}" ;;
esac
