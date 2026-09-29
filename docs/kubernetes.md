# Kubernetes

Build or load the image into the target cluster, then apply the resources:

```powershell
docker build -t travelraft:latest .
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/storage.yaml
kubectl apply -f k8s/node1.yaml -f k8s/node1-service.yaml
kubectl apply -f k8s/node2.yaml -f k8s/node2-service.yaml
kubectl apply -f k8s/node3.yaml -f k8s/node3-service.yaml
```

The manifests provide three independently addressed pods and a persistent volume claim. The current Go implementation is a single-node-capable Raft prototype; peer RPC and shared durable log replication remain required before treating these pods as a fault-tolerant cluster.
