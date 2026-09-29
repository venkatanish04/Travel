package raft

import (
	"sync"
)

type StateMachine interface{ Apply([]byte) error }

type Node struct {
	mu          sync.RWMutex
	id          string
	state       State
	log         Log
	commitIndex uint64
	machine     StateMachine
}

func NewNode(id string, machine StateMachine) *Node {
	return &Node{id: id, state: State{Role: Leader, LeaderID: id}, machine: machine}
}

func (n *Node) ID() string { return n.id }

func (n *Node) State() State { n.mu.RLock(); defer n.mu.RUnlock(); return n.state }

func (n *Node) Entries() []Entry { n.mu.RLock(); defer n.mu.RUnlock(); return n.log.Entries() }

func (n *Node) Apply(command []byte) (Entry, error) {
	if len(command) == 0 {
		return Entry{}, ErrEmptyCommand
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state.Role != Leader {
		return Entry{}, ErrNotLeader
	}
	entry := n.log.Append(n.state.CurrentTerm, command)
	if err := n.machine.Apply(command); err != nil {
		return Entry{}, err
	}
	n.commitIndex = entry.Index
	return entry, nil
}
