# Fluxera

Fluxera is a project management backend for organizing projects, tasks, comments, and project activity in one workspace.

It provides a focused API for teams and products that need user authentication, owner-scoped projects, task workflows, task discussions, and an activity feed that tracks important changes across a project.

## Overview

Fluxera is built around a simple workflow:

- users register and sign in;
- authenticated users create projects;
- projects contain tasks with statuses and priorities;
- tasks support comments;
- important project events are recorded in an activity feed.

The backend keeps a clear separation between request handling, business logic, database access, authentication, and infrastructure setup.

## Features

### Authentication

- User registration
- User login
- Password hashing with bcrypt
- JWT-based authentication
- Protected routes
- Current user endpoint

### Projects

- Create projects
- List projects owned by the current user
- Get project details
- Delete projects
- Owner-based access checks

### Tasks

- Create tasks inside projects
- List project tasks
- Filter tasks by status
- Sort tasks by creation or update time
- Update task details
- Change task status
- Delete tasks

### Comments

- Add comments to tasks
- List task comments
- Update own comments
- Delete own comments
- Project owners can remove comments in their projects

### Activity Feed

- Store project activity events
- Track task creation
- Track task status changes
- Track comment creation
- Read project activity feed

## Tech Stack

- Go
- chi
- PostgreSQL
- JWT
- bcrypt
- Docker
- docker-compose

Planned infrastructure additions:

- Kafka
- Redis
- Next.js frontend

## Architecture

Fluxera follows a layered backend structure:

```text
HTTP request
  -> handler
  -> service
  -> repository
  -> PostgreSQL
```

Handlers decode requests and write responses. Services contain business rules, validation, access checks, and activity logging. Repositories isolate SQL and database access.

Authentication is handled through JWT middleware. Protected handlers receive the current user id from the request context.

Project access is owner-scoped. Tasks, comments, and activity feed access are checked through the project that owns the resource.

## Project Structure

```text
fluxera/
  cmd/
    api/
      main.go
  internal/
    auth/
    config/
    db/
    handlers/
    middleware/
    models/
    repositories/
    service/
    storage/
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
```

## Running Locally

Start the local environment:

```bash
docker compose up -d
```

Health check:

```bash
curl http://localhost:8080/healthz
```

## API

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

Supported status values:

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

Fluxera currently stores:

- users
- projects
- tasks
- comments
- activity logs

Activity logs use a JSONB payload to store event-specific details while keeping a consistent event schema.

## License

No license has been specified yet.
