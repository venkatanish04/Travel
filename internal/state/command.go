package state

import "encoding/json"

type BookingCommandType string

const (
	CommandBook   BookingCommandType = "BOOK"
	CommandCancel BookingCommandType = "CANCEL"
)

type BookingCommand struct {
	Type          BookingCommandType `json:"type"`
	PNR           string             `json:"pnr"`
	UserID        int                `json:"user_id"`
	VehicleID     int                `json:"vehicle_id"`
	SeatID        int                `json:"seat_id"`
	PassengerName string             `json:"passenger_name,omitempty"`
}

func EncodeBookingCommand(command BookingCommand) ([]byte, error) {
	return json.Marshal(command)
}

func DecodeBookingCommand(data []byte) (BookingCommand, error) {
	var command BookingCommand
	err := json.Unmarshal(data, &command)
	return command, err
}
