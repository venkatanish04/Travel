package main

import (
	"fmt"
	"log"
	"net"

	pb "travelraft/api/proto"
	"travelraft/internal/booking"
	"travelraft/internal/storage"

	"google.golang.org/grpc"
)

const grpcPort = ":50051"

func main() {
	fmt.Println("========================================")
	fmt.Println("          TravelRaft Server")
	fmt.Println("========================================")

	db, err := storage.OpenDatabase("data/travelraft.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := storage.InitializeSchema(db); err != nil {
		log.Fatal(err)
	}
	if err := booking.SeedDatabase(db); err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, booking.NewGRPCServer(booking.NewDatabaseService(db), db))
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println()
	fmt.Println("Database: CONNECTED")
	fmt.Println("Travel Data: READY")
	fmt.Println("gRPC Server: LISTENING")
	fmt.Println("Address:", grpcPort)
	fmt.Println()
	fmt.Println("Server Status: RUNNING")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
