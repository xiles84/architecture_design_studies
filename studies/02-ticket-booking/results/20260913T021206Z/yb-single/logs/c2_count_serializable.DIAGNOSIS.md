# Why `yb-single / c2_count_serializable` failed

Diagnosed on 2026-09-13 from the logs captured in `c2_count_serializable.yb-single.log` and
from the full tserver logs inside the still-running `yb-single` container.

**Symptom.** During the sell-out race, every connection failed with
`FATAL: the database system is shutting down (SQLSTATE 57P03)`; the harness could not even
read the final count of race event 503. The same failure happened in dev check
`devcheck-3` — two out of two attempts on YugabyteDB 1-node.

**Cause, in order:**

1. C2 is count-then-insert at SERIALIZABLE. With 32 buyers on one event, the transactions'
   read and write locks conflict with each other continuously: the tserver WARNING log holds
   **78 438** `deadlock_detector ... Timed out` entries for the cell.
2. The node runs under a 2-CPU quota and was throttled throughout (see
   `c2_count_serializable.dbcpu.txt`). Starved of CPU, the tserver's own RPCs to the master
   (port 7100) and its PG client service (port 9100) went unanswered for 28–32 s against a
   15 s timeout (`yb_rpc.cc:548 ... Rpc timeout, passed: 29.388s, timeout: 15.000s`).
3. The YSQL lease could not be renewed:
   `ysql_lease_manager.cc:280] Lease has expired, killing pg sessions.` (04:06:45 UTC)
4. The tserver terminated every PostgreSQL-layer backend ("terminating connection due to
   administrator command") and refused new connections while YSQL restarted.

**What it is and is not.** This is YugabyteDB protecting consistency when a node loses
contact with its master; the node lost contact because it was overloaded. It is not a
correctness failure of design C2 (every audit before the failure was consistent), and it is
not a harness bug. It is a real outcome of SERIALIZABLE check-then-insert under a hot drop
**on a 2-CPU node**, and must not be generalised to adequately provisioned nodes without a
run that shows it there.
