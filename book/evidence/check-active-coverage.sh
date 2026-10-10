#!/usr/bin/env bash
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLAIMS="${1:-$HERE/v6/claims.json}"
COVERAGE="${2:-$HERE/v6/coverage.json}"

[[ -f "$CLAIMS" && -f "$COVERAGE" ]] || {
  echo "coverage: missing registry or coverage file" >&2
  exit 1
}

EXPECTED="$(jq -cS '
  def freq(f):
    map(f) | sort | group_by(.) |
    map({key: .[0], value: length}) | from_entries;

  .claims | {
    total: length,
    strength: freq(.strength),
    evidence_form: freq(.evidence_form),
    control_coverage: freq(.control_coverage),
    replication_tier: freq(.replication.tier // "none"),
    prediction_verdict: freq(.prediction.verdict // "none"),
    review_state: freq(.review.state // "none"),
    attribution_withheld:
      (map(select(.attribution.withheld == true)) | length)
  }
' "$CLAIMS")"

ACTUAL="$(jq -cS '.active_claim_summary' "$COVERAGE")"

if [[ "$EXPECTED" != "$ACTUAL" ]]; then
  echo "coverage: active summary differs from claims.json" >&2
  echo "expected: $EXPECTED" >&2
  echo "recorded: $ACTUAL" >&2
  exit 1
fi

echo "coverage: active summary consistent"
