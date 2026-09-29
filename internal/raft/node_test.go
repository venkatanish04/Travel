package raft

import (
	"context"
	"testing"

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
}
