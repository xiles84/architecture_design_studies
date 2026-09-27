# Diagnosis: why `l3_section_sharded` never finishes its load (study 03)

**Status: answered, and the answer is bigger than the cell.** The diagnostic run itself never started — it
was refused before measuring — and the reason for the refusal is the reason `l3` kept dying.

## What was found

**The podman machine has 2 GiB of memory.** Every container in this project runs inside it, while
`infra/versions.env` allots the database **3 GB** and the client **2 GB**:

```
podman machine inspect podman-machine-default --format '{{.Resources.Memory}} {{.Resources.CPUs}} {{.State}}'
2048 4 stopped
```

So **5 GB of declared container budgets on a 2 GiB machine**. No container limit can be enforced; the VM's
own ceiling is the real one, and each container's declared budget has always been fiction. This is the
mechanism behind the three lost cells: the database's own processes were **killed with signal 9** (122 such
lines in one cell's log), which is memory pressure at the VM level, not a design behaving badly.

It also explains why raising the container limits changed nothing that mattered:

| Pass | Container budget | tserver flag | Outcome for `l3` |
|---|---|---|---|
| 1 and 2 | 3 GB (never available) | 1.5 GiB | killed mid-`COPY`, SQLSTATE 57P01 |
| 3 | 5 GB (never available) | 3 GiB | `connect: connection refused` |

**The machine was also stopped when this diagnosis began**, and that produced a third symptom that had
nothing to do with `l3`: the runner refused with *"another benchmark holds this machine's lock: (lock volume
vanished)"*, and the engine's own client errored with exit code 125. A stopped engine makes the lock volume
unreadable, the lock check cannot pass, and every run is refused — correctly.

## What the diagnostic did not do

It did not run `l3` alone, because the engine was down. **D1, D2 and D3 are therefore unevaluated**, and no
statement is made about whether `l3`'s load would succeed in isolation. That question remains open and is
cheap to answer once the sizing decision below is made — it should not be answered before, because a 2 GiB
machine cannot support a conclusion about a load that needs more than 1.5 GiB of database memory.

## Consequences for the corpus, stated plainly

1. **Every published run on this host ran under an unenforced budget.** The manifests record
   `budget_per_db_node: cpus=2 memory=3g` because that is what the runner was told; the machine could not
   provide it. Results are comparable with each other, since the same ceiling applied throughout, but the
   declared envelope does not describe the environment, and any comparison against a reader's own database
   must say so.
2. **Study 07's topology work cannot be attempted here at all**: a 3-node cluster with a 3 GB per-node
   budget does not fit a 2 GiB machine, and this is the sizing calculation the methodology requires before
   such a study runs.
3. **`v2-11`'s registration is unaffected in its direction** — the 27.8x/24.4x agreement came from two runs
   under the same unenforced ceiling, and the ratio is what the claim states — but the sizings in its
   protocol are labelled rather than validated, which is a coverage gap the owner's rules require to be named.

## The decision this needs (escalated, not taken)

Two options, both environment changes, neither of which an executing agent should choose silently:

- **Resize the machine**: `podman machine set --memory 8192` and restart, so the declared budgets become
  enforceable. This changes the environment for every study; runs before and after are not comparable and
  the change must be recorded as an environment change, with the studies' manifests noting the old ceiling.
- **Shrink the declared budgets** to fit 2 GiB (for example database 1 GB, client 512 MB), which makes the
  documentation honest at the cost of making the database containers much smaller than the published studies
  assume.

## A smaller lesson, from the same hour

My own lock check — `podman volume inspect ads-run-lock >/dev/null 2>&1 && echo HELD || echo FREE` — printed
**FREE** while the engine was stopped, because a *failed* inspection and an *absent* volume both fall into
the `||` branch. The runner's own `run_lock_acquire` got this right (it refused, naming the vanished volume);
my ad-hoc version reported the machine as idle and ready. **A failure and an absence are not the same
answer**, which is the third time this pattern has appeared in this corpus (a null from my own query read as
an absent field; a grep for a fingerprint matching the wrong field). The check should distinguish three
states — held, free, unreadable — and refuse to proceed in the third.
