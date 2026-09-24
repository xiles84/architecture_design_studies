# Progress — pin the PlantUML renderer and unify the diagram render helper

**Task:** `task-20260923T235910Z-diagram-renderer-pin`
**Branch:** `repo/pin-plantuml-render`
**Attempt:** `attempt-20260924T003519Z-d9d98e` (claim `claim-c5e3a0d780d00ffb`, epoch 1)
**Capability/role:** HIGH session executing a LOW task / executor (model `deepseek-flash`)
**Required tag:** `repo/plantuml-pinned`

Infrastructure-only. No database started, no benchmark lock taken, no measurement.

## Decision Log

1. **Pin by manifest digest, not a floating or exact tag.** `PLANTUML_IMAGE` in
   `infra/versions.env` is `docker.io/plantuml/plantuml@sha256:9b9ee6af…` (the digest the corpus was
   rendered with; tag `1.2026.8`, PlantUML 1.2026.8 / 149874a). Digest pinning matches the Typst
   toolchain and is stronger than an exact tag. Confidence: High.
2. **Pin to the version the committed SVGs already came from, then prove it.** Re-rendering all 30
   committed sources with the pinned digest is byte-identical for **30/30** SVGs, so there is no
   digest migration and no committed SVG is changed. Confidence: High.
3. **One shared helper at `infra/diagram-render.sh`; the three study scripts delegate.** The helper
   owns the image, the path conversion and the provenance record, so no renderer logic is
   duplicated. Confidence: High.
4. **The helper uses an `engine_path()` policy (cygpath → `hostpath`, else `winpath`).** The old
   scripts were inconsistent: 01 used `hostpath` (a no-op under WSL), 02 passed the raw path, and
   03's `cygpath` branch is dead under WSL; from this host none of them could bind-mount the sources
   into Windows `podman.exe`. Copying `book/build.sh`'s policy fixes that without touching
   `infra/lib.sh`. Confidence: High.
5. **Provenance record at `diagrams/RENDERER.md`, beside `rendered/`.** The acceptance wants the
   resolved digest recorded next to the rendered output, but `studies/*/diagrams/rendered` is a
   forbidden path for this task, so the record sits in the diagrams directory next to it. Written by
   the helper, never by hand. Confidence: High.
6. **The helper can render to `--to DIR` and record `--record FILE`.** This is what let the task
   compare a fresh pinned render against the committed SVGs without writing into the forbidden
   `rendered/`, and it is what the later fidelity-audit task will use to stage corrections.
   Confidence: High.

## What was produced

- `infra/versions.env` — `PLANTUML_IMAGE` digest pin with the tag named in a comment
- `infra/diagram-render.sh` — shared renderer helper (pinned image, path conversion, provenance)
- `studies/{01,02,03}-*/diagrams/render.sh` — each a 12-line delegation to the helper
- `studies/{01,02,03}-*/diagrams/RENDERER.md` — the digest/version record for each rendered set
- `CONTEXT.md` — the diagram-renderer infrastructure bullet

## Re-render comparison (all three studies, pinned digest)

```text
study              sources  committed SVGs  byte-identical  differing
01-charity-tree       20         20              20            0
02-ticket-booking      5          5               5            0
03-reserved-seating    5          5               5            0
                      --         --              30            0
```

## Residual risks / limitations

- PlantUML output is deterministic at one pinned version in this corpus, but the check is
  byte-identity, not a semantic diff; a renderer upgrade should re-run it and review any diff.
- `RENDERER.md` is written by the helper and will change its "rendered at" line on every run; later
  tasks should treat the digest/version fields as the provenance and regenerate the file with the
  render, not hand-edit it.
- The pinned digest is the amd64 manifest the local engine resolved; a different architecture would
  need its own digest recorded (not needed on this host).
