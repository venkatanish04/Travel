package tests

import (
	"testing"
	"time"

	"travelraft/internal/booking"
)

func TestConcurrentPerformance(t *testing.T) {
	const workers = 20
	service := booking.NewService("performance-node")
	start := time.Now()
	results := make(chan error, workers)

	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			requestStart := time.Now()
			_, err := service.Reserve("customer", "performance-trip", string(rune('A'+worker)))
			t.Logf("worker %d latency: %s", worker, time.Since(requestStart))
			results <- err
		}(worker)
	}

	for worker := 0; worker < workers; worker++ {
		if err := <-results; err != nil {
			t.Fatalf("worker %d failed: %v", worker, err)
		}
	}
	t.Logf("total execution time: %s", time.Since(start))
}
