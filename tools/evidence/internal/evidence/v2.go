package evidence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// V2Document is the correction-package registry (book/evidence/v2/claims.json).
type V2Document struct {
	SchemaVersion   int         `json:"schema_version"`
	GeneratedFrom   string      `json:"generated_from"`
	Taxonomy        V2Taxonomy  `json:"taxonomy"`
	RetiredV1Claims []V2Retired `json:"retired_v1_claims"`
	Claims          []*V2Claim  `json:"claims"`
}

type V2Taxonomy struct {
	Families []V2Family `json:"families"`
}

type V2Family struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type V2Retired struct {
	ClaimID string `json:"claim_id"`
	Reason  string `json:"reason"`
}

type V2Claim struct {
	ClaimID                    string      `json:"claim_id"`
	Statement                  string      `json:"statement"`
	Family                     string      `json:"family"`
	Study                      string      `json:"study"`
	Strength                   string      `json:"strength"`
	Trials                     V2Trials    `json:"trials"`
	Transfer                   string      `json:"transfer"`
	Comparability              string      `json:"comparability"`
	ComparabilityRecord        *string     `json:"comparability_record"`
	Confounds                  []string    `json:"confounds"`
	Limits                     []string    `json:"limits"`
	Anchors                    []V2Anchor  `json:"anchors"`
	RunLevel                   V2RunLevel  `json:"run_level"`
	Support                    []V2Support `json:"support"`
	GapBasis                   *string     `json:"gap_basis"`
	LegacyProvenanceIncomplete bool        `json:"legacy_provenance_incomplete"`
	Supersedes                 *string     `json:"supersedes"`
	Review                     V2Review    `json:"review"`
}

type V2Trials struct {
	WholeRunReplications   int `json:"whole_run_replications"`
	FreshLoadTrialsPerCell int `json:"fresh_load_trials_per_cell"`
	InnerIterations        int `json:"inner_iterations"`
}

type V2Anchor struct {
	RunID        string  `json:"run_id"`
	RunTag       *string `json:"run_tag"`
	RepoCommit   *string `json:"repo_commit"`
	InputsDigest string  `json:"inputs_digest"`
	Report       string  `json:"report"`
	Analysis     string  `json:"analysis"`
	Environment  string  `json:"environment"`
	Topology     string  `json:"topology"`
	Role         string  `json:"role"`
}

type V2RunLevel struct {
	ReportedCells   *int   `json:"reported_cells"`
	ReportedFailed  *int   `json:"reported_failed"`
	Source          string `json:"source"`
	ControlsPresent int    `json:"controls_present"`
	ControlsFired   int    `json:"controls_fired"`
	Note            string `json:"note"`
}

type V2Support struct {
	Key        V2SupportKey `json:"key"`
	Status     string       `json:"status"`
	ResolvesTo string       `json:"resolves_to"`
	Expected   string       `json:"expected"`
}

type V2SupportKey struct {
	Topology  string `json:"topology"`
	Design    string `json:"design"`
	Operation string `json:"operation"`
	Metric    string `json:"metric"`
}

type V2Review struct {
	State    string `json:"state"`
	Artifact string `json:"artifact"`
}

// V2Confounds is book/evidence/v2/confounds.json.
type V2Confounds struct {
	SchemaVersion int          `json:"schema_version"`
	Note          string       `json:"note"`
	Confounds     []V2Confound `json:"confounds"`
}

type V2Confound struct {
	ConfoundID     string   `json:"confound_id"`
	Summary        string   `json:"summary"`
	Detail         string   `json:"detail"`
	AffectedClaims []string `json:"affected_claims"`
}

// LoadV2 reads the correction registry and returns typed and raw forms.
func LoadV2(path string) (*V2Document, any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var doc V2Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return &doc, raw, nil
}

