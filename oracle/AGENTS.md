# 🧠 Oracle — Service Guide for AI Agents

## 🎯 Purpose & Identity

**Oracle** is the **central control plane** and **decision-maker** for Isolet.

**Core Responsibilities:**
- User authentication & authorization
- Challenge lifecycle management
- Instance orchestration (create/start/stop/extend)
- Scoreboard & team scoring
- Event processing (facts from Herald via Kafka)
- Real-time configuration management
- Team & user management

**Key Philosophy:**
- Never crash during a live CTF event
- Event-driven architecture
- Dual-mode operation (REST API + Kafka Consumer)
- Isolation-first design

---

## 🏗️ Architecture Overview

### Deployment Modes

Oracle runs in **two distinct identities** (configured via `IDENTITY` env var):

1. **`rest`** — HTTP API server for user interactions
2. **`consumer`** — Kafka consumer processing facts from Herald

Both modes share the same codebase but initialize different delivery layers.

### Layered Architecture

```
┌─────────────────────────────────────────────────────┐
│  Delivery Layer (REST/Consumer)                     │
│  • REST: Fiber HTTP handlers + middleware           │
│  • Consumer: Kafka message processing               │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  Usecase Layer (Business Logic)                     │
│  • Orchestrates domain logic                        │
│  • Coordinates between repos/infra/external         │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  Domain Layer (Business Entities + Interfaces)      │
│  • Domain models & validation                       │
│  • Usecase/Repository/Service interfaces            │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  Repository Layer (Data Persistence)                │
│  • PostgreSQL via GORM                              │
│  • Cache abstraction (Valkey)                       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  Infrastructure & External Services                 │
│  • Infra: JWT, CNC (pub/sub), K8s client           │
│  • External: Instance service (→ Tide)              │
└─────────────────────────────────────────────────────┘
```

---

## 📦 Domain Model

Oracle manages **13 core domains** (see `internal/domain/`):

### 🔐 User & Access Control

| Domain | Purpose | Key Operations |
|--------|---------|----------------|
| **user** | User accounts | Create, get by ID/email, verify, update password |
| **auth** | Authentication | Login, register, verify email, password reset, token generation |
| **team** | Team management | Create, join via invite, get members, update team |
| **token** | JWT lifecycle | Generate/validate auth/realtime/email/password/invite tokens |

### 🎮 Challenge & Scoring

| Domain | Purpose | Key Operations |
|--------|---------|----------------|
| **challenge** | Challenge CRUD | Get challenges, check requirements, unlock hints, validate flags |
| **score** | Scoreboard & scoring | Get scoreboard/graph, update team scores, track solves |
| **instance** | Container lifecycle | Start/stop/extend instances, track endpoints |
| **manifest** | Deployment specs | Define container resources, ports, images |
| **fact** | Event processing | Handle instance state updates from Herald |

### ⚙️ System & Support

| Domain | Purpose | Key Operations |
|--------|---------|----------------|
| **event** | Event metadata | Get event info (start/end times, visibility) |
| **profile** | User profiles | Get user/team/challenge data for dashboards |
| **email** | Email dispatch | Send verification, password reset, team invites |
| **configvars** | Dynamic config | Fetch event timings, flags, settings from DB/cache |

---

## 🔌 Key Integrations

### External Services (`external/`)

**Instance Service** (`external/instance/`)
- Communicates with **Tide** (K8s operator)
- Operations: `CreateInstance`, `DeleteInstance`, `GetInstanceStatus`
- Used by Instance usecase to orchestrate container lifecycle

### Infrastructure Services (`infra/`)

**JWT** (`infra/jwt/`)
- Signs and validates JWT tokens
- Used for auth, team invites, email verification

**CNC** (Command & Control, `infra/cnc/`)
- Pub/sub over Valkey
- Broadcasts config updates, event changes
- Real-time coordination between Oracle instances

**K8s Client** (`infra/k8s/`)
- Direct Kubernetes API access (if needed)
- Typically used via Instance service

**SMTP** (`infra/smtp/`)
- Email sending via background workers
- Retry logic + timeout handling

### Data Layer

