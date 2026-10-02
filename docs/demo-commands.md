# TravelRaft Demo Commands

```cmd
cd "C:\Users\venka\OneDrive\Desktop\os project\TravelRaft"
minikube start --driver=docker
kubectl config use-context minikube
minikube image build -t travelraft:latest .
kubectl apply -f k8s
kubectl rollout status deployment/travelraft-mysql -n travelraft
kubectl rollout status deployment/travelraft-node1 -n travelraft
kubectl rollout status deployment/travelraft-node2 -n travelraft
kubectl rollout status deployment/travelraft-node3 -n travelraft
kubectl get pods -n travelraft
kubectl get services -n travelraft
```

```text
travelraft-mysql    1/1   Running
travelraft-node1    1/1   Running
travelraft-node2    1/1   Running
travelraft-node3    1/1   Running
```

```cmd
kubectl logs deployment/travelraft-node1 -n travelraft | findstr LEADER
kubectl logs deployment/travelraft-node2 -n travelraft | findstr LEADER
kubectl logs deployment/travelraft-node3 -n travelraft | findstr LEADER
```

```cmd
REM Run this in a separate cmd window.
REM Replace node3 and 15057:50053 with the current leader Service and port.
kubectl port-forward service/node3 15057:50053 -n travelraft
```

```cmd
REM Run this in another cmd window while port-forward is running.
go run ./cmd/client -address localhost:15057
```

```text
BOOKING CONFIRMED
PNR: TRF...
Lookup found: true
```

```cmd
kubectl exec deployment/travelraft-mysql -n travelraft -- mysql -uroot -proot -e "SELECT id,pnr,seat_id,status FROM travelraft.bookings ORDER BY created_at DESC;"
go test -count=1 ./...
go vet ./...
go build ./...
```

```cmd
REM Replace node1 in all commands below with the current leader node.
kubectl delete pod -l app=travelraft,node=node1 -n travelraft
kubectl rollout status deployment/travelraft-node1 -n travelraft
kubectl get pods -n travelraft
kubectl logs deployment/travelraft-node2 -n travelraft | findstr LEADER
kubectl logs deployment/travelraft-node3 -n travelraft | findstr LEADER
```

```cmd
kubectl scale deployment travelraft-node1 --replicas=0 -n travelraft
kubectl scale deployment travelraft-node2 --replicas=0 -n travelraft
kubectl get pods -n travelraft
kubectl scale deployment travelraft-node1 --replicas=1 -n travelraft
kubectl scale deployment travelraft-node2 --replicas=1 -n travelraft
kubectl rollout status deployment/travelraft-node1 -n travelraft
kubectl rollout status deployment/travelraft-node2 -n travelraft
kubectl get pods -n travelraft
```
