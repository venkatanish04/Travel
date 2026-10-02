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

func EncodeCommand(command interface{}) ([]byte, error) {
	switch typed := command.(type) {
	case BookingCommand:
		return EncodeBookingCommand(typed)
	case Command:
		return json.Marshal(typed)
	default:
		return json.Marshal(command)
	}
}

func DecodeBookingCommand(data []byte) (BookingCommand, error) {
	var command BookingCommand
	err := json.Unmarshal(data, &command)
	return command, err
}

func DecodeCommand(data []byte) (interface{}, error) {
	var booking BookingCommand
	if err := json.Unmarshal(data, &booking); err == nil && (booking.Type == CommandBook || booking.Type == CommandCancel || booking.PNR != "" || booking.UserID != 0 || booking.VehicleID != 0 || booking.SeatID != 0 || booking.PassengerName != "") {
		return booking, nil
	}
	var legacy Command
	if err := json.Unmarshal(data, &legacy); err == nil {
		return legacy, nil
	}
	return nil, json.Unmarshal(data, &struct{}{})
}
