package tests

import (
	"sync"
	"testing"

	"travelraft/internal/booking"
)

func TestConcurrentReservationsRejectDoubleBooking(t *testing.T) {
	service := booking.NewService("double-booking-node")
	const attempts = 10

	var wait sync.WaitGroup
	results := make(chan error, attempts)
	for index := 0; index < attempts; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			_, err := service.Reserve("customer", "train-1", "A5")
			results <- err
		}(index)
	}
	wait.Wait()
	close(results)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful bookings for one seat = %d, want 1", successes)
	}
}
