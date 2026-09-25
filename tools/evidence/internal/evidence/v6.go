package evidence

import "fmt"

// V6 method metadata. These types are declared here rather than in v2.go so the
// v2-v5 validators stay readable: the fields are carried on V2Claim (one loader
// reads every versioned file) but only ValidateV6 reads them.

// V2Attribution records whether a claim's headline delta may be named to a
// component, and which isolation pairs support doing so. `Withheld` true means
// the delta must not be attributed; `IsolatedBy` names pairs that do support a
// component-level statement.
type V2Attribution struct {
	Withheld   bool     `json:"withheld"`
	IsolatedBy []string `json:"isolated_by"`
	Note       string   `json:"note"`
}

// V2Prediction is a prediction recorded before measurement.
type V2Prediction struct {
	Text    string `json:"text"`
	Source  string `json:"source"`
	Verdict string `json:"verdict"`
}

// V2Replication is how many whole runs support a claim, in what order, and at
// which confidence tier. A same-machine repeat is `run-level`; only separate
// hosts are `hardware-independent`.
type V2Replication struct {
	Runs               int      `json:"runs"`
	Order              string   `json:"order"`
	AcrossRunSpreadPct *float64 `json:"across_run_spread_pct"`
	Tier               string   `json:"tier"`
}

// V2Crossover records a sweep that locates where one design overtakes another.
type V2Crossover struct {
	Dimension        string  `json:"dimension"`
	Levels           int     `json:"levels"`
	ReversalObserved bool    `json:"reversal_observed"`
	CrossingInterval *string `json:"crossing_interval"`
}

