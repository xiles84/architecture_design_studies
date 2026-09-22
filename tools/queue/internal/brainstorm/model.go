// Package brainstorm contains the durable multi-leader deliberation model.
//
// Brainstorms are coordination records, not implementation tasks. Their
// immutable specification, contributions and chained events live under
// docs/ai-work/brainstorms; short-lived slot claims live in Git refs.
package brainstorm

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"adsqueue/internal/model"
)

const (
	SchemaVersion = 1

	StageCollectingPositions = "collecting_positions"
	StageCrossReview         = "cross_review"
	StageSynthesis           = "synthesis"
	StageConcluded           = "concluded"
	StageTasked              = "tasked"
	StageCancelled           = "cancelled"
	StageSuperseded          = "superseded"

	SlotPosition  = "position"
	SlotCritique  = "critique"
	SlotSynthesis = "synthesis"
)

// Spec is the immutable brainstorm.json.
type Spec struct {
	SchemaVersion    int            `json:"schema_version"`
	BrainstormID     string         `json:"brainstorm_id"`
	GoalID           string         `json:"goal_id,omitempty"`
	RootTaskID       string         `json:"root_task_id,omitempty"`
	ParentTaskID     string         `json:"parent_task_id,omitempty"`
	Title            string         `json:"title"`
	Question         string         `json:"question"`
	DecisionCriteria []string       `json:"decision_criteria"`
	EvidencePaths    []string       `json:"evidence_paths"`
	PositionTarget   int            `json:"position_target"`
	CritiqueTarget   int            `json:"critique_target"`
	CreatedAt        string         `json:"created_at"`
	CreatedBy        model.Identity `json:"created_by"`
}

func (s *Spec) Validate() error {
	if s.SchemaVersion != SchemaVersion {
		return fmt.Errorf("brainstorm %s: unsupported schema_version %d", s.BrainstormID, s.SchemaVersion)
	}
	if !ValidID(s.BrainstormID) {
		return fmt.Errorf("invalid brainstorm_id %q", s.BrainstormID)
	}
	if strings.TrimSpace(s.Title) == "" || strings.TrimSpace(s.Question) == "" {
		return errors.New("brainstorm title and question are required")
	}
	if s.PositionTarget < 2 || s.PositionTarget > 20 {
		return fmt.Errorf("position_target must be within 2..20, got %d", s.PositionTarget)
	}
	if s.CritiqueTarget < 1 || s.CritiqueTarget > 20 {
		return fmt.Errorf("critique_target must be within 1..20, got %d", s.CritiqueTarget)
	}
	if _, err := model.ParseTime(s.CreatedAt); err != nil {
		return fmt.Errorf("created_at: %w", err)
	}
	if err := s.CreatedBy.Validate(); err != nil {
		return fmt.Errorf("created_by: %w", err)
	}
	if s.CreatedBy.SessionCapability != model.CapHIGH {
		return errors.New("only HIGH capability may create a brainstorm")
	}
	if s.GoalID != "" && !model.ValidID(s.GoalID, "goal-") {
		return fmt.Errorf("invalid goal_id %q", s.GoalID)
	}
	for name, id := range map[string]string{"root_task_id": s.RootTaskID, "parent_task_id": s.ParentTaskID} {
		if id != "" && !model.ValidID(id, "task-") {
			return fmt.Errorf("invalid %s %q", name, id)
		}
	}
	return nil
}

// Event is one immutable brainstorm transition.
type Event struct {
	SchemaVersion   int            `json:"schema_version"`
	EventID         string         `json:"event_id"`
	PreviousEventID *string        `json:"previous_event_id"`
	Sequence        int            `json:"sequence"`
	BrainstormID    string         `json:"brainstorm_id"`
	OccurredAt      string         `json:"occurred_at"`
	Actor           model.Identity `json:"actor"`
	EventType       string         `json:"event_type"`
	ResultingStage  string         `json:"resulting_stage"`
	SlotID          string         `json:"slot_id,omitempty"`
	SlotKind        string         `json:"slot_kind,omitempty"`
	ContributionID  string         `json:"contribution_id,omitempty"`
	Summary         string         `json:"summary"`
	Disagreements   string         `json:"disagreements,omitempty"`
	Confidence      string         `json:"confidence,omitempty"`
	LinkedTaskIDs   []string       `json:"linked_task_ids,omitempty"`
	SupersededBy    string         `json:"superseded_by,omitempty"`
}

