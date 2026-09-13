# Replicating these studies

Everything runs in containers. You need **Podman** and nothing else — no database, no Go
toolchain, no `psql`, no Java.

## Prerequisites

| | |
|---|---|
| Podman | 5.x (tested on 5.8.1) |
| Disk | ~8 GB for images, plus room for data volumes |
| Memory | ≥ 8 GiB available to the container VM; **16 GiB for the 3-node cluster** |
| CPU | 8 cores recommended (the 3-node topology asks for 2 per node plus 2 for the client) |

No compose provider is needed. Orchestration is plain `podman` CLI in POSIX shell
scripts, which work identically under Linux, macOS, WSL and Git Bash on Windows.

### macOS and Windows

Podman runs containers inside a Linux VM. Start it and give it enough resources:

```bash
podman machine init --cpus 8 --memory 16384 --disk-size 100
```

```bash
podman machine start
```

If the machine already exists but is too small, `podman machine set --cpus 8 --memory 16384`
then restart it. Check what the VM actually got:

```bash
podman info --format '{{.Host.CPUs}} cpus / {{.Host.MemTotal}} bytes'
```

> **Windows note.** The scripts set `MSYS_NO_PATHCONV=1` and convert paths with `cygpath`
> before handing them to Podman. Without that, Git Bash rewrites `/results` into a
> Windows path and Podman silently mounts an empty directory instead of yours — the run
> appears to succeed and produces no files. If you invoke the benchmark container by hand
> from Git Bash, pass Windows-style paths (`C:/path`) for bind mounts.

## Run a study

```bash
cd studies/01-charity-tree
```

```bash
./run-study.sh --scale small --duration 10s
```

That brings each topology up in turn, runs every design against it, tears it down, and
writes results to `results/<run-id>/`. Expect roughly:

| Scale | Approx. wall time for the full matrix |
|---|---|
| `tiny` | 15 min — for checking the pipeline works |
| `small` | ~2 h |
| `medium` | ~5 h |

### Useful slices

```bash
# one engine, two designs — the foreign-key comparison on its own
./run-study.sh --topologies pg-single --designs d3_flattened_fk,d8_flattened_nofk
```

```bash
# skip the writes, which are the slowest part
./run-study.sh --topologies pg-single --write-ops ""
```

```bash
# leave the database running afterwards so you can poke at it
./run-study.sh --topologies pg-single --designs d4_rollup_trigger --keep-up
```

### Study 02

```bash
cd studies/02-ticket-booking
```

```bash
./run-study.sh --scale small --tag
```

The study 02 image is built from the **repository root** (it includes the shared `platform/`
module) and is named `localhost/ticketbench:1`, so it never replaces another study's client.
YugabyteDB nodes are started with `yb_enable_read_committed_isolation=true` (from
`study.env`). The runner generates `reports/<run-id>.md` at the end of the matrix.

To reproduce a tagged run exactly:

```bash
git checkout run/02-ticket-booking/<run-id>
```

Expect the `small` matrix to take several hours: YugabyteDB sells far fewer seats per second
than PostgreSQL on this hardware, so its large-event races run to their timeout.

## Poke at a database by hand

```bash
../../infra/pg-single.sh up
```

```bash
../../infra/pg-single.sh psql
```

The YugabyteDB equivalents are `infra/yb-single.sh` and `infra/yb-cluster3.sh`, whose
SQL shell subcommand is `ysqlsh`. Each script also accepts `down`, `logs` and `dsn`.

> YugabyteDB binds YSQL to its `--advertise_address`, **not** to `127.0.0.1`. Connecting
> to loopback inside the container is refused; use the container name as the host.

## Generate the report

```bash
podman run --rm -v "$PWD/results:/results" localhost/charitybench:1 \
  -cmd report -results /results/<run-id> -report-out /results/<run-id>/report.md
```

## Regenerate the schema diagrams

```bash
./diagrams/render.sh svg
```

## What a run produces

```
results/<run-id>/
├── manifest.yaml                  what was asked for, and which cells failed
├── pg-single/
│   ├── topology.yaml              what was actually running (image, cpus, memory)
│   ├── d1_normalized_minimal.json result: timings, plans, sizes, verification
│   └── plans/
│       └── d1_normalized_minimal.txt   readable EXPLAIN output per query
├── yb-single/
└── yb-cluster3/
```

The JSON is the source of truth. Each file carries everything needed to interpret it
without external context: the engine version banner, the dataset fingerprint, the
benchmark settings, the correctness results and the plans.

## Clean up

```bash
../../infra/pg-single.sh down
../../infra/yb-single.sh down
../../infra/yb-cluster3.sh down
```

Each `down` removes the container **and its data volume** — a benchmark must never start
from a previous run's pages, statistics or bloat. To remove the images too:

```bash
podman rmi localhost/charitybench:1 docker.io/library/postgres:17.11 docker.io/yugabytedb/yugabyte:2025.2.6.0-b111
```

## If your numbers differ from the published ones

They will, and that is expected — see [methodology](methodology.md) rule 2. Results are
only comparable within an environment. If you want to publish a comparison from your own
machine, add a page under [`environments/`](environments/) describing it, including the
measurement hazards specific to that hardware, and tag your run with that environment id:

```bash
BENCH_ENVIRONMENT=host-my-machine ./run-study.sh --scale small
```
