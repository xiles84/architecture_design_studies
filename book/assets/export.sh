#!/usr/bin/env bash
# Figure export and provenance gate for the book.
#
# The book embeds generated SVG assets; it never holds an independently editable
# copy of a figure. Each figure's canonical editable source is a PlantUML file
# kept in the study that owns it (book/assets/figures.json). This script renders
# every source with the pinned renderer (infra/diagram-render.sh), then either
# writes the assets and manifest or verifies them:
#
#   book/assets/export.sh --write   # re-render, refresh assets + figure-manifest.json
#   book/assets/export.sh --check   # fail on a missing asset, changed source or stale render
#
# `book/build.sh` runs --check before compiling, so a figure whose source moved
# without a re-export aborts the build. No database is started; the only
# container used is the pinned PlantUML renderer.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"          # book/assets
BOOK_DIR="$(cd "$HERE/.." && pwd)"                             # book
REPO_ROOT="$(cd "$BOOK_DIR/.." && pwd)"                        # repository root
# shellcheck source=../../infra/lib.sh
source "$REPO_ROOT/infra/lib.sh"
: "${PLANTUML_IMAGE:?PLANTUML_IMAGE is not set (expected in infra/versions.env)}"

FIGURES="$HERE/figures.json"
MANIFEST="$HERE/figure-manifest.json"

MODE="${1:---check}"
case "$MODE" in
  --check|--write) ;;
  -h|--help) sed -n '2,18p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
  *) die "usage: book/assets/export.sh [--check|--write]" ;;
esac

[[ -f "$FIGURES" ]] || die "figure registry missing: $FIGURES"
need_podman

RENDERER_DIGEST="$(podman image inspect "$PLANTUML_IMAGE" --format '{{.Digest}}')"
RENDERER_ID="$(podman image inspect "$PLANTUML_IMAGE" --format '{{.Id}}')"

# Postprocess step, recorded in the manifest: a change to this step invalidates
# every asset, exactly as a renderer change does.
#
# PlantUML emits textLength and lengthAdjust="spacing" on every <text> element, and
# Typst honours them, so each label is stretched to the renderer's guessed width
# instead of being laid out in the real font. Every asset relied on this, which is
# why a long-labelled diagram looked letter-spaced while the short ones looked
# acceptable. Stripping it at export, rather than correcting one figure, is what
# makes the fix hold for figures nobody has drawn yet.
POSTPROCESS="strip-textlength-v1"

postprocess_svg() {
  local f="$1" before after
  before="$(grep -cE 'textLength=|lengthAdjust=' "$f" || true)"
  sed -i -E 's/ textLength="[^"]*"//g; s/ lengthAdjust="[^"]*"//g' "$f"
  after="$(grep -cE 'textLength=|lengthAdjust=' "$f" || true)"
  [[ "$after" == "0" ]] || die "postprocess left $after stretching attribute(s) in $f"
  printf '  postprocess: %s on %s (%s line(s) stripped)\n' "$POSTPROCESS" "$(basename "$f")" "$before" >&2
}

# A Windows-visible temp dir matters: the podman client cannot bind-mount a WSL
# /tmp path. book/dist already exists for generated output.
TMP="$BOOK_DIR/dist/.figures-tmp"
rm -rf "$TMP"
mkdir -p "$TMP"
trap 'rm -rf "$TMP"' EXIT

# One render per diagram directory, however many figures it holds.
declare -A DIR_RENDER
render_dir() {
  local dir="$1" key out
  key="$(printf '%s' "$dir" | sha256sum | cut -c1-12)"
  out="$TMP/dir-$key"
  if [[ ! -d "$out" ]]; then
    mkdir -p "$out"
    "$REPO_ROOT/infra/diagram-render.sh" --to "$out" "$dir" >/dev/null
  fi
  printf '%s' "$out"
}

COUNT="$(jq '.figures | length' "$FIGURES")"
ENTRIES="$TMP/entries.ndjson"
: > "$ENTRIES"

