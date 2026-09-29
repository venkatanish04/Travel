package storage

import (
	"database/sql"
	"fmt"

	"travelraft/internal/state"
)

func InitializeSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS vehicles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	vehicle_type TEXT NOT NULL,
	vehicle_number TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	source TEXT NOT NULL,
	destination TEXT NOT NULL,
	departure_time TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS seats (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	vehicle_id INTEGER NOT NULL,
	seat_number TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'AVAILABLE',
	FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
	UNIQUE(vehicle_id, seat_number)
);
CREATE TABLE IF NOT EXISTS bookings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	pnr TEXT NOT NULL UNIQUE,
	user_id INTEGER NOT NULL,
	vehicle_id INTEGER NOT NULL,
	seat_id INTEGER NOT NULL,
	passenger_name TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'CONFIRMED',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
	FOREIGN KEY (seat_id) REFERENCES seats(id)
);`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("initialize database schema: %w", err)
	}
	return nil
}

type ReservationRepository struct{ database *Database }

func NewReservationRepository(database *Database) *ReservationRepository {
	return &ReservationRepository{database: database}
}
func (r *ReservationRepository) Save(reservations []state.Reservation) error {
	return r.database.Save(reservations)
}
func (r *ReservationRepository) Load() ([]state.Reservation, error) {
	var reservations []state.Reservation
	err := r.database.Load(&reservations)
	return reservations, err
}
