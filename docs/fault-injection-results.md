# TravelRaft Fault-Injection Results

Date: 2026-10-02
Environment: Minikube, Kubernetes v1.37.0

## Baseline

Result: PASS

All four pods were running: MySQL and node1, node2, and node3. Node1 was the latest observed leader before the first fault. A baseline booking through node1 committed successfully:

- PNR: `TRF621500101`
- Seat: `1`
- Status: `CONFIRMED`

## Test 1 - Leader Failure

Action: Deleted the node1 pod.

Expected: Kubernetes recreates node1 and the remaining majority elects a leader.

Result: PASS

Node1 was recreated as `1/1 Running`. Node2 and node3 remained available. Node2 became leader at term 92.

## Test 2 - Booking After Leader Failure

Action: Submitted a booking through node2 after the new election.

Expected: The new leader commits the booking through the remaining majority.

Result: PASS

- PNR: `TRF62274501701`
- Seat: `2`
- Status: `CONFIRMED`

The booking was confirmed and immediately retrievable through gRPC.

## Test 3 - Failed Node Recovery

Action: Observed the node1 Deployment after pod deletion.

Expected: The failed node is recreated and rejoins the cluster.

Result: PASS

Node1 returned to `1/1 Running`. Its replacement used the Kubernetes peer address `node1:50051` and connected to MySQL successfully.

## Test 4 - One-Follower Failure

Action: Deleted the node3 pod while node1 and node2 remained available.

Expected: The remaining two-node majority can still commit a booking.

Result: PASS

Node3 was recreated. During the brief election transition node1 became leader at term 95. A booking through node1 committed successfully:

- PNR: `TRF20114697601`
- Seat: `3`
- Status: `CONFIRMED`

## Test 5 - Loss of Majority

Action: Scaled node1 and node2 to zero replicas, leaving only node3.

Expected: One node cannot commit a new Raft command because a 2-of-3 majority is unavailable.

Result: PASS

A booking attempt through node3 failed with `BOOKING_COMMIT_TIMEOUT`. The three-node cluster had only one active node, so no majority commit was possible.

## Test 6 - Cluster Recovery

Action: Scaled node1 and node2 back to one replica each.

Expected: All three Raft nodes return to service.

Result: PASS

Node1, node2, and node3 all returned to `1/1 Running`.

## Test 7 - MySQL Restart

Action: Deleted the MySQL pod managed by `travelraft-mysql`.

Expected: Kubernetes recreates MySQL and the persistent booking data remains available.

Result: PASS

The MySQL Deployment recreated a ready pod. The `mysql-pvc` remained bound, and the database still contained all three confirmed bookings:

- `TRF621500101`
- `TRF62274501701`
- `TRF20114697601`

## Final Verification

The final cluster had four ready pods:

- `travelraft-mysql`: `1/1 Running`
- `travelraft-node1`: `1/1 Running`
- `travelraft-node2`: `1/1 Running`
- `travelraft-node3`: `1/1 Running`

The Kubernetes Services `mysql`, `node1`, `node2`, and `node3` remained available. Detailed `CommitIndex` and `LastApplied` values are not externally exposed by the current Raft gRPC API; Raft package tests cover those internal state transitions.
