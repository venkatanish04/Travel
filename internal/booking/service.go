package booking

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"travelraft/internal/raft"
	"travelraft/internal/state"
)

type Service struct {
	node     *raft.Node
	machine  *state.Machine
	sequence uint64
}

func NewService(nodeID string) *Service {
	machine := state.NewMachine()
	return &Service{node: raft.NewNode(nodeID, machine), machine: machine}
}
func (s *Service) Node() *raft.Node { return s.node }

func (s *Service) Reserve(customerID, tripID, seatID string) (state.Reservation, error) {
	reservation := state.Reservation{ID: fmt.Sprintf("r-%d", atomic.AddUint64(&s.sequence, 1)), CustomerID: customerID, TripID: tripID, SeatID: seatID, Status: "confirmed"}
	command, err := json.Marshal(state.Command{Type: "reserve", Reservation: reservation})
	if err != nil {
		return state.Reservation{}, err
	}
	if _, err = s.node.Apply(command); err != nil {
		return state.Reservation{}, err
	}
	return reservation, nil
}

func (s *Service) Cancel(id string) (state.Reservation, error) {
	command, err := json.Marshal(state.Command{Type: "cancel", ReservationID: id})
	if err != nil {
		return state.Reservation{}, err
	}
	if _, err = s.node.Apply(command); err != nil {
		return state.Reservation{}, err
	}
	reservation, ok := s.machine.Get(id)
	if !ok {
		return state.Reservation{}, fmt.Errorf("reservation %s not found", id)
	}
	return reservation, nil
}

func (s *Service) Get(id string) (state.Reservation, bool) { return s.machine.Get(id) }
func (s *Service) List() []state.Reservation               { return s.machine.List() }
