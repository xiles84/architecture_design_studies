package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Analysis provenance
//
// A results directory holds numbers. What those numbers MEAN is a separate
// artefact, written by a person or a model, and this repository expects several
// of them: different models bring different priors, and disagreement between two
// analyses of the same data is a finding in itself.
//
// Two problems follow from that, and both are solved here:
//
//	1. Attribution — which analyst said this, when, looking at what?
//	2. Duplication — has this exact data already been analysed by someone?
//
// (2) is answered by hashing the inputs. The digest covers the content of every
// result file, so an analyst can ask "has digest abc123 already been analysed?"
// and get a reliable yes/no. A run id alone cannot do that: a run can be
// extended with extra cells after the fact, which is exactly what happened to
// this study's first matrix.
// ---------------------------------------------------------------------------

// ResultsDigest is a stable content hash over every result file in a run.
//
// Stable means: independent of filesystem order, of absolute paths, and of the
// machine computing it. Two people who hold the same results get the same digest.
func ResultsDigest(dir string) (string, int, error) {
	type entry struct{ rel, sum string }
	var entries []entry

	err := filepath.WalkDir(dir, func(path string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		entries = append(entries, entry{
			rel: filepath.ToSlash(rel),
			sum: hex.EncodeToString(sum[:]),
		})
		return nil
	})
	if err != nil {
		return "", 0, err
	}
	if len(entries) == 0 {
		return "", 0, fmt.Errorf("no result files under %s", dir)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })

	h := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(h, "%s %s\n", e.rel, e.sum)
	}
	return hex.EncodeToString(h.Sum(nil))[:16], len(entries), nil
}

// Analysis is the frontmatter of one written analysis of a run.
type Analysis struct {
	File       string
	AnalysisID string
	RunID      string
	Analyst    string
	Kind       string
	Version    string
	AnalyzedAt string
	Digest     string
	// Inputs lists every run the analysis draws on, as run id -> digest. An
	// analysis spanning several runs is indexed by each of them, and each report
	// checks staleness against the digest recorded for ITS run.
	Inputs   map[string]string
	Status   string
	Headline string
}

// LoadAnalyses reads the YAML-ish frontmatter of every analysis written about a
// run. The parser is deliberately minimal -- flat `key: value` pairs between two
// `---` fences -- so that adding an analysis never requires a YAML dependency or
// a build step. An analyst only has to be able to write a text file.
func LoadAnalyses(dir, runID string) ([]Analysis, error) {
	var out []Analysis
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "TEMPLATE") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		a := parseFrontmatter(string(b))
		a.File = name
		if runID != "" && a.RunID != runID {
			d, ok := a.Inputs[runID]
			if !ok {
				continue
			}
			a.Digest = d
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AnalyzedAt != out[j].AnalyzedAt {
			return out[i].AnalyzedAt < out[j].AnalyzedAt
		}
		return out[i].Analyst < out[j].Analyst
	})
	return out, nil
}

func parseFrontmatter(src string) Analysis {
	var a Analysis
	lines := strings.Split(src, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return a
	}
	for _, line := range lines[1:] {
		t := strings.TrimSpace(line)
		if t == "---" {
			break
		}
		k, v, ok := strings.Cut(t, ":")
		if !ok {
			continue
		}
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		switch strings.TrimSpace(k) {
		case "analysis_id":
			a.AnalysisID = v
		case "run_id":
			a.RunID = v
		case "analyst":
			a.Analyst = v
		case "analyst_kind":
			a.Kind = v
		case "analyst_version":
			a.Version = v
		case "analyzed_at":
			a.AnalyzedAt = v
		case "inputs_digest":
			a.Digest = v
		case "inputs":
			a.Inputs = map[string]string{}
			for _, pair := range strings.Split(v, ",") {
				if id, dg, ok := strings.Cut(strings.TrimSpace(pair), "@"); ok {
					a.Inputs[strings.TrimSpace(id)] = strings.TrimSpace(dg)
				}
			}
		case "status":
			a.Status = v
		case "headline":
			a.Headline = v
		}
	}
	return a
}

