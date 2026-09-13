#!/usr/bin/env bash
# Single-node PostgreSQL under test.
#
# Usage: pg-single.sh up | down | psql | logs
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

NAME="pg-single"
VOLUME="ads-pg-single-data"
DB=bench
USER=bench
PASS=bench

# Configuration is set explicitly rather than left at the image defaults.
# The image ships shared_buffers=128MB, which would make every result a
# measurement of a starved buffer cache instead of a measurement of the design.
# These values are the conventional starting points for a 3 GiB node; they are
# deliberately ordinary, not tuned per design.
pg_settings=(
  -c shared_buffers=768MB
  -c effective_cache_size=2GB
  -c work_mem=32MB
  -c maintenance_work_mem=256MB
  -c max_connections=200
  -c random_page_cost=1.1          # NVMe, not a spinning disk
  -c effective_io_concurrency=200
  -c max_wal_size=4GB
  -c checkpoint_completion_target=0.9
  -c track_io_timing=on            # makes EXPLAIN (BUFFERS) report real I/O time
  -c track_functions=all
  -c autovacuum_naptime=10s        # keep bloat visible but not runaway
)

cmd_up() {
  run_lock_guard
  need_podman
  net_ensure
  rm_container "$NAME"
  log "starting $NAME ($PG_IMAGE, ${DB_CPUS} cpus, ${DB_MEMORY})"
  podman run -d --name "$NAME" \
    --network "$NETWORK" \
    --cpus "$DB_CPUS" --memory "$DB_MEMORY" \
    -e POSTGRES_USER="$USER" -e POSTGRES_PASSWORD="$PASS" -e POSTGRES_DB="$DB" \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "$VOLUME":/var/lib/postgresql/data \
    -p "${PG_PORT}:5432" \
    "$PG_IMAGE" "${pg_settings[@]}" >/dev/null
  wait_sql "$NAME" psql 127.0.0.1 5432 "$USER" "$DB" 180
  log "$NAME ready on container port 5432 (host ${PG_PORT})"
}

cmd_down() {
  run_lock_guard
  rm_container "$NAME"
  # The volume is removed too: a benchmark must never start from a previous
  # run's pages, statistics or bloat.
  podman volume rm -f "$VOLUME" >/dev/null 2>&1 || true
  log "$NAME removed (including its data volume)"
}

cmd_psql() { podman exec -it "$NAME" psql -U "$USER" -d "$DB" "$@"; }
cmd_logs() { podman logs --tail "${1:-80}" "$NAME"; }

case "${1:-up}" in
  up) cmd_up ;;
  down) cmd_down ;;
  psql) shift; cmd_psql "$@" ;;
  logs) shift; cmd_logs "$@" ;;
  dsn) echo "postgres://${USER}:${PASS}@${NAME}:5432/${DB}?sslmode=disable" ;;
  *) die "usage: $0 {up|down|psql|logs|dsn}" ;;
esac
