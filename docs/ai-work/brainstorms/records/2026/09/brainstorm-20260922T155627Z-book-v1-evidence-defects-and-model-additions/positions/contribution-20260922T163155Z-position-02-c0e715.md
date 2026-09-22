# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution | `contribution-20260922T163155Z-position-02-c0e715` |
| Slot | `position-02` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=unknown capability=HIGH role=analyst |
| Capability input | leader |
| Submitted | `2026-09-22T16:31:55Z` |
| Confidence | high |

## Summary

The tag/commit half of the evidence contract is unenforced: 11 of 13 numeric claims omit a run tag that exists in git and tools/evidence still reports 0 errors, while claim-01 has no commit anywhere in its chain (manifest, report and analysis alike) and must be published with that gap labelled rather than re-run. Four measured families sit outside the book's six - arbitration of a contended resource (study 02, nine arbiters, two controls, 4-8x), expiry and clock authority (study 03, six designs), unit of change / aggregate boundary, and derived current value (study 01 D18-D24) - plus three decision models (cost of maintained derived state, how to read these numbers, questions a schema cannot answer). All of it needs writing, not measuring.

## Disagreements

Anticipates three objections: that optional run_tag is deliberate because study 01 predates the tagging rule (answer: record the reason, do not re-run, since a null meaning two different things is the defect); that naming families pre-empts book-synthesis-v1 (answer: the synthesis owns the taxonomy, this position only names which measured designs lack a home); and that claim-07's disagreement between trials is a harness fault rather than a reading rule (answer: the analysis publishes it as a limit, so it is evidence for the reading chapter).

## Contribution

