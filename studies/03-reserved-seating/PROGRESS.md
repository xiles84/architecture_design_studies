# Study 03 — execution progress

Executor log for [HANDOFF.md](HANDOFF.md). Facts only: what was done, which commit, which
mapped decisions were taken and why. Interpretation belongs to the analysis phase.

**Executor:** Claude Opus 5 (`claude-opus-5`), Claude Code desktop. Sessions from 2026-09-14.

## Steps

| Step | What | State | Commit / tag |
|---|---|---|---|
| 1 | Scaffold, README, PROGRESS | done | (this commit) |
| 2 | Study 02 terminology (Part A) | pending | |
| 3 | Engine probe | pending | |
| 4 | Platform: ErrLockNotAvailable | pending | |
| 5 | SQL catalogue, 13 designs | pending | |
| 6 | Harness | pending | |
| 7 | Diagrams | pending | |
| 8 | Dev checks and calibration | pending | |
| 9 | Main matrix (`small`) | pending | |
| 10 | Repeated race trials | pending | |
| 11 | Context, lessons, README; ready for analysis | pending | |

## Mapped decisions (handoff §11)

| # | Decision | Reason |
|---|---|---|

## Observations for the analysis phase

Facts recorded during dev checks and runs (handoff §13).

## Environment notes

- 2026-09-14: the Podman machine was stopped (the host had restarted); started with
  `podman machine start`. 8 CPUs and 15.4 GiB are visible, as the environment page records.
  Images from other projects (`dnakit-*`) exist on this Podman machine. No containers of
  theirs were running when this study's work started.
