package selectbench

import (
	"fmt"
	"sort"
	"strings"
)

// Status is the disposition of one measured cell.
type Status string

const (
	StatusValid       Status = "valid"
	StatusInvalid     Status = "invalid"
	StatusUnsupported Status = "unsupported"
	StatusTimeout     Status = "timeout"
	StatusError       Status = "error"
)

// Provenance travels with every result: a number whose environment, image, seed
// and settings cannot be named cannot be compared with anything.
type Provenance struct {
	RepoCommit          string
	ImageDigest         string
	EngineImage         string
	Seed                int64
	DatasetFingerprint  string
	Settings            map[string]string
	RunTag              string
	Environment         string
	EngineImageDigest   string
	EngineVersion       string
	ClientCPUStatBefore string
	ClientCPUStatAfter  string
	ClientCPUQuota      string
	ClientMemoryLimit   string
}

// Cell is one (statement, topology) measurement, with the plan that produced it.
type Cell struct {
	RunID           string
	Topology        string
	Statement       string
	Group           string
	Axis            string
	Params          string
	Status          Status
	MedianMS        float64
	SpreadPct       float64
	Trials          int
	Plan            string
	Note            string
	TimingMeasured  bool
	SamplesMS       []float64
	ResultRows      int
	ResultJSONBytes int
	PlanningMS      *float64
	ExecutionMS     *float64
}

func (c Cell) structureSig() string { return StructureSignature(c.Plan) }

// ReportInput is everything a generated comparison report may contain.
// Interpretation is deliberately absent: a generated report holds measurements,
// and the recommendation belongs in a signed analysis beside it.
type ReportInput struct {
	Provenance     Provenance
	Predictions    []string
	Cells          []Cell
	Interpretation string // must be empty; present only so a caller that tries is refused
	Controls       []ControlResult
	InputsDigest   string
}

func comparableCells(a, b Cell) (bool, string) {
	if a.Group == ControlGroup || b.Group == ControlGroup {
		return false, "a deliberately wrong control may never be pooled with a template"
	}
	if a.Group != b.Group {
		return false, "different semantic groups answer different questions"
	}
	if a.Topology != b.Topology {
		return false, "different topology"
	}
	if a.Statement == b.Statement {
		return false, "same statement"
	}
	return true, ""
}

