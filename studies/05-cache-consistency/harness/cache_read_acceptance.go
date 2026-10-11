package main

import (
	"fmt"
	"strings"
)

// rejectImpossibleReadEvidence enforces the committed-state
// floor using every available result source.
//
// The counters are deliberately NOT added together. Phase
// snapshots may overlap with cell-wide counters, and adding
// them would double-count the same invalid read.
//
// Any positive observation rejects the entire cell.
func rejectImpossibleReadEvidence(res *CellResult) error {
	if res == nil {
		return fmt.Errorf("IMPOSSIBLE CACHE VALUE: missing cell result")
	}

	var observations []string

	if res.ImpossibleValues > 0 {
		observations = append(observations,
			fmt.Sprintf("cell=%d", res.ImpossibleValues))
	}

	if res.Churn != nil && res.Churn.Impossible > 0 {
		observations = append(observations,
			fmt.Sprintf("churn=%d", res.Churn.Impossible))
	}

	for _, phase := range res.Wrong {
		if phase.Summary.ImpossibleValues > 0 {
			observations = append(observations,
				fmt.Sprintf("phase[%s]=%d",
					phase.Phase,
					phase.Summary.ImpossibleValues))
		}
	}

	for _, instance := range res.Instances {
		if instance.Impossible > 0 {
			observations = append(observations,
				fmt.Sprintf("instances[%d]=%d",
					instance.Instances,
					instance.Impossible))
		}
	}

	if len(observations) == 0 {
		return nil
	}

	return fmt.Errorf(
		"IMPOSSIBLE CACHE VALUE: cell rejected; recorded evidence: %s",
		strings.Join(observations, "; "),
	)
}
