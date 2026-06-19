# Go Mentorship Platform (Backend)

Go mentorship platform backend. Built with Go, Gin, GORM, and PostgreSQL. Provides a high-performance, concurrent REST API for students, buddies (mentors), and administrators to manage learning roadmaps, training progress, bonus economics, achievements, interviews, and calendars.

## Prerequisites

- Go v1.22+ (for local compilation)
- Docker v24.0+
- Docker Compose v2.20+

## Quick Start

1. Clone the repository and navigate to the project root directory:

   ```bash
   git clone https://github.com
   cd mentorship-backend
   ```

2. Configure infrastructure environment variables:

   ```bash
   cp .env.example .env
   ```

3. Initialize local containers ecosystem:
   ```bash
   docker compose up -d --build
   ```
   The backend service will be available at http://localhost:8080.

## API Architecture

### Health Check

Verify the service network lifecycle using the liveness probe:

- **GET** `/ping` — Returns HTTP 200 `{"message": "pong"}` if the service runtime is healthy

### Authentication Endpoints

- **POST** `/api/auth/register` — Provision a new account profile
- **POST** `/api/auth/login` — Authenticate credentials and obtain a stateful JWT Bearer token

_Note: For the comprehensive routing registry including `/api/admin/_`boundaries, inspect structural definitions inside the`internal/handlers/` package.\*

### Pre-Seeded Local Environment Setup

Execute these explicit payloads via terminal to populate your local PostgreSQL database instance before logging in:

```bash
# Register Admin Profile
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"admin123","display_name":"Platform Administrator","roles":["admin"]}'

# Register Buddy Profile
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"buddy1","password":"123","display_name":"Senior Mentor","roles":["buddy"]}'

# Register Student Profile
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"student1","password":"123","display_name":"Junior Engineer","roles":["student"]}'
```

## Automated Testing Suite

Execute the integration test script suite to validate core authentication state machines, database transaction isolation, and bonus accounting mechanics:

```bash
chmod +x test_backend.sh
./test_backend.sh
```

## Infrastructure Teardown

To gracefully terminate operational runtimes, isolate networks, and flush transient execution layers:

```bash
docker compose down
```

## Repository Layout (Clean Architecture)

- `cmd/api/main.go` — Application entry point. Resolves configuration bindings, initializes connection pools, and mounts the HTTP server engine.
- `internal/config/config.go` — Atomic environment parser mappings (Database pooling settings, JWT credentials secret, GIN network flags).
- `internal/models/` — Domain entity abstractions and relational mappings: users, blocks, materials, progress trackers, bonus logs.
- `internal/repositories/` — Data Access Object (DAO) layer encapsulating database mutation logic.
- `internal/services/` — Pure business domain execution layer isolated from networking transport protocols.
- `internal/handlers/` — Input transport validation, route binding contexts, and output JSON contract formatting.
- `internal/middleware/auth.go` — Token validation boundary intercepting contextual user extraction and role claims authorization.
- `pkg/db/postgres.go` — Connection lifecycle broker establishing thread-safe connection pools and automating migrations.
- `docker-compose.yml` — Local environment orchestration schema linking stateless Go binary runtime and PostgreSQL instances.
- `Dockerfile` — High-efficiency multi-stage container file building stripped down static lightweight scratch binaries.
