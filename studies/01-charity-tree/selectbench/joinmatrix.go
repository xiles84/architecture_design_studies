package selectbench

import (
	"fmt"
	"math"
)

// Join/CTE matrix definitions.
//
// Two contracts are measured separately and never pooled:
//
//   - the SELECTIVITY matrix: the ten bag spellings (six default inner-join
//     forms plus the two CTE spellings under each forced policy) at six exact
//     rank thresholds;
//   - the REUSE matrix: the two CTE spellings under the three policies at five
//     reference counts and two per-consumer selections.
//
// One-use spellings and reuse contracts have separate predictions (P1-P3 vs
// P4), so they get separate result tables.

// joinPercents are the reviewed selectivity levels.
var joinPercents = []float64{0.01, 0.1, 1, 10, 50, 90}

// JoinStatements are the ten bag spellings, in report order.
var JoinStatements = []string{
	"j01_join_where",
	"j02_join_on",
	"j03_filtered_derived",
	"j04_cte_filtered",
	"j04m_cte_filtered_mat",
	"j04n_cte_filtered_notmat",
	"j05_cte_joined_outer",
	"j05m_cte_joined_outer_mat",
	"j05n_cte_joined_outer_notmat",
	"j06_comma_join",
}

// JoinThresholds maps the reviewed percentages onto exact rank counts for a
// fact-table size. On the reviewed 200,000 rows these are the protocol's
// 20 / 200 / 2,000 / 20,000 / 100,000 / 180,000.
func JoinThresholds(factCount int) []int {
	out := make([]int, 0, len(joinPercents))
	for _, pct := range joinPercents {
		t := int(math.Round(float64(factCount) * pct / 100))
		if t < 1 {
			t = 1
		}
		out = append(out, t)
	}
	return out
}

// JoinThresholdPct returns the reviewed percentage label for a threshold.
func JoinThresholdPct(threshold, factCount int) string {
	if factCount <= 0 {
		return fmt.Sprintf("threshold=%d", threshold)
	}
	pct := 100 * float64(threshold) / float64(factCount)
	return fmt.Sprintf("%g%%", pct)
}

func shortJoinStatement(name string) string {
	switch name {
	case "j01_join_where":
		return "join_where"
	case "j02_join_on":
		return "join_on"
	case "j03_filtered_derived":
		return "filtered_derived"
	case "j04_cte_filtered":
		return "cte_filtered"
	case "j04m_cte_filtered_mat":
		return "cte_filtered_mat"
	case "j04n_cte_filtered_notmat":
		return "cte_filtered_notmat"
	case "j05_cte_joined_outer":
		return "cte_joined_outer"
	case "j05m_cte_joined_outer_mat":
		return "cte_joined_outer_mat"
	case "j05n_cte_joined_outer_notmat":
		return "cte_joined_outer_notmat"
	case "j06_comma_join":
		return "comma_join"
	}
	return name
}

// JoinBagAxis names one selectivity cell.
func JoinBagAxis(statement string, threshold int) string {
	return fmt.Sprintf("family=bag;stmt=%s;threshold=%d", shortJoinStatement(statement), threshold)
}

// JoinReuseAxis names one reuse cell. It is deliberately a different family id
// from the bag cells so the two contracts can never be pooled.
func JoinReuseAxis(policy string, refs, consumerPct int) string {
	return fmt.Sprintf("family=reuse;policy=%s;refs=%d;consumer=%d%%", policy, refs, consumerPct)
}

// BuildJoinSpecs is the reviewed matrix. base carries the timing budget; every
// cell overrides only its own axis parameters.
func BuildJoinSpecs(base JoinArgs, factCount, trials int) []JoinSpec {
	var specs []JoinSpec
	for _, t := range JoinThresholds(factCount) {
		for _, name := range JoinStatements {
			a := base
			a.Threshold = t
			a.Refs = 1
			a.Trials = trials
			specs = append(specs, JoinSpec{Statement: name, Axis: JoinBagAxis(name, t), Args: a})
		}
	}
	// The reuse contract fixes one selective CTE (1% of the facts) and varies
	// the reference count and the per-consumer selection.
	reuseThreshold := int(math.Round(float64(factCount) * 1 / 100))
	if reuseThreshold < 1 {
		reuseThreshold = 1
	}
	for _, policy := range []string{"default", "mat", "notmat"} {
		for _, refs := range ReuseLevels() {
			for _, pct := range ReuseConsumerPercents() {
				a := base
				a.Threshold = reuseThreshold
				a.Refs = refs
				a.ConsumerMax = int(math.Round(float64(factCount) * float64(pct) / 100))
				if a.ConsumerMax < 1 {
					a.ConsumerMax = 1
				}
				a.Trials = trials
				specs = append(specs, JoinSpec{Statement: ReuseStatementName(policy, refs, pct),
					Axis: JoinReuseAxis(policy, refs, pct), Args: a})
			}
		}
	}
	return specs
}

