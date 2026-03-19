# Taskify API (Backend)

The Taskify backend is a high-performance RESTful API written in Go, focusing on secure authentication, granular RBAC, and efficient task management. It follows the principles of Clean Architecture to ensure separation of concerns between HTTP handlers, business logic Services, and data persistence Repositories.

## ⚙️ Architecture

The backend is architected in several layers:
- **`cmd/api`**: Entry point for the Gin server and dependency injection (wire).
- **`routes`**: Centralized HTTP routing and middleware configuration.
- **`handlers`**: Request parsing, validation, and response formatting (Gin).
- **`service`**: Core business business logic and coordination.
- **`repository`**: Direct data interaction with GORM/PostgreSQL and logic for pagination/filtering.
- **`pkg`**: Shared infrastructure utilities (JWT Manager, Argon2 Hasher, Cache Manager).
- **`models`**: Definitions for core domain entities (User, Task, Role, Permission).
- **`dto`**: Data Transfer Objects for decoupled API request and response structures.

## 🚀 Key Features

### 🛡️ Secure Authentication & Authorization
- **JWT (RS256)**: Public/Private key signing for decoupled verification.
- **Token Rotation**: Seamless refresh token mechanism with Redis-backed session tracking.
- **Casbin RBAC**: Declarative policy-based authorization enforcing resource/action permissions.

### 📊 Optimized Data Access
- **Pagination**: Standardized `PaginationRequest` and `PaginatedResponse` across all list endpoints.
- **GORM Hooks**: Integrated `updated_at` time stamping and preloading associations.
- **Safe Queries**: Input sanitization and parameterized queries protecting against SQL injection.

### 📄 API Documentation
- **Swagger/OpenAPI**: Automated documentation generation via `swag init`.
- Access UI: `http://localhost:8080/swagger/index.html`.

## 🛠️ Installation & Run

1. **Prerequisites**: [Go 1.24](https://go.dev/doc/install), [PostgreSQL](https://www.postgresql.org/), [Redis](https://redis.io/).
2. **Environment**: Copy `.env.example` to `.env` and configure accordingly.
3. **Dependencies**: `go mod download`.
4. **Migrations**: `make migration-up`.
5. **Start**: `make start`.

## 🧪 Testing

- **Unit Tests**: `make test-unit` to run isolated logic tests.
- **Integration Tests**: `make test-integration` for full API and DB interaction.
- **Coverage**: `make test-cover` to generate an HTML coverage report.
