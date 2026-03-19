/*
Package middleware contains the Gin engine plugins that intercept and transform HTTP requests and responses.

Key middlewares include:
  - Auth: Validates JWT access tokens using RSA public keys (asymmetric signature).
  - RBAC: Enforces granular permission checks using the Casbin enforcer.
  - RateLimiting: Protects against brute-force and DDoS via token-bucket algorithms and Redis.
  - Tracing: Injects OpenTelemetry context for end-to-end request visibility.
  - ErrorHandling: Unified Gin-to-app-error translation and response normalization.
*/
package middleware
