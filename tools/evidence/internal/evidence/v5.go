package evidence

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidateV5 resolves the v5 registry. v5 carries v4 forward with five claims
// re-scoped, so a supersession or retirement may name a claim in any frozen
// predecessor: v1, v2, v3 or v4.
//
// On top of the shared resolver and the v4 lifecycle rules it adds one lint for
// the class the third review round exposed: a claim whose statement or limits
// cite more repetitions than the claim declares, with no corroborating anchor to
// back it. v2-07 said "stable across three runs" beside `Trials: 1` and a single
// anchor, which is the shape this catches.
func ValidateV5(repo string, schema map[string]any, doc *V2Document, raw any, conf *V2Confounds, v1 *Document, v2, v3, v4 *V2Document) []error {
	predecessors := map[string]bool{}
	if v1 != nil {
		for _, c := range v1.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	for _, d := range []*V2Document{v2, v3, v4} {
		if d == nil {
			continue
		}
		for _, c := range d.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	errs := validateRegistry(repo, schema, doc, raw, conf, predecessors, "predecessor")
	errs = append(errs, validateLifecycle(doc)...)
	return append(errs, validateRepetition(doc)...)
}

// repetitionWords are the number words a claim might write instead of a digit.
var repetitionWords = map[string]int{
	"one": 1, "two": 2, "three": 3, "four": 4, "five": 5,
	"six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10,
}

// repetitionRe finds "<count> <noun>" where the noun is a unit of repetition.
// It is deliberately shallow: it catches the v2-07 shape and nothing subtler, and
// §5 of book/evidence/v5/ANALYSIS.md says so. A claim asserting a cost, a
// direction or a physical outcome its limits disclaim is not detectable this way
// and still needs a reader.
var repetitionRe = regexp.MustCompile(`(?i)\b(one|two|three|four|five|six|seven|eight|nine|ten|[0-9]+)[ -](run|runs|replication|replications|trial|trials|measurement|measurements)\b`)

// validateRepetition fails a claim that cites more repetitions than it declares,
// unless it carries a corroborating anchor. Corroboration is the escape hatch on
// purpose: a claim that really was measured repeatedly should say so with anchors,
// not with a limit sentence.
func validateRepetition(doc *V2Document) []error {
	var errs []error
	for _, c := range doc.Claims {
		declared := c.Trials.WholeRunReplications
		if c.Trials.FreshLoadTrialsPerCell > declared {
			declared = c.Trials.FreshLoadTrialsPerCell
		}
		if c.Trials.InnerIterations > declared {
			declared = c.Trials.InnerIterations
		}
		corroborating := 0
		for _, a := range c.Anchors {
			if a.Role == "corroborating" {
				corroborating++
			}
		}

		text := c.Statement
		for _, l := range c.Limits {
			text += " " + l
		}
		seen := map[string]bool{}
		for _, m := range repetitionRe.FindAllStringSubmatch(text, -1) {
			count := 0
			if n, ok := repetitionWords[strings.ToLower(m[1])]; ok {
				count = n
			} else if n, err := strconv.Atoi(m[1]); err == nil {
				count = n
			}
			if count <= declared || corroborating > 0 {
				continue
			}
			phrase := strings.Join(m, "")
			if seen[phrase] {
				continue
			}
			seen[phrase] = true
			errs = append(errs, fmt.Errorf(
				"%s: cites %q but declares trials %d whole-run / %d fresh-load / %d inner with no corroborating anchor; either support the repetition with an anchor or state what is verifiable",
				c.ClaimID, phrase, c.Trials.WholeRunReplications, c.Trials.FreshLoadTrialsPerCell, c.Trials.InnerIterations))
		}
	}
	return errs
}
