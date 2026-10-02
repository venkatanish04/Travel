package raft

import (
	"testing"
	"time"
)

func TestStep8LeaderHelpers(t *testing.T) {
	node := NewNode("node1")
	node.BecomeLeader()

	accepted, index, err := node.ProposeCommand("BOOK")
	if err != nil {
		t.Fatalf("ProposeCommand returned error: %v", err)
	}
	if !accepted {
		t.Fatal("ProposeCommand should accept a leader write")
	}
	if index != 0 {
		t.Fatalf("ProposeCommand index = %d, want 0", index)
	}
	if !node.IsLeader() {
		t.Fatal("node should report as leader")
	}
	if node.CurrentLeader() != "node1" {
		t.Fatalf("CurrentLeader() = %q, want %q", node.CurrentLeader(), "node1")
	}
	if !node.WaitForCommit(0, time.Second) {
		t.Fatal("WaitForCommit should succeed for a leader-local log entry")
	}
}
