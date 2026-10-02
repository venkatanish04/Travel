# Raft Notes

The current implementation covers node state, terms, elections, log entries, commit index, replication, and peer transport. RequestVote and AppendEntries are exposed through Raft gRPC and used for peer elections, heartbeats, log replication, and majority commits. A new node starts as the single leader only when constructed for local in-memory tests; clustered server processes begin as followers and elect a leader. Durable Raft term/log persistence and production-grade secret management remain future hardening work.
