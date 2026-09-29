package raft

// Replicate applies a committed command to this node's state machine.
// A production cluster would send the entry to peer transports before committing.
func (n *Node) Replicate(command []byte) (Entry, error) { return n.Apply(command) }
