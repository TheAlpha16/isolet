# Tide — Agent Guide

## Purpose & Identity

**Tide** is the infrastructure execution layer of Isolet. It is a Kubernetes operator built with [Kubebuilder](https://book.kubebuilder.io/) that manages the lifecycle of CTF challenge instances.

Upstream services (Oracle, Herald) describe what they need by creating `Instance` custom resources. Tide reconciles those CRs into real Kubernetes workloads and keeps their status up to date.

Tide is **stateless and domain-agnostic**. It has no knowledge of teams, scores, or CTF rules. Those live in Oracle.

---

## What Tide Creates Per Instance

For every `Instance` CR, Tide manages:

| Resource | Name pattern | Notes |
|----------|-------------|-------|
| `Deployment` | `deployment-{instance-name}` | Single container named `"challenge"` running the challenge image |
| `Service` | `svc-{instance-name}` | `ClusterIP`, one port per endpoint |
| `IngressRoute` | `ir-{instance-name}` | Created **only** for `http`/`https` endpoints |

For `nc`/`ssh` endpoints: a Service port is created, but **no IngressRoute**. Traefik TCP routing is not currently wired up — `nc`/`ssh` challenges are cluster-internal only via the Service.

---

## Instance Lifecycle

```
              ┌─────────────┐
              │   Pending   │◄────────────────────┐
              └──────┬──────┘                     │ (deployment becomes unready)
                     │ deployment + service +      │
                     │ ingress ready               │
                     ├──────────────────────┐      │
                     │                      │      │
                     │ AvailableAt           │ no  │
                     │ in future            │ AvailableAt
                     ▼                      ▼      │
              ┌─────────────┐       ┌─────────────┐
              │   Staged    │──────►│   Running   │
              └──────┬──────┘ AvailableAt └──┬──┬─┘
                     │        passed         │  │
                     │                       │  │ ExpiresAt passed
                     │ pod failure           │  ▼
                     ▼                    ┌─────────────┐
              ┌─────────────┐             │   Expired   │──► deleted
              │   Failed    │◄────────────┘             │
              └──────┬──────┘ pod failure   DeletionTimestamp set
                     │                      ▼
                     │ all resources  ┌─────────────┐
                     └───────────────►│ Terminated  │──► finalizer removed
                       become ready   └─────────────┘
```

| Phase | Meaning |
|-------|---------|
| `Pending` | CR created or deployment not yet ready |
| `Staged` | Deployment ready but `availableAt` is still in the future |
| `Running` | All resources ready and `availableAt` has passed (or not set) |
| `Failed` | Pod entered a failure state (see failure reasons below) |
| `Expired` | `expiresAt` was reached; CR is deleted shortly after |
| `Terminated` | `DeletionTimestamp` was set (manual deletion or SDK `DeleteInstance`) |

**Phase recovery:** if an instance is `Failed` but all resources subsequently become ready (e.g., transient image pull failure resolved), it transitions back to `Running`.

**Requeue scheduling:** Tide self-schedules requeue based on `availableAt` and `expiresAt`. No external TTL controller or cron job is needed.

### Pod failure reasons that set `PhaseFailed`

`CrashLoopBackOff`, `ErrImagePull`, `ImagePullBackOff`, `CreateContainerConfigError`, `InvalidImageName`, `CreateContainerError`, or non-zero exit code.

---

## Challenge Types

| Type | Description |
|------|-------------|
| `dynamic` | Shared, long-running instance available to all teams. `Team` field is nil. Label `challenges.isolet.dev/team = "dynamic"`. |
| `on-demand` | Isolated, short-lived instance for a specific team. `Team` field must be set (validated by webhook). |

---

## Protocol Support

| Protocol | Resources created | Notes |
|----------|------------------|-------|
| `http` | Deployment + Service + IngressRoute | Traefik routes on `Host()` rule |
| `https` | Deployment + Service + IngressRoute | Same as http; TLS cert from secret `challenge-certs` |
| `nc` | Deployment + Service | No IngressRoute; cluster-internal only |
| `ssh` | Deployment + Service | No IngressRoute; cluster-internal only |

### IngressRoute details

- **Entry points:** always `["web", "websecure"]`
- **TLS:** always configured with `secretName: "challenge-certs"` (must exist in the instance namespace)
- **Hostname format:**
  - Single HTTP/HTTPS endpoint: `{instance-name}.{challenge-slug}.{domain}`
  - Multiple HTTP/HTTPS endpoints: `{endpoint-name}-{instance-name}.{challenge-slug}.{domain}`
- IngressRoute is created with `r.Create` / `r.Update`, **not** Server-Side Apply (unlike Deployment and Service)

### Flag injection

If `spec.challenge.flag` is set, it is injected into the container as the `FLAG` environment variable.

---

## Webhooks

Two webhooks fire on Instance create/update:

### Mutating (defaulter)
- If `spec.lifecycle` is nil, initializes it with `allowExtension: true`

### Validating
Enforces on create:
- `challenge.domain` must be set
- `on-demand` challenges must have `team` set
- Endpoint names must be unique
- Endpoint ports must be in range 1–65535, no duplicates
- Endpoint protocol must be one of `http`, `https`, `nc`, `ssh`
- `lifecycle.expiresAt` must be in the future if set
- `lifecycle.availableAt` must be before `lifecycle.expiresAt` if both are set

Enforces on update (immutability):
- `challenge.id`, `challenge.slug`, `challenge.image`, `challenge.type`, `challenge.domain` are **immutable**
- `team` is **immutable**
- `lifecycle.expiresAt` can only be extended if `lifecycle.allowExtension == true`

---

## Status Conditions

Three conditions are actively maintained:

| Condition | Set `True` when | Set `False` when |
|-----------|----------------|-----------------|
| `DeploymentReady` | Deployment has desired replicas available | Deployment not ready or reconciliation failed |
| `ServiceReady` | Service applied (or no endpoints → not needed) | Service reconciliation failed |
| `IngressReady` | IngressRoute applied (or no HTTP/HTTPS endpoints) | IngressRoute reconciliation failed |

`ConditionReady` is defined as a constant but is not currently set by the reconciler.

---

## Resource Naming and Labels

All child resources carry these labels:

| Label | Value |
|-------|-------|
| `app.kubernetes.io/name` | Instance name |
| `app.kubernetes.io/part-of` | `"instance"` |
| `app.kubernetes.io/managed-by` | `"tide-controller"` |
| `app.kubernetes.io/component` | `"deployment"` / `"service"` / `"ingress"` |
| `challenges.isolet.dev/id` | Instance name |
| `challenges.isolet.dev/challenge` | Challenge slug |
| `challenges.isolet.dev/challenge-id` | Challenge numeric ID |
| `challenges.isolet.dev/type` | `"dynamic"` or `"on-demand"` |
| `challenges.isolet.dev/team` | Team ID or `"dynamic"` |

---

## SDK

External services interact with Tide through `sdk/` only. The `sdk.Handler` interface is the only supported way to manage Instance CRs from outside the controller.

```go
type Handler interface {
    CreateInstance(ctx, inst *Instance) error
    GetInstance(ctx, name, namespace string) (*Instance, error)
    ListInstances(ctx, opts ...ListOption) ([]Instance, error)
    UpdateInstance(ctx, inst *Instance) error
    DeleteInstance(ctx, name, namespace string) error
}
```

`sdk.New(k8s client.Client) Handler` constructs a client. Always define new methods in `sdk/interface.go` before implementing in `sdk/client.go`.

---

## Development

```sh
# Run the full test suite
make test

# Lint
make lint

# Auto-fix lint issues
make lint-fix

# Build the manager binary to bin/manager
make build

# Regenerate CRDs, RBAC, DeepCopy after changing api/v1/ types or RBAC annotations
make manifests generate

# Run controller locally against the active kubeconfig cluster
make run

# Run e2e tests with a Kind cluster
make test-e2e
```

> After any change to `api/v1/` types or RBAC annotations, always run `make manifests generate`. Never edit `config/crd/bases/` or `api/v1/zz_generated.deepcopy.go` by hand.

---

## Deploying to a Cluster

```sh
# Install CRDs
make install

# Build and push the controller image
make docker-build docker-push IMG=<registry>/tide:tag

# Deploy the controller
make deploy IMG=<registry>/tide:tag

# Apply sample Instance CRs
kubectl apply -k config/samples/

# Uninstall
kubectl delete -k config/samples/
make undeploy
make uninstall
```

---

## Key Conventions

- **Server-Side Apply** — Deployment and Service use `client.Apply` with `FieldOwner: "tide-controller"`. IngressRoute uses `Create`/`Update` (does not use SSA).
- **Owner references** — all child resources are owned by their `Instance` CR via `controllerutil.SetControllerReference`, enabling automatic garbage collection when the CR is deleted.
- **Finalizer** — `challenges.isolet.dev/instance-finalizer` is added before any reconciliation work, removed only after all cleanup (deletion flow) is complete.
- **Status updates** — use `r.Status().Patch()` (not `Update()`). On conflict, requeue with `StatusUpdateRetryInterval` (3s).
- **Logging** — use `logf.FromContext(ctx)` within the reconciler. No module-level loggers.
- **Constants** — all event reasons, label keys, condition types, and the finalizer name live in `utils/constants.go`. Never hardcode these strings.
- **Pointer helpers** — use `utils.Ptr[T]` instead of taking the address of a temporary variable.

---

## Code Navigation

| Path | Contents |
|------|----------|
| `api/v1/instance_types.go` | Full CRD spec: Instance, InstanceSpec, InstanceStatus, Phase, Protocol, ChallengeType |
| `internal/controller/instance_controller.go` | Main reconciler: lifecycle logic, Deployment/Service/IngressRoute reconciliation, phase determination |
| `internal/webhook/v1/instance_webhook.go` | Mutating (defaulter) and validating webhooks |
| `sdk/interface.go` | `Handler` interface |
| `sdk/client.go` | SDK implementation |
| `utils/constants.go` | Event reasons, label keys, condition types, finalizer name, field owner |
| `utils/helpers.go` | `Ptr[T]` and other shared helpers |
| `cmd/main.go` | Manager setup, scheme registration, webhook/controller wiring |

---

## API Reference

- **API group:** `challenges.isolet.dev`
- **Version:** `v1`
- **Kind:** `Instance` (short name: `inst`)
- **kubectl columns:** `Challenge`, `Team`, `Type`, `Phase`, `ExpiresAt`, `Age`

---

## Related Services

- **Oracle** — calls the Tide SDK to create/delete/update Instance CRs
- **Herald** — watches Instance CRDs in K8s (k8s_source mode) and emits `InstanceExpired` facts when phase transitions to Expired
- **Traefik** — must be installed in the cluster; Tide creates `IngressRoute` resources that Traefik watches
