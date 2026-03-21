# Pulse ⚡

**Pulse** is Isolet’s realtime delivery service.

It consumes domain events from Kafka and broadcasts **realtime updates** to connected clients via WebSockets.

---

## 🧠 What Pulse Is

Pulse is a **stateless, event-driven realtime gateway**.

It:

* maintains WebSocket connections with clients
* authenticates users via JWT
* consumes **Notification facts** from Kafka
* derives a **semantic event name**
* forwards events to the correct audience (team / global)

It does **not**:

* store state
* implement business logic
* act as a source of truth
* reshape or interpret payloads

> Pulse is a **delivery layer**, not a decision-maker.

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
    "data": {
      "id": 200,
      "challenge_id": 5,
      "team_id": 9
    }
  },
  "action": "created",
  "team_ids": [9]
}
```

These are:

* canonical
* system-oriented
* multi-consumer

---

### 2️⃣ Event Derivation

Pulse derives a **semantic event name** from the Notification:

```
event = entity.name + "." + action
```

Example:

```
instance.created
instance.updated
instance.deleted
```

For message-based notifications:

```
event = "notification." + severity
```

Examples:

```
notification.info
notification.warning
notification.success
```

---

### 3️⃣ Delivered Event Shape

Pulse **does not modify the payload**.
It only adds an `event` field.

```json
{
  "type": "Notification",
  "id": "uuid",
  "at": "timestamp",

  "event": "instance.created",

  "entity": {
    "name": "instance",
    "id": 200,
    "data": {
      "id": 200,
      "challenge_id": 5,
      "team_id": 9
    }
  },

  "action": "created",
  "message": null,
  "severity": null,
  "team_ids": [9]
}
```

---

### 4️⃣ Channels

Pulse uses Phoenix Channels for routing.

* `team:<team_id>` → team-scoped events
* `global` → broadcast events

Clients subscribe to relevant topics and receive updates instantly.

---

## 🔐 Authentication

* Clients authenticate via **JWT on connection**
* Unauthorized connections are rejected
* Team context is derived from token claims

---

## 🔁 Event Flow

1. Instance changes state

2. Herald emits a `Notification` fact

3. Kafka partitions event by `team:<id>`

4. Pulse consumes event

5. Pulse:

   * derives `event` (`entity.name + action`)
   * forwards event as-is
   * broadcasts to `team:<id>`

6. Clients interpret and update UI

---

## 📈 Ordering Guarantees

Pulse relies on Kafka partitioning:

* Events for a team share the same key: `team:<id>`
* Kafka preserves ordering within a partition
* Pulse processes events sequentially per partition

> This ensures **consistent ordering per team** without additional coordination.

---

## 🎯 Design Goals

* **Low latency** — near-instant updates
* **Scalable** — thousands of concurrent connections
* **Lightweight** — minimal resource usage
* **Decoupled** — no dependency on core logic
* **Stable** — no changes required for new entities

---

## 🚫 Non-Goals

* No persistence
* No business rules
* No payload transformation
* No state reconciliation
* No guaranteed delivery

Missed events can always be recovered via API refresh.

---

## 🔌 Example Client Usage

```javascript
import { Socket } from "phoenix"

const socket = new Socket("/socket", {
  params: { token: "<JWT>" }
})

socket.connect()

const channel = socket.channel("team:9")

channel.on("notification", (event) => {
  dispatch(event)
})

channel.join()
```

---

## 📊 Scaling Model

* Horizontally scalable (multiple Pulse instances)
* Stateless — no shared state between nodes
* Kafka handles distribution
* Phoenix PubSub handles fan-out

---

## 🧭 Philosophy

Pulse follows a simple principle:

> **Facts are canonical. Events are derived. UI assigns meaning.**

* Kafka carries **facts**
* Pulse adds **event context**
* UI defines **behavior**
* API remains the **source of truth**

---

## 🧩 Role in Isolet

Pulse complements:

* **Oracle** → API & business logic
* **Herald** → fact production
* **Tide** → instance orchestration

---

## 🏁 Summary

Pulse is:

* simple
* fast
* stable
* predictable

It exists to do one thing well:

> **deliver realtime facts with minimal transformation and let clients decide what they mean.**
