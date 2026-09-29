package raft

import (
	"context"
	"sync"
	"time"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (n *Node) Replicate(command []byte) (Entry, error) {
	if len(command) == 0 {
		return Entry{}, ErrEmptyCommand
	}
	n.mu.Lock()
	if n.state.Role != Leader {
		n.mu.Unlock()
		return Entry{}, ErrNotLeader
	}
	entry := n.log.Append(uint64(n.CurrentTerm), command)
	n.mu.Unlock()

	if len(n.Peers) == 0 {
		n.mu.Lock()
		n.CommitIndex = len(n.log.Entries) - 1
		n.applyCommittedLocked()
		n.mu.Unlock()
		return entry, nil
	}
	n.broadcastAppendEntries(context.Background())
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.advanceCommitLocked() < int(entry.Index-1) {
		return entry, ErrReplicationFailed
	}
	return entry, nil
}

func (n *Node) AppendEntries(_ context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if req.GetTerm() < int32(n.CurrentTerm) {
		return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm), MatchIndex: int32(n.log.LastIndex())}, nil
	}
	if req.GetTerm() > int32(n.CurrentTerm) {
		n.becomeFollowerLocked(int(req.GetTerm()))
	}
	n.state.Role = Follower
	n.state.LeaderID = req.GetLeaderId()
	if req.GetPrevLogIndex() >= 0 {
		if req.GetPrevLogIndex() >= int32(len(n.log.Entries)) || n.log.Entries[req.GetPrevLogIndex()].Term != uint64(req.GetPrevLogTerm()) {
			return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm), MatchIndex: int32(n.log.LastIndex())}, nil
		}
	}
	insertAt := int(req.GetPrevLogIndex()) + 1
	for offset, entry := range req.GetEntries() {
		position := insertAt + offset
		if position < len(n.log.Entries) {
			if n.log.Entries[position].Term != uint64(entry.GetTerm()) || string(n.log.Entries[position].Command) != entry.GetCommand() {
				n.log.Entries = n.log.Entries[:position]
			}
		}
		if position >= len(n.log.Entries) {
			n.log.Append(uint64(entry.GetTerm()), []byte(entry.GetCommand()))
		}
	}
	if req.GetLeaderCommit() > int32(n.CommitIndex) {
		lastIndex := len(n.log.Entries) - 1
		if req.GetLeaderCommit() < int32(lastIndex) {
			n.CommitIndex = int(req.GetLeaderCommit())
		} else {
			n.CommitIndex = lastIndex
		}
	}
	n.applyCommittedLocked()
	n.signalElectionReset()
	return &pb.AppendEntriesResponse{Term: int32(n.CurrentTerm), Success: true}, nil
}

func (n *Node) broadcastAppendEntries(ctx context.Context) {
	n.mu.RLock()
	if n.state.Role != Leader {
		n.mu.RUnlock()
		return
	}
	term := n.CurrentTerm
	commitIndex := n.CommitIndex
	peers := make(map[string]string, len(n.Peers))
	for id, address := range n.Peers {
		peers[id] = address
	}
	n.mu.RUnlock()

	var wait sync.WaitGroup
	for peerID, address := range peers {
		wait.Add(1)
		go func(peerID, address string) {
			defer wait.Done()
			request := n.appendRequest(peerID, term, commitIndex)
			callContext, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return
			}
			defer conn.Close()
			response, err := pb.NewRaftServiceClient(conn).AppendEntries(callContext, request)
			if err != nil {
				return
			}
			n.mu.Lock()
			defer n.mu.Unlock()
			if response.GetTerm() > int32(n.CurrentTerm) {
				n.becomeFollowerLocked(int(response.GetTerm()))
				return
			}
			if response.GetSuccess() {
				n.MatchIndex[peerID] = int(response.GetMatchIndex())
				n.NextIndex[peerID] = int(response.GetMatchIndex()) + 1
			} else if n.NextIndex[peerID] > 0 {
				n.NextIndex[peerID]--
			}
		}(peerID, address)
	}
	wait.Wait()
	n.mu.Lock()
	n.advanceCommitLocked()
	n.mu.Unlock()
}

func (n *Node) advanceCommitLocked() int {
	n.updateCommitIndexLocked()
	return n.CommitIndex
}

func (n *Node) appendRequest(peerID string, term, commitIndex int) *pb.AppendEntriesRequest {
	n.mu.RLock()
	defer n.mu.RUnlock()
	next := n.NextIndex[peerID]
	if next < 0 {
		next = 0
	}
	if next > len(n.log.Entries) {
		next = len(n.log.Entries)
	}
	request := &pb.AppendEntriesRequest{Term: int32(term), LeaderId: n.id, PrevLogIndex: int32(next - 1), LeaderCommit: int32(commitIndex)}
	if next > 0 && next <= len(n.log.Entries) {
		request.PrevLogTerm = int32(n.log.Entries[next-1].Term)
	}
	for _, entry := range n.log.Entries[next:] {
		request.Entries = append(request.Entries, &pb.LogEntry{Term: int32(entry.Term), Command: string(entry.Command)})
	}
	return request
}

func (n *Node) applyCommittedLocked() {
	for n.LastApplied < n.CommitIndex {
		position := n.LastApplied + 1
		if position >= len(n.log.Entries) {
			return
		}
		if n.machine != nil {
			_ = n.machine.Apply(n.log.Entries[position].Command)
		}
		n.LastApplied = position
	}
}
