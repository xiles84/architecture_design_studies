// Package archive reads and writes the durable task archive under
// docs/ai-work/. The committed event files, not the live Git refs, are the
// permanent record: a fresh clone with no live refs must reconstruct identical
// state, which is why every read path here derives state by replaying events.
package archive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"adsqueue/internal/model"
)

// TasksRoot is the repository-relative place task records live.
const TasksRoot = "docs/ai-work/tasks"

// SessionsRoot holds per-session capability records.
const SessionsRoot = "docs/ai-work/sessions"

// TaskDir returns the repository-relative directory for a task id, derived from
// the timestamp embedded in the id: task-YYYYMMDDThhmmssZ-slug.
func TaskDir(taskID string) (string, error) {
	m := regexp.MustCompile(`^task-([0-9]{8})T[0-9]{6}Z-`).FindStringSubmatch(taskID)
	if m == nil {
		return "", fmt.Errorf("cannot derive a directory from task id %q", taskID)
	}
	y, mo := m[1][:4], m[1][4:6]
	return filepath.ToSlash(filepath.Join(TasksRoot, y, mo, taskID)), nil
}

// TaskDirAbs returns the absolute task directory.
func TaskDirAbs(root, taskID string) (string, error) {
	rel, err := TaskDir(taskID)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, filepath.FromSlash(rel)), nil
}

// LoadTask reads task.json and returns the task and its raw bytes.
func LoadTask(dir string) (*model.Task, []byte, error) {
	path := filepath.Join(dir, "task.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var t model.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return &t, data, nil
}

// LoadEvents reads every event file in dir/events, sorted by sequence.
func LoadEvents(dir string) ([]*model.Event, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "events", "*.json"))
	if err != nil {
		return nil, err
	}
	var events []*model.Event
	seen := map[string]string{}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var ev model.Event
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if prev, dup := seen[ev.EventID]; dup {
			return nil, fmt.Errorf("duplicate event_id %s in %s and %s", ev.EventID, prev, p)
		}
		seen[ev.EventID] = p
		events = append(events, &ev)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Sequence != events[j].Sequence {
			return events[i].Sequence < events[j].Sequence
		}
		return events[i].OccurredAt < events[j].OccurredAt
	})
	return events, nil
}

// LoadTaskRecord loads a task directory into a record with derived state.
func LoadTaskRecord(dir string) (*model.TaskRecord, error) {
	t, _, err := LoadTask(dir)
	if err != nil {
		return nil, err
	}
	events, err := LoadEvents(dir)
	if err != nil {
		return nil, err
	}
	return &model.TaskRecord{Task: t, Dir: dir, State: DeriveState(events), Events: events}, nil
}

// LoadAll scans the archive for tasks. A directory with a malformed task.json or
// event chain is returned alongside its error so `audit` can report it without
// hiding the rest of the queue.
func LoadAll(root string) ([]*model.TaskRecord, []error) {
	pattern := filepath.Join(root, filepath.FromSlash(TasksRoot), "*", "*", "*", "task.json")
	paths, _ := filepath.Glob(pattern)
	var (
		records []*model.TaskRecord
		errs    []error
	)
	for _, p := range paths {
		rec, err := LoadTaskRecord(filepath.Dir(p))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, rec)
	}
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Task.TaskID < records[j].Task.TaskID
	})
	return records, errs
}

// DeriveState replays the event chain and returns the resulting state.
func DeriveState(events []*model.Event) string {
	if len(events) == 0 {
		return model.StateProposed
	}
	return events[len(events)-1].ResultingState
}

// ValidateChain checks one task's immutable event history. It is the core of
// `audit` and of the `publish`/transition preconditions. now is passed in so a
// future-clock event can be rejected deterministically in tests.
func ValidateChain(t *model.Task, events []*model.Event, now time.Time) []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	if len(events) == 0 {
		return nil // a proposed task need not have an event yet
	}

	seen := map[string]bool{}
	var prev *model.Event
	state := model.StateProposed
	for i, ev := range events {
		if ev.SchemaVersion != model.SchemaVersion {
			add("event %s: unsupported schema_version %d", ev.EventID, ev.SchemaVersion)
		}
		if seen[ev.EventID] {
			add("duplicate event_id %s", ev.EventID)
		}
		seen[ev.EventID] = true
		if ev.TaskID != t.TaskID {
			add("event %s: task_id %q != %q", ev.EventID, ev.TaskID, t.TaskID)
		}
		if ev.Sequence != i+1 {
			add("event %s: sequence %d, expected %d (gap or duplicate)", ev.EventID, ev.Sequence, i+1)
		}
		if i == 0 {
			if ev.PreviousEventID != nil {
				add("event %s: first event must have a null previous_event_id", ev.EventID)
			}
		} else {
			if ev.PreviousEventID == nil || *ev.PreviousEventID != prev.EventID {
				got := "<nil>"
				if ev.PreviousEventID != nil {
					got = *ev.PreviousEventID
				}
				add("event %s: previous_event_id %s does not name %s (fork or gap)", ev.EventID, got, prev.EventID)
			}
		}
		at, err := model.ParseTime(ev.OccurredAt)
		if err != nil {
			add("event %s: occurred_at: %v", ev.EventID, err)
		} else {
			if at.After(now.Add(5 * time.Minute)) {
				add("event %s: occurred_at %s is in the future", ev.EventID, ev.OccurredAt)
			}
			if prev != nil {
				if pat, perr := model.ParseTime(prev.OccurredAt); perr == nil && at.Before(pat) {
					add("event %s: occurred_at goes backwards", ev.EventID)
				}
			}
		}
		if err := ev.Actor.Validate(); err != nil {
			add("event %s: actor: %v", ev.EventID, err)
		}
		if !validState(ev.ResultingState) {
			add("event %s: unknown resulting_state %q", ev.EventID, ev.ResultingState)
		} else if !model.AllowedTransition(state, ev.ResultingState) {
			add("event %s: impossible transition %s -> %s", ev.EventID, state, ev.ResultingState)
		}
		state = ev.ResultingState
		prev = ev
	}
	return errs
}

