package raft

import (
	"context"
	"time"
)

func (n *Node) StartApplyLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				n.ApplyCommittedEntries()
			case <-ctx.Done():
				return
			case <-n.stopCh:
				return
			}
		}
	}()
}
