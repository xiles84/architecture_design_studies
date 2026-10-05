//go:build ignore

// Reader audit for Edition 3 preview 4 (task select-book-high). It reads the
// text that pdftotext extracted from the built PDF and checks that the new
// SELECT material rendered, that every SELECT claim is cited and carded, and
// that the stale wording it replaced is gone. It reads text only.
//
//	podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/book-select-high-20261005/check-reader.go docs/book-select-high-20261005/rendered.txt
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	// pdftotext breaks lines and may hyphenate; compare on whitespace-collapsed text.
	// Typst typesets straight apostrophes as typographic ones; compare on the plain form.
	text := strings.Join(strings.Fields(strings.ReplaceAll(string(b), "’", "'")), " ")
	fails := 0
	check := func(ok bool, format string, args ...any) {
		status := "ok  "
		if !ok {
			status = "FAIL"
			fails++
		}
		fmt.Printf("%s %s\n", status, fmt.Sprintf(format, args...))
	}

	for _, s := range []string{
		"Edition 3 preview 4",
		"Variation: join donations and filter them",
		"Variation: reuse one expensive CTE",
		"Variation: who has, or has not, donated",
		"Variation: donor totals and each donor's top donation",
		"Variation: page through donations, or match two ranges",
		"Variation: place the joined donations on three nodes",
		"Choosing a SELECT template",
		"repeated — no prior prediction",
		"Access paths for a disjunction.",
		"Charity placement variation",
	} {
		check(strings.Contains(text, s), "renders %q", s)
	}
	for _, s := range []string{
		"Unverified development build",
		"Planned variation: join donations",
		"no new SELECT ratios are asserted here",
		"repeated within one run",
		"does not yet promote the newer SELECT measurement families",
	} {
		check(!strings.Contains(text, s), "no longer renders %q", s)
	}
	ids := []string{
		"v6-sel-01-inner-join-spellings-same-plan",
		"v6-sel-02-materialized-joined-cte-penalty",
		"v6-sel-03-cte-reuse-crossing",
		"v6-sel-04-cte-reuse-postgresql-selective-no-crossing",
		"v6-sel-05-membership-join-distinct-penalty",
		"v6-sel-06-membership-in-exists-no-difference",
		"v6-sel-07-nullable-not-in-empty",
		"v6-sel-08-anti-join-not-exists-postgresql",
		"v6-sel-09-aggregate-correlated-postgresql-crossover",
		"v6-sel-10-aggregate-yugabytedb-no-reversal",
		"v6-sel-11-topn-lateral-probe-crossover",
		"v6-sel-12-topn-distinct-on-postgresql",
		"v6-sel-13-or-not-beaten-by-union",
		"v6-sel-14-expression-index-postgresql",
		"v6-sel-15-keyset-offset-crossover",
		"v6-sel-16-placement-splitting-effect",
		"v6-sel-17-placement-colocated-vs-distributed",
		"v6-sel-gap-01-yugabytedb-scan-only-and-single-node-families",
	}
	// A long id can wrap inside a badge; count occurrences after removing spaces
	// that a line break inserted into the hyphenated id.
	squeezed := strings.ReplaceAll(text, "- ", "-")
	for _, id := range ids {
		n := strings.Count(squeezed, id)
		check(n >= 2, "%s appears %d times (badge and card need at least 2)", id, n)
	}
	fmt.Printf("\n%d failure(s)\n", fails)
	if fails > 0 {
		os.Exit(1)
	}
}
