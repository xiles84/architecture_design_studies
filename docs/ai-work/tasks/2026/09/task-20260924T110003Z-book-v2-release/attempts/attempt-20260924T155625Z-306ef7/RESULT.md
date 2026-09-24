# Result — Edition 2 release

**Task:** `task-20260924T110003Z-book-v2-release`
**Branch:** `repo/book-v2-release`
**Required tag:** `repo/data-architecture-book-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — the only container started is the pinned Typst image

## What was delivered

The tracked artefact in `book/dist/` was the **49-page Edition 1 build of a v3 registry**, committed
before any of this work started. It is now **Edition 2**, built once from fully merged sources and
committed with its manifest.

**Contents relative to Edition 1:** the v4 and v5 registries (stale placement/endpoint gaps retired, five
claims re-scoped to their own limits), the cache chapter rebuilt as mechanism cards with a corrected
relaxed-freshness taxonomy, the terminology and page-3 evidence-legend pass, figure numbering owned by
Typst with the renderer's text stretching stripped at export, a part opener, a 26-term glossary, three new
figures (decision path, freshness timeline, lease-versus-fence) plus the benchmark-reading checklist as a
table, and repeat citations as capsules carrying every limit with the full card in the registry index.

**Two release blockers found while producing it**, both in the provenance records rather than the page:

1. `book/build.sh` hard-coded `"edition": "Edition 1 (draft)"` while `lib/config.typ` said Edition 2, so the
   manifest described a different book from the PDF. The edition is now read from its single source, and
   the build **fails** if that string does not appear in the extracted page text. Every existing gate
   passed on the wrong manifest.
2. `book/assets/export.sh` made figure staleness depend on `source_revision` — a property of the
   repository, not the artefact. Written before its own commit it recorded `uncommitted`; checked after it
   recorded a hash, and the build failed on a figure nobody had touched. The field stays as documentation
   and is excluded from the comparison, which now uses the source and asset content hashes, the renderer
   identity and the postprocess step.

## Verified from committed files

```text
file sha256 == manifest pdf_sha256     58bea049… → 684318f3… (final build), equal
pages                                  54
claims indexed                         29/29
fonts embedded                         6/6
evidence digest in pdf                 true
unverified banner in pdf               false
unresolved tokens in pdf               false
edition on page                        true — "Edition 2 (draft)", zero occurrences of "Edition 1"
manifest                               edition "Edition 2 (draft)", release true, registry v5, 8 figures
check.sh                               OK, no waivers
book/assets/export.sh --check          passes for all eight figures
```

## Where this task is weak

- **The PDF is not byte-reproducible** across runs: the build clock is an input, so `pdf_sha256` describes
  one build. The source-tree hash and the evidence digest are stable, and they are the fields that identify
  *inputs*; the artefact hash identifies the file.
- **`source_dirty` is true in this manifest**, because the release commit itself is the commit that carries
  the PDF and its build.sh/export.sh fixes. A future release built from a clean tree will report false.
- **The reviewer has not seen it yet.** Their round-2 and round-3 findings are what this addresses; whether
  they consider it releasable is the verdict this artefact exists to obtain.
- **Tranche 2 of the visuals is still unimplemented** (V2, V3, V4–V6, V9, V10, V12), and the `Edition 2`
  label is still marked `(draft)`, which is deliberate until an external reader signs off.
- **The Edition 1 PDF is replaced, not archived** in `dist/`. Its state remains recoverable from the tags
  `repo/data-architecture-book-v1` and `repo/data-architecture-book-v1-draft`.

## Conclusion

The book's tracked artefact and its provenance record now describe the same document, and the release is
ready to be judged by someone other than its author.
