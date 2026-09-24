#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card, registry-label

#heading("Evidence and reproduction")

#marker("chapter", "evidence and reproduction")
Every number in this book resolves to a claim in #raw(registry-label), and every claim to
a measured cell in a cited report. This chapter explains the confidence cards, the run-tag
provenance model, and how to reproduce a run in Podman.
#heading(level: 2, "The provenance model")
#direct[
  A run tag marks the producing code state, not a tree that contains the run's results. Check out
  the tag for the code, then read the committed results directory for the data.
]
#registry-card("v2-gap-06-native-datastore-families")
#registry-card("v2-gap-01-ledger-load-multiplier-unstable")

#heading(level: 2, "Reproducing the book itself")
The book is built in Podman from a pinned toolchain; the manifest records every input and the
resulting hash.

```bash
book/build.sh          # builds the pinned image, compiles, verifies, writes the manifest
cat book/dist/build-manifest.json
```

#heading(level: 2, "Reproducing a measured run")
```bash
git checkout run/<study>/<run-id>     # the producing code state
infra/<topology>.sh up                # take the benchmark lock first
cd studies/<study> && ./run-<engine>.sh
```
A run's results live in `studies/<study>/results/<run-id>/`; the report beside them is generated,
never hand-edited, and its inputs digest identifies the data.
