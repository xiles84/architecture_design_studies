// Visual system for the data-architecture reference book.
//
// Everything visual is defined here and nowhere else: page geometry, the type
// scale, the running header/footer, and the callout/label vocabulary the
// chapters share. Chapters import this file and use the functions; a chapter
// that sets its own fonts or margins would break the one-look rule.

#let palette = (
  ink: rgb("#111417"),
  muted: rgb("#5b6470"),
  rule: rgb("#d5d9de"),
  direct: rgb("#1f6f4a"),
  analogy: rgb("#8a5a00"),
  gap: rgb("#9a2f3b"),
  accent: rgb("#2b4b8f"),
  surface: rgb("#f5f6f8"),
)

// Typst 0.15 ships these fonts, so the PDF embeds them without any host font
// installation. Changing this tuple changes the PDF; the build manifest records
// the pinned image digest that supplies them.
#let fonts = ("Libertinus Serif", "New Computer Modern")
#let mono-font = "DejaVu Sans Mono"

#let book-title = "Data Architecture Reference"
#let book-edition = "Edition 3 preview 5"

// ---- page geometry and running matter --------------------------------

#let page-setup(body) = {
  set page(
    paper: "a4",
    margin: (top: 22mm, bottom: 20mm, inside: 24mm, outside: 20mm),
    numbering: "1",
    header: context {
      let n = counter(page).get().first()
      if n > 1 [
        #set text(size: 8.5pt, fill: palette.muted)
        #grid(
          columns: (1fr, auto),
          align: (left, right),
          [#book-title · #book-edition],
          [Evidence-backed reference · draft],
        )
        #v(1mm)
        #line(length: 100%, stroke: 0.4pt + palette.rule)
      ]
    },
    footer: context {
      let n = counter(page).get().first()
      if n > 1 [
        #line(length: 100%, stroke: 0.4pt + palette.rule)
        #v(1mm)
        #set text(size: 8.5pt, fill: palette.muted)
        #grid(
          columns: (1fr, auto),
          align: (left, right),
          [Do not quote a number without its evidence card.],
          [#n],
        )
      ]
    },
  )
  set text(font: fonts, size: 10.5pt, fill: palette.ink, lang: "en")
  set par(justify: true, leading: 0.72em)
  set heading(numbering: "1.1")
  // Give each level enough room to read as a new thought. These gaps apply in every
  // chapter, rather than relying on local vertical spacers that drift as text changes.
  show heading.where(level: 1): it => {
    pagebreak(weak: true)
    block(above: 0.8em, below: 1.15em)[
      #set text(size: 19pt, weight: "bold", fill: palette.accent)
      #it
      #v(1mm)
      #line(length: 100%, stroke: 1.2pt + palette.accent)
    ]
  }
  show heading.where(level: 2): it => {
    block(above: 1.3em, below: 0.7em)[
    #set text(size: 14pt, weight: "bold", fill: palette.ink)
    #it
    ]
  }
  show heading.where(level: 3): it => block(above: 0.95em, below: 0.55em)[
    #set text(size: 14pt, weight: "bold", fill: palette.accent)
    #it
  ]
  // Accessible code: monospace, generous leading, a labelled surface, and no
  show heading.where(level: 4): it => block(above: 0.8em, below: 0.4em)[
    #text(size: 10.5pt, weight: "bold", fill: palette.accent)[#it]
  ]
  // forced hyphenation of commands. `raw` is the only block code the book uses.
  show raw.where(block: true): it => block(
    width: 100%,
    inset: 8pt,
    radius: 3pt,
    stroke: 0.5pt + palette.rule,
    fill: palette.surface,
  )[
    #set text(font: mono-font, size: 9pt, hyphenate: false)
    #it
  ]
  body
}

// ---- callouts and evidence vocabulary --------------------------------

#let callout(title, color, body) = block(
  width: 100%,
  // A callout's label states how the reader must interpret its body. Keep that
  // label with the complete warning instead of splitting the box at a page foot.
  breakable: false,
  inset: 8pt,
  radius: 3pt,
  stroke: 0.6pt + color,
  fill: color.transparentize(94%),
)[
  #set text(size: 9.5pt)
  #text(weight: "bold", fill: color)[#title]
  #v(2pt)
  #body
]

#let direct(body) = callout("Direct evidence", palette.direct, body)
#let analogy(body) = callout("Analogy — not measured here", palette.analogy, body)
#let gap(body) = callout("Coverage gap", palette.gap, body)

