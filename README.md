# kube-pod-rightsizer

watches your kubernetes pods and tells you which ones are wasting resources.

## why

i got tired of seeing pods with 1Gi memory requests using 50Mi. this thing compares what pods actually use vs what they request and spits out recommendations.

it doesn't change anything - just watches and logs.

## quick start

needs metrics-server in your cluster.

```bash
# if you're on kind
kind create cluster
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'

# deploy
kubectl apply -f deploy/
```

or helm:

```bash
helm install kube-pod-rightsizer ./charts/kube-pod-rightsizer -n kube-pod-rightsizer --create-namespace
```

## building

```bash
go build -o kube-pod-rightsizer ./cmd
```

docker:

```bash
docker build -t kube-pod-rightsizer:latest .
```

## config

env vars:

- `SCAN_INTERVAL` - default 30s
- `NAMESPACES` - comma separated, empty = all namespaces
- `LOG_LEVEL` - debug/info/warn/error

## how it works

every 30s (or whatever you set):
1. grabs metrics from metrics-server
2. compares against pod requests
3. if usage is way below requests, logs a recommendation

adds 20% buffer to recommendations so you don't cut it too close.

## output looks like

```
INF recommendation pod=nginx-xyz namespace=default current_cpu=500m recommended_cpu=60m cpu_savings=88%
```

## structure

```
cmd/main.go              - entrypoint
internal/controller/     - main loop
internal/metrics/        - talks to metrics-server
internal/recommender/    - calculates recommendations
deploy/                  - k8s manifests
charts/                  - helm
```

## todo

- [ ] prometheus metrics endpoint
- [ ] webhook for slack/teams notifications
- [ ] historical data tracking