// JoinTopologyThresholds are the three named selectivities of the reviewed
// bounded single/cluster subset: 0.01%, 1% and 90%.
func JoinTopologyThresholds(factCount int) []int {
	out := make([]int, 0, 3)
	for _, pct := range []float64{0.01, 1, 90} {
		t := int(math.Round(float64(factCount) * pct / 100))
		if t < 1 {
			t = 1
		}
		out = append(out, t)
	}
	return out
}

// BuildJoinTopologySpecs is the bounded single/cluster subset the review
// approved: JOIN+WHERE and the joined-output CTE under its three policies at the
// three named selectivities, plus the reuse ladder at {1,4,16} references. It is
// deliberately smaller than the full 91-cell matrix so the topology arms fit the
// run budget; the axes use the same family ids so the pair tables render the
// same decisions.
func BuildJoinTopologySpecs(base JoinArgs, factCount, trials int) []JoinSpec {
	var specs []JoinSpec
	bag := []string{"j01_join_where", "j05_cte_joined_outer", "j05m_cte_joined_outer_mat", "j05n_cte_joined_outer_notmat"}
	for _, t := range JoinTopologyThresholds(factCount) {
		for _, name := range bag {
			a := base
			a.Threshold = t
			a.Refs = 1
			a.Trials = trials
			specs = append(specs, JoinSpec{Statement: name, Axis: JoinBagAxis(name, t), Args: a})
		}
	}
	reuseThreshold := int(math.Round(float64(factCount) * 1 / 100))
	if reuseThreshold < 1 {
		reuseThreshold = 1
	}
	for _, policy := range []string{"default", "mat", "notmat"} {
		for _, refs := range []int{1, 4, 16} {
			for _, pct := range ReuseConsumerPercents() {
				a := base
				a.Threshold = reuseThreshold
				a.Refs = refs
				a.ConsumerMax = int(math.Round(float64(factCount) * float64(pct) / 100))
				if a.ConsumerMax < 1 {
					a.ConsumerMax = 1
				}
				a.Trials = trials
				specs = append(specs, JoinSpec{Statement: ReuseStatementName(policy, refs, pct),
					Axis: JoinReuseAxis(policy, refs, pct), Args: a})
			}
		}
	}
	return specs
}

// JoinControlSpecs are the controls the family's gate must see fire.
func JoinControlSpecs(base JoinArgs, factCount int) []JoinSpec {
	var specs []JoinSpec
	for _, name := range []string{"x01_optional_where", "x02_join_limit_before"} {
		a := base
		a.Threshold = int(math.Round(float64(factCount) * 10 / 100))
		if a.Threshold < 1 {
			a.Threshold = 1
		}
		a.Trials = 1
		specs = append(specs, JoinSpec{Statement: name, Axis: "control=" + name, Args: a})
	}
	return specs
}

// JoinControlResult is one control's disposition. The control must be DETECTED
// (its answer differs from the contract it pretends to serve).
type JoinControlResult struct {
	Statement string `json:"statement"`
	Group     string `json:"group"`
	Detected  bool   `json:"detected"`
	Rows      int    `json:"rows"`
	Error     string `json:"error,omitempty"`
}

// JoinSequence reports the axis of each join spec in plan order.
func (p OrderPlan) JoinSequence(specs []JoinSpec) []string {
	out := make([]string, 0, len(p.Order))
	for _, i := range p.Order {
		out = append(out, specs[i].Axis)
	}
	return out
}