// Claim is the live CAS blob for one contribution slot.
type Claim struct {
	SchemaVersion int            `json:"schema_version"`
	BrainstormID  string         `json:"brainstorm_id"`
	SlotID        string         `json:"slot_id"`
	SlotKind      string         `json:"slot_kind"`
	ClaimID       string         `json:"claim_id"`
	ClaimEpoch    int            `json:"claim_epoch"`
	Worker        model.Identity `json:"worker"`
	ClaimedAt     string         `json:"claimed_at"`
	HeartbeatAt   string         `json:"heartbeat_at"`
	ExpiresAt     string         `json:"expires_at"`
	Stage         string         `json:"stage"`
}

func (c *Claim) Validate() error {
	if c.SchemaVersion != SchemaVersion {
		return fmt.Errorf("claim: unsupported schema_version %d", c.SchemaVersion)
	}
	if !ValidID(c.BrainstormID) {
		return fmt.Errorf("claim: invalid brainstorm_id %q", c.BrainstormID)
	}
	if c.ClaimID == "" || c.ClaimEpoch < 1 {
		return errors.New("claim: claim_id and positive claim_epoch are required")
	}
	if err := c.Worker.Validate(); err != nil {
		return fmt.Errorf("claim worker: %w", err)
	}
	if c.Worker.SessionCapability != model.CapHIGH {
		return errors.New("claim worker must be HIGH")
	}
	wantStage := ""
	switch c.SlotKind {
	case SlotPosition:
		wantStage = StageCollectingPositions
	case SlotCritique:
		wantStage = StageCrossReview
	case SlotSynthesis:
		wantStage = StageSynthesis
	default:
		return fmt.Errorf("claim: invalid slot_kind %q", c.SlotKind)
	}
	if !regexp.MustCompile("^" + c.SlotKind + "-[0-9]{2}$").MatchString(c.SlotID) {
		return fmt.Errorf("claim: invalid slot_id %q", c.SlotID)
	}
	if c.Stage != wantStage {
		return fmt.Errorf("claim: stage %q does not match %s slot", c.Stage, c.SlotKind)
	}
	claimed, err := model.ParseTime(c.ClaimedAt)
	if err != nil {
		return fmt.Errorf("claim claimed_at: %w", err)
	}
	if _, err := model.ParseTime(c.HeartbeatAt); err != nil {
		return fmt.Errorf("claim heartbeat_at: %w", err)
	}
	expires, err := model.ParseTime(c.ExpiresAt)
	if err != nil {
		return fmt.Errorf("claim expires_at: %w", err)
	}
	if !expires.After(claimed) {
		return errors.New("claim: expires_at must follow claimed_at")
	}
	return nil
}

// Record is a loaded spec plus its derived event state.
type Record struct {
	Spec   *Spec
	Dir    string
	Events []*Event
	Stage  string
}

func (r *Record) LastEvent() *Event {
	if len(r.Events) == 0 {
		return nil
	}
	return r.Events[len(r.Events)-1]
}

func (r *Record) SubmittedSlots() map[string]bool {
	out := map[string]bool{}
	for _, ev := range r.Events {
		if ev.SlotID != "" && (ev.EventType == "position_submitted" || ev.EventType == "critique_submitted" || ev.EventType == "synthesis_submitted") {
			out[ev.SlotID] = true
		}
	}
	return out
}

func (r *Record) Counts() (positions, critiques int, synthesized bool) {
	for _, ev := range r.Events {
		switch ev.EventType {
		case "position_submitted":
			positions++
		case "critique_submitted":
			critiques++
		case "synthesis_submitted":
			synthesized = true
		}
	}
	return
}

func (r *Record) OverallState() string {
	switch r.Stage {
	case StageCollectingPositions, StageCrossReview, StageSynthesis:
		return "ongoing"
	case StageConcluded:
		return "concluded"
	default:
		return r.Stage
	}
}

func (r *Record) Progress() string {
	p, c, s := r.Counts()
	synth := "pending"
	if s {
		synth = "complete"
	}
	return fmt.Sprintf("positions %d/%d; critiques %d/%d; synthesis %s", p, r.Spec.PositionTarget, c, r.Spec.CritiqueTarget, synth)
}

func (r *Record) Summary() string {
	if ev := r.LastEvent(); ev != nil && strings.TrimSpace(ev.Summary) != "" {
		return strings.TrimSpace(ev.Summary)
	}
	return strings.TrimSpace(r.Spec.Question)
}

func (r *Record) Disagreements() string {
	for i := len(r.Events) - 1; i >= 0; i-- {
		if s := strings.TrimSpace(r.Events[i].Disagreements); s != "" {
			return s
		}
	}
	return "none recorded"
}

