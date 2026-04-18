# isolet

Isolet is a self-hosted CTF platform designed to be organizer-friendly and resource-efficient. It supports on-demand challenge instances, real-time updates, team management, dynamic scoring, and full Kubernetes-native deployment.

Originally built in college to run Bandit-style wargames. Evolved into a production platform after hosting live CTF events in 2024–2025 and learning what actually breaks under load.

## Architecture

```
                         ┌─────────────────────────────────────────┐
                         │              Kubernetes Cluster          │
                         │                                         │
  Browser ─── Proxy ─── Oracle (REST/consumer) ── PostgreSQL      │
              (Nginx)  │     └── Valkey (cache/sessions/pubsub)   │
                       │     └── Tide (K8s operator) ─── Challenge Pods
                       │                                           │
                       └── Pulse (WebSocket) ── Valkey (PubSub)   │
                                                                   │
  Herald (event watcher) ─── Kafka ──► Oracle (consumer)          │
                                  └──► Pulse (consumer)            │
                                                                   │
  Phoros (file server) — stores downloadable challenge files       │
                         └─────────────────────────────────────────┘
```

**Request path:** Browser → Proxy → UI (Next.js) or Oracle API or Pulse WebSocket

**Instance lifecycle:** Oracle calls Tide SDK → Tide reconciles `Instance` CR → K8s Deployment + Service + Traefik IngressRoute → Herald watches pod events → emits facts to Kafka → Oracle consumer updates instance state → Pulse pushes real-time update to user's browser.

## Services

| Service | Language | Role |
|---------|----------|------|
| [oracle](./oracle/) | Go (Fiber) | Central API + Kafka consumer. Auth, challenges, instances, scoring. |
| [tide](./tide/) | Go (Kubebuilder) | Kubernetes operator managing `Instance` CRDs → pods. |
| [pulse](./pulse/) | Elixir/Phoenix | Real-time WebSocket gateway. Kafka consumer → Redis PubSub → browser. |
| [herald](./herald/) | Go | Event emitter. Watches K8s events and PostgreSQL, publishes facts to Kafka. |
| [ui](./ui/) | Next.js / TypeScript | CTF frontend: challenges, scoreboard, team management, instance controls. |
| [proxy](./proxy/) | Nginx | Reverse proxy. Routes `/api` → Oracle, `/socket` → Pulse, `/files` → Phoros, `/` → UI. |
| [phoros](https://github.com/TheAlpha16/phoros) | — | File server for challenge attachments. External service. |

## Tech Stack

- **Backend:** Go, Elixir/Phoenix
- **Frontend:** Next.js 14, TypeScript, Tailwind CSS, Radix UI
- **Data:** PostgreSQL, Valkey (Redis-compatible)
- **Messaging:** Apache Kafka (Redpanda in dev)
- **Infrastructure:** Kubernetes, Helm, Traefik, Nginx
- **Observability:** Prometheus, Grafana, OpenTelemetry, Sentry

## Local Development

Requires Docker and Docker Compose. This setup runs the full application stack but **does not include Tide or a real Kubernetes cluster** — challenge instance spawning will not work locally without a cluster.

```sh
docker compose up
```

This starts: Oracle, UI, Proxy, Valkey, PostgreSQL, Redpanda (Kafka), Redpanda Console, Prometheus, Grafana.

| Service | URL |
|---------|-----|
| UI + API | http://localhost |
| Redpanda Console | http://localhost:8085 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3005 (admin/admin) |

For full instance functionality, run Tide against a local cluster:

```sh
cd tide
make run   # runs controller against your active kubeconfig
```

## Production Deployment

Isolet ships a Helm chart in [`charts/`](./charts/).

### Prerequisites

- Kubernetes cluster (1.24+)
- Helm 3
- Traefik ingress controller (for challenge instance routing)
- cert-manager (for TLS, optional but recommended)
- A container registry with the built images

### Install

```sh
# Install CRDs first (Tide Instance CRD)
kubectl apply -f charts/crds/

# Install the chart
helm install isolet ./charts \
  --namespace isolet \
  --create-namespace \
  --values charts/values.yaml \
  --values my-values.yaml
```

### Key values.yaml sections

```yaml
event:
  name: "My CTF 2025"
  start: "1739145600"       # unix timestamp
  end: "1741564800"
  domain: ctf.example.com
  publicURL: https://ctf.example.com
  teamSize: 4
  registration:
    emailVerification: true
    passwordReset: true

instances:
  maxConcurrent: 2
  resources:
    requests:
      cpu: "15m"
      memory: "32Mi"
    limits:
      cpu: "50m"
      memory: "128Mi"
  lifetime:
    default: "30m"
    max: "24h"

secrets:
  signing:
    token: <base64>          # JWT signing key
    instance: <base64>       # instance name derivation key
  smtp:
    host: <base64>
    port: <base64>
    user: <base64>
    password: <base64>
  database:
    user: <base64>
    password: <base64>
    host: <base64>
    port: <base64>

certManager:
  dnsProvider: cloudflare
  email: your@email.com

registry:
  type: public               # or "private"
  url: docker.io/thealpha16
```

See [`charts/values.yaml`](./charts/values.yaml) for the full reference with comments.

### Build and push images

```sh
# Build and push a specific service
just build-push oracle
just build-push tide
just build-push pulse

# Build only
just docker-build oracle v2.1.0
```

Each service has a `VERSION` file that is used as the default tag when no tag is specified.

## Configuration

Each service is configured via environment variables. See the service-level documentation for the full list:

- **Oracle** — [`oracle/AGENTS.md`](./oracle/AGENTS.md) (configuration section)
- **Tide** — [`tide/AGENTS.md`](./tide/AGENTS.md)
- **Pulse** — [`pulse/AGENTS.md`](./pulse/AGENTS.md)
- **Herald** — [`herald/AGENTS.md`](./herald/AGENTS.md)
- **Proxy** — [`proxy/AGENTS.md`](./proxy/AGENTS.md)

## Contributing

Each service has an `AGENTS.md` that covers architecture, key patterns, development workflow, and common pitfalls. Start there before touching a service.

```
oracle/AGENTS.md   — API server, Kafka consumer, domain model
tide/AGENTS.md     — Kubernetes operator, CRD, reconciler
pulse/AGENTS.md    — Real-time WebSocket gateway
herald/AGENTS.md   — Event emitter, Kafka producer
ui/AGENTS.md       — Next.js frontend
proxy/AGENTS.md    — Nginx reverse proxy
```

## License

[MIT](./LICENSE)
