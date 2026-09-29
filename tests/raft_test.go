package tests

import (
	"testing"
	"travelraft/internal/raft"
	"travelraft/internal/state"
)

func TestNodeElectionAndCommit(t *testing.T) {
	machine := state.NewMachine()
	node := raft.NewNode("node1", machine)
	if !node.IsLeader() {
		t.Fatal("new standalone node should be leader")
	}
	term := node.StartElection().CurrentTerm
	if term != 1 {
		t.Fatalf("term = %d, want 1", term)
	}
	entry, err := node.Apply([]byte(`{"type":"reserve","reservation":{"id":"r1","customer_id":"alice","trip_id":"train-1","seat_id":"A1","status":"confirmed"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if entry.Index != 1 || len(node.Entries()) != 1 {
		t.Fatalf("entry = %#v, log length = %d", entry, len(node.Entries()))
	}
}
