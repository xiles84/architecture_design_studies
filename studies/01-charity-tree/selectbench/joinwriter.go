package selectbench

import "fmt"

// Matched blended read/write phase (review R2 / handoff C5).
//
// The blended phase runs one writer that increments an independent maintenance
// counter on a fixed set of existing child rows. The counter is not in the
// catalogue projection, is not part of the rank or index set and is not in the
// MD5 reuse token, so a concurrent update cannot change a measured read answer:
// the independent oracles stay valid in both phases. The isolated phase loads
// exactly the same column and data with the writer disabled.
//
// The writer is deliberately NOT an optimistic-versus-pessimistic comparison;
// it is a fixed maintenance load, held constant across spellings.

// WriterTargetRate is the offered rate, in committed updates per second.
const WriterTargetRate = 10

// MaintenanceKeys is the fixed set of child ids the writer rotates through, in
// order. Using the lowest ids keeps the set identical for every spelling and
// every run, so a mutation count is comparable across cells.
func MaintenanceKeys(factCount, n int) []int {
	if n < 1 {
		n = 1
	}
	if n > factCount {
		n = factCount
	}
	out := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		out = append(out, i)
	}
	return out
}

// WriterStats is the writer's own account of the phase: what it offered, what
// committed, how long each update took, and the per-key ledger reconciliation
// against the database.
type WriterStats struct {
	Phase         string         `json:"phase"`
	TargetRate    int            `json:"target_rate_per_sec"`
	Offered       int64          `json:"offered"`
	Completed     int64          `json:"completed"`
	Failed        int64          `json:"failed"`
	WallSeconds   float64        `json:"wall_seconds"`
	OfferedRate   float64        `json:"offered_rate_per_sec"`
	CompletedRate float64        `json:"completed_rate_per_sec"`
	LatencyMS     []float64      `json:"latency_ms,omitempty"`
	PerKey        map[string]int `json:"per_key,omitempty"`
	// EndpointOps is how many updates the writer sent to each endpoint, so a
	// cluster arm cannot be silently single-ended.
	EndpointOps []int `json:"endpoint_ops,omitempty"`
	// LedgerKeys is how many keys were reconciled; LedgerMismatch how many
	// database counter values disagreed with the writer's intended increments.
	LedgerKeys     int    `json:"ledger_keys"`
	LedgerMismatch int    `json:"ledger_mismatch"`
	LedgerDetail   string `json:"ledger_detail,omitempty"`
	Error          string `json:"error,omitempty"`
}

// joinPerKeySpread describes how evenly the writer's updates landed across its
// fixed key set; a writer stuck on one key would not exercise row maintenance.
func joinPerKeySpread(perKey map[string]int) string {
	if len(perKey) == 0 {
		return "no keys touched"
	}
	mn, mx := -1, 0
	for _, v := range perKey {
		if mn < 0 || v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return fmt.Sprintf("%d keys, min %d, max %d per key", len(perKey), mn, mx)
}

// joinLatencySummary summarizes the writer's per-update latency samples.
func joinLatencySummary(lat []float64) string {
	if len(lat) == 0 {
		return "no completed updates"
	}
	var sum float64
	mn, mx := lat[0], lat[0]
	for _, v := range lat {
		sum += v
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return fmt.Sprintf("avg %.3f ms, min %.3f, max %.3f over %d", sum/float64(len(lat)), mn, mx, len(lat))
}
