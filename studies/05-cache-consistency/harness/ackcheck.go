package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"adsplatform/adapters/pgxdb"
)

// The acknowledgement check. This is the experiment that decides whether the
// surviving strict violations are a harness defect or a property of the system.
//
// The failing samples were database reads (fill and bypass) returning a state exactly
// one behind the ledger's requirement, which resolves before a phase-end audit runs.
// Two explanations fit that shape:
//
//   - the harness acknowledges a write before the database makes the commit visible;
//   - the ledger's requirement advances for a state the database has not committed.
//
// Both are testable in one place: write through the adapter (which fences, commits and
// acknowledges), and IMMEDIATELY afterwards read the portal view from a SECOND database
// session and compare it with the requirement. A second session is essential -- reading
// through the writer's own pool could return the writer's own connection and see the
// commit trivially, turning a real defect into a clean result.
//
// The check runs in its own phase, outside any measured phase, so it costs the
// benchmark nothing: a verification read inside the measured write path would slow the
// thing it is verifying.
func (c *cell) runAckCheck(ctx context.Context) error {
	res := AckCheck{Passed: true}
	// The verification session: one connection, used by nothing else, so it cannot
	// inherit a transaction or a cached view from the writing pool.
	vdb, err := pgxdb.Open(ctx, pgxdb.Config{
		DSN: c.dsn, MaxConns: 1, StatementTimeoutMS: 60000,
		ApplicationName: "study05-ackcheck",
	})
	if err != nil {
		return fmt.Errorf("ack check could not open its verification session: %w", err)
	}
	defer vdb.Close()

	verify := func(m mutation, ackAtMS int64, part string) {
		res.Attempts++
		keys := m.keys()
		if len(keys) == 0 {
			return
		}
		keyID := keys[0]
		// The requirement is captured BEFORE the verification read, exactly as a real
		// read captures it before touching the database. Capturing it afterwards -- as
		// this helper did first -- makes the check manufacture its own violation: a
		// concurrent writer that commits between the read and the capture moves the
		// requirement forward, and the read is then judged against a state that did not
		// exist when it began. That is the same ordering mistake this study had already
		// fixed in the adapter, reproduced in the checker.
		reqHash, reqSeq := c.orc.Required(keyID)
		content, _, err := readPortalSQL(ctx, vdb, c.cat, c.d.usesVersion(), keyID)
		if err != nil {
			res.Errors++
			if len(res.Examples) < 8 {
				res.Examples = append(res.Examples, "verification read failed: "+err.Error())
			}
			return
		}
		dbHash := content.ContentHash()
		if dbHash == reqHash {
			return
		}
		kind, _, behind, _ := c.orc.Classify(keyID, dbHash, reqHash, reqSeq)
		if kind != KindStale {
			return
		}
		res.Violations++
		res.Passed = false
		if part == "sequential" {
			res.SequentialViolations++
		} else {
			res.ConcurrentViolations++
		}
		if len(res.Examples) < 8 {
			res.Examples = append(res.Examples, fmt.Sprintf(
				"%s: acknowledged donor %d and a SECOND session still sees a state %d behind the requirement %d ms later (requirements seq %d)",
				part, keyID, behind, nowMS()-ackAtMS, reqSeq))
		}
	}

	// Part 1: sequential. A single writer, so any violation is an ordering defect and
	// not a race.
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 31))
	for i := 0; i < 150; i++ {
		m, err := c.buildMutation(r)
		if err != nil {
			continue
		}
		if err := c.ad.Write(ctx, c.instIdx(r), m); err != nil {
			continue
		}
		res.Sequential++
		verify(m, nowMS(), "sequential")
	}

	// Part 2: concurrent. Eight writers, each verifying its own acknowledgement. This
	// separates "the ack is early" from "a concurrent commit wins the ordering".
	var wg sync.WaitGroup
	var mu sync.Mutex
	const writers, each = 4, 40
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			wr := rand.New(rand.NewSource(c.opts.FaultSeed + 100 + int64(w)*7919))
			for i := 0; i < each; i++ {
				m, err := c.buildMutation(wr)
				if err != nil {
					continue
				}
				if err := c.ad.Write(ctx, c.instIdx(wr), m); err != nil {
					continue
				}
				mu.Lock()
				res.Concurrent++
				mu.Unlock()
				verify(m, nowMS(), "concurrent")
			}
		}(w)
	}
	wg.Wait()

	c.res.AckCheck = &res
	if !res.Passed {
		return fmt.Errorf("ACK CHECK FAILED: %d of %d acknowledgements were not visible to a second session immediately afterwards (%s)",
			res.Violations, res.Attempts, strings.Join(res.Examples[:minInt(len(res.Examples), 2)], "; "))
	}
	return nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AckCheck reports the acknowledgement-visibility experiment.
type AckCheck struct {
	Sequential int `json:"sequential_writes"`
	Concurrent int `json:"concurrent_writes"`
	Attempts   int `json:"verifications"`
	Violations int `json:"violations"`
	// The two counts are the whole point of the experiment: a violation with ONE
	// writer means the acknowledgement is issued before the commit is visible; a
	// violation only under concurrency means the ledger's history order and the
	// database's commit order disagree.
	SequentialViolations int      `json:"violations_sequential"`
	ConcurrentViolations int      `json:"violations_concurrent"`
	Errors               int      `json:"verification_errors"`
	Examples             []string `json:"examples,omitempty"`
	Passed               bool     `json:"passed"`
}
