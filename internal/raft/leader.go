package raft

func (n *Node) LeaderID() string { return n.State().LeaderID }
