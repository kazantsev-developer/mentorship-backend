# Go Mentorship Platform (Backend)

Internal platform for managing the Go mentorship program. It provides a REST API for students, buddies, and administrators to track roadmaps, progress, bonuses, achievements, and 1:1 meetings.

## Prerequisites

- Docker v24.0+
- Docker Compose v2.20+

## Quick Start

1. Clone the repository and navigate to the project root directory:
   git clone https://github.com
   cd mentorship-backend

2. Configure environment variables:
   Create a local .env file based on .env.example:
   cp .env.example .env

3. Start the application:
   docker compose up -d
   The backend service will be available at http://localhost:8080.

## API

### Health Check

Verify the service status using the liveness probe:

- GET /ping — Returns HTTP 200 if the service is running

### Authentication Endpoints

- POST /api/auth/register — Register a new account
- POST /api/auth/login — Authenticate and obtain a JWT Bearer token
- Note: For a complete list of endpoints, see the route definitions in the internal/handlers package.

### Demo Accounts

Register these accounts using the POST /api/auth/register endpoint before authenticating:

- Student:
  { "login": "student1", "password": "123", "roles": ["student"] }

- Buddy (Mentor):
  { "login": "buddy1", "password": "123", "roles": ["buddy"] }

- Admin:
  { "login": "admin1", "password": "admin123", "roles": ["admin"] }

## Tests

Run the automated E2E script to validate core authentication, roadmap tracking, and bonus logic:

chmod +x test_backend.sh
./test_backend.sh

## Teardown

To stop the containers and remove networks, execute:

docker compose down

## Layout

- cmd/api/main.go — Application entry point. Initializes configuration, database, services, and starts the Gin HTTP server.
- internal/config/config.go — Loads environment variables: DB, JWT, port.
- internal/models/ — Domain and GORM models: users, blocks, materials, progress, bonuses.
- internal/repositories/ — Data access layer, one repository per aggregate root.
- internal/services/ — Business logic: progress, bonuses, achievements.
- internal/handlers/ — HTTP handlers, request validation, and JSON responses.
- internal/middleware/auth.go — Gin middleware for JWT validation and role extraction.
- pkg/db/postgres.go — PostgreSQL client pool initialization and automated GORM migrations.
- docker-compose.yml — Orchestration file for local development (PostgreSQL and backend service).
- Dockerfile — Multi-stage Docker build config for production Go binaries.