**PostgreSQL** (`infra/database/postgres/`)
- Primary data store
- GORM-based models
- Auto-migration in local environment

**Valkey** (`infra/database/valkey/`)
- High-performance cache
- Used for: sessions, scoreboard, challenge cache, instance state
- Pub/sub for CNC

**Cache Abstraction** (`infra/cache/`)
- Wrapper over Valkey
- Provides: Get, Set, Delete, Publish operations
- TTL management

---

## 🚀 REST API Endpoints

All routes are prefixed with `/api/v1` (configurable via `REST_API_VERSION_PREFIX`).

### Authentication (`/auth`)
- `POST /login` — User login
- `POST /register` — User registration
- `POST /verify` — Verify email
- `POST /forgot-password` — Request password reset
- `POST /reset-password` — Reset password
- `POST /logout` — Logout

### Teams (`/team`)
- `POST /create` — Create new team
- `POST /join` — Join team via invite token
- `GET /members` — Get team members
- `PATCH /update` — Update team details

### Challenges (`/challenge`)
- `GET /list` — List visible challenges
- `GET /:id` — Get challenge details
- `POST /:id/unlock-hint` — Unlock hint
- `POST /:id/submit` — Submit flag

### Instances (`/instance`)
- `GET /list` — List team instances
- `POST /start` — Start challenge instance
- `POST /stop` — Stop instance
- `POST /extend` — Extend instance lifetime

### Scoreboard (`/score`)
- `GET /scoreboard` — Get current scoreboard
- `GET /graph` — Get score graph (top 10 teams)

### Profile (`/profile`)
- `GET /` — Get user profile data

### Event (`/event`)
- `GET /info` — Get event info (timings, visibility)

### Health
- `GET /ping` — Health check

---

## 📡 Kafka Consumer (Facts)

When running in **consumer mode** (`IDENTITY=consumer`), Oracle:

1. **Subscribes to Kafka topics** containing facts from Herald
2. **Deserializes facts** using msgpack codec
3. **Routes to Fact usecase** for processing
4. **Updates instance state** in database

### Fact Types

Defined in `herald/pkg/facts/`:

- **`InstanceExpiredFactType`** — Instance TTL expired
- More fact types can be added for instance health, errors, etc.

### Processing Flow

```
Herald → Kafka → Oracle Consumer → Fact Usecase → Instance Usecase → DB + Cache
```

**Error Handling:**
- Failed facts are logged and sent to Sentry
- Metrics track success/failure rates per fact type
- No message reprocessing (at-most-once delivery)

---

## 🔧 Configuration

Oracle is configured via **environment variables** (see `utils/config.go`).

### Critical Settings

| Variable | Default | Purpose |
|----------|---------|---------|
| `IDENTITY` | `rest` | Runtime mode: `rest` or `consumer` |
| `ENVIRONMENT` | `local` | `local`, `dev`, or `prod` |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `VALKEY_ADDRESS` | `valkey://localhost:6379` | Valkey connection |
| `TOKEN_SIGNING_KEY` | `trustmebro` | JWT signing key (⚠️ change in prod!) |
| `INSTANCE_NAMESPACE` | `isolet` | K8s namespace for instances |
| `INSTANCE_START_TIMEOUT` | `30s` | Max time for instance start |
| `INSTANCE_STOP_TIMEOUT` | `10s` | Max time for instance stop |
| `SENTRY_DSN` | (empty) | Sentry error tracking |
| `LOG_LEVEL` | `DEBUG` | Logging verbosity |

### Dynamic Configuration (ConfigVars)

Many settings are **fetched from the database** at runtime:
- Event start/end times
- Event visibility flags
- Max team size
- Instance TTL
- Email templates

These are cached and refreshed periodically via `ConfigVars` usecase.

---

## 🛠️ Development Workflow

### Running Locally

**REST mode:**
```bash
cd oracle
export IDENTITY=rest
go run main.go
```

**Consumer mode:**
```bash
export IDENTITY=consumer
go run main.go
```

**With Air (hot reload):**
```bash
air  # Uses .air.toml configuration
```

### Database Migrations

Auto-migration runs in **local environment only**:
```go
if utils.GetConfig().Environment == utils.LOCAL {
    err = repository.AutoMigrate(dbPool)
}
```

