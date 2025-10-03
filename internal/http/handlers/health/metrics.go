package health

import (
	"sync/atomic"
	"time"
)

// Metrics captures application level telemetry useful for health responses.
type Metrics struct {
	startedAt time.Time
	requests  atomic.Uint64
}

// NewMetrics constructs a metrics tracker initialised with the provided start time.
func NewMetrics(startedAt time.Time) *Metrics {
	return &Metrics{startedAt: startedAt}
}

// IncrementRequests increases the handled request counter.
func (m *Metrics) IncrementRequests() {
	if m == nil {
		return
	}
	m.requests.Add(1)
}

// RequestCount returns the total number of handled HTTP requests.
func (m *Metrics) RequestCount() uint64 {
	if m == nil {
		return 0
	}
	return m.requests.Load()
}

// StartedAt exposes when the application began serving traffic.
func (m *Metrics) StartedAt() time.Time {
	if m == nil {
		return time.Time{}
	}
	return m.startedAt
}

// Uptime reports how long the application has been running.
func (m *Metrics) Uptime() time.Duration {
	if m == nil || m.startedAt.IsZero() {
		return 0
	}
	return time.Since(m.startedAt)
}
