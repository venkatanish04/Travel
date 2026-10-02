# Kubernetes

Start Minikube with Docker, make `travelraft:latest` available in its runtime, then apply the resources:

```powershell
minikube start --driver=docker
minikube image build -t travelraft:latest .
kubectl apply -f k8s/
```

Check deployment status with:

```powershell
kubectl get pods -n travelraft
kubectl get services -n travelraft
kubectl logs deployment/travelraft-node1 -n travelraft
```

The manifests provide three independently addressed Raft pods, one shared MySQL deployment, and a persistent MySQL volume. The application uses Kubernetes service DNS names (`node1`, `node2`, and `node3`) when `TRAVELRAFT_ENV=kubernetes`; MySQL is reached at `mysql:3306`.

Verified locally on 2026-10-02 with Minikube Kubernetes v1.37.0:

- `travelraft-mysql`, `travelraft-node1`, `travelraft-node2`, and `travelraft-node3` reached `1/1 Running`.
- Services `mysql`, `node1`, `node2`, and `node3` were created on ports `3306`, `50051`, `50052`, and `50053`.
- `mysql-pvc` reached `Bound`, and the `travelraft` schema contained `bookings`, `seats`, `users`, and `vehicles`.
- Node1 became leader for the latest observed term after the cluster's startup elections.
