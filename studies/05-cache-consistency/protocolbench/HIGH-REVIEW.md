# Premeasurement HIGH review

Task `task-20260923T235925Z-cache-untested-protocols`; GPT-6 via Codex desktop,
declared HIGH, effort unknown. Reviewed 2026-10-03 before any database trial.
This is a self-review, not an independent analyst's endorsement.

The scientific contract is unchanged. Expected payloads, absence and committed
history come from the Go operation ledger; candidates cannot admit their own
cached value as an expected answer. The paused-fill, impossible-value,
disabled-coherence, lost-update, outbox-barrier and expiring-fence controls
exercise distinct wrong paths. A failed complete cell retains its error and
control record while removing all its trial timings.

Arm A holds source/outbox SQL, indexes, cache, client and service budgets
constant. CAS and explicit row locks are crossed with three publication
protocols. The client's per-key ordering preserves the established oracle
contract; it prevents an intra-key contention/scalability claim. Arm B compares
deletion packages with one publication protocol and retains a separate durable
generation. Its result cannot show that soft deletion alone ensures coherence.

Three whole runs restart the database and cache with fresh state. Each phase
has three trials. Seeded candidate order, repeated first/last position controls,
rotated deletion-package order, per-container throttling and Redis eviction
records bound comparisons. EXPLAIN ANALYZE is captured for the read, version
and relay statements. Short closed-loop phases establish supplied throughput,
not a demand knee. Configured and guest-observed resources are both recorded.

Required-comparison table in the precommitted protocol is accepted with the
explicit single-node/cluster and placement coverage gaps. No existing result,
signed conclusion or book claim is replaced. No cluster run is authorized by
this bounded envelope. The tiny keyed projection's fixed payload is a declared
scope limit and must accompany the analysis.

Verdict: proceed to the committed harness build and controlled matrix after
the image's oracle race tests and vet pass. Any failing runtime control stops
admission of that cell's timing; correct the implementation transparently,
retain the failed records and create a new source/run tag.

NEXT MODEL: HIGH
