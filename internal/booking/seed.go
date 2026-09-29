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
		{"FLIGHT", "FL301", "Domestic Flight", "Vijayawada", "Delhi", "06:30"},
		{"FLIGHT", "FL302", "Domestic Flight", "Vijayawada", "Hyderabad", "14:00"},
	}

	for _, item := range data {
		id, err := storage.AddVehicle(db, item.vehicleType, item.number, item.name, item.source, item.destination, item.departure)
		if err != nil {
			return err
		}
		for seat := 1; seat <= 20; seat++ {
			if err := storage.AddSeat(db, id, fmt.Sprintf("A%d", seat)); err != nil {
				return err
			}
		}
	}
	return nil
}
