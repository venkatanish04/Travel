package raft

import (
	"context"
	"math/rand"
	"time"
)

const (
	heartbeatInterval      = 100 * time.Millisecond
	minimumElectionTimeout = 300 * time.Millisecond
	maximumElectionTimeout = 500 * time.Millisecond
)

func (n *Node) Start(ctx context.Context) {
	n.startOnce.Do(func() {
		go n.electionLoop(ctx)
		go n.heartbeatLoop(ctx)
	})
}

func (n *Node) Stop() {
	n.stopOnce.Do(func() { close(n.stopCh) })
}

func (n *Node) electionLoop(ctx context.Context) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	timer := time.NewTimer(randomElectionTimeout(random))
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.stopCh:
			return
		case <-n.electionReset:
			resetTimer(timer, randomElectionTimeout(random))
		case <-timer.C:
			if n.GetState() != Leader {
				n.StartElection()
			}
			resetTimer(timer, randomElectionTimeout(random))
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

func randomElectionTimeout(random *rand.Rand) time.Duration {
	return minimumElectionTimeout + time.Duration(random.Int63n(int64(maximumElectionTimeout-minimumElectionTimeout)))
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
