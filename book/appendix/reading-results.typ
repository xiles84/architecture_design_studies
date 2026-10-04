#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-label
#set heading(numbering: none)
#heading(level: 2, numbering: none, "How to read and reproduce a result")

Read the comparison, its workload and its limit before using a summary number.
All measurements share one host; they do not model a real inter-node network.

#table(columns: (1fr, 1.7fr), inset: 6pt, stroke: 0.4pt + palette.rule,
  table.header([Check], [What to establish]),
  [Answer], [Same grain, duplicates, NULLs, ties and ordering; independent correctness gate.],
  [Comparison], [One controlled decision, with matching data and resource framing.],
  [Instrument], [Trials, whole-run repeats, measurement order and drift.],
  [Boundary], [Observed regime, expected failure condition, or not measured.],
  [Provenance], [Claim card, inputs digest, producing run tag and committed results.],
)

#figure-evidence("../assets/fig-decision-tree.svg", "structure", "conceptual illustration",
  "Start with the required answer and the measured bottleneck.",
  [Choose an access path, information placement or coordination mechanism only after identifying
   the contract and the expensive work. The branches are a selection aid, not benchmark rankings.])

#heading(level: 3, numbering: none, "Facts, derived state and history")
The two pictures answer different questions: where a reusable answer is maintained,
and whether the system must retain what happened.

#figure-evidence("../assets/fig-derived-vs-history.svg", "structure", "conceptual illustration",
  "A maintained answer and an event history have different roles.",
  [Derived state summarizes source facts. History retains changes. A projection can make history
   cheap to read, but its maintenance needs a correctness rule.])

#pagebreak(weak: true)
#heading(level: 3, numbering: none, "Detailed coordination reference")
The theory lessons use short diagrams. This expanded sequence shows the interleaving
behind the optimistic and pessimistic write strategies.

#figure-evidence("../assets/fig-optimistic-vs-pessimistic.svg", "sequence", "conceptual illustration",
  "A conditional write detects change; a lock orders writers before their decision.",
  [Optimistic writers reject a mismatched version and re-read. Pessimistic writers wait before reading.
   Both still recheck the business invariant.], image-width: 82%)

#heading(level: 3, numbering: none, "Reproduce the book")
```bash
book/build.sh
cat book/dist/build-manifest.json
```
The pinned Podman build records its source commit, source hash, registry digest, image and PDF hash.

#heading(level: 3, numbering: none, "Reproduce a measured run")
- Open the cited claim card and follow its reports, input digest and run tags.
- A run tag identifies producing code; committed results live separately in the study's results directory.
- Use the named study runner. It acquires the benchmark lock; never start a competing database manually.

The full claim cards below are the audit trail. Their detailed limits are separate from the teaching chapters.
