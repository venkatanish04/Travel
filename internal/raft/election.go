package raft

import (
	"context"
	"sync"
	"time"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (n *Node) StartElection() State {
	n.mu.Lock()
	n.CurrentTerm++
	n.state.CurrentTerm = uint64(n.CurrentTerm)
	n.state.Role = Candidate
	n.state.LeaderID = ""
	n.state.VotedFor = n.id
	n.VotedFor = n.id
	term := n.CurrentTerm
	lastIndex := n.log.LastIndex()
	lastTerm := n.log.LastTerm()
	peers := make(map[string]string, len(n.Peers))
	for id, address := range n.Peers {
		peers[id] = address
	}
	n.mu.Unlock()

	votes := 1
	var votesMu sync.Mutex
	var wait sync.WaitGroup
	for peerID, address := range peers {
		wait.Add(1)
		go func(peerID, address string) {
			defer wait.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return
			}
			defer conn.Close()
			response, err := pb.NewRaftServiceClient(conn).RequestVote(ctx, &pb.RequestVoteRequest{
				Term: int32(term), CandidateId: n.id, LastLogIndex: int32(lastIndex), LastLogTerm: int32(lastTerm),
			})
			if err != nil {
				return
			}
			n.mu.Lock()
			if response.GetTerm() > int32(n.CurrentTerm) {
				n.becomeFollowerLocked(int(response.GetTerm()))
			}
			n.mu.Unlock()
			if response.GetVoteGranted() {
				votesMu.Lock()
				votes++
				votesMu.Unlock()
			}
			_ = peerID
		}(peerID, address)
	}
	wait.Wait()

	n.mu.Lock()
	defer n.mu.Unlock()
	majority := (len(peers)+1)/2 + 1
	if n.state.Role == Candidate && n.CurrentTerm == term && votes >= majority {
		n.state.Role = Leader
		n.state.LeaderID = n.id
	}
	return n.state
}

func (n *Node) RequestVote(_ context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	if err := validateVoteRequest(req); err != nil {
		return nil, err
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if req.GetTerm() < int32(n.CurrentTerm) {
		return &pb.RequestVoteResponse{Term: int32(n.CurrentTerm)}, nil
	}
	if req.GetTerm() > int32(n.CurrentTerm) {
		n.becomeFollowerLocked(int(req.GetTerm()))
	}
	canVote := n.VotedFor == "" || n.VotedFor == req.GetCandidateId()
	upToDate := req.GetLastLogTerm() > int32(n.log.LastTerm()) ||
		(req.GetLastLogTerm() == int32(n.log.LastTerm()) && req.GetLastLogIndex() >= int32(n.log.LastIndex()))
	if !canVote || !upToDate {
		return &pb.RequestVoteResponse{Term: int32(n.CurrentTerm)}, nil
	}
	n.VotedFor = req.GetCandidateId()
	n.state.VotedFor = n.VotedFor
	n.state.Role = Follower
	return &pb.RequestVoteResponse{Term: int32(n.CurrentTerm), VoteGranted: true}, nil
}

func validateVoteRequest(req *pb.RequestVoteRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.GetCandidateId() == "" {
		return status.Error(codes.InvalidArgument, "candidate ID is required")
	}
	return nil
}
