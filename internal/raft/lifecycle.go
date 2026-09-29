package raft

import (
	"context"
	"math/rand"
	"time"
)

const (
	heartbeatInterval  = 500 * time.Millisecond
	minElectionTimeout = 1500 * time.Millisecond
	maxElectionTimeout = 3000 * time.Millisecond
)

func (n *Node) Start(ctx context.Context) {
	n.startOnce.Do(func() {
		go n.electionLoop(ctx)
		go n.heartbeatLoop(ctx)
	})
}

func (n *Node) StartElectionLoop() { n.Start(context.Background()) }

func (n *Node) StartLeaderLoop() { n.Start(context.Background()) }

func (n *Node) Stop() {
	n.stopOnce.Do(func() { close(n.stopCh) })
}

func (n *Node) electionLoop(ctx context.Context) {
	timer := time.NewTimer(randomElectionTimeout())
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.stopCh:
			return
		case <-n.electionReset:
			resetTimer(timer, randomElectionTimeout())
		case <-timer.C:
			if n.GetState() != Leader {
				n.StartElection()
			}
			resetTimer(timer, randomElectionTimeout())
		}
	}
}

func (n *Node) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.stopCh:
			return
		case <-ticker.C:
			if n.GetState() == Leader {
				n.broadcastAppendEntries(ctx)
			}
		}
	}
}

func randomElectionTimeout() time.Duration {
	return minElectionTimeout + time.Duration(rand.Int63n(int64(maxElectionTimeout-minElectionTimeout)))
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(duration)
}

func (n *Node) signalElectionReset() {
	select {
	case n.electionReset <- struct{}{}:
	default:
	}
}
