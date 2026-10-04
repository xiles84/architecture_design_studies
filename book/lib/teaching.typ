#import "config.typ": *
#let flow-picture(nodes, caption, equivalent) = block(breakable: false)[
  #figure(
    grid(columns: (1fr, auto, 1fr, auto, 1fr), align: center + horizon,
      block(inset: 8pt, stroke: 0.6pt + palette.accent, fill: palette.surface)[#nodes.at(0)],
      [→],
      block(inset: 8pt, stroke: 0.6pt + palette.accent, fill: palette.surface)[#nodes.at(1)],
      [→],
      block(inset: 8pt, stroke: 0.6pt + palette.accent, fill: palette.surface)[#nodes.at(2)],
    ),
    caption: [#caption \ #text(size: 8.5pt, fill: palette.muted)[Form: sequence · Evidence: conceptual illustration]],
  )
  #text(size: 8.5pt)[*Text equivalent.* #equivalent]
]
#let strategy(title, id, idea, problem, steps, fits, avoids, costs, picture, check, case, measured: true) = {
  pagebreak(weak: true)
  heading(level: 3, numbering: none, title)
  label(id)
  marker("strategy", title)
  callout("Key idea", palette.accent)[#idea]
  heading(level: 4, numbering: none, "The problem")
  problem
  picture
  heading(level: 4, numbering: none, "How it works")
  list(..steps)
  heading(level: 4, numbering: none, "Thrive or die: the expected boundary")
  fit-table(use: fits, avoid: avoids)
  block(breakable: false)[
    #heading(level: 4, numbering: none, "What it costs")
    #cost-table(reads: costs.at(0), writes: costs.at(1), storage: costs.at(2), correctness: costs.at(3), operations: costs.at(4))
  ]
  block(breakable: false)[
    #heading(level: 4, numbering: none, "In practice")
    #if measured [Follow #link(label(case))[the matching study case] for tested variations and registered results.]
    else [These are design expectations. #link(label(case))[The study-case coverage map] identifies what has not been measured.]
    #text(size: 9pt, fill: palette.muted)[A plausible boundary is not a measured crossover.]
  ]
  block(breakable: false)[
    #heading(level: 4, numbering: none, "The checklist")
    #checklist-table((
      ("What is the invariant?", check.at(0), "mechanism"),
      ("Where does the information live?", check.at(1), "mechanism"),
      ("Where is the contention?", check.at(2), "mechanism"),
      ("What work moved from reads to writes?", check.at(3), "mechanism"),
      ("What state is derived?", check.at(4), "mechanism"),
      ("What can become stale?", check.at(5), "mechanism"),
      ("Under which workload does it stop working?", check.at(6), "gap"),
      ("Which experiment would prove it?", check.at(7), "gap"),
    ))
  ]
}