i=0
while [[ "$i" -lt "$COUNT" ]]; do
  id="$(jq -r ".figures[$i].id" "$FIGURES")"
  src="$(jq -r ".figures[$i].source" "$FIGURES")"
  asset="$(jq -r ".figures[$i].asset" "$FIGURES")"
  src_abs="$REPO_ROOT/$src"
  [[ -f "$src_abs" ]] || die "figure $id: canonical source does not exist: $src"

  out="$(render_dir "$(dirname "$src_abs")")"
  stem="$(basename "$src" .puml)"
  rendered="$out/$stem.svg"
  [[ -f "$rendered" ]] || die "figure $id: the renderer produced no $stem.svg from $src"

  postprocess_svg "$rendered"

  src_sha="$(sha256sum "$src_abs" | cut -d' ' -f1)"
  # HEAD:$src is the blob at the last commit, which is not the file's content when
  # a change has not been committed yet — and a manifest written in a dirty tree
  # records a revision that --check can never reproduce, because by then the commit
  # exists. Record the committed blob only when it actually describes the file.
  # Purely informational: which commit last carried this content, or uncommitted.
  # It is deliberately excluded from the staleness comparison below, because the
  # question "is this asset stale?" is answered by content hashes (source_sha256,
  # asset_sha256) and by the renderer and postprocess identity. A revision is a
  # property of the repository, not of the artefact, and using it here made a
  # manifest written before its own commit fail the check after it.
  src_rev="$(git -C "$REPO_ROOT" rev-parse "HEAD:$src" 2>/dev/null || echo uncommitted)"
  if [[ "$src_rev" != "uncommitted" ]] && [[ "$(git -C "$REPO_ROOT" hash-object "$src_abs")" != "$src_rev" ]]; then
    src_rev="uncommitted"
  fi
  svg_sha="$(sha256sum "$rendered" | cut -d' ' -f1)"
  bytes="$(wc -c < "$rendered" | tr -d ' ')"

  if [[ "$MODE" == "--write" ]]; then
    cp "$rendered" "$HERE/$asset"
  else
    [[ -f "$HERE/$asset" ]] || die "figure $id: asset missing: book/assets/$asset (run book/assets/export.sh --write)"
    cmp -s "$HERE/$asset" "$rendered" || die "figure $id: book/assets/$asset is stale relative to its source $src (run book/assets/export.sh --write)"
  fi

  jq -n \
    --arg id "$id" --arg source "$src" --arg source_revision "$src_rev" \
    --arg source_sha256 "$src_sha" --arg renderer_image "$PLANTUML_IMAGE" \
    --arg renderer_digest "$RENDERER_DIGEST" --arg renderer_id "$RENDERER_ID" \
    --arg asset "$asset" --arg asset_sha256 "$svg_sha" --argjson bytes "$bytes" \
    --arg postprocess "$POSTPROCESS" \
    '{id:$id, source:$source, source_revision:$source_revision, source_sha256:$source_sha256,
      renderer_image:$renderer_image, renderer_digest:$renderer_digest, renderer_id:$renderer_id,
      postprocess:$postprocess,
      asset:$asset, asset_sha256:$asset_sha256, bytes:$bytes}' >> "$ENTRIES"

  i=$((i + 1))
done

FRESH="$TMP/figure-manifest.json"
jq -s --arg ri "$PLANTUML_IMAGE" --arg rd "$RENDERER_DIGEST" --arg rid "$RENDERER_ID" \
  --arg pp "$POSTPROCESS" \
  '{schema_version:2, renderer_image:$ri, renderer_digest:$rd, renderer_id:$rid,
    postprocess:$pp, figures:.}' \
  "$ENTRIES" > "$FRESH"

if [[ "$MODE" == "--write" ]]; then
  cp "$FRESH" "$MANIFEST"
  log "wrote $MANIFEST ($COUNT figure(s) from $FIGURES)"
  exit 0
fi

[[ -f "$MANIFEST" ]] || die "figure manifest missing: book/assets/figure-manifest.json (run book/assets/export.sh --write)"
if ! diff -q <(jq -S '.figures | map(del(.source_revision))' "$MANIFEST") <(jq -S '.figures | map(del(.source_revision))' "$FRESH") >/dev/null; then
  diff -u <(jq -S '.figures | map(del(.source_revision))' "$MANIFEST") <(jq -S '.figures | map(del(.source_revision))' "$FRESH") >&2 || true
  die "figure manifest is out of date (run book/assets/export.sh --write)"
fi

log "figure assets and manifest verified ($COUNT figure(s), renderer digest ${RENDERER_DIGEST:0:19}…)"
