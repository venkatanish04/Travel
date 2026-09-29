package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "travelraft/api/proto"
	"travelraft/internal/booking"
	"travelraft/internal/raft"
	"travelraft/internal/state"
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

	databaseConfig := storage.MySQLConfigFromEnv()
	db, err := storage.NewDatabase(
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := storage.InitializeMySQLSchema(db); err != nil {
		log.Fatal(err)
	}
	if err := booking.SeedDatabase(db); err != nil {
		log.Fatal(err)
	}

	node := raft.NewNode(config.ID, config.Address, config.Peers)
	node.SetStateMachine(state.NewMySQLMachine(storage.NewBookingStore(db)))
	node.StartElectionLoop()
	node.StartLeaderLoop()
	node.StartApplyLoop(context.Background())
	defer node.Stop()
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, booking.NewGRPCServer(booking.NewDatabaseService(db), db, node))
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
