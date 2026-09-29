package raft

func (n *Node) updateCommitIndex() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.updateCommitIndexLocked()
}

func (n *Node) updateCommitIndexLocked() {
	if n.state.Role != Leader {
		return
	}
	lastIndex := n.log.LastIndex()
	for index := lastIndex; index > n.CommitIndex; index-- {
		replicated := 1
		for peerID, match := range n.MatchIndex {
			if peerID != n.id && match >= index {
				replicated++
			}
		}
		if replicated >= (len(n.Peers)+1)/2+1 && int(n.log.Entries[index].Term) == n.CurrentTerm {
			n.CommitIndex = index
			n.applyCommittedLocked()
			return
		}
	}
}
