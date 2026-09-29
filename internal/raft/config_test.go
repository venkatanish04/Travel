package raft

import "testing"

func TestDefaultConfigs(t *testing.T) {
	configs := DefaultConfigs()
	if len(configs) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(configs))
	}
	if configs["node1"].Address != "localhost:50051" {
		t.Fatalf("unexpected node1 address: %s", configs["node1"].Address)
	}
	if len(configs["node1"].Peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(configs["node1"].Peers))
	}
}
