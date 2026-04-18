# UI — Agent Guide

## Purpose & Identity

The **UI** is the CTF platform frontend. It is a Next.js 14 App Router application with TypeScript, Tailwind CSS, and Radix UI. It communicates with Oracle via REST and with Pulse via WebSocket for real-time updates.

---

## Pages

| Route | Description |
|-------|-------------|
| `/` | Home — event countdown, sponsor logos |
| `/login` | User login |
| `/onboard/user` | User registration |
| `/onboard/team` | Team creation/join |
| `/challenges` | Challenge browser |
| `/scoreboard` | Live scoreboard |
| `/profile` | User and team profile |
| `/forgot-password` | Password reset request |
| `/reset-password` | Password reset confirmation |

### Next.js API Routes

| Route | Description |
|-------|-------------|
| `/api/auth/token` | Token management (refresh / session) |

---

## Stack

- **Next.js 14** — App Router, server and client components
- **TypeScript** — strict mode
- **Tailwind CSS** — utility-first styling
- **Radix UI** — headless primitives (`@radix-ui/react-*`)
- **shadcn/ui** conventions — components in `components/ui/`, built on Radix primitives
- **Zustand** — client-side state management
- **Recharts** — scoreboard/score graph charts
- **Lucide React** — icon library
- **react-toastify** — toast notifications
- **Phoenix JS** — WebSocket client for Pulse channels

---

## Directory Layout

```
app/                  Next.js App Router pages and API routes
components/
  ui/                 shadcn/ui base components (Button, Dialog, etc.)
  challenges/         Challenge cards, submission forms, hint panels
  instances/          Instance status cards, start/stop controls
  charts/             Scoreboard/score graph charts (Recharts)
  hints/              Hint unlock components
  profile/            User and team profile components
  NavBar.tsx          Top navigation
  NotificationContainer.tsx  Real-time notification toasts
services/             API client wrappers (one file per Oracle resource)
  base.ts             BaseService with response handling and toast notifications
  auth.ts             Login, register, verify, password reset
  challenge.ts        List, get, submit flag, unlock hint
  instance.ts         Start, stop, extend instances
  score.ts            Scoreboard, score graph
  team.ts             Create, join, update team
  event.ts            Event info (timings, visibility)
  profile.ts          User profile data
realtime/
  socket.ts           Phoenix socket initialization (connects to /socket)
  client.ts           Channel management (global, team:{id})
  dispatcher.ts       Routes incoming events to handlers
  handlers/           Per-event-type handlers:
                        instance.{created,updated,deleted} → instanceStore
                        endpoint.{created,updated,deleted} → instanceStore
                        notification.{info,warning,error,success} → toast
  types.ts            Realtime event type definitions
store/                Zustand stores for client-side state
  challenge.ts, event.ts, instance.ts, profile.ts, score.ts
models/               TypeScript types for API data shapes
  instance.ts, endpoint.ts, hint.ts, score.ts, event.ts
hooks/                Custom React hooks (auth/onboarding flows)
  useLogin.ts, useForgotPassword.ts, useResetPassword.ts
  useUserOnboard.ts, useTeamOnboard.ts, useTeamInvite.ts
  useRealtimeClient.ts
styles/               Global styles
```

---

## API Integration

All Oracle REST calls go through `services/`. The base URL is `/api/v1` (proxied by Nginx to Oracle).

```ts
// services/base.ts — response handling pattern
protected static async handleResponse<TData>(
  apiCall: CancelablePromise<Response & { data?: TData }>
): Promise<TData | undefined>
```

Each service class wraps generated API client methods and normalizes responses with toast feedback.

---

## Real-time Integration

Pulse WebSocket connection via Phoenix JS:

```ts
// realtime/socket.ts — connects to ws://{host}/socket
// realtime/client.ts — joins channels

startRealtime(teamId)
// → joins "global" channel (broadcast events)
// → joins "team:{teamId}" channel (team-specific events)
// Both channels listen on the "notification" event and route through dispatcher
```

Event shape from Pulse:
```ts
{
  event: "instance.created" | "instance.updated" | "notification.warning" | ...,
  // rest of Notification fact payload
}
```

Handlers in `realtime/handlers/` update stores or show toasts based on event type.

---

## Development

```sh
cd ui

# Install dependencies
npm install

# Start dev server
npm run dev          # http://localhost:3000

# Type check
npx tsc --noEmit

# Lint
npm run lint

# Build
npm run build
```

For local dev, the proxy is required to forward `/api` and `/socket` to Oracle and Pulse. Use `docker compose up` at the repo root for the full stack. If you run the Next.js dev server directly, you must configure local proxying or Next.js rewrites for `/api` and `/socket`.

---

## Conventions

- **Radix UI only** — use Radix primitives from `components/ui/`. Do not import other component libraries.
- **Icons** — use `lucide-react`. There are no heroicons in this project.
- **Tailwind for all styling** — no inline styles, no CSS Modules, no styled-components.
- **Server vs client components** — prefer server components. Use `"use client"` only when you need browser APIs, hooks, or event handlers.
- **State** — ephemeral UI state lives in component state; shared state lives in Zustand stores in `store/`.
- **API errors** — surface via `react-toastify` using the `BaseService.handleResponse` pattern. Do not throw uncaught errors in service methods.
- **Types** — all API data shapes go in `models/`. Do not define inline types in components.

---

## Adding a New Page

1. Create directory under `app/` with `page.tsx`
2. Add the route to this doc's page table
3. Create service method in `services/` if the page needs API data
4. Add TypeScript types to `models/` if needed
5. Add store slice in `store/` if state needs to be shared

---

## Adding a New Real-time Event

1. Add the event string key (e.g. `"instance.created"`) to `realtime/handlers/` in the appropriate handler file
2. Register the handler function in the `handlers` map exported from `realtime/handlers/index.ts`
3. Add any new type constants to `realtime/types.ts` if needed
4. The dispatcher (`realtime/dispatcher.ts`) routes by `message.event` automatically — no changes needed there
5. Test by emitting the corresponding `Notification` fact from Herald

---

## Related Services

- **Oracle** — REST API at `/api/v1`
- **Pulse** — WebSocket at `/socket`
- **Proxy** — routes all requests; must be running for API and WebSocket to work
