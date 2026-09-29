package raft

func (n *Node) StartElection() State {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state.CurrentTerm++
	n.state.Role = Leader
	n.state.LeaderID = n.id
	n.state.VotedFor = n.id
	return n.state
}