// RenderComparison writes the measurements-only report: every cell with its
// status, the invalid/unsupported/timed-out cells named and excluded from any
// ranking, and each comparable pair classified as same-plan, no-resolved-
// difference, or a resolved difference.
func RenderComparison(in ReportInput) (string, error) {
	if strings.TrimSpace(in.Interpretation) != "" {
		return "", fmt.Errorf("a generated report holds measurements only; write the interpretation in a signed analysis")
	}
	var b strings.Builder

	b.WriteString("# SELECT comparison (generated measurements)\n\n")
	b.WriteString("_This file contains measurements only. Interpretation belongs in a signed analysis._\n\n")
	b.WriteString("**TL;DR:** a validation run supplies correctness and plan evidence. Within-run samples do not establish a replicated performance winner. Missing measurements are N/A.\n\n")

	b.WriteString("## Provenance\n\n| Field | Value |\n|---|---|\n")
	p := in.Provenance
	fmt.Fprintf(&b, "| repo commit | `%s` |\n", orUnknown(p.RepoCommit))
	fmt.Fprintf(&b, "| inputs digest (SHA-256 of result JSON bytes) | `%s` |\n", orUnknown(in.InputsDigest))
	fmt.Fprintf(&b, "| benchmark image | `%s` |\n", orUnknown(p.EngineImage))
	fmt.Fprintf(&b, "| image digest | `%s` |\n", orUnknown(p.ImageDigest))
	fmt.Fprintf(&b, "| run tag | `%s` |\n| environment | `%s` |\n| engine image digest | `%s` |\n| engine version | %s |\n", orUnknown(p.RunTag), orUnknown(p.Environment), orUnknown(p.EngineImageDigest), orUnknown(p.EngineVersion))
	fmt.Fprintf(&b, "| dataset seed | %d |\n", p.Seed)
	fmt.Fprintf(&b, "| dataset fingerprint | `%s` |\n", orUnknown(p.DatasetFingerprint))
	for _, k := range sortedKeys(p.Settings) {
		fmt.Fprintf(&b, "| setting: %s | %s |\n", k, p.Settings[k])
	}
	b.WriteString("\n")

	// Every cell, as reported.
	b.WriteString("## Cells\n\n| run | topology | statement | group | status | median ms | spread % | trials | plan structure |\n")
	b.WriteString("|---|---|---|---|---|---|---|---|---|\n")
	cells := append([]Cell(nil), in.Cells...)
	sort.SliceStable(cells, func(i, j int) bool {
		if cells[i].Group != cells[j].Group {
			return cells[i].Group < cells[j].Group
		}
		return cells[i].Statement < cells[j].Statement
	})
	for _, c := range cells {
		ms, spread := "N/A", "N/A"
		if c.TimingMeasured && c.Status == StatusValid {
			ms, spread = fmt.Sprintf("%.6f", c.MedianMS), fmt.Sprintf("%.1f", c.SpreadPct)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %d | `%s` |\n",
			c.RunID, c.Topology, c.Statement, c.Group, c.Status, ms, spread, c.Trials, c.structureSig())
	}
	b.WriteString("\n")

	// Excluded cells, named.
	var excluded []Cell
	for _, c := range cells {
		if c.Status != StatusValid {
			excluded = append(excluded, c)
		}
	}
	b.WriteString("## Excluded cells\n\n")
	if len(excluded) == 0 {
		b.WriteString("None: every cell is valid.\n\n")
	} else {
		for _, c := range excluded {
			fmt.Fprintf(&b, "- `%s` on %s — **%s**%s\n", c.Statement, c.Topology, c.Status, dash(c.Note))
		}
		b.WriteString("\nA cell that is invalid, unsupported or timed out is never ranked below.\n\n")
	}

	// Comparable pairs.
	b.WriteString("## Comparisons\n\n")
	valid := make([]Cell, 0, len(cells))
	for _, c := range cells {
		if c.Status == StatusValid {
			valid = append(valid, c)
		}
	}
	type pair struct{ a, b Cell }
	var samePlan, noDiff []string
	var resolved []string
	var incomparable []string
	for i := 0; i < len(valid); i++ {
		for j := i + 1; j < len(valid); j++ {
			a, c := valid[i], valid[j]
			ok, why := comparableCells(a, c)
			if !ok {
				incomparable = append(incomparable, fmt.Sprintf("- `%s` vs `%s`: not comparable — %s", a.Statement, c.Statement, why))
				continue
			}
			label := fmt.Sprintf("`%s` vs `%s` on %s", a.Statement, c.Statement, a.Topology)
			if a.structureSig() == c.structureSig() {
				samePlan = append(samePlan, label+": same plan structure")
				continue
			}
			diff := c.MedianMS - a.MedianMS
			tol := (a.SpreadPct + c.SpreadPct) / 100 * a.MedianMS
			if abs(diff) <= tol {
				noDiff = append(noDiff, fmt.Sprintf("%s: medians differ by %.3f ms, inside the %.1f%%/%.1f%% spread — no resolved difference",
					label, abs(diff), a.SpreadPct, c.SpreadPct))
			} else {
				resolved = append(resolved, fmt.Sprintf("%s: within-run medians %.6f / %.6f ms; whole-run replication required — no resolved difference", label, a.MedianMS, c.MedianMS))
			}
		}
	}
	writeSection(&b, "Same plan (no plan difference to explain)", samePlan)
	writeSection(&b, "No resolved difference", noDiff)
	writeSection(&b, "Within-run differences (not replicated)", resolved)
	b.WriteString("## Live negative controls\n\n")
	for _, control := range in.Controls {
		fmt.Fprintf(&b, "- `%s`: detected=%t, rows=%d%s\n", control.Statement, control.Detected, control.Rows, dash(control.Error))
	}
	b.WriteString("\n## Captured plans and result sizes\n\n")
	for _, c := range cells {
		fmt.Fprintf(&b, "### %s\n\nRows: %d; serialized decoded-row JSON bytes: %d (not wire bytes). Planning/execution: %s / %s ms.\n\n```text\n%s```\n\n", c.Statement, c.ResultRows, c.ResultJSONBytes, optionalMS(c.PlanningMS), optionalMS(c.ExecutionMS), c.Plan+"\n")
	}
	writeSection(&b, "Not comparable (never ranked)", incomparable)

	if len(in.Predictions) > 0 {
		b.WriteString("## Predictions this run informs\n\n")
		for _, p := range in.Predictions {
			fmt.Fprintf(&b, "- %s\n", p)
		}
		b.WriteString("\nVerdicts are assigned in the signed analysis, not here.\n")
	}
	return b.String(), nil
}

func writeSection(b *strings.Builder, title string, lines []string) {
	fmt.Fprintf(b, "### %s\n\n", title)
	if len(lines) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	for _, l := range lines {
		fmt.Fprintf(b, "- %s\n", l)
	}
	b.WriteString("\n")
}

func orUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "unknown"
	}
	return s
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	return " — " + s
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
