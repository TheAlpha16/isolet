# Pulse ⚡

**Pulse** is Isolet’s realtime delivery service.

It consumes domain events from Kafka and broadcasts **realtime updates** to connected clients via WebSockets.

---

## 🧠 What Pulse Is

Pulse is a **stateless, event-driven realtime gateway**.

It:

* Maintains WebSocket connections with clients
* Authenticates users via JWT + Redis session validation (Oracle invariants)
* Consumes **Notification facts** from Kafka
* Derives a **semantic event name** dynamically
* Forwards events to the correct audience (team / global)

It does **not**:

* Store state
* Implement business logic
* Act as a source of truth
* Reshape or interpret payloads

---

## 🏗️ Architecture

```
Herald → Kafka (Notification facts) → Pulse → Clients (WebSocket)
```

* **Herald** emits canonical domain events (`Notification`)
* **Kafka** ensures ordering and delivery
* **Pulse** adds minimal context and broadcasts
* **Clients** interpret and act on events

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

* `team:<team_id>` → team-scoped events
* `global` → broadcast events

All channels are **read-only** for clients. Clients cannot push business events directly to Pulse.

---

## 🔐 Authentication

* **JWT Verification**: Clients must provide a valid JWT.
* **Redis Session Validation**: Pulse verifies the token against Redis (`redix`) to ensure the session is still valid in Oracle.
* **Context**: User ID, Team ID, and Role are extracted from claims and assigned to the socket.

---

## 🔁 Event Flow

1.  **Event Emitted**: Herald publishes a `Notification` to Kafka.
2.  **Consumption**: Pulse (`:brod`) consumes the message.
3.  **Processing**: Pulse decodes JSON and derives the semantic `event`.
4.  **Broadcast**: Pulse publishes the enrichment payload to the appropriate Phoenix Channel.
5.  **Delivery**: Connected WebSocket clients receive the update.

---

## 🎯 Design Goals

* **Low latency** — near-instant updates
* **Scalable** — thousands of concurrent connections
* **Decoupled** — 100% domain-agnostic; zero maintenance for new business rules
* **Stable** — resilient to malformed messages or missing fields

---

## 🔌 Example Client Usage

```javascript
import { Socket } from "phoenix"

const socket = new Socket("/socket", {
  params: { token: "<JWT>" }
})

socket.connect()

const channel = socket.channel("team:9")
channel.on("notification", (payload) => {
  console.log("Received event:", payload.event, payload)
})

channel.join()
```
