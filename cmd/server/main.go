package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "travelraft/api/proto"
	"travelraft/internal/booking"
	"travelraft/internal/raft"
	"travelraft/internal/storage"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("          TravelRaft Server")
	fmt.Println("========================================")

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "node1"
	}
	config, ok := raft.DefaultConfigs()[nodeID]
	if !ok {
		log.Fatalf("unknown node ID %q", nodeID)
	}
	listenAddress := config.Address
	if port := os.Getenv("PORT"); port != "" {
		listenAddress = "localhost:" + port
	}

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

	node := raft.NewNode(nodeID, listenAddress, config.Peers)
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, booking.NewGRPCServer(booking.NewDatabaseService(db), db))
	pb.RegisterRaftServiceServer(grpcServer, raft.NewGRPCServer(node))
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println()
	fmt.Println("Database: CONNECTED")
	fmt.Println("Travel Data: READY")
	fmt.Println("gRPC Server: LISTENING")
	fmt.Println("Node:", nodeID)
	fmt.Println("Address:", listenAddress)
	fmt.Println()
	fmt.Println("Server Status: RUNNING")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
