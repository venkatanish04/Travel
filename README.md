# TravelRaft

TravelRaft is a Go booking service prototype with a Raft-shaped replicated state core. It supports atomic reservations for trains, buses, and flights through a JSON HTTP API.

## Run

```powershell
go test ./...
go run ./cmd/server
go run ./cmd/client -customer alice -trip train-1 -seat A1
```

The server exposes `GET /health`, `GET /reservations`, `POST /reservations`, and `DELETE /reservations/{id}`. The protobuf files under `api/proto` define the intended wire contracts; the runnable prototype uses JSON HTTP to avoid requiring code generation during setup.

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
