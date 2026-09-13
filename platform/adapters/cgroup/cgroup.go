// Package cgroup reads the CPU accounting of the container the harness runs in.
//
// Every container in these studies is limited with a CFS quota (--cpus), never
// pinned. A quota has a failure mode of its own: a burst that uses the period's
// allowance early is throttled until the next period starts (100 ms by default),
// which shows up as a cluster of tail latencies near a few tens of milliseconds
// that belong to no design. Recording the throttling beside the measurement turns
// that suspicion into a number a reader can check.
package cgroup

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// CPUStat is the subset of cgroup v2 cpu.stat the studies report.
type CPUStat struct {
	UsageUsec     int64 `json:"usage_usec"`
	NrPeriods     int64 `json:"nr_periods"`
	NrThrottled   int64 `json:"nr_throttled"`
	ThrottledUsec int64 `json:"throttled_usec"`
}

// Read returns the current counters, or ok=false where cgroup v2 accounting is
// not visible (the harness run outside a container, or cgroup v1).
func Read() (CPUStat, bool) {
	f, err := os.Open("/sys/fs/cgroup/cpu.stat")
	if err != nil {
		return CPUStat{}, false
	}
	defer f.Close()
	var s CPUStat
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		k, v, ok := strings.Cut(sc.Text(), " ")
		if !ok {
			continue
		}
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			continue
		}
		switch k {
		case "usage_usec":
			s.UsageUsec = n
		case "nr_periods":
			s.NrPeriods = n
		case "nr_throttled":
			s.NrThrottled = n
		case "throttled_usec":
			s.ThrottledUsec = n
		}
	}
	return s, true
}

// Sub returns the counters accumulated between an earlier reading and this one.
func (s CPUStat) Sub(earlier CPUStat) CPUStat {
	return CPUStat{
		UsageUsec:     s.UsageUsec - earlier.UsageUsec,
		NrPeriods:     s.NrPeriods - earlier.NrPeriods,
		NrThrottled:   s.NrThrottled - earlier.NrThrottled,
		ThrottledUsec: s.ThrottledUsec - earlier.ThrottledUsec,
	}
}
