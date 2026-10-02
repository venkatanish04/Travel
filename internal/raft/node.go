package raft

import (
	"context"
	"fmt"
	"sync"
)

type StateMachine interface{ Apply([]byte) error }

type Node struct {
	mu            sync.RWMutex
	id            string
	Address       string
	Peers         map[string]string
	state         State
	log           Log
	CommitIndex   int
	LastApplied   int
	CurrentTerm   int
	VotedFor      string
	machine       StateMachine
	NextIndex     map[string]int
	MatchIndex    map[string]int
	electionReset chan struct{}
	startOnce     sync.Once
	stopOnce      sync.Once
	stopCh        chan struct{}
}

func NewNode(id string, args ...interface{}) *Node {
	node := &Node{
		id:            id,
		Peers:         make(map[string]string),
		state:         State{Role: Follower},
		CommitIndex:   -1,
		LastApplied:   -1,
		NextIndex:     make(map[string]int),
		MatchIndex:    make(map[string]int),
		electionReset: make(chan struct{}, 1),
		stopCh:        make(chan struct{}),
	}

	if len(args) == 1 {
		node.machine, _ = args[0].(StateMachine)
		node.state = State{Role: Leader, LeaderID: id}
	}
	if len(args) >= 2 {
		node.Address, _ = args[0].(string)
		if peers, ok := args[1].(map[string]string); ok {
			node.Peers = peers
		}
	}
	return node
}

func (n *Node) ID() string { return n.id }

func (n *Node) SetStateMachine(machine StateMachine) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.machine = machine
}

func (n *Node) State() State { n.mu.RLock(); defer n.mu.RUnlock(); return n.state }

func (n *Node) GetState() Role { return n.State().Role }

func (n *Node) Entries() []Entry {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]Entry, len(n.log.Entries))
	copy(result, n.log.Entries)
	return result
}

func (n *Node) BecomeFollower(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.becomeFollowerLocked(term)
}

func (n *Node) becomeFollowerLocked(term int) {
	n.state.Role = Follower
	n.state.LeaderID = ""
	n.state.CurrentTerm = uint64(term)
	n.state.VotedFor = ""
	n.CurrentTerm = term
	n.VotedFor = ""
}

func (n *Node) BecomeCandidate() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.CurrentTerm++
	n.state.Role = Candidate
	n.state.LeaderID = ""
	n.state.CurrentTerm = uint64(n.CurrentTerm)
	n.state.VotedFor = n.id
	n.VotedFor = n.id
}

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state.Role = Leader
	n.state.LeaderID = n.id
	n.initializeReplicationLocked()
}

func (n *Node) initializeReplicationLocked() {
	for peerID := range n.Peers {
		n.NextIndex[peerID] = len(n.log.Entries)
		n.MatchIndex[peerID] = -1
	}
}

func (n *Node) PrintStatus() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	fmt.Println("----------------------------------------")
	fmt.Println("Node ID:", n.id)
	fmt.Println("Address:", n.Address)
	fmt.Println("State:", n.state.Role)
	fmt.Println("Current Term:", n.CurrentTerm)
	fmt.Println("Voted For:", n.VotedFor)
	fmt.Println("Commit Index:", n.CommitIndex)
	fmt.Println("Last Applied:", n.LastApplied)
	fmt.Println("Log Entries:", len(n.log.Entries))
	fmt.Println("----------------------------------------")
}

func (n *Node) Apply(command []byte) (Entry, error) {
	if len(command) == 0 {
		return Entry{}, ErrEmptyCommand
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state.Role != Leader {
		return Entry{}, ErrNotLeader
	}
	entry := n.log.Append(uint64(n.CurrentTerm), command)
	if n.machine != nil {
		if err := n.machine.Apply(command); err != nil {
			if lastIndex := n.log.LastIndex(); lastIndex >= 0 {
				n.log.DeleteFrom(lastIndex)
			}
			return Entry{}, err
		}
	}
	n.CommitIndex = int(entry.Index)
	n.LastApplied = n.CommitIndex
	return entry, nil
}

func (n *Node) AddCommand(command string) bool {
	if command == "" {
		return false
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state.Role != Leader {
		return false
	}
	n.log.Append(uint64(n.CurrentTerm), []byte(command))
	n.MatchIndex[n.id] = n.log.LastIndex()
	n.NextIndex[n.id] = n.log.LastIndex() + 1
	return true
}

func (n *Node) ProposeCommand(command string) (bool, int, error) {
	if command == "" {
		return false, -1, ErrEmptyCommand
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state.Role != Leader {
		return false, -1, nil
	}

	entry := LogEntry{Term: uint64(n.CurrentTerm), Command: []byte(command)}
	n.log.Append(entry)

	index := n.log.LastIndex()
	n.MatchIndex[n.id] = index
	n.NextIndex[n.id] = index + 1

	if len(n.Peers) == 0 {
		n.CommitIndex = index
		n.applyCommittedLocked()
	}
	return true, index, nil
}

func (n *Node) ReplicateNow() {
	if n.GetState() != Leader {
		return
	}
	n.broadcastAppendEntries(context.Background())
}

func (n *Node) ApplyCommittedEntries() {
	n.mu.Lock()
	n.applyCommittedLocked()
	n.mu.Unlock()
}
