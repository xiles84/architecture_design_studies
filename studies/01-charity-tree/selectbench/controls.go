package selectbench

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
)

type ControlResult struct {
	Statement string
	Group     string
	Detected  bool
	Rows      int
	Observed  [][]any
	Error     string
}

var controlContracts = map[string]string{
	"x01_optional_where":         "sel_optional_match",
	"x02_not_in_nullable":        "sel_absence",
	"x03_distinct_on_name":       "sel_membership",
	"x04_limit_before_join":      "sel_full_rows",
	"x05_aggregate_double_count": "sel_aggregate_sel",
	"x06_topn_rank_ties":         "sel_topn_sel",
}

// DetectControlRows requires a successful query with the same column shape as
// its intended contract. Query/decoding errors are not successful controls.
func DetectControlRows(c *Catalog, name string, fx Fixture, a Args, rows [][]any) error {
	group, ok := controlContracts[name]
	if !ok {
		return fmt.Errorf("unknown control %s", name)
	}
	got, err := Canonicalize(group, rows)
	if err != nil {
		return err
	}
	return DetectControl(c, name, fx, a, got, group)
}

func RunControls(ctx context.Context, q Querier, c *Catalog, fx Fixture, a Args) []ControlResult {
	var results []ControlResult
	for _, name := range c.Order {
		group, ok := controlContracts[name]
		if !ok {
			continue
		}
		s := c.Statements[name]
		result := ControlResult{Statement: name, Group: group}
		bound, _, err := Bind(s.SQL, a.ParamMap(s))
		if err == nil {
			result.Observed, err = q.Query(ctx, bound)
		}
		result.Rows = len(result.Observed)
		if err == nil {
			err = DetectControlRows(c, name, fx, a, result.Observed)
		}
		result.Detected = err == nil
		if err != nil {
			result.Error = err.Error()
		}
		results = append(results, result)
	}
	return results
}

var planningTime = regexp.MustCompile(`(?m)^Planning Time: ([0-9.]+) ms`)
var executionTime = regexp.MustCompile(`(?m)^Execution Time: ([0-9.]+) ms`)

func planTimes(plan string) (*float64, *float64) {
	read := func(re *regexp.Regexp) *float64 {
		match := re.FindStringSubmatch(plan)
		if len(match) != 2 {
			return nil
		}
		v, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return nil
		}
		return &v
	}
	return read(planningTime), read(executionTime)
}

func optionalMS(v *float64) string {
	if v == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.6f", *v)
}
