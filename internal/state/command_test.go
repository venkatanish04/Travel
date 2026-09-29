package state

import "testing"

type recordingApplier struct {
	booked    bool
	cancelled bool
}

func (a *recordingApplier) ApplyBooking(string, int, int, int, string) error {
	a.booked = true
	return nil
}

func (a *recordingApplier) ApplyCancellation(string) error {
	a.cancelled = true
	return nil
}

func TestMySQLMachineAppliesBookingCommands(t *testing.T) {
	applier := &recordingApplier{}
	machine := NewMySQLMachine(applier)
	data, err := EncodeBookingCommand(BookingCommand{
		Type: CommandBook, PNR: "TRF-1", UserID: 1, VehicleID: 2, SeatID: 3, PassengerName: "Venkat",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := machine.Apply(data); err != nil {
		t.Fatal(err)
	}
	if !applier.booked {
		t.Fatal("booking command was not applied")
	}
}

func TestMySQLMachineAppliesCancellationCommands(t *testing.T) {
	applier := &recordingApplier{}
	machine := NewMySQLMachine(applier)
	data, err := EncodeBookingCommand(BookingCommand{Type: CommandCancel, PNR: "TRF-1"})
	if err != nil {
		t.Fatal(err)
	}
	if err := machine.Apply(data); err != nil {
		t.Fatal(err)
	}
	if !applier.cancelled {
		t.Fatal("cancellation command was not applied")
	}
}
