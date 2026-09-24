#!/usr/bin/env bash
# Structural checks over the book's Typst sources and the active evidence
# package. `book/build.sh` runs this before compiling; it is also safe to run by
# hand (it needs only grep, jq and the repo).
#
# Usage: book/evidence/check.sh [--book DIR] [--registry-version v4]
#
# What it can and cannot do: these are *source* rules, so they catch the defects
# that made Edition 1 print a stale registry path and two retired gaps. They do
# not judge prose; a rule that cannot be decided mechanically is not here. Where
# a rule currently fails for a known reason, the exception is listed in
# check-waivers.txt, so the debt is visible rather than silent.

set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOK_DIR="$(cd "$HERE/.." && pwd)"
REGISTRY_VERSION="v4"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --book) BOOK_DIR="$2"; shift 2 ;;
    --registry-version) REGISTRY_VERSION="$2"; shift 2 ;;
    *) printf 'unknown argument: %s\n' "$1" >&2; exit 2 ;;
  esac
done

REGISTRY="$BOOK_DIR/evidence/$REGISTRY_VERSION/claims.json"
WAIVERS="$HERE/check-waivers.txt"
[[ -f "$REGISTRY" ]] || { printf 'check: registry not found: %s\n' "$REGISTRY" >&2; exit 1; }

failures=0
failure() {
  printf '  [FAIL] %s: %s\n' "$1" "$2"
  failures=$((failures + 1))
}

# A rule is waived only by an explicit, owned line in check-waivers.txt.
waived() {
  local rule="$1" subject="$2"
  [[ -f "$WAIVERS" ]] || return 1
  grep -qE "^${rule}[[:space:]]+${subject}([[:space:]]|$)" "$WAIVERS"
}

waiver_note() {
  local rule="$1" subject="$2"
  grep -E "^${rule}[[:space:]]+${subject}([[:space:]]|$)" "$WAIVERS" | head -1 |
    sed -E 's/^[^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+—[[:space:]]*//'
}

printf 'check: book sources against %s\n' "book/evidence/$REGISTRY_VERSION/claims.json"

# --- registry-literal ----------------------------------------------------
# Exactly one place may know the registry version: the build input, derived into
# a path. Edition 1 printed "book/evidence/v2/claims.json" in a chapter while the
# cover said v3, and a reader who checks provenance notices that first.
while IFS= read -r hit; do
  [[ -z "$hit" ]] && continue
  file="${hit%%:*}"
  if waived "registry-literal" "$(basename "$file")"; then
    printf '  [waived] registry-literal %s — %s\n' "$(basename "$file")" "$(waiver_note registry-literal "$(basename "$file")")"
    continue
  fi
  failure "registry-literal" "$hit names a fixed registry version; derive it from the build input instead"
done < <(grep -rEn 'evidence/v[0-9]' "$BOOK_DIR" --include='*.typ' --include='*.md' --include='*.sh' 2>/dev/null | grep -v '^'"$HERE"'/' || true)

# --- quoted claim ids ----------------------------------------------------
# Every id a chapter renders must exist in the active package and must not be
# retired or superseded. The Typst loader also panics on an unknown id; this runs
# earlier and names the retired case specifically.
quoted="$(grep -rhoE '(registry-card|claim-by-id)\("[^"]+"\)' "$BOOK_DIR" --include='*.typ' 2>/dev/null |
  sed -E 's/.*\("([^"]+)"\)/\1/' | sort -u || true)"
while IFS= read -r id; do
  [[ -z "$id" ]] && continue
  if ! jq -e --arg id "$id" '.claims[] | select(.claim_id == $id)' "$REGISTRY" >/dev/null; then
    if jq -e --arg id "$id" '.retired_predecessor_claims[] | select(.claim_id == $id)' "$REGISTRY" >/dev/null; then
      failure "retired-claim-quoted" "$id is retired in $REGISTRY_VERSION; cite its successor instead"
    else
      failure "unknown-claim-id" "$id is not a claim in $REGISTRY_VERSION"
    fi
    continue
  fi
  status="$(jq -r --arg id "$id" '.claims[] | select(.claim_id == $id) | .status // "active"' "$REGISTRY")"
  if [[ "$status" == "retired" || "$status" == "superseded" ]]; then
    failure "retired-claim-quoted" "$id has status $status and may not be rendered as evidence"
  fi
done <<<"$quoted"

# --- gap-kind ------------------------------------------------------------
# A gap must say what sort of gap it is, or the book cannot badge a schema
# limitation ("no representation can answer this") differently from a regime
# nobody ran.
while IFS= read -r id; do
  [[ -z "$id" ]] && continue
  failure "gap-kind" "$id is a gap with no gap_kind"
done < <(jq -r '.claims[] | select(.strength == "gap") | select((.gap_kind // []) | length == 0) | .claim_id' "$REGISTRY")

# --- figure ownership ----------------------------------------------------
# Typst owns figure numbering. An asset that draws its own "Figure 3" produces
# two competing numbers on one page, which is what Edition 1 shipped.
while IFS= read -r asset; do
  [[ -z "$asset" ]] && continue
  name="$(basename "$asset")"
  if waived "figure-number" "$name"; then
    printf '  [waived] figure-number %s — %s\n' "$name" "$(waiver_note figure-number "$name")"
    continue
  fi
  failure "figure-number" "$name embeds its own figure number; Typst numbers figures"
done < <(grep -rlE 'Figure [0-9]' "$BOOK_DIR/assets" --include='*.svg' 2>/dev/null || true)

# --- orphan-figure -------------------------------------------------------
# Every registered figure must be embedded by some Typst source, so a figure
# cannot sit in the registry (and in the manifest, and in the source hash)
# without appearing in the book.
while IFS= read -r asset; do
  [[ -z "$asset" ]] && continue
  if ! grep -rqF "$asset" "$BOOK_DIR" --include='*.typ' 2>/dev/null; then
    failure "orphan-figure" "$asset is registered but no Typst source embeds it"
  fi
done < <(jq -r '.figures[].asset' "$BOOK_DIR/assets/figures.json" 2>/dev/null || true)

if [[ "$failures" -gt 0 ]]; then
  printf 'check: %d failure(s)\n' "$failures" >&2
  exit 1
fi
printf 'check: OK\n'
