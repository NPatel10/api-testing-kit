# API Testing Kit - Implementation Tasks

This file is the execution plan for building the project in parallel.

The goal is to let multiple subagents work at the same time with minimal overlap:
- one track owns the app shell and design system
- one track owns the guest-authenticated workspace
- one track owns backend APIs and persistence
- one track owns security, rate limiting, and request execution
- one track owns public pages and product docs
- one optional later track owns monetization and entitlement groundwork
- one track owns deployment and test infrastructure

All tasks below are aligned with `plan.md`, `ui-pages.md`, and `design-guideline.md`.

---

## Verification Notes

- Verification on `2026-03-26` found follow-up gaps after the initial completion pass.
- A task marked `completed` can still have a later verification follow-up if the underlying feature exists but is not yet wired end-to-end or fully covered by tests/CI.

---

## Work Rules

- Do not start feature work before the shared scaffold exists.
- Do not let multiple agents edit the same file unless the task explicitly says so.
- Treat the database schema, route map, and design tokens as shared contracts.
- Keep guest mode on `/app` limited by backend rules, not only UI state.
- Prefer small, independently mergeable changes.

### Assumed Post-Scaffold Ownership Slices

Adapt these to the scaffold you create, but keep ownership boundaries similar:

- `apps/web`: SvelteKit routes, layouts, load functions, and marketing pages
- `server/cmd/api` and `server/internal`: Go backend modules, auth, runner, rate limiting, and admin endpoints
- `apps/web/src/lib` or a future `packages/ui`: shared app shell, themed components, and UI primitives
- `schema.sql`, `goose` migrations, and any future Go data layer: schema, migrations, seeders, query helpers
- `infra` or deployment folder: Docker, Compose, CI, deployment scripts

---

## Blockers First

These items should be completed before most parallel feature work starts.

### BL1. Project Scaffold
Owner: `Agent Group A`
Status: `completed` on `2026-03-26`

Scope:
- Create the base SvelteKit app structure.
- Wire Tailwind CSS and any chosen Svelte-compatible component primitive layer.
- Add core folder conventions for `lib`, `routes`, `components`, `server`, and `styles`.
- Add environment variable scaffolding and example `.env` shape.
- Add basic scripts for dev, build, lint, typecheck, and test.

Depends on:
- none

Unblocks:
- all frontend tracks
- all route-level implementation work

### BL2. Database Schema
Owner: `Agent Group C`
Status: `completed` on `2026-03-26`

Scope:
- Create `schema.sql`.
- Choose `goose` as the migration tool and convert the draft schema into an initial migration set when DB work begins.
- Define core entities for users, sessions, collections, saved requests, templates, history, usage, abuse, blocked targets, and billing/entitlements.
- Add indexes, constraints, and comments.

Depends on:
- none

Unblocks:
- data layer integration
- backend API work
- auth work
- collections/history work

### BL3. Route and App Contract
Owner: `Agent Group A`
Status: `completed` on `2026-03-26`

Scope:
- Confirm the `/app` shared guest/authenticated model.
- Define route-level guards and locked states.
- Establish common app shell structure and route layout ownership.

Depends on:
- BL1
- `ui-pages.md`

Unblocks:
- guest gating work
- shell/layout work
- workspace pages

---

## Track A. Frontend Shell And Design System

### A1. Theme Tokens And Global Styles
Owner: `Agent Group A`
Status: `completed` on `2026-03-25`

Scope:
- Translate `design-guideline.md` into Tailwind theme tokens.
- Add global CSS variables for shell, canvas, panel, success, warning, danger, and shadows.
- Configure typography defaults and code font pairing.

Depends on:
- BL1
- `design-guideline.md`

Parallel With:
- backend tracks
- auth tracks
- route content tracks

### A2. shadcn Component Styling Pass
Owner: `Agent Group A`
Status: `completed` on `2026-03-25`

Scope:
- Re-theme button, card, badge, input, tabs, dialog, sheet, dropdown, tooltip, separator, scroll area, and table primitives.
- Standardize radius, border, shadow, and color behavior.

Depends on:
- A1

Parallel With:
- workspace feature work
- marketing pages

### A3. App Shell Layout
Owner: `Agent Group A`
Status: `completed` on `2026-03-25`

Scope:
- Build the shared shell for `/app`.
- Add sidebar, top toolbar, central workspace canvas, and utility drawer region.
- Support responsive collapse for smaller screens.

Depends on:
- A1
- BL3

Unblocks:
- request builder
- response viewer
- history
- collections

### A4. Marketing Layout System
Owner: `Agent Group A`
Status: `completed` on `2026-03-25`

Scope:
- Build the landing page layout system that matches the app visual language.
- Add reusable hero, feature grid, CTA, and content section patterns.

Depends on:
- A1

