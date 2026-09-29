package raft

import (
	"context"

	pb "travelraft/api/proto"
)

func (n *Node) Replicate(command []byte) (Entry, error) { return n.Apply(command) }

func (n *Node) AppendEntries(_ context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if req.GetTerm() < int32(n.CurrentTerm) {
		return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm)}, nil
	}
	if req.GetTerm() > int32(n.CurrentTerm) {
		n.becomeFollowerLocked(int(req.GetTerm()))
	}
	n.state.Role = Follower
	n.state.LeaderID = req.GetLeaderId()
	if req.GetPrevLogIndex() >= 0 {
		if req.GetPrevLogIndex() >= int32(len(n.log.Entries)) || n.log.Entries[req.GetPrevLogIndex()].Term != uint64(req.GetPrevLogTerm()) {
			return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm)}, nil
		}
	}
	for _, entry := range req.GetEntries() {
		n.log.Append(uint64(entry.GetTerm()), []byte(entry.GetCommand()))
	}
	if req.GetLeaderCommit() > int32(n.CommitIndex) {
		lastIndex := len(n.log.Entries) - 1
		if req.GetLeaderCommit() < int32(lastIndex) {
			n.CommitIndex = int(req.GetLeaderCommit())
		} else {
			n.CommitIndex = lastIndex
		}
	}
	return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm), Success: true}, nil
}
