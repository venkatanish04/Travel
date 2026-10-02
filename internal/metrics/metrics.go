package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                  sync.RWMutex
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	TotalLatency        time.Duration
	RaftReplicationTime time.Duration
	DBOperationTime     time.Duration
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordRequest(success bool, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++
	m.TotalLatency += latency
	if success {
		m.SuccessfulRequests++
	} else {
		m.FailedRequests++
	}
}

func (m *Metrics) RecordRaftReplication(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RaftReplicationTime += duration
}

func (m *Metrics) RecordDBOperation(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DBOperationTime += duration
}

func (m *Metrics) AverageLatency() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.TotalRequests == 0 {
		return 0
	}
	return m.TotalLatency / time.Duration(m.TotalRequests)
}
