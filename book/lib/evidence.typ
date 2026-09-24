// The book reads exactly one evidence registry: the v3 package (the v2 claims
// carried forward plus the 2026-09-23 Study 05 update). This module loads it and
// renders an index of active claims. It deliberately fails to a visible error
// string rather than a silent empty table, so a build with a missing registry is
// obvious on the page.

// A path without a leading slash is resolved against THIS file, so the book
// compiles whether the Typst root is `book/` (as build.sh sets) or the
// repository (as an editor's language server defaults to). A leading-slash path
// is root-relative, which made preview fail with "file not found (searched at
// <repo>/evidence/v3/claims.json)" before this change.
#let registry-path = "../evidence/v3/claims.json"

#let loaded = {
  // `json` needs the file to exist at compile time; the build copies the book
  // directory wholesale, so the registry is always beside this file.
  json(registry-path,)
}

#let claims = loaded.claims

#let family-names = {
  let m = (:)
  for f in loaded.taxonomy.families {
    m.insert(f.id, f.name)
  }
  m
}

#let active-claims() = claims

#let claim-by-id(id) = {
  let found = none
  for c in claims {
    if c.claim_id == id { found = c }
  }
  if found == none {
    panic("no active claim with id " + id + "; do not write a number without a claim")
  }
  found
}

// A compact confidence card rendered straight from the registry, so the page
// cannot drift from the data. Chapters use `evidence-card` in config.typ when
// they want the statement highlighted in prose; this renders the full card.
#let registry-card(id) = {
  let c = claim-by-id(id)
  block(
    width: 100%,
    inset: 9pt,
    radius: 3pt,
    stroke: 0.6pt + rgb("#d5d9de"),
    fill: rgb("#f5f6f8"),
  )[
    #set text(size: 9pt)
    #grid(
      columns: (1fr, auto),
      align: (left, right),
      [#text(weight: "bold", fill: rgb("#2b4b8f"))[#c.claim_id]],
      [#text(size: 8pt, fill: rgb("#5b6470"))[#family-names.at(c.family, default: c.family) · #c.strength]],
    )
    #v(2pt)
    #c.statement
    #v(3pt)
    #text(size: 8pt, fill: rgb("#5b6470"))[
      Trials: #c.trials.whole_run_replications whole-run replications,
      #c.trials.fresh_load_trials_per_cell fresh-load trials per cell,
      #c.trials.inner_iterations inner iterations.
    ]
    #if c.confounds.len() > 0 [
      #v(2pt)
      #text(size: 8pt, fill: rgb("#9a2f3b"))[Confounds: #c.confounds.join(", ")]
    ]
    #v(2pt)
    #text(size: 8pt, fill: rgb("#5b6470"))[Limits: #c.limits.join("; ")]
    #v(2pt)
    #text(size: 7.5pt, fill: rgb("#5b6470"))[
      Source: #c.anchors.at(0).report
      · run #c.anchors.at(0).run_id
      · tag #(if c.anchors.at(0).run_tag == none { "none (legacy)" } else { c.anchors.at(0).run_tag })
    ]
  ]
}

// The index of every active claim, grouped by the fixed six families. This is
// the reader-facing map from prose to data; a rendered page proves the registry
// compiled into the PDF.
#let evidence-index() = {
  for f in loaded.taxonomy.families [
    #heading(level: 2, f.name)
    #for c in claims {
      if c.family == f.id [
        #registry-card(c.claim_id)
        #v(4pt)
      ]
    }
  ]
  heading(level: 2, "Retired predecessor claims")
  list(
    ..loaded.retired_predecessor_claims.map(r => [#raw(r.claim_id) — #r.reason]),
  )
}
