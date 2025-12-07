**Purpose**
- **Scope:** Guidance for AI coding agents working on this small Go monolith (HTTP API + frontend). Focus on making safe, minimal changes and preserving the service layering: transport -> service -> repository -> db.

**Quick Start (developer)**
- **Run locally:** copy environment example and run the app from the project root:
  - `cp env.example .env`
  - `APP_PORT=8080 DATABASE_URL="postgres://postgres:password@localhost:5432/myapp?sslmode=disable" go run .`
- **Build binary:** `go build -o bin/server .`
- **Frontend (production asset build):** `cd frontend && npm install && npm run build` (output expected at `./frontend/dist` and served by `ServeStaticSPA`).
- **Database:** `schema.sql` contains the `http_requests` table used by request logging. The Go code uses GORM for main models and `sqlx` (pgx) for request logging; they use separate DB connections.

**Big-picture architecture**
- **Monolith layout:** files live in the same `package main` and are split by role rather than packages: `internal_config_config.go`, `internal_db_db.go`, `internal_repository_user_repository.go`, `internal_service_user_service.go`, `internal_transport_http_*.go`.
- **Layering (do not invert):**
  - Transport/HTTP (`internal_transport_http_*.go`) depends on Service.
  - Service (`internal_service_user_service.go`) depends on Repository.
  - Repository (`internal_repository_user_repository.go`) depends on GORM DB (`internal_db_db.go`).
  - Request logging uses a separate `sqlx` DB connection in `request_logger.go` (background worker writing to `http_requests`).

**Key files to reference**
- `cmd_server_main.go` — app wiring: `Load()` config, `New()` DB, `AutoMigrate()`, repository/service/handler construction, server start/shutdown.
- `internal_config_config.go` — loads `.env`, default values (app port, log level, default `DATABASE_URL`).
- `internal_db_db.go` — GORM initialization and connection pool settings; callers should call `AutoMigrate` as shown in `main`.
- `internal_repository_user_repository.go` — GORM `User` model and repository implementation. Uses sentinel errors: `ErrNotFound`.
- `internal_service_user_service.go` — business logic, DTOs and request structs, `ErrAlreadyExists` sentinel. Uses `bcrypt` hashing and `go-playground/validator` for input validation.
- `internal_transport_http_handler.go` — Gin routes, validation, how errors map to HTTP status codes (handlers use `JSONError`).
- `request_logger.go` — `sqlx` connection and `StartRequestLogWorker` pattern (buffered channel + non-blocking send) to insert logs into `http_requests` table.

**Project-specific conventions and patterns**
- **Single package `main`:** Files are separated by filename prefixes (e.g., `internal_service_`) rather than Go packages. Avoid scattering new packages without refactor discussion.
- **Factory `New*` functions:** Construction uses `NewGormUserRepository`, `NewUserService`, `NewHandler`, `NewServer`. Follow this style when adding components.
- **Error sentinels used for control flow:** `ErrNotFound` and `ErrAlreadyExists` are used by repository/service layers and mapped to HTTP statuses in handlers — preserve these semantics when changing behavior.
- **DB duality:** Primary app DB uses GORM (`gorm.DB`); request logging uses `sqlx` + `pgx` in `request_logger.go`. Don't replace one with the other without ensuring compatibility and preserving the separate connection and schema (`schema.sql`).
- **Non-blocking logging:** The request logger uses a buffered channel and a non-blocking send; when the buffer is full logs are dropped. Changes should maintain non-blocking behavior or consciously change the trade-off.
- **Validation & error mapping:** Handlers call `validator` on request structs and map specific domain errors to HTTP statuses (e.g., `ErrAlreadyExists` -> 409 Conflict).

**Testing & debug notes (no tests in repo)**
- There are no unit tests in the repo. When adding tests, adhere to the current layering: test services by providing fake `UserRepository` implementations; test handlers with `httptest` using Gin router returned by `Handler.Router()`.
- For debugging startup issues, check `LOG_LEVEL` and `.env` values; `cmd_server_main.go` uses `AutoMigrate` at startup which will modify the DB schema.

**Safe-edit checklist for AI edits**
- Preserve dependency directions: transport -> service -> repository -> db. Don't introduce circular imports.
- If changing DB schema or model tags (`User` struct), update `AutoMigrate` implications and `schema.sql` where relevant (request logs table).
- Keep `NewSQLXDB` and `StartRequestLogWorker` semantics (separate connection and buffered, non-blocking worker) unless explicitly refactoring the logging design.
- Map domain errors to HTTP responses in handlers rather than returning raw errors to clients.
- Use existing helper functions: `JSONError`, `hashPassword`, `clientIP`.

**Common commands & environment**
- Start with `.env` or `env.example`. Default DB DSN used in code: `postgres://postgres:password@localhost:5432/myapp?sslmode=disable`.
- Run server: `go run .` (from project root). Build: `go build -o bin/server .`.
- Build frontend: `cd frontend && npm install && npm run build` -> files served from `./frontend/dist`.

**When to ask the human**
- If you need to introduce package-level refactors (new packages, breaking import layout), ask before changing — this repo intentionally keeps everything in `package main`.
- If changing how logs are persisted (schema or sync behavior) or altering `AutoMigrate` usage, confirm DB migration strategy with the maintainer.

Please review this guidance and tell me which sections you'd like clarified or expanded (examples, more commands, or stricter edit rules).