For production, use explicit migration scripts.

### Testing

(Check `test/` directory for test structure)

---

## 📊 Observability

### Metrics

Exposed on `/metrics` endpoint (default port: `6969`).

**Custom metrics** (defined in `utils/tracer/`):
- `facts_total` — Count of processed facts by topic/status/type
- `fact_processing_duration_seconds` — Fact processing latency
- HTTP request metrics via middleware

### Tracing

**OpenTelemetry** integration:
- Traces HTTP requests (via `otelfiber` middleware)
- Traces span context propagated through `context.Context`
- Consumer events traced independently

### Error Tracking

**Sentry** integration:
- Automatic error capture
- Context-aware error reporting via `errorDom.RaiseToSentry()`
- Sample rates configurable per environment

### Logging

**Structured logging** via Zap:
- Request/response logging
- Error logging with context
- Configurable log levels

---

## 🧩 Key Patterns & Conventions

### Error Handling

Oracle uses **domain-specific errors** (`internal/domain/errors/`):

```go
// Raise an error with context
err := errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "instance not found", nil, common.ExtraData{"instance_id": id})

// Send to Sentry
errorDom.RaiseToSentry(ctx, err)
```

### Context Propagation

Every usecase/handler accepts `context.Context`:
- Carries request ID, user info, span context
- Used for cancellation and timeouts
- Set via middleware in REST mode

### Cache Keys

Standardized cache key patterns:
- `instance:{teamID}:{challengeID}` — Instance data
- `challenge:{id}` — Challenge details
- `scoreboard` — Full scoreboard
- See domain files for `*CacheKey()` functions

### Middleware Stack (REST)

Applied in order:
1. `ContextMiddleware` — Injects request context
2. `MetricsMiddleware` — Records HTTP metrics
3. `otelfiber.Middleware` — OpenTelemetry tracing
4. `IPMiddleware` — Extracts client IP
5. `SentryMiddleware` — Captures errors
6. `LoggingMiddleware` — Logs requests/responses
7. `ErrorMiddleware` — Standardizes error responses

**Route-specific middleware:**
- `AuthMiddleware` — Validates JWT, injects user
- `RequireTeamMiddleware` — Ensures user is in a team
- `CheckTimingsMiddleware` — Validates event is active
- `ContextDeadlineMiddleware` — Sets timeout for long operations

---

## 🚨 Critical Paths (Don't Break!)

### Instance Start Flow

**Most complex operation** — must be reliable:

1. **Validate request** (team has access, challenge exists, etc.)
2. **Check existing instance** (prevent duplicates)
3. **Get/create manifest** (deployment spec)
4. **Generate flag** (if dynamic)
5. **Call external Instance service** (→ Tide → K8s)
6. **Create DB records** (instance + endpoints)
7. **Update cache**
8. **Return endpoints to user**

**Failure points:**
- K8s API timeout → handled by deadline middleware
- Duplicate instance → checked via cache + DB
- Manifest missing → validated before creation

### Fact Processing (Consumer)

**Must be idempotent** — facts may arrive multiple times:

1. **Deserialize fact** from Kafka message
2. **Route to Fact usecase** (`HandleEvent`)
3. **Update instance state** (e.g., mark expired)
4. **Update cache** (invalidate stale data)
5. **Log + metrics**

**Failure modes:**
- Deserialization error → skip message, log error
- DB update failure → retry not implemented (investigate if needed)
- Cache update failure → eventually consistent

### Scoreboard Updates

**High contention** — multiple teams solving simultaneously:

1. **Update team score in cache** (atomic increment)
2. **Update challenge solve count**
3. **Invalidate scoreboard cache**
4. **Background refresh** if cache fails
5. **Broadcast update via CNC** (notify Pulse)

---

## 🐛 Common Issues & Fixes

### "Instance already exists"
- **Cause:** Cache/DB out of sync
- **Fix:** Check Valkey connection, verify instance cleanup

### "Failed to connect to valkey"
- **Cause:** Valkey not running or wrong address
- **Fix:** Check `VALKEY_ADDRESS`, verify Valkey is reachable

