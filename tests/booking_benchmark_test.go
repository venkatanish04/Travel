package tests

import (
	"testing"

	"travelraft/internal/booking"
)

func BenchmarkReservationCommand(b *testing.B) {
	for b.Loop() {
		service := booking.NewService("benchmark-node")
		if _, err := service.Reserve("customer", "trip", "seat"); err != nil {
			b.Fatal(err)
		}
	}
}
