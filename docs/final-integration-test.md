# TravelRaft Final Integration Testing

Date: 2026-10-02

## Automated validation

### Test 1 - Repository test suite

Result: PASS

Command:

```powershell
go test -count=1 travelraft/...
```

All packages passed, including the booking, Raft, state, storage, and integration tests.

### Test 2 - Static validation

Result: PASS

Command:

```powershell
go vet travelraft/...
```

No diagnostics were reported.

### Test 3 - Repeatable database lifecycle test

Result: PASS

The lifecycle test now selects the first available seeded seat instead of assuming that A1 is unused. This allows repeated runs against the persistent MySQL database.

### Test 4 - Benchmark smoke test

Result: PASS

Command:

```powershell
go test -run=^$ -bench=BenchmarkReservationCommand travelraft/tests
```

`BenchmarkReservationCommand` completed successfully at approximately 3,872 ns/op.

### Test 5 - Race test suite

Result: BLOCKED BY ENVIRONMENT

The Windows ThreadSanitizer runtime failed to allocate its shadow memory for every tested package. No race result was produced.

## Distributed environment validation

### Test 6 - MySQL prerequisite

Result: PASS

All three server processes connected to the local `travelraft` database and initialized successfully.

### Test 7 - Docker prerequisite

Result: PASS

Docker Desktop daemon responded with server version `29.8.1`. No containers were required for this local-process test.

## Live three-node scenarios

### Test 8 - Initial leader election

Result: PASS

Nodes started on ports `50051`, `50052`, and `50053`. Node1 became leader for term 13; node2 and node3 remained running as followers.

### Test 9 - Search and seat availability

Result: PASS

The gRPC client found `TR101 - Vande Bharat Express` for the Vijayawada to Hyderabad search and listed the seat statuses.

### Test 10 - Booking and lookup

Result: PASS

The client booked an available seat through node1 and successfully looked it up. Observed PNR: `TRF59642450001`.

### Test 11 - Leader failure and re-election

Result: PASS

Node1 was stopped. Node3 became leader for term 14 while node2 remained available, preserving the two-node majority.

### Test 12 - Booking after leader failure

Result: PASS

The client connected to node3 and completed a booking and lookup after node1 stopped. Observed PNR: `TRF14122340001`.

### Test 13 - Failed-node restart

Result: PASS

Node1 restarted and connected to MySQL as a running cluster member. Its console did not emit a follower status line, so detailed log-index convergence was not separately verified.

### Test 14 - One-follower failure

Result: PASS

Node2 was stopped while node1 and node3 remained available. A booking through node3 succeeded. Observed PNR: `TRF58403220002`.

### Test 15 - Two-node failure

Result: PASS

Only node2 was started. A new client booking was rejected with `NOT_LEADER`; no majority was available to commit the command.

### Test 16 - Duplicate booking rejection

Result: PASS

The live gRPC check attempted to reserve an already-booked seat. The request was rejected immediately with `seat is already booked or does not belong to vehicle`, without an apply timeout.

### Test 17 - Cancellation and seat release

Result: PASS

The live gRPC check cancelled a confirmed booking successfully. A subsequent lookup returned `CANCELLED`, and direct MySQL inspection showed the seat available again.

### Test 18 - Direct MySQL verification

Result: PASS

The MySQL container was queried directly. The `travelraft` database contained the expected tables, and the booking query returned both confirmed and cancelled rows.

### Test 19 - Raft status and index convergence

Result: PARTIALLY VERIFIED

Leader election, replicated booking, failover, and recovery succeeded across the three live nodes. Detailed `CommitIndex`, `LastApplied`, and `LogLength` values cannot currently be queried from a running node because the Raft gRPC API exposes only vote and append RPCs. The existing Raft package tests cover these internal indexes.

## Current conclusion

The automated suite and live Docker cluster passed leader election, gRPC search, booking, lookup, duplicate rejection, cancellation, direct database verification, leader failover, recovery startup, majority operation, and no-majority rejection. Detailed Raft index convergence remains covered by package tests rather than a live status endpoint.
