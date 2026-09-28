# User Management System

A distributed user management system built in Go, following microservice and
event-driven architecture principles. Built as a training assignment covering
REST/WebSocket API design, NATS-based inter-service communication, SQLC-driven
database access, transparent client-side caching, and full CI/CD + Docker
delivery.

## Architecture

The system consists of exactly two services plus a shared client library:

- **Gateway Service** — exposes REST and WebSocket APIs to external clients.
  Contains no direct database access; all user operations go through the User
  Client library.
- **User Service** — owns the `users` PostgreSQL table, exposes CRUD via NATS
  RPC, and publishes user lifecycle events (create/update/delete).
- **User Client library** (`user-service/pkg/user-client`) — a reusable Go
  package living inside the User Service module. Wraps NATS RPC calls and
  event subscriptions behind a plain Go interface, and owns an internal cache
  that is fully transparent to consumers (Gateway included).

```
user-management-system/
├── gateway-service/
│   ├── cmd/gateway-service/   # entrypoint
│   ├── internal/              # gateway-only code, not importable elsewhere
│   └── go.mod
├── user-service/
│   ├── cmd/user-service/      # entrypoint
│   ├── internal/              # user-service-only code (db, rpc handlers)
│   ├── pkg/user-client/       # reusable client library (importable)
│   └── go.mod
├── go.work                    # local multi-module workspace
├── .golangci.yml
├── lefthook.yml
└── docker-compose.yml
```

Each service is its own Go module. Locally, `go.work` lets them resolve each
other's imports (e.g. Gateway importing the User Client library) without
needing a published/tagged dependency.

## Tech Stack

| Concern         | Choice                          |
|------------------|----------------------------------|
| Language         | Go                               |
| Database         | PostgreSQL                       |
| SQL generation   | SQLC                             |
| Messaging        | NATS (RPC + pub/sub)             |
| REST router      | Chi                              |
| WebSocket        | Gorilla WebSocket                |
| Validation       | go-playground/validator          |
| Logging          | `log/slog` (structured, JSON)    |
| Linting          | golangci-lint                    |
| API docs         | OpenAPI (REST), AsyncAPI (WS)    |
| Deployment       | Docker + Docker Compose          |

## Getting Started

### Prerequisites
- Go (see `go.work` for version)
- Docker + Docker Compose
- [lefthook](https://github.com/evilmartians/lefthook)
- [golangci-lint](https://golangci-lint.run/)

### Setup

```bash
git clone <repo-url>
cd user-management-system
lefthook install
```

### Building

```bash
go build ./gateway-service/... ./user-service/...
```

### Running a service locally

```bash
cd gateway-service && go run ./cmd/gateway-service
cd user-service && go run ./cmd/user-service
```

(Full multi-service startup via Docker Compose is documented once available —
see Milestones below.)

### Linting

```bash
golangci-lint run ./gateway-service/...
golangci-lint run ./user-service/...
```

## Development Workflow

- **Commit style:** [Conventional Commits](https://www.conventionalcommits.org/)
  (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`, ...), enforced via
  a `commit-msg` git hook.
- **Pre-commit checks:** golangci-lint runs automatically before each commit
  via [lefthook](https://github.com/evilmartians/lefthook). Hooks are wired up
  by `lefthook install` (see Setup above); run them manually at any time with:
  ```bash
  lefthook run pre-commit
  ```
- **Branching:** feature branches + pull requests for each milestone, with a
  short design note in the PR description.

## Internal NATS Design

_To be documented here once the RPC/event subject conventions are finalized
(Milestone 2–3)._

## Milestones

- [x] **M0** — Repo, tooling, hooks, CI skeleton, base Compose (Postgres + NATS)
- [ ] **M1** — Data model & SQLC foundation
- [ ] **M2** — NATS RPC + client library skeleton
- [ ] **M3** — User events + subscribe/unsubscribe
- [ ] **M4** — Client library caching (transparent)
- [ ] **M5** — Gateway REST API + OpenAPI draft
- [ ] **M6** — Gateway WebSocket API + AsyncAPI draft
- [ ] **M7** — Containerization + full system Compose
- [ ] **M8** — Hardening & release CI complete

## API Documentation

- REST (OpenAPI): `/docs/openapi.yaml` _(added in M5)_
- WebSocket (AsyncAPI): `/docs/asyncapi.yaml` _(added in M6)_