// JoinFeasibility is the per-engine feasibility evidence the handoff requires
// before the matrix may run: statement support, actual materialization, exact
// selective counts, index use, expensive-work support and drained results.
type JoinFeasibility struct {
	FactCount        int               `json:"fact_count"`
	Customers        int               `json:"customers"`
	Thresholds       []int             `json:"thresholds"`
	ThresholdCounts  map[string]int    `json:"threshold_counts"`
	StatementSupport map[string]string `json:"statement_support"`
	PolicyEvidence   map[string]string `json:"policy_evidence"`
	RankIndexUsed    map[string]bool   `json:"rank_index_used"`
	// DefaultIndexGate is the physical-arm gate (review M1): for the default join
	// at the three most selective thresholds, whether the plan used the rank
	// index, with the plan text kept when it did not.
	DefaultIndexGate    map[string]string `json:"default_index_gate,omitempty"`
	DefaultIndexPlans   map[string]string `json:"default_index_plans,omitempty"`
	ReuseExprSupported  bool              `json:"reuse_expr_supported"`
	ReuseExprError      string            `json:"reuse_expr_error,omitempty"`
	MaxThresholdDrained int               `json:"max_threshold_drained_rows"`
	TokenEvidence       string            `json:"token_evidence,omitempty"`
	Notes               []string          `json:"notes,omitempty"`
}

// JoinResult is one whole join/CTE run.
type JoinResult struct {
	RunID                string              `json:"run_id"`
	Engine               string              `json:"engine"`
	Topology             string              `json:"topology"`
	Scale                string              `json:"scale"`
	Arm                  string              `json:"arm,omitempty"`
	FactCount            int                 `json:"fact_count"`
	Customers            int                 `json:"customers"`
	Seed                 int64               `json:"seed"`
	Thresholds           []int               `json:"thresholds"`
	Provenance           Provenance          `json:"provenance"`
	Cells                []Cell              `json:"cells"`
	Controls             []JoinControlResult `json:"controls,omitempty"`
	Order                RunOrder            `json:"order"`
	Feasibility          *JoinFeasibility    `json:"feasibility,omitempty"`
	ExplainOptions       string              `json:"explain_options,omitempty"`
	Endpoints            []string            `json:"endpoints,omitempty"`
	EndpointOps          []int               `json:"endpoint_ops,omitempty"`
	EndpointConns        []int               `json:"endpoint_conns,omitempty"`
	Phase                string              `json:"phase,omitempty"`
	Writer               *WriterStats        `json:"writer,omitempty"`
	WriterSamples        int                 `json:"writer_samples,omitempty"`
	WriterOverlapSamples int                 `json:"writer_overlap_samples,omitempty"`
	Summary              map[string]int      `json:"summary"`
	// IndexDefs and Placement are the engine's own physical evidence: the
	// indexdef text of every domain index, and (YugabyteDB) the colocation and
	// tablet properties of each domain table. They make the physical arm
	// checkable rather than asserted.
	IndexDefs []IndexDef       `json:"index_defs,omitempty"`
	Placement []TablePlacement `json:"placement,omitempty"`
	// NodeThrottle is per-node cgroup CPU throttling for the run, computed by the
	// runner from the before/after cpu.stat files. It is loaded at render time,
	// because the "after" file exists only once the run has finished.
	NodeThrottle []NodeThrottleStat `json:"node_throttle,omitempty"`
	Notes        []string           `json:"notes,omitempty"`
	Error        string             `json:"error,omitempty"`
}

// NodeThrottleStat is one container's CPU-throttling delta over the run.
type NodeThrottleStat struct {
	Node          string `json:"node"`
	NrPeriods     int64  `json:"nr_periods"`
	NrThrottled   int64  `json:"nr_throttled"`
	ThrottledUsec int64  `json:"throttled_usec"`
	CPUQuota      string `json:"cpu_quota,omitempty"`
	MemoryMax     string `json:"memory_max,omitempty"`
}

// IndexDef is one index's engine-reported definition.
type IndexDef struct {
	Name string `json:"name"`
	Def  string `json:"def"`
}

