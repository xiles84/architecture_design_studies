//go:build ignore

// register expands claims-spec.json with the anchors that evidence.go -anchors
// prints and inserts the result into book/evidence/v6/claims.json, after the
// last existing claim. It writes the claims in the registry's own key order and
// refuses to run twice (an id already present aborts). It is a one-off editing
// aid for the SELECT synthesis, kept beside the spec as provenance; the
// registry's own validator (tools/evidence) is the gate.
//
//	podman run --rm --network none -v "<worktree>:/repo" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/select-query-templates/synthesis-20261005/register.go <anchors.json>
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	specPath   = "docs/select-query-templates/synthesis-20261005/claims-spec.json"
	claimsPath = "book/evidence/v6/claims.json"
	reviewRef  = "studies/01-charity-tree/reports/analyses/20261005T-select-synthesis--claude-opus-5-5--2026-10-05.md"
)

type Anchor struct {
	RunID        string `json:"run_id"`
	RunTag       string `json:"run_tag"`
	RepoCommit   string `json:"repo_commit"`
	InputsDigest string `json:"inputs_digest"`
	Report       string `json:"report"`
	Analysis     string `json:"analysis"`
	Environment  string `json:"environment"`
	Topology     string `json:"topology"`
	Role         string `json:"role"`
}

type Trials struct {
	WholeRunReplications   int `json:"whole_run_replications"`
	FreshLoadTrialsPerCell int `json:"fresh_load_trials_per_cell"`
	InnerIterations        int `json:"inner_iterations"`
}

type RunLevel struct {
	ReportedCells   *int   `json:"reported_cells"`
	ReportedFailed  *int   `json:"reported_failed"`
	Source          string `json:"source"`
	ControlsPresent int    `json:"controls_present"`
	ControlsFired   int    `json:"controls_fired"`
	Note            string `json:"note"`
}

type Support struct {
	Key struct {
		Topology  string `json:"topology"`
		Design    string `json:"design"`
		Operation string `json:"operation"`
		Metric    string `json:"metric"`
	} `json:"key"`
	Status     string `json:"status"`
	ResolvesTo string `json:"resolves_to"`
	Expected   string `json:"expected"`
}

type Review struct {
	State    string `json:"state"`
	Artifact string `json:"artifact"`
}

type Attribution struct {
	Withheld   bool     `json:"withheld"`
	IsolatedBy []string `json:"isolated_by"`
	Note       string   `json:"note"`
}

type Replication struct {
	Runs               int      `json:"runs"`
	Order              string   `json:"order"`
	AcrossRunSpreadPct *float64 `json:"across_run_spread_pct"`
	Tier               string   `json:"tier"`
}

type Crossover struct {
	Dimension        string  `json:"dimension"`
	Levels           int     `json:"levels"`
	ReversalObserved bool    `json:"reversal_observed"`
	CrossingInterval *string `json:"crossing_interval"`
}

type Prediction struct {
	Text    string `json:"text"`
	Source  string `json:"source"`
	Verdict string `json:"verdict"`
}

// Claim mirrors the v6 schema in the registry's key order.
type Claim struct {
	ClaimID                    string       `json:"claim_id"`
	Statement                  string       `json:"statement"`
	Family                     string       `json:"family"`
	Study                      string       `json:"study"`
	Strength                   string       `json:"strength"`
	Trials                     Trials       `json:"trials"`
	Transfer                   string       `json:"transfer"`
	Comparability              string       `json:"comparability"`
	ComparabilityRecord        *string      `json:"comparability_record"`
	Confounds                  []string     `json:"confounds"`
	Limits                     []string     `json:"limits"`
	Anchors                    []Anchor     `json:"anchors"`
	RunLevel                   RunLevel     `json:"run_level"`
	Support                    []Support    `json:"support"`
	GapBasis                   *string      `json:"gap_basis,omitempty"`
	LegacyProvenanceIncomplete bool         `json:"legacy_provenance_incomplete"`
	Supersedes                 *string      `json:"supersedes"`
	Review                     Review       `json:"review"`
	Status                     string       `json:"status"`
	GapKind                    []string     `json:"gap_kind,omitempty"`
	EvidenceForm               string       `json:"evidence_form"`
	Attribution                Attribution  `json:"attribution"`
	ControlCoverage            string       `json:"control_coverage"`
	Replication                *Replication `json:"replication,omitempty"`
	Crossover                  *Crossover   `json:"crossover,omitempty"`
	Prediction                 *Prediction  `json:"prediction,omitempty"`
}

type spec struct {
	Claim
	AnchorSets []string `json:"anchor_sets"`
	PrimaryRun string   `json:"primary_run"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: register <anchors.json>")
		os.Exit(2)
	}
	var anchors map[string][]Anchor
	mustRead(os.Args[1], &anchors)
	var specs []spec
	mustRead(specPath, &specs)

	text, err := os.ReadFile(claimsPath)
	if err != nil {
		panic(err)
	}
	var entries []string
	for _, s := range specs {
		if bytes.Contains(text, []byte(`"claim_id": "`+s.ClaimID+`"`)) {
			fail("%s is already in %s", s.ClaimID, claimsPath)
		}
		c := s.Claim
		c.Family, c.Study = "normalized-facts", "01-charity-tree"
		c.Transfer, c.Comparability = "direct", "within-run"
		c.Status = "active"
		c.Review = Review{State: "open", Artifact: reviewRef}
		if c.Strength == "gap" {
			c.Transfer, c.Comparability = "not-applicable", "not-applicable"
		}
		if c.Confounds == nil {
			c.Confounds = []string{}
		}
		if c.Support == nil {
			c.Support = []Support{}
		}
		primary := false
		for _, set := range s.AnchorSets {
			as, ok := anchors[set]
			if !ok || len(as) == 0 {
				fail("%s: unknown or empty anchor set %s", s.ClaimID, set)
			}
			for _, a := range as {
				a.Role = "corroborating"
				if a.RunID == s.PrimaryRun {
					a.Role = "primary"
					primary = true
				}
				c.Anchors = append(c.Anchors, a)
			}
		}
		if !primary {
			fail("%s: primary run %s is not in its anchor sets", s.ClaimID, s.PrimaryRun)
		}
		// The primary anchor goes first, as in the existing claims.
		for i, a := range c.Anchors {
			if a.Role == "primary" && i > 0 {
				c.Anchors = append([]Anchor{a}, append(c.Anchors[:i:i], c.Anchors[i+1:]...)...)
				break
			}
		}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("    ", "  ")
		if err := enc.Encode(c); err != nil {
			panic(err)
		}
		entries = append(entries, "    "+strings.TrimRight(buf.String(), "\n"))
	}

	marker := "\n  ],\n  \"retired_predecessor_claims\""
	idx := bytes.Index(text, []byte(marker))
	if idx < 0 {
		fail("could not find the end of the claims array in %s", claimsPath)
	}
	out := string(text[:idx]) + ",\n" + strings.Join(entries, ",\n") + string(text[idx:])
	if err := os.WriteFile(claimsPath, []byte(out), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("inserted %d claims into %s\n", len(entries), claimsPath)
}

func mustRead(path string, v any) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		fail("%s: %v", path, err)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
