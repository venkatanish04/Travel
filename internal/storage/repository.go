package storage

import (
	"database/sql"
	"fmt"
)

func InitializeMySQLSchema(db *sql.DB) error {
	statements := []string{`
CREATE TABLE IF NOT EXISTS users (
	 id BIGINT AUTO_INCREMENT PRIMARY KEY,
	 name VARCHAR(100) NOT NULL,
	 email VARCHAR(150) NOT NULL UNIQUE,
	 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`, `
CREATE TABLE IF NOT EXISTS vehicles (
	 id BIGINT AUTO_INCREMENT PRIMARY KEY,
	 vehicle_type VARCHAR(20) NOT NULL,
	 vehicle_number VARCHAR(50) NOT NULL UNIQUE,
	 name VARCHAR(150) NOT NULL,
	 source VARCHAR(100) NOT NULL,
	 destination VARCHAR(100) NOT NULL,
	 departure_time VARCHAR(32) NOT NULL,
	 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`, `
CREATE TABLE IF NOT EXISTS seats (
	 id BIGINT AUTO_INCREMENT PRIMARY KEY,
	 vehicle_id BIGINT NOT NULL,
	 seat_number VARCHAR(20) NOT NULL,
	 status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
	 CONSTRAINT fk_seats_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE,
	 CONSTRAINT unique_vehicle_seat UNIQUE (vehicle_id, seat_number)
);`, `
CREATE TABLE IF NOT EXISTS bookings (
	 id BIGINT AUTO_INCREMENT PRIMARY KEY,
	 pnr VARCHAR(100) NOT NULL UNIQUE,
	 user_id BIGINT NOT NULL,
	 vehicle_id BIGINT NOT NULL,
	 seat_id BIGINT NOT NULL,
	 passenger_name VARCHAR(150) NOT NULL,
	 status VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED',
	 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	 CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id),
	 CONSTRAINT fk_booking_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
	 CONSTRAINT fk_booking_seat FOREIGN KEY (seat_id) REFERENCES seats(id)
);`}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("initialize MySQL schema: %w", err)
		}
	}
	return nil
}
