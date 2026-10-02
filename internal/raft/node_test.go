package raft

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "travelraft/api/proto"
)

func TestConfiguredNodeStartsAsFollower(t *testing.T) {
	node := NewNode("node1", "localhost:50051", map[string]string{"node2": "localhost:50052"})
	if node.GetState() != Follower {
		t.Fatalf("state = %s, want %s", node.GetState(), Follower)
	}
	if node.Address != "localhost:50051" || len(node.Peers) != 1 {
		t.Fatalf("node configuration was not retained: %#v", node)
	}
}

func TestRequestVoteAndAppendEntries(t *testing.T) {
	node := NewNode("node2", "localhost:50052", map[string]string{})
	vote, err := node.RequestVote(context.Background(), &pb.RequestVoteRequest{Term: 1, CandidateId: "node1", LastLogIndex: -1, LastLogTerm: 0})
	if err != nil || !vote.GetVoteGranted() {
		t.Fatalf("vote = %#v, err = %v", vote, err)
	}
	appendResponse, err := node.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
		Term: 1, LeaderId: "node1", PrevLogIndex: -1,
		Entries: []*pb.LogEntry{{Term: 1, Command: "ADD_TRAIN"}}, LeaderCommit: 0,
	})
	if err != nil || !appendResponse.GetSuccess() {
		t.Fatalf("append response = %#v, err = %v", appendResponse, err)
	}
	if len(node.Entries()) != 1 || node.State().LeaderID != "node1" {
		t.Fatalf("node did not append leader entry: entries=%d state=%#v", len(node.Entries()), node.State())
	}
	if appendResponse.GetMatchIndex() != 0 {
		t.Fatalf("match index = %d, want 0", appendResponse.GetMatchIndex())
	}
}

func TestAddCommandOnlyAppendsOnLeader(t *testing.T) {
	node := NewNode("node1", "localhost:50051", map[string]string{})
	if node.AddCommand("BOOK|TR101|A1|user123") {
		t.Fatal("follower accepted a command")
	}
	node.BecomeLeader()
	if !node.AddCommand("BOOK|TR101|A1|user123") {
		t.Fatal("leader rejected a command")
	}
	if len(node.Entries()) != 1 || node.NextIndex[node.ID()] != 1 {
		t.Fatalf("log/index state = %d/%d, want 1/1", len(node.Entries()), node.NextIndex[node.ID()])
	}
}

type failingStateMachine struct{}

func (f *failingStateMachine) Apply(_ []byte) error {
	return errors.New("rejected")
}

func TestApplyRevertsInvalidCommands(t *testing.T) {
	node := NewNode("node1")
	node.BecomeLeader()
	node.SetStateMachine(&failingStateMachine{})

	if _, err := node.Apply([]byte("INVALID_COMMAND")); err == nil {
		t.Fatal("expected invalid command to fail")
	}
	if len(node.Entries()) != 0 {
		t.Fatalf("invalid command left a log entry behind: %d entries", len(node.Entries()))
	}
	if node.CommitIndex != -1 || node.LastApplied != -1 {
		t.Fatalf("invalid command changed commit state: commit=%d applied=%d", node.CommitIndex, node.LastApplied)
	}
}

type recordingStateMachine struct {
	commands [][]byte
}

func (m *recordingStateMachine) Apply(command []byte) error {
	m.commands = append(m.commands, append([]byte(nil), command...))
	return nil
}

func TestProposeCommandAppliesAfterSingleNodeCommit(t *testing.T) {
	machine := &recordingStateMachine{}
	node := NewNode("node1", machine)

	accepted, index, err := node.ProposeCommand("BOOK")
	if err != nil || !accepted {
		t.Fatalf("proposal accepted=%v err=%v", accepted, err)
	}
	if !node.WaitForApply(index, time.Second) {
		t.Fatal("proposal was not applied")
	}
	if len(machine.commands) != 1 || string(machine.commands[0]) != "BOOK" {
		t.Fatalf("applied commands = %q, want BOOK", machine.commands)
	}
}
