package tests

import (
	"sync"
	"testing"
	"travelraft/internal/booking"
)

func TestConcurrentReservationsUseUniqueIDs(t *testing.T) {
	service := booking.NewService("test-node")
	const count = 20
	var wait sync.WaitGroup
	ids := make(chan string, count)
	for i := 0; i < count; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			reservation, err := service.Reserve("customer", "trip", string(rune('A'+index)))
			if err == nil {
				ids <- reservation.ID
			}
		}(i)
	}
	wait.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatalf("duplicate id %s", id)
		}
		seen[id] = true
	}
	if len(seen) != count {
		t.Fatalf("reserved %d seats, want %d", len(seen), count)
	}
}
