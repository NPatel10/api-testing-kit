# API Testing Kit

API Testing Kit is a monorepo for a developer-focused API testing product.

## Current Repository State

- `apps/web` is the SvelteKit frontend scaffold and still contains mostly starter routes.
- `server` is the Go HTTP scaffold. Only health endpoints are implemented today.
- `schema.sql` is a future-state PostgreSQL schema draft and is ahead of the current backend implementation.
- `goose` is the chosen database migration tool once persistence work starts. The current `schema.sql` should be converted into the initial migration set.
- `doc/` contains the target product, UX, and execution docs. Those docs should explicitly separate implemented behavior from planned behavior.

## Repository Structure

```text
api-testing-kit/
  apps/
    web/
  server/
  doc/
  schema.sql
  package.json
  README.md
```

## Packages

### Frontend
- Path: `apps/web`
- Stack: SvelteKit, TypeScript, Tailwind CSS
- Current status: starter scaffold that still needs to be aligned to the `/app` product direction

### Backend
- Path: `server`
- Stack: Go
- Current status: minimal HTTP scaffold with health endpoints

### Documentation
- Path: `doc`
- Includes:
  - `basic.md` for the project summary and locked decisions
  - `plan.md` for the target product direction and phased scope
  - `ui-pages.md` for route and workspace planning
  - `design-guideline.md` for the visual system
  - `tasks.md` for implementation sequencing
  - `test-suite-proposal.md` for test planning
  - `apis.md` for implemented endpoints and planned API contracts

## Root Commands

### Monorepo

```bash
bun run dev
bun run build
bun run preview
bun run lint
bun run typecheck
bun run test
```

### Frontend

```bash
bun run dev:web
bun run build:web
bun run preview:web
bun run check:web
bun run lint:web
bun run typecheck:web
bun run test:web
```

### Backend

```bash
bun run dev:server
bun run lint:server
bun run typecheck:server
bun run test:server
```

## Environment

Copy values from `.env.example` before local development.

Current variables:

```bash
INTERNAL_API_BASE_URL=http://api:8080
PUBLIC_API_BASE_URL=http://localhost:8080
API_PORT=8080
WEB_PORT=4173
POSTGRES_PORT=5432
POSTGRES_DB=api_testing_kit
POSTGRES_USER=api_testing_kit
POSTGRES_PASSWORD=api_testing_kit_dev
DATABASE_URL=postgres://api_testing_kit:api_testing_kit_dev@db:5432/api_testing_kit?sslmode=disable
```

## Local Docker Stack

Start the frontend, Go API, and PostgreSQL with:

```bash
docker compose up --build
```

The stack exposes:

- `http://localhost:4173` for the web app
- `http://localhost:8080` for the Go API
- `localhost:5432` for PostgreSQL

The web container uses `INTERNAL_API_BASE_URL` for server-side fetches so SvelteKit can reach the API container directly, while `PUBLIC_API_BASE_URL` remains host-friendly for browser links and local non-Docker development.

PostgreSQL bootstrap SQL lives in `docker/postgres/init/001-bootstrap.sql`.
Goose migrations should be added under `server/migrations`, then run with:

```bash
docker compose --profile migrations run --rm migrate
```

## Deployment

The deployment order and service boundaries are documented in [`doc/deployment.md`](./doc/deployment.md).

In short:

1. Start PostgreSQL.
2. Run `migrate` once against the target database.
3. Start the API after the database is healthy.
4. Start the web app after the API is healthy.

The compose stack now enforces that order with health checks for `db`, `api`, and `web`.

## Backend Endpoints

The current Go scaffold exposes:

- `/`
- `/healthz`
- `/api/v1/health`

## Source Of Truth

- Current implementation state: this README plus the code under `apps/web`, `server`, and `schema.sql`
- Target product and UX direction: `doc/plan.md`, `doc/ui-pages.md`, and `doc/design-guideline.md`
- Execution backlog: `doc/tasks.md`
- Backend contract decisions: `doc/apis.md`
