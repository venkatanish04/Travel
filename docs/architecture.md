# Architecture

`cmd/server` owns the HTTP process. `internal/booking` translates booking operations into serialized commands. `internal/raft` appends and commits commands through a node state machine. `internal/state` applies commands atomically and owns reservation state. `internal/storage` provides a small JSON persistence adapter.
