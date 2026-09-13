// Package markdown holds the formatting rules shared by every generated report.
// The rules are about honesty as much as looks: precision that matches the
// noise, and a dash rather than a zero where nothing was measured.
package markdown

import (
	"fmt"
	"strings"
)

// Ops renders throughput at a precision that matches how much it varies: four
// significant figures on a number like 48 090 is false precision.
func Ops(v float64) string {
	switch {
	case v <= 0:
		return "—"
	case v >= 10000:
		return fmt.Sprintf("%.0fk", v/1000)
	case v >= 1000:
		return fmt.Sprintf("%.1fk", v/1000)
	case v >= 100:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%.1f", v)
	}
}

// MS renders a latency; zero means "not measured / not enough samples".
func MS(v float64) string {
	switch {
	case v <= 0:
		return "—"
	case v < 10:
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%.0f", v)
	}
}

func Bytes(b int64) string {
	if b <= 0 {
		return "—"
	}
	const u = 1024
	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGT"[exp])
}

// Ratio expresses b relative to a as "Nx faster/slower": the effects these
// studies measure span orders of magnitude, where percentages stop reading.
func Ratio(a, b float64) string {
	if a <= 0 || b <= 0 {
		return "—"
	}
	if b >= a {
		return fmt.Sprintf("%.1fx faster", b/a)
	}
	return fmt.Sprintf("%.1fx slower", a/b)
}

// Table accumulates a markdown table.
type Table struct {
	head  []string
	align []string
	rows  [][]string
}

// NewTable takes column headers; a header ending in ">" is right-aligned (the
// marker is stripped), which is what every numeric column wants.
func NewTable(headers ...string) *Table {
	t := &Table{}
	for _, h := range headers {
		if strings.HasSuffix(h, ">") {
			t.head = append(t.head, strings.TrimSuffix(h, ">"))
			t.align = append(t.align, "---:")
		} else {
			t.head = append(t.head, h)
			t.align = append(t.align, "---")
		}
	}
	return t
}

func (t *Table) Row(cells ...string) { t.rows = append(t.rows, cells) }

func (t *Table) Len() int { return len(t.rows) }

func (t *Table) Write(b *strings.Builder) {
	fmt.Fprintf(b, "| %s |\n|%s|\n", strings.Join(t.head, " | "), strings.Join(t.align, "|"))
	for _, r := range t.rows {
		fmt.Fprintf(b, "| %s |\n", strings.Join(r, " | "))
	}
	b.WriteString("\n")
}
