# Execution brief — Typst toolchain and shell

Create the exact `book/` structure in the goal. Pin the stable Typst image/version selected
at implementation time by immutable image digest, and pin every redistributable font.
Compile only in Podman. Use modular `.typ` sources; operational documentation remains
Markdown and evidence/manifest state remains JSON.

Implement the visual system, title/edition/provenance pages, table of contents, consistent
family/variant/topology/scenario callouts, evidence labels, source links, page headers,
footers, tables, and accessible code blocks. Populate chapter shells and navigation but
do not invent scientific prose assigned to the HIGH synthesis task.

Generate the draft PDF and build manifest from a clean commit. Verify the rendered pages,
embedded fonts, links, evidence digest, and PDF hash. Use actual clock/version outputs;
never hand-set them.

NEXT MODEL: LOW
