# Herald — Agent Guide

## Purpose & Identity

**Herald** is the event backbone of Isolet. It watches infrastructure state changes and emits typed facts to Kafka for downstream consumers (Oracle consumer, Pulse).

Herald is **domain-agnostic and stateless** — it observes what happens and reports it. It has no business logic, no database writes, and no direct user-facing API.

---

## Architecture

Herald runs as one of two completely separate pipeline configurations, chosen at startup via `IDENTITY`. Each configuration has its own source, channel, worker pool, and emitter pointing to a single Kafka topic.

```
  IDENTITY=k8s_source                  IDENTITY=postgres_source

  K8s informer (Instance CRDs)         PostgreSQL logical replication
         │ InstanceFact                        │ Notification
         ▼                                     ▼
  Pipeline (workers)                    Pipeline (workers)
         │                                     │
         ▼                                     ▼
  Kafka emitter                          Kafka emitter
         │                                     │
         ▼                                     ▼
  herald.instance.lifecycle          herald.notifications
  (consumed by Oracle consumer)      (consumed by Pulse)
```

---

## Runtime Identity

Herald runs in one of two modes, set via the `IDENTITY` environment variable:

| `IDENTITY` | Source | Emits | Topic |
|------------|--------|-------|-------|
| `k8s_source` (default) | Kubernetes — controller-runtime informers on `Instance` CRDs | `InstanceExpired` | `herald.instance.lifecycle` |
| `postgres_source` | PostgreSQL — logical replication on `instances` and `endpoints` tables | `Notification` | `herald.notifications` |

The two modes target **different Kafka topics** and serve different downstream consumers. They are not interchangeable — deploy both if you need both topics.

---

## Fact Types

### `InstanceExpired`

Published to `herald.instance.lifecycle`. Consumed by Oracle (consumer mode).

