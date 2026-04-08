# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sonic is a high-performance blogging platform written in Go, using Cloudwego Hertz as its HTTP framework, GORM for database access, and Uber FX for dependency injection.

## Common Commands

```bash
# Build
go build ./...

# Run
go run main.go

# Run tests
go test ./...
go test ./path/to/package/... -run TestName  # single test

# Lint
golangci-lint run --config=.golangci.yml
```

## Architecture

### Request Flow

```
HTTP Request
→ Hertz Router (handler/router.go)
→ Global Middleware (timeout, locale, request ID)
→ Route-specific Middleware (auth, CSRF, rate limit)
→ Handler (admin/ or content/)
→ Service (service/impl/)
→ DAL (dal/) via GORM
→ Response
```

### Key Layers

- **`handler/`** — HTTP handlers and routing
  - `handler/admin/` — Admin console API (authenticated endpoints)
  - `handler/content/` — Public-facing content API
  - `handler/middleware/` — Auth, CSRF, rate limiting, timeout, logging
  - `handler/web/` — Hertz adapter/abstraction (request/response helpers)
  - `handler/binding/` — Request parsing and validation
- **`service/impl/`** — Business logic implementations
- **`dal/`** — Data access layer; GORM query generation (generated code lives here)
- **`model/`** — Data structures: `entity/` (GORM models), `dto/` (request/response), `vo/` (value objects), `param/`, `projection/`
- **`injection/`** — Uber FX dependency injection wiring
- **`template/extension/`** — Custom Go template functions for content rendering
- **`event/`** + **`event/listener/`** — Synchronous event bus for system events
- **`config/`** — Configuration structs and loading (Viper-based)
- **`resources/`** — Bundled static assets, admin HTML templates, default theme

### Database

Supports SQLite3 (default), MySQL, and PostgreSQL. GORM with generated query code in `dal/`. Context timeouts are enforced: 5s default, 10s+ for complex queries. See `CONTEXT_TIMEOUT_GUIDE.md` for patterns.

### Dependency Injection

All services, handlers, and repositories are wired via Uber FX in `injection/`. Avoid using `fx.Populate` (global mutable state anti-pattern — was removed). Add new components as FX `fx.Provide` entries.

### Security Patterns

- CSRF: `X-CSRF-Token` header required for all state-changing operations
- Auth: JWT middleware in `handler/middleware/`
- Rate limiting: login endpoint limited to 5 req/min
- Passwords: bcrypt-hashed (including category passwords)
- Path traversal: validate all file paths from user input before file operations

### Storage

Object storage abstraction in `service/storage/`. Supports MinIO, AWS S3, Aliyun OSS, Google Cloud. Local filesystem is the default.

### Configuration

Config files in `conf/config.yaml` (prod) and `conf/config.dev.yaml` (dev). Key settings: server port (default 8080), work directory (logs/DB/uploads/templates), database DSN.

## gstack

Use the `/browse` skill from gstack for all web browsing. Never use `mcp__claude-in-chrome__*` tools.

Available skills: `/office-hours`, `/plan-ceo-review`, `/plan-eng-review`, `/plan-design-review`, `/design-consultation`, `/design-shotgun`, `/design-html`, `/review`, `/ship`, `/land-and-deploy`, `/canary`, `/benchmark`, `/browse`, `/connect-chrome`, `/qa`, `/qa-only`, `/design-review`, `/setup-browser-cookies`, `/setup-deploy`, `/retro`, `/investigate`, `/document-release`, `/codex`, `/cso`, `/autoplan`, `/plan-devex-review`, `/devex-review`, `/careful`, `/freeze`, `/guard`, `/unfreeze`, `/gstack-upgrade`, `/learn`
