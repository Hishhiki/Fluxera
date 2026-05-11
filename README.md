# Fluxera

Fluxera is a lightweight project management system for teams that need a clear place to organize projects, track work, and follow activity across their workspace.

It provides a backend foundation for creating projects, managing task workflows, attaching discussions to work items, and building an activity feed around everything that happens in a project.

## Overview

Fluxera helps users manage work through a simple structure:

- users create accounts and sign in securely;
- authenticated users create and manage projects;
- projects contain tasks with statuses and priorities;
- tasks can have comments and activity history;
- project activity can be processed asynchronously through events.

The backend is designed around clear boundaries between HTTP handlers, business services, repositories, storage, middleware, and infrastructure integrations.

## Features

### Authentication

- User registration
- User login
- Password hashing with bcrypt
- JWT-based authentication
- Protected API routes
- Current user endpoint

### Projects

- Create projects
- List projects owned by the current user
- Get a project by id
- Delete projects
- Owner-based access checks

### Planned Capabilities

- Task management
- Task status updates
- Comments
- Activity feed
- Kafka-driven event processing
- Redis caching
- Worker pools for asynchronous jobs

## Tech Stack

- Go
- chi
- PostgreSQL
- JWT
- bcrypt
- Docker
- docker-compose

Planned infrastructure:

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
  -> storage
  -> PostgreSQL
```

Authentication uses JWT middleware to protect private routes. The middleware validates the token, extracts the current user id, and passes it through the request context to downstream handlers.

Project access is scoped by owner id, so users can only read and mutate their own project resources.

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


Health check:

```bash
curl http://localhost:8080/healthz
```

## API

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
```

Protected routes require:

```http
Authorization: Bearer <token>
```

## License

No license has been specified yet.
