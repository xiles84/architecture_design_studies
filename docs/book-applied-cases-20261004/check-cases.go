// Source audit for the five teaching cases. Run in the repository's Go container.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	data, err := os.ReadFile("book/evidence/v6/claims.json")
	if err != nil { panic(err) }
	var registry struct { Claims []struct { ID string `json:"claim_id"` } `json:"claims"` }
	if err := json.Unmarshal(data, &registry); err != nil { panic(err) }
	active := map[string]bool{}
	for _, c := range registry.Claims { active[c.ID] = true }
	seen := map[string]bool{}
	badge := regexp.MustCompile(`claim-badge\("([^"]+)"\)`)
	label := regexp.MustCompile(`^== [^\n]+ <case-[^>]+>`)
	headings := 0
	files, err := filepath.Glob("book/cases/*.typ")
	if err != nil { panic(err) }
	if len(files) != 5 { panic(fmt.Sprintf("want five cases, have %d", len(files))) }
	for _, path := range files {
		b, err := os.ReadFile(path); if err != nil { panic(err) }
		s := string(b)
		if strings.Contains(s, "registry-card") { panic(path + ": full evidence card in teaching source") }
		local := map[string]bool{}
		for _, match := range badge.FindAllStringSubmatch(s, -1) {
			if !active[match[1]] { panic(path + ": unknown claim " + match[1]) }
			seen[match[1]], local[match[1]] = true, true
		}
		for _, line := range strings.Split(s, "\n") {
			if label.MatchString(line) { headings++ }
			if strings.HasPrefix(line, "=== Variation:") { continue }
			if strings.HasPrefix(line, "*Thrives") || strings.HasPrefix(line, "*Dies") || strings.HasPrefix(line, "*Observed") || strings.HasPrefix(line, "*Business") || strings.HasPrefix(line, "*Reproduce") || strings.HasPrefix(line, "Theory:") || (len(line) > 0 && line[0] >= 'A' && line[0] <= 'Z') {
				if count := len(strings.Fields(line)); count > 60 { panic(fmt.Sprintf("%s: paragraph has %d words", path, count)) }
			}
		}
		fmt.Printf("%s: %d unique valid badges; no full cards; labelled case heading\n", path, len(local))
	}
	if headings != 5 { panic(fmt.Sprintf("want five labelled case headings, have %d", headings)) }
	for id := range active { if !seen[id] { panic("uncited active claim: " + id) } }
	fmt.Printf("PASS: %d/%d active claims cited across five cases; audited prose lines <=60 words\n", len(seen), len(active))
}
