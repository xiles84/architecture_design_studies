package evidence

// Plan exhibits for the book's query-plan appendix.
//
// Every SELECT cell in Study 01 stores the text of one EXPLAIN ANALYZE execution,
// captured after the cell's answer passed the correctness gate and before its
// timed samples. The appendix prints a trimmed copy of selected plans beside the
// claim each one illustrates. This file makes that copy mechanical and checkable,
// so the appendix cannot drift from the evidence it quotes:
//
//   - an exhibit may only quote a run that is one of its claim's registry anchors;
//   - the result file's sha256 must equal that anchor's inputs_digest;
//   - the quoted cell must be valid and carry a plan;
//   - every plan-shape statement a reading note makes ("the same node tree as",
//     "a different node tree from") is declared in the spec and recomputed here.
//
// The trim keeps the node tree and the attributes that decide a plan's cost and
// drops the bulk that hides them. The untrimmed text stays in the result file.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PlanSpecSchema and PlanDataSchema version the authored spec and the generated data.
const (
	PlanSpecSchema = "ads-plan-exhibits-spec/1"
	PlanDataSchema = "ads-plan-exhibits/1"
)

// PlanSpec is the authored list of exhibits (book/appendix/query-plans.spec.json).
type PlanSpec struct {
	Schema   string        `json:"schema"`
	Note     string        `json:"note"`
	Exhibits []ExhibitSpec `json:"exhibits"`
}

// ExhibitSpec selects plans from one result file of one anchored run.
type ExhibitSpec struct {
	ID     string       `json:"id"`
	Claim  string       `json:"claim"`
	RunID  string       `json:"run_id"`
	File   string       `json:"file"`
	Plans  []PlanRef    `json:"plans"`
	Shapes []ShapeCheck `json:"shapes,omitempty"`
}

// PlanRef names one cell by its axis.
type PlanRef struct {
	Axis  string `json:"axis"`
	Label string `json:"label"`
}

// ShapeCheck is a plan-shape statement the reading note relies on: the cell at
// Axis has the same (or a different) node tree as the cell at As. AsRun names
// another anchored run of the same claim when the comparison crosses runs, as a
// placement comparison does (one topology per run).
type ShapeCheck struct {
	Axis   string `json:"axis"`
	Label  string `json:"label"`
	As     string `json:"as"`
	AsRun  string `json:"as_run,omitempty"`
	Expect string `json:"expect"`
}

// PlanData is the generated file the Typst appendix renders.
type PlanData struct {
	Schema      string        `json:"schema"`
	GeneratedBy string        `json:"generated_by"`
	Spec        string        `json:"spec"`
	Registry    string        `json:"registry"`
	Trim        TrimRule      `json:"trim"`
	Exhibits    []PlanExhibit `json:"exhibits"`
}

// TrimRule records what the trim keeps, so the appendix can state it from data.
type TrimRule struct {
	Description    string   `json:"description"`
	KeptAttributes []string `json:"kept_attributes"`
}

// PlanExhibit is one exhibit with its provenance resolved and verified.
type PlanExhibit struct {
	ID           string     `json:"id"`
	Claim        string     `json:"claim"`
	RunID        string     `json:"run_id"`
	RunTag       string     `json:"run_tag"`
	RepoCommit   string     `json:"repo_commit"`
	InputsDigest string     `json:"inputs_digest"`
	Topology     string     `json:"topology"`
	ResultsPath  string     `json:"results_path"`
	Plans        []PlanOut  `json:"plans"`
	Shapes       []ShapeOut `json:"shapes"`
}

// PlanOut is one trimmed plan.
type PlanOut struct {
	Axis       string `json:"axis"`
	Label      string `json:"label"`
	Statement  string `json:"statement"`
	Text       string `json:"text"`
	LinesTotal int    `json:"lines_total"`
	LinesKept  int    `json:"lines_kept"`
}

// ShapeOut is a verified plan-shape statement.
type ShapeOut struct {
	Axis      string `json:"axis"`
	Label     string `json:"label"`
	Statement string `json:"statement"`
	As        string `json:"as"`
	AsRun     string `json:"as_run"`
	Expect    string `json:"expect"`
}

// keptAttributes are the non-node plan lines the trim keeps, matched as a prefix
// of the trimmed line. Each one changes how a reader explains the plan's cost:
// which predicate an index serves, how many rows a filter threw away, how often
// a YugabyteDB scan went to storage. Output, Buffers, the per-request RocksDB
// metrics, the planner's cost estimates and memory figures are dropped.
var keptAttributes = []string{
	"Index Cond:", "Recheck Cond:", "Filter:", "Rows Removed by Filter:",
	"Rows Removed by Index Recheck:", "Join Filter:", "Rows Removed by Join Filter:",
	"Hash Cond:", "Merge Cond:", "One-Time Filter:", "Run Condition:",
	"Storage Filter:", "Storage Index Filter:", "Partial Aggregate:",
	"Sort Key:", "Presorted Key:", "Sort Method:", "Group Key:",
	"Cache Key:", "Hits:", "Heap Fetches:", "Workers Launched:",
	"Storage Table Read Requests:", "Storage Index Read Requests:",
	"Storage Table Rows Scanned:", "Storage Index Rows Scanned:",
	"SubPlan", "InitPlan", "CTE ",
	"Planning Time:", "Execution Time:", "Storage Read Requests:", "Storage Rows Scanned:",
}

