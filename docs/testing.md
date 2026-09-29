# Testing

Run `go test ./...` for booking lifecycle, duplicate-seat, concurrent reservation, and log-commit checks. The concurrency test uses distinct seats and verifies that the service's atomic ID generation and Raft-protected apply path preserve all successful requests.
