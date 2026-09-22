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
	claims := fs.String("claims", "book/evidence/v2/claims.json", "registry file (repo-relative; default is the active v2 correction package)")
	schemaPath := fs.String("schema", "book/evidence/v2/schema.json", "JSON Schema file (repo-relative)")
	confoundsPath := fs.String("confounds", "book/evidence/v2/confounds.json", "confound register for v2 (repo-relative)")
	v1Claims := fs.String("v1-claims", "book/evidence/claims.json", "historical v1 registry (repo-relative)")
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

	// The default registry is the v2 correction package. A caller inspecting the
	// historical v1 registry points --claims/--schema at book/evidence/claims.json
	// and book/evidence/claims.schema.json; the command dispatches on the file's
	// own schema_version so there is exactly one active source by default.
	peek, err := os.ReadFile(filepath.Join(root, *claims))
	if err != nil {
		fmt.Fprintf(errOut, "load claims: %v\n", err)
		return 1
	}
	var header struct {
		SchemaVersion int `json:"schema_version"`
	}
	_ = json.Unmarshal(peek, &header)

	if header.SchemaVersion == 2 {
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
