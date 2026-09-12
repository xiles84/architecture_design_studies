// Package charitytree embeds the SQL catalogue for study 01 so the benchmark
// binary is self-contained: the container needs no bind mount to run, and the
// SQL that produced a result is guaranteed to be the SQL that shipped with the
// binary.
package charitytree

import "embed"

//go:embed sql
var SQL embed.FS
