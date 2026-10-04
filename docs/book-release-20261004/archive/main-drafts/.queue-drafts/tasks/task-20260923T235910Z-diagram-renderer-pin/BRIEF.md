# Execution brief — Pin the PlantUML renderer and unify the diagram render helper

Read the visual-explanations brainstorm conclusion
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies/CONCLUSION.md`),
`infra/versions.env`, `infra/lib.sh`, and the three `studies/0{1,2,3}-*/diagrams/render.sh`.

## Goal

All three diagram render scripts use `IMAGE="docker.io/plantuml/plantuml:latest"`, and `PLANTUML_IMAGE` is
absent from `infra/versions.env`. That breaks the repository's own pinning rule for the visual layer while
every measurement image is pinned. Pin the renderer and make the three studies share one helper.

## Decisions already made (do not re-litigate)

- Keep **PlantUML**; do not switch renderer.
- The image goes in `infra/versions.env` (single pinning location), referenced by a shared `infra/` helper.
- A byte-different SVG after pinning is not automatically a defect: PlantUML output can move with the renderer
  version. Record the migration and review the visual diff; only unrecorded change is a failure.
- Do not edit `.puml` sources or committed SVGs in this task; that is the fidelity-audit task.

## Acceptance criteria

- `PLANTUML_IMAGE` pinned to an exact tag or digest in `infra/versions.env`.
- One shared helper in `infra/`; the three study scripts delegate to it and carry no duplicated renderer logic.
- Every study re-renders; output is byte-identical to the committed SVGs, or each differing file is reviewed and recorded.
- The resolved image digest is recorded with the rendered output.
- No host Java/PlantUML needed; everything renders in Podman.

## Constraints

- No measurement; this task takes no benchmark lock, but it must not start or stop databases.
- Own only the render scripts and infra pinning; do not touch diagram content, book sources or study results.
- Commit explicit paths, tag `repo/plantuml-pinned`, merge through the queue lifecycle.

## Escalate only if

- the PlantUML image cannot be resolved or pulled on this host (then the pin is impossible and the honest gap must be documented).

NEXT MODEL: LOW
