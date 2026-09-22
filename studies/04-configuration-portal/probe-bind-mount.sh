#!/usr/bin/env bash
# Phase B — the WSL bind-mount probe (HANDOFF.md, experiment B).
#
# Every result this study produces is written through a bind mount from the
# repository into the benchmark container. On Windows that mount source must be a
# Windows path; under WSL the repository lives at /mnt/c/... while the podman
# client is a Windows binary. So "which spelling does it accept?" is a real
# question with a *silent* failure mode: an unconverted or wrong-form source
# mounts an empty directory instead of erroring, and a matrix would then run to
# completion producing nothing, with no message (LESSONS_LEARNED, "Git Bash
# rewrites POSIX paths before Podman sees them").
#
# This probe answers it before any measurement:
#
#   1. a sentinel written in the task worktree is readable inside the container
#      at the expected path, byte for byte;
#   2. a file the container writes into the mounted results path survives on the
#      host and is readable from WSL;
#   3. which source spelling actually works -- /mnt/c/... (hostpath) or C:/...
#      (winpath) -- so the runner can use the proven one and record the choice.
#
# It starts a container, so it takes the benchmark lock like every other runner.
# It is read-only with respect to the databases: no database container is started.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
# lib.sh sets -e; a probe must record failures, not abort on the first one.
set +e

PROBE_IMAGE="${PROBE_IMAGE:-docker.io/library/golang:1.26-bookworm}"
OUT="$STUDY_DIR/results/devchecks/phase-b-bind-mount"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$OUT"

run_lock_acquire "study 04-configuration-portal probe-bind-mount.sh ${RUN_ID}"
need_podman

log "probe image: $PROBE_IMAGE"
podman image exists "$PROBE_IMAGE" || die "probe image $PROBE_IMAGE is not present locally; refusing to pull during a measurement session"

# The sentinel lives in the study tree, which is exactly the kind of path the
# runner mounts. Its content is fixed so the container's echo can be compared.
SENTINEL_DIR="$STUDY_DIR/results/devchecks/phase-b-bind-mount"
SENTINEL="$SENTINEL_DIR/sentinel.txt"
SENTINEL_TEXT="study-04 phase-b sentinel ${RUN_ID}"
printf '%s\n' "$SENTINEL_TEXT" > "$SENTINEL"
SENTINEL_SHA="$(sha256sum "$SENTINEL" | cut -d' ' -f1)"

REPORT="$OUT/report-${RUN_ID}.txt"
: > "$REPORT"
{
  echo "phase: B (WSL bind-mount probe)"
  echo "probe_started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "run_id: $RUN_ID"
  echo "worktree: $REPO"
  echo "study_dir: $STUDY_DIR"
  echo "image: $PROBE_IMAGE"
  echo "podman_client: $ADS_PODMAN"
  echo "podman_version: $(podman --version)"
  echo "sentinel_host_path: $SENTINEL"
  echo "sentinel_sha256: $SENTINEL_SHA"
  echo "sentinel_text: $SENTINEL_TEXT"
} >> "$REPORT"

