# Result — pin the PlantUML renderer and unify the diagram render helper

**Task:** `task-20260923T235910Z-diagram-renderer-pin`
**Branch:** `repo/pin-plantuml-render`
**Required tag:** `repo/plantuml-pinned`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`)
**Benchmark:** none — infrastructure only

## What was delivered

1. **A digest pin.** `infra/versions.env` gains
   `PLANTUML_IMAGE="docker.io/plantuml/plantuml@sha256:9b9ee6af54a86ea6ac53805e9da4da7fd913711cfe6c19b6ac1a41c398223aeb"`
   with the matching tag (`1.2026.8`, PlantUML 1.2026.8 / 149874a) named in the comment. The three
   study scripts no longer mention `:latest`.
2. **One shared helper.** `infra/diagram-render.sh` owns the image, the format, the bind-mount path
   conversion (`engine_path`: cygpath → `hostpath`, else `winpath`) and the provenance record, and
   supports `--format`, `--to DIR` and `--record FILE`. No renderer logic remains in the studies.
3. **Delegating study scripts.** `studies/0{1,2,3}-*/diagrams/render.sh` are 12-line delegations to
   the helper; their CLI (`./render.sh [svg|png]`) is unchanged.
4. **Per-study renderer provenance.** `studies/0{1,2,3}-*/diagrams/RENDERER.md` records the image,
   resolved digest, resolved image id, renderer version, format, source count and UTC time of the
   render that produced the committed SVGs.
5. `CONTEXT.md` records the new infrastructure.

## Re-render comparison (acceptance check)

Every study re-rendered all of its sources with the pinned digest; the fresh output was compared
byte-for-byte against the committed SVGs.

```text
study              sources  committed SVGs  byte-identical  differing
01-charity-tree       20         20              20            0
02-ticket-booking      5          5               5            0
03-reserved-seating    5          5               5            0
                      --         --              30            0
```

No committed SVG changed, so no digest migration and no unreviewed visual diff. The comparison is
recorded here; the delegated `render.sh` runs were also executed in place and produced no SVG diff.

## Validation

```text
bash -n infra/diagram-render.sh          ok
bash -n studies/*/diagrams/render.sh     ok
render 01: 20 diagrams -> rendered/      ok, provenance written
render 02:  5 diagrams -> rendered/      ok, provenance written
render 03:  5 diagrams -> rendered/      ok, provenance written
cmp committed vs pinned render           30 identical / 0 differing
host Java / PlantUML required            none
```

## Conclusion

The diagram renderer is now pinned in one place and every study gets it through one helper; the
committed SVG corpus is byte-identical to a render under the pin. No measurement was performed and
no database or benchmark lock was touched.