// A mechanism card: one named coordination mechanism, with the race it closes kept
// separate from what it does not solve.
//
// Why a card and not a table row: Edition 1 put five mechanisms in a three-column
// table declared as `columns: (auto, 1fr, auto)`. The long evidence-status cell
// claimed almost the whole text width, the middle column collapsed to a sliver, and
// five pages of hyphenated fragments — "se - ri - alises", "asyn - chro - nous" —
// carried the book's most important conceptual section. A card cannot collapse that
// way, and the mandatory `not-solved` line stops "this mechanism" being read as
// "this problem is solved".
#let mechanism-card(name, closes, how, guarantee, not-solved, evidence) = block(
  width: 100%,
  inset: 8pt,
  radius: 3pt,
  stroke: 0.5pt + palette.rule,
  fill: palette.surface,
  above: 0.7em,
  below: 0.7em,
)[
  #set text(size: 9pt)
  #text(weight: "bold", size: 10.5pt, fill: palette.accent)[#name]
  #v(3pt)
  #grid(
    columns: (56pt, 1fr),
    row-gutter: 3pt,
    column-gutter: 8pt,
    align: (left, left),
    text(size: 8pt, fill: palette.muted)[Race closed], closes,
    text(size: 8pt, fill: palette.muted)[Mechanism], how,
    text(size: 8pt, fill: palette.muted)[Guarantee], guarantee,
    text(size: 8pt, fill: palette.muted)[Does not solve], not-solved,
  )
  #v(3pt)
  #text(size: 8pt, fill: palette.muted)[Evidence: #evidence]
]

// Mechanisms fall into two layers, and conflating them is how "the database lock
// makes the cache consistent" gets written. A source mechanism arbitrates the
// mutation; a copy mechanism orders what leaves the transaction boundary.
#let mechanism-group(title-text, note) = block(above: 1em, below: 0.1em)[
  #text(weight: "bold", size: 11pt)[#title-text]
  #v(1pt)
  #text(size: 8.5pt, fill: palette.muted)[#note]
]

// A family / variant / topology / scenario marker used at the head of a section.
#let marker(kind, name) = block(below: 0.6em)[
  #set text(size: 8.5pt, tracking: 0.3pt, fill: palette.muted)
  #upper(kind) · #upper(name)
]

#let source(url, label: none) = link(url)[#if label == none { url } else { label }]

// A book figure. Every figure carries TWO independent labels — its form
// (structure / sequence / state) and its evidence status (conceptual illustration /
// implemented design contract / negative control / observed result) — plus an
// adjacent text equivalent for a reader who cannot see the image. An
// observed-result figure names the active claim it illustrates; it never embeds a
// measured number of its own.
// `image-width` exists because a tall sequence diagram rendered at the full text
// width can be taller than a page and overflow into the running footer. A tall
// figure is centred at a reduced width instead; its caption and text equivalent
// stay full width.
#let figure-evidence(asset, form, status, caption, equivalent, claim: none, image-width: 100%) = block(
  width: 100%,
  breakable: false,
  above: 1em,
  below: 1em,
)[
  #figure(
    image(asset, width: image-width),
    caption: {
      text(weight: "bold")[#caption]
      linebreak()
      text(size: 8.5pt, fill: palette.muted)[
        Form: #form · Evidence: #status#if claim != none [ · illustrates #claim]
      ]
    },
  )
  #block(inset: (left: 8pt), stroke: (left: 0.8pt + palette.rule))[
    #text(size: 8.5pt)[*Text equivalent.* #equivalent]
  ]
]

// The progressive structure every family/variant section follows (REQUEST.md).
#let progressive-structure(depth: 1) = {
  heading(level: depth, "Quick choice in plain language")
  heading(level: depth, "Where it thrives and where it perishes")
  heading(level: depth, "Common scenarios")
  heading(level: depth, "Reusable mechanism")
  heading(level: depth, "Evidence and reproduction")
}

// The registry is data, so its prose keeps ASCII hyphens in numeric ranges, as a
// JSON file should. The page typesets a range with an en dash. This converts one to
// the other at render time, so the chapter prose and the generated cards agree
// without editing the evidence: the registry stays the verbatim record, and nothing
// here changes a digit.
#let typeset-ranges(t) = t.replace(
  regex("([0-9])-([0-9])"),
  m => m.captures.at(0) + "–" + m.captures.at(1),
)

// ---- Edition 3: the concept-first didactic components ------------------
//
// Owner decision 2026-10-01 (book/STYLE.md). Every chapter teaches one
// concept, the case studies show when it fits, and the text stays short.
// These components exist so a writer cannot accidentally publish a chapter
// that breaks the style: each one validates its input and panics with a
// message that names the rule, the value it got, and where to look. A panic
// is the right failure here — a rule nobody can check is a rule that drifts.


// Plain text of a content value, for the word-count guards below. A content
// sequence (for example one containing an em dash) has no `.text`, so walk
// the children and join what they say.
#let content-text(c) = {
  if type(c) == str { c }
  else if type(c) == content {
    let f = c.fields()
    if f.keys().contains("text") and type(f.text) == str { f.text }
    else if f.keys().contains("children") { f.children.map(content-text).join(" ") }
    else { "" }
  } else { "" }
}