# Independent position — evidence corrections and model additions for book v1

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions` |
| Contribution / slot | assigned by the CLI on submission · `position-02` |
| Actor | deepseek-flash, WSL2 shell + `tools/queue`, effort not exposed by the client; canonical capability `HIGH` (raw input `leader`), role `analyst` |
| Submitted | UTC timestamp recorded by the CLI at submission |
| Confidence | **high** on the provenance findings (each is reproduced by a command below); **medium** on the taxonomy proposals |

## Evidence examined

Everything below was read or run on 2026-09-22 from the main checkout, no database started:

- `book/evidence/claims.json`, `book/evidence/claims.schema.json`, `book/evidence/README.md`.
- `tools/evidence`: `go run . validate --repo ../..` → `16 claims, 0 errors, 16 book inputs`;
  `--list-inputs` → the sixteen ids. The validator's own documented rules (README §What
  `validate` enforces) are quoted where they matter.
- `git tag -l 'run/*'` (17 tags) compared, claim by claim, with each claim's
  `provenance.run_id`; and `git rev-parse -q --verify refs/tags/<tag>` per comparison.
- All 71 `studies/*/results/**/manifest.yaml`, extracting `repo_commit`, `repo_dirty` and
  `run_tag` from every file.
- The five study `README.md` design tables and controlled-pair tables; `CONTEXT.md`;
  `LESSONS_LEARNED.md` section headings; the book goal's `GOAL.md` and `BACKLOG.md`;
  `docs/methodology.md`; `docs/environments/host-zenbook-ux5406sa.md`.
- `studies/01-charity-tree/results/20260912-small/manifest.yaml`,
  `studies/01-charity-tree/reports/20260912-small.md`,
  `studies/01-charity-tree/reports/analyses/20260912-study01--claude-opus-5--2026-09-12.md`.
- `studies/02-ticket-booking/reports/analyses/20260913T021206Z--claude-opus-5--2026-09-13.md`
  (the line that excludes the aborted run).
- `studies/01-charity-tree/reports/analyses/20260913T125342Z-v3--gpt-6--2026-09-13.md` front-matter.
- Environment check: `podman info` → 8 CPUs / 16,496,418,816 B, matching the environment page
  exactly; `podman version --format` → client 5.8.1, server 5.8.5.

## Recommendation

First, the two things I would not publish v1 without.

**(a1) Treat the tag/commit half of the evidence contract as unenforced, and fix it.**
`book/evidence/claims.schema.json` requires only
`[run_id, environment, topology, inputs_digest, report, analysis]` in `provenance`;
`run_tag` and `repo_commit` are typed `["string","null"]` and are **not required**. The
validator checks a tag *if given* ("a `run_tag`, if given, exists") and never that an existing
tag *is* given. Consequences I verified:

- **11 of the 13 numeric claims record `run_tag: null` although `run/<study>/<run_id>` exists
  in this repository right now** — claims 02, 03, 04, 05, 06, 07, 08, 09, 10, `gap-01`,
  `gap-03`, `gap-04`. Only claims 11, 12 and `gap-02` name a tag.
- `validate` still exits `0 errors`. A green registry therefore says nothing about whether a
  claim can be resolved to a tag, which is the property the book's promise rests on.

Correction, entirely inside the registry and the validator: make `run_tag` a **required
resolution** rather than an optional field — for every non-`gap` claim, `run_tag` must name a
tag that exists, or the claim must carry an explicit reason it cannot (see a2). This touches
no report, analysis, discussion, run or tag.

**(a2) `claim-01-indexes-and-rollups` has no commit anywhere in its chain.**
`studies/01-charity-tree/results/20260912-small/manifest.yaml` records run id, environment,
scale, durations, connections, images, budgets and podman version, but **no `repo_commit`, no
`repo_dirty`, no `run_tag`**; no tag `run/01-charity-tree/20260912-small` exists; the report
adds none; and the analysis front-matter has **no `repo_commit:` field at all** (unlike every
later analysis, e.g. the v3 one). The same is true of the other 18 pre-v3 study-01 run
directories. So the book's first family claim — stored parent aggregates, indexes, and
"embedding never came out ahead on reads" — rests on a run that cannot be tied to any commit,
while `claim-03`, which asserts the same family's read advantage, carries five fresh loads, a
commit (`c084c73b…`) and a tag.

I am **not** proposing a re-run: that is a measurement, out of scope for v1. I am proposing
that (i) the book states claim-01's provenance gap next to the claim, and (ii) the registry
distinguishes *"resolved to a tag"* from *"untagged, reason recorded"* so the two are never
rendered identically. This is the cheapest place in the corpus where a reader's trust is
currently unearned.

Everything else, strongest first:

- **(a3) `run/*` tags are not a results index.** `run/02-ticket-booking/20260921T182546Z`
  points at `facb2c0` ("Context: phase 3b resumed by a new session") and there is **no**
  `studies/02-ticket-booking/results/20260921T182546Z/`. The book must not derive "the runs"
  from `git tag -l 'run/*'`.
- **(a4) The manifest's `run_tag` field is unreliable in both directions.** It asserts
  `run/02-ticket-booking/20260913T021010Z`, which does not exist (that run's commit is also
  `unknown` — the known path-conversion failure); and `results/20260921T-survey/manifest.yaml`
  says `tag=none` while `run/05-cache-consistency/20260921T-survey` exists. Any tag resolution
  must go through git, never through the manifest field. The aborted run itself is handled
  correctly: the study-02 analysis calls it "aborted, INVALID" and excludes it from every
  conclusion, so no published number is affected.
- **(a5) Dirty-tree runs exist and are correctly quarantined.** `03/20260916T000706Z` and
  `05/20260921T-survey2` are `repo_dirty: true` with no tag, and no claim rests on either.
  The book should say that the pipeline produces and quarantines dirty runs, rather than let a
  reader assume every committed result came from a clean tree. Study 01's pre-v3 runs cannot
  even be audited for this, because the field is absent (a2).
- **(a6) Three claims are validation-open.** `claim-05`, `claim-06` and `claim-07` (all
  phase 3b) carry `claim-gap-04`. Their measured status must be rendered inline in the book,
  not as a footnote.

**(b) Models the corpus can already support, strongest first.** The book's six families are
normalized facts, rolldown/duplication, rollups/materialized aggregates, embedded documents,
append-only history/snapshots, and external read copies/caches. Four measured families in the
corpus sit outside all six, and three decision models are supported by committed evidence:

- **(b1) Arbitration of a contended resource — the strongest addition.** Study 02 measures
  nine arbiters for one scarcer-than-demand resource, each isolated by a controlled pair:
  locking the lowest free row (P1), `SKIP LOCKED` (P2), compare-and-set (P3), a guarded counter
  on the parent row (C4/P4), moving the counter off the hot row (R1), sharding it (R2),
  a unique-index arbiter (C5), a pre-created narrow slot pool (R3), and a serializable
  read-then-write (C2) — with two negative controls that must overbook (C1, H0), and a measured
  4–8× spread between the arbiters. This is a design family in its own right, not a variant of
  "normalized facts", and its cost is a correctness cost as much as a throughput cost.
- **(b2) Expiry, lease and clock authority.** Study 03 varies who makes a deadline take effect
  and whose clock judges it: lazy check at use (S1), a sweeper (E1), the deadline stored once
  on a cart row instead of on every seat (E2), a bounded payment window that extends the lease
  (K1), a single retry at confirmation (S1r) — against a negative control that judges expiry by
  the application clock under injected skew (E0) and one that confirms without checking the hold
  (K0). Six designs, one measured mechanism, and the book's own brief already asks for an
  "expiry and clock authority" chapter.
- **(b3) Unit of change / aggregate boundary.** What one transaction replaces: a key (n1), a
  whole document row (04 `d1_doc_row`), a section document under version compare-and-set
  (03 `l2_section_document`), or a bounded embedded slice (01 D9/D10 vs D6). Study 04's
  `d1`↔`n1` pair is a clean one-decision contrast, and study 05's owned model adds the
  transactional outbox to the same boundary question.
- **(b4) Derived current value (recency), with its own write path.** Study 01 D18–D24 maintain
  "the donor's latest donation" four ways — a flag on the child plus a partial index with a
  parent-row lock (D20), the same index on an already-stored rollup column (D22/D23), a
  parent-driven probe (D18), and a window-first check of the same candidates (D19) — with a
  negative control that removes the lock (D21), for a measured 2.9× (PostgreSQL) and 7.9×
  (YugabyteDB) gain. This is not a rollup: it derives a *latest* value and needs a lock order
  the sum-rollups do not.
- **(b5) Minor models, cheap to state:** the cost of enforcing referential integrity
  (01 D3→D8 — I verified the pair differs *only* in the three `REFERENCES` clauses once
  comments are stripped, so it is a clean one-decision contrast) and the covering index as the
  access path for a derived aggregate (01 D14/D15 with `INCLUDE`, and D22/D23).
- **(b6) Decision model — what maintaining derived state costs.** One model spanning three
  studies: parent aggregate (01 D4/D5), hot counter (02 C4), off-row counter (R1), sharded
  counter (R2), trigger vs application maintenance (01 D4/D5, 04 n3/n4), and the answerability
  price of an append-only ledger (02 X1: 35 % storage, 8× load time, and a measured question
  nothing else can answer).
- **(b7) Decision model — how to read these numbers.** The corpus supplies its own error bars:
  byte-identical-SQL pairs that measure the harness noise floor (01 D4/D5, 02 C1/C2/C3,
  03 q05/q06 through twelve designs), five-load repetition (claim-03), trial spread that
  reverses a direction (claim-07: P3's PostgreSQL trials disagree by 92–283 %), and controls
  whose failure invalidates the run. This is the book's most transferable chapter and needs no
  new measurement.
- **(b8) Decision model — the questions a schema cannot answer.** Two measured instances:
  refund reporting absent in 14 of 15 study-02 designs (`r05`), and the hold-window sizing
  report unanswerable in every study-03 design. Both are committed as `gap` claims.

## Reasoning and trade-offs

The provenance items are worth doing first because they are the only findings here that make
the book **wrong** rather than incomplete: a reader who checks a claim against a tag, or who
assumes `validate` proves traceability, is misled by a repository that currently looks green.
They are also the cheapest — registry data, one schema requirement, and one validator rule —
and they touch no signed artefact.

The model additions are the largest quality gain available without measuring. Four of them
(b1–b4) already have measured designs, controlled pairs and negative controls; naming them
costs book writing, not benchmarking. I deliberately put b1 first: it has the most designs, the
cleanest controls, and the largest measured effect, and it is the only one of the four whose
absence from the book's taxonomy would leave a reader without any model for the most common
production question in the corpus ("who wins the last one?").

The cost of this direction is honesty-maintenance: promoting a family means every claim inside
it inherits that family's weakest evidence (single run, one topology, laptop), so the taxes in
(a1)–(a6) come due at the same time.

## Assumptions and evidence gaps

- I did **not** recompute any `inputs_digest`; I verified digest presence in the report and
  analysis, not the hash over result bytes.
- I read the study READMEs and the front-matter/headline of the analyses, not all 56 reports.
  A defect that only appears inside a report body may be missing from this position.
- Design counts and controlled pairs come from the READMEs; I did not re-verify each pair
  against `harness/designs.go`.
- The map from designs to the book's six families is my reading of the book brief, not a
  committed taxonomy.
- I inferred nothing from the `podman machine list` line: the guest was measured with
  `podman info` and matches the environment page, so there is no environment defect here.

## Alternatives

- *Re-run study 01's matrix to replace claim-01 with a tagged run.* Rejected: it is a new
  measurement, it is not authorized by this brainstorm, and it would not help the book's
  honesty in v1 — a label does that today, for free.
- *Leave the registry as it is and resolve tags at book-build time only.* Rejected: the same
  resolver would then silently mis-handle `claim-01`, which has no tag at all, and the
  registry is the artefact the book's claims index is generated from.
- *Treat the four candidate families as variants inside the existing six.* Rejected for b1
  (arbitration is orthogonal to every one of the six: each of P1–P4 uses the same normalized
  schema) and for b2 (expiry is a policy over time, not a storage shape). Accepted for b5,
  which genuinely is a variant.
- *Add the backlog's six studies instead.* Rejected: out of scope by the owner's instruction
  and by the boundary that this brainstorm authorizes no measurement.

## Disagreements

Anticipated, without reading the other position:

- That optional `run_tag` is deliberate, because study 01 predates the tagging rule. My answer:
  predating is a legitimate historical fact, and the fix is to *record* it, not to re-run;
  a null that means two different things ("no tag exists" vs "not filled in") is the defect.
- That taxonomy work belongs to `book-synthesis-v1` and should not be pre-empted here. My
  answer: the synthesis task is the natural *owner* of the taxonomy; this brainstorm's value is
  naming which measured designs lack a home before it starts, which is cheaper now than after.
- That `claim-07`'s disagreement between trials is a harness problem rather than a reading
  rule. My reading: the analysis itself publishes it as a limit, so it is evidence for b7.

## Falsification tests

- **a1/a2**: a resolver that already reads git tags instead of `claims.json.run_tag`, plus a
  decision that the field is decorative — then a1 is cosmetic and only a2 (claim-01) stands.
- **a2**: any committed artefact that ties `20260912-small` to a commit — e.g.
  `git log --all --format=%H -S'20260912-small'` yielding a commit that recorded it, or a
  worktree manifest from the pre-rewrite history under `docs/history/`. If it exists, the item
  becomes a pointer to fix rather than a label to add.
- **a3/a4**: a documented statement that `run/*` tags may exist without result directories.
  Then the item is already-known convention and should be dropped.
- **b1**: if `sql/p1..p4` differ in the table definition as well as the arbiter, the family is
  confounded and the README's claim that the pair isolates the arbiter is wrong — a check that
  costs one `diff -r`. *Ran it on P1 vs P2: `indexes.sql` identical, `schema.sql` differs in
  comments only, `writes.sql` differs in the locking clause plus the exact non-locking fallback
  `SKIP LOCKED` makes necessary. The pair isolates the arbiter. The same check on the other
  arbiter pairs is worth one command each before the book repeats the table.*
- **b3**: if `04 d1_doc_row` and `01 D6` share the write-amplification mechanism, then b3 and
  the existing embedded-documents family are the same model and one should be merged.
