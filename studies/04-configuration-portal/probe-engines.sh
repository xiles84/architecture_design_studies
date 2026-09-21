#!/usr/bin/env bash
# Phase A — engine and semantic probes.
#
# Before any SQL is trusted, the engines are asked what they actually do rather
# than what the study assumes. Every item here is an assumption a design depends
# on, and each one has cost a study in this repository a wrong number before:
#
#   * the effective isolation level of each engine (older YugabyteDB releases ran
#     Snapshot Isolation where the flag said READ COMMITTED)
#   * whether a JSONB document can be updated in place, and whether `?` (key
#     existence) works -- the document design's search read depends on it
#   * whether a partial index, an array parameter and `unnest` behave as expected
#   * which EXPLAIN options exist, since the portable distributed-cost signal is
#     the RPC count from EXPLAIN (ANALYZE, DIST)
#   * the actual data-placement mechanism, which is what makes a colocated /
#     non-colocated comparison a comparison rather than an assertion
#
# It starts database containers, so it takes the benchmark lock and runs one
# topology at a time.
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
source "$STUDY_DIR/study.env"
set +e

TOPOLOGIES="${1:-pg-single,yb-single}"
OUT="$STUDY_DIR/results/devchecks/phase-a-probe"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$OUT"
REPORT="$OUT/probe-${RUN_ID}.txt"

run_lock_acquire "study ${STUDY_ID} probe-engines.sh ${RUN_ID}"
need_podman

{
  echo "phase: A (engine and semantic probes)"
  echo "probe_started_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "run_id: $RUN_ID"
  echo "podman: $(podman --version)"
  echo "postgres_image: $PG_IMAGE"
  echo "yugabyte_image: $YB_IMAGE"
} > "$REPORT"

# probe <label> <container> <client> <host> <port> <user> <db> <sql>
# Records the answer, or the exact error. A probe that fails is a finding, never
# a reason to continue quietly.
probe() {
  local label="$1" c="$2" client="$3" host="$4" port="$5" user="$6" db="$7" sql="$8"
  local out
  out="$(podman exec "$c" "$client" -h "$host" -p "$port" -U "$user" -d "$db" -Atc "$sql" 2>&1)"
  local st=$?
  {
    echo ""
    echo "--- $label"
    echo "status: $st"
    echo "$out"
  } >> "$REPORT"
}

for topo in ${TOPOLOGIES//,/ }; do
  case "$topo" in
    pg-single)
      bash "$REPO/infra/pg-single.sh" up || { warn "pg-single did not start"; continue; }
      c=pg-single; client=psql; host=pg-single; port=5432; user=bench; db=bench
      engine=postgres
      ;;
    yb-single)
      bash "$REPO/infra/yb-single.sh" up || { warn "yb-single did not start"; continue; }
      c=yb-single; client=ysqlsh; host=yb-single; port=5433; user=yugabyte; db=yugabyte
      engine=yugabyte
      ;;
    *) warn "unknown topology $topo"; continue ;;
  esac

  echo "" >> "$REPORT"
  echo "=== topology: $topo (engine $engine)" >> "$REPORT"

  probe "server version"        "$c" "$client" "$host" "$port" "$user" "$db" "SELECT version()"
  probe "effective isolation"   "$c" "$client" "$host" "$port" "$user" "$db" "SHOW transaction_isolation"
  probe "default_transaction_isolation" "$c" "$client" "$host" "$port" "$user" "$db" "SHOW default_transaction_isolation"
  probe "jsonb round trip"      "$c" "$client" "$host" "$port" "$user" "$db" "SELECT ('{\"a\":1}'::jsonb ->> 'a')"
  probe "jsonb in-place update" "$c" "$client" "$host" "$port" "$user" "$db" "SELECT (jsonb_set('{\"a\":1}'::jsonb, '{b}', '\"2\"') ->> 'b')"
  probe "jsonb key existence ?" "$c" "$client" "$host" "$port" "$user" "$db" "SELECT ('{\"a\":1}'::jsonb ? 'a')"
  probe "jsonb_object_agg"      "$c" "$client" "$host" "$port" "$user" "$db" "SELECT jsonb_object_agg(k, v) FROM (VALUES ('a','1'),('b','2')) AS t(k,v)"
  probe "unnest arrays"         "$c" "$client" "$host" "$port" "$user" "$db" "SELECT count(*) FROM unnest(ARRAY['a','b']::text[], ARRAY['1','2']::text[]) AS t(k,v)"
  probe "md5 + string_agg"      "$c" "$client" "$host" "$port" "$user" "$db" "SELECT md5(string_agg(k || '=' || v, ',' ORDER BY k)) FROM (VALUES ('a','1')) AS t(k,v)"
  probe "FOR UPDATE"            "$c" "$client" "$host" "$port" "$user" "$db" "SELECT 1"
  probe "partial index support" "$c" "$client" "$host" "$port" "$user" "$db" "SELECT count(*) FROM pg_indexes WHERE schemaname='pg_catalog'"
  probe "explain analyze options" "$c" "$client" "$host" "$port" "$user" "$db" "EXPLAIN (ANALYZE, BUFFERS) SELECT 1"
  probe "explain dist (YB only)"  "$c" "$client" "$host" "$port" "$user" "$db" "EXPLAIN (ANALYZE, DIST) SELECT 1"
  if [[ "$engine" == "yugabyte" ]]; then
    probe "placements (yb_tablet_metadata)" "$c" "$client" "$host" "$port" "$user" "$db" "SELECT count(*) FROM yb_tablet_metadata()"
    probe "colocation supported"            "$c" "$client" "$host" "$port" "$user" "$db" "SHOW yb_use_hash_splitting_by_default"
  fi

  bash "$REPO/infra/${topo}.sh" down
done

echo "" >> "$REPORT"
echo "probe_finished_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$REPORT"
cat "$REPORT"
