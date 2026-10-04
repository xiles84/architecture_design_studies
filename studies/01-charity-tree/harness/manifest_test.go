package main

import "testing"

func TestManifestReadsRecordedSELECTRunTag(t *testing.T) {
	for _, input := range []string{
		`{"provenance":{"RunTag":"run/study-01/corrected"}}`,
		`{"run_tag":"run/study-01/corrected"}`,
	} {
		if got := extractString([]byte(input), "run_tag"); got != "run/study-01/corrected" {
			t.Fatalf("producing tag lost: %q", got)
		}
	}
}
