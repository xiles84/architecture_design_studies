#!/usr/bin/env bash
# Build the data-architecture reference book in Podman and write its manifest.
#
# Everything happens in the pinned image: nothing on the host needs Typst, a
# font, or a PDF tool. The script records actual values only — commit, describe,
# dirty state, build clock, Typst version, image digest, evidence-registry
# digest, source-tree hash, PDF hash, page count and embedded fonts — and it
# refuses to claim a number it did not measure.
#
# Usage: book/build.sh [--out dist/data-architecture-reference.pdf]

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$HERE/.." && pwd)"
BOOK_DIR="$HERE"

# shellcheck source=../infra/lib.sh
source "$REPO_ROOT/infra/lib.sh"

# engine_path: a path the podman binary accepts on this platform (Git Bash:
# cygpath form; WSL: the Windows form; native Linux/macOS: unchanged).
engine_path() {
  if command -v cygpath >/dev/null 2>&1; then
    hostpath "$1"
  else
    winpath "$1"
  fi
}

BOOK_IMAGE="localhost/ads-book:1"
TYPST_DIGEST="sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f"
TYPST_IMAGE="ghcr.io/typst/typst@${TYPST_DIGEST}"

OUT_REL="dist/data-architecture-reference.pdf"
if [[ "${1:-}" == "--out" ]]; then
  OUT_REL="$2"
fi
MANIFEST_REL="dist/build-manifest.json"

need_podman

if ! podman image exists "$BOOK_IMAGE"; then
  log "building $BOOK_IMAGE from book/Containerfile"
  podman build -t "$BOOK_IMAGE" -f "$(engine_path "$BOOK_DIR/Containerfile")" "$(engine_path "$BOOK_DIR")"
fi

# ---- actual build inputs ------------------------------------------------
COMMIT="$(git -C "$REPO_ROOT" rev-parse HEAD)"
DESCRIBE="$(git -C "$REPO_ROOT" describe --tags --always 2>/dev/null || echo "$COMMIT")"
if [[ -n "$(git -C "$REPO_ROOT" status --porcelain --untracked-files=no)" ]]; then
  DIRTY="true"
else
  DIRTY="false"
fi
BUILT_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
TYPST_VERSION="$(podman run --rm "$BOOK_IMAGE" --version | sed 's/^typst //')"
EVIDENCE_DIGEST="$(sha256sum "$BOOK_DIR/evidence/v3/claims.json" | cut -d' ' -f1)"

# ---- figure provenance gate ---------------------------------------------
# Every embedded figure is a generated asset exported from one canonical study
# .puml source. This aborts the build on a missing asset, a source that changed
# without a re-export, or an asset that no longer matches a fresh pinned render.
FIGURES_JSON="$BOOK_DIR/assets/figures.json"
FIGURES_MANIFEST="$BOOK_DIR/assets/figure-manifest.json"
[[ -f "$FIGURES_JSON" ]] || { log "figure registry missing: $FIGURES_JSON"; exit 1; }
bash "$BOOK_DIR/assets/export.sh" --check
FIGURES_TOTAL="$(jq '.figures | length' "$FIGURES_JSON")"
FIGURES_MANIFEST_SHA="$(sha256sum "$FIGURES_MANIFEST" | cut -d' ' -f1)"

# The source-tree hash covers the Typst sources, the whole figure layer
# (registry, manifest, assets) and each referenced canonical figure source, so a
# changed figure source cannot leave the hash unchanged.
SOURCE_HASH="$(
  {
    find "$BOOK_DIR" -name '*.typ' -type f
    find "$BOOK_DIR/assets" -type f
    jq -r '.figures[].source' "$FIGURES_JSON" | while read -r src; do printf '%s\n' "$REPO_ROOT/$src"; done
  } | sort -u | xargs sha256sum | sha256sum | cut -d' ' -f1
)"

log "compiling $OUT_REL (typst $TYPST_VERSION, commit ${COMMIT:0:12}, tree dirty=$DIRTY)"
mkdir -p "$BOOK_DIR/$(dirname "$OUT_REL")"

podman run --rm \
  -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" \
  -w "$BOOK_DIR" \
  "$BOOK_IMAGE" \
  compile --root "$BOOK_DIR" \
    --input "commit=$COMMIT" \
    --input "describe=$DESCRIBE" \
    --input "dirty=$DIRTY" \
    --input "built_at=$BUILT_AT" \
    --input "typst_version=$TYPST_VERSION" \
    --input "typst_digest=$TYPST_DIGEST" \
    --input "evidence_digest=$EVIDENCE_DIGEST" \
    --input "source_hash=$SOURCE_HASH" \
    main.typ "$OUT_REL"

