/*
Package config handles the centralization and loading of environment-based configuration for the Taskify application.

The Load() function reads environment variables (using a .env file locally as a fallback) and maps them to a
structured Config struct used across the entire application's backend.

Configuration includes:
  - Application Metadata: Port, Environment (dev/prod), App Name.
  - Database: Host, Port, Credentials, DB Name.
  - Redis: Cache and session management connection details.
  - Security: Paths to the RSA public and private keys for JWT signing.
*/
package config