func (r *Record) Confidence() string {
	for i := len(r.Events) - 1; i >= 0; i-- {
		if s := strings.TrimSpace(r.Events[i].Confidence); s != "" {
			return s
		}
	}
	return "not yet assessed"
}

func (r *Record) LinkedTasks() []string {
	seen := map[string]bool{}
	var out []string
	for _, ev := range r.Events {
		for _, id := range ev.LinkedTaskIDs {
			if !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}
	sort.Strings(out)
	return out
}

func (r *Record) NextAction() string {
	switch r.Stage {
	case StageCollectingPositions, StageCrossReview, StageSynthesis:
		return "CONTINUE BRAINSTORM " + r.Spec.BrainstormID
	case StageConcluded:
		return "CREATE TASKS FROM BRAINSTORM " + r.Spec.BrainstormID
	case StageTasked:
		return "Review linked task progress"
	default:
		return "none"
	}
}

// NextSlots returns all slots that may be claimed in the current stage.
func (r *Record) NextSlots() []Slot {
	submitted := r.SubmittedSlots()
	var out []Slot
	add := func(kind string, n int) {
		for i := 1; i <= n; i++ {
			id := fmt.Sprintf("%s-%02d", kind, i)
			if !submitted[id] {
				out = append(out, Slot{ID: id, Kind: kind})
			}
		}
	}
	switch r.Stage {
	case StageCollectingPositions:
		add(SlotPosition, r.Spec.PositionTarget)
	case StageCrossReview:
		add(SlotCritique, r.Spec.CritiqueTarget)
	case StageSynthesis:
		add(SlotSynthesis, 1)
	}
	return out
}

type Slot struct{ ID, Kind string }

func StageAfterSubmission(r *Record, kind string) (string, error) {
	p, c, _ := r.Counts()
	switch kind {
	case SlotPosition:
		if r.Stage != StageCollectingPositions {
			return "", fmt.Errorf("a position cannot be submitted during %s", r.Stage)
		}
		if p+1 >= r.Spec.PositionTarget {
			return StageCrossReview, nil
		}
		return StageCollectingPositions, nil
	case SlotCritique:
		if r.Stage != StageCrossReview {
			return "", fmt.Errorf("a critique cannot be submitted during %s", r.Stage)
		}
		if c+1 >= r.Spec.CritiqueTarget {
			return StageSynthesis, nil
		}
		return StageCrossReview, nil
	case SlotSynthesis:
		if r.Stage != StageSynthesis {
			return "", fmt.Errorf("a synthesis cannot be submitted during %s", r.Stage)
		}
		return StageConcluded, nil
	default:
		return "", fmt.Errorf("unknown slot kind %q", kind)
	}
}

var idRe = regexp.MustCompile(`^brainstorm-[0-9]{8}T[0-9]{6}Z-[a-z0-9][a-z0-9-]*$`)

func ValidID(id string) bool { return idRe.MatchString(id) }

func Slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 48 {
		s = strings.Trim(s[:48], "-")
	}
	if s == "" {
		return "discussion"
	}
	return s
}

