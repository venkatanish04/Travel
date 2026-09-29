package tests

import (
	"testing"
	"travelraft/internal/booking"
)

func TestBookingServiceNodeCommitsLog(t *testing.T) {
	service := booking.NewService("integration-node")
	if _, err := service.Reserve("alice", "flight-1", "12A"); err != nil {
		t.Fatal(err)
	}
	if got := len(service.Node().Entries()); got != 1 {
		t.Fatalf("log entries = %d, want 1", got)
	}
}