Parallel With:
- `/app` workspace implementation
- backend APIs

---

## Track B. Guest Gating And Workspace UX

### UX1. Guest Mode State Model
Owner: `Agent Group B`
Status: `completed` on `2026-03-25`

Scope:
- Define guest vs authenticated UI state on `/app`.
- Create shared access rules for locked actions and upgrade prompts.
- Map visible but disabled controls for guests.

Depends on:
- BL3

Parallel With:
- backend auth
- collections/templates

### UX2. Request Builder UI
Owner: `Agent Group B`
Status: `completed` on `2026-03-25`; live runner wiring follow-up opened on `2026-03-26`

Scope:
- Implement method selector, URL bar, params, headers, auth, body, and send action.
- Support JSON, raw text, and form-urlencoded body modes.
- Add request validation feedback and empty states.

Depends on:
- A3
- UX1

Unblocks:
- response viewer
- history
- collections save flows

### UX3. Response Viewer UI
Owner: `Agent Group B`
Status: `completed` on `2026-03-25`; live response hydration follow-up opened on `2026-03-26`

Scope:
- Implement pretty/raw/headers tabs.
- Show status, duration, size, content type, and error states.
- Add code block styling for response data.

Depends on:
- A3
- UX2

### UX4. Guest Lock Surfaces
Owner: `Agent Group B`
Status: `completed` on `2026-03-25`

Scope:
- Add locked overlays for save, history persistence, env vars, custom URLs, and advanced tools.
- Make upgrade/sign-in CTAs visible but non-intrusive.

Depends on:
- UX1
- A3

### UX5. Templates And Example Collections UI
Owner: `Agent Group B`
Status: `completed` on `2026-03-25`

Scope:
- Build the templates browser and example request launch points.
- Make guest-safe templates available inside `/app`.

Depends on:
- UX1
- backend template endpoints

### UX6. Live `/app` Request Execution Wiring
Owner: `Agent Group B`
Status: `completed` on `2026-03-26`

Scope:
- Replace default mock request/response previews in `/app` with live execution flows where backend contracts already exist.
- Wire guest sends to `/api/v1/guest-runs`.
- Wire authenticated sends to `/api/v1/runs`.
- Hydrate the response viewer from actual execution results, not only seeded examples.
- Keep preview/mock state only as an explicit fallback or demo mode, not the default authenticated workspace path.

Depends on:
- UX2
- UX3
- C2
- C5
- D1
- D3

Unblocks:
- end-to-end `/app` verification
- real guest/authenticated request execution checks

### UX7. Live History And Collection Surfaces
Owner: `Agent Group B`
Status: `completed` on `2026-03-26`

Scope:
- Replace mock-backed `/app/history` and `/app/collections/[id]` loaders with live backend data.
- Load authenticated collections, saved requests, and history from API responses.
- Preserve guest locked/preview states without implying persistence that does not exist.
- Keep route-level UX aligned with the shared `/app` contract while using real persistence for signed-in sessions.

Depends on:
- C3
- C4
- C7
- C8

Unblocks:
- persisted workspace verification
- authenticated reuse and rerun flows

---

## Track C. Backend API And Persistence

### C1. DB Layer And Migrations
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`

Scope:
- Connect the app to PostgreSQL using the chosen Go data layer.
- Create a `goose` migration workflow.
- Add typed repository helpers for core tables.

Depends on:
- BL2

Parallel With:
- frontend shell
- auth

### C2. Auth And Sessions
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`

Scope:
- Implement signup, login, logout, and session handling.
- Add account verification and basic account metadata.
- Expose session identity to the frontend.

Depends on:
- BL2
- C1

### C3. Collections API
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`; collection detail surface follow-up completed on `2026-03-26`; live `/app` consumer follow-up completed on `2026-03-26`

Scope:
- CRUD for collections.
- Nested request grouping if supported by the schema.
- Ownership checks and read-only access rules.

Depends on:
- C1
- C2

### C4. Saved Requests And History API
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`; authenticated history surface follow-up completed on `2026-03-26`; live `/app` consumer follow-up completed on `2026-03-26`

Scope:
- CRUD for saved requests.
- Request run history persistence.
- Store execution snapshots and summary metadata.

Depends on:
- C1
- C2

