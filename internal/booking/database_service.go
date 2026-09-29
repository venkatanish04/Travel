package booking

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"travelraft/internal/storage"
)

type DatabaseService struct {
	DB       *sql.DB
	sequence uint64
}

func NewDatabaseService(db *sql.DB) *DatabaseService {
	return &DatabaseService{DB: db}
}

func (s *DatabaseService) Search(vehicleType, source, destination string) ([]storage.Vehicle, error) {
	return storage.GetVehicles(s.DB, vehicleType, source, destination)
}

func (s *DatabaseService) Seats(vehicleID int) ([]storage.Seat, error) {
	return storage.GetSeats(s.DB, vehicleID)
}

func (s *DatabaseService) RegisterUser(name, email string) (int64, error) {
	user, err := storage.GetUserByEmail(s.DB, email)
	if err == nil {
		return int64(user.ID), nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	return storage.AddUser(s.DB, name, email)
}

func (s *DatabaseService) Book(userID, vehicleID, seatID int, passengerName string) (string, error) {
	pnr := s.NewPNR()
	if err := storage.CreateBooking(s.DB, pnr, userID, vehicleID, seatID, passengerName); err != nil {
		return "", err
	}
	return pnr, nil
}

func (s *DatabaseService) NewPNR() string {
	return fmt.Sprintf("TRF%d%02d", time.Now().UnixNano()%1000000000, atomic.AddUint64(&s.sequence, 1)%100)
}

func (s *DatabaseService) GetBooking(pnr string) (*storage.Booking, error) {
	return storage.GetBookingByPNR(s.DB, pnr)
}

func (s *DatabaseService) Cancel(pnr string) error {
	return storage.CancelBooking(s.DB, pnr)
}
