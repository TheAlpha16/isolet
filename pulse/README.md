# Pulse ⚡

**Pulse** is Isolet’s realtime delivery service.

It consumes domain events from Kafka and broadcasts **realtime updates** to connected clients via WebSockets.

---

## 🧠 What Pulse Is

Pulse is a **stateless, event-driven realtime gateway**.

It:

- Maintains WebSocket connections with clients
- Authenticates users via JWT + Redis session validation (Oracle invariants)
- Consumes **Notification facts** from Kafka
- Derives a **semantic event name** dynamically
- Forwards events to the correct audience (team / global)

It does **not**:

- Store state
- Implement business logic
- Act as a source of truth
- Reshape or interpret payloads

---

## 🏗️ Architecture

```
Herald → Kafka (Notification facts) → Pulse Pod A (consumer)
                                          ↓
                                    Redis PubSub
                                    /     |     \
                                   /      |      \
                              Pod A    Pod B    Pod C
                                |        |        |
                          WebSocket  WebSocket  WebSocket
                            Clients   Clients   Clients
```

- **Herald** emits canonical domain events (`Notification`)
- **Kafka** ensures ordering and at-least-once delivery to consumer group
- **Pulse Pod** (single pod in consumer group) consumes event and publishes to Redis PubSub
- **Redis PubSub** fans out the event to **all Pulse pods** in the cluster
- **Each Pulse Pod** broadcasts to its locally connected WebSocket clients
- **Clients** interpret and act on events

### Why Redis PubSub?

In Kubernetes with multiple Pulse pods:

- Kafka consumer groups ensure **only one pod** receives each event
- WebSocket clients may be connected to **different pods**
- Redis PubSub provides **cluster-wide distribution** without Erlang clustering
- All pods receive all events and push to their local clients

This architecture ensures:

- ✅ No missed events across pods
- ✅ Horizontal scalability
- ✅ Works in Kubernetes without node discovery
- ✅ Clean separation: Kafka for ingestion, PubSub for distribution

---

## 📡 Core Concepts

### 1️⃣ Notification Facts (Input)

Pulse consumes generic domain events of type `Notification`.

Example:

```json
{
  "type": "Notification",
  "id": "uuid",
  "at": "timestamp",
  "entity": {
    "name": "instance",
    "id": 200,
    "data": { ... }
  },
  "action": "created",
  "team_ids": [9]
}
```

These are **canonical**, **system-oriented**, and **multi-consumer**.

---

### 2️⃣ Dynamic Event Derivation

Pulse derives a **semantic event name** dynamically from the Notification:

1.  If `severity` is present: `notification.<severity>` (e.g., `notification.warning`)
2.  If `entity` and `action` are present: `<entity>.<action>` (e.g., `instance.created`)
3.  Fallback: `unknown.unknown`

This design is **domain-agnostic**; new entities or actions added to Herald are automatically supported by Pulse without code changes.

---

### 3️⃣ Delivered Event Shape

Pulse **does not modify the payload**. It only adds an `event` field for client-side routing.

---

### 4️⃣ Channels

- `team:<team_id>` → team-scoped events
- `global` → broadcast events

All channels are **read-only** for clients. Clients cannot push business events directly to Pulse.

---

## 🔐 Authentication

- **JWT Verification**: Clients must provide a valid JWT.
- **Redis Session Validation**: Pulse verifies the token against Redis (`redix`) to ensure the session is still valid in Oracle.
- **Context**: User ID, Team ID, and Role are extracted from claims and assigned to the socket.

---

## 🔁 Event Flow

1.  **Event Emitted**: Herald publishes a `Notification` to Kafka.
2.  **Consumption**: One Pulse pod (`:brod` consumer group) consumes the message.
3.  **Processing**: That pod decodes JSON and derives the semantic `event`.
4.  **PubSub Broadcast**: That pod publishes to Redis PubSub (`Phoenix.PubSub.Redis`).
5.  **Cluster Distribution**: All Pulse pods receive the PubSub message.
6.  **Local Broadcast**: Each pod broadcasts to its Phoenix Channels (`team:<id>` or `global`).
7.  **Delivery**: Connected WebSocket clients on all pods receive the update.

---

## 🎯 Design Goals

- **Low latency** — near-instant updates
- **Scalable** — thousands of concurrent connections
- **Decoupled** — 100% domain-agnostic; zero maintenance for new business rules
- **Stable** — resilient to malformed messages or missing fields

---

## 🔌 Example Client Usage

```javascript
import { Socket } from "phoenix";

const socket = new Socket("/socket", {
	params: { token: "<JWT>" },
});

socket.connect();

const channel = socket.channel("team:9");
channel.on("notification", (payload) => {
	console.log("Received event:", payload.event, payload);
});

channel.join();
```

---

## 🚀 Deployment & Configuration

### Environment Variables

- `REDIS_URL` — Redis connection URL (required, used for sessions and PubSub)
- `HOSTNAME` or `POD_NAME` — Unique identifier for this pod/instance (auto-set in Kubernetes, used for Redis PubSub node identification)
- `KAFKA_BROKERS` — Kafka broker addresses (default: `localhost:9092`)
- `KAFKA_TOPIC` — Topic to consume (default: `herald.notifications`)
- `KAFKA_GROUP_ID` — Consumer group ID (default: `pulse`)
- `JWT_SECRET` — Secret for JWT verification (required)

> **Note**: Redis PubSub requires each node to have a unique identifier. In Kubernetes, `HOSTNAME` is automatically set to the pod name. For local development, a unique identifier is generated automatically.

### Multi-Pod Kubernetes Deployment

Pulse is designed to run with multiple replicas in Kubernetes:

1. **Single Kafka Consumer Group**: All pods share the same `KAFKA_GROUP_ID`, ensuring each event is consumed by exactly one pod
2. **Redis PubSub Distribution**: The consuming pod publishes to Redis, which fans out to all pods
3. **No Erlang Clustering Required**: Uses Redis PubSub instead of distributed Erlang, simplifying Kubernetes deployment

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: pulse
spec:
    replicas: 3 # Multiple pods for HA
    template:
        spec:
            containers:
                - name: pulse
                  env:
                      - name: REDIS_URL
                        value: "redis://redis-service:6379"
                      - name: KAFKA_BROKERS
                        value: "kafka-broker:9092"
```

### Local Development

For local development, tests use a **local PubSub adapter** (no Redis required for PubSub):

```elixir
# config/test.exs
config :pulse, pubsub_adapter: :local
```

In production/dev, the Redis adapter is used by default:

```elixir
# lib/pulse/application.ex
# Configures Redis PubSub unless overridden
pubsub_config = [
  name: Pulse.PubSub,
  adapter: Phoenix.PubSub.Redis,
  url: redis_url
]
```

---

## 🔧 Monitoring

Access LiveDashboard in development:

```
http://localhost:4000/dev/dashboard
```

View:

- Phoenix metrics (requests, channels, sockets)
- VM metrics (memory, processes, schedulers)
- Telemetry data
