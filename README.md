# kube-pod-rightsizer

Kubernetes controller that analyzes pod resource usage and recommends right-sizing. Compares actual metrics against requests and logs savings opportunities.

Read-only - observes and reports, doesn't modify anything.

## Building

```bash
go build -o kube-pod-rightsizer ./cmd
```

## Running

Requires metrics-server in your cluster.

```bash
./kube-pod-rightsizer
```

## Config

| Variable | Default | Description |
|----------|---------|-------------|
| `SCAN_INTERVAL` | 30s | How often to analyze pods |
| `NAMESPACES` | all | Comma separated list |
| `LOG_LEVEL` | info | debug/info/warn/error |

## How it works

1. Fetches pod metrics from metrics-server
2. Compares usage against resource requests
3. Logs recommendations for over-provisioned pods

Recommendations include 20% buffer as safety margin.

## Output

```
INF recommendation pod=nginx-xyz namespace=default current_cpu=500m recommended_cpu=60m cpu_savings=88%
```

## Todo

- [ ] Dockerfile and k8s manifests
- [ ] Helm chart
- [ ] Prometheus metrics endpoint
