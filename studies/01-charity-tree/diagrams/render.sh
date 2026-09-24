#!/usr/bin/env bash
# Regenerate rendered/*.svg from the .puml sources with the repository-pinned
# renderer. All renderer logic lives in infra/diagram-render.sh (pinned image,
# bind-mount path handling, provenance record); this script only points it at
# this study's diagrams, so no host Java or PlantUML install is needed.
#
# Usage: ./render.sh [svg|png]
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$DIR/../../../infra/diagram-render.sh" \
  --format "${1:-svg}" \
  --record "$DIR/RENDERER.md" \
  "$DIR"
