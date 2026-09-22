package model

import (
	"testing"
	"time"
)

func mk(id string, priority int, created string, notBefore string, minCap string, deps ...string) *TaskRecord {
	return &TaskRecord{
		Task: &Task{
			TaskID:            id,
			Priority:          priority,
			CreatedAt:         created,
			NotBefore:         notBefore,
			MinimumCapability: minCap,
			Dependencies:      deps,
		},
		State: StateReady,
	}
}

func TestQueueOrdering(t *testing.T) {
	base := "2026-09-22T00:00:00Z"
	now, _ := ParseTime(base)
	records := []*TaskRecord{
		mk("task-20260922T000001Z-a", 10, "2026-09-22T00:00:01Z", base, CapLOW),
		mk("task-20260922T000002Z-b", 90, "2026-09-22T00:00:02Z", base, CapLOW),
		mk("task-20260922T000003Z-c", 90, "2026-09-22T00:00:01Z", base, CapLOW), // same priority, earlier
		mk("task-20260922T000004Z-d", 90, "2026-09-22T00:00:01Z", base, CapLOW), // same priority+time, lower id
		mk("task-20260922T000004Z-zz", 90, "2026-09-22T00:00:01Z", base, CapLOW),
	}
	// Force the two same-time tasks to a deterministic id order.
	records[3].Task.TaskID = "task-20260922T000004Z-a"
	records[4].Task.TaskID = "task-20260922T000004Z-b"
	got := Next(records, now, CapLOW)
	if got == nil {
		t.Fatal("expected a task")
	}
	if got.Task.TaskID != "task-20260922T000003Z-c" {
		t.Fatalf("expected priority then earliest creation: got %s", got.Task.TaskID)
	}
	ordered := Eligible(records, now, CapLOW)
	want := []string{
		"task-20260922T000003Z-c",
		"task-20260922T000004Z-a",
		"task-20260922T000004Z-b",
		"task-20260922T000002Z-b",
		"task-20260922T000001Z-a",
	}
	for i, w := range want {
		if ordered[i].Task.TaskID != w {
			t.Fatalf("order[%d] = %s, want %s", i, ordered[i].Task.TaskID, w)
		}
	}
}

func TestNotBeforeBoundary(t *testing.T) {
	now, _ := ParseTime("2026-09-22T12:00:00Z")
	at := mk("task-20260922T000001Z-at", 50, "2026-09-22T00:00:00Z", "2026-09-22T12:00:00Z", CapLOW)
	future := mk("task-20260922T000002Z-future", 50, "2026-09-22T00:00:00Z", "2026-09-22T12:00:01Z", CapLOW)
	if got := Next([]*TaskRecord{at}, now, CapLOW); got == nil {
		t.Fatal("not_before exactly at now must be eligible (UTC boundary)")
	}
	if got := Next([]*TaskRecord{future}, now, CapLOW); got != nil {
		t.Fatal("not_before one second ahead must not be eligible")
	}
}

func TestCapabilityFiltering(t *testing.T) {
	now, _ := ParseTime("2026-09-22T12:00:00Z")
	low := mk("task-20260922T000001Z-low", 50, "2026-09-22T00:00:00Z", "2026-09-22T00:00:00Z", CapLOW)
	high := mk("task-20260922T000002Z-high", 50, "2026-09-22T00:00:00Z", "2026-09-22T00:00:00Z", CapHIGH)
	if CapabilityEligible(CapLOW, CapHIGH) {
		t.Fatal("LOW must not claim a HIGH-only task")
	}
	if !CapabilityEligible(CapHIGH, CapLOW) {
		t.Fatal("HIGH must be able to claim a LOW task")
	}
	if got := Next([]*TaskRecord{low, high}, now, CapLOW); got == nil || got.Task.TaskID != low.Task.TaskID {
		t.Fatalf("LOW session picked %v", got)
	}
	if got := Next([]*TaskRecord{low, high}, now, CapHIGH); got == nil || got.Task.TaskID != high.Task.TaskID {
		// Both priority 50, high created later, so FIFO picks low first.
		if got == nil || got.Task.TaskID != low.Task.TaskID {
			t.Fatalf("HIGH session picked %v", got)
		}
	}
}

