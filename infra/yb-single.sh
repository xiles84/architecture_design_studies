#!/usr/bin/env bash
# Single-node YugabyteDB (replication factor 1).
#
# This is the control for the sharding experiment: same engine, same SQL, same
# storage layer as the three-node cluster, but nothing to distribute to. Any
# difference between this and yb-cluster3 is the cost or benefit of distribution
# itself, not of switching databases.
#
# Usage: yb-single.sh up | down | ysqlsh | logs | dsn
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

NAME="yb-single"
VOLUME="ads-yb-single-data"
DB=yugabyte
USER=yugabyte

# Tablet count is pinned instead of left to auto-detection. YugabyteDB otherwise
# derives shards-per-table from the CPU count it sees, which on this host varies
# with how the container happened to be scheduled -- and a benchmark whose shard
# count moves between runs is not a benchmark.
#
# memory_limit_hard_bytes is also pinned: YugabyteDB sizes its block cache from
# the machine's memory, and inside a memory-limited container that leads it to
# plan for far more RAM than the cgroup will actually hand out.
TSERVER_FLAGS="memory_limit_hard_bytes=1610612736,ysql_num_shards_per_tserver=2,yb_num_shards_per_tserver=2"
MASTER_FLAGS="memory_limit_hard_bytes=536870912"

# A study may need extra tserver flags; the default (unset) leaves the node
# exactly as study 01 measured it. Study 02 sets
# yb_enable_read_committed_isolation=true. On this image READ COMMITTED is
# already effective without it (checked: yb_effective_transaction_isolation_level
# reports "read committed"), but older releases silently ran RC as Snapshot
# Isolation, so a study whose designs depend on RC semantics pins it explicitly.
TSERVER_FLAGS="${TSERVER_FLAGS}${YB_EXTRA_TSERVER_FLAGS:+,$YB_EXTRA_TSERVER_FLAGS}"

cmd_up() {
  run_lock_guard
  need_podman
  net_ensure
  rm_container "$NAME"
  log "starting $NAME ($YB_IMAGE, ${DB_CPUS} cpus, ${DB_MEMORY})"
  podman run -d --name "$NAME" --hostname "$NAME" \
    --network "$NETWORK" \
    --cpus "$DB_CPUS" --memory "$DB_MEMORY" \
    -v "$VOLUME":/home/yugabyte/yb_data \
    -p "${YB_PORT}:5433" -p "${YB_UI_PORT}:15433" \
    "$YB_IMAGE" \
    bin/yugabyted start --background=false \
      --advertise_address="$NAME" \
      --base_dir=/home/yugabyte/yb_data \
      --callhome=false \
      --tserver_flags="$TSERVER_FLAGS" \
      --master_flags="$MASTER_FLAGS" >/dev/null

  # Note the host: yugabyted binds YSQL to its --advertise_address, not to
  # loopback, so probing 127.0.0.1 inside its own container is refused.
  wait_sql "$NAME" ysqlsh "$NAME" 5433 "$USER" "$DB" 300
  log "$NAME ready on container port 5433 (host ${YB_PORT})"
  sql_in "$NAME" ysqlsh "$NAME" 5433 "$USER" "$DB" "SELECT count(*) FROM yb_servers()" \
    | sed 's/^/    tservers: /'
}

cmd_down() {
  run_lock_guard
  rm_container "$NAME"
  podman volume rm -f "$VOLUME" >/dev/null 2>&1 || true
  log "$NAME removed (including its data volume)"
}

cmd_ysqlsh() { podman exec -it "$NAME" bash -lc "ysqlsh -h $NAME -p 5433 -U $USER -d $DB $*"; }
cmd_logs()   { podman logs --tail "${1:-80}" "$NAME"; }

case "${1:-up}" in
  up) cmd_up ;;
  down) cmd_down ;;
  ysqlsh) shift; cmd_ysqlsh "$@" ;;
  logs) shift; cmd_logs "$@" ;;
  dsn) echo "postgres://${USER}@${NAME}:5433/${DB}?sslmode=disable" ;;
  *) die "usage: $0 {up|down|ysqlsh|logs|dsn}" ;;
esac
