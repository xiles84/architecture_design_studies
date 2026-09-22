package evidence

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Document is the whole registry file.
type Document struct {
	SchemaVersion int      `json:"schema_version"`
	GeneratedFrom string   `json:"generated_from"`
	Claims        []*Claim `json:"claims"`
}

// Claim is one book-level statement and its evidence.
type Claim struct {
	ClaimID             string     `json:"claim_id"`
	Statement           string     `json:"statement"`
	Study               string     `json:"study"`
	Strength            string     `json:"strength"`
	Trials              int        `json:"trials"`
	Transfer            string     `json:"transfer"`
	Comparability       string     `json:"comparability"`
	ComparabilityRecord *string    `json:"comparability_record"`
	Supersedes          *string    `json:"supersedes"`
	SupersededBy        *string    `json:"superseded_by"`
	Limits              []string   `json:"limits"`
	Provenance          Provenance `json:"provenance"`
}

// Provenance names the exact run, digest, report and analysis behind a claim.
type Provenance struct {
	RunID        string   `json:"run_id"`
	RunTag       *string  `json:"run_tag"`
	Environment  string   `json:"environment"`
	Topology     string   `json:"topology"`
	RepoCommit   *string  `json:"repo_commit"`
	InputsDigest string   `json:"inputs_digest"`
	Report       string   `json:"report"`
	Analysis     string   `json:"analysis"`
	FailedCells  int      `json:"failed_cells"`
	Cells        []Cell   `json:"cells"`
	WinnerCells  []string `json:"winner_cells"`
}

// Cell is one measured cell and whether it is usable evidence.
type Cell struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

var numericStrengths = map[string]bool{
	"repeated_controlled":    true,
	"single_run_directional": true,
	"mechanism_supported":    true,
}

// Load reads the registry and returns both typed and raw forms.
func Load(path string) (*Document, any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return &doc, raw, nil
}

