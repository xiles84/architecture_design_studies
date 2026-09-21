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

# Explicit alternative resource conditions. Defaults preserve the published
# per-node experiment; each alternative is inspected and named in its results.
DB_CPUS="${ADS_DB_CPUS:-$DB_CPUS}"
DB_MEMORY="${ADS_DB_MEMORY:-$DB_MEMORY}"

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m warn\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[1;31mfatal\033[0m %s\n' "$*" >&2; exit 1; }

# ---------------------------------------------------------------------------
# Container engine resolution.
#
# On Linux and macOS `podman` is a native binary, and under Git Bash on Windows
# it is podman.exe on PATH. WSL is the awkward case: the distro has no podman of
# its own, while the Windows podman.exe drives the *same* podman-machine-default
# engine -- the one that already holds every study's images, named volumes and
# the benchmark lock. From WSL it is therefore the only client that leaves us on
# the same engine, and the same lock, as every other session.
#
# A WSL-local podman is never an option, and this file must not make one easy:
# a second engine has different images, different volumes and *no shared lock*,
# so two studies could measure at once on the same cores. A second runner has to
# fail on the lock; it must never succeed somewhere else instead.
#
# Precedence is explicit and never guesses: ADS_PODMAN > PODMAN > `podman` on
# PATH > the Windows client. `podman` is then a function, so every call site in
# every script -- and `command -v podman` -- keeps working unchanged.
# ---------------------------------------------------------------------------
ADS_PODMAN="${ADS_PODMAN:-${PODMAN:-}}"
if [[ -z "$ADS_PODMAN" ]]; then
  if command -v podman >/dev/null 2>&1; then
    ADS_PODMAN="podman"
  elif [[ -x "/mnt/c/Program Files/RedHat/Podman/podman.exe" ]]; then
    ADS_PODMAN="/mnt/c/Program Files/RedHat/Podman/podman.exe"
  else
    ADS_PODMAN="podman"
  fi
fi
export ADS_PODMAN

podman() { command "$ADS_PODMAN" "$@"; }

need_podman() {
  if [[ "$ADS_PODMAN" == */* || "$ADS_PODMAN" == *.exe ]]; then
    [[ -x "$ADS_PODMAN" ]] || die "podman client not executable: $ADS_PODMAN (set PODMAN=... to point at one)"
  elif ! command -v "$ADS_PODMAN" >/dev/null 2>&1; then
    die "podman not found in PATH (set PODMAN=... to point at a client)"
  fi
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

# winpath converts a WSL path (/mnt/c/extra/...) into the Windows form
# (C:/extra/...) that the Windows podman.exe may require as a bind-mount source.
# It is a no-op under Git Bash (where cygpath already produced C:/...), on Linux
# and on macOS.
#
# hostpath() is deliberately left alone rather than extended. Its contract is
# "a path the podman binary accepts", and under Git Bash that is already the
# Windows form; but studies 01-03 also pass its result to git, so widening its
# meaning would change code that already produced published results. Study 04
# calls winpath() only where its phase-B bind-mount probe proves the Windows
# form is needed.
winpath() {
  local p="$1"
  if [[ "$p" =~ ^/mnt/([a-zA-Z])/(.*)$ ]]; then
    printf '%s:/%s' "${BASH_REMATCH[1]^^}" "${BASH_REMATCH[2]}"
  else
    printf '%s' "$p"
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

# ---------------------------------------------------------------------------
# Benchmark lock: one measurement on this machine at a time.
#
# Several AI sessions may work at once, in separate git worktrees. Code and
# analyses can be written in parallel; measurements cannot. The infra scripts
# reuse container names in every study and worktree (a second `up` destroys the
# first session's database), and even with unique names two matrices would
# share the same eight cores and silently change each other's numbers.
#
# The lock is a podman named volume. `podman volume create` fails atomically when
# the volume exists, and the podman machine is shared by every worktree and every
# shell (Git Bash, PowerShell, WSL, any vendor's agent) -- the lock lives exactly
# where the contended resource lives. Labels record who holds it.
#
# A holder killed with SIGKILL leaves a stale lock. Remove it by hand only after
# `podman ps` shows no benchmark or database container:  podman volume rm ads-run-lock
# ---------------------------------------------------------------------------
RUN_LOCK_VOLUME="ads-run-lock"

# run_lock_acquire <description>
# Takes the lock for the rest of this script, or dies naming the holder. Nested
# calls (a runner that calls another runner) inherit the parent's lock.
run_lock_acquire() {
  local desc="$1"
  if [[ -n "${ADS_RUN_LOCK_HELD:-}" ]]; then
    return 0
  fi
  if podman volume create \
      --label "ads.holder=${desc}" \
      --label "ads.started=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
      --label "ads.worktree=$(hostpath "$REPO_ROOT")" \
      --label "ads.branch=$(git -C "$(hostpath "$REPO_ROOT")" rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)" \
      "$RUN_LOCK_VOLUME" >/dev/null 2>&1; then
    export ADS_RUN_LOCK_HELD="$desc"
    ADS_RUN_LOCK_OWNER_PID="$$"
    # EXIT traps do not run on an untrapped signal; turning signals into exits
    # makes Ctrl+C and `kill` release the lock too. Bash runs a trap only after
    # the current foreground command returns: Ctrl+C reaches the whole process
    # group and releases at once; `kill <runner-pid>` releases when the running
    # cell's podman command ends. (Tested 2026-09-13.)
    trap 'run_lock_release' EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    log "benchmark lock acquired: ${desc}"
    return 0
  fi
  warn "another benchmark holds this machine's lock:"
  run_lock_describe >&2
  die "refusing to run concurrently: measurements on a shared machine would corrupt each other (see infra/lib.sh)"
}

run_lock_release() {
  if [[ "${ADS_RUN_LOCK_OWNER_PID:-}" == "$$" ]]; then
    podman volume rm -f "$RUN_LOCK_VOLUME" >/dev/null 2>&1 || true
    unset ADS_RUN_LOCK_OWNER_PID ADS_RUN_LOCK_HELD
  fi
}

run_lock_describe() {
  podman volume inspect "$RUN_LOCK_VOLUME" \
    --format '  holder:   {{index .Labels "ads.holder"}}
  started:  {{index .Labels "ads.started"}}
  worktree: {{index .Labels "ads.worktree"}}
  branch:   {{index .Labels "ads.branch"}}' 2>/dev/null || echo "  (lock volume vanished)"
}

# run_lock_guard: for the infra topology scripts. Starting or stopping a database
# while someone else's benchmark holds the lock would destroy their run. Scripts
# invoked by a lock holder pass through; a person poking at a database by hand
# while nothing is measured passes through; ADS_IGNORE_RUN_LOCK=1 overrides.
run_lock_guard() {
  [[ -n "${ADS_RUN_LOCK_HELD:-}" || "${ADS_IGNORE_RUN_LOCK:-}" == "1" ]] && return 0
  if podman volume exists "$RUN_LOCK_VOLUME" 2>/dev/null; then
    warn "a benchmark is running on this machine:"
    run_lock_describe >&2
    die "not touching database containers while it runs (set ADS_IGNORE_RUN_LOCK=1 only if you are certain it is stale)"
  fi
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
