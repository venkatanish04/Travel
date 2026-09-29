# Raft Notes

The current prototype models the Raft boundaries: node state, terms, elections, log entries, commit index, replication, and peer transport. A new node starts as the single leader so local development works without a cluster. The next production step is implementing RPC-backed RequestVote and AppendEntries between peers, then persisting term and log entries before acknowledging commits.
