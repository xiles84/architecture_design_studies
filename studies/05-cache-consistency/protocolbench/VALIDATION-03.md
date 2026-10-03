# Runtime control validation

2026-10-03, GPT-6 via Codex desktop. Producing commit `8a6e7f09`, immutable
run tag `run/study-05/cache-v5-controls-20261003-r1`, input digest
`02900bf8b0429be4`. One short control run, three 500-ms trials per phase;
it is not the final three-whole-run performance comparison.

All 10 cells passed their complete gates; all 57 controls passed. The generated
report and results are retained under the run id. Read-only source-query
contrasts were exactly one per version-validated hit and zero per synchronous
or outbox hit. No candidate cell observed stale-after-ack or impossible reads.
The disabled protocol, fabricated payload, lost-update, omitted outbox barrier
and expiring-fence negative controls fired. Tombstone expiry was observed by GET.

The completed runner removed only its own database/cache containers and volume,
then released the shared lock. A portability artifact `NUL` was created by the
Git Bash SSH executable when native Podman passed Windows' null path as its
known-hosts destination. It contained only this machine's public host key. The
runner now obtains guest `/proc/meminfo` and online CPU count from the source
container, avoiding SSH entirely; the exact task-owned artifact is removed.

Guest-observed memory was approximately 16 GiB, distinct from the saved 2-GiB
machine configuration. Effective container limits remain explicitly inspected.
The final matrix still has the predeclared single-node and fixed-payload scope.

Added final whole-run reduction tooling and accurate audit-error counters;
the next producing commit and full matrix will rebuild and retest those changes.

NEXT MODEL: HIGH
