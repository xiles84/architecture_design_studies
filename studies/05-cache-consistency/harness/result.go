package main

import (
	"time"

	"adsplatform/adapters/cgroup"
	"adsplatform/core/measure"
)

// The result types. A result file must name the environment, the commit, the image
// identities, the seeds and the resource conditions, because a number without them
// cannot be reproduced or compared (methodology 2 and 12a).

type Options struct {
	Scale         string        `json:"scale"`
	Seed          int64         `json:"seed"`
	FaultSeed     int64         `json:"fault_seed"`
	Workers       int           `json:"workers"`
	Writers       int           `json:"writer_workers"`
	Duration      time.Duration `json:"duration_per_measurement"`
	Warmup        time.Duration `json:"warmup"`
	Trials        int           `json:"trials"`
	Retries       int           `json:"bounded_retries"`
	Phases        string        `json:"phases"`
	Instances     int           `json:"logical_instances"`
	ChurnFor      time.Duration `json:"churn_duration"`
	Stampede      int           `json:"stampede_readers"`
	StampedeN     int           `json:"stampede_keys"`
	HotKeys       int           `json:"hot_keys"`
	CapacityKB    int           `json:"cache_capacity_kib"`
	CacheFit      bool          `json:"cache_capacity_fits_working_set"`
	RedisAddr     string        `json:"redis_addr,omitempty"`
	RedisMaxMB    int           `json:"redis_maxmemory_mb,omitempty"`
	ResourceFrame string        `json:"resource_framing"`
	ExtFraction   int           `json:"external_write_fraction_pct"`
}

type EngineInfo struct {
	Version            string `json:"version"`
	ServerVersion      string `json:"server_version"`
	EffectiveIsolation string `json:"effective_isolation"`
}

type DatasetInfo struct {
	Seed            int64   `json:"seed"`
	Charities       int     `json:"charities"`
	People          int     `json:"people"`
	Donations       int     `json:"donations"`
	SerializedBytes int64   `json:"serialized_bytes"`
	WorkingSetBytes int64   `json:"logical_working_set_bytes"`
	MeanPayloadB    float64 `json:"mean_payload_bytes"`
}

// ClientCPUInfo records the client container's own CFS throttling across the cell.
// A throttled client puts milliseconds into tails that belong to no scenario, so the
// counters travel with every result rather than being assumed absent
// (methodology 7).
type ClientCPUInfo struct {
	Before cgroup.CPUStat `json:"before"`
	After  cgroup.CPUStat `json:"after"`
	Delta  cgroup.CPUStat `json:"delta"`
	Seen   bool           `json:"cgroup_visible"`
}

// CacheStats is the policy's own accounting, kept apart from the backend's so that
// "the broker evicted 400 keys" and "this design never filled the key" are two
// different observations.
type CacheStats struct {
	Backend                 string           `json:"backend"`
	CapacityBytes           int64            `json:"capacity_bytes"`
	ResidentBytes           int64            `json:"resident_bytes"`
	Items                   int64            `json:"items"`
	Evictions               int64            `json:"evictions"`
	EvictedKeys             int64            `json:"evicted_keys_reported_by_backend"`
	Fills                   int64            `json:"fills"`
	FillErrors              int64            `json:"fill_errors"`
	Publishes               int64            `json:"publishes"`
	PublishFenced           int64            `json:"publishes_refused_by_fence"`
	PublishFailed           int64            `json:"publish_failures"`
	Invalidations           int64            `json:"invalidations"`
	Tombstones              int64            `json:"tombstones_installed_pre_commit"`
	ValidationQueries       int64            `json:"authoritative_version_validations"`
	BypassReads             int64            `json:"cache_bypass_reads"`
	FallbackReads           int64            `json:"lease_fallback_reads"`
	ExternalWrites          int64            `json:"external_writes_bypassing_adapter"`
	AmbiguousWrites         int64            `json:"ambiguous_writes_invalidated"`
	SuppressedInvalidations int64            `json:"suppressed_invalidations"`
	UnrecordedConfirmed     int64            `json:"unrecorded_committed_states_confirmed"`
	BackendStats            map[string]int64 `json:"backend_stats,omitempty"`
	PayloadBytes            int64            `json:"logical_payload_bytes"`
	MetadataBytes           int64            `json:"metadata_bytes"`
}

// PhaseWrong is one phase's wrong-read account. Per phase, because a design that is
// stale only while writes are in flight is a different finding from one that is
// stale during warm reads, and a single cell-wide number would hide which.
type PhaseWrong struct {
	Phase   string           `json:"phase"`
	Summary WrongReadSummary `json:"summary"`
}

type StampedeResult struct {
	Readers        int     `json:"readers"`
	Keys           int     `json:"keys"`
	ElapsedMS      float64 `json:"elapsed_ms"`
	DatabaseLoads  int64   `json:"database_loads"`
	LoadsPerKey    float64 `json:"database_loads_per_key"`
	LeaseAcquired  int64   `json:"lease_acquisitions"`
	LeaseContended int64   `json:"lease_contended"`
	Fallbacks      int64   `json:"fallbacks"`
	DuplicateFills int64   `json:"duplicate_fills"`
	WrongReads     int64   `json:"wrong_reads"`
	Impossible     int64   `json:"impossible_values"`
	Note           string  `json:"note"`
}

