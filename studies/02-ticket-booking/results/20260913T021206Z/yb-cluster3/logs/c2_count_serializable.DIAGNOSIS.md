# Why `yb-cluster3 / c2_count_serializable` failed

Diagnosed on 2026-09-13 from the tserver logs inside the still-running cluster containers.

Same failure, same cause as `yb-single/logs/c2_count_serializable.DIAGNOSIS.md`:

- During the sell-out race, connections to `yb-n1` were refused (the YSQL layer was
  restarting); the harness could not read the final count of race event 502.
- `yb-n1` tserver WARNING log: **49 140** `deadlock_detector ... Timed out` entries, then
  `ysql_lease_manager.cc:280] Lease has expired, killing pg sessions.` at 07:02:48 UTC.
- `yb-n2` and `yb-n3`: **0** deadlock-detector timeouts and no lease expiry.

**Why only one node.** The harness connects every client — buyers, organiser, audits — to
`yb-n1`, as study 01's harness did. The PostgreSQL layer of that one node therefore carries
all SQL processing and all lock-conflict resolution for the cluster, while the other two
nodes hold tablet replicas and leaders but no sessions. The lease expiry on `yb-n1` is the
same overload-induced self-protection seen on the single node. It says nothing about how C2
would behave with connections spread across all three nodes; that would be a separate run.
