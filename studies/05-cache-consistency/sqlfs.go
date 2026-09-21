// Package cacheconsistency embeds the SQL catalogue for study 05 so the benchmark
// binary is self-contained, and the SQL that produced a result is provably the SQL
// that shipped with the binary.
package cacheconsistency

import "embed"

//go:embed sql
var SQL embed.FS
