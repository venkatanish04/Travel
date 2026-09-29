package raft

import (
	"context"

	pb "travelraft/api/proto"
)

type Peer struct {
	ID      string
	Address string
}

type Transport struct{ Peers []Peer }

func NewTransport(peers []Peer) *Transport { return &Transport{Peers: peers} }

type GRPCServer struct {
	pb.UnimplementedRaftServiceServer
	Node *Node
}

func NewGRPCServer(node *Node) *GRPCServer { return &GRPCServer{Node: node} }

func (s *GRPCServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	if err := validateVoteRequest(req); err != nil {
		return nil, err
	}
	return s.Node.RequestVote(ctx, req)
}

func (s *GRPCServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	return s.Node.AppendEntries(ctx, req)
}
