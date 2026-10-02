package booking

import (
	"testing"

	"travelraft/internal/storage"
)

func TestDatabaseBookingLifecycle(t *testing.T) {
	config := storage.MySQLConfigFromEnv()
	db, err := storage.NewDatabase(config.Username, config.Password, config.Host, config.Port, config.Database)
	if err != nil {
		t.Skipf("MySQL is unavailable: %v", err)
	}
	defer db.Close()
	if err := storage.InitializeMySQLSchema(db); err != nil {
		t.Fatal(err)
	}
	if err := SeedDatabase(db); err != nil {
		t.Fatal(err)
	}

	service := NewDatabaseService(db)
	userID, err := service.RegisterUser("Venkat", "venkat@example.com")
	if err != nil {
		t.Fatal(err)
	}
	vehicles, err := service.Search("TRAIN", "Vijayawada", "Hyderabad")
	if err != nil || len(vehicles) != 1 {
		t.Fatalf("search = %d vehicles, err = %v", len(vehicles), err)
	}
	seats, err := service.Seats(vehicles[0].ID)
	if err != nil || len(seats) != 20 {
		t.Fatalf("seats = %d, err = %v", len(seats), err)
	}
	var availableSeat storage.Seat
	for _, seat := range seats {
		if seat.Status == "AVAILABLE" {
			availableSeat = seat
			break
		}
	}
	if availableSeat.ID == 0 {
		t.Fatal("no available seat")
	}
	pnr, err := service.Book(int(userID), vehicles[0].ID, availableSeat.ID, "Venkat")
	if err != nil {
		t.Fatal(err)
	}
	booking, err := service.GetBooking(pnr)
	if err != nil || booking.Status != "CONFIRMED" || booking.SeatNumber != availableSeat.SeatNumber {
		t.Fatalf("booking = %#v, err = %v", booking, err)
	}
	if err := service.Cancel(pnr); err != nil {
		t.Fatal(err)
	}
	booking, err = service.GetBooking(pnr)
	if err != nil || booking.Status != "CANCELLED" {
		t.Fatalf("cancelled booking = %#v, err = %v", booking, err)
	}
}
