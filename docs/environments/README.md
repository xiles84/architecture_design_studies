# Environments

One page per machine that results were produced on.

**Results are comparable only within an environment** ([methodology rule 2](../methodology.md)).
Hardware moves benchmark numbers by more than most design decisions do, so a result from
one machine must never be compared against a result from another. Every result file names
its environment id, and that id must have a page here.

## Adding yours

Copy an existing page and fill in:

1. **Identity** — a stable id (`host-<something-distinctive>`), used verbatim in results.
2. **Hardware** — CPU model, physical vs logical cores, **core topology** (heterogeneous
   performance/efficiency cores change everything about scheduling), RAM type and speed,
   storage type, OS build.
3. **Container runtime** — Podman version, VM provider, and what the VM *actually* got, not
   what you asked for. `podman info` is the authority; a 32 GiB laptop can present 15 GiB.
4. **Images under test** — and where they are pinned.
5. **Known measurement hazards** — the most important section, and the one people skip.

Then tag your runs:

```bash
BENCH_ENVIRONMENT=host-my-machine ./run-study.sh --scale small
```

## On the hazards section

This is not a disclaimer. It is the part that tells a reader which conclusions your numbers
can carry and which they cannot.

Write down the things that would make someone misread your results — client and database
sharing cores, no real network between "nodes", a laptop that thermally throttles under
sustained load, a filesystem translation layer, cores the guest kernel cannot tell apart.
For each one, say what was done about it and what it means for interpretation.

A hazard you know about and state is a limitation. The same hazard unstated is an error.
