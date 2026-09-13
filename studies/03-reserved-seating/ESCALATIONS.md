# Study 03 — Escalations

Decisions the [Execution Handoff](HANDOFF.md) did not map, raised by an executor session
(Opus 5, `high`) for the planning session (Opus 5, `ultracode`) to decide.

## Protocol

1. **Executor:** append a new `ER-NN` entry using the template below. Never edit an existing
   entry. Commit it with a message starting `Study 03: Escalation Required ER-NN`. Continue
   with unblocked work. If nothing is unblocked, end the session with:
   `Escalation Required: ER-NN — switch to Opus 5, ultracode`.
2. **Ultra session:** append a `Decision` block under the entry. Never rewrite the
   executor's text. If the decision changes the specification, also append an amendment to
   `HANDOFF.md` §16 and tag it `study-03/v0.N-handoff-amendment-NN`. Commit.
3. **Executor, next session:** read every entry with a new decision before resuming.

A decision records the model that took it, and when. Entries stay in this file for good:
they explain why the study looks the way it does.

## Template

```markdown
## ER-NN — <short title>

- raised_by: <model id, setting, tool>
- raised_at: <YYYY-MM-DD HH:MM UTC>
- handoff_sections: <§ numbers>
- trigger: <§12 item number, or "unmapped">
- status: open

**Blocked:** what cannot proceed.

**Context:** what was observed, with file paths, commands and outputs (quote them; link
logs under results/devchecks/).

**Options:**
1. <option> — consequence for what is measured and how results are judged
2. <option> — consequence

**Executor's recommendation:** <option and why, or "none">

**Work continuing meanwhile:** <steps>

### Decision

- decided_by: <model id, setting, tool>
- decided_at: <YYYY-MM-DD HH:MM UTC>
- decision: <option or new instruction>
- rationale:
- handoff_amendment: <AM-NN, or none>
```

---

No escalations yet.
