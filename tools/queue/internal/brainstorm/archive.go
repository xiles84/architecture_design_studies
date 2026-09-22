package brainstorm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"adsqueue/internal/archive"
	"adsqueue/internal/model"
)

const RecordsRoot = "docs/ai-work/brainstorms/records"

func Dir(root, id string) (string, error) {
	if !ValidID(id) {
		return "", fmt.Errorf("invalid brainstorm id %q", id)
	}
	date := strings.SplitN(id, "-", 3)[1]
	return filepath.Join(root, filepath.FromSlash(RecordsRoot), date[:4], date[4:6], id), nil
}

func Load(root, id string) (*Record, error) {
	dir, err := Dir(root, id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "brainstorm.json"))
	if err != nil {
		return nil, err
	}
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	events, err := LoadEvents(dir)
	if err != nil {
		return nil, err
	}
	stage := ""
	if len(events) > 0 {
		stage = events[len(events)-1].ResultingStage
	}
	return &Record{Spec: &spec, Dir: dir, Events: events, Stage: stage}, nil
}

func LoadAll(root string) ([]*Record, []error) {
	pattern := filepath.Join(root, filepath.FromSlash(RecordsRoot), "*", "*", "*", "brainstorm.json")
	paths, _ := filepath.Glob(pattern)
	var records []*Record
	var errs []error
	for _, p := range paths {
		var spec Spec
		data, err := os.ReadFile(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := json.Unmarshal(data, &spec); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", p, err))
			continue
		}
		r, err := Load(root, spec.BrainstormID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, r)
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Spec.BrainstormID < records[j].Spec.BrainstormID })
	return records, errs
}

func LoadEvents(dir string) ([]*Event, error) {
	paths, _ := filepath.Glob(filepath.Join(dir, "events", "*.json"))
	var out []*Event
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var ev Event
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		out = append(out, &ev)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Sequence < out[j].Sequence })
	return out, nil
}

func WriteSpec(dir string, spec *Spec) (string, error) {
	if err := spec.Validate(); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(dir, "events"), 0o755); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "brainstorm.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func WriteEvent(dir string, ev *Event) (string, error) {
	data, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "events", fmt.Sprintf("%s-%s.json", compact(ev.OccurredAt), ev.EventID))
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("event already exists: %s", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func NewEventID(kind string, now time.Time, events []*Event) string {
	base := "brainstorm-event-" + archive.CompactTime(now) + "-" + Slug(kind)
	seen := map[string]bool{}
	for _, ev := range events {
		seen[ev.EventID] = true
	}
	id := base
	for n := 2; seen[id]; n++ {
		id = fmt.Sprintf("%s-%d", base, n)
	}
	return id
}

func AppendEvent(r *Record, ev *Event) {
	r.Events = append(r.Events, ev)
	r.Stage = ev.ResultingStage
}

func compact(s string) string {
	if t, err := model.ParseTime(s); err == nil {
		return archive.CompactTime(t)
	}
	s = strings.ReplaceAll(s, ":", "")
	return strings.ReplaceAll(s, "-", "")
}

func WriteQuestion(dir string, spec *Spec) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# Brainstorm question\n\n")
	fmt.Fprintf(&b, "**Brainstorm:** `%s`  \n", spec.BrainstormID)
	fmt.Fprintf(&b, "**Title:** %s  \n", spec.Title)
	fmt.Fprintf(&b, "**Created:** `%s`\n\n", spec.CreatedAt)
	fmt.Fprintf(&b, "## Question\n\n%s\n", strings.TrimSpace(spec.Question))
	if len(spec.DecisionCriteria) > 0 {
		fmt.Fprintf(&b, "\n## Decision criteria\n\n")
		for _, s := range spec.DecisionCriteria {
			fmt.Fprintf(&b, "- %s\n", s)
		}
	}
	if len(spec.EvidencePaths) > 0 {
		fmt.Fprintf(&b, "\n## Evidence packet\n\n")
		for _, s := range spec.EvidencePaths {
			fmt.Fprintf(&b, "- `%s`\n", s)
		}
	}
	path := filepath.Join(dir, "QUESTION.md")
	return path, os.WriteFile(path, []byte(b.String()), 0o644)
}

func WriteContribution(r *Record, ev *Event, content []byte) (string, error) {
	var rel string
	switch ev.SlotKind {
	case SlotPosition:
		rel = filepath.Join("positions", ev.ContributionID+".md")
	case SlotCritique:
		rel = filepath.Join("critiques", ev.ContributionID+".md")
	case SlotSynthesis:
		rel = "CONCLUSION.md"
	default:
		return "", fmt.Errorf("unknown contribution kind %q", ev.SlotKind)
	}
	path := filepath.Join(r.Dir, rel)
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("contribution already exists: %s", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	var b strings.Builder
	title := map[string]string{SlotPosition: "Independent position", SlotCritique: "Cross-review", SlotSynthesis: "Brainstorm conclusion"}[ev.SlotKind]
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "| Field | Value |\n|---|---|\n")
	fmt.Fprintf(&b, "| Brainstorm | `%s` |\n", ev.BrainstormID)
	fmt.Fprintf(&b, "| Contribution | `%s` |\n", ev.ContributionID)
	fmt.Fprintf(&b, "| Slot | `%s` |\n", ev.SlotID)
	fmt.Fprintf(&b, "| Actor | %s |\n", ev.Actor.Verbose())
	if ev.Actor.SessionCapabilityInput != "" {
		fmt.Fprintf(&b, "| Capability input | %s |\n", ev.Actor.SessionCapabilityInput)
	}
	fmt.Fprintf(&b, "| Submitted | `%s` |\n", ev.OccurredAt)
	fmt.Fprintf(&b, "| Confidence | %s |\n", first(ev.Confidence, "not stated"))
	fmt.Fprintf(&b, "\n## Summary\n\n%s\n", ev.Summary)
	fmt.Fprintf(&b, "\n## Disagreements\n\n%s\n", first(ev.Disagreements, "None recorded."))
	fmt.Fprintf(&b, "\n## Contribution\n\n%s", strings.TrimSpace(string(content)))
	fmt.Fprintf(&b, "\n")
	return path, os.WriteFile(path, []byte(b.String()), 0o644)
}

func first(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}