// TablePlacement is one table's engine-reported physical placement.
type TablePlacement struct {
	Table            string `json:"table"`
	NumTablets       int    `json:"num_tablets"`
	NumHashKeyCols   int    `json:"num_hash_key_columns"`
	IsColocated      bool   `json:"is_colocated"`
	TablegroupOID    string `json:"tablegroup_oid,omitempty"`
	ColocationID     string `json:"colocation_id,omitempty"`
	TabletCountLocal int    `json:"tablet_count_local,omitempty"`
}

// Valid counts cells whose answer matched the oracle and whose timing measured.
func (r JoinResult) Valid() int { return r.Summary["valid"] }

// Find returns the cell for an axis, if present.
func (r JoinResult) Find(axis string) (Cell, bool) {
	for _, c := range r.Cells {
		if c.Axis == axis {
			return c, true
		}
	}
	return Cell{}, false
}

// JoinPair is one controlled comparison between two bag spellings at one
// threshold.
type JoinPair struct {
	Threshold int
	A, B      Cell
	RatioBA   float64
	Label     string
}

// BagPairs compares every non-baseline bag spelling with the canonical join
// spelling at the same threshold.
func (r JoinResult) BagPairs() []JoinPair {
	var out []JoinPair
	baseline := "j01_join_where"
	for _, t := range r.Thresholds {
		a, okA := r.Find(JoinBagAxis(baseline, t))
		if !okA {
			continue
		}
		for _, name := range JoinStatements {
			if name == baseline {
				continue
			}
			b, okB := r.Find(JoinBagAxis(name, t))
			if !okB {
				continue
			}
			out = append(out, newJoinPair(t, a, b, shortJoinStatement(name)+" vs "+shortJoinStatement(baseline)))
		}
	}
	return out
}

// PolicyPairs compares each CTE spelling's forced policies with its default
// spelling at the same threshold (P3).
func (r JoinResult) PolicyPairs() []JoinPair {
	var out []JoinPair
	families := [][3]string{
		{"j04_cte_filtered", "j04m_cte_filtered_mat", "j04n_cte_filtered_notmat"},
		{"j05_cte_joined_outer", "j05m_cte_joined_outer_mat", "j05n_cte_joined_outer_notmat"},
	}
	for _, t := range r.Thresholds {
		for _, f := range families {
			def, okD := r.Find(JoinBagAxis(f[0], t))
			if !okD {
				continue
			}
			for _, name := range []string{f[1], f[2]} {
				b, okB := r.Find(JoinBagAxis(name, t))
				if !okB {
					continue
				}
				out = append(out, newJoinPair(t, def, b, shortJoinStatement(name)+" vs "+shortJoinStatement(f[0])))
			}
		}
	}
	return out
}

// ReusePair is one policy comparison inside the reuse contract.
type ReusePair struct {
	Policy  string
	Refs    int
	Pct     int
	A, B    Cell
	RatioBA float64
	Label   string
}

// ReusePairs compares materialized and not-materialized against the default
// policy at each reuse level and consumer selection.
func (r JoinResult) ReusePairs() []ReusePair {
	var out []ReusePair
	for _, refs := range ReuseLevels() {
		for _, pct := range ReuseConsumerPercents() {
			def, okD := r.Find(JoinReuseAxis("default", refs, pct))
			if !okD {
				continue
			}
			for _, policy := range []string{"mat", "notmat"} {
				b, okB := r.Find(JoinReuseAxis(policy, refs, pct))
				if !okB {
					continue
				}
				p := ReusePair{Policy: policy, Refs: refs, Pct: pct, A: def, B: b,
					Label: policy + " vs default"}
				if def.MedianMS > 0 {
					p.RatioBA = b.MedianMS / def.MedianMS
				}
				out = append(out, p)
			}
		}
	}
	return out
}

func newJoinPair(threshold int, a, b Cell, label string) JoinPair {
	p := JoinPair{Threshold: threshold, A: a, B: b, Label: label}
	if a.MedianMS > 0 {
		p.RatioBA = b.MedianMS / a.MedianMS
	}
	return p
}

// JoinPositionPair returns the repeated-design position control's first and last
// cells, if the run carried it.
func (r JoinResult) JoinPositionPair() (JoinPair, bool) {
	first, last, ok := positionPairCells(r.Cells)
	if !ok {
		return JoinPair{}, false
	}
	return newJoinPair(0, first, last, "position control first vs last"), true
}