// ValidateChain reconstructs and checks one brainstorm event chain.
func ValidateChain(spec *Spec, events []*Event, now time.Time) []error {
	var errs []error
	add := func(f string, a ...any) { errs = append(errs, fmt.Errorf(f, a...)) }
	if err := spec.Validate(); err != nil {
		add("%v", err)
	}
	stage := ""
	positions, critiques := 0, 0
	seenIDs, seenSlots := map[string]bool{}, map[string]bool{}
	var prev *Event
	for i, ev := range events {
		if ev.SchemaVersion != SchemaVersion {
			add("event %s: unsupported schema_version %d", ev.EventID, ev.SchemaVersion)
		}
		if seenIDs[ev.EventID] {
			add("duplicate event_id %s", ev.EventID)
		}
		seenIDs[ev.EventID] = true
		if ev.BrainstormID != spec.BrainstormID {
			add("event %s: brainstorm_id mismatch", ev.EventID)
		}
		if ev.Sequence != i+1 {
			add("event %s: sequence %d, expected %d", ev.EventID, ev.Sequence, i+1)
		}
		if i == 0 {
			if ev.PreviousEventID != nil {
				add("event %s: first event has a predecessor", ev.EventID)
			}
		} else if ev.PreviousEventID == nil || *ev.PreviousEventID != prev.EventID {
			add("event %s: predecessor does not name %s", ev.EventID, prev.EventID)
		}
		at, err := model.ParseTime(ev.OccurredAt)
		if err != nil {
			add("event %s: %v", ev.EventID, err)
		} else if at.After(now.Add(5 * time.Minute)) {
			add("event %s: occurred_at is in the future", ev.EventID)
		}
		if err := ev.Actor.Validate(); err != nil {
			add("event %s: actor: %v", ev.EventID, err)
		} else if ev.Actor.SessionCapability != model.CapHIGH {
			add("event %s: brainstorm actor must be HIGH", ev.EventID)
		}

		expected := stage
		switch ev.EventType {
		case "started":
			if i != 0 {
				add("event %s: started is not first", ev.EventID)
			}
			expected = StageCollectingPositions
		case "position_submitted":
			if stage != StageCollectingPositions || ev.SlotKind != SlotPosition {
				add("event %s: position submitted during %s", ev.EventID, stage)
			}
			positions++
			expected = StageCollectingPositions
			if positions >= spec.PositionTarget {
				expected = StageCrossReview
			}
		case "critique_submitted":
			if stage != StageCrossReview || ev.SlotKind != SlotCritique {
				add("event %s: critique submitted during %s", ev.EventID, stage)
			}
			critiques++
			expected = StageCrossReview
			if critiques >= spec.CritiqueTarget {
				expected = StageSynthesis
			}
		case "synthesis_submitted":
			if stage != StageSynthesis || ev.SlotKind != SlotSynthesis {
				add("event %s: synthesis submitted during %s", ev.EventID, stage)
			}
			expected = StageConcluded
		case "tasks_linked":
			if stage != StageConcluded && stage != StageTasked {
				add("event %s: tasks linked during %s", ev.EventID, stage)
			}
			expected = StageTasked
		case "cancelled":
			if stage != StageCollectingPositions && stage != StageCrossReview && stage != StageSynthesis {
				add("event %s: cancellation during %s", ev.EventID, stage)
			}
			expected = StageCancelled
		case "superseded":
			if stage != StageConcluded && stage != StageTasked {
				add("event %s: superseded during %s", ev.EventID, stage)
			}
			expected = StageSuperseded
		default:
			add("event %s: unknown event_type %q", ev.EventID, ev.EventType)
		}
		if ev.SlotID != "" {
			if seenSlots[ev.SlotID] {
				add("event %s: duplicate slot %s", ev.EventID, ev.SlotID)
			}
			seenSlots[ev.SlotID] = true
		}
		if ev.ResultingStage != expected {
			add("event %s: resulting_stage %s, expected %s", ev.EventID, ev.ResultingStage, expected)
		}
		stage = ev.ResultingStage
		prev = ev
	}
	return errs
}

// Intent is the deliberately narrow natural-language trigger recognized by
// repository agents. Ordinary questions return Kind="".
type Intent struct {
	Kind          string
	BrainstormID  string
	Question      string
	ReplacementID string
}

func ParseIntent(input string) Intent {
	s := strings.TrimSpace(input)
	if s == "" {
		return Intent{}
	}
	firstLine := s
	if i := strings.IndexByte(firstLine, '\n'); i >= 0 {
		firstLine = strings.TrimSpace(firstLine[:i])
	}
	upper := strings.ToUpper(firstLine)
	fields := strings.Fields(firstLine)
	switch {
	case strings.HasPrefix(upper, "START BRAINSTORM:"):
		return Intent{Kind: "start", Question: strings.TrimSpace(s[len("START BRAINSTORM:"):])}
	case strings.HasPrefix(upper, "CONTINUE BRAINSTORM "):
		if len(fields) == 3 {
			return Intent{Kind: "continue", BrainstormID: fields[2]}
		}
	case upper == "LIST BRAINSTORMS" || upper == "LIST BRAINSTORMS.":
		return Intent{Kind: "list"}
	case strings.HasPrefix(upper, "BRAINSTORM STATUS "):
		if len(fields) == 3 {
			return Intent{Kind: "status", BrainstormID: fields[2]}
		}
	case strings.HasPrefix(upper, "CREATE TASKS FROM BRAINSTORM "):
		if len(fields) == 5 {
			return Intent{Kind: "create_tasks", BrainstormID: fields[4]}
		}
	case strings.HasPrefix(upper, "CANCEL BRAINSTORM "):
		if len(fields) == 3 {
			return Intent{Kind: "cancel", BrainstormID: fields[2]}
		}
	case strings.HasPrefix(upper, "SUPERSEDE BRAINSTORM "):
		if len(fields) == 5 && strings.EqualFold(fields[3], "WITH") {
			return Intent{Kind: "supersede", BrainstormID: fields[2], ReplacementID: fields[4]}
		}
	}
	return Intent{}
}
