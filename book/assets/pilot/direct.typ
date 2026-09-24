// Pilot variant (a): import the canonical study SVG directly by repository-root
// path. This only compiles when Typst's --root is the repository, because a
// leading-slash path is root-relative. Kept as the recorded negative control for
// the import-mechanism pilot in book/FIGURES.md.
#set page(width: 90mm, height: 70mm, margin: 4mm)
#image("/studies/01-charity-tree/diagrams/rendered/00_overview.svg", width: 100%)
