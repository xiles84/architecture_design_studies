package app

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"adsqueue/internal/archive"
	bs "adsqueue/internal/brainstorm"
	"adsqueue/internal/gitx"
	"adsqueue/internal/model"
)

const brainstormLivePrefix = "refs/ads-brainstorms/live/"

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return errors.New("value may not be empty")
	}
	*s = append(*s, v)
	return nil
}

func requireHigh(id model.Identity, command string) error {
	if err := id.Validate(); err != nil {
		return failf(2, "%s: %v", command, err)
	}
	if id.SessionCapability != model.CapHIGH {
		return failf(2, "%s: only a HIGH session may change brainstorm state", command)
	}
	return nil
}

func brainstormRef(id, slot string) string { return brainstormLivePrefix + id + "/" + slot }

func (c *cmdCtx) brainstormMain() (*gitx.Repo, error) { return c.mainWorktree() }

func (c *cmdCtx) lockBrainstormMain(id string, actor model.Identity) (*gitx.Repo, func(), error) {
	main, err := c.brainstormMain()
	if err != nil {
		return nil, nil, err
	}
	oid, err := c.acquireIntegrationLock(id, "", actor)
	if err != nil {
		return nil, nil, err
	}
	released := false
	release := func() {
		if !released {
			_ = c.Repo.DeleteCAS(IntegrationLockRef, oid, "brainstorm archive release "+id)
			released = true
		}
	}
	if st, err := main.StatusPorcelainTracked(); err != nil {
		release()
		return nil, nil, err
	} else if strings.TrimSpace(st) != "" {
		release()
		return nil, nil, failf(1, "brainstorm: local main has tracked edits; refusing archive mutation:\n%s", st)
	}
	return main, release, nil
}

func relPaths(repo *gitx.Repo, paths ...string) ([]string, error) {
	var out []string
	for _, p := range paths {
		r, err := repo.RepoRelative(p)
		if err != nil {
			return nil, err
		}
		out = append(out, filepath.ToSlash(r))
	}
	return out, nil
}

func newBrainstormEvent(r *bs.Record, now time.Time, actor model.Identity, kind, stage, summary string) *bs.Event {
	ev := &bs.Event{
		SchemaVersion:  bs.SchemaVersion,
		EventID:        bs.NewEventID(kind, now, r.Events),
		Sequence:       len(r.Events) + 1,
		BrainstormID:   r.Spec.BrainstormID,
		OccurredAt:     model.FormatTime(now),
		Actor:          actor,
		EventType:      kind,
		ResultingStage: stage,
		Summary:        strings.TrimSpace(summary),
	}
	if last := r.LastEvent(); last != nil {
		id := last.EventID
		ev.PreviousEventID = &id
	}
	return ev
}

func renderBrainstormStatus(r *bs.Record) string {
	var b strings.Builder
	fmt.Fprintf(&b, "BRAINSTORM STATUS\n")
	fmt.Fprintf(&b, "ID: %s\n", r.Spec.BrainstormID)
	fmt.Fprintf(&b, "State: %s\n", r.OverallState())
	fmt.Fprintf(&b, "Stage: %s\n", r.Stage)
	fmt.Fprintf(&b, "Summary:\n  %s\n  Progress: %s\n", oneLine(r.Summary()), r.Progress())
	fmt.Fprintf(&b, "Disagreements: %s\n", oneLine(r.Disagreements()))
	fmt.Fprintf(&b, "Next action: %s\n", r.NextAction())
	return b.String()
}

func renderBrainstormEnded(r *bs.Record) string {
	return fmt.Sprintf("BRAINSTORM ENDED\nID: %s\nConclusion: %s\nRemaining dissent: %s\nConfidence: %s\nTasks: NOT CREATED\n\nTo proceed:\nCREATE TASKS FROM BRAINSTORM %s\n",
		r.Spec.BrainstormID, oneLine(r.Summary()), oneLine(r.Disagreements()), r.Confidence(), r.Spec.BrainstormID)
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func renderBrainstormIndex(title string, records []*bs.Record) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s brainstorms\n\n", title)
	fmt.Fprintf(&b, "Generated from immutable brainstorm records; do not hand-edit.\n")
	if len(records) == 0 {
		fmt.Fprintf(&b, "\n(none)\n")
		return b.String()
	}
	for _, r := range records {
		fmt.Fprintf(&b, "\n## `%s` — %s\n\n", r.Spec.BrainstormID, r.Spec.Title)
		fmt.Fprintf(&b, "- State/stage: `%s` / `%s`\n", r.OverallState(), r.Stage)
		fmt.Fprintf(&b, "- Summary: %s\n", oneLine(r.Summary()))
		fmt.Fprintf(&b, "- Progress: %s\n", r.Progress())
		fmt.Fprintf(&b, "- Disagreements: %s\n", oneLine(r.Disagreements()))
		fmt.Fprintf(&b, "- Next: `%s`\n", r.NextAction())
		if tasks := r.LinkedTasks(); len(tasks) > 0 {
			fmt.Fprintf(&b, "- Linked tasks: `%s`\n", strings.Join(tasks, "`, `"))
		}
	}
	return b.String()
}