const trimDescription = "Node lines are kept with their actual time, rows and loops; the planner's cost " +
	"estimate is removed from each node line. Of the other lines only the listed attributes are kept. " +
	"Output column lists, buffer counts, YugabyteDB per-request storage metrics and estimates, " +
	"and memory figures are removed. Indentation and order are unchanged."

var (
	planNodeLine = regexp.MustCompile(`\((cost|actual)[= ]`)
	planCost     = regexp.MustCompile(`  \(cost=[^)]*\)`)
	planDigits   = regexp.MustCompile(`[0-9]+`)
)

// TrimPlan applies the appendix trim. It returns the kept text and the number of
// non-empty lines before and after.
func TrimPlan(plan string) (string, int, int) {
	var kept []string
	total := 0
	for _, line := range strings.Split(strings.ReplaceAll(plan, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		total++
		if planNodeLine.MatchString(line) {
			kept = append(kept, strings.TrimRight(planCost.ReplaceAllString(line, ""), " "))
			continue
		}
		t := strings.TrimSpace(line)
		for _, a := range keptAttributes {
			if strings.HasPrefix(t, a) {
				kept = append(kept, strings.TrimRight(line, " "))
				break
			}
		}
	}
	return strings.Join(kept, "\n"), total, len(kept)
}

// PlanShape reduces a plan to its node tree: depth, node type, relation and index
// names, with numbers removed. It is the same normalisation the SELECT synthesis
// used for "plan shape equal", so the appendix and the synthesis agree on what
// "the same plan" means.
func PlanShape(plan string) string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(plan, "\r\n", "\n"), "\n") {
		if !planNodeLine.MatchString(line) {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		node := strings.TrimPrefix(strings.TrimSpace(line), "->")
		for _, cut := range []string{"  (", " (cost=", " (actual"} {
			if i := strings.Index(node, cut); i >= 0 {
				node = node[:i]
				break
			}
		}
		out = append(out, fmt.Sprintf("%d:%s", indent/2, planDigits.ReplaceAllString(strings.TrimSpace(node), "#")))
	}
	return strings.Join(out, "|")
}

type planCell struct {
	Axis      string `json:"Axis"`
	Statement string `json:"Statement"`
	Status    string `json:"Status"`
	Plan      string `json:"Plan"`
}

type planRun struct {
	Cells []planCell `json:"cells"`
}

// LoadPlanSpec reads the authored exhibit list.
func LoadPlanSpec(path string) (*PlanSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s PlanSpec
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&s); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if s.Schema != PlanSpecSchema {
		return nil, fmt.Errorf("%s: schema %q, want %q", path, s.Schema, PlanSpecSchema)
	}
	return &s, nil
}

