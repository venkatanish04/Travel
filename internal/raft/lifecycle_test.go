package raft

import (
	"context"
	"net"
	"testing"
	"time"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
)

func TestSingleNodeStartsElectionAutomatically(t *testing.T) {
	node := NewNode("node1", "localhost:50051", map[string]string{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node.Start(ctx)
	defer node.Stop()

	deadline := time.After(5 * time.Second)
	for node.GetState() != Leader {
		select {
		case <-deadline:
			t.Fatal("node did not become leader before election timeout")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestThreeNodeClusterElectsAndReplicates(t *testing.T) {
	type runningNode struct {
		node   *Node
		server *grpc.Server
		stop   func()
	}
	addresses := make([]string, 3)
	listeners := make([]net.Listener, 3)
	for index := range listeners {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		listeners[index] = listener
		addresses[index] = listener.Addr().String()
	}

	ids := []string{"node1", "node2", "node3"}
	running := make([]runningNode, 3)
	for index, id := range ids {
		peers := make(map[string]string)
		for peerIndex, peerID := range ids {
			if peerIndex != index {
				peers[peerID] = addresses[peerIndex]
			}
		}
		node := NewNode(id, addresses[index], peers)
		server := grpc.NewServer()
		pb.RegisterRaftServiceServer(server, NewGRPCServer(node))
		go server.Serve(listeners[index])
		node.Start(context.Background())
		running[index] = runningNode{node: node, server: server, stop: func() { node.Stop(); server.Stop(); listeners[index].Close() }}
	}
	defer func() {
		for _, item := range running {
			item.stop()
		}
	}()

	var leader *Node
	deadline := time.After(8 * time.Second)
	for leader == nil {
		for _, item := range running {
			if item.node.GetState() == Leader {
				if leader != nil && leader != item.node {
					t.Fatal("cluster elected multiple leaders")
				}
				leader = item.node
			}
		}
		if leader != nil {
			break
		}
		select {
		case <-deadline:
			t.Fatal("cluster did not elect a leader")
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}

	entry, err := leader.Replicate([]byte("ADD_TRAIN"))
	if err != nil {
		t.Fatalf("replication failed: %v", err)
	}
	if leader.CommitIndex != 0 || entry.Index != 1 {
		t.Fatalf("leader commit=%d entry=%d, want 0 and 1", leader.CommitIndex, entry.Index)
	}
	for _, item := range running {
		if len(item.node.Entries()) != 1 {
			t.Fatalf("node %s has %d entries, want 1", item.node.ID(), len(item.node.Entries()))
		}
	}
}
