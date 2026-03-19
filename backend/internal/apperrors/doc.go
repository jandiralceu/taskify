/*
Package apperrors defines the standard error types and handling logic for the Taskify application.

It centralizes common error definitions such as:
  - ErrNotFound: Resource does not exist.
  - ErrConflict: Resource already exists.
  - ErrUnauthorized: Authentication is missing or invalid.
  - ErrForbidden: Authenticated user lacks sufficient permissions.
  - ErrInternal: Unexpected server-side failure.

Using standardized errors allows the handlers to consistently map internal failures to HTTP status codes 
and JSON Problem Details (RFC 7807) responses.
*/
package apperrors
