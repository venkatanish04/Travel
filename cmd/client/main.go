package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	customer := flag.String("customer", "demo", "customer id")
	trip := flag.String("trip", "train-1", "trip id")
	seat := flag.String("seat", "A1", "seat id")
	server := flag.String("server", "http://localhost:8080", "server URL")
	flag.Parse()
	body, _ := json.Marshal(map[string]string{"customer_id": *customer, "trip_id": *trip, "seat_id": *seat})
	response, err := http.Post(*server+"/reservations", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var result any
	_ = json.NewDecoder(response.Body).Decode(&result)
	fmt.Printf("status=%s\n%v\n", response.Status, result)
}
