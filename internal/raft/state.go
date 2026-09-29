package raft

import "errors"

type Role string

const (
	Follower  Role = "follower"
	Candidate Role = "candidate"
	Leader    Role = "leader"
)

var ErrNotLeader = errors.New("raft node is not the leader")
var ErrEmptyCommand = errors.New("raft command is empty")

type State struct {
	CurrentTerm uint64
	VotedFor    string
	Role        Role
	LeaderID    string
}