func TestDependencies(t *testing.T) {
	now, _ := ParseTime("2026-09-22T12:00:00Z")
	dep := mk("task-20260922T000001Z-dep", 50, "2026-09-22T00:00:00Z", "2026-09-22T00:00:00Z", CapLOW)
	dep.State = StateInProgress
	blocked := mk("task-20260922T000002Z-blocked", 99, "2026-09-22T00:00:00Z", "2026-09-22T00:00:00Z", CapLOW, dep.Task.TaskID)
	if got := Next([]*TaskRecord{dep, blocked}, now, CapLOW); got != nil {
		t.Fatalf("incomplete dependency must block; got %v", got)
	}
	dep.State = StateCompleted
	if got := Next([]*TaskRecord{dep, blocked}, now, CapLOW); got == nil || got.Task.TaskID != blocked.Task.TaskID {
		t.Fatalf("completed dependency should unblock; got %v", got)
	}
}

func TestDurations(t *testing.T) {
	cases := map[string]time.Duration{
		"7d":    7 * 24 * time.Hour,
		"36h":   36 * time.Hour,
		"1d12h": 36 * time.Hour,
		"90m":   90 * time.Minute,
		"30s":   30 * time.Second,
	}
	for in, want := range cases {
		got, err := ParseDuration(in)
		if err != nil {
			t.Fatalf("ParseDuration(%q): %v", in, err)
		}
		if got != want {
			t.Fatalf("ParseDuration(%q) = %v, want %v", in, got, want)
		}
	}
	for _, bad := range []string{"", "7", "d", "-3d", "1x"} {
		if _, err := ParseDuration(bad); err == nil {
			t.Fatalf("ParseDuration(%q) should fail", bad)
		}
	}
}

func TestTaskDigestStableAcrossLineEndings(t *testing.T) {
	unix := []byte("{\n  \"a\": 1\n}\n")
	dos := []byte("{\r\n  \"a\": 1\r\n}\r\n")
	if TaskDigest(unix) != TaskDigest(dos) {
		t.Fatal("CRLF checkout must not change the task digest")
	}
}

func TestTransitions(t *testing.T) {
	allowed := []Transition{
		{StateProposed, StateReady},
		{StateReady, StateClaimed},
		{StateClaimed, StateInProgress},
		{StateClaimed, StateAwaitingReview},
		{StateInProgress, StateAwaitingReview},
		{StateAwaitingReview, StateApprovedForIntegration},
		{StateAwaitingReview, StateChangesRequested},
		{StateApprovedForIntegration, StateCompleted},
		{StateInProgress, StateReady},
		{StateBlockedHigh, StateReady},
		{StateInProgress, StateClaimed},
		// The pre-CLI bootstrap planning record's shape; see model.go.
		{StateProposed, StateInProgress},
		{StateInProgress, StateCompleted},
	}
	for _, tr := range allowed {
		if !AllowedTransition(tr.From, tr.To) {
			t.Errorf("expected %s -> %s allowed", tr.From, tr.To)
		}
	}
	forbidden := []Transition{
		{StateProposed, StateCompleted},
		{StateReady, StateCompleted},
		{StateCompleted, StateClaimed},
		{StateAwaitingReview, StateInProgress},
	}
	for _, tr := range forbidden {
		if AllowedTransition(tr.From, tr.To) {
			t.Errorf("expected %s -> %s forbidden", tr.From, tr.To)
		}
	}
}

func TestIdentityValidation(t *testing.T) {
	if err := (Identity{SessionCapability: CapLOW, WorkRole: RoleExecutor}).Validate(); err != nil {
		t.Fatalf("valid identity rejected: %v", err)
	}
	if err := (Identity{SessionCapability: "", WorkRole: RoleExecutor}).Validate(); err == nil {
		t.Fatal("missing capability must be rejected")
	}
	if err := (Identity{SessionCapability: CapLOW, WorkRole: "wizard"}).Validate(); err == nil {
		t.Fatal("unknown role must be rejected")
	}
}
