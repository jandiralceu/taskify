/*
Package pkg provides a set of generalized infrastructures and utilities shared across the backend.

Utilities included here:
  - Cache: Redis-backed cache manager for sessions, rate limits, and caching.
  - JWT: RSA-signed (RS256) token generation and validation.
  - Logger: Structured logging via slog (JSON for production, tinted for dev).
  - Password: Argon2-id hashing with salt and parameter tuning.
  - Telemetry: OpenTelemetry initialization for tracing and performance monitoring.
*/
package pkg
