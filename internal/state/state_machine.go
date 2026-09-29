package state

import (
	"encoding/json"
	"errors"
	"sync"
)

type Reservation struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	TripID     string `json:"trip_id"`
	SeatID     string `json:"seat_id"`
	Status     string `json:"status"`
}

type Command struct {
	Type          string      `json:"type"`
	Reservation   Reservation `json:"reservation"`
	ReservationID string      `json:"reservation_id"`
}

type Machine struct {
	mu           sync.RWMutex
	reservations map[string]Reservation
	occupied     map[string]string
}

func NewMachine() *Machine {
	return &Machine{reservations: make(map[string]Reservation), occupied: make(map[string]string)}
}

func (m *Machine) Apply(data []byte) error {
	var command Command
	if err := json.Unmarshal(data, &command); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	switch command.Type {
	case "reserve":
		key := command.Reservation.TripID + ":" + command.Reservation.SeatID
		if _, exists := m.occupied[key]; exists {
			return errors.New("seat is already reserved")
		}
		m.reservations[command.Reservation.ID] = command.Reservation
		m.occupied[key] = command.Reservation.ID
	case "cancel":
		reservation, exists := m.reservations[command.ReservationID]
		if !exists {
			return errors.New("reservation not found")
		}
		reservation.Status = "cancelled"
		m.reservations[reservation.ID] = reservation
		delete(m.occupied, reservation.TripID+":"+reservation.SeatID)
	default:
		return errors.New("unknown state command")
	}
	return nil
}

func (m *Machine) Get(id string) (Reservation, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.reservations[id]
	return value, ok
}
func (m *Machine) List() []Reservation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Reservation, 0, len(m.reservations))
	for _, value := range m.reservations {
		result = append(result, value)
	}
	return result
}
