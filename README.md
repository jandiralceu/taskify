# Taskify 🚀

Taskify is a robust, full-stack Task Management System designed for high-performance teams. Built with a modern Go backend and a Svelte 5 frontend, it features a comprehensive Role-Based Access Control (RBAC) system, secure JWT authentication with token rotation, and integration with PostgreSQL and Redis.

## 🛠️ Technology Stack

### Backend
- **Go 1.24**: Core programming language.
- **Gin Web Framework**: High-performance HTTP routing.
- **GORM**: Object-Relational Mapping (ORM) for PostgreSQL.
- **Casbin**: Powerful Authorization engine with GORM adapter.
- **Argon2**: Secure password hashing.
- **JWT (RS256)**: Asymmetric token signing with RSA keys.
- **Redis**: Fast caching and refresh token session management.
- **OpenTelemetry**: Tracing and performance monitoring.
- **Swagger**: Automated API documentation.

### Frontend
- **Svelte 5 (Runes)**: Next-generation reactive UI framework.
- **TypeScript**: Type-safe development.
- **TanStack Svelte Query**: Efficient server-state management.
- **Vanilla CSS**: Premium, custom-styled interface with modern aesthetics.

### Infrastructure
- **Docker & Docker Compose**: Seamless environment orchestration.
- **PostgreSQL 17**: Relational database for core tasks and user data.

## 🔐 Core Features

- **Granular RBAC**: Dynamic roles (`admin`, `employee`) and permissions based on resources and actions.
- **Secure Authentication**: Sign-in, sign-up, and automated token rotation via refresh tokens.
- **Task Management**: Full CRUD operations for tasks with priorities, status tracking, blockers, and attachments.
- **Clean Architecture**: Decoupled handlers, services, and repositories for maximum maintainability.
- **Performance Optimized**: Pagination, sorting, and caching built directly into the data layer.

## 📈 Getting Started

1. **Prerequisites**: Ensure you have Docker and Docker Compose installed.
2. **Setup Environment**: Copy `.env.example` to `.env` in the `backend` directory.
3. **Launch Project**:
   ```bash
   make docker-up-all
   ```
4. **API Documentation**: Once running, access the Swagger UI at `http://localhost:8080/swagger/index.html`.

## 📜 Project Rubric

This project was developed following the strict guidelines in the [RUBRIC.md](RUBRIC.md), ensuring code quality, security standards (OWASP), and best architectural practices.
