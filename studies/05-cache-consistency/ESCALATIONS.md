# Escalations Required — Study 05

Unmapped decisions are recorded here, with evidence, blocked steps, options and the decision
needed. HIGH decides and appends an `AM-NN` amendment to `HANDOFF.md`, then the dependent work
resumes. This agent is both HIGH and the executor, so each entry records the decision explicitly
and is tagged as an immutable handoff-amendment milestone.

---

## 05-ER-01 — Engine access from WSL

- Raised by: DeepSeek HIGH selected by user (`deepseek-flash`; effort and tool identity not
  exposed by the session — HIGH role)
- Raised at: 2026-09-21 preflight
- Handoff ID / revision / producing commit: `EH-05` revision 1 @ `study-05/v0-handoff`
- Status: **decided — no handoff change needed**
- Blocked steps: none
- Next setting: **HIGH — DeepSeek HIGH continues the mapped work**

### Observation and evidence

There is no native WSL `podman` binary: `podman` is not on `PATH`, and
`source infra/lib.sh` resolves `ADS_PODMAN=/mnt/c/Program Files/RedHat/Podman/podman.exe`
(`podman version 5.8.1`), which already reports the same images, volumes and benchmark-lock
state as every other session.

### Unmapped decision

Whether Study 05 needs its own engine-bridge decision.

### Decision

**No new decision is required.** Study 04's `04-ER-01` already resolved this
identically and its resolver is merged in `infra/lib.sh` at the starting revision. EH-05 §2.9
therefore adopts it unchanged: WSL drives the documented Windows engine through `podman.exe`;
a WSL-local podman is never initialised, because a second engine would carry no shared benchmark
lock. This entry exists only so the observation is on the record rather than looking unexamined.

---

## 05-ER-02 — Redis image pin (methodology 3)

- Raised by: DeepSeek HIGH selected by user (`deepseek-flash`; HIGH role)
- Raised at: 2026-09-21 preflight
- Handoff ID / revision / producing commit: `EH-05` revision 1 @ `study-05/v0-handoff`
- Status: **decided — no handoff change needed**
- Blocked steps: none
- Next setting: **HIGH — DeepSeek HIGH continues the mapped work**

### Observation and evidence

`infra/versions.env` pinned PostgreSQL, YugabyteDB and the Go builder, but no cache engine; the
repository rule (methodology 3) is that every image is pinned in `infra/versions.env`.

### Unmapped decision

Which Redis tag, and where the pin lives.

### Decision

`REDIS_IMAGE="docker.io/library/redis:7.4.11-alpine"`, appended to `infra/versions.env` — the
newest patch of the 7.4 line at the time of planning, pinned to an exact tag, never `latest`.
Its image id is recorded in every run's manifest and in the Redis capability probe. Append-only:
no existing value in that file changed, so studies 01–04 are unaffected.
