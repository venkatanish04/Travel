package storage

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type Vehicle struct {
	ID            int
	VehicleType   string
	VehicleNumber string
	Name          string
	Source        string
	Destination   string
	DepartureTime string
}

type Seat struct {
	ID         int
	VehicleID  int
	SeatNumber string
	Status     string
}

type Booking struct {
	ID            int
	PNR           string
	UserID        int
	VehicleID     int
	SeatID        int
	PassengerName string
	Status        string
	VehicleName   string
	VehicleNumber string
	SeatNumber    string
}

func AddUser(db *sql.DB, name, email string) (int64, error) {
	result, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, name, email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	err := db.QueryRow(`SELECT id, name, email FROM users WHERE email = ?`, email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func AddVehicle(db *sql.DB, vehicleType, vehicleNumber, name, source, destination, departureTime string) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO vehicles (vehicle_type, vehicle_number, name, source, destination, departure_time)
		VALUES (?, ?, ?, ?, ?, ?)`, vehicleType, vehicleNumber, name, source, destination, departureTime)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func AddSeat(db *sql.DB, vehicleID int64, seatNumber string) error {
	_, err := db.Exec(`INSERT INTO seats (vehicle_id, seat_number, status) VALUES (?, ?, 'AVAILABLE')`, vehicleID, seatNumber)
	return err
}

func GetVehicles(db *sql.DB, vehicleType, source, destination string) ([]Vehicle, error) {
	rows, err := db.Query(`
		SELECT id, vehicle_type, vehicle_number, name, source, destination, departure_time
		FROM vehicles WHERE vehicle_type = ? AND source = ? AND destination = ? ORDER BY departure_time`, vehicleType, source, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVehicles(rows)
}

func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	rows, err := db.Query(`
		SELECT id, vehicle_type, vehicle_number, name, source, destination, departure_time
		FROM vehicles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVehicles(rows)
}

func scanVehicles(rows *sql.Rows) ([]Vehicle, error) {
	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.VehicleType, &vehicle.VehicleNumber, &vehicle.Name, &vehicle.Source, &vehicle.Destination, &vehicle.DepartureTime); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, rows.Err()
}

func GetSeats(db *sql.DB, vehicleID int) ([]Seat, error) {
	rows, err := db.Query(`SELECT id, vehicle_id, seat_number, status FROM seats WHERE vehicle_id = ? ORDER BY id`, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []Seat
	for rows.Next() {
		var seat Seat
		if err := rows.Scan(&seat.ID, &seat.VehicleID, &seat.SeatNumber, &seat.Status); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, rows.Err()
}

func CreateBooking(db *sql.DB, pnr string, userID, vehicleID, seatID int, passengerName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`UPDATE seats SET status = 'BOOKED' WHERE id = ? AND vehicle_id = ? AND status = 'AVAILABLE'`, seatID, vehicleID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("seat is already booked or does not belong to vehicle")
	}

	if _, err := tx.Exec(`
		INSERT INTO bookings (pnr, user_id, vehicle_id, seat_id, passenger_name, status)
		VALUES (?, ?, ?, ?, ?, 'CONFIRMED')`, pnr, userID, vehicleID, seatID, passengerName); err != nil {
		return err
	}
	return tx.Commit()
}

func GetBookingByPNR(db *sql.DB, pnr string) (*Booking, error) {
	var booking Booking
	err := db.QueryRow(`
		SELECT b.id, b.pnr, b.user_id, b.vehicle_id, b.seat_id, b.passenger_name, b.status,
		       v.name, v.vehicle_number, s.seat_number
		FROM bookings b JOIN vehicles v ON b.vehicle_id = v.id JOIN seats s ON b.seat_id = s.id
		WHERE b.pnr = ?`, pnr).Scan(
		&booking.ID, &booking.PNR, &booking.UserID, &booking.VehicleID, &booking.SeatID,
		&booking.PassengerName, &booking.Status, &booking.VehicleName, &booking.VehicleNumber, &booking.SeatNumber)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func CancelBooking(db *sql.DB, pnr string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var seatID int
	var status string
	if err := tx.QueryRow(`SELECT seat_id, status FROM bookings WHERE pnr = ?`, pnr).Scan(&seatID, &status); err != nil {
		return err
	}
	if status == "CANCELLED" {
		return fmt.Errorf("booking is already cancelled")
	}
	if _, err := tx.Exec(`UPDATE bookings SET status = 'CANCELLED' WHERE pnr = ?`, pnr); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE seats SET status = 'AVAILABLE' WHERE id = ?`, seatID); err != nil {
		return err
	}
	return tx.Commit()
}
