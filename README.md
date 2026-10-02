# TravelRaft

TravelRaft is a Go booking service prototype with gRPC booking and Raft services. It supports trains, buses, flights, seats, bookings, and majority-based Raft log replication.

## MySQL setup

Create a MySQL database named `travelraft`, then set these environment variables before starting a node:

```powershell
$env:MYSQL_USER="root"
$env:MYSQL_PASSWORD="your-password"
$env:MYSQL_HOST="127.0.0.1"
$env:MYSQL_PORT="3306"
$env:MYSQL_DATABASE="travelraft"
```

The server creates the required tables on startup. MySQL must be installed and running locally.

## Run and test

```powershell
go test ./...
go vet ./...
go build ./...
go run ./cmd/server --id=node1 --port=50051
```

Start node2 and node3 on ports `50052` and `50053` in separate terminals. The gRPC client connects to node1 with `go run ./cmd/client`.

The production server uses MySQL. Set the `MYSQL_*` environment variables before starting the server.

## Docker

```cmd
docker build -t travelraft:latest .
docker compose up -d
docker compose ps
```

## Kubernetes

```powershell
minikube start --driver=docker
minikube image build -t travelraft:latest .
kubectl apply -f k8s/
kubectl get pods -n travelraft
```

The Kubernetes resources run three independently addressed nodes and one shared MySQL instance backed by a persistent volume. Raft peer RPC and replicated-log coordination are included in the application; durable Raft term/log persistence and production-grade secret management remain future hardening work.
