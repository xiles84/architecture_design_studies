# Cache publication and deletion experiment

This is the bounded v5 extension of Study 05. It compares authoritative version
checks, synchronous invalidation fences and acknowledged outbox delivery, then
tests deletion/recreation and absence-marker expiry. The independent Go ledger
checks committed payloads, absence and the strict acknowledgement contract.

Read `PROTOCOL.md` and `HIGH-REVIEW.md` before running. It uses 128 keyed records,
not the earlier donor portal. SQL and the oracle are compiled into a distinct
client image; every candidate shares the same source transactions and outbox
append. A failed control or candidate read removes that cell's timing.

From a clean, committed task worktree:

```bash
bash studies/05-cache-consistency/protocolbench/run.sh <new-run-prefix> 3s 3
```

The runner acquires `ads-run-lock`, builds and tests the Go client in Podman,
uses fresh PostgreSQL/Redis containers for each whole run, and creates annotated
producing run tags. It refuses an existing prefix, service name or data volume.
It records three trials per phase, seeded order and a same-candidate first/last
position control. Results are exclusive writes under `results/<prefix>-rN/`;
generated factual reports go under `reports/`. Only the runner's own services
are removed at exit. Do not run a second database measurement on this machine.

Afterward, the same client image can generate the fixed whole-run reduction:

```bash
podman run --rm -v <host-study-results>:/results:ro -v <host-reports>:/reports \
  localhost/cache-protocol-v5:20261003 --cmd summary --inputs /results \
  --run-id <prefix> --out /reports/<prefix>-whole-runs.md
```

Use the image id recorded in the manifests for reproduction. Summary reduces
three trials within each whole run before calculating the across-run median
and min/max; position-last cells remain separate. Interpretation is a signed
analysis indexed by each exact result digest. Reports are never overwritten:
move a superseded generated report to `reports/outdated/` before regeneration.
Optional fault timings are recorded in JSON only when measured; generated
`0.000` placeholders elsewhere mean no such timing was recorded.

Single source endpoint, fixed payload, per-key client writer ordering and
short closed-loop phases limit the conclusions. Cluster, multi-host placement,
uncoordinated hot-key contention and saturation are unmeasured. Saved Podman
machine configuration differs from guest-observed WSL resources on this host;
use the retained guest and container records, not configuration alone.

NEXT MODEL: HIGH