// BuildPlanData resolves and verifies every exhibit. specRel and registryRel are
// repo-relative and are recorded in the output; every failure is returned, so one
// run of the generator names every stale exhibit at once.
func BuildPlanData(root, specRel, registryRel string, spec *PlanSpec, reg *V2Document) (*PlanData, []error) {
	out := &PlanData{
		Schema:      PlanDataSchema,
		GeneratedBy: "tools/evidence plans",
		Spec:        filepath.ToSlash(specRel),
		Registry:    filepath.ToSlash(registryRel),
		Trim:        TrimRule{Description: trimDescription, KeptAttributes: keptAttributes},
	}
	var errs []error
	fail := func(id, format string, args ...any) {
		errs = append(errs, fmt.Errorf("exhibit %s: %s", id, fmt.Sprintf(format, args...)))
	}
	claims := map[string]*V2Claim{}
	for _, c := range reg.Claims {
		claims[c.ClaimID] = c
	}
	seen := map[string]bool{}
	for _, ex := range spec.Exhibits {
		if ex.ID == "" || seen[ex.ID] {
			fail(ex.ID, "empty or duplicate id")
			continue
		}
		seen[ex.ID] = true
		c, ok := claims[ex.Claim]
		if !ok {
			fail(ex.ID, "claim %s is not in the registry", ex.Claim)
			continue
		}
		if c.Status != "active" {
			fail(ex.ID, "claim %s has status %q; only active claims may be illustrated", ex.Claim, c.Status)
		}
		if c.Strength == "gap" {
			fail(ex.ID, "claim %s is a gap; a gap has no measured plan to show", ex.Claim)
		}
		// loadRun opens one anchored result file of this claim and verifies its
		// digest; an exhibit quotes nothing that the registry does not cite.
		type loadedRun struct {
			anchor     *V2Anchor
			resultsRel string
			digest     string
			cells      map[string]planCell
		}
		loadRun := func(runID string) (*loadedRun, bool) {
			var anchor *V2Anchor
			for i := range c.Anchors {
				if c.Anchors[i].RunID == runID {
					anchor = &c.Anchors[i]
					break
				}
			}
			if anchor == nil {
				fail(ex.ID, "run %s is not an anchor of %s", runID, ex.Claim)
				return nil, false
			}
			resultsRel := filepath.ToSlash(filepath.Join("studies", c.Study, "results", runID, ex.File))
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(resultsRel)))
			if err != nil {
				fail(ex.ID, "read %s: %v", resultsRel, err)
				return nil, false
			}
			sum := sha256.Sum256(data)
			digest := hex.EncodeToString(sum[:])
			if digest != anchor.InputsDigest {
				fail(ex.ID, "sha256(%s) = %s, but the %s anchor records %s", resultsRel, digest, ex.Claim, anchor.InputsDigest)
				return nil, false
			}
			var r planRun
			if err := json.Unmarshal(data, &r); err != nil {
				fail(ex.ID, "parse %s: %v", resultsRel, err)
				return nil, false
			}
			lr := &loadedRun{anchor: anchor, resultsRel: resultsRel, digest: digest, cells: map[string]planCell{}}
			for _, cl := range r.Cells {
				lr.cells[cl.Axis] = cl
			}
			return lr, true
		}
		get := func(lr *loadedRun, axis string) (planCell, bool) {
			cl, ok := lr.cells[axis]
			switch {
			case !ok:
				fail(ex.ID, "no cell with axis %q in %s", axis, lr.resultsRel)
			case cl.Status != "valid":
				fail(ex.ID, "cell %q has status %q; only cells that passed the correctness gate are shown", axis, cl.Status)
			case strings.TrimSpace(cl.Plan) == "":
				fail(ex.ID, "cell %q has no plan", axis)
			default:
				return cl, true
			}
			return cl, false
		}
		run, ok := loadRun(ex.RunID)
		if !ok {
			continue
		}
		anchor := run.anchor
		pe := PlanExhibit{ID: ex.ID, Claim: ex.Claim, RunID: ex.RunID, InputsDigest: run.digest,
			Topology: anchor.Topology, ResultsPath: run.resultsRel, Shapes: []ShapeOut{}}
		if anchor.RunTag != nil {
			pe.RunTag = *anchor.RunTag
		}
		if anchor.RepoCommit != nil {
			pe.RepoCommit = *anchor.RepoCommit
		}
		if len(ex.Plans) == 0 {
			fail(ex.ID, "no plans selected")
		}
		for _, p := range ex.Plans {
			cl, ok := get(run, p.Axis)
			if !ok {
				continue
			}
			text, total, kept := TrimPlan(cl.Plan)
			pe.Plans = append(pe.Plans, PlanOut{Axis: p.Axis, Label: p.Label, Statement: cl.Statement,
				Text: text, LinesTotal: total, LinesKept: kept})
		}
		for _, s := range ex.Shapes {
			other := run
			if s.AsRun != "" && s.AsRun != ex.RunID {
				if other, ok = loadRun(s.AsRun); !ok {
					continue
				}
			}
			a, okA := get(run, s.Axis)
			b, okB := get(other, s.As)
			if !okA || !okB {
				continue
			}
			same := PlanShape(a.Plan) == PlanShape(b.Plan)
			switch {
			case s.Expect != "same" && s.Expect != "different":
				fail(ex.ID, "shape check %q: expect must be \"same\" or \"different\", got %q", s.Axis, s.Expect)
			case same && s.Expect == "different":
				fail(ex.ID, "shape check: %q has the same node tree as %q, but the spec expects a different one", s.Axis, s.As)
			case !same && s.Expect == "same":
				fail(ex.ID, "shape check: %q has a different node tree from %q, but the spec expects the same one", s.Axis, s.As)
			}
			asRun := s.AsRun
			if asRun == "" {
				asRun = ex.RunID
			}
			pe.Shapes = append(pe.Shapes, ShapeOut{Axis: s.Axis, Label: s.Label, Statement: a.Statement, As: s.As, AsRun: asRun, Expect: s.Expect})
		}
		out.Exhibits = append(out.Exhibits, pe)
	}
	return out, errs
}

// EncodePlanData renders the data deterministically: two-space indent, no HTML
// escaping, LF line endings and a final newline.
func EncodePlanData(d *PlanData) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
