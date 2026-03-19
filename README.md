# Taskify

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?logo=svelte&logoColor=white)](https://svelte.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?logo=kubernetes&logoColor=white)](https://kubernetes.io/)

A full-stack, production-ready Task Management & Kanban application. Taskify features a highly reactive frontend built with Svelte 5 and a robust Go backend designed with clean architecture principles. It includes multi-role RBAC, distributed tracing, and automated deployment manifests.

## Overview

Taskify is a comprehensive demonstration of modern web development, combining high-performance backend logic with a fluid, state-driven user interface:

- **Full-stack Task Management** — Kanban-style task board with drag-and-drop status updates, priority indicators, and detailed task views.
- **Authentication & RBAC** — Secure JWT authentication using RSA key pairs and Argon2id password hashing. Role-Based Access Control via [Casbin](https://casbin.org/) supporting `admin` and `employee` roles.
- **Interactive Features** — Task comments (notes), file attachments, and user search with real-time feedback.
- **Advanced State Management** — Frontend powered by Svelte 5 Runes for reactivity and TanStack Query for efficient server-state synchronization and caching.
- **Observability** — End-to-end distributed tracing using OpenTelemetry (OTLP) with Jaeger, and structured logging via `slog`.
- **Database Migrations** — Version-controlled schema management with `golang-migrate`, ensuring consistency across environments.
- **Deployment Ready** — Multi-stage Docker builds and complete Kubernetes manifests (Deployments, Services, ConfigMaps, Secrets) for scalable orchestration.
- **Premium UI/UX** — Modern design utilizing Skeleton UI, Lucide icons, and Tailwind CSS for a professional, responsive experience.

## Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | Go 1.24, Gin, GORM |
| **Frontend** | Svelte 5 (Runes), Vite, TanStack Query |
| **Styling** | Tailwind CSS, Skeleton UI |
| **Database** | PostgreSQL 17 |
| **Cache** | Redis 7 |
| **Auth & Security** | JWT (RSA), Argon2id, Casbin (RBAC) |
| **Tracing** | OpenTelemetry + Jaeger |
| **Documentation** | Swagger (swag) |
| **Infrastructure** | Docker, Docker Compose, Kubernetes |

## Architecture

The project maintains a clear separation between the presentation layer and the core business logic:

```
┌─────────────────────────────────────────────────────┐
│                  Frontend (Svelte)                  │
│   Reactive UI, TanStack Query, Skeleton Components   │
├─────────────────────────────────────────────────────┤
│                    API Gateway                      │
│   Middleware (Auth, RBAC, Tracing, Rate Limiting)   │
├─────────────────────────────────────────────────────┤
│                  Backend (Go API)                   │
│   Handlers, Services, Repositories, Domain Models    │
├─────────────────────────────────────────────────────┤
│                 Infrastructure Layer                │
│   PostgreSQL, Redis, OpenTelemetry, S3-style Uploads │
└─────────────────────────────────────────────────────┘
```

## Project Structure

```
.
├── backend/                   # Go API source code
│   ├── cmd/api/               # Server entrypoint
│   ├── internal/              # Core logic (Handlers, Services, Repositories)
│   ├── migrations/            # SQL migration files
│   └── Dockerfile             # Backend production build
├── frontend/                  # Svelte frontend application
│   ├── src/                   # Application source (Routes, Components, State)
│   ├── static/                # Static assets
│   └── Dockerfile             # Frontend production build (Nginx)
├── deployment/                # Deployment-related configurations
│   ├── compose.yaml           # Docker Compose for local full-stack
│   └── k8s/                   # Kubernetes manifests
│       ├── base/              # API & Web base manifests
│       └── local/             # Infrastructure (PostgreSQL, Redis, Jaeger)
├── Makefile                   # Developer workflow automation
├── model.conf                 # Casbin RBAC model
└── policy.csv                 # Casbin RBAC initial policies
```

## Environment Variables

### Backend (.env)
| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_USER` | `taskuser` | PostgreSQL username |
| `DB_PASSWORD` | — | PostgreSQL password |
| `DB_NAME` | `taskify` | PostgreSQL database name |
| `APP_PORT` | `8080` | Backend API port |
| `REDIS_HOST` | `localhost` | Redis host |
| `OTLP_ENDPOINT`| `localhost:4317` | OpenTelemetry collector endpoint |

### Frontend (.env)
| Variable | Default | Description |
|---|---|---|
| `PUBLIC_API_URL`| `http://localhost:8080` | Backend API base URL |

## Getting Started

### Option 1 — Local Development with Docker Compose

Ensure you have Docker and Docker Compose installed.

```bash
# 1. Clone the repository
git clone https://github.com/jandiralceu/taskify.git
cd taskify

# 2. Generate RSA keys for JWT
# (Run in backend directory if using local scripts, or use the project Makefile)
make generate-keys

# 3. Start the entire stack
docker compose -f deployment/compose.yaml up --build
```

The app will be available at `http://localhost:3000` and the API at `http://localhost:8080`.

### Option 2 — Manual Setup

1.  **Backend:**
    *   Initialize PostgreSQL and Redis.
    *   Run `make migration-up` in the `backend/` folder.
    *   Start with `go run cmd/api/main.go`.
2.  **Frontend:**
    *   `cd frontend && npm install`
    *   `npm run dev`

## Kubernetes Deployment

Production-ready manifests are available in `deployment/k8s/`.

```bash
# 1. Create the namespace
kubectl apply -f deployment/k8s/base/taskify-namespace.yaml

# 2. Deploy local infrastructure (Postgres, Redis, Jaeger)
kubectl apply -f deployment/k8s/local/

# 3. Deploy the application (API & Web)
kubectl apply -f deployment/k8s/base/
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
