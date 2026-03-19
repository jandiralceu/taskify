/*
Package database handles the application's connection to persistent storage.

It focuses on:
  - Initializing a shared connection pool via GORM.
  - Configuring connection settings like SSL mode, port, and credentials.
  - Enabling OpenTelemetry query tracing via the otelgorm plugin.
  - Managing connection life cycles and health checks (ping).
*/
package database
