# Architecture

`cmd/server` hosts the BookingService and RaftService gRPC endpoints. `internal/booking` handles search, seats, users, bookings, and seed data. `internal/raft` elects a leader, replicates log entries with AppendEntries, and advances commits after a majority acknowledgement. `internal/storage` uses MySQL for production node persistence and retains a small JSON/SQLite compatibility adapter for offline tests.
