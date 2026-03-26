# API Testing Kit - API Contracts

## Purpose

This file is the backend contract document for the project.

It must always distinguish between:
- implemented endpoints that exist in the current repo
- planned endpoints that describe the target product shape

## Current Implemented Endpoints

These endpoints exist today in the Go scaffold under `server`.

### `GET /`

Returns a basic JSON message:

```json
{
  "message": "API Testing Kit server scaffold is running"
}
```

### `GET /healthz`

Returns:

```json
{
  "status": "ok",
  "service": "api-testing-kit-server",
  "timestamp": "RFC3339 UTC timestamp"
}
```

### `GET /api/v1/health`

Returns the same payload as `/healthz`.

### `GET /api/v1/admin/abuse`

Returns an admin-only abuse monitoring payload with:

- grouped abuse summaries
- recent abuse events
- blocked target records

### `GET /api/v1/admin/blocked-targets`

Returns an admin-only blocked target list.

## Contract Decisions Locked Now

- Backend runtime: Go in `server`, not Bun.
- App-facing JSON endpoints should live under `/api/v1`.
- Browser auth should use server-managed sessions with HTTP-only cookies when auth is introduced.
- Guest execution must never accept arbitrary raw target URLs from the browser.
- Guest execution should use allowlisted templates plus safe field overrides only.
- Authenticated execution may accept a full request definition, but the server must still revalidate destination, DNS resolution, redirects, ports, body size, response size, and timeouts.
- Billing and monetization routes are later-phase and should not block the MVP API design.

## Planned Route Groups

The routes below are target-state planning only. They are not implemented yet.

### Public / Guest

- `GET /api/v1/templates`
  - List guest-safe template collections and requests.
- `GET /api/v1/templates/:slug`
  - Return a single template collection or request.
- `POST /api/v1/guest-runs`
  - Execute an allowlisted guest template.
  - Request body should reference a template ID or slug plus explicitly allowlisted overrides.
  - Request body must not contain a raw arbitrary target URL.

Suggested guest run payload shape:

```json
{
  "templateId": "weather-demo",
  "overrides": {
    "query": {
      "city": "Kolkata"
    }
  }
}
```

### Authenticated

- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/me`
- `POST /api/v1/runs`
- `GET /api/v1/history`
- `GET /api/v1/collections`
- `POST /api/v1/collections`
- `PATCH /api/v1/collections/:id`
- `DELETE /api/v1/collections/:id`
- `GET /api/v1/requests/:id`
- `POST /api/v1/requests`
- `PATCH /api/v1/requests/:id`
- `DELETE /api/v1/requests/:id`
- `GET /api/v1/settings`
- `PUT /api/v1/settings`

Suggested authenticated run payload shape:

```json
{
  "method": "GET",
  "url": "https://api.example.com/users",
  "queryParams": [],
  "headers": [],
  "auth": {
    "scheme": "none"
  },
  "body": {
    "mode": "none"
  }
}
```

### Admin / Internal

- `GET /api/v1/admin/abuse`
- `GET /api/v1/admin/blocked-targets`
- `POST /api/v1/admin/blocked-targets`
- `PATCH /api/v1/admin/blocked-targets/:id`

## Response Conventions

- JSON responses by default
- RFC3339 timestamps
- Stable top-level error shape once error handling is introduced
- Clear distinction between:
  - validation failure
  - blocked/safety rejection
  - upstream request failure
  - quota/rate-limit rejection

Suggested error shape:

```json
{
  "error": {
    "code": "blocked_target",
    "message": "Requests to private network targets are not allowed."
  }
}
```

## Important Non-Goals For The Contract

- Do not design guest APIs as a public generic proxy.
- Do not couple the initial contract to billing or paid plans.
- Do not assume GraphQL, WebSocket, or team-workspace routes in V1.