// ValidateV6 resolves the v6 registry. v6 carries every still-active v5 claim
// forward with its statement, limits, anchors, confounds, trials and lifecycle
// unchanged, and adds the empirical-architecture method metadata: evidence_form,
// attribution, control_coverage, and optional prediction/replication/crossover.
//
// On top of the shared resolver, the v4 lifecycle rules and the v5 repetition
// lint, it enforces the method rules the 2026-09-25 charter states:
//
//   - a package comparison may not name a component unless an isolation pair
//     supports it, and a claim listing a confound must either withhold
//     attribution or name that pair;
//   - a control that was present but never fired is not coverage, so a read-only
//     phase can no longer imply its designs were correct;
//   - a claim at strength run_replicated or hardware_replicated must record a
//     prediction, at least two whole runs, and a non-fixed order, and the two
//     tiers may never be conflated;
//   - a crossover that observed a reversal must state at least five levels and
//     the crossing interval, because a two-point comparison cannot bound a
//     crossing.
func ValidateV6(repo string, schema map[string]any, doc *V2Document, raw any, conf *V2Confounds, v1 *Document, v2, v3, v4, v5 *V2Document) []error {
	predecessors := map[string]bool{}
	if v1 != nil {
		for _, c := range v1.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	for _, d := range []*V2Document{v2, v3, v4, v5} {
		if d == nil {
			continue
		}
		for _, c := range d.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	errs := validateRegistry(repo, schema, doc, raw, conf, predecessors, "predecessor")
	errs = append(errs, validateLifecycle(doc)...)
	errs = append(errs, validateRepetition(doc)...)
	errs = append(errs, validateEvidenceForm(doc)...)
	errs = append(errs, validateControlCoverage(doc)...)
	errs = append(errs, validateReplication(doc)...)
	return append(errs, validateCrossover(doc)...)
}

func validateEvidenceForm(doc *V2Document) []error {
	var errs []error
	for _, c := range doc.Claims {
		iso := len(c.Attribution.IsolatedBy)
		switch c.EvidenceForm {
		case "isolate", "interaction", "package", "not_applicable":
		default:
			errs = append(errs, fmt.Errorf("%s: unknown evidence_form %q", c.ClaimID, c.EvidenceForm))
			continue
		}
		if c.EvidenceForm == "package" && !c.Attribution.Withheld && iso == 0 {
			errs = append(errs, fmt.Errorf(
				"%s: evidence_form package neither withholds attribution nor names an isolation pair; a package result may not be attributed to a component",
				c.ClaimID))
		}
		if len(c.Confounds) > 0 && !c.Attribution.Withheld && iso == 0 {
			errs = append(errs, fmt.Errorf(
				"%s: lists confound(s) %v but neither withholds attribution nor names an isolation pair",
				c.ClaimID, c.Confounds))
		}
		if c.EvidenceForm == "not_applicable" && iso > 0 {
			errs = append(errs, fmt.Errorf(
				"%s: evidence_form not_applicable cannot name isolation pairs", c.ClaimID))
		}
	}
	return errs
}

func validateControlCoverage(doc *V2Document) []error {
	var errs []error
	for _, c := range doc.Claims {
		switch c.ControlCoverage {
		case "present", "absent", "not_applicable":
		default:
			errs = append(errs, fmt.Errorf("%s: unknown control_coverage %q", c.ClaimID, c.ControlCoverage))
			continue
		}
		if c.ControlCoverage == "present" && c.RunLevel.ControlsFired == 0 {
			errs = append(errs, fmt.Errorf(
				"%s: control_coverage present but controls_fired is %d; a control that never fired is not coverage",
				c.ClaimID, c.RunLevel.ControlsFired))
		}
		if c.ControlCoverage == "present" && c.RunLevel.ControlsPresent == 0 {
			errs = append(errs, fmt.Errorf(
				"%s: control_coverage present but the run declares no control present", c.ClaimID))
		}
	}
	return errs
}

func validateReplication(doc *V2Document) []error {
	var errs []error
	needsRun := func(s string) bool { return s == "run_replicated" || s == "hardware_replicated" }
	for _, c := range doc.Claims {
		if needsRun(c.Strength) && c.Prediction == nil {
			errs = append(errs, fmt.Errorf(
				"%s: strength %s requires a recorded prediction with a verdict", c.ClaimID, c.Strength))
		}
		r := c.Replication
		if r == nil {
			if needsRun(c.Strength) {
				errs = append(errs, fmt.Errorf(
					"%s: strength %s requires a replication record", c.ClaimID, c.Strength))
			}
			continue
		}
		switch r.Order {
		case "fixed", "abba", "randomized", "interleaved":
		default:
			errs = append(errs, fmt.Errorf("%s: unknown replication order %q", c.ClaimID, r.Order))
		}
		switch r.Tier {
		case "single-run", "run-level", "hardware-independent":
		default:
			errs = append(errs, fmt.Errorf("%s: unknown replication tier %q", c.ClaimID, r.Tier))
		}
		if r.Tier == "run-level" && r.Runs < 2 {
			errs = append(errs, fmt.Errorf("%s: tier run-level needs at least 2 whole runs, has %d", c.ClaimID, r.Runs))
		}
		if r.Tier == "single-run" && r.Runs > 1 {
			errs = append(errs, fmt.Errorf("%s: tier single-run contradicts %d whole runs", c.ClaimID, r.Runs))
		}
		if r.Tier == "hardware-independent" && r.Runs < 2 {
			errs = append(errs, fmt.Errorf("%s: tier hardware-independent needs at least 2 whole runs, has %d", c.ClaimID, r.Runs))
		}
		if c.Strength == "run_replicated" {
			if r.Runs < 2 {
				errs = append(errs, fmt.Errorf("%s: strength run_replicated needs at least 2 whole runs, has %d", c.ClaimID, r.Runs))
			}
			if r.Order == "fixed" {
				errs = append(errs, fmt.Errorf(
					"%s: strength run_replicated needs randomized, interleaved or AB/BA order, not fixed", c.ClaimID))
			}
		}
		if c.Strength == "hardware_replicated" && r.Tier != "hardware-independent" {
			errs = append(errs, fmt.Errorf(
				"%s: strength hardware_replicated requires tier hardware-independent, has %q", c.ClaimID, r.Tier))
		}
	}
	return errs
}

func validateCrossover(doc *V2Document) []error {
	var errs []error
	for _, c := range doc.Claims {
		cr := c.Crossover
		if cr == nil {
			continue
		}
		if cr.ReversalObserved && cr.Levels < 5 {
			errs = append(errs, fmt.Errorf(
				"%s: a crossover with a reversal needs at least 5 levels to bound the crossing, has %d", c.ClaimID, cr.Levels))
		}
		if cr.ReversalObserved && (cr.CrossingInterval == nil || *cr.CrossingInterval == "") {
			errs = append(errs, fmt.Errorf(
				"%s: a crossover with a reversal must state the crossing interval", c.ClaimID))
		}
	}
	return errs
}
