package booking

import (
	"database/sql"
	"fmt"

	"travelraft/internal/storage"
)

func SeedDatabase(db *sql.DB) error {
	vehicles, err := storage.GetAllVehicles(db)
	if err != nil {
		return err
	}
	if len(vehicles) > 0 {
		return nil
	}

	data := []struct {
		vehicleType string
		number      string
		name        string
		source      string
		destination string
		departure   string
	}{
		{"TRAIN", "TR101", "Vande Bharat Express", "Vijayawada", "Hyderabad", "06:00"},
		{"TRAIN", "TR102", "Intercity Express", "Vijayawada", "Chennai", "08:30"},
		{"BUS", "BU201", "APSRTC Super Luxury", "Vijayawada", "Hyderabad", "07:00"},
		{"BUS", "BU202", "Garuda Express", "Vijayawada", "Visakhapatnam", "09:30"},
		{"TRAIN", "TR103", "Charminar Express", "Hyderabad", "Chennai", "16:00"},
		{"TRAIN", "TR104", "Godavari Express", "Rajahmundry", "Vijayawada", "05:45"},
		{"TRAIN", "TR105", "Coromandel Express", "Chennai", "Visakhapatnam", "07:15"},
		{"TRAIN", "TR106", "Kacheguda Express", "Hyderabad", "Bengaluru", "20:30"},
		{"BUS", "BU203", "Orange Sleeper", "Hyderabad", "Bengaluru", "21:00"},
		{"BUS", "BU204", "Kaveri Express", "Chennai", "Bengaluru", "10:30"},
		{"BUS", "BU205", "Simhapuri Travels", "Nellore", "Vijayawada", "06:45"},
		{"BUS", "BU206", "Metro Luxury", "Visakhapatnam", "Rajahmundry", "08:15"},
		{"FLIGHT", "FL301", "Domestic Flight", "Vijayawada", "Delhi", "06:30"},
		{"FLIGHT", "FL302", "Domestic Flight", "Vijayawada", "Hyderabad", "14:00"},
	}

	for _, item := range data {
		vehicle, err := storage.GetVehicleByNumber(db, item.number)
		if err == sql.ErrNoRows {
			id, addErr := storage.AddVehicle(db, item.vehicleType, item.number, item.name, item.source, item.destination, item.departure)
			if addErr != nil {
				return addErr
			}
			vehicle = &storage.Vehicle{ID: int(id)}
		} else if err != nil {
			return err
		}
		for seat := 1; seat <= 20; seat++ {
			seatNumber := fmt.Sprintf("A%d", seat)
			if err := storage.AddSeatIfMissing(db, int64(vehicle.ID), seatNumber); err != nil {
				return err
			}
		}
	}
	return nil
}