// writeProvenance renders the closing section of a generated report: what this
// data fingerprints to, who has already interpreted it, and how to add another
// interpretation without duplicating one that exists.
func writeProvenance(b *strings.Builder, resultsDir, analysesDir, runID string) {
	digest, n, err := ResultsDigest(resultsDir)
	if err != nil {
		digest = "unavailable: " + err.Error()
	}

	fmt.Fprintf(b, "\n## Conclusions and analysis provenance\n\n")
	fmt.Fprintf(b, "This file is **generated from the measurements** and contains no interpretation.\n")
	fmt.Fprintf(b, "Conclusions live in separate signed analyses, so that several analysts — including\n")
	fmt.Fprintf(b, "different AI models — can read the same numbers and each record what they make of\n")
	fmt.Fprintf(b, "them. Where two analyses disagree, that disagreement is itself a finding and is\n")
	fmt.Fprintf(b, "left visible rather than resolved by editing one of them.\n\n")

	fmt.Fprintf(b, "| | |\n|---|---|\n")
	fmt.Fprintf(b, "| Run id | `%s` |\n", runID)
	fmt.Fprintf(b, "| Result files | %d |\n", n)
	fmt.Fprintf(b, "| **Inputs digest** | `%s` |\n", digest)
	fmt.Fprintf(b, "\nThe digest is a content hash of every result file in the run. **Quote it in any\n")
	fmt.Fprintf(b, "analysis.** A run id alone is not enough to identify what was analysed — a run can\n")
	fmt.Fprintf(b, "be extended with extra cells afterwards — so the digest is what tells a later\n")
	fmt.Fprintf(b, "analyst whether they are looking at the same data someone else already wrote about.\n\n")

	analyses, err := LoadAnalyses(analysesDir, runID)
	if err != nil {
		fmt.Fprintf(b, "> Could not read the analyses directory: %v\n\n", err)
		return
	}

	fmt.Fprintf(b, "### Analyses of this run\n\n")
	if len(analyses) == 0 {
		fmt.Fprintf(b, "**None yet.** This run has measurements but no interpretation.\n\n")
	} else {
		fmt.Fprintf(b, "| Analyst | Kind | Date | Digest analysed | Status | Headline |\n")
		fmt.Fprintf(b, "|---|---|---|---|---|---|\n")
		for _, a := range analyses {
			match := "✅"
			if a.Digest != "" && a.Digest != digest {
				// The data moved after this analysis was written. Its conclusions
				// may still hold, but they were not drawn from what is here now.
				match = "⚠️ stale"
			}
			fmt.Fprintf(b, "| [%s](analyses/%s) | %s | %s | `%s` %s | %s | %s |\n",
				orDash(a.Analyst), a.File, orDash(a.Kind), orDash(a.AnalyzedAt),
				shortDigest(a.Digest), match, orDash(a.Status), orDash(a.Headline))
		}
		fmt.Fprintf(b, "\n")
	}

	fmt.Fprintf(b, "### Adding an analysis\n\n")
	fmt.Fprintf(b, "Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you\n")
	fmt.Fprintf(b, "conclude. Then regenerate this report so the table above picks it up.\n\n")
	fmt.Fprintf(b, "Before writing one, check the table: if an analysis already exists for digest\n")
	fmt.Fprintf(b, "`%s`, read it first. Add a new analysis to **disagree, extend, or bring a\n", digest)
	fmt.Fprintf(b, "different perspective** — not to restate what is already there. Never edit another\n")
	fmt.Fprintf(b, "analyst's file; write your own and reference theirs.\n")
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func shortDigest(s string) string {
	if s == "" {
		return "—"
	}
	if len(s) > 16 {
		return s[:16]
	}
	return s
}
