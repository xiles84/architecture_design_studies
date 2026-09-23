# Progress — Study 05 v2 equal-total, engines and placement

**Task:** `task-20260922T200831Z-study05-v2-resources-engines`
**Attempt:** `attempt-20260923T111955Z-b1e692` (claim `claim-434ba2982ca22c00`, epoch 1)
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-05/v2-resources-engines`

Measured: eight cells (equal-total vs db-only; strict/relaxed on yb-single and yb-cluster3;
colocated/non-colocated). One benchmark lock at a time.

## Decision Log

1. **Fix the multi-endpoint defect rather than drop the cluster arm.** Every yb-cluster3 cell failed
   at connect (`sslmode is invalid`) because a comma DSN list went to one pgx pool. Using only the
   first endpoint would have measured one node and called it a cluster, which the methodology
   forbids. Added `pgxdb.OpenSpread` (one pool per endpoint, round robin) and recorded
   `connection_nodes`. Confidence: High.
2. **Re-ran the whole matrix after the fix** so the reported numbers come from one coherent commit.
   The pre-fix failed runs are kept and excluded. Confidence: High.
3. **Equal-total is reported as a separately labelled arm**, never pooled with db-only. Confidence: High.
4. **The placement distribution is reported as a gap.** `placement-evidence.txt` holds only the
   server list; `yb-admin list_tablets` produced nothing on this image, and the command had named a
   Study 01 table. The runner command was corrected for future runs, but the evidence for *this*
   task is not a tablet/leader distribution. Confidence: High.
5. **One host is not colocation evidence**, and the study's own report says so; the y1/y2 numbers are
   reported as key-choice locality on one machine. Confidence: High.

## Residual risks

- Placement distribution unevidenced; one host; single runs; per-endpoint operation counts not instrumented.
