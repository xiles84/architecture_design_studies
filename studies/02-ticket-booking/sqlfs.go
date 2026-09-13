// Package ticketbooking embeds the SQL catalogue for study 02 so the benchmark
// binary is self-contained, and the SQL that produced a result is provably the
// SQL that shipped with the binary.
package ticketbooking

import "embed"

//go:embed sql
var SQL embed.FS
