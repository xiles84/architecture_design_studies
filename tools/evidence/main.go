// Command evidence validates the book claim-to-evidence registry and can dump
// candidate claims from signed analyses.
//
//	evidence validate [--repo DIR] [--claims FILE] [--schema FILE] [--list-inputs]
//	evidence extract  [--repo DIR]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"adsevidence/internal/evidence"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out, errOut io.Writer) int {
	cmd := "validate"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd, args = args[0], args[1:]
	}
	switch cmd {
	case "validate":
		return runValidate(args, out, errOut)
	case "extract":
		return runExtract(args, out, errOut)
	case "validate-v2":
		return runValidateV2(args, out, errOut)
	case "validate-v3":
		return runValidateV3(args, out, errOut)
	case "validate-v4":
		return runValidateV4(args, out, errOut)
	case "validate-v5":
		return runValidateV5(args, out, errOut)
	case "v1-diagnostic":
		return runV1Diagnostic(args, out, errOut)
	default:
		fmt.Fprintf(errOut, "unknown command %q\n", cmd)
		return 2
	}
}

func runValidate(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	claims := fs.String("claims", "book/evidence/v5/claims.json", "registry file (repo-relative; default is the active v5 package)")
	schemaPath := fs.String("schema", "book/evidence/v5/schema.json", "JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v5/confounds.json", "confound register for the active package (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "historical v1 registry (repo-relative)")
	v2Claims := fs.String("v2-claims", "book/evidence/v2/claims.json", "frozen v2 registry (repo-relative)")
	v3Claims := fs.String("v3-claims", "book/evidence/v3/claims.json", "frozen v3 registry, the v4 predecessor (repo-relative)")
	v4Claims := fs.String("v4-claims", "book/evidence/v4/claims.json", "frozen v4 registry, the v5 predecessor (repo-relative)")
	listInputs := fs.Bool("list-inputs", false, "print the claim ids that may appear in the book")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	rawSchema, err := os.ReadFile(filepath.Join(root, *schemaPath))
	if err != nil {
		fmt.Fprintf(errOut, "read schema: %v\n", err)
		return 2
	}
	var schema map[string]any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		fmt.Fprintf(errOut, "parse schema: %v\n", err)
		return 2
	}

	// The default registry is the v5 package (v4 with five claims re-scoped to
	// their own limits). A caller inspecting an older registry points
	// --claims/--schema at that package; the command dispatches on the file's own
	// schema_version so there is exactly one active source by default, and a v5
	// supersession resolves against v1, v2, v3 or v4.
	peek, err := os.ReadFile(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load claims: %v\n", err)
		return 1
	}
	var header struct {
		SchemaVersion int `json:"schema_version"`
	}
	_ = json.Unmarshal(peek, &header)

	if header.SchemaVersion == 2 || header.SchemaVersion == 3 || header.SchemaVersion == 4 || header.SchemaVersion == 5 {
		doc, raw, err := evidence.LoadV2(filepath.Join(root, *claims))
		if err != nil {
			fmt.Fprintf(errOut, "load versioned claims: %v\n", err)
			return 1
		}
		conf, err := evidence.LoadConfounds(filepath.Join(root, *confoundsPath))
		if err != nil {
			fmt.Fprintf(errOut, "load confounds: %v\n", err)
			return 1
		}
		v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
		if err != nil {
			fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
			return 1
		}
		var errs []error
		switch header.SchemaVersion {
		case 5:
			v2doc, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v2 claims: %v\n", err)
				return 1
			}
			v3doc, _, err := evidence.LoadV2(filepath.Join(root, *v3Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v3 claims: %v\n", err)
				return 1
			}
			v4doc, _, err := evidence.LoadV2(filepath.Join(root, *v4Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v4 claims: %v\n", err)
				return 1
			}
			errs = evidence.ValidateV5(root, schema, doc, raw, conf, v1, v2doc, v3doc, v4doc)
		case 4:
			v2doc, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v2 claims: %v\n", err)
				return 1
			}
			v3doc, _, err := evidence.LoadV2(filepath.Join(root, *v3Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v3 claims: %v\n", err)
				return 1
			}
			errs = evidence.ValidateV4(root, schema, doc, raw, conf, v1, v2doc, v3doc)
		case 3:
			v2doc, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
			if err != nil {
				fmt.Fprintf(errOut, "load predecessor v2 claims: %v\n", err)
				return 1
			}
			errs = evidence.ValidateV3(root, schema, doc, raw, conf, v1, v2doc)
		default:
			errs = evidence.ValidateV2(root, schema, doc, raw, conf, v1)
		}
		inputs := make([]string, 0, len(doc.Claims))
		for _, c := range doc.Claims {
			inputs = append(inputs, c.ClaimID)
		}
		if *listInputs {
			for _, id := range inputs {
				fmt.Fprintln(out, id)
			}
		}
		if *jsonOut {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			_ = enc.Encode(map[string]any{"claims": len(doc.Claims), "book_inputs": inputs, "errors": errStrings(errs)})
		} else if len(errs) == 0 {
			fmt.Fprintf(out, "evidence: %d claims, 0 errors, %d book inputs (active source %s)\n", len(doc.Claims), len(inputs), *claims)
		} else {
			for _, e := range errs {
				fmt.Fprintf(out, "  [ERROR] %v\n", e)
			}
			fmt.Fprintf(out, "evidence: %d claims, %d errors\n", len(doc.Claims), len(errs))
		}
		if len(errs) > 0 {
			return 1
		}
		return 0
	}

	doc, raw, err := evidence.Load(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load claims: %v\n", err)
		return 1
	}
	errs := evidence.Validate(root, schema, doc, raw)

	if *listInputs {
		for _, id := range doc.BookInputs() {
			fmt.Fprintln(out, id)
		}
	}
	if *jsonOut {
		report := map[string]any{
			"claims":      len(doc.Claims),
			"book_inputs": doc.BookInputs(),
			"errors":      errStrings(errs),
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(report)
	} else if len(errs) == 0 {
		fmt.Fprintf(out, "evidence: %d claims, 0 errors, %d book inputs (historical v1 %s)\n", len(doc.Claims), len(doc.BookInputs()), *claims)
	} else {
		for _, e := range errs {
			fmt.Fprintf(out, "  [ERROR] %v\n", e)
		}
		fmt.Fprintf(out, "evidence: %d claims, %d errors\n", len(doc.Claims), len(errs))
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

func runValidateV2(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("validate-v2", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	claims := fs.String("claims", "book/evidence/v2/claims.json", "v2 registry file (repo-relative)")
	schemaPath := fs.String("schema", "book/evidence/v2/schema.json", "v2 JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v2/confounds.json", "confound register (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "v1 registry to supersede/retire against (repo-relative)")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	rawSchema, err := os.ReadFile(filepath.Join(root, *schemaPath))
	if err != nil {
		fmt.Fprintf(errOut, "read schema: %v\n", err)
		return 2
	}
	var schema map[string]any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		fmt.Fprintf(errOut, "parse schema: %v\n", err)
		return 2
	}
	doc, raw, err := evidence.LoadV2(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v2 claims: %v\n", err)
		return 1
	}
	conf, err := evidence.LoadConfounds(filepath.Join(root, *confoundsPath))
	if err != nil {
		fmt.Fprintf(errOut, "load confounds: %v\n", err)
		return 1
	}
	v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
		return 1
	}
	errs := evidence.ValidateV2(root, schema, doc, raw, conf, v1)
	if *jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{"claims": len(doc.Claims), "errors": errStrings(errs)})
	} else if len(errs) == 0 {
		fmt.Fprintf(out, "evidence-v2: %d claims, 0 errors, active source book/evidence/v2/claims.json\n", len(doc.Claims))
	} else {
		for _, e := range errs {
			fmt.Fprintf(out, "  [ERROR] %v\n", e)
		}
		fmt.Fprintf(out, "evidence-v2: %d claims, %d errors\n", len(doc.Claims), len(errs))
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

// runValidateV5 resolves the v5 package explicitly, with all three frozen
// predecessors (v2, v3, v4) named and the v5 repetition lint enabled.
func runValidateV5(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("validate-v5", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	claims := fs.String("claims", "book/evidence/v5/claims.json", "v5 registry file (repo-relative)")
	schemaPath := fs.String("schema", "book/evidence/v5/schema.json", "v5 JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v5/confounds.json", "confound register (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "v1 registry, a predecessor (repo-relative)")
	v2Claims := fs.String("v2-claims", "book/evidence/v2/claims.json", "v2 registry, a predecessor (repo-relative)")
	v3Claims := fs.String("v3-claims", "book/evidence/v3/claims.json", "v3 registry, a predecessor (repo-relative)")
	v4Claims := fs.String("v4-claims", "book/evidence/v4/claims.json", "v4 registry, the direct predecessor (repo-relative)")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	rawSchema, err := os.ReadFile(filepath.Join(root, *schemaPath))
	if err != nil {
		fmt.Fprintf(errOut, "read schema: %v\n", err)
		return 2
	}
	var schema map[string]any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		fmt.Fprintf(errOut, "parse schema: %v\n", err)
		return 2
	}
	doc, raw, err := evidence.LoadV2(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v5 claims: %v\n", err)
		return 1
	}
	conf, err := evidence.LoadConfounds(filepath.Join(root, *confoundsPath))
	if err != nil {
		fmt.Fprintf(errOut, "load confounds: %v\n", err)
		return 1
	}
	v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
		return 1
	}
	v2, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v2 claims: %v\n", err)
		return 1
	}
	v3, _, err := evidence.LoadV2(filepath.Join(root, *v3Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v3 claims: %v\n", err)
		return 1
	}
	v4, _, err := evidence.LoadV2(filepath.Join(root, *v4Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v4 claims: %v\n", err)
		return 1
	}
	errs := evidence.ValidateV5(root, schema, doc, raw, conf, v1, v2, v3, v4)
	if *jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{"claims": len(doc.Claims), "errors": errStrings(errs)})
	} else if len(errs) == 0 {
		fmt.Fprintf(out, "evidence-v5: %d claims, 0 errors, active source book/evidence/v5/claims.json\n", len(doc.Claims))
	} else {
		for _, e := range errs {
			fmt.Fprintf(out, "  [ERROR] %v\n", e)
		}
		fmt.Fprintf(out, "evidence-v5: %d claims, %d errors\n", len(doc.Claims), len(errs))
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

// runValidateV4 resolves the v4 package explicitly. It is the same work the
// default `validate` performs, with all paths pinned to v4 and both frozen
// predecessors (v2 and v3) named, since a v4 retirement or supersession may
// point at a claim in either.
func runValidateV4(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("validate-v4", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	claims := fs.String("claims", "book/evidence/v4/claims.json", "v4 registry file (repo-relative)")
	schemaPath := fs.String("schema", "book/evidence/v4/schema.json", "v4 JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v4/confounds.json", "confound register (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "v1 registry, a predecessor (repo-relative)")
	v2Claims := fs.String("v2-claims", "book/evidence/v2/claims.json", "v2 registry, a predecessor (repo-relative)")
	v3Claims := fs.String("v3-claims", "book/evidence/v3/claims.json", "v3 registry, the direct predecessor (repo-relative)")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	rawSchema, err := os.ReadFile(filepath.Join(root, *schemaPath))
	if err != nil {
		fmt.Fprintf(errOut, "read schema: %v\n", err)
		return 2
	}
	var schema map[string]any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		fmt.Fprintf(errOut, "parse schema: %v\n", err)
		return 2
	}
	doc, raw, err := evidence.LoadV2(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v4 claims: %v\n", err)
		return 1
	}
	conf, err := evidence.LoadConfounds(filepath.Join(root, *confoundsPath))
	if err != nil {
		fmt.Fprintf(errOut, "load confounds: %v\n", err)
		return 1
	}
	v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
		return 1
	}
	v2, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v2 claims: %v\n", err)
		return 1
	}
	v3, _, err := evidence.LoadV2(filepath.Join(root, *v3Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v3 claims: %v\n", err)
		return 1
	}
	errs := evidence.ValidateV4(root, schema, doc, raw, conf, v1, v2, v3)
	if *jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{"claims": len(doc.Claims), "errors": errStrings(errs)})
	} else if len(errs) == 0 {
		fmt.Fprintf(out, "evidence-v4: %d claims, 0 errors, active source book/evidence/v4/claims.json\n", len(doc.Claims))
	} else {
		for _, e := range errs {
			fmt.Fprintf(out, "  [ERROR] %v\n", e)
		}
		fmt.Fprintf(out, "evidence-v4: %d claims, %d errors\n", len(doc.Claims), len(errs))
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

// runValidateV3 resolves the v3 package explicitly. It is the same work the
// default `validate` performs, with all paths pinned to v3 and the frozen v2
// registry named as the predecessor.
func runValidateV3(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("validate-v3", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	claims := fs.String("claims", "book/evidence/v3/claims.json", "v3 registry file (repo-relative)")
	schemaPath := fs.String("schema", "book/evidence/v3/schema.json", "v3 JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v3/confounds.json", "confound register (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "v1 registry, a predecessor (repo-relative)")
	v2Claims := fs.String("v2-claims", "book/evidence/v2/claims.json", "v2 registry, the direct predecessor (repo-relative)")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	rawSchema, err := os.ReadFile(filepath.Join(root, *schemaPath))
	if err != nil {
		fmt.Fprintf(errOut, "read schema: %v\n", err)
		return 2
	}
	var schema map[string]any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		fmt.Fprintf(errOut, "parse schema: %v\n", err)
		return 2
	}
	doc, raw, err := evidence.LoadV2(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v3 claims: %v\n", err)
		return 1
	}
	conf, err := evidence.LoadConfounds(filepath.Join(root, *confoundsPath))
	if err != nil {
		fmt.Fprintf(errOut, "load confounds: %v\n", err)
		return 1
	}
	v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
		return 1
	}
	v2, _, err := evidence.LoadV2(filepath.Join(root, *v2Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v2 claims: %v\n", err)
		return 1
	}
	errs := evidence.ValidateV3(root, schema, doc, raw, conf, v1, v2)
	if *jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{"claims": len(doc.Claims), "errors": errStrings(errs)})
	} else if len(errs) == 0 {
		fmt.Fprintf(out, "evidence-v3: %d claims, 0 errors, active source book/evidence/v3/claims.json\n", len(doc.Claims))
	} else {
		for _, e := range errs {
			fmt.Fprintf(out, "  [ERROR] %v\n", e)
		}
		fmt.Fprintf(out, "evidence-v3: %d claims, %d errors\n", len(doc.Claims), len(errs))
	}
	if len(errs) > 0 {
		return 1
	}
	return 0
}

func runV1Diagnostic(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("v1-diagnostic", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "v1 registry (repo-relative)")
	jsonOut := fs.Bool("json", false, "emit a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	v1, _, err := evidence.Load(filepath.Join(root, *v1Claims))
	if err != nil {
		fmt.Fprintf(errOut, "load v1 claims: %v\n", err)
		return 1
	}
	rep, err := evidence.V1Diagnostic(root, v1)
	if err != nil {
		fmt.Fprintf(errOut, "v1 diagnostic: %v\n", err)
		return 1
	}
	if *jsonOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
		return 0
	}
	for _, r := range rep.Claims {
		fmt.Fprintf(out, "%s: %d cells, %d unresolved (%s)\n", r.ClaimID, r.Cells, r.Unresolved, r.Report)
	}
	fmt.Fprintf(out, "v1-diagnostic: %d numeric claims, %d with a missing name, %d/%d names unresolved\n",
		rep.NumericClaims, rep.ClaimsWithMissing, rep.NamesMissing, rep.Names)
	return 0
}

func runExtract(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("extract", flag.ContinueOnError)
	fs.SetOutput(errOut)
	repo := fs.String("repo", ".", "repository root")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	pattern := filepath.Join(root, "studies", "*", "reports", "analyses", "*.md")
	files, _ := filepath.Glob(pattern)
	sort.Strings(files)
	type row struct {
		Study    string `json:"study"`
		RunID    string `json:"run_id"`
		Analysis string `json:"analysis"`
		Digest   string `json:"inputs_digest"`
		Commit   string `json:"repo_commit"`
		Headline string `json:"headline"`
	}
	var rows []row
	for _, f := range files {
		base := filepath.Base(f)
		if base == "TEMPLATE.md" || strings.Contains(f, string(filepath.Separator)+"outdated"+string(filepath.Separator)) {
			continue
		}
		text, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(root, f)
		parts := strings.Split(rel, string(filepath.Separator))
		study := ""
		if len(parts) > 1 {
			study = parts[1]
		}
		rows = append(rows, row{
			Study:    study,
			RunID:    frontmatter(string(text), "run_id"),
			Analysis: rel,
			Digest:   frontmatter(string(text), "inputs_digest"),
			Commit:   frontmatter(string(text), "repo_commit"),
			Headline: frontmatter(string(text), "headline"),
		})
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	_ = enc.Encode(rows)
	return 0
}

func frontmatter(text, key string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			break
		}
		if strings.HasPrefix(line, key+":") {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, key+":")), `"'`)
		}
	}
	return ""
}

func errStrings(errs []error) []string {
	var out []string
	for _, e := range errs {
		out = append(out, e.Error())
	}
	return out
}
