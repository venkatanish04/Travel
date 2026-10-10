package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	pb "travelraft/api/proto"
	"travelraft/internal/booking"
)

type Server struct {
	booking *booking.GRPCService
}

func NewServer(service *booking.GRPCService) http.Handler {
	server := &Server{booking: service}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", server.search)
	mux.HandleFunc("/api/seats", server.seats)
	mux.HandleFunc("/api/register", server.register)
	mux.HandleFunc("/api/book", server.book)
	mux.HandleFunc("/api/booking", server.bookingByPNR)
	return withCORS(mux)
}

type searchRequest struct {
	VehicleType string `json:"vehicleType"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type registerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type bookRequest struct {
	PassengerName string `json:"passengerName"`
	Email         string `json:"email"`
	VehicleID     int32  `json:"vehicleId"`
	SeatID        int32  `json:"seatId"`
}

func decode(w http.ResponseWriter, r *http.Request, value any) bool {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var request searchRequest
	if !decode(w, r, &request) {
		return
	}
	response, err := s.booking.Search(r.Context(), &pb.SearchRequest{
		VehicleType: request.VehicleType,
		Source:      request.Source,
		Destination: request.Destination,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) seats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var vehicleID int32
	if _, err := fmt.Sscan(r.URL.Query().Get("vehicleId"), &vehicleID); err != nil {
		writeError(w, http.StatusBadRequest, "vehicleId is required")
		return
	}
	response, err := s.booking.GetSeats(r.Context(), &pb.GetSeatsRequest{VehicleId: vehicleID})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var request registerRequest
	if !decode(w, r, &request) {
		return
	}
	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Email) == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}
	id, err := s.booking.Service.RegisterUser(request.Name, request.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "name": request.Name, "email": request.Email})
}

func (s *Server) book(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var request bookRequest
	if !decode(w, r, &request) {
		return
	}
	response, err := s.booking.Book(r.Context(), &pb.BookRequest{
		PassengerName: request.PassengerName,
		Email:         request.Email,
		VehicleId:     request.VehicleID,
		SeatId:        request.SeatID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) bookingByPNR(w http.ResponseWriter, r *http.Request) {
	pnr := strings.TrimSpace(r.URL.Query().Get("pnr"))
	if pnr == "" {
		writeError(w, http.StatusBadRequest, "pnr is required")
		return
	}
	if r.Method == http.MethodGet {
		response, err := s.booking.GetBooking(r.Context(), &pb.GetBookingRequest{Pnr: pnr})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, response)
		return
	}
	if r.Method == http.MethodDelete {
		response, err := s.booking.Cancel(r.Context(), &pb.CancelRequest{Pnr: pnr})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, response)
		return
	}
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
