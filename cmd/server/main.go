package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "travelraft/api/proto"
	"travelraft/config"
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
	configMap := raft.ConfigsForEnvironment(os.Getenv("TRAVELRAFT_ENV"))
	nodeConfig, ok := configMap[*nodeID]
	if !ok {
		log.Fatalf("unknown node ID %q", *nodeID)
	}
	listenAddress := ":" + *port
	environment := os.Getenv("TRAVELRAFT_ENV")
	if environment == "" || environment == "local" {
		nodeConfig.Address = "localhost:" + *port
	} else if environment == "docker" {
		nodeConfig.Address = nodeConfig.ID + ":" + *port
	}

	databaseConfig := config.LoadDatabaseConfig()
	db, err := storage.NewDatabase(
		databaseConfig.User,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Name,
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

	node := raft.NewNode(nodeConfig.ID, nodeConfig.Address, nodeConfig.Peers)
	node.SetStateMachine(state.NewMySQLMachine(storage.NewBookingStore(db)))
	node.StartElectionLoop()
	node.StartLeaderLoop()
	node.StartApplyLoop(context.Background())
	defer node.Stop()
	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, booking.NewGRPCServer(booking.NewDatabaseService(db), db, node))
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
	fmt.Println("Node:", nodeConfig.ID)
	fmt.Println("Address:", nodeConfig.Address)
	fmt.Println()
	fmt.Println("Server Status: RUNNING")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
