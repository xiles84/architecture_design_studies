#!/usr/bin/env bash
# Optional host-side cross-platform check for the queue.
#
# It runs the worktree-portability test against a repository placed under
# /mnt/<drive>/ so that a *native Windows* git and the shell's WSL git operate
# on the same common Git directory. The canonical test path is the Podman image
# (`tools/queue/queue --build`); inside a Linux container there is no second git
# binary, so this script exists for the shared Windows/WSL development host.
#
# It needs `go` and `git` on the host. It never starts a database and never
# measures, so it does not touch the benchmark lock.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$HERE/../.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "go not found on this host; run the Podman path instead: tools/queue/queue --build" >&2
  exit 2
fi
if ! command -v git.exe >/dev/null 2>&1 && [[ ! -x "/mnt/c/Program Files/Git/cmd/git.exe" ]]; then
  echo "no native Windows git found; nothing cross-platform to check here" >&2
  exit 0
fi

# A repository under /mnt/<drive>/ is addressable as C:/... by git.exe.
ROOT="$REPO_ROOT/.worktrees/queue-crossplatform"
mkdir -p "$ROOT"
trap 'rm -rf "$ROOT"' EXIT

cd "$REPO_ROOT/tools/queue"
ADS_QUEUE_TEST_ROOT="$ROOT" go test ./internal/app/ \
  -run 'TestWorktreeRegistrationIsPortable' -count=1 -v