# probe_form <label> <mount source>
# Mounts the study tree read-write at /probe and records what the container can
# see and write. A form "works" only if the sentinel is present with the right
# content *and* a file it writes comes back to the host with the right content.
probe_form() {
  local label="$1" src="$2"
  local cname="ads-probe-bind-${label}"
  local host_out="$SENTINEL_DIR/container-wrote-${label}.txt"
  local echo_out="$SENTINEL_DIR/container-read-${label}.txt"
  # The container's own output goes to a per-form file first. Grepping the
  # cumulative report instead would let one failing form mark every later form as
  # an empty mount -- the check would then agree with whatever ran first.
  local form_log="$SENTINEL_DIR/container-log-${label}.txt"
  rm -f "$host_out" "$echo_out" "$form_log"

  {
    echo ""
    echo "--- form: $label"
    echo "mount_source_as_passed: $src"
  } >> "$REPORT"

  podman rm -f "$cname" >/dev/null 2>&1 || true
  podman run --rm --name "$cname" \
    -v "$src:/probe" \
    "$PROBE_IMAGE" \
    bash -c "set -u
      echo container_cwd=\$(pwd)
      echo sentinel_visible=\$([ -f /probe/results/devchecks/phase-b-bind-mount/sentinel.txt ] && echo yes || echo no)
      if [ -f /probe/results/devchecks/phase-b-bind-mount/sentinel.txt ]; then
        cat /probe/results/devchecks/phase-b-bind-mount/sentinel.txt > /probe/results/devchecks/phase-b-bind-mount/container-read-${label}.txt
      fi
      echo written_from_container=${RUN_ID} > /probe/results/devchecks/phase-b-bind-mount/container-wrote-${label}.txt
      echo probe_container_uid=\$(id -u)
    " > "$form_log" 2>&1
  local st=$?
  cat "$form_log" >> "$REPORT"
  echo "container_exit_status: $st" >> "$REPORT"

  local host_seen="no" read_ok="no" write_ok="no"
  if [[ -f "$host_out" ]]; then
    host_seen="yes"
    [[ "$(cat "$host_out")" == "written_from_container=${RUN_ID}" ]] && write_ok="yes"
  fi
  if [[ -f "$echo_out" ]]; then
    [[ "$(sha256sum "$echo_out" | cut -d' ' -f1)" == "$SENTINEL_SHA" ]] && read_ok="yes"
  fi
  # An empty bind mount is the failure this probe exists to catch: the directory
  # exists but has no sentinel.
  local empty_mount="no"
  grep -q '^sentinel_visible=no$' "$form_log" && empty_mount="yes"

  {
    echo "host_saw_written_file: $host_seen"
    echo "sentinel_read_back_identical: $read_ok"
    echo "host_write_survived: $write_ok"
    echo "empty_mount_suspected: $empty_mount"
    if [[ "$read_ok" == "yes" && "$write_ok" == "yes" ]]; then
      echo "form_verdict: WORKS"
    else
      echo "form_verdict: FAILS"
    fi
  } >> "$REPORT"
}

# hostpath is what the existing scripts use; winpath is the WSL-specific
# Windows form. Both are tried, in both directions, so the runner can be told
# which one to use rather than guessing.
probe_form "hostpath" "$(hostpath "$STUDY_DIR")"
probe_form "winpath"  "$(winpath  "$STUDY_DIR")"

# A control: a source that cannot contain the sentinel must NOT be reported as
# working. Without this, a probe that always says "WORKS" would look convincing.
probe_form "bogus-control" "$(hostpath "$STUDY_DIR")/definitely-not-here-${RUN_ID}"

{
  echo ""
  echo "verdicts:"
  echo "  hostpath_works: $(grep -A12 '^--- form: hostpath$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"
  echo "  winpath_works:  $(grep -A12 '^--- form: winpath$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"
  echo "  bogus_works:    $(grep -A12 '^--- form: bogus-control$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"
  echo "probe_finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
} >> "$REPORT"

# The runner needs a yes/no, and a probe that cannot distinguish must say so
# rather than defaulting to the convenient answer. The window has to reach past
# the per-form diagnostics to the verdict line itself.
HOSTPATH_OK="$(grep -A12 '^--- form: hostpath$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"
WINPATH_OK="$(grep -A12 '^--- form: winpath$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"
BOGUS_OK="$(grep -A12 '^--- form: bogus-control$' "$REPORT" | awk -F': ' '/form_verdict/{print $2}')"

cat "$REPORT"

echo
if [[ "$BOGUS_OK" == "WORKS" ]]; then
  warn "the bogus control reported WORKS -- the probe cannot tell an empty mount from a real one; result unusable"
  exit 1
fi
if [[ "$HOSTPATH_OK" == "WORKS" ]]; then
  log "decision: hostpath() is sufficient on this host; the runner passes -v \$(hostpath ...) unchanged"
  echo "RUNNER_DECISION=hostpath" >> "$REPORT"
  exit 0
fi
if [[ "$WINPATH_OK" == "WORKS" ]]; then
  log "decision: hostpath() is NOT sufficient here; study 04's runner must pass winpath() as the bind source"
  echo "RUNNER_DECISION=winpath" >> "$REPORT"
  exit 0
fi
warn "neither /mnt/c nor C:/ mounted the study tree; measurement is impossible until this is resolved"
echo "RUNNER_DECISION=none" >> "$REPORT"
exit 1
