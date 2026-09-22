# Review — book synthesis v1 (draft)

**Task:** `task-20260922T025912Z-book-synthesis-v1`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Reviewed:** `attempt-20260922T191804Z-49e273`, submitted at `579d98a`.
**Decision:** **approve for integration** as the v1 draft. Two recommendations are recorded for the
release review rather than treated as blocking.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | Write only from the 23 active v2 claims; cite through evidence cards | Accept; verified mechanically below. |
| 2 | Keep the six families and the progressive structure | Accept; each family section has all required sub-sections. |
| 3 | Mechanisms live in concepts, families point to them | Accept; no duplicated mechanism exposition. |
| 4 | Confounds, gaps and analogies on the page | Accept; all six confounds are listed and the analogy transfers are labelled. |
| 5 | No JSONB-as-native-document claim, no Redis-as-authoritative-store claim | Accept; both are explicitly disclaimed. |
| 6 | No laptop-topology capacity guidance | Accept; the caveat recurs in topologies, every family boundary and the provenance page. |
| 7 | Numbers kept with their caveats | Accept; the 128-buyer scope, the unsettled PostgreSQL ratio and the upper-bound status are stated where the figures appear. |
| 8 | The draft is generated and verified | Accept; the build result is recorded below. |

## Independent checks

```text
registry-card ids used in prose: 21; ids that do not resolve: 0
v2 claims not cited in prose:    v2-gap-03 (Study 04 missing designs),
                                 v2-gap-04 (Study 05 unmeasured regimes)
committed PDF blob sha256 == manifest pdf_sha256:
   8490e4d89e25876a974956583648852a38dc26cf9996ecaa427fad54ffaad59a  (match)
manifest source_dirty=false; pages 37; fonts 4/4; claims indexed 23/23; digest in pdf true
```

The build itself guarantees the first check (an unknown claim id aborts compilation); the reviewer
re-ran the extraction to confirm no prose figure is cited without an existing card.

## Recommendations for `book-release-review-v1` (non-blocking)

1. The two family-coverage gaps (`v2-gap-03`, `v2-gap-04`) appear in the evidence index but are not
   cited in any chapter's prose. The release review should either cite them or state that they are
   deliberately index-only.
2. A claim-to-page map is not rendered; a future edition could confirm mechanically that every claim
   is reachable from at least one page, not only from the index.

## Result reviewed against the brief

- **Six families with the full progressive structure**: present.
- **Variants / technology / topology / workload as separate axes**: present; the variants chapter
  marks which pairs are controlled and which are named but unmeasured.
- **Reuse concepts instead of repeating mechanisms**: done.
- **Named domains with every transfer labelled**: done (ticketing, seating, configuration, donor
  portal direct; supermarket/auction/e-commerce labelled analogy).
- **No native-family or cache-authority over-claim**: confirmed.
- **Draft PDF plus Decision Log submitted**: done.

## Integration

Approve and integrate into local `main`; tag the integrated state
`repo/data-architecture-book-v1-draft`.
