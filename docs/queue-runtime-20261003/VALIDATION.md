# Queue runtime repair validation

Host: xiles Windows/Podman environment; Go image `docker.io/library/golang:1.26-bookworm`. Source base `be397e9be97f227267f37bd18abb831f848fab7a`; isolated branch `repo/queue-runtime-20261003`. No database or benchmark is started.

- `go vet ./... && go test ./...` in Podman: all five tested packages passed, including multi-process claims, integration locking, recovery, dirty-file preservation, and stale-epoch rejection.
- Added `TestRecoveredClaimCanSubmitWithoutExplicitFlags`: two-step recovery in a temporary task worktree persists the new epoch/event, guard succeeds, committed RESULT submits to awaiting_review. Passed.
- Added `TestImageSourceRejectsDriftAndNormalizesCRLF`: matching digest accepted; LF/CRLF accepted as identical; changed Go source and unmarked image refused; incomplete source refused. Passed.
- `bash -n tools/queue/queue`: passed. `git diff --check`: passed.
- Repaired wrapper invoked against installed image `9481de7c07cb748ba90b059d3e139f1ac358ccd5d908daf3b840b120530878d4`: refused its unknown `image-check` before dispatching the requested list, with an explicit rebuild instruction.

Rebuild through `tools/queue/queue --build` passed its full vet/test gate and
produced image `06e706e59ff488249fe1c463a4857915f4e0893b98a714ce1d02f7f492f4c3da`.
Embedded/current source digest:
`1c6c46eb05da53419a89eed0f04a964344b29e49bb8ebb43fa75a5198d1c5712`.
The production wrapper's `guard` accepted the matching image and current repair
claim epoch 1. The image also accepted the bind-mounted Windows checkout directly,
establishing CRLF portability in addition to the unit regression. This check
report establishes queue behavior, not scientific measurement validity.
