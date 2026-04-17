# Proxy — Agent Guide

## Purpose

The **Proxy** is an Nginx reverse proxy that serves as the single entry point to all Isolet services. It routes incoming HTTP and WebSocket traffic, handles trusted IP header forwarding (Cloudflare), and enables zero-config TLS termination when placed behind a CDN or load balancer.

---

## Routing Table

| Path | Upstream | Notes |
|------|----------|-------|
| `/` | UI (`UI_URL`) | Next.js app |
| `/api` | Oracle (`API_URL`) | REST API (`/api/v1/...`) |
| `/socket` | Pulse (`PULSE_URL`) | WebSocket upgrade (Phoenix Channels) |
| `/files` | Phoros (`PHOROS_URL`) | Challenge file downloads |

---

## Configuration

All values are injected at container startup via `envsubst` into `default.conf.template`.

| Variable | Example | Description |
|----------|---------|-------------|
| `UI_URL` | `ui:3000` | UI service host:port |
| `API_URL` | `oracle:8000` | Oracle service host:port |
| `PULSE_URL` | `pulse:4000` | Pulse service host:port |
| `PHOROS_URL` | `phoros:8080` | Phoros file server host:port |
| `PROXY_SERVER_NAME` | `isolet.dev localhost` | Nginx `server_name` directive value |
| `TRUSTED_IPS` | See below | One `set_real_ip_from` directive per line |
| `TRUSTED_HEADER` | `real_ip_header CF-Connecting-IP;` | Which header contains the real client IP |

### Trusted IP configuration (Cloudflare example)

```yaml
TRUSTED_IPS: |
  set_real_ip_from 173.245.48.0/20;
  set_real_ip_from 103.21.244.0/22;
  set_real_ip_from 104.16.0.0/13;
TRUSTED_HEADER: "real_ip_header CF-Connecting-IP;"
```

When not behind a trusted proxy, pass empty strings for both variables.

---

## TLS

The proxy listens on port 80 only. TLS is expected to be terminated upstream — by Cloudflare, a load balancer, or cert-manager with an external ingress controller. Do not add TLS config to Nginx unless you are handling TLS inside the cluster.

---

## WebSocket Handling

The `/socket` location block sets the required upgrade headers:

```nginx
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection 'Upgrade';
```

This is required for Phoenix Channel WebSocket connections to work correctly.

---

## Local Development

The proxy is included in `docker-compose.yaml` and exposes port 80:

```sh
docker compose up proxy
```

In local dev, `PULSE_URL` defaults to `ui:3000` (proxying through the UI dev server) because Pulse is not included in the compose stack. For full WebSocket testing, run Pulse separately and update `PULSE_URL`.

---

## Changing Routes

1. Edit `default.conf.template` — add or modify `location` blocks
2. If the new upstream is a new service, add its URL as an env var
3. Update `docker-compose.yaml` to pass the new env var
4. Update `charts/values.yaml` and the Helm templates for production

---

## Files

| File | Description |
|------|-------------|
| `default.conf.template` | Nginx config template with `${VAR}` placeholders |
| `Dockerfile` | Extends `nginx:alpine`, copies template, runs `envsubst` on start |
| `VERSION` | Image version tag |