#let word-count(c) = {
  let t = content-text(c).trim()
  if t == "" { 0 } else { t.split(" ").filter(w => w != "").len() }
}

#let enforce(cond, message) = if not cond { panic(message) }

// 1. Opener. `idea` is the key idea, at most 25 words (STYLE.md, template 1).
#let concept-opener(number, title, idea) = {
  let n = word-count(idea)
  enforce(n <= 25, "concept-opener " + title + ": the key idea has " + str(n) +
    " words; STYLE.md allows 25. Shorten it, or move the detail to the prose below.")
  // The number belongs above the title, but a level-1 heading breaks the page first
  // (page-setup's show rule), so an eyebrow emitted before it lands on the previous
  // page. It is emitted after the heading instead, under the rule, where it reads as
  // the chapter's label without risking the contents entry.
  heading(level: 1, numbering: none, title)
  block(above: 0em, below: 0.2em)[
    #text(size: 8.5pt, tracking: 0.4pt, fill: palette.muted)[CHAPTER #number]
  ]
  block(
    width: 100%,
    inset: 8pt,
    radius: 3pt,
    stroke: 0.6pt + palette.accent,
    fill: palette.accent.transparentize(95%),
    below: 0.9em,
  )[
    #set text(size: 10.5pt)
    #text(weight: "bold", fill: palette.accent)[Key idea: ]#idea
  ]
}

// 2. Fit table. STYLE.md wants 3-4 rows per side; more means the chapter is
// hedging and the reader leaves with nothing to decide on.
#let fit-table(use: (), avoid: ()) = {
  enforce(use.len() >= 3 and use.len() <= 4,
    "fit-table: \"fits when\" needs 3-4 rows, got " + str(use.len()) + " (STYLE.md, template 5)")
  enforce(avoid.len() >= 3 and avoid.len() <= 4,
    "fit-table: \"does not fit when\" needs 3-4 rows, got " + str(avoid.len()) + " (STYLE.md, template 5)")
  let rows = ()
  for i in range(calc.max(use.len(), avoid.len())) {
    rows.push(if i < use.len() { use.at(i) } else { [] })
    rows.push(if i < avoid.len() { avoid.at(i) } else { [] })
  }
  block(width: 100%, breakable: false, above: 0.7em, below: 0.7em)[
    #set text(size: 9.5pt)
    #table(
      columns: (1fr, 1fr),
      inset: 7pt,
      stroke: 0.4pt + palette.rule,
      align: (left + top, left + top),
      table.header(
        text(weight: "bold", fill: palette.direct)[It fits when],
        text(weight: "bold", fill: palette.gap)[It does not fit when],
      ),
      ..rows,
    )
  ]
}

// 3. Cost table: five rows, one phrase each, reads first because that is what
// the decision buys.
#let cost-table(reads: none, writes: none, storage: none, correctness: none, operations: none) = {
  let cells = (
    ("Reads", reads), ("Writes", writes), ("Storage", storage),
    ("Correctness", correctness), ("Operations", operations),
  )
  for (name, v) in cells {
    enforce(v != none, "cost-table: the \"" + name + "\" row is missing; STYLE.md gives the table five fixed rows")
    enforce(word-count(v) <= 12,
      "cost-table: the \"" + name + "\" row has " + str(word-count(v)) + " words; " +
      "STYLE.md wants one phrase per row (12 words is the ceiling here)")
  }
  block(width: 100%, above: 0.7em, below: 0.7em)[
    #set text(size: 9.5pt)
    #table(
      columns: (auto, 1fr),
      inset: 6pt,
      stroke: 0.4pt + palette.rule,
      align: (left + top, left + top),
      ..cells.map(((name, v)) => (
        text(weight: "bold", fill: palette.muted, size: 8.5pt)[#name], v,
      )).flatten(),
    )
  ]
}

