package raft

func (n *Node) IsLeader() bool   { return n.State().Role == Leader }
func (n *Node) LeaderID() string { return n.State().LeaderID }
