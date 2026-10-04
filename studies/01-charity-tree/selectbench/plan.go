package selectbench

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
)

// Plan handling for the SELECT families.
//
// A timing is only interpretable beside the plan that produced it, and two cells
// may only be called "the same plan" when they agree on the *structure*: node
// types, relations and index names, and materialization boundaries. Per-node
// costs, row estimates, actual times and buffer counts change between runs of the
// same plan and must not enter a signature, or a difference would be reported
// where the plan never changed.
//
// The signature is a hash of the normalized structure; it is evidence that two
// cells ran the same shape, never a claim that they are semantically equivalent
// (two different statements can share a shape and mean different things).

// volatile removes the parenthetical cost/actual metadata PostgreSQL prints
// after each node, e.g. "(cost=0.00..1.07 rows=7 width=8) (actual time=...)".
var volatile = regexp.MustCompile(`\((?:cost|actual)[^)]*\)`)

// nodeExtras are per-node lines that depend on the bound parameters, the plan's
// row estimates or buffer state rather than on the plan's shape. They are kept
// in Normalize (they are part of what the plan means) and dropped in Structure
// (they are not part of its shape).
var nodeExtras = regexp.MustCompile(`(?i)^\s*(Filter|Index Cond|Recheck Cond|Rows Removed by [A-Za-z ]+|Heap Fetches|Sort Key|Sort Method|Hash Cond|Join Filter|Merge Cond|Buckets|Batches|Memory Usage|Peak Memory Usage|One-Time Filter|Index Recheck|Rows Removed by Filter):`)

// runLevel removes the statements that never describe plan shape.
var runLevel = regexp.MustCompile(`(?i)^\s*(Planning Time|Execution Time|JIT|Buffers|I/O Timings|WAL|Settings):`)

// SQL aliases are labels, not an access-path or join-tree change. Keep the
// relation/index identity while normalizing aliases in our simple catalogue.
var scanAlias = regexp.MustCompile(`(Scan(?: Backward)?(?: using [^ ]+)? on [^ ]+) [^ ]+$`)

// Normalize returns the plan's step lines with volatile metadata removed, in
// execution (print) order. Semantic lines such as Filter are kept.
func Normalize(explain string) []string {
	var out []string
	for _, raw := range strings.Split(strings.ReplaceAll(explain, "\r\n", "\n"), "\n") {
		if runLevel.MatchString(raw) {
			continue
		}
		line := strings.TrimSpace(volatile.ReplaceAllString(raw, ""))
		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

// Structure preserves node order and indentation: outer/inner roles and nested
// joins are part of a plan's shape. A sorted node inventory loses that evidence.
func Structure(explain string) []string {
	var out []string
	for _, raw := range strings.Split(strings.ReplaceAll(explain, "\r\n", "\n"), "\n") {
		if nodeExtras.MatchString(raw) || runLevel.MatchString(raw) || strings.TrimSpace(raw) == "Planning:" {
			continue
		}
		line := strings.Join(strings.Fields(volatile.ReplaceAllString(raw, "")), " ")
		line = scanAlias.ReplaceAllString(line, "$1")
		if line == "" {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " \t"))
		out = append(out, strings.Repeat(" ", indent)+line)
	}
	return out
}

func hashLines(lines []string) string {
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(sum[:])[:16]
}

// Signature identifies the full normalized plan (shape plus filters and index
// conditions).
func Signature(explain string) string { return hashLines(Normalize(explain)) }

// StructureSignature identifies the plan's shape only.
func StructureSignature(explain string) string { return hashLines(Structure(explain)) }

// SameStructure reports whether two plans have the same shape. Two statements
// that agree here ran the same node/index structure; it says nothing about
// whether they answer the same question — that is the oracle's job.
func SameStructure(a, b string) bool { return StructureSignature(a) == StructureSignature(b) }

// AccessMetrics is the engine-reported access-path evidence parsed from one plan.
// It exists so the report can show the physical regime each cell ran under
// (review M1): a hash-sharded YugabyteDB index produces `Storage Filter` over
// the whole table, a range index produces `Index Cond` and a small
// `Storage Rows Scanned`.
type AccessMetrics struct {
	StorageRowsScanned  int  `json:"storage_rows_scanned,omitempty"`
	IndexRowsScanned    int  `json:"index_rows_scanned,omitempty"`
	PlanIndexScan       bool `json:"plan_index_scan,omitempty"`
	PlanSeqScan         bool `json:"plan_seq_scan,omitempty"`
	BufSharedHit        int  `json:"buf_shared_hit,omitempty"`
	BufSharedRead       int  `json:"buf_shared_read,omitempty"`
	RowsRemovedByFilter int  `json:"rows_removed_by_filter,omitempty"`
}

var (
	storageRowsRE = regexp.MustCompile(`Storage Rows Scanned:\s*(\d+)`)
	indexRowsRE   = regexp.MustCompile(`Storage Index Rows Scanned:\s*(\d+)`)
	seqScanRE     = regexp.MustCompile(`(?m)^\s*(?:->\s+)?Seq Scan on `)
	indexScanRE   = regexp.MustCompile(`(?m)^\s*(?:->\s+)?Index (?:Only )?Scan using `)
	buffersRE     = regexp.MustCompile(`Buffers: shared hit=(\d+)(?: read=(\d+))?`)
	rowsRemovedRE = regexp.MustCompile(`Rows Removed by Filter:\s*(\d+)`)
)

// ParseAccessMetrics reads the engine-reported access evidence out of a plan
// text. Missing fields stay zero; a cell with none is reported as unparsed.
func ParseAccessMetrics(plan string) AccessMetrics {
	var m AccessMetrics
	if v := storageRowsRE.FindStringSubmatch(plan); v != nil {
		m.StorageRowsScanned, _ = strconv.Atoi(v[1])
	}
	if v := indexRowsRE.FindStringSubmatch(plan); v != nil {
		m.IndexRowsScanned, _ = strconv.Atoi(v[1])
	}
	m.PlanIndexScan = indexScanRE.MatchString(plan)
	m.PlanSeqScan = seqScanRE.MatchString(plan)
	if v := buffersRE.FindStringSubmatch(plan); v != nil {
		m.BufSharedHit, _ = strconv.Atoi(v[1])
		if v[2] != "" {
			m.BufSharedRead, _ = strconv.Atoi(v[2])
		}
	}
	for _, v := range rowsRemovedRE.FindAllStringSubmatch(plan, -1) {
		if n, _ := strconv.Atoi(v[1]); n > m.RowsRemovedByFilter {
			m.RowsRemovedByFilter = n
		}
	}
	return m
}
