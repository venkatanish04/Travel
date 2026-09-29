package raft

import (
	"net"
	"testing"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
)

func TestStartElectionWinsMajorityOverGRPC(t *testing.T) {
	followerA, stopA := startTestRaftServer(t, "node2")
	defer stopA()
	followerB, stopB := startTestRaftServer(t, "node3")
	defer stopB()

	candidate := NewNode("node1", "", map[string]string{
		"node2": followerA,
		"node3": followerB,
	})
	state := candidate.StartElection()
	if state.Role != Leader {
		t.Fatalf("state = %s, want %s", state.Role, Leader)
	}
	if state.CurrentTerm != 1 {
		t.Fatalf("term = %d, want 1", state.CurrentTerm)
	}
}

func startTestRaftServer(t *testing.T, id string) (string, func()) {
	t.Helper()
	node := NewNode(id, "", map[string]string{})
	server := grpc.NewServer()
	pb.RegisterRaftServiceServer(server, NewGRPCServer(node))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go server.Serve(listener)
	return listener.Addr().String(), func() {
		server.Stop()
		listener.Close()
	}
}