// Validate enforces the schema and every rule the book relies on, resolving
// paths, digests and tags inside repo.
func Validate(repo string, schema map[string]any, doc *Document, raw any) []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	for _, e := range ValidateSchema(schema, raw) {
		add("schema %s", e.Error())
	}

	byID := map[string]*Claim{}
	for _, c := range doc.Claims {
		if _, dup := byID[c.ClaimID]; dup {
			add("%s: duplicate claim_id", c.ClaimID)
		}
		byID[c.ClaimID] = c
	}

	for _, c := range doc.Claims {
		p := c.Provenance
		report := filepath.Join(repo, p.Report)
		analysis := filepath.Join(repo, p.Analysis)
		reportText, rerr := os.ReadFile(report)
		if rerr != nil {
			add("%s: report %s does not resolve", c.ClaimID, p.Report)
		}
		analysisText, aerr := os.ReadFile(analysis)
		if aerr != nil {
			add("%s: analysis %s does not resolve", c.ClaimID, p.Analysis)
		}
		if env := filepath.Join(repo, "docs/environments", p.Environment+".md"); !exists(env) {
			add("%s: environment page %s.md does not resolve", c.ClaimID, p.Environment)
		}
		if numericStrengths[c.Strength] {
			runDir := filepath.Join(repo, "studies", c.Study, "results", p.RunID)
			if !exists(runDir) {
				add("%s: run %s has no results directory %s", c.ClaimID, p.RunID, runDir)
			}
		}
		// The digest must be exactly the one recorded in both the report and
		// the signed analysis; a re-run changes it, which is how stale evidence
		// is detected.
		if rerr == nil && !strings.Contains(string(reportText), p.InputsDigest) {
			add("%s: inputs_digest %s does not appear in report %s (stale digest?)", c.ClaimID, p.InputsDigest, p.Report)
		}
		if aerr == nil && !strings.Contains(string(analysisText), p.InputsDigest) {
			add("%s: inputs_digest %s does not appear in analysis %s (stale digest?)", c.ClaimID, p.InputsDigest, p.Analysis)
		}
		if p.RunTag != nil && *p.RunTag != "" && !tagExists(repo, *p.RunTag) {
			add("%s: run tag %s does not exist", c.ClaimID, *p.RunTag)
		}

		switch c.Strength {
		case "repeated_controlled":
			if c.Trials < 2 {
				add("%s: repeated_controlled needs trials >= 2, got %d", c.ClaimID, c.Trials)
			}
		case "single_run_directional":
			if c.Trials != 1 {
				add("%s: single_run_directional needs exactly 1 trial, got %d", c.ClaimID, c.Trials)
			}
		case "mechanism_supported":
			if c.Trials < 1 {
				add("%s: mechanism_supported needs at least 1 trial", c.ClaimID)
			}
		case "gap", "analogy":
			if c.Trials != 0 {
				add("%s: %s claims carry no trials, got %d", c.ClaimID, c.Strength, c.Trials)
			}
		}
		if c.Strength == "analogy" && c.Transfer != "analogy" {
			add("%s: analogy strength requires transfer=analogy, got %q", c.ClaimID, c.Transfer)
		}
		if numericStrengths[c.Strength] && c.Transfer == "analogy" {
			add("%s: a numeric claim cannot be labelled an analogy", c.ClaimID)
		}
		if c.Comparability == "cross-study" {
			if c.ComparabilityRecord == nil || *c.ComparabilityRecord == "" {
				add("%s: cross-study numeric comparison needs an explicit comparability record", c.ClaimID)
			} else if !exists(filepath.Join(repo, *c.ComparabilityRecord)) {
				add("%s: comparability record %s does not resolve", c.ClaimID, *c.ComparabilityRecord)
			}
		}

		// A numeric claim's winners must be cells that passed every gate.
		if numericStrengths[c.Strength] {
			if len(p.Cells) == 0 {
				add("%s: numeric claim has no cells", c.ClaimID)
			}
			if len(p.WinnerCells) == 0 {
				add("%s: numeric claim names no winner cell", c.ClaimID)
			}
			status := map[string]string{}
			for _, cell := range p.Cells {
				status[cell.Name] = cell.Status
			}
			bad := 0
			for _, cell := range p.Cells {
				if cell.Status == "failed" || cell.Status == "invalid" {
					bad++
				}
			}
			if bad != p.FailedCells {
				add("%s: failed_cells=%d but %d cells are failed/invalid", c.ClaimID, p.FailedCells, bad)
			}
			for _, w := range p.WinnerCells {
				st, ok := status[w]
				if !ok {
					add("%s: winner cell %q is not listed in cells", c.ClaimID, w)
					continue
				}
				if st != "valid" {
					add("%s: winner cell %q has status %q; a failed or invalid cell cannot be a winner", c.ClaimID, w, st)
				}
			}
		}

		if c.Supersedes != nil && *c.Supersedes != "" {
			target, ok := byID[*c.Supersedes]
			if !ok {
				add("%s: supersedes %s, which is not in the registry", c.ClaimID, *c.Supersedes)
			} else if target.SupersededBy == nil || *target.SupersededBy != c.ClaimID {
				add("%s: supersedes %s but that claim is not marked superseded_by %s", c.ClaimID, *c.Supersedes, c.ClaimID)
			}
		}
		// If this analysis itself declares it supersedes another analysis that
		// is present as a claim, that claim must be marked.
		if aerr == nil {
			if sup := frontMatter(string(analysisText), "supersedes"); sup != "" {
				for _, other := range doc.Claims {
					if other.ClaimID == c.ClaimID {
						continue
					}
					if strings.Contains(other.Provenance.Analysis, sup) &&
						(other.SupersededBy == nil || *other.SupersededBy != c.ClaimID) {
						add("%s: supersedes analysis %s, but claim %s is not marked superseded_by %s", c.ClaimID, sup, other.ClaimID, c.ClaimID)
					}
				}
			}
		}
	}
	return errs
}

// BookInputs returns the claim ids that may appear in the book: every claim
// that has not been superseded.
func (d *Document) BookInputs() []string {
	var out []string
	for _, c := range d.Claims {
		if c.SupersededBy == nil || *c.SupersededBy == "" {
			out = append(out, c.ClaimID)
		}
	}
	return out
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func tagExists(repo, tag string) bool {
	cmd := exec.Command("git", "-C", repo, "show-ref", "--verify", "--quiet", "refs/tags/"+tag)
	return cmd.Run() == nil
}

// frontMatter returns a scalar value from a leading `---` YAML block.
func frontMatter(text, key string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			break
		}
		if strings.HasPrefix(line, key+":") {
			v := strings.TrimSpace(strings.TrimPrefix(line, key+":"))
			return strings.Trim(v, `"'`)
		}
	}
	return ""
}
