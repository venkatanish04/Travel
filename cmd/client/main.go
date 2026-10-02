package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddress := flag.String("address", "localhost:50051", "TravelRaft gRPC server address")
	flag.Parse()

	conn, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("connection error:", err)
	}
	defer conn.Close()

	client := pb.NewBookingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("========================================")
	fmt.Println("       TravelRaft gRPC Client")
	fmt.Println("========================================")
	vehicles, err := search(client, ctx, "TRAIN", "Vijayawada", "Hyderabad")
	if err != nil {
		log.Fatal("search error:", err)
	}
	if len(vehicles) == 0 {
		fmt.Println("No transport found.")
		return
	}
	vehicle := vehicles[0]
	fmt.Printf("\nSelected: %s - %s\n", vehicle.GetVehicleNumber(), vehicle.GetName())

	seats, err := getSeats(client, ctx, vehicle.GetId())
	if err != nil {
		log.Fatal("seat error:", err)
	}
	var selectedSeat *pb.Seat
	for _, seat := range seats {
		fmt.Printf("ID: %d | %s | %s\n", seat.GetId(), seat.GetSeatNumber(), seat.GetStatus())
		if selectedSeat == nil && seat.GetStatus() == "AVAILABLE" {
			selectedSeat = seat
		}
	}
	if selectedSeat == nil {
		fmt.Println("No available seats.")
		return
	}

	response, err := client.Book(ctx, &pb.BookRequest{
		PassengerName: "Venkat",
		Email:         "venkat-grpc@example.com",
		VehicleId:     vehicle.GetId(),
		SeatId:        selectedSeat.GetId(),
	})
	if err != nil {
		log.Fatal("booking error:", err)
	}
	if !response.GetSuccess() {
		log.Fatal("booking failed:", response.GetMessage())
	}
	fmt.Println("\n========================================")
	fmt.Println("       BOOKING CONFIRMED")
	fmt.Println("========================================")
	fmt.Println("PNR:", response.GetPnr())
	fmt.Println("Message:", response.GetMessage())

	booking, err := client.GetBooking(ctx, &pb.GetBookingRequest{Pnr: response.GetPnr()})
	if err != nil {
		log.Fatal("lookup error:", err)
	}
	fmt.Println("Lookup found:", booking.GetFound())
}

func search(client pb.BookingServiceClient, ctx context.Context, vehicleType, source, destination string) ([]*pb.Vehicle, error) {
	response, err := client.Search(ctx, &pb.SearchRequest{VehicleType: vehicleType, Source: source, Destination: destination})
	if err != nil {
		return nil, err
	}
	return response.GetVehicles(), nil
}

func getSeats(client pb.BookingServiceClient, ctx context.Context, vehicleID int32) ([]*pb.Seat, error) {
	response, err := client.GetSeats(ctx, &pb.GetSeatsRequest{VehicleId: vehicleID})
	if err != nil {
		return nil, err
	}
	return response.GetSeats(), nil
}
