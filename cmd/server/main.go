package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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

	nodeID := flag.String("id", "node1", "Raft node ID")
	port := flag.String("port", "50051", "gRPC port")
	flag.Parse()
	config, ok := raft.DefaultConfigs()[*nodeID]
	if !ok {
		log.Fatalf("unknown node ID %q", *nodeID)
	}
	config.Address = "localhost:" + *port

	db, err := storage.OpenDatabase("data/" + config.ID + ".sqlite.db")
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

	node := raft.NewNode(config.ID, config.Address, config.Peers)
	node.StartElectionLoop()
	node.StartLeaderLoop()
	defer node.Stop()
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, booking.NewGRPCServer(booking.NewDatabaseService(db), db))
	pb.RegisterRaftServiceServer(grpcServer, raft.NewGRPCServer(node))
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println()
	fmt.Println("Database: CONNECTED")
	fmt.Println("Travel Data: READY")
	fmt.Println("gRPC Server: LISTENING")
	fmt.Println("Node:", config.ID)
	fmt.Println("Address:", config.Address)
	fmt.Println()
	fmt.Println("Server Status: RUNNING")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