```json
{
  "type": "InstanceExpired",
  "at": "2025-10-15T12:00:00Z",
  "id": "instance-uuid",
  "challenge_id": 42,
  "team_id": 9
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Always `"InstanceExpired"` |
| `at` | timestamp | When the event occurred |
| `id` | string | Instance UUID |
| `challenge_id` | int | Challenge this instance belongs to |
| `team_id` | int\|null | Team that owns the instance (`null` for dynamic challenges) |

---

### `Notification`

Published to `herald.notifications`. Consumed by Pulse for real-time delivery to browser clients.

```json
{
  "type": "Notification",
  "at": "2025-10-15T12:00:00Z",
  "id": "uuid",
  "entity": {
    "name": "instance",
    "id": 200,
    "data": { ... }
  },
  "action": "created",
  "team_ids": [9]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Always `"Notification"` |
| `id` | string | Unique fact UUID |
| `entity` | object\|null | Entity affected (`instance` or `endpoint`) |
| `entity.name` | string | `"instance"` or `"endpoint"` |
| `entity.id` | int | Entity database ID |
| `entity.data` | object | Payload (entity-specific) |
| `action` | string\|null | `"created"`, `"updated"`, or `"deleted"` |
| `message` | string\|null | Human-readable message (used for severity-only notifications) |
| `severity` | string | `"info"`, `"warning"`, or `"success"` (defaults to `"info"` when message is set) |
| `team_ids` | []int | Teams to notify. Empty = broadcast to all. |

Either `entity+action` or `message+severity` must be present. Not both, not neither.

**Kafka partition key logic:**
- Single team → `team:{team_id}` (ordering per team)
- Broadcast → `broadcast` (or `broadcast:{entity}:{id}`)
- Multiple teams → `{entity}:{id}` (or random for distribution)

---

## Configuration

All settings via environment variables (see `utils/config.go`):

| Variable | Default | Description |
|----------|---------|-------------|
| `IDENTITY` | `k8s_source` | Runtime mode: `k8s_source` or `postgres_source` |
| `ENVIRONMENT` | `local` | `local`, `dev`, or `prod` |
| `LOG_LEVEL` | `DEBUG` | Logging verbosity |
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated Kafka broker addresses |
| `KAFKA_RETRIES` | `3` | Kafka producer retry count |
| `INSTANCE_LIFECYCLE_KAFKA_TOPIC` | `herald.instance.lifecycle` | Topic for `InstanceExpired` facts |
| `INSTANCE_LIFECYCLE_WORKERS` | `10` | Pipeline worker count for lifecycle facts |
| `INSTANCE_LIFECYCLE_FACT_CHANNEL_SIZE` | `1024` | Channel buffer size |
| `NOTIFICATION_KAFKA_TOPIC` | `herald.notifications` | Topic for `Notification` facts |
| `NOTIFICATION_WORKERS` | `10` | Pipeline worker count for notifications |
| `NOTIFICATION_FACT_CHANNEL_SIZE` | `1024` | Channel buffer size |
| `EMITTER_RETRIES` | `5` | Kafka emit retry count |
| `EMITTER_RETRY_INTERVAL` | `1s` | Delay between emit retries |
| `K8S_NAMESPACES` | `isolet,dynamic` | Namespaces to watch (k8s_source only) |
| `K8S_KUBE_CONFIG_FILE_PATH` | `` | Path to kubeconfig (empty = in-cluster config) |
| `K8S_DEDUPE_WINDOW` | `60s` | Suppress duplicate events within this window |
| `VALKEY_ADDRESS` | `valkey://localhost:6379` | Valkey for dedup cache |
| `POSTGRES_DSN` | `postgres://...@localhost/isolet` | PostgreSQL DSN (postgres_source only) |
| `POSTGRES_REPLICATION_DSN` | `...?replication=database` | Replication slot DSN (postgres_source only) |
| `POSTGRES_STANDBY_TIMEOUT` | `5s` | Replication standby timeout |
| `SENTRY_DSN` | `` | Sentry error reporting (optional) |

---

## Development

```sh
cd herald

# k8s source (default) — requires a kubeconfig with access to Instance CRDs
go run main.go

# Postgres source
IDENTITY=postgres_source go run main.go

# With hot reload (requires air)
air
```

Herald will exit if it cannot connect to Kafka or initialize its source. Check connectivity before running.

---

## Key Patterns

### k8s source — what triggers an emit

The K8s informer receives CREATE, UPDATE, and DELETE events but **only UPDATE events are processed** (create/delete are explicitly ignored). Of those, only transitions to `Phase == Expired` result in a fact. The source does nothing for any other phase change.

### De-duplication (k8s source)

Because K8s can deliver the same update event multiple times, the k8s source uses a Valkey `SET NX` with TTL (`K8S_DEDUPE_WINDOW`, default 60s) keyed as `herald:k8s:Expired:{instance-uid}`. If the key already exists, the event is silently dropped.

### postgres source — setup requirements

The postgres source uses PostgreSQL logical replication and requires:
- `wal_level = logical` in `postgresql.conf`
- A replication slot named `herald_slot` (auto-created by Herald on startup)
- A publication named `herald_pub` covering `instances` and `endpoints` tables (auto-created by Herald on startup)

Herald uses the `pgoutput` plugin. The `POSTGRES_REPLICATION_DSN` must include `?replication=database`.

### Pipeline workers

Facts flow from the source channel through a configurable pool of workers, each calling the emitter. Channel size is configurable to absorb bursts. Workers for the two modes are sized independently (`INSTANCE_LIFECYCLE_WORKERS` and `NOTIFICATION_WORKERS`).

### Retry logic

The pipeline retries Kafka emission up to `EMITTER_RETRIES` times with a **progressive backoff** of `attempt × EMITTER_RETRY_INTERVAL` (e.g. 1s, 2s, 3s, ...). After all retries are exhausted, the fact is dropped and the error is sent to Sentry.

---

## Code Navigation

| Path | Contents |
|------|----------|
| `main.go` | Entry point, identity switch, wiring |
| `internal/sources/k8s/` | K8s informer source, Instance CR handlers |
| `internal/sources/postgres/` | PostgreSQL logical replication source |
| `internal/pipeline/` | Worker pool, validation, routing |
| `internal/emitter/kafka/` | Kafka producer |
| `pkg/facts/` | Fact types: `InstanceFact`, `Notification`, `BaseFact` |
| `utils/config.go` | Full configuration struct |
| `utils/constants.go` | Identity constants (`SourceK8s`, `SourcePostgres`) |

---

## Adding a New Fact Type

1. Define the struct in `pkg/facts/` implementing the `Fact` interface (`FactType`, `Key`, `OccuredAt`, `Validate`, `Marshal`)
2. Add the new type to `allowedTypes` in `pkg/facts/fact.go`
3. Emit the fact from the appropriate source handler (`internal/sources/k8s/handlers.go` or postgres equivalent)
4. If the new fact needs its own Kafka topic, add a new `AppConfig` block in `main.go` with a new `RunPipeline` call — the pipeline itself does not route; each pipeline is hardwired to one topic
5. Add the topic env var to `utils/config.go`
6. Update downstream consumers (Oracle, Pulse) to handle the new type

---

## Related Services

- **Tide** — creates `Instance` CRs that Herald watches (k8s source)
- **Oracle (consumer)** — consumes `InstanceExpired` facts from `herald.instance.lifecycle`
- **Pulse** — consumes `Notification` facts from `herald.notifications` and delivers them to browser clients
