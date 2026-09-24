# Review — pin the PlantUML renderer and unify the diagram render helper

**Task:** `task-20260923T235910Z-diagram-renderer-pin`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Verdict:** **approved for integration** — no rework required
**Required tag:** `repo/plantuml-pinned`

## Decision Log reviewed first

| # | LOW decision | Verdict |
|---|---|---|
| 1 | pin by manifest digest (tag `1.2026.8` named in the comment) | accepted — strongest reproducibility; no `:latest` remains in any render script |
| 2 | pin to the version the corpus came from and prove byte-identity | accepted — reviewer re-ran the comparison: 30 identical, 0 differing |
| 3 | one shared helper, three delegations | accepted — no duplicated renderer logic; study CLI unchanged |
| 4 | `engine_path` policy copied from `build.sh` rather than editing `infra/lib.sh` | accepted — the old scripts could not mount under WSL; `lib.sh` is out of owned scope |
| 5 | `RENDERER.md` beside `rendered/` rather than inside it | accepted — `studies/*/diagrams/rendered` is a forbidden path; a record next to it satisfies the intent |
| 6 | helper supports `--to`/`--record`/`--format` | accepted — used for the comparison and by the fidelity-audit task |

No decision materially degraded correctness, architecture or reproducibility.

## Independent checks re-run by the reviewer

```text
grep plantuml:latest studies/*/diagrams/render.sh infra/diagram-render.sh   0 matches
cmp committed SVGs vs pinned re-render                                     30 identical / 0 differing
bash -n on the helper and all three render scripts                         ok
PLANTUML_IMAGE present in infra/versions.env as a digest                   yes
diagrams/RENDERER.md present in all three studies                          yes
```

Acceptance criteria:

- `PLANTUML_IMAGE` is pinned to a digest in `infra/versions.env` — met, with the tag recorded.
- One shared helper; the three study scripts delegate and carry no renderer logic — met.
- Every study re-rendered; output byte-identical or reviewed — met (all identical).
- Resolved digest recorded beside the rendered output — met (`diagrams/RENDERER.md`).
- No host Java/PlantUML — met (Podman only).

## Findings (non-blocking)

1. `RENDERER.md`'s "rendered at" line changes per run; later tasks must regenerate it with a render,
   not hand-edit it (recorded in the executor's limits).
2. The pin is the amd64 manifest digest; other architectures would need their own digest.

## Verdict

Approved for integration. The pin is real, the helper is the single implementation, and the
committed SVG corpus is provably unchanged under the pin.
