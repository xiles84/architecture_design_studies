#!/usr/bin/env bash
# Shared helpers for the podman orchestration scripts.
#
# Sourced, never executed directly.

# Git Bash on Windows rewrites arguments that look like POSIX paths into Windows
# paths before podman sees them, which turns "-w /src" into "-w C:/Program Files/Git/src".
# Disabling that conversion is what makes these scripts work identically on
# Windows, Linux and macOS.
export MSYS_NO_PATHCONV=1
export MSYS2_ARG_CONV_EXCL='*'

set -euo pipefail

INFRA_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$INFRA_DIR/.." && pwd)"
# shellcheck source=./versions.env
source "$INFRA_DIR/versions.env"

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m warn\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[1;31mfatal\033[0m %s\n' "$*" >&2; exit 1; }

need_podman() {
  command -v podman >/dev/null 2>&1 || die "podman not found in PATH"
  podman info >/dev/null 2>&1 || die "podman is installed but not reachable. On Windows/macOS run: podman machine start"
}

net_ensure() {
  podman network exists "$NETWORK" 2>/dev/null || podman network create "$NETWORK" >/dev/null
}

container_exists() { podman container exists "$1" 2>/dev/null; }

# hostpath converts a shell path into one the podman binary accepts as a bind
# mount source. Under Git Bash on Windows, $(pwd) yields /c/extra/... but podman
# is a native Windows binary expecting C:/extra/..., so an unconverted path
# silently creates a container-side empty directory instead of mounting the
# host's -- results vanish with no error. No-op everywhere else.
hostpath() {
  if command -v cygpath >/dev/null 2>&1; then
    cygpath -m "$1"
  else
    printf '%s' "$1"
  fi
}

rm_container() {
  if container_exists "$1"; then
    podman rm -f "$1" >/dev/null 2>&1 || true
  fi
}

# sql_in <container> <client> <host> <port> <user> <db> <sql>
# Runs one statement from inside a container and prints the bare result.
#
# The client binary differs by image (postgres ships psql, yugabyte ships ysqlsh)
# and so does the address: yugabyted binds YSQL to its --advertise_address, NOT
# to 127.0.0.1, so a loopback probe inside its own container is refused.
sql_in() {
  local c="$1" client="$2" host="$3" port="$4" user="$5" db="$6" sql="$7"
  podman exec "$c" bash -lc "$client -h $host -p $port -U $user -d $db -Atc \"$sql\""
}

# wait_sql <container> <client> <host> <port> <user> <db> [timeout_s]
# Polls with a real query, so readiness means "accepts queries", not merely
# "port is open". A port-open check races: PostgreSQL binds before it finishes
# recovery, and YugabyteDB binds before its tablets have elected leaders.
wait_sql() {
  local c="$1" client="$2" host="$3" port="$4" user="$5" db="$6" timeout="${7:-180}"
  local start now
  start=$(date +%s)
  while true; do
    if sql_in "$c" "$client" "$host" "$port" "$user" "$db" "SELECT 1" >/dev/null 2>&1; then
      return 0
    fi
    now=$(date +%s)
    if (( now - start > timeout )); then
      warn "last 40 log lines from $c:"
      podman logs --tail 40 "$c" >&2 || true
      die "$c not ready after ${timeout}s"
    fi
    sleep 2
  done
}

# Records what was actually running, so a result can be tied to a topology
# after the fact rather than trusted from a script argument.
record_topology() {
  local out="$1"; shift
  mkdir -p "$(dirname "$out")"
  {
    echo "captured_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo "podman_version: $(podman --version)"
    echo "containers:"
    for c in "$@"; do
      if container_exists "$c"; then
        podman inspect "$c" --format \
          '  - name: {{.Name}}
    image: {{.ImageName}}
    cpus: {{.HostConfig.NanoCpus}}
    memory_bytes: {{.HostConfig.Memory}}
    started: {{.State.StartedAt}}'
      fi
    done
  } > "$out"
}