type ChurnResult struct {
	DurationS     float64 `json:"duration_s"`
	TTLBoundaries float64 `json:"hard_ttl_boundaries_crossed"`
	Reads         int64   `json:"reads"`
	Writes        int64   `json:"writes"`
	HardExpired   int64   `json:"expired_hard"`
	ProbExpired   int64   `json:"expired_probabilistic"`
	WrongReads    int64   `json:"wrong_reads"`
	Impossible    int64   `json:"impossible_values"`
	Note          string  `json:"note"`
}

type InstanceResult struct {
	Instances      int       `json:"logical_instances"`
	Workers        int       `json:"total_workers"`
	Backend        string    `json:"backend"`
	AppOpsPerSec   float64   `json:"total_application_ops_per_sec"`
	AppTrials      []float64 `json:"total_application_trials_ops_per_sec,omitempty"`
	AppSpreadPct   float64   `json:"total_application_spread_pct"`
	CacheOpsPerSec float64   `json:"cacheable_endpoint_ops_per_sec"`
	WrongReads     int64     `json:"wrong_reads"`
	Impossible     int64     `json:"impossible_values"`
	BypassReads    int64     `json:"cache_bypass_reads"`
	LocalLRUItems  int64     `json:"local_lru_items_final"`
	Note           string    `json:"note"`
}

type FaultResult struct {
	Name        string `json:"name"`
	Reproduced  bool   `json:"reproduced"`
	Detail      string `json:"detail"`
	Correctness string `json:"correctness"`
	WrongReads  int64  `json:"wrong_reads_observed"`
	Impossible  int64  `json:"impossible_values_observed"`
}

type AuditInfo struct {
	Passed              bool           `json:"passed"`
	Checks              int            `json:"checks"`
	KeyValueMismatches  int            `json:"portal_content_mismatches"`
	DonationSetMismatch int            `json:"donation_set_mismatches"`
	AggregateMismatch   int            `json:"aggregate_mismatches"`
	RecentSliceMismatch int            `json:"recent_slice_mismatches"`
	DesignChecks        map[string]int `json:"design_checks,omitempty"`
	Failures            []string       `json:"failures,omitempty"`
}

type GateInfo struct {
	Passed     bool     `json:"passed"`
	Checks     int      `json:"checks"`
	Failures   []string `json:"failures,omitempty"`
	Statements []string `json:"statements,omitempty"`
}

type CellResult struct {
	Study         string `json:"study"`
	Scenario      string `json:"scenario"`
	ScenarioShort string `json:"scenario_short"`
	Group         string `json:"group"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Risk          string `json:"risk"`
	Pair          string `json:"controlled_pair,omitempty"`
	Negative      bool   `json:"negative_control,omitempty"`
	// Every dimension of the scenario is retained separately as well as in the id,
	// so a report can never show a number without its contract.
	ModelDim     string `json:"dim_model"`
	VersionDim   string `json:"dim_version"`
	BackendDim   string `json:"dim_backend"`
	StrategyDim  string `json:"dim_strategy"`
	FreshnessDim string `json:"dim_freshness"`
	WritersDim   string `json:"dim_writers"`
	StrictPolicy string `json:"strict_freshness_policy"`

	Topology     string `json:"topology"`
	Engine       string `json:"engine"`
	RunID        string `json:"run_id"`
	Environment  string `json:"environment"`
	RepoCommit   string `json:"repo_commit"`
	RepoDesc     string `json:"repo_describe"`
	RepoDirty    bool   `json:"repo_dirty"`
	RunTag       string `json:"run_tag,omitempty"`
	BenchImage   string `json:"bench_image,omitempty"`
	BenchImageID string `json:"bench_image_id,omitempty"`

	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`

	Options    Options        `json:"options"`
	EngineInfo EngineInfo     `json:"engine_info"`
	Dataset    DatasetInfo    `json:"dataset"`
	Load       LoadInfo       `json:"load"`
	Gate       GateInfo       `json:"gate"`
	ClientCPU  *ClientCPUInfo `json:"client_cpu,omitempty"`

	Cache   CacheStats        `json:"cache"`
	Redis   map[string]string `json:"redis_info,omitempty"`
	Storage map[string]int64  `json:"storage_bytes,omitempty"`

	Warm   []measure.Result `json:"warm_reads,omitempty"`
	Mixed  []measure.Result `json:"cacheable_endpoint,omitempty"`
	AppMix []measure.Result `json:"total_application,omitempty"`
	Writes []measure.Result `json:"writes,omitempty"`
	Wrong  []PhaseWrong     `json:"wrong_reads,omitempty"`

	Stampede  *StampedeResult  `json:"stampede,omitempty"`
	Churn     *ChurnResult     `json:"churn,omitempty"`
	Instances []InstanceResult `json:"instances,omitempty"`
	Faults    []FaultResult    `json:"faults,omitempty"`

	Audits           []AuditInfo                 `json:"audits,omitempty"`
	StrictViolations int64                       `json:"strict_contract_violations"`
	ImpossibleValues int64                       `json:"impossible_cache_values"`
	Lease            LeaseStats                  `json:"lease"`
	ExpiryByAge      map[string]map[string]int64 `json:"expiry_by_age_bucket,omitempty"`

	Error string `json:"error,omitempty"`
}
