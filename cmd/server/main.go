package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"travelraft/internal/booking"
	"travelraft/internal/storage"
)

type reserveRequest struct {
	CustomerID string `json:"customer_id"`
	TripID     string `json:"trip_id"`
	SeatID     string `json:"seat_id"`
}

func main() {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "node1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	service := booking.NewService(nodeID)
	database, err := storage.OpenDatabase("data/travelraft.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if err := storage.InitializeSchema(database); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "leader": service.Node().LeaderID()})
	})
	mux.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, service.List())
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request reserveRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reservation, err := service.Reserve(request.CustomerID, request.TripID, request.SeatID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		writeJSON(w, http.StatusCreated, reservation)
	})
	mux.HandleFunc("/reservations/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/reservations/")
		reservation, err := service.Cancel(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, reservation)
	})
	address := ":" + port
	log.Printf("TravelRaft server %s listening on %s", nodeID, address)
	log.Fatal(http.ListenAndServe(address, mux))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