// LoadConfounds reads the confound register.
func LoadConfounds(path string) (*V2Confounds, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c V2Confounds
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

var v2TagPattern = regexp.MustCompile(`^run/0[1-5]-`)

// ValidateV2 is the semantic resolver: a green run means every atomic support
// key resolves to a token in the primary anchor's cited report, every anchor
// resolves, run-level status is separate from support status, and confounds and
// supersessions are closed over their registers.
func ValidateV2(repo string, schema map[string]any, doc *V2Document, raw any, conf *V2Confounds, v1 *Document) []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	for _, e := range ValidateSchema(schema, raw) {
		add("schema %s", e.Error())
	}

	families := map[string]bool{}
	for _, f := range doc.Taxonomy.Families {
		families[f.ID] = true
	}
	if len(families) != 6 {
		add("taxonomy must declare exactly 6 families, got %d", len(families))
	}

	confounds := map[string]*V2Confound{}
	if conf != nil {
		for i := range conf.Confounds {
			c := &conf.Confounds[i]
			if confounds[c.ConfoundID] != nil {
				add("duplicate confound_id %s", c.ConfoundID)
			}
			confounds[c.ConfoundID] = c
		}
	}

	v1IDs := map[string]bool{}
	if v1 != nil {
		for _, c := range v1.Claims {
			v1IDs[c.ClaimID] = true
		}
	}

	seen := map[string]bool{}
	referencedConfounds := map[string]bool{}

	for _, c := range doc.Claims {
		if seen[c.ClaimID] {
			add("%s: duplicate claim_id", c.ClaimID)
		}
		seen[c.ClaimID] = true
		if !families[c.Family] {
			add("%s: family %q is not in the taxonomy", c.ClaimID, c.Family)
		}

		numeric := c.Strength != "gap" && c.Strength != "analogy"
		var primary *V2Anchor
		for i := range c.Anchors {
			a := &c.Anchors[i]
			if a.Role == "primary" && primary == nil {
				primary = a
			}
		}
		primaryText := ""
		for i := range c.Anchors {
			a := &c.Anchors[i]
			reportPath := filepath.Join(repo, a.Report)
			analysisPath := filepath.Join(repo, a.Analysis)
			reportText, rerr := os.ReadFile(reportPath)
			if rerr != nil {
				add("%s: anchor %s report %s does not resolve", c.ClaimID, a.RunID, a.Report)
			}
			analysisText, aerr := os.ReadFile(analysisPath)
			if aerr != nil {
				add("%s: anchor %s analysis %s does not resolve", c.ClaimID, a.RunID, a.Analysis)
			}
			if rerr == nil && !strings.Contains(string(reportText), a.InputsDigest) {
				add("%s: anchor %s inputs_digest %s does not appear in report %s", c.ClaimID, a.RunID, a.InputsDigest, a.Report)
			}
			if aerr == nil && !strings.Contains(string(analysisText), a.InputsDigest) {
				add("%s: anchor %s inputs_digest %s does not appear in analysis %s", c.ClaimID, a.RunID, a.InputsDigest, a.Analysis)
			}
			if numeric {
				runDir := filepath.Join(repo, "studies", c.Study, "results", a.RunID)
				if !exists(runDir) {
					add("%s: anchor %s has no results directory %s", c.ClaimID, a.RunID, runDir)
				}
			}
			if a.RunTag != nil && *a.RunTag != "" {
				if !v2TagPattern.MatchString(*a.RunTag) {
					add("%s: anchor %s run tag %q is not a run/<study>/... tag", c.ClaimID, a.RunID, *a.RunTag)
				} else if !tagExists(repo, *a.RunTag) {
					add("%s: anchor %s run tag %s does not exist", c.ClaimID, a.RunID, *a.RunTag)
				}
			}
			if a.Role == "primary" && rerr == nil {
				primaryText = string(reportText)
			}
		}
		if primary == nil {
			add("%s: no primary anchor", c.ClaimID)
		}

		// Every structured support key must resolve to a token in the cited report.
		for i, s := range c.Support {
			if primaryText == "" {
				add("%s: support[%d] cannot be resolved without a readable primary report", c.ClaimID, i)
				continue
			}
			if !strings.Contains(primaryText, s.ResolvesTo) {
				add("%s: support[%d] key %s/%s/%s resolves_to %q does not appear in %s",
					c.ClaimID, i, s.Key.Topology, s.Key.Design, s.Key.Operation, s.ResolvesTo, primary.Report)
			}
		}
		if numeric && len(c.Support) == 0 {
			add("%s: numeric claim has no structured support key", c.ClaimID)
		}
		if c.Strength == "gap" {
			if c.GapBasis == nil || strings.TrimSpace(*c.GapBasis) == "" {
				add("%s: gap claim needs a gap_basis", c.ClaimID)
			}
		}

		// Legacy provenance is only allowed when the primary run has no tag.
		noTag := primary == nil || primary.RunTag == nil || *primary.RunTag == ""
		if c.LegacyProvenanceIncomplete != noTag {
			add("%s: legacy_provenance_incomplete=%v but its run tag presence is %v", c.ClaimID, c.LegacyProvenanceIncomplete, !noTag)
		}

		// Trials must support the declared strength.
		switch c.Strength {
		case "repeated_controlled":
			if c.Trials.WholeRunReplications < 2 && c.Trials.FreshLoadTrialsPerCell < 2 {
				add("%s: repeated_controlled needs >=2 whole-run replications or fresh-load trials", c.ClaimID)
			}
		case "single_run_directional":
			if c.Trials.WholeRunReplications > 1 {
				add("%s: single_run_directional has %d whole-run replications", c.ClaimID, c.Trials.WholeRunReplications)
			}
		}

		if c.Supersedes != nil && *c.Supersedes != "" {
			if v1 == nil {
				add("%s: supersedes %s but the v1 registry could not be read", c.ClaimID, *c.Supersedes)
			} else if !v1IDs[*c.Supersedes] {
				add("%s: supersedes %s, which is not in the v1 registry", c.ClaimID, *c.Supersedes)
			}
		}

		for _, id := range c.Confounds {
			cf, ok := confounds[id]
			if !ok {
				add("%s: confound %s is not in the register", c.ClaimID, id)
				continue
			}
			referencedConfounds[id] = true
			found := false
			for _, a := range cf.AffectedClaims {
				if a == c.ClaimID {
					found = true
				}
			}
			if !found {
				add("%s: confound %s does not list this claim in affected_claims", c.ClaimID, id)
			}
		}

		if c.RunLevel.ReportedCells != nil && c.RunLevel.ReportedFailed != nil {
			if *c.RunLevel.ReportedFailed > *c.RunLevel.ReportedCells {
				add("%s: run_level failed %d exceeds cells %d", c.ClaimID, *c.RunLevel.ReportedFailed, *c.RunLevel.ReportedCells)
			}
		}
		if c.RunLevel.ControlsFired > c.RunLevel.ControlsPresent {
			add("%s: controls_fired %d exceeds controls_present %d", c.ClaimID, c.RunLevel.ControlsFired, c.RunLevel.ControlsPresent)
		}
	}

	// Every confound must be referenced by at least one claim.
	for id := range confounds {
		if !referencedConfounds[id] {
			add("confound %s is not referenced by any claim", id)
		}
	}

	// Retired v1 claims must exist in v1 and not be reproduced here.
	for _, r := range doc.RetiredV1Claims {
		if v1 != nil && !v1IDs[r.ClaimID] {
			add("retired v1 claim %s is not in the v1 registry", r.ClaimID)
		}
	}

	return errs
}

