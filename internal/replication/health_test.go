package replication

import (
	"sync"
	"testing"
	"time"
)

// =============================================================================
// WHAT WE'RE TESTING: Pure logic in HealthMonitor
// WHY TESTABLE: No I/O, no goroutines, deterministic behavior
// =============================================================================

func TestHealthMonitor_InitialState(t *testing.T) {
	h := NewHealthMonitor(5*time.Second, 3)

	if !h.IsHealthy() {
		t.Error("new HealthMonitor should be healthy")
	}
}

func TestHealthMonitor_FailureThreshold(t *testing.T) {
	h := NewHealthMonitor(5*time.Second, 3) // 3 failures = unhealthy

	// 1 failure - still healthy
	h.RecordFailure()
	if !h.IsHealthy() {
		t.Error("should be healthy with 1 failure")
	}

	// 2 failures - still healthy
	h.RecordFailure()
	if !h.IsHealthy() {
		t.Error("should be healthy with 2 failures")
	}

	// 3 failures - NOW unhealthy (>= maxMissedHeartbeats)
	h.RecordFailure()
	if h.IsHealthy() {
		t.Error("should be unhealthy with 3 failures")
	}

	// 4 failures - still unhealthy
	h.RecordFailure()
	if h.IsHealthy() {
		t.Error("should still be unhealthy with 4 failures")
	}
}

func TestHealthMonitor_SuccessResetsFailures(t *testing.T) {
	h := NewHealthMonitor(5*time.Second, 3)

	// Accumulate 2 failures
	h.RecordFailure()
	h.RecordFailure()

	// Success resets counter
	h.RecordSuccess()

	// Should be healthy again
	if !h.IsHealthy() {
		t.Error("should be healthy after success")
	}

	// Now need 3 more failures to become unhealthy
	h.RecordFailure()
	h.RecordFailure()
	if !h.IsHealthy() {
		t.Error("should still be healthy with 2 failures after reset")
	}

	h.RecordFailure()
	if h.IsHealthy() {
		t.Error("should be unhealthy with 3 failures")
	}
}

func TestHealthMonitor_RecoveryFromUnhealthy(t *testing.T) {
	h := NewHealthMonitor(5*time.Second, 3)

	// Become unhealthy
	h.RecordFailure()
	h.RecordFailure()
	h.RecordFailure()
	if h.IsHealthy() {
		t.Fatal("should be unhealthy")
	}

	// Single success recovers
	h.RecordSuccess()
	if !h.IsHealthy() {
		t.Error("should recover after success")
	}
}

func TestHealthMonitor_DifferentThresholds(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		failures  int
		healthy   bool
	}{
		{"threshold 1, 0 failures", 1, 0, true},
		{"threshold 1, 1 failure", 1, 1, false},
		{"threshold 5, 4 failures", 5, 4, true},
		{"threshold 5, 5 failures", 5, 5, false},
		{"threshold 10, 9 failures", 10, 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHealthMonitor(5*time.Second, tt.threshold)
			for i := 0; i < tt.failures; i++ {
				h.RecordFailure()
			}
			if h.IsHealthy() != tt.healthy {
				t.Errorf("expected healthy=%v, got %v", tt.healthy, h.IsHealthy())
			}
		})
	}
}

// =============================================================================
// CONCURRENCY TEST: Verify thread-safety
// WHY TESTABLE: We can spawn goroutines and check for races
// RUN WITH: go test -race ./internal/replication/...
// =============================================================================

func TestHealthMonitor_ConcurrentAccess(t *testing.T) {
	h := NewHealthMonitor(5*time.Second, 100) // High threshold so we don't hit it

	var wg sync.WaitGroup
	iterations := 1000

	// Spawn writers (RecordFailure)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			h.RecordFailure()
		}
	}()

	// Spawn writers (RecordSuccess)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			h.RecordSuccess()
		}
	}()

	// Spawn readers (IsHealthy)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = h.IsHealthy()
		}
	}()

	wg.Wait()
	// If we get here without race detector complaining, locks are working
}

// =============================================================================
// WHAT WE'RE NOT TESTING HERE (and why):
//
// 1. StartHeartbeatForSlave() - Uses time.Ticker, time.After, goroutines
//    WHY HARD: Would need to mock time or wait real seconds
//    ALTERNATIVE: Integration test (TestMasterSlaveReplication)
//
// 2. ListenForPongs() - Blocks on bufio.Scanner
//    WHY HARD: Needs real TCP connection or mock net.Conn
//    ALTERNATIVE: Integration test
//
// 3. broadcast() - Spawns goroutines per slave
//    WHY HARD: Need to verify async behavior
//    ALTERNATIVE: Integration test (TestMultipleSlaves)
//
// 4. Actual PING/PONG over network
//    WHY HARD: Timing, network, multiple goroutines
//    ALTERNATIVE: Manual test (we did this!) or integration test
// =============================================================================

