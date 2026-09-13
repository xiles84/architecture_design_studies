// Package provenance ties numbers to the data, the code and the analysts that
// interpreted them.
//
// Three questions a later reader must be able to answer about any report:
//
//  1. Which exact data is this?          -> Digest, a content hash of the results
//  2. Which exact code produced it?      -> RepoVersion, recorded at run time
//  3. Who has already interpreted it?    -> the signed analyses index
//
// (1) exists because a run id cannot identify data: runs get extended with
// extra cells after the fact. (2) exists because studies are revisited with new
// designs and questions over time, and a number is only reproducible from the
// commit that produced it. (3) exists because several analysts, including
// different AI models, are expected to read the same data.
package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// RepoVersion is the state of the repository a run was produced from. It is
// captured by the runner script (git is a host tool; the benchmark container
// has no repository) and passed into the harness, which copies it into every
// result file.
type RepoVersion struct {
	Commit string `json:"commit"`
	// Describe is `git describe --tags --always --dirty`.
	Describe string `json:"describe,omitempty"`
	// Dirty means uncommitted changes existed. A dirty run is not reproducible
	// from its commit alone, and reports say so rather than hiding it.
	Dirty bool `json:"dirty"`
	// Tag is the run tag created for this run, if any.
	Tag string `json:"tag,omitempty"`
}

// Digest is a stable content hash over every *.json file under fsys: independent
// of filesystem order, absolute paths and the machine computing it.
func Digest(fsys fs.FS) (string, int, error) {
	type entry struct{ rel, sum string }
	var entries []entry
	err := fs.WalkDir(fsys, ".", func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || path.Ext(p) != ".json" {
			return err
		}
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		entries = append(entries, entry{rel: p, sum: hex.EncodeToString(sum[:])})
		return nil
	})
	if err != nil {
		return "", 0, err
	}
	if len(entries) == 0 {
		return "", 0, fmt.Errorf("no result files")
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	h := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(h, "%s %s\n", e.rel, e.sum)
	}
	return hex.EncodeToString(h.Sum(nil))[:16], len(entries), nil
}

// Analysis is the frontmatter of one signed analysis.
type Analysis struct {
	File       string
	AnalysisID string
	RunID      string
	Analyst    string
	Kind       string
	Version    string
	AnalyzedAt string
	Digest     string
	// Inputs lists every run the analysis draws on (run id -> digest), so an
	// analysis spanning runs is indexed by each of them.
	Inputs   map[string]string
	Commit   string
	Status   string
	Headline string
}

// LoadAnalyses reads the frontmatter of every analysis in fsys that concerns
// runID (either as its primary run or listed in its inputs).
func LoadAnalyses(fsys fs.FS, runID string) ([]Analysis, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, nil // no analyses directory yet is a normal state
	}
	var out []Analysis
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "TEMPLATE") {
			continue
		}
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		a := ParseFrontmatter(string(b))
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

// ParseFrontmatter reads flat `key: value` pairs between two `---` fences. It is
// deliberately minimal so that adding an analysis never needs a YAML library:
// an analyst only has to be able to write a text file.
func ParseFrontmatter(src string) Analysis {
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
		case "repo_commit":
			a.Commit = v
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

// WriteSection renders the closing section of a generated report: the code
// version, the data fingerprint, and who has already interpreted it.
func WriteSection(b *strings.Builder, runID string, versions []RepoVersion, results fs.FS, analyses fs.FS) {
	digest, n, err := Digest(results)
	if err != nil {
		digest = "unavailable: " + err.Error()
	}

	fmt.Fprintf(b, "\n## Provenance: code, data and analyses\n\n")
	fmt.Fprintf(b, "This file is **generated from the measurements** and contains no interpretation.\n")
	fmt.Fprintf(b, "Conclusions live in separate signed analyses, so several analysts — including\n")
	fmt.Fprintf(b, "different AI models — can read the same numbers and each record what they make of\n")
	fmt.Fprintf(b, "them. Where two analyses disagree, the disagreement is a finding and stays visible.\n\n")

	fmt.Fprintf(b, "| | |\n|---|---|\n")
	fmt.Fprintf(b, "| Run id | `%s` |\n", runID)
	fmt.Fprintf(b, "| Result files | %d |\n", n)
	fmt.Fprintf(b, "| **Inputs digest** | `%s` |\n", digest)
	for _, v := range uniqueVersions(versions) {
		state := "clean"
		if v.Dirty {
			state = "⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone"
		}
		desc := v.Describe
		if desc == "" {
			desc = v.Commit
		}
		fmt.Fprintf(b, "| Repository version | `%s` (`%s`), %s |\n", desc, shortSHA(v.Commit), state)
		if v.Tag != "" {
			fmt.Fprintf(b, "| Run tag | `%s` |\n", v.Tag)
		}
	}
	if len(uniqueVersions(versions)) > 1 {
		fmt.Fprintf(b, "\n> ⚠️ Cells in this run were produced by **more than one repository version**. Compare\n")
		fmt.Fprintf(b, "> cells across versions only if the diff between them does not touch the harness or the SQL.\n")
	}
	fmt.Fprintf(b, "\nQuote the digest in any analysis: it identifies the data. The repository version\n")
	fmt.Fprintf(b, "identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or\n")
	fmt.Fprintf(b, "`git diff <commit>` to see what has changed in the study since.\n\n")

	list, err := LoadAnalyses(analyses, runID)
	if err != nil {
		fmt.Fprintf(b, "> Could not read the analyses directory: %v\n\n", err)
		return
	}
	fmt.Fprintf(b, "### Analyses of this run\n\n")
	if len(list) == 0 {
		fmt.Fprintf(b, "**None yet.** This run has measurements but no interpretation.\n\n")
	} else {
		fmt.Fprintf(b, "| Analyst | Kind | Date | Digest analysed | Status | Headline |\n|---|---|---|---|---|---|\n")
		for _, a := range list {
			match := "✅"
			if a.Digest != "" && a.Digest != digest {
				match = "⚠️ stale"
			}
			fmt.Fprintf(b, "| [%s](analyses/%s) | %s | %s | `%s` %s | %s | %s |\n",
				dash(a.Analyst), a.File, dash(a.Kind), dash(a.AnalyzedAt), dash(a.Digest), match, dash(a.Status), dash(a.Headline))
		}
		fmt.Fprintf(b, "\n")
	}
	fmt.Fprintf(b, "### Adding an analysis\n\n")
	fmt.Fprintf(b, "Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.\n")
	fmt.Fprintf(b, "If an analysis already exists for digest `%s`, read it first and write a new one only to\n", digest)
	fmt.Fprintf(b, "**disagree, extend, or bring a different perspective**. Never edit another analyst's file.\n")
}

func uniqueVersions(vs []RepoVersion) []RepoVersion {
	seen := map[string]bool{}
	var out []RepoVersion
	for _, v := range vs {
		k := fmt.Sprintf("%s|%v", v.Commit, v.Dirty)
		if v.Commit == "" || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, v)
	}
	return out
}

func shortSHA(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
