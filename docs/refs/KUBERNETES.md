# Déploiement Kubernetes local (Kind)

Prérequis : Docker, `kind`, `kubectl`, Go 1.26 et Node 24.

```bash
kind create cluster --name pulse --config deployments/kind/cluster.yaml
docker build -t pulse-backend:dev -f deployments/docker/backend.Dockerfile .
docker build -t pulse-frontend:dev --build-arg VITE_OTEL_COLLECTOR_URL=/otlp-http/v1/traces -f deployments/docker/frontend.Dockerfile frontend
kind load docker-image pulse-backend:dev pulse-frontend:dev --name pulse
kubectl apply -k deployments/kubernetes/kind
kubectl -n pulse rollout status deployment/pulse-backend
kubectl -n pulse rollout status deployment/pulse-frontend
```

L'application est disponible sur `http://localhost:8080`; sa santé frontend est sur `/healthz`. Le backend expose `/livez`, `/readyz` et `/healthz`. Les traces reçues par le Collector de développement sont visibles avec `kubectl -n pulse logs deployment/otel-collector`.