func writeBrainstormViews(root string, changed *bs.Record) ([]string, error) {
	statePath := filepath.Join(changed.Dir, "STATE.md")
	if err := os.WriteFile(statePath, []byte("# Current brainstorm state\n\n```text\n"+renderBrainstormStatus(changed)+"```\n"), 0o644); err != nil {
		return nil, err
	}
	records, errs := bs.LoadAll(root)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	var ongoing, concluded []*bs.Record
	for _, r := range records {
		if r.OverallState() == "ongoing" {
			ongoing = append(ongoing, r)
		} else {
			concluded = append(concluded, r)
		}
	}
	dir := filepath.Join(root, "docs", "ai-work", "brainstorms")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	ongoingPath := filepath.Join(dir, "ONGOING.md")
	concludedPath := filepath.Join(dir, "CONCLUDED.md")
	if err := os.WriteFile(ongoingPath, []byte(renderBrainstormIndex("Ongoing", ongoing)), 0o644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(concludedPath, []byte(renderBrainstormIndex("Concluded", concluded)), 0o644); err != nil {
		return nil, err
	}
	return []string{statePath, ongoingPath, concludedPath}, nil
}

func cmdBrainstormStart(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-start")
	var cf commonFlags
	cf.register(fs, true)
	idFlag := fs.String("id", "", "brainstorm id (default: generated)")
	title := fs.String("title", "", "short title (required)")
	question := fs.String("question", "", "question text")
	questionFile := fs.String("question-file", "", "file containing the question")
	positions := fs.Int("positions", 3, "independent position target")
	critiques := fs.Int("critiques", 2, "cross-review target")
	goalID := fs.String("goal", "", "related goal id")
	rootTaskID := fs.String("root-task", "", "related root task id")
	parentTaskID := fs.String("parent-task", "", "related parent task id")
	var criteria, evidence stringList
	fs.Var(&criteria, "criterion", "decision criterion (repeatable)")
	fs.Var(&evidence, "evidence", "repository evidence path (repeatable)")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	actor := cf.identity()
	if err := requireHigh(actor, "brainstorm-start"); err != nil {
		return err
	}
	if strings.TrimSpace(*title) == "" {
		return failf(2, "brainstorm-start: --title is required")
	}
	q := strings.TrimSpace(*question)
	if *questionFile != "" {
		data, err := os.ReadFile(*questionFile)
		if err != nil {
			return err
		}
		q = strings.TrimSpace(string(data))
	}
	if q == "" {
		return failf(2, "brainstorm-start: --question or --question-file is required")
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	id := *idFlag
	if id == "" {
		id = "brainstorm-" + archive.CompactTime(c.Now) + "-" + bs.Slug(*title)
	}
	spec := &bs.Spec{SchemaVersion: bs.SchemaVersion, BrainstormID: id, GoalID: *goalID, RootTaskID: *rootTaskID, ParentTaskID: *parentTaskID,
		Title: strings.TrimSpace(*title), Question: q, DecisionCriteria: criteria, EvidencePaths: evidence,
		PositionTarget: *positions, CritiqueTarget: *critiques, CreatedAt: model.FormatTime(c.Now), CreatedBy: actor}
	if err := spec.Validate(); err != nil {
		return failf(2, "%v", err)
	}
	main, release, err := c.lockBrainstormMain(id, actor)
	if err != nil {
		return err
	}
	defer release()
	dir, _ := bs.Dir(main.Root, id)
	if _, err := os.Stat(dir); err == nil {
		return failf(1, "brainstorm %s already exists", id)
	}
	specPath, err := bs.WriteSpec(dir, spec)
	if err != nil {
		return err
	}
	questionPath, err := bs.WriteQuestion(dir, spec)
	if err != nil {
		return err
	}
	r := &bs.Record{Spec: spec, Dir: dir}
	ev := newBrainstormEvent(r, c.Now, actor, "started", bs.StageCollectingPositions, q)
	eventPath, err := bs.WriteEvent(dir, ev)
	if err != nil {
		return err
	}
	bs.AppendEvent(r, ev)
	views, err := writeBrainstormViews(main.Root, r)
	if err != nil {
		return err
	}
	paths, err := relPaths(main, append([]string{specPath, questionPath, eventPath}, views...)...)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Brainstorm: start %s\n\n%s\n\nAgent: %s\n", id, oneLine(q), actor.Verbose())
	if _, err := main.CommitPaths(msg, paths...); err != nil {
		return err
	}
	release()
	if cf.json {
		return c.jsonOut(map[string]any{"brainstorm_id": id, "state": r.OverallState(), "stage": r.Stage, "next_action": r.NextAction()})
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	return nil
}

type brainstormRow struct {
	BrainstormID  string   `json:"brainstorm_id"`
	Title         string   `json:"title"`
	State         string   `json:"state"`
	Stage         string   `json:"stage"`
	Summary       string   `json:"summary"`
	Progress      string   `json:"progress"`
	Disagreements string   `json:"disagreements"`
	NextAction    string   `json:"next_action"`
	LinkedTasks   []string `json:"linked_tasks,omitempty"`
}

func rowFor(r *bs.Record) brainstormRow {
	return brainstormRow{r.Spec.BrainstormID, r.Spec.Title, r.OverallState(), r.Stage, r.Summary(), r.Progress(), r.Disagreements(), r.NextAction(), r.LinkedTasks()}
}

func cmdBrainstormList(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-list")
	var cf commonFlags
	cf.register(fs, false)
	state := fs.String("state", "", "ongoing, concluded, cancelled or superseded")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	main, err := c.brainstormMain()
	if err != nil {
		return err
	}
	records, errs := bs.LoadAll(main.Root)
	if len(errs) > 0 {
		return errs[0]
	}
	var ongoing, ended []brainstormRow
	for _, r := range records {
		row := rowFor(r)
		if *state != "" && row.State != *state && row.Stage != *state {
			continue
		}
		if row.State == "ongoing" {
			ongoing = append(ongoing, row)
		} else {
			ended = append(ended, row)
		}
	}
	if cf.json {
		return c.jsonOut(map[string]any{"ongoing": ongoing, "concluded": ended})
	}
	printRows := func(title string, rows []brainstormRow) {
		fmt.Fprintf(out, "%s\n\n", title)
		if len(rows) == 0 {
			fmt.Fprintln(out, "(none)")
			return
		}
		for _, r := range rows {
			fmt.Fprintf(out, "%s — %s\nStage: %s\nSummary:\n  %s\n  Progress: %s\nDisagreements: %s\nNext action: %s\n\n",
				r.BrainstormID, r.Title, r.Stage, oneLine(r.Summary), r.Progress, oneLine(r.Disagreements), r.NextAction)
		}
	}
	printRows("ONGOING", ongoing)
	fmt.Fprintln(out)
	printRows("CONCLUDED", ended)
	return nil
}

func cmdBrainstormStatus(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-status")
	var cf commonFlags
	cf.register(fs, false)
	id := fs.String("brainstorm", "", "brainstorm id")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	main, err := c.brainstormMain()
	if err != nil {
		return err
	}
	r, err := bs.Load(main.Root, *id)
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(rowFor(r))
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	if r.Stage == bs.StageConcluded {
		fmt.Fprintln(out)
		fmt.Fprint(out, renderBrainstormEnded(r))
	}
	return nil
}

func readBrainstormClaim(repo *gitx.Repo, id, slot string) (*bs.Claim, string, error) {
	oid, ok := repo.RefOID(brainstormRef(id, slot))
	if !ok {
		return nil, "", nil
	}
	data, err := repo.ReadBlob(oid)
	if err != nil {
		return nil, oid, err
	}
	var cl bs.Claim
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, oid, err
	}
	if err := cl.Validate(); err != nil {
		return nil, oid, err
	}
	return &cl, oid, nil
}

func writeBrainstormClaim(repo *gitx.Repo, cl *bs.Claim, oldOID string) (string, error) {
	data, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return "", err
	}
	oid, err := repo.WriteBlob(append(data, '\n'))
	if err != nil {
		return "", err
	}
	if err := repo.CAS(brainstormRef(cl.BrainstormID, cl.SlotID), oid, oldOID, "brainstorm claim "+cl.BrainstormID+" "+cl.SlotID); err != nil {
		return "", err
	}
	return oid, nil
}

func brainstormClaimContextPath(repo *gitx.Repo, id string) string {
	return filepath.Join(repo.GitDir, "ads-brainstorms", "claims", id+".json")
}

func saveBrainstormClaimContext(repo *gitx.Repo, cl *bs.Claim) error {
	path := brainstormClaimContextPath(repo, cl.BrainstormID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(cl, "", "  ")
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func loadBrainstormClaimContext(repo *gitx.Repo, id string) (*bs.Claim, error) {
	data, err := os.ReadFile(brainstormClaimContextPath(repo, id))
	if err != nil {
		return nil, err
	}
	var cl bs.Claim
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, err
	}
	return &cl, nil
}

func cmdBrainstormClaim(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-claim")
	var cf commonFlags
	cf.register(fs, true)
	id := fs.String("brainstorm", "", "brainstorm id")
	leaseText := fs.String("lease", "2h", "claim lease")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	actor := cf.identity()
	if err := requireHigh(actor, "brainstorm-claim"); err != nil {
		return err
	}
	lease, err := model.ParseDuration(*leaseText)
	if err != nil {
		return failf(2, "brainstorm-claim: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	main, err := c.brainstormMain()
	if err != nil {
		return err
	}
	r, err := bs.Load(main.Root, *id)
	if err != nil {
		return err
	}
	if r.OverallState() != "ongoing" {
		return failf(1, "brainstorm %s is %s, not claimable", *id, r.Stage)
	}
	for _, slot := range r.NextSlots() {
		old, oid, err := readBrainstormClaim(c.Repo, *id, slot.ID)
		if err != nil {
			return err
		}
		epoch := 1
		if old != nil {
			exp, xerr := model.ParseTime(old.ExpiresAt)
			if xerr == nil && exp.After(c.Now) {
				continue
			}
			epoch = old.ClaimEpoch + 1
		}
		cl := &bs.Claim{SchemaVersion: bs.SchemaVersion, BrainstormID: *id, SlotID: slot.ID, SlotKind: slot.Kind,
			ClaimID: "brainstorm-claim-" + randHex(8), ClaimEpoch: epoch, Worker: actor,
			ClaimedAt: model.FormatTime(c.Now), HeartbeatAt: model.FormatTime(c.Now), ExpiresAt: model.FormatTime(c.Now.Add(lease)), Stage: r.Stage}
		if _, err := writeBrainstormClaim(c.Repo, cl, oid); err != nil {
			if errors.Is(err, gitx.ErrCASConflict) {
				continue
			}
			return err
		}
		if err := saveBrainstormClaimContext(c.Repo, cl); err != nil {
			return err
		}
		if cf.json {
			return c.jsonOut(cl)
		}
		fmt.Fprintf(out, "claimed brainstorm %s slot %s (%s)\nclaim_id: %s\nepoch: %d\nexpires: %s\n", *id, slot.ID, slot.Kind, cl.ClaimID, cl.ClaimEpoch, cl.ExpiresAt)
		if slot.Kind == bs.SlotPosition {
			fmt.Fprintf(out, "Read: %s and its evidence packet only. Do not read positions/ before this round closes.\n", filepath.ToSlash(filepath.Join(r.Dir, "QUESTION.md")))
		} else if slot.Kind == bs.SlotCritique {
			fmt.Fprintln(out, "Read: QUESTION.md and every file under positions/.")
		} else {
			fmt.Fprintln(out, "Read: QUESTION.md, positions/, and critiques/; preserve dissent in the conclusion.")
		}
		return nil
	}
	return failf(1, "brainstorm %s has no unclaimed slot ready in stage %s", *id, r.Stage)
}

func resolveBrainstormClaim(c *cmdCtx, id, slot, claimID string, epoch int) (*bs.Claim, string, error) {
	if slot == "" || claimID == "" || epoch == 0 {
		ctx, err := loadBrainstormClaimContext(c.Repo, id)
		if err != nil {
			return nil, "", failf(3, "no brainstorm claim context for %s", id)
		}
		if slot == "" {
			slot = ctx.SlotID
		}
		if claimID == "" {
			claimID = ctx.ClaimID
		}
		if epoch == 0 {
			epoch = ctx.ClaimEpoch
		}
	}
	cl, oid, err := readBrainstormClaim(c.Repo, id, slot)
	if err != nil {
		return nil, "", err
	}
	if cl == nil || oid == "" {
		return nil, "", failf(3, "no live claim for %s slot %s", id, slot)
	}
	if cl.ClaimID != claimID || cl.ClaimEpoch != epoch {
		return nil, "", failf(3, "brainstorm claim superseded: live %s epoch %d, caller %s epoch %d", cl.ClaimID, cl.ClaimEpoch, claimID, epoch)
	}
	exp, err := model.ParseTime(cl.ExpiresAt)
	if err != nil || !exp.After(c.Now) {
		return nil, "", failf(3, "brainstorm claim expired at %s", cl.ExpiresAt)
	}
	return cl, oid, nil
}

func brainstormClaimFlags(fs *flag.FlagSet) (id, slot, claimID *string, epoch *int) {
	id = fs.String("brainstorm", "", "brainstorm id")
	slot = fs.String("slot", "", "slot id (default: saved claim context)")
	claimID = fs.String("claim-id", "", "claim id (default: saved claim context)")
	epoch = fs.Int("claim-epoch", 0, "claim epoch (default: saved claim context)")
	return
}

func cmdBrainstormGuard(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-guard")
	var cf commonFlags
	cf.register(fs, false)
	id, slot, claimID, epoch := brainstormClaimFlags(fs)
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, _, err := resolveBrainstormClaim(c, *id, *slot, *claimID, *epoch)
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(cl)
	}
	fmt.Fprintf(out, "brainstorm guard ok: %s %s claim %s epoch %d expires %s\n", cl.BrainstormID, cl.SlotID, cl.ClaimID, cl.ClaimEpoch, cl.ExpiresAt)
	return nil
}

func cmdBrainstormHeartbeat(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-heartbeat")
	var cf commonFlags
	cf.register(fs, false)
	id, slot, claimID, epoch := brainstormClaimFlags(fs)
	leaseText := fs.String("lease", "2h", "lease extension")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	lease, err := model.ParseDuration(*leaseText)
	if err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := resolveBrainstormClaim(c, *id, *slot, *claimID, *epoch)
	if err != nil {
		return err
	}
	cl.HeartbeatAt = model.FormatTime(c.Now)
	cl.ExpiresAt = model.FormatTime(c.Now.Add(lease))
	if _, err := writeBrainstormClaim(c.Repo, cl, oid); err != nil {
		return err
	}
	_ = saveBrainstormClaimContext(c.Repo, cl)
	if cf.json {
		return c.jsonOut(cl)
	}
	fmt.Fprintf(out, "brainstorm heartbeat %s %s -> %s\n", cl.BrainstormID, cl.SlotID, cl.ExpiresAt)
	return nil
}

func cmdBrainstormSubmit(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-submit")
	var cf commonFlags
	cf.register(fs, false)
	id, slot, claimID, epoch := brainstormClaimFlags(fs)
	file := fs.String("file", "", "contribution markdown file")
	summary := fs.String("summary", "", "two-to-four-line state summary")
	disagreements := fs.String("disagreements", "none", "remaining disagreements")
	confidence := fs.String("confidence", "medium", "high, medium or low")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *file == "" || strings.TrimSpace(*summary) == "" {
		return failf(2, "brainstorm-submit: --file and --summary are required")
	}
	if *confidence != "high" && *confidence != "medium" && *confidence != "low" {
		return failf(2, "brainstorm-submit: --confidence must be high, medium or low")
	}
	content, err := os.ReadFile(*file)
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(content)) == "" {
		return failf(2, "brainstorm-submit: contribution file may not be empty")
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := resolveBrainstormClaim(c, *id, *slot, *claimID, *epoch)
	if err != nil {
		return err
	}
	main, release, err := c.lockBrainstormMain(*id, cl.Worker)
	if err != nil {
		return err
	}
	defer release()
	// Revalidate under the archive lock, then extend the claim before writing.
	// This closes the window in which a lease could expire and be recovered
	// after the first guard but before the contribution commit.
	cl, oid, err = resolveBrainstormClaim(c, *id, cl.SlotID, cl.ClaimID, cl.ClaimEpoch)
	if err != nil {
		return err
	}
	cl.HeartbeatAt = model.FormatTime(c.Now)
	cl.ExpiresAt = model.FormatTime(c.Now.Add(model.DefaultLease))
	oid, err = writeBrainstormClaim(c.Repo, cl, oid)
	if err != nil {
		return err
	}
	r, err := bs.Load(main.Root, *id)
	if err != nil {
		return err
	}
	if r.SubmittedSlots()[cl.SlotID] {
		return failf(1, "slot %s was already submitted", cl.SlotID)
	}
	stage, err := bs.StageAfterSubmission(r, cl.SlotKind)
	if err != nil {
		return err
	}
	eventType := cl.SlotKind + "_submitted"
	ev := newBrainstormEvent(r, c.Now, cl.Worker, eventType, stage, *summary)
	ev.SlotID, ev.SlotKind = cl.SlotID, cl.SlotKind
	ev.ContributionID = "contribution-" + archive.CompactTime(c.Now) + "-" + cl.SlotID + "-" + randHex(3)
	ev.Disagreements, ev.Confidence = strings.TrimSpace(*disagreements), *confidence
	contributionPath, err := bs.WriteContribution(r, ev, content)
	if err != nil {
		return err
	}
	eventPath, err := bs.WriteEvent(r.Dir, ev)
	if err != nil {
		return err
	}
	bs.AppendEvent(r, ev)
	views, err := writeBrainstormViews(main.Root, r)
	if err != nil {
		return err
	}
	paths, err := relPaths(main, append([]string{contributionPath, eventPath}, views...)...)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Brainstorm: %s %s (%s)\n\n%s\n\nAgent: %s\n", *id, cl.SlotID, eventType, oneLine(*summary), cl.Worker.Verbose())
	if _, err := main.CommitPaths(msg, paths...); err != nil {
		return err
	}
	if err := c.Repo.DeleteCAS(brainstormRef(*id, cl.SlotID), oid, "brainstorm submit "+*id+" "+cl.SlotID); err != nil {
		return err
	}
	_ = os.Remove(brainstormClaimContextPath(c.Repo, *id))
	release()
	if cf.json {
		return c.jsonOut(map[string]any{"brainstorm_id": *id, "contribution_id": ev.ContributionID, "state": r.OverallState(), "stage": r.Stage, "next_action": r.NextAction()})
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	if r.Stage == bs.StageConcluded {
		fmt.Fprintln(out)
		fmt.Fprint(out, renderBrainstormEnded(r))
	}
	return nil
}

func activeBrainstormRefs(c *cmdCtx, id string) ([]string, error) {
	out, err := c.Repo.Run("for-each-ref", "--format=%(refname)", brainstormLivePrefix+id+"/")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(out) == "" {
		return nil, nil
	}
	return strings.Split(strings.TrimSpace(out), "\n"), nil
}

func rejectActiveAndClearExpiredBrainstormClaims(c *cmdCtx, id string) error {
	refs, err := activeBrainstormRefs(c, id)
	if err != nil {
		return err
	}
	for _, ref := range refs {
		slot := strings.TrimPrefix(ref, brainstormLivePrefix+id+"/")
		cl, oid, err := readBrainstormClaim(c.Repo, id, slot)
		if err != nil {
			return err
		}
		expires, err := model.ParseTime(cl.ExpiresAt)
		if err != nil {
			return err
		}
		if expires.After(c.Now) {
			return failf(1, "brainstorm %s has active contribution claim %s until %s", id, slot, cl.ExpiresAt)
		}
		if err := c.Repo.DeleteCAS(ref, oid, "brainstorm clear expired claim "+id+" "+slot); err != nil {
			return err
		}
	}
	return nil
}

func commitAdministrativeEvent(c *cmdCtx, id string, actor model.Identity, kind, summary, replacement string, taskIDs []string) (*bs.Record, error) {
	main, release, err := c.lockBrainstormMain(id, actor)
	if err != nil {
		return nil, err
	}
	defer release()
	if err := rejectActiveAndClearExpiredBrainstormClaims(c, id); err != nil {
		return nil, err
	}
	r, err := bs.Load(main.Root, id)
	if err != nil {
		return nil, err
	}
	stage := r.Stage
	switch kind {
	case "cancelled":
		if r.OverallState() != "ongoing" {
			return nil, failf(1, "only an ongoing brainstorm can be cancelled")
		}
		stage = bs.StageCancelled
	case "superseded":
		if r.Stage != bs.StageConcluded && r.Stage != bs.StageTasked {
			return nil, failf(1, "only a concluded or tasked brainstorm can be superseded")
		}
		if _, err := bs.Load(main.Root, replacement); err != nil {
			return nil, fmt.Errorf("replacement brainstorm: %w", err)
		}
		stage = bs.StageSuperseded
	case "tasks_linked":
		if r.Stage != bs.StageConcluded && r.Stage != bs.StageTasked {
			return nil, failf(1, "tasks may be linked only after conclusion")
		}
		for _, taskID := range taskIDs {
			dir, err := archive.TaskDirAbs(main.Root, taskID)
			if err != nil {
				return nil, err
			}
			t, _, err := archive.LoadTask(dir)
			if err != nil {
				return nil, err
			}
			if t.OriginatingBrainstormID != id {
				return nil, failf(1, "task %s does not declare originating_brainstorm_id %s", taskID, id)
			}
		}
		stage = bs.StageTasked
	}
	ev := newBrainstormEvent(r, c.Now, actor, kind, stage, summary)
	ev.SupersededBy, ev.LinkedTaskIDs = replacement, taskIDs
	eventPath, err := bs.WriteEvent(r.Dir, ev)
	if err != nil {
		return nil, err
	}
	bs.AppendEvent(r, ev)
	views, err := writeBrainstormViews(main.Root, r)
	if err != nil {
		return nil, err
	}
	paths, err := relPaths(main, append([]string{eventPath}, views...)...)
	if err != nil {
		return nil, err
	}
	if _, err := main.CommitPaths(fmt.Sprintf("Brainstorm: %s %s\n\n%s\n\nAgent: %s\n", id, kind, summary, actor.Verbose()), paths...); err != nil {
		return nil, err
	}
	release()
	return r, nil
}

func cmdBrainstormCancel(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-cancel")
	var cf commonFlags
	cf.register(fs, true)
	id := fs.String("brainstorm", "", "brainstorm id")
	reason := fs.String("reason", "Cancelled by owner direction.", "reason")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	actor := cf.identity()
	if err := requireHigh(actor, "brainstorm-cancel"); err != nil {
		return err
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	r, err := commitAdministrativeEvent(c, *id, actor, "cancelled", *reason, "", nil)
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(rowFor(r))
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	return nil
}

func cmdBrainstormSupersede(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-supersede")
	var cf commonFlags
	cf.register(fs, true)
	id := fs.String("brainstorm", "", "old brainstorm id")
	replacement := fs.String("with", "", "replacement brainstorm id")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	actor := cf.identity()
	if err := requireHigh(actor, "brainstorm-supersede"); err != nil {
		return err
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	r, err := commitAdministrativeEvent(c, *id, actor, "superseded", "Superseded by "+*replacement+".", *replacement, nil)
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(rowFor(r))
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	return nil
}

func cmdBrainstormLinkTasks(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-link-tasks")
	var cf commonFlags
	cf.register(fs, true)
	id := fs.String("brainstorm", "", "brainstorm id")
	var tasks stringList
	fs.Var(&tasks, "task", "published task id (repeatable)")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if len(tasks) == 0 {
		return failf(2, "brainstorm-link-tasks: at least one --task is required")
	}
	actor := cf.identity()
	if err := requireHigh(actor, "brainstorm-link-tasks"); err != nil {
		return err
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	r, err := commitAdministrativeEvent(c, *id, actor, "tasks_linked", "Published and linked implementation tasks: "+strings.Join(tasks, ", ")+".", "", tasks)
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(rowFor(r))
	}
	fmt.Fprint(out, renderBrainstormStatus(r))
	return nil
}

func cmdBrainstormAudit(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-audit")
	var cf commonFlags
	cf.register(fs, false)
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	main, err := c.brainstormMain()
	if err != nil {
		return err
	}
	records, loadErrs := bs.LoadAll(main.Root)
	type issue struct {
		BrainstormID string `json:"brainstorm_id,omitempty"`
		Detail       string `json:"detail"`
	}
	var issues []issue
	for _, err := range loadErrs {
		issues = append(issues, issue{Detail: err.Error()})
	}
	for _, r := range records {
		for _, err := range bs.ValidateChain(r.Spec, r.Events, c.Now) {
			issues = append(issues, issue{r.Spec.BrainstormID, err.Error()})
		}
		for _, ev := range r.Events {
			if ev.EventType == "tasks_linked" {
				for _, taskID := range ev.LinkedTaskIDs {
					dir, err := archive.TaskDirAbs(main.Root, taskID)
					if err != nil {
						issues = append(issues, issue{r.Spec.BrainstormID, "invalid linked task " + taskID + ": " + err.Error()})
						continue
					}
					task, _, err := archive.LoadTask(dir)
					if err != nil {
						issues = append(issues, issue{r.Spec.BrainstormID, "missing linked task " + taskID})
						continue
					}
					if task.OriginatingBrainstormID != r.Spec.BrainstormID {
						issues = append(issues, issue{r.Spec.BrainstormID, "linked task " + taskID + " has a different brainstorm origin"})
					}
				}
			}
			if ev.ContributionID == "" {
				continue
			}
			var path string
			switch ev.SlotKind {
			case bs.SlotPosition:
				path = filepath.Join(r.Dir, "positions", ev.ContributionID+".md")
			case bs.SlotCritique:
				path = filepath.Join(r.Dir, "critiques", ev.ContributionID+".md")
			case bs.SlotSynthesis:
				path = filepath.Join(r.Dir, "CONCLUSION.md")
			}
			if _, err := os.Stat(path); err != nil {
				issues = append(issues, issue{r.Spec.BrainstormID, "missing contribution " + path})
			}
		}
	}
	byID := make(map[string]*bs.Record, len(records))
	for _, r := range records {
		byID[r.Spec.BrainstormID] = r
	}
	refs, err := c.Repo.Run("for-each-ref", "--format=%(refname)", brainstormLivePrefix)
	if err != nil {
		return err
	}
	for _, ref := range strings.Fields(refs) {
		rest := strings.TrimPrefix(ref, brainstormLivePrefix)
		parts := strings.Split(rest, "/")
		if len(parts) != 2 {
			issues = append(issues, issue{Detail: "malformed live ref " + ref})
			continue
		}
		cl, _, err := readBrainstormClaim(c.Repo, parts[0], parts[1])
		if err != nil {
			issues = append(issues, issue{BrainstormID: parts[0], Detail: "invalid live claim: " + err.Error()})
			continue
		}
		r := byID[parts[0]]
		switch {
		case cl == nil:
			issues = append(issues, issue{BrainstormID: parts[0], Detail: "live ref has no claim"})
		case cl.BrainstormID != parts[0] || cl.SlotID != parts[1]:
			issues = append(issues, issue{BrainstormID: parts[0], Detail: "live claim does not match ref path"})
		case r == nil:
			issues = append(issues, issue{BrainstormID: parts[0], Detail: "orphan live claim"})
		case r.SubmittedSlots()[parts[1]]:
			issues = append(issues, issue{BrainstormID: parts[0], Detail: "live claim names an already submitted slot " + parts[1]})
		case r.Stage != cl.Stage:
			issues = append(issues, issue{BrainstormID: parts[0], Detail: fmt.Sprintf("live claim stage %s differs from archive stage %s", cl.Stage, r.Stage)})
		}
	}
	if cf.json {
		return c.jsonOut(map[string]any{"brainstorms": len(records), "errors": len(issues), "issues": issues})
	}
	fmt.Fprintf(out, "brainstorm audit: %d records, %d errors\n", len(records), len(issues))
	for _, i := range issues {
		fmt.Fprintf(out, "  [ERROR] %s %s\n", i.BrainstormID, i.Detail)
	}
	if len(issues) > 0 {
		return failf(1, "brainstorm audit found %d error(s)", len(issues))
	}
	return nil
}

func cmdBrainstormIntent(args []string, out io.Writer) error {
	fs := newFlagSet("brainstorm-intent")
	var cf commonFlags
	cf.register(fs, false)
	text := fs.String("text", "", "user message")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	intent := bs.ParseIntent(*text)
	if cf.json {
		c, err := cf.open(out)
		if err != nil {
			return err
		}
		return c.jsonOut(intent)
	}
	if intent.Kind == "" {
		fmt.Fprintln(out, "(no brainstorm intent)")
		return nil
	}
	fmt.Fprintf(out, "%s\t%s\t%s\n", intent.Kind, intent.BrainstormID, intent.Question)
	return nil
}
