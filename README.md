# Fluxera

Fluxera is an event-driven project management backend for organizing projects, tasks, comments, and activity history inside a focused workspace.

It combines a clean HTTP API with owner-scoped access control, JWT authentication, PostgreSQL persistence, and Kafka-based event processing for project activity.

```text
Projects. Tasks. Comments. Activity.
One backend flow built around explicit events.
```

## What Fluxera Does

Fluxera gives authenticated users a workspace where they can create projects, manage task workflows, discuss work through comments, and track meaningful changes through an activity feed.

The activity feed is not written as a side effect hidden inside handlers. Domain services publish events, Kafka transports them, and consumers persist activity logs. This keeps the request flow, event flow, and persistence logic separated.

## Core Capabilities

### Authentication

- User registration
- User login
- Password hashing with bcrypt
- JWT token generation
- JWT middleware for protected routes
- Current user endpoint

### Projects

- Create projects
- List projects owned by the current user
- Get project details
- Delete projects
- Enforce owner-based access checks

### Tasks

- Create tasks inside projects
- List tasks by project
- Filter tasks by status
- Sort tasks by creation or update time
- Update task details
- Change task status
- Delete tasks

### Comments

- Add comments to tasks
- Read task comments
- Update own comments
- Delete own comments
- Allow project owners to remove comments in their projects

### Activity Feed

- Store project activity logs
- Track project creation
- Track task creation
- Track task updates
- Track task status changes
- Track comment creation
- Read activity by project

## Architecture

Fluxera uses a layered backend structure with a separate event pipeline.

```text
HTTP Request
    |
    v
Handler
    |
    v
Service
    |
    +--------------------+
    |                    |
    v                    v
Repository          Event Publisher
    |                    |
    v                    v
PostgreSQL            Kafka
                         |
                         v
                  Kafka Consumer
                         |
                         v
                 Activity Service
                         |
                         v
                 activity_logs
```

Handlers are responsible for HTTP input and output.

Services contain validation, business rules, ownership checks, and event publishing.

Repositories isolate SQL and database access.

Kafka publishers and consumers move domain events between the application flow and asynchronous activity processing.

## Event Flow

Fluxera represents important domain changes as events:

```text
project.created
task.created
task.updated
task.status_changed
comment.created
```

Each event has a stable envelope:

```json
{
  "id": "bab3808b-e163-4b9a-bb8e-f4b4bd87d93b",
  "type": "task.created",
  "project_id": 10,
  "user_id": 41,
  "payload": {
    "task_id": 7,
    "title": "Check Kafka activity flow"
  },
  "created_at": "2026-05-17T13:23:24Z"
}
```

The service layer publishes events through a small interface:

```go
type Publisher interface {
	Publish(ctx context.Context, event models.Event) error
}
```

That keeps services independent from Kafka-specific code. Kafka is an implementation detail behind the publisher interface.

## Tech Stack

| Area | Technology |
|---|---|
| Language | Go |
| Router | chi |
| Authentication | JWT, bcrypt |
| Database | PostgreSQL |
| Events | Kafka |
| Event UI | Kafka UI |
| Containers | Docker, Docker Compose |
| Cache | Redis planned |
| Frontend | Next.js planned |

## Project Structure

```text
fluxera/
  cmd/
    api/
      main.go
  internal/
    auth/             JWT helpers
    config/           Environment configuration
    db/               SQL schema files
    events/           Event model helpers, Kafka publisher, Kafka consumer
    handlers/         HTTP handlers
    middleware/       Auth middleware
    models/           Domain models
    repositories/     PostgreSQL repositories
    service/          Business logic
    storage/          Database connection setup
  docker-compose.yml
  go.mod
  README.md
```

## Configuration

Fluxera reads configuration from environment variables.

Example `.env`:

```env
HTTP_ADDR=:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5433/fluxera?sslmode=disable
JWT_SECRET=change-me
KAFKA_BROKERS=localhost:9092
```

## Infrastructure

Docker Compose provides the local backing services:

- PostgreSQL
- Kafka
- Kafka UI

Start infrastructure:

```bash
docker compose up -d
```

Kafka UI:

```text
http://localhost:8081
```

Required Kafka topics:

```text
project.created
task.created
task.updated
task.status_changed
comment.created
```

Each topic can be created with:

```text
Partitions: 3
Replication factor: 1
```

## API Overview

Protected routes require:

```http
Authorization: Bearer <token>
```

### Health

```http
GET /healthz
```

### Authentication

```http
POST /auth/register
POST /auth/login
GET /me
```

### Projects

```http
POST /projects
GET /projects
GET /projects/{id}
DELETE /projects/{id}
GET /projects/{id}/activity
```

### Tasks

```http
POST /projects/{projectID}/tasks
GET /projects/{projectID}/tasks
PUT /tasks/{id}
PATCH /tasks/{id}/status
DELETE /tasks/{id}
```

Task list query parameters:

```http
GET /projects/{projectID}/tasks?status=todo&sort=created_at_desc
```

Supported statuses:

```text
todo
in_progress
done
```

Supported sort values:

```text
created_at_desc
created_at_asc
updated_at_desc
updated_at_asc
```

### Comments

```http
POST /tasks/{id}/comments
GET /tasks/{id}/comments
PUT /comments/{id}
DELETE /comments/{id}
```

## Data Model

Fluxera stores the core workspace data in PostgreSQL:

| Table | Purpose |
|---|---|
| `users` | Accounts and password hashes |
| `projects` | Owner-scoped project containers |
| `tasks` | Work items inside projects |
| `comments` | Task discussion messages |
| `activity_logs` | Event-backed activity feed |

Activity logs use `JSONB` payloads so each event can keep event-specific details while sharing one consistent activity table.

## Reliability

Fluxera includes:

- request logging middleware
- panic recovery middleware
- JWT-protected routes
- ownership checks for project-scoped data
- Kafka consumer retry on temporary fetch errors
- graceful shutdown for HTTP server and Kafka consumers

## Current Direction

Fluxera is ready for the next infrastructure layer: Redis.

Redis will be used for:

- activity feed cache
- project cache
- task list cache
- TTL-based cache expiration
- cache invalidation after Kafka events

## License

No license has been specified yet.
