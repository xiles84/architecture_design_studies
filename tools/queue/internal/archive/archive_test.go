package archive

import (
	"strings"
	"testing"
	"time"

	"adsqueue/internal/model"
)

func actor() model.Identity {
	return model.Identity{Model: "test", Tool: "test", Effort: "unknown", SessionCapability: model.CapLOW, WorkRole: model.RoleExecutor, SessionID: "s"}
}

func ev(id string, seq int, prev *string, at, state, eventType string) *model.Event {
	return &model.Event{
		SchemaVersion:   model.SchemaVersion,
		EventID:         id,
		PreviousEventID: prev,
		Sequence:        seq,
		TaskID:          "task-20260922T000000Z-demo",
		GoalID:          "goal-20260922T000000Z-demo",
		RootTaskID:      "task-20260922T000000Z-demo",
		OccurredAt:      at,
		Actor:           actor(),
		EventType:       eventType,
		ResultingState:  state,
		EvidencePaths:   []string{},
	}
}

func testTask() *model.Task {
	return &model.Task{TaskID: "task-20260922T000000Z-demo"}
}

func TestValidateChainHappyPath(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2026-09-22T12:00:00Z")
	a := "event-1"
	b := "event-2"
	c := "event-3"
	events := []*model.Event{
		ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published"),
		ev("event-2", 2, &a, "2026-09-22T10:01:00Z", model.StateClaimed, "claimed"),
		ev("event-3", 3, &b, "2026-09-22T10:02:00Z", model.StateInProgress, "checkpoint"),
	}
	_ = c
	if errs := ValidateChain(testTask(), events, now); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := DeriveState(events); got != model.StateInProgress {
		t.Fatalf("DeriveState = %s", got)
	}
	if got := DeriveState(nil); got != model.StateProposed {
		t.Fatalf("empty chain should be proposed, got %s", got)
	}
}

func TestValidateChainFaults(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2026-09-22T12:00:00Z")
	a := "event-1"
	wrong := "event-nope"

	tests := []struct {
		name   string
		events []*model.Event
		want   string
	}{
		{
			name: "sequence gap",
			events: []*model.Event{
				ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published"),
				ev("event-3", 3, &a, "2026-09-22T10:01:00Z", model.StateClaimed, "claimed"),
			},
			want: "sequence",
		},
		{
			name: "fork / wrong predecessor",
			events: []*model.Event{
				ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published"),
				ev("event-2", 2, &wrong, "2026-09-22T10:01:00Z", model.StateClaimed, "claimed"),
			},
			want: "previous_event_id",
		},
		{
			name: "duplicate id",
			events: []*model.Event{
				ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published"),
				ev("event-1", 2, &a, "2026-09-22T10:01:00Z", model.StateClaimed, "claimed"),
			},
			want: "duplicate",
		},
		{
			name: "malformed timestamp",
			events: []*model.Event{
				ev("event-1", 1, nil, "yesterday", model.StateReady, "published"),
			},
			want: "occurred_at",
		},
		{
			name: "future clock",
			events: []*model.Event{
				ev("event-1", 1, nil, "2026-09-22T18:00:00Z", model.StateReady, "published"),
			},
			want: "future",
		},
		{
			name: "impossible transition",
			events: []*model.Event{
				ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published"),
				ev("event-2", 2, &a, "2026-09-22T10:01:00Z", model.StateCompleted, "completed"),
			},
			want: "impossible transition",
		},
		{
			name: "malformed actor",
			events: []*model.Event{
				func() *model.Event {
					e := ev("event-1", 1, nil, "2026-09-22T10:00:00Z", model.StateReady, "published")
					e.Actor = model.Identity{SessionCapability: "SIDEWAYS", WorkRole: model.RoleExecutor}
					return e
				}(),
			},
			want: "actor",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateChain(testTask(), tc.events, now)
			if len(errs) == 0 {
				t.Fatal("expected an error")
			}
			joined := ""
			for _, e := range errs {
				joined += e.Error() + " | "
			}
			if !strings.Contains(joined, tc.want) {
				t.Fatalf("errors %q do not mention %q", joined, tc.want)
			}
		})
	}
}

func TestEventIDAndFilename(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2026-09-22T03:19:27Z")
	taken := map[string]bool{}
	id := EventID("published", now, func(s string) bool { return taken[s] })
	if id != "event-20260922T031927Z-published" {
		t.Fatalf("EventID = %s", id)
	}
	e := ev(id, 1, nil, model.FormatTime(now), model.StateReady, "published")
	if got := EventFilename(e); got != "20260922T031927Z-event-20260922T031927Z-published.json" {
		t.Fatalf("EventFilename = %s", got)
	}
	taken[id] = true
	if got := EventID("published", now, func(s string) bool { return taken[s] }); got != "event-20260922T031927Z-published-2" {
		t.Fatalf("collision suffix = %s", got)
	}
}

func TestTaskDirDerivation(t *testing.T) {
	got, err := TaskDir("task-20260922T025912Z-queue-v1")
	if err != nil {
		t.Fatal(err)
	}
	if got != "docs/ai-work/tasks/2026/09/task-20260922T025912Z-queue-v1" {
		t.Fatalf("TaskDir = %s", got)
	}
	if _, err := TaskDir("not-a-task"); err == nil {
		t.Fatal("expected error for a malformed id")
	}
}
