// Package evidence validates the book claim-to-evidence registry against its
// committed JSON Schema and the repository it cites.
package evidence

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// SchemaError is one structural violation of the committed JSON Schema.
type SchemaError struct {
	Path string
	Msg  string
}

func (e SchemaError) Error() string { return fmt.Sprintf("%s: %s", e.Path, e.Msg) }

// ValidateSchema checks doc against a JSON-Schema subset: type, const, enum,
// pattern, minLength, minimum, minItems, required, additionalProperties,
// properties, items and local $ref into #/$defs. That subset is everything the
// committed claims schema uses, and it is evaluated here so validation needs no
// external dependency and works offline.
func ValidateSchema(schema map[string]any, doc any) []SchemaError {
	var errs []SchemaError
	walk(schema, schema, doc, "$", &errs)
	return errs
}

func walk(root, sch map[string]any, doc any, path string, errs *[]SchemaError) {
	if ref, ok := sch["$ref"].(string); ok {
		target := resolveRef(root, ref)
		if target == nil {
			*errs = append(*errs, SchemaError{path, "unresolvable $ref " + ref})
			return
		}
		walk(root, target, doc, path, errs)
		return
	}
	if t, ok := sch["type"]; ok && !typeMatches(t, doc) {
		*errs = append(*errs, SchemaError{path, fmt.Sprintf("expected type %v, got %s", t, typeName(doc))})
		return
	}
	if c, ok := sch["const"]; ok && !equalJSON(c, doc) {
		*errs = append(*errs, SchemaError{path, fmt.Sprintf("must equal %v", c)})
	}
	if en, ok := sch["enum"].([]any); ok {
		found := false
		for _, v := range en {
			if equalJSON(v, doc) {
				found = true
				break
			}
		}
		if !found {
			*errs = append(*errs, SchemaError{path, fmt.Sprintf("value %v not in enum %v", doc, en)})
		}
	}
	if s, ok := doc.(string); ok {
		if p, ok := sch["pattern"].(string); ok {
			if re, err := regexp.Compile(p); err == nil && !re.MatchString(s) {
				*errs = append(*errs, SchemaError{path, fmt.Sprintf("%q does not match %s", s, p)})
			}
		}
		if ml, ok := num(sch["minLength"]); ok && float64(len(s)) < ml {
			*errs = append(*errs, SchemaError{path, fmt.Sprintf("shorter than minLength %v", ml)})
		}
	}
	if n, ok := num(doc); ok {
		if m, ok := num(sch["minimum"]); ok && n < m {
			*errs = append(*errs, SchemaError{path, fmt.Sprintf("%v below minimum %v", n, m)})
		}
	}
	switch d := doc.(type) {
	case map[string]any:
		props, _ := sch["properties"].(map[string]any)
		if req, ok := sch["required"].([]any); ok {
			for _, r := range req {
				name, _ := r.(string)
				if _, present := d[name]; !present {
					*errs = append(*errs, SchemaError{path, "missing required property " + name})
				}
			}
		}
		keys := make([]string, 0, len(d))
		for k := range d {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ps, ok := props[k].(map[string]any)
			if !ok {
				if ap, ok := sch["additionalProperties"].(bool); ok && !ap {
					*errs = append(*errs, SchemaError{path, "unexpected property " + k})
				}
				continue
			}
			walk(root, ps, d[k], path+"."+k, errs)
		}
	case []any:
		if mi, ok := num(sch["minItems"]); ok && float64(len(d)) < mi {
			*errs = append(*errs, SchemaError{path, fmt.Sprintf("fewer than minItems %v", mi)})
		}
		if items, ok := sch["items"].(map[string]any); ok {
			for i, v := range d {
				walk(root, items, v, fmt.Sprintf("%s[%d]", path, i), errs)
			}
		}
	}
}

func resolveRef(root map[string]any, ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}
	var cur any = root
	for _, part := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur, ok = m[part]
		if !ok {
			return nil
		}
	}
	m, _ := cur.(map[string]any)
	return m
}

func typeMatches(t any, doc any) bool {
	names := []string{}
	switch v := t.(type) {
	case string:
		names = []string{v}
	case []any:
		for _, x := range v {
			names = append(names, fmt.Sprint(x))
		}
	}
	for _, n := range names {
		switch n {
		case "object":
			if _, ok := doc.(map[string]any); ok {
				return true
			}
		case "array":
			if _, ok := doc.([]any); ok {
				return true
			}
		case "string":
			if _, ok := doc.(string); ok {
				return true
			}
		case "integer":
			if n, ok := num(doc); ok && n == float64(int64(n)) {
				return true
			}
		case "number":
			if _, ok := num(doc); ok {
				return true
			}
		case "boolean":
			if _, ok := doc.(bool); ok {
				return true
			}
		case "null":
			if doc == nil {
				return true
			}
		}
	}
	return false
}

func num(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}

func equalJSON(a, b any) bool {
	if af, ok := num(a); ok {
		if bf, ok := num(b); ok {
			return af == bf
		}
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		if _, ok := num(v); ok {
			return "number"
		}
		return "unknown"
	}
}