// 4. Case box: one study, one lesson about this concept. `verdict` is the
// only controlled vocabulary; the colours match the callouts a reader has
// already met, so "does not fit" is the same red as a gap.
#let case-verdicts = ("fits", "does not fit", "depends")
#let case-box(study: "", verdict: "", situation: none, compared: none, result: none,
              lesson: none, watch-out: none, evidence: none) = {
  enforce(verdict in case-verdicts,
    "case-box " + study + ": verdict \"" + verdict + "\" is not one of " +
    case-verdicts.map(str).join(", ") + " (STYLE.md, case box)")
  for (name, v) in (("situation", situation), ("compared", compared), ("result", result),
                    ("lesson", lesson), ("watch-out", watch-out), ("evidence", evidence)) {
    enforce(v != none, "case-box " + study + ": the \"" + name + "\" field is missing; " +
      "STYLE.md gives the case box eight fields")
  }
  let colour = if verdict == "fits" { palette.direct }
    else if verdict == "does not fit" { palette.gap }
    else { palette.analogy }
  block(
    width: 100%, inset: 9pt, radius: 3pt, breakable: false,
    stroke: 0.6pt + palette.rule, fill: palette.surface,
    above: 0.8em, below: 0.8em,
  )[
    #set text(size: 9.5pt)
    #grid(
      columns: (1fr, auto),
      align: (left, right),
      text(weight: "bold", fill: palette.accent)[#study],
      box(inset: (x: 5pt, y: 2pt), radius: 2pt, fill: colour.transparentize(85%))[
        #text(size: 8.5pt, weight: "bold", fill: colour)[#upper(verdict)]
      ],
    )
    #v(4pt)
    #grid(
      columns: (58pt, 1fr), row-gutter: 4pt, column-gutter: 8pt, align: (left + top, left + top),
      text(size: 8.5pt, fill: palette.muted)[Situation], situation,
      text(size: 8.5pt, fill: palette.muted)[Compared], compared,
      text(size: 8.5pt, fill: palette.muted)[Result], result,
      text(size: 8.5pt, fill: palette.muted)[Lesson], lesson,
      text(size: 8.5pt, fill: palette.muted)[Watch out], watch-out,
    )
    #v(4pt)
    #evidence
  ]
}

// 5. Remember: at most three bullets. A fourth means the chapter has not
// decided what it is about.
#let remember(..bullets) = {
  let items = bullets.pos()
  enforce(items.len() >= 1 and items.len() <= 3,
    "remember: needs 1-3 bullets, got " + str(items.len()) + " (STYLE.md, template 9)")
  block(width: 100%, inset: 8pt, radius: 3pt, stroke: 0.5pt + palette.rule,
        fill: palette.surface, above: 0.8em, below: 0.8em)[
    #set text(size: 9.5pt)
    #text(weight: "bold", fill: palette.accent)[Remember]
    #v(3pt)
    #list(..items)
  ]
}

// 6. Checklist table: exactly eight rows (methodology 21), each ending in an
// evidence status. The status is the answer's honesty label, so an unknown
// one must stop the build rather than render unlabelled.
#let checklist-statuses = ("direct", "mechanism", "analogy", "gap")
#let checklist-status-label(s) = {
  if s == "direct" { "direct" }
  else if s == "mechanism" { "mechanism" }
  else if s == "analogy" { "analogy" }
  else if s == "gap" { "gap" }
  else { panic("checklist-table: unknown status \"" + s + "\"; use one of " +
    checklist-statuses.map(str).join(", ") + " (book/concepts/CHECKLIST.md)") }
}
#let checklist-table(rows) = {
  enforce(rows.len() == 8, "checklist-table: needs exactly 8 rows, got " + str(rows.len()) +
    " (methodology 21; the eight questions are fixed)")
  let body = ()
  for i in range(rows.len()) {
    let row = rows.at(i)
    enforce(row.len() == 3, "checklist-table row " + str(i + 1) +
      ": expected (question, answer, status), got " + str(row.len()) + " items")
    let (question, answer, status) = row
    let label = checklist-status-label(status)
    enforce(word-count(answer) <= 15, "checklist-table row " + str(i + 1) + " (\"" + question +
      "\"): the answer has " + str(word-count(answer)) + " words; STYLE.md allows 15")
    body.push(text(weight: "bold", size: 9pt)[#question])
    body.push(text(size: 9pt)[#answer])
    body.push(text(size: 8pt, fill: palette.muted)[#label])
  }
  block(width: 100%, above: 0.7em, below: 0.7em)[
    #set text(size: 9.5pt)
    #table(
      columns: (1.1fr, 1.6fr, auto),
      inset: 6pt,
      stroke: 0.4pt + palette.rule,
      align: (left + top, left + top, left + top),
      table.header([Question], [Answer], [Evidence]),
      ..body,
    )
  ]
}

// 7. Part opener: one page per part, so the structure of the book is visible
// in the contents and on the page.
#let part-opener(number, title, blurb) = {
  pagebreak(weak: true)
  block(above: 12mm)[
    #text(size: 8.5pt, tracking: 0.4pt, fill: palette.muted)[PART #number]
    #v(2pt)
    #text(size: 22pt, weight: "bold", fill: palette.accent)[#title]
    #v(2mm)
    #line(length: 60%, stroke: 1.2pt + palette.accent)
    #v(3mm)
    #block(width: 85%)[
      #set text(size: 10.5pt)
      #blurb
    ]
  ]
}
