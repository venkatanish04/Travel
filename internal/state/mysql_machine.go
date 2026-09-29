package state

import (
	"fmt"
)

type BookingApplier interface {
	ApplyBooking(pnr string, userID, vehicleID, seatID int, passengerName string) error
	ApplyCancellation(pnr string) error
}

type MySQLMachine struct {
	store BookingApplier
}

func NewMySQLMachine(store BookingApplier) *MySQLMachine {
	return &MySQLMachine{store: store}
}

func (m *MySQLMachine) Apply(data []byte) error {
	command, err := DecodeBookingCommand(data)
	if err != nil {
		return err
	}
	switch command.Type {
	case CommandBook:
		return m.store.ApplyBooking(command.PNR, command.UserID, command.VehicleID, command.SeatID, command.PassengerName)
	case CommandCancel:
		return m.store.ApplyCancellation(command.PNR)
	default:
		return fmt.Errorf("unknown booking command: %s", command.Type)
	}
}
