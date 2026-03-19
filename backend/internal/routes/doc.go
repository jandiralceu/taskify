/*
Package routes centralizes the HTTP API definition and routing configuration.

It is responsible for:
  - Initializing the Gin engine with standard middlewares (CORS, Trace, Logger).
  - Registering endpoints for Auth, User, and Task operations.
  - Applying route-level guards for Auth and RBAC (Casbin).
  - Exposing Swagger/OpenAPI documentation and static file serving (for uploads).
*/
package routes
