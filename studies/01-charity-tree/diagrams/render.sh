#!/usr/bin/env bash
# Regenerate rendered/*.svg from the .puml sources.
#
# Runs PlantUML in a container, so no Java and no PlantUML install is needed on
# the host -- same rule as the rest of this repository.
#
# Usage: ./render.sh [svg|png]
set -euo pipefail
export MSYS_NO_PATHCONV=1

FMT="${1:-svg}"
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE="docker.io/plantuml/plantuml:latest"

mkdir -p "$DIR/rendered"
# _style.puml is an include, not a diagram, so it is listed explicitly nowhere.
mapfile -t files < <(cd "$DIR" && ls *.puml | grep -v '^_')

podman run --rm -v "$DIR:/data" -w /data "$IMAGE" "-t$FMT" -o /data/rendered "${files[@]}"
echo "rendered ${#files[@]} diagrams to $DIR/rendered as .$FMT"
