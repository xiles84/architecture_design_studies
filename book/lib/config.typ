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
#let book-edition = "Edition 1 (draft)"

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
  show heading.where(level: 1): it => {
    pagebreak(weak: true)
    block(above: 0.6em, below: 0.8em)[
      #set text(size: 19pt, weight: "bold", fill: palette.accent)
      #it
      #v(1mm)
      #line(length: 100%, stroke: 1.2pt + palette.accent)
    ]
  }
  show heading.where(level: 2): it => block(above: 0.9em, below: 0.5em)[
    #set text(size: 14pt, weight: "bold", fill: palette.ink)
    #it
  ]
  show heading.where(level: 3): it => block(above: 0.7em, below: 0.4em)[
    #set text(size: 11.5pt, weight: "bold", fill: palette.muted)
    #it
  ]
  // Accessible code: monospace, generous leading, a labelled surface, and no
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

// A family / variant / topology / scenario marker used at the head of a section.
#let marker(kind, name) = block(below: 0.6em)[
  #set text(size: 8.5pt, tracking: 0.8pt, fill: palette.muted)
  #upper(kind) · #upper(name)
]

// An evidence card names a v2 claim and is the only way a number may enter the
// prose. `strength` and `family` must match the registry; the validator does not
// check the book's copies, so writers copy them by hand *and* cite the id.
#let evidence-card(id, family, strength, statement, limits: (), confounds: ()) = block(
  width: 100%,
  inset: 9pt,
  radius: 3pt,
  stroke: 0.6pt + palette.rule,
  fill: palette.surface,
)[
  #set text(size: 9pt)
  #grid(
    columns: (1fr, auto),
    align: (left, right),
    [#text(weight: "bold", fill: palette.accent)[#id]],
    [#text(size: 8pt, fill: palette.muted)[#family · #strength]],
  )
  #v(2pt)
  #statement
  #if limits.len() > 0 [
    #v(3pt)
    #text(size: 8pt, fill: palette.muted)[Limits: #limits.join("; ")]
  ]
  #if confounds.len() > 0 [
    #v(2pt)
    #text(size: 8pt, fill: palette.gap)[Confounds: #confounds.join(", ")]
  ]
]

#let source(url, label: none) = link(url)[#if label == none { url } else { label }]

// A book figure. Every figure carries TWO independent labels — its form
// (structure / sequence / state) and its evidence status (conceptual illustration /
// implemented design contract / negative control / observed result) — plus an
// adjacent text equivalent for a reader who cannot see the image. An
// observed-result figure names the active claim it illustrates; it never embeds a
// measured number of its own.
#let figure-evidence(asset, form, status, caption, equivalent, claim: none) = block(
  width: 100%,
  breakable: false,
  above: 1em,
  below: 1em,
)[
  #figure(
    image(asset, width: 100%),
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
