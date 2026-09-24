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
#                      [--manifest dist/build-manifest.json]
#                      [--registry-version v4]
#
# --out and --manifest exist so a task that changes the sources can build and
# verify into a scratch path without committing a new dist/ artefact; the tracked
# PDF is regenerated once, from fully merged sources.

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
MANIFEST_REL="dist/build-manifest.json"
REGISTRY_VERSION="v5"

# The edition is stated once, in lib/config.typ, and read here. A second copy in this
# script is how the manifest came to say "Edition 1" beside a title page that said
# "Edition 2". The assertion below then checks that the value really reaches the page.
EDITION="$(sed -nE 's/^#let book-edition = "(.*)"$/\1/p' "$BOOK_DIR/lib/config.typ")"
[[ -n "$EDITION" ]] || die "could not read book-edition from lib/config.typ"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --out) OUT_REL="$2"; shift 2 ;;
    --manifest) MANIFEST_REL="$2"; shift 2 ;;
    --registry-version) REGISTRY_VERSION="$2"; shift 2 ;;
    *) printf 'unknown argument: %s\n' "$1" >&2; exit 2 ;;
  esac
done

case "$OUT_REL" in
  /*) die "--out must be a path inside book/: the build mounts only book/" ;;
esac
case "$MANIFEST_REL" in
  /*) die "--manifest must be a path inside book/: the build mounts only book/" ;;
esac

# The book build is not a measurement, so it does not take the benchmark lock.
# It does share the localhost/ads-book:1 image tag with every other worktree and
# agent, though, so two concurrent builds must not race on it. This lock is
# separate from ads-run-lock and never touches it.
BOOK_BUILD_LOCK_VOLUME="ads-book-build-lock"

book_build_lock_acquire() {
  local desc="$1"
  if podman volume create \
      --label "ads.holder=${desc}" \
      --label "ads.started=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
      --label "ads.worktree=$(hostpath "$REPO_ROOT")" \
      "$BOOK_BUILD_LOCK_VOLUME" >/dev/null 2>&1; then
    trap 'book_build_lock_release' EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    log "book build lock acquired: ${desc}"
    return 0
  fi
  warn "another book build holds ${BOOK_BUILD_LOCK_VOLUME}:"
  podman volume inspect "$BOOK_BUILD_LOCK_VOLUME" \
    --format '  holder:   {{index .Labels "ads.holder"}}
  started:  {{index .Labels "ads.started"}}
  worktree: {{index .Labels "ads.worktree"}}' 2>/dev/null || echo "  (lock volume vanished)"
  die "refusing to build concurrently: both builds would share the localhost/ads-book:1 image tag"
}

book_build_lock_release() {
  podman volume rm -f "$BOOK_BUILD_LOCK_VOLUME" >/dev/null 2>&1 || true
}

need_podman
book_build_lock_acquire "book build in ${REPO_ROOT} (pid $$)"

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
# Typst reports its version as "0.15.1 (unknown commit)" on a release build: the
# parenthetical is Typst's own build stamp, not this book's provenance. Record
# the raw string in the manifest and print only the version, so a reader cannot
# mistake it for one of our fields being unpopulated.
TYPST_VERSION_RAW="$(podman run --rm "$BOOK_IMAGE" --version | sed 's/^typst //')"
TYPST_VERSION="${TYPST_VERSION_RAW%% *}"

REGISTRY_REL="book/evidence/${REGISTRY_VERSION}/claims.json"
REGISTRY_FILE="$BOOK_DIR/evidence/${REGISTRY_VERSION}/claims.json"
[[ -f "$REGISTRY_FILE" ]] || die "registry ${REGISTRY_REL} not found"
EVIDENCE_DIGEST="$(sha256sum "$REGISTRY_FILE" | cut -d' ' -f1)"

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

# ---- book source gate ---------------------------------------------------
# Structural rules over the Typst sources and the registry: one active registry
# path, no retired claim rendered as evidence, no asset claiming its own figure
# number. Fatal, with the standing exceptions listed in check-waivers.txt so the
# debt is visible instead of silent.
bash "$BOOK_DIR/evidence/check.sh" --book "$BOOK_DIR" --registry-version "$REGISTRY_VERSION"

# The source-tree hash covers the Typst sources, the whole figure layer
# (registry, manifest, assets), the active evidence package and each referenced
# canonical figure source, so a changed registry or figure source cannot leave
# the hash unchanged.
SOURCE_HASH="$(
  {
    find "$BOOK_DIR" -name '*.typ' -type f
    find "$BOOK_DIR/assets" -type f
    find "$BOOK_DIR/evidence/${REGISTRY_VERSION}" -type f
    jq -r '.figures[].source' "$FIGURES_JSON" | while read -r src; do printf '%s\n' "$REPO_ROOT/$src"; done
  } | sort -u | xargs sha256sum | sha256sum | cut -d' ' -f1
)"

log "compiling $OUT_REL (typst $TYPST_VERSION, commit ${COMMIT:0:12}, tree dirty=$DIRTY, registry $REGISTRY_VERSION)"
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
    --input "registry_version=$REGISTRY_VERSION" \
    --input "release=1" \
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
TOTAL_CLAIMS="$(jq -r '.claims | length' "$REGISTRY_FILE")"
# Count only registry ids that actually render, so the number is a check that the
# evidence index compiled in, not a count of look-alike tokens.
CLAIM_COUNT="$(jq -r '.claims[].claim_id' "$REGISTRY_FILE" | while read -r id; do grep -q "$id" <<<"$TEXT" && echo "$id"; done | wc -l | tr -d ' ')"
DIGEST_IN_PDF="false"
# A here-string, not a pipe: grep -q exits on first match, and under pipefail the
# upstream printf's SIGPIPE would flip the test to false.
if grep -q "$EVIDENCE_DIGEST" <<<"$TEXT"; then DIGEST_IN_PDF="true"; fi
LINK_ANNOTS="$(grep -c '/URI' "$PDF_PATH" || true)"
# A release build passes every provenance input, so the unverified-build banner
# must not be on the page. This is what stops a preview build from being
# mistaken for a release artefact.
# A template token on the page means a source wrote `#name` somewhere the
# evaluator does not run: inside backticks (raw text), or inside a plain string
# argument. Edition 2's draft read "Evidence registry #registry-label" on page 2
# for exactly that reason, so this is fatal rather than a warning.
EDITION_IN_PDF="false"
if grep -qF "$EDITION" <<<"$TEXT"; then EDITION_IN_PDF="true"; fi

UNRESOLVED_IN_PDF="false"
UNRESOLVED_SEEN=""
for tok in '#registry-' '#build-' '#source-' '${' '{{' 'UNKNOWN_PLACEHOLDER'; do
  if grep -qF "$tok" <<<"$TEXT"; then
    UNRESOLVED_IN_PDF="true"
    UNRESOLVED_SEEN="${UNRESOLVED_SEEN}${UNRESOLVED_SEEN:+ }${tok}"
  fi
done

UNVERIFIED_BANNER="Unverified development build"
PREVIEW_IN_PDF="false"
if grep -q "$UNVERIFIED_BANNER" <<<"$TEXT"; then PREVIEW_IN_PDF="true"; fi

cat > "$BOOK_DIR/$MANIFEST_REL" <<JSON
{
  "artefact": "$OUT_REL",
  "book": "Data Architecture Reference",
  "edition": "$EDITION",
  "release": true,
  "source_commit": "$COMMIT",
  "source_describe": "$DESCRIBE",
  "source_dirty": $DIRTY,
  "built_at_utc": "$BUILT_AT",
  "typst_version": "$TYPST_VERSION",
  "typst_version_raw": "$TYPST_VERSION_RAW",
  "typst_image": "$TYPST_IMAGE",
  "typst_image_digest": "$TYPST_DIGEST",
  "evidence_registry": "$REGISTRY_REL",
  "evidence_registry_version": "$REGISTRY_VERSION",
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
  "unverified_banner_in_pdf": $PREVIEW_IN_PDF,
  "unresolved_tokens_in_pdf": $UNRESOLVED_IN_PDF,
  "edition_in_pdf": $EDITION_IN_PDF,
  "uri_annotations": $LINK_ANNOTS,
  "build_command": "book/build.sh (typst compile --root book main.typ --input commit=$COMMIT ... --input registry_version=$REGISTRY_VERSION --input release=1)"
}
JSON

log "wrote $MANIFEST_REL"
printf '  artefact:       %s\n' "$OUT_REL"
printf '  pdf sha256:     %s\n' "$PDF_SHA256"
printf '  pages:          %s\n' "${PAGES:-0}"
printf '  fonts embedded: %s/%s\n' "$FONTS_EMBEDDED" "$FONTS_TOTAL"
printf '  claims indexed: %s/%s\n' "$CLAIM_COUNT" "$TOTAL_CLAIMS"
printf '  evidence digest in pdf: %s\n' "$DIGEST_IN_PDF"
printf '  unverified banner in pdf: %s\n' "$PREVIEW_IN_PDF"
printf '  unresolved tokens in pdf: %s%s\n' "$UNRESOLVED_IN_PDF" "${UNRESOLVED_SEEN:+ ($UNRESOLVED_SEEN)}"
printf '  edition on page: %s (%s)\n' "$EDITION_IN_PDF" "$EDITION"
printf '  uri annotations: %s\n' "$LINK_ANNOTS"

if [[ "${PAGES:-0}" -lt 10 || "$FONTS_TOTAL" -eq 0 || "$FONTS_EMBEDDED" -ne "$FONTS_TOTAL" || "$DIGEST_IN_PDF" != "true" || "$CLAIM_COUNT" -ne "$TOTAL_CLAIMS" ]]; then
  log "verification FAILED: pages=$PAGES fonts=$FONTS_EMBEDDED/$FONTS_TOTAL digest_in_pdf=$DIGEST_IN_PDF claims=$CLAIM_COUNT/$TOTAL_CLAIMS"
  exit 1
fi
if [[ "$EDITION_IN_PDF" != "true" ]]; then
  log "verification FAILED: the manifest edition \"$EDITION\" does not appear in the PDF"
  exit 1
fi
if [[ "$UNRESOLVED_IN_PDF" != "false" ]]; then
  log "verification FAILED: unresolved template token(s) rendered in the PDF: ${UNRESOLVED_SEEN}"
  exit 1
fi
if [[ "$PREVIEW_IN_PDF" != "false" ]]; then
  log "verification FAILED: the unverified-build banner rendered in a release build"
  exit 1
fi
log "verified"