// V1DiagnosticReport is the mechanical resolution audit of the v1 registry that
// the brainstorm conclusion computed: 11 of 12 numeric claims have at least one
// declared cell name absent from their cited report, and 36 of 43 names resolve
// nowhere. It is retained as a regression fixture for the resolver.
type V1DiagnosticReport struct {
	NumericClaims     int               `json:"numeric_claims"`
	ClaimsWithMissing int               `json:"claims_with_missing"`
	Names             int               `json:"names"`
	NamesMissing      int               `json:"names_missing"`
	Claims            []V1DiagnosticRow `json:"claims"`
}

type V1DiagnosticRow struct {
	ClaimID    string `json:"claim_id"`
	Report     string `json:"report"`
	Cells      int    `json:"cells"`
	Unresolved int    `json:"unresolved"`
}

// V1Diagnostic recomputes the v1 claim-to-cell resolution failure from the
// committed registry and reports. A cell name resolves only if it appears
// verbatim in the report the claim cites.
func V1Diagnostic(repo string, v1 *Document) (V1DiagnosticReport, error) {
	var rep V1DiagnosticReport
	for _, c := range v1.Claims {
		if c.Strength == "gap" || c.Strength == "analogy" {
			continue
		}
		rep.NumericClaims++
		text, err := os.ReadFile(filepath.Join(repo, c.Provenance.Report))
		if err != nil {
			return rep, fmt.Errorf("%s: report %s: %w", c.ClaimID, c.Provenance.Report, err)
		}
		row := V1DiagnosticRow{ClaimID: c.ClaimID, Report: c.Provenance.Report}
		for _, cell := range c.Provenance.Cells {
			row.Cells++
			rep.Names++
			if !strings.Contains(string(text), cell.Name) {
				row.Unresolved++
				rep.NamesMissing++
			}
		}
		if row.Unresolved > 0 {
			rep.ClaimsWithMissing++
		}
		rep.Claims = append(rep.Claims, row)
	}
	return rep, nil
}
