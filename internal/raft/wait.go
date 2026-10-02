package raft

import (
	"time"
)

func (n *Node) WaitForCommit(index int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		n.mu.Lock()
		committed := n.CommitIndex >= index
		n.mu.Unlock()
		if committed {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

func (n *Node) WaitForApply(index int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		n.mu.RLock()
		applied := n.LastApplied >= index
		n.mu.RUnlock()
		if applied {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}
