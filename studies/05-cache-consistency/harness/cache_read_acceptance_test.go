package main

import (
	"strings"
	"testing"
)

func TestImpossibleReadEvidenceRejectsEverySource(t *testing.T) {
	cases := []struct {
		name string
		res  CellResult
		fail bool
	}{
		{
			name: "clean",
			res:  CellResult{},
			fail: false,
		},
		{
			name: "authoritative atomic count",
			res:  CellResult{ImpossibleValues: 1},
			fail: true,
		},
		{
			name: "churn despite zero cell total",
			res: CellResult{
				ImpossibleValues: 0,
				Churn: &ChurnResult{
					Impossible: 1,
				},
			},
			fail: true,
		},
		{
			name: "phase despite zero cell total",
			res: CellResult{
				Wrong: []PhaseWrong{
					{
						Phase: "churn_90_reads",
						Summary: WrongReadSummary{
							ImpossibleValues: 1,
						},
					},
				},
			},
			fail: true,
		},
		{
			name: "instance despite zero cell total",
			res: CellResult{
				Instances: []InstanceResult{
					{
						Instances:  3,
						Impossible: 1,
					},
				},
			},
			fail: true,
		},
		{
			name: "overlapping evidence",
			res: CellResult{
				ImpossibleValues: 1,
				Churn: &ChurnResult{
					Impossible: 1,
				},
			},
			fail: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := rejectImpossibleReadEvidence(&tc.res)

			if tc.fail && err == nil {
				t.Fatal("invalid cell was accepted")
			}

			if !tc.fail && err != nil {
				t.Fatalf("clean cell rejected: %v", err)
			}

			if err != nil &&
				!strings.Contains(err.Error(), "IMPOSSIBLE CACHE VALUE") {
				t.Fatalf("unexpected failure: %v", err)
			}
		})
	}
}

func TestImpossibleReadEvidenceMissingResultFailsClosed(t *testing.T) {
	if err := rejectImpossibleReadEvidence(nil); err == nil {
		t.Fatal("missing result was accepted")
	}
}
