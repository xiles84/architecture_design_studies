#!/usr/bin/env bash
# Experiment E — the regimes in which each design's limitation stops mattering.
#
#   ./run-regimes.sh
#
# Every design in this study loses somewhere. The survey measures the general case
# (unbounded donation history, ten large charities, a mixed dashboard-and-portal
# workload), and in that case several designs look plainly worse. But a conclusion
# of the form "D6 is bad" is incomplete in a way that does real damage: it leads
# people whose domain happens to avoid D6's weakness to reject the design that
# would have served them best.
#
# So for each limitation, this runs the regime that removes it, and the comparison
# reports ask one question: does the regime change which design wins?
#
#   limitation                                   regime that removes it     stage
#   -------------------------------------------  -------------------------  -----
#   D6 rewrites the whole array per append;      donations per donor         1
#     D9's cache covers only 20; D4's delete     capped at 20 (same volume)
#     trigger recomputes over the whole history;
#     D7's per-donor tablet grows unbounded
#   D4/D5 serialise every write in a charity     1 000 small charities       2
#     on ONE rollup row; D6 unnests a whole      instead of 10 large ones
#     charity for charity-scoped questions       (same people, same volume)
#   D6 must unnest every array for any           donor-portal workload that  3
#     question that crosses donors               never crosses donors
#   D1 has no secondary indexes                  a small dataset             4
#
# Regimes already covered by earlier experiments, recorded here so the map is
# complete:
#   D4's hot row at LOW write concurrency       -> experiment C, 1 and 4 writers
#   D5 maintains rollups on insert only         -> experiment C (insert path); fine
#                                                  for an append-only ledger
#   D8 enforces no referential integrity        -> not a performance regime; it is
#                                                  acceptable when one writer owns
#                                                  the invariant, which no benchmark
#                                                  can establish
#
# Unbounded remains the study's primary profile: a charity's donation history
# genuinely grows. These regimes are alternatives, not replacements.
#
# Stage 0 comes first because it replaces invalid data: the earlier repeated-trials
# and parent-write passes measured deletes with a harness bug (see
# LESSONS_LEARNED.md, "Finite pools").
set -uo pipefail

STUDY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$STUDY_DIR/../.." && pwd)"
source "$REPO/infra/lib.sh"
set +e
run_lock_acquire "study 01 run-regimes.sh $*"

DATE="${REGIME_DATE:-$(date -u +%Y%m%d)}"
ALL_WRITES="insert,update,delete,update_person,delete_person"
S="$STUDY_DIR"

stage() { log "==================== $* ===================="; }

# --- Stage 0: isolated writes, fixed delete measurement ----------------------
stage "0/5 isolated writes (replaces the invalid trials and parent-write passes)"
bash "$S/run-study.sh" --scale small --duration 6s --warmup 2s --trials 3 \
  --topologies pg-single,yb-cluster3 \
  --designs d3_flattened_fk,d8_flattened_nofk,d4_rollup_trigger,d6_embedded_jsonb,d9_embedded_hybrid \
  --write-ops "$ALL_WRITES" --extra "-queries __none__" \
  --run-id "$DATE-writes-isolated"

# --- Stage 1: bounded donation history ----------------------------------------
for prof in unbounded cap20; do
  extra=""
  [[ "$prof" == "cap20" ]] && extra="-max-per-person 20"
  stage "1/5 donation history: $prof"
  bash "$S/run-study.sh" --scale small --duration 8s --warmup 3s \
    --topologies pg-single \
    --designs d3_flattened_fk,d4_rollup_trigger,d6_embedded_jsonb,d9_embedded_hybrid \
    --write-ops "$ALL_WRITES" --extra "$extra" --run-id "$DATE-history-$prof"
  bash "$S/run-study.sh" --scale small --duration 8s --warmup 3s \
    --topologies yb-cluster3 \
    --designs d3_flattened_fk,d6_embedded_jsonb,d7_yb_child_colocated \
    --write-ops "insert,update_person,delete_person" --extra "$extra" --run-id "$DATE-history-$prof"
done

# --- Stage 2: many small charities --------------------------------------------
stage "2/5 top of the tree: 1000 charities"
bash "$S/run-study.sh" --scale small --duration 8s --warmup 3s \
  --topologies pg-single \
  --designs d3_flattened_fk,d4_rollup_trigger,d6_embedded_jsonb,d9_embedded_hybrid \
  --write-ops "$ALL_WRITES" --extra "-charities 1000" --run-id "$DATE-top-1000charities"

# The hot-row question is a concurrency question, so it is also asked the way
# experiment C asks it -- at 8 and 32 writers, for both charity counts, in the
# same pass so the two are directly comparable.
for ch in 10 1000; do
  stage "2/5 contention with $ch charities"
  bash "$S/run-concurrency.sh" --topology pg-single --scale small --duration 8s \
    --conns-list 8,32 --strategies forupdate,optimistic --isolations rc \
    --extra "-charities $ch" --run-id "$DATE-contention-${ch}charities"
done

# --- Stage 3: a workload that never crosses donors ----------------------------
for prof in unbounded cap20; do
  extra="-read-mix portal"
  [[ "$prof" == "cap20" ]] && extra="$extra -max-per-person 20"
  stage "3/5 donor-portal workload: $prof"
  bash "$S/run-mixed.sh" --topology pg-single --scale small --duration 10s \
    --splits 8:0,6:2 --designs d3_flattened_fk,d6_embedded_jsonb,d9_embedded_hybrid \
    --extra "$extra" --run-id "$DATE-portal-$prof"
done

# --- Stage 4: a small dataset ---------------------------------------------------
for sc in small tiny; do
  stage "4/5 dataset size: $sc"
  bash "$S/run-study.sh" --scale "$sc" --duration 8s --warmup 3s \
    --topologies pg-single --designs d1_normalized_minimal,d2_normalized_indexed \
    --write-ops "$ALL_WRITES" --run-id "$DATE-size-$sc"
done

# --- Stage 5: reports -----------------------------------------------------------
stage "5/5 reports"
rep() {
  podman run --rm \
    -v "$(hostpath "$S/results"):/results" -v "$(hostpath "$S/reports"):/reports" \
    "$BENCH_IMAGE" "$@" || warn "report failed: $*"
}
rep -cmd report -results "/results/$DATE-writes-isolated" -report-out "/reports/$DATE-writes-isolated.md"
for id in history-unbounded history-cap20 top-1000charities size-small size-tiny; do
  rep -cmd report -results "/results/$DATE-$id" -report-out "/reports/$DATE-$id.md"
done
rep -cmd report-compare -results "/results/$DATE-history-unbounded" -compare "/results/$DATE-history-cap20" \
  -label-a unbounded -label-b cap20 -report-out "/reports/$DATE-regime-history.md"
rep -cmd report-compare -results "/results/$DATE-history-unbounded" -compare "/results/$DATE-top-1000charities" \
  -label-a "10 charities" -label-b "1000 charities" -report-out "/reports/$DATE-regime-charities.md"
rep -cmd report-compare -results "/results/$DATE-size-small" -compare "/results/$DATE-size-tiny" \
  -label-a small -label-b tiny -report-out "/reports/$DATE-regime-size.md"
for ch in 10 1000; do
  rep -cmd report-concurrency -results "/results/$DATE-contention-${ch}charities" \
    -report-out "/reports/$DATE-contention-${ch}charities.md"
done
for prof in unbounded cap20; do
  rep -cmd report-mixed -results "/results/$DATE-portal-$prof" -report-out "/reports/$DATE-portal-$prof.md"
done
log "experiment E complete"