### C5. Templates API
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`

Scope:
- Serve guest-safe template collections.
- Return example payloads and metadata for `/app`.

Depends on:
- C1

### C6. Usage And Abuse Storage
Owner: `Agent Group C`
Status: `completed` on `2026-03-25`; admin monitoring query surface follow-up completed on `2026-03-26`

Scope:
- Persist usage events, blocked attempts, abuse events, and denylist entries.
- Expose admin-friendly query shapes for monitoring.

Depends on:
- BL2
- C1

### C7. Collection Detail Query Surface
Owner: `Agent Group C`
Status: `completed` on `2026-03-26`

Scope:
- Add a backend query surface for `/app/collections/[id]`.
- Return collection metadata plus the saved requests that belong to the collection in stable display order.
- Preserve ownership checks and any read-only sharing rules already defined by the data model.
- Avoid forcing the frontend to reconstruct collection detail from seeded mocks.

Depends on:
- C3
- C4
- C2

Unblocks:
- UX7
- authenticated collection detail routing

### C8. History Filters And Re-run Query Surface
Owner: `Agent Group C`
Status: `completed` on `2026-03-26`

Scope:
- Extend history APIs with filters for status, domain, method, and date.
- Return enough request snapshot detail to support rerun and save-to-collection flows from history.
- Make pagination and limit behavior explicit for authenticated history views.
- Keep access checks and ownership boundaries unchanged while broadening the query contract.

Depends on:
- C4
- C2
- D1

Unblocks:
- UX7
- history filters and reuse flows

---

## Track D. Outbound Runner, SSRF, And Rate Limiting

### D1. Request Runner Core
Owner: `Agent Group D`
Status: `completed` on `2026-03-25`

Scope:
- Build the server-side request execution pipeline.
- Support method, URL, headers, query params, auth, and body forwarding.
- Return structured metadata plus truncated/previewed response bodies.

Depends on:
- C1
- UX2

### D2. Destination Validation
Owner: `Agent Group D`
Status: `completed` on `2026-03-25`

Scope:
- Block localhost, private ranges, metadata IPs, and unsupported protocols.
- Validate DNS resolution and redirect hops.
- Enforce port and redirect restrictions.

Depends on:
- D1

### D3. Guest Restrictions
Owner: `Agent Group D`
Status: `completed` on `2026-03-26`

Scope:
- Enforce allowlisted endpoints for guests on `/app`.
- Deny arbitrary outbound requests from guest sessions.
- Enforce guest size and timeout limits.

Depends on:
- C2
- D1

### D4. Rate Limiting Engine
Owner: `Agent Group D`
Status: `completed` on `2026-03-25` for limiter primitives and guest wiring; authenticated runner follow-up completed on `2026-03-26`

Scope:
- Implement IP and user quota enforcement.
- Add burst protection, cooldowns, and redirect limits.
- Add domain-level throttling hooks.

Depends on:
- C6
- C2

### D5. Abuse Detection Hooks
Owner: `Agent Group D`
Status: `completed` on `2026-03-25`; admin abuse monitor surface follow-up completed on `2026-03-26`

Scope:
- Log suspicious request patterns.
- Emit blocked event records for admin review.
- Create reusable checks for scanning/spam-relay behavior.

Depends on:
- D1
- C6

### D6. Authenticated Runner Limit Enforcement
Owner: `Agent Group D`
Status: `completed` on `2026-03-26`

Scope:
- Apply IP, user, and domain limiter checks to authenticated `/api/v1/runs`.
- Enforce authenticated size, timeout, concurrency, and redirect limits through the live runner path, not only guest mode or UI copy.
- Record quota denials, cooldown hits, and suspicious spikes in usage/abuse storage for admin review.
- Keep the limiter behavior consistent with the rate-limit rules documented in `ui-pages.md`.

Depends on:
- D1
- D4
- C6
- C2

Unblocks:
- end-to-end authenticated runner verification
- accurate abuse and quota monitoring

---

## Track E. Optional Monetization And Entitlement Groundwork

This track is intentionally optional.

Do not start it until the guest/authenticated request runner, collections/history, and core safety controls are stable.

If monetization is not part of the current milestone, skip this track entirely.

### E1. Plan And Entitlement Model
Owner: `Agent Group E`
Status: `completed` on `2026-03-26`

Scope:
- Define plan tiers and feature entitlements.
- Add backend data structures for limits, trial state, and upgrades.
- Keep the model provider-agnostic for Stripe, Paddle, or Lemon Squeezy.

Depends on:
- BL2

### E2. Billing Customer And Subscription Storage
Owner: `Agent Group E`
Status: `completed` on `2026-03-26`

Scope:
- Store billing customer IDs, subscription status, renewal state, and plan mapping.
- Add invoice and payment event persistence.

Depends on:
- BL2
- E1

### E3. Access Control By Plan
Owner: `Agent Group E`
Status: `completed` on `2026-03-26`

Scope:
- Gate custom URL execution, saved history depth, env vars, and shared links by entitlement.
- Feed plan state into backend authorization and frontend locked states.

Depends on:
- E1
- C2
- C3

### E4. Billing Provider Integration Stub
Owner: `Agent Group E`
Status: `completed` on `2026-03-26`

Scope:
- Add webhook endpoint skeleton.
- Add checkout success/cancel flow hooks.
- Keep implementation provider-agnostic until provider choice is finalized.

Depends on:
- C2
- E2

---

## Track F. Public Pages And Product Content

### F1. Landing Page
Owner: `Agent Group F`
Status: `completed` on `2026-03-25`

Scope:
- Build the marketing landing page.
- Match the design system and drive users into `/app`.

Depends on:
- A4

### F2. Templates Page
Owner: `Agent Group F`
Status: `completed` on `2026-03-25`

Scope:
- Build template browsing and category filters.
- Connect to backend template data.

Depends on:
- C5
- A4

### F3. Features Page
Owner: `Agent Group F`
Status: `completed` on `2026-03-25`

Scope:
- Document product capabilities, guest limits, and authenticated unlocks.

Depends on:
- A4

### F4. Docs Page
Owner: `Agent Group F`
Status: `completed` on `2026-03-25`

Scope:
- Add quick start content and guest/auth usage guidance.

Depends on:
- A4

### F5. Case Study Page
Owner: `Agent Group F`
Status: `completed` on `2026-03-25`

Scope:
- Add architecture and engineering narrative.

Depends on:
- A4
- backend architecture decisions

---

## Track G. Deployment, DevOps, And Test Infrastructure

### G1. Local Dev Environment
Owner: `Agent Group G`
Status: `completed` on `2026-03-25`

Scope:
- Docker and Docker Compose setup.
- Local PostgreSQL bootstrap.
- `goose` migration execution in local setup.
- Optional Redis bootstrap if used for rate limiting.

Depends on:
- BL1
- BL2

### G2. Application Deployment
Owner: `Agent Group G`
Status: `completed` on `2026-03-25`

Scope:
- Add production build and runtime instructions.
- Prepare environment configuration for deployment.
- Define service boundaries and startup order.

Depends on:
- G1
- C1
- D1

### G3. Test Harness
Owner: `Agent Group G`
Status: `completed` on `2026-03-25`; browser verification follow-up completed on `2026-03-26`

Scope:
- Add unit and integration test setup.
- Add request runner tests, auth tests, and basic UI smoke tests.

Depends on:
- BL1
- C1
- D1

### G4. CI Checks
Owner: `Agent Group G`
Status: `completed` on `2026-03-26`; browser coverage follow-up completed on `2026-03-26`

Scope:
- Add lint, typecheck, test, and build checks.
- Keep the checks fast enough for parallel feature work.

Depends on:
- G3

### G5. Browser Verification And CI Coverage
Owner: `Agent Group G`
Status: `completed` on `2026-03-26`

Scope:
- Repair stale Playwright smoke and e2e assertions so browser suites pass again.
- Make the docs, `/app`, and related public routes under test align with the current UI copy and structure.
- Run the agreed browser checks in CI, or explicitly split them into required and visible non-blocking jobs if runtime is a concern.
- Keep the verification layer strong enough to catch drift between docs, route content, and shipped UI.

Depends on:
- G3
- G4
- F4

Unblocks:
- trustworthy browser regression coverage
- CI parity with local verification

---

## Parallel Execution Map

These tracks can run at the same time after the blockers are in place:

- `Agent Group A`: design system, app shell, marketing shell
- `Agent Group B`: guest gating, request builder, response viewer, templates UI
- `Agent Group C`: DB layer, auth, collections, history, templates API, abuse storage
- `Agent Group D`: runner, SSRF, guest restrictions, rate limiting, abuse hooks
- `Agent Group E`: optional monetization, entitlement storage, billing scaffolding
- `Agent Group F`: landing page, docs, features, case study
- `Agent Group G`: Docker, CI, deployment, test harness

The only hard blockers are:
- BL1 scaffold
- BL2 schema
- BL3 route/app contract

After those are in place, the rest can proceed in parallel with careful file ownership.

---

## Suggested Merge Order

1. Merge scaffold, schema, and route contract.
2. Merge design system and shell.
3. Merge backend base APIs and DB wiring.
4. Merge guest workspace UI and response viewer.
5. Merge request runner, SSRF, and rate limits.
6. Merge collections/history/templates.
7. Merge marketing pages and docs.
8. Merge deployment and test infrastructure.
9. Merge optional monetization groundwork, only if it is in scope.

---

## Definition Of Done

The implementation phase is done when:
- `/app` works in guest and authenticated modes with live backend-backed execution, not only seeded previews
- requests can be sent safely from the backend
- collections and history are persisted and surfaced from live APIs
- abuse protection and rate limits are active on both guest and authenticated runner paths
- the UI matches the approved design system
- the app can be built and deployed locally with Docker
- browser smoke/e2e checks are green for the agreed suite, and CI runs the required verification set
- the task list supports multiple subagents without overlapping ownership
