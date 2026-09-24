#!/usr/bin/env bash
# Render one study's PlantUML diagram set with the repository-pinned renderer.
#
# Every study's diagrams/render.sh delegates here, so the renderer version, the
# bind-mount path handling and the renderer-provenance record have exactly one
# implementation. The image is pinned by digest in infra/versions.env; nothing
# runs on the host (no Java, no PlantUML install), and no database is started.
#
# Usage:
#   infra/diagram-render.sh [--format svg|png] [--to DIR] [--record FILE] <diagrams-dir>
#
#   <diagrams-dir>  directory holding *.puml sources; files whose name starts
#                   with "_" are includes, not diagrams, and are skipped
#   --format FMT    output format: svg (default) or png
#   --to DIR        output directory; default <diagrams-dir>/rendered
#   --record FILE   also write a renderer-provenance markdown file

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib.sh
source "$HERE/lib.sh"
: "${PLANTUML_IMAGE:?PLANTUML_IMAGE is not set (expected in infra/versions.env)}"

# engine_path: the form of a bind-mount source the podman binary accepts. Git
# Bash wants the cygpath form; WSL has no cygpath and needs the C:/ form; native
# Linux and macOS are unchanged. Same policy as book/build.sh and the queue
# wrapper -- an unconverted WSL path mounts an empty directory without erroring.
engine_path() {
  if command -v cygpath >/dev/null 2>&1; then
    hostpath "$1"
  else
    winpath "$1"
  fi
}

usage() { sed -n '2,14p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; }

FORMAT="svg"
TO=""
RECORD=""
DIAGRAM_DIR=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --format) FORMAT="${2:?--format needs a value}"; shift 2 ;;
    --to)     TO="${2:?--to needs a value}"; shift 2 ;;
    --record) RECORD="${2:?--record needs a value}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    -*) die "unknown option: $1" ;;
    *) [[ -z "$DIAGRAM_DIR" ]] || die "unexpected argument: $1"; DIAGRAM_DIR="$1"; shift ;;
  esac
done
case "$FORMAT" in
  svg|png) ;;
  *) die "unsupported format: $FORMAT (svg or png)" ;;
esac
[[ -n "$DIAGRAM_DIR" ]] || { usage >&2; die "missing <diagrams-dir>"; }
[[ -d "$DIAGRAM_DIR" ]] || die "not a directory: $DIAGRAM_DIR"
DIAGRAM_DIR="$(cd "$DIAGRAM_DIR" && pwd)"
[[ -n "$TO" ]] || TO="$DIAGRAM_DIR/rendered"
mkdir -p "$TO"
TO="$(cd "$TO" && pwd)"

# _style.puml (and any other _-prefixed file) is an include, not a diagram.
mapfile -t files < <(cd "$DIAGRAM_DIR" && ls *.puml 2>/dev/null | grep -v '^_' || true)
[[ ${#files[@]} -gt 0 ]] || die "no .puml sources in $DIAGRAM_DIR"

need_podman
log "rendering ${#files[@]} diagram(s) from $DIAGRAM_DIR with $PLANTUML_IMAGE"
podman run --rm \
  -v "$(engine_path "$DIAGRAM_DIR"):/data" \
  -v "$(engine_path "$TO"):/out" \
  -w /data \
  "$PLANTUML_IMAGE" \
  "-t$FORMAT" -o /out "${files[@]}"

if [[ -n "$RECORD" ]]; then
  IMAGE_ID="$(podman image inspect "$PLANTUML_IMAGE" --format '{{.Id}}')"
  IMAGE_DIGEST="$(podman image inspect "$PLANTUML_IMAGE" --format '{{.Digest}}')"
  # No pipe to head: under `set -o pipefail` the podman side of a pipe that is
  # closed early dies with SIGPIPE, which would append a spurious "unknown".
  PUML_RAW="$(podman run --rm "$PLANTUML_IMAGE" -version 2>&1 || true)"
  PUML_VERSION="${PUML_RAW%%$'\n'*}"
  [[ -n "$PUML_VERSION" ]] || PUML_VERSION="unknown"
  {
    printf '# Diagram renderer provenance\n\n'
    printf 'Written by `infra/diagram-render.sh`; do not edit by hand.\n\n'
    printf '| Field | Value |\n|---|---|\n'
    printf '| Image | `%s` |\n' "$PLANTUML_IMAGE"
    printf '| Resolved digest | `%s` |\n' "$IMAGE_DIGEST"
    printf '| Resolved image id | `%s` |\n' "$IMAGE_ID"
    printf '| Renderer version | %s |\n' "$PUML_VERSION"
    printf '| Format | `%s` |\n' "$FORMAT"
    printf '| Sources | %d (`*.puml`, excluding `_`-prefixed includes) |\n' "${#files[@]}"
    printf '| Rendered at (UTC) | %s |\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    printf '\nOutput lives in `rendered/`. The image is pinned once in `infra/versions.env`\n'
    printf '(entry `PLANTUML_IMAGE`); the digest above is what this render actually used.\n'
  } > "$RECORD"
  log "wrote renderer provenance $RECORD"
fi

log "rendered ${#files[@]} diagram(s) to $TO as .$FORMAT"
