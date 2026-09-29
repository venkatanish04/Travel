package tests

import (
	"testing"
	"travelraft/internal/booking"
)

func TestReserveAndCancel(t *testing.T) {
	service := booking.NewService("test-node")
	reservation, err := service.Reserve("alice", "train-1", "A1")
	if err != nil {
		t.Fatal(err)
	}
	if reservation.Status != "confirmed" {
		t.Fatalf("status = %q", reservation.Status)
	}
	cancelled, err := service.Cancel(reservation.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("status = %q", cancelled.Status)
	}
}

func TestDuplicateSeatRejected(t *testing.T) {
	service := booking.NewService("test-node")
	if _, err := service.Reserve("alice", "bus-1", "B2"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Reserve("bob", "bus-1", "B2"); err == nil {
		t.Fatal("expected duplicate seat error")
	}
}
