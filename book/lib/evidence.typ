// The book reads exactly one evidence registry: the active versioned package
// under book/evidence/<version>/. This module loads it and renders an index of
// the claims that are allowed to appear as evidence.
//
// The version is a build input, derived once, and never written into prose. That
// is deliberate: Edition 1 hard-coded a v2 registry path in one chapter while the
// cover named v3, which is the provenance defect a reader notices first.
// `book/evidence/check.sh` fails the build if a version literal reappears outside
// this derivation.

// The one canonical value. build.sh passes it; the default keeps an editor
// preview compiling against the committed active package.
#let registry-version = sys.inputs.at("registry_version", default: "v4")

// The source-relative path `json` needs. A path without a leading slash is
// resolved against THIS file, so the book compiles whether the Typst root is
// `book/` (as build.sh sets) or the repository (as an editor's language server
// defaults to). A leading-slash path is root-relative, which made preview fail
// with "file not found" before this change.
#let registry-path = "../evidence/" + registry-version + "/claims.json"

// The repo-relative path a human reads. Derived from the same value, never typed.
#let registry-label = "book/evidence/" + registry-version + "/claims.json"

#let loaded = {
  // `json` needs the file to exist at compile time; the build copies the book
  // directory wholesale, so the registry is always beside this file.
  json(registry-path,)
}

#let all-claims = loaded.claims

// A claim may appear as evidence only when it is still standing. v4 states this
// explicitly with `status`; a package without the key (v1-v3) is treated as all
// active so this module keeps working on the frozen predecessors.
#let claim-status(c) = c.at("status", default: "active")
#let is-renderable(c) = claim-status(c) != "retired" and claim-status(c) != "superseded"

#let claims = all-claims.filter(is-renderable)

#let family-names = {
  let m = (:)
  for f in loaded.taxonomy.families {
    m.insert(f.id, f.name)
  }
  m
}

#let active-claims() = claims

// The gap vocabulary. A schema limitation ("no representation can answer this")
// is a different epistemic object from a regime nobody ran, and the registry
// says which one a gap is; the page must not flatten them into "Coverage gap".
#let gap-kind-label(k) = {
  if k == "schema_limitation" { "schema limitation" }
  else if k == "instrument" { "instrument gap" }
  else if k == "unstable_measurement" { "unstable measurement" }
  else { "coverage gap" }
}

#let gap-kind-text(c) = {
  let kinds = c.at("gap_kind", default: ())
  if kinds.len() == 0 { none } else { kinds.map(gap-kind-label).join(" + ") }
}

#let claim-by-id(id) = {
  let found = none
  for c in all-claims {
    if c.claim_id == id { found = c }
  }
  if found == none {
    // A retired claim is not in claims[] any more, so look for it in the
    // retirement chain: an author who cites it needs the successor's id, not a
    // generic "unknown claim".
    let successor = none
    for r in loaded.retired_predecessor_claims {
      if r.claim_id == id {
        let by = r.at("superseded_by", default: ())
        if by.len() > 0 { successor = by.join(", ") }
      }
    }
    if successor != none {
      panic("claim " + id + " is retired in " + registry-label + "; cite its successor instead: " + successor)
    }
    panic("no claim with id " + id + " in " + registry-label + "; do not write a number without a claim")
  }
  if not is-renderable(found) {
    panic("claim " + id + " is " + claim-status(found) + " in " + registry-label + " and may not be rendered as active evidence; cite its successor instead")
  }
  found
}

// One trial line, pluralised per noun. Edition 1 printed "Trials: 1 whole-run
// replications, 1 fresh-load trials per cell, 1 inner iterations" on every card.
// Declared before the card that uses it: Typst resolves names in order.
#let plural(n, singular, plural-form) = if n == 1 { singular } else { plural-form }

#let trials-line(t) = [
  Trials: #t.whole_run_replications #plural(t.whole_run_replications, "whole-run replication", "whole-run replications")
  · #t.fresh_load_trials_per_cell #plural(t.fresh_load_trials_per_cell, "fresh-load trial", "fresh-load trials") per cell
  · #t.inner_iterations #plural(t.inner_iterations, "inner iteration", "inner iterations")
]

// A confidence card rendered straight from the registry, so the page cannot
// drift from the data. This is the only way a number enters the prose.
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
    // The id sits on its own line, full width: a claim id is the handle a reader
    // copies, and putting a long meta note beside it squeezed the id column until
    // the id broke across lines. build.sh's claims_indexed gate catches that.
    #text(weight: "bold", fill: rgb("#2b4b8f"))[#c.claim_id]
    #v(1pt)
    #text(size: 8pt, fill: rgb("#5b6470"))[
      #family-names.at(c.family, default: c.family) · #c.strength
      #if claim-status(c) == "partially_superseded" [
        · partially superseded by #c.superseded_by.join(", ")
      ]
      #if c.at("gap_kind", default: ()).len() > 0 [
        · #gap-kind-text(c)
      ]
    ]
    #v(2pt)
    #c.statement
    #if claim-status(c) == "partially_superseded" [
      #v(3pt)
      #text(size: 8pt, fill: rgb("#5b6470"))[
        Still open: #c.remaining_dimensions.join("; ")
      ]
    ]
    #v(3pt)
    #text(size: 8pt, fill: rgb("#5b6470"))[
      #trials-line(c.trials)
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

// The index of every claim that may appear as evidence, grouped by the fixed six
// families. This is the reader-facing map from prose to data; a rendered page
// proves the registry compiled into the PDF.
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
    ..loaded.retired_predecessor_claims.map(r => [
      #raw(r.claim_id) — #r.reason
      #if r.at("closed_dimensions", default: ()).len() > 0 [
        _Closed:_ #r.closed_dimensions.join("; ").
      ]
    ]),
  )
}
