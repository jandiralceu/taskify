/*
Package handlers contains the Gin HTTP entry points for the application.

Handlers are responsible for:
  - Extracting path parameters, query strings, and request bodies.
  - Handing over business logic calls to the appropriate Service interfaces.
  - Processing Service errors and mapping them to HTTP status codes.
  - Formatting JSON responses using standardized Problem Details (RFC 7807) for failures.
  - Hosting Swag/Swagger annotations for automated documentation.
*/
package handlers
