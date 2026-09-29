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

The production server uses MySQL. The legacy SQLite adapter remains available for offline package tests until MySQL is installed in the development environment.

## Docker

```powershell
docker build -t travelraft:latest .
docker run --rm -p 8080:8080 -e NODE_ID=node1 travelraft:latest
```

## Kubernetes

```powershell
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/storage.yaml
kubectl apply -f k8s/node1.yaml -f k8s/node1-service.yaml
kubectl apply -f k8s/node2.yaml -f k8s/node2-service.yaml
kubectl apply -f k8s/node3.yaml -f k8s/node3-service.yaml
```

The Kubernetes resources run three independently addressed nodes. Peer RPC and durable replicated-log coordination are still required before these pods provide production Raft fault tolerance; the current implementation is a single-node-capable prototype.
