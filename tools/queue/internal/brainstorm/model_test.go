package brainstorm

import (
	"strings"
	"testing"
	"time"

	"adsqueue/internal/model"
)

func TestParseIntentRequiresDeliberateFirstLine(t *testing.T) {
	tests := []struct {
		input, kind, id, replacement string
	}{
		{"START BRAINSTORM: How should we structure this?", "start", "", ""},
		{"\nCONTINUE BRAINSTORM brainstorm-20260922T100000Z-demo\nUse the repository evidence.", "continue", "brainstorm-20260922T100000Z-demo", ""},
		{"LIST BRAINSTORMS", "list", "", ""},
		{"BRAINSTORM STATUS brainstorm-20260922T100000Z-demo", "status", "brainstorm-20260922T100000Z-demo", ""},
		{"CREATE TASKS FROM BRAINSTORM brainstorm-20260922T100000Z-demo", "create_tasks", "brainstorm-20260922T100000Z-demo", ""},
		{"CANCEL BRAINSTORM brainstorm-20260922T100000Z-demo", "cancel", "brainstorm-20260922T100000Z-demo", ""},
		{"SUPERSEDE BRAINSTORM brainstorm-20260922T100000Z-old WITH brainstorm-20260922T110000Z-new", "supersede", "brainstorm-20260922T100000Z-old", "brainstorm-20260922T110000Z-new"},
		{"What can we improve in this document?", "", "", ""},
		{"Please START BRAINSTORM: maybe", "", "", ""},
		{"CONTINUE BRAINSTORM", "", "", ""},
		{"SUPERSEDE BRAINSTORM x", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt.input, " ", "_"), func(t *testing.T) {
			got := ParseIntent(tt.input)
			if got.Kind != tt.kind || got.BrainstormID != tt.id || got.ReplacementID != tt.replacement {
				t.Fatalf("ParseIntent(%q) = %#v", tt.input, got)
			}
		})
	}
	if got := ParseIntent("START BRAINSTORM: First line\nSecond line"); got.Question != "First line\nSecond line" {
		t.Fatalf("multiline start question = %q", got.Question)
	}
}

func TestValidateChainReconstructsStages(t *testing.T) {
	now := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)
	actor := model.Identity{Model: "test", Tool: "gotest", Effort: "unknown", SessionCapability: model.CapHIGH, WorkRole: model.RolePlanner, SessionID: "s0"}
	spec := &Spec{
		SchemaVersion: SchemaVersion, BrainstormID: "brainstorm-20260922T100000Z-demo", Title: "Demo", Question: "Choose.",
		PositionTarget: 2, CritiqueTarget: 1, CreatedAt: model.FormatTime(now), CreatedBy: actor,
	}
	events := []*Event{
		{SchemaVersion: 1, EventID: "e1", Sequence: 1, BrainstormID: spec.BrainstormID, OccurredAt: model.FormatTime(now), Actor: actor, EventType: "started", ResultingStage: StageCollectingPositions, Summary: "started"},
		{SchemaVersion: 1, EventID: "e2", PreviousEventID: ptr("e1"), Sequence: 2, BrainstormID: spec.BrainstormID, OccurredAt: model.FormatTime(now), Actor: actor, EventType: "position_submitted", ResultingStage: StageCollectingPositions, SlotID: "position-01", SlotKind: SlotPosition, ContributionID: "c2", Confidence: "high", Summary: "one"},
		{SchemaVersion: 1, EventID: "e3", PreviousEventID: ptr("e2"), Sequence: 3, BrainstormID: spec.BrainstormID, OccurredAt: model.FormatTime(now), Actor: actor, EventType: "position_submitted", ResultingStage: StageCrossReview, SlotID: "position-02", SlotKind: SlotPosition, ContributionID: "c3", Confidence: "high", Summary: "two"},
		{SchemaVersion: 1, EventID: "e4", PreviousEventID: ptr("e3"), Sequence: 4, BrainstormID: spec.BrainstormID, OccurredAt: model.FormatTime(now), Actor: actor, EventType: "critique_submitted", ResultingStage: StageSynthesis, SlotID: "critique-01", SlotKind: SlotCritique, ContributionID: "c4", Confidence: "medium", Summary: "reviewed"},
		{SchemaVersion: 1, EventID: "e5", PreviousEventID: ptr("e4"), Sequence: 5, BrainstormID: spec.BrainstormID, OccurredAt: model.FormatTime(now), Actor: actor, EventType: "synthesis_submitted", ResultingStage: StageConcluded, SlotID: "synthesis-01", SlotKind: SlotSynthesis, ContributionID: "c5", Confidence: "high", Summary: "done"},
	}
	if errs := ValidateChain(spec, events, now); len(errs) != 0 {
		t.Fatalf("valid chain rejected: %v", errs)
	}
	events[3].ResultingStage = StageCrossReview
	if errs := ValidateChain(spec, events, now); len(errs) == 0 {
		t.Fatal("invalid transition was accepted")
	}
	events[3].ResultingStage = StageSynthesis
	events[3].Actor.SessionCapability = model.CapLOW
	if errs := ValidateChain(spec, events, now); len(errs) == 0 {
		t.Fatal("LOW event actor was accepted")
	}
}

func ptr(s string) *string { return &s }
