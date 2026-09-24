package evidence

// ValidateV3 resolves the v3 registry. v3 carries the still-active v2 claims
// forward and adds the Study 05 update, so a supersession or retirement may name
// a claim in either frozen predecessor: the v1 registry or the v2 registry. The
// resolver itself is shared with v2; only the allowed-predecessor set and the
// label in error messages differ.
func ValidateV3(repo string, schema map[string]any, doc *V2Document, raw any, conf *V2Confounds, v1 *Document, v2 *V2Document) []error {
	predecessors := map[string]bool{}
	if v1 != nil {
		for _, c := range v1.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	if v2 != nil {
		for _, c := range v2.Claims {
			predecessors[c.ClaimID] = true
		}
	}
	return validateRegistry(repo, schema, doc, raw, conf, predecessors, "predecessor")
}