func validState(s string) bool {
	switch s {
	case model.StateProposed, model.StateReady, model.StateClaimed, model.StateInProgress,
		model.StateYielded, model.StateBlockedHigh, model.StateAwaitingReview,
		model.StateChangesRequested, model.StateApprovedForIntegration,
		model.StateCompleted, model.StateCancelled, model.StateSuperseded:
		return true
	}
	return false
}

// CompactTime renders a UTC time in the filename form.
func CompactTime(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

var slugClean = regexp.MustCompile(`[^a-z0-9]+`)

// EventID builds a stable, sortable event id: event-<utc>-<type-slug>. A
// counter suffix disambiguates two events written in the same second.
func EventID(eventType string, now time.Time, taken func(string) bool) string {
	slug := slugClean.ReplaceAllString(strings.ToLower(eventType), "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "event"
	}
	base := fmt.Sprintf("event-%s-%s", CompactTime(now), slug)
	id := base
	for i := 2; taken(id); i++ {
		id = base + "-" + strconv.Itoa(i)
	}
	return id
}

// EventFilename is the on-disk name for an event, mirroring the existing
// archive's `<compact-time>-<event-id>.json`.
func EventFilename(ev *model.Event) string {
	at, err := model.ParseTime(ev.OccurredAt)
	ts := strings.ReplaceAll(ev.OccurredAt, ":", "")
	ts = strings.ReplaceAll(ts, "-", "")
	if err == nil {
		ts = CompactTime(at)
	}
	return fmt.Sprintf("%s-%s.json", ts, ev.EventID)
}

// WriteEvent writes an event file and returns its repository-relative-ish path.
func WriteEvent(dir string, ev *model.Event) (string, error) {
	eventsDir := filepath.Join(dir, "events")
	if err := os.MkdirAll(eventsDir, 0o755); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		return "", err
	}
	data = append(data, '\n')
	path := filepath.Join(eventsDir, EventFilename(ev))
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("event file %s already exists", path)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// EventIDs returns the set of existing event ids for a task.
func EventIDs(dir string) (map[string]bool, error) {
	events, err := LoadEvents(dir)
	if err != nil {
		return nil, err
	}
	out := map[string]bool{}
	for _, ev := range events {
		out[ev.EventID] = true
	}
	return out, nil
}

// RenderStatus renders the human-readable STATUS.md snapshot. The immutable
// events remain authoritative; this is a convenience view.
func RenderStatus(rec *model.TaskRecord, now time.Time, nextCap, nextRole, note string) string {
	last := ""
	if ev := rec.LastEvent(); ev != nil {
		last = ev.EventID
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Task status\n\n")
	fmt.Fprintf(&b, "| Field | Value |\n|---|---|\n")
	fmt.Fprintf(&b, "| Task | `%s` |\n", rec.Task.TaskID)
	fmt.Fprintf(&b, "| State | `%s` |\n", rec.State)
	fmt.Fprintf(&b, "| Sequence | `%d` |\n", len(rec.Events))
	fmt.Fprintf(&b, "| Active attempt | %s |\n", activeAttempt(rec))
	fmt.Fprintf(&b, "| Active worker | %s |\n", activeWorker(rec))
	fmt.Fprintf(&b, "| Last event | `%s` |\n", last)
	fmt.Fprintf(&b, "| Updated at | `%s` |\n", model.FormatTime(now))
	fmt.Fprintf(&b, "| Next capability | %s |\n", nextCap)
	fmt.Fprintf(&b, "| Next work role | %s |\n", nextRole)
	fmt.Fprintf(&b, "| Branch | `%s` |\n", rec.Task.CanonicalBranch)
	fmt.Fprintf(&b, "| Worktree | `%s` |\n", rec.Task.CanonicalWorktree)
	if note != "" {
		fmt.Fprintf(&b, "\n%s\n", strings.TrimSpace(note))
	}
	return b.String()
}

func activeAttempt(rec *model.TaskRecord) string {
	ev := rec.LastEvent()
	if ev == nil || ev.AttemptID == nil {
		return "none"
	}
	return "`" + *ev.AttemptID + "`"
}

func activeWorker(rec *model.TaskRecord) string {
	ev := rec.LastEvent()
	if ev == nil {
		return "none"
	}
	return "`" + ev.Actor.Verbose() + "`"
}

// WriteStatus writes STATUS.md.
func WriteStatus(dir string, content string) (string, error) {
	path := filepath.Join(dir, "STATUS.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
