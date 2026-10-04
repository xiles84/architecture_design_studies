# Release Decision Log — GPT-6 / Codex desktop / HIGH

1. **Edition 3 preview 1.** The approved pilot is only a slice of the full rewrite;
   label partial conversion honestly while releasing the combined source. High
   confidence. Alternative: falsely label the full Edition 3 conversion complete.
2. **Preserve then clean.** Archive loose/ignored drafts and pilot outputs, retain
   old run evidence at its original paths, and prove raw SHA-256 values in Git.
   High confidence; affects every completed-worktree cleanup.
3. **Correct the exporter and qualify inference.** Fix the RunTag reader; add an
   independent exact manifest and signed interpretation companion. Preserve the
   original reports/analysis. High confidence; affects SELECT review and synthesis.
4. **Portable source checkout.** Convert Git Bash paths for Git, pin LF for book
   orchestration and canonical PlantUML sources, and use existing WSL orchestration
   with the shared Windows Podman engine. No host toolchain installation. High
   confidence; figure/source identities were verified, not silently regenerated.
5. **Short explanations with useful visuals.** Split the Rollup opener and two
   dense legacy passages, shorten the figure equivalent, correct Typst emphasis
   and keep the five-row cost table together. High confidence; numeric values and
   evidence package unchanged. Full legacy conversion remains unfinished.
6. **Keep yielded work.** Merge join/CTE's committed plan/queue state but retain
   its ready worktree and branch. High confidence; no measurement or completion
   is asserted for that family.

Release source and rendered proof accepted. Cleanup follows main integration;
final operations are captured in `docs/book-release-20261004/CLEANUP-RESULT.json`.

NEXT MODEL: HIGH