### "Database connection failed"
- **Cause:** PostgreSQL not running or wrong credentials
- **Fix:** Verify `DB_HOST`, `DB_USER`, `DB_PASSWORD`

### "Token verification failed"
- **Cause:** Mismatched signing keys between instances
- **Fix:** Ensure `TOKEN_SIGNING_KEY` is consistent across all Oracle instances

### Consumer not processing facts
- **Cause:** Kafka broker unreachable or topic doesn't exist
- **Fix:** Verify Kafka connection, check topic creation

---

## 📚 Code Navigation Tips

### Finding a Feature

1. **Start with domain** (`internal/domain/{feature}/`)
   - Check `interface.go` for contracts
   - Check domain model files for entities

2. **Check usecase** (`internal/usecase/{feature}/`)
   - Business logic lives here
   - Orchestrates repos, infra, external services

3. **Find delivery** (`delivery/rest/handler/{feature}/` or `delivery/consumer/`)
   - REST handlers or Kafka consumers
   - Input validation and response formatting

### Adding a New Endpoint

1. **Define DTOs** in `internal/domain/{feature}/dto.go`
2. **Add method to Usecase interface** in `internal/domain/{feature}/interface.go`
3. **Implement in usecase** in `internal/usecase/{feature}/`
4. **Create handler** in `delivery/rest/handler/{feature}/`
5. **Register route** in `delivery/rest/routes/{feature}.go`
6. **Add middleware** as needed (auth, timing checks, etc.)

### Adding a New Fact Type

1. **Define fact struct** in `herald/pkg/facts/`
2. **Add deserializer** in `delivery/consumer/codec.go`
3. **Handle in Fact usecase** (`internal/usecase/fact/`)
4. **Update instance state** as needed

---

## 🧪 Testing Strategy

### Unit Tests
- Test usecase logic in isolation
- Mock repository/infra dependencies
- Focus on business logic validation

### Integration Tests
- Test full request flow (handler → usecase → repo)
- Use test database/cache
- Verify side effects (DB writes, cache updates)

### E2E Tests
- Run full stack (API + Consumer + Dependencies)
- Test critical paths (instance start, scoreboard updates)
- Validate observability (metrics, traces, logs)

---

## 🎓 Onboarding Checklist for New AI Agents

- [ ] Understand the **dual-mode** architecture (REST vs Consumer)
- [ ] Know the **13 core domains** and their responsibilities
- [ ] Trace the **instance start flow** end-to-end
- [ ] Understand **fact processing** from Herald
- [ ] Review **error handling** patterns
- [ ] Check **middleware stack** in REST mode
- [ ] Learn **cache key conventions**
- [ ] Explore **ConfigVars** for dynamic configuration
- [ ] Review **observability setup** (metrics, traces, logs)
- [ ] Read **critical paths** section to avoid breaking prod

---

## 🔗 Related Services

- **Herald** — Emits facts about instance lifecycle events
- **Tide** — Kubernetes operator, manages actual instance pods
- **Pulse** — Real-time WebSocket server, receives CNC broadcasts
- **UI** — Frontend, consumes Oracle REST API

---

## 📝 TODO / Future Enhancements

(From main.go comments)
- Add metrics to herald, tide
- Don't init full usecases in case of listener and consumer (optimize initialization)
- Implement retry logic for fact processing
- Add circuit breaker for external Instance service calls
- Consider read replicas for scoreboard queries

---

## ✅ Summary for AI Agents

**When working on Oracle:**

1. **Identify the mode** (REST or Consumer) — different entry points
2. **Start at the domain layer** — understand the entities and interfaces
3. **Check usecases** — business logic orchestration
4. **Trace dependencies** — repos, infra, external services
5. **Test critical paths** — instance start, fact processing, scoreboard updates
6. **Preserve idempotency** — especially in consumer mode
7. **Follow error patterns** — use `errorDom.Raise()` and `RaiseToSentry()`
8. **Update observability** — add metrics, traces, logs for new features
9. **Document config** — new env vars or ConfigVars entries
10. **Never break prod** — test thoroughly, especially instance lifecycle

Oracle is the **brain** — respect its critical role in the platform. 🧠