# ---- verify the rendered PDF -------------------------------------------
# The verification runs in the same pinned image; a claim in the manifest is a
# value these commands actually printed.
PDF_PATH="$BOOK_DIR/$OUT_REL"
PDF_SHA256="$(sha256sum "$PDF_PATH" | cut -d' ' -f1)"
PAGES="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" --entrypoint pdfinfo "$BOOK_IMAGE" "$BOOK_DIR/$OUT_REL" | sed -n 's/^Pages:[[:space:]]*//p')"
FONTS_RAW="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" --entrypoint pdffonts "$BOOK_IMAGE" "$BOOK_DIR/$OUT_REL")"
FONTS_TOTAL="$(printf '%s\n' "$FONTS_RAW" | awk 'NR>2 && NF>0' | wc -l | tr -d ' ')"
# pdffonts' `emb` column moves with the font type; the embedded triple is the
# stable marker ("yes yes yes" = embedded, subset, unicode).
FONTS_EMBEDDED="$(printf '%s\n' "$FONTS_RAW" | awk 'NR>2 && /yes[[:space:]]+yes[[:space:]]+yes/' | wc -l | tr -d ' ')"
TEXT="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" --entrypoint pdftotext "$BOOK_IMAGE" "$BOOK_DIR/$OUT_REL" -)"
TOTAL_CLAIMS="$(jq -r '.claims | length' "$BOOK_DIR/evidence/v3/claims.json")"
# Count only registry ids that actually render, so the number is a check that the
# evidence index compiled in, not a count of look-alike tokens.
CLAIM_COUNT="$(jq -r '.claims[].claim_id' "$BOOK_DIR/evidence/v3/claims.json" | while read -r id; do grep -q "$id" <<<"$TEXT" && echo "$id"; done | wc -l | tr -d ' ')"
DIGEST_IN_PDF="false"
# A here-string, not a pipe: grep -q exits on first match, and under pipefail the
# upstream printf's SIGPIPE would flip the test to false.
if grep -q "$EVIDENCE_DIGEST" <<<"$TEXT"; then DIGEST_IN_PDF="true"; fi
LINK_ANNOTS="$(grep -c '/URI' "$PDF_PATH" || true)"

cat > "$BOOK_DIR/$MANIFEST_REL" <<JSON
{
  "artefact": "$OUT_REL",
  "book": "Data Architecture Reference",
  "edition": "Edition 1 (draft)",
  "source_commit": "$COMMIT",
  "source_describe": "$DESCRIBE",
  "source_dirty": $DIRTY,
  "built_at_utc": "$BUILT_AT",
  "typst_version": "$TYPST_VERSION",
  "typst_image": "$TYPST_IMAGE",
  "typst_image_digest": "$TYPST_DIGEST",
  "evidence_registry": "book/evidence/v3/claims.json",
  "evidence_digest_sha256": "$EVIDENCE_DIGEST",
  "figure_registry": "book/assets/figures.json",
  "figure_manifest": "book/assets/figure-manifest.json",
  "figure_manifest_sha256": "$FIGURES_MANIFEST_SHA",
  "figures_total": $FIGURES_TOTAL,
  "source_tree_hash_sha256": "$SOURCE_HASH",
  "pdf_sha256": "$PDF_SHA256",
  "pdf_pages": ${PAGES:-0},
  "fonts_total": $FONTS_TOTAL,
  "fonts_embedded": $FONTS_EMBEDDED,
  "claims_total": $TOTAL_CLAIMS,
  "claims_indexed": $CLAIM_COUNT,
  "evidence_digest_in_pdf": $DIGEST_IN_PDF,
  "uri_annotations": $LINK_ANNOTS,
  "build_command": "book/build.sh (typst compile --root book main.typ --input commit=$COMMIT ... --input evidence_digest=$EVIDENCE_DIGEST)"
}
JSON

log "wrote $MANIFEST_REL"
printf '  artefact:       %s\n' "$OUT_REL"
printf '  pdf sha256:     %s\n' "$PDF_SHA256"
printf '  pages:          %s\n' "${PAGES:-0}"
printf '  fonts embedded: %s/%s\n' "$FONTS_EMBEDDED" "$FONTS_TOTAL"
printf '  claims indexed: %s/%s\n' "$CLAIM_COUNT" "$TOTAL_CLAIMS"
printf '  evidence digest in pdf: %s\n' "$DIGEST_IN_PDF"
printf '  uri annotations: %s\n' "$LINK_ANNOTS"

if [[ "${PAGES:-0}" -lt 10 || "$FONTS_TOTAL" -eq 0 || "$FONTS_EMBEDDED" -ne "$FONTS_TOTAL" || "$DIGEST_IN_PDF" != "true" || "$CLAIM_COUNT" -ne "$TOTAL_CLAIMS" ]]; then
  log "verification FAILED: pages=$PAGES fonts=$FONTS_EMBEDDED/$FONTS_TOTAL digest_in_pdf=$DIGEST_IN_PDF claims=$CLAIM_COUNT/$TOTAL_CLAIMS"
  exit 1
fi
log "verified"
