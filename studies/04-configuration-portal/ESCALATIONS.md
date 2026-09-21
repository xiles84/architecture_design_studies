# Escalations — Study 04, configuration portal

Items the committed handoff did not map, recorded here before any dependent work continues.
Raised and resolved items stay visible; amendments are appended to `HANDOFF.md`, never written
over it.

---

## 04-ER-01 — No native WSL `podman`: the documented engine is reachable only as `podman.exe`

- Raised by: DeepSeek HIGH (`deepseek-flash`), effort setting not exposed, tool identity not exposed
- Raised at: 2026-09-21T11:38Z
- Handoff ID / revision / producing commit: `EH-04` revision 1, commit pending (`study-04/v0-handoff`)
- Status: **decided — see HIGH decision below (AM-01)**
- Blocked steps: every container operation — engine probes (phase A), the WSL bind-mount probe
  (phase B), dev checks (phase C), the matrix (phase D) and everything after it
- Next setting: **HIGH — decide the engine path**

### Observation and evidence

The task's preflight assumes a native WSL `podman` client that shares the Windows sessions'
service and storage. The live checks disagree with that assumption.

```
$ command -v podman
(no output; exit 1)
$ ls /usr/bin/podman* /usr/local/bin/podman*
ls: cannot access '/usr/bin/podman*': No such file or directory
$ command -v docker
(no output; exit 1)
$ command -v podman.exe
/mnt/c/Program Files/RedHat/Podman/podman.exe
```

The Windows client **does** reach the documented engine, and from WSL it already shows the same
storage the Windows sessions use:

```
$ podman.exe version
Client:  Podman Engine      Version: 5.8.1    OS/Arch: windows/amd64
Server:  Podman Engine      Version: 5.8.5    OS/Arch: linux/amd64

$ podman.exe system connection list
podman-machine-default   ssh://user@127.0.0.1:60709/run/user/1000/podman/podman.sock   (default, rw)

$ podman.exe images
localhost/ticketbench  1   4849d5c3cc1a   (study 02's client image)
localhost/seatbench    1   e0a79946b106   (study 03's client image)
docker.io/library/golang  1.26-bookworm
$ podman.exe ps --all
2531b0135335  postgres:17.11  pg-scratch  Exited (0)
$ podman.exe volume inspect ads-run-lock
[]  Error: no such volume "ads-run-lock"          # the benchmark lock is free
```

The machine is also WSL2 on this same host and sees the repository at the path WSL uses:

```
$ podman.exe machine ssh -- ls /mnt/c/extra/code/architecture_design_studies
AGENTS.md  CLAUDE.md  CONTEXT.md  LESSONS_LEARNED.md  LICENSE  README.md  docs  infra
$ podman.exe machine ssh -- 'nproc; free -b | head -2; uname -r'
8
Mem:  16496422912   ...
6.6.87.2-microsoft-standard-WSL2
```

That is exactly the guest in `docs/environments/host-zenbook-ux5406sa.md` (8 CPUs, ~15.36 GiB,
same kernel), so the environment id remains valid and the measurements stay comparable with
studies 01–03.

Two secondary observations, recorded because they look like discrepancies and are not:

- `podman.exe machine inspect` reports `Resources.CPUs = 4`, `Resources.Memory = 2048`. That is
  the **stale value recorded at `machine init`**, not the live VM. The live guest is 8 CPUs /
  16 496 422 912 bytes, matching the environment page. No VM was created, resized or started.
- `hostpath()` in `infra/lib.sh` converts paths using `cygpath`, which exists under Git Bash but
  **not** in WSL. From WSL it is therefore a no-op and returns `/mnt/c/...`.

### Unmapped decision

The handoff assumes a `podman` on `PATH`. It cannot choose between the options below, because the
choice changes how every script reaches the engine and how bind-mount sources are spelled.

| Option | Tradeoff |
|---|---|
| **A. Drive `podman.exe` from WSL** | Same engine, same storage, same `ads-run-lock` — the lock keeps meaning what it means. Needs an explicit resolver in `infra/lib.sh` and a decided rule for bind-mount source paths. Deviates from the preflight's assumption of a native client, so it must be an explicit, attributed decision. |
| B. Initialise a WSL-local podman | **Rejected on sight.** A second engine has different images, different volumes and **no shared benchmark lock**. Two studies could then measure at once on the same 8 cores, which is the exact failure the lock was built to prevent. |
| C. Stop and report the environment gap | Honest but leaves Study 04 unstarted, when the documented engine is reachable and healthy. |

### Independent work

None that touches containers. The study's design catalogue, dataset definition, invariants and
this document were written while the decision was pending; no measurement depends on them yet.

### HIGH decision

- Decided by / at: DeepSeek HIGH (`deepseek-flash`), effort not exposed, 2026-09-21
- Decision and rationale: **Option A.** WSL drives the documented engine through `podman.exe`.
  `infra/lib.sh` gains an explicit, additive engine resolver (`PODMAN`/`ADS_PODMAN` → `podman` on
  `PATH` → the Windows client path) plus a `podman()` shell function, so every existing call site
  and `command -v podman` keep working unchanged; study 01–03 runners are not edited. A WSL-local
  engine is never initialised (option B is excluded). `winpath()` is added for bind-mount sources
  that need the `C:/…` form, and **phase B decides** whether it is actually required — if
  `podman.exe` accepts `/mnt/c/…`, `hostpath()`'s existing behaviour is left untouched and only
  the resolver is added.
- Handoff amendment: **AM-01** (§21 of `HANDOFF.md`)
- Next setting: **HIGH — execute the handoff with the same DeepSeek HIGH model in WSL**, starting
  with the phase B bind-mount probe.
