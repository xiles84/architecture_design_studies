package evidence

import "fmt"

// ValidateV4 resolves the v4 registry. v4 carries the still-active v3 claims
// forward and adds claim lifecycle metadata, so a supersession or retirement may
// name a claim in any frozen predecessor: v1, v2 or v3. On top of the shared
// resolver it enforces the lifecycle rules that v4 exists to make mechanical,
// because the defect v4 repairs was not a bad number: v3 added the evidence that
// closed clauses of two active gaps and left the gaps printed beside it.
func ValidateV4(repo string, schema map[string]any, doc *V2Document, raw any, conf *V2Confounds, v1 *Document, v2, v3 *V2Document) []error {
	predecessors := map[string]bool{}
	if v1 != nil {
		for _, c := range v1.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	for _, d := range []*V2Document{v2, v3} {
		if d == nil {
			continue
		}
		for _, c := range d.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	errs := validateRegistry(repo, schema, doc, raw, conf, predecessors, "predecessor")
	return append(errs, validateLifecycle(doc)...)
}

// gapKinds are the v4 gap categories. They exist so the book can badge a
// schema limitation ("no representation can answer this") differently from a
// regime nobody ran, an instrument the harness lacks, or a measurement that
// exists but is too unstable to publish.
var gapKinds = map[string]bool{
	"coverage":             true,
	"schema_limitation":    true,
	"instrument":           true,
	"unstable_measurement": true,
}

// validateLifecycle enforces the v4-only rules:
//
//   - a claim in claims[] is active or partially superseded; retired and
//     superseded states are expressed by moving the claim into
//     retired_predecessor_claims, so a retired claim can never render as active
//     evidence in the book;
//   - a partially superseded claim names its successors and what remains open;
//   - a gap claim declares its kind, and a numeric claim does not pretend to be
//     one;
//   - every successor id resolves to an active claim in this same registry, so a
//     retirement cannot point at nothing;
//   - a retired predecessor is not also an active claim (no double accounting).
func validateLifecycle(doc *V2Document) []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	active := map[string]bool{}
	for _, c := range doc.Claims {
		active[c.ClaimID] = true
	}

	for _, c := range doc.Claims {
		switch c.Status {
		case "active":
		case "partially_superseded":
			if len(c.SupersededBy) == 0 {
				add("%s: status partially_superseded requires superseded_by", c.ClaimID)
			}
			if len(c.RemainingDimensions) == 0 {
				add("%s: status partially_superseded requires remaining_dimensions", c.ClaimID)
			}
		case "superseded", "retired":
			add("%s: status %q must not appear in claims[]; move it to retired_predecessor_claims", c.ClaimID, c.Status)
		case "":
			add("%s: missing status (v4 requires one of active, partially_superseded)", c.ClaimID)
		default:
			add("%s: status %q is not a v4 lifecycle state", c.ClaimID, c.Status)
		}

		if c.Strength == "gap" {
			if len(c.GapKind) == 0 {
				add("%s: gap claim must declare gap_kind", c.ClaimID)
			}
			for _, k := range c.GapKind {
				if !gapKinds[k] {
					add("%s: gap_kind %q is not one of coverage, schema_limitation, instrument, unstable_measurement", c.ClaimID, k)
				}
			}
		} else if len(c.GapKind) > 0 {
			add("%s: gap_kind is only allowed on gap claims (strength is %q)", c.ClaimID, c.Strength)
		}

		for _, id := range c.SupersededBy {
			if !active[id] {
				add("%s: superseded_by names %s, which is not an active claim in this registry", c.ClaimID, id)
			}
		}
	}

	for _, r := range doc.RetiredV1Claims {
		if active[r.ClaimID] {
			add("retired predecessor claim %s is also an active claim in this registry", r.ClaimID)
		}
	}
	for _, r := range doc.RetiredPredecessorClaims {
		if active[r.ClaimID] {
			add("retired predecessor claim %s is also an active claim in this registry", r.ClaimID)
		}
		for _, id := range r.SupersededBy {
			if !active[id] {
				add("retired predecessor claim %s: superseded_by names %s, which is not an active claim in this registry", r.ClaimID, id)
			}
		}
	}

	return errs
}
