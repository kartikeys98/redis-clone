package replication

import (
	"sync"
	"time"
)

type HealthMonitor struct {
	lastHeartbeat    time.Time
	mu               sync.RWMutex
	missedHeartbeats int
	threshold        time.Duration
	maxMissedHeartbeats int
}

func NewHealthMonitor(threshold time.Duration, maxMissedHeartbeats int) *HealthMonitor {
	return &HealthMonitor{
		lastHeartbeat:    time.Now(),
		missedHeartbeats: 0,
		threshold:        threshold,
		maxMissedHeartbeats: maxMissedHeartbeats,
		mu:               sync.RWMutex{},
	}
}

func(h *HealthMonitor) IsHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.missedHeartbeats < h.maxMissedHeartbeats //&& h.timeSinceLastHeartbeat() < h.threshold
}

// func(h *HealthMonitor) timeSinceLastHeartbeat() time.Duration { // Why unexported?
// 	return time.Since(h.lastHeartbeat)
// }

func(h *HealthMonitor) RecordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastHeartbeat = time.Now()
	h.missedHeartbeats = 0
}

func(h *HealthMonitor) RecordFailure() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.missedHeartbeats++
}