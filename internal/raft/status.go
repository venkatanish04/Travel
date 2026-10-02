package raft

type Status struct {
	NodeID      string
	State       string
	Term        int
	CommitIndex int
	LastApplied int
	LogLength   int
}

func (n *Node) IsLeader() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state.Role == Leader
}

func (n *Node) CurrentLeader() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.state.Role == Leader {
		return n.id
	}
	return ""
}

func (n *Node) GetStatus() Status {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return Status{
		NodeID:      n.id,
		State:       n.state.Role.String(),
		Term:        n.CurrentTerm,
		CommitIndex: n.CommitIndex,
		LastApplied: n.LastApplied,
		LogLength:   len(n.log.Entries),
	}
}
