// Pilot variant (b): import the generated, read-only book asset with a
// file-relative path. This compiles under either Typst root, because a path
// without a leading slash resolves against the file that names it.
#set page(width: 90mm, height: 70mm, margin: 4mm)
#image("../pilot-00-overview.svg", width: 100%)
