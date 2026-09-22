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
TYPST_VERSION="$(podman run --rm "$BOOK_IMAGE" /bin/typst --version | sed 's/^typst //')"
EVIDENCE_DIGEST="$(sha256sum "$BOOK_DIR/evidence/v2/claims.json" | cut -d' ' -f1)"
SOURCE_HASH="$(
  find "$BOOK_DIR" -name '*.typ' -type f -print0 \
    | sort -z \
    | xargs -0 sha256sum \
    | sha256sum | cut -d' ' -f1
)"

log "compiling $OUT_REL (typst $TYPST_VERSION, commit ${COMMIT:0:12}, tree dirty=$DIRTY)"
mkdir -p "$BOOK_DIR/$(dirname "$OUT_REL")"

podman run --rm \
  -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" \
  -w "$BOOK_DIR" \
  "$BOOK_IMAGE" \
  /bin/typst compile --root "$BOOK_DIR" \
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
PAGES="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" "$BOOK_IMAGE" pdfinfo "$BOOK_DIR/$OUT_REL" | sed -n 's/^Pages:[[:space:]]*//p')"
FONTS_RAW="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" "$BOOK_IMAGE" pdffonts "$BOOK_DIR/$OUT_REL")"
FONTS_TOTAL="$(printf '%s\n' "$FONTS_RAW" | awk 'NR>2 && NF>0' | wc -l | tr -d ' ')"
FONTS_EMBEDDED="$(printf '%s\n' "$FONTS_RAW" | awk 'NR>2 && $4=="yes"' | wc -l | tr -d ' ')"
TEXT="$(podman run --rm -v "$(engine_path "$BOOK_DIR"):$BOOK_DIR" "$BOOK_IMAGE" pdftotext "$BOOK_DIR/$OUT_REL" -)"
CLAIM_COUNT="$(printf '%s' "$TEXT" | grep -oE 'v2-[a-z0-9-]+' | sort -u | wc -l | tr -d ' ')"
DIGEST_IN_PDF="false"
if printf '%s' "$TEXT" | grep -q "$EVIDENCE_DIGEST"; then DIGEST_IN_PDF="true"; fi
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
  "evidence_registry": "book/evidence/v2/claims.json",
  "evidence_digest_sha256": "$EVIDENCE_DIGEST",
  "source_tree_hash_sha256": "$SOURCE_HASH",
  "pdf_sha256": "$PDF_SHA256",
  "pdf_pages": ${PAGES:-0},
  "fonts_total": $FONTS_TOTAL,
  "fonts_embedded": $FONTS_EMBEDDED,
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
printf '  claims indexed: %s\n' "$CLAIM_COUNT"
printf '  evidence digest in pdf: %s\n' "$DIGEST_IN_PDF"
printf '  uri annotations: %s\n' "$LINK_ANNOTS"

if [[ "${PAGES:-0}" -lt 10 || "$FONTS_TOTAL" -eq 0 || "$FONTS_EMBEDDED" -ne "$FONTS_TOTAL" || "$DIGEST_IN_PDF" != "true" ]]; then
  log "verification FAILED: pages=$PAGES fonts=$FONTS_EMBEDDED/$FONTS_TOTAL digest_in_pdf=$DIGEST_IN_PDF"
  exit 1
fi
log "verified"
