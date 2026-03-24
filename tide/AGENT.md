# Tide — Agent Guidelines

This is a Kubernetes operator written in Go using the [Kubebuilder](https://book.kubebuilder.io/) framework and [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime).

## Project overview

Tide is the infrastructure execution layer of Isolet. It manages **Challenge Instances** — isolated, ephemeral Kubernetes environments for CTF challenges. For every `Instance` custom resource (CR), Tide creates and reconciles:

- A **Deployment** running the challenge container (`deployment-{uuid}`)
- A **Service** exposing it inside the cluster (`svc-{uuid}`)
- A **Traefik IngressRoute** for HTTP/HTTPS routing (`ir-{uuid}`)

Tide is stateless — it works off CR state only and carries no business logic. Upstream services (Oracle, Herald) create Instance CRs; Tide executes them.

## Project guidelines

- Run `make test` to run the full test suite (includes `manifests`, `generate`, `fmt`, `vet`, and `envtest`).
- Run `make lint` to run the linter (`golangci-lint`). Run `make lint-fix` to auto-fix linting issues.
- Run `make build` to build the manager binary to `bin/manager`.
- After any change to CRD types in `api/v1/` or RBAC annotations, run `make manifests generate` to regenerate CRDs, RBAC roles, and DeepCopy methods. **Never** edit files under `config/crd/bases/` or `api/v1/zz_generated.deepcopy.go` by hand.
- Use `make run` to run the controller locally against the active kubeconfig cluster.
- Use `make test-e2e` to run end-to-end tests with a Kind cluster.

## Go guidelines

- **Always** use the `utils.Ptr[T]` generic helper for pointer values of any type instead of taking the address of a temporary variable.
- All constants (event reasons, label keys, condition types, finalizer name) live in `utils/constants.go`. **Never** hardcode these strings elsewhere.
- All Kubernetes events must use the constants defined in `utils/constants.go` (e.g., `utils.EventReasonCreated`, `utils.EventReasonInstanceRunning`).
- Use `logf.FromContext(ctx)` for contextual logging within the reconciler. Do not create module-level loggers.
- Follow the existing kubebuilder RBAC annotation pattern (`// +kubebuilder:rbac:...`) above the `Reconcile` method. After adding or changing RBAC annotations, run `make manifests` to regenerate `ClusterRole` manifests.
- The reconciler uses **Server-Side Apply** via `client.Apply` with `FieldOwner` set to `utils.FieldOwner` (`"tide-controller"`). Always use this pattern when creating or patching child resources — do not use `Create` followed by `Update`.
- Use `controllerutil.SetControllerReference` to set the Instance as the owner of all child resources (Deployment, Service, IngressRoute). This enables automatic garbage collection.
- **Always** guard against missing finalizers: add `utils.InstanceFinalizer` before any other reconciliation work, and remove it only after all child resources have been cleaned up.
- Status updates must go through the status subresource client (`r.Status().Update()`), not through the main client. Conflicts on status update are expected; retry with the `StatusUpdateRetryInterval` constant.

## API & CRD guidelines

- The API group is `challenges.isolet.dev`, version `v1`, and the CRD kind is `Instance` (short name: `inst`).
- The `Phase` field is the primary signal of instance lifecycle. Valid transitions are: `Pending` → `Staged` | `Running` | `Failed`, `Staged` → `Running` | `Failed`, `Running` → `Expired` | `Terminated`. Never skip a phase or transition backwards.
- `ChallengeType` is either `"dynamic"` (shared, long-running) or `"on-demand"` (isolated, per-team). Controller behavior (especially team label and hostname generation) differs between the two.
- `Protocol` values (`http`, `https`, `nc`, `ssh`) drive how IngressRoutes are generated. `http`/`https` endpoints get Traefik IngressRoute rules; `nc`/`ssh` endpoints use TCP routing.
- Hostname format: single endpoint → `{instance-uuid}.{challenge-slug}.{domain}`; multiple endpoints → `{endpoint-name}-{instance-uuid}.{challenge-slug}.{domain}`.
- The `Lifecycle` struct defaults (`AllowExtension = true`) are applied by the mutating webhook. Do not duplicate default logic in the controller.

## SDK guidelines

- External services (Oracle, Herald) interact with Tide through the `sdk/` package, which implements the `sdk.Handler` interface.
- The `sdk.Handler` interface methods (`CreateInstance`, `GetInstance`, `ListInstances`, `UpdateInstance`, `DeleteInstance`) are the **only** supported way to manage Instance CRs from outside the controller. Do not interact with the Kubernetes API directly from consuming services.
- When adding new operations to the SDK, always define the method signature in `sdk/interface.go` first, then implement it in `sdk/client.go`.

## Tide project guidelines

- **Isolation**: Every Instance must run in its own Deployment with its own Service and IngressRoute. Never share Kubernetes resources across Instances.
- **Expiry**: The controller schedules its own requeueing based on `lifecycle.expiresAt` and `lifecycle.availableAt`. Do not rely on external cron jobs or TTL controllers for lifecycle management.
- **Failure detection**: The controller actively polls Pod conditions for `ImagePullBackOff`, `CrashLoopBackOff`, `OOMKilled`, and similar failure reasons to set `PhaseFailed`. Maintain this proactive detection — do not rely solely on Deployment `.status.availableReplicas`.
- **No business logic**: Tide must remain domain-agnostic. It executes what Instance CRs specify. Business rules (who can create an instance, time limits, scoring) belong in Oracle or Herald, not here.
- **Webhook validation**: The validating webhook (`internal/webhook/v1/`) is the enforcement point for Instance CR constraints. Keep validation logic in the webhook, not scattered in the reconciler.
- **Conditions**: Always keep the four standard conditions (`DeploymentReady`, `ServiceReady`, `IngressReady`, `Ready`) up to date in `instance.Status.Conditions`. Use `meta.SetStatusCondition` for idempotent updates.